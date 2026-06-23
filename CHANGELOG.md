# Changelog

All notable changes to this project will be documented in this file. See [Conventional Commits](https://conventionalcommits.org) for commit guidelines.

## [1.11.0](https://github.com/alfin-efendy/lua-bundler/compare/v1.10.1...v1.11.0) (2026-06-23)

### 🚀 Features

* accept Luau types, type aliases, if-expr, interpolation ([0c75c91](https://github.com/alfin-efendy/lua-bundler/commit/0c75c91d1477a7eac7307e4a0ba4c0688105dab4))
* add Luau AST node definitions ([3b7e6d1](https://github.com/alfin-efendy/lua-bundler/commit/3b7e6d13d038842e6faff64ba1254123f36338a9))
* lexical scope resolver attaching bindings to names ([0004ceb](https://github.com/alfin-efendy/lua-bundler/commit/0004cebe0aa11e4c84fc7db6119db88245cdaa7c))
* minifying AST printer with re-lex-safe spacing ([9b051b8](https://github.com/alfin-efendy/lua-bundler/commit/9b051b854c5d79fc31b1c0b304f3217e2c1da660))
* parse full Lua statement grammar + compound-assign tokens ([6efb2fc](https://github.com/alfin-efendy/lua-bundler/commit/6efb2fccf8d2017f309e2c00b91ff76007259b69))
* parse Lua expressions with precedence climbing ([2cfd2e5](https://github.com/alfin-efendy/lua-bundler/commit/2cfd2e53689e043d6ee38a6b8650abfd5ff13be4))
* scope-aware renamer for local bindings ([8c5e8be](https://github.com/alfin-efendy/lua-bundler/commit/8c5e8be12102fff1e7472ee2acfac78907414834))
* support file:// sources in downloadHTTP ([ac53074](https://github.com/alfin-efendy/lua-bundler/commit/ac530749c3b0da3a4e4b80de96afa99e6004fabc))

### 🐛 Bug Fixes

* printer emits Luau local attributes ([975fa8c](https://github.com/alfin-efendy/lua-bundler/commit/975fa8cc5af36ad1f16b457437b21e810a41b28a))

### ♻️ Code Refactoring

* drop dead Luau local-attribute code; fix log mojibake ([a9c0ad0](https://github.com/alfin-efendy/lua-bundler/commit/a9c0ad08f15b69ab2df5cd6ac1f20cddff182bb5))
* extract lexer+minify into internal/lua package ([a061ad7](https://github.com/alfin-efendy/lua-bundler/commit/a061ad79345e0ab06dee47a0df145286e101f840))
* rewire obfuscator onto internal/lua parser+renamer ([6658b99](https://github.com/alfin-efendy/lua-bundler/commit/6658b99b4759baa1e05b79223b80e26ae9cca24e))

## [1.10.1](https://github.com/alfin-efendy/lua-bundler/compare/v1.10.0...v1.10.1) (2026-06-22)

### 🐛 Bug Fixes

* add Lua lexer for string-aware minification ([7d01055](https://github.com/alfin-efendy/lua-bundler/commit/7d0105538fed41c910ab57793a59d384e7afb2cf))
* keep space between number and concat to avoid malformed number ([5b0080d](https://github.com/alfin-efendy/lua-bundler/commit/5b0080d1b49243c7d780501d36a023c385008a63))
* point example to current ez-rbx-ui release asset URL ([009b93c](https://github.com/alfin-efendy/lua-bundler/commit/009b93c099a66d3b395f4562ed53cc62b98d82b5))
* stop minifier corrupting string literals with keyword spacing ([fa777f5](https://github.com/alfin-efendy/lua-bundler/commit/fa777f544a65635f0718da38a781de6155c1a771))

### ♻️ Code Refactoring

* drop redundant removeComments now handled by minifyCode ([adcc36e](https://github.com/alfin-efendy/lua-bundler/commit/adcc36ebfa9044053436c04a649e2d24e84679ab))

## [1.10.0](https://github.com/alfin-efendy/lua-bundler/compare/v1.9.0...v1.10.0) (2026-04-03)

### 🚀 Features

* add environment variable support for Lua bundling ([d8a5161](https://github.com/alfin-efendy/lua-bundler/commit/d8a51618c05def4964d2eee6cf6462851dda340d))

## [1.9.0](https://github.com/alfin-efendy/lua-bundler/compare/v1.8.3...v1.9.0) (2025-11-30)

### 🚀 Features

* add smart HttpGet bundling logic and tests to preserve function call integrity ([552757e](https://github.com/alfin-efendy/lua-bundler/commit/552757efc672e3e2037781cffc94cea2984511e4))

## [1.8.3](https://github.com/alfin-efendy/lua-bundler/compare/v1.8.2...v1.8.3) (2025-11-14)

### 🐛 Bug Fixes

* improve obfuscation logic to preserve module names and enhance identifier handling ([b5c93a2](https://github.com/alfin-efendy/lua-bundler/commit/b5c93a213d0a3948546d6dbed6d5e3260d555940))
* update release mode checks to ensure comments and print statements are removed while preserving EmbeddedModules ([4a948e0](https://github.com/alfin-efendy/lua-bundler/commit/4a948e009641e3e1ce21ec99d0d270ce52fe6da2))

## [1.8.2](https://github.com/alfin-efendy/lua-bundler/compare/v1.8.1...v1.8.2) (2025-11-14)

### 🐛 Bug Fixes

* enhance identifier renaming to preserve require paths in obfuscation ([b9f4ba1](https://github.com/alfin-efendy/lua-bundler/commit/b9f4ba1027767bc21c387af8fadde337de0d86a3))

## [1.8.1](https://github.com/alfin-efendy/lua-bundler/compare/v1.8.0...v1.8.1) (2025-11-02)

### 🐛 Bug Fixes

* add build metadata to binaries including build date and git commit hash ([8f6e963](https://github.com/alfin-efendy/lua-bundler/commit/8f6e9637f2e9f0705f484ad4725611a74efd07ed))

### 📚 Documentation

* enhance installation instructions and add custom repository guidelines ([5d1c045](https://github.com/alfin-efendy/lua-bundler/commit/5d1c0450b358d1bdee515df66e9b843fe80f17a8))

## [1.8.0](https://github.com/alfin-efendy/lua-bundler/compare/v1.7.0...v1.8.0) (2025-11-02)

### 🚀 Features

* implement HTTP caching mechanism with 24-hour expiry and update CLI options ([ec304c2](https://github.com/alfin-efendy/lua-bundler/commit/ec304c289fa5f10eb2c53a6891ac09e9a4c18081))

## [1.7.0](https://github.com/alfin-efendy/lua-bundler/compare/v1.6.0...v1.7.0) (2025-11-02)

### 🚀 Features

* add built-in HTTP server to serve bundled files ([ecfd345](https://github.com/alfin-efendy/lua-bundler/commit/ecfd3455b81cbfe861175d3a8e6a70055e7e9009))

## [1.6.0](https://github.com/alfin-efendy/lua-bundler/compare/v1.5.0...v1.6.0) (2025-10-06)

### 🚀 Features

* Add critical executor HTTP request functions to identifier mapping in obfuscator ([043b76e](https://github.com/alfin-efendy/lua-bundler/commit/043b76ebc176ab92de75e06b2304548a3d419205))

## [1.5.0](https://github.com/alfin-efendy/lua-bundler/compare/v1.4.1...v1.5.0) (2025-10-05)

### 🚀 Features

* Add obfuscation option and improve minification in Lua bundler ([287ea05](https://github.com/alfin-efendy/lua-bundler/commit/287ea05db42c3ad8a9cefa05aa4808f2fda18f39))

## [1.4.1](https://github.com/alfin-efendy/lua-bundler/compare/v1.4.0...v1.4.1) (2025-10-05)

### 🐛 Bug Fixes

* Update APT publish workflow to handle repository_dispatch and workflow_dispatch events ([fa50a30](https://github.com/alfin-efendy/lua-bundler/commit/fa50a30521b6abfbec1257ae2d897dd41a71d52b))

## [1.4.0](https://github.com/alfin-efendy/lua-bundler/compare/v1.3.0...v1.4.0) (2025-10-05)

### 🚀 Features

* Update APT publish workflow to require release tag input ([f7e95f3](https://github.com/alfin-efendy/lua-bundler/commit/f7e95f3d3d65b314c61e6035b6cccf6a5af8b081))

## [1.3.0](https://github.com/alfin-efendy/lua-bundler/compare/v1.2.0...v1.3.0) (2025-10-05)

### 🚀 Features

* Implement Lua bundler with module handling and obfuscation ([68778cd](https://github.com/alfin-efendy/lua-bundler/commit/68778cd51a8a9cb30bb7dffb9d4709c1bb1d7b2a))

### ♻️ Code Refactoring

* Remove unused string encoding and control flow functions from obfuscator ([c801913](https://github.com/alfin-efendy/lua-bundler/commit/c801913ca64a9484885ad6acb3769c40f1d752b8))

### 📚 Documentation

* Update README.md ([9774180](https://github.com/alfin-efendy/lua-bundler/commit/9774180cb3ca34eb7c5c86d79233bd79dab8b460))

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
