# Orchestrator - SETUP_INTEGRATION_INFRASTRUCTURE State Rules

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

## 🔴🔴🔴 CRITICAL: INTEGRATION BASE BRANCH DETERMINATION (R308) 🔴🔴🔴

**VIOLATION = -100% AUTOMATIC FAILURE**

### YOU MUST DETERMINE THE CORRECT BASE BRANCH PER R308:
- **Wave Integration**: Base on the PREVIOUS integration or main (if Wave 1)
- **Phase Integration**: Base on the last wave integration of the phase
- **Project Integration**: Base on the last phase integration

### INTEGRATION BRANCHES FOLLOW R308 INCREMENTAL STRATEGY:
```
Phase 1, Wave 1 integration: from main
Phase 1, Wave 2 integration: from phase1-wave1-integration (NOT main!)
Phase 2, Wave 1 integration: from phase1-integration (NOT main!)
Phase 2, Wave 2 integration: from phase2-wave1-integration (NOT main!)
```

### 🔴 CRITICAL EXAMPLE:
```bash
# Phase 2, Wave 1 Integration
# WRONG - AUTOMATIC FAILURE:
BASE_BRANCH="main"  # ❌ NEVER for Phase 2!

# CORRECT:
BASE_BRANCH="phase1-integration"  # ✅ From previous phase!
```

**Acknowledge: "I understand integration branches MUST follow R308 incremental strategy"**

## 🔴🔴🔴 STOP! STATE RULE READING IS ABSOLUTELY FIRST! 🔴🔴🔴

**YOU HAVE ENTERED SETUP_INTEGRATION_INFRASTRUCTURE STATE - YOU MUST READ AND ACKNOWLEDGE ALL STATE RULES BEFORE DOING ANY STATE WORK!**

## 🔴🔴🔴 R290 VERIFICATION REQUIREMENT 🔴🔴🔴

**R290 ENFORCEMENT: CREATE VERIFICATION MARKER AFTER READING**

After reading and acknowledging all state rules, you MUST create a verification marker:

```bash
# MANDATORY: Create verification marker after reading rules
touch .state_rules_read_orchestrator_SETUP_INTEGRATION_INFRASTRUCTURE
echo "$(date +%s) - Rules read and acknowledged for SETUP_INTEGRATION_INFRASTRUCTURE" > .state_rules_read_orchestrator_SETUP_INTEGRATION_INFRASTRUCTURE
```

**FAILURE TO CREATE MARKER = AUTOMATIC -100% PENALTY**

The system will check for this marker. No marker = Immediate failure.

### ❌ DO NOT DO ANY SETUP_INTEGRATION_INFRASTRUCTURE WORK UNTIL RULES ARE READ:
- ❌ Start creating integration workspace
- ❌ Start creating integration branch
- ❌ Start cloning repository
- ❌ Start pushing to remote
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
   ❌ WRONG: "I acknowledge R308, R250, R034..."
   (No Read tool calls detected = AUTOMATIC FAILURE)
   ```

2. **Bulk Acknowledgment**:
   ```
   ❌ WRONG: "I acknowledge all SETUP_INTEGRATION_INFRASTRUCTURE rules"
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
   ❌ WRONG: "I know R308 requires incremental bases..."
   (Must READ from file, not recall from memory)
   ```

5. **Skipping Rules in PRIMARY DIRECTIVES**:
   ```
   ❌ WRONG: Reading only some rules from the list
   (ALL rules in PRIMARY DIRECTIVES are MANDATORY)
   ```

### ✅ CORRECT PATTERN FOR SETUP_INTEGRATION_INFRASTRUCTURE:
```
1. READ: $CLAUDE_PROJECT_DIR/rule-library/[first-rule-file].md
2. "I acknowledge [Rule#] - [Rule Name]: [Brief description]"
3. READ: $CLAUDE_PROJECT_DIR/rule-library/[second-rule-file].md  
4. "I acknowledge [Rule#] - [Rule Name]: [Brief description]"
[Continue for EVERY rule in PRIMARY DIRECTIVES...]
5. "Ready to execute SETUP_INTEGRATION_INFRASTRUCTURE work"
```

### 🚨 NO WORK UNTIL ACKNOWLEDGMENT COMPLETE 🚨
**You may NOT begin ANY SETUP_INTEGRATION_INFRASTRUCTURE work until:**
1. ✅ ALL rules in PRIMARY DIRECTIVES have been READ
2. ✅ ALL rules have been individually ACKNOWLEDGED
3. ✅ You have stated readiness to execute SETUP_INTEGRATION_INFRASTRUCTURE work
1. **READ** every rule file listed in PRIMARY DIRECTIVES below
2. **ACKNOWLEDGE** each rule individually with number and description
3. **ONLY THEN** proceed with SETUP_INTEGRATION_INFRASTRUCTURE work

### 🚨 FAILURE TO READ STATE RULES FIRST = IMMEDIATE EXIT 🚨
**If you do ANY SETUP_INTEGRATION_INFRASTRUCTURE work before reading and acknowledging rules:**
- **STOP ALL WORK IMMEDIATELY**
- **EXIT WITH FAILURE STATUS**
- **YOU HAVE VIOLATED STATE COMPLIANCE**

**THE SYSTEM IS MONITORING YOUR READ TOOL CALLS!**

## 📋 PRIMARY DIRECTIVES FOR SETUP_INTEGRATION_INFRASTRUCTURE

**YOU MUST READ EACH RULE LISTED HERE. YOUR READ TOOL CALLS ARE BEING MONITORED.**

### 🔴🔴🔴 R308 - INCREMENTAL BRANCHING STRATEGY (SUPREME LAW!)
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R308-incremental-branching-strategy.md`
**Criticality**: 🔴🔴🔴 SUPREME LAW - Integration branches MUST be incremental!
**Summary**: Wave N+1 integration based on Wave N integration, Phase N+1 based on Phase N
**CRITICAL**: Phase 2 Wave 1 integration MUST use phase1-integration, NOT main!

### 🚨🚨🚨 R250 - Integration Isolation Requirement
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R250-integration-isolation-requirement.md`
**Criticality**: BLOCKING - Integration must use separate target clone
**Summary**: Integration must happen under /efforts/ directory structure

### 🚨🚨🚨 R034 - Integration Requirements
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R034-integration-requirements.md`
**Criticality**: BLOCKING - Required for wave approval
**Summary**: Complete integration protocol with testing and validation

### 🚨🚨🚨 R014 - Branch Naming Convention
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R014-branch-naming-convention.md`
**Criticality**: BLOCKING - Mandatory project prefix for all branches
**Summary**: Use project prefix for all integration branches

### 🚨🚨🚨 R271 - Mandatory Production-Ready Validation
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R271-mandatory-production-ready-validation.md`
**Criticality**: BLOCKING - Full checkouts required for integration
**Summary**: Integration must use full repository clones, no sparse checkouts

### 🚨🚨🚨 R006 - Orchestrator NEVER Writes Code [BLOCKING]
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R006-orchestrator-never-writes-code.md`
**Criticality**: BLOCKING - Any code operation = -100% IMMEDIATE FAILURE
**Summary**: Orchestrator coordinates but NEVER implements or fixes code

### 🚨🚨🚨 R329 - Orchestrator NEVER Performs Git Merges [BLOCKING]
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R329-orchestrator-never-performs-merges.md`
**Criticality**: BLOCKING - Any merge operation = -100% IMMEDIATE FAILURE
**Summary**: Orchestrator MUST spawn Integration Agent for ALL merges - NO EXCEPTIONS

### 🚨🚨🚨 R307 - Independent Branch Mergeability [BLOCKING]
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R307-independent-branch-mergeability.md`
**Criticality**: BLOCKING - Must verify branches are mergeable before attempting
**Summary**: Check for conflicts and mergeability before integration operations

### 🚨🚨🚨 R216 - Bash Execution Syntax Protocol (BLOCKING)
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R216-bash-execution-syntax.md`
**Criticality**: BLOCKING - Incorrect syntax causes failures
**Summary**: Use parentheses for subshells, proper variable syntax

### 🚨🚨🚨 R288 - State File Update and Commit Protocol
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R288-state-file-update-and-commit-protocol.md`
**Criticality**: SUPREME LAW - Update on every transition
**Summary**: Update orchestrator-state.json with integration metadata

## 🚨 SETUP_INTEGRATION_INFRASTRUCTURE IS A VERB - CREATE INFRASTRUCTURE NOW! 🚨

### 🔴🔴🔴 CRITICAL: YOU ARE ALREADY IN SETUP_INTEGRATION_INFRASTRUCTURE STATE! 🔴🔴🔴

**If current_state = "SETUP_INTEGRATION_INFRASTRUCTURE" in orchestrator-state.json, you MUST:**
1. **IMMEDIATELY** start creating integration infrastructure
2. **NO ANNOUNCEMENTS** - just start working
3. **NO WAITING** - infrastructure creation begins NOW

### IMMEDIATE ACTIONS UPON ENTERING SETUP_INTEGRATION_INFRASTRUCTURE

**THE MOMENT YOU SEE current_state: SETUP_INTEGRATION_INFRASTRUCTURE, YOU MUST:**
1. Determine the integration type (wave/phase/project) NOW
2. Determine the correct incremental base branch per R308
3. Create integration working directory immediately
4. Clone repository with FULL checkout (R271)
5. Create integration branch following R308
6. Push integration branch to remote
7. Update state file with integration metadata
8. Transition to appropriate next state

**FORBIDDEN - AUTOMATIC FAILURE:**
- ❌ "STATE TRANSITION COMPLETE: Now in SETUP_INTEGRATION_INFRASTRUCTURE" [stops]
- ❌ "Successfully entered SETUP_INTEGRATION_INFRASTRUCTURE state" [waits]
- ❌ "Ready to setup integration infrastructure" [pauses]
- ❌ "I'm in integration infrastructure state" [does nothing]
- ❌ "Preparing to create integration workspace..." [delays]
- ❌ "I see we're in SETUP_INTEGRATION_INFRASTRUCTURE state..." [announces]

**REQUIRED - IMMEDIATE ACTION:**
- ✅ "SETUP_INTEGRATION_INFRASTRUCTURE: Determining integration type NOW..."
- ✅ "Creating wave integration infrastructure at /efforts/phase${X}/wave${Y}/integration-workspace..."
- ✅ "Applying R308: Phase 2 Wave 1 integration will use phase1-integration as base..."

## State Context
You are CREATING INFRASTRUCTURE for integration, not executing merges. Your responsibilities:
1. **DETERMINE** correct incremental base branch per R308
2. **CREATE** integration workspace directory
3. **CLONE** target repository with FULL checkout
4. **CREATE** integration branch with proper naming
5. **PUSH** integration branch to establish remote
6. **UPDATE** state file with integration metadata
7. **TRANSITION** to next appropriate state

**YOU MUST NEVER (R329 + R006 ENFORCEMENT):**
- ❌ Execute any git merges (R329 VIOLATION)
- ❌ Resolve any conflicts (R329 VIOLATION)
- ❌ Run any builds or tests (R006 VIOLATION)
- ❌ Write any code (R006 VIOLATION)

## 🔴🔴🔴 CRITICAL: R308 INCREMENTAL BASE DETERMINATION 🔴🔴🔴

### Integration Base Branch Function
```bash
determine_integration_base_branch() {
    local INTEGRATION_TYPE="$1"  # wave, phase, or project
    local PHASE="$2"
    local WAVE="$3"
    
    echo "🔴 R308: Determining incremental base for $INTEGRATION_TYPE integration"
    
    case "$INTEGRATION_TYPE" in
        "wave")
            # Wave integration bases
            if [[ $PHASE -eq 1 && $WAVE -eq 1 ]]; then
                echo "📌 R308: Phase 1 Wave 1 integration → base: main"
                echo "main"
            elif [[ $WAVE -eq 1 ]]; then
                # First wave of new phase: from previous phase integration
                PREV_PHASE=$((PHASE - 1))
                BASE="phase${PREV_PHASE}-integration"
                echo "🔴 R308: Phase $PHASE Wave 1 integration → base: $BASE (NOT main!)"
                echo "$BASE"
            else
                # Subsequent waves: from previous wave integration
                PREV_WAVE=$((WAVE - 1))
                BASE="phase${PHASE}-wave${PREV_WAVE}-integration"
                echo "📌 R308: Phase $PHASE Wave $WAVE integration → base: $BASE"
                echo "$BASE"
            fi
            ;;
            
        "phase")
            # Phase integration: from last wave of the phase
            LAST_WAVE=$(jq ".phases[] | select(.number == $PHASE) | .waves | length" phase-plans/PROJECT-PHASES.json)
            BASE="phase${PHASE}-wave${LAST_WAVE}-integration"
            echo "📌 R308: Phase $PHASE integration → base: $BASE"
            echo "$BASE"
            ;;
            
        "project")
            # Project integration: from last phase integration
            LAST_PHASE=$(jq '.phases | length' phase-plans/PROJECT-PHASES.json)
            BASE="phase${LAST_PHASE}-integration"
            echo "📌 R308: Project integration → base: $BASE"
            echo "$BASE"
            ;;
    esac
}
```

## Infrastructure Setup Protocol (DETERMINISTIC)

```bash
# 🔴🔴🔴 DETERMINISTIC INTEGRATION INFRASTRUCTURE SETUP 🔴🔴🔴
setup_integration_infrastructure() {
    local INTEGRATION_TYPE="$1"  # wave, phase, or project
    
    echo "═══════════════════════════════════════════════════════"
    echo "🔧 SETTING UP INTEGRATION INFRASTRUCTURE (DETERMINISTIC)"
    echo "Type: $INTEGRATION_TYPE"
    echo "═══════════════════════════════════════════════════════"
    
    # 0. MUST be in SF instance directory
    SF_INSTANCE_DIR=$(pwd)
    if [ ! -f "${SF_INSTANCE_DIR}/orchestrator-state.json" ]; then
        echo "❌ ERROR: Not in SF instance directory!"
        exit 1
    fi
    
    # 1. Determine phase/wave from state
    PHASE=$(jq -r '.current_phase' orchestrator-state.json)
    WAVE=$(jq -r '.current_wave' orchestrator-state.json)
    
    # 2. Source branch naming helpers
    if [ -f "$SF_INSTANCE_DIR/utilities/branch-naming-helpers.sh" ]; then
        source "$SF_INSTANCE_DIR/utilities/branch-naming-helpers.sh"
    fi
    PROJECT_PREFIX=$(yq '.branch_naming.project_prefix' "$SF_INSTANCE_DIR/target-repo-config.yaml")
    
    # 3. CRITICAL: Determine incremental base per R308
    echo "🔴🔴🔴 R308 ENFORCEMENT: Determining INCREMENTAL base branch"
    BASE_BRANCH=$(determine_integration_base_branch "$INTEGRATION_TYPE" "$PHASE" "$WAVE")
    
    # CRITICAL VALIDATION
    if [[ "$INTEGRATION_TYPE" == "wave" && $PHASE -gt 1 && $WAVE -eq 1 && "$BASE_BRANCH" == "main" ]]; then
        echo "🔴🔴🔴 R308 VIOLATION: Phase $PHASE Wave 1 CANNOT use main!"
        echo "FATAL ERROR: Must use phase$((PHASE-1))-integration"
        exit 308
    fi
    
    echo "✅ R308 VALIDATED: Using incremental base: $BASE_BRANCH"
    
    # 4. Verify base branch exists
    TARGET_REPO_URL=$(yq '.target_repository.url' "$SF_INSTANCE_DIR/target-repo-config.yaml")
    if ! git ls-remote --heads "$TARGET_REPO_URL" "$BASE_BRANCH" > /dev/null 2>&1; then
        echo "❌ R308 FATAL: Base branch not found: $BASE_BRANCH"
        echo "Previous integration must be completed first!"
        exit 308
    fi
    
    # 5. DETERMINISTIC integration workspace path (R250 compliant)
    case "$INTEGRATION_TYPE" in
        "wave")
            INTEGRATION_WORKSPACE="${SF_INSTANCE_DIR}/efforts/phase${PHASE}/wave${WAVE}/integration-workspace"
            INTEGRATION_BRANCH=$(get_wave_integration_branch_name "$PHASE" "$WAVE" 2>/dev/null || echo "phase${PHASE}-wave${WAVE}-integration")
            ;;
        "phase")
            INTEGRATION_WORKSPACE="${SF_INSTANCE_DIR}/efforts/phase${PHASE}/phase-integration-workspace"
            INTEGRATION_BRANCH=$(get_phase_integration_branch_name "$PHASE" 2>/dev/null || echo "phase${PHASE}-integration")
            ;;
        "project")
            INTEGRATION_WORKSPACE="${SF_INSTANCE_DIR}/efforts/project-integration-workspace"
            INTEGRATION_BRANCH=$(get_project_integration_branch_name 2>/dev/null || echo "project-integration")
            ;;
    esac
    
    # DETERMINISTIC: The actual repo goes INTO integration-workspace as 'repo'
    INTEGRATION_DIR="${INTEGRATION_WORKSPACE}/repo"
    
    echo "📍 Integration workspace: $INTEGRATION_WORKSPACE"
    echo "📍 Repository will be at: $INTEGRATION_DIR"
    
    # 6. Handle re-integration (DETERMINISTIC ARCHIVING)
    if [ -d "$INTEGRATION_WORKSPACE" ]; then
        ARCHIVE_NUM=1
        while [ -d "${INTEGRATION_WORKSPACE}-archived-${ARCHIVE_NUM}" ]; do
            ARCHIVE_NUM=$((ARCHIVE_NUM + 1))
        done
        echo "📦 Archiving previous integration to: ${INTEGRATION_WORKSPACE}-archived-${ARCHIVE_NUM}"
        mv "$INTEGRATION_WORKSPACE" "${INTEGRATION_WORKSPACE}-archived-${ARCHIVE_NUM}"
    fi
    
    # 7. Create fresh integration workspace
    mkdir -p "$INTEGRATION_WORKSPACE"
    
    # 8. Load and validate target repository URL
    if [ -z "$TARGET_REPO_URL" ] || [ "$TARGET_REPO_URL" = "null" ]; then
        echo "🔴 ERROR: No target repository URL in config!"
        exit 191
    fi
    
    # 9. SINGLE-BRANCH FULL clone INTO workspace/repo (R271 + R250)
    echo "📦 Cloning target repo INTO: ${INTEGRATION_WORKSPACE}/repo"
    echo "   From branch: $BASE_BRANCH"
    
    git clone \
        --single-branch \
        --branch "$BASE_BRANCH" \
        "$TARGET_REPO_URL" \
        "$INTEGRATION_DIR"
    
    if [ $? -ne 0 ]; then
        echo "❌ Clone failed! Check if base branch '$BASE_BRANCH' exists"
        exit 1
    fi
    
    cd "$INTEGRATION_DIR"
    
    # 8. Verify FULL checkout (R271 compliance)
    if [ -f ".git/info/sparse-checkout" ]; then
        echo "🔴🔴🔴 SUPREME LAW VIOLATION: Sparse checkout detected in integration!"
        exit 271
    fi
    echo "✅ Full codebase available for integration from $BASE_BRANCH"
    
    # 9. Create integration branch
    echo "Creating integration branch: $INTEGRATION_BRANCH"
    git checkout -b "$INTEGRATION_BRANCH"
    
    # 10. Push integration branch to establish remote
    git push -u origin "$INTEGRATION_BRANCH"
    
    # 11. Create integration metadata file
    cat > INTEGRATION-METADATA.md << EOF
# Integration Infrastructure Metadata

## Integration Details
- **Type**: $INTEGRATION_TYPE
- **Phase**: $PHASE
- **Wave**: $WAVE
- **Branch**: $INTEGRATION_BRANCH
- **Base Branch**: $BASE_BRANCH
- **Created**: $(date)

## R308 Incremental Branching Compliance
- **Rule Applied**: Integration branch properly based on $BASE_BRANCH
- **Verification**: This integration builds on all previous integrated work
$(if [[ $PHASE -eq 2 && $WAVE -eq 1 ]]; then
    echo "- **CRITICAL**: Phase 2 Wave 1 correctly based on phase1-integration (NOT main)"
elif [[ $PHASE -gt 1 ]]; then
    echo "- **Incremental**: Building on previous phase/wave integrations as required"
fi)

## Next Steps
1. Spawn Code Reviewer to create merge plan
2. Spawn Integration Agent to execute merges
3. Monitor integration progress
4. Spawn Code Reviewer for validation
EOF
    
    git add INTEGRATION-METADATA.md
    git commit -m "chore: initialize $INTEGRATION_TYPE integration infrastructure with R308 compliance"
    git push
    
    # 12. Update state file with integration metadata
    cd "$SF_INSTANCE_DIR"
    
    jq --arg type "$INTEGRATION_TYPE" \
       --arg branch "$INTEGRATION_BRANCH" \
       --arg base "$BASE_BRANCH" \
       --arg dir "$INTEGRATION_DIR" \
       --arg phase "$PHASE" \
       --arg wave "$WAVE" \
       '.integration_infrastructure = {
           "type": $type,
           "phase": ($phase | tonumber),
           "wave": ($wave | tonumber),
           "branch": $branch,
           "base_branch": $base,
           "directory": $dir,
           "status": "infrastructure_ready",
           "created_at": now | strftime("%Y-%m-%dT%H:%M:%SZ")
       }' orchestrator-state.json > tmp.json && mv tmp.json orchestrator-state.json
    
    echo "✅ Integration infrastructure ready"
    echo "📋 Details:"
    echo "   - Type: $INTEGRATION_TYPE"
    echo "   - Branch: $INTEGRATION_BRANCH"
    echo "   - Base: $BASE_BRANCH (R308 compliant)"
    echo "   - Directory: $INTEGRATION_DIR"
}
```

## State Transitions

From SETUP_INTEGRATION_INFRASTRUCTURE state:
- **SPAWN_CODE_REVIEWER_MERGE_PLAN** - Infrastructure ready, spawn Code Reviewer for merge plan
- **SPAWN_CODE_REVIEWER_PHASE_MERGE_PLAN** - For phase integration
- **SPAWN_CODE_REVIEWER_PROJECT_MERGE_PLAN** - For project integration
- **ERROR_RECOVERY** - Infrastructure creation failed

**HOW YOU GOT HERE (Valid paths):**
- WAVE_COMPLETE → SETUP_INTEGRATION_INFRASTRUCTURE (wave integration)
- PHASE_INTEGRATION → SETUP_INTEGRATION_INFRASTRUCTURE (phase integration)
- PROJECT_INTEGRATION → SETUP_INTEGRATION_INFRASTRUCTURE (project integration)

## Common Mistakes to Avoid

❌ **WRONG:** Using main as base for Phase 2+ integrations
❌ **WRONG:** Creating integration branches without verifying base exists
❌ **WRONG:** Executing merges in this state (R329 violation)
❌ **WRONG:** Skipping R308 incremental base determination

✅ **CORRECT:** Follow R308 strictly for all integration base branches
✅ **CORRECT:** Create infrastructure only, delegate merges to Integration Agent
✅ **CORRECT:** Verify base branch exists before cloning

## Key Points

1. **R308 is PARAMOUNT** - Integration branches MUST be incremental
2. **Phase 2 Wave 1 MUST use phase1-integration** - NOT main!
3. **Create infrastructure ONLY** - No merges in this state
4. **Full checkouts required** - No sparse checkouts (R271)
5. **Push immediately** - Establish remote tracking
6. **Update state file** - Track integration metadata

## Grading Impact

- **Wrong base branch (violating R308)**: -100%
- **Using main for Phase 2+ integration**: -100%
- **Executing merges (R329 violation)**: -100%
- **Sparse checkout**: -100%
- **Not pushing to remote**: -50%
- **Not updating state file**: -50%

## R322 VIOLATION DETECTION

If you find yourself:
- Starting work for a new state without /continue-orchestrating
- Transitioning without stopping after state file commit
- Continuing after completing state work

**STOP IMMEDIATELY - You are violating R322!**

### 🔴🔴🔴 MANDATORY VALIDATION REQUIREMENT 🔴🔴🔴

**Per R288 and R324**: ALL state file updates MUST be validated before commit:

```bash
# After ANY update to orchestrator-state.json:
"$CLAUDE_PROJECT_DIR/tools/validate-state.sh" orchestrator-state.json || {
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

---

**REMEMBER**: This state creates integration INFRASTRUCTURE with proper R308 incremental bases. It does NOT execute merges - that's for the Integration Agent!