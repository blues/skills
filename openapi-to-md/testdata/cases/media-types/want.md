# MCP execute reference · Media 1 · OpenAPI 3.0.3
# call_tool('<name>', {param: value, ...})  ·  *required  ·  header keys use the OpenAPI name exactly
# !write / !destructive on a param escalates the call when that param is used (a boolean: when true)
# each {x} in a path is a required param x, a string unless an `x  path` line defines it  ·  `a, b  query  string` = params with the same spec

## untagged

put_blob  write  PUT /blobs/{name}
  *body  body  string(binary)  — application/octet-stream: The bytes to store.

append_blob  write  POST /blobs/{name}
  body  body  string  — text/plain

submit_form  write  POST /forms
  data  body  object

post_event  write  POST /events
  *body  body  object  — Any JSON object.

put_events  write  PUT /events
  body  body  array[map[string]string]
