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

5. **Skipping Rules in PRIMARY DIRECTIVES**:
   ```
   ❌ WRONG: Reading only some rules from the list
   (ALL rules in PRIMARY DIRECTIVES are MANDATORY)
   ```

### ✅ CORRECT PATTERN FOR INIT:
```
1. READ: $CLAUDE_PROJECT_DIR/rule-library/[first-rule-file].md
2. "I acknowledge [Rule#] - [Rule Name]: [Brief description]"
3. READ: $CLAUDE_PROJECT_DIR/rule-library/[second-rule-file].md  
4. "I acknowledge [Rule#] - [Rule Name]: [Brief description]"
[Continue for EVERY rule in PRIMARY DIRECTIVES...]
5. "Ready to execute INIT work"
```

### 🚨 NO WORK UNTIL ACKNOWLEDGMENT COMPLETE 🚨
**You may NOT begin ANY INIT work until:**
1. ✅ ALL rules in PRIMARY DIRECTIVES have been READ
2. ✅ ALL rules have been individually ACKNOWLEDGED
3. ✅ You have stated readiness to execute INIT work
1. **READ** every rule file listed in PRIMARY DIRECTIVES below
2. **ACKNOWLEDGE** each rule individually with number and description
3. **ONLY THEN** proceed with INIT work

### 🚨 FAILURE TO READ STATE RULES FIRST = IMMEDIATE EXIT 🚨
**If you do ANY INIT work before reading and acknowledging rules:**
- **STOP ALL WORK IMMEDIATELY**
- **EXIT WITH FAILURE STATUS**
- **YOU HAVE VIOLATED STATE COMPLIANCE**

**THE SYSTEM IS MONITORING YOUR READ TOOL CALLS!**

## ⚠️⚠️⚠️ MANDATORY RULE READING AND ACKNOWLEDGMENT ⚠️⚠️⚠️

**YOU MUST READ EACH RULE FILE LISTED IN PRIMARY DIRECTIVES. YOUR READ TOOL CALLS ARE BEING MONITORED.**

### ❌ ANTI-PATTERNS THAT WILL CAUSE FAILURE:
1. Fake acknowledgment without reading
2. Bulk acknowledgment
3. Reading from memory

### ✅ CORRECT PATTERN:
1. READ each rule file
2. Acknowledge individually with rule number and description

## 📋 PRIMARY DIRECTIVES FOR INIT STATE

### 🚨🚨🚨 R191 - Target Repository Configuration
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R191-target-repo-config.md`
**Criticality**: BLOCKING - Must load target-repo-config.yaml
**Summary**: Load target repository configuration from specified file

### 🚨🚨🚨 R192 - Repository Separation Enforcement
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R192-repository-separation.md`
**Criticality**: BLOCKING - Keep Software Factory and target repos separate
**Summary**: Maintain strict separation between factory and target code

### 🚨🚨🚨 R203 - State-Aware Agent Startup
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R203-state-aware-agent-startup.md`
**Criticality**: BLOCKING - Must follow startup sequence
**Summary**: Load config, determine state, load state rules, acknowledge

### 🚨🚨🚨 R287 - TODO Persistence Suite
**Files**:
- R287: `$CLAUDE_PROJECT_DIR/rule-library/R287-mandatory-todo-save-triggers.md`
- R287: `$CLAUDE_PROJECT_DIR/rule-library/R287-todo-save-frequency-requirements.md`
- R287: `$CLAUDE_PROJECT_DIR/rule-library/R287-todo-file-commit-protocol.md`
- R287: `$CLAUDE_PROJECT_DIR/rule-library/R287-todo-recovery-verification.md`
**Criticality**: BLOCKING - TODO loss = -50% to -100% penalty
**Summary**: Save TODOs within 30s, every 10 messages/15 min, commit within 60s

### 🔴🔴🔴 R235 - Mandatory Pre-Flight Verification (SUPREME LAW)
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R235-mandatory-pre-flight-verification.md`
**Criticality**: SUPREME LAW - Wrong location = -100% failure
**Summary**: Verify correct directory and branch before ANY work

### 🔴🔴🔴 R288 - State File Update and Commit Protocol (SUPREME LAW)
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R288-state-file-update-and-commit-protocol.md`
**Criticality**: SUPREME LAW - Update on every transition
**Summary**: Update orchestrator-state.yaml on all state changes

### 🚨🚨🚨 R288 - State File Update and Commit Protocol
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R288-state-file-update-and-commit-protocol.md`
**Criticality**: BLOCKING - Push within 60 seconds
**Summary**: Commit and push state file immediately after updates

### 🔴🔴🔴 R281 - Mandatory Complete State File Initialization (SUPREME LAW #7)
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R281-initial-state-file-creation.md`
**Criticality**: SUPREME LAW #7 - Incomplete state = -100% failure
**Summary**: Create COMPLETE state file with ALL phases, waves, efforts from plan

### 🔴🔴🔴 R232 - TodoWrite Pending Items Override (SUPREME LAW)
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R232-todowrite-pending-items-override.md`
**Criticality**: SUPREME LAW - Pending items are COMMANDS
**Summary**: Any pending TODO items must be executed immediately

## 📋 RULE ENFORCEMENT SUMMARY FOR INIT STATE

### Critical Requirements:
1. Check for compaction IMMEDIATELY - Penalty: -100%
2. Load orchestrator-state.yaml NOW - Penalty: -50%
3. Load target-repo-config.yaml NOW - Penalty: -50%
4. Run R203 startup sequence - Penalty: -100%
5. Check TodoWrite for pending items - Penalty: -30%
6. Save TODOs within 30 seconds - Penalty: -20%
7. Create COMPLETE state file (R281) - Penalty: -100%

### Success Criteria:
- ✅ Compaction check completed
- ✅ All configs loaded successfully
- ✅ TODOs checked and processed
- ✅ State file updated
- ✅ Next state determined

### Failure Triggers:
- ❌ Skip compaction check = IMMEDIATE STOP
- ❌ Skip R203 startup = -100% penalty
- ❌ Forget TODO persistence = -20% per violation
- ❌ Stop after entering state = AUTOMATIC FAILURE

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
# READ: $CLAUDE_PROJECT_DIR/rule-library/R287-mandatory-todo-save-triggers.md
# READ: $CLAUDE_PROJECT_DIR/rule-library/R287-todo-save-frequency-requirements.md  
# READ: $CLAUDE_PROJECT_DIR/rule-library/R287-todo-commit-protocol.md
# READ: $CLAUDE_PROJECT_DIR/rule-library/R287-todo-recovery-verification.md

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
3. **R191**: Load target repository configuration
4. **R281**: Create COMPLETE state file with ALL phases/waves/efforts (SUPREME LAW #7)
5. **R252**: Verify/update state file integrity
6. **R232**: Check and execute pending TodoWrite items
7. **R287**: Save TODOs within 30 seconds of changes
8. **R287**: Save TODOs every 10 messages/15 minutes
9. **R287**: Commit TODOs within 60 seconds
10. **R287**: Verify TODO recovery after compaction
11. **R235**: Run pre-flight verification before ANY work

**Remember**: All SUPREME LAWs override other rules. Violation of any BLOCKING rule = immediate failure.
