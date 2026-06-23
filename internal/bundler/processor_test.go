package bundler

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIsLocalModule(t *testing.T) {
	b, err := NewBundler("test.lua", false, false)
	require.NoError(t, err, "NewBundler should not fail")

	tests := []struct {
		name       string
		modulePath string
		want       bool
	}{
		{
			name:       "relative path with dot",
			modulePath: "./module.lua",
			want:       true,
		},
		{
			name:       "relative path with double dot",
			modulePath: "../module.lua",
			want:       true,
		},
		{
			name:       "absolute path from base",
			modulePath: "/core.lua",
			want:       true,
		},
		{
			name:       "subdirectory path",
			modulePath: "utils/helper.lua",
			want:       true,
		},
		{
			name:       "lua extension",
			modulePath: "module.lua",
			want:       true,
		},
		{
			name:       "no extension, no special chars",
			modulePath: "localmodule",
			want:       true,
		},
		{
			name:       "external module with dot",
			modulePath: "game.Workspace",
			want:       false,
		},
		{
			name:       "external module with colon",
			modulePath: "game::HttpService",
			want:       false,
		},
		{
			name:       "roblox service",
			modulePath: "ReplicatedStorage",
			want:       true, // Current implementation treats this as local
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := b.isLocalModule(tt.modulePath)
			assert.Equal(t, tt.want, got, "isLocalModule(%q) should return %v", tt.modulePath, tt.want)
		})
	}
}

func TestResolveModulePath(t *testing.T) {
	b, err := NewBundler("/base/main.lua", false, false)
	require.NoError(t, err, "NewBundler should not fail")
	b.baseDir = "/base"

	tests := []struct {
		name        string
		currentFile string
		modulePath  string
		want        string
	}{
		{
			name:        "relative path same directory",
			currentFile: "/base/main.lua",
			modulePath:  "helper",
			want:        "/base/helper.lua",
		},
		{
			name:        "relative path with extension",
			currentFile: "/base/main.lua",
			modulePath:  "helper.lua",
			want:        "/base/helper.lua",
		},
		{
			name:        "relative path subdirectory",
			currentFile: "/base/main.lua",
			modulePath:  "utils/helper",
			want:        "/base/utils/helper.lua",
		},
		{
			name:        "relative path parent directory",
			currentFile: "/base/sub/file.lua",
			modulePath:  "../core",
			want:        "/base/core.lua",
		},
		{
			name:        "absolute path from base",
			currentFile: "/base/sub/file.lua",
			modulePath:  "/core.lua",
			want:        "/base/core.lua",
		},
		{
			name:        "quoted module path",
			currentFile: "/base/main.lua",
			modulePath:  "'helper'",
			want:        "/base/helper.lua",
		},
		{
			name:        "double quoted module path",
			currentFile: "/base/main.lua",
			modulePath:  `"helper"`,
			want:        "/base/helper.lua",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := b.resolveModulePath(tt.currentFile, tt.modulePath)
			assert.Equal(t, tt.want, got, "resolveModulePath(%q, %q) should return correct path", tt.currentFile, tt.modulePath)
		})
	}
}

func TestCanonicalKey(t *testing.T) {
	dir := t.TempDir()
	// layout: <dir>/main.lua, <dir>/core/theme.lua, <dir>/components/button.lua
	require.NoError(t, os.MkdirAll(filepath.Join(dir, "core"), 0o755))
	require.NoError(t, os.MkdirAll(filepath.Join(dir, "components"), 0o755))
	for _, p := range []string{"main.lua", "core/theme.lua", "components/button.lua"} {
		require.NoError(t, os.WriteFile(filepath.Join(dir, p), []byte("return {}"), 0o644))
	}
	b, err := NewBundler(filepath.Join(dir, "main.lua"), false, false)
	require.NoError(t, err)

	entry := filepath.Join(dir, "main.lua")
	button := filepath.Join(dir, "components/button.lua")
	cases := []struct{ cur, mod, want string }{
		{entry, "core/theme", "core/theme"},
		{entry, "core/theme.lua", "core/theme"},
		{button, "../core/theme", "core/theme"},  // same module, different caller/spelling
		{button, "./sibling", "components/sibling"},
		{entry, "/core/theme", "core/theme"}, // base-relative spelling
	}
	for _, c := range cases {
		if got := b.canonicalKey(c.cur, c.mod); got != c.want {
			t.Errorf("canonicalKey(%q,%q)=%q want %q", c.cur, c.mod, got, c.want)
		}
	}
}

func TestProcessFile_KeysByCanonical(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.MkdirAll(filepath.Join(dir, "core"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "core/util.lua"), []byte("return {}"), 0o644))
	// entry requires it WITH a .lua suffix; key must canonicalize to "core/util"
	require.NoError(t, os.WriteFile(filepath.Join(dir, "main.lua"),
		[]byte(`local U = require("core/util.lua")`+"\nreturn U"), 0o644))

	b, err := NewBundler(filepath.Join(dir, "main.lua"), false, false)
	require.NoError(t, err)
	_, err = b.Bundle(false)
	require.NoError(t, err)

	mods := b.GetModules()
	if _, ok := mods["core/util"]; !ok {
		t.Fatalf("expected canonical key %q, got keys: %v", "core/util", keysOf(mods))
	}
	if _, ok := mods["core/util.lua"]; ok {
		t.Fatalf("must not key by the as-written .lua spelling")
	}
}

func keysOf(m map[string]string) []string {
	var k []string
	for key := range m {
		k = append(k, key)
	}
	return k
}
