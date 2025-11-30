# Lua Bundler

[![CI/CD Pipeline](https://github.com/alfin-efendy/lua-bundler/actions/workflows/ci.yml/badge.svg)](https://github.com/alfin-efendy/lua-bundler/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/alfin-efendy/lua-bundler)](https://goreportcard.com/report/github.com/alfin-efendy/lua-bundler)
[![GitHub release](https://img.shields.io/github/release/alfin-efendy/lua-bundler.svg)](https://github.com/alfin-efendy/lua-bundler/releases)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A powerful Lua script bundler specifically designed for Roblox development. Automatically resolves dependencies and bundles multiple Lua files into a single executable script.

## ✨ Features

- 🔄 **Dependency Resolution**: Automatically resolves local `require()` statements
- 🌐 **HTTP Support**: Bundles `loadstring(game:HttpGet(...))()` patterns  
- � **Smart Caching**: Automatic caching of HTTP scripts with 24-hour expiry
- �📁 **Complex Paths**: Handles relative paths, subdirectories, and parent directories
- 🚀 **Release Mode**: Removes debug statements (`print`, `warn`) for production
- 🔒 **Code Obfuscation**: 3-level obfuscation system to protect your code
- 🖥️ **HTTP Server**: Serve bundled files via HTTP for easy Roblox integration
- 🎨 **Modern CLI**: Beautiful command-line interface with Cobra and Lipgloss styling
- 🏗️ **Cross-platform**: Supports Linux, macOS, and Windows

## 📦 Installation

### Package Managers

#### Homebrew (macOS/Linux)
```bash
# Install from official tap
brew install alfin-efendy/tap/lua-bundler

# Or install from custom tap/fork
brew tap YOUR-USERNAME/tap
brew install YOUR-USERNAME/tap/lua-bundler
```

#### Winget (Windows)
```bash
# Install from official source
winget install alfin-efendy.lua-bundler

# Or install from custom manifest repository
winget install --source "YOUR-MANIFEST-REPO-URL" lua-bundler
```

#### APT (Ubuntu/Debian)
```bash
# Quick install (recommended)
curl -fsSL https://alfin-efendy.github.io/lua-bundler/install.sh | sudo bash

# Or manual installation
echo "deb [trusted=yes] https://alfin-efendy.github.io/lua-bundler/ stable main" | sudo tee /etc/apt/sources.list.d/lua-bundler.list
sudo apt update && sudo apt install lua-bundler

# Install from custom APT repository
echo "deb [trusted=yes] https://YOUR-DOMAIN/repo/ stable main" | sudo tee /etc/apt/sources.list.d/lua-bundler-custom.list
sudo apt update && sudo apt install lua-bundler
```

#### Docker
```bash
docker pull alfin-efendy/lua-bundler:latest
docker run --rm -v $(pwd):/app alfin-efendy/lua-bundler -entry /app/main.lua -output /app/bundle.lua

# Or use custom image
docker pull YOUR-USERNAME/lua-bundler:latest
docker run --rm -v $(pwd):/app YOUR-USERNAME/lua-bundler -entry /app/main.lua -output /app/bundle.lua
```

### Direct Download

Download the latest binary for your platform from [GitHub Releases](https://github.com/alfin-efendy/lua-bundler/releases).

### Build from Source

```bash
git clone https://github.com/alfin-efendy/lua-bundler.git
cd lua-bundler
make build

# Binary will be available at ./build/lua-bundler
# Optionally, move to your PATH
sudo mv build/lua-bundler /usr/local/bin/
```

### Install from Custom Source/Fork

If you want to install from a custom repository or fork:

```bash
# Install directly using Go
go install github.com/YOUR-USERNAME/lua-bundler@latest

# Or build from custom source
git clone https://github.com/YOUR-USERNAME/lua-bundler.git
cd lua-bundler
make build
sudo mv build/lua-bundler /usr/local/bin/
```

### Creating Your Own Package Repository

If you're maintaining a fork or custom version, you can create your own package repositories:

#### Homebrew Tap
```bash
# 1. Create a tap repository: homebrew-tap
# 2. Add Formula/lua-bundler.rb with your custom URL
# 3. Users install with:
brew tap YOUR-USERNAME/tap
brew install YOUR-USERNAME/tap/lua-bundler
```

**Formula Structure:**
```ruby
class LuaBundler < Formula
  desc "Your custom Lua script bundler"
  homepage "https://github.com/YOUR-USERNAME/lua-bundler"
  url "https://github.com/YOUR-USERNAME/lua-bundler/archive/refs/tags/v1.0.0.tar.gz"
  sha256 "YOUR-SHA256-HASH"
  
  depends_on "go" => :build
  
  def install
    system "make", "build"
    bin.install "build/lua-bundler"
  end
end
```

#### Winget Manifest
```bash
# 1. Fork https://github.com/microsoft/winget-pkgs
# 2. Create manifest in manifests/y/YourName/LuaBundler/
# 3. Submit PR or host your own manifest repository
# 4. Users install with custom source:
winget install --source "YOUR-MANIFEST-URL" lua-bundler
```

**Manifest Files:**
- `YourName.LuaBundler.yaml` (package metadata)
- `YourName.LuaBundler.installer.yaml` (installer info)
- `YourName.LuaBundler.locale.en-US.yaml` (localization)

#### APT Repository
```bash
# 1. Create debian package structure
# 2. Sign packages with GPG
# 3. Host repository (GitHub Pages, S3, etc.)
# 4. Users add your repo:
echo "deb [trusted=yes] https://YOUR-DOMAIN/repo/ stable main" | sudo tee /etc/apt/sources.list.d/lua-bundler-custom.list
sudo apt update && sudo apt install lua-bundler
```

**See [PACKAGING.md](PACKAGING.md) for detailed packaging guidelines.**

### Install Specific Version

```bash
# Using Go (replace v1.0.0 with desired version)
go install github.com/alfin-efendy/lua-bundler@v1.0.0

# Using Homebrew (pin specific version)
brew install lua-bundler
brew pin lua-bundler

# Download specific release from GitHub
wget https://github.com/alfin-efendy/lua-bundler/releases/download/v1.0.0/lua-bundler-linux-amd64
chmod +x lua-bundler-linux-amd64
sudo mv lua-bundler-linux-amd64 /usr/local/bin/lua-bundler
```

# Or build from custom source
git clone https://github.com/YOUR-USERNAME/lua-bundler.git
cd lua-bundler
make build
sudo mv build/lua-bundler /usr/local/bin/
```

## 🔄 Updating

### Homebrew
```bash
# Update from official tap
brew update
brew upgrade lua-bundler

# Update from custom tap
brew update
brew upgrade YOUR-USERNAME/tap/lua-bundler
```

### Winget
```bash
# Update from official source
winget upgrade alfin-efendy.lua-bundler

# Update from custom source
winget upgrade --source "YOUR-MANIFEST-URL" lua-bundler
```

### APT (Ubuntu/Debian)
```bash
# Update from official repository
sudo apt update
sudo apt upgrade lua-bundler

# Update from custom repository (if added)
# Same command works for all configured repositories
sudo apt update
sudo apt upgrade lua-bundler
```

### Go Install
```bash
# Update to latest version
go install github.com/alfin-efendy/lua-bundler@latest

# Or from custom repository
go install github.com/YOUR-USERNAME/lua-bundler@latest
```

### Manual Update (Direct Download/Build from Source)
```bash
# Download latest release from GitHub
# Or rebuild from source
cd lua-bundler
git pull origin main
make build
sudo mv build/lua-bundler /usr/local/bin/
```

### Check Current Version
```bash
lua-bundler --version
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

# Serve bundled file via HTTP (useful for Roblox game:HttpGet())
lua-bundler -e main.lua -o bundle.lua --serve --port 8080

# Show help with beautiful CLI interface
lua-bundler --help

# Check version
lua-bundler --version

# Disable cache for fresh downloads
lua-bundler -e main.lua -o bundle.lua --no-cache
```

#### Quick Reference

| Command | Description |
|---------|-------------|
| `lua-bundler -e main.lua -o out.lua` | Basic bundling |
| `lua-bundler -e main.lua -o out.lua -r` | Release mode (remove debug) |
| `lua-bundler -e main.lua -o out.lua -O 2` | With obfuscation |
| `lua-bundler -e main.lua -o out.lua -s` | Bundle and serve via HTTP |
| `lua-bundler -e main.lua -o out.lua -n` | Disable cache |
| `lua-bundler -e main.lua -o out.lua -v` | Verbose output |
| `lua-bundler --version` | Check version |
| `lua-bundler --help` | Show help |

#### CLI Flags

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--entry` | `-e` | Entry point Lua file | `main.lua` |
| `--output` | `-o` | Output bundled file | `bundle.lua` |
| `--release` | `-r` | Release mode: remove print and warn statements | `false` |
| `--obfuscate` | `-O` | Obfuscation level (0-3): 0=none, 1=basic, 2=medium, 3=heavy | `0` |
| `--verbose` | `-v` | Enable verbose output | `false` |
| `--serve` | `-s` | Start HTTP server to serve the output file | `false` |
| `--port` | `-p` | Port for HTTP server (used with --serve) | `8080` |
| `--no-cache` | `-n` | Disable HTTP cache for remote scripts | `false` |
| `--help` | `-h` | Show help information | - |

### 💾 HTTP Cache

Lua Bundler automatically caches downloaded HTTP scripts to improve build times and reduce network requests.

**Features:**
- ✅ Automatic caching of `game:HttpGet()` scripts
- ✅ Cache expiry after 24 hours
- ✅ Stored in `~/.lua-bundler-cache/`
- ✅ MD5-based cache keys for URL uniqueness

**Usage:**

```bash
# Default behavior - cache enabled
lua-bundler -e main.lua -o bundle.lua

# First run - downloads and caches HTTP scripts
# 📥 Downloading: https://example.com/script.lua

# Second run - uses cached version
# 💾 Using cached: https://example.com/script.lua

# Disable cache for always-fresh downloads
lua-bundler -e main.lua -o bundle.lua --no-cache
```

**When to use `--no-cache`:**
- 🔄 During active development when remote scripts change frequently
- 🐛 When debugging issues with remote dependencies
- ✅ When you need to ensure the latest version is fetched

### 🌍 HTTP Server

Lua Bundler includes a built-in HTTP server to serve your bundled files, making it easy to load them into Roblox using `game:HttpGet()`.

#### Basic Usage

```bash
# Bundle and serve on default port (8080)
lua-bundler -e main.lua -o bundle.lua --serve

# Bundle and serve on custom port
lua-bundler -e main.lua -o bundle.lua --serve --port 3000
```

#### Features

- 📄 **Direct File Access**: Access the bundled file directly at `http://localhost:PORT/filename.lua`
- 📋 **Directory Listing**: View all `.lua` files in the output directory at `http://localhost:PORT/`
- 🔄 **Live Serving**: Server keeps running until you stop it (Ctrl+C)
- 🌐 **CORS Enabled**: Cross-Origin Resource Sharing enabled for easy integration
- 📝 **Request Logging**: All HTTP requests are logged with timestamps

#### Using in Roblox

Once the server is running, you can load the bundled script in Roblox:

```lua
-- Load the bundled script from local HTTP server
loadstring(game:HttpGet("http://localhost:8080/bundle.lua"))()
```

**Note**: For production, you should host your bundled files on a public server. The built-in HTTP server is primarily for development and testing purposes.

### 🎯 Smart HttpGet Bundling

Lua Bundler intelligently determines which `loadstring(game:HttpGet(...))()` calls should be bundled and which should remain unchanged.

#### Bundling Behavior

**✅ BUNDLED** - Standalone HttpGet calls:
```lua
-- These are embedded into the bundle
local RemoteLib = loadstring(game:HttpGet('https://example.com/lib.lua'))()
local UI = loadstring(game:HttpGet('https://cdn.example.com/ui.lua'))()
```

**❌ NOT BUNDLED** - HttpGet inside function calls:
```lua
-- These remain unchanged (not bundled)
queue_on_teleport("loadstring(game:HttpGet('https://example.com/loader.lua'))()")
syn.queue_on_teleport("loadstring(game:HttpGet('https://cdn.example.com/script.lua'))()")
task.spawn("loadstring(game:HttpGet('https://example.com/async.lua'))()")
```

#### Why This Matters

Some Roblox functions like `queue_on_teleport` require the HttpGet URL as a string parameter to execute later. Bundling these would break the functionality since they need to dynamically load the script at execution time.

**Example - Safe Usage:**
```lua
-- This lib will be bundled
local MyLib = loadstring(game:HttpGet('https://example.com/mylib.lua'))()

-- This will NOT be bundled (stays as-is for proper teleport handling)
queue_on_teleport("loadstring(game:HttpGet('https://example.com/loader.lua'))()")

MyLib.initialize()
```

**Bundled Output:**
```lua
-- MyLib is embedded in EmbeddedModules
local MyLib = loadModule("https://example.com/mylib.lua")

-- This line remains unchanged
queue_on_teleport("loadstring(game:HttpGet('https://example.com/loader.lua'))()")

MyLib.initialize()
```

This smart detection ensures your scripts work correctly in all scenarios!

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

## ❓ FAQ / Troubleshooting

### How do I check my current version?
```bash
lua-bundler --version
```

### How do I update to the latest version?

**Homebrew:**
```bash
brew update && brew upgrade lua-bundler
```

**Winget:**
```bash
winget upgrade alfin-efendy.lua-bundler
```

**APT:**
```bash
sudo apt update && sudo apt upgrade lua-bundler
```

**Manual/Go Install:**
```bash
go install github.com/alfin-efendy/lua-bundler@latest
```

### How do I switch between official and custom package sources?

**Homebrew - Switch to custom tap:**
```bash
# Remove official version
brew uninstall lua-bundler

# Add custom tap and install
brew tap YOUR-USERNAME/tap
brew install YOUR-USERNAME/tap/lua-bundler
```

**Homebrew - Switch back to official:**
```bash
# Remove custom version
brew uninstall YOUR-USERNAME/tap/lua-bundler

# Untap custom repository
brew untap YOUR-USERNAME/tap

# Install official version
brew install alfin-efendy/tap/lua-bundler
```

**Winget - Use specific source:**
```bash
# List available sources
winget source list

# Install from specific source
winget install --source "SOURCE-NAME" lua-bundler

# Update from specific source
winget upgrade --source "SOURCE-NAME" lua-bundler
```

**APT - Multiple repositories:**
```bash
# Both official and custom repos can coexist
# Priority is determined by version numbers
# To prefer custom repo, use higher version number

# Check which version is available
apt policy lua-bundler

# Install specific version
sudo apt install lua-bundler=1.0.0-custom
```

### Cache is not working or giving errors

Clear the cache and try again:
```bash
# Remove cache directory
rm -rf ~/.lua-bundler-cache/

# Or use --no-cache flag
lua-bundler -e main.lua -o bundle.lua --no-cache
```

### HTTP downloads are failing

1. Check your internet connection
2. Try with `--no-cache` flag
3. Verify the URL is accessible
4. Check if you need a proxy configuration

### Command not found after installation

Make sure the binary is in your PATH:
```bash
# Check if it exists
which lua-bundler

# If using manual install, add to PATH
export PATH="$PATH:/usr/local/bin"

# Or for permanent, add to ~/.bashrc or ~/.zshrc
echo 'export PATH="$PATH:/usr/local/bin"' >> ~/.bashrc
```

### Permission denied when running

Make the binary executable:
```bash
chmod +x /usr/local/bin/lua-bundler
# Or for local build
chmod +x ./build/lua-bundler
```

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