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
| Projects | `project_uid` | Verified account, human-readable Project name, associated ProductUIDs, other Project metadata, enrollment state and cache location |
| Devices | `project_uid`, `device_uid` | Current Device data, FleetUID memberships as returned, raw JSON and deletion/freshness state. The API lists Smart Fleet members by the plain FleetUID and omits deleted Fleets; use `GetFleets` (each Fleet's smart rule) to tell which memberships are rule-driven |
| Sessions | `project_uid`, `session_uid` | Full Session data and the DeviceUID of the one Device that generated it |
| Events | `project_uid`, `event_uid` | Raw Event data, DeviceUID of the one Device that generated it, SessionUID and Notefile references, capture (`when`) and arrival (`received`) times, the batch cursor that delivered them, and tombstones. Events carry no per-row modification time |
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
| Full Sessions | `GetSessions`, one project-wide stream with a modified cursor; upsert by SessionUID in `session` with the DeviceUID in `device`; updates replace prior session estimates; same recent retention. Use `GetDeviceSessions` only for a targeted single-device read |
| Notefile Schemas | Complete `GetNotefileSchemas` snapshot on enrollment and daily by default, sooner when needed; no modified cursor; swap only after a successful full response |
| Product skills | Read the catalog and needed content, or a complete content snapshot; check metadata periodically, reuse unchanged documents and recheck dependent views/monitors when meaning changes |

A start cursor in epoch microseconds bounds **modification time**, not capture time.
Bootstrap Events and Sessions from a synthesized start cursor **30 days ago** (never page
through older data the cache will not use; the 45-day figure is only the purge buffer)
where the contract supports it, then verify that the requested capture-time window is
covered. If complete recent coverage cannot be established from that bound, use a complete
scan or documented capture-time backfill plus overlapping modification scan. State service
retention gaps. Do not assume advancing a cursor proves that all historical events exist.

On the Essentials plan, the Events API (cursor reads included) returns only Events Notehub
stored or deleted in the last 7 days, plus the Project's 100 newest Events, with no signal
that older ones exist. A 30-day bootstrap there yields about 7 days, and the cache builds
longer history only from enrollment forward. Record that as service retention, not cache
coverage. Retention purges leave no tombstones, so an Event's absence from a feed is not a
deletion.

One project-wide `GetSessions` stream covers every Device's Sessions, including Devices
that appear later; do not enumerate Devices to fetch Sessions. Do not silently discard
retained history when a Device tombstone arrives. `_session.qo` is supplementary Event
evidence, not full Sessions.

Consult the live schema for each stream's tombstone support. In particular, do not invent
a Session `deleted` field or interpret absence in an incremental response as deletion.
Schema snapshot replacement is a separate complete-snapshot operation, not a cursor scan.

## 4. Commit incremental changes with their checkpoints

Implement the batch transaction below using the endpoint's current cursor contract.

| For the acquisition contract | Read | Then |
|---|---|---|
| Poll for changes in Events, Sessions, Devices, and more | [Modified cursors](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-project-chat/modified-cursors.md) | Apply exact cursor, batching and deletion semantics |

Namespace checkpoints by account, ProjectUID, entity, filters and
projection contract. Each maintenance task's fill phase rereads the registry and
budgets work fairly across active projects. Commit a batch's upserts/tombstones and
its exact returned cursor in one transaction. Bulk-load JSON; don't issue per-row
requests or SQL inserts. A retry upserts by the ProjectUID and DeviceUID, EventUID
or SessionUID keys defined above. Persist last attempt, last completed scan, cursor,
partial/error state, coverage boundaries and the time basis for those boundaries. An
empty successful final batch advances freshness; a failed request never does.

Bound work per invocation. The Notehub API allows 7 requests per minute per billing
account (Organization), shared by every client of that Organization; a request over the
limit returns 429 with `Retry-After` and `X-RateLimit-*` headers. Budget pages per run
against that quota: this is why Devices, Events and Sessions each page as one project-wide
stream. If a backlog exceeds the budget, retain the committed cursor, mark incomplete and
schedule a continuation through the real scheduler. On 429 wait at least `Retry-After`,
with jitter; the `notehub` CLI already waits it out and retries. Do not let a retry loop
consume the entire scheduler or flood the API. Keep the previous usable data on failure,
labeled stale or partial. Log counts and failure classes rather than raw customer records
or credentials.

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
snapshots, retention and stream orchestration as part of the maintenance task.

Return the verified model or coverage result to the procedure that called this resource.

## Typed side table, refresh and query logging (Grok Bot reference)

This is the measured layout used on Grok Bot Cloud Computers (cache folder
`/workspace/notehub-cache`, venv `/home/box/venvs/duckdb`). Adapt the promoted
columns to the fields your own `query_log.ndjson` shows are extracted repeatedly.

**Why.** On ~4.5M Events, JSON extraction from `events.data` dominated query time.
A typed side table made real scripts 3–17× faster and per-device latest-value
aggregates 100–190× faster with identical results. Re-sorting `events` gave little,
because ingest appends in roughly time order anyway. ART indexes on `device` or `event`
gave no benefit to the observed queries and made each ingest batch ~5× slower, so
`events_x` has **no primary key and no index**. A full sort/rebuild must run with a
memory cap (`memory_limit='4GB'`); an uncapped sort was killed by the OOM killer.

**Table.** `events_x` mirrors `events` row for row:

```sql
CREATE TABLE events_x AS
SELECT project_uid, event, device, file, captured, fetched_at,
  coalesce(json_extract_string(data,'$.sn'), json_extract_string(data,'$.best_id')) AS sn,
  json_extract_string(data,'$.best_country')  AS best_country,
  json_extract_string(data,'$.best_location') AS best_location,
  try_cast(json_extract_string(data,'$.body.cpm')   AS DOUBLE) AS cpm,
  try_cast(json_extract_string(data,'$.body.usv')   AS DOUBLE) AS usv,
  try_cast(json_extract_string(data,'$.body.csecs') AS DOUBLE) AS csecs,
  try_cast(json_extract_string(data,'$.best_lat') AS DOUBLE) AS best_lat,
  try_cast(json_extract_string(data,'$.best_lon') AS DOUBLE) AS best_lon,
  json_extract_string(data,'$.best_location_type') AS best_location_type,
  try_cast(json_extract_string(data,'$.where_lat') AS DOUBLE) AS where_lat,
  try_cast(json_extract_string(data,'$.where_lon') AS DOUBLE) AS where_lon
FROM events ORDER BY project_uid, file, captured;
```

`best_lat`/`best_lon` are Notehub's chosen best fix, which already falls back from GPS
(`where_*`) to triangulation (`tri_*`), then to the cell tower (`tower_*`); in sampled
data no Event lacked `best_lat` while having another fix, so no `coalesce` is needed.
`where_lat`/`where_lon` keep the GPS-only fix for track analysis. Promote a field when the
query log shows bots joining back to `events.data` for it: Event UIDs are random, so that
join scans every `events.data` row (about 3.8 CPU-seconds on 4M rows; one such query took
3.7 s and fell to 8 ms once `best_lat`/`best_lon` were promoted). If the refresh finds a
listed column missing, it rebuilds the table.

**Refresh after every ingest.** `step.sh` runs `refresh_events_x.py` immediately after
`handoff.py ingest`. It takes the same exclusive flock on `notehub.duckdb.lock` as
ingest and runs one transaction (about 80–130 ms for a normal 15-minute fill). The
lookback absorbs files ingested slightly out of fetch order:

```sql
-- w = SELECT max(fetched_at) FROM events_x   (build the table if missing)
CREATE TEMP TABLE _nx AS SELECT <same column list> FROM events
  WHERE fetched_at > $w - INTERVAL 30 MINUTE;
DELETE FROM events_x WHERE (project_uid, event) IN (SELECT project_uid, event FROM _nx);
INSERT INTO events_x SELECT * FROM _nx;
-- self-heal: if count(events) <> count(events_x), delete keys missing from events
-- and insert keys missing from events_x
```

The refresh is fail-safe: on any error it appends to `errors.log` and exits 0, so ingest
and the cursor advance are never affected. It writes timings to `refresh_timing.csv`, and
`--rebuild` recreates the table (~16 s under the lock). The daily maintenance purge in
`worker.py` applies the same 45-day rule to `events_x`
(`DELETE FROM events_x WHERE fetched_at < now() - INTERVAL 45 DAY`).

**Readers and query logging.** `nhq.py` in the cache folder wraps `duckdb.connect`:
`nhq.connect(bot)` takes a shared flock on `notehub.duckdb.lock`, so fills wait for
readers and readers wait for fills, opens read-only, and retries on lock errors. Each
`execute()` appends `{ts, bot, script, sql, nparams, ms, rows, error}` to
`query_log.ndjson`; logging never raises. Readers should keep queries short and close
the connection. This log is the query-usage record for
[query performance](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-cloud-agent/query-performance.md) reviews.
Document both in the cache folder's README so every bot finds them.

## Batch size and cursor starts for acquisition

Cache fill **always** uses **`pageSize=10000`** (10K-row batches) for Devices, Events,
and Sessions wherever the API accepts `pageSize`. Smaller pages are an anti-pattern
for catch-up and for steady-state incremental polls.

| Stream | Time-windowed? | Initial cursor |
|---|---|---|
| Devices | No | **`0`** (full scan) |
| Sessions | Yes (same epoch-microsecond modified cursor as Events) | One project-wide `GetSessions` stream. Synthesize start as Unix epoch **microseconds** 30 days ago, as for Events; resume with opaque `next_cursor` thereafter. |
| Events | Yes (modified cursor is epoch microseconds) | Synthesize start as Unix epoch **microseconds** 30 days ago (the bootstrap above); **do not** start Events at `0` when filling a recent window. Resume with opaque `next_cursor` thereafter. |
| Notefile Schemas | N/A | Full snapshot; no modified cursor |

Loop while `has_more` is true. After catch-up, periodic (~15m) runs only fetch
incremental changes from the saved cursors.

## One call_tool per execute for MCP batch fetch

When acquisition uses Notehub MCP `execute` (agent-turn fetch handoff), each fetch batch
must be **one** `execute` block with **exactly one** `await call_tool(...)` (typically
`get_events` / `get_devices` / `get_sessions` at `pageSize=10000`, or the largest size
the transport allows). Do not pack multiple pages into one execute. Write the envelope
to a file, ingest locally, then start a **new** execute for the next batch while
`has_more` is true.

