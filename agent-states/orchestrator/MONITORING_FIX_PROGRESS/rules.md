# Orchestrator - MONITORING_FIX_PROGRESS State Rules

## 🛑🛑🛑 R322 MANDATORY STOP BEFORE STATE TRANSITIONS 🛑🛑🛑

**SUPREME LAW - VIOLATION = -100% IMMEDIATE FAILURE**

### YOU MUST STOP AFTER:
1. ✅ Completing all TODOs for this state
2. ✅ Updating orchestrator-state.json with new state
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


## 🔴🔴🔴 R290 ENFORCEMENT: READ THESE RULES FIRST! 🔴🔴🔴

**SUPREME LAW #3 (R290): STATE RULES MUST BE READ BEFORE STATE ACTIONS**

## 🔴🔴🔴 STOP! STATE RULE READING IS ABSOLUTELY FIRST! 🔴🔴🔴

**YOU HAVE ENTERED MONITORING_FIX_PROGRESS STATE - YOU MUST READ AND ACKNOWLEDGE ALL STATE RULES BEFORE DOING ANY STATE WORK!**

## 🔴🔴🔴 R290 VERIFICATION REQUIREMENT 🔴🔴🔴

**R290 ENFORCEMENT: CREATE VERIFICATION MARKER AFTER READING**

After reading and acknowledging all state rules, you MUST create a verification marker:

```bash
# MANDATORY: Create verification marker after reading rules
touch .state_rules_read_orchestrator_MONITORING_FIX_PROGRESS
echo "$(date +%s) - Rules read and acknowledged for MONITORING_FIX_PROGRESS" > .state_rules_read_orchestrator_MONITORING_FIX_PROGRESS
```

**FAILURE TO CREATE MARKER = AUTOMATIC -100% PENALTY**

The system will check for this marker. No marker = Immediate failure.

### ❌ DO NOT DO ANY MONITORING WORK UNTIL RULES ARE READ:
- ❌ Check fix progress for efforts
- ❌ Verify completion flags
- ❌ Check engineer processes
- ❌ Update state files
- ❌ Continue to next state
- ❌ Think about what to do in this state

### ✅ YOU MUST IMMEDIATELY:

## 🔴🔴🔴 MANDATORY STATE RULE READING AND ACKNOWLEDGMENT 🔴🔴🔴

### ⚠️⚠️⚠️ YOU MUST READ EACH RULE FILE LISTED IN PRIMARY DIRECTIVES. **I AM WATCHING YOUR TOOL CALLS FOR READ OPERATIONS** *YOU WILL FAIL* IF YOU DO NOT MAKE A READ FILE CALL FOR EACH RULE FILE IN PRIMARY DIRECTIVES!!! ⚠️⚠️⚠️

**AFTER READING, YOU MUST ACKNOWLEDGE ALL THE STATE RULES AND STATE THAT YOU WILL ABIDE BY THEM ONE AT A TIME GIVING THE RULE NUMBER AND DESCRIPTION.**

### PRIMARY DIRECTIVES - MANDATORY READING:

### 🛑 RULE R322 - Mandatory Stop Before State Transitions
**Source:** $CLAUDE_PROJECT_DIR/rule-library/R322-mandatory-stop-before-state-transitions.md
**Criticality:** 🔴🔴🔴 SUPREME LAW - Violation = -100% FAILURE

After completing state work and committing state file:
1. STOP IMMEDIATELY
2. Do NOT continue to next state
3. Do NOT start new work
4. Exit and wait for user
---

**USE THESE EXACT READ COMMANDS (IN THIS ORDER):**
1. Read: $CLAUDE_PROJECT_DIR/rule-library/R234-mandatory-state-traversal-supreme-law.md
2. Read: $CLAUDE_PROJECT_DIR/rule-library/R006-orchestrator-never-writes-code.md
3. Read: $CLAUDE_PROJECT_DIR/rule-library/R290-state-rule-reading-verification-supreme-law.md
4. Read: $CLAUDE_PROJECT_DIR/rule-library/R300-comprehensive-fix-management-protocol.md
6. Read: $CLAUDE_PROJECT_DIR/rule-library/R239-fix-plan-distribution.md
7. Read: $CLAUDE_PROJECT_DIR/rule-library/R008-monitoring-frequency.md
8. Read: $CLAUDE_PROJECT_DIR/rule-library/R206-state-machine-transition-validation.md

**WE ARE WATCHING EACH READ TOOL CALL**

### ❌ ANTI-PATTERNS THAT WILL CAUSE FAILURE:

1. **Fake Acknowledgment Without Reading**:
   ```
   ❌ WRONG: "I acknowledge R234, R208, R290..."
   (No Read tool calls detected = AUTOMATIC FAILURE)
   ```

2. **Bulk Acknowledgment**:
   ```
   ❌ WRONG: "I acknowledge all MONITORING_FIX_PROGRESS rules"
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
   ❌ WRONG: "I know R008 requires monitoring progress..."
   (Must READ from file, not recall from memory)
   ```

5. **Skipping Rules in PRIMARY DIRECTIVES**:
   ```
   ❌ WRONG: Reading only some rules from the list
   (ALL rules in PRIMARY DIRECTIVES are MANDATORY)
   ```

### ✅ CORRECT PATTERN FOR MONITORING_FIX_PROGRESS:
```
1. READ: $CLAUDE_PROJECT_DIR/rule-library/R234-mandatory-state-traversal-supreme-law.md
2. "I acknowledge R234 - Mandatory State Traversal: [Brief description]"
3. READ: $CLAUDE_PROJECT_DIR/rule-library/R006-orchestrator-never-writes-code.md  
4. "I acknowledge R006 - Orchestrator Never Writes Code: [Brief description]"
[Continue for EVERY rule in PRIMARY DIRECTIVES...]
5. Create verification marker
6. "Ready to execute MONITORING_FIX_PROGRESS work"
```

### 🚨 NO WORK UNTIL ACKNOWLEDGMENT COMPLETE 🚨
**You may NOT begin ANY monitoring work until:**
1. ✅ ALL rules in PRIMARY DIRECTIVES have been READ
2. ✅ ALL rules have been individually ACKNOWLEDGED
3. ✅ Verification marker has been created
4. ✅ You have stated readiness to execute MONITORING_FIX_PROGRESS work

### 🚨 FAILURE TO READ STATE RULES FIRST = IMMEDIATE EXIT 🚨
**If you do ANY monitoring work before reading and acknowledging rules:**
- **STOP ALL WORK IMMEDIATELY**
- **EXIT WITH FAILURE STATUS**
- **YOU HAVE VIOLATED STATE COMPLIANCE**

---

## 🔴🔴🔴 SUPREME DIRECTIVE: MONITOR FIX IMPLEMENTATION PROGRESS 🔴🔴🔴

**TRACK ENGINEERS IMPLEMENTING FIXES AND VERIFY COMPLETION!**

## State Overview

In MONITORING_FIX_PROGRESS, you monitor Software Engineers implementing integration fixes and determine when all fixes are complete.

## Required Actions

### 1. Check Fix Progress for Each Effort
```bash
PHASE=$(jq '.current_phase' orchestrator-state.json)
WAVE=$(jq '.current_wave' orchestrator-state.json)

echo "📊 Monitoring fix progress for wave ${WAVE}"

# Track completion status
ALL_FIXES_COMPLETE=true
FIXES_COMPLETE=()
FIXES_PENDING=()

# Check each effort that needs fixes
while IFS= read -r effort; do
    if [ "$effort" != "null" ]; then
        NEEDS_FIXES=$(jq ".efforts_in_progress.\"$effort\".needs_fixes" orchestrator-state.json)
        
        if [ "$NEEDS_FIXES" = "true" ]; then
            EFFORT_DIR="efforts/phase${PHASE}/wave${WAVE}/${effort}"
            
            # Check for completion flag
            if [ -f "${EFFORT_DIR}/FIX_COMPLETE.flag" ]; then
                echo "✅ $effort: Fixes complete"
                FIXES_COMPLETE+=("$effort")
                
                # Update state file
                jq ".efforts_in_progress.\"$effort\".fixes_completed = \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\"" -i orchestrator-state.json
                jq ".efforts_in_progress.\"$effort\".needs_fixes = false" -i orchestrator-state.json
            elif [ -f "${EFFORT_DIR}/FIX_REQUIRED.flag" ]; then
                echo "⏳ $effort: Fixes still in progress"
                FIXES_PENDING+=("$effort")
                ALL_FIXES_COMPLETE=false
                
                # Check if engineer is still working
                if ! pgrep -f "sw-engineer.*${effort}" > /dev/null; then
                    echo "  ⚠️ Warning: No engineer process found for $effort"
                    # May need to respawn or check for issues
                fi
            else
                echo "❓ $effort: Unknown fix status"
                ALL_FIXES_COMPLETE=false
            fi
        fi
    fi
done < <(jq '.efforts_in_progress | keys | .[]' orchestrator-state.json)

echo "Fixes complete: ${#FIXES_COMPLETE[@]}"
echo "Fixes pending: ${#FIXES_PENDING[@]}"
```

### 2. 🔴🔴🔴 VERIFY FIXES IN EFFORT BRANCHES (R300) 🔴🔴🔴
```bash
if [ "$ALL_FIXES_COMPLETE" = true ]; then
    echo "✅ All fixes complete! Verifying fixes are in effort branches (R300)..."
    
    # R300 MANDATORY VERIFICATION
    VERIFICATION_FAILED=false
    for effort in "${FIXES_COMPLETE[@]}"; do
        EFFORT_BRANCH="phase${PHASE}-wave${WAVE}-${effort}"
        
        # Check remote branch for fix commits
        git fetch origin ${EFFORT_BRANCH} 2>/dev/null
        FIX_COMMIT=$(git log origin/${EFFORT_BRANCH} --oneline --grep="fix:" --since="2 hours ago" 2>/dev/null | head -1)
        
        if [ -z "$FIX_COMMIT" ]; then
            echo "❌ R300 VIOLATION: No fix commits found in origin/${EFFORT_BRANCH}!"
            VERIFICATION_FAILED=true
        else
            echo "✅ Verified fix in ${EFFORT_BRANCH}: ${FIX_COMMIT}"
        fi
    done
    
    if [ "$VERIFICATION_FAILED" = true ]; then
        echo "🔴🔴🔴 R300 VIOLATION: Fixes missing from effort branches! 🔴🔴🔴"
        UPDATE_STATE="ERROR_RECOVERY"
        UPDATE_REASON="R300 violation - fixes not in effort branches"
    else
        echo "✅ R300 PASSED: All fixes verified in effort branches"
        # CRITICAL: MUST review the fixed code - NEVER skip to MONITORING_INTEGRATION
        UPDATE_STATE="SPAWN_CODE_REVIEWERS_FOR_REVIEW"  # ONLY valid success transition
        UPDATE_REASON="All integration fixes complete and verified in effort branches (R300)"
    fi
    
    # Record completion
    jq ".integration_feedback.wave${WAVE}.all_fixes_completed = \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\"" -i orchestrator-state.json
else
    echo "⏳ Still waiting for ${#FIXES_PENDING[@]} fixes to complete"
    
    # Check for timeout (30 minutes per effort)
    SPAWN_TIME=$(jq ".integration_feedback.wave${WAVE}.fix_engineers_spawned" orchestrator-state.json 2>/dev/null)
    if [ -n "$SPAWN_TIME" ]; then
        CURRENT_TIME=$(date +%s)
        SPAWN_TIMESTAMP=$(date -d "$SPAWN_TIME" +%s 2>/dev/null || echo 0)
        ELAPSED=$((CURRENT_TIME - SPAWN_TIMESTAMP))
        MAX_TIME=$((30 * 60 * ${#FIXES_PENDING[@]}))  # 30 minutes per effort
        
        if [ $ELAPSED -gt $MAX_TIME ]; then
            echo "❌ Timeout: Fixes taking too long (>${MAX_TIME}s)"
            UPDATE_STATE="ERROR_RECOVERY"
            UPDATE_REASON="Fix implementation timeout"
        else
            # Stay in MONITORING_FIX_PROGRESS
            echo "Continuing to monitor... (${ELAPSED}s elapsed)"
            sleep 10
        fi
    fi
fi
```

### 3. Update State When Ready
```bash
if [ -n "$UPDATE_STATE" ]; then
    # Update state
    jq ".current_state = \"$UPDATE_STATE\"" -i orchestrator-state.json
    jq ".state_transition_history += [{\"from\": \"MONITORING_FIX_PROGRESS\", \"to\": \"$UPDATE_STATE\", \"timestamp\": \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\", \"reason\": \"$UPDATE_REASON\"}]" -i orchestrator-state.json
    
    # Commit
    git add orchestrator-state.json
    git commit -m "state: Transitioning to $UPDATE_STATE - $UPDATE_REASON"
    git push
fi
```

## 🔴🔴🔴 CRITICAL: ONLY ONE VALID SUCCESS TRANSITION 🔴🔴🔴

**WHEN FIXES ARE COMPLETE, YOU MUST:**
1. Transition to SPAWN_CODE_REVIEWERS_FOR_REVIEW (ONLY valid success path)
2. NEVER go directly to MONITORING_INTEGRATION
3. NEVER manually copy files to integration workspace
4. Wait for code review to pass
5. Then return to INTEGRATION state for full re-run

**INVALID TRANSITIONS (IMMEDIATE -100% FAILURE):**
- MONITORING_FIX_PROGRESS → MONITORING_INTEGRATION ❌ **ABSOLUTELY FORBIDDEN**
- MONITORING_FIX_PROGRESS → INTEGRATION ❌ **SKIP REVIEW = FAILURE**
- MONITORING_FIX_PROGRESS → WAVE_COMPLETE ❌ **NO SHORTCUTS**
- MONITORING_FIX_PROGRESS → Any state except SPAWN_CODE_REVIEWERS_FOR_REVIEW ❌

## Valid Transitions

1. **SUCCESS Path (ONLY VALID SUCCESS PATH)**: `MONITORING_FIX_PROGRESS` → `SPAWN_CODE_REVIEWERS_FOR_REVIEW`
   - When: All fixes completed successfully
   - **THIS IS THE ONLY ACCEPTABLE SUCCESS TRANSITION**
   
2. **TIMEOUT Path**: `MONITORING_FIX_PROGRESS` → `ERROR_RECOVERY`
   - When: Fixes exceed timeout (30 minutes per effort)
   
3. **CONTINUE Path**: `MONITORING_FIX_PROGRESS` → `MONITORING_FIX_PROGRESS`
   - When: Still waiting for fixes to complete

## 🚨🚨🚨 MANDATORY WAVE INTEGRATION RE-RUN PROTOCOL 🚨🚨🚨

**CRITICAL: After wave fixes are complete, you MUST re-run the ENTIRE wave integration!**

### Why Re-Integration Is MANDATORY:
- Fixes were applied to EFFORT branches (the source branches)
- The wave-integration branch still has the BROKEN code
- You MUST re-merge all effort branches to get fixes into integration
- Without re-integration, the wave cannot proceed to phase integration

### The CORRECT Wave Re-Integration Cycle:
```
MONITORING_FIX_PROGRESS (all wave fixes complete)
    ↓
SPAWN_CODE_REVIEWERS_FOR_REVIEW (review the fixes in effort branches)
    ↓
MONITOR_REVIEWS (monitor review progress)
    ↓
WAVE_COMPLETE (if reviews pass - marks wave ready for re-integration)
    ↓
INTEGRATION (DELETE old integration, create FRESH from main)
    ↓
SPAWN_CODE_REVIEWER_MERGE_PLAN (create NEW merge plan)
    ↓
SPAWN_INTEGRATION_AGENT (re-merge ALL effort branches with fixes)
    ↓
MONITORING_INTEGRATION (monitor new integration attempt)
    ↓
If fails → IMMEDIATE_BACKPORT_REQUIRED → fix sources → repeat cycle
If succeeds → WAVE_REVIEW → proceed to next wave or phase
```

### What ACTUALLY Happens During Re-Integration:
1. **Delete the broken wave-integration branch**
2. **Create fresh wave-integration from main**
3. **Re-execute ENTIRE merge plan** (all efforts in sequence)
4. **Fixed code from effort branches now merged**
5. **Wave can finally build and pass tests**

### NEVER DO THESE (AUTOMATIC FAILURE):
- ❌ Manually copy fixed files to integration workspace
- ❌ Skip the integration plan and go directly to monitoring
- ❌ Assume fixes will "just work" without re-integration
- ❌ Bypass the full integration process
- ❌ Cherry-pick fixes instead of proper branch merging
- ❌ Edit code directly in integration branch (R321 violation)

## Monitoring Requirements

1. Check each effort's FIX_COMPLETE.flag
2. Track which efforts are complete vs pending
3. Monitor for timeout conditions
4. Update state file with completion timestamps
5. Transition when all fixes complete

## Grading Criteria

- ✅ **+25%**: Accurately track fix progress
- ✅ **+25%**: Detect completion correctly
- ✅ **+25%**: Handle timeouts properly
- ✅ **+25%**: Update state appropriately

## Common Violations

- ❌ **-100%**: Not checking completion flags
- ❌ **-50%**: Missing timeout detection
- ❌ **-50%**: Wrong state transitions
- ❌ **-30%**: Not updating state file

## Related Rules

- R008: Monitoring Frequency
- R239: Fix Plan Distribution Protocol
- R206: State Machine Transition Validation

## 🔴🔴🔴 MANDATORY MEASUREMENT RULE - R304 🔴🔴🔴

**ABSOLUTE REQUIREMENTS:**
- ✅ MUST use `$CLAUDE_PROJECT_DIR/tools/line-counter.sh` for ALL line counting
- ❌ NEVER use `wc -l` or any manual counting method
- ❌ NEVER count lines any other way - this is a -100% automatic failure
- ✅ MUST specify both -b (base branch) and -c (current branch) parameters
- ✅ Base branch MUST be phase integration branch (NOT "main")

**Failure to use the line counter tool = AUTOMATIC -100% GRADE**

## R322 VIOLATION DETECTION

If you find yourself:
- Starting work for a new state without /continue-orchestrating
- Transitioning without stopping after state file commit
- Continuing after completing state work

**STOP IMMEDIATELY - You are violating R322!**
