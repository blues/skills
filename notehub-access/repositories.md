# Repositories and analytical queries

Use a repository for historical data beyond a cloud cache, or fast Q&A when no suitable
cache exists. Discover the repository tools the **selected authenticated connection**
actually offers. In MCP use the server's search function and inspect the returned tool
schemas. The public OpenAPI may contain repository-related types without exposing the
server's analytical query functions. Do not invent operation names or assume CLI parity.

## Choose the relevant reference

Read the rows that apply to the current step; reuse context already established.

| If you need to | Read | Then |
|---|---|---|
| Resolve the product meaning needed by the query | [Reading product knowledge](../notehub-project-chat/product-knowledge.md) | Load the relevant recipe and meaning dependencies |
| Interpret native fields in the dataset | [Platform](../notehub-platform/SKILL.md) | Verify the data surface and its semantics |
| Choose live evidence when repository coverage is unsuitable | [Query planning](../notehub-project-chat/querying.md) | Use the selected transport and disclose remaining coverage limits |

Bind the repository/dataset to the verified project. Inspect its schema and bounded
sample before writing a query: Event envelopes, flattened bodies, full Sessions and
transformed views are different evidence. Check ingestion freshness, time coverage,
filters, deletions and any retention limits. A dataset name is not a completeness claim.

Build typed, bounded queries using product recipes and native platform definitions.
Use stable device IDs for joins and serials/names for display. Partition by project and
account; use UTC unless the user specifies another analytical timezone. Parameterize
inputs, inspect a small result and aggregate server-side before returning large data.
Show missing populations or uncovered periods separately. A query returning zero rows
is meaningful only when the dataset covers the requested scope and time.

Repository provisioning, transformations and replacement are writes, not prerequisites
for answering. Do not create or replace a dataset just to discover whether one exists.
If suitable repository capability is absent, use bounded live API calls through the same
transport and state any coverage limitation. Don't seek an `X-Repository-Key`, database
password or raw API token. Fresh cloud caches take priority when they fully support the
question; direct API reads still serve current state and requested operational changes.
