# Fix Cascade Orchestrator - FIX_CASCADE_INIT State Rules

## 🔴🔴🔴 SUB-STATE MACHINE ENTRY POINT 🔴🔴🔴

**YOU HAVE ENTERED THE FIX CASCADE SUB-STATE MACHINE**

This is a DIVERSION from the main state machine to handle fix cascade operations.

## 🔴🔴🔴 R375 ENFORCEMENT - DUAL STATE FILES 🔴🔴🔴

**MANDATORY FIRST ACTION:**
1. Create fix-specific state file per R375
2. NEVER pollute main orchestrator-state.json with fix details
3. Track high-level progress in main state only

```bash
# Create fix state file IMMEDIATELY
FIX_ID="[extract from context or generate]"
cat > orchestrator-${FIX_ID}-state.json << 'EOF'
{
  "sub_state_type": "FIX_CASCADE",
  "fix_identifier": "${FIX_ID}",
  "current_state": "FIX_CASCADE_INIT",
  "created_at": "$(date -u +%Y-%m-%dT%H:%M:%SZ)",
  "parent_state_machine": {
    "state_file": "orchestrator-state.json",
    "return_state": "[from main state]",
    "nested_level": 1
  }
}
EOF
```

## 📋 PRIMARY DIRECTIVES FOR FIX_CASCADE_INIT

### State Purpose
Initialize the fix cascade sub-state machine and prepare for fix operations.

### Immediate Actions Upon Entry
1. **Check Parent State Machine**
   - Verify sub_state_machine.active = true in main state
   - Record return state for later
   - Note trigger reason

2. **Create Fix State File** (R375 MANDATORY)
   - Generate unique fix identifier
   - Create orchestrator-[fix-id]-state.json
   - Initialize with fix cascade structure

3. **Load Fix Requirements**
   - Determine fix type (HOTFIX, BACKPORT, FORWARD_PORT)
   - Identify source branch and commit
   - List target branches

4. **Set Up Tracking**
   - Initialize quality gate tracking
   - Prepare cascade plan structure
   - Set up validation checklist

## 🔴🔴🔴 R376 QUALITY GATES AWARENESS 🔴🔴🔴

**UNDERSTAND THE GATES YOU MUST PASS:**
1. **Gate 1**: Post-Backport Review (after EACH backport)
2. **Gate 2**: Post-Forward-Port Review (after EACH forward port)
3. **Gate 3**: Conflict Resolution Review (after ANY conflict)
4. **Gate 4**: Comprehensive Validation (before completion)

**VIOLATION OF QUALITY GATES = -100% AUTOMATIC FAILURE**

## State Transition Rules

### Valid Transitions From FIX_CASCADE_INIT
- → FIX_CASCADE_ANALYSIS (always, after initialization)

### Entry Conditions
- From main ERROR_RECOVERY with fix cascade trigger
- From main MONITORING with hotfix requirement
- From /fix-cascade command

### Exit Conditions
- Fix state file created successfully
- Initial requirements loaded
- Ready for analysis

## Required Validations
1. ✅ Fix state file created and committed
2. ✅ Main state updated with sub-state reference
3. ✅ Fix identifier is unique
4. ✅ Parent state properly recorded

## Enforcement Checklist
```markdown
- [ ] Created fix-specific state file
- [ ] Named per R375 convention
- [ ] Committed within 30 seconds
- [ ] Main state references sub-state
- [ ] Return state recorded
- [ ] Fix requirements documented
```

## Example State Transition
```bash
# Update fix state to transition to ANALYSIS
jq '.current_state = "FIX_CASCADE_ANALYSIS" |
    .previous_state = "FIX_CASCADE_INIT" |
    .transition_time = now' orchestrator-${FIX_ID}-state.json > tmp.json && \
    mv tmp.json orchestrator-${FIX_ID}-state.json

# Commit the transition
git add orchestrator-${FIX_ID}-state.json
git commit -m "fix-cascade: ${FIX_ID} - transition to ANALYSIS"
git push
```

## Integration with Main State
The main orchestrator-state.json should show:
```json
{
  "sub_state_machine": {
    "active": true,
    "type": "FIX_CASCADE",
    "state_file": "orchestrator-[fix-id]-state.json",
    "current_state": "FIX_CASCADE_INIT",
    "return_state": "[original state]",
    "started_at": "[timestamp]",
    "trigger_reason": "[why fix cascade started]"
  }
}
```

## Common Errors to Avoid
- ❌ Putting fix details in main state file
- ❌ Not creating fix state file immediately
- ❌ Forgetting to record return state
- ❌ Not committing state changes
- ❌ Skipping quality gate setup

## Notes
- This is the entry point to fix cascade operations
- All subsequent fix work uses the fix-specific state file
- Main state only tracks that we're in a sub-state machine
- When complete, will archive fix state and return to main