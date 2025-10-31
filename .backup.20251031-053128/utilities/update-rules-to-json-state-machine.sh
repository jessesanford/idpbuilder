#!/bin/bash

# Update all rule references to use JSON state machines instead of markdown

set -e

echo "========================================="
echo "UPDATING RULES TO JSON STATE MACHINE"
echo "========================================="
echo

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

# Counter
CHANGES=0

# Function to update a rule file
update_rule_file() {
    local file="$1"
    local local_changes=0

    echo "Updating: $file"

    # Create backup
    cp "$file" "${file}.bak"

    # Replace markdown references with JSON
    # Main state machine references
    sed -i 's|SOFTWARE-FACTORY-STATE-MACHINE\.md|software-factory-3.0-state-machine.json|g' "$file"

    # Update grep patterns to use jq (use different delimiter to avoid issues)
    sed -i 's#grep -q "STATE: ${next_state}" software-factory-state-machine\.json#jq -e ".agents.${AGENT_TYPE}.states.${next_state}" software-factory-3.0-state-machine.json >/dev/null 2>\&1#g' "$file"
    sed -i 's#grep -q "STATE: ${TARGET_STATE}" software-factory-state-machine\.json#jq -e ".agents.${AGENT_TYPE}.states.${TARGET_STATE}" software-factory-3.0-state-machine.json >/dev/null 2>\&1#g' "$file"
    sed -i 's#grep -q "STATE: $current_state" software-factory-state-machine\.json#jq -e ".agents.orchestrator.states.${current_state}" software-factory-3.0-state-machine.json >/dev/null 2>\&1#g' "$file"

    # Update transition validation patterns
    sed -i 's#grep -q "$current_state.*->.*$new_state" software-factory-state-machine\.json#jq -e ".transition_matrix.orchestrator.${current_state} | index(\"${new_state}\")" software-factory-3.0-state-machine.json >/dev/null 2>\&1#g' "$file"

    # Update state extraction patterns
    sed -i 's#grep.*"## Orchestrator States" software-factory-state-machine\.json.*#jq -r ".agents.orchestrator.states | keys[]" software-factory-3.0-state-machine.json#g' "$file"
    sed -i 's#grep.*"## SW Engineer States" software-factory-state-machine\.json.*#jq -r ".agents.\"sw-engineer\".states | keys[]" software-factory-3.0-state-machine.json#g' "$file"
    sed -i 's#grep.*"## Code Reviewer States" software-factory-state-machine\.json.*#jq -r ".agents.\"code-reviewer\".states | keys[]" software-factory-3.0-state-machine.json#g' "$file"
    sed -i 's#grep.*"## Architect States" software-factory-state-machine\.json.*#jq -r ".agents.architect.states | keys[]" software-factory-3.0-state-machine.json#g' "$file"

    # Update orphaned state detection patterns
    sed -i "s#grep '\^\- \\\*\\\*\[A-Z_\]\*\\\*\\\*' software-factory-state-machine\.json#jq -r '.agents | to_entries[] | .value.states | keys[]' software-factory-3.0-state-machine.json#g" "$file"
    sed -i "s#grep -E '→' software-factory-state-machine\.json#jq -r '.transition_matrix | to_entries[] | .value | to_entries[] | .value[]' software-factory-3.0-state-machine.json#g" "$file"

    # Check if file changed
    if ! diff -q "$file" "${file}.bak" >/dev/null 2>&1; then
        local_changes=1
        CHANGES=$((CHANGES + 1))
        echo -e "  ${GREEN}✓${NC} Updated"
        rm -f "${file}.bak"
    else
        echo "  No changes needed"
        rm -f "${file}.bak"
    fi
}

# Update specific rules with more complex patterns
echo "1. Updating R206 - State Machine Transition Validation..."
cat > /tmp/r206_update.md << 'EOF'
## Rule Statement
Agents MUST NEVER attempt to transition to states that are not defined in the software-factory-3.0-state-machine.json. Before ANY state transition, agents MUST read the state machine definition, verify the target state exists for their agent type, and ONLY update current_state if the transition is valid.

## Implementation

### State Validation Using JSON
```bash
validate_state_transition() {
    local AGENT_TYPE="$1"
    local CURRENT_STATE="$2"
    local TARGET_STATE="$3"

    # STEP 1: Read state machine definition
    STATE_MACHINE_FILE="software-factory-3.0-state-machine.json"
    if [ ! -f "$STATE_MACHINE_FILE" ]; then
        echo "❌ FATAL: State machine missing!"
        return 1
    fi

    # STEP 2: Verify target state exists for this agent type
    if ! jq -e ".agents.\"${AGENT_TYPE}\".states.\"${TARGET_STATE}\"" "$STATE_MACHINE_FILE" >/dev/null 2>&1; then
        echo "❌ Invalid state '${TARGET_STATE}' for agent '${AGENT_TYPE}'"
        return 1
    fi

    # STEP 3: Check if transition is valid
    VALID_TRANSITIONS=$(jq -r ".transition_matrix.\"${AGENT_TYPE}\".\"${CURRENT_STATE}\"[]?" "$STATE_MACHINE_FILE" 2>/dev/null)
    if [ -z "$VALID_TRANSITIONS" ]; then
        echo "❌ No valid transitions from '${CURRENT_STATE}'"
        return 1
    fi

    if ! echo "$VALID_TRANSITIONS" | grep -q "^${TARGET_STATE}$"; then
        echo "❌ Invalid transition: ${CURRENT_STATE} -> ${TARGET_STATE}"
        echo "Valid transitions: $(echo $VALID_TRANSITIONS | tr '\n' ', ')"
        return 1
    fi

    echo "✅ Valid transition: ${CURRENT_STATE} -> ${TARGET_STATE}"
    return 0
}
```

### List Valid States for Agent
```bash
list_valid_states() {
    local AGENT_TYPE="$1"

    if [ ! -f "software-factory-3.0-state-machine.json" ]; then
        echo "❌ FATAL: State machine definition missing!";
        return 1
    fi

    echo "Valid states for ${AGENT_TYPE}:"
    jq -r ".agents.\"${AGENT_TYPE}\".states | keys[]" software-factory-3.0-state-machine.json 2>/dev/null | sed 's/^/  /'
}
```
EOF

# Apply R206 update if the file exists
if [ -f "rule-library/R206-state-machine-transition-validation.md" ]; then
    # Keep the rest of the file but update the implementation section
    head -n 6 rule-library/R206-state-machine-transition-validation.md > /tmp/r206_new.md
    cat /tmp/r206_update.md | tail -n +7 >> /tmp/r206_new.md

    # Append the rest of the original file from a safe point
    awk '/## Violations/{p=1} p' rule-library/R206-state-machine-transition-validation.md >> /tmp/r206_new.md

    mv /tmp/r206_new.md rule-library/R206-state-machine-transition-validation.md
    echo -e "  ${GREEN}✓${NC} Updated R206 with JSON implementation"
    CHANGES=$((CHANGES + 1))
fi
echo

# Update all rule files
echo "2. Updating all rule library files..."
for file in rule-library/*.md; do
    update_rule_file "$file"
done
echo

# Update agent-state rules
echo "3. Updating agent-state rules..."
for file in agent-states/*/*/rules.md; do
    if [ -f "$file" ]; then
        update_rule_file "$file"
    fi
done
echo

# Create a reference guide for JSON state machine usage
cat > rule-library/JSON-STATE-MACHINE-REFERENCE.md << 'EOF'
# JSON State Machine Reference Guide

## Overview
The Software Factory state machine is now maintained in JSON format for improved programmatic access and validation.

## Files
- **Main State Machine**: `$CLAUDE_PROJECT_DIR/state-machines/software-factory-3.0-state-machine.json`
- **Sub-State Machines**: `$CLAUDE_PROJECT_DIR/state-machines/[agent].json`
  - orchestrator.json
  - sw-engineer.json
  - code-reviewer.json
  - architect.json

## Usage Examples

### Check if a state exists
```bash
# Using jq
jq -e '.agents.orchestrator.states.INIT' software-factory-3.0-state-machine.json

# In a script
if jq -e ".agents.${AGENT_TYPE}.states.${STATE}" software-factory-3.0-state-machine.json >/dev/null 2>&1; then
    echo "State exists"
fi
```

### Get valid transitions from a state
```bash
# List all valid transitions
jq -r '.transition_matrix.orchestrator.INIT[]' software-factory-3.0-state-machine.json

# Check if specific transition is valid
jq -e '.transition_matrix.orchestrator.INIT | index("WAVE_START")' software-factory-3.0-state-machine.json
```

### Get state metadata
```bash
# Get state type
jq -r '.agents.orchestrator.states.INIT.type' software-factory-3.0-state-machine.json

# Get state description
jq -r '.agents.orchestrator.states.INIT.description' software-factory-3.0-state-machine.json

# Get entry conditions
jq -r '.agents.orchestrator.states.INIT.entry_conditions[]' software-factory-3.0-state-machine.json
```

### List all states for an agent
```bash
jq -r '.agents.orchestrator.states | keys[]' software-factory-3.0-state-machine.json
```

### Count states by type
```bash
jq '.agents.orchestrator.states | to_entries | group_by(.value.type) | map({type: .[0].value.type, count: length})' software-factory-3.0-state-machine.json
```

## Helper Functions
Source the helper functions for easier usage:
```bash
source utilities/state-machine-jq-helpers.sh

# Then use functions like:
state_exists "orchestrator" "INIT"
get_valid_transitions "orchestrator" "INIT"
is_valid_transition "orchestrator" "INIT" "WAVE_START"
```

## Migration from Markdown
All references to `software-factory-3.0-state-machine.json` should be updated to use `software-factory-3.0-state-machine.json` with appropriate jq queries instead of grep patterns.

### Common Replacements
| Old (Markdown/grep) | New (JSON/jq) |
|-------------------|---------------|
| `grep -q "STATE: $STATE" software-factory-3.0-state-machine.json` | `jq -e ".agents.${AGENT_TYPE}.states.${STATE}" software-factory-3.0-state-machine.json` |
| `grep "→" software-factory-3.0-state-machine.json` | `jq -r '.transition_matrix | to_entries[] | "\(.key) -> \(.value[])"' software-factory-3.0-state-machine.json` |
| `grep "## Orchestrator States" -A 20` | `jq -r '.agents.orchestrator.states | keys[]' software-factory-3.0-state-machine.json` |
EOF

echo -e "${GREEN}✓${NC} Created JSON State Machine Reference Guide"
CHANGES=$((CHANGES + 1))

echo
echo "========================================="
echo -e "${GREEN}UPDATE COMPLETE${NC}"
echo "========================================="
echo "Total changes made: $CHANGES"
echo
echo "All rules have been updated to use the JSON state machine format."
echo "Reference guide created at: rule-library/JSON-STATE-MACHINE-REFERENCE.md"