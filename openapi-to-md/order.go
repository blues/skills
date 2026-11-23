package main

import (
	"sort"

	"github.com/getkin/kin-openapi/openapi3"
	"gopkg.in/yaml.v3"
)

// Spec order comes from source positions: the loader records where each
// path, operation and property key appears.

type position struct{ line, col int }

func keyPosition(o *openapi3.Origin) (position, bool) {
	if o == nil || o.Key == nil || o.Key.Line == 0 {
		return position{}, false
	}
	return position{o.Key.Line, o.Key.Column}, true
}

// orderedByPosition sorts keys by the source position of their values;
// values without one go last, in tieBreak order.
func orderedByPosition(keys []string, pos func(string) (position, bool), tieBreak func(a, b string) bool) []string {
	sort.SliceStable(keys, func(i, j int) bool {
		pi, oki := pos(keys[i])
		pj, okj := pos(keys[j])
		switch {
		case oki && okj && pi != pj:
			if pi.line != pj.line {
				return pi.line < pj.line
			}
			return pi.col < pj.col
		case oki != okj:
			return oki
		}
		return tieBreak(keys[i], keys[j])
	})
	return keys
}

func byName(a, b string) bool { return a < b }

// orderedPaths lists the paths in the order the root document gives them.
// Positions come from the raw document, because a $ref'd path item's
// recorded position is in the file that defines it.
func (r *renderer) orderedPaths() []string {
	if r.doc.Paths == nil {
		return nil
	}
	order := map[string]int{}
	if root := r.raw; root != nil {
		if root.Kind == yaml.DocumentNode && len(root.Content) == 1 {
			root = root.Content[0]
		}
		if paths := mapValue(root, "paths"); paths != nil {
			for i := 0; i+1 < len(paths.Content); i += 2 {
				order[paths.Content[i].Value] = i/2 + 1
			}
		}
	}
	m := r.doc.Paths.Map()
	return orderedByPosition(sortedKeys(m), func(k string) (position, bool) {
		if i, ok := order[k]; ok {
			return position{line: i}, true
		}
		return keyPosition(m[k].Origin)
	}, byName)
}

type methodOp struct {
	method string
	op     *openapi3.Operation
}

var methodOrder = map[string]int{"GET": 0, "PUT": 1, "POST": 2, "DELETE": 3, "OPTIONS": 4, "HEAD": 5, "PATCH": 6, "TRACE": 7}

func orderedMethods(item *openapi3.PathItem) []methodOp {
	if item == nil {
		return nil
	}
	ops := item.Operations()
	keys := orderedByPosition(sortedKeys(ops), func(k string) (position, bool) {
		return keyPosition(ops[k].Origin)
	}, func(a, b string) bool {
		ra, oka := methodOrder[a]
		rb, okb := methodOrder[b]
		if oka && okb {
			return ra < rb
		}
		if oka != okb {
			return oka
		}
		return a < b
	})
	out := make([]methodOp, 0, len(keys))
	for _, k := range keys {
		out = append(out, methodOp{k, ops[k]})
	}
	return out
}

func orderedProperties(props openapi3.Schemas) []string {
	return orderedByPosition(sortedKeys(props), func(k string) (position, bool) {
		if p := props[k]; p != nil {
			return keyPosition(p.Origin)
		}
		return position{}, false
	}, byName)
}

func sortedKeys[V any](m map[string]V) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func appendUnique(list []string, items ...string) []string {
	for _, s := range items {
		dup := false
		for _, x := range list {
			if x == s {
				dup = true
				break
			}
		}
		if !dup {
			list = append(list, s)
		}
	}
	return list
}
