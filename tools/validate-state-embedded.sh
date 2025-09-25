#!/bin/bash

# validate-state-embedded.sh - Self-contained orchestrator state validator
# This version has the schema embedded directly in the script for maximum portability
# Usage: tools/validate-state-embedded.sh [path/to/orchestrator-state.json]

set -euo pipefail

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Get state file from argument or use default
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
DEFAULT_STATE_FILE="$PROJECT_ROOT/orchestrator-state.json"
STATE_FILE="${1:-$DEFAULT_STATE_FILE}"

# Check if state file exists
if [ ! -f "$STATE_FILE" ]; then
    echo -e "${RED}❌ State file not found: $STATE_FILE${NC}"
    exit 1
fi

# Create temporary schema file
TEMP_SCHEMA=$(mktemp /tmp/orchestrator-state-schema.XXXXXX.json)
trap "rm -f $TEMP_SCHEMA" EXIT

# Embedded schema (automatically updated from orchestrator-state.schema.json)
cat > "$TEMP_SCHEMA" << 'EMBEDDED_SCHEMA'
{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "title": "Orchestrator State",
  "description": "State tracking for Software Factory 2.0 orchestration",
  "type": "object",
  "required": [
    "current_phase",
    "current_wave", 
    "current_state",
    "previous_state",
    "transition_time",
    "phases_planned",
    "waves_per_phase",
    "efforts_completed",
    "efforts_in_progress",
    "efforts_pending",
    "project_info"
  ],
  "properties": {
    "current_phase": {
      "type": "integer",
      "minimum": 0,
      "description": "Current phase number (0-based)"
    },
    "current_wave": {
      "type": "integer",
      "minimum": 0,
      "description": "Current wave within phase (0-based)"
    },
    "current_state": {
      "type": "string",
      "enum": [
        "INIT",
        "PLANNING",
        "SETUP_EFFORT_INFRASTRUCTURE",
        "ANALYZE_CODE_REVIEWER_PARALLELIZATION",
        "SPAWN_CODE_REVIEWERS_EFFORT_PLANNING",
        "WAITING_FOR_EFFORT_PLANS",
        "ANALYZE_IMPLEMENTATION_PARALLELIZATION",
        "SPAWN_AGENTS",
        "MONITOR",
        "MONITOR_IMPLEMENTATION",
        "MONITOR_REVIEWS",
        "WAVE_START",
        "WAVE_COMPLETE",
        "WAVE_REVIEW",
        "INTEGRATION",
        "PHASE_INTEGRATION",
        "PROJECT_INTEGRATION",
        "PHASE_COMPLETE",
        "ERROR_RECOVERY",
        "SUCCESS",
        "HARD_STOP"
      ],
      "description": "Current state in the state machine"
    },
    "previous_state": {
      "type": ["string", "null"],
      "description": "Previous state before last transition"
    },
    "transition_time": {
      "type": "string",
      "format": "date-time",
      "description": "ISO 8601 timestamp of last state transition"
    },
    "phases_planned": {
      "type": "integer",
      "minimum": 1,
      "description": "Total number of phases planned"
    },
    "waves_per_phase": {
      "type": "object",
      "additionalProperties": {
        "type": "integer",
        "minimum": 1
      },
      "description": "Number of waves planned per phase"
    },
    "efforts_completed": {
      "type": "array",
      "items": {
        "type": "object",
        "required": ["effort_id", "phase", "wave"],
        "properties": {
          "effort_id": {"type": "string"},
          "phase": {"type": "integer"},
          "wave": {"type": "integer"},
          "branch": {"type": "string"},
          "completion_time": {"type": "string", "format": "date-time"},
          "pr_url": {"type": "string"},
          "review_status": {"type": "string"}
        }
      },
      "description": "List of completed efforts"
    },
    "efforts_in_progress": {
      "type": "array",
      "items": {
        "type": "object",
        "required": ["effort_id", "phase", "wave", "agent_id"],
        "properties": {
          "effort_id": {"type": "string"},
          "phase": {"type": "integer"},
          "wave": {"type": "integer"},
          "agent_id": {"type": "string"},
          "agent_type": {"type": "string"},
          "working_copy": {"type": "string"},
          "branch": {"type": "string"},
          "start_time": {"type": "string", "format": "date-time"},
          "last_update": {"type": "string", "format": "date-time"},
          "status": {"type": "string"}
        }
      },
      "description": "List of efforts currently in progress"
    },
    "efforts_pending": {
      "type": "array",
      "items": {
        "type": "object",
        "required": ["effort_id", "phase", "wave"],
        "properties": {
          "effort_id": {"type": "string"},
          "phase": {"type": "integer"},
          "wave": {"type": "integer"},
          "dependencies": {
            "type": "array",
            "items": {"type": "string"}
          }
        }
      },
      "description": "List of pending efforts"
    },
    "project_info": {
      "type": "object",
      "required": ["name", "repository"],
      "properties": {
        "name": {"type": "string"},
        "repository": {"type": "string"},
        "description": {"type": "string"},
        "main_branch": {"type": "string"},
        "integration_branch": {"type": "string"}
      },
      "description": "Project metadata"
    },
    "phase_branches": {
      "type": "object",
      "additionalProperties": {"type": "string"},
      "description": "Integration branches per phase"
    },
    "wave_branches": {
      "type": "object",
      "additionalProperties": {"type": "string"},
      "description": "Integration branches per wave"
    },
    "error_log": {
      "type": "array",
      "items": {
        "type": "object",
        "properties": {
          "timestamp": {"type": "string", "format": "date-time"},
          "error": {"type": "string"},
          "recovery_action": {"type": "string"}
        }
      },
      "description": "Log of errors and recovery actions"
    },
    "metrics": {
      "type": "object",
      "properties": {
        "total_efforts": {"type": "integer"},
        "efforts_completed": {"type": "integer"},
        "efforts_failed": {"type": "integer"},
        "average_effort_duration": {"type": "number"},
        "total_lines_changed": {"type": "integer"}
      },
      "description": "Project metrics"
    },
    "last_checkpoint": {
      "type": "string",
      "format": "date-time",
      "description": "Last checkpoint timestamp"
    },
    "notes": {
      "type": "string",
      "description": "Free-form notes about current state"
    }
  }
}
EMBEDDED_SCHEMA

# Function to validate using Python jsonschema
validate_with_python() {
    python3 -c "
import json
import sys
from jsonschema import validate, ValidationError, Draft7Validator

try:
    # Load schema
    with open('$TEMP_SCHEMA', 'r') as f:
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
    
    # Check if current_state is in allowed list
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
    echo -e "${YELLOW}   Note: This is using embedded schema (self-contained version)${NC}"
    return 0
}

# Main validation logic
echo "🔍 Validating orchestrator state file (embedded schema)..."
echo "   State file: $STATE_FILE"
echo "   Schema: Embedded in script"
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