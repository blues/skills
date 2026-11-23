# MCP execute reference · Common 1 · OpenAPI 3.0.3
# call_tool('<name>', {param: value, ...})  ·  *required  ·  header keys use the OpenAPI name exactly
# !write / !destructive on a param escalates the call when that param is used (a boolean: when true)
# each {x} in a path is a required param x, a string unless an `x  path` line defines it  ·  `a, b  query  string` = params with the same spec
# `+ a, *b` = params defined under ## common params; a param line under a call applies to that call only

## common params
shopId  path  string  — The shop's handle, e.g. `acme-store`, as shown in its admin URL.
limit  query  integer(int32)  default=100  — Most items to return per page, between 1 and 1000.
cursor  query  string  — Opaque paging cursor from the previous response's `next` field, passed back unchanged.
since  query  integer(int64)  — Only items changed at or after this Unix time, in seconds.

## orders

list_orders  read  GET /shops/{shopId}/orders
  + limit, cursor, since
  status  query  string[open|paid|shipped]  — Only orders in this state.

list_refunds  read  GET /shops/{shopId}/refunds
  + limit, cursor
  *since  query  integer(int64)  — Only refunds issued at or after this Unix time, in seconds.

## customers

list_customers  read  GET /shops/{shopId}/customers
  + cursor
  limit  query  integer  default=10  — Customers per page, at most 50.
  purge  query  boolean  !destructive  — If true, delete the customers this call returns.

export_customers  write  POST /shops/{shopId}/customers/export
  + since
  purge  query  boolean  !destructive  — If true, delete the customers this call returns.
  page_cursor  query  string  — Opaque paging cursor from the previous response's `next` field, passed back unchanged.
