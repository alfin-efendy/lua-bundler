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
		{
			name: "Replace variable in require path but preserve module name",
			code: `local Core = _core
local ShopData = require(Core.ReplicatedStorage.Data.EventShopData)`,
			check: func(t *testing.T, result string) {
				// Both variable names should be obfuscated
				assert.NotContains(t, result, "local Core =")
				assert.NotContains(t, result, "local ShopData =")
				assert.Contains(t, result, "local _0x")

				// The variable "Core" in require path should be replaced
				assert.NotContains(t, result, "require(Core.ReplicatedStorage")

				// But the module name "EventShopData" should be preserved
				assert.Contains(t, result, ".EventShopData)")

				// Print the result for manual verification
				t.Logf("Input:  local Core = _core\\nlocal ShopData = require(Core.ReplicatedStorage.Data.EventShopData)")
				t.Logf("Output: %s", result)
			},
		},
		{
			name: "Complex scenario with variable used in require",
			code: `local PetEggs = require(Core.ReplicatedStorage.Data.PetRegistry.PetEggs)
local result = PetEggs.getAll()`,
			check: func(t *testing.T, result string) {
				// Variable name should be obfuscated
				assert.NotContains(t, result, "local PetEggs =")

				// Module name in require should be preserved
				assert.Contains(t, result, ".PetEggs)")

				// Usage of the variable should be obfuscated
				assert.NotContains(t, result, "PetEggs.getAll")

				t.Logf("Output: %s", result)
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

func TestPreservePropertyAccess(t *testing.T) {
	obf := NewObfuscator(2)

	tests := []struct {
		name  string
		code  string
		check func(t *testing.T, result string)
	}{
		{
			name: "Preserve property names in object access",
			code: `local Rarity = require(RarityModule)
local data = {}
local value = data.Rarity`,
			check: func(t *testing.T, result string) {
				// Variable name should be obfuscated
				assert.NotContains(t, result, "local Rarity =")

				// But property access should NOT be obfuscated
				assert.Contains(t, result, ".Rarity")

				t.Logf("Output: %s", result)
			},
		},
		{
			name: "Complex property access scenario",
			code: `local items = SeasonPassShop:GetItemRepository()
local sortedList = {}

for itemName, data in pairs(items) do
    data._name = itemName
    table.insert(sortedList, data)
end

table.sort(sortedList, function(a, b)
    local rarityA = Rarity.RarityOrder[a.Rarity] or 99
    local rarityB = Rarity.RarityOrder[b.Rarity] or 99
    
    if rarityA == rarityB then
        if a.LayoutOrder == b.LayoutOrder then
            return a._name < b._name
        else
            return a.LayoutOrder < b.LayoutOrder
        end
    end
    
    return rarityA < rarityB
end)`,
			check: func(t *testing.T, result string) {
				// Variables should be obfuscated
				assert.NotContains(t, result, "local items =")
				assert.NotContains(t, result, "local sortedList =")
				assert.NotContains(t, result, "local rarityA =")
				assert.NotContains(t, result, "local rarityB =")

				// But property names should NOT be obfuscated
				assert.Contains(t, result, ".Rarity")
				assert.Contains(t, result, ".RarityOrder")
				assert.Contains(t, result, ".LayoutOrder")
				assert.Contains(t, result, "._name")

				// Property access should preserve the property name
				assert.Contains(t, result, "a.Rarity")
				assert.Contains(t, result, "b.Rarity")
				assert.Contains(t, result, "a.LayoutOrder")
				assert.Contains(t, result, "b.LayoutOrder")
				assert.Contains(t, result, "a._name")
				assert.Contains(t, result, "b._name")

				t.Logf("Output: %s", result)
			},
		},
		{
			name: "Property access in table literal",
			code: `local data = getData()
local rarity = data.Rarity or "Unknown"
local name = data._name or "Unnamed"
table.insert(itemNames, {text= "[" .. rarity .. "] " .. name, value=name})`,
			check: func(t *testing.T, result string) {
				// Variables should be obfuscated
				assert.NotContains(t, result, "local rarity =")
				assert.NotContains(t, result, "local name =")
				assert.NotContains(t, result, "local data =")

				// Property names should NOT be obfuscated
				assert.Contains(t, result, ".Rarity")
				assert.Contains(t, result, "._name")

				t.Logf("Output: %s", result)
			},
		},
		{
			name: "Real world scenario from user",
			code: `local items = SeasonPassShop:GetItemRepository()
local sortedList = {}
local itemNames = {}

for itemName, data in pairs(items) do
    data._name = itemName
    table.insert(sortedList, data)
end

table.sort(sortedList, function(a, b)
    local rarityA = Rarity.RarityOrder[a.Rarity] or 99
    local rarityB = Rarity.RarityOrder[b.Rarity] or 99

    if rarityA == rarityB then
        if a.LayoutOrder == b.LayoutOrder then
            return a._name < b._name
        else
            return a.LayoutOrder < b.LayoutOrder
        end
    end

    return rarityA < rarityB
end)

for _, data in pairs(sortedList) do
    local rarity = data.Rarity or "Unknown"
    local name = data._name or "Unnamed"
    table.insert(itemNames, {text= "[" .. rarity .. "] " .. name, value=name})
end`,
			check: func(t *testing.T, result string) {
				// Variables should be obfuscated
				assert.NotContains(t, result, "local items =")
				assert.NotContains(t, result, "local sortedList =")
				assert.NotContains(t, result, "local itemNames =")
				assert.NotContains(t, result, "local rarityA =")
				assert.NotContains(t, result, "local rarityB =")
				assert.NotContains(t, result, "local rarity =")
				assert.NotContains(t, result, "local name =")

				// Property names should NOT be obfuscated - they should remain as is
				// (the variable name before the dot will be obfuscated, but the property name after the dot should not)
				assert.Contains(t, result, ".Rarity")
				assert.Contains(t, result, ".LayoutOrder")
				assert.Contains(t, result, "._name")
				assert.Contains(t, result, ".RarityOrder")

				// Specifically check that property access is preserved
				assert.Contains(t, result, "a.Rarity")
				assert.Contains(t, result, "b.Rarity")
				assert.Contains(t, result, "a.LayoutOrder")
				assert.Contains(t, result, "b.LayoutOrder")
				assert.Contains(t, result, "a._name")
				assert.Contains(t, result, "b._name")

				t.Logf("Full Output:\n%s", result)
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

func TestPreserveTableKeys(t *testing.T) {
	obf := NewObfuscator(2)

	tests := []struct {
		name  string
		code  string
		check func(t *testing.T, result string)
	}{
		{
			name: "Preserve table keys in constructor",
			code: `local text = "Hello"
local value = 123
local obj = {text = text, value = value}`,
			check: func(t *testing.T, result string) {
				// Table keys should NOT be obfuscated
				assert.Contains(t, result, "{text =")
				assert.Contains(t, result, "value =")

				// But the values (variable references) should be obfuscated when they match local vars
				// Note: text and value are NOT obfuscated in declaration because regex doesn't match
				// string/number literals on the right side

				t.Logf("Output: %s", result)
			},
		},
		{
			name: "Preserve Button key in table",
			code: `local Button = createButton()
local result = {Button = Button}
result.Button.MouseButton1Click:Connect(handler)`,
			check: func(t *testing.T, result string) {
				// Table key should NOT be obfuscated
				assert.Contains(t, result, "{Button =")

				// Property access should NOT be obfuscated
				assert.Contains(t, result, ".Button.")
				assert.Contains(t, result, ".MouseButton1Click")

				t.Logf("Output: %s", result)
			},
		},
		{
			name: "Complex table constructor",
			code: `local items = getData()
table.insert(items, {text = "[Rare] Item", value = "item123"})`,
			check: func(t *testing.T, result string) {
				// Table keys should NOT be obfuscated
				assert.Contains(t, result, "{text =")
				assert.Contains(t, result, "value =")

				t.Logf("Output: %s", result)
			},
		},
		{
			name: "Discord webhook nested tables",
			code: `local content = "Hello"
local title = "Test"
local name = "Field"
local value = "Data"
local data = {
	content = content,
	embeds = {{
		title = title,
		type = 'rich',
		color = tonumber("0xfa0c0c"),
		fields = {{
			name = name,
			value = value,
			inline = false
		}}
	}}
}`,
			check: func(t *testing.T, result string) {
				// All table keys should NOT be obfuscated
				assert.Contains(t, result, "content =")
				assert.Contains(t, result, "embeds =")
				assert.Contains(t, result, "title =")
				assert.Contains(t, result, "type =")
				assert.Contains(t, result, "color =")
				assert.Contains(t, result, "fields =")
				assert.Contains(t, result, "name =")
				assert.Contains(t, result, "value =")
				assert.Contains(t, result, "inline =")

				// Variables on the left side of assignment should be obfuscated
				assert.Contains(t, result, "local _0x")

				t.Logf("Output: %s", result)
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
