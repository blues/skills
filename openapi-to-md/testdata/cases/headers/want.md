# MCP execute reference · Headers 1 · OpenAPI 3.0.3
# call_tool('<name>', {param: value, ...})  ·  *required  ·  header keys use the OpenAPI name exactly
# !write / !destructive on a param escalates the call when that param is used (a boolean: when true)
# each {x} in a path is a required param x, a string unless an `x  path` line defines it  ·  `a, b  query  string` = params with the same spec

## untagged

put_doc  write  PUT /docs/{docName}
  *If-Match  header  string  — Current ETag of the doc.
  X-Request-ID  header  string(uuid)
  session  cookie  string
  *body  body  string  — text/markdown: The document text.
