#!/bin/bash

# Script to fix missing rule references in orchestrator state files

echo "🔧 Starting fix for missing rule references..."

# Fix R020 -> Should likely be R322 (state transitions)
echo "Fixing R020 references..."
sed -i 's/RULE R020 - State Transition Requirements/RULE R322 - Mandatory Stop Before State Transitions/g' /home/vscode/software-factory-template/agent-states/software-factory/orchestrator/INTEGRATE_WAVE_EFFORTS/rules.md
sed -i 's|R020-state-transitions.md|R322-mandatory-stop-before-state-transitions.md|g' /home/vscode/software-factory-template/agent-states/software-factory/orchestrator/INTEGRATE_WAVE_EFFORTS/rules.md

# Fix R040 -> Should likely be R035 (phase completion testing)
echo "Fixing R040 references..."
sed -i 's/RULE R040/RULE R035/g' /home/vscode/software-factory-template/agent-states/software-factory/orchestrator/COMPLETE_PHASE/rules.md
sed -i 's|R040-[^.]*\.md|R035-phase-completion-testing.md|g' /home/vscode/software-factory-template/agent-states/software-factory/orchestrator/COMPLETE_PHASE/rules.md

# Fix R013 -> Likely should be R108 (code review protocol)
echo "Fixing R013 references..."
sed -i 's/RULE R013/RULE R108/g' /home/vscode/software-factory-template/agent-states/software-factory/orchestrator/SPAWN_SW_ENGINEERS/rules.md
sed -i 's|R013-[^.]*\.md|R108-code-review-protocol.md|g' /home/vscode/software-factory-template/agent-states/software-factory/orchestrator/SPAWN_SW_ENGINEERS/rules.md

# Fix R017 -> Likely should be R152 (implementation speed)
echo "Fixing R017 references..."
sed -i 's/RULE R017/RULE R152/g' /home/vscode/software-factory-template/agent-states/software-factory/orchestrator/SPAWN_SW_ENGINEERS/rules.md
sed -i 's|R017-[^.]*\.md|R152-implementation-speed.md|g' /home/vscode/software-factory-template/agent-states/software-factory/orchestrator/SPAWN_SW_ENGINEERS/rules.md

# Fix R060 -> Likely should be R054 (implementation plan creation)
echo "Fixing R060 references..."
sed -i 's/RULE R060/RULE R054/g' /home/vscode/software-factory-template/agent-states/software-factory/orchestrator/SPAWN_SW_ENGINEERS/rules.md
sed -i 's|R060-[^.]*\.md|R054-implementation-plan-creation.md|g' /home/vscode/software-factory-template/agent-states/software-factory/orchestrator/SPAWN_SW_ENGINEERS/rules.md

echo "✅ Fixes applied. Verifying changes..."