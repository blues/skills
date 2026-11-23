package main

import (
	"fmt"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

// options control rendering.
type options struct {
	includeHidden bool
	// commonMin is how many operations must share a param, with the same
	// spec, for it to be defined once under "common params"; 0 never does.
	commonMin int
}

// render produces the Markdown reference. Warnings are returned even when
// err is non-nil.
func render(src *source, ov overrides, opts options) ([]byte, []string, error) {
	r := newRenderer(src, ov, opts)
	ops, err := r.operations()
	if err != nil {
		return nil, r.warnings, err
	}
	common := pickCommon(ops, opts.commonMin)
	types := r.usedTypes(ops)

	var b strings.Builder
	writeHeader(&b, src.doc, len(common) > 0, len(types) > 0)
	writeCommon(&b, ops, common)
	for _, g := range groupByTag(ops) {
		heading := g.tag
		switch strings.ToLower(heading) {
		case "common params", "types": // keep the reserved section names unique
			heading = "tag: " + heading
		}
		fmt.Fprintf(&b, "\n## %s\n", heading)
		for _, o := range g.ops {
			b.WriteByte('\n')
			writeOperation(&b, o, common)
		}
	}
	if len(types) > 0 {
		b.WriteString("\n## types\n")
		for _, t := range types {
			b.WriteByte('\n')
			b.WriteString(t.name)
			if t.alias != "" {
				b.WriteString(" = " + t.alias)
			}
			if t.note != "" {
				b.WriteString("  — " + t.note)
			}
			b.WriteByte('\n')
			writeParams(&b, t.fields, classRead, "  ")
		}
	}
	return []byte(b.String()), r.warnings, nil
}

func writeHeader(b *strings.Builder, doc *openapi3.T, common, types bool) {
	title, version := "", ""
	if doc.Info != nil {
		title, version = oneLine(doc.Info.Title), oneLine(doc.Info.Version)
	}
	fmt.Fprintf(b, "# MCP execute reference · %s %s · OpenAPI %s\n", title, version, doc.OpenAPI)
	b.WriteString("# call_tool('<name>', {param: value, ...})  ·  *required  ·  header keys use the OpenAPI name exactly\n")
	b.WriteString("# !write / !destructive on a param escalates the call when that param is used (a boolean: when true)\n")
	b.WriteString("# each {x} in a path is a required param x, a string unless an `x  path` line defines it  ·  `a, b  query  string` = params with the same spec\n")
	if common {
		b.WriteString("# `+ a, *b` = params defined under ## common params; a param line under a call applies to that call only\n")
	}
	if types {
		b.WriteString("# a capitalized type such as Foo is defined under ## types; an inline object's fields are indented below it\n")
	}
}

// Common params must pay for themselves: each definition must save at least
// minDefSaving characters over repeating its line, and all of them together
// minSectionSaving, or the section (and its legend line) is left out.
const (
	minDefSaving     = 40
	minSectionSaving = 200
)

// pickCommon chooses, for each key that at least min operations use with
// the same spec, that spec as the key's common definition; it returns
// key -> sig. Escalating and aliased params always stay with their call.
func pickCommon(ops []*operation, min int) map[string]string {
	common := map[string]string{}
	if min <= 0 {
		return common
	}
	type use struct{ count, first int }
	uses := map[string]map[string]*use{} // key -> sig -> use
	n := 0
	for _, o := range ops {
		for _, p := range o.params {
			if p.impliedBy(o) || p.escalates() || p.key != p.name {
				continue
			}
			sig := p.sig()
			if uses[p.key] == nil {
				uses[p.key] = map[string]*use{}
			}
			u := uses[p.key][sig]
			if u == nil {
				u = &use{first: n}
				n++
				uses[p.key][sig] = u
			}
			u.count++
		}
	}
	total := 0
	for key, sigs := range uses {
		var best string
		var bu *use
		for sig, u := range sigs {
			if bu == nil || u.count > bu.count || u.count == bu.count && u.first < bu.first {
				best, bu = sig, u
			}
		}
		// Repeating the line costs about len(sig) per extra use; a
		// reference costs about the key per use.
		saving := (bu.count-1)*len(best) - bu.count*(len(key)+2)
		if bu.count >= max(min, 2) && saving >= minDefSaving {
			common[key] = best
			total += saving
		}
	}
	if total < minSectionSaving {
		return map[string]string{}
	}
	return common
}

// isCommon reports whether p is the common definition of its key.
func (p *param) isCommon(common map[string]string) bool {
	sig, ok := common[p.key]
	return ok && !p.escalates() && p.key == p.name && sig == p.sig()
}

// escalates reports whether p, or any field nested in it, raises a call's
// class: such params always stay with their call.
func (p *param) escalates() bool {
	if p.class > classRead {
		return true
	}
	for _, f := range p.fields {
		if f.escalates() {
			return true
		}
	}
	return false
}

// shadowedPath reports whether common params define p's key as a different
// path param, which the agent would apply to p if p had no line of its own.
func (p *param) shadowedPath(common map[string]string) bool {
	sig, ok := common[p.key]
	return ok && strings.HasPrefix(sig, "path  ") && sig != p.sig()
}

// writeCommon defines the common params, in order of first use.
func writeCommon(b *strings.Builder, ops []*operation, common map[string]string) {
	if len(common) == 0 {
		return
	}
	var defs []*param
	done := map[string]bool{}
	for _, o := range ops {
		for _, p := range o.params {
			if done[p.key] || p.impliedBy(o) || !p.isCommon(common) {
				continue
			}
			done[p.key] = true
			q := *p
			q.required = false // requiredness belongs to each call
			defs = append(defs, &q)
		}
	}
	b.WriteString("\n## common params\n")
	writeParams(b, defs, classRead, "")
}

func writeOperation(b *strings.Builder, o *operation, common map[string]string) {
	fmt.Fprintf(b, "%s  %s  %s %s", o.name, o.class, o.method, o.callPath())
	if o.deprecated {
		b.WriteString("  deprecated")
	}
	if o.hidden {
		b.WriteString("  hidden")
	}
	b.WriteByte('\n')
	if o.summary != "" {
		fmt.Fprintf(b, "  %s\n", o.summary)
	}
	var refs []string
	var own []*param
	for _, p := range o.params {
		switch {
		case p.impliedBy(o) && !p.shadowedPath(common):
		case p.isCommon(common):
			if p.in == "path" && o.hasPathParam(p.name) {
				continue // the path shows it and common params define it
			}
			refs = append(refs, p.label())
		default:
			own = append(own, p)
		}
	}
	if len(refs) > 0 {
		fmt.Fprintf(b, "  + %s\n", strings.Join(refs, ", "))
	}
	writeParams(b, own, o.class, "  ")
}

// writeParams writes one line per param; params whose lines would be the
// same apart from the key share a line. An inline object's fields follow it,
// indented.
func writeParams(b *strings.Builder, params []*param, opClass class, indent string) {
	type line struct {
		keys  []string
		tail  string
		owner *param // set when the line has fields to expand
	}
	var lines []*line
	byTail := map[string]*line{}
	for _, p := range params {
		tail := p.tail(opClass)
		if len(p.fields) > 0 {
			lines = append(lines, &line{keys: []string{p.label()}, tail: tail, owner: p})
			continue
		}
		if l := byTail[tail]; l != nil {
			l.keys = append(l.keys, p.label())
			continue
		}
		l := &line{keys: []string{p.label()}, tail: tail}
		byTail[tail] = l
		lines = append(lines, l)
	}
	for _, l := range lines {
		fmt.Fprintf(b, "%s%s  %s\n", indent, strings.Join(l.keys, ", "), l.tail)
		if l.owner != nil {
			writeParams(b, l.owner.fields, opClass, indent+"  ")
		}
	}
}

// label is the key as listed: "*key" when required.
func (p *param) label() string {
	if p.required {
		return "*" + p.key
	}
	return p.key
}

// impliedBy reports whether o's path template alone says everything about
// p: a required string param named as in {name}.
func (p *param) impliedBy(o *operation) bool {
	return p.in == "path" && p.typ == "string" && p.key == p.name && p.def == "" &&
		p.desc == "" && !p.deprecated && p.class <= o.class && o.hasPathParam(p.name)
}

func (o *operation) hasPathParam(name string) bool {
	return strings.Contains(o.path, "{"+name+"}")
}

// callPath is the path with each aliased path param's placeholder showing
// its call key, so that {x} always names the key to send.
func (o *operation) callPath() string {
	p := o.path
	for _, q := range o.params {
		if q.in == "path" && q.key != q.name {
			p = strings.ReplaceAll(p, "{"+q.name+"}", "{"+q.key+"}")
		}
	}
	return p
}

// spec renders p without its key, requiredness or escalation marker.
func (p *param) spec() string {
	var b strings.Builder
	if p.in != "" {
		b.WriteString(p.in)
		b.WriteString("  ")
	}
	b.WriteString(p.typ)
	if p.def != "" {
		b.WriteString("  default=")
		b.WriteString(p.def)
	}
	if p.deprecated {
		b.WriteString("  deprecated")
	}
	if p.desc != "" {
		b.WriteString("  — ")
		b.WriteString(p.desc)
	}
	return b.String()
}

// sig is spec plus an inline object's fields: params with the same key and
// sig mean the same thing.
func (p *param) sig() string {
	s := p.spec()
	for _, f := range p.fields {
		s += "\n\t" + f.label() + "  " + f.tail(classRead) + strings.ReplaceAll(strings.TrimPrefix(f.sig(), f.spec()), "\n", "\n\t")
	}
	return s
}

// sig is a type's whole content, to recognize identical definitions.
func (t *typeDef) sig() string {
	s := t.alias + "\x00" + t.note
	for _, f := range t.fields {
		s += "\n" + f.label() + "  " + f.tail(classRead) + f.sig()
	}
	return s
}

// tail is the line after the key: spec, with the escalation marker before
// the description when p raises opClass.
func (p *param) tail(opClass class) string {
	s := p.spec()
	if p.class <= opClass {
		return s
	}
	marker := "  !" + p.class.String()
	if i := strings.Index(s, "  — "); i >= 0 {
		return s[:i] + marker + s[i:]
	}
	return s + marker
}

// usedTypes returns the named types that the operations use, directly or
// through other types, in order of first use.
func (r *renderer) usedTypes(ops []*operation) []*typeDef {
	var out []*typeDef
	seen := map[string]bool{}
	add := func(names []string) {
		for _, n := range names {
			if t := r.typeDefs[n]; t != nil && !seen[n] {
				seen[n] = true
				out = append(out, t)
			}
		}
	}
	for _, o := range ops {
		for _, p := range o.params {
			add(p.refs)
		}
	}
	for i := 0; i < len(out); i++ {
		add(out[i].refs)
	}
	return out
}

type tagGroup struct {
	tag string
	ops []*operation
}

// groupByTag groups operations under their first tag, in order of first use.
func groupByTag(ops []*operation) []tagGroup {
	index := map[string]int{}
	var out []tagGroup
	for _, o := range ops {
		i, ok := index[o.tag]
		if !ok {
			i = len(out)
			index[o.tag] = i
			out = append(out, tagGroup{tag: o.tag})
		}
		out[i].ops = append(out[i].ops, o)
	}
	return out
}
