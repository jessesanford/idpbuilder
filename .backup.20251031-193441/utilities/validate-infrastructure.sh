#!/bin/bash

# 🚨🚨🚨 INFRASTRUCTURE VALIDATION SCRIPT - R507/R508 ENFORCEMENT
# Purpose: Validate ALL infrastructure against authoritative sources
# Exit codes:
#   0 = All validation passed
#   1 = Validation failure (recoverable)
#   911 = CATASTROPHIC repository mismatch (R508 violation)

set -euo pipefail

# Set CLAUDE_PROJECT_DIR if not already set
: "${CLAUDE_PROJECT_DIR:=/home/vscode/software-factory-template}"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo "========================================="
echo "🔍 INFRASTRUCTURE VALIDATION STARTING"
echo "========================================="
echo "Time: $(date '+%Y-%m-%d %H:%M:%S %Z')"
echo ""

# Track validation results
VALIDATION_PASSED=true
ERRORS=()
CATASTROPHIC=false

# 1. LOAD CONFIGURATION FILES
echo "📋 Loading configuration files..."

# Check for target-repo-config.yaml
if [ ! -f "$CLAUDE_PROJECT_DIR/target-repo-config.yaml" ]; then
    echo -e "${RED}❌ CRITICAL: target-repo-config.yaml not found!${NC}"
    echo "Cannot validate without target repository configuration"
    exit 1
fi

# Check for orchestrator-state-v3.json
if [ ! -f "$CLAUDE_PROJECT_DIR/orchestrator-state-v3.json" ]; then
    echo -e "${YELLOW}⚠️ WARNING: orchestrator-state-v3.json not found${NC}"
    echo "Cannot validate branch names without state file"
fi

# Extract target repository from config (supports both Python yq and Go yq)
# NOTE: The field is named 'url' not 'repository_url' in target-repo-config.yaml
if command -v yq &> /dev/null; then
    # Try Python yq syntax first
    TARGET_REPO=$(cat "$CLAUDE_PROJECT_DIR/target-repo-config.yaml" | yq -r '.url' 2>/dev/null || echo "")

    # If that didn't work, try Go yq syntax
    if [ -z "$TARGET_REPO" ]; then
        TARGET_REPO=$(yq eval '.url' "$CLAUDE_PROJECT_DIR/target-repo-config.yaml" 2>/dev/null || echo "")
    fi
else
    # Fallback to grep if yq is not available
    TARGET_REPO=$(grep "url:" "$CLAUDE_PROJECT_DIR/target-repo-config.yaml" | sed 's/.*url:[ ]*//' | tr -d '"' | tr -d "'" || echo "")
fi

if [ -z "$TARGET_REPO" ]; then
    echo -e "${RED}❌ CRITICAL: No url in target-repo-config.yaml${NC}"
    exit 1
fi

echo "✅ Target repository: $TARGET_REPO"
echo ""

# 2. VALIDATE EFFORT DIRECTORIES
echo "📂 Validating effort directories..."

# Find all effort directories
EFFORT_DIRS=$(find "$CLAUDE_PROJECT_DIR/efforts" -type d -name ".git" 2>/dev/null | xargs -I {} dirname {} || true)

if [ -z "$EFFORT_DIRS" ]; then
    echo "No effort directories found (may be OK if pre-infrastructure)"
else
    while IFS= read -r effort_dir; do
        echo ""
        echo "Checking: $effort_dir"

        # A. VALIDATE REMOTE REPOSITORY (R508 - SUPREME LAW)
        cd "$effort_dir"

        # Check origin remote
        ORIGIN_REMOTE=$(git remote get-url origin 2>/dev/null || echo "")
        TARGET_REMOTE=$(git remote get-url target 2>/dev/null || echo "")

        # Determine which remote is configured
        if [ -n "$ORIGIN_REMOTE" ]; then
            ACTUAL_REMOTE="$ORIGIN_REMOTE"
            REMOTE_NAME="origin"
        elif [ -n "$TARGET_REMOTE" ]; then
            ACTUAL_REMOTE="$TARGET_REMOTE"
            REMOTE_NAME="target"
        else
            echo -e "${RED}❌ No remote configured!${NC}"
            ERRORS+=("$effort_dir: No remote configured")
            VALIDATION_PASSED=false
            continue
        fi

        # SUPREME LAW CHECK - R508
        if [ "$ACTUAL_REMOTE" != "$TARGET_REPO" ]; then
            echo -e "${RED}🔴🔴🔴 CATASTROPHIC FAILURE: WRONG REPOSITORY!${NC}"
            echo -e "${RED}  Expected: $TARGET_REPO${NC}"
            echo -e "${RED}  Actual:   $ACTUAL_REMOTE${NC}"
            echo -e "${RED}  THIS IS A SUPREME LAW VIOLATION (R508)${NC}"
            CATASTROPHIC=true
            VALIDATION_PASSED=false
            ERRORS+=("$effort_dir: CATASTROPHIC - Wrong repository!")
        else
            echo -e "${GREEN}✅ Remote repository correct${NC}"
        fi

        # B. VALIDATE BRANCH NAME
        CURRENT_BRANCH=$(git branch --show-current 2>/dev/null || echo "")

        if [ -z "$CURRENT_BRANCH" ]; then
            echo -e "${RED}❌ Not on any branch (detached HEAD?)${NC}"
            ERRORS+=("$effort_dir: Not on any branch")
            VALIDATION_PASSED=false
        else
            # Try to find expected branch from orchestrator-state-v3.json and pre_planned_infrastructure
            if [ -f "$CLAUDE_PROJECT_DIR/orchestrator-state-v3.json" ]; then
                # Extract effort key from path (e.g., efforts/phase1/wave1/effort1)
                EFFORT_KEY=$(echo "$effort_dir" | sed "s|$CLAUDE_PROJECT_DIR/efforts/||" | sed 's|/|_|g')

                # Try to find pre-planned branch name - check both pre_planned_infrastructure and efforts_in_progress
                PLANNED_BRANCH=$(jq -r ".pre_planned_infrastructure.efforts.\"$EFFORT_KEY\".branch_name // .efforts_in_progress.\"$EFFORT_KEY\".branch // empty" \
                    "$CLAUDE_PROJECT_DIR/orchestrator-state-v3.json" 2>/dev/null || echo "")

                # Also check for the expected remote and tracking info
                PLANNED_REMOTE=$(jq -r ".pre_planned_infrastructure.efforts.\"$EFFORT_KEY\".target_remote // empty" \
                    "$CLAUDE_PROJECT_DIR/orchestrator-state-v3.json" 2>/dev/null || echo "")

                PLANNED_REMOTE_BRANCH=$(jq -r ".pre_planned_infrastructure.efforts.\"$EFFORT_KEY\".remote_branch // empty" \
                    "$CLAUDE_PROJECT_DIR/orchestrator-state-v3.json" 2>/dev/null || echo "")

                if [ -n "$PLANNED_BRANCH" ]; then
                    if [ "$CURRENT_BRANCH" != "$PLANNED_BRANCH" ]; then
                        echo -e "${RED}❌ Wrong branch name!${NC}"
                        echo -e "${RED}  Expected: $PLANNED_BRANCH${NC}"
                        echo -e "${RED}  Actual:   $CURRENT_BRANCH${NC}"
                        ERRORS+=("$effort_dir: Wrong branch name")
                        VALIDATION_PASSED=false
                    else
                        echo -e "${GREEN}✅ Branch name correct: $CURRENT_BRANCH${NC}"

                        # Now validate remote tracking configuration
                        if [ -n "$PLANNED_REMOTE" ] && [ -n "$PLANNED_REMOTE_BRANCH" ]; then
                            # Get the actual tracking info
                            TRACKING_INFO=$(git branch -vv | grep "^\\*" | sed 's/.*\[\(.*\)\].*/\1/' | cut -d: -f1 2>/dev/null || echo "")

                            if [ -z "$TRACKING_INFO" ]; then
                                echo -e "${RED}❌ Branch is not tracking any remote!${NC}"
                                echo -e "${RED}  Expected to track: $PLANNED_REMOTE_BRANCH${NC}"
                                echo -e "${YELLOW}  This likely means the -u flag was missing in git push${NC}"
                                echo -e "${YELLOW}  Correct command: git push -u origin \"\$BRANCH_NAME\"${NC}"
                                echo -e "${YELLOW}  Recovery: git branch --set-upstream-to=$PLANNED_REMOTE_BRANCH${NC}"
                                ERRORS+=("$effort_dir: Branch not tracking remote (missing -u flag?)")
                                VALIDATION_PASSED=false
                            else
                                # Check if tracking the correct remote branch
                                if [ "$TRACKING_INFO" != "$PLANNED_REMOTE_BRANCH" ]; then
                                    echo -e "${RED}❌ Branch tracking wrong remote!${NC}"
                                    echo -e "${RED}  Expected: $PLANNED_REMOTE_BRANCH${NC}"
                                    echo -e "${RED}  Actual:   $TRACKING_INFO${NC}"
                                    ERRORS+=("$effort_dir: Wrong tracking remote")
                                    VALIDATION_PASSED=false
                                else
                                    echo -e "${GREEN}✅ Branch tracking correct remote: $TRACKING_INFO${NC}"
                                fi
                            fi
                        fi
                    fi
                else
                    echo -e "${YELLOW}⚠️ Cannot verify branch name (not in pre_planned_infrastructure)${NC}"
                    echo "  Current branch: $CURRENT_BRANCH"
                fi
            else
                echo "  Current branch: $CURRENT_BRANCH (cannot validate)"
            fi
        fi

        # C. VALIDATE DIRECTORY PATH
        # Check if directory matches expected path pattern
        if [[ "$effort_dir" =~ ^.*/efforts/phase[0-9]+/(wave[0-9]+/)?effort[0-9]+(-split[0-9]+)?$ ]]; then
            echo -e "${GREEN}✅ Directory path follows correct pattern${NC}"
        else
            echo -e "${YELLOW}⚠️ Directory path may not follow standard pattern${NC}"
            echo "  Path: $effort_dir"
        fi

    done <<< "$EFFORT_DIRS"
fi

# 3. VALIDATE SPLIT DIRECTORIES
echo ""
echo "📂 Validating split directories..."

# Find all split directories (within efforts)
SPLIT_DIRS=$(find "$CLAUDE_PROJECT_DIR/efforts" -type d -name "*-split*" -path "*/.git" 2>/dev/null | xargs -I {} dirname {} || true)

if [ -z "$SPLIT_DIRS" ]; then
    echo "No split directories found (may be OK if no splits needed)"
else
    while IFS= read -r split_dir; do
        echo ""
        echo "Checking split: $split_dir"

        # A. VALIDATE REMOTE REPOSITORY (R508 - SUPREME LAW)
        cd "$split_dir"

        # Check origin remote
        ORIGIN_REMOTE=$(git remote get-url origin 2>/dev/null || echo "")
        TARGET_REMOTE=$(git remote get-url target 2>/dev/null || echo "")

        # Determine which remote is configured
        if [ -n "$ORIGIN_REMOTE" ]; then
            ACTUAL_REMOTE="$ORIGIN_REMOTE"
            REMOTE_NAME="origin"
        elif [ -n "$TARGET_REMOTE" ]; then
            ACTUAL_REMOTE="$TARGET_REMOTE"
            REMOTE_NAME="target"
        else
            echo -e "${RED}❌ No remote configured for split!${NC}"
            ERRORS+=("$split_dir: No remote configured")
            VALIDATION_PASSED=false
            continue
        fi

        # SUPREME LAW CHECK - R508
        if [ "$ACTUAL_REMOTE" != "$TARGET_REPO" ]; then
            echo -e "${RED}🔴🔴🔴 CATASTROPHIC FAILURE: SPLIT ON WRONG REPOSITORY!${NC}"
            echo -e "${RED}  Expected: $TARGET_REPO${NC}"
            echo -e "${RED}  Actual:   $ACTUAL_REMOTE${NC}"
            echo -e "${RED}  THIS IS A SUPREME LAW VIOLATION (R508)${NC}"
            CATASTROPHIC=true
            VALIDATION_PASSED=false
            ERRORS+=("$split_dir: CATASTROPHIC - Split on wrong repository!")
        else
            echo -e "${GREEN}✅ Split remote repository correct${NC}"
        fi

        # B. VALIDATE BRANCH NAME
        CURRENT_BRANCH=$(git branch --show-current 2>/dev/null || echo "")

        if [ -z "$CURRENT_BRANCH" ]; then
            echo -e "${RED}❌ Split not on any branch (detached HEAD?)${NC}"
            ERRORS+=("$split_dir: Not on any branch")
            VALIDATION_PASSED=false
        else
            # Validate split branch naming convention
            if [[ "$CURRENT_BRANCH" =~ -split-[0-9]+$ ]]; then
                echo -e "${GREEN}✅ Split branch name follows convention: $CURRENT_BRANCH${NC}"
            else
                echo -e "${YELLOW}⚠️ Split branch name may not follow convention: $CURRENT_BRANCH${NC}"
                echo "  Expected format: *-split-N"
            fi
        fi

    done <<< "$SPLIT_DIRS"
fi

# 4. VALIDATE PRE-PLANNED INTEGRATE_WAVE_EFFORTS INFRASTRUCTURE (R504)
echo ""
echo "📂 Validating pre-planned integration infrastructure..."

# Check if orchestrator-state-v3.json has integration pre-planning
if [ -f "$CLAUDE_PROJECT_DIR/orchestrator-state-v3.json" ]; then
    HAS_INTEGRATE_WAVE_EFFORTSS=$(jq -r 'has("pre_planned_infrastructure") and .pre_planned_infrastructure | has("integrations")' "$CLAUDE_PROJECT_DIR/orchestrator-state-v3.json")

    if [ "$HAS_INTEGRATE_WAVE_EFFORTSS" = "true" ]; then
        echo "✅ Integration infrastructure pre-planning found"

        # Validate all pre-planned integration URLs match target repo
        WRONG_URLS=$(jq -r '
            [.pre_planned_infrastructure.integrations.wave_integrations[]?.target_repo_url // empty,
             .pre_planned_infrastructure.integrations.phase_integrations[]?.target_repo_url // empty]
            | map(select(. != null and . != "'"$TARGET_REPO"'"))
            | unique[]' "$CLAUDE_PROJECT_DIR/orchestrator-state-v3.json" 2>/dev/null)

        if [ -n "$WRONG_URLS" ]; then
            echo -e "${RED}🔴🔴🔴 CATASTROPHIC: Pre-planned integrations have wrong repository URLs!${NC}"
            echo -e "${RED}Expected: $TARGET_REPO${NC}"
            echo -e "${RED}Found wrong URLs: $WRONG_URLS${NC}"
            CATASTROPHIC=true
            VALIDATION_PASSED=false
        else
            echo "✅ All pre-planned integration URLs correct"
        fi
    else
        echo -e "${YELLOW}⚠️ WARNING: No integration pre-planning found in orchestrator-state-v3.json${NC}"
        echo "This will cause failures in SETUP_INTEGRATE_WAVE_EFFORTS_INFRASTRUCTURE state (R504)"
    fi
fi

# 5. VALIDATE ACTUAL INTEGRATE_WAVE_EFFORTS DIRECTORIES
echo ""
echo "📂 Validating existing integration directories..."

# Find all integration directories
INTEGRATE_WAVE_EFFORTS_DIRS=$(find "$CLAUDE_PROJECT_DIR" -type d -name "*-integration" -o -name "*-merge" | grep -v ".git" || true)

if [ -z "$INTEGRATE_WAVE_EFFORTS_DIRS" ]; then
    echo "No integration directories found (may be OK if no integrations yet)"
else
    while IFS= read -r integration_dir; do
        # Skip if not a git directory
        if [ ! -d "$integration_dir/.git" ]; then
            continue
        fi

        echo ""
        echo "Checking integration: $integration_dir"

        # A. VALIDATE REMOTE REPOSITORY (R508 - SUPREME LAW)
        cd "$integration_dir"

        # Check origin remote
        ORIGIN_REMOTE=$(git remote get-url origin 2>/dev/null || echo "")
        TARGET_REMOTE=$(git remote get-url target 2>/dev/null || echo "")

        # Determine which remote is configured
        if [ -n "$ORIGIN_REMOTE" ]; then
            ACTUAL_REMOTE="$ORIGIN_REMOTE"
            REMOTE_NAME="origin"
        elif [ -n "$TARGET_REMOTE" ]; then
            ACTUAL_REMOTE="$TARGET_REMOTE"
            REMOTE_NAME="target"
        else
            echo -e "${RED}❌ No remote configured for integration!${NC}"
            ERRORS+=("$integration_dir: No remote configured")
            VALIDATION_PASSED=false
            continue
        fi

        # SUPREME LAW CHECK - R508
        if [ "$ACTUAL_REMOTE" != "$TARGET_REPO" ]; then
            echo -e "${RED}🔴🔴🔴 CATASTROPHIC FAILURE: INTEGRATE_WAVE_EFFORTS ON WRONG REPOSITORY!${NC}"
            echo -e "${RED}  Expected: $TARGET_REPO${NC}"
            echo -e "${RED}  Actual:   $ACTUAL_REMOTE${NC}"
            echo -e "${RED}  THIS IS A SUPREME LAW VIOLATION (R508)${NC}"
            CATASTROPHIC=true
            VALIDATION_PASSED=false
            ERRORS+=("$integration_dir: CATASTROPHIC - Integration on wrong repository!")
        else
            echo -e "${GREEN}✅ Integration remote repository correct${NC}"
        fi

        # B. VALIDATE BRANCH NAME
        CURRENT_BRANCH=$(git branch --show-current 2>/dev/null || echo "")

        if [ -z "$CURRENT_BRANCH" ]; then
            echo -e "${RED}❌ Integration not on any branch (detached HEAD?)${NC}"
            ERRORS+=("$integration_dir: Not on any branch")
            VALIDATION_PASSED=false
        else
            # Validate integration branch naming convention
            if [[ "$CURRENT_BRANCH" =~ -(integration|merge)(-|$) ]]; then
                echo -e "${GREEN}✅ Integration branch name follows convention: $CURRENT_BRANCH${NC}"
            else
                echo -e "${YELLOW}⚠️ Integration branch name may not follow convention: $CURRENT_BRANCH${NC}"
                echo "  Expected to contain: -integration or -merge"
            fi
        fi

    done <<< "$INTEGRATE_WAVE_EFFORTS_DIRS"
fi

echo ""
echo "========================================="
echo "📊 VALIDATION SUMMARY"
echo "========================================="

if $CATASTROPHIC; then
    echo -e "${RED}🔴🔴🔴 CATASTROPHIC FAILURE DETECTED!${NC}"
    echo -e "${RED}SUPREME LAW VIOLATION (R508): Wrong repository configured${NC}"
    echo ""
    echo "Errors found:"
    for error in "${ERRORS[@]}"; do
        echo "  - $error"
    done
    echo ""
    echo -e "${RED}IMMEDIATE ACTION REQUIRED: Transition to ERROR_RECOVERY${NC}"
    exit 911  # Special exit code for catastrophic failure
elif $VALIDATION_PASSED; then
    echo -e "${GREEN}✅ ALL VALIDATION CHECKS PASSED${NC}"
    echo "Infrastructure is correctly configured"
    exit 0
else
    echo -e "${RED}❌ VALIDATION FAILED${NC}"
    echo ""
    echo "Errors found:"
    for error in "${ERRORS[@]}"; do
        echo "  - $error"
    done
    echo ""
    echo "Action required: Fix infrastructure issues before proceeding"
    exit 1
fi