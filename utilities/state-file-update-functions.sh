#!/bin/bash

# state-file-update-functions.sh - Helper functions for safe state file updates
# Source this file in orchestrator scripts for validated state management

set -euo pipefail

# Color codes for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to validate state file
validate_state_file() {
    local state_file="${1:-orchestrator-state.json}"
    
    if [ ! -f "$state_file" ]; then
        echo -e "${RED}❌ State file not found: $state_file${NC}"
        return 1
    fi
    
    # Use the validation tool if available
    if [ -f "$CLAUDE_PROJECT_DIR/tools/validate-state.sh" ]; then
        "$CLAUDE_PROJECT_DIR/tools/validate-state.sh" "$state_file"
        return $?
    else
        # Basic validation fallback
        if ! jq empty "$state_file" 2>/dev/null; then
            echo -e "${RED}❌ State file is not valid JSON!${NC}"
            return 1
        fi
        
        # Check for required fields
        local required_fields=(
            "current_state"
            "previous_state"
            "transition_time"
            "current_phase"
            "current_wave"
        )
        
        for field in "${required_fields[@]}"; do
            if ! jq -e "has(\"$field\")" "$state_file" >/dev/null 2>&1; then
                echo -e "${RED}❌ Missing required field: $field${NC}"
                return 1
            fi
        done
        
        echo -e "${GREEN}✅ State file passed basic validation${NC}"
        return 0
    fi
}

# Function to safely update a single field in state file
safe_update_field() {
    local field="$1"
    local value="$2"
    local state_file="${3:-orchestrator-state.json}"
    
    # Create backup
    cp "$state_file" "${state_file}.backup"
    
    # Update field
    if [[ "$value" =~ ^[0-9]+$ ]]; then
        # Numeric value
        jq ".$field = $value" "$state_file" > tmp.json && mv tmp.json "$state_file"
    elif [[ "$value" == "true" || "$value" == "false" || "$value" == "null" ]]; then
        # Boolean or null
        jq ".$field = $value" "$state_file" > tmp.json && mv tmp.json "$state_file"
    else
        # String value
        jq ".$field = \"$value\"" "$state_file" > tmp.json && mv tmp.json "$state_file"
    fi
    
    # Validate after update
    if ! validate_state_file "$state_file"; then
        echo -e "${RED}❌ Validation failed after updating $field${NC}"
        echo "Restoring backup..."
        mv "${state_file}.backup" "$state_file"
        return 1
    fi
    
    # Remove backup on success
    rm -f "${state_file}.backup"
    return 0
}

# Function to perform safe state transition with validation
safe_state_transition() {
    local new_state="$1"
    local reason="${2:-State transition}"
    local state_file="${3:-orchestrator-state.json}"
    
    # Get current state
    local current_state=$(jq -r '.current_state' "$state_file")
    
    if [ "$current_state" == "$new_state" ]; then
        echo -e "${YELLOW}⚠️  Already in state $new_state${NC}"
        return 0
    fi
    
    echo "🔄 Transitioning: $current_state → $new_state"
    
    # Create backup before transition
    cp "$state_file" "${state_file}.pre-transition"
    
    # Update all transition fields
    local timestamp=$(date -u +%Y-%m-%dT%H:%M:%SZ)
    
    jq ".current_state = \"$new_state\" |
        .previous_state = \"$current_state\" |
        .transition_time = \"$timestamp\" |
        .transition_reason = \"$reason\"" "$state_file" > tmp.json && mv tmp.json "$state_file"
    
    # Add to transition history if the field exists
    if jq -e 'has("transition_history")' "$state_file" >/dev/null 2>&1; then
        local transition_entry="{\"from\": \"$current_state\", \"to\": \"$new_state\", \"time\": \"$timestamp\", \"reason\": \"$reason\"}"
        jq ".transition_history = (.transition_history // []) + [$transition_entry]" "$state_file" > tmp.json && mv tmp.json "$state_file"
    fi
    
    # Validate the updated state
    if ! validate_state_file "$state_file"; then
        echo -e "${RED}❌❌❌ State validation failed after transition!${NC}"
        echo "Restoring pre-transition state..."
        mv "${state_file}.pre-transition" "$state_file"
        return 1
    fi
    
    # Commit and push immediately (R288)
    git add "$state_file"
    git commit -m "state: $current_state → $new_state - $reason [R288/R324]"
    
    if ! git push; then
        echo -e "${RED}❌ Failed to push state transition!${NC}"
        echo "Attempting force push with lease..."
        git push --force-with-lease || {
            echo -e "${RED}❌❌❌ CRITICAL: Could not push state change!${NC}"
            return 1
        }
    fi
    
    # Clean up backup
    rm -f "${state_file}.pre-transition"
    
    echo -e "${GREEN}✅ State transition complete: $current_state → $new_state${NC}"
    return 0
}

# Function to update and validate effort status
update_effort_status() {
    local effort_name="$1"
    local status="$2"
    local phase="${3:-}"
    local wave="${4:-}"
    local state_file="${5:-orchestrator-state.json}"
    
    echo "📝 Updating effort '$effort_name' status to '$status'"
    
    # Find and update the effort in the appropriate array
    local updated=false
    
    # Try efforts_in_progress first
    if jq -e ".efforts_in_progress[] | select(.name == \"$effort_name\")" "$state_file" >/dev/null 2>&1; then
        jq "(.efforts_in_progress[] | select(.name == \"$effort_name\") | .status) = \"$status\"" "$state_file" > tmp.json && mv tmp.json "$state_file"
        updated=true
    fi
    
    # Try efforts_completed if not found
    if [ "$updated" = false ] && jq -e ".efforts_completed[] | select(.name == \"$effort_name\")" "$state_file" >/dev/null 2>&1; then
        jq "(.efforts_completed[] | select(.name == \"$effort_name\") | .status) = \"$status\"" "$state_file" > tmp.json && mv tmp.json "$state_file"
        updated=true
    fi
    
    if [ "$updated" = false ]; then
        echo -e "${YELLOW}⚠️  Effort '$effort_name' not found in state file${NC}"
        return 1
    fi
    
    # Validate after update
    if ! validate_state_file "$state_file"; then
        echo -e "${RED}❌ State validation failed after updating effort status${NC}"
        return 1
    fi
    
    # Commit the change
    git add "$state_file"
    git commit -m "state: effort '$effort_name' status → $status [R288]"
    git push
    
    echo -e "${GREEN}✅ Effort status updated and validated${NC}"
    return 0
}

# Function to add violation record with validation
add_violation() {
    local effort="$1"
    local lines="$2"
    local limit="${3:-800}"
    local state_file="${4:-orchestrator-state.json}"
    
    echo "⚠️  Recording violation for effort '$effort': $lines lines (limit: $limit)"
    
    # Create violation entry
    local violation_entry="{
        \"effort\": \"$effort\",
        \"lines\": $lines,
        \"limit\": $limit,
        \"requires_split\": true,
        \"recorded_at\": \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\"
    }"
    
    # Add to violations array
    jq ".violations = (.violations // []) + [$violation_entry]" "$state_file" > tmp.json && mv tmp.json "$state_file"
    
    # Validate
    if ! validate_state_file "$state_file"; then
        echo -e "${RED}❌ State validation failed after adding violation${NC}"
        return 1
    fi
    
    # Commit
    git add "$state_file"
    git commit -m "state: violation recorded for '$effort' - $lines lines [R288]"
    git push
    
    echo -e "${GREEN}✅ Violation recorded and validated${NC}"
    return 0
}

# Function to verify state machine transition is valid
verify_transition_allowed() {
    local from_state="$1"
    local to_state="$2"
    local state_machine_file="${3:-$CLAUDE_PROJECT_DIR/SOFTWARE-FACTORY-STATE-MACHINE.md}"
    
    # Check if transition exists in state machine
    if grep -q "$from_state --> $to_state" "$state_machine_file"; then
        echo -e "${GREEN}✅ Transition allowed: $from_state → $to_state${NC}"
        return 0
    else
        echo -e "${RED}❌ Invalid transition: $from_state → $to_state${NC}"
        echo "This transition is not defined in the state machine!"
        return 1
    fi
}

# Function to perform emergency state recovery
emergency_state_recovery() {
    local state_file="${1:-orchestrator-state.json}"
    
    echo -e "${YELLOW}🚨 Attempting emergency state recovery...${NC}"
    
    # Check for backups
    if [ -f "${state_file}.backup" ]; then
        echo "Found backup file, validating..."
        if validate_state_file "${state_file}.backup"; then
            cp "${state_file}.backup" "$state_file"
            echo -e "${GREEN}✅ Restored from backup${NC}"
            return 0
        fi
    fi
    
    # Check git history for last valid version
    echo "Checking git history for last valid state..."
    git checkout HEAD -- "$state_file"
    
    if validate_state_file "$state_file"; then
        echo -e "${GREEN}✅ Restored from git${NC}"
        return 0
    fi
    
    echo -e "${RED}❌ Could not recover valid state file${NC}"
    echo "Manual intervention required!"
    return 1
}

# Export functions for use in other scripts
export -f validate_state_file
export -f safe_update_field
export -f safe_state_transition
export -f update_effort_status
export -f add_violation
export -f verify_transition_allowed
export -f emergency_state_recovery

echo "✅ State file update functions loaded"
echo "Available functions:"
echo "  - validate_state_file [file]"
echo "  - safe_update_field <field> <value> [file]"
echo "  - safe_state_transition <new_state> [reason] [file]"
echo "  - update_effort_status <effort> <status> [phase] [wave] [file]"
echo "  - add_violation <effort> <lines> [limit] [file]"
echo "  - verify_transition_allowed <from> <to> [state_machine_file]"
echo "  - emergency_state_recovery [file]"