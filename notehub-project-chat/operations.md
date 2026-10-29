# Requested operations

Use the live schema and current selected transport for project, device, fleet, route,
configuration, note and firmware operations. Don't constrain the agent to Q&A simply
because the original skills were written for it.

## Choose the relevant reference

Read the rows that apply to the current step; reuse context already established.

| If you need to | Read | Then |
|---|---|---|
| Perform the change through MCP | [Notehub MCP: an authenticated request proxy](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-access/mcp.md) | Use the correct execution tier and operation-specific guards |
| Perform the change through the CLI | [Notehub CLI: an authenticated request proxy](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-access/cli.md) | Use the installed agent guide and authenticated API proxy |
| Resolve the current request or response contract | [API discovery](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-access/api.md) | Check operation parameters and verification evidence |
| Configure a requested monitor or alert | [Monitors and alerts](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-project-chat/monitors.md) | Define it and follow activation through verified execution |

1. Resolve the exact target and inspect its state. For a bulk action, enumerate or count
   the selected population and verify the filter. An empty filter must never silently
   mean all devices. Keep product, project and parent skill-owner scope distinct.
2. Check the intended change against the user's request and prior authorization. For an
   ambiguous or consequential action outside that scope, present the proposed targets,
   settings and effect before asking. Do not repeatedly ask for work already authorized.
3. Use dry runs, conditional revisions and dedicated single-device operations where
   available. Follow MCP guards or the installed
   CLI guide. Fetch the current contract; do not invent a `confirm` argument.
4. Execute a bounded change. Record resource IDs, returned status and partial failures,
   without credentials. A batch is not a transaction. After a timeout, inspect the
   actual resource before retrying; a duplicate request can have a second effect.
5. Read back server state. For a device change, distinguish accepted, queued, delivered,
   acknowledged and effective according to the actual operation. Firmware target,
   last-reported firmware and DFU result are separate evidence. Report what is verified
   and what still awaits the device; don't claim success from inactivity.

Project Chat monitoring uses harness-owned scripts and the monitor workflow above,
rather than Notehub-native Monitors/Alerts APIs.
