# openapi-to-md (design only)

## Purpose

Agents call an MCP server's code-mode `execute` tool directly as `call_tool('<name>', {params})`. `openapi-to-md` turns any OpenAPI 3 spec into a compact Markdown reference that lets an agent make correct calls — tool names, parameters, and read/write/destructive class — without a search round-trip. The program must be generic: no product-specific logic in the Go. Product policy is carried in vendor extensions or an optional overrides file. The Markdown is meant to be regenerated in CI whenever the spec changes.

## CLI

```
openapi-to-md [flags] <spec>
  <spec>            path, http(s) URL, or - (stdin)
  -overrides PATH   optional YAML (same fields as the extensions, keyed by operationId)
  -o PATH           output file (default stdout)
  -include-hidden   also emit hidden ops, marked hidden (default: omit)
```

Deterministic output (spec order, grouped under `## <first tag>`).

## Naming

`operationId` → snake_case by inserting `_` before every uppercase letter except the first character, then lowercasing; digits and existing `_` stay. Examples: `GetDevices` → `get_devices`, `GetDbNote` → `get_db_note`, `OAuth2ClientCredentials` → `o_auth2_client_credentials`. The same rule splits acronyms: `GetDeviceUID` → `get_device_u_i_d`. This must match the MCP server's own operationId-to-tool-name rule exactly; implementers should confirm against the server and keep the conversion in one function. Operations with no `operationId` are skipped (warned). Duplicate `operationId` values are a fatal error.

## Classification

Default from the HTTP method only: `GET|HEAD|OPTIONS` → `read`; `POST|PUT|PATCH` → `write`; `DELETE` → `destructive`. Anything else → `write`.

Override with `x-mcp-class: read|write|destructive`:

- On the **operation**: replaces the method default.
- On a **parameter** or request-body property: *escalates* the call when that argument is used. Booleans escalate only when the value is `true`; other types escalate when the argument is present and non-null. Order is `read < write < destructive` (a param cannot downgrade a DELETE).

Unknown `x-*` keys (including any product permission metadata) are ignored. `x-mcp-hidden: true` on an operation omits it unless `-include-hidden`. No name denylist in Go.

## Overrides file

Same fields as the extensions, keyed by original `operationId`. File wins over the spec; spec wins over method defaults. Optional `x-mcp-name` on an operation or parameter aliases the call or param key. For specs you cannot edit:

```yaml
OAuth2ClientCredentials:
  x-mcp-hidden: true
GetNotefile:
  parameters:
    delete: { x-mcp-class: destructive }  # name is enough; use "query.delete" only if ambiguous
```

## Parameter rendering

Merge path-item + operation parameters (operation wins). Resolve `$ref` on parameters, request bodies, and schemas before rendering. Keys are the OpenAPI `name` **exactly**, including header casing and hyphens (`X-Auth-Token`, `If-Match`). Do not snake_case param names. A parameter may set `x-mcp-name` if a server requires an identifier-safe key.

Request body, after `$ref` resolution:

- `allOf`: flatten/merge properties and `required` from each member (and the schema itself). If properties remain, flatten them as `body` params; if none remain, treat as unstructured.
- `oneOf` / `anyOf` at the top level: a single param named `body`, with a short note naming the alternatives (do not flatten a union).
- Else if the schema (object or form) has `properties`, flatten each one as a `body` param (`*` if required).
- Otherwise — empty object, free-form `additionalProperties` only, non-object, `octet-stream`, etc. — one param named `body`.

Types: keep formats (`integer(int32)`, `array[string]`, `string[json|csv]`). OpenAPI 3.1 type arrays and 3.0 `nullable: true` are rendered explicitly: `[string, "null"]` or `nullable: true` → `string (nullable)`; `[string, integer]` → `string|integer`.

Render: `[*]name  path|query|header|body  type  [default=…] [enum=a|b] [!write|!destructive]  — short desc`. Descriptions are the first sentence, trimmed to ~120 chars; skip empties. Include defaults and enums, not examples. Responses are out of scope (optional later flag).

## Output format (real lines from the current Notehub spec)

```
# MCP execute reference · <title> <version> · OpenAPI <openapi>
# call_tool('<name>', {param: value, ...})  ·  *required  ·  header keys use the OpenAPI name exactly
# !write / !destructive on a param escalates the call when that param is used

## device
get_devices  read  GET /v1/projects/{projectOrProductUID}/devices
  List devices in a project
  *projectOrProductUID  path  string
  pageSize  query  integer(int32)  default=1000
  pageNum  query  integer(int32)  default=1
  cursor  query  string  — If supplied, results are returned in ascending order of modification time.
  deviceUID  query  array[string]  — A DeviceUID.
  # also query array[string]: tag, serialNumber, fleetUID, notecardFirmware, location, hostFirmware, productUID, sku

get_events  read  GET /v1/projects/{projectOrProductUID}/events
  List events in a project
  *projectOrProductUID  path  string
  pageSize  query  integer(int32)  default=1000
  pageNum  query  integer(int32)  default=1
  cursor  query  string  — If supplied, results are returned in ascending order of modification time.
  sortBy  query  string[best_id|device_serial|device_uid|captured|modified|device_location|tower_location|triangulated_location|best_location]  default=captured
  # also: deviceUID, sensorUID, sortOrder enum=asc|desc, startDate, endDate, dateType,
  # systemFilesOnly, files, format enum=json|csv, serialNumber, fleetUID, sessionUID, eventUID, selectFields

get_notefile  read  GET /v1/projects/{projectOrProductUID}/devices/{deviceUID}/notes/{notefileID}
  Get the notes in a notefile
  *projectOrProductUID  path  string
  *deviceUID  path  string
  *notefileID  path  string  — Notefile name, including its extension, e.g. `config.db` or `commands.qi`.
  max  query  integer(int32)  — The most notes to return.
  deleted  query  boolean  — If true, also return deleted notes.
  delete  query  boolean  !destructive  — If true, delete the notes this request returns. Requires the Owner or Developer role.

delete_device  destructive  DELETE /v1/projects/{projectOrProductUID}/devices/{deviceUID}
  Delete a device
  *projectOrProductUID  path  string
  *deviceUID  path  string

delete_skill  destructive  DELETE /v1/projects/{projectOrProductUID}/skills/{skillName}
  Delete a project skill
  *projectOrProductUID  path  string
  *skillName  path  string  — Relative .md name, case-insensitive, at most 512 UTF-8 bytes.
  *If-Match  header  string  — Required; use one strong ETag or quoted skill revision. Wildcards and weak tags are rejected.

get_device_environment_variables_by_pin  read  GET /v1/products/{productUID}/devices/{deviceUID}/environment_variables_with_pin
  Get device environment variables using a PIN
  *productUID  path  string
  *deviceUID  path  string
  *X-Auth-Token  header  string  — The device's PIN, used instead of an API token.
```

After the annotations below, `get_db_note` is the same pattern as `get_notefile` (`delete` is `!destructive`). Unstructured bodies (e.g. `CreateLegacyWebhookEvent`) become a single `*body  body  object`. `OAuth2ClientCredentials` is omitted once hidden. The `!destructive` marker on `delete` appears only once the spec or overrides file carries `x-mcp-class: destructive`.

## Annotations for Notehub's spec (YAML, not Go)

These cover the known gaps. The current public spec (136 operations) has **no** `Login` operation.

- `OAuth2ClientCredentials` (operation): `x-mcp-hidden: true`
- `GetNotefile` query param `delete`: `x-mcp-class: destructive`
- `GetDbNote` query param `delete`: `x-mcp-class: destructive`
- `CancelJobRun` (operation): `x-mcp-class: destructive` (a POST that cancels a run)
- `Login`, if it is ever re-added: `x-mcp-hidden: true`

Nothing else is required: DELETE already classifies as destructive; `X-Auth-Token` / `If-Match` / `If-None-Match` and unstructured bodies fall out of the generic rules.

## Implementation

~250–350 lines, stdlib `flag`, `net/http` for URLs, [getkin/kin-openapi](https://github.com/getkin/kin-openapi) to load/resolve the spec (3.0 and 3.1), `gopkg.in/yaml.v3` for overrides. No Notehub URLs, operationIds, or parameter names appear in the Go source (testdata fixtures are exempt).

## Acceptance checks

- **Golden output:** commit a pinned copy of the Notehub spec as testdata, fetched from `https://api.notehub.io/openapi.yaml`, plus a golden Markdown generated from that pin as-is. A second golden (or the same test with an overrides fixture) covers the annotations above (`OAuth2ClientCredentials` hidden; `delete` on `GetNotefile` / `GetDbNote` marked `!destructive`).
- **Synthetic specs** (small fixtures) covering: `$ref` resolution for parameters and bodies; `allOf` bodies flatten/merge properties; `oneOf`/`anyOf` bodies become a single `body` param with a note; OpenAPI 3.1 type arrays and `nullable`; duplicate `operationId` values (fatal error); missing `operationId` (skip + warn); path-item vs operation parameter merge (operation wins); header params keep exact names; `x-mcp-class` operation replacement and parameter escalation; `x-mcp-hidden` omits an op, and `-include-hidden` emits it marked hidden; overrides file precedence over spec extensions over method defaults.
- **Determinism:** running the tool twice on the same inputs produces byte-identical output.
- **Hygiene:** `go vet` is clean. A grep over Go source excluding testdata finds no product names (`notehub`, `blues`, or Notehub operationIds). Test code may refer to fixture paths only.
