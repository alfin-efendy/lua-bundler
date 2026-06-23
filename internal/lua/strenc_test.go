package lua

import "testing"

func TestUnquoteLuaString(t *testing.T) {
	cases := []struct {
		in   string
		want string
		ok   bool
	}{
		{`"abc"`, "abc", true},
		{`'a\tb'`, "a\tb", true},
		{`"a\nb"`, "a\nb", true},
		{`"quote\"x"`, `quote"x`, true},
		{`"\065\066"`, "AB", true},
		{`"\x41\x42"`, "AB", true},
		{"[[raw]]", "raw", true},
		{"[==[a]==]", "a", true},
		{"[[\nleading]]", "leading", true}, // first newline dropped
		{"`interp{x}`", "", false},          // backtick not handled
		{`"bad\q"`, "", false},              // unknown escape
	}
	for _, c := range cases {
		got, ok := unquoteLuaString(c.in)
		if ok != c.ok || (ok && got != c.want) {
			t.Errorf("unquoteLuaString(%q) = (%q,%v), want (%q,%v)", c.in, got, ok, c.want, c.ok)
		}
	}
}

func TestEncodeStringRoundTrip(t *testing.T) {
	// encodeString emits "\ddd..." of value[i]^key; decoding those escapes and
	// XOR-ing by key again must recover value.
	key := byte(0x5a)
	value := "Hello\tWorld\000\255"
	lit := encodeString(value, key)
	decoded, ok := unquoteLuaString(lit)
	if !ok {
		t.Fatalf("encodeString produced un-decodable literal: %q", lit)
	}
	if len(decoded) != len(value) {
		t.Fatalf("length mismatch: %d vs %d", len(decoded), len(value))
	}
	var got []byte
	for i := 0; i < len(decoded); i++ {
		got = append(got, decoded[i]^key)
	}
	if string(got) != value {
		t.Fatalf("round-trip failed: %q != %q", string(got), value)
	}
}

func TestDecoderPrelude_NoBit32(t *testing.T) {
	p := DecoderPrelude(137)
	if p == "" {
		t.Fatal("expected non-empty prelude")
	}
	// Must define _d and must NOT use bit32 (lua5.1 harness lacks it).
	if !contains(p, "_d") {
		t.Fatalf("prelude must define _d: %q", p)
	}
	if contains(p, "bit32") {
		t.Fatalf("prelude must not use bit32: %q", p)
	}
	if !contains(p, "137") {
		t.Fatalf("prelude must embed the key: %q", p)
	}
}

func contains(s, sub string) bool {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}

func encStr(t *testing.T, src string, key byte) string {
	t.Helper()
	c, err := Parse(src)
	if err != nil {
		t.Fatalf("parse %q: %v", src, err)
	}
	EncryptStrings(c, key)
	return c.Print()
}

func TestEncryptStrings_PlainEncrypted(t *testing.T) {
	out := encStr(t, `local s = "hello"`, 0x33)
	if contains(out, `"hello"`) {
		t.Fatalf("plaintext leaked: %q", out)
	}
	if !contains(out, "_d(") {
		t.Fatalf("expected _d() decoder call: %q", out)
	}
}

func TestEncryptStrings_RequireExcluded(t *testing.T) {
	out := encStr(t, `local M = require("core/theme")`, 0x33)
	if !contains(out, `require("core/theme")`) {
		t.Fatalf("require path must NOT be encrypted: %q", out)
	}
}

func TestEncryptStrings_HttpGetExcluded(t *testing.T) {
	out := encStr(t, `loadstring(game:HttpGet("https://x.test/a.lua"))()`, 0x33)
	if !contains(out, `"https://x.test/a.lua"`) {
		t.Fatalf("HttpGet URL must NOT be encrypted: %q", out)
	}
}

func TestEncryptStrings_GetServiceEncrypted(t *testing.T) {
	// Non-protected calls: string args ARE encrypted (runtime decodes them).
	out := encStr(t, `local p = game:GetService("Players")`, 0x33)
	if contains(out, `"Players"`) {
		t.Fatalf("GetService arg should be encrypted: %q", out)
	}
	if !contains(out, "_d(") {
		t.Fatalf("expected _d() call: %q", out)
	}
}

func TestEncryptStrings_BacktickLeftAlone(t *testing.T) {
	out := encStr(t, "local s = `hi {x}`", 0x33)
	if !contains(out, "`hi {x}`") {
		t.Fatalf("interp string must be left as-is: %q", out)
	}
}

func TestEncryptStrings_ParenthesizedRequireExcluded(t *testing.T) {
	out := encStr(t, `local M = require(("core/theme"))`, 0x33)
	if !contains(out, `"core/theme"`) {
		t.Fatalf("parenthesized require path must NOT be encrypted: %q", out)
	}
}

func TestEncryptStrings_ParenthesizedHttpGetExcluded(t *testing.T) {
	out := encStr(t, `loadstring(game:HttpGet(("https://x.test/a.lua")))()`, 0x33)
	if !contains(out, `"https://x.test/a.lua"`) {
		t.Fatalf("parenthesized HttpGet URL must NOT be encrypted: %q", out)
	}
}

func TestEncryptStrings_LoadModuleExcluded(t *testing.T) {
	out := encStr(t, `local M = loadModule("core/theme")`, 0x33)
	if !contains(out, `"core/theme"`) {
		t.Fatalf("loadModule key must NOT be encrypted: %q", out)
	}
}
