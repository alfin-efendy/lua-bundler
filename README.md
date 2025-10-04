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
- 🔧 **CLI Interface**: Simple command-line interface with customizable options
- 📋 **Clipboard Integration**: Auto-copy bundled output to clipboard
- 🏗️ **Cross-platform**: Supports Linux, macOS, and Windows

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
lua-bundler -entry main.lua -output bundle.lua

# Bundle in release mode (removes debug statements)
lua-bundler -entry main.lua -output bundle.lua -release
```

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
local EzUI = loadstring(game:HttpGet('https://example.com/ui.lua'))()

-- Your main code here
```

**Command:**
```bash
lua-bundler -entry main.lua -output bundle.lua -release
```

**Result:** A single `bundle.lua` file with all dependencies embedded and debug statements removed.

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
[![APT Repository](https://github.com/alfin-efendy/lua-bundler/actions/workflows/apt-repository.yml/badge.svg)](https://github.com/alfin-efendy/lua-bundler/actions/workflows/apt-repository.yml)

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

- Built for the Roblox development community
- Inspired by modern bundling tools like Webpack and Rollup
- Thanks to all contributors and users!