// Command openapi-to-md turns an OpenAPI 3.0/3.1 spec into a compact Markdown
// reference for agents that call an MCP server's code-mode execute tool as
// call_tool('<name>', {params}). See DESIGN.md and README.md.
package main

import (
	"bytes"
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/getkin/kin-openapi/openapi3"
	"gopkg.in/yaml.v3"
)

const maxFetchBytes = 64 << 20

func main() {
	os.Exit(run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr))
}

// run is main without the process: it returns the exit code (0 ok, 1 error,
// 2 usage) so tests can drive the whole CLI.
func run(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("openapi-to-md", flag.ContinueOnError)
	fs.SetOutput(stderr)
	overridesPath := fs.String("overrides", "", "optional YAML of x-mcp-* fields keyed by operationId")
	outPath := fs.String("o", "", "output file (default or - : stdout)")
	includeHidden := fs.Bool("include-hidden", false, "also emit hidden ops, marked hidden")
	commonMin := fs.Int("common", 2, "define a param once, under common params, when this many operations share it (0: never)")
	fs.Usage = func() {
		fmt.Fprintln(stderr, "usage: openapi-to-md [flags] <spec>")
		fmt.Fprintln(stderr, "  <spec>  path, http(s) URL, or - (stdin)")
		fs.PrintDefaults()
	}
	// Flags may come before or after the spec; everything after "--" is
	// positional.
	var rest []string
	for i, a := range args {
		if a == "--" {
			args, rest = args[:i], args[i+1:]
			break
		}
	}
	var positional []string
	for {
		if err := fs.Parse(args); err != nil {
			if errors.Is(err, flag.ErrHelp) {
				return 0
			}
			return 2
		}
		if fs.NArg() == 0 {
			break
		}
		positional = append(positional, fs.Arg(0))
		args = fs.Args()[1:]
	}
	positional = append(positional, rest...)
	if len(positional) != 1 {
		fs.Usage()
		return 2
	}

	fail := func(err error) int {
		fmt.Fprintf(stderr, "openapi-to-md: %v\n", err)
		return 1
	}
	src, err := loadSpec(positional[0], stdin)
	if err != nil {
		return fail(err)
	}
	var ov overrides
	if *overridesPath != "" {
		if ov, err = loadOverrides(*overridesPath); err != nil {
			return fail(err)
		}
	}
	out, warnings, err := render(src, ov, options{includeHidden: *includeHidden, commonMin: *commonMin})
	for _, w := range warnings {
		fmt.Fprintf(stderr, "openapi-to-md: warning: %s\n", w)
	}
	if err != nil {
		return fail(err)
	}
	if *outPath == "" || *outPath == "-" {
		_, err = stdout.Write(out)
	} else {
		err = os.WriteFile(*outPath, out, 0o644)
	}
	if err != nil {
		return fail(err)
	}
	return 0
}

// source is a loaded spec plus its root document as written.
type source struct {
	doc *openapi3.T
	raw *yaml.Node
}

// loadSpec loads and resolves a spec from a file path, an http(s) URL, or "-"
// (stdin). $refs may name http(s) documents, or local files by relative path
// from a local document; a remote document never reads local files.
func loadSpec(arg string, stdin io.Reader) (*source, error) {
	remote := strings.HasPrefix(arg, "http://") || strings.HasPrefix(arg, "https://")
	refs := &refPolicy{remote: remote, local: map[string]bool{}, client: &http.Client{Timeout: 60 * time.Second}}
	var data []byte
	var base *url.URL
	var err error
	switch {
	case arg == "-":
		data, err = readAll(stdin)
	case remote:
		var u *url.URL
		if u, err = url.Parse(arg); err == nil {
			refs.client = fetchClient(u.Hostname())
			data, base, err = fetch(refs.client, arg)
		}
	default:
		data, err = os.ReadFile(arg)
		base = &url.URL{Path: filepath.ToSlash(arg)}
		refs.local[path.Clean(base.Path)] = true
	}
	if err != nil {
		return nil, fmt.Errorf("loading %s: %w", arg, err)
	}
	raw, err := parseYAML(data)
	if err != nil {
		return nil, fmt.Errorf("loading %s: %w", arg, err)
	}
	loader := openapi3.NewLoader()
	loader.IncludeOrigin = true // positions give spec order for operations and properties
	loader.JoinFunc = refs.join
	loader.ReadFromURIFunc = openapi3.URIMapCache(refs.read)
	var doc *openapi3.T
	if base == nil {
		doc, err = loader.LoadFromData(data)
	} else {
		doc, err = loader.LoadFromDataWithPath(data, base)
	}
	if err != nil {
		return nil, fmt.Errorf("loading %s: %w", arg, err)
	}
	if doc.OpenAPI == "" {
		return nil, fmt.Errorf("loading %s: not an OpenAPI 3 document (no openapi field)", arg)
	}
	return &source{doc: doc, raw: raw}, nil
}

// refPolicy decides how $refs resolve and which documents they may read.
type refPolicy struct {
	remote bool            // the root spec came from a URL
	local  map[string]bool // local files reachable by relative refs from local documents
	client *http.Client
}

// join resolves a relative $ref: against a URL by RFC 3986 (so a redirected
// base and its query behave), against a local file by path.
func (p *refPolicy) join(base, rel *url.URL) *url.URL {
	if base != nil && (base.Scheme == "http" || base.Scheme == "https") {
		return base.ResolveReference(rel)
	}
	out := *rel
	if base != nil {
		out = *base
		out.Path = path.Join(path.Dir(base.Path), rel.Path)
	}
	p.local[path.Clean(out.Path)] = true
	return &out
}

func (p *refPolicy) read(_ *openapi3.Loader, u *url.URL) ([]byte, error) {
	var data []byte
	var err error
	switch {
	case u.Scheme == "http" || u.Scheme == "https":
		data, _, err = fetch(p.client, u.String())
	case u.Scheme == "" && !p.remote && p.local[path.Clean(u.Path)]:
		data, err = os.ReadFile(filepath.FromSlash(u.Path))
	case u.Scheme == "" || u.Scheme == "file":
		return nil, fmt.Errorf("refusing to read %s: $refs may name local files only by relative path from a local spec", u.Path)
	default:
		return nil, fmt.Errorf("unsupported $ref %q", u)
	}
	if err == nil {
		_, err = parseYAML(data)
	}
	return data, err
}

func readAll(r io.Reader) ([]byte, error) {
	data, err := io.ReadAll(io.LimitReader(r, maxFetchBytes+1))
	if err == nil && len(data) > maxFetchBytes {
		err = errors.New("input exceeds 64 MiB")
	}
	return data, err
}

// fetch GETs a document and returns it with its final (post-redirect) URL.
func fetch(client *http.Client, u string) ([]byte, *url.URL, error) {
	resp, err := client.Get(u)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode/100 != 2 {
		return nil, nil, fmt.Errorf("GET %s: %s", u, resp.Status)
	}
	data, err := readAll(resp.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("GET %s: %w", u, err)
	}
	return data, resp.Request.URL, nil
}

// fetchClient returns the HTTP client for a spec served from host. A spec on
// a public host may not lead the tool (by $ref or redirect) to loopback,
// private or link-local addresses, so it can't pull internal content into
// the output; a spec served from such a host itself may.
func fetchClient(host string) *http.Client {
	if ips, err := net.LookupIP(host); err == nil {
		for _, ip := range ips {
			if !public(ip) {
				return &http.Client{Timeout: 60 * time.Second}
			}
		}
	}
	return publicOnlyClient()
}

// publicOnlyClient connects only to public addresses, checked as it dials so
// DNS answers can't steer it elsewhere.
func publicOnlyClient() *http.Client {
	dialer := &net.Dialer{Timeout: 30 * time.Second, Control: func(_, address string, _ syscall.RawConn) error {
		h, _, err := net.SplitHostPort(address)
		if err != nil {
			return err
		}
		if ip := net.ParseIP(h); ip == nil || !public(ip) {
			return fmt.Errorf("refusing to connect to non-public address %s", h)
		}
		return nil
	}}
	tr := http.DefaultTransport.(*http.Transport).Clone()
	tr.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
		return dialer.DialContext(ctx, network, addr)
	}
	return &http.Client{Timeout: 60 * time.Second, Transport: tr}
}

func public(ip net.IP) bool {
	return !(ip.IsLoopback() || ip.IsPrivate() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() ||
		ip.IsInterfaceLocalMulticast() || ip.IsUnspecified() || ip.IsMulticast())
}

// parseYAML parses a document (YAML or JSON) without expanding aliases, and
// refuses one whose aliases would expand it enormously: the loader decodes
// each schema separately, which bypasses the YAML library's own guard.
func parseYAML(data []byte) (*yaml.Node, error) {
	var n yaml.Node
	if err := yaml.NewDecoder(bytes.NewReader(data)).Decode(&n); err != nil && !errors.Is(err, io.EOF) {
		return nil, err
	}
	const limit = 1_000_000
	literal := countNodes(&n, nil, 1<<62)
	if expanded := countNodes(&n, map[*yaml.Node]int{}, limit); expanded > limit && expanded > 20*literal {
		return nil, errors.New("document expands to over a million nodes through YAML aliases")
	}
	return &n, nil
}

// countNodes counts n's nodes, following aliases when memo is non-nil, and
// stops counting past limit.
func countNodes(n *yaml.Node, memo map[*yaml.Node]int, limit int) int {
	if n == nil {
		return 0
	}
	if memo != nil {
		if c, ok := memo[n]; ok {
			return c
		}
		memo[n] = limit + 1 // an alias cycle counts as too big
	}
	c := 1
	if n.Kind == yaml.AliasNode && memo != nil {
		c += countNodes(n.Alias, memo, limit)
	}
	for _, child := range n.Content {
		if c > limit {
			break
		}
		c += countNodes(child, memo, limit)
	}
	c = min(c, limit+1)
	if memo != nil {
		memo[n] = c
	}
	return c
}
