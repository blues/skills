# MCP execute reference · Tags 1 · OpenAPI 3.0.3
# call_tool('<name>', {param: value, ...})  ·  *required  ·  header keys use the OpenAPI name exactly
# !write / !destructive on a param escalates the call when that param is used (a boolean: when true)
# each {x} in a path is a required param x, a string unless an `x  path` line defines it  ·  `a, b  query  string` = params with the same spec

## alpha

get_one  read  GET /one

get_four  read  GET /four

## untagged

get_two  read  GET /two

get_five  read  GET /five

## zeta

get_three  read  GET /three
