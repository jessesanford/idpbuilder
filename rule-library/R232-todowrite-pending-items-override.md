# 🔴🔴🔴 RULE R232 - TODOWRITE PENDING ITEMS OVERRIDE (SUPREME LAW)

**Criticality:** SUPREME LAW - ABSOLUTE HIGHEST PRIORITY  
**Grading Impact:** -100% IMMEDIATE FAILURE for violation  
**Enforcement:** CONTINUOUS - Before EVERY response completion

## SUPREME LAW STATEMENT

**TODOWRITE PENDING ITEMS ARE COMMANDS, NOT SUGGESTIONS. IF TODOWRITE HAS PENDING ITEMS, YOU MUST CONTINUE WORKING. STOPPING WITH PENDING TODOS IS AUTOMATIC FAILURE.**

## 🚨🚨🚨 THE ABSOLUTE MANDATE 🚨🚨🚨

### TODOWRITE PENDING ITEMS ARE:
```markdown
✅ COMMANDS that MUST be executed
✅ WORK QUEUE that MUST be processed
✅ MANDATORY ACTIONS that CANNOT be ignored
✅ BLOCKING REQUIREMENTS for any stopping decision
✅ ABSOLUTE OVERRIDE of any completion urge
```

### TODOWRITE PENDING ITEMS ARE NOT:
```markdown
❌ Suggestions or recommendations
❌ Optional future work
❌ Things to do "later"
❌ Items that can wait for next session
❌ Ideas or possibilities
```

## 🔴🔴🔴 ENFORCEMENT MECHANISM 🔴🔴🔴

### MANDATORY CHECK BEFORE ANY RESPONSE COMPLETION
```bash
check_todowrite_before_stopping() {
    echo "🔴🔴🔴 R232 TODOWRITE OVERRIDE CHECK 🔴🔴🔴"
    
    # Check TodoWrite for pending items
    pending_count=$(get_pending_todos_count)
    
    if [ $pending_count -gt 0 ]; then
        echo "❌❌❌ VIOLATION DETECTED ❌❌❌"
        echo "TodoWrite has $pending_count PENDING items!"
        echo "YOU ARE FORBIDDEN FROM STOPPING!"
        echo ""
        echo "PENDING TODOS:"
        list_pending_todos
        echo ""
        echo "🚨 MUST CONTINUE WORKING ON THESE ITEMS NOW!"
        return 1  # CANNOT STOP - MUST CONTINUE
    fi
    
    echo "✅ No pending TodoWrite items - may proceed with stop check"
    return 0
}

# THIS CHECK RUNS BEFORE R021 CHECK
# THIS CHECK RUNS BEFORE ANY STOP DECISION
# THIS CHECK IS ABSOLUTE AND CANNOT BE OVERRIDDEN
```

## 🚨🚨🚨 "I WILL..." STATEMENTS ARE LIES 🚨🚨🚨

### THE DEADLY PATTERN TO ELIMINATE:
```markdown
❌ CATASTROPHIC VIOLATION:
"I will spawn Code Reviewer for E4.1.3 split planning"
[Response ends without spawning]

✅ CORRECT BEHAVIOR:
"Spawning Code Reviewer for E4.1.3 split planning now..."
[Actually spawns the agent]
[Response continues with the spawn]
```

### THE TRUTH ABOUT "I WILL..." STATEMENTS:
1. **"I will..." = LIE** - You're saying you'll do it but not doing it
2. **"I am..." = TRUTH** - You're actually doing it right now
3. **"Spawning..." = ACTION** - You're executing the command
4. **"Will continue..." = DECEPTION** - You're stopping instead

### ENFORCEMENT:
```bash
detect_false_promises() {
    # Pattern detection
    if message_contains "I will" && !action_taken; then
        echo "❌❌❌ R232 VIOLATION: FALSE PROMISE DETECTED ❌❌❌"
        echo "You said 'I will' but didn't do it!"
        echo "AUTOMATIC FAILURE - DERELICTION OF DUTY"
        return 1
    fi
}
```

## 📊 TODOWRITE INTEGRATION REQUIREMENTS

### TodoWrite Tool Must Be Checked:
1. **BEFORE** any response ends
2. **BEFORE** any summary is provided
3. **BEFORE** any "completion" message
4. **BEFORE** transitioning to COMPLETE state
5. **ALWAYS** when considering stopping

### TodoWrite Pending Items Mean:
```yaml
pending_status: "in_progress"
implication: "MUST ACT NOW"
priority: "IMMEDIATE"
deferrable: false
ignorable: false
override_stop: true
```

## ❌ VIOLATION EXAMPLES

### VIOLATION 1: Stopping with Pending TODOs
```bash
# TodoWrite shows:
# - ⏳ Spawn Code Reviewer for E4.1.3 split planning
# - ⏳ Spawn Code Reviewer for E4.1.2 review

# Orchestrator says:
"Here's the current status. I'll continue monitoring..."
[RESPONSE ENDS]

# ❌❌❌ CATASTROPHIC VIOLATION ❌❌❌
# GRADE: -100% IMMEDIATE FAILURE
```

### VIOLATION 2: "I Will" Without Action
```bash
# Orchestrator says:
"I will spawn the Code Reviewers for the remaining tasks."
[RESPONSE ENDS WITHOUT SPAWNING]

# ❌❌❌ CATASTROPHIC VIOLATION ❌❌❌
# GRADE: -100% IMMEDIATE FAILURE
```

### VIOLATION 3: Summary Instead of Action
```bash
# TodoWrite has pending items
# Orchestrator provides summary of what needs to be done
# Doesn't actually do any of it

# ❌❌❌ CATASTROPHIC VIOLATION ❌❌❌
# GRADE: -100% IMMEDIATE FAILURE
```

## ✅ CORRECT BEHAVIOR EXAMPLES

### CORRECT 1: Acting on Pending TODOs
```bash
# TodoWrite shows pending items
echo "📋 Processing pending TODOs from TodoWrite..."
echo "🚀 Spawning Code Reviewer for E4.1.3 split planning..."
/spawn code-reviewer --effort E4.1.3 --task split-planning
echo "✅ Spawned. Continuing with next TODO..."
# CONTINUES WORKING
```

### CORRECT 2: "I Am" Instead of "I Will"
```bash
# Instead of "I will spawn..."
echo "🚀 I am spawning Code Reviewer now..."
[ACTUAL SPAWN COMMAND]
echo "✅ Spawned successfully"
# ACTION TAKEN IMMEDIATELY
```

### CORRECT 3: Processing All TODOs Before Stopping
```bash
while has_pending_todos(); do
    next_todo=$(get_next_pending_todo)
    echo "📋 Processing: $next_todo"
    execute_todo "$next_todo"
    mark_todo_complete "$next_todo"
done
echo "✅ All TODOs complete - may now check stop conditions"
```

## 🔴 GRADING IMPACT

### Violations Result In:
```yaml
violation_grades:
  stopped_with_pending_todos: -100%  # IMMEDIATE FAILURE
  false_promise_i_will: -100%        # IMMEDIATE FAILURE
  ignored_todowrite: -100%           # IMMEDIATE FAILURE
  summary_instead_of_action: -100%   # IMMEDIATE FAILURE
  deferred_pending_work: -100%       # IMMEDIATE FAILURE

classification:
  violation_type: "DERELICTION OF DUTY"
  severity: "SUPREME LAW VIOLATION"
  recovery: "IMPOSSIBLE - AUTOMATIC FAILURE"
```

## 🎯 INTEGRATION WITH OTHER RULES

### R232 + R021 (Orchestrator Never Stops):
```bash
# Check order is CRITICAL:
1. FIRST: Check R232 - Any pending TODOs?
   - If YES: MUST CONTINUE (no further checks needed)
   - If NO: Proceed to R021 check
   
2. SECOND: Check R021 - Valid reason to stop?
   - Only if R232 passed (no pending TODOs)
```

### R232 + R231 (Continuous Operation):
- R232 enforces the mechanism
- R231 provides the philosophy
- Together: UNSTOPPABLE ORCHESTRATION

### R232 + State Machine:
- ANY STATE + Pending TODOs = MUST CONTINUE
- COMPLETE state unreachable with pending TODOs
- State transitions blocked by pending work

## 📢 THE TODOWRITE MANTRA

### Repeat Before Every Response Ends:
1. **"Check TodoWrite for pending items"**
2. **"Pending items are commands, not suggestions"**
3. **"I cannot stop with pending work"**
4. **"I will becomes I am - action not promise"**
5. **"Process all TODOs before considering stop"**

## 🛡️ IMPLEMENTATION CHECKLIST

### For Every Agent Response:
```bash
[ ] Check TodoWrite pending items
[ ] If pending > 0: MUST CONTINUE
[ ] No "I will" statements - only "I am"
[ ] Execute actions immediately
[ ] Update TODOs as completed
[ ] Only check stop conditions when pending = 0
```

## 💪 POSITIVE REINFORCEMENT

### Excellence Indicators:
```yaml
todowrite_compliance:
  processed_all_todos: "EXCELLENT"
  no_false_promises: "PROFESSIONAL"
  immediate_action: "PERFECT EXECUTION"
  zero_pending_stops: "EXEMPLARY"
  will_to_am_conversion: "OUTSTANDING"
```

## 📜 THE TODOWRITE OATH

```
I, the Orchestrator, swear by R232:

TodoWrite pending items are COMMANDS, not suggestions.
I WILL NOT stop while TODOs are pending.
I WILL NOT make false promises with "I will..."
I WILL take immediate action with "I am..."
I WILL process every TODO before considering completion.
I WILL check TodoWrite before EVERY response ends.

My TODOs are my work queue.
My work queue MUST be processed.
Pending means DO IT NOW.
Not later. Not "I will." NOW.

This is SUPREME LAW.
Violation means FAILURE.
I WILL EXECUTE, NOT PROMISE.
```

---
**Remember:** An orchestrator that stops with pending TODOs has FAILED. An orchestrator that says "I will" without doing has LIED. TodoWrite is your COMMAND QUEUE, not a suggestion box.