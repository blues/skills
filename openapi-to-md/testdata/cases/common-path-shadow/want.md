# MCP execute reference · Shadow 1 · OpenAPI 3.0.3
# call_tool('<name>', {param: value, ...})  ·  *required  ·  header keys use the OpenAPI name exactly
# !write / !destructive on a param escalates the call when that param is used (a boolean: when true)
# each {x} in a path is a required param x, a string unless an `x  path` line defines it  ·  `a, b  query  string` = params with the same spec
# `+ a, *b` = params defined under ## common params; a param line under a call applies to that call only

## common params
id  path  integer(int64)  — Numeric user ID, as returned when the user was created.

## untagged

get_user  read  GET /users/{id}

delete_user  destructive  DELETE /users/{id}

get_user_roles  read  GET /users/{id}/roles

get_user_keys  read  GET /users/{id}/keys

get_user_logs  read  GET /users/{id}/logs

get_tag  read  GET /tags/{id}
  *id  path  string

get_item  read  GET /items/{item_id}
  *item_id  path  string
