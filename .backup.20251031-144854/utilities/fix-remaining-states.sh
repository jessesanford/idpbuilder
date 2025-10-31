#!/bin/bash
# Fix the 12 remaining states that didn't get the checklist

TEMPLATE_FILE="/home/vscode/software-factory-template/templates/state-completion-checklist-template.md"
TEMPLATE_CONTENT=$(cat "$TEMPLATE_FILE")

# States that need the checklist
STATES=(
    "WAVE_COMPLETE"
    "WAITING_FOR_PROJECT_REVIEW_WAVE_INTEGRATION"
    "WAITING_FOR_PROJECT_VALIDATION"
    "WAITING_FOR_PROJECT_FIX_PLANS"
    "WAITING_FOR_PHASE_TEST_PLAN"
    "WAVE_START"
    "WAITING_FOR_WAVE_TEST_PLAN"
    "WAITING_FOR_PROJECT_TEST_PLAN"
    "REVIEW_WAVE_ARCHITECTURE"
    "WAITING_FOR_PHASE_PLANS"
    "WAITING_FOR_PHASE_MERGE_PLAN"
    "WAITING_FOR_PROJECT_MERGE_PLAN"
)

for STATE_NAME in "${STATES[@]}"; do
    STATE_FILE="agent-states/software-factory/orchestrator/$STATE_NAME/rules.md"

    echo "Processing: $STATE_NAME..."

    # Create customized checklist
    CHECKLIST=$(echo "$TEMPLATE_CONTENT" | sed "s/{{STATE_NAME}}/$STATE_NAME/g")

    # Find R405 line
    R405_LINE=$(grep -n "^## 🔴🔴🔴 R405 - MANDATORY AUTOMATION CONTINUATION FLAG 🔴🔴🔴" "$STATE_FILE" | cut -d: -f1)

    if [ -n "$R405_LINE" ]; then
        # Insert before R405
        {
            head -n $((R405_LINE - 1)) "$STATE_FILE"
            echo ""
            echo "$CHECKLIST"
            echo ""
            tail -n +$R405_LINE "$STATE_FILE"
        } > "${STATE_FILE}.tmp"

        mv "${STATE_FILE}.tmp" "$STATE_FILE"
        echo "  ✅ Updated"
    else
        echo "  ❌ No R405 section found"
    fi
done

echo ""
echo "✅ Complete - added checklist to 12 remaining states"
