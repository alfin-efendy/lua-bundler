# Lua Bundler

[![CI/CD Pipeline](https://github.com/alfin-efendy/lua-bundler/actions/workflows/ci.yml/badge.svg)](https://github.com/alfin-efendy/lua-bundler/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/alfin-efendy/lua-bundler)](https://goreportcard.com/report/github.com/alfin-efendy/lua-bundler)
[![GitHub release](https://img.shields.io/github/release/alfin-efendy/lua-bundler.svg)](https://github.com/alfin-efendy/lua-bundler/releases)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A powerful Lua script bundler specifically designed for Roblox development. Automatically resolves dependencies and bundles multiple Lua files into a single executable script.

## ✨ Features

- 🔄 **Dependency Resolution**: Automatically resolves local `require()` statements
- 🌐 **HTTP Support**: Bundles `loadstring(game:HttpGet(...))()` patterns  
- 📁 **Complex Paths**: Handles relative paths, subdirectories, and parent directories
- 🚀 **Release Mode**: Removes debug statements (`print`, `warn`) for production
- 🔒 **Code Obfuscation**: 3-level obfuscation system to protect your code
- 🎨 **Modern CLI**: Beautiful command-line interface with Cobra and Lipgloss styling
- 🏗️ **Cross-platform**: Supports Linux, macOS, 

## 📦 Installation

### Package Managers

#### Homebrew (macOS/Linux)
```bash
brew install alfin-efendy/tap/lua-bundler
```

#### Winget (Windows)
```bash
winget install alfin-efendy.lua-bundler
```

#### APT (Ubuntu/Debian)
```bash
# Quick install (recommended)
curl -fsSL https://alfin-efendy.github.io/lua-bundler/install.sh | sudo bash

# Or manual installation
echo "deb [trusted=yes] https://alfin-efendy.github.io/lua-bundler/ stable main" | sudo tee /etc/apt/sources.list.d/lua-bundler.list
sudo apt update && sudo apt install lua-bundler
```

#### Docker
```bash
docker pull alfin-efendy/lua-bundler:latest
docker run --rm -v $(pwd):/app alfin-efendy/lua-bundler -entry /app/main.lua -output /app/bundle.lua
```

### Direct Download

Download the latest binary for your platform from [GitHub Releases](https://github.com/alfin-efendy/lua-bundler/releases).

### Build from Source

```bash
git clone https://github.com/alfin-efendy/lua-bundler.git
cd lua-bundler
make build
```

## 🚀 Quick Start

### Basic Usage

```bash
# Bundle with default settings
lua-bundler -e main.lua -o bundle.lua

# Bundle in release mode (removes debug statements)  
lua-bundler -e main.lua -o bundle.lua --release

# Bundle with code obfuscation (level 2 - medium)
lua-bundler -e main.lua -o bundle.lua --obfuscate 2

# Bundle with release mode AND heavy obfuscation
lua-bundler -e main.lua -o bundle.lua --release --obfuscate 3

# Enable verbose output for debugging
lua-bundler -e main.lua -o bundle.lua --verbose

# Show help with beautiful CLI interface
lua-bundler --help
```

#### CLI Flags

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--entry` | `-e` | Entry point Lua file | `main.lua` |
| `--output` | `-o` | Output bundled file | `bundle.lua` |
| `--release` | `-r` | Release mode: remove print and warn statements | `false` |
| `--obfuscate` | `-O` | Obfuscation level (0-3): 0=none, 1=basic, 2=medium, 3=heavy | `0` |
| `--verbose` | `-v` | Enable verbose output | `false` |
| `--help` | `-h` | Show help information | - |

### 🔒 Code Obfuscation

Lua Bundler includes a powerful 3-level obfuscation system to protect your code:

#### Obfuscation Levels

| Level | Name | Description | Features |
|-------|------|-------------|----------|
| **0** | None | No obfuscation (default) | Original readable code |
| **1** | Basic | Light obfuscation | • Removes all comments<br>• Minifies whitespace<br>• Keeps code structure |
| **2** | Medium | Moderate protection | • All Level 1 features<br>• Renames local variables<br>• Renames functions<br>• Preserves string literals |
| **3** | Heavy | Maximum protection | • All Level 2 features<br>• Aggressive minification<br>• Single-line output<br>• Minimal size |

#### Obfuscation Examples

**Original Code:**
```lua
-- This is a greeting function
local function greet(name)
    local message = "Hello, " .. name
    print(message)
    return message
end

local userName = "World"
greet(userName)
```

**Level 1 (Basic):**
```lua
local function greet(name)
local message="Hello, "..name
print(message)
return message
end
local userName="World"
greet(userName)
```

**Level 2 (Medium):**
```lua
local function _0x4a2f8c(_0x1b3e9d)
local _0x5c7a2e="Hello, ".._0x1b3e9d
print(_0x5c7a2e)
return _0x5c7a2e
end
local _0x8f1d4b="World"
_0x4a2f8c(_0x8f1d4b)
```

**Level 3 (Heavy):**
```lua
local function _0x4a2f8c(_0x1b3e9d) local _0x5c7a2e="Hello, ".._0x1b3e9d print(_0x5c7a2e) return _0x5c7a2e end local _0x8f1d4b="World" _0x4a2f8c(_0x8f1d4b)
```

#### String Preservation

The obfuscator is **string-aware** and preserves all string literals:
- ✅ Service names: `game:GetService("HttpService")`
- ✅ Remote event names: `game:GetService("ReplicatedStorage"):WaitForChild("RemoteEvent")`
- ✅ All quoted strings remain intact
- ✅ No breaking of game functionality

#### Usage Examples

```bash
# Basic obfuscation (comments removed, whitespace minified)
lua-bundler -e main.lua -o bundle.lua -O 1

# Medium obfuscation (+ identifier renaming)
lua-bundler -e main.lua -o bundle.lua -O 2

# Heavy obfuscation (+ single-line minification)
lua-bundler -e main.lua -o bundle.lua -O 3

# Combine with release mode for maximum optimization
lua-bundler -e main.lua -o bundle.lua --release --obfuscate 3
```

#### When to Use Obfuscation

| Use Case | Recommended Level |
|----------|-------------------|
| Open-source projects | 0 (None) |
| Development/Testing | 0-1 |
| Private projects | 1-2 |
| Commercial products | 2-3 |
| Premium scripts | 3 (Heavy) |

> **Note:** Obfuscation is not encryption. It makes code harder to read but doesn't provide complete security. Always use server-side validation for critical logic.

### Using Makefile (Development)

```bash
# Build and run example
make example

# Build, run, and copy to clipboard  
make example-copy

# Custom entry file
make run ENTRY_FILE=src/main.lua OUTPUT_FILE=output/my_bundle.lua

# Development workflow with copy
make run-copy ENTRY_FILE=src/main.lua
```

## 📋 Example

Given this project structure:
```
project/
├── main.lua
├── ui.lua
├── utils/
│   └── fancy_print.lua
└── core.lua
```

**main.lua:**
```lua
local UI = require('ui.lua')
local FancyPrint = require('utils/fancy_print.lua') 
local Core = require('../core.lua')
local EzUI = loadstring(game:HttpGet('https://raw.githubusercontent.com/alfin-efendy/ez-rbx-ui/refs/heads/main/ui.lua'))()

-- Your main code here
```

**Command:**
```bash
# Basic bundle
lua-bundler -entry main.lua -output bundle.lua

# With release mode and obfuscation
lua-bundler -entry main.lua -output bundle.lua --release --obfuscate 2
```

**Result:** A single `bundle.lua` file with all dependencies embedded, debug statements removed, and code obfuscated for protection.

## 🛠️ Development

### Prerequisites
- Go 1.24+
- Make (optional, for convenience)

### Development Workflow

```bash
# Install dependencies
make deps

# Format, vet, and lint
make check

# Build
make build

# Test with examples
make example

# Run tests
make test

# Development with auto-copy
make run-copy ENTRY_FILE=your_script.lua
```

### Contributing

We welcome contributions! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## 🏷️ Automated Versioning & Releases

This project uses **Semantic Release** with **Conventional Commits** for fully automated versioning and releases.

### ✨ Key Features
- 🤖 **Automatic version calculation** based on commit messages
- 📋 **Generated changelogs** from conventional commits  
- 🏷️ **Git tags** and **GitHub releases** created automatically
- 📦 **Multi-platform binaries** built and attached to releases
- 🔄 **Package managers** updated automatically (Homebrew, APT, Winget)

### 📝 Commit Format
All commits must follow [Conventional Commits](https://conventionalcommits.org/):

```bash
feat: add new bundling option     # Minor release (1.1.0)
fix: resolve bundling crash       # Patch release (1.0.1)  
feat!: change CLI argument format # Major release (2.0.0)
docs: update README examples      # No release
test: add integration tests       # No release
chore: update dependencies        # No release
```

### 🚀 Release Process
1. **Make changes** using conventional commit messages
2. **Create PR** → Automated validation runs
3. **Merge to main** → Release automatically triggered
4. **New version** tagged and published with binaries
5. **Package managers** updated automatically

See **[VERSIONING.md](VERSIONING.md)** for detailed documentation.

## 📖 Documentation

- **[Examples](example/)** - Sample projects and use cases  
- **[Contributing Guide](CONTRIBUTING.md)** - How to contribute
- **[Changelog](CHANGELOG.md)** - Release notes and changes

## 🚀 Development

### Building from Source
```bash
# Clone the repository
git clone https://github.com/alfin-efendy/lua-bundler.git
cd lua-bundler

# Install dependencies
go mod download

# Build for your platform
make build

# Run tests
make test

# Bundle example and test
make example
```

### Creating Releases

```bash
# Use the automated release script (recommended)
./scripts/create-release.sh

# Or manually:
# 1. Update CHANGELOG.md with new version
# 2. Create and push tag
git tag v1.0.1
git push origin v1.0.1
```

**Release Process:**
1. 📝 **Update CHANGELOG.md** - Document changes for the new version
2. 🏷️ **Create Tag** - Use semantic versioning (v1.0.0, v1.1.0, v2.0.0)
3. 🔄 **CI Quality Gate** - Release waits for CI pipeline to pass
4. 🚀 **Automated Publishing** - All package managers update automatically

The release pipeline automatically:
1. ✅ **Builds** binaries for all platforms
2. 📦 **Creates** GitHub release with assets
3. 🔄 **Updates** all package managers:
   - APT repository (immediate)
   - Homebrew tap (immediate)
   - Winget (requires approval)

## 🏗️ CI/CD Pipeline

[![Release Pipeline](https://github.com/alfin-efendy/lua-bundler/actions/workflows/release.yml/badge.svg)](https://github.com/alfin-efendy/lua-bundler/actions/workflows/release.yml)
[![APT Repository](https://github.com/alfin-efendy/lua-bundler/actions/workflows/apt-publish.yml/badge.svg)](https://github.com/alfin-efendy/lua-bundler/actions/workflows/apt-publish.yml)

**Separated Workflows:**
- 🧪 **[CI Pipeline](/.github/workflows/ci.yml)**: Testing and quality checks
- � **[Release Pipeline](/.github/workflows/release.yml)**: Automated releases  
- 📦 **[Package Publishers](/.github/workflows/)**: APT, Homebrew, Winget

**Features:**
- ✅ Multi-platform builds (Linux, macOS, Windows)
- 🔄 Automatic package manager updates
- � Comprehensive testing and integration
- 🎯 One-command releases

See **[Release Pipeline Documentation](RELEASE_PIPELINE.md)** for detailed information.

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [Cobra](https://github.com/spf13/cobra) - CLI interactions
- [Lipgloss](https://github.com/charmbracelet/lipgloss) - Style terminal layouts
- [Testify](https://github.com/stretchr/testify) - Unit test lib
- Built for the Roblox development community
- Inspired by modern bundling tools like Webpack and Rollup
- Thanks to all contributors and users!