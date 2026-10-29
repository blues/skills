---
name: notehub-access
description: Establish or reuse secure Notehub access before project selection, Chat or Training. Install or connect one MCP or notehub CLI, authenticate and prepare the execution environment, including the user's Cloud Computer when available. Return the verified connection and capabilities to the caller.
---

# Prepare the connection and execution environment

Complete missing access prerequisites, then return to the step that called this skill.
Reuse verified setup from the same environment. Access supplies the connection and
execution capabilities; Project Selection binds the task to a project, and the entry
guide chooses Chat or Training. Follow [working style](../notehub/working-style.md).

## Connect and authenticate once

Choose one authenticated request proxy: prefer Notehub MCP; use the `notehub` CLI when
MCP is unavailable or the user chooses it.
Both own their credentials—the agent does not retrieve tokens or construct Authorization
headers. If the hub/account and invocation method are already verified, reuse them.

| If you need to | Read | Then |
|---|---|---|
| Register/connect MCP, authenticate or learn its invocation method | [Notehub MCP: an authenticated request proxy](mcp.md) | Use the harness's supported connection and browser authorization; verify access |
| Install or authenticate the CLI | [Notehub CLI: an authenticated request proxy](cli.md) | Follow official installation guidance and `notehub --help-agent`; complete its browser/polling login flow |

Respect a cancelled authorization. Change transport only for a concrete need or the
user's choice, explicitly verifying the intended hub/account. A project permission
failure is not a reason to seek another credential source. Do not repeat login or API
discovery at each skill boundary.

## Prepare the execution environment

When running on the user's Cloud Computer, prepare two central deliverables during Access:
**a Notehub Cache Database for efficient queries, and a Notehub Cache Maintenance Task**
to fill, refresh, purge and optimize it. The current database of choice is DuckDB.
Set up a durable project registry so every project the user first references can be
enrolled into maintained caching. Verify the same authenticated connection works unattended.
Daily maintenance refreshes project schemas and optimizes queries using a record of the
fields actually queried. Discover app/messaging delivery capabilities for later requested
monitors and engagement.

| Operating capability | Read | Then |
|---|---|---|
| The user has their own Cloud Computer, such as Grok Bot, Muse, Dots, or 'claws' | [Cloud Computer preparation](../notehub-cloud-agent/SKILL.md#prepare-the-cloud-computer-for-access) | Create and verify the Notehub Cache Database and Notehub Cache Maintenance Task, ready for first-use project enrollment |
| Foreground session without continuing execution | The selected transport guide above | Retain interactive access without persistent pollers |
| Environment already verified | No repeated preparation | Reuse its actual workspace/invocation paths and capability record |

A terminal alone does not establish continuing execution. Record any missing cloud
capability and keep supported work available; a missing delivery channel need not block
interactive access.

## Return the access context

| Context | What to retain |
|---|---|
| Connection | Selected MCP connection or CLI environment, verified hub/account identity and invocation method |
| Environment | Foreground or Cloud Computer; cache database engine/paths, project registry, maintenance task IDs/run receipts, and actual scheduling/delivery capabilities |
| Status | What is verified, what remains unavailable or unverified, and limitations relevant to the caller |

Return this nonsecret context to the step that called Access. Keep an existing project
binding if its hub/account still match. Public guidance can be read without authentication;
API work uses the verified connection.
