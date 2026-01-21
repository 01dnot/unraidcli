#!/bin/bash
set -e

# unraidcli installer script
# Usage: curl -sSL https://raw.githubusercontent.com/01dnot/unraidcli/main/install.sh | bash

VERSION="${UNRAIDCLI_VERSION:-latest}"
INSTALL_DIR="${UNRAIDCLI_INSTALL_DIR:-/usr/local/bin}"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

print_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Detect OS and architecture
detect_platform() {
    OS=$(uname -s | tr '[:upper:]' '[:lower:]')
    ARCH=$(uname -m)

    case "$OS" in
        linux)
            OS="linux"
            ;;
        darwin)
            OS="darwin"
            ;;
        *)
            print_error "Unsupported operating system: $OS"
            exit 1
            ;;
    esac

    case "$ARCH" in
        x86_64|amd64)
            ARCH="amd64"
            ;;
        aarch64|arm64)
            ARCH="arm64"
            ;;
        *)
            print_error "Unsupported architecture: $ARCH"
            exit 1
            ;;
    esac

    PLATFORM="${OS}-${ARCH}"
    print_info "Detected platform: $PLATFORM"
}

# Get latest version from GitHub
get_latest_version() {
    if [ "$VERSION" = "latest" ]; then
        print_info "Fetching latest version..."
        VERSION=$(curl -sSL https://api.github.com/repos/01dnot/unraidcli/releases/latest | grep '"tag_name"' | sed -E 's/.*"([^"]+)".*/\1/')
        if [ -z "$VERSION" ]; then
            print_error "Failed to fetch latest version"
            exit 1
        fi
    fi
    print_info "Installing version: $VERSION"
}

# Download and install
install_binary() {
    BINARY_NAME="unraidcli-${PLATFORM}"
    DOWNLOAD_URL="https://github.com/01dnot/unraidcli/releases/download/${VERSION}/${BINARY_NAME}"

    print_info "Downloading from: $DOWNLOAD_URL"

    # Create temporary directory
    TMP_DIR=$(mktemp -d)
    trap "rm -rf $TMP_DIR" EXIT

    # Download binary
    if ! curl -sSL -o "${TMP_DIR}/unraidcli" "$DOWNLOAD_URL"; then
        print_error "Failed to download binary"
        exit 1
    fi

    # Make executable
    chmod +x "${TMP_DIR}/unraidcli"

    # Install to destination
    print_info "Installing to ${INSTALL_DIR}/unraidcli"

    if [ -w "$INSTALL_DIR" ]; then
        mv "${TMP_DIR}/unraidcli" "${INSTALL_DIR}/unraidcli"
    else
        print_warn "Need sudo permissions to install to $INSTALL_DIR"
        sudo mv "${TMP_DIR}/unraidcli" "${INSTALL_DIR}/unraidcli"
    fi

    print_info "✓ Installation complete!"
}

# Verify installation
verify_installation() {
    if command -v unraidcli >/dev/null 2>&1; then
        VERSION_OUTPUT=$(unraidcli --version 2>&1 || echo "unknown")
        print_info "✓ unraidcli installed successfully"
        echo ""
        echo "Get started by configuring your server:"
        echo "  unraidcli config set --url http://YOUR_SERVER --apikey YOUR_API_KEY"
        echo ""
        echo "For help:"
        echo "  unraidcli --help"
    else
        print_error "Installation verification failed"
        echo "Make sure $INSTALL_DIR is in your PATH"
        exit 1
    fi
}

# Main installation flow
main() {
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo "  unraidcli installer"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo ""

    detect_platform
    get_latest_version
    install_binary
    verify_installation
}

main
