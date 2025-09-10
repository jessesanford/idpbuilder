#!/bin/bash

# validate-state.sh - Validates orchestrator state file against JSON schema
# Usage: tools/validate-state.sh [path/to/orchestrator-state.json]
# Returns: 0 for valid, 1 for invalid

set -euo pipefail

# Color codes for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Get the script directory (tools/)
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

# Function to find schema file in multiple locations
find_schema_file() {
    local locations=(
        # 1. First check in current project directory
        "$PROJECT_ROOT/orchestrator-state.schema.json"
        # 2. Check in SOFTWARE_FACTORY_DIR if set
        "${SOFTWARE_FACTORY_DIR:-}/orchestrator-state.schema.json"
        # 3. Check in CLAUDE_PROJECT_DIR if set
        "${CLAUDE_PROJECT_DIR:-}/orchestrator-state.schema.json"
        # 4. Check in known software-factory-template location
        "/home/vscode/software-factory-template/orchestrator-state.schema.json"
        # 5. Check relative to workspaces
        "/home/vscode/workspaces/software-factory-2.0-template/orchestrator-state.schema.json"
        # 6. Check parent directory (for when run from subdirs)
        "$(dirname "$PROJECT_ROOT")/orchestrator-state.schema.json"
    )
    
    for location in "${locations[@]}"; do
        if [ -n "$location" ] && [ -f "$location" ]; then
            echo "$location"
            return 0
        fi
    done
    
    return 1
}

# Find schema file
SCHEMA_FILE=$(find_schema_file || echo "")

# Default state file location
DEFAULT_STATE_FILE="$PROJECT_ROOT/orchestrator-state.json"

# Get state file from argument or use default
STATE_FILE="${1:-$DEFAULT_STATE_FILE}"

# Check if schema file exists
if [ -z "$SCHEMA_FILE" ] || [ ! -f "$SCHEMA_FILE" ]; then
    echo -e "${RED}❌ Schema file not found in any of these locations:${NC}"
    echo "   - $PROJECT_ROOT/orchestrator-state.schema.json"
    echo "   - \${SOFTWARE_FACTORY_DIR}/orchestrator-state.schema.json"
    echo "   - \${CLAUDE_PROJECT_DIR}/orchestrator-state.schema.json"
    echo "   - /home/vscode/software-factory-template/orchestrator-state.schema.json"
    echo ""
    echo -e "${YELLOW}💡 To fix this:${NC}"
    echo "   1. Set SOFTWARE_FACTORY_DIR or CLAUDE_PROJECT_DIR environment variable"
    echo "   2. OR copy the schema file to your project:"
    echo "      cp /home/vscode/software-factory-template/orchestrator-state.schema.json $PROJECT_ROOT/"
    exit 1
fi

# Check if state file exists
if [ ! -f "$STATE_FILE" ]; then
    echo -e "${RED}❌ State file not found: $STATE_FILE${NC}"
    exit 1
fi

# Function to validate using Python jsonschema
validate_with_python() {
    python3 -c "
import json
import sys
from jsonschema import validate, ValidationError, Draft7Validator

try:
    # Load schema
    with open('$SCHEMA_FILE', 'r') as f:
        schema = json.load(f)
    
    # Load state file
    with open('$STATE_FILE', 'r') as f:
        state = json.load(f)
    
    # Create validator for better error messages
    validator = Draft7Validator(schema)
    
    # Check if valid
    errors = sorted(validator.iter_errors(state), key=lambda e: e.path)
    
    if errors:
        print('❌ State file validation failed!')
        print('')
        for error in errors:
            if error.path:
                path = '.'.join(str(p) for p in error.path)
                print(f'  Field: {path}')
            else:
                print('  Field: (root)')
            print(f'  Error: {error.message}')
            if error.validator == 'enum' and hasattr(error, 'validator_value'):
                print(f'  Valid values: {error.validator_value[:5]}...' if len(error.validator_value) > 5 else f'  Valid values: {error.validator_value}')
            print('')
        sys.exit(1)
    else:
        print('✅ State file is valid!')
        sys.exit(0)
        
except json.JSONDecodeError as e:
    print(f'❌ Invalid JSON in state file: {e}')
    sys.exit(1)
except Exception as e:
    print(f'❌ Validation error: {e}')
    sys.exit(1)
"
    return $?
}

# Function to validate using jq (basic validation)
validate_with_jq() {
    # First check if file is valid JSON
    if ! jq empty "$STATE_FILE" 2>/dev/null; then
        echo -e "${RED}❌ State file is not valid JSON!${NC}"
        jq empty "$STATE_FILE" 2>&1 | sed 's/^/  /'
        return 1
    fi
    
    # Check required fields
    local missing_fields=()
    local required_fields=(
        "current_phase"
        "current_wave"
        "current_state"
        "previous_state"
        "transition_time"
        "phases_planned"
        "waves_per_phase"
        "efforts_completed"
        "efforts_in_progress"
        "efforts_pending"
        "project_info"
    )
    
    for field in "${required_fields[@]}"; do
        if ! jq -e "has(\"$field\")" "$STATE_FILE" >/dev/null 2>&1; then
            missing_fields+=("$field")
        fi
    done
    
    if [ ${#missing_fields[@]} -gt 0 ]; then
        echo -e "${RED}❌ State file validation failed!${NC}"
        echo -e "${RED}Missing required fields:${NC}"
        for field in "${missing_fields[@]}"; do
            echo "  - $field"
        done
        return 1
    fi
    
    # Check current_state is not empty
    local current_state=$(jq -r '.current_state' "$STATE_FILE")
    if [ -z "$current_state" ] || [ "$current_state" = "null" ]; then
        echo -e "${RED}❌ State file validation failed!${NC}"
        echo -e "${RED}Field 'current_state' is empty or null${NC}"
        return 1
    fi
    
    # Check if current_state is in allowed list (simplified check)
    local valid_states=(
        "INIT" "PLANNING" "SETUP_EFFORT_INFRASTRUCTURE"
        "ANALYZE_CODE_REVIEWER_PARALLELIZATION" "SPAWN_CODE_REVIEWERS_EFFORT_PLANNING"
        "WAITING_FOR_EFFORT_PLANS" "ANALYZE_IMPLEMENTATION_PARALLELIZATION"
        "SPAWN_AGENTS" "MONITOR" "MONITOR_IMPLEMENTATION" "MONITOR_REVIEWS"
        "WAVE_START" "WAVE_COMPLETE" "WAVE_REVIEW" "INTEGRATION"
        "PHASE_INTEGRATION" "PROJECT_INTEGRATION" "PHASE_COMPLETE"
        "ERROR_RECOVERY" "SUCCESS" "HARD_STOP"
    )
    
    local state_valid=false
    for valid_state in "${valid_states[@]}"; do
        if [ "$current_state" = "$valid_state" ]; then
            state_valid=true
            break
        fi
    done
    
    if [ "$state_valid" = false ]; then
        echo -e "${YELLOW}⚠️  Warning: current_state '$current_state' may not be in state machine${NC}"
        echo -e "${YELLOW}   (This is a simplified check - use Python validation for full verification)${NC}"
    fi
    
    echo -e "${GREEN}✅ State file passed basic validation!${NC}"
    echo -e "${YELLOW}   Note: For complete schema validation, Python jsonschema is recommended${NC}"
    return 0
}

# Main validation logic
echo "🔍 Validating orchestrator state file..."
echo "   State file: $STATE_FILE"
echo "   Schema: $SCHEMA_FILE"
echo ""

# Try Python validation first (most comprehensive)
if command -v python3 >/dev/null 2>&1; then
    # Check if jsonschema is installed
    if python3 -c "import jsonschema" 2>/dev/null; then
        validate_with_python
        exit $?
    else
        echo -e "${YELLOW}⚠️  Python jsonschema module not found${NC}"
        echo "   Install with: pip3 install jsonschema"
        echo "   Falling back to basic validation..."
        echo ""
    fi
fi

# Fall back to jq validation
if command -v jq >/dev/null 2>&1; then
    validate_with_jq
    exit $?
else
    echo -e "${RED}❌ Neither Python jsonschema nor jq is available${NC}"
    echo "   Install one of:"
    echo "     - Python: pip3 install jsonschema"
    echo "     - jq: apt-get install jq (or brew install jq on macOS)"
    exit 1
fi