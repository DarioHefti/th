#!/bin/bash
set -e

REPO="DarioHefti/th"
BINARY_NAME="th"
INSTALL_DIR="${INSTALL_DIR:-$HOME/.local/bin}"
FORCE=false

usage() {
    cat <<EOF
Usage: install.sh [OPTIONS]

Install th (Terminal Help) CLI

OPTIONS:
    -d, --dir DIR       Installation directory (default: ~/.local/bin)
    -f, --force         Force reinstall
    -h, --help          Show this help message

EXAMPLES:
    install.sh
    install.sh -d /usr/local/bin
    install.sh -f
EOF
}

parse_args() {
    while [[ $# -gt 0 ]]; do
        case $1 in
            -d|--dir)
                INSTALL_DIR="$2"
                shift 2
                ;;
            -f|--force)
                FORCE=true
                shift
                ;;
            -h|--help)
                usage
                exit 0
                ;;
            *)
                echo "Unknown option: $1"
                usage
                exit 1
                ;;
        esac
    done
}

get_os() {
    case "$(uname -s)" in
        Linux*)     echo "linux" ;;
        Darwin*)    echo "darwin" ;;
        *)          echo "unknown" ;;
    esac
}

get_arch() {
    case "$(uname -m)" in
        x86_64)     echo "amd64" ;;
        aarch64)    echo "arm64" ;;
        arm64)      echo "arm64" ;;
        *)          echo "amd64" ;;
    esac
}

get_latest_version() {
    curl -sSL "https://api.github.com/repos/$REPO/releases/latest" | grep -o '"tag_name": "v[^"]*"' | cut -d'"' -f4 | sed 's/v//'
}

download_binary() {
    local os=$1
    local arch=$2
    local version=$3
    local extension=""

    if [ "$os" = "windows" ]; then
        extension=".exe"
    fi

    local filename="th-${os}-${arch}${extension}"
    local url="https://github.com/$REPO/releases/download/v${version}/${filename}"

    echo "Downloading $url..."
    curl -sSL "$url" -o "$INSTALL_DIR/th${extension}"

    if [ "$os" != "windows" ]; then
        chmod +x "$INSTALL_DIR/th${extension}"
    fi
}

main() {
    parse_args "$@"

    if [ -f "$INSTALL_DIR/$BINARY_NAME" ] && [ "$FORCE" = false ]; then
        echo "$BINARY_NAME is already installed at $INSTALL_DIR"
        echo "Use -f to force reinstall"
        exit 0
    fi

    mkdir -p "$INSTALL_DIR"

    local os=$(get_os)
    local arch=$(get_arch)
    local version=$(get_latest_version)

    if [ -z "$version" ]; then
        echo "Failed to get latest version"
        exit 1
    fi

    echo "Installing th v$version for $os/$arch..."

    download_binary "$os" "$arch" "$version"

    echo ""
    echo "✓ Installed to $INSTALL_DIR/th"
    echo ""
    echo "Add to PATH if not already added:"
    echo "  echo 'export PATH=\"\$PATH:$INSTALL_DIR\"' >> ~/.bashrc"
    echo "  source ~/.bashrc"
}

main "$@"
