# MCP execute reference · Unions 1 · OpenAPI 3.0.3
# call_tool('<name>', {param: value, ...})  ·  *required  ·  header keys use the OpenAPI name exactly
# !write / !destructive on a param escalates the call when that param is used (a boolean: when true)
# each {x} in a path is a required param x, a string unless an `x  path` line defines it  ·  `a, b  query  string` = params with the same spec
# a capitalized type such as Foo is defined under ## types; an inline object's fields are indented below it

## untagged

pay  write  POST /pay
  *body  body  object  — one of Card | BankAccount

create_alert  write  POST /alerts
  *target  body  Target1|Target2
  options  body  object (nullable)
    retries  integer
  items  body  Item|array[Item]

create_contact  write  POST /contacts
  body requires one of: email | phone
  *name, email, phone  body  string
  status  body  string[new|active|done]

add_pet  write  POST /pets
  body plus the fields of one of: Dog | Cat; petType: Dog=Dog | Cat=Cat
  *name  body  string

create_payment  write  POST /payments
  body  body  object  — one of CardPayment | BankPayment; method: bank_transfer=BankPayment | card=CardPayment

create_shape  write  POST /shapes
  shape  body  Shape

## types

Card
  *card  string

BankAccount
  *iban  string

Target1
  *email  string(email)

Target2
  *url  string(uri)

Item
  *sku  string

Dog
  *petType, *bark  string

Cat
  *petType  string
  indoor  boolean

CardPayment
  *method, *number  string

BankPayment
  *method, *iban  string

Shape  — kind: circle=Circle | rect=Rect
  *kind  string

Circle
  *kind  string
  radius  number

Rect
  *kind  string
  w, h  number
