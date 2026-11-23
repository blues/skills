# MCP execute reference · Refs 1 · OpenAPI 3.0.3
# call_tool('<name>', {param: value, ...})  ·  *required  ·  header keys use the OpenAPI name exactly
# !write / !destructive on a param escalates the call when that param is used (a boolean: when true)
# each {x} in a path is a required param x, a string unless an `x  path` line defines it  ·  `a, b  query  string` = params with the same spec

## pets

add_pet  write  POST /stores/{storeId}/pets
  *storeId  path  integer(int64)  — Store that owns the pet.
  dryRun  query  boolean  default=true
  *name  body  string  — The pet's name.
  status  body  string[available|sold]  default=available
  ownerId  body  integer(int64)
