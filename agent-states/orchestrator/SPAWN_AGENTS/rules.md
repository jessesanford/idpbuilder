# Orchestrator - SPAWN_AGENTS State Rules

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

1. **🚨🚨🚨 R151** - Parallel Spawning Timestamp Requirement
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R151-parallel-agent-spawning-timing.md`
   - Criticality: CRITICAL - Timestamps must be within 5s for parallel agents (50% of orchestrator grade)
   - Summary: All parallel agents must emit timestamps within 5 seconds, acknowledge plan creator's decision

2. **R052** - Agent Spawning Protocol
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R052-agent-spawning-protocol.md`
   - Criticality: CRITICAL - Complete context and deliverables required
   - Summary: Provide complete context, startup requirements, deliverables, size limits to each agent

3. **R197** - One Agent Per Effort
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R197-one-agent-per-effort.md`
   - Criticality: BLOCKING - Never spawn multiple agents for same effort

4. **R255** - Post-Agent Work Verification
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R255-POST-AGENT-WORK-VERIFICATION.md`
   - Criticality: BLOCKING - Verify correct locations after completion

**Note**: R208 (CD before spawn), R221 (bash reset) are already in orchestrator.md Supreme Laws.

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

### R287-R287 TODO PERSISTENCE
```bash
# After all SW Engineers spawned
echo "💾 R287: Saving TODOs after spawning SW Engineers..."
save_todos "SPAWN_AGENTS complete - all SW Engineers spawned"

# R287: Commit within 60 seconds
cd $CLAUDE_PROJECT_DIR
git add todos/*.todo orchestrator-state.yaml
git commit -m "todo: SW Engineers spawned, entering monitoring"
git push

# Transition to MONITOR
echo "➡️ Transitioning to MONITOR state..."
```

## Spawn Message Template WITH EXPLICIT DIRECTORY NAVIGATION

```markdown
# SPAWN SW ENGINEER WITH MANDATORY DIRECTORY NAVIGATION:
Task sw-engineer:
PURPOSE: Implement {effort_id} - {effort_name}

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
- Follow IMPLEMENTATION-PLAN.md exactly
- Size limit: {limit} lines
- Test coverage: {X}% minimum
- Update work-log.md every checkpoint

STARTUP REQUIREMENTS:
1. Print: "AGENT STARTUP: $(date '+%Y-%m-%d %H:%M:%S %Z')"
2. Verify pwd matches WORKING_DIR
3. Verify branch matches BRANCH
4. Acknowledge rules R054, R007, R013, R060, R017, R152

DELIVERABLES:
- Implementation complete
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

