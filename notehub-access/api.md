# Discovering and using the Notehub API

The primary contract is the [Notehub OpenAPI specification](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-access/openapi.yaml).
This unmodified mirror of the production schema is checked for updates every 30 minutes;
scheduled checks can be delayed, so the copy can lag behind the service.
Download the complete YAML without summarization, cache it for the task, and record its
retrieval time/hub. Refresh it when changing hubs or when a response contradicts the
contract. Public schema retrieval needs no credential. It does not authorize API writes.

First survey **all** paths, methods, operationIds and tags so the tool surface is not
limited to remembered operations. For each selected operation, resolve `$ref`s in
path-level and operation-level parameters, body, responses and component schemas.
Record required inputs, filters, ordering, limits, timestamps, response shapes, response
headers, retrieval controls and permissions. Then map it to MCP or CLI.
Use current operations for projects, devices, fleets, notes, events, sessions, schemas,
environment variables, routes, firmware and skills as the task needs them. Specialized
MCP functions are discovered from that server, not guessed from OpenAPI components.

Use bounded discovery calls before broad retrieval. Distinguish an error, an empty
successful response and a truncated response. Obey Retry-After where present, otherwise
use bounded exponential backoff with jitter; stop and report sustained failures. Retry
reads safely. For writes with an uncertain outcome, read current state before retrying.
Do not silently substitute a different hub, account, project or time filter.

## Choose the relevant reference

Read the rows that apply to the current step; reuse context already established.

| If you need to | Read | Then |
|---|---|---|
| Invoke the resolved operation through the selected MCP connection | [Notehub MCP: an authenticated request proxy](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-access/mcp.md) | Use 'code mode' to map OperationId and parameters to an 'execute' MCP tool |
| Invoke the resolved operation through the selected CLI environment | [Notehub CLI: an authenticated request proxy](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-access/cli.md) | Use the tool's CURL-like syntax as an API proxy |
