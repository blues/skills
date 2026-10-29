---
name: notehub-project-firmware
description: Explore and manage Notecard and host firmware updates for Devices and Fleets in a selected Notehub Project. Inspect available and reported versions, investigate DFU status and history, carry out authorized updates or cancellations, and verify results through the established MCP or CLI connection.
---

# Explore and manage Firmware Updates

Use this cookbook within Project Chat for firmware questions and authorized changes.
Reuse the selected ProjectUID, MCP or CLI connection, product knowledge and any prepared
cache. Both access methods support this work. Follow
[working style](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub/working-style.md).

Keep **Notecard firmware** and **host firmware** separate. For each, distinguish the
reported running version, desired version and progress of an update attempt. An uploaded
binary or an accepted update request does not prove installation.

| The request concerns | Follow | Return with |
|---|---|---|
| Versions, rollout progress or a failure | Steps 1–2, then Step 6 | Findings from the relevant firmware evidence |
| Updating Devices, canceling updates or changing desired firmware settings | Steps 1–6 | The authorized change and verified or pending Device outcomes |
| Uploading, downloading, editing metadata or deleting a firmware file | [Manage firmware files](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-project-firmware/files.md) | The verified file or catalog operation, without initiating a Device update |

## 1. Identify the question and target Devices

Resolve whether the user wants an inventory, comparison, investigation, update or
cancellation. Identify host or Notecard firmware and the DeviceUIDs or FleetUIDs within
the selected Project. Use Product Skills to understand
customer product types, hardware variants and the host's update mechanism.

For a Fleet, inspect current membership and applicable ProductUID or hardware filters.
Devices can belong to multiple Fleets, and Smart Fleet membership can change. Keep the
actual Device population and observation time with the answer. A request to inspect a
Fleet's firmware does not request changing it.

**Complete when:** the firmware type, Project and intended population are clear.

## 2. Inspect versions, desired configuration and update progress

Use cached versions for an overview when sufficiently fresh. For current update status,
troubleshooting or pre-change verification, read the relevant API state through the
selected connection. Use these operations from the OpenAPI contract:

| Need | Operation or reference | Establish |
|---|---|---|
| Available firmware files and metadata | `GetFirmwareInfo` | Exact stored filenames, firmware type, target, version and relevant metadata |
| One Device's update status or history | `GetDeviceDfuStatus`, `GetDeviceDfuHistory` | Current report and the relevant attempt for `host` or `notecard` |
| Status or history across Devices or a Fleet | `GetDevicesDfuStatus`, `GetDevicesDfuHistory` | The selected population's outcomes, using the documented Device, Fleet, ProductUID and other filters |
| Desired firmware and configuration inherited by a Device | `GetDeviceEnvironmentHierarchy`; the relevant Device, Fleet or Project environment reads | Effective settings and which scope supplies them |
| Interpret firmware reports and DFU evidence | [Firmware meanings](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-platform/firmware.md) | Correctly distinguish running version, requested version, progress and missing evidence |
| Resolve exact parameters and request bodies | [OpenAPI specification](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-access/openapi.yaml) | The current contract for the operation being used |

Read all matching batches needed for a complete population comparison. For repeated
reads, follow the endpoint's modified-cursor contract. The DFU collection endpoints
explicitly warn that Fleet or Project firmware-setting changes do not appear in cursor
results until the Device itself is next updated. After a configuration change, read back
the relevant settings and obtain current status directly; an unchanged cursor scan does
not establish that the desired configuration is unchanged.

Compare versions within compatible product and hardware groups. Keep unreported or stale
Devices separate from Devices confirmed to run an older version. Firmware metadata's
`product` label is not ProjectUID or ProductUID. For an investigation, relate the update
timeline to connectivity, last activity, effective configuration and the product's host
update mechanism; use the Platform reference for the meaning of each field.

**Complete when:** the relevant versions, configuration, progress and evidence gaps are
understood. For a question or investigation, go to Step 6 with the findings. Requested
Device updates or configuration changes continue through Step 3.

## 3. Define the authorized change

Resolve the exact image and intended Device population before an update. Match the
image's firmware type and hardware target, and any product compatibility requirements.
Use the stored filename returned by the API, preserving suffixes; a displayed version
label alone is not a file identifier. For host firmware, use the product's documented
installation mechanism and prerequisites.

Distinguish an action on selected Devices from changing a Fleet's ongoing desired
firmware settings. Fleet settings can affect Devices that join later. Inspect inherited
settings and overrides before changing their owning scope. Recheck membership before a
bulk action; choose explicit DeviceUIDs if the request covers only the observed members.
Never let a missing filter expand a Fleet or Device request to the entire Project.

Reuse the user's authorization for the stated image, population and action. Ask only
when a consequential target, image or scope is unresolved or falls outside that request.
Apply an existing staged-rollout policy when one is part of the product's procedures or
the user's request; preserve its criteria for proceeding or stopping.

**Complete when:** the exact change and affected scope are understood and authorized.

## 4. Perform the requested change

Choose the operation for the actual request, through the established authenticated
proxy. Verify the current parameter and body definitions before calling it.

| Requested change | Operation or procedure | Verify immediately |
|---|---|---|
| Request an update for selected Devices, including a Fleet-filtered population | `PerformDfuAction` with `action=update`, the correct `firmwareType`, explicit filters and `filename` body | Acceptance for the intended scope, then read current status |
| Cancel requested updates | `PerformDfuAction` with `action=cancel` and the intended scope | Current request/configuration and Device status; cancellation does not establish a rollback |
| Change ongoing desired firmware or update windows at Device or Fleet scope | The documented environment-variable operations and [DFU configuration](https://dev.blues.io/notehub/host-firmware-updates/environment-variables-for-dfu/) | Settings at the changed scope and their effective values for affected Devices |

Desired host and Notecard firmware use `_fw` and `_fwc`, respectively. Preserve unrelated
environment values and existing timing or power conditions. Use the documented semantics
for the chosen scope; a cancellation, removed setting and completed rollback are different
outcomes. A rollback requires a compatible prior image and an authorized update request.

Keep the target population, image, action time and returned result for verification. If a
request times out, read the resulting state before retrying. Stop or narrow work when an
unexpected result invalidates the intended scope; report partial results accurately.

**Complete when:** the requested service change has been accepted and read back, or a
specific failure is understood. Device installation may still be pending.

## 5. Verify the Device outcomes

Read status and, when useful, history for the selected firmware type. Compare the latest
reported `current.version` with the intended image and inspect the relevant attempt and
reporting time. The attempt's `current_version` describes its starting version. Neither
100% download, an inactive `dfu_in_progress`, nor silence proves that the target is running.

Summarize counts and exceptions: reported running target, awaiting contact, downloading,
ready/installing, failed, canceled or unconfirmed, as supported by actual evidence. Use
the same population as the request and identify any later membership changes. Investigate
failures using Step 2; do not silently retry, change the image or expand the rollout.

Updates can span future Device contacts. Report the current verified state and remaining
work. If the user requests ongoing notification, reuse the normal monitor procedure with
this rollout's population, target and completion rules.

**Complete when:** verified outcomes and pending or failed Devices are distinguished.

## 6. Return the findings to Project Chat

Give the caller the answer or change made, identifying host or Notecard firmware,
Project, Device/Fleet scope, image/version and evidence time. Keep exceptions, pending
work and the calculation or operation context for the next question. In the Chat loop,
continue at Step 3 to explain the result.

| Continue with | Procedure | Result |
|---|---|---|
| Follow-up questions or another requested action | [Project Chat](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-project-chat/SKILL.md#3-explain-the-verified-result) | Carry the firmware findings into the ongoing conversation |
| Notify on progress, failures, completion or a periodic summary | [Create or change a monitor](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-project-chat/monitors.md#create-or-change-a-monitor) | Watch the rollout and notify according to the user's preferences |
