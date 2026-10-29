# Interpretation checks for native data

Use the platform references as the definitions for native Notehub data. These checks
resolve specific ambiguities in a question or response; they are not a prerequisite
to every answer.

## Common distinctions that change answers

| If you need to interpret | Read | Apply this distinction |
|---|---|---|
| Session `moved` | [Sessions](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-platform/sessions.md) | Last-motion Unix seconds, distinct from tracking/motion body counts |
| `power_mah` | [Power](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-platform/power.md) | Signed, resettable Mojo charge accumulator in mAh; do not sum snapshots or call it energy |
| DFU `current.product` | [Firmware](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-platform/firmware.md) | Firmware product label, not an `app:` or ProductUID join key |
| DFU attempt `current_version` | [Firmware](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-platform/firmware.md) | The attempt's starting version; top-level `current.version` is the latest report |
| DFU finished or inactive state | [Firmware](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-platform/firmware.md) | Require target-version and outcome evidence before declaring success |
| A journey's start and stop | [Tracking](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-platform/tracking.md) | Use the device's tracking mode and configuration to interpret its journey markers |
| Journey detail results | [Tracking](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-platform/tracking.md) | The endpoint filters coordinate-less and consecutive duplicate-position Events; use raw Events for status counts, stops and sequence completeness |
| Location “best” | [Location](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-platform/location.md) | Selected source and observation time matter; best does not mean fresh or most accurate |
| Zero or absent template fields | [Objects and time](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-platform/objects.md) | Distinguish measured zero from an encoding default or unavailable observation |
| Secondary-SIM usage | [Sessions](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-platform/sessions.md) | Secondary bytes are a breakdown already included in total cellular bytes |

## Historical records and device-specific meaning

Notehub serves records from devices with different firmware, hardware and configuration.
Use those distinctions where they affect a calculation. For example, apply the legacy
SINR conversion only to firmware that reported raw modem values; do not convert readings
already expressed in dB. Use the firmware at the time of the observation when known,
rather than applying today's Device version to all history.

A missing optional field does not establish zero or lack of activity. Check the
relevant feature/configuration and the actual returned fields before drawing that
conclusion. Reuse established context instead of repeating these checks for every query.

## Resolve an actual discrepancy

First identify the exact endpoint, JSON path, producer, unit/scale, observation time
and any transformation between the original record and the queried dataset. Check the
current API contract and the relevant field reference. A schema's type or a field's
English name does not override its documented meaning.

If the discrepancy remains unresolved, qualify the particular conclusion it affects;
continue with supported parts of the answer. For a product-specific exception, keep the field, version/cohort and evidence explicit
in the answer. For an uncovered field, establish its meaning from current documentation
rather than guessing.

| If you need to | Read | Then |
|---|---|---|
| Resolve an endpoint or response contract | [API discovery](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-access/api.md) | Use the live schema for the selected operation |
