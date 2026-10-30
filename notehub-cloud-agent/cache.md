# Build and verify the cache data model

Use Steps 1–2 while creating or verifying the database. Steps 3–6 supply the data
rules for the maintenance worker and project enrollment. Apply the section needed by
the calling procedure, then return there. Preserve an existing sound organization.

## 1. Choose and verify database access

Use DuckDB in the worker runtime established during harness preparation. Cached source
records and projections can be rebuilt; private state and shared product meaning have
separate owners. Use the selected authenticated connection for acquisition.

For concurrent Q&A and polling, choose an actual database-access design: one owning
process handles reads and writes, cooperating processes serialize database opens under
a lock, or the writer publishes stable read snapshots. Merely opening another process
read-only does not make it compatible with a process holding the database for writing.
See [DuckDB concurrency](https://www.duckdb.org/docs/2.0/connect/concurrency) for the runtime contract.

## 2. Create or verify the tables and keys

Reuse the shared cache established during Cache preparation. Read its folder's README
before querying, then check live coverage/freshness records and follow its access mechanism.
All of the user's bots in the same verified Notehub scope use this cache and maintenance
task; a bot or channel does not get its own copy. Choose physical organization for the
workload and record it in the registry. These logical tables are a useful starting point; adapt names and physical
organization while preserving the identities and relationships. All keys are additionally
scoped to the verified account, either in columns or by an isolated database binding.

| Logical table | Suggested key | Main contents |
|---|---|---|
| Projects | `project_uid` | AccountUID, human-readable Project name, associated ProductUIDs, other Project metadata, enrollment state and cache location |
| Devices | `project_uid`, `device_uid` | Current Device data, FleetUID memberships preserving any `smart:` prefix, raw JSON and deletion/freshness state |
| Sessions | `project_uid`, `session_uid` | Full Session data and the DeviceUID of the one Device that generated it |
| Events | `project_uid`, `event_uid` | Raw Event data, DeviceUID of the one Device that generated it, SessionUID and Notefile references, capture/modified times and tombstones |
| Schemas | `project_uid`, `notefile_id` | Current full schema, fetch time and change hash, with history where needed |
| Monitors | `project_uid`, `monitor_id` | What is watched, requested or inferred origin, checking and notification preferences, pause state and progress needed to resume |

Here `project_uid` stores ProjectUID (also called AppUID); those names refer to the
same unique project GUID. The project's multiple ProductUID strings identify
customer product types or SKUs feeding that project and do not create separate
project cache entries. Use separate typed key columns, such as `project_uid` and
`device_uid`, rather than concatenating ProjectUID and DeviceUID into one string. A
Project can have many monitors; ProjectUID alone cannot identify each one. Keep private
monitor preferences and progress distinct from Notehub source records and shared product
skills. Enforce idempotent ingestion by the actual stable keys;
declare constraints where needed for correctness rather than assuming each composite
key improves query speed.

Retain scoped stream checkpoints, coverage, schema/meaning versions and query-usage
records as those workflows run. Source refreshes must
update only the source-owned fields; project enrollment, pauses and private preferences
are not replaced by an API metadata snapshot. Raw JSON plus typed native columns and
workload-driven product projections is a practical combination. Choose the physical layout
to suit the workload while preserving transactional checkpoints and raw evidence.

### State that survives rebuilds

Start with the project registry. As workers and queries run, add only the records they
need: exact acquisition checkpoints, coverage/freshness, task definitions and run results.
When monitors are used, retain their definitions, last valid comparisons, check and
notification progress, interests and pauses. Let the agent choose the organization;
there is no separate required assistant-state schema. Reuse information reliably
available from the harness.

Keep structured state in the database rather than a parallel directory of JSON files.
Worker scripts, engine files and temporary transfers still belong in the filesystem.
The harness needs only enough bootstrap information to locate the database and invoke
the worker. Keep credentials in the selected tool.

Commit source records with their checkpoints, and monitor comparisons with any pending
notifications. Rebuilding source caches must preserve operational state, user choices and
delivery history. Task registration and successful runs/delivery are distinct evidence.
The folder README and shared memory locate the cache; live state stays in the database.

## 3. Define streams and verify initial coverage

| Entity | Acquisition and retention |
|---|---|
| Projects | Retain metadata for enrolled projects; refresh on first use and during daily maintenance |
| Devices | Complete Project scan with modified cursor `0`; upsert by DeviceUID in `uid`, apply documented `deleted` tombstones; keep the full current population |
| Events | `GetEvents`, upsert by `event`, apply its `deleted` tombstones; retain at least the latest 30 days and prune only after 45 days |
| Full Sessions | `GetDeviceSessions` for each DeviceUID, upsert by SessionUID in `session` within the ProjectUID and DeviceUID; updates replace prior session estimates; same recent retention |
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

## 4. Commit incremental changes with their checkpoints

Implement the batch transaction below using the endpoint's current cursor contract.

| For the acquisition contract | Read | Then |
|---|---|---|
| Poll for changes in Events, Sessions, Devices, and more | [Modified cursors](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-project-chat/modified-cursors.md) | Apply exact cursor, batching and deletion semantics |

Namespace checkpoints by account, ProjectUID, entity, DeviceUID, filters and
projection contract. Each maintenance task's fill phase rereads the registry and
budgets work fairly across active projects. Commit a batch's upserts/tombstones and
its exact returned cursor in one transaction. Bulk-load JSON; don't issue per-row
requests or SQL inserts. A retry upserts by the ProjectUID and DeviceUID, EventUID
or SessionUID keys defined above. Persist last attempt, last completed scan, cursor,
partial/error state, coverage boundaries and the time basis for those boundaries. An
empty successful final batch advances freshness; a failed request never does.

Bound work per invocation. If a backlog exceeds the budget, retain the committed cursor,
mark incomplete and schedule a continuation through the real scheduler. Back off rate
limits with jitter, respecting Retry-After. Do not let a retry loop consume the entire
scheduler or flood the API. Keep the previous usable data on failure, labeled stale or
partial. Log counts and failure classes rather than raw customer records or credentials.

## 5. Build typed views from verified meaning

Read the relevant product knowledge before choosing product projections. Create or
verify only views required by actual questions; raw records can be acquired first.

Partition by ProjectUID and entity. Join Devices by DeviceUID, Events by EventUID
and Sessions by SessionUID; use names and serial numbers for display. Build typed
per-Notefile views using confirmed product units, scaling, sentinel rules and
version/cohort boundaries. A JSON type conversion failure is a quality flag, not
zero. Preserve raw suspect records; do not delete outliers to make a chart smooth.
Cache coverage and measurement freshness are distinct; a poll can succeed while a
sensor has not reported for days.

Equal values are not duplicates. Deduplicate replay only when EventUID or a
confirmed source/session/sequence mechanism supports it. Track counter resets and
corrected full Session values; don't add a new estimate to an old cumulative
estimate. Current-state queries select latest first and then validity, while period
aggregates use their documented recipe. Product skill changes invalidate derived
projections and any monitor thresholds that depended on the changed rule.

## 6. Apply retention without losing required state

Maintain at least 30 days of recent Events/Sessions; keep a 45-day buffer and purge only
records strictly older than that bound using the documented analytical timestamp. Don't
purge rows with unknown/invalid time as if they were old. Preserve ongoing sessions,
evidence still needed by monitors and required comparison state beyond that window when
needed. Pruning a cache record is not a Notehub deletion. A late modification can reinsert
an older record; apply it idempotently, then evaluate retention independently.

Devices and schemas are full current snapshots, not 45-day histories. Store historical
bindings separately when recipes require them and the evidence supports them. A request
beyond cache coverage uses the Project Repository when MCP access and an enabled
Repository are available, or the ordinary Notehub API. Never label the cache's oldest
row as the service's retention policy. Verify bootstrap completeness, schema
snapshots, retention and all-device orchestration as part of the maintenance task.

Return the verified model or coverage result to the procedure that called this resource.

## Batch size and cursor starts for acquisition

Cache fill **always** uses **`pageSize=10000`** (10K-row batches) for Devices, Events,
and Sessions wherever the API accepts `pageSize`. Smaller pages are an anti-pattern
for catch-up and for steady-state incremental polls.

| Stream | Time-windowed? | Initial cursor |
|---|---|---|
| Devices | No | **`0`** (full scan) |
| Sessions (per device) | No | **`0`** (full scan per DeviceUID) |
| Events | Yes (modified cursor is epoch microseconds) | Synthesize start as Unix epoch **microseconds** for the desired lookback (e.g. 14 days); **do not** start Events at `0` when filling a recent window. Resume with opaque `next_cursor` thereafter. |
| Notefile Schemas | N/A | Full snapshot; no modified cursor |

Loop while `has_more` is true. After catch-up, periodic (~15m) runs only fetch
incremental changes from the saved cursors.

