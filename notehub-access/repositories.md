# Authenticated access to a Project Repository

This is the MCP interface reference used by the Project Repository cookbook. It requires
MCP access and a Repository enabled for the selected Project. Reuse the existing
connection; authentication stays inside MCP.

| For the question workflow | Procedure | Result |
|---|---|---|
| Select tables, establish meaning and coverage, then calculate an answer | [Query the Project Repository](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-project-repository/SKILL.md) | Return a verified result for the requested population and history |

## Discover the supported functions

Use MCP search and read the returned schemas for Repository functions. These can include
`list_repositories`, `describe_schema` (tables and their columns), `get_repository_dataset`
(one dataset's definition, meaning its Notefiles and projection, and its load status),
`get_table_sample` and `query_repository`; use what the connected server exposes. These
functions may exist beyond the public OpenAPI paths. Reuse interface definitions already
discovered while they remain valid.

A Repository covers either a whole Project (`project_uids` holds its ProjectUID) or some
of its Fleets (`project_uids` is empty and `fleet_uids` holds FleetUIDs). Match the
selected ProjectUID against `project_uids` and the Project's FleetUIDs from `GetFleets`
against `fleet_uids`. When several match, keep each one's fleet scope and choose the one
that covers the devices in question. Retain RepositoryUID from `uid` and pass it
explicitly under the tool's actual parameter name, such as `repositoryUID`; do not rely on
a tool defaulting to the first Repository. Each table is one dataset of the same name:
pass that name as `dataset_name` to read the definition and load status the cookbook's
Step 2 checks.

## Apply the query engine's contract

Inspect actual schemas, supported syntax and response metadata. When MCP exposes
ClickHouse SQL, these details matter:

| Detail | Apply |
|---|---|
| Table `time` column | Read its projection's time expression. Use efficient time bounds without substituting upload time for the requested measurement time |
| Generated `location` Point, when present | Its coordinate order is longitude, latitude; verify the actual columns before spatial calculations |
| Load status from `get_repository_dataset` | Only `ready` or `ready_optimized` means fully loaded; under `loading` or `not_loading` the table is partial or empty. Replayed rows can repeat until background merges remove them, even at `ready_optimized`: for an exact count, count distinct Event identifiers or query with `FINAL` |
| `truncated` on the first (status) line of a `query_repository` result | If true, a cap (50 MB) cut the result; aggregate, add a `LIMIT` or split the requested range before treating the rows as complete |

“Recent” retains the user's intended wall-clock window. Do not shift it to the latest
available observations to hide ingestion lag. Follow the cookbook for full-history
coverage and interpretation of transformed rows.

The existing Repository supplies read queries. Provisioning, replacing tables or changing
projections is separate requested work. If query access is unavailable, return that
limitation to the caller rather than obtaining a second connection, a Repository key or
raw credentials.
