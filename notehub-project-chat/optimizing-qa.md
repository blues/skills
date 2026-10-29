# Make Project Chat answers fast and reliable

Optimize the path from the user's question to verified evidence. Reuse access, project
context, relevant platform definitions and current product knowledge. Preserve meaning,
coverage and freshness while reducing repeated setup, downloads and computation.

## Choose the optimization that fits the environment

| If you need to | Read | Then |
|---|---|---|
| Choose evidence for the current question | [Query planning](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-project-chat/querying.md) | Use a sufficient cache, repository or bounded live read |
| Poll for changes in Events, Sessions, Devices, and more | [Modified cursors](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-project-chat/modified-cursors.md) | Resume saved progress and retrieve changes in batches |
| Resolve native field meaning | [Platform](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-platform/SKILL.md) | Load only the relevant native reference |
| Resolve custom product meaning | [Reading product knowledge](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-project-chat/product-knowledge.md) | Reuse the catalog and fetched documents; read only missing or changed content and consequential dependencies |
| Prepare reusable calculations for repeated questions | [Prepared question plans](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-cloud-agent/questions.md) | Bind parameters, meanings, evidence requirements and result shape |
| Create cache infrastructure or enroll another project on the Cloud Computer | [Cache preparation](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-cloud-agent/cache-preparation.md) | Reuse the Notehub Cache Database and Notehub Cache Maintenance Task |
| Optimize recurring queries using actual field usage | [Query performance](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-cloud-agent/query-performance.md) | Record queried fields and review schemas, typed projections, physical ordering and indexes daily |
| Find historical or analytical data outside available cache coverage | [Repositories](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-access/repositories.md) | Verify the actual schema, scope and freshness before aggregating |

## Reuse the work, refresh the evidence

1. Recognize the intent and reuse an existing validated plan when its quantity,
   population, time basis and validity rules match. New wording need not mean a new
   query; superficially similar questions can still need different calculations.
2. Retain a compact catalog of useful questions, their platform/product dependencies,
   parameters, required fields, freshness, coverage and verified query/script. A native
   question can use platform definitions without requiring a trained product recipe.
3. Run computation near the data: aggregate/filter in DuckDB or a suitable repository,
   select needed fields and return a compact result. Use typed views derived from
   verified meanings and retain raw records for later reinterpretation. Batch compatible
   reads through the chosen transport instead of making a separate call for each row.
4. Reuse plans, not old numeric answers. Check coverage and last successful acquisition,
   run the calculation with the new parameters, and return current evidence. Refresh
   missing/stale inputs through bounded reads when needed.
5. When a product revision, schema or relevant platform interpretation changes, identify
   affected plans and rebuild/recheck their derivations. Preserve unrelated working
   plans. An answer can be fast and wrong if its meaning changed beneath it.
6. Record fields used and their roles, query frequency and observed performance in the
   private workload history. Daily maintenance uses this evidence to improve physical
   layout and indexes; a field merely appearing in a schema does not make it a useful index.

## Cloud Computer and foreground execution

Access prepares the Notehub Cache Database and Notehub Cache Maintenance Task; DuckDB
is the current database of choice. Project Selection
enrolls each project the user first references, including during Training. Reuse this
infrastructure and honor paused work. Maintain the configured baseline: all Devices and Schemas,
recent Events and full per-device Sessions, with 15-minute modified-cursor polls where
supported, at least 30 days of recent history and pruning only after 45 days. Use the
cloud setup and cache contracts for exact checkpoint, concurrency and coverage behavior.
Refresh schema snapshots daily and review actual query usage to choose useful projections
and indexes. Verify maintenance coverage for every enrolled project. A monitor needing faster evidence
requires appropriately faster or bounded live acquisition.

Answer a supported first question while the larger cache initializes. State its actual
window and population; do not wait for every historical record or claim completeness
from a fresh but partial cache. Use a repository or live fallback for uncovered evidence.

In a foreground session, reuse knowledge and plans within the available session and
use repository aggregation or bounded live reads. Temporary analysis files are fine;
do not set up a persistent synchronized cache or promise jobs after the session ends.

## Carry a useful answer forward

Keep the calculation, filters, units, time window, result shape and meaning references
available for the next turn. If the person says “notify me when this changes,” those
become the basis of the monitor rather than a loosely related new query.

| If the next step is | Read | Then |
|---|---|---|
| Suggest or create a periodic follow-up | [Monitors](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-project-chat/monitors.md) | Reuse the answered question and define meaningful change and delivery |
| Provide ongoing status or contextual updates | [Engagement](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-project-chat/engagement.md) | Use current monitor/query results and the user's interests |
