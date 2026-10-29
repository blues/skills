# Notehub CLI: an authenticated request proxy

Use this when MCP is unavailable, or when the user has chosen the CLI. Run
`notehub --help-agent` and follow the installed version's guide; this page is a workflow,
not a replacement for that guide. If the executable is absent, use the official
[Note CLI installation instructions](https://github.com/blues/note-cli) for the host.
Do not guess download URLs or run a privileged install without a need.

## Authenticate without handling credentials

Run `notehub auth status --json`. Inspect `signed_in`, the reported `user` and any
`reason`/`action`; exit zero alone does not mean signed in. A network failure does not
mean the user must log in again. Reuse a working login in the same OS user/environment.

When login is needed, set `--name` to the name of the agent harness you are running in,
such as `Grok Bot`, `Muse`, `Claude` or `Codex`, so the user can recognize which agent is
requesting access. Start `notehub login --name "YOUR_HARNESS_NAME" --json`, replacing
the placeholder with that actual name, in a persistent process and read its JSONL
stdout as it arrives. Immediately give the user the exact
approval URL the process returns. The user opens it in their browser; the agent keeps
that **same process** alive while it polls. Progress is not success: await the terminal
`success: true` result and successful process exit, then verify status. Respect the
reported expiry (normally about ten minutes), cancellation and failure. Do not start
parallel login attempts or retry a cancelled login. An expired login needs a new attempt
only if the user still wants to connect.

The CLI stores and refreshes credentials itself. Never run `auth token`, inspect its
credential files, supply `--token`, print an inherited token, or add Authorization headers.
If status reports an environment credential override, report the source without its
value and follow the guide; do not silently change the user's environment. Log out only
when requested. Unattended jobs must run in an environment where this same authorized
login works and refreshes; do not copy secrets into job configuration.

## Call the API

Use the selected operation's method, path and parameters from the current schema.
Reuse the schema already obtained during API discovery; `notehub api --spec` retrieves
it when needed.

```sh
notehub api GET /v1/projects
notehub api GET '/v1/projects/{project}/devices' --project 'app:EXAMPLE'
notehub api GET '/v1/projects/{project}/events' --project 'app:EXAMPLE' \
  --raw-field cursor='EXACT_CURSOR_FROM_PREVIOUS_RESPONSE'

# Skills API examples: metadata, one document, or all content.
notehub api GET '/v1/projects/{project}/skills' --project 'app:EXAMPLE'
# Substitute an actual catalog links.self for this example document path.
notehub api GET '/v1/projects/app:EXAMPLE/skills/references%2Fdevices.md'
notehub api GET '/v1/projects/{project}/skills' --project 'app:EXAMPLE' \
  --field includeContent=true
```

The method is explicit. The path is an API path, not a full URL. `--project`, `--product`
and `--device` fill supported path placeholders; they do not add filters. Use the exact
operation's parameter names and permitted values.

| Need | CLI behavior |
|---|---|
| Typed scalar or file contents | `--field key=value` interprets numbers, booleans and null; `@file` or `@-` reads text from a file or stdin |
| Exact string, especially an opaque cursor | `--raw-field key=value` |
| Query parameters on GET/HEAD | Fields become query parameters |
| Fields on methods other than GET/HEAD | Fields become JSON body properties; put query parameters in the encoded path |
| Nested/complete JSON | `--input request.json` or `--input -`; cannot mix with field flags |
| Repeated query keys | Encode them in the path; repeated fields are not a substitute |
| Conditional read | `--header 'If-None-Match: "ETAG_FROM_THIS_REPRESENTATION"'` reuses an unchanged cached response; never add an Authorization header |
| Response headers | `--include` precedes the body with status/headers; parse separately |

Stdout is the response body, stderr diagnostics; HTTP failures exit nonzero. A successful
response can have an empty body. Avoid `--verbose` unless needed; logs can contain
customer data. `--pretty` buffers output, so avoid it for streams. For binary firmware,
use the documented `notehub firmware get` command: the API proxy may append a newline.

`--paginate` follows modified cursors (starting at `0` unless supplied), emits successive
responses rather than one JSON document, and fails on incomplete cursor metadata or
nonadvancing cursors with more data. Omit it when the calling worker needs control over
each batch and its checkpoint; use individual API requests instead.

The CLI's API proxy covers the OpenAPI surface. The special Project Repository query
capability requires MCP access and a Repository enabled for the Project; it is unavailable
through the CLI. Use the Notehub Cache or ordinary Notehub APIs for supported questions.
Do not establish a second authentication path to obtain Repository access.

On 401 check status/refresh; on 403 check permissions; on 400/404 check the current
contract and resource. Back off for rate limits. After an uncertain write, inspect the
resource before retrying: a lost response does not mean the write failed.
