#!/usr/bin/env bash

# Software Factory 2.0 - Initialization Validation Script
# Validates that initialization completed successfully and all files are ready

set -euo pipefail

PROJECT_PREFIX="${1:-}"
CLAUDE_PROJECT_DIR="${CLAUDE_PROJECT_DIR:-/home/vscode/software-factory-template}"

if [ -z "$PROJECT_PREFIX" ]; then
    echo "Usage: $0 <project_prefix>"
    echo "Example: $0 idpbuilder-registry"
    exit 1
fi

echo "========================================="
echo "SF2.0 INITIALIZATION VALIDATION"
echo "Project: $PROJECT_PREFIX"
echo "========================================="
echo ""

# Track validation results
ERRORS=0
WARNINGS=0

# Color codes for output
RED='\033[0;31m'
YELLOW='\033[1;33m'
GREEN='\033[0;32m'
NC='\033[0m' # No Color

# Helper functions
check_file() {
    local file="$1"
    local description="$2"

    if [ -f "$CLAUDE_PROJECT_DIR/$file" ]; then
        echo -e "${GREEN}✅${NC} $description exists"
        return 0
    else
        echo -e "${RED}❌${NC} $description missing: $file"
        ((ERRORS++))
        return 1
    fi
}

check_directory() {
    local dir="$1"
    local description="$2"

    if [ -d "$CLAUDE_PROJECT_DIR/$dir" ]; then
        echo -e "${GREEN}✅${NC} $description exists"
        return 0
    else
        echo -e "${RED}❌${NC} $description missing: $dir"
        ((ERRORS++))
        return 1
    fi
}

check_yaml_field() {
    local file="$1"
    local field="$2"
    local description="$3"

    if [ ! -f "$CLAUDE_PROJECT_DIR/$file" ]; then
        return 1
    fi

    # Use grep to check for field (simple check)
    if grep -q "^[[:space:]]*${field}:" "$CLAUDE_PROJECT_DIR/$file" 2>/dev/null; then
        # Check if value is not empty or placeholder
        local value=$(grep "^[[:space:]]*${field}:" "$CLAUDE_PROJECT_DIR/$file" | cut -d: -f2- | tr -d ' ')
        if [[ "$value" == *"["*"]"* ]] || [[ "$value" == *"TODO"* ]] || [ -z "$value" ]; then
            echo -e "${YELLOW}⚠️${NC}  $description has placeholder value"
            ((WARNINGS++))
            return 1
        else
            echo -e "${GREEN}✅${NC} $description configured"
            return 0
        fi
    else
        echo -e "${RED}❌${NC} $description missing field: $field"
        ((ERRORS++))
        return 1
    fi
}

check_git_repo() {
    local path="$1"
    local description="$2"

    if [ -d "$CLAUDE_PROJECT_DIR/$path/.git" ]; then
        echo -e "${GREEN}✅${NC} $description is git repository"

        # Check for commits
        cd "$CLAUDE_PROJECT_DIR/$path"
        if git log --oneline -1 &>/dev/null; then
            echo -e "${GREEN}✅${NC} $description has commits"
        else
            echo -e "${YELLOW}⚠️${NC}  $description has no commits"
            ((WARNINGS++))
        fi
        cd - > /dev/null

        return 0
    else
        echo -e "${RED}❌${NC} $description not a git repository"
        ((ERRORS++))
        return 1
    fi
}

# Start validation
echo "1. CHECKING CORE FILES"
echo "----------------------"
check_file "IMPLEMENTATION-PLAN.md" "Implementation plan"
check_file "efforts/$PROJECT_PREFIX/configs/setup-config.yaml" "Setup configuration"
check_file "efforts/$PROJECT_PREFIX/configs/target-repo-config.yaml" "Target repo configuration"
echo ""

echo "2. CHECKING DIRECTORIES"
echo "------------------------"
check_directory "efforts/$PROJECT_PREFIX" "Project directory"
check_directory "efforts/$PROJECT_PREFIX/configs" "Configs directory"

# Check for repository based on type
if [ -d "$CLAUDE_PROJECT_DIR/efforts/$PROJECT_PREFIX/target-repo" ]; then
    check_directory "efforts/$PROJECT_PREFIX/target-repo" "Target repository"
elif [ -d "$CLAUDE_PROJECT_DIR/efforts/$PROJECT_PREFIX/project" ]; then
    check_directory "efforts/$PROJECT_PREFIX/project" "Project directory"
fi
echo ""

echo "3. VALIDATING CONFIGURATION FILES"
echo "----------------------------------"

# Check setup-config.yaml
if [ -f "$CLAUDE_PROJECT_DIR/efforts/$PROJECT_PREFIX/configs/setup-config.yaml" ]; then
    echo "Checking setup-config.yaml:"
    check_yaml_field "efforts/$PROJECT_PREFIX/configs/setup-config.yaml" "name" "Project name"
    check_yaml_field "efforts/$PROJECT_PREFIX/configs/setup-config.yaml" "prefix" "Project prefix"
    check_yaml_field "efforts/$PROJECT_PREFIX/configs/setup-config.yaml" "primary_language" "Primary language"
    check_yaml_field "efforts/$PROJECT_PREFIX/configs/setup-config.yaml" "build_system" "Build system"
    check_yaml_field "efforts/$PROJECT_PREFIX/configs/setup-config.yaml" "framework" "Test framework"
fi
echo ""

# Check target-repo-config.yaml
if [ -f "$CLAUDE_PROJECT_DIR/efforts/$PROJECT_PREFIX/configs/target-repo-config.yaml" ]; then
    echo "Checking target-repo-config.yaml:"
    check_yaml_field "efforts/$PROJECT_PREFIX/configs/target-repo-config.yaml" "type" "Repository type"

    # Check for upstream fork specific fields
    if grep -q "upstream_fork" "$CLAUDE_PROJECT_DIR/efforts/$PROJECT_PREFIX/configs/target-repo-config.yaml"; then
        check_yaml_field "efforts/$PROJECT_PREFIX/configs/target-repo-config.yaml" "url" "Repository URL"
    fi
fi
echo ""

echo "4. VALIDATING IMPLEMENTATION PLAN"
echo "---------------------------------"
if [ -f "$CLAUDE_PROJECT_DIR/IMPLEMENTATION-PLAN.md" ]; then
    # Check for required sections
    for section in "Project Overview" "Goals and Objectives" "Technical Architecture" "Phase 1" "Phase 2" "Phase 3" "Success Criteria" "Risk Mitigation"; do
        if grep -q "$section" "$CLAUDE_PROJECT_DIR/IMPLEMENTATION-PLAN.md"; then
            echo -e "${GREEN}✅${NC} Has section: $section"
        else
            echo -e "${RED}❌${NC} Missing section: $section"
            ((ERRORS++))
        fi
    done

    # Check for effort structure
    if grep -q "Effort [0-9]\.[0-9]\.[0-9]" "$CLAUDE_PROJECT_DIR/IMPLEMENTATION-PLAN.md"; then
        echo -e "${GREEN}✅${NC} Has properly numbered efforts"
    else
        echo -e "${YELLOW}⚠️${NC}  Efforts may not be properly numbered"
        ((WARNINGS++))
    fi
fi
echo ""

echo "5. CHECKING GIT REPOSITORIES"
echo "----------------------------"
if [ -d "$CLAUDE_PROJECT_DIR/efforts/$PROJECT_PREFIX/target-repo" ]; then
    check_git_repo "efforts/$PROJECT_PREFIX/target-repo" "Target repository"

    # Check for upstream remote
    cd "$CLAUDE_PROJECT_DIR/efforts/$PROJECT_PREFIX/target-repo"
    if git remote | grep -q upstream; then
        echo -e "${GREEN}✅${NC} Upstream remote configured"
    else
        echo -e "${YELLOW}⚠️${NC}  No upstream remote configured"
        ((WARNINGS++))
    fi
    cd - > /dev/null

elif [ -d "$CLAUDE_PROJECT_DIR/efforts/$PROJECT_PREFIX/project" ]; then
    check_git_repo "efforts/$PROJECT_PREFIX/project" "Project repository"
fi
echo ""

echo "6. CHECKING AGENT CUSTOMIZATIONS"
echo "--------------------------------"
if grep -q "Project-Specific Expertise" "$CLAUDE_PROJECT_DIR/.claude/agents/sw-engineer.md" 2>/dev/null; then
    echo -e "${GREEN}✅${NC} SW Engineer agent customized"

    # Check if language matches config
    if [ -f "$CLAUDE_PROJECT_DIR/efforts/$PROJECT_PREFIX/configs/setup-config.yaml" ]; then
        lang=$(grep "primary_language:" "$CLAUDE_PROJECT_DIR/efforts/$PROJECT_PREFIX/configs/setup-config.yaml" | cut -d: -f2 | tr -d ' ' || true)
        if [ -n "$lang" ] && grep -q "$lang" "$CLAUDE_PROJECT_DIR/.claude/agents/sw-engineer.md"; then
            echo -e "${GREEN}✅${NC} Language expertise matches: $lang"
        else
            echo -e "${YELLOW}⚠️${NC}  Language expertise may not match configuration"
            ((WARNINGS++))
        fi
    fi
else
    echo -e "${YELLOW}⚠️${NC}  SW Engineer agent not customized"
    ((WARNINGS++))
fi

if grep -q "Project-Specific Review Points" "$CLAUDE_PROJECT_DIR/.claude/agents/code-reviewer.md" 2>/dev/null; then
    echo -e "${GREEN}✅${NC} Code Reviewer agent customized"
else
    echo -e "${YELLOW}⚠️${NC}  Code Reviewer agent not customized"
    ((WARNINGS++))
fi
echo ""

echo "7. CHECKING ORCHESTRATOR STATE"
echo "------------------------------"
if [ -f "$CLAUDE_PROJECT_DIR/orchestrator-state-v3.json" ]; then
    echo -e "${GREEN}✅${NC} Orchestrator state file exists"

    # Check if it's ready for normal operation
    if grep -q '"current_state"[[:space:]]*:[[:space:]]*"INIT"' "$CLAUDE_PROJECT_DIR/orchestrator-state-v3.json"; then
        echo -e "${GREEN}✅${NC} Ready for /continue-orchestrating"
    else
        echo -e "${YELLOW}⚠️${NC}  State may not be INIT"
        ((WARNINGS++))
    fi
else
    echo -e "${YELLOW}⚠️${NC}  Orchestrator state file not created yet"
    echo "    (Will be created during INIT_HANDOFF)"
fi
echo ""

# Summary
echo "========================================="
echo "VALIDATION SUMMARY"
echo "========================================="

if [ $ERRORS -eq 0 ] && [ $WARNINGS -eq 0 ]; then
    echo -e "${GREEN}✅ VALIDATION PASSED${NC}"
    echo "All initialization requirements met!"
    echo ""
    echo "Next step: Run /continue-orchestrating"
    exit 0
elif [ $ERRORS -eq 0 ]; then
    echo -e "${YELLOW}⚠️  VALIDATION PASSED WITH WARNINGS${NC}"
    echo "Errors: 0"
    echo "Warnings: $WARNINGS"
    echo ""
    echo "The project can proceed but review warnings above."
    echo "Next step: Run /continue-orchestrating"
    exit 0
else
    echo -e "${RED}❌ VALIDATION FAILED${NC}"
    echo "Errors: $ERRORS"
    echo "Warnings: $WARNINGS"
    echo ""
    echo "Please fix the errors above before proceeding."
    echo "You may need to run error recovery or restart initialization."
    exit 1
fi