package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMain_Integration(t *testing.T) {
	// Build the binary first
	binaryName := "lua-bundler-test"
	err := exec.Command("go", "build", "-o", binaryName, ".").Run()
	require.NoError(t, err, "Failed to build binary")
	defer os.Remove(binaryName)

	// Create temporary test directory
	tempDir, err := os.MkdirTemp("", "integration-test")
	require.NoError(t, err, "Failed to create temp dir")
	defer os.RemoveAll(tempDir)

	// Create test files
	mainFile := filepath.Join(tempDir, "main.lua")
	moduleFile := filepath.Join(tempDir, "helper.lua")
	outputFile := filepath.Join(tempDir, "bundle.lua")

	mainContent := `-- Main script
local helper = require('./helper.lua')
print("Starting application")
helper.greet("World")
print("Application finished")`

	moduleContent := `-- Helper module
local helper = {}

function helper.greet(name)
    print("Hello, " .. name .. "!")
end

return helper`

	err = os.WriteFile(mainFile, []byte(mainContent), 0644)
	require.NoError(t, err, "Failed to write main file")

	err = os.WriteFile(moduleFile, []byte(moduleContent), 0644)
	require.NoError(t, err, "Failed to write module file")

	tests := []struct {
		name        string
		args        []string
		expectError bool
		checkOutput func(t *testing.T, output string)
		checkFile   func(t *testing.T, filePath string)
	}{
		{
			name: "help command",
			args: []string{"--help"},
			checkOutput: func(t *testing.T, output string) {
				assert.Contains(t, output, "lua-bundler", "Help output should contain 'lua-bundler'")
				assert.Contains(t, output, "Usage:", "Help output should contain 'Usage:'")
			},
		},
		{
			name: "basic bundling",
			args: []string{"-e", mainFile, "-o", outputFile},
			checkOutput: func(t *testing.T, output string) {
				assert.Contains(t, output, "Successfully bundled", "Output should contain success message")
				assert.Contains(t, output, "Modules embedded: 1", "Output should show 1 embedded module")
			},
			checkFile: func(t *testing.T, filePath string) {
				content, err := os.ReadFile(filePath)
				require.NoError(t, err, "Should be able to read output file")

				bundleContent := string(content)
				assert.Contains(t, bundleContent, "Bundled Lua Script", "Bundle should contain header")
				assert.Contains(t, bundleContent, "EmbeddedModules", "Bundle should contain EmbeddedModules")
				assert.Contains(t, bundleContent, "loadModule", "Bundle should contain loadModule function")
			},
		},
		{
			name: "verbose bundling",
			args: []string{"-e", mainFile, "-o", outputFile, "--verbose"},
			checkOutput: func(t *testing.T, output string) {
				assert.Contains(t, output, "Verbose: Enabled", "Verbose mode should be indicated in output")
				assert.Contains(t, output, "Successfully bundled", "Output should contain success message")
			},
		},
		{
			name: "release mode bundling",
			args: []string{"-e", mainFile, "-o", outputFile, "--release"},
			checkOutput: func(t *testing.T, output string) {
				assert.Contains(t, output, "Release (debug statements removed)", "Release mode should be indicated in output")
			},
			checkFile: func(t *testing.T, filePath string) {
				content, err := os.ReadFile(filePath)
				require.NoError(t, err, "Should be able to read output file")

				bundleContent := string(content)
				// Should have fewer print statements in release mode
				originalPrintCount := strings.Count(mainContent+moduleContent, "print(")
				bundlePrintCount := strings.Count(bundleContent, "print(")
				assert.Less(t, bundlePrintCount, originalPrintCount, "Release mode should remove some print statements")
			},
		},
		{
			name:        "nonexistent file",
			args:        []string{"-e", "nonexistent.lua", "-o", outputFile},
			expectError: true,
			checkOutput: func(t *testing.T, output string) {
				assert.Contains(t, output, "Bundling failed", "Should show bundling failed message")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clean up output file between tests
			os.Remove(outputFile)

			// Run the command
			cmd := exec.Command("./"+binaryName, tt.args...)
			output, err := cmd.CombinedOutput()

			// Check if error expectation matches
			if tt.expectError {
				assert.Error(t, err, "Expected error but command succeeded")
			} else {
				assert.NoError(t, err, "Unexpected error")
			}

			// Check output if function provided
			if tt.checkOutput != nil {
				tt.checkOutput(t, string(output))
			}

			// Check output file if function provided and file should exist
			if tt.checkFile != nil && !tt.expectError {
				assert.FileExists(t, outputFile, "Output file should exist")
				tt.checkFile(t, outputFile)
			}
		})
	}
}

func TestMain_Version(t *testing.T) {
	// Test version injection (this is a compile-time feature)
	// We can't easily test this in unit tests, but we can verify
	// the binary builds correctly with version flags

	binaryName := "lua-bundler-version-test"
	version := "test-version"

	// Build with version
	err := exec.Command("go", "build",
		"-ldflags", "-X main.Version="+version,
		"-o", binaryName, ".").Run()

	if err != nil {
		t.Fatalf("Failed to build binary with version: %v", err)
	}
	defer os.Remove(binaryName)

	// The binary should build successfully
	// Version testing would require more complex setup to verify
	// the version is actually set correctly
}
