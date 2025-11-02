package bundler

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// downloadHTTP downloads content from HTTP URL
func (b *Bundler) downloadHTTP(url string) (string, error) {
	// Check cache first
	if b.cache.IsEnabled() {
		if content, found, err := b.cache.Get(url); err == nil && found {
			if b.verbose {
				fmt.Printf("� Using cached: %s\n", url)
			}
			return content, nil
		}
	}

	if b.verbose {
		fmt.Printf("�📥 Downloading: %s\n", url)
	}

	resp, err := b.httpClient.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to download %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to download %s: status %d", url, resp.StatusCode)
	}

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response from %s: %w", url, err)
	}

	contentStr := string(content)

	// Store in cache
	if b.cache.IsEnabled() {
		if err := b.cache.Set(url, contentStr); err != nil {
			// Log warning but don't fail
			if b.verbose {
				fmt.Printf("⚠️  Failed to cache %s: %v\n", url, err)
			}
		}
	}

	return contentStr, nil
}

// isLocalModule checks if a module path refers to a local file
func (b *Bundler) isLocalModule(modulePath string) bool {
	// Module dianggap lokal jika:
	// 1. Dimulai dengan "." (relatif)
	// 2. Dimulai dengan "/" (absolut dari base)
	// 3. Berisi "/" (subdirectory)
	// 4. Berakhir dengan ".lua"
	// 5. Tidak berisi karakter yang mengindikasikan external module
	return strings.HasPrefix(modulePath, ".") ||
		strings.HasPrefix(modulePath, "/") ||
		strings.Contains(modulePath, "/") ||
		strings.HasSuffix(modulePath, ".lua") ||
		(!strings.Contains(modulePath, ".") && !strings.Contains(modulePath, "::"))
}

// resolveModulePath resolves relative module paths to absolute paths
func (b *Bundler) resolveModulePath(currentFile, modulePath string) string {
	modulePath = strings.Trim(modulePath, "'\"")

	// Handle absolute paths from base directory
	if strings.HasPrefix(modulePath, "/") {
		resolvedPath := filepath.Join(b.baseDir, strings.TrimPrefix(modulePath, "/"))
		if !strings.HasSuffix(resolvedPath, ".lua") {
			resolvedPath += ".lua"
		}
		return resolvedPath
	}

	// Handle relative paths
	currentDir := filepath.Dir(currentFile)
	resolvedPath := filepath.Join(currentDir, modulePath)

	// Add .lua extension if not present
	if !strings.HasSuffix(resolvedPath, ".lua") {
		resolvedPath += ".lua"
	}

	// Clean the path to resolve .. and . components
	resolvedPath = filepath.Clean(resolvedPath)

	return resolvedPath
}

// processFile recursively processes a file and its dependencies
func (b *Bundler) processFile(filePath string, content string) error {
	// Regex patterns
	requireRegex := regexp.MustCompile(`require\s*\(\s*['"]([^'"]+)['"]\s*\)`)
	httpGetRegex := regexp.MustCompile(`loadstring\s*\(\s*game:HttpGet\s*\(\s*['"]([^'"]+)['"]\s*\)\s*\)\s*\(\s*\)`)

	lines := strings.Split(content, "\n")

	for _, line := range lines {
		// Check for loadstring(game:HttpGet(...))()
		if matches := httpGetRegex.FindStringSubmatch(line); len(matches) > 1 {
			url := matches[1]

			// Skip if already processed
			if _, exists := b.modules[url]; exists {
				continue
			}

			// Download content from URL
			httpContent, err := b.downloadHTTP(url)
			if err != nil {
				return err
			}

			b.modules[url] = httpContent

			// Process downloaded content (might have requires in it)
			if err := b.processFile(url, httpContent); err != nil {
				return err
			}
		}

		// Check for local require()
		if matches := requireRegex.FindStringSubmatch(line); len(matches) > 1 {
			modulePath := matches[1]

			// Process local files (relative, absolute from base, or subdirectory)
			if b.isLocalModule(modulePath) {
				resolvedPath := b.resolveModulePath(filePath, modulePath)

				// Skip if already processed
				if _, exists := b.modules[modulePath]; exists {
					continue
				}

				// Read local file
				fileContent, err := os.ReadFile(resolvedPath)
				if err != nil {
					return fmt.Errorf("failed to read file %s: %w", resolvedPath, err)
				}

				b.modules[modulePath] = string(fileContent)

				if b.verbose {
					fmt.Printf("📄 Processed: %s\n", modulePath)
				}

				// Process file recursively
				if err := b.processFile(resolvedPath, string(fileContent)); err != nil {
					return err
				}
			}
		}
	}

	return nil
}
