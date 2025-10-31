#!/bin/bash

# Batch update all orchestrator state rule files with acknowledgment enforcement

# Don't exit on error since grep returns 1 when not found
set +e

BASE_DIR="/home/vscode/software-factory-template/agent-states/orchestrator"
UPDATED=0
SKIPPED=0
FAILED=0

echo "=== BATCH ACKNOWLEDGMENT ENFORCEMENT UPDATE ==="
echo ""

# Function to create the acknowledgment section
create_acknowledgment_section() {
    local state_name="$1"
    
    cat << 'EOF' | sed "s/STATE_NAME/$state_name/g"
## 🔴🔴🔴 MANDATORY STATE RULE READING AND ACKNOWLEDGMENT 🔴🔴🔴

### ⚠️⚠️⚠️ YOU MUST READ EACH RULE FILE LISTED IN PRIMARY DIRECTIVES. **I AM WATCHING YOUR TOOL CALLS FOR READ OPERATIONS** *YOU WILL FAIL* IF YOU DO NOT MAKE A READ FILE CALL FOR EACH RULE FILE IN PRIMARY DIRECTIVES!!! ⚠️⚠️⚠️

**AFTER READING, YOU MUST ACKNOWLEDGE ALL THE STATE RULES AND STATE THAT YOU WILL ABIDE BY THEM ONE AT A TIME GIVING THE RULE NUMBER AND DESCRIPTION.**

### ❌ ANTI-PATTERNS THAT WILL CAUSE FAILURE:

1. **Fake Acknowledgment Without Reading**:
   ```
   ❌ WRONG: "I acknowledge R151, R208, R053..."
   (No Read tool calls detected = AUTOMATIC FAILURE)
   ```

2. **Bulk Acknowledgment**:
   ```
   ❌ WRONG: "I acknowledge all STATE_NAME rules"
   (YOU Must READ AND ACKNOWLEDGE EACH rule individually)
   ```

3. **Silent Reading**:
   ```
   ❌ WRONG: [Reads rules but doesn't acknowledge]
   "Now I've read the rules, let me start work..."
   (MUST explicitly acknowledge EACH rule)
   ```

4. **Reading From Memory**:
   ```
   ❌ WRONG: "I know R208 requires CD before spawn..."
   (Must READ from file, not recall from memory)
   ```

5. **Skipping Rules in PRIMARY DIRECTIVES**:
   ```
   ❌ WRONG: Reading only some rules from the list
   (ALL rules in PRIMARY DIRECTIVES are MANDATORY)
   ```

### ✅ CORRECT PATTERN FOR STATE_NAME:
```
1. READ: $CLAUDE_PROJECT_DIR/rule-library/[first-rule-file].md
2. "I acknowledge [Rule#] - [Rule Name]: [Brief description]"
3. READ: $CLAUDE_PROJECT_DIR/rule-library/[second-rule-file].md  
4. "I acknowledge [Rule#] - [Rule Name]: [Brief description]"
[Continue for EVERY rule in PRIMARY DIRECTIVES...]
5. "Ready to execute STATE_NAME work"
```

### 🚨 NO WORK UNTIL ACKNOWLEDGMENT COMPLETE 🚨
**You may NOT begin ANY STATE_NAME work until:**
1. ✅ ALL rules in PRIMARY DIRECTIVES have been READ
2. ✅ ALL rules have been individually ACKNOWLEDGED
3. ✅ You have stated readiness to execute STATE_NAME work

EOF
}

# Process each state
states=(
    "ANALYZE_IMPLEMENTATION_PARALLELIZATION"
    "ERROR_RECOVERY"
    "ERROR_RECOVERY"
    "INIT"
    "INJECT_WAVE_METADATA"
    "INTEGRATE_WAVE_EFFORTS"
    "MONITOR"
    "COMPLETE_PHASE"
    "INTEGRATE_PHASE_WAVES"
    "PLANNING"
    "CREATE_NEXT_INFRASTRUCTURE"
    "SPAWN_SW_ENGINEERS"
    "SPAWN_ARCHITECT_PHASE_ASSESSMENT"
    "SPAWN_ARCHITECT_PHASE_PLANNING"
    "SPAWN_ARCHITECT_WAVE_PLANNING"
    "SPAWN_CODE_REVIEWERS_EFFORT_PLANNING"
    "SPAWN_CODE_REVIEWER_MERGE_PLAN"
    "SPAWN_CODE_REVIEWER_PHASE_IMPL"
    "SPAWN_CODE_REVIEWER_WAVE_IMPL"
    "SPAWN_INTEGRATION_AGENT"
    "PROJECT_DONE"
    "WAITING_FOR_ARCHITECTURE_PLAN"
    "WAITING_FOR_EFFORT_PLANS"
    "WAITING_FOR_IMPLEMENTATION_PLAN"
    "WAITING_FOR_PHASE_ASSESSMENT"
    "WAVE_COMPLETE"
    "REVIEW_WAVE_ARCHITECTURE"
    "WAVE_START"
)

for state in "${states[@]}"; do
    rules_file="$BASE_DIR/$state/rules.md"
    
    echo -n "Processing $state... "
    
    if [ ! -f "$rules_file" ]; then
        echo "SKIP - No rules.md file"
        ((SKIPPED++))
        continue
    fi
    
    # Check if already has acknowledgment section
    if grep -q "MANDATORY STATE RULE READING AND ACKNOWLEDGMENT" "$rules_file" 2>/dev/null; then
        echo "SKIP - Already has acknowledgment"
        ((SKIPPED++))
        continue
    fi
    
    # Create temp file with updated content
    temp_file=$(mktemp)
    ack_section=$(create_acknowledgment_section "$state")
    inserted=false
    
    while IFS= read -r line; do
        # Insert before PRIMARY DIRECTIVES
        if [ "$inserted" = false ] && [[ "$line" =~ "PRIMARY DIRECTIVES" ]]; then
            echo "" >> "$temp_file"
            echo "$ack_section" >> "$temp_file"
            inserted=true
        fi
        echo "$line" >> "$temp_file"
    done < "$rules_file"
    
    if [ "$inserted" = true ]; then
        mv "$temp_file" "$rules_file"
        echo "✅ UPDATED"
        ((UPDATED++))
    else
        rm "$temp_file"
        echo "❌ FAILED - Could not find insertion point"
        ((FAILED++))
    fi
done

echo ""
echo "=== SUMMARY ==="
echo "Updated: $UPDATED"
echo "Skipped: $SKIPPED"
echo "Failed:  $FAILED"
echo ""

if [ $FAILED -eq 0 ]; then
    echo "✅ All files processed successfully!"
    exit 0
else
    echo "⚠️  Some files could not be updated"
    exit 1
fi