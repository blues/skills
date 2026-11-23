# MCP execute reference · Ref descriptions 1 · OpenAPI 3.1.0
# call_tool('<name>', {param: value, ...})  ·  *required  ·  header keys use the OpenAPI name exactly
# !write / !destructive on a param escalates the call when that param is used (a boolean: when true)
# each {x} in a path is a required param x, a string unless an `x  path` line defines it  ·  `a, b  query  string` = params with the same spec

## untagged

post_a  write  POST /a
  limit  query  integer  — How many items to return.
  body  body  string(binary)  — application/octet-stream: Generic bytes.

post_b  write  POST /b
  limit  query  integer  — At most 10 for this call.
  body  body  string(binary)  — application/octet-stream: A gzip-compressed firmware image.

get_c  read  GET /c
  limit  query  integer  — How many items to return.
