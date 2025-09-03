# Orchestrator - WAITING_FOR_EFFORT_PLANS State Rules

## 🔴🔴🔴 STOP! STATE RULE READING IS ABSOLUTELY FIRST! 🔴🔴🔴

**YOU HAVE ENTERED WAITING_FOR_EFFORT_PLANS STATE - YOU MUST READ AND ACKNOWLEDGE ALL STATE RULES BEFORE DOING ANY STATE WORK!**

## 🔴🔴🔴 R290 VERIFICATION REQUIREMENT 🔴🔴🔴

**R290 ENFORCEMENT: CREATE VERIFICATION MARKER AFTER READING**

After reading and acknowledging all state rules, you MUST create a verification marker:

```bash
# MANDATORY: Create verification marker after reading rules
touch .state_rules_read_orchestrator_WAITING_FOR_EFFORT_PLANS
echo "$(date +%s) - Rules read and acknowledged for WAITING_FOR_EFFORT_PLANS" > .state_rules_read_orchestrator_WAITING_FOR_EFFORT_PLANS
```

**FAILURE TO CREATE MARKER = AUTOMATIC -100% PENALTY**

The system will check for this marker. No marker = Immediate failure.

### ❌ DO NOT DO ANY WAITING_FOR_EFFORT_PLANS WORK UNTIL RULES ARE READ:
- ❌ Start check effort plan status
- ❌ Start monitor reviewer progress
- ❌ Start collect completed plans
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
   ❌ WRONG: "I acknowledge all WAITING_FOR_EFFORT_PLANS rules"
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

### ✅ CORRECT PATTERN FOR WAITING_FOR_EFFORT_PLANS:
```
1. READ: $CLAUDE_PROJECT_DIR/rule-library/[first-rule-file].md
2. "I acknowledge [Rule#] - [Rule Name]: [Brief description]"
3. READ: $CLAUDE_PROJECT_DIR/rule-library/[second-rule-file].md  
4. "I acknowledge [Rule#] - [Rule Name]: [Brief description]"
[Continue for EVERY rule in PRIMARY DIRECTIVES...]
5. "Ready to execute WAITING_FOR_EFFORT_PLANS work"
```

### 🚨 NO WORK UNTIL ACKNOWLEDGMENT COMPLETE 🚨
**You may NOT begin ANY WAITING_FOR_EFFORT_PLANS work until:**
1. ✅ ALL rules in PRIMARY DIRECTIVES have been READ
2. ✅ ALL rules have been individually ACKNOWLEDGED
3. ✅ You have stated readiness to execute WAITING_FOR_EFFORT_PLANS work
1. **READ** every rule file listed in PRIMARY DIRECTIVES below
2. **ACKNOWLEDGE** each rule individually with number and description
3. **ONLY THEN** proceed with WAITING_FOR_EFFORT_PLANS work

### 🚨 FAILURE TO READ STATE RULES FIRST = IMMEDIATE EXIT 🚨
**If you do ANY WAITING_FOR_EFFORT_PLANS work before reading and acknowledging rules:**
- **STOP ALL WORK IMMEDIATELY**
- **EXIT WITH FAILURE STATUS**
- **YOU HAVE VIOLATED STATE COMPLIANCE**

**THE SYSTEM IS MONITORING YOUR READ TOOL CALLS!**

## 📋 RULE SUMMARY FOR WAITING_FOR_EFFORT_PLANS STATE

### Rules Enforced in This State:
- R234: Mandatory State Traversal [SUPREME LAW - Part of sequence]
- R255: Post-Agent Work Verification [BLOCKING - Check every completion]
- R021: Never Stop Monitoring [SUPREME LAW - Keep checking]
- R287: TODO Save Frequency [BLOCKING - Every 10 messages/15 min]
- R288: State File Update and Commit [SUPREME LAW - Track progress]

### Critical Requirements:
1. Actively poll for plans NOW - Penalty: -30%
2. Check every 5-10 seconds - Penalty: -20%
3. Verify R255 for each completion - Penalty: -100%
4. Save TODOs every 15 minutes - Penalty: -15%
5. Transition to ANALYZE_IMPLEMENTATION_PARALLELIZATION - Penalty: -100%

### Success Criteria:
- ✅ All IMPLEMENTATION-PLAN.md files created
- ✅ All plans in correct directories (R255)
- ✅ All plans committed and pushed
- ✅ Work logs updated

### Failure Triggers:
- ❌ Skip to SPAWN_AGENTS = -100% R234 VIOLATION
- ❌ Accept plans in wrong location = R255 VIOLATION
- ❌ Stop monitoring = R021 VIOLATION
- ❌ Forget TODO saves = -15% per violation

## 🚨 WAITING_FOR_EFFORT_PLANS IS A VERB - START ACTIVELY CHECKING IMMEDIATELY! 🚨

### IMMEDIATE ACTIONS UPON ENTERING WAITING_FOR_EFFORT_PLANS

**THE MOMENT YOU ENTER THIS STATE, YOU MUST:**
1. Poll effort directories for IMPLEMENTATION-PLAN.md NOW
2. Check every 5-10 seconds for completion
3. Check TodoWrite for pending items and process them
4. Report status of each effort immediately

**FORBIDDEN - AUTOMATIC FAILURE:**
- ❌ "STATE TRANSITION COMPLETE: Now in WAITING_FOR_EFFORT_PLANS" [stops]
- ❌ "Successfully entered WAITING_FOR_EFFORT_PLANS state" [waits]
- ❌ "Ready to start actively checking" [pauses]
- ❌ "I'm in WAITING_FOR_EFFORT_PLANS state" [does nothing]

**REQUIRED - IMMEDIATE ACTION:**
- ✅ "Entering WAITING_FOR_EFFORT_PLANS, Poll effort directories for IMPLEMENTATION-PLAN.md NOW..."
- ✅ "START ACTIVELY CHECKING, check every 5-10 seconds for completion..."
- ✅ "WAITING_FOR_EFFORT_PLANS: Report status of each effort immediately..."

## State Context
You are waiting for Code Reviewers to complete individual effort implementation plans.

## 🔴🔴🔴 SUPREME LAW R234 - STAY IN SEQUENCE 🔴🔴🔴

### YOUR POSITION IN THE MANDATORY SEQUENCE:
```
SETUP_EFFORT_INFRASTRUCTURE (✓ completed)
    ↓
ANALYZE_CODE_REVIEWER_PARALLELIZATION (✓ completed)
    ↓
SPAWN_CODE_REVIEWERS_EFFORT_PLANNING (✓ completed)
    ↓
WAITING_FOR_EFFORT_PLANS (👈 YOU ARE HERE)
    ↓ (MUST GO HERE NEXT)
ANALYZE_IMPLEMENTATION_PARALLELIZATION
    ↓
SPAWN_AGENTS
```

**NOW:** Actively monitor Code Reviewers
**NEXT:** You MUST go to ANALYZE_IMPLEMENTATION_PARALLELIZATION
**FORBIDDEN:** Skipping analysis to go directly to SPAWN_AGENTS = -100%

## Monitoring Requirements

```bash
# Check status of effort plans
check_effort_plan_status() {
    local PHASE=$1 WAVE=$2
    local ALL_COMPLETE=true
    
    echo "📊 Checking effort plan status..."
    
    # Check each effort directory for IMPLEMENTATION-PLAN.md
    for effort_dir in efforts/phase${PHASE}/wave${WAVE}/*/; do
        EFFORT=$(basename "$effort_dir")
        PLAN_FILE="$effort_dir/IMPLEMENTATION-PLAN.md"
        
        if [ -f "$PLAN_FILE" ]; then
            echo "✅ $EFFORT: Plan complete"
        else
            echo "⏳ $EFFORT: Plan in progress"
            ALL_COMPLETE=false
        fi
    done
    
    if [ "$ALL_COMPLETE" = true ]; then
        echo "✅ All effort plans complete!"
        return 0
    else
        echo "⏳ Waiting for remaining plans..."
        return 1
    fi
}
```

## Validation Before Proceeding

Before transitioning to SPAWN_AGENTS, verify:

1. **All Plans Exist:**
   ```bash
   for effort_dir in efforts/phase${PHASE}/wave${WAVE}/*/; do
       [ -f "$effort_dir/IMPLEMENTATION-PLAN.md" ] || exit 1
   done
   ```

2. **Plans Include Required Sections:**
   - Implementation approach
   - Test requirements
   - Size limits
   - Dependencies
   - File structure

3. **Work Logs Updated:**
   ```bash
   for effort_dir in efforts/phase${PHASE}/wave${WAVE}/*/; do
       grep -q "Planning complete" "$effort_dir/work-log.md" || echo "Missing"
   done
   ```

## State Transition

Once ALL effort plans are complete:
1. Update orchestrator-state.yaml
2. Record effort plan locations
3. **MANDATORY: Transition to ANALYZE_IMPLEMENTATION_PARALLELIZATION (R234)**
   - DO NOT skip directly to SPAWN_AGENTS
   - MUST analyze SW Engineer parallelization first!
   - This is the MANDATORY sequence

### R287 MONITORING CHECKPOINT
```bash
# Every 15 minutes while monitoring
TIME_SINCE_SAVE=$(($(date +%s) - LAST_TODO_SAVE))
if [ $TIME_SINCE_SAVE -gt 900 ]; then
    echo "⚠️ R287: 15 minutes elapsed - saving TODOs..."
    save_todos "WAITING_FOR_EFFORT_PLANS checkpoint"
    LAST_TODO_SAVE=$(date +%s)
fi

# Every 10 messages
if [ $((MESSAGE_COUNT % 10)) -eq 0 ]; then
    echo "💾 R287: 10 messages - saving TODOs..."
    save_todos "Message checkpoint"
fi
```

### BEFORE TRANSITION
```bash
# R287: State transition trigger
echo "💾 R287: Saving TODOs before state transition..."
save_todos "All effort plans complete"

# R287: Commit within 60 seconds
cd $CLAUDE_PROJECT_DIR
git add todos/*.todo orchestrator-state.yaml
git commit -m "todo: effort plans complete, ready for analysis"
git push
```

## Timeout Handling

If plans not complete within reasonable time:
- Check for blocked Code Reviewers
- Review error logs
- Consider ERROR_RECOVERY state

## Do NOT Proceed If:
- ❌ Any effort missing IMPLEMENTATION-PLAN.md
- ❌ Plans are incomplete or malformed
- ❌ Infrastructure issues detected
- ❌ Code Reviewers report blocking issues

