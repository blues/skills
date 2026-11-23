package main

import (
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/getkin/kin-openapi/openapi3"
)

// maxDesc is the longest description kept, in characters.
const maxDesc = 120

// describe returns the first sentence of text that says something the name
// does not, with operationId cross-references rewritten to tool names and
// trimmed to about maxDesc characters. A nil name keeps the first sentence.
func (r *renderer) describe(text string, name map[string]bool) string {
	for _, s := range sentences(text) {
		s = r.rewriteRefs(s)
		if name != nil && isEcho(s, name) {
			continue
		}
		return trim(s, maxDesc)
	}
	return ""
}

// describeParam describes a param or field: the first sentence that says
// something the name doesn't, plus the first later sentence that states a
// requirement, a prohibition or what omitting it does (or gives a unit the
// first lacks), since a correct call often depends on it. When the schema has
// no enum but the description lists `code`-led values, they are appended.
func (r *renderer) describeParam(text string, name map[string]bool, schema *openapi3.SchemaRef) string {
	ss := sentences(text)
	d, next := "", 0
	for i, s := range ss {
		if s = r.rewriteRefs(s); !isEcho(s, name) {
			d, next = trim(s, maxDesc), i+1
			break
		}
	}
	for _, t := range ss[next:] {
		if noiseRE.MatchString(t) || isEcho(t, name) || strings.HasPrefix(t, "`") {
			continue // a `code`-led sentence is a list of values, not a constraint
		}
		if constraintRE.MatchString(t) || unitRE.MatchString(t) && !unitRE.MatchString(d) {
			d = joinSentences(d, trim(r.rewriteRefs(t), maxDesc))
			break
		}
	}
	if hasEnum(schema) || !scalar(schema) {
		return d
	}
	lead, vals := openingList(text)
	for _, v := range vals {
		if !strings.Contains(d, "`"+v+"`") {
			label := "Values: "
			if exampleLeadRE.MatchString(lead) {
				label = "e.g. "
			}
			return joinSentences(d, label+strings.Join(vals, "|"))
		}
	}
	return d
}

// opNote returns the first sentence of an operation's (or its request
// body's) description that says what leaving fields out does or how far the
// call reaches, e.g. "Omitted fields are cleared." or "With no filters, that
// is every device."
func (r *renderer) opNote(texts ...string) string {
	for _, text := range texts {
		for _, s := range sentences(text) {
			if omissionRE.MatchString(s) && !noiseRE.MatchString(s) {
				return trim(r.rewriteRefs(s), maxDesc)
			}
		}
	}
	return ""
}

var (
	constraintRE  = regexp.MustCompile(`(?i)\b(required|must|not allowed|cannot|can't|don't|do not|always send|if omitted|omitting|omit (it|or)|missing|leav(e|ing) (it )?out|left out)\b`)
	unitRE        = regexp.MustCompile(`(?i)\b((micro|milli|nano)?seconds|minutes|hours|days|bytes)\b`)
	noiseRE       = regexp.MustCompile(`(?i)\b(role|returns? [1-5]\d\d)\b`)
	omissionRE    = regexp.MustCompile("(?i)(\\bomit(s|ted fields|ting)?\\b [^(]|\\bmissing\\b|\\bleav(e|ing) (it |them |`\\w+` )?out\\b|\\bfields you leave out\\b|\\bdon't send\\b|\\bsend every\\b|\\balways send\\b|\\bwith no (filter|param|argument)s?\\b)")
	exampleLeadRE = regexp.MustCompile(`(?i)(\bfor example\b|\bexamples?\s*:|\be\.g\.|\bsuch as\b|\bfor instance\b)`)
	skipLeadRE    = regexp.MustCompile(`(?i)\b(errors?|status codes?|responses?|returns?|fields|properties|keys)\b`)
)

// joinSentences appends sentence t to d.
func joinSentences(d, t string) string {
	switch {
	case d == "":
		return t
	case strings.ContainsAny(d[len(d)-1:], ".!?:;…"):
		return d + " " + t
	}
	return d + ". " + t
}

// openingList returns the values documented by a Markdown list that opens
// the description: one in its first paragraph, or one that starts the next
// block with no lead-in. lead is the text that introduces it. Lists led in by
// errors, responses or fields don't list values.
func openingList(text string) (lead string, vals []string) {
	blocks := strings.Split(strings.ReplaceAll(strings.TrimSpace(text), "\r\n", "\n"), "\n\n")
	var leadLines, items []string
	for bi := 0; bi < len(blocks) && bi < 2 && len(items) == 0; bi++ {
		var lines []string
		for _, line := range strings.Split(blocks[bi], "\n") {
			if line = strings.TrimSpace(line); line != "" {
				lines = append(lines, line)
			}
		}
		if bi == 1 {
			if len(lines) == 0 {
				break
			}
			if _, ok := listItem(lines[0]); !ok {
				break // a later list must open its block
			}
		}
		for _, line := range lines {
			if item, ok := listItem(line); ok {
				items = append(items, item)
			} else if len(items) == 0 {
				leadLines = append(leadLines, line)
			}
		}
	}
	if len(items) == 0 {
		return "", nil
	}
	lead = strings.Join(leadLines, " ")
	if skipLeadRE.MatchString(lead) {
		return lead, nil
	}
	for _, it := range items {
		heads := leadingCode(it)
		if len(heads) == 0 {
			return lead, nil
		}
		vals = appendUnique(vals, heads...)
	}
	return lead, vals
}

var leadingCodeRE = regexp.MustCompile("^`([^`]+)`")

// leadingCode returns the code spans an item starts with, separated by
// commas, semicolons, slashes or "or", skipping a parenthetical after each.
func leadingCode(item string) []string {
	var out []string
	s := item
	for {
		m := leadingCodeRE.FindStringSubmatch(s)
		if m == nil {
			return out
		}
		out = appendUnique(out, m[1])
		s = strings.TrimSpace(s[len(m[0]):])
		if strings.HasPrefix(s, "(") {
			if i := strings.IndexByte(s, ')'); i >= 0 {
				s = strings.TrimSpace(s[i+1:])
			}
		}
		switch {
		case s == "":
			return out
		case strings.IndexByte(",;/|", s[0]) >= 0:
			s = strings.TrimSpace(s[1:])
		case !strings.HasPrefix(s, "or "):
			return out
		}
		s = strings.TrimSpace(strings.TrimPrefix(s, "or "))
	}
}

var (
	codeSpanRE  = regexp.MustCompile("`[^`]+`")
	camelWordRE = regexp.MustCompile(`\b[A-Z][a-z0-9]*[A-Z][A-Za-z0-9]*\b`)
	mdLinkRE    = regexp.MustCompile(`!?\[([^\]]+)\]\([^)\s]+\)`)
	headingRE   = regexp.MustCompile(`^#{1,6}\s+[\p{Lu}\d]`) // "# of items" is prose
)

// stripLinks replaces Markdown links and images with their text, outside
// code spans.
func stripLinks(line string) string {
	parts := strings.Split(line, "`")
	for i := 0; i < len(parts); i += 2 { // even parts are outside code
		parts[i] = mdLinkRE.ReplaceAllString(parts[i], "$1")
	}
	return strings.Join(parts, "`")
}

// rewriteRefs replaces operationIds mentioned in descriptions, `backticked`
// or as bare CamelCase words, with their tool names.
func (r *renderer) rewriteRefs(s string) string {
	s = codeSpanRE.ReplaceAllStringFunc(s, func(m string) string {
		if n, ok := r.names[m[1:len(m)-1]]; ok {
			return "`" + n + "`"
		}
		return m
	})
	return camelWordRE.ReplaceAllStringFunc(s, func(m string) string {
		if n, ok := r.names[m]; ok {
			return n
		}
		return m
	})
}

// sentences splits Markdown-ish text into sentences. Soft-wrapped lines are
// joined; blank lines end a sentence; list items are joined with "; ",
// continuing a sentence that ends with ":"; fenced code is dropped.
func sentences(text string) []string {
	var paras []string
	var cur strings.Builder
	inList, inFence := false, false
	flush := func() {
		if cur.Len() > 0 {
			paras = append(paras, cur.String())
			cur.Reset()
		}
		inList = false
	}
	for _, line := range strings.Split(text, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "```") || strings.HasPrefix(line, "~~~") {
			inFence = !inFence
			flush()
			continue
		}
		if inFence {
			continue
		}
		if line == "" {
			flush()
			continue
		}
		if headingRE.MatchString(line) { // a Markdown heading is not a sentence
			flush()
			continue
		}
		line = stripLinks(line)
		item, isItem := listItem(line)
		switch {
		case isItem && inList:
			cur.WriteString("; ")
			cur.WriteString(item)
		case isItem && strings.HasSuffix(cur.String(), ":"):
			cur.WriteString(" ")
			cur.WriteString(item)
			inList = true
		case isItem:
			flush()
			cur.WriteString(item)
			inList = true
		default:
			if cur.Len() > 0 {
				cur.WriteByte(' ')
			}
			cur.WriteString(line)
		}
	}
	flush()
	var out []string
	for _, p := range paras {
		out = append(out, splitSentences(strings.Join(strings.Fields(p), " "))...)
	}
	return out
}

var listMarkerRE = regexp.MustCompile(`^(?:[-*+]|\d+[.)])\s+`)

func listItem(line string) (string, bool) {
	if m := listMarkerRE.FindString(line); m != "" {
		return line[len(m):], true
	}
	return line, false
}

// splitSentences splits one paragraph at '.', '!' or '?' followed by a space
// or the end, outside `code`, ignoring common abbreviations.
func splitSentences(p string) []string {
	var out []string
	start, inCode := 0, false
	for i := 0; i < len(p); i++ {
		c := p[i]
		if c == '`' {
			inCode = !inCode
			continue
		}
		if inCode || (c != '.' && c != '!' && c != '?') {
			continue
		}
		end := i + 1
		for end < len(p) && strings.IndexByte(`)"'`, p[end]) >= 0 {
			end++
		}
		if end < len(p) && p[end] != ' ' {
			continue
		}
		if c == '.' && isAbbrev(p[start:i+1]) {
			continue
		}
		if s := strings.TrimSpace(p[start:end]); s != "" {
			out = append(out, s)
		}
		start, i = end, end-1
	}
	if s := strings.TrimSpace(p[start:]); s != "" {
		out = append(out, s)
	}
	return out
}

var abbrevs = []string{"e.g.", "i.e.", "etc.", "vs.", "cf.", "approx.", "incl."}

func isAbbrev(s string) bool {
	s = strings.ToLower(s)
	for _, a := range abbrevs {
		if strings.HasSuffix(s, a) && (len(s) == len(a) || !isWordByte(s[len(s)-len(a)-1])) {
			return true
		}
	}
	return false
}

func isWordByte(c byte) bool {
	return c == '_' || c >= '0' && c <= '9' || c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z'
}

// trim shortens s to at most max characters at a word boundary, adding "…".
func trim(s string, max int) string {
	if utf8.RuneCountInString(s) <= max {
		return s
	}
	rs := []rune(s)
	cut := max - 1
	i := cut
	for i > max/2 && rs[i] != ' ' {
		i--
	}
	if rs[i] != ' ' {
		i = cut
	}
	out := strings.TrimRight(string(rs[:i]), " ,;:-—(")
	if strings.Count(out, "`")%2 == 1 { // don't leave a code span open
		out = strings.TrimRight(out[:strings.LastIndexByte(out, '`')], " ,;:-—(")
	}
	return out + "…"
}

// fillerWords never carry information on their own in a description.
var fillerWords = map[string]bool{
	"a": true, "an": true, "the": true, "of": true, "in": true, "on": true, "for": true,
	"to": true, "by": true, "with": true, "from": true, "at": true, "and": true, "or": true,
	"its": true, "this": true, "that": true, "is": true, "are": true, "be": true, "as": true,
	"filter": true, "specified": true, "given": true, "identifier": true,
}

// isEcho reports whether every word of s is filler or part of the name, so
// that s tells an agent nothing the name doesn't.
func isEcho(s string, name map[string]bool) bool {
	known := func(w string) bool {
		w = stem(strings.ToLower(w))
		return fillerWords[w] || name[w]
	}
	s = strings.NewReplacer("'s ", " ", "’s ", " ", "'s.", ".", "’s.", ".").Replace(s + " ")
	for _, tok := range strings.FieldsFunc(s, notAlnum) {
		if known(tok) {
			continue
		}
		parts := identWords(tok)
		if len(parts) < 2 {
			return false
		}
		for _, p := range parts {
			if !known(p) {
				return false
			}
		}
	}
	return true
}

// oneLine collapses all whitespace runs to single spaces.
func oneLine(s string) string { return strings.Join(strings.Fields(s), " ") }
