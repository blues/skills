# openapi-reference

Generates two views of the mirrored [`openapi.yaml`](../../notehub-access/openapi.yaml),
which remains the contract:

- [`notehub-access/openapi-operations.md`](../../notehub-access/openapi-operations.md), a
  compact index of every operation, so an agent can survey the API without reading the
  whole specification.
- [`notehub-access/operations/`](../../notehub-access/operations), one contract file per
  operation, linked from the index: the operation and every component it references,
  directly or indirectly, cut unchanged from the specification (order, descriptions and
  comments included). An agent reads the contracts of the operations it uses, including
  what their response fields mean, instead of the whole specification.

The program reads any single-file OpenAPI 3 document (3.0, 3.1 or 3.2) and has no
knowledge of Notehub. For each operation in `paths` it prints the operationId, the call
name (the operationId in snake_case), the effect, the method and path, any tags beyond the
first, the summary, and the parameters by location. A request body is marked optional
unless the specification requires it, and a body with several media types or `oneOf` /
`anyOf` alternatives says so. Deprecated parameters, `readOnly` body properties, and the
`Accept`, `Content-Type` and `Authorization` headers OpenAPI says to ignore are left out.

It doesn't list webhooks, callbacks or OpenAPI 3.2 `additionalOperations`; the reference's
header says how many of each it left out. It follows local JSON Pointer references
(including path items that define operations beside their `$ref`), but not references to
other files, `$anchor` references or `$dynamicRef`. Those, references that point to nothing,
and cycles of references are reported on standard error rather than silently skipped.
Date-shaped values such as `2024-01-01` stay the text they are, as in JSON. A contract
file can't include a reference target reached through a list index, and says so.

An operation's effect is `read`, `write` or `destructive`, from its HTTP method, unless
the operation says otherwise with an `x-effect` extension (and optionally `x-effect-note`),
or [`overrides.yaml`](overrides.yaml) does. The mirrored specification is never edited, so
`overrides.yaml` holds what the specification doesn't state yet; an entry there wins.
When the specification gains an operation's `x-effect`, delete its entry.

The scheduled `update-openapi` workflow regenerates both whenever it runs, in a
`reference` job that runs after the job that mirrors the specification, so a problem here
never holds the mirror back. Pull requests that touch the generator, the overrides, the
specification or the generated files run `check-openapi-reference`, which fails if the
committed files are out of date. To regenerate them yourself, from the repository root:

```sh
go run -C .github/openapi-reference . \
  -overrides overrides.yaml \
  -spec-url 'https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-access/openapi.yaml' \
  -operations-dir ../../notehub-access/operations \
  -operations-url 'https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-access/operations' \
  ../../notehub-access/openapi.yaml > notehub-access/openapi-operations.md
```

The contract files are written only once the index renders. Generated files that no longer
belong to an operation are removed; any other file in the directory is left alone.

Warnings, such as an override naming an operation the specification no longer has, go
to standard error. `go test -C .github/openapi-reference ./...` runs the tests, which use
a small made-up API.
