#!/bin/bash

# upgrade-orchestrator-state.sh
# Upgrades orchestrator state file to include:
# - Architecture planning states (R210-R214)
# - Parallelization analysis states (R231-R233)
# - Phase assessment gateway (R256-R259)
# - Three-agent integration workflow (R269-R270)
# - Integration agent tracking
# - Wave and phase review report tracking
# Backs up existing state and preserves all current data

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
echo -e "${BOLD}ORCHESTRATOR STATE UPGRADE TOOL${NC}"
echo "Adding new architecture planning states (R210, R211)"
echo -e "${GREEN}✅ IMPLEMENTATION-PLAN.md files are protected from overwrites${NC}"
echo "════════════════════════════════════════════════════════════════"
echo ""

# Function to backup state file
backup_state_file() {
    local STATE_FILE="$1"
    local BACKUP_FILE="${STATE_FILE}.backup.$(date +%Y%m%d-%H%M%S)"
    
    echo -e "${BLUE}Creating backup...${NC}"
    cp "$STATE_FILE" "$BACKUP_FILE"
    echo -e "${GREEN}✅ Backup created: $BACKUP_FILE${NC}"
    return 0
}

# Function to check if new fields already exist
check_existing_fields() {
    local STATE_FILE="$1"
    
    if grep -q "architecture_plans:" "$STATE_FILE" 2>/dev/null; then
        echo -e "${YELLOW}⚠️  Architecture plans section already exists${NC}"
        return 1
    fi
    
    if grep -q "phase_architecture_plans:" "$STATE_FILE" 2>/dev/null; then
        echo -e "${YELLOW}⚠️  Phase architecture plans already exists${NC}"
        return 1
    fi
    
    if grep -q "wave_architecture_plans:" "$STATE_FILE" 2>/dev/null; then
        echo -e "${YELLOW}⚠️  Wave architecture plans already exists${NC}"
        return 1
    fi
    
    return 0
}

# Function to add new architecture planning sections
add_architecture_sections() {
    local STATE_FILE="$1"
    local TEMP_FILE="${STATE_FILE}.tmp"
    
    echo -e "${BLUE}Adding new architecture planning sections...${NC}"
    
    # First, convert any tabs to spaces in the original file
    sed -i 's/\t/  /g' "$STATE_FILE"
    
    # Read the entire file
    local CONTENT=$(cat "$STATE_FILE")
    
    # Create the new sections to add
    local NEW_SECTIONS=$(cat <<'EOF'

# ════════════════════════════════════════════════════════════════
# NEW: Architecture Planning States (R210, R211)
# ════════════════════════════════════════════════════════════════

# Phase Architecture Plans (Created by Architect after PHASE_ASSESSMENT passes)
phase_architecture_plans:
  # Example:
  # - phase: 2
  #   status: "created"
  #   file: "phase-plans/PHASE-2-ARCHITECTURE-PLAN.md"
  #   created_by: "architect"
  #   created_at: "2025-08-25T10:00:00Z"
  #   architectural_decisions:
  #     - "Use event-driven architecture for notifications"
  #     - "Implement CQRS for read/write separation"
  #   parallelization_strategy: "Waves 1-2 concurrent, Wave 3 depends on both"

# Wave Architecture Plans (Created by Architect after wave review)
wave_architecture_plans:
  # Example:
  # - phase: 1
  #   wave: 3
  #   status: "created"
  #   file: "phase-plans/PHASE-1-WAVE-3-ARCHITECTURE-PLAN.md"
  #   created_by: "architect"
  #   created_at: "2025-08-25T11:00:00Z"
  #   effort_contracts:
  #     - effort: "api-gateway"
  #       interfaces: ["RouteHandler", "AuthMiddleware"]
  #     - effort: "rate-limiter"
  #       interfaces: ["RateLimiter", "TokenBucket"]
  #   mvp_features: ["Basic routing", "JWT validation"]
  #   nice_to_have: ["Request caching", "Response compression"]

# Phase Implementation Plans (Created by Code Reviewer from Architecture)
phase_implementation_plans:
  # Example:
  # - phase: 2
  #   status: "created"
  #   file: "phase-plans/PHASE-2-IMPLEMENTATION-PLAN.md"
  #   created_by: "code-reviewer"
  #   created_at: "2025-08-25T12:00:00Z"
  #   source_architecture: "phase-plans/PHASE-2-ARCHITECTURE-PLAN.md"
  #   wave_count: 4
  #   total_estimated_lines: 3200

# Wave Implementation Plans (Created by Code Reviewer from Architecture)
wave_implementation_plans:
  # Example:
  # - phase: 1
  #   wave: 3
  #   status: "created"
  #   file: "phase-plans/PHASE-1-WAVE-3-IMPLEMENTATION-PLAN.md"
  #   created_by: "code-reviewer"
  #   created_at: "2025-08-25T13:00:00Z"
  #   source_architecture: "phase-plans/PHASE-1-WAVE-3-ARCHITECTURE-PLAN.md"
  #   effort_count: 3
  #   r213_metadata_injected: true
  #   wave_root: "efforts/phase1/wave3"

# R212 Phase Directory Acknowledgments
phase_directory_acknowledgments:
  # Example:
  # - phase: 1
  #   acknowledged_by: ["architect", "code-reviewer"]
  #   phase_root: "efforts/phase1"
  #   acknowledgment_time: "2025-08-25T09:00:00Z"

# R213 Wave Metadata Injections (Orchestrator as master)
wave_metadata_injections:
  # Example:
  # - phase: 1
  #   wave: 2
  #   metadata_source: "ORCHESTRATOR"
  #   wave_root: "efforts/phase1/wave2"
  #   effort_count: 3
  #   injection_time: "2025-08-25T08:00:00Z"
  #   acknowledged_by: ["code-reviewer"]

# R214 Code Reviewer Wave Acknowledgments
code_reviewer_wave_acknowledgments:
  # Example:
  # - phase: 1
  #   wave: 2
  #   acknowledged: true
  #   wave_root: "efforts/phase1/wave2"
  #   acknowledgment_time: "2025-08-25T14:00:00Z"
  #   effort_plans_created: 3

# New Architect States (for state machine)
architect_review_states:
  phase_assessments:
    # Example:
    # - phase: 1
    #   assessment_time: "2025-08-25T15:00:00Z"
    #   result: "PASS"
    #   architecture_plan_created: true
  wave_reviews:
    # Example:
    # - phase: 1
    #   wave: 1
    #   review_time: "2025-08-25T16:00:00Z"
    #   result: "PROCEED"
    #   architecture_plan_created: true

# New Code Reviewer Planning States
code_reviewer_planning_states:
  # Track when code reviewers are in planning mode vs review mode
  current_mode: "PLANNING"  # PLANNING | REVIEWING
  planning_queue:
    # Example:
    # - phase: 1
    #   wave: 3
    #   type: "wave_implementation"
    #   source_architecture: "phase-plans/PHASE-1-WAVE-3-ARCHITECTURE-PLAN.md"
    #   status: "pending"

# Extended State Machine States
extended_states:
  # New states added for architecture-driven planning
  orchestrator_states:
    - "SPAWN_ARCHITECT_PHASE_PLANNING"
    - "SPAWN_ARCHITECT_WAVE_PLANNING"
    - "SPAWN_CODE_REVIEWER_PHASE_IMPL"
    - "SPAWN_CODE_REVIEWER_WAVE_IMPL"
    - "WAITING_FOR_ARCHITECTURE_PLAN"
    - "WAITING_FOR_IMPLEMENTATION_PLAN"
    # Parallelization analysis states
    - "ANALYZE_CODE_REVIEWER_PARALLELIZATION"
    - "ANALYZE_IMPLEMENTATION_PARALLELIZATION"
    # Phase assessment states
    - "SPAWN_ARCHITECT_PHASE_ASSESSMENT"
    - "WAITING_FOR_PHASE_ASSESSMENT"
    - "PHASE_COMPLETE"
    - "PHASE_INTEGRATION"
    # Integration workflow states
    - "SPAWN_CODE_REVIEWER_MERGE_PLANNING"
    - "WAITING_FOR_MERGE_PLAN"
    - "SPAWN_INTEGRATION_AGENT"
    - "MONITORING_INTEGRATION"
  architect_states:
    - "PHASE_ARCHITECTURE_PLANNING"
    - "WAVE_ARCHITECTURE_PLANNING"
    - "ARCHITECTURE_VALIDATION"
    - "PHASE_ASSESSMENT"
    - "WAVE_REVIEW"
  code_reviewer_states:
    - "PHASE_IMPLEMENTATION_PLANNING"
    - "WAVE_IMPLEMENTATION_PLANNING"
    - "WAVE_DIRECTORY_ACKNOWLEDGMENT"
    - "WAVE_MERGE_PLANNING"
    - "PHASE_MERGE_PLANNING"
  integration_agent_states:
    - "INIT"
    - "MERGING"
    - "VALIDATION"
    - "COMPLETED"

# Parallelization Analysis Tracking
parallelization_analysis:
  code_reviewer_parallelization:
    # Example:
    # - phase: 3
    #   wave: 1
    #   analysis_complete: true
    #   blocking_efforts: ["E3.1.1"]
    #   parallel_efforts: ["E3.1.2", "E3.1.3", "E3.1.4", "E3.1.5"]
  implementation_parallelization:
    # Example:
    # - phase: 3
    #   wave: 1
    #   effort: "E3.1.1"
    #   can_parallelize: false

# Integration Workflow Tracking
integration_workflow:
  wave_integrations:
    # Example:
    # - phase: 3
    #   wave: 1
    #   merge_plan_created: true
    #   merge_plan_file: "integrations/phase3/wave1/WAVE-MERGE-PLAN.md"
    #   integration_agent_spawned: true
    #   integration_complete: true
    #   integration_branch: "phase3-wave1-integration-20250827-143000"
  phase_integrations:
    # Example:
    # - phase: 3
    #   merge_plan_created: true
    #   merge_plan_file: "integrations/phase3/PHASE-MERGE-PLAN.md"
    #   includes_waves: [1, 2, 3]
    #   integration_branch: "phase3-integration-20250827-153000"

# Phase Assessment Tracking
phase_assessments:
  # Example:
  # - phase: 3
  #   assessment_requested: "2025-08-27T15:00:00Z"
  #   assessment_report: "phase-assessments/phase3/PHASE-3-ASSESSMENT-REPORT.md"
  #   decision: "PHASE_COMPLETE"
  #   architect_signoff: "2025-08-27T15:30:00Z"
  #   issues_identified: 0

# Wave Review Tracking  
wave_reviews:
  # Example:
  # - phase: 3
  #   wave: 1
  #   review_requested: "2025-08-27T14:00:00Z"
  #   review_report: "wave-reviews/phase3/wave1/PHASE-3-WAVE-1-REVIEW-REPORT.md"
  #   decision: "PROCEED_NEXT_WAVE"
  #   architect_signoff: "2025-08-27T14:30:00Z"

EOF
)
    
    # Append the new sections to the file
    echo "$NEW_SECTIONS" >> "$STATE_FILE"
    
    # Convert any tabs to spaces in the entire file
    sed -i 's/\t/  /g' "$STATE_FILE"
    
    echo -e "${GREEN}✅ New architecture planning sections added${NC}"
    return 0
}

# Function to validate the upgraded file
validate_upgraded_file() {
    local STATE_FILE="$1"
    
    echo -e "${BLUE}Validating upgraded state file...${NC}"
    
    # Check YAML syntax (basic check)
    echo -e "${BLUE}Checking YAML structure...${NC}"
    
    # Basic YAML validation - check for common issues
    if grep -P '\t' "$STATE_FILE" 2>/dev/null; then
        echo -e "${RED}❌ YAML contains tabs (should use spaces)${NC}"
        return 1
    fi
    
    # Check for unbalanced quotes
    local SINGLE_QUOTES=$(grep -o "'" "$STATE_FILE" | wc -l)
    local DOUBLE_QUOTES=$(grep -o '"' "$STATE_FILE" | wc -l)
    
    if [ $((SINGLE_QUOTES % 2)) -ne 0 ]; then
        echo -e "${YELLOW}⚠️  Possible unbalanced single quotes${NC}"
    fi
    
    if [ $((DOUBLE_QUOTES % 2)) -ne 0 ]; then
        echo -e "${YELLOW}⚠️  Possible unbalanced double quotes${NC}"
    fi
    
    echo -e "${GREEN}✅ Basic YAML structure looks valid${NC}"
    
    # Check for required new sections
    local REQUIRED_SECTIONS=(
        "phase_architecture_plans:"
        "wave_architecture_plans:"
        "phase_implementation_plans:"
        "wave_implementation_plans:"
        "phase_directory_acknowledgments:"
        "wave_metadata_injections:"
        "code_reviewer_wave_acknowledgments:"
        "architect_review_states:"
        "code_reviewer_planning_states:"
        "extended_states:"
        "parallelization_analysis:"
        "integration_workflow:"
        "phase_assessments:"
        "wave_reviews:"
    )
    
    for section in "${REQUIRED_SECTIONS[@]}"; do
        if grep -q "$section" "$STATE_FILE"; then
            echo -e "${GREEN}✓ Section found: $section${NC}"
        else
            echo -e "${RED}✗ Missing section: $section${NC}"
            return 1
        fi
    done
    
    return 0
}

# Main upgrade process
main() {
    local STATE_FILE="$1"
    
    # If no file specified, try common locations
    if [ -z "$STATE_FILE" ]; then
        if [ -f "orchestrator-state.yaml" ]; then
            STATE_FILE="orchestrator-state.yaml"
        elif [ -f "orchestrator-state.yaml.example" ]; then
            STATE_FILE="orchestrator-state.yaml.example"
        else
            echo -e "${RED}❌ No state file specified and none found in current directory${NC}"
            echo "Usage: $0 <orchestrator-state.yaml>"
            exit 1
        fi
    fi
    
    # Check if file exists
    if [ ! -f "$STATE_FILE" ]; then
        echo -e "${RED}❌ State file not found: $STATE_FILE${NC}"
        exit 1
    fi
    
    echo "Processing: $STATE_FILE"
    echo ""
    
    # Check if already upgraded
    if ! check_existing_fields "$STATE_FILE"; then
        echo -e "${YELLOW}State file appears to be already upgraded${NC}"
        echo "Do you want to continue anyway? (y/N)"
        read -r response
        if [[ ! "$response" =~ ^[Yy]$ ]]; then
            echo "Upgrade cancelled"
            exit 0
        fi
    fi
    
    # Backup the file
    backup_state_file "$STATE_FILE"
    
    # Add new sections
    add_architecture_sections "$STATE_FILE"
    
    # Validate the result
    if validate_upgraded_file "$STATE_FILE"; then
        echo ""
        echo "════════════════════════════════════════════════════════════════"
        echo -e "${GREEN}${BOLD}🎉 UPGRADE SUCCESSFUL!${NC}"
        echo "════════════════════════════════════════════════════════════════"
        echo ""
        echo "The orchestrator state file has been upgraded with:"
        echo "✅ Phase architecture planning support (R210)"
        echo "✅ Wave architecture planning support (R210)"
        echo "✅ Implementation from architecture support (R211)"
        echo "✅ Phase directory acknowledgments (R212)"
        echo "✅ Wave metadata injections (R213)"
        echo "✅ Code reviewer wave acknowledgments (R214)"
        echo "✅ Parallelization analysis states (R231-R233)"
        echo "✅ Phase assessment gateway (R256-R257)"
        echo "✅ Wave review reports (R258)"
        echo "✅ Phase integration after fixes (R259)"
        echo "✅ Three-agent integration workflow (R269-R270)"
        echo "✅ Integration agent tracking"
        echo "✅ Extended state machine states"
        echo ""
        echo "New States Added:"
        echo "• ANALYZE_CODE_REVIEWER_PARALLELIZATION"
        echo "• ANALYZE_IMPLEMENTATION_PARALLELIZATION"
        echo "• SPAWN_ARCHITECT_PHASE_ASSESSMENT"
        echo "• WAITING_FOR_PHASE_ASSESSMENT"
        echo "• PHASE_COMPLETE"
        echo "• PHASE_INTEGRATION"
        echo "• SPAWN_CODE_REVIEWER_MERGE_PLANNING"
        echo "• WAITING_FOR_MERGE_PLAN"
        echo "• SPAWN_INTEGRATION_AGENT"
        echo "• MONITORING_INTEGRATION"
        echo "• WAVE_MERGE_PLANNING (Code Reviewer)"
        echo "• PHASE_MERGE_PLANNING (Code Reviewer)"
        echo ""
        echo "Next steps:"
        echo "1. Review the upgraded file: $STATE_FILE"
        echo "2. Update orchestrator to use new states"
        echo "3. Test parallelization analysis workflow"
        echo "4. Test phase assessment gateway"
        echo "5. Test three-agent integration workflow"
    else
        echo ""
        echo -e "${RED}❌ Validation failed!${NC}"
        echo "Please check the file for errors"
        exit 1
    fi
}

# Run main function with all arguments
main "$@"