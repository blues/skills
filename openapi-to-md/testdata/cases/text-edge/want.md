# MCP execute reference · Text edges 1 · OpenAPI 3.0.3
# call_tool('<name>', {param: value, ...})  ·  *required  ·  header keys use the OpenAPI name exactly
# !write / !destructive on a param escalates the call when that param is used (a boolean: when true)
# each {x} in a path is a required param x, a string unless an `x  path` line defines it  ·  `a, b  query  string` = params with the same spec

## untagged

list_items  read  GET /items
  limit  query  integer  — # of items to return per page.
  host  query  string  — `[a-z](?:x)` must match; see the guide.
  badge  query  string  — status Badge style.
  greeting  query  string  default="Hello\nWorld"
  sep  query  string["a\tb"|c]
  nothing  query  null
  notes  query  string  — Free text, up to 1 KB.
