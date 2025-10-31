#!/bin/bash

# Fix PRIMARY DIRECTIVES placement in remaining SPAWN states

STATES=(
    "SPAWN_SW_ENGINEERS"
    "SPAWN_ARCHITECT_MASTER_PLANNING"
    "SPAWN_ARCHITECT_PHASE_PLANNING"
    "SPAWN_ARCHITECT_WAVE_PLANNING"
    "SPAWN_SW_ENGINEERS"
    "SPAWN_INTEGRATION_AGENT"
    "SPAWN_INTEGRATION_AGENT_PHASE"
    "SPAWN_INTEGRATION_AGENT_PROJECT"
    "SPAWN_SW_ENGINEER_BACKPORT_FIXES"
    "SPAWN_SW_ENGINEER_PHASE_FIXES"
    "SPAWN_SW_ENGINEERS"
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
    if [[ "$STATE" == "SPAWN_SW_ENGINEERS" ]]; then
        echo "2. **R151** - Parallel Agent Timestamp Requirement (CRITICAL)" >> "$TEMP_FILE"
        echo "   - File: \`\$CLAUDE_PROJECT_DIR/rule-library/R151-parallel-agent-timestamp-requirement.md\`" >> "$TEMP_FILE"
        echo "" >> "$TEMP_FILE"
        echo "3. **R208** - Orchestrator Spawn Directory Protocol (SUPREME LAW)" >> "$TEMP_FILE"
        echo "   - File: \`\$CLAUDE_PROJECT_DIR/rule-library/R208-orchestrator-spawn-directory-protocol.md\`" >> "$TEMP_FILE"
        echo "" >> "$TEMP_FILE"
    fi

    if [[ "$STATE" == *"ARCHITECT"* ]]; then
        echo "2. **R308** - Architect Integration Role" >> "$TEMP_FILE"
        echo "   - File: \`\$CLAUDE_PROJECT_DIR/rule-library/R308-architect-integration-role.md\`" >> "$TEMP_FILE"
        echo "" >> "$TEMP_FILE"
    fi

    if [[ "$STATE" == *"INTEGRATE_WAVE_EFFORTS_AGENT"* ]]; then
        echo "2. **R269** - WAVE Integration Merge Plan Protocol (SUPREME LAW)" >> "$TEMP_FILE"
        echo "   - File: \`\$CLAUDE_PROJECT_DIR/rule-library/R269-WAVE-integration-merge-plan-protocol.md\`" >> "$TEMP_FILE"
        echo "" >> "$TEMP_FILE"
        echo "3. **R270** - PHASE Integration Merge Plan Protocol (SUPREME LAW)" >> "$TEMP_FILE"
        echo "   - File: \`\$CLAUDE_PROJECT_DIR/rule-library/R270-PHASE-integration-merge-plan-protocol.md\`" >> "$TEMP_FILE"
        echo "" >> "$TEMP_FILE"
    fi

    if [[ "$STATE" == *"SW_ENGINEER"* ]] || [[ "$STATE" == *"ENGINEERS_FOR_FIXES"* ]]; then
        echo "2. **R232** - Line Counting Requirements (SUPREME LAW)" >> "$TEMP_FILE"
        echo "   - File: \`\$CLAUDE_PROJECT_DIR/rule-library/R232-line-counting-requirements.md\`" >> "$TEMP_FILE"
        echo "" >> "$TEMP_FILE"
        echo "3. **R220** - Size Limit Compliance" >> "$TEMP_FILE"
        echo "   - File: \`\$CLAUDE_PROJECT_DIR/rule-library/R220-size-limit-compliance.md\`" >> "$TEMP_FILE"
        echo "" >> "$TEMP_FILE"
    fi

    if [[ "$STATE" == *"BACKPORT"* ]]; then
        echo "4. **R256** - Fix Planning Protocol" >> "$TEMP_FILE"
        echo "   - File: \`\$CLAUDE_PROJECT_DIR/rule-library/R256-fix-planning-protocol.md\`" >> "$TEMP_FILE"
        echo "" >> "$TEMP_FILE"
    fi

    # Add common rules for all states
    echo "5. **R287** - TODO Persistence Requirements (SUPREME LAW)" >> "$TEMP_FILE"
    echo "   - File: \`\$CLAUDE_PROJECT_DIR/rule-library/R287-todo-persistence-comprehensive.md\`" >> "$TEMP_FILE"
    echo "" >> "$TEMP_FILE"
    echo "6. **R288** - State File Update Requirements (SUPREME LAW)" >> "$TEMP_FILE"
    echo "   - File: \`\$CLAUDE_PROJECT_DIR/rule-library/R288-state-file-update-and-commit-protocol.md\`" >> "$TEMP_FILE"
    echo "" >> "$TEMP_FILE"
    echo "7. **R304** - Mandatory Line Counter Usage (SUPREME LAW)" >> "$TEMP_FILE"
    echo "   - File: \`\$CLAUDE_PROJECT_DIR/rule-library/R304-mandatory-line-counter-usage.md\`" >> "$TEMP_FILE"
    echo "" >> "$TEMP_FILE"
    echo "8. **R322** - Mandatory Stop After Spawn States (SUPREME LAW)" >> "$TEMP_FILE"
    echo "   - File: \`\$CLAUDE_PROJECT_DIR/rule-library/R322-mandatory-stop-after-spawn.md\`" >> "$TEMP_FILE"
    echo "" >> "$TEMP_FILE"
    echo "9. **R324** - State Transition Validation (SUPREME LAW)" >> "$TEMP_FILE"
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

echo "Done! All remaining SPAWN states have been updated."