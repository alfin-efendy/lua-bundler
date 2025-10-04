package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestMain_Integration(t *testing.T) {
	// Build the binary first
	binaryName := "lua-bundler-test"
	if err := exec.Command("go", "build", "-o", binaryName, ".").Run(); err != nil {
		t.Fatalf("Failed to build binary: %v", err)
	}
	defer os.Remove(binaryName)

	// Create temporary test directory
	tempDir, err := os.MkdirTemp("", "integration-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
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

	if err := os.WriteFile(mainFile, []byte(mainContent), 0644); err != nil {
		t.Fatalf("Failed to write main file: %v", err)
	}

	if err := os.WriteFile(moduleFile, []byte(moduleContent), 0644); err != nil {
		t.Fatalf("Failed to write module file: %v", err)
	}

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
				if !strings.Contains(output, "lua-bundler") {
					t.Errorf("Help output should contain 'lua-bundler'")
				}
				if !strings.Contains(output, "Usage:") {
					t.Errorf("Help output should contain 'Usage:'")
				}
			},
		},
		{
			name: "basic bundling",
			args: []string{"-e", mainFile, "-o", outputFile},
			checkOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "Successfully bundled") {
					t.Errorf("Output should contain success message")
				}
				if !strings.Contains(output, "Modules embedded: 1") {
					t.Errorf("Output should show 1 embedded module")
				}
			},
			checkFile: func(t *testing.T, filePath string) {
				content, err := os.ReadFile(filePath)
				if err != nil {
					t.Errorf("Failed to read output file: %v", err)
					return
				}

				bundleContent := string(content)
				if !strings.Contains(bundleContent, "Bundled Lua Script") {
					t.Errorf("Bundle should contain header")
				}
				if !strings.Contains(bundleContent, "EmbeddedModules") {
					t.Errorf("Bundle should contain EmbeddedModules")
				}
				if !strings.Contains(bundleContent, "loadModule") {
					t.Errorf("Bundle should contain loadModule function")
				}
			},
		},
		{
			name: "verbose bundling",
			args: []string{"-e", mainFile, "-o", outputFile, "--verbose"},
			checkOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "Verbose: Enabled") {
					t.Errorf("Verbose mode should be indicated in output")
				}
				if !strings.Contains(output, "Successfully bundled") {
					t.Errorf("Output should contain success message")
				}
			},
		},
		{
			name: "release mode bundling",
			args: []string{"-e", mainFile, "-o", outputFile, "--release"},
			checkOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "Release (debug statements removed)") {
					t.Errorf("Release mode should be indicated in output")
				}
			},
			checkFile: func(t *testing.T, filePath string) {
				content, err := os.ReadFile(filePath)
				if err != nil {
					t.Errorf("Failed to read output file: %v", err)
					return
				}

				bundleContent := string(content)
				// Should have fewer print statements in release mode
				if strings.Count(bundleContent, "print(") >= strings.Count(mainContent+moduleContent, "print(") {
					t.Errorf("Release mode should remove some print statements")
				}
			},
		},
		{
			name:        "nonexistent file",
			args:        []string{"-e", "nonexistent.lua", "-o", outputFile},
			expectError: true,
			checkOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "Bundling failed") {
					t.Errorf("Should show bundling failed message")
				}
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
			if tt.expectError && err == nil {
				t.Errorf("Expected error but command succeeded")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v\nOutput: %s", err, output)
			}

			// Check output if function provided
			if tt.checkOutput != nil {
				tt.checkOutput(t, string(output))
			}

			// Check output file if function provided and file should exist
			if tt.checkFile != nil && !tt.expectError {
				if _, err := os.Stat(outputFile); os.IsNotExist(err) {
					t.Errorf("Output file should exist but doesn't")
				} else {
					tt.checkFile(t, outputFile)
				}
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
