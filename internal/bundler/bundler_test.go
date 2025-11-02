package bundler

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
			b, err := NewBundler(tt.entryFile, tt.verbose, false)

			if tt.wantErr {
				assert.Error(t, err, "NewBundler() should return error")
				return
			}

			require.NoError(t, err, "NewBundler() should not return error")
			require.NotNil(t, b, "NewBundler() should return non-nil bundler")

			assert.Equal(t, tt.entryFile, b.entryFile, "entryFile should match")
			assert.Equal(t, tt.verbose, b.verbose, "verbose flag should match")
			assert.NotNil(t, b.modules, "modules map should not be nil")
			assert.NotNil(t, b.httpClient, "httpClient should not be nil")
		})
	}
}

func TestBundle(t *testing.T) {
	// Create temporary test files
	tempDir, err := os.MkdirTemp("", "bundler-test")
	require.NoError(t, err, "Failed to create temp dir")
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

	err = os.WriteFile(mainFile, []byte(mainContent), 0644)
	require.NoError(t, err, "Failed to write main file")

	err = os.WriteFile(moduleFile, []byte(moduleContent), 0644)
	require.NoError(t, err, "Failed to write module file")

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
			b, err := NewBundler(tt.entryFile, false, false)
			require.NoError(t, err, "NewBundler() should not fail")

			result, err := b.Bundle(tt.release)

			if tt.wantErr {
				assert.Error(t, err, "Bundle() should return error")
				return
			}

			require.NoError(t, err, "Bundle() should not return error")
			assert.NotEmpty(t, result, "Bundle() should return non-empty result")

			if tt.checkOutput != nil {
				assert.True(t, tt.checkOutput(result), "Bundle() output validation should pass")
			}
		})
	}
}

func TestBundle_NonexistentFile(t *testing.T) {
	b, err := NewBundler("nonexistent.lua", false, false)
	require.NoError(t, err, "NewBundler() should not fail")

	_, err = b.Bundle(false)
	assert.Error(t, err, "Bundle() should return error for nonexistent file")
}

func TestGetModules(t *testing.T) {
	b, err := NewBundler("test.lua", false, false)
	require.NoError(t, err, "NewBundler() should not fail")

	// Initially should be empty
	modules := b.GetModules()
	assert.Empty(t, modules, "GetModules() should return empty map initially")

	// Add a module manually for testing
	b.modules["test"] = "content"
	modules = b.GetModules()

	assert.Len(t, modules, 1, "GetModules() should return map with 1 item")
	assert.Equal(t, "content", modules["test"], "GetModules() should return correct content")
}
