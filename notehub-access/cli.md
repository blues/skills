# Notehub CLI: an authenticated request proxy

Use this guide when Access determined MCP cannot be connected for this harness (no
Notehub MCP tools and no way to add an MCP server, like Muse), or when the CLI is chosen or preferred
among connectable options. Run
`notehub --help-agent` and follow the installed version's guide; this page is a workflow,
not a replacement for that guide. If the executable is absent, install it from the
official [note-cli releases](https://github.com/blues/note-cli/releases) for the host OS
and CPU. Do not guess download URLs, do not assume Homebrew exists (especially on Linux),
and do not run a privileged install without a need.

## Install when `notehub` is missing

**Linux and other non-macOS hosts (no Homebrew, no sudo).** Read the latest release's
asset list rather than guessing names. Archives are named
`note-cli_<version>_<os>_<arch>.tar.gz` (for example `linux_amd64`, `linux_arm64`,
`linux_armv6`, `windows_amd64`) and contain the `notehub` and `notecard` binaries; install both:

```sh
# Map `uname -m` to the asset arch: x86_64 -> amd64, aarch64/arm64 -> arm64, armv6l/armv7l -> armv6
ARCH=amd64
URL=$(curl -fsSL https://api.github.com/repos/blues/note-cli/releases/latest \
  | grep -o "https://[^\"]*note-cli_[^\"]*_linux_${ARCH}\.tar\.gz")
mkdir -p ~/.local/bin && curl -fsSL "$URL" | tar -xz -C ~/.local/bin notehub notecard
export PATH="$HOME/.local/bin:$PATH"   # persist in the shell profile if not already on PATH
notehub --version && notehub --help-agent
```

Optionally verify the archive against the release's `checksums.txt`. Debian/Ubuntu
`.deb` packages (`notehub_<version>_linux_<arch>.deb`) are also published, but installing
them needs root; prefer the user-local archive.

**macOS.** `brew install --cask blues/note-cli/note-cli` if Homebrew is present, or use
the `darwin_amd64`/`darwin_arm64` archive with the same steps as above.

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

If your turn ends after you give the user the approval URL (nothing wakes you when a
detached login process finishes), do not promise to keep watching. Ask the user to reply
"done" once they have approved, or set a short check if your harness supports scheduled
or background wakeups. When you resume, run `notehub auth status --json` and confirm
`signed_in` and the expected `user` before continuing.

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
