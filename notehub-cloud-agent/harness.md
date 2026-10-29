# Verify the harness before promising ongoing work

Use this during Access's Cloud Computer preparation or step 1 of project Setup. Record
the result privately and return to the caller. Before project selection, retain an
environment capability record so Access can create the Notehub Cache Database and register
the Notehub Cache Maintenance Task;
bind individual project cache state on first use.

## Choose the relevant reference

Read the rows that apply to the current step; reuse context already established.

| If you need to | Read | Then |
|---|---|---|
| Record verified capabilities and their actual invocation paths | [Database state](workspace.md) | Keep durable nonsecret configuration and receipts |
| Resolve project scope for a project-bound capability check | [Project Selection](../notehub-access/projects.md) | Return with the verified project; general environment preparation can use an account-level read |
| Create the cache database and maintenance task during Access | [Cache preparation](cache-preparation.md) | Complete database installation and scheduled execution checks |
| Configure project-specific capabilities requested by Chat | [Setup](setup.md) | Resume the first incomplete step using existing context |

Establish this capability record with actual tool documentation and a harmless read/test:

| Capability | Evidence needed |
|---|---|
| Durable storage | The user's standard shared workspace where available, with verified visibility across bots/channels and persistence beyond the creating bot's deletion; otherwise record the supported storage scope |
| Shared memory and task ownership | A supported way to publish the cache location to every bot, and whether a shared/user-level maintenance task can survive deletion of its creating bot |
| Unattended authenticated acquisition | A successful read from a scheduled execution environment using the already selected Notehub identity/hub |
| Scheduler | A real job API with cadence/timezone, next-run status, failure reporting, cancellation and overlap control |
| Notehub Cache Database | Current implementation: DuckDB installed in the worker environment, a persistent scratch write/read/reopen check, recorded interpreter/version and bounded resource limits |
| Delivery | A supported user-facing destination and receipt/failure signal; recipients and external sends within the user's authorization |

MCP: the harness must supply a documented scheduled agent/tool context or programmatic
bridge to the same authorized MCP connection. A scheduled agent with MCP and workspace
tools is a usable path, with its model-execution cost recorded; prefer a deterministic
worker bridge where available. A tool available in conversation is not automatically an executable shell
command or a background credential. Don't export MCP credentials, invent SDK methods or
silently add a CLI login. CLI: run `notehub --help-agent`, verify `auth status --json` in
the job's OS environment, and test its normal stored-login refresh. Don't copy tokens.

Partial capability is useful. Durable acquisition plus scheduling can maintain a cache
without notifications. Delivery without unattended acquisition cannot make a live monitor.
A foreground session can answer robustly with repositories or bounded API reads, but
must not promise persistent cache refresh or periodic engagement. State the actual missing
capability and continue the useful supported work. Do not require a cloud migration just
to answer a question.

Return the observed capabilities, invocation paths and gaps to the caller. This capability
check is one part of preparation: Access must continue through Cache preparation to
create and verify the database and maintenance task. Project Setup later reuses that
infrastructure for enrolled projects and requested monitors.
