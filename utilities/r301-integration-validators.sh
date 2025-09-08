#!/bin/bash
# R301 Integration Branch Current Tracking Validators
# SUPREME LAW: Only ONE current integration allowed per phase

# Validate that we're using the current integration branch
validate_using_current_integration() {
    local PHASE=$1
    local BRANCH_TO_USE=$2
    
    # Get the CURRENT integration branch
    CURRENT=$(yq '.current_integration | select(.phase == env(PHASE)).branch' orchestrator-state.json)
    CURRENT_STATUS=$(yq '.current_integration | select(.phase == env(PHASE)).status' orchestrator-state.json)
    
    # FATAL if not using current
    if [[ "$BRANCH_TO_USE" != "$CURRENT" ]]; then
        echo "🔴🔴🔴 SUPREME LAW VIOLATION: R301 - Using deprecated integration branch!"
        echo "  Current: $CURRENT (status: $CURRENT_STATUS)"
        echo "  Attempted: $BRANCH_TO_USE"
        echo "  PENALTY: -100% AUTOMATIC FAILURE"
        return 1
    fi
    
    # FATAL if current is not active
    if [[ "$CURRENT_STATUS" != "active" ]]; then
        echo "🔴🔴🔴 FATAL: Current integration is not active!"
        echo "  Branch: $CURRENT"
        echo "  Status: $CURRENT_STATUS"
        return 1
    fi
    
    echo "✅ Using current integration: $CURRENT"
    return 0
}

# Get the current integration branch for a phase
get_current_integration_branch() {
    local PHASE=$1
    
    CURRENT=$(yq '.current_integration | select(.phase == env(PHASE)).branch' orchestrator-state.json)
    STATUS=$(yq '.current_integration | select(.phase == env(PHASE)).status' orchestrator-state.json)
    
    if [ -z "$CURRENT" ]; then
        echo "❌ ERROR: No current integration for phase $PHASE" >&2
        return 1
    fi
    
    if [ "$STATUS" != "active" ]; then
        echo "❌ ERROR: Current integration is not active (status: $STATUS)" >&2
        return 1
    fi
    
    echo "$CURRENT"
    return 0
}

# Deprecate old integration and set new current
set_current_integration() {
    local PHASE=$1
    local NEW_BRANCH=$2
    local INTEGRATION_TYPE=$3  # "initial", "post_fixes", "wave", "phase"
    local REASON=${4:-"New integration created"}
    
    echo "📝 Updating current_integration per R301..."
    
    # First, move any existing current to deprecated
    EXISTING=$(yq '.current_integration | select(.phase == env(PHASE))' orchestrator-state.json)
    if [ "$EXISTING" != "" ] && [ "$EXISTING" != "null" ]; then
        echo "  Deprecating old integration..."
        
        # Add to deprecated list with updated fields
        yq -i '.deprecated_integrations += [{
            "phase": env(PHASE),
            "branch": .current_integration.branch,
            "status": "deprecated",
            "deprecated_at": "'$(date -u +%Y-%m-%dT%H:%M:%SZ)'",
            "reason": env(REASON)
        }]' orchestrator-state.json
    fi
    
    # Set the new current integration
    echo "  Setting new current: $NEW_BRANCH"
    yq -i '.current_integration = {
        "phase": env(PHASE),
        "branch": env(NEW_BRANCH),
        "status": "active",
        "created_at": "'$(date -u +%Y-%m-%dT%H:%M:%SZ)'",
        "type": env(INTEGRATION_TYPE)
    }' orchestrator-state.json
    
    echo "✅ Current integration updated: $NEW_BRANCH"
    return 0
}

# Validate no multiple active integrations exist
validate_single_active_integration() {
    local PHASE=$1
    
    # Count active integrations for this phase
    ACTIVE_COUNT=$(yq '.current_integration | select(.phase == env(PHASE) and .status == "active")' orchestrator-state.json | grep -c "phase:")
    
    if [ "$ACTIVE_COUNT" -gt 1 ]; then
        echo "🔴🔴🔴 FATAL: Multiple active integrations detected for phase $PHASE!"
        echo "  Count: $ACTIVE_COUNT"
        echo "  R301 VIOLATION: Only ONE current integration allowed"
        return 1
    fi
    
    if [ "$ACTIVE_COUNT" -eq 0 ]; then
        echo "⚠️ WARNING: No active integration for phase $PHASE"
        return 1
    fi
    
    echo "✅ Single active integration confirmed for phase $PHASE"
    return 0
}

# List all deprecated integrations for a phase
list_deprecated_integrations() {
    local PHASE=$1
    
    echo "📋 Deprecated integrations for phase $PHASE:"
    yq '.deprecated_integrations[] | select(.phase == env(PHASE))' orchestrator-state.json | \
        yq '{branch: .branch, deprecated_at: .deprecated_at, reason: .reason}'
}

# Check if a branch is deprecated
is_branch_deprecated() {
    local BRANCH=$1
    
    DEPRECATED=$(yq '.deprecated_integrations[] | select(.branch == env(BRANCH))' orchestrator-state.json)
    
    if [ -n "$DEPRECATED" ]; then
        echo "⚠️ Branch '$BRANCH' is DEPRECATED"
        echo "$DEPRECATED" | yq '.reason'
        return 0  # Yes, it's deprecated
    fi
    
    return 1  # Not deprecated
}

# Migrate from old phase_integration_branches to new structure
migrate_to_r301_structure() {
    local PHASE=$1
    
    echo "🔄 Migrating to R301 structure for phase $PHASE..."
    
    # Check if already migrated
    if [ "$(yq '.current_integration' orchestrator-state.json)" != "null" ]; then
        echo "  Already using R301 structure"
        return 0
    fi
    
    # Find the most recent integration branch from old structure
    LATEST_BRANCH=$(yq ".phase_integration_branches[] | select(.phase == $PHASE)" orchestrator-state.json | \
                    yq '.branch' | tail -1)
    
    if [ -z "$LATEST_BRANCH" ]; then
        echo "  No existing integration branches to migrate"
        return 0
    fi
    
    echo "  Latest branch found: $LATEST_BRANCH"
    
    # Set as current
    yq -i '.current_integration = {
        "phase": env(PHASE),
        "branch": env(LATEST_BRANCH),
        "status": "active",
        "created_at": "'$(date -u +%Y-%m-%dT%H:%M:%SZ)'",
        "type": "migrated"
    }' orchestrator-state.json
    
    # Move others to deprecated
    yq -i '.deprecated_integrations = (.phase_integration_branches[] | 
           select(.branch != env(LATEST_BRANCH)) | 
           . + {"status": "deprecated", "deprecated_at": "'$(date -u +%Y-%m-%dT%H:%M:%SZ)'", "reason": "migrated from old structure"})' \
           orchestrator-state.json
    
    # Remove old structure
    yq -i 'del(.phase_integration_branches)' orchestrator-state.json
    
    echo "✅ Migration complete"
    return 0
}

# Pre-architect validation
validate_before_architect_assessment() {
    local PHASE=$1
    
    echo "🔍 R301 Pre-Architect Validation for Phase $PHASE"
    
    # Check current integration exists
    CURRENT=$(get_current_integration_branch $PHASE)
    if [ $? -ne 0 ]; then
        echo "❌ No current integration branch!"
        return 1
    fi
    
    # Validate it's active
    validate_single_active_integration $PHASE
    if [ $? -ne 0 ]; then
        return 1
    fi
    
    # Check branch exists in git
    if ! git ls-remote --heads origin "$CURRENT" | grep -q "$CURRENT"; then
        echo "❌ Current integration branch not found on remote: $CURRENT"
        return 1
    fi
    
    echo "✅ Ready for architect assessment of: $CURRENT"
    return 0
}

# Export functions for use in other scripts
export -f validate_using_current_integration
export -f get_current_integration_branch
export -f set_current_integration
export -f validate_single_active_integration
export -f list_deprecated_integrations
export -f is_branch_deprecated
export -f migrate_to_r301_structure
export -f validate_before_architect_assessment

# If sourced with arguments, run the requested function
if [ $# -gt 0 ]; then
    "$@"
fi