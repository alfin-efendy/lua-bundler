package lua

import "strings"

// tokenKind classifies a Lua lexical token.
type tokenKind int

const (
	tkName    tokenKind = iota // identifier or keyword
	tkNumber                   // numeric literal
	tkString                   // short or long string literal
	tkComment                  // line or block comment
	tkOp                       // operator or punctuation
)

// token is one lexical unit of Lua source. text is the original bytes.
type token struct {
	kind tokenKind
	text string
}

// lex scans Lua source into tokens. Whitespace is consumed but never emitted.
// The lexer assumes reasonably well-formed input; an unterminated string or
// long bracket is consumed to end of input rather than panicking.
func lex(src string) []token {
	var tokens []token
	i, n := 0, len(src)

	for i < n {
		c := src[i]

		if c == ' ' || c == '\t' || c == '\n' || c == '\r' || c == '\v' || c == '\f' {
			i++
			continue
		}

		start := i
		switch {
		case c == '-' && i+1 < n && src[i+1] == '-':
			i += 2
			if level, ok := longBracketLevel(src, i); ok {
				i = scanLongBracket(src, i, level)
			} else {
				for i < n && src[i] != '\n' {
					i++
				}
			}
			tokens = append(tokens, token{tkComment, src[start:i]})

		case c == '[' && isLongBracketOpen(src, i):
			level, _ := longBracketLevel(src, i)
			i = scanLongBracket(src, i, level)
			tokens = append(tokens, token{tkString, src[start:i]})

		case c == '"' || c == '\'':
			i = scanShortString(src, i)
			tokens = append(tokens, token{tkString, src[start:i]})

		case isDigit(c) || (c == '.' && i+1 < n && isDigit(src[i+1])):
			i = scanNumber(src, i)
			tokens = append(tokens, token{tkNumber, src[start:i]})

		case isNameStart(c):
			i++
			for i < n && isNameCont(src[i]) {
				i++
			}
			tokens = append(tokens, token{tkName, src[start:i]})

		default:
			i = scanOp(src, i)
			tokens = append(tokens, token{tkOp, src[start:i]})
		}
	}

	return tokens
}

func isDigit(c byte) bool { return c >= '0' && c <= '9' }

func isHexDigit(c byte) bool {
	return isDigit(c) || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')
}

func isNameStart(c byte) bool {
	return c == '_' || (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
}

func isNameCont(c byte) bool { return isNameStart(c) || isDigit(c) }

// longBracketLevel reports the number of '=' in a long-bracket opener
// "[" "="* "[" starting at src[i]. ok is false if src[i:] is not such an opener.
func longBracketLevel(src string, i int) (level int, ok bool) {
	if i >= len(src) || src[i] != '[' {
		return 0, false
	}
	j := i + 1
	for j < len(src) && src[j] == '=' {
		j++
	}
	if j < len(src) && src[j] == '[' {
		return j - (i + 1), true
	}
	return 0, false
}

func isLongBracketOpen(src string, i int) bool {
	_, ok := longBracketLevel(src, i)
	return ok
}

// scanLongBracket consumes a long-bracket body (string or comment) whose opener
// starts at src[i] with the given level, returning the index just past the
// closer "]" "="*level "]". Returns len(src) if the closer is missing.
func scanLongBracket(src string, i, level int) int {
	i += 2 + level // skip "[" "="*level "["
	closer := "]" + strings.Repeat("=", level) + "]"
	if idx := strings.Index(src[i:], closer); idx >= 0 {
		return i + idx + len(closer)
	}
	return len(src)
}

// scanShortString consumes a "..." or '...' literal starting at the quote and
// returns the index just past the closing quote. Backslash escapes are skipped.
// An unterminated literal stops at the newline or end of input.
func scanShortString(src string, i int) int {
	quote := src[i]
	i++
	for i < len(src) {
		switch src[i] {
		case '\\':
			i += 2
			continue
		case quote:
			return i + 1
		case '\n':
			return i
		}
		i++
	}
	return len(src)
}

// scanNumber consumes a numeric literal starting at i and returns the index just
// past it. Handles decimal and hex integers/floats with exponents.
func scanNumber(src string, i int) int {
	n := len(src)
	if src[i] == '0' && i+1 < n && (src[i+1] == 'x' || src[i+1] == 'X') {
		i += 2
		for i < n && (isHexDigit(src[i]) || src[i] == '.') {
			i++
		}
		if i < n && (src[i] == 'p' || src[i] == 'P') {
			i++
			if i < n && (src[i] == '+' || src[i] == '-') {
				i++
			}
			for i < n && isDigit(src[i]) {
				i++
			}
		}
		return i
	}
	for i < n && (isDigit(src[i]) || src[i] == '.') {
		i++
	}
	if i < n && (src[i] == 'e' || src[i] == 'E') {
		i++
		if i < n && (src[i] == '+' || src[i] == '-') {
			i++
		}
		for i < n && isDigit(src[i]) {
			i++
		}
	}
	return i
}

// multiCharOps are Lua operators of length >= 2, longest first so greedy
// matching prefers "..." over "..=" over ".." over ".".
var multiCharOps = []string{
	"...", "..=", "//=", "<<", ">>",
	"..", "==", "~=", "<=", ">=", "::", "//",
	"+=", "-=", "*=", "/=", "%=", "^=",
}

// scanOp consumes one operator/punctuation token starting at i.
func scanOp(src string, i int) int {
	for _, op := range multiCharOps {
		if strings.HasPrefix(src[i:], op) {
			return i + len(op)
		}
	}
	return i + 1
}
