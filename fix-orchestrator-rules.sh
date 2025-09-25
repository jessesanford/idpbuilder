#!/bin/bash

# Script to fix missing rule references in orchestrator state files
# Strategy: Remove references to non-existent rules

echo "🔧 Starting comprehensive fix for missing rule references..."
echo ""

FIXES_MADE=0

# Fix INTEGRATION state - Remove R020 reference (non-existent)
echo "📝 Fixing INTEGRATION state..."
if grep -q "R020" /home/vscode/software-factory-template/agent-states/orchestrator/INTEGRATION/rules.md; then
    # Remove the entire R020 section (3 lines)
    sed -i '/RULE R020 - State Transition Requirements/,+2d' /home/vscode/software-factory-template/agent-states/orchestrator/INTEGRATION/rules.md
    FIXES_MADE=$((FIXES_MADE + 1))
    echo "   ✅ Removed R020 reference from INTEGRATION"
fi

# Fix PHASE_COMPLETE state - Remove R040 reference (non-existent)
echo "📝 Fixing PHASE_COMPLETE state..."
if grep -q "R040" /home/vscode/software-factory-template/agent-states/orchestrator/PHASE_COMPLETE/rules.md; then
    # Remove the entire R040 section (4 lines)
    sed -i '/R040 - Documentation Requirements/,+3d' /home/vscode/software-factory-template/agent-states/orchestrator/PHASE_COMPLETE/rules.md
    FIXES_MADE=$((FIXES_MADE + 1))
    echo "   ✅ Removed R040 reference from PHASE_COMPLETE"
fi

# Fix SPAWN_AGENTS state - Remove R013, R017, R060 from acknowledgment list
echo "📝 Fixing SPAWN_AGENTS state..."
if grep -q "R013\|R017\|R060" /home/vscode/software-factory-template/agent-states/orchestrator/SPAWN_AGENTS/rules.md; then
    # Update the acknowledgment line to remove non-existent rules
    sed -i 's/Acknowledge rules R054, R007, R013, R060, R017, R152, R295/Acknowledge rules R054, R007, R152, R295/g' /home/vscode/software-factory-template/agent-states/orchestrator/SPAWN_AGENTS/rules.md
    FIXES_MADE=$((FIXES_MADE + 1))
    echo "   ✅ Removed R013, R017, R060 from SPAWN_AGENTS acknowledgment list"
fi

echo ""
echo "═══════════════════════════════════════════════════"
echo "📊 FIX SUMMARY:"
echo "   Total fixes applied: $FIXES_MADE"
echo "═══════════════════════════════════════════════════"