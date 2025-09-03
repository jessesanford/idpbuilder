#!/bin/bash

# Ensure all agent state directories and rules files exist
# This implements the requirement that every state must have a rules.md file

AGENT_STATES_DIR="/workspaces/software-factory-2.0-template/agent-states"

# Define all states for each agent type based on SOFTWARE-FACTORY-STATE-MACHINE.md

# Orchestrator states
ORCHESTRATOR_STATES=(
    "INIT"
    "WAVE_START"
    "SETUP_EFFORT_INFRASTRUCTURE"
    "SPAWN_CODE_REVIEWERS_EFFORT_PLANNING"
    "WAITING_FOR_EFFORT_PLANS"
    "SPAWN_AGENTS"
    "SPAWN_ARCHITECT_PHASE_PLANNING"
    "SPAWN_ARCHITECT_WAVE_PLANNING"
    "SPAWN_CODE_REVIEWER_PHASE_IMPL"
    "SPAWN_CODE_REVIEWER_WAVE_IMPL"
    "WAITING_FOR_ARCHITECTURE_PLAN"
    "WAITING_FOR_IMPLEMENTATION_PLAN"
    "INJECT_WAVE_METADATA"
    "MONITOR"
    "WAVE_COMPLETE"
    "INTEGRATION"
    "WAVE_REVIEW"
    "ERROR_RECOVERY"
    "PLANNING"
    "SUCCESS"
    "HARD_STOP"
)

# SW Engineer states
SW_ENGINEER_STATES=(
    "INIT"
    "IMPLEMENTATION"
    "MEASURE_SIZE"
    "FIX_ISSUES"
    "SPLIT_IMPLEMENTATION"
    "TEST_WRITING"
    "REQUEST_REVIEW"
    "COMPLETED"
    "BLOCKED"
)

# Code Reviewer states
CODE_REVIEWER_STATES=(
    "INIT"
    "PHASE_IMPLEMENTATION_PLANNING"
    "WAVE_IMPLEMENTATION_PLANNING"
    "WAVE_DIRECTORY_ACKNOWLEDGMENT"
    "EFFORT_PLAN_CREATION"
    "CODE_REVIEW"
    "CREATE_SPLIT_PLAN"
    "SPLIT_REVIEW"
    "VALIDATION"
    "COMPLETED"
)

# Architect states
ARCHITECT_STATES=(
    "INIT"
    "PHASE_ARCHITECTURE_PLANNING"
    "WAVE_ARCHITECTURE_PLANNING"
    "PHASE_DIRECTORY_ACKNOWLEDGMENT"
    "WAVE_REVIEW"
    "PHASE_ASSESSMENT"
    "INTEGRATION_REVIEW"
    "ARCHITECTURE_AUDIT"
    "ARCHITECTURE_VALIDATION"
    "DECISION"
)

# Function to create a minimal rules.md file
create_minimal_rules_file() {
    local agent_type="$1"
    local state="$2"
    local file_path="$3"
    
    cat > "$file_path" << EOF
# ${agent_type^} - $state State Rules

## State Context
This is the $state state for the ${agent_type}.

## Acknowledgment Required
Thank you for reading the rules file for the $state state.

**IMPORTANT**: Please report that you have successfully read the $state rules file.

Say: "✅ Successfully read $state rules for ${agent_type}"

## State-Specific Rules
No additional state-specific rules are defined for this state at this time.

## General Responsibilities
Follow all general ${agent_type} rules and the Software Factory state machine.

## Next Steps
Proceed with the standard workflow for the $state state as defined in the state machine.
EOF
    
    echo "   Created minimal rules.md"
}

# Function to ensure state directory and rules file exist
ensure_state_rules() {
    local agent_type="$1"
    local state="$2"
    
    local dir_path="$AGENT_STATES_DIR/$agent_type/$state"
    local rules_file="$dir_path/rules.md"
    
    # Create directory if it doesn't exist
    if [ ! -d "$dir_path" ]; then
        mkdir -p "$dir_path"
        echo "✅ Created directory: $agent_type/$state/"
    fi
    
    # Create rules.md if it doesn't exist
    if [ ! -f "$rules_file" ]; then
        echo "📝 Creating rules.md for $agent_type/$state..."
        create_minimal_rules_file "$agent_type" "$state" "$rules_file"
    else
        echo "✓ Rules file exists: $agent_type/$state/rules.md"
    fi
}

echo "=========================================="
echo "Ensuring All Agent State Rules Files Exist"
echo "=========================================="
echo ""

# Process Orchestrator states
echo "🎯 ORCHESTRATOR States:"
echo "------------------------"
for state in "${ORCHESTRATOR_STATES[@]}"; do
    ensure_state_rules "orchestrator" "$state"
done
echo ""

# Process SW Engineer states
echo "💻 SW-ENGINEER States:"
echo "----------------------"
for state in "${SW_ENGINEER_STATES[@]}"; do
    ensure_state_rules "sw-engineer" "$state"
done
echo ""

# Process Code Reviewer states
echo "📋 CODE-REVIEWER States:"
echo "------------------------"
for state in "${CODE_REVIEWER_STATES[@]}"; do
    ensure_state_rules "code-reviewer" "$state"
done
echo ""

# Process Architect states
echo "🏗️ ARCHITECT States:"
echo "--------------------"
for state in "${ARCHITECT_STATES[@]}"; do
    ensure_state_rules "architect" "$state"
done
echo ""

# Final verification
echo "=========================================="
echo "VERIFICATION SUMMARY"
echo "=========================================="
echo ""

# Count total states and verify
total_expected=$((${#ORCHESTRATOR_STATES[@]} + ${#SW_ENGINEER_STATES[@]} + ${#CODE_REVIEWER_STATES[@]} + ${#ARCHITECT_STATES[@]}))
total_found=$(find "$AGENT_STATES_DIR" -name "rules.md" | wc -l)

echo "Expected total state rules files: $total_expected"
echo "Found total state rules files: $total_found"
echo ""

if [ "$total_found" -ge "$total_expected" ]; then
    echo "✅ SUCCESS: All required state rules files exist!"
else
    echo "⚠️  WARNING: Some rules files may be missing"
    echo "   Expected: $total_expected"
    echo "   Found: $total_found"
fi

echo ""
echo "You can verify manually with:"
echo "  find $AGENT_STATES_DIR -type d -name '*' | while read dir; do"
echo "    [ ! -f \"\$dir/rules.md\" ] && [ -d \"\$dir\" ] && echo \"Missing: \$dir/rules.md\""
echo "  done"