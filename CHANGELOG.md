# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

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

[Unreleased]: https://github.com/alfin-efendy/lua-bundler/compare/v1.0.1...HEAD
[1.0.1]: https://github.com/alfin-efendy/lua-bundler/compare/v1.0.0...v1.0.1
[1.0.0]: https://github.com/alfin-efendy/lua-bundler/releases/tag/v1.0.0