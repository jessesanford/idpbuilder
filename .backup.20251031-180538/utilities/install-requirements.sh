#!/bin/bash

# Software Factory 2.0 - Requirements Installation Script
# Automatically installs required dependencies based on your system

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo "======================================"
echo "Software Factory 2.0 Requirements Installer"
echo "======================================"
echo ""

# Detect OS and package manager
detect_os() {
    if [[ "$OSTYPE" == "linux-gnu"* ]]; then
        if command -v apt-get &> /dev/null; then
            PKG_MANAGER="apt"
            OS="debian"
        elif command -v yum &> /dev/null; then
            PKG_MANAGER="yum"
            OS="redhat"
        elif command -v dnf &> /dev/null; then
            PKG_MANAGER="dnf"
            OS="fedora"
        elif command -v pacman &> /dev/null; then
            PKG_MANAGER="pacman"
            OS="arch"
        elif command -v apk &> /dev/null; then
            PKG_MANAGER="apk"
            OS="alpine"
        else
            echo -e "${RED}❌ Unsupported Linux distribution${NC}"
            exit 1
        fi
    elif [[ "$OSTYPE" == "darwin"* ]]; then
        PKG_MANAGER="brew"
        OS="macos"
    else
        echo -e "${RED}❌ Unsupported OS: $OSTYPE${NC}"
        exit 1
    fi
    
    echo -e "${GREEN}✅ Detected OS: $OS${NC}"
    echo -e "${GREEN}✅ Package manager: $PKG_MANAGER${NC}"
    echo ""
}

# Check if running with sufficient privileges
check_privileges() {
    if [[ "$PKG_MANAGER" == "apt" || "$PKG_MANAGER" == "yum" || "$PKG_MANAGER" == "dnf" || "$PKG_MANAGER" == "pacman" ]]; then
        if [[ $EUID -ne 0 ]]; then
            echo -e "${YELLOW}⚠️  This script needs sudo privileges to install packages${NC}"
            echo -e "${YELLOW}   Re-running with sudo...${NC}"
            exec sudo "$0" "$@"
        fi
    fi
}

# Function to check if a command exists
command_exists() {
    command -v "$1" &> /dev/null
}

# Function to get version of installed tool
get_version() {
    local cmd="$1"
    case "$cmd" in
        yq)
            yq --version 2>/dev/null | grep -oE '[0-9]+\.[0-9]+' | head -1 || echo "0.0"
            ;;
        git)
            git --version | grep -oE '[0-9]+\.[0-9]+' | head -1 || echo "0.0"
            ;;
        jq)
            jq --version | grep -oE '[0-9]+\.[0-9]+' | head -1 || echo "0.0"
            ;;
        bash)
            bash --version | head -1 | grep -oE '[0-9]+\.[0-9]+' | head -1 || echo "0.0"
            ;;
        gh)
            gh --version | head -1 | grep -oE '[0-9]+\.[0-9]+' | head -1 || echo "0.0"
            ;;
        rg)
            rg --version | head -1 | grep -oE '[0-9]+\.[0-9]+' | head -1 || echo "0.0"
            ;;
        *)
            echo "0.0"
            ;;
    esac
}

# Version comparison
version_ge() {
    [ "$(printf '%s\n' "$1" "$2" | sort -V | head -n1)" = "$2" ]
}

# Install functions for different package managers
install_with_apt() {
    local package="$1"
    local apt_package="$2"
    
    echo -e "${BLUE}Installing $package with apt...${NC}"
    apt-get update -qq
    apt-get install -y "$apt_package"
}

install_with_brew() {
    local package="$1"
    local brew_package="$2"
    
    echo -e "${BLUE}Installing $package with brew...${NC}"
    
    # Check if brew is installed
    if ! command_exists brew; then
        echo -e "${YELLOW}Installing Homebrew first...${NC}"
        /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
    fi
    
    brew install "$brew_package"
}

install_with_yum() {
    local package="$1"
    local yum_package="$2"
    
    echo -e "${BLUE}Installing $package with yum...${NC}"
    yum install -y "$yum_package"
}

install_with_dnf() {
    local package="$1"
    local dnf_package="$2"
    
    echo -e "${BLUE}Installing $package with dnf...${NC}"
    dnf install -y "$dnf_package"
}

install_with_pacman() {
    local package="$1"
    local pacman_package="$2"
    
    echo -e "${BLUE}Installing $package with pacman...${NC}"
    pacman -S --noconfirm "$pacman_package"
}

install_with_apk() {
    local package="$1"
    local apk_package="$2"
    
    echo -e "${BLUE}Installing $package with apk...${NC}"
    apk add "$apk_package"
}

# Install yq (special case - often needs manual installation)
install_yq() {
    echo -e "${BLUE}Installing yq...${NC}"
    
    if [[ "$PKG_MANAGER" == "brew" ]]; then
        brew install yq
    else
        # Install via binary download for Linux
        YQ_VERSION="v4.35.2"
        ARCH=$(uname -m)
        case "$ARCH" in
            x86_64)
                YQ_BINARY="yq_linux_amd64"
                ;;
            aarch64|arm64)
                YQ_BINARY="yq_linux_arm64"
                ;;
            *)
                echo -e "${RED}❌ Unsupported architecture for yq: $ARCH${NC}"
                return 1
                ;;
        esac
        
        curl -sL "https://github.com/mikefarah/yq/releases/download/${YQ_VERSION}/${YQ_BINARY}" -o /tmp/yq
        chmod +x /tmp/yq
        mv /tmp/yq /usr/local/bin/yq
        echo -e "${GREEN}✅ yq installed successfully${NC}"
    fi
}

# Install GitHub CLI
install_gh() {
    echo -e "${BLUE}Installing GitHub CLI...${NC}"
    
    case "$PKG_MANAGER" in
        apt)
            curl -fsSL https://cli.github.com/packages/githubcli-archive-keyring.gpg | dd of=/usr/share/keyrings/githubcli-archive-keyring.gpg
            chmod go+r /usr/share/keyrings/githubcli-archive-keyring.gpg
            echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/githubcli-archive-keyring.gpg] https://cli.github.com/packages stable main" | tee /etc/apt/sources.list.d/github-cli.list > /dev/null
            apt update
            apt install gh -y
            ;;
        yum|dnf)
            dnf install -y 'dnf-command(config-manager)'
            dnf config-manager --add-repo https://cli.github.com/packages/rpm/gh-cli.repo
            dnf install -y gh
            ;;
        brew)
            brew install gh
            ;;
        pacman)
            pacman -S --noconfirm github-cli
            ;;
        *)
            echo -e "${YELLOW}⚠️  Please install GitHub CLI manually from: https://cli.github.com${NC}"
            return 1
            ;;
    esac
}

# Install ripgrep
install_ripgrep() {
    echo -e "${BLUE}Installing ripgrep...${NC}"
    
    case "$PKG_MANAGER" in
        apt)
            install_with_apt "ripgrep" "ripgrep"
            ;;
        brew)
            install_with_brew "ripgrep" "ripgrep"
            ;;
        yum|dnf)
            install_with_dnf "ripgrep" "ripgrep"
            ;;
        pacman)
            install_with_pacman "ripgrep" "ripgrep"
            ;;
        apk)
            install_with_apk "ripgrep" "ripgrep"
            ;;
        *)
            # Install via binary
            RIPGREP_VERSION="14.0.3"
            ARCH=$(uname -m)
            case "$ARCH" in
                x86_64)
                    RG_ARCH="x86_64"
                    ;;
                aarch64|arm64)
                    RG_ARCH="aarch64"
                    ;;
                *)
                    echo -e "${YELLOW}⚠️  Cannot auto-install ripgrep for architecture: $ARCH${NC}"
                    return 1
                    ;;
            esac
            
            curl -sL "https://github.com/BurntSushi/ripgrep/releases/download/${RIPGREP_VERSION}/ripgrep-${RIPGREP_VERSION}-${RG_ARCH}-unknown-linux-musl.tar.gz" | tar xz -C /tmp
            mv /tmp/ripgrep-*/rg /usr/local/bin/
            ;;
    esac
}

# Main installation function
install_package() {
    local tool="$1"
    local required_version="$2"
    local package_map="$3"
    
    echo ""
    echo "Checking $tool..."
    
    if command_exists "$tool"; then
        current_version=$(get_version "$tool")
        echo -e "${GREEN}✅ $tool is installed (version: $current_version)${NC}"
        
        if [ -n "$required_version" ] && ! version_ge "$current_version" "$required_version"; then
            echo -e "${YELLOW}⚠️  Version $current_version is below required $required_version${NC}"
            echo -e "${YELLOW}   Attempting upgrade...${NC}"
        else
            return 0
        fi
    else
        echo -e "${YELLOW}⚠️  $tool is not installed${NC}"
    fi
    
    # Special cases
    case "$tool" in
        yq)
            install_yq
            ;;
        gh)
            install_gh
            ;;
        rg|ripgrep)
            install_ripgrep
            ;;
        *)
            # Standard package installation
            case "$PKG_MANAGER" in
                apt)
                    install_with_apt "$tool" "${package_map:-$tool}"
                    ;;
                brew)
                    install_with_brew "$tool" "${package_map:-$tool}"
                    ;;
                yum)
                    install_with_yum "$tool" "${package_map:-$tool}"
                    ;;
                dnf)
                    install_with_dnf "$tool" "${package_map:-$tool}"
                    ;;
                pacman)
                    install_with_pacman "$tool" "${package_map:-$tool}"
                    ;;
                apk)
                    install_with_apk "$tool" "${package_map:-$tool}"
                    ;;
            esac
            ;;
    esac
    
    # Verify installation
    if command_exists "$tool"; then
        new_version=$(get_version "$tool")
        echo -e "${GREEN}✅ $tool installed successfully (version: $new_version)${NC}"
    else
        echo -e "${RED}❌ Failed to install $tool${NC}"
        return 1
    fi
}

# Main execution
main() {
    detect_os
    check_privileges
    
    echo "======================================"
    echo "Installing Required Tools"
    echo "======================================"
    
    # Required tools with minimum versions
    declare -A REQUIRED_TOOLS=(
        ["yq"]="4.30"
        ["git"]="2.25"
        ["jq"]="1.6"
        ["bash"]="4.4"
        ["curl"]=""
        ["gh"]=""
    )
    
    # Package name mappings for different package managers
    # Format: tool:apt_name:brew_name:yum_name
    declare -A PACKAGE_MAPS=(
        ["git"]="git"
        ["jq"]="jq"
        ["bash"]="bash"
        ["curl"]="curl"
    )
    
    FAILED_INSTALLS=()
    
    # Install required tools
    for tool in "${!REQUIRED_TOOLS[@]}"; do
        if ! install_package "$tool" "${REQUIRED_TOOLS[$tool]}" "${PACKAGE_MAPS[$tool]:-$tool}"; then
            FAILED_INSTALLS+=("$tool")
        fi
    done
    
    echo ""
    echo "======================================"
    echo "Installing Optional Tools"
    echo "======================================"
    
    # Optional tools
    declare -A OPTIONAL_TOOLS=(
        ["rg"]=""
        ["fd"]=""
        ["tree"]=""
        ["make"]=""
    )
    
    for tool in "${!OPTIONAL_TOOLS[@]}"; do
        install_package "$tool" "${OPTIONAL_TOOLS[$tool]}" "${PACKAGE_MAPS[$tool]:-$tool}" || true
    done
    
    echo ""
    echo "======================================"
    echo "Installation Summary"
    echo "======================================"
    
    if [ ${#FAILED_INSTALLS[@]} -eq 0 ]; then
        echo -e "${GREEN}✅ All required tools installed successfully!${NC}"
        echo ""
        echo "You can now run: ./setup.sh"
    else
        echo -e "${RED}❌ Some required tools failed to install:${NC}"
        for tool in "${FAILED_INSTALLS[@]}"; do
            echo -e "${RED}   - $tool${NC}"
        done
        echo ""
        echo -e "${YELLOW}Please install these tools manually before running setup.sh${NC}"
        exit 1
    fi
    
    # Run verification
    echo ""
    echo "Running verification..."
    if [ -f "$(dirname "$0")/check-requirements.sh" ]; then
        bash "$(dirname "$0")/check-requirements.sh"
    fi
}

# Run main function
main "$@"