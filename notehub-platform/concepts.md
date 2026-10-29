# Notehub concepts

Use these relationships to choose the right data surface, then check the live API for
its exact contract and the product's skills for its particular meaning.

## Choose the relevant reference

Read the rows that apply to the current step; reuse context already established.

| If you need to | Read | Then |
|---|---|---|
| Understand entity relationships or find an endpoint | [API discovery](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-access/api.md) | Review the entity model, then resolve the exact operation and response shape |
| Interpret native record fields or clocks | [Objects and time](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-platform/objects.md) | Choose the correct data surface and time basis |

| Concept | Meaning and practical distinction |
|---|---|
| Notecard | The communications module associated with a DeviceUID; it can report its own state and carry host-generated data |
| Device | Identified by DeviceUID and tagged with its FleetUID memberships; represents a Notecard and its last-known metadata/configuration, not proof that the host machine is healthy or its latest measurement is fresh |
| Notehub | The cloud service managing devices, project data, routing and APIs |
| Account | Identified by AccountUID; contains Projects |
| Project | Identified by ProjectUID, synonymous with AppUID, and has a human-readable name; contains Devices, Fleets, Routes, Events and Sessions |
| ProductUID | A unique text string representing a customer product type or SKU whose data flows into a project; each project has multiple ProductUIDs, distinct from its AppUID/ProjectUID and display name |
| Fleet | Identified by FleetUID; a list of DeviceUIDs belonging to that Fleet. Present membership does not establish historical membership |
| Smart Fleet | A Fleet whose Devices are automatically added or removed based on dynamic characteristics; its FleetUID has a `smart:` prefix when listed in a Device's memberships |
| Route | Identified by RouteUID; defines processing and delivery of Events within a Project |
| Note | A JSON document stored/synchronized in a Notefile; database-like Notes and queued Notes have different lifecycles |
| Notefile | A named collection with direction and storage semantics, often `.qo` for outbound queued data, `.qi` for inbound, `.db` for database-like data; verify the actual producer and use |
| Body | The structured JSON portion of a Note/Event; a field's type does not define its units or product meaning |
| Payload | Optional binary content distinct from the JSON body; preserve encoding and API transport rules |
| Event | Identified by EventUID and tagged with the DeviceUID of the one Device that generated it; contains metadata plus body/payload when applicable and is distinct from the current Note or Device state |
| Session | Identified by SessionUID and tagged with the DeviceUID of the one Device that generated it; contains timing, connectivity and usage data. Full Sessions and `_session.qo` Events are different surfaces |
| Notefile schema | A description of fields/types observed or defined for a Notefile; absence from a sample is not proof a field never exists, and schema does not explain semantics |
| Modified cursor | An incremental checkpoint for changes, initialized with epoch microseconds or `0`, then retained as an opaque returned string; not an EventUID or a capture-time filter |
| Product skills | Shared project-owned Markdown knowledge, API `kinds` and revisions; inherited from the parent for subprojects, separate from these public workflow skills |
| Project Repository | Notehub's historical query capability, sometimes called a “Project Dataset”; requires MCP access and a Repository enabled for the Project. Use for whole-dataset questions or when no cache is available. Slower than a local Notehub Cache, faster than direct Notehub API queries |

A cloud cache is another representation, not a new authority on meaning. It retains raw
API records and typed views derived from product skills, with explicit coverage and
freshness. A person's monitor preferences and delivery history are private harness state,
not product skills. Use API discovery for endpoints
and objects and time for record-level interpretation.
