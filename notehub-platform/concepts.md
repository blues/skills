# Notehub concepts

Use these relationships to choose the right data surface, then check the live API for
its exact contract and the product's skills for its particular meaning.

## Choose the relevant reference

Read the rows that apply to the current step; reuse context already established.

| If you need to | Read | Then |
|---|---|---|
| Find an endpoint and its invocation contract | [API discovery](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-access/api.md) | Resolve the exact operation and response shape |
| Interpret native record fields or clocks | [Objects and time](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-platform/objects.md) | Choose the correct data surface and time basis |

| Concept | Meaning and practical distinction |
|---|---|
| Notecard | The communications module associated with a stable device UID; it can report its own state and carry host-generated data |
| Device | Notehub's representation of a Notecard and its last-known metadata/configuration; not proof that the host machine is healthy or that its latest measurement is fresh |
| Notehub | The cloud service managing devices, project data, routing and APIs |
| Product | A product identity used across deployments; a product UID is not necessarily a single concrete project |
| Project | A data/configuration and access scope; use its verified UID in queries and joins |
| Fleet | A grouping of devices for management and configuration; present membership does not establish historical membership |
| Smart Fleet | Membership derived from configured rules; inspect the effective rule and evaluation semantics, not just the fleet's name |
| Note | A JSON document stored/synchronized in a Notefile; database-like Notes and queued Notes have different lifecycles |
| Notefile | A named collection with direction and storage semantics, often `.qo` for outbound queued data, `.qi` for inbound, `.db` for database-like data; verify the actual producer and use |
| Body | The structured JSON portion of a Note/Event; a field's type does not define its units or product meaning |
| Payload | Optional binary content distinct from the JSON body; preserve encoding and API transport rules |
| Event | An occurrence recorded by Notehub, with metadata plus body/payload when applicable; it is not the same as the current Note or Device state |
| Session | A communication record with timing, connectivity and usage data; full Sessions and `_session.qo` Events are different surfaces |
| Notefile schema | A description of fields/types observed or defined for a Notefile; absence from a sample is not proof a field never exists, and schema does not explain semantics |
| Modified cursor | An incremental checkpoint for changes, initialized with epoch microseconds or `0`, then retained as an opaque returned string; not an Event UID or a capture-time filter |
| Product skills | Shared project-owned Markdown knowledge, API `kinds` and revisions; inherited from the parent for subprojects, separate from these public workflow skills |
| Repository | An analytical/historical data surface, when configured and exposed through the chosen connection; inspect schema, coverage and freshness before querying |

A cloud cache is another representation, not a new authority on meaning. It retains raw
API records and typed views derived from product skills, with explicit coverage and
freshness. A person's monitor preferences and delivery history are private harness state,
not product skills. Use API discovery for endpoints
and objects and time for record-level interpretation.
