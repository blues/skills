# MCP execute reference · Unions 1 · OpenAPI 3.0.3
# call_tool('<name>', {param: value, ...})  ·  *required  ·  header keys use the OpenAPI name exactly
# !write / !destructive on a param escalates the call when that param is used (a boolean: when true)
# each {x} in a path is a required param x, a string unless an `x  path` line defines it  ·  `a, b  query  string` = params with the same spec
# a capitalized type such as Foo is defined under ## types; an inline object's fields are indented below it

## untagged

create_payment  write  POST /payments
  *body  body  object  — one of CardPayment | BankPayment; The payment to make.

put_value  write  PUT /values
  body  body  string|integer(int32)

## types

CardPayment
  card  string

BankPayment
  iban  string
