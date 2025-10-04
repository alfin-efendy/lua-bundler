package bundler

import (
	"testing"
)

func TestIsLocalModule(t *testing.T) {
	b, _ := NewBundler("test.lua", false)

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
			if got != tt.want {
				t.Errorf("isLocalModule(%q) = %v, want %v", tt.modulePath, got, tt.want)
			}
		})
	}
}

func TestResolveModulePath(t *testing.T) {
	b, _ := NewBundler("/base/main.lua", false)
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
			if got != tt.want {
				t.Errorf("resolveModulePath(%q, %q) = %v, want %v", tt.currentFile, tt.modulePath, got, tt.want)
			}
		})
	}
}
