# MCP execute reference · External 1 · OpenAPI 3.0.3
# call_tool('<name>', {param: value, ...})  ·  *required  ·  header keys use the OpenAPI name exactly
# !write / !destructive on a param escalates the call when that param is used (a boolean: when true)
# each {x} in a path is a required param x, a string unless an `x  path` line defines it  ·  `a, b  query  string` = params with the same spec

## untagged

update_boat  write  PUT /boats/{boatId}
  *boatId  path  string  — Hull number, e.g. `HN-42`.
  *name  body  string
  length  body  number(float)  — Length in metres.
