# MCP execute reference · Hidden 1 · OpenAPI 3.0.3
# call_tool('<name>', {param: value, ...})  ·  *required  ·  header keys use the OpenAPI name exactly
# !write / !destructive on a param escalates the call when that param is used (a boolean: when true)
# each {x} in a path is a required param x, a string unless an `x  path` line defines it  ·  `a, b  query  string` = params with the same spec

## auth

sign_in  write  POST /login  hidden
  Exchange a password for a session token
  *user  body  string
  *password  body  string(password)

who_am_i  read  GET /me
  Show the signed-in user
