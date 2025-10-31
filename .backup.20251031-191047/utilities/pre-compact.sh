#!/bin/bash
################################################################################
# SOFTWARE FACTORY 2.0 - PRE-COMPACTION HOOK
################################################################################
#
# This script runs before Claude Code context compaction to preserve critical
# state information and TODO lists that must survive memory compression.
#
# WHEN THIS RUNS:
# - Automatically before Claude Code compaction (manual or automatic)
# - When context window approaches limits
# - Before agent transitions that may trigger compaction
#
# WHAT THIS PRESERVES:
# - Current TODO state from all agents
# - State machine position and context  
# - Agent identity and active work
# - Critical file checkpoints
# - Integration branch status
#
################################################################################

set -euo pipefail

# Constants
readonly SCRIPT_NAME="$(basename "$0")"
readonly HOOKS_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
readonly MARKER_FILE="/tmp/compaction_marker.txt"
readonly TODOS_BACKUP="/tmp/todos-precompact.txt"
readonly STATE_BACKUP="/tmp/state-precompact.json"
readonly LOG_FILE="/tmp/pre-compact.log"

# Logging function
log() {
    local level="$1"
    shift
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] [$level] $*" | tee -a "$LOG_FILE"
}

# Error handling
error_exit() {
    log "ERROR" "$1"
    exit 1
}

# Cleanup function
cleanup() {
    local exit_code=$?
    if [ $exit_code -ne 0 ]; then
        log "ERROR" "Script failed with exit code $exit_code"
    fi
    return $exit_code
}

trap cleanup EXIT

################################################################################
# MAIN PRE-COMPACTION LOGIC
################################################################################

main() {
    log "INFO" "🚨 PRE-COMPACTION HOOK STARTING 🚨"
    log "INFO" "Context compaction detected - preserving critical state..."
    
    # Initialize marker file with header
    cat > "$MARKER_FILE" << 'EOF'
################################################################################
# SOFTWARE FACTORY 2.0 - COMPACTION RECOVERY MARKER
################################################################################
# This file indicates that context compaction occurred and critical state
# has been preserved for recovery by the post-compaction hook.
################################################################################

COMPACTION_TIMESTAMP="$(date '+%Y-%m-%d %H:%M:%S %Z')"
COMPACTION_TRIGGER="pre-compact.sh"
FACTORY_VERSION="2.0"

EOF
    
    # Detect and preserve current working directory
    preserve_working_directory
    
    # Detect and preserve git context
    preserve_git_context
    
    # Detect agent identity from current context
    detect_agent_identity
    
    # Preserve TODO state from todos directory
    preserve_todo_state
    
    # Preserve state machine context
    preserve_state_machine_context
    
    # Create complete state snapshot
    create_state_snapshot
    
    # Finalize marker file
    finalize_marker_file
    
    log "INFO" "✅ PRE-COMPACTION PRESERVATION COMPLETE"
    log "INFO" "Recovery marker created: $MARKER_FILE"
    log "INFO" "TODO backup created: $TODOS_BACKUP"
    log "INFO" "State backup created: $STATE_BACKUP"
    
    return 0
}

################################################################################
# PRESERVATION FUNCTIONS
################################################################################

preserve_working_directory() {
    log "INFO" "Preserving working directory context..."
    
    local current_dir="$(pwd)"
    echo "WORKING_DIRECTORY=\"$current_dir\"" >> "$MARKER_FILE"
    
    # Detect if we're in a specific project context
    if [[ "$current_dir" =~ /agent-states/ ]]; then
        local agent_state=$(basename "$current_dir")
        echo "AGENT_STATE_CONTEXT=\"$agent_state\"" >> "$MARKER_FILE"
        log "INFO" "Agent state context detected: $agent_state"
    fi
    
    # Check for project root indicators
    if [ -f "$current_dir/package.json" ] || [ -f "$current_dir/go.mod" ] || [ -f "$current_dir/Cargo.toml" ]; then
        echo "PROJECT_ROOT=\"$current_dir\"" >> "$MARKER_FILE"
        log "INFO" "Project root detected: $current_dir"
    fi
}

preserve_git_context() {
    log "INFO" "Preserving git context..."
    
    if git rev-parse --git-dir > /dev/null 2>&1; then
        local branch=$(git branch --show-current 2>/dev/null || echo "detached")
        local commit=$(git rev-parse HEAD 2>/dev/null || echo "unknown")
        local status=$(git status --porcelain 2>/dev/null | wc -l)
        local remote_url=$(git remote get-url origin 2>/dev/null || echo "no-remote")
        
        echo "GIT_BRANCH=\"$branch\"" >> "$MARKER_FILE"
        echo "GIT_COMMIT=\"$commit\"" >> "$MARKER_FILE"
        echo "GIT_DIRTY_FILES=\"$status\"" >> "$MARKER_FILE"
        echo "GIT_REMOTE_URL=\"$remote_url\"" >> "$MARKER_FILE"
        
        log "INFO" "Git context: branch=$branch, commit=${commit:0:8}, dirty_files=$status"
    else
        echo "GIT_CONTEXT=\"none\"" >> "$MARKER_FILE"
        log "WARN" "No git repository detected"
    fi
}

detect_agent_identity() {
    log "INFO" "Detecting agent identity..."
    
    # Try to detect agent from various context clues
    local agent_type="unknown"
    local current_dir="$(pwd)"
    
    # Check current directory for agent hints
    if [[ "$current_dir" =~ /orchestrator/ ]]; then
        agent_type="orchestrator"
    elif [[ "$current_dir" =~ /sw-engineer/ ]]; then
        agent_type="sw-engineer"
    elif [[ "$current_dir" =~ /code-reviewer/ ]]; then
        agent_type="code-reviewer"
    elif [[ "$current_dir" =~ /architect/ ]]; then
        agent_type="architect"
    fi
    
    # Check for agent-specific files
    if [ -f "./orchestrator-state-v3.json" ]; then
        agent_type="orchestrator"
    elif [ -f "./IMPLEMENTATION-PLAN.md" ]; then
        agent_type="sw-engineer"
    elif [ -f "./REVIEW-FEEDBACK.md" ]; then
        agent_type="code-reviewer"
    fi
    
    echo "DETECTED_AGENT=\"$agent_type\"" >> "$MARKER_FILE"
    log "INFO" "Agent identity detected: $agent_type"
    
    # Save agent-specific context
    case $agent_type in
        "orchestrator")
            preserve_orchestrator_context
            ;;
        "sw-engineer")
            preserve_sw_engineer_context
            ;;
        "code-reviewer")
            preserve_code_reviewer_context
            ;;
        "architect")
            preserve_architect_context
            ;;
        *)
            log "WARN" "Unknown agent type, using generic preservation"
            ;;
    esac
}

preserve_todo_state() {
    log "INFO" "Preserving TODO state..."
    
    # Look for todos directory in various locations
    # Use CLAUDE_PROJECT_DIR if set, otherwise check relative paths
    local project_dir="${CLAUDE_PROJECT_DIR:-/workspaces/software-factory-2.0-template}"
    local todos_dirs=(
        "./todos"
        "../todos" 
        "../../todos"
        "${project_dir}/todos"
    )
    
    local found_todos=false
    
    for todos_dir in "${todos_dirs[@]}"; do
        if [ -d "$todos_dir" ]; then
            log "INFO" "Found todos directory: $todos_dir"
            
            # Find latest TODO file by modification time
            local latest_todo=$(find "$todos_dir" -name "*.todo" -type f -printf '%T@ %p\n' 2>/dev/null | sort -nr | head -1 | cut -d' ' -f2-)
            
            if [ -n "$latest_todo" ] && [ -f "$latest_todo" ]; then
                log "INFO" "Latest TODO file: $latest_todo"
                
                # Copy to backup location
                cp "$latest_todo" "$TODOS_BACKUP"
                echo "TODO_STATE_SAVED=\"$latest_todo\"" >> "$MARKER_FILE"
                echo "TODO_BACKUP_PATH=\"$TODOS_BACKUP\"" >> "$MARKER_FILE"
                
                # Count TODO files
                local todo_count=$(find "$todos_dir" -name "*.todo" -type f | wc -l)
                echo "TOTAL_TODO_FILES=\"$todo_count\"" >> "$MARKER_FILE"
                
                found_todos=true
                break
            fi
        fi
    done
    
    if [ "$found_todos" = false ]; then
        echo "TODO_STATE_SAVED=\"none\"" >> "$MARKER_FILE"
        log "WARN" "No TODO files found in any expected location"
    fi
}

preserve_state_machine_context() {
    log "INFO" "Preserving state machine context..."
    
    # Look for state machine indicators
    local state_files=(
        "./orchestrator-state-v3.json"
        "../orchestrator-state-v3.json"
        "./state.json"
        "./current-state.md"
    )
    
    for state_file in "${state_files[@]}"; do
        if [ -f "$state_file" ]; then
            log "INFO" "Found state file: $state_file"
            
            # Extract current phase/wave if possible
            if [[ "$state_file" =~ \.yaml$ ]]; then
                if command -v yq > /dev/null 2>&1; then
                    local phase=$(jq '.current_phase // "unknown"' "$state_file" 2>/dev/null || echo "unknown")
                    local wave=$(jq '.current_wave // "unknown"' "$state_file" 2>/dev/null || echo "unknown")
                    echo "STATE_MACHINE_PHASE=\"$phase\"" >> "$MARKER_FILE"
                    echo "STATE_MACHINE_WAVE=\"$wave\"" >> "$MARKER_FILE"
                fi
            fi
            
            echo "STATE_MACHINE_FILE=\"$state_file\"" >> "$MARKER_FILE"
            break
        fi
    done
    
    # Check for specific state directories
    local current_dir="$(basename "$(pwd)")"
    if [[ "$current_dir" =~ ^[A-Z_]+$ ]]; then
        echo "STATE_MACHINE_STATE=\"$current_dir\"" >> "$MARKER_FILE"
        log "INFO" "State machine state detected from directory: $current_dir"
    fi
}

preserve_orchestrator_context() {
    log "INFO" "Preserving orchestrator-specific context..."
    
    if [ -f "./orchestrator-state-v3.json" ]; then
        echo "ORCHESTRATOR_STATE_FILE=\"./orchestrator-state-v3.json\"" >> "$MARKER_FILE"
    fi
    
    # Check for wave/phase context
    if [ -d "./efforts_in_progress" ]; then
        local effort_count=$(find ./efforts_in_progress -type d | wc -l)
        echo "EFFORTS_IN_PROGRESS=\"$effort_count\"" >> "$MARKER_FILE"
    fi
}

preserve_sw_engineer_context() {
    log "INFO" "Preserving software engineer-specific context..."
    
    if [ -f "./IMPLEMENTATION-PLAN.md" ]; then
        echo "SW_ENG_PLAN_FILE=\"./IMPLEMENTATION-PLAN.md\"" >> "$MARKER_FILE"
    fi
    
    if [ -f "./work-log.md" ]; then
        echo "SW_ENG_WORKLOG=\"./work-log.md\"" >> "$MARKER_FILE"
        
        # Extract last work entry timestamp if possible
        local last_entry=$(tail -5 "./work-log.md" | grep -E '^\[.*\]' | tail -1 || echo "")
        if [ -n "$last_entry" ]; then
            echo "SW_ENG_LAST_WORK=\"$last_entry\"" >> "$MARKER_FILE"
        fi
    fi
}

preserve_code_reviewer_context() {
    log "INFO" "Preserving code reviewer-specific context..."
    
    if [ -f "./REVIEW-FEEDBACK.md" ]; then
        echo "CODE_REVIEWER_FEEDBACK=\"./REVIEW-FEEDBACK.md\"" >> "$MARKER_FILE"
    fi
    
    if [ -f "./SPLIT-SUMMARY.md" ]; then
        echo "CODE_REVIEWER_SPLIT_PLAN=\"./SPLIT-SUMMARY.md\"" >> "$MARKER_FILE"
    fi
}

preserve_architect_context() {
    log "INFO" "Preserving architect-specific context..."
    
    # Look for architect review files
    if [ -f "./WAVE-REVIEW.md" ]; then
        echo "ARCHITECT_REVIEW_WAVE_ARCHITECTURE=\"./WAVE-REVIEW.md\"" >> "$MARKER_FILE"
    fi
    
    if [ -f "./PHASE-ASSESSMENT.md" ]; then
        echo "ARCHITECT_PHASE_ASSESSMENT=\"./PHASE-ASSESSMENT.md\"" >> "$MARKER_FILE"
    fi
}

create_state_snapshot() {
    log "INFO" "Creating complete state snapshot..."
    
    # Create JSON snapshot with all preserved data
    cat > "$STATE_BACKUP" << EOF
{
  "timestamp": "$(date '+%Y-%m-%d %H:%M:%S %Z')",
  "factory_version": "2.0",
  "working_directory": "$(pwd)",
  "git_branch": "$(git branch --show-current 2>/dev/null || echo 'unknown')",
  "git_commit": "$(git rev-parse HEAD 2>/dev/null || echo 'unknown')",
  "detected_agent": "$(grep 'DETECTED_AGENT=' "$MARKER_FILE" | cut -d'"' -f2)",
  "environment": {
    "hostname": "$(hostname)",
    "user": "$(whoami)",
    "shell": "$SHELL",
    "path": "$PATH"
  },
  "files_present": [
EOF
    
    # Add important files that exist
    local important_files=(
        "./orchestrator-state-v3.json"
        "./IMPLEMENTATION-PLAN.md"
        "./work-log.md"
        "./REVIEW-FEEDBACK.md"
        "./SPLIT-SUMMARY.md"
        "./package.json"
        "./go.mod"
        "./Makefile"
    )
    
    local first=true
    for file in "${important_files[@]}"; do
        if [ -f "$file" ]; then
            if [ "$first" = true ]; then
                first=false
            else
                echo "," >> "$STATE_BACKUP"
            fi
            echo "    \"$file\"" >> "$STATE_BACKUP"
        fi
    done
    
    echo "" >> "$STATE_BACKUP"
    echo "  ]" >> "$STATE_BACKUP"
    echo "}" >> "$STATE_BACKUP"
    
    echo "STATE_SNAPSHOT=\"$STATE_BACKUP\"" >> "$MARKER_FILE"
}

finalize_marker_file() {
    log "INFO" "Finalizing recovery marker..."
    
    cat >> "$MARKER_FILE" << 'EOF'

################################################################################
# RECOVERY INSTRUCTIONS
################################################################################
# 1. Run post-compact.sh to get recovery assistance
# 2. Check TODO_BACKUP_PATH for preserved TODO state
# 3. Use detected agent type to determine required context files
# 4. Load TODO state using TodoWrite tool (CRITICAL - not just read!)
# 5. Resume work from preserved state machine position
################################################################################

RECOVERY_UTILITY="/workspaces/software-factory-2.0-template/utilities/post-compact.sh"
RECOVERY_ASSISTANT="/workspaces/software-factory-2.0-template/utilities/recovery-assistant.sh"

EOF
    
    echo "MARKER_COMPLETE=\"true\"" >> "$MARKER_FILE"
}

################################################################################
# SCRIPT EXECUTION
################################################################################

# Ensure we're being called in the right context
if [ "${BASH_SOURCE[0]}" = "${0}" ]; then
    main "$@"
fi