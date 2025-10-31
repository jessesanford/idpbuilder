#!/bin/bash

# Script to add validation requirements to all state update examples in orchestrator rules

set -euo pipefail

echo "🔍 Adding validation requirements to orchestrator state rules..."

# Function to add validation note to a file
add_validation_note() {
    local file="$1"
    
    # Check if file already mentions validation
    if grep -q "validate-state.sh\|MANDATORY VALIDATION" "$file"; then
        echo "  ✓ Already has validation: $(basename $(dirname "$file"))"
        return
    fi
    
    # Check if file mentions updating orchestrator-state-v3.json
    if ! grep -q "orchestrator-state-v3.json" "$file"; then
        return
    fi
    
    echo "  📝 Updating: $(basename $(dirname "$file"))"
    
    # Add validation note after state update sections
    # This is a simplified approach - for production, would need more careful editing
    
    # Create a temporary marker file with the validation note
    cat > /tmp/validation_note.txt << 'EOF'

### 🔴🔴🔴 MANDATORY VALIDATION REQUIREMENT 🔴🔴🔴

**Per R288 and R324**: ALL state file updates MUST be validated before commit:

```bash
# After ANY update to orchestrator-state-v3.json:
"$CLAUDE_PROJECT_DIR/tools/validate-state.sh" orchestrator-state-v3.json || {
    echo "❌ State file validation failed!"
    exit 288
}
```

**Use helper functions for automatic validation:**
```bash
# Source the helper functions
source "$CLAUDE_PROJECT_DIR/utilities/state-file-update-functions.sh"

# Use safe functions that include validation:
safe_state_transition "NEW_STATE" "reason"
safe_update_field "field_name" "value"
```
EOF
    
    # Check if the file has a section about state updates
    if grep -q "git add orchestrator-state-v3.json\|git commit.*orchestrator-state-v3.json\|jq.*orchestrator-state-v3.json" "$file"; then
        # Append validation note at the end of the file
        echo "" >> "$file"
        cat /tmp/validation_note.txt >> "$file"
    fi
    
    rm -f /tmp/validation_note.txt
}

# Process all orchestrator state rule files
for rules_file in /home/vscode/software-factory-template/agent-states/software-factory/orchestrator/*/rules.md; do
    add_validation_note "$rules_file"
done

echo ""
echo "✅ Validation requirements added to orchestrator state rules"
echo ""
echo "📋 Summary of changes:"
echo "  - Added validation requirement notes to state rule files"
echo "  - Referenced validate-state.sh tool"
echo "  - Referenced state-file-update-functions.sh helpers"
echo ""
echo "⚠️  Note: This was a bulk update. Review individual files for accuracy."