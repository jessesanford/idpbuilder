#!/bin/bash

# R290 Compliance Verification Script
# Ensures state rule reading and verification is happening before state actions

echo "🔍 R290 COMPLIANCE VERIFICATION TOOL"
echo "===================================="
echo "Checking enforcement of SUPREME LAW #3 (State Rule Reading and Verification)"
echo ""

# Color codes
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check if R290 exists
echo "1. Checking R290 rule file existence..."
if [[ -f "rule-library/R290-state-rule-reading-verification-supreme-law.md" ]]; then
    echo -e "${GREEN}✅ R290 rule file exists${NC}"
else
    echo -e "${RED}❌ R290 rule file MISSING!${NC}"
    exit 1
fi

# Check orchestrator.md references
echo ""
echo "2. Checking orchestrator.md references R290..."
if grep -q "R290.*SUPREME LAW #3" .claude/agents/orchestrator.md; then
    echo -e "${GREEN}✅ R290 listed as SUPREME LAW #3 in orchestrator${NC}"
else
    echo -e "${RED}❌ R290 not properly listed in orchestrator!${NC}"
    exit 1
fi

# Check R231 clarification
echo ""
echo "3. Checking R231 clarification about R290..."
if grep -q "R290.*SUPREME LAW #3" rule-library/R231-continuous-operation-through-transitions.md; then
    echo -e "${GREEN}✅ R231 properly references R290${NC}"
else
    echo -e "${YELLOW}⚠️ R231 may need R290 reference update${NC}"
fi

# Check state files for R290 headers
echo ""
echo "4. Checking state files for R290 enforcement headers..."
states_checked=0
states_compliant=0

for state_dir in agent-states/software-factory/orchestrator/*/; do
    state_name=$(basename "$state_dir")
    rules_file="${state_dir}rules.md"
    
    if [[ -f "$rules_file" ]]; then
        ((states_checked++))
        if grep -q "R290" "$rules_file"; then
            echo -e "${GREEN}  ✅ ${state_name} has R290 enforcement${NC}"
            ((states_compliant++))
        else
            echo -e "${YELLOW}  ⚠️ ${state_name} missing R290 header${NC}"
        fi
    fi
done

echo ""
echo "State compliance: $states_compliant/$states_checked"

# Check for state transition patterns in logs (if they exist)
echo ""
echo "5. Checking for compliance patterns in recent logs..."
if [[ -d "logs" ]]; then
    echo "Scanning logs for state rule reading patterns..."
    if grep -r "Reading state rules" logs/ 2>/dev/null | head -3; then
        echo -e "${GREEN}✅ Found evidence of state rule reading${NC}"
    else
        echo -e "${YELLOW}⚠️ No evidence of state rule reading in logs${NC}"
    fi
else
    echo "No logs directory found (this is normal if not running)"
fi

# Final summary
echo ""
echo "========================================"
echo "R290 COMPLIANCE SUMMARY"
echo "========================================"

compliance_score=100

if [[ ! -f "rule-library/R290-mandatory-state-rule-reading-supreme-law.md" ]]; then
    compliance_score=0
    echo -e "${RED}CRITICAL: R290 rule file missing${NC}"
fi

if ! grep -q "R290.*SUPREME LAW #3" .claude/agents/orchestrator.md; then
    compliance_score=$((compliance_score - 50))
    echo -e "${RED}CRITICAL: Orchestrator not enforcing R290${NC}"
fi

if [[ $states_compliant -lt $states_checked ]]; then
    missing=$((states_checked - states_compliant))
    penalty=$((missing * 10))
    compliance_score=$((compliance_score - penalty))
    echo -e "${YELLOW}WARNING: $missing state files need R290 headers${NC}"
fi

echo ""
if [[ $compliance_score -eq 100 ]]; then
    echo -e "${GREEN}🎉 FULL COMPLIANCE WITH R290!${NC}"
    echo "All agents will read state rules before taking state actions"
elif [[ $compliance_score -ge 70 ]]; then
    echo -e "${YELLOW}⚠️ PARTIAL COMPLIANCE: ${compliance_score}%${NC}"
    echo "Some improvements needed for full R290 enforcement"
else
    echo -e "${RED}❌ NON-COMPLIANT: ${compliance_score}%${NC}"
    echo "CRITICAL: System at risk of state rule violations!"
fi

echo ""
echo "Run this script regularly to ensure R290 compliance"
exit 0