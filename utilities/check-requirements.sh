#!/bin/bash

# Software Factory 2.0 - Requirements Checker
# Verifies all required dependencies are installed and meet version requirements

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo "======================================"
echo "Software Factory 2.0 Requirements Check"
echo "======================================"
echo ""

# Track failures
MISSING_REQUIRED=()
MISSING_OPTIONAL=()
VERSION_ISSUES=()
TOTAL_ISSUES=0

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
            jq --version 2>/dev/null | grep -oE '[0-9]+\.[0-9]+' | head -1 || echo "0.0"
            ;;
        bash)
            bash --version | head -1 | grep -oE '[0-9]+\.[0-9]+' | head -1 || echo "0.0"
            ;;
        gh)
            gh --version 2>/dev/null | head -1 | grep -oE '[0-9]+\.[0-9]+' | head -1 || echo "0.0"
            ;;
        rg|ripgrep)
            rg --version 2>/dev/null | head -1 | grep -oE '[0-9]+\.[0-9]+' | head -1 || echo "0.0"
            ;;
        fd)
            fd --version 2>/dev/null | grep -oE '[0-9]+\.[0-9]+' | head -1 || echo "0.0"
            ;;
        tree)
            tree --version 2>/dev/null | grep -oE '[0-9]+\.[0-9]+' | head -1 || echo "0.0"
            ;;
        make)
            make --version 2>/dev/null | head -1 | grep -oE '[0-9]+\.[0-9]+' | head -1 || echo "0.0"
            ;;
        *)
            echo "0.0"
            ;;
    esac
}

# Version comparison (returns 0 if version1 >= version2)
version_ge() {
    [ "$(printf '%s\n' "$1" "$2" | sort -V | head -n1)" = "$2" ]
}

# Check a single requirement
check_requirement() {
    local tool="$1"
    local required_version="$2"
    local is_required="$3"
    local description="$4"
    
    echo -n "Checking $tool... "
    
    if command_exists "$tool"; then
        if [ -n "$required_version" ]; then
            current_version=$(get_version "$tool")
            if version_ge "$current_version" "$required_version"; then
                echo -e "${GREEN}✅ $tool v$current_version (>= v$required_version) - $description${NC}"
            else
                echo -e "${YELLOW}⚠️  $tool v$current_version (need >= v$required_version) - $description${NC}"
                VERSION_ISSUES+=("$tool needs upgrade: v$current_version -> v$required_version")
                if [ "$is_required" = "true" ]; then
                    ((TOTAL_ISSUES++))
                fi
            fi
        else
            echo -e "${GREEN}✅ $tool installed - $description${NC}"
        fi
    else
        if [ "$is_required" = "true" ]; then
            echo -e "${RED}❌ $tool NOT FOUND (REQUIRED) - $description${NC}"
            MISSING_REQUIRED+=("$tool")
            ((TOTAL_ISSUES++))
        else
            echo -e "${YELLOW}⚠️  $tool not found (optional) - $description${NC}"
            MISSING_OPTIONAL+=("$tool")
        fi
    fi
}

# Check Git configuration
check_git_config() {
    echo ""
    echo "Checking Git configuration..."
    
    if command_exists git; then
        USER_NAME=$(git config --global user.name 2>/dev/null || echo "")
        USER_EMAIL=$(git config --global user.email 2>/dev/null || echo "")
        
        if [ -z "$USER_NAME" ] || [ -z "$USER_EMAIL" ]; then
            echo -e "${YELLOW}⚠️  Git user configuration incomplete${NC}"
            echo "   Run these commands to configure:"
            [ -z "$USER_NAME" ] && echo "     git config --global user.name \"Your Name\""
            [ -z "$USER_EMAIL" ] && echo "     git config --global user.email \"your.email@example.com\""
        else
            echo -e "${GREEN}✅ Git configured for: $USER_NAME <$USER_EMAIL>${NC}"
        fi
    fi
}

# Check GitHub CLI authentication
check_gh_auth() {
    if command_exists gh; then
        echo ""
        echo "Checking GitHub CLI authentication..."
        
        if gh auth status &>/dev/null; then
            echo -e "${GREEN}✅ GitHub CLI authenticated${NC}"
        else
            echo -e "${YELLOW}⚠️  GitHub CLI not authenticated${NC}"
            echo "   Run: gh auth login"
        fi
    fi
}

# Check for specific Software Factory tools
check_sf_tools() {
    echo ""
    echo "Checking Software Factory specific tools..."
    
    # Check if tmc-pr-line-counter.sh exists (if specified in target-repo-config.yaml)
    if [ -f "target-repo-config.yaml" ]; then
        LINE_COUNTER=$(yq '.tools.line_counter' target-repo-config.yaml 2>/dev/null || echo "")
        if [ -n "$LINE_COUNTER" ] && [ "$LINE_COUNTER" != "null" ]; then
            if [ -f "$LINE_COUNTER" ]; then
                echo -e "${GREEN}✅ Line counter tool found: $LINE_COUNTER${NC}"
            else
                echo -e "${YELLOW}⚠️  Line counter tool not found: $LINE_COUNTER${NC}"
                echo "   This tool is needed for size compliance checking"
            fi
        fi
    fi
    
    # Check if /efforts directory exists or can be created
    if [ -d "/efforts" ]; then
        echo -e "${GREEN}✅ /efforts directory exists${NC}"
    else
        if mkdir -p /efforts 2>/dev/null; then
            echo -e "${GREEN}✅ /efforts directory created${NC}"
            rmdir /efforts 2>/dev/null  # Clean up test
        else
            echo -e "${YELLOW}⚠️  Cannot create /efforts directory (may need sudo)${NC}"
            echo "   Run: sudo mkdir -p /efforts && sudo chown $USER:$USER /efforts"
        fi
    fi
}

# Main execution
main() {
    echo "======================================"
    echo "Required Tools"
    echo "======================================"
    
    # Check required tools
    check_requirement "yq" "4.30" "true" "YAML processor for state management"
    check_requirement "git" "2.25" "true" "Version control for state persistence"
    check_requirement "jq" "1.6" "true" "JSON processor for configuration"
    check_requirement "bash" "4.4" "true" "Shell for orchestration scripts"
    check_requirement "curl" "" "true" "For downloading dependencies"
    check_requirement "gh" "" "true" "GitHub CLI for PR creation"
    
    echo ""
    echo "======================================"
    echo "Optional Tools"
    echo "======================================"
    
    # Check optional tools
    check_requirement "rg" "" "false" "Fast file searching (ripgrep)"
    check_requirement "fd" "" "false" "Modern find alternative"
    check_requirement "tree" "" "false" "Directory structure visualization"
    check_requirement "make" "" "false" "Build automation"
    
    # Additional checks
    check_git_config
    check_gh_auth
    check_sf_tools
    
    echo ""
    echo "======================================"
    echo "Summary"
    echo "======================================"
    
    if [ ${#MISSING_REQUIRED[@]} -eq 0 ] && [ ${#VERSION_ISSUES[@]} -eq 0 ]; then
        echo -e "${GREEN}✅ All required dependencies are satisfied!${NC}"
        echo ""
        echo "You're ready to run: ./setup.sh"
    else
        if [ ${#MISSING_REQUIRED[@]} -gt 0 ]; then
            echo -e "${RED}❌ Missing required tools:${NC}"
            for tool in "${MISSING_REQUIRED[@]}"; do
                echo -e "${RED}   - $tool${NC}"
            done
            echo ""
        fi
        
        if [ ${#VERSION_ISSUES[@]} -gt 0 ]; then
            echo -e "${YELLOW}⚠️  Version issues:${NC}"
            for issue in "${VERSION_ISSUES[@]}"; do
                echo -e "${YELLOW}   - $issue${NC}"
            done
            echo ""
        fi
        
        echo -e "${YELLOW}Run './utilities/install-requirements.sh' to install missing dependencies${NC}"
        exit 1
    fi
    
    if [ ${#MISSING_OPTIONAL[@]} -gt 0 ]; then
        echo ""
        echo -e "${BLUE}ℹ️  Optional tools not installed:${NC}"
        for tool in "${MISSING_OPTIONAL[@]}"; do
            echo "   - $tool"
        done
        echo ""
        echo "These tools enhance functionality but are not required."
    fi
    
    # Performance tips
    echo ""
    echo "======================================"
    echo "Performance Tips"
    echo "======================================"
    echo "• Use ripgrep (rg) for faster file searching"
    echo "• Use fd for faster file finding"
    echo "• Ensure /efforts directory has fast disk access"
    echo "• Configure Git to use SSH keys for authentication"
    echo "• Set up GitHub CLI with proper scopes for repo access"
    
    exit 0
}

# Run main function
main "$@"