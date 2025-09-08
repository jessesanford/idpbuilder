#!/bin/bash

# CRITICAL: Enforce R322 - Mandatory Stop Before State Transitions
# This script updates ALL orchestrator state files with R322 enforcement

set -euo pipefail

echo "🚨🚨🚨 ENFORCING R322 - MANDATORY STOP BEFORE STATE TRANSITIONS 🚨🚨🚨"
echo "================================================================"

PROJECT_DIR="/home/vscode/software-factory-template"
STATES_DIR="$PROJECT_DIR/agent-states/orchestrator"

# Function to add R322 section to a rules file
add_r322_section() {
    local state_dir="$1"
    local state_name=$(basename "$state_dir")
    local rules_file="$state_dir/rules.md"
    
    if [ ! -f "$rules_file" ]; then
        echo "⚠️ Skipping $state_name - no rules.md file"
        return
    fi
    
    echo "📝 Processing $state_name..."
    
    # Check if R322 already exists
    if grep -q "R322" "$rules_file"; then
        echo "  ✓ R322 already mentioned - updating to ensure compliance"
    fi
    
    # Create temporary file with R322 enforcement at the TOP
    local temp_file="/tmp/rules_${state_name}.md"
    
    cat > "$temp_file" << 'EOF'
# Orchestrator - STATE_NAME State Rules

## 🛑🛑🛑 R322 MANDATORY STOP BEFORE STATE TRANSITIONS 🛑🛑🛑

**SUPREME LAW - VIOLATION = -100% IMMEDIATE FAILURE**

### YOU MUST STOP AFTER:
1. ✅ Completing all TODOs for this state
2. ✅ Updating orchestrator-state.json with new state
3. ✅ Committing and pushing the state file  
4. ✅ Providing work summary

### YOU MUST NOT:
- ❌ Continue to the next state automatically
- ❌ Start work for the new state
- ❌ Spawn agents for the new state
- ❌ Assume permission to continue

### STOP PROTOCOL:
```markdown
## 🛑 STATE TRANSITION CHECKPOINT: CURRENT_STATE → NEXT_STATE

### ✅ Current State Work Completed:
- [List completed work]

### 📊 Current Status:
- Current State: CURRENT_STATE
- Next State: NEXT_STATE
- TODOs Completed: X/Y
- State Files: Updated and committed ✅

### ⏸️ STOPPED - Awaiting User Continuation
Ready to transition to NEXT_STATE. Please use /continue-orchestrating.
```

**STOP MEANS STOP - Exit and wait for /continue-orchestrating**

---

EOF
    
    # Replace STATE_NAME placeholder
    sed -i "s/STATE_NAME/$state_name/g" "$temp_file"
    
    # Get the rest of the original file (skip the header if it exists)
    local skip_lines=0
    if head -n 1 "$rules_file" | grep -q "^# Orchestrator"; then
        skip_lines=1
    fi
    
    # Append the rest of the original content
    tail -n +$((skip_lines + 1)) "$rules_file" >> "$temp_file"
    
    # Now ensure PRIMARY DIRECTIVES section includes R322
    if grep -q "PRIMARY DIRECTIVES" "$temp_file"; then
        # Insert R322 right after PRIMARY DIRECTIVES header
        sed -i '/## PRIMARY DIRECTIVES/a\
\
### 🛑 RULE R322 - Mandatory Stop Before State Transitions\
**Source:** $CLAUDE_PROJECT_DIR/rule-library/R322-mandatory-stop-before-state-transitions.md\
**Criticality:** 🔴🔴🔴 SUPREME LAW - Violation = -100% FAILURE\
\
After completing state work and committing state file:\
1. STOP IMMEDIATELY\
2. Do NOT continue to next state\
3. Do NOT start new work\
4. Exit and wait for user\
---' "$temp_file"
    fi
    
    # Add STATE TRANSITION section if it doesn't exist
    if ! grep -q "STATE TRANSITION" "$temp_file"; then
        echo "" >> "$temp_file"
        echo "## STATE TRANSITION PROTOCOL (R322 ENFORCEMENT)" >> "$temp_file"
        echo "" >> "$temp_file"
        echo "When transitioning to next state:" >> "$temp_file"
        echo "1. Complete all state work" >> "$temp_file"
        echo "2. Update orchestrator-state.json" >> "$temp_file"
        echo "3. Commit and push state file" >> "$temp_file"
        echo "4. 🛑 **STOP HERE** - Per R322, exit and wait for /continue-orchestrating 🛑" >> "$temp_file"
        echo "" >> "$temp_file"
        echo "**DO NOT AUTOMATICALLY CONTINUE TO NEXT STATE!**" >> "$temp_file"
    fi
    
    # Add R322 VIOLATION DETECTION section
    if ! grep -q "R322 VIOLATION DETECTION" "$temp_file"; then
        echo "" >> "$temp_file"
        echo "## R322 VIOLATION DETECTION" >> "$temp_file"
        echo "" >> "$temp_file"
        echo "If you find yourself:" >> "$temp_file"
        echo "- Starting work for a new state without /continue-orchestrating" >> "$temp_file"
        echo "- Transitioning without stopping after state file commit" >> "$temp_file"
        echo "- Continuing after completing state work" >> "$temp_file"
        echo "" >> "$temp_file"
        echo "**STOP IMMEDIATELY - You are violating R322!**" >> "$temp_file"
    fi
    
    # Replace the original file
    mv "$temp_file" "$rules_file"
    echo "  ✅ Updated $state_name with R322 enforcement"
}

# Process all state directories
echo ""
echo "📂 Processing all orchestrator state directories..."
echo ""

for state_dir in "$STATES_DIR"/*; do
    if [ -d "$state_dir" ]; then
        add_r322_section "$state_dir"
    fi
done

echo ""
echo "✅ R322 enforcement added to all orchestrator state files"
echo ""
echo "📋 Summary of changes:"
echo "- Added R322 mandatory stop section at TOP of each file"
echo "- Added R322 to PRIMARY DIRECTIVES sections"
echo "- Added STATE TRANSITION protocol sections"
echo "- Added R322 VIOLATION DETECTION sections"
echo ""
echo "🔴🔴🔴 R322 IS NOW SUPREME LAW ACROSS ALL STATES 🔴🔴🔴"