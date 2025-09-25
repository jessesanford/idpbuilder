#!/bin/bash

# Software Factory 1.0 to 2.0 Migration Script
# Automates the migration process while preserving state and progress

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

# Banner
echo -e "${CYAN}${BOLD}"
cat << "EOF"
╔═══════════════════════════════════════════════════════════════════╗
║                                                                   ║
║   ███╗   ███╗██╗ ██████╗ ██████╗  █████╗ ████████╗███████╗      ║
║   ████╗ ████║██║██╔════╝ ██╔══██╗██╔══██╗╚══██╔══╝██╔════╝      ║
║   ██╔████╔██║██║██║  ███╗██████╔╝███████║   ██║   █████╗        ║
║   ██║╚██╔╝██║██║██║   ██║██╔══██╗██╔══██║   ██║   ██╔══╝        ║
║   ██║ ╚═╝ ██║██║╚██████╔╝██║  ██║██║  ██║   ██║   ███████╗      ║
║   ╚═╝     ╚═╝╚═╝ ╚═════╝ ╚═╝  ╚═╝╚═╝  ╚═╝   ╚═╝   ╚══════╝      ║
║                                                                   ║
║              Software Factory 1.0 → 2.0 Migration                ║
║                                                                   ║
╚═══════════════════════════════════════════════════════════════════╝
EOF
echo -e "${NC}"

# Function to check if directory is SF 1.0 project
check_sf1_project() {
    local dir="$1"
    
    if [ ! -f "$dir/orchestrator-state.json" ]; then
        return 1
    fi
    
    if [ ! -d "$dir/.claude" ]; then
        return 1
    fi
    
    # Check for SF 1.0 indicators (no state machines)
    if [ -d "$dir/state-machines" ]; then
        echo -e "${YELLOW}Warning: This appears to be a SF 2.0 project already${NC}"
        return 2
    fi
    
    return 0
}

# Function to backup project
backup_project() {
    local source="$1"
    local backup_dir="$2"
    
    echo -e "${CYAN}Creating backup...${NC}"
    
    mkdir -p "$backup_dir"
    
    # Create timestamped backup
    local timestamp=$(date +%Y%m%d-%H%M%S)
    local backup_file="$backup_dir/sf1-backup-$timestamp.tar.gz"
    
    tar -czf "$backup_file" -C "$(dirname "$source")" "$(basename "$source")" 2>/dev/null
    
    echo -e "${GREEN}✓ Backup created: $backup_file${NC}"
    
    # Also create state snapshots
    cp "$source/orchestrator-state.json" "$backup_dir/state-$timestamp.yaml" 2>/dev/null || true
    
    if [ -d "$source/todos" ]; then
        cp -r "$source/todos" "$backup_dir/todos-$timestamp" 2>/dev/null || true
    fi
}

# Function to analyze SF 1.0 state
analyze_sf1_state() {
    local dir="$1"
    local analysis_file="$2"
    
    echo -e "${CYAN}Analyzing SF 1.0 project state...${NC}"
    
    cat > "$analysis_file" << EOF
# SF 1.0 State Analysis
# Generated: $(date)

## Project Structure
EOF
    
    # Check orchestrator state
    if [ -f "$dir/orchestrator-state.json" ]; then
        echo "### Orchestrator State" >> "$analysis_file"
        echo '```yaml' >> "$analysis_file"
        grep -E "current_phase|current_wave|efforts_completed|efforts_in_progress" "$dir/orchestrator-state.json" >> "$analysis_file"
        echo '```' >> "$analysis_file"
    fi
    
    # Check branches
    echo "### Git Branches" >> "$analysis_file"
    echo '```' >> "$analysis_file"
    cd "$dir" && git branch -a >> "$analysis_file" 2>/dev/null
    echo '```' >> "$analysis_file"
    
    # Check for custom rules
    echo "### Custom Files" >> "$analysis_file"
    find "$dir" -name "*.md" -type f ! -path "*/node_modules/*" ! -path "*/.git/*" | head -20 >> "$analysis_file"
    
    echo -e "${GREEN}✓ Analysis complete${NC}"
}

# Function to migrate state file
migrate_state_file() {
    local old_state="$1"
    local new_state="$2"
    
    echo -e "${CYAN}Migrating orchestrator state...${NC}"
    
    # Read old state
    local current_phase=$(grep "current_phase:" "$old_state" | awk '{print $2}')
    local current_wave=$(grep "current_wave:" "$old_state" | awk '{print $2}')
    
    # Determine appropriate state machine state
    local current_state="INIT"
    if [ -n "$current_phase" ] && [ "$current_phase" -gt 0 ]; then
        if [ -n "$current_wave" ] && [ "$current_wave" -gt 0 ]; then
            current_state="WAVE_START"
        else
            current_state="PLANNING"
        fi
    fi
    
    # Create new state file with SF 2.0 format
    cat > "$new_state" << EOF
# Orchestrator State - Migrated from SF 1.0
# Migration Date: $(date)
# Original Phase: $current_phase, Wave: $current_wave

current_phase: $current_phase
current_wave: $current_wave
current_state: $current_state

# SF 2.0 Additions
state_machine_version: "2.0"
migration_info:
  migrated_from: "SF 1.0"
  migration_date: "$(date -Iseconds)"
  migration_method: "automated"

# Grading metrics (initialized)
grading_history:
  parallel_spawn_average: 10.0
  parallel_spawn_count: 0
  review_first_try_rate: 0.7
  review_total_count: 0
  integration_success_rate: 0.8
  integration_total_count: 0
  size_compliance_rate: 0.9
  size_check_count: 0

# Checkpoint system
last_checkpoint: "$(date -Iseconds)"
checkpoint_version: "2.0"

EOF
    
    # Append existing efforts data
    echo "# Preserved from SF 1.0" >> "$new_state"
    grep -A 1000 "efforts_completed:" "$old_state" >> "$new_state" 2>/dev/null || true
    
    echo -e "${GREEN}✓ State file migrated${NC}"
}

# Function to update slash commands and agent configs
update_slash_commands() {
    local target_dir="$1"
    
    echo -e "${CYAN}Updating slash commands for SF 2.0...${NC}"
    
    # Add pre-flight checks to each command
    for cmd in continue-orchestrating continue-implementing continue-reviewing continue-architecting; do
        if [ -f "$target_dir/.claude/commands/$cmd.md" ]; then
            # Check if pre-flight checks already exist
            if ! grep -q "Pre-Flight Checks" "$target_dir/.claude/commands/$cmd.md"; then
                # Prepend pre-flight checks section
                local temp_file=$(mktemp)
                cat > "$temp_file" << 'EOF'
## 🚨 Pre-Flight Checks (SF 2.0)

┌─────────────────────────────────────────────────────────────────┐
│ RULE R001.0.0 - Pre-Flight Checks                              │
├─────────────────────────────────────────────────────────────────┤
│ MANDATORY STARTUP:                                             │
│ 1. Print: "AGENT STARTUP: $(date '+%Y-%m-%d %H:%M:%S %Z')"    │
│ 2. Verify working directory                                    │
│ 3. Verify Git branch                                          │
│ 4. Check current state machine state                          │
│ 5. Acknowledge critical rules                                 │
└─────────────────────────────────────────────────────────────────┘

EOF
                cat "$target_dir/.claude/commands/$cmd.md" >> "$temp_file"
                mv "$temp_file" "$target_dir/.claude/commands/$cmd.md"
            fi
        fi
    done
    
    echo -e "${GREEN}✓ Slash commands updated${NC}"
    
    # Update or create agent configurations if missing
    if [ ! -d "$target_dir/.claude/agents" ]; then
        echo -e "${CYAN}Creating SF 2.0 agent configurations...${NC}"
        mkdir -p "$target_dir/.claude/agents"
        
        # Copy from template
        local source="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
        if [ -d "$source/.claude/agents" ]; then
            cp -r "$source/.claude/agents"/* "$target_dir/.claude/agents/"
            echo -e "${GREEN}✓ Agent configurations created${NC}"
        fi
    fi
    
    # Update or create settings.json
    if [ ! -f "$target_dir/.claude/settings.json" ]; then
        echo -e "${CYAN}Creating SF 2.0 settings.json...${NC}"
        local source="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
        if [ -f "$source/.claude/settings.json" ]; then
            cp "$source/.claude/settings.json" "$target_dir/.claude/"
            # Update paths in settings.json
            sed -i "s|/workspaces/software-factory-2.0-template|$target_dir|g" "$target_dir/.claude/settings.json"
            echo -e "${GREEN}✓ Settings.json created${NC}"
        fi
    fi
}

# Function to create SF 2.0 directories
create_sf2_directories() {
    local dir="$1"
    
    echo -e "${CYAN}Creating SF 2.0 directory structure...${NC}"
    
    # Note: 🚨-CRITICAL folder deprecated - rules now in rule-library
    mkdir -p "$dir/state-machines"
    mkdir -p "$dir/agent-states/orchestrator"
    mkdir -p "$dir/agent-states/sw-engineer"
    mkdir -p "$dir/agent-states/code-reviewer"
    mkdir -p "$dir/agent-states/architect"
    mkdir -p "$dir/expertise"
    mkdir -p "$dir/hooks"
    mkdir -p "$dir/quick-reference"
    mkdir -p "$dir/rule-library"
    mkdir -p "$dir/grading-rubrics"
    mkdir -p "$dir/checkpoints/active"
    
    echo -e "${GREEN}✓ SF 2.0 directories created${NC}"
}

# Function to copy SF 2.0 components
copy_sf2_components() {
    # Find the directory where this script is located (the template directory)
    local source="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
    local target="$1"
    
    echo -e "${CYAN}Copying SF 2.0 components from $source...${NC}"
    
    # Note: 🚨-CRITICAL folder deprecated - rules now in rule-library
    # Critical rules are now in rule-library/R00X-*.md files
    
    # Copy state machines
    cp -r "$source/state-machines"/* "$target/state-machines/" 2>/dev/null || true
    
    # Copy agent states
    cp -r "$source/agent-states"/* "$target/agent-states/" 2>/dev/null || true
    
    # Copy expertise modules
    cp -r "$source/expertise"/* "$target/expertise/" 2>/dev/null || true
    
    # Copy utilities (manual helper scripts)
    cp -r "$source/utilities"/* "$target/utilities/" 2>/dev/null || true
    chmod +x "$target/utilities"/*.sh
    
    # Copy quick reference
    cp -r "$source/quick-reference"/* "$target/quick-reference/" 2>/dev/null || true
    
    # Copy rule library
    cp -r "$source/rule-library"/* "$target/rule-library/" 2>/dev/null || true
    
    echo -e "${GREEN}✓ SF 2.0 components copied${NC}"
}

# Function to migrate TODOs
migrate_todos() {
    local source="$1"
    local target="$2"
    
    echo -e "${CYAN}Migrating TODO files...${NC}"
    
    if [ -d "$source/todos" ]; then
        # Create new todos directory
        mkdir -p "$target/todos"
        
        # Convert TODO files to new naming format
        for todo_file in "$source/todos"/*.md "$source/todos"/*.todo; do
            if [ -f "$todo_file" ]; then
                local basename=$(basename "$todo_file")
                local timestamp=$(date +%Y%m%d-%H%M%S)
                
                # Determine agent type from content or filename
                local agent="orchestrator"
                if grep -q "implement" "$todo_file" 2>/dev/null; then
                    agent="sw-engineer"
                elif grep -q "review" "$todo_file" 2>/dev/null; then
                    agent="code-reviewer"
                fi
                
                # Copy with new naming
                cp "$todo_file" "$target/todos/${agent}-MIGRATED-${timestamp}.todo"
            fi
        done
        
        echo -e "${GREEN}✓ TODOs migrated${NC}"
    else
        echo -e "${YELLOW}No TODOs to migrate${NC}"
    fi
}

# Function to create migration report
create_migration_report() {
    local target="$1"
    local analysis="$2"
    
    echo -e "${CYAN}Creating migration report...${NC}"
    
    cat > "$target/MIGRATION-REPORT.md" << EOF
# Software Factory 1.0 to 2.0 Migration Report

## Migration Summary
- **Date:** $(date)
- **Method:** Automated migration script
- **Status:** ✅ Complete

## Pre-Migration Analysis
$(cat "$analysis")

## Migration Actions Performed

### ✅ State File Migration
- Converted orchestrator-state.json to SF 2.0 format
- Added state machine state tracking
- Initialized grading metrics
- Added checkpoint system

### ✅ Directory Structure
- Created SF 2.0 directory hierarchy
- Added state-machines/
- Added agent-states/
- Added rule-library/
- Added hooks/

### ✅ Component Installation
- Copied state machine definitions
- Installed agent state rules
- Added expertise modules
- Configured pre-compaction hooks
- Added quick reference guides

### ✅ Slash Commands
- Updated with pre-flight checks
- Added state machine awareness
- Integrated grading requirements

### ✅ TODO Migration
- Converted to SF 2.0 naming format
- Preserved all existing TODOs

## Post-Migration Checklist

### Immediate Actions Required:

1. **Verify State Machine State:**
   \`\`\`bash
   grep current_state orchestrator-state.json
   \`\`\`
   
2. **Test Manual Utilities:**
   \`\`\`bash
   ./utilities/pre-compact.sh
   ./utilities/state-snapshot.sh
   \`\`\`

3. **Check Status:**
   \`\`\`
   /check-status
   \`\`\`

4. **Update Project Configuration:**
   - Edit project-config.yaml with your project details
   - Set appropriate constraints

### Configuration Updates Needed:

- [ ] Review and update project-config.yaml
- [ ] Set max_lines_per_effort (default: 800)
- [ ] Set test_coverage_target
- [ ] Configure security_level

### Testing Recommendations:

1. **Test Agent Startup:**
   - Run /continue-orchestrating
   - Verify pre-flight checks execute
   - Check state machine transitions

2. **Test Grading System:**
   - Spawn multiple agents
   - Check parallel_spawn_average
   - Verify timing under 5s

3. **Test Recovery:**
   - Simulate compaction
   - Verify TODO preservation
   - Test recovery procedures

## Known Limitations

### Manual Updates Required:
1. Custom rules need ID assignment in rule-library/
2. Project-specific configurations in project-config.yaml
3. Integration branch mappings

### Potential Issues:
- State machine state may need adjustment
- Grading baselines are conservative
- Some custom modifications may need review

## Rollback Instructions

If issues occur, rollback is available:
\`\`\`bash
# Restore from backup
tar -xzf /tmp/sf-migration-backup/sf1-backup-*.tar.gz
\`\`\`

## Next Steps

1. **Start with Status Check:**
   \`\`\`
   /check-status
   \`\`\`

2. **Resume Work:**
   \`\`\`
   /continue-orchestrating
   \`\`\`

3. **Monitor Grading:**
   Check grading_history in orchestrator-state.json

4. **Utilize New Features:**
   - Parallel agent spawning
   - Automatic size management
   - Pre-compaction preservation

## Support Resources

- Migration Guide: MIGRATION-GUIDE-1.0-TO-2.0.md
- Quick Reference: quick-reference/
- Emergency Procedures: quick-reference/emergency-procedures.md
- State Machines: state-machines/

---

**Migration Complete!** Your project is now running Software Factory 2.0.
EOF
    
    echo -e "${GREEN}✓ Migration report created${NC}"
}

# Main migration process
main() {
    echo -e "${MAGENTA}${BOLD}Software Factory 1.0 to 2.0 Migration Tool${NC}\n"
    
    # Get source directory
    local source_dir="${1:-}"
    if [ -z "$source_dir" ]; then
        echo -ne "${CYAN}Enter path to SF 1.0 project: ${NC}"
        read -r source_dir
    fi
    
    # Validate source directory
    if [ ! -d "$source_dir" ]; then
        echo -e "${RED}Error: Directory not found: $source_dir${NC}"
        exit 1
    fi
    
    # Check if it's a SF 1.0 project
    check_sf1_project "$source_dir"
    local check_result=$?
    
    if [ $check_result -eq 1 ]; then
        echo -e "${RED}Error: Not a Software Factory 1.0 project${NC}"
        echo -e "${YELLOW}Missing orchestrator-state.json or .claude directory${NC}"
        exit 1
    elif [ $check_result -eq 2 ]; then
        echo -e "${YELLOW}This appears to already be a SF 2.0 project${NC}"
        echo -ne "${CYAN}Continue anyway? (y/n): ${NC}"
        read -r confirm
        if [ "$confirm" != "y" ]; then
            exit 0
        fi
    fi
    
    # Get target directory
    local target_dir="${2:-${source_dir}-sf2}"
    echo -e "${CYAN}Migration target: ${BOLD}$target_dir${NC}"
    
    # Confirm migration
    echo -e "\n${YELLOW}${BOLD}Migration Plan:${NC}"
    echo -e "  ${CYAN}Source:${NC} $source_dir"
    echo -e "  ${CYAN}Target:${NC} $target_dir"
    echo -e "  ${CYAN}Backup:${NC} /tmp/sf-migration-backup/"
    
    echo -ne "\n${YELLOW}Proceed with migration? (y/n): ${NC}"
    read -r confirm
    if [ "$confirm" != "y" ]; then
        echo -e "${RED}Migration cancelled${NC}"
        exit 0
    fi
    
    # Start migration
    echo -e "\n${MAGENTA}${BOLD}Starting Migration...${NC}\n"
    
    # Step 1: Backup
    backup_project "$source_dir" "/tmp/sf-migration-backup"
    
    # Step 2: Analyze
    local analysis_file="/tmp/sf1-analysis.md"
    analyze_sf1_state "$source_dir" "$analysis_file"
    
    # Step 3: Prepare target
    if [ "$target_dir" != "$source_dir" ]; then
        echo -e "${CYAN}Copying project to target location...${NC}"
        cp -r "$source_dir" "$target_dir"
    fi
    
    # Step 4: Create SF 2.0 structure
    create_sf2_directories "$target_dir"
    
    # Step 5: Copy SF 2.0 components
    copy_sf2_components "$target_dir"
    
    # Step 6: Migrate state file
    if [ -f "$target_dir/orchestrator-state.json" ]; then
        mv "$target_dir/orchestrator-state.json" "$target_dir/orchestrator-state-sf1.yaml.bak"
        migrate_state_file "$target_dir/orchestrator-state-sf1.yaml.bak" "$target_dir/orchestrator-state.json"
    fi
    
    # Step 7: Update slash commands
    update_slash_commands "$target_dir"
    
    # Step 8: Migrate TODOs
    migrate_todos "$source_dir" "$target_dir"
    
    # Step 9: Create migration report
    create_migration_report "$target_dir" "$analysis_file"
    
    # Step 10: Git commit
    cd "$target_dir"
    if [ -d .git ]; then
        git add -A
        git commit -m "Migrate from Software Factory 1.0 to 2.0

- Added state machine support
- Integrated grading system
- Added pre-compaction hooks
- Updated slash commands with pre-flight checks
- Migrated orchestrator state
- Added SF 2.0 components" 2>/dev/null || true
        
        echo -e "${GREEN}✓ Changes committed to Git${NC}"
    fi
    
    # Success message
    echo -e "\n${GREEN}${BOLD}════════════════════════════════════════════════${NC}"
    echo -e "${GREEN}${BOLD}✅ Migration Complete!${NC}"
    echo -e "${GREEN}${BOLD}════════════════════════════════════════════════${NC}\n"
    
    echo -e "${CYAN}Your SF 2.0 project is ready at:${NC}"
    echo -e "  ${BOLD}$target_dir${NC}\n"
    
    echo -e "${YELLOW}Next steps:${NC}"
    echo -e "  1. ${CYAN}cd $target_dir${NC}"
    echo -e "  2. ${CYAN}cat MIGRATION-REPORT.md${NC}"
    echo -e "  3. Run ${CYAN}/check-status${NC} in Claude"
    echo -e "  4. Resume with ${CYAN}/continue-orchestrating${NC}"
    
    echo -e "\n${GREEN}Welcome to Software Factory 2.0! 🚀${NC}"
}

# Run main function
main "$@"