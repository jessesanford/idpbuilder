#!/bin/bash
# SF 3.0 JSON Schema Validation Script
# Validates state files against their corresponding JSON schemas
#
# Usage: ./validate-state-file.sh <state-file>
#
# Supported file types:
#   - orchestrator-state-v3.json (or .example)
#   - bug-tracking.json (or .example)
#   - integration-containers.json (or .example)
#   - fix-cascade-state.json (or .example)
#
# Exit codes:
#   0 = validation successful
#   1 = validation failed
#   2 = invalid usage or file not found

set -euo pipefail

# Color codes for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Script directory (to find schemas)
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
SCHEMAS_DIR="$PROJECT_ROOT/schemas"

# Usage function
usage() {
    echo "Usage: $0 <state-file>"
    echo ""
    echo "Validates a state file against its JSON schema."
    echo ""
    echo "Supported file types:"
    echo "  - orchestrator-state-v3.json[.example]"
    echo "  - bug-tracking.json[.example]"
    echo "  - integration-containers.json[.example]"
    echo "  - fix-cascade-state.json[.example]"
    echo ""
    echo "Example:"
    echo "  $0 orchestrator-state-v3.json.example"
    exit 2
}

# Check arguments
if [ $# -ne 1 ]; then
    echo -e "${RED}ERROR: Missing state file argument${NC}" >&2
    usage
fi

STATE_FILE="$1"

# Check if file exists
if [ ! -f "$STATE_FILE" ]; then
    echo -e "${RED}ERROR: File not found: $STATE_FILE${NC}" >&2
    exit 2
fi

# Extract basename for pattern matching
FILENAME="$(basename "$STATE_FILE")"

# Detect file type and determine schema
SCHEMA_FILE=""
FILE_TYPE=""

case "$FILENAME" in
    orchestrator-state-v3.json*)
        # SF 3.0 ONLY - No more hybrid/v2 support
        # All orchestrator-state-v3.json files MUST be pure SF 3.0 format
        # Required structure: project_progression, state_machine, references
        SCHEMA_FILE="$SCHEMAS_DIR/orchestrator-state-v3.schema.json"
        FILE_TYPE="orchestrator-state-v3 (SF 3.0)"
        ;;
    bug-tracking.json*)
        SCHEMA_FILE="$SCHEMAS_DIR/bug-tracking.schema.json"
        FILE_TYPE="bug-tracking"
        ;;
    integration-containers.json*)
        SCHEMA_FILE="$SCHEMAS_DIR/integration-containers.schema.json"
        FILE_TYPE="integration-containers"
        ;;
    fix-cascade-state.json*)
        SCHEMA_FILE="$SCHEMAS_DIR/fix-cascade-state.schema.json"
        FILE_TYPE="fix-cascade-state"
        ;;
    *)
        echo -e "${RED}ERROR: Unknown state file type: $FILENAME${NC}" >&2
        echo -e "${YELLOW}Supported patterns:${NC}" >&2
        echo "  - orchestrator-state-v3.json*" >&2
        echo "  - bug-tracking.json*" >&2
        echo "  - integration-containers.json*" >&2
        echo "  - fix-cascade-state.json*" >&2
        exit 2
        ;;
esac

# Check if schema exists
if [ ! -f "$SCHEMA_FILE" ]; then
    echo -e "${RED}ERROR: Schema file not found: $SCHEMA_FILE${NC}" >&2
    exit 2
fi

# Validate using Python's jsonschema library
echo "Validating $FILE_TYPE state file..."
echo "  State file: $STATE_FILE"
echo "  Schema: $SCHEMA_FILE"
echo ""

# Python validation script
VALIDATION_RESULT=$(python3 -c "
import json
import sys
from jsonschema import validate, ValidationError, SchemaError

try:
    # Load state file
    with open('$STATE_FILE', 'r') as f:
        state_data = json.load(f)

    # Load schema
    with open('$SCHEMA_FILE', 'r') as f:
        schema = json.load(f)

    # Validate
    validate(instance=state_data, schema=schema)

    print('✅ VALIDATION PASSED')
    sys.exit(0)

except json.JSONDecodeError as e:
    print(f'❌ JSON SYNTAX ERROR: {e}', file=sys.stderr)
    sys.exit(1)

except ValidationError as e:
    print(f'❌ VALIDATION FAILED:', file=sys.stderr)
    print(f'   Path: {\" -> \".join(str(p) for p in e.absolute_path)}', file=sys.stderr)
    print(f'   Message: {e.message}', file=sys.stderr)
    if e.schema_path:
        print(f'   Schema path: {\" -> \".join(str(p) for p in e.schema_path)}', file=sys.stderr)
    sys.exit(1)

except SchemaError as e:
    print(f'❌ SCHEMA ERROR: {e}', file=sys.stderr)
    sys.exit(1)

except Exception as e:
    print(f'❌ ERROR: {e}', file=sys.stderr)
    sys.exit(1)
" 2>&1)

VALIDATION_EXIT_CODE=$?

# Print validation result
echo "$VALIDATION_RESULT"

# Exit with appropriate code
exit $VALIDATION_EXIT_CODE
