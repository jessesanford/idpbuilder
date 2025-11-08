#!/bin/bash
# R510 Exit Validation for INTEGRATE_WAVE_EFFORTS State
# Ensures integration agent was actually spawned before state exit

set -e

PHASE=${1:-$(jq -r '.project_progression.current_phase.phase_number' orchestrator-state-v3.json)}
WAVE=${2:-$(jq -r '.project_progression.current_wave.wave_number' orchestrator-state-v3.json)}
ITERATION=${3:-$(jq -r '.project_progression.current_wave.iteration' orchestrator-state-v3.json)}

echo "🔍 Validating INTEGRATE_WAVE_EFFORTS exit criteria..."
echo "   Phase: $PHASE, Wave: $WAVE, Iteration: $ITERATION"

INTEGRATION_WORKSPACE="efforts/phase${PHASE}/wave${WAVE}/integration-workspace"
REPORT_FILE="${INTEGRATION_WORKSPACE}/.software-factory/phase${PHASE}/wave${WAVE}/integration/WAVE-${PHASE}.${WAVE}-INTEGRATION-REPORT-ITERATION-${ITERATION}.md"

# Check 1: Integration workspace exists
if [ ! -d "$INTEGRATION_WORKSPACE" ]; then
    echo "❌ VIOLATION: Integration workspace does not exist!"
    echo "   Expected: $INTEGRATION_WORKSPACE"
    echo "   This means integration agent was NEVER spawned"
    echo ""
    echo "Per R510: BLOCKING items must be completed before exit"
    echo "Per R329: Orchestrator must spawn integration agent"
    exit 510
fi

# Check 2: Integration report exists for current iteration
if [ ! -f "$REPORT_FILE" ]; then
    echo "❌ VIOLATION: Integration report for iteration $ITERATION does not exist!"
    echo "   Expected: $REPORT_FILE"
    echo "   This means integration agent was NOT spawned for this iteration"
    echo ""
    echo "Found reports:"
    find "$INTEGRATION_WORKSPACE/.software-factory" -name "*INTEGRATION-REPORT*" 2>/dev/null || echo "   (none)"
    echo ""
    echo "Per R510: Integration agent must be spawned (checklist item #5)"
    echo "Per R329: Orchestrator cannot perform merges - must delegate"
    exit 510
fi

# Check 3: Report indicates successful integration
INTEGRATION_STATUS=$(grep -o "Status.*: .*" "$REPORT_FILE" | head -1 | cut -d: -f2 | tr -d ' ')
if [ "$INTEGRATION_STATUS" != "✅CLEAN" ] && [ "$INTEGRATION_STATUS" != "CLEAN" ]; then
    echo "⚠️  WARNING: Integration status is not CLEAN: $INTEGRATION_STATUS"
    echo "   This may require fix cascade - check report for details"
fi

# Check 4: Integration branch has commits from this iteration
cd "$INTEGRATION_WORKSPACE"
COMMITS_THIS_ITERATION=$(git log --oneline --since="1 hour ago" | wc -l)
if [ "$COMMITS_THIS_ITERATION" -eq 0 ]; then
    echo "⚠️  WARNING: No commits in integration workspace in last hour"
    echo "   Integration may be stale - verify agent actually ran"
fi

echo "✅ INTEGRATE_WAVE_EFFORTS exit validation passed"
echo "   - Integration workspace exists"
echo "   - Integration report for iteration $ITERATION exists"
echo "   - Integration agent was properly spawned"
exit 0
