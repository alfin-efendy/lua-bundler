package lua

import (
	"strings"
	"testing"
)

func obf(t *testing.T, src string) string {
	t.Helper()
	c, err := Parse(src)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	Rename(c)
	return c.Print()
}

func TestRename_LocalsRenamed(t *testing.T) {
	out := obf(t, "local myVar = 10\nreturn myVar + 1")
	if strings.Contains(out, "myVar") {
		t.Fatalf("local not renamed: %q", out)
	}
	if !strings.Contains(out, "_0x") {
		t.Fatalf("expected obfuscated name: %q", out)
	}
}

func TestRename_GlobalsAndFieldsPreserved(t *testing.T) {
	out := obf(t, "local t = game:GetService(\"Players\")\nreturn t.LocalPlayer")
	if !strings.Contains(out, "game") {
		t.Fatalf("global game must be preserved: %q", out)
	}
	if !strings.Contains(out, "GetService") {
		t.Fatalf("method name must be preserved: %q", out)
	}
	if !strings.Contains(out, ".LocalPlayer") {
		t.Fatalf("field must be preserved: %q", out)
	}
	if !strings.Contains(out, `"Players"`) {
		t.Fatalf("string must be preserved: %q", out)
	}
}

func TestRename_RequireStringPreserved(t *testing.T) {
	out := obf(t, `local M = require("core/theme")`+"\n"+`return M`)
	if !strings.Contains(out, `require("core/theme")`) {
		t.Fatalf("require path must be preserved: %q", out)
	}
}

func TestRename_TableKeysPreserved(t *testing.T) {
	out := obf(t, "local name = 1\nreturn {name = name}")
	if !strings.Contains(out, "{name=") && !strings.Contains(out, "{name =") {
		t.Fatalf("table key must be preserved: %q", out)
	}
}

func TestRename_InterpolationBailsOut(t *testing.T) {
	out := obf(t, "local n = 1\nreturn `x {n}`")
	if !strings.Contains(out, "local n") {
		t.Fatalf("must NOT rename when interpolation present: %q", out)
	}
}

func TestRename_ParamsAndForVars(t *testing.T) {
	out := obf(t, "local function f(p) for i = 1, p do end end\nreturn f")
	if strings.Contains(out, "(p)") {
		t.Fatalf("param not renamed: %q", out)
	}
}
