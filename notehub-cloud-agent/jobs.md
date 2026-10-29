# Jobs that work without conversation memory

Access creates the Notehub Cache Maintenance Task using the acquisition and daily review
contracts below; first-use project enrollment supplies the scopes it maintains. Implement them
with the harness's actual worker/scheduler tools and database state.
Use deterministic code and SQL for routine work, in a language supported by the harness.
The pseudocode below describes behavior; its names are not callable APIs or a prescribed
database schema. Implement and validate it using the selected authenticated transport,
database and scheduler.

## Choose the relevant reference

Read the rows that apply to the current step; reuse context already established.

| If you need to | Read | Then |
|---|---|---|
| Prepare project queries or requested monitor jobs | [Setup](setup.md) | Reuse the cache workers and verify the requested project capabilities |
| Create the cache database and maintenance task during Access | [Cache preparation](cache-preparation.md) | Verify database persistence and scheduled execution before project enrollment |
| Recover durable execution inputs and state | [Database state](workspace.md) | Use the recorded scope and configuration rather than conversation memory |
| Implement monitor comparisons and alert delivery | [Monitor pseudocode](../notehub-project-chat/monitors.md#evaluation-and-delivery-pseudocode) | Adapt the algorithm to the requested monitor and actual delivery tools |
| Establish bootstrap, coverage, concurrency or retention rules | [Cache](cache.md) | Preserve per-stream coverage and atomic checkpoints |
| Fetch incremental changes for the acquisition job | [Modified cursors](../notehub-project-chat/modified-cursors.md) | Use the exact saved cursor, process successive batches and persist progress |
| Refresh the knowledge used by questions and monitors | [Read product skills](../notehub-access/project-skills.md) | Check the metadata catalog and fetch needed changes without rewriting shared skills |
| Define or revise a calculation before evaluating it | [Prepared questions](questions.md) | Use its validated meaning revision and evidence requirements |
| Choose relevant findings and tone for engagement | [Engagement](../notehub-project-chat/engagement.md) | Produce due status reports or useful updates, respecting user preferences |
| Refresh schemas and optimize recurring queries | [Query performance](query-performance.md) | Use daily schema checks, field-usage records and measured query plans |

## Common execution context

Store the worker's actual executable/script path or tool invocation, database locator,
project/monitor scope and configuration version in database job records. Give the harness
enough bootstrap information to open the database safely and find that job record.
The worker reads configuration and product revisions from the database,
verifies the selected identity/hub, obtains the appropriate lock, performs bounded work,
commits state and returns a small machine-readable run result. No credential is placed in
a command, task prompt, script or database record.

The Notehub Cache Maintenance Task reads the current enrolled-project registry on
every run, including enrollments made after it was created. Process active scopes
only; preserve pauses/exclusions and isolate failures by project. An empty registry is
an idle success. Keep per-project budgets and continuation state so a large backlog does
not monopolize the service. By default the task runs every 15 minutes, with per-project
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
`paused`), relevant project/monitor, counts, evidence freshness/coverage and continuation
need. Include error class/action, not secrets or raw customer dumps. A success marker
written before all required steps complete is not a successful run. Persist registration
and execution receipts separately.

## Acquisition

The maintenance task's fill phase normally runs every 15 minutes across enrolled projects. For each due project,
serialize writes per database, keep stream definitions stable, and budget work across
streams/devices so a large Event backlog cannot starve other entities or projects.

1. Validate context. An account/hub/project mismatch stops this scope; do not reuse its
   cursor under another identity. A network failure does not initiate repeated login.
2. Poll Devices through the saved modified cursor using transactional batch ingestion.
   Enroll new device identities into Session streams. Preserve historical records when
   a Device is deleted and track historical-device coverage gaps separately.
3. Poll Events and each due full-Session stream with their own saved exact cursors.
   Commit each batch with its checkpoint. A batch budget leaves a resumable partial stream.
   Apply only the entity's documented deletion markers; equal values are not duplicates.
4. On first fill, fetch project metadata and the complete Notefile Schema snapshot to
   staging. Validate and commit it before marking that part ready. Daily maintenance
   refreshes these thereafter; a new field need can trigger an earlier refresh. A failed
   or partial response retains the previous snapshot and freshness. An authoritative
   empty snapshot is different from an error. Schemas do not acquire an invented cursor.
5. Read product-skill metadata to check for changes. Reuse unchanged documents; fetch
   needed new or changed content individually, or get a full snapshot when simpler.
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

Apply this to each due modified-cursor stream. Use the database access design chosen in
Cache, including serialization of overlapping work. Each attempt belongs to the verified
hub/account/project and stable stream definition.

```text
cursor = saved exact cursor, or the stream's fixed bootstrap cursor
seen = {cursor}
while this stream has time and batch budget:
    claim the current attempt for this stream and cursor; mark the scan in progress
    fetch a batch through the selected authenticated connection
    retain the actual fetch-completion time
    validate records, stable identities, deletion markers and cursor metadata
    reject a repeated/cyclic cursor, except an unchanged cursor with has_more=false

    in one database transaction:
        verify this attempt and its scope/cursor are still current
        bulk upsert records and documented tombstones by stable identity
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

Normalize the endpoint's documented body or header metadata without inventing missing
values. Reject malformed records and conflicting duplicate identities within a batch;
equal measurements with different identities remain separate. An empty final batch may
keep the same cursor or advance it; commit its successful freshness either way. A failed
or unfinished scan retains the last complete success and its coverage limits.

When only scheduled agent turns can invoke MCP, the turn performs the fetch through that
same connection and passes the response to its database worker using a supported harness
mechanism. Carry the scope, requested cursor, current attempt identity and actual fetch
time with the response. Reject superseded results and late failures, including when a
newer empty scan kept the same cursor. Never relabel an old response with ingestion time.
Resolve abandoned attempts after the worker timeout. Choose a direct bridge when supported;
do not assume MCP tools are local imports or add a CLI login to work around that limitation.

### Validate the implemented acquisition worker

Use synthetic batches and a disposable database in the actual worker runtime before
activation. These cases test behavior without customer credentials:

| Case | Expected result |
|---|---|
| Updates, documented tombstones and equal values under different IDs | Updates replace by identity, tombstones apply, distinct records survive |
| Empty final batch, including an unchanged cursor | Successful scan and freshness are committed; no extra empty request |
| Missing metadata, conflicting identities, cursor cycle or failed transaction | Failing batch advances nothing; earlier commits survive |
| Budget exhaustion, process restart or changed project/filter scope | Resume exactly from the committed cursor; a different scope uses its own checkpoint |
| Late MCP result/failure or replay after a newer completion | No overwrite, freshness regression or degradation of the newer success |

## Evaluation

For each due monitor, periodic status report or engagement candidate, load its prepared question,
meaning revision, latest coverage and persisted state. Check that its data satisfies the
specified freshness, scope and time window; fetch bounded fresh evidence through the
same transport if needed and authorized. A five-minute monitor cannot obtain five-minute
freshness by querying an unchanged 15-minute cache.

Compute the result and its quality/coverage metadata. Record the fields and query plan
used by each evaluator for daily performance review. For a condition monitor, the
predicate yields true, false or unknown. Apply product-specific dwell/hysteresis before the incident transition;
unknown/stale evidence cannot close a condition. Follow the Monitor pseudocode for the
requested comparison mode. Persist the resulting incident state and
pending candidate in the **same transaction** so a crash cannot lose a notification while
advancing the state that would suppress it on retry.

For a result-change subscription, compare the current valid result with its saved
baseline using the configured equality/significance rule. Every meaningful difference
gets its own change ID, even on consecutive polls; atomically persist the new baseline
and outbox candidate. Unknown results do not advance the baseline. Condition transitions
and result comparisons use their respective rules.

A candidate carries stable monitor/incident/transition IDs, relevant evidence timestamps,
query and meaning revisions, population/window/units, severity, unknowns, relevance expiry
and intended destination. Deduplication includes project and monitor scope. Repeated
unchanged conditions ordinarily produce no candidate; reminders or renewed escalation
follow the configured rule. A due status report is a separate candidate keyed to its
reporting interval; unchanged incident state does not suppress it. Do not call the model just to decide whether a number crossed
a confirmed threshold. Model reasoning is useful when a candidate needs explanation,
source reconciliation or a consequential user question.

## Delivery

Use the scheduler's completion hook or a bounded delivery worker to drain pending items.
Apply pauses, quiet hours, authorization, suppression and expiry before composition.
Recheck fresh evidence and incident state immediately before sending; cancel/supersede
an obsolete closeout or an alert whose premise changed. Recompute reported quantities
with the original threshold/population/window, not a neighboring statistic.

Compose a useful message from verified evidence when needed. Deliver to the recorded
surface using the actual harness API, with a stable idempotency key if supported. Mark
acknowledged only after the provider's receipt; a timeout leaves an uncertain send to
reconcile, not a new incident. Persist attempted sends and provider IDs. Exactly-once
delivery requires provider support; otherwise use the best available reconciliation and
state the limitation when it affects behavior. Expired candidates are not replayed as
fresh news after a prolonged outage.

Closeouts and coverage-failure notices follow the same durable delivery path. Third-party
recipients need applicable explicit authorization. A useful correction updates dependent
state and prior conclusions, not just the next message's wording.

## Engagement

At each configured local-time opportunity, consult the interest/exclusion profile and
recent feedback/delivery history, run the relevant prepared queries, and check for due
status reports as well as meaningful changes. Use Engagement for relevance and tone.
A requested periodic status report can contain unchanged results, a graph or both;
report missing/stale evidence honestly rather than asserting stability. Bundle related
findings. Outside a requested status cadence, no useful candidate means no message.
Explicit monitor urgency/cadence remains separate
from these check-in windows. Write down outcomes so repeated quiet runs do not become
repeated suggestions or stale questions.

## Maintenance

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
preserve unknown-time rows, active sessions, open incident evidence and needed baselines.
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
