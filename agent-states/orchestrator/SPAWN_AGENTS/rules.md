# Orchestrator - SPAWN_AGENTS State Rules

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

**YOU HAVE ENTERED SPAWN_AGENTS STATE - YOU MUST READ AND ACKNOWLEDGE ALL STATE RULES BEFORE DOING ANY STATE WORK!**

## 🔴🔴🔴 R290 VERIFICATION REQUIREMENT 🔴🔴🔴

**R290 ENFORCEMENT: CREATE VERIFICATION MARKER AFTER READING**

After reading and acknowledging all state rules, you MUST create a verification marker:

```bash
# MANDATORY: Create verification marker after reading rules
touch .state_rules_read_orchestrator_SPAWN_AGENTS
echo "$(date +%s) - Rules read and acknowledged for SPAWN_AGENTS" > .state_rules_read_orchestrator_SPAWN_AGENTS
```

**FAILURE TO CREATE MARKER = AUTOMATIC -100% PENALTY**

The system will check for this marker. No marker = Immediate failure.

### ❌ DO NOT DO ANY SPAWN_AGENTS WORK UNTIL RULES ARE READ:
- ❌ Start spawn software engineer agents
- ❌ Start assign effort work
- ❌ Start distribute implementation tasks
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
   ❌ WRONG: "I acknowledge all SPAWN_AGENTS rules"
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

### ✅ CORRECT PATTERN FOR SPAWN_AGENTS:
```
1. READ: $CLAUDE_PROJECT_DIR/rule-library/[first-rule-file].md
2. "I acknowledge [Rule#] - [Rule Name]: [Brief description]"
3. READ: $CLAUDE_PROJECT_DIR/rule-library/[second-rule-file].md  
4. "I acknowledge [Rule#] - [Rule Name]: [Brief description]"
[Continue for EVERY rule in PRIMARY DIRECTIVES...]
5. "Ready to execute SPAWN_AGENTS work"
```

### 🚨 NO WORK UNTIL ACKNOWLEDGMENT COMPLETE 🚨
**You may NOT begin ANY SPAWN_AGENTS work until:**
1. ✅ ALL rules in PRIMARY DIRECTIVES have been READ
2. ✅ ALL rules have been individually ACKNOWLEDGED
3. ✅ You have stated readiness to execute SPAWN_AGENTS work
1. **READ** every rule file listed in PRIMARY DIRECTIVES below
2. **ACKNOWLEDGE** each rule individually with number and description
3. **ONLY THEN** proceed with SPAWN_AGENTS work

### 🚨 FAILURE TO READ STATE RULES FIRST = IMMEDIATE EXIT 🚨
**If you do ANY SPAWN_AGENTS work before reading and acknowledging rules:**
- **STOP ALL WORK IMMEDIATELY**
- **EXIT WITH FAILURE STATUS**
- **YOU HAVE VIOLATED STATE COMPLIANCE**

**THE SYSTEM IS MONITORING YOUR READ TOOL CALLS!**

## 📋 PRIMARY DIRECTIVES FOR SPAWN_AGENTS

**YOU MUST READ EACH RULE LISTED HERE. YOUR READ TOOL CALLS ARE BEING MONITORED.**

### State-Specific Rules (NOT in orchestrator.md):

1. **🚨🚨🚨 R006** - ORCHESTRATOR NEVER WRITES CODE OR PERFORMS FILE OPERATIONS (BLOCKING)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R006-orchestrator-never-writes-code.md`
   - Criticality: BLOCKING - Automatic termination, 0% grade
   - Summary: NEVER write, copy, move, or manipulate ANY code files - delegate ALL to agents
   - **CRITICAL**: Copying files is NOT infrastructure - it's implementation work!

2. **🔴🔴🔴 R322** - MANDATORY STOP BEFORE STATE TRANSITIONS (SUPREME LAW)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R322-mandatory-stop-after-spawn.md`
   - Criticality: SUPREME - Automatic -100% failure for violation
   - Summary: MUST STOP IMMEDIATELY after spawning agents to preserve context
   - **CRITICAL**: Record what was spawned, save state, EXIT with clear message

3. **🔴🔴🔴 R251** - UNIVERSAL REPOSITORY SEPARATION LAW (PARAMOUNT)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R251-REPOSITORY-SEPARATION-LAW.md`
   - Criticality: PARAMOUNT - Automatic -100% failure for violation
   - Summary: Software Factory = Planning ONLY, Target Repo = Code ONLY
   - **CRITICAL**: Verify agents will work in TARGET repo clones, not SF repo

4. **🔴🔴🔴 R309** - NEVER Create Efforts in SF Repo (PARAMOUNT LAW)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R309-never-create-efforts-in-sf-repo.md`
   - Criticality: PARAMOUNT - Automatic -100% failure for violation
   - Summary: NEVER spawn agents to work in Software Factory repo
   - **CRITICAL**: Verify effort directories are TARGET clones before spawning!

5. **🚨🚨🚨 R318** - Agent Failure Escalation Protocol (BLOCKING)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R318-agent-failure-escalation-protocol.md`
   - Criticality: BLOCKING - -40% for attempting forbidden fixes
   - Summary: NEVER fix agent failures directly - respawn with better instructions or escalate
   - **CRITICAL**: If agent fails, respawn or escalate - NEVER attempt DIY fixes!

6. **🚨🚨🚨 R151** - Parallel Spawning Timestamp Requirement
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R151-parallel-agent-spawning-timing.md`
   - Criticality: CRITICAL - Timestamps must be within 5s for parallel agents (50% of orchestrator grade)
   - Summary: All parallel agents must emit timestamps within 5 seconds, acknowledge plan creator's decision

7. **🔴🔴🔴 R295** - SW Engineer Spawn Clarity Protocol
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R295-sw-engineer-spawn-clarity-protocol.md`
   - Criticality: SUPREME - MUST specify exact state, plan file, and clear instructions
   - Summary: Every spawn MUST include state name, exact plan file, and warnings about old plans

8. **R052** - Agent Spawning Protocol
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R052-agent-spawning-protocol.md`
   - Criticality: CRITICAL - Complete context and deliverables required
   - Summary: Provide complete context, startup requirements, deliverables, size limits to each agent

9. **R197** - One Agent Per Effort
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R197-one-agent-per-effort.md`
   - Criticality: BLOCKING - Never spawn multiple agents for same effort

10. **R255** - Post-Agent Work Verification
    - File: `$CLAUDE_PROJECT_DIR/rule-library/R255-POST-AGENT-WORK-VERIFICATION.md`
    - Criticality: BLOCKING - Verify correct locations after completion

11. **🚨🚨🚨 R216** - Bash Execution Syntax Protocol (BLOCKING)
    - File: `$CLAUDE_PROJECT_DIR/rule-library/R216-bash-execution-syntax-protocol.md`
    - Criticality: BLOCKING - Incorrect syntax causes failures
    - Summary: Use parentheses for subshells, proper variable syntax
    
12. **🚨🚨🚨 R235** - Pre-flight Verification Checklist (BLOCKING)
    - File: `$CLAUDE_PROJECT_DIR/rule-library/R235-pre-flight-verification-checklist.md`
    - Criticality: BLOCKING - Must verify environment before spawning
    - Summary: Check directories, permissions, branches before agent spawn

**Note**: R208 (CD before spawn), R221 (bash reset) are already in orchestrator.md Supreme Laws.

## 🔴🔴🔴 MANDATORY SPAWN DIRECTORY VERIFICATION PROTOCOL 🔴🔴🔴

### R208 ENFORCEMENT: Pre-Spawn Directory Verification
**EVERY AGENT SPAWN MUST FOLLOW THIS EXACT SEQUENCE:**

```bash
# 1. DETERMINE target directory for the agent
TARGET_DIR="/efforts/phase${PHASE}/wave${WAVE}/${EFFORT_NAME}"
echo "🎯 Target directory for agent: $TARGET_DIR"

# 2. VERIFY directory exists (create if needed)
if [ ! -d "$TARGET_DIR" ]; then
    echo "📁 Creating target directory: $TARGET_DIR"
    mkdir -p "$TARGET_DIR"
fi

# 3. CD to that directory (R208 MANDATORY - NO EXCEPTIONS)
echo "🔄 Changing to target directory..."
cd "$TARGET_DIR" || {
    echo "❌❌❌ FATAL: Cannot change to $TARGET_DIR"
    echo "❌❌❌ R208 VIOLATION: Cannot spawn without CD!"
    exit 208
}

# 4. VERIFY pwd output shows correct directory
ACTUAL_DIR=$(pwd)
echo "✅ Now in: $ACTUAL_DIR"
if [ "$ACTUAL_DIR" != "$TARGET_DIR" ]; then
    echo "❌❌❌ FATAL: Directory mismatch!"
    echo "Expected: $TARGET_DIR"
    echo "Actual: $ACTUAL_DIR"
    exit 208
fi

# 5. SPAWN the agent (ONLY after successful CD and verification)
echo "🚀 Spawning agent from verified directory: $(pwd)"
# [Actual spawn command here]

# 6. RETURN to orchestrator directory
cd "$ORCHESTRATOR_DIR"
echo "📍 Returned to orchestrator directory: $(pwd)"
```

**VIOLATIONS = AUTOMATIC -100% FAILURE:**
- ❌ Spawning without CD'ing first
- ❌ Skipping pwd verification
- ❌ Assuming agent will CD itself
- ❌ Using --working-directory flag instead of CD
- ❌ Not returning to orchestrator directory after spawn

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

## 🚨 SPAWN_AGENTS IS A VERB - START SPAWNING AGENTS IMMEDIATELY! 🚨

**See R151 for immediate action requirements when entering this state.**

The SPAWN_AGENTS state requires IMMEDIATE spawning action - no pausing or waiting.
See rule R151 which includes requirements for immediate state execution.

## State Context
You are spawning SW Engineers to implement efforts based on the implementation plans created by Code Reviewers.

## 🔴🔴🔴 PREREQUISITES FOR SPAWN_AGENTS 🔴🔴🔴

**BEFORE ENTERING THIS STATE, YOU MUST ALREADY HAVE:**
1. ✅ All effort directories created (done in SETUP_EFFORT_INFRASTRUCTURE)
2. ✅ All git clones and branches ready (done in SETUP_EFFORT_INFRASTRUCTURE) 
3. ✅ All effort IMPLEMENTATION-PLAN.md files (created by Code Reviewers)
4. ✅ All work-log.md files initialized
5. ✅ **PARALLELIZATION ANALYSIS COMPLETE (ANALYZE_IMPLEMENTATION_PARALLELIZATION)**
6. ✅ **SW Engineer parallelization plan in orchestrator-state.yaml**

**IF PARALLELIZATION NOT ANALYZED, GO BACK TO ANALYZE_IMPLEMENTATION_PARALLELIZATION!**
**Infrastructure was created BEFORE Code Reviewers made plans!**
**Now you're just spawning SW Engineers to implement using the PRE-ANALYZED strategy.**


## Parallel Spawning


## ✅ Infrastructure Already Ready

Infrastructure was set up in SETUP_EFFORT_INFRASTRUCTURE state:
- Effort directories exist at: `efforts/phase{X}/wave{Y}/{effort-name}`
- Git branches created with project prefix from target-repo-config.yaml
- Remote tracking configured
- IMPLEMENTATION-PLAN.md files created by Code Reviewers
- work-log.md files initialized
- **SW Engineer parallelization plan in orchestrator-state.yaml**

**Just CD to directories and spawn SW Engineers per the analyzed plan!**

### R287 TODO PERSISTENCE + R322 MANDATORY STOP
```bash
# After all SW Engineers spawned
echo "💾 R287: Saving TODOs after spawning SW Engineers..."
save_todos "SPAWN_AGENTS complete - all SW Engineers spawned"

# R287: Commit within 60 seconds
cd $CLAUDE_PROJECT_DIR
git add todos/*.todo orchestrator-state.yaml
git commit -m "todo: SW Engineers spawned, stopping per R322"
git push

# 🔴🔴🔴 R322 MANDATORY STOP AFTER SPAWNING 🔴🔴🔴
echo "
🛑 STOPPING PER R322 - CONTEXT PRESERVATION
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Agents spawned: [List all SW Engineers]
State saved to: orchestrator-state.yaml
Next state: MONITOR

To continue after agents complete:
  claude --continue

This stop preserves context and prevents rule loss.
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
"
exit 0  # MANDATORY EXIT PER R322
```

## Spawn Message Template WITH R295 COMPLIANCE

```markdown
# SPAWN SW ENGINEER WITH MANDATORY CLARITY (R295):
Task sw-engineer:
PURPOSE: Implement {effort_id} - {effort_name}

🔴🔴🔴 CRITICAL STATE INFORMATION (R295):
YOU ARE IN STATE: IMPLEMENTATION
This means you should: Implement the features defined in IMPLEMENTATION-PLAN.md
🔴🔴🔴

📋 YOUR INSTRUCTIONS (R295):
FOLLOW ONLY: IMPLEMENTATION-PLAN.md
LOCATION: In your effort directory
IGNORE: Any files named *-COMPLETED-*.md or other plan files

🎯 CONTEXT:
- EFFORT: {effort_name}
- WAVE: {wave_number}
- PHASE: {phase_number}
- YOUR TASK: Implement features as specified in IMPLEMENTATION-PLAN.md

🔴🔴🔴 CRITICAL: YOU WILL NOT BE IN THE RIGHT DIRECTORY! 🔴🔴🔴
YOU MUST NAVIGATE TO YOUR EFFORT DIRECTORY IMMEDIATELY!

TARGET_DIRECTORY: /efforts/phase{X}/wave{Y}/{effort-name}
EXPECTED_BRANCH: {PROJECT_PREFIX}/phase{X}/wave{Y}/{effort-name}

YOUR MANDATORY FIRST ACTIONS:
1. Echo your current directory: pwd
2. Navigate to effort directory: cd /efforts/phase{X}/wave{Y}/{effort-name}
3. Verify you're now in correct directory: pwd
4. Verify branch: git branch --show-current
5. If directory doesn't exist or branch is wrong:
   - STOP IMMEDIATELY
   - Report: "❌ ENVIRONMENT ERROR: Directory or branch incorrect"
   - Request orchestrator correction
6. Run R209 directory isolation protocol
7. Set readonly EFFORT_ISOLATION_DIR environment variable

BRANCH: {PROJECT_PREFIX}/phase{X}/wave{Y}/effort{Z}-{name}  # Include project prefix from target-repo-config.yaml!

REQUIREMENTS:
- Follow IMPLEMENTATION-PLAN.md exactly (R295: This is your ONLY plan)
- Size limit: {limit} lines
- Test coverage: {X}% minimum
- Update work-log.md every checkpoint

STARTUP REQUIREMENTS:
1. Print: "AGENT STARTUP: $(date '+%Y-%m-%d %H:%M:%S %Z')"
2. Verify pwd matches WORKING_DIR
3. Verify branch matches BRANCH
4. Acknowledge rules R054, R007, R013, R060, R017, R152, R295

DELIVERABLES:
- Implementation complete per IMPLEMENTATION-PLAN.md
- Tests passing at required coverage
- Size under limit
- Work log updated
- Code committed and pushed
```

## Parallelization Matrix

**See R053-parallelization-decisions.md in rule-library for parallelization decisions guidance.**

Key points:
- Can parallelize: Independent efforts, no shared dependencies
- Must serialize: Dependent efforts, splits, shared files

## Recording Spawn Times

```yaml
# Update in orchestrator-state.yaml
parallel_spawn_records:
  wave{X}_group{Y}:
    spawned_at: "2025-08-23T14:30:45Z"
    agents:
      - name: "sw-engineer-effort1"
        timestamp: "2025-08-23T14:30:47Z"
      - name: "sw-engineer-effort2"
        timestamp: "2025-08-23T14:30:49Z"
      - name: "sw-engineer-effort3"
        timestamp: "2025-08-23T14:30:51Z"
    deltas: [2, 2]
    average_delta: 2.0
    grade: "PASS"
```

## Common Spawn Patterns

### Pattern 1: Wave Start (All Planning)
```
Spawn all Code Reviewers for planning → Wait → Spawn all SW Engineers
```

### Pattern 2: Post-Implementation (All Reviews)
```
Spawn all Code Reviewers for review → Process decisions → Handle splits/fixes
```

### Pattern 3: Mixed Dependencies
```
Spawn independent efforts → Monitor → Spawn dependent efforts as prerequisites complete
```

## Grading Calculation

```python
def calculate_spawn_grade(timestamps):
    if len(timestamps) < 2:
        return "PASS"  # Single spawn
    
    deltas = []
    for i in range(1, len(timestamps)):
        delta = (timestamps[i] - timestamps[i-1]).total_seconds()
        deltas.append(delta)
    
    avg = sum(deltas) / len(deltas)
    grade = "PASS" if avg < 5.0 else "FAIL"
    
    print(f"Spawn Grade: {grade}")
    print(f"Average Delta: {avg:.2f}s (target: <5s)")
    
    return grade
```


## R322 VIOLATION DETECTION

If you find yourself:
- Starting work for a new state without /continue-orchestrating
- Transitioning without stopping after state file commit
- Continuing after completing state work

**STOP IMMEDIATELY - You are violating R322!**
