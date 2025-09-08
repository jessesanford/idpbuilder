# Orchestrator - INIT State Rules

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


## 🔴🔴🔴 R290 ENFORCEMENT: STATE RULES PRECEDENCE 🔴🔴🔴

**YOU HAVE ENTERED INIT STATE**

Per SUPREME LAW #3 (R290), you MUST:
1. Have already read the 15 SUPREME LAWS from orchestrator.md
2. Have already read the 4 MANDATORY RULES from orchestrator.md  
3. Now read these INIT-specific rules before ANY state work

## 🔴🔴🔴 R290 VERIFICATION REQUIREMENT 🔴🔴🔴

After acknowledging state rules, create verification marker:

```bash
# MANDATORY: Create verification marker after reading rules
touch .state_rules_read_orchestrator_INIT
echo "$(date +%s) - Rules read and acknowledged for INIT" > .state_rules_read_orchestrator_INIT
```

**FAILURE TO CREATE MARKER = AUTOMATIC -100% PENALTY**

## Prerequisites

**Before reading this file, you MUST have already:**
- ✅ Read all 15 SUPREME LAWS from orchestrator.md
- ✅ Read all 4 MANDATORY RULES from orchestrator.md
- ✅ Acknowledged each rule individually

**Do NOT re-read rules already covered in orchestrator.md**

## 🔴🔴🔴 MANDATORY MEASUREMENT RULE - R304 🔴🔴🔴

**ABSOLUTE REQUIREMENTS:**
- ✅ MUST use `$CLAUDE_PROJECT_DIR/tools/line-counter.sh` for ALL line counting
- ❌ NEVER use `wc -l` or any manual counting method
- ❌ NEVER count lines any other way - this is a -100% automatic failure
- ✅ MUST specify both -b (base branch) and -c (current branch) parameters
- ✅ Base branch MUST be phase integration branch (NOT "main")

**Failure to use the line counter tool = AUTOMATIC -100% GRADE**

## INIT-Specific Rules (2 UNIQUE FILES)

### 1. 🚨🚨🚨 R191 - Target Repository Configuration
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R191-target-repo-config.md`
**Criticality**: BLOCKING - Missing config = cannot proceed
**Summary**: Load and validate target-repo-config.yaml for repository settings
**INIT Relevance**: Must load configuration to determine repository structure

### 2. 🚨🚨🚨 R192 - Repository Separation Protocol  
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R192-repo-separation.md`
**Criticality**: BLOCKING - Mixing repositories = corruption
**Summary**: Keep Software Factory and target repository completely separate
**INIT Relevance**: Must verify separation before initializing workspace

## INIT State Immediate Actions

### THE MOMENT YOU ENTER THIS STATE:

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

# STEP 1: LOAD STATE FILE OR CREATE NEW
STATE_FILE="$CLAUDE_PROJECT_DIR/orchestrator-state.json"
if [ -f "$STATE_FILE" ]; then
    CURRENT_STATE=$(jq '.current_state' "$STATE_FILE")
    echo "✅ Found existing state: $CURRENT_STATE"
    # Transition to that state
else
    echo "🔴🔴🔴 R281: Creating COMPLETE initial state file..."
    # Parse PROJECT-IMPLEMENTATION-PLAN.md
    # Create state with ALL phases/waves/efforts
    # Use templates/initial-state-template.yaml
    # Validate with utilities/validate-state-completeness.sh
fi

# STEP 2: LOAD TARGET CONFIG (R191 - INIT-specific)
CONFIG_FILE="$CLAUDE_PROJECT_DIR/target-repo-config.yaml"
if [ ! -f "$CONFIG_FILE" ]; then
    echo "🔴🔴🔴 CRITICAL: target-repo-config.yaml NOT FOUND!"
    exit 191
fi

# STEP 3: VERIFY REPOSITORY SEPARATION (R192 - INIT-specific)
# Ensure Software Factory and target repo are separate
if [[ "$PWD" == *"software-factory"* ]]; then
    echo "✅ R192: In Software Factory directory"
else
    echo "❌ R192: Not in Software Factory directory!"
    exit 192
fi

# STEP 4: CHECK TODOS (R232 from orchestrator.md)
# Check TodoWrite for any pending initialization tasks

# STEP 5: VERIFY ENVIRONMENT (R235 from orchestrator.md)
pwd
git branch --show-current

# STEP 6: DETERMINE AND TRANSITION
# Determine next state based on state file existence
```

## State Context
Initial state when orchestrator starts. Primary purpose is to:
1. Load or create orchestrator-state.json
2. Verify environment configuration
3. Transition to appropriate next state

## Primary Actions
1. Check for compaction and recover if needed
2. Load target-repo-config.yaml (R191 - INIT-specific)
3. Verify repository separation (R192 - INIT-specific)  
4. Check/create orchestrator-state.json (R281 from orchestrator.md)
5. If creating new state file:
   - Parse PROJECT-IMPLEMENTATION-PLAN.md completely
   - Extract ALL phases, waves, and efforts
   - Create COMPLETE state file with every item
   - Validate completeness
6. Determine if resuming or starting fresh
7. Transition to appropriate state immediately

## State Transitions
- If orchestrator-state.json exists with current_state → Resume from that state
- If no state file → Create COMPLETE initial state (R281), then go to SPAWN_ARCHITECT_PHASE_PLANNING
- If state file incomplete → IMMEDIATE FAILURE (-100%)
- If config missing → HARD_STOP

## Critical Enforcement Order

**From orchestrator.md (already read):**
- R283, R288, R232, R322, R322, R281, R280, SOFTWARE-FACTORY-STATE-MACHINE.md
- R235, R221, R290, R208, R234, R307, R308
- R287, R216, R206, R203

**INIT-specific (read now):**
1. **R191**: Load target repository configuration
2. **R192**: Verify repository separation

## Remember
- All SUPREME LAWs from orchestrator.md take precedence
- INIT is a VERB - start initializing IMMEDIATELY (R322)
- No idling, no pausing, no waiting (R322, R280)
- Save TODOs per R287 requirements
- Update state file per R288 requirements
## R322 VIOLATION DETECTION

If you find yourself:
- Starting work for a new state without /continue-orchestrating
- Transitioning without stopping after state file commit
- Continuing after completing state work

**STOP IMMEDIATELY - You are violating R322!**
