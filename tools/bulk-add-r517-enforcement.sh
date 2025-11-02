#!/bin/bash
# Bulk update all state rules files with R517 State Manager enforcement section
# Created: 2025-11-01
# Purpose: Fulfill user mandate for complete State Manager enforcement

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
AGENT_STATES_DIR="$PROJECT_ROOT/agent-states"

echo "🏭 R517 ENFORCEMENT BULK UPDATE"
echo "================================"
echo "Project: $PROJECT_ROOT"
echo "Started: $(date -u '+%Y-%m-%d %H:%M:%S UTC')"
echo ""

# Counters
TOTAL_FILES=0
UPDATED_FILES=0
SKIPPED_FILES=0
FAILED_FILES=0
FAILED_LIST=()

# Create the enforcement section content
create_enforcement_section() {
    cat <<'ENFORCEMENT_EOF'

## 🔴🔴🔴 ABSOLUTE REQUIREMENT: STATE MANAGER CONSULTATION 🔴🔴🔴

**SUPREME LAW - R517 - UNIVERSAL STATE MANAGER CONSULTATION LAW**

**BEFORE EXITING THIS STATE, YOU MUST:**

1. **SPAWN STATE MANAGER FOR SHUTDOWN_CONSULTATION** (MANDATORY - NO EXCEPTIONS)
2. **NEVER UPDATE STATE FILES DIRECTLY** (orchestrator-state-v3.json, bug-tracking.json, etc.)
3. **NEVER COMMIT STATE FILES YOURSELF** (State Manager does this atomically)

**FAILURE TO CONSULT STATE MANAGER = IMMEDIATE SYSTEM HALT (-100% GRADE)**

### Enforcement Mechanism

If you attempt to exit this state without spawning State Manager:
- ❌ Pre-commit hooks will REJECT your commit
- ❌ Validation tools will FAIL the build
- ❌ Grading system will assign -100% penalty
- ❌ System will transition to ERROR_RECOVERY

### Required Pattern (COPY THIS EXACTLY)

```bash
# At end of state work, BEFORE any state file updates:

echo "🔴 MANDATORY: Spawning State Manager for SHUTDOWN_CONSULTATION"

# Spawn State Manager (REQUIRED - NOT OPTIONAL)
# Task: state-manager
# State: SHUTDOWN_CONSULTATION
# Current State: [YOUR_CURRENT_STATE]
# Proposed Next State: [YOUR_PROPOSED_NEXT_STATE]
# Work Summary: [SUMMARY_OF_WORK_COMPLETED]

# State Manager will:
# 1. Validate proposed transition against state machine
# 2. Update all 4 state files atomically
# 3. Commit with [R288] tag
# 4. Return REQUIRED next state (may differ from proposal)

# Wait for State Manager response
# Follow State Manager's directive (REQUIRED next state)
# DO NOT proceed until State Manager confirms
```

**YOU MUST NEVER:**
- ❌ Update orchestrator-state-v3.json yourself
- ❌ Update bug-tracking.json yourself
- ❌ Update integration-containers.json yourself
- ❌ Use `jq` to modify state files
- ❌ Use `sed/awk` to modify state files
- ❌ Set `validated_by: "orchestrator"` (must be "state-manager")
- ❌ Skip State Manager consultation "just this once"
- ❌ Think "I'll validate it manually"

**ONLY State Manager may update state files. This is NON-NEGOTIABLE.**

See: `rule-library/R517-universal-state-manager-consultation-law.md`

ENFORCEMENT_EOF
}

# Function to update a single file
update_file() {
    local file="$1"

    # Check if file already has R517 enforcement
    if grep -q "R517 - UNIVERSAL STATE MANAGER CONSULTATION LAW" "$file"; then
        echo "  ✅ SKIP: Already has R517 enforcement"
        ((SKIPPED_FILES++))
        return 0
    fi

    # Create backup
    cp "$file" "$file.backup-$(date +%Y%m%d-%H%M%S)"

    # Find insertion point (after PRIMARY DIRECTIVES section)
    if ! grep -q "# PRIMARY DIRECTIVES\|# CORE DIRECTIVES\|# STATE-SPECIFIC RULES" "$file"; then
        echo "  ⚠️  WARNING: No clear directive section found, inserting after first header"
        # Insert after first # header
        sed -i '0,/^# /{ /^# /a\
'"$(create_enforcement_section)"'
}' "$file"
    else
        # Insert after PRIMARY DIRECTIVES section
        # Find the line with PRIMARY DIRECTIVES and insert after the next blank line
        awk '
        /# PRIMARY DIRECTIVES|# CORE DIRECTIVES|# STATE-SPECIFIC RULES/ {
            found=1;
            print;
            next
        }
        found && /^$/ && !inserted {
            print;
            system("cat <<\"ENFORCEMENT_EOF\"\n\n## 🔴🔴🔴 ABSOLUTE REQUIREMENT: STATE MANAGER CONSULTATION 🔴🔴🔴\n\n**SUPREME LAW - R517 - UNIVERSAL STATE MANAGER CONSULTATION LAW**\n\n**BEFORE EXITING THIS STATE, YOU MUST:**\n\n1. **SPAWN STATE MANAGER FOR SHUTDOWN_CONSULTATION** (MANDATORY - NO EXCEPTIONS)\n2. **NEVER UPDATE STATE FILES DIRECTLY** (orchestrator-state-v3.json, bug-tracking.json, etc.)\n3. **NEVER COMMIT STATE FILES YOURSELF** (State Manager does this atomically)\n\n**FAILURE TO CONSULT STATE MANAGER = IMMEDIATE SYSTEM HALT (-100% GRADE)**\n\n### Enforcement Mechanism\n\nIf you attempt to exit this state without spawning State Manager:\n- ❌ Pre-commit hooks will REJECT your commit\n- ❌ Validation tools will FAIL the build\n- ❌ Grading system will assign -100% penalty\n- ❌ System will transition to ERROR_RECOVERY\n\n### Required Pattern (COPY THIS EXACTLY)\n\n\`\`\`bash\n# At end of state work, BEFORE any state file updates:\n\necho \"🔴 MANDATORY: Spawning State Manager for SHUTDOWN_CONSULTATION\"\n\n# Spawn State Manager (REQUIRED - NOT OPTIONAL)\n# Task: state-manager\n# State: SHUTDOWN_CONSULTATION\n# Current State: [YOUR_CURRENT_STATE]\n# Proposed Next State: [YOUR_PROPOSED_NEXT_STATE]\n# Work Summary: [SUMMARY_OF_WORK_COMPLETED]\n\n# State Manager will:\n# 1. Validate proposed transition against state machine\n# 2. Update all 4 state files atomically\n# 3. Commit with [R288] tag\n# 4. Return REQUIRED next state (may differ from proposal)\n\n# Wait for State Manager response\n# Follow State Manager'\''s directive (REQUIRED next state)\n# DO NOT proceed until State Manager confirms\n\`\`\`\n\n**YOU MUST NEVER:**\n- ❌ Update orchestrator-state-v3.json yourself\n- ❌ Update bug-tracking.json yourself\n- ❌ Update integration-containers.json yourself\n- ❌ Use \`jq\` to modify state files\n- ❌ Use \`sed/awk\` to modify state files\n- ❌ Set \`validated_by: \"orchestrator\"\` (must be \"state-manager\")\n- ❌ Skip State Manager consultation \"just this once\"\n- ❌ Think \"I'\''ll validate it manually\"\n\n**ONLY State Manager may update state files. This is NON-NEGOTIABLE.**\n\nSee: \`rule-library/R517-universal-state-manager-consultation-law.md\`\n\nENFORCEMENT_EOF");
            inserted=1;
            next
        }
        { print }
        ' "$file" > "$file.tmp" && mv "$file.tmp" "$file"
    fi

    # Verify the insertion worked
    if grep -q "R517 - UNIVERSAL STATE MANAGER CONSULTATION LAW" "$file"; then
        echo "  ✅ UPDATED successfully"
        ((UPDATED_FILES++))
        return 0
    else
        echo "  ❌ FAILED to insert enforcement section"
        ((FAILED_FILES++))
        FAILED_LIST+=("$file")
        # Restore backup
        mv "$file.backup-"* "$file" 2>/dev/null || true
        return 1
    fi
}

# Main processing loop
echo "📋 Processing state rules files..."
echo ""

while IFS= read -r file; do
    ((TOTAL_FILES++))

    # Get relative path for display
    rel_path="${file#$PROJECT_ROOT/}"

    echo "[$TOTAL_FILES] $rel_path"

    update_file "$file"

done < <(find "$AGENT_STATES_DIR" -name "rules.md" -type f | sort)

echo ""
echo "================================"
echo "📊 BULK UPDATE SUMMARY"
echo "================================"
echo "Total files processed: $TOTAL_FILES"
echo "✅ Files updated: $UPDATED_FILES"
echo "⏭️  Files skipped (already had R517): $SKIPPED_FILES"
echo "❌ Files failed: $FAILED_FILES"
echo ""

if [ $FAILED_FILES -gt 0 ]; then
    echo "⚠️  FAILED FILES:"
    for failed_file in "${FAILED_LIST[@]}"; do
        echo "  - ${failed_file#$PROJECT_ROOT/}"
    done
    echo ""
fi

# Calculate success rate
SUCCESS_RATE=$(( (UPDATED_FILES + SKIPPED_FILES) * 100 / TOTAL_FILES ))
echo "Success rate: $SUCCESS_RATE%"
echo ""

if [ $SUCCESS_RATE -eq 100 ]; then
    echo "🎉 PERFECT: 100% compliance achieved!"
    exit 0
else
    echo "⚠️  WARNING: Some files failed to update"
    exit 1
fi
