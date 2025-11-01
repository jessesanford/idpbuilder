# REVIEW_WAVE_INTEGRATION State Rules

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

## 🔴🔴🔴 MANDATORY STATE RULE READING AND ACKNOWLEDGMENT 🔴🔴🔴

### ⚠️⚠️⚠️ YOU MUST READ EACH RULE FILE LISTED IN PRIMARY DIRECTIVES. **I AM WATCHING YOUR TOOL CALLS FOR READ OPERATIONS** *YOU WILL FAIL* IF YOU DO NOT MAKE A READ FILE CALL FOR EACH RULE FILE IN PRIMARY DIRECTIVES!!! ⚠️⚠️⚠️

**AFTER READING, YOU MUST ACKNOWLEDGE ALL THE STATE RULES AND STATE THAT YOU WILL ABIDE BY THEM ONE AT A TIME GIVING THE RULE NUMBER AND DESCRIPTION.**

## 📋 PRIMARY DIRECTIVES FOR REVIEW_WAVE_INTEGRATION STATE

### Core Mandatory Rules (ALL orchestrator states must have these):

1. **🚨🚨🚨 R006** - ORCHESTRATOR NEVER WRITES CODE OR PERFORMS FILE OPERATIONS (BLOCKING)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R006-orchestrator-never-writes-code.md`
   - Criticality: BLOCKING - Automatic termination, 0% grade
   - Summary: NEVER write, copy, move, or manipulate ANY code files - delegate ALL to agents

2. **🔴🔴🔴 R287** - TODO PERSISTENCE COMPREHENSIVE (SUPREME LAW)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R287-todo-persistence-comprehensive.md`
   - Criticality: SUPREME - -20% to -100% penalty for violations
   - Summary: MUST save TODOs within 30s after write, every 10 messages, before transitions

3. **🔴🔴🔴 R288** - STATE FILE UPDATE REQUIREMENTS (SUPREME LAW)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R288-state-file-update-and-commit-protocol.md`
   - Criticality: SUPREME - State updates required for all transitions
   - Summary: MUST update orchestrator-state-v3.json before EVERY state transition

4. **🔴🔴🔴 R322 Part A** - Mandatory Stop After Spawn States
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R322-mandatory-stop-before-state-transitions.md`
   - Criticality: SUPREME LAW - Must stop after spawning
   - Summary: ALL spawn states require STOP after spawning agents

### State-Specific Rules:

5. **🔴🔴🔴 R304** - Mandatory Line Counting Tool
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R304-mandatory-line-counter-enforcement.md`
   - Criticality: SUPREME LAW - Line counting requirements
   - Summary: MUST use tools/line-counter.sh for ALL measurements

## State Purpose
Spawn Code Reviewer to perform quality review of integrated code after wave integration is complete. This state ensures the merged code maintains quality standards and catches integration-introduced issues before architect review.

## Entry Criteria
- **From**: MONITORING_INTEGRATE_WAVE_EFFORTS
- **Condition**: Integration Agent has completed wave integration successfully
- **Required**: 
  - Integration workspace exists with merged code
  - All effort branches merged into wave integration branch
  - Basic merge conflicts resolved by Integration Agent

## State Actions

### 1. IMMEDIATE: Spawn Code Reviewer for Integration Review
```bash
# Spawn Code Reviewer to review integrated code
/spawn agent-code-reviewer REVIEW_WAVE_INTEGRATION \
  --wave "W${current_wave}" \
  --branch "${wave_integration_branch}" \
  --focus "integration-quality"
```

### 2. Code Reviewer Responsibilities
The spawned Code Reviewer will:
- Review the quality of integrated code
- Check for integration-introduced bugs
- Validate merge conflict resolutions
- Ensure code consistency across merged efforts
- Verify tests still pass after integration
- Check for duplicate code or functionality conflicts
- **Verify demo scripts are present and executable (R330/R291)**
- **Check demo coverage of integrated features**
- **Validate demo documentation completeness**
- **Review demo execution results from integration testing**
- Create REVIEW_WAVE_INTEGRATION_REPORT.md

### 3. Update State File
```json
{
  "current_state": "REVIEW_WAVE_INTEGRATION",
  "phase": "integration",
  "review_status": "code_review_in_progress",
  "integration_review": {
    "reviewer": "agent-code-reviewer",
    "branch": "${wave_integration_branch}",
    "focus": "integration_quality",
    "started_at": "timestamp"
  }
}
```

## Exit Criteria

### Success Path → WAITING_FOR_REVIEW_WAVE_INTEGRATION
- Code Reviewer spawned successfully
- State file updated with review details
- Transition to waiting state for results

### Failure Scenarios
- **Spawn Failure** → ERROR_RECOVERY
- **Invalid Integration State** → ERROR_RECOVERY

## Key Differences from Other Reviews
- **Focus**: Integration quality, not individual effort quality
- **Scope**: Merged code interactions and conflicts
- **Timing**: After integration, before architect review
- **Purpose**: Catch integration-specific issues

## Rules Enforced
- R233: Immediate action upon state entry
- R313: Stop after spawning agent
- R238: Monitor for review reports
- R321: Any fixes require immediate backport

## Report Expected
The Code Reviewer will create:
- `REVIEW_WAVE_INTEGRATION_REPORT.md` with:
  - Integration quality assessment
  - Merge conflict resolution review
  - Cross-effort consistency check
  - Test status after integration
  - **Demo verification results (R330/R291)**
    - Demo scripts present: YES/NO
    - Demo scripts executable: YES/NO
    - Demo coverage adequate: YES/NO
    - Demo documentation complete: YES/NO
    - Demo execution results: PASSED/FAILED
  - Issues found (if any)
  - Recommendation: PASS/FAIL

### Demo Verification Checklist (R330/R291)
```markdown
## Demo Verification (R330/R291 Compliance)
- [ ] All effort demo scripts present per R330 plans
- [ ] Wave integration demo script exists
- [ ] Demo documentation (DEMO.md files) present
- [ ] Demo scripts are executable (chmod +x)
- [ ] Demo test data files available
- [ ] Demo execution logs captured in demo-results/
- [ ] Demo failures documented with root cause
- [ ] Demo coverage matches integrated features
```

## Transition Rules
- **ALWAYS** → WAITING_FOR_REVIEW_WAVE_INTEGRATION (after spawn)
- **NEVER** skip directly to REVIEW_WAVE_ARCHITECTURE
- **NEVER** proceed without Code Reviewer spawn


## 🔴🔴🔴 STATE COMPLETION CHECKLIST 🔴🔴🔴

**Execute these steps IN ORDER to properly complete REVIEW_WAVE_INTEGRATION:**

### ✅ Step 1: Complete State-Specific Work
**Refer to "Primary Actions" and "Mandatory Validations" sections above for state-specific tasks.**

Once all state work is complete, proceed to mandatory exit protocol:

---

### ✅ Step 2: Set Proposed Next State and Transition Reason
```bash
# Based on state work results, set variables for State Manager
PROPOSED_NEXT_STATE="[DETERMINE_FROM_STATE_LOGIC]"
TRANSITION_REASON="[REASON_FOR_TRANSITION]"
echo "Proposed next state: $PROPOSED_NEXT_STATE"
echo "Transition reason: $TRANSITION_REASON"
```

---

### ✅ Step 3: Spawn State Manager for State Transition (R288 - SUPREME LAW)
```bash
# State Manager handles ALL state file updates atomically
/spawn state-manager SHUTDOWN_CONSULTATION \
  --current-state "CURRENT_STATE_NAME" \
  --proposed-next-state "$PROPOSED_NEXT_STATE" \
  --transition-reason "$TRANSITION_REASON"

# State Manager will:
# 1. Validate the transition against state machine
# 2. Update all 4 state tracking locations atomically:
#    - orchestrator-state-v3.json
#    - orchestrator-state-demo.json
#    - .cascade-state-backup.json
#    - .orchestrator-state-v3.json
# 3. Commit and push all changes
# 4. Return control

echo "✅ State Manager completed transition"
```

---

### ✅ Step 4: Save TODOs (R287 - SUPREME LAW)
```bash
# Save TODO state before exit (R287 trigger)
save_todos "STATE_COMPLETE"

# Commit TODOs within 60 seconds (R287)
cd "$CLAUDE_PROJECT_DIR"
git add todos/*.todo
git commit -m "todo: orchestrator - state complete [R287]"
git push
echo "✅ TODOs saved and committed"
```

---

### ✅ Step 5: Output Continuation Flag (R405 - SUPREME LAW) ⚠️ MANDATORY
```bash
# Output continuation flag as LAST action (R405)
# Use TRUE for normal completion, FALSE only for catastrophic errors

echo "CONTINUE-SOFTWARE-FACTORY=TRUE"
```

**⚠️ THIS MUST BE THE ABSOLUTE LAST LINE OF OUTPUT BEFORE EXIT! ⚠️**

---

### ✅ Step 6: Stop Processing (R322 - SUPREME LAW)
```bash
# Stop for context preservation (R322)
echo "🛑 Stopping for context preservation - use /continue-orchestrating to resume"
exit 0
```

---

## 🚨 CHECKLIST ENFORCEMENT 🚨

**Skipping ANY step in this checklist = FAILURE:**
- Missing Step 2: No next state = stuck forever
- Missing Step 3: No State Manager spawn = state machine broken (R288 violation, -100%)
- Missing Step 4: No TODO save = work lost (R287 violation, -20% to -100%)
- Missing Step 5: No continuation flag = automation stops (R405 violation, -100%)
- Missing Step 6: No exit = R322 violation (-100%)

**ALL 6 STEPS ARE MANDATORY - NO EXCEPTIONS**
**NOTE**: State file validation and commits now handled by State Manager (SF 3.0 pattern)


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

### THE PATTERN AT R322 CHECKPOINTS (SF 3.0)

```bash
# 1. Complete state work
echo "✅ State work complete"

# 2. Set proposed next state
PROPOSED_NEXT_STATE="NEXT_STATE"
TRANSITION_REASON="State work complete"

# 3. Spawn State Manager for state transition
/spawn state-manager SHUTDOWN_CONSULTATION \
  --current-state "CURRENT_STATE" \
  --proposed-next-state "$PROPOSED_NEXT_STATE" \
  --transition-reason "$TRANSITION_REASON"
# State Manager updates all 4 state files atomically

# 4. Save TODOs
save_todos "R322_CHECKPOINT"

# 5. Factory continues (operational status)
echo "CONTINUE-SOFTWARE-FACTORY=TRUE"

# 6. Agent stops (technical requirement)
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
**See: $CLAUDE_PROJECT_DIR/rule-library/R405-automation-continuation-flag.md**

