package bundler

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRemoveDebugStatements(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "simple print statement",
			input: `local x = 5
print("Hello World")
local y = 10`,
			expected: `local x = 5
local y = 10`,
		},
		{
			name: "simple warn statement",
			input: `local x = 5
warn("This is a warning")
local y = 10`,
			expected: `local x = 5
local y = 10`,
		},
		{
			name: "multiple print statements",
			input: `print("Start")
local function test()
    print("Inside function")
end
print("End")`,
			expected: `local function test()
end`,
		},
		{
			name: "indented print statements",
			input: `if true then
    print("Indented print")
    local x = 5
end`,
			expected: `if true then
    local x = 5
end`,
		},
		{
			name: "multiline print statement",
			input: `print(
    "This is a multiline",
    "print statement"
)
local x = 5`,
			expected: `local x = 5`,
		},
		{
			name: "print with complex arguments",
			input: `print("Value:", x, "Result:", calculateSomething())
local result = true`,
			expected: `local result = true`,
		},
		{
			name: "nested parentheses in print",
			input: `print("Nested:", (x + y) * (a + b))
local done = true`,
			expected: `local done = true`,
		},
		{
			name: "print inside string should not be removed",
			input: `local message = "This contains print() in string"
print("Real print")
return message`,
			expected: `local message = "This contains print() in string"
return message`,
		},
		{
			name: "warn with multiline arguments",
			input: `warn(
    "Warning message",
    variable
)
local continue = true`,
			expected: `local continue = true`,
		},
		{
			name: "mixed print and warn statements",
			input: `print("Starting")
local x = getValue()
warn("Got value:", x)
process(x)
print("Done")`,
			expected: `local x = getValue()
process(x)`,
		},
		{
			name: "empty lines preservation",
			input: `local a = 1

print("Debug")

local b = 2`,
			expected: `local a = 1


local b = 2`,
		},
		{
			name: "no debug statements",
			input: `local function calculate()
    return x + y
end
return calculate()`,
			expected: `local function calculate()
    return x + y
end
return calculate()`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := removeDebugStatements(tt.input)

			// Normalize line endings for comparison
			expected := strings.ReplaceAll(tt.expected, "\r\n", "\n")
			result = strings.ReplaceAll(result, "\r\n", "\n")

			assert.Equal(t, expected, result, "removeDebugStatements() should match expected output for %s", tt.name)
		})
	}
}

func TestRemoveDebugStatements_ComplexCases(t *testing.T) {
	input := `-- Complex test case
local function initialize()
    print("Initializing...")
    local config = {
        debug = true,
        version = "1.0"
    }
    
    if config.debug then
        print("Debug mode enabled")
        warn("This is debug output")
    end
    
    return config
end

print("Application starting")
local app = initialize()
print("Configuration loaded:", app)

-- This should remain
local function processData(data)
    return data * 2
end

warn(
    "Processing complete",
    "with result:",
    processData(42)
)`

	expected := `-- Complex test case
local function initialize()
    local config = {
        debug = true,
        version = "1.0"
    }
    
    if config.debug then
    end
    
    return config
end

local app = initialize()

-- This should remain
local function processData(data)
    return data * 2
end
`

	result := removeDebugStatements(input)

	// Normalize line endings
	expected = strings.ReplaceAll(expected, "\r\n", "\n")
	result = strings.ReplaceAll(result, "\r\n", "\n")

	assert.Equal(t, expected, result, "removeDebugStatements() complex case should match expected output")
}
