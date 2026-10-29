# Authenticated access to a Project Repository

This is the MCP interface reference used by the Project Repository cookbook. It requires
MCP access and a Repository enabled for the selected Project. Reuse the existing
connection; authentication stays inside MCP.

| For the question workflow | Procedure | Result |
|---|---|---|
| Select tables, establish meaning and coverage, then calculate an answer | [Query the Project Repository](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-project-repository/SKILL.md) | Return a verified result for the requested population and history |

## Discover the supported functions

Use MCP search and read the returned schemas for Repository functions. These can include
`list_repositories`, `describe_schema`, `get_table_sample` and `query_repository`; use
what the connected server exposes. These functions may exist beyond the public OpenAPI
paths. Reuse interface definitions already discovered while they remain valid.

Match the selected ProjectUID to the Repository's `project_uids`, and retain RepositoryUID
from `uid`. Pass it explicitly under the tool's actual parameter name, such as
`repository_uid`. Keep any FleetUID restrictions with that scope; do not rely on a tool
defaulting to the first Repository.

## Apply the query engine's contract

Inspect actual schemas, supported syntax and response metadata. When MCP exposes
ClickHouse SQL, these details matter:

| Detail | Apply |
|---|---|
| Table `time` column | Read its projection's time expression. Use efficient time bounds without substituting upload time for the requested measurement time |
| Generated `location` Point, when present | Its coordinate order is longitude, latitude; verify the actual columns before spatial calculations |
| Loading and optimization state | Account for unfinished backfill or replay deduplication before trusting counts |
| `statistics.truncated`, when returned | Check it before treating returned detail as complete; aggregate or split the requested range when needed |

“Recent” retains the user's intended wall-clock window. Do not shift it to the latest
available observations to hide ingestion lag. Follow the cookbook for full-history
coverage and interpretation of transformed rows.

The existing Repository supplies read queries. Provisioning, replacing tables or changing
projections is separate requested work. If query access is unavailable, return that
limitation to the caller rather than obtaining a second connection, a Repository key or
raw credentials.
