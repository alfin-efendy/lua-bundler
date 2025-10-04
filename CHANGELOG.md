# Changelog

All notable changes to this project will be documented in this file. See [Conventional Commits](https://conventionalcommits.org) for commit guidelines.

## [1.2.0](https://github.com/alfin-efendy/lua-bundler/compare/v1.1.0...v1.2.0) (2025-10-04)

### 🚀 Features

* Add comprehensive test suite for Lua bundler including unit, integration, and command tests ([8bc77ce](https://github.com/alfin-efendy/lua-bundler/commit/8bc77ce4f02f298daec90bf44f078b545e244d92))
* Add pull request template and enhance test cases with improved assertions ([b455390](https://github.com/alfin-efendy/lua-bundler/commit/b45539004046ea512d54b68a52310e6b4372b69c))
* Enhance integration tests for Lua bundler with detailed checks and success messages ([ec09442](https://github.com/alfin-efendy/lua-bundler/commit/ec094426a776f4b5bfc279e6610fd2292301bb74))
* Enhance PR title validation and update contributing guidelines for commit message format ([32ffaa7](https://github.com/alfin-efendy/lua-bundler/commit/32ffaa70cf13af1611e360d864834a9a028c9ba7))
* Implement automated versioning and release process with Semantic Release ([16f0050](https://github.com/alfin-efendy/lua-bundler/commit/16f0050107ae22273785a62d1d474ce87840101e))
* Increase maximum line length limits for commit messages and bodies ([58faf37](https://github.com/alfin-efendy/lua-bundler/commit/58faf371c8db1e80f0ce7b0b744a354356ea621d))
* Update commitlint configuration to improve PR title validation and streamline commit message checks ([2b14945](https://github.com/alfin-efendy/lua-bundler/commit/2b1494561bbcccc3422b6ff6849234e923b37bcf))

### 🐛 Bug Fixes

* Update CI/CD pipeline name and enhance CI wait logic with skip option ([b34eb19](https://github.com/alfin-efendy/lua-bundler/commit/b34eb19298a57e428e7dce91e30894d250512b81))
* update permissions and token handling in semantic-release workflow ([b5e4dc8](https://github.com/alfin-efendy/lua-bundler/commit/b5e4dc8e826aa3204cc6cf1dee5dcd89d1e373ac))

# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [1.1.0] - 2025-10-04

### Added
- **Modern CLI Framework**: Implemented Cobra CLI framework for professional command-line interface
- **Beautiful Terminal Styling**: Added Lipgloss library for colored, styled output with emojis
- **Modular Architecture**: Refactored codebase into clean, separated packages for better maintainability
- **Enhanced User Experience**: Progress indicators, styled help text, and improved error messages
- **Verbose Mode**: Added `--verbose` flag for detailed processing information
- **Professional Help System**: Comprehensive help documentation with examples and feature descriptions

### Changed
- **CLI Flag Format**: Modernized to use `-e`/`--entry`, `-o`/`--output`, `--release`, `--verbose`
- **Code Organization**: Split monolithic `main.go` into modular packages:
  - `cmd/`: CLI interface and command handling with Cobra
  - `bundler/`: Core functionality split across specialized files
    - `bundler.go`: Main bundler struct and Bundle() method
    - `processor.go`: File processing, module resolution, HTTP downloading
    - `generator.go`: Bundle generation and module replacement logic
    - `utils.go`: Debug statement removal utilities
- **Dependencies**: Updated to Go 1.24 with Cobra v1.8.1 and Lipgloss v0.13.1
- **Terminal Output**: Enhanced with colors, emojis, and professional formatting
- **Documentation**: Updated README.md with new CLI interface and architecture details

### Fixed
- **Code Maintainability**: Separated concerns make the codebase easier to understand and extend
- **User Interface**: Improved clarity with styled output and better error handling
- **Development Workflow**: Updated Makefile to use new CLI flag format while maintaining compatibility

### Infrastructure
- **Package Structure**: Clean separation between CLI and core bundling logic
- **Modern Dependencies**: Professional-grade libraries for CLI and terminal styling
- **Backward Compatibility**: All existing functionality preserved while adding improvements

## [1.0.1] - 2025-10-04

### Added
- **Separated Release Pipeline**: Dedicated `release.yml` workflow for professional release management
- **Automated Package Publishing**: Automatic updates to APT, Homebrew, and Winget on release
- **CI Quality Gate**: Release pipeline now requires CI to pass before creating releases
- **Interactive Release Tool**: `scripts/create-release.sh` for guided release creation
- **Repository Dispatch System**: Orchestrated package manager updates via events
- **APT Repository**: Full Debian package repository hosted on GitHub Pages
- **Homebrew Tap**: Automatic formula updates in `alfin-efendy/homebrew-tap`
- **Winget Manifests**: Automated submission to Microsoft's winget-pkgs repository
- **Professional Documentation**: Complete pipeline documentation and troubleshooting guides

### Changed
- **Refactored CI/CD Architecture**: Separated concerns between testing (CI) and releasing
- **Enhanced Workflow Names**: Consistent naming across all workflows (`APT Publish`, etc.)
- **Improved Error Handling**: Better version detection and fallback mechanisms
- **Package Manager Integration**: Parallel publishing to all supported package managers
- **Release Process**: From manual to fully automated with quality gates

### Fixed
- **Version Detection Issues**: Enhanced version parsing for all trigger types
- **APT Repository Metadata**: Proper Release file format with correct hash entries
- **GZIP File Conflicts**: Force overwrite flags to prevent existing file errors
- **GitHub Actions Permissions**: Proper token authentication and repository access
- **Package Installation**: Resolved "unable to locate package" issues

### Infrastructure
- **Workflow Files**: 5 dedicated workflows for different aspects of CI/CD
- **Automation Scripts**: Professional tooling for release management
- **Documentation**: Comprehensive guides for development and release processes
- **Quality Assurance**: CI requirements ensure only tested code gets released

## [1.0.0] - 2025-10-04

### Added
- Initial release of lua-bundler
- Support for local file requires with relative paths
- HTTP loadstring bundling support  
- Release mode with debug statement removal
- Multi-platform binary builds
- Comprehensive Makefile with development workflow
- Example project structure
- Clipboard integration for bundled output

### Features
- Bundle Lua scripts with dependency resolution
- Support for `require()` statements with local files
- Support for `loadstring(game:HttpGet(...))()` patterns
- Automatic module embedding and replacement
- Release mode removes `print()` and `warn()` statements
- Cross-platform builds (Linux, macOS, Windows)
- Command-line interface with customizable entry and output files

[Unreleased]: https://github.com/alfin-efendy/lua-bundler/compare/v1.1.0...HEAD
[1.1.0]: https://github.com/alfin-efendy/lua-bundler/compare/v1.0.1...v1.1.0
[1.0.1]: https://github.com/alfin-efendy/lua-bundler/compare/v1.0.0...v1.0.1
[1.0.0]: https://github.com/alfin-efendy/lua-bundler/releases/tag/v1.0.0
