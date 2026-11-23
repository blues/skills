# openapi-reference

Generates [`notehub-access/openapi-operations.md`](../../notehub-access/openapi-operations.md),
a compact index of every operation in the mirrored
[`openapi.yaml`](../../notehub-access/openapi.yaml). The YAML remains the contract; the
index lets an agent survey every operation without reading the whole specification, and
points it back to the specification for each operation's full contract.

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
Date-shaped values such as `2024-01-01` stay the text they are, as in JSON.

An operation's effect is `read`, `write` or `destructive`, from its HTTP method, unless
the operation says otherwise with an `x-effect` extension (and optionally `x-effect-note`),
or [`overrides.yaml`](overrides.yaml) does. The mirrored specification is never edited, so
`overrides.yaml` holds what the specification doesn't state yet; an entry there wins.
When the specification gains an operation's `x-effect`, delete its entry.

The scheduled `update-openapi` workflow regenerates the reference whenever it runs, in a
`reference` job that runs after the job that mirrors the specification, so a problem here
never holds the mirror back. Pull requests that touch the generator, the overrides, the
specification or the reference run `check-openapi-reference`, which fails if the committed
reference is out of date. To regenerate it yourself, from the repository root:

```sh
go run -C .github/openapi-reference . \
  -overrides overrides.yaml \
  -spec-url 'https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-access/openapi.yaml' \
  ../../notehub-access/openapi.yaml > notehub-access/openapi-operations.md
```

Warnings, such as an override naming an operation the specification no longer has, go
to standard error. `go test -C .github/openapi-reference ./...` runs the tests, which use
a small made-up API.
