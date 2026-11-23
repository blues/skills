# MCP execute reference · Classes 1 · OpenAPI 3.0.3
# call_tool('<name>', {param: value, ...})  ·  *required  ·  header keys use the OpenAPI name exactly
# !write / !destructive on a param escalates the call when that param is used (a boolean: when true)
# each {x} in a path is a required param x, a string unless an `x  path` line defines it  ·  `a, b  query  string` = params with the same spec

## untagged

list_files  read  GET /files
  purge  query  boolean  !destructive  — If true, delete the files this call returns.
  touch  query  string(date-time)  !write  — Set the access time of the returned files to this time.
  depth  query  integer

search_files  read  POST /files
  query  body  string
  purgeMatches  body  boolean  !destructive  — If true, delete every match.

check_files  read  HEAD /files

file_options  read  OPTIONS /files

replace_file  destructive  PUT /files/{fileId}

patch_file  write  PATCH /files/{fileId}

delete_file  destructive  DELETE /files/{fileId}
  reason  query  string

trace_file  write  TRACE /files/{fileId}
