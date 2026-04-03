package bundler

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/alfin-efendy/lua-bundler/internal/cache"
	"github.com/alfin-efendy/lua-bundler/internal/obfuscator"
)

type Bundler struct {
	modules        map[string]string // path -> content
	httpModules    map[string]bool   // track which modules are from HTTP
	baseDir        string
	entryFile      string
	httpClient     *http.Client
	cache          *cache.Cache
	verbose        bool
	obfuscator     *obfuscator.Obfuscator
	obfuscateLevel int
	envVars        map[string]string // env var substitutions for {{VAR_NAME}}
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
		modules:     make(map[string]string),
		httpModules: make(map[string]bool),
		baseDir:     baseDir,
		entryFile:   entryFile,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		cache:          c,
		verbose:        verbose,
		obfuscateLevel: 0,
		envVars:        make(map[string]string),
	}, nil
}

// SetEnvVars sets the environment variable map used for {{VAR_NAME}} substitution.
func (b *Bundler) SetEnvVars(vars map[string]string) {
	b.envVars = vars
}

// SetObfuscationLevel sets the obfuscation level for local modules
func (b *Bundler) SetObfuscationLevel(level int) {
	b.obfuscateLevel = level
	if level > 0 {
		b.obfuscator = obfuscator.NewObfuscator(level)
	}
}

func (b *Bundler) Bundle(releaseMode bool) (string, error) {
	// Read entry file
	content, err := os.ReadFile(b.entryFile)
	if err != nil {
		return "", fmt.Errorf("failed to read entry file: %w", err)
	}

	mainContent := string(content)

	// Apply env var substitution to entry file
	mainContent = substituteEnvVars(mainContent, b.envVars, b.verbose)

	// Process all dependencies
	if b.verbose {
		fmt.Println("🔍 Processing dependencies...")
	}
	if err := b.processFile(b.entryFile, mainContent); err != nil {
		return "", err
	}

	// Obfuscate main content (entry file) if obfuscation is enabled
	if b.obfuscateLevel > 0 && b.obfuscator != nil {
		mainContent = b.obfuscator.Obfuscate(mainContent)
	}

	// Generate bundle
	bundleOutput := b.generateBundle(mainContent)

	// Apply release mode if enabled
	if releaseMode {
		if b.verbose {
			fmt.Println("🚀 Applying release mode...")
			fmt.Println("  - Removing print/warn statements...")
		}
		bundleOutput = removeDebugStatements(bundleOutput)

		if b.verbose {
			fmt.Println("  - Removing comments...")
		}
		bundleOutput = removeComments(bundleOutput)

		if b.verbose {
			fmt.Println("  - Minifying to single line...")
		}
		bundleOutput = minifyCode(bundleOutput)
	}

	return bundleOutput, nil
}

func (b *Bundler) GetModules() map[string]string {
	return b.modules
}
