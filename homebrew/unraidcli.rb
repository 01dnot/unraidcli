class Unraidcli < Formula
  desc "Command-line interface for managing Unraid servers"
  homepage "https://github.com/01dnot/unraidcli"
  version "0.1.0"

  on_macos do
    if Hardware::CPU.intel?
      url "https://github.com/01dnot/unraidcli/releases/download/v0.1.0/unraidcli-darwin-amd64"
      sha256 "" # Will be filled automatically by release process
    else
      url "https://github.com/01dnot/unraidcli/releases/download/v0.1.0/unraidcli-darwin-arm64"
      sha256 "" # Will be filled automatically by release process
    end
  end

  on_linux do
    if Hardware::CPU.intel?
      url "https://github.com/01dnot/unraidcli/releases/download/v0.1.0/unraidcli-linux-amd64"
      sha256 "" # Will be filled automatically by release process
    else
      url "https://github.com/01dnot/unraidcli/releases/download/v0.1.0/unraidcli-linux-arm64"
      sha256 "" # Will be filled automatically by release process
    end
  end

  def install
    bin.install Dir["unraidcli*"].first => "unraidcli"
  end

  test do
    system "#{bin}/unraidcli", "--version"
  end
end
