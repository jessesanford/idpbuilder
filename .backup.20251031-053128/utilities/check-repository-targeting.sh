#!/bin/bash

# 🔴🔴🔴 REPOSITORY TARGETING VERIFICATION SCRIPT - R508 ENFORCEMENT
# Purpose: FORCIBLY verify that ALL infrastructure is on the CORRECT repository
# Exit codes:
#   0 = All repositories correct
#   911 = CATASTROPHIC - Wrong repository detected (R508 violation)

set -euo pipefail

# Set CLAUDE_PROJECT_DIR if not already set
: "${CLAUDE_PROJECT_DIR:=/home/vscode/software-factory-template}"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo "🔴🔴🔴 REPOSITORY TARGETING VERIFICATION (R508 SUPREME LAW)"
echo "========================================="
echo "Time: $(date '+%Y-%m-%d %H:%M:%S %Z')"
echo ""

# Track results
ALL_CORRECT=true
CATASTROPHIC_FAILURES=()

# 1. LOAD TARGET REPOSITORY CONFIGURATION
echo "📋 Loading target repository configuration..."

if [ ! -f "$CLAUDE_PROJECT_DIR/target-repo-config.yaml" ]; then
    echo -e "${RED}❌ CRITICAL: target-repo-config.yaml not found!${NC}"
    exit 911
fi

# Extract target repository (supports both Python yq and Go yq)
if command -v yq &> /dev/null; then
    TARGET_REPO=$(cat "$CLAUDE_PROJECT_DIR/target-repo-config.yaml" | yq -r '.repository_url' 2>/dev/null || echo "")
    if [ -z "$TARGET_REPO" ]; then
        TARGET_REPO=$(yq eval '.repository_url' "$CLAUDE_PROJECT_DIR/target-repo-config.yaml" 2>/dev/null || echo "")
    fi
else
    TARGET_REPO=$(grep "repository_url:" "$CLAUDE_PROJECT_DIR/target-repo-config.yaml" | sed 's/.*repository_url:[ ]*//' | tr -d '"' | tr -d "'" || echo "")
fi

if [ -z "$TARGET_REPO" ]; then
    echo -e "${RED}❌ CRITICAL: No repository_url in target-repo-config.yaml${NC}"
    exit 911
fi

echo "✅ Target repository configured: $TARGET_REPO"
echo ""

# Get the current repository's remote URL for comparison
CURRENT_DIR=$(pwd)
if [ -d ".git" ]; then
    CURRENT_REMOTE=$(git remote get-url origin 2>/dev/null || git remote get-url target 2>/dev/null || echo "")
    echo "📍 Current directory remote: $CURRENT_REMOTE"

    if [ "$CURRENT_REMOTE" == "$TARGET_REPO" ]; then
        echo -e "${YELLOW}⚠️ WARNING: Currently in target repository!${NC}"
        echo "   This might be correct for some operations"
    else
        echo "   Current directory is in planning/SF repository (expected)"
    fi
fi
echo ""

# 2. CHECK ALL EFFORT DIRECTORIES
echo "🔍 Checking all effort directories..."

EFFORT_COUNT=0
WRONG_REPO_COUNT=0

# Find all git repositories under efforts/
while IFS= read -r git_dir; do
    REPO_DIR=$(dirname "$git_dir")
    EFFORT_COUNT=$((EFFORT_COUNT + 1))

    echo -n "Checking: $REPO_DIR ... "

    cd "$REPO_DIR"

    # Get remote URL
    REMOTE_URL=$(git remote get-url origin 2>/dev/null || git remote get-url target 2>/dev/null || echo "")

    if [ -z "$REMOTE_URL" ]; then
        echo -e "${YELLOW}No remote configured${NC}"
    elif [ "$REMOTE_URL" != "$TARGET_REPO" ]; then
        echo -e "${RED}❌ WRONG REPOSITORY!${NC}"
        echo -e "${RED}   Expected: $TARGET_REPO${NC}"
        echo -e "${RED}   Actual:   $REMOTE_URL${NC}"
        CATASTROPHIC_FAILURES+=("$REPO_DIR: Wrong repository!")
        ALL_CORRECT=false
        WRONG_REPO_COUNT=$((WRONG_REPO_COUNT + 1))
    else
        echo -e "${GREEN}✅ Correct${NC}"
    fi

    cd "$CLAUDE_PROJECT_DIR"
done < <(find "$CLAUDE_PROJECT_DIR/efforts" -type d -name ".git" 2>/dev/null || true)

echo ""
echo "Total efforts checked: $EFFORT_COUNT"
echo "Wrong repository: $WRONG_REPO_COUNT"
echo ""

# 3. CHECK INTEGRATE_WAVE_EFFORTS DIRECTORIES
echo "🔍 Checking integration directories..."

INTEGRATE_WAVE_EFFORTS_COUNT=0
WRONG_INTEGRATE_WAVE_EFFORTS_COUNT=0

# Find all potential integration directories
for pattern in "*-integration" "*-merge" "integration-workspace"; do
    while IFS= read -r integ_dir; do
        if [ -d "$integ_dir/.git" ]; then
            INTEGRATE_WAVE_EFFORTS_COUNT=$((INTEGRATE_WAVE_EFFORTS_COUNT + 1))

            echo -n "Checking: $integ_dir ... "

            cd "$integ_dir"

            # Get remote URL
            REMOTE_URL=$(git remote get-url origin 2>/dev/null || git remote get-url target 2>/dev/null || echo "")

            if [ -z "$REMOTE_URL" ]; then
                echo -e "${YELLOW}No remote configured${NC}"
            elif [ "$REMOTE_URL" != "$TARGET_REPO" ]; then
                echo -e "${RED}❌ WRONG REPOSITORY!${NC}"
                echo -e "${RED}   Expected: $TARGET_REPO${NC}"
                echo -e "${RED}   Actual:   $REMOTE_URL${NC}"
                CATASTROPHIC_FAILURES+=("$integ_dir: Integration on wrong repository!")
                ALL_CORRECT=false
                WRONG_INTEGRATE_WAVE_EFFORTS_COUNT=$((WRONG_INTEGRATE_WAVE_EFFORTS_COUNT + 1))
            else
                echo -e "${GREEN}✅ Correct${NC}"
            fi

            cd "$CLAUDE_PROJECT_DIR"
        fi
    done < <(find "$CLAUDE_PROJECT_DIR" -type d -name "$pattern" 2>/dev/null || true)
done

echo ""
echo "Total integrations checked: $INTEGRATE_WAVE_EFFORTS_COUNT"
echo "Wrong repository: $WRONG_INTEGRATE_WAVE_EFFORTS_COUNT"
echo ""

# 4. FINAL VERDICT
echo "========================================="
echo "🔴🔴🔴 REPOSITORY TARGETING VERDICT"
echo "========================================="

if $ALL_CORRECT; then
    echo -e "${GREEN}✅ ALL REPOSITORIES CORRECTLY TARGETED${NC}"
    echo "All infrastructure is on the configured target repository"
    echo "R508 Compliance: PASSED"
    exit 0
else
    echo -e "${RED}🔴🔴🔴 CATASTROPHIC FAILURE DETECTED!${NC}"
    echo -e "${RED}SUPREME LAW VIOLATION (R508): Infrastructure on WRONG repository!${NC}"
    echo ""
    echo "CATASTROPHIC FAILURES:"
    for failure in "${CATASTROPHIC_FAILURES[@]}"; do
        echo "  🔴 $failure"
    done
    echo ""
    echo -e "${RED}IMMEDIATE ACTION REQUIRED:${NC}"
    echo "1. STOP ALL WORK IMMEDIATELY"
    echo "2. TRANSITION TO ERROR_RECOVERY"
    echo "3. DELETE ALL INCORRECT INFRASTRUCTURE"
    echo "4. RECREATE ON TARGET REPOSITORY: $TARGET_REPO"
    echo ""
    echo "THIS IS A -100% AUTOMATIC FAILURE!"
    echo "CONTINUE-SOFTWARE-FACTORY=FALSE"
    exit 911
fi