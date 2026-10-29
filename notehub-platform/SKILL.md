---
name: notehub-platform
description: Interpret native Blues Notecard and Notehub data when answering questions about firmware updates, device health, connectivity, sessions, signal strength, data usage, journeys, location, motion, power or battery life. Maps user concepts to API objects, Event metadata and system Notefile bodies. Complements the project's own product and schema knowledge.
---

# Notehub platform data

Use this skill to choose **where the evidence lives and what it means**. An OpenAPI
type tells you how to read a value; it does not tell you whether that value measures
the thing the person asked about. Native fields need this interpretation just as
customer-defined bodies need the project's own skills.

For live project work, reuse the verified project and authorization.

Use [working style](../notehub/working-style.md). Familiar identities and already verified
field meanings do not require a glossary detour.

| If you need to | Read | Then |
|---|---|---|
| Answer a live project question or perform a requested operation | [Project chat](../notehub-project-chat/SKILL.md) | Carry the platform interpretation into the project's question workflow |
| Understand unfamiliar object relationships | [Notehub concepts](concepts.md) | Return with the relevant identities and relationships understood |
| Resolve an easily confused field or a mismatch in observed data | [Interpretation checks](verification.md) | Apply the correct meaning and isolate any specific discrepancy |

## Route the question

Load the relevant page below. Read its applicable rules in full; don't load the whole
platform manual or repeat meanings already verified in the current context.

| The person asks about | Start with | Read |
|---|---|---|
| Choosing a native data surface or interpreting unfamiliar Event metadata/time | Endpoint response, producer and clock | [Objects, metadata and time](objects.md) |
| Firmware version, rollout, update stuck or failed | Device's last-reported firmware, target-specific DFU status/history, effective desired configuration | [Firmware](firmware.md) |
| Online, offline, check-ins, sessions, weak signal, connection failures | Device freshness and Session records; `_session.qo` for event history | [Sessions and connectivity](sessions.md), [Health and system files](health.md) |
| Notecard data usage, traffic, plan balance | Session period counters for traffic; SIM usage for plan accounting; Usage API for Notehub consumption | [Sessions and connectivity](sessions.md) |
| Position, GPS versus cell tower or Wi-Fi location | Event location families for history; Device locations for latest known position | [Location](location.md) |
| Journeys, speed, bearing, distance, geofences | Tracking bodies together with Event location metadata; journey API only with its filtering understood | [Tracking and journeys](tracking.md), [Location](location.md) |
| Modem power cycles, battery life, energy, charging | Session high-power counters, Mojo measurements, `_power.qo` transitions; product battery model | [Power](power.md) |
| Reboots, diagnostics, motion, temperature, button events | `_health.qo`, `_health_host.qo`, `_log.qo` and the relevant system file | [Health and system files](health.md) |
| Project, Fleet, Device or Event fields; a field not listed here | Exact endpoint response and schema, then its producer and project meaning | [Objects, metadata and time](objects.md) |

## Use the platform meaning with the product meaning

1. Identify the question's quantity, population and time basis. A customer's “session”
   may mean a machine operation; a Notehub Session means a communication record.
   “Power cycles” might mean the modem, the host, the Notecard or an external supply.
2. Use the project's skills for product roles, sensor wiring, thresholds, custom
   Notefiles and transformations. Use this skill for native semantics. A project may
   redirect tracking to another file, attach custom body fields, or store a transformed
   copy. Check the actual field path and producer before applying either definition.
3. Inspect a small response and, for a repository or cache, the actual dataset schema. Establish whether
   it contains Event envelopes, decoded bodies, full Session records or a derived view.
   A dataset named “sessions” is not proof that all Session API fields are present.
4. Check field presence, firmware generation, configuration, timestamps and coverage
   before deriving an answer. A missing optional field, template-produced zero, stale
   snapshot or filtered history cannot establish that something never happened.
5. Distinguish observed values, derived estimates and unavailable quantities in the
   answer. Give the as-of time where it matters. For an uncovered field, preserve the
   uncertainty and consult its current authoritative definition; do not infer meaning
   from the field's English name alone.

## Use the current service contract

Use these references for native field meanings and the current
[API reference](https://dev.blues.io/api-reference/notehub-api/) and
[OpenAPI specification](https://api.notefile.net/openapi.yaml)
for endpoint contracts. Apply the documented meanings directly. Check device firmware,
hardware or configuration when a specific field depends on them, especially when
interpreting historical records. If an actual
response contradicts the expected contract, investigate that discrepancy and keep its
effect on the answer explicit.
