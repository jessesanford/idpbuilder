#!/bin/bash

# Software Factory 2.0 Setup Script
# This script helps you quickly set up a new Software Factory 2.0 project
#
# Features:
# - Creates project structure and configuration files
# - Installs Python dependencies (jsonschema) for validation
# - Sets up pre-commit hooks for state validation
# - Copies all necessary Software Factory 2.0 files
#
# Usage:
#   Interactive: ./setup.sh
#   Non-interactive: ./setup.sh --config setup-config.yaml

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

# Check for --config flag
if [[ "$1" == "--config" ]] || [[ "$1" == "-c" ]]; then
    # Redirect to non-interactive setup script
    SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
    if [ -f "$SCRIPT_DIR/setup-noninteractive.sh" ]; then
        echo -e "${CYAN}Redirecting to non-interactive setup...${NC}"
        exec "$SCRIPT_DIR/setup-noninteractive.sh" "$@"
    else
        echo -e "${RED}Error: setup-noninteractive.sh not found${NC}"
        echo -e "${YELLOW}Please ensure setup-noninteractive.sh is in the same directory${NC}"
        exit 1
    fi
fi

# Check for --help flag
if [[ "$1" == "--help" ]] || [[ "$1" == "-h" ]]; then
    echo "Software Factory 2.0 Setup Script"
    echo ""
    echo "Usage:"
    echo "  $0                         # Interactive setup"
    echo "  $0 --config <file>         # Non-interactive setup with config file"
    echo "  $0 --help                  # Show this help message"
    echo ""
    echo "For non-interactive setup, create a configuration file based on"
    echo "setup-config-example.yaml and run with --config flag."
    exit 0
fi

# Banner
echo -e "${CYAN}${BOLD}"
cat << "EOF"
╔═══════════════════════════════════════════════════════════════════╗
║                                                                   ║
║   ███████╗ ██████╗ ███████╗████████╗██╗    ██╗ █████╗ ██████╗   ║
║   ██╔════╝██╔═══██╗██╔════╝╚══██╔══╝██║    ██║██╔══██╗██╔══██╗  ║
║   ███████╗██║   ██║█████╗     ██║   ██║ █╗ ██║███████║██████╔╝  ║
║   ╚════██║██║   ██║██╔══╝     ██║   ██║███╗██║██╔══██║██╔══██╗  ║
║   ███████║╚██████╔╝██║        ██║   ╚███╔███╔╝██║  ██║██║  ██║  ║
║   ╚══════╝ ╚═════╝ ╚═╝        ╚═╝    ╚══╝╚══╝ ╚═╝  ╚═╝╚═╝  ╚═╝  ║
║                                                                   ║
║            ███████╗ █████╗  ██████╗████████╗ ██████╗ ██████╗     ║
║            ██╔════╝██╔══██╗██╔════╝╚══██╔══╝██╔═══██╗██╔══██╗    ║
║            █████╗  ███████║██║        ██║   ██║   ██║██████╔╝    ║
║            ██╔══╝  ██╔══██║██║        ██║   ██║   ██║██╔══██╗    ║
║            ██║     ██║  ██║╚██████╗   ██║   ╚██████╔╝██║  ██║    ║
║            ╚═╝     ╚═╝  ╚═╝ ╚═════╝   ╚═╝    ╚═════╝ ╚═╝  ╚═╝    ║
║                                                                   ║
║                            Version 2.0                           ║
║                                                                   ║
╚═══════════════════════════════════════════════════════════════════╝
EOF
echo -e "${NC}"

# Function to prompt with default
prompt_with_default() {
    local prompt="$1"
    local default="$2"
    local var_name="$3"
    
    echo -ne "${CYAN}${prompt}${NC}"
    if [ -n "$default" ]; then
        echo -ne " ${YELLOW}[$default]${NC}: "
    else
        echo -ne ": "
    fi
    
    read -r response
    if [ -z "$response" ] && [ -n "$default" ]; then
        # Use printf to safely assign the value
        printf -v "$var_name" "%s" "$default"
    else
        # Use printf to safely assign the value, escaping single quotes
        printf -v "$var_name" "%s" "$response"
    fi
}

# Function to select from options
select_option() {
    local prompt="$1"
    shift
    local options=("$@")
    
    echo -e "${CYAN}${prompt}${NC}"
    for i in "${!options[@]}"; do
        echo -e "  ${BOLD}$((i+1)))${NC} ${options[$i]}"
    done
    
    local choice
    while true; do
        echo -ne "${CYAN}Select option (1-${#options[@]}): ${NC}"
        read -r choice
        if [[ "$choice" =~ ^[0-9]+$ ]] && [ "$choice" -ge 1 ] && [ "$choice" -le "${#options[@]}" ]; then
            # Set global variable instead of using return
            SELECTED_INDEX=$((choice-1))
            return 0
        fi
        echo -e "${RED}Invalid selection. Please try again.${NC}"
    done
}

# Function to select multiple options
select_multiple() {
    local prompt="$1"
    shift
    local options=("$@")
    local selected=()
    
    echo -e "${CYAN}${prompt}${NC}"
    echo -e "${YELLOW}(Enter numbers separated by spaces, or 'all' for all options)${NC}"
    
    for i in "${!options[@]}"; do
        echo -e "  ${BOLD}$((i+1)))${NC} ${options[$i]}"
    done
    
    echo -ne "${CYAN}Select options: ${NC}"
    read -r choices
    
    if [ "$choices" = "all" ]; then
        selected=("${options[@]}")
    else
        for choice in $choices; do
            if [[ "$choice" =~ ^[0-9]+$ ]] && [ "$choice" -ge 1 ] && [ "$choice" -le "${#options[@]}" ]; then
                selected+=("${options[$((choice-1))]}")
            fi
        done
    fi
    
    echo "${selected[@]}"
}

# Step 1: Project Information
echo -e "\n${MAGENTA}${BOLD}═══ Step 1: Project Information ═══${NC}\n"

prompt_with_default "Project name" "" "PROJECT_NAME"
while [ -z "$PROJECT_NAME" ]; do
    echo -e "${RED}Project name cannot be empty!${NC}"
    prompt_with_default "Project name" "" "PROJECT_NAME"
done

prompt_with_default "Project description" "" "PROJECT_DESC"
prompt_with_default "Target directory" "/workspaces/$PROJECT_NAME" "TARGET_DIR"

# Validate target directory path
if [[ ! "$TARGET_DIR" =~ ^/ ]]; then
    echo -e "${YELLOW}⚠ Converting relative path to absolute path${NC}"
    TARGET_DIR="$(pwd)/$TARGET_DIR"
    echo -e "${CYAN}Using: $TARGET_DIR${NC}"
fi

# Check if parent directory exists and is writable
TARGET_PARENT=$(dirname "$TARGET_DIR")
if [ ! -d "$TARGET_PARENT" ]; then
    echo -e "${YELLOW}⚠ Parent directory does not exist: $TARGET_PARENT${NC}"
    echo -ne "${CYAN}Create it? (y/n): ${NC}"
    read -r create_parent
    if [ "$create_parent" = "y" ]; then
        mkdir -p "$TARGET_PARENT" 2>/dev/null || {
            echo -e "${RED}❌ Cannot create parent directory. Permission denied.${NC}"
            echo -e "${YELLOW}Try using a different location or run with appropriate permissions.${NC}"
            exit 1
        }
        echo -e "${GREEN}✓ Parent directory created${NC}"
    else
        echo -e "${RED}Setup cancelled.${NC}"
        exit 1
    fi
elif [ ! -w "$TARGET_PARENT" ]; then
    echo -e "${RED}❌ No write permission in parent directory: $TARGET_PARENT${NC}"
    echo -e "${YELLOW}Try using a different location or run with appropriate permissions.${NC}"
    exit 1
fi

prompt_with_default "GitHub repository URL (optional)" "" "GITHUB_URL"

# Step 1.5: Target Repository Configuration
echo -e "\n${MAGENTA}${BOLD}═══ Step 1.5: Target Repository ═══${NC}\n"
echo -e "${YELLOW}The target repository is the actual project code you'll be working on.${NC}"
echo -e "${YELLOW}This is different from the Software Factory instance repository.${NC}\n"

prompt_with_default "Target repository URL (e.g., https://github.com/owner/repo.git)" "" "TARGET_REPO_URL"
while [ -z "$TARGET_REPO_URL" ]; do
    echo -e "${RED}Target repository URL is required!${NC}"
    echo -e "${YELLOW}This is the repository where your code will be developed.${NC}"
    prompt_with_default "Target repository URL" "" "TARGET_REPO_URL"
done

# Validate URL format
if [[ ! "$TARGET_REPO_URL" =~ ^(https://|git@|ssh://) ]]; then
    echo -e "${YELLOW}⚠ Warning: URL doesn't start with https://, git@, or ssh://${NC}"
    echo -e "${YELLOW}Make sure this is a valid git repository URL.${NC}"
fi

prompt_with_default "Target repository base branch" "main" "TARGET_BASE_BRANCH"
prompt_with_default "Clone depth (0 for full history)" "100" "TARGET_CLONE_DEPTH"

# Step 2: Technology Stack
echo -e "\n${MAGENTA}${BOLD}═══ Step 2: Technology Stack ═══${NC}\n"

echo -e "${CYAN}What is your primary programming language?${NC}"
languages=("Go" "Python" "TypeScript" "Java" "Rust" "C++" "Other")
select_option "Select primary language:" "${languages[@]}"
PRIMARY_LANG="${languages[$SELECTED_INDEX]}"

if [ "$PRIMARY_LANG" = "Other" ]; then
    prompt_with_default "Enter language name" "" "PRIMARY_LANG"
fi

# Technology selection based on language
case "$PRIMARY_LANG" in
    "Go")
        echo -e "\n${CYAN}Select Go frameworks/libraries to include:${NC}"
        go_tech=("Kubernetes/KCP" "Gin Web Framework" "GORM" "Cobra CLI" "gRPC" "Prometheus" "None")
        TECH_STACK=$(select_multiple "Select technologies:" "${go_tech[@]}")
        ;;
    "Python")
        echo -e "\n${CYAN}Select Python frameworks/libraries to include:${NC}"
        python_tech=("Django" "FastAPI" "Flask" "SQLAlchemy" "Celery" "PyTorch" "NumPy/Pandas" "None")
        TECH_STACK=$(select_multiple "Select technologies:" "${python_tech[@]}")
        ;;
    "TypeScript")
        echo -e "\n${CYAN}Select TypeScript/JavaScript frameworks to include:${NC}"
        ts_tech=("React" "Next.js" "Node.js/Express" "NestJS" "Vue.js" "Angular" "None")
        TECH_STACK=$(select_multiple "Select technologies:" "${ts_tech[@]}")
        ;;
    *)
        prompt_with_default "Enter frameworks/libraries (comma-separated)" "" "TECH_STACK"
        ;;
esac

# Step 3: Agent Configuration
echo -e "\n${MAGENTA}${BOLD}═══ Step 3: Agent Configuration ═══${NC}\n"

echo -e "${CYAN}Which agents do you need for your project?${NC}"
agents=("Orchestrator (Required)" "Software Engineer" "Code Reviewer" "Architect" "Test Engineer" "DevOps Engineer")
SELECTED_AGENTS=$(select_multiple "Select agents:" "${agents[@]}")

# Always include Orchestrator
if [[ ! "$SELECTED_AGENTS" =~ "Orchestrator" ]]; then
    SELECTED_AGENTS="Orchestrator (Required) $SELECTED_AGENTS"
fi

# Expertise areas
echo -e "\n${CYAN}What expertise areas should agents have?${NC}"
echo -e "${YELLOW}You can either:${NC}"
echo -e "${YELLOW}1) Select from the list below (enter numbers)${NC}"
echo -e "${YELLOW}2) Enter 'custom' to provide your own expertise areas${NC}"

expertise=("Cloud Architecture" "Security" "Performance Optimization" "Database Design" "API Design" "Testing Strategies" "CI/CD" "Monitoring/Observability")

# Show options
for i in "${!expertise[@]}"; do
    echo -e "  ${BOLD}$((i+1)))${NC} ${expertise[$i]}"
done

echo -ne "${CYAN}Select options (numbers separated by spaces) or 'custom': ${NC}"
read -r expertise_choice

if [ "$expertise_choice" = "custom" ]; then
    echo -ne "${CYAN}Enter expertise areas (comma-separated): ${NC}"
    read -r EXPERTISE_AREAS
else
    # Process number selections
    EXPERTISE_AREAS=""
    for choice in $expertise_choice; do
        if [[ "$choice" =~ ^[0-9]+$ ]] && [ "$choice" -ge 1 ] && [ "$choice" -le "${#expertise[@]}" ]; then
            if [ -n "$EXPERTISE_AREAS" ]; then
                EXPERTISE_AREAS="$EXPERTISE_AREAS, ${expertise[$((choice-1))]}"
            else
                EXPERTISE_AREAS="${expertise[$((choice-1))]}"
            fi
        fi
    done
fi

# Step 4: Implementation Planning
echo -e "\n${MAGENTA}${BOLD}═══ Step 4: Implementation Planning ═══${NC}\n"

select_option "Do you have an existing implementation plan?" "No, generate one" "Yes, I'll provide it" "Use IDPBuilder OCI example (recommended for demo)"
HAS_PLAN=$SELECTED_INDEX

if [ $HAS_PLAN -eq 0 ]; then
    echo -e "\n${CYAN}What type of project are you building?${NC}"
    project_types=("Kubernetes Controller/Operator" "Web Application" "CLI Tool" "API Service" "Data Pipeline" "Library/SDK" "Other")
    select_option "Select project type:" "${project_types[@]}"
    PROJECT_TYPE="${project_types[$SELECTED_INDEX]}"
    
    if [ "$PROJECT_TYPE" = "Other" ]; then
        prompt_with_default "Describe your project type" "" "PROJECT_TYPE"
    fi
    
    prompt_with_default "Estimated lines of code" "5000" "ESTIMATED_LOC"
    prompt_with_default "Number of phases" "3" "NUM_PHASES"
    prompt_with_default "Test coverage target (%)" "80" "TEST_COVERAGE"
elif [ $HAS_PLAN -eq 2 ]; then
    # Use IDPBuilder OCI example
    echo -e "${GREEN}✓ Using IDPBuilder OCI Build & Push example project${NC}"
    echo -e "${CYAN}This will set up a project to add container build capabilities to IDPBuilder${NC}"
    PROJECT_TYPE="Kubernetes Controller Extension"
    ESTIMATED_LOC="8000"
    NUM_PHASES="5"
    TEST_COVERAGE="80"
    USE_EXAMPLE_PROMPT="true"
else
    prompt_with_default "Path to existing plan (or press Enter to add later)" "" "PLAN_PATH"
fi

# Step 5: Development Constraints
echo -e "\n${MAGENTA}${BOLD}═══ Step 5: Development Constraints ═══${NC}\n"

prompt_with_default "Maximum lines per effort" "800" "MAX_LINES"
prompt_with_default "Maximum parallel agents" "3" "MAX_PARALLEL"
prompt_with_default "Code review requirement" "mandatory" "REVIEW_REQUIREMENT"

select_option "Security level:" "Standard" "Enhanced" "Maximum"
SECURITY_LEVEL=$SELECTED_INDEX

# Step 6: Setup Confirmation
echo -e "\n${MAGENTA}${BOLD}═══ Setup Summary ═══${NC}\n"

echo -e "${BOLD}Project Configuration:${NC}"
echo -e "  ${CYAN}Name:${NC} $PROJECT_NAME"
echo -e "  ${CYAN}Description:${NC} $PROJECT_DESC"
echo -e "  ${CYAN}Directory:${NC} $TARGET_DIR"
echo -e "  ${CYAN}Language:${NC} $PRIMARY_LANG"
echo -e "  ${CYAN}Technologies:${NC} $TECH_STACK"
echo -e "  ${CYAN}Agents:${NC} $SELECTED_AGENTS"
echo -e "  ${CYAN}Expertise:${NC} $EXPERTISE_AREAS"

echo -e "\n${YELLOW}Proceed with setup? (y/n): ${NC}"
read -r confirm
if [ "$confirm" != "y" ]; then
    echo -e "${RED}Setup cancelled.${NC}"
    exit 1
fi

# Step 7: Create Project Structure
echo -e "\n${MAGENTA}${BOLD}═══ Creating Project Structure ═══${NC}\n"

# Check if target directory exists
if [ -d "$TARGET_DIR" ]; then
    echo -e "${YELLOW}⚠ Directory already exists: $TARGET_DIR${NC}"
    
    # Check if backup exists
    if [ -d "${TARGET_DIR}-bak" ]; then
        echo -e "${YELLOW}⚠ Backup directory also exists: ${TARGET_DIR}-bak${NC}"
        echo -e "\n${CYAN}What would you like to do?${NC}"
        echo -e "  ${BOLD}1)${NC} Delete the backup and move current to backup"
        echo -e "  ${BOLD}2)${NC} Move backup to ${TARGET_DIR}-bak-$(date +%Y%m%d-%H%M%S)"
        echo -e "  ${BOLD}3)${NC} Cancel setup"
        
        echo -ne "${CYAN}Select option (1-3): ${NC}"
        read -r backup_choice
        
        case "$backup_choice" in
            1)
                echo -e "${CYAN}Removing old backup...${NC}"
                rm -rf "${TARGET_DIR}-bak"
                echo -e "${CYAN}Moving current directory to backup...${NC}"
                mv "$TARGET_DIR" "${TARGET_DIR}-bak"
                echo -e "${GREEN}✓ Backup created${NC}"
                ;;
            2)
                new_backup="${TARGET_DIR}-bak-$(date +%Y%m%d-%H%M%S)"
                echo -e "${CYAN}Moving old backup to $new_backup...${NC}"
                mv "${TARGET_DIR}-bak" "$new_backup"
                echo -e "${CYAN}Moving current directory to backup...${NC}"
                mv "$TARGET_DIR" "${TARGET_DIR}-bak"
                echo -e "${GREEN}✓ Backups reorganized${NC}"
                ;;
            3)
                echo -e "${RED}Setup cancelled.${NC}"
                exit 1
                ;;
            *)
                echo -e "${RED}Invalid selection. Setup cancelled.${NC}"
                exit 1
                ;;
        esac
    else
        # No backup exists, offer to create one
        echo -e "\n${CYAN}Would you like to:${NC}"
        echo -e "  ${BOLD}1)${NC} Move existing directory to ${TARGET_DIR}-bak"
        echo -e "  ${BOLD}2)${NC} Delete existing directory (dangerous!)"
        echo -e "  ${BOLD}3)${NC} Cancel setup"
        
        echo -ne "${CYAN}Select option (1-3): ${NC}"
        read -r dir_choice
        
        case "$dir_choice" in
            1)
                echo -e "${CYAN}Moving existing directory to backup...${NC}"
                mv "$TARGET_DIR" "${TARGET_DIR}-bak"
                echo -e "${GREEN}✓ Backup created at ${TARGET_DIR}-bak${NC}"
                ;;
            2)
                echo -e "${RED}⚠ WARNING: This will permanently delete $TARGET_DIR${NC}"
                echo -ne "${YELLOW}Type 'DELETE' to confirm: ${NC}"
                read -r confirm_delete
                if [ "$confirm_delete" = "DELETE" ]; then
                    echo -e "${CYAN}Removing existing directory...${NC}"
                    rm -rf "$TARGET_DIR"
                    echo -e "${GREEN}✓ Directory removed${NC}"
                else
                    echo -e "${RED}Deletion not confirmed. Setup cancelled.${NC}"
                    exit 1
                fi
                ;;
            3)
                echo -e "${RED}Setup cancelled.${NC}"
                exit 1
                ;;
            *)
                echo -e "${RED}Invalid selection. Setup cancelled.${NC}"
                exit 1
                ;;
        esac
    fi
fi

# Create target directory
echo -e "${CYAN}Creating directory: $TARGET_DIR${NC}"
mkdir -p "$TARGET_DIR"

# Find the directory where this script is located (the template directory)
TEMPLATE_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Copy Software Factory 2.0 template
echo -e "${CYAN}Copying Software Factory 2.0 template from $TEMPLATE_DIR...${NC}"

# Copy all files except the scripts and .git
for item in "$TEMPLATE_DIR"/*; do
    basename_item=$(basename "$item")
    # Skip setup scripts, migration files, and template marker
    if [[ "$basename_item" != "setup.sh" && \
          "$basename_item" != "migrate-from-1.0.sh" && \
          "$basename_item" != "migrate-planning-only.sh" && \
          "$basename_item" != "MIGRATION-GUIDE-1.0-TO-2.0.md" && \
          "$basename_item" != "SF-1.0-VS-2.0-COMPARISON.md" && \
          "$basename_item" != "IDPBUILDER-MIGRATION-STRATEGY.md" && \
          "$basename_item" != "README.md" && \
          "$basename_item" != ".template-repository" ]]; then
        cp -r "$item" "$TARGET_DIR/"
    fi
done

# CRITICAL: Ensure .template-repository is NEVER copied to project instances
if [ -f "$TARGET_DIR/.template-repository" ]; then
    echo -e "${YELLOW}⚠ Removing .template-repository from project instance${NC}"
    rm -f "$TARGET_DIR/.template-repository"
fi

# Copy templates to project root for easy access
if [ -d "$TEMPLATE_DIR/templates" ]; then
    echo -e "${CYAN}Copying implementation plan templates...${NC}"
    cp -r "$TEMPLATE_DIR/templates" "$TARGET_DIR/"
    echo -e "${GREEN}✓ Templates copied to $TARGET_DIR/templates/${NC}"
fi

# Copy entire .claude folder as single source of truth
if [ -d "$TEMPLATE_DIR/.claude" ]; then
    echo -e "${CYAN}Copying SF 2.0 .claude configuration directory...${NC}"
    cp -r "$TEMPLATE_DIR/.claude" "$TARGET_DIR/"
    
    # Update configs with project-specific information if needed
    if [ -d "$TARGET_DIR/.claude" ]; then
        # Update agent configs
        for agent_file in "$TARGET_DIR/.claude/agents"/*.md; do
            [ -f "$agent_file" ] || continue
            sed -i "s/\[PROJECT_NAME\]/$PROJECT_NAME/g" "$agent_file" 2>/dev/null || true
            sed -i "s/\[PRIMARY_LANG\]/$PRIMARY_LANG/g" "$agent_file" 2>/dev/null || true
        done
        # Update command configs
        for cmd_file in "$TARGET_DIR/.claude/commands"/*.md; do
            [ -f "$cmd_file" ] || continue
            sed -i "s/\[PROJECT_NAME\]/$PROJECT_NAME/g" "$cmd_file" 2>/dev/null || true
            sed -i "s/\[PRIMARY_LANG\]/$PRIMARY_LANG/g" "$cmd_file" 2>/dev/null || true
        done
        echo -e "${GREEN}✓ .claude directory configured${NC}"
        
        # Ensure CLAUDE.md is present and configured
        if [ -f "$TARGET_DIR/.claude/CLAUDE.md" ]; then
            echo -e "${GREEN}✓ CLAUDE.md project configuration installed${NC}"
            echo -e "${YELLOW}   Contains grading criteria, TODO persistence rules, and agent configs${NC}"
        else
            echo -e "${YELLOW}⚠ CLAUDE.md not found - copying from template${NC}"
            if [ -f "$TEMPLATE_DIR/.claude/CLAUDE.md" ]; then
                cp "$TEMPLATE_DIR/.claude/CLAUDE.md" "$TARGET_DIR/.claude/"
                echo -e "${GREEN}✓ CLAUDE.md copied from template${NC}"
            fi
        fi
        
        # Settings.json uses relative paths, no updates needed
        if [ -f "$TARGET_DIR/.claude/settings.json" ]; then
            echo -e "${GREEN}✓ Claude Code hooks configured (PreCompact only)${NC}"
            echo -e "${YELLOW}   Note: PreCompact hook will create /tmp/compaction_marker.txt${NC}"
        fi
    fi
fi

# Copy essential Software Factory files
echo -e "\n${CYAN}Copying essential Software Factory 2.0 files...${NC}"

# Copy state machine definition
if [ -f "$TEMPLATE_DIR/state-machines/software-factory-3.0-state-machine.json" ]; then
    mkdir -p "$TARGET_DIR/state-machines"
    cp "$TEMPLATE_DIR/state-machines/software-factory-3.0-state-machine.json" "$TARGET_DIR/state-machines/"
    echo -e "${GREEN}✓ State machine definition copied${NC}"
fi

# Copy rule library
if [ -d "$TEMPLATE_DIR/rule-library" ]; then
    cp -r "$TEMPLATE_DIR/rule-library" "$TARGET_DIR/"
    echo -e "${GREEN}✓ Rule library copied${NC}"
fi

# Copy and configure utility scripts
if [ -d "$TEMPLATE_DIR/utilities" ]; then
    echo -e "${CYAN}Installing SF 2.0 utility scripts...${NC}"
    
    # Copy to project directory
    cp -r "$TEMPLATE_DIR/utilities" "$TARGET_DIR/"
    chmod +x "$TARGET_DIR/utilities"/*.sh
    echo -e "${GREEN}✓ Utility scripts installed to project directory${NC}"
    
    # ALSO install to user's home directory for consistent access
    echo -e "${CYAN}Installing utilities to ~/.claude/utilities for global access...${NC}"
    mkdir -p "$HOME/.claude/utilities"
    cp "$TEMPLATE_DIR/utilities"/*.sh "$HOME/.claude/utilities/"
    chmod +x "$HOME/.claude/utilities"/*.sh
    echo -e "${GREEN}✓ Utility scripts installed to ~/.claude/utilities${NC}"
    echo -e "${YELLOW}   Agents will automatically use utilities from this location${NC}"
    
    # Create required directories for utilities
    mkdir -p "$TARGET_DIR/todos"
    mkdir -p "$TARGET_DIR/checkpoints"
    mkdir -p "$TARGET_DIR/snapshots"
    echo -e "${GREEN}✓ Support directories created (todos/, checkpoints/, snapshots/)${NC}"
    
    echo -e "${YELLOW}Note: Utilities are available in two locations:${NC}"
    echo -e "${YELLOW}  - Project: ./utilities/ (for manual use)${NC}"
    echo -e "${YELLOW}  - Global: ~/.claude/utilities/ (agents use this automatically)${NC}"
fi

# Install .claude directory to user's home for global access
echo -e "\n${CYAN}Installing Claude configurations globally...${NC}"

if [ -d "$TEMPLATE_DIR/.claude" ]; then
    # Copy entire .claude directory to home
    mkdir -p "$HOME/.claude"
    cp -r "$TEMPLATE_DIR/.claude/agents" "$HOME/.claude/" 2>/dev/null || true
    cp -r "$TEMPLATE_DIR/.claude/commands" "$HOME/.claude/" 2>/dev/null || true
    echo -e "${GREEN}✓ Claude configurations installed to ~/.claude/${NC}"
    echo -e "${GREEN}  - Agent configurations: ~/.claude/agents/${NC}"
    echo -e "${GREEN}  - Command configurations: ~/.claude/commands/${NC}"
fi

# Copy settings.json if it exists
if [ -f "$TEMPLATE_DIR/.claude/settings.json" ]; then
    mkdir -p "$HOME/.claude"
    # Check if settings.json already exists
    if [ -f "$HOME/.claude/settings.json" ]; then
        echo -e "${YELLOW}⚠ ~/.claude/settings.json already exists${NC}"
        echo -ne "${CYAN}Overwrite it? (y/n): ${NC}"
        read -r overwrite_settings
        if [ "$overwrite_settings" = "y" ]; then
            cp "$TEMPLATE_DIR/.claude/settings.json" "$HOME/.claude/"
            echo -e "${GREEN}✓ Settings installed to ~/.claude/settings.json${NC}"
        else
            echo -e "${YELLOW}✓ Keeping existing settings.json${NC}"
        fi
    else
        cp "$TEMPLATE_DIR/.claude/settings.json" "$HOME/.claude/"
        echo -e "${GREEN}✓ Settings installed to ~/.claude/settings.json${NC}"
    fi
fi

echo -e "${YELLOW}Note: Claude configurations are now globally available:${NC}"
echo -e "${YELLOW}  - Agents: ~/.claude/agents/${NC}"
echo -e "${YELLOW}  - Commands: ~/.claude/commands/${NC}"
echo -e "${YELLOW}  - Settings: ~/.claude/settings.json${NC}"

# Step 8: Configure Agents
echo -e "\n${MAGENTA}${BOLD}═══ Configuring Agents ═══${NC}\n"

# Create project-specific configuration
# Escape quotes in description for YAML
ESCAPED_DESC="${PROJECT_DESC//\"/\\\"}"
ESCAPED_DESC="${ESCAPED_DESC//\'/\'\'}"

cat > "$TARGET_DIR/project-config.yaml" << EOF
# Software Factory 2.0 Project Configuration
# Generated: $(date)

project:
  name: "$PROJECT_NAME"
  description: "$ESCAPED_DESC"
  language: "$PRIMARY_LANG"
  repository: "$GITHUB_URL"
  
technology_stack:
$(echo "$TECH_STACK" | tr ' ' '\n' | sed 's/^/  - /')

agents:
$(echo "$SELECTED_AGENTS" | tr ' ' '\n' | sed 's/^/  - /')

expertise_areas:
$(echo "$EXPERTISE_AREAS" | tr ',' '\n' | sed 's/^[[:space:]]*//;s/[[:space:]]*$//' | sed 's/^/  - /')

constraints:
  max_lines_per_effort: $MAX_LINES
  max_parallel_agents: $MAX_PARALLEL
  code_review: $REVIEW_REQUIREMENT
  test_coverage_target: $TEST_COVERAGE
  security_level: $SECURITY_LEVEL

implementation:
  estimated_loc: ${ESTIMATED_LOC:-TBD}
  phases: ${NUM_PHASES:-TBD}
  project_type: "${PROJECT_TYPE:-TBD}"
EOF

# Create target repository configuration
echo -e "${CYAN}Creating target repository configuration...${NC}"

cat > "$TARGET_DIR/target-repo-config.yaml" << EOF
# Target Repository Configuration
# This file specifies the actual project repository that Software Factory will work on
# Generated: $(date)

target_repository:
  # The upstream repository URL
  url: "$TARGET_REPO_URL"
  
  # The default base branch to work from
  base_branch: "$TARGET_BASE_BRANCH"
  
  # Clone depth for sparse checkouts (0 for full history)
  clone_depth: $TARGET_CLONE_DEPTH
  
  # Authentication method
  auth_method: "https"

# Branch naming configuration
branch_naming:
  # Project prefix for all branches (lowercase kebab-case)
  # This will be prepended to all branch names with a slash
  project_prefix: "$(echo "$PROJECT_NAME" | tr '[:upper:]' '[:lower:]' | sed 's/[^a-z0-9-]/-/g' | sed 's/--*/-/g' | sed 's/^-//' | sed 's/-$//')"
  
  # Format for effort branches
  # Variables: {prefix}, {phase}, {wave}, {effort_name}
  # Note: {prefix} includes the trailing slash if project_prefix is set
  effort_format: "{prefix}phase{phase}/wave{wave}/{effort_name}"
  
  # Format for integration branches
  # Variables: {prefix}, {phase}, {wave}
  integration_format: "{prefix}phase{phase}/wave{wave}/integration"
  
  # Format for phase integration branches
  # Variables: {prefix}, {phase}
  phase_integration_format: "{prefix}phase{phase}/integration"

# Workspace configuration
workspace:
  # Root directory for effort workspaces
  efforts_root: "efforts"
  
  # Structure under efforts_root
  effort_path: "phase{phase}/wave{wave}/{effort_name}"
  
  # Sparse checkout patterns (empty for full checkout)
  sparse_patterns: []

# Repository separation enforcement
separation:
  # Prevent agents from modifying SF instance repository
  protect_sf_instance: true
  
  # List of protected paths in SF instance
  protected_paths:
    - "rule-library/"
    - "utilities/"
    - ".claude/"
    - "*.yaml"
    - "*.md"
  
  # Allowed operations in SF instance repo
  allowed_in_sf_instance:
    - "read"
    - "state"

# Validation rules
validation:
  # Require all effort branches to track remote
  require_remote_tracking: true
  
  # Require push verification after commits
  require_push_verification: true
  
  # Maximum time between push attempts (minutes)
  max_push_delay: 5
  
  # Require clean git status before switching efforts
  require_clean_status: true
EOF

echo -e "${GREEN}✓ Target repository configured: $TARGET_REPO_URL${NC}"

# Update agent configurations with project specifics
echo -e "${CYAN}Updating agent configurations...${NC}"

# Replace placeholders in critical files
find "$TARGET_DIR" -type f -name "*.md" -exec sed -i \
    -e "s/\[project\]/$PROJECT_NAME/g" \
    -e "s/\[LANG\]/$PRIMARY_LANG/g" \
    -e "s/\[TEST_COVERAGE\]/$TEST_COVERAGE/g" \
    {} \;

# Create initial orchestrator state
cat > "$TARGET_DIR/orchestrator-state-v3.json" << EOF
# Orchestrator State for $PROJECT_NAME
# Generated: $(date)

current_phase: 0
current_wave: 0
current_state: INIT

project_info:
  name: "$PROJECT_NAME"
  description: "$ESCAPED_DESC"
  start_date: "$(date +%Y-%m-%d)"
  
phases_planned: $NUM_PHASES
waves_per_phase: []

efforts_completed: []
efforts_in_progress: []
efforts_planned: []

integration_branches: []

grading_history:
  parallel_spawn_average: 0.0
  review_first_try_rate: 0.0
  integration_success_rate: 0.0
  
last_checkpoint: "$(date -Iseconds)"
EOF

# Step 9: Initialize Git Repository
echo -e "\n${MAGENTA}${BOLD}═══ Initializing Git Repository ═══${NC}\n"

cd "$TARGET_DIR"

if [ ! -d .git ]; then
    git init
    git add .
    git commit -m "Initial Software Factory 2.0 setup for $PROJECT_NAME"
    echo -e "${GREEN}✓ Git repository initialized${NC}"
fi

# Create 2.0 branch
git checkout -b software-factory-2.0
echo -e "${GREEN}✓ Created software-factory-2.0 branch${NC}"

# Install Python dependencies required by validation scripts
echo -e "\n${CYAN}Installing Python dependencies for validation...${NC}"

# Function to ensure jsonschema is installed
ensure_jsonschema_installed() {
    # Check if Python3 is available
    if ! command -v python3 &> /dev/null; then
        echo -e "${YELLOW}⚠ Python3 not found, skipping jsonschema installation${NC}"
        echo -e "${YELLOW}  Some validation features will be limited${NC}"
        return 1
    fi

    # Check if jsonschema is already installed and accessible in clean environment
    if python3 -c "import jsonschema" &> /dev/null; then
        # Test in clean environment (like agent-spawned shells)
        if env -i HOME="$HOME" PATH="/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin" \
           bash -c "python3 -c 'import jsonschema' 2>/dev/null"; then
            echo -e "${GREEN}✓ Python jsonschema module already installed and globally accessible${NC}"
            return 0
        else
            echo -e "${YELLOW}⚠ jsonschema found but not accessible in clean environments${NC}"
        fi
    fi

    echo -e "${BLUE}ℹ jsonschema module not found or not globally accessible, attempting to install...${NC}"

    # Check for Software Factory venv approach
    local SF_VENV_DIR="${HOME}/.software-factory-venv"
    local use_venv=false

    # Ask user preference for installation method
    echo -e "${CYAN}Choose jsonschema installation method:${NC}"
    echo -e "  ${BOLD}1)${NC} Install to user site-packages (default)"
    echo -e "  ${BOLD}2)${NC} Create dedicated Software Factory virtual environment (most reliable)"
    echo -ne "${CYAN}Select option (1-2) [1]: ${NC}"
    read -r install_choice

    if [ "$install_choice" = "2" ]; then
        use_venv=true
    fi

    if [ "$use_venv" = true ]; then
        # Create Software Factory virtual environment
        echo -e "${CYAN}Creating Software Factory virtual environment...${NC}"

        # Remove old venv if exists
        [ -d "$SF_VENV_DIR" ] && rm -rf "$SF_VENV_DIR"

        # Create new venv
        if python3 -m venv "$SF_VENV_DIR"; then
            # Install jsonschema in venv
            if "$SF_VENV_DIR/bin/pip" install --upgrade pip jsonschema &> /dev/null; then
                echo -e "${GREEN}✓ jsonschema installed in Software Factory venv${NC}"

                # Add activation to shell profiles
                local activation_snippet='
# Software Factory 2.0 - Auto-activate virtual environment for jsonschema
if [ -f "$HOME/.software-factory-venv/bin/activate" ] && [ -z "$VIRTUAL_ENV" ]; then
    source "$HOME/.software-factory-venv/bin/activate" 2>/dev/null
    export SF_VENV_ACTIVE=1
fi'

                # Add to .bashrc if not already there
                if [ -f "$HOME/.bashrc" ]; then
                    if ! grep -q "Software Factory 2.0 - Auto-activate virtual environment" "$HOME/.bashrc" 2>/dev/null; then
                        echo "$activation_snippet" >> "$HOME/.bashrc"
                        echo -e "${GREEN}✓ Added venv activation to ~/.bashrc${NC}"
                    fi
                fi

                # Add to .profile for sh compatibility
                if [ -f "$HOME/.profile" ] || [ ! -f "$HOME/.bashrc" ]; then
                    if ! grep -q "Software Factory 2.0 - Auto-activate virtual environment" "$HOME/.profile" 2>/dev/null; then
                        echo "$activation_snippet" >> "$HOME/.profile"
                        echo -e "${GREEN}✓ Added venv activation to ~/.profile${NC}"
                    fi
                fi

                echo -e "${GREEN}✓ Software Factory venv configured successfully${NC}"
                echo -e "${YELLOW}  Note: Run 'source ~/.bashrc' to activate in current shell${NC}"
                return 0
            fi
        fi

        echo -e "${YELLOW}⚠ Failed to create virtual environment, falling back to user install${NC}"
    fi

    # Standard user installation approach
    if command -v pip3 &> /dev/null; then
        if pip3 install --user jsonschema &> /dev/null; then
            echo -e "${GREEN}✓ jsonschema module installed successfully${NC}"
            return 0
        else
            # Try with --break-system-packages flag for newer systems
            if pip3 install --user --break-system-packages jsonschema &> /dev/null 2>&1; then
                echo -e "${GREEN}✓ jsonschema module installed successfully${NC}"
                return 0
            fi
        fi
    fi

    # Try to install using pip if pip3 failed
    if command -v pip &> /dev/null; then
        if pip install --user jsonschema &> /dev/null; then
            echo -e "${GREEN}✓ jsonschema module installed successfully${NC}"
            return 0
        else
            # Try with --break-system-packages flag
            if pip install --user --break-system-packages jsonschema &> /dev/null 2>&1; then
                echo -e "${GREEN}✓ jsonschema module installed successfully${NC}"
                return 0
            fi
        fi
    fi

    echo -e "${YELLOW}⚠ Could not install jsonschema automatically${NC}"
    echo -e "${YELLOW}  You can manually install it with one of:${NC}"
    echo -e "${YELLOW}    pip3 install --user jsonschema${NC}"
    echo -e "${YELLOW}    python3 -m venv ~/.software-factory-venv && ~/.software-factory-venv/bin/pip install jsonschema${NC}"
    echo -e "${YELLOW}  Some validation features will be limited without it${NC}"
    return 1
}

# Ensure jsonschema is installed for state validation
ensure_jsonschema_installed || true

# Install pre-commit hooks using the automated installation script
echo -e "\n${CYAN}Installing pre-commit hooks for repository separation enforcement...${NC}"

# First ensure hook templates are available
if [ ! -d "$TARGET_DIR/.git-hooks" ]; then
    mkdir -p "$TARGET_DIR/.git-hooks"
    if [ -f "$TEMPLATE_DIR/.git-hooks/planning-pre-commit.sh" ]; then
        cp "$TEMPLATE_DIR/.git-hooks/planning-pre-commit.sh" "$TARGET_DIR/.git-hooks/"
        echo -e "${GREEN}✓ Planning pre-commit hook template copied${NC}"
    fi
    if [ -f "$TEMPLATE_DIR/.git-hooks/effort-pre-commit.sh" ]; then
        cp "$TEMPLATE_DIR/.git-hooks/effort-pre-commit.sh" "$TARGET_DIR/.git-hooks/"
        echo -e "${GREEN}✓ Effort pre-commit hook template copied${NC}"
    fi
fi

# Check if the automated installation script exists
if [ -f "$TARGET_DIR/utilities/install-pre-commit-hooks.sh" ]; then
    echo -e "${BLUE}ℹ Using automated pre-commit hook installation${NC}"

    # Run the installation script
    cd "$TARGET_DIR"
    if bash utilities/install-pre-commit-hooks.sh; then
        echo -e "${GREEN}✓ Pre-commit hooks installed successfully${NC}"
    else
        echo -e "${YELLOW}⚠ Pre-commit hook installation reported warnings${NC}"
    fi
else
    # Fallback to simple inline hook if automated script not available
    echo -e "${YELLOW}⚠ Automated installer not found, using basic hook${NC}"
    mkdir -p .git/hooks

    # Try to use the planning hook template if available
    if [ -f "$TARGET_DIR/.git-hooks/planning-pre-commit.sh" ]; then
        cp "$TARGET_DIR/.git-hooks/planning-pre-commit.sh" .git/hooks/pre-commit
        chmod +x .git/hooks/pre-commit
        echo -e "${GREEN}✓ Planning pre-commit hook installed${NC}"
    else
        # Ultimate fallback to inline basic hook
        cat > .git/hooks/pre-commit << 'HOOK_EOF'
#!/bin/bash

current_branch="$(git branch --show-current)"
if [[ "software-factory-2.0" != "$current_branch" ]]; then
    echo "ERROR: DO NOT COMMIT WORK TO THE PLANNING REPO! YOU ARE ON THE WRONG BRANCH FOR PLANS! USE THE TARGET REPO FOR ALL CODE RELATED WORK! IMMEDIATELY READ $CLAUDE_PROJECT_DIR/rule-library/R251-REPOSITORY-SEPARATION-LAW.md"
    exit 1
fi

exit 0
HOOK_EOF
        chmod +x .git/hooks/pre-commit
        echo -e "${GREEN}✓ Basic pre-commit hook installed${NC}"
    fi
fi

echo -e "${GREEN}✓ Pre-commit hooks configured (enforces R251 - Repository Separation)${NC}"
echo -e "${YELLOW}  Hooks will prevent commits to incorrect branches/repositories${NC}"

# Add remote if provided
if [ -n "$GITHUB_URL" ]; then
    git remote add origin "$GITHUB_URL" 2>/dev/null || echo -e "${YELLOW}Remote already exists${NC}"
    echo -e "${GREEN}✓ Added remote: $GITHUB_URL${NC}"
fi

# Step 10: Generate Initial Plan (if needed)
if [ $HAS_PLAN -eq 0 ] || [ "$USE_EXAMPLE_PROMPT" = "true" ]; then
    echo -e "\n${MAGENTA}${BOLD}═══ Generating Implementation Plan ═══${NC}\n"
    
    # If using example, copy the architect prompt
    if [ "$USE_EXAMPLE_PROMPT" = "true" ]; then
        echo -e "${CYAN}Copying IDPBuilder OCI example architect prompt...${NC}"
        if [ -f "$TEMPLATE_DIR/ARCHITECT-PROMPT-IDPBUILDER-OCI.md" ]; then
            cp "$TEMPLATE_DIR/ARCHITECT-PROMPT-IDPBUILDER-OCI.md" "$TARGET_DIR/"
            echo -e "${GREEN}✓ Architect prompt copied${NC}"
        fi
        if [ -f "$TEMPLATE_DIR/ARCHITECT-TASK-COMMAND.md" ]; then
            cp "$TEMPLATE_DIR/ARCHITECT-TASK-COMMAND.md" "$TARGET_DIR/"
            echo -e "${GREEN}✓ Architect task command copied${NC}"
        fi
        echo -e "${YELLOW}⚠ The orchestrator will spawn an architect to create the master plan on first run${NC}"
    fi
    
    # Copy master template and customize it (but don't overwrite existing)
    if [ -f "$TARGET_DIR/templates/MASTER-IMPLEMENTATION-PLAN.md" ] && [ "$USE_EXAMPLE_PROMPT" != "true" ]; then
        echo -e "${CYAN}Using comprehensive master plan template...${NC}"
        
        # Check if IMPLEMENTATION-PLAN.md already exists
        if [ -f "$TARGET_DIR/IMPLEMENTATION-PLAN.md" ]; then
            echo -e "${YELLOW}⚠ IMPLEMENTATION-PLAN.md already exists${NC}"
            echo -e "${CYAN}Creating new template as IMPLEMENTATION-PLAN-NEW.md instead${NC}"
            cp "$TARGET_DIR/templates/MASTER-IMPLEMENTATION-PLAN.md" "$TARGET_DIR/IMPLEMENTATION-PLAN-NEW.md"
            PLAN_FILE="$TARGET_DIR/IMPLEMENTATION-PLAN-NEW.md"
        else
            cp "$TARGET_DIR/templates/MASTER-IMPLEMENTATION-PLAN.md" "$TARGET_DIR/IMPLEMENTATION-PLAN.md"
            PLAN_FILE="$TARGET_DIR/IMPLEMENTATION-PLAN.md"
        fi
        
        # Replace common placeholders
        sed -i "s/\[PROJECT_NAME\]/$PROJECT_NAME/g" "$PLAN_FILE"
        sed -i "s/\[LANGUAGE\]/$PRIMARY_LANG/g" "$PLAN_FILE"
        sed -i "s/\[TEST_COVERAGE\]/$TEST_COVERAGE/g" "$PLAN_FILE"
        sed -i "s/\[NUMBER\]/$ESTIMATED_LOC/g" "$PLAN_FILE"
        
        echo -e "${GREEN}✓ Master implementation plan created from template${NC}"
        echo -e "${YELLOW}⚠ Please review and complete: $PLAN_FILE${NC}"
        
        # Also create phase plans
        mkdir -p "$TARGET_DIR/phase-plans"
        for phase in {1..3}; do
            if [ -f "$TARGET_DIR/templates/PHASE-IMPLEMENTATION-PLAN.md" ]; then
                cp "$TARGET_DIR/templates/PHASE-IMPLEMENTATION-PLAN.md" "$TARGET_DIR/phase-plans/PHASE-$phase-PLAN.md"
                sed -i "s/\[NUMBER\]/$phase/g" "$TARGET_DIR/phase-plans/PHASE-$phase-PLAN.md"
                sed -i "s/\[PROJECT_NAME\]/$PROJECT_NAME/g" "$TARGET_DIR/phase-plans/PHASE-$phase-PLAN.md"
            fi
        done
        echo -e "${GREEN}✓ Phase plan templates created${NC}"
    else
        # Fallback to simple template if comprehensive template not found
        cat > "$TARGET_DIR/IMPLEMENTATION-PLAN-TEMPLATE.md" << EOF
# Implementation Plan for $PROJECT_NAME

## Project Overview
**Project:** $PROJECT_NAME
**Description:** ${PROJECT_DESC//\$/\\\$}
**Type:** $PROJECT_TYPE
**Language:** $PRIMARY_LANG
**Estimated Size:** $ESTIMATED_LOC lines

## Technology Stack
$(echo "$TECH_STACK" | tr ' ' '\n' | sed 's/^/- /')

## Phase Structure

### Phase 1: Foundation (Target: $(($ESTIMATED_LOC / $NUM_PHASES)) lines)
#### Wave 1: Core Infrastructure
- **Effort 1**: Project structure and configuration
- **Effort 2**: Base models/types
- **Effort 3**: Core utilities

#### Wave 2: [Define based on project needs]

### Phase 2: Core Features (Target: $(($ESTIMATED_LOC / $NUM_PHASES)) lines)
[Define waves and efforts]

### Phase 3: Advanced Features (Target: $(($ESTIMATED_LOC / $NUM_PHASES)) lines)
[Define waves and efforts]

## Constraints
- Maximum ${MAX_LINES} lines per effort
- Test coverage minimum: ${TEST_COVERAGE}%
- Code review: ${REVIEW_REQUIREMENT}

## Success Criteria
- All efforts under size limit
- Test coverage met
- All reviews passed
- Clean integration

---
**TODO**: Have an Architect agent review and complete this plan
EOF

        echo -e "${YELLOW}⚠ Basic implementation plan created at:${NC}"
        echo -e "  ${CYAN}$TARGET_DIR/IMPLEMENTATION-PLAN-TEMPLATE.md${NC}"
        echo -e "${YELLOW}  Please have an Architect agent review and complete it.${NC}"
    fi
fi

# Step 11: Create Quick Start Guide
echo -e "\n${MAGENTA}${BOLD}═══ Creating Quick Start Guide ═══${NC}\n"

cat > "$TARGET_DIR/QUICK-START.md" << EOF
# Quick Start Guide for $PROJECT_NAME

## 🚀 Your Software Factory 2.0 is Ready!

### Next Steps:

1. **Review the implementation plan:**
   \`\`\`bash
   cd $TARGET_DIR
   cat IMPLEMENTATION-PLAN-TEMPLATE.md
   \`\`\`

2. **Start the Orchestrator:**
   \`\`\`bash
   # Use the slash command in Claude:
   /continue-orchestrating
   \`\`\`

3. **The Orchestrator will:**
   - Load the current state
   - Spawn Code Reviewer for planning
   - Spawn SW Engineer for implementation
   - Manage the entire workflow

### Available Commands:

- \`/continue-orchestrating\` - Resume orchestration
- \`/continue-implementing\` - Resume implementation work
- \`/continue-reviewing\` - Resume code review
- \`/continue-architecting\` - Resume architecture review
- \`/check-status\` - Check current status
- \`/reset-state\` - Reset state (use carefully!)

### Key Files:

- \`project-config.yaml\` - Your project configuration
- \`orchestrator-state-v3.json\` - Current orchestration state
- \`.claude/commands/\` - Available slash commands
- \`rule-library/\` - Rule definitions and registry
- \`quick-reference/\` - Quick reference guides

### Monitoring Progress:

\`\`\`bash
# Check current state
cat orchestrator-state-v3.json

# View TODOs
ls -la todos/

# Check branch status
git status
\`\`\`

### Emergency Procedures:

If something goes wrong:
1. Run \`/check-status\` to diagnose
2. Check \`quick-reference/emergency-procedures.md\`
3. Use \`/reset-state\` if necessary (Level 1 first)

### Remember:

- **Never exceed $MAX_LINES lines per effort**
- **Always maintain ${TEST_COVERAGE}% test coverage**
- **Code review is ${REVIEW_REQUIREMENT}**
- **Maximum $MAX_PARALLEL parallel agents**

Good luck with your project!
EOF

# Step 12: Copy planning directory structure
echo -e "\n${CYAN}Setting up planning directory structure...${NC}"
if [ -d "$TEMPLATE_DIR/planning" ]; then
    cp -r "$TEMPLATE_DIR/planning" "$TARGET_DIR/"
    echo -e "${GREEN}✓ Planning directory structure copied with examples${NC}"
    echo -e "${YELLOW}   Includes example plans for project/phase/wave/effort levels${NC}"
else
    # Create basic planning structure if template doesn't exist
    echo -e "${YELLOW}⚠ Planning template not found, creating basic structure...${NC}"
    mkdir -p "$TARGET_DIR/planning/project"
    mkdir -p "$TARGET_DIR/planning/phase1/wave1"
    mkdir -p "$TARGET_DIR/planning/phase1/wave2"
    mkdir -p "$TARGET_DIR/planning/phase2/wave1"
    echo -e "${GREEN}✓ Basic planning directory structure created${NC}"
fi

# Step 13: Copy tools directory
echo -e "\n${CYAN}Copying tools directory...${NC}"
cp -r "$SCRIPT_DIR/tools" "$TARGET_DIR/"
chmod +x "$TARGET_DIR/tools"/*.sh
echo -e "${GREEN}✓ Tools copied and made executable${NC}"

# Step 12a: Copy orchestrator state schema file (critical for validation)
echo -e "\n${CYAN}Copying orchestrator state schema...${NC}"
if [ -f "$SCRIPT_DIR/orchestrator-state.schema.json" ]; then
    cp "$SCRIPT_DIR/orchestrator-state.schema.json" "$TARGET_DIR/"
    echo -e "${GREEN}✓ Orchestrator state schema copied${NC}"
else
    echo -e "${YELLOW}⚠ Warning: orchestrator-state.schema.json not found in template${NC}"
fi

# Step 13: Final Setup
echo -e "\n${MAGENTA}${BOLD}═══ Finalizing Setup ═══${NC}\n"

# Create necessary directories
mkdir -p "$TARGET_DIR/todos"
mkdir -p "$TARGET_DIR/checkpoints/active"
mkdir -p "$TARGET_DIR/efforts"
mkdir -p "$TARGET_DIR/integration"

# Utilities are already made executable above in the utilities section
# No need to update paths as utilities use relative paths

# Commit everything
git add .
git commit -m "Configure Software Factory 2.0 for $PROJECT_NAME

- Language: $PRIMARY_LANG
- Technology: $TECH_STACK
- Agents: $SELECTED_AGENTS
- Constraints: ${MAX_LINES} lines/effort, ${TEST_COVERAGE}% coverage"

echo -e "${GREEN}${BOLD}✅ Software Factory 2.0 Setup Complete!${NC}\n"

# Display summary
echo -e "${CYAN}${BOLD}═══════════════════════════════════════════════════════${NC}"
echo -e "${BOLD}Your project is ready at: ${GREEN}$TARGET_DIR${NC}"
echo -e "${CYAN}${BOLD}═══════════════════════════════════════════════════════${NC}\n"

echo -e "${YELLOW}📋 Next Steps:${NC}"
echo -e "  1. ${CYAN}cd $TARGET_DIR${NC}"
echo -e "  2. Review ${CYAN}QUICK-START.md${NC}"
echo -e "  3. Complete ${CYAN}IMPLEMENTATION-PLAN-TEMPLATE.md${NC}"
echo -e "  4. Run ${CYAN}/continue-orchestrating${NC} in Claude"

if [ -n "$GITHUB_URL" ]; then
    echo -e "\n${YELLOW}📦 To push to GitHub:${NC}"
    echo -e "  ${CYAN}git push -u origin software-factory-2.0${NC}"
fi

echo -e "\n${GREEN}${BOLD}Happy coding with Software Factory 2.0! 🚀${NC}"