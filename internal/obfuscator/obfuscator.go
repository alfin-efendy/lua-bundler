package obfuscator

import (
	"crypto/rand"
	"math/big"
	"regexp"
	"strings"
)

// Obfuscator handles Lua code obfuscation
type Obfuscator struct {
	identifierMap map[string]string
	level         int // 1 = basic, 2 = medium, 3 = heavy
}

// NewObfuscator creates a new obfuscator instance
func NewObfuscator(level int) *Obfuscator {
	if level < 1 {
		level = 1
	}
	if level > 3 {
		level = 3
	}
	return &Obfuscator{
		identifierMap: make(map[string]string),
		level:         level,
	}
}

// Obfuscate applies obfuscation to Lua code
func (o *Obfuscator) Obfuscate(code string) string {
	result := code

	switch o.level {
	case 1:
		// Basic obfuscation: minify and remove comments
		result = o.removeComments(result)
		result = o.minifyWhitespace(result)
	case 2:
		// Medium obfuscation: basic + rename variables
		result = o.removeComments(result)
		result = o.renameIdentifiers(result)
		result = o.minifyWhitespace(result)
	case 3:
		// Heavy obfuscation: medium + extra aggressive minification
		result = o.removeComments(result)
		result = o.renameIdentifiers(result)
		result = o.minifyWhitespace(result)
		result = o.aggressiveMinify(result)
	}

	return result
}

// removeComments removes Lua comments from code
func (o *Obfuscator) removeComments(code string) string {
	// Remove multi-line comments --[[ ... ]]
	multiLineComment := regexp.MustCompile(`--\[\[[\s\S]*?\]\]`)
	code = multiLineComment.ReplaceAllString(code, "")

	// Remove single-line comments
	lines := strings.Split(code, "\n")
	var result []string
	for _, line := range lines {
		// Check if line contains a string to avoid removing -- inside strings
		if idx := strings.Index(line, "--"); idx != -1 {
			beforeComment := line[:idx]
			// Simple check: if there's an odd number of quotes before --, it's in a string
			if strings.Count(beforeComment, "\"")%2 == 0 && strings.Count(beforeComment, "'")%2 == 0 {
				line = beforeComment
			}
		}
		if strings.TrimSpace(line) != "" {
			result = append(result, line)
		}
	}
	return strings.Join(result, "\n")
}

// minifyWhitespace removes unnecessary whitespace
func (o *Obfuscator) minifyWhitespace(code string) string {
	// Replace multiple spaces with single space
	multiSpace := regexp.MustCompile(`[ \t]+`)
	code = multiSpace.ReplaceAllString(code, " ")

	// Remove spaces around operators
	code = regexp.MustCompile(`\s*([=+\-*/<>~])\s*`).ReplaceAllString(code, "$1")

	// Remove empty lines
	code = regexp.MustCompile(`\n\s*\n`).ReplaceAllString(code, "\n")

	return strings.TrimSpace(code)
}

// aggressiveMinify applies aggressive minification
func (o *Obfuscator) aggressiveMinify(code string) string {
	// Remove all newlines and replace with spaces
	code = strings.ReplaceAll(code, "\n", " ")

	// Remove spaces after specific keywords
	code = regexp.MustCompile(`\b(then|do|repeat)\s+`).ReplaceAllString(code, "$1 ")

	// Remove spaces before specific keywords
	code = regexp.MustCompile(`\s+(end|until|else)\b`).ReplaceAllString(code, " $1")

	// Collapse multiple spaces
	code = regexp.MustCompile(`\s{2,}`).ReplaceAllString(code, " ")

	return strings.TrimSpace(code)
}

// renameIdentifiers renames local variables and functions
func (o *Obfuscator) renameIdentifiers(code string) string {
	// Find local variable and function declarations
	localVarRegex := regexp.MustCompile(`\blocal\s+([a-zA-Z_][a-zA-Z0-9_]*)\b`)
	localFuncRegex := regexp.MustCompile(`\blocal\s+function\s+([a-zA-Z_][a-zA-Z0-9_]*)\b`)

	matches := localVarRegex.FindAllStringSubmatch(code, -1)
	funcMatches := localFuncRegex.FindAllStringSubmatch(code, -1)

	// Reserved Lua keywords that should not be renamed
	reserved := map[string]bool{
		"and": true, "break": true, "do": true, "else": true, "elseif": true,
		"end": true, "false": true, "for": true, "function": true, "if": true,
		"in": true, "local": true, "nil": true, "not": true, "or": true,
		"repeat": true, "return": true, "then": true, "true": true, "until": true,
		"while": true, "goto": true,
	}

	// Create mapping for identifiers
	for _, match := range matches {
		identifier := match[1]
		if !reserved[identifier] && identifier != "function" && o.identifierMap[identifier] == "" {
			o.identifierMap[identifier] = o.generateObfuscatedName()
		}
	}

	// Add function names to mapping
	for _, match := range funcMatches {
		identifier := match[1]
		if !reserved[identifier] && o.identifierMap[identifier] == "" {
			o.identifierMap[identifier] = o.generateObfuscatedName()
		}
	}

	// Replace identifiers, but preserve them inside strings
	// Split by strings first to avoid replacing inside string literals
	result := o.replaceOutsideStrings(code, o.identifierMap)

	return result
}

// replaceOutsideStrings replaces identifiers only outside of string literals
func (o *Obfuscator) replaceOutsideStrings(code string, replacements map[string]string) string {
	var result strings.Builder
	i := 0

	for i < len(code) {
		// Check if we're at the start of a string
		if code[i] == '"' || code[i] == '\'' {
			quote := code[i]
			result.WriteByte(quote)
			i++

			// Copy everything inside the string without modification
			for i < len(code) {
				if code[i] == '\\' && i+1 < len(code) {
					// Escaped character
					result.WriteByte(code[i])
					i++
					if i < len(code) {
						result.WriteByte(code[i])
						i++
					}
				} else if code[i] == quote {
					// End of string
					result.WriteByte(code[i])
					i++
					break
				} else {
					result.WriteByte(code[i])
					i++
				}
			}
		} else {
			// We're outside a string, look for identifiers to replace
			foundReplacement := false

			// Try each replacement
			for original, replacement := range replacements {
				if i+len(original) <= len(code) && code[i:i+len(original)] == original {
					// Check word boundaries
					isWordBoundaryBefore := i == 0 || !isAlphaNumOrUnderscore(code[i-1])
					isWordBoundaryAfter := i+len(original) >= len(code) || !isAlphaNumOrUnderscore(code[i+len(original)])

					if isWordBoundaryBefore && isWordBoundaryAfter {
						result.WriteString(replacement)
						i += len(original)
						foundReplacement = true
						break
					}
				}
			}

			if !foundReplacement {
				result.WriteByte(code[i])
				i++
			}
		}
	}

	return result.String()
}

// isAlphaNumOrUnderscore checks if a character is alphanumeric or underscore
func isAlphaNumOrUnderscore(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '_'
}

// generateObfuscatedName generates a random obfuscated identifier
func (o *Obfuscator) generateObfuscatedName() string {
	// Generate random identifier like _0x1a2b3c
	prefix := "_0x"
	length := 6
	chars := "0123456789abcdef"

	result := make([]byte, length)
	for i := range result {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		result[i] = chars[n.Int64()]
	}

	return prefix + string(result)
}
