#!/bin/bash

# Software Factory 2.0 - Ensure jsonschema Module Availability
# This script ensures the jsonschema Python module is available for all
# Software Factory operations, including agent-spawned bash sessions.
#
# Strategy:
#   1. Check if jsonschema is available
#   2. If not, create a dedicated venv for Software Factory
#   3. Install jsonschema in the venv
#   4. Configure shell initialization to activate venv
#
# Usage: source tools/ensure-jsonschema.sh

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
BOLD='\033[1m'
NC='\033[0m' # No Color

# Software Factory virtual environment location
SF_VENV_DIR="${HOME}/.software-factory-venv"
SF_VENV_MARKER="${SF_VENV_DIR}/.sf-venv-marker"

# Function to print colored messages
print_error() {
    echo -e "${RED}${BOLD}❌ ERROR:${NC} ${RED}$1${NC}" >&2
}

print_success() {
    echo -e "${GREEN}${BOLD}✅ PROJECT_DONE:${NC} ${GREEN}$1${NC}"
}

print_warning() {
    echo -e "${YELLOW}${BOLD}⚠ WARNING:${NC} ${YELLOW}$1${NC}"
}

print_info() {
    echo -e "${BLUE}${BOLD}ℹ INFO:${NC} ${BLUE}$1${NC}"
}

# Function to check if jsonschema is available
check_jsonschema() {
    python3 -c "import jsonschema" 2>/dev/null
}

# Function to test jsonschema in clean environment
test_clean_environment() {
    env -i HOME="$HOME" PATH="/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin" \
        bash -c "python3 -c 'import jsonschema' 2>/dev/null"
}

# Function to create and setup Software Factory venv
setup_sf_venv() {
    print_info "Creating Software Factory virtual environment at $SF_VENV_DIR..."

    # Remove old venv if it exists
    if [ -d "$SF_VENV_DIR" ]; then
        print_warning "Removing existing venv..."
        rm -rf "$SF_VENV_DIR"
    fi

    # Create new venv
    python3 -m venv "$SF_VENV_DIR"

    # Activate venv and install jsonschema
    source "$SF_VENV_DIR/bin/activate"

    print_info "Installing jsonschema in venv..."
    pip install --upgrade pip >/dev/null 2>&1
    pip install jsonschema >/dev/null 2>&1

    # Create marker file with installation info
    cat > "$SF_VENV_MARKER" << EOF
# Software Factory Virtual Environment
created_at: $(date -u +"%Y-%m-%dT%H:%M:%SZ")
python_version: $(python3 --version 2>&1)
jsonschema_version: $(python3 -c "import jsonschema; import importlib.metadata; print(importlib.metadata.version('jsonschema'))" 2>/dev/null || echo "unknown")
purpose: Ensures jsonschema availability for Software Factory operations
EOF

    # Deactivate for now
    deactivate 2>/dev/null || true

    print_success "Virtual environment created successfully"
}

# Function to add venv activation to shell initialization
configure_shell_init() {
    local shell_rc=""
    local shell_profile=""

    # Determine which files to update based on shell
    if [ -n "$BASH_VERSION" ]; then
        shell_rc="$HOME/.bashrc"
        shell_profile="$HOME/.bash_profile"
    elif [ -n "$ZSH_VERSION" ]; then
        shell_rc="$HOME/.zshrc"
        shell_profile="$HOME/.zprofile"
    else
        shell_rc="$HOME/.profile"
        shell_profile="$HOME/.profile"
    fi

    # Software Factory venv activation snippet
    local activation_snippet='
# Software Factory 2.0 - Auto-activate virtual environment for jsonschema
if [ -f "$HOME/.software-factory-venv/bin/activate" ]; then
    # Check if we are not already in a venv to avoid conflicts
    if [ -z "$VIRTUAL_ENV" ]; then
        # Silently activate Software Factory venv
        source "$HOME/.software-factory-venv/bin/activate" 2>/dev/null
        # Export marker so we know SF venv is active
        export SF_VENV_ACTIVE=1
    fi
fi'

    # Function to add snippet to file if not already present
    add_to_file() {
        local file="$1"

        if [ ! -f "$file" ]; then
            print_info "Creating $file..."
            touch "$file"
        fi

        # Check if snippet already exists
        if grep -q "Software Factory 2.0 - Auto-activate virtual environment" "$file" 2>/dev/null; then
            print_info "Activation already configured in $file"
        else
            print_info "Adding venv activation to $file..."
            echo "$activation_snippet" >> "$file"
            print_success "Updated $file"
        fi
    }

    # Add to both rc and profile files for maximum compatibility
    for file in "$shell_rc" "$shell_profile"; do
        if [ -n "$file" ] && [ "$file" != "$HOME/" ]; then
            add_to_file "$file"
        fi
    done

    # Also add to .profile for sh compatibility
    if [ "$shell_rc" != "$HOME/.profile" ]; then
        add_to_file "$HOME/.profile"
    fi
}

# Function to test the final setup
test_final_setup() {
    print_info "Testing jsonschema availability in various contexts..."

    local all_pass=true

    # Test 1: Current shell
    echo -n "  Current shell: "
    if check_jsonschema; then
        echo -e "${GREEN}✓${NC}"
    else
        echo -e "${RED}✗${NC}"
        all_pass=false
    fi

    # Test 2: New bash session
    echo -n "  New bash session: "
    if bash -c "source $HOME/.bashrc 2>/dev/null; python3 -c 'import jsonschema' 2>/dev/null"; then
        echo -e "${GREEN}✓${NC}"
    else
        echo -e "${RED}✗${NC}"
        all_pass=false
    fi

    # Test 3: Login shell
    echo -n "  Login shell: "
    if bash -l -c "python3 -c 'import jsonschema' 2>/dev/null"; then
        echo -e "${GREEN}✓${NC}"
    else
        echo -e "${RED}✗${NC}"
        all_pass=false
    fi

    # Test 4: sh shell
    echo -n "  sh shell: "
    if sh -c ". $HOME/.profile 2>/dev/null; python3 -c 'import jsonschema' 2>/dev/null"; then
        echo -e "${GREEN}✓${NC}"
    else
        echo -e "${RED}✗${NC}"
        all_pass=false
    fi

    if [ "$all_pass" = true ]; then
        print_success "All tests passed!"
        return 0
    else
        print_warning "Some tests failed. Manual intervention may be needed."
        return 1
    fi
}

# Main execution
main() {
    echo -e "${CYAN}${BOLD}========================================${NC}"
    echo -e "${CYAN}${BOLD}Software Factory jsonschema Setup${NC}"
    echo -e "${CYAN}${BOLD}========================================${NC}"
    echo ""

    # Step 1: Check current availability
    print_info "Checking current jsonschema availability..."

    if check_jsonschema; then
        print_success "jsonschema is available in current environment"

        # Test clean environment
        if test_clean_environment; then
            print_success "jsonschema is available in clean environment"
            echo ""
            print_info "No additional setup needed - jsonschema is properly installed"

            # Still offer to create venv for extra reliability
            echo ""
            echo -e "${YELLOW}Would you like to create a dedicated Software Factory venv anyway?${NC}"
            echo -e "${YELLOW}This ensures maximum compatibility across all environments.${NC}"
            echo -ne "${CYAN}Create venv? (y/n) [n]: ${NC}"
            read -r response

            if [ "$response" = "y" ] || [ "$response" = "Y" ]; then
                setup_sf_venv
                configure_shell_init
                test_final_setup
            else
                print_info "Keeping current setup"
            fi
        else
            print_warning "jsonschema not available in clean environment"
            print_info "Setting up Software Factory venv for reliability..."
            setup_sf_venv
            configure_shell_init
            test_final_setup
        fi
    else
        print_warning "jsonschema not found in current environment"

        # Check if Python is available
        if ! command -v python3 >/dev/null 2>&1; then
            print_error "Python3 is not installed!"
            echo "Please install Python3 first:"
            echo "  apt-get update && apt-get install -y python3 python3-venv python3-pip"
            return 1
        fi

        print_info "Setting up Software Factory venv..."
        setup_sf_venv
        configure_shell_init
        test_final_setup
    fi

    echo ""
    echo -e "${CYAN}${BOLD}========================================${NC}"
    print_success "Setup complete!"
    echo ""
    echo -e "${YELLOW}Next steps:${NC}"
    echo "  1. Source your shell configuration:"
    echo "     ${GREEN}source ~/.bashrc${NC}  (or source ~/.zshrc for zsh)"
    echo "  2. Verify jsonschema is available:"
    echo "     ${GREEN}python3 -c 'import jsonschema; print(\"jsonschema ready!\")'${NC}"
    echo "  3. Test state validation:"
    echo "     ${GREEN}bash tools/validate-state.sh${NC}"
    echo ""

    # Check if we should also update setup.sh and upgrade.sh
    if [ -f "$(dirname "${BASH_SOURCE[0]}")/setup.sh" ]; then
        echo -e "${YELLOW}Note: You may want to integrate this into setup.sh and upgrade.sh${NC}"
        echo -e "${YELLOW}See the implementation in ensure-jsonschema.sh for details.${NC}"
    fi
}

# Run main function
main "$@"