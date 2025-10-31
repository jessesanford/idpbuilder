#!/bin/bash
# Add STATE COMPLETION CHECKLIST to all SF 3.0 orchestrator states
# This ensures every state has clear exit protocol including R405 continuation flag

# Don't exit on error - we want to process all files
set +e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"
TEMPLATE_FILE="$PROJECT_DIR/templates/state-completion-checklist-template.md"
ORCHESTRATOR_DIR="$PROJECT_DIR/agent-states/software-factory/orchestrator"

echo "════════════════════════════════════════════════════════════"
echo "Adding STATE COMPLETION CHECKLIST to SF 3.0 Orchestrator States"
echo "════════════════════════════════════════════════════════════"
echo ""

# Check template exists
if [ ! -f "$TEMPLATE_FILE" ]; then
    echo "❌ Error: Template file not found: $TEMPLATE_FILE"
    exit 1
fi

# Read template content
TEMPLATE_CONTENT=$(cat "$TEMPLATE_FILE")

# Counter
UPDATED=0
SKIPPED=0
ERRORS=0

# Find all state rules.md files (excluding DEPRECATED and _TEMPLATES)
while IFS= read -r state_file; do
    # Extract state name from path
    STATE_NAME=$(basename "$(dirname "$state_file")")

    echo "Processing: $STATE_NAME..."

    # Check if checklist already exists
    if grep -q "STATE COMPLETION CHECKLIST" "$state_file"; then
        echo "  ⏭️  Skipped - checklist already exists"
        ((SKIPPED++))
        continue
    fi

    # Check if file has R405 section (where we'll insert before)
    if ! grep -q "🔴🔴🔴 R405 - MANDATORY AUTOMATION CONTINUATION FLAG 🔴🔴🔴" "$state_file"; then
        echo "  ⚠️  Warning - No R405 section found, appending to end"
    fi

    # Create customized checklist for this state
    CHECKLIST=$(echo "$TEMPLATE_CONTENT" | sed "s/{{STATE_NAME}}/$STATE_NAME/g")

    # Create backup
    cp "$state_file" "${state_file}.backup"

    # Find the line number of the R405 section
    R405_LINE=$(grep -n "^## 🔴🔴🔴 R405 - MANDATORY AUTOMATION CONTINUATION FLAG 🔴🔴🔴" "$state_file" | cut -d: -f1 || echo "")

    if [ -n "$R405_LINE" ]; then
        # Insert checklist before R405 section
        {
            head -n $((R405_LINE - 1)) "$state_file"
            echo ""
            echo "$CHECKLIST"
            echo ""
            tail -n +$R405_LINE "$state_file"
        } > "${state_file}.tmp"

        mv "${state_file}.tmp" "$state_file"
        echo "  ✅ Updated - checklist inserted before R405 section"
        ((UPDATED++))
    else
        # Append to end if no R405 section found
        {
            cat "$state_file"
            echo ""
            echo "$CHECKLIST"
            echo ""
        } > "${state_file}.tmp"

        mv "${state_file}.tmp" "$state_file"
        echo "  ✅ Updated - checklist appended to end"
        ((UPDATED++))
    fi

    # Remove backup if successful
    rm -f "${state_file}.backup"

done < <(find "$ORCHESTRATOR_DIR" -name "rules.md" ! -path "*/DEPRECATED/*" ! -path "*/_TEMPLATES/*" | sort)

echo ""
echo "════════════════════════════════════════════════════════════"
echo "Update Complete"
echo "════════════════════════════════════════════════════════════"
echo "✅ Updated: $UPDATED files"
echo "⏭️  Skipped: $SKIPPED files (already had checklist)"
echo "❌ Errors: $ERRORS files"
echo ""
echo "Next step: Review changes and commit"
echo "  git diff agent-states/software-factory/orchestrator/"
echo "  git add agent-states/software-factory/orchestrator/"
echo "  git commit -m 'feat: Add STATE COMPLETION CHECKLIST to all SF 3.0 orchestrator states [R405]'"
echo ""
