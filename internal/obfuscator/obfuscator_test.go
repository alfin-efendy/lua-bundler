package obfuscator

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewObfuscator_ClampsLevel(t *testing.T) {
	assert.Equal(t, 1, NewObfuscator(0).level)
	assert.Equal(t, 3, NewObfuscator(9).level)
}

func TestObfuscate_Level1_MinifiesStringSafe(t *testing.T) {
	o := NewObfuscator(1)
	out := o.Obfuscate(`local s = "a   b" -- c` + "\nreturn s")
	assert.NotContains(t, out, "-- c")
	assert.Contains(t, out, `"a   b"`, "string spaces must survive")
}

func TestObfuscate_Level2_RenamesLocals(t *testing.T) {
	o := NewObfuscator(2)
	out := o.Obfuscate("local myVar = 1\nreturn myVar")
	assert.NotContains(t, out, "myVar")
	assert.Contains(t, out, "_0x")
}

func TestObfuscate_Level2_PreservesGlobalsAndRequire(t *testing.T) {
	o := NewObfuscator(2)
	out := o.Obfuscate(`local M = require("core/x")` + "\nreturn M.y")
	assert.Contains(t, out, `require("core/x")`)
	assert.Contains(t, out, ".y")
}

func TestObfuscate_FallsBackOnParseError(t *testing.T) {
	o := NewObfuscator(2)
	// Deliberately broken syntax: must not panic, must return minified-ish text.
	out := o.Obfuscate("local = = = 5 $$$")
	require.NotEmpty(t, out)
}

func TestObfuscate_PerUnitIsolation(t *testing.T) {
	o := NewObfuscator(2)
	a := o.Obfuscate("local Config = 1\nreturn Config")
	b := o.Obfuscate("return Config") // Config is GLOBAL here
	assert.NotContains(t, a, "Config")
	assert.Contains(t, b, "Config", "global Config in a second unit must not be renamed")
}

func TestObfuscate_DoesNotCorruptLongStrings(t *testing.T) {
	o := NewObfuscator(2)
	out := o.Obfuscate("local s = [[ it's a test ]]\nreturn s")
	assert.Contains(t, out, "[[ it's a test ]]")
}
