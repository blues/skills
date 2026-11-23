package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
	"unicode"

	"gopkg.in/yaml.v3"
)

// The x-mcp-* extension keys. The same keys are used in an overrides file.
const (
	extClass  = "x-mcp-class"
	extHidden = "x-mcp-hidden"
	extName   = "x-mcp-name"
)

// class is the read/write/destructive class of a call, ordered so that a
// larger class is stronger.
type class int

const (
	classRead class = iota
	classWrite
	classDestructive
)

var classNames = [...]string{"read", "write", "destructive"}

func (c class) String() string { return classNames[c] }

func parseClass(s string) (class, error) {
	for i, n := range classNames {
		if s == n {
			return class(i), nil
		}
	}
	return 0, fmt.Errorf("invalid %s %q (want read, write or destructive)", extClass, s)
}

// methodClass is the default class for an HTTP method.
func methodClass(method string) class {
	switch method {
	case "GET", "HEAD", "OPTIONS":
		return classRead
	case "DELETE":
		return classDestructive
	}
	return classWrite
}

// mcpFields holds x-mcp-* values from a spec extension map or an overrides
// file entry. Nil means "not set".
type mcpFields struct {
	Class  *class
	Hidden *bool
	Name   *string
}

// overrides maps an original operationId to its overrides.
type overrides map[string]*opOverride

type opOverride struct {
	mcpFields
	// params is keyed by parameter name, or "in.name" when the name alone is
	// ambiguous (e.g. "query.delete", "body.force").
	params map[string]*mcpFields
}

// Raw shapes of the overrides file; unknown keys are rejected so a typo
// cannot silently drop a safety annotation.
type rawOpOverride struct {
	Class      *string                      `yaml:"x-mcp-class"`
	Hidden     *bool                        `yaml:"x-mcp-hidden"`
	Name       *string                      `yaml:"x-mcp-name"`
	Parameters map[string]*rawParamOverride `yaml:"parameters"`
}

type rawParamOverride struct {
	Class *string `yaml:"x-mcp-class"`
	Name  *string `yaml:"x-mcp-name"`
}

// goTypeRE matches the Go type names yaml.v3 puts in its errors.
var goTypeRE = regexp.MustCompile(` (in type|into) \*?(map\[string\]\*?)?main\.\w+`)

func loadOverrides(path string) (overrides, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	ov, err := parseOverrides(data)
	if err != nil {
		return nil, fmt.Errorf("overrides %s: %w", path, err)
	}
	return ov, nil
}

func parseOverrides(data []byte) (overrides, error) {
	var raw map[string]*rawOpOverride
	dec := yaml.NewDecoder(bytes.NewReader(data))
	dec.KnownFields(true)
	if err := dec.Decode(&raw); err != nil && !errors.Is(err, io.EOF) {
		return nil, errors.New(goTypeRE.ReplaceAllString(err.Error(), ""))
	}
	var extra yaml.Node
	if err := dec.Decode(&extra); !errors.Is(err, io.EOF) {
		return nil, errors.New("more than one YAML document; put every entry in one")
	}
	ov := overrides{}
	for _, id := range sortedKeys(raw) {
		o := &opOverride{params: map[string]*mcpFields{}}
		r := raw[id]
		if r == nil { // an empty entry still names an operation
			ov[id] = o
			continue
		}
		var err error
		if o.mcpFields, err = toFields(r.Class, r.Hidden, r.Name); err != nil {
			return nil, fmt.Errorf("%s: %w", id, err)
		}
		for _, key := range sortedKeys(r.Parameters) {
			f := mcpFields{}
			if p := r.Parameters[key]; p != nil {
				if f, err = toFields(p.Class, nil, p.Name); err != nil {
					return nil, fmt.Errorf("%s: parameter %s: %w", id, key, err)
				}
			}
			o.params[key] = &f
		}
		ov[id] = o
	}
	return ov, nil
}

func toFields(cls *string, hidden *bool, name *string) (mcpFields, error) {
	f := mcpFields{Hidden: hidden}
	if cls != nil {
		c, err := parseClass(*cls)
		if err != nil {
			return f, err
		}
		f.Class = &c
	}
	if name != nil {
		if err := checkName(*name); err != nil {
			return f, err
		}
		f.Name = name
	}
	return f, nil
}

// checkName rejects an x-mcp-name that would break the output's lines.
func checkName(s string) error {
	if s == "" || strings.ContainsFunc(s, func(r rune) bool { return unicode.IsSpace(r) || unicode.IsControl(r) }) {
		return fmt.Errorf("%s %q must be non-empty, without spaces", extName, s)
	}
	return nil
}

// extFields reads the x-mcp-* keys from a spec extension map; x-mcp-hidden
// is allowed only on an operation. Other x-* keys are ignored.
func extFields(ext map[string]any, operation bool) (mcpFields, error) {
	var f mcpFields
	if v, ok := ext[extClass]; ok {
		s, ok := v.(string)
		if !ok {
			return f, fmt.Errorf("%s must be a string, got %v", extClass, v)
		}
		c, err := parseClass(s)
		if err != nil {
			return f, err
		}
		f.Class = &c
	}
	if v, ok := ext[extHidden]; ok {
		if !operation {
			return f, fmt.Errorf("%s applies only to operations", extHidden)
		}
		b, ok := v.(bool)
		if !ok {
			return f, fmt.Errorf("%s must be a boolean, got %v", extHidden, v)
		}
		f.Hidden = &b
	}
	if v, ok := ext[extName]; ok {
		s, ok := v.(string)
		if !ok {
			return f, fmt.Errorf("%s must be a string, got %v", extName, v)
		}
		if err := checkName(s); err != nil {
			return f, err
		}
		f.Name = &s
	}
	return f, nil
}

// merge returns f with every field that is set in g replaced by g's value.
func (f mcpFields) merge(g *mcpFields) mcpFields {
	if g == nil {
		return f
	}
	if g.Class != nil {
		f.Class = g.Class
	}
	if g.Hidden != nil {
		f.Hidden = g.Hidden
	}
	if g.Name != nil {
		f.Name = g.Name
	}
	return f
}
