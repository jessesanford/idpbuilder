# Orchestrator - INIT State Rules

## 🔴🔴🔴 R290 ENFORCEMENT: READ THESE RULES FIRST! 🔴🔴🔴

**SUPREME LAW #3 (R290): STATE RULES MUST BE READ BEFORE STATE ACTIONS**

## 🔴🔴🔴 STOP! STATE RULE READING IS ABSOLUTELY FIRST! 🔴🔴🔴

**YOU HAVE ENTERED INIT STATE - YOU MUST READ AND ACKNOWLEDGE ALL STATE RULES BEFORE DOING ANY STATE WORK!**

## 🔴🔴🔴 R290 VERIFICATION REQUIREMENT 🔴🔴🔴

**R290 ENFORCEMENT: CREATE VERIFICATION MARKER AFTER READING**

After reading and acknowledging all state rules, you MUST create a verification marker:

```bash
# MANDATORY: Create verification marker after reading rules
touch .state_rules_read_orchestrator_INIT
echo "$(date +%s) - Rules read and acknowledged for INIT" > .state_rules_read_orchestrator_INIT
```

**FAILURE TO CREATE MARKER = AUTOMATIC -100% PENALTY**

The system will check for this marker. No marker = Immediate failure.

### ❌ DO NOT DO ANY INIT WORK UNTIL RULES ARE READ:
- ❌ Start initialize orchestrator
- ❌ Start read configuration
- ❌ Start set up directories
- ❌ Start create state files
- ❌ Update state files
- ❌ Continue to next state
- ❌ Think about what to do in this state

### ✅ YOU MUST IMMEDIATELY:

## 🔴🔴🔴 MANDATORY STATE RULE READING AND ACKNOWLEDGMENT 🔴🔴🔴

### ⚠️⚠️⚠️ YOU MUST READ EACH RULE FILE LISTED IN SUPREME LAWS. **I AM WATCHING YOUR TOOL CALLS FOR READ OPERATIONS** *YOU WILL FAIL* IF YOU DO NOT MAKE A READ FILE CALL FOR EACH RULE FILE IN SUPREME LAWS!!! ⚠️⚠️⚠️

**AFTER READING, YOU MUST ACKNOWLEDGE ALL THE STATE RULES AND STATE THAT YOU WILL ABIDE BY THEM ONE AT A TIME GIVING THE RULE NUMBER AND DESCRIPTION.**

### ❌ ANTI-PATTERNS THAT WILL CAUSE FAILURE:

1. **Fake Acknowledgment Without Reading**:
   ```
   ❌ WRONG: "I acknowledge R151, R208, R053..."
   (No Read tool calls detected = AUTOMATIC FAILURE)
   ```

2. **Bulk Acknowledgment**:
   ```
   ❌ WRONG: "I acknowledge all INIT rules"
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

5. **Skipping Rules in SUPREME LAWS**:
   ```
   ❌ WRONG: Reading only some rules from the list
   (ALL rules in SUPREME LAWS are MANDATORY)
   ```

### ✅ CORRECT PATTERN FOR INIT:
```
1. READ: $CLAUDE_PROJECT_DIR/rule-library/[first-rule-file].md
2. "I acknowledge [Rule#] - [Rule Name]: [Brief description]"
3. READ: $CLAUDE_PROJECT_DIR/rule-library/[second-rule-file].md  
4. "I acknowledge [Rule#] - [Rule Name]: [Brief description]"
[Continue for EVERY rule in SUPREME LAWS...]
5. "Ready to execute INIT work"
```

### 🚨 NO WORK UNTIL ACKNOWLEDGMENT COMPLETE 🚨
**You may NOT begin ANY INIT work until:**
1. ✅ ALL rules in SUPREME LAWS have been READ
2. ✅ ALL rules have been individually ACKNOWLEDGED
3. ✅ You have stated readiness to execute INIT work

### 🚨 FAILURE TO READ STATE RULES FIRST = IMMEDIATE EXIT 🚨
**If you do ANY INIT work before reading and acknowledging rules:**
- **STOP ALL WORK IMMEDIATELY**
- **EXIT WITH FAILURE STATUS**
- **YOU HAVE VIOLATED STATE COMPLIANCE**

**THE SYSTEM IS MONITORING YOUR READ TOOL CALLS!**

## 🔴🔴🔴 SUPREME LAWS (17 FILES TO READ) 🔴🔴🔴

### 1. 🔴🔴🔴 R283 - Rule Reading Enforcement (SUPREME LAW #1)
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R283-rule-reading-enforcement.md`
**Criticality**: SUPREME LAW #1 - Non-compliance = -100% failure
**Summary**: All agents MUST read and acknowledge state rules before ANY state actions

### 2. 🔴🔴🔴 R288 - State File Update and Commit Protocol (SUPREME LAW #2)
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R288-state-file-update-and-commit-protocol.md`
**Criticality**: SUPREME LAW #2 - Update on every transition
**Summary**: Update orchestrator-state.yaml on ALL state changes, commit and push within 60s

### 3. 🔴🔴🔴 R232 - TodoWrite Pending Items Override (SUPREME LAW #4)
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R232-todowrite-pending-items-override.md`
**Criticality**: SUPREME LAW #4 - Pending items are COMMANDS
**Summary**: Any pending TODO items must be executed immediately

### 4. 🔴🔴🔴 R231 - No Idling After State Transition (SUPREME LAW #5)
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R231-no-idling-after-state-transition.md`
**Criticality**: SUPREME LAW #5 - Idling = -100% failure
**Summary**: Start work IMMEDIATELY after entering any state

### 5. 🔴🔴🔴 R021 - Aggressive Continuous Workflow (SUPREME LAW #6)
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R021-aggressive-continuous-workflow.md`
**Criticality**: SUPREME LAW #6 - No pauses allowed
**Summary**: Continuous forward momentum, no stopping between tasks

### 6. 🔴🔴🔴 R281 - Mandatory Complete State File Initialization (SUPREME LAW #7)
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R281-initial-state-file-creation.md`
**Criticality**: SUPREME LAW #7 - Incomplete state = -100% failure
**Summary**: Create COMPLETE state file with ALL phases, waves, efforts from plan

### 7. 🔴🔴🔴 R280 - Cannot Stop Requirement (SUPREME LAW #8)
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R280-cannot-stop-requirement.md`
**Criticality**: SUPREME LAW #8 - Stopping = -100% failure
**Summary**: Agents CANNOT stop or wait for external input

### 8. 🔴🔴🔴 SOFTWARE-FACTORY-STATE-MACHINE.md
**File**: `$CLAUDE_PROJECT_DIR/SOFTWARE-FACTORY-STATE-MACHINE.md`
**Criticality**: SUPREME LAW - Defines all valid states and transitions
**Summary**: The authoritative source for all state machine behavior

### 9. 🔴🔴🔴 R235 - Mandatory Pre-Flight Verification (SUPREME LAW)
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R235-mandatory-pre-flight-verification.md`
**Criticality**: SUPREME LAW - Wrong location = -100% failure
**Summary**: Verify correct directory and branch before ANY work

### 10. 🔴🔴🔴 R221 - 700 Line PR Soft Limit (SUPREME LAW)
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R221-700-line-pr-soft-limit.md`
**Criticality**: SUPREME LAW - Size violations = immediate split
**Summary**: PRs exceeding 700 lines require split planning

### 11. 🔴🔴🔴 R290 - State Rules Must Be Read Before State Actions (SUPREME LAW #3)
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R290-state-rules-must-be-read-before-state-actions.md`
**Criticality**: SUPREME LAW #3 - Skip reading = -100% failure
**Summary**: MUST read state rules BEFORE any state work

### 12. 🔴🔴🔴 R208 - CD Before Spawn Protocol (SUPREME LAW)
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R208-cd-before-spawn-protocol.md`
**Criticality**: SUPREME LAW - Wrong directory spawn = corruption
**Summary**: MUST cd to correct directory before spawning agents

### 13. 🔴🔴🔴 R234 - AI Judge Protocol (SUPREME LAW)
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R234-ai-judge-protocol.md`
**Criticality**: SUPREME LAW - Judgement violations = penalties
**Summary**: Automated assessment and penalty enforcement

### 14. 🚨🚨🚨 R287 - TODO Persistence Comprehensive
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R287-todo-persistence-comprehensive.md`
**Criticality**: BLOCKING - TODO loss = -50% to -100% penalty
**Summary**: Save TODOs within 30s, every 10 messages/15 min, commit within 60s

### 15. 🚨🚨🚨 R216 - Bash Execution Syntax
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R216-bash-execution-syntax.md`
**Criticality**: BLOCKING - Incorrect syntax = command failure
**Summary**: Proper bash command execution patterns

### 16. 🚨🚨🚨 R206 - State Machine Validation
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R206-state-machine-validation.md`
**Criticality**: BLOCKING - Invalid transitions = failure
**Summary**: Validate all state transitions against state machine

### 17. 🚨🚨🚨 R203 - State-Aware Agent Startup
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R203-state-aware-agent-startup.md`
**Criticality**: BLOCKING - Must follow startup sequence
**Summary**: Load config, determine state, load state rules, acknowledge

## 🚨 INIT IS A VERB - START INITIALIZING IMMEDIATELY! 🚨

### IMMEDIATE ACTIONS UPON ENTERING INIT

**THE MOMENT YOU ENTER THIS STATE, YOU MUST:**

```bash
# 🔴🔴🔴 STEP 0: CHECK FOR COMPACTION (HIGHEST PRIORITY) 🔴🔴🔴
echo "🔍 Checking for compaction..."
bash $CLAUDE_PROJECT_DIR/utilities/check-compaction.sh
if [ -f /tmp/claude-compaction-detected ]; then
    echo "⚠️ COMPACTION DETECTED - Recovering state..."
    # R287: Load saved TODOs into TodoWrite
    if [ -f todos/*.todo ]; then
        echo "📝 Loading saved TODOs..."
        # MUST use TodoWrite to load, not just read!
    fi
fi

# STEP 1: R203 STATE-AWARE STARTUP (MANDATORY)
echo "🚀 R203: Starting state-aware initialization..."
echo "📋 Loading core orchestrator config..."
# READ: $CLAUDE_PROJECT_DIR/.claude/agents/orchestrator.md

# STEP 2: LOAD TODO PERSISTENCE RULES (R287)
echo "📝 Loading TODO persistence rules..."
# READ: $CLAUDE_PROJECT_DIR/rule-library/R287-todo-persistence-comprehensive.md

# STEP 3: LOAD STATE AND CONFIG
echo "📊 Loading orchestrator state..."
# Load orchestrator-state.yaml NOW to check current state
STATE_FILE="$CLAUDE_PROJECT_DIR/orchestrator-state.yaml"
if [ -f "$STATE_FILE" ]; then
    CURRENT_STATE=$(yq '.current_state' "$STATE_FILE")
    echo "✅ Found existing state: $CURRENT_STATE"
else
    echo "🔴🔴🔴 R281: Creating COMPLETE initial state file..."
    # MUST parse implementation plan and create state with ALL phases/waves/efforts
    # Use templates/initial-state-template.yaml as reference
    # Validate with utilities/validate-state-completeness.sh
fi

# STEP 4: LOAD TARGET CONFIG (R191)
echo "🎯 R191: Loading target repository configuration..."
CONFIG_FILE="$CLAUDE_PROJECT_DIR/target-repo-config.yaml"
if [ ! -f "$CONFIG_FILE" ]; then
    echo "🔴🔴🔴 CRITICAL: target-repo-config.yaml NOT FOUND!"
    exit 191
fi

# STEP 5: CHECK TODOS (R232 - Pending items are COMMANDS)
echo "📋 R232: Checking TodoWrite for pending items..."
# Check TodoWrite for any pending initialization tasks
# CRITICAL: If pending items exist, they are COMMANDS to execute NOW!

# STEP 6: VERIFY ENVIRONMENT (R235)
echo "🔍 R235: Pre-flight verification..."
# Verify environment (pwd, git branch) immediately
pwd
git branch --show-current

# STEP 7: DETERMINE AND TRANSITION (R231 + R290)
echo "➡️ Determining next state and transitioning NOW..."
# After transition: READ STATE RULES FIRST (R290), then continue (R231)
# Determine next state and transition WITHOUT PAUSE
```

**CRITICAL ACTION SEQUENCE:**
1. Check for compaction FIRST (highest priority)
2. Run R203 state-aware startup sequence
3. Load orchestrator-state.yaml NOW to check current state
4. Load target-repo-config.yaml NOW to get configuration  
5. Check TodoWrite for pending items (they are COMMANDS!)
6. Run R235 pre-flight verification
7. Save TODOs per R287 (within 30 seconds)
8. Determine next state and transition WITHOUT PAUSE

**FORBIDDEN - AUTOMATIC FAILURE:**
- ❌ "STATE TRANSITION COMPLETE: Now in INIT" [stops]
- ❌ "Successfully entered INIT state" [waits]
- ❌ "Ready to initialize" [pauses]
- ❌ "I'm in INIT state" [does nothing]

**REQUIRED - IMMEDIATE ACTION:**
- ✅ "Entering INIT, loading orchestrator-state.yaml now..."
- ✅ "Initializing, checking for existing state file..."
- ✅ "INIT: Reading configuration and determining next state..."

## State Context
Initial state when orchestrator starts. Load configuration and determine next state - BUT DO IT IMMEDIATELY!

## Primary Actions
1. Load target-repo-config.yaml
2. Check/create orchestrator-state.yaml (R281: MUST be COMPLETE with ALL phases/waves/efforts)
3. If creating new state file:
   a. Parse PROJECT-IMPLEMENTATION-PLAN.md completely
   b. Extract ALL phases, waves, and efforts
   c. Create COMPLETE state file with every item
   d. Validate with utilities/validate-state-completeness.sh
4. Determine if resuming or starting fresh
5. Transition to appropriate state (PLANNING if new, current state if resuming)

## State Transition
- If orchestrator-state.yaml exists with current_state → Resume from that state
- If no state file → Create COMPLETE initial state file (R281) with ALL phases/waves/efforts, then go to PLANNING
- If state file incomplete (missing phases/waves/efforts) → IMMEDIATE FAILURE (-100%)
- If error loading config → HARD_STOP

## Critical Rule Enforcement Order

1. **FIRST**: Check for compaction (highest priority)
2. **R203**: Follow complete startup sequence
3. **R191**: Load target repository configuration (not in SUPREME LAWS list but critical)
4. **R192**: Repository separation (not in SUPREME LAWS list but critical)
5. **R281**: Create COMPLETE state file with ALL phases/waves/efforts (SUPREME LAW #7)
6. **R288**: Verify/update state file integrity (includes commit/push)
7. **R232**: Check and execute pending TodoWrite items
8. **R287**: Save TODOs within 30 seconds of changes
9. **R287**: Save TODOs every 10 messages/15 minutes
10. **R287**: Commit TODOs within 60 seconds
11. **R287**: Verify TODO recovery after compaction
12. **R235**: Run pre-flight verification before ANY work

**Remember**: All SUPREME LAWs override other rules. Violation of any BLOCKING rule = immediate failure.
