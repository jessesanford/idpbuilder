#!/bin/bash

# validate-effort-isolation.sh
# Validates R209 effort directory isolation protocol compliance
# Checks implementation plans for required metadata and validates SW engineer isolation

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Counters
TOTAL_EFFORTS=0
VALID_EFFORTS=0
INVALID_EFFORTS=0
WARNINGS=0

echo "════════════════════════════════════════════════════════════════"
echo "🔍 R209: Effort Directory Isolation Validation"
echo "════════════════════════════════════════════════════════════════"

# Function to validate implementation plan metadata
validate_implementation_plan() {
    local PLAN="$1"
    local EFFORT_DIR=$(dirname "$PLAN")
    local EFFORT_NAME=$(basename "$EFFORT_DIR")
    local ERRORS=0
    
    echo ""
    echo "Checking: $PLAN"
    echo "─────────────────────────────────────────────────────────"
    
    # Check if plan exists
    if [ ! -f "$PLAN" ]; then
        echo -e "${RED}❌ FAIL: Implementation plan not found${NC}"
        return 1
    fi
    
    # Required metadata fields per R209
    local REQUIRED_FIELDS=(
        "EFFORT INFRASTRUCTURE METADATA"
        "WORKING_DIRECTORY"
        "BRANCH"
        "ISOLATION_BOUNDARY"
        "EFFORT_NAME"
        "PHASE"
        "WAVE"
        "CRITICAL ISOLATION RULES"
    )
    
    # Check each required field
    for field in "${REQUIRED_FIELDS[@]}"; do
        if grep -q "$field" "$PLAN"; then
            echo -e "${GREEN}✅ Found: $field${NC}"
        else
            echo -e "${RED}❌ Missing: $field${NC}"
            ((ERRORS++))
        fi
    done
    
    # Extract and validate metadata values
    if grep -q "WORKING_DIRECTORY" "$PLAN"; then
        WORKING_DIR=$(grep "**WORKING_DIRECTORY**:" "$PLAN" 2>/dev/null | cut -d: -f2- | xargs || true)
        if [ -n "$WORKING_DIR" ]; then
            # Check if working directory matches effort directory
            if [[ "$WORKING_DIR" == *"$EFFORT_DIR"* ]] || [[ "$EFFORT_DIR" == *"$(basename $WORKING_DIR)"* ]]; then
                echo -e "${GREEN}✅ Working directory correctly references effort${NC}"
            else
                echo -e "${RED}❌ Working directory mismatch!${NC}"
                echo "   Expected pattern: */$EFFORT_NAME"
                echo "   Found: $WORKING_DIR"
                ((ERRORS++))
            fi
        fi
    fi
    
    if grep -q "BRANCH" "$PLAN"; then
        BRANCH=$(grep "**BRANCH**:" "$PLAN" 2>/dev/null | cut -d: -f2- | xargs || true)
        if [ -n "$BRANCH" ]; then
            # Extract phase/wave from path
            if [[ "$EFFORT_DIR" =~ phase([0-9]+)/wave([0-9]+) ]]; then
                PHASE="${BASH_REMATCH[1]}"
                WAVE="${BASH_REMATCH[2]}"
                EXPECTED_BRANCH="phase${PHASE}/wave${WAVE}/${EFFORT_NAME}"
                
                if [[ "$BRANCH" == "$EXPECTED_BRANCH"* ]]; then
                    echo -e "${GREEN}✅ Branch naming follows convention${NC}"
                else
                    echo -e "${YELLOW}⚠️  Branch naming may not follow convention${NC}"
                    echo "   Expected: $EXPECTED_BRANCH"
                    echo "   Found: $BRANCH"
                    ((WARNINGS++))
                fi
            fi
        fi
    fi
    
    # Check for isolation rules section
    if grep -q "CRITICAL ISOLATION RULES" "$PLAN"; then
        # Check for the 5 critical rules
        ISOLATION_RULES=(
            "ALL work MUST happen in:"
            "ALL code MUST go in:"
            "NEVER cd out of this directory"
            "NEVER create files outside"
            "This is YOUR isolated workspace"
        )
        
        local RULES_FOUND=0
        for rule in "${ISOLATION_RULES[@]}"; do
            if grep -q "$rule" "$PLAN"; then
                ((RULES_FOUND++))
            fi
        done
        
        if [ $RULES_FOUND -eq 5 ]; then
            echo -e "${GREEN}✅ All 5 isolation rules present${NC}"
        else
            echo -e "${YELLOW}⚠️  Only $RULES_FOUND/5 isolation rules found${NC}"
            ((WARNINGS++))
        fi
    fi
    
    return $ERRORS
}

# Function to check SW engineer work logs for isolation violations
check_work_log() {
    local WORK_LOG="$1"
    local EFFORT_DIR=$(dirname "$WORK_LOG")
    local VIOLATIONS=0
    
    if [ ! -f "$WORK_LOG" ]; then
        return 0  # No work log yet is OK
    fi
    
    echo ""
    echo "Checking work log: $WORK_LOG"
    
    # Look for dangerous patterns
    DANGEROUS_PATTERNS=(
        "cd \.\."
        "cd /"
        "cd ~"
        "../.."
        "cd /workspaces"
        "cd /home"
    )
    
    for pattern in "${DANGEROUS_PATTERNS[@]}"; do
        if grep -q "$pattern" "$WORK_LOG" 2>/dev/null; then
            echo -e "${RED}❌ VIOLATION: Found '$pattern' in work log!${NC}"
            ((VIOLATIONS++))
        fi
    done
    
    # Check for files created outside pkg/
    if grep -E "touch |echo.*>|cat.*>" "$WORK_LOG" 2>/dev/null | grep -v "pkg/" | grep -v "work-log" | grep -v "IMPLEMENTATION-PLAN" > /dev/null; then
        echo -e "${YELLOW}⚠️  Warning: May have created files outside pkg/${NC}"
        ((WARNINGS++))
    fi
    
    if [ $VIOLATIONS -eq 0 ]; then
        echo -e "${GREEN}✅ No isolation violations detected in work log${NC}"
    fi
    
    return $VIOLATIONS
}

# Function to validate effort directory structure
validate_effort_structure() {
    local EFFORT_DIR="$1"
    local ERRORS=0
    
    echo ""
    echo "Validating structure: $EFFORT_DIR"
    
    # Check required directories
    if [ -d "$EFFORT_DIR/pkg" ]; then
        echo -e "${GREEN}✅ pkg/ directory exists${NC}"
    else
        echo -e "${YELLOW}⚠️  pkg/ directory not yet created${NC}"
    fi
    
    # Check for R209 acknowledgment file
    if [ -f "$EFFORT_DIR/.r209-acknowledged" ]; then
        echo -e "${GREEN}✅ R209 acknowledgment file found${NC}"
        echo "   Last acknowledgment:"
        tail -3 "$EFFORT_DIR/.r209-acknowledged" | sed 's/^/   /'
    else
        echo -e "${YELLOW}⚠️  No R209 acknowledgment file - SW Engineer may not have started work yet${NC}"
    fi
    
    # Check for files outside allowed locations
    local FILES_OUTSIDE=$(find "$EFFORT_DIR" -maxdepth 1 -type f \
        ! -name "IMPLEMENTATION-PLAN.md" \
        ! -name "work-log.md" \
        ! -name "REVIEW-FEEDBACK.md" \
        ! -name "SPLIT-*.md" \
        ! -name ".gitignore" \
        ! -name ".r209-acknowledged" 2>/dev/null | wc -l)
    
    if [ $FILES_OUTSIDE -gt 0 ]; then
        echo -e "${YELLOW}⚠️  Found $FILES_OUTSIDE unexpected files in effort root${NC}"
        find "$EFFORT_DIR" -maxdepth 1 -type f \
            ! -name "IMPLEMENTATION-PLAN.md" \
            ! -name "work-log.md" \
            ! -name "REVIEW-FEEDBACK.md" \
            ! -name "SPLIT-*.md" \
            ! -name ".gitignore" \
            ! -name ".r209-acknowledged" 2>/dev/null | while read file; do
            echo "   - $(basename $file)"
        done
        ((WARNINGS++))
    fi
    
    return $ERRORS
}

# Main validation loop
echo ""
echo "Scanning for effort directories..."
echo "════════════════════════════════════════════════════════════════"

# Find all effort directories
for PHASE_DIR in efforts/phase*/; do
    if [ ! -d "$PHASE_DIR" ]; then
        continue
    fi
    
    for WAVE_DIR in ${PHASE_DIR}wave*/; do
        if [ ! -d "$WAVE_DIR" ]; then
            continue
        fi
        
        for EFFORT_DIR in ${WAVE_DIR}*/; do
            if [ ! -d "$EFFORT_DIR" ]; then
                continue
            fi
            
            # Skip split directories
            if [[ "$(basename $EFFORT_DIR)" == split-* ]]; then
                continue
            fi
            
            ((TOTAL_EFFORTS++))
            
            echo ""
            echo "════════════════════════════════════════════════════════════════"
            echo "📁 Effort: $EFFORT_DIR"
            echo "════════════════════════════════════════════════════════════════"
            
            # Validate implementation plan
            PLAN_ERRORS=0
            if validate_implementation_plan "${EFFORT_DIR}IMPLEMENTATION-PLAN.md"; then
                :  # Success
            else
                PLAN_ERRORS=$?
            fi
            
            # Check work log for violations
            WORK_LOG_VIOLATIONS=0
            if check_work_log "${EFFORT_DIR}work-log.md"; then
                :  # Success
            else
                WORK_LOG_VIOLATIONS=$?
            fi
            
            # Validate directory structure
            STRUCTURE_ERRORS=0
            if validate_effort_structure "$EFFORT_DIR"; then
                :  # Success
            else
                STRUCTURE_ERRORS=$?
            fi
            
            # Tally results
            TOTAL_ERRORS=$((PLAN_ERRORS + WORK_LOG_VIOLATIONS + STRUCTURE_ERRORS))
            
            if [ $TOTAL_ERRORS -eq 0 ]; then
                echo ""
                echo -e "${GREEN}✅ EFFORT VALID: All R209 requirements met${NC}"
                ((VALID_EFFORTS++))
            else
                echo ""
                echo -e "${RED}❌ EFFORT INVALID: $TOTAL_ERRORS errors found${NC}"
                ((INVALID_EFFORTS++))
            fi
        done
    done
done

# Summary report
echo ""
echo "════════════════════════════════════════════════════════════════"
echo "📊 R209 Validation Summary"
echo "════════════════════════════════════════════════════════════════"
echo "Total Efforts Scanned: $TOTAL_EFFORTS"
echo -e "${GREEN}Valid Efforts: $VALID_EFFORTS${NC}"
echo -e "${RED}Invalid Efforts: $INVALID_EFFORTS${NC}"
echo -e "${YELLOW}Warnings: $WARNINGS${NC}"
echo ""

if [ $INVALID_EFFORTS -eq 0 ] && [ $WARNINGS -eq 0 ]; then
    echo -e "${GREEN}🎉 PROJECT_DONE: All efforts comply with R209 isolation protocol!${NC}"
    exit 0
elif [ $INVALID_EFFORTS -eq 0 ]; then
    echo -e "${YELLOW}⚠️  MOSTLY COMPLIANT: No errors but $WARNINGS warnings found${NC}"
    exit 0
else
    echo -e "${RED}❌ VALIDATION FAILED: $INVALID_EFFORTS efforts violate R209${NC}"
    echo ""
    echo "Required Actions:"
    echo "1. Orchestrator must inject metadata into all implementation plans"
    echo "2. SW Engineers must read metadata and stay in effort directories"
    echo "3. All code must go under pkg/ directory"
    echo "4. No navigation outside effort boundaries"
    echo ""
    echo "Run orchestrator with R209 enforcement to fix these issues."
    exit 1
fi