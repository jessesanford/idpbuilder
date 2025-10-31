# PROJECT_REVIEW_WAVE_INTEGRATION State Rules


## 🔴🔴🔴 MANDATORY STATE RULE READING AND ACKNOWLEDGMENT 🔴🔴🔴

### ⚠️⚠️⚠️ YOU MUST READ EACH RULE FILE LISTED IN PRIMARY DIRECTIVES. **I AM WATCHING YOUR TOOL CALLS FOR READ OPERATIONS** *YOU WILL FAIL* IF YOU DO NOT MAKE A READ FILE CALL FOR EACH RULE FILE IN PRIMARY DIRECTIVES!!! ⚠️⚠️⚠️

**AFTER READING, YOU MUST ACKNOWLEDGE ALL THE STATE RULES AND STATE THAT YOU WILL ABIDE BY THEM ONE AT A TIME GIVING THE RULE NUMBER AND DESCRIPTION.**

## 📋 PRIMARY DIRECTIVES FOR PROJECT_REVIEW_WAVE_INTEGRATION STATE

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
Spawn Code Reviewer to perform comprehensive quality review of the fully integrated project after all phases are merged. This is the final code quality gate before project completion.

## Entry Criteria
- **From**: MONITORING_PROJECT_INTEGRATE_WAVE_EFFORTS
- **Condition**: Integration Agent has completed project integration successfully
- **Required**: 
  - Project integration workspace exists with all phases merged
  - All phase branches merged into project branch
  - Project-level merge conflicts resolved
  - No build failures from integration

## State Actions

### 1. IMMEDIATE: Spawn Code Reviewer for Project Integration Review
```bash
# Spawn Code Reviewer for comprehensive project review
/spawn agent-code-reviewer PROJECT_REVIEW_WAVE_INTEGRATION \
  --project "${project_name}" \
  --branch "${project_integration_branch}" \
  --focus "project-integration-quality" \
  --comprehensive true
```

### 2. Code Reviewer Responsibilities
The spawned Code Reviewer will:
- Perform comprehensive project-wide code review
- Validate all phases integrate correctly
- Check for cross-phase conflicts or duplications
- Ensure architectural consistency across entire project
- Verify all project requirements are met
- Validate test suite completeness and passing
- Check for performance regressions
- Review security implications
- Assess technical debt introduced
- **Verify end-to-end demo scenarios (R330/R291)**
- **Validate production readiness demonstrations**
- **Check comprehensive feature coverage in demos**
- **Review full project demo execution results**
- **Ensure demo documentation is production-ready**
- Create PROJECT_REVIEW_WAVE_INTEGRATION_REPORT.md

### 3. Update State File
```json
{
  "current_state": "PROJECT_REVIEW_WAVE_INTEGRATION",
  "phase": "project_integration",
  "review_status": "project_code_review_in_progress",
  "project_integration_review": {
    "reviewer": "agent-code-reviewer",
    "branch": "${project_integration_branch}",
    "focus": "project_integration_quality",
    "phases_integrated": ["P1", "P2", "P3"],
    "comprehensive_review": true,
    "started_at": "timestamp"
  }
}
```

## Exit Criteria

### Success Path → WAITING_FOR_PROJECT_REVIEW_WAVE_INTEGRATION
- Code Reviewer spawned successfully
- State file updated with project review details
- Transition to waiting state for comprehensive results

### Failure Scenarios
- **Spawn Failure** → ERROR_RECOVERY
- **Invalid Project State** → ERROR_RECOVERY

## Exit Conditions and Continuation Flag

**⚠️ READ THIS:** R405-CONTINUATION-FLAG-MASTER-GUIDE.md before setting flag!

### ✅ CONTINUE-SOFTWARE-FACTORY=TRUE (Standard Operation - DEFAULT)

Use TRUE when:
- ✅ Code Reviewer spawned successfully
- ✅ State file updated correctly
- ✅ Transitioning to WAITING state normally
- ✅ Ready to receive review results

**THIS IS THE NORMAL PATH - ALWAYS USE TRUE FOR PROJECT_DONEFUL SPAWN**

### 🔴 CONTINUE-SOFTWARE-FACTORY=FALSE (Exceptional - EXTREMELY RARE)

Use FALSE ONLY when:
- ❌ Cannot spawn Code Reviewer (infrastructure completely broken)
- ❌ State file corruption detected
- ❌ Project integration workspace missing/corrupt
- ❌ State machine corruption detected
- ❌ Truly unrecoverable error

**DO NOT set FALSE because:**
- ❌ "User might want to see spawn" (NO! Spawn is normal! Use TRUE!)
- ❌ R322 checkpoint (stop ≠ FALSE flag! Use TRUE!)
- ❌ Waiting for results (THIS IS NORMAL! Use TRUE!)

## CRITICAL: Review Results Are Handled in WAITING State

**DO NOT CONFUSE THIS STATE WITH THE NEXT STATE:**

**This state (PROJECT_REVIEW_WAVE_INTEGRATION):**
- Purpose: Spawn the Code Reviewer
- Action: Spawn agent and transition to WAITING
- Flag: TRUE (spawn successful = normal operation)

**Next state (WAITING_FOR_PROJECT_REVIEW_WAVE_INTEGRATION):**
- Purpose: Handle review results
- Action: Process review outcome (APPROVED, NEEDS_FIXES, etc.)
- Flag: TRUE for ALL review outcomes (system has fix protocol!)

**Review finding issues is handled in the WAITING state, not this one!**

## Grading Impact

- Setting FALSE for successful spawn: -20% (defeats automation)
- Pattern of FALSE for normal operations: -50%
- Complete automation defeat: -100%

## References

- R405-CONTINUATION-FLAG-MASTER-GUIDE.md (THE definitive guide)
- R405: Automation Flag Continuation Principle
- R322: Mandatory Stop Before State Transitions

## Key Differences from Phase Integration Review
- **Scope**: Entire project (all phases)
- **Depth**: Most comprehensive review
- **Focus**: Project-wide consistency and completeness
- **Validation**: All requirements met
- **Final Gate**: Last quality check before completion

## Rules Enforced
- R233: Immediate action upon state entry
- R313: Stop after spawning agent
- R238: Monitor for review reports
- R283: Project must include all phases
- R266: Comprehensive validation required
- R321: Any fixes require immediate backport

## Report Expected
The Code Reviewer will create:
- `PROJECT_REVIEW_WAVE_INTEGRATION_REPORT.md` with:
  - Project-wide quality assessment
  - Cross-phase integration verification
  - Requirements completion checklist
  - Architectural consistency review
  - Performance baseline comparison
  - Security audit results
  - Technical debt assessment
  - Test coverage report
  - Critical issues (if any)
  - Risk assessment
  - Recommendation: PASS/FAIL/CONDITIONAL_PASS

## Special Considerations
- This is the FINAL code review gate
- More stringent than phase reviews
- May require multiple iterations
- Focus on production readiness

## Transition Rules
- **ALWAYS** → WAITING_FOR_PROJECT_REVIEW_WAVE_INTEGRATION (after spawn)
- **NEVER** skip directly to SPAWN_CODE_REVIEWER_DEMO_VALIDATION
- **NEVER** proceed without comprehensive review


## 🔴🔴🔴 STATE COMPLETION CHECKLIST 🔴🔴🔴

**Execute these steps IN ORDER to properly complete PROJECT_REVIEW_WAVE_INTEGRATION:**

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

