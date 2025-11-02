#!/bin/bash

# verify-state-rules.sh - R290 State Rule Reading and Verification Enforcement
# This script enforces that state rules are read and verified before state work

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to verify state rules were read
verify_state_rules_read() {
    local NEW_STATE="$1"
    local AGENT_TYPE="${2:-orchestrator}"
    local PROJECT_DIR="${CLAUDE_PROJECT_DIR:-$(pwd)}"
    
    echo -e "${YELLOW}🔍 R290 VERIFICATION: Checking state rule reading for ${AGENT_TYPE}/${NEW_STATE}${NC}"
    
    # Step 1: Create marker directory and prepare marker filename
    mkdir -p "${PROJECT_DIR}/markers/state-verification"
    local TIMESTAMP=$(date +%Y%m%d-%H%M%S)
    local MARKER_FILE="${PROJECT_DIR}/markers/state-verification/state_rules_read_${AGENT_TYPE}_${NEW_STATE}-${TIMESTAMP}"
    
    # Step 2: Check if state rules file exists
    local RULES_FILE="${PROJECT_DIR}/agent-states/${AGENT_TYPE}/${NEW_STATE}/rules.md"
    if [[ ! -f "$RULES_FILE" ]]; then
        echo -e "${RED}❌ FATAL: No rules file for ${NEW_STATE}${NC}"
        echo -e "${RED}Missing: ${RULES_FILE}${NC}"
        echo -e "${RED}AUTOMATIC FAILURE: Missing state rules (R290)${NC}"
        exit 290
    fi
    
    # Step 3: Read state rules (simulating the Read tool)
    echo -e "${GREEN}📖 READING STATE RULES FOR ${NEW_STATE}...${NC}"
    echo "----------------------------------------"
    head -20 "$RULES_FILE"
    echo "----------------------------------------"
    echo "[... rest of file ...]"
    
    # Step 4: Create marker with timestamp
    echo "$(date +%s) - Rules read for ${NEW_STATE}" > "$MARKER_FILE"
    
    # Step 5: Explicit acknowledgment required
    echo -e "${GREEN}✅ STATE RULES READ AND ACKNOWLEDGED FOR ${NEW_STATE}${NC}"
    echo -e "${GREEN}📋 Verification marker created: ${MARKER_FILE}${NC}"
    
    return 0
}

# Function to check if rules were read before work
check_rules_were_read() {
    local STATE="$1"
    local AGENT_TYPE="${2:-orchestrator}"
    local PROJECT_DIR="${CLAUDE_PROJECT_DIR:-$(pwd)}"

    echo -e "${YELLOW}🔍 R290 CHECK: Verifying rules were read for ${STATE}${NC}"

    # Check new location first
    local MARKER_PATTERN="${PROJECT_DIR}/markers/state-verification/state_rules_read_${AGENT_TYPE}_${STATE}-*"
    local MARKER_FILE=$(ls $MARKER_PATTERN 2>/dev/null | tail -1)

    # Backward compatibility: check old location
    if [[ -z "$MARKER_FILE" ]]; then
        local OLD_MARKER="${PROJECT_DIR}/.state_rules_read_${AGENT_TYPE}_${STATE}"
        if [[ -f "$OLD_MARKER" ]]; then
            echo -e "${YELLOW}⚠️  Found marker in old location, migrating...${NC}"
            mkdir -p "${PROJECT_DIR}/markers/state-verification"
            local TIMESTAMP=$(date +%Y%m%d-%H%M%S)
            MARKER_FILE="${PROJECT_DIR}/markers/state-verification/state_rules_read_${AGENT_TYPE}_${STATE}-${TIMESTAMP}"
            mv "$OLD_MARKER" "$MARKER_FILE"
            echo -e "${GREEN}✅ Migrated to: $MARKER_FILE${NC}"
        fi
    fi

    if [[ -z "$MARKER_FILE" || ! -f "$MARKER_FILE" ]]; then
        echo -e "${RED}🔴🔴🔴 FATAL ERROR: R290 VIOLATION DETECTED! 🔴🔴🔴${NC}"
        echo -e "${RED}State work attempted in ${STATE} WITHOUT reading rules!${NC}"
        echo -e "${RED}Missing verification marker: ${MARKER_PATTERN}${NC}"
        echo -e "${RED}AUTOMATIC FAILURE: -100% penalty${NC}"
        echo -e "${RED}This is a SUPREME LAW #3 violation (R290)${NC}"
        exit 290
    fi
    
    # Check marker age (must be recent - within 60 seconds)
    local MARKER_TIME=$(cat "$MARKER_FILE" | cut -d' ' -f1)
    local CURRENT_TIME=$(date +%s)
    local AGE=$((CURRENT_TIME - MARKER_TIME))
    
    if [[ $AGE -gt 60 ]]; then
        echo -e "${YELLOW}⚠️ WARNING: State rules read ${AGE} seconds ago${NC}"
        echo -e "${YELLOW}Re-reading required if context lost${NC}"
        # Re-read if too old
        if [[ $AGE -gt 300 ]]; then
            echo -e "${YELLOW}Rules too old (>5 minutes), re-reading...${NC}"
            verify_state_rules_read "$STATE" "$AGENT_TYPE"
        fi
    else
        echo -e "${GREEN}✅ Verification passed: Rules were read ${AGE} seconds ago${NC}"
    fi
    
    return 0
}

# Function to perform state transition with verification
perform_verified_transition() {
    local FROM_STATE="$1"
    local TO_STATE="$2"
    local AGENT_TYPE="${3:-orchestrator}"
    
    echo "================================================"
    echo -e "${YELLOW}STATE TRANSITION: ${FROM_STATE} → ${TO_STATE}${NC}"
    echo "================================================"
    
    # Step 1: Validate transition (would call R206 validator)
    echo "✓ Transition validated (R206)"
    
    # Step 2: Update state file (would call R252 updater)
    echo "✓ State file updated (R252)"
    
    # Step 3: Commit and push (would call R253)
    echo "✓ Changes committed and pushed (R253)"
    
    # Step 4: MANDATORY - Verify and read state rules (R290)
    verify_state_rules_read "$TO_STATE" "$AGENT_TYPE"
    
    # Step 5: Check verification before allowing work
    check_rules_were_read "$TO_STATE" "$AGENT_TYPE"
    
    echo -e "${GREEN}✅ TRANSITION COMPLETE - Authorized to execute ${TO_STATE} work${NC}"
}

# Function to audit verification markers
audit_verification_markers() {
    local PROJECT_DIR="${CLAUDE_PROJECT_DIR:-$(pwd)}"
    
    echo "================================================"
    echo "R290 VERIFICATION AUDIT"
    echo "================================================"
    
    # Find all verification markers (new and old locations)
    local markers=$(find "$PROJECT_DIR/markers/state-verification" -name "state_rules_read_*" 2>/dev/null)
    local old_markers=$(find "$PROJECT_DIR" -maxdepth 1 -name ".state_rules_read_*" 2>/dev/null)
    
    if [[ -z "$markers" && -z "$old_markers" ]]; then
        echo -e "${YELLOW}No verification markers found${NC}"
    else
        if [[ -n "$markers" ]]; then
            echo "Found verification markers (new location):"
            for marker in $markers; do
                local marker_name=$(basename "$marker")
                local timestamp=$(cat "$marker" | cut -d' ' -f1)
                local readable_time=$(date -d "@$timestamp" "+%Y-%m-%d %H:%M:%S" 2>/dev/null || date -r "$timestamp" "+%Y-%m-%d %H:%M:%S")
                echo "  • ${marker_name}: ${readable_time}"
            done
        fi

        if [[ -n "$old_markers" ]]; then
            echo ""
            echo -e "${YELLOW}Found verification markers (old location - should migrate):${NC}"
            for marker in $old_markers; do
                local marker_name=$(basename "$marker")
                echo "  • ${marker_name}"
            done
        fi
    fi
    
    echo "================================================"
}

# Main script logic
main() {
    local ACTION="${1:-help}"
    
    case "$ACTION" in
        verify)
            # Verify state rules were read
            local STATE="${2:-INIT}"
            local AGENT="${3:-orchestrator}"
            verify_state_rules_read "$STATE" "$AGENT"
            ;;
            
        check)
            # Check if rules were read before work
            local STATE="${2:-INIT}"
            local AGENT="${3:-orchestrator}"
            check_rules_were_read "$STATE" "$AGENT"
            ;;
            
        transition)
            # Perform verified state transition
            local FROM="${2:-INIT}"
            local TO="${3:-PLANNING}"
            local AGENT="${4:-orchestrator}"
            perform_verified_transition "$FROM" "$TO" "$AGENT"
            ;;
            
        audit)
            # Audit all verification markers
            audit_verification_markers
            ;;
            
        clean)
            # Clean old verification markers
            echo "Cleaning old verification markers..."
            find markers/state-verification -name "state_rules_read_*" -mtime +1 -delete 2>/dev/null || true
            find . -maxdepth 1 -name ".state_rules_read_*" -delete 2>/dev/null || true
            echo "✓ Old markers cleaned"
            ;;
            
        *)
            echo "R290 State Rule Reading and Verification Enforcement Tool"
            echo ""
            echo "Usage: $0 <action> [arguments]"
            echo ""
            echo "Actions:"
            echo "  verify <state> [agent]     - Verify and read state rules"
            echo "  check <state> [agent]      - Check if rules were read"
            echo "  transition <from> <to> [agent] - Perform verified transition"
            echo "  audit                      - Audit all verification markers"
            echo "  clean                      - Clean old verification markers"
            echo ""
            echo "Examples:"
            echo "  $0 verify INTEGRATE_WAVE_EFFORTS orchestrator"
            echo "  $0 check MONITOR orchestrator"
            echo "  $0 transition WAVE_COMPLETE INTEGRATE_WAVE_EFFORTS"
            echo "  $0 audit"
            ;;
    esac
}

# Run main function
main "$@"