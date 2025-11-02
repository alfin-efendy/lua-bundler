package bundler

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/alfin-efendy/lua-bundler/internal/cache"
)

type Bundler struct {
	modules    map[string]string // path -> content
	baseDir    string
	entryFile  string
	httpClient *http.Client
	cache      *cache.Cache
	verbose    bool
}

func NewBundler(entryFile string, verbose bool, useCache bool) (*Bundler, error) {
	baseDir := filepath.Dir(entryFile)
	if baseDir == "." {
		var err error
		baseDir, err = os.Getwd()
		if err != nil {
			return nil, fmt.Errorf("failed to get working directory: %w", err)
		}
	}

	// Initialize cache
	c, err := cache.NewCache(useCache)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize cache: %w", err)
	}

	return &Bundler{
		modules:   make(map[string]string),
		baseDir:   baseDir,
		entryFile: entryFile,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		cache:   c,
		verbose: verbose,
	}, nil
}

func (b *Bundler) Bundle(releaseMode bool) (string, error) {
	// Read entry file
	content, err := os.ReadFile(b.entryFile)
	if err != nil {
		return "", fmt.Errorf("failed to read entry file: %w", err)
	}

	mainContent := string(content)

	// Process all dependencies
	if b.verbose {
		fmt.Println("🔍 Processing dependencies...")
	}
	if err := b.processFile(b.entryFile, mainContent); err != nil {
		return "", err
	}

	// Generate bundle
	bundleOutput := b.generateBundle(mainContent)

	// Apply release mode if enabled
	if releaseMode {
		if b.verbose {
			fmt.Println("🚀 Applying release mode (removing print/warn statements)...")
		}
		bundleOutput = removeDebugStatements(bundleOutput)
	}

	return bundleOutput, nil
}

func (b *Bundler) GetModules() map[string]string {
	return b.modules
}
