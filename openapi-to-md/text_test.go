package main

import (
	"reflect"
	"strings"
	"testing"
)

func TestToolName(t *testing.T) {
	for in, want := range map[string]string{
		"ListPets":            "list_pets",
		"GetPetByID":          "get_pet_by_id",
		"GetHTTPServerStatus": "get_http_server_status",
		"OAuth2Token":         "o_auth2_token",
		"getV2Users":          "get_v2_users",
		"list_orders_v1":      "list_orders_v1",
		"Fetch3DModel":        "fetch3_d_model",
		"S3BucketInfo":        "s3_bucket_info",
		"UID":                 "uid",
		"a":                   "a",
		"":                    "",
		"ÜberCafé":            "über_café",
	} {
		if got := toolName(in); got != want {
			t.Errorf("toolName(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestSentences(t *testing.T) {
	for _, c := range []struct {
		in   string
		want []string
	}{
		{"One. Two.", []string{"One.", "Two."}},
		{"Use e.g. `a.b` or i.e. this. Next", []string{"Use e.g. `a.b` or i.e. this.", "Next"}},
		{"Version 1.5 is out! Really?", []string{"Version 1.5 is out!", "Really?"}},
		{"Wrapped\nline. Second", []string{"Wrapped line.", "Second"}},
		{"Title\n\nBody text.", []string{"Title", "Body text."}},
		{"Pick one:\n- `a`: first\n- `b`: second", []string{"Pick one: `a`: first; `b`: second"}},
		{"Intro.\n- x\n- y", []string{"Intro.", "x; y"}},
		{"## Heading\nText with #hash.", []string{"Text with #hash."}},
		{"# of items per page.", []string{"# of items per page."}},
		{"Use `[a-z](x)` or [docs](https://example.com) ![logo](l.png) here.", []string{"Use `[a-z](x)` or docs logo here."}},
		{"See [the docs](https://example.com/x) now.", []string{"See the docs now."}},
		{"Code:\n```\nnot. a. sentence.\n```\nAfter.", []string{"Code:", "After."}},
		{`He said "stop." Then left.`, []string{`He said "stop."`, "Then left."}},
		{"", nil},
	} {
		if got := sentences(c.in); !reflect.DeepEqual(got, c.want) {
			t.Errorf("sentences(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestIsEcho(t *testing.T) {
	for _, c := range []struct {
		name, text string
		want       bool
	}{
		{"ownerUID", "An OwnerUID.", true},
		{"groupUID", "Filter by GroupUID", true},
		{"label", "Label filter.", true},
		{"phoneNumber", "Phone number filter.", true},
		{"reportId", "The report's ID.", true},
		{"list_pets", "List pets", true},
		{"delete_pet", "Delete a pet", true},
		{"list_pets", "List pets in a store", false},
		{"limit", "Limit the number of results returned", false},
		{"deleted", "If true, also return deleted items.", false},
		{"orderBy", "Field to order by", false},
	} {
		if got := isEcho(c.text, nameWords(c.name)); got != c.want {
			t.Errorf("isEcho(%q, %q) = %v, want %v", c.text, c.name, got, c.want)
		}
	}
}

func TestOpeningList(t *testing.T) {
	for _, c := range []struct {
		in, lead string
		want     []string
	}{
		{"Kind.\n- `a`: first\n- `b`, `c`: others", "Kind.", []string{"a", "b", "c"}},
		{"- `x` (`obj`); `y` (`obj2`)\n- `z` or `w`", "", []string{"x", "y", "z", "w"}},
		{"Kind.\n- `a`: first\n- plain item", "Kind.", nil},
		{"- `a`: wrapped\n  continuation\n- `b`", "", []string{"a", "b"}},
		{"Sort order.\n\n- `asc`: oldest first\n- `desc`: newest first", "Sort order.", []string{"asc", "desc"}},
		{"URL to call.\n\nErrors:\n- `400`: not https", "", nil},
		{"Kind.\n- plain\n\n- `p`\n- `q`", "Kind.", nil},
		{"Display settings. Recognized keys:\n- `theme`: light or dark", "Display settings. Recognized keys:", nil},
		{"How often. For example:\n- `10m`\n- `1h`", "How often. For example:", []string{"10m", "1h"}},
		{"No list here, `a` and `b`.", "", nil},
	} {
		lead, got := openingList(c.in)
		if !reflect.DeepEqual(got, c.want) || got != nil && lead != c.lead {
			t.Errorf("openingList(%q) = %q, %q; want %q, %q", c.in, lead, got, c.lead, c.want)
		}
	}
}

func TestTrim(t *testing.T) {
	long := strings.Repeat("word ", 40)
	if got := trim(long, 120); len([]rune(got)) > 120 || !strings.HasSuffix(got, "…") || strings.Contains(got, " …") {
		t.Errorf("trim(long) = %q", got)
	}
	if got := trim("short", 120); got != "short" {
		t.Errorf("trim(short) = %q", got)
	}
	code := strings.Repeat("x", 100) + " `some code span that is long` tail"
	if got := trim(code, 120); strings.Count(got, "`")%2 != 0 {
		t.Errorf("trim left a code span open: %q", got)
	}
}

func TestLiteral(t *testing.T) {
	for _, c := range []struct {
		in   any
		want string
	}{
		{"plain", "plain"},
		{"two words", `"two words"`},
		{"a|b", `"a|b"`},
		{"", `""`},
		{float64(1000), "1000"},
		{1.5, "1.5"},
		{true, "true"},
		{nil, "null"},
		{[]any{"a", "b"}, `["a","b"]`},
	} {
		if got := literal(c.in); got != c.want {
			t.Errorf("literal(%#v) = %s, want %s", c.in, got, c.want)
		}
	}
}

func TestDescribeRewritesOperationIDs(t *testing.T) {
	r := &renderer{names: map[string]string{"ListThings": "list_things", "GetThing": "fetch_one"}}
	got := r.describe("As listed by `ListThings` or GetThing, not GetThings.", nil)
	want := "As listed by `list_things` or fetch_one, not GetThings."
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestParseOverrides(t *testing.T) {
	ov, err := parseOverrides([]byte(`
OpA:
  x-mcp-class: destructive
  x-mcp-hidden: true
  x-mcp-name: op_a
  parameters:
    query.p: {x-mcp-class: write, x-mcp-name: pp}
OpB:
`))
	if err != nil {
		t.Fatal(err)
	}
	a := ov["OpA"]
	if a == nil || *a.Class != classDestructive || !*a.Hidden || *a.Name != "op_a" {
		t.Fatalf("OpA = %+v", a)
	}
	if p := a.params["query.p"]; p == nil || *p.Class != classWrite || *p.Name != "pp" {
		t.Fatalf("OpA query.p = %+v", p)
	}
	for _, bad := range []string{
		"OpA: {x-mcp-class: nope}",
		"OpA: {x-mcp-name: ''}",
		"OpA: {parameters: {p: {x-mcp-hidden: true}}}",
		"OpA: {typo: 1}",
		"- not a map",
	} {
		if _, err := parseOverrides([]byte(bad)); err == nil {
			t.Errorf("parseOverrides(%q) succeeded", bad)
		}
	}
	if ov, err := parseOverrides(nil); err != nil || len(ov) != 0 {
		t.Errorf("empty overrides: %v, %v", ov, err)
	}
}
