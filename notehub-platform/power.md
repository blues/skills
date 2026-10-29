# Power cycles, charge, energy and battery life

First identify the physical quantity. A modem power cycle, a Notecard reboot, USB
removal, charger state change and host-machine power cycle are different events.
There is no universal native field for “battery life remaining”.

| Question | Evidence | Essential limitation |
|---|---|---|
| Modem on/off cycles and on-time | Session `hp_cycles_*` and `hp_secs_*` snapshots | Cumulative counters, not exact on/off timestamps |
| Notecard rebooted | `_health.qo` boot diagnostics, supported reset counters | Different from modem cycling; logs depend on capture/upload |
| USB or charging state changed | `_power.qo`, tracking `usb`/`charging`, Session power flags | Sensed source state, not whole-product power consumption |
| Charge used | Mojo-backed `power_mah` / `milliamp_hours` | mAh accumulator with reset, sign and wiring caveats |
| Energy used | Charge/current plus appropriate voltage history and measurement scope | mAh alone is not Wh |
| Battery state/lifetime | Product battery model plus relevant measured history | Supply voltage alone is not state of charge or a runtime guarantee |

## Choose the relevant reference

Read the rows that apply to the current step; reuse context already established.

| If you need to | Read | Then |
|---|---|---|
| Find full Session evidence for high-power activity | [Sessions](sessions.md) | Use the full Session API/dataset rather than assuming Event parity |
| Distinguish charge, energy and cumulative readings | [Interpretation checks](verification.md) | Apply the field's units and reset semantics |

## High-power counters

Full Session records expose a snapshot at the start of communication:

- `hp_cycles_total` and `hp_secs_total`: cumulative completed modem power cycles and
  time in seconds. Counters update when modem power is
  turned off; the current active interval may not yet be counted.
- `hp_cycles_data` / `hp_secs_data` and `hp_cycles_gps` / `hp_secs_gps`: those intervals
  tagged for data and GPS work. An interval can have both tags, so the categories can
  overlap. Do not require `data + gps = total` or count their sum as disjoint activity.
- These are cumulative persistent statistics, not per-session quantities, and can be
  explicitly reset. They are not a measurement of the host's power cycles. Hardware
  generations with different GPS/modem arrangements also need separate interpretation.

For an interval, take differences between ordered snapshots for the same device,
counter generation and reset segment. Retrieve a baseline before the desired window
and a later snapshot if needed; describe the boundary uncertainty. Do not sum cumulative
snapshots. A decrease indicates reset/wrap/changed context or another discontinuity
to investigate, not negative usage. No later communication means no later snapshot,
not zero additional modem activity. Counts alone cannot recover individual on/off times.

The `_session.qo` Event is not guaranteed to expose these high-power fields; inspect
the full Session API/dataset. See sessions.

## Mojo charge measurements

Session/Event `power_mah` and tracking/health/log
`body.milliamp_hours` derive from the Mojo charge accumulator, observed at different
times. They are **milliamp-hours (mAh)**. They are not a universal software estimate
of all energy consumed by any Notecard.

The hardware coulomb counter supplies signed charge increments using the configured
sense resistance. The accumulator can decrease while charging and can be reset
through power-counter controls/configuration. Persistence does not make it an
ever-increasing lifetime measure. Missing/zero does not prove a device used no charge:
Mojo may be absent, not initialized, reset, unreported, or template-defaulted.

Establish which load and current path the sense circuit measures. Do not label this
“whole product consumption” unless the project's wiring supports that claim. Keep
charging/discharging intervals and resets separate. For comparable observations
`Q2 - Q1` is a net charge change over the interval; it is not automatically gross
discharge. A decrease during charging must not be silently discarded as bad data.

For a supported discharge interval, average current in mA is `delta_mAh / hours`.
Energy requires voltage weighting: with sufficiently representative voltage samples,
`Wh ≈ sum(volts × delta_mAh / 1000)`. State the sampling and load assumptions. One
latest voltage multiplied by a lifetime charge total is not a measured energy history.

### Supply and charging indicators

`power_usb` is detected USB presence; `power_charging` is normalized charging indication
where sensing is available; `power_primary` reflects the configured primary-battery
mode. It is not a measured battery chemistry identification. Tracking has `usb` and
`charging` at body level. Boolean defaults/omission can hide unsupported sensing.
`voltage` is volts; `voltage_mode` classifies the configured voltage operating regime,
not battery percentage. The same physical voltage can have different significance
with different battery configuration, temperature, load and wiring.

## `_power.qo` is a source-transition log

For records produced by `card.power.notify`, the Notecard observes USB and/or carrier
charging sense independently and emits debounced source transitions. It does not produce
a generic energy-usage time series or a record of every modem on/off cycle.

| Body field | Meaning |
|---|---|
| `source` string | `usb` or `charging`. Analyze each source separately. |
| `from-state`, `to-state` string | `unknown`, `off` or `on`. Unknown means no established/observable state. |
| `minutes` integer | Rounded observed duration since the preceding report, credited across committed states. Normally the prior state's dwell; if reporting was delayed across transitions it can cover intermediate states too. It is not time in the new state or remaining battery life. |
| `voltage` number | Supply volts sampled for the notification, using Mojo voltage when active. |
| `voltage_mode` string | Current voltage-regime classification; depends on configuration. |

The state machine starts unknown after boot. `unknown → off` is establishment of the
first observed level, not evidence that an external supply failed at that moment.
Charging sense becoming unavailable can report a transition to unknown. Delayed
processing can coalesce transitions into a net state change. Short changes
can be suppressed by debounce and minimum dwell; reporting need not capture every
physical edge. Duration uses heartbeat observation, is rounded, and saturates for
very long dwell (roughly 49 days). Do not use it as a
precision outage stopwatch. Emission and sync configuration also differ: a queued
transition may arrive later.

## Battery questions need product knowledge

For “how much battery is left?”, seek a calibrated state-of-charge/coulomb model or
a product-defined voltage-to-charge relationship, battery capacity/chemistry, reset
baseline, temperature and load. If those are unavailable, report observed voltage and
charge trends, not a manufactured percentage.

For “how long will it last?”, a forecast needs usable remaining capacity and a
representative future load/duty-cycle model, with charging treated separately. Radio
conditions, retry behavior, sensing, host consumption and temperature can change it.
Use the project's documented assumptions and show uncertainty. High modem-on time can
explain a drain pattern; it cannot by itself quantify the drain in mAh.

For “why is consumption higher?”, align valid counter/charge deltas with sessions,
failed connections, signal, GPS seeks and product operating modes. Correlation may
suggest a cause; do not claim the host, battery or antenna failed from one counter.

Reference: [Notecard card requests](https://dev.blues.io/api-reference/notecard-api/card-requests/).
