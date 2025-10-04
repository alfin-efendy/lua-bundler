#!/bin/bash
set -e

echo "🏗️ Installing lua-bundler from APT repository..."

# Check if running as root
if [[ $EUID -eq 0 ]]; then
    SUDO=""
else
    SUDO="sudo"
fi

# Add repository
echo "📦 Adding lua-bundler APT repository..."
echo "deb [trusted=yes] https://alfin-efendy.github.io/lua-bundler/ stable main" | $SUDO tee /etc/apt/sources.list.d/lua-bundler.list

# Update package list
echo "🔄 Updating package list..."
$SUDO apt update

# Install package
echo "⬇️ Installing lua-bundler..."
$SUDO apt install -y lua-bundler

echo "✅ lua-bundler installed successfully!"
echo ""
echo "Usage: lua-bundler -entry main.lua -output bundle.lua"
echo "Help:  lua-bundler -help"
