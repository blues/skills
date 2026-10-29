# Create or verify the Notehub Cache Database

This is Cloud Computer cookbook Step 2. Start with the runtime, storage location and
access mechanisms verified in Step 1. Complete the steps below for one shared DuckDB
cache. Reuse a working database and preserve its data and settings. Task registration
is the next cookbook step; project enrollment happens during Project Selection.

## 1. Find the existing cache or choose its shared location

Check shared memory, the standard shared workspace and existing task definitions. If
there is a cache for this Notehub account, read its README and use it. Otherwise choose
a durable location all of the user's bots would naturally find, such as a
`notehub-cache/` folder in the shared workspace. Coordinate creation so two bots cannot
create separate caches at once.

Where the harness supports sharing, keep the folder outside any individual bot's
private or temporary workspace and verify it survives deletion of that bot. Use the
supported storage scope when sharing is unavailable, and record that limitation.

**Complete when:** the database location and the bots/workers that can access it are known.

## 2. Install or verify DuckDB in the worker runtime

Use the interpreter and dependency mechanism chosen in Step 1. For Python, install the
`duckdb` package there if absent. Record the runtime and DuckDB version. Apply memory,
thread and temporary-storage limits appropriate to the available resources.

| If installation help is needed | Resource | Then |
|---|---|---|
| Install the Python package | [DuckDB installation](https://duckdb.org/docs/2.0/clients/python/overview) | Verify it imports in the actual worker environment |

**Complete when:** the worker can use DuckDB, including after its runtime is restored.

## 3. Create or verify the database and its access method

Choose one access method for every bot, worker and query: a single owning process,
serialized database opens, or published read snapshots. Reuse the existing method if
sound. Shared filesystem access alone does not coordinate independent DuckDB processes.

Create the database and a minimal project registry if missing. Keep it empty until a
project is selected. Choose the tables needed for native data and progress, preserving
existing contents and private state; add monitor and query records as those features
are used. Do not create a parallel directory of configuration files for database state.

| For this step | Resource | Then |
|---|---|---|
| Choose access, keys and minimal state | [Cache model, Steps 1–2](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-cloud-agent/cache.md#1-choose-and-verify-database-access) | Apply the suggested tables and integrity rules needed now; later sections support ingestion |

**Complete when:** the database opens through the agreed method and can store project scope.

## 4. Verify persistence and worker access

From the actual worker environment, write and read a disposable test record, close and
reopen the database, and verify it remains. Remove the test record. Verify other bots
or query workers use the same access method; keep tests separate from customer data.
Reuse a verified background Notehub read, or leave that connection check for the first
maintenance-task run. Record database and connection verification separately.

**Complete when:** worker read/write/reopen succeeds and any remaining access gap is explicit.

## 5. Make the cache discoverable

Write or update a short README beside the database covering:

- What each table contains, its suggested keys and project scope.
- How another bot opens and queries it safely.
- The maintenance routine and cadence, once configured in the next step.
- Timestamped freshness/coverage and how to read current freshness from the database.

Save a shared memory fact, where supported, identifying the actual folder and telling
all bots to read its README before querying. Verify it is shared rather than private
to the creating bot. Keep the README current as tables, tasks and freshness change.

**Complete when:** another bot can find the cache and learn how to use it without coaching.

Return the database location, runtime and verification result to the cookbook. Continue
with Step 3 to create or verify its maintenance task; no project data needs to be loaded yet.
