# MCP execute reference · Precedence 1 · OpenAPI 3.0.3
# call_tool('<name>', {param: value, ...})  ·  *required  ·  header keys use the OpenAPI name exactly
# !write / !destructive on a param escalates the call when that param is used (a boolean: when true)
# each {x} in a path is a required param x, a string unless an `x  path` line defines it  ·  `a, b  query  string` = params with the same spec

## untagged

get_task  read  GET /tasks/{taskId}
  purge  query  boolean  !destructive  — If true, delete the task after reading it.
  v  query  boolean

run_task  write  POST /tasks/{taskId}
  mode  query  string  !destructive
  body_mode  body  string  — Mode in the body.

delete_task  destructive  DELETE /tasks/{taskId}

list_all_tasks  read  GET /tasks
