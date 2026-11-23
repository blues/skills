# MCP execute reference · Descriptions 1 · OpenAPI 3.0.3
# call_tool('<name>', {param: value, ...})  ·  *required  ·  header keys use the OpenAPI name exactly
# !write / !destructive on a param escalates the call when that param is used (a boolean: when true)
# each {x} in a path is a required param x, a string unless an `x  path` line defines it  ·  `a, b  query  string` = params with the same spec

## untagged

list_reports  read  GET /reports
  reportType, ownerUID  query  string
  format  query  string  — Output format, e.g. `csv` or `json.gz`.
  owner  query  string  — Matches the owner's email address, case-insensitive.
  window  query  string  — `day`: the last 24 hours; `week`: the last 7 days
  scope  query  string  — What to include. Values: mine|team|all
  sort  query  string[asc|desc]  — Sort order.
  cursor  query  string  — An opaque cursor returned by `list_reports` or by get_report_page, which must be passed back unchanged, byte for byte…

get_report_page  read  GET /reports/{reportId}
  Get one page of a report.

delete_report  destructive  DELETE /reports/{reportId}
  *reportId  path  string  — Report ID, as listed by `list_reports`.
