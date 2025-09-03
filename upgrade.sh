#!/bin/bash

# Software Factory 2.0 Upgrade Script
# This script updates an existing SF 2.0 instance with the latest rules and configurations
# while preserving work in progress (efforts, todos, state)
#
# Usage:
#   ./upgrade.sh /path/to/existing/instance [--config setup-config.yaml]

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

# Script directory (where the template lives)
TEMPLATE_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

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
║                    UPGRADE TOOL                                  ║
║            Update Rules While Preserving Work                    ║
║                                                                   ║
╚═══════════════════════════════════════════════════════════════════╝
EOF
echo -e "${NC}"

# Function to show usage
show_usage() {
    echo "Usage: $0 /path/to/existing/instance [OPTIONS]"
    echo ""
    echo "Options:"
    echo "  --config FILE    Use configuration file for variable substitution"
    echo "  --dry-run        Show what would be updated without making changes"
    echo "  --force          Skip confirmation prompts"
    echo "  --backup         Create backup before upgrading (default: true)"
    echo "  --no-backup      Skip backup creation"
    echo "  --help           Show this help message"
    echo ""
    echo "Example:"
    echo "  $0 /workspaces/my-project --config setup-config.yaml"
    echo "  $0 /workspaces/my-project --dry-run"
}

# Parse arguments
if [ $# -lt 1 ]; then
    show_usage
    exit 1
fi

TARGET_DIR="$1"
shift

# Default options
CONFIG_FILE=""
DRY_RUN=false
FORCE=false
CREATE_BACKUP=true

# Parse options
while [[ $# -gt 0 ]]; do
    case $1 in
        --config)
            CONFIG_FILE="$2"
            shift 2
            ;;
        --dry-run)
            DRY_RUN=true
            shift
            ;;
        --force)
            FORCE=true
            shift
            ;;
        --backup)
            CREATE_BACKUP=true
            shift
            ;;
        --no-backup)
            CREATE_BACKUP=false
            shift
            ;;
        --help|-h)
            show_usage
            exit 0
            ;;
        *)
            echo -e "${RED}Unknown option: $1${NC}"
            show_usage
            exit 1
            ;;
    esac
done

# Validate target directory
if [ ! -d "$TARGET_DIR" ]; then
    echo -e "${RED}❌ Target directory does not exist: $TARGET_DIR${NC}"
    exit 1
fi

# Check if it's a valid SF 2.0 instance
if [ ! -f "$TARGET_DIR/project-config.yaml" ] && [ ! -f "$TARGET_DIR/.sf-version" ]; then
    echo -e "${YELLOW}⚠️  Warning: Target doesn't appear to be a Software Factory instance${NC}"
    echo -e "${YELLOW}   Missing project-config.yaml or .sf-version${NC}"
    
    if [ "$FORCE" != true ]; then
        echo -ne "${CYAN}Continue anyway? (y/n): ${NC}"
        read -r response
        if [ "$response" != "y" ]; then
            echo -e "${RED}Upgrade cancelled${NC}"
            exit 1
        fi
    fi
fi

echo -e "${BLUE}═══════════════════════════════════════════════════════${NC}"
echo -e "${BOLD}Upgrade Configuration${NC}"
echo -e "${BLUE}═══════════════════════════════════════════════════════${NC}"
echo -e "Template Source: ${CYAN}$TEMPLATE_DIR${NC}"
echo -e "Target Instance: ${CYAN}$TARGET_DIR${NC}"
if [ -n "$CONFIG_FILE" ]; then
    echo -e "Config File:     ${CYAN}$CONFIG_FILE${NC}"
fi
echo -e "Dry Run:         ${CYAN}$DRY_RUN${NC}"
echo -e "Create Backup:   ${CYAN}$CREATE_BACKUP${NC}"
echo -e "${BLUE}═══════════════════════════════════════════════════════${NC}"

# Load configuration if provided
PROJECT_NAME=""
PROJECT_DESC=""
GITHUB_URL=""
TARGET_REPO_URL=""
TARGET_BASE_BRANCH="main"

if [ -n "$CONFIG_FILE" ] && [ -f "$CONFIG_FILE" ]; then
    echo -e "\n${CYAN}Loading configuration from $CONFIG_FILE...${NC}"
    
    # Source the parse-yaml.sh utility if available
    if [ -f "$TEMPLATE_DIR/utilities/parse-yaml.sh" ]; then
        source "$TEMPLATE_DIR/utilities/parse-yaml.sh"
        
        # Parse the config file
        eval $(parse_yaml "$CONFIG_FILE")
        
        # Extract values
        PROJECT_NAME="${project_name:-}"
        PROJECT_DESC="${project_description:-}"
        GITHUB_URL="${project_github_url:-}"
        TARGET_REPO_URL="${target_repository_url:-}"
        TARGET_BASE_BRANCH="${target_repository_base_branch:-main}"
        
        echo -e "${GREEN}✓ Configuration loaded${NC}"
    else
        echo -e "${YELLOW}⚠️  parse-yaml.sh not found, skipping variable substitution${NC}"
    fi
fi

# Function to create backup
create_backup() {
    if [ "$DRY_RUN" = true ]; then
        echo -e "${CYAN}[DRY RUN] Would create backup at: $TARGET_DIR.backup.$(date +%Y%m%d-%H%M%S)${NC}"
        return
    fi
    
    local backup_dir="$TARGET_DIR.backup.$(date +%Y%m%d-%H%M%S)"
    echo -e "${CYAN}Creating backup at: $backup_dir${NC}"
    
    # Create backup excluding large/temporary directories
    rsync -av --exclude='efforts/' \
              --exclude='*.git' \
              --exclude='node_modules' \
              --exclude='__pycache__' \
              --exclude='*.pyc' \
              "$TARGET_DIR/" "$backup_dir/"
    
    # Save list of efforts for reference
    if [ -d "$TARGET_DIR/efforts" ]; then
        echo "Efforts present (not backed up to save space):" > "$backup_dir/efforts-list.txt"
        find "$TARGET_DIR/efforts" -type d -name ".git" | while read -r gitdir; do
            effort_dir=$(dirname "$gitdir")
            echo "  - $effort_dir" >> "$backup_dir/efforts-list.txt"
        done
    fi
    
    echo -e "${GREEN}✓ Backup created at: $backup_dir${NC}"
}

# Function to preserve important files
preserve_files() {
    local temp_preserve_dir="/tmp/sf-upgrade-preserve-$$"
    mkdir -p "$temp_preserve_dir"
    
    echo -e "${CYAN}Preserving work in progress...${NC}"
    
    # Files/directories to preserve
    local preserve_items=(
        "orchestrator-state.yaml"
        "todos"
        "efforts"
        "checkpoints"
        "snapshots"
        "work-logs"
        "project-config.yaml"
        "target-repo-config.yaml"
        ".git"
    )
    
    for item in "${preserve_items[@]}"; do
        if [ -e "$TARGET_DIR/$item" ]; then
            echo -e "  Preserving: $item"
            cp -a "$TARGET_DIR/$item" "$temp_preserve_dir/" 2>/dev/null || true
        fi
    done
    
    echo "$temp_preserve_dir"
}

# Function to restore preserved files
restore_preserved() {
    local temp_preserve_dir="$1"
    
    if [ -d "$temp_preserve_dir" ]; then
        echo -e "${CYAN}Restoring preserved work...${NC}"
        
        for item in "$temp_preserve_dir"/*; do
            if [ -e "$item" ]; then
                local basename=$(basename "$item")
                echo -e "  Restoring: $basename"
                
                # Special handling for directories
                if [ -d "$item" ]; then
                    # Don't overwrite efforts directory, just ensure it exists
                    if [ "$basename" = "efforts" ]; then
                        if [ ! -d "$TARGET_DIR/$basename" ]; then
                            mv "$item" "$TARGET_DIR/"
                        fi
                    else
                        rm -rf "$TARGET_DIR/$basename" 2>/dev/null || true
                        mv "$item" "$TARGET_DIR/"
                    fi
                else
                    mv -f "$item" "$TARGET_DIR/"
                fi
            fi
        done
        
        rm -rf "$temp_preserve_dir"
        echo -e "${GREEN}✓ Work in progress restored${NC}"
    fi
}

# Function to update files with variable substitution
update_with_substitution() {
    local src_file="$1"
    local dest_file="$2"
    
    if [ "$DRY_RUN" = true ]; then
        echo -e "${CYAN}[DRY RUN] Would update: $dest_file${NC}"
        return
    fi
    
    # If no config, just copy
    if [ -z "$CONFIG_FILE" ]; then
        cp -f "$src_file" "$dest_file"
        return
    fi
    
    # Perform variable substitution
    cp "$src_file" "$dest_file.tmp"
    
    # Replace variables if they're set
    if [ -n "$PROJECT_NAME" ]; then
        sed -i "s/\${PROJECT_NAME}/$PROJECT_NAME/g" "$dest_file.tmp"
        sed -i "s/{{PROJECT_NAME}}/$PROJECT_NAME/g" "$dest_file.tmp"
    fi
    
    if [ -n "$TARGET_REPO_URL" ]; then
        sed -i "s|\${TARGET_REPO_URL}|$TARGET_REPO_URL|g" "$dest_file.tmp"
        sed -i "s|{{TARGET_REPO_URL}}|$TARGET_REPO_URL|g" "$dest_file.tmp"
    fi
    
    if [ -n "$TARGET_BASE_BRANCH" ]; then
        sed -i "s/\${TARGET_BASE_BRANCH}/$TARGET_BASE_BRANCH/g" "$dest_file.tmp"
        sed -i "s/{{TARGET_BASE_BRANCH}}/$TARGET_BASE_BRANCH/g" "$dest_file.tmp"
    fi
    
    mv "$dest_file.tmp" "$dest_file"
}

# Main upgrade process
main_upgrade() {
    echo -e "\n${BOLD}Starting Upgrade Process...${NC}\n"
    
    # Create backup if requested
    if [ "$CREATE_BACKUP" = true ] && [ "$DRY_RUN" = false ]; then
        create_backup
    fi
    
    # Preserve important files
    if [ "$DRY_RUN" = false ]; then
        PRESERVE_DIR=$(preserve_files)
    fi
    
    # Update state machine definition
    echo -e "${CYAN}Updating state machine definition...${NC}"
    if [ -f "$TEMPLATE_DIR/SOFTWARE-FACTORY-STATE-MACHINE.md" ]; then
        if [ "$DRY_RUN" = true ]; then
            echo -e "${CYAN}[DRY RUN] Would update SOFTWARE-FACTORY-STATE-MACHINE.md${NC}"
        else
            echo -e "  Updating: SOFTWARE-FACTORY-STATE-MACHINE.md"
            update_with_substitution "$TEMPLATE_DIR/SOFTWARE-FACTORY-STATE-MACHINE.md" "$TARGET_DIR/SOFTWARE-FACTORY-STATE-MACHINE.md"
        fi
        echo -e "${GREEN}✓ State machine definition updated${NC}"
    fi
    
    # Update rule library
    echo -e "${CYAN}Updating rule library...${NC}"
    if [ -d "$TEMPLATE_DIR/rule-library" ]; then
        if [ "$DRY_RUN" = true ]; then
            echo -e "${CYAN}[DRY RUN] Would update rule-library${NC}"
        else
            mkdir -p "$TARGET_DIR/rule-library"
            
            # Copy all rule files
            for rule_file in "$TEMPLATE_DIR/rule-library"/*.md; do
                if [ -f "$rule_file" ]; then
                    basename=$(basename "$rule_file")
                    echo -e "  Updating: rule-library/$basename"
                    update_with_substitution "$rule_file" "$TARGET_DIR/rule-library/$basename"
                fi
            done
        fi
        echo -e "${GREEN}✓ Rule library updated${NC}"
    fi
    
    # Update agent configurations
    echo -e "${CYAN}Updating agent configurations...${NC}"
    if [ -d "$TEMPLATE_DIR/.claude/agents" ]; then
        if [ "$DRY_RUN" = true ]; then
            echo -e "${CYAN}[DRY RUN] Would update .claude/agents${NC}"
        else
            mkdir -p "$TARGET_DIR/.claude/agents"
            
            for agent_file in "$TEMPLATE_DIR/.claude/agents"/*.md; do
                if [ -f "$agent_file" ]; then
                    basename=$(basename "$agent_file")
                    echo -e "  Updating: .claude/agents/$basename"
                    update_with_substitution "$agent_file" "$TARGET_DIR/.claude/agents/$basename"
                fi
            done
            
            # Also update global ~/.claude/agents
            echo -e "${CYAN}Updating global agent configurations...${NC}"
            mkdir -p "$HOME/.claude/agents"
            for agent_file in "$TARGET_DIR/.claude/agents"/*.md; do
                if [ -f "$agent_file" ]; then
                    basename=$(basename "$agent_file")
                    cp -f "$agent_file" "$HOME/.claude/agents/$basename"
                    echo -e "  Global: ~/.claude/agents/$basename"
                fi
            done
        fi
        echo -e "${GREEN}✓ Agent configurations updated${NC}"
    fi
    
    # Update CLAUDE.md configuration file
    echo -e "${CYAN}Updating CLAUDE.md configuration...${NC}"
    if [ -f "$TEMPLATE_DIR/.claude/CLAUDE.md" ]; then
        if [ "$DRY_RUN" = true ]; then
            echo -e "${CYAN}[DRY RUN] Would update .claude/CLAUDE.md${NC}"
        else
            mkdir -p "$TARGET_DIR/.claude"
            echo -e "  Updating: .claude/CLAUDE.md"
            update_with_substitution "$TEMPLATE_DIR/.claude/CLAUDE.md" "$TARGET_DIR/.claude/CLAUDE.md"
            echo -e "${GREEN}✓ CLAUDE.md updated (grading criteria, TODO rules, agent configs)${NC}"
        fi
    fi
    
    # Update command configurations
    echo -e "${CYAN}Updating command configurations...${NC}"
    if [ -d "$TEMPLATE_DIR/.claude/commands" ]; then
        if [ "$DRY_RUN" = true ]; then
            echo -e "${CYAN}[DRY RUN] Would update .claude/commands${NC}"
        else
            mkdir -p "$TARGET_DIR/.claude/commands"
            
            for cmd_file in "$TEMPLATE_DIR/.claude/commands"/*.md; do
                if [ -f "$cmd_file" ]; then
                    basename=$(basename "$cmd_file")
                    echo -e "  Updating: .claude/commands/$basename"
                    update_with_substitution "$cmd_file" "$TARGET_DIR/.claude/commands/$basename"
                fi
            done
            
            # Also update global ~/.claude/commands
            echo -e "${CYAN}Updating global command configurations...${NC}"
            mkdir -p "$HOME/.claude/commands"
            for cmd_file in "$TARGET_DIR/.claude/commands"/*.md; do
                if [ -f "$cmd_file" ]; then
                    basename=$(basename "$cmd_file")
                    cp -f "$cmd_file" "$HOME/.claude/commands/$basename"
                    echo -e "  Global: ~/.claude/commands/$basename"
                fi
            done
        fi
        echo -e "${GREEN}✓ Command configurations updated${NC}"
    fi
    
    # Update utilities
    echo -e "${CYAN}Updating utility scripts...${NC}"
    if [ -d "$TEMPLATE_DIR/utilities" ]; then
        if [ "$DRY_RUN" = true ]; then
            echo -e "${CYAN}[DRY RUN] Would update utilities${NC}"
        else
            mkdir -p "$TARGET_DIR/utilities"
            
            for util_file in "$TEMPLATE_DIR/utilities"/*.sh; do
                if [ -f "$util_file" ]; then
                    basename=$(basename "$util_file")
                    echo -e "  Updating: utilities/$basename"
                    cp -f "$util_file" "$TARGET_DIR/utilities/$basename"
                    chmod +x "$TARGET_DIR/utilities/$basename"
                fi
            done
            
            # Also update in ~/.claude/utilities if it exists
            if [ -d "$HOME/.claude/utilities" ]; then
                echo -e "${CYAN}Updating global utilities...${NC}"
                for util_file in "$TEMPLATE_DIR/utilities"/*.sh; do
                    if [ -f "$util_file" ]; then
                        basename=$(basename "$util_file")
                        cp -f "$util_file" "$HOME/.claude/utilities/$basename"
                        chmod +x "$HOME/.claude/utilities/$basename"
                    fi
                done
            fi
        fi
        echo -e "${GREEN}✓ Utilities updated${NC}"
    fi
    
    # Update tools directory (line-counter.sh and other tools)
    echo -e "${CYAN}Updating tools...${NC}"
    if [ -d "$TEMPLATE_DIR/tools" ]; then
        if [ "$DRY_RUN" = true ]; then
            echo -e "${CYAN}[DRY RUN] Would update tools${NC}"
        else
            # Check if tools directory exists and has files
            if [ -d "$TARGET_DIR/tools" ] && [ "$(ls -A $TARGET_DIR/tools 2>/dev/null)" ]; then
                if [ "$FORCE" = true ]; then
                    echo -e "  ${YELLOW}Overwriting existing tools...${NC}"
                    rm -rf "$TARGET_DIR/tools"
                    mkdir -p "$TARGET_DIR/tools"
                else
                    echo -e "${YELLOW}⚠️  Tools directory exists with files${NC}"
                    echo -ne "${CYAN}Overwrite existing tools? (y/n): ${NC}"
                    read -r overwrite_tools
                    if [ "$overwrite_tools" = "y" ]; then
                        echo -e "  ${YELLOW}Backing up old tools to tools.bak...${NC}"
                        rm -rf "$TARGET_DIR/tools.bak"
                        mv "$TARGET_DIR/tools" "$TARGET_DIR/tools.bak"
                        mkdir -p "$TARGET_DIR/tools"
                    else
                        echo -e "  ${CYAN}Merging with existing tools...${NC}"
                    fi
                fi
            else
                mkdir -p "$TARGET_DIR/tools"
            fi
            
            # Copy all tool files
            for tool_file in "$TEMPLATE_DIR/tools"/*; do
                if [ -f "$tool_file" ]; then
                    basename=$(basename "$tool_file")
                    echo -e "  Updating: tools/$basename"
                    cp -f "$tool_file" "$TARGET_DIR/tools/$basename"
                    # Make shell scripts executable
                    if [[ "$basename" == *.sh ]]; then
                        chmod +x "$TARGET_DIR/tools/$basename"
                    fi
                fi
            done
            
            # Special handling for line-counter.sh
            if [ -f "$TARGET_DIR/tools/line-counter.sh" ]; then
                echo -e "  ${GREEN}✓ line-counter.sh updated with fix for parameter parsing${NC}"
            fi
        fi
        echo -e "${GREEN}✓ Tools updated${NC}"
    fi
    
    # Update agent-states directory (state machine rules)
    echo -e "${CYAN}Updating agent state rules...${NC}"
    if [ -d "$TEMPLATE_DIR/agent-states" ]; then
        if [ "$DRY_RUN" = true ]; then
            echo -e "${CYAN}[DRY RUN] Would update agent-states${NC}"
        else
            # Recursively copy the entire agent-states structure
            # This includes all agent types and their state-specific rules
            echo -e "  Syncing agent-states directory structure..."
            
            # Create base directory
            mkdir -p "$TARGET_DIR/agent-states"
            
            # Copy orchestrator states
            if [ -d "$TEMPLATE_DIR/agent-states/orchestrator" ]; then
                echo -e "  Updating orchestrator states..."
                cp -r "$TEMPLATE_DIR/agent-states/orchestrator" "$TARGET_DIR/agent-states/"
            fi
            
            # Copy sw-engineer states
            if [ -d "$TEMPLATE_DIR/agent-states/sw-engineer" ]; then
                echo -e "  Updating sw-engineer states..."
                cp -r "$TEMPLATE_DIR/agent-states/sw-engineer" "$TARGET_DIR/agent-states/"
            fi
            
            # Copy code-reviewer states
            if [ -d "$TEMPLATE_DIR/agent-states/code-reviewer" ]; then
                echo -e "  Updating code-reviewer states..."
                cp -r "$TEMPLATE_DIR/agent-states/code-reviewer" "$TARGET_DIR/agent-states/"
            fi
            
            # Copy architect states
            if [ -d "$TEMPLATE_DIR/agent-states/architect" ]; then
                echo -e "  Updating architect states..."
                cp -r "$TEMPLATE_DIR/agent-states/architect" "$TARGET_DIR/agent-states/"
            fi
            
            # Count total state directories updated
            state_count=$(find "$TARGET_DIR/agent-states" -name "rules.md" | wc -l)
            echo -e "  Updated $state_count state rule files"
        fi
        echo -e "${GREEN}✓ Agent state rules updated${NC}"
    fi
    
    # Update critical files
    echo -e "${CYAN}Updating critical files...${NC}"
    critical_files=(
        "rule-library/R002-agent-acknowledgment.md"
        "rule-library/R003-performance-grading.md"
        "DELEGATION-AND-AGENT-PROMPT-GUIDE.md"
        "CRITICAL-IMPLEMENTATION-GUIDANCE.md"
        "AGENT-ACKNOWLEDGMENT-EXAMPLES.md"
        "SOFTWARE-FACTORY-STATE-MACHINE.md"
    )
    
    for file_path in "${critical_files[@]}"; do
        if [ -f "$TEMPLATE_DIR/$file_path" ]; then
            if [ "$DRY_RUN" = true ]; then
                echo -e "${CYAN}[DRY RUN] Would update: $file_path${NC}"
            else
                # Create directory if needed
                dir_path=$(dirname "$file_path")
                if [ "$dir_path" != "." ]; then
                    mkdir -p "$TARGET_DIR/$dir_path"
                fi
                
                echo -e "  Updating: $file_path"
                update_with_substitution "$TEMPLATE_DIR/$file_path" "$TARGET_DIR/$file_path"
            fi
        fi
    done
    echo -e "${GREEN}✓ Critical files updated${NC}"
    
    # Update claude settings if present
    if [ -f "$TEMPLATE_DIR/.claude/settings.json" ]; then
        if [ "$DRY_RUN" = true ]; then
            echo -e "${CYAN}[DRY RUN] Would update .claude/settings.json${NC}"
        else
            echo -e "${CYAN}Updating Claude settings...${NC}"
            mkdir -p "$TARGET_DIR/.claude"
            cp -f "$TEMPLATE_DIR/.claude/settings.json" "$TARGET_DIR/.claude/settings.json"
            
            # Also update global ~/.claude/settings.json
            echo -e "${CYAN}Updating global Claude settings...${NC}"
            mkdir -p "$HOME/.claude"
            
            # Check if global settings exists and prompt
            if [ -f "$HOME/.claude/settings.json" ]; then
                if [ "$FORCE" = true ]; then
                    cp -f "$TARGET_DIR/.claude/settings.json" "$HOME/.claude/settings.json"
                    echo -e "  Global: ~/.claude/settings.json (overwritten)"
                else
                    echo -e "${YELLOW}⚠️  ~/.claude/settings.json exists${NC}"
                    echo -ne "${CYAN}Overwrite global settings? (y/n): ${NC}"
                    read -r overwrite_settings
                    if [ "$overwrite_settings" = "y" ]; then
                        cp -f "$TARGET_DIR/.claude/settings.json" "$HOME/.claude/settings.json"
                        echo -e "  Global: ~/.claude/settings.json (overwritten)"
                    else
                        echo -e "  Global: ~/.claude/settings.json (kept existing)"
                    fi
                fi
            else
                cp -f "$TARGET_DIR/.claude/settings.json" "$HOME/.claude/settings.json"
                echo -e "  Global: ~/.claude/settings.json"
            fi
            
            echo -e "${GREEN}✓ Claude settings updated${NC}"
        fi
    fi
    
    # Restore preserved files
    if [ "$DRY_RUN" = false ] && [ -n "$PRESERVE_DIR" ]; then
        restore_preserved "$PRESERVE_DIR"
    fi
    
    # Create version marker with template repo commit hash
    if [ "$DRY_RUN" = false ]; then
        echo -e "${CYAN}Creating version marker with template commit info...${NC}"
        
        # Get the current commit hash of the template repository
        TEMPLATE_COMMIT=""
        if [ -d "$TEMPLATE_DIR/.git" ]; then
            TEMPLATE_COMMIT=$(cd "$TEMPLATE_DIR" && git rev-parse HEAD 2>/dev/null || echo "unknown")
            TEMPLATE_BRANCH=$(cd "$TEMPLATE_DIR" && git branch --show-current 2>/dev/null || echo "unknown")
            TEMPLATE_REMOTE=$(cd "$TEMPLATE_DIR" && git remote get-url origin 2>/dev/null || echo "unknown")
        else
            TEMPLATE_COMMIT="not-a-git-repo"
            TEMPLATE_BRANCH="N/A"
            TEMPLATE_REMOTE="N/A"
        fi
        
        # Create the .sf-version file with comprehensive version info
        {
            echo "version: 2.0"
            echo "upgraded_at: $(date +%Y%m%d-%H%M%S)"
            echo "upgraded_from: $TEMPLATE_DIR"
            echo "template_commit: $TEMPLATE_COMMIT"
            echo "template_branch: $TEMPLATE_BRANCH"
            echo "template_remote: $TEMPLATE_REMOTE"
            echo "upgrade_date: $(date -u +"%Y-%m-%dT%H:%M:%SZ")"
        } > "$TARGET_DIR/.sf-version"
        
        echo -e "${GREEN}✓ Version marker created with commit: ${CYAN}$TEMPLATE_COMMIT${NC}"
    fi
}

# Confirmation prompt
if [ "$FORCE" != true ] && [ "$DRY_RUN" != true ]; then
    echo -e "\n${YELLOW}⚠️  This will update rules and configurations in:${NC}"
    echo -e "${YELLOW}   $TARGET_DIR${NC}"
    echo -e "\n${GREEN}The following will be preserved:${NC}"
    echo -e "  • efforts/ directory (all work in progress)"
    echo -e "  • todos/ directory"
    echo -e "  • orchestrator-state.yaml"
    echo -e "  • project-config.yaml"
    echo -e "  • target-repo-config.yaml"
    echo -e "  • .git repository"
    echo ""
    echo -ne "${CYAN}Proceed with upgrade? (y/n): ${NC}"
    read -r response
    if [ "$response" != "y" ]; then
        echo -e "${RED}Upgrade cancelled${NC}"
        exit 1
    fi
fi

# Run the upgrade
main_upgrade

# Final report
echo -e "\n${BLUE}═══════════════════════════════════════════════════════${NC}"
echo -e "${GREEN}${BOLD}✓ Upgrade Complete!${NC}"
echo -e "${BLUE}═══════════════════════════════════════════════════════${NC}"

if [ "$DRY_RUN" = true ]; then
    echo -e "${YELLOW}This was a dry run. No changes were made.${NC}"
    echo -e "${YELLOW}Run without --dry-run to apply changes.${NC}"
else
    echo -e "${GREEN}Your Software Factory instance has been upgraded.${NC}"
    echo -e ""
    echo -e "${CYAN}What was updated:${NC}"
    echo -e "  • Rule library (all R*.md files)"
    echo -e "  • Agent configurations (.claude/agents/)"
    echo -e "  • CLAUDE.md configuration (.claude/CLAUDE.md)"
    echo -e "  • Command configurations (.claude/commands/)"
    echo -e "  • Utility scripts (utilities/)"
    echo -e "  • Critical documentation files"
    echo -e ""
    echo -e "${CYAN}Global updates (~/.claude/):${NC}"
    echo -e "  • Agent configurations → ~/.claude/agents/"
    echo -e "  • Command configurations → ~/.claude/commands/"
    echo -e "  • Utility scripts → ~/.claude/utilities/"
    echo -e "  • Settings → ~/.claude/settings.json"
    echo -e ""
    echo -e "${CYAN}What was preserved:${NC}"
    echo -e "  • All work in efforts/"
    echo -e "  • TODO state in todos/"
    echo -e "  • Orchestrator state"
    echo -e "  • Project configurations"
    echo -e ""
    
    if [ "$CREATE_BACKUP" = true ]; then
        echo -e "${CYAN}Backup location:${NC}"
        echo -e "  $TARGET_DIR.backup.*"
    fi
    
    # Display version info if created
    if [ -f "$TARGET_DIR/.sf-version" ]; then
        echo -e "${CYAN}Version Information:${NC}"
        echo -e "  Template Commit: $(grep '^template_commit:' "$TARGET_DIR/.sf-version" | cut -d' ' -f2)"
        echo -e "  Template Branch: $(grep '^template_branch:' "$TARGET_DIR/.sf-version" | cut -d' ' -f2)"
        echo -e ""
    fi
    
    echo -e "${YELLOW}Recommended next steps:${NC}"
    echo -e "  1. Review updated rules in rule-library/"
    echo -e "  2. Check agent configurations in .claude/agents/"
    echo -e "  3. Test with a simple command to ensure everything works"
    echo -e "  4. Resume your work where you left off"
fi

echo -e "${BLUE}═══════════════════════════════════════════════════════${NC}"