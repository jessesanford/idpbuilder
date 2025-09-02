#!/bin/bash
################################################################################
# SOFTWARE FACTORY 2.0 - RECOVERY ASSISTANT
################################################################################
#
# This script provides interactive recovery assistance after context compaction
# or when agents lose track of their state. It guides users through proper
# recovery procedures and helps restore working context.
#
# FEATURES:
# - Interactive recovery workflow
# - Agent-specific guidance
# - TODO state recovery assistance  
# - Context file identification
# - State machine position detection
# - Validation of recovery completeness
#
# USAGE:
#   ./recovery-assistant.sh                # Interactive mode
#   ./recovery-assistant.sh --agent <type> # Agent-specific recovery
#   ./recovery-assistant.sh --validate     # Validate current state
#   ./recovery-assistant.sh --emergency    # Emergency recovery mode
#
################################################################################

set -euo pipefail

# Constants
readonly SCRIPT_NAME="$(basename "$0")"
readonly HOOKS_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
readonly FACTORY_ROOT="/workspaces/software-factory-2.0-template"
readonly LOG_FILE="/tmp/recovery-assistant.log"
readonly RECOVERY_STATE_FILE="/tmp/recovery-session.state"

# Colors for output
readonly RED='\033[0;31m'
readonly GREEN='\033[0;32m'
readonly YELLOW='\033[1;33m'
readonly BLUE='\033[0;34m'
readonly MAGENTA='\033[0;35m'
readonly CYAN='\033[0;36m'
readonly BOLD='\033[1m'
readonly NC='\033[0m' # No Color

# Recovery session state
declare -A RECOVERY_STATE
RECOVERY_STATE[agent_type]=""
RECOVERY_STATE[current_directory]="$(pwd)"
RECOVERY_STATE[todos_loaded]="false"
RECOVERY_STATE[context_files_read]="false"
RECOVERY_STATE[validation_passed]="false"

# Logging function
log() {
    local level="$1"
    shift
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] [$level] $*" | tee -a "$LOG_FILE"
}

# Output functions
print_banner() {
    echo -e "${CYAN}"
    echo "################################################################################"
    echo "#                 SOFTWARE FACTORY 2.0 - RECOVERY ASSISTANT                   #"
    echo "################################################################################"
    echo -e "${NC}"
}

print_header() {
    echo -e "${BOLD}${CYAN}━━━ $1 ━━━${NC}"
}

print_step() {
    echo -e "${BOLD}${BLUE}🔄 STEP $1:${NC} $2"
}

print_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

print_info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

print_critical() {
    echo -e "${BOLD}${RED}🚨 $1${NC}"
}

print_instruction() {
    echo -e "${MAGENTA}📋 $1${NC}"
}

# Interactive input functions
prompt_yes_no() {
    local question="$1"
    local default="${2:-n}"
    
    while true; do
        if [ "$default" = "y" ]; then
            read -p "$question (Y/n): " -r response
            response=${response:-y}
        else
            read -p "$question (y/N): " -r response
            response=${response:-n}
        fi
        
        case "$response" in
            [Yy]|[Yy][Ee][Ss]) return 0 ;;
            [Nn]|[Nn][Oo]) return 1 ;;
            *) echo "Please answer yes or no." ;;
        esac
    done
}

prompt_choice() {
    local question="$1"
    shift
    local choices=("$@")
    
    echo "$question"
    for i in "${!choices[@]}"; do
        echo "  $((i+1))) ${choices[$i]}"
    done
    
    while true; do
        read -p "Enter choice [1-${#choices[@]}]: " -r choice
        if [[ "$choice" =~ ^[0-9]+$ ]] && [ "$choice" -ge 1 ] && [ "$choice" -le "${#choices[@]}" ]; then
            echo "${choices[$((choice-1))]}"
            return 0
        fi
        echo "Invalid choice. Please enter a number between 1 and ${#choices[@]}."
    done
}

################################################################################
# USAGE AND HELP
################################################################################

usage() {
    cat << EOF
SOFTWARE FACTORY 2.0 - RECOVERY ASSISTANT

USAGE:
    $SCRIPT_NAME [options]

OPTIONS:
    --agent <type>      Start recovery for specific agent type
    --validate          Validate current recovery state
    --emergency         Emergency recovery mode (minimal prompts)
    --help              Show this help message

AGENT TYPES:
    orchestrator        orchestrator
    sw-engineer         sw-engineer
    code-reviewer       code-reviewer  
    architect           architect

RECOVERY PROCESS:
    1. Identify agent type and current context
    2. Locate and load TODO state files
    3. Read required context files
    4. Validate state machine position
    5. Provide continuation guidance

EXAMPLES:
    $SCRIPT_NAME                        # Interactive recovery
    $SCRIPT_NAME --agent orchestrator   # Orchestrator-specific recovery
    $SCRIPT_NAME --validate             # Check current state

EOF
}

################################################################################
# MAIN RECOVERY WORKFLOW
################################################################################

interactive_recovery() {
    print_banner
    echo "Welcome to the Software Factory 2.0 Recovery Assistant!"
    echo "This tool will help you recover your context after compaction or state loss."
    echo
    
    log "INFO" "Starting interactive recovery session"
    
    # Save recovery session state
    save_recovery_state
    
    # Step 1: Identify context
    print_step 1 "CONTEXT IDENTIFICATION"
    identify_context
    echo
    
    # Step 2: Detect agent type
    print_step 2 "AGENT TYPE DETECTION"
    detect_agent_type
    echo
    
    # Step 3: TODO state recovery
    print_step 3 "TODO STATE RECOVERY"
    recover_todo_state
    echo
    
    # Step 4: Context files recovery
    print_step 4 "CONTEXT FILES RECOVERY"
    recover_context_files
    echo
    
    # Step 5: State validation
    print_step 5 "RECOVERY VALIDATION"
    validate_recovery
    echo
    
    # Step 6: Continuation guidance
    print_step 6 "CONTINUATION GUIDANCE"
    provide_continuation_guidance
    echo
    
    # Complete recovery
    complete_recovery
}

identify_context() {
    print_info "Analyzing current environment..."
    
    local current_dir="$(pwd)"
    local git_branch="$(git branch --show-current 2>/dev/null || echo 'not-a-repo')"
    local git_repo="$(git rev-parse --show-toplevel 2>/dev/null || echo 'not-a-repo')"
    
    echo "Current Directory: $current_dir"
    echo "Git Repository: $git_repo"
    echo "Git Branch: $git_branch"
    
    RECOVERY_STATE[current_directory]="$current_dir"
    RECOVERY_STATE[git_branch]="$git_branch"
    RECOVERY_STATE[git_repo]="$git_repo"
    
    # Check for compaction markers
    if [ -f "/tmp/compaction_marker.txt" ]; then
        print_success "Compaction marker found - loading preserved context"
        source "/tmp/compaction_marker.txt" 2>/dev/null || true
        
        if [ -n "${DETECTED_AGENT:-}" ]; then
            RECOVERY_STATE[agent_type]="$DETECTED_AGENT"
            print_info "Pre-compaction agent detected: $DETECTED_AGENT"
        fi
        
        if [ -n "${WORKING_DIRECTORY:-}" ]; then
            print_info "Pre-compaction directory: $WORKING_DIRECTORY"
            
            if [ "$current_dir" != "$WORKING_DIRECTORY" ]; then
                print_warning "Directory mismatch detected!"
                echo "  Current: $current_dir"
                echo "  Expected: $WORKING_DIRECTORY"
                
                if prompt_yes_no "Navigate to expected directory?"; then
                    cd "$WORKING_DIRECTORY" || {
                        print_error "Failed to change directory"
                        print_warning "Continuing with current directory"
                    }
                fi
            fi
        fi
    else
        print_info "No compaction marker found - manual recovery needed"
    fi
}

detect_agent_type() {
    local agent_type="${RECOVERY_STATE[agent_type]}"
    
    if [ -n "$agent_type" ] && [ "$agent_type" != "unknown" ]; then
        print_success "Agent type already detected: $agent_type"
        return 0
    fi
    
    print_info "Detecting agent type from context clues..."
    
    # Try to detect from directory structure
    local current_dir="$(pwd)"
    if [[ "$current_dir" =~ /orchestrator ]]; then
        agent_type="orchestrator"
    elif [[ "$current_dir" =~ /sw-engineer ]]; then
        agent_type="sw-engineer"
    elif [[ "$current_dir" =~ /code-reviewer ]]; then
        agent_type="code-reviewer"
    elif [[ "$current_dir" =~ /architect ]]; then
        agent_type="architect"
    fi
    
    # Try to detect from available files
    if [ -z "$agent_type" ]; then
        if [ -f "./orchestrator-state.yaml" ]; then
            agent_type="orchestrator"
        elif [ -f "./IMPLEMENTATION-PLAN.md" ]; then
            agent_type="sw-engineer"
        elif [ -f "./REVIEW-FEEDBACK.md" ] || [ -f "./SPLIT-SUMMARY.md" ]; then
            agent_type="code-reviewer"
        elif [ -f "./WAVE-REVIEW.md" ] || [ -f "./PHASE-ASSESSMENT.md" ]; then
            agent_type="architect"
        fi
    fi
    
    if [ -n "$agent_type" ]; then
        print_success "Auto-detected agent type: $agent_type"
        RECOVERY_STATE[agent_type]="$agent_type"
    else
        print_warning "Could not auto-detect agent type"
        
        agent_type=$(prompt_choice "Please select your agent type:" \
            "orchestrator (orchestrator)" \
            "sw-engineer (sw-engineer)" \
            "code-reviewer (code-reviewer)" \
            "architect (architect)")
        
        # Extract just the agent name
        agent_type=$(echo "$agent_type" | cut -d' ' -f1)
        RECOVERY_STATE[agent_type]="$agent_type"
    fi
    
    print_info "Confirmed agent type: ${RECOVERY_STATE[agent_type]}"
}

recover_todo_state() {
    local agent_type="${RECOVERY_STATE[agent_type]}"
    
    print_info "Searching for TODO state files for agent: $agent_type"
    
    # Look for TODO files in various locations
    # Use CLAUDE_PROJECT_DIR if set
    local project_dir="${CLAUDE_PROJECT_DIR:-/workspaces/software-factory-2.0-template}"
    local todo_locations=(
        "$FACTORY_ROOT/todos"
        "./todos"
        "../todos"
        "../../todos"
        "${project_dir}/todos"
    )
    
    local found_todos=()
    
    for location in "${todo_locations[@]}"; do
        if [ -d "$location" ]; then
            while IFS= read -r -d '' file; do
                if [[ "$(basename "$file")" =~ ^${agent_type}- ]]; then
                    found_todos+=("$file")
                fi
            done < <(find "$location" -name "${agent_type}-*.todo" -print0 2>/dev/null)
        fi
    done
    
    if [ ${#found_todos[@]} -eq 0 ]; then
        print_warning "No TODO files found for agent: $agent_type"
        
        if prompt_yes_no "Create a new TODO state?"; then
            create_initial_todo_state "$agent_type"
        else
            print_info "Continuing without TODO state recovery"
        fi
        return 0
    fi
    
    # Sort by modification time (newest first)
    IFS=$'\n' found_todos=($(printf '%s\n' "${found_todos[@]}" | while read file; do
        printf '%s %s\n' "$(stat -c %Y "$file")" "$file"
    done | sort -nr | cut -d' ' -f2-))
    
    print_success "Found ${#found_todos[@]} TODO files for $agent_type"
    
    # Show available files
    echo "Available TODO files:"
    for i in "${!found_todos[@]}"; do
        local file="${found_todos[$i]}"
        local mod_time=$(stat -c '%Y' "$file")
        local readable_time=$(date -d "@$mod_time" '+%Y-%m-%d %H:%M:%S')
        printf "  %d) %-40s %s\n" $((i+1)) "$(basename "$file")" "$readable_time"
    done
    
    # Let user choose which file to load
    local choice
    if [ ${#found_todos[@]} -eq 1 ]; then
        choice="1"
        print_info "Using the only available file"
    else
        read -p "Which TODO file should we load? [1-${#found_todos[@]}] (default: 1): " choice
        choice=${choice:-1}
    fi
    
    if [[ "$choice" =~ ^[0-9]+$ ]] && [ "$choice" -ge 1 ] && [ "$choice" -le "${#found_todos[@]}" ]; then
        local selected_file="${found_todos[$((choice-1))]}"
        print_success "Selected: $(basename "$selected_file")"
        
        # Display the file content
        print_instruction "TODO STATE CONTENT:"
        echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
        cat "$selected_file"
        echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
        echo
        
        print_critical "CRITICAL: You MUST now use the TodoWrite tool!"
        print_instruction "Required steps:"
        echo "1. Parse the TODO items from the content above"
        echo "2. Use TodoWrite tool to populate your working TODO list"
        echo "3. Set appropriate status for each item (pending/in_progress/completed)"
        echo "4. Verify all tasks are loaded correctly"
        echo
        
        if prompt_yes_no "Have you completed the TodoWrite loading?"; then
            RECOVERY_STATE[todos_loaded]="true"
            print_success "TODO state recovery marked as complete"
        else
            print_warning "TODO state recovery incomplete - you must complete this manually"
        fi
    else
        print_error "Invalid choice"
        RECOVERY_STATE[todos_loaded]="false"
    fi
}

create_initial_todo_state() {
    local agent_type="$1"
    
    print_info "Creating initial TODO state for $agent_type"
    
    # Use the todo-preservation script if available
    if [ -x "$HOOKS_DIR/todo-preservation.sh" ]; then
        print_info "Using todo-preservation utility..."
        "$HOOKS_DIR/todo-preservation.sh" save "$agent_type" "RECOVERY_INIT"
        print_success "Initial TODO state created"
    else
        print_warning "Todo-preservation utility not available"
        print_instruction "You'll need to create your TODO state manually using TodoWrite"
    fi
}

recover_context_files() {
    local agent_type="${RECOVERY_STATE[agent_type]}"
    
    print_info "Identifying required context files for agent: $agent_type"
    
    case "$agent_type" in
        "orchestrator")
            recover_orchestrator_context
            ;;
        "sw-engineer")
            recover_sw_engineer_context
            ;;
        "code-reviewer")
            recover_code_reviewer_context
            ;;
        "architect")
            recover_architect_context
            ;;
        *)
            print_warning "Unknown agent type for context recovery"
            ;;
    esac
    
    if prompt_yes_no "Have you read all required context files using the Read tool?"; then
        RECOVERY_STATE[context_files_read]="true"
        print_success "Context file recovery marked as complete"
    else
        print_warning "Context file recovery incomplete"
    fi
}

recover_orchestrator_context() {
    print_instruction "ORCHESTRATOR CONTEXT RECOVERY:"
    echo
    echo "MANDATORY files to read with Read tool:"
    echo "  1. orchestrator-state.yaml (if exists)"
    echo "  2. SOFTWARE-FACTORY-STATE-MACHINE.md"
    echo "  3. Current phase-specific plan files"
    echo
    echo "Check these locations:"
    
    local required_files=(
        "./orchestrator-state.yaml"
        "../orchestrator-state.yaml"
        "$FACTORY_ROOT/state-machines/orchestrator.md"
    )
    
    for file in "${required_files[@]}"; do
        if [ -f "$file" ]; then
            print_success "✓ Found: $file"
        else
            print_warning "✗ Missing: $file"
        fi
    done
    
    print_instruction "After reading these files, determine:"
    echo "  - Current wave and phase position"
    echo "  - Any efforts in progress"
    echo "  - Integration branches needed"
    echo "  - Architect reviews pending"
}

recover_sw_engineer_context() {
    print_instruction "SOFTWARE ENGINEER CONTEXT RECOVERY:"
    echo
    echo "MANDATORY files to read with Read tool:"
    echo "  1. IMPLEMENTATION-PLAN.md"
    echo "  2. work-log.md"
    echo "  3. TEST-DRIVEN-VALIDATION-REQUIREMENTS.md"
    echo
    echo "Check these locations:"
    
    local required_files=(
        "./IMPLEMENTATION-PLAN.md"
        "./work-log.md"
        "$FACTORY_ROOT/expertise/testing-strategies.md"
    )
    
    for file in "${required_files[@]}"; do
        if [ -f "$file" ]; then
            print_success "✓ Found: $file"
        else
            print_warning "✗ Missing: $file"
        fi
    done
    
    print_instruction "After reading these files, check:"
    echo "  - Current implementation progress"
    echo "  - Line count compliance (<800 lines)"
    echo "  - Test coverage requirements"
    echo "  - Any pending review feedback"
    
    # Check line count if we have the tool
    if command -v $PROJECT_ROOT/tools/line-counter.sh > /dev/null 2>&1; then
        local current_branch=$(git branch --show-current 2>/dev/null || echo "no-branch")
        if [ "$current_branch" != "no-branch" ]; then
            print_info "Checking current line count..."
            $PROJECT_ROOT/tools/line-counter.sh -c "$current_branch" || true
        fi
    fi
}

recover_code_reviewer_context() {
    print_instruction "CODE REVIEWER CONTEXT RECOVERY:"
    echo
    echo "MANDATORY files to read with Read tool:"
    echo "  1. KCP-CODE-REVIEWER-COMPREHENSIVE-GUIDE.md"
    echo "  2. IMPLEMENTATION-PLAN.md (if exists)"
    echo "  3. REVIEW-FEEDBACK.md (if exists)"
    echo "  4. SPLIT-SUMMARY.md (if exists)"
    echo
    echo "Check these locations:"
    
    local required_files=(
        "./IMPLEMENTATION-PLAN.md"
        "./work-log.md"
        "./REVIEW-FEEDBACK.md"
        "./SPLIT-SUMMARY.md"
        "$FACTORY_ROOT/rule-library/agents/code-reviewer/"
    )
    
    for file in "${required_files[@]}"; do
        if [ -f "$file" ]; then
            print_success "✓ Found: $file"
        elif [ -d "$file" ]; then
            print_success "✓ Found directory: $file"
        else
            print_warning "✗ Missing: $file"
        fi
    done
    
    print_instruction "After reading these files, determine:"
    echo "  - Current review state"
    echo "  - Line count compliance status"
    echo "  - Split planning requirements"
    echo "  - KCP pattern compliance issues"
}

recover_architect_context() {
    print_instruction "ARCHITECT CONTEXT RECOVERY:"
    echo
    echo "MANDATORY files to read with Read tool:"
    echo "  1. orchestrator-state.yaml"
    echo "  2. WAVE-COMPLETION-ARCHITECT-REVIEW-PROTOCOL.md"
    echo "  3. Integration branch status"
    echo
    echo "Check these locations:"
    
    local required_files=(
        "./orchestrator-state.yaml"
        "../orchestrator-state.yaml"
        "./WAVE-REVIEW.md"
        "./PHASE-ASSESSMENT.md"
        "$FACTORY_ROOT/rule-library/agents/architect/"
    )
    
    for file in "${required_files[@]}"; do
        if [ -f "$file" ]; then
            print_success "✓ Found: $file"
        elif [ -d "$file" ]; then
            print_success "✓ Found directory: $file"
        else
            print_warning "✗ Missing: $file"
        fi
    done
    
    print_instruction "After reading these files, assess:"
    echo "  - Wave completion status"
    echo "  - Integration readiness"
    echo "  - Architectural compliance"
    echo "  - Performance implications"
}

validate_recovery() {
    print_info "Validating recovery completeness..."
    
    local validation_passed=true
    
    # Check agent type
    if [ -z "${RECOVERY_STATE[agent_type]}" ]; then
        print_error "Agent type not identified"
        validation_passed=false
    else
        print_success "✓ Agent type identified: ${RECOVERY_STATE[agent_type]}"
    fi
    
    # Check TODO state
    if [ "${RECOVERY_STATE[todos_loaded]}" = "true" ]; then
        print_success "✓ TODO state recovery completed"
    else
        print_warning "⚠ TODO state recovery incomplete"
        if ! prompt_yes_no "Continue without complete TODO recovery?"; then
            print_error "Recovery validation failed - TODO state required"
            validation_passed=false
        fi
    fi
    
    # Check context files
    if [ "${RECOVERY_STATE[context_files_read]}" = "true" ]; then
        print_success "✓ Context files recovery completed"
    else
        print_warning "⚠ Context files recovery incomplete"
        if ! prompt_yes_no "Continue without complete context recovery?"; then
            print_error "Recovery validation failed - context files required"
            validation_passed=false
        fi
    fi
    
    # Environment validation
    validate_environment
    
    if [ "$validation_passed" = true ]; then
        RECOVERY_STATE[validation_passed]="true"
        print_success "🎉 Recovery validation PASSED"
    else
        print_error "Recovery validation FAILED"
        if ! prompt_yes_no "Continue with incomplete recovery?"; then
            exit 1
        fi
    fi
}

validate_environment() {
    print_info "Validating environment setup..."
    
    # Check working directory
    local current_dir="$(pwd)"
    print_info "Current directory: $current_dir"
    
    # Check git status
    if git rev-parse --git-dir > /dev/null 2>&1; then
        local branch=$(git branch --show-current)
        print_success "✓ Git repository: $branch branch"
    else
        print_warning "⚠ Not in a git repository"
    fi
    
    # Check for critical tools
    local tools=(
        "$PROJECT_ROOT/tools/line-counter.sh"
    )
    
    for tool in "${tools[@]}"; do
        if [ -x "$tool" ]; then
            print_success "✓ Tool available: $(basename "$tool")"
        else
            print_warning "⚠ Tool missing: $(basename "$tool")"
        fi
    done
}

provide_continuation_guidance() {
    local agent_type="${RECOVERY_STATE[agent_type]}"
    
    print_header "CONTINUATION GUIDANCE FOR $agent_type"
    
    case "$agent_type" in
        "orchestrator")
            provide_orchestrator_guidance
            ;;
        "sw-engineer")
            provide_sw_engineer_guidance
            ;;
        "code-reviewer")
            provide_code_reviewer_guidance
            ;;
        "architect")
            provide_architect_guidance
            ;;
    esac
    
    print_instruction "NEXT IMMEDIATE ACTIONS:"
    echo "1. Verify your TodoWrite tool contains all recovered tasks"
    echo "2. Review your current state machine position"
    echo "3. Check for any blocking dependencies"
    echo "4. Resume work from where you left off"
    echo "5. Save TODO state periodically to prevent future loss"
}

provide_orchestrator_guidance() {
    echo "As the Orchestrator, your next steps should be:"
    echo
    echo "🎯 IMMEDIATE PRIORITIES:"
    echo "  - Review orchestrator-state.yaml for current wave/phase"
    echo "  - Check efforts_in_progress for any blocking issues"
    echo "  - Identify if any integration branches need creation"
    echo "  - Determine if architect reviews are pending"
    echo
    echo "🔄 STATE MACHINE TRANSITIONS:"
    echo "  - If WAVE_COMPLETE: Create integration branch and spawn architect"
    echo "  - If CHANGES_REQUIRED: Spawn SW engineer for fixes"
    echo "  - If INIT: Continue wave planning and agent spawning"
    echo
    echo "⚠️ CRITICAL REMINDERS:"
    echo "  - NEVER write code yourself - always delegate to SW Engineer"
    echo "  - Ensure every effort gets line count compliance (<800 lines)"
    echo "  - All efforts must be reviewed before wave completion"
}

provide_sw_engineer_guidance() {
    echo "As the Software Engineer, your next steps should be:"
    echo
    echo "🔨 IMMEDIATE PRIORITIES:"
    echo "  - Check your IMPLEMENTATION-PLAN.md for current objectives"
    echo "  - Review work-log.md for progress made so far"
    echo "  - Measure current line count for compliance"
    echo "  - Continue implementation from last checkpoint"
    echo
    echo "📏 LINE COUNT COMPLIANCE:"
    echo "  - Run: $PROJECT_ROOT/tools/line-counter.sh -c \$(git branch --show-current)"
    echo "  - If >800 lines: STOP and request split from code reviewer"
    echo "  - Measure every 200 lines during development"
    echo
    echo "🧪 TEST REQUIREMENTS:"
    echo "  - Follow TDD practices per validation requirements"
    echo "  - Ensure proper test coverage for all new code"
    echo "  - Update work-log.md with progress"
}

provide_code_reviewer_guidance() {
    echo "As the Code Reviewer, your next steps should be:"
    echo
    echo "👁️ IMMEDIATE PRIORITIES:"
    echo "  - Determine current review state (planning/reviewing/splitting)"
    echo "  - Check line count compliance for any pending reviews"
    echo "  - Review KCP patterns and multi-tenancy requirements"
    echo "  - Complete any pending review feedback"
    echo
    echo "📊 REVIEW PROCESS:"
    echo "  - If planning: Create IMPLEMENTATION-PLAN.md and work-log.md"
    echo "  - If reviewing: Check compliance, patterns, tests"
    echo "  - If splitting: Design logical groups <700 lines each"
    echo
    echo "🔧 TOOLS AND VALIDATION:"
    echo "  - Use line-counter.sh for accurate measurements"
    echo "  - Verify test coverage requirements are met"
    echo "  - Ensure KCP architectural patterns compliance"
}

provide_architect_guidance() {
    echo "As the Architect, your next steps should be:"
    echo
    echo "🏗️ IMMEDIATE PRIORITIES:"
    echo "  - Review current wave completion status"
    echo "  - Check all efforts for architectural compliance"
    echo "  - Assess integration readiness and conflicts"
    echo "  - Provide review decision (PROCEED/CHANGES_REQUIRED/STOP)"
    echo
    echo "🔍 REVIEW CRITERIA:"
    echo "  - All efforts are <800 lines compliant"
    echo "  - KCP patterns properly implemented"
    echo "  - Multi-tenancy correctly handled"
    echo "  - No integration conflicts detected"
    echo
    echo "📋 DECISION OUTCOMES:"
    echo "  - PROCEED: Wave is ready for integration"
    echo "  - CHANGES_REQUIRED: Specify needed fixes"
    echo "  - STOP: Critical architectural issues found"
}

complete_recovery() {
    print_header "RECOVERY COMPLETE"
    
    # Save final recovery state
    save_recovery_state
    
    print_success "🎉 Recovery process completed successfully!"
    print_info "Recovery session logged to: $LOG_FILE"
    
    if [ -f "$RECOVERY_STATE_FILE" ]; then
        print_info "Recovery state saved to: $RECOVERY_STATE_FILE"
    fi
    
    echo
    print_instruction "FINAL CHECKLIST:"
    echo "  ✓ Agent type identified: ${RECOVERY_STATE[agent_type]}"
    echo "  ✓ TODO state recovery: ${RECOVERY_STATE[todos_loaded]}"  
    echo "  ✓ Context files read: ${RECOVERY_STATE[context_files_read]}"
    echo "  ✓ Environment validated: ${RECOVERY_STATE[validation_passed]}"
    echo
    
    print_critical "REMEMBER:"
    echo "  - Your TodoWrite tool should contain all recovered tasks"
    echo "  - You've read all required context files"
    echo "  - You know your current state machine position"
    echo "  - Save TODO state regularly to prevent future issues"
    echo
    
    if prompt_yes_no "Clean up recovery files?"; then
        rm -f "$RECOVERY_STATE_FILE" "/tmp/compaction_marker.txt"
        print_info "Recovery files cleaned up"
    fi
    
    print_success "You're ready to resume your Software Factory work! 🚀"
}

################################################################################
# UTILITY FUNCTIONS
################################################################################

save_recovery_state() {
    cat > "$RECOVERY_STATE_FILE" << EOF
# Software Factory 2.0 - Recovery Session State
# Generated: $(date '+%Y-%m-%d %H:%M:%S %Z')

AGENT_TYPE="${RECOVERY_STATE[agent_type]}"
CURRENT_DIRECTORY="${RECOVERY_STATE[current_directory]}"
GIT_BRANCH="${RECOVERY_STATE[git_branch]:-}"
GIT_REPO="${RECOVERY_STATE[git_repo]:-}"
TODOS_LOADED="${RECOVERY_STATE[todos_loaded]}"
CONTEXT_FILES_READ="${RECOVERY_STATE[context_files_read]}"
VALIDATION_PASSED="${RECOVERY_STATE[validation_passed]}"

EOF
}

load_recovery_state() {
    if [ -f "$RECOVERY_STATE_FILE" ]; then
        source "$RECOVERY_STATE_FILE"
        
        RECOVERY_STATE[agent_type]="${AGENT_TYPE:-}"
        RECOVERY_STATE[current_directory]="${CURRENT_DIRECTORY:-$(pwd)}"
        RECOVERY_STATE[git_branch]="${GIT_BRANCH:-}"
        RECOVERY_STATE[git_repo]="${GIT_REPO:-}"
        RECOVERY_STATE[todos_loaded]="${TODOS_LOADED:-false}"
        RECOVERY_STATE[context_files_read]="${CONTEXT_FILES_READ:-false}"
        RECOVERY_STATE[validation_passed]="${VALIDATION_PASSED:-false}"
        
        print_info "Loaded previous recovery session state"
    fi
}

agent_specific_recovery() {
    local agent_type="$1"
    
    print_banner
    print_info "Starting agent-specific recovery for: $agent_type"
    
    RECOVERY_STATE[agent_type]="$agent_type"
    save_recovery_state
    
    print_step 1 "TODO STATE RECOVERY"
    recover_todo_state
    echo
    
    print_step 2 "CONTEXT FILES RECOVERY"
    recover_context_files
    echo
    
    print_step 3 "CONTINUATION GUIDANCE"
    provide_continuation_guidance
    echo
    
    print_success "Agent-specific recovery complete for: $agent_type"
}

validate_current_state() {
    print_banner
    print_info "Validating current recovery state..."
    
    load_recovery_state
    validate_recovery
    
    if [ "${RECOVERY_STATE[validation_passed]}" = "true" ]; then
        print_success "Current state validation PASSED"
        return 0
    else
        print_error "Current state validation FAILED"
        return 1
    fi
}

emergency_recovery() {
    print_banner
    print_critical "EMERGENCY RECOVERY MODE"
    print_warning "Minimal prompts - following automated recovery procedure"
    
    # Quick context detection
    identify_context
    detect_agent_type
    
    # Load TODOs without extensive prompting
    local agent_type="${RECOVERY_STATE[agent_type]}"
    if [ -n "$agent_type" ]; then
        print_info "Emergency TODO recovery for: $agent_type"
        
        # Find latest TODO file
        local latest_todo=$(find "$FACTORY_ROOT/todos" -name "${agent_type}-*.todo" -type f -printf '%T@ %p\n' 2>/dev/null | sort -nr | head -1 | cut -d' ' -f2-)
        
        if [ -n "$latest_todo" ] && [ -f "$latest_todo" ]; then
            print_success "Emergency TODO file found: $(basename "$latest_todo")"
            print_instruction "TODO CONTENT:"
            cat "$latest_todo"
            print_critical "Use TodoWrite tool to load these items immediately!"
        else
            print_error "No TODO files found for emergency recovery"
        fi
    fi
    
    # Provide minimal continuation guidance
    provide_continuation_guidance
    
    print_success "Emergency recovery procedure complete"
    print_warning "Verify your state and continue work carefully"
}

################################################################################
# MAIN COMMAND DISPATCH
################################################################################

main() {
    # Load any existing recovery state
    load_recovery_state
    
    # Parse command line arguments
    case "${1:-}" in
        "--agent")
            if [ $# -lt 2 ]; then
                print_error "Agent type required"
                usage
                exit 1
            fi
            agent_specific_recovery "$2"
            ;;
        "--validate")
            validate_current_state
            ;;
        "--emergency")
            emergency_recovery
            ;;
        "--help"|"-h"|"help")
            usage
            ;;
        "")
            # Interactive mode
            interactive_recovery
            ;;
        *)
            print_error "Unknown option: $1"
            usage
            exit 1
            ;;
    esac
}

################################################################################
# SCRIPT EXECUTION
################################################################################

# Ensure we're being called in the right context
if [ "${BASH_SOURCE[0]}" = "${0}" ]; then
    main "$@"
fi