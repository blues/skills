package main

import (
	"fmt"
	"strings"
	"testing"
)

func TestSnakeCase(t *testing.T) {
	for id, want := range map[string]string{
		"GetDevicesDfuStatus":     "get_devices_dfu_status",
		"GetDbNote":               "get_db_note",
		"GetAWSRoleConfig":        "get_aws_role_config",
		"OAuth2ClientCredentials": "o_auth2_client_credentials",
		"listPets":                "list_pets",
	} {
		if got := snakeCase(id); got != want {
			t.Errorf("snakeCase(%q) = %q, want %q", id, got, want)
		}
	}
}

// A small OpenAPI 3.0 API that exercises each rule; nothing in it resembles any real API.
const petstore = `
openapi: 3.0.3
info: {title: Petstore, version: 2.1.0}
tags:
  - name: pets
  - name: owners
paths:
  x-settings:                     # an extension, not a path
    get: {operationId: NotAnEndpoint, summary: Not an endpoint}
  /pets/{petId}:
    parameters:
      - {name: petId, in: path, schema: {type: string}}
      - {name: verbose, in: query, schema: {type: boolean}}
    get:
      operationId: GetPet
      tags: [pets]
      summary: Get a pet
      parameters:
        - {name: verbose, in: query, required: true, schema: {type: boolean}}
        - {name: legacy, in: query, deprecated: true, schema: {type: string}}
        - {name: Authorization, in: header, required: true, schema: {type: string}}
        - $ref: '#/components/parameters/Fields'
      responses:
        '200':
          description: OK
          content:
            application/json:
              examples:
                sample: {value: {$ref: './this-is-data.yaml'}}
    delete:
      operationId: RemovePet
      tags: [pets]
      summary: Remove a pet
      x-effect: write
      x-effect-note: It only archives the pet.
    put:
      operationId: ReplacePet
      tags: [pets]
      description: |
        Replace a pet entirely.
        More detail on a second line.
      x-effect: destructive
      requestBody:
        required: true
        content:
          application/json:
            schema:
              required: [owner]           # required, with no property schema
              allOf:
                - $ref: '#/components/schemas/PetFields'
                - type: object
                  required: [name]
                  properties:
                    name: {type: string}
                    id: {type: string, readOnly: true}
                    # In 3.0 the keywords beside a $ref are ignored, so this stays writable.
                    nickname: {$ref: '#/components/schemas/Name', readOnly: true}
                - type: object
                  properties:
                    id: {type: string, maxLength: 8}   # still read-only from the branch above
  /pets/{petId}/photo:
    post:
      operationId: UploadPhoto
      tags: [pets]
      summary: Upload a photo
      deprecated: true
      requestBody:                      # not required: optional
        content:
          image/png: {schema: {type: string, format: binary}}
          image/jpeg: {schema: {type: string, format: binary}}
  /pets/{petId}/notes:
    patch:
      operationId: AddNote
      tags: [pets]
      summary: Add a note
      requestBody:
        required: true
        content:
          text/plain: {schema: {type: string}}
          application/json:
            schema: {type: object, required: [text], properties: {text: {type: string}}}
  /pets/search:
    post:
      operationId: SearchPets
      tags: [pets]
      summary: Search pets
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              properties: {limit: {type: integer}}
              oneOf:
                - {required: [name], properties: {name: {type: string}}}
                - {required: [species], properties: {species: {type: string}}}
  /owners:
    get:
      operationId: ListOwners
      tags: [owners, billing]
      summary: List owners
      parameters:
        - name: sort
          in: query
          schema: {type: array, items: {type: string, enum: [name, age]}}
        - name: filter
          in: query
          schema: {enum: [null, 'null', 'a|b']}
  /health:
    get:
      operationId: GetHealth
      summary: Health check
  /stats:
    get:
      operationId: GetStats
      tags: [admin]
      summary: Statistics
components:
  parameters:
    Fields:
      name: fields
      in: query
      schema: {type: string, enum: [all, basic]}
  schemas:
    Name: {type: string}
    PetFields:
      type: object
      required: [species]
      properties:
        species: {type: string, enum: [cat, dog]}
        tags: {type: array, items: {type: string}}
`

func renderDoc(t *testing.T, doc string, overrides map[string]override) (string, []string) {
	t.Helper()
	spec, _, err := decodeSpec([]byte(doc))
	if err != nil {
		t.Fatal(err)
	}
	out, warnings, err := render(spec, overrides, "https://example.com/openapi.yaml", "")
	if err != nil {
		t.Fatal(err)
	}
	return out, warnings
}

// Like expect, but ignoring how the text is wrapped, for the header's wrapped paragraph.
func expectText(t *testing.T, out string, wants ...string) {
	t.Helper()
	flat := strings.Join(strings.Fields(out), " ")
	for _, want := range wants {
		if !strings.Contains(flat, want) {
			t.Errorf("output lacks %q\n--- output ---\n%s", want, out)
		}
	}
}

func expect(t *testing.T, out string, wants ...string) {
	t.Helper()
	for _, want := range wants {
		if !strings.Contains(out, want) {
			t.Errorf("output lacks %q\n--- output ---\n%s", want, out)
		}
	}
}

func TestRender(t *testing.T) {
	out, warnings := renderDoc(t, petstore, nil)
	if len(warnings) != 0 {
		t.Errorf("unexpected warnings: %v", warnings)
	}
	expect(t, out,
		"# Petstore operations\n",

		// Path-level parameters, replaced by operation-level ones; deprecated ones and the
		// Authorization header left out; path parameters always required; a referenced
		// parameter and its allowed values.
		"- `GetPet` `get_pet` read `GET /pets/{petId}`\n  Get a pet\n  path: petId*; query: verbose*, fields: all|basic\n",

		// x-effect and x-effect-note; an operation with no parameters of its own still
		// lists the path-level ones.
		"- `RemovePet` `remove_pet` write `DELETE /pets/{petId}`\n  Remove a pet. Note: It only archives the pet.\n  path: petId*; query: verbose\n",

		// Summary from the description's first line; a required JSON body's properties
		// from allOf, required ones marked (including one with no property schema),
		// read-only ones left out even when another allOf branch also describes them.
		"- `ReplacePet` `replace_pet` destructive `PUT /pets/{petId}`\n  Replace a pet entirely.\n  path: petId*; query: verbose; body: name*, nickname, owner*, species*: cat|dog, tags[]\n",

		// A body that isn't required is optional; a non-JSON body gives its media types;
		// deprecated operations are marked.
		"- `UploadPhoto` `upload_photo` write `POST /pets/{petId}/photo` (deprecated)\n  Upload a photo\n  body (optional): (image/jpeg or image/png)\n",

		// A body with several media types says whose fields are listed.
		"  body (fields of application/json; also text/plain): text*\n",

		// Alternative schemas are counted, not flattened.
		"  body: limit, plus the fields of one of 2 alternative schemas\n",

		// Array parameters with allowed item values; values that aren't plain words as JSON;
		// tags beyond the first.
		"- `ListOwners` `list_owners` read `GET /owners` (also tagged billing)\n  List owners\n  query: sort[]: name|age, filter: null|\"null\"|\"a|b\"\n",
	)
	expectText(t, out, "The 9 operations in the paths of [the OpenAPI specification](https://example.com/openapi.yaml) (Petstore 2.1.0), grouped by tag. This is a complete index")
	if strings.Contains(out, "NotAnEndpoint") {
		t.Errorf("a paths extension was listed as an operation:\n%s", out)
	}

	// Tags in the specification's order, then unlisted tags, then untagged operations.
	pets, owners, admin, untagged := strings.Index(out, "## pets"), strings.Index(out, "## owners"),
		strings.Index(out, "## admin"), strings.Index(out, "## Untagged")
	if !(pets < owners && owners < admin && admin < untagged) || pets < 0 {
		t.Errorf("tag sections out of order: pets %d, owners %d, admin %d, untagged %d", pets, owners, admin, untagged)
	}
}

// More OpenAPI 3.0 cases: path item references, allOf enums, media types, alternatives,
// type lists, YAML dates, long reference chains, and references the reference can't follow.
const more = `
openapi: 3.0.3
info: {title: More, version: '1'}
x-shared:
  pets:
    get: {operationId: ListPets, summary: List pets}
paths:
  /pets:
    $ref: '#/x-shared/pets'
    post:                             # beside the path item's $ref: listed too
      operationId: AddPet
      summary: Add a pet
      requestBody:
        required: true
        content:
          Application/JSON:           # media types are case-insensitive
            schema:
              allOf:
                - properties: {color: {enum: [red, green]}}
                - properties: {color: {enum: [green, blue]}}   # only green satisfies both
              anyOf:
                - {required: [name], properties: {name: {type: string}}}
                - {required: [tag], properties: {tag: {type: string}}}
  /pets/stream:
    post:
      operationId: StreamPets
      summary: Stream pets
      requestBody:
        required: true
        content:
          application/json-seq: {schema: {type: object, properties: {id: {type: string}}}}
  /pets/{petId}:
    put:
      operationId: ReplacePet
      summary: Replace a pet
      parameters:
        - {name: petId, in: path, schema: {type: string}}
        - name: born
          in: query
          schema: {type: string, enum: [2024-01-01]}      # a string, not a timestamp
        - name: kind
          in: query
          schema: {type: [array, string], items: {enum: [a, b]}}   # a string is valid too
      requestBody:
        required: true
        content:
          application/json:
            schema:
              allOf:
                - oneOf: [{properties: {a: {type: string}}}, {properties: {b: {type: string}}}]
                - oneOf: [{properties: {c: {type: string}}}, {properties: {d: {type: string}}}]
  /chain:
    get:
      operationId: GetChain
      summary: Long chain
      parameters:
        - {name: deep, in: query, schema: {$ref: '#/components/schemas/L0'}}
      responses:
        default: {$ref: './responses.yaml'}       # under responses, default is a response
        '200':
          description: OK
          content:
            application/json:
              schema:
                properties:
                  example: {$ref: '#/components/schemas/Gone'}   # a property named example
                  space: {$ref: '#/components/schemas/a%20b'}    # percent-encoded pointer
                  anchored: {$ref: '#node'}
                  dynamic: {$dynamicRef: '#node'}
                  loop: {$ref: '#/components/schemas/LoopA'}
components:
  schemas:
    'a b': {type: string}
    LoopA: {$ref: '#/components/schemas/LoopB'}
    LoopB: {$ref: '#/components/schemas/LoopA'}
`

// A chain of 40 references, longer than any fixed depth limit the generator used to have.
func longChain() string {
	var b strings.Builder
	for i := 0; i < 40; i++ {
		fmt.Fprintf(&b, "    L%d: {$ref: '#/components/schemas/L%d'}\n", i, i+1)
	}
	b.WriteString("    L40: {type: string, enum: [end]}\n")
	return b.String()
}

func TestMore(t *testing.T) {
	out, warnings := renderDoc(t, more+longChain(), nil)
	expect(t, out,
		"The 5 operations in the paths of",
		"- `ListPets` `list_pets` read `GET /pets`\n",
		"- `AddPet` `add_pet` write `POST /pets`\n  Add a pet\n  body: color: green, plus the fields of one or more of 2 alternative schemas\n",
		"- `StreamPets` `stream_pets` write `POST /pets/stream`\n  Stream pets\n  body: (application/json-seq)\n",
		// A date-shaped value stays the text it is; a [array, string] type isn't an array;
		// alternatives in separate allOf parts are not one choice.
		"  path: petId*; query: born: 2024-01-01, kind; body: (2 sets of alternative schemas; see the specification)\n",
		"  query: deep: end\n",
	)
	joined := strings.Join(warnings, "\n")
	for _, want := range []string{
		`"./responses.yaml" points outside the document`,
		`"#/components/schemas/Gone" points to nothing`,
		`"#node" names an anchor`,
		`dynamic reference "#node" is not followed`,
		"is part of a cycle of references",
	} {
		if !strings.Contains(joined, want) {
			t.Errorf("warnings lack %q: %v", want, warnings)
		}
	}
	if strings.Contains(joined, "a%20b") {
		t.Errorf("a percent-encoded pointer was not followed: %v", warnings)
	}
}

// OpenAPI 3.1 and 3.2 features, and what the reference reports instead of listing.
const later = `
openapi: 3.2.0
info: {title: Later, version: '1'}
webhooks:
  newPet: {post: {operationId: NewPetHook}}
paths:
  /pets:
    query:
      operationId: QueryPets
      summary: Query pets
      parameters:
        - name: filter
          in: querystring
          required: true
          content: {application/json: {$ref: '#/components/mediaTypes/Filter'}}
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              properties:
                # From 3.1, the keywords beside a $ref apply too.
                id: {$ref: '#/components/schemas/Id', readOnly: true}
                tags: {type: [array, 'null'], items: {enum: [a, b]}}
                pairs: {type: array, items: {type: string, enum: [x, y]}, enum: [[x], [y]]}
                owner: {$ref: './owner.yaml'}
                breed: {$ref: '#/components/schemas/Missing'}
    additionalOperations:
      LINK: {operationId: LinkPets}
    post:
      operationId: AddPet
      summary: Add a pet
      callbacks:
        added: {'{$request.body#/url}': {post: {operationId: AddedHook}}}
      requestBody:
        required: true
        content:
          application/json: {$ref: '#/components/mediaTypes/NewPet'}
  /owners:
    post:
      operationId: AddOwner
      summary: Add an owner
      requestBody:
        required: true
        content:
          application/json:
            schema:
              # The referenced schema's required id and the sibling's required name both apply.
              $ref: '#/components/schemas/Base'
              required: [name]
              properties: {name: {type: string}}
components:
  schemas:
    Id: {type: string}
    Base: {type: object, required: [id], properties: {id: {type: string}}}
  mediaTypes:
    NewPet: {schema: {type: object, required: [species], properties: {species: {type: string}}}}
    Filter: {schema: {type: string, enum: [all, mine]}}
`

func TestLaterVersions(t *testing.T) {
	out, warnings := renderDoc(t, later, nil)
	expect(t, out,
		"- `QueryPets` `query_pets` read `QUERY /pets`\n  Query pets\n  querystring: filter*: all|mine; body: breed, owner, pairs[]: [\"x\"]|[\"y\"], tags[]: a|b\n",
		"- `AddPet` `add_pet` write `POST /pets`\n  Add a pet\n  body: species*\n",
		"- `AddOwner` `add_owner` write `POST /owners`\n  Add an owner\n  body: id*, name*\n",
	)
	// The published header, not only the warnings, says what isn't listed.
	expectText(t, out, "The 3 operations in the paths of",
		"grouped by tag. Not listed: 1 webhook, 1 callback and 1 additionalOperations entry; read them in the specification. Apart from those, this is a complete index")
	joined := strings.Join(warnings, "\n")
	for _, want := range []string{
		"webhooks (newPet) are not listed",
		"/pets has additionalOperations (LINK)",
		"AddPet has callbacks (added)",
		`"./owner.yaml" points outside the document`,
		`"#/components/schemas/Missing" points to nothing`,
	} {
		if !strings.Contains(joined, want) {
			t.Errorf("warnings lack %q: %v", want, warnings)
		}
	}
}

func TestNonASCIIOperationID(t *testing.T) {
	_, warnings := renderDoc(t, "openapi: 3.0.3\ninfo: {title: X, version: '1'}\npaths:\n  /item:\n    get: {operationId: Getİtem}\n", nil)
	if !strings.Contains(strings.Join(warnings, "\n"), "Getİtem is not plain ASCII") {
		t.Errorf("no warning for a non-ASCII operationId: %v", warnings)
	}
}

func TestOverrides(t *testing.T) {
	out, warnings := renderDoc(t, petstore, map[string]override{
		"RemovePet":  {Effect: "destructive", Note: "Deletes the pet."}, // beats the spec's x-effect
		"ListOwners": {Effect: "sometimes"},                             // not an effect: ignored
		"FeedPet":    {Effect: "write"},                                 // not an operation
	})
	expect(t, out,
		"- `RemovePet` `remove_pet` destructive `DELETE /pets/{petId}`\n  Remove a pet. Note: Deletes the pet.\n",
		"- `ListOwners` `list_owners` read `GET /owners`",
	)
	joined := strings.Join(warnings, "\n")
	for _, want := range []string{`"sometimes"`, "FeedPet"} {
		if !strings.Contains(joined, want) {
			t.Errorf("warnings lack %q: %v", want, warnings)
		}
	}
}

func TestRejectsNonOpenAPI3(t *testing.T) {
	if _, _, err := render(map[string]any{"swagger": "2.0"}, nil, "", ""); err == nil {
		t.Error("expected an error for a Swagger 2.0 document")
	}
}
