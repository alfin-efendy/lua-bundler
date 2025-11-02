package bundler

import (
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
