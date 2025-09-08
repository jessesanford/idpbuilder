#!/bin/bash
# verify-phase-integration-branch.sh
# Verifies R259 compliance - phase integration branch after ERROR_RECOVERY fixes

set -e

# Source helper functions
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "${SCRIPT_DIR}/state-file-update-functions.sh" 2>/dev/null || true

echo "🔍 Verifying Phase Integration Branch (R259 Compliance)"
echo "========================================================="

# Check if orchestrator-state.json exists
if [ ! -f "orchestrator-state.json" ]; then
    echo "❌ No orchestrator-state.json found"
    exit 1
fi

# Get current phase
CURRENT_PHASE=$(jq '.current_phase' orchestrator-state.json)
CURRENT_STATE=$(jq '.current_state' orchestrator-state.json)
PREVIOUS_STATE=$(jq '.previous_state' orchestrator-state.json)

echo "📊 Current Status:"
echo "  - Phase: ${CURRENT_PHASE}"
echo "  - State: ${CURRENT_STATE}"
echo "  - Previous: ${PREVIOUS_STATE}"
echo ""

# Function to verify phase integration branch
verify_phase_integration_branch() {
    local PHASE=$1
    local BRANCH_NAME="phase${PHASE}-integration-*"
    local POST_FIXES_BRANCH="phase${PHASE}-post-fixes-integration-*"
    
    echo "🔍 Checking for phase integration branches..."
    
    # Check for post-fixes integration branch (after ERROR_RECOVERY)
    if git branch -r | grep -q "$POST_FIXES_BRANCH"; then
        local BRANCH=$(git branch -r | grep "$POST_FIXES_BRANCH" | head -1 | sed 's/origin\///')
        echo "✅ Found post-fixes integration branch: $BRANCH"
        
        # Verify it contains ERROR_RECOVERY fixes
        echo "🔍 Verifying branch contains fixes..."
        if git log "origin/$BRANCH" --oneline | grep -q -E "fix|ERROR_RECOVERY|Merge.*fix"; then
            echo "✅ Branch contains ERROR_RECOVERY fixes"
        else
            echo "⚠️ Warning: Branch may not contain expected fixes"
        fi
        
        # Check if recorded in state file
        local RECORDED=$(jq ".phase_integration_branches[] | select(.phase == $PHASE).branch" orchestrator-state.json)
        if [ -n "$RECORDED" ]; then
            echo "✅ Branch recorded in state file: $RECORDED"
        else
            echo "⚠️ Warning: Branch not recorded in orchestrator-state.json"
        fi
        
        return 0
    fi
    
    # Check for regular phase integration branch
    if git branch -r | grep -q "$BRANCH_NAME"; then
        local BRANCH=$(git branch -r | grep "$BRANCH_NAME" | head -1 | sed 's/origin\///')
        echo "ℹ️ Found regular phase integration branch: $BRANCH"
        return 0
    fi
    
    echo "❌ No phase integration branch found for phase ${PHASE}"
    return 1
}

# Function to check if phase integration is required
check_if_integration_required() {
    local PHASE=$1
    
    # Check if we're coming from ERROR_RECOVERY with phase fixes
    local ERROR_RECOVERY_TYPE=$(jq '.error_recovery_reason' orchestrator-state.json)
    
    if [ "$ERROR_RECOVERY_TYPE" = "PHASE_ASSESSMENT_NEEDS_WORK" ]; then
        echo "⚠️ Phase assessment fixes detected - integration branch REQUIRED (R259)"
        return 0
    fi
    
    # Check if phase assessment report exists with NEEDS_WORK
    local ASSESSMENT_REPORT="phase-assessments/phase${PHASE}/PHASE-${PHASE}-ASSESSMENT-REPORT.md"
    if [ -f "$ASSESSMENT_REPORT" ]; then
        if grep -q "NEEDS_WORK\|CHANGES_REQUIRED" "$ASSESSMENT_REPORT"; then
            echo "⚠️ Phase assessment requires fixes - integration branch needed"
            return 0
        fi
    fi
    
    return 1
}

# Main verification logic
echo "🔍 Performing R259 Compliance Check..."
echo ""

# Check if phase integration is required
if check_if_integration_required $CURRENT_PHASE; then
    echo "📋 Phase integration branch is REQUIRED per R259"
    
    if ! verify_phase_integration_branch $CURRENT_PHASE; then
        echo ""
        echo "❌❌❌ R259 VIOLATION DETECTED ❌❌❌"
        echo "Phase assessment fixes completed but no integration branch exists!"
        echo ""
        echo "Required Action:"
        echo "1. Transition to PHASE_INTEGRATION state"
        echo "2. Create phase${CURRENT_PHASE}-post-fixes-integration-* branch"
        echo "3. Merge all wave integration branches"
        echo "4. Merge all ERROR_RECOVERY fix branches"
        echo "5. Only then transition to SPAWN_ARCHITECT_PHASE_ASSESSMENT"
        echo ""
        echo "Command to fix:"
        echo "  jq '.current_state = \"PHASE_INTEGRATION\"' orchestrator-state.json"
        exit 1
    fi
else
    echo "ℹ️ Phase integration branch not currently required"
    
    # Still check if one exists
    if verify_phase_integration_branch $CURRENT_PHASE; then
        echo "✅ Phase integration branch exists (good practice)"
    fi
fi

# Check state transition compliance
if [ "$CURRENT_STATE" = "SPAWN_ARCHITECT_PHASE_ASSESSMENT" ]; then
    if [ "$PREVIOUS_STATE" = "ERROR_RECOVERY" ]; then
        echo ""
        echo "⚠️ WARNING: Direct transition from ERROR_RECOVERY to reassessment"
        echo "Should have gone through PHASE_INTEGRATION first (R259)"
    elif [ "$PREVIOUS_STATE" = "PHASE_INTEGRATION" ]; then
        echo "✅ Correct transition: ERROR_RECOVERY → PHASE_INTEGRATION → reassessment"
    fi
fi

echo ""
echo "✅ R259 Compliance Check Complete"

# Generate compliance report
cat > /tmp/r259-compliance-report.yaml << EOF
r259_compliance_report:
  timestamp: "$(date -u +%Y-%m-%dT%H:%M:%SZ)"
  phase: $CURRENT_PHASE
  current_state: "$CURRENT_STATE"
  previous_state: "$PREVIOUS_STATE"
  
  phase_integration_branch:
    required: $(check_if_integration_required $CURRENT_PHASE && echo "true" || echo "false")
    exists: $(verify_phase_integration_branch $CURRENT_PHASE 2>/dev/null && echo "true" || echo "false")
    branch_name: "$(git branch -r | grep "phase${CURRENT_PHASE}.*integration" | head -1 | sed 's/.*origin\///' || echo "none")"
    
  compliance:
    r259_compliant: $(verify_phase_integration_branch $CURRENT_PHASE 2>/dev/null && echo "true" || echo "false")
    correct_transitions: $([ "$PREVIOUS_STATE" = "PHASE_INTEGRATION" ] && echo "true" || echo "false")
    
  recommendations:
    - "Always transition through PHASE_INTEGRATION after phase fixes"
    - "Create integration branch before reassessment"
    - "Document all fixes in integration commit messages"
EOF

echo ""
echo "📄 Report saved to: /tmp/r259-compliance-report.yaml"