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

Prefer ordinary Python workers that call the API, update the cache and evaluate
conditions without invoking an LLM. Reuse the connection chosen during Access:

| Selected connection | How a worker calls Notehub |
|---|---|
| CLI | Python invokes `notehub api` in the same authorized OS environment and parses its response. The CLI supplies authentication and handles supported credential renewal. |
| MCP | Use the harness's documented programmatic bridge to the existing authenticated connection. An interactive MCP tool alone does not give a Python process this access. |

Identify the actual invocation and verify a harmless read from the worker environment
when available; otherwise verify it during the first task run. If only scheduled agent
turns can invoke the connection, record that fallback and its model cost before building
the workers. Keep credentials in the selected tool and retain one authenticated path.

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
If the harness only schedules agent turns, use that supported path and state the cost;
do not describe it as LLM-free polling.

Return to the cookbook with the script/runtime location, available storage and warning
capabilities, authenticated invocation, scheduling commands, restart mechanism and
agent wake-up/delivery path identified.
Record gaps and pending checks. The following steps implement and verify these paths
as they create the cache and tasks; no separate configuration hierarchy is required.
