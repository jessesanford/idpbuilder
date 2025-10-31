#!/bin/bash

# Script to fix R322 misinterpretations across all state files
# Created by Software Factory Manager to enable automatic continuation

echo "🏭 SOFTWARE FACTORY MANAGER - R322 FIX SCRIPT"
echo "=============================================="
echo "Fixing incorrect R322 interpretations that cause unnecessary stops"
echo ""

# States that ACTUALLY require stops per R322
STATES_REQUIRING_STOPS=(
    # Part A: After spawning (context preservation)
    "SPAWN_SW_ENGINEERS"
    "SPAWN_CODE_REVIEWERS_EFFORT_PLANNING"
    "SPAWN_CODE_REVIEWERS_EFFORT_REVIEW"
    "SPAWN_SW_ENGINEERS"
    "SPAWN_INTEGRATION_AGENT"
    "SPAWN_INTEGRATION_AGENT_PHASE"
    "SPAWN_INTEGRATION_AGENT_PROJECT"
    "SPAWN_ARCHITECT_*"
    "SPAWN_CODE_REVIEWER_*"
    "SPAWN_SW_ENGINEER_*"

    # Part B: Critical integration plan reviews (NOT normal operations)
    "WAITING_FOR_MERGE_PLAN"  # → SPAWN_INTEGRATION_AGENT
    "WAITING_FOR_PHASE_MERGE_PLAN"  # → SPAWN_INTEGRATION_AGENT_PHASE
    "WAITING_FOR_PROJECT_MERGE_PLAN"  # → SPAWN_INTEGRATION_AGENT_PROJECT

    # Part C: Assessment → Action
    "WAITING_FOR_PHASE_ASSESSMENT"  # → COMPLETE_PHASE
    "WAITING_FOR_PROJECT_VALIDATION"  # → CREATE_INTEGRATE_WAVE_EFFORTS_TESTING

    # Part D: Exceptional monitoring stops
    "MONITORING_INTEGRATE_WAVE_EFFORTS"  # → REVIEW_WAVE_ARCHITECTURE
    "MONITORING_INTEGRATE_PHASE_WAVES"  # → SPAWN_ARCHITECT_PHASE_ASSESSMENT
)

# States that should CONTINUE AUTOMATICALLY (normal operations)
STATES_CONTINUING_AUTOMATICALLY=(
    "WAVE_COMPLETE"  # → INTEGRATE_WAVE_EFFORTS (normal)
    "COMPLETE_PHASE"  # → START_PHASE_ITERATION (normal)
    "WAVE_START"  # → Next state (normal)
    "INTEGRATE_WAVE_EFFORTS"  # → Next state (normal)
    "MONITORING_SWE_PROGRESS"  # → Next state (normal)
    "MONITORING_EFFORT_FIXES"  # → Next state (normal)
    "WAITING_FOR_FIX_PLANS"  # FIX PLANS ARE NORMAL OPERATIONS
    "WAITING_FOR_BACKPORT_PLAN"  # BACKPORTS ARE NORMAL OPERATIONS
    "BUILD_VALIDATION"  # Normal validation
    "INIT"  # Initial state should continue
    "ERROR_RECOVERY"  # Should continue if recovery succeeds
)

# Function to check if a state requires stops
requires_stop() {
    local state="$1"
    for stop_state in "${STATES_REQUIRING_STOPS[@]}"; do
        if [[ "$state" == "$stop_state" ]] || [[ "$state" =~ ^${stop_state}\* ]]; then
            return 0  # True - requires stop
        fi
    done
    return 1  # False - should continue automatically
}

# Function to fix a rules.md file
fix_rules_file() {
    local file="$1"
    local state_name=$(basename $(dirname "$file"))

    echo "Processing: $state_name"

    # Check if this state should have stops
    if requires_stop "$state_name"; then
        echo "  ✓ $state_name correctly requires stops per R322"
        return 0
    fi

    # This state should NOT have mandatory stops - fix it!
    echo "  🔧 Fixing $state_name to allow automatic continuation..."

    # Create a temporary file with the fix
    cat > /tmp/r322_fix_template.txt << 'EOF'
# REPLACE_STATE_NAME State Rules

## ✅ NORMAL OPERATION - AUTOMATIC CONTINUATION

**Per R322: Normal state transitions do NOT require stops!**
**This is NORMAL SOFTWARE DEVELOPMENT FLOW**

### AUTOMATIC FLOW:
1. ✅ Complete all work for this state
2. ✅ Update orchestrator-state-v3.json to next state
3. ✅ Commit and push the state file
4. ✅ Output CONTINUE-SOFTWARE-FACTORY=TRUE
5. ✅ System continues automatically

### R322 CLARIFICATION:
- ✅ This state transition is NOT in R322's checkpoint list
- ✅ This is NORMAL software development flow
- ✅ Should continue AUTOMATICALLY for full automation
- ✅ Only EXCEPTIONAL situations listed in R322 require stops

### NORMAL COMPLETION PROTOCOL:
```bash
# Complete state work
echo "✅ State work completed successfully"

# Update to next state
NEXT_STATE="[DETERMINE_NEXT_STATE]"
jq --arg state "$NEXT_STATE" '.current_state = $state' orchestrator-state-v3.json > tmp.json && mv tmp.json orchestrator-state-v3.json

# Commit state transition
git add orchestrator-state-v3.json
git commit -m "state: automatic transition to $NEXT_STATE"
git push

# Output continuation flag for automation
echo "✅ Continuing to $NEXT_STATE automatically"
echo "CONTINUE-SOFTWARE-FACTORY=TRUE"  # ENABLE AUTOMATIC CONTINUATION
```

**CONTINUE AUTOMATICALLY - Full factory automation enabled!**

---

EOF

    # Replace the state name in the template
    sed -i "s/REPLACE_STATE_NAME/$state_name/" /tmp/r322_fix_template.txt

    # Check if file has the wrong R322 section
    if grep -q "R322 MANDATORY STOP BEFORE STATE TRANSITIONS" "$file"; then
        # Extract content before and after the R322 section
        # Find line numbers
        START_LINE=$(grep -n "^## 🛑🛑🛑 R322 MANDATORY STOP" "$file" | cut -d: -f1 | head -1)
        END_LINE=$(grep -n "^---$" "$file" | while read line_num; do
            num=$(echo "$line_num" | cut -d: -f1)
            if [ "$num" -gt "$START_LINE" ]; then
                echo "$num"
                break
            fi
        done)

        if [ -n "$START_LINE" ] && [ -n "$END_LINE" ]; then
            # Create new file with fixed content
            TEMP_FILE="/tmp/fixed_rules_$$.md"

            # Get content before R322 section (if any)
            if [ "$START_LINE" -gt 1 ]; then
                head -n $((START_LINE - 1)) "$file" > "$TEMP_FILE"
            else
                touch "$TEMP_FILE"
            fi

            # Add fixed R322 section
            cat /tmp/r322_fix_template.txt >> "$TEMP_FILE"

            # Add content after R322 section
            TOTAL_LINES=$(wc -l < "$file")
            if [ "$END_LINE" -lt "$TOTAL_LINES" ]; then
                tail -n $((TOTAL_LINES - END_LINE)) "$file" >> "$TEMP_FILE"
            fi

            # Replace original file
            mv "$TEMP_FILE" "$file"
            echo "  ✅ Fixed $state_name rules.md"
        else
            echo "  ⚠️ Could not determine R322 section boundaries in $file"
        fi
    else
        echo "  ℹ️ No wrong R322 section found in $file"
    fi
}

# Main execution
echo ""
echo "Scanning all orchestrator state files..."
echo ""

# Process all rules.md files in orchestrator states
for rules_file in /home/vscode/software-factory-template/agent-states/software-factory/orchestrator/*/rules.md; do
    if [ -f "$rules_file" ]; then
        fix_rules_file "$rules_file"
    fi
done

echo ""
echo "🏭 R322 Fix Script Complete!"
echo ""
echo "Summary:"
echo "- Fixed states to allow automatic continuation"
echo "- Preserved required stops for spawn and critical review states"
echo "- Enabled full Software Factory automation"
echo ""
echo "Next steps:"
echo "1. Review the changes"
echo "2. Test with a sample orchestrator run"
echo "3. Commit and push the fixes"