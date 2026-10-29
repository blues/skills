# Sessions, connectivity and data usage

## Choose the relevant reference

Read the rows that apply to the current step; reuse context already established.

| If you need to | Read | Then |
|---|---|---|
| Difference high-power counters or interpret charging and Mojo fields | [Power](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-platform/power.md) | Apply reset, generation and charging semantics |
| Interpret session tower, triangulated or GPS location | [Location](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-platform/location.md) | Preserve each source and its observation time |

## Which record answers the question

Use full Session records for connection history and per-session counters. Use
`_session.qo` Events for the corresponding event-stream evidence and joins to application
Events. Use the Device object for latest known activity. These surfaces overlap, but
neither a session-start Event nor a Device snapshot contains a complete Session record.

For a repository or cache, inspect whether the dataset contains Session objects or just `_session.qo`
Events. In particular, high-power counters available in the Session API need not be
present on session-begin Events. Use the API when required fields or history are absent.

## Session field meanings

Each Session has a SessionUID and is tagged with the DeviceUID of the one Device that
generated it.

| Fields | Meaning and unit |
|---|---|
| `session`, `device`, `sn`, `product`, `fleets` | SessionUID, DeviceUID, device serial number, ProductUID and FleetUIDs. Join Events by ProjectUID, DeviceUID and SessionUID. Count distinct sessions, not start plus end rows. |
| `session_began`, `session_ended` | Session open and close times, Unix seconds. Missing/zero close is not a zero-duration completed session. |
| `why_session_opened`, `why_session_closed` | Descriptive reasons. Preserve the text and its context; do not assume a permanent exhaustive enum. |
| `continuous`, `tls`, `transport` | Connection mode, TLS use and transport description. Persistent/packet transports need not follow one TCP connection per session. |
| `when` | Last known capture time of a note routed through this session. Not the session start. |
| `work`, `events` | Last service work time (Unix seconds) and number of Events routed through the session. Neither is a count of physical sensor measurements. |
| `voltage`, `temp` | Supply volts and device temperature in degrees Celsius, reported snapshot. |
| `moved`, `orientation` | Time of last reported accelerometer motion (Unix seconds), and orientation label relative to the device. `moved` is **not a motion count**; orientation is **not compass bearing**. |
| `hp_secs_total`, `hp_secs_data`, `hp_secs_gps`, `hp_cycles_total`, `hp_cycles_data`, `hp_cycles_gps` | Cumulative high-power activity snapshots at session start. See power before differencing. |
| `power_charging`, `power_usb`, `power_primary`, `power_mah` | Charging/USB/configured primary-battery indications and Mojo accumulator. See power. |
| `failed_connects`, `penalty_secs` | Accumulated unsuccessful connection attempts and associated penalty seconds before this successful communication, not failures during the row's session. Reset boundaries and firmware matter. |
| `cell` | Serving-cell identity encoded as `mcc,mnc,lac,cellid`. Event session metadata instead uses `cellid`. It is not a coordinate. |
| `iccid`, `apn`, `ip` | SIM identity, access point name and connection network address. ICCID can change across SIM selection; IP is not a physical location. |
| `bssid`, `ssid` | Wi-Fi access point identity and network name, when applicable. SSID alone is not a unique site identifier. |
| `rat`, `bearer` | Radio access technology and bearer description. Compare signal and usage within compatible technologies. |
| `scan`, `triangulate` | Scanned-radio input to location lookup, when exposed; `scan` can be an encoded string. Not a completed location result. |
| `tower`, `tri`, `where*` | Session tower/triangulated and last-known note GPS context. See location. |
| `period`, `usage_actual` | Usage during the accounting period and whether device/transport actuals have replaced the service estimate. See below. |

Internal or other transport serializers may expose `session_updated`, `time_bounded`,
`packet_usage` or socket fields. The public v1 Session converter does not necessarily
expose those. Do not build a query from an internal model without verifying the actual
API/dataset. Internal lifetime snapshots `this` and `next` are not the public `period`.

### Signal quality

| Field | Interpretation |
|---|---|
| `rssi` | Received signal strength in dBm; less negative generally means stronger. Applicable meaning depends on radio technology. |
| `rsrp` | LTE reference-signal received power in dBm. |
| `rsrq` | LTE reference-signal received quality in dB. |
| `sinr` | Signal to interference-plus-noise ratio in dB on modern firmware. Some firmware before build 15323 / 3.5.2 reported raw modem values 0–250; the documented source correction for that generation is `-20 + raw/5`. Verify the version at the measurement before converting. |
| `bars` | Device-reported or service-derived ordinal signal category (typically 1–4); 0/absence can mean unavailable. Not a linear physical measurement. |

Do not average bars as though they were dBm, or combine raw legacy SINR with scaled
SINR. A zero often means unavailable in these encodings and can also be a real SINR;
check companion fields and firmware. A strong signal on the last successful connection
cannot describe signal during a later outage. Successful-session history has selection
bias: failures without a connection cannot send their own Session records.

For “which devices have poor coverage?”, compare compatible RAT/modem cohorts over a
specified time window, report samples and recency, and use failures/penalties and sync
delay as additional evidence. Do not diagnose the antenna from RSSI alone.

## Usage counters and their lifecycle

`period.since` is the period start in Unix seconds; `period.duration` is its recorded
duration in seconds. `period.bytes_sent` is uplink from the device and `bytes_rcvd` is
downlink to it. For cellular traffic these totals **already include
secondary-SIM traffic**. `bytes_sent_secondary` and `bytes_rcvd_secondary` break out
the secondary portion; do not add them to the totals. Preserve the breakdown and
verify this relationship for another transport or transformed dataset before deriving
primary-only usage by subtraction.

For example, 1,000 sent bytes, 500 received bytes and 200 secondary sent bytes mean
1,500 total bytes, with 200 of the sent bytes attributed to the secondary SIM. They
do not mean 1,700 total bytes.

`period.notes_sent` and `notes_rcvd` count transferred notes, not JSON fields or Event
rows. `sessions_tcp` and `sessions_tls` count connections by type, not modem cycles;
do not add them and label the sum “unique sessions” without checking their relationship
on that transport.

For ordinary sessions, usage starts as a service-side estimate. At a later connection,
device cumulative statistics can replace the previous period's byte, note and connection
counts with deltas, and `usage_actual` becomes true. These are **replacement values**,
not another period to add. A closed session can therefore change. Negative cumulative
deltas from resets are clamped to zero during service reconciliation; an actual-marked
zero does not necessarily demonstrate no traffic. Missing device totals can also limit
the correction. Persistent packet sessions can obtain actual usage by a different
transport-specific path.

Refresh recent mutable Session records before answering, rather than treating an
append-only Event cache as an immutable Session store. Keep `usage_actual` and report
provisional totals where relevant. Distinguish three quantities:

| User meaning | Evidence |
|---|---|
| Device traffic over a period | Deduplicated Session `period` counters, their transport and actual/estimated status |
| SIM plan balance / carrier-accounted consumption | Device `cellular_usage` or plans API, per ICCID and accounting window, with `last_updated` |
| Notehub service consumption or charges | Usage API's specific metric, dimensions and billing interval |

Event count times mean body size is not a traffic total: it omits protocol overhead,
downloads, retries and compression. Session counters are not a substitute for carrier
billing or the Notehub Usage API. Plans endpoints can use `kb_*` while Device
`cellular_usage` uses bytes; read the actual unit before combining them.

## `_session.qo` bodies and Event metadata

Notehub generates these Events; no host `note.add` is necessary.

| Shape | Where the useful data lives |
|---|---|
| `req="session.begin"`, `body.opened=true`, `body.why` | Start reason in the body. Radio, supply, motion, power and firmware context principally in the **Event envelope**, not the body. |
| `req="session.end"`, `body.closed=true`, `body.why` | End reason in the body. Service usage fields can include `hub_last_work_done`, `hub_duration_secs`, `hub_events_routed`, `hub_rcvd_bytes`, `hub_sent_bytes`, `hub_tcp_sessions`, `hub_tls_sessions`, `hub_rcvd_notes`, `hub_sent_notes`. Omission is not a complete zero report. |
| Device-derived usage correction, where that surface emits it | `card_rcvd_bytes`, `card_sent_bytes`, their `_secondary` variants, `card_tcp_sessions`, `card_tls_sessions`, `card_rcvd_notes`, `card_sent_notes`, associated with the prior session. Verify the emitted operation/shape instead of assuming every deployment stores a third Event. |

`hub_` versus `card_` distinguishes the reporting source, not two disjoint streams of
traffic. Prefer the refreshed Session record for the reconciled period. Do not count
each phase/correction Event as a separate connection or sum `hub_*` and `card_*`.

## “Has this device gone dark?”

Compare Device activity, Session history and customer-note capture/upload times against
the product's expected cadence, mode and permitted silence. Distinguish “has not
communicated”, “communicates but sends no fresh measurements” and “measurements arrived
late”. Without an expected cadence, report last observed activity and elapsed time;
do not invent an offline threshold. A continuous session or recent platform Event does
not prove the host application is healthy.

Endpoint details: [Device API](https://dev.blues.io/api-reference/notehub-api/device-api/)
and [Usage API](https://dev.blues.io/api-reference/notehub-api/usage-api/).
