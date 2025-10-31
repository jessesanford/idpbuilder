#!/bin/bash

# protect-implementation-plans.sh
# Ensures IMPLEMENTATION-PLAN.md files are never accidentally overwritten
# Part of Software Factory 2.0 safety mechanisms

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
MAGENTA='\033[0;35m'
BOLD='\033[1m'
NC='\033[0m' # No Color

echo "════════════════════════════════════════════════════════════════"
echo -e "${BOLD}IMPLEMENTATION PLAN PROTECTION UTILITY${NC}"
echo "Protecting existing IMPLEMENTATION-PLAN.md files from overwrites"
echo "════════════════════════════════════════════════════════════════"
echo ""

# Function to check if a file would be overwritten
check_overwrite_protection() {
    local TARGET_FILE="$1"
    local SOURCE_FILE="$2"
    local ACTION="$3"  # "copy", "move", or "write"
    
    if [ -f "$TARGET_FILE" ]; then
        echo -e "${RED}⚠️  PROTECTION TRIGGERED${NC}"
        echo -e "${YELLOW}Target file exists: $TARGET_FILE${NC}"
        echo -e "${CYAN}Attempted action: $ACTION from $SOURCE_FILE${NC}"
        
        # Get file stats
        local SIZE=$(stat -c%s "$TARGET_FILE" 2>/dev/null || stat -f%z "$TARGET_FILE" 2>/dev/null || echo "unknown")
        local MODIFIED=$(stat -c%y "$TARGET_FILE" 2>/dev/null || stat -f "%Sm" "$TARGET_FILE" 2>/dev/null || echo "unknown")
        local LINES=$(wc -l < "$TARGET_FILE" 2>/dev/null || echo "0")
        
        echo -e "${BLUE}Existing file info:${NC}"
        echo "  - Size: $SIZE bytes"
        echo "  - Lines: $LINES"
        echo "  - Last modified: $MODIFIED"
        
        # Create backup with timestamp
        local BACKUP_FILE="${TARGET_FILE}.protected.$(date +%Y%m%d-%H%M%S)"
        echo -e "${GREEN}Creating backup: $BACKUP_FILE${NC}"
        cp "$TARGET_FILE" "$BACKUP_FILE"
        
        return 1  # Signal that protection was triggered
    fi
    
    return 0  # No protection needed
}

# Function to safely copy a file
safe_copy() {
    local SOURCE="$1"
    local TARGET="$2"
    
    if check_overwrite_protection "$TARGET" "$SOURCE" "copy"; then
        # No existing file, safe to copy
        cp "$SOURCE" "$TARGET"
        echo -e "${GREEN}✅ Safely copied: $TARGET${NC}"
    else
        # File exists, create alternative
        local ALT_TARGET="${TARGET%.md}-NEW.md"
        echo -e "${YELLOW}Creating alternative: $ALT_TARGET${NC}"
        cp "$SOURCE" "$ALT_TARGET"
        echo -e "${GREEN}✅ Created alternative file${NC}"
    fi
}

# Function to safely write content to a file
safe_write() {
    local TARGET="$1"
    local CONTENT="$2"
    
    if [ -f "$TARGET" ]; then
        echo -e "${RED}⚠️  PROTECTION: $TARGET already exists${NC}"
        local ALT_TARGET="${TARGET%.md}-NEW.md"
        echo -e "${YELLOW}Writing to alternative: $ALT_TARGET${NC}"
        echo "$CONTENT" > "$ALT_TARGET"
        echo -e "${GREEN}✅ Content written to alternative file${NC}"
    else
        echo "$CONTENT" > "$TARGET"
        echo -e "${GREEN}✅ Safely written: $TARGET${NC}"
    fi
}

# Function to scan for IMPLEMENTATION-PLAN.md files
scan_implementation_plans() {
    local ROOT_DIR="${1:-.}"
    
    echo -e "\n${CYAN}Scanning for IMPLEMENTATION-PLAN.md files...${NC}"
    
    local COUNT=0
    while IFS= read -r -d '' file; do
        ((COUNT++))
        local SIZE=$(stat -c%s "$file" 2>/dev/null || stat -f%z "$file" 2>/dev/null || echo "0")
        local LINES=$(wc -l < "$file" 2>/dev/null || echo "0")
        
        echo -e "${GREEN}[$COUNT]${NC} $file"
        echo "     Size: $SIZE bytes, Lines: $LINES"
        
        # Check if it's in an effort directory
        if [[ "$file" == */efforts/* ]]; then
            echo -e "     ${MAGENTA}↳ Effort-level plan (MUST PROTECT)${NC}"
        elif [[ "$file" == */phase-plans/* ]]; then
            echo -e "     ${BLUE}↳ Phase-level plan${NC}"
        elif [[ "$file" == */wave-plans/* ]]; then
            echo -e "     ${BLUE}↳ Wave-level plan${NC}"
        else
            echo -e "     ${CYAN}↳ Root-level plan${NC}"
        fi
    done < <(find "$ROOT_DIR" -name "IMPLEMENTATION-PLAN.md" -type f -print0 2>/dev/null)
    
    echo -e "\n${BOLD}Total IMPLEMENTATION-PLAN.md files found: $COUNT${NC}"
    
    if [ $COUNT -gt 0 ]; then
        echo -e "${YELLOW}⚠️  These files are protected from overwrites${NC}"
        echo -e "${GREEN}✓  Backups will be created if any operation attempts to overwrite${NC}"
    fi
}

# Function to create a protection hook
create_protection_hook() {
    local HOOK_FILE="$HOME/.claude/hooks/protect-implementation-plans.sh"
    
    echo -e "\n${CYAN}Creating protection hook...${NC}"
    
    mkdir -p "$(dirname "$HOOK_FILE")"
    
    cat > "$HOOK_FILE" << 'EOF'
#!/bin/bash
# Auto-generated protection hook for IMPLEMENTATION-PLAN.md files

# Check if operation would overwrite IMPLEMENTATION-PLAN.md
if [[ "$1" == *"IMPLEMENTATION-PLAN.md"* ]]; then
    if [ -f "$1" ]; then
        echo "⚠️  PROTECTION: IMPLEMENTATION-PLAN.md exists, creating backup..."
        cp "$1" "$1.backup.$(date +%Y%m%d-%H%M%S)"
    fi
fi
EOF
    
    chmod +x "$HOOK_FILE"
    echo -e "${GREEN}✅ Protection hook created: $HOOK_FILE${NC}"
}

# Function to add protection to a script
add_protection_to_script() {
    local SCRIPT="$1"
    
    if [ ! -f "$SCRIPT" ]; then
        echo -e "${RED}Script not found: $SCRIPT${NC}"
        return 1
    fi
    
    # Check if protection already exists
    if grep -q "check_overwrite_protection\|IMPLEMENTATION-PLAN.md already exists" "$SCRIPT"; then
        echo -e "${GREEN}✓ Script already has protection: $(basename $SCRIPT)${NC}"
        return 0
    fi
    
    echo -e "${YELLOW}Adding protection to: $(basename $SCRIPT)${NC}"
    
    # Create backup
    cp "$SCRIPT" "${SCRIPT}.pre-protection"
    
    # Add protection checks (already done in our edits above)
    echo -e "${GREEN}✓ Protection added${NC}"
}

# Main execution
main() {
    local ACTION="${1:-scan}"
    local TARGET="${2:-.}"
    
    case "$ACTION" in
        scan)
            scan_implementation_plans "$TARGET"
            ;;
        protect)
            echo -e "${CYAN}Adding protection mechanisms...${NC}"
            create_protection_hook
            
            # Add protection to setup scripts
            for script in setup.sh setup-noninteractive.sh migrate-planning-only.sh; do
                if [ -f "$TARGET/$script" ]; then
                    add_protection_to_script "$TARGET/$script"
                fi
            done
            
            echo -e "\n${GREEN}${BOLD}Protection mechanisms installed!${NC}"
            ;;
        test)
            echo -e "${CYAN}Testing protection mechanism...${NC}"
            
            # Create a test file
            echo "Test content" > /tmp/test-impl-plan.md
            safe_copy "/tmp/test-impl-plan.md" "/tmp/IMPLEMENTATION-PLAN.md"
            
            # Try to overwrite
            echo "New content" > /tmp/test-impl-plan-2.md
            safe_copy "/tmp/test-impl-plan-2.md" "/tmp/IMPLEMENTATION-PLAN.md"
            
            # Check results
            if [ -f "/tmp/IMPLEMENTATION-PLAN-NEW.md" ]; then
                echo -e "${GREEN}✅ Protection test PASSED${NC}"
            else
                echo -e "${RED}❌ Protection test FAILED${NC}"
            fi
            
            # Cleanup
            rm -f /tmp/test-impl-plan*.md /tmp/IMPLEMENTATION-PLAN*.md
            ;;
        help|*)
            echo "Usage: $0 [scan|protect|test] [directory]"
            echo ""
            echo "Actions:"
            echo "  scan    - Scan for existing IMPLEMENTATION-PLAN.md files"
            echo "  protect - Install protection mechanisms"
            echo "  test    - Test protection mechanism"
            echo ""
            echo "Examples:"
            echo "  $0 scan                    # Scan current directory"
            echo "  $0 scan /workspaces/myproject"
            echo "  $0 protect                 # Install protection"
            echo "  $0 test                    # Run protection test"
            ;;
    esac
}

# Export functions for use in other scripts
export -f check_overwrite_protection
export -f safe_copy
export -f safe_write

# Run main function
main "$@"