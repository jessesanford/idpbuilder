#!/bin/bash

# R300 Compliance Verification Script
# Verifies that all fixes are properly applied to effort branches

set -euo pipefail

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo "🔍 R300 COMPLIANCE VERIFICATION SCRIPT"
echo "======================================"
echo "Checking that all fixes are in effort branches..."
echo ""

# Function to check for fixes in integration branches
check_integration_branches() {
    echo "Checking for forbidden fixes in integration branches..."
    
    # Find all integration branches
    INTEGRATION_BRANCHES=$(git branch -a | grep -E "integration-|integration_" || true)
    
    if [ -z "$INTEGRATION_BRANCHES" ]; then
        echo -e "${GREEN}✅ No integration branches found${NC}"
        return 0
    fi
    
    VIOLATIONS_FOUND=false
    
    for branch in $INTEGRATION_BRANCHES; do
        # Clean branch name
        branch=$(echo "$branch" | sed 's/^[* ]*//' | sed 's/remotes\///')
        
        # Check for fix commits in integration branch
        FIX_COMMITS=$(git log "$branch" --oneline --grep="^fix:" 2>/dev/null || true)
        
        if [ -n "$FIX_COMMITS" ]; then
            echo -e "${RED}❌ R300 VIOLATION: Fixes found in integration branch $branch!${NC}"
            echo "$FIX_COMMITS" | head -5
            VIOLATIONS_FOUND=true
        fi
    done
    
    if [ "$VIOLATIONS_FOUND" = false ]; then
        echo -e "${GREEN}✅ No fixes in integration branches - R300 compliant${NC}"
    fi
    
    return $([ "$VIOLATIONS_FOUND" = false ])
}

# Function to verify fixes are in effort branches
verify_effort_branch_fixes() {
    echo ""
    echo "Verifying fixes are properly in effort branches..."
    
    # Check orchestrator state for efforts that had fixes
    if [ -f "orchestrator-state.json" ]; then
        EFFORTS_WITH_FIXES=$(jq '.efforts_with_fixes[]' orchestrator-state.json 2>/dev/null || true)
        
        if [ -z "$EFFORTS_WITH_FIXES" ]; then
            echo -e "${YELLOW}⚠️ No efforts marked as having fixes in state file${NC}"
            return 0
        fi
        
        ALL_VERIFIED=true
        
        for effort in $EFFORTS_WITH_FIXES; do
            echo "Checking effort: $effort"
            
            # Find effort branch
            EFFORT_BRANCH=$(git branch -a | grep -E "effort-$effort|feature/$effort" | head -1 || true)
            
            if [ -z "$EFFORT_BRANCH" ]; then
                echo -e "${RED}❌ Cannot find branch for effort $effort${NC}"
                ALL_VERIFIED=false
                continue
            fi
            
            # Clean branch name
            EFFORT_BRANCH=$(echo "$EFFORT_BRANCH" | sed 's/^[* ]*//' | sed 's/remotes\///')
            
            # Check for recent fix commits
            RECENT_FIXES=$(git log "$EFFORT_BRANCH" --oneline --grep="^fix:" --since="24 hours ago" 2>/dev/null || true)
            
            if [ -z "$RECENT_FIXES" ]; then
                echo -e "${RED}❌ No recent fixes found in $EFFORT_BRANCH${NC}"
                ALL_VERIFIED=false
            else
                echo -e "${GREEN}✅ Fixes found in $EFFORT_BRANCH${NC}"
                echo "$RECENT_FIXES" | head -3
            fi
            
            # Verify pushed to remote
            LOCAL_SHA=$(git rev-parse "$EFFORT_BRANCH" 2>/dev/null || echo "none")
            REMOTE_SHA=$(git rev-parse "origin/$EFFORT_BRANCH" 2>/dev/null || echo "none")
            
            if [ "$LOCAL_SHA" != "$REMOTE_SHA" ]; then
                echo -e "${RED}❌ Branch $EFFORT_BRANCH not pushed to remote!${NC}"
                ALL_VERIFIED=false
            fi
        done
        
        if [ "$ALL_VERIFIED" = true ]; then
            echo -e "${GREEN}✅ All effort branches have fixes and are pushed - R300 compliant${NC}"
        else
            echo -e "${RED}❌ Some effort branches missing fixes or not pushed - R300 violation${NC}"
            return 1
        fi
    else
        echo -e "${YELLOW}⚠️ No orchestrator-state.json found${NC}"
    fi
}

# Function to check for R300 violations in code
check_code_for_violations() {
    echo ""
    echo "Checking for R300 violations in code..."
    
    # Look for patterns that indicate fixing in wrong place
    VIOLATION_PATTERNS=(
        "cd.*integration.*&&.*vim"
        "cd.*integration.*&&.*edit"
        "checkout.*integration.*&&.*fix"
        "git commit.*integration.*fix"
    )
    
    VIOLATIONS_FOUND=false
    
    for pattern in "${VIOLATION_PATTERNS[@]}"; do
        if grep -r "$pattern" agent-states/ 2>/dev/null | grep -v "WRONG\|FORBIDDEN\|VIOLATION" > /dev/null; then
            echo -e "${RED}❌ Found potential R300 violation pattern: $pattern${NC}"
            VIOLATIONS_FOUND=true
        fi
    done
    
    if [ "$VIOLATIONS_FOUND" = false ]; then
        echo -e "${GREEN}✅ No R300 violation patterns found in code${NC}"
    fi
}

# Function to verify R300 is referenced correctly
check_r300_references() {
    echo ""
    echo "Verifying R300 is properly referenced..."
    
    # Check that deprecated rules are not used (except in deprecated files)
    DEPRECATED_RULES="R299 R240 R292 R298"
    
    for rule in $DEPRECATED_RULES; do
        # Find references excluding rule-library and git directories
        REFERENCES=$(grep -r "$rule" --exclude-dir=.git --exclude-dir=rule-library agent-states/ .claude/agents/ 2>/dev/null || true)
        
        if [ -n "$REFERENCES" ]; then
            echo -e "${RED}❌ Found references to deprecated rule $rule that should use R300:${NC}"
            echo "$REFERENCES" | head -3
        fi
    done
    
    # Check that R300 is present in key files
    KEY_FILES=(
        "agent-states/orchestrator/ERROR_RECOVERY/rules.md"
        "agent-states/orchestrator/MONITORING_FIX_PROGRESS/rules.md"
        "agent-states/sw-engineer/FIX_ISSUES/rules.md"
        "agent-states/sw-engineer/FIX_INTEGRATION_ISSUES/rules.md"
    )
    
    ALL_PRESENT=true
    for file in "${KEY_FILES[@]}"; do
        if [ -f "$file" ]; then
            if ! grep -q "R300" "$file"; then
                echo -e "${RED}❌ R300 not referenced in $file${NC}"
                ALL_PRESENT=false
            fi
        fi
    done
    
    if [ "$ALL_PRESENT" = true ]; then
        echo -e "${GREEN}✅ R300 properly referenced in all key files${NC}"
    fi
}

# Main execution
main() {
    echo "Starting R300 compliance verification..."
    echo ""
    
    OVERALL_COMPLIANCE=true
    
    # Run all checks
    if ! check_integration_branches; then
        OVERALL_COMPLIANCE=false
    fi
    
    if ! verify_effort_branch_fixes; then
        OVERALL_COMPLIANCE=false
    fi
    
    if ! check_code_for_violations; then
        OVERALL_COMPLIANCE=false
    fi
    
    if ! check_r300_references; then
        OVERALL_COMPLIANCE=false
    fi
    
    echo ""
    echo "======================================"
    
    if [ "$OVERALL_COMPLIANCE" = true ]; then
        echo -e "${GREEN}✅ SYSTEM IS R300 COMPLIANT${NC}"
        echo "All fixes are properly managed in effort branches"
        exit 0
    else
        echo -e "${RED}❌ R300 COMPLIANCE VIOLATIONS DETECTED${NC}"
        echo "Review the violations above and correct them"
        exit 1
    fi
}

# Run main function
main "$@"