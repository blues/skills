# Discovering and using the Notehub API

The primary contract is the [Notehub OpenAPI specification](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-access/openapi.yaml).
This YAML file is an unmodified mirror of the production schema which is updated upon every service deployment.
Download the complete YAML without summarization, cache it for the task, and record its
retrieval time. Refresh it when a response contradicts the contract.

## Understand the entities and their relationships

| Entity | Identifier | Relationship |
|---|---|---|
| Account | AccountUID | Contains Projects |
| Project | ProjectUID, also called AppUID | Contains Devices, Fleets, Routes, Events and Sessions; ProjectUID and AppUID are synonymous |
| Device | DeviceUID | Tagged with the FleetUIDs of the many Fleets to which it belongs |
| Fleet | FleetUID | A list of DeviceUIDs belonging to that Fleet |
| Smart Fleet | FleetUID | A Fleet whose Devices are automatically added or removed based on dynamic characteristics; its FleetUID has a `smart:` prefix when listed in a Device's memberships |
| Route | RouteUID | Defines processing and delivery of Events within a Project |
| Event | EventUID | Tagged with the DeviceUID of the one Device that generated it |
| Session | SessionUID | Tagged with the DeviceUID of the one Device that generated it |

Every Project also has a human-readable name and multiple unique ProductUID text strings
representing customer product types or SKUs whose data flows into that Project.
ProductUID is distinct from ProjectUID/AppUID.

Use DeviceUID to relate Events and Sessions to their Device within the selected
ProjectUID. EventUID and SessionUID identify the individual Event and Session records.
Preserve the `smart:` prefix in returned Device memberships, including when caching them.

## Discover operations and their contracts

First survey **all** paths, methods, operationIds and tags so the tool surface is not
limited to remembered operations. For each selected operation, resolve `$ref`s in
path-level and operation-level parameters, body, responses and component schemas.
Record required inputs, filters, ordering, limits, timestamps, response shapes, response
headers, retrieval controls and permissions.
Use current operations for projects, devices, fleets, notes, events, sessions, schemas,
environment variables, routes, firmware and skills as the task needs them.

Use bounded discovery calls before broad retrieval. Distinguish an error, an empty
successful response and a truncated response. Obey Retry-After where present, otherwise
use bounded exponential backoff with jitter; stop and report sustained failures. Retry
reads safely. For writes with an uncertain outcome, read current state before retrying.
Do not silently substitute a different account, project or time filter.
