#!/bin/bash
#
# enforce-state-validation.sh
# Automated enforcement of R406: Orchestrator State Schema Validation
# This script provides comprehensive validation and auto-fix capabilities
#

set -euo pipefail

# Color codes for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Get script directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

# Configuration
STATE_FILE="${1:-$PROJECT_ROOT/orchestrator-state-v3.json}"
SCHEMA_FILE="$PROJECT_ROOT/orchestrator-state.schema.json"
STATE_MACHINE="$PROJECT_ROOT/state-machines/software-factory-3.0-state-machine.json"
BACKUP_DIR="$PROJECT_ROOT/state-backups"
VALIDATION_LOG="$PROJECT_ROOT/validation.log"

# Create backup directory if needed
mkdir -p "$BACKUP_DIR"

# Function to log messages
log_message() {
    local level="$1"
    local message="$2"
    local timestamp=$(date '+%Y-%m-%d %H:%M:%S')
    echo -e "${level}${message}${NC}"
    echo "[$timestamp] $level: $message" >> "$VALIDATION_LOG"
}

# Function to create backup
backup_state() {
    local backup_file="$BACKUP_DIR/orchestrator-state-$(date +%Y%m%d-%H%M%S).json"
    if [ -f "$STATE_FILE" ]; then
        cp "$STATE_FILE" "$backup_file"
        log_message "$BLUE" "📁 Backup created: $backup_file"

        # Keep only last 50 backups
        ls -t "$BACKUP_DIR"/orchestrator-state-*.json 2>/dev/null | tail -n +51 | xargs -r rm
    fi
}

# Function to validate JSON syntax
validate_json_syntax() {
    local file="$1"

    if ! jq empty "$file" 2>/dev/null; then
        log_message "$RED" "❌ Invalid JSON syntax in $file"
        jq empty "$file" 2>&1 | sed 's/^/   /'
        return 1
    fi

    log_message "$GREEN" "✅ JSON syntax valid"
    return 0
}

# Function to validate required fields
validate_required_fields() {
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
        log_message "$RED" "❌ Missing required fields:"
        for field in "${missing_fields[@]}"; do
            echo "   - $field"
        done
        return 1
    fi

    log_message "$GREEN" "✅ All required fields present"
    return 0
}

# Function to validate state against state machine
validate_state_value() {
    local current_state=$(jq -r '.current_state' "$STATE_FILE")

    if [ "$current_state" = "null" ] || [ -z "$current_state" ]; then
        log_message "$RED" "❌ current_state is null or empty"
        return 1
    fi

    # Check if state exists in schema enum list
    local valid_state=$(jq -r ".properties.current_state.enum[] | select(. == \"$current_state\")" "$SCHEMA_FILE" 2>/dev/null)

    if [ -z "$valid_state" ]; then
        # Also check in state machine mermaid diagram as fallback
        if ! grep -q "$current_state" "$STATE_MACHINE" 2>/dev/null; then
            log_message "$RED" "❌ Invalid state '$current_state' not in schema or state machine"
            echo "   Valid states include: INIT, PLANNING, WAVE_START, WAVE_COMPLETE, etc."
            echo "   Check orchestrator-state.schema.json for full list"
            return 1
        fi
    fi

    log_message "$GREEN" "✅ State '$current_state' is valid"
    return 0
}

# Function to validate with schema
validate_with_schema() {
    if [ -f "$PROJECT_ROOT/tools/validate-state.sh" ]; then
        if "$PROJECT_ROOT/tools/validate-state.sh" "$STATE_FILE" >/dev/null 2>&1; then
            log_message "$GREEN" "✅ Schema validation passed"
            return 0
        else
            log_message "$RED" "❌ Schema validation failed"
            "$PROJECT_ROOT/tools/validate-state.sh" "$STATE_FILE" 2>&1 | sed 's/^/   /'
            return 1
        fi
    else
        log_message "$YELLOW" "⚠️  Schema validator not found, skipping schema validation"
        return 0
    fi
}

# Function to auto-fix missing fields
auto_fix_missing_fields() {
    log_message "$YELLOW" "🔧 Attempting to auto-fix missing fields..."

    local temp_file="$STATE_FILE.fix.tmp"

    # Add missing required fields with defaults
    jq '. + {
        "project_info": (.project_info // {
            "name": "unnamed-project",
            "description": "No description provided",
            "start_date": (now | strftime("%Y-%m-%d"))
        }),
        "current_phase": (.current_phase // 1),
        "current_wave": (.current_wave // 0),
        "current_state": (.current_state // "INIT"),
        "previous_state": (.previous_state // null),
        "transition_time": (.transition_time // (now | strftime("%Y-%m-%dT%H:%M:%SZ"))),
        "phases_planned": (.phases_planned // 1),
        "waves_per_phase": (.waves_per_phase // [1]),
        "efforts_completed": (.efforts_completed // []),
        "efforts_in_progress": (.efforts_in_progress // []),
        "efforts_pending": (.efforts_pending // [])
    }' "$STATE_FILE" > "$temp_file"

    if validate_json_syntax "$temp_file"; then
        mv "$temp_file" "$STATE_FILE"
        log_message "$GREEN" "✅ Missing fields added with defaults"
        return 0
    else
        rm -f "$temp_file"
        log_message "$RED" "❌ Auto-fix failed"
        return 1
    fi
}

# Function to fix invalid state value
auto_fix_invalid_state() {
    local current_state=$(jq -r '.current_state' "$STATE_FILE")

    if ! grep -q "^STATE: $current_state" "$STATE_MACHINE" 2>/dev/null; then
        log_message "$YELLOW" "🔧 Fixing invalid state '$current_state' -> 'INIT'"

        local temp_file="$STATE_FILE.fix.tmp"
        jq '.current_state = "INIT" | .previous_state = null' "$STATE_FILE" > "$temp_file"

        if validate_json_syntax "$temp_file"; then
            mv "$temp_file" "$STATE_FILE"
            log_message "$GREEN" "✅ State reset to INIT"
            return 0
        else
            rm -f "$temp_file"
            return 1
        fi
    fi
    return 0
}

# Function to validate state transitions
validate_transition() {
    local from_state="$1"
    local to_state="$2"

    # Extract valid transitions from state machine
    local valid_transitions=$(grep -A 30 "^STATE: $from_state" "$STATE_MACHINE" 2>/dev/null | \
                              grep "^### TRANSITIONS TO:" -A 20 | \
                              grep "^- " | sed 's/^- //')

    if echo "$valid_transitions" | grep -q "^$to_state"; then
        log_message "$GREEN" "✅ Valid transition: $from_state -> $to_state"
        return 0
    else
        log_message "$RED" "❌ Invalid transition: $from_state -> $to_state"
        echo "   Valid transitions from $from_state:"
        echo "$valid_transitions" | sed 's/^/      - /'
        return 1
    fi
}

# Function to generate validation report
generate_report() {
    echo ""
    echo "════════════════════════════════════════════════════════════"
    echo "     ORCHESTRATOR STATE VALIDATION REPORT"
    echo "════════════════════════════════════════════════════════════"
    echo ""
    echo "File: $STATE_FILE"
    echo "Schema: $SCHEMA_FILE"
    echo "Timestamp: $(date '+%Y-%m-%d %H:%M:%S')"
    echo ""

    local all_valid=true

    # Check file exists
    if [ ! -f "$STATE_FILE" ]; then
        log_message "$RED" "❌ State file not found!"
        return 1
    fi

    # Run validations
    echo "1. JSON SYNTAX VALIDATION"
    echo "────────────────────────────────────────────────────────────"
    validate_json_syntax "$STATE_FILE" || all_valid=false
    echo ""

    echo "2. REQUIRED FIELDS VALIDATION"
    echo "────────────────────────────────────────────────────────────"
    validate_required_fields || all_valid=false
    echo ""

    echo "3. STATE MACHINE VALIDATION"
    echo "────────────────────────────────────────────────────────────"
    validate_state_value || all_valid=false
    echo ""

    echo "4. SCHEMA VALIDATION"
    echo "────────────────────────────────────────────────────────────"
    validate_with_schema || all_valid=false
    echo ""

    # Summary
    echo "════════════════════════════════════════════════════════════"
    if [ "$all_valid" = true ]; then
        log_message "$GREEN" "✅ ALL VALIDATIONS PASSED!"
        return 0
    else
        log_message "$RED" "❌ VALIDATION FAILED - See errors above"
        return 1
    fi
}

# Function to setup git hooks
setup_git_hooks() {
    local hook_file="$PROJECT_ROOT/.git/hooks/pre-commit"

    if [ ! -d "$PROJECT_ROOT/.git/hooks" ]; then
        log_message "$YELLOW" "⚠️  Git hooks directory not found"
        return 1
    fi

    cat > "$hook_file" << 'EOF'
#!/bin/bash
# Pre-commit hook for R406 enforcement

# Check if orchestrator-state-v3.json is being committed
if git diff --cached --name-only | grep -q "orchestrator-state-v3.json"; then
    echo "🔍 Validating orchestrator-state-v3.json before commit..."

    # Run validation
    if utilities/enforce-state-validation.sh; then
        echo "✅ State file validation passed"
    else
        echo "❌ COMMIT BLOCKED: State file validation failed"
        echo "   Run: utilities/enforce-state-validation.sh --fix"
        exit 1
    fi
fi
EOF

    chmod +x "$hook_file"
    log_message "$GREEN" "✅ Git pre-commit hook installed"
}

# Main execution
main() {
    local mode="${2:-validate}"

    case "$mode" in
        --validate|validate)
            backup_state
            generate_report
            exit $?
            ;;

        --fix|fix)
            backup_state
            log_message "$BLUE" "🔧 Running auto-fix mode..."

            # Try to fix issues
            if ! validate_json_syntax "$STATE_FILE"; then
                log_message "$RED" "❌ Cannot auto-fix invalid JSON syntax"
                exit 1
            fi

            if ! validate_required_fields; then
                auto_fix_missing_fields
            fi

            if ! validate_state_value; then
                auto_fix_invalid_state
            fi

            # Validate after fixes
            generate_report
            exit $?
            ;;

        --setup-hooks|setup-hooks)
            setup_git_hooks
            ;;

        --transition|transition)
            # Validate a state transition
            if [ $# -lt 4 ]; then
                echo "Usage: $0 <state-file> --transition <from-state> <to-state>"
                exit 1
            fi
            validate_transition "$3" "$4"
            ;;

        --help|help)
            cat << EOF
Usage: $0 [state-file] [mode]

Modes:
  --validate    Validate state file (default)
  --fix         Attempt to auto-fix validation issues
  --setup-hooks Install git pre-commit hooks
  --transition  Validate state transition (requires from/to states)
  --help        Show this help message

Examples:
  $0                                    # Validate default state file
  $0 orchestrator-state-v3.json --fix     # Auto-fix issues
  $0 any --transition INIT PLANNING    # Validate transition

Environment Variables:
  PROJECT_ROOT    Override project root directory
  STATE_FILE      Override default state file path
  SCHEMA_FILE     Override schema file path

Exit Codes:
  0  Validation passed
  1  Validation failed

EOF
            ;;

        *)
            echo "Unknown mode: $mode"
            echo "Use --help for usage information"
            exit 1
            ;;
    esac
}

# Run main function
main "$@"