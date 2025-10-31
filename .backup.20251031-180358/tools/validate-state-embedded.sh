#!/bin/bash

# validate-state-embedded.sh - Self-contained orchestrator state validator
# This version has the schema embedded directly in the script for maximum portability
# Usage: tools/validate-state-embedded.sh [path/to/orchestrator-state-v3.json]

set -euo pipefail

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Get state file from argument or use default
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
DEFAULT_STATE_FILE="$PROJECT_ROOT/orchestrator-state-v3.json"
STATE_FILE="${1:-$DEFAULT_STATE_FILE}"

# Check if state file exists
if [ ! -f "$STATE_FILE" ]; then
    echo -e "${RED}❌ State file not found: $STATE_FILE${NC}"
    exit 1
fi

# Create temporary schema file
TEMP_SCHEMA=$(mktemp /tmp/orchestrator-state-schema.XXXXXX.json)
trap "rm -f $TEMP_SCHEMA" EXIT

# Embedded schema (SF 3.0 - extracted from schemas/orchestrator-state-v3.schema.json)
cat "$PROJECT_ROOT/schemas/orchestrator-state-v3.schema.json" > "$TEMP_SCHEMA"

# Function to validate using Python jsonschema
validate_with_python() {
    python3 << PYEOF
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
PYEOF
    return $?
}

# Main validation logic
echo "🔍 Validating orchestrator state file (embedded schema)..."
echo "   State file: $STATE_FILE"
echo "   Schema: SF 3.0 (from schemas/orchestrator-state-v3.schema.json)"
echo ""

# Try Python validation (most comprehensive)
if command -v python3 >/dev/null 2>&1; then
    # Check if jsonschema is installed
    if python3 -c "import jsonschema" 2>/dev/null; then
        validate_with_python
        exit $?
    else
        echo -e "${YELLOW}⚠️  Python jsonschema module not found${NC}"
        echo "   Install with: pip3 install jsonschema"
        echo ""
        echo "Note: This SF 3.0 validator requires Python jsonschema for full validation."
        echo "      For simpler validation, use tools/validate-state-file.sh instead."
        exit 1
    fi
else
    echo -e "${RED}❌ Python 3 is required for SF 3.0 state validation${NC}"
    echo "   Install Python 3, then: pip3 install jsonschema"
    exit 1
fi
