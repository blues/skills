# Prepare the Notehub Cache Database and Notehub Cache Maintenance Task

Cloud Computer access preparation has two central deliverables:

| Deliverable | Purpose |
|---|---|
| Notehub Cache Database | Durable, project-scoped data organized for efficient queries. The current database of choice is DuckDB. |
| Notehub Cache Maintenance Task | Recurring work that fills and refreshes the database, purges eligible history, refreshes schemas and optimizes queries. |

Complete both during Access so every project the user later works with can gain a
maintained cache on first use, in Chat or Training. Reuse the selected MCP or CLI
connection and existing infrastructure; retain pauses, limits and chosen schedules.

## Choose the relevant reference

| If you need to | Read | Then |
|---|---|---|
| Verify durable execution, storage or the authenticated background path | [Harness capabilities](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-cloud-agent/harness.md) | Establish the actual environment before registering the task |
| Decide what must survive between runs | [Database state](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-cloud-agent/workspace.md) | Store the minimum needed to resume work; choose the representation |
| Define tables, identities, bootstrap and retention | [Cache](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-cloud-agent/cache.md) | Preserve project boundaries, raw evidence and exact checkpoints |
| Implement the maintenance task's fill and daily review phases | [Job contracts](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-cloud-agent/jobs.md) | Connect acquisition, purge, schema refresh and optimization to real scheduled execution |
| Keep schemas and query performance current | [Query performance](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-cloud-agent/query-performance.md) | Record queried fields and use daily schema and workload review to optimize |
| Prepare product-specific questions or requested monitors | [Project setup](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-cloud-agent/setup.md) | Build on the installed runtime and enrolled cache |

## Find a shared home for one Notehub Cache

Some Cloud Computers support multiple bots or channels for concurrent conversations with
the same user, and may share a filesystem and memory among them. Where sharing is
supported, create and maintain **one Notehub Cache and one maintenance task across those
bots**, within the verified Notehub scope. A new bot or channel reuses them.

Before creating anything, check shared memory, the harness's normal shared workspace and
existing scheduled tasks for the cache. Read its README and reuse a matching cache and
task. Coordinate creation through the shared access/locking mechanism so two bots starting
together do not each create a database or poller.

Choose the harness's standard durable shared workspace: a place all bots would naturally
look first without user coaching. A clearly named `notehub-cache/` folder there is a useful
default. Keep the database outside any individual bot's private or temporary workspace,
and verify from the harness's storage lifecycle that the folder survives deletion of
the bot that created it. Use shared or user-level ownership for the maintenance task
where supported, so its execution
can survive too. All bots must use the same database access design; shared storage alone
does not permit independent DuckDB readers and writers to bypass concurrency rules.

After creating or reusing the database, write or update a short `README.md` in its folder covering:

- What each table contains, its keys and project scope.
- How often it is refreshed, and the maintenance routine/task responsible.
- How fresh the data is: timestamped last successful refreshes, coverage or gaps, and
  how to query the live freshness records.
- How another bot should open/query the database using the chosen access mechanism.

Keep the README current as tables, routines and freshness change. The database remains
the source of live state; a timestamped README summary is not proof that subsequent runs
succeeded.

Then save a short fact using the harness's **shared memory visible to every bot**, for example:

> The Notehub Cache exists at `<shared-folder>/notehub-cache/`. Read `README.md` in that
> folder before querying, and reuse its database and maintenance task.

Use the actual folder path and verify that the fact is shared, rather than saved only to
the creating bot's private memory. Other bots should discover this fact and read the
README before querying or preparing another cache. If shared storage, shared memory or
task ownership is unavailable, use the supported facilities and state the specific limit;
do not invent a shared path or claim that bot deletion cannot affect the task.

## Create and verify the Notehub Cache Database

Use DuckDB as the current implementation. The database's purpose and scope remain the
same regardless of engine; the installation, concurrency and indexing guidance here
describes DuckDB specifically.

1. Reuse or choose the shared database location above and the runtime used by scheduled jobs. Install
   DuckDB in that environment, using its managed dependency facilities. With Python,
   install the `duckdb` package through that interpreter's package manager; follow the
   [DuckDB Python installation guide](https://duckdb.org/docs/2.0/clients/python/overview).
   Installing only in the interactive session is insufficient. Record the interpreter,
   DuckDB version, database paths and worker invocation in the registry.
2. Create the database and a minimal project registry; add other state as work needs it.
   Start with an empty registry when the user has not referenced a project yet. All bots
   use the same cache and project registry; choose its physical organization for the workload
   and preserve it across restarts. Store only the minimum database locator and
   safe worker invocation in the harness's task configuration; no parallel directory/JSON
   registry is required. Configure sensible memory, thread and temporary storage limits.
3. Establish the database access design: one owning process, serialized database opens,
   or published read snapshots. All fill, purge, query and optimization paths must obey
   it. Use the concurrency contract in Cache; a writer-only lock is insufficient when
   independent query processes also open the database.
4. From the actual worker environment, perform a scratch transaction, query its result,
   close and reopen the database, and verify persistence. Keep the smoke test separate
   from customer tables. Verify the same authenticated background connection with a
   harmless read. Record both results; an import check alone is not readiness.

## Create and verify the Notehub Cache Maintenance Task

Reuse the shared maintenance task if present. Otherwise implement the Job contracts'
pseudocode using the harness's supported runtime and register the task with its real scheduler.
Include project dispatch, snapshots, retention and
due-work scheduling, and verify actual execution. A desired schedule alone is incomplete.

| Task phase | Default cadence | Work and completion evidence |
|---|---|---|
| Fill and refresh | Every 15 minutes | Read the current project registry, fill due active projects through modified cursors, commit batches and coverage, and retain a run receipt plus next run |
| Purge, schema refresh and optimization | Daily at a suitable local off-peak time | Purge eligible history, refresh each enrolled project's schemas and metadata, review queried fields and query performance, and reconcile task health |

By default, one recurring task wakes every 15 minutes, fills due projects and performs
daily work when its per-project due time arrives. It covers all enrolled projects for
the verified hub/account. If the harness or workload needs separate scheduled subtasks,
record their IDs and responsibilities under this maintenance task and keep the same
registry and enrollment behavior. Do not register overlapping work for the same scope.
Set timeouts, retry/backoff, overlap control and per-project budgets. One large or failing
project must not starve the others. No enrolled projects means a successful idle run,
not discovery and caching of every accessible project.

Run the task through its actual unattended path and exercise both phases. Verify scheduler IDs, next runs,
locking, registry reads and durable receipts. With no projects yet, verify idle behavior;
exercise real fill and daily review on the first enrolled project. Record database
readiness and maintenance-task readiness separately from each project's cache coverage.
Missing unattended authentication or scheduling
leaves that part incomplete; keep supported foreground work available and report the gap.
Update the cache README with the verified task identity, cadence and observed freshness,
and confirm the shared memory fact points to its folder.

## Enroll a project on first use

Whenever the user first references or starts working with a project on this Cloud
Computer, resolve its canonical UID and verify read access through Project Selection.
Then enroll it as part of that workflow, without waiting for a monitor request or a
separate caching request. Apply this to every project the user interacts with, including
projects used during Training. A name merely encountered in a listing, stored document
or another project's data is not a user request to enroll it.

1. Look up the durable hub/account/project binding. Reuse an existing active enrollment;
   preserve cursors, coverage, jobs and query history. A paused or explicitly excluded
   project remains so until the user changes that instruction.
2. For a new binding, retain enough project scope and status for the task to maintain its
   cache. Create data and progress records as bootstrap proceeds; query-usage records grow
   as queries run. Keep an inherited skill owner separate from the data scope; inheritance
   does not enroll the parent's data.
3. Queue an immediate bounded bootstrap: Projects/Devices metadata, recent Events, full
   per-device Sessions and the complete Notefile Schema snapshot. Read product skills
   through the shared product-reading workflow. Use the cache's history and retention
   rules, and record the first schema refresh plus its next daily due time.
4. Ensure the Notehub Cache Maintenance Task discovers this enrollment on its next run;
   it must not retain a fixed project list copied from its creation prompt. If using
   project-specific subtasks, create and verify only the missing ones.
5. Record initializing, ready, degraded or paused state per project. Answer a supported
   question or continue the interview while history loads. Switching the conversation to
   another project enrolls that project too; it does not stop or retarget earlier caches.

Return the runtime/enrollment receipts to the caller. The cache infrastructure does not
create user-facing monitors or messages; those follow the user's requested Chat workflow.
