#!/bin/bash

# utilities/validate-state-creation.sh
# Validates state rules.md files comply with R516 Part 3 requirements

set -euo pipefail

# Usage: validate-state-creation.sh <STATE_NAME> <AGENT_NAME>
# Example: validate-state-creation.sh SETUP_PHASE_INFRASTRUCTURE orchestrator

if [ $# -ne 2 ]; then
    echo "Usage: $0 <STATE_NAME> <AGENT_NAME>"
    echo "Example: $0 SETUP_PHASE_INFRASTRUCTURE orchestrator"
    exit 1
fi

STATE_NAME="$1"
AGENT_NAME="$2"

# Determine state file path based on agent
if [ "$AGENT_NAME" = "orchestrator" ]; then
    STATE_FILE="agent-states/software-factory/orchestrator/${STATE_NAME}/rules.md"
else
    STATE_FILE="agent-states/${AGENT_NAME}/${STATE_NAME}/rules.md"
fi

# Check if file exists
if [ ! -f "$STATE_FILE" ]; then
    echo "❌ FAIL: State file not found: $STATE_FILE"
    exit 1
fi

echo "🔍 Validating R516 compliance for: $STATE_FILE"
echo ""

# Track failures
FAILURES=0

# R516 Part 3: Required Sections
REQUIRED_SECTIONS=(
    "# ${STATE_NAME} State Rules"
    "## 🔴🔴🔴 MANDATORY STATE RULE READING AND ACKNOWLEDGMENT 🔴🔴🔴"
    "## 📋 PRIMARY DIRECTIVES FOR ${STATE_NAME} STATE"
    "## 🔴🔴🔴 MANDATORY EXECUTION CHECKLIST"
    "## State Purpose"
    "## Entry Criteria"
    "## State Actions"
    "## Exit Criteria"
    "## Rules Enforced"
    "## Transition Rules"
    "## 🔴🔴🔴 R405 - MANDATORY AUTOMATION CONTINUATION FLAG 🔴🔴🔴"
)

# Check each required section
for SECTION in "${REQUIRED_SECTIONS[@]}"; do
    if grep -qF "$SECTION" "$STATE_FILE"; then
        echo "✅ PASS: Section found: $SECTION"
    else
        echo "❌ FAIL: Missing required section: $SECTION"
        FAILURES=$((FAILURES + 1))
    fi
done

echo ""

# Check for core mandatory rules in PRIMARY DIRECTIVES section
echo "🔍 Checking Core Mandatory Rules..."

CORE_RULES=(
    "R510"
    "R288"
    "R287"
    "R405"
)

for RULE in "${CORE_RULES[@]}"; do
    if grep -q "${RULE}" "$STATE_FILE"; then
        echo "✅ PASS: Core rule referenced: ${RULE}"
    else
        echo "❌ FAIL: Missing core rule: ${RULE}"
        FAILURES=$((FAILURES + 1))
    fi
done

echo ""

# Check for BLOCKING REQUIREMENTS subsection in checklist
if grep -q "### BLOCKING REQUIREMENTS" "$STATE_FILE"; then
    echo "✅ PASS: BLOCKING REQUIREMENTS subsection found"
else
    echo "⚠️  WARNING: No BLOCKING REQUIREMENTS subsection (may be valid if no blocking items)"
fi

# Check for STANDARD EXECUTION TASKS subsection in checklist
if grep -q "### STANDARD EXECUTION TASKS" "$STATE_FILE"; then
    echo "✅ PASS: STANDARD EXECUTION TASKS subsection found"
else
    echo "⚠️  WARNING: No STANDARD EXECUTION TASKS subsection"
fi

# Check for EXIT REQUIREMENTS subsection in checklist
if grep -q "### EXIT REQUIREMENTS" "$STATE_FILE"; then
    echo "✅ PASS: EXIT REQUIREMENTS subsection found"
else
    echo "❌ FAIL: Missing EXIT REQUIREMENTS subsection"
    FAILURES=$((FAILURES + 1))
fi

echo ""

# Final result
if [ $FAILURES -eq 0 ]; then
    echo "════════════════════════════════════════════════════════"
    echo "✅ PROJECT_DONE: $STATE_FILE passes R516 validation"
    echo "════════════════════════════════════════════════════════"
    exit 0
else
    echo "════════════════════════════════════════════════════════"
    echo "❌ FAILURE: $STATE_FILE has $FAILURES R516 violations"
    echo "════════════════════════════════════════════════════════"
    exit 1
fi
