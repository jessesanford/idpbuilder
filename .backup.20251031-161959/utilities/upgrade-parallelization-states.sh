#!/bin/bash

# upgrade-parallelization-states.sh
# Adds new mandatory parallelization analysis states to prevent R151 violations
# Created in response to orchestrator spawning all agents in parallel when some should be sequential

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
BOLD='\033[1m'
NC='\033[0m' # No Color

echo "════════════════════════════════════════════════════════════════"
echo -e "${BOLD}PARALLELIZATION STATES UPGRADE TOOL${NC}"
echo "Adding ANALYZE_CODE_REVIEWER_PARALLELIZATION and"
echo "ANALYZE_IMPLEMENTATION_PARALLELIZATION states"
echo "════════════════════════════════════════════════════════════════"
echo ""

# Function to backup state file
backup_state_file() {
    local STATE_FILE="$1"
    local BACKUP_FILE="${STATE_FILE}.backup.parallelization.$(date +%Y%m%d-%H%M%S)"
    
    echo -e "${BLUE}Creating backup...${NC}"
    cp "$STATE_FILE" "$BACKUP_FILE"
    echo -e "${GREEN}✅ Backup created: $BACKUP_FILE${NC}"
    return 0
}

# Function to check if new fields already exist
check_existing_fields() {
    local STATE_FILE="$1"
    
    if grep -q "code_reviewer_parallelization_plan:" "$STATE_FILE" 2>/dev/null; then
        echo -e "${YELLOW}⚠️  Code Reviewer parallelization plan already exists${NC}"
        return 1
    fi
    
    if grep -q "sw_engineer_parallelization_plan:" "$STATE_FILE" 2>/dev/null; then
        echo -e "${YELLOW}⚠️  SW Engineer parallelization plan already exists${NC}"
        return 1
    fi
    
    return 0
}

# Function to add parallelization sections
add_parallelization_sections() {
    local STATE_FILE="$1"
    local TEMP_FILE="${STATE_FILE}.tmp"
    
    echo -e "${BLUE}Adding parallelization analysis sections...${NC}"
    
    # Create the new sections to add
    cat >> "$STATE_FILE" <<'EOF'

# ════════════════════════════════════════════════════════════════
# PARALLELIZATION ANALYSIS (R151, R218, R219 Compliance)
# ════════════════════════════════════════════════════════════════
# These sections are populated by the new mandatory analysis states:
# - ANALYZE_CODE_REVIEWER_PARALLELIZATION (before spawning reviewers)
# - ANALYZE_IMPLEMENTATION_PARALLELIZATION (before spawning engineers)

# Code Reviewer Parallelization Plan (from wave plan analysis)
code_reviewer_parallelization_plan:
  wave: null
  phase: null
  analysis_timestamp: null
  blocking_efforts: []
  # Example blocking effort:
  # - effort_id: "E3.1.1"
  #   name: "sync-engine-foundation"
  #   reason: "blocks all other efforts"
  #   can_parallelize: false
  
  parallel_groups: {}
  # Example parallel group:
  # group_1:
  #   can_start_after: ["E3.1.1"]
  #   efforts:
  #     - effort_id: "E3.1.2"
  #       name: "webhook-framework"
  
  spawn_sequence: []
  # Example spawn step:
  # - step: 1
  #   action: "spawn_sequential"
  #   efforts: ["E3.1.1"]
  #   wait_for_completion: true

# SW Engineer Parallelization Plan (from implementation plans analysis)
sw_engineer_parallelization_plan:
  wave: null
  phase: null
  analysis_timestamp: null
  blocking_implementations: []
  # Example blocking implementation:
  # - effort_id: "E3.1.1"
  #   name: "sync-engine-foundation"
  #   implementation_plan: "efforts/phase3/wave1/sync-engine-foundation/IMPLEMENTATION-PLAN.md"
  #   can_parallelize: false
  #   reason: "Core foundation - blocks all others"
  #   dependencies: []
  
  parallel_groups: {}
  # Example parallel group:
  # group_1:
  #   can_start_after: ["E3.1.1"]
  #   efforts:
  #     - effort_id: "E3.1.2"
  #       name: "webhook-framework"
  #       implementation_plan: "efforts/phase3/wave1/webhook-framework/IMPLEMENTATION-PLAN.md"
  
  spawn_sequence: []
  # Example spawn step:
  # - step: 1
  #   action: "spawn_sequential"
  #   agent_type: "sw-engineer"
  #   efforts: ["E3.1.1"]
  #   wait_for_completion: true
  #   expected_duration: "2-3 hours"

# Parallelization Violation Tracking (for grading)
parallelization_violations:
  # Tracks any violations of R151, R218, R219
  # Example:
  # - timestamp: "2025-08-26T15:30:00Z"
  #   violation_type: "spawned_blocking_in_parallel"
  #   rule_violated: "R151"
  #   details: "Spawned E3.1.1 with others when it should be sequential"
  #   grade_impact: "-50%"

EOF
    
    echo -e "${GREEN}✅ Parallelization sections added${NC}"
}

# Function to verify the upgrade
verify_upgrade() {
    local STATE_FILE="$1"
    
    echo -e "${BLUE}Verifying upgrade...${NC}"
    
    # Check for new sections
    local checks_passed=0
    local total_checks=2
    
    if grep -q "code_reviewer_parallelization_plan:" "$STATE_FILE"; then
        echo -e "${GREEN}  ✅ Code Reviewer parallelization plan section found${NC}"
        ((checks_passed++))
    else
        echo -e "${RED}  ❌ Code Reviewer parallelization plan section missing${NC}"
    fi
    
    if grep -q "sw_engineer_parallelization_plan:" "$STATE_FILE"; then
        echo -e "${GREEN}  ✅ SW Engineer parallelization plan section found${NC}"
        ((checks_passed++))
    else
        echo -e "${RED}  ❌ SW Engineer parallelization plan section missing${NC}"
    fi
    
    echo ""
    if [ $checks_passed -eq $total_checks ]; then
        echo -e "${GREEN}${BOLD}✅ All verification checks passed!${NC}"
        return 0
    else
        echo -e "${RED}${BOLD}❌ Some verification checks failed${NC}"
        return 1
    fi
}

# Main upgrade logic
main() {
    local STATE_FILE="${1:-orchestrator-state-v3.json}"
    
    # Check if file exists
    if [ ! -f "$STATE_FILE" ]; then
        echo -e "${YELLOW}State file not found: $STATE_FILE${NC}"
        echo "Creating new state file with parallelization sections..."
        
        # Create a new state file with basic structure
        cat > "$STATE_FILE" <<'EOF'
# Orchestrator State File
# Auto-generated with parallelization analysis support

orchestrator_state:
  current_state: "INIT"
  previous_state: null
  transition_time: null
  transition_reason: "Initial state"
  current_phase: null
  current_wave: null

efforts_in_progress: []
efforts_completed: []

EOF
        echo -e "${GREEN}✅ Created new state file${NC}"
    fi
    
    # Create backup
    backup_state_file "$STATE_FILE"
    
    # Check if already upgraded
    if ! check_existing_fields "$STATE_FILE"; then
        echo -e "${YELLOW}State file appears to already have parallelization sections${NC}"
        echo -e "${BLUE}Checking current state...${NC}"
        verify_upgrade "$STATE_FILE"
        echo ""
        echo -e "${CYAN}If you want to force re-add the sections, delete them manually first.${NC}"
        exit 0
    fi
    
    # Add new sections
    add_parallelization_sections "$STATE_FILE"
    
    # Verify upgrade
    echo ""
    verify_upgrade "$STATE_FILE"
    
    echo ""
    echo "════════════════════════════════════════════════════════════════"
    echo -e "${GREEN}${BOLD}UPGRADE COMPLETE!${NC}"
    echo ""
    echo -e "${CYAN}New mandatory states added to prevent parallelization violations:${NC}"
    echo "  1. ANALYZE_CODE_REVIEWER_PARALLELIZATION"
    echo "     - Analyzes wave plan before spawning reviewers"
    echo "     - Enforces R218 compliance"
    echo ""
    echo "  2. ANALYZE_IMPLEMENTATION_PARALLELIZATION"
    echo "     - Analyzes effort plans before spawning engineers"
    echo "     - Enforces R151 compliance"
    echo ""
    echo -e "${YELLOW}IMPORTANT: The orchestrator MUST transition through these states${NC}"
    echo -e "${YELLOW}or face grading penalties for R151/R218 violations!${NC}"
    echo "════════════════════════════════════════════════════════════════"
}

# Run the upgrade
main "$@"