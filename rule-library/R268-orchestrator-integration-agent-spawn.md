# 🚨🚨🚨 RULE R268: Orchestrator Integration Agent Spawn Protocol [SUPERSEDED BY R329] 🚨🚨🚨

## ⚠️⚠️⚠️ IMPORTANT: THIS RULE HAS BEEN SUPERSEDED BY R329 ⚠️⚠️⚠️

**R329 makes Integration Agent MANDATORY for ALL merges - NO EXCEPTIONS**
See: `$CLAUDE_PROJECT_DIR/rule-library/R329-orchestrator-never-performs-merges.md`

## Rule Definition
**Criticality:** BLOCKING (Updated per R329) - MANDATORY for ALL integrations
**Category:** Agent-Specific
**Applies To:** orchestrator, integration-agent
**Status:** SUPERSEDED - R329 takes precedence

## When to Spawn Integration Agent

### ALWAYS Spawn Integration Agent (R329 MANDATORY):
Per R329, the orchestrator MUST spawn Integration Agent for:
1. **ALL Merges** - Even single branch merges
2. **Simple Merges** - No exceptions for "simple" cases
3. **Complex Merges** - Multiple branches with dependencies
4. **Test Merges** - Even quick validation merges
5. **Any git merge operation** - 100% delegation required

### NEVER Handle Directly (R329 VIOLATION):
The orchestrator MUST NEVER:
1. ❌ Execute "simple merges" directly
2. ❌ Perform "quick integrations" itself
3. ❌ Handle "linear history" merges
4. ❌ Do "test integrations" directly
5. ❌ Execute ANY merge operation itself

## Spawning Protocol

```bash
# Orchestrator decision logic
BRANCH_COUNT=$(echo "$EFFORT_BRANCHES" | wc -w)
HAS_SPLITS=$(echo "$EFFORT_BRANCHES" | grep -c "split" || true)

if [[ $BRANCH_COUNT -ge 3 || $HAS_SPLITS -gt 0 ]]; then
    echo "📋 Complex integration detected - spawning integration agent"
    
    # Prepare integration agent task
    cat > integration-task.md << 'EOF'
# Integration Task

## Target Branches
$(for branch in $EFFORT_BRANCHES; do echo "- $branch"; done)

## Integration Goals
Merge all wave ${WAVE} efforts into a single integration branch while:
- Preserving complete git history
- Resolving all conflicts
- Documenting every operation
- Testing the integrated result

## Expected Outputs
- INTEGRATION-PLAN.md
- work-log.md (replayable)
- INTEGRATION-REPORT.md
- Integration branch: phase${PHASE}/wave${WAVE}-integration
EOF

    # Spawn integration agent
    /integrate-branches \
        TARGET_BRANCHES="$EFFORT_BRANCHES" \
        TARGET_BASE="main" \
        INTEGRATION_GOALS="Wave ${WAVE} integration"
        
    # Transition to waiting state
    NEXT_STATE="WAITING_FOR_INTEGRATION"
else
    echo "📋 Simple integration - handling directly"
    # Continue with direct integration per existing rules
fi
```

## Integration Agent Response

When spawned, the integration agent will:
1. Acknowledge grading criteria (50% completeness, 50% documentation)
2. Create comprehensive INTEGRATION-PLAN.md
3. Execute integration with meticulous tracking
4. Never modify original branches
5. Document all upstream bugs (not fix them)
6. Produce final INTEGRATION-REPORT.md

## Orchestrator Verification

After integration agent completes:

```bash
# Verify integration agent outputs
verify_integration_outputs() {
    local integration_dir="$1"
    
    # Check required documents exist
    for doc in INTEGRATION-PLAN.md work-log.md INTEGRATION-REPORT.md; do
        if [[ ! -f "$integration_dir/$doc" ]]; then
            echo "❌ Missing required document: $doc"
            return 1
        fi
    done
    
    # Check integration branch exists
    local integration_branch="phase${PHASE}/wave${WAVE}-integration"
    if ! git branch -r | grep -q "$integration_branch"; then
        echo "❌ Integration branch not found on remote"
        return 1
    fi
    
    # Review integration report for issues
    if grep -q "CRITICAL" "$integration_dir/INTEGRATION-REPORT.md"; then
        echo "⚠️ Critical issues found in integration report"
        # Transition to ERROR_RECOVERY if needed
    fi
    
    echo "✅ Integration complete and verified"
}
```

## State Machine Integration

### New States (Optional)
- `SPAWN_INTEGRATION_AGENT` - Spawning integration specialist
- `WAITING_FOR_INTEGRATION` - Waiting for integration completion

### State Transitions
```
WAVE_COMPLETE 
    → [Complex?] → SPAWN_INTEGRATION_AGENT
    → [Simple?] → INTEGRATION (direct)

SPAWN_INTEGRATION_AGENT 
    → WAITING_FOR_INTEGRATION

WAITING_FOR_INTEGRATION
    → [Success] → WAVE_REVIEW
    → [Failed] → ERROR_RECOVERY
```

## Benefits of Using Integration Agent

1. **Expertise** - Specialized in git operations and merging
2. **Documentation** - Meticulous tracking of every operation
3. **Replayability** - Work log can reproduce exact integration
4. **Safety** - Never modifies original branches
5. **Clarity** - Comprehensive reports on issues found
6. **Grading** - Clear success metrics (50/50 split)

## Related Rules
- R260-R267 - Integration Agent rules
- R034 - Integration Requirements
- R250 - Integration Location Isolation
- R258 - Mandatory Wave Review Report