#!/bin/bash

# Software Factory 2.0 - Valid States Synchronization Utility
# This script synchronizes valid states from the state machine to the orchestrator-state-v3.json schema
# ensuring that the schema always enforces only valid states

set -e

# Color codes for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
BOLD='\033[1m'
NC='\033[0m' # No Color

# Function to print colored messages
print_error() {
    echo -e "${RED}${BOLD}❌ ERROR:${NC} ${RED}$1${NC}" >&2
}

print_success() {
    echo -e "${GREEN}${BOLD}✅ PROJECT_DONE:${NC} ${GREEN}$1${NC}"
}

print_info() {
    echo -e "${BLUE}${BOLD}ℹ️  INFO:${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}${BOLD}⚠️  WARNING:${NC} ${YELLOW}$1${NC}"
}

# Get script directory and project root
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
cd "$PROJECT_ROOT"

# Paths
STATE_MACHINE="state-machines/software-factory-3.0-state-machine.json"
SCHEMA_FILE="orchestrator-state.schema.json"
TEMP_STATES="/tmp/valid-states.json"
TEMP_SCHEMA="/tmp/updated-schema.json"

echo -e "${BOLD}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${BOLD}  Software Factory 2.0 - State Synchronization${NC}"
echo -e "${BOLD}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""

# Check if state machine exists
if [ ! -f "$STATE_MACHINE" ]; then
    print_error "State machine file not found at $STATE_MACHINE"
    exit 1
fi

# Check if schema file exists
if [ ! -f "$SCHEMA_FILE" ]; then
    print_warning "Schema file not found at $SCHEMA_FILE - creating new one"
fi

print_info "Extracting valid states from state machine..."

# Extract valid states from state machine using Python
python3 << 'EOF' > "$TEMP_STATES"
import json
import sys

try:
    with open('state-machines/software-factory-3.0-state-machine.json') as f:
        data = json.load(f)

    all_states = {}

    # Extract states for each agent from transition_matrix
    for agent in ['orchestrator', 'sw-engineer', 'code-reviewer', 'architect']:
        if 'transition_matrix' in data and agent in data['transition_matrix']:
            states = list(data['transition_matrix'][agent].keys())
            all_states[agent] = sorted(states)

    # Output valid states as JSON
    print(json.dumps(all_states, indent=2))

except Exception as e:
    print(f"Error extracting states: {e}", file=sys.stderr)
    sys.exit(1)
EOF

if [ $? -ne 0 ]; then
    print_error "Failed to extract valid states from state machine"
    rm -f "$TEMP_STATES"
    exit 1
fi

# Show statistics
ORCHESTRATOR_COUNT=$(jq '.orchestrator | length' "$TEMP_STATES")
SW_ENGINEER_COUNT=$(jq '."sw-engineer" | length' "$TEMP_STATES")
CODE_REVIEWER_COUNT=$(jq '."code-reviewer" | length' "$TEMP_STATES")
ARCHITECT_COUNT=$(jq '.architect | length' "$TEMP_STATES")

echo ""
print_info "States extracted successfully:"
echo "    • Orchestrator:  $ORCHESTRATOR_COUNT states"
echo "    • SW-Engineer:   $SW_ENGINEER_COUNT states"
echo "    • Code-Reviewer: $CODE_REVIEWER_COUNT states"
echo "    • Architect:     $ARCHITECT_COUNT states"
echo ""

print_info "Updating orchestrator-state.schema.json..."

# Update the schema file with new valid states
python3 << 'EOF'
import json
import sys
from datetime import datetime

try:
    # Load valid states
    with open('/tmp/valid-states.json') as f:
        valid_states = json.load(f)

    # Load existing schema or create new one
    try:
        with open('orchestrator-state.schema.json') as f:
            schema = json.load(f)
    except FileNotFoundError:
        # Create minimal schema if it doesn't exist
        schema = {
            "$schema": "http://json-schema.org/draft-07/schema#",
            "title": "Software Factory 2.0 - Orchestrator State Schema",
            "description": "Schema for orchestrator-state-v3.json with valid state enforcement",
            "type": "object",
            "required": ["current_phase", "current_wave", "current_state"],
            "properties": {}
        }

    # Update current_state property with new enum
    if 'properties' not in schema:
        schema['properties'] = {}

    orchestrator_states = valid_states.get('orchestrator', [])

    schema['properties']['current_state'] = {
        "type": "string",
        "description": "Current state of the orchestrator - MUST be a valid state from state machine",
        "enum": orchestrator_states
    }

    # Update metadata
    if '_metadata' not in schema:
        schema['_metadata'] = {}

    schema['_metadata'].update({
        "last_updated": datetime.now().strftime("%Y-%m-%d %H:%M:%S"),
        "auto_generated": True,
        "source": "state-machines/software-factory-3.0-state-machine.json",
        "orchestrator_states_count": len(orchestrator_states),
        "note": "This schema is auto-generated from the state machine. Run utilities/sync-valid-states.sh to update."
    })

    # Write updated schema
    with open('/tmp/updated-schema.json', 'w') as f:
        json.dump(schema, f, indent=2)

    print(f"Schema updated with {len(orchestrator_states)} valid orchestrator states")

except Exception as e:
    print(f"Error updating schema: {e}", file=sys.stderr)
    sys.exit(1)
EOF

if [ $? -ne 0 ]; then
    print_error "Failed to update schema"
    rm -f "$TEMP_STATES" "$TEMP_SCHEMA"
    exit 1
fi

# Backup existing schema if it exists
if [ -f "$SCHEMA_FILE" ]; then
    BACKUP_FILE="${SCHEMA_FILE}.backup.$(date +%Y%m%d-%H%M%S)"
    cp "$SCHEMA_FILE" "$BACKUP_FILE"
    print_info "Backed up existing schema to $BACKUP_FILE"
fi

# Move updated schema to final location
mv "$TEMP_SCHEMA" "$SCHEMA_FILE"

# Clean up temp files
rm -f "$TEMP_STATES"

print_success "Schema successfully updated with valid states!"
echo ""

# Validate current orchestrator-state-v3.json against new schema
if [ -f "orchestrator-state-v3.json" ]; then
    print_info "Validating current orchestrator-state-v3.json against updated schema..."

    # Extract current state
    CURRENT_STATE=$(jq -r '.current_state' orchestrator-state-v3.json 2>/dev/null || echo "")

    if [ -n "$CURRENT_STATE" ]; then
        # Check if current state is in the valid states in the schema
        if jq -e --arg state "$CURRENT_STATE" '.properties.current_state.enum | index($state)' "$SCHEMA_FILE" > /dev/null 2>&1; then
            print_success "Current state '$CURRENT_STATE' is valid!"
        else
            print_error "Current state '$CURRENT_STATE' is NOT valid!"
            print_warning "You need to update orchestrator-state-v3.json with a valid state"
            echo ""
            echo "Valid states:"
            jq -r '.properties.current_state.enum[]' "$SCHEMA_FILE" | head -20
            echo "... and more. Run: jq '.properties.current_state.enum[]' $SCHEMA_FILE"
        fi
    fi
fi

echo ""
echo -e "${BOLD}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
print_success "State synchronization complete!"
echo ""
echo "Next steps:"
echo "  1. Review the updated schema: $SCHEMA_FILE"
echo "  2. Commit the changes: git add $SCHEMA_FILE && git commit -m 'sync: update schema with valid states from state machine'"
echo "  3. The pre-commit hook will now validate against these states"
echo ""