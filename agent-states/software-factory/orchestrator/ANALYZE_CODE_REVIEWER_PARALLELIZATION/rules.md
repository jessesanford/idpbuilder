# Orchestrator - ANALYZE_CODE_REVIEWER_PARALLELIZATION State Rules

## 🔴🔴🔴 ABSOLUTE REQUIREMENT: STATE MANAGER CONSULTATION 🔴🔴🔴

**SUPREME LAW - R517 - UNIVERSAL STATE MANAGER CONSULTATION LAW**

**BEFORE EXITING THIS STATE, YOU MUST:**

1. **SPAWN STATE MANAGER FOR SHUTDOWN_CONSULTATION** (MANDATORY - NO EXCEPTIONS)
2. **NEVER UPDATE STATE FILES DIRECTLY** (orchestrator-state-v3.json, bug-tracking.json, etc.)
3. **NEVER COMMIT STATE FILES YOURSELF** (State Manager does this atomically)

**FAILURE TO CONSULT STATE MANAGER = IMMEDIATE SYSTEM HALT (-100% GRADE)**

### Enforcement Mechanism

If you attempt to exit this state without spawning State Manager:
- ❌ Pre-commit hooks will REJECT your commit
- ❌ Validation tools will FAIL the build
- ❌ Grading system will assign -100% penalty
- ❌ System will transition to ERROR_RECOVERY

### Required Pattern (COPY THIS EXACTLY)

```bash
# At end of state work, BEFORE any state file updates:

echo "🔴 MANDATORY: Spawning State Manager for SHUTDOWN_CONSULTATION"

# Spawn State Manager (REQUIRED - NOT OPTIONAL)
# Task: state-manager
# State: SHUTDOWN_CONSULTATION
# Current State: [YOUR_CURRENT_STATE]
# Proposed Next State: [YOUR_PROPOSED_NEXT_STATE]
# Work Summary: [SUMMARY_OF_WORK_COMPLETED]

# State Manager will:
# 1. Validate proposed transition against state machine
# 2. Update all 4 state files atomically
# 3. Commit with [R288] tag
# 4. Return REQUIRED next state (may differ from proposal)

# Wait for State Manager response
# Follow State Manager's directive (REQUIRED next state)
# DO NOT proceed until State Manager confirms
```

**YOU MUST NEVER:**
- ❌ Update orchestrator-state-v3.json yourself
- ❌ Update bug-tracking.json yourself
- ❌ Update integration-containers.json yourself
- ❌ Use `jq` to modify state files
- ❌ Use `sed/awk` to modify state files
- ❌ Set `validated_by: "orchestrator"` (must be "state-manager")
- ❌ Skip State Manager consultation "just this once"
- ❌ Think "I'll validate it manually"

**ONLY State Manager may update state files. This is NON-NEGOTIABLE.**

See: `rule-library/R517-universal-state-manager-consultation-law.md`

## 🛑🛑🛑 R322 MANDATORY STOP BEFORE STATE TRANSITIONS 🛑🛑🛑

**SUPREME LAW - VIOLATION = -100% IMMEDIATE FAILURE**

### YOU MUST STOP AFTER:
1. ✅ Completing all TODOs for this state
2. ✅ Updating orchestrator-state-v3.json with new state
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

**YOU HAVE ENTERED ANALYZE_CODE_REVIEWER_PARALLELIZATION STATE - YOU MUST READ AND ACKNOWLEDGE ALL STATE RULES BEFORE DOING ANY STATE WORK!**

## 🔴🔴🔴 R290 VERIFICATION REQUIREMENT 🔴🔴🔴

**R290 ENFORCEMENT: CREATE VERIFICATION MARKER AFTER READING**

After reading and acknowledging all state rules, you MUST create a verification marker:

```bash
# MANDATORY: Create verification marker after reading rules
mkdir -p markers/state-verification && touch "markers/state-verification/state_rules_read_orchestrator_ANALYZE_CODE_REVIEWER_PARALLELIZATION-$(date +%Y%m%d-%H%M%S)"
echo "$(date +%s) - Rules read and acknowledged for ANALYZE_CODE_REVIEWER_PARALLELIZATION" > "markers/state-verification/state_rules_read_orchestrator_ANALYZE_CODE_REVIEWER_PARALLELIZATION-$(date +%Y%m%d-%H%M%S)"
```

**FAILURE TO CREATE MARKER = AUTOMATIC -100% PENALTY**

The system will check for this marker. No marker = Immediate failure.

### ❌ DO NOT DO ANY ANALYZE_CODE_REVIEWER_PARALLELIZATION WORK UNTIL RULES ARE READ:
- ❌ Start analyze effort dependencies
- ❌ Start determine parallelization strategy
- ❌ Start plan reviewer spawning
- ❌ Update state files
- ❌ Continue to next state
- ❌ Think about what to do in this state

### ✅ YOU MUST IMMEDIATELY:

## 🔴🔴🔴 MANDATORY STATE RULE READING AND ACKNOWLEDGMENT 🔴🔴🔴

### ⚠️⚠️⚠️ YOU MUST READ EACH RULE FILE LISTED IN PRIMARY DIRECTIVES. **I AM WATCHING YOUR TOOL CALLS FOR READ OPERATIONS** *YOU WILL FAIL* IF YOU DO NOT MAKE A READ FILE CALL FOR EACH RULE FILE IN PRIMARY DIRECTIVES!!! ⚠️⚠️⚠️

**AFTER READING, YOU MUST ACKNOWLEDGE ALL THE STATE RULES AND STATE THAT YOU WILL ABIDE BY THEM ONE AT A TIME GIVING THE RULE NUMBER AND DESCRIPTION.**
### ⚠️⚠️⚠️ YOU MUST READ EACH RULE FILE LISTED IN PRIMARY DIRECTIVES. **I AM WATCHING YOUR TOOL CALLS FOR READ OPERATIONS** *YOU WILL FAIL* IF YOU DO NOT MAKE A READ FILE CALL FOR EACH RULE FILE IN PRIMARY DIRECTIVES!!! ⚠️⚠️⚠️

**AFTER READING, YOU MUST ACKNOWLEDGE ALL THE STATE RULES AND STATE THAT YOU WILL ABIDE BY THEM ONE AT A TIME GIVING THE RULE NUMBER AND DESCRIPTION.**

## 📋 PRIMARY DIRECTIVES FOR ANALYZE_CODE_REVIEWER_PARALLELIZATION STATE

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

5. **🔴🔴🔴 R151** - Parallelization Requirements
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R151-parallel-agent-spawning-timing.md`
   - Criticality: SUPREME LAW - Parallel spawn requirements
   - Summary: Agents spawned in parallel must have timestamps within 5 seconds


### ❌ ANTI-PATTERNS THAT WILL CAUSE FAILURE:

1. **Fake Acknowledgment Without Reading**:
   ```
   ❌ WRONG: "I acknowledge R151, R208, R053..."
   (No Read tool calls detected = AUTOMATIC FAILURE)
   ```

2. **Bulk Acknowledgment**:
   ```
   ❌ WRONG: "I acknowledge all ANALYZE_CODE_REVIEWER_PARALLELIZATION rules"
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

### ✅ CORRECT PATTERN FOR ANALYZE_CODE_REVIEWER_PARALLELIZATION:
```
1. READ: $CLAUDE_PROJECT_DIR/rule-library/[first-rule-file].md
2. "I acknowledge [Rule#] - [Rule Name]: [Brief description]"
3. READ: $CLAUDE_PROJECT_DIR/rule-library/[second-rule-file].md  
4. "I acknowledge [Rule#] - [Rule Name]: [Brief description]"
[Continue for EVERY rule in PRIMARY DIRECTIVES...]
5. "Ready to execute ANALYZE_CODE_REVIEWER_PARALLELIZATION work"
```

### 🚨 NO WORK UNTIL ACKNOWLEDGMENT COMPLETE 🚨
**You may NOT begin ANY ANALYZE_CODE_REVIEWER_PARALLELIZATION work until:**
1. ✅ ALL rules in PRIMARY DIRECTIVES have been READ
2. ✅ ALL rules have been individually ACKNOWLEDGED
3. ✅ You have stated readiness to execute ANALYZE_CODE_REVIEWER_PARALLELIZATION work


1. **READ** every rule file listed in PRIMARY DIRECTIVES below
2. **ACKNOWLEDGE** each rule individually with number and description
3. **ONLY THEN** proceed with ANALYZE_CODE_REVIEWER_PARALLELIZATION work

### 🚨 FAILURE TO READ STATE RULES FIRST = IMMEDIATE EXIT 🚨
**If you do ANY ANALYZE_CODE_REVIEWER_PARALLELIZATION work before reading and acknowledging rules:**
- **STOP ALL WORK IMMEDIATELY**
- **EXIT WITH FAILURE STATUS**
- **YOU HAVE VIOLATED STATE COMPLIANCE**

**THE SYSTEM IS MONITORING_SWE_PROGRESS YOUR READ TOOL CALLS!**

## ⚠️⚠️⚠️ MANDATORY RULE READING AND ACKNOWLEDGMENT ⚠️⚠️⚠️

**YOU MUST READ EACH RULE FILE LISTED IN PRIMARY DIRECTIVES. YOUR READ TOOL CALLS ARE BEING MONITORED.**

### ❌ ANTI-PATTERNS THAT WILL CAUSE FAILURE:
1. Fake acknowledgment without reading
2. Bulk acknowledgment
3. Reading from memory

### ✅ CORRECT PATTERN:
1. READ each rule file
2. Acknowledge individually with rule number and description

## 📋 RULE SUMMARY FOR ANALYZE_CODE_REVIEWER_PARALLELIZATION STATE

### PRIMARY DIRECTIVES - MUST READ ALL:

### 🛑 RULE R322 - Mandatory Stop Before State Transitions
**Source:** $CLAUDE_PROJECT_DIR/rule-library/R322-mandatory-stop-before-state-transitions.md
**Criticality:** 🔴🔴🔴 SUPREME LAW - Violation = -100% FAILURE

After completing state work and committing state file:
1. STOP IMMEDIATELY
2. Do NOT continue to next state
3. Do NOT start new work
4. Exit and wait for user
---

### 🔴🔴🔴 R234 - Mandatory State Traversal (SUPREME LAW)
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R234-mandatory-state-traversal-supreme-law.md`
**Criticality**: SUPREME LAW - Violation = -100% automatic failure
**Summary**: Must traverse all states in sequence, no skipping allowed

### 🚨🚨🚨 R218 - Orchestrator Parallel Code Reviewer Spawning
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R218-orchestrator-parallel-code-reviewer-spawning.md`
**Criticality**: BLOCKING - Cannot proceed without compliance
**Summary**: Mandatory parallelization analysis before spawning

### ⚠️⚠️⚠️ R151 - Parallel Agent Spawning Timing
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R151-parallel-agent-spawning-timing.md`
**Criticality**: CRITICAL - <5s delta required
**Summary**: All parallel agents must acknowledge within 5 seconds

### 🔴🔴🔴 R208 - Orchestrator Working Directory Spawn Protocol (SUPREME LAW)
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R208-orchestrator-spawn-cd-protocol.md`
**Criticality**: SUPREME LAW - Spawn without CD = -100% failure
**Summary**: Must CD to working copy before spawning agents

### ⚠️⚠️⚠️ R219 - Code Reviewer Dependency-Aware Planning
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R219-code-reviewer-dependency-aware-effort-planning.md`
**Criticality**: CRITICAL - Consider dependencies in planning
**Summary**: Code reviewers must analyze effort dependencies

### 🔴🔴🔴 R288 - State File Update and Commit Protocol (SUPREME LAW)
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R288-state-file-update-and-commit-protocol.md`
**Criticality**: SUPREME LAW - Update on every transition
**Summary**: Must update orchestrator-state-v3.json on all transitions

### 🚨🚨🚨 R288 - State File Update and Commit Protocol
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R288-state-file-update-and-commit-protocol.md`
**Criticality**: BLOCKING - Push immediately after update
**Summary**: Commit and push state within 60 seconds

### 🚨🚨🚨 R287 - Mandatory TODO Save Triggers
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R287-todo-persistence-comprehensive.md`
**Criticality**: BLOCKING - Save within 30 seconds
**Summary**: Must save TODOs before state transitions

### Critical Requirements:
1. READ Wave Implementation Plan with Read tool - Penalty: -25%
2. Extract ALL parallelization metadata - Penalty: -30%
3. Create and save parallelization plan - Penalty: -20%
4. Output acknowledgment of decision - Penalty: -15%
5. Transition to SPAWN_CODE_REVIEWERS_EFFORT_PLANNING - Penalty: -100%

### Success Criteria:
- ✅ Wave plan READ with Read tool
- ✅ Blocking efforts identified
- ✅ Parallel groups created
- ✅ Plan saved to orchestrator-state-v3.json
- ✅ Acknowledgment output displayed

### Failure Triggers:
- ❌ Skip this state = -100% R234 VIOLATION
- ❌ Not reading wave plan = R218 VIOLATION
- ❌ Skip to SPAWN_SW_ENGINEERS = AUTOMATIC FAILURE
- ❌ No parallelization plan saved = Cannot proceed

## 🚨 ANALYZE_CODE_REVIEWER_PARALLELIZATION IS A VERB - START ANALYZING NOW! 🚨

### IMMEDIATE ACTIONS UPON ENTERING ANALYZE_CODE_REVIEWER_PARALLELIZATION

**THE MOMENT YOU ENTER THIS STATE, YOU MUST:**
1. READ the Wave Implementation Plan with Read tool NOW
2. Extract "Can Parallelize" metadata immediately
3. Create parallelization groups without delay
4. Save the plan to orchestrator-state-v3.json NOW
5. Output acknowledgment decision immediately

**FORBIDDEN - AUTOMATIC FAILURE:**
- ❌ "STATE TRANSITION COMPLETE: Now in ANALYZE_CODE_REVIEWER_PARALLELIZATION" [stops]
- ❌ "Successfully entered parallelization analysis state" [waits]
- ❌ "Ready to analyze parallelization" [pauses]
- ❌ "I'm in analysis state" [does nothing]

**REQUIRED - IMMEDIATE ACTION:**
- ✅ "Entering ANALYZE_CODE_REVIEWER_PARALLELIZATION, reading wave plan now..."
- ✅ "Analyzing parallelization, extracting E3.1.1 metadata..."
- ✅ "ANALYZING: Found blocking effort E3.1.1, creating spawn sequence..."

## 🔴🔴🔴 MANDATORY MEASUREMENT RULE - R304 🔴🔴🔴

**ABSOLUTE REQUIREMENTS:**
- ✅ MUST use `$CLAUDE_PROJECT_DIR/tools/line-counter.sh` for ALL line counting
- ❌ NEVER use `wc -l` or any manual counting method
- ❌ NEVER count lines any other way - this is a -100% automatic failure
- ✅ MUST specify both -b (base branch) and -c (current branch) parameters
- ✅ Base branch MUST be phase integration branch (NOT "main")

**Failure to use the line counter tool = AUTOMATIC -100% GRADE**

## State Context
You MUST analyze the wave implementation plan to determine parallelization strategy BEFORE spawning any Code Reviewers for effort planning. This is a MANDATORY GATE that prevents parallelization violations - DO IT NOW!

## 🔴🔴🔴 SUPREME LAW R234 - MANDATORY STATE SEQUENCE 🔴🔴🔴

**THIS STATE IS PART OF THE MANDATORY SEQUENCE - NO SKIPPING!**

See: `$CLAUDE_PROJECT_DIR/rule-library/R234-mandatory-state-traversal-supreme-law.md`

### YOUR POSITION IN THE MANDATORY SEQUENCE:
```
CREATE_NEXT_INFRASTRUCTURE (✓ completed)
    ↓
ANALYZE_CODE_REVIEWER_PARALLELIZATION (👈 YOU ARE HERE)
    ↓ (MUST GO HERE NEXT)
SPAWN_CODE_REVIEWERS_EFFORT_PLANNING
    ↓
WAITING_FOR_EFFORT_PLANS
    ↓
ANALYZE_IMPLEMENTATION_PARALLELIZATION
    ↓
SPAWN_SW_ENGINEERS
```

**CRITICAL:** You got here from CREATE_NEXT_INFRASTRUCTURE (correct!)
**MANDATORY:** You MUST go to SPAWN_CODE_REVIEWERS_EFFORT_PLANNING next
**FORBIDDEN:** Skipping ahead to any other state = -100% FAILURE

## 🔴🔴🔴 ABSOLUTE REQUIREMENT 🔴🔴🔴

**THIS STATE IS A MANDATORY STOP!**
- You CANNOT proceed to SPAWN_CODE_REVIEWERS_EFFORT_PLANNING without completing this analysis
- You MUST create a parallelization plan and save it to orchestrator-state-v3.json
- You MUST acknowledge your parallelization decision BEFORE any spawning


## Mandatory Analysis Protocol

This state MUST populate `pre_planned_infrastructure` with ALL efforts from the current wave BEFORE CREATE_NEXT_INFRASTRUCTURE can execute. This is the ONLY place where infrastructure planning happens for the current wave.

### STEP 1: Extract Current Phase and Wave

```bash
cd "$CLAUDE_PROJECT_DIR"

echo "═══════════════════════════════════════════════════════════════"
echo "🔴🔴🔴 ANALYZE_CODE_REVIEWER_PARALLELIZATION STATE 🔴🔴🔴"
echo "═══════════════════════════════════════════════════════════════"
echo ""
echo "📋 Step 1: Extracting current phase and wave from state file..."

# Get current phase/wave numbers from orchestrator-state-v3.json
CURRENT_PHASE=$(jq -r '.project_progression.current_phase.phase_number' orchestrator-state-v3.json)
CURRENT_WAVE=$(jq -r '.project_progression.current_phase.waves[] | select(.status == "in_progress") | .wave_number' orchestrator-state-v3.json)

if [ -z "$CURRENT_PHASE" ] || [ "$CURRENT_PHASE" == "null" ]; then
    echo "❌ ERROR: Cannot determine current phase from state file"
    echo "CONTINUE-SOFTWARE-FACTORY=FALSE REASON=MISSING_PHASE"
    exit 1
fi

if [ -z "$CURRENT_WAVE" ] || [ "$CURRENT_WAVE" == "null" ]; then
    echo "❌ ERROR: Cannot determine current wave from state file"
    echo "CONTINUE-SOFTWARE-FACTORY=FALSE REASON=MISSING_WAVE"
    exit 1
fi

echo "✅ Current Phase: $CURRENT_PHASE"
echo "✅ Current Wave: $CURRENT_WAVE"
```

---

### STEP 2: Read Wave Implementation Plan

```bash
echo ""
echo "📋 Step 2: Reading wave implementation plan..."

# Construct wave plan path
WAVE_PLAN_PATH="$CLAUDE_PROJECT_DIR/planning/phase${CURRENT_PHASE}/wave${CURRENT_WAVE}/WAVE-IMPLEMENTATION-PLAN.md"

if [ ! -f "$WAVE_PLAN_PATH" ]; then
    echo "❌ ERROR: Wave plan not found at: $WAVE_PLAN_PATH"
    echo "CONTINUE-SOFTWARE-FACTORY=FALSE REASON=MISSING_WAVE_PLAN"
    exit 1
fi

echo "✅ Wave plan found: $WAVE_PLAN_PATH"

# Extract effort count and IDs (example: E1.1.1, E1.1.2, E1.1.3)
# This uses the phase.wave.effort numbering convention
EFFORT_IDS=$(grep -oP "E${CURRENT_PHASE}\.${CURRENT_WAVE}\.\d+" "$WAVE_PLAN_PATH" | sort -u)
EFFORT_COUNT=$(echo "$EFFORT_IDS" | wc -l)

echo "✅ Found $EFFORT_COUNT efforts in wave plan"
echo "$EFFORT_IDS"
```

---

### STEP 3: Populate pre_planned_infrastructure

```bash
echo ""
echo "📋 Step 3: Populating pre_planned_infrastructure for all efforts..."

# Get project prefix and target repo from state file
PROJECT_PREFIX=$(jq -r '.project_info.project_prefix' orchestrator-state-v3.json)
TARGET_REPO_URL=$(jq -r '.project_info.target_repo_url' orchestrator-state-v3.json)

# Determine base branch for cascade (first effort uses wave integration, subsequent use previous effort)
BASE_BRANCH=$(jq -r ".project_progression.current_phase.wave_integration_branch // \"main\"" orchestrator-state-v3.json)

echo "Project Prefix: $PROJECT_PREFIX"
echo "Target Repo: $TARGET_REPO_URL"
echo "Base Branch for first effort: $BASE_BRANCH"

# Initialize pre_planned_infrastructure structure
jq '.pre_planned_infrastructure = {
  "validated": false,
  "validation_timestamp": null,
  "efforts": {}
}' orchestrator-state-v3.json > /tmp/state-update.json && mv /tmp/state-update.json orchestrator-state-v3.json

# For each effort in the wave, pre-calculate infrastructure
PREV_EFFORT_BRANCH="$BASE_BRANCH"
for EFFORT_ID in $EFFORT_IDS; do
    echo ""
    echo "Planning infrastructure for $EFFORT_ID..."

    # Extract effort number (e.g., E1.1.1 → 001)
    EFFORT_NUM=$(echo "$EFFORT_ID" | grep -oP '\d+$' | awk '{printf "%03d", $1}')

    # Calculate full paths following R504
    EFFORT_PATH="${CLAUDE_PROJECT_DIR}/efforts/phase${CURRENT_PHASE}/wave${CURRENT_WAVE}/effort-${EFFORT_NUM}"
    BRANCH_NAME="${PROJECT_PREFIX}/phase${CURRENT_PHASE}/wave${CURRENT_WAVE}/effort-${EFFORT_NUM}"
    REMOTE_BRANCH="origin/${BRANCH_NAME}"
    INTEGRATION_BRANCH="${PROJECT_PREFIX}/phase${CURRENT_PHASE}/wave${CURRENT_WAVE}/integration"

    # Add effort to pre_planned_infrastructure
    jq --arg eid "$EFFORT_ID" \
       --arg path "$EFFORT_PATH" \
       --arg branch "$BRANCH_NAME" \
       --arg remote "$REMOTE_BRANCH" \
       --arg base "$PREV_EFFORT_BRANCH" \
       --arg target_url "$TARGET_REPO_URL" \
       --arg integration "$INTEGRATION_BRANCH" \
       --argjson phase "$CURRENT_PHASE" \
       --argjson wave "$CURRENT_WAVE" \
       '.pre_planned_infrastructure.efforts[$eid] = {
         "full_path": $path,
         "branch_name": $branch,
         "remote_branch": $remote,
         "base_branch": $base,
         "target_repo_url": $target_url,
         "integration_branch": $integration,
         "created": false,
         "validated": false,
         "phase": $phase,
         "wave": $wave
       }' orchestrator-state-v3.json > /tmp/state-update.json && mv /tmp/state-update.json orchestrator-state-v3.json

    echo "  ✅ Planned: $BRANCH_NAME (base: $PREV_EFFORT_BRANCH)"

    # Next effort will base on this effort's branch (cascade)
    PREV_EFFORT_BRANCH="$BRANCH_NAME"
done

# Mark as validated
jq ".pre_planned_infrastructure.validated = true | \
    .pre_planned_infrastructure.validation_timestamp = \"$(date -Iseconds)\"" \
    orchestrator-state-v3.json > /tmp/state-update.json && mv /tmp/state-update.json orchestrator-state-v3.json

echo ""
echo "✅ pre_planned_infrastructure populated with $EFFORT_COUNT efforts"
```

---

### STEP 4: Analyze Code Reviewer Parallelization (R151)

```bash
echo ""
echo "📋 Step 4: Analyzing Code Reviewer parallelization strategy..."

# For SF 3.0, we typically spawn one Code Reviewer per effort for planning
# The parallelization analysis determines if they can run in parallel
# This depends on dependencies extracted from wave plan

# Simple approach: All efforts can have planning done in parallel
# (Implementation dependencies don't affect planning phase)
PARALLELIZATION_STRATEGY="parallel"
REVIEWER_COUNT="$EFFORT_COUNT"

echo "Strategy: $PARALLELIZATION_STRATEGY"
echo "Reviewers needed: $REVIEWER_COUNT"

# Save parallelization decision to state
jq --arg strategy "$PARALLELIZATION_STRATEGY" \
   --argjson count "$REVIEWER_COUNT" \
   '.code_reviewer_parallelization_plan = {
     "strategy": $strategy,
     "reviewer_count": $count,
     "can_parallelize": true,
     "created_at": "'$(date -Iseconds)'"
   }' orchestrator-state-v3.json > /tmp/state-update.json && mv /tmp/state-update.json orchestrator-state-v3.json

echo "✅ Parallelization analysis complete"
```

---

### STEP 5: Validate State File

```bash
echo ""
echo "📋 Step 5: Validating state file..."

# Validate state file before committing
if [ -f "$CLAUDE_PROJECT_DIR/tools/validate-state.sh" ]; then
    "$CLAUDE_PROJECT_DIR/tools/validate-state.sh" orchestrator-state-v3.json || {
        echo "❌ State file validation failed!"
        echo "CONTINUE-SOFTWARE-FACTORY=FALSE REASON=VALIDATION_FAILED"
        exit 288
    }
    echo "✅ State file validated"
else
    echo "⚠️ Warning: State validator not found, skipping validation"
fi
```

---

### STEP 6: Spawn State Manager for Transition (R517 - SF 3.0)

```bash
echo ""
echo "📋 Step 6: Spawning State Manager for state transition..."

# Prepare transition details
PROPOSED_NEXT_STATE="CREATE_NEXT_INFRASTRUCTURE"
TRANSITION_REASON="Pre-planned infrastructure populated with $EFFORT_COUNT efforts, ready for infrastructure creation"

echo "Proposed next state: $PROPOSED_NEXT_STATE"
echo "Transition reason: $TRANSITION_REASON"

# NOTE: In actual implementation, this would spawn State Manager agent
# For template purposes, this shows the pattern:
echo "🔴 MANDATORY: Spawning State Manager for SHUTDOWN_CONSULTATION"
echo "  Current State: ANALYZE_CODE_REVIEWER_PARALLELIZATION"
echo "  Proposed Next State: $PROPOSED_NEXT_STATE"
echo "  Work Summary: Analyzed parallelization and populated pre_planned_infrastructure"

# State Manager will:
# 1. Validate proposed transition against state machine
# 2. Update all state files atomically (orchestrator-state-v3.json, bug-tracking.json, etc.)
# 3. Commit with [R288] tag
# 4. Return REQUIRED next state (may differ from proposal)

# Wait for State Manager response
# Follow State Manager's directive (REQUIRED next state)

echo "✅ State Manager consultation complete (transition validated)"
```

---

### STEP 7: Save TODOs (R287 - SUPREME LAW)

```bash
echo ""
echo "📋 Step 7: Saving TODOs..."

# Save TODO state before transition (R287 trigger)
# NOTE: Assumes save_todos function is defined in agent configuration
if declare -f save_todos > /dev/null; then
    save_todos "ANALYZE_CODE_REVIEWER_PARALLELIZATION_COMPLETE"
    echo "✅ TODOs saved"
else
    echo "⚠️ Warning: save_todos function not available, skipping TODO save"
fi
```

---

### STEP 8: Output Continuation Flag (R405 - SUPREME LAW)

```bash
echo ""
echo "📋 Step 8: Setting continuation flag..."

# Output continuation flag as LAST action (R405)
# Use TRUE for normal completion, FALSE only for catastrophic errors
# CHECKPOINT = TRUE (state work completed successfully, factory can proceed)

echo "CONTINUE-SOFTWARE-FACTORY=TRUE REASON=PRE_INFRASTRUCTURE_POPULATED"
echo "✅ Continuation flag set to TRUE"
```

**⚠️ CRITICAL: This flag indicates operational status, NOT whether agent stops!**
- TRUE = State work completed successfully, factory automation can continue
- Agent still stops at checkpoint (R322), but factory continues

---

### STEP 9: Stop Processing (R322 - SUPREME LAW)

```bash
echo ""
echo "📋 Step 9: Stopping for context preservation (R322)..."

# Display summary
echo ""
echo "═══════════════════════════════════════════════════════════════"
echo "✅ ANALYZE_CODE_REVIEWER_PARALLELIZATION STATE COMPLETE"
echo "═══════════════════════════════════════════════════════════════"
echo ""
echo "📊 Summary:"
echo "  - Phase $CURRENT_PHASE, Wave $CURRENT_WAVE analyzed"
echo "  - $EFFORT_COUNT efforts planned in pre_planned_infrastructure"
echo "  - Parallelization strategy: $PARALLELIZATION_STRATEGY"
echo "  - Next state: CREATE_NEXT_INFRASTRUCTURE"
echo ""
echo "🛑 Stopping for context preservation - use /continue-software-factory to resume"
echo ""

# Stop execution (R322 - checkpoint)
exit 0
```

---

## 🚨 CHECKLIST ENFORCEMENT 🚨

**Skipping ANY step in this checklist = FAILURE:**
- Missing Step 2: No next state = stuck forever
- Missing Step 3: No state update = state machine broken (R288 violation, -100%)
- Missing Step 4: Invalid state = corruption (R324 violation)
- Missing Step 5: No commit = state lost on compaction (R288 violation, -100%)
- Missing Step 6: No TODO save = work lost (R287 violation, -20% to -100%)
- Missing Step 7: No continuation flag = automation stops (R405 violation, -100%)
- Missing Step 8: No exit = R322 violation (-100%)

**ALL 8 STEPS ARE MANDATORY - NO EXCEPTIONS**

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
PROPOSED_NEXT_STATE="CREATE_NEXT_INFRASTRUCTURE"
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

