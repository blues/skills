# Location: source, time and uncertainty

For “where is it now?”, report the latest known location **with its observation time
and source**. For “where was this reading taken?”, use that Event's location context
and compare its fix time with capture time. Today's Device object cannot supply the
historical location of all its Events.

## Choose the relevant reference

Read the rows that apply to the current step; reuse context already established.

| If you need to | Read | Then |
|---|---|---|
| Interpret journeys, fences or movement history | [Tracking](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-platform/tracking.md) | Combine tracking-body meaning with location provenance and time |

## Equivalent concepts have different paths

| Meaning | Event | Device API | Session API |
|---|---|---|---|
| Device/note-associated GPS location | `where_lat`, `where_lon`, `where_when`, `where_olc` | `gps_location.latitude`, `.longitude`, `.when` | `where_lat`, `where_lon`, `where_when`, `where` (OLC) |
| Serving tower lookup | `tower_lat`, `tower_lon`, `tower_when`, `tower_id` | `tower_location`; identity in `tower_info` | `tower` plus `cell` |
| Radio geolocation result | `tri_lat`, `tri_lon`, `tri_when`, `tri_points` | `triangulated_location` | `tri` |
| Service-selected location | `best_lat`, `best_lon`, `best_location_when`, `best_location_type` | `best_location` | No equivalent public `best_*` field set on v1 Session records |

Event families also have location name, country and timezone fields: `where_location`,
`where_country`, `where_timezone`, with parallel `tower_`, `tri_` and `best_` names.
Names/countries/timezones are geographic lookup labels, not a customer's installation
name or guaranteed street address. `best_location` on an Event is a **name string**;
on Device it is a **nested object**.

Device location objects use `name`, `country`, `timezone`, `latitude`, `longitude` and
RFC3339 `when`. Session `tower`/`tri` use compact keys: `time` (Unix seconds), `n` (name),
`c` (country), `lat`, `lon`, `zone` (timezone name), `l` (Open Location Code), `z`
(timezone identifier), `source`, `mcc`, `mnc`, `lac`, `cid` and `towers` (lookup point
count). A legacy `count` in other serializers is deprecated; it is not a sample count.
Device `tower_info` instead spells the final identifier `cell_id`.

### Source is part of the measurement

- `where_*` carries the location attached to the note/device, commonly GPS. Explicit
  host-provided locations, external GPS and configured fixed locations can also feed
  this path. The field family alone does not prove a fresh fix by the onboard receiver.
- `tower_*` is a serving-cell lookup location, not a GPS position of the device.
  Coverage can extend far from a tower; a precise-looking decimal is not precise
  device positioning. Cell identity is the MCC/MNC/LAC/cell tuple, not cell ID alone.
- `tri_*` is a service lookup from scanned Wi-Fi access points and/or cellular radios.
  “Triangulated” does not promise geometric three-tower triangulation, a fixed radius,
  or better accuracy than every GPS fix. `tri_points`/`towers` counts lookup inputs,
  not satellites. `triangulate` and Session `scan` describe inputs, not the result.
- Location attached during upload/processing can be newer than a delayed reading.
  Always compare source timestamps to `when`; a recent tower fix at sync time cannot
  precisely locate all measurements captured offline.

## What “best” means

Notehub selects among available location families, generally GPS, then triangulated,
then tower. It is a selection policy, not an accuracy bound or a guarantee of freshness.
The selection rules include these cases:

- Device/Session GPS staleness is firmware-feature-dependent. Service thresholds
  default to three days and can be configured; older-device behavior can preserve GPS.
- A fresh tower sufficiently near the last GPS position can preserve that GPS position
  despite its age. Optional spoof avoidance can prefer a newer, distant radio result.
- When all candidates are stale/unusable, the service can fall back to an old position.
  A nonempty `best_location` is therefore not proof the device has checked in recently.
- Events deliberately retain note-associated GPS without the same age rejection used
  for current device state, because Events may have been captured offline. Do not apply
  the Device selection algorithm blindly to historical Events.

Read returned `best_location_type` (`gps`, `triangulated`, `tower`, or empty/absent
when unavailable) and time where present. If Device does not expose
the type, avoid assigning provenance from equality alone when several locations match;
show the underlying candidates when the distinction affects the answer.

## Tracking coordinates and derived travel

The compact tracking template uses reserved `_lat`, `_lon`, `_ltime`, `_time` slots
to transport metadata. Notehub lifts them into the Event envelope and removes them
from the decoded body. Query `where_lat/where_lon/where_when`, not presumed
`body._lat/body._lon` or `body.lat/body.lon`. A custom host body may independently
contain coordinates with its own meaning.

Privacy/fuzz configuration can coarsen visible coordinates. Device-derived distance
and velocity can use unfuzzed fixes, so recalculating from displayed coordinates can
differ. Do not silently “correct” one with the other.

For a map, retain source/time and distinguish radio estimates from GPS points. For
distance or speed, restrict to compatible, credible fixes, positive time deltas and
acceptable gaps; explain the policy. Straight-line point-to-point distance is sampled
displacement, not road distance. Tower switching and GPS jitter can mimic travel;
mixing GPS and tower coordinates can create a false high-speed jump. Bearing from a
segment is direction of travel, not physical device orientation.

For journeys and fences, continue with tracking. For API shapes consult
the [Device API](https://dev.blues.io/api-reference/notehub-api/device-api/).
