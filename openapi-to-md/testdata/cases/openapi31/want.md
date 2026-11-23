# MCP execute reference · ThreeOne 1 · OpenAPI 3.1.0
# call_tool('<name>', {param: value, ...})  ·  *required  ·  header keys use the OpenAPI name exactly
# !write / !destructive on a param escalates the call when that param is used (a boolean: when true)
# each {x} in a path is a required param x, a string unless an `x  path` line defines it  ·  `a, b  query  string` = params with the same spec

## untagged

patch_item  write  PATCH /items/{itemId}
  *itemId  path  string|integer
  note  query  string (nullable)
  *kind  body  string[gadget]
  weight  body  number(float) (nullable)
  labels  body  array[string] (nullable)
  mode  body  string[fast|slow] (nullable)
