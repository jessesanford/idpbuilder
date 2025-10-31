#!/bin/bash

# Script to update all SPAWN states with R313 requirement

SPAWN_STATES=(
    "SPAWN_CODE_REVIEWERS_EFFORT_PLANNING"
    "SPAWN_CODE_REVIEWERS_EFFORT_REVIEW"
    "SPAWN_SW_ENGINEERS"
    "SPAWN_INTEGRATION_AGENT"
    "SPAWN_ARCHITECT_PHASE_PLANNING"
    "SPAWN_ARCHITECT_PHASE_ASSESSMENT"
    "SPAWN_ARCHITECT_WAVE_PLANNING"
    "SPAWN_CODE_REVIEWER_FIX_PLAN"
    "CREATE_PHASE_FIX_PLAN"
    "SPAWN_CODE_REVIEWER_MERGE_PLAN"
    "SPAWN_CODE_REVIEWER_PHASE_IMPL"
    "SPAWN_CODE_REVIEWER_WAVE_IMPL"
)

for state in "${SPAWN_STATES[@]}"; do
    rules_file="/home/vscode/software-factory-template/agent-states/software-factory/orchestrator/$state/rules.md"
    
    if [ -f "$rules_file" ]; then
        echo "Updating $state..."
        
        # Check if R313 is already mentioned
        if grep -q "R313" "$rules_file"; then
            echo "  ✓ R313 already present in $state"
            continue
        fi
        
        # Add R313 as first rule in PRIMARY DIRECTIVES if not present
        # This is a simplified approach - manual review recommended
        echo "  → Adding R313 reference to $state"
        
        # Create a backup
        cp "$rules_file" "${rules_file}.backup.r313"
        
        # Add R313 warning at the top of the file if it contains spawn logic
        cat > /tmp/r313_header.txt << 'EOF'
## 🔴🔴🔴 R313 MANDATORY: STOP AFTER SPAWNING AGENTS 🔴🔴🔴

**CRITICAL REQUIREMENT PER R313:**
After spawning ANY agents in this state, you MUST:
1. Record what was spawned in state file
2. Save TODOs per R287
3. Commit and push state changes
4. Display stop message with continuation instructions
5. EXIT immediately with code 0

**VIOLATION = AUTOMATIC -100% FAILURE**

See: `$CLAUDE_PROJECT_DIR/rule-library/R313-mandatory-stop-after-spawn.md`

---

EOF
        
        # Prepend the R313 header to the file
        cat /tmp/r313_header.txt "$rules_file" > /tmp/updated_rules.md
        mv /tmp/updated_rules.md "$rules_file"
        
        echo "  ✓ Updated $state with R313 requirement"
    else
        echo "  ⚠ State $state rules file not found"
    fi
done

echo ""
echo "Summary:"
echo "========="
echo "Updated ${#SPAWN_STATES[@]} spawn states with R313 requirement"
echo "Backups created with .backup.r313 extension"
echo ""
echo "Please review the changes manually to ensure proper integration"