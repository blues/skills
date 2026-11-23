# MCP execute reference · Deprecated 1 · OpenAPI 3.0.3
# call_tool('<name>', {param: value, ...})  ·  *required  ·  header keys use the OpenAPI name exactly
# !write / !destructive on a param escalates the call when that param is used (a boolean: when true)
# each {x} in a path is a required param x, a string unless an `x  path` line defines it  ·  `a, b  query  string` = params with the same spec

## untagged

get_old  read  GET /old  deprecated

get_new  read  GET /new
  legacyFilter  query  string  deprecated
  filter  query  string
