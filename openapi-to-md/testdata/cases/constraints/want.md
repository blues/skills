# MCP execute reference · Constraints 1 · OpenAPI 3.0.3
# call_tool('<name>', {param: value, ...})  ·  *required  ·  header keys use the OpenAPI name exactly
# !write / !destructive on a param escalates the call when that param is used (a boolean: when true)
# each {x} in a path is a required param x, a string unless an `x  path` line defines it  ·  `a, b  query  string` = params with the same spec

## untagged

update_fleet  write  PUT /fleets/{fleetId}
  Always send `smart_rule_enabled`; leaving it out turns the smart rule off.
  smart_rule_enabled  body  boolean  — Whether the smart rule runs. Always send it; omitting it sets `false`.
  watchdog_mins  body  integer  — Watchdog timeout. Omit or `0` to keep the current value; a negative value turns it off.
  threshold  body  number  — The value to compare against. For heartbeat monitors this is seconds of inactivity.
  label  body  string  — Fleet name. Required.
  level  body  string[low|high]  — Alert level, in seconds of delay.

start_update  write  POST /devices/dfu
  Start firmware updates · With no filters, that is every device in the project.
  tag  query  string
