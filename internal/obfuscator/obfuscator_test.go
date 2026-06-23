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
	assert.Contains(t, out, "_0x") // the local M must be renamed to an obfuscated name
}

func TestObfuscate_FallsBackOnParseError(t *testing.T) {
	broken := "local = = = 5 $$$"
	out := NewObfuscator(2).Obfuscate(broken)
	require.NotEmpty(t, out)
	// On parse failure level 2 must fall back to the level-1 (minify) result.
	assert.Equal(t, NewObfuscator(1).Obfuscate(broken), out)
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

func TestObfuscate_Level3_EncryptsStrings(t *testing.T) {
	o := NewObfuscator(3)
	out := o.Obfuscate(`local s = "secret"` + "\nreturn s")
	assert.NotContains(t, out, `"secret"`, "string literal must be encrypted")
	assert.Contains(t, out, "_d(", "expected decoder call")
}

func TestObfuscate_Level3_PreservesRequire(t *testing.T) {
	o := NewObfuscator(3)
	out := o.Obfuscate(`local M = require("core/x")` + "\nreturn M")
	assert.Contains(t, out, `require("core/x")`, "require path must stay literal")
}

func TestObfuscate_Level3_DecoderPreludeNonEmpty(t *testing.T) {
	o := NewObfuscator(3)
	p := o.DecoderPrelude()
	assert.NotEmpty(t, p)
	assert.Contains(t, p, "_d")
	assert.NotContains(t, p, "bit32")
}

func TestObfuscate_Level2_NoDecoderPrelude(t *testing.T) {
	assert.Empty(t, NewObfuscator(2).DecoderPrelude())
	out := NewObfuscator(2).Obfuscate(`local s = "x"` + "\nreturn s")
	assert.NotContains(t, out, "_d(", "level 2 must not encrypt strings")
}
