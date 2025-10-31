#!/bin/bash
# Validates that project integration is required and not skipped (R283)

set -euo pipefail

echo "================================================================"
echo "🔍 R283 PROJECT INTEGRATE_WAVE_EFFORTS VALIDATION"
echo "================================================================"
echo ""
echo "Checking that multi-phase projects MUST perform project integration..."
echo ""

STATE_MACHINE="/home/vscode/software-factory-template/state-machines/software-factory-3.0-state-machine.json"
COMPLETE_PHASE_RULES="/home/vscode/software-factory-template/agent-states/software-factory/orchestrator/COMPLETE_PHASE/rules.md"
VIOLATIONS=0

# Check 1: State Machine - No direct COMPLETE_PHASE → PROJECT_DONE for multi-phase
echo "📋 Check 1: State Machine Transitions"
echo "-------------------------------------"

# Check if COMPLETE_PHASE can transition to PROJECT_DONE
if jq -e '.transition_matrix.orchestrator.COMPLETE_PHASE[] | select(. == "PROJECT_DONE")' "$STATE_MACHINE" > /dev/null 2>&1; then
    echo "⚠️  WARNING: COMPLETE_PHASE has PROJECT_DONE in transition list"
    echo "   This is only acceptable if it's conditional on single-phase projects"
    echo "   Checking state rules for proper conditions..."
else
    echo "✅ State machine: COMPLETE_PHASE does NOT directly transition to PROJECT_DONE"
fi

# Check if COMPLETE_PHASE can transition to PROJECT_INTEGRATE_WAVE_EFFORTS
if jq -e '.transition_matrix.orchestrator.COMPLETE_PHASE[] | select(. == "PROJECT_INTEGRATE_WAVE_EFFORTS")' "$STATE_MACHINE" > /dev/null 2>&1; then
    echo "✅ State machine: COMPLETE_PHASE → PROJECT_INTEGRATE_WAVE_EFFORTS transition exists"
else
    echo "🔴 VIOLATION: COMPLETE_PHASE missing PROJECT_INTEGRATE_WAVE_EFFORTS transition!"
    ((VIOLATIONS++))
fi

echo ""

# Check 2: COMPLETE_PHASE State Rules
echo "📋 Check 2: COMPLETE_PHASE State Rules"
echo "---------------------------------------"

if [ ! -f "$COMPLETE_PHASE_RULES" ]; then
    echo "🔴 VIOLATION: COMPLETE_PHASE state rules file missing!"
    ((VIOLATIONS++))
else
    # Check for correct transition language
    if grep -q "Multi-phase.*last phase.*PROJECT_INTEGRATE_WAVE_EFFORTS" "$COMPLETE_PHASE_RULES"; then
        echo "✅ State rules: Correct multi-phase → PROJECT_INTEGRATE_WAVE_EFFORTS transition"
    else
        echo "🔴 VIOLATION: State rules missing proper PROJECT_INTEGRATE_WAVE_EFFORTS transition!"
        ((VIOLATIONS++))
    fi

    # Check for R283 reference
    if grep -q "R283" "$COMPLETE_PHASE_RULES"; then
        echo "✅ State rules: R283 reference found"
    else
        echo "⚠️  WARNING: State rules don't reference R283"
    fi

    # Check for prohibition of direct PROJECT_DONE for multi-phase
    if grep -qi "PROHIBITED.*multi-phase\|multi-phase.*PROHIBITED" "$COMPLETE_PHASE_RULES"; then
        echo "✅ State rules: Prohibition of direct PROJECT_DONE documented"
    else
        echo "⚠️  WARNING: State rules should explicitly prohibit direct PROJECT_DONE for multi-phase"
    fi
fi

echo ""

# Check 3: PROJECT_INTEGRATE_WAVE_EFFORTS State Exists
echo "📋 Check 3: PROJECT_INTEGRATE_WAVE_EFFORTS State Chain"
echo "--------------------------------------------"

REQUIRED_STATES=(
    "PROJECT_INTEGRATE_WAVE_EFFORTS"
    "SETUP_PROJECT_INTEGRATE_WAVE_EFFORTS_INFRASTRUCTURE"
    "SPAWN_INTEGRATION_AGENT_PROJECT"
    "MONITORING_PROJECT_INTEGRATE_WAVE_EFFORTS"
)

for state in "${REQUIRED_STATES[@]}"; do
    if jq -e ".states[] | select(. == \"$state\")" "$STATE_MACHINE" > /dev/null 2>&1; then
        echo "✅ State exists: $state"
    else
        echo "🔴 VIOLATION: Required state missing: $state"
        ((VIOLATIONS++))
    fi
done

echo ""

# Check 4: Verify Path to PROJECT_DONE
echo "📋 Check 4: Path to PROJECT_DONE"
echo "---------------------------"

# Only PR_PLAN_CREATION should transition to PROJECT_DONE
PROJECT_DONE_SOURCES=$(jq -r '.transition_matrix.orchestrator | to_entries[] | select(.value[] == "PROJECT_DONE") | .key' "$STATE_MACHINE")
echo "States that can transition to PROJECT_DONE:"
echo "$PROJECT_DONE_SOURCES" | while read -r state; do
    if [ "$state" == "PR_PLAN_CREATION" ]; then
        echo "✅ $state → PROJECT_DONE (correct final state)"
    else
        echo "⚠️  $state → PROJECT_DONE (verify this is appropriate)"
    fi
done

echo ""

# Check 5: R283 Rule File Exists
echo "📋 Check 5: R283 Rule Documentation"
echo "-----------------------------------"

R283_FILE="/home/vscode/software-factory-template/rule-library/R283-project-integration-protocol.md"
if [ -f "$R283_FILE" ]; then
    echo "✅ R283 rule file exists"

    # Check for state machine enforcement section
    if grep -q "State Machine Enforcement" "$R283_FILE"; then
        echo "✅ R283 rule includes state machine enforcement section"
    else
        echo "⚠️  WARNING: R283 should include state machine enforcement section"
    fi

    # Check for -100% penalty
    if grep -q "\-100%" "$R283_FILE"; then
        echo "✅ R283 rule documents -100% penalty"
    else
        echo "⚠️  WARNING: R283 should document -100% penalty for violations"
    fi
else
    echo "🔴 VIOLATION: R283 rule file missing!"
    ((VIOLATIONS++))
fi

echo ""
echo "================================================================"
echo "📊 VALIDATION SUMMARY"
echo "================================================================"

if [ $VIOLATIONS -eq 0 ]; then
    echo "✅ PASS: Project integration requirement properly enforced"
    echo ""
    echo "Multi-phase projects are correctly required to perform project-level"
    echo "integration before reaching PROJECT_DONE state."
    exit 0
else
    echo "🔴 FAIL: $VIOLATIONS CRITICAL VIOLATIONS DETECTED"
    echo ""
    echo "Project integration requirement is NOT properly enforced!"
    echo "This would allow multi-phase projects to skip project integration,"
    echo "resulting in incomplete projects and -100% grading penalty."
    echo ""
    echo "Fix required to ensure R283 compliance."
    exit 283
fi
