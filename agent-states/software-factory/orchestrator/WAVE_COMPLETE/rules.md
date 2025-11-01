# Orchestrator - WAVE_COMPLETE State Rules

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
## 🟢 NORMAL OPERATION - STOP WITH CONTINUATION FLAG TRUE 🟢

**CRITICAL: You MUST stop after this state, but set flag to TRUE for normal operations**

### MANDATORY STOP PROTOCOL (R322):
1. Complete all wave completion tasks
2. Update orchestrator-state-v3.json to next state
3. Commit and push the state file
4. Output: CONTINUE-SOFTWARE-FACTORY=TRUE (for normal operations)
5. **STOP PROCESSING** (mandatory for context preservation)

### KEY UNDERSTANDING:
- ✅ Stopping is MANDATORY (prevents context overflow)
- ✅ TRUE flag means "all is well, restart me automatically"
- ✅ External automation sees TRUE and restarts Claude Code
- ✅ The stop preserves context while allowing continuous operation

### 🛑 R322 CHECKPOINT CLARIFICATION - CRITICAL UNDERSTANDING

**AGENT STOPS ≠ FACTORY STOPS** (These are TWO INDEPENDENT operations!)

At R322 checkpoints in WAVE_COMPLETE:
- ✅ **Agent STOPS work** = Agent completes state tasks, exits process, preserves context
- ✅ **Factory CONTINUES** = Set CONTINUE-SOFTWARE-FACTORY=TRUE for automation
- ✅ **User resumes** = Run /continue-orchestrating to start next state
- ✅ **This is NORMAL** = Expected workflow, not an error condition

**KEY INSIGHT:**
- "MUST STOP" in R322 = Agent exits to save context (technical requirement)
- "CONTINUE-SOFTWARE-FACTORY=TRUE" = Factory should proceed (normal operation)
- **BOTH happen together at successful checkpoints!**

**THE CONFUSION TO AVOID:**
- ❌ WRONG: "I must stop (R322) → Manual intervention needed → FALSE"
- ✅ RIGHT: "I must stop (R322) → Normal checkpoint → TRUE → User continues me later"

**R322 checkpoints are NORMAL, not exceptional!**
- ✅ WAVE_COMPLETE requires a stop for context preservation (R322)
- ✅ Moving to INTEGRATE_WAVE_EFFORTS is NORMAL (use TRUE flag per R405)
- ✅ TRUE flag enables automatic continuation after restart
- ✅ Stop with TRUE = normal operation that continues automatically
- ✅ Stop with FALSE = catastrophic corruption preventing any progress

### ❌ COMMON ORCHESTRATOR MISTAKES TO AVOID:
- ❌ Continuing without stopping (causes context overflow)
- ❌ Setting flag to FALSE for normal operations (blocks automation)
- ❌ Not stopping at all (violates R322)
- ❌ Confusing "stop" with "failure" (stop is normal!)

### NORMAL COMPLETION PROTOCOL:
```bash
# Complete wave validation
echo "✅ Wave complete - all efforts finished and reviewed"

# Update to next state (INTEGRATE_WAVE_EFFORTS is default)
NEXT_STATE="INTEGRATE_WAVE_EFFORTS"

# Commit state transition
git add orchestrator-state-v3.json
git commit -m "state: transition from WAVE_COMPLETE to $NEXT_STATE"
git push

# Stop with TRUE flag for automatic restart
echo "✅ Stopping for context preservation (normal operation)"
echo "CONTINUE-SOFTWARE-FACTORY=TRUE REASON=STATE_COMPLETE"  # TRUE = auto-restart allowed
exit 0  # STOP processing (mandatory)
```

**STOP PROCESSING but set TRUE flag to signal automatic restart by external automation!**

---


## 🔴🔴🔴 STOP! STATE RULE READING IS ABSOLUTELY FIRST! 🔴🔴🔴

**YOU HAVE ENTERED WAVE_COMPLETE STATE - YOU MUST READ AND ACKNOWLEDGE ALL STATE RULES BEFORE DOING ANY STATE WORK!**

## 🔴🔴🔴 R290 VERIFICATION REQUIREMENT 🔴🔴🔴

**R290 ENFORCEMENT: CREATE VERIFICATION MARKER AFTER READING**

After reading and acknowledging all state rules, you MUST create a verification marker:

```bash
# MANDATORY: Create verification marker after reading rules
touch .state_rules_read_orchestrator_WAVE_COMPLETE
echo "$(date +%s) - Rules read and acknowledged for WAVE_COMPLETE" > .state_rules_read_orchestrator_WAVE_COMPLETE
```

**FAILURE TO CREATE MARKER = AUTOMATIC -100% PENALTY**

The system will check for this marker. No marker = Immediate failure.

### ❌ DO NOT DO ANY WAVE_COMPLETE WORK UNTIL RULES ARE READ:
- ❌ Start finalize wave efforts
- ❌ Start collect implementation results
- ❌ Start prepare integration
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
   ❌ WRONG: "I acknowledge all WAVE_COMPLETE rules"
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

### ✅ CORRECT PATTERN FOR WAVE_COMPLETE:
```
1. READ: $CLAUDE_PROJECT_DIR/rule-library/[first-rule-file].md
2. "I acknowledge [Rule#] - [Rule Name]: [Brief description]"
3. READ: $CLAUDE_PROJECT_DIR/rule-library/[second-rule-file].md  
4. "I acknowledge [Rule#] - [Rule Name]: [Brief description]"
[Continue for EVERY rule in PRIMARY DIRECTIVES...]
5. "Ready to execute WAVE_COMPLETE work"
```

### 🚨 NO WORK UNTIL ACKNOWLEDGMENT COMPLETE 🚨
**You may NOT begin ANY WAVE_COMPLETE work until:**
1. ✅ ALL rules in PRIMARY DIRECTIVES have been READ
2. ✅ ALL rules have been individually ACKNOWLEDGED
3. ✅ You have stated readiness to execute WAVE_COMPLETE work
1. **READ** every rule file listed in PRIMARY DIRECTIVES below
2. **ACKNOWLEDGE** each rule individually with number and description
3. **ONLY THEN** proceed with WAVE_COMPLETE work

### 🚨 FAILURE TO READ STATE RULES FIRST = IMMEDIATE EXIT 🚨
**If you do ANY WAVE_COMPLETE work before reading and acknowledging rules:**
- **STOP ALL WORK IMMEDIATELY**
- **EXIT WITH FAILURE STATUS**
- **YOU HAVE VIOLATED STATE COMPLIANCE**

**THE SYSTEM IS MONITORING_SWE_PROGRESS YOUR READ TOOL CALLS!**

### ✅ Step 2: Set Proposed Next State and Transition Reason
```bash
# Based on state work results, determine proposed next state
PROPOSED_NEXT_STATE="[DETERMINE_FROM_STATE_LOGIC]"
TRANSITION_REASON="WAVE_COMPLETE complete - [accomplishment description]"
echo "Proposed next state: $PROPOSED_NEXT_STATE"
echo "Transition reason: $TRANSITION_REASON"
```

---

### ✅ Step 3: Spawn State Manager for SHUTDOWN_CONSULTATION
```bash
# State Manager validates transition and updates state files (SF 3.0 Pattern)
echo "🔄 Spawning State Manager for SHUTDOWN_CONSULTATION..."

# Prepare work results summary
WORK_RESULTS=$(cat <<EOF
{
  "state_completed": "WAVE_COMPLETE",
  "work_accomplished": [
    "Verified all efforts complete and passed reviews",
    "Updated orchestrator-state-v3.json with completion",
    "Checked for stale integrations requiring CASCADE",
    "Determined next state (INTEGRATE_WAVE_EFFORTS or CASCADE_REINTEGRATION)"
  ],
  "proposed_next_state": "$PROPOSED_NEXT_STATE",
  "transition_reason": "$TRANSITION_REASON"
}
EOF
)

# Spawn State Manager
/spawn state-manager SHUTDOWN_CONSULTATION \
  --current-state "WAVE_COMPLETE" \
  --proposed-next-state "$PROPOSED_NEXT_STATE" \
  --transition-reason "$TRANSITION_REASON" \
  --work-results "$WORK_RESULTS"

# State Manager will:
# 1. Validate PROPOSED_NEXT_STATE exists and transition is valid
# 2. Update all 4 state files atomically (R288)
# 3. Commit and push state files
# 4. Return REQUIRED_NEXT_STATE (usually same as proposed unless invalid)

echo "✅ State Manager consultation complete"
echo "✅ State files updated by State Manager"
```

---

### ✅ Step 4: Save TODOs (R287 - SUPREME LAW)
```bash
# Save TODO state before transition (R287 trigger)
save_todos "WAVE_COMPLETE_COMPLETE"

# Commit TODOs within 60 seconds (R287)
cd "$CLAUDE_PROJECT_DIR"
git add todos/*.todo

if ! git commit -m "todo: orchestrator - WAVE_COMPLETE complete [R287]"; then
    echo "❌ ERROR: Failed to commit TODO files"
    echo "This is non-fatal but TODOs may be lost in compaction"
    echo "Proceeding with state execution..."
    # Don't exit - TODO commit failure is not fatal
fi

git push || echo "⚠️ WARNING: TODO push failed - committed locally"
echo "✅ TODOs saved and committed"
git push
echo "✅ TODOs saved and committed"
```

---

### ✅ Step 5: Output Continuation Flag (R405 - SUPREME LAW) ⚠️ MANDATORY
```bash
# Output continuation flag as LAST action (R405)
# Use TRUE for normal completion, FALSE only for catastrophic errors

echo "CONTINUE-SOFTWARE-FACTORY=TRUE REASON=STATE_COMPLETE"
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
- Missing Step 2: No proposed next state = State Manager can't proceed
- Missing Step 3: No State Manager consultation = bypassing bookend pattern (-100%)
- Missing Step 4: No TODO save = work lost (R287 violation, -20% to -100%)
- Missing Step 5: No continuation flag = automation stops (R405 violation, -100%)
- Missing Step 6: No exit = R322 violation (-100%)

**ALL 6 STEPS ARE MANDATORY - NO EXCEPTIONS**
**NOTE**: State file validation and commits now handled by State Manager (SF 3.0 pattern)

## ✅ EXIT CHECKLIST (R405 - Continuation Flag Protocol) 🔴🔴🔴

**SUPREME LAW - PENALTY FOR VIOLATION: -100% GRADE**

### 🟢🟢🟢 WAVE_COMPLETE ALWAYS SETS TRUE FOR AUTOMATION! 🟢🟢🟢

**WAVE_COMPLETE IS NORMAL OPERATION - ALWAYS CONTINUE AUTOMATICALLY!**

### YOU MUST OUTPUT THE CONTINUATION FLAG AS YOUR LAST ACTION

**FOR WAVE_COMPLETE STATE - ALWAYS OUTPUT TRUE!**
```bash
# WAVE_COMPLETE → INTEGRATE_WAVE_EFFORTS is NORMAL FLOW per R322
# This is NOT an exceptional checkpoint!
# Always enable automation to continue:
echo "CONTINUE-SOFTWARE-FACTORY=TRUE REASON=STATE_COMPLETE"  # ALWAYS TRUE FOR WAVE_COMPLETE!
```

### ⚠️⚠️⚠️ WARNING: COMMON MISTAKE ⚠️⚠️⚠️
**DO NOT STOP FOR "MILESTONES" OR "HUMAN REVIEW"!**
- ❌ WRONG: "Wave complete, stopping for user review" → FALSE
- ✅ RIGHT: "Wave complete, continuing to integration" → TRUE

### CRITICAL REQUIREMENTS:
1. **ABSOLUTE LAST OUTPUT**: This MUST be the very last line of output before state completion
2. **NO TEXT AFTER**: No explanations, summaries, or any text after the flag
3. **EXACTLY AS SHOWN**: Use exact format - no variations like CONTINUE-ORCHESTRATING
4. **ALWAYS REQUIRED**: Every single state must output this flag
5. **GREPPABLE**: Must be on its own line for automation parsing

### ✅ WHEN TO USE TRUE IN WAVE_COMPLETE (This state = 99.9% TRUE):
- ✅ Wave validation complete → TRUE (NORMAL R322 checkpoint)
- ✅ All efforts reviewed and passed → TRUE (PROJECT_DONE!)
- ✅ State updated to INTEGRATE_WAVE_EFFORTS → TRUE (NORMAL workflow)
- ✅ Ready to create integration branch → TRUE (EXPECTED)
- ✅ No blockers detected → TRUE (ALL GOOD)
- ✅ Minor issues exist → TRUE (handle in next state)
- ✅ **Agent stopping per R322** → TRUE (NORMAL CHECKPOINT!)
- ✅ **Waiting for /continue-orchestrating** → TRUE (EXPECTED!)

### ❌ WHEN TO USE FALSE IN WAVE_COMPLETE (Almost never!):
**ONLY USE FALSE FOR CATASTROPHIC SYSTEM FAILURES:**
- ❌ State file corrupted beyond repair (cannot record completion)
- ❌ Git repository completely broken (cannot commit)
- ❌ Data corruption that would cascade to other efforts
- ❌ **NOT for R322 stops (agent stopping is NORMAL → use TRUE)**
- ❌ **NOT for "milestone reached" (use TRUE)**
- ❌ **NOT for "user might want to review" (use TRUE)**
- ❌ **NOT for "being cautious" (use TRUE)**
- ❌ **NOT for "important transition" (use TRUE)**

### 🔴 CRITICAL: UNDERSTAND THE TWO INDEPENDENT DECISIONS

**Decision 1: Should Agent Stop? (R322)**
- Answer for WAVE_COMPLETE: **YES, ALWAYS** (context preservation)

**Decision 2: Should Factory Continue? (R405)**
- Answer for WAVE_COMPLETE: **TRUE, ALMOST ALWAYS** (normal operation)

**These are DIFFERENT questions with DIFFERENT answers!**
- Agent stops = Technical requirement (save memory/context)
- Factory continues = Operational status (is system healthy?)

### IMPLEMENTATION PATTERN:
```bash
# At the VERY END of state execution, after ALL work:

# Determine success/failure
if [[ "$STATE_COMPLETED_PROJECT_DONEFULLY" == "true" ]]; then
    echo "✅ State work completed successfully"
    echo "CONTINUE-SOFTWARE-FACTORY=TRUE REASON=STATE_COMPLETE"
else
    echo "❌ State encountered issues requiring intervention"
    echo "CONTINUE-SOFTWARE-FACTORY=FALSE REASON=MANUAL_INTERVENTION_REQUIRED"
fi

# DO NOT OUTPUT ANYTHING AFTER THE FLAG!
```

### GRADING IMPACT:
- **Missing flag**: -100% AUTOMATIC FAILURE
- **Text after flag**: -50% penalty
- **Wrong format**: -100% AUTOMATIC FAILURE
- **Multiple flags**: -50% penalty

**See: $CLAUDE_PROJECT_DIR/rule-library/R405-automation-continuation-flag.md**

