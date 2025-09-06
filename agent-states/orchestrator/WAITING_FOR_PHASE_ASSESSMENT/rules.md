# Orchestrator - WAITING_FOR_PHASE_ASSESSMENT State Rules

## 🛑🛑🛑 R322 MANDATORY STOP BEFORE STATE TRANSITIONS 🛑🛑🛑

**SUPREME LAW - VIOLATION = -100% IMMEDIATE FAILURE**

### YOU MUST STOP AFTER:
1. ✅ Completing all TODOs for this state
2. ✅ Updating orchestrator-state.yaml with new state
3. ✅ Committing and pushing the state file  
4. ✅ Providing work summary

### YOU MUST NOT:
- ❌ Continue to the next state automatically
- ❌ Start work for the new state
- ❌ Spawn agents for the new state
- ❌ Assume permission to continue

### STOP PROTOCOL:
```markdown
## 🛑 STATE TRANSITION CHECKPOINT: CURRENT_STATE → NEXT_STATE

### ✅ Current State Work Completed:
- [List completed work]

### 📊 Current Status:
- Current State: CURRENT_STATE
- Next State: NEXT_STATE
- TODOs Completed: X/Y
- State Files: Updated and committed ✅

### ⏸️ STOPPED - Awaiting User Continuation
Ready to transition to NEXT_STATE. Please use /continue-orchestrating.
```

**STOP MEANS STOP - Exit and wait for /continue-orchestrating**

---


## 🔴🔴🔴 STOP! STATE RULE READING IS ABSOLUTELY FIRST! 🔴🔴🔴

**YOU HAVE ENTERED WAITING_FOR_PHASE_ASSESSMENT STATE - YOU MUST READ AND ACKNOWLEDGE ALL STATE RULES BEFORE DOING ANY STATE WORK!**

## 🔴🔴🔴 R290 VERIFICATION REQUIREMENT 🔴🔴🔴

**R290 ENFORCEMENT: CREATE VERIFICATION MARKER AFTER READING**

After reading and acknowledging all state rules, you MUST create a verification marker:

```bash
# MANDATORY: Create verification marker after reading rules
touch .state_rules_read_orchestrator_WAITING_FOR_PHASE_ASSESSMENT
echo "$(date +%s) - Rules read and acknowledged for WAITING_FOR_PHASE_ASSESSMENT" > .state_rules_read_orchestrator_WAITING_FOR_PHASE_ASSESSMENT
```

**FAILURE TO CREATE MARKER = AUTOMATIC -100% PENALTY**

The system will check for this marker. No marker = Immediate failure.

### ❌ DO NOT DO ANY WAITING_FOR_PHASE_ASSESSMENT WORK UNTIL RULES ARE READ:
- ❌ Start monitor assessment progress
- ❌ Start check architect status
- ❌ Start wait for phase evaluation
- ❌ Update state files
- ❌ Continue to next state
- ❌ Think about what to do in this state

### ✅ YOU MUST IMMEDIATELY:

## 🔴🔴🔴 MANDATORY STATE RULE READING AND ACKNOWLEDGMENT 🔴🔴🔴

### ⚠️⚠️⚠️ YOU MUST READ EACH RULE FILE LISTED IN PRIMARY DIRECTIVES. **I AM WATCHING YOUR TOOL CALLS FOR READ OPERATIONS** *YOU WILL FAIL* IF YOU DO NOT MAKE A READ FILE CALL FOR EACH RULE FILE IN PRIMARY DIRECTIVES!!! ⚠️⚠️⚠️

**AFTER READING, YOU MUST ACKNOWLEDGE ALL THE STATE RULES AND STATE THAT YOU WILL ABIDE BY THEM ONE AT A TIME GIVING THE RULE NUMBER AND DESCRIPTION.**

### ❌ ANTI-PATTERNS THAT WILL CAUSE FAILURE:

1. **Fake Acknowledgment Without Reading**:
   ```
   ❌ WRONG: "I acknowledge R151, R208, R053..."
   (No Read tool calls detected = AUTOMATIC FAILURE)
   ```

2. **Bulk Acknowledgment**:
   ```
   ❌ WRONG: "I acknowledge all WAITING_FOR_PHASE_ASSESSMENT rules"
   (YOU Must READ AND ACKNOWLEDGE EACH rule individually)
   ```

3. **Silent Reading**:
   ```
   ❌ WRONG: [Reads rules but doesn't acknowledge]
   "Now I've read the rules, let me start work..."
   (MUST explicitly acknowledge EACH rule)
   ```

4. **Reading From Memory**:
   ```
   ❌ WRONG: "I know R208 requires CD before spawn..."
   (Must READ from file, not recall from memory)
   ```

5. **Skipping Rules in PRIMARY DIRECTIVES**:
   ```
   ❌ WRONG: Reading only some rules from the list
   (ALL rules in PRIMARY DIRECTIVES are MANDATORY)
   ```

### ✅ CORRECT PATTERN FOR WAITING_FOR_PHASE_ASSESSMENT:
```
1. READ: $CLAUDE_PROJECT_DIR/rule-library/[first-rule-file].md
2. "I acknowledge [Rule#] - [Rule Name]: [Brief description]"
3. READ: $CLAUDE_PROJECT_DIR/rule-library/[second-rule-file].md  
4. "I acknowledge [Rule#] - [Rule Name]: [Brief description]"
[Continue for EVERY rule in PRIMARY DIRECTIVES...]
5. "Ready to execute WAITING_FOR_PHASE_ASSESSMENT work"
```

### 🚨 NO WORK UNTIL ACKNOWLEDGMENT COMPLETE 🚨
**You may NOT begin ANY WAITING_FOR_PHASE_ASSESSMENT work until:**
1. ✅ ALL rules in PRIMARY DIRECTIVES have been READ
2. ✅ ALL rules have been individually ACKNOWLEDGED
3. ✅ You have stated readiness to execute WAITING_FOR_PHASE_ASSESSMENT work
1. **READ** every rule file listed in PRIMARY DIRECTIVES below
2. **ACKNOWLEDGE** each rule individually with number and description
3. **ONLY THEN** proceed with WAITING_FOR_PHASE_ASSESSMENT work

### 🚨 FAILURE TO READ STATE RULES FIRST = IMMEDIATE EXIT 🚨
**If you do ANY WAITING_FOR_PHASE_ASSESSMENT work before reading and acknowledging rules:**
- **STOP ALL WORK IMMEDIATELY**
- **EXIT WITH FAILURE STATUS**
- **YOU HAVE VIOLATED STATE COMPLIANCE**

**THE SYSTEM IS MONITORING YOUR READ TOOL CALLS!**

## 🚨 WAITING_FOR_PHASE_ASSESSMENT IS A VERB - ACTIVELY MONITOR! 🚨

## 🚨🚨🚨 MANDATORY: VERIFY PHASE ASSESSMENT REPORT EXISTS [R257] 🚨🚨🚨

**CRITICAL REQUIREMENT**: The architect MUST create a `PHASE-{N}-ASSESSMENT-REPORT.md` file before you can proceed!

### IMMEDIATE ACTIONS UPON ENTERING WAITING_FOR_PHASE_ASSESSMENT

**THE MOMENT YOU ENTER THIS STATE, YOU MUST:**
1. Check if architect has already responded
2. **VERIFY ASSESSMENT REPORT FILE EXISTS** at `phase-assessments/phase{N}/PHASE-{N}-ASSESSMENT-REPORT.md`
3. Read and validate the assessment report content
4. Extract the DECISION from the report
5. Monitor for architect decision actively
6. Process any pending TodoWrite items
7. Prepare for phase completion or error recovery

**FORBIDDEN - AUTOMATIC FAILURE:**
- ❌ "STATE TRANSITION COMPLETE: Now waiting" [stops]
- ❌ "I'll wait for the architect" [passive waiting]
- ❌ "Standing by for assessment" [does nothing]

**REQUIRED - IMMEDIATE ACTION:**
- ✅ "Checking for architect response..."
- ✅ "Monitoring phase assessment status..."
- ✅ "Processing pending tasks while waiting..."

## 🔴🔴🔴 MANDATORY MEASUREMENT RULE - R304 🔴🔴🔴

### R304: Mandatory Line Counter Tool Enforcement
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R304-mandatory-line-counter-enforcement.md`
**Criticality**: BLOCKING - Manual counting = AUTOMATIC -100% FAILURE

**ABSOLUTE REQUIREMENTS:**
- ✅ MUST use `$CLAUDE_PROJECT_DIR/tools/line-counter.sh` for ALL line counting
- ❌ NEVER use `wc -l` or any manual counting method
- ❌ NEVER count lines any other way - this is a -100% automatic failure
- ✅ MUST specify both -b (base branch) and -c (current branch) parameters
- ✅ Base branch MUST be phase integration branch (NOT "main")

**Failure to use the line counter tool = AUTOMATIC -100% GRADE**

## State Context

Waiting for architect's phase-level assessment that determines if the phase can be marked complete (SUCCESS) or needs more work.

## Primary Purpose

The WAITING_FOR_PHASE_ASSESSMENT state is for:
1. Actively monitoring architect assessment progress
2. Processing architect's phase decision
3. Determining transition to PHASE_COMPLETE or ERROR_RECOVERY
4. Ensuring no premature SUCCESS without approval

## Processing Architect Decision

### MANDATORY FIRST STEP - Verify Assessment Report:
```bash
# MUST verify assessment report exists per R257
PHASE="$CURRENT_PHASE"
REPORT_FILE="phase-assessments/phase${PHASE}/PHASE-${PHASE}-ASSESSMENT-REPORT.md"

if [ ! -f "$REPORT_FILE" ]; then
    echo "❌ CRITICAL: No phase assessment report found!"
    echo "❌ Expected: $REPORT_FILE"
    echo "❌ Cannot proceed without assessment report per R257"
    transition_to "ERROR_RECOVERY"
    exit 1
fi

# Extract decision from report
DECISION=$(grep "^\*\*DECISION\*\*:" "$REPORT_FILE" | cut -d: -f2 | xargs)
if [ -z "$DECISION" ]; then
    echo "❌ Invalid assessment report - no decision found!"
    transition_to "ERROR_RECOVERY"
    exit 1
fi

echo "✅ Phase assessment report verified: $REPORT_FILE"
echo "📊 Decision: $DECISION"
```

### If PHASE_COMPLETE (Assessment Passed):
```bash
# Update state file with report location
yq -i ".phase_assessment.status = \"COMPLETE\"" orchestrator-state.yaml
yq -i ".phase_assessment.completed_at = \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\"" orchestrator-state.yaml
yq -i ".phase_assessment.decision = \"PHASE_COMPLETE\"" orchestrator-state.yaml
yq -i ".phase_assessment.report_file = \"$REPORT_FILE\"" orchestrator-state.yaml
yq -i ".phase_assessment.feedback = \"$FEEDBACK\"" orchestrator-state.yaml

# Phase approved - can now complete
transition_to "PHASE_COMPLETE"  # Handle phase completion tasks
```

### If NEEDS_WORK (Fixes Required):
```bash
# Update state file with decision and report location
yq -i ".phase_assessment.status = \"NEEDS_WORK\"" orchestrator-state.yaml
yq -i ".phase_assessment.decision = \"NEEDS_WORK\"" orchestrator-state.yaml
yq -i ".phase_assessment.report_file = \"$REPORT_FILE\"" orchestrator-state.yaml

# CRITICAL: Extract specific issues from the report (R257)
echo "📋 Extracting required fixes from assessment report..."
PRIORITY_1_FIXES=$(sed -n '/### Priority 1/,/### Priority 2/p' "$REPORT_FILE" | grep "^- \[")
ISSUES_IDENTIFIED=$(sed -n '/## Issues Identified/,/## Required Fixes/p' "$REPORT_FILE")

# Store issues in state file for ERROR_RECOVERY to process
yq -i ".phase_assessment.priority_1_fixes = \"$PRIORITY_1_FIXES\"" orchestrator-state.yaml
yq -i ".phase_assessment.issues_identified = \"$ISSUES_IDENTIFIED\"" orchestrator-state.yaml
yq -i ".error_recovery.source = \"PHASE_ASSESSMENT_NEEDS_WORK\"" orchestrator-state.yaml
yq -i ".error_recovery.report_to_read = \"$REPORT_FILE\"" orchestrator-state.yaml

echo "❌ Phase assessment returned NEEDS_WORK - transitioning to ERROR_RECOVERY"
echo "📄 ERROR_RECOVERY must read report: $REPORT_FILE"

# Must fix before phase can complete
transition_to "ERROR_RECOVERY"  # ERROR_RECOVERY will read report and coordinate fixes
```

### If PHASE_FAILED (Cannot Complete):
```bash
# Update state file
yq -i ".phase_assessment.status = \"FAILED\"" orchestrator-state.yaml
yq -i ".phase_assessment.failure_reason = \"$REASON\"" orchestrator-state.yaml

# Phase cannot complete - major rework needed
transition_to "HARD_STOP"  # Critical failure
```

## State Transitions

From WAITING_FOR_PHASE_ASSESSMENT:
- **PHASE_COMPLETE decision** → PHASE_COMPLETE (proceed to completion)
- **NEEDS_WORK decision** → ERROR_RECOVERY (fix issues)
- **PHASE_FAILED decision** → HARD_STOP (critical failure)
- **Timeout** → ERROR_RECOVERY (handle timeout)

**CRITICAL**: This state BLOCKS the path to SUCCESS. No phase can complete without passing through here and getting approval.

## Timeout Handling

```bash
# Check for timeout (e.g., 45 minutes for phase assessment)
REQUESTED_AT=$(yq '.phase_assessment.requested_at' orchestrator-state.yaml)
TIMEOUT_MINUTES=45

if timeout_exceeded "$REQUESTED_AT" "$TIMEOUT_MINUTES"; then
    yq -i ".phase_assessment.status = \"TIMEOUT\"" orchestrator-state.yaml
    log_error "Phase assessment timeout after $TIMEOUT_MINUTES minutes"
    transition_to "ERROR_RECOVERY"
fi
```

## Active Monitoring Requirements

While waiting, you must:
1. Check for architect response every interaction
2. Process any TodoWrite pending items
3. Monitor elapsed time for timeout
4. Prepare phase completion documentation
5. NOT passively wait doing nothing

## Phase Assessment Tracking

```yaml
# State file structure for phase assessment
phase_assessment:
  phase: 1
  requested_at: "2025-08-25T23:00:00Z"
  architect_spawned: "architect-789012"
  phase_branch: "tmc-workspace/phase1-complete"
  wave_count: 4
  report_file: "phase-assessments/phase1/PHASE-1-ASSESSMENT-REPORT.md"  # R257 MANDATORY
  decision: "PHASE_COMPLETE"  # Extracted from report file
  score: 92  # Extracted from report file
  feedback: |
    - All features implemented successfully
    - Architecture consistent across waves
    - APIs stable and well-documented
    - Test coverage exceeds requirements (85%)
    - Ready for production deployment
  completed_at: "2025-08-25T23:30:00Z"
  next_action: "PHASE_COMPLETE"
```

## Success Criteria

Successfully process assessment when:
- [ ] Architect provides clear decision
- [ ] Decision properly recorded in state file
- [ ] Appropriate next state determined
- [ ] Feedback captured for audit trail
- [ ] Metrics updated

## Common Issues

1. **Timeout**: Architect takes too long - escalate
2. **Unclear Decision**: Request clarification
3. **Missing Context**: Provide additional information
4. **Partial Assessment**: Request complete review

## Required State File Updates

Before leaving WAITING_FOR_PHASE_ASSESSMENT:
- [ ] Assessment decision recorded
- [ ] Completion timestamp set
- [ ] Feedback/issues documented
- [ ] Next state determined
- [ ] Phase metrics finalized

## Grading Impact

- Processing decision promptly: +10 points
- Proper state transitions: +20 points
- Skipping to SUCCESS without this: -100 points (CRITICAL)
- Timeout handling: +5 points

## R322 VIOLATION DETECTION

If you find yourself:
- Starting work for a new state without /continue-orchestrating
- Transitioning without stopping after state file commit
- Continuing after completing state work

**STOP IMMEDIATELY - You are violating R322!**
