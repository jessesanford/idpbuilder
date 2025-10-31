#!/bin/bash

TEMPLATE_FILE="/home/vscode/software-factory-template/templates/state-completion-checklist-template.md"

for state in "WAITING_FOR_PHASE_MERGE_PLAN" "WAITING_FOR_PROJECT_MERGE_PLAN"; do
  STATE_FILE="agent-states/software-factory/orchestrator/$state/rules.md"
  echo "Processing: $state"

  # Create customized checklist
  CHECKLIST=$(cat "$TEMPLATE_FILE" | sed "s/{{STATE_NAME}}/$state/g")

  # Find R405 line
  R405_LINE=$(grep -n "^## 🔴🔴🔴 R405 - MANDATORY AUTOMATION CONTINUATION FLAG 🔴🔴🔴" "$STATE_FILE" | cut -d: -f1)

  if [ -n "$R405_LINE" ]; then
    {
      head -n $((R405_LINE - 1)) "$STATE_FILE"
      echo ""
      echo "$CHECKLIST"
      echo ""
      tail -n +$R405_LINE "$STATE_FILE"
    } > "${STATE_FILE}.tmp"
    mv "${STATE_FILE}.tmp" "$STATE_FILE"
    echo "  ✅ Updated"
  fi
done

echo "✅ Complete"
