---
name: notehub-access
description: Establish or reuse secure Notehub access before project selection, Chat or Training. Install or connect one MCP or notehub CLI, authenticate and verify the connection. Return the selected tool, account and invocation method to the caller.
---

# Establish secure Notehub access

Complete missing connection and authentication prerequisites, then return to the step
that called this skill. Reuse verified setup from the same environment. The entry
workflow handles API learning and Cloud Computer setup before project selection.
Follow [working style](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub/working-style.md).

## Connect and authenticate once

Choose one authenticated request proxy: prefer Notehub MCP; use the `notehub` CLI when
MCP is unavailable or the user chooses it.
Both own their credentials—the agent does not retrieve tokens or construct Authorization
headers. If the account and invocation method are already verified, reuse them.

| If you need to | Read | Then |
|---|---|---|
| Register/connect MCP, authenticate or learn its invocation method | [Notehub MCP: an authenticated request proxy](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-access/mcp.md) | Use the harness's supported connection and browser authorization; verify access |
| Install or authenticate the CLI | [Notehub CLI: an authenticated request proxy](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-access/cli.md) | Follow official installation guidance and `notehub --help-agent`; complete its browser/polling login flow |

Verify the selected account and a harmless authenticated read through that tool; retain
useful metadata for later project selection. Respect a cancelled authorization. Change
transport only for a concrete need or the user's choice, explicitly verifying the intended
account. A project permission failure is not a reason to seek another credential source.
Do not repeat login or API discovery at each skill boundary.

## Return the access context

| Context | What to retain |
|---|---|
| Connection | Selected MCP connection or CLI environment, verified account identity and invocation method |
| Verification | Successful authenticated read and useful metadata already obtained |
| Status | What is verified, what remains unavailable or unverified, and limitations relevant to the caller |

Return this nonsecret context to the step that called Access. Keep the selected
ProjectUID if its account still matches. Public guidance can be read without
authentication; API work uses the verified connection.
