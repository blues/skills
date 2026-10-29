# Query planning and honest coverage

Choose the source by the evidence needed, not by a fixed preference for one query tool.

| Situation | Source | Read when the source or contract is unfamiliar |
|---|---|---|
| Cloud cache supports the fields, population, window and freshness | Query that cache | [Cache coverage](../notehub-cloud-agent/cache.md) |
| History extends beyond cache, or a foreground agent needs analytical queries | Suitable repository through the chosen connection | [Repositories](../notehub-access/repositories.md) |
| Missing field, current resource state, bounded investigation | Live API through that same connection | [API discovery](../notehub-access/api.md) |
| Read successive batches or changes to Events, Sessions, Devices or another supported entity | Modified-cursor reads through the selected connection | [Modified cursors](modified-cursors.md) |
| Improve repeated or slow cloud queries | Field-usage records and measured query plans | [Query performance](../notehub-cloud-agent/query-performance.md) |
| No source covers the request | Explain the missing evidence and offer the useful bounded answer | No additional reference |

Check coverage **before** querying: hub/project binding, entity, filters, complete scan,
last successful poll, capture-time window, ingestion delays and deletion handling. Caches
can be fresh but incomplete, or complete but stale. A repository can contain Events but
not full Sessions. If combining sources, reconcile overlapping stable identities and
revision/freshness rules; never simply add overlapping totals. Show unknown populations
separately from zeros. Inspect failed filters and request errors before calling a result
empty. Never make a failed query appear successful by silently removing a filter.

## Read changes with modified cursors

Use modified cursors for supported live scans and incremental refreshes, through either
MCP or CLI. Read batches using the exact returned `next_cursor` and continue while
`has_more` is true. Cloud cache workers save progress and resume it at each poll; a
foreground analysis stops after its bounded task. These cursors order modifications,
so establish the question's capture-time window and coverage separately. Use the shared
cursor guide above for the complete contract.

## Interpret before aggregating

| If you need to | Read | Then |
|---|---|---|
| Resolve custom field meaning or recipe dependencies | [Reading product knowledge](product-knowledge.md) | Apply the owning product definitions and recipe status |
| Interpret native records or system fields | [Platform](../notehub-platform/SKILL.md) | Choose the reference for the actual data surface |

Use the project's confirmed vocabulary and field inventory, including units, scaling,
array positions, sentinels, missing/zero distinctions, conditional emission and firmware
scope. For native records, use the platform route above. In particular:

- A Note, its Event history and a Note's deletion are different objects.
- Full Session records may be corrected as a session completes; incremental upserts
  replace earlier estimates. `_session.qo` contains a different set of fields.
- Equal measurement values are not duplicate packets. Use Event identity and, when
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

## Keep acquisition proportionate

Survey the schema and fetch a small bounded sample first. Aggregate in the repository or
cache where possible. For a large authorized read, use the operation's documented limits
and modified cursors, bounded retries and durable checkpoints if in cloud mode. Bulk-load
JSON into DuckDB; don't use one network request or one SQL insertion per row. Keep raw
records alongside typed per-Notefile views so revised product rules can be recomputed.
Avoid a single generic table that silently coerces unrelated body fields into one type.
Record the fields and query plan actually used, including live/repository reads that
could benefit from a future cached projection. Feed this usage into daily cloud maintenance.

A foreground agent may use temporary files for a bounded analysis or export, but does not
start a persistent synchronized database or background job. Explain the actual time and
coverage cost of a large question. Offer a narrower window only when it materially helps;
do not change the requested scope without saying so.

The answer should make the decision easier: use the user's units and reasonable precision,
compare against a confirmed baseline when useful, name the as-of time for changing state,
and say what remains unknown. Do not dump every row just because a tool returned it.
