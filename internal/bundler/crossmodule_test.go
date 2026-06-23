package bundler

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func findLua() string {
	for _, n := range []string{"lua5.1", "lua"} {
		if p, err := exec.LookPath(n); err == nil {
			return p
		}
	}
	return ""
}

func TestBundle_RemoteModuleBodyRewritten(t *testing.T) {
	dir := t.TempDir()
	write := func(p, s string) {
		require.NoError(t, os.WriteFile(filepath.Join(dir, p), []byte(s), 0o644))
	}
	leaf := "file://" + filepath.Join(dir, "leaf.lua")
	remote := "file://" + filepath.Join(dir, "remote.lua")
	write("leaf.lua", `return { tag = "LEAF" }`)
	// remote module's BODY pulls the leaf via HttpGet — must be rewritten to loadModule.
	write("remote.lua", `local L = loadstring(game:HttpGet("`+leaf+`"))()
return { leaf = L.tag }`)
	write("main.lua", `local R = loadstring(game:HttpGet("`+remote+`"))()
print("remote.leaf="..R.leaf)`)

	b, err := NewBundler(filepath.Join(dir, "main.lua"), false, false)
	require.NoError(t, err)
	out, err := b.Bundle(false)
	require.NoError(t, err)

	// Both remote and leaf embedded (discovery), and the remote BODY must call loadModule(leaf), not raw HttpGet.
	assert.Contains(t, b.GetModules(), remote)
	assert.Contains(t, b.GetModules(), leaf)
	assert.Contains(t, out, `loadModule("`+leaf+`")`, "remote module body's HttpGet must be rewritten to loadModule")

	if lua := findLua(); lua != "" {
		bp := filepath.Join(dir, "bundle.lua")
		require.NoError(t, os.WriteFile(bp, []byte(out), 0o644))
		got, rerr := exec.Command(lua, bp).CombinedOutput()
		require.NoErrorf(t, rerr, "bundle failed to run: %s", got)
		assert.Contains(t, string(got), "remote.leaf=LEAF",
			"remote module must resolve its leaf via the embedded loadModule (got: %s)", got)
	}
}

func TestBundle_CrossModuleRequire_RunsAndMemoizes(t *testing.T) {
	dir := t.TempDir()
	write := func(p, s string) {
		require.NoError(t, os.MkdirAll(filepath.Dir(filepath.Join(dir, p)), 0o755))
		require.NoError(t, os.WriteFile(filepath.Join(dir, p), []byte(s), 0o644))
	}
	// b.lua: shared module with a load counter; entry does NOT require it directly.
	write("b.lua", `_G.bloads = (_G.bloads or 0) + 1
return { n = _G.bloads }`)
	// a.lua and c.lua both require b (cross-module).
	write("a.lua", `local B = require("b")
return { v = function() return B.n end }`)
	write("c.lua", `local B = require("b")
return { v = function() return B.n end }`)
	// entry requires a and c only; prints whether both see the same single load.
	write("main.lua", `local A = require("a")
local C = require("c")
print("A="..A.v().." C="..C.v().." loads=".._G.bloads)`)

	b, err := NewBundler(filepath.Join(dir, "main.lua"), false, false)
	require.NoError(t, err)
	out, err := b.Bundle(false)
	require.NoError(t, err)

	// b.lua must be embedded once under canonical key "b" even though only A/C require it.
	assert.Contains(t, b.GetModules(), "b", "shared module must be discovered via cross-module require")

	lua := findLua()
	if lua == "" {
		t.Skip("no lua5.1/lua on PATH; cannot run cross-module bundle")
	}
	bundlePath := filepath.Join(dir, "bundle.lua")
	require.NoError(t, os.WriteFile(bundlePath, []byte(out), 0o644))
	got, rerr := exec.Command(lua, bundlePath).CombinedOutput()
	require.NoErrorf(t, rerr, "cross-module bundle failed to run: %s", got)
	// Memoization: b loaded once → A and C both see n=1, loads=1.
	assert.Contains(t, string(got), "A=1 C=1 loads=1",
		"shared module must load exactly once (got: %s)", got)
}
