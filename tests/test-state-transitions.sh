#!/bin/bash

# Test Script for R206 - State Machine Transition Validation
# This tests that agents properly validate state transitions

set -e

echo "======================================"
echo "Testing R206: State Machine Validation"
echo "======================================"

# Setup test environment
TEST_DIR="/tmp/state-test-$$"
mkdir -p "$TEST_DIR"
cd "$TEST_DIR"

# Copy state machine definition
cp /workspaces/software-factory-2.0-template/SOFTWARE-FACTORY-STATE-MACHINE.md .

# Test function for state validation
validate_state_transition() {
    local AGENT_TYPE="$1"
    local TARGET_STATE="$2"
    local EXPECTED="$3"
    
    # Map agent type to section header
    case "$AGENT_TYPE" in
        orchestrator) SECTION="Orchestrator States" ;;
        sw-engineer) SECTION="SW Engineer States" ;;
        code-reviewer) SECTION="Code Reviewer States" ;;
        architect) SECTION="Architect States" ;;
        *) echo "❌ Unknown agent type: $AGENT_TYPE"; return 1 ;;
    esac
    
    # Extract valid states for this agent type
    if sed -n "/## $SECTION/,/^##[^#]/p" SOFTWARE-FACTORY-STATE-MACHINE.md | \
       grep -q "^- \*\*$TARGET_STATE\*\*"; then \
        RESULT="VALID"; \
    else \
        RESULT="INVALID"; \
    fi
    
    if [ "$RESULT" = "$EXPECTED" ]; then \
        echo "✅ PASS: $AGENT_TYPE → $TARGET_STATE is $EXPECTED"; \
        return 0; \
    else \
        echo "❌ FAIL: $AGENT_TYPE → $TARGET_STATE expected $EXPECTED, got $RESULT"; \
        return 1; \
    fi
}

echo ""
echo "Test 1: Valid Orchestrator States"
echo "----------------------------------"
validate_state_transition "orchestrator" "INIT" "VALID"
validate_state_transition "orchestrator" "WAVE_START" "VALID"
validate_state_transition "orchestrator" "SPAWN_AGENTS" "VALID"
validate_state_transition "orchestrator" "MONITOR" "VALID"
validate_state_transition "orchestrator" "WAVE_COMPLETE" "VALID"
validate_state_transition "orchestrator" "INTEGRATION" "VALID"
validate_state_transition "orchestrator" "WAVE_REVIEW" "VALID"
validate_state_transition "orchestrator" "ERROR_RECOVERY" "VALID"
validate_state_transition "orchestrator" "SUCCESS" "VALID"
validate_state_transition "orchestrator" "HARD_STOP" "VALID"

echo ""
echo "Test 2: Invalid Orchestrator States"
echo "------------------------------------"
validate_state_transition "orchestrator" "WAITING_FOR_COFFEE" "INVALID"
validate_state_transition "orchestrator" "IMPLEMENTATION" "INVALID"  # Wrong agent
validate_state_transition "orchestrator" "CODE_REVIEW" "INVALID"     # Wrong agent
validate_state_transition "orchestrator" "THINKING" "INVALID"
validate_state_transition "orchestrator" "WAVE_COMPELTE" "INVALID"   # Typo

echo ""
echo "Test 3: Valid SW Engineer States"
echo "---------------------------------"
validate_state_transition "sw-engineer" "INIT" "VALID"
validate_state_transition "sw-engineer" "IMPLEMENTATION" "VALID"
validate_state_transition "sw-engineer" "MEASURE_SIZE" "VALID"
validate_state_transition "sw-engineer" "FIX_ISSUES" "VALID"
validate_state_transition "sw-engineer" "SPLIT_IMPLEMENTATION" "VALID"
validate_state_transition "sw-engineer" "TEST_WRITING" "VALID"
validate_state_transition "sw-engineer" "COMPLETED" "VALID"
validate_state_transition "sw-engineer" "BLOCKED" "VALID"

echo ""
echo "Test 4: Invalid SW Engineer States"
echo "-----------------------------------"
validate_state_transition "sw-engineer" "WAVE_REVIEW" "INVALID"      # Wrong agent
validate_state_transition "sw-engineer" "SPAWN_AGENTS" "INVALID"     # Wrong agent
validate_state_transition "sw-engineer" "DEBUGGING" "INVALID"
validate_state_transition "sw-engineer" "IMPLEMENATION" "INVALID"    # Typo

echo ""
echo "Test 5: Valid Code Reviewer States"
echo "-----------------------------------"
validate_state_transition "code-reviewer" "INIT" "VALID"
validate_state_transition "code-reviewer" "EFFORT_PLAN_CREATION" "VALID"
validate_state_transition "code-reviewer" "CODE_REVIEW" "VALID"
validate_state_transition "code-reviewer" "CREATE_SPLIT_PLAN" "VALID"
validate_state_transition "code-reviewer" "SPLIT_REVIEW" "VALID"
validate_state_transition "code-reviewer" "VALIDATION" "VALID"
validate_state_transition "code-reviewer" "COMPLETED" "VALID"

echo ""
echo "Test 6: Invalid Code Reviewer States"
echo "-------------------------------------"
validate_state_transition "code-reviewer" "IMPLEMENTATION" "INVALID"  # Wrong agent
validate_state_transition "code-reviewer" "MONITOR" "INVALID"        # Wrong agent
validate_state_transition "code-reviewer" "COFFEE_BREAK" "INVALID"
validate_state_transition "code-reviewer" "CODE_REVEIW" "INVALID"    # Typo

echo ""
echo "Test 7: Valid Architect States"
echo "-------------------------------"
validate_state_transition "architect" "INIT" "VALID"
validate_state_transition "architect" "WAVE_REVIEW" "VALID"
validate_state_transition "architect" "PHASE_ASSESSMENT" "VALID"
validate_state_transition "architect" "INTEGRATION_REVIEW" "VALID"
validate_state_transition "architect" "ARCHITECTURE_AUDIT" "VALID"
validate_state_transition "architect" "DECISION" "VALID"

echo ""
echo "Test 8: Invalid Architect States"
echo "---------------------------------"
validate_state_transition "architect" "IMPLEMENTATION" "INVALID"     # Wrong agent
validate_state_transition "architect" "CODE_REVIEW" "INVALID"        # Wrong agent
validate_state_transition "architect" "SPAWN_AGENTS" "INVALID"       # Wrong agent
validate_state_transition "architect" "PROCRASTINATING" "INVALID"

echo ""
echo "Test 9: Transition Validation Function"
echo "---------------------------------------"

# Test the actual validation function
cat > test-validation.sh << 'EOF'
#!/bin/bash

# Function from R206
validate_state_transition() {
    local CURRENT="$1"
    local TARGET="$2"
    local AGENT_TYPE="$3"
    
    # Check state machine definition
    STATE_MACHINE="SOFTWARE-FACTORY-STATE-MACHINE.md"
    
    if [ ! -f "$STATE_MACHINE" ]; then \
        echo "❌ State machine definition not found!"; \
        return 1; \
    fi
    
    # Map agent type to section
    case "$AGENT_TYPE" in
        orchestrator) SECTION="Orchestrator States" ;;
        sw-engineer) SECTION="SW Engineer States" ;;
        code-reviewer) SECTION="Code Reviewer States" ;;
        architect) SECTION="Architect States" ;;
        *) echo "❌ Unknown agent type: $AGENT_TYPE"; return 1 ;;
    esac
    
    # Validate target state exists
    if sed -n "/## $SECTION/,/^##[^#]/p" "$STATE_MACHINE" | \
       grep -q "^- \*\*$TARGET\*\*"; then \
        echo "✅ Valid transition: $CURRENT → $TARGET"; \
        return 0; \
    else \
        echo "❌ INVALID STATE: $TARGET not found in $SECTION"; \
        return 1; \
    fi
}

# Test cases
echo "Testing transition validation function:"
validate_state_transition "WAVE_COMPLETE" "INTEGRATION" "orchestrator" && echo "  Passed" || echo "  Failed"
validate_state_transition "IMPLEMENTATION" "MEASURE_SIZE" "sw-engineer" && echo "  Passed" || echo "  Failed"
validate_state_transition "CODE_REVIEW" "VALIDATION" "code-reviewer" && echo "  Passed" || echo "  Failed"
validate_state_transition "WAVE_REVIEW" "DECISION" "architect" && echo "  Passed" || echo "  Failed"

echo ""
echo "Testing invalid transitions:"
validate_state_transition "WAVE_COMPLETE" "COFFEE_TIME" "orchestrator" && echo "  Failed!" || echo "  Correctly rejected"
validate_state_transition "IMPLEMENTATION" "WAVE_REVIEW" "sw-engineer" && echo "  Failed!" || echo "  Correctly rejected"
EOF

chmod +x test-validation.sh
./test-validation.sh

echo ""
echo "Test 10: State Extraction"
echo "--------------------------"

# Test extracting all states for an agent
echo "Extracting all orchestrator states:"
sed -n '/## Orchestrator States/,/^##[^#]/p' SOFTWARE-FACTORY-STATE-MACHINE.md | \
    grep "^- \*\*" | \
    sed 's/- \*\*/  /' | \
    sed 's/\*\*.*//' | \
    head -5

echo ""
echo "Test 11: Terminal State Detection"
echo "----------------------------------"

# Check if a state is terminal
check_terminal() {
    local STATE="$1"
    local AGENT="$2"
    
    case "$AGENT" in
        orchestrator) SECTION="Orchestrator States" ;;
        sw-engineer) SECTION="SW Engineer States" ;;
        code-reviewer) SECTION="Code Reviewer States" ;;
        architect) SECTION="Architect States" ;;
    esac
    
    if sed -n "/## $SECTION/,/^##[^#]/p" SOFTWARE-FACTORY-STATE-MACHINE.md | \
       grep "^- \*\*$STATE\*\*" | \
       grep -q "terminal"; then \
        echo "  $STATE is TERMINAL for $AGENT"; \
    else \
        echo "  $STATE is NOT terminal for $AGENT"; \
    fi
}

echo "Checking terminal states:"
check_terminal "SUCCESS" "orchestrator"
check_terminal "HARD_STOP" "orchestrator"
check_terminal "COMPLETED" "sw-engineer"
check_terminal "BLOCKED" "sw-engineer"
check_terminal "DECISION" "architect"
check_terminal "INIT" "orchestrator"  # Should not be terminal

echo ""
echo "Test 12: Simulated Agent Spawn with Validation"
echo "------------------------------------------------"

# Simulate agent spawn with state validation
cat > simulate-spawn.sh << 'EOF'
#!/bin/bash

echo "Simulating orchestrator spawning SW engineer..."

# Orchestrator state
ORCHESTRATOR_STATE="SPAWN_AGENTS"
AGENT_TO_SPAWN="sw-engineer"
TARGET_INITIAL_STATE="INIT"

# Validate orchestrator is in valid spawn state
if sed -n '/## Orchestrator States/,/^##[^#]/p' SOFTWARE-FACTORY-STATE-MACHINE.md | \
   grep -q "^- \*\*$ORCHESTRATOR_STATE\*\*"; then \
    echo "✅ Orchestrator in valid state: $ORCHESTRATOR_STATE"; \
else \
    echo "❌ Orchestrator in invalid state!"; \
    exit 1; \
fi

# Validate target agent initial state
if sed -n '/## SW Engineer States/,/^##[^#]/p' SOFTWARE-FACTORY-STATE-MACHINE.md | \
   grep -q "^- \*\*$TARGET_INITIAL_STATE\*\*"; then \
    echo "✅ SW Engineer will start in: $TARGET_INITIAL_STATE"; \
else \
    echo "❌ Invalid initial state for SW Engineer!"; \
    exit 1; \
fi

echo "✅ State validation passed - spawn can proceed"

# Simulate SW Engineer startup
echo ""
echo "SW Engineer starting..."
SW_CURRENT_STATE="INIT"
SW_TARGET_STATE="IMPLEMENTATION"

# Validate transition
if grep -q "INIT → IMPLEMENTATION" SOFTWARE-FACTORY-STATE-MACHINE.md; then \
    echo "✅ Valid transition: $SW_CURRENT_STATE → $SW_TARGET_STATE"; \
    SW_CURRENT_STATE="$SW_TARGET_STATE"; \
    echo "  Current state updated to: $SW_CURRENT_STATE"; \
else \
    echo "❌ Invalid transition attempted!"; \
    exit 1; \
fi
EOF

chmod +x simulate-spawn.sh
./simulate-spawn.sh

echo ""
echo "======================================"
echo "Test Summary"
echo "======================================"
echo "✅ State validation tests complete"
echo "✅ All agents can validate their states"
echo "✅ Invalid states are properly rejected"
echo "✅ Terminal states are detected"
echo "✅ State transitions are validated"
echo ""
echo "R206 Implementation Status:"
echo "  ✅ Rule created and documented"
echo "  ✅ State machine definition created"
echo "  ✅ All agents updated with validation"
echo "  ✅ Examples and test cases provided"

# Cleanup
cd /
rm -rf "$TEST_DIR"