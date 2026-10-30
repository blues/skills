# Notehub MCP: an authenticated request proxy

Use the hosted server at `https://mcp.notehub.io/mcp`. Use the harness's supported MCP
connection and browser authorization flow; keep credentials inside the connection. Reuse
an existing session. If MCP registration is not available to the agent, explain the
specific connection step the person must take. Do not fabricate harness configuration
commands. `?mode=tools` is an alternative for harnesses that cannot run code mode or need
per-function tool permissions; discover the tools actually exposed in that mode.

After connecting, discover `search`, `execute`, `execute_write` and
`execute_destructive` and read their current schemas. Tool names may carry a harness
prefix. Search finds server functions, including repository functions absent from the
public OpenAPI paths. For the general Notehub API, use the current OpenAPI contract to
map the selected operation's `operationId` to `call_tool`. When the operationId and
resolved parameters are known, skip per-operation `search`/`get_schema` calls; code mode
can combine bounded related calls and return just the needed evidence:

```python
# Inside the server's code-mode execute tool, not a local Python process:
result = await call_tool("get_devices", {"projectOrProductUID": "app:EXAMPLE"})
return result
```

## Exact mapping

Convert the operationId to snake_case in two passes, then lowercase:
`re.sub(r"([A-Z]+)([A-Z][a-z])", r"\1_\2", name)`, followed by
`re.sub(r"([a-z0-9])([A-Z])", r"\1_\2", name)`. Thus `GetDevicesDfuStatus`
becomes `get_devices_dfu_status` and `GetDbNote` becomes `get_db_note`.

- Path parameters, query parameters and JSON object body properties become top-level
  keys with their **specification names**, not snake_case parameter names.
- Free-form JSON or a non-JSON body uses a single `body` argument. Inspect the schema;
  do not wrap a normal object's properties under `body`.
- Headers retain their exact names, e.g. `"If-Match": "\"revision\""`.
- If a query parameter and body property share a name, use `<name>__query` for the query.
- Never send `X-Repository-Key`; the authenticated session supplies repository access.
- Return the result explicitly. A final expression also works; `print()` output is not
  the result. Lists and strings may be wrapped as `{"result": ...}`; unwrap only the
  transport envelope, not a similarly named property in a resource.

## Choose the tool by effect

| Tool | Available effects |
|---|---|
| `execute` | Reads (GET) |
| `execute_write` | Reads and POST/PUT/PATCH writes |
| `execute_destructive` | Those operations plus DELETE and `replace_repository_dataset` |

Exceptions: `delete_device_from_fleets` is a write; `set_global_event_transformation`
is destructive. Confirm current tool metadata. Tool annotations and the harness govern
approval; there is no invented `confirm` argument, and not every harness presents the
same prompt. Stay within the user's existing authorization. A multi-call program is not
a transaction: earlier calls can succeed before a later failure. Inspect before retrying.
Keep outputs small and relevant, but never hide incomplete coverage behind a summary.

The code-mode contract excludes `login`, `o_auth2_client_credentials`,
`get_device_environment_variables_by_pin`, `set_device_environment_variables_by_pin`,
`delete_project`, `set_global_description`, `delete_global_description`,
`upload_firmware` and `download_firmware`. Do not assume every OpenAPI operation is
exposed unchanged. If needed, report a concrete limitation; don't reach for raw tokens.

Specific guards still apply:

- `get_notefile` and `get_db_note` do not take a `delete` argument; use `delete_note`.
- `perform_dfu_action` without a device filter requires explicit `all_devices=true`.
- Use a dry run before `run_job`; actual execution requires `real_run=true`.
- Setting firmware environment variables `_fw`, `_fwc`, `_fwm`, `_fwc_retry` or
  `_fwc_min_version` requires `update_firmware=true`. For removal use the explicit
  delete operation, not an empty string or `-`.
- Use dedicated single-device enable/disable/provision operations rather than
  constructing a `create_job` for one device.

## Grok Bot only

*This section applies only to agents running as Grok Bot. Other harnesses should skip it.*

- Connect the Notehub MCP as a Grok Bot connector (a custom connector for
  `https://mcp.notehub.io/mcp` when it is not in the catalog). The user authorizes it
  through the connect card. On desktop, the browser ends on a `localhost` callback page.
  That page is normal and means sign-in completed.
- The OAuth session lives with the connector, not on the Cloud Computer. Notehub tools
  are callable only from agent turns, not from box scripts. Do not install the CLI,
  store a token or set a token environment variable as a second path.
- Large results can be saved to a file with the dynamic-tool call's `write_to_file`
  option. This is how the cache handoff passes fetched batches to its worker.
