package main

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"gopkg.in/yaml.v3"
)

// operation is one callable tool.
type operation struct {
	name, method, path, tag string
	class                   class
	hidden, deprecated      bool
	summary                 string
	params                  []*param
}

// param is one key of the call_tool params dict, or one field of an object
// (in == "").
type param struct {
	name       string // OpenAPI name; overrides match on it
	key        string // call key: x-mcp-name, else name
	in         string // path, query, header, cookie or body; "" for a field
	required   bool
	typ        string
	def        string
	class      class // escalation; classRead never escalates
	deprecated bool
	desc       string
	fields     []*param // fields of an inline object
	refs       []string // named types used by typ or fields
}

// typeDef is a named type used inside a request: an object's fields, or for
// a recursive non-object type, alias = its type expression.
type typeDef struct {
	name, base  string
	alias, note string
	fields      []*param
	refs        []string
	built       bool
}

// renderer turns a loaded spec into the reference.
type renderer struct {
	doc       *openapi3.T
	raw       *yaml.Node // the root document as written, for spec order and descriptions
	ov        overrides
	opts      options
	warnings  []string
	err       error                       // first error found while building nested schemas
	names     map[string]string           // operationId -> tool name
	typeDefs  map[string]*typeDef         // by type name
	typeNames map[*openapi3.Schema]string // schema -> type name
	active    map[*openapi3.Schema]bool   // schemas being rendered, to cut recursion
	recursive map[*openapi3.Schema]bool   // schemas found to contain themselves
}

func newRenderer(src *source, ov overrides, opts options) *renderer {
	return &renderer{doc: src.doc, raw: src.raw, ov: ov, opts: opts,
		names: map[string]string{}, typeDefs: map[string]*typeDef{},
		typeNames: map[*openapi3.Schema]string{}, active: map[*openapi3.Schema]bool{},
		recursive: map[*openapi3.Schema]bool{}}
}

func (r *renderer) warnf(format string, args ...any) {
	r.warnings = append(r.warnings, fmt.Sprintf(format, args...))
}

// fail records the first error found where no error can be returned.
func (r *renderer) fail(err error) {
	if r.err == nil {
		r.err = err
	}
}

type rawOp struct {
	path, method string
	item         *openapi3.PathItem
	op           *openapi3.Operation
	fields       mcpFields
}

var mcpToolNameRE = regexp.MustCompile(`^[A-Za-z0-9_.-]{1,128}$`)

// operations collects the operations in spec order and builds those that
// will be emitted, applying extensions and then overrides.
func (r *renderer) operations() ([]*operation, error) {
	var raws []rawOp
	where := map[string]string{}     // operationId -> "METHOD path"
	toolOwner := map[string]string{} // tool name -> operationId
	for _, path := range r.orderedPaths() {
		item := r.doc.Paths.Value(path)
		for _, m := range orderedMethods(item) {
			at := m.method + " " + path
			id := m.op.OperationID
			if id == "" {
				r.warnf("skipping %s: no operationId", at)
				continue
			}
			if prev, dup := where[id]; dup {
				return nil, fmt.Errorf("duplicate operationId %q (%s and %s)", id, prev, at)
			}
			where[id] = at
			f, err := extFields(m.op.Extensions, true)
			if err != nil {
				return nil, fmt.Errorf("%s: %w", id, err)
			}
			if o := r.ov[id]; o != nil {
				f = f.merge(&o.mcpFields)
			}
			name := toolName(id)
			if f.Name != nil {
				name = *f.Name
			}
			if other, dup := toolOwner[name]; dup {
				return nil, fmt.Errorf("operationIds %q and %q both map to tool name %q", other, id, name)
			}
			toolOwner[name] = id
			r.names[id] = name
			if !mcpToolNameRE.MatchString(name) {
				r.warnf("%s: tool name %q has characters MCP tool names can't use; set %s", id, name, extName)
			}
			raws = append(raws, rawOp{path: path, method: m.method, item: item, op: m.op, fields: f})
		}
	}
	for _, id := range sortedKeys(r.ov) {
		if _, ok := where[id]; !ok {
			r.warnf("overrides: operationId %q not found in spec", id)
		}
	}
	var ops []*operation
	for _, raw := range raws {
		if raw.fields.Hidden != nil && *raw.fields.Hidden && !r.opts.includeHidden {
			continue // not built, so it reserves no type names
		}
		o, err := r.operation(raw)
		if err == nil && r.err != nil {
			err = fmt.Errorf("%s: %w", raw.op.OperationID, r.err)
		}
		if err != nil {
			return nil, err
		}
		ops = append(ops, o)
	}
	return ops, nil
}

func (r *renderer) operation(raw rawOp) (*operation, error) {
	op, id := raw.op, raw.op.OperationID
	o := &operation{
		name:       r.names[id],
		method:     raw.method,
		path:       raw.path,
		class:      methodClass(raw.method),
		hidden:     raw.fields.Hidden != nil && *raw.fields.Hidden,
		deprecated: op.Deprecated,
		tag:        "untagged",
	}
	if raw.fields.Class != nil {
		o.class = *raw.fields.Class
	}
	if len(op.Tags) > 0 && strings.TrimSpace(op.Tags[0]) != "" {
		o.tag = oneLine(op.Tags[0])
	}
	text := op.Summary
	if text == "" {
		text = op.Description
	}
	o.summary = r.describe(text, nameWords(o.name))

	params, err := r.params(raw)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", id, err)
	}
	body, bodyNote, err := r.bodyParams(op.RequestBody)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", id, err)
	}
	var bodyDesc string
	if op.RequestBody != nil && op.RequestBody.Value != nil {
		bodyDesc = r.refDescription(op.RequestBody.Ref, op.RequestBody.Description, op.RequestBody.Value.Description)
	}
	note := r.opNote(op.Description, bodyDesc)
	if bodyNote != "" {
		note = joinNotes(note, "body "+bodyNote)
	}
	if note != "" && !strings.Contains(o.summary, note) {
		if o.summary != "" {
			o.summary += " · "
		}
		o.summary += note
	}
	o.params = append(params, body...)
	if err := r.applyParamOverrides(id, o.params); err != nil {
		return nil, err
	}
	seen := map[string]string{}
	for _, p := range o.params {
		if prev, dup := seen[p.key]; dup {
			r.warnf("%s: params %s and %s share the key %q; alias one with %s", id, prev, p.in+"."+p.name, p.key, extName)
		}
		seen[p.key] = p.in + "." + p.name
	}
	return o, nil
}

// params merges path-item and operation parameters (the operation wins),
// path params first in path-template order, the rest in spec order.
func (r *renderer) params(raw rawOp) ([]*param, error) {
	var merged []*openapi3.ParameterRef
	index := map[string]int{}
	add := func(refs openapi3.Parameters) {
		for _, ref := range refs {
			if ref == nil || ref.Value == nil {
				continue
			}
			p := ref.Value
			k := p.In + "\x00" + p.Name
			if p.In == openapi3.ParameterInHeader {
				k = strings.ToLower(k)
				switch strings.ToLower(p.Name) { // OpenAPI says these are ignored
				case "accept", "content-type", "authorization":
					continue
				}
			}
			if i, ok := index[k]; ok {
				merged[i] = ref
				continue
			}
			index[k] = len(merged)
			merged = append(merged, ref)
		}
	}
	add(raw.item.Parameters)
	add(raw.op.Parameters)

	template := pathTemplateParams(raw.path)
	sort.SliceStable(merged, func(i, j int) bool {
		return paramRank(merged[i].Value, template) < paramRank(merged[j].Value, template)
	})
	declared := map[string]bool{}
	out := make([]*param, 0, len(merged))
	for _, ref := range merged {
		p := ref.Value
		if p.In == openapi3.ParameterInPath {
			declared[p.Name] = true
		}
		schema := p.Schema
		if schema == nil {
			if _, mt, _ := pickMedia(p.Content); mt != nil {
				schema = mt.Schema
			}
		}
		desc := r.refDescription(ref.Ref, ref.Description, p.Description)
		if desc == "" && schema != nil && schema.Value != nil {
			desc = schema.Value.Description
		}
		sh := r.shape(schema, 0, p.Name)
		q := &param{
			name:       p.Name,
			key:        p.Name,
			in:         p.In,
			required:   p.Required || p.In == openapi3.ParameterInPath,
			typ:        sh.top(),
			def:        defaultOf(schema),
			deprecated: p.Deprecated || schema != nil && schema.Value != nil && schema.Value.Deprecated,
			desc:       joinNotes(r.describeParam(desc, nameWords(p.Name), schema), sh.note),
			fields:     sh.fields,
			refs:       sh.refs,
		}
		// x-mcp-* on the parameter's schema, then on the parameter itself.
		var f mcpFields
		for _, ext := range schemaExtensions(schema) {
			g, err := extFields(ext, false)
			if err != nil {
				return nil, fmt.Errorf("parameter %s: %w", p.Name, err)
			}
			f = f.merge(&g)
		}
		g, err := extFields(p.Extensions, false)
		if err != nil {
			return nil, fmt.Errorf("parameter %s: %w", p.Name, err)
		}
		q.apply(f.merge(&g))
		out = append(out, q)
	}
	for _, name := range template {
		if !declared[name] {
			r.warnf("%s: path parameter {%s} is not declared", raw.op.OperationID, name)
		}
	}
	return out, nil
}

// schemaExtensions returns the extension maps of a schema and of the $ref
// pointing at it (keywords beside a $ref win, so they come last).
func schemaExtensions(ref *openapi3.SchemaRef) []map[string]any {
	if ref == nil || ref.Value == nil {
		return nil
	}
	return []map[string]any{ref.Value.Extensions, ref.Extensions}
}

// refDescription returns a component's description as used at one $ref. A
// 3.1 Reference Object may override it; the loader then writes the override
// into the shared component, so other uses read the original from the raw
// document.
func (r *renderer) refDescription(ref string, override *string, value string) string {
	if ref == "" {
		return value
	}
	if override != nil {
		return *override
	}
	if d, ok := r.rawString(ref, "description"); ok {
		return d
	}
	return value
}

// rawString reads key from the object a local "#/…" pointer names in the raw
// root document; ok is false when the pointer can't be followed there.
func (r *renderer) rawString(pointer, key string) (string, bool) {
	if r.raw == nil || !strings.HasPrefix(pointer, "#/") {
		return "", false
	}
	n := r.raw
	if n.Kind == yaml.DocumentNode && len(n.Content) == 1 {
		n = n.Content[0]
	}
	for _, seg := range strings.Split(strings.TrimPrefix(pointer, "#/"), "/") {
		n = mapValue(n, strings.NewReplacer("~1", "/", "~0", "~").Replace(seg))
		if n == nil {
			return "", false
		}
	}
	if n.Kind != yaml.MappingNode {
		return "", false
	}
	if v := mapValue(n, key); v != nil && v.Kind == yaml.ScalarNode {
		return v.Value, true
	}
	return "", true
}

func mapValue(n *yaml.Node, key string) *yaml.Node {
	if n == nil || n.Kind != yaml.MappingNode {
		return nil
	}
	for i := 0; i+1 < len(n.Content); i += 2 {
		if n.Content[i].Value == key {
			return n.Content[i+1]
		}
	}
	return nil
}

// apply sets the param's escalation class and call key from x-mcp fields.
func (p *param) apply(f mcpFields) {
	if f.Class != nil {
		p.class = *f.Class
	}
	if f.Name != nil {
		p.key = *f.Name
	}
}

// paramRank orders path params by their position in the path template, then
// everything else (the sort is stable).
func paramRank(p *openapi3.Parameter, template []string) int {
	if p.In == openapi3.ParameterInPath {
		for i, n := range template {
			if n == p.Name {
				return i
			}
		}
		return len(template)
	}
	return len(template) + 1
}

var templateParamRE = regexp.MustCompile(`\{([^{}]+)\}`)

func pathTemplateParams(path string) []string {
	var names []string
	for _, m := range templateParamRE.FindAllStringSubmatch(path, -1) {
		names = append(names, m[1])
	}
	return names
}

// bodyParams turns a request body into "body" params. An object's properties
// (allOf merged) are flattened, and a oneOf/anyOf beside them comes back as a
// note for the operation; a body that is only a union, or has no properties,
// becomes one param named body. A body that must be empty adds nothing.
func (r *renderer) bodyParams(ref *openapi3.RequestBodyRef) ([]*param, string, error) {
	if ref == nil || ref.Value == nil {
		return nil, "", nil
	}
	rb := ref.Value
	media, mt, others := pickMedia(rb.Content)
	if mt == nil {
		return nil, "", nil
	}
	if len(others) > 1 {
		media = strings.Join(others, "|")
	}
	whole := &param{name: "body", key: "body", in: "body", required: rb.Required}
	desc := r.refDescription(ref.Ref, ref.Description, rb.Description)
	if mt.Schema == nil || mt.Schema.Value == nil {
		whole.typ = "string"
		if isJSON(media) {
			whole.typ = "any"
		}
		whole.desc = r.describe(mediaNote(media, desc), nil)
		return []*param{whole}, "", nil
	}
	s := effective(mt.Schema).Value
	if desc == "" {
		desc = mt.Schema.Value.Description
	}
	if !hasStructure(s) {
		if emptyObject(s) {
			return nil, "", nil
		}
		sh := r.shape(mt.Schema, 0, "body")
		whole.typ, whole.fields, whole.refs = sh.top(), sh.fields, sh.refs
		if _, kind := unionOf(s); kind != "" && r.allObjectTypes(sh.typ) {
			// A union of objects: the note names the alternatives (under types).
			whole.typ = "object"
			if sh.nullable {
				whole.typ += " (nullable)"
			}
			sh.note = joinNotes(kind+" "+strings.ReplaceAll(sh.typ, "|", " | "), sh.note)
		}
		whole.desc = joinNotes(sh.note, r.describe(mediaNote(media, desc), nil))
		return []*param{whole}, "", nil
	}
	props, required, unions := flattenObject(mt.Schema)
	own := map[string]bool{}
	out := make([]*param, 0, len(props))
	for _, pr := range props {
		own[pr.name] = true
		q := r.field(pr, required[pr.name], 0)
		q.in = "body"
		out = append(out, q)
	}
	note, refs := r.constraintNote(unions, own, 0, "body", nil)
	if len(refs) > 0 && len(out) > 0 {
		out[0].refs = appendUnique(out[0].refs, refs...)
	}
	return out, note, nil
}

// allObjectTypes reports whether a union type lists only object types
// defined under types.
func (r *renderer) allObjectTypes(union string) bool {
	parts := strings.Split(union, "|")
	for _, p := range parts {
		if t := r.typeDefs[p]; t == nil || t.alias != "" {
			return false
		}
	}
	return len(parts) > 1
}

// emptyObject reports whether a schema admits only {}: an object with no
// properties, unions or composition, and additionalProperties: false.
func emptyObject(s *openapi3.Schema) bool {
	ap := s.AdditionalProperties
	return s.Type.Is("object") && len(s.Properties) == 0 && len(s.PatternProperties) == 0 &&
		!isUnion(s) && len(s.AllOf) == 0 && ap.Has != nil && !*ap.Has && ap.Schema == nil
}

// mediaNote prefixes a non-JSON media type to a body description, since the
// agent must send that kind of content.
func mediaNote(media, desc string) string {
	if media == "" || isJSON(media) {
		return desc
	}
	if desc == "" {
		return media
	}
	return media + ": " + desc
}

// field builds the param for one object property. An invalid x-mcp-* value
// is recorded as the renderer's error.
func (r *renderer) field(pr property, required bool, depth int) *param {
	v := pr.ref.Value
	sh := r.shape(pr.ref, depth, pr.name)
	q := &param{
		name:       pr.name,
		key:        pr.name,
		required:   required,
		typ:        sh.top(),
		def:        defaultOf(pr.ref),
		deprecated: v.Deprecated,
		desc:       joinNotes(r.describeParam(v.Description, nameWords(pr.name), pr.ref), sh.note),
		fields:     sh.fields,
		refs:       sh.refs,
	}
	var f mcpFields
	for _, ext := range schemaExtensions(pr.ref) {
		g, err := extFields(ext, false)
		if err != nil {
			r.fail(fmt.Errorf("property %s: %w", pr.name, err))
		}
		f = f.merge(&g)
	}
	q.apply(f)
	return q
}

type property struct {
	name string
	ref  *openapi3.SchemaRef
}

// flattenObject merges the properties and required lists of a schema and
// its allOf members, recursively, in spec order, and returns the schemas
// among them that carry a oneOf/anyOf. A later member that redeclares a
// property narrows it (its enum, const, default, format and description
// win). readOnly properties are dropped: they never belong in a request.
func flattenObject(root *openapi3.SchemaRef) ([]property, map[string]bool, []*openapi3.SchemaRef) {
	var props []property
	at := map[string]int{}
	required := map[string]bool{}
	var unions []*openapi3.SchemaRef
	visited := map[*openapi3.Schema]bool{}
	var walk func(*openapi3.SchemaRef)
	walk = func(ref *openapi3.SchemaRef) {
		if ref == nil || ref.Value == nil || visited[ref.Value] {
			return
		}
		s := ref.Value
		visited[s] = true
		for _, m := range s.AllOf {
			walk(m)
		}
		for _, name := range orderedProperties(s.Properties) {
			p := s.Properties[name]
			if p == nil || p.Value == nil || p.Value.ReadOnly {
				continue
			}
			if i, ok := at[name]; ok {
				props[i].ref = narrow(props[i].ref, p)
				continue
			}
			at[name] = len(props)
			props = append(props, property{name, p})
		}
		for _, n := range s.Required {
			required[n] = true
		}
		if isUnion(s) {
			unions = append(unions, ref)
		}
	}
	walk(root)
	return props, required, unions
}

// narrow applies a later redeclaration of a property to the earlier one. A
// redeclared object replaces it; scalar keywords are laid over a copy.
func narrow(base, later *openapi3.SchemaRef) *openapi3.SchemaRef {
	b, l := base.Value, later.Value
	if hasStructure(l) || len(l.Type.Slice()) > 0 && len(b.Type.Slice()) == 0 {
		return later
	}
	c := *b
	if len(l.Type.Slice()) > 0 {
		c.Type = l.Type
	}
	if l.Format != "" {
		c.Format = l.Format
	}
	if len(l.Enum) > 0 {
		c.Enum = l.Enum
	}
	if l.Const != nil {
		c.Const = l.Const
	}
	if l.Default != nil {
		c.Default = l.Default
	}
	if l.Description != "" {
		c.Description = l.Description
	}
	return &openapi3.SchemaRef{Value: &c, Extensions: later.Extensions}
}

// applyParamOverrides applies an overrides entry's parameters, matched by
// OpenAPI name, or by "in.name" when the name alone is ambiguous.
func (r *renderer) applyParamOverrides(id string, params []*param) error {
	o := r.ov[id]
	if o == nil {
		return nil
	}
	for _, key := range sortedKeys(o.params) {
		in, name := "", key
		if i := strings.IndexByte(key, '.'); i > 0 {
			switch key[:i] {
			case "path", "query", "header", "cookie", "body":
				in, name = key[:i], key[i+1:]
			}
		}
		var hits []*param
		for _, p := range params {
			if (in == "" || p.in == in) && (p.name == name || p.in == "header" && strings.EqualFold(p.name, name)) {
				hits = append(hits, p)
			}
		}
		switch len(hits) {
		case 0:
			r.warnf("overrides: %s: parameter %q not found", id, key)
		case 1:
			hits[0].apply(*o.params[key])
		default:
			return fmt.Errorf("overrides: %s: parameter %q is ambiguous; use in.name (e.g. %s.%s)", id, key, hits[0].in, name)
		}
	}
	return nil
}
