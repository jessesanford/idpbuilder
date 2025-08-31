#!/bin/bash
# diagnose-target-repo.sh - Diagnose target repository configuration issues
# This script helps identify why an orchestrator might be confused about which repo to clone

set -euo pipefail

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
MAGENTA='\033[0;35m'
BOLD='\033[1m'
NC='\033[0m' # No Color

echo -e "${MAGENTA}${BOLD}════════════════════════════════════════════════════════${NC}"
echo -e "${MAGENTA}${BOLD}   TARGET REPOSITORY CONFIGURATION DIAGNOSTIC${NC}"
echo -e "${MAGENTA}${BOLD}════════════════════════════════════════════════════════${NC}\n"

# Function to check if we're in an SF instance
check_sf_instance() {
    echo -e "${CYAN}1. Checking if this is a Software Factory instance...${NC}"
    
    if [ -f "orchestrator-state.yaml" ] || [ -f "orchestrator-state.yaml.example" ]; then
        echo -e "   ${GREEN}✅ Found orchestrator-state.yaml - This IS an SF instance${NC}"
        return 0
    else
        echo -e "   ${RED}❌ No orchestrator-state.yaml - Not an SF instance${NC}"
        echo -e "   ${YELLOW}Run this script from the Software Factory root directory${NC}"
        return 1
    fi
}

# Function to check target-repo-config.yaml
check_target_config() {
    echo -e "\n${CYAN}2. Checking for target-repo-config.yaml...${NC}"
    
    if [ ! -f "target-repo-config.yaml" ]; then
        echo -e "   ${RED}❌ CRITICAL: target-repo-config.yaml NOT FOUND!${NC}"
        echo -e "   ${YELLOW}This is likely why the orchestrator is confused!${NC}"
        echo ""
        echo -e "   ${BOLD}SOLUTION:${NC}"
        echo -e "   Create target-repo-config.yaml with:"
        echo ""
        cat << 'EOF'
target_repository:
  url: "https://github.com/OWNER/YOUR-ACTUAL-PROJECT.git"
  base_branch: "main"
  
branch_naming:
  project_prefix: ""  # Optional: e.g., "sf-generated"
  effort_format: "{prefix}phase{phase}/wave{wave}/{effort_name}"
  
workspace:
  efforts_root: "efforts"
  effort_path: "phase{phase}/wave{wave}/{effort_name}"
EOF
        return 1
    else
        echo -e "   ${GREEN}✅ Found target-repo-config.yaml${NC}"
        return 0
    fi
}

# Function to validate target URL
validate_target_url() {
    echo -e "\n${CYAN}3. Validating target repository URL...${NC}"
    
    if [ ! -f "target-repo-config.yaml" ]; then
        echo -e "   ${YELLOW}⚠ Skipping - no config file${NC}"
        return 1
    fi
    
    TARGET_URL=$(yq '.target_repository.url' target-repo-config.yaml 2>/dev/null || echo "")
    
    if [ -z "$TARGET_URL" ] || [ "$TARGET_URL" = "null" ]; then
        echo -e "   ${RED}❌ No target URL configured!${NC}"
        echo -e "   ${YELLOW}Edit target-repo-config.yaml and add the repository URL${NC}"
        return 1
    fi
    
    echo -e "   Target URL: ${BOLD}$TARGET_URL${NC}"
    
    # Check if it's a placeholder
    if [[ "$TARGET_URL" == *"OWNER"* ]] || [[ "$TARGET_URL" == *"REPO"* ]] || [[ "$TARGET_URL" == *"your-"* ]]; then
        echo -e "   ${RED}❌ URL looks like a placeholder!${NC}"
        echo -e "   ${YELLOW}Replace with your actual project repository${NC}"
        return 1
    fi
    
    # Check if URL format is valid
    if [[ ! "$TARGET_URL" =~ ^(https://|git@|ssh://) ]]; then
        echo -e "   ${RED}❌ Invalid URL format!${NC}"
        echo -e "   ${YELLOW}URL should start with https://, git@, or ssh://${NC}"
        return 1
    fi
    
    echo -e "   ${GREEN}✅ URL format looks valid${NC}"
    return 0
}

# Function to check for self-reference
check_self_reference() {
    echo -e "\n${CYAN}4. Checking for self-reference...${NC}"
    
    if [ ! -f "target-repo-config.yaml" ]; then
        echo -e "   ${YELLOW}⚠ Skipping - no config file${NC}"
        return 1
    fi
    
    TARGET_URL=$(yq '.target_repository.url' target-repo-config.yaml 2>/dev/null || echo "")
    
    # Try to get SF instance remote URL
    if git rev-parse --git-dir > /dev/null 2>&1; then
        SF_URL=$(git remote get-url origin 2>/dev/null || echo "")
        
        if [ -n "$SF_URL" ] && [ -n "$TARGET_URL" ]; then
            echo -e "   SF Instance URL: ${BOLD}$SF_URL${NC}"
            echo -e "   Target Repo URL: ${BOLD}$TARGET_URL${NC}"
            
            # Normalize URLs for comparison (remove .git, trailing slashes)
            SF_NORMALIZED="${SF_URL%.git}"
            SF_NORMALIZED="${SF_NORMALIZED%/}"
            TARGET_NORMALIZED="${TARGET_URL%.git}"
            TARGET_NORMALIZED="${TARGET_NORMALIZED%/}"
            
            if [ "$SF_NORMALIZED" = "$TARGET_NORMALIZED" ]; then
                echo -e "   ${RED}❌ CRITICAL: Target repo is same as SF instance!${NC}"
                echo -e "   ${YELLOW}The target should be your PROJECT, not Software Factory!${NC}"
                echo ""
                echo -e "   ${BOLD}SOLUTION:${NC}"
                echo -e "   Edit target-repo-config.yaml and change URL to your actual project"
                return 1
            else
                echo -e "   ${GREEN}✅ Target is different from SF instance (correct!)${NC}"
            fi
        fi
    else
        echo -e "   ${YELLOW}⚠ Not a git repository - cannot check self-reference${NC}"
    fi
    
    return 0
}

# Function to test repository access
test_repo_access() {
    echo -e "\n${CYAN}5. Testing target repository accessibility...${NC}"
    
    if [ ! -f "target-repo-config.yaml" ]; then
        echo -e "   ${YELLOW}⚠ Skipping - no config file${NC}"
        return 1
    fi
    
    TARGET_URL=$(yq '.target_repository.url' target-repo-config.yaml 2>/dev/null || echo "")
    
    if [ -z "$TARGET_URL" ] || [ "$TARGET_URL" = "null" ]; then
        echo -e "   ${YELLOW}⚠ No URL to test${NC}"
        return 1
    fi
    
    echo -e "   Testing: $TARGET_URL"
    
    if git ls-remote "$TARGET_URL" HEAD > /dev/null 2>&1; then
        echo -e "   ${GREEN}✅ Repository is accessible${NC}"
        
        # Check base branch
        BASE_BRANCH=$(yq '.target_repository.base_branch' target-repo-config.yaml 2>/dev/null || echo "main")
        if git ls-remote "$TARGET_URL" "refs/heads/$BASE_BRANCH" > /dev/null 2>&1; then
            echo -e "   ${GREEN}✅ Base branch '$BASE_BRANCH' exists${NC}"
        else
            echo -e "   ${RED}❌ Base branch '$BASE_BRANCH' not found!${NC}"
            echo -e "   ${YELLOW}Check if branch name is correct in config${NC}"
        fi
    else
        echo -e "   ${RED}❌ Cannot access repository!${NC}"
        echo -e "   ${YELLOW}Possible causes:${NC}"
        echo -e "   - Repository doesn't exist"
        echo -e "   - No network connection"
        echo -e "   - Authentication required (run: gh auth login)"
        echo -e "   - Private repository without access"
    fi
}

# Function to check existing efforts
check_existing_efforts() {
    echo -e "\n${CYAN}6. Checking for existing effort clones...${NC}"
    
    if [ -d "efforts" ]; then
        echo -e "   ${GREEN}✅ Found efforts directory${NC}"
        
        # Count effort directories
        EFFORT_COUNT=$(find efforts -type d -name ".git" 2>/dev/null | wc -l)
        
        if [ "$EFFORT_COUNT" -gt 0 ]; then
            echo -e "   Found ${BOLD}$EFFORT_COUNT${NC} cloned efforts"
            
            # Check what they're clones of
            echo -e "\n   Checking what's been cloned:"
            find efforts -type d -name ".git" -exec dirname {} \; 2>/dev/null | while read -r effort_dir; do
                if [ -d "$effort_dir/.git" ]; then
                    REMOTE_URL=$(cd "$effort_dir" && git remote get-url origin 2>/dev/null || echo "unknown")
                    BRANCH=$(cd "$effort_dir" && git branch --show-current 2>/dev/null || echo "unknown")
                    echo -e "   - ${BOLD}$(basename "$effort_dir")${NC}"
                    echo -e "     Repo: $REMOTE_URL"
                    echo -e "     Branch: $BRANCH"
                fi
            done
        else
            echo -e "   ${YELLOW}No effort clones found yet${NC}"
        fi
    else
        echo -e "   ${YELLOW}No efforts directory found${NC}"
    fi
}

# Function to show configuration summary
show_summary() {
    echo -e "\n${MAGENTA}${BOLD}════════════════════════════════════════════════════════${NC}"
    echo -e "${MAGENTA}${BOLD}   DIAGNOSTIC SUMMARY${NC}"
    echo -e "${MAGENTA}${BOLD}════════════════════════════════════════════════════════${NC}\n"
    
    if [ -f "target-repo-config.yaml" ]; then
        TARGET_URL=$(yq '.target_repository.url' target-repo-config.yaml 2>/dev/null || echo "NOT SET")
        BASE_BRANCH=$(yq '.target_repository.base_branch' target-repo-config.yaml 2>/dev/null || echo "NOT SET")
        PREFIX=$(yq '.branch_naming.project_prefix' target-repo-config.yaml 2>/dev/null || echo "(none)")
        
        echo -e "${BOLD}Current Configuration:${NC}"
        echo -e "  Target URL: $TARGET_URL"
        echo -e "  Base Branch: $BASE_BRANCH"
        echo -e "  Branch Prefix: $PREFIX"
    else
        echo -e "${RED}${BOLD}NO CONFIGURATION FILE FOUND!${NC}"
        echo -e "This is why the orchestrator is confused."
    fi
    
    echo -e "\n${BOLD}Required Actions:${NC}"
    
    if [ ! -f "target-repo-config.yaml" ]; then
        echo -e "  1. ${RED}Create target-repo-config.yaml with your actual project URL${NC}"
    elif [[ "$(yq '.target_repository.url' target-repo-config.yaml 2>/dev/null)" == *"OWNER"* ]]; then
        echo -e "  1. ${RED}Replace placeholder URL with actual project repository${NC}"
    else
        echo -e "  1. ${GREEN}Configuration appears correct${NC}"
    fi
    
    echo -e "\n${BOLD}Quick Fix Command:${NC}"
    echo -e "${CYAN}cat > target-repo-config.yaml << 'EOF'"
    echo "target_repository:"
    echo "  url: \"https://github.com/YOUR-ORG/YOUR-PROJECT.git\"  # <-- CHANGE THIS!"
    echo "  base_branch: \"main\""
    echo ""
    echo "branch_naming:"
    echo "  project_prefix: \"\"  # Optional prefix for branches"
    echo "  effort_format: \"{prefix}phase{phase}/wave{wave}/{effort_name}\""
    echo ""
    echo "workspace:"
    echo "  efforts_root: \"efforts\""
    echo "  effort_path: \"phase{phase}/wave{wave}/{effort_name}\""
    echo "EOF${NC}"
}

# Main execution
main() {
    local has_errors=0
    
    check_sf_instance || has_errors=1
    
    if [ $has_errors -eq 0 ]; then
        check_target_config || has_errors=1
        validate_target_url || has_errors=1
        check_self_reference || has_errors=1
        test_repo_access || has_errors=1
        check_existing_efforts
    fi
    
    show_summary
    
    if [ $has_errors -eq 0 ]; then
        echo -e "\n${GREEN}${BOLD}✅ Target repository configuration looks good!${NC}"
        exit 0
    else
        echo -e "\n${RED}${BOLD}❌ Configuration issues detected - see above for fixes${NC}"
        exit 1
    fi
}

# Run main
main