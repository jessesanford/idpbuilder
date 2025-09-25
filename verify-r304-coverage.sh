#!/bin/bash

echo "================================================================"
echo "R304 LINE COUNTER ENFORCEMENT - COMPREHENSIVE COVERAGE REPORT"
echo "================================================================"
echo ""

BASE_DIR="/home/vscode/software-factory-template/agent-states"
TOTAL_STATES=0
STATES_WITH_R304=0
STATES_WITHOUT_R304=0

# Track which states are missing R304
MISSING_STATES=()

for agent_dir in "$BASE_DIR"/*/; do
    agent_name=$(basename "$agent_dir")
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo "AGENT: $agent_name"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    
    agent_total=0
    agent_with=0
    agent_without=0
    
    for state_dir in "$agent_dir"*/; do
        if [ -d "$state_dir" ]; then
            state_name=$(basename "$state_dir")
            rules_file="$state_dir/rules.md"
            
            if [ -f "$rules_file" ]; then
                TOTAL_STATES=$((TOTAL_STATES + 1))
                agent_total=$((agent_total + 1))
                
                # Check for R304 reference or line-counter.sh mention
                if grep -q "R304\|line-counter.sh" "$rules_file" 2>/dev/null; then
                    echo "  ✅ $state_name - HAS R304 enforcement"
                    STATES_WITH_R304=$((STATES_WITH_R304 + 1))
                    agent_with=$((agent_with + 1))
                else
                    echo "  ❌ $state_name - MISSING R304 enforcement"
                    STATES_WITHOUT_R304=$((STATES_WITHOUT_R304 + 1))
                    agent_without=$((agent_without + 1))
                    MISSING_STATES+=("$agent_name/$state_name")
                fi
            fi
        fi
    done
    
    echo ""
    echo "  Agent Summary: $agent_with/$agent_total states have R304"
    if [ $agent_without -gt 0 ]; then
        echo "  ⚠️ WARNING: $agent_without states missing R304!"
    fi
    echo ""
done

echo "================================================================"
echo "OVERALL SUMMARY"
echo "================================================================"
echo "Total States: $TOTAL_STATES"
echo "States with R304: $STATES_WITH_R304"
echo "States without R304: $STATES_WITHOUT_R304"
echo ""

if [ $STATES_WITHOUT_R304 -eq 0 ]; then
    echo "🎉 SUCCESS: ALL STATES HAVE R304 ENFORCEMENT! 🎉"
else
    echo "⚠️ WARNING: $STATES_WITHOUT_R304 states are still missing R304!"
    echo ""
    echo "States missing R304:"
    for missing in "${MISSING_STATES[@]}"; do
        echo "  - $missing"
    done
fi

echo ""
echo "================================================================"
echo "CRITICAL STATES CHECK (states that MUST have R304)"
echo "================================================================"

CRITICAL_STATES=(
    "sw-engineer/MEASURE_SIZE"
    "sw-engineer/IMPLEMENTATION"
    "sw-engineer/SPLIT_IMPLEMENTATION"
    "code-reviewer/CODE_REVIEW"
    "code-reviewer/CREATE_SPLIT_PLAN"
    "code-reviewer/SPLIT_REVIEW"
    "code-reviewer/VALIDATION"
    "orchestrator/MONITOR"
    "orchestrator/PHASE_INTEGRATION"
    "orchestrator/WAVE_COMPLETE"
)

all_critical_ok=true
for critical in "${CRITICAL_STATES[@]}"; do
    file="$BASE_DIR/$critical/rules.md"
    if [ -f "$file" ]; then
        if grep -q "R304\|line-counter.sh" "$file" 2>/dev/null; then
            echo "  ✅ $critical - OK"
        else
            echo "  🔴 $critical - CRITICAL FAILURE: Missing R304!"
            all_critical_ok=false
        fi
    else
        echo "  ⚠️ $critical - File not found"
    fi
done

echo ""
if [ "$all_critical_ok" = true ]; then
    echo "✅ All critical states have R304 enforcement"
else
    echo "🔴 CRITICAL: Some measurement-critical states are missing R304!"
fi

echo ""
echo "================================================================"