# Bounded reconnaissance for training

Reuse access and the verified project. Never retrieve a
token or build an Authorization header. Read [the live schema](https://api.notefile.net/openapi.yaml)
for current contracts. The agent needs API mechanics; the project's product knowledge
needs the durable meaning learned from them.

## Choose the relevant reference

Read the rows that apply to the current step; reuse context already established.

| If you need to | Read | Then |
|---|---|---|
| Establish a missing connection or execution environment | [Access](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-access/SKILL.md) | Return with the selected authenticated transport |
| Resolve missing or changed project scope | [Project Selection](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-access/projects.md) | Return with the verified data project and skill owner |

Start with schemas, a few recent event bodies, fleet/device descriptions and relevant
nonsecret configuration. Look at full per-device Sessions when the question needs them.
Request only relevant route fields; connection settings may contain secrets. Avoid dumps
of environment values or route credentials into transcripts or drafts. Product names and
field shapes can guide questions but cannot prove what a sensor measures.

- A batch is not a population. Use each operation's documented retrieval controls, preferring
  modified cursors for incremental scans. Verify returned identities, times and counts;
  an ignored or unsupported filter must not silently widen a claim.
- Not observed in a bounded sample is not never emitted. Inventory the union of schemas,
  sample records, firmware and stored knowledge, including conditional and system fields.
- Earliest observed data is not a retention policy. Record filters/window/date and keep
  documented retention, observed coverage and unverified coverage separate.
- Configured is not applied. Fleet membership, environment precedence and route behavior
  need evidence of their actual effect. A rule that exists may never fire.
- Empty has several explanations: permissions, filters, clock choice, ingestion lag,
  retention or absence of measurements. A denied or failed request is not an empty result.
- Inspect real responses and the contract. Capture time, arrival time, modification time
  and an endpoint's filter clock differ. Native `temperature`, `voltage` and location
  need not represent the customer's probe or machine.

Use a fresh cloud cache or available repository for suitable bounded analysis; inspect
its coverage, schema and freshness first. Compare one bounded case to source records
before trusting an unfamiliar representation. Training does not provision repositories,
change datasets or exercise devices. A missing harness capability and a project without
a configured repository are different findings.

Keep connectivity, ingestion lag, cache checkpoints and current capability observations
in private session notes. Shared project knowledge records the durable question, relevant
Notefiles/fields, required coverage and known interpretation limits. Record a temporary
infrastructure failure as a task limitation, not a permanent product rule.
