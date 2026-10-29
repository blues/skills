# Poll for changes in Events, Sessions, Devices, and more

Modified cursors let Project Chat read successive batches of created, updated or deleted
records, and let cloud workers keep a cache current by fetching only changes since their
last successful checkpoint. The same mechanism works through MCP and the Notehub CLI.
For a foreground question, use it for the bounded read and stop when the work is done.
On the user's Cloud Computer, resume saved cursors at each periodic acquisition run.

## 1. Verify the stream and cursor contract

Use modified cursors for incremental synchronization wherever the operation
documents them. Start with `0` for a complete initial entity scan, or an epoch time
in **microseconds** when a deliberately bounded starting point is appropriate. A
returned cursor is an **opaque string**, potentially including a tie-breaker:
preserve it exactly. Never parse it into a float, increment it, derive it from the
last row's timestamp, or replace it with an EventUID. Old EventUID cursors and
`/events-cursor` are not the sync mechanism.

Check the exact response: Events, Devices and Device Sessions expose `next_cursor` and
`has_more`; some operations return arrays and cursor metadata in `X-Next-Cursor` and
`X-Has-More` headers. There is no universal envelope. Preserve headers when needed.

`next_cursor` identifies where to resume; `has_more` says whether another batch is
available now. Request a large batch within the operation's documented limit and the
worker's response capacity. Where supported, `pageSize` is the API parameter for batch
size; modified-cursor scans advance with `cursor`, never `pageNum`.

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
