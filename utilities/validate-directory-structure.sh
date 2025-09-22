#!/bin/bash
# validate-directory-structure.sh - Prevent nested effort directory bugs

set -euo pipefail

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo "======================================================================"
echo "Directory Structure Validation Tool"
echo "Prevents nested effort/split directory creation bugs"
echo "======================================================================"

# Function to check for nested effort structures
check_nested_efforts() {
    local dir="${1:-$(pwd)}"
    
    echo -e "\n${YELLOW}Checking for nested effort structures in: $dir${NC}"
    
    # Look for the catastrophic pattern: /efforts/*/efforts/
    if find "$dir" -type d -path "*/efforts/*/efforts/*" 2>/dev/null | grep -q .; then
        echo -e "${RED}🔴🔴🔴 CRITICAL: NESTED EFFORT STRUCTURE DETECTED!${NC}"
        echo "Found directories with duplicate '/efforts/' paths:"
        find "$dir" -type d -path "*/efforts/*/efforts/*" 2>/dev/null | while read -r nested; do
            echo -e "  ${RED}❌ $nested${NC}"
        done
        echo ""
        echo "This is the CATASTROPHIC nesting bug!"
        echo "Splits/efforts must be SIBLINGS, not nested!"
        return 1
    else
        echo -e "${GREEN}✅ No nested effort structures found${NC}"
        return 0
    fi
}

# Function to validate split placement
validate_split_placement() {
    local phase="$1"
    local wave="$2"
    local effort_name="$3"
    
    echo -e "\n${YELLOW}Validating split placement for $effort_name${NC}"
    
    local base_dir="${CLAUDE_PROJECT_DIR}/efforts/phase${phase}/wave${wave}"
    local effort_dir="${base_dir}/${effort_name}"
    
    # Check if splits are siblings (correct) or nested (wrong)
    local correct_splits=0
    local wrong_splits=0
    
    # Correct: splits at same level as effort
    for split in "${base_dir}/${effort_name}-split-"*; do
        if [ -d "$split" ]; then
            echo -e "  ${GREEN}✅ Correct: $(basename "$split") is sibling to $effort_name${NC}"
            ((correct_splits++))
        fi
    done
    
    # Wrong: splits nested inside effort
    if [ -d "$effort_dir" ]; then
        for nested in "${effort_dir}/"*split*; do
            if [ -d "$nested" ]; then
                echo -e "  ${RED}❌ WRONG: $(basename "$nested") is INSIDE $effort_name!${NC}"
                ((wrong_splits++))
            fi
        done
        
        # Check for nested efforts directory
        if [ -d "${effort_dir}/efforts" ]; then
            echo -e "  ${RED}❌ CATASTROPHIC: Found 'efforts' directory inside $effort_name!${NC}"
            ((wrong_splits++))
        fi
    fi
    
    if [ $wrong_splits -gt 0 ]; then
        echo -e "${RED}FAILURE: Found $wrong_splits incorrectly placed splits/directories${NC}"
        return 1
    elif [ $correct_splits -gt 0 ]; then
        echo -e "${GREEN}SUCCESS: All $correct_splits splits correctly placed as siblings${NC}"
        return 0
    else
        echo "No splits found for $effort_name"
        return 0
    fi
}

# Function to show correct vs wrong structure
show_structure_examples() {
    echo -e "\n${YELLOW}Directory Structure Examples:${NC}"
    
    cat << 'EOF'

✅ CORRECT STRUCTURE (Splits as siblings):
────────────────────────────────────────
efforts/
└── phase2/
    └── wave1/
        ├── gitea-client/              # Original effort
        ├── gitea-client-split-001/    # Split 1 - SIBLING
        ├── gitea-client-split-002/    # Split 2 - SIBLING
        ├── gitea-client-split-003/    # Split 3 - SIBLING
        └── other-effort/              # Another effort

❌ WRONG STRUCTURE (Nested splits):
──────────────────────────────────
efforts/
└── phase2/
    └── wave1/
        └── gitea-client/                              # Original effort
            └── efforts/                               # ❌ DUPLICATE PATH!
                └── phase2/                            # ❌ NESTED!
                    └── wave1/                         # ❌ WRONG!
                        └── gitea-client-split-001/   # ❌ DEEPLY NESTED!

EOF
}

# Function to fix common path mistakes
suggest_fixes() {
    local current_dir="$1"
    
    echo -e "\n${YELLOW}Suggested fixes for common mistakes:${NC}"
    
    if [[ "$current_dir" == *"/efforts/"*"/efforts/"* ]]; then
        echo -e "${RED}You have a nested structure!${NC}"
        echo "Current directory has duplicate '/efforts/' segments"
        echo ""
        echo "To fix:"
        echo "1. cd to SF root: cd \$CLAUDE_PROJECT_DIR"
        echo "2. Create splits at correct level:"
        echo "   mkdir -p efforts/phaseX/waveY/effort-split-XXX"
        echo "3. Never cd into an effort before creating splits!"
    fi
    
    if [[ "$current_dir" == */efforts/phase*/wave*/* ]] && [[ "$current_dir" != *split* ]]; then
        echo -e "${YELLOW}You're inside an effort directory${NC}"
        echo "To create splits correctly:"
        echo "1. Return to SF root: cd \$CLAUDE_PROJECT_DIR"
        echo "2. Create split as sibling, not child"
        echo "3. Use absolute paths from SF root"
    fi
}

# Main validation logic
main() {
    local mode="${1:-check}"
    
    case "$mode" in
        check)
            # Check current directory structure
            check_nested_efforts "${CLAUDE_PROJECT_DIR}/efforts"
            ;;
            
        validate-split)
            # Validate specific split placement
            if [ $# -lt 4 ]; then
                echo "Usage: $0 validate-split <phase> <wave> <effort-name>"
                exit 1
            fi
            validate_split_placement "$2" "$3" "$4"
            ;;
            
        show-examples)
            # Show correct vs wrong structures
            show_structure_examples
            ;;
            
        suggest-fix)
            # Suggest fixes for current directory
            suggest_fixes "$(pwd)"
            ;;
            
        full)
            # Run all checks
            echo "Running full validation..."
            check_nested_efforts "${CLAUDE_PROJECT_DIR}/efforts"
            show_structure_examples
            
            # Check all existing efforts
            if [ -d "${CLAUDE_PROJECT_DIR}/efforts" ]; then
                find "${CLAUDE_PROJECT_DIR}/efforts" -mindepth 3 -maxdepth 3 -type d | while read -r wave_dir; do
                    if [[ "$wave_dir" =~ phase([0-9]+)/wave([0-9]+) ]]; then
                        phase="${BASH_REMATCH[1]}"
                        wave="${BASH_REMATCH[2]}"
                        
                        for effort_dir in "$wave_dir"/*; do
                            if [ -d "$effort_dir" ] && [[ ! "$(basename "$effort_dir")" == *"-split-"* ]]; then
                                effort_name=$(basename "$effort_dir")
                                validate_split_placement "$phase" "$wave" "$effort_name" || true
                            fi
                        done
                    fi
                done
            fi
            ;;
            
        *)
            echo "Usage: $0 {check|validate-split|show-examples|suggest-fix|full}"
            echo ""
            echo "Commands:"
            echo "  check           - Check for nested effort structures"
            echo "  validate-split  - Validate split placement for specific effort"
            echo "  show-examples   - Show correct vs wrong directory structures"
            echo "  suggest-fix     - Suggest fixes for current directory"
            echo "  full            - Run all validation checks"
            exit 1
            ;;
    esac
}

# Run main function
main "$@"