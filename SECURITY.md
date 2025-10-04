# Security Policy

## Supported Versions

We support the following versions of lua-bundler with security updates:

| Version | Supported          |
| ------- | ------------------ |
| 1.x.x   | :white_check_mark: |
| < 1.0   | :x:                |

## Reporting a Vulnerability

We take the security of lua-bundler seriously. If you believe you have found a security vulnerability, please report it to us as described below.

### How to Report

**Please do not report security vulnerabilities through public GitHub issues.**

Instead, please report them via email to: **security@[your-domain].com** or create a private security advisory on GitHub.

### What to Include

Please include the following information in your report:

- **Type of issue** (e.g. buffer overflow, SQL injection, cross-site scripting, etc.)
- **Full paths of source file(s)** related to the manifestation of the issue
- **The location of the affected source code** (tag/branch/commit or direct URL)
- **Any special configuration required** to reproduce the issue
- **Step-by-step instructions to reproduce the issue**
- **Proof-of-concept or exploit code** (if possible)
- **Impact of the issue**, including how an attacker might exploit the issue

### Response Timeline

- **Initial response**: Within 48 hours
- **Detailed response**: Within 1 week
- **Fix timeline**: Depends on complexity, but we aim for 2-4 weeks

### Process

1. **Receipt acknowledgment**: We'll acknowledge receipt of your vulnerability report within 48 hours
2. **Investigation**: We'll investigate and confirm the vulnerability
3. **Fix development**: We'll develop and test a fix
4. **Release**: We'll release the security fix
5. **Disclosure**: We'll publicly disclose the vulnerability after the fix is available

### Recognition

We maintain a security hall of fame for responsible disclosure. With your permission, we'll:

- Credit you in our security advisories
- Add you to our hall of fame
- Potentially offer a bounty (case by case basis)

## Security Best Practices

When using lua-bundler:

### For Users
- Always download from official sources (GitHub releases, official package managers)
- Verify checksums of downloaded binaries
- Keep lua-bundler updated to the latest version
- Review bundled output before deploying to production

### For Developers
- Validate all input files before processing
- Use the latest Go version for security patches
- Review dependencies regularly
- Follow secure coding practices

## Security Features

lua-bundler includes these security considerations:

- **Input validation**: All file paths and URLs are validated
- **Sandboxing**: No arbitrary code execution during bundling
- **Path traversal protection**: Prevents directory traversal attacks
- **Memory safety**: Built with Go's memory-safe runtime
- **No external dependencies**: Minimal attack surface

## Known Security Considerations

### HTTP Downloads
- lua-bundler can download and execute Lua code from HTTP URLs
- This is intentional behavior for Roblox loadstring patterns
- Users should only bundle from trusted sources
- Consider using HTTPS URLs when possible

### File System Access
- lua-bundler reads local files during bundling
- It respects file permissions and doesn't modify source files
- Output files are created with standard permissions (0644)

## Vulnerability Disclosure Policy

We follow responsible disclosure principles:

1. **Coordinated disclosure**: We work with reporters to ensure proper fixes
2. **Public disclosure**: After fixes are available and users have time to update
3. **CVE assignment**: For significant vulnerabilities when appropriate
4. **Credit**: We provide credit to security researchers who follow responsible disclosure

## Contact

- **Email**: security@[your-domain].com
- **GitHub Security Advisories**: [Create private advisory](https://github.com/alfin-efendy/lua-bundler/security/advisories/new)
- **GPG Key**: Available on request for encrypted communication

Thank you for helping keep lua-bundler and its users safe!