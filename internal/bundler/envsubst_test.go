package bundler

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSubstituteEnvVars(t *testing.T) {
	tests := []struct {
		name    string
		content string
		envVars map[string]string
		want    string
	}{
		{
			name:    "simple substitution",
			content: `Hello {{NAME}}`,
			envVars: map[string]string{"NAME": "World"},
			want:    `Hello World`,
		},
		{
			name:    "multiple placeholders",
			content: `{{A}} and {{B}}`,
			envVars: map[string]string{"A": "foo", "B": "bar"},
			want:    `foo and bar`,
		},
		{
			name:    "missing var left as-is",
			content: `key={{MISSING}}`,
			envVars: map[string]string{},
			want:    `key={{MISSING}}`,
		},
		{
			name:    "uppercase var name",
			content: `{{API_KEY}}`,
			envVars: map[string]string{"API_KEY": "secret"},
			want:    `secret`,
		},
		{
			name:    "lowercase var name",
			content: `{{api_key}}`,
			envVars: map[string]string{"api_key": "val"},
			want:    `val`,
		},
		{
			name:    "mixed case var name",
			content: `{{MyVar}}`,
			envVars: map[string]string{"MyVar": "test"},
			want:    `test`,
		},
		{
			name:    "var with numbers",
			content: `{{VAR_123}}`,
			envVars: map[string]string{"VAR_123": "num"},
			want:    `num`,
		},
		{
			name:    "no placeholders",
			content: `plain lua code`,
			envVars: map[string]string{"X": "y"},
			want:    `plain lua code`,
		},
		{
			name:    "empty content",
			content: ``,
			envVars: map[string]string{"X": "y"},
			want:    ``,
		},
		{
			name:    "empty envVars map",
			content: `{{X}}`,
			envVars: map[string]string{},
			want:    `{{X}}`,
		},
		{
			name:    "value with special chars (URL)",
			content: `{{URL}}`,
			envVars: map[string]string{"URL": "https://a.b/c?d=e&f=g"},
			want:    `https://a.b/c?d=e&f=g`,
		},
		{
			name:    "multiline content",
			content: "line1\n{{VAR}}\nline3",
			envVars: map[string]string{"VAR": "line2"},
			want:    "line1\nline2\nline3",
		},
		{
			name:    "consecutive placeholders",
			content: `{{A}}{{B}}`,
			envVars: map[string]string{"A": "x", "B": "y"},
			want:    `xy`,
		},
		{
			name:    "single brace not matched",
			content: `{VAR}`,
			envVars: map[string]string{"VAR": "val"},
			want:    `{VAR}`,
		},
		{
			name:    "digit-start not matched",
			content: `{{1VAR}}`,
			envVars: map[string]string{"1VAR": "val"},
			want:    `{{1VAR}}`,
		},
		{
			name:    "value containing braces (no double substitution)",
			content: `{{X}}`,
			envVars: map[string]string{"X": "{{Y}}", "Y": "should-not-appear"},
			want:    `{{Y}}`,
		},
		{
			name:    "lua string with placeholder",
			content: `local key = "{{API_KEY}}"`,
			envVars: map[string]string{"API_KEY": "my-secret"},
			want:    `local key = "my-secret"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := substituteEnvVars(tt.content, tt.envVars, false)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSubstituteEnvVars_VerboseWarning(t *testing.T) {
	// Just verify it doesn't panic and still returns the placeholder
	result := substituteEnvVars("{{MISSING}}", map[string]string{}, true)
	assert.Equal(t, "{{MISSING}}", result)
}

func TestLoadEnvFile(t *testing.T) {
	t.Run("existing .env file", func(t *testing.T) {
		dir := t.TempDir()
		envPath := filepath.Join(dir, ".env")
		err := os.WriteFile(envPath, []byte("KEY=value\nOTHER=123\n"), 0644)
		require.NoError(t, err)

		vars, err := loadEnvFile(envPath)
		require.NoError(t, err)
		assert.Equal(t, "value", vars["KEY"])
		assert.Equal(t, "123", vars["OTHER"])
	})

	t.Run("nonexistent file returns empty map silently", func(t *testing.T) {
		vars, err := loadEnvFile("/tmp/this-file-should-not-exist-lua-bundler.env")
		require.NoError(t, err)
		assert.Empty(t, vars)
	})

	t.Run("empty path defaults to .env and silently skips if absent", func(t *testing.T) {
		// Change to a temp dir with no .env file
		orig, _ := os.Getwd()
		dir := t.TempDir()
		require.NoError(t, os.Chdir(dir))
		t.Cleanup(func() { os.Chdir(orig) })

		vars, err := loadEnvFile("")
		require.NoError(t, err)
		assert.Empty(t, vars)
	})
}

func TestBuildEnvVars_OSOverridesFile(t *testing.T) {
	dir := t.TempDir()
	envPath := filepath.Join(dir, ".env")
	err := os.WriteFile(envPath, []byte("MY_TEST_VAR=from_file\n"), 0644)
	require.NoError(t, err)

	t.Setenv("MY_TEST_VAR", "from_os")

	vars, err := BuildEnvVars(envPath)
	require.NoError(t, err)
	assert.Equal(t, "from_os", vars["MY_TEST_VAR"])
}

func TestBundleWithEnvVars(t *testing.T) {
	dir := t.TempDir()
	mainLua := filepath.Join(dir, "main.lua")
	err := os.WriteFile(mainLua, []byte(`local key = "{{API_KEY}}"`+"\n"), 0644)
	require.NoError(t, err)

	b, err := NewBundler(mainLua, false, false)
	require.NoError(t, err)

	b.SetEnvVars(map[string]string{"API_KEY": "my-secret-key"})

	result, err := b.Bundle(false)
	require.NoError(t, err)
	assert.Contains(t, result, "my-secret-key")
	assert.NotContains(t, result, "{{API_KEY}}")
}
