# Firmware: reported version, desired version and update attempts

“What is running?”, “what should be running?” and “did the update succeed?” require
different evidence. Keep **Notecard firmware** and **host firmware** separate. A host
version is only as current and meaningful as what the host has reported to the Notecard.
An uploaded binary or desired environment value is not an installed version.

## The public structures

| Surface / field | Semantic meaning |
|---|---|
| Device `firmware_notecard`, `firmware_host` | Latest reported version strings known to Notehub, not a fresh interrogation of hardware. A host can report a plain string or richer version metadata. |
| Event `firmware_notecard`, `firmware_host` | Versions known when that context was attached. A session-begin Event can precede processing of newly reported version state; delayed application Events need not carry their capture-time firmware. |
| Device `dfu.card`, `dfu.user` | Most recently reported update state for Notecard and host respectively. `user` means host firmware here, not a Notehub user account. Missing object means no state exposed, not “never updated”. |
| DFU status/history `current` | Last reported running version and available build metadata for the selected target. This is the place to compare the reported version with the target. |
| Status `status` / history `history[]` | Latest requested attempt / list of attempts. Within an attempt, `current_version` is the **version before that request**; it is not synonymous with top-level `current.version`. |
| Attempt `requested_version` | Desired version derived from the selected upload's metadata. May fall back to a filename, a deleted-file indication or unknown. Not necessarily a comparable semantic version. |
| Attempt `initiated` | Time the request/attempt was established, RFC3339; late reports can affect reconstructed timing. Not proof download began then. |
| Attempt `updates[]` | Timeline entries with `phase`, explanatory `description`, `status` text and RFC3339 `datetime`. Can include service-synthesized pending, download-started/completed, installing or cancellation entries. |
| Status `dfu_in_progress` | Service assessment of ongoing work. False alone is not proof of success: an error, cancellation or another nonrunning state can also be terminal. |

`current.version` is the reported version; `current.organization`, `description`,
`product`, `built` and `builder` are firmware build metadata when available.
**`current.product` is a firmware product label, not the Notehub ProductUID**, despite
terse descriptions that can suggest otherwise. `built` describes the build, not the
installation date. Metadata can be missing even when a version string is present.

For example, top-level `current.version="2.0"` with an attempt's
`current_version="1.0"` and `requested_version="2.0"` is consistent with an update
from 1.0 to a now-reported 2.0. The attempt's starting version does not contradict the
latest report.

The Device API's `dfu.card` / `dfu.user` fields have a different shape:

| Field | Meaning |
|---|---|
| `type`, `file` | Firmware/image type and selected file identifier/name reported for the operation. Neither is proof of what is currently running. |
| `length`, `crc32`, `md5` | Image size in bytes and integrity identifiers, not version or rollout percentage. |
| `mode` | Reported update phase after API normalization; see below. |
| `status` | Human-readable device progress/error text. Do not treat arbitrary text as a stable enum. |
| `began`, `updated` | Device-reported operation start and last update times, Unix seconds. Clock validity and delayed sync affect interpretation. |
| `retry`, `errors` | Retry identifier/count associated with the attempt, and consecutive errors. `errors` is not a lifetime failure count. The same file with a different retry value can be a new attempt. |
| `read` | Bytes read from the service for the operation; download progress is not installation confirmation. Interpret against `length` only when both describe the same image/attempt. |
| `version` | Reported running version information, encoded as a string; can itself contain JSON metadata. It is not automatically the file's requested version. Parse conditionally, retain the original if plain text. |

Recognized raw modes include `idle`, `error`, `downloading`, `sideloading`, `ready`,
`ready-retry`, `updating` and `completed`. `ready` means an image is ready for the next
step; `completed` must be interpreted for the target and reconciled with running-version
evidence. Unknown/newer modes can become an empty mode through API normalization.
Neither an empty mode nor missing data means success.

Do not collapse the raw Device DFU object and the status/history timeline into one
enum. The latter synthesizes phases such as `pending`, `download_started`,
`download_completed`, `installing` and `canceled`. Their text may be lowercased by
the API. The internal service knows additional firmware targets and state fields
that the public Device converter does not expose; do not invent public fields or
endpoints for modem/other firmware from that internal knowledge.

## Answering rollout questions

For “which devices are still on the old firmware?”, resolve the target (host or
Notecard), get latest reported versions and last activity, then compare within the
appropriate product/hardware cohort. Keep unreported and stale devices separate from
confirmed old versions. Parse known version formats; lexical string sorting can put
`10.x` before `9.x`, and some host strings are not versions suitable for ordering.

For “did this update work?”, identify the attempt and target, inspect its timeline,
compare top-level `current.version` with the intended image/version, and check the
reporting time. Say “requested”, “awaiting sync”, “downloading”, “downloaded/ready”,
“reported running target”, “failed” or “not yet confirmed” according to the evidence.
Downloading 100%, a finished attempt or `dfu_in_progress=false` alone is insufficient.
Use `_health.qo`/`_health_host.qo` and connection history to explain gaps where available.

For “why is it stuck?”, compare the last phase/time, effective desired firmware
configuration, connectivity, compatible hardware, supply conditions and host update
mechanism. The host may need to consume/install a downloaded image; different DFU
mechanisms have different handoffs. Do not assume the Notecard being online proves
the host is participating. Read the product's skill for that mechanism.

Desired firmware commonly uses `_fw` for host and `_fwc` for Notecard. Read effective
configuration across scopes and compare with reported state. Answering these questions
does not itself authorize changing environment variables, retrying DFU, uploading an
image, rebooting or canceling an attempt.

Endpoint and procedure references: [Project API](https://dev.blues.io/api-reference/notehub-api/project-api/),
[Notecard firmware updates](https://dev.blues.io/notecard/notecard-walkthrough/updating-notecard-firmware/),
and [host DFU](https://dev.blues.io/notehub/host-firmware-updates/notecard-api-requests-for-dfu/).
