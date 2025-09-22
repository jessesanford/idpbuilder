#!/bin/bash
# R346: Orchestrator Recovery Validation Script
# Validates and repairs state consistency after crashes/restarts

set -e

STATE_FILE="${1:-orchestrator-state.json}"
PROJECT_DIR="${PROJECT_DIR:-/home/vscode/software-factory-template}"
RECOVERY_LOG="${2:-recovery-validation.log}"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Start logging
exec > >(tee -a "$RECOVERY_LOG")
exec 2>&1

echo -e "${BLUE}═══════════════════════════════════════════════════════════════${NC}"
echo -e "${BLUE}       ORCHESTRATOR RECOVERY VALIDATION (R346)                    ${NC}"
echo -e "${BLUE}═══════════════════════════════════════════════════════════════${NC}"
echo "Timestamp: $(date -Iseconds)"
echo "State file: $STATE_FILE"
echo ""

# Function to clean orphaned stale tracking
clean_orphaned_stale_tracking() {
    echo -e "${YELLOW}Checking for orphaned stale tracking...${NC}"
    
    local cleaned_count=0
    
    # Get all stale integrations that claim to need recreation
    local stale_ids=$(jq -r '.stale_integration_tracking.stale_integrations[]? | 
                             select(.recreation_completed == false) | 
                             .integration_id' "$STATE_FILE" 2>/dev/null)
    
    for id in $stale_ids; do
        echo "Checking stale integration: $id"
        
        # Check if this integration is actually complete
        local is_complete=false
        
        # Check in waves_completed
        if jq -e --arg id "$id" '.waves_completed | .. | 
                                 select(.integration_branch? == $id)' \
                                 "$STATE_FILE" > /dev/null 2>&1; then
            is_complete=true
            echo -e "${YELLOW}  Found in waves_completed - marking as recreated${NC}"
        fi
        
        # Check in phases_completed
        if jq -e --arg id "$id" '.phases_completed | .. | 
                                 select(.integration_branch? == $id)' \
                                 "$STATE_FILE" > /dev/null 2>&1; then
            is_complete=true
            echo -e "${YELLOW}  Found in phases_completed - marking as recreated${NC}"
        fi
        
        # Check in integration_status
        local status_key="wave_$id"
        local status=$(jq -r --arg key "$status_key" '.integration_status[$key].status // "none"' "$STATE_FILE")
        if [[ "$status" == "COMPLETE" ]]; then
            is_complete=true
            echo -e "${YELLOW}  Integration status is COMPLETE - marking as recreated${NC}"
        fi
        
        if $is_complete; then
            # Clean up the orphaned tracking
            jq --arg id "$id" --arg time "$(date -Iseconds)" '
                (.stale_integration_tracking.stale_integrations[]? | 
                 select(.integration_id == $id)) |= 
                 . + {
                    "recreation_completed": true,
                    "recreation_at": $time,
                    "recreation_required": false,
                    "cleanup_reason": "Orphaned tracking cleaned during recovery"
                 }
            ' "$STATE_FILE" > tmp.json && mv tmp.json "$STATE_FILE"
            
            echo -e "${GREEN}  ✅ Cleaned orphaned stale tracking for $id${NC}"
            ((cleaned_count++))
        fi
    done
    
    if [ $cleaned_count -gt 0 ]; then
        echo -e "${GREEN}Cleaned $cleaned_count orphaned stale tracking entries${NC}"
    else
        echo -e "${GREEN}No orphaned stale tracking found${NC}"
    fi
}

# Function to validate current integration consistency
validate_current_integrations() {
    echo -e "${YELLOW}Validating current integration consistency...${NC}"
    
    local fixed_count=0
    
    for type in wave phase project; do
        # Check if integration exists
        local has_integration=$(jq -r --arg type "$type" '
            has("current_" + $type + "_integration") and 
            .["current_" + $type + "_integration"] != null' "$STATE_FILE")
        
        if [[ "$has_integration" == "true" ]]; then
            echo "Checking current $type integration..."
            
            # Fix: is_stale=false but has staleness_reason
            local needs_fix=$(jq -r --arg type "$type" '
                .["current_" + $type + "_integration"] | 
                select(.is_stale == false and .staleness_reason != null) | 
                "fix_needed"' "$STATE_FILE" 2>/dev/null)
            
            if [[ "$needs_fix" == "fix_needed" ]]; then
                echo -e "${YELLOW}  Fixing: not stale but has staleness reason${NC}"
                jq --arg type "$type" '
                    .["current_" + $type + "_integration"].staleness_reason = null |
                    .["current_" + $type + "_integration"].stale_since = null |
                    .["current_" + $type + "_integration"].stale_due_to_fixes = []
                ' "$STATE_FILE" > tmp.json && mv tmp.json "$STATE_FILE"
                ((fixed_count++))
            fi
            
            # Fix: is_stale=true but no staleness_reason
            needs_fix=$(jq -r --arg type "$type" '
                .["current_" + $type + "_integration"] | 
                select(.is_stale == true and .staleness_reason == null) | 
                "fix_needed"' "$STATE_FILE" 2>/dev/null)
            
            if [[ "$needs_fix" == "fix_needed" ]]; then
                echo -e "${YELLOW}  Fixing: stale but no reason - clearing stale flag${NC}"
                jq --arg type "$type" '
                    .["current_" + $type + "_integration"].is_stale = false
                ' "$STATE_FILE" > tmp.json && mv tmp.json "$STATE_FILE"
                ((fixed_count++))
            fi
        fi
    done
    
    if [ $fixed_count -gt 0 ]; then
        echo -e "${GREEN}Fixed $fixed_count integration consistency issues${NC}"
    else
        echo -e "${GREEN}All current integrations are consistent${NC}"
    fi
}

# Function to check state transition validity
check_state_validity() {
    echo -e "${YELLOW}Checking state transition validity...${NC}"
    
    local current_state=$(jq -r '.current_state' "$STATE_FILE")
    local previous_state=$(jq -r '.previous_state' "$STATE_FILE")
    
    echo "Current state: $current_state"
    echo "Previous state: $previous_state"
    
    # Check for stuck states
    if [[ "$current_state" == "$previous_state" ]]; then
        echo -e "${YELLOW}⚠️ WARNING: current_state same as previous_state${NC}"
        echo "This might indicate a failed transition"
    fi
    
    # Validate state name
    local valid_states=(
        "INIT" "PLANNING" "SETUP_EFFORT_INFRASTRUCTURE" 
        "SPAWN_AGENTS" "MONITOR" "WAVE_COMPLETE" 
        "INTEGRATION" "MONITORING_INTEGRATION" 
        "INTEGRATION_CODE_REVIEW" "INTEGRATION_FEEDBACK_REVIEW"
        "ERROR_RECOVERY" "WAVE_START" "WAVE_REVIEW"
        "PHASE_INTEGRATION" "PROJECT_INTEGRATION"
    )
    
    if [[ ! " ${valid_states[@]} " =~ " ${current_state} " ]]; then
        echo -e "${RED}❌ ERROR: Invalid current_state: $current_state${NC}"
        echo "This requires manual intervention to fix"
        return 1
    fi
    
    echo -e "${GREEN}State names are valid${NC}"
}

# Function to reconcile effort tracking
reconcile_effort_tracking() {
    echo -e "${YELLOW}Reconciling effort tracking...${NC}"
    
    # Check for duplicates between completed and in_progress
    local duplicates=$(jq -r '
        [.efforts_completed[].name] as $completed |
        [.efforts_in_progress[].name] as $progress |
        $completed | map(select(. as $item | $progress | index($item))) | 
        .[]' "$STATE_FILE" 2>/dev/null)
    
    if [[ -n "$duplicates" && "$duplicates" != "null" ]]; then
        echo -e "${YELLOW}Found efforts in both completed and in_progress:${NC}"
        for dup in $duplicates; do
            echo "  - $dup"
            # Remove from in_progress since it's completed
            jq --arg name "$dup" '
                .efforts_in_progress |= map(select(.name != $name))
            ' "$STATE_FILE" > tmp.json && mv tmp.json "$STATE_FILE"
            echo -e "${GREEN}  ✅ Removed $dup from in_progress (keeping in completed)${NC}"
        done
    else
        echo -e "${GREEN}No duplicate efforts found${NC}"
    fi
}

# Function to add recovery metadata
add_recovery_metadata() {
    echo -e "${YELLOW}Adding recovery metadata...${NC}"
    
    jq --arg time "$(date -Iseconds)" '
        .recovery_history = (.recovery_history // []) + [{
            "timestamp": $time,
            "recovery_type": "orchestrator_restart",
            "validator": "orchestrator-recovery-validation.sh",
            "actions_taken": []
        }] |
        # Keep only last 5 recovery entries
        .recovery_history |= .[-5:]
    ' "$STATE_FILE" > tmp.json && mv tmp.json "$STATE_FILE"
    
    echo -e "${GREEN}Recovery metadata added${NC}"
}

# Main recovery process
echo -e "${BLUE}Starting recovery validation process...${NC}"
echo ""

# Step 1: Clean orphaned stale tracking
clean_orphaned_stale_tracking

# Step 2: Validate current integrations
validate_current_integrations

# Step 3: Check state validity
check_state_validity

# Step 4: Reconcile effort tracking
reconcile_effort_tracking

# Step 5: Run comprehensive validation
echo ""
echo -e "${YELLOW}Running comprehensive state validation...${NC}"
if [ -f "$PROJECT_DIR/utilities/validate-state-consistency.sh" ]; then
    if "$PROJECT_DIR/utilities/validate-state-consistency.sh" "$STATE_FILE"; then
        echo -e "${GREEN}✅ State validation passed${NC}"
    else
        echo -e "${RED}❌ State validation failed - manual intervention may be needed${NC}"
    fi
else
    echo -e "${YELLOW}Comprehensive validator not found, skipping${NC}"
fi

# Step 6: Add recovery metadata
add_recovery_metadata

# Step 7: Commit recovery changes
echo ""
echo -e "${YELLOW}Committing recovery changes...${NC}"
if git diff --quiet "$STATE_FILE"; then
    echo "No changes needed - state file was already consistent"
else
    git add "$STATE_FILE"
    git commit -m "recovery: R346 state consistency validation and cleanup" \
               -m "- Cleaned orphaned stale tracking" \
               -m "- Fixed integration consistency issues" \
               -m "- Reconciled effort tracking" \
               -m "- Added recovery metadata"
    git push
    echo -e "${GREEN}✅ Recovery changes committed and pushed${NC}"
fi

# Final summary
echo ""
echo -e "${BLUE}═══════════════════════════════════════════════════════════════${NC}"
echo -e "${BLUE}                    RECOVERY SUMMARY                              ${NC}"
echo -e "${BLUE}═══════════════════════════════════════════════════════════════${NC}"

# Get current state info for summary
CURRENT_STATE=$(jq -r '.current_state' "$STATE_FILE")
CURRENT_PHASE=$(jq -r '.current_phase' "$STATE_FILE")
CURRENT_WAVE=$(jq -r '.current_wave' "$STATE_FILE")

echo "Current State: $CURRENT_STATE"
echo "Current Phase: $CURRENT_PHASE"
echo "Current Wave: $CURRENT_WAVE"
echo ""

# Check for any active integrations that might be stale
WAVE_INT_STALE=$(jq -r '.current_wave_integration.is_stale // false' "$STATE_FILE")
if [[ "$WAVE_INT_STALE" == "true" ]]; then
    echo -e "${YELLOW}⚠️ WARNING: Current wave integration is marked as stale${NC}"
    echo "You may need to recreate the integration branch"
fi

# Provide next steps
echo ""
echo -e "${GREEN}Recovery validation complete!${NC}"
echo ""
echo "Next steps:"
echo "1. Review the current state: $CURRENT_STATE"
echo "2. Check for any pending work in this state"
echo "3. Use /continue-orchestrating to resume"

# Check for error recovery state
if [[ "$CURRENT_STATE" == "ERROR_RECOVERY" ]]; then
    echo ""
    echo -e "${YELLOW}⚠️ System is in ERROR_RECOVERY state${NC}"
    echo "Review the error condition before continuing"
fi

echo ""
echo "Recovery log saved to: $RECOVERY_LOG"

exit 0