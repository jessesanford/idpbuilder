#!/bin/bash

# Validate State Rules Migration Script
# Ensures all states have the rules they need after migration

echo "═══════════════════════════════════════════════════════════════"
echo "     STATE RULES MIGRATION VALIDATION"
echo "═══════════════════════════════════════════════════════════════"
echo ""

ERRORS=0
WARNINGS=0

# Define rule mappings based on migration plan
declare -A RULE_MAPPINGS=(
    ["R234"]="SETUP_EFFORT_INFRASTRUCTURE ANALYZE_CODE_REVIEWER_PARALLELIZATION ANALYZE_IMPLEMENTATION_PARALLELIZATION"
    ["R208"]="SPAWN_AGENTS SPAWN_CODE_REVIEWERS_EFFORT_PLANNING SPAWN_CODE_REVIEWERS_FOR_SPLITS SPAWN_SW_ENGINEERS_FOR_FIXES SPAWN_ARCHITECT_FOR_WAVE_REVIEW SPAWN_ARCHITECT_FOR_PHASE_ASSESSMENT SPAWN_ARCHITECT_FOR_PROJECT_ASSESSMENT"
    ["R151"]="SPAWN_AGENTS SPAWN_CODE_REVIEWERS_EFFORT_PLANNING SPAWN_CODE_REVIEWERS_FOR_SPLITS SPAWN_SW_ENGINEERS_FOR_FIXES"
    ["R281"]="INIT"
    ["R221"]="SPAWN_AGENTS SETUP_EFFORT_INFRASTRUCTURE CREATE_NEXT_SPLIT_INFRASTRUCTURE"
    ["R235"]="SETUP_EFFORT_INFRASTRUCTURE SPAWN_AGENTS"
    ["R280"]="INTEGRATION PHASE_INTEGRATION PROJECT_INTEGRATION FINAL_INTEGRATION"
    ["R307"]="INTEGRATION PHASE_INTEGRATION PROJECT_INTEGRATION FINAL_INTEGRATION"
    ["R308"]="SETUP_EFFORT_INFRASTRUCTURE"
    ["R309"]="SETUP_EFFORT_INFRASTRUCTURE SPAWN_AGENTS SPAWN_CODE_REVIEWERS_EFFORT_PLANNING"
    ["R319"]="MONITOR_IMPLEMENTATION MONITOR_EFFORT_PLANNING MONITOR_SIZE_VALIDATION MONITOR_CODE_REVIEW MONITOR_ARCHITECT_REVIEW MONITOR_FIX_IMPLEMENTATION MONITOR_TESTING"
    ["R321"]="INTEGRATION PHASE_INTEGRATION PROJECT_INTEGRATION FINAL_INTEGRATION IMMEDIATE_BACKPORT_REQUIRED"
    ["R216"]="SPAWN_AGENTS SETUP_EFFORT_INFRASTRUCTURE"
)

# Rules that should be in bootstrap (not moved)
BOOTSTRAP_RULES=("R283" "R290" "R203" "R206" "R288" "R287" "R322" "R309" "R006")

echo "📋 Checking Bootstrap Rules Remain..."
echo "────────────────────────────────────"
for rule in "${BOOTSTRAP_RULES[@]}"; do
    if grep -q "$rule" /home/vscode/software-factory-template/ORCHESTRATOR-BOOTSTRAP-RULES.md 2>/dev/null; then
        echo "✅ $rule - Present in bootstrap"
    else
        echo "❌ $rule - MISSING from bootstrap!"
        ((ERRORS++))
    fi
done
echo ""

echo "📋 Validating State-Specific Rule Distribution..."
echo "────────────────────────────────────────────────"

for rule in "${!RULE_MAPPINGS[@]}"; do
    states=(${RULE_MAPPINGS[$rule]})
    echo ""
    echo "🔍 Checking $rule distribution:"
    
    for state in "${states[@]}"; do
        state_file="/home/vscode/software-factory-template/agent-states/orchestrator/$state/rules.md"
        
        if [ ! -f "$state_file" ]; then
            echo "  ⚠️  $state - State file doesn't exist!"
            ((WARNINGS++))
            continue
        fi
        
        if grep -q "$rule" "$state_file" 2>/dev/null; then
            echo "  ✅ $state - Has $rule"
        else
            echo "  ❌ $state - MISSING $rule!"
            ((ERRORS++))
        fi
    done
done

echo ""
echo "📋 Checking Critical States Have Essential Rules..."
echo "──────────────────────────────────────────────────"

# Check ERROR_RECOVERY has recovery rules
echo ""
echo "🔍 ERROR_RECOVERY State:"
ERROR_RECOVERY_RULES=("R019" "R156" "R010" "R258" "R257" "R259" "R300")
for rule in "${ERROR_RECOVERY_RULES[@]}"; do
    if grep -q "$rule" /home/vscode/software-factory-template/agent-states/orchestrator/ERROR_RECOVERY/rules.md 2>/dev/null; then
        echo "  ✅ Has $rule (recovery rule)"
    else
        echo "  ❌ MISSING $rule!"
        ((ERRORS++))
    fi
done

# Check INIT has initialization rules
echo ""
echo "🔍 INIT State:"
INIT_RULES=("R191" "R192" "R281" "R304")
for rule in "${INIT_RULES[@]}"; do
    if grep -q "$rule" /home/vscode/software-factory-template/agent-states/orchestrator/INIT/rules.md 2>/dev/null; then
        echo "  ✅ Has $rule (init rule)"
    else
        echo "  ❌ MISSING $rule!"
        ((ERRORS++))
    fi
done

# Check all states have R322 (stop before transitions)
echo ""
echo "📋 Verifying All States Have R322 (Stop Before Transitions)..."
echo "──────────────────────────────────────────────────────────────"
state_count=0
states_with_r322=0

for state_dir in /home/vscode/software-factory-template/agent-states/orchestrator/*/; do
    if [ -d "$state_dir" ]; then
        state_name=$(basename "$state_dir")
        state_file="$state_dir/rules.md"
        ((state_count++))
        
        if [ -f "$state_file" ] && grep -q "R322" "$state_file"; then
            ((states_with_r322++))
        else
            echo "  ❌ $state_name - Missing R322!"
            ((ERRORS++))
        fi
    fi
done

echo "  📊 States with R322: $states_with_r322/$state_count"

echo ""
echo "═══════════════════════════════════════════════════════════════"
echo "                    VALIDATION SUMMARY"
echo "═══════════════════════════════════════════════════════════════"
echo ""
echo "🔴 Critical Errors: $ERRORS"
echo "⚠️  Warnings: $WARNINGS"
echo ""

if [ $ERRORS -gt 0 ]; then
    echo "❌ MIGRATION NOT SAFE - Critical rules missing from states!"
    echo ""
    echo "🔧 Required Actions:"
    echo "  1. Add missing rules to the states identified above"
    echo "  2. Ensure R208, R221 are added to all SPAWN states"
    echo "  3. Ensure R321 is in all integration states"
    echo "  4. Re-run this validation script"
    exit 1
else
    echo "✅ Migration appears safe - all critical rules accounted for"
    echo ""
    echo "📋 Next Steps:"
    echo "  1. Update orchestrator.md to use minimal bootstrap"
    echo "  2. Test critical scenarios"
    echo "  3. Create migration report"
fi