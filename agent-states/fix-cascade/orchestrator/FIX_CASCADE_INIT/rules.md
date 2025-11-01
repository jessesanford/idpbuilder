# Fix Cascade Orchestrator - FIX_CASCADE_INIT State Rules

## 🔴🔴🔴 ABSOLUTE REQUIREMENT: STATE MANAGER CONSULTATION 🔴🔴🔴

**SUPREME LAW - R517 - UNIVERSAL STATE MANAGER CONSULTATION LAW**

**BEFORE EXITING THIS STATE, YOU MUST:**

1. **SPAWN STATE MANAGER FOR SHUTDOWN_CONSULTATION** (MANDATORY - NO EXCEPTIONS)
2. **NEVER UPDATE STATE FILES DIRECTLY** (orchestrator-state-v3.json, bug-tracking.json, etc.)
3. **NEVER COMMIT STATE FILES YOURSELF** (State Manager does this atomically)

**FAILURE TO CONSULT STATE MANAGER = IMMEDIATE SYSTEM HALT (-100% GRADE)**

### Enforcement Mechanism

If you attempt to exit this state without spawning State Manager:
- ❌ Pre-commit hooks will REJECT your commit
- ❌ Validation tools will FAIL the build
- ❌ Grading system will assign -100% penalty
- ❌ System will transition to ERROR_RECOVERY

### Required Pattern (COPY THIS EXACTLY)

```bash
# At end of state work, BEFORE any state file updates:

echo "🔴 MANDATORY: Spawning State Manager for SHUTDOWN_CONSULTATION"

# Spawn State Manager (REQUIRED - NOT OPTIONAL)
# Task: state-manager
# State: SHUTDOWN_CONSULTATION
# Current State: [YOUR_CURRENT_STATE]
# Proposed Next State: [YOUR_PROPOSED_NEXT_STATE]
# Work Summary: [SUMMARY_OF_WORK_COMPLETED]

# State Manager will:
# 1. Validate proposed transition against state machine
# 2. Update all 4 state files atomically
# 3. Commit with [R288] tag
# 4. Return REQUIRED next state (may differ from proposal)

# Wait for State Manager response
# Follow State Manager's directive (REQUIRED next state)
# DO NOT proceed until State Manager confirms
```

**YOU MUST NEVER:**
- ❌ Update orchestrator-state-v3.json yourself
- ❌ Update bug-tracking.json yourself
- ❌ Update integration-containers.json yourself
- ❌ Use `jq` to modify state files
- ❌ Use `sed/awk` to modify state files
- ❌ Set `validated_by: "orchestrator"` (must be "state-manager")
- ❌ Skip State Manager consultation "just this once"
- ❌ Think "I'll validate it manually"

**ONLY State Manager may update state files. This is NON-NEGOTIABLE.**

See: `rule-library/R517-universal-state-manager-consultation-law.md`
## 🔴🔴🔴 SUB-STATE MACHINE ENTRY POINT 🔴🔴🔴

**YOU HAVE ENTERED THE FIX CASCADE SUB-STATE MACHINE**

This is a DIVERSION from the main state machine to handle fix cascade operations.

## 🔴🔴🔴 R375 ENFORCEMENT - DUAL STATE FILES 🔴🔴🔴

**MANDATORY FIRST ACTION:**
1. Create fix-specific state file per R375
2. NEVER pollute main orchestrator-state-v3.json with fix details
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
    "state_file": "orchestrator-state-v3.json",
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
- From main MONITORING_SWE_PROGRESS with hotfix requirement
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
jq '.state_machine.current_state = "FIX_CASCADE_ANALYSIS" |
    .state_machine.previous_state = "FIX_CASCADE_INIT" |
    .transition_time = now' orchestrator-${FIX_ID}-state.json > tmp.json && \
    mv tmp.json orchestrator-${FIX_ID}-state.json

# Commit the transition
git add orchestrator-${FIX_ID}-state.json
git commit -m "fix-cascade: ${FIX_ID} - transition to ANALYSIS"
git push
```

## Integration with Main State
The main orchestrator-state-v3.json should show:
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

## 🔴🔴🔴 R405 - MANDATORY AUTOMATION CONTINUATION FLAG 🔴🔴🔴

**SUPREME LAW - PENALTY FOR VIOLATION: -100% GRADE**


### 🚨 CRITICAL DISTINCTION: AGENT STOPS ≠ FACTORY STOPS 🚨

**TWO INDEPENDENT DECISIONS - DO NOT CONFUSE THEM:**

#### 1. Should Agent Stop Work? (R322 Technical Requirement)
- Agent completes current state
- Agent saves TODOs and commits state
- Agent exits with `exit 0` (preserves context)
- User runs /continue-orchestrating to resume
- **This is NORMAL at checkpoints**

#### 2. Should Factory Continue? (R405 Operational Status)
- Even though agent stopped, can automation proceed?
- TRUE = Healthy completion, automation can continue
- FALSE = Catastrophic failure, must halt everything
- **R322 checkpoints = TRUE (99.9% of cases)**

### THE PATTERN AT R322 CHECKPOINTS

```bash
# 1. Complete state work
echo "✅ State work complete"

# 2. Update state file
jq '.state_machine.current_state = "NEXT_STATE"' orchestrator-state-v3.json > tmp && mv tmp orchestrator-state-v3.json

# 3. Save TODOs
save_todos "R322_CHECKPOINT"

# 4. Factory continues (operational status)
echo "CONTINUE-SOFTWARE-FACTORY=TRUE"

# 5. Agent stops (technical requirement)
exit 0
```

**Both happen together! Agent stops AND factory continues!**

### WHEN TO USE EACH FLAG VALUE

**TRUE (99.9%):**
- ✅ R322 checkpoint reached
- ✅ State work completed successfully
- ✅ Ready for /continue-orchestrating
- ✅ Waiting for user to continue (NORMAL)
- ✅ Plan ready for review (agent done, factory proceeds)

**FALSE (0.1%):**
- ❌ CATASTROPHIC unrecoverable error
- ❌ Data corruption spreading
- ❌ Critical security violation
- ❌ NOT for R322 checkpoints
- ❌ NOT for user review needs
### 🚨 CASCADE OPERATIONS - CRITICAL R405 USAGE 🚨

**CASCADE transitions are a common source of incorrect FALSE usage!**

**CASCADE operations are AUTOMATED:**
```bash
# After deleting stale integration
# Transitioning to SETUP to recreate
# R322 checkpoint
# Flag? → MUST BE TRUE!

# Why? The system knows:
# - Current state from state file
# - What to do next
# - NO HUMAN DECISION NEEDED!

echo "CONTINUE-SOFTWARE-FACTORY=TRUE"  # CASCADE CONTINUES!
```
**See: $CLAUDE_PROJECT_DIR/rule-library/R405-automation-continuation-flag.md**

