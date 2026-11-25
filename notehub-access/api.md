# Discovering and using the Notehub API

The primary contract is the [Notehub OpenAPI specification](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-access/openapi.yaml).
This YAML file is an unmodified mirror of the production schema, checked every 30 minutes and updated when it changes.
Download the complete YAML without summarization, cache it for the task, and record its
retrieval time. Refresh it when a response contradicts the contract.

## Understand the entities and their relationships

| Entity | Identifier | Relationship |
|---|---|---|
| Organization | organizationUID (`billing_account_uid` in project bodies) | Owns and is billed for Projects; people are Project members with a role |
| Project | ProjectUID, also called AppUID | Contains Devices, Fleets, Routes, Events and Sessions; ProjectUID and AppUID are synonymous |
| Device | DeviceUID | Tagged with the FleetUIDs of the many Fleets to which it belongs |
| Fleet | FleetUID | A group of Devices; membership lives on each Device's `fleet_uids`, and `GetFleetDevices` lists a Fleet's Devices |
| Smart Fleet | FleetUID | A Fleet whose Devices are automatically added or removed based on dynamic characteristics; API responses list its members by the bare FleetUID, once |
| Route | RouteUID | Defines processing and delivery of Events within a Project |
| Event | EventUID | Tagged with the DeviceUID of the one Device that generated it |
| Session | SessionUID | Tagged with the DeviceUID of the one Device that generated it |

Every Project also has a human-readable name and zero or more unique ProductUID text
strings (a Project created through the API has none until `CreateProduct`) representing
customer product types or SKUs whose data flows into that Project.
ProductUID is distinct from ProjectUID/AppUID.

Use DeviceUID to relate Events and Sessions to their Device within the selected
ProjectUID. EventUID and SessionUID identify the individual Event and Session records.
Routed Events can still carry the stored `smart:fleet:…` spelling; strip the `smart:`
prefix before joining those to API FleetUIDs.

## Discover operations and their contracts

First survey **all** paths, methods, operationIds and tags so the tool surface is not
limited to remembered operations. For each selected operation, resolve `$ref`s in
path-level and operation-level parameters, body, responses and component schemas.
Record required inputs, filters, ordering, limits, timestamps, response shapes, response
headers, retrieval controls and permissions.
Use current operations for projects, devices, fleets, notes, events, sessions, schemas,
environment variables, routes, firmware and skills as the task needs them.

Use bounded discovery calls before broad retrieval. The Notehub API allows 7 requests per
minute per billing account (Organization), shared by all of its clients; a request over
the limit returns 429 with `Retry-After` and `X-RateLimit-*` headers. Plan large reads
within that budget. Distinguish an error, an empty successful response and a truncated
response. Obey Retry-After where present (the `notehub` CLI already waits it out and
retries), otherwise use bounded exponential backoff with jitter; stop and report
sustained failures. Retry reads safely. For writes with an uncertain outcome, read
current state before retrying.
Do not silently substitute a different account, project or time filter.
