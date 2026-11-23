# MCP execute reference · Merge 1 · OpenAPI 3.0.3
# call_tool('<name>', {param: value, ...})  ·  *required  ·  header keys use the OpenAPI name exactly
# !write / !destructive on a param escalates the call when that param is used (a boolean: when true)
# each {x} in a path is a required param x, a string unless an `x  path` line defines it  ·  `a, b  query  string` = params with the same spec

## untagged

get_team  read  GET /orgs/{orgId}/teams/{teamId}
  verbose  query  string[low|high]  — Operation-level description wins.
  limit  query  integer
  expand  query  array[string]

delete_team  destructive  DELETE /orgs/{orgId}/teams/{teamId}
  verbose  query  boolean  — Path-level description.
  limit  query  integer
