---
name: notehub-cloud-agent
description: Prepare the Notehub Cache Database and Notehub Cache Maintenance Task during Access on the user's Cloud Computer, currently using DuckDB. Enroll projects on first use, maintain query-ready data and support Project Chat's Q&A, monitors, alerts and engagement.
---

# Cloud Computer support for Access and Project Chat

Access uses this skill to prepare the Notehub Cache Database and Notehub Cache Maintenance
Task. DuckDB is the current database of choice. Project Selection enrolls projects on
first use in Chat or Training. Project Chat
uses the maintained evidence and runs the user's periodic work. Keep the conversation and its question,
monitor or engagement intent in Chat; return there after completing the needed setup.
Prefer cloud mode where supported. Follow [working style](../notehub/working-style.md).

## Prepare the Cloud Computer for Access

**The two central deliverables are the Notehub Cache Database and the Notehub Cache
Maintenance Task.** Create the database for efficient queries, using DuckDB as the current
implementation. Verify durable read/write/reopen behavior and register the maintenance
task that fills, refreshes, purges and optimizes it. Prepare a project registry even when
no project has been selected yet. Reuse the selected MCP or CLI connection and verify
authenticated background reads without exporting credentials or establishing a second login.
Where the user's bots or channels share a workspace, reuse one shared cache and maintenance
task. Make them discoverable through the cache folder's README and a shared memory fact.

| If you need to | Read | Then |
|---|---|---|
| Create the Notehub Cache Database and Notehub Cache Maintenance Task | [Cache preparation](cache-preparation.md) | Verify database persistence, register the task and retain actual execution receipts |
| Locate or prepare the cache for multiple bots or channels | [Shared cache location](cache-preparation.md#find-a-shared-home-for-one-notehub-cache) | Reuse one durable shared cache; publish its README and shared memory location |
| Establish the environment's capabilities and invocation paths | [Harness capabilities](harness.md) | Record supported storage, runtime, unattended access, scheduling and delivery |
| Decide what to persist between runs | [Database state](workspace.md) | Keep state minimal and choose its organization to suit the work |

Before a project is selected, verify an appropriate harmless account-level read and the
maintenance task's idle behavior. On each project's first user reference, verify its binding and
enroll it for immediate bootstrap and ongoing maintenance. This includes projects used
during Training. Default fill runs are every 15 minutes; daily maintenance purges eligible
history, refreshes schemas and reviews queried fields for useful physical optimizations.
Return installation, registry and job receipts to Access. User-facing monitors and
notifications follow the requested Chat workflow.

## Choose your starting point

| Current situation | Read | Then |
|---|---|---|
| The user first references a project on the Cloud Computer | [Project enrollment](cache-preparation.md#enroll-a-project-on-first-use) | Verify its identity/read access and add it to the cache maintenance task's registry |
| An enrolled project needs Chat preparation | [Setup](setup.md) | Reuse its cache infrastructure and prepare useful queries and requested recurring work |
| Existing configured assistant, ordinary question | [Project Chat](../notehub-project-chat/SKILL.md) | Return with the workspace and prepared plans; answer from suitable evidence |
| New or changed monitor/alert | [Monitors](../notehub-project-chat/monitors.md), then [Activation](setup.md#activate-recurring-work) | Reuse setup, define the monitor, install it and verify it runs |
| Scheduled execution | [Job contracts](jobs.md) | Execute the relevant contract with the durable workspace |
| Failed job, changed meaning, pause, removal or handoff | [Lifecycle](lifecycle.md) | Apply the change and verify the resulting state |

## Read implementation guidance when needed

| If you need to | Read | Then |
|---|---|---|
| Establish a missing connection or execution environment | [Access](../notehub-access/SKILL.md) | Return to setup using the same verified connection |
| Resolve missing or changed project scope | [Project Selection](../notehub-access/projects.md) | Return with the verified data project and skill owner |
| Load or refresh custom product meaning | [Reading product knowledge](../notehub-project-chat/product-knowledge.md) | Ground the relevant questions in the project's recipes and dependencies |
| Interpret native records used by a question | [Platform](../notehub-platform/SKILL.md) | Apply the relevant native definitions |
| Design maintained evidence and coverage | [Cache](cache.md) | Implement its transaction, concurrency and retention requirements |
| Refresh schemas or improve query performance | [Query performance](query-performance.md) | Record queried fields and perform daily schema and measured optimization review |
| Poll for changes in Events, Sessions, Devices, and more | [Modified cursors](../notehub-project-chat/modified-cursors.md) | Retrieve changes in batches through the selected MCP or CLI connection |
| Implement incremental ingestion | [Batch ingestion pseudocode](jobs.md#batch-ingestion-pseudocode) | Adapt the algorithm to the harness and verify checkpoint/restart behavior |
| Implement monitor comparisons and delivery | [Monitor pseudocode](../notehub-project-chat/monitors.md#evaluation-and-delivery-pseudocode) | Validate the requested behavior using the harness's own tools |
| Configure useful proactive engagement | [Engagement](../notehub-project-chat/engagement.md) | Apply its relevance and timing rules during setup or scheduled execution |

Setup reuses Access and the relevant platform/product knowledge; it does not
start a second authentication flow. Verify durable storage, scheduled execution and an
unattended path through the **same selected MCP or CLI connection**. Add delivery where
notifications are authorized. A terminal or an interactive MCP tool alone proves neither.

Access creates the database and maintenance task; first-use project enrollment starts acquisition.
Respect user-requested pauses and schedules. Prepare
all Devices and Notefile Schemas, recent Events and **full per-device Sessions**
in DuckDB, polling every 15 minutes with modified cursors where documented. Maintain at
least 30 days of recent history and prune only after 45 days. Use the relevant meanings to
build typed views and useful query plans. Apply the cache and job guidance above.

Scripts handle acquisition, calculation and routine state changes; the model handles
interpretation and useful communication. User-requested monitors use the requested cadence
and verified freshness. Use harness-owned jobs, avoiding Notehub operations with `Monitors`
or `Alerts` in their operationId. Engagement uses three daily
local-time evaluation opportunities unless specified otherwise, with messages only when
useful. Keep personal preferences and delivery history in private harness state.

Complete supported setup within the user's request and existing authorization. If a
capability is missing, configure the supported parts, identify the precise gap and keep
foreground Q&A useful. Report separately what is ready, initializing or blocked; don't
call an untested schedule active or wait for every history row before answering a question.
Training keeps its dedicated interview workflow while its project's cache initializes;
neither Training nor initial Q&A waits for every historical record to load.
