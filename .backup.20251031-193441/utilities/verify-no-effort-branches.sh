#!/bin/bash
# verify-no-effort-branches.sh
# R309 Enforcement: Detect and report effort branches in SF repo

set -euo pipefail

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo "═══════════════════════════════════════════════════════════════"
echo "🔍 R309 ENFORCEMENT: Checking for effort branches in SF repo"
echo "═══════════════════════════════════════════════════════════════"

# Ensure we're in the SF root
SF_ROOT="${CLAUDE_PROJECT_DIR:-$(pwd)}"
cd "$SF_ROOT"

# Verify this IS the SF repo (should have these markers)
if [ ! -f ".claude/CLAUDE.md" ] && [ ! -f "rule-library/RULE-REGISTRY.md" ]; then
    echo -e "${YELLOW}⚠️ WARNING: This doesn't appear to be the SF repo root${NC}"
    echo "Expected to find .claude/CLAUDE.md or rule-library/"
    exit 1
fi

echo "📍 Checking repository: $(pwd)"
echo "🔗 Remote: $(git remote get-url origin 2>/dev/null || echo 'no remote')"
echo ""

# Check for effort-related branches
echo "🌿 Scanning for effort/wave/split branches..."
EFFORT_BRANCHES=$(git branch -a 2>/dev/null | grep -E "(effort|wave|split|phase[0-9])" | grep -v "main\|master" || true)

if [ -n "$EFFORT_BRANCHES" ]; then
    echo -e "${RED}🔴🔴🔴 R309 VIOLATION DETECTED! 🔴🔴🔴${NC}"
    echo -e "${RED}═══════════════════════════════════════════════════════${NC}"
    echo -e "${RED}CRITICAL: Found effort branches in Software Factory repo!${NC}"
    echo -e "${RED}This is the PLANNING repo - efforts go in TARGET clones!${NC}"
    echo ""
    echo "Polluting branches found:"
    echo "$EFFORT_BRANCHES" | while read branch; do
        echo -e "  ${RED}❌ $branch${NC}"
    done
    echo ""
    echo -e "${RED}THESE MUST BE DELETED IMMEDIATELY!${NC}"
    echo ""
    echo "To clean up this pollution:"
    echo "1. Delete local branches:"
    echo "$EFFORT_BRANCHES" | grep -v "remotes/" | while read branch; do
        branch_clean=$(echo "$branch" | sed 's/^[* ]*//')
        echo "   git branch -D '$branch_clean'"
    done
    echo ""
    echo "2. Delete remote branches:"
    echo "$EFFORT_BRANCHES" | grep "remotes/origin/" | while read branch; do
        branch_clean=$(echo "$branch" | sed 's|remotes/origin/||')
        echo "   git push origin --delete '$branch_clean'"
    done
    echo ""
    echo "3. Then clone TARGET repo to efforts/ for actual work"
    echo ""
    echo -e "${RED}═══════════════════════════════════════════════════════${NC}"
    echo -e "${RED}GRADING IMPACT: -100% AUTOMATIC FAILURE${NC}"
    exit 309
else
    echo -e "${GREEN}✅ PASS: No effort branches found in SF repo${NC}"
fi

# Check for code files that shouldn't be in SF root
echo ""
echo "📁 Scanning for misplaced code files..."
MISPLACED_CODE=$(find . -path "./efforts" -prune -o \
    -path "./.git" -prune -o \
    -path "./node_modules" -prune -o \
    -path "./venv" -prune -o \
    \( -name "*.go" -o -name "*.py" -o -name "*.js" -o -name "*.ts" -o -name "*.java" \) \
    -type f -print 2>/dev/null | head -20 || true)

if [ -n "$MISPLACED_CODE" ]; then
    echo -e "${YELLOW}⚠️ WARNING: Found code files in SF repo:${NC}"
    echo "$MISPLACED_CODE"
    echo ""
    echo "These might indicate effort work happening in wrong location."
    echo "Implementation code should be in efforts/phaseX/waveY/effort-name/"
else
    echo -e "${GREEN}✅ PASS: No misplaced code files in SF root${NC}"
fi

# Check target-repo-config.yaml
echo ""
echo "📋 Checking target repository configuration..."
if [ -f "target-repo-config.yaml" ]; then
    TARGET_URL=$(yq '.target_repository.url' target-repo-config.yaml 2>/dev/null || echo "")
    if [[ "$TARGET_URL" == *"software-factory"* ]]; then
        echo -e "${RED}🔴 ERROR: Target repo points to Software Factory!${NC}"
        echo "Target URL: $TARGET_URL"
        echo "This will cause recursive cloning!"
        echo "Fix target-repo-config.yaml to point to actual project repo"
        exit 309
    else
        echo -e "${GREEN}✅ PASS: Target repo is not SF template${NC}"
        echo "Target: $TARGET_URL"
    fi
else
    echo -e "${YELLOW}⚠️ WARNING: No target-repo-config.yaml found${NC}"
fi

# Summary
echo ""
echo "═══════════════════════════════════════════════════════════════"
echo "📊 R309 VERIFICATION COMPLETE"
echo "═══════════════════════════════════════════════════════════════"
echo ""
echo "Remember the fundamental separation:"
echo "  • SF Repo = Planning/Orchestration (here)"
echo "  • Target Repo = Implementation (efforts/)"
echo "  • NEVER mix them up!"
echo ""
echo -e "${GREEN}✅ Verification passed - SF repo is clean${NC}"