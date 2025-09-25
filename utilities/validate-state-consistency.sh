#!/bin/bash
# R346: State Metadata Synchronization Validator
# R349: Atomic Status Flag Cleanup Validator
# Ensures orchestrator-state.json has no contradictions or inconsistencies

set -e

STATE_FILE="${1:-orchestrator-state.json}"
PROJECT_DIR="${PROJECT_DIR:-/home/vscode/software-factory-template}"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Track validation results
ERRORS=()
WARNINGS=()
FIXED=()

echo -e "${YELLOW}═══════════════════════════════════════════════════════════════${NC}"
echo -e "${YELLOW}          R346: STATE CONSISTENCY VALIDATION                       ${NC}"
echo -e "${YELLOW}          R327/R348: CASCADE VALIDATION                           ${NC}"
echo -e "${YELLOW}═══════════════════════════════════════════════════════════════${NC}"
echo ""

# Function to check stale tracking consistency
check_stale_tracking() {
    echo -e "${YELLOW}Checking stale tracking consistency...${NC}"
    
    # Check for contradictions: recreation_required=true AND recreation_completed=true
    local contradictions=$(jq -r '.stale_integration_tracking.stale_integrations[]? | 
                                  select(.recreation_required == true and .recreation_completed == true) | 
                                  .integration_id' "$STATE_FILE" 2>/dev/null)
    
    if [[ -n "$contradictions" ]]; then
        for id in $contradictions; do
            ERRORS+=("❌ Stale tracking contradiction for $id: both required and completed")
        done
    fi
    
    # Check for orphaned stale tracking (integration completed but still marked stale)
    local orphans=$(jq -r '.stale_integration_tracking.stale_integrations[]? | 
                          select(.recreation_completed == false) | 
                          .integration_id' "$STATE_FILE" 2>/dev/null)
    
    for orphan in $orphans; do
        # Check if this integration is actually complete
        local is_complete=$(jq -r --arg id "$orphan" '
            (.waves_completed | .. | select(.integration_branch? == $id)) // 
            (.phases_completed | .. | select(.integration_branch? == $id)) // 
            "not_found"' "$STATE_FILE")
        
        if [[ "$is_complete" != "not_found" && "$is_complete" != "null" ]]; then
            WARNINGS+=("⚠️ Orphaned stale tracking for completed integration: $orphan")
            FIXED+=("Will clean up stale tracking for $orphan")
            
            # Fix the orphan
            jq --arg id "$orphan" --arg time "$(date -Iseconds)" '
                (.stale_integration_tracking.stale_integrations[]? | 
                 select(.integration_id == $id)) |= 
                 . + {
                    "recreation_completed": true,
                    "recreation_at": $time,
                    "recreation_required": false
                 }
            ' "$STATE_FILE" > tmp.json && mv tmp.json "$STATE_FILE"
        fi
    done
}

# Function to check current integration consistency
check_current_integrations() {
    echo -e "${YELLOW}Checking current integration consistency...${NC}"
    
    for type in wave phase project; do
        # Check if is_stale=false but has staleness_reason
        local inconsistent=$(jq -r --arg type "$type" '
            .["current_" + $type + "_integration"] | 
            select(.is_stale == false and .staleness_reason != null) | 
            "found"' "$STATE_FILE" 2>/dev/null)
        
        if [[ "$inconsistent" == "found" ]]; then
            ERRORS+=("❌ Current $type integration: not stale but has staleness reason")
        fi
        
        # Check if is_stale=true but no staleness_reason
        inconsistent=$(jq -r --arg type "$type" '
            .["current_" + $type + "_integration"] | 
            select(.is_stale == true and .staleness_reason == null) | 
            "found"' "$STATE_FILE" 2>/dev/null)
        
        if [[ "$inconsistent" == "found" ]]; then
            ERRORS+=("❌ Current $type integration: marked stale but no reason given")
        fi
        
        # Check if stale_due_to_fixes is empty but is_stale=true
        inconsistent=$(jq -r --arg type "$type" '
            .["current_" + $type + "_integration"] | 
            select(.is_stale == true and (.stale_due_to_fixes | length) == 0) | 
            "found"' "$STATE_FILE" 2>/dev/null)
        
        if [[ "$inconsistent" == "found" ]]; then
            WARNINGS+=("⚠️ Current $type integration: stale but no triggering fixes listed")
        fi
    done
}

# Function to check integration status consistency
check_integration_status() {
    echo -e "${YELLOW}Checking integration status consistency...${NC}"
    
    # Get all integration status entries
    local statuses=$(jq -r '.integration_status | keys[]?' "$STATE_FILE" 2>/dev/null)
    
    for status_key in $statuses; do
        local status=$(jq -r ".integration_status[\"$status_key\"].status" "$STATE_FILE")
        local integration_id="${status_key#*_}" # Remove prefix
        
        if [[ "$status" == "COMPLETE" ]]; then
            # Check if still marked as stale
            local still_stale=$(jq -r --arg id "$integration_id" '
                .stale_integration_tracking.stale_integrations[]? | 
                select(.integration_id == $id and .recreation_completed == false) | 
                "found"' "$STATE_FILE" 2>/dev/null)
            
            if [[ "$still_stale" == "found" ]]; then
                ERRORS+=("❌ Integration $integration_id marked COMPLETE but still in stale tracking")
            fi
        fi
    done
}

# Function to check effort completion consistency
check_effort_completion() {
    echo -e "${YELLOW}Checking effort completion consistency...${NC}"
    
    # Check efforts_completed for consistency
    local completed_count=$(jq '.efforts_completed | length' "$STATE_FILE" 2>/dev/null)
    local in_progress_count=$(jq '.efforts_in_progress | length' "$STATE_FILE" 2>/dev/null)
    
    # Check for duplicates between completed and in_progress
    local duplicates=$(jq -r '
        [.efforts_completed[].name] as $completed |
        [.efforts_in_progress[].name] as $progress |
        $completed | map(select(. as $item | $progress | index($item))) | 
        join(", ")' "$STATE_FILE" 2>/dev/null)
    
    if [[ -n "$duplicates" && "$duplicates" != "null" ]]; then
        ERRORS+=("❌ Efforts in both completed and in_progress: $duplicates")
    fi
}

# Function to check metadata locations
check_metadata_locations() {
    echo -e "${YELLOW}Checking metadata location consistency...${NC}"
    
    # Check if metadata locations point to existing integrations
    local metadata_integrations=$(jq -r '.metadata_locations.integration_reports | keys[]?' "$STATE_FILE" 2>/dev/null)
    
    for meta_key in $metadata_integrations; do
        local file_path=$(jq -r ".metadata_locations.integration_reports[\"$meta_key\"].file_path // \"none\"" "$STATE_FILE")
        
        if [[ "$file_path" != "none" && "$file_path" != "null" ]]; then
            # Check if corresponding integration exists
            local integration_exists=$(jq -r --arg key "$meta_key" '
                (.current_wave_integration.branch | contains($key)) or
                (.current_phase_integration.branch | contains($key)) or
                (.integration_branches[].branch | contains($key)) | 
                any' "$STATE_FILE" 2>/dev/null)
            
            if [[ "$integration_exists" == "false" ]]; then
                WARNINGS+=("⚠️ Metadata location for non-existent integration: $meta_key")
            fi
        fi
    done
}

# Function to check fix tracking consistency
check_fix_tracking() {
    echo -e "${YELLOW}Checking fix tracking consistency...${NC}"
    
    # Check fixes marked as integrated but integrations are stale
    local fixes=$(jq -r '.stale_integration_tracking.fix_tracking.fixes_applied[]?' "$STATE_FILE" 2>/dev/null)
    
    if [[ -n "$fixes" ]]; then
        jq -c '.stale_integration_tracking.fix_tracking.fixes_applied[]?' "$STATE_FILE" 2>/dev/null | while read -r fix; do
            local fix_id=$(echo "$fix" | jq -r '.fix_id')
            local wave_integrated=$(echo "$fix" | jq -r '.integrated_into.wave')
            local phase_integrated=$(echo "$fix" | jq -r '.integrated_into.phase')
            
            # Check if marked as integrated but integration is stale
            if [[ "$wave_integrated" == "true" ]]; then
                local wave_stale=$(jq -r '.current_wave_integration.is_stale' "$STATE_FILE")
                if [[ "$wave_stale" == "true" ]]; then
                    WARNINGS+=("⚠️ Fix $fix_id marked as wave-integrated but wave is stale")
                fi
            fi
            
            if [[ "$phase_integrated" == "true" ]]; then
                local phase_stale=$(jq -r '.current_phase_integration.is_stale' "$STATE_FILE")
                if [[ "$phase_stale" == "true" ]]; then
                    WARNINGS+=("⚠️ Fix $fix_id marked as phase-integrated but phase is stale")
                fi
            fi
        done
    fi
}

# Function to validate state transitions
check_state_transitions() {
    echo -e "${YELLOW}Checking state transition consistency...${NC}"
    
    # Check current_state vs previous_state
    local current=$(jq -r '.current_state' "$STATE_FILE" 2>/dev/null)
    local previous=$(jq -r '.previous_state' "$STATE_FILE" 2>/dev/null)
    
    if [[ "$current" == "$previous" ]]; then
        WARNINGS+=("⚠️ current_state same as previous_state: $current")
    fi
    
    # Check for valid state names
    local valid_states=(
        "INIT" "PLANNING" "SETUP_EFFORT_INFRASTRUCTURE" 
        "SPAWN_AGENTS" "MONITOR" "WAVE_COMPLETE" 
        "INTEGRATION" "MONITORING_INTEGRATION" 
        "INTEGRATION_CODE_REVIEW" "INTEGRATION_FEEDBACK_REVIEW"
        "ERROR_RECOVERY" "WAVE_START" "WAVE_REVIEW"
    )
    
    if [[ ! " ${valid_states[@]} " =~ " ${current} " ]]; then
        ERRORS+=("❌ Invalid current_state: $current")
    fi
}

# R349: Function to check for stale status flags
check_stale_status_flags() {
    echo -e "${YELLOW}R349: Checking for stale status flags...${NC}"
    
    # Check for completed rebases with stale flags
    local stale_rebase=$(jq -r '.. | objects | select(
        has("completed_at") and has("needs_rebase") and
        .completed_at != null and .needs_rebase == true
    ) | "found"' "$STATE_FILE" 2>/dev/null | head -1)
    
    if [[ "$stale_rebase" == "found" ]]; then
        ERRORS+=("❌ R349 VIOLATION: Rebase completed but needs_rebase flag still true")
        
        # Try to identify which effort/section has the issue
        local problem_sections=$(jq -r 'paths(objects) as $p | 
            getpath($p) | 
            select(type == "object" and has("completed_at") and has("needs_rebase") and
                   .completed_at != null and .needs_rebase == true) | 
            $p | join(".")' "$STATE_FILE" 2>/dev/null)
        
        if [[ -n "$problem_sections" ]]; then
            ERRORS+=("   Problem in: $problem_sections")
        fi
    fi
    
    # Check for completed splits with stale flags
    local stale_splits=$(jq -r '.violations[]? | select(
        has("split_completed") and .split_completed == true and 
        .requires_split == true
    ) | .effort' "$STATE_FILE" 2>/dev/null)
    
    if [[ -n "$stale_splits" ]]; then
        for effort in $stale_splits; do
            ERRORS+=("❌ R349 VIOLATION: Split completed but requires_split still true for: $effort")
        done
    fi
    
    # Check for contradictory staleness flags
    local stale_contradiction=$(jq -r '.. | objects | select(
        has("is_stale") and has("staleness_reason") and
        .is_stale == false and .staleness_reason != null and
        .staleness_reason != ""
    ) | "found"' "$STATE_FILE" 2>/dev/null | head -1)
    
    if [[ "$stale_contradiction" == "found" ]]; then
        ERRORS+=("❌ R349 VIOLATION: Integration not stale but has staleness_reason")
    fi
    
    # Check for contradictory recreation flags
    local recreation_contradiction=$(jq -r '.stale_integration_tracking.stale_integrations[]? | select(
        .recreation_completed == true and .recreation_required == true
    ) | .integration_id' "$STATE_FILE" 2>/dev/null)
    
    if [[ -n "$recreation_contradiction" ]]; then
        for id in $recreation_contradiction; do
            ERRORS+=("❌ R349 VIOLATION: Recreation both required and completed for: $id")
        done
    fi
    
    # Check base_branch_tracking for stale rebase flags
    local stale_base_tracking=$(jq -r '.. | .base_branch_tracking? | select(. != null) | select(
        .last_rebase != null and .requires_rebase == true
    ) | "found"' "$STATE_FILE" 2>/dev/null | head -1)
    
    if [[ "$stale_base_tracking" == "found" ]]; then
        ERRORS+=("❌ R349 VIOLATION: base_branch_tracking has last_rebase but requires_rebase still true")
    fi
    
    # Check for efforts with contradictory fix flags
    local stale_fixes=$(jq -r '.. | objects | select(
        has("has_fixes_applied") and has("needs_fix") and
        .has_fixes_applied == true and .needs_fix == true
    ) | .name // "unknown"' "$STATE_FILE" 2>/dev/null)
    
    if [[ -n "$stale_fixes" && "$stale_fixes" != "null" ]]; then
        for effort in $stale_fixes; do
            ERRORS+=("❌ R349 VIOLATION: Fixes applied but needs_fix still true for: $effort")
        done
    fi
    
    # Check for stale split status in split_tracking
    local stale_split_status=$(jq -r '.split_tracking | to_entries[]? | select(
        .value.status == "COMPLETED" and 
        (.value.splits[]? | select(.status != "COMPLETED") | .number)
    ) | .key' "$STATE_FILE" 2>/dev/null)
    
    if [[ -n "$stale_split_status" ]]; then
        for effort in $stale_split_status; do
            WARNINGS+=("⚠️ R349 WARNING: Split tracking shows COMPLETED but has incomplete splits: $effort")
        done
    fi
}

# R327/R348 CASCADE validation
check_cascade_requirements() {
    echo -e "${YELLOW}Checking R327/R348 CASCADE requirements...${NC}"
    
    # Check if in CASCADE_REINTEGRATION state
    local current_state=$(jq -r '.current_state' "$STATE_FILE" 2>/dev/null)
    
    if [[ "$current_state" == "CASCADE_REINTEGRATION" ]]; then
        # Must have pending cascades
        local pending_cascades=$(jq -r '.stale_integration_tracking.staleness_cascade[]? | 
                                       select(.cascade_status != "completed")' "$STATE_FILE" 2>/dev/null | wc -l)
        
        if [[ "$pending_cascades" -eq 0 ]]; then
            ERRORS+=("❌ R327 VIOLATION: In CASCADE_REINTEGRATION but no pending cascades")
        fi
        
        # Check dependency graph exists
        local has_graph=$(jq -r '.dependency_graph | if . then "yes" else "no" end' "$STATE_FILE" 2>/dev/null)
        if [[ "$has_graph" != "yes" ]]; then
            WARNINGS+=("⚠️ R350 WARNING: CASCADE_REINTEGRATION state but no dependency graph")
        fi
    fi
    
    # Check for incomplete cascades
    local incomplete_cascades=$(jq -r '.dependency_graph.cascade_chains[]? | 
                                      select(.cascade_status == "in_progress" or .cascade_status == "pending") | 
                                      .cascade_id' "$STATE_FILE" 2>/dev/null)
    
    if [[ -n "$incomplete_cascades" ]]; then
        for cascade in $incomplete_cascades; do
            # Check if this cascade is stale (older than 1 hour)
            local triggered_at=$(jq -r --arg id "$cascade" '.dependency_graph.cascade_chains[]? | 
                                       select(.cascade_id == $id) | 
                                       .triggered_at' "$STATE_FILE" 2>/dev/null)
            
            if [[ -n "$triggered_at" && "$triggered_at" != "null" ]]; then
                local age_seconds=$(($(date +%s) - $(date -d "$triggered_at" +%s 2>/dev/null || echo 0)))
                if [[ "$age_seconds" -gt 3600 ]]; then
                    WARNINGS+=("⚠️ R327 WARNING: Cascade $cascade has been pending for over 1 hour")
                fi
            fi
        done
    fi
    
    # Check for stale integrations without cascade tracking
    local untracked_stale=$(jq -r '.stale_integration_tracking.stale_integrations[]? | 
                                  select(.recreation_required == true and .recreation_completed == false) | 
                                  .integration_id' "$STATE_FILE" 2>/dev/null)
    
    for stale in $untracked_stale; do
        # Check if this has a cascade chain
        local has_cascade=$(jq -r --arg id "$stale" '.dependency_graph.cascade_chains[]? | 
                                  select(.operations[]?.target == $id) | 
                                  "found"' "$STATE_FILE" 2>/dev/null)
        
        if [[ "$has_cascade" != "found" ]]; then
            ERRORS+=("❌ R327 VIOLATION: Stale integration $stale has no cascade chain")
        fi
    done
    
    # Check for fixes without cascade triggers
    local recent_fixes=$(jq -r '.stale_integration_tracking.fix_tracking.fixes_applied[]? | 
                               select(.integrated_into.wave == false or 
                                      .integrated_into.phase == false or 
                                      .integrated_into.project == false) | 
                               .fix_id' "$STATE_FILE" 2>/dev/null)
    
    if [[ -n "$recent_fixes" ]]; then
        for fix in $recent_fixes; do
            WARNINGS+=("⚠️ R327 WARNING: Fix $fix not fully integrated (needs cascade)")
        done
    fi
}

# Run all checks
check_stale_tracking
check_current_integrations
check_integration_status
check_effort_completion
check_metadata_locations
check_fix_tracking
check_state_transitions
check_stale_status_flags  # R349 validation
check_cascade_requirements  # R327/R348 CASCADE validation

# Report results
echo ""
echo -e "${YELLOW}═══════════════════════════════════════════════════════════════${NC}"
echo -e "${YELLOW}                     VALIDATION RESULTS                           ${NC}"
echo -e "${YELLOW}═══════════════════════════════════════════════════════════════${NC}"
echo ""

# Report errors
if [ ${#ERRORS[@]} -gt 0 ]; then
    echo -e "${RED}ERRORS FOUND (${#ERRORS[@]})${NC}"
    for error in "${ERRORS[@]}"; do
        echo -e "$error"
    done
    echo ""
fi

# Report warnings
if [ ${#WARNINGS[@]} -gt 0 ]; then
    echo -e "${YELLOW}WARNINGS FOUND (${#WARNINGS[@]})${NC}"
    for warning in "${WARNINGS[@]}"; do
        echo -e "$warning"
    done
    echo ""
fi

# Report fixes
if [ ${#FIXED[@]} -gt 0 ]; then
    echo -e "${GREEN}AUTO-FIXED (${#FIXED[@]})${NC}"
    for fix in "${FIXED[@]}"; do
        echo -e "✅ $fix"
    done
    echo ""
fi

# Final result
if [ ${#ERRORS[@]} -eq 0 ]; then
    if [ ${#WARNINGS[@]} -eq 0 ]; then
        echo -e "${GREEN}✅ STATE IS FULLY CONSISTENT - NO ISSUES FOUND${NC}"
    else
        echo -e "${YELLOW}⚠️ STATE IS MOSTLY CONSISTENT - WARNINGS ONLY${NC}"
    fi
    
    # Add validation timestamp to state file
    jq --arg time "$(date -Iseconds)" '
        .stale_integration_tracking.validation_history += [{
            "timestamp": $time,
            "result": "validated",
            "errors": 0,
            "warnings": '${#WARNINGS[@]}',
            "auto_fixed": '${#FIXED[@]}',
            "validator": "validate-state-consistency.sh"
        }] |
        # Keep only last 10 validation entries
        .stale_integration_tracking.validation_history |= .[-10:]
    ' "$STATE_FILE" > tmp.json && mv tmp.json "$STATE_FILE"
    
    echo ""
    echo -e "${GREEN}Validation recorded in state file${NC}"
    exit 0
else
    echo -e "${RED}❌ STATE INCONSISTENCY DETECTED - R346 VIOLATION!${NC}"
    echo -e "${RED}Cannot proceed with ${#ERRORS[@]} critical errors${NC}"
    
    # Add failed validation to history
    jq --arg time "$(date -Iseconds)" '
        .stale_integration_tracking.validation_history += [{
            "timestamp": $time,
            "result": "failed",
            "errors": '${#ERRORS[@]}',
            "warnings": '${#WARNINGS[@]}',
            "auto_fixed": '${#FIXED[@]}',
            "validator": "validate-state-consistency.sh"
        }] |
        # Keep only last 10 validation entries
        .stale_integration_tracking.validation_history |= .[-10:]
    ' "$STATE_FILE" > tmp.json && mv tmp.json "$STATE_FILE"
    
    exit 346
fi