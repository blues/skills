# MCP execute reference · Naming 1 · OpenAPI 3.0.3
# call_tool('<name>', {param: value, ...})  ·  *required  ·  header keys use the OpenAPI name exactly
# !write / !destructive on a param escalates the call when that param is used (a boolean: when true)
# each {x} in a path is a required param x, a string unless an `x  path` line defines it  ·  `a, b  query  string` = params with the same spec

## untagged

list_pets  read  GET /a

get_pet_by_id  read  GET /b

get_http_server_status  read  GET /c

o_auth2_token  read  GET /d

get_v2_users  read  GET /e

list_orders_v1  read  GET /f

fetch3_d_model  read  GET /g

s3_bucket_info  read  GET /h

io_stats  read  GET /i

search  read  GET /j
  query  query  string
  page-size  query  integer
