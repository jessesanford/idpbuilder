#!/bin/bash

# Software Factory 3.0 Upgrade Script
# This script updates an existing SF 2.0 or SF 3.0 instance with the latest rules and configurations
# while preserving work in progress (efforts, todos, state)
#
# Supports:
#   - Software Factory 2.0 (orchestrator-state-v3.json)
#   - Software Factory 3.0 (orchestrator-state-v3.json + 4-file structure)
#   - Auto-detection of version
#   - Unified hook system (master-pre-commit.sh)
#
# CRITICAL: This script now creates COMPREHENSIVE BACKUPS of effort directories
#           including complete git history to prevent data loss from branch deletion
#
# Backup Features:
#   - Full copy of all effort directories with .git folders
#   - Preserves all branches, commits, and uncommitted changes
#   - Creates timestamped backups to prevent overwrites
#   - Generates manifest files listing all backed-up efforts
#   - Uses backup-all-efforts.sh utility when available
#   - Fallback manual backup if utility is missing
#
# Usage:
#   ./upgrade.sh /path/to/existing/instance [--config setup-config.yaml]
#
# Recovery:
#   If branches are deleted, restore from:
#   - $TARGET.backup.TIMESTAMP/efforts-backup-TIMESTAMP/
#   - $TARGET/backups/efforts-upgrade-TIMESTAMP/

set -e

# Run entire script with nice priority to reduce system load during upgrades
# All child processes (rsync, cp, git, etc.) inherit this nice level automatically
renice -n 10 $$ > /dev/null 2>&1 || true

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
MAGENTA='\033[0;35m'
CYAN='\033[0;36m'
BOLD='\033[1m'
NC='\033[0m' # No Color

# Script directory (where upgrade.sh lives)
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# Template directory (parent of tools/)
TEMPLATE_DIR="$(dirname "$SCRIPT_DIR")"

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

# Check for help first
for arg in "$@"; do
    if [[ "$arg" == "--help" || "$arg" == "-h" ]]; then
        show_usage
        exit 0
    fi
done

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

# Check if upgrading the template itself
if [ "$TEMPLATE_DIR" = "$TARGET_DIR" ]; then
    echo -e "${YELLOW}⚠️  Warning: You are running upgrade on the template directory itself${NC}"
    echo -e "${YELLOW}   Source and destination are the same: $TARGET_DIR${NC}"
    echo -e "${YELLOW}   This is typically not necessary as the template is already up-to-date.${NC}"
    echo -e "${YELLOW}   Same-file copy operations will be skipped automatically.${NC}"

    if [ "$FORCE" != true ] && [ "$DRY_RUN" != true ]; then
        echo -ne "${CYAN}Continue anyway? (y/n): ${NC}"
        read -r response
        if [ "$response" != "y" ]; then
            echo -e "${RED}Upgrade cancelled${NC}"
            exit 1
        fi
    fi
fi

# Detect Software Factory version
SF_VERSION="unknown"
if [ -f "$TARGET_DIR/orchestrator-state-v3.json" ] || [ -f "$TARGET_DIR/bug-tracking.json" ] || [ -f "$TARGET_DIR/integration-containers.json" ]; then
    SF_VERSION="3.0"
elif [ -f "$TARGET_DIR/orchestrator-state-v3.json" ]; then
    SF_VERSION="2.0"
elif [ -f "$TARGET_DIR/project-config.yaml" ] || [ -f "$TARGET_DIR/.sf-version" ]; then
    # Has SF markers but no state file - assume 2.0
    SF_VERSION="2.0"
else
    echo -e "${YELLOW}⚠️  Warning: Target doesn't appear to be a Software Factory instance${NC}"
    echo -e "${YELLOW}   Missing state files or SF markers${NC}"

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
echo -e "SF Version:      ${CYAN}$SF_VERSION${NC}"
if [ -n "$CONFIG_FILE" ]; then
    echo -e "Config File:     ${CYAN}$CONFIG_FILE${NC}"
fi
echo -e "Dry Run:         ${CYAN}$DRY_RUN${NC}"
echo -e "${BLUE}═══════════════════════════════════════════════════════${NC}"

# Interactive backup prompt if not already specified and not in force/dry-run mode
if [ "$FORCE" != true ] && [ "$DRY_RUN" != true ]; then
    # Only prompt if backup wasn't explicitly set via flags
    if [ "$CREATE_BACKUP" = true ]; then
        echo ""
        echo -e "${YELLOW}This will upgrade your project with the latest rules from the template.${NC}"
        echo ""
        echo -ne "${CYAN}Do you want to create a backup before upgrading? (yes/no): ${NC}"
        read -r backup_response
        if [[ ! "$backup_response" =~ ^[Yy]([Ee][Ss])?$ ]]; then
            CREATE_BACKUP=false
            echo -e "${YELLOW}⏭️  Skipping backup (as requested)${NC}"
        fi
    fi
fi

# Default: don't compress unless user asks
COMPRESS_BACKUP=false

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
        echo -e "${CYAN}[DRY RUN] Would create backup in: backups/backup-$(date +%Y%m%d-%H%M%S)${NC}"
        echo -e "${CYAN}[DRY RUN] Would also backup effort directories with full git history${NC}"
        return
    fi

    local backup_timestamp=$(date +%Y%m%d-%H%M%S)
    local backup_dir="$TARGET_DIR/backups/backup-$backup_timestamp"

    echo ""
    echo -e "${CYAN}📦 Creating backup at: ${backup_dir#$TARGET_DIR/}${NC}"

    # Create backups directory if it doesn't exist
    mkdir -p "$TARGET_DIR/backups"

    # Create backup excluding large/temporary directories and the backups folder itself
    # CRITICAL: --exclude='backups/' prevents recursive backup-of-backups!
    ionice -c 2 -n 7 nice -n 10 rsync -av \
              --exclude='backups/' \
              --exclude='efforts/' \
              --exclude='.backup.*' \
              --exclude='*.backup-*' \
              --exclude='.git/' \
              --exclude='node_modules/' \
              --exclude='vendor/' \
              --exclude='__pycache__/' \
              --exclude='*.pyc' \
              --exclude='*.log' \
              "$TARGET_DIR/" "$backup_dir/"

    # CRITICAL: Backup effort directories with complete git history
    if [ -d "$TARGET_DIR/efforts" ]; then
        echo -e "${CYAN}Creating comprehensive effort backup...${NC}"
        echo -e "${YELLOW}⚠️  This preserves complete git history for disaster recovery${NC}"

        # Create efforts backup directory
        local efforts_backup_dir="$backup_dir/efforts-backup-$backup_timestamp"
        mkdir -p "$efforts_backup_dir"

        # Use the backup-all-efforts.sh utility if available
        if [ -f "$TARGET_DIR/utilities/backup-all-efforts.sh" ]; then
            echo -e "${CYAN}Using backup-all-efforts.sh utility...${NC}"
            (cd "$TARGET_DIR" && bash utilities/backup-all-efforts.sh "upgrade-$backup_timestamp")

            # Note: We no longer copy existing backups to prevent recursive backup growth
            # The backup-all-efforts.sh utility creates its own backup in $TARGET_DIR/backups/
        else
            # Fallback: Manual comprehensive effort backup
            echo -e "${YELLOW}backup-all-efforts.sh not found, using manual backup...${NC}"

            # Copy entire efforts directory preserving all git data using rsync
            echo -e "${CYAN}Copying all efforts with complete git history using rsync...${NC}"
            ionice -c 2 -n 7 nice -n 10 rsync -av "$TARGET_DIR/efforts/" "$efforts_backup_dir/efforts/"

            # Create manifest of backed up efforts
            echo "Effort Backup Manifest - $(date)" > "$efforts_backup_dir/manifest.txt"
            echo "================================" >> "$efforts_backup_dir/manifest.txt"
            echo "" >> "$efforts_backup_dir/manifest.txt"

            find "$efforts_backup_dir/efforts" -type d -name ".git" | while read -r gitdir; do
                effort_dir=$(dirname "$gitdir")
                effort_name=$(basename "$effort_dir")

                echo "Effort: $effort_name" >> "$efforts_backup_dir/manifest.txt"
                echo "  Path: $effort_dir" >> "$efforts_backup_dir/manifest.txt"

                # Get current branch and commit
                if cd "$effort_dir" 2>/dev/null; then
                    branch=$(git branch --show-current 2>/dev/null || echo "unknown")
                    commit=$(git rev-parse --short HEAD 2>/dev/null || echo "no commits")
                    echo "  Branch: $branch" >> "$efforts_backup_dir/manifest.txt"
                    echo "  Latest commit: $commit" >> "$efforts_backup_dir/manifest.txt"

                    # List all branches
                    echo "  All branches:" >> "$efforts_backup_dir/manifest.txt"
                    git branch -a 2>/dev/null | sed 's/^/    /' >> "$efforts_backup_dir/manifest.txt"

                    cd - > /dev/null
                fi
                echo "" >> "$efforts_backup_dir/manifest.txt"
            done

            echo -e "${GREEN}✓ Efforts backed up to: $efforts_backup_dir${NC}"
            echo -e "${GREEN}✓ Manifest created: $efforts_backup_dir/manifest.txt${NC}"
        fi

        # Also save a quick reference list
        echo "Efforts present (full backup available in efforts-backup-$backup_timestamp/):" > "$backup_dir/efforts-list.txt"
        find "$TARGET_DIR/efforts" -type d -name ".git" | while read -r gitdir; do
            effort_dir=$(dirname "$gitdir")
            echo "  - $effort_dir" >> "$backup_dir/efforts-list.txt"
        done
    fi

    echo -e "${GREEN}✓ Backup created: ${backup_dir#$TARGET_DIR/}${NC}"
    echo ""

    # Ask about compression (interactive mode only)
    if [ "$FORCE" != true ]; then
        echo -ne "${CYAN}Do you want to compress the PROJECT FILES backup to save space? (yes/no): ${NC}"
        read -r compress_response

        if [[ "$compress_response" =~ ^[Yy]([Ee][Ss])?$ ]]; then
            COMPRESS_BACKUP=true
            echo -e "${CYAN}🗜️  Compressing PROJECT FILES backup...${NC}"

            # Change to backups directory for cleaner tar paths
            local original_dir="$(pwd)"
            cd "$TARGET_DIR/backups"

            # Compress the backup
            if tar -czf "backup-$backup_timestamp.tar.gz" "backup-$backup_timestamp" 2>/dev/null; then
                # Remove uncompressed backup
                rm -rf "backup-$backup_timestamp"

                # Show size savings
                local compressed_size=$(du -h "backup-$backup_timestamp.tar.gz" | cut -f1)
                echo -e "${GREEN}✓ Backup compressed: backups/backup-$backup_timestamp.tar.gz${NC}"
                echo -e "${CYAN}   Compressed size: $compressed_size${NC}"
            else
                echo -e "${RED}❌ ERROR: Compression failed!${NC}"
                echo -e "${YELLOW}   Keeping uncompressed backup${NC}"
                COMPRESS_BACKUP=false
            fi

            cd "$original_dir"
        else
            echo -e "${YELLOW}⏭️  Keeping backup uncompressed${NC}"
        fi
    fi

    echo ""
    echo -e "${CYAN}═══════════════════════════════════════════════════════════════════════════${NC}"
    echo -e "${BOLD}📦 BACKUP SUMMARY${NC}"
    echo -e "${CYAN}═══════════════════════════════════════════════════════════════════════════${NC}"
    echo -e "${GREEN}✓ Two separate backups created:${NC}"
    echo -e "${YELLOW}  1. EFFORTS backup (git repos with full history)${NC}"
    echo -e "     Location: backups/efforts-upgrade-*/"
    echo -e "     Contains: Complete effort working copies with all branches"
    echo -e "${YELLOW}  2. PROJECT FILES backup (configs, rules, state files)${NC}"
    echo -e "     Location: backups/backup-*/"
    echo -e "     Contains: Rules, configs, state files (excludes efforts/)"
    echo -e "${CYAN}═══════════════════════════════════════════════════════════════════════════${NC}"
}

# Function to preserve important files
preserve_files() {
    local temp_preserve_dir="/tmp/sf-upgrade-preserve-$$"
    mkdir -p "$temp_preserve_dir"
    
    echo -e "${CYAN}Preserving work in progress...${NC}"
    
    # Files/directories to preserve
    local preserve_items=(
        "orchestrator-state-v3.json"
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
            # Use rsync for reliable preservation with proper trailing slash handling
            if [ -d "$TARGET_DIR/$item" ]; then
                ionice -c 2 -n 7 nice -n 10 rsync -a "$TARGET_DIR/$item/" "$temp_preserve_dir/$item/" 2>/dev/null || true
            else
                ionice -c 2 -n 7 nice -n 10 rsync -a "$TARGET_DIR/$item" "$temp_preserve_dir/" 2>/dev/null || true
            fi
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
    
    # Migrate to JSON state machine format
    echo -e "${CYAN}Migrating to JSON state machine format...${NC}"

    if [ "$DRY_RUN" = true ]; then
        echo -e "${CYAN}[DRY RUN] Would migrate state machines to JSON:${NC}"
        [ -f "$TARGET_DIR/SOFTWARE-FACTORY-STATE-MACHINE.md" ] && echo "  - Remove: SOFTWARE-FACTORY-STATE-MACHINE.md"
        [ -f "$TARGET_DIR/SOFTWARE-FACTORY-INITIALIZATION-STATE-MACHINE.md" ] && echo "  - Remove: SOFTWARE-FACTORY-INITIALIZATION-STATE-MACHINE.md"
        echo "  - Install: software-factory-3.0-state-machine.json"
        echo "  - Install: state-machines/*.json"
        ls "$TARGET_DIR/state-machines"/*.md 2>/dev/null | while read f; do
            echo "  - Remove: state-machines/$(basename "$f")"
        done
    else
        # Remove old markdown state machines
        if [ -f "$TARGET_DIR/SOFTWARE-FACTORY-STATE-MACHINE.md" ]; then
            echo -e "  ${YELLOW}Removing old:${NC} SOFTWARE-FACTORY-STATE-MACHINE.md"
            rm -f "$TARGET_DIR/SOFTWARE-FACTORY-STATE-MACHINE.md"
        fi

        if [ -f "$TARGET_DIR/SOFTWARE-FACTORY-INITIALIZATION-STATE-MACHINE.md" ]; then
            echo -e "  ${YELLOW}Removing old:${NC} SOFTWARE-FACTORY-INITIALIZATION-STATE-MACHINE.md"
            rm -f "$TARGET_DIR/SOFTWARE-FACTORY-INITIALIZATION-STATE-MACHINE.md"
        fi

        # Main state machine is in state-machines/, symlinked from root
        # This block is kept for backward compatibility but no longer needed
        # as the file is copied in the state-machines directory block below

        # Update state-machines directory
        if [ -d "$TEMPLATE_DIR/state-machines" ]; then
            mkdir -p "$TARGET_DIR/state-machines"

            # Copy JSON state machines
            for json_file in "$TEMPLATE_DIR/state-machines"/*.json; do
                if [ -f "$json_file" ]; then
                    basename=$(basename "$json_file")
                    dest_file="$TARGET_DIR/state-machines/$basename"

                    # Check if source and destination are the same file
                    if [ "$json_file" = "$dest_file" ]; then
                        echo -e "  ${CYAN}Skipping:${NC} state-machines/$basename (source and destination are the same)"
                    else
                        echo -e "  ${GREEN}Installing:${NC} state-machines/$basename"
                        cp -f "$json_file" "$dest_file"
                    fi
                fi
            done

            # Remove old markdown state machines
            for md_file in "$TARGET_DIR/state-machines"/*.md; do
                if [ -f "$md_file" ]; then
                    echo -e "  ${YELLOW}Removing old:${NC} state-machines/$(basename "$md_file")"
                    rm -f "$md_file"
                fi
            done

            # Clean up backups
            rm -f "$TARGET_DIR/state-machines"/*.backup
            rm -f "$TARGET_DIR/state-machines"/*.bak
            rm -f "$TARGET_DIR/state-machines"/*.pre-reconciliation
        fi

        # Symlinks no longer needed - all references updated to use state-machines/ directory

        echo -e "${GREEN}✓ State machines migrated to JSON format${NC}"
        echo -e "${YELLOW}  📝 Note: All state operations now use JSON (jq) instead of markdown (grep)${NC}"
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

            # Copy shell scripts
            for util_file in "$TEMPLATE_DIR/utilities"/*.sh; do
                if [ -f "$util_file" ]; then
                    basename=$(basename "$util_file")
                    dest_file="$TARGET_DIR/utilities/$basename"

                    # Check if source and destination are the same file
                    if [ "$util_file" = "$dest_file" ]; then
                        echo -e "  ${CYAN}Skipping:${NC} utilities/$basename (source and destination are the same)"
                    else
                        echo -e "  Updating: utilities/$basename"
                        cp -f "$util_file" "$dest_file"
                    fi
                    chmod +x "$dest_file"
                fi
            done

            # Copy Python scripts
            for py_file in "$TEMPLATE_DIR/utilities"/*.py; do
                if [ -f "$py_file" ]; then
                    basename=$(basename "$py_file")
                    dest_file="$TARGET_DIR/utilities/$basename"

                    # Check if source and destination are the same file
                    if [ "$py_file" = "$dest_file" ]; then
                        echo -e "  ${CYAN}Skipping:${NC} utilities/$basename (source and destination are the same)"
                    else
                        echo -e "  Updating: utilities/$basename"
                        cp -f "$py_file" "$dest_file"
                    fi
                    chmod +x "$dest_file"
                fi
            done
            
            # Also update global ~/.claude/utilities (always create if missing)
            echo -e "${CYAN}Updating global utilities...${NC}"
            mkdir -p "$HOME/.claude/utilities"
            for util_file in "$TEMPLATE_DIR/utilities"/*.sh; do
                if [ -f "$util_file" ]; then
                    basename=$(basename "$util_file")
                    cp -f "$util_file" "$HOME/.claude/utilities/$basename"
                    chmod +x "$HOME/.claude/utilities/$basename"
                    echo -e "  Global: ~/.claude/utilities/$basename"
                fi
            done
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
                    dest_file="$TARGET_DIR/tools/$basename"

                    # Check if source and destination are the same file
                    if [ "$tool_file" = "$dest_file" ]; then
                        echo -e "  ${CYAN}Skipping:${NC} tools/$basename (source and destination are the same)"
                    else
                        echo -e "  Updating: tools/$basename"
                        cp -f "$tool_file" "$dest_file"
                    fi

                    # Make shell scripts executable
                    if [[ "$basename" == *.sh ]]; then
                        chmod +x "$dest_file"
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
    
    # Update planning directory structure with examples
    echo -e "${CYAN}Updating planning directory structure...${NC}"
    if [ -d "$TEMPLATE_DIR/planning" ]; then
        if [ "$DRY_RUN" = true ]; then
            echo -e "${CYAN}[DRY RUN] Would update planning directory structure${NC}"
        else
            # Check if planning directory exists with user content
            if [ -d "$TARGET_DIR/planning" ] && [ "$(find $TARGET_DIR/planning -type f -name "*.md" ! -name "*.example.md" 2>/dev/null | head -1)" ]; then
                echo -e "${YELLOW}⚠️  Planning directory exists with user content${NC}"
                echo -e "  Preserving existing plans and adding new examples..."

                # Copy only example files to preserve user's actual plans
                find "$TEMPLATE_DIR/planning" -name "*.example.md" -o -name "README.md" | while read example_file; do
                    rel_path="${example_file#$TEMPLATE_DIR/planning/}"
                    dest_dir="$TARGET_DIR/planning/$(dirname "$rel_path")"
                    mkdir -p "$dest_dir"
                    cp -f "$example_file" "$TARGET_DIR/planning/$rel_path"
                    echo -e "  Updated: planning/$rel_path"
                done
            else
                # No user content, safe to copy everything
                echo -e "  Installing complete planning directory structure..."
                mkdir -p "$TARGET_DIR/planning"
                rsync -av "$TEMPLATE_DIR/planning/" "$TARGET_DIR/planning/"
                echo -e "  ${GREEN}✓ Planning structure with examples installed${NC}"
            fi
        fi
        echo -e "${GREEN}✓ Planning directory updated${NC}"
        echo -e "${YELLOW}  📁 Project plans: planning/project/${NC}"
        echo -e "${YELLOW}  📁 Phase plans: planning/phase*/${NC}"
        echo -e "${YELLOW}  📁 Wave plans: planning/phase*/wave*/${NC}"
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
            
            # Copy the entire agent-states structure including all sub-state machines
            echo -e "  Syncing complete agent-states directory structure..."

            # Copy ALL state machine directories (software-factory, initialization, etc.)
            # CRITICAL: Exclude ARCHIVED directories - they should NEVER exist in active projects
            if [ "$TEMPLATE_DIR" != "$TARGET_DIR" ]; then
                # Use rsync to intelligently sync the entire structure
                ionice -c 2 -n 7 nice -n 10 rsync -av --delete \
                    --exclude='ARCHIVED' \
                    --exclude='*/ARCHIVED' \
                    --exclude='*/*/ARCHIVED' \
                    "$TEMPLATE_DIR/agent-states/" \
                    "$TARGET_DIR/agent-states/"
                echo -e "    ${GREEN}✓ Synced all state machines and agent states (ARCHIVED excluded)${NC}"
            else
                echo -e "    ${CYAN}Skipping (source and destination are the same)${NC}"
            fi

            # List what was updated for clarity
            echo -e "  Updated state machines:"
            for state_dir in "$TARGET_DIR/agent-states"/*; do
                if [ -d "$state_dir" ]; then
                    dirname=$(basename "$state_dir")
                    count=$(find "$state_dir" -name "rules.md" 2>/dev/null | wc -l)
                    if [ $count -gt 0 ]; then
                        echo -e "    • $dirname: $count states"
                    fi
                fi
            done
            
            # Count total state directories updated
            state_count=$(find "$TARGET_DIR/agent-states" -name "rules.md" | wc -l)
            echo -e "  Updated $state_count state rule files"
        fi
        echo -e "${GREEN}✓ Agent state rules updated${NC}"
    fi

    # Software Factory 3.0 Specific Files
    if [ "$SF_VERSION" = "3.0" ]; then
        echo -e "${CYAN}Updating Software Factory 3.0 specific files...${NC}"

        if [ "$DRY_RUN" = true ]; then
            echo -e "${CYAN}[DRY RUN] Would update SF 3.0 files${NC}"
        else
            # SF 3.0 example state files
            sf3_state_examples=(
                "orchestrator-state-v3.json.example"
                "bug-tracking.json.example"
                "integration-containers.json.example"
                "fix-cascade-state.json.example"
            )

            for example_file in "${sf3_state_examples[@]}"; do
                if [ -f "$TEMPLATE_DIR/$example_file" ]; then
                    echo -e "  Updating: $example_file"
                    cp -f "$TEMPLATE_DIR/$example_file" "$TARGET_DIR/$example_file"
                fi
            done

            # SF 3.0 validation and state management tools
            sf3_tools=(
                "tools/validate-state-file.sh"
                "tools/atomic-state-update.sh"
            )

            for tool_file in "${sf3_tools[@]}"; do
                if [ -f "$TEMPLATE_DIR/$tool_file" ]; then
                    echo -e "  Updating: $tool_file"
                    mkdir -p "$(dirname "$TARGET_DIR/$tool_file")"
                    cp -f "$TEMPLATE_DIR/$tool_file" "$TARGET_DIR/$tool_file"
                    chmod +x "$TARGET_DIR/$tool_file"
                fi
            done

            # SF 3.0 templates directory
            if [ -d "$TEMPLATE_DIR/templates" ]; then
                echo -e "  Syncing templates directory..."
                mkdir -p "$TARGET_DIR/templates"
                ionice -c 2 -n 7 nice -n 10 rsync -av --update "$TEMPLATE_DIR/templates/" "$TARGET_DIR/templates/" | sed 's/^/    /'
            fi

            # SF 3.0 state machines
            if [ -f "$TEMPLATE_DIR/state-machines/software-factory-3.0-state-machine.json" ]; then
                echo -e "  Updating: state-machines/software-factory-3.0-state-machine.json"
                mkdir -p "$TARGET_DIR/state-machines"
                cp -f "$TEMPLATE_DIR/state-machines/software-factory-3.0-state-machine.json" \
                      "$TARGET_DIR/state-machines/software-factory-3.0-state-machine.json"
            fi

            # SF 3.0 schemas directory
            if [ -d "$TEMPLATE_DIR/schemas" ]; then
                echo -e "  Syncing schemas directory..."
                mkdir -p "$TARGET_DIR/schemas"

                if [ "$TEMPLATE_DIR" != "$TARGET_DIR" ]; then
                    # Copy all schema files from schemas/ directory
                    for schema_file in "$TEMPLATE_DIR/schemas"/*.json; do
                        if [ -f "$schema_file" ]; then
                            basename=$(basename "$schema_file")
                            dest_file="$TARGET_DIR/schemas/$basename"

                            # Check if source and destination are the same file
                            if [ "$schema_file" = "$dest_file" ]; then
                                echo -e "    ${CYAN}Skipping:${NC} schemas/$basename (source and destination are the same)"
                            else
                                echo -e "    Updating: schemas/$basename"
                                cp -f "$schema_file" "$dest_file"
                            fi
                        fi
                    done
                else
                    echo -e "    ${CYAN}Skipping (source and destination are the same)${NC}"
                fi
            fi

            # SF 3.0 schema files (legacy root location - for backward compatibility)
            # Note: Schemas are now primarily in schemas/ directory above
            sf3_schemas=(
                "orchestrator-state-v3.schema.json"
                "bug-tracking.schema.json"
                "integration-containers.schema.json"
                "fix-cascade-state.schema.json"
            )

            for schema_file in "${sf3_schemas[@]}"; do
                if [ -f "$TEMPLATE_DIR/$schema_file" ]; then
                    dest_file="$TARGET_DIR/$schema_file"

                    # Check if source and destination are the same file
                    if [ "$TEMPLATE_DIR/$schema_file" = "$dest_file" ]; then
                        echo -e "  ${CYAN}Skipping:${NC} $schema_file (source and destination are the same)"
                    else
                        echo -e "  Updating: $schema_file"
                        cp -f "$TEMPLATE_DIR/$schema_file" "$dest_file"
                    fi
                fi
            done

            # SF 3.0 utilities
            sf3_utilities=(
                "utilities/init-software-factory.sh"
            )

            for util_file in "${sf3_utilities[@]}"; do
                if [ -f "$TEMPLATE_DIR/$util_file" ]; then
                    echo -e "  Updating: $util_file"
                    mkdir -p "$(dirname "$TARGET_DIR/$util_file")"
                    cp -f "$TEMPLATE_DIR/$util_file" "$TARGET_DIR/$util_file"
                    chmod +x "$TARGET_DIR/$util_file"
                fi
            done

            echo -e "${GREEN}✓ SF 3.0 specific files updated${NC}"
        fi
    fi

    # Update critical files
    echo -e "${CYAN}Updating critical files...${NC}"
    critical_files=(
        "rule-library/R002-agent-acknowledgment.md"
        "rule-library/R003-performance-grading.md"
        "DELEGATION-AND-AGENT-PROMPT-GUIDE.md"
        "CRITICAL-IMPLEMENTATION-GUIDANCE.md"
        "AGENT-ACKNOWLEDGMENT-EXAMPLES.md"
        "software-factory-3.0-state-machine.json"
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

    # CRITICAL: Remove .template-repository if it somehow exists in project instance
    if [ -f "$TARGET_DIR/.template-repository" ]; then
        if [ "$DRY_RUN" = true ]; then
            echo -e "${CYAN}[DRY RUN] Would remove .template-repository from project instance${NC}"
        else
            echo -e "${YELLOW}⚠ Removing .template-repository from project instance${NC}"
            echo -e "${YELLOW}  This file should only exist in the template repository${NC}"
            rm -f "$TARGET_DIR/.template-repository"
            echo -e "${GREEN}✓ .template-repository removed${NC}"
        fi
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
            echo "version: $SF_VERSION"
            echo "upgraded_at: $(date +%Y%m%d-%H%M%S)"
            echo "upgraded_from: $TEMPLATE_DIR"
            echo "template_commit: $TEMPLATE_COMMIT"
            echo "template_branch: $TEMPLATE_BRANCH"
            echo "template_remote: $TEMPLATE_REMOTE"
            echo "upgrade_date: $(date -u +"%Y-%m-%dT%H:%M:%SZ")"
        } > "$TARGET_DIR/.sf-version"
        
        echo -e "${GREEN}✓ Version marker created with commit: ${CYAN}$TEMPLATE_COMMIT${NC}"
    fi
    
    # Install/Update pre-commit hooks using the unified hook system
    echo -e "\n${CYAN}Installing/Updating pre-commit hooks...${NC}"
    if [ "$DRY_RUN" != true ]; then
        # Ensure tools/git-commit-hooks directory and master hook are up to date
        echo -e "  Syncing git-commit-hooks directory..."
        mkdir -p "$TARGET_DIR/tools/git-commit-hooks"

        # Always use rsync for reliability - it's a standard tool on all systems
        ionice -c 2 -n 7 nice -n 10 rsync -av --update "$TEMPLATE_DIR/tools/git-commit-hooks/" "$TARGET_DIR/tools/git-commit-hooks/" | sed 's/^/    /'

        # Make master-pre-commit.sh executable
        if [ -f "$TARGET_DIR/tools/git-commit-hooks/master-pre-commit.sh" ]; then
            chmod +x "$TARGET_DIR/tools/git-commit-hooks/master-pre-commit.sh"
            echo -e "${GREEN}✓ Unified hook system updated${NC}"
        fi

        # Check if the new installation script exists
        if [ -f "$TARGET_DIR/utilities/install-pre-commit-hooks.sh" ]; then
            echo -e "${BLUE}ℹ Using automated pre-commit hook installation script${NC}"

            # Change to target directory for proper execution
            local original_dir="$(pwd)"
            cd "$TARGET_DIR"

            # Run the installation script
            if bash utilities/install-pre-commit-hooks.sh; then
                echo -e "${GREEN}✓ Pre-commit hooks installed successfully${NC}"
                hook_installed=true
            else
                echo -e "${YELLOW}⚠ Pre-commit hook installation reported warnings${NC}"
                hook_installed=true  # Still mark as installed even with warnings
            fi

            cd "$original_dir"
        else
            # Fallback to direct installation
            echo -e "${YELLOW}⚠ Automated hook installer not found, installing directly${NC}"

            if [ -d "$TARGET_DIR/.git" ]; then
                mkdir -p "$TARGET_DIR/.git/hooks"

                # Install unified master-pre-commit.sh hook
                local hook_installed=false

                if [ -f "$TARGET_DIR/tools/git-commit-hooks/master-pre-commit.sh" ]; then
                    cp -f "$TARGET_DIR/tools/git-commit-hooks/master-pre-commit.sh" \
                          "$TARGET_DIR/.git/hooks/pre-commit"
                    chmod +x "$TARGET_DIR/.git/hooks/pre-commit"
                    hook_installed=true
                    echo -e "${GREEN}✓ Pre-commit hook installed (unified master-pre-commit.sh)${NC}"
                    echo -e "  ${CYAN}Hook will auto-detect SF 2.0/3.0 and repo type${NC}"
                else
                    echo -e "${RED}✗ Could not find master-pre-commit.sh${NC}"
                    echo -e "${YELLOW}  Please ensure tools/git-commit-hooks/ is properly synced${NC}"
                fi
            fi
        fi

        # Always sync the hooks directory structure (ensure new hooks are copied)
        if [ -d "$TEMPLATE_DIR/tools/git-commit-hooks" ]; then
            echo -e "${CYAN}Syncing git commit hooks library...${NC}"
            mkdir -p "$TARGET_DIR/tools/git-commit-hooks"

            # Copy hooks with rsync to update existing and add new ones
            if [ "$TEMPLATE_DIR" != "$TARGET_DIR" ]; then
                # Always use rsync - it's standard on all modern systems
                ionice -c 2 -n 7 nice -n 10 rsync -av --update "$TEMPLATE_DIR/tools/git-commit-hooks/" "$TARGET_DIR/tools/git-commit-hooks/" | sed 's/^/  /'
                echo -e "${GREEN}✓ Hook library synced to tools/git-commit-hooks/${NC}"
            else
                echo -e "${CYAN}  Skipping (source and destination are the same)${NC}"
            fi
        fi

        # If no hook was installed yet, create inline fallback
        if [ "$hook_installed" != true ]; then
            # Fallback: Create the hook inline if no template found
            echo -e "${YELLOW}⚠ No template hook found, creating comprehensive hook inline...${NC}"
            cat > "$TARGET_DIR/.git/hooks/pre-commit" << 'HOOK_EOF'
#!/bin/bash

# Pre-commit hook for Software Factory 3.0
# This hook validates orchestrator-state-v3.json before allowing commits
# It ensures the state file is always valid and follows the proper schema

set -euo pipefail

# Color codes for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
BOLD='\033[1m'
NC='\033[0m' # No Color

# Get the git root directory
GIT_ROOT="$(git rev-parse --show-toplevel)"
cd "$GIT_ROOT"

# Function to print colored messages
print_error() {
    echo -e "${RED}${BOLD}❌ ERROR:${NC} ${RED}$1${NC}" >&2
}

print_success() {
    echo -e "${GREEN}${BOLD}✅ PROJECT_DONE:${NC} ${GREEN}$1${NC}"
}

print_info() {
    echo -e "${BLUE}${BOLD}ℹ️  INFO:${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}${BOLD}⚠️  WARNING:${NC} ${YELLOW}$1${NC}"
}

# Function to check if file is staged for commit
is_file_staged() {
    local file="$1"
    git diff --cached --name-only | grep -q "^${file}$"
}

# Function to ensure jsonschema is installed
ensure_jsonschema_installed() {
    # Check if Python3 is available
    if ! command -v python3 &> /dev/null; then
        print_warning "Python3 not found, skipping jsonschema installation"
        return 1
    fi

    # Check if jsonschema is already installed and accessible in clean environment
    if python3 -c "import jsonschema" &> /dev/null; then
        # Test in clean environment (like agent-spawned shells)
        if env -i HOME="$HOME" PATH="/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin" \
           bash -c "python3 -c 'import jsonschema' 2>/dev/null"; then
            # Already working globally, nothing to do
            return 0
        else
            print_warning "jsonschema found but not accessible in clean environments"
        fi
    fi

    print_info "Ensuring jsonschema is globally accessible..."

    # Check for Software Factory venv
    local SF_VENV_DIR="${HOME}/.software-factory-venv"

    # If venv exists and has jsonschema, ensure it's activated in profiles
    if [ -f "$SF_VENV_DIR/bin/python" ] && \
       "$SF_VENV_DIR/bin/python" -c "import jsonschema" 2>/dev/null; then
        print_info "Software Factory venv found with jsonschema"

        # Ensure activation in profiles
        local activation_snippet='
# Software Factory 3.0 - Auto-activate virtual environment for jsonschema
if [ -f "$HOME/.software-factory-venv/bin/activate" ] && [ -z "$VIRTUAL_ENV" ]; then
    source "$HOME/.software-factory-venv/bin/activate" 2>/dev/null
    export SF_VENV_ACTIVE=1
fi'

        # Add to .bashrc if not already there
        if [ -f "$HOME/.bashrc" ]; then
            if ! grep -q "Software Factory 3.0 - Auto-activate virtual environment" "$HOME/.bashrc" 2>/dev/null; then
                echo "$activation_snippet" >> "$HOME/.bashrc"
                print_success "Added venv activation to ~/.bashrc"
            fi
        fi

        # Add to .profile for sh compatibility
        if [ -f "$HOME/.profile" ] || [ ! -f "$HOME/.bashrc" ]; then
            if ! grep -q "Software Factory 3.0 - Auto-activate virtual environment" "$HOME/.profile" 2>/dev/null; then
                echo "$activation_snippet" >> "$HOME/.profile"
                print_success "Added venv activation to ~/.profile"
            fi
        fi

        return 0
    fi

    # Try standard user installation first
    if command -v pip3 &> /dev/null; then
        if pip3 install --user jsonschema &> /dev/null; then
            print_success "jsonschema module installed successfully"
            return 0
        else
            # Try with --break-system-packages flag for newer systems
            if pip3 install --user --break-system-packages jsonschema &> /dev/null 2>&1; then
                print_success "jsonschema module installed successfully"
                return 0
            fi
        fi
    fi

    # If standard install failed, try creating venv automatically (non-interactive for upgrade)
    print_info "Creating Software Factory virtual environment for reliable jsonschema access..."

    # Remove old venv if exists
    [ -d "$SF_VENV_DIR" ] && rm -rf "$SF_VENV_DIR"

    # Create new venv
    if python3 -m venv "$SF_VENV_DIR"; then
        # Install jsonschema in venv
        if "$SF_VENV_DIR/bin/pip" install --upgrade pip jsonschema &> /dev/null; then
            print_success "jsonschema installed in Software Factory venv"

            # Add activation to shell profiles
            local activation_snippet='
# Software Factory 3.0 - Auto-activate virtual environment for jsonschema
if [ -f "$HOME/.software-factory-venv/bin/activate" ] && [ -z "$VIRTUAL_ENV" ]; then
    source "$HOME/.software-factory-venv/bin/activate" 2>/dev/null
    export SF_VENV_ACTIVE=1
fi'

            # Add to .bashrc if not already there
            if [ -f "$HOME/.bashrc" ]; then
                if ! grep -q "Software Factory 3.0 - Auto-activate virtual environment" "$HOME/.bashrc" 2>/dev/null; then
                    echo "$activation_snippet" >> "$HOME/.bashrc"
                    print_success "Added venv activation to ~/.bashrc"
                fi
            fi

            # Add to .profile for sh compatibility
            if [ -f "$HOME/.profile" ] || [ ! -f "$HOME/.bashrc" ]; then
                if ! grep -q "Software Factory 3.0 - Auto-activate virtual environment" "$HOME/.profile" 2>/dev/null; then
                    echo "$activation_snippet" >> "$HOME/.profile"
                    print_success "Added venv activation to ~/.profile"
                fi
            fi

            print_success "Software Factory venv configured successfully"
            print_info "Note: Run 'source ~/.bashrc' to activate in current shell"
            return 0
        fi
    fi

    print_warning "Could not ensure jsonschema availability"
    print_info "You can manually run: bash tools/ensure-jsonschema.sh"
    return 1
}

# Main validation logic
main() {
    local validation_failed=false

    # Ensure jsonschema is installed for better validation
    ensure_jsonschema_installed || true

    # Check if orchestrator-state-v3.json is being committed
    if is_file_staged "orchestrator-state-v3.json"; then
        echo ""
        print_info "Detected changes to orchestrator-state-v3.json"
        echo -e "${BOLD}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"

        # Check if validation script exists
        if [ ! -f "tools/validate-state.sh" ]; then
            print_error "Validation script tools/validate-state.sh not found!"
            print_info "Please ensure the Software Factory validation tools are properly installed."
            exit 1
        fi

        # Run the validation
        print_info "Running orchestrator-state-v3.json validation..."
        echo ""

        if tools/validate-state.sh orchestrator-state-v3.json; then
            echo ""
            print_success "orchestrator-state-v3.json validation passed!"
            echo -e "${BOLD}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
        else
            echo ""
            print_error "orchestrator-state-v3.json validation failed!"
            echo -e "${BOLD}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
            echo ""
            print_error "The orchestrator-state-v3.json file has validation errors."
            print_info "Please fix the validation errors before committing."
            print_info "You can run 'tools/validate-state.sh orchestrator-state-v3.json' to see detailed errors."
            echo ""
            validation_failed=true
        fi
    fi

    # Also validate state machine files if they're being committed
    local state_machine_files=(
        "state-machines/software-factory-3.0-state-machine.json"
        "state-machines/initialization-state-machine.json"
        "state-machines/pr-ready-state-machine.json"
        "state-machines/fix-cascade-state-machine.json"
        "state-machines/integration-state-machine.json"
        "state-machines/splitting-state-machine.json"
    )

    local any_state_machine_staged=false
    for file in "${state_machine_files[@]}"; do
        if is_file_staged "$file"; then
            any_state_machine_staged=true
            break
        fi
    done

    if [ "$any_state_machine_staged" = true ]; then
        echo ""
        print_info "Detected changes to state machine files"
        echo -e "${BOLD}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"

        if [ -f "tools/validate-state-machines.sh" ]; then
            print_info "Running state machine validation..."
            echo ""

            if tools/validate-state-machines.sh; then
                echo ""
                print_success "State machine validation passed!"
                echo -e "${BOLD}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
            else
                echo ""
                print_error "State machine validation failed!"
                echo -e "${BOLD}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
                echo ""
                print_error "One or more state machine files have validation errors."
                print_info "Please fix the validation errors before committing."
                echo ""
                validation_failed=true
            fi
        else
            print_warning "State machine validation script not found, skipping validation"
        fi
    fi

    # Check for R251 compliance (no commits on main branch in effort repos)
    current_branch="$(git branch --show-current)"
    remote_url="$(git config --get remote.origin.url 2>/dev/null || true)"

    # Check if this is an effort repository (not the software-factory-template)
    if [[ "$remote_url" != *"software-factory-template"* ]] && [[ "$current_branch" == "main" ]]; then
        if [[ -d "../software-factory-template" ]] || [[ "$GIT_ROOT" =~ efforts/ ]]; then
            echo ""
            print_error "R251 VIOLATION: Cannot commit to main branch in effort repository!"
            echo -e "${BOLD}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
            echo ""
            print_info "You are attempting to commit to the main branch in an effort repository."
            print_info "Per R251 (Repository Separation), all work must be done in feature branches."
            echo ""
            print_info "Please create and switch to a feature branch:"
            echo -e "  ${BOLD}git checkout -b your-feature-branch${NC}"
            echo ""
            validation_failed=true
        fi
    fi

    # Exit with appropriate code
    if [ "$validation_failed" = true ]; then
        exit 1
    fi

    exit 0
}

# Run main function
main "$@"
HOOK_EOF
            chmod +x "$TARGET_DIR/.git/hooks/pre-commit"
            hook_installed=true
            echo -e "${GREEN}✓ Comprehensive pre-commit hook created inline${NC}"
        fi

        # Display hook features if successfully installed
        if [ "$hook_installed" = true ]; then
            echo -e "${YELLOW}  Hook features:${NC}"
            echo -e "${YELLOW}    • Validates orchestrator-state-v3.json with schema${NC}"
            echo -e "${YELLOW}    • Validates state machine JSON files${NC}"
            echo -e "${YELLOW}    • Auto-configures PYTHONPATH for jsonschema${NC}"
            echo -e "${YELLOW}    • Auto-installs Python jsonschema if needed${NC}"
            echo -e "${YELLOW}    • Enforces R251 - Repository Separation${NC}"
        fi

        # Ensure Python jsonschema is properly installed
        echo -e "${CYAN}Ensuring Python jsonschema module is installed...${NC}"
            if command -v python3 &> /dev/null; then
                # Always try to ensure it's installed, even if it appears to be present
                # This handles cases where it might be in a different Python path

                # First check current status
                if python3 -c "import jsonschema" &> /dev/null 2>&1; then
                    echo -e "${GREEN}  ✓ jsonschema module found${NC}"
                    # But still ensure it's in the user site-packages
                    if command -v pip3 &> /dev/null; then
                        # Force reinstall to user directory to ensure it's accessible
                        echo -e "${CYAN}  Ensuring jsonschema is in user site-packages...${NC}"
                        if pip3 install --user --force-reinstall jsonschema &> /dev/null 2>&1; then
                            echo -e "${GREEN}  ✓ jsonschema reinstalled to user directory${NC}"
                        elif pip3 install --user --force-reinstall --break-system-packages jsonschema &> /dev/null 2>&1; then
                            echo -e "${GREEN}  ✓ jsonschema reinstalled to user directory${NC}"
                        else
                            # If force reinstall fails, at least try upgrade
                            pip3 install --user --upgrade jsonschema &> /dev/null 2>&1 || \
                            pip3 install --user --upgrade --break-system-packages jsonschema &> /dev/null 2>&1 || \
                            echo -e "${YELLOW}  ⚠ Could not reinstall, but module is available${NC}"
                        fi
                    fi
                else
                    # Module not found, need to install
                    echo -e "${YELLOW}  jsonschema module not found, installing...${NC}"
                    if command -v pip3 &> /dev/null; then
                        if pip3 install --user jsonschema &> /dev/null 2>&1; then
                            echo -e "${GREEN}  ✓ jsonschema module installed successfully${NC}"
                        elif pip3 install --user --break-system-packages jsonschema &> /dev/null 2>&1; then
                            echo -e "${GREEN}  ✓ jsonschema module installed successfully${NC}"
                        else
                            echo -e "${RED}  ❌ Failed to install jsonschema${NC}"
                            echo -e "${YELLOW}  The git hook will attempt to install it on first use${NC}"
                        fi
                    else
                        echo -e "${RED}  ❌ pip3 not found, cannot install jsonschema${NC}"
                        echo -e "${YELLOW}  Please install manually: pip3 install jsonschema${NC}"
                    fi
                fi

                # Final verification
                if python3 -c "import jsonschema; import importlib.metadata; print(f'    Version: {importlib.metadata.version(\"jsonschema\")}')" 2>/dev/null; then
                    echo -e "${GREEN}  ✓ jsonschema is ready for use${NC}"
                else
                    echo -e "${YELLOW}  ⚠ jsonschema status unclear - hook will handle installation if needed${NC}"
                fi
            else
                echo -e "${YELLOW}  ⚠ Python3 not found - jsonschema will be installed when needed${NC}"
            fi
    else
        echo -e "${YELLOW}[DRY RUN] Would install comprehensive pre-commit hook${NC}"
    fi
}

# Install base branch validation hooks in all effort branches
install_effort_repo_hooks() {
    echo -e "\n${CYAN}Installing effort repository hooks...${NC}"

    # Use the new install-effort-hooks.sh script
    local hook_installer="$TEMPLATE_DIR/tools/install-effort-hooks.sh"

    if [ ! -f "$hook_installer" ]; then
        echo -e "${YELLOW}⚠ Hook installer script not found in template${NC}"
        echo -e "${YELLOW}  Expected: $hook_installer${NC}"
        echo -e "${YELLOW}  Effort repositories will not have pre-commit validation!${NC}"
        return
    fi

    if [ "$DRY_RUN" != true ]; then
        local hooks_installed=0
        local hooks_failed=0

        # Skip main project repo - it should have planning hooks, not effort hooks
        echo -e "${BLUE}ℹ Main project repository uses planning hooks (installed separately)${NC}"

        # Find and install in all effort working copies
        if [ -d "$TARGET_DIR/efforts" ]; then
            echo -e "${BLUE}Scanning for effort working copies...${NC}"

            while IFS= read -r -d '' git_dir; do
                local repo_dir="${git_dir%/.git}"

                # Check if this is an effort repository (phase/wave pattern)
                if [[ "$repo_dir" =~ efforts/phase[0-9]+/wave[0-9]+ ]]; then
                    echo -e "${BLUE}Installing hooks in: ${repo_dir#$TARGET_DIR/}${NC}"

                    # Run the install-effort-hooks.sh script
                    if bash "$hook_installer" "$repo_dir" 2>&1 | sed 's/^/  /'; then
                        ((hooks_installed++))
                        echo -e "${GREEN}✓ Hooks installed successfully in: ${repo_dir#$TARGET_DIR/}${NC}"
                    else
                        ((hooks_failed++))
                        echo -e "${RED}✗ Hook installation failed in: ${repo_dir#$TARGET_DIR/}${NC}"
                    fi
                    echo ""
                fi
            done < <(find "$TARGET_DIR/efforts" -type d -name ".git" -print0 2>/dev/null)
        else
            echo -e "${BLUE}ℹ No efforts/ directory found - no effort repositories to update${NC}"
        fi

        echo ""
        echo -e "${BOLD}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
        echo -e "${CYAN}Effort Repository Hook Installation Summary:${NC}"
        echo -e "  ${GREEN}✓ Successfully installed: $hooks_installed${NC}"
        if [ $hooks_failed -gt 0 ]; then
            echo -e "  ${RED}✗ Failed installations: $hooks_failed${NC}"
        fi
        echo -e "${BOLD}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
        echo ""

        if [ $hooks_installed -gt 0 ]; then
            echo -e "${GREEN}Effort repositories now have:${NC}"
            echo -e "  ${GREEN}• R383 - Metadata placement validation (SUPREME LAW)${NC}"
            echo -e "  ${GREEN}• R343 - Metadata directory standardization${NC}"
            echo -e "  ${GREEN}• R506 - Absolute prohibition on --no-verify${NC}"
            echo -e "  ${GREEN}• R500 - Branch HEAD tracking synchronization${NC}"
            echo ""
        fi
    else
        echo -e "${YELLOW}[DRY RUN] Would install effort-specific hooks in all effort repositories using install-effort-hooks.sh${NC}"
    fi
}

# Confirmation prompt
if [ "$FORCE" != true ] && [ "$DRY_RUN" != true ]; then
    echo -e "\n${YELLOW}⚠️  This will update rules and configurations in:${NC}"
    echo -e "${YELLOW}   $TARGET_DIR${NC}"
    echo -e "\n${GREEN}The following will be preserved:${NC}"
    echo -e "  • efforts/ directory (all work in progress)"
    echo -e "  • todos/ directory"
    echo -e "  • orchestrator-state-v3.json"
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

# Install base branch validation hooks
install_effort_repo_hooks

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
    echo -e "  • Git pre-commit hook with comprehensive validation:"
    echo -e "      - Validates orchestrator-state-v3.json with schema"
    echo -e "      - Validates state machine JSON files"
    echo -e "      - Auto-installs Python jsonschema module"
    echo -e "      - Enforces R251 - Repository Separation"
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
        echo -e "${CYAN}Backup information:${NC}"
        if [ "$COMPRESS_BACKUP" = true ]; then
            echo -e "  💾 Backup location: backups/backup-*.tar.gz (compressed)"
            echo -e "  📁 Backup contents: All project files (excluding backups/, efforts/, .git/)"
        else
            echo -e "  💾 Backup location: backups/backup-*/ (uncompressed)"
            echo -e "  📁 Backup contents: All project files (excluding backups/, efforts/, .git/)"
        fi
        echo -e "  🔄 Effort backups: backups/efforts-upgrade-*/ (with full git history)"
        echo -e ""
        echo -e "${YELLOW}⚠️  IMPORTANT: Effort backups contain complete git history${NC}"
        echo -e "${YELLOW}   Keep these backups until you're certain all branches are safe${NC}"
        echo ""
        echo -e "${CYAN}Manage backups:${NC}"
        echo -e "  • View all backups: ls -lh backups/"
        echo -e "  • Cleanup old backups: bash tools/cleanup-backups.sh"
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