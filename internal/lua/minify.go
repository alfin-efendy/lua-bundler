package lua

import "strings"

// Minify collapses Lua source to a single line, removing comments and
// unnecessary whitespace while preserving the byte-for-byte contents of string
// literals. It is string- and comment-aware via the lexer, so keywords and
// punctuation inside string literals are never altered. It cannot fail.
func Minify(src string) string {
	tokens := lex(src)

	var b strings.Builder
	var prev *token
	for i := range tokens {
		t := &tokens[i]
		if t.kind == tkComment {
			continue // drop comments
		}
		if prev != nil && needsSpace(*prev, *t) {
			b.WriteByte(' ')
		}
		b.WriteString(t.text)
		prev = t
	}
	return b.String()
}

// isWordChar reports whether c can be part of a Lua name or number, so that two
// adjacent word chars would lex as a single token.
func isWordChar(c byte) bool {
	return c == '_' || (c >= '0' && c <= '9') ||
		(c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
}

// mergePairs are two-character sequences that, written adjacently, would lex as
// a single longer operator or start a comment — so a space must separate them.
var mergePairs = map[string]bool{
	"--": true, "..": true, "[[": true, "[=": true,
	"==": true, "~=": true, "<=": true, ">=": true,
	"<<": true, ">>": true, "//": true, "::": true,
}

// needsSpace reports whether a separating space is required between the already-
// emitted token prev and the next token cur so that re-lexing yields the same
// two tokens.
func needsSpace(prev, cur token) bool {
	a, b := prev.text, cur.text
	if a == "" || b == "" {
		return false
	}
	la, fb := a[len(a)-1], b[0]
	if isWordChar(la) && isWordChar(fb) {
		return true
	}
	if prev.kind == tkNumber && fb == '.' {
		return true
	}
	return mergePairs[string([]byte{la, fb})]
}
