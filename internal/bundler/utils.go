package bundler

import (
	"regexp"
	"strings"
)

// removeDebugStatements removes print() and warn() statements for release mode
func removeDebugStatements(content string) string {
	lines := strings.Split(content, "\n")
	var result []string

	// Regex patterns for detecting print and warn statements
	printRegex := regexp.MustCompile(`^\s*print\s*\(`)
	warnRegex := regexp.MustCompile(`^\s*warn\s*\(`)

	inMultilineStatement := false
	parenDepth := 0

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Skip empty lines
		if trimmed == "" {
			result = append(result, line)
			continue
		}

		// Check if this line starts a print/warn statement
		if !inMultilineStatement && (printRegex.MatchString(line) || warnRegex.MatchString(line)) {
			inMultilineStatement = true
			parenDepth = 0

			// Count parentheses in this line
			for _, char := range line {
				if char == '(' {
					parenDepth++
				} else if char == ')' {
					parenDepth--
				}
			}

			// If statement ends on the same line
			if parenDepth <= 0 {
				inMultilineStatement = false
			}
			continue // Skip this line
		}

		// If we're in a multiline print/warn statement
		if inMultilineStatement {
			for _, char := range line {
				if char == '(' {
					parenDepth++
				} else if char == ')' {
					parenDepth--
				}
			}

			if parenDepth <= 0 {
				inMultilineStatement = false
			}
			continue // Skip this line
		}

		// Add lines that are not print/warn
		result = append(result, line)
	}

	return strings.Join(result, "\n")
}

// minifyCode collapses Lua source to a single line, removing comments and
// unnecessary whitespace while preserving the byte-for-byte contents of string
// literals. It is string- and comment-aware via a real Lua lexer (see lex), so
// keywords and punctuation inside string literals are never altered.
func minifyCode(content string) string {
	tokens := lex(content)

	var b strings.Builder
	prev := ""
	for _, t := range tokens {
		if t.kind == tkComment {
			continue // drop comments
		}
		if prev != "" && needsSpace(prev, t.text) {
			b.WriteByte(' ')
		}
		b.WriteString(t.text)
		prev = t.text
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

// needsSpace reports whether a separating space is required between already-
// emitted token text a and the next token text b so that re-lexing yields the
// same two tokens.
func needsSpace(a, b string) bool {
	if a == "" || b == "" {
		return false
	}
	la, fb := a[len(a)-1], b[0]
	if isWordChar(la) && isWordChar(fb) {
		return true
	}
	return mergePairs[string([]byte{la, fb})]
}
