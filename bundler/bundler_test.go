package bundler

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNewBundler(t *testing.T) {
	tests := []struct {
		name      string
		entryFile string
		verbose   bool
		wantErr   bool
	}{
		{
			name:      "valid entry file",
			entryFile: "test.lua",
			verbose:   false,
			wantErr:   false,
		},
		{
			name:      "verbose mode",
			entryFile: "test.lua",
			verbose:   true,
			wantErr:   false,
		},
		{
			name:      "entry file in subdirectory",
			entryFile: "subdir/test.lua",
			verbose:   false,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := NewBundler(tt.entryFile, tt.verbose)

			if tt.wantErr {
				if err == nil {
					t.Errorf("NewBundler() expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("NewBundler() unexpected error: %v", err)
				return
			}

			if b == nil {
				t.Errorf("NewBundler() returned nil bundler")
				return
			}

			if b.entryFile != tt.entryFile {
				t.Errorf("NewBundler() entryFile = %v, want %v", b.entryFile, tt.entryFile)
			}

			if b.verbose != tt.verbose {
				t.Errorf("NewBundler() verbose = %v, want %v", b.verbose, tt.verbose)
			}

			if b.modules == nil {
				t.Errorf("NewBundler() modules map is nil")
			}

			if b.httpClient == nil {
				t.Errorf("NewBundler() httpClient is nil")
			}
		})
	}
}

func TestBundle(t *testing.T) {
	// Create temporary test files
	tempDir, err := os.MkdirTemp("", "bundler-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create test files
	mainFile := filepath.Join(tempDir, "main.lua")
	moduleFile := filepath.Join(tempDir, "module.lua")

	mainContent := `-- Main script
local module = require('./module.lua')
print("Hello from main")
module.greet()
`

	moduleContent := `-- Module script
local m = {}

function m.greet()
    print("Hello from module")
end

return m
`

	if err := os.WriteFile(mainFile, []byte(mainContent), 0644); err != nil {
		t.Fatalf("Failed to write main file: %v", err)
	}

	if err := os.WriteFile(moduleFile, []byte(moduleContent), 0644); err != nil {
		t.Fatalf("Failed to write module file: %v", err)
	}

	tests := []struct {
		name        string
		entryFile   string
		release     bool
		wantErr     bool
		checkOutput func(string) bool
	}{
		{
			name:      "basic bundling",
			entryFile: mainFile,
			release:   false,
			wantErr:   false,
			checkOutput: func(output string) bool {
				return strings.Contains(output, "Bundled Lua Script") &&
					strings.Contains(output, "EmbeddedModules") &&
					strings.Contains(output, "loadModule")
			},
		},
		{
			name:      "release mode bundling",
			entryFile: mainFile,
			release:   true,
			wantErr:   false,
			checkOutput: func(output string) bool {
				return strings.Contains(output, "Bundled Lua Script") &&
					!strings.Contains(output, "print(") // Should remove print statements
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := NewBundler(tt.entryFile, false)
			if err != nil {
				t.Fatalf("NewBundler() failed: %v", err)
			}

			result, err := b.Bundle(tt.release)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Bundle() expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("Bundle() unexpected error: %v", err)
				return
			}

			if result == "" {
				t.Errorf("Bundle() returned empty result")
				return
			}

			if tt.checkOutput != nil && !tt.checkOutput(result) {
				t.Errorf("Bundle() output validation failed")
			}
		})
	}
}

func TestBundle_NonexistentFile(t *testing.T) {
	b, err := NewBundler("nonexistent.lua", false)
	if err != nil {
		t.Fatalf("NewBundler() failed: %v", err)
	}

	_, err = b.Bundle(false)
	if err == nil {
		t.Errorf("Bundle() expected error for nonexistent file, got nil")
	}
}

func TestGetModules(t *testing.T) {
	b, err := NewBundler("test.lua", false)
	if err != nil {
		t.Fatalf("NewBundler() failed: %v", err)
	}

	// Initially should be empty
	modules := b.GetModules()
	if len(modules) != 0 {
		t.Errorf("GetModules() expected empty map, got %d items", len(modules))
	}

	// Add a module manually for testing
	b.modules["test"] = "content"

	modules = b.GetModules()
	if len(modules) != 1 {
		t.Errorf("GetModules() expected 1 item, got %d", len(modules))
	}

	if modules["test"] != "content" {
		t.Errorf("GetModules() expected 'content', got %v", modules["test"])
	}
}
