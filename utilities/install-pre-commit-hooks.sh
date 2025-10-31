#!/bin/bash

################################################################################
# Pre-Commit Hook Installation Script
#
# This script installs the unified pre-commit hook (master-pre-commit.sh) which:
# - Auto-detects SF 2.0 vs SF 3.0
# - Auto-detects repo type (planning/effort/template)
# - Applies appropriate validation rules
#
# Usage: bash utilities/install-pre-commit-hooks.sh
################################################################################

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
BOLD='\033[1m'
NC='\033[0m' # No Color

echo -e "${BLUE}${BOLD}═══════════════════════════════════════════════════════${NC}"
echo -e "${BLUE}${BOLD}     Pre-Commit Hook Installation Script${NC}"
echo -e "${BLUE}${BOLD}═══════════════════════════════════════════════════════${NC}"

# Find the project directory (where orchestrator-state-v3.json or orchestrator-state-v3.json lives)
if [ -f "orchestrator-state-v3.json" ] || [ -f "orchestrator-state-v3.json" ]; then
    PROJECT_DIR="$(pwd)"
elif [ -f "../orchestrator-state-v3.json" ] || [ -f "../orchestrator-state-v3.json" ]; then
    PROJECT_DIR="$(dirname $(pwd))"
elif [ -n "$CLAUDE_PROJECT_DIR" ]; then
    PROJECT_DIR="$CLAUDE_PROJECT_DIR"
else
    echo -e "${RED}❌ ERROR: Cannot find orchestrator state file${NC}"
    echo "Please run this script from the project root or set CLAUDE_PROJECT_DIR"
    echo "Looking for: orchestrator-state-v3.json (SF 2.0) or orchestrator-state-v3.json (SF 3.0)"
    exit 1
fi

echo -e "${GREEN}✓${NC} Project directory: $PROJECT_DIR"

# Detect SF version
SF_VERSION="unknown"
if [ -f "$PROJECT_DIR/orchestrator-state-v3.json" ] || [ -f "$PROJECT_DIR/bug-tracking.json" ]; then
    SF_VERSION="3.0"
    echo -e "${GREEN}✓${NC} Detected Software Factory 3.0"
elif [ -f "$PROJECT_DIR/orchestrator-state-v3.json" ]; then
    SF_VERSION="2.0"
    echo -e "${GREEN}✓${NC} Detected Software Factory 2.0 (will upgrade to 3.0)"
else
    echo -e "${YELLOW}⚠${NC} Cannot determine SF version"
fi

# Check if master pre-commit hook exists
MASTER_HOOK="$PROJECT_DIR/tools/git-commit-hooks/master-pre-commit.sh"

if [ ! -f "$MASTER_HOOK" ]; then
    echo -e "${RED}❌ ERROR: Master pre-commit hook not found at:${NC}"
    echo "   $MASTER_HOOK"
    echo ""
    echo -e "${YELLOW}This hook is required for SF 2.0/3.0 validation.${NC}"
    echo "Please ensure you have the latest template files."
    exit 1
fi

echo -e "${GREEN}✓${NC} Master pre-commit hook found"

# Function to install unified hook
install_hook() {
    local repo_path="$1"
    local hook_path="$repo_path/.git/hooks/pre-commit"

    if [ ! -d "$repo_path/.git" ]; then
        echo -e "${YELLOW}⚠${NC} Not a git repository, skipping: $repo_path"
        return 1
    fi

    # Create hooks directory if it doesn't exist
    mkdir -p "$repo_path/.git/hooks"

    # Copy unified hook
    cp "$MASTER_HOOK" "$hook_path"
    chmod +x "$hook_path"

    echo -e "${GREEN}✓${NC} Hook installed: $repo_path"
    return 0
}

echo -e "\n${BLUE}${BOLD}════════════════════════════════════════════════════════${NC}"
echo -e "${BLUE}${BOLD}     Step 1: Install Planning/Main Repository Hook${NC}"
echo -e "${BLUE}${BOLD}════════════════════════════════════════════════════════${NC}"

# Install hook for the planning repository (current project)
install_hook "$PROJECT_DIR"

echo -e "\n${BLUE}${BOLD}════════════════════════════════════════════════════════${NC}"
echo -e "${BLUE}${BOLD}     Step 2: Install Effort Repository Hooks${NC}"
echo -e "${BLUE}${BOLD}════════════════════════════════════════════════════════${NC}"

# Count installed hooks
EFFORT_COUNT=0

# Check for efforts in the efforts/ directory structure (SF 2.0 pattern)
if [ -d "$PROJECT_DIR/efforts" ]; then
    echo -e "\n${BLUE}Searching for effort repositories...${NC}"

    # Find all directories that look like effort repos (have .git directories)
    while IFS= read -r -d '' git_dir; do
        effort_dir="$(dirname "$git_dir")"
        if install_hook "$effort_dir"; then
            ((EFFORT_COUNT++))
        fi
    done < <(find "$PROJECT_DIR/efforts" -type d -name ".git" -print0 2>/dev/null)

    if [ $EFFORT_COUNT -eq 0 ]; then
        echo -e "${YELLOW}⚠${NC} No effort repositories found in efforts/"
    fi
else
    echo -e "${YELLOW}⚠${NC} No efforts/ directory found"
fi

# Also check if there are effort paths in pre_planned_infrastructure (SF 2.0)
if [ "$SF_VERSION" = "2.0" ] && command -v jq &> /dev/null; then
    echo -e "\n${BLUE}Checking pre_planned_infrastructure for additional repositories...${NC}"

    EFFORT_PATHS=$(jq -r '.pre_planned_infrastructure.efforts[]?.full_path // empty' "$PROJECT_DIR/orchestrator-state-v3.json" 2>/dev/null)

    if [ -n "$EFFORT_PATHS" ]; then
        while IFS= read -r effort_path; do
            # Make path absolute if it's relative
            if [[ ! "$effort_path" = /* ]]; then
                effort_path="$PROJECT_DIR/$effort_path"
            fi

            # Remove trailing slash
            effort_path="${effort_path%/}"

            if [ -d "$effort_path" ] && [ -d "$effort_path/.git" ]; then
                if install_hook "$effort_path"; then
                    ((EFFORT_COUNT++))
                fi
            elif [ -d "$effort_path" ]; then
                echo -e "${YELLOW}⚠${NC} Directory exists but not a git repo yet: $effort_path"
            fi
        done <<< "$EFFORT_PATHS"
    fi
fi

echo -e "\n${BLUE}${BOLD}════════════════════════════════════════════════════════${NC}"
echo -e "${BLUE}${BOLD}     Step 3: Check for Target Repository${NC}"
echo -e "${BLUE}${BOLD}════════════════════════════════════════════════════════${NC}"

# Check if there's a separate target repository configured
if [ -f "$PROJECT_DIR/target-repo-config.yaml" ]; then
    if command -v yq &> /dev/null; then
        TARGET_REPO_PATH=$(yq -r '.repository_path // empty' "$PROJECT_DIR/target-repo-config.yaml" 2>/dev/null)

        if [ -n "$TARGET_REPO_PATH" ] && [ -d "$TARGET_REPO_PATH" ]; then
            echo -e "\n${BLUE}Found target repository at: $TARGET_REPO_PATH${NC}"

            # Target repository should get the hook
            install_hook "$TARGET_REPO_PATH"

            # Also install hooks for any effort directories within the target repo
            if [ -d "$TARGET_REPO_PATH/efforts" ]; then
                while IFS= read -r -d '' git_dir; do
                    effort_dir="$(dirname "$git_dir")"
                    if install_hook "$effort_dir"; then
                        ((EFFORT_COUNT++))
                    fi
                done < <(find "$TARGET_REPO_PATH/efforts" -type d -name ".git" -print0 2>/dev/null)
            fi
        fi
    else
        echo -e "${YELLOW}⚠${NC} yq not installed, cannot parse target-repo-config.yaml"
    fi
else
    echo -e "${YELLOW}⚠${NC} No target-repo-config.yaml found (not required)"
fi

echo -e "\n${BLUE}${BOLD}════════════════════════════════════════════════════════${NC}"
echo -e "${BLUE}${BOLD}     Installation Summary${NC}"
echo -e "${BLUE}${BOLD}════════════════════════════════════════════════════════${NC}"

echo -e "${GREEN}✓${NC} Software Factory Version: $SF_VERSION"
echo -e "${GREEN}✓${NC} Planning/main repository hook installed: 1"
echo -e "${GREEN}✓${NC} Effort repository hooks installed: $EFFORT_COUNT"

# Verify installations
echo -e "\n${BLUE}Verifying installations...${NC}"

# Check planning repo
PLANNING_INSTALLED="NO"
if [ -f "$PROJECT_DIR/.git/hooks/pre-commit" ]; then
    if grep -qi "master.pre-commit\|Software Factory.*Validation" "$PROJECT_DIR/.git/hooks/pre-commit" 2>/dev/null; then
        PLANNING_INSTALLED="YES"
        echo -e "${GREEN}✓${NC} Planning/main repository hook verified"
    else
        echo -e "${RED}✗${NC} Planning/main repository has unexpected hook content"
    fi
else
    echo -e "${RED}✗${NC} Planning/main repository hook not found"
fi

echo -e "\n${GREEN}${BOLD}═══════════════════════════════════════════════════════${NC}"
echo -e "${GREEN}${BOLD}     Pre-Commit Hook Installation Complete${NC}"
echo -e "${GREEN}${BOLD}═══════════════════════════════════════════════════════${NC}"

if [ "$PLANNING_INSTALLED" = "YES" ]; then
    echo -e "\n${GREEN}✓ All hooks installed successfully!${NC}"
    echo -e "\nThe following protections are now in place:"
    echo -e "  ${BOLD}Unified Hook Features:${NC}"
    echo -e "  • Auto-detects Software Factory version (2.0 or 3.0)"
    echo -e "  • Auto-detects repository type (planning/effort/template)"
    echo -e "  • SF 2.0: Validates orchestrator-state-v3.json + legacy hooks"
    echo -e "  • SF 3.0: Validates 4-file structure (orchestrator-state-v3.json, bug-tracking.json, etc.)"
    echo -e "  • R506: Enforces no --no-verify usage"
    echo -e "  • R288: Validates atomic state updates"
    echo ""
    exit 0
else
    echo -e "\n${YELLOW}⚠ Some hooks may not have been installed correctly${NC}"
    echo -e "Please check the output above for any errors."
    exit 1
fi
