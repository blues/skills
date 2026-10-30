# Establish how background work runs

Before creating the cache, determine how this Cloud Computer will run durable,
efficient background work. Use the harness's actual tools and documentation, and
reuse any working setup. Leave this step with concrete paths, commands and hooks.

## 1. Locate persistent code and its runtime

Choose where worker scripts and their dependencies live so scheduled jobs can find
them after a restart. Use the standard shared workspace when available, independent
of any one bot. Identify the Python interpreter or other supported runtime, how to
install dependencies, and the exact command and working directory a job will use.
Distinguish durable storage from temporary execution space. If the runtime is recreated
on restart, determine how its dependencies are restored before jobs resume.

## 2. Check storage capacity and low-space warnings

Check the capacity and currently available space in the durable workspace and temporary
storage, including any account/workspace quotas. Determine whether other bots share
those limits and how to inspect usage. Record the actual measurements and any limits
the harness does not expose.

Identify the harness's low-space or quota-warning mechanism, its warning thresholds,
and how a warning can wake the agent to investigate and free space. Check whether that
notification path still works when local storage cannot accept writes. If there is no
native warning, identify how a lightweight periodic check could detect low space and
trigger the agent. This step verifies available capabilities; the later Cache
Maintenance Task incorporates storage monitoring and the low-space response.

## 3. Establish programmatic Notehub access

Local database work (ingestion, purge, schema review, SQL evaluation, storage checks)
always runs in ordinary Python workers, without invoking an LLM. Notehub fetches reuse
the connection chosen during Access:

| Selected connection | How a worker calls Notehub |
|---|---|
| CLI | Python invokes `notehub api` in the same authorized OS environment and parses its response. The CLI supplies authentication and handles supported credential renewal. |
| MCP | **Ideal:** a simple, script-callable Python-to-MCP programmatic bridge that lets a script work through the MCP server without the model's involvement. Use it whenever the harness provides one. **If not:** use the **agent-turn fetch handoff** below. An interactive MCP tool alone does not give a Python process access. |

Identify the actual invocation and verify a harmless read from the worker environment
when available; otherwise verify it during the first task run. Keep credentials in the
selected tool and retain one authenticated path.

### Ideal: a script-callable MCP bridge

If a simple, script-callable Python-to-MCP programmatic bridge is available, it is ideal:
it lets a script do its work through the MCP server, sharing the existing authenticated
session, without the model's involvement. Prefer it whenever the harness provides it. The
worker then fetches, ingests and evaluates entirely in code, and no agent turn runs on
routine polls.

### If not: agent-turn fetch handoff

If no such bridge exists, this is the required acquisition path. Build
it without asking the user, and do not ask them to weigh the model cost. Do not offer
or request a second credential (personal access token, CLI login, environment token),
and do not substitute ad-hoc workarounds such as internal harness sockets or copied
session credentials.

- A Python handoff worker owns the database. It exposes `plan` (the due fetches with
  their saved cursors), `ingest` (validate one fetched batch and commit its records and
  next cursor in one transaction, rejecting superseded cursors), `enroll` (add a
  ProjectUID to the registry) and `status` (coverage and freshness).
- A scheduled agent turn at the refresh interval runs `plan` and stops quietly when it
  is empty. For each due item, it fetches one batch with a read-only MCP call. It saves
  the raw response, wrapped in an envelope (`project_uid`, `stream`, `device_uid` for
  sessions, `requested_cursor`, `fetched_at`), to a file, then runs `ingest` on that
  file. It repeats while `has_more` is true, within a per-run batch budget.
- Pass responses through files, not through the turn's own text, so batch size does not
  drive model usage.
- The pure-Python scheduler keeps running storage checks, purge and optimization. It
  never fetches from Notehub itself.

Record this path as the selected acquisition method in the task record and the cache
README. Switch to direct worker calls when a programmatic bridge becomes available.

#### Writing the handoff worker

Write the handoff worker locally on the Cloud Computer, next to the database. The skills
repository does not ship code. The examples below show the contract; adapt the table
layout to the [Cache model](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-cloud-agent/cache.md).

**Envelope.** Each fetch turn returns exactly this object from one read-only MCP block
and saves it to a file:

```json
{"project_uid": "app:...", "stream": "devices|events|sessions|schemas",
 "device_uid": "dev:... (sessions only)", "requested_cursor": "0",
 "fetched_at": "2026-10-30T14:00:00Z", "response": {"devices": [], "has_more": false, "next_cursor": "..."}}
```

**Fetch block** (runs inside the MCP code-mode `execute` tool, not as a local script):

```python
import datetime
p = "app:..."; cur = "0"   # values taken from `plan`
r = await call_tool("get_events", {"projectOrProductUID": p, "cursor": cur, "pageSize": 500})
return {"project_uid": p, "stream": "events", "requested_cursor": cur,
        "fetched_at": datetime.datetime.now(datetime.timezone.utc).isoformat(), "response": r}
```

**Commands.**

| Command | Does |
|---|---|
| `plan` | For each active enrolled project, list devices and events with their saved cursors (`"0"` if none), sessions for each known device, and schemas when a daily refresh is due |
| `ingest FILE...` | Parse each envelope. Ignore any text wrapper around the JSON. Validate and commit it |
| `enroll APP_UID LABEL` | Add the project to the registry; the next fetch turn performs its first fill |
| `status` | Per project: counts, last successful ingest and cursor positions |

**Ingest rules**, the core of the worker:

```python
def ingest(con, env):
    p, s, dev = env["project_uid"], env["stream"], env.get("device_uid")
    if not enrolled(con, p):
        return "rejected:not-enrolled"
    key = s if s != "sessions" else f"sessions:{dev}"
    saved = saved_cursor(con, p, key) or "0"
    if s != "schemas" and str(env["requested_cursor"]) != saved:
        return "rejected:superseded"          # a newer batch already committed
    r = env["response"]
    rows = r if s == "schemas" else r.get(s, [])
    nxt, more = (None, False) if s == "schemas" else (r.get("next_cursor"), bool(r.get("has_more")))
    if more and (not nxt or nxt == saved):
        return "rejected:bad-cursor"          # has_more without progress
    con.execute("BEGIN")
    try:
        upsert(con, p, s, rows, env["fetched_at"])   # key by uid / event / session / notefile
        if nxt is not None:
            save_cursor(con, p, key, nxt, env["fetched_at"])
        log_ingest(con, env, nxt, more, len(rows))
        con.execute("COMMIT")
    except Exception:
        con.execute("ROLLBACK"); raise
    return f"ok:{len(rows)} has_more={more}"
```

Keep cursors as opaque strings. Use the actual fetch time, never the ingest time.
Serialize writers with a file lock. Test with synthetic envelopes in a disposable
database first: a replayed envelope must be rejected, and a new device must make
`plan` add its sessions stream.

## 4. Determine scheduling and restart recovery

Identify how to register, list, invoke, update and cancel periodic jobs. Determine the
supported cadence, timezone, runtime limits, run-status reporting and overlap behavior.
Prefer a scheduler that launches the worker directly rather than starting an LLM turn
for every poll.

Check what survives a Cloud Computer restart. If schedules persist, verify how they
resume. If registrations are lost, identify a supported startup hook that can restore
them from durable definitions. Restoration must reuse matching jobs, register missing
active jobs and preserve pauses/cancellations. The recovery hook must work even when
all periodic registrations are gone. Record any missing automatic recovery capability.

## 5. Determine how workers wake the agent

Find the harness's supported way for code to enqueue or trigger an LLM task in the
right user/project context, and its notification mechanism for the user. Ordinary
polls, SQL calculations and comparisons stay in code. Wake the LLM when a finding
needs reasoning, a message or graph, or a question for the user. An unchanged routine
check should finish without a model turn.

A worker saves a pending finding before requesting that wake-up. Pass a compact job or
finding ID and evidence location so the agent can recover the context. Identify how
failed wake-ups are retried without losing the finding or duplicating user messages.
If the harness only schedules agent turns, use that path for the wake-up without asking
the user. Record it as agent-turn scheduling, not LLM-free polling.

Return to the cookbook with the script/runtime location, available storage and warning
capabilities, authenticated invocation, scheduling commands, restart mechanism and
agent wake-up/delivery path identified.
Record gaps and pending checks. The following steps implement and verify these paths
as they create the cache and tasks; no separate configuration hierarchy is required.


### One call_tool per execute (MCP fetch handoff)

Each agent-turn fetch must run **one** Notehub MCP `execute` containing **exactly one**
`call_tool` for the planned batch (`pageSize` per cache/modified-cursor rules). Do not
chain multiple fetches in one execute. Persist the envelope to a file, run local
`ingest`, then issue another execute for the next batch.

## Grok Bot only

*This section applies only to agents running as Grok Bot on its Cloud Computer (the "box").
Other harnesses should skip it.*

- **No MCP bridge yet.** Box scripts cannot call connected MCP servers, so use the
  agent-turn fetch handoff. Do not probe internal agent sockets or local service ports
  as a substitute.
- **Storage and runtime.** Use a `notehub-cache/` folder in the shared workspace, which
  every bot of the user can reach. Put DuckDB in a shared Python virtual environment.
- **Local scheduler.** The box has no cron or systemd. Run the pure-Python worker from a
  `nohup setsid` loop script with a lock. It does not survive a box restart, so record
  that limitation and the restart command in the README.
- **Agent-turn fetch.** Create one routine, for example "Notehub cache refresh", every 15
  minutes. Its prompt runs `handoff.py plan`. It calls the Notehub MCP's read-only
  `execute` with `write_to_file` set, so the envelope lands in the agent-tools folder
  rather than in the conversation. Then it runs `handoff.py ingest <file>`, deletes the
  file, and stays silent on success. Name no credential in the prompt.
- **Discovery.** Save a shared memory fact pointing every bot to the cache folder's README.
- **Wake-up and delivery.** Box scripts cannot start an agent turn. Workers write findings
  to a pending folder, and a scheduled routine delivers them (see Monitor preparation,
  Grok Bot only).
- **Creating the routines.** Grok Bot creates scheduled agent turns as routines. If you
  cannot create routines yourself (for example, you are a subagent without the routine
  tool), do not skip the step. Hand the exact routine name, schedule and prompt to the
  agent that can create it. The same applies to monitor and engagement routines.
  Routines are Grok Bot's mechanism; other harnesses use their own scheduler and should
  not assume this tool exists.
