#!/bin/bash

# Script to add mandatory rule reading enforcement to all orchestrator state rules.md files

BASE_DIR="/home/vscode/software-factory-template/agent-states/orchestrator"
TIMESTAMP=$(date '+%Y%m%d-%H%M%S')

# Define the enforcement text templates for each state
declare -A STATE_SPECIFIC_ACTIONS
STATE_SPECIFIC_ACTIONS["INIT"]="initialize orchestrator|read configuration|set up directories|create state files"
STATE_SPECIFIC_ACTIONS["PLANNING"]="load planning templates|spawn architects|request implementation plans|create phase plans"
STATE_SPECIFIC_ACTIONS["CREATE_NEXT_INFRASTRUCTURE"]="create effort directories|set up branches|initialize effort tracking|configure worktrees"
STATE_SPECIFIC_ACTIONS["ANALYZE_CODE_REVIEWER_PARALLELIZATION"]="analyze effort dependencies|determine parallelization strategy|plan reviewer spawning"
STATE_SPECIFIC_ACTIONS["SPAWN_CODE_REVIEWERS_EFFORT_PLANNING"]="spawn code reviewer agents|request effort plans|assign efforts to reviewers"
STATE_SPECIFIC_ACTIONS["WAITING_FOR_EFFORT_PLANS"]="check effort plan status|monitor reviewer progress|collect completed plans"
STATE_SPECIFIC_ACTIONS["ANALYZE_IMPLEMENTATION_PARALLELIZATION"]="analyze implementation dependencies|determine SWE parallelization|plan agent allocation"
STATE_SPECIFIC_ACTIONS["SPAWN_SW_ENGINEERS"]="spawn software engineer agents|assign effort work|distribute implementation tasks"
STATE_SPECIFIC_ACTIONS["MONITOR"]="check agent progress|monitor size limits|track implementation status|collect metrics"
STATE_SPECIFIC_ACTIONS["WAVE_COMPLETE"]="finalize wave efforts|collect implementation results|prepare integration"
STATE_SPECIFIC_ACTIONS["INTEGRATE_WAVE_EFFORTS"]="create integration branch|merge effort branches|resolve conflicts|run tests"
STATE_SPECIFIC_ACTIONS["REVIEW_WAVE_ARCHITECTURE"]="spawn architecture reviewer|request wave assessment|collect review feedback"
STATE_SPECIFIC_ACTIONS["INTEGRATE_PHASE_WAVES"]="merge wave branches|create phase branch|integrate wave work"
STATE_SPECIFIC_ACTIONS["COMPLETE_PHASE"]="finalize phase work|update documentation|prepare next phase"
STATE_SPECIFIC_ACTIONS["ERROR_RECOVERY"]="diagnose errors|recover from failures|restart failed efforts"
STATE_SPECIFIC_ACTIONS["PROJECT_DONE"]="finalize all work|generate reports|clean up resources"
STATE_SPECIFIC_ACTIONS["ERROR_RECOVERY"]="emergency shutdown|save critical state|preserve work in progress"
STATE_SPECIFIC_ACTIONS["INJECT_WAVE_METADATA"]="inject wave metadata|update tracking files|configure wave settings"
STATE_SPECIFIC_ACTIONS["SPAWN_ARCHITECT_PHASE_ASSESSMENT"]="spawn architect agent|request phase assessment|evaluate phase completion"
STATE_SPECIFIC_ACTIONS["SPAWN_ARCHITECT_PHASE_PLANNING"]="spawn architect agent|request phase planning|generate phase strategy"
STATE_SPECIFIC_ACTIONS["SPAWN_ARCHITECT_WAVE_PLANNING"]="spawn architect agent|request wave planning|generate wave strategy"
STATE_SPECIFIC_ACTIONS["SPAWN_CODE_REVIEWER_MERGE_PLAN"]="spawn code reviewer|request merge strategy|plan integration approach"
STATE_SPECIFIC_ACTIONS["SPAWN_CODE_REVIEWER_PHASE_IMPL"]="spawn code reviewer|review phase implementation|validate phase work"
STATE_SPECIFIC_ACTIONS["SPAWN_CODE_REVIEWER_WAVE_IMPL"]="spawn code reviewer|review wave implementation|validate wave work"
STATE_SPECIFIC_ACTIONS["SPAWN_INTEGRATION_AGENT"]="spawn integration specialist|coordinate merging|manage integration"
STATE_SPECIFIC_ACTIONS["WAITING_FOR_ARCHITECTURE_PLAN"]="monitor architect progress|check planning status|wait for architecture plans"
STATE_SPECIFIC_ACTIONS["WAITING_FOR_IMPLEMENTATION_PLAN"]="monitor planning progress|check plan status|wait for implementation plans"
STATE_SPECIFIC_ACTIONS["WAITING_FOR_PHASE_ASSESSMENT"]="monitor assessment progress|check architect status|wait for phase evaluation"
STATE_SPECIFIC_ACTIONS["WAVE_START"]="initialize wave|set up wave infrastructure|prepare wave execution"

# Function to generate enforcement text for a specific state
generate_enforcement_text() {
    local state_name="$1"
    local actions="${STATE_SPECIFIC_ACTIONS[$state_name]:-execute tasks|perform work|continue processing}"
    
    # Replace | with proper formatting
    local formatted_actions=""
    IFS='|' read -ra ACTION_ARRAY <<< "$actions"
    for action in "${ACTION_ARRAY[@]}"; do
        formatted_actions="${formatted_actions}- ❌ Start ${action}
"
    done
    
    cat <<EOF

## 🔴🔴🔴 STOP! STATE RULE READING IS ABSOLUTELY FIRST! 🔴🔴🔴

**YOU HAVE ENTERED ${state_name} STATE - YOU MUST READ AND ACKNOWLEDGE ALL STATE RULES BEFORE DOING ANY STATE WORK!**

### ❌ DO NOT DO ANY ${state_name} WORK UNTIL RULES ARE READ:
${formatted_actions}- ❌ Update state files
- ❌ Continue to next state
- ❌ Think about what to do in this state

### ✅ YOU MUST IMMEDIATELY:
1. **READ** every rule file listed in PRIMARY DIRECTIVES below
2. **ACKNOWLEDGE** each rule individually with number and description
3. **ONLY THEN** proceed with ${state_name} work

### 🚨 FAILURE TO READ STATE RULES FIRST = IMMEDIATE EXIT 🚨
**If you do ANY ${state_name} work before reading and acknowledging rules:**
- **STOP ALL WORK IMMEDIATELY**
- **EXIT WITH FAILURE STATUS**
- **YOU HAVE VIOLATED STATE COMPLIANCE**

**THE SYSTEM IS MONITORING YOUR READ TOOL CALLS!**

EOF
}

# Function to update a single rules.md file
update_rules_file() {
    local state_dir="$1"
    local state_name=$(basename "$state_dir")
    local rules_file="${state_dir}/rules.md"
    
    # Skip if rules.md doesn't exist
    if [[ ! -f "$rules_file" ]]; then
        echo "⚠️  No rules.md found in ${state_name}"
        return
    fi
    
    # Check if enforcement text already exists
    if grep -q "🔴🔴🔴 STOP! STATE RULE READING IS ABSOLUTELY FIRST!" "$rules_file" 2>/dev/null; then
        echo "✓ ${state_name} already has enforcement section"
        return
    fi
    
    echo "📝 Updating ${state_name}/rules.md"
    
    # Create backup
    cp "$rules_file" "${rules_file}.backup.${TIMESTAMP}"
    
    # Generate enforcement text for this state
    local enforcement_text=$(generate_enforcement_text "$state_name")
    
    # Read the original file
    local original_content=$(cat "$rules_file")
    
    # Find the position after the main heading (first line)
    local heading_line=$(echo "$original_content" | head -n1)
    local rest_content=$(echo "$original_content" | tail -n +2)
    
    # Combine: heading + enforcement + rest of content
    {
        echo "$heading_line"
        echo "$enforcement_text"
        echo "$rest_content"
    } > "$rules_file"
    
    echo "✅ Updated ${state_name}/rules.md"
}

# Main execution
echo "🏭 SOFTWARE FACTORY ORCHESTRATOR RULE ENFORCEMENT UPDATE"
echo "========================================================="
echo "Starting at: $(date '+%Y-%m-%d %H:%M:%S %Z')"
echo ""

# Process all orchestrator state directories
for state_dir in "$BASE_DIR"/*; do
    if [[ -d "$state_dir" ]]; then
        update_rules_file "$state_dir"
    fi
done

echo ""
echo "✅ Update complete!"
echo "Finished at: $(date '+%Y-%m-%d %H:%M:%S %Z')"

# Show summary
echo ""
echo "📊 Summary:"
echo "----------"
updated_count=$(find "$BASE_DIR" -name "rules.md" -newer "$BASE_DIR" -mmin -1 2>/dev/null | wc -l)
total_count=$(find "$BASE_DIR" -name "rules.md" 2>/dev/null | wc -l)
echo "Updated: ${updated_count} files"
echo "Total: ${total_count} files"
echo ""
echo "Backups created with suffix: .backup.${TIMESTAMP}"