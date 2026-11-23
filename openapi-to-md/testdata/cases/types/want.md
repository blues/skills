# MCP execute reference · Types 1 · OpenAPI 3.0.3
# call_tool('<name>', {param: value, ...})  ·  *required  ·  header keys use the OpenAPI name exactly
# !write / !destructive on a param escalates the call when that param is used (a boolean: when true)
# each {x} in a path is a required param x, a string unless an `x  path` line defines it  ·  `a, b  query  string` = params with the same spec
# a capitalized type such as Foo is defined under ## types; an inline object's fields are indented below it

## untagged

plant_tree  write  POST /trees
  *species  body  string
  *location  body  Location
  root  body  Node
  waypoints  body  array[Location]
  labels  body  map[string]Label
  care  body  object  — Care plan.
    *interval  integer  — Days between waterings.
    fertilize  boolean
    log  object
      notes  string
  notify  body  array[EmailTarget|WebhookTarget]

replant_tree  write  PUT /trees
  location  body  Location

## types

Location
  *lat, *lon  number(double)
  clear  boolean  !destructive  — If true, clear the location's history.

Node
  name  string
  children  array[Node]

Label
  text  string
  color  string[red|green]

EmailTarget
  *email  string(email)

WebhookTarget
  *url  string(uri)
  secret  string
