package bundler

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/joho/godotenv"
)

var envVarRegex = regexp.MustCompile(`\{\{([A-Za-z_][A-Za-z0-9_]*)\}\}`)

// substituteEnvVars replaces {{VAR_NAME}} placeholders with values from envVars.
// Missing variables are left as-is and optionally warned about.
func substituteEnvVars(content string, envVars map[string]string, verbose bool) string {
	return envVarRegex.ReplaceAllStringFunc(content, func(match string) string {
		varName := envVarRegex.FindStringSubmatch(match)[1]
		if val, ok := envVars[varName]; ok {
			return val
		}
		if verbose {
			fmt.Printf("⚠️  Env var not found: %s\n", varName)
		}
		return match
	})
}

// loadEnvFile loads variables from a .env file.
// If the file does not exist, returns an empty map (silent no-op).
func loadEnvFile(envFilePath string) (map[string]string, error) {
	if envFilePath == "" {
		envFilePath = ".env"
	}
	vars, err := godotenv.Read(envFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return make(map[string]string), nil
		}
		return nil, err
	}
	return vars, nil
}

// BuildEnvVars loads a .env file and merges with OS env vars (OS wins).
// Exported so cmd package can call it.
func BuildEnvVars(envFilePath string) (map[string]string, error) {
	fileVars, err := loadEnvFile(envFilePath)
	if err != nil {
		return nil, err
	}
	result := fileVars
	for _, entry := range os.Environ() {
		parts := strings.SplitN(entry, "=", 2)
		if len(parts) == 2 {
			result[parts[0]] = parts[1]
		}
	}
	return result, nil
}
