#!/bin/bash
set -e

REPO="DarioHefti/th"
BINARY_NAME="th"
INSTALL_DIR="${INSTALL_DIR:-$HOME/.local/bin}"
FORCE=false
MAX_RETRIES=3
RETRY_DELAY=2

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

log_info() {
    echo "[INFO] $1"
}

log_error() {
    echo "[ERROR] $1" >&2
}

parse_args() {
    while [[ $# -gt 0 ]]; do
        case $1 in
            -d|--dir)
                if [[ -z "$2" || "$2" == -* ]]; then
                    log_error "Option $1 requires a directory argument"
                    usage
                    exit 1
                fi
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
                log_error "Unknown option: $1"
                usage
                exit 1
                ;;
        esac
    done
}

check_dependencies() {
    local missing=()
    
    if ! command -v curl &> /dev/null; then
        missing+=("curl")
    fi
    
    if [[ ${#missing[@]} -gt 0 ]]; then
        log_error "Missing required dependency: ${missing[*]}"
        log_error "Please install the missing dependency and try again"
        exit 1
    fi
}

get_os() {
    case "$(uname -s)" in
        Linux*)     echo "linux" ;;
        Darwin*)    echo "darwin" ;;
        *)          log_error "Unsupported operating system: $(uname -s)"; exit 1 ;;
    esac
}

get_arch() {
    case "$(uname -m)" in
        x86_64)     echo "amd64" ;;
        aarch64|arm64)  echo "arm64" ;;
        *)          log_error "Unsupported architecture: $(uname -m)"; exit 1 ;;
    esac
}

get_latest_version() {
    local version
    version=$(curl -sSL --fail --connect-timeout 10 "https://api.github.com/repos/$REPO/releases/latest" | grep -o '"tag_name": "v[^"]*"' | cut -d'"' -f4 | sed 's/v//')
    
    if [[ -z "$version" ]]; then
        return 1
    fi
    
    if [[ ! "$version" =~ ^[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
        log_error "Invalid version format: $version"
        return 1
    fi
    
    echo "$version"
}

check_url() {
    local url="$1"
    local http_code
    http_code=$(curl -sSL -o /dev/null -w "%{http_code}" --head --fail --connect-timeout 10 "$url" 2>/dev/null || echo "000")
    [[ "$http_code" == "200" ]]
}

download_binary() {
    local os=$1
    local arch=$2
    local version=$3
    local extension=""

    if [[ "$os" == "windows" ]]; then
        extension=".exe"
    fi

    local filename="th-${os}-${arch}${extension}"
    local url="https://github.com/$REPO/releases/download/v${version}/${filename}"
    local output_path="$INSTALL_DIR/th${extension}"
    
    log_info "Downloading $filename..."
    
    if ! check_url "$url"; then
        log_error "Release not found for $os/$arch version $version"
        log_error "URL: $url"
        exit 1
    fi
    
    local attempt=1
    local success=false
    
    while [[ $attempt -le $MAX_RETRIES ]]; do
        log_info "Download attempt $attempt of $MAX_RETRIES..."
        
        if curl -sSL --fail --connect-timeout 30 --retry 2 "$url" -o "$output_path" 2>/dev/null; then
            if [[ -s "$output_path" ]]; then
                success=true
                break
            else
                log_error "Downloaded file is empty"
            fi
        fi
        
        log_info "Retrying in ${RETRY_DELAY}s..."
        sleep "$RETRY_DELAY"
        ((attempt++))
    done
    
    if [[ "$success" != "true" ]]; then
        log_error "Failed to download after $MAX_RETRIES attempts"
        rm -f "$output_path"
        exit 1
    fi
    
    if [[ "$os" != "windows" ]]; then
        chmod +x "$output_path"
    fi
    
    log_info "Download complete"
}

verify_binary() {
    local binary_path="$1"
    
    if [[ ! -f "$binary_path" ]]; then
        log_error "Binary not found at $binary_path"
        return 1
    fi
    
    if [[ ! -s "$binary_path" ]]; then
        log_error "Binary is empty"
        return 1
    fi
    
    if [[ "$binary_path" == *.exe ]]; then
        return 0
    fi
    
    if [[ ! -x "$binary_path" ]]; then
        log_error "Binary is not executable"
        return 1
    fi
    
    return 0
}

main() {
    parse_args "$@"
    check_dependencies
    
    if [[ -f "$INSTALL_DIR/$BINARY_NAME" ]] && [[ "$FORCE" != "true" ]]; then
        log_info "$BINARY_NAME is already installed at $INSTALL_DIR"
        log_info "Use -f to force reinstall"
        exit 0
    fi
    
    if [[ ! -d "$INSTALL_DIR" ]]; then
        mkdir -p "$INSTALL_DIR" || {
            log_error "Failed to create installation directory: $INSTALL_DIR"
            exit 1
        }
    fi
    
    if [[ ! -w "$INSTALL_DIR" ]]; then
        log_error "Installation directory is not writable: $INSTALL_DIR"
        exit 1
    fi
    
    local os
    local arch
    local version
    
    os=$(get_os)
    arch=$(get_arch)
    
    log_info "Fetching latest version..."
    version=$(get_latest_version) || {
        log_error "Failed to get latest version"
        log_error "Please check your network connection and try again"
        exit 1
    }
    
    log_info "Installing th v$version for $os/$arch..."
    
    download_binary "$os" "$arch" "$version"
    
    local binary_path="$INSTALL_DIR/th"
    [[ "$os" == "windows" ]] && binary_path+=".exe"
    
    if ! verify_binary "$binary_path"; then
        log_error "Binary verification failed"
        rm -f "$binary_path"
        exit 1
    fi
    
    echo ""
    echo "✓ Installed to $binary_path"
    echo ""
    echo "Add to PATH if not already added:"
    echo "  echo 'export PATH=\"\$PATH:$INSTALL_DIR\"' >> ~/.bashrc"
    echo "  source ~/.bashrc"
}

main "$@"
