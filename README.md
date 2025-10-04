# lua-bundler APT Repository

This is the APT repository for lua-bundler, hosted on GitHub Pages.

## Installation

```bash
echo "deb [trusted=yes] https://alfin-efendy.github.io/lua-bundler/ stable main" | sudo tee /etc/apt/sources.list.d/lua-bundler.list
sudo apt update
sudo apt install lua-bundler
```

## Quick Install

```bash
curl -fsSL https://alfin-efendy.github.io/lua-bundler/install.sh | sudo bash
```

## Repository Structure

- `dists/` - APT repository metadata
- `pool/` - Package files (.deb)
- `index.html` - Repository homepage
- `install.sh` - Quick installation script

This repository is automatically maintained by GitHub Actions.
