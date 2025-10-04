# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- GitHub Actions CI/CD pipeline
- Multi-platform builds (Linux, macOS, Windows)
- Package publishing to Homebrew, Winget, and APT
- Docker image support
- Comprehensive test suite
- Integration tests with example bundling

### Changed
- Improved Makefile with copy functionality
- Enhanced error handling and logging

### Fixed
- Module resolution for complex relative paths
- Bundling of subdirectory modules (utils/fancy_print.lua)

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

[Unreleased]: https://github.com/alfin-efendy/lua-bundler/compare/v1.0.0...HEAD
[1.0.0]: https://github.com/alfin-efendy/lua-bundler/releases/tag/v1.0.0