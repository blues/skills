# Health, diagnostics and other system Notefiles

“Device health” spans the Notecard, host application, sensors, connectivity and power.
Identify the component and symptom before selecting a file. A successful sync does
not prove the host or sensor is healthy; an empty diagnostic file does not prove it
has no faults. Logging, feature configuration, transport suppression, retention and
delayed uploads affect what can be seen.

## Choose the relevant reference

Read the rows that apply to the current step; reuse context already established.

| If you need to | Read | Then |
|---|---|---|
| Interpret Mojo charge in diagnostic messages | [Power](power.md) | Account for resets and charging before deriving usage |
| Determine whether firmware actually updated | [DFU state](firmware.md) | Combine log evidence with target-specific status and reported firmware |
| Resolve an ambiguous native field or a mismatch in observed data | [Interpretation checks](verification.md) | Apply the relevant field meaning and isolate the discrepancy |

## System-file map

| File | Producer / use | Where its meaning is defined |
|---|---|---|
| `_health.qo` | Notecard operational/health messages, including boot and DFU diagnostics | Log schema below |
| `_health_host.qo` | Host-submitted health messages through the Notecard | Same basic log shape; host's method/text vocabulary requires project knowledge |
| `_log.qo` | Configured Notecard diagnostics, often more verbose and subsystem-specific | Log schema below; emission is conditional |
| `_session.qo` | Notehub-generated communication start/end Events | [Sessions](sessions.md); much of the data is Event metadata, not body |
| `_track.qo` or configured tracking file | Notecard tracking and supported Notehub fence transitions | [Tracking](tracking.md), [location](location.md) |
| `_power.qo` | Supported Notecard USB/charging source transitions | [Power](power.md) |
| `_motion.qo` or configured motion file | Notecard accelerometer reporting | Motion schema below |
| `_temp.qo` | Configured automatic temperature/diagnostic reporting | Temperature schema below |
| `_button.qo` or configured button file | Notecard button / AUX GPIO reporting | Button schema below |
| `_geolocate.qo` | Notehub radio-geolocation result/diagnostic | Body can include `location`, `country`, `radios` or `err`; coordinate result is in Event `tri_*` context, with optional `triangulate` inputs |
| `_tower.qo` | Notehub serving-tower update | Body `location`, `country`, `tower_lat`, `tower_lon`, `tower_time`, `tower` (cell tuple); optional `weather`/`weather_err` follow their own response schema |
| `_socket.qo` | Notehub socket usage report | Body `bytes_sent`, `bytes_rcvd` for that socket operation; not whole-device lifetime/SIM usage |
| `_env.dbs`, `_env.qis` | Environment/configuration and device-reported management state | Per-note management schemas; use effective environment and public Device/DFU APIs rather than treating every management body as a sensor reading |
| `_fence.dbs`, `_fence.dbx` | Synchronized device fences / service-local fences on supporting versions | Named fence notes with `body.values` alternating latitude/longitude polygon vertices; implicit closing edge. Definitions are configuration, not crossing Events |
| `_req.qis`, `_rsp.qos` | Remote Notecard request/response transport | Operation-specific request/response schemas, not a stable telemetry body |
| `_web.qo`, `_signal.qo`, `_watchdog.qo`, `_synclog.qi`, `_air.qo` and other encountered system files | Web/signal/watchdog/sync diagnostics or specialized sensor features | Inspect the producer's current API documentation, generation and sample. This skill does not claim a complete body dictionary for these specialized families |

An underscore is a naming clue, not enough to identify a producer or schema. Use the
actual file, producer and configuration to interpret unfamiliar system records.

## Health and log bodies

| Field | Meaning |
|---|---|
| `text` string | Diagnostic message. Interpret the whole message with method, firmware and time; wording is not a permanent machine enum. |
| `method` string | Subsystem/event category, such as boot or DFU where supplied. Host-defined values need the project's vocabulary. |
| `alert` boolean | Message flagged for attention by its producer. Not a universal severity scale, safety alarm or diagnosis. |
| `voltage` number | Supply volts sampled for health/power diagnostics where available. |
| `voltage_mode` string | Configured voltage operating classification. Not charge percentage. |
| `milliamp_hours` number | Mojo accumulator when present. Use power, including resets and charging. |

Fields vary by method/build/template. The log template can materialize defaults even
when a message did not supply a measurement. There is no guarantee that all log
messages contain valid voltage or charge readings. Timestamp and location usually
live in the Event envelope.

For reboot history, inspect boot-related health messages with Event capture time and
firmware context. Distinguish a Notecard boot from modem cycling or a host restart.
For firmware failures, correlate DFU messages with the actual target/attempt and
DFU state; do not infer success from a single “download complete” line.

The Device health-log endpoint reads a **bounded recent cache** (up to 100 entries), with supported filters such as boot and DFU-completed
message matching. Embedded Device `health_log` can apply different filtering and omit
method-tagged entries. Date filters do not turn either cache into full history. Use
Event APIs or repositories over `_health.qo` and `_health_host.qo` for retained historical analysis,
and state the coverage. Device latest-event caches are also snapshots updated at sync
boundaries/periodically, not guaranteed instantaneous live streams.

## Motion body

| Field | Meaning |
|---|---|
| `motion` integer | Accelerometer interrupt/event count since the preceding motion report. Not displacement or a timestamp. |
| `tilt` integer | Count of detected orientation changes since the preceding report. Not an angle in degrees. |
| `orientation` string | Orientation category relative to the device's axes. Not compass heading. Mounting determines its product meaning. |
| `movements` string | Encoded time buckets of motion events, newest first. In this encoding `0`–`9` represent 0–9, `A`–`P` represent 10–25, and `Z` represents 26 or more. Bucket duration comes from motion configuration; do not read it as a decimal number or infer a duration from string length alone. |
| `voltage`, `temperature` number | Supply volts and onboard temperature in degrees Celsius. Missing/default and sensor-failure rules still apply. |
| `alert` boolean | Whether this report was flagged as an alert by the motion producer/configuration. |

Cadence, sensitivity, orientation triggers and motion thresholds are configurable.
Motion events can be vibration while the product stays in place, and a moving vehicle
can have sparse accelerometer reports. Project mounting and use decide the inference.
Session/Event `moved` is instead a **last-motion timestamp**.

## Temperature body

The default `temperature` is normally the Notecard's onboard sensor, but an attached
supported environmental sensor can replace it and supply `humidity` (%RH) and
`pressure` (Pa for BME measurements). It is not automatically ambient air or the
customer's process temperature. Onboard errors can leave `temperature=0`.

Other fields include `voltage` (V), `usb` (presence), `count` (report-generation counter,
not degrees or seconds) and `sn` (attached sensor serial in the BME path, not Event
device `sn`). Configured diagnostics can add `temp_mcu`, `temp_accel`, `temp_modem`
(degrees Celsius), `cycles` (cumulative modem data-power cycles), `resets` (cumulative
hard-reset count) and `milliamp_hours` (Mojo accumulator). `motion` can be constructed
as a motion-event delta but is absent from some default templates; verify the emitted
schema. Counters have their own initialization/reset scope and are not all lifetime
counts. Emission can be periodic or change-triggered; silence is not necessarily a
missing heartbeat.

## Button and AUX GPIO body

`power` is a boolean for USB presence in this body, despite the different field name.
It is not watts. GPIO reporting can add `state[]` entries with `high`, `low` and
`count[]`. Positions identify configured AUX pins and mode-specific count slots;
read the corresponding pin/mode configuration before mapping them to a physical
switch or counting presses. A low pin is not universally an unpressed switch because
wiring and polarity vary. Some configurations report GPIO state through environment
variables instead of this Notefile. Do not promise a note for every physical edge or
identify double/long presses solely from the default body.

Public baseline: [System Notefiles](https://dev.blues.io/api-reference/system-notefiles/).
