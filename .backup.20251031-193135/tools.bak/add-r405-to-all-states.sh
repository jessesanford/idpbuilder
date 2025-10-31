#!/bin/bash
# Script to add R405 automation continuation flag requirement to all state rule files

set -euo pipefail

CLAUDE_PROJECT_DIR="/home/vscode/software-factory-template"
cd "$CLAUDE_PROJECT_DIR"

# R405 content to add to each state rule file
R405_CONTENT='
## 🔴🔴🔴 R405 - MANDATORY AUTOMATION CONTINUATION FLAG 🔴🔴🔴

**SUPREME LAW - PENALTY FOR VIOLATION: -100% GRADE**

### YOU MUST OUTPUT THE CONTINUATION FLAG AS YOUR LAST ACTION

**EVERY STATE COMPLETION MUST END WITH EXACTLY ONE OF:**
```bash
# If state completed successfully and factory should continue:
echo "CONTINUE-SOFTWARE-FACTORY=TRUE"

# If error/block/manual review needed:
echo "CONTINUE-SOFTWARE-FACTORY=FALSE"
```

### CRITICAL REQUIREMENTS:
1. **ABSOLUTE LAST OUTPUT**: This MUST be the very last line of output before state completion
2. **NO TEXT AFTER**: No explanations, summaries, or any text after the flag
3. **EXACTLY AS SHOWN**: Use exact format - no variations like CONTINUE-ORCHESTRATING
4. **ALWAYS REQUIRED**: Every single state must output this flag
5. **GREPPABLE**: Must be on its own line for automation parsing

### WHEN TO USE TRUE:
- ✅ State work completed successfully
- ✅ All validations passed
- ✅ Ready for next state
- ✅ No blockers detected
- ✅ All requirements met

### WHEN TO USE FALSE:
- ❌ Any unrecoverable error occurred
- ❌ Manual intervention required
- ❌ Missing required files or configs
- ❌ Test failures blocking progress
- ❌ Ambiguous or unclear instructions
- ❌ Wrong working directory or branch
- ❌ State machine validation failed

### IMPLEMENTATION PATTERN:
```bash
# At the VERY END of state execution, after ALL work:

# Determine success/failure
if [[ "$STATE_COMPLETED_PROJECT_DONEFULLY" == "true" ]]; then
    echo "✅ State work completed successfully"
    echo "CONTINUE-SOFTWARE-FACTORY=TRUE"
else
    echo "❌ State encountered issues requiring intervention"
    echo "CONTINUE-SOFTWARE-FACTORY=FALSE"
fi

# DO NOT OUTPUT ANYTHING AFTER THE FLAG!
```

### GRADING IMPACT:
- **Missing flag**: -100% AUTOMATIC FAILURE
- **Text after flag**: -50% penalty
- **Wrong format**: -100% AUTOMATIC FAILURE
- **Multiple flags**: -50% penalty

**See: $CLAUDE_PROJECT_DIR/rule-library/R405-automation-continuation-flag.md**
'

# Count files to process
TOTAL_FILES=$(find agent-states -name "rules.md" -type f | wc -l)
echo "🔍 Found $TOTAL_FILES state rule files to update"
echo ""

# Process each file
UPDATED=0
SKIPPED=0
ERRORS=0

find agent-states -name "rules.md" -type f | while read -r file; do
    # Extract agent and state from path
    # Example: agent-states/main/orchestrator/INIT/rules.md
    rel_path="${file#agent-states/}"
    context=$(echo "$rel_path" | cut -d'/' -f1)
    agent=$(echo "$rel_path" | cut -d'/' -f2)
    state=$(echo "$rel_path" | cut -d'/' -f3)

    echo -n "📝 Processing: $context/$agent/$state... "

    # Check if R405 already exists
    if grep -q "R405\|CONTINUE-SOFTWARE-FACTORY" "$file"; then
        echo "✅ Already has R405 (skipping)"
        ((SKIPPED++)) || true
        continue
    fi

    # Create backup
    cp "$file" "${file}.backup-pre-r405" 2>/dev/null || {
        echo "❌ Failed to backup"
        ((ERRORS++)) || true
        continue
    }

    # Add R405 content before the last section or at the end
    # Try to add it before "## Related Rules" or "## Common Violations" if they exist
    if grep -q "^## Related Rules" "$file"; then
        # Add before Related Rules section
        awk -v content="$R405_CONTENT" '
            /^## Related Rules/ { print content; print ""; }
            { print }
        ' "$file" > "${file}.tmp" && mv "${file}.tmp" "$file"
        echo "✅ Added before Related Rules"
    elif grep -q "^## Common Violations" "$file"; then
        # Add before Common Violations section
        awk -v content="$R405_CONTENT" '
            /^## Common Violations/ { print content; print ""; }
            { print }
        ' "$file" > "${file}.tmp" && mv "${file}.tmp" "$file"
        echo "✅ Added before Common Violations"
    else
        # Add at the end of file
        echo "$R405_CONTENT" >> "$file"
        echo "✅ Added at end of file"
    fi

    ((UPDATED++)) || true
done

echo ""
echo "📊 SUMMARY:"
echo "  ✅ Updated: $UPDATED files"
echo "  ⏭️  Skipped: $SKIPPED files (already had R405)"
echo "  ❌ Errors: $ERRORS files"
echo ""
echo "✅ R405 automation flag requirement added to all state rule files!"
echo ""
echo "🔍 Verifying changes..."

# Verify the changes
VERIFIED=$(find agent-states -name "rules.md" -type f -exec grep -l "R405\|CONTINUE-SOFTWARE-FACTORY" {} \; | wc -l)
echo "  Files with R405: $VERIFIED/$TOTAL_FILES"

if [ "$VERIFIED" -eq "$TOTAL_FILES" ]; then
    echo "  ✅ ALL state files now have R405!"
else
    MISSING=$((TOTAL_FILES - VERIFIED))
    echo "  ⚠️ Warning: $MISSING files still missing R405"
    echo "  Run verification: find agent-states -name 'rules.md' -type f -exec grep -L 'R405' {} \;"
fi

echo ""
echo "📝 Next steps:"
echo "  1. Review a sample of updated files to verify correctness"
echo "  2. Run: git diff agent-states/ | head -100"
echo "  3. Commit changes: git add agent-states/ && git commit -m 'fix: Add R405 automation flag to all state rule files'"
echo "  4. Push changes: git push"