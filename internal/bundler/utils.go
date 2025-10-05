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
