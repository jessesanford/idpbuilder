# Orchestrator - DISTRIBUTE_FIX_PLANS State Rules

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


## 🔴🔴🔴 R290 ENFORCEMENT: READ THESE RULES FIRST! 🔴🔴🔴

**SUPREME LAW #3 (R290): STATE RULES MUST BE READ BEFORE STATE ACTIONS**

## 🔴🔴🔴 STOP! STATE RULE READING IS ABSOLUTELY FIRST! 🔴🔴🔴

**YOU HAVE ENTERED DISTRIBUTE_FIX_PLANS STATE - YOU MUST READ AND ACKNOWLEDGE ALL STATE RULES BEFORE DOING ANY STATE WORK!**

## 🔴🔴🔴 R290 VERIFICATION REQUIREMENT 🔴🔴🔴

**R290 ENFORCEMENT: CREATE VERIFICATION MARKER AFTER READING**

After reading and acknowledging all state rules, you MUST create a verification marker:

```bash
# MANDATORY: Create verification marker after reading rules
touch .state_rules_read_orchestrator_DISTRIBUTE_FIX_PLANS
echo "$(date +%s) - Rules read and acknowledged for DISTRIBUTE_FIX_PLANS" > .state_rules_read_orchestrator_DISTRIBUTE_FIX_PLANS
```

**FAILURE TO CREATE MARKER = AUTOMATIC -100% PENALTY**

The system will check for this marker. No marker = Immediate failure.

### ❌ DO NOT DO ANY DISTRIBUTION WORK UNTIL RULES ARE READ:
- ❌ Load fix plan summaries
- ❌ Copy plans to directories
- ❌ Create marker files
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
5. Read: $CLAUDE_PROJECT_DIR/rule-library/R239-fix-plan-distribution.md
6. Read: $CLAUDE_PROJECT_DIR/rule-library/R209-effort-directory-isolation-protocol.md
7. Read: $CLAUDE_PROJECT_DIR/rule-library/R194-remote-branch-tracking.md
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
   ❌ WRONG: "I acknowledge all DISTRIBUTE_FIX_PLANS rules"
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
   ❌ WRONG: "I know R239 requires distribution..."
   (Must READ from file, not recall from memory)
   ```

5. **Skipping Rules in PRIMARY DIRECTIVES**:
   ```
   ❌ WRONG: Reading only some rules from the list
   (ALL rules in PRIMARY DIRECTIVES are MANDATORY)
   ```

### ✅ CORRECT PATTERN FOR DISTRIBUTE_FIX_PLANS:
```
1. READ: $CLAUDE_PROJECT_DIR/rule-library/R234-mandatory-state-traversal-supreme-law.md
2. "I acknowledge R234 - Mandatory State Traversal: [Brief description]"
3. READ: $CLAUDE_PROJECT_DIR/rule-library/R006-orchestrator-never-writes-code.md  
4. "I acknowledge R006 - Orchestrator Never Writes Code: [Brief description]"
[Continue for EVERY rule in PRIMARY DIRECTIVES...]
5. Create verification marker
6. "Ready to execute DISTRIBUTE_FIX_PLANS work"
```

### 🚨 NO WORK UNTIL ACKNOWLEDGMENT COMPLETE 🚨
**You may NOT begin ANY DISTRIBUTE_FIX_PLANS work until:**
1. ✅ ALL rules in PRIMARY DIRECTIVES have been READ
2. ✅ ALL rules have been individually ACKNOWLEDGED
3. ✅ Verification marker has been CREATED
4. ✅ You have stated readiness to execute DISTRIBUTE_FIX_PLANS work

### 🚨 FAILURE TO READ STATE RULES FIRST = IMMEDIATE EXIT 🚨
**If you do ANY DISTRIBUTE_FIX_PLANS work before reading and acknowledging rules:**
- **STOP ALL WORK IMMEDIATELY**
- **EXIT WITH FAILURE STATUS**
- **YOU HAVE VIOLATED STATE COMPLIANCE**

**After reading ALL rules, acknowledge them:**
□ I have read R234 - Mandatory State Traversal (SUPREME LAW #1)
□ I have read R208 - No Implementation in Orchestrator (SUPREME LAW #2)
□ I have read R290 - State Rule Reading and Verification (SUPREME LAW #3)
□ I have read R239 - Fix Plan Distribution Protocol
□ I have read R209 - Effort Directory Isolation Protocol
□ I have read R194 - Remote Branch Tracking
□ I have read R206 - State Machine Transition Validation

**CRITICAL**: You must have made 8 actual Read tool calls. Count them!

---

## 🔴🔴🔴 SUPREME DIRECTIVE: DISTRIBUTE FIX PLANS TO EFFORT DIRECTORIES 🔴🔴🔴

**COPY FIX PLANS TO EACH EFFORT'S WORKING DIRECTORY!**

## State Overview

In DISTRIBUTE_FIX_PLANS, you copy the fix plans created by Code Reviewer to each effort's working directory so engineers can access them.

## Required Actions

### 1. Load Fix Plan Summary
```bash
PHASE=$(yq '.current_phase' orchestrator-state.yaml)
WAVE=$(yq '.current_wave' orchestrator-state.yaml)
FIX_PLAN_DIR="efforts/phase${PHASE}/wave${WAVE}/fix-plans"
SUMMARY_FILE="${FIX_PLAN_DIR}/FIX_PLAN_SUMMARY.yaml"
TIMESTAMP=$(date +%Y%m%d-%H%M%S)

echo "📦 Distributing fix plans to effort directories"

# Parse fix plans from summary
readarray -t FIX_PLANS < <(yq '.fix_plans_created[].effort' "$SUMMARY_FILE")
readarray -t PLAN_FILES < <(yq '.fix_plans_created[].plan_file' "$SUMMARY_FILE")
```

### 2. Distribute Fix Plans to Efforts
```bash
DISTRIBUTION_LOG="efforts/phase${PHASE}/wave${WAVE}/FIX_DISTRIBUTION_LOG_${TIMESTAMP}.md"
echo "# Fix Plan Distribution Log" > "$DISTRIBUTION_LOG"
echo "Timestamp: $(date -u +%Y-%m-%dT%H:%M:%SZ)" >> "$DISTRIBUTION_LOG"
echo "" >> "$DISTRIBUTION_LOG"

for i in "${!FIX_PLANS[@]}"; do
    EFFORT="${FIX_PLANS[$i]}"
    PLAN_FILE="${PLAN_FILES[$i]}"
    SOURCE_FILE="${FIX_PLAN_DIR}/${PLAN_FILE}"
    EFFORT_DIR="efforts/phase${PHASE}/wave${WAVE}/${EFFORT}"
    DEST_FILE="${EFFORT_DIR}/INTEGRATION_FIX_PLAN_${TIMESTAMP}.md"
    
    echo "📋 Distributing fix plan for effort: $EFFORT"
    
    # Verify effort directory exists
    if [ ! -d "$EFFORT_DIR" ]; then
        echo "❌ ERROR: Effort directory not found: $EFFORT_DIR"
        echo "- ERROR: $EFFORT directory not found" >> "$DISTRIBUTION_LOG"
        continue
    fi
    
    # Copy fix plan to effort directory
    cp "$SOURCE_FILE" "$DEST_FILE"
    
    # Add marker file for engineer to find
    echo "INTEGRATION_FIX_REQUIRED" > "${EFFORT_DIR}/FIX_REQUIRED.flag"
    echo "$DEST_FILE" > "${EFFORT_DIR}/FIX_PLAN_LOCATION.txt"
    
    # Commit to effort branch
    cd "$EFFORT_DIR"
    BRANCH="phase${PHASE}-wave${WAVE}-${EFFORT}"
    git checkout "$BRANCH" 2>/dev/null || git checkout -b "$BRANCH"
    git add "INTEGRATION_FIX_PLAN_${TIMESTAMP}.md" "FIX_REQUIRED.flag" "FIX_PLAN_LOCATION.txt"
    git commit -m "fix-plan: Integration fixes required [${TIMESTAMP}]" \
               -m "Fix plan distributed from Code Reviewer" \
               -m "See INTEGRATION_FIX_PLAN_${TIMESTAMP}.md for instructions"
    git push origin "$BRANCH"
    
    echo "✅ Fix plan distributed to: $DEST_FILE"
    echo "- ✅ $EFFORT: $DEST_FILE" >> "$DISTRIBUTION_LOG"
    
    # Return to orchestrator directory
    cd - > /dev/null
done

echo "" >> "$DISTRIBUTION_LOG"
echo "Total efforts with fix plans: ${#FIX_PLANS[@]}" >> "$DISTRIBUTION_LOG"
```

### 3. Update State File
```bash
# Record distribution
yq eval ".integration_feedback.wave${WAVE}.fix_plans_distributed = \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\"" -i orchestrator-state.yaml
yq eval ".integration_feedback.wave${WAVE}.distribution_log = \"$DISTRIBUTION_LOG\"" -i orchestrator-state.yaml

# Mark efforts as needing fixes
for effort in "${FIX_PLANS[@]}"; do
    yq eval ".efforts_in_progress.\"$effort\".needs_fixes = true" -i orchestrator-state.yaml
    yq eval ".efforts_in_progress.\"$effort\".fix_plan_distributed = \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\"" -i orchestrator-state.yaml
done

# Transition to spawn engineers
yq eval ".current_state = \"SPAWN_ENGINEERS_FOR_FIXES\"" -i orchestrator-state.yaml
yq eval ".state_transition_history += [{\"from\": \"DISTRIBUTE_FIX_PLANS\", \"to\": \"SPAWN_ENGINEERS_FOR_FIXES\", \"timestamp\": \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\", \"reason\": \"Fix plans distributed to ${#FIX_PLANS[@]} efforts\"}]" -i orchestrator-state.yaml

# Commit
git add orchestrator-state.yaml "$DISTRIBUTION_LOG"
git commit -m "distribute: Fix plans sent to ${#FIX_PLANS[@]} efforts"
git push
```

## Valid Transitions

1. **ALWAYS**: `DISTRIBUTE_FIX_PLANS` → `SPAWN_ENGINEERS_FOR_FIXES`
   - Transition after distribution complete

## Distribution Requirements

For each effort needing fixes:
1. Copy fix plan to effort directory
2. Create FIX_REQUIRED.flag marker
3. Create FIX_PLAN_LOCATION.txt pointer
4. Commit to effort's branch
5. Push to remote

## Grading Criteria

- ✅ **+25%**: Copy plans to correct directories
- ✅ **+25%**: Create marker files
- ✅ **+25%**: Commit and push to effort branches
- ✅ **+25%**: Update state file properly

## Common Violations

- ❌ **-100%**: Not distributing fix plans
- ❌ **-50%**: Wrong effort directories
- ❌ **-50%**: Not committing to branches
- ❌ **-30%**: Missing marker files

## Related Rules

- R239: Fix Plan Distribution Protocol
- R209: Effort Directory Isolation Protocol
- R194: Remote Branch Tracking
- R206: State Machine Transition Validation

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

## R322 VIOLATION DETECTION

If you find yourself:
- Starting work for a new state without /continue-orchestrating
- Transitioning without stopping after state file commit
- Continuing after completing state work

**STOP IMMEDIATELY - You are violating R322!**
