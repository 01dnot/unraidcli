# Homebrew Setup Guide

This guide explains how to set up automatic Homebrew formula updates using GoReleaser.

## What GoReleaser Does

When you push a new tag (e.g., `v0.2.0`), GoReleaser automatically:
- ✅ Builds binaries for all platforms (Linux, macOS, Windows)
- ✅ Creates GitHub release with all binaries
- ✅ Calculates SHA256 checksums
- ✅ Updates Homebrew formula in `homebrew-unraidcli` repo
- ✅ Users can install/update with `brew install unraidcli`

## Setup Steps

### 1. Create Homebrew Tap Repository

On GitHub, create a **new public repository**:
- Name: `homebrew-unraidcli`
- URL: `https://github.com/01dnot/homebrew-unraidcli`
- Description: "Homebrew tap for unraidcli"
- Public repository
- Initialize with README

### 2. Create GitHub Personal Access Token

1. Go to: https://github.com/settings/tokens
2. Click "Generate new token" → "Generate new token (classic)"
3. Give it a name: `HOMEBREW_TAP_TOKEN`
4. Select scopes:
   - ✅ `repo` (Full control of private repositories)
5. Click "Generate token"
6. **Copy the token** (you won't see it again!)

### 3. Add Secret to unraidcli Repository

1. Go to: https://github.com/01dnot/unraidcli/settings/secrets/actions
2. Click "New repository secret"
3. Name: `HOMEBREW_TAP_GITHUB_TOKEN`
4. Value: Paste the token you just created
5. Click "Add secret"

### 4. Done!

That's it! Now when you create a new release:

```bash
git tag v0.2.0
git push origin v0.2.0
```

GitHub Actions will automatically:
- Build all binaries
- Create GitHub release
- Update the Homebrew formula in `homebrew-unraidcli`

## User Installation

Once set up, users can install with:

```bash
# Add the tap
brew tap 01dnot/unraidcli

# Install
brew install unraidcli

# Update
brew upgrade unraidcli

# Uninstall
brew uninstall unraidcli
```

Or use the shorthand:
```bash
brew install 01dnot/unraidcli/unraidcli
```

## Verification

After your next release, check that:
1. The release appears at: https://github.com/01dnot/unraidcli/releases
2. The formula was updated at: https://github.com/01dnot/homebrew-unraidcli/blob/main/Formula/unraidcli.rb
3. You can install with: `brew tap 01dnot/unraidcli && brew install unraidcli`

## Troubleshooting

### "Failed to update Homebrew formula"

**Cause:** Missing or incorrect `HOMEBREW_TAP_GITHUB_TOKEN` secret

**Fix:**
1. Verify the secret exists in repository settings
2. Ensure the token has `repo` scope
3. Check the token hasn't expired

### Formula not updating

**Cause:** The `homebrew-unraidcli` repository doesn't exist

**Fix:**
1. Create the repository: https://github.com/new
2. Name it exactly: `homebrew-unraidcli`
3. Make it public
4. Re-run the release workflow

## Manual Formula Update (Fallback)

If automation fails, you can manually update the formula:

1. Clone the tap repo:
   ```bash
   git clone https://github.com/01dnot/homebrew-unraidcli.git
   cd homebrew-unraidcli
   mkdir -p Formula
   ```

2. Copy the formula:
   ```bash
   cp /path/to/unraidcli/homebrew/unraidcli.rb Formula/
   ```

3. Update version and SHA256 checksums in `Formula/unraidcli.rb`

4. Commit and push:
   ```bash
   git add Formula/unraidcli.rb
   git commit -m "Update to v0.2.0"
   git push
   ```

## Getting into Homebrew Core

Once your project is mature and popular, you can submit to official Homebrew:

1. Your tap must be stable for a while
2. Must have significant usage/stars
3. Submit PR to: https://github.com/Homebrew/homebrew-core

Benefits of Homebrew Core:
- Users can install with just: `brew install unraidcli`
- No need to add tap
- More visibility and trust

For now, stick with the tap approach - it's easier and gives you full control!
