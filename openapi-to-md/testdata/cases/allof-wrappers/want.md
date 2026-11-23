# MCP execute reference · Wrappers 1 · OpenAPI 3.0.3
# call_tool('<name>', {param: value, ...})  ·  *required  ·  header keys use the OpenAPI name exactly
# !write / !destructive on a param escalates the call when that param is used (a boolean: when true)
# each {x} in a path is a required param x, a string unless an `x  path` line defines it  ·  `a, b  query  string` = params with the same spec
# a capitalized type such as Foo is defined under ## types; an inline object's fields are indented below it

## untagged

create_order  write  POST /orders
  sort  query  string[asc|desc]  default=asc
  limit  query  integer(int32)
  *status  body  string[pending|paid|shipped]  — Initial status.
  priority  body  integer[1|2|3] (nullable)
  shipTo  body  Address  — Where to ship.
  tags  body  array[string[pending|paid|shipped]]
  labels  body  map[string]string
  born  body  string(date-time)  — When it was born.
  narrowed  body  string[paid|shipped]

create_circle  write  POST /circles
  *kind  body  string[circle]  — The shape kind.
  radius  body  number

## types

Address
  *line1, city  string
