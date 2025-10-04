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
# Add repository
curl -fsSL https://apt.example.com/gpg | sudo apt-key add -
echo "deb https://apt.example.com stable main" | sudo tee /etc/apt/sources.list.d/lua-bundler.list

# Install
sudo apt update
sudo apt install lua-bundler
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

- **[API Documentation](docs/API.md)** - Detailed usage and options
- **[Examples](example/)** - Sample projects and use cases  
- **[Contributing Guide](CONTRIBUTING.md)** - How to contribute
- **[Changelog](CHANGELOG.md)** - Release notes and changes

## 🏗️ CI/CD

This project uses GitHub Actions for:
- ✅ **Continuous Integration**: Tests on every push/PR
- 🏗️ **Multi-platform Builds**: Linux, macOS, Windows (AMD64 + ARM64)  
- 📦 **Automated Publishing**: Homebrew, Winget, APT, Docker
- 🔄 **Release Automation**: Automatic releases on git tags

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Built for the Roblox development community
- Inspired by modern bundling tools like Webpack and Rollup
- Thanks to all contributors and users!