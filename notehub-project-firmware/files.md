# Manage firmware files

Use this for a requested upload, download, metadata edit or deletion. Reuse the selected
ProjectUID and authenticated MCP or CLI connection. File management changes the firmware
catalog; Device updates follow the separate firmware cookbook.

## 1. Identify and inspect the file

Use `GetFirmwareInfo` to inspect the relevant entries. Match firmware type, target,
version and metadata to the request. For an existing file, retain the full stored
filename returned by the API, including any timestamp suffix; the displayed version
is not a file identifier. For an upload, identify the supplied binary and its intended
firmware type and metadata.

## 2. Choose and carry out the requested operation

Use the current contract and selected tool's binary-transfer facilities. Reuse the user's
authorization for an upload, edit or deletion. Before deleting a file, check whether the
request accounts for existing desired-version settings or pending updates that refer to
it; resolve an unintended consequence before making that change.

| Request | Operation | Apply |
|---|---|---|
| Download the selected binary | `DownloadFirmware` | Retrieve the exact stored file through the authenticated tool |
| Upload a supplied binary | `UploadFirmware` | Send binary content using the documented request type; retain the returned filename and metadata |
| Edit an existing host firmware entry's metadata | `UpdateFirmware` | Supply only the intended metadata changes; this edits the catalog entry, not running Device firmware |
| Delete a stored host binary | `DeleteFirmware` | Use its exact full stored filename and inspect the result |
| Resolve parameters, supported firmware types or binary request/response details | [OpenAPI specification](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-access/openapi.yaml) | Follow the operation's current contract through the established proxy |

Do not upload a different image, initiate an update or remove other files as a side effect
of catalog work. After an uncertain write response, inspect the catalog before retrying.

## 3. Verify the file or catalog result

Read back the uploaded or edited entry and compare its filename, relevant metadata and
available integrity information with the intended result. Verify a download's returned
file and a deletion's removal from the catalog. An empty listing after a failed request
does not verify deletion.

## 4. Return the result

Return the exact filename, completed operation and any limitation to the calling workflow.
Uploading a binary leaves Device installation unchanged. If the user also requested a
rollout, continue with the selected image in the firmware cookbook.

| Next requested work | Procedure | Result |
|---|---|---|
| Deploy the verified image to Devices or Fleets | [Firmware updates](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-project-firmware/SKILL.md#1-identify-the-question-and-target-devices) | Resolve the target population and carry the authorized update through verification |
| Continue the conversation | [Project Chat](https://raw.githubusercontent.com/blues/skills/refs/heads/master/notehub-project-chat/SKILL.md#3-explain-the-verified-result) | Explain the verified result and retain context for follow-ups |
