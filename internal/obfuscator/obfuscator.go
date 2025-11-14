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

	// Skip operator minification to avoid breaking syntax
	// It's safer to preserve spaces around operators

	// Remove empty lines
	code = regexp.MustCompile(`\n\s*\n`).ReplaceAllString(code, "\n")

	return strings.TrimSpace(code)
}

// aggressiveMinify applies aggressive minification
func (o *Obfuscator) aggressiveMinify(code string) string {
	// Remove all newlines and replace with spaces (single line output)
	code = strings.ReplaceAll(code, "\n", " ")

	// Remove spaces after specific keywords that don't need them
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

	// Also capture identifiers from multi-line local declarations without assignment
	// Pattern: lines starting with "local" followed by just identifier (no =, no function)
	multiLineLocalRegex := regexp.MustCompile(`(?m)^\s*local\s+([a-zA-Z_][a-zA-Z0-9_]*)\s*$`)
	multiLineMatches := multiLineLocalRegex.FindAllStringSubmatch(code, -1)

	// Reserved Lua keywords and Roblox globals that should not be renamed
	reserved := map[string]bool{
		// Lua keywords
		"and": true, "break": true, "do": true, "else": true, "elseif": true,
		"end": true, "false": true, "for": true, "function": true, "if": true,
		"in": true, "local": true, "nil": true, "not": true, "or": true,
		"repeat": true, "return": true, "then": true, "true": true, "until": true,
		"while": true, "goto": true,
		// Roblox/Luau globals that should be preserved
		"game": true, "workspace": true, "script": true, "task": true, "wait": true,
		"spawn": true, "delay": true, "tick": true, "time": true, "typeof": true,
		"pcall": true, "xpcall": true, "pairs": true, "ipairs": true, "next": true,
		"getmetatable": true, "setmetatable": true, "rawget": true, "rawset": true,
		"print": true, "warn": true, "error": true, "assert": true, "select": true,
		"unpack": true, "tostring": true, "tonumber": true, "type": true,
		"loadstring": true, "require": true, "getfenv": true, "setfenv": true,
		// Additional Roblox services and globals
		"Instance": true, "Vector3": true, "Vector2": true, "CFrame": true,
		"Color3": true, "UDim2": true, "Enum": true, "Random": true, "Region3": true,
		"TweenInfo": true, "BrickColor": true, "Ray": true, "Faces": true,
		// Executor HTTP request functions (CRITICAL - must not be renamed)
		"request": true, "http_request": true, "http": true, "syn": true,
		"fluxus": true, "HttpService": true,
		// Common executor globals
		"getgenv": true, "getrenv": true, "getrawmetatable": true,
		"hookfunction": true, "hookmetamethod": true, "islclosure": true,
		"newcclosure": true, "checkcaller": true, "cloneref": true,
		"compareinstances": true, "Drawing": true, "WebSocket": true,
		"crypt": true, "base64": true, "base64_encode": true, "base64_decode": true,
		"readfile": true, "writefile": true, "appendfile": true, "makefolder": true,
		"isfolder": true, "isfile": true, "listfiles": true, "delfile": true, "delfolder": true,
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

	// Add multi-line local declarations (e.g., "local Core" on its own line)
	for _, match := range multiLineMatches {
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

// replaceOutsideStrings replaces identifiers only outside of string literals and require paths
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
		} else if o.isAtRequirePath(code, i) {
			// We're at a require() call
			// Write "require"
			result.WriteString("require")
			i += 7

			// Skip whitespace and find the opening parenthesis
			for i < len(code) && (code[i] == ' ' || code[i] == '\t') {
				result.WriteByte(code[i])
				i++
			}

			if i < len(code) && code[i] == '(' {
				result.WriteByte(code[i]) // Write '('
				i++

				// Process the content inside require()
				// We need to replace variables but preserve the last component (file/module name)
				parenDepth := 1
				requireContent := strings.Builder{}

				// Collect everything inside require()
				for i < len(code) && parenDepth > 0 {
					if code[i] == '(' {
						parenDepth++
						requireContent.WriteByte(code[i])
						i++
					} else if code[i] == ')' {
						parenDepth--
						if parenDepth > 0 {
							requireContent.WriteByte(code[i])
							i++
						}
					} else {
						requireContent.WriteByte(code[i])
						i++
					}
				}

				// Process the require content to preserve only the last component
				processedContent := o.processRequireContent(requireContent.String(), replacements)
				result.WriteString(processedContent)

				// Write the closing parenthesis
				if i < len(code) && code[i] == ')' {
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

					// Don't replace if this is a property access (preceded by a dot)
					// This prevents replacing identifiers in expressions like "data.Rarity" or "obj.propertyName"
					isPropertyAccess := i > 0 && code[i-1] == '.'

					// Don't replace if this is a table key (identifier followed = in a table constructor)
					// This prevents replacing table keys like in {text = "value", Button = obj}
					isTableKeyResult := o.isTableKey(code, i, i+len(original))

					if isWordBoundaryBefore && isWordBoundaryAfter && !isPropertyAccess && !isTableKeyResult {
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

// isAtRequirePath checks if the current position is at the start of a require() call
func (o *Obfuscator) isAtRequirePath(code string, pos int) bool {
	// Check if we're at "require"
	if pos+7 > len(code) {
		return false
	}
	if code[pos:pos+7] != "require" {
		return false
	}
	// Check word boundary before
	if pos > 0 && isAlphaNumOrUnderscore(code[pos-1]) {
		return false
	}
	// Check word boundary after
	if pos+7 < len(code) && isAlphaNumOrUnderscore(code[pos+7]) {
		return false
	}
	return true
}

// isTableKey checks if the identifier at the given position is a table key
// (i.e., it's followed by = and we're inside a table constructor)
func (o *Obfuscator) isTableKey(code string, startPos, currentPos int) bool {
	pos := currentPos

	// Skip whitespace after the identifier
	for pos < len(code) && (code[pos] == ' ' || code[pos] == '\t' || code[pos] == '\n' || code[pos] == '\r') {
		pos++
	}

	// Check if followed by = (table key assignment)
	// NOTE: We do NOT check for : here because:
	// - Core:Method() is a method call (should be obfuscated)
	// - {key: value} is rare Lua 5.2+ syntax not commonly used
	if pos >= len(code) || code[pos] != '=' {
		return false
	}

	// Make sure it's not ==, <=, >=, ~=
	if code[pos] == '=' {
		if pos+1 < len(code) && code[pos+1] == '=' {
			return false
		}
		// Check if it's preceded by <, >, ~
		if pos > 0 && (code[pos-1] == '<' || code[pos-1] == '>' || code[pos-1] == '~') {
			return false
		}
	}

	// Now check if we're inside a table constructor by looking backwards
	// We need to find an opening { before we find a closing } or statement boundary
	braceDepth := 0
	parenDepth := 0

	for i := startPos - 1; i >= 0; i-- {
		ch := code[i]

		// Track parentheses depth
		if ch == ')' {
			parenDepth++
		} else if ch == '(' {
			parenDepth--
		}

		// Track brace depth
		if ch == '}' {
			braceDepth++
		} else if ch == '{' {
			if braceDepth == 0 {
				// Found an unmatched opening brace, we're in a table constructor
				return true
			}
			braceDepth--
		}

		// Only check for statement boundaries when we're at the outermost level (not inside any braces/parens)
		if braceDepth == 0 && parenDepth == 0 {
			// Check for 'local' keyword which indicates variable declaration
			if i >= 4 {
				startIdx := i - 4
				if startIdx < 0 {
					startIdx = 0
				}
				if i+1 <= len(code) && code[startIdx:i+1] == "local" {
					// Make sure 'local' has word boundary before it
					if startIdx == 0 || !isAlphaNumOrUnderscore(code[startIdx-1]) {
						// Make sure 'local' is followed by whitespace
						if i+1 >= len(code) || !isAlphaNumOrUnderscore(code[i+1]) {
							// We found 'local' at the same level - this is variable declaration, not table key
							return false
						}
					}
				}
			}
		}
	}

	return false
}

// processRequireContent processes the content inside require() to replace variables
// but preserve the last component (module/file name)
func (o *Obfuscator) processRequireContent(content string, replacements map[string]string) string {
	content = strings.TrimSpace(content)

	// Handle string literals in require - don't process them
	if strings.HasPrefix(content, "\"") || strings.HasPrefix(content, "'") {
		return content
	}

	// Split by dots to get path components
	// We need to be careful about nested parentheses or function calls
	parts := strings.Split(content, ".")
	if len(parts) == 0 {
		return content
	}

	// Process each part except the last one (which is the module/file name)
	var processedParts []string
	for i, part := range parts {
		part = strings.TrimSpace(part)

		// Last component (module/file name) should not be replaced
		if i == len(parts)-1 {
			processedParts = append(processedParts, part)
			continue
		}

		// For other parts, try to replace with obfuscated name
		replaced := false
		for original, replacement := range replacements {
			if part == original {
				processedParts = append(processedParts, replacement)
				replaced = true
				break
			}
		}

		if !replaced {
			processedParts = append(processedParts, part)
		}
	}

	return strings.Join(processedParts, ".")
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
