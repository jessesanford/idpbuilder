# Orchestrator - INIT State Rules

## 🔴🔴🔴 MANDATORY STATE RULE READING AND ACKNOWLEDGMENT 🔴🔴🔴

### ⚠️⚠️⚠️ YOU MUST READ EACH RULE FILE LISTED IN PRIMARY DIRECTIVES. **I AM WATCHING YOUR TOOL CALLS FOR READ OPERATIONS** *YOU WILL FAIL* IF YOU DO NOT MAKE A READ FILE CALL FOR EACH RULE FILE IN PRIMARY DIRECTIVES!!! ⚠️⚠️⚠️

**AFTER READING, YOU MUST ACKNOWLEDGE ALL THE STATE RULES AND STATE THAT YOU WILL ABIDE BY THEM ONE AT A TIME GIVING THE RULE NUMBER AND DESCRIPTION.**

## 📋 PRIMARY DIRECTIVES FOR INIT STATE

### Core Mandatory Rules (ALL orchestrator states must have these):

1. **🚨🚨🚨 R006** - ORCHESTRATOR NEVER WRITES CODE OR PERFORMS FILE OPERATIONS (BLOCKING)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R006-orchestrator-never-writes-code.md`
   - Criticality: BLOCKING - Automatic termination, 0% grade
   - Summary: NEVER write, copy, move, or manipulate ANY code files - delegate ALL to agents

2. **🔴🔴🔴 R287** - TODO PERSISTENCE COMPREHENSIVE (SUPREME LAW)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R287-todo-persistence-comprehensive.md`
   - Criticality: SUPREME - -20% to -100% penalty for violations
   - Summary: MUST save TODOs within 30s after write, every 10 messages, before transitions

3. **🔴🔴🔴 R288** - STATE FILE UPDATE REQUIREMENTS (SUPREME LAW)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R288-state-file-update-and-commit-protocol.md`
   - Criticality: SUPREME - State updates required for all transitions
   - Summary: MUST update orchestrator-state-v3.json before EVERY state transition

4. **🔴🔴🔴 R322 Part A** - Mandatory Stop After Spawn States
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R322-mandatory-stop-before-state-transitions.md`
   - Criticality: SUPREME LAW - Must stop after spawning
   - Summary: ALL spawn states require STOP after spawning agents

### State-Specific Rules:

5. **🚨🚨🚨 R191** - Target Repository Configuration
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R191-target-repo-config.md`
   - Criticality: BLOCKING - Missing config = cannot proceed
   - Summary: Load and validate target-repo-config.yaml for repository settings

6. **🚨🚨🚨 R192** - Repository Separation Protocol
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R192-repo-separation.md`
   - Criticality: BLOCKING - Mixing repositories = corruption
   - Summary: Keep Software Factory and target repository completely separate

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
STATE_FILE="$CLAUDE_PROJECT_DIR/orchestrator-state-v3.json"
if [ -f "$STATE_FILE" ]; then
    CURRENT_STATE=$(jq '.state_machine.current_state' "$STATE_FILE")
    echo "✅ Found existing state: $CURRENT_STATE"

    # Extract PRD status fields (PRIORITY 1 - MANDATORY SEQUENCE)
    PRD_EXISTS=$(jq -r '.project_progression.current_project.prd_exists' "$STATE_FILE")
    PRD_PRE_EXISTS=$(jq -r '.project_progression.current_project.prd_pre_exists' "$STATE_FILE")
    PRD_VALIDATED=$(jq -r '.project_progression.current_project.prd_validated' "$STATE_FILE")

    # Extract architecture status fields (PRIORITY 2) - R550 compliant
    # Check both new (planning_files) and legacy (planning_artifacts) locations
    PROJECT_ARCH=$(jq -r '.planning_files.project.architecture_plan // .planning_artifacts.master_architecture_file // "null"' "$STATE_FILE")

    # Extract phase planning status fields (PRIORITY 3)
    PHASE_ARCH=$(jq -r '.project_progression.current_phase.architecture_file' "$STATE_FILE")

    echo "📊 PRD Status: exists=$PRD_EXISTS, pre_exists=$PRD_PRE_EXISTS, validated=$PRD_VALIDATED"
    echo "📊 Architecture Status: project=$PROJECT_ARCH, phase=$PHASE_ARCH (R550 compliant)"

    # Determine next state following mandatory_sequences.project_initialization
    if [ "$PRD_EXISTS" == "false" ] && [ "$PRD_PRE_EXISTS" == "false" ]; then
        PROPOSED_NEXT_STATE="SPAWN_PRODUCT_MANAGER_PRD_CREATION"
        echo "🔴 PRIORITY 1: No PRD exists - must create PRD first"
    elif [ "$PROJECT_ARCH" == "null" ] || [ ! -f "$PROJECT_ARCH" ]; then
        PROPOSED_NEXT_STATE="SPAWN_ARCHITECT_MASTER_PLANNING"
        echo "🔴 PRIORITY 2: PRD exists, need project architecture (R550)"
    elif [ "$PHASE_ARCH" == "null" ] || [ ! -f "$PHASE_ARCH" ]; then
        PROPOSED_NEXT_STATE="SPAWN_ARCHITECT_PHASE_PLANNING"
        echo "🔴 PRIORITY 3: Master arch exists, need phase planning"
    else
        PROPOSED_NEXT_STATE="WAVE_START"
        echo "✅ All planning complete - ready for wave execution"
    fi

    echo "→ Proposed next state: $PROPOSED_NEXT_STATE"
else
    echo "🔴🔴🔴 R281: Creating COMPLETE initial state file..."
    # Parse PROJECT-IMPLEMENTATION-PLAN.md
    # Create state with ALL phases/waves/efforts
    # Use templates/initial-state-template.yaml
    # Validate with utilities/validate-state-completeness.sh
    # MUST include PRD fields: prd_exists, prd_pre_exists, prd_validated
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
# Transition to PROPOSED_NEXT_STATE determined in STEP 1
```

## State Context
Initial state when orchestrator starts. Primary purpose is to:
1. Load or create orchestrator-state-v3.json
2. Verify environment configuration
3. Transition to appropriate next state

## Primary Actions
1. Check for compaction and recover if needed
2. Load target-repo-config.yaml (R191 - INIT-specific)
3. Verify repository separation (R192 - INIT-specific)  
4. Check/create orchestrator-state-v3.json (R281 from orchestrator.md)
5. If creating new state file:
   - Parse PROJECT-IMPLEMENTATION-PLAN.md completely
   - Extract ALL phases, waves, and efforts
   - Create COMPLETE state file with every item
   - Validate completeness
6. Determine if resuming or starting fresh
7. Transition to appropriate state immediately

## State Transitions

### VALID TRANSITIONS FROM INIT (from state-machines/software-factory-3.0-state-machine.json):
Following the mandatory_sequences.project_initialization order:
1. **SPAWN_PRODUCT_MANAGER_PRD_CREATION** - When no PRD exists (new project without pre-existing PRD)
2. **SPAWN_ARCHITECT_MASTER_PLANNING** - When PRD exists/validated but no master plan exists
3. **SPAWN_ARCHITECT_PHASE_PLANNING** - When master plan exists but phase plan doesn't
4. **WAVE_START** - When resuming existing wave work or all planning complete

### Decision Logic (Following mandatory_sequences.project_initialization):

**STEP 1: Check if resuming existing state**
- If orchestrator-state-v3.json exists with current_state → Resume from that state
  - Usually WAVE_START or appropriate continuation state

**STEP 2: If creating new state file (R281), determine first state needed:**

**Priority 1: Check PRD (FIRST IN MANDATORY SEQUENCE)**
- Extract from state file:
  - `prd_exists` (boolean - PRD file created by factory)
  - `prd_pre_exists` (boolean - PRD existed before factory initialization)
  - `prd_validated` (boolean - PRD passed validation)

- **Decision logic:**
  - If `prd_exists == false` AND `prd_pre_exists == false`:
    → **SPAWN_PRODUCT_MANAGER_PRD_CREATION** (must create PRD from scratch)
  - If `prd_exists == true` OR `prd_pre_exists == true`:
    → Skip to Priority 2 (PRD exists, check architecture)

**Priority 2: Check Master Architecture (SECOND IN MANDATORY SEQUENCE)**
- Extract from state file: `master_architecture_file`
- **Decision logic:**
  - If master_architecture_file is null or file doesn't exist:
    → **SPAWN_ARCHITECT_MASTER_PLANNING** (need master architecture)
  - If master_architecture_file exists:
    → Skip to Priority 3 (check phase planning)

**Priority 3: Check Phase Planning (THIRD IN MANDATORY SEQUENCE)**
- Extract from state file: current phase's architecture_file
- **Decision logic:**
  - If no phase plan exists:
    → **SPAWN_ARCHITECT_PHASE_PLANNING** (need phase architecture)
  - If phase plan exists:
    → **WAVE_START** (all planning complete, start execution)

**STEP 3: Error conditions**
- If state file incomplete → IMMEDIATE FAILURE (-100%)
- If config missing → ERROR_RECOVERY

## Critical Enforcement Order

**From orchestrator.md (already read):**
- R283, R288, R232, R322, R322, R281, R280, software-factory-3.0-state-machine.json
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



## 🔴🔴🔴 STATE COMPLETION CHECKLIST 🔴🔴🔴

**Execute these steps IN ORDER to properly complete INIT:**

### ✅ Step 1: Complete State-Specific Work
**Refer to "Primary Actions" and "Mandatory Validations" sections above for state-specific tasks.**

Once all state work is complete, proceed to mandatory exit protocol:

---

### ✅ Step 2: Set Proposed Next State and Transition Reason
```bash
# Based on state work results, determine proposed next state
PROPOSED_NEXT_STATE="[DETERMINE_FROM_STATE_LOGIC]"
TRANSITION_REASON="INIT complete - [describe what was accomplished]"
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
  "state_completed": "INIT",
  "work_accomplished": [
    "Loaded target-repo-config.yaml",
    "Verified repository separation",
    "Created/loaded orchestrator-state-v3.json",
    "Determined next state"
  ],
  "proposed_next_state": "$PROPOSED_NEXT_STATE",
  "transition_reason": "$TRANSITION_REASON"
}
EOF
)

# Spawn State Manager
/spawn state-manager SHUTDOWN_CONSULTATION \
  --current-state "INIT" \
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
save_todos "INIT_COMPLETE"

# Commit TODOs within 60 seconds (R287)
cd "$CLAUDE_PROJECT_DIR"
git add todos/*.todo

if ! git commit -m "todo: orchestrator - INIT complete [R287]"; then
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


### 🚨 CRITICAL DISTINCTION: AGENT STOPS ≠ FACTORY STOPS 🚨

**TWO INDEPENDENT DECISIONS - DO NOT CONFUSE THEM:**

#### 1. Should Agent Stop Work? (R322 Technical Requirement)
- Agent completes current state
- Agent saves TODOs and commits state
- Agent exits with `exit 0` (preserves context)
- User runs /continue-orchestrating to resume
- **This is NORMAL at checkpoints**

#### 2. Should Factory Continue? (R405 Operational Status)
- Even though agent stopped, can automation proceed?
- TRUE = Healthy completion, automation can continue
- FALSE = Catastrophic failure, must halt everything
- **R322 checkpoints = TRUE (99.9% of cases)**

### THE PATTERN AT R322 CHECKPOINTS (SF 3.0)

```bash
# 1. Complete state work
echo "✅ State work complete"

# 2. Set proposed next state
PROPOSED_NEXT_STATE="NEXT_STATE"
TRANSITION_REASON="State work complete"

# 3. Spawn State Manager for state transition
/spawn state-manager SHUTDOWN_CONSULTATION \
  --current-state "CURRENT_STATE" \
  --proposed-next-state "$PROPOSED_NEXT_STATE" \
  --transition-reason "$TRANSITION_REASON"
# State Manager updates all 4 state files atomically

# 4. Save TODOs
save_todos "R322_CHECKPOINT"

# 5. Factory continues (operational status)
echo "CONTINUE-SOFTWARE-FACTORY=TRUE REASON=STATE_COMPLETE"

# 6. Agent stops (technical requirement)
exit 0
```

**Both happen together! Agent stops AND factory continues!**

### WHEN TO USE EACH FLAG VALUE

**TRUE (99.9%):**
- ✅ R322 checkpoint reached
- ✅ State work completed successfully
- ✅ Ready for /continue-orchestrating
- ✅ Waiting for user to continue (NORMAL)
- ✅ Plan ready for review (agent done, factory proceeds)

**FALSE (0.1%):**
- ❌ CATASTROPHIC unrecoverable error
- ❌ Data corruption spreading
- ❌ Critical security violation
- ❌ NOT for R322 checkpoints
- ❌ NOT for user review needs
**See: $CLAUDE_PROJECT_DIR/rule-library/R405-automation-continuation-flag.md**

