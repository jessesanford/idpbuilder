# Orchestrator - SPAWN_SW_ENGINEERS State Rules

## 🔴🔴🔴 CRITICAL: STATE FILES PROPOSE, STATE MANAGER DECIDES 🔴🔴🔴

**R288 SUPREME LAW - YOU NEVER CALL update_state DIRECTLY!**

This state file uses the **BOOKEND PATTERN** (R600):
1. **START**: Set `PROPOSED_NEXT_STATE` and `TRANSITION_REASON` variables
2. **WORK**: Execute state-specific logic
3. **END**: Follow 10-step completion checklist to exit properly

**NEVER CALL update_state() DIRECTLY - IT IS PROHIBITED!**
- ❌ `update_state "NEXT_STATE" "reason"` = SYSTEM VIOLATION
- ✅ `PROPOSED_NEXT_STATE="NEXT_STATE"` = CORRECT
- ✅ `TRANSITION_REASON="reason"` = CORRECT

The State Manager (`run-software-factory.sh`) handles ALL state transitions.

**See:**
- `$CLAUDE_PROJECT_DIR/rule-library/R288-state-transition-authority.md` (SUPREME LAW)
- `$CLAUDE_PROJECT_DIR/rule-library/R600-orchestrator-bookend-pattern.md`

---

# PRIMARY DIRECTIVES

You MUST read and acknowledge these rules:

1. **R006** - Orchestrator cannot write code (BLOCKING)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R006-orchestrator-never-writes-code.md`

2. **R151** - Parallel Agent Timestamp Requirement (CRITICAL)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R151-parallel-agent-spawning-timing.md`

3. **R208** - Orchestrator Spawn Directory Protocol (SUPREME LAW)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R208-orchestrator-spawn-cd-protocol.md`

4. **R502** - Mandatory Plan Validation Gates (BLOCKING)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R502-mandatory-plan-validation-gates.md`
   - **CRITICAL: Must validate phase/wave plans exist before spawning!**

5. **R287** - TODO Persistence Requirements (SUPREME LAW)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R287-todo-persistence-comprehensive.md`

6. **R288** - State File Update Requirements (SUPREME LAW)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R288-state-file-update-and-commit-protocol.md`

7. **R304** - Mandatory Line Counter Usage (SUPREME LAW)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R304-mandatory-line-counter-enforcement.md`

8. **R322** - Mandatory Stop After Spawn States (SUPREME LAW)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R322-mandatory-stop-before-state-transitions.md`

9. **R324** - State Transition Validation (SUPREME LAW)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R324-mandatory-line-counter-auto-detection.md`


## 🔴🔴🔴 STOP! STATE RULE READING IS ABSOLUTELY FIRST! 🔴🔴🔴

**YOU HAVE ENTERED SPAWN_SW_ENGINEERS STATE - YOU MUST READ AND ACKNOWLEDGE ALL STATE RULES BEFORE DOING ANY STATE WORK!**

## 🔴🔴🔴 R290 VERIFICATION REQUIREMENT 🔴🔴🔴

**R290 ENFORCEMENT: CREATE VERIFICATION MARKER AFTER READING**

After reading and acknowledging all state rules, you MUST create a verification marker:

```bash
# MANDATORY: Create verification marker after reading rules
touch .state_rules_read_orchestrator_SPAWN_SW_ENGINEERS
echo "$(date +%s) - Rules read and acknowledged for SPAWN_SW_ENGINEERS" > .state_rules_read_orchestrator_SPAWN_SW_ENGINEERS
```

**FAILURE TO CREATE MARKER = AUTOMATIC -100% PENALTY**

The system will check for this marker. No marker = Immediate failure.

### ❌ DO NOT DO ANY SPAWN_SW_ENGINEERS WORK UNTIL RULES ARE READ:
- ❌ Start spawn software engineer agents
- ❌ Start assign effort work
- ❌ Start distribute implementation tasks
- ❌ Update state files
- ❌ Continue to next state
- ❌ Think about what to do in this state

### ✅ YOU MUST IMMEDIATELY:

## 🔴🔴🔴 MANDATORY STATE RULE READING AND ACKNOWLEDGMENT 🔴🔴🔴

### ❌ ANTI-PATTERNS THAT WILL CAUSE FAILURE:

1. **Fake Acknowledgment Without Reading**:
   ```
   ❌ WRONG: "I acknowledge R151, R208, R053..."
   (No Read tool calls detected = AUTOMATIC FAILURE)
   ```

2. **Bulk Acknowledgment**:
   ```
   ❌ WRONG: "I acknowledge all SPAWN_SW_ENGINEERS rules"
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

### ✅ CORRECT PATTERN FOR SPAWN_SW_ENGINEERS:
```
1. READ: $CLAUDE_PROJECT_DIR/rule-library/[first-rule-file].md
2. "I acknowledge [Rule#] - [Rule Name]: [Brief description]"
3. READ: $CLAUDE_PROJECT_DIR/rule-library/[second-rule-file].md  
4. "I acknowledge [Rule#] - [Rule Name]: [Brief description]"
### 🚨 NO WORK UNTIL ACKNOWLEDGMENT COMPLETE 🚨
**You may NOT begin ANY SPAWN_SW_ENGINEERS work until:**
### 🚨 FAILURE TO READ STATE RULES FIRST = IMMEDIATE EXIT 🚨
**If you do ANY SPAWN_SW_ENGINEERS work before reading and acknowledging rules:**
- **STOP ALL WORK IMMEDIATELY**
- **EXIT WITH FAILURE STATUS**
- **YOU HAVE VIOLATED STATE COMPLIANCE**

**THE SYSTEM IS MONITORING_SWE_PROGRESS YOUR READ TOOL CALLS!**

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

**ABSOLUTE REQUIREMENTS:**
- ✅ MUST use `$CLAUDE_PROJECT_DIR/tools/line-counter.sh` for ALL line counting
- ❌ NEVER use `wc -l` or any manual counting method
- ❌ NEVER count lines any other way - this is a -100% automatic failure
- ✅ MUST specify both -b (base branch) and -c (current branch) parameters
- ✅ Base branch MUST be phase integration branch (NOT "main")

**Failure to use the line counter tool = AUTOMATIC -100% GRADE**

## 🚨 SPAWN_SW_ENGINEERS IS A VERB - START SPAWNING AGENTS IMMEDIATELY! 🚨

**See R151 for immediate action requirements when entering this state.**

The SPAWN_SW_ENGINEERS state requires IMMEDIATE spawning action - no pausing or waiting.
See rule R151 which includes requirements for immediate state execution.

## Helper Functions

```bash
# Generate metadata file path with R383/R343 compliance
sf_metadata_path() {
    local file_type="$1"  # IMPLEMENTATION-PLAN, CODE-REVIEW-REPORT, etc.
    local phase="$2"
    local wave="$3"
    local effort="$4"
    local timestamp="${5:-$(date +%Y%m%d-%H%M%S)}"

    echo ".software-factory/phase${phase}/wave${wave}/${effort}/${file_type}--${timestamp}.md"
}
```

## State Context
You are spawning SW Engineers to implement efforts based on the implementation plans created by Code Reviewers.

## 🔴🔴🔴 PREREQUISITES FOR SPAWN_SW_ENGINEERS 🔴🔴🔴

### 🚨🚨🚨 R502 MANDATORY PLAN VALIDATION (CRITICAL) 🚨🚨🚨

```bash
echo "🔍 R502: VALIDATING PHASE/WAVE PLANS EXIST..."

# Extract current phase and wave
CURRENT_PHASE=$(jq -r '.current_phase' orchestrator-state-v3.json)
CURRENT_WAVE=$(jq -r '.current_wave' orchestrator-state-v3.json)

# PHASE PLANS (MANDATORY)
PHASE_ARCH="$CLAUDE_PROJECT_DIR/planning/phase${CURRENT_PHASE}/PHASE-${CURRENT_PHASE}-ARCHITECTURE-PLAN.md"
PHASE_IMPL="$CLAUDE_PROJECT_DIR/planning/phase${CURRENT_PHASE}/PHASE-${CURRENT_PHASE}-IMPLEMENTATION-PLAN.md"

if [ ! -f "$PHASE_ARCH" ] || [ ! -f "$PHASE_IMPL" ]; then
    echo "🚨🚨🚨 R502 VIOLATION: Phase ${CURRENT_PHASE} plans missing!"
    [ ! -f "$PHASE_ARCH" ] && echo "❌ Missing: $PHASE_ARCH"
    [ ! -f "$PHASE_IMPL" ] && echo "❌ Missing: $PHASE_IMPL"
    echo "BLOCKING: Cannot spawn agents without phase plans!"
    PROPOSED_NEXT_STATE="ERROR_RECOVERY"
    TRANSITION_REASON="R502: Missing phase plans"
    exit 1
fi

# WAVE PLANS (MANDATORY)
WAVE_ARCH="$CLAUDE_PROJECT_DIR/planning/phase${CURRENT_PHASE}/wave${CURRENT_WAVE}/WAVE-${CURRENT_PHASE}-${CURRENT_WAVE}-ARCHITECTURE-PLAN.md"
WAVE_IMPL="$CLAUDE_PROJECT_DIR/planning/phase${CURRENT_PHASE}/wave${CURRENT_WAVE}/WAVE-${CURRENT_PHASE}-${CURRENT_WAVE}-IMPLEMENTATION-PLAN.md"

if [ ! -f "$WAVE_ARCH" ] || [ ! -f "$WAVE_IMPL" ]; then
    echo "🚨🚨🚨 R502 VIOLATION: Wave ${CURRENT_PHASE}-${CURRENT_WAVE} plans missing!"
    [ ! -f "$WAVE_ARCH" ] && echo "❌ Missing: $WAVE_ARCH"
    [ ! -f "$WAVE_IMPL" ] && echo "❌ Missing: $WAVE_IMPL"
    echo "BLOCKING: Cannot spawn agents without wave plans!"
    PROPOSED_NEXT_STATE="ERROR_RECOVERY"
    TRANSITION_REASON="R502: Missing wave plans"
    exit 1
fi

echo "✅ R502: All required plans validated!"
```

**BEFORE ENTERING THIS STATE, YOU MUST ALREADY HAVE:**
1. ✅ Phase and Wave plans validated (R502 - checked above)
2. ✅ All effort directories created (done in CREATE_NEXT_INFRASTRUCTURE)
3. ✅ All git clones and branches ready (done in CREATE_NEXT_INFRASTRUCTURE)
4. ✅ All effort IMPLEMENTATION-PLAN--*.md files (created by Code Reviewers in .software-factory/)
5. ✅ All work-log--*.log files initialized (in .software-factory/ per R383/R343)
6. ✅ **PARALLELIZATION ANALYSIS COMPLETE (ANALYZE_IMPLEMENTATION_PARALLELIZATION)**
7. ✅ **SW Engineer parallelization plan in orchestrator-state-v3.json**

**IF PARALLELIZATION NOT ANALYZED, GO BACK TO ANALYZE_IMPLEMENTATION_PARALLELIZATION!**
**Infrastructure was created BEFORE Code Reviewers made plans!**
**Now you're just spawning SW Engineers to implement using the PRE-ANALYZED strategy.**


## Parallel Spawning


## ✅ Infrastructure Already Ready

Infrastructure was set up in CREATE_NEXT_INFRASTRUCTURE state:
- Effort directories exist at: `efforts/phase{X}/wave{Y}/{effort-name}`
- Git branches created with project prefix from target-repo-config.yaml
- Remote tracking configured
- IMPLEMENTATION-PLAN--*.md files created by Code Reviewers (in .software-factory/)
- work-log--*.log files initialized (in .software-factory/)
- **SW Engineer parallelization plan in orchestrator-state-v3.json**

**Just CD to directories and spawn SW Engineers per the analyzed plan!**

### R287 TODO PERSISTENCE + R322 MANDATORY STOP
```bash
# After all SW Engineers spawned
echo "💾 R287: Saving TODOs after spawning SW Engineers..."
save_todos "SPAWN_SW_ENGINEERS complete - all SW Engineers spawned"

# R287: Commit within 60 seconds
cd $CLAUDE_PROJECT_DIR
git add todos/*.todo orchestrator-state-v3.json
git commit -m "todo: SW Engineers spawned, stopping per R322"
git push

# 🔴🔴🔴 R322 MANDATORY STOP AFTER SPAWNING 🔴🔴🔴
echo "
🛑 STOPPING PER R322 - CONTEXT PRESERVATION
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Agents spawned: [List all SW Engineers]
State saved to: orchestrator-state-v3.json
Next state: MONITORING_SWE_PROGRESS

To continue after agents complete:
  claude --continue

This stop preserves context and prevents rule loss.
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
"
exit 0  # MANDATORY EXIT PER R322
```

## 🔴🔴🔴 SPECIAL HANDLING FOR SPLITS 🔴🔴🔴

### When Spawning SW Engineer for Split Implementation:

```bash
# 1. Read split information from state file (USE ABSOLUTE PATHS!)
EFFORT_NAME="gitea-client"  # Example
SPLIT_NUM=2  # Current split being worked

# 🔴🔴🔴 CRITICAL: GET PATHS FROM STATE FILE! 🔴🔴🔴
SPLIT_PLAN_PATH=$(jq '.split_tracking.\"$EFFORT_NAME\".splits[$((SPLIT_NUM-1))].split_plan_path' orchestrator-state-v3.json)
INFRASTRUCTURE_DIR=$(jq '.split_tracking.\"$EFFORT_NAME\".splits[$((SPLIT_NUM-1))].infrastructure_dir' orchestrator-state-v3.json)

if [ -z "$SPLIT_PLAN_PATH" ] || [ -z "$INFRASTRUCTURE_DIR" ]; then
    echo "❌ FATAL: Split paths not found in state file!"
    echo "Expected .split_tracking.$EFFORT_NAME.splits[$((SPLIT_NUM-1))].split_plan_path"
    echo "This should have been set by CREATE_NEXT_INFRASTRUCTURE"
    exit 1
fi

echo "✅ Split paths loaded from state:"
echo "   Plan: $SPLIT_PLAN_PATH"
echo "   Infrastructure: $INFRASTRUCTURE_DIR"

# 2. Verify split plan exists
if [ ! -f "$SPLIT_PLAN_PATH" ]; then
    echo "❌ FATAL: Split plan not found at: $SPLIT_PLAN_PATH"
    exit 1
fi

# 3. CD to infrastructure directory (USE ABSOLUTE PATH!)
cd "$INFRASTRUCTURE_DIR" || {
    echo "❌ FATAL: Cannot cd to split infrastructure: $INFRASTRUCTURE_DIR"
    exit 1
}

# 4. Verify we're in the right place
pwd
git branch --show-current

# 5. Now spawn with EXPLICIT paths
```

### Split Spawn Message Template:

```markdown
# SPAWN SW ENGINEER FOR SPLIT IMPLEMENTATION:
Task software-engineer:
PURPOSE: Implement Split-{SPLIT_NUM} of {effort_name}

🔴🔴🔴 CRITICAL SPLIT INFORMATION:
YOU ARE IMPLEMENTING A SPLIT!
SPLIT NUMBER: {SPLIT_NUM}
SPLIT PLAN PATH: {SPLIT_PLAN_FULL_PATH}
INFRASTRUCTURE DIR: {INFRASTRUCTURE_DIR}
🔴🔴🔴

🔴🔴🔴 CRITICAL FILE PLACEMENT WARNING (R326) 🔴🔴🔴
DO NOT CREATE split-{SPLIT_NUM}/ SUBDIRECTORY!
Files go DIRECTLY in standard project directories:
✅ CORRECT: pkg/registry/auth.go
❌ WRONG: split-{SPLIT_NUM}/pkg/registry/auth.go

Your working directory is ALREADY split-specific:
{INFRASTRUCTURE_DIR}

Creating split subdirectories causes CATASTROPHIC measurement errors!
Git will see files as different, measuring 1989 lines instead of 400!
🔴🔴🔴

📋 YOUR INSTRUCTIONS:
FOLLOW ONLY: The split plan at {SPLIT_PLAN_FULL_PATH}
LOCATION: {INFRASTRUCTURE_DIR}
IGNORE: Original effort plans, other splits

🔴🔴🔴 NAVIGATION REQUIREMENTS:
YOU MUST USE ABSOLUTE PATHS!
TARGET_DIRECTORY: {INFRASTRUCTURE_DIR}
SPLIT_PLAN: {SPLIT_PLAN_FULL_PATH}

YOUR MANDATORY FIRST ACTIONS:
1. Echo current directory: pwd
2. Navigate using ABSOLUTE path: cd "{INFRASTRUCTURE_DIR}"
3. Verify correct directory: pwd
4. Verify branch: git branch --show-current
5. Read split plan: cat "{SPLIT_PLAN_FULL_PATH}"
6. If paths don't exist:
   - STOP IMMEDIATELY
   - Report exact error with full paths

REQUIREMENTS:
- Follow split plan at {SPLIT_PLAN_FULL_PATH} exactly
- Size limit: {split_limit} lines for this split
- Base branch: {base_branch} (from state file)
- Tests passing for split scope

DELIVERABLES:
- Split implementation complete per plan
- Tests for split passing
- Size under split limit
- Code committed and pushed
```

## Spawn Message Template WITH R295 COMPLIANCE (Regular Efforts)

```markdown
# SPAWN SW ENGINEER WITH MANDATORY CLARITY (R295):
Task software-engineer:
PURPOSE: Implement {effort_id} - {effort_name}

🔴🔴🔴 CRITICAL STATE INFORMATION (R295):
YOU ARE IN STATE: IMPLEMENTATION
This means you should: Implement the features defined in your IMPLEMENTATION-PLAN (in .software-factory/)
🔴🔴🔴

📋 YOUR INSTRUCTIONS (R295):
FOLLOW ONLY: The most recent IMPLEMENTATION-PLAN--*.md in .software-factory/
LOCATION: .software-factory/phase{X}/wave{Y}/{effort-name}/
IGNORE: Any files named *-COMPLETED-*.md or other plan files

🎯 CONTEXT:
- EFFORT: {effort_name}
- WAVE: {wave_number}
- PHASE: {phase_number}
- YOUR TASK: Implement features as specified in the IMPLEMENTATION-PLAN found in .software-factory/

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
- Follow the IMPLEMENTATION-PLAN--*.md in .software-factory/ exactly (R295: This is your ONLY plan)
- Size limit: {limit} lines
- Test coverage: {X}% minimum
- Update work-log--*.log in .software-factory/ every checkpoint

STARTUP REQUIREMENTS:
1. Print: "AGENT STARTUP: $(date '+%Y-%m-%d %H:%M:%S %Z')"
2. Verify pwd matches WORKING_DIR
3. Verify branch matches BRANCH
4. Acknowledge rules R054, R007, R152, R295

DELIVERABLES:
- Implementation complete per IMPLEMENTATION-PLAN--*.md in .software-factory/
- Tests passing at required coverage
- Size under limit
- Work log (work-log--*.log) updated in .software-factory/
- Code committed and pushed
```

## Parallelization Matrix

**See R053-parallelization-decisions.md in rule-library for parallelization decisions guidance.**

Key points:
- Can parallelize: Independent efforts, no shared dependencies
- Must serialize: Dependent efforts, splits, shared files

## Recording Spawn Times

```yaml
# Update in orchestrator-state-v3.json
parallel_spawn_records:
  wave{X}_group{Y}:
    spawned_at: "2025-08-23T14:30:45Z"
    agents:
      - name: "software-engineer-effort1"
        timestamp: "2025-08-23T14:30:47Z"
      - name: "software-engineer-effort2"
        timestamp: "2025-08-23T14:30:49Z"
      - name: "software-engineer-effort3"
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

## 🚨🚨🚨 STATE TRANSITION PROTOCOL (R324/R325) 🚨🚨🚨

**AFTER ALL AGENTS ARE SPAWNED, YOU MUST UPDATE current_state BEFORE STOPPING!**

```bash
# 🔴🔴🔴 MANDATORY: Execute this AFTER all spawns complete! 🔴🔴🔴

echo "✅ All SW Engineers spawned successfully"

# CRITICAL: Update state file FIRST (R324 requirement)
echo "🔴 R324: Updating current_state to prevent infinite loop..."
jq '.transition_time = \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\"' orchestrator-state-v3.json > tmp.json && mv tmp.json orchestrator-state-v3.json

# Verify the update
echo "✅ State updated to:"
grep "current_state:" orchestrator-state-v3.json

# Commit and push IMMEDIATELY
git add orchestrator-state-v3.json
git commit -m "state: transition to MONITORING_SWE_PROGRESS from SPAWN_SW_ENGINEERS (R324)"
git push

# NOW stop per R322
echo "🛑 STATE TRANSITION CHECKPOINT: SPAWN_SW_ENGINEERS → MONITORING_SWE_PROGRESS"
echo "📊 State file updated to: MONITORING_SWE_PROGRESS ✅"
echo "⏸️ STOPPED - Ready to continue in MONITORING_SWE_PROGRESS"
echo "When restarted, will monitor agent progress"
# EXIT HERE
```

**⚠️ FAILURE TO UPDATE current_state = INFINITE LOOP BUG! ⚠️**


### 🔴🔴🔴 RULE R340: Capture and Save Planning File Metadata 🔴🔴🔴

**File**: `$CLAUDE_PROJECT_DIR/rule-library/R340-planning-file-metadata-tracking.md`
**Criticality**: BLOCKING - Must track all planning files for agent discovery

**WHEN AGENTS REPORT PLANNING FILES:**

When Code Reviewer or other agents report "📋 PLANNING FILE CREATED", you MUST:

1. **Parse the metadata from their report**
2. **Update orchestrator-state-v3.json immediately**
3. **Commit the change**

**EXAMPLE HANDLING:**

```bash
# When Code Reviewer reports an effort plan
yq eval '.effort_repo_files.effort_plans["buildah-builder-interface"] = {
  "file_path": "/efforts/phase1/wave2/buildah-builder-interface/.software-factory/phase1/wave2/buildah-builder-interface/IMPLEMENTATION-PLAN-20250120-100000.md",
  "created_at": "2025-01-20T10:00:00Z",
  "created_by": "code-reviewer",
  "target_branch": "phase1/wave2/buildah-builder-interface",
  "phase": 1,
  "wave": 2,
  "effort_name": "buildah-builder-interface",
  "status": "active",
  "replaced_by": null
}' -i orchestrator-state-v3.json

# When Code Reviewer reports split plans
for split in 001 002; do
  yq eval ".effort_repo_files.split_plans[\"oci-types-split-${split}\"] = {
    \"file_path\": \"/efforts/phase1/wave1/oci-types/.software-factory/splits/oci-types-split-${split}/SPLIT-PLAN-20250120-110000.md\",
    \"created_at\": \"2025-01-20T11:00:00Z\",
    \"created_by\": \"code-reviewer\",
    \"target_branch\": \"phase1/wave1/oci-types-split-${split}\",
    \"parent_effort\": \"oci-types\",
    \"split_number\": ${split#00},
    \"total_splits\": 2,
    \"status\": \"active\"
  }" -i orchestrator-state-v3.json
done

# When Code Reviewer reports merge plans
yq eval '.planning_files.merge_plans.wave["phase1_wave2"] = {
  "file_path": "/efforts/phase1/wave2/integration-workspace/WAVE-MERGE-PLAN.md",
  "created_at": "2025-01-20T12:00:00Z",
  "created_by": "code-reviewer",
  "target_branch": "phase1-wave2-integration",
  "phase": 1,
  "wave": 2,
  "efforts_count": 3,
  "status": "active"
}' -i orchestrator-state-v3.json

# Commit immediately
git add orchestrator-state-v3.json
git commit -m "state: track planning file metadata per R340"
git push
```

**VERIFICATION:**
```bash
# Verify plans are tracked
echo "📋 Tracked planning files:"
yq '.planning_files' orchestrator-state-v3.json

# Confirm SW Engineers can find their plans
for effort in $(yq '.efforts_in_progress[].name' orchestrator-state-v3.json); do
  plan_path=$(yq ".effort_repo_files.effort_plans[\"$effort\"].file_path" orchestrator-state-v3.json)
  if [ "$plan_path" != "null" ]; then
    echo "✅ $effort plan tracked: $plan_path"
  else
    echo "❌ WARNING: No plan tracked for $effort"
  fi
done
```

### 🔴🔴🔴 MANDATORY VALIDATION REQUIREMENT 🔴🔴🔴

**Per R288 and R324**: ALL state file updates MUST be validated before commit:

```bash
# After ANY update to orchestrator-state-v3.json:
"$CLAUDE_PROJECT_DIR/tools/validate-state.sh" orchestrator-state-v3.json || {
    echo "❌ State file validation failed!"
    exit 288
}
```

**Use helper functions for automatic validation:**
```bash
# Source the helper functions
source "$CLAUDE_PROJECT_DIR/utilities/state-file-update-functions.sh"

# Use safe functions that include validation:
safe_state_transition "NEW_STATE" "reason"
safe_update_field "field_name" "value"
```


## Next State Transitions

### VALID TRANSITIONS FROM SPAWN_SW_ENGINEERS (from state-machines/software-factory-3.0-state-machine.json):
1. **MONITORING_SWE_PROGRESS** - After spawning SW Engineers for implementation (ONLY valid transition)

### Transition Logic:
- After successfully spawning agents → MONITORING_SWE_PROGRESS (mandatory per R322 stop then continue)
- If spawn fails → ERROR_RECOVERY

**IMPORTANT**: Per R322, you MUST stop after spawning agents and wait for /continue-orchestrating


## 🔴🔴🔴 STATE COMPLETION CHECKLIST 🔴🔴🔴

**Execute these steps IN ORDER to properly complete SPAWN_SW_ENGINEERS:**

### ✅ Step 1: Complete State-Specific Work
**Refer to "Primary Actions" and "Mandatory Validations" sections above for state-specific tasks.**

**CRITICAL**: During state work, set these variables (DO NOT call update_state):
```bash
PROPOSED_NEXT_STATE="[state determined by logic]"
TRANSITION_REASON="[reason for transition]"
```

Once all state work is complete, proceed to mandatory exit protocol:

---

### ✅ Step 2: Verify Proposed Transition Set
```bash
# Verify PROPOSED_NEXT_STATE was set during state work
if [ -z "$PROPOSED_NEXT_STATE" ]; then
    echo "❌ CRITICAL: PROPOSED_NEXT_STATE not set!"
    echo "State work must set PROPOSED_NEXT_STATE and TRANSITION_REASON"
    exit 288
fi

echo "✅ Proposed transition: SPAWN_SW_ENGINEERS → $PROPOSED_NEXT_STATE"
echo "   Reason: $TRANSITION_REASON"
```

---

### ✅ Step 3: Save Proposed Transition to File (R288)
```bash
# Write proposed transition to file for State Manager
cat > "$CLAUDE_PROJECT_DIR/.proposed-transition" <<EOF
PROPOSED_NEXT_STATE="$PROPOSED_NEXT_STATE"
TRANSITION_REASON="$TRANSITION_REASON"
CURRENT_STATE="SPAWN_SW_ENGINEERS"
TIMESTAMP="$(date -u +%Y-%m-%dT%H:%M:%SZ)"
EOF

echo "✅ Proposed transition saved to .proposed-transition"
```

---

### ✅ Step 4: Validate Current State File (R324)
```bash
# Validate state file integrity BEFORE proposing transition
"$CLAUDE_PROJECT_DIR/tools/validate-state.sh" orchestrator-state-v3.json || {
    echo "❌ State file validation failed!"
    exit 288
}
echo "✅ Current state file validated"
```

---

### ✅ Step 5: Commit Work Products (if any)
```bash
# Commit any work products created during this state
# (State file itself is committed by State Manager)
if [ -n "$(git status --porcelain | grep -v orchestrator-state-v3.json)" ]; then
    git add .
    git commit -m "work: SPAWN_SW_ENGINEERS work products [R288]"
    git push
    echo "✅ Work products committed and pushed"
else
    echo "✅ No work products to commit"
fi
```

---

### ✅ Step 4: Save TODOs (R287 - SUPREME LAW)
```bash
# Save TODO state before transition (R287 trigger)
save_todos "SPAWN_SW_ENGINEERS_COMPLETE"

# Commit TODOs within 60 seconds (R287)
cd "$CLAUDE_PROJECT_DIR"
git add todos/*.todo
git commit -m "todo: orchestrator - SPAWN_SW_ENGINEERS complete [R287]"
git push
echo "✅ TODOs saved and committed"
```

---

### ✅ Step 7: Output Transition Proposal (R288)
```bash
# Output the proposed transition for State Manager
echo "════════════════════════════════════════════════════════"
echo "PROPOSED STATE TRANSITION (R288):"
echo "  FROM: SPAWN_SW_ENGINEERS"
echo "  TO:   $PROPOSED_NEXT_STATE"
echo "  REASON: $TRANSITION_REASON"
echo "════════════════════════════════════════════════════════"
```

---

### ✅ Step 8: Output Continuation Flag (R405 - SUPREME LAW)
```bash
# Output continuation flag (R405)
# TRUE = state work complete, transition proposed
# FALSE = error occurred, manual intervention needed

if [ "$PROPOSED_NEXT_STATE" = "ERROR_RECOVERY" ]; then
    echo "CONTINUE-SOFTWARE-FACTORY=FALSE"
else
    echo "CONTINUE-SOFTWARE-FACTORY=TRUE"
fi
```

**⚠️ THIS MUST BE THE LAST OUTPUT BEFORE STOP! ⚠️**

---

### ✅ Step 9: State Manager Transition Checkpoint
```
🔄 STATE MANAGER TAKES CONTROL HERE (R288)

The State Manager (run-software-factory.sh) will:
1. Read .proposed-transition file
2. Validate proposed transition against state machine
3. Update orchestrator-state-v3.json with new state
4. Commit and push state file
5. Re-invoke orchestrator in new state (if CONTINUE-SOFTWARE-FACTORY=TRUE)

DO NOT PROCEED PAST THIS POINT - STATE MANAGER HANDLES TRANSITION
```

---

### ✅ Step 10: Stop Processing (R322 - SUPREME LAW)
```bash
# Stop for State Manager handoff (R322)
echo "🛑 State work complete - State Manager will handle transition"
exit 0
```

---

## 🚨 CHECKLIST ENFORCEMENT 🚨

**Skipping ANY step in this checklist = FAILURE:**
- Missing Step 1: State work incomplete
- Missing Step 2: No proposed transition = stuck forever (R288 violation, -100%)
- Missing Step 3: State Manager can't read proposal = broken automation (R288 violation, -100%)
- Missing Step 4: Invalid state = corruption (R324 violation)
- Missing Step 5: Work products lost (data loss)
- Missing Step 6: TODOs lost on compaction (R287 violation, -20% to -100%)
- Missing Step 7: No transition visibility (R288 violation)
- Missing Step 8: Automation doesn't know if it can continue (R405 violation, -100%)
- Missing Step 9: State Manager handoff failed (R288 violation, -100%)
- Missing Step 10: Context corruption (R322 violation, -100%)

- Missing Step 2: No proposed next state = State Manager can't proceed
- Missing Step 3: No State Manager consultation = bypassing bookend pattern (-100%)
- Missing Step 4: No TODO save = work lost (R287 violation, -20% to -100%)
- Missing Step 5: No continuation flag = automation stops (R405 violation, -100%)
- Missing Step 6: No exit = R322 violation (-100%)

**ALL 6 STEPS ARE MANDATORY - NO EXCEPTIONS**
**NOTE**: State file validation and commits now handled by State Manager (SF 3.0 pattern) - NO EXCEPTIONS**

## 🔴🔴🔴 R405 - MANDATORY AUTOMATION CONTINUATION FLAG 🔴🔴🔴

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
echo "CONTINUE-SOFTWARE-FACTORY=TRUE"

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
### 🚨 SPAWN STATE PATTERN - R322 + R405 USAGE 🚨

**Spawning operations require R322 stop for context preservation:**
```bash
# After spawning agent(s)
echo "✅ Spawned agents for work"

# R322 checkpoint (context preservation)
echo "🛑 R322: Stopping after spawn for context preservation"

# Flag? → MUST BE TRUE (normal operation!)
echo "CONTINUE-SOFTWARE-FACTORY=TRUE"

# Stop inference
exit 0
```

**Why TRUE is correct:**
- Spawning is NORMAL operation
- System knows next state
- Automation can continue
- **Context preservation ≠ manual intervention needed!**
**See: $CLAUDE_PROJECT_DIR/rule-library/R405-automation-continuation-flag.md**

