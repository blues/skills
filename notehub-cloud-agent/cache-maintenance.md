# Create or verify the Notehub Cache Maintenance Task

This is Cloud Computer cookbook Step 3. Use the verified database and the execution,
scheduling, restart and wake-up mechanisms already established. Complete missing
work and preserve working schedules, enrolled ProjectUIDs, pauses and checkpoints.

## 1. Find the existing task and recover its definition

Inspect the database's task records and the actual scheduler. Reuse a matching task;
shared bots use the same task and project registry. If it is already working, retain its
invocation and receipts and verify its current health in Step 4. For a new task, use the
durable worker location and selected authenticated connection from harness preparation.

**Complete when:** there is one identified task definition for this cache, existing or new.

## 2. Implement or verify its worker

Aim to run one task every 15 minutes. If the background scheduler cannot meet that
goal, use its fastest supported interval so changes can reach the user promptly.
Each run checks storage, fills due projects and runs daily work when due. Read the
live registry each time so newly enrolled projects join automatically. Preserve
pauses/exclusions. With no projects yet, check storage and skip
acquisition. Use Python or the supported runtime for routine work; wake the agent only
when a finding or failure needs reasoning or user interaction.

Implement the phases below, or verify the existing worker performs them. The resources
contain algorithms and data rules; return to this procedure after the needed section.

| Phase | Implementation resource | Verify |
|---|---|---|
| Execution inputs and durable state | [Worker context](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-cloud-agent/jobs.md#common-execution-context) | Exact invocation, database locator, scope and restartable progress |
| Storage monitoring on every run | [Storage monitoring](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-cloud-agent/jobs.md#storage-monitoring) | Free space/quota checks, established housekeeping and a working low-space wake-up path |
| Fill and refresh | [Acquisition](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-cloud-agent/jobs.md#acquisition) | Devices, Events, full per-device Sessions, metadata and schema snapshots, with atomic cursor progress |
| Daily purge, schema refresh and optimization | [Daily maintenance](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-cloud-agent/jobs.md#maintenance) | Retention, daily schema checks and improvements based on fields actually queried |

Set per-project work budgets, timeouts, retry/backoff and overlap control. One large or
failing project must not starve others. Apply native storage warnings alongside periodic
checks. If separate scheduled workers are needed, retain their IDs and responsibilities
under this maintenance task and use the same registry.

When fetches run as agent turns (no script-callable MCP bridge), follow the fast-loop
rules in [Harness: one call_tool per execute](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-cloud-agent/harness.md#one-call_tool-per-execute-mcp-fetch-handoff):
pageSize 10000, two tool calls per page, per-phase timing in a CSV, and a per-run page
budget large enough for backfills (hundreds of pages) rather than a small default cap.

**Complete when:** the worker implements all phases and resumes from durable state.

## 3. Register or verify scheduling and restart recovery

Register the exact worker invocation at the chosen cadence and timezone, or verify the
existing registration. Local phases run as direct code execution. Ideally, fetches do
too, through a script-callable MCP bridge. Where no bridge exists and fetches use the
agent-turn fetch handoff, register that turn too, at the same cadence, without
asking the user about its cost. Use shared/user-level ownership
where supported so deleting one bot does not remove the user's maintenance task.

Apply the restart mechanism established in harness preparation: persistent schedules,
or a supported startup hook that reconciles durable definitions. Preserve existing
jobs and pauses, restore missing active jobs, and avoid duplicate registrations.

| For restart recovery | Resource | Then |
|---|---|---|
| Implement and verify restoration | [Schedule recovery](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-cloud-agent/lifecycle.md#restore-schedules-after-a-restart) | Test recovery idempotently with a disposable task |

**Complete when:** the scheduler reports the task and its next run, with a verified recovery path or an explicit limitation.

## 4. Exercise the actual background path

Run the registered task through the path it will use unattended. Verify authenticated
access, database access/locking, registry reads, bounded execution and saved run results.
Exercise the acquisition algorithm with disposable batches and the storage warning with
simulated low space, including failed local writes. Verify worker-to-agent handoff without
sending an unsolicited notification. Reuse current verification evidence when valid.

With no project selected, a successful storage check and idle acquisition are expected.
Verify real fill and daily review as Project Selection enrolls the first project. Keep
runtime readiness separate from each project's coverage; initial history can fill in
bounded runs while the agent answers supported questions or continues Training.

**Complete when:** actual runs, next execution and remaining gaps are recorded; a saved script alone is insufficient.

## 5. Save the result and return

Store the task identity, invocation, cadence/timezone, restart mechanism and verification
receipts. Keep storage, acquisition, schema refresh, purge and optimization status
separate. Update the cache README with the responsible routine and observed freshness;
verify the shared memory fact still points to the folder.

Return to the cookbook with the database and task ready, or the precise unsupported
capability. Project Selection later adds projects to this live registry; it does not
create another cache or task.
