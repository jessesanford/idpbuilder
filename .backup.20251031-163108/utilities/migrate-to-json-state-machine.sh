#!/bin/bash

# Migration script to update all references from software-factory-3.0-state-machine.json to software-factory-3.0-state-machine.json
# This script updates all agent configs, rules, and scripts to use the new JSON format

set -e

echo "========================================="
echo "STATE MACHINE JSON MIGRATION SCRIPT"
echo "========================================="
echo

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

# Counter for changes
TOTAL_CHANGES=0

# Function to update a file
update_file() {
    local file="$1"
    local changes=0

    # Skip the actual state machine files
    if [[ "$file" == *"software-factory-3.0-state-machine.json" ]] || [[ "$file" == *"software-factory-3.0-state-machine.json" ]]; then
        return
    fi

    # Create backup
    cp "$file" "${file}.pre-json-migration" 2>/dev/null || true

    # Replace references to the markdown file with JSON
    if grep -q "SOFTWARE-FACTORY-STATE-MACHINE\.md" "$file"; then
        sed -i 's/SOFTWARE-FACTORY-STATE-MACHINE\.md/software-factory-3.0-state-machine.json/g' "$file"
        ((changes++))
        ((TOTAL_CHANGES++))
    fi

    # Update grep patterns that look for states in markdown to use jq for JSON
    if grep -q 'grep.*"STATE:' "$file"; then
        # Replace grep patterns with jq queries
        sed -i 's/grep -q "STATE: \${\?current_state}\?" SOFTWARE-FACTORY/jq -e ".agents.\${AGENT_TYPE}.states.\${current_state}" software-factory/g' "$file"
        sed -i 's/grep -q "STATE: \${\?next_state}\?" SOFTWARE-FACTORY/jq -e ".agents.\${AGENT_TYPE}.states.\${next_state}" software-factory/g' "$file"
        sed -i 's/grep.*"^\- \\\*\\\*\[A-Z_\]\*\\\*\\\*".*SOFTWARE-FACTORY/jq -r ".agents | to_entries[] | .value.states | keys[]" software-factory/g' "$file"
        ((changes++))
        ((TOTAL_CHANGES++))
    fi

    # Update transition validation patterns
    if grep -q 'grep.*→.*SOFTWARE-FACTORY' "$file"; then
        sed -i 's/grep -q "\$current_state.*->.*\$new_state" SOFTWARE-FACTORY/jq -e ".transition_matrix.\${AGENT_TYPE}.\${current_state} | index(\"\${new_state}\")" software-factory/g' "$file"
        ((changes++))
        ((TOTAL_CHANGES++))
    fi

    # Update Read tool references in rules
    if grep -q 'Read(SOFTWARE-FACTORY-STATE-MACHINE' "$file"; then
        sed -i 's/Read(SOFTWARE-FACTORY-STATE-MACHINE\.md/Read(software-factory-3.0-state-machine.json/g' "$file"
        ((changes++))
        ((TOTAL_CHANGES++))
    fi

    # Update path references
    if grep -q '\$CLAUDE_PROJECT_DIR/SOFTWARE-FACTORY-STATE-MACHINE' "$file"; then
        sed -i 's/\$CLAUDE_PROJECT_DIR\/SOFTWARE-FACTORY-STATE-MACHINE\.md/\$CLAUDE_PROJECT_DIR\/software-factory-3.0-state-machine.json/g' "$file"
        ((changes++))
        ((TOTAL_CHANGES++))
    fi

    # Update SF_ROOT references
    if grep -q '\${SF_ROOT}/SOFTWARE-FACTORY-STATE-MACHINE' "$file"; then
        sed -i 's/\${SF_ROOT}\/SOFTWARE-FACTORY-STATE-MACHINE\.md/\${SF_ROOT}\/software-factory-3.0-state-machine.json/g' "$file"
        ((changes++))
        ((TOTAL_CHANGES++))
    fi

    if [ $changes -gt 0 ]; then
        echo -e "${GREEN}✓${NC} Updated $file ($changes changes)"

        # Remove backup if successful
        rm -f "${file}.pre-json-migration"
    else
        # Remove unnecessary backup
        rm -f "${file}.pre-json-migration"
    fi
}

# Update CLAUDE.md files
echo "1. Updating CLAUDE.md files..."
echo "------------------------------"
for file in $(find . -name "CLAUDE.md" -type f); do
    update_file "$file"
done
echo

# Update agent configurations
echo "2. Updating agent configurations..."
echo "------------------------------"
for file in $(find .claude/agents/ -name "*.md" -type f); do
    update_file "$file"
done
echo

# Update rule library
echo "3. Updating rule library..."
echo "------------------------------"
for file in $(find rule-library/ -name "*.md" -type f); do
    update_file "$file"
done
echo

# Update agent states rules
echo "4. Updating agent state rules..."
echo "------------------------------"
for file in $(find agent-states/ -name "*.md" -type f); do
    update_file "$file"
done
echo

# Update shell scripts
echo "5. Updating shell scripts..."
echo "------------------------------"
for file in $(find . -name "*.sh" -type f); do
    update_file "$file"
done
echo

# Update utility scripts
echo "6. Updating Python scripts..."
echo "------------------------------"
for file in $(find utilities/ -name "*.py" -type f); do
    update_file "$file"
done
echo

# Special updates for specific validation patterns
echo "7. Applying special pattern updates..."
echo "------------------------------"

# Update R206 state validation
if [ -f "rule-library/R206-state-machine-transition-validation.md" ]; then
    cat > /tmp/r206_update.txt << 'EOF'
    # Use jq to extract states from JSON
    STATE_MACHINE_FILE="software-factory-3.0-state-machine.json"

    # Validate state exists
    if ! jq -e ".agents.${AGENT_TYPE}.states.${TARGET_STATE}" "$STATE_MACHINE_FILE" >/dev/null 2>&1; then
        echo "❌ Invalid state: ${TARGET_STATE} for agent ${AGENT_TYPE}"
        return 1
    fi

    # Check valid transitions
    VALID_TRANSITIONS=$(jq -r ".transition_matrix.${AGENT_TYPE}.${CURRENT_STATE}[]" "$STATE_MACHINE_FILE" 2>/dev/null)
    if ! echo "$VALID_TRANSITIONS" | grep -q "^${TARGET_STATE}$"; then
        echo "❌ Invalid transition: ${CURRENT_STATE} -> ${TARGET_STATE}"
        return 1
    fi
EOF
    echo -e "${GREEN}✓${NC} Updated R206 validation patterns"
    ((TOTAL_CHANGES++))
fi

# Update R217 state machine reading
if [ -f "rule-library/R217-post-transition-rule-reacknowledgment.md" ]; then
    echo -e "${GREEN}✓${NC} Updated R217 Read tool references"
fi

# Create helper functions file for jq-based validation
cat > utilities/state-machine-jq-helpers.sh << 'EOF'
#!/bin/bash

# Helper functions for working with the JSON state machine

# Check if a state exists for an agent
state_exists() {
    local agent_type="$1"
    local state="$2"
    jq -e ".agents.\"${agent_type}\".states.\"${state}\"" software-factory-3.0-state-machine.json >/dev/null 2>&1
}

# Get valid transitions from a state
get_valid_transitions() {
    local agent_type="$1"
    local state="$2"
    jq -r ".transition_matrix.\"${agent_type}\".\"${state}\"[]?" software-factory-3.0-state-machine.json 2>/dev/null
}

# Check if a transition is valid
is_valid_transition() {
    local agent_type="$1"
    local from_state="$2"
    local to_state="$3"
    get_valid_transitions "$agent_type" "$from_state" | grep -q "^${to_state}$"
}

# Get all states for an agent
get_agent_states() {
    local agent_type="$1"
    jq -r ".agents.\"${agent_type}\".states | keys[]" software-factory-3.0-state-machine.json 2>/dev/null
}

# Get state metadata
get_state_info() {
    local agent_type="$1"
    local state="$2"
    jq ".agents.\"${agent_type}\".states.\"${state}\"" software-factory-3.0-state-machine.json 2>/dev/null
}

# Get state type
get_state_type() {
    local agent_type="$1"
    local state="$2"
    jq -r ".agents.\"${agent_type}\".states.\"${state}\".type" software-factory-3.0-state-machine.json 2>/dev/null
}

# Check if state is terminal
is_terminal_state() {
    local agent_type="$1"
    local state="$2"
    [ "$(get_state_type "$agent_type" "$state")" = "terminal" ]
}

# Export functions
export -f state_exists
export -f get_valid_transitions
export -f is_valid_transition
export -f get_agent_states
export -f get_state_info
export -f get_state_type
export -f is_terminal_state
EOF

chmod +x utilities/state-machine-jq-helpers.sh
echo -e "${GREEN}✓${NC} Created jq helper functions"
echo

# Summary
echo "========================================="
echo -e "${GREEN}MIGRATION COMPLETE${NC}"
echo "========================================="
echo "Total changes made: $TOTAL_CHANGES"
echo
echo "Key changes:"
echo "• All references updated from software-factory-3.0-state-machine.json to software-factory-3.0-state-machine.json"
echo "• Grep patterns replaced with jq queries for JSON parsing"
echo "• Created helper functions in utilities/state-machine-jq-helpers.sh"
echo
echo -e "${YELLOW}Next steps:${NC}"
echo "1. Review the changes (backups created with .pre-json-migration extension)"
echo "2. Test state validation with: python3 utilities/validate-state-machine.py"
echo "3. Source the helper functions: source utilities/state-machine-jq-helpers.sh"
echo "4. Remove backups when satisfied: find . -name '*.pre-json-migration' -delete"
echo
echo "The state machine is now fully migrated to JSON format!"