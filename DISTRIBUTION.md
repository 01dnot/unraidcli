# Distribution Guide

This document explains how to set up and maintain package manager distributions for unraidcli.

## Homebrew Tap (macOS/Linux)

### Setup

1. **Create a separate tap repository:**
   ```bash
   # On GitHub, create a new repository named: homebrew-unraidcli
   # Repository URL: https://github.com/01dnot/homebrew-unraidcli
   ```

2. **Initialize the tap repository:**
   ```bash
   mkdir homebrew-unraidcli
   cd homebrew-unraidcli
   git init
   ```

3. **Copy the formula:**
   ```bash
   cp /path/to/unraidcli/homebrew/unraidcli.rb Formula/unraidcli.rb
   ```

4. **Update SHA256 checksums after each release:**
   ```bash
   # Download each binary and calculate SHA256
   shasum -a 256 unraidcli-darwin-amd64
   shasum -a 256 unraidcli-darwin-arm64
   shasum -a 256 unraidcli-linux-amd64
   shasum -a 256 unraidcli-linux-arm64

   # Update the formula with the checksums
   ```

5. **Push to GitHub:**
   ```bash
   git add Formula/unraidcli.rb
   git commit -m "Add unraidcli formula"
   git push origin main
   ```

### Usage

Users can then install with:
```bash
brew tap 01dnot/unraidcli
brew install unraidcli
```

### Automation

Consider using [GoReleaser](https://goreleaser.com/) to automate:
- Building binaries for all platforms
- Creating GitHub releases
- Updating Homebrew formula
- Calculating checksums

Example `.goreleaser.yml`:
```yaml
project_name: unraidcli

builds:
  - env:
      - CGO_ENABLED=0
    goos:
      - linux
      - darwin
      - windows
    goarch:
      - amd64
      - arm64
    ignore:
      - goos: windows
        goarch: arm64

archives:
  - format: binary
    name_template: "{{ .ProjectName }}-{{ .Os }}-{{ .Arch }}"

brews:
  - repository:
      owner: 01dnot
      name: homebrew-unraidcli
    folder: Formula
    homepage: https://github.com/01dnot/unraidcli
    description: Command-line interface for managing Unraid servers
    install: |
      bin.install "unraidcli"
    test: |
      system "#{bin}/unraidcli --version"

release:
  github:
    owner: 01dnot
    name: unraidcli
```

## Other Package Managers

### Scoop (Windows)

1. Create a bucket repository: `scoop-bucket`
2. Add manifest for unraidcli
3. Users install with: `scoop bucket add 01dnot https://github.com/01dnot/scoop-bucket && scoop install unraidcli`

### AUR (Arch Linux)

1. Create `PKGBUILD` file
2. Submit to AUR
3. Community maintains the package
4. Users install with: `yay -S unraidcli` or `paru -S unraidcli`

### apt/deb (Debian/Ubuntu)

Consider using a PPA or hosting your own apt repository.

### Snap (Cross-platform Linux)

Package as a snap for easy installation across Linux distributions.

## Install Script

The `install.sh` script provides a universal installer that:
- Detects OS and architecture
- Downloads the appropriate binary
- Installs to `/usr/local/bin`
- Works on macOS, Linux (x64, ARM64)

Users can install with:
```bash
curl -sSL https://raw.githubusercontent.com/01dnot/unraidcli/main/install.sh | bash
```

## Distribution Checklist

- [x] GitHub Releases with binaries
- [x] Install script
- [ ] Homebrew tap (requires separate repo)
- [ ] GoReleaser automation
- [ ] Scoop bucket (Windows)
- [ ] AUR package (Arch Linux)
- [ ] Docker image
- [ ] Snap package
- [ ] Flatpak

## Versioning

Follow [Semantic Versioning](https://semver.org/):
- **MAJOR**: Incompatible API changes
- **MINOR**: Backwards-compatible functionality
- **PATCH**: Backwards-compatible bug fixes

Tag releases as `vMAJOR.MINOR.PATCH` (e.g., `v0.1.0`, `v1.0.0`)
