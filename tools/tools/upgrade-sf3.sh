#!/bin/bash

# Software Factory 3.0 Upgrade Script
# This script updates an existing SF 3.0 instance with the latest rules, state machines,
# and agent configurations while preserving work in progress.
#
# CRITICAL FEATURES:
#   - SF 3.0 specific validation and updates
#   - State Manager agent validation
#   - Integration hierarchy compliance checks
#   - Comprehensive backup with git history preservation
#   - Dry-run mode for safety
#   - Detailed change reporting
#
# Usage:
#   ./upgrade-sf3.sh [--dry-run] [--verbose] [--no-backup]
#
# Recovery:
#   Backups are created at: .upgrade-backup/TIMESTAMP/

set -euo pipefail

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
MAGENTA='\033[0;35m'
CYAN='\033[0;36m'
BOLD='\033[1m'
NC='\033[0m' # No Color

# Script directory (where upgrade-sf3.sh lives)
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# Template directory (parent of tools/)
TEMPLATE_DIR="$(dirname "$SCRIPT_DIR")"

# Target is current working directory by default
TARGET_DIR="$(pwd)"

# Default options
DRY_RUN=false
VERBOSE=false
CREATE_BACKUP=true
FORCE=false

# Parse arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        --dry-run)
            DRY_RUN=true
            shift
            ;;
        --verbose)
            VERBOSE=true
            shift
            ;;
        --no-backup)
            CREATE_BACKUP=false
            shift
            ;;
        --force)
            FORCE=true
            shift
            ;;
        --help|-h)
            echo "Usage: $0 [OPTIONS]"
            echo ""
            echo "Options:"
            echo "  --dry-run        Show what would be updated without making changes"
            echo "  --verbose        Show detailed output during upgrade"
            echo "  --no-backup      Skip backup creation (NOT RECOMMENDED)"
            echo "  --force          Skip confirmation prompts"
            echo "  --help           Show this help message"
            echo ""
            echo "This script should be run from your project directory."
            exit 0
            ;;
        *)
            echo -e "${RED}Unknown option: $1${NC}"
            echo "Use --help for usage information"
            exit 1
            ;;
    esac
done

# Banner
echo -e "${CYAN}${BOLD}"
cat << "EOF"
╔═══════════════════════════════════════════════════════════════════╗
║                                                                   ║
║   ███████╗███████╗    ██████╗     ██████╗                        ║
║   ██╔════╝██╔════╝    ╚════██╗   ██╔═████╗                       ║
║   ███████╗█████╗       █████╔╝   ██║██╔██║                       ║
║   ╚════██║██╔══╝      ██╔═══╝    ████╔╝██║                       ║
║   ███████║██║         ███████╗██╗╚██████╔╝                       ║
║   ╚══════╝╚═╝         ╚══════╝╚═╝ ╚═════╝                        ║
║                                                                   ║
║              SOFTWARE FACTORY 3.0 UPGRADE TOOL                   ║
║         Update to Latest Rules & State Machine Fixes             ║
║                                                                   ║
╚═══════════════════════════════════════════════════════════════════╝
EOF
echo -e "${NC}"

# Function to log verbose output
log_verbose() {
    if [ "$VERBOSE" = true ]; then
        echo -e "${CYAN}[VERBOSE]${NC} $1"
    fi
}

# Function to check if this is an SF 3.0 project
validate_sf3_project() {
    echo -e "${BLUE}═══════════════════════════════════════════════════════${NC}"
    echo -e "${BOLD}Validating Software Factory 3.0 Project${NC}"
    echo -e "${BLUE}═══════════════════════════════════════════════════════${NC}"

    local is_sf3=true
    local warnings=()

    # Check for SF 3.0 state file
    if [ ! -f "$TARGET_DIR/orchestrator-state-v3.json" ]; then
        echo -e "${RED}✗ orchestrator-state-v3.json not found${NC}"
        is_sf3=false
    else
        echo -e "${GREEN}✓ orchestrator-state-v3.json found${NC}"

        # Check state machine version
        local version=$(jq -r '.state_machine.version // "unknown"' "$TARGET_DIR/orchestrator-state-v3.json" 2>/dev/null || echo "unknown")
        if [[ "$version" == "3.0"* ]]; then
            echo -e "${GREEN}✓ State machine version: $version${NC}"
        else
            echo -e "${YELLOW}⚠ State machine version: $version (will be updated)${NC}"
            warnings+=("State machine version needs update")
        fi
    fi

    # Check for SF 3.0 companion files
    if [ -f "$TARGET_DIR/bug-tracking.json" ]; then
        echo -e "${GREEN}✓ bug-tracking.json found${NC}"
    else
        echo -e "${YELLOW}⚠ bug-tracking.json not found (optional)${NC}"
    fi

    if [ -f "$TARGET_DIR/integration-containers.json" ]; then
        echo -e "${GREEN}✓ integration-containers.json found${NC}"
    else
        echo -e "${YELLOW}⚠ integration-containers.json not found (optional)${NC}"
    fi

    # Check for state machine files
    if [ -f "$TARGET_DIR/state-machines/software-factory-3.0-state-machine.json" ]; then
        echo -e "${GREEN}✓ SF 3.0 state machine found${NC}"
    else
        echo -e "${YELLOW}⚠ SF 3.0 state machine not found (will be installed)${NC}"
        warnings+=("State machine will be installed")
    fi

    # Check for State Manager agent
    if [ -f "$TARGET_DIR/.claude/agents/state-manager.md" ]; then
        echo -e "${GREEN}✓ State Manager agent found${NC}"
    else
        echo -e "${YELLOW}⚠ State Manager agent not found (will be installed)${NC}"
        warnings+=("State Manager agent will be installed")
    fi

    if [ "$is_sf3" = false ]; then
        echo -e "\n${RED}ERROR: This does not appear to be a Software Factory 3.0 project${NC}"
        echo -e "${YELLOW}Required: orchestrator-state-v3.json${NC}"
        exit 1
    fi

    if [ ${#warnings[@]} -gt 0 ]; then
        echo -e "\n${YELLOW}Warnings:${NC}"
        for warning in "${warnings[@]}"; do
            echo -e "  ${YELLOW}• $warning${NC}"
        done
    fi

    echo -e "${BLUE}═══════════════════════════════════════════════════════${NC}"
    echo ""
}

# Function to create comprehensive backup
create_backup() {
    if [ "$DRY_RUN" = true ]; then
        echo -e "${CYAN}[DRY RUN] Would create backup at: .upgrade-backup/$(date +%Y%m%d-%H%M%S)${NC}"
        return
    fi

    local backup_timestamp=$(date +%Y%m%d-%H%M%S)
    local backup_dir=".upgrade-backup/$backup_timestamp"

    echo -e "${CYAN}Creating comprehensive backup...${NC}"
    mkdir -p "$backup_dir"

    # Backup critical SF 3.0 files
    local critical_files=(
        "orchestrator-state-v3.json"
        "bug-tracking.json"
        "integration-containers.json"
        "fix-cascade-state.json"
    )

    for file in "${critical_files[@]}"; do
        if [ -f "$TARGET_DIR/$file" ]; then
            cp "$TARGET_DIR/$file" "$backup_dir/" 2>/dev/null || true
            log_verbose "Backed up: $file"
        fi
    done

    # Backup configuration directories
    local config_dirs=(
        ".claude"
        "agent-states"
        "rule-library"
        "state-machines"
        "tools"
        "utilities"
    )

    for dir in "${config_dirs[@]}"; do
        if [ -d "$TARGET_DIR/$dir" ]; then
            cp -r "$TARGET_DIR/$dir" "$backup_dir/" 2>/dev/null || true
            log_verbose "Backed up: $dir/"
        fi
    done

    # Special handling for efforts with git history
    if [ -d "$TARGET_DIR/efforts" ]; then
        echo -e "${CYAN}Backing up efforts with complete git history...${NC}"
        local efforts_backup="$backup_dir/efforts-backup"
        mkdir -p "$efforts_backup"

        # Create manifest
        echo "Effort Backup Manifest - $backup_timestamp" > "$efforts_backup/manifest.txt"
        echo "================================" >> "$efforts_backup/manifest.txt"

        # Find all git repositories in efforts
        find "$TARGET_DIR/efforts" -type d -name ".git" | while read -r gitdir; do
            local effort_dir=$(dirname "$gitdir")
            local rel_path="${effort_dir#$TARGET_DIR/efforts/}"

            echo "  Backing up: efforts/$rel_path"
            mkdir -p "$efforts_backup/$(dirname "$rel_path")"
            cp -a "$effort_dir" "$efforts_backup/$rel_path"

            # Add to manifest
            echo "" >> "$efforts_backup/manifest.txt"
            echo "Effort: $rel_path" >> "$efforts_backup/manifest.txt"
            if cd "$effort_dir" 2>/dev/null; then
                echo "  Branch: $(git branch --show-current 2>/dev/null || echo 'unknown')" >> "$efforts_backup/manifest.txt"
                echo "  Commit: $(git rev-parse --short HEAD 2>/dev/null || echo 'no commits')" >> "$efforts_backup/manifest.txt"
                cd - > /dev/null
            fi
        done
    fi

    echo -e "${GREEN}✓ Backup created at: $backup_dir${NC}"
    echo -e "${YELLOW}  Keep this backup until you verify the upgrade was successful${NC}"
    echo ""
}

# Function to update with diff reporting
update_file() {
    local src_file="$1"
    local dest_file="$2"
    local file_type="$3"

    if [ "$DRY_RUN" = true ]; then
        if [ -f "$dest_file" ]; then
            if ! diff -q "$src_file" "$dest_file" > /dev/null 2>&1; then
                echo -e "${CYAN}[DRY RUN] Would update: $file_type${NC}"
                if [ "$VERBOSE" = true ]; then
                    echo "  Changes:"
                    diff -u "$dest_file" "$src_file" | head -20 || true
                fi
            fi
        else
            echo -e "${CYAN}[DRY RUN] Would create: $file_type${NC}"
        fi
        return
    fi

    # Create parent directory if needed
    mkdir -p "$(dirname "$dest_file")"

    # Check if file changed
    if [ -f "$dest_file" ]; then
        if ! diff -q "$src_file" "$dest_file" > /dev/null 2>&1; then
            cp -f "$src_file" "$dest_file"
            echo -e "  ${GREEN}Updated:${NC} $file_type"
            return 0
        else
            log_verbose "  Unchanged: $file_type"
            return 1
        fi
    else
        cp -f "$src_file" "$dest_file"
        echo -e "  ${GREEN}Created:${NC} $file_type"
        return 0
    fi
}

# Function to update critical SF 3.0 files
update_critical_sf3_files() {
    echo -e "\n${BOLD}Updating Critical SF 3.0 Files${NC}"
    echo -e "${BLUE}────────────────────────────────────────${NC}"

    local files_updated=0

    # State machine (MOST CRITICAL - has our bug fixes!)
    echo -e "\n${CYAN}State Machine:${NC}"
    if update_file \
        "$TEMPLATE_DIR/state-machines/software-factory-3.0-state-machine.json" \
        "$TARGET_DIR/state-machines/software-factory-3.0-state-machine.json" \
        "software-factory-3.0-state-machine.json"; then
        ((files_updated++))
        echo -e "    ${YELLOW}→ Contains integration hierarchy fixes!${NC}"
        echo -e "    ${YELLOW}→ Contains State Manager validation fixes!${NC}"
    fi

    # State Manager agent
    echo -e "\n${CYAN}State Manager Agent:${NC}"
    if update_file \
        "$TEMPLATE_DIR/.claude/agents/state-manager.md" \
        "$TARGET_DIR/.claude/agents/state-manager.md" \
        "state-manager.md"; then
        ((files_updated++))
        echo -e "    ${YELLOW}→ Critical for SF 3.0 state transitions!${NC}"
    fi

    # Orchestrator agent (has State Manager integration)
    if update_file \
        "$TEMPLATE_DIR/.claude/agents/orchestrator.md" \
        "$TARGET_DIR/.claude/agents/orchestrator.md" \
        "orchestrator.md"; then
        ((files_updated++))
    fi

    # Critical rules
    echo -e "\n${CYAN}Critical Rules:${NC}"
    local critical_rules=(
        "R288-state-manager-consultation-protocol.md"
        "R322-integration-hierarchy-mandatory.md"
        "R405-orchestrator-state-transitions-only.md"
        "R506-ABSOLUTE-PROHIBITION-PRE-COMMIT-BYPASS-SUPREME-LAW.md"
    )

    for rule in "${critical_rules[@]}"; do
        if update_file \
            "$TEMPLATE_DIR/rule-library/$rule" \
            "$TARGET_DIR/rule-library/$rule" \
            "$rule"; then
            ((files_updated++))

            case "$rule" in
                *R288*)
                    echo -e "    ${YELLOW}→ State Manager consultation requirements${NC}"
                    ;;
                *R322*)
                    echo -e "    ${YELLOW}→ Integration hierarchy enforcement${NC}"
                    ;;
                *R405*)
                    echo -e "    ${YELLOW}→ State transition automation${NC}"
                    ;;
                *R506*)
                    echo -e "    ${YELLOW}→ Pre-commit bypass prohibition${NC}"
                    ;;
            esac
        fi
    done

    echo -e "\n${GREEN}✓ Updated $files_updated critical files${NC}"
}

# Function to update all rule library files
update_rule_library() {
    echo -e "\n${BOLD}Updating Rule Library${NC}"
    echo -e "${BLUE}────────────────────────────────────────${NC}"

    if [ "$DRY_RUN" = true ]; then
        echo -e "${CYAN}[DRY RUN] Would sync all rules from template${NC}"
        return
    fi

    mkdir -p "$TARGET_DIR/rule-library"

    local rules_updated=0
    local rules_created=0

    for rule_file in "$TEMPLATE_DIR/rule-library"/*.md; do
        if [ -f "$rule_file" ]; then
            local basename=$(basename "$rule_file")
            local dest_file="$TARGET_DIR/rule-library/$basename"

            if [ -f "$dest_file" ]; then
                if ! diff -q "$rule_file" "$dest_file" > /dev/null 2>&1; then
                    cp -f "$rule_file" "$dest_file"
                    ((rules_updated++))
                    log_verbose "  Updated: $basename"
                fi
            else
                cp -f "$rule_file" "$dest_file"
                ((rules_created++))
                log_verbose "  Created: $basename"
            fi
        fi
    done

    echo -e "${GREEN}✓ Rules: $rules_updated updated, $rules_created created${NC}"
}

# Function to update agent states
update_agent_states() {
    echo -e "\n${BOLD}Updating Agent State Rules${NC}"
    echo -e "${BLUE}────────────────────────────────────────${NC}"

    if [ "$DRY_RUN" = true ]; then
        echo -e "${CYAN}[DRY RUN] Would sync agent-states directory${NC}"
        return
    fi

    # Use rsync for efficient sync
    if command -v rsync &> /dev/null; then
        rsync -av --delete \
            "$TEMPLATE_DIR/agent-states/" \
            "$TARGET_DIR/agent-states/" \
            --quiet
        echo -e "${GREEN}✓ Agent states synchronized${NC}"
    else
        # Fallback to cp
        rm -rf "$TARGET_DIR/agent-states.old" 2>/dev/null || true
        [ -d "$TARGET_DIR/agent-states" ] && mv "$TARGET_DIR/agent-states" "$TARGET_DIR/agent-states.old"
        cp -r "$TEMPLATE_DIR/agent-states" "$TARGET_DIR/"
        echo -e "${GREEN}✓ Agent states updated${NC}"
    fi

    # Count state rules
    local state_count=$(find "$TARGET_DIR/agent-states" -name "rules.md" 2>/dev/null | wc -l)
    echo -e "  ${CYAN}Total state rules: $state_count${NC}"
}

# Function to update utilities and tools
update_utilities_tools() {
    echo -e "\n${BOLD}Updating Utilities and Tools${NC}"
    echo -e "${BLUE}────────────────────────────────────────${NC}"

    if [ "$DRY_RUN" = true ]; then
        echo -e "${CYAN}[DRY RUN] Would update utilities and tools${NC}"
        return
    fi

    # Update utilities
    mkdir -p "$TARGET_DIR/utilities"
    for util_file in "$TEMPLATE_DIR/utilities"/*.sh "$TEMPLATE_DIR/utilities"/*.py; do
        if [ -f "$util_file" ]; then
            local basename=$(basename "$util_file")
            cp -f "$util_file" "$TARGET_DIR/utilities/$basename"
            chmod +x "$TARGET_DIR/utilities/$basename"
            log_verbose "  Updated: utilities/$basename"
        fi
    done

    # Update tools (including line-counter.sh)
    mkdir -p "$TARGET_DIR/tools"
    for tool_file in "$TEMPLATE_DIR/tools"/*.sh; do
        if [ -f "$tool_file" ]; then
            local basename=$(basename "$tool_file")
            # Skip this script itself
            if [ "$basename" != "upgrade-sf3.sh" ]; then
                cp -f "$tool_file" "$TARGET_DIR/tools/$basename"
                chmod +x "$TARGET_DIR/tools/$basename"
                log_verbose "  Updated: tools/$basename"
            fi
        fi
    done

    # Special check for line-counter.sh
    if [ -f "$TARGET_DIR/tools/line-counter.sh" ]; then
        echo -e "${GREEN}✓ line-counter.sh updated${NC}"
    fi

    echo -e "${GREEN}✓ Utilities and tools updated${NC}"
}

# Function to validate the upgrade
validate_upgrade() {
    echo -e "\n${BOLD}Validating Upgrade${NC}"
    echo -e "${BLUE}════════════════════════════════════════${NC}"

    local validation_passed=true

    # Check state machine
    if [ -f "$TARGET_DIR/state-machines/software-factory-3.0-state-machine.json" ]; then
        echo -e "${GREEN}✓ SF 3.0 state machine present${NC}"

        # Validate it has integration hierarchy states
        if jq -e '.integration_hierarchy' "$TARGET_DIR/state-machines/software-factory-3.0-state-machine.json" > /dev/null 2>&1; then
            echo -e "${GREEN}✓ Integration hierarchy defined${NC}"
        else
            echo -e "${RED}✗ Integration hierarchy missing${NC}"
            validation_passed=false
        fi
    else
        echo -e "${RED}✗ SF 3.0 state machine missing${NC}"
        validation_passed=false
    fi

    # Check State Manager agent
    if [ -f "$TARGET_DIR/.claude/agents/state-manager.md" ]; then
        echo -e "${GREEN}✓ State Manager agent present${NC}"
    else
        echo -e "${RED}✗ State Manager agent missing${NC}"
        validation_passed=false
    fi

    # Check critical rules
    if [ -f "$TARGET_DIR/rule-library/R288-state-manager-consultation-protocol.md" ]; then
        echo -e "${GREEN}✓ R288 (State Manager consultation) present${NC}"
    else
        echo -e "${RED}✗ R288 missing${NC}"
        validation_passed=false
    fi

    if [ -f "$TARGET_DIR/rule-library/R322-integration-hierarchy-mandatory.md" ]; then
        echo -e "${GREEN}✓ R322 (Integration hierarchy) present${NC}"
    else
        echo -e "${RED}✗ R322 missing${NC}"
        validation_passed=false
    fi

    # Check orchestrator state version
    if [ -f "$TARGET_DIR/orchestrator-state-v3.json" ]; then
        local version=$(jq -r '.state_machine.version // "unknown"' "$TARGET_DIR/orchestrator-state-v3.json" 2>/dev/null)
        if [[ "$version" == "3.0"* ]]; then
            echo -e "${GREEN}✓ State file version: $version${NC}"
        else
            echo -e "${YELLOW}⚠ State file version needs manual update to 3.0.0${NC}"
            echo -e "  Run: jq '.state_machine.version = \"3.0.0\"' orchestrator-state-v3.json > tmp && mv tmp orchestrator-state-v3.json"
        fi
    fi

    if [ "$validation_passed" = true ]; then
        echo -e "\n${GREEN}${BOLD}✓ All validations passed!${NC}"
    else
        echo -e "\n${RED}${BOLD}✗ Some validations failed - manual intervention needed${NC}"
    fi

    echo -e "${BLUE}════════════════════════════════════════${NC}"
}

# Function to generate upgrade report
generate_report() {
    echo -e "\n${BOLD}Upgrade Report${NC}"
    echo -e "${BLUE}════════════════════════════════════════${NC}"

    echo -e "\n${CYAN}What was updated:${NC}"
    echo -e "  • SF 3.0 state machine with integration fixes"
    echo -e "  • State Manager agent configuration"
    echo -e "  • Rule library (all rules including R288, R322, R405, R506)"
    echo -e "  • Agent state rules for all agents"
    echo -e "  • Utilities and tools"

    echo -e "\n${CYAN}What was preserved:${NC}"
    echo -e "  • orchestrator-state-v3.json (your project state)"
    echo -e "  • efforts/ directory (all your work)"
    echo -e "  • todos/ directory"
    echo -e "  • All project-specific files"

    echo -e "\n${CYAN}Critical improvements:${NC}"
    echo -e "  ${GREEN}• Integration hierarchy enforcement fixed${NC}"
    echo -e "  ${GREEN}• State Manager validation at BUILD_VALIDATION level${NC}"
    echo -e "  ${GREEN}• Enhanced test coverage for state transitions${NC}"
    echo -e "  ${GREEN}• R288 State Manager consultation protocol${NC}"
    echo -e "  ${GREEN}• R322 Integration hierarchy mandatory${NC}"

    echo -e "\n${YELLOW}Recommended next steps:${NC}"
    echo -e "  1. Review the State Manager agent configuration"
    echo -e "  2. Verify integration hierarchy in your current state"
    echo -e "  3. Run state validation: tools/validate-state-file.sh"
    echo -e "  4. Continue with /continue-software-factory"

    if [ "$CREATE_BACKUP" = true ] && [ "$DRY_RUN" = false ]; then
        echo -e "\n${CYAN}Backup location:${NC}"
        echo -e "  .upgrade-backup/[timestamp]/"
        echo -e "  ${YELLOW}Keep this backup until you verify everything works${NC}"
    fi

    # Create version marker
    if [ "$DRY_RUN" = false ]; then
        local template_commit=$(cd "$TEMPLATE_DIR" && git rev-parse HEAD 2>/dev/null || echo "unknown")
        local template_branch=$(cd "$TEMPLATE_DIR" && git branch --show-current 2>/dev/null || echo "unknown")

        {
            echo "version: 3.0.0"
            echo "upgraded_at: $(date +%Y%m%d-%H%M%S)"
            echo "upgraded_from: $TEMPLATE_DIR"
            echo "template_commit: $template_commit"
            echo "template_branch: $template_branch"
            echo "upgrade_tool: upgrade-sf3.sh"
            echo "critical_fixes:"
            echo "  - integration_hierarchy_enforcement"
            echo "  - state_manager_validation"
            echo "  - build_validation_level_checks"
        } > "$TARGET_DIR/.sf-version"

        echo -e "\n${CYAN}Version info:${NC}"
        echo -e "  Template commit: ${GREEN}$template_commit${NC}"
        echo -e "  Template branch: ${GREEN}$template_branch${NC}"
    fi

    echo -e "${BLUE}════════════════════════════════════════${NC}"
}

# Main upgrade process
main() {
    # Validate this is an SF 3.0 project
    validate_sf3_project

    # Confirmation prompt
    if [ "$FORCE" != true ] && [ "$DRY_RUN" != true ]; then
        echo -e "${YELLOW}This will upgrade your SF 3.0 project with latest fixes.${NC}"
        echo -e "${YELLOW}Your work will be preserved and backed up.${NC}"
        echo ""
        echo -ne "${CYAN}Proceed with upgrade? (y/n): ${NC}"
        read -r response
        if [ "$response" != "y" ]; then
            echo -e "${RED}Upgrade cancelled${NC}"
            exit 1
        fi
        echo ""
    fi

    # Create backup if requested
    if [ "$CREATE_BACKUP" = true ] && [ "$DRY_RUN" = false ]; then
        create_backup
    fi

    echo -e "${BOLD}Starting SF 3.0 Upgrade...${NC}"
    echo -e "${BLUE}════════════════════════════════════════${NC}"

    # Update critical SF 3.0 files first
    update_critical_sf3_files

    # Update all rules
    update_rule_library

    # Update agent states
    update_agent_states

    # Update utilities and tools
    update_utilities_tools

    # Update other agent configs
    echo -e "\n${BOLD}Updating Other Agent Configurations${NC}"
    echo -e "${BLUE}────────────────────────────────────────${NC}"

    if [ "$DRY_RUN" = false ]; then
        for agent_file in "$TEMPLATE_DIR/.claude/agents"/*.md; do
            if [ -f "$agent_file" ]; then
                local basename=$(basename "$agent_file")
                if [ "$basename" != "state-manager.md" ] && [ "$basename" != "orchestrator.md" ]; then
                    update_file "$agent_file" "$TARGET_DIR/.claude/agents/$basename" "$basename"
                fi
            fi
        done
    else
        echo -e "${CYAN}[DRY RUN] Would update all agent configurations${NC}"
    fi

    # Update commands
    echo -e "\n${BOLD}Updating Command Configurations${NC}"
    echo -e "${BLUE}────────────────────────────────────────${NC}"

    if [ "$DRY_RUN" = false ]; then
        mkdir -p "$TARGET_DIR/.claude/commands"
        for cmd_file in "$TEMPLATE_DIR/.claude/commands"/*.md; do
            if [ -f "$cmd_file" ]; then
                local basename=$(basename "$cmd_file")
                update_file "$cmd_file" "$TARGET_DIR/.claude/commands/$basename" "$basename"
            fi
        done
    else
        echo -e "${CYAN}[DRY RUN] Would update command configurations${NC}"
    fi

    # Validate the upgrade
    if [ "$DRY_RUN" = false ]; then
        validate_upgrade
    fi

    # Generate report
    generate_report

    if [ "$DRY_RUN" = true ]; then
        echo -e "\n${YELLOW}This was a dry run. No changes were made.${NC}"
        echo -e "${YELLOW}Run without --dry-run to apply the upgrade.${NC}"
    else
        echo -e "\n${GREEN}${BOLD}✓ SF 3.0 Upgrade Complete!${NC}"
        echo -e "${GREEN}Your project has been updated with the latest fixes.${NC}"
    fi
}

# Run main upgrade
main