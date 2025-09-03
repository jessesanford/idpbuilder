#!/bin/bash

# update-orchestrator-state-section.sh
# Helper script to update specific sections of orchestrator state file
# R215: ONLY THE ORCHESTRATOR can use this script!

set -e

# R215 ENFORCEMENT: Check if caller is orchestrator
check_orchestrator_only() {
    # Check if we're being called by the orchestrator
    if [[ "${AGENT_TYPE}" != "orchestrator" ]] && [[ "${FORCE_ORCHESTRATOR}" != "true" ]]; then
        echo "❌ FATAL: R215 VIOLATION - Only orchestrator can update state file!"
        echo "Current agent: ${AGENT_TYPE:-unknown}"
        echo "This script can only be used by the orchestrator."
        echo ""
        echo "Other agents should:"
        echo "1. Write to agent-status.yaml"
        echo "2. Let orchestrator read status and update its state"
        exit 1
    fi
}

# Enforce R215 immediately
check_orchestrator_only

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to update phase architecture plans
update_phase_architecture_plans() {
    local PHASE="$1"
    local STATUS="$2"
    local FILE="$3"
    local STATE_FILE="${4:-orchestrator-state.yaml}"
    
    echo -e "${BLUE}Updating phase_architecture_plans section...${NC}"
    
    # Check if section exists
    if ! grep -q "^phase_architecture_plans:" "$STATE_FILE"; then
        echo -e "${YELLOW}Section doesn't exist, run upgrade script first${NC}"
        exit 1
    fi
    
    # Add new entry (using sed to insert after the section header)
    TIMESTAMP=$(date -Iseconds)
    ENTRY="  - phase: $PHASE\n    status: \"$STATUS\"\n    file: \"$FILE\"\n    created_by: \"architect\"\n    created_at: \"$TIMESTAMP\""
    
    sed -i "/^phase_architecture_plans:/a\\$ENTRY" "$STATE_FILE"
    echo -e "${GREEN}✅ Added phase $PHASE architecture plan${NC}"
}

# Function to update wave architecture plans
update_wave_architecture_plans() {
    local PHASE="$1"
    local WAVE="$2"
    local STATUS="$3"
    local FILE="$4"
    local STATE_FILE="${5:-orchestrator-state.yaml}"
    
    echo -e "${BLUE}Updating wave_architecture_plans section...${NC}"
    
    if ! grep -q "^wave_architecture_plans:" "$STATE_FILE"; then
        echo -e "${YELLOW}Section doesn't exist, run upgrade script first${NC}"
        exit 1
    fi
    
    TIMESTAMP=$(date -Iseconds)
    ENTRY="  - phase: $PHASE\n    wave: $WAVE\n    status: \"$STATUS\"\n    file: \"$FILE\"\n    created_by: \"architect\"\n    created_at: \"$TIMESTAMP\""
    
    sed -i "/^wave_architecture_plans:/a\\$ENTRY" "$STATE_FILE"
    echo -e "${GREEN}✅ Added phase $PHASE wave $WAVE architecture plan${NC}"
}

# Function to update phase implementation plans
update_phase_implementation_plans() {
    local PHASE="$1"
    local STATUS="$2"
    local FILE="$3"
    local SOURCE_ARCH="$4"
    local STATE_FILE="${5:-orchestrator-state.yaml}"
    
    echo -e "${BLUE}Updating phase_implementation_plans section...${NC}"
    
    if ! grep -q "^phase_implementation_plans:" "$STATE_FILE"; then
        echo -e "${YELLOW}Section doesn't exist, run upgrade script first${NC}"
        exit 1
    fi
    
    TIMESTAMP=$(date -Iseconds)
    ENTRY="  - phase: $PHASE\n    status: \"$STATUS\"\n    file: \"$FILE\"\n    created_by: \"code-reviewer\"\n    created_at: \"$TIMESTAMP\"\n    source_architecture: \"$SOURCE_ARCH\""
    
    sed -i "/^phase_implementation_plans:/a\\$ENTRY" "$STATE_FILE"
    echo -e "${GREEN}✅ Added phase $PHASE implementation plan${NC}"
}

# Function to update wave implementation plans
update_wave_implementation_plans() {
    local PHASE="$1"
    local WAVE="$2"
    local STATUS="$3"
    local FILE="$4"
    local SOURCE_ARCH="$5"
    local STATE_FILE="${6:-orchestrator-state.yaml}"
    
    echo -e "${BLUE}Updating wave_implementation_plans section...${NC}"
    
    if ! grep -q "^wave_implementation_plans:" "$STATE_FILE"; then
        echo -e "${YELLOW}Section doesn't exist, run upgrade script first${NC}"
        exit 1
    fi
    
    TIMESTAMP=$(date -Iseconds)
    WAVE_ROOT="efforts/phase${PHASE}/wave${WAVE}"
    ENTRY="  - phase: $PHASE\n    wave: $WAVE\n    status: \"$STATUS\"\n    file: \"$FILE\"\n    created_by: \"code-reviewer\"\n    created_at: \"$TIMESTAMP\"\n    source_architecture: \"$SOURCE_ARCH\"\n    r213_metadata_injected: true\n    wave_root: \"$WAVE_ROOT\""
    
    sed -i "/^wave_implementation_plans:/a\\$ENTRY" "$STATE_FILE"
    echo -e "${GREEN}✅ Added phase $PHASE wave $WAVE implementation plan${NC}"
}

# Function to add wave metadata injection (R213)
add_wave_metadata_injection() {
    local PHASE="$1"
    local WAVE="$2"
    local EFFORT_COUNT="$3"
    local STATE_FILE="${4:-orchestrator-state.yaml}"
    
    echo -e "${BLUE}Adding wave metadata injection (R213)...${NC}"
    
    if ! grep -q "^wave_metadata_injections:" "$STATE_FILE"; then
        echo -e "${YELLOW}Section doesn't exist, run upgrade script first${NC}"
        exit 1
    fi
    
    TIMESTAMP=$(date -Iseconds)
    WAVE_ROOT="efforts/phase${PHASE}/wave${WAVE}"
    ENTRY="  - phase: $PHASE\n    wave: $WAVE\n    metadata_source: \"ORCHESTRATOR\"\n    wave_root: \"$WAVE_ROOT\"\n    effort_count: $EFFORT_COUNT\n    injection_time: \"$TIMESTAMP\""
    
    sed -i "/^wave_metadata_injections:/a\\$ENTRY" "$STATE_FILE"
    echo -e "${GREEN}✅ Added R213 metadata injection for phase $PHASE wave $WAVE${NC}"
}

# Function to add code reviewer wave acknowledgment (R214)
add_code_reviewer_acknowledgment() {
    local PHASE="$1"
    local WAVE="$2"
    local EFFORT_PLANS="$3"
    local STATE_FILE="${4:-orchestrator-state.yaml}"
    
    echo -e "${BLUE}Adding code reviewer wave acknowledgment (R214)...${NC}"
    
    if ! grep -q "^code_reviewer_wave_acknowledgments:" "$STATE_FILE"; then
        echo -e "${YELLOW}Section doesn't exist, run upgrade script first${NC}"
        exit 1
    fi
    
    TIMESTAMP=$(date -Iseconds)
    WAVE_ROOT="efforts/phase${PHASE}/wave${WAVE}"
    ENTRY="  - phase: $PHASE\n    wave: $WAVE\n    acknowledged: true\n    wave_root: \"$WAVE_ROOT\"\n    acknowledgment_time: \"$TIMESTAMP\"\n    effort_plans_created: $EFFORT_PLANS"
    
    sed -i "/^code_reviewer_wave_acknowledgments:/a\\$ENTRY" "$STATE_FILE"
    echo -e "${GREEN}✅ Added R214 acknowledgment for phase $PHASE wave $WAVE${NC}"
}

# Main script logic
case "$1" in
    phase-arch)
        update_phase_architecture_plans "$2" "$3" "$4" "$5"
        ;;
    wave-arch)
        update_wave_architecture_plans "$2" "$3" "$4" "$5" "$6"
        ;;
    phase-impl)
        update_phase_implementation_plans "$2" "$3" "$4" "$5" "$6"
        ;;
    wave-impl)
        update_wave_implementation_plans "$2" "$3" "$4" "$5" "$6" "$7"
        ;;
    wave-metadata)
        add_wave_metadata_injection "$2" "$3" "$4" "$5"
        ;;
    wave-ack)
        add_code_reviewer_acknowledgment "$2" "$3" "$4" "$5"
        ;;
    *)
        echo "Usage: $0 {phase-arch|wave-arch|phase-impl|wave-impl|wave-metadata|wave-ack} [args...]"
        echo ""
        echo "Examples:"
        echo "  $0 phase-arch PHASE STATUS FILE [STATE_FILE]"
        echo "  $0 wave-arch PHASE WAVE STATUS FILE [STATE_FILE]"
        echo "  $0 phase-impl PHASE STATUS FILE SOURCE_ARCH [STATE_FILE]"
        echo "  $0 wave-impl PHASE WAVE STATUS FILE SOURCE_ARCH [STATE_FILE]"
        echo "  $0 wave-metadata PHASE WAVE EFFORT_COUNT [STATE_FILE]"
        echo "  $0 wave-ack PHASE WAVE EFFORT_PLANS [STATE_FILE]"
        exit 1
        ;;
esac