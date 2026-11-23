package main

import (
	"strings"
	"unicode"
)

// toolName converts an operationId to the MCP tool name, snake_case: "_"
// goes before an uppercase letter that follows a lowercase letter or a digit,
// or that starts a word after an acronym; then everything is lowercased.
// Digits and existing "_" stay. ListPets -> list_pets, GetHTTPStatus ->
// get_http_status, OAuth2Token -> o_auth2_token, getV2Users -> get_v2_users.
// This must match the MCP server's own rule; keep it the only conversion.
func toolName(id string) string {
	rs := []rune(id)
	var b strings.Builder
	for i, c := range rs {
		if i > 0 && unicode.IsUpper(c) {
			prev := rs[i-1]
			nextLower := i+1 < len(rs) && unicode.IsLower(rs[i+1])
			if unicode.IsLower(prev) || unicode.IsDigit(prev) || unicode.IsUpper(prev) && nextLower {
				b.WriteByte('_')
			}
		}
		b.WriteRune(unicode.ToLower(c))
	}
	return b.String()
}

// identWords splits an identifier at non-alphanumerics and at the word
// boundaries toolName uses.
func identWords(s string) []string {
	var words []string
	for _, part := range strings.FieldsFunc(s, notAlnum) {
		for _, w := range strings.Split(toolName(part), "_") {
			if w != "" {
				words = append(words, w)
			}
		}
	}
	return words
}

// nameWords returns the stemmed words of an identifier plus all of them
// joined, so that "ownerUID" covers "owner", "uid" and "OwnerUID".
func nameWords(name string) map[string]bool {
	set := map[string]bool{}
	var joined strings.Builder
	for _, w := range identWords(name) {
		set[stem(w)] = true
		joined.WriteString(w)
	}
	set[stem(joined.String())] = true
	return set
}

// stem crudely singularizes a lowercase word.
func stem(w string) string {
	switch {
	case len(w) > 4 && strings.HasSuffix(w, "ies"):
		return w[:len(w)-3] + "y"
	case len(w) > 3 && strings.HasSuffix(w, "s") && !strings.HasSuffix(w, "ss"):
		return w[:len(w)-1]
	}
	return w
}

func notAlnum(r rune) bool { return !unicode.IsLetter(r) && !unicode.IsDigit(r) }
