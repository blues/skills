# MCP execute reference · Values 1 · OpenAPI 3.0.3
# call_tool('<name>', {param: value, ...})  ·  *required  ·  header keys use the OpenAPI name exactly
# !write / !destructive on a param escalates the call when that param is used (a boolean: when true)
# each {x} in a path is a required param x, a string unless an `x  path` line defines it  ·  `a, b  query  string` = params with the same spec

## untagged

create_job  write  POST /jobs
  interval  body  string  — How often the job runs. e.g. 10m|1h30m
  webhook  body  string  — URL to call when the job finishes.
  settings  body  object  — Display settings.
  order  body  string  — Sort order. Values: asc|desc
  notify  body  boolean  — Whether to notify.
  scope  body  array[string]  — What the job may touch: `devices`; `fleets`, `routes`
