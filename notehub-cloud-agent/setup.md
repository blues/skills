# Set up the assistant from start to working jobs

Use this to prepare an enrolled project's questions and requested recurring work.
Access creates the Notehub Cache Database and Notehub Cache Maintenance Task; Project Selection
establishes the binding and enrolls the project on first use. Initial product reading supplies the product context;
reuse it here and acquire further meaning only when the work needs it.
The maintenance task covers all enrolled projects, including those used during Training.
Find the existing cache through shared memory or the harness's normal shared workspace
where supported, read its README, and reuse it across the user's bots/channels.
For first project setup or an existing assistant, read its database state
and resume the first incomplete step; don't recreate jobs, rebuild good caches or repeat
login. Preserve user-requested pauses and schedules. Carry forward existing authorization;
ask only for missing decisions that change scope, recipients or meaning. User-facing
monitors and engagement are configured for the requested or already established scope.

The completion target is concrete: durable database state, the selected connection working
in background execution, product knowledge and useful question plans, maintained evidence,
registered jobs, verified runs, and a short readiness report. Use the harness's actual
tools to achieve it. These instructions describe what to configure, not a universal vendor API.

## 1. Recover context and verify the execution environment

| If you need to | Read | Then |
|---|---|---|
| Establish a missing connection or execution environment | [Access](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-access/SKILL.md) | Reuse one verified transport and environment |
| Resolve missing or changed project scope | [Project Selection](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-access/projects.md) | Retain the verified project binding and relevant permissions |
| Load or refresh meaning for custom product data | [Reading product knowledge](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-project-chat/product-knowledge.md) | Retain the index, relevant recipes and meaning dependencies |
| Interpret native records used by a question | [Platform](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-platform/SKILL.md) | Retain the applicable native definitions |
| Repair the cache database, maintenance task or project enrollment | [Cache preparation](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-cloud-agent/cache-preparation.md) | Complete only the missing infrastructure and preserve existing jobs/checkpoints |

Reuse completed access and product reading. Inspect the existing project and job registry
in the Notehub Cache Database. Verify project/account binding, storage persistence, the DuckDB runtime, scheduling,
overlap control and user-facing delivery as applicable. Record exact commands or tool
entrypoints and returned IDs so the next scheduled run can use them.
Reuse Access's database and maintenance task receipts; do not create another cache task
merely because the conversation switched projects.

When that environment or binding has changed, test one small read from the **background execution environment**. For CLI, verify
`auth status --json` there and reuse its stored login. For MCP, the harness may run a
scheduled agent with the same tools, or provide a supported programmatic bridge; use
what it actually exposes. Code mode alone is not a local Python client. If there is no
unattended path, mark acquisition blocked and continue supported foreground work rather
than quietly adding CLI credentials to an MCP workflow.

If the harness schedules only model turns, retain that limitation explicitly: let the
scheduled turn invoke deterministic acquisition/evaluation scripts and avoid reasoning
when results are unchanged. Prefer a script/worker scheduler when available. Never label
model-scheduled polling as token-free. Don't install a cron job in an ephemeral container
and assume it survives after the session.

## 2. Reuse and update database records

| If you need to | Read or use | Then |
|---|---|---|
| Decide what must survive between runs | [Database state](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-cloud-agent/workspace.md) | Keep only the state needed by the configured work |

Reuse the database and registry established during Access, preserving other projects,
settings and job IDs. Persist what the configured work needs to resume without this
conversation; let the agent choose the representation and add records as needed.
Keep credentials in the selected tool. Record field usage as queries run so daily
maintenance can adapt projections and indexes to the actual workload.

Resolve the user's timezone and delivery surface from existing harness/user context;
ask only if unknown and necessary for scheduling. Defaults for an authorized continuing
assistant are 15-minute cache polls, three local engagement checks at 09:00/13:00/18:00,
and daily maintenance at a suitable local off-peak time. Record chosen times and quiet-hour
preferences; a user-specified schedule wins. Analytical windows remain UTC unless requested
otherwise. A schedule is desired state until the harness returns a real job ID.

## 3. Verify the project's bootstrap and periodic maintenance

Use modified cursors to fetch changes to Devices, Events and full per-device Sessions
on each scheduled run. Each stream resumes its own saved checkpoint through the selected
MCP or CLI connection. Establish this raw-data acquisition before refining product-specific
queries and typed views; missing product meaning does not prevent caching native records.

| If you need to | Read | Then |
|---|---|---|
| Poll for changes in Events, Sessions, Devices, and more | [Modified cursors](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-project-chat/modified-cursors.md) | Drain successive batches, commit each returned cursor and resume it on the next poll |
| Implement the acquisition worker | [Acquisition contract](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-cloud-agent/jobs.md#acquisition) | Fetch through the selected transport and persist per-stream progress |
| Define coverage, concurrency and retention | [Cache](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-cloud-agent/cache.md) | Preserve raw evidence, exact checkpoints and honest coverage |
| Keep schema snapshots and query performance current | [Query performance](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-cloud-agent/query-performance.md) | Verify daily refresh and a working record of queried fields |
| Implement and validate batch ingestion | [Batch ingestion pseudocode](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-cloud-agent/jobs.md#batch-ingestion-pseudocode) | Use the harness's supported runtime and authenticated background path |
| Configure retention, schema refresh and optimization | [Maintenance contract](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-cloud-agent/jobs.md#maintenance) | Verify the task's daily review covers this project |

Verify enrollment has configured:

- A project metadata record and its database/worker binding.
- A complete Devices scan with cursor `0`, and incremental upserts/tombstones thereafter.
- Recent Events, with stable bootstrap bounds and a separate capture-time coverage record.
- Full Sessions per discovered device, including new-device enrollment and honest gaps for
  inaccessible/historical devices. Keep these distinct from `_session.qo` events.
- A complete Notefile Schemas snapshot on enrollment, refreshed daily and replaced
  atomically only after a successful fetch.
- A coherent product-skill snapshot/revision check; changed meaning invalidates dependent
  views, prepared queries and monitors until they are rechecked.

Fix initial stream definitions once and keep them stable; a rolling “now minus 45 days”
config changes scope and would reset a hashed checkpoint. Retention advances separately.
Run bounded batches, commit progress, and schedule continuation for a backlog. Register the
project in the existing maintenance task's registry; verify its fill/daily review schedule,
IDs, invocation, timezone, retry/overlap policy and next run. Create a missing task through
Cache preparation, rather than duplicating a functioning one. Use the cache's database access design for concurrent Q&A and polling: one owning
process, serialized opens, or published read snapshots. A writer lock alone does not
make a separate reader process compatible with an open write connection.

Invoke acquisition through its actual unattended execution path. Inspect committed
records, cursors and per-stream progress, then verify the registered job's next run.
Retain the scheduler ID and run receipt; a saved script or configuration alone is not
periodic caching. Reuse these receipts for a healthy existing setup. Large initial scans
can continue in bounded runs while Q&A uses the evidence already available.

Verify maintenance covers the 45-day purge boundary, protected ongoing
sessions/incidents, daily schema refresh, field-usage review, projection/index optimization
and state health. Don't delete unknown-time rows
or move a cursor during retention. Verify schemas, all-device inventory and per-device
Session progress separately from the Event stream. An initializing cache is useful only
for the coverage it has already verified; use live/repository evidence for the rest.

## 4. Prepare the question catalog before choosing projections

| If you need to | Read | Then |
|---|---|---|
| Turn platform/product questions into executable plans | [Prepared questions](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-cloud-agent/questions.md) | Record meaning, evidence requirements, query inputs and verification status |

Read the relevant platform definitions and, for custom product meaning, the product's
index, useful recipes and all their meaning dependencies. Select the questions the
person actually asks or the available knowledge supports; record conditional/blocked
questions honestly. Link each executable query to its definitions/recipe, units,
population, time basis, field paths, applicable revisions and freshness/coverage needs.

Build typed views for useful product fields, while retaining raw JSON and native IDs.
Don't invent a universal wide table or index every field without a question it helps.
Prepare at least a representative current-state question and a period/population question
when the product supports them. When one doesn't apply, record why. Missing product
semantics block that query or monitor, not the rest of setup.

## 5. Prove that a new question can be answered

| If you need to | Read | Then |
|---|---|---|
| Find evidence beyond the cache's history or fields | [Repositories](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-access/repositories.md) | Verify actual dataset coverage and freshness before querying |
| Choose among cache, repository and live evidence | [Query planning](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-project-chat/querying.md) | Answer only for the supported scope and show gaps |

Execute the prepared current-state and period/population queries that apply. For each,
compare one bounded case with source records through the same transport, then try a
new natural-language wording using only the product index, recipe and its dependencies.
Check units, latest-before-validity where relevant, zero/absence, population exclusions,
unknown counts, window boundaries, duplicates and source freshness. Record results and
limitations in the query plan. A passing SQL statement is not evidence of correct meaning.

Probe a repository when needed for longer history or an uncovered question. Inspect its
actual schema, coverage and freshness, and record a dated capability observation in the
database state. Don't provision a dataset just to establish readiness. If no repository
is available, record the supported live API path and its cost/coverage instead.

## Activate recurring work

| If you need to | Read | Then |
|---|---|---|
| Define a new or changed monitor/alert | [Monitors](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-project-chat/monitors.md) | Resolve the calculation, state transitions, freshness and delivery contract |
| Configure authorized proactive engagement | [Engagement](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-project-chat/engagement.md) | Set useful local-time opportunities and suppression rules |
| Supply durable execution inputs and completion behavior | [Job contracts](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-cloud-agent/jobs.md) | Apply the relevant evaluation, delivery or engagement contract |

For each requested monitor, complete the monitor definition using the
prepared question and applicable platform/product rules. Store the evaluator reference, state
and delivery settings in the database. A monitoring request is not fulfilled by a prose suggestion or a saved
configuration alone:

1. Resolve population, condition, units, window, cadence, freshness, action and destination
   from the request and verified knowledge. If a necessary threshold is unknown, keep
   that monitor a draft and ask the focused question; continue unrelated setup.
2. Run the evaluator against representative normal, active and unknown/stale evidence,
   and test opening, escalation, repeated condition, recovery and restart. Use synthetic
   cases where real evidence lacks a branch, clearly labeled. Do not change device
   readings or production settings to manufacture an alert.
3. Persist incident state and a pending notification transactionally. Exercise delivery
   in dry-run mode first. Verify the destination through the harness; a test notification
   is sent only if the setup request/standing scope covers it. A supported delivery
   receipt is evidence; a draft message is not. Keep delivery unverified if it cannot be
   tested within scope, without stopping the evaluator or pretending alerts are ready.
4. Register the monitor evaluator at the requested cadence and the delivery worker/hook.
   If the monitor is faster than the 15-minute cache or needs fresher data, fetch its
   bounded inputs live or increase acquisition frequency for those streams within scope.
   Running the same stale query more often does not improve detection latency.
5. Register engagement checks for authorized continuing assistance, using explicit or
   proportionate inferred interests. With no useful interests/candidates, checks stay
   silent. Preserve suppression, quiet hours, feedback and closeout history.
6. Invoke each registered job through the same execution path it will use unattended.
   Inspect run status, persisted state and next run; verify repeat invocation reuses
   checkpoints and suppresses duplicate messages. Keep real scheduler IDs and receipts.

Use the job contracts above to give each execution all of its inputs. Alerts are delivered
monitor findings, not a separate authentication system or a reason to use native Notehub
Monitors/Alerts APIs. Background jobs need their own failure status and a way to pause.

## Report verified readiness

Read the database records and scheduler back, then give a short product-specific account:

- Which questions work now, and any material window/coverage limitation.
- Cache state per project: ready or initializing, latest successful poll and known gaps.
- Notehub Cache Database readiness (currently DuckDB) and Notehub Cache Maintenance Task
  readiness; daily schema/purge/optimization status and
  confirmation that new project enrollments are included in scheduled runs.
- Active monitors: condition/population in the user's terms, cadence, destination and
  verified delivery status; list drafts or blocked items separately.
- Engagement times/timezone and the fact that quiet checks produce no message.
- How to pause/change the work, plus any unfinished dependency and next concrete step.

A setup receipt records query checks, snapshot/stream coverage, job IDs, last successful
runs, next runs, delivery verification and remaining gaps. If the harness supports all
required capabilities, complete these steps rather than stopping at a plan. If it does
not, describe the exact supported subset and keep Q&A available. Return to the user's
pending question once setup work relevant to it is complete.
