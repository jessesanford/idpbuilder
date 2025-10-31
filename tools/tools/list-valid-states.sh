#!/bin/bash

# Software Factory 2.0 - List Valid States Tool
# Lists all valid states from the state machine for easy reference
# Useful when validation fails and you need to see valid options

set -e

# Color codes for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
BOLD='\033[1m'
NC='\033[0m' # No Color

# Get the project root (handle both direct and symlinked calls)
if [[ -n "$CLAUDE_PROJECT_DIR" ]]; then
    PROJECT_ROOT="$CLAUDE_PROJECT_DIR"
else
    # Try to find project root from script location
    SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
    PROJECT_ROOT="$( cd "$SCRIPT_DIR/.." && pwd )"
fi

STATE_MACHINE="$PROJECT_ROOT/state-machines/software-factory-3.0-state-machine.json"

# Check if state machine exists
if [ ! -f "$STATE_MACHINE" ]; then
    echo -e "${RED}❌ State machine file not found at: $STATE_MACHINE${NC}" >&2
    echo "Please ensure you're in a Software Factory 2.0 project directory" >&2
    exit 1
fi

# Function to print header
print_header() {
    echo ""
    echo -e "${BOLD}${BLUE}═══════════════════════════════════════════════════════════════${NC}"
    echo -e "${BOLD}${BLUE}         SOFTWARE FACTORY 2.0 - VALID STATES REFERENCE         ${NC}"
    echo -e "${BOLD}${BLUE}═══════════════════════════════════════════════════════════════${NC}"
    echo ""
}

# Function to print agent section
print_agent_states() {
    local agent="$1"
    local color="$2"

    echo -e "${BOLD}${color}📋 ${agent^^} STATES:${NC}"
    echo -e "${color}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"

    # Extract and group states
    local states=$(jq -r ".transition_matrix.\"$agent\" | keys | .[]" "$STATE_MACHINE" 2>/dev/null | sort)

    if [ -z "$states" ]; then
        echo -e "  ${RED}No states found for $agent${NC}"
        return
    fi

    # Group states by prefix for better readability
    local current_prefix=""
    local group_started=false

    while IFS= read -r state; do
        # Get prefix (first part before underscore)
        local prefix=$(echo "$state" | cut -d'_' -f1)

        if [ "$prefix" != "$current_prefix" ]; then
            if [ "$group_started" = true ]; then
                echo "" # Add spacing between groups
            fi
            current_prefix="$prefix"
            group_started=true
            echo -e "\n  ${BOLD}$prefix States:${NC}"
        fi

        # Print the state with proper indentation
        echo -e "    ${color}•${NC} $state"
    done <<< "$states"

    echo ""
}

# Function to print common state patterns
print_common_patterns() {
    echo -e "${BOLD}${CYAN}🔍 COMMON STATE PATTERNS:${NC}"
    echo -e "${CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo ""
    echo -e "  ${BOLD}Initialization:${NC}"
    echo "    • INIT - Initial state for all agents"
    echo ""
    echo -e "  ${BOLD}Planning & Setup:${NC}"
    echo "    • PLANNING - Creating implementation plans"
    echo "    • SETUP_* - Infrastructure and environment setup"
    echo "    • CREATE_* - Creating branches, structures, etc."
    echo ""
    echo -e "  ${BOLD}Execution:${NC}"
    echo "    • SPAWN_* - Spawning other agents"
    echo "    • IMPLEMENTATION - Active development work"
    echo "    • MONITORING_* - Tracking progress"
    echo "    • WAITING_* - Awaiting results"
    echo ""
    echo -e "  ${BOLD}Review & Validation:${NC}"
    echo "    • CODE_REVIEW - Reviewing code changes"
    echo "    • VALIDATION - Validating completeness"
    echo "    • INTEGRATE_WAVE_EFFORTS_* - Integration activities"
    echo ""
    echo -e "  ${BOLD}Completion:${NC}"
    echo "    • *_COMPLETE - Task/wave/phase completion"
    echo "    • PROJECT_DONE - Successful completion"
    echo "    • HANDOFF - Transitioning to next stage"
    echo ""
    echo -e "  ${BOLD}Error Handling:${NC}"
    echo "    • ERROR_* - Error states"
    echo "    • FIX_* - Fixing issues"
    echo "    • RECOVERY - Recovery procedures"
    echo ""
}

# Function to print usage tips
print_usage_tips() {
    echo -e "${BOLD}${YELLOW}💡 USAGE TIPS:${NC}"
    echo -e "${YELLOW}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo ""
    echo "  1. When validation fails, check the current_state and previous_state"
    echo "  2. Ensure states match EXACTLY (case-sensitive)"
    echo "  3. Use states appropriate for your agent type"
    echo "  4. Follow the state machine transitions (don't skip states)"
    echo ""
    echo -e "  ${BOLD}To validate your state file:${NC}"
    echo -e "    ${GREEN}bash \$CLAUDE_PROJECT_DIR/tools/validate-state-embedded.sh orchestrator-state-v3.json${NC}"
    echo ""
    echo -e "  ${BOLD}To see state transitions:${NC}"
    echo -e "    ${GREEN}jq '.transition_matrix.orchestrator.YOUR_STATE' $STATE_MACHINE${NC}"
    echo ""
}

# Main execution
main() {
    print_header

    # Print states for each agent
    print_agent_states "orchestrator" "$GREEN"
    print_agent_states "sw-engineer" "$BLUE"
    print_agent_states "code-reviewer" "$YELLOW"
    print_agent_states "architect" "$CYAN"

    print_common_patterns
    print_usage_tips

    echo -e "${BOLD}${BLUE}═══════════════════════════════════════════════════════════════${NC}"
    echo ""
}

# Run main
main