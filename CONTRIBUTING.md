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

1. Fork the repo and create your branch from `main`.
2. If you've added code that should be tested, add tests.
3. If you've changed APIs, update the documentation.
4. Ensure the test suite passes.
5. Make sure your code lints.
6. Issue that pull request!

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
   make deps
   ```

4. **Run the development workflow**
   ```bash
   # Format, vet, and lint
   make check
   
   # Build
   make build
   
   # Test with examples
   make example
   
   # Run tests
   make test
   ```

### Code Style

- Use `gofmt` to format your code
- Run `make check` before committing
- Follow Go best practices and conventions
- Add comments for exported functions and types

### Testing

- Write unit tests for new functionality
- Run `make test` to execute all tests
- Test with the provided examples using `make example`
- Ensure integration tests pass

### Commit Messages

- Use the present tense ("Add feature" not "Added feature")
- Use the imperative mood ("Move cursor to..." not "Moves cursor to...")
- Limit the first line to 72 characters or less
- Reference issues and pull requests liberally after the first line

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

## Release Process

1. Update version in relevant files
2. Update CHANGELOG.md
3. Create a git tag: `git tag v1.x.x`
4. Push tag: `git push origin v1.x.x`
5. Create GitHub Release
6. CI/CD will automatically build and publish packages

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