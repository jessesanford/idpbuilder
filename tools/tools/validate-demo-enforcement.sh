#!/bin/bash
# Validates demo enforcement in state machine (R291 Gate 4)

echo "🔍 Validating Demo Enforcement (R291 Gate 4)"
echo "============================================="

STATE_MACHINE="state-machines/software-factory-3.0-state-machine.json"
ERRORS=0

# Check for demo validation states
echo ""
echo "📋 Checking for required states..."

if ! jq -e '.states[] | select(. == "SPAWN_CODE_REVIEWER_DEMO_VALIDATION")' "$STATE_MACHINE" > /dev/null 2>&1; then
    echo "🔴 MISSING: SPAWN_CODE_REVIEWER_DEMO_VALIDATION state!"
    ((ERRORS++))
else
    echo "✅ SPAWN_CODE_REVIEWER_DEMO_VALIDATION state present"
fi

if ! jq -e '.states[] | select(. == "WAITING_FOR_DEMO_VALIDATION")' "$STATE_MACHINE" > /dev/null 2>&1; then
    echo "🔴 MISSING: WAITING_FOR_DEMO_VALIDATION state!"
    ((ERRORS++))
else
    echo "✅ WAITING_FOR_DEMO_VALIDATION state present"
fi

# Check for illegal bypass transitions
echo ""
echo "📋 Checking for illegal bypass transitions..."

# WAITING_FOR_REVIEW_WAVE_INTEGRATION should NOT transition directly to REVIEW_WAVE_ARCHITECTURE
if jq -e '.agents.orchestrator.states.WAITING_FOR_REVIEW_WAVE_INTEGRATION.valid_transitions[] | select(. == "REVIEW_WAVE_ARCHITECTURE")' "$STATE_MACHINE" > /dev/null 2>&1; then
    echo "🔴 VIOLATION: Direct transition WAITING_FOR_REVIEW_WAVE_INTEGRATION → REVIEW_WAVE_ARCHITECTURE bypasses demo validation!"
    ((ERRORS++))
else
    echo "✅ No direct WAITING_FOR_REVIEW_WAVE_INTEGRATION → REVIEW_WAVE_ARCHITECTURE bypass"
fi

# WAITING_FOR_REVIEW_WAVE_INTEGRATION should NOT transition directly to COMPLETE_PHASE
if jq -e '.agents.orchestrator.states.WAITING_FOR_REVIEW_WAVE_INTEGRATION.valid_transitions[] | select(. == "COMPLETE_PHASE")' "$STATE_MACHINE" > /dev/null 2>&1; then
    echo "🔴 VIOLATION: Direct transition WAITING_FOR_REVIEW_WAVE_INTEGRATION → COMPLETE_PHASE bypasses demo validation!"
    ((ERRORS++))
else
    echo "✅ No direct WAITING_FOR_REVIEW_WAVE_INTEGRATION → COMPLETE_PHASE bypass"
fi

# Check for correct transitions
echo ""
echo "📋 Checking for correct transitions..."

# WAITING_FOR_REVIEW_WAVE_INTEGRATION should transition to SPAWN_CODE_REVIEWER_DEMO_VALIDATION
if jq -e '.agents.orchestrator.states.WAITING_FOR_REVIEW_WAVE_INTEGRATION.valid_transitions[] | select(. == "SPAWN_CODE_REVIEWER_DEMO_VALIDATION")' "$STATE_MACHINE" > /dev/null 2>&1; then
    echo "✅ Correct transition: WAITING_FOR_REVIEW_WAVE_INTEGRATION → SPAWN_CODE_REVIEWER_DEMO_VALIDATION"
else
    echo "🔴 MISSING: Transition WAITING_FOR_REVIEW_WAVE_INTEGRATION → SPAWN_CODE_REVIEWER_DEMO_VALIDATION"
    ((ERRORS++))
fi

# SPAWN_CODE_REVIEWER_DEMO_VALIDATION should transition to WAITING_FOR_DEMO_VALIDATION
if jq -e '.agents.orchestrator.states.SPAWN_CODE_REVIEWER_DEMO_VALIDATION.valid_transitions[] | select(. == "WAITING_FOR_DEMO_VALIDATION")' "$STATE_MACHINE" > /dev/null 2>&1; then
    echo "✅ Correct transition: SPAWN_CODE_REVIEWER_DEMO_VALIDATION → WAITING_FOR_DEMO_VALIDATION"
else
    echo "🔴 MISSING: Transition SPAWN_CODE_REVIEWER_DEMO_VALIDATION → WAITING_FOR_DEMO_VALIDATION"
    ((ERRORS++))
fi

# WAITING_FOR_DEMO_VALIDATION should transition to ERROR_RECOVERY (for failures)
if jq -e '.agents.orchestrator.states.WAITING_FOR_DEMO_VALIDATION.valid_transitions[] | select(. == "ERROR_RECOVERY")' "$STATE_MACHINE" > /dev/null 2>&1; then
    echo "✅ Correct transition: WAITING_FOR_DEMO_VALIDATION → ERROR_RECOVERY (for failures)"
else
    echo "🔴 MISSING: Transition WAITING_FOR_DEMO_VALIDATION → ERROR_RECOVERY"
    ((ERRORS++))
fi

# Check for state rule files
echo ""
echo "📋 Checking for state rule files..."

if [ -f "agent-states/software-factory/orchestrator/SPAWN_CODE_REVIEWER_DEMO_VALIDATION/rules.md" ]; then
    echo "✅ SPAWN_CODE_REVIEWER_DEMO_VALIDATION rules.md exists"
else
    echo "🔴 MISSING: agent-states/software-factory/orchestrator/SPAWN_CODE_REVIEWER_DEMO_VALIDATION/rules.md"
    ((ERRORS++))
fi

if [ -f "agent-states/software-factory/orchestrator/WAITING_FOR_DEMO_VALIDATION/rules.md" ]; then
    echo "✅ WAITING_FOR_DEMO_VALIDATION rules.md exists"
else
    echo "🔴 MISSING: agent-states/software-factory/orchestrator/WAITING_FOR_DEMO_VALIDATION/rules.md"
    ((ERRORS++))
fi

if [ -f "agent-states/code-reviewer/DEMO_VALIDATION/rules.md" ]; then
    echo "✅ DEMO_VALIDATION rules.md exists"
else
    echo "🔴 MISSING: agent-states/code-reviewer/DEMO_VALIDATION/rules.md"
    ((ERRORS++))
fi

# Check orchestrator state schema (SF 3.0)
echo ""
echo "📋 Checking orchestrator state schema (SF 3.0)..."

SCHEMA_FILE="schemas/orchestrator-state-v3.schema.json"

if [ ! -f "$SCHEMA_FILE" ]; then
    echo "⚠️ WARNING: SF 3.0 schema not found at $SCHEMA_FILE (skipping schema checks)"
    echo "   This is expected if running from older template version"
else
    if jq -e '.properties.demo_validation' "$SCHEMA_FILE" > /dev/null 2>&1; then
        echo "✅ demo_validation property in $SCHEMA_FILE"
    else
        echo "⚠️ INFO: demo_validation property not in $SCHEMA_FILE (may not be required)"
    fi

    if jq -e '.properties.state_machine.properties.current_state.enum[] | select(. == "SPAWN_CODE_REVIEWER_DEMO_VALIDATION")' "$SCHEMA_FILE" > /dev/null 2>&1; then
        echo "✅ SPAWN_CODE_REVIEWER_DEMO_VALIDATION in current_state enum"
    else
        echo "⚠️ INFO: SPAWN_CODE_REVIEWER_DEMO_VALIDATION not in current_state enum (may not be required)"
    fi

    if jq -e '.properties.state_machine.properties.current_state.enum[] | select(. == "WAITING_FOR_DEMO_VALIDATION")' "$SCHEMA_FILE" > /dev/null 2>&1; then
        echo "✅ WAITING_FOR_DEMO_VALIDATION in current_state enum"
    else
        echo "⚠️ INFO: WAITING_FOR_DEMO_VALIDATION not in current_state enum (may not be required)"
    fi
fi

# Check R291 updates
echo ""
echo "📋 Checking R291 updates..."

if grep -q "STATE MACHINE ENFORCEMENT" rule-library/R291-integration-demo-requirement.md; then
    echo "✅ R291 includes state machine enforcement section"
else
    echo "🔴 MISSING: R291 state machine enforcement section"
    ((ERRORS++))
fi

# Final summary
echo ""
echo "============================================="
if [ $ERRORS -eq 0 ]; then
    echo "✅✅✅ ALL VALIDATIONS PASSED ✅✅✅"
    echo ""
    echo "Demo enforcement is properly configured!"
    echo "R291 Gate 4 cannot be bypassed via state machine."
    echo ""
    exit 0
else
    echo "🔴🔴🔴 VALIDATION FAILED: $ERRORS ERRORS 🔴🔴🔴"
    echo ""
    echo "Demo enforcement has configuration issues."
    echo "Please fix the errors listed above."
    echo ""
    exit 1
fi
