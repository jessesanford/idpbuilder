#!/bin/bash

# Fix PRIMARY DIRECTIVES placement in all SPAWN_CODE_REVIEWER states

STATES=(
    "SPAWN_CODE_REVIEWER_BACKPORT_PLAN"
    "SPAWN_CODE_REVIEWER_FIX_PLAN"
    "SPAWN_CODE_REVIEWER_INTEGRATE_WAVE_EFFORTS_FIX_PLAN"
    "CREATE_PHASE_FIX_PLAN"
    "SPAWN_CODE_REVIEWER_PHASE_IMPL"
    "SPAWN_CODE_REVIEWER_PHASE_MERGE_PLAN"
    "SPAWN_CODE_REVIEWER_PHASE_TEST_PLANNING"
    "SPAWN_CODE_REVIEWER_PROJECT_FIX_PLANNING"
    "SPAWN_CODE_REVIEWER_PROJECT_MERGE_PLAN"
    "SPAWN_CODE_REVIEWER_PROJECT_TEST_PLANNING"
    "SPAWN_CODE_REVIEWER_DEMO_VALIDATION"
    "SPAWN_CODE_REVIEWER_WAVE_IMPL"
    "SPAWN_CODE_REVIEWER_WAVE_TEST_PLANNING"
    "SPAWN_CODE_REVIEWERS_EFFORT_PLANNING"
    "SPAWN_CODE_REVIEWERS_EFFORT_REVIEW"
)

for STATE in "${STATES[@]}"; do
    FILE="agent-states/software-factory/orchestrator/${STATE}/rules.md"

    if [ ! -f "$FILE" ]; then
        echo "⚠️ File not found: $FILE"
        continue
    fi

    echo "Processing: $STATE"

    # Check if PRIMARY DIRECTIVES already exists in top 30 lines
    if head -30 "$FILE" | grep -q "^# PRIMARY DIRECTIVES"; then
        echo "  ✓ Already has PRIMARY DIRECTIVES at top"
        continue
    fi

    # Create temp file with new content
    TEMP_FILE=$(mktemp)

    # Get the title line
    TITLE=$(head -1 "$FILE")

    # Write new content
    echo "$TITLE" > "$TEMP_FILE"
    echo "" >> "$TEMP_FILE"
    echo "# PRIMARY DIRECTIVES" >> "$TEMP_FILE"
    echo "" >> "$TEMP_FILE"
    echo "You MUST read and acknowledge these rules:" >> "$TEMP_FILE"
    echo "" >> "$TEMP_FILE"
    echo "1. **R006** - Orchestrator cannot write code (BLOCKING)" >> "$TEMP_FILE"
    echo "   - File: \`\$CLAUDE_PROJECT_DIR/rule-library/R006-orchestrator-never-writes-code.md\`" >> "$TEMP_FILE"
    echo "" >> "$TEMP_FILE"

    # Add state-specific rules based on the state name
    if [[ "$STATE" == *"MERGE_PLAN"* ]]; then
        echo "2. **R269** - WAVE Integration Merge Plan Protocol (SUPREME LAW)" >> "$TEMP_FILE"
        echo "   - File: \`\$CLAUDE_PROJECT_DIR/rule-library/R269-WAVE-integration-merge-plan-protocol.md\`" >> "$TEMP_FILE"
        echo "" >> "$TEMP_FILE"
        echo "3. **R270** - PHASE Integration Merge Plan Protocol (SUPREME LAW)" >> "$TEMP_FILE"
        echo "   - File: \`\$CLAUDE_PROJECT_DIR/rule-library/R270-PHASE-integration-merge-plan-protocol.md\`" >> "$TEMP_FILE"
        echo "" >> "$TEMP_FILE"
    fi

    if [[ "$STATE" == *"FIX_PLAN"* ]] || [[ "$STATE" == *"BACKPORT"* ]]; then
        echo "2. **R256** - Fix Planning Protocol" >> "$TEMP_FILE"
        echo "   - File: \`\$CLAUDE_PROJECT_DIR/rule-library/R256-fix-planning-protocol.md\`" >> "$TEMP_FILE"
        echo "" >> "$TEMP_FILE"
    fi

    if [[ "$STATE" == *"TEST_PLANNING"* ]]; then
        echo "2. **R355** - Code Reviewer Test Planning" >> "$TEMP_FILE"
        echo "   - File: \`\$CLAUDE_PROJECT_DIR/rule-library/R355-code-reviewer-test-planning.md\`" >> "$TEMP_FILE"
        echo "" >> "$TEMP_FILE"
    fi

    if [[ "$STATE" == *"EFFORT_PLANNING"* ]] || [[ "$STATE" == *"FOR_REVIEW"* ]]; then
        echo "2. **R251** - Initial Effort Planning Protocol (SUPREME LAW)" >> "$TEMP_FILE"
        echo "   - File: \`\$CLAUDE_PROJECT_DIR/rule-library/R251-initial-effort-planning.md\`" >> "$TEMP_FILE"
        echo "" >> "$TEMP_FILE"
        echo "3. **R309** - Primary Implementation Effort Planning (SUPREME LAW)" >> "$TEMP_FILE"
        echo "   - File: \`\$CLAUDE_PROJECT_DIR/rule-library/R309-primary-implementation-effort-planning.md\`" >> "$TEMP_FILE"
        echo "" >> "$TEMP_FILE"
    fi

    # Add common rules for all states
    echo "4. **R287** - TODO Persistence Requirements (SUPREME LAW)" >> "$TEMP_FILE"
    echo "   - File: \`\$CLAUDE_PROJECT_DIR/rule-library/R287-todo-persistence-comprehensive.md\`" >> "$TEMP_FILE"
    echo "" >> "$TEMP_FILE"
    echo "5. **R288** - State File Update Requirements (SUPREME LAW)" >> "$TEMP_FILE"
    echo "   - File: \`\$CLAUDE_PROJECT_DIR/rule-library/R288-state-file-update-and-commit-protocol.md\`" >> "$TEMP_FILE"
    echo "" >> "$TEMP_FILE"
    echo "6. **R304** - Mandatory Line Counter Usage (SUPREME LAW)" >> "$TEMP_FILE"
    echo "   - File: \`\$CLAUDE_PROJECT_DIR/rule-library/R304-mandatory-line-counter-usage.md\`" >> "$TEMP_FILE"
    echo "" >> "$TEMP_FILE"
    echo "7. **R322** - Mandatory Stop After Spawn States (SUPREME LAW)" >> "$TEMP_FILE"
    echo "   - File: \`\$CLAUDE_PROJECT_DIR/rule-library/R322-mandatory-stop-after-spawn.md\`" >> "$TEMP_FILE"
    echo "" >> "$TEMP_FILE"
    echo "8. **R324** - State Transition Validation (SUPREME LAW)" >> "$TEMP_FILE"
    echo "   - File: \`\$CLAUDE_PROJECT_DIR/rule-library/R324-state-transition-validation.md\`" >> "$TEMP_FILE"
    echo "" >> "$TEMP_FILE"

    # Add the rest of the original file, skipping the old PRIMARY DIRECTIVES section if it exists
    IN_OLD_DIRECTIVES=0
    while IFS= read -r line; do
        # Skip the title line (already added)
        if [[ "$line" == "$TITLE" ]]; then
            continue
        fi

        # Check if we're entering the old PRIMARY DIRECTIVES section
        if [[ "$line" == *"PRIMARY DIRECTIVES"* ]]; then
            IN_OLD_DIRECTIVES=1
            continue
        fi

        # If we're in the old directives section, skip until we find the next major section
        if [ $IN_OLD_DIRECTIVES -eq 1 ]; then
            if [[ "$line" == "## "* ]] || [[ "$line" == "### "* ]] && [[ "$line" != "### "*"Rules"* ]]; then
                IN_OLD_DIRECTIVES=0
                # Don't skip this line, it's the start of a new section
            else
                continue
            fi
        fi

        # Write the line to temp file
        echo "$line" >> "$TEMP_FILE"
    done < <(tail -n +2 "$FILE")

    # Replace the original file
    mv "$TEMP_FILE" "$FILE"
    echo "  ✓ Fixed PRIMARY DIRECTIVES placement"
done

echo "Done! All SPAWN_CODE_REVIEWER states have been updated."