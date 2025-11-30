package bundler

import (
	"strings"
	"testing"
)

func TestQueueOnTeleportNotBundled(t *testing.T) {
	tests := []struct {
		name           string
		content        string
		shouldBeBundle bool
		description    string
	}{
		{
			name: "queue_on_teleport should not be bundled",
			content: `queue_on_teleport("loadstring(game:HttpGet('https://example.com/script.lua'))()")
print("test")`,
			shouldBeBundle: false,
			description:    "HttpGet inside queue_on_teleport should remain unchanged",
		},
		{
			name: "syn.queue_on_teleport should not be bundled",
			content: `syn.queue_on_teleport("loadstring(game:HttpGet('https://example.com/script.lua'))()")
print("test")`,
			shouldBeBundle: false,
			description:    "HttpGet inside syn.queue_on_teleport should remain unchanged",
		},
		{
			name: "standalone loadstring should be bundled",
			content: `local lib = loadstring(game:HttpGet('https://example.com/lib.lua'))()
print("test")`,
			shouldBeBundle: true,
			description:    "Standalone HttpGet should be converted to loadModule",
		},
		{
			name:           "any_function_call with HttpGet should not be bundled",
			content:        `some_function("loadstring(game:HttpGet('https://example.com/script.lua'))()")`,
			shouldBeBundle: false,
			description:    "HttpGet inside any function call (single line) should remain unchanged",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Bundler{
				modules:     make(map[string]string),
				httpModules: make(map[string]bool),
				baseDir:     "/tmp",
				verbose:     false,
			}

			result := b.replaceModuleCalls(tt.content)

			// Check if HttpGet pattern was replaced with loadModule
			containsLoadModule := strings.Contains(result, `loadModule("https://example.com/`)
			containsHttpGet := strings.Contains(result, "game:HttpGet")

			if tt.shouldBeBundle {
				// Should be replaced with loadModule
				if !containsLoadModule {
					t.Errorf("%s: Expected HttpGet to be replaced with loadModule, but it wasn't", tt.description)
				}
			} else {
				// Should remain as HttpGet (not bundled)
				if !containsHttpGet {
					t.Errorf("%s: Expected HttpGet to remain unchanged, but it was replaced", tt.description)
				}
				if containsLoadModule {
					t.Errorf("%s: Expected HttpGet NOT to be replaced with loadModule, but it was", tt.description)
				}
			}
		})
	}
}

func TestQueueOnTeleportVariations(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "queue_on_teleport with string argument",
			input: `queue_on_teleport("loadstring(game:HttpGet('https://cdn.example.com/loader.lua'))()")`,
			want:  `queue_on_teleport("loadstring(game:HttpGet('https://cdn.example.com/loader.lua'))()")`,
		},
		{
			name:  "syn.queue_on_teleport with string argument",
			input: `syn.queue_on_teleport("loadstring(game:HttpGet('https://cdn.example.com/loader.lua'))()")`,
			want:  `syn.queue_on_teleport("loadstring(game:HttpGet('https://cdn.example.com/loader.lua'))()")`,
		},
		{
			name:  "standalone HttpGet should be replaced",
			input: `local module = loadstring(game:HttpGet('https://cdn.example.com/lib.lua'))()`,
			want:  `local module = loadModule("https://cdn.example.com/lib.lua")`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Bundler{
				modules:     make(map[string]string),
				httpModules: make(map[string]bool),
				baseDir:     "/tmp",
				verbose:     false,
			}

			got := b.replaceModuleCalls(tt.input)
			if got != tt.want {
				t.Errorf("replaceModuleCalls() got = %v, want %v", got, tt.want)
			}
		})
	}
}
