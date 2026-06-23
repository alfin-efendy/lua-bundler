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
