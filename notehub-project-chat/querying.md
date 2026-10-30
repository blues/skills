# Answer a question using the best available evidence

Use this inside Project Chat's conversation loop. Reuse the selected ProjectUID,
authenticated connection and product knowledge already internalized. Work through the
steps below, reusing a verified calculation when its meaning still matches the question.

## 1. Establish what the question means

Identify the quantity, population, time basis and intended decision. Use the conversation
and loaded skills to resolve product terminology, NotefileID, Event `body` fields and
database columns. Ask only about ambiguity that changes the answer. Analytical times
default to UTC unless the user specifies otherwise.

| If you need to | Read | Then |
|---|---|---|
| Recheck a product rule after a change or recover missing knowledge | [Reading product knowledge](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-project-chat/product-knowledge.md) | Reuse the full skill set already read; refresh or recover it only when needed |
| Interpret native records or system fields | [Platform](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-platform/SKILL.md) | Choose the reference for the actual data surface |

Use the project's confirmed vocabulary and field inventory, including units, scaling,
array positions, sentinels, missing/zero distinctions, conditional emission and firmware
scope. For native records, use the platform route above.

For example, “radiation data in the last two weeks within 100 km of Kyiv” needs product
knowledge for the radiation field, units, validity and interpretation, plus platform
knowledge for capture time, location provenance and geographic filtering. Establish the
requested spatial reference, select the two-week window and associate the appropriate
location with each measurement. Do not substitute today's Device location for historical
Event location or assume every radiation field is a comparable dose-rate measurement.

The producer determines which knowledge applies: native system Notefiles also have
bodies, and a familiar field name does not establish customer measurement semantics.

## 2. Choose the cache, Repository or Notehub API

**Before choosing a source:** if this agent is on a Cloud Computer and the selected
Project is not yet in the Notehub Cache registry, **stop and enroll it now** via
[Project selection Step 3](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-access/projects.md#3-enroll-the-project-and-check-its-cloud-cache),
begin fill immediately, and report that to the user. Do not answer the question over the
live API as a shortcut.

**ALWAYS route Q&A through the Notehub Cache when it is present** (Project enrolled).
That is the default for ordinary questions — **including while fill is still in
progress**. Answer from verified available coverage and last successful acquisition;
inspect those before answering. Do not wait for fill to finish if the cache already
covers the fields needed for this question.

Only when the cache is **absent** (not enrolled and cannot be) or the question needs
history/fields the cache cannot supply, use the Project Repository **if both MCP access
and an enabled Repository are available for this Project**. Its procedure handles
introspection and query execution; bring its result back here for any further calculation
or combination with other evidence.

**MUST NOT** use live Notehub list/get APIs for ordinary fleet, device, event, session,
or location analytics when the Notehub Cache is present / enrolled (including while the
cache is initializing). Live API is only for evidence the cache and Repository cannot
supply, or to verify current service/control-plane state.

**Local Notehub Cache queries are fastest; Repository queries are slower than the cache
but faster than direct Notehub API queries.** Use the APIs for required evidence that
neither available source supplies, or to verify current service state.

| Situation | Source | Read when the source or contract is unfamiliar |
|---|---|---|
| The question concerns the whole dataset, no cache is available, or required history extends beyond it; MCP access and an enabled Repository are available | Introspect and query the Project Repository through MCP | [Project Repository](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-project-repository/SKILL.md) |
| An ordinary question fits the Notehub Cache's fields, population, window and freshness | Query the local cache | [Cache coverage](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-cloud-agent/cache.md) |
| The cache and MCP Repository cannot supply the required evidence, or current resource state must be verified | Notehub API through the established MCP or CLI connection | [OpenAPI contract](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-access/openapi.yaml) |
| No source covers the request | Explain the missing evidence and offer the useful bounded answer | No additional reference |

Check coverage **before** querying: account, ProjectUID, entity, filters, complete
scan, last successful poll, capture-time window, ingestion delays and deletion
handling. Caches can be fresh but incomplete, or complete but stale. A repository
can contain Events but not full Sessions. If combining sources, match Devices by
DeviceUID, Events by EventUID and Sessions by SessionUID within the same ProjectUID,
then reconcile revisions and freshness; never simply add overlapping totals. Show
unknown populations separately from zeros. Inspect failed filters and request errors
before calling a result empty. Never make a failed query appear successful by
silently removing a filter.

Answer from verified available evidence while initial history fills; do not wait for
the entire cache when it already covers the question. Preserve maintenance pauses and
report any narrower window or population that limits the answer.

## 3. Acquire only the missing evidence

Skip acquisition when the selected source already contains what the question needs.
For live reads, use the same authenticated connection established during Access.

| If a live read needs | Read | Then |
|---|---|---|
| Successive batches or changes to Events, Sessions, Devices or another supported entity | [Modified cursors](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-project-chat/modified-cursors.md) | Apply the exact cursor, batch and checkpoint contract |

Use the cursor guide for supported scans through either MCP or CLI. Cursors order
modifications, so establish the question's capture-time window and coverage separately.
Cache maintenance owns ongoing acquisition; a question fetches only the evidence it lacks.

For an unfamiliar source, inspect its schema and a small bounded sample first. Filter
and aggregate near the data, batch compatible reads and use the documented response
limits. A small sample helps interpret fields but cannot stand in for the complete
population. Keep enough source evidence to verify or revise the calculation.

A foreground agent may use temporary files for a bounded analysis or export, but does not
start a persistent synchronized database or background job. Explain the actual time and
coverage cost of a large question. Offer a narrower window only when it materially helps;
do not change the requested scope without saying so.

## 4. Calculate and verify the answer

Reuse a validated query when its quantity, population, time basis and interpretation
match. Run it with the current parameters and evidence. Relate Events and Sessions to
Devices by DeviceUID within ProjectUID; EventUID and SessionUID identify their individual
records. Use names and serial numbers for display. Check representative rows and the
edge cases that could change an unfamiliar result.

- A Note, its Event history and a Note's deletion are different objects.
- Full Session records may be corrected as a session completes; incremental upserts
  replace earlier estimates. `_session.qo` contains a different set of fields.
- Equal measurement values are not duplicate packets. Use EventUID and, when
  investigating replay, the actual source/session/sequence/time mechanism.
- For current state, select the latest record first, **then** apply validity; otherwise
  an old good measurement can conceal a new fault. A period statistic has a different
  validity rule and must not blindly inherit this ordering.
- Separate per-note reduction (e.g. maximum of three phase currents) from reduction
  over devices or a period. State the counting unit and denominator.
- A span between two observed events is not automatically time in a state. Silence,
  sparse samples, offline buffering, clock errors and late arrival affect that inference.
- Historical fleet membership, configuration and serial-to-device binding may differ
  from today's snapshot. Mark unavailable historical bindings instead of projecting
  today's state backward.

Follow recipe status: answerable, conditional or blocked. An assumption stays visible in
the answer and cannot become a silent alarm threshold. If two sources disagree, check
version, cohort, time basis and completeness before treating the disagreement as a product
contradiction. Preserve a useful correction; don't overwrite shared knowledge merely
because a new sample looks different.

## 5. Return the result and retain the calculation

Return the calculated result, units, population, time window, freshness and evidence
limits to the caller. In the Chat loop, continue at Step 3 to explain it. Keep the
calculation and meaning revisions for follow-ups and possible monitors.

On a Cloud Computer, record the fields and query plan actually used, including live or
Repository reads that could benefit from a future cached projection. Maintenance uses
this workload history; answering a question does not require running maintenance now.
