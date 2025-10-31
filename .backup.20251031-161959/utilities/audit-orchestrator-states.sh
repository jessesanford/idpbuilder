#!/bin/bash
# ORCHESTRATOR STATE AUDIT SCRIPT
# Comprehensive analysis of R322 and R405 compliance across all states

set -e

TOTAL_FILES=0
R322_REFERENCE=0
R405_SECTION=0
AGENT_FACTORY_CLARIFICATION=0
PROPER_DEMARCATION=0
ISSUES_FOUND=()

echo "=========================================="
echo "ORCHESTRATOR STATE AUDIT"
echo "Date: $(date '+%Y-%m-%d %H:%M:%S')"
echo "=========================================="
echo ""

# Find all orchestrator state files
STATE_FILES=$(find /home/vscode/software-factory-template/agent-states -type f -name "rules.md" | grep -E "(orchestrator|software-factory/orchestrator)" | sort)

for STATE_FILE in $STATE_FILES; do
    ((TOTAL_FILES++))

    STATE_NAME=$(basename $(dirname "$STATE_FILE"))
    STATE_PATH=$(echo "$STATE_FILE" | sed 's|/home/vscode/software-factory-template/||')

    echo "[$TOTAL_FILES] $STATE_PATH"
    echo "    State: $STATE_NAME"

    # Check R322 reference
    if grep -q "R322" "$STATE_FILE"; then
        ((R322_REFERENCE++))
        echo "    ✅ R322 reference found"
    else
        echo "    ❌ R322 reference MISSING"
        ISSUES_FOUND+=("$STATE_PATH: Missing R322 reference")
    fi

    # Check R405 section
    if grep -q "R405.*AUTOMATION CONTINUATION FLAG" "$STATE_FILE"; then
        ((R405_SECTION++))
        echo "    ✅ R405 section found"
    else
        echo "    ❌ R405 section MISSING"
        ISSUES_FOUND+=("$STATE_PATH: Missing R405 section")
    fi

    # Check for AGENT STOPS ≠ FACTORY STOPS clarification
    if grep -q "AGENT STOP.*FACTORY STOP" "$STATE_FILE" || \
       grep -q "stop inference.*flag.*TRUE" "$STATE_FILE" || \
       grep -q "Stop inference.*CONTINUE.*TRUE" "$STATE_FILE"; then
        ((AGENT_FACTORY_CLARIFICATION++))
        echo "    ✅ Agent≠Factory clarification found"
    else
        echo "    ⚠️  Agent≠Factory clarification MISSING"
        ISSUES_FOUND+=("$STATE_PATH: Missing Agent≠Factory clarification")
    fi

    # Check for proper rule demarcation
    if grep -q "### 🛑.*R322" "$STATE_FILE" || \
       grep -q "## 🔴🔴🔴.*R322" "$STATE_FILE" || \
       grep -q "<!-- RULE:R322" "$STATE_FILE"; then
        ((PROPER_DEMARCATION++))
        echo "    ✅ Proper R322 demarcation"
    else
        echo "    ⚠️  Weak R322 demarcation"
    fi

    echo ""
done

echo "=========================================="
echo "AUDIT SUMMARY"
echo "=========================================="
echo ""
echo "Total orchestrator state files: $TOTAL_FILES"
echo ""
echo "R322 Coverage:"
echo "  Files with R322 reference: $R322_REFERENCE/$TOTAL_FILES"
echo "  Percentage: $(( R322_REFERENCE * 100 / TOTAL_FILES ))%"
echo ""
echo "R405 Coverage:"
echo "  Files with R405 section: $R405_SECTION/$TOTAL_FILES"
echo "  Percentage: $(( R405_SECTION * 100 / TOTAL_FILES ))%"
echo ""
echo "Clarification Coverage:"
echo "  Files with Agent≠Factory clarification: $AGENT_FACTORY_CLARIFICATION/$TOTAL_FILES"
echo "  Percentage: $(( AGENT_FACTORY_CLARIFICATION * 100 / TOTAL_FILES ))%"
echo ""
echo "Demarcation Quality:"
echo "  Files with proper R322 demarcation: $PROPER_DEMARCATION/$TOTAL_FILES"
echo "  Percentage: $(( PROPER_DEMARCATION * 100 / TOTAL_FILES ))%"
echo ""

if [ ${#ISSUES_FOUND[@]} -gt 0 ]; then
    echo "=========================================="
    echo "ISSUES DETECTED (${#ISSUES_FOUND[@]} total)"
    echo "=========================================="
    for ISSUE in "${ISSUES_FOUND[@]}"; do
        echo "  - $ISSUE"
    done
    echo ""
fi

echo "=========================================="
echo "COMPLIANCE SCORE"
echo "=========================================="
TOTAL_CHECKS=$((TOTAL_FILES * 3))  # R322, R405, Clarification
PASSED_CHECKS=$((R322_REFERENCE + R405_SECTION + AGENT_FACTORY_CLARIFICATION))
SCORE=$(( PASSED_CHECKS * 100 / TOTAL_CHECKS ))
echo "Overall: $SCORE% ($PASSED_CHECKS/$TOTAL_CHECKS checks passed)"
echo ""

if [ $SCORE -lt 100 ]; then
    echo "❌ SYNCHRONIZATION REQUIRED - Score below 100%"
    exit 1
else
    echo "✅ FULL SYNCHRONIZATION ACHIEVED"
    exit 0
fi
