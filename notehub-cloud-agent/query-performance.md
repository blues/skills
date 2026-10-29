# Maintain schemas and optimize the queries people actually use

Run this review daily for enrolled projects, and when a newly needed field, schema change
or slow query warrants it. It is part of the Notehub Cache Maintenance Task created during Access.
The agent chooses the physical layout using actual workload evidence, retaining correct
product meaning and the raw data needed to rebuild derived structures.

## Choose the relevant reference

| If you need to | Read | Then |
|---|---|---|
| Create the cache database and maintenance task | [Cache preparation](cache-preparation.md) | Verify both deliverables and enroll projects on first use |
| Choose entity identities and retention rules | [Cache](cache.md) | Keep records project-scoped and preserve raw evidence |
| Run the task's daily review phase | [Maintenance](jobs.md#maintenance) | Refresh schemas, purge eligible history and record optimization results |
| Interpret a queried product field | [Reading product knowledge](../notehub-project-chat/product-knowledge.md) | Apply units, validity and recipe dependencies before changing a projection |

## Keep a record of fields queried

Record usage when Q&A, training analysis, prepared questions or monitor evaluators query
project data and store the usage records in the Notehub Cache Database. Include bounded
live/repository analyses so future cache projections can
support those needs. Aggregate by project, entity/Notefile, field path, semantic version
and query fingerprint; do not retain unnecessary literal values or raw customer results.

| Record | Why it matters |
|---|---|
| Field path and source entity/Notefile | Identify the actual field, including native metadata and nested body paths |
| Use as filter, join, grouping, ordering, projection or aggregate | Distinguish a useful lookup or scan key from a field merely displayed |
| Query fingerprint/plan, execution count and last use | Find repeated Q&A and periodic monitor workloads |
| Time window, population and selectivity/cardinality estimates when available | Choose useful physical ordering and distinguish broad scans from selective lookups |
| Elapsed time, rows scanned/returned and available plan diagnostics | Identify and measure actual bottlenecks |
| Schema/product revision and applied layout/index changes | Recheck correctness after meaning changes and track the effects of optimization |

Use a wrapper around executed queries or the harness's supported query instrumentation;
update the usage record as part of execution, rather than expecting the daily review to
reconstruct it from conversation history. For generated SQL, record the resolved field
mapping when constructing the plan; do not guess it from arbitrary SQL using a fragile
regular expression. Bound telemetry retention and preserve aggregate usage history.

## Refresh each project's schemas daily

Fetch `GetNotefileSchemas` on enrollment and at least daily by default, using the selected
authenticated connection. A user's cadence takes precedence. Refresh sooner when a
question needs an unknown field or actual data contradicts the stored schema. Refresh
project metadata too; neither resource acquires an invented modified cursor.

Stage and validate the full schema response before replacing its snapshot. Store project
UID, Notefile ID, exact schema, fetch time and a local content hash; preserve relevant
schema history or change records. A failed or partial fetch leaves the prior snapshot
available with its real freshness. Compare added, removed and changed fields/types with
both the query-usage records and product definitions.

A schema describes structure, not the customer's units or business meaning. Keep raw
values and historical interpretations when fields disappear or change type. Mark affected
views, query plans and monitors for reconciliation; do not silently cast failures to zero,
rewrite shared product skills or discard old evidence. Build replacements in staging and
validate their results before switching readers.

## Choose and verify useful optimizations

Start with repeated or costly queries and their actual plans. Measure representative,
bounded read-only queries with `EXPLAIN` and, when useful, `EXPLAIN ANALYZE`; the latter
executes the query and reports operator timing. Compare equivalent results and resource
cost before and after changes. See [DuckDB query profiling](https://duckdb.org/docs/2.0/guides/meta/explain_analyze).

Promote frequently queried JSON fields into appropriately typed columns or derived tables
when repeated extraction is costly. Ordinary views organize SQL but do not materialize
results; persist a derived table when measurements justify it and keep it synchronized
with updates, tombstones, retention and meaning revisions. Choose ordering/partitioning
for the observed project, device, Notefile and time predicates rather than one universal
wide table for unrelated products.

DuckDB automatically maintains min-max indexes; physical ordering can improve their
ability to skip irrelevant data. Explicit ART indexes are useful for sufficiently
selective lookups and add memory and write costs. Composite identity constraints ensure
uniqueness but do not imply a fast index scan for analytical queries. Check index eligibility
for the installed version, and use explicit indexes where query plans and measurements
show a benefit. See [DuckDB indexing](https://duckdb.org/docs/2.0/guides/performance/indexing).

After substantial inserts, updates or purges, run `ANALYZE` when refreshing join-order
statistics will help. It updates optimizer statistics; it does not choose or create indexes.
See [DuckDB ANALYZE](https://duckdb.org/docs/2.0/sql/statements/analyze).

Keep the useful constraints needed for idempotent ingestion. Schedule costly rebuilds
under the same database access design as ingestion and readers, with bounded work and
recoverable staging. Retain the old valid layout until the replacement is verified;
defer work that would monopolize the runtime or exhaust temporary storage. Maintenance
may legitimately make no physical changes when current plans already perform well.

Record which schema/fields were reviewed, optimizations applied or deferred, before/after
results, coverage and meaning checks, and the next review time. Query speed never substitutes
for correct units, population, history coverage or freshness.
