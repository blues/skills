package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/fs"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"testing"
)

var update = flag.Bool("update", false, "rewrite the want.* files under testdata/cases")

// TestCases runs the CLI on every directory under testdata/cases. Each holds
// an args file (CLI arguments, one per line, relative to the directory) and
// the expected results: want.md (stdout), want.err (stderr, if any) and
// want.code (exit status, if not 0).
func TestCases(t *testing.T) {
	dirs, err := filepath.Glob("testdata/cases/*")
	if err != nil {
		t.Fatal(err)
	}
	if len(dirs) == 0 {
		t.Fatal("no cases")
	}
	for _, dir := range dirs {
		t.Run(filepath.Base(dir), func(t *testing.T) {
			args := readArgs(t, dir)
			t.Chdir(dir)
			stdout, stderr, code := runCLI(args, "")
			golden(t, "want.md", stdout)
			golden(t, "want.err", stderr)
			c := ""
			if code != 0 {
				c = strconv.Itoa(code) + "\n"
			}
			golden(t, "want.code", c)
		})
	}
}

// TestDeterminism renders every successful case several times; Go randomizes
// map iteration, so any order leak shows up as a difference.
func TestDeterminism(t *testing.T) {
	dirs, _ := filepath.Glob("testdata/cases/*")
	for _, dir := range dirs {
		t.Run(filepath.Base(dir), func(t *testing.T) {
			args := readArgs(t, dir)
			t.Chdir(dir)
			first, _, code := runCLI(args, "")
			if code != 0 {
				t.Skip("case fails by design")
			}
			for range 5 {
				if again, _, _ := runCLI(args, ""); again != first {
					t.Fatal("output differs between runs")
				}
			}
		})
	}
}

func TestStdinAndOutputFile(t *testing.T) {
	spec, err := os.ReadFile("testdata/cases/refs/spec.yaml")
	if err != nil {
		t.Fatal(err)
	}
	want, err := os.ReadFile("testdata/cases/refs/want.md")
	if err != nil {
		t.Fatal(err)
	}
	stdout, stderr, code := runCLI([]string{"-"}, string(spec))
	if code != 0 || stdout != string(want) {
		t.Fatalf("stdin: code %d, stderr %q, output differs: %v", code, stderr, stdout != string(want))
	}
	out := filepath.Join(t.TempDir(), "out.md")
	// Flags may also follow the spec argument.
	stdout, stderr, code = runCLI([]string{"testdata/cases/refs/spec.yaml", "-o", out}, "")
	if code != 0 || stdout != "" || stderr != "" {
		t.Fatalf("-o: code %d, stdout %q, stderr %q", code, stdout, stderr)
	}
	if got, _ := os.ReadFile(out); !bytes.Equal(got, want) {
		t.Fatal("-o file differs from want.md")
	}
}

func TestURL(t *testing.T) {
	dir := "testdata/cases/external-refs"
	want, err := os.ReadFile(filepath.Join(dir, "want.md"))
	if err != nil {
		t.Fatal(err)
	}
	srv := httptest.NewServer(http.FileServer(http.Dir(dir)))
	defer srv.Close()
	stdout, stderr, code := runCLI([]string{srv.URL + "/spec.yaml"}, "")
	if code != 0 || stdout != string(want) {
		t.Fatalf("code %d, stderr %q, output differs: %v", code, stderr, stdout != string(want))
	}
	if _, stderr, code = runCLI([]string{srv.URL + "/missing.yaml"}, ""); code != 1 || !strings.Contains(stderr, "404") {
		t.Fatalf("missing spec: code %d, stderr %q", code, stderr)
	}
}

func TestRemoteSpecCannotReadLocalFiles(t *testing.T) {
	local, err := filepath.Abs("testdata/cases/external-refs/parts/schemas.yaml")
	if err != nil {
		t.Fatal(err)
	}
	spec := `openapi: 3.0.3
info: {title: t, version: "1"}
paths:
  /x:
    post:
      operationId: PostX
      requestBody:
        content:
          application/json:
            schema: {$ref: 'file://` + filepath.ToSlash(local) + `#/Boat'}
      responses: {'200': {description: ok}}
`
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte(spec))
	}))
	defer srv.Close()
	_, stderr, code := runCLI([]string{srv.URL + "/spec.yaml"}, "")
	if code != 1 || !strings.Contains(stderr, "refusing to read") {
		t.Fatalf("code %d, stderr %q", code, stderr)
	}
}

// A remote document reached from a local spec may not read local files,
// whether it names them by absolute path or file:// URL.
func TestRemoteRefCannotReadLocalFiles(t *testing.T) {
	secret, err := filepath.Abs("testdata/cases/external-refs/parts/schemas.yaml")
	if err != nil {
		t.Fatal(err)
	}
	for _, ref := range []string{filepath.ToSlash(secret), "file://" + filepath.ToSlash(secret)} {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Write([]byte("type: object\nproperties:\n  inner: {$ref: '" + ref + "#/Boat'}\n"))
		}))
		spec := filepath.Join(t.TempDir(), "spec.yaml")
		os.WriteFile(spec, []byte(`openapi: 3.0.3
info: {title: t, version: "1"}
paths:
  /x:
    post:
      operationId: PostX
      requestBody: {content: {application/json: {schema: {$ref: '`+srv.URL+`/remote.yaml'}}}}
      responses: {'200': {description: ok}}
`), 0o644)
		stdout, stderr, code := runCLI([]string{spec}, "")
		srv.Close()
		if code != 1 || !strings.Contains(stderr, "refusing to read") || strings.Contains(stdout, "name") {
			t.Errorf("ref %s: code %d, stderr %q", ref, code, stderr)
		}
	}
}

// Relative refs resolve against the spec's final URL after redirects, and
// don't inherit its query string.
func TestRedirectedSpecResolvesRefs(t *testing.T) {
	dir := "testdata/cases/external-refs"
	files := http.FileServer(http.Dir(dir))
	var queries []string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/latest" {
			http.Redirect(w, r, "/real/spec.yaml?token=secret", http.StatusFound)
			return
		}
		queries = append(queries, r.URL.Path+"?"+r.URL.RawQuery)
		r.URL.Path = strings.TrimPrefix(r.URL.Path, "/real")
		files.ServeHTTP(w, r)
	}))
	defer srv.Close()
	want, _ := os.ReadFile(filepath.Join(dir, "want.md"))
	stdout, stderr, code := runCLI([]string{srv.URL + "/latest"}, "")
	if code != 0 || stdout != string(want) {
		t.Fatalf("code %d, stderr %q", code, stderr)
	}
	for _, q := range queries[1:] {
		if strings.Contains(q, "token") {
			t.Errorf("ref fetch %s inherited the spec's query", q)
		}
	}
}

func TestPublicOnlyClient(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {}))
	defer srv.Close()
	if _, err := publicOnlyClient().Get(srv.URL); err == nil || !strings.Contains(err.Error(), "non-public") {
		t.Fatalf("loopback fetch: %v", err)
	}
	for ip, want := range map[string]bool{"8.8.8.8": true, "127.0.0.1": false, "10.1.2.3": false, "169.254.169.254": false, "::1": false, "192.168.0.1": false} {
		if got := public(net.ParseIP(ip)); got != want {
			t.Errorf("public(%s) = %v", ip, got)
		}
	}
}

// A small spec whose YAML aliases expand exponentially is refused.
func TestAliasBomb(t *testing.T) {
	var b strings.Builder
	b.WriteString("openapi: 3.0.3\ninfo: {title: b, version: '1'}\ncomponents:\n  schemas:\n    L0: &a0 {type: string}\n")
	for i := 1; i <= 7; i++ {
		var props []string
		for j := range 10 {
			props = append(props, fmt.Sprintf("p%d: *a%d", j, i-1))
		}
		fmt.Fprintf(&b, "    L%d: &a%d {type: object, properties: {%s}}\n", i, i, strings.Join(props, ", "))
	}
	b.WriteString("paths: {}\n")
	_, stderr, code := runCLI([]string{"-"}, b.String())
	if code != 1 || !strings.Contains(stderr, "aliases") {
		t.Fatalf("code %d, stderr %q", code, stderr)
	}
}

func TestUsage(t *testing.T) {
	for _, args := range [][]string{nil, {"a.yaml", "b.yaml"}, {"-nope", "a.yaml"}} {
		if _, stderr, code := runCLI(args, ""); code != 2 || !strings.Contains(stderr, "usage") {
			t.Errorf("%q: code %d, stderr %q", args, code, stderr)
		}
	}
	if _, stderr, code := runCLI([]string{"-h"}, ""); code != 0 || !strings.Contains(stderr, "usage") {
		t.Errorf("-h: code %d, stderr %q", code, stderr)
	}
	if _, stderr, code := runCLI([]string{"testdata/no-such-file.yaml"}, ""); code != 1 || stderr == "" {
		t.Errorf("missing file: code %d, stderr %q", code, stderr)
	}
}

// TestHygiene enforces that no product-specific names reach the Go source:
// testdata/hygiene.txt lists forbidden words ("word x", any case) and specs
// ("spec path") whose operationIds, their tool names, and their distinctive
// parameter names (with an uppercase letter or a hyphen) are forbidden.
func TestHygiene(t *testing.T) {
	data, err := os.ReadFile("testdata/hygiene.txt")
	if err != nil {
		t.Fatal(err)
	}
	var words []string
	names := map[string]bool{}
	for _, line := range strings.Split(string(data), "\n") {
		kind, arg, _ := strings.Cut(strings.TrimSpace(line), " ")
		switch kind {
		case "word":
			words = append(words, strings.ToLower(arg))
		case "spec":
			src, err := loadSpec(filepath.Join("testdata", arg), nil)
			if err != nil {
				t.Fatal(err)
			}
			for _, item := range src.doc.Paths.Map() {
				for _, op := range item.Operations() {
					names[op.OperationID] = true
					names[toolName(op.OperationID)] = true
					for _, p := range append(item.Parameters, op.Parameters...) {
						if n := p.Value.Name; strings.ContainsAny(n, "ABCDEFGHIJKLMNOPQRSTUVWXYZ-") {
							names[n] = true
						}
					}
				}
			}
		}
	}
	if len(words) == 0 || len(names) == 0 {
		t.Fatal("hygiene.txt lists no words or no specs")
	}
	err = filepath.WalkDir(".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() && path != "." && (d.Name() == "testdata" || strings.HasPrefix(d.Name(), ".")) {
			return filepath.SkipDir
		}
		if d.IsDir() || !strings.HasSuffix(path, ".go") {
			return nil
		}
		src, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		lower := strings.ToLower(string(src))
		for _, w := range words {
			if strings.Contains(lower, w) {
				t.Errorf("%s mentions %q", path, w)
			}
		}
		for _, n := range sortedKeys(names) {
			if regexp.MustCompile(`(^|[^\w-])` + regexp.QuoteMeta(n) + `($|[^\w-])`).Match(src) {
				t.Errorf("%s mentions %q", path, n)
			}
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}

func runCLI(args []string, stdin string) (stdout, stderr string, code int) {
	var out, errOut bytes.Buffer
	code = run(args, strings.NewReader(stdin), &out, &errOut)
	return out.String(), errOut.String(), code
}

func readArgs(t *testing.T, dir string) []string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(dir, "args"))
	if err != nil {
		t.Fatal(err)
	}
	return strings.Fields(string(data))
}

// golden compares got with the named file (a missing file means empty), or
// rewrites the file with -update.
func golden(t *testing.T, name, got string) {
	t.Helper()
	if *update {
		if got == "" {
			os.Remove(name)
			return
		}
		if err := os.WriteFile(name, []byte(got), 0o644); err != nil {
			t.Fatal(err)
		}
		return
	}
	want, err := os.ReadFile(name)
	if err != nil && !os.IsNotExist(err) {
		t.Fatal(err)
	}
	if got != string(want) {
		t.Errorf("%s mismatch (rerun with -update to accept):\n%s", name, diff(string(want), got))
	}
}

// diff shows the first differing lines, enough to see what changed.
func diff(want, got string) string {
	w, g := strings.Split(want, "\n"), strings.Split(got, "\n")
	for i := 0; i < len(w) || i < len(g); i++ {
		var a, b string
		if i < len(w) {
			a = w[i]
		}
		if i < len(g) {
			b = g[i]
		}
		if a != b {
			return "line " + strconv.Itoa(i+1) + ":\n  want: " + a + "\n  got:  " + b
		}
	}
	return "(identical lines; trailing difference)"
}
