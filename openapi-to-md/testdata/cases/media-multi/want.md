# MCP execute reference · Media 1 · OpenAPI 3.1.0
# call_tool('<name>', {param: value, ...})  ·  *required  ·  header keys use the OpenAPI name exactly
# !write / !destructive on a param escalates the call when that param is used (a boolean: when true)
# each {x} in a path is a required param x, a string unless an `x  path` line defines it  ·  `a, b  query  string` = params with the same spec

## untagged

upload_image  write  POST /images
  *body  body  string(binary)  — image/jpeg|image/png

upload_file  write  POST /files
  *file  body  string(binary)
  thumb  body  string(base64)
  caption  body  string

get_customer  read  GET /customers/{id}
