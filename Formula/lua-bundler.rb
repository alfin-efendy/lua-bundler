# Homebrew Formula for lua-bundler
class LuaBundler < Formula
  desc "Lua script bundler for Roblox development"
  homepage "https://github.com/alfin-efendy/lua-bundler"
  url "https://github.com/alfin-efendy/lua-bundler/releases/download/v1.0.0/lua-bundler-darwin-amd64"
  sha256 "SHA256_PLACEHOLDER"
  license "MIT"
  version "1.0.0"

  depends_on "git"

  def install
    bin.install "lua-bundler-darwin-amd64" => "lua-bundler"
  end

  test do
    # Test basic functionality
    system "#{bin}/lua-bundler", "-help"
  end
end