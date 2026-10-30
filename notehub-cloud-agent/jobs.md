# Jobs that work without conversation memory

Cloud Computer setup creates the Notehub Cache Maintenance Task using the acquisition and daily review
contracts below; first-use project enrollment supplies the scopes it maintains. Implement them
with the harness's actual worker/scheduler tools and database state.
Use deterministic code and SQL for routine work, in a language supported by the harness.
Local work runs in directly scheduled Python workers. Ideally, Notehub fetches run in the
same workers through a script-callable Python-to-MCP bridge, with no model involvement.
If no bridge exists, LLM turns are used for Notehub fetches (see the agent-turn fetch handoff in
[Background execution](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-cloud-agent/harness.md)) and for findings
that need reasoning or user interaction.
The pseudocode below describes behavior; its names are not callable APIs or a prescribed
database schema. Implement and validate it using the selected authenticated transport,
database and scheduler.

This is an implementation resource, not another setup workflow. Implement or verify
the sections named by the calling cookbook step, then return there. The pseudocode
and verification cases below remain the contracts for the workers you create.

## Common execution context

Store the worker's actual executable/script path or tool invocation, database
locator, ProjectUID, monitor ID and configuration version in database job records.
Give the harness enough bootstrap information to open the database safely and find
that job record. The worker reads configuration and product revisions from the
database, verifies the selected account and ProjectUID, obtains the appropriate
lock, performs bounded work, commits state and returns a small machine-readable run
result. No credential is placed in a command, task prompt, script or database
record.

The Notehub Cache Maintenance Task reads the current enrolled-project registry on
every run, including enrollments made after it was created. Process active scopes
only; preserve pauses/exclusions and isolate failures by project. With an empty registry,
run the storage check and skip acquisition. Keep per-project budgets and continuation state so a large backlog does
not monopolize the service. Aim for a 15-minute refresh interval; if the background
scheduler cannot meet that goal, use its fastest supported interval. Keep per-project
due times for daily review. If the implementation uses separate scheduled subtasks,
record them under the maintenance task, reconcile the same registry and avoid duplicate work.
Bots sharing the workspace reuse this task and its registry. After each cache run, refresh
the folder README's timestamped freshness summary from committed state, including partial
or failed refreshes. Update its table/routine descriptions when they change, and update
the shared memory locator if the cache folder moves.

If the scheduler launches agent turns instead of processes, its task instruction should
specify the database locator, safe worker invocation, job ID/type and authorized scope,
then direct the turn to load its database records and run it. Return quietly on an unchanged successful result;
invoke reasoning only for a meaningful candidate or unresolved failure. Ensure that task
has both the chosen connection and database/worker access. Don't paste an entire manual or rely
on chat history in every scheduled prompt.

A run summary records job ID, start/finish, status (`complete`, `partial`, `failed`,
`paused`), ProjectUID and monitor ID when applicable, counts, evidence
freshness/coverage and continuation need. Include error class/action, not secrets or
raw customer dumps. A success marker written before all required steps complete is
not a successful run. Persist registration and execution receipts separately.

## Storage monitoring

Storage monitoring is an implicit part of the Notehub Cache Maintenance Task, requiring
no separate user monitor. On every run, before acquisition or other substantial writes,
check free space and applicable quotas for durable and temporary storage. Use the
mechanisms identified during harness verification. Set warning reserves for expected
growth and temporary working space, accounting for other bots sharing the quota.

Normal checks run in code without waking an LLM. On low space, apply established
housekeeping within its retention rules and verify the space actually recovered.
Preserve checkpoints, pending findings/deliveries, protected history and unrelated files.
If space remains low or an allocation fails, defer further growth and wake the agent
with the measurements and cleanup results so it can investigate and free space.
Deduplicate repeated warnings while surfacing material worsening or recovery.

If local writes are already failing, use the harness warning/wake-up path identified
in Step 1; reporting a full disk must not depend on writing to that same disk. Test
the response with simulated low-space readings, including unavailable local writes.
Keep any native harness storage warnings enabled as an additional signal between runs.

## Wake the agent only for work that needs it

Use the worker-to-agent mechanism identified during harness preparation. The ordinary
path is code calling the authenticated API, updating the cache and evaluating queries.
Only a finding needing interpretation, composition or user input wakes an LLM task.
If code cannot wake the agent or reach the delivery channel, a scheduled agent turn reads
the pending findings and delivers them. This fallback is required, not a choice to put to
the user.
Already prepared routine deliveries can use the delivery worker without another model turn.

```text
run due acquisition and query work in the worker
save new findings needing reasoning or user interaction with stable finding IDs and evidence
for each due pending finding, including retries from earlier runs:
    request an agent task through the harness's supported wake-up mechanism
    pass the finding ID and context to locate its project, meaning and evidence
    on wake-up failure: retain pending work and retry with the same finding ID
if no findings need attention: finish quietly
```

The agent claims the pending work, rereads its current state and evidence, then
reasons, composes or asks the user as needed. Repeated wake-ups must not process the
same finding concurrently or create new delivery IDs. Scheduling a turn does not
mark a finding handled or a message delivered; use the existing durable evaluation
and delivery state. Test a quiet result, a useful candidate, a failed/repeated
wake-up and recovery of pending work using a disposable case and no external
notification.

## Acquisition

| Needed for this worker | Resource | Apply |
|---|---|---|
| Stream definitions and initial coverage | [Cache model, Step 3](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-cloud-agent/cache.md#3-define-streams-and-verify-initial-coverage) | ProjectUID, entity keys and required history |
| Poll for changes in Events, Sessions, Devices, and more | [Modified cursors](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-project-chat/modified-cursors.md) | Exact endpoint cursor, batching and deletion contract |
| Refresh product knowledge used by queries | [Read product skills](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-access/project-skills.md) | Reuse unchanged knowledge and fetch relevant changes |

The maintenance task's fill phase runs at the configured refresh interval across
enrolled projects. For each due project, serialize writes per database, keep stream
definitions stable, and budget work across streams/devices so a large Event backlog
cannot starve other entities or projects.

1. Validate the account and ProjectUID. A mismatch stops work for that Project; do not
   reuse its cursor under another account or ProjectUID. A network failure does not
   initiate repeated login.
2. Poll Devices through the saved modified cursor using transactional batch ingestion.
   Enroll new DeviceUIDs into Session streams. Preserve historical records when
   a Device is deleted and track historical-device coverage gaps separately.
3. Poll Events and each due full-Session stream with their own saved exact cursors.
   Commit each batch with its checkpoint. A batch budget leaves a resumable partial stream.
   Apply only the entity's documented deletion markers; equal values are not duplicates.
4. On first fill, fetch project metadata and the complete Notefile Schema snapshot to
   staging. Validate and commit it before marking that part ready. Daily maintenance
   refreshes these thereafter; a new field need can trigger an earlier refresh. A failed
   or partial response retains the previous snapshot and freshness. An authoritative
   empty snapshot is different from an error. Schemas do not acquire an invented cursor.
5. On first acquisition, fetch all product skills with `GetSkills` and
   `includeContent=true`. For later refreshes, compare metadata and reuse unchanged
   documents; fetch all new or changed content individually, or get a full snapshot.
   Reconcile the private copy only after validation; mark affected query plans/views/
   monitors for review instead of silently running old semantics. This is a read workflow.
6. Publish coverage/freshness from committed state. Update typed projections only under
   the meaning revision they implement. Return per-stream progress and continuation needs.

Use Cache to establish bootstrap and retention rules. `has_more=false` proves
completion of a scan, not capture-time coverage or all-device Session completeness. Each
stream can succeed or fail independently; a project summary retains that distinction.
Use bounded backoff and scheduler continuation rather than a busy loop. A healthy poller
and a healthy product are separate conclusions.

### Batch ingestion pseudocode

Apply this to each due modified-cursor stream. Use the database access design chosen
in Cache, including serialization of overlapping work. Each attempt belongs to the
verified account, ProjectUID and saved stream definition. Use DeviceUID for Device
upserts, EventUID for Event upserts and SessionUID for Session upserts, retaining
the associated DeviceUID for per-device Session streams. These are the record keys
below.

```text
cursor = saved exact cursor, or the stream's fixed bootstrap cursor
seen = {cursor}
while this stream has time and batch budget:
    claim the current attempt for this stream and cursor; mark the scan in progress
    fetch a batch through the selected authenticated connection
    retain the actual fetch-completion time
    validate records, record keys, deletion markers and cursor metadata
    reject a repeated/cyclic cursor, except an unchanged cursor with has_more=false

    in one database transaction:
        verify this attempt and its scope/cursor are still current
        bulk upsert records and documented tombstones by ProjectUID and record key
        save the exact next_cursor and actual fetch time
        mark the scan complete only if has_more is false
        retire this attempt
    on fetch, validation or transaction failure:
        roll back uncommitted work; keep prior committed rows and cursor
        if this attempt is still current:
            mark incomplete, record the failure and arrange a bounded retry
        stop without changing a newer attempt's state

    cursor = next_cursor
    if has_more is false: return complete
    add cursor to seen
retain the committed cursor, mark incomplete and arrange continuation; return partial
```

Normalize the endpoint's documented body or header metadata without inventing
missing values. Reject malformed records and conflicting duplicate keys within a
batch; equal measurements with different EventUIDs remain separate. An empty final
batch may keep the same cursor or advance it; commit its successful freshness either
way. A failed or unfinished scan retains the last complete success and its coverage
limits.

When only scheduled agent turns can invoke MCP, the turn performs the fetch through
that same connection and passes the response to its database worker as a file, using
the agent-turn fetch handoff in Background execution. This is the required path, not
an option to put to the user. Carry the ProjectUID, stream definition, requested
cursor, current attempt ID and actual fetch time with the response. Reject
superseded results and late failures, including when a newer empty scan kept the
same cursor. Never relabel an old response with ingestion time. Resolve abandoned
attempts after the worker timeout. Choose a direct bridge when supported; do not
assume MCP tools are local imports or add a CLI login to work around that
limitation.

### Validate the implemented acquisition worker

Use synthetic batches and a disposable database in the actual worker runtime before
activation. These cases test behavior without customer credentials:

| Case | Expected result |
|---|---|
| Updates, documented tombstones and equal values under different EventUIDs | Updates replace by the record key, tombstones apply, distinct Events survive |
| Empty final batch, including an unchanged cursor | Successful scan and freshness are committed; no extra empty request |
| Missing metadata, conflicting duplicate keys, cursor cycle or failed transaction | Failing batch advances nothing; earlier commits survive |
| Budget exhaustion, process restart or changed project/filter scope | Resume exactly from the committed cursor; a different scope uses its own checkpoint |
| Late MCP result/failure or replay after a newer completion | No overwrite, freshness regression or degradation of the newer success |

## Evaluation

Run each due monitor using its prepared calculation, meaning revision and preferences.
Use the cache when its scope, coverage and freshness suffice; otherwise acquire bounded
fresh inputs through the selected connection. Record fields and query performance during
execution. A faster evaluator cannot make stale evidence fresh.

| For this routine | Resource | Apply |
|---|---|---|
| Calculate the question | [Prepared questions](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-project-chat/prepared-questions.md) | Reuse the verified meaning, parameters and evidence requirements |
| Define what to watch and when and where to notify | [Monitor preferences](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-project-chat/monitors.md#2-set-the-monitoring-preferences) | Use the same definition for requested and interest-driven monitors |

### Evaluation and delivery pseudocode

Implement the following behavior with the harness's runtime and scheduler. This is one
monitoring routine; the harness may split execution where needed. Keep concurrent runs
from evaluating or sending the same work twice.

```text
load the monitor's current definition, preferences and saved progress
if paused or removed: stop

if a check or periodic notification is due:
    evaluate the question; record evidence time, coverage and meaning revision
    compare valid results with the last valid result using the agreed rule
    on the first valid result, follow the initial-send preference
    if evidence is unknown, preserve the last valid comparison
    require valid evidence for a condition or change notification
    identify notifications due under the preferences, including periodic updates
    combine overlapping reasons when one message serves them
    save the new comparison, check progress and pending notifications together

for each due pending notification, including retries:
    recheck current preferences, pause state, quiet hours and relevance
    cancel obsolete messages; refresh evidence needed for claims about current state
    explain unknown evidence in a due status update; never claim an unverified recovery
    compose from verified results, waking the agent only when useful
    send to the configured channel and record the actual outcome
    on known failure, retain pending work for a bounded retry
    on uncertain delivery, reconcile provider status before retrying
```

Preserving comparison progress and pending notifications together prevents a restart
from losing a change that has already advanced the baseline. Use a stable notification
identity across retries and the provider's idempotency support when available. Keep due
periodic updates identifiable by interval so unchanged results are delivered when requested
and the same update is not sent twice. Supersede stale messages after a prolonged outage.

## Delivery

Use the notification path verified during monitor preparation, as part of the routine
above or through the harness's existing hook. A queued or composed message is not a send;
retain its actual outcome and reconcile retries as the pseudocode specifies.

| For message content and replies | Resource | Apply |
|---|---|---|
| Explain the finding or status and handle the user's response | [Monitor notifications](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-project-chat/monitors.md#notify-the-user-and-handle-replies) | Keep the message tied to its evidence, calculation and monitoring preferences |

## Engagement

At each configured opportunity, read the current interests and their monitors, including
those added since registration. Skip paused, declined or expired interests; an empty set
is a successful idle check. Run due calculations through the routine above. Requested
monitor cadence and urgency continue independently of engagement opportunities.

Use recent results and notification history to identify due or potentially useful updates.
Wake the agent only when judgment, composition or interaction is needed. Retain outcomes
so quiet runs do not become repeated suggestions or stale questions.

| When an update needs judgment | Resource | Apply |
|---|---|---|
| Decide what to bring back | [Engagement opportunities](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-project-chat/engagement.md#3-decide-what-to-bring-back-at-each-opportunity) | Honor requested periodic updates and send inferred-interest updates when useful |

## Maintenance

| Needed for this phase | Resource | Apply |
|---|---|---|
| Daily schema and performance review | [Query performance](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-cloud-agent/query-performance.md) | Usage recording, schema refresh, measured optimization and its result |
| Retention boundaries | [Cache retention](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-cloud-agent/cache.md#6-apply-retention-without-losing-required-state) | Preserve protected evidence and operational state |

This is the task's daily review phase: purge, schema refresh, query optimization and
health checks. Run daily by default for every active enrolled project, with additional review after
meaning/schema changes or a slow query. Read the registry again so newly enrolled projects
are included. Isolate per-project failures and bound each stage. Reconcile registered
jobs with desired scope, last success and next runs; surface sustained acquisition or
delivery failures through the existing authorized destination. A watchdog relying only
on its own broken scheduler cannot guarantee failure notification—use native scheduler
failure reporting when available.

```text
for each active enrolled project, with a fair work budget:
    fill due streams and refresh needed product knowledge
    if daily review is due:
        refresh project metadata and the complete Notefile Schema snapshot
        review schema changes and actual query usage; validate useful optimizations
        purge eligible history while preserving required evidence and private state
        record each phase's outcome and next due time independently
    keep last-good data on failure; continue other eligible work and projects
```

For each project, refresh its metadata and full `GetNotefileSchemas` snapshot through
the chosen connection. Stage, validate and compare schema hashes before replacing the
current snapshot; preserve history needed by older data and flag affected projections.
Record schema freshness separately from Event/Device/Session freshness. A schema refresh
failure does not erase a useful cache or make a successful fill appear fully maintained.

Review the recorded fields queried, query fingerprints, frequencies, timings and scan
patterns. Follow Query performance to choose typed projections, physical ordering and
selective indexes when justified. Compare correctness and measured plans before/after,
refresh optimizer statistics after substantial changes when useful, and record applied
or deferred work. Schema discovery alone does not show which columns deserve an index.

Prune Events/Sessions only beyond 45 days using the documented analytical timestamp and
preserve unknown-time rows, active sessions, evidence still needed by monitors and needed baselines.
Retain complete current Devices and Schemas. Pruning is a local cache operation and does
not change an API checkpoint. Rebuild invalidated projections in staging, verify them,
then switch readers atomically. Reconcile affected monitor rules before reactivating them.

Keep database state restartable: bounded run records, appropriate private-state retention,
recoverable configuration and validated job IDs. Cache pruning preserves operational
state, pending deliveries and history needed for recovery. A paused or removed monitor does not
quietly recreate itself from an old job definition. After maintenance, update the readiness
receipt with actual status and remaining gaps rather than resetting every component to green.
Retain separate schema-refresh, purge and optimization outcomes and their next due times;
a single successful subtask does not mark the whole maintenance run complete.
