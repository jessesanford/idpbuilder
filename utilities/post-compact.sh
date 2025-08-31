#!/bin/bash
################################################################################
# SOFTWARE FACTORY 2.0 - POST-COMPACTION RECOVERY HOOK  
################################################################################
#
# This script runs after Claude Code context compaction to help recover
# critical state information and guide the user through proper recovery steps.
#
# WHEN THIS RUNS:
# - After Claude Code context compaction 
# - When compaction marker is detected
# - Can be run manually for recovery assistance
#
# WHAT THIS DOES:
# - Detects compaction markers and preserved state
# - Provides agent-specific recovery instructions
# - Helps load TODO state properly
# - Guides through context restoration
# - Cleans up temporary files
#
################################################################################

set -euo pipefail

# Constants
readonly SCRIPT_NAME="$(basename "$0")"
readonly HOOKS_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
readonly MARKER_FILE="/tmp/compaction_marker.txt"
readonly TODOS_BACKUP="/tmp/todos-precompact.txt" 
readonly STATE_BACKUP="/tmp/state-precompact.json"
readonly LOG_FILE="/tmp/post-compact.log"

# Colors for output
readonly RED='\033[0;31m'
readonly GREEN='\033[0;32m'
readonly YELLOW='\033[1;33m'
readonly BLUE='\033[0;34m'
readonly MAGENTA='\033[0;35m'
readonly CYAN='\033[0;36m'
readonly NC='\033[0m' # No Color

# Logging function
log() {
    local level="$1"
    shift
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] [$level] $*" | tee -a "$LOG_FILE"
}

# Colored output functions
print_header() {
    echo -e "${CYAN}################################################################################${NC}"
    echo -e "${CYAN}# $1${NC}"
    echo -e "${CYAN}################################################################################${NC}"
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
    echo -e "${RED}🚨 $1${NC}"
}

################################################################################
# MAIN RECOVERY LOGIC
################################################################################

main() {
    print_header "SOFTWARE FACTORY 2.0 - POST-COMPACTION RECOVERY"
    log "INFO" "Post-compaction recovery hook starting..."
    
    # Check if compaction actually occurred
    if ! check_compaction_occurred; then
        print_info "No compaction detected. Nothing to recover."
        return 0
    fi
    
    # Load and display preserved state
    load_preserved_state
    
    # Provide agent-specific recovery guidance
    provide_recovery_guidance
    
    # Help with TODO state recovery
    assist_todo_recovery
    
    # Provide final instructions
    provide_final_instructions
    
    # Clean up (with user confirmation)
    cleanup_recovery_files
    
    print_success "Post-compaction recovery assistance complete"
    return 0
}

################################################################################
# DETECTION AND LOADING FUNCTIONS
################################################################################

check_compaction_occurred() {
    if [ ! -f "$MARKER_FILE" ]; then
        return 1
    fi
    
    print_success "Compaction marker detected: $MARKER_FILE"
    return 0
}

load_preserved_state() {
    print_header "PRESERVED STATE ANALYSIS"
    
    if [ -f "$MARKER_FILE" ]; then
        log "INFO" "Loading preserved state from marker file..."
        
        # Source the marker file to get variables
        set +u  # Allow undefined variables temporarily
        source "$MARKER_FILE"
        set -u
        
        # Display key information
        print_info "Compaction Timestamp: ${COMPACTION_TIMESTAMP:-unknown}"
        print_info "Factory Version: ${FACTORY_VERSION:-unknown}"
        print_info "Working Directory: ${WORKING_DIRECTORY:-unknown}"
        
        if [ -n "${GIT_BRANCH:-}" ]; then
            print_info "Git Branch: $GIT_BRANCH"
            print_info "Git Commit: ${GIT_COMMIT:0:8}..."
            
            if [ "${GIT_DIRTY_FILES:-0}" -gt 0 ]; then
                print_warning "Git repository has $GIT_DIRTY_FILES uncommitted changes"
            fi
        fi
        
        if [ -n "${DETECTED_AGENT:-}" ] && [ "$DETECTED_AGENT" != "unknown" ]; then
            print_success "Detected Agent Type: $DETECTED_AGENT"
        else
            print_warning "Agent type could not be determined automatically"
        fi
        
        # Display state machine context
        if [ -n "${STATE_MACHINE_STATE:-}" ]; then
            print_info "State Machine State: $STATE_MACHINE_STATE"
        fi
        
        if [ -n "${STATE_MACHINE_PHASE:-}" ] && [ "$STATE_MACHINE_PHASE" != "unknown" ]; then
            print_info "Phase: $STATE_MACHINE_PHASE"
        fi
        
        if [ -n "${STATE_MACHINE_WAVE:-}" ] && [ "$STATE_MACHINE_WAVE" != "unknown" ]; then
            print_info "Wave: $STATE_MACHINE_WAVE"
        fi
        
    else
        print_error "Marker file not found - unable to load preserved state"
        return 1
    fi
}

provide_recovery_guidance() {
    print_header "RECOVERY GUIDANCE"
    
    local agent_type="${DETECTED_AGENT:-unknown}"
    
    print_critical "CRITICAL RECOVERY STEPS FOR AGENT: $agent_type"
    echo
    
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
        *)
            provide_generic_guidance
            ;;
    esac
}

provide_orchestrator_guidance() {
    echo -e "${YELLOW}🎯 ORCHESTRATOR RECOVERY PROTOCOL:${NC}"
    echo
    echo "1. 📋 RESTORE TODO STATE (CRITICAL):"
    echo "   - Use Read tool on your latest orchestrator-*.todo file"
    echo "   - Use TodoWrite tool to populate your working TODO list"
    echo "   - DO NOT just read - you must LOAD into TodoWrite!"
    echo
    echo "2. 📊 READ STATE FILES:"
    echo "   - orchestrator-state.yaml (current wave/phase status)"
    echo "   - SOFTWARE-FACTORY-STATE-MACHINE.md (state transitions)"
    echo "   - Current phase/wave specific plan files"
    echo
    echo "3. 🔍 CHECK CURRENT POSITION:"
    echo "   - Review efforts_in_progress for blocking issues"
    echo "   - Check integration_branches for pending work"
    echo "   - Determine if wave is complete and needs integration"
    echo
    echo "4. 🎬 RESUME ORCHESTRATION:"
    echo "   - If WAVE_COMPLETE: Create integration branch first"
    echo "   - If CHANGES_REQUIRED: Spawn SW engineer for fixes"
    echo "   - If INIT: Continue from current wave planning"
    echo
    
    if [ -n "${ORCHESTRATOR_STATE_FILE:-}" ]; then
        print_info "State file available: $ORCHESTRATOR_STATE_FILE"
    fi
    
    if [ -n "${EFFORTS_IN_PROGRESS:-}" ]; then
        print_warning "Efforts in progress detected: $EFFORTS_IN_PROGRESS"
    fi
}

provide_sw_engineer_guidance() {
    echo -e "${YELLOW}⚙️ SOFTWARE ENGINEER RECOVERY PROTOCOL:${NC}"
    echo
    echo "1. 📋 RESTORE TODO STATE (CRITICAL):"
    echo "   - Use Read tool on your latest sw-eng-*.todo file"  
    echo "   - Use TodoWrite tool to populate your working TODO list"
    echo "   - Verify all implementation tasks are loaded"
    echo
    echo "2. 📖 READ CONTEXT FILES:"
    echo "   - IMPLEMENTATION-PLAN.md (what you're building)"
    echo "   - work-log.md (progress tracking)"
    echo "   - TEST-DRIVEN-VALIDATION-REQUIREMENTS.md (coverage rules)"
    echo
    echo "3. 📏 CHECK LINE COUNT COMPLIANCE:"
    echo "   - Run: /workspaces/kcp-shared-tools/tmc-pr-line-counter.sh -c \$(git branch --show-current)"
    echo "   - If >800 lines: STOP and request split"
    echo "   - If <800 lines: Continue implementation"
    echo
    echo "4. 🔨 RESUME IMPLEMENTATION:"
    echo "   - Update work-log.md as you progress"
    echo "   - Follow TDD practices per validation requirements"
    echo "   - Measure every 200 lines added"
    echo
    
    if [ -n "${SW_ENG_PLAN_FILE:-}" ]; then
        print_info "Implementation plan available: $SW_ENG_PLAN_FILE"
    fi
    
    if [ -n "${SW_ENG_WORKLOG:-}" ]; then
        print_info "Work log available: $SW_ENG_WORKLOG"
        if [ -n "${SW_ENG_LAST_WORK:-}" ]; then
            print_info "Last work entry: $SW_ENG_LAST_WORK"
        fi
    fi
}

provide_code_reviewer_guidance() {
    echo -e "${YELLOW}👁️ CODE REVIEWER RECOVERY PROTOCOL:${NC}"
    echo
    echo "1. 📋 RESTORE TODO STATE (CRITICAL):"
    echo "   - Use Read tool on your latest code-reviewer-*.todo file"
    echo "   - Use TodoWrite tool to populate your working TODO list"
    echo "   - Check for pending reviews or split planning tasks"
    echo
    echo "2. 📖 READ CONTEXT FILES:"
    echo "   - KCP-CODE-REVIEWER-COMPREHENSIVE-GUIDE.md (review criteria)"
    echo "   - IMPLEMENTATION-PLAN.md (what's being built)"
    echo "   - work-log.md (implementation progress)"
    echo
    echo "3. 📏 CHECK LINE COUNT STATUS:"
    echo "   - Run: /workspaces/kcp-shared-tools/tmc-pr-line-counter.sh -c \$(git branch --show-current)"
    echo "   - If >800 lines: Create split plan immediately"
    echo "   - If <800 lines: Complete review process"
    echo
    echo "4. 🔍 RESUME REVIEW WORK:"
    echo "   - If reviewing: Check KCP patterns, tests, compliance"
    echo "   - If splitting: Design logical groupings <700 lines each"
    echo "   - If planning: Create IMPLEMENTATION-PLAN.md and work-log.md"
    echo
    
    if [ -n "${CODE_REVIEWER_FEEDBACK:-}" ]; then
        print_info "Review feedback available: $CODE_REVIEWER_FEEDBACK"
    fi
    
    if [ -n "${CODE_REVIEWER_SPLIT_PLAN:-}" ]; then
        print_info "Split plan available: $CODE_REVIEWER_SPLIT_PLAN"
    fi
}

provide_architect_guidance() {
    echo -e "${YELLOW}🏗️ ARCHITECT RECOVERY PROTOCOL:${NC}"
    echo
    echo "1. 📋 RESTORE TODO STATE (CRITICAL):"
    echo "   - Use Read tool on your latest architect-*.todo file"
    echo "   - Use TodoWrite tool to populate your working TODO list"
    echo "   - Check for pending wave or phase reviews"
    echo
    echo "2. 📖 READ CONTEXT FILES:"
    echo "   - orchestrator-state.yaml (current wave status)"
    echo "   - WAVE-COMPLETION-ARCHITECT-REVIEW-PROTOCOL.md"
    echo "   - Integration branch status and conflicts"
    echo
    echo "3. 🔍 ASSESS CURRENT REVIEW STATE:"
    echo "   - Are all efforts in wave complete?"
    echo "   - Are all splits <800 lines compliant?"
    echo "   - Are there integration conflicts?"
    echo
    echo "4. 📋 PROVIDE REVIEW DECISION:"
    echo "   - PROCEED: All good, continue to next wave"
    echo "   - CHANGES_REQUIRED: Specify what needs fixing"
    echo "   - STOP: Critical issues that block progress"
    echo
    
    if [ -n "${ARCHITECT_WAVE_REVIEW:-}" ]; then
        print_info "Wave review available: $ARCHITECT_WAVE_REVIEW"
    fi
    
    if [ -n "${ARCHITECT_PHASE_ASSESSMENT:-}" ]; then
        print_info "Phase assessment available: $ARCHITECT_PHASE_ASSESSMENT"
    fi
}

provide_generic_guidance() {
    echo -e "${YELLOW}🔧 GENERIC RECOVERY PROTOCOL:${NC}"
    echo
    echo "1. 🆔 IDENTIFY YOUR AGENT ROLE:"
    echo "   - Check your current prompt for @agent-* references"
    echo "   - Look at your working directory context"
    echo "   - Review what files are available"
    echo
    echo "2. 📋 RESTORE TODO STATE:"
    echo "   - Find your agent's latest .todo file in todos directory"
    echo "   - Use Read tool to examine the content"
    echo "   - Use TodoWrite tool to load tasks (not just read!)"
    echo
    echo "3. 📖 READ REQUIRED CONTEXT:"
    echo "   - Check SOFTWARE-FACTORY-STATE-MACHINE.md for your role"
    echo "   - Read agent-specific instruction files"
    echo "   - Load current work context files"
    echo
    echo "4. 🎯 RESUME WORK:"
    echo "   - Follow your agent's state machine protocol"
    echo "   - Check for any blocking issues"
    echo "   - Continue from preserved state position"
}

assist_todo_recovery() {
    print_header "TODO STATE RECOVERY ASSISTANCE"
    
    if [ -n "${TODO_STATE_SAVED:-}" ] && [ "$TODO_STATE_SAVED" != "none" ]; then
        print_success "TODO state was preserved: $TODO_STATE_SAVED"
        
        if [ -f "$TODOS_BACKUP" ]; then
            print_info "TODO backup available at: $TODOS_BACKUP"
            echo
            print_critical "MANDATORY TODO RECOVERY STEPS:"
            echo "1. 📖 Use Read tool on: $TODO_STATE_SAVED"
            echo "2. 📝 Use TodoWrite tool to populate your working TODO list"
            echo "3. ✅ Verify TodoWrite now contains all recovered TODOs"
            echo "4. 🔄 Deduplicate any duplicate tasks"
            echo
            print_warning "DO NOT just read the file - you MUST load into TodoWrite!"
        else
            print_warning "TODO backup file not found, but original may still exist"
        fi
        
        # Show todo files that exist
        echo
        print_info "Available TODO files in todos directory:"
        
        # Use CLAUDE_PROJECT_DIR if set
        local project_dir="${CLAUDE_PROJECT_DIR:-/workspaces/software-factory-2.0-template}"
        local todos_dirs=(
            "./todos"
            "../todos" 
            "../../todos"
            "${project_dir}/todos"
        )
        
        for todos_dir in "${todos_dirs[@]}"; do
            if [ -d "$todos_dir" ]; then
                local todo_files=$(find "$todos_dir" -name "*.todo" -type f 2>/dev/null | head -5)
                if [ -n "$todo_files" ]; then
                    echo "$todo_files" | while read -r file; do
                        local mod_time=$(stat -c '%Y' "$file" 2>/dev/null || echo "0")
                        local readable_time=$(date -d "@$mod_time" '+%Y-%m-%d %H:%M:%S' 2>/dev/null || echo "unknown")
                        print_info "  - $(basename "$file") (modified: $readable_time)"
                    done
                    break
                fi
            fi
        done
        
    else
        print_warning "No TODO state was preserved during compaction"
        print_info "You may need to recreate your TODO list from context"
    fi
}

provide_final_instructions() {
    print_header "FINAL RECOVERY INSTRUCTIONS"
    
    print_critical "⚠️⚠️⚠️ CRITICAL NEXT STEPS ⚠️⚠️⚠️"
    echo
    echo "Now you MUST:"
    echo "1. 🆔 Confirm your agent identity"
    echo "2. 📖 READ your TODO file using Read tool"
    echo "3. 📝 LOAD TODOs using TodoWrite tool (CRITICAL!)"
    echo "4. 🔄 DEDUPLICATE any overlapping tasks"
    echo "5. ✅ VERIFY TodoWrite contains all recovered TODOs"
    echo "6. 📚 Read your required context files per agent guidance above"
    echo "7. 🎯 Resume work from your preserved state machine position"
    echo
    
    # Show environment verification
    print_info "ENVIRONMENT VERIFICATION:"
    echo "   Current Directory: $(pwd)"
    echo "   Git Branch: $(git branch --show-current 2>/dev/null || echo 'not a git repo')"
    echo "   Working Directory Expected: ${WORKING_DIRECTORY:-unknown}"
    echo "   Git Branch Expected: ${GIT_BRANCH:-unknown}"
    echo
    
    local current_dir="$(pwd)"
    local expected_dir="${WORKING_DIRECTORY:-}"
    
    if [ -n "$expected_dir" ] && [ "$current_dir" != "$expected_dir" ]; then
        print_error "⚠️ DIRECTORY MISMATCH!"
        print_error "Current: $current_dir"
        print_error "Expected: $expected_dir"
        print_critical "You may need to navigate to correct directory"
    else
        print_success "✅ Working directory matches preserved state"
    fi
}

cleanup_recovery_files() {
    print_header "CLEANUP RECOVERY FILES"
    
    echo -e "${YELLOW}The following recovery files can be cleaned up:${NC}"
    echo "  - $MARKER_FILE"
    echo "  - $TODOS_BACKUP (if exists)"
    echo "  - $STATE_BACKUP (if exists)"
    echo "  - $LOG_FILE"
    echo
    
    read -p "Clean up recovery files? (y/N): " -r
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        rm -f "$MARKER_FILE" "$TODOS_BACKUP" "$STATE_BACKUP" 
        print_success "Recovery files cleaned up"
    else
        print_info "Recovery files preserved for manual inspection"
        print_info "You can clean them up later by running: rm -f /tmp/*compact*"
    fi
}

################################################################################
# SCRIPT EXECUTION
################################################################################

# Ensure we're being called in the right context
if [ "${BASH_SOURCE[0]}" = "${0}" ]; then
    main "$@"
fi