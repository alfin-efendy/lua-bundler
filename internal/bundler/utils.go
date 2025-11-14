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

// removeComments removes all Lua comments (-- and --[[ ]]) from code
func removeComments(content string) string {
	lines := strings.Split(content, "\n")
	var result []string
	inMultilineComment := false

	for _, line := range lines {
		if inMultilineComment {
			// Check for end of multiline comment
			if strings.Contains(line, "--]]") || strings.Contains(line, "]]") {
				inMultilineComment = false
				// Get content after the closing ]]
				if idx := strings.Index(line, "]]"); idx != -1 {
					remaining := line[idx+2:]
					if strings.TrimSpace(remaining) != "" {
						result = append(result, remaining)
					}
				}
			}
			continue
		}

		// Check for start of multiline comment
		if strings.Contains(line, "--[[") {
			inMultilineComment = true
			// Get content before the opening --[[
			if idx := strings.Index(line, "--[["); idx > 0 {
				before := line[:idx]
				if strings.TrimSpace(before) != "" {
					result = append(result, before)
				}
			}
			// Check if comment ends on same line
			if strings.Contains(line, "]]") {
				inMultilineComment = false
			}
			continue
		}

		// Remove single line comments
		if idx := strings.Index(line, "--"); idx != -1 {
			// Check if -- is inside a string
			inString := false
			stringChar := rune(0)
			for i, char := range line {
				if i == idx && !inString {
					// Found comment outside string
					line = line[:idx]
					break
				}
				if char == '"' || char == '\'' {
					if !inString {
						inString = true
						stringChar = char
					} else if char == stringChar {
						inString = false
					}
				}
			}
		}

		// Add non-empty lines
		if strings.TrimSpace(line) != "" {
			result = append(result, line)
		}
	}

	return strings.Join(result, "\n")
}

// minifyCode converts code to single line by removing unnecessary whitespace
func minifyCode(content string) string {
	lines := strings.Split(content, "\n")
	var result []string

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}

	// Join with space to maintain separation between statements
	minified := strings.Join(result, " ")

	// Clean up excessive spaces
	minified = regexp.MustCompile(`\s+`).ReplaceAllString(minified, " ")

	// Remove spaces around operators and punctuation where safe
	minified = regexp.MustCompile(`\s*([,;=(){}[\]])\s*`).ReplaceAllString(minified, "$1")

	// Add space after keywords that require it
	keywords := []string{"local", "function", "if", "then", "else", "elseif", "end", "for", "while", "do", "return", "in", "and", "or", "not"}
	for _, keyword := range keywords {
		minified = regexp.MustCompile(`\b`+keyword+`\b`).ReplaceAllStringFunc(minified, func(match string) string {
			return match + " "
		})
	}

	return minified
}
