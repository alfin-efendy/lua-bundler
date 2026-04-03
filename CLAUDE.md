# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
# Dependencies
make deps          # go mod download && go mod tidy

# Build
make build         # Build binary for current platform → ./build/lua-bundler
make build-all     # Build for Linux, Windows, macOS
make release       # Optimized release build (CGO_ENABLED=0)

# Code quality
make check         # Run fmt + vet + lint
make fmt           # go fmt ./...
make vet           # go vet ./...

# Testing
make test                               # go test -v ./...
go test -v ./internal/bundler/...       # Single package
go test -v -run TestName ./internal/... # Single test

# Run example
make run           # Build and run with example/myscript/main.lua → output/example_bundle.lua
make example       # Run with release mode + obfuscation level 3
```

## Architecture

This is a **Go CLI tool** (Cobra) that bundles Lua scripts for Roblox. It resolves local `require()` calls and remote `loadstring(game:HttpGet('url'))()` calls into a single self-contained output file.

### Processing Pipeline

```
Entry .lua file
  → processor.go   (recursive DFS: detect require/HttpGet, download, cache)
  → obfuscator.go  (optional, levels 0–3)
  → generator.go   (emit EmbeddedModules table + loadModule() + main code)
  → utils.go       (optional release mode: strip print/warn/comments, minify)
  → Output .lua file (optionally served via http/server.go)
```

### Key Modules

| Path | Role |
|------|------|
| `cmd/root.go` | CLI flags, orchestration |
| `internal/bundler/bundler.go` | Bundler struct, sets obfuscation level |
| `internal/bundler/processor.go` | Recursive dependency resolution; downloads HTTP modules |
| `internal/bundler/generator.go` | Generates the final bundle with `EmbeddedModules` |
| `internal/bundler/utils.go` | Release-mode stripping and minification |
| `internal/obfuscator/obfuscator.go` | 4-level obfuscation (comments → minify → rename vars → aggressive) |
| `internal/cache/cache.go` | HTTP cache in `~/.lua-bundler-cache/`, 24-hour expiry, MD5-keyed |
| `internal/http/server.go` | Serves bundled output over HTTP for live testing |

### Output Format

The generated bundle wraps each dependency in a named closure inside an `EmbeddedModules` table, then replaces `require()` and `loadstring(game:HttpGet())` calls with `loadModule()`:

```lua
local EmbeddedModules = {}
EmbeddedModules["ui.lua"] = function() ... end
EmbeddedModules["https://example.com/mod.lua"] = function() ... end

local function loadModule(url)
    if EmbeddedModules[url] then return EmbeddedModules[url]() end
    return require(url)
end
```

### Smart HttpGet Bundling

A critical design decision: `loadstring(game:HttpGet('url'))()` is **only bundled when it appears as a direct statement**, not when wrapped in a function call like `queue_on_teleport(...)`. The regex `\w+\s*\([^)]*loadstring\s*\(\s*game:HttpGet` detects function-wrapped calls and skips them, because those functions need the raw string to execute later (e.g., on teleport).

### Module Key Strategy

- Local modules: key = relative path as written in the `require()` call
- HTTP modules: key = the full URL string

### Obfuscation Levels

- **0**: No obfuscation
- **1**: Remove comments + minify whitespace
- **2**: Level 1 + rename local variables and functions (string-aware, preserves Lua keywords and Roblox/executor globals)
- **3**: Level 2 + aggressive single-line minification

### Release Process

Uses **Semantic Release** with **Conventional Commits**. Versioning and multi-platform binaries are produced automatically by CI. Version/build metadata is injected via LDFLAGS at build time.
