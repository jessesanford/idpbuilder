#!/bin/bash
# Verify architect reports are in correct locations
# Per R257 and R258 requirements

echo "========================================="
echo "ARCHITECT REPORT LOCATION VERIFICATION"
echo "========================================="
echo ""

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Counters
CORRECT_PHASE_REPORTS=0
WRONG_PHASE_REPORTS=0
CORRECT_WAVE_REPORTS=0
WRONG_WAVE_REPORTS=0

echo "Checking Phase Assessment Reports (R257)..."
echo "Expected location: phase-assessments/phase{N}/PHASE-{N}-ASSESSMENT-REPORT.md"
echo ""

# Find all phase assessment reports
find . -name "PHASE-*-ASSESSMENT-REPORT.md" -type f 2>/dev/null | grep -v ".git" | while read report; do
    # Remove leading ./
    report_path="${report#./}"
    
    # Check if it matches the correct pattern
    if [[ "$report_path" =~ ^phase-assessments/phase[0-9]+/PHASE-[0-9]+-ASSESSMENT-REPORT\.md$ ]]; then
        echo -e "${GREEN}✅ CORRECT${NC}: $report_path"
        ((CORRECT_PHASE_REPORTS++))
    else
        echo -e "${RED}❌ WRONG LOCATION${NC}: $report_path"
        # Extract phase number from filename
        phase_num=$(echo "$report_path" | grep -o "PHASE-[0-9]*" | grep -o "[0-9]*")
        echo -e "   ${YELLOW}Should be at${NC}: phase-assessments/phase${phase_num}/$(basename "$report_path")"
        ((WRONG_PHASE_REPORTS++))
    fi
done

echo ""
echo "Checking Wave Review Reports (R258)..."
echo "Expected location: wave-reviews/phase{N}/wave{W}/PHASE-{N}-WAVE-{W}-REVIEW-REPORT.md"
echo ""

# Find all wave review reports
find . -name "PHASE-*-WAVE-*-REVIEW-REPORT.md" -type f 2>/dev/null | grep -v ".git" | while read report; do
    # Remove leading ./
    report_path="${report#./}"
    
    # Check if it matches the correct pattern
    if [[ "$report_path" =~ ^wave-reviews/phase[0-9]+/wave[0-9]+/PHASE-[0-9]+-WAVE-[0-9]+-REVIEW-REPORT\.md$ ]]; then
        echo -e "${GREEN}✅ CORRECT${NC}: $report_path"
        ((CORRECT_WAVE_REPORTS++))
    else
        echo -e "${RED}❌ WRONG LOCATION${NC}: $report_path"
        # Extract phase and wave numbers from filename
        phase_num=$(echo "$report_path" | grep -o "PHASE-[0-9]*" | grep -o "[0-9]*")
        wave_num=$(echo "$report_path" | grep -o "WAVE-[0-9]*" | grep -o "[0-9]*")
        echo -e "   ${YELLOW}Should be at${NC}: wave-reviews/phase${phase_num}/wave${wave_num}/$(basename "$report_path")"
        ((WRONG_WAVE_REPORTS++))
    fi
done

echo ""
echo "========================================="
echo "SUMMARY"
echo "========================================="

# Check for reports in root directory (common mistake)
ROOT_REPORTS=$(ls -1 *ASSESSMENT-REPORT.md *REVIEW-REPORT.md 2>/dev/null | wc -l)
if [ "$ROOT_REPORTS" -gt 0 ]; then
    echo -e "${RED}⚠️ WARNING${NC}: Found $ROOT_REPORTS report(s) in root directory!"
    echo "These should be moved to their proper locations:"
    ls -1 *ASSESSMENT-REPORT.md *REVIEW-REPORT.md 2>/dev/null | while read report; do
        echo "  - $report"
    done
fi

# Function to verify specific report location
verify_phase_report() {
    local PHASE=$1
    local EXPECTED="phase-assessments/phase${PHASE}/PHASE-${PHASE}-ASSESSMENT-REPORT.md"
    
    if [[ -f "$EXPECTED" ]]; then
        echo -e "${GREEN}✅ Phase $PHASE assessment report in correct location${NC}"
        return 0
    else
        echo -e "${RED}❌ Phase $PHASE assessment report NOT FOUND at expected location${NC}"
        echo "   Expected: $EXPECTED"
        
        # Check if it exists elsewhere
        local FOUND=$(find . -name "PHASE-${PHASE}-ASSESSMENT-REPORT.md" -type f 2>/dev/null | grep -v ".git")
        if [ -n "$FOUND" ]; then
            echo -e "   ${YELLOW}Found at${NC}: $FOUND"
            echo -e "   ${YELLOW}Move it to${NC}: $EXPECTED"
        fi
        return 1
    fi
}

verify_wave_report() {
    local PHASE=$1
    local WAVE=$2
    local EXPECTED="wave-reviews/phase${PHASE}/wave${WAVE}/PHASE-${PHASE}-WAVE-${WAVE}-REVIEW-REPORT.md"
    
    if [[ -f "$EXPECTED" ]]; then
        echo -e "${GREEN}✅ Phase $PHASE Wave $WAVE review report in correct location${NC}"
        return 0
    else
        echo -e "${RED}❌ Phase $PHASE Wave $WAVE review report NOT FOUND at expected location${NC}"
        echo "   Expected: $EXPECTED"
        
        # Check if it exists elsewhere
        local FOUND=$(find . -name "PHASE-${PHASE}-WAVE-${WAVE}-REVIEW-REPORT.md" -type f 2>/dev/null | grep -v ".git")
        if [ -n "$FOUND" ]; then
            echo -e "   ${YELLOW}Found at${NC}: $FOUND"
            echo -e "   ${YELLOW}Move it to${NC}: $EXPECTED"
        fi
        return 1
    fi
}

# If specific phase/wave provided as arguments, verify those
if [ $# -eq 1 ]; then
    echo ""
    echo "Verifying Phase $1 Assessment Report..."
    verify_phase_report $1
elif [ $# -eq 2 ]; then
    echo ""
    echo "Verifying Phase $1 Wave $2 Review Report..."
    verify_wave_report $1 $2
fi

echo ""
echo "========================================="
echo "ENFORCEMENT REMINDERS"
echo "========================================="
echo ""
echo "Per R257 (Phase Assessment Reports):"
echo "  - MUST be at: phase-assessments/phase{N}/PHASE-{N}-ASSESSMENT-REPORT.md"
echo "  - Wrong location = -50% grading penalty"
echo "  - Orchestrator cannot proceed without correct location"
echo ""
echo "Per R258 (Wave Review Reports):"
echo "  - MUST be at: wave-reviews/phase{N}/wave{W}/PHASE-{N}-WAVE-{W}-REVIEW-REPORT.md"
echo "  - Wrong location = -50% grading penalty"
echo "  - Orchestrator cannot proceed without correct location"
echo ""
echo "Architect MUST run verification after creating reports!"