#!/bin/bash

# Software Factory 1.0 to 2.0 Planning Migration Script
# Keeps architectural planning, discards implementations

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
MAGENTA='\033[0;35m'
CYAN='\033[0;36m'
BOLD='\033[1m'
NC='\033[0m' # No Color

echo -e "${CYAN}${BOLD}"
cat << "EOF"
╔═══════════════════════════════════════════════════════════════════╗
║                                                                   ║
║   Planning Migration - SF 1.0 → 2.0                              ║
║   Keep Architecture, Regenerate Implementation                   ║
║                                                                   ║
╚═══════════════════════════════════════════════════════════════════╝
EOF
echo -e "${NC}"

# Function to extract planning from SF 1.0
extract_planning() {
    local source="$1"
    local target="$2"
    
    echo -e "${CYAN}Extracting planning artifacts from SF 1.0...${NC}"
    
    # Create planning directory
    mkdir -p "$target/planning/original-sf1"
    
    # Copy phase plans if they exist
    if [ -d "$source/phase-plans" ]; then
        cp -r "$source/phase-plans" "$target/planning/original-sf1/"
        echo -e "${GREEN}✓ Phase plans copied${NC}"
    fi
    
    # Copy high-level plans
    for file in PROJECT-IMPLEMENTATION-PLAN*.md PLANNING*.md HOW-TO-PLAN.md; do
        if [ -f "$source/$file" ]; then
            cp "$source/$file" "$target/planning/original-sf1/"
            echo -e "${GREEN}✓ $file copied${NC}"
        fi
    done
    
    # Extract orchestrator state for phase/wave structure
    if [ -f "$source/orchestrator-state-v3.json" ]; then
        cp "$source/orchestrator-state-v3.json" "$target/planning/original-sf1/state-snapshot.yaml"
        echo -e "${GREEN}✓ State structure extracted${NC}"
    fi
}

# Function to analyze phase structure
analyze_phase_structure() {
    local state_file="$1"
    local output_file="$2"
    
    echo -e "${CYAN}Analyzing phase structure...${NC}"
    
    cat > "$output_file" << 'EOF'
# Extracted Phase Structure from SF 1.0

## Project Overview
EOF
    
    # Extract project info
    if [ -f "$state_file" ]; then
        echo '```yaml' >> "$output_file"
        grep -A 10 "^project:" "$state_file" >> "$output_file" 2>/dev/null || true
        echo '```' >> "$output_file"
        
        echo -e "\n## Phase Structure\n" >> "$output_file"
        echo '```yaml' >> "$output_file"
        grep -A 100 "^phases:" "$state_file" >> "$output_file" 2>/dev/null || true
        echo '```' >> "$output_file"
        
        echo -e "\n## Wave Structure\n" >> "$output_file"
        echo '```yaml' >> "$output_file"
        grep -A 200 "^waves:" "$state_file" >> "$output_file" 2>/dev/null || true
        echo '```' >> "$output_file"
    fi
    
    echo -e "${GREEN}✓ Phase structure analyzed${NC}"
}

# Function to create master plan
create_master_plan() {
    local planning_dir="$1"
    local target="$2"
    
    echo -e "${CYAN}Creating master implementation plan...${NC}"
    
    cat > "$target/MASTER-IMPLEMENTATION-PLAN.md" << 'EOF'
# Master Implementation Plan
# Migrated from Software Factory 1.0

## Migration Notes
- **Date:** $(date)
- **Approach:** Keep planning, regenerate implementation
- **Reason:** Want SF 2.0 quality and automation

## Project Structure
[Imported from SF 1.0 - See planning/original-sf1/]

## Phase Breakdown

### Phase 1: [Name]
**Goal:** [From original plan]
**Waves:** [X]
**Efforts:** [Y]
**Test Coverage:** [Z]%

#### Wave 1.1: [Name]
- Effort 1.1.1: [Name] (est. [X] lines)
- Effort 1.1.2: [Name] (est. [X] lines)

#### Wave 1.2: [Name]
[Continue...]

## Dependency Graph
[Import from original planning]

## Success Criteria
[Import from original planning]

## Notes for SF 2.0 Regeneration

### What to Keep:
1. Phase/Wave/Effort decomposition
2. Dependency relationships
3. Test coverage targets
4. Success criteria

### What to Regenerate:
1. Effort-level IMPLEMENTATION-PLAN.md
2. Work logs
3. Review criteria
4. Test scaffolding

## Instructions for Orchestrator

When starting:
1. Load this master plan
2. For each effort, spawn Code Reviewer to create:
   - Detailed IMPLEMENTATION-PLAN.md
   - work-log.md template
   - Test scaffolding
3. Use SF 2.0 grading and enforcement
EOF
    
    echo -e "${GREEN}✓ Master plan template created${NC}"
    echo -e "${YELLOW}NOTE: Edit MASTER-IMPLEMENTATION-PLAN.md to add specific details${NC}"
}

# Function to create initial orchestrator state
create_initial_state() {
    local target="$1"
    local old_state="$2"
    
    echo -e "${CYAN}Creating initial SF 2.0 orchestrator state...${NC}"
    
    # Extract phase count from old state
    local num_phases=$(grep -c "^  - phase:" "$old_state" 2>/dev/null || echo "3")
    
    cat > "$target/orchestrator-state-v3.json" << EOF
# Orchestrator State - SF 2.0
# Migrated Planning from SF 1.0
# Implementation will be regenerated

current_phase: 1
current_wave: 1
current_state: INIT

# SF 2.0 Configuration
state_machine_version: "2.0"
migration_info:
  migrated_from: "SF 1.0"
  migration_date: "$(date -Iseconds)"
  migration_type: "planning-only"
  implementation_status: "to-be-regenerated"

# Grading metrics (fresh start)
grading_history:
  parallel_spawn_average: 0.0
  parallel_spawn_count: 0
  review_first_try_rate: 0.0
  review_total_count: 0
  integration_success_rate: 0.0
  integration_total_count: 0
  size_compliance_rate: 0.0
  size_check_count: 0

# Planning imported from SF 1.0
planning_source: "planning/original-sf1/"

# Phases structure (to be filled from master plan)
phases_planned: $num_phases

# All efforts will be regenerated
efforts_completed: []
efforts_in_progress: []
efforts_pending: []  # To be populated from master plan

# Integration branches (fresh)
integration_branches: []

# Checkpoint
last_checkpoint: "$(date -Iseconds)"
checkpoint_version: "2.0"

# Notes
notes:
  - timestamp: "$(date -Iseconds)"
    note: "Planning migrated from SF 1.0, implementations to be regenerated with SF 2.0 quality"
EOF
    
    echo -e "${GREEN}✓ Initial state created${NC}"
}

# Main function
main() {
    echo -e "${MAGENTA}${BOLD}Planning-Only Migration Tool${NC}\n"
    
    # Get source directory
    local source_dir="${1:-}"
    if [ -z "$source_dir" ]; then
        echo -ne "${CYAN}Enter path to SF 1.0 project: ${NC}"
        read -r source_dir
    fi
    
    # Validate source
    if [ ! -d "$source_dir" ]; then
        echo -e "${RED}Error: Directory not found: $source_dir${NC}"
        exit 1
    fi
    
    if [ ! -f "$source_dir/orchestrator-state-v3.json" ]; then
        echo -e "${RED}Error: Not a Software Factory project (no orchestrator-state-v3.json)${NC}"
        exit 1
    fi
    
    # Get target directory
    echo -ne "${CYAN}Enter name for new SF 2.0 project: ${NC}"
    read -r project_name
    
    local target_dir="/workspaces/${project_name}-sf2"
    
    echo -e "\n${YELLOW}${BOLD}Migration Plan:${NC}"
    echo -e "  ${CYAN}Source:${NC} $source_dir"
    echo -e "  ${CYAN}Target:${NC} $target_dir"
    echo -e "  ${CYAN}Type:${NC} Planning only (discard implementations)"
    
    echo -ne "\n${YELLOW}Proceed? (y/n): ${NC}"
    read -r confirm
    if [ "$confirm" != "y" ]; then
        echo -e "${RED}Migration cancelled${NC}"
        exit 0
    fi
    
    # Run setup.sh to create SF 2.0 structure
    echo -e "\n${CYAN}Creating SF 2.0 project structure...${NC}"
    
    # Find the directory where this script is located
    local script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
    
    # Check if we can find setup.sh in the same directory
    local setup_script="$script_dir/setup.sh"
    if [ ! -f "$setup_script" ]; then
        echo -e "${RED}Error: setup.sh not found at $setup_script${NC}"
        echo -e "${YELLOW}Make sure you're running this from the SF 2.0 template directory${NC}"
        exit 1
    fi
    
    # Create project using setup (automated responses for common case)
    echo -e "${CYAN}Running SF 2.0 setup...${NC}"
    echo -e "${YELLOW}Please answer the setup wizard questions${NC}"
    
    # Let user run setup interactively from the script's directory
    cd "$script_dir"
    ./setup.sh
    
    # Now work with the created project
    if [ ! -d "$target_dir" ]; then
        # Try to find where it was created
        echo -ne "${CYAN}Enter the path where SF 2.0 project was created: ${NC}"
        read -r target_dir
    fi
    
    if [ ! -d "$target_dir" ]; then
        echo -e "${RED}Error: Target directory not found${NC}"
        exit 1
    fi
    
    # Extract planning from SF 1.0
    extract_planning "$source_dir" "$target_dir"
    
    # Ensure SF 2.0 agent configs and settings are properly configured
    echo -e "\n${CYAN}Verifying SF 2.0 agent configurations...${NC}"
    
    # Check if .claude/agents exists and has content
    if [ -d "$target_dir/.claude/agents" ] && [ "$(ls -A $target_dir/.claude/agents)" ]; then
        echo -e "${GREEN}✓ Agent configurations present${NC}"
    else
        echo -e "${YELLOW}⚠ Agent configurations missing - copying from template${NC}"
        local script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
        if [ -d "$script_dir/.claude/agents" ]; then
            cp -r "$script_dir/.claude/agents" "$target_dir/.claude/"
            echo -e "${GREEN}✓ Agent configurations restored${NC}"
        fi
    fi
    
    # Check settings.json
    if [ -f "$target_dir/.claude/settings.json" ]; then
        # Update paths to point to new project
        sed -i "s|/workspaces/software-factory-2.0-template|$target_dir|g" "$target_dir/.claude/settings.json"
        echo -e "${GREEN}✓ Settings.json configured${NC}"
    else
        echo -e "${YELLOW}⚠ Settings.json missing - copying from template${NC}"
        local script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
        if [ -f "$script_dir/.claude/settings.json" ]; then
            cp "$script_dir/.claude/settings.json" "$target_dir/.claude/"
            sed -i "s|/workspaces/software-factory-2.0-template|$target_dir|g" "$target_dir/.claude/settings.json"
            echo -e "${GREEN}✓ Settings.json created${NC}"
        fi
    fi
    
    # Analyze phase structure
    analyze_phase_structure \
        "$target_dir/planning/original-sf1/state-snapshot.yaml" \
        "$target_dir/planning/EXTRACTED-STRUCTURE.md"
    
    # Create master plan
    create_master_plan "$target_dir/planning" "$target_dir"
    
    # Create initial state
    create_initial_state "$target_dir" "$source_dir/orchestrator-state-v3.json"
    
    # Create migration report
    cat > "$target_dir/MIGRATION-REPORT.md" << EOF
# Planning Migration Report

## Migration Summary
- **Date:** $(date)
- **Type:** Planning-only migration
- **Source:** $source_dir
- **Target:** $target_dir

## What Was Migrated

### ✅ Kept from SF 1.0:
- Phase/Wave/Effort structure
- Dependency relationships
- Test coverage requirements
- Success criteria
- Project goals

### 🔄 To Be Regenerated in SF 2.0:
- All effort implementations
- Effort-level IMPLEMENTATION-PLAN.md files
- Work logs
- Review criteria
- Test scaffolding

## Next Steps

1. **Review and edit the master plan:**
   \`\`\`bash
   cat $target_dir/MASTER-IMPLEMENTATION-PLAN.md
   # Edit to add specific effort details
   \`\`\`

2. **Start the orchestrator:**
   \`\`\`
   /continue-orchestrating
   \`\`\`

3. **Let SF 2.0 generate better plans:**
   - Orchestrator reads MASTER-IMPLEMENTATION-PLAN.md
   - Spawns Code Reviewer for each effort
   - Code Reviewer creates detailed plans
   - Plans include grading, tests, checkpoints

## File Locations

- **Original SF 1.0 Plans:** planning/original-sf1/
- **Extracted Structure:** planning/EXTRACTED-STRUCTURE.md
- **Master Plan:** MASTER-IMPLEMENTATION-PLAN.md
- **Orchestrator State:** orchestrator-state-v3.json

## Benefits of This Approach

1. **Preserve Good Architecture** - Your planning was solid
2. **Better Implementation Quality** - SF 2.0 enforcement
3. **Automatic Size Management** - No manual checking
4. **Superior Test Coverage** - TDD templates
5. **Faster Development** - Parallel spawning

## Important Notes

- No SF 1.0 code was copied (fresh start)
- All implementations will be regenerated
- Grading starts from zero (clean baseline)
- State machine begins at INIT

---

**Migration Complete!** Review MASTER-IMPLEMENTATION-PLAN.md and start orchestrating.
EOF
    
    # Git commit
    cd "$target_dir"
    git add -A
    git commit -m "Migrate planning from SF 1.0 to SF 2.0

- Imported phase/wave/effort structure
- Preserved architectural planning
- Ready for SF 2.0 implementation generation
- Type: Planning-only migration" 2>/dev/null || true
    
    # Success message
    echo -e "\n${GREEN}${BOLD}════════════════════════════════════════════════${NC}"
    echo -e "${GREEN}${BOLD}✅ Planning Migration Complete!${NC}"
    echo -e "${GREEN}${BOLD}════════════════════════════════════════════════${NC}\n"
    
    echo -e "${CYAN}Your SF 2.0 project is ready at:${NC}"
    echo -e "  ${BOLD}$target_dir${NC}\n"
    
    echo -e "${YELLOW}Next steps:${NC}"
    echo -e "  1. ${CYAN}cd $target_dir${NC}"
    echo -e "  2. ${CYAN}Review and edit MASTER-IMPLEMENTATION-PLAN.md${NC}"
    echo -e "  3. ${CYAN}Add specific effort details from SF 1.0 plans${NC}"
    echo -e "  4. Run ${CYAN}/continue-orchestrating${NC} in Claude"
    
    echo -e "\n${GREEN}Your architecture is preserved, implementations will be better! 🚀${NC}"
}

# Run main
main "$@"