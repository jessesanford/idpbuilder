#!/bin/bash

# validate-factory-compliance.sh
# Master validation script for Software Factory 2.0 compliance
# Runs all validation checks and provides comprehensive report

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
MAGENTA='\033[0;35m'
CYAN='\033[0;36m'
BOLD='\033[1m'
NC='\033[0m' # No Color

# Script directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Counters
TOTAL_VALIDATIONS=0
PASSED_VALIDATIONS=0
FAILED_VALIDATIONS=0
WARNINGS=0

echo "════════════════════════════════════════════════════════════════"
echo -e "${BOLD}SOFTWARE FACTORY 2.0 - MASTER COMPLIANCE VALIDATOR${NC}"
echo "════════════════════════════════════════════════════════════════"
echo ""
echo "Running comprehensive validation of all factory rules..."
echo ""

# Function to run a validation script
run_validation() {
    local SCRIPT="$1"
    local NAME="$2"
    local SCRIPT_PATH="$SCRIPT_DIR/$SCRIPT"
    
    echo -e "${CYAN}════════════════════════════════════════════════════════════════${NC}"
    echo -e "${CYAN}Running: $NAME${NC}"
    echo -e "${CYAN}════════════════════════════════════════════════════════════════${NC}"
    
    ((TOTAL_VALIDATIONS++))
    
    if [ ! -f "$SCRIPT_PATH" ]; then
        echo -e "${YELLOW}⚠️  Script not found: $SCRIPT${NC}"
        ((WARNINGS++))
        return 1
    fi
    
    # Make executable
    chmod +x "$SCRIPT_PATH"
    
    # Run and capture exit code
    if "$SCRIPT_PATH"; then
        echo -e "${GREEN}✅ $NAME: PASSED${NC}"
        ((PASSED_VALIDATIONS++))
    else
        echo -e "${RED}❌ $NAME: FAILED${NC}"
        ((FAILED_VALIDATIONS++))
    fi
    
    echo ""
}

# Function to check for rule existence
check_rule_exists() {
    local RULE_FILE="$1"
    local RULE_NAME="$2"
    
    if [ -f "../rule-library/$RULE_FILE" ]; then
        echo -e "${GREEN}✓ $RULE_NAME exists${NC}"
        return 0
    else
        echo -e "${RED}✗ $RULE_NAME missing${NC}"
        return 1
    fi
}

# Check rule library
echo -e "${BLUE}════════════════════════════════════════════════════════════════${NC}"
echo -e "${BLUE}📚 Rule Library Check${NC}"
echo -e "${BLUE}════════════════════════════════════════════════════════════════${NC}"

# Rules are in the rule-library directory
RULE_DIR="../rule-library"
if [ -f "$RULE_DIR/R209-effort-directory-isolation-protocol.md" ]; then
    echo -e "${GREEN}✓ R209 - Effort Directory Isolation exists${NC}"
else
    echo -e "${YELLOW}⚠ R209 - Effort Directory Isolation not yet created${NC}"
fi

if [ -f "$RULE_DIR/R210-architect-architecture-planning-protocol.md" ]; then
    echo -e "${GREEN}✓ R210 - Architect Architecture Planning exists${NC}"
else
    echo -e "${YELLOW}⚠ R210 - Architect Architecture Planning not yet created${NC}"
fi

if [ -f "$RULE_DIR/R211-code-reviewer-implementation-from-architecture.md" ]; then
    echo -e "${GREEN}✓ R211 - Implementation from Architecture exists${NC}"
else
    echo -e "${YELLOW}⚠ R211 - Implementation from Architecture not yet created${NC}"
fi

if [ -f "$RULE_DIR/R212-phase-directory-isolation-protocol.md" ]; then
    echo -e "${GREEN}✓ R212 - Phase Directory Isolation exists${NC}"
else
    echo -e "${YELLOW}⚠ R212 - Phase Directory Isolation not yet created${NC}"
fi

if [ -f "$RULE_DIR/R213-wave-and-effort-metadata-protocol.md" ]; then
    echo -e "${GREEN}✓ R213 - Wave and Effort Metadata exists${NC}"
else
    echo -e "${YELLOW}⚠ R213 - Wave and Effort Metadata not yet created${NC}"
fi

if [ -f "$RULE_DIR/R214-code-reviewer-wave-directory-acknowledgment.md" ]; then
    echo -e "${GREEN}✓ R214 - Wave Directory Acknowledgment exists${NC}"
else
    echo -e "${YELLOW}⚠ R214 - Wave Directory Acknowledgment not yet created${NC}"
fi

echo ""

# Run all validation scripts in order
echo -e "${MAGENTA}════════════════════════════════════════════════════════════════${NC}"
echo -e "${MAGENTA}🔍 Running Validation Scripts${NC}"
echo -e "${MAGENTA}════════════════════════════════════════════════════════════════${NC}"
echo ""

# Core validation scripts
run_validation "validate-effort-isolation.sh" "R209 - Effort Isolation Validation"
run_validation "validate-phase-isolation.sh" "R212 - Phase Isolation Validation"
run_validation "validate-r214-compliance.sh" "R214 - Wave Acknowledgment Validation"
run_validation "validate-all-metadata.sh" "Comprehensive Metadata Validation"

# Check template existence
echo -e "${BLUE}════════════════════════════════════════════════════════════════${NC}"
echo -e "${BLUE}📄 Template Check${NC}"
echo -e "${BLUE}════════════════════════════════════════════════════════════════${NC}"

TEMPLATES=(
    "PHASE-ARCHITECTURE-PLAN.md"
    "WAVE-ARCHITECTURE-PLAN.md"
    "PHASE-IMPLEMENTATION-PLAN.md"
    "WAVE-IMPLEMENTATION-PLAN.md"
    "EFFORT-PLANNING-TEMPLATE.md"
    "SPLIT-PLANNING-TEMPLATE.md"
    "WORK-LOG-TEMPLATE.md"
)

for template in "${TEMPLATES[@]}"; do
    if [ -f "../templates/$template" ]; then
        echo -e "${GREEN}✓ Template exists: $template${NC}"
    else
        echo -e "${YELLOW}⚠ Template missing: $template${NC}"
        ((WARNINGS++))
    fi
done

echo ""

# Check agent configurations
echo -e "${BLUE}════════════════════════════════════════════════════════════════${NC}"
echo -e "${BLUE}🤖 Agent Configuration Check${NC}"
echo -e "${BLUE}════════════════════════════════════════════════════════════════${NC}"

AGENTS=(
    "orchestrator.md"
    "sw-engineer.md"
    "architect-reviewer.md"
    "code-reviewer.md"
)

for agent in "${AGENTS[@]}"; do
    AGENT_FILE="../.claude/agents/$agent"
    if [ ! -f "$AGENT_FILE" ]; then
        AGENT_FILE="../agents/$agent"
    fi
    
    if [ -f "$AGENT_FILE" ]; then
        echo -e "${GREEN}✓ Agent config exists: $agent${NC}"
        
        # Check for critical rule references
        if grep -q "R209" "$AGENT_FILE" 2>/dev/null; then
            echo "  ✓ References R209 (Effort Isolation)"
        fi
        if grep -q "R212" "$AGENT_FILE" 2>/dev/null; then
            echo "  ✓ References R212 (Phase Isolation)"
        fi
        if grep -q "R213" "$AGENT_FILE" 2>/dev/null; then
            echo "  ✓ References R213 (Wave Metadata)"
        fi
        if grep -q "R214" "$AGENT_FILE" 2>/dev/null; then
            echo "  ✓ References R214 (Wave Acknowledgment)"
        fi
    else
        echo -e "${RED}✗ Agent config missing: $agent${NC}"
        ((FAILED_VALIDATIONS++))
    fi
done

echo ""

# Check for acknowledgment files
echo -e "${BLUE}════════════════════════════════════════════════════════════════${NC}"
echo -e "${BLUE}📋 Acknowledgment Files Check${NC}"
echo -e "${BLUE}════════════════════════════════════════════════════════════════${NC}"

ACK_FILES=(
    ".r209-acknowledged"
    ".r212-phase-acknowledged"
    ".r214-wave-acknowledged"
)

for ack_file in "${ACK_FILES[@]}"; do
    if [ -f "../$ack_file" ]; then
        echo -e "${GREEN}✓ Acknowledgment file exists: $ack_file${NC}"
        echo "  Recent entries:"
        tail -3 "../$ack_file" 2>/dev/null | sed 's/^/    /'
    else
        echo -e "${YELLOW}⚠ No acknowledgments yet: $ack_file${NC}"
    fi
done

echo ""

# Final summary
echo "════════════════════════════════════════════════════════════════"
echo -e "${BOLD}📊 FINAL VALIDATION SUMMARY${NC}"
echo "════════════════════════════════════════════════════════════════"
echo "Total Validations Run: $TOTAL_VALIDATIONS"
echo -e "${GREEN}Passed: $PASSED_VALIDATIONS${NC}"
echo -e "${RED}Failed: $FAILED_VALIDATIONS${NC}"
echo -e "${YELLOW}Warnings: $WARNINGS${NC}"
echo ""

# Calculate compliance score
if [ $TOTAL_VALIDATIONS -gt 0 ]; then
    COMPLIANCE_SCORE=$((PASSED_VALIDATIONS * 100 / TOTAL_VALIDATIONS))
else
    COMPLIANCE_SCORE=0
fi

echo -e "${BOLD}Compliance Score: ${COMPLIANCE_SCORE}%${NC}"
echo ""

# Overall assessment
if [ $FAILED_VALIDATIONS -eq 0 ] && [ $WARNINGS -eq 0 ]; then
    echo -e "${GREEN}${BOLD}🎉 PERFECT COMPLIANCE!${NC}"
    echo "The Software Factory 2.0 is fully compliant with all rules!"
    echo ""
    echo "Key achievements:"
    echo "✅ Orchestrator is master of all directory structures"
    echo "✅ Complete metadata hierarchy (Phase → Wave → Effort)"
    echo "✅ All agents have isolation protocols"
    echo "✅ Architecture-driven planning workflow established"
    echo "✅ Code reviewers acknowledge wave directories"
    exit 0
elif [ $FAILED_VALIDATIONS -eq 0 ]; then
    echo -e "${YELLOW}${BOLD}⚠️  MOSTLY COMPLIANT with $WARNINGS warnings${NC}"
    echo ""
    echo "Recommendations:"
    echo "• Create missing templates if any"
    echo "• Generate acknowledgment files during execution"
    echo "• Ensure all metadata explicitly marks ORCHESTRATOR as source"
    exit 0
else
    echo -e "${RED}${BOLD}❌ COMPLIANCE FAILURES DETECTED!${NC}"
    echo ""
    echo "Critical issues to address:"
    echo "• Fix failed validations immediately"
    echo "• Ensure orchestrator injects metadata at all levels"
    echo "• Verify agents acknowledge their isolation boundaries"
    echo "• Check that all rules are properly implemented"
    exit 1
fi