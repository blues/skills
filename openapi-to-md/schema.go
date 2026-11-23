package main

import (
	"encoding/json"
	"fmt"
	"path"
	"sort"
	"strconv"
	"strings"
	"unicode"

	"github.com/getkin/kin-openapi/openapi3"
)

// maxDepth is a last-resort bound on nesting; recursion is cut by the
// renderer's cycle guard long before it.
const maxDepth = 32

// shaped is a schema rendered for a param or field line.
type shaped struct {
	typ      string   // without nullability
	nullable bool     // rendered " (nullable)" at the top, "|null" when nested
	fields   []*param // fields of an inline object (or of the items it holds)
	refs     []string // named types used
	note     string   // union, discriminator or constraint note for the description
}

// top is the type as a param line shows it.
func (s shaped) top() string {
	if s.nullable {
		return s.typ + " (nullable)"
	}
	return s.typ
}

// inner is the type as it shows inside array[...], map[string]... or a union.
func (s shaped) inner() string {
	if s.nullable {
		return s.typ + "|null"
	}
	return s.typ
}

// shape renders a schema: integer(int32), array[string], string[a|b] (enum),
// map[string]integer, string|integer, Foo (a named type). A $ref'd object
// with properties renders as its type name and is defined once under types;
// an inline one renders as "object" with its fields returned for indented
// display. hint names inline union alternatives that have no title.
func (r *renderer) shape(ref *openapi3.SchemaRef, depth int, hint string) shaped {
	if ref == nil || ref.Value == nil {
		return shaped{typ: "any"}
	}
	wrapperNullable := ref.Value.Nullable
	ref = effective(ref)
	s := ref.Value
	if r.active[s] { // recursion: refer to the type by name
		if ref.Ref != "" {
			r.recursive[s] = true
			name := r.typeName(ref)
			return shaped{typ: name, refs: []string{name}}
		}
		return shaped{typ: "any"}
	}
	if depth > maxDepth {
		return shaped{typ: "any"}
	}
	r.active[s] = true
	defer delete(r.active, s)

	var out shaped
	switch {
	case hasStructure(s) && ref.Ref != "":
		name := r.namedType(ref)
		out = shaped{typ: name, refs: []string{name}}
	case hasStructure(s):
		out = r.inlineObject(ref, depth, hint)
	case isUnion(s):
		out = r.union(ref, depth, hint)
	default:
		out = r.plain(ref, depth, hint)
	}
	out.nullable = out.nullable || wrapperNullable || s.Nullable || s.Type.Includes("null")
	if vals, hasNull := enumValues(s); len(vals) > 0 {
		out.typ += "[" + strings.Join(vals, "|") + "]"
		out.nullable = out.nullable || hasNull
	}
	if r.recursive[s] && ref.Ref != "" && !hasStructure(s) {
		// A recursive non-object type is defined once, as Name = <type>.
		name := r.typeName(ref)
		def := r.typeDefs[name]
		def.alias, def.fields, def.refs, def.built = out.inner(), out.fields, out.refs, true
		return shaped{typ: name, refs: []string{name}, note: out.note}
	}
	return out
}

// plain renders a scalar, array or map schema.
func (r *renderer) plain(ref *openapi3.SchemaRef, depth int, hint string) shaped {
	s := ref.Value
	var out shaped
	var parts []string
	for _, t := range inferTypes(s) {
		switch t {
		case "array":
			items := r.shape(s.Items, depth+1, hint)
			out.fields = append(out.fields, items.fields...)
			out.refs = appendUnique(out.refs, items.refs...)
			parts = append(parts, "array["+items.inner()+"]")
		case "object":
			if ap := s.AdditionalProperties.Schema; ap != nil {
				v := r.shape(ap, depth+1, hint)
				out.fields = append(out.fields, v.fields...)
				out.refs = appendUnique(out.refs, v.refs...)
				vt := v.inner()
				if strings.Contains(vt, "|") {
					vt = "(" + vt + ")"
				}
				parts = append(parts, "map[string]"+vt)
			} else {
				parts = append(parts, "object")
			}
		default:
			if f := formatOf(s); f != "" {
				t += "(" + f + ")"
			}
			parts = append(parts, t)
		}
	}
	out.typ = strings.Join(parts, "|")
	return out
}

// inlineObject renders an object without a type name: its fields, plus a
// note for a oneOf/anyOf combined with its properties.
func (r *renderer) inlineObject(ref *openapi3.SchemaRef, depth int, hint string) shaped {
	out := shaped{typ: "object"}
	props, required, unions := flattenObject(ref)
	own := map[string]bool{}
	for _, pr := range props {
		own[pr.name] = true
		f := r.field(pr, required[pr.name], depth+1)
		out.fields = append(out.fields, f)
		out.refs = appendUnique(out.refs, f.refs...)
	}
	out.note, out.refs = r.constraintNote(unions, own, depth, hint, out.refs)
	return out
}

// union renders oneOf/anyOf alternatives: "A|B", with inline object
// alternatives defined as types named by their title (or hint and index).
// A single non-null alternative renders as itself, nullable.
func (r *renderer) union(ref *openapi3.SchemaRef, depth int, hint string) shaped {
	s := ref.Value
	alts, _ := unionOf(s)
	var out shaped
	var live openapi3.SchemaRefs
	for _, a := range alts {
		switch {
		case a == nil || a.Value == nil:
		case isNullType(a.Value):
			out.nullable = true
		default:
			live = append(live, a)
		}
	}
	switch len(live) {
	case 0:
		return shaped{typ: "null"}
	case 1: // X|null is just a nullable X
		one := r.shape(live[0], depth+1, hint)
		one.nullable = one.nullable || out.nullable
		return one
	}
	var parts []string
	for i, a := range live {
		sh := r.shape(a, depth+1, hint)
		if sh.typ == "object" && len(sh.fields) > 0 {
			// Several inline objects can't share one indented block.
			name := r.syntheticType(a, hint, i, sh)
			sh = shaped{typ: name, refs: []string{name}, nullable: sh.nullable}
		}
		parts = appendUnique(parts, sh.typ)
		out.refs = appendUnique(out.refs, sh.refs...)
		out.nullable = out.nullable || sh.nullable
	}
	out.typ = strings.Join(parts, "|")
	out.note = r.discriminatorNote(s, alts)
	return out
}

// constraintNote describes oneOf/anyOf found beside an object's properties
// (or in its allOf members). Alternatives that only list required properties
// become "requires one of: a | b+c"; others are named types the object's
// fields combine with: "plus the fields of one of: A | B".
func (r *renderer) constraintNote(unions []*openapi3.SchemaRef, own map[string]bool, depth int, hint string, refs []string) (string, []string) {
	var notes []string
	for _, u := range unions {
		alts, _ := unionOf(u.Value)
		var reqs, names []string
		for i, a := range alts {
			if a == nil || a.Value == nil {
				continue
			}
			if a.Ref == "" && refinesOnly(a.Value, own) {
				// Narrowed values show in the property's own type already.
				if len(a.Value.Required) > 0 {
					reqs = append(reqs, strings.Join(a.Value.Required, "+"))
				}
				continue
			}
			sh := r.shape(a, depth+1, hint)
			if sh.typ == "object" && len(sh.fields) > 0 {
				name := r.syntheticType(a, hint, i, sh)
				sh = shaped{typ: name, refs: []string{name}}
			}
			names = appendUnique(names, sh.typ)
			refs = appendUnique(refs, sh.refs...)
		}
		if len(reqs) > 0 {
			notes = append(notes, "requires one of: "+strings.Join(reqs, " | "))
		}
		if len(names) > 0 {
			n := "plus the fields of one of: " + strings.Join(names, " | ")
			if d := r.discriminatorNote(u.Value, alts); d != "" {
				n += "; " + d
			}
			notes = append(notes, n)
		}
	}
	return strings.Join(notes, "; "), refs
}

// refinesOnly reports whether a union alternative only requires or narrows
// properties the object already has.
func refinesOnly(s *openapi3.Schema, own map[string]bool) bool {
	if isUnion(s) || len(s.AllOf) > 0 {
		return false
	}
	for name := range s.Properties {
		if !own[name] {
			return false
		}
	}
	for _, name := range s.Required {
		if !own[name] {
			return false
		}
	}
	return true
}

// discriminatorNote lists the values a discriminator property takes:
// "petType: cat=Cat | dog=Dog". Without a mapping the value is the
// alternative's component name.
func (r *renderer) discriminatorNote(s *openapi3.Schema, alts openapi3.SchemaRefs) string {
	d := s.Discriminator
	if d == nil || d.PropertyName == "" {
		return ""
	}
	var pairs []string
	if len(d.Mapping) > 0 {
		for _, v := range sortedKeys(d.Mapping) {
			pairs = append(pairs, literal(v)+"="+r.mappedName(d.Mapping[v].Ref))
		}
	} else {
		for _, a := range alts {
			if a != nil && a.Ref != "" {
				pairs = append(pairs, literal(refName(a.Ref))+"="+r.typeName(effective(a)))
			}
		}
	}
	if len(pairs) == 0 {
		return ""
	}
	return d.PropertyName + ": " + strings.Join(pairs, " | ")
}

// mappedName is the type name for a discriminator mapping target, which may
// be a $ref or a bare component name.
func (r *renderer) mappedName(target string) string {
	name := target
	if !strings.Contains(target, "/") {
		target = "#/components/schemas/" + target
	}
	if comps := r.doc.Components; comps != nil && strings.HasPrefix(target, "#/components/schemas/") {
		if sr := comps.Schemas[strings.TrimPrefix(target, "#/components/schemas/")]; sr != nil && sr.Value != nil {
			ref := &openapi3.SchemaRef{Ref: target, Value: sr.Value}
			if hasStructure(sr.Value) {
				return r.namedType(ref)
			}
			return r.shape(ref, 1, "").typ
		}
	}
	return cleanTypeName(name)
}

// namedType returns the type name of a $ref'd object schema, building its
// definition (fields, and a discriminator note) on first use.
func (r *renderer) namedType(ref *openapi3.SchemaRef) string {
	name := r.typeName(ref)
	def := r.typeDefs[name]
	if def.built {
		return name
	}
	def.built = true // before its fields, so cycles end here
	props, required, unions := flattenObject(ref)
	own := map[string]bool{}
	for _, pr := range props {
		own[pr.name] = true
		f := r.field(pr, required[pr.name], 1)
		def.fields = append(def.fields, f)
		def.refs = appendUnique(def.refs, f.refs...)
	}
	def.note, def.refs = r.constraintNote(unions, own, 1, name, def.refs)
	if d := ref.Value.Discriminator; d != nil && len(d.Mapping) > 0 && len(unions) == 0 {
		var pairs []string
		for _, v := range sortedKeys(d.Mapping) {
			t := r.mappedName(d.Mapping[v].Ref)
			pairs = append(pairs, literal(v)+"="+t)
			if r.typeDefs[t] != nil && t != name {
				def.refs = appendUnique(def.refs, t)
			}
		}
		def.note = joinNotes(def.note, d.PropertyName+": "+strings.Join(pairs, " | "))
	}
	return r.dedupe(ref, def)
}

// syntheticType defines an inline union alternative as a named type.
func (r *renderer) syntheticType(alt *openapi3.SchemaRef, hint string, index int, sh shaped) string {
	if n, ok := r.typeNames[alt.Value]; ok {
		return n
	}
	base := identifier(alt.Value.Title)
	if base == "" {
		base = identifier(hint) + strconv.Itoa(index+1)
		if base == strconv.Itoa(index+1) {
			base = "Option" + base
		}
	}
	name := r.reserveName(base)
	r.typeNames[alt.Value] = name
	r.typeDefs[name] = &typeDef{name: name, base: base, fields: sh.fields, refs: sh.refs, note: sh.note, built: true}
	return name
}

// typeName returns (reserving on first use) the name of a $ref'd schema.
// Types are keyed by schema, so same-named schemas from different files get
// distinct names; names are the component name with its first letter
// capitalized, so they never read as primitive types.
func (r *renderer) typeName(ref *openapi3.SchemaRef) string {
	if n, ok := r.typeNames[ref.Value]; ok {
		return n
	}
	base := cleanTypeName(ref.Ref)
	name := r.reserveName(base)
	r.typeNames[ref.Value] = name
	r.typeDefs[name] = &typeDef{name: name, base: base}
	return name
}

func (r *renderer) reserveName(base string) string {
	name := base
	for i := 2; r.typeDefs[name] != nil; i++ {
		name = base + strconv.Itoa(i)
	}
	return name
}

// dedupe folds a definition into an earlier one with the same base name and
// the same content, as OpenAPI 3.1 $ref siblings produce per-use copies.
func (r *renderer) dedupe(ref *openapi3.SchemaRef, def *typeDef) string {
	if def.name == def.base {
		return def.name
	}
	same := r.typeDefs[def.base]
	if same == nil || same.sig() != def.sig() || strings.Contains(def.sig(), def.name) {
		return def.name
	}
	delete(r.typeDefs, def.name)
	r.typeNames[ref.Value] = same.name
	return same.name
}

// cleanTypeName derives a type name from a $ref: the component name after
// the last JSON pointer segment that isn't a keyword or an index, or a
// file's base name without its extension; first letter capitalized.
func cleanTypeName(ref string) string {
	file, frag, _ := strings.Cut(ref, "#")
	if frag == "" || frag == "/" {
		base := path.Base(file)
		return identifier(strings.TrimSuffix(base, path.Ext(base)))
	}
	segs := strings.Split(strings.TrimPrefix(frag, "/"), "/")
	for i := range segs {
		segs[i] = strings.NewReplacer("~1", "/", "~0", "~").Replace(segs[i])
	}
	switch { // canonical component locations name the schema directly
	case len(segs) == 3 && segs[0] == "components" && segs[1] == "schemas",
		len(segs) == 2 && (segs[0] == "$defs" || segs[0] == "definitions"),
		len(segs) == 1:
		if n := identifier(segs[len(segs)-1]); n != "" {
			return n
		}
	}
	for i := len(segs) - 1; i >= 0; i-- {
		if pointerKeywords[segs[i]] || isIndex(segs[i]) || identifier(segs[i]) == "" {
			continue
		}
		name := identifier(segs[i])
		for _, rest := range segs[i+1:] {
			name += identifier(rest)
		}
		return name
	}
	return "Type"
}

var pointerKeywords = map[string]bool{
	"components": true, "schemas": true, "$defs": true, "definitions": true, "properties": true,
	"items": true, "allOf": true, "oneOf": true, "anyOf": true, "not": true, "additionalProperties": true,
	"schema": true, "content": true, "parameters": true, "requestBody": true, "requestBodies": true,
	"responses": true, "paths": true, "patternProperties": true, "prefixItems": true,
}

func isIndex(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}

// identifier turns text into a type-name-like word: letters, digits, '_',
// '-' and '.', words joined CamelCase, first letter capitalized.
func identifier(s string) string {
	var b strings.Builder
	upper := true
	for _, c := range s {
		switch {
		case unicode.IsLetter(c) || unicode.IsDigit(c) || c == '_' || c == '-' || c == '.':
			if upper {
				c = unicode.ToUpper(c)
				upper = false
			}
			b.WriteRune(c)
		default:
			upper = true
		}
	}
	return b.String()
}

// effective returns the schema an allOf composition amounts to when it adds
// no properties: the 3.0 idiom {allOf: [{$ref: X}], description: …} becomes
// X itself (so a named type, enum or format survives), and scalar members are
// merged (types, formats, intersected enums), as OpenAPI 3.1 $ref siblings
// produce. Compositions with properties are left to flattenObject.
func effective(ref *openapi3.SchemaRef) *openapi3.SchemaRef {
	for range maxDepth {
		if ref == nil || ref.Value == nil || len(ref.Value.AllOf) == 0 || len(ref.Value.Properties) > 0 {
			return ref
		}
		s := ref.Value
		var members openapi3.SchemaRefs
		for _, m := range s.AllOf {
			if m != nil && m.Value != nil && !annotationOnly(m.Value) {
				members = append(members, m)
			}
		}
		if len(members) == 1 && !addsShape(s, members[0].Value) {
			ref = members[0] // the wrapper only annotates its one member
			continue
		}
		for _, m := range members {
			if hasStructure(m.Value) {
				return ref // an object composition: flattenObject merges it
			}
		}
		merged := *s
		merged.AllOf = nil
		for _, m := range members {
			v := effective(m).Value
			if len(merged.Type.Slice()) == 0 {
				merged.Type = v.Type
			}
			if merged.Format == "" {
				merged.Format = v.Format
			}
			if merged.Items == nil {
				merged.Items = v.Items
			}
			if merged.AdditionalProperties.Schema == nil {
				merged.AdditionalProperties = v.AdditionalProperties
			}
			if len(merged.OneOf)+len(merged.AnyOf) == 0 {
				merged.OneOf, merged.AnyOf, merged.Discriminator = v.OneOf, v.AnyOf, v.Discriminator
			}
			if merged.Const == nil {
				merged.Const = v.Const
			}
			merged.Enum = intersect(merged.Enum, v.Enum)
			merged.Nullable = merged.Nullable || v.Nullable
		}
		return &openapi3.SchemaRef{Value: &merged}
	}
	return ref
}

// addsShape reports whether an allOf wrapper constrains more than its one
// member: a type of its own (other than the member's), a format, values,
// items, map values or a union.
func addsShape(wrapper, member *openapi3.Schema) bool {
	if ts := wrapper.Type.Slice(); len(ts) > 0 {
		mt := inferTypes(member)
		if len(ts) != len(mt) || ts[0] != mt[0] {
			return true
		}
	}
	return wrapper.Format != "" || len(wrapper.Enum) > 0 || wrapper.Const != nil || wrapper.Items != nil ||
		wrapper.AdditionalProperties.Schema != nil || isUnion(wrapper)
}

// annotationOnly reports whether a schema says nothing about shape: only a
// description, default, nullable, deprecated, readOnly, examples and the like.
func annotationOnly(s *openapi3.Schema) bool {
	return len(s.Type.Slice()) == 0 && s.Format == "" && len(s.Enum) == 0 && s.Const == nil &&
		s.Items == nil && len(s.Properties) == 0 && s.AdditionalProperties.Schema == nil &&
		len(s.AllOf)+len(s.OneOf)+len(s.AnyOf) == 0 && len(s.Required) == 0
}

func intersect(a, b []any) []any {
	switch {
	case len(a) == 0:
		return b
	case len(b) == 0:
		return a
	}
	var out []any
	for _, x := range a {
		for _, y := range b {
			if fmt.Sprint(x) == fmt.Sprint(y) {
				out = append(out, x)
				break
			}
		}
	}
	return out
}

func isUnion(s *openapi3.Schema) bool { return len(s.OneOf)+len(s.AnyOf) > 0 }

func isNullType(s *openapi3.Schema) bool {
	return s.Type.Is("null") || (len(s.Enum) == 1 && s.Enum[0] == nil)
}

func unionOf(s *openapi3.Schema) (openapi3.SchemaRefs, string) {
	switch {
	case len(s.OneOf) > 0:
		return s.OneOf, "one of"
	case len(s.AnyOf) > 0:
		return s.AnyOf, "any of"
	}
	return nil, ""
}

// hasStructure reports whether a schema declares properties, itself or
// through allOf.
func hasStructure(s *openapi3.Schema) bool { return structured(s, map[*openapi3.Schema]bool{}) }

func structured(s *openapi3.Schema, seen map[*openapi3.Schema]bool) bool {
	if s == nil || seen[s] {
		return false
	}
	seen[s] = true
	if len(s.Properties) > 0 {
		return true
	}
	for _, m := range s.AllOf {
		if m != nil && structured(m.Value, seen) {
			return true
		}
	}
	return false
}

// hasEnum reports whether a schema, or the items it holds, lists its values.
func hasEnum(ref *openapi3.SchemaRef) bool {
	seen := map[*openapi3.Schema]bool{}
	for ref != nil && ref.Value != nil && !seen[ref.Value] {
		seen[ref.Value] = true
		ref = effective(ref)
		s := ref.Value
		if len(s.Enum) > 0 || s.Const != nil {
			return true
		}
		ref = s.Items
	}
	return false
}

// scalar reports whether a schema is a string or number (or an array of
// them): the kinds whose values a description can list.
func scalar(ref *openapi3.SchemaRef) bool {
	seen := map[*openapi3.Schema]bool{}
	for ref != nil && ref.Value != nil && !seen[ref.Value] {
		seen[ref.Value] = true
		ref = effective(ref)
		s := ref.Value
		if hasStructure(s) || isUnion(s) {
			return false
		}
		ts := inferTypes(s)
		if len(ts) != 1 {
			return false
		}
		switch ts[0] {
		case "string", "integer", "number":
			return true
		case "array":
			ref = s.Items
			continue
		}
		return false
	}
	return false
}

// inferTypes returns the declared non-null types, or a type inferred from
// the schema's shape when none is declared.
func inferTypes(s *openapi3.Schema) []string {
	if ts := nonNull(s.Type); len(ts) > 0 {
		return ts
	}
	switch {
	case hasStructure(s) || s.AdditionalProperties.Schema != nil:
		return []string{"object"}
	case s.Items != nil:
		return []string{"array"}
	case len(s.Enum) > 0 || s.Const != nil:
		v := s.Const
		for _, e := range s.Enum {
			if v == nil {
				v = e
			}
		}
		switch v.(type) {
		case string:
			return []string{"string"}
		case bool:
			return []string{"boolean"}
		case float64, int, int64:
			return []string{"number"}
		}
	}
	return []string{"any"}
}

// formatOf is the format, or for OpenAPI 3.1 the content encoding (base64)
// or a non-text content media type (binary).
func formatOf(s *openapi3.Schema) string {
	switch {
	case s.Format != "":
		return s.Format
	case s.ContentEncoding != "":
		return s.ContentEncoding
	case s.ContentMediaType != "" && !strings.HasPrefix(s.ContentMediaType, "text/") && !isJSON(s.ContentMediaType):
		return "binary"
	}
	return ""
}

func nonNull(t *openapi3.Types) []string {
	var out []string
	for _, s := range t.Slice() {
		if s != "null" {
			out = append(out, s)
		}
	}
	return out
}

// enumValues renders a schema's enum (or 3.1 const) values; null is reported
// separately so it renders as nullability.
func enumValues(s *openapi3.Schema) (vals []string, hasNull bool) {
	src := s.Enum
	if s.Const != nil {
		src = []any{s.Const}
	}
	for _, v := range src {
		if v == nil {
			hasNull = true
			continue
		}
		vals = append(vals, literal(v))
	}
	return vals, hasNull
}

// defaultOf renders a schema's default, omitting false and "", which are
// what an absent value means anyway.
func defaultOf(ref *openapi3.SchemaRef) string {
	if ref == nil || ref.Value == nil {
		return ""
	}
	d := ref.Value.Default
	if d == nil {
		d = effective(ref).Value.Default
	}
	switch x := d.(type) {
	case nil:
		return ""
	case bool:
		if !x {
			return ""
		}
	case string:
		if x == "" {
			return ""
		}
	}
	return literal(d)
}

// literal renders an enum or default value; strings are quoted only when
// they would otherwise be ambiguous or break the line.
func literal(v any) string {
	switch x := v.(type) {
	case nil:
		return "null"
	case string:
		if x == "" || strings.ContainsAny(x, " |[](),=\"'`") || strings.TrimSpace(x) != x ||
			strings.ContainsFunc(x, unicode.IsControl) {
			return strconv.Quote(x)
		}
		return x
	case bool:
		return strconv.FormatBool(x)
	case float64:
		return strconv.FormatFloat(x, 'f', -1, 64)
	}
	data, err := json.Marshal(v)
	if err != nil {
		return fmt.Sprint(v)
	}
	return string(data)
}

func isJSON(media string) bool {
	m := strings.ToLower(strings.TrimSpace(strings.SplitN(media, ";", 2)[0]))
	return m == "application/json" || strings.HasSuffix(m, "+json")
}

func mediaRank(k string) int {
	m := strings.ToLower(strings.TrimSpace(strings.SplitN(k, ";", 2)[0]))
	switch {
	case m == "application/json":
		return 0
	case isJSON(m):
		return 1
	case m == "application/x-www-form-urlencoded":
		return 2
	case m == "multipart/form-data":
		return 3
	}
	return 4
}

// pickMedia chooses the media type an MCP server most likely maps to call
// params: JSON first, then form encodings, then any other. It also returns
// every other-kind media type, which all accept the same param.
func pickMedia(c openapi3.Content) (string, *openapi3.MediaType, []string) {
	keys := sortedKeys(c)
	sort.SliceStable(keys, func(i, j int) bool { return mediaRank(keys[i]) < mediaRank(keys[j]) })
	if len(keys) == 0 {
		return "", nil, nil
	}
	var others []string
	if mediaRank(keys[0]) == 4 {
		others = keys
	}
	return keys[0], c[keys[0]], others
}

func refName(ref string) string { return ref[strings.LastIndexByte(ref, '/')+1:] }

func joinNotes(a, b string) string {
	switch {
	case a == "":
		return b
	case b == "":
		return a
	}
	return a + "; " + b
}
