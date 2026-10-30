# Poll for changes in Events, Sessions, Devices, and more

Modified cursors let Project Chat read successive batches of created, updated or deleted
records, and let cloud workers keep a cache current by fetching only changes since their
last successful checkpoint. The same mechanism works through MCP and the Notehub CLI.
For a foreground question, use it for the bounded read and stop when the work is done.
On the user's Cloud Computer, resume saved cursors at each periodic acquisition run.

## 1. Verify the stream and cursor contract

Use modified cursors for incremental synchronization wherever the operation
documents them.

**Cursor start rules (required):**
- Entities that are **not time-windowed** (Devices, and full per-device Sessions
  scans): always start at **`0`** for a complete initial scan. Do not synthesize a
  start time for those streams.
- Entities that **are time-windowed** (Events, and any other stream whose modified
  cursor is documented as an epoch-microsecond time cursor): **never** need to start
  at `0` for catch-up of “recent” history. Start at a **synthesized Unix epoch
  microseconds** value for “how far back I want to begin” (for example ~14 days ago
  in µs since epoch). The cursor **is** documented as an epoch-microsecond cursor;
  after the first response, resume only with the returned opaque `next_cursor`.
  Do not invent mid-stream cursors by parsing floats or bumping timestamps yourself.

A returned cursor is an **opaque string**, potentially including a tie-breaker:
preserve it exactly once issued. Never parse a returned cursor into a float, increment
it, derive it from the last row's timestamp, or replace it with an EventUID. Old
EventUID cursors and `/events-cursor` are not the sync mechanism. (Synthesizing the
*initial* start as epoch microseconds for time-windowed streams is allowed and
required when bounding how far back to fill; that is not parsing a returned cursor.)

Check the exact response: Events, Devices and Device Sessions expose `next_cursor` and
`has_more`; some operations return arrays and cursor metadata in `X-Next-Cursor` and
`X-Has-More` headers. There is no universal envelope. Preserve headers when needed.

`next_cursor` identifies where to resume; `has_more` says whether another batch is
available now. **One `call_tool` per MCP `execute` block (required for batch fetch):**
When filling or polling via Notehub MCP code-mode `execute`, each `execute` block must
contain **exactly one** `await call_tool(...)` for the fetch. Do **not** chain multiple
`get_events` / `get_devices` / `GetDeviceSessions` calls inside one execute (the sandbox
caps `call_tool` invocations per block, and packing pages invites timeouts/OOM). Ingest
and cursor advance happen **outside** execute (agent-turn handoff / local worker) after
that single response is written to a file. A new execute is required for the next batch.

**Always request `pageSize=10000` (10K batches)** for modified-cursor
scans wherever the operation supports `pageSize`, up to the operation’s documented
maximum. Do not default to 500 or other small pages for cache fill. Where supported,
`pageSize` is the API parameter for batch size; modified-cursor scans advance with
`cursor`, never `pageNum`.

Keep each record under its ProjectUID and the endpoint's documented key:

| Entity | Identifier in the response |
|---|---|
| Device | DeviceUID in `uid` |
| Event | EventUID in `event` |
| Session | SessionUID in `session`, associated with DeviceUID in `device` |

For other entities, use the key named in the operation's response schema.

## 2. Read and commit successive batches

1. Scope a checkpoint to the account, ProjectUID, entity, filters and contract version;
   include DeviceUID for per-device Session streams.
2. Fetch one batch with the saved cursor through the chosen authenticated transport.
3. Validate the response and its cursor metadata. With `has_more=true`, a missing,
   repeated or cyclic next cursor is an error, not completion.
4. Apply records by the keys above, including documented deletion tombstones, and
   save the returned cursor only after the batch is accepted. A persistent cache commits
   records and cursor in the **same transaction**. A foreground analysis can accumulate
   its bounded results in memory or temporary storage; it needs no persistent database.
5. Fetch the next batch immediately while `has_more=true`. When `has_more=false`,
   the scan is complete; keep the returned cursor for the next scheduled poll. There
   is no need for an extra request just to obtain an empty batch. An empty successful
   batch can still advance the checkpoint and freshness. Never advance it on a failed
   fetch, invalid batch or failed transaction.

## 3. Verify what the completed scan covers

A cursor orders **modifications**, not capture time. Old records may be corrected now;
late events may describe old measurements. Retention and analytical windows are separate
from the cursor. A Note's `note.delete` event is not automatically a deletion of that
Event row: follow the entity's tombstone contract. Do not deduplicate by equal sensor
values or drop raw suspect evidence.

Devices and Events support project-level modified scans. Full Sessions are retrieved
per device with `GetDeviceSessions`; `_session.qo` is not a replacement. Notefile schemas
use `GetNotefileSchemas`; where no modified cursor is documented, refresh a complete
snapshot and compare only after a successful complete fetch. Never invent a cursor
parameter or infer deletion from a partial snapshot.

Return the saved checkpoint and actual coverage to the calling query or acquisition worker.
