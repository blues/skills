# MCP execute reference · Notehub API 1.2.0 · OpenAPI 3.0.3
# call_tool('<name>', {param: value, ...})  ·  *required  ·  header keys use the OpenAPI name exactly
# !write / !destructive on a param escalates the call when that param is used (a boolean: when true)
# each {x} in a path is a required param x, a string unless an `x  path` line defines it  ·  `a, b  query  string` = params with the same spec
# `+ a, *b` = params defined under ## common params; a param line under a call applies to that call only
# a capitalized type such as Foo is defined under ## types; an inline object's fields are indented below it

## common params
descriptionName  path  string  — Record name, 1-128 letters, digits, `.`, `_` or `-`, e.g. `_health.qo`.
startDate  query  integer(int64)  — Start date for filtering results, specified as a Unix timestamp
endDate  query  integer(int64)  — End date for filtering results, specified as a Unix timestamp
environment_variable_notes  body  map[string]string  — Optional per-variable annotations, keyed by variable name.
cell  body  string  — Cell ID where the session originated and quality ("mcc,mnc,lac,cellid")
failed_connects  body  integer(int64)  — Number of failed connection attempts in the prior session
penalty_secs  body  integer(int64)  — Number of seconds in penalty in the prior session
when  body  integer(int64)  — Last known capture time of a note routed through this session in Unix timestamp
where  body  string  — Open Location Code from last GPS location
work  body  integer(int64)  — Unix timestamp of the last time work was done for this session
webhookUID  path  string  — Webhook ID, e.g. `ttn-uplink`.
disabled  body  boolean  — If true, posts to the webhook return an error and create no events.
id  body  string  — Webhook ID, also the default notefile name (`<id>.qo`).
transform  body  string  — JSONata expression applied to the event, with the post under `body`. Its output replaces the event and must set `body`.
pageSize  query  integer(int32)  default=1000
pageNum  query  integer(int32)  default=1
cursor  query  string  — If supplied, results are returned in ascending order of modification time. A cursor is a timestamp: Unix epoch microseconds (UTC), optionally followed by `.N`.
deviceUID, fleetUID  query  array[string]
tag  query  array[string]  — Matches the whole tag, case-insensitive.
serialNumber, location, hostFirmware, sku  query  array[string]  — Matches any value containing the filter text.
notecardFirmware  query  array[string]  — Firmware version filter.
productUID  query  array[string]  — Matches any value containing the filter text, case-insensitive.
firmwareType  path  string[host|notecard]  — Which firmware: `host` (the host microcontroller's) or `notecard`.
notefileID  path  string  — Notefile name, including its extension, e.g. `config.db` or `commands.qi`.
body  body  object  — Arbitrary user-defined JSON for the note.
payload  body  string(byte)  — Optional base64-encoded payload.
noteID  path  string  — Note ID within the notefile, e.g. `my_note`.
firstSync  query  boolean  — If true, return only sessions whose `why_session_opened` begins with "first sync".
sortBy  query  string[best_id|device_serial|device_uid|captured|modified|device_location|tower_location|triangulated_location|best_location]  default=captured
sortOrder  query  string[asc|desc]  default=asc
key  path  string  — Name of the environment variable to delete.
sensorUID  query  array[string]  — SensorUID to filter by, matched exactly.
dateType  query  string[captured|uploaded]  default=captured  — Which time `startDate` and `endDate` filter on.
systemFilesOnly  query  boolean  — If true, return only events from system notefiles, such as `_session.qo` and `_health.qo`.
files  query  string  — Comma-separated notefiles, e.g. `data.qo,_health.qo`. Don't combine with `systemFilesOnly`.
format  query  string[json|csv]  default=json  — Response format.
selectFields  query  string  — Comma-separated event fields to return, e.g. `device,when,body.temp`.
smart_rule  body  string  — JSONata expression for smart fleet membership: `1` adds the device, `-1` removes it, and `0` or an error leaves it…
jobUID  path  string  — Batch job UID, as returned by `create_job`.
reportUID  path  string  — Job run UID, as returned by `run_job`.
aggregate_function  body  string[none|avg|max|min|sum|count]  — How values in `aggregate_window` are combined before the comparison.
aggregate_window  body  string  — The time window for `aggregate_function`, such as `10m` or `5h30m40s`.
alert  body  boolean  — If true, the monitor is in alert state.
alert_routes  body  array[SlackWebHookNotification|SlackBearerNotification|EmailNotification|TwilioNotification]  — Where to send notifications, by email, Slack or Twilio. Required to create a monitor.
condition_type  body  string[greater_than|greater_than_or_equal_to|less_than|less_than_or_equal_to|equal_to|not_equal_to]  — How the value is compared with `threshold`. Required for event and usage monitors; heartbeat monitors always use `greater_than_or_equal_to`.
fleet_filter  body  array[string]  — Watch only devices in these fleets (FleetUIDs).
name  body  string  — Monitor name. Required to create a monitor.
notefile_filter  body  array[string]  — Watch only events from these notefiles, e.g. `data.qo`. Not allowed for heartbeat monitors.
per_device  body  boolean  — With an `aggregate_function`, aggregate each device on its own, so one device can raise an alert.
routing_cooldown_period  body  string  — How long to wait after notifying before notifying again, such as `10m` or `5h30m`.
silenced  body  boolean  — If true, alerts will be created, but no notifications will be sent.
source_selector  body  string  — Dot path to a numeric value inside the event body, such as `temperature`, not `body.temperature`. Required for event monitors unless `aggregate_function` is `count`.
source_type  body  string[event|heartbeat|usage]  default=event  — What the monitor watches.
threshold  body  number(double)  — The value that condition_type compares against. For heartbeat monitors this is seconds of inactivity; for usage monitors it is bytes.
uid  body  string  — MonitorUID, `monitor:` followed by a UUID.
usage_scope  body  string  — For usage monitors, what usage is summed over. Values: device|fleet
usage_type  body  string  — For usage monitors, which data to watch. Values: cellular|satellite
usage_window  body  integer(int32)  — For usage monitors, the rolling window to sum usage over, in days, e.g. `30`.
monitorUID  path  string  — MonitorUID, `monitor:` followed by a UUID.
notes  body  string  — Optional free-form text for annotating the route.
type  body  string  — Route type. Values: http|proxy|mqtt|google-function|aws-lambda|aws-lambda-with-access-key|aws-sqs|aws-sqs-with-access-key|aws-sqs-fifo|aws-sqs-fifo-with-access-key|aws-sqs-with-role|aws-iot-analytics|azure-function|azure-function-with-key|azure-service-bus-with-sas-token|snowflake|snowpipe-streaming|s3-archive|slack-bearer|slack-webhook|radnote-radresp-fixed-survey|radnote-radresp-field-survey|radnote-radresp-mobile-survey|radnote-radresp-mobile-equipment-survey|radnote-missionedge-fixed-survey|thingworx|twilio|datacake|blynk|qubitro
routeUID  path  string  — RouteUID, e.g. `route:cbd20093cba58392c9f9bbdd0cdeb1a0`.
secretName  path  string  — Name of the secret, as listed by `get_project_secrets`.
skillName  path  string  — Relative .md name, case-insensitive, at most 512 UTF-8 bytes.
limit  query  integer(int32)  default=200000  — The most data points to return.
period  query  string[day|week|month]  — The length of time each data point covers.
aggregate  query  string[device|fleet|project]  default=device  — How to group the counts, such as per device or one total for the project.
skipRecentData  query  boolean  — If true, return only summarized data, without filling in the latest day.

## description

list_global_descriptions  read  GET /v1/description
  List global description records

delete_global_description  destructive  DELETE /v1/description/{descriptionName}
  Delete a global description record

get_global_description  read  GET /v1/description/{descriptionName}
  Get a global description record

set_global_description  write  POST /v1/description/{descriptionName}
  Create or replace a global description record
  *body  body  string(binary)  — application/octet-stream: The document to store, up to 10 MB.

list_project_descriptions  read  GET /v1/projects/{projectOrProductUID}/description
  List a project's description records

delete_project_description  destructive  DELETE /v1/projects/{projectOrProductUID}/description/{descriptionName}
  Delete a project description record

get_project_description  read  GET /v1/projects/{projectOrProductUID}/description/{descriptionName}
  Get a project description record

set_project_description  write  POST /v1/projects/{projectOrProductUID}/description/{descriptionName}
  Create or replace a project description record · The request's Content-Type is stored and served back by `get_project_description`; if it's missing or…
  *body  body  string(binary)  — application/octet-stream: The document to store, up to 10 MB.

## organization

get_organizations  read  GET /v1/organizations
  List your organizations

get_organization  read  GET /v1/organizations/{organizationUID}
  *organizationUID  path  string  — List with `get_organizations`.

get_organization_balance_history  read  GET /v1/organizations/{organizationUID}/balance-history
  Get an organization's daily balance
  + startDate, endDate
  *organizationUID  path  string  — List with `get_organizations`.

## device

get_device_environment_variables_by_pin  read  GET /v1/products/{productUID}/devices/{deviceUID}/environment_variables_with_pin
  Get device environment variables using a PIN
  *X-Auth-Token  header  string  — The device's PIN, used instead of an API token.

set_device_environment_variables_by_pin  write  PUT /v1/products/{productUID}/devices/{deviceUID}/environment_variables_with_pin
  Set device environment variables using a PIN
  + environment_variable_notes
  *X-Auth-Token  header  string  — The device's PIN, used instead of an API token.
  *environment_variables  body  map[string]string

get_devices  read  GET /v1/projects/{projectOrProductUID}/devices
  List devices in a project
  + pageSize, pageNum, cursor, deviceUID, tag, serialNumber, fleetUID, notecardFirmware, location, hostFirmware, productUID, sku

delete_device  destructive  DELETE /v1/projects/{projectOrProductUID}/devices/{deviceUID}

get_device  read  GET /v1/projects/{projectOrProductUID}/devices/{deviceUID}

disable_device  write  POST /v1/projects/{projectOrProductUID}/devices/{deviceUID}/disable

enable_device  write  POST /v1/projects/{projectOrProductUID}/devices/{deviceUID}/enable

get_device_environment_hierarchy  read  GET /v1/projects/{projectOrProductUID}/devices/{deviceUID}/environment_hierarchy

get_device_environment_variables  read  GET /v1/projects/{projectOrProductUID}/devices/{deviceUID}/environment_variables

set_device_environment_variables  write  PUT /v1/projects/{projectOrProductUID}/devices/{deviceUID}/environment_variables
  Add or change a device's environment variables; variables you don't send keep their values.
  + environment_variable_notes
  *environment_variables  body  map[string]string

delete_device_environment_variable  destructive  DELETE /v1/projects/{projectOrProductUID}/devices/{deviceUID}/environment_variables/{key}
  *key  path  string  — The environment variable key to delete.

delete_notefiles  destructive  DELETE /v1/projects/{projectOrProductUID}/devices/{deviceUID}/files
  files  body  array[string]  — Names of the notefiles to delete, e.g. `config.db`.

list_notefiles  read  GET /v1/projects/{projectOrProductUID}/devices/{deviceUID}/files
  List a device's notefiles and notes
  files  query  array[string]  — Notefiles to include, e.g. `config.db`. If omitted, all are listed.
  pending  query  boolean  — show only files that are pending sync to the Notecard

get_device_health_log  read  GET /v1/projects/{projectOrProductUID}/devices/{deviceUID}/health-log
  + startDate, endDate
  log_type  query  array[string[boot|dfu_completed]]  — Return only these kinds of entry. If omitted, all entries are returned.

get_device_journeys  read  GET /v1/projects/{projectOrProductUID}/devices/{deviceUID}/journeys
  List a device's journeys
  + startDate, endDate

get_device_journey  read  GET /v1/projects/{projectOrProductUID}/devices/{deviceUID}/journeys/{journeyID}
  Get a journey and its events
  + pageSize, pageNum
  *journeyID  path  integer(int64)  — Identifier of the journey, taken from the numeric `journey` field in the event body (a Unix timestamp marking the start…

get_device_latest_events  read  GET /v1/projects/{projectOrProductUID}/devices/{deviceUID}/latest
  Get a device's latest event per notefile

create_notefile  write  POST /v1/projects/{projectOrProductUID}/devices/{deviceUID}/notefiles/{notefileID}
  Create an empty notefile

get_notefile  read  GET /v1/projects/{projectOrProductUID}/devices/{deviceUID}/notes/{notefileID}
  Get the notes in a notefile
  max  query  integer(int32)  — The most notes to return.
  deleted  query  boolean  — If true, also return deleted notes.
  delete  query  boolean  !destructive  — If true, delete the notes this request returns.

add_qi_note  write  POST /v1/projects/{projectOrProductUID}/devices/{deviceUID}/notes/{notefileID}
  Queue a note for a Notecard
  + body, payload

delete_note  destructive  DELETE /v1/projects/{projectOrProductUID}/devices/{deviceUID}/notes/{notefileID}/{noteID}

get_db_note  read  GET /v1/projects/{projectOrProductUID}/devices/{deviceUID}/notes/{notefileID}/{noteID}
  delete  query  boolean  !destructive  — If true, delete the note after returning it.
  deleted  query  boolean  — If true, also return the note if it has been deleted.

add_db_note  write  POST /v1/projects/{projectOrProductUID}/devices/{deviceUID}/notes/{notefileID}/{noteID}
  Add a note to a .db notefile
  + body, payload

update_db_note  write  PUT /v1/projects/{projectOrProductUID}/devices/{deviceUID}/notes/{notefileID}/{noteID}
  Update a note in a .db notefile
  + body, payload

get_device_plans  read  GET /v1/projects/{projectOrProductUID}/devices/{deviceUID}/plans
  Get a device's cellular data plans

provision_device  write  POST /v1/projects/{projectOrProductUID}/devices/{deviceUID}/provision
  Provision a Notecard
  device_sn  body  string  — Serial number to give the device.
  fleet_uids  body  array[string] (nullable)  — FleetUIDs to put a new device in. If omitted or empty, the product's auto-provision fleets are used.
  *product_uid  body  string  — ProductUID to provision the device under, e.g. `com.blues.bridge:sensors`. Must belong to this project.

get_device_public_key  read  GET /v1/projects/{projectOrProductUID}/devices/{deviceUID}/public-key

get_device_sessions  read  GET /v1/projects/{projectOrProductUID}/devices/{deviceUID}/sessions
  List a device's sessions
  + pageSize, pageNum, cursor, startDate, endDate, firstSync

signal_device  write  POST /v1/projects/{projectOrProductUID}/devices/{deviceUID}/signal
  Send a signal to a connected Notecard
  body  body  object  — The JSON object to send to the Notecard. Must not have a `type` field, which the Notecard reserves.

get_device_public_keys  read  GET /v1/projects/{projectOrProductUID}/devices/public-keys
  List public keys of all devices
  + pageSize, pageNum, cursor

get_fleet_devices  read  GET /v1/projects/{projectOrProductUID}/fleets/{fleetUID}/devices
  List devices in a fleet
  + pageSize, pageNum, cursor, deviceUID, tag, serialNumber, notecardFirmware, location, hostFirmware, productUID, sku

get_sessions  read  GET /v1/projects/{projectOrProductUID}/sessions
  List sessions in a project
  + pageSize, pageNum, cursor, startDate, endDate, firstSync, deviceUID, fleetUID

## webhook

create_legacy_webhook_event  write  POST /v1/products/{productUID}/devices/{deviceUID}/webhook-event
  Add a note-shaped event from a device (legacy)
  *body  body  object  — A Note-shaped event.

update_legacy_webhook_session  write  PUT /v1/products/{productUID}/devices/{deviceUID}/webhook-session
  Open or close a webhook device's session (legacy)
  body  body  object  — Session fields, such as `why_session_opened`, `why_session_closed` and signal or location values.

create_webhook_device_event_by_product  write  POST /v1/products/{productUID}/webhooks/{webhookUID}/devices/{deviceUID}/event
  Send a device's event through a webhook
  *body  body  object  — Any JSON object; it becomes the event body, or the transform's input.

create_webhook_event_by_product  write  POST /v1/products/{productUID}/webhooks/{webhookUID}/event
  Send an event through a webhook
  *body  body  object  — Any JSON object; it becomes the event body, or the transform's input.

get_webhook_settings_by_product  read  GET /v1/products/{productUID}/webhooks/{webhookUID}/settings
  Get a webhook by ProductUID

update_webhook_settings_by_product  write  PUT /v1/products/{productUID}/webhooks/{webhookUID}/settings
  Update a webhook by ProductUID
  + disabled, id, transform

get_webhooks  read  GET /v1/projects/{projectOrProductUID}/webhooks
  List webhooks

delete_webhook  destructive  DELETE /v1/projects/{projectOrProductUID}/webhooks/{webhookUID}

get_webhook  read  GET /v1/projects/{projectOrProductUID}/webhooks/{webhookUID}

create_webhook  write  POST /v1/projects/{projectOrProductUID}/webhooks/{webhookUID}
  + disabled, id, transform

update_webhook  write  PUT /v1/projects/{projectOrProductUID}/webhooks/{webhookUID}
  Send every field: a missing `transform` is cleared and a missing `disabled` enables the webhook.
  + disabled, id, transform

## external devices

create_event_ext_device  write  POST /v1/products/{productUID}/ext-devices/{deviceUID}/event
  Add an event from an external device
  app  body  string  — AppUID (globally unique)
  bars, rsrp, rsrq, sinr  body  integer(int32)  — Only available on _session.qo events.
  best_country, best_location, best_timezone, tower_country, tower_id, tower_location, tower_timezone, tri_country, tri_location, tri_timezone, where_country, where_location, where_timezone  body  string
  best_id  body  string  — The device serial number, or the DeviceUID if the serial number is not set
  best_lat, tower_lat, tri_lat, where_lat  body  number(double)  — Latitude
  best_location_type  body  string  — One of "gps", "triangulated", or "tower"
  best_location_when, tower_when, tri_when, where_when  body  integer(int64)  — Unix timestamp
  best_lon, tower_lon, tri_lon, where_lon  body  number(double)  — Longitude
  body  body  object  — A JSON object containing event details
  bssid, ordering_code, rat, sku, ssid  body  string  — Only available on _session.qo events.
  deleted  body  boolean  — True if this event has been deleted.
  device  body  string  — DeviceUID (globally unique)
  environment  body  object  — Routed environment variables beginning with "$".
  event  body  string  — EventUID (globally unique)
  file  body  string  — The notefile associated with this event
  moved  body  integer(int64)  — The number of times the device was sensed to have moved between the last session and this session.
  note  body  string  — The note ID in the notefile
  orientation  body  string  — The orientation of the device.
  payload  body  string  — A base64-encoded binary payload
  platform  body  boolean  — Whether this is a platform event (administrative, e.g. _health.qo, _session.qo) rather than user data
  product  body  string  — ProductUID (globally unique)
  received  body  number(double)  — The unix timestamp when the event was received
  req  body  string  — The notecard request
  rssi  body  integer(int32)  — Received Signal Strength Indicator (RSSI) is an estimated measurement of how well a device can receive signals.
  sensor  body  string  — SensorUID, for events originating from a sensor attached to the device
  session  body  string  — SessionUID (globally unique)
  sn  body  string  — The device serial number
  temp  body  number(double)  — Device temperature.
  tls  body  boolean  — Whether TLS was used on the connection between the device and notehub.
  transport  body  string  — The transport used for this event, e.g., "cellular", "wifi", ", etc.
  tri_points  body  integer(int32)  — Triangulation points
  updates  body  integer(int32)
  voltage  body  number(double)  — Device voltage.
  when  body  integer(int64)  — When the event was captured on the device
  when_ms  body  integer(int64)  — Millisecond-accurate Unix epoch timestamp for when the event was captured on the device
  where_olc  body  string  — Open Location Code
  wire_length  body  integer(int64)  — Approximate number of bytes that the note's body and payload occupied on the wire when transmitted from the device to…

ext_device_session_close  write  POST /v1/products/{productUID}/ext-devices/{deviceUID}/session/close
  Close an external device's session
  + cell, failed_connects, penalty_secs, when, where, work
  apn, bearer, bssid, iccid, ip, orientation, rat, ssid, where_country, where_location, where_timezone  body  string
  bars, rsrp, rsrq, rssi, sinr  body  integer(int)
  continuous  body  boolean  — Was this a continuous connection?
  device  body  string  — DeviceUID
  events  body  integer(int64)  — Number of events routed
  fleets  body  array[string]  — Array of FleetUIDs
  hp_cycles_data, hp_cycles_gps, hp_cycles_total, hp_secs_data, hp_secs_gps, moved  body  integer(int64)
  hp_secs_total  body  integer(int64)  — Total number of seconds in high power mode
  period  body  DeviceUsage
  power_charging, power_primary, power_usb, usage_actual  body  boolean
  power_mah, temp, voltage, where_lat, where_lon  body  number(double)
  product  body  string  — ProductUID
  scan  body  string(byte)
  session  body  string  — SessionUID
  session_began  body  integer(int64)  — UNIX timestamp of session start
  session_ended  body  integer(int64)  — UNIX timestamp of session end
  sn  body  string  — Device Serial Number
  tls  body  boolean  — Was TLS used?
  tower, tri  body  TowerLocation
  transport  body  string  — Type of network transport
  triangulate  body  object
  where_when  body  integer(int64)  — Unix timestamp of last GPS location
  why_session_closed  body  string  — Reason for session closing
  why_session_opened  body  string  — Reason for session opening

ext_device_session_open  write  POST /v1/products/{productUID}/ext-devices/{deviceUID}/session/open
  Open an external device's session
  + cell, failed_connects, penalty_secs, when, where, work
  apn, bearer, bssid, iccid, ip, orientation, rat, ssid, where_country, where_location, where_timezone  body  string
  bars, rsrp, rsrq, rssi, sinr  body  integer(int)
  continuous  body  boolean  — Was this a continuous connection?
  device  body  string  — DeviceUID
  events  body  integer(int64)  — Number of events routed
  fleets  body  array[string]  — Array of FleetUIDs
  hp_cycles_data, hp_cycles_gps, hp_cycles_total, hp_secs_data, hp_secs_gps, moved  body  integer(int64)
  hp_secs_total  body  integer(int64)  — Total number of seconds in high power mode
  period  body  DeviceUsage
  power_charging, power_primary, power_usb, usage_actual  body  boolean
  power_mah, temp, voltage, where_lat, where_lon  body  number(double)
  product  body  string  — ProductUID
  scan  body  string(byte)
  session  body  string  — SessionUID
  session_began  body  integer(int64)  — UNIX timestamp of session start
  session_ended  body  integer(int64)  — UNIX timestamp of session end
  sn  body  string  — Device Serial Number
  tls  body  boolean  — Was TLS used?
  tower, tri  body  TowerLocation
  transport  body  string  — Type of network transport
  triangulate  body  object
  where_when  body  integer(int64)  — Unix timestamp of last GPS location
  why_session_closed  body  string  — Reason for session closing
  why_session_opened  body  string  — Reason for session opening

## project

get_project_by_product  read  GET /v1/products/{productUID}/project
  Get the project for a ProductUID
  *productUID  path  string  — ProductUID, e.g. `com.blues.airnote`, with or without the `product:` prefix.

get_projects  read  GET /v1/projects
  List projects

create_project  write  POST /v1/projects
  *billing_account_uid  body  string  — UID of the organization that owns and is billed for the project.
  *label  body  string  — Project name shown in Notehub.

delete_project  destructive  DELETE /v1/projects/{projectOrProductUID}

get_project  read  GET /v1/projects/{projectOrProductUID}

get_aws_role_config  read  GET /v1/projects/{projectOrProductUID}/aws-role-config
  Get AWS role configuration for role-based authentication

clone_project  write  POST /v1/projects/{projectOrProductUID}/clone
  The new project's name and organization, and what to leave out of the copy.
  *billing_account_uid  body  string  — UID of the organization that owns the new project.
  disable_clone_fleets  body  boolean  — If `true`, the new project gets one default fleet instead of copies of this project's fleets.
  disable_clone_routes  body  boolean  — If `true`, the new project starts with no routes.
  *label  body  string  — Name of the new project.

get_device_dfu_history  read  GET /v1/projects/{projectOrProductUID}/devices/{deviceUID}/dfu/{firmwareType}/history
  Get a device's firmware update history

get_device_dfu_status  read  GET /v1/projects/{projectOrProductUID}/devices/{deviceUID}/dfu/{firmwareType}/status
  Get a device's firmware update status

delete_device_from_fleets  write  DELETE /v1/projects/{projectOrProductUID}/devices/{deviceUID}/fleets
  Remove a device from fleets
  *fleet_uids  body  array[string]  — FleetUIDs to remove the device from.

get_device_fleets  read  GET /v1/projects/{projectOrProductUID}/devices/{deviceUID}/fleets
  List a device's fleets

add_device_to_fleets  write  PUT /v1/projects/{projectOrProductUID}/devices/{deviceUID}/fleets
  *fleet_uids  body  array[string]  — FleetUIDs to add the device to.

perform_dfu_action  write  POST /v1/projects/{projectOrProductUID}/dfu/{firmwareType}/{action}
  Start or cancel firmware updates · With no filters, that is every device in the project.
  + deviceUID, tag, serialNumber, notecardFirmware, location, hostFirmware, productUID, sku
  *action  path  string[update|cancel]  — `update`: start updating to the firmware in `filename`; `cancel`: cancel any update in progress
  fleetUID  query  string
  filename  body  string  — Stored firmware filename to install, as listed by `get_firmware_info`. Required for `update`; ignored for `cancel`.

get_devices_dfu_history  read  GET /v1/projects/{projectOrProductUID}/dfu/{firmwareType}/history
  List devices' firmware update history
  + pageSize, pageNum, cursor, sortBy, sortOrder, deviceUID, tag, serialNumber, notecardFirmware, location, hostFirmware, productUID, sku
  fleetUID  query  string

get_devices_dfu_status  read  GET /v1/projects/{projectOrProductUID}/dfu/{firmwareType}/status
  List devices' firmware status
  + pageSize, pageNum, cursor, sortBy, sortOrder, deviceUID, tag, serialNumber, notecardFirmware, location, hostFirmware, productUID, sku
  fleetUID  query  string

get_project_environment_hierarchy  read  GET /v1/projects/{projectOrProductUID}/environment_hierarchy

get_project_environment_variables  read  GET /v1/projects/{projectOrProductUID}/environment_variables

set_project_environment_variables  write  PUT /v1/projects/{projectOrProductUID}/environment_variables
  Add or change project environment variables; variables you don't send keep their values.
  + environment_variable_notes
  *environment_variables  body  map[string]string

delete_project_environment_variable  destructive  DELETE /v1/projects/{projectOrProductUID}/environment_variables/{key}

get_firmware_info  read  GET /v1/projects/{projectOrProductUID}/firmware
  List firmware
  product  query  string  — Only firmware whose product contains this text.
  firmwareType  query  string  — Only list `host` or `notecard` firmware.
  version  query  string  — Only firmware whose version contains this text.
  target  query  string  — Only firmware whose target contains this text.
  filename  query  string  — Only firmware whose stored filename contains this text.
  md5  query  string  — Only firmware whose MD5 hash contains this text.
  unpublished  query  boolean  — If `true`, also list unpublished Notecard firmware.
  sortBy  query  string[created|name|version|length]  default=created  — Field to sort by
  sortOrder  query  string[asc|desc]  default=desc  — Sort order (asc for ascending, desc for descending)

delete_firmware  destructive  DELETE /v1/projects/{projectOrProductUID}/firmware/{firmwareType}/{filename}
  Delete host firmware
  *firmwareType  path  string[host]

download_firmware  read  GET /v1/projects/{projectOrProductUID}/firmware/{firmwareType}/{filename}
  Download host firmware

update_firmware  write  POST /v1/projects/{projectOrProductUID}/firmware/{firmwareType}/{filename}
  Update host firmware metadata · Change a host firmware file's `version`, `notes` or `info`; fields you leave out or send empty keep their values.
  *firmwareType  path  string[host]
  info  body  object  — JSON metadata for this firmware.
  notes  body  string  — New notes for this firmware.
  version  body  string  — New version string.

upload_firmware  write  PUT /v1/projects/{projectOrProductUID}/firmware/{firmwareType}/{filename}
  *filename  path  string  — Name to store the file under.
  version  query  string  — Firmware version (optional).
  notes  query  string  — Optional notes describing what's different about this firmware version
  *body  body  string(binary)  — application/octet-stream: contents of the firmware binary

get_fleets  read  GET /v1/projects/{projectOrProductUID}/fleets
  List fleets

create_fleet  write  POST /v1/projects/{projectOrProductUID}/fleets
  + smart_rule
  connectivity_assurance  body  FleetConnectivityAssurance (nullable)
  label  body  string  — Fleet name shown in Notehub. Required.
  smart_rule_enabled  body  boolean  — If `true`, `smart_rule` adds and removes devices.

delete_fleet  destructive  DELETE /v1/projects/{projectOrProductUID}/fleets/{fleetUID}

get_fleet  read  GET /v1/projects/{projectOrProductUID}/fleets/{fleetUID}

update_fleet  write  PUT /v1/projects/{projectOrProductUID}/fleets/{fleetUID}
  Always send `smart_rule_enabled`; leaving it out turns the smart rule off.
  + smart_rule
  addDevices  body  array[string]  — DeviceUIDs to add to the fleet. Each must be in the project and not also in `removeDevices`.
  connectivity_assurance  body  FleetConnectivityAssurance (nullable)
  label  body  string  — New fleet name.
  removeDevices  body  array[string]  — DeviceUIDs to remove from the fleet. Each must be in the project and not also in `addDevices`.
  smart_rule_enabled  body  boolean  — Whether `smart_rule` adds and removes devices. Always send it; omitting it sets `false`.
  watchdog_mins  body  integer(int64)  — Minutes of inactivity before Notehub generates a watchdog event. Omit or `0` to keep the current value; a negative value turns it off.

get_fleet_environment_hierarchy  read  GET /v1/projects/{projectOrProductUID}/fleets/{fleetUID}/environment_hierarchy

get_fleet_environment_variables  read  GET /v1/projects/{projectOrProductUID}/fleets/{fleetUID}/environment_variables

set_fleet_environment_variables  write  PUT /v1/projects/{projectOrProductUID}/fleets/{fleetUID}/environment_variables
  Add or change a fleet's environment variables; variables you don't send keep their values.
  + environment_variable_notes
  *environment_variables  body  map[string]string

delete_fleet_environment_variable  destructive  DELETE /v1/projects/{projectOrProductUID}/fleets/{fleetUID}/environment_variables/{key}

set_global_event_transformation  destructive  POST /v1/projects/{projectOrProductUID}/global-transformation
  *body  body  string  — text/plain: JSONata expression as plain text.

disable_global_event_transformation  write  POST /v1/projects/{projectOrProductUID}/global-transformation/disable

enable_global_event_transformation  write  POST /v1/projects/{projectOrProductUID}/global-transformation/enable

get_project_members  read  GET /v1/projects/{projectOrProductUID}/members
  List project members

get_products  read  GET /v1/projects/{projectOrProductUID}/products
  List a project's products

create_product  write  POST /v1/projects/{projectOrProductUID}/products
  auto_provision_fleets  body  array[string]  — FleetUIDs that devices are added to when first provisioned with this ProductUID.
  disable_devices_by_default  body  boolean  — If `true`, devices provisioned with this ProductUID start out disabled.
  *label  body  string  — Product name shown in Notehub.
  *product_uid  body  string  — Suffix of the new ProductUID. Must start with a lowercase letter and contain only lowercase letters, digits, `.`, `_` or `-`.

delete_product  destructive  DELETE /v1/projects/{projectOrProductUID}/products/{productUID}

get_notefile_schemas  read  GET /v1/projects/{projectOrProductUID}/schemas
  List Notefile schemas

get_project_secrets  read  GET /v1/projects/{projectOrProductUID}/secrets
  List project secrets

create_project_secret  write  POST /v1/projects/{projectOrProductUID}/secrets
  *name  body  string  — 1 to 64 letters, digits or underscores.
  *value  body  string  — The secret value.

delete_project_secret  destructive  DELETE /v1/projects/{projectOrProductUID}/secrets/{secretName}

update_project_secret  write  PUT /v1/projects/{projectOrProductUID}/secrets/{secretName}
  *value  body  string  — The new secret value.

## alert

get_alerts  read  GET /v1/projects/{projectOrProductUID}/alerts
  List alerts
  + pageSize, pageNum, cursor
  monitorUID  query  string  — Return only this monitor's alerts.

## event

get_events  read  GET /v1/projects/{projectOrProductUID}/events
  List events in a project
  + pageSize, pageNum, cursor, deviceUID, sensorUID, sortBy, sortOrder, startDate, endDate, dateType, systemFilesOnly, files, format, fleetUID, selectFields
  serialNumber, sessionUID, eventUID  query  array[string]  — Repeat for more than one.

get_route_logs_by_event  read  GET /v1/projects/{projectOrProductUID}/events/{eventUID}/route-logs
  *eventUID  path  string  — EventUID, e.g. `4506f411-dea6-44a0-9743-1130f57d7747`.

get_fleet_events  read  GET /v1/projects/{projectOrProductUID}/fleets/{fleetUID}/events
  List events in a fleet
  + pageSize, pageNum, cursor, deviceUID, sensorUID, sortBy, sortOrder, startDate, endDate, dateType, systemFilesOnly, files, format, selectFields
  serialNumber, sessionUID, eventUID  query  array[string]  — Repeat for more than one.

## jobs

get_jobs  read  GET /v1/projects/{projectOrProductUID}/jobs
  List batch jobs

create_job  write  POST /v1/projects/{projectOrProductUID}/jobs
  Create a batch job
  *name  query  string  — Job name, shown in lists.
  comment  body  string  — Human-readable description of the job
  default_requests  body  BatchJobRequests  — Operations to apply to a device
  device_requests  body  map[string]BatchJobRequests  — Device-specific request overrides, keyed by DeviceUID
  report_options  body  object  — Controls what data is included in the job report
    app_fleets  boolean  — Include project fleets in the report
    app_info  boolean  — Include project info in the report
    app_vars  boolean  — Include project environment variables in the report
    comment  string
    device_activity  boolean  — Include device activity data in the report
    device_health  boolean  — Include device health data in the report
    device_info  boolean  — Include device info in the report
    device_vars  boolean  — Include device environment variables in the report
  select  body  object  — Device selection criteria
    all_devices  boolean  — Select all devices in the project
    comment  string
    devices  array[string]  — Specific DeviceUIDs to include
    devices_by_sn  array[string]  — Serial number patterns to match (supports glob wildcards *, ?, [...])
    devices_in_fleets  array[string]  — FleetUIDs whose devices should be included

delete_job  destructive  DELETE /v1/projects/{projectOrProductUID}/jobs/{jobUID}
  Delete a batch job

get_job  read  GET /v1/projects/{projectOrProductUID}/jobs/{jobUID}
  Get a batch job

run_job  write  POST /v1/projects/{projectOrProductUID}/jobs/{jobUID}/run
  Run a batch job
  dry_run  query  boolean  — If true, report what the job would change without changing anything.

get_job_runs  read  GET /v1/projects/{projectOrProductUID}/jobs/{jobUID}/runs
  List a job's runs
  status  query  string  — Return only runs with this status.
  dry_run  query  boolean (nullable)  — If set, return only dry runs (`true`) or only real runs (`false`).

delete_job_run  destructive  DELETE /v1/projects/{projectOrProductUID}/jobs/runs/{reportUID}

get_job_run  read  GET /v1/projects/{projectOrProductUID}/jobs/runs/{reportUID}
  view  query  string[summary|detail]  default=summary  — Controls the level of detail returned: 'summary' returns metadata only, 'detail' returns the full result payload

cancel_job_run  destructive  POST /v1/projects/{projectOrProductUID}/jobs/runs/{reportUID}/cancel

## monitor

get_monitors  read  GET /v1/projects/{projectOrProductUID}/monitors
  List monitors

create_monitor  write  POST /v1/projects/{projectOrProductUID}/monitors
  + aggregate_function, aggregate_window, alert, *alert_routes, condition_type, fleet_filter, *name, notefile_filter, per_device, routing_cooldown_period, silenced, source_selector, source_type, *threshold, uid, usage_scope, usage_type, usage_window
  description  body  string  — Free-form description.
  disabled  body  boolean  — If true, the monitor will not be evaluated.
  last_routed_at  body  string  — When the monitor last sent notifications.

delete_monitor  destructive  DELETE /v1/projects/{projectOrProductUID}/monitors/{monitorUID}

get_monitor  read  GET /v1/projects/{projectOrProductUID}/monitors/{monitorUID}

update_monitor  write  PUT /v1/projects/{projectOrProductUID}/monitors/{monitorUID}
  Omitted fields are cleared or reset, except `source_type`, `fleet_filter`, `notefile_filter`, `condition_type`…
  + aggregate_function, aggregate_window, alert, alert_routes, condition_type, fleet_filter, name, notefile_filter, per_device, routing_cooldown_period, silenced, source_selector, source_type, threshold, uid, usage_scope, usage_type, usage_window
  description  body  string  — Free-form description.
  disabled  body  boolean  — If true, the monitor will not be evaluated.
  last_routed_at  body  string  — When the monitor last sent notifications.

## route

get_routes  read  GET /v1/projects/{projectOrProductUID}/routes
  List routes

create_route  write  POST /v1/projects/{projectOrProductUID}/routes
  Omit `fleets` to send every fleet's events.
  + notes, type
  aws  body  AwsRoute
  azure  body  AzureRoute
  blynk  body  BlynkRoute
  datacake  body  DatacakeRoute
  disabled  body  boolean  — If true, the route sends no events.
  google  body  GoogleRoute
  http  body  HttpRoute
  label  body  string  — Route name. Required.
  mqtt  body  MqttRoute
  proxy  body  ProxyRoute
  qubitro  body  QubitroRoute
  radnote  body  RadRoute
  s3archive  body  S3ArchiveRoute
  slack  body  SlackRoute
  snowflake  body  SnowflakeRoute
  snowpipe_streaming  body  SnowpipeStreamingRoute
  thingworx  body  ThingworxRoute
  twilio  body  TwilioRoute

delete_route  destructive  DELETE /v1/projects/{projectOrProductUID}/routes/{routeUID}

get_route  read  GET /v1/projects/{projectOrProductUID}/routes/{routeUID}

update_route  write  PUT /v1/projects/{projectOrProductUID}/routes/{routeUID}
  Fields you leave out keep their values, and a secret sent back as `(omitted)` is kept, but leaving out `notes` clears…
  + notes, type
  aws  body  AwsRoute
  azure  body  AzureRoute
  blynk  body  BlynkRoute
  datacake  body  DatacakeRoute
  disabled  body  boolean  — If true, the route sends no events.
  google  body  GoogleRoute
  http  body  HttpRoute
  label  body  string  — Route name. Required.
  mqtt  body  MqttRoute
  proxy  body  ProxyRoute
  qubitro  body  QubitroRoute
  radnote  body  RadRoute
  s3archive  body  S3ArchiveRoute
  slack  body  SlackRoute
  snowflake  body  SnowflakeRoute
  snowpipe_streaming  body  SnowpipeStreamingRoute
  thingworx  body  ThingworxRoute
  twilio  body  TwilioRoute

get_route_logs_by_route  read  GET /v1/projects/{projectOrProductUID}/routes/{routeUID}/route-logs
  List a route's logs
  + pageSize, pageNum, cursor, deviceUID, startDate, endDate, systemFilesOnly, files
  sortBy  query  string[date|event]  default=date  — Field to sort by.
  sortOrder  query  string[asc|desc]  default=desc  — Sort direction.
  mostRecentOnly  query  boolean  — If true, return only the latest attempt for each event.
  routingStatus  query  array[string[success|failure]]  — Return only attempts with this outcome.
  responseStatus  query  array[string]  — Response status codes to include, e.g. `500`.

## Skills

get_skills  read  GET /v1/projects/{projectOrProductUID}/skills
  List project skills
  includeContent  query  boolean  — Use true to include exact Markdown content. Omit or use false for metadata only.
  If-None-Match  header  string  — Return 304 if an ETag from get_skills matches, or use * to match any current representation.

delete_skill  destructive  DELETE /v1/projects/{projectOrProductUID}/skills/{skillName}
  Delete a project skill
  *If-Match  header  string  — Required; use one strong ETag or quoted skill revision.

get_skill  read  GET /v1/projects/{projectOrProductUID}/skills/{skillName}
  Get a project skill
  If-Match  header  string  — Optionally require one strong ETag or quoted skill revision; wildcards and weak tags are rejected.
  If-None-Match  header  string  — Return 304 if an ETag for this skill matches, or use * to match any existing skill at this name.

put_skill  write  PUT /v1/projects/{projectOrProductUID}/skills/{skillName}
  Create or replace a project skill · Send `content` and optional `kinds`; omit `kinds` to keep them, or send `[]` to clear them.
  If-Match  header  string  — Required for replacement; use one strong ETag or quoted skill revision.
  If-None-Match  header  string  — Use * to create only if the name is absent; cannot be combined with If-Match.
  *content  body  string  — Exact UTF-8, at most 1 MiB decoded.
  kinds  body  array[string]  — Lowercase kinds; publish is reserved.

## usage

get_api_usage  read  GET /v1/projects/{projectOrProductUID}/usage/api
  Get billable API usage
  + startDate, endDate, limit, *period

get_data_usage  read  GET /v1/projects/{projectOrProductUID}/usage/data
  Get cellular and satellite data usage
  + startDate, endDate, deviceUID, fleetUID, limit, *period, aggregate

get_events_usage  read  GET /v1/projects/{projectOrProductUID}/usage/events
  + startDate, endDate, deviceUID, fleetUID, limit, *period, aggregate, skipRecentData
  notefile  query  array[string]  — Count only events from these notefiles.
  includeNotefiles  query  boolean  — If true, include event counts for each notefile.

get_route_logs_usage  read  GET /v1/projects/{projectOrProductUID}/usage/route-logs
  + startDate, endDate, limit, *period, skipRecentData
  routeUID  query  array[string]
  aggregate  query  string[route|project]  default=route  — How to group the counts, per route or one total for the project.

get_sessions_usage  read  GET /v1/projects/{projectOrProductUID}/usage/sessions
  + startDate, endDate, deviceUID, fleetUID, limit, *period, aggregate, skipRecentData

## types

DeviceUsage
  bytes_rcvd, bytes_rcvd_secondary, bytes_sent, bytes_sent_secondary, notes_rcvd, notes_sent, sessions_tcp, sessions_tls  integer(int64)
  duration  integer(int64)  — Duration in seconds
  since  integer(int64)  — Unix timestamp

TowerLocation
  c  string  — Country code
  cid  integer(int32)  — Cell ID
  l  string  — Open Location Code
  lac  integer(int32)  — Location Area Code
  lat  number(double)  — Latitude
  lon  number(double)  — Longitude
  mcc  integer(int32)  — Mobile Country Code
  mnc  integer(int32)  — Mobile Network Code
  n  string  — Name of the location
  source  string  — The source of this location
  time  integer(int64)  — Unix timestamp when this location was ascertained
  towers  integer(int32)  — Number of triangulation points
  z  integer(int32)  — Timezone ID
  zone  string  — Timezone name

FleetConnectivityAssurance
  enabled  boolean (nullable)  — Whether Connectivity Assurance is enabled for this fleet

BatchJobRequests
  comment  string
  connectivity_assurance_disable  boolean  — Disable connectivity assurance for the device
  connectivity_assurance_enable  boolean  — Enable connectivity assurance for the device
  disable  boolean  — Disable the device
  enable  boolean  — Enable the device
  fleets_to_default  array[string]  — FleetUIDs to assign to the device if it has no fleets
  fleets_to_join  array[string]  — FleetUIDs to add the device to
  fleets_to_leave  array[string]  — FleetUIDs to remove the device from
  note_reqs  array[BatchJobNoteRequest]  — note.add/note.update/note.delete requests to perform against the device's own notefiles.
  provision_product  string  — ProductUID to provision the device with if not already provisioned
  sn_to_default  string  — Set the device serial number only if not already set
  sn_to_set  string  — Set the device serial number ("-" to clear)
  vars_to_default  map[string]string  — Environment variables to set only if not already set
  vars_to_set  map[string]string  — Environment variables to set (use "-" as value to clear)

SlackWebHookNotification
  message_type  string[text|blocks]  — How `text` is read.
  text  string  — The message, or Block Kit JSON.
  *url  string  — The URL of the Slack webhook.

SlackBearerNotification
  channel  string  — The channel to send the message to.
  message_type  string[text|blocks]  — How `text` is read.
  text  string  — The message, or Block Kit JSON.
  *token  string  — The bearer token for the Slack app.

EmailNotification
  *email  string  — Email address, or several separated by commas.

TwilioNotification
  *account_sid  string  — The Twilio account SID.
  auth_token  string  — The Twilio auth token.
  from  string  — The number to send the SMS from, in E.164 format, e.g. `+15551234567`.
  message  string  — The text of the SMS.
  to  string  — The number to send the SMS to, in E.164 format.

AwsRoute
  access_key_id, access_key_secret, channel, message_deduplication_id, message_group_id, region  string
  disable_http_headers  boolean
  filter  Filter  — Filter applied to route data.
  fleets  array[string]
  http_headers  map[string]string
  role_arn  string  — IAM Role ARN for role-based authentication via STS AssumeRole
  throttle_ms, timeout  integer(int32)
  transform  RouteTransformSettings  — Settings for transforming route payloads before delivery.
  url  string(uri)

AzureRoute
  filter  Filter  — Filter applied to route data.
  fleets  array[string]
  functions_key_secret, sas_policy_key, sas_policy_name  string
  throttle_ms, timeout  integer(int32)
  transform  RouteTransformSettings  — Settings for transforming route payloads before delivery.
  url  string(uri)

BlynkRoute
  filter  Filter  — Filter applied to route data.
  fleets  array[string]
  region  string
  throttle_ms, timeout  integer(int32)
  transform  RouteTransformSettings  — Settings for transforming route payloads before delivery.

DatacakeRoute
  disable_http_headers  boolean
  filter  Filter  — Filter applied to route data.
  fleets  array[string]
  http_headers  map[string]string
  throttle_ms, timeout  integer(int32)
  transform  RouteTransformSettings  — Settings for transforming route payloads before delivery.
  url  string(uri)

GoogleRoute
  filter  Filter  — Filter applied to route data.
  fleets  array[string]
  throttle_ms, timeout  integer(int32)
  token  string
  transform  RouteTransformSettings  — Settings for transforming route payloads before delivery.
  url  string(uri)

HttpRoute
  disable_http_headers  boolean
  filter  Filter  — Filter applied to route data.
  fleets  array[string]  — If non-empty, applies only to the listed fleets.
  http_headers  map[string]string
  throttle_ms, timeout  integer(int32)
  transform  RouteTransformSettings  — Settings for transforming route payloads before delivery.
  url  string(uri)

MqttRoute
  broker, certificate, certificate_name, client_id, key, port, private_key_name, topic, username  string
  filter  Filter  — Filter applied to route data.
  fleets  array[string]
  password  string(password)
  throttle_ms, timeout  integer(int32)
  transform  RouteTransformSettings  — Settings for transforming route payloads before delivery.

ProxyRoute
  alias  string
  fleets  array[string]
  http_headers  map[string]string
  timeout  integer(int32)
  transform  RouteTransformSettings  — Settings for transforming route payloads before delivery.
  url  string(uri)

QubitroRoute
  filter  Filter  — Filter applied to route data.
  fleets  array[string]
  project_id, webhook_signing_key  string
  throttle_ms, timeout  integer(int32)
  transform  RouteTransformSettings  — Settings for transforming route payloads before delivery.

RadRoute
  client_id, data_feed_key  string
  client_secret  string(password)
  event_id, throttle_ms  integer(int32)
  fleets  array[string]
  test_api  boolean

S3ArchiveRoute
  archive_count_exceeds, archive_every_mins, throttle_ms, timeout  integer(int32)
  archive_id, bucket_endpoint, bucket_name, bucket_region, file_access, file_folder, file_format, key_id, key_secret  string
  filter  Filter  — Filter applied to route data.
  fleets  array[string]
  transform  RouteTransformSettings  — Settings for transforming route payloads before delivery.
  url  string(uri)

SlackRoute
  bearer, blocks, channel, text  string
  filter  Filter  — Filter applied to route data.
  fleets  array[string]
  throttle_ms, timeout  integer(int32)
  transform  RouteTransformSettings  — Settings for transforming route payloads before delivery.
  webhook_url  string(uri)

SnowflakeRoute
  account_name, organization_name, pem, private_key_name, user_name  string
  filter  Filter  — Filter applied to route data.
  fleets  array[string]
  timeout  integer(int32)
  transform  RouteTransformSettings  — Settings for transforming route payloads before delivery.

SnowpipeStreamingRoute
  account_name, database_name, organization_name, pem, pipe_name, private_key_name, schema_name, user_name  string
  filter  Filter  — Filter applied to route data.
  fleets  array[string]
  timeout  integer(int32)
  transform  RouteTransformSettings  — Settings for transforming route payloads before delivery.

ThingworxRoute
  app_key  string
  filter  Filter  — Filter applied to route data.
  fleets  array[string]
  throttle_ms, timeout  integer(int32)
  transform  RouteTransformSettings  — Settings for transforming route payloads before delivery.
  url  string(uri)

TwilioRoute
  account_sid, from, message, to  string
  auth_token  string(password)
  filter  Filter  — Filter applied to route data.
  fleets  array[string]
  throttle_ms, timeout  integer(int32)

BatchJobNoteRequest
  body  object  — The note's JSON body (used by note.add and note.update)
  *file  string  — The notefile to operate on (e.g. data.qi, config.dbs)
  note  string  — The note ID. Required for note.update and note.delete, and for note.add against a database (.dbs/.db) notefile.
  *req  string[note.add|note.update|note.delete]  — The note operation to perform

Filter
  files  array[string]  — List of notefile names or patterns to filter on.
  system_notefiles  boolean  — Whether system notefiles should be included.
  type  string  — Type of filter to apply (corresponds to `hublib.NotefileFilterType`).

RouteTransformSettings
  format  string  — Output format for transformed data (e.g., "json", "xml", "text").
  jsonata  string  — JSONata expression used to transform the data payload (outgoing).
  jsonata_in  string  — JSONata expression used to transform the data payload (incoming).
