#!/bin/bash

# 🔧 EFFORT WORKSPACE SETUP UTILITY
# Part of Software Factory 2.0 - RULE R271 (Single-Branch Full Checkout)
#
# Purpose: Set up complete git workspace with FULL code for an effort
# Usage: ./setup-effort-workspace.sh <phase> <wave> <effort-id> <effort-description>
#
# Example: ./setup-effort-workspace.sh 1 1 "1.1.1" "write-command-tests"
#
# This script MUST be run by the orchestrator BEFORE spawning agents
# R271 SUPREME LAW: Full single-branch checkouts only (NO SPARSE)

set -e

echo "═══════════════════════════════════════════════════════════════"
echo "🔧 EFFORT WORKSPACE SETUP"
echo "═══════════════════════════════════════════════════════════════"
echo "Timestamp: $(date '+%Y-%m-%d %H:%M:%S %Z')"
echo "═══════════════════════════════════════════════════════════════"

# Get the script's directory to find target-repo-config.yaml
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"

# Parse arguments
PHASE="${1}"
WAVE="${2}"
EFFORT_ID="${3}"
EFFORT_DESCRIPTION="${4}"

# Read repository configuration from target-repo-config.yaml
CONFIG_FILE="${PROJECT_DIR}/target-repo-config.yaml"
if [ ! -f "$CONFIG_FILE" ]; then
    echo "❌ ERROR: target-repo-config.yaml not found at: $CONFIG_FILE"
    echo "   Please ensure target-repo-config.yaml exists in the project root"
    exit 1
fi

# Extract repository URL and base branch from config
REPO_URL=$(grep "url:" "$CONFIG_FILE" | head -1 | sed 's/.*url: *//' | tr -d '"' | tr -d "'")
BASE_BRANCH=$(grep "base_branch:" "$CONFIG_FILE" | head -1 | sed 's/.*base_branch: *//' | tr -d '"' | tr -d "'")

# Set defaults if not found in config
if [ -z "$REPO_URL" ] || [ "$REPO_URL" = "https://github.com/OWNER/REPO.git" ]; then
    echo "❌ ERROR: Repository URL not configured in target-repo-config.yaml"
    echo "   Please update target-repo-config.yaml with your target repository URL"
    exit 1
fi

if [ -z "$BASE_BRANCH" ]; then
    BASE_BRANCH="main"
    echo "⚠️ WARNING: Base branch not specified in config, using default: main"
fi

# Validate arguments
if [ -z "$PHASE" ] || [ -z "$WAVE" ] || [ -z "$EFFORT_ID" ] || [ -z "$EFFORT_DESCRIPTION" ]; then
    echo "❌ ERROR: Missing required arguments"
    echo "Usage: $0 <phase> <wave> <effort-id> <effort-description>"
    echo "Example: $0 1 1 \"1.1.1\" \"write-command-tests\""
    exit 1
fi

# Construct paths using both ID and description
EFFORT_NAME="${EFFORT_ID}-${EFFORT_DESCRIPTION}"
EFFORT_DIR="efforts/phase${PHASE}/wave${WAVE}/${EFFORT_NAME}"
BRANCH_NAME="phase${PHASE}/wave${WAVE}/effort-${EFFORT_NAME}"

echo "Configuration:"
echo "  Phase: ${PHASE}"
echo "  Wave: ${WAVE}"
echo "  Effort ID: ${EFFORT_ID}"
echo "  Effort Description: ${EFFORT_DESCRIPTION}"
echo "  Directory: ${EFFORT_DIR}"
echo "  Branch: ${BRANCH_NAME}"
echo "  Repository: ${REPO_URL}"
echo "  Base Branch: ${BASE_BRANCH}"
echo ""

# Function to determine appropriate base branch (R308 Incremental Branching)
determine_base_branch() {
    local phase="$1"
    local wave="$2"
    local effort="$3"
    local default_base="$4"
    
    echo "🧠 THINKING: What base branch should $effort use? (R308)"
    
    # R308: Incremental Branching Strategy
    # Phase 1, Wave 1: Use main/default
    if [[ $phase -eq 1 && $wave -eq 1 ]]; then
        echo "   Decision: Phase 1, Wave 1 - using main branch"
        echo "$default_base"
        return
    fi
    
    # First wave of new phase: Use previous phase integration
    if [[ $wave -eq 1 ]]; then
        local prev_phase=$((phase - 1))
        local integration_branch="phase${prev_phase}-integration"
        
        # Check if integration branch exists
        if git ls-remote --heads origin "$integration_branch" > /dev/null 2>&1; then
            echo "   Decision: Wave 1 of Phase $phase - using previous phase integration"
            echo "$integration_branch"
            return
        else
            echo "   ⚠️ Previous phase integration not found, checking for last wave..."
            # Fallback: try to find last wave of previous phase
            for w in 5 4 3 2 1; do
                local alt_branch="phase${prev_phase}-wave${w}-integration"
                if git ls-remote --heads origin "$alt_branch" > /dev/null 2>&1; then
                    echo "   Decision: Using last available wave integration from Phase $prev_phase"
                    echo "$alt_branch"
                    return
                fi
            done
        fi
    fi
    
    # Subsequent waves: Use previous wave integration
    local prev_wave=$((wave - 1))
    local wave_integration="phase${phase}-wave${prev_wave}-integration"
    
    # Check if previous wave integration exists
    if git ls-remote --heads origin "$wave_integration" > /dev/null 2>&1; then
        echo "   Decision: Wave $wave - using previous wave integration (R308)"
        echo "$wave_integration"
        return
    else
        echo "   ⚠️ WARNING: Previous wave integration not found: $wave_integration"
        echo "   ⚠️ This violates R308 - Incremental Branching!"
        echo "   Using default base as fallback (not recommended)"
        echo "$default_base"
        return
    fi
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

# Check 4: No sparse checkout (R271 compliance)
if [ -f ".git/info/sparse-checkout" ]; then
    echo "❌ Sparse checkout detected - R271 VIOLATION!"
    CHECKS_PASSED=false
else
    echo "✅ No sparse checkout - full clone verified (R271 compliant)"
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