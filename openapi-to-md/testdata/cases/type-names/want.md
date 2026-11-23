# MCP execute reference · Type names 1 · OpenAPI 3.0.3
# call_tool('<name>', {param: value, ...})  ·  *required  ·  header keys use the OpenAPI name exactly
# !write / !destructive on a param escalates the call when that param is used (a boolean: when true)
# each {x} in a path is a required param x, a string unless an `x  path` line defines it  ·  `a, b  query  string` = params with the same spec
# a capitalized type such as Foo is defined under ## types; an inline object's fields are indented below it

## tag: types

list_kinds  read  GET /types

## untagged

create_order_v1  write  POST /v1/orders
  item  body  Item

create_order_v2  write  POST /v2/orders
  item  body  Item2

create_route  write  POST /routes
  http  body  Http-route
  meta  body  Object
  home  body  Address
  first  body  PairAllOf0

## types

Item
  sku  string

Item2
  *productId  string(uuid)
  *quantity  integer

Http-route
  *url  string(uri)

Object
  *kind  string

Address
  line1  string

PairAllOf0
  a  string
