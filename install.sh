#!/bin/bash
set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
REPO="Px8-fi/skene-cli"
BINARY_NAME="skene"
INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"
VERSION="${VERSION:-latest}"
USE_LOCAL="${USE_LOCAL:-false}"

# Detect OS and Architecture
detect_platform() {
    local os=""
    local arch=""
    
    case "$(uname -s)" in
        Linux*)     os="linux" ;;
        Darwin*)    os="darwin" ;;
        CYGWIN*)    os="windows" ;;
        MINGW*)     os="windows" ;;
        MSYS*)      os="windows" ;;
        *)          echo -e "${RED}Error: Unsupported OS: $(uname -s)${NC}" >&2; exit 1 ;;
    esac
    
    case "$(uname -m)" in
        x86_64|amd64)   arch="amd64" ;;
        arm64|aarch64)  arch="arm64" ;;
        *)              echo -e "${RED}Error: Unsupported architecture: $(uname -m)${NC}" >&2; exit 1 ;;
    esac
    
    echo "${os}-${arch}"
}

# Get latest release version from GitHub
get_latest_version() {
    if [ "$VERSION" = "latest" ]; then
        local version=$(curl -s "https://api.github.com/repos/${REPO}/releases/latest" 2>/dev/null | \
        grep '"tag_name":' | \
        sed -E 's/.*"([^"]+)".*/\1/')
        
        if [ -z "$version" ]; then
            echo ""
        else
            echo "$version"
        fi
    else
        echo "$VERSION"
    fi
}

# Check if local build exists
check_local_build() {
    local build_dir="build"
    if [ -f "$build_dir/$BINARY_NAME" ]; then
        echo "$build_dir/$BINARY_NAME"
        return 0
    fi
    return 1
}

# Build locally
build_local() {
    echo -e "${YELLOW}Building ${BINARY_NAME} locally...${NC}"
    
    if ! command -v go >/dev/null 2>&1; then
        echo -e "${RED}Error: Go is required to build locally. Install Go or use a pre-built release.${NC}" >&2
        return 1
    fi
    
    # Check if we're in a Go project directory
    if [ ! -f "go.mod" ]; then
        echo -e "${RED}Error: go.mod not found. Please run this script from the project root directory.${NC}" >&2
        return 1
    fi
    
    if [ -f "Makefile" ]; then
        # Try 'make build' first, fall back to 'make' if build target doesn't exist
        if make build 2>/dev/null; then
            :
        elif make 2>/dev/null; then
            :
        else
            echo -e "${YELLOW}Makefile found but build target failed. Trying go build directly...${NC}"
            mkdir -p build
            if ! go build -o "build/$BINARY_NAME" ./cmd/skene 2>&1; then
                echo -e "${RED}Error: Build failed${NC}" >&2
                return 1
            fi
        fi
    else
        mkdir -p build
        if ! go build -o "build/$BINARY_NAME" ./cmd/skene 2>&1; then
            echo -e "${RED}Error: Build failed. Make sure you're in the project root directory.${NC}" >&2
            return 1
        fi
    fi
    
    if [ -f "build/$BINARY_NAME" ]; then
        echo "build/$BINARY_NAME"
        return 0
    else
        echo -e "${RED}Error: Build completed but binary not found at build/$BINARY_NAME${NC}" >&2
        return 1
    fi
}

# Download binary
download_binary() {
    local platform=$1
    local version=$2
    local url=""
    local filename=""
    
    if [ "$platform" = "linux-amd64" ]; then
        filename="${BINARY_NAME}-linux-amd64"
    elif [ "$platform" = "darwin-amd64" ]; then
        filename="${BINARY_NAME}-darwin-amd64"
    elif [ "$platform" = "darwin-arm64" ]; then
        filename="${BINARY_NAME}-darwin-arm64"
    elif [ "$platform" = "windows-amd64" ]; then
        filename="${BINARY_NAME}-windows-amd64.exe"
    else
        echo -e "${RED}Error: Unsupported platform: ${platform}${NC}" >&2
        exit 1
    fi
    
    url="https://github.com/${REPO}/releases/download/${version}/${filename}"
    
    echo -e "${YELLOW}Downloading ${BINARY_NAME} ${version} for ${platform}...${NC}"
    
    if command -v curl >/dev/null 2>&1; then
        if ! curl -L -f -o "$filename" "$url" 2>/dev/null; then
            return 1
        fi
    elif command -v wget >/dev/null 2>&1; then
        if ! wget -O "$filename" "$url" 2>/dev/null; then
            return 1
        fi
    else
        echo -e "${RED}Error: curl or wget is required${NC}" >&2
        exit 1
    fi
    
    echo "$filename"
}

# Install binary
install_binary() {
    local binary_file=$1
    local target_path="${INSTALL_DIR}/${BINARY_NAME}"
    
    if [ -z "$binary_file" ] || [ ! -f "$binary_file" ]; then
        echo -e "${RED}Error: Binary file not found: ${binary_file}${NC}" >&2
        exit 1
    fi
    
    # Create install directory if it doesn't exist
    if [ ! -d "$INSTALL_DIR" ]; then
        echo -e "${YELLOW}Creating directory ${INSTALL_DIR}...${NC}"
        if [ -w "$(dirname "$INSTALL_DIR")" ]; then
            mkdir -p "$INSTALL_DIR"
        else
            sudo mkdir -p "$INSTALL_DIR"
        fi
    fi
    
    # Make binary executable (for Unix-like systems)
    if [[ "$binary_file" != *.exe ]]; then
        chmod +x "$binary_file"
    fi
    
    # Install binary
    echo -e "${YELLOW}Installing to ${target_path}...${NC}"
    if [ -w "$INSTALL_DIR" ]; then
        cp "$binary_file" "$target_path"
    else
        sudo cp "$binary_file" "$target_path"
    fi
    
    # Clean up downloaded file (only if it's in current directory and not the build dir)
    if [ -f "$binary_file" ] && [ "$(dirname "$binary_file")" = "." ]; then
        rm -f "$binary_file"
    fi
    
    echo -e "${GREEN}✓ Successfully installed ${BINARY_NAME} to ${target_path}${NC}"
}

# Verify installation
verify_installation() {
    if command -v "$BINARY_NAME" >/dev/null 2>&1; then
        local version=$("$BINARY_NAME" --version 2>/dev/null || echo "unknown")
        echo -e "${GREEN}✓ Installation verified${NC}"
        echo -e "  Run '${BINARY_NAME}' to get started"
        return 0
    else
        echo -e "${YELLOW}⚠ Warning: ${BINARY_NAME} not found in PATH${NC}"
        echo -e "  You may need to add ${INSTALL_DIR} to your PATH"
        return 1
    fi
}

# Main installation flow
main() {
    echo -e "${GREEN}Installing ${BINARY_NAME}...${NC}"
    
    # Detect platform
    local platform=$(detect_platform)
    echo -e "Detected platform: ${platform}"
    
    local binary_file=""
    local temp_dir=""
    
    # If USE_LOCAL is set or we're in the repo directory, try local build first
    if [ "$USE_LOCAL" = "true" ] || [ -f "go.mod" ]; then
        if check_local_build >/dev/null 2>&1; then
            echo -e "${BLUE}Using existing local build...${NC}"
            binary_file=$(check_local_build)
        elif [ -f "go.mod" ]; then
            echo -e "${BLUE}No local build found. Building from source...${NC}"
            if ! binary_file=$(build_local); then
                if [ "$USE_LOCAL" = "true" ]; then
                    echo -e "${RED}Error: Failed to build locally${NC}" >&2
                    exit 1
                fi
                # Clear binary_file so we can try GitHub download
                binary_file=""
            fi
        elif [ "$USE_LOCAL" = "true" ]; then
            echo -e "${RED}Error: go.mod not found. Please run this script from the project root directory.${NC}" >&2
            echo -e "${YELLOW}Current directory: $(pwd)${NC}" >&2
            exit 1
        fi
    fi
    
    # If no local binary and USE_LOCAL is not set, try downloading from GitHub
    if [ -z "$binary_file" ] && [ "$USE_LOCAL" != "true" ]; then
        # Get version
        local version=$(get_latest_version)
        
        if [ -z "$version" ]; then
            echo -e "${YELLOW}⚠ No GitHub releases found.${NC}"
            if [ -f "go.mod" ]; then
                echo -e "${BLUE}Building from source instead...${NC}"
                binary_file=$(build_local)
            else
                echo -e "${YELLOW}No releases available. Attempting to clone and build from source...${NC}"
                
                # Try to clone and build
                temp_dir=$(mktemp -d)
                clone_dir="${temp_dir}/skene-cli"
                original_dir=$(pwd)
                
                echo -e "${BLUE}Cloning repository to ${clone_dir}...${NC}"
                if git clone --depth 1 --branch Rust-impelementation https://github.com/${REPO}.git "${clone_dir}" 2>/dev/null; then
                    cd "${clone_dir}"
                    if [ -f "go.mod" ]; then
                        echo -e "${BLUE}Building from source...${NC}"
                        if binary_file=$(build_local); then
                            echo -e "${GREEN}Build successful!${NC}"
                            # Make path absolute
                            if [ ! "$(echo "$binary_file" | cut -c1)" = "/" ]; then
                                binary_file="${clone_dir}/${binary_file}"
                            fi
                            # Return to original directory
                            cd "$original_dir"
                        else
                            echo -e "${RED}Build failed.${NC}" >&2
                            cd "$original_dir"
                            rm -rf "$temp_dir"
                            exit 1
                        fi
                    else
                        echo -e "${RED}Error: Repository cloned but go.mod not found${NC}" >&2
                        cd "$original_dir"
                        rm -rf "$temp_dir"
                        exit 1
                    fi
                else
                    echo -e "${RED}Error: Failed to clone repository${NC}" >&2
                    echo -e "${YELLOW}Please either:${NC}"
                    echo -e "  1. Create a GitHub release with binaries"
                    echo -e "  2. Clone the repository manually: git clone https://github.com/${REPO}"
                    echo -e "  3. Run this script from the repository directory"
                    rm -rf "$temp_dir"
                    exit 1
                fi
            fi
        else
            echo -e "Installing version: ${version}"
            binary_file=$(download_binary "$platform" "$version")
            
            # If download failed, try building locally as fallback (if in repo)
            if [ $? -ne 0 ] && [ -f "go.mod" ]; then
                echo -e "${YELLOW}Download failed. Building from source instead...${NC}"
                binary_file=$(build_local)
            fi
        fi
    fi
    
    if [ -z "$binary_file" ] || [ ! -f "$binary_file" ]; then
        echo -e "${RED}Error: Could not obtain binary${NC}" >&2
        # Clean up temp directory if we cloned the repo
        if [ -n "$temp_dir" ] && [ -d "$temp_dir" ]; then
            rm -rf "$temp_dir"
        fi
        exit 1
    fi
    
    # Install binary
    install_binary "$binary_file"
    
    # Clean up temp directory if we cloned the repo
    if [ -n "$temp_dir" ] && [ -d "$temp_dir" ]; then
        echo -e "${BLUE}Cleaning up temporary files...${NC}"
        rm -rf "$temp_dir"
    fi
    
    # Verify installation
    verify_installation
    
    echo -e "${GREEN}Installation complete!${NC}"
}

# Run main function
main "$@"
