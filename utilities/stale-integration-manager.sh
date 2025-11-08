#!/bin/bash
# Stale Integration Manager - Comprehensive tracking and management

set -e

STATE_FILE="${STATE_FILE:-orchestrator-state-v3.json}"
PROJECT_DIR="${PROJECT_DIR:-/home/vscode/software-factory-template}"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to detect staleness for an integration
detect_staleness() {
    local integration_type="$1"  # wave, phase, or project
    local integration_id="$2"
    
    echo -e "${YELLOW}Checking freshness of ${integration_type} integration: ${integration_id}${NC}"
    
    # Get integration creation time
    local created_at=$(jq -r ".current_${integration_type}_integration.created_at // \"none\"" "$STATE_FILE")
    
    if [[ "$created_at" == "none" || "$created_at" == "null" ]]; then
        echo "No ${integration_type} integration exists yet"
        return 0
    fi
    
    # Find fixes applied after integration creation
    local stale_fixes=$(jq -r --arg created "$created_at" '
        .stale_integration_tracking.fix_tracking.fixes_applied[]? | 
        select(.applied_at > $created) | 
        select(.integrated_into.'$integration_type' == false) |
        .fix_id' "$STATE_FILE" | tr '\n' ' ')
    
    if [[ -n "$stale_fixes" ]]; then
        echo -e "${RED}✗ STALE: Integration created at ${created_at}${NC}"
        echo -e "${RED}  Fixes applied after creation: ${stale_fixes}${NC}"
        
        # Record staleness
        record_staleness "$integration_type" "$integration_id" "$stale_fixes"
        return 1
    else
        echo -e "${GREEN}✓ FRESH: No fixes after creation${NC}"
        
        # Update freshness check timestamp
        jq --arg type "$integration_type" --arg time "$(date -Iseconds)" '
            .["current_" + $type + "_integration"].last_freshness_check = $time |
            .["current_" + $type + "_integration"].is_stale = false
        ' "$STATE_FILE" > tmp.json && mv tmp.json "$STATE_FILE"
        
        return 0
    fi
}

# Function to record staleness in state file
record_staleness() {
    local integration_type="$1"
    local integration_id="$2"
    local stale_fixes="$3"
    
    echo "Recording staleness for ${integration_id}..."
    
    # Get fix details
    local fix_details=$(jq -r --arg fixes "$stale_fixes" '
        .stale_integration_tracking.fix_tracking.fixes_applied[]? |
        select(.fix_id | IN($fixes | split(" ")[])) |
        {fix_id, commit, branch, applied_at, description}
    ' "$STATE_FILE" | jq -s .)
    
    # Update state file with staleness info
    jq --arg type "$integration_type" \
       --arg id "$integration_id" \
       --arg time "$(date -Iseconds)" \
       --argjson fixes "$fix_details" '
        # Mark integration as stale
        .["current_" + $type + "_integration"].is_stale = true |
        .["current_" + $type + "_integration"].staleness_reason = "Fixes applied after integration creation" |
        .["current_" + $type + "_integration"].stale_since = $time |
        .["current_" + $type + "_integration"].stale_due_to_fixes = ($fixes | map(.fix_id)) |
        
        # Add to stale integrations list if not already there
        .stale_integration_tracking.stale_integrations |= 
        if any(.integration_id == $id) then
            map(if .integration_id == $id then
                .became_stale_at = $time |
                .triggering_fixes = $fixes |
                .recreation_completed = false
            else . end)
        else
            . + [{
                "integration_type": $type,
                "integration_id": $id,
                "became_stale_at": $time,
                "stale_reason": "Fixes applied after integration creation",
                "triggering_fixes": $fixes,
                "recreation_required": true,
                "recreation_completed": false
            }]
        end
    ' "$STATE_FILE" > tmp.json && mv tmp.json "$STATE_FILE"
}

# Function to track a fix being applied
track_fix() {
    local fix_id="$1"
    local commit="$2"
    local branch="$3"
    local effort_name="$4"
    local fix_type="$5"  # build_fix, test_fix, review_fix
    local description="$6"
    
    echo "Tracking fix ${fix_id} in ${branch}..."
    
    jq --arg id "$fix_id" \
       --arg commit "$commit" \
       --arg branch "$branch" \
       --arg effort "$effort_name" \
       --arg type "$fix_type" \
       --arg desc "$description" \
       --arg time "$(date -Iseconds)" '
        .stale_integration_tracking.fix_tracking.fixes_applied += [{
            "fix_id": $id,
            "commit": $commit,
            "branch": $branch,
            "effort_name": $effort,
            "applied_at": $time,
            "type": $type,
            "description": $desc,
            "integrated_into": {
                "wave": false,
                "phase": false,
                "project": false
            },
            "made_stale": []
        }]
    ' "$STATE_FILE" > tmp.json && mv tmp.json "$STATE_FILE"
    
    # Check what this fix makes stale
    check_cascade_impact "$branch"
}

# Function to check cascade impact of a fix
check_cascade_impact() {
    local fix_branch="$1"
    
    echo "Checking cascade impact of fix in ${fix_branch}..."
    
    # Determine affected levels from branch name
    local phase=$(echo "$fix_branch" | grep -oP 'phase\K\d+' || echo "")
    local wave=$(echo "$fix_branch" | grep -oP 'wave\K\d+' || echo "")
    
    if [[ -n "$phase" && -n "$wave" ]]; then
        echo -e "${YELLOW}Fix affects: Phase ${phase}, Wave ${wave}${NC}"
        echo -e "${RED}CASCADE REQUIRED: Wave → Phase → Project integrations${NC}"
        
        # Record cascade requirement
        jq --arg branch "$fix_branch" \
           --arg phase "$phase" \
           --arg wave "$wave" \
           --arg time "$(date -Iseconds)" '
            .stale_integration_tracking.staleness_cascade += [{
                "trigger": {
                    "type": "effort_fix",
                    "branch": $branch,
                    "timestamp": $time
                },
                "cascade_sequence": [
                    {
                        "level": "wave",
                        "integration": "phase" + $phase + "-wave" + $wave + "-integration",
                        "became_stale": $time,
                        "must_recreate": true,
                        "recreation_status": "pending"
                    },
                    {
                        "level": "phase",
                        "integration": "phase" + $phase + "-integration",
                        "became_stale": $time,
                        "must_recreate": true,
                        "recreation_status": "pending"
                    },
                    {
                        "level": "project",
                        "integration": "project-integration",
                        "became_stale": $time,
                        "must_recreate": true,
                        "recreation_status": "pending"
                    }
                ],
                "cascade_status": "pending",
                "cascade_started_at": null,
                "cascade_completed_at": null
            }]
        ' "$STATE_FILE" > tmp.json && mv tmp.json "$STATE_FILE"
    fi
}

# Function to mark integration as recreated
mark_recreated() {
    local integration_type="$1"
    local integration_id="$2"
    
    echo -e "${GREEN}Marking ${integration_id} as recreated${NC}"
    
    jq --arg id "$integration_id" \
       --arg type "$integration_type" \
       --arg time "$(date -Iseconds)" '
        # Update stale integration record
        (.stale_integration_tracking.stale_integrations[]? | 
         select(.integration_id == $id)) |= 
         . + {"recreation_completed": true, "recreation_at": $time} |
        
        # Clear staleness from current integration
        .["current_" + $type + "_integration"].is_stale = false |
        .["current_" + $type + "_integration"].staleness_reason = null |
        .["current_" + $type + "_integration"].stale_since = null |
        .["current_" + $type + "_integration"].stale_due_to_fixes = [] |
        .["current_" + $type + "_integration"].created_at = $time |
        
        # Update cascade status
        (.stale_integration_tracking.staleness_cascade[]? | 
         select(.cascade_sequence[].integration == $id) |
         .cascade_sequence[] | select(.integration == $id)) |= 
         . + {"recreation_status": "completed", "recreated_at": $time}
    ' "$STATE_FILE" > tmp.json && mv tmp.json "$STATE_FILE"
}

# Function to generate staleness report
generate_report() {
    echo -e "${YELLOW}Generating Staleness Report...${NC}"
    
    cat > STALENESS-REPORT.md << EOF
# STALE INTEGRATE_WAVE_EFFORTS REPORT
Generated: $(date -Iseconds)

## Currently Stale Integrations
$(jq -r '.stale_integration_tracking.stale_integrations[]? | 
        select(.recreation_completed == false) | 
        "### " + .integration_type + ": " + .integration_id + "\n" +
        "- **Stale since**: " + .became_stale_at + "\n" +
        "- **Reason**: " + .stale_reason + "\n" +
        "- **Triggering fixes**:\n" +
        ((.triggering_fixes[]? | "  - Fix " + .fix_id + " in " + .branch + ": " + .description) // "  None")' \
        "$STATE_FILE")

## Pending Cascades
$(jq -r '.stale_integration_tracking.staleness_cascade[]? | 
        select(.cascade_status != "completed") | 
        "### Cascade from " + .trigger.branch + "\n" +
        "- **Triggered at**: " + .trigger.timestamp + "\n" +
        "- **Status**: " + .cascade_status + "\n" +
        "- **Required recreations**:\n" +
        (.cascade_sequence[]? | "  - " + .level + " (" + .integration + "): " + .recreation_status)' \
        "$STATE_FILE")

## Fixes Tracking
### Applied Fixes
$(jq -r '.stale_integration_tracking.fix_tracking.fixes_applied[]? | 
        "- **" + .fix_id + "** in " + .branch + " at " + .applied_at + "\n" +
        "  - Description: " + .description + "\n" +
        "  - Integrated: Wave=" + (.integrated_into.wave|tostring) + 
        ", Phase=" + (.integrated_into.phase|tostring) + 
        ", Project=" + (.integrated_into.project|tostring)' \
        "$STATE_FILE")

### Pending Integration
$(jq -r '.stale_integration_tracking.fix_tracking.fixes_pending_integration[]? | 
        "- **" + .fix_id + "** in " + .branch + "\n" +
        "  - Pending since: " + .pending_since + "\n" +
        "  - Priority: " + .priority' \
        "$STATE_FILE")

## Validation History (Last 5)
$(jq -r '.stale_integration_tracking.validation_history[-5:][]? | 
        "- **" + .timestamp + "**: " + .integration_type + " " + .integration_id + " = " + .result' \
        "$STATE_FILE")

## Action Required
$(if jq -e '.stale_integration_tracking.stale_integrations[]? | select(.recreation_completed == false)' "$STATE_FILE" > /dev/null; then
    echo "⚠️ **IMMEDIATE ACTION REQUIRED**: Stale integrations detected. Run cascade recreation process."
else
    echo "✅ All integrations are fresh."
fi)
EOF
    
    echo -e "${GREEN}Report generated: STALENESS-REPORT.md${NC}"
}

# Function to check all integrations
check_all() {
    echo -e "${YELLOW}=== COMPREHENSIVE STALENESS CHECK ===${NC}"
    echo "Timestamp: $(date -Iseconds)"
    echo ""
    
    local any_stale=false
    
    for level in wave phase project; do
        # Get integration ID from state
        local integration_id=$(jq -r ".current_${level}_integration.branch // \"none\"" "$STATE_FILE")
        
        if [[ "$integration_id" != "none" && "$integration_id" != "null" ]]; then
            if ! detect_staleness "$level" "$integration_id"; then
                any_stale=true
            fi
        fi
        echo ""
    done
    
    if $any_stale; then
        echo -e "${RED}⚠️ STALE INTEGRATE_WAVE_EFFORTSS DETECTED - CASCADE RECREATION REQUIRED${NC}"
        generate_report
        return 1
    else
        echo -e "${GREEN}✓ All integrations are fresh${NC}"
        return 0
    fi
}

# Main command handler
case "${1:-}" in
    check)
        shift
        if [[ -n "${1:-}" ]]; then
            detect_staleness "$1" "${2:-}"
        else
            check_all
        fi
        ;;
    track-fix)
        track_fix "$2" "$3" "$4" "$5" "$6" "$7"
        ;;
    mark-recreated)
        mark_recreated "$2" "$3"
        ;;
    report)
        generate_report
        ;;
    cascade-check)
        check_cascade_impact "$2"
        ;;
    *)
        echo "Usage: $0 {check [type id]|track-fix|mark-recreated|report|cascade-check}"
        echo ""
        echo "Commands:"
        echo "  check [type id]  - Check staleness (all if no args)"
        echo "  track-fix        - Track a new fix"
        echo "  mark-recreated   - Mark integration as recreated"
        echo "  report           - Generate staleness report"
        echo "  cascade-check    - Check cascade impact of a branch"
        echo ""
        echo "Examples:"
        echo "  $0 check                    # Check all integrations"
        echo "  $0 check wave phase1-wave2  # Check specific integration"
        echo "  $0 track-fix FIX-001 abc123 phase1/wave2/effort3 auth-module build_fix 'Fixed import'"
        echo "  $0 mark-recreated wave phase1-wave2-integration"
        echo "  $0 report                   # Generate STALENESS-REPORT.md"
        exit 1
        ;;
esac