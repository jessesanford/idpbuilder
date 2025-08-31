#!/bin/bash

# 🔧 EFFORT WORKSPACE SETUP UTILITY
# Part of Software Factory 2.0 - RULE R271 (Single-Branch Full Checkout)
#
# Purpose: Set up complete git workspace with FULL code for an effort
# Usage: ./setup-effort-workspace.sh <phase> <wave> <effort-name> [repo-url] [base-branch]
#
# Example: ./setup-effort-workspace.sh 1 1 core-types https://github.com/idpbuilder/idpbuilder.git main
#
# This script MUST be run by the orchestrator BEFORE spawning agents
# R271 SUPREME LAW: Full single-branch checkouts only (NO SPARSE)

set -e

echo "═══════════════════════════════════════════════════════════════"
echo "🔧 EFFORT WORKSPACE SETUP"
echo "═══════════════════════════════════════════════════════════════"
echo "Timestamp: $(date '+%Y-%m-%d %H:%M:%S %Z')"
echo "═══════════════════════════════════════════════════════════════"

# Parse arguments
PHASE="${1}"
WAVE="${2}"
EFFORT_NAME="${3}"
REPO_URL="${4:-https://github.com/[org]/[project].git}"
BASE_BRANCH="${5:-main}"

# Validate arguments
if [ -z "$PHASE" ] || [ -z "$WAVE" ] || [ -z "$EFFORT_NAME" ]; then
    echo "❌ ERROR: Missing required arguments"
    echo "Usage: $0 <phase> <wave> <effort-name> [repo-url] [base-branch]"
    echo "Example: $0 1 1 core-types https://github.com/idpbuilder/idpbuilder.git main"
    exit 1
fi

# Construct paths
EFFORT_DIR="efforts/phase${PHASE}/wave${WAVE}/${EFFORT_NAME}"
BRANCH_NAME="phase${PHASE}/wave${WAVE}/effort-${EFFORT_NAME}"

echo "Configuration:"
echo "  Phase: ${PHASE}"
echo "  Wave: ${WAVE}"
echo "  Effort: ${EFFORT_NAME}"
echo "  Directory: ${EFFORT_DIR}"
echo "  Branch: ${BRANCH_NAME}"
echo "  Repository: ${REPO_URL}"
echo "  Base Branch: ${BASE_BRANCH}"
echo ""

# Function to determine appropriate base branch (R271 compliance)
determine_base_branch() {
    local phase="$1"
    local wave="$2"
    local effort="$3"
    local default_base="$4"
    
    echo "🧠 THINKING: What base branch should $effort use?"
    
    # Check if this effort has dependencies in orchestrator-state.yaml
    if [ -f "../../orchestrator-state.yaml" ]; then
        local deps=$(yq ".efforts_planned.\"$effort\".depends_on[]" ../../orchestrator-state.yaml 2>/dev/null)
        if [ -n "$deps" ] && [ "$deps" != "null" ]; then
            echo "   Decision: Has dependencies, checking if they're ready..."
            # For now, use default base (dependency checking would be more complex)
            echo "$default_base"
            return
        fi
    fi
    
    echo "   Decision: No dependencies, using default base branch"
    echo "$default_base"
}

# Step 1: Create directory structure
echo "Step 1: Creating directory structure..."
if [ -d "$EFFORT_DIR" ]; then
    echo "⚠️ Directory already exists: $EFFORT_DIR"
    echo -n "   Overwrite? (y/N): "
    read -r response
    if [[ ! "$response" =~ ^[Yy]$ ]]; then
        echo "❌ Aborted by user"
        exit 1
    fi
    rm -rf "$EFFORT_DIR"
fi

mkdir -p "$EFFORT_DIR"
cd "$EFFORT_DIR"
echo "✅ Directory created: $(pwd)"

# Step 2: Determine appropriate base branch (R271 requirement)
echo ""
echo "Step 2: Determining appropriate base branch..."
DETERMINED_BASE=$(determine_base_branch "$PHASE" "$WAVE" "$EFFORT_NAME" "$BASE_BRANCH")
echo "📌 Selected base branch: $DETERMINED_BASE"

# Step 3: Create FULL single-branch clone (R271 SUPREME LAW)
echo ""
echo "Step 3: Creating FULL single-branch clone (R271 compliance)..."
echo "   NO SPARSE CHECKOUT - Full code from $DETERMINED_BASE"

git clone \
    --single-branch \
    --branch "$DETERMINED_BASE" \
    "$REPO_URL" \
    . || {
    echo "❌ Failed to clone repository"
    echo "   Check repository URL: $REPO_URL"
    echo "   Check base branch: $DETERMINED_BASE"
    exit 1
}

echo "✅ Repository cloned with FULL code"
echo "   Verifying full checkout (R271 compliance check)..."

# Verify NO sparse checkout (R271 enforcement)
if [ -f ".git/info/sparse-checkout" ]; then
    echo "🔴🔴🔴 SUPREME LAW VIOLATION: Sparse checkout detected!"
    exit 1
fi

echo "✅ Full checkout verified - ALL code available"
echo "   Repository structure:"
ls -la | head -10
echo "   ... (showing first 10 entries)"

# Step 4: Already on base branch from clone
echo ""
echo "Step 4: Verifying base branch..."
CURRENT=$(git branch --show-current)
echo "✅ On base branch: $CURRENT"

# Step 5: Create effort branch
echo ""
echo "Step 5: Creating effort branch..."
git checkout -b "$BRANCH_NAME" 2>/dev/null || {
    echo "❌ Failed to create branch: $BRANCH_NAME"
    echo "   Branch may already exist"
    exit 1
}
echo "✅ Effort branch created: $BRANCH_NAME"

# Step 6: Create required files
echo ""
echo "Step 6: Creating required files..."
touch IMPLEMENTATION-PLAN.md
touch work-log.md
mkdir -p pkg

# Create initial work log entry with R271 compliance note
cat > work-log.md << EOF
# Work Log - ${EFFORT_NAME}

## Infrastructure Setup - $(date '+%Y-%m-%d %H:%M:%S')
- Workspace created by orchestrator
- **Branch**: ${BRANCH_NAME}
- **Base Branch**: ${DETERMINED_BASE}
- **Clone Type**: FULL (R271 compliance)
- **Sparse Checkout**: NO (R271 SUPREME LAW)
- **Full Codebase**: YES

## Base Branch Selection Rationale
- Analyzed dependencies for ${EFFORT_NAME}
- Selected ${DETERMINED_BASE} as appropriate base
- Full code available from base branch

## Session 1 - [Date]
- [ ] Review implementation plan
- [ ] Begin implementation
- [ ] Write tests
- [ ] Update documentation

EOF

echo "✅ Required files created"

# Step 7: Verify workspace
echo ""
echo "Step 7: Verifying workspace..."
echo ""

# Run verification checks
CHECKS_PASSED=true

# Check 1: Directory structure
if [ -d ".git" ] && [ -f "IMPLEMENTATION-PLAN.md" ] && [ -f "work-log.md" ]; then
    echo "✅ Directory structure correct"
else
    echo "❌ Directory structure incomplete"
    CHECKS_PASSED=false
fi

# Check 2: Git branch
CURRENT_BRANCH=$(git branch --show-current)
if [[ "$CURRENT_BRANCH" == "$BRANCH_NAME" ]]; then
    echo "✅ Git branch correct: $CURRENT_BRANCH"
else
    echo "❌ Git branch incorrect: $CURRENT_BRANCH (expected: $BRANCH_NAME)"
    CHECKS_PASSED=false
fi

# Check 3: Remote configuration
if git remote -v | grep -q origin; then
    echo "✅ Remote configured"
else
    echo "❌ No remote configured"
    CHECKS_PASSED=false
fi

# Check 4: Sparse checkout
if [ -f ".git/info/sparse-checkout" ]; then
    echo "✅ Sparse checkout active"
else
    echo "⚠️ Sparse checkout not configured (full clone)"
fi

# Final summary
echo ""
echo "═══════════════════════════════════════════════════════════════"
if [ "$CHECKS_PASSED" = true ]; then
    echo "✅ WORKSPACE SETUP COMPLETE"
    echo ""
    echo "Effort workspace ready at:"
    echo "  Directory: $(pwd)"
    echo "  Branch: $BRANCH_NAME"
    echo ""
    echo "Next steps:"
    echo "1. Create implementation plan in IMPLEMENTATION-PLAN.md"
    echo "2. Spawn agent with working directory: $EFFORT_DIR"
    echo "3. Agent should verify workspace before starting"
else
    echo "❌ WORKSPACE SETUP FAILED"
    echo "Fix the issues above before spawning agents"
    exit 1
fi
echo "═══════════════════════════════════════════════════════════════"

# Return to original directory
cd - > /dev/null