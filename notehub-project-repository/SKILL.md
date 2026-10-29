---
name: notehub-project-repository
description: Query a Project Repository through MCP when the Project has a Repository enabled and a question concerns the whole dataset or no Notehub Cache is available. Introspect the Repository through MCP and combine Notehub Skills with Product Skills to query the appropriate tables, fields and history.
---

# Query the Project Repository

A **Project Repository** is Notehub's historical storage and analytical query capability
for a Project. It holds archived Events and exposes queryable tables derived from them.
Use it to answer questions over the full available history, beyond the period retained
in a Notehub Cache. It is sometimes called a “Project Dataset”; use **Project Repository**
for this capability. The API also uses `DataSet` for a particular table and its projection
within a Repository; that table is distinct from the Repository as a whole.

**This capability is available only when MCP is the access method AND the Project has
a Repository enabled.** Use it when the question concerns the whole dataset, including
100% of historical data, or when no Notehub Cache is available. It also supplies history
beyond the cache's coverage. **Repository queries are slower than querying a local
Notehub Cache, but faster than querying the Notehub APIs directly.**

This is part of Project Chat. Reuse the selected ProjectUID, authenticated connection,
Notehub Skills, Product Skills and question context. The Repository is a Notehub service
capability; querying it does not require the agent to run on a Cloud Computer.
Follow [working style](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub/working-style.md).

## 1. Verify MCP access and find the enabled Repository

First verify that the selected access method is MCP. If it is the CLI, this query
capability is unavailable; use the cache or ordinary Notehub APIs for supported work.
Keep the established access method unless the user chooses to change it.

With MCP, reuse the Repository information retained during Project Selection, or use
its Repository read/list capability. Verify that a Repository is enabled for this
Project. Match the returned `project_uids` to the selected ProjectUID and retain the
RepositoryUID returned in `uid`. Pass it explicitly
when querying; do not rely on a tool defaulting to the first available Repository.

| If you need to | Read | Then |
|---|---|---|
| Discover Repository access and query functions through MCP | [Repository access](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-access/repositories.md) | Verify the available read operations and keep authentication inside MCP |

Distinguish an absent Repository, denied access, an unavailable query capability and a
Repository still loading. Ordinary questions use an existing Repository; provisioning
or replacing it is a separate requested operation.

**Complete when:** MCP access and the selected Project's enabled Repository are verified,
or their precise unavailability is known. If unavailable, return that limitation so the
caller can select another suitable source.

## 2. Inspect the tables, meaning and historical coverage

**Introspect the Repository through MCP before crafting the query.** Read its actual
table and column schemas, then inspect a small sample when a table is unfamiliar.
Retain exact column names, types and nullability. Use both the loaded **Notehub Skills**
for native entities, metadata and time semantics, and **Product Skills** for customer
terminology, Event `body` meanings and their mapping to Repository columns.
A Notefile's JSON schema helps interpret its Event body; it does not by itself describe
a transformed Repository table.

For the tables needed by the question, establish:

| Check | What it establishes |
|---|---|
| ProjectUID, FleetUID or NotefileID filters | Which population and Events are included |
| Projection and row meaning | Whether a row represents an Event, an individual measurement expanded from an Event, or another derived quantity |
| Time expression and available timestamp columns | Whether the query measures capture, upload or another documented time; use the actual schema |
| Archive bounds, table bounds, loading progress and latest ingestion | Whether the complete requested history is queryable now |
| Available records and fields | Whether the question can use Event data, or needs full Session records or other entities from another source |

An `_session.qo` Event table is not the full Session API. Preserve that distinction when
combining Repository history with other evidence. Check configured filters, retention,
deletions and unfinished backfill before claiming that a result covers all history.

**Complete when:** the table definitions support the intended product meaning and the
available population, time range, freshness and any gaps are known.

## 3. Build the query from both kinds of knowledge

For a request about **100% of historical data**, use the entire applicable Repository
history. Do not inherit the Notehub Cache's recent-data window. Use explicit bounds
covering that history and preserve the user's chosen time basis and timezone; use UTC
when none was specified.

When using the Repository because no Notehub Cache is available, use the period the
question actually asks about. Repository availability does not make every question an
all-history query. Combine Notehub and product meanings to choose the appropriate
tables, filters, calculations and output.

Use the Repository's advanced query capabilities to compute near the data:

| Question needs | Query approach, when supported by the discovered interface |
|---|---|
| Totals, ranges, distributions or population comparisons | Server-side filtering, aggregates and grouping |
| Trends or changes over time | Time buckets and window calculations over the full requested period |
| Most recent state or distinct Devices | Latest-record or distinct calculations using DeviceUID and the product's validity rules |
| Geographic conditions combined with product measurements | Spatial filtering using the actual location columns and product measurement definitions |
| Relationships across tables | Joins on verified ProjectUID, DeviceUID, EventUID or SessionUID relationships, accounting for each table's row meaning |

Use the query language and functions exposed through MCP. Parameterize
values where supported, validate identifiers against the schema, and select only needed
columns. Preserve units, scaling, missing values, sentinels and firmware-dependent rules.
Names and serial numbers are for display; DeviceUID identifies a Device across queries.

**Complete when:** the query expresses the requested quantity, population and full time
range, using verified fields and supported query syntax.

## 4. Execute and verify completeness

Run the read query through MCP. Inspect the result's column metadata,
errors, limits and truncation indicators. A row cap limits returned detail; it is not
permission to report a partial result as the complete historical answer. Prefer aggregates;
when necessary, query nonoverlapping time intervals spanning the entire requested period
and combine them according to the calculation's semantics.

Check representative results against source evidence when validating a new or changed
calculation; reuse checks whose assumptions still hold. Preserve distinct Events by
EventUID and distinguish multiple measurements within one Event from replay duplicates.
Account for any Repository loading or deduplication state before trusting counts.
A zero result is meaningful only when the query covered the intended population and time.

| If verification needs | Read | Then |
|---|---|---|
| Native Event, Session, location or Device meanings | [Platform](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-platform/SKILL.md) | Apply the correct record and time semantics alongside product rules |

**Complete when:** the result is verified and complete for the stated scope, or any
remaining coverage or execution limitation is explicit.

## 5. Return the answer and reusable calculation

Return the result, units, ProjectUID, population, time window, evidence freshness and
coverage limits to the calling procedure. If entered from the question procedure, resume
there to combine or verify evidence; if entered directly from Chat, continue at Step 3
of its loop. State whether the result covers the full requested history.

Retain the query, RepositoryUID, table definitions and product-meaning revisions for
follow-ups. A later monitor can reuse the calculation while applying its own cadence and
freshness requirements. Keep Repository query access inside MCP; no database
password, API token or Repository key belongs in the answer or working records.

**Complete when:** the caller has the grounded historical answer and enough context to explain
it, investigate a follow-up or define a requested monitor.
