#!/bin/bash

# enforce-state-validation.sh - Comprehensive orchestrator state validation with effort metadata checks
# Usage: tools/enforce-state-validation.sh [--fix] [path/to/orchestrator-state-v3.json]
# Returns: 0 for valid, 1 for invalid
# The --fix option attempts to auto-fix certain issues (like missing timestamps)

set -euo pipefail

# Color codes for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
MAGENTA='\033[0;35m'
NC='\033[0m' # No Color

# Get the script directory (tools/)
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
CLAUDE_PROJECT_DIR="${CLAUDE_PROJECT_DIR:-$PROJECT_ROOT}"

# Check for --fix option
FIX_MODE=false
if [[ "${1:-}" == "--fix" ]]; then
    FIX_MODE=true
    shift
fi

# Default state file location
DEFAULT_STATE_FILE="$PROJECT_ROOT/orchestrator-state-v3.json"

# Get state file from argument or use default
STATE_FILE="${1:-$DEFAULT_STATE_FILE}"

# Track validation failures
VALIDATION_FAILED=false
CRITICAL_FAILURES=0
WARNINGS=0

# Print section header
print_section() {
    echo ""
    echo -e "${BLUE}$1${NC}"
    echo "────────────────────────────────────────────────────────────"
}

# Print error
print_error() {
    echo -e "${RED}❌ $1${NC}"
    VALIDATION_FAILED=true
    ((CRITICAL_FAILURES++))
}

# Print warning
print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
    ((WARNINGS++))
}

# Print success
print_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

# Print info
print_info() {
    echo -e "${MAGENTA}ℹ️  $1${NC}"
}

# Check if state file exists
if [ ! -f "$STATE_FILE" ]; then
    print_error "State file not found: $STATE_FILE"
    exit 1
fi

# Load target repository configuration
TARGET_REPO_CONFIG="$CLAUDE_PROJECT_DIR/target-repo-config.yaml"
if [ ! -f "$TARGET_REPO_CONFIG" ]; then
    print_error "target-repo-config.yaml not found at: $TARGET_REPO_CONFIG"
    exit 1
fi

# Extract target repository URL
TARGET_REPO=$(grep "url:" "$TARGET_REPO_CONFIG" | head -1 | sed 's/.*url:[[:space:]]*"//' | sed 's/".*//')
if [ -z "$TARGET_REPO" ]; then
    print_error "Could not extract repository URL from target-repo-config.yaml"
    exit 1
fi

# Extract project prefix from state file (if exists)
PROJECT_PREFIX=$(jq -r '.project_info.project_prefix // empty' "$STATE_FILE" 2>/dev/null || echo "")

# If not in state file, check branch_naming section in config
if [ -z "$PROJECT_PREFIX" ]; then
    PROJECT_PREFIX=$(grep "project_prefix:" "$TARGET_REPO_CONFIG" 2>/dev/null | sed 's/.*project_prefix:[[:space:]]*"//' | sed 's/".*//' | sed 's/#.*//' | xargs || echo "")
fi

echo "🔍 COMPREHENSIVE ORCHESTRATOR STATE VALIDATION"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "State file: $STATE_FILE"
echo "Fix mode: $FIX_MODE"
echo "Target repository: $TARGET_REPO"
echo "Project prefix: ${PROJECT_PREFIX:-<none>}"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

# 1. JSON VALIDITY CHECK
print_section "1. JSON VALIDITY CHECK"
if jq empty "$STATE_FILE" 2>/dev/null; then
    print_success "State file is valid JSON"
else
    print_error "State file is not valid JSON!"
    jq empty "$STATE_FILE" 2>&1 | sed 's/^/  /'
    exit 1
fi

# 2. REQUIRED FIELDS CHECK
print_section "2. REQUIRED FIELDS CHECK"
REQUIRED_FIELDS=(
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

for field in "${REQUIRED_FIELDS[@]}"; do
    if jq -e "has(\"$field\")" "$STATE_FILE" >/dev/null 2>&1; then
        print_success "Required field '$field' present"
    else
        print_error "Missing required field: $field"
    fi
done

# 3. STATE MACHINE VALIDATION
print_section "3. STATE MACHINE VALIDATION"
CURRENT_STATE=$(jq -r '.current_state' "$STATE_FILE")
STATE_MACHINE_JSON="$CLAUDE_PROJECT_DIR/state-machines/software-factory-3.0-state-machine.json"
STATE_MACHINE_MD="$CLAUDE_PROJECT_DIR/state-machines/software-factory-3.0-state-machine.json"

# Try JSON file first (more reliable)
if [ -f "$STATE_MACHINE_JSON" ]; then
    # Get the agent type from the state file (default to orchestrator)
    AGENT_TYPE="orchestrator"

    # Extract valid states from JSON for the specific agent
    VALID_STATES=$(jq -r ".transition_matrix.${AGENT_TYPE} | keys[]" "$STATE_MACHINE_JSON" 2>/dev/null)

    if [ -z "$VALID_STATES" ]; then
        # Try getting all orchestrator states if the above fails
        VALID_STATES=$(jq -r '.transition_matrix.orchestrator | keys[]' "$STATE_MACHINE_JSON" 2>/dev/null)
    fi

    if echo "$VALID_STATES" | grep -q "^${CURRENT_STATE}$"; then
        print_success "Current state '$CURRENT_STATE' is valid in state machine"
    else
        # Also check other agents' states in case it's a cross-agent state
        ALL_STATES=$(jq -r '.transition_matrix | .[] | keys[]' "$STATE_MACHINE_JSON" 2>/dev/null | sort -u)
        if echo "$ALL_STATES" | grep -q "^${CURRENT_STATE}$"; then
            print_success "Current state '$CURRENT_STATE' is valid in state machine"
        else
            print_error "Current state '$CURRENT_STATE' not found in state machine!"
            echo "  Valid orchestrator states include:"
            echo "$VALID_STATES" | grep -i "monitor\|spawn\|init" | head -10 | sed 's/^/    - /'
            echo "    ... (showing sample states)"
        fi
    fi
elif [ -f "$STATE_MACHINE_MD" ]; then
    # Fallback to markdown file
    # Extract valid states from the state machine file (look for STATE: patterns)
    VALID_STATES=$(grep -E "^STATE:" "$STATE_MACHINE_MD" | sed 's/STATE:[[:space:]]*//' | sort -u)

    if echo "$VALID_STATES" | grep -q "^${CURRENT_STATE}$"; then
        print_success "Current state '$CURRENT_STATE' is valid in state machine"
    else
        print_error "Current state '$CURRENT_STATE' not found in state machine!"
        echo "  Valid states are:"
        echo "$VALID_STATES" | head -10 | sed 's/^/    - /'
        echo "    ... (showing first 10)"
    fi
else
    print_warning "State machine files not found, skipping state validation"
fi

# 4. SCHEMA VALIDATION (if Python available)
print_section "4. SCHEMA VALIDATION (SF 3.0)"
SCHEMA_FILE="$PROJECT_ROOT/schemas/orchestrator-state-v3.schema.json"

# Ensure jsonschema is available through venv
VENV_DIR="$HOME/.software-factory-venv"
VENV_PYTHON=""

if [ -f "$SCHEMA_FILE" ] && command -v python3 >/dev/null 2>&1; then
    # Check if jsonschema is available in system Python
    if ! python3 -c "import jsonschema" 2>/dev/null; then
        # Try to use or create venv with jsonschema
        if [ -d "$VENV_DIR" ]; then
            # Venv exists, activate it
            if [ -f "$VENV_DIR/bin/python3" ]; then
                VENV_PYTHON="$VENV_DIR/bin/python3"
                # Check if jsonschema is installed in venv
                if ! "$VENV_PYTHON" -c "import jsonschema" 2>/dev/null; then
                    print_info "Installing jsonschema in existing venv..."
                    "$VENV_DIR/bin/pip" install jsonschema >/dev/null 2>&1 || {
                        print_warning "Failed to install jsonschema in venv, proceeding without schema validation"
                        VENV_PYTHON=""
                    }
                fi
            else
                print_warning "Venv directory exists but no python3 found, attempting recreation..."
                rm -rf "$VENV_DIR"
                python3 -m venv "$VENV_DIR" && \
                "$VENV_DIR/bin/pip" install jsonschema >/dev/null 2>&1 && \
                VENV_PYTHON="$VENV_DIR/bin/python3" || {
                    print_warning "Failed to recreate venv, proceeding without schema validation"
                }
            fi
        else
            # Create new venv with jsonschema
            print_info "Creating Software Factory venv with jsonschema..."
            if python3 -m venv "$VENV_DIR" 2>/dev/null; then
                if "$VENV_DIR/bin/pip" install jsonschema >/dev/null 2>&1; then
                    VENV_PYTHON="$VENV_DIR/bin/python3"
                    print_success "Created venv with jsonschema at $VENV_DIR"
                else
                    print_warning "Failed to install jsonschema in new venv"
                    rm -rf "$VENV_DIR"
                fi
            else
                print_warning "Failed to create venv, proceeding without schema validation"
            fi
        fi
    else
        # System Python has jsonschema
        VENV_PYTHON="python3"
    fi

    # Now run schema validation if we have a working Python with jsonschema
    if [ -n "$VENV_PYTHON" ]; then
        if "$VENV_PYTHON" -c "
import json
import sys
from jsonschema import validate, ValidationError, Draft7Validator

try:
    with open('$SCHEMA_FILE', 'r') as f:
        schema = json.load(f)
    with open('$STATE_FILE', 'r') as f:
        state = json.load(f)

    validator = Draft7Validator(schema)
    errors = list(validator.iter_errors(state))

    if errors:
        sys.exit(1)
    else:
        sys.exit(0)
except Exception as e:
    sys.exit(1)
" 2>/dev/null; then
            print_success "State file passes JSON schema validation"
        else
            print_warning "State file has schema validation issues (non-critical for metadata check)"
        fi
    else
        print_info "jsonschema not available, skipping detailed schema validation"
    fi
else
    print_info "Schema file or Python not available, skipping schema validation"
fi

# 5. EFFORT METADATA VALIDATION
print_section "5. EFFORT METADATA VALIDATION"

# Function to validate ISO timestamp
validate_timestamp() {
    local timestamp="$1"
    local field_name="$2"

    if [ -z "$timestamp" ] || [ "$timestamp" == "null" ]; then
        return 1
    fi

    # Check if it matches ISO format (basic check)
    if echo "$timestamp" | grep -qE '^[0-9]{4}-[0-9]{2}-[0-9]{2}T[0-9]{2}:[0-9]{2}:[0-9]{2}'; then
        return 0
    else
        return 1
    fi
}

# Function to validate branch naming format
validate_branch_format() {
    local branch="$1"
    local effort_name="$2"
    local phase="$3"
    local wave="$4"

    # Determine branch type and validate accordingly
    local is_valid=false
    local branch_type=""

    # Check for different branch types based on R014 and other rules
    if [ -n "$PROJECT_PREFIX" ]; then
        # With project prefix
        if echo "$branch" | grep -qE "^${PROJECT_PREFIX}/phase[0-9]+/wave[0-9]+/integration$"; then
            # Wave integration branch: {prefix}/phase{X}/wave{Y}/integration
            branch_type="wave-integration"
            is_valid=true
        elif echo "$branch" | grep -qE "^${PROJECT_PREFIX}/phase[0-9]+/integration$"; then
            # Phase integration branch: {prefix}/phase{X}/integration
            branch_type="phase-integration"
            is_valid=true
        elif echo "$branch" | grep -qE "^${PROJECT_PREFIX}/integration$"; then
            # Project integration branch: {prefix}/integration
            branch_type="project-integration"
            is_valid=true
        elif echo "$branch" | grep -qE "^${PROJECT_PREFIX}/phase[0-9]+/wave[0-9]+/[-a-z0-9.-]+--split-[0-9]{3}$"; then
            # Split branch: {prefix}/phase{X}/wave{Y}/{effort-name}--split-{NNN}
            branch_type="split"
            is_valid=true
        elif echo "$branch" | grep -qE "^${PROJECT_PREFIX}/phase[0-9]+/wave[0-9]+/[-a-z0-9.-]+-fix$"; then
            # Fix branch: {prefix}/phase{X}/wave{Y}/{effort-name}-fix (single dash only)
            branch_type="fix"
            is_valid=true
        elif echo "$branch" | grep -qE "^${PROJECT_PREFIX}/phase[0-9]+/wave[0-9]+/[-a-z0-9.-]+$"; then
            # Standard effort branch: {prefix}/phase{X}/wave{Y}/{effort-name}
            branch_type="effort"
            is_valid=true
        fi
    else
        # Without project prefix
        if echo "$branch" | grep -qE "^phase[0-9]+/wave[0-9]+/integration$"; then
            # Wave integration branch: phase{X}/wave{Y}/integration
            branch_type="wave-integration"
            is_valid=true
        elif echo "$branch" | grep -qE "^phase[0-9]+/integration$"; then
            # Phase integration branch: phase{X}/integration
            branch_type="phase-integration"
            is_valid=true
        elif echo "$branch" | grep -qE "^integration$"; then
            # Project integration branch: integration
            branch_type="project-integration"
            is_valid=true
        elif echo "$branch" | grep -qE "^phase[0-9]+/wave[0-9]+/[-a-z0-9.-]+--split-[0-9]{3}$"; then
            # Split branch: phase{X}/wave{Y}/{effort-name}--split-{NNN}
            branch_type="split"
            is_valid=true
        elif echo "$branch" | grep -qE "^phase[0-9]+/wave[0-9]+/[-a-z0-9.-]+-fix$"; then
            # Fix branch: phase{X}/wave{Y}/{effort-name}-fix (single dash only)
            branch_type="fix"
            is_valid=true
        elif echo "$branch" | grep -qE "^phase[0-9]+/wave[0-9]+/[-a-z0-9.-]+$"; then
            # Standard effort branch: phase{X}/wave{Y}/{effort-name}
            branch_type="effort"
            is_valid=true
        fi
    fi

    if [ "$is_valid" = true ]; then
        # For effort/split/fix branches, verify phase/wave numbers match
        if [[ "$branch_type" =~ ^(effort|split|fix)$ ]]; then
            local branch_phase=$(echo "$branch" | sed 's/.*phase\([0-9]*\).*/\1/')
            local branch_wave=$(echo "$branch" | sed 's/.*wave\([0-9]*\).*/\1/')

            if [ "$branch_phase" == "$phase" ] && [ "$branch_wave" == "$wave" ]; then
                return 0
            else
                return 2  # Phase/wave mismatch
            fi
        else
            # For integration branches, just check if phase number matches (if applicable)
            if [[ "$branch_type" =~ ^(wave|phase)-integration$ ]]; then
                local branch_phase=$(echo "$branch" | sed 's/.*phase\([0-9]*\).*/\1/')
                if [ -n "$phase" ] && [ "$branch_phase" != "$phase" ]; then
                    return 2  # Phase mismatch
                fi
            fi
            return 0
        fi
    else
        return 1  # Format mismatch
    fi
}

# Check efforts_completed
echo ""
echo "Checking efforts_completed..."
COMPLETED_COUNT=$(jq '.efforts_completed | length' "$STATE_FILE")
echo "  Found $COMPLETED_COUNT completed efforts"

if [ "$COMPLETED_COUNT" -gt 0 ]; then
    for i in $(seq 0 $((COMPLETED_COUNT - 1))); do
        EFFORT=$(jq -r ".efforts_completed[$i]" "$STATE_FILE")
        EFFORT_NAME=$(echo "$EFFORT" | jq -r '.name // "unknown"')
        EFFORT_PHASE=$(echo "$EFFORT" | jq -r '.phase // 0')
        EFFORT_WAVE=$(echo "$EFFORT" | jq -r '.wave // 0')

        echo ""
        echo "  Effort: $EFFORT_NAME (Phase $EFFORT_PHASE, Wave $EFFORT_WAVE)"

        # Check repository URL
        EFFORT_REPO=$(echo "$EFFORT" | jq -r '.repository_url // empty')
        if [ -z "$EFFORT_REPO" ]; then
            if $FIX_MODE; then
                print_warning "Missing repository_url for $EFFORT_NAME - auto-fixing"
                # Add repository_url field
                jq ".efforts_completed[$i].repository_url = \"$TARGET_REPO\"" "$STATE_FILE" > "$STATE_FILE.tmp" && mv "$STATE_FILE.tmp" "$STATE_FILE"
                print_success "Added repository_url: $TARGET_REPO"
            else
                print_error "Missing required field 'repository_url' for effort: $EFFORT_NAME"
                echo "    Expected: $TARGET_REPO"
            fi
        elif [ "$EFFORT_REPO" != "$TARGET_REPO" ]; then
            print_error "Repository mismatch for effort: $EFFORT_NAME"
            echo "    Found: $EFFORT_REPO"
            echo "    Expected: $TARGET_REPO"
            echo "    This is a SECURITY VIOLATION - cannot auto-fix"
        else
            print_success "Repository URL valid: matches target"
        fi

        # Check branch naming
        CURRENT_BRANCH=$(echo "$EFFORT" | jq -r '.current_branch // .branch // empty')
        if [ -z "$CURRENT_BRANCH" ]; then
            if $FIX_MODE; then
                # Try to use existing branch field
                OLD_BRANCH=$(echo "$EFFORT" | jq -r '.branch // empty')
                if [ -n "$OLD_BRANCH" ]; then
                    jq ".efforts_completed[$i].current_branch = \"$OLD_BRANCH\"" "$STATE_FILE" > "$STATE_FILE.tmp" && mv "$STATE_FILE.tmp" "$STATE_FILE"
                    print_warning "Migrated 'branch' to 'current_branch' field"
                    CURRENT_BRANCH="$OLD_BRANCH"
                else
                    print_error "Missing both 'current_branch' and 'branch' fields - cannot auto-fix"
                fi
            else
                print_error "Missing required field 'current_branch' for effort: $EFFORT_NAME"
            fi
        fi

        if [ -n "$CURRENT_BRANCH" ]; then
            # Capture function result without triggering set -e
            if validate_branch_format "$CURRENT_BRANCH" "$EFFORT_NAME" "$EFFORT_PHASE" "$EFFORT_WAVE"; then
                BRANCH_RESULT=0
            else
                BRANCH_RESULT=$?
            fi
            if [ $BRANCH_RESULT -eq 0 ]; then
                print_success "Branch format valid: $CURRENT_BRANCH"
            elif [ $BRANCH_RESULT -eq 2 ]; then
                print_error "Branch phase/wave mismatch: $CURRENT_BRANCH"
                echo "    Expected phase $EFFORT_PHASE, wave $EFFORT_WAVE"
            else
                print_error "Invalid branch format: $CURRENT_BRANCH"
                echo "    Valid formats for phase $EFFORT_PHASE, wave $EFFORT_WAVE:"
                if [ -n "$PROJECT_PREFIX" ]; then
                    echo "      - ${PROJECT_PREFIX}/phase${EFFORT_PHASE}/wave${EFFORT_WAVE}/<effort-name> (standard effort)"
                    echo "      - ${PROJECT_PREFIX}/phase${EFFORT_PHASE}/wave${EFFORT_WAVE}/<effort>--split-001 (split branch)"
                    echo "      - ${PROJECT_PREFIX}/phase${EFFORT_PHASE}/wave${EFFORT_WAVE}/<effort>-fix (fix branch)"
                    echo "      - ${PROJECT_PREFIX}/phase${EFFORT_PHASE}/wave${EFFORT_WAVE}/integration (wave integration)"
                    echo "      - ${PROJECT_PREFIX}/phase${EFFORT_PHASE}/integration (phase integration)"
                    echo "      - ${PROJECT_PREFIX}/integration (project integration)"
                else
                    echo "      - phase${EFFORT_PHASE}/wave${EFFORT_WAVE}/<effort-name> (standard effort)"
                    echo "      - phase${EFFORT_PHASE}/wave${EFFORT_WAVE}/<effort>--split-001 (split branch)"
                    echo "      - phase${EFFORT_PHASE}/wave${EFFORT_WAVE}/<effort>-fix (fix branch)"
                    echo "      - phase${EFFORT_PHASE}/wave${EFFORT_WAVE}/integration (wave integration)"
                    echo "      - phase${EFFORT_PHASE}/integration (phase integration)"
                    echo "      - integration (project integration)"
                fi
            fi
        fi

        # Check deprecated branches
        DEPRECATED_BRANCHES=$(echo "$EFFORT" | jq -r '.deprecated_branches // empty')
        if [ -z "$DEPRECATED_BRANCHES" ] || [ "$DEPRECATED_BRANCHES" == "null" ]; then
            if $FIX_MODE; then
                # Initialize empty array
                jq ".efforts_completed[$i].deprecated_branches = []" "$STATE_FILE" > "$STATE_FILE.tmp" && mv "$STATE_FILE.tmp" "$STATE_FILE"
                print_warning "Added empty deprecated_branches array (recommended field)"
            else
                print_warning "Missing 'deprecated_branches' field (recommended for cleanup tracking)"
            fi
        else
            DEPRECATED_COUNT=$(echo "$DEPRECATED_BRANCHES" | jq 'length')
            print_info "Has $DEPRECATED_COUNT deprecated branches tracked"
        fi

        # Check timestamps
        CREATED_AT=$(echo "$EFFORT" | jq -r '.created_at // empty')
        LAST_WORKED=$(echo "$EFFORT" | jq -r '.last_worked_on // empty')
        COMPLETED_AT=$(echo "$EFFORT" | jq -r '.completed_at // .completion_time // empty')

        # Validate created_at
        if ! validate_timestamp "$CREATED_AT" "created_at"; then
            if $FIX_MODE; then
                # Use transition time as fallback
                DEFAULT_TIME=$(jq -r '.transition_time' "$STATE_FILE")
                jq ".efforts_completed[$i].created_at = \"$DEFAULT_TIME\"" "$STATE_FILE" > "$STATE_FILE.tmp" && mv "$STATE_FILE.tmp" "$STATE_FILE"
                print_warning "Added created_at timestamp (using transition_time as default)"
            else
                print_error "Missing or invalid 'created_at' timestamp"
            fi
        else
            print_success "Timestamp 'created_at' valid"
        fi

        # Validate last_worked_on
        if ! validate_timestamp "$LAST_WORKED" "last_worked_on"; then
            if $FIX_MODE; then
                # Use completed_at or current time
                if validate_timestamp "$COMPLETED_AT" "completed_at"; then
                    jq ".efforts_completed[$i].last_worked_on = \"$COMPLETED_AT\"" "$STATE_FILE" > "$STATE_FILE.tmp" && mv "$STATE_FILE.tmp" "$STATE_FILE"
                    print_warning "Added last_worked_on timestamp (using completed_at)"
                else
                    CURRENT_TIME=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
                    jq ".efforts_completed[$i].last_worked_on = \"$CURRENT_TIME\"" "$STATE_FILE" > "$STATE_FILE.tmp" && mv "$STATE_FILE.tmp" "$STATE_FILE"
                    print_warning "Added last_worked_on timestamp (using current time)"
                fi
            else
                print_error "Missing or invalid 'last_worked_on' timestamp"
            fi
        else
            print_success "Timestamp 'last_worked_on' valid"
        fi

        # Validate completed_at
        if ! validate_timestamp "$COMPLETED_AT" "completed_at"; then
            if $FIX_MODE; then
                # For completed efforts, this should exist
                CURRENT_TIME=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
                jq ".efforts_completed[$i].completed_at = \"$CURRENT_TIME\"" "$STATE_FILE" > "$STATE_FILE.tmp" && mv "$STATE_FILE.tmp" "$STATE_FILE"
                print_warning "Added completed_at timestamp (using current time)"
            else
                print_error "Missing or invalid 'completed_at' timestamp"
            fi
        else
            print_success "Timestamp 'completed_at' valid"
        fi

        # Validate chronological order
        if validate_timestamp "$CREATED_AT" "created_at" && validate_timestamp "$LAST_WORKED" "last_worked_on" && validate_timestamp "$COMPLETED_AT" "completed_at"; then
            # Convert to epoch for comparison (if date command supports it)
            if command -v date >/dev/null 2>&1; then
                CREATED_EPOCH=$(date -d "$CREATED_AT" +%s 2>/dev/null || echo 0)
                WORKED_EPOCH=$(date -d "$LAST_WORKED" +%s 2>/dev/null || echo 0)
                COMPLETED_EPOCH=$(date -d "$COMPLETED_AT" +%s 2>/dev/null || echo 0)

                if [ "$CREATED_EPOCH" -gt 0 ] && [ "$WORKED_EPOCH" -gt 0 ] && [ "$COMPLETED_EPOCH" -gt 0 ]; then
                    if [ "$CREATED_EPOCH" -le "$WORKED_EPOCH" ] && [ "$WORKED_EPOCH" -le "$COMPLETED_EPOCH" ]; then
                        print_success "Timestamp chronology valid (created <= worked <= completed)"
                    else
                        print_error "Timestamp chronology invalid!"
                        echo "    Created: $CREATED_AT"
                        echo "    Worked: $LAST_WORKED"
                        echo "    Completed: $COMPLETED_AT"
                    fi
                fi
            fi
        fi
    done
fi

# Check efforts_in_progress
echo ""
echo "Checking efforts_in_progress..."
IN_PROGRESS_COUNT=$(jq '.efforts_in_progress | length' "$STATE_FILE")
echo "  Found $IN_PROGRESS_COUNT in-progress efforts"

if [ "$IN_PROGRESS_COUNT" -gt 0 ]; then
    for i in $(seq 0 $((IN_PROGRESS_COUNT - 1))); do
        EFFORT=$(jq -r ".efforts_in_progress[$i]" "$STATE_FILE")
        EFFORT_NAME=$(echo "$EFFORT" | jq -r '.name // "unknown"')
        EFFORT_PHASE=$(echo "$EFFORT" | jq -r '.phase // 0')
        EFFORT_WAVE=$(echo "$EFFORT" | jq -r '.wave // 0')

        echo ""
        echo "  Effort: $EFFORT_NAME (Phase $EFFORT_PHASE, Wave $EFFORT_WAVE)"

        # Check repository URL
        EFFORT_REPO=$(echo "$EFFORT" | jq -r '.repository_url // empty')
        if [ -z "$EFFORT_REPO" ]; then
            if $FIX_MODE; then
                jq ".efforts_in_progress[$i].repository_url = \"$TARGET_REPO\"" "$STATE_FILE" > "$STATE_FILE.tmp" && mv "$STATE_FILE.tmp" "$STATE_FILE"
                print_success "Added repository_url: $TARGET_REPO"
            else
                print_error "Missing required field 'repository_url'"
            fi
        elif [ "$EFFORT_REPO" != "$TARGET_REPO" ]; then
            print_error "Repository mismatch!"
            echo "    Found: $EFFORT_REPO"
            echo "    Expected: $TARGET_REPO"
        else
            print_success "Repository URL valid"
        fi

        # Check branch fields
        CURRENT_BRANCH=$(echo "$EFFORT" | jq -r '.current_branch // .branch // empty')
        if [ -z "$CURRENT_BRANCH" ]; then
            if $FIX_MODE; then
                OLD_BRANCH=$(echo "$EFFORT" | jq -r '.branch // empty')
                if [ -n "$OLD_BRANCH" ]; then
                    jq ".efforts_in_progress[$i].current_branch = \"$OLD_BRANCH\"" "$STATE_FILE" > "$STATE_FILE.tmp" && mv "$STATE_FILE.tmp" "$STATE_FILE"
                    print_warning "Migrated 'branch' to 'current_branch'"
                    CURRENT_BRANCH="$OLD_BRANCH"
                fi
            else
                print_error "Missing 'current_branch' field"
            fi
        fi

        if [ -n "$CURRENT_BRANCH" ]; then
            # Capture function result without triggering set -e
            if validate_branch_format "$CURRENT_BRANCH" "$EFFORT_NAME" "$EFFORT_PHASE" "$EFFORT_WAVE"; then
                BRANCH_RESULT=0
            else
                BRANCH_RESULT=$?
            fi
            if [ $BRANCH_RESULT -eq 0 ]; then
                print_success "Branch format valid: $CURRENT_BRANCH"
            elif [ $BRANCH_RESULT -eq 2 ]; then
                print_error "Branch phase/wave mismatch: $CURRENT_BRANCH"
                echo "    Expected phase $EFFORT_PHASE, wave $EFFORT_WAVE"
            else
                print_error "Invalid branch format: $CURRENT_BRANCH"
                echo "    Valid formats for phase $EFFORT_PHASE, wave $EFFORT_WAVE:"
                if [ -n "$PROJECT_PREFIX" ]; then
                    echo "      - ${PROJECT_PREFIX}/phase${EFFORT_PHASE}/wave${EFFORT_WAVE}/<effort-name> (standard effort)"
                    echo "      - ${PROJECT_PREFIX}/phase${EFFORT_PHASE}/wave${EFFORT_WAVE}/<effort>--split-001 (split branch)"
                    echo "      - ${PROJECT_PREFIX}/phase${EFFORT_PHASE}/wave${EFFORT_WAVE}/<effort>-fix (fix branch)"
                    echo "      - ${PROJECT_PREFIX}/phase${EFFORT_PHASE}/wave${EFFORT_WAVE}/integration (wave integration)"
                    echo "      - ${PROJECT_PREFIX}/phase${EFFORT_PHASE}/integration (phase integration)"
                    echo "      - ${PROJECT_PREFIX}/integration (project integration)"
                else
                    echo "      - phase${EFFORT_PHASE}/wave${EFFORT_WAVE}/<effort-name> (standard effort)"
                    echo "      - phase${EFFORT_PHASE}/wave${EFFORT_WAVE}/<effort>--split-001 (split branch)"
                    echo "      - phase${EFFORT_PHASE}/wave${EFFORT_WAVE}/<effort>-fix (fix branch)"
                    echo "      - phase${EFFORT_PHASE}/wave${EFFORT_WAVE}/integration (wave integration)"
                    echo "      - phase${EFFORT_PHASE}/integration (phase integration)"
                    echo "      - integration (project integration)"
                fi
            fi
        fi

        # Check timestamps for in-progress efforts
        CREATED_AT=$(echo "$EFFORT" | jq -r '.created_at // empty')
        LAST_WORKED=$(echo "$EFFORT" | jq -r '.last_worked_on // .last_updated_at // empty')

        if ! validate_timestamp "$CREATED_AT" "created_at"; then
            if $FIX_MODE; then
                CURRENT_TIME=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
                jq ".efforts_in_progress[$i].created_at = \"$CURRENT_TIME\"" "$STATE_FILE" > "$STATE_FILE.tmp" && mv "$STATE_FILE.tmp" "$STATE_FILE"
                print_warning "Added created_at timestamp"
            else
                print_error "Missing 'created_at' timestamp"
            fi
        else
            print_success "Timestamp 'created_at' valid"
        fi

        if ! validate_timestamp "$LAST_WORKED" "last_worked_on"; then
            if $FIX_MODE; then
                CURRENT_TIME=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
                jq ".efforts_in_progress[$i].last_worked_on = \"$CURRENT_TIME\"" "$STATE_FILE" > "$STATE_FILE.tmp" && mv "$STATE_FILE.tmp" "$STATE_FILE"
                print_warning "Added last_worked_on timestamp"
            else
                print_error "Missing 'last_worked_on' timestamp"
            fi
        else
            print_success "Timestamp 'last_worked_on' valid"
        fi
    done
fi

# 6. PROJECT PREFIX VALIDATION
print_section "6. PROJECT PREFIX VALIDATION"
if [ -n "$PROJECT_PREFIX" ]; then
    print_info "Project prefix configured: $PROJECT_PREFIX"

    # Check if all branches follow the prefix
    ALL_BRANCHES=$(jq -r '[.efforts_completed[].current_branch // .efforts_completed[].branch // empty, .efforts_in_progress[].current_branch // .efforts_in_progress[].branch // empty] | map(select(. != "")) | unique | .[]' "$STATE_FILE" 2>/dev/null)

    if [ -n "$ALL_BRANCHES" ]; then
        INVALID_PREFIX_COUNT=0
        while IFS= read -r branch; do
            if [ -n "$branch" ] && [ "$branch" != "null" ]; then
                if ! echo "$branch" | grep -q "^${PROJECT_PREFIX}/"; then
                    print_warning "Branch missing project prefix: $branch"
                    ((INVALID_PREFIX_COUNT++))
                fi
            fi
        done <<< "$ALL_BRANCHES"

        if [ $INVALID_PREFIX_COUNT -eq 0 ]; then
            print_success "All branches use correct project prefix"
        fi
    fi
else
    print_info "No project prefix configured"
fi

# 7. FINAL SUMMARY
print_section "7. VALIDATION SUMMARY"
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
if [ "$VALIDATION_FAILED" = true ]; then
    echo -e "${RED}❌ VALIDATION FAILED${NC}"
    echo ""
    echo "Critical failures: $CRITICAL_FAILURES"
    echo "Warnings: $WARNINGS"
    echo ""
    if $FIX_MODE; then
        echo "Some issues were auto-fixed. Please review the changes and run validation again."
    else
        echo "Run with --fix flag to attempt auto-fixing some issues:"
        echo "  $0 --fix $STATE_FILE"
    fi
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    exit 1
else
    if [ $WARNINGS -gt 0 ]; then
        echo -e "${YELLOW}✅ VALIDATION PASSED WITH WARNINGS${NC}"
        echo ""
        echo "Warnings: $WARNINGS"
    else
        echo -e "${GREEN}✅ VALIDATION PASSED${NC}"
        echo ""
        echo "All effort metadata checks passed!"
    fi
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    exit 0
fi