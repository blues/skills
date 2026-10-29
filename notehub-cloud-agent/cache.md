# Notehub Cache Database: efficient queries with explicit coverage

This is the data contract for the Notehub Cache Database and the Notehub Cache Maintenance
Task that keeps it current. Access prepares both; each newly referenced project joins
their registry. Use prepared questions to choose useful projections.

The current database of choice is DuckDB. Cached Notehub records and derived projections
can be rebuilt; private assistant state stored alongside them must be preserved.
Notehub remains the owner of shared product meaning. Install DuckDB in a harness-managed
environment and keep one writer per database.
Use the selected authenticated connection for every acquisition. Read product skills
before choosing projections; retain raw JSON and stable IDs so corrections can rebuild
views without reacquiring everything.

For concurrent Q&A and polling, choose an actual database-access design: one owning
process handles reads and writes, cooperating processes serialize database opens under
a lock, or the writer publishes stable read snapshots. Merely opening another process
read-only does not make it compatible with a process holding the database for writing.
See [DuckDB concurrency](https://www.duckdb.org/docs/2.0/connect/concurrency) for the runtime contract.

## Choose the relevant reference

Read the rows that apply to the current step; reuse context already established.

| If you need to | Read | Then |
|---|---|---|
| Create the database, register its maintenance task or enroll another project | [Cache preparation](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-cloud-agent/cache-preparation.md) | Reuse the shared infrastructure and verify per-project readiness |
| Decide what to persist beyond cached data | [Database state](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-cloud-agent/workspace.md) | Add only what later runs need and preserve it during cache rebuilds |
| Run the acquisition worker | [Acquisition job contract](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-cloud-agent/jobs.md#acquisition) | Execute bounded work with persisted scope and checkpoints |
| Choose typed projections from product meaning | [Prepared questions](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-cloud-agent/questions.md) | Project only grounded fields and calculations needed by useful questions |
| Acquire or refresh stored product knowledge | [Read product skills](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-access/project-skills.md) | List names and kinds, read relevant documents or the full set, and reuse unchanged content |
| Poll for changes in Events, Sessions, Devices, and more | [Modified cursors](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-project-chat/modified-cursors.md) | Use each endpoint's exact cursor and deletion contract |
| Maintain schemas and optimize the fields actually queried | [Query performance](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-cloud-agent/query-performance.md) | Track field usage and review physical design daily |

## Choose the database organization

Reuse the shared cache established during Cache preparation. Read its folder's README
before querying, then check live coverage/freshness records and follow its access mechanism.
All of the user's bots in the same verified Notehub scope use this cache and maintenance
task; a bot or channel does not get its own copy. Choose physical organization for the
workload and record it in the registry. These logical tables are a useful starting point; adapt names and physical
organization while preserving the identities and relationships. All keys are additionally
scoped to the verified hub/account, either in columns or by an isolated database binding.

| Logical table | Suggested key | Main contents |
|---|---|---|
| Projects | `project_uid` | Project metadata, enrollment state and cache location |
| Devices | `project_uid`, `device_uid` | Current native Device data, raw JSON and deletion/freshness state |
| Sessions | `project_uid`, `session_uid` | Full Session data and its `device_uid`; include device identity in the key if required by the source's uniqueness contract |
| Events | `project_uid`, `event_uid` | Raw Event data, device/session/Notefile references, capture/modified times and tombstones |
| Schemas | `project_uid`, `notefile_id` | Current full schema, fetch time and change hash, with history where needed |
| Monitor/alert state | `project_uid`, `monitor_id`, plus incident/delivery identity where needed | Definitions, baselines, transitions and delivery state for requested monitors |

Use separate typed key columns, rather than building a `projectUID:entityUID` string.
A project can have many monitors and alerts; project UID alone cannot identify each one.
Keep private monitor/incident/outbox state distinct from Notehub source records and shared
product skills. Enforce idempotent ingestion by the actual stable keys; declare constraints
where needed for correctness rather than assuming each composite key improves query speed.

Retain scoped stream checkpoints, coverage, schema/meaning versions and query-usage
records as those workflows run. Source refreshes must
update only the source-owned fields; project enrollment, pauses and private preferences
are not replaced by an API metadata snapshot. Raw JSON plus typed native columns and
workload-driven product projections is a practical combination. Choose the physical layout
to suit the workload while preserving transactional checkpoints and raw evidence.

## Streams and initial coverage

| Entity | Acquisition and retention |
|---|---|
| Projects | Retain metadata for enrolled projects; refresh on first use and during daily maintenance |
| Devices | Complete project scan with modified cursor `0`; upsert by `uid`, apply documented `deleted` tombstones; keep the full current population |
| Events | `GetEvents`, upsert by `event`, apply its `deleted` tombstones; retain at least the latest 30 days and prune only after 45 days |
| Full Sessions | `GetDeviceSessions` for each device, upsert by `session` within project/device scope; updates replace prior session estimates; same recent retention |
| Notefile Schemas | Complete `GetNotefileSchemas` snapshot on enrollment and daily by default, sooner when needed; no modified cursor; swap only after a successful full response |
| Product skills | Read the catalog and needed content, or a complete content snapshot; check metadata periodically, reuse unchanged documents and recheck dependent views/monitors when meaning changes |

A start cursor in epoch microseconds bounds **modification time**, not capture time.
Prefer a bounded recent bootstrap of 45 days where the contract and retention support
it, then verify that the requested capture-time window is covered. If complete recent
coverage cannot be established from that bound, use a complete scan or documented
capture-time backfill plus overlapping modification scan. State service retention gaps.
Do not assume advancing a cursor proves that all historical events exist.

Discover all devices before claiming Session coverage; maintain per-device progress,
including newly discovered devices. Failures on one device leave that part of the Session
population incomplete. A current-device scan may omit removed devices whose historic
Sessions still matter: reconcile identities from retained Events/repository history where
available and state the remaining limit. Do not silently discard retained history when
a Device tombstone arrives. `_session.qo` is supplementary Event evidence, not full Sessions.

Consult the live schema for each stream's tombstone support. In particular, do not invent
a Session `deleted` field or interpret absence in an incremental response as deletion.
Schema snapshot replacement is a separate complete-snapshot operation, not a cursor scan.

## Incremental work, every 15 minutes

Use modified-cursor rules.
Namespace checkpoints by hub/account/project/entity/device/filter and projection contract.
Each maintenance task's fill phase rereads the registry and budgets work fairly across active projects.
Commit a batch's upserts/tombstones and its exact returned cursor in one transaction.
Bulk-load JSON; don't issue per-row requests or SQL inserts. A retry replays safely by
stable identity. Persist last attempt, last completed scan, cursor, partial/error state,
coverage boundaries and the time basis for those boundaries. An empty successful final
batch advances freshness; a failed request never does.

Bound work per invocation. If a backlog exceeds the budget, retain the committed cursor,
mark incomplete and schedule a continuation through the real scheduler. Back off rate
limits with jitter, respecting Retry-After. Do not let a retry loop consume the entire
scheduler or flood the API. Keep the previous usable data on failure, labeled stale or
partial. Log counts and failure classes rather than raw customer records or credentials.

## Raw evidence and typed views

Partition by project and entity. Keep native IDs for joins and names/serials for display.
Build typed per-Notefile views using confirmed product units, scaling, sentinel rules and
version/cohort boundaries. A JSON type conversion failure is a quality flag, not zero.
Preserve raw suspect records; do not delete outliers to make a chart smooth. Cache
coverage and measurement freshness are distinct; a poll can succeed while a sensor has
not reported for days.

Equal values are not duplicates. Deduplicate replay only when the actual identity or
confirmed source/session/sequence mechanism supports it. Track counter resets and
corrected full Session values; don't add a new estimate to an old cumulative estimate.
Current-state queries select latest first and then validity, while period aggregates use
their documented recipe. Product skill changes invalidate derived projections and any
monitor thresholds that depended on the changed rule.

## Retention

Maintain at least 30 days of recent Events/Sessions; keep a 45-day buffer and purge only
records strictly older than that bound using the documented analytical timestamp. Don't
purge rows with unknown/invalid time as if they were old. Preserve ongoing sessions,
monitor evidence, open incidents and required baseline state beyond that window when
needed. Pruning a cache record is not a Notehub deletion. A late modification can reinsert
an older record; apply it idempotently, then evaluate retention independently.

Devices and schemas are full current snapshots, not 45-day histories. Store historical
bindings separately when recipes require them and the evidence supports them. A request
beyond cache coverage goes to a suitable repository or live API; never label the cache's
oldest row as the service's retention policy. Verify bootstrap completeness, schema
snapshots, retention and all-device orchestration as part of the maintenance task.
