#!/bin/bash
# Verify Sequential Split Branching
# Ensures splits follow mandatory sequential pattern (each based on previous)

set -euo pipefail

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print colored output
print_error() { echo -e "${RED}❌ $1${NC}"; }
print_success() { echo -e "${GREEN}✅ $1${NC}"; }
print_warning() { echo -e "${YELLOW}⚠️  $1${NC}"; }
print_info() { echo -e "📊 $1"; }

# Function to verify sequential splits for an effort
verify_sequential_splits() {
    local effort_name="$1"
    local split_count="${2:-0}"
    
    echo ""
    echo "════════════════════════════════════════════════════════════"
    print_info "Verifying Sequential Branching for: $effort_name"
    echo "════════════════════════════════════════════════════════════"
    
    # If split count not provided, try to detect
    if [ "$split_count" -eq 0 ]; then
        print_info "Auto-detecting split count..."
        
        # Try to find splits by branch pattern
        local splits_found=$(git branch -a | grep -c "${effort_name}--split-" || true)
        
        if [ "$splits_found" -eq 0 ]; then
            print_error "No splits found for $effort_name"
            return 1
        fi
        
        split_count="$splits_found"
        print_info "Found $split_count splits"
    fi
    
    # Check first split (should be based on phase integration or main)
    local first_split="${effort_name}--split-001"
    
    if ! git rev-parse --verify "$first_split" >/dev/null 2>&1; then
        print_error "First split branch not found: $first_split"
        return 1
    fi
    
    print_success "Split-001 exists: $first_split"
    
    # Verify sequential relationship for remaining splits
    local all_sequential=true
    
    for i in $(seq 2 $split_count); do
        local current_num=$(printf "%03d" $i)
        local previous_num=$(printf "%03d" $((i-1)))
        
        local current_branch="${effort_name}--split-${current_num}"
        local previous_branch="${effort_name}--split-${previous_num}"
        
        # Check if branches exist
        if ! git rev-parse --verify "$current_branch" >/dev/null 2>&1; then
            print_warning "Split branch not found: $current_branch (skipping)"
            continue
        fi
        
        # Check if current is based on previous (sequential)
        if git merge-base --is-ancestor "$previous_branch" "$current_branch" 2>/dev/null; then
            print_success "split-${current_num} correctly based on split-${previous_num}"
            
            # Show line count if line-counter.sh exists
            if [ -f "tools/line-counter.sh" ] || [ -f "../tools/line-counter.sh" ]; then
                local line_counter=""
                [ -f "tools/line-counter.sh" ] && line_counter="tools/line-counter.sh"
                [ -f "../tools/line-counter.sh" ] && line_counter="../tools/line-counter.sh"
                
                local lines=$($line_counter -b "$previous_branch" -c "$current_branch" 2>/dev/null | grep "Total" | awk '{print $NF}' || echo "?")
                print_info "  Lines added in split-${current_num}: $lines"
                
                if [ "$lines" != "?" ] && [ "$lines" -gt 800 ]; then
                    print_error "  WARNING: Split exceeds 800 line limit!"
                fi
            fi
        else
            print_error "split-${current_num} NOT based on split-${previous_num}!"
            print_error "  This violates MANDATORY sequential branching requirement!"
            
            # Try to detect what it's actually based on
            local actual_base=$(git merge-base "$current_branch" "$previous_branch" 2>/dev/null || echo "unknown")
            print_error "  Appears to be based on: $actual_base"
            
            all_sequential=false
        fi
    done
    
    echo ""
    if [ "$all_sequential" = true ]; then
        print_success "ALL SPLITS FOLLOW SEQUENTIAL BRANCHING ✅"
        echo "Integration order: split-001 → split-002 → ... → split-${split_count}"
        return 0
    else
        print_error "SEQUENTIAL BRANCHING VIOLATIONS DETECTED ❌"
        echo ""
        echo "Required pattern:"
        echo "  split-001 based on phase-integration"
        echo "  split-002 based on split-001"
        echo "  split-003 based on split-002"
        echo "  etc..."
        return 1
    fi
}

# Function to check all efforts in orchestrator-state.json
check_all_split_efforts() {
    local state_file="${1:-orchestrator-state.json}"
    
    if [ ! -f "$state_file" ]; then
        print_warning "State file not found: $state_file"
        return 1
    fi
    
    print_info "Checking all split efforts in $state_file..."
    
    # Extract efforts with splits from state file
    local split_efforts=$(yq '.split_tracking | keys | .[]' "$state_file" 2>/dev/null || echo "")
    
    if [ -z "$split_efforts" ]; then
        print_info "No split efforts found in state file"
        return 0
    fi
    
    local total_efforts=0
    local valid_efforts=0
    
    for effort in $split_efforts; do
        ((total_efforts++))
        
        # Get split count from state file
        local split_count=$(yq ".split_tracking.\"$effort\".split_count" "$state_file" 2>/dev/null || echo "0")
        
        # Get branch prefix from state file
        local branch_prefix=$(yq ".split_tracking.\"$effort\".original_branch" "$state_file" 2>/dev/null || echo "$effort")
        
        if verify_sequential_splits "$branch_prefix" "$split_count"; then
            ((valid_efforts++))
        fi
    done
    
    echo ""
    echo "════════════════════════════════════════════════════════════"
    print_info "SUMMARY: $valid_efforts/$total_efforts efforts follow sequential branching"
    echo "════════════════════════════════════════════════════════════"
    
    if [ "$valid_efforts" -eq "$total_efforts" ]; then
        print_success "All split efforts are correctly sequential!"
        return 0
    else
        print_error "Some efforts have branching violations!"
        return 1
    fi
}

# Function to show usage
show_usage() {
    cat << EOF
Usage: $0 [OPTIONS] [EFFORT_NAME] [SPLIT_COUNT]

Verify that split branches follow MANDATORY sequential pattern.

OPTIONS:
    -a, --all           Check all split efforts in orchestrator-state.json
    -s, --state FILE    Use alternate state file (default: orchestrator-state.json)
    -h, --help          Show this help message

ARGUMENTS:
    EFFORT_NAME         Full branch name of the effort (without --split-XXX)
    SPLIT_COUNT         Number of splits (auto-detected if not provided)

EXAMPLES:
    $0 project/phase1/wave1/authentication 3
        Verify authentication effort has 3 sequential splits
    
    $0 --all
        Check all split efforts in orchestrator-state.json
    
    $0 tmc-workspace/phase2/wave1/api-types
        Auto-detect and verify splits for api-types

SEQUENTIAL BRANCHING RULES:
    - Split-001 must be based on phase-integration (or main)
    - Split-002 must be based on split-001
    - Split-003 must be based on split-002
    - And so on...

This is MANDATORY for correct line counting and integration!
EOF
}

# Main script logic
main() {
    local check_all=false
    local state_file="orchestrator-state.json"
    local effort_name=""
    local split_count=0
    
    # Parse arguments
    while [[ $# -gt 0 ]]; do
        case $1 in
            -a|--all)
                check_all=true
                shift
                ;;
            -s|--state)
                state_file="$2"
                shift 2
                ;;
            -h|--help)
                show_usage
                exit 0
                ;;
            *)
                if [ -z "$effort_name" ]; then
                    effort_name="$1"
                elif [ "$split_count" -eq 0 ]; then
                    split_count="$1"
                fi
                shift
                ;;
        esac
    done
    
    # Execute appropriate check
    if [ "$check_all" = true ]; then
        check_all_split_efforts "$state_file"
    elif [ -n "$effort_name" ]; then
        verify_sequential_splits "$effort_name" "$split_count"
    else
        print_error "No effort specified and --all not used"
        echo ""
        show_usage
        exit 1
    fi
}

# Run main function
main "$@"