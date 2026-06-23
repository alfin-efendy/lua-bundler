package bundler

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ezRbxUIEntry is the ez-rbx-ui submodule entry, relative to this package dir.
const ezRbxUIEntry = "../../testdata/ez-rbx-ui/main.lua"

// findLuac returns a Lua compiler on PATH, or "" if none is available.
func findLuac() string {
	for _, name := range []string{"luac5.1", "luac"} {
		if p, err := exec.LookPath(name); err == nil {
			return p
		}
	}
	return ""
}

// TestBundle_EzRbxUI_Integration bundles the real ez-rbx-ui library and checks
// that recursive require resolution works and that --release minification does
// not corrupt string literals (the regression guard for the keyword-spacing bug).
func TestBundle_EzRbxUI_Integration(t *testing.T) {
	if _, err := os.Stat(ezRbxUIEntry); err != nil {
		t.Skip("ez-rbx-ui submodule not checked out: run git submodule update --init --recursive")
	}

	// Normal bundle: resolves all local modules recursively.
	bn, err := NewBundler(ezRbxUIEntry, false, false)
	require.NoError(t, err)
	normal, err := bn.Bundle(false)
	require.NoError(t, err)

	// These assertions track the pinned ez-rbx-ui submodule (currently 35 modules:
	// 13 core + 22 components). Bumping the submodule may require updating the
	// representative keys or the >=30 floor below.
	modules := bn.GetModules()
	for _, key := range []string{"core/theme", "core/signal", "components/button", "components/slider"} {
		assert.Contains(t, modules, key, "expected module %q to be embedded", key)
	}
	assert.GreaterOrEqual(t, len(modules), 30, "expected the real library to embed many modules")
	assert.Contains(t, normal, "EmbeddedModules")

	// Release bundle: string literals must survive minification.
	br, err := NewBundler(ezRbxUIEntry, false, false)
	require.NoError(t, err)
	release, err := br.Bundle(true)
	require.NoError(t, err)
	assert.Contains(t, release, `"function"`, "string literal must survive minification")
	assert.NotContains(t, release, `"function "`, "minifier must not inject a space inside a string literal")

	// Best-effort syntax check when a Lua compiler is available.
	if luac := findLuac(); luac != "" {
		tmp := filepath.Join(t.TempDir(), "release.lua")
		require.NoError(t, os.WriteFile(tmp, []byte(release), 0o644))
		out, cerr := exec.Command(luac, "-p", tmp).CombinedOutput()
		assert.NoErrorf(t, cerr, "release bundle failed to parse with %s: %s", luac, out)
	} else {
		t.Log("no Lua compiler (luac5.1/luac) on PATH; skipping syntax check (CI Layer B covers it)")
	}
}

// TestBundle_EzRbxUI_Obfuscated bundles the real library with rename obfuscation
// and verifies the result still parses (and runs, if the mocked-Roblox harness
// is wired). This is the regression guard for "obfuscated bundles don't run".
func TestBundle_EzRbxUI_Obfuscated(t *testing.T) {
	if _, err := os.Stat(ezRbxUIEntry); err != nil {
		t.Skip("ez-rbx-ui submodule not checked out")
	}
	b, err := NewBundler(ezRbxUIEntry, false, false)
	require.NoError(t, err)
	b.SetObfuscationLevel(2)
	out, err := b.Bundle(true)
	require.NoError(t, err)

	// Renaming must have happened somewhere.
	assert.Contains(t, out, "_0x", "expected obfuscated identifiers")

	if luac := findLuac(); luac != "" {
		tmp := filepath.Join(t.TempDir(), "obf.lua")
		require.NoError(t, os.WriteFile(tmp, []byte(out), 0o644))
		cmdOut, cerr := exec.Command(luac, "-p", tmp).CombinedOutput()
		assert.NoErrorf(t, cerr, "obfuscated bundle failed to parse: %s", cmdOut)
	} else {
		t.Log("no luac on PATH; CI Layer B / make verify-ezui covers runtime")
	}
}
