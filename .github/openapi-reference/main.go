// Command openapi-reference prints a compact Markdown reference to every
// operation in the paths of a single-file OpenAPI 3 document: its operationId
// and call name, its effect, its method and path, its summary, and its
// parameters by location.
//
// It knows nothing about any particular API. An operation's effect follows
// its HTTP method unless the operation's x-effect extension, or an entry in
// an optional overrides file keyed by operationId, says otherwise. What it
// doesn't list (webhooks, callbacks, additionalOperations) is named in the
// reference's header, and references it can't follow are reported as
// warnings on standard error.
//
//	openapi-reference [-overrides overrides.yaml] [-spec-url URL] openapi.yaml > reference.md
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

// Effects an operation can have, and the HTTP methods that have each one by default.
const (
	effectRead        = "read"
	effectWrite       = "write"
	effectDestructive = "destructive"
)

var defaultEffects = map[string]string{
	"get": effectRead, "head": effectRead, "options": effectRead, "trace": effectRead, "query": effectRead,
	"post": effectWrite, "put": effectWrite, "patch": effectWrite,
	"delete": effectDestructive,
}

var (
	// The order operations and parameter locations are listed in. The QUERY method and
	// the querystring location are OpenAPI 3.2's.
	methods   = []string{"get", "head", "options", "trace", "query", "post", "put", "patch", "delete"}
	locations = []string{"path", "query", "querystring", "header", "cookie"}

	// Header parameters OpenAPI says to ignore: content negotiation and authorization.
	ignoredHeaders = map[string]bool{"accept": true, "content-type": true, "authorization": true}

	// Keywords whose values are data, not OpenAPI objects, so a "$ref" inside isn't a reference.
	dataKeys = map[string]bool{"example": true, "default": true, "enum": true, "const": true}

	// Maps whose keys are names the author chose rather than keywords, so an entry there
	// named "default" or "example" is an ordinary one.
	namedMaps = map[string]bool{
		"paths": true, "webhooks": true, "components": true, "schemas": true, "responses": true,
		"parameters": true, "examples": true, "requestBodies": true, "headers": true,
		"securitySchemes": true, "links": true, "callbacks": true, "pathItems": true,
		"mediaTypes": true, "content": true, "encoding": true, "properties": true,
		"patternProperties": true, "dependentSchemas": true, "$defs": true, "definitions": true,
	}
)

// What an overrides file, or an operation's own extensions, can say about it.
type override struct {
	Effect string `yaml:"effect"`
	Note   string `yaml:"note"`
}

func main() {
	overridesFile := flag.String("overrides", "", "YAML file of effects and notes keyed by operationId")
	specURL := flag.String("spec-url", "", "URL of the specification, linked from the reference")
	operationsDir := flag.String("operations-dir", "", "directory to write one contract file per operation into")
	operationsURL := flag.String("operations-url", "", "URL of that directory, to link each operation from the reference")
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "usage: openapi-reference [-overrides file] [-spec-url URL] [-operations-dir dir [-operations-url URL]] openapi.yaml")
		flag.PrintDefaults()
	}
	flag.Parse()
	if flag.NArg() != 1 || (*operationsURL != "" && *operationsDir == "") {
		flag.Usage()
		os.Exit(2)
	}

	data, err := os.ReadFile(flag.Arg(0))
	if err != nil {
		fail(err)
	}
	spec, root, err := decodeSpec(data)
	if err != nil {
		fail(fmt.Errorf("%s: %w", flag.Arg(0), err))
	}
	overrides := map[string]override{}
	if *overridesFile != "" {
		data, err := os.ReadFile(*overridesFile)
		if err == nil {
			err = yaml.Unmarshal(data, &overrides)
		}
		if err != nil {
			fail(fmt.Errorf("%s: %w", *overridesFile, err))
		}
	}

	// Without a URL, link each operation relative to a reference written beside the directory.
	linkURL := *operationsURL
	if *operationsDir != "" && linkURL == "" {
		linkURL = filepath.Base(*operationsDir)
	}
	reference, warnings, err := render(spec, overrides, *specURL, linkURL)
	for _, warning := range warnings {
		fmt.Fprintln(os.Stderr, "warning:", warning)
	}
	if err != nil {
		fail(err)
	}
	if *operationsDir != "" {
		files, warnings := contracts(root, spec, *specURL)
		for _, warning := range warnings {
			fmt.Fprintln(os.Stderr, "warning:", warning)
		}
		if err := writeContracts(*operationsDir, files); err != nil {
			fail(err)
		}
	}
	// A partial write must fail, so that a caller keeps its previous complete reference.
	if _, err := os.Stdout.WriteString(reference); err != nil {
		fail(fmt.Errorf("writing the reference: %w", err))
	}
}

func fail(err error) {
	fmt.Fprintln(os.Stderr, "openapi-reference:", err)
	os.Exit(1)
}

// Decode a YAML or JSON specification into plain maps, lists and scalars. Unlike decoding
// straight into Go values, a date-shaped value such as 2024-01-01 stays the text it is, as
// it would in JSON, rather than becoming a timestamp.
// It also returns the document's root node, from which contract files are cut unchanged.
func decodeSpec(data []byte) (map[string]any, *yaml.Node, error) {
	var doc yaml.Node
	if err := yaml.Unmarshal(data, &doc); err != nil {
		return nil, nil, err
	}
	v, err := plain(&doc, 0)
	if err != nil {
		return nil, nil, err
	}
	spec, ok := v.(map[string]any)
	if !ok || len(doc.Content) == 0 {
		return nil, nil, fmt.Errorf("not a mapping")
	}
	return spec, doc.Content[0], nil
}

func plain(n *yaml.Node, depth int) (any, error) {
	if depth > 10000 {
		return nil, fmt.Errorf("nested too deeply, or an alias refers to itself (line %d)", n.Line)
	}
	switch n.Kind {
	case yaml.DocumentNode:
		if len(n.Content) == 0 {
			return nil, nil
		}
		return plain(n.Content[0], depth+1)
	case yaml.AliasNode:
		return plain(n.Alias, depth+1)
	case yaml.SequenceNode:
		list := make([]any, 0, len(n.Content))
		for _, item := range n.Content {
			v, err := plain(item, depth+1)
			if err != nil {
				return nil, err
			}
			list = append(list, v)
		}
		return list, nil
	case yaml.MappingNode:
		m := make(map[string]any, len(n.Content)/2)
		var merges []map[string]any
		for i := 0; i+1 < len(n.Content); i += 2 {
			key, value := n.Content[i], n.Content[i+1]
			v, err := plain(value, depth+1)
			if err != nil {
				return nil, err
			}
			if key.Tag == "!!merge" { // <<: keys the mapping doesn't set itself
				switch merged := v.(type) {
				case map[string]any:
					merges = append(merges, merged)
				case []any:
					for _, each := range merged {
						if mm, ok := each.(map[string]any); ok {
							merges = append(merges, mm)
						}
					}
				}
				continue
			}
			m[key.Value] = v
		}
		for _, merged := range merges {
			for k, v := range merged {
				if _, set := m[k]; !set {
					m[k] = v
				}
			}
		}
		return m, nil
	case yaml.ScalarNode:
		if n.Tag == "!!timestamp" {
			return n.Value, nil
		}
		var v any
		if err := n.Decode(&v); err != nil {
			return nil, err
		}
		return v, nil
	}
	return nil, fmt.Errorf("unexpected YAML at line %d", n.Line)
}

// One operation, ready to print.
type operation struct {
	id, call, effect, method, path, summary, note string
	link                                          string // the operation's contract file, if any
	deprecated                                    bool
	otherTags                                     []string // tags beyond the one it is grouped under
	parameters                                    []string // "path: a*", "query: b, c[]", "body: d*", ...
}

// Resolves references within one document, and reports what it can't follow.
type resolver struct {
	spec map[string]any
	// From OpenAPI 3.1, keywords beside a schema's $ref apply too; in 3.0 they are ignored.
	siblings bool
	warn     func(string, ...any)
	warned   map[string]bool
}

func (r *resolver) warnOnce(format string, args ...any) {
	message := fmt.Sprintf(format, args...)
	if !r.warned[message] {
		r.warned[message] = true
		r.warn("%s", message)
	}
}

// Render the reference, with warnings about anything it had to skip or ignore.
// With operationsURL, each operationId links to its contract file there.
func render(spec map[string]any, overrides map[string]override, specURL, operationsURL string) (string, []string, error) {
	var warnings []string
	warn := func(format string, args ...any) { warnings = append(warnings, fmt.Sprintf(format, args...)) }

	version, _ := spec["openapi"].(string)
	if !strings.HasPrefix(version, "3.") {
		return "", warnings, fmt.Errorf("not an OpenAPI 3 document")
	}
	r := &resolver{spec: spec, siblings: !strings.HasPrefix(version, "3.0"), warn: warn, warned: map[string]bool{}}
	r.checkReferences(spec, "")

	info, _ := spec["info"].(map[string]any)
	title, _ := info["title"].(string)
	apiVersion := fmt.Sprint(info["version"])
	var unlisted []string // what the header says is not listed
	if webhooks, ok := spec["webhooks"].(map[string]any); ok && len(webhooks) > 0 {
		unlisted = append(unlisted, plural(len(webhooks), "webhook"))
		warn("the specification's webhooks (%s) are not listed", strings.Join(sortedKeys(webhooks), ", "))
	}
	paths, _ := spec["paths"].(map[string]any)

	groups := map[string][]operation{}
	seen := map[string]bool{}
	count, callbacks, additional := 0, 0, 0
	for _, path := range sortedKeys(paths) {
		if !strings.HasPrefix(path, "/") { // an x- extension, not a path
			continue
		}
		item := r.pathItem(paths[path])
		if extra, ok := item["additionalOperations"].(map[string]any); ok && len(extra) > 0 {
			additional += len(extra)
			warn("%s has additionalOperations (%s), which are not listed", path, strings.Join(sortedKeys(extra), ", "))
		}
		for _, method := range methods {
			raw, ok := r.resolve(item[method]).(map[string]any)
			if !ok {
				continue
			}
			op := r.describe(path, method, item, raw)
			label := op.id
			if label == "" {
				label = op.method + " " + path
				warn("%s has no operationId", label)
			} else if strings.IndexFunc(op.id, func(c rune) bool { return c > 127 }) >= 0 {
				warn("%s is not plain ASCII, so its call name may differ from Python's lowercasing in the snake_case rule", op.id)
			}
			seen[op.id] = true
			if operationsURL != "" && op.id != "" {
				op.link = strings.TrimSuffix(operationsURL, "/") + "/" + contractFile(op.id)
			}
			if cb, ok := raw["callbacks"].(map[string]any); ok && len(cb) > 0 {
				callbacks += len(cb)
				warn("%s has callbacks (%s), which are not listed", label, strings.Join(sortedKeys(cb), ", "))
			}

			if ext, ok := raw["x-effect"].(string); ok {
				op.effect = checkEffect(ext, label, "x-effect", op.effect, warn)
			}
			if ext, ok := raw["x-effect-note"].(string); ok {
				op.note = oneLine(ext)
			}
			if o, ok := overrides[op.id]; ok && op.id != "" {
				if o.Effect != "" {
					op.effect = checkEffect(o.Effect, label, "the overrides file", op.effect, warn)
				}
				if o.Note != "" {
					op.note = oneLine(o.Note)
				}
			}

			tag := ""
			if tags, ok := raw["tags"].([]any); ok && len(tags) > 0 {
				tag = fmt.Sprint(tags[0])
				for _, t := range tags[1:] {
					op.otherTags = append(op.otherTags, fmt.Sprint(t))
				}
			}
			groups[tag] = append(groups[tag], op)
			count++
		}
	}
	if callbacks > 0 {
		unlisted = append(unlisted, plural(callbacks, "callback"))
	}
	if additional > 0 {
		unlisted = append(unlisted, plural(additional, "additionalOperations entry", "additionalOperations entries"))
	}
	if count == 0 {
		warn("the specification has no operations in its paths")
	}
	for _, id := range sortedKeys(overrides) {
		if !seen[id] {
			warn("the overrides file names %s, which is not an operation in the specification", id)
		}
	}

	var b strings.Builder
	writeHeader(&b, title, apiVersion, specURL, count, unlisted, operationsURL != "")
	for _, tag := range tagOrder(spec, groups) {
		name := tag
		if name == "" {
			name = "Untagged"
		}
		fmt.Fprintf(&b, "\n## %s\n\n", name)
		for _, op := range groups[tag] {
			writeOperation(&b, op)
		}
	}
	return b.String(), warnings, nil
}

// Describe one operation from its path item and operation object.
func (r *resolver) describe(path, method string, item, raw map[string]any) operation {
	op := operation{
		method:     strings.ToUpper(method),
		path:       path,
		effect:     defaultEffects[method],
		deprecated: raw["deprecated"] == true,
	}
	op.id, _ = raw["operationId"].(string)
	if op.id != "" {
		op.call = snakeCase(op.id)
	}
	op.summary, _ = raw["summary"].(string)
	if op.summary == "" {
		description, _ := raw["description"].(string)
		op.summary, _, _ = strings.Cut(strings.TrimSpace(description), "\n")
	}
	op.summary = oneLine(op.summary)

	// Operation-level parameters replace path-level ones with the same name and location.
	type entry struct{ in, text string }
	entries := map[string]*entry{}
	var order []string
	for _, list := range []any{item["parameters"], raw["parameters"]} {
		params, _ := list.([]any)
		for _, p := range params {
			param, _ := r.resolve(p).(map[string]any)
			name, _ := param["name"].(string)
			in, _ := param["in"].(string)
			if name == "" || in == "" {
				r.warnOnce("%s %s has a parameter with no name or location", op.method, path)
				continue
			}
			if in == "header" && ignoredHeaders[strings.ToLower(name)] {
				continue
			}
			key := in + "\x00" + name
			if _, ok := entries[key]; !ok {
				order = append(order, key)
			}
			if param["deprecated"] == true {
				entries[key] = nil
				continue
			}
			required := in == "path" || param["required"] == true
			entries[key] = &entry{in, r.formatParameter(name, required, r.mergeSchemas([]any{r.parameterSchema(param)}))}
		}
	}
	lists := map[string][]string{}
	for _, key := range order {
		if e := entries[key]; e != nil {
			lists[e.in] = append(lists[e.in], e.text)
		}
	}
	for _, in := range locations {
		if len(lists[in]) > 0 {
			op.parameters = append(op.parameters, in+": "+strings.Join(lists[in], ", "))
		}
	}
	for _, in := range sortedKeys(lists) {
		if !contains(locations, in) {
			r.warnOnce("%s %s has parameters in %q, which is not a parameter location", op.method, path, in)
		}
	}
	if body := r.describeBody(raw["requestBody"]); body != "" {
		op.parameters = append(op.parameters, body)
	}
	return op
}

// A parameter's schema, from `schema` or else the first `content` entry.
func (r *resolver) parameterSchema(param map[string]any) any {
	if schema, ok := param["schema"]; ok {
		return schema
	}
	content, _ := param["content"].(map[string]any)
	for _, mediaType := range sortedKeys(content) {
		media, _ := r.resolve(content[mediaType]).(map[string]any) // 3.2 media types can be references
		return media["schema"]
	}
	return nil
}

// "name", with * if required, [] if an array, and ": a|b" if its values are listed.
func (r *resolver) formatParameter(name string, required bool, schema map[string]any) string {
	s := name
	if required {
		s += "*"
	}
	values, hasEnum := schema["enum"].([]any) // for an array, the permitted whole arrays
	if isArray(schema) {
		s += "[]"
		if !hasEnum {
			items := r.mergeSchemas([]any{schema["items"]})
			values, _ = items["enum"].([]any)
		}
	}
	if len(values) > 0 {
		s += ": " + formatValues(values)
	}
	return s
}

// Whether a schema is an array: type array, alone or with only null in a 3.1 type list. A
// list such as [array, string] is not, since a plain string is also valid.
func isArray(schema map[string]any) bool {
	switch t := schema["type"].(type) {
	case string:
		return t == "array"
	case []any:
		array := false
		for _, each := range t {
			switch each {
			case "array":
				array = true
			case "null":
			default:
				return false
			}
		}
		return array
	}
	return false
}

var plainValue = regexp.MustCompile(`^[A-Za-z0-9_.:/@+-]+$`)

// Values joined with |: bare when every one is a plain word, otherwise each written as JSON.
func formatValues(values []any) string {
	bare := true
	for _, v := range values {
		s, ok := v.(string)
		if !ok || !plainValue.MatchString(s) || s == "null" || s == "true" || s == "false" {
			bare = false
			break
		}
		if _, err := strconv.ParseFloat(s, 64); err == nil {
			bare = false
			break
		}
	}
	parts := make([]string, len(values))
	for i, v := range values {
		if bare {
			parts[i] = v.(string)
		} else {
			parts[i] = jsonText(v)
		}
	}
	return strings.Join(parts, "|")
}

func jsonText(v any) string {
	encoded, err := json.Marshal(v)
	if err != nil {
		encoded, _ = json.Marshal(fmt.Sprint(v))
	}
	return string(encoded)
}

// "body: a*, b", with whether the body is optional, its media types and any alternatives.
func (r *resolver) describeBody(raw any) string {
	body, ok := r.resolve(raw).(map[string]any)
	if !ok {
		return ""
	}
	content, _ := body["content"].(map[string]any)
	if len(content) == 0 {
		return ""
	}
	types := sortedKeys(content)
	var qualifiers []string
	if body["required"] != true {
		qualifiers = append(qualifiers, "optional")
	}
	label := func(fields string) string {
		if len(qualifiers) == 0 {
			return "body: " + fields
		}
		return "body (" + strings.Join(qualifiers, "; ") + "): " + fields
	}

	// List the fields of the JSON representation, preferring application/json.
	jsonType := ""
	for _, t := range types {
		if mediaTypeName(t) == "application/json" {
			jsonType = t
			break
		}
	}
	for _, t := range types {
		if jsonType == "" && isJSON(t) {
			jsonType = t
		}
	}
	if jsonType == "" {
		return label("(" + strings.Join(types, " or ") + ")")
	}
	if len(types) > 1 {
		var others []string
		for _, t := range types {
			if t != jsonType {
				others = append(others, t)
			}
		}
		qualifiers = append(qualifiers, "fields of "+jsonType+"; also "+strings.Join(others, " or "))
	}

	media, _ := r.resolve(content[jsonType]).(map[string]any) // 3.2 media types can be references
	properties, required, alternatives := r.objectProperties(media["schema"], nil)
	names := map[string]bool{}
	for name := range properties {
		names[name] = true
	}
	for name := range required { // a required key can have no property schema
		names[name] = true
	}
	var parts []string
	for _, name := range sortedKeys(names) {
		schema := r.mergeSchemas(properties[name])
		if schema["readOnly"] == true { // returned by the server, never sent
			continue
		}
		parts = append(parts, r.formatParameter(name, required[name], schema))
	}
	fields := strings.Join(parts, ", ")
	withFields, alone := describeAlternatives(alternatives)
	switch {
	case alone != "" && fields == "":
		fields = "(" + alone + ")"
	case alone != "":
		fields += ", plus " + withFields
	case fields == "":
		schema := r.mergeSchemas([]any{media["schema"]})
		description := jsonType
		if t, ok := schema["type"].(string); ok && t != "object" {
			description += " " + t
		}
		fields = "(" + description + ")"
	}
	return label(fields)
}

// A oneOf or anyOf in a request body schema: its keyword and how many schemas it offers.
type alternatives struct {
	kind  string
	count int
}

// How to describe a body's alternatives beside its common fields, and with none.
func describeAlternatives(groups []alternatives) (withFields, alone string) {
	switch {
	case len(groups) == 0:
		return "", ""
	case len(groups) > 1:
		alone = fmt.Sprintf("%d sets of alternative schemas; see the specification", len(groups))
		return "fields from " + alone, alone
	case groups[0].kind == "anyOf":
		alone = fmt.Sprintf("one or more of %d alternative schemas", groups[0].count)
	default:
		alone = fmt.Sprintf("one of %d alternative schemas", groups[0].count)
	}
	return "the fields of " + alone, alone
}

// A media type's type/subtype, lowercased and without parameters.
func mediaTypeName(t string) string {
	name, _, _ := strings.Cut(t, ";")
	return strings.ToLower(strings.TrimSpace(name))
}

// Whether a media type is JSON: application/json, or a structured +json type.
func isJSON(t string) bool {
	_, subtype, _ := strings.Cut(mediaTypeName(t), "/")
	return subtype == "json" || strings.HasSuffix(subtype, "+json")
}

// The properties of an object schema, including those from allOf, with every schema given for
// each property; which properties are required; and its oneOf/anyOf alternatives.
func (r *resolver) objectProperties(raw any, visiting map[string]bool) (map[string][]any, map[string]bool, []alternatives) {
	properties, required := map[string][]any{}, map[string]bool{}
	var groups []alternatives
	if ref := refOf(raw); ref != "" {
		if visiting[ref] {
			r.warnOnce("the schema %q includes itself through allOf", ref)
			return properties, required, groups
		}
		visiting = with(visiting, ref)
	}
	schema, ok := r.resolve(raw).(map[string]any)
	if !ok {
		return properties, required, groups
	}
	if own, ok := schema["properties"].(map[string]any); ok {
		for _, name := range sortedKeys(own) {
			properties[name] = append(properties[name], own[name])
		}
	}
	if names, ok := schema["required"].([]any); ok {
		for _, name := range names {
			required[fmt.Sprint(name)] = true
		}
	}
	for _, kind := range []string{"oneOf", "anyOf"} {
		if list, ok := schema[kind].([]any); ok && len(list) > 0 {
			groups = append(groups, alternatives{kind, len(list)})
		}
	}
	if parts, ok := schema["allOf"].([]any); ok {
		for _, part := range parts {
			p, req, alt := r.objectProperties(part, visiting)
			for _, name := range sortedKeys(p) {
				properties[name] = append(properties[name], p[name]...)
			}
			for name := range req {
				required[name] = true
			}
			groups = append(groups, alt...)
		}
	}
	return properties, required, groups
}

// Combine schemas that all apply to one value, including their allOf parts. The first to give
// a keyword wins, except that readOnly from any of them makes the value read-only, and enums
// intersect, since a value must satisfy each.
func (r *resolver) mergeSchemas(schemas []any) map[string]any {
	merged := map[string]any{}
	var add func(raw any, visiting map[string]bool)
	add = func(raw any, visiting map[string]bool) {
		if ref := refOf(raw); ref != "" {
			if visiting[ref] {
				r.warnOnce("the schema %q includes itself through allOf", ref)
				return
			}
			visiting = with(visiting, ref)
		}
		schema, _ := r.resolve(raw).(map[string]any)
		for key, value := range schema {
			switch key {
			case "allOf":
			case "readOnly":
				if value == true {
					merged[key] = true
				}
			case "enum":
				if existing, ok := merged[key].([]any); ok {
					merged[key] = intersect(existing, value)
				} else {
					merged[key] = value
				}
			default:
				if _, set := merged[key]; !set {
					merged[key] = value
				}
			}
		}
		if parts, ok := schema["allOf"].([]any); ok {
			for _, part := range parts {
				add(part, visiting)
			}
		}
	}
	for _, raw := range schemas {
		add(raw, nil)
	}
	return merged
}

// The values in a that are also in b.
func intersect(a []any, b any) []any {
	list, ok := b.([]any)
	if !ok {
		return a
	}
	in := map[string]bool{}
	for _, v := range list {
		in[jsonText(v)] = true
	}
	both := []any{}
	for _, v := range a {
		if in[jsonText(v)] {
			both = append(both, v)
		}
	}
	return both
}

// Follow local references until reaching a value that isn't one. A reference this can't
// follow is returned unchanged; checkReferences warns about it.
func (r *resolver) resolve(v any) any {
	visited := map[string]bool{}
	for {
		m, ok := v.(map[string]any)
		if !ok {
			return v
		}
		ref, ok := m["$ref"].(string)
		if !ok {
			return v
		}
		target, ok := r.lookup(ref)
		if !ok {
			return v
		}
		if visited[ref] {
			r.warnOnce("the reference %q is part of a cycle of references", ref)
			return nil
		}
		visited[ref] = true
		if r.siblings && len(m) > 1 {
			target = withSiblings(target, m)
		}
		v = target
	}
}

// A 3.1 reference with keywords beside it. Beside a schema, both the keywords and the target
// apply, so the result is the two as allOf; summary and description alone override the
// target's, as they do for any reference.
func withSiblings(target any, m map[string]any) any {
	t, ok := target.(map[string]any)
	if !ok {
		return target
	}
	siblings := map[string]any{}
	annotationsOnly := true
	for key, value := range m {
		if key == "$ref" {
			continue
		}
		siblings[key] = value
		if key != "summary" && key != "description" {
			annotationsOnly = false
		}
	}
	if !annotationsOnly {
		return map[string]any{"allOf": []any{t, siblings}}
	}
	merged := make(map[string]any, len(t)+len(siblings))
	for key, value := range t {
		merged[key] = value
	}
	for key, value := range siblings {
		merged[key] = value
	}
	return merged
}

// A path item, including operations defined beside a $ref to a shared one. A path item's $ref
// is not a plain reference, so in every version the fields beside it apply as well.
func (r *resolver) pathItem(v any) map[string]any {
	m, _ := v.(map[string]any)
	ref, ok := m["$ref"].(string)
	if !ok {
		return m
	}
	merged := map[string]any{}
	if target, ok := r.resolve(map[string]any{"$ref": ref}).(map[string]any); ok && refOf(target) == "" {
		for key, value := range target {
			merged[key] = value
		}
	}
	for key, value := range m {
		if key != "$ref" {
			merged[key] = value
		}
	}
	return merged
}

// The value a local JSON Pointer reference ("#/a/b", or "#" for the whole document) names.
// Pointer tokens are percent-decoded, as a URI fragment's are, before ~1 and ~0.
func (r *resolver) lookup(ref string) (any, bool) {
	if ref == "#" {
		return r.spec, true
	}
	if !strings.HasPrefix(ref, "#/") {
		return nil, false // another document, or an anchor
	}
	var node any = r.spec
	for _, token := range strings.Split(ref[2:], "/") {
		if unescaped, err := url.PathUnescape(token); err == nil {
			token = unescaped
		}
		token = strings.ReplaceAll(strings.ReplaceAll(token, "~1", "/"), "~0", "~")
		switch parent := node.(type) {
		case map[string]any:
			child, ok := parent[token]
			if !ok {
				return nil, false
			}
			node = child
		case []any:
			i, err := strconv.Atoi(token)
			if err != nil || i < 0 || i >= len(parent) {
				return nil, false
			}
			node = parent[i]
		default:
			return nil, false
		}
	}
	return node, true
}

// Warn once about each reference that can't be followed: to another document, to an anchor,
// a dynamic reference, or one that points to nothing. key is the key v sits under.
func (r *resolver) checkReferences(v any, key string) {
	switch node := v.(type) {
	case map[string]any:
		if ref, ok := node["$ref"].(string); ok {
			switch {
			case !strings.HasPrefix(ref, "#"):
				r.warnOnce("the reference %q points outside the document and is not followed", ref)
			case ref != "#" && !strings.HasPrefix(ref, "#/"):
				r.warnOnce("the reference %q names an anchor, which is not followed", ref)
			default:
				if _, ok := r.lookup(ref); !ok {
					r.warnOnce("the reference %q points to nothing", ref)
				} else {
					r.resolve(map[string]any{"$ref": ref}) // warns if the chain is a cycle
				}
			}
		}
		if ref, ok := node["$dynamicRef"].(string); ok {
			r.warnOnce("the dynamic reference %q is not followed", ref)
		}
		for _, k := range sortedKeys(node) {
			child := node[k]
			if !namedMaps[key] {
				if _, isList := child.([]any); dataKeys[k] || (k == "examples" && isList) {
					continue // data, including a JSON Schema examples list
				}
			}
			if key == "examples" { // an Example Object, whose value is data
				if example, ok := child.(map[string]any); ok {
					child = withoutKey(example, "value")
				}
			}
			r.checkReferences(child, k)
		}
	case []any:
		for _, each := range node {
			r.checkReferences(each, "")
		}
	}
}

// The $ref of a reference object, or "".
func refOf(v any) string {
	m, _ := v.(map[string]any)
	ref, _ := m["$ref"].(string)
	return ref
}

// A copy of set with one more member.
func with(set map[string]bool, member string) map[string]bool {
	copied := make(map[string]bool, len(set)+1)
	for k := range set {
		copied[k] = true
	}
	copied[member] = true
	return copied
}

func withoutKey(m map[string]any, key string) map[string]any {
	copied := make(map[string]any, len(m))
	for k, v := range m {
		if k != key {
			copied[k] = v
		}
	}
	return copied
}

func contains(list []string, s string) bool {
	for _, each := range list {
		if each == s {
			return true
		}
	}
	return false
}

func checkEffect(effect, label, source, current string, warn func(string, ...any)) string {
	switch effect {
	case effectRead, effectWrite, effectDestructive:
		return effect
	}
	warn("%s gives %s the effect %q; expected read, write or destructive", source, label, effect)
	return current
}

var (
	acronymBoundary = regexp.MustCompile(`([A-Z]+)([A-Z][a-z])`)
	wordBoundary    = regexp.MustCompile(`([a-z0-9])([A-Z])`)
)

// GetDevicesDfuStatus → get_devices_dfu_status; OAuth2ClientCredentials → o_auth2_client_credentials.
func snakeCase(id string) string {
	s := acronymBoundary.ReplaceAllString(id, "${1}_${2}")
	s = wordBoundary.ReplaceAllString(s, "${1}_${2}")
	return strings.ToLower(s)
}

// Tags in the order the specification lists them, then any others alphabetically, then untagged.
func tagOrder(spec map[string]any, groups map[string][]operation) []string {
	var order []string
	listed := map[string]bool{}
	if tags, ok := spec["tags"].([]any); ok {
		for _, t := range tags {
			tag, _ := t.(map[string]any)
			name := fmt.Sprint(tag["name"])
			if _, used := groups[name]; used && !listed[name] {
				order = append(order, name)
				listed[name] = true
			}
		}
	}
	for _, name := range sortedKeys(groups) {
		if !listed[name] && name != "" {
			order = append(order, name)
		}
	}
	if _, ok := groups[""]; ok {
		order = append(order, "")
	}
	return order
}

func writeHeader(b *strings.Builder, title, version, specURL string, count int, unlisted []string, linked bool) {
	if title == "" {
		title = "API"
	}
	spec := "the OpenAPI specification"
	if specURL != "" {
		spec = "[the OpenAPI specification](" + specURL + ")"
	}
	intro := fmt.Sprintf("The %s in the paths of %s (%s %s), grouped by tag.", plural(count, "operation"), spec, title, version)
	if len(unlisted) > 0 {
		intro += " Not listed: " + joinAnd(unlisted) + "; read them in the specification. Apart from" +
			" those, this is a complete index generated from the specification, not a summary of it."
	} else {
		intro += " This is a complete index generated from the specification, not a summary of it."
	}
	if linked {
		intro += " Each operationId links to that operation's complete contract, taken unchanged from the" +
			" specification: its parameters, request body and responses, and every component they" +
			" reference, including what the fields of a response mean."
	} else {
		intro += " For an operation's full contract, including its request and response schemas and what" +
			" the fields of a response mean, read that operation in the specification."
	}
	intro += " If this index and the specification ever disagree, the specification is right."

	fmt.Fprintf(b, "# %s operations\n\n", title)
	b.WriteString(generatedComment + "\n")
	b.WriteString(wrap(intro, 90) + "\n\n")
	b.WriteString("Each entry gives:\n\n")
	b.WriteString("- The operationId, its call name (the operationId in snake_case), its effect, its method\n")
	b.WriteString("  and path, and any tags beyond the one it is grouped under.\n")
	b.WriteString("- Its summary, and a note when its effect needs one.\n")
	b.WriteString("- Its parameters by location: `path`, `query`, `querystring`, `header`, `cookie` and\n")
	b.WriteString("  `body`. `*` marks a required parameter, `[]` an array, and `: a|b` the allowed values,\n")
	b.WriteString("  written as JSON when any of them isn't a plain word. `body` lists the top-level\n")
	b.WriteString("  properties of a JSON request body, or gives the media type of any other body. A body\n")
	b.WriteString("  the request may omit is marked optional, and `*` then marks the properties required\n")
	b.WriteString("  when a body is sent; a body with other media types or alternative schemas says so.\n")
	b.WriteString("  Deprecated parameters, read-only body properties, and `Accept`, `Content-Type` and\n")
	b.WriteString("  `Authorization` headers are left out.\n\n")
	b.WriteString("An effect is `read`, `write` or `destructive`. It follows the HTTP method (GET, HEAD,\n")
	b.WriteString("OPTIONS, TRACE and QUERY read; POST, PUT and PATCH write; DELETE is destructive) unless\n")
	b.WriteString("the entry says otherwise.\n")
}

func writeOperation(b *strings.Builder, op operation) {
	id := "`" + op.id + "`"
	if op.link != "" {
		id = "[" + id + "](" + op.link + ")"
	}
	if op.call != "" {
		id += " `" + op.call + "`"
	}
	suffix := ""
	if op.deprecated {
		suffix += " (deprecated)"
	}
	if len(op.otherTags) > 0 {
		suffix += " (also tagged " + strings.Join(op.otherTags, ", ") + ")"
	}
	fmt.Fprintf(b, "- %s %s `%s %s`%s\n", id, op.effect, op.method, op.path, suffix)
	line := op.summary
	if op.note != "" {
		if line != "" && !strings.ContainsAny(line[len(line)-1:], ".!?") {
			line += "."
		}
		line = strings.TrimSpace(line + " Note: " + op.note)
	}
	if line != "" {
		fmt.Fprintf(b, "  %s\n", line)
	}
	if len(op.parameters) > 0 {
		fmt.Fprintf(b, "  %s\n", strings.Join(op.parameters, "; "))
	}
}

// "1 operation", "2 operations"; an irregular plural can be given.
func plural(n int, singular string, pluralForm ...string) string {
	if n == 1 {
		return "1 " + singular
	}
	if len(pluralForm) > 0 {
		return fmt.Sprintf("%d %s", n, pluralForm[0])
	}
	return fmt.Sprintf("%d %ss", n, singular)
}

// "a", "a and b", "a, b and c".
func joinAnd(items []string) string {
	if len(items) < 2 {
		return strings.Join(items, "")
	}
	return strings.Join(items[:len(items)-1], ", ") + " and " + items[len(items)-1]
}

// Wrap text at width columns, breaking only at spaces.
func wrap(text string, width int) string {
	var lines []string
	line := ""
	for _, word := range strings.Fields(text) {
		if line != "" && len(line)+1+len(word) > width {
			lines = append(lines, line)
			line = word
		} else if line == "" {
			line = word
		} else {
			line += " " + word
		}
	}
	if line != "" {
		lines = append(lines, line)
	}
	return strings.Join(lines, "\n")
}

func oneLine(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

func sortedKeys[V any](m map[string]V) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
