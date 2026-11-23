# MCP execute reference · Recursive 1 · OpenAPI 3.0.3
# call_tool('<name>', {param: value, ...})  ·  *required  ·  header keys use the OpenAPI name exactly
# !write / !destructive on a param escalates the call when that param is used (a boolean: when true)
# each {x} in a path is a required param x, a string unless an `x  path` line defines it  ·  `a, b  query  string` = params with the same spec
# a capitalized type such as Foo is defined under ## types; an inline object's fields are indented below it

## untagged

put_tree  write  PUT /trees
  filter  query  Nested
  tree  body  Tree
  meta  body  JsonValue
  root  body  Node
  loop  body  any

## types

Nested = array[Nested]

Tree = array[Tree]

JsonValue = string|number|boolean|array[JsonValue]|map[string]JsonValue

Node
  name  string
  left  Node  — Left child.
  right  Node  — Right child.
