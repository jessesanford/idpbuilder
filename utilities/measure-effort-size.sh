#!/bin/bash
# measure-effort-size.sh - Correctly measure effort size in separate git repositories
# 
# CRITICAL: This script MUST be run from within an effort directory
# Efforts are SEPARATE git repositories with their own branches!

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}=== EFFORT SIZE MEASUREMENT TOOL ===${NC}"
echo ""

# Step 1: Verify we're in a git repository
if [ ! -d ".git" ]; then
    echo -e "${RED}ERROR: Not in a git repository!${NC}"
    echo "This script must be run from within an effort directory."
    echo "Efforts are SEPARATE git repositories, not just directories."
    echo ""
    echo "Example:"
    echo "  cd efforts/phase2/wave1/go-containerregistry-image-builder"
    echo "  ../../utilities/measure-effort-size.sh"
    exit 1
fi

# Step 2: Check for uncommitted changes
if [[ -n $(git status --porcelain) ]]; then
    echo -e "${YELLOW}WARNING: Uncommitted changes detected!${NC}"
    echo "The line counter uses 'git diff' which requires committed code."
    echo ""
    echo "To commit your changes:"
    echo "  git add -A"
    echo "  git commit -m 'feat: implementation ready for measurement'"
    echo "  git push"
    echo ""
    read -p "Continue anyway? (y/n) " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        exit 1
    fi
fi

# Step 3: Get current branch name (NOT directory name!)
CURRENT_BRANCH=$(git branch --show-current)
if [ -z "$CURRENT_BRANCH" ]; then
    echo -e "${RED}ERROR: Could not determine current branch!${NC}"
    echo "Make sure you're on a branch, not in detached HEAD state."
    exit 1
fi
echo -e "Current branch: ${GREEN}$CURRENT_BRANCH${NC}"

# Step 4: Find base/integration branch
echo ""
echo "Available branches in this repository:"
git branch -a | grep -E "integration|main" || true
echo ""

# Try to auto-detect base branch
BASE_BRANCH=""
if git branch -a | grep -q "phase.*/integration"; then
    BASE_BRANCH=$(git branch -a | grep -E "phase.*/integration" | head -1 | sed 's/.*\///')
elif git branch -a | grep -q "origin/main"; then
    BASE_BRANCH="main"
elif git branch -a | grep -q "main"; then
    BASE_BRANCH="main"
fi

if [ -n "$BASE_BRANCH" ]; then
    echo -e "Auto-detected base branch: ${GREEN}$BASE_BRANCH${NC}"
    read -p "Use this base branch? (y/n) " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        read -p "Enter base branch name: " BASE_BRANCH
    fi
else
    read -p "Enter base branch name (e.g., phase2/integration): " BASE_BRANCH
fi

# Step 5: Find the line-counter.sh tool
TOOL=""
for path in "../../tools/line-counter.sh" "../../../tools/line-counter.sh" "../../../../tools/line-counter.sh" "/home/vscode/software-factory-template/tools/line-counter.sh"; do
    if [ -f "$path" ]; then
        TOOL="$path"
        break
    fi
done

if [ -z "$TOOL" ]; then
    echo -e "${RED}ERROR: Cannot find line-counter.sh tool!${NC}"
    echo "Expected locations:"
    echo "  ../../tools/line-counter.sh"
    echo "  ../../../tools/line-counter.sh"
    exit 1
fi

echo ""
echo -e "${YELLOW}Running measurement...${NC}"
echo "Command: $TOOL -b \"$BASE_BRANCH\" -c \"$CURRENT_BRANCH\""
echo ""

# Step 6: Run the measurement
if $TOOL -b "$BASE_BRANCH" -c "$CURRENT_BRANCH"; then
    echo ""
    echo -e "${GREEN}Measurement complete!${NC}"
    
    # Extract the line count
    SIZE=$($TOOL -b "$BASE_BRANCH" -c "$CURRENT_BRANCH" 2>/dev/null | grep "Total" | awk '{print $NF}')
    if [ -n "$SIZE" ]; then
        echo ""
        if [ "$SIZE" -le 800 ]; then
            echo -e "${GREEN}✅ Size compliance: $SIZE lines (within 800 line limit)${NC}"
        else
            echo -e "${RED}❌ SIZE VIOLATION: $SIZE lines (exceeds 800 line limit!)${NC}"
            echo "IMMEDIATE ACTION REQUIRED: Split this effort!"
        fi
    fi
else
    echo -e "${RED}ERROR: Line counter failed!${NC}"
    echo "Common issues:"
    echo "1. Base branch doesn't exist in this repository"
    echo "2. Current branch not pushed to remote"
    echo "3. No actual changes between branches"
    exit 1
fi

echo ""
echo "=== MEASUREMENT NOTES ==="
echo "1. This tool measures ONLY committed changes"
echo "2. It excludes generated code (*.pb.go, *_generated.go, etc.)"
echo "3. It excludes vendor directories and test files"
echo "4. The base branch should be the integration branch for this phase"
echo "5. NEVER use directory names as branch names!"