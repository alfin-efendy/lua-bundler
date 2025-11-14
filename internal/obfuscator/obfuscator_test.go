package obfuscator

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewObfuscator(t *testing.T) {
	tests := []struct {
		name          string
		level         int
		expectedLevel int
	}{
		{
			name:          "Level 1 - Basic",
			level:         1,
			expectedLevel: 1,
		},
		{
			name:          "Level 2 - Medium",
			level:         2,
			expectedLevel: 2,
		},
		{
			name:          "Level 3 - Heavy",
			level:         3,
			expectedLevel: 3,
		},
		{
			name:          "Level too low defaults to 1",
			level:         0,
			expectedLevel: 1,
		},
		{
			name:          "Level too high caps at 3",
			level:         5,
			expectedLevel: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obf := NewObfuscator(tt.level)
			require.NotNil(t, obf)
			assert.Equal(t, tt.expectedLevel, obf.level)
		})
	}
}

func TestRemoveComments(t *testing.T) {
	obf := NewObfuscator(1)

	tests := []struct {
		name        string
		input       string
		contains    []string
		notContains []string
	}{
		{
			name: "Remove single line comments",
			input: `local x = 10 -- this is a comment
local y = 20`,
			contains:    []string{"local x = 10", "local y = 20"},
			notContains: []string{"this is a comment"},
		},
		{
			name: "Remove multi-line comments",
			input: `--[[ 
This is a multi-line comment
that spans multiple lines
]]
local z = 30`,
			contains:    []string{"local z = 30"},
			notContains: []string{"multi-line comment"},
		},
		{
			name:     "Preserve comments in strings",
			input:    `local msg = "Hello -- World"`,
			contains: []string{`"Hello -- World"`},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := obf.removeComments(tt.input)

			for _, str := range tt.contains {
				assert.Contains(t, result, str)
			}

			for _, str := range tt.notContains {
				assert.NotContains(t, result, str)
			}
		})
	}
}

func TestMinifyWhitespace(t *testing.T) {
	obf := NewObfuscator(1)

	tests := []struct {
		name  string
		input string
		check func(t *testing.T, result string)
	}{
		{
			name:  "Remove extra spaces",
			input: "local  x    =    10",
			check: func(t *testing.T, result string) {
				assert.NotContains(t, result, "  ")
			},
		},
		{
			name:  "Remove empty lines",
			input: "local x = 10\n\n\nlocal y = 20",
			check: func(t *testing.T, result string) {
				assert.Equal(t, 1, strings.Count(result, "\n"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := obf.minifyWhitespace(tt.input)
			tt.check(t, result)
		})
	}
}

func TestObfuscateBasic(t *testing.T) {
	obf := NewObfuscator(1)

	code := `
-- This is a comment
local x = 10
local y = 20
return x + y
`

	result := obf.Obfuscate(code)

	assert.NotContains(t, result, "comment")
	assert.Contains(t, result, "local x")
	assert.Contains(t, result, "return")
}

func TestObfuscateMedium(t *testing.T) {
	obf := NewObfuscator(2)

	code := `
local myVariable = 10
local anotherVar = 20
return myVariable + anotherVar
`

	result := obf.Obfuscate(code)

	// Variables should be renamed
	assert.NotContains(t, result, "myVariable")
	assert.NotContains(t, result, "anotherVar")
	assert.Contains(t, result, "_0x") // Obfuscated names start with _0x
}

func TestObfuscateHeavy(t *testing.T) {
	obf := NewObfuscator(3)

	code := `
local greeting = "Hello World"
local count = 42
return greeting
`

	result := obf.Obfuscate(code)

	// Should have control flow and obfuscated identifiers
	assert.Contains(t, result, "_0x")
	assert.NotContains(t, result, "greeting")
}

func TestGenerateObfuscatedName(t *testing.T) {
	obf := NewObfuscator(1)

	name1 := obf.generateObfuscatedName()
	name2 := obf.generateObfuscatedName()

	// Should generate different names
	assert.NotEqual(t, name1, name2)

	// Should start with _0x
	assert.True(t, strings.HasPrefix(name1, "_0x"))
	assert.True(t, strings.HasPrefix(name2, "_0x"))

	// Should have correct length (_0x + 6 chars = 9)
	assert.Equal(t, 9, len(name1))
	assert.Equal(t, 9, len(name2))
}

func TestRenameIdentifiers(t *testing.T) {
	obf := NewObfuscator(2)

	code := `
local function myFunc()
	local var1 = 10
	local var2 = 20
	return var1 + var2
end
`

	result := obf.renameIdentifiers(code)

	// Local variables should be renamed
	assert.NotContains(t, result, "myFunc")
	assert.NotContains(t, result, "var1")
	assert.NotContains(t, result, "var2")

	// Reserved keywords should remain
	assert.Contains(t, result, "local")
	assert.Contains(t, result, "function")
	assert.Contains(t, result, "return")
}

func TestRenameIdentifiersPreserveRequirePaths(t *testing.T) {
	obf := NewObfuscator(2)

	tests := []struct {
		name  string
		code  string
		check func(t *testing.T, result string)
	}{
		{
			name: "Preserve path in require statement",
			code: `local PetEggs = require(Core.ReplicatedStorage.Data.PetRegistry.PetEggs)`,
			check: func(t *testing.T, result string) {
				// Variable name should be obfuscated
				assert.NotContains(t, result, "local PetEggs =")
				assert.Contains(t, result, "local _0x")

				// But the path in require() should preserve "PetEggs"
				assert.Contains(t, result, "require(Core.ReplicatedStorage.Data.PetRegistry.PetEggs)")

				// The entire path should remain intact
				assert.Contains(t, result, "Core.ReplicatedStorage.Data.PetRegistry.PetEggs")

				// Print the result for manual verification
				t.Logf("Input:  %s", `local PetEggs = require(Core.ReplicatedStorage.Data.PetRegistry.PetEggs)`)
				t.Logf("Output: %s", result)
			},
		},
		{
			name: "Preserve complex require path",
			code: `local MyModule = require(game.ReplicatedStorage.Modules.MyModule)
local result = MyModule.doSomething()`,
			check: func(t *testing.T, result string) {
				// Variable name should be obfuscated
				assert.NotContains(t, result, "local MyModule =")

				// But the module name in require path should be preserved
				assert.Contains(t, result, "require(game.ReplicatedStorage.Modules.MyModule)")

				// Usage of the variable should be obfuscated
				assert.NotContains(t, result, "MyModule.doSomething")
			},
		},
		{
			name: "Multiple requires with same name",
			code: `local Utils = require(script.Parent.Utils)
local Core = require(game.ServerStorage.Core)
local value = Utils.getValue()`,
			check: func(t *testing.T, result string) {
				// Variable names should be obfuscated
				assert.NotContains(t, result, "local Utils =")
				assert.NotContains(t, result, "local Core =")

				// But require paths should preserve original names
				assert.Contains(t, result, "require(script.Parent.Utils)")
				assert.Contains(t, result, "require(game.ServerStorage.Core)")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := obf.renameIdentifiers(tt.code)
			tt.check(t, result)
		})
	}
}
