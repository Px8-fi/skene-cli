#!/bin/bash
set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

REPO="Px8-fi/skene-cli"
BINARY_NAME="skene"
INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"
VERSION="${VERSION:-v030}"

detect_platform() {
    local os arch
    case "$(uname -s)" in
        Linux*)   os="linux" ;;
        Darwin*)  os="darwin" ;;
        CYGWIN*|MINGW*|MSYS*) os="windows" ;;
        *) echo -e "${RED}Error: Unsupported OS: $(uname -s)${NC}" >&2; exit 1 ;;
    esac
    case "$(uname -m)" in
        x86_64|amd64)  arch="amd64" ;;
        arm64|aarch64) arch="arm64" ;;
        *) echo -e "${RED}Error: Unsupported architecture: $(uname -m)${NC}" >&2; exit 1 ;;
    esac
    echo "${os}-${arch}"
}

download_binary() {
    local platform="$1"
    local tmpdir
    tmpdir=$(mktemp -d)
    local asset="${BINARY_NAME}-${platform}"
    local url="https://github.com/${REPO}/releases/download/${VERSION}/${asset}.tar.gz"

    echo -e "${BLUE}Downloading ${BINARY_NAME} ${VERSION} for ${platform}...${NC}"
    if ! curl -fSL -o "${tmpdir}/${asset}.tar.gz" "$url"; then
        rm -rf "$tmpdir"
        echo -e "${RED}Error: Download failed. Check https://github.com/${REPO}/releases${NC}" >&2
        exit 1
    fi

    tar -xzf "${tmpdir}/${asset}.tar.gz" -C "$tmpdir"
    mv "${tmpdir}/${asset}" "${tmpdir}/${BINARY_NAME}"
    chmod +x "${tmpdir}/${BINARY_NAME}"

    # Strip macOS quarantine flag
    xattr -d com.apple.quarantine "${tmpdir}/${BINARY_NAME}" 2>/dev/null || true

    echo "${tmpdir}/${BINARY_NAME}"
}

install_binary() {
    local binary="$1"

    if [ ! -d "$INSTALL_DIR" ]; then
        echo -e "${YELLOW}Creating ${INSTALL_DIR}...${NC}"
        if [ -w "$(dirname "$INSTALL_DIR")" ]; then
            mkdir -p "$INSTALL_DIR"
        else
            sudo mkdir -p "$INSTALL_DIR"
        fi
    fi

    echo -e "${YELLOW}Installing to ${INSTALL_DIR}/${BINARY_NAME}...${NC}"
    if [ -w "$INSTALL_DIR" ]; then
        cp "$binary" "${INSTALL_DIR}/${BINARY_NAME}"
    else
        sudo cp "$binary" "${INSTALL_DIR}/${BINARY_NAME}"
    fi

    rm -f "$binary"
    rmdir "$(dirname "$binary")" 2>/dev/null || true
}

main() {
    echo -e "${GREEN}Installing ${BINARY_NAME}...${NC}"

    local platform
    platform=$(detect_platform)
    echo -e "Detected platform: ${platform}"

    local binary
    binary=$(download_binary "$platform")
    install_binary "$binary"

    if command -v "$BINARY_NAME" >/dev/null 2>&1; then
        echo -e "${GREEN}✓ Installed ${BINARY_NAME} to ${INSTALL_DIR}/${BINARY_NAME}${NC}"
        echo -e "  Run '${BINARY_NAME}' to get started"
    else
        echo -e "${GREEN}✓ Installed to ${INSTALL_DIR}/${BINARY_NAME}${NC}"
        echo -e "${YELLOW}⚠ ${INSTALL_DIR} may not be in your PATH${NC}"
    fi
}

main "$@"
