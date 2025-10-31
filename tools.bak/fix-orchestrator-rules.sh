#!/bin/bash

# Script to fix missing rule references in orchestrator state files
# Strategy: Remove references to non-existent rules

echo "🔧 Starting comprehensive fix for missing rule references..."
echo ""

FIXES_MADE=0

# Fix INTEGRATE_WAVE_EFFORTS state - Remove R020 reference (non-existent)
echo "📝 Fixing INTEGRATE_WAVE_EFFORTS state..."
if grep -q "R020" /home/vscode/software-factory-template/agent-states/software-factory/orchestrator/INTEGRATE_WAVE_EFFORTS/rules.md; then
    # Remove the entire R020 section (3 lines)
    sed -i '/RULE R020 - State Transition Requirements/,+2d' /home/vscode/software-factory-template/agent-states/software-factory/orchestrator/INTEGRATE_WAVE_EFFORTS/rules.md
    FIXES_MADE=$((FIXES_MADE + 1))
    echo "   ✅ Removed R020 reference from INTEGRATE_WAVE_EFFORTS"
fi

# Fix COMPLETE_PHASE state - Remove R040 reference (non-existent)
echo "📝 Fixing COMPLETE_PHASE state..."
if grep -q "R040" /home/vscode/software-factory-template/agent-states/software-factory/orchestrator/COMPLETE_PHASE/rules.md; then
    # Remove the entire R040 section (4 lines)
    sed -i '/R040 - Documentation Requirements/,+3d' /home/vscode/software-factory-template/agent-states/software-factory/orchestrator/COMPLETE_PHASE/rules.md
    FIXES_MADE=$((FIXES_MADE + 1))
    echo "   ✅ Removed R040 reference from COMPLETE_PHASE"
fi

# Fix SPAWN_SW_ENGINEERS state - Remove R013, R017, R060 from acknowledgment list
echo "📝 Fixing SPAWN_SW_ENGINEERS state..."
if grep -q "R013\|R017\|R060" /home/vscode/software-factory-template/agent-states/software-factory/orchestrator/SPAWN_SW_ENGINEERS/rules.md; then
    # Update the acknowledgment line to remove non-existent rules
    sed -i 's/Acknowledge rules R054, R007, R013, R060, R017, R152, R295/Acknowledge rules R054, R007, R152, R295/g' /home/vscode/software-factory-template/agent-states/software-factory/orchestrator/SPAWN_SW_ENGINEERS/rules.md
    FIXES_MADE=$((FIXES_MADE + 1))
    echo "   ✅ Removed R013, R017, R060 from SPAWN_SW_ENGINEERS acknowledgment list"
fi

echo ""
echo "═══════════════════════════════════════════════════"
echo "📊 FIX SUMMARY:"
echo "   Total fixes applied: $FIXES_MADE"
echo "═══════════════════════════════════════════════════"