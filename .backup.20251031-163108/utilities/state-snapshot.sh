#!/bin/bash
################################################################################
# SOFTWARE FACTORY 2.0 - COMPLETE STATE SNAPSHOT UTILITY
################################################################################
#
# This utility creates comprehensive snapshots of the entire Software Factory
# state for backup, recovery, and analysis purposes.
#
# FEATURES:
# - Complete agent state preservation
# - Git repository state capture
# - File system structure analysis
# - Checkpoint and rule state backup
# - Integration branch tracking
# - Performance metrics collection
#
# USAGE:
#   ./state-snapshot.sh create [snapshot-name]
#   ./state-snapshot.sh restore <snapshot-name>
#   ./state-snapshot.sh list
#   ./state-snapshot.sh compare <snapshot1> <snapshot2>
#   ./state-snapshot.sh cleanup [days-old]
#
################################################################################

set -euo pipefail

# Constants
readonly SCRIPT_NAME="$(basename "$0")"
readonly HOOKS_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
readonly FACTORY_ROOT="/workspaces/software-factory-2.0-template"
readonly SNAPSHOTS_DIR="$FACTORY_ROOT/snapshots"
readonly BACKUP_DIR="/tmp/factory-snapshots"
readonly LOG_FILE="/tmp/state-snapshot.log"

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

# Output functions
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

# Error handling
error_exit() {
    print_error "$1"
    exit 1
}

################################################################################
# USAGE AND HELP
################################################################################

usage() {
    cat << EOF
SOFTWARE FACTORY 2.0 - STATE SNAPSHOT UTILITY

USAGE:
    $SCRIPT_NAME <command> [arguments]

COMMANDS:
    create [name]           Create a complete state snapshot
    restore <name>          Restore from a state snapshot (WIP)
    list                    List available snapshots
    compare <s1> <s2>       Compare two snapshots
    cleanup [days]          Clean old snapshots (default: 7 days)
    analyze <name>          Analyze snapshot contents
    help                    Show this help message

SNAPSHOT CONTENTS:
    - All agent states and TODO lists
    - Git repository information
    - File system structure
    - Checkpoint states
    - Rule configurations 
    - Integration branch status
    - Performance metrics

EXAMPLES:
    $SCRIPT_NAME create pre-wave-2
    $SCRIPT_NAME create "before-phase-transition"
    $SCRIPT_NAME list
    $SCRIPT_NAME compare baseline current-state
    $SCRIPT_NAME cleanup 3

EOF
}

################################################################################
# CORE SNAPSHOT FUNCTIONS
################################################################################

create_snapshot() {
    local snapshot_name="${1:-auto-$(date '+%Y%m%d-%H%M%S')}"
    
    log "INFO" "Creating state snapshot: $snapshot_name"
    print_header "CREATING STATE SNAPSHOT: $snapshot_name"
    
    # Create snapshot directory
    local snapshot_dir="$SNAPSHOTS_DIR/$snapshot_name"
    mkdir -p "$snapshot_dir"
    
    # Create snapshot manifest
    create_snapshot_manifest "$snapshot_dir" "$snapshot_name"
    
    # Capture all state components
    print_info "📊 Capturing system state..."
    capture_system_state "$snapshot_dir"
    
    print_info "🗂️ Capturing agent states..."
    capture_agent_states "$snapshot_dir"
    
    print_info "📋 Capturing TODO states..."
    capture_todo_states "$snapshot_dir" 
    
    print_info "🔗 Capturing git state..."
    capture_git_state "$snapshot_dir"
    
    print_info "📁 Capturing file structure..."
    capture_file_structure "$snapshot_dir"
    
    print_info "🏗️ Capturing checkpoint states..."
    capture_checkpoint_states "$snapshot_dir"
    
    print_info "📏 Capturing metrics..."
    capture_metrics "$snapshot_dir"
    
    # Create snapshot summary
    create_snapshot_summary "$snapshot_dir" "$snapshot_name"
    
    # Compress the snapshot
    compress_snapshot "$snapshot_dir" "$snapshot_name"
    
    print_success "State snapshot created: $snapshot_name"
    print_info "Location: $snapshot_dir"
    
    return 0
}

create_snapshot_manifest() {
    local snapshot_dir="$1"
    local snapshot_name="$2"
    
    cat > "$snapshot_dir/MANIFEST.json" << EOF
{
  "snapshot_name": "$snapshot_name",
  "timestamp": "$(date '+%Y-%m-%d %H:%M:%S %Z')",
  "factory_version": "2.0",
  "creator": "state-snapshot.sh",
  "environment": {
    "hostname": "$(hostname)",
    "user": "$(whoami)",
    "working_directory": "$(pwd)",
    "shell": "$SHELL"
  },
  "components": [
    "system_state",
    "agent_states", 
    "todo_states",
    "git_state",
    "file_structure",
    "checkpoint_states",
    "metrics"
  ]
}
EOF
}

capture_system_state() {
    local snapshot_dir="$1"
    local system_dir="$snapshot_dir/system"
    mkdir -p "$system_dir"
    
    # System information
    cat > "$system_dir/system_info.txt" << EOF
################################################################################
# SYSTEM STATE SNAPSHOT
################################################################################
Timestamp: $(date '+%Y-%m-%d %H:%M:%S %Z')
Hostname: $(hostname)
User: $(whoami)
Shell: $SHELL
Working Directory: $(pwd)
PATH: $PATH

# System Resources
$(df -h)

# Process Information
$(ps aux | head -20)

# Environment Variables (filtered)
$(env | grep -E '^(HOME|USER|SHELL|PATH|PWD)=' | sort)

EOF

    # Capture network status if available
    if command -v netstat > /dev/null 2>&1; then
        echo "# Network Connections" >> "$system_dir/system_info.txt"
        netstat -tuln 2>/dev/null | head -10 >> "$system_dir/system_info.txt" || true
    fi
}

capture_agent_states() {
    local snapshot_dir="$1"
    local agents_dir="$snapshot_dir/agents"
    mkdir -p "$agents_dir"
    
    # Capture each agent state directory
    local agent_types=("orchestrator" "sw-engineer" "code-reviewer" "architect")
    
    for agent in "${agent_types[@]}"; do
        local agent_state_dir="$FACTORY_ROOT/agent-states/$agent"
        if [ -d "$agent_state_dir" ]; then
            print_info "  Capturing $agent state..."
            cp -r "$agent_state_dir" "$agents_dir/" 2>/dev/null || true
            
            # Create agent summary
            cat > "$agents_dir/$agent/SNAPSHOT_INFO.md" << EOF
# Agent State Snapshot: $agent

Captured: $(date '+%Y-%m-%d %H:%M:%S %Z')
Source: $agent_state_dir

## Directory Structure:
$(find "$agent_state_dir" -type f -name "*.md" | head -20)

## Recent Activity:
$(find "$agent_state_dir" -type f -newer /tmp/24hours 2>/dev/null || echo "No recent activity markers found")

EOF
        else
            print_warning "Agent state directory not found: $agent_state_dir"
        fi
    done
}

capture_todo_states() {
    local snapshot_dir="$1" 
    local todos_dir="$snapshot_dir/todos"
    mkdir -p "$todos_dir"
    
    # Copy TODO files if they exist
    if [ -d "$FACTORY_ROOT/todos" ]; then
        cp -r "$FACTORY_ROOT/todos"/* "$todos_dir/" 2>/dev/null || true
    fi
    
    # Also check for TODO files in other locations
    # Use CLAUDE_PROJECT_DIR if set
    local project_dir="${CLAUDE_PROJECT_DIR:-/workspaces/software-factory-2.0-template}"
    local todo_locations=(
        "${project_dir}/todos"
    )
    
    for location in "${todo_locations[@]}"; do
        if [ -d "$location" ]; then
            local basename_loc=$(basename "$location")
            mkdir -p "$todos_dir/external/$basename_loc"
            cp "$location"/*.todo "$todos_dir/external/$basename_loc/" 2>/dev/null || true
        fi
    done
    
    # Create TODO state summary
    cat > "$todos_dir/TODO_SUMMARY.md" << EOF
# TODO State Summary

Captured: $(date '+%Y-%m-%d %H:%M:%S %Z')

## TODO Files Found:
$(find "$todos_dir" -name "*.todo" | wc -l) total files

### By Agent:
$(find "$todos_dir" -name "*.todo" | xargs -I {} basename {} | cut -d'-' -f1 | sort | uniq -c)

### Recent Files (last 24 hours):
$(find "$todos_dir" -name "*.todo" -mtime -1 2>/dev/null | xargs -I {} basename {} || echo "None")

EOF
}

capture_git_state() {
    local snapshot_dir="$1"
    local git_dir="$snapshot_dir/git"
    mkdir -p "$git_dir"
    
    # Change to factory root for git operations
    local current_dir="$(pwd)"
    cd "$FACTORY_ROOT"
    
    if git rev-parse --git-dir > /dev/null 2>&1; then
        # Basic git information
        cat > "$git_dir/git_info.txt" << EOF
################################################################################
# GIT STATE SNAPSHOT  
################################################################################
Timestamp: $(date '+%Y-%m-%d %H:%M:%S %Z')
Repository Root: $(git rev-parse --show-toplevel)

# Current Branch and Commit
Branch: $(git branch --show-current)
Commit: $(git rev-parse HEAD)
Commit Message: $(git log -1 --pretty=format:"%s")
Author: $(git log -1 --pretty=format:"%an <%ae>")
Date: $(git log -1 --pretty=format:"%ad")

# Repository Status
$(git status --porcelain)

# Branch Information
$(git branch -a)

# Recent Commits (last 10)
$(git log --oneline -10)

# Remote Information
$(git remote -v)

# Stash Information
$(git stash list)

EOF

        # Capture diff if there are changes
        if ! git diff --quiet; then
            echo "# Working Directory Changes" >> "$git_dir/git_info.txt"
            git diff >> "$git_dir/working_changes.diff" 2>/dev/null || true
        fi
        
        # Capture staged changes
        if ! git diff --cached --quiet; then
            echo "# Staged Changes" >> "$git_dir/git_info.txt"  
            git diff --cached >> "$git_dir/staged_changes.diff" 2>/dev/null || true
        fi
        
    else
        echo "Not a git repository" > "$git_dir/git_info.txt"
    fi
    
    cd "$current_dir"
}

capture_file_structure() {
    local snapshot_dir="$1"
    local files_dir="$snapshot_dir/file_structure"
    mkdir -p "$files_dir"
    
    # Create directory tree
    if command -v tree > /dev/null 2>&1; then
        tree "$FACTORY_ROOT" > "$files_dir/directory_tree.txt" 2>/dev/null || true
    else
        find "$FACTORY_ROOT" -type d | sort > "$files_dir/directories.txt"
        find "$FACTORY_ROOT" -type f | sort > "$files_dir/files.txt"
    fi
    
    # File statistics
    cat > "$files_dir/file_stats.txt" << EOF
################################################################################
# FILE STRUCTURE STATISTICS
################################################################################  
Timestamp: $(date '+%Y-%m-%d %H:%M:%S %Z')

# Counts
Total Directories: $(find "$FACTORY_ROOT" -type d | wc -l)
Total Files: $(find "$FACTORY_ROOT" -type f | wc -l)

# File Types
$(find "$FACTORY_ROOT" -type f -name "*.*" | sed 's/.*\.//' | sort | uniq -c | sort -nr | head -10)

# Largest Files
$(find "$FACTORY_ROOT" -type f -exec ls -la {} \; | sort -k5 -nr | head -10)

# Recently Modified (last 24 hours)
$(find "$FACTORY_ROOT" -type f -mtime -1 | head -20)

EOF
}

capture_checkpoint_states() {
    local snapshot_dir="$1"
    local checkpoints_dir="$snapshot_dir/checkpoints"
    mkdir -p "$checkpoints_dir"
    
    # Copy checkpoint directories if they exist
    if [ -d "$FACTORY_ROOT/checkpoints" ]; then
        cp -r "$FACTORY_ROOT/checkpoints" "$checkpoints_dir/factory_checkpoints" 2>/dev/null || true
    fi
    
    # Look for checkpoints in agent states
    for agent_dir in "$FACTORY_ROOT"/agent-states/*/; do
        if [ -d "$agent_dir" ]; then
            local agent_name=$(basename "$agent_dir")
            find "$agent_dir" -name "checkpoint.md" -o -name "*.checkpoint" | while read -r checkpoint; do
                local rel_path=${checkpoint#$FACTORY_ROOT/}
                mkdir -p "$checkpoints_dir/$(dirname "$rel_path")"
                cp "$checkpoint" "$checkpoints_dir/$rel_path" 2>/dev/null || true
            done
        fi
    done
    
    # Create checkpoint summary
    cat > "$checkpoints_dir/CHECKPOINT_SUMMARY.md" << EOF
# Checkpoint State Summary

Captured: $(date '+%Y-%m-%d %H:%M:%S %Z')

## Checkpoint Files Found:
$(find "$checkpoints_dir" -name "*.md" -o -name "*.checkpoint" | wc -l) total files

### By Type:
$(find "$checkpoints_dir" -name "*.md" | wc -l) Markdown checkpoints
$(find "$checkpoints_dir" -name "*.checkpoint" | wc -l) Binary checkpoints

### Recent Checkpoints:
$(find "$checkpoints_dir" -name "checkpoint.md" -mtime -1 | head -10)

EOF
}

capture_metrics() {
    local snapshot_dir="$1"
    local metrics_dir="$snapshot_dir/metrics"
    mkdir -p "$metrics_dir"
    
    # System metrics
    cat > "$metrics_dir/performance.txt" << EOF
################################################################################
# PERFORMANCE METRICS
################################################################################
Timestamp: $(date '+%Y-%m-%d %H:%M:%S %Z')

# Disk Usage
$(du -sh "$FACTORY_ROOT")
$(du -sh "$FACTORY_ROOT"/* | sort -hr)

# Memory Usage (if available)
$(free -h 2>/dev/null || echo "Memory info not available")

# Load Average (if available)
$(uptime 2>/dev/null || echo "Load info not available")

EOF

    # Project metrics
    cat > "$metrics_dir/project_stats.txt" << EOF
# Project Statistics

## Markdown Files:
$(find "$FACTORY_ROOT" -name "*.md" | wc -l) total
$(find "$FACTORY_ROOT" -name "*.md" -exec wc -l {} \; | awk '{sum+=$1} END {print sum}') total lines

## Configuration Files:
$(find "$FACTORY_ROOT" -name "*.yaml" -o -name "*.yml" -o -name "*.json" | wc -l) total

## Scripts:
$(find "$FACTORY_ROOT" -name "*.sh" -o -name "*.bash" | wc -l) total

## Agent State Files:
$(find "$FACTORY_ROOT"/agent-states -name "*.md" | wc -l) total

EOF

    # TODO metrics
    if [ -d "$FACTORY_ROOT/todos" ]; then
        cat > "$metrics_dir/todo_stats.txt" << EOF
# TODO Statistics

## TODO Files:
$(find "$FACTORY_ROOT/todos" -name "*.todo" | wc -l) total files

## By Agent:
$(find "$FACTORY_ROOT/todos" -name "*.todo" | xargs -I {} basename {} | cut -d'-' -f1 | sort | uniq -c)

## File Sizes:
$(find "$FACTORY_ROOT/todos" -name "*.todo" -exec ls -la {} \; | awk '{print $5 " " $9}' | sort -nr)

EOF
    fi
}

create_snapshot_summary() {
    local snapshot_dir="$1"
    local snapshot_name="$2"
    
    cat > "$snapshot_dir/SNAPSHOT_SUMMARY.md" << EOF
# Software Factory 2.0 - State Snapshot Summary

**Snapshot Name:** $snapshot_name  
**Created:** $(date '+%Y-%m-%d %H:%M:%S %Z')  
**Creator:** $(whoami)@$(hostname)  
**Factory Version:** 2.0  

## Snapshot Contents

### 📊 System State
- System information and resources
- Environment variables
- Process information

### 🤖 Agent States
- Orchestrator state files
- Software Engineer working directories
- Code Reviewer feedback and plans
- Architect review documents

### 📋 TODO States
- All agent TODO files
- External TODO locations
- TODO state summaries

### 🔗 Git State
- Current branch and commit information
- Working directory changes
- Staged changes
- Remote configuration

### 📁 File Structure
- Complete directory tree
- File statistics and counts
- Recently modified files

### 🏗️ Checkpoint States
- Agent checkpoint files
- Factory checkpoint configuration
- Checkpoint summaries

### 📏 Metrics
- Performance statistics
- Project metrics
- TODO statistics

## Quick Stats

- **Total Size:** $(du -sh "$snapshot_dir" | cut -f1)
- **Files Captured:** $(find "$snapshot_dir" -type f | wc -l)
- **Directories:** $(find "$snapshot_dir" -type d | wc -l)

## Usage

To analyze this snapshot:
\`\`\`bash
./state-snapshot.sh analyze $snapshot_name
\`\`\`

To compare with another snapshot:
\`\`\`bash  
./state-snapshot.sh compare $snapshot_name other-snapshot
\`\`\`

## Recovery Notes

This snapshot can be used to understand the state of the Software Factory
at the time of capture. It includes all necessary information to restore
context after compaction or to analyze system state changes over time.

EOF
}

compress_snapshot() {
    local snapshot_dir="$1"
    local snapshot_name="$2"
    
    if command -v tar > /dev/null 2>&1; then
        local archive_name="$SNAPSHOTS_DIR/${snapshot_name}.tar.gz"
        
        print_info "Compressing snapshot to: $(basename "$archive_name")"
        
        # Create compressed archive
        cd "$SNAPSHOTS_DIR"
        tar -czf "${snapshot_name}.tar.gz" "$snapshot_name" 2>/dev/null || {
            print_warning "Compression failed, keeping uncompressed snapshot"
            return 0
        }
        
        # Verify archive
        if tar -tzf "${snapshot_name}.tar.gz" > /dev/null 2>&1; then
            print_success "Snapshot compressed successfully"
            
            # Ask if user wants to remove uncompressed version
            read -p "Remove uncompressed snapshot directory? (y/N): " -r
            if [[ $REPLY =~ ^[Yy]$ ]]; then
                rm -rf "$snapshot_dir"
                print_info "Uncompressed snapshot removed"
            fi
        else
            print_error "Archive verification failed"
            rm -f "${snapshot_name}.tar.gz"
        fi
    else
        print_info "tar not available, keeping uncompressed snapshot"
    fi
}

################################################################################
# SNAPSHOT MANAGEMENT FUNCTIONS
################################################################################

list_snapshots() {
    print_header "SOFTWARE FACTORY 2.0 - AVAILABLE SNAPSHOTS"
    
    if [ ! -d "$SNAPSHOTS_DIR" ]; then
        print_info "No snapshots directory found"
        return 0
    fi
    
    # List directories (uncompressed snapshots)
    local dirs=($(find "$SNAPSHOTS_DIR" -maxdepth 1 -type d -not -path "$SNAPSHOTS_DIR" | sort))
    
    # List archives (compressed snapshots)
    local archives=($(find "$SNAPSHOTS_DIR" -maxdepth 1 -name "*.tar.gz" | sort))
    
    if [ ${#dirs[@]} -eq 0 ] && [ ${#archives[@]} -eq 0 ]; then
        print_info "No snapshots found"
        return 0
    fi
    
    print_info "Uncompressed Snapshots:"
    for dir in "${dirs[@]}"; do
        local name=$(basename "$dir")
        local size=$(du -sh "$dir" | cut -f1)
        local mod_time=$(stat -c '%Y' "$dir")
        local readable_time=$(date -d "@$mod_time" '+%Y-%m-%d %H:%M:%S')
        printf "  %-30s %10s  %s\n" "$name" "$size" "$readable_time"
    done
    
    if [ ${#archives[@]} -gt 0 ]; then
        echo
        print_info "Compressed Snapshots:"
        for archive in "${archives[@]}"; do
            local name=$(basename "$archive" .tar.gz)
            local size=$(du -sh "$archive" | cut -f1)
            local mod_time=$(stat -c '%Y' "$archive")
            local readable_time=$(date -d "@$mod_time" '+%Y-%m-%d %H:%M:%S')
            printf "  %-30s %10s  %s (compressed)\n" "$name" "$size" "$readable_time"
        done
    fi
}

analyze_snapshot() {
    local snapshot_name="$1"
    
    print_header "ANALYZING SNAPSHOT: $snapshot_name"
    
    local snapshot_dir="$SNAPSHOTS_DIR/$snapshot_name"
    local archive_file="$SNAPSHOTS_DIR/${snapshot_name}.tar.gz"
    
    # Check if snapshot exists
    if [ -d "$snapshot_dir" ]; then
        print_success "Found uncompressed snapshot"
    elif [ -f "$archive_file" ]; then
        print_info "Found compressed snapshot, extracting for analysis..."
        cd "$SNAPSHOTS_DIR"
        tar -xzf "${snapshot_name}.tar.gz" || error_exit "Failed to extract snapshot"
    else
        error_exit "Snapshot not found: $snapshot_name"
    fi
    
    # Display summary if available
    if [ -f "$snapshot_dir/SNAPSHOT_SUMMARY.md" ]; then
        print_info "Snapshot Summary:"
        cat "$snapshot_dir/SNAPSHOT_SUMMARY.md"
    fi
    
    # Show detailed analysis
    echo
    print_info "Detailed Analysis:"
    
    if [ -f "$snapshot_dir/MANIFEST.json" ]; then
        print_info "Manifest Information:"
        cat "$snapshot_dir/MANIFEST.json" | jq . 2>/dev/null || cat "$snapshot_dir/MANIFEST.json"
    fi
    
    # File structure analysis
    echo
    print_info "Directory Structure:"
    if command -v tree > /dev/null 2>&1; then
        tree "$snapshot_dir" -L 3
    else
        find "$snapshot_dir" -type d | sort | head -20
    fi
    
    # Size analysis
    echo  
    print_info "Size Breakdown:"
    du -sh "$snapshot_dir"/* | sort -hr
}

compare_snapshots() {
    local snapshot1="$1"
    local snapshot2="$2"
    
    print_header "COMPARING SNAPSHOTS: $snapshot1 vs $snapshot2"
    
    # Ensure both snapshots exist
    local dir1="$SNAPSHOTS_DIR/$snapshot1"  
    local dir2="$SNAPSHOTS_DIR/$snapshot2"
    
    for snapshot in "$snapshot1" "$snapshot2"; do
        local dir="$SNAPSHOTS_DIR/$snapshot"
        local archive="$SNAPSHOTS_DIR/${snapshot}.tar.gz"
        
        if [ ! -d "$dir" ] && [ -f "$archive" ]; then
            print_info "Extracting compressed snapshot: $snapshot"
            cd "$SNAPSHOTS_DIR"
            tar -xzf "${snapshot}.tar.gz" || error_exit "Failed to extract $snapshot"
        elif [ ! -d "$dir" ]; then
            error_exit "Snapshot not found: $snapshot"
        fi
    done
    
    # Compare manifests
    if [ -f "$dir1/MANIFEST.json" ] && [ -f "$dir2/MANIFEST.json" ]; then
        print_info "Manifest Comparison:"
        echo "Snapshot 1 created: $(jq -r '.timestamp' "$dir1/MANIFEST.json" 2>/dev/null || echo 'unknown')"
        echo "Snapshot 2 created: $(jq -r '.timestamp' "$dir2/MANIFEST.json" 2>/dev/null || echo 'unknown')"
    fi
    
    # Size comparison
    echo
    print_info "Size Comparison:"
    local size1=$(du -sb "$dir1" | cut -f1)
    local size2=$(du -sb "$dir2" | cut -f1)
    echo "Snapshot 1: $(du -sh "$dir1" | cut -f1)"
    echo "Snapshot 2: $(du -sh "$dir2" | cut -f1)"
    
    if [ $size1 -gt $size2 ]; then
        echo "Snapshot 1 is larger by $(awk "BEGIN {printf \"%.1f\", ($size1-$size2)/1024/1024}") MB"
    elif [ $size2 -gt $size1 ]; then
        echo "Snapshot 2 is larger by $(awk "BEGIN {printf \"%.1f\", ($size2-$size1)/1024/1024}") MB"
    else
        echo "Snapshots are the same size"
    fi
    
    # File count comparison
    echo
    print_info "File Count Comparison:"
    local files1=$(find "$dir1" -type f | wc -l)
    local files2=$(find "$dir2" -type f | wc -l)
    echo "Snapshot 1: $files1 files"
    echo "Snapshot 2: $files2 files"
    
    # TODO comparison if available
    if [ -d "$dir1/todos" ] && [ -d "$dir2/todos" ]; then
        echo
        print_info "TODO State Comparison:"
        local todos1=$(find "$dir1/todos" -name "*.todo" | wc -l)
        local todos2=$(find "$dir2/todos" -name "*.todo" | wc -l)
        echo "Snapshot 1 TODO files: $todos1"
        echo "Snapshot 2 TODO files: $todos2"
    fi
}

cleanup_snapshots() {
    local days_old="${1:-7}"
    
    print_header "CLEANING UP OLD SNAPSHOTS (older than $days_old days)"
    
    if [ ! -d "$SNAPSHOTS_DIR" ]; then
        print_info "No snapshots directory found"
        return 0
    fi
    
    # Find old snapshots
    local old_dirs=($(find "$SNAPSHOTS_DIR" -maxdepth 1 -type d -mtime +$days_old -not -path "$SNAPSHOTS_DIR"))
    local old_archives=($(find "$SNAPSHOTS_DIR" -maxdepth 1 -name "*.tar.gz" -mtime +$days_old))
    
    local total_old=$((${#old_dirs[@]} + ${#old_archives[@]}))
    
    if [ $total_old -eq 0 ]; then
        print_info "No old snapshots found to clean up"
        return 0
    fi
    
    print_warning "Found $total_old snapshots older than $days_old days:"
    
    for dir in "${old_dirs[@]}"; do
        local name=$(basename "$dir")
        local size=$(du -sh "$dir" | cut -f1)
        echo "  Directory: $name ($size)"
    done
    
    for archive in "${old_archives[@]}"; do
        local name=$(basename "$archive")
        local size=$(du -sh "$archive" | cut -f1)
        echo "  Archive: $name ($size)"
    done
    
    echo
    read -p "Proceed with cleanup? (y/N): " -r
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        for dir in "${old_dirs[@]}"; do
            rm -rf "$dir"
            print_info "Removed directory: $(basename "$dir")"
        done
        
        for archive in "${old_archives[@]}"; do
            rm -f "$archive"
            print_info "Removed archive: $(basename "$archive")"
        done
        
        print_success "Cleanup complete"
    else
        print_info "Cleanup cancelled"
    fi
}

################################################################################
# MAIN COMMAND DISPATCH
################################################################################

main() {
    # Ensure snapshots directory exists
    mkdir -p "$SNAPSHOTS_DIR"
    
    # Parse command
    local command="${1:-help}"
    
    case "$command" in
        "create")
            create_snapshot "${2:-}"
            ;;
        "list")
            list_snapshots
            ;;
        "analyze")
            if [ $# -lt 2 ]; then
                error_exit "Usage: $SCRIPT_NAME analyze <snapshot-name>"
            fi
            analyze_snapshot "$2"
            ;;
        "compare")
            if [ $# -lt 3 ]; then
                error_exit "Usage: $SCRIPT_NAME compare <snapshot1> <snapshot2>"
            fi
            compare_snapshots "$2" "$3"
            ;;
        "cleanup")
            cleanup_snapshots "${2:-7}"
            ;;
        "restore")
            print_error "Restore functionality not yet implemented"
            error_exit "Use 'analyze' to inspect snapshots manually"
            ;;
        "help"|"--help"|"-h")
            usage
            ;;
        *)
            print_error "Unknown command: $command"
            echo
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