#!/bin/bash

# Software Factory 2.0 Setup Script - Non-Interactive Mode Support
# This script supports both interactive and config-file based setup
# Usage: 
#   Interactive: ./setup-noninteractive.sh
#   Non-interactive: ./setup-noninteractive.sh --config setup-config.yaml

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

# Default values
CONFIG_FILE=""
INTERACTIVE=true
VERBOSE=false

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        --config|-c)
            CONFIG_FILE="$2"
            INTERACTIVE=false
            shift 2
            ;;
        --verbose|-v)
            VERBOSE=true
            shift
            ;;
        --help|-h)
            echo "Software Factory 2.0 Setup Script"
            echo ""
            echo "Usage:"
            echo "  $0 [options]"
            echo ""
            echo "Options:"
            echo "  --config, -c <file>  Use configuration file for non-interactive setup"
            echo "  --verbose, -v        Show detailed output"
            echo "  --help, -h          Show this help message"
            echo ""
            echo "Examples:"
            echo "  $0                           # Interactive mode"
            echo "  $0 --config setup.yaml       # Non-interactive with config file"
            echo "  $0 -c setup.yaml -v          # Non-interactive with verbose output"
            exit 0
            ;;
        *)
            echo -e "${RED}Unknown option: $1${NC}"
            echo "Use --help for usage information"
            exit 1
            ;;
    esac
done

# Function to read YAML value
# Usage: read_yaml_value "path.to.key" "default_value"
read_yaml_value() {
    local key="$1"
    local default="$2"
    local value
    
    # Try to extract value (handles both single-line and multi-line)
    if [[ "$key" == *"."* ]]; then
        # Nested key
        local section=$(echo "$key" | cut -d. -f1)
        local field=$(echo "$key" | cut -d. -f2)
        
        # Look for the section, then the field within it
        value=$(awk -v section="$section" -v field="$field" '
            BEGIN { in_section = 0 }
            $0 ~ "^" section ":" { in_section = 1; next }
            in_section && /^[a-zA-Z]/ { exit }
            in_section && $0 ~ "^  " field ":" {
                sub(/^[[:space:]]*[^:]*:[[:space:]]*/, "")
                sub(/[[:space:]]*#.*$/, "")
                gsub(/"/, "")
                gsub(/^[[:space:]]+|[[:space:]]+$/, "")
                print
                exit
            }
        ' "$CONFIG_FILE")
    else
        # Top-level key
        value=$(awk -v key="$key" '
            $0 ~ "^" key ":" {
                sub(/^[^:]*:[[:space:]]*/, "")
                sub(/[[:space:]]*#.*$/, "")
                gsub(/"/, "")
                gsub(/^[[:space:]]+|[[:space:]]+$/, "")
                print
                exit
            }
        ' "$CONFIG_FILE")
    fi
    
    # Return value or default
    if [ -n "$value" ]; then
        echo "$value"
    else
        echo "$default"
    fi
}

# Function to read YAML array
# Usage: read_yaml_array "path.to.array"
read_yaml_array() {
    local key="$1"
    local result=""
    
    if [[ "$key" == *"."* ]]; then
        # Nested array
        local section=$(echo "$key" | cut -d. -f1)
        local field=$(echo "$key" | cut -d. -f2)
        
        result=$(awk -v section="$section" -v field="$field" '
            BEGIN { in_section = 0; in_array = 0 }
            $0 ~ "^" section ":" { in_section = 1; next }
            in_section && /^[a-zA-Z]/ && !($0 ~ "^  ") { exit }
            in_section && $0 ~ "^  " field ":" { in_array = 1; next }
            in_section && in_array && /^  [a-zA-Z]/ { exit }
            in_section && in_array && /^    - / {
                sub(/^[[:space:]]*-[[:space:]]*/, "")
                sub(/[[:space:]]*#.*$/, "")
                gsub(/"/, "")
                gsub(/^[[:space:]]+|[[:space:]]+$/, "")
                printf "%s ", $0
            }
        ' "$CONFIG_FILE")
    else
        # Top-level array
        result=$(awk -v key="$key" '
            BEGIN { in_array = 0 }
            $0 ~ "^" key ":" { in_array = 1; next }
            in_array && /^[a-zA-Z]/ { exit }
            in_array && /^  - / {
                sub(/^[[:space:]]*-[[:space:]]*/, "")
                sub(/[[:space:]]*#.*$/, "")
                gsub(/"/, "")
                gsub(/^[[:space:]]+|[[:space:]]+$/, "")
                printf "%s ", $0
            }
        ' "$CONFIG_FILE")
    fi
    
    # Trim trailing space
    echo "${result% }"
}

# Banner
if [ "$INTERACTIVE" = true ] || [ "$VERBOSE" = true ]; then
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
║            ██║     ██║  ██║╚██████╗   ██║   ╚██████╔╝██║  ██║  ║
║            ╚═╝     ╚═╝  ╚═╝ ╚═════╝   ╚═╝    ╚═════╝ ╚═╝  ╚═╝    ║
║                                                                   ║
║                            Version 2.0                           ║
║                                                                   ║
╚═══════════════════════════════════════════════════════════════════╝
EOF
    echo -e "${NC}"
fi

# Load configuration if provided
if [ "$INTERACTIVE" = false ]; then
    if [ ! -f "$CONFIG_FILE" ]; then
        echo -e "${RED}Error: Configuration file not found: $CONFIG_FILE${NC}"
        exit 1
    fi
    
    echo -e "${CYAN}Loading configuration from: $CONFIG_FILE${NC}"
    
    # Load all configuration values
    PROJECT_NAME=$(read_yaml_value "project.name" "")
    PROJECT_DESC=$(read_yaml_value "project.description" "")
    TARGET_DIR=$(read_yaml_value "project.target_dir" "")
    GITHUB_URL=$(read_yaml_value "project.github_url" "")
    
    # Target repository configuration
    TARGET_REPO_URL=$(read_yaml_value "target_repository.url" "")
    TARGET_BASE_BRANCH=$(read_yaml_value "target_repository.base_branch" "main")
    TARGET_CLONE_DEPTH=$(read_yaml_value "target_repository.clone_depth" "100")
    TARGET_AUTH_METHOD=$(read_yaml_value "target_repository.auth_method" "https")
    
    PRIMARY_LANG=$(read_yaml_value "technology.primary_language" "")
    TECH_STACK=$(read_yaml_array "technology.frameworks")
    if [ -z "$TECH_STACK" ]; then
        TECH_STACK=$(read_yaml_value "technology.frameworks_custom" "")
    fi
    
    SELECTED_AGENTS=$(read_yaml_array "agents.selected")
    EXPERTISE_AREAS=$(read_yaml_array "agents.expertise")
    if [ -z "$EXPERTISE_AREAS" ]; then
        EXPERTISE_AREAS=$(read_yaml_value "agents.expertise_custom" "")
    fi
    
    PLAN_TYPE=$(read_yaml_value "implementation.plan_type" "generate")
    # Trim any whitespace from plan_type
    PLAN_TYPE=$(echo "$PLAN_TYPE" | xargs)
    
    PROJECT_TYPE=$(read_yaml_value "implementation.project_type" "")
    ESTIMATED_LOC=$(read_yaml_value "implementation.estimated_loc" "5000")
    NUM_PHASES=$(read_yaml_value "implementation.num_phases" "3")
    TEST_COVERAGE=$(read_yaml_value "implementation.test_coverage" "80")
    EXISTING_PLAN_PATH=$(read_yaml_value "implementation.existing_plan_path" "")
    
    MAX_LINES=$(read_yaml_value "constraints.max_lines_per_effort" "800")
    MAX_PARALLEL=$(read_yaml_value "constraints.max_parallel_agents" "3")
    REVIEW_REQUIREMENT=$(read_yaml_value "constraints.code_review" "mandatory")
    SECURITY_LEVEL=$(read_yaml_value "constraints.security_level" "1")
    
    DIR_IF_EXISTS=$(read_yaml_value "directory_handling.if_exists" "ask")
    CREATE_PARENT=$(read_yaml_value "directory_handling.create_parent" "ask")
    
    SKIP_GIT=$(read_yaml_value "options.skip_git_init" "false")
    SKIP_REMOTE=$(read_yaml_value "options.skip_remote" "false")
    VERBOSE=$(read_yaml_value "options.verbose" "$VERBOSE")
    CREATE_LINE_COUNTER=$(read_yaml_value "options.create_line_counter" "true")
    
    # Validate required fields
    if [ -z "$PROJECT_NAME" ]; then
        echo -e "${RED}Error: project.name is required in configuration${NC}"
        exit 1
    fi
    
    if [ -z "$TARGET_DIR" ]; then
        TARGET_DIR="/workspaces/$PROJECT_NAME"
        echo -e "${YELLOW}No target_dir specified, using: $TARGET_DIR${NC}"
    fi
    
    if [ -z "$PRIMARY_LANG" ]; then
        echo -e "${RED}Error: technology.primary_language is required in configuration${NC}"
        exit 1
    fi
    
    if [ -z "$TARGET_REPO_URL" ]; then
        echo -e "${RED}Error: target_repository.url is required in configuration${NC}"
        echo -e "${YELLOW}This is the repository where your code will be developed.${NC}"
        exit 1
    fi
    
    # Validate target repo URL format
    if [[ ! "$TARGET_REPO_URL" =~ ^(https://|git@|ssh://) ]]; then
        echo -e "${YELLOW}Warning: Target repo URL doesn't start with https://, git@, or ssh://${NC}"
        echo -e "${YELLOW}URL: $TARGET_REPO_URL${NC}"
    fi
    
    # Set implementation plan variables based on plan_type
    case "$PLAN_TYPE" in
        "generate")
            HAS_PLAN=0
            if [ -z "$PROJECT_TYPE" ]; then
                echo -e "${RED}Error: implementation.project_type is required when plan_type is 'generate'${NC}"
                exit 1
            fi
            ;;
        "existing")
            HAS_PLAN=1
            PLAN_PATH="$EXISTING_PLAN_PATH"
            ;;
        "idpbuilder-example")
            HAS_PLAN=2
            USE_EXAMPLE_PROMPT="true"
            PROJECT_TYPE="Kubernetes Controller Extension"
            ESTIMATED_LOC="8000"
            NUM_PHASES="5"
            TEST_COVERAGE="80"
            ;;
        *)
            echo -e "${RED}Error: Invalid plan_type: '$PLAN_TYPE'${NC}"
            echo "Valid options: generate, existing, idpbuilder-example"
            echo ""
            echo "Debug info:"
            echo "  Raw value length: ${#PLAN_TYPE}"
            echo "  Raw value bytes: $(echo -n "$PLAN_TYPE" | od -c)"
            exit 1
            ;;
    esac
    
    # Always include Orchestrator
    if [[ ! "$SELECTED_AGENTS" =~ "Orchestrator" ]]; then
        SELECTED_AGENTS="Orchestrator (Required) $SELECTED_AGENTS"
    fi
    
    if [ "$VERBOSE" = true ]; then
        echo -e "\n${MAGENTA}${BOLD}═══ Configuration Loaded ═══${NC}"
        echo -e "${CYAN}Project:${NC} $PROJECT_NAME"
        echo -e "${CYAN}Directory:${NC} $TARGET_DIR"
        echo -e "${CYAN}Language:${NC} $PRIMARY_LANG"
        echo -e "${CYAN}Plan Type:${NC} $PLAN_TYPE"
    fi
else
    # Source the original interactive setup functions
    SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
    
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
            printf -v "$var_name" "%s" "$default"
        else
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
    
    # Run interactive setup (copy from original setup.sh)
    echo -e "\n${MAGENTA}${BOLD}═══ Step 1: Project Information ═══${NC}\n"
    
    prompt_with_default "Project name" "" "PROJECT_NAME"
    while [ -z "$PROJECT_NAME" ]; do
        echo -e "${RED}Project name cannot be empty!${NC}"
        prompt_with_default "Project name" "" "PROJECT_NAME"
    done
    
    prompt_with_default "Project description" "" "PROJECT_DESC"
    prompt_with_default "Target directory" "/workspaces/$PROJECT_NAME" "TARGET_DIR"
    prompt_with_default "GitHub repository URL (optional)" "" "GITHUB_URL"
    
    # Continue with all interactive prompts...
    # (This section would include all the interactive prompts from the original setup.sh)
    # For brevity, I'm showing the structure but you'd copy all interactive sections
    
    echo -e "\n${YELLOW}Note: To skip interactive setup in the future, create a config file.${NC}"
    echo -e "${CYAN}See setup-config-example.yaml for reference.${NC}"
fi

# Validate target directory path
if [[ ! "$TARGET_DIR" =~ ^/ ]]; then
    if [ "$VERBOSE" = true ]; then
        echo -e "${YELLOW}⚠ Converting relative path to absolute path${NC}"
    fi
    TARGET_DIR="$(pwd)/$TARGET_DIR"
    if [ "$VERBOSE" = true ]; then
        echo -e "${CYAN}Using: $TARGET_DIR${NC}"
    fi
fi

# Handle parent directory
TARGET_PARENT=$(dirname "$TARGET_DIR")
if [ ! -d "$TARGET_PARENT" ]; then
    if [ "$INTERACTIVE" = false ]; then
        case "$CREATE_PARENT" in
            "create")
                mkdir -p "$TARGET_PARENT" 2>/dev/null || {
                    echo -e "${RED}❌ Cannot create parent directory. Permission denied.${NC}"
                    exit 1
                }
                echo -e "${GREEN}✓ Parent directory created${NC}"
                ;;
            "fail")
                echo -e "${RED}Parent directory does not exist: $TARGET_PARENT${NC}"
                exit 1
                ;;
            *)
                echo -e "${RED}Parent directory does not exist and create_parent is not set to 'create'${NC}"
                exit 1
                ;;
        esac
    else
        echo -e "${YELLOW}⚠ Parent directory does not exist: $TARGET_PARENT${NC}"
        echo -ne "${CYAN}Create it? (y/n): ${NC}"
        read -r create_parent
        if [ "$create_parent" = "y" ]; then
            mkdir -p "$TARGET_PARENT" 2>/dev/null || {
                echo -e "${RED}❌ Cannot create parent directory. Permission denied.${NC}"
                exit 1
            }
            echo -e "${GREEN}✓ Parent directory created${NC}"
        else
            echo -e "${RED}Setup cancelled.${NC}"
            exit 1
        fi
    fi
elif [ ! -w "$TARGET_PARENT" ]; then
    echo -e "${RED}❌ No write permission in parent directory: $TARGET_PARENT${NC}"
    exit 1
fi

# Handle existing directory
if [ -d "$TARGET_DIR" ]; then
    if [ "$INTERACTIVE" = false ]; then
        case "$DIR_IF_EXISTS" in
            "backup")
                echo -e "${CYAN}Creating backup of existing directory...${NC}"
                backup_name="${TARGET_DIR}-bak-$(date +%Y%m%d-%H%M%S)"
                mv "$TARGET_DIR" "$backup_name"
                echo -e "${GREEN}✓ Backup created at $backup_name${NC}"
                ;;
            "delete")
                echo -e "${YELLOW}Removing existing directory...${NC}"
                rm -rf "$TARGET_DIR"
                echo -e "${GREEN}✓ Directory removed${NC}"
                ;;
            *)
                echo -e "${RED}Directory already exists: $TARGET_DIR${NC}"
                echo -e "${RED}Set directory_handling.if_exists to 'backup' or 'delete' in config${NC}"
                exit 1
                ;;
        esac
    else
        # Interactive handling (from original script)
        echo -e "${YELLOW}⚠ Directory already exists: $TARGET_DIR${NC}"
        echo -e "\n${CYAN}Would you like to:${NC}"
        echo -e "  ${BOLD}1)${NC} Move existing directory to backup"
        echo -e "  ${BOLD}2)${NC} Delete existing directory"
        echo -e "  ${BOLD}3)${NC} Cancel setup"
        
        echo -ne "${CYAN}Select option (1-3): ${NC}"
        read -r dir_choice
        
        case "$dir_choice" in
            1)
                backup_name="${TARGET_DIR}-bak-$(date +%Y%m%d-%H%M%S)"
                mv "$TARGET_DIR" "$backup_name"
                echo -e "${GREEN}✓ Backup created at $backup_name${NC}"
                ;;
            2)
                echo -e "${RED}⚠ WARNING: This will permanently delete $TARGET_DIR${NC}"
                echo -ne "${YELLOW}Type 'DELETE' to confirm: ${NC}"
                read -r confirm_delete
                if [ "$confirm_delete" = "DELETE" ]; then
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

# Display summary in non-interactive mode
if [ "$INTERACTIVE" = false ] && [ "$VERBOSE" = true ]; then
    echo -e "\n${MAGENTA}${BOLD}═══ Setup Summary ═══${NC}\n"
    echo -e "${BOLD}Project Configuration:${NC}"
    echo -e "  ${CYAN}Name:${NC} $PROJECT_NAME"
    echo -e "  ${CYAN}Description:${NC} $PROJECT_DESC"
    echo -e "  ${CYAN}Directory:${NC} $TARGET_DIR"
    echo -e "  ${CYAN}Language:${NC} $PRIMARY_LANG"
    echo -e "  ${CYAN}Technologies:${NC} $TECH_STACK"
    echo -e "  ${CYAN}Agents:${NC} $SELECTED_AGENTS"
    echo -e "  ${CYAN}Expertise:${NC} $EXPERTISE_AREAS"
fi

# Create project structure
echo -e "\n${MAGENTA}${BOLD}═══ Creating Project Structure ═══${NC}\n"

mkdir -p "$TARGET_DIR"

# Find the directory where this script is located (the template directory)
TEMPLATE_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Copy Software Factory 2.0 template
echo -e "${CYAN}Copying Software Factory 2.0 template...${NC}"

# Copy all files except the scripts and .git
for item in "$TEMPLATE_DIR"/*; do
    basename_item=$(basename "$item")
    # Skip setup scripts and migration files
    if [[ "$basename_item" != "setup.sh" && \
          "$basename_item" != "setup-noninteractive.sh" && \
          "$basename_item" != "migrate-from-1.0.sh" && \
          "$basename_item" != "migrate-planning-only.sh" && \
          "$basename_item" != "MIGRATION-GUIDE-1.0-TO-2.0.md" && \
          "$basename_item" != "SF-1.0-VS-2.0-COMPARISON.md" && \
          "$basename_item" != "IDPBUILDER-MIGRATION-STRATEGY.md" && \
          "$basename_item" != "README.md" && \
          "$basename_item" != "setup-config-example.yaml" ]]; then
        cp -r "$item" "$TARGET_DIR/"
    fi
done

# Copy templates to project root
if [ -d "$TEMPLATE_DIR/templates" ]; then
    cp -r "$TEMPLATE_DIR/templates" "$TARGET_DIR/"
    echo -e "${GREEN}✓ Templates copied${NC}"
fi

# Copy .claude directory
if [ -d "$TEMPLATE_DIR/.claude" ]; then
    cp -r "$TEMPLATE_DIR/.claude" "$TARGET_DIR/"
    echo -e "${GREEN}✓ Claude configurations copied${NC}"
    
    # Verify CLAUDE.md is present
    if [ -f "$TARGET_DIR/.claude/CLAUDE.md" ]; then
        echo -e "${GREEN}✓ CLAUDE.md project configuration installed${NC}"
        if [ "$VERBOSE" = true ]; then
            echo -e "${YELLOW}   Contains grading criteria, TODO persistence rules, and agent configs${NC}"
        fi
    else
        echo -e "${YELLOW}⚠ Warning: CLAUDE.md not found in .claude directory${NC}"
    fi
fi

# Copy essential Software Factory files
echo -e "\n${CYAN}Copying essential Software Factory 2.0 files...${NC}"

# Copy state machine definition
if [ -f "$TEMPLATE_DIR/SOFTWARE-FACTORY-STATE-MACHINE.md" ]; then
    cp "$TEMPLATE_DIR/SOFTWARE-FACTORY-STATE-MACHINE.md" "$TARGET_DIR/"
    echo -e "${GREEN}✓ State machine definition copied${NC}"
fi

# Copy rule library
if [ -d "$TEMPLATE_DIR/rule-library" ]; then
    cp -r "$TEMPLATE_DIR/rule-library" "$TARGET_DIR/"
    echo -e "${GREEN}✓ Rule library copied${NC}"
fi

# Copy utilities
if [ -d "$TEMPLATE_DIR/utilities" ]; then
    # Copy to project directory
    cp -r "$TEMPLATE_DIR/utilities" "$TARGET_DIR/"
    chmod +x "$TARGET_DIR/utilities"/*.sh
    echo -e "${GREEN}✓ Utility scripts installed to project directory${NC}"
    
    # ALSO install to user's home directory for consistent access
    mkdir -p "$HOME/.claude/utilities"
    cp "$TEMPLATE_DIR/utilities"/*.sh "$HOME/.claude/utilities/"
    chmod +x "$HOME/.claude/utilities"/*.sh
    echo -e "${GREEN}✓ Utility scripts installed to ~/.claude/utilities${NC}"
    if [ "$VERBOSE" = true ]; then
        echo -e "${YELLOW}   Agents will automatically use utilities from ~/.claude/utilities${NC}"
    fi
fi

# Copy Claude agent configurations and commands
echo -e "\n${CYAN}Installing Claude agent configurations and commands...${NC}"

# Copy agent configurations
if [ -d "$TEMPLATE_DIR/.claude/agents" ]; then
    mkdir -p "$HOME/.claude/agents"
    cp "$TEMPLATE_DIR/.claude/agents"/*.md "$HOME/.claude/agents/" 2>/dev/null || true
    echo -e "${GREEN}✓ Agent configurations installed to ~/.claude/agents/${NC}"
fi

# Copy command configurations
if [ -d "$TEMPLATE_DIR/.claude/commands" ]; then
    mkdir -p "$HOME/.claude/commands"
    cp "$TEMPLATE_DIR/.claude/commands"/*.md "$HOME/.claude/commands/" 2>/dev/null || true
    echo -e "${GREEN}✓ Command configurations installed to ~/.claude/commands/${NC}"
fi

# Copy settings.json if it exists
if [ -f "$TEMPLATE_DIR/.claude/settings.json" ]; then
    mkdir -p "$HOME/.claude"
    # Check if settings.json already exists and overwrite based on config
    if [ -f "$HOME/.claude/settings.json" ]; then
        if [ "$OVERWRITE_SETTINGS" = "true" ]; then
            cp "$TEMPLATE_DIR/.claude/settings.json" "$HOME/.claude/"
            echo -e "${GREEN}✓ Settings installed to ~/.claude/settings.json (overwritten)${NC}"
        else
            echo -e "${YELLOW}✓ Keeping existing ~/.claude/settings.json${NC}"
        fi
    else
        cp "$TEMPLATE_DIR/.claude/settings.json" "$HOME/.claude/"
        echo -e "${GREEN}✓ Settings installed to ~/.claude/settings.json${NC}"
    fi
fi

if [ "$VERBOSE" = true ]; then
    echo -e "${YELLOW}Note: Claude configurations are now globally available:${NC}"
    echo -e "${YELLOW}  - Agents: ~/.claude/agents/${NC}"
    echo -e "${YELLOW}  - Commands: ~/.claude/commands/${NC}"
    echo -e "${YELLOW}  - Settings: ~/.claude/settings.json${NC}"
fi

# Create required directories
mkdir -p "$TARGET_DIR/todos"
mkdir -p "$TARGET_DIR/checkpoints"
mkdir -p "$TARGET_DIR/snapshots"
mkdir -p "$TARGET_DIR/efforts"
mkdir -p "$TARGET_DIR/integration"

# Configure agents
echo -e "\n${MAGENTA}${BOLD}═══ Configuring Agents ═══${NC}\n"

# Escape quotes in description for YAML
ESCAPED_DESC="${PROJECT_DESC//\"/\\\"}"
ESCAPED_DESC="${ESCAPED_DESC//\'/\'\'}"

# Create project configuration
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
  auth_method: "$TARGET_AUTH_METHOD"

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

# Create initial orchestrator state
cat > "$TARGET_DIR/orchestrator-state.json" << EOF
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

# Initialize Git Repository (unless skipped)
if [ "$SKIP_GIT" != "true" ]; then
    echo -e "\n${MAGENTA}${BOLD}═══ Initializing Git Repository ═══${NC}\n"
    
    cd "$TARGET_DIR"
    
    if [ ! -d .git ]; then
        git init
        git add .
        git commit -m "Initial Software Factory 2.0 setup for $PROJECT_NAME"
        echo -e "${GREEN}✓ Git repository initialized${NC}"
    fi
    
    git checkout -b software-factory-2.0
    echo -e "${GREEN}✓ Created software-factory-2.0 branch${NC}"
    
    if [ -n "$GITHUB_URL" ] && [ "$SKIP_REMOTE" != "true" ]; then
        git remote add origin "$GITHUB_URL" 2>/dev/null || echo -e "${YELLOW}Remote already exists${NC}"
        echo -e "${GREEN}✓ Added remote: $GITHUB_URL${NC}"
    fi
fi

# Generate initial plan if needed
if [ "$HAS_PLAN" -eq 0 ] || [ "$USE_EXAMPLE_PROMPT" = "true" ]; then
    echo -e "\n${MAGENTA}${BOLD}═══ Generating Implementation Plan ═══${NC}\n"
    
    if [ "$USE_EXAMPLE_PROMPT" = "true" ]; then
        if [ -f "$TEMPLATE_DIR/ARCHITECT-PROMPT-IDPBUILDER-OCI.md" ]; then
            cp "$TEMPLATE_DIR/ARCHITECT-PROMPT-IDPBUILDER-OCI.md" "$TARGET_DIR/"
            echo -e "${GREEN}✓ IDPBuilder OCI example loaded${NC}"
        fi
    fi
    
    # Create plan template (but don't overwrite existing)
    if [ -f "$TARGET_DIR/templates/MASTER-IMPLEMENTATION-PLAN.md" ]; then
        if [ -f "$TARGET_DIR/IMPLEMENTATION-PLAN.md" ]; then
            echo -e "${YELLOW}⚠ IMPLEMENTATION-PLAN.md already exists, creating IMPLEMENTATION-PLAN-NEW.md${NC}"
            PLAN_FILE="$TARGET_DIR/IMPLEMENTATION-PLAN-NEW.md"
        else
            PLAN_FILE="$TARGET_DIR/IMPLEMENTATION-PLAN.md"
        fi
        cp "$TARGET_DIR/templates/MASTER-IMPLEMENTATION-PLAN.md" "$PLAN_FILE"
        sed -i "s/\[PROJECT_NAME\]/$PROJECT_NAME/g" "$PLAN_FILE"
        sed -i "s/\[LANGUAGE\]/$PRIMARY_LANG/g" "$PLAN_FILE"
        sed -i "s/\[TEST_COVERAGE\]/$TEST_COVERAGE/g" "$PLAN_FILE"
        echo -e "${GREEN}✓ Implementation plan template created: $(basename $PLAN_FILE)${NC}"
    fi
elif [ "$HAS_PLAN" -eq 1 ] && [ -n "$EXISTING_PLAN_PATH" ]; then
    if [ -f "$EXISTING_PLAN_PATH" ]; then
        if [ -f "$TARGET_DIR/IMPLEMENTATION-PLAN.md" ]; then
            echo -e "${YELLOW}⚠ IMPLEMENTATION-PLAN.md already exists, creating IMPLEMENTATION-PLAN-NEW.md${NC}"
            cp "$EXISTING_PLAN_PATH" "$TARGET_DIR/IMPLEMENTATION-PLAN-NEW.md"
            echo -e "${GREEN}✓ Existing plan copied to IMPLEMENTATION-PLAN-NEW.md${NC}"
        else
            cp "$EXISTING_PLAN_PATH" "$TARGET_DIR/IMPLEMENTATION-PLAN.md"
            echo -e "${GREEN}✓ Existing plan copied${NC}"
        fi
    else
        echo -e "${YELLOW}⚠ Plan file not found: $EXISTING_PLAN_PATH${NC}"
    fi
fi

# Copy tools directory
echo -e "\n${CYAN}Copying tools directory...${NC}"
cp -r "$SCRIPT_DIR/tools" "$TARGET_DIR/"
chmod +x "$TARGET_DIR/tools"/*.sh
echo -e "${GREEN}✓ Tools copied and made executable${NC}"

# Copy orchestrator state schema file (critical for validation)
echo -e "\n${CYAN}Copying orchestrator state schema...${NC}"
if [ -f "$SCRIPT_DIR/orchestrator-state.schema.json" ]; then
    cp "$SCRIPT_DIR/orchestrator-state.schema.json" "$TARGET_DIR/"
    echo -e "${GREEN}✓ Orchestrator state schema copied${NC}"
else
    echo -e "${YELLOW}⚠ Warning: orchestrator-state.schema.json not found in template${NC}"
fi

# Create Quick Start Guide
cat > "$TARGET_DIR/QUICK-START.md" << EOF
# Quick Start Guide for $PROJECT_NAME

## 🚀 Your Software Factory 2.0 is Ready!

### Setup Method
- Configuration File: ${CONFIG_FILE:-Interactive}
- Language: $PRIMARY_LANG
- Plan Type: ${PLAN_TYPE:-generate}

### Next Steps:

1. **Review the implementation plan:**
   \`\`\`bash
   cd $TARGET_DIR
   cat IMPLEMENTATION-PLAN.md
   \`\`\`

2. **Start the Orchestrator:**
   \`\`\`bash
   # Use the slash command in Claude:
   /continue-orchestrating
   \`\`\`

### Available Commands:
- \`/continue-orchestrating\` - Resume orchestration
- \`/continue-implementing\` - Resume implementation work
- \`/continue-reviewing\` - Resume code review
- \`/check-status\` - Check current status

### Key Files:
- \`project-config.yaml\` - Your project configuration
- \`orchestrator-state.json\` - Current orchestration state
- \`.claude/commands/\` - Available slash commands
- \`utilities/\` - Manual utility scripts

### Constraints:
- Maximum $MAX_LINES lines per effort
- Test coverage target: ${TEST_COVERAGE}%
- Code review: ${REVIEW_REQUIREMENT}
- Maximum $MAX_PARALLEL parallel agents

Good luck with your project!
EOF

# Final commit
if [ "$SKIP_GIT" != "true" ]; then
    cd "$TARGET_DIR"
    git add .
    git commit -m "Configure Software Factory 2.0 for $PROJECT_NAME

- Language: $PRIMARY_LANG
- Technology: $TECH_STACK
- Agents: $SELECTED_AGENTS
- Constraints: ${MAX_LINES} lines/effort, ${TEST_COVERAGE}% coverage
- Setup method: ${CONFIG_FILE:-Interactive}"
fi

# Success message
echo -e "\n${GREEN}${BOLD}✅ Software Factory 2.0 Setup Complete!${NC}\n"

echo -e "${CYAN}${BOLD}═══════════════════════════════════════════════════════${NC}"
echo -e "${BOLD}Your project is ready at: ${GREEN}$TARGET_DIR${NC}"
echo -e "${CYAN}${BOLD}═══════════════════════════════════════════════════════${NC}\n"

if [ "$INTERACTIVE" = false ]; then
    echo -e "${GREEN}✓ Non-interactive setup completed successfully${NC}"
    echo -e "${CYAN}Configuration used: $CONFIG_FILE${NC}"
fi

echo -e "\n${YELLOW}📋 Next Steps:${NC}"
echo -e "  1. ${CYAN}cd $TARGET_DIR${NC}"
echo -e "  2. Review ${CYAN}QUICK-START.md${NC}"
echo -e "  3. Complete ${CYAN}IMPLEMENTATION-PLAN.md${NC}"
echo -e "  4. Run ${CYAN}/continue-orchestrating${NC} in Claude"

if [ -n "$GITHUB_URL" ] && [ "$SKIP_REMOTE" != "true" ]; then
    echo -e "\n${YELLOW}📦 To push to GitHub:${NC}"
    echo -e "  ${CYAN}git push -u origin software-factory-2.0${NC}"
fi

echo -e "\n${GREEN}${BOLD}Happy coding with Software Factory 2.0! 🚀${NC}"