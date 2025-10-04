# Contributing to Lua Bundler

We love your input! We want to make contributing to lua-bundler as easy and transparent as possible, whether it's:

- Reporting a bug
- Discussing the current state of the code
- Submitting a fix
- Proposing new features
- Becoming a maintainer

## Development Process

We use GitHub to host code, to track issues and feature requests, as well as accept pull requests.

### Pull Requests

1. **Fork the repo** and create your branch from `main`
2. **Use conventional commits** for all commit messages
3. **Add tests** if you've added code that should be tested
4. **Update documentation** if you've changed APIs
5. **Ensure test suite passes**: `make test`
6. **Make sure code lints**: `make check`
7. **Create PR** with conventional commit format in title

#### PR Title Format
Your PR title must follow conventional commits format:
```bash
feat: add new bundling feature
fix: resolve memory leak in processor  
docs: update installation guide
```

#### Automated Checks
When you create a PR, the following validations run automatically:
- ✅ **Commit message format** validation
- ✅ **PR title format** validation  
- ✅ **Test suite** execution
- ✅ **Code quality** checks (lint, vet, format)
- ✅ **Release preview** (shows what version would be released)

#### What Happens After Merge
- 🤖 **Semantic release** analyzes your commits
- 📈 **Version bump** calculated automatically
- 📋 **Changelog** updated with your changes
- 🚀 **Release created** if version bump is needed
- 📦 **Packages published** to all platforms

### Development Setup

1. **Clone the repository**
   ```bash
   git clone https://github.com/alfin-efendy/lua-bundler.git
   cd lua-bundler
   ```

2. **Install Go 1.24+**
   - Download from [golang.org](https://golang.org/dl/)

3. **Install development dependencies**
   ```bash
   # Go dependencies
   go mod tidy
   
   # Optional: Node.js for commit helpers
   npm install  # For commitizen and semantic-release tools
   ```

4. **Run the development workflow**
   ```bash
   # Format, vet, and lint
   make check
   
   # Build binary
   make build
   
   # Run comprehensive tests (includes Testify)
   make test
   
   # Test with examples
   make example
   
   # Check test coverage
   go test -cover ./...
   ```

5. **Optional: Use commit helpers**
   ```bash
   # Interactive conventional commit creator
   npm run commit
   
   # Or just use git with conventional format
   git commit -m "feat: your new feature"
   ```

### Code Style

- Use `gofmt` to format your code
- Run `make check` before committing
- Follow Go best practices and conventions
- Add comments for exported functions and types

### Testing

This project uses **[Testify](https://github.com/stretchr/testify)** for comprehensive unit testing with 78.7% coverage.

#### Test Requirements
- **Write unit tests** for new functionality using Testify assertions
- **Run tests locally**: `make test` or `go test ./...`
- **Test with examples**: `make example`
- **Check coverage**: `go test -cover ./...`
- **Integration tests** must pass

#### Test Structure
```go
func TestNewFeature(t *testing.T) {
    // Use testify assertions
    assert.Equal(t, expected, actual, "descriptive message")
    require.NoError(t, err, "should not return error")
    assert.Contains(t, result, "expected content")
}
```

#### Test Categories
- **Unit tests**: `*_test.go` files alongside source code
- **Integration tests**: `main_test.go` with real binary execution
- **CLI tests**: `cmd/root_test.go` for command-line interface
- **Coverage target**: Maintain >75% test coverage

### Commit Messages

**Important:** This project uses [Conventional Commits](https://conventionalcommits.org/) for automated versioning and releases.

All commit messages must follow this format:
```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

#### Commit Types & Version Impact

| Type | Description | Version Bump | Example |
|------|-------------|--------------|---------|
| `feat` | New feature | **Minor** (1.1.0) | `feat: add HTTP module support` |
| `fix` | Bug fix | **Patch** (1.0.1) | `fix: resolve bundling crash` |
| `perf` | Performance improvement | **Patch** (1.0.1) | `perf: optimize module loading` |
| `feat!` | Breaking change | **Major** (2.0.0) | `feat!: change CLI argument format` |
| `fix!` | Breaking bug fix | **Major** (2.0.0) | `fix!: remove deprecated API` |
| `docs` | Documentation only | **No release** | `docs: update README examples` |
| `test` | Adding/updating tests | **No release** | `test: add integration tests` |
| `chore` | Maintenance tasks | **No release** | `chore: update dependencies` |
| `ci` | CI/CD changes | **No release** | `ci: add automated testing` |
| `build` | Build system changes | **No release** | `build: update Makefile` |
| `refactor` | Code refactoring | **Patch** (1.0.1) | `refactor: simplify bundler logic` |
| `style` | Code style changes | **No release** | `style: fix formatting` |
| `revert` | Revert previous commit | **Patch** (1.0.1) | `revert: undo breaking change` |

#### Commit Rules
- **Type**: Must be lowercase
- **Description**: Must not start with uppercase, no period at end
- **Header**: Maximum 100 characters
- **Body**: Wrap at 100 characters per line
- **Breaking Changes**: Use `!` after type or `BREAKING CHANGE:` in footer

#### Good Examples ✅
```bash
feat: add support for nested module dependencies
fix: resolve crash when entry file not found
feat!: change CLI flags to use single dash format
docs: add installation instructions for macOS
test: add comprehensive bundler unit tests
perf: optimize module resolution algorithm
chore: update Go version to 1.24
```

#### Bad Examples ❌
```bash
Add feature              # Missing type
feat Add new feature     # Missing colon
Fix bug                  # Not descriptive enough
FEAT: add feature        # Wrong case
feat: Add feature        # Description starts with uppercase
feat: add feature.       # Ends with period
```

#### Validation
- Commits are validated automatically on PRs
- PR titles must also follow conventional format
- Use `npm run commit` for guided commit creation (optional)

### Issue and Bug Reports

**Great Bug Reports** tend to have:

- A quick summary and/or background
- Steps to reproduce
  - Be specific!
  - Give sample code if you can
- What you expected would happen
- What actually happens
- Notes (possibly including why you think this might be happening, or stuff you tried that didn't work)

### Feature Requests

We track feature requests as GitHub issues. Provide:

- **Clear title and description**
- **Use case**: Explain why this feature would be useful
- **Proposed solution**: If you have ideas on implementation
- **Alternatives**: Any alternative solutions you've considered

## Automated Workflows

This project includes several automated GitHub Actions workflows:

### On Pull Request
- **[validate-commits.yml](.github/workflows/validate-commits.yml)**: Validates commit messages and PR title
- **[ci.yml](.github/workflows/ci.yml)**: Runs tests, builds, and quality checks
- **Release preview**: Shows what version would be released

### On Merge to Main  
- **[semantic-release.yml](.github/workflows/semantic-release.yml)**: Automated versioning and release
- **Multi-platform builds**: Linux, macOS, Windows (amd64/arm64)
- **Package publishing**: Homebrew, APT, Winget

### Quality Gates
- ✅ All tests must pass
- ✅ Code must be properly formatted (`gofmt`)
- ✅ No linting errors (`golangci-lint`)
- ✅ Conventional commit format enforced
- ✅ Minimum test coverage maintained

## Release Process

**🎉 Fully Automated!** This project uses [Semantic Release](https://semantic-release.gitbook.io/) for completely automated versioning and publishing.

### How It Works
1. **Make changes** using conventional commit messages
2. **Create Pull Request** → Automated validation runs
3. **Merge to main** → Release automatically triggered if needed
4. **Version calculated** automatically from commit types
5. **Changelog generated** from commit messages
6. **Git tag created** and **GitHub release** published
7. **Binaries built** for all platforms (Linux, macOS, Windows)
8. **Package managers updated** (Homebrew, APT, Winget)

### Release Types Triggered By Commits
- **Major Release (2.0.0)**: `feat!:`, `fix!:`, or `BREAKING CHANGE:`
- **Minor Release (1.1.0)**: `feat:` commits
- **Patch Release (1.0.1)**: `fix:`, `perf:`, `refactor:` commits  
- **No Release**: `docs:`, `test:`, `chore:`, `ci:`, `build:`, `style:`

### Manual Release (Emergency Only)
```bash
# Trigger workflow manually in GitHub Actions
gh workflow run semantic-release.yml

# Or create manual tag (not recommended)
git tag v1.2.3
git push origin v1.2.3
```

### What Gets Generated
- 🏷️ **Git tag** (e.g., `v1.2.0`)
- 📋 **GitHub release** with auto-generated changelog
- 📦 **Cross-platform binaries** (attached to release)
- 📝 **Updated CHANGELOG.md**
- 🔄 **Package manager updates** (Homebrew, APT, Winget)

See **[VERSIONING.md](VERSIONING.md)** for detailed documentation.

## Code of Conduct

### Our Pledge

In the interest of fostering an open and welcoming environment, we as contributors and maintainers pledge to making participation in our project and our community a harassment-free experience for everyone.

### Our Standards

Examples of behavior that contributes to creating a positive environment include:

- Using welcoming and inclusive language
- Being respectful of differing viewpoints and experiences
- Gracefully accepting constructive criticism
- Focusing on what is best for the community
- Showing empathy towards other community members

### Enforcement

Instances of abusive, harassing, or otherwise unacceptable behavior may be reported by contacting the project team. All complaints will be reviewed and investigated and will result in a response that is deemed necessary and appropriate to the circumstances.

## License

By contributing, you agree that your contributions will be licensed under the same license as the project (MIT License).