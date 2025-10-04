#!/bin/bash

# Script to create and publish a new release

set -e

echo "=== Release Creator for lua-bundler ==="
echo ""

# Check if we're in correct repository
REPO_NAME=$(basename "$(git rev-parse --show-toplevel)" 2>/dev/null || echo "unknown")
if [ "$REPO_NAME" != "lua-bundler" ]; then
    echo "❌ Not in lua-bundler repository"
    exit 1
fi

echo "✅ In lua-bundler repository"

# Check if GitHub CLI is available
if ! command -v gh &> /dev/null; then
    echo "❌ GitHub CLI not found. Install from: https://cli.github.com/"
    exit 1
fi

echo "✅ GitHub CLI found"

# Check authentication
if ! gh auth status &> /dev/null; then
    echo "❌ Not authenticated with GitHub CLI"
    echo "Run: gh auth login"
    exit 1
fi

echo "✅ GitHub CLI authenticated"

# Check for uncommitted changes
if ! git diff-index --quiet HEAD --; then
    echo "⚠️ You have uncommitted changes. Please commit or stash them first."
    git status --short
    exit 1
fi

echo "✅ Working directory clean"

# Check CHANGELOG.md
echo ""
echo "📝 CHANGELOG.md Check:"
if [ ! -f "CHANGELOG.md" ]; then
    echo "⚠️ CHANGELOG.md not found. Please create one to document changes."
else
    echo "✅ CHANGELOG.md exists"
    
    # Check if there are unreleased changes
    if grep -q "## \[Unreleased\]" CHANGELOG.md; then
        echo "📋 Found [Unreleased] section in CHANGELOG.md"
        
        # Show unreleased changes
        echo ""
        echo "🔍 Current unreleased changes:"
        echo "---"
        sed -n '/## \[Unreleased\]/,/## \[/p' CHANGELOG.md | head -n -1
        echo "---"
        
        echo ""
        read -p "❓ Have you updated CHANGELOG.md with changes for this release? (y/n): " changelog_updated
        if [[ "$changelog_updated" != "y" && "$changelog_updated" != "Y" ]]; then
            echo ""
            echo "📝 Please update CHANGELOG.md before creating a release:"
            echo "1. Move changes from [Unreleased] to new version section"
            echo "2. Add proper version number and date"
            echo "3. Document all Added/Changed/Fixed items"
            echo ""
            echo "Example format:"
            echo "## [1.0.1] - $(date +%Y-%m-%d)"
            echo ""
            exit 1
        fi
    else
        echo "⚠️ No [Unreleased] section found in CHANGELOG.md"
        echo "Consider adding an [Unreleased] section for future changes"
    fi
fi

# Get current version/tag
CURRENT_TAG=$(git describe --tags --abbrev=0 2>/dev/null || echo "")
echo ""
echo "📋 Current latest tag: ${CURRENT_TAG:-"none"}"

# Suggest next version
if [ -n "$CURRENT_TAG" ]; then
    # Extract version numbers
    VERSION_NUM=${CURRENT_TAG#v}
    IFS='.' read -ra VERSION_PARTS <<< "$VERSION_NUM"
    MAJOR=${VERSION_PARTS[0]:-0}
    MINOR=${VERSION_PARTS[1]:-0}
    PATCH=${VERSION_PARTS[2]:-0}
    
    # Suggest versions
    PATCH_VERSION="v$MAJOR.$MINOR.$((PATCH + 1))"
    MINOR_VERSION="v$MAJOR.$((MINOR + 1)).0"
    MAJOR_VERSION="v$((MAJOR + 1)).0.0"
else
    # First release
    PATCH_VERSION="v1.0.0"
    MINOR_VERSION="v1.0.0"
    MAJOR_VERSION="v1.0.0"
fi

echo ""
echo "🏷️ Suggested versions:"
echo "1. Patch release: $PATCH_VERSION (bug fixes)"
echo "2. Minor release: $MINOR_VERSION (new features)"
echo "3. Major release: $MAJOR_VERSION (breaking changes)"
echo "4. Custom version"
echo ""

read -p "Select version type (1-4): " version_choice

case $version_choice in
    1)
        NEW_VERSION="$PATCH_VERSION"
        ;;
    2)
        NEW_VERSION="$MINOR_VERSION"
        ;;
    3)
        NEW_VERSION="$MAJOR_VERSION"
        ;;
    4)
        read -p "Enter custom version (e.g., v1.2.3): " NEW_VERSION
        if [[ ! "$NEW_VERSION" =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
            echo "❌ Invalid version format. Use v1.2.3 format."
            exit 1
        fi
        ;;
    *)
        echo "❌ Invalid selection"
        exit 1
        ;;
esac

echo ""
echo "🎯 Selected version: $NEW_VERSION"

# Check if tag already exists
if git tag -l | grep -q "^$NEW_VERSION$"; then
    echo "❌ Tag $NEW_VERSION already exists"
    exit 1
fi

# Get release notes
echo ""
echo "📝 Release notes (press Ctrl+D when finished, Ctrl+C to cancel):"
echo "---"
RELEASE_NOTES=$(cat)
echo "---"

if [ -z "$RELEASE_NOTES" ]; then
    RELEASE_NOTES="Release $NEW_VERSION

## What's Changed

- Bug fixes and improvements
- See commit history for detailed changes

## Installation

### Direct Download
Download the appropriate binary for your platform from the assets below.

### Package Managers
- **Homebrew**: \`brew install alfin-efendy/tap/lua-bundler\`
- **APT**: \`curl -fsSL https://alfin-efendy.github.io/lua-bundler/install.sh | sudo bash\`

## Checksums

Verify your download with the provided \`.sha256\` files."
fi

echo ""
echo "🚀 Creating release $NEW_VERSION..."

# Create and push tag
echo "Creating tag..."
git tag -a "$NEW_VERSION" -m "Release $NEW_VERSION"
git push origin "$NEW_VERSION"

echo "✅ Tag pushed to GitHub"

# Wait a moment for Release pipeline to start
echo "⏳ Waiting for Release pipeline to start..."
sleep 10

# Monitor Release pipeline
echo ""
echo "📊 Release Pipeline Status:"
gh run list --workflow=release.yml --limit=3

echo ""
echo "🎯 Next Steps:"
echo "1. Monitor Release pipeline: gh run watch"
echo "2. Check build status: gh run list --workflow=release.yml"
echo "3. Release will be created automatically with binaries"
echo "4. Package managers will auto-update:"
echo "   • APT repository (immediate)"
echo "   • Homebrew tap (immediate)" 
echo "   • Winget (requires manual approval)"
echo ""

# Ask if user wants to create release now
echo "📋 Options:"
echo "1. Wait for CI and create release manually"
echo "2. Create release now (CI artifacts will be added later)"
echo "3. Exit and handle manually"
echo ""

read -p "Select option (1-3): " release_option

case $release_option in
    1)
        echo "✅ Tag created successfully. Monitor CI and create release when ready."
        echo "Command to create release later:"
        echo "gh release create $NEW_VERSION --title \"Release $NEW_VERSION\" --notes-file <(echo \"$RELEASE_NOTES\")"
        ;;
    2)
        echo "Creating GitHub release..."
        echo "$RELEASE_NOTES" | gh release create "$NEW_VERSION" \
            --title "Release $NEW_VERSION" \
            --notes-file -
        
        echo "✅ Release created! CI will add artifacts when build completes."
        ;;
    3)
        echo "✅ Tag created. Handle release creation manually."
        ;;
    *)
        echo "✅ Tag created successfully."
        ;;
esac

echo ""
echo "🎉 Release process initiated!"
echo ""
echo "� Next Steps:"
echo "1. Monitor CI and Release pipelines"
echo "2. Update CHANGELOG.md if not already done"
echo "3. Commit changelog updates if needed"
echo ""
echo "�📊 Monitor progress:"
echo "- CI Pipeline: https://github.com/alfin-efendy/lua-bundler/actions"
echo "- Releases: https://github.com/alfin-efendy/lua-bundler/releases"
echo "- APT Repo: https://alfin-efendy.github.io/lua-bundler/"
echo ""
echo "💡 Tip: Keep CHANGELOG.md updated for better release documentation!"