# openapi-to-md

Turns an OpenAPI 3.0 or 3.1 spec into a compact Markdown reference. With it, an agent can call an MCP server's code-mode `execute` tool as `call_tool('<name>', {params})` correctly, without a search round-trip. The reference gives tool names, exact param keys and types, nested request shapes, and the read/write/destructive class. See [DESIGN.md](DESIGN.md) for the design. This README covers usage and where the implementation goes beyond the design.

The Go source has no product knowledge. Product policy lives in `x-mcp-*` extensions in the spec, or in an overrides file.

## In this repo

The scheduled workflow `.github/workflows/update-openapi.yml` mirrors Notehub's spec to `notehub-access/openapi.yaml`. A follow-on `reference` job then regenerates `notehub-access/openapi.md` from it, using Notehub's annotations in `notehub-access/openapi-overrides.yaml`. That job runs this module's tests first, and commits only when the output changes. To regenerate by hand from the repository root:

```
go run -C openapi-to-md . -overrides ../notehub-access/openapi-overrides.yaml ../notehub-access/openapi.yaml > notehub-access/openapi.md
```

## Usage

```
go build            # or: go install
openapi-to-md [flags] <spec>
  <spec>             path, http(s) URL, or - (stdin)
  -overrides PATH    optional YAML of x-mcp-* fields keyed by operationId
  -o PATH            output file (default, or -, is stdout)
  -include-hidden    also emit hidden ops, marked hidden
  -common N          define a param once, under "common params", when N or more
                     operations share it (default 2; 0 keeps every param with its call)
```

Flags can go before or after the spec. Everything after `--` is positional. Warnings go to stderr: skipped operations, overrides that match nothing, tool names MCP can't use, param keys that collide. Errors exit 1: duplicate operationIds, tool-name collisions, invalid `x-mcp-*` values (including on nested fields), and unknown keys or extra YAML documents in the overrides file. The output is deterministic, so CI can regenerate it and diff the result.

## Output

From the pinned Notehub spec (`testdata/cases/notehub-annotated/want.md`):

```
# MCP execute reference · Notehub API 1.2.0 · OpenAPI 3.0.3
# call_tool('<name>', {param: value, ...})  ·  *required  ·  header keys use the OpenAPI name exactly
# !write / !destructive on a param escalates the call when that param is used (a boolean: when true)
# each {x} in a path is a required param x, a string unless an `x  path` line defines it  ·  `a, b  query  string` = params with the same spec
# `+ a, *b` = params defined under ## common params; a param line under a call applies to that call only
# a capitalized type such as Foo is defined under ## types; an inline object's fields are indented below it

## common params
pageSize  query  integer(int32)  default=1000
cursor  query  string  — If supplied, results are returned in ascending order of modification time. A cursor is a timestamp: Unix epoch microseconds (UTC), optionally followed by `.N`.
...

## device
get_db_note  read  GET /v1/projects/{projectOrProductUID}/devices/{deviceUID}/notes/{notefileID}/{noteID}
  delete  query  boolean  !destructive  — If true, delete the note after returning it.
  deleted  query  boolean  — If true, also return the note if it has been deleted.

## project
update_fleet  write  PUT /v1/projects/{projectOrProductUID}/fleets/{fleetUID}
  Always send `smart_rule_enabled`; leaving it out turns the smart rule off.
  + smart_rule
  label  body  string  — New fleet name.
  smart_rule_enabled  body  boolean  — Whether `smart_rule` adds and removes devices. Always send it; omitting it sets `false`.
  ...

## types
HttpRoute
  filter  Filter  — Filter applied to route data.
  http_headers  map[string]string
  throttle_ms, timeout  integer(int32)
  url  string(uri)
```

Each operation has a line `name  class  METHOD path`, then a summary, then one line per param: `[*]key  in  type  [default=…]  [!class]  [deprecated]  — description`. Types keep formats and enums: `integer(int32)`, `array[string]`, `string[json|csv]`, `map[string]string`, `map[string](string|integer)`, `string|integer`, `string (nullable)` (`|null` when nested). A named type is a `$ref`'d object, its name capitalized, defined once under `## types`. A recursive non-object type is defined as an alias: `JsonValue = string|number|boolean|array[JsonValue]|map[string]JsonValue`.

## Classes, hidden ops, names

| Extension | On | Effect |
|---|---|---|
| `x-mcp-class: read\|write\|destructive` | operation | Replaces the method default: `GET`/`HEAD`/`OPTIONS` are read, `DELETE` is destructive, anything else is write. |
| `x-mcp-class` | parameter, its schema, body property, nested field | Escalates the call when used (a boolean only when `true`). Shown as `!write` / `!destructive` when it is above the call's class. |
| `x-mcp-hidden: true` | operation only | Omitted unless `-include-hidden` is set. |
| `x-mcp-name: alias` | operation, parameter, property | The tool name or param key the server uses instead. A path placeholder shows the alias. |

Other `x-*` keys are ignored. An overrides file carries the same fields, keyed by the original operationId. Parameters are keyed by their OpenAPI name, or by `in.name` (`query.delete`, `body.force`) when the name alone is ambiguous. The file wins over the spec, and the spec wins over method defaults:

```yaml
OAuth2ClientCredentials:
  x-mcp-hidden: true
GetNotefile:
  parameters:
    delete: { x-mcp-class: destructive }
CancelJobRun:
  x-mcp-class: destructive
```

A tool name is its operationId in snake_case. A `_` goes before an uppercase letter that follows a lowercase letter or a digit, or that starts a word after an acronym. So `ListPets` becomes `list_pets`, `GetHTTPStatus` becomes `get_http_status`, and `OAuth2Token` becomes `o_auth2_token`. The rule lives in `toolName` (names.go).

## How it stays small, and what it adds back

The rules are generic and none of them checks for product names.

- **Descriptions get one sentence, plus one that matters.** A param gets the first sentence that says something its name doesn't. `tag — Tag filter.` is an echo, so it is skipped and "Matches the whole tag, case-insensitive." is used instead. It also gets the first later sentence that states a requirement, a prohibition, what omitting it does, or a unit that the first sentence lacks: "Always send it; omitting it sets `false`.", "Required.", "Unix epoch microseconds". An operation gets its summary, unless the summary only restates the tool name. It also gets the first description sentence that says what leaving fields out does, or how far the call reaches: "Omitted fields are cleared or reset…", "With no filters, that is every device in the project." Sentences about roles and status codes are skipped.
- **Path params are implied.** Each `{x}` in the path is a required string param `x`. It gets its own line only when it has a description, a non-string type or an alias, or when a different `x  path` definition under common params would otherwise apply to it.
- **Identical lines merge.** Params whose lines differ only by key share one line: `serialNumber, location, sku  query  array[string]  — …`.
- **Common params.** A param that several calls use with the same spec is defined once under `## common params` and referenced as `+ a, *b`. A definition is kept only if it saves tokens, so small specs don't get the section. A param that escalates, even through a nested field, stays with its call. `-common 0` turns the section off.
- **Nested request objects are defined once.** A `$ref`'d object renders as its name. An inline object's fields are indented under its param. Inline union alternatives become types named by their `title`, or by the param name and a number. A union of objects in a body gets a note: `one of Card | BankAccount`. A union beside an object's properties also becomes a note: `body requires one of: email | phone`, or `plus the fields of one of: Dog | Cat`. A discriminator's values are listed: `method: card=CardPayment | bank_transfer=BankPayment`. Only types reachable from emitted calls are listed. A hidden op reserves no type names.
- **Compositions are resolved.** `{allOf: [{$ref: X}], description: …}` renders as `X`, keeping its enum, format and named type. Scalar `allOf` members merge, with enums intersected. A redeclared property in a later member narrows the earlier one. In OpenAPI 3.1, `$ref` siblings, `contentEncoding` (`string(base64)`), `contentMediaType` and Reference Object descriptions are honored.
- **Noise is dropped.** Dropped items are `default=false`, `default=""`, `readOnly` properties in requests, bodies that admit only `{}`, and the `Accept`, `Content-Type` and `Authorization` header params (which OpenAPI ignores).
- **Values that are written only in prose are surfaced.** A schema may have no `enum` while its description opens with a list of `` `code` ``-led values. If it is a string or a number, those values are appended (`Values: http|proxy|mqtt|…`). A list introduced by "for example" is labeled `e.g.` instead. Lists of errors, fields or keys are ignored. OperationIds mentioned in descriptions are rewritten to tool names (``List with `get_organizations` ``).

On the pinned Notehub spec (136 operations), the spec itself is about 73.8k tokens and a literal rendering of the design's line format is about 17.5k. This output is about 13.5k (o200k token counts), and that includes every nested request type, which the literal format leaves out.

## Loading and safety

The spec is read from a file, a URL or stdin, and its root document is parsed as written. That parse gives spec order (even when path items are `$ref`'d to other files), original 3.1 descriptions, and an alias check. A document whose YAML aliases would expand past a million nodes is refused. Without that check, a 1 KB spec could take gigabytes, because kin-openapi decodes each schema separately and that bypasses the YAML library's own guard.

`$ref`s may name http(s) documents. They may name local files only by a relative path from a local document. Absolute paths and `file://` URLs are refused, and a remote document never reads local files. Relative refs resolve against the spec's final URL after redirects, per RFC 3986, so they don't inherit its query string. When the spec URL is on a public host, neither its `$ref`s nor its redirects may reach loopback, private or link-local addresses. That address is checked as the connection is dialed. Each fetch is capped at 64 MiB.

## Departures from DESIGN.md

- **Naming rule.** The design's rule puts a `_` before every uppercase letter, which turns `GetAWSRoleConfig` into `get_a_w_s_role_config`. The live Notehub MCP server calls it `get_aws_role_config`. The acronym-aware rule here matches the server for all 122 spec operations it exposes; the design's rule gets this one wrong. The design's other examples (`get_devices`, `get_db_note`, `o_auth2_client_credentials`) come out the same.
- **`## types`, inline objects and union notes.** The design renders nested bodies as `object`. In the eval below, every variant without nested types got nested-body tasks wrong, for example `headers` for `http_headers` and `secret_access_key` for `access_key_secret`.
- **More than one sentence per description, and operation notes.** These hold requirements, omission effects and units. Without them, every agent in the first eval renamed a fleet without resending `smart_rule_enabled`. The spec says that silently turns the smart rule off.
- **Required `*` comes from the schema** even when `requestBody.required` is false, because flattened params can't express "send no body at all".
- **`## common params` and `+` lines, echo skipping, implied path params and merged lines** (all above). The design's `# also` lines were a lossy form of merging. Merged lines here keep every type.
- **Enums** render in the type (`string[a|b]`), as in the design's examples, not as `enum=a|b`.
- **Group order** follows first use of each first tag, which is the design's "spec order". A tag named `types` or `common params` gets the heading `## tag: …`.

### Eval

Agents got only one generated reference each and wrote calls for realistic tasks: filters, headers, escalation, nested routes, monitors, jobs, update semantics and units. A script checked each call against the real spec and the expected values. It checked the tool, class, required keys, unknown keys, types, enums and nested fields. In the first round, 6 doc variants were tested on Opus and Haiku with 26 tasks. Variants with the types section scored 26/26 on both models, including the compact common-params variant. The design-literal variant and the flat variant without types scored 22–25/26. A second round added 3 tasks that depend on a description's later sentences: turning a watchdog off (a negative value; `0` keeps it), a usage threshold measured in bytes, and a cursor in microseconds. The final output scored 29/29 on both models. The earlier first-sentence-only output scored 26/29 on both and failed exactly those 3 tasks.

## Tests

```
go test ./...                    # goldens, unit tests, determinism, URL/stdin/-o, safety, hygiene
go test -run TestCases -update   # accept new output after reviewing the diff
go vet ./...
```

Each directory under `testdata/cases/` holds `args` (CLI args, one per line), its inputs, and the expected `want.md` (stdout), `want.err` (stderr) and `want.code` (exit status). `notehub/` holds a pinned copy of `https://api.notehub.io/openapi.yaml`, rendered as-is. `notehub-annotated/` renders the same spec with the design's annotations, which are in `overrides.yaml`. Forty-nine small synthetic cases cover the following:

- refs, `allOf` and its wrappers, unions and discriminators, recursion
- OpenAPI 3.1 types and `nullable`
- duplicate and missing operationIds, param merging, headers
- classes, hidden ops, override precedence and errors, naming and type names
- descriptions, constraint sentences, value lists, media types
- external and split refs, JSON input, common params and shadowed path params
- warnings and CLI edge cases

`TestHygiene` fails if any Go file outside `testdata` mentions a word or spec named in `testdata/hygiene.txt`. For a named spec, that means any of its operationIds, its tool names, or its distinctive parameter names.
