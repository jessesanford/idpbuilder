#!/bin/bash

# Software Factory 2.0 Installation Verification Script
# Checks that all components are properly installed and configured

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
BOLD='\033[1m'
NC='\033[0m' # No Color

echo -e "${CYAN}${BOLD}"
cat << "EOF"
╔═══════════════════════════════════════════════════════════════════╗
║         Software Factory 2.0 Installation Verification           ║
╚═══════════════════════════════════════════════════════════════════╝
EOF
echo -e "${NC}"

# Get project directory
PROJECT_DIR="${1:-$(pwd)}"

echo -e "${CYAN}Checking SF 2.0 installation in: ${BOLD}$PROJECT_DIR${NC}\n"

# Track issues
ISSUES=0
WARNINGS=0

# Function to check file exists
check_file() {
    local file="$1"
    local description="$2"
    local required="${3:-true}"
    
    if [ -f "$PROJECT_DIR/$file" ]; then
        echo -e "${GREEN}✓${NC} $description"
        return 0
    else
        if [ "$required" = "true" ]; then
            echo -e "${RED}✗${NC} $description ${RED}(MISSING)${NC}"
            ((ISSUES++))
        else
            echo -e "${YELLOW}⚠${NC} $description ${YELLOW}(Optional)${NC}"
            ((WARNINGS++))
        fi
        return 1
    fi
}

# Function to check directory exists
check_dir() {
    local dir="$1"
    local description="$2"
    local required="${3:-true}"
    
    if [ -d "$PROJECT_DIR/$dir" ]; then
        echo -e "${GREEN}✓${NC} $description"
        return 0
    else
        if [ "$required" = "true" ]; then
            echo -e "${RED}✗${NC} $description ${RED}(MISSING)${NC}"
            ((ISSUES++))
        else
            echo -e "${YELLOW}⚠${NC} $description ${YELLOW}(Optional)${NC}"
            ((WARNINGS++))
        fi
        return 1
    fi
}

# Function to check for hardcoded paths
check_paths() {
    local file="$1"
    local description="$2"
    
    if [ -f "$PROJECT_DIR/$file" ]; then
        if grep -q "/workspaces/software-factory-2.0-template" "$PROJECT_DIR/$file" 2>/dev/null; then
            echo -e "${YELLOW}⚠${NC} $description has template paths ${YELLOW}(Should be updated)${NC}"
            ((WARNINGS++))
            return 1
        else
            echo -e "${GREEN}✓${NC} $description paths configured"
            return 0
        fi
    fi
}

echo -e "${BOLD}Core Structure:${NC}"
check_dir "rule-library" "Rule library (critical rules)"
check_dir "agent-states" "Agent state configurations"
check_dir "state-machines" "State machine definitions"
check_dir "expertise" "Expertise modules"
check_dir "hooks" "Pre/post-compaction hooks"
check_dir "quick-reference" "Quick reference guides"
check_dir "checkpoints" "Checkpoint directory"
check_dir "todos" "TODO preservation directory"

echo -e "\n${BOLD}Claude Configuration:${NC}"
check_dir ".claude" "Claude configuration directory"
check_dir ".claude/commands" "Slash commands"
check_dir ".claude/agents" "Agent configurations"
check_file ".claude/settings.json" "Settings file"

echo -e "\n${BOLD}Agent Configurations:${NC}"
check_file ".claude/agents/orchestrator.md" "Orchestrator config"
check_file ".claude/agents/sw-engineer.md" "SW Engineer config"
check_file ".claude/agents/code-reviewer.md" "Code Reviewer config"
check_file ".claude/agents/architect.md" "Architect config"

echo -e "\n${BOLD}Slash Commands:${NC}"
check_file ".claude/commands/continue-orchestrating.md" "Orchestrator command"
check_file ".claude/commands/continue-implementing.md" "SW Engineer command"
check_file ".claude/commands/continue-reviewing.md" "Code Reviewer command"
check_file ".claude/commands/continue-architecting.md" "Architect command"
check_file ".claude/commands/check-status.md" "Status check command"
check_file ".claude/commands/reset-state.md" "Reset command"

echo -e "\n${BOLD}Utility Scripts:${NC}"
check_file "utilities/pre-compact.sh" "Pre-compaction utility"
check_file "utilities/post-compact.sh" "Post-compaction utility"
check_file "utilities/todo-preservation.sh" "TODO preservation"
check_file "utilities/state-snapshot.sh" "State snapshot"
check_file "utilities/recovery-assistant.sh" "Recovery assistant"

# Check if utilities are executable
if [ -f "$PROJECT_DIR/utilities/pre-compact.sh" ]; then
    if [ -x "$PROJECT_DIR/utilities/pre-compact.sh" ]; then
        echo -e "${GREEN}✓${NC} Utilities are executable"
    else
        echo -e "${RED}✗${NC} Utilities are not executable ${RED}(Run: chmod +x utilities/*.sh)${NC}"
        ((ISSUES++))
    fi
fi

echo -e "\n${BOLD}Configuration Files:${NC}"
check_file "orchestrator-state-v3.json" "Orchestrator state"
check_file "project-config.yaml" "Project configuration" false

echo -e "\n${BOLD}Path Configuration:${NC}"
check_paths ".claude/settings.json" "Settings.json"
check_paths "utilities/pre-compact.sh" "Pre-compact utility"
check_paths "utilities/post-compact.sh" "Post-compact utility"

# Check settings.json for proper hook paths
if [ -f "$PROJECT_DIR/.claude/settings.json" ]; then
    echo -e "\n${BOLD}Settings Validation:${NC}"
    
    # Check if PreCompact hook is configured
    if grep -q "PreCompact" "$PROJECT_DIR/.claude/settings.json" 2>/dev/null; then
        echo -e "${GREEN}✓${NC} PreCompact hook configured"
    else
        echo -e "${YELLOW}⚠${NC} PreCompact hook may need configuration"
        ((WARNINGS++))
    fi
    
    # Check factory version
    if grep -q '"factory_version": "2.0"' "$PROJECT_DIR/.claude/settings.json" 2>/dev/null; then
        echo -e "${GREEN}✓${NC} Factory version 2.0 confirmed"
    else
        echo -e "${RED}✗${NC} Factory version not set to 2.0"
        ((ISSUES++))
    fi
fi

# Check for Git repository
echo -e "\n${BOLD}Git Configuration:${NC}"
if [ -d "$PROJECT_DIR/.git" ]; then
    echo -e "${GREEN}✓${NC} Git repository initialized"
    
    # Check for SF 2.0 branch
    cd "$PROJECT_DIR"
    if git branch | grep -q "software-factory-2.0"; then
        echo -e "${GREEN}✓${NC} SF 2.0 branch exists"
    else
        echo -e "${YELLOW}⚠${NC} No software-factory-2.0 branch"
        ((WARNINGS++))
    fi
else
    echo -e "${RED}✗${NC} No Git repository"
    ((ISSUES++))
fi

# Summary
echo -e "\n${BOLD}═══════════════════════════════════════════════════════${NC}"
echo -e "${BOLD}Verification Summary:${NC}"

if [ $ISSUES -eq 0 ] && [ $WARNINGS -eq 0 ]; then
    echo -e "${GREEN}${BOLD}✅ Perfect! Software Factory 2.0 is fully configured.${NC}"
elif [ $ISSUES -eq 0 ]; then
    echo -e "${GREEN}✓ Installation complete with $WARNINGS warnings${NC}"
    echo -e "${YELLOW}Review warnings above for optional improvements${NC}"
else
    echo -e "${RED}✗ Found $ISSUES critical issues and $WARNINGS warnings${NC}"
    echo -e "${RED}Please fix critical issues before proceeding${NC}"
fi

echo -e "${BOLD}═══════════════════════════════════════════════════════${NC}"

# Exit with error if critical issues found
if [ $ISSUES -gt 0 ]; then
    exit 1
fi

exit 0