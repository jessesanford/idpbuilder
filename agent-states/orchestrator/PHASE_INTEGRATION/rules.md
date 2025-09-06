# Orchestrator - PHASE_INTEGRATION State Rules

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

**YOU HAVE ENTERED PHASE_INTEGRATION STATE - YOU MUST READ AND ACKNOWLEDGE ALL STATE RULES BEFORE DOING ANY STATE WORK!**

## 🔴🔴🔴 R290 VERIFICATION REQUIREMENT 🔴🔴🔴

**R290 ENFORCEMENT: CREATE VERIFICATION MARKER AFTER READING**

After reading and acknowledging all state rules, you MUST create a verification marker:

```bash
# MANDATORY: Create verification marker after reading rules
touch .state_rules_read_orchestrator_PHASE_INTEGRATION
echo "$(date +%s) - Rules read and acknowledged for PHASE_INTEGRATION" > .state_rules_read_orchestrator_PHASE_INTEGRATION
```

**FAILURE TO CREATE MARKER = AUTOMATIC -100% PENALTY**

The system will check for this marker. No marker = Immediate failure.

### ❌ DO NOT DO ANY PHASE_INTEGRATION WORK UNTIL RULES ARE READ:
- ❌ Start merge wave branches
- ❌ Start create phase branch
- ❌ Start integrate wave work
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
   ❌ WRONG: "I acknowledge all PHASE_INTEGRATION rules"
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

### ✅ CORRECT PATTERN FOR PHASE_INTEGRATION:
```
1. READ: $CLAUDE_PROJECT_DIR/rule-library/[first-rule-file].md
2. "I acknowledge [Rule#] - [Rule Name]: [Brief description]"
3. READ: $CLAUDE_PROJECT_DIR/rule-library/[second-rule-file].md  
4. "I acknowledge [Rule#] - [Rule Name]: [Brief description]"
[Continue for EVERY rule in PRIMARY DIRECTIVES...]
5. "Ready to execute PHASE_INTEGRATION work"
```

### 🚨 NO WORK UNTIL ACKNOWLEDGMENT COMPLETE 🚨
**You may NOT begin ANY PHASE_INTEGRATION work until:**
1. ✅ ALL rules in PRIMARY DIRECTIVES have been READ
2. ✅ ALL rules have been individually ACKNOWLEDGED
3. ✅ You have stated readiness to execute PHASE_INTEGRATION work
1. **READ** every rule file listed in PRIMARY DIRECTIVES below
2. **ACKNOWLEDGE** each rule individually with number and description
3. **ONLY THEN** proceed with PHASE_INTEGRATION work

### 🚨 FAILURE TO READ STATE RULES FIRST = IMMEDIATE EXIT 🚨
**If you do ANY PHASE_INTEGRATION work before reading and acknowledging rules:**
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

## 📋 PRIMARY DIRECTIVES FOR PHASE_INTEGRATION STATE

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

### 🔴🔴🔴 R301 - Integration Branch Current Tracking (SUPREME LAW)
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R301-integration-branch-current-tracking.md`
**Criticality**: SUPREME LAW - Only ONE current integration allowed
**Summary**: Track current vs deprecated integrations, prevent wrong branch usage

### 🚨🚨🚨 R285 - Mandatory Phase Integration Before Assessment  
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R285-mandatory-phase-integration-before-assessment.md`
**Criticality**: BLOCKING - Must integrate before assessment
**Summary**: Phase integration required in normal flow (from WAVE_REVIEW)

### 🚨🚨🚨 R259 - Mandatory Phase Integration After Fixes
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R259-mandatory-phase-integration-after-fixes.md`
**Criticality**: BLOCKING - Must create integration branch after fixes
**Summary**: Create phase-level integration after ERROR_RECOVERY fixes

### 🚨🚨🚨 R257 - Mandatory Phase Assessment Report
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R257-mandatory-phase-assessment-report.md`
**Criticality**: BLOCKING - Required for phase completion
**Summary**: Verify all assessment issues are addressed

### 🚨🚨🚨 R014 - Branch Naming Convention
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R014-branch-naming-convention.md`
**Criticality**: BLOCKING - Mandatory project prefix for all branches
**Summary**: Use project prefix for phase integration branches

### 🚨🚨🚨 R296 - Deprecated Branch Marking Protocol
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R296-deprecated-branch-marking-protocol.md`
**Criticality**: BLOCKING - Prevents integration of wrong branches
**Summary**: Check for and prevent integration of deprecated split branches

### 🚨🚨🚨 R271 - Mandatory Production-Ready Validation
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R271-mandatory-production-ready-validation.md`
**Criticality**: BLOCKING - Full checkouts required
**Summary**: Phase integration must use full repository clones

### 🚨🚨🚨 R269 - Code Reviewer Merge Plan No Execution
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R269-code-reviewer-merge-plan-no-execution.md`
**Criticality**: BLOCKING - Code Reviewer only plans
**Summary**: Code Reviewer creates plan, Integration Agent executes

### 🚨🚨🚨 R260 - Integration Agent Core Requirements
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R260-integration-agent-core-requirements.md`
**Criticality**: BLOCKING - Integration Agent requirements
**Summary**: Integration Agent must acknowledge INTEGRATION_DIR

### 🔴🔴🔴 R321 - Immediate Backport During Integration Protocol
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R321-immediate-backport-during-integration.md`
**Criticality**: SUPREME LAW - Immediate backporting required
**Summary**: ANY fix during integration MUST be immediately backported to source branches before continuing

### 🚨🚨🚨 R280 - Main Branch Protection Protocol [BLOCKING]
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R280-main-branch-protection-protocol.md`
**Criticality**: BLOCKING - Direct commits to main/master are forbidden
**Summary**: All changes must go through PR process with proper reviews

### 🚨🚨🚨 R307 - Branch Mergeability Check [BLOCKING]
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R307-branch-mergeability-check.md`
**Criticality**: BLOCKING - Must verify branches are mergeable before attempting
**Summary**: Check for conflicts and mergeability before integration operations

### 🔴🔴🔴 R233 - All States Require Immediate Action (CRITICAL)
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R233-all-states-immediate-action.md`
**Criticality**: CRITICAL - States are verbs
**Summary**: PHASE_INTEGRATION means START INTEGRATING NOW

### 🔴🔴🔴 R288 - State File Update and Commit Protocol (SUPREME LAW)
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R288-state-file-update-and-commit-protocol.md`
**Criticality**: SUPREME LAW - Update on every transition
**Summary**: Update orchestrator-state.yaml with integration details

### 🚨🚨🚨 R288 - State File Update and Commit Protocol
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R288-state-file-update-and-commit-protocol.md`
**Criticality**: BLOCKING - Push within 60 seconds
**Summary**: Commit and push state file immediately

## 🚨 PHASE_INTEGRATION IS A VERB - START INTEGRATING IMMEDIATELY! 🚨

### IMMEDIATE ACTIONS UPON ENTERING PHASE_INTEGRATION

**THE MOMENT YOU ENTER THIS STATE, YOU MUST:**
1. Create phase-level integration branch NOW
2. Merge all wave integration branches into phase branch
3. Merge all ERROR_RECOVERY fix branches  
4. Verify all architect-identified issues are addressed
5. Push phase integration branch for re-review

**FORBIDDEN - AUTOMATIC FAILURE:**
- ❌ "STATE TRANSITION COMPLETE: Now in PHASE_INTEGRATION" [stops]
- ❌ "Successfully entered PHASE_INTEGRATION state" [waits]
- ❌ "Ready to start phase integration" [pauses]
- ❌ "I'm in PHASE_INTEGRATION state" [does nothing]

**REQUIRED - IMMEDIATE ACTION:**
- ✅ "Entering PHASE_INTEGRATION, creating phase3-integration branch NOW..."
- ✅ "START PHASE_INTEGRATION, merging all wave branches immediately..."
- ✅ "PHASE_INTEGRATION: Merging ERROR_RECOVERY fixes from fix branches..."

## Primary Purpose

Create a clean phase-level integration branch that includes:
- All wave integration branches from the phase
- All fixes from ERROR_RECOVERY addressing phase assessment issues
- Ready for architect phase reassessment
- Comprehensive integration of all phase work

## 🔴 CRITICAL: Locating Effort Branches for Integration

### Effort Branch Locations (Per R193/R191)
All effort branches are located in specific directories with predictable patterns:

#### Directory Structure:
```bash
# Effort workspaces follow this pattern:
/efforts/phase${PHASE}/wave${WAVE}/${EFFORT_NAME}/

# Example for Phase 2, Wave 1 with 3 efforts:
/efforts/phase2/wave1/auth-system/       # Contains effort branch
/efforts/phase2/wave1/user-management/   # Contains effort branch
/efforts/phase2/wave1/api-gateway/       # Contains effort branch
```

#### Branch Naming Convention:
```bash
# Effort branches follow naming from target-repo-config.yaml:
# Pattern: phase${PHASE}-wave${WAVE}-${EFFORT_NAME}
# Or with project prefix: ${PREFIX}/phase${PHASE}-wave${WAVE}-${EFFORT_NAME}

# Examples without prefix:
phase2-wave1-auth-system
phase2-wave1-user-management
phase2-wave1-api-gateway

# Examples with prefix (e.g., "tmc-workspace"):
tmc-workspace/phase2-wave1-auth-system
tmc-workspace/phase2-wave1-user-management
tmc-workspace/phase2-wave1-api-gateway
```

### Finding Efforts to Integrate

**MANDATORY: Before integration, locate all effort branches:**

```bash
#!/bin/bash
# Script to find all effort branches for current phase

PHASE=$(yq '.current_phase' orchestrator-state.yaml)
SF_INSTANCE_DIR=$(pwd)

echo "🔍 Locating effort branches for Phase ${PHASE} integration"
echo "================================================="

# Find all waves in the phase
for wave_num in $(seq 1 10); do
    WAVE_DIR="/efforts/phase${PHASE}/wave${wave_num}"
    
    if [ ! -d "$WAVE_DIR" ]; then
        continue  # Wave doesn't exist, skip
    fi
    
    echo "\n📁 Wave ${wave_num} efforts:"
    echo "-------------------"
    
    # List all effort directories in this wave
    for effort_dir in "$WAVE_DIR"/*/; do
        if [ ! -d "$effort_dir/.git" ]; then
            continue  # Not a git repository, skip
        fi
        
        EFFORT_NAME=$(basename "$effort_dir")
        cd "$effort_dir"
        
        # Get current branch and remote info
        CURRENT_BRANCH=$(git branch --show-current)
        REMOTE_URL=$(git remote get-url origin 2>/dev/null || echo "No remote")
        
        echo "  ✅ Effort: $EFFORT_NAME"
        echo "     Directory: $effort_dir"
        echo "     Branch: $CURRENT_BRANCH"
        echo "     Remote: $REMOTE_URL"
        
        # Check if branch exists on remote
        if git ls-remote --heads origin "$CURRENT_BRANCH" 2>/dev/null | grep -q "$CURRENT_BRANCH"; then
            echo "     Status: ✅ Branch exists on remote"
        else
            echo "     Status: ⚠️ Branch not found on remote!"
        fi
    done
done

cd "$SF_INSTANCE_DIR"
echo "\n================================================="
echo "📊 Effort branch discovery complete"
```

### Integration Source Requirements (R300)

**Per R300 (Comprehensive Fix Management Protocol):**
- ALL fixes must be in effort branches before integration
- Integration branches are temporary and recreated from main
- Effort branches are the SOURCE OF TRUTH that become PRs
- NEVER apply fixes directly to integration branches

### Verifying Effort Branches Before Integration

```bash
#!/bin/bash
# Verify all effort branches are ready for integration

verify_effort_branches() {
    local PHASE=$1
    local ERRORS=0
    
    echo "🔍 Verifying effort branches for Phase ${PHASE}"
    
    # Check orchestrator-state.yaml for expected efforts
    EXPECTED_EFFORTS=$(yq ".phases.phase_${PHASE}.efforts[]" orchestrator-state.yaml 2>/dev/null)
    
    if [ -z "$EXPECTED_EFFORTS" ]; then
        echo "⚠️ No efforts recorded in orchestrator-state.yaml for phase ${PHASE}"
        echo "   Searching filesystem for actual efforts..."
    fi
    
    # Scan filesystem for actual effort branches
    for wave_num in $(seq 1 10); do
        WAVE_DIR="/efforts/phase${PHASE}/wave${wave_num}"
        
        if [ ! -d "$WAVE_DIR" ]; then
            continue
        fi
        
        for effort_dir in "$WAVE_DIR"/*/; do
            if [ ! -d "$effort_dir/.git" ]; then
                continue
            fi
            
            EFFORT_NAME=$(basename "$effort_dir")
            cd "$effort_dir"
            
            # Verify branch is pushed
            CURRENT_BRANCH=$(git branch --show-current)
            if ! git ls-remote --heads origin "$CURRENT_BRANCH" 2>/dev/null | grep -q "$CURRENT_BRANCH"; then
                echo "❌ Effort '$EFFORT_NAME' branch not pushed to remote!"
                ((ERRORS++))
            else
                echo "✅ Effort '$EFFORT_NAME' ready for integration"
            fi
        done
    done
    
    cd "$SF_INSTANCE_DIR"
    
    if [ $ERRORS -gt 0 ]; then
        echo "❌ Found $ERRORS effort(s) not ready for integration"
        return 1
    else
        echo "✅ All effort branches verified and ready"
        return 0
    fi
}

# Run verification
PHASE=$(yq '.current_phase' orchestrator-state.yaml)
verify_effort_branches $PHASE
```

## 🔴🔴🔴 R321 ENFORCEMENT: IMMEDIATE BACKPORTING 🔴🔴🔴

### CRITICAL: Phase Integration Branches Are READ-ONLY
**Per R321, if ANY issue is found during phase integration:**

1. **STOP IMMEDIATELY** - Do not continue integration
2. **IDENTIFY** which wave/effort branch needs fixing
3. **SPAWN SW ENGINEER** to fix the source branch
4. **WAIT** for fix to be applied and pushed
5. **VERIFY** source branch works independently
6. **ONLY THEN** retry phase integration with fixed sources

### Validation Before Phase Integration
```bash
# R321 MANDATORY: Verify all wave branches work independently
validate_wave_branches_before_phase_integration() {
    echo "🔍 R321 Validation: Checking all wave branches"
    
    for wave_num in $(seq 1 10); do
        WAVE_DIR="/efforts/phase${PHASE}/wave${wave_num}"
        if [ ! -d "$WAVE_DIR" ]; then
            continue
        fi
        
        # Each wave must have working effort branches
        for effort_dir in "$WAVE_DIR"/*/; do
            cd "$effort_dir"
            BRANCH=$(git branch --show-current)
            
            echo "Testing $BRANCH independently..."
            if ! npm run build; then
                echo "❌ R321 VIOLATION: $BRANCH doesn't build!"
                echo "Must fix in source before phase integration"
                exit 1
            fi
            
            if ! npm test; then
                echo "❌ R321 VIOLATION: $BRANCH tests fail!"
                echo "Must fix in source before phase integration"
                exit 1
            fi
        done
    done
    
    echo "✅ R321 Validated: All wave sources work independently"
}
```

### Detection of Phase Integration Branch Modifications
```bash
# R321 ENFORCEMENT: Check for forbidden direct edits
check_phase_integration_purity() {
    NON_MERGE=$(git log --oneline --no-merges origin/main..HEAD)
    
    if [ -n "$NON_MERGE" ]; then
        echo "🔴🔴🔴 R321 VIOLATION DETECTED!"
        echo "Direct edits found in phase integration branch:"
        echo "$NON_MERGE"
        echo "ALL fixes must go to source branches!"
        exit 1
    fi
}
```

## State Context

### CRITICAL: Detect Entry Context (R285 vs R259)

You MUST determine HOW you entered PHASE_INTEGRATION:

1. **FROM WAVE_REVIEW (Normal Flow - R285)**:
   - This is standard phase completion after last wave
   - Integrate all wave branches for the phase
   - No ERROR_RECOVERY fixes to include
   - Branch name: `phase-{N}-integration`
   - Purpose: Prepare integrated phase for architect assessment

2. **FROM ERROR_RECOVERY (Fix Flow - R259)**:
   - This is after phase assessment returned NEEDS_WORK
   - Integrate all wave branches PLUS fix branches
   - Must address issues from assessment report
   - Branch name: `phase{N}-post-fixes-integration-{TIMESTAMP}`
   - Purpose: Prepare fixed phase for architect reassessment

Check previous_state in orchestrator-state.yaml to determine context!

### 🚨🚨🚨 RULE R259 - Phase Integration Requirements
**SEE**: `$CLAUDE_PROJECT_DIR/rule-library/R259-mandatory-phase-integration-after-fixes.md`

### 🚨🚨🚨 RULE R257 - Assessment Report Verification
**SEE**: `$CLAUDE_PROJECT_DIR/rule-library/R257-mandatory-phase-assessment-report.md`

## Branch Naming Convention (R014 MANDATORY)

```bash
# 🔴 CRITICAL: Use branch naming helpers for project prefix support
source "$SF_INSTANCE_DIR/utilities/branch-naming-helpers.sh"

# Determine context to choose naming convention
PREVIOUS_STATE=$(yq '.previous_state' orchestrator-state.yaml)

if [ "$PREVIOUS_STATE" = "WAVE_REVIEW" ]; then
    # Normal flow (R285) - standard phase integration
    PHASE_BRANCH=$(get_phase_integration_branch_name "$PHASE")
    # Example with prefix: tmc-workspace/phase3-integration
    # Example without: phase3-integration
elif [ "$PREVIOUS_STATE" = "ERROR_RECOVERY" ]; then
    # Fix flow (R259) - post-fixes integration with timestamp
    BASE_BRANCH=$(get_phase_integration_branch_name "$PHASE")
    PHASE_BRANCH="${BASE_BRANCH}-post-fixes-$(date +%Y%m%d-%H%M%S)"
    # Example with prefix: tmc-workspace/phase3-integration-post-fixes-20250827-153000
    # Example without: phase3-integration-post-fixes-20250827-153000
else
    echo "❌ Unexpected previous state: $PREVIOUS_STATE"
    exit 1
fi
```

## Integration Process

### 1. Create Phase Integration Infrastructure

```bash
#!/bin/bash
# Setup phase integration workspace and branch

SF_INSTANCE_DIR=$(pwd)  # Save SF instance location
PHASE=$(yq '.current_phase' orchestrator-state.yaml)
TIMESTAMP=$(date +%Y%m%d-%H%M%S)

# Source branch naming helpers (R014 MANDATORY)
source "$SF_INSTANCE_DIR/utilities/branch-naming-helpers.sh"

# Create phase integration workspace
PHASE_DIR="${CLAUDE_PROJECT_DIR}/efforts/phase${PHASE}"
INTEGRATION_DIR="${PHASE_DIR}/phase-integration-workspace"
echo "Creating phase integration workspace at: $INTEGRATION_DIR"
mkdir -p "$(dirname "$INTEGRATION_DIR")"

# Determine base branch for phase integration (R271)
echo "🧠 THINKING: Phase integration needs clean base from main branch"
BASE_BRANCH=$(yq '.target_repository.base_branch' "$SF_INSTANCE_DIR/target-repo-config.yaml")
if [ -z "$BASE_BRANCH" ] || [ "$BASE_BRANCH" = "null" ]; then
    BASE_BRANCH="main"
fi
echo "📌 Decision: Using '$BASE_BRANCH' for phase integration (clean start)"

# SINGLE-BRANCH FULL clone of TARGET repository (R271 Supreme Law)
echo "📦 Creating FULL phase integration clone from branch: $BASE_BRANCH"
TARGET_REPO_URL=$(yq '.target_repository.url' "$SF_INSTANCE_DIR/target-repo-config.yaml")

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

# Verify FULL checkout (R271 compliance)
if [ -f ".git/info/sparse-checkout" ]; then
    echo "🔴🔴🔴 SUPREME LAW VIOLATION: Sparse checkout in phase integration!"
    exit 1
fi
echo "✅ Full codebase available for phase integration from $BASE_BRANCH"

# Create phase integration branch with proper naming (R014)
BASE_BRANCH=$(get_phase_integration_branch_name "$PHASE")
BRANCH_NAME="${BASE_BRANCH}-post-fixes-${TIMESTAMP}"
echo "🔀 Creating phase integration branch: $BRANCH_NAME"
git checkout -b "$BRANCH_NAME"

# Push to establish remote tracking
git push -u origin "$BRANCH_NAME"

# R301 MANDATORY: Update current_phase_integration and deprecate old
echo "📝 Updating current_phase_integration per R301..."

# First, move any existing current phase integration to deprecated
EXISTING_PHASE=$(yq ".current_phase_integration | select(.phase == env(PHASE))" "$SF_INSTANCE_DIR/orchestrator-state.yaml")
if [ ! -z "$EXISTING_PHASE" ]; then
    yq -i '.deprecated_phase_integrations += (.current_phase_integration | select(.phase == env(PHASE)))' "$SF_INSTANCE_DIR/orchestrator-state.yaml"
fi

# Set the new current phase integration
yq -i '.current_phase_integration = {
  "phase": env(PHASE),
  "branch": env(BRANCH_NAME),
  "status": "active",
  "created_at": "'$(date -u +%Y-%m-%dT%H:%M:%SZ)'",
  "type": "post_fixes"
}' "$SF_INSTANCE_DIR/orchestrator-state.yaml"

echo "✅ Phase integration infrastructure ready: $BRANCH_NAME"
echo "✅ Current phase integration updated per R301"
```

### 2. Spawn Code Reviewer for Phase Merge Plan

```bash
#!/bin/bash
# Spawn Code Reviewer to create PHASE-MERGE-PLAN.md

PHASE=$(yq '.current_phase' orchestrator-state.yaml)
INTEGRATION_DIR="${CLAUDE_PROJECT_DIR}/efforts/phase${PHASE}/phase-integration-workspace"
BRANCH_NAME=$(git branch --show-current)

# First, collect all effort branches for Code Reviewer
EFFORT_BRANCHES=""
for wave_num in $(seq 1 10); do
    WAVE_DIR="/efforts/phase${PHASE}/wave${wave_num}"
    if [ -d "$WAVE_DIR" ]; then
        for effort_dir in "$WAVE_DIR"/*/; do
            if [ -d "$effort_dir/.git" ]; then
                cd "$effort_dir"
                BRANCH=$(git branch --show-current)
                EFFORT_BRANCHES="${EFFORT_BRANCHES}\n- ${BRANCH} (from ${effort_dir})"
            fi
        done
    fi
done

cat > /tmp/code-reviewer-phase-merge-plan-task.md << EOF
Create PHASE MERGE PLAN for Phase ${PHASE} integration.

CRITICAL REQUIREMENTS:
1. Use ONLY original wave effort branches - NO integration branches!
2. Include ALL ERROR_RECOVERY fix branches from phase assessment
3. Analyze branch bases to determine correct merge order
4. Create PHASE-MERGE-PLAN.md with exact merge instructions
5. DO NOT execute merges - only plan them!

Integration Directory: ${INTEGRATION_DIR}
Target Branch: ${BRANCH_NAME}

Effort Branches Found:${EFFORT_BRANCHES}

Phase Assessment Report: phase-assessments/phase${PHASE}/PHASE-${PHASE}-ASSESSMENT-REPORT.md
EOF

# Spawn Code Reviewer
Task: subagent_type="code-reviewer" \
      prompt="$(cat /tmp/code-reviewer-phase-merge-plan-task.md)" \
      description="Create Phase ${PHASE} Merge Plan"

echo "📋 Waiting for Code Reviewer to create PHASE-MERGE-PLAN.md..."
```

### 3. Spawn Integration Agent for Phase Integration

```bash
#!/bin/bash
# After Code Reviewer creates PHASE-MERGE-PLAN.md

PHASE=$(yq '.current_phase' orchestrator-state.yaml)
INTEGRATION_DIR="${CLAUDE_PROJECT_DIR}/efforts/phase${PHASE}/phase-integration-workspace"

# CD into phase integration directory
cd "$INTEGRATION_DIR"

# Verify merge plan exists
if [ ! -f "PHASE-MERGE-PLAN.md" ]; then
    echo "❌ Cannot spawn Integration Agent - no phase merge plan!"
    exit 1
fi

# Spawn Integration Agent
Task: subagent_type="integration-agent" \
      prompt="Execute phase integration merges for Phase ${PHASE}.
      
      CRITICAL REQUIREMENTS:
      1. You are in INTEGRATION_DIR: ${INTEGRATION_DIR}
      2. Acknowledge and set INTEGRATION_DIR variable
      3. Read and follow PHASE-MERGE-PLAN.md EXACTLY
      4. Execute merges in specified order
      5. Handle conflicts as directed in plan
      6. Run phase-level tests after all merges
      
      Your working directory has been set to: ${INTEGRATION_DIR}
      The merge plan is: PHASE-MERGE-PLAN.md" \
      description="Execute Phase ${PHASE} Integration"

echo "🎯 Integration Agent spawned for phase integration"
```

### 4. Verify Against Assessment Report

```bash
#!/bin/bash
# Verify all issues from phase assessment are addressed

verify_phase_integration() {
    local PHASE=$1
    local REPORT="phase-assessments/phase${PHASE}/PHASE-${PHASE}-ASSESSMENT-REPORT.md"
    
    if [ ! -f "$REPORT" ]; then
        echo "❌ Phase assessment report not found!"
        return 1
    fi
    
    echo "📋 Verifying fixes against assessment report..."
    
    # Extract Priority 1 issues
    PRIORITY_1=$(sed -n '/### Priority 1/,/### Priority 2/p' "$REPORT" | grep "^- \[" | wc -l)
    
    echo "📊 Assessment requirements:"
    echo "  - Priority 1 (Must Fix): $PRIORITY_1 items"
    
    # Check for evidence of fixes
    echo "🔍 Checking for fix evidence in commits..."
    
    # Look for commits referencing the assessment
    FIX_COMMITS=$(git log --oneline | grep -i "phase.*assess\|priority.*1\|must.*fix" | wc -l)
    
    echo "  - Found $FIX_COMMITS commits referencing fixes"
    
    if [ $FIX_COMMITS -ge $PRIORITY_1 ]; then
        echo "✅ Sufficient fix commits found"
        return 0
    else
        echo "⚠️ May need additional verification"
        return 1
    fi
}

PHASE=$(yq '.current_phase' orchestrator-state.yaml)
verify_phase_integration $PHASE
```

### 5. 🏗️ MANDATORY BUILD VERIFICATION

**NO PHASE INTEGRATION IS COMPLETE WITHOUT A WORKING BUILD:**

```bash
#!/bin/bash
# Mandatory build verification for phase integration

verify_phase_build() {
    local PHASE=$1
    
    echo "🏗️ Running phase-level build verification..."
    
    # Clean previous builds
    rm -rf dist/ build/ out/
    
    # Run the build
    if npm run build 2>&1 | tee "phase${PHASE}-build.log"; then
        echo "✅ Phase build successful"
    else
        echo "❌ Phase build failed - integration incomplete!"
        return 1
    fi
    
    # Verify build artifacts
    if [ -d "dist" ] || [ -d "build" ] || [ -d "out" ]; then
        echo "✅ Build artifacts created"
        ls -la dist/ build/ out/ 2>/dev/null | tee "phase${PHASE}-artifacts.log"
    else
        echo "❌ No build artifacts found!"
        return 1
    fi
    
    # Verify executable/runnable
    if [ -f "dist/index.js" ] || [ -f "build/main" ] || [ -f "out/app.jar" ]; then
        echo "✅ Executable artifacts verified"
    else
        echo "⚠️ Warning: No executable found - verify manually"
    fi
    
    return 0
}

PHASE=$(yq '.current_phase' orchestrator-state.yaml)
verify_phase_build $PHASE || exit 1
```

### 6. 🧪 MANDATORY TEST HARNESS

**EVERY PHASE INTEGRATION MUST HAVE A COMPREHENSIVE TEST HARNESS:**

```bash
#!/bin/bash
# Create and run phase-level test harness

create_phase_test_harness() {
    local PHASE=$1
    
    cat > "phase${PHASE}-test-harness.sh" << 'EOF'
#!/bin/bash
# Phase Integration Test Harness
echo "🧪 Starting Phase Integration Test Suite"
echo "========================================="

PHASE=${1:-1}

# Unit tests
echo "📦 Running ALL unit tests..."
if npm test 2>&1 | tee "phase${PHASE}-unit-tests.log"; then
    echo "✅ Unit tests passed"
else
    echo "❌ Unit tests failed"
    exit 1
fi

# Integration tests
echo "🔗 Running integration tests..."
if npm run test:integration 2>&1 | tee "phase${PHASE}-integration-tests.log"; then
    echo "✅ Integration tests passed"
else
    echo "❌ Integration tests failed"
    exit 1
fi

# Phase-specific tests
echo "🎯 Running phase-specific tests..."
if [ -f "tests/phase${PHASE}/run-all.sh" ]; then
    bash "tests/phase${PHASE}/run-all.sh" 2>&1 | tee "phase${PHASE}-specific-tests.log"
else
    echo "ℹ️ No phase-specific tests found"
fi

# Regression tests
echo "🔄 Running regression tests..."
if npm run test:regression 2>&1 | tee "phase${PHASE}-regression-tests.log"; then
    echo "✅ Regression tests passed"
else
    echo "⚠️ Regression test issues detected"
fi

echo "========================================="
echo "✅ PHASE TEST SUITE COMPLETE!"
echo "Results saved to phase${PHASE}-*-tests.log files"
EOF
    
    chmod +x "phase${PHASE}-test-harness.sh"
    
    # Run the test harness
    "./phase${PHASE}-test-harness.sh" "$PHASE"
}

PHASE=$(yq '.current_phase' orchestrator-state.yaml)
create_phase_test_harness $PHASE
```

### 7. 🎬 MANDATORY PHASE DEMO

**DEMONSTRATE THE WORKING PHASE INTEGRATION:**

```bash
#!/bin/bash
# Create comprehensive phase demo

create_phase_demo() {
    local PHASE=$1
    
    # Create demo documentation
    cat > "PHASE-${PHASE}-DEMO.md" << EOF
# Phase ${PHASE} Integration Demo

## Build & Test Status
- Build: ✅ PASSING (see phase${PHASE}-build.log)
- Unit Tests: ✅ ALL PASSING
- Integration Tests: ✅ ALL PASSING
- Phase Tests: ✅ COMPLETE
- Test Harness: phase${PHASE}-test-harness.sh

## Integrated Waves
$(yq ".phases.phase_${PHASE}.waves[]" orchestrator-state.yaml | sed 's/^/- Wave /')

## Features Delivered in Phase ${PHASE}
$(grep "^- " "phase-plans/phase${PHASE}/PHASE-PLAN.md" | head -10)

## How to Run Full Demo
\`\`\`bash
# 1. Start the application
npm start

# 2. Run the automated demo
./phase${PHASE}-demo.sh

# 3. Verify all features
./verify-phase${PHASE}-features.sh
\`\`\`

## Manual Verification Steps
1. Open application at http://localhost:3000
2. Navigate to [Feature Area]
3. Verify [Specific Functionality]
4. Check logs for successful operations

## Evidence & Artifacts
- Build Log: phase${PHASE}-build.log
- Build Artifacts: phase${PHASE}-artifacts.log
- Test Results: phase${PHASE}-*-tests.log
- Demo Script: phase${PHASE}-demo.sh
- Screenshots: demos/phase${PHASE}/

## Phase Metrics
- Total Lines Added: $(git diff main --numstat | awk '{sum+=$1} END {print sum}')
- Test Coverage: $(npm run coverage --silent | grep "All files" | awk '{print $10}')
- Features Completed: $(yq ".phases.phase_${PHASE}.features_completed" orchestrator-state.yaml)
EOF
    
    # Create automated demo script
    cat > "phase${PHASE}-demo.sh" << 'EOF'
#!/bin/bash
echo "🎬 Phase Integration Demo"
echo "========================="
echo "Demonstrating all integrated functionality from Phase ${1:-1}"
echo ""

# Start application if needed
if ! curl -s http://localhost:3000/health > /dev/null 2>&1; then
    echo "Starting application..."
    npm start &
    sleep 5
fi

# Demo each major feature
echo "📍 Feature 1: [Name]"
# Add actual demo commands
curl -X GET http://localhost:3000/api/feature1
echo ""

echo "📍 Feature 2: [Name]"
# Add actual demo commands
curl -X POST http://localhost:3000/api/feature2 -d '{"test": "data"}'
echo ""

echo "📍 Feature 3: [Name]"
# Add actual demo commands
curl -X GET http://localhost:3000/api/feature3
echo ""

echo "========================="
echo "✅ Phase Demo Complete!"
EOF
    
    chmod +x "phase${PHASE}-demo.sh"
    
    # Create feature verification script
    cat > "verify-phase${PHASE}-features.sh" << 'EOF'
#!/bin/bash
echo "🔍 Verifying Phase Features"
FAILURES=0

# Check each feature endpoint/functionality
echo -n "Checking Feature 1... "
if curl -s http://localhost:3000/api/feature1 | grep -q "expected"; then
    echo "✅"
else
    echo "❌"
    ((FAILURES++))
fi

echo -n "Checking Feature 2... "
if curl -s http://localhost:3000/api/feature2 | grep -q "expected"; then
    echo "✅"
else
    echo "❌"
    ((FAILURES++))
fi

if [ $FAILURES -eq 0 ]; then
    echo "✅ All features verified!"
    exit 0
else
    echo "❌ $FAILURES features failed verification"
    exit 1
fi
EOF
    
    chmod +x "verify-phase${PHASE}-features.sh"
    
    echo "📄 Phase demo artifacts created:"
    echo "  - PHASE-${PHASE}-DEMO.md"
    echo "  - phase${PHASE}-demo.sh"
    echo "  - verify-phase${PHASE}-features.sh"
}

PHASE=$(yq '.current_phase' orchestrator-state.yaml)
create_phase_demo $PHASE
```

### 6. Create Integration Summary

```bash
#!/bin/bash
# Create comprehensive phase integration summary

create_phase_integration_summary() {
    local PHASE=$1
    local BRANCH=$2
    local TIMESTAMP=$(date -u +%Y-%m-%dT%H:%M:%SZ)
    
    cat > "phase-integrations/phase${PHASE}/integration-summary.md" << EOF
# Phase ${PHASE} Integration Summary

**Integration Branch:** ${BRANCH}
**Created At:** ${TIMESTAMP}
**Integration Type:** Post-Assessment-Fixes

## Included Components

### Wave Integration Branches
$(git log --oneline | grep "Merge wave integration" | sed 's/^/- /')

### ERROR_RECOVERY Fix Branches  
$(git log --oneline | grep "Merge ERROR_RECOVERY" | sed 's/^/- /')

## Issues Addressed

### From Phase Assessment Report
$(grep "Priority 1" -A10 "phase-assessments/phase${PHASE}/PHASE-${PHASE}-ASSESSMENT-REPORT.md" | grep "^- \[")

## Validation Results

- ✅ All wave branches merged successfully
- ✅ All fix branches integrated
- ✅ No merge conflicts remaining
- ✅ Phase-level tests passing
- ✅ Ready for architect reassessment

## Next Steps

1. Push integration branch to remote
2. Transition to SPAWN_ARCHITECT_PHASE_ASSESSMENT
3. Await architect reassessment with fixes integrated
EOF
    
    echo "📄 Integration summary created"
}
```

## State Tracking Updates (R301 Compliant)

```yaml
# Update orchestrator-state.yaml per R301
current_phase_integration:
  phase: 3
  branch: "phase3-post-fixes-integration-20250827-143000"
  status: "active"  # MUST be "active" for current
  created_at: "2025-08-27T14:30:00Z"
  type: "post_fixes"
  metadata:
    includes_waves: [1, 2, 3, 4]
    includes_fixes: 
      - "phase3-fix-kcp-patterns-20250827-120000"
      - "phase3-fix-api-compatibility-20250827-130000"  
      - "phase3-fix-test-coverage-20250827-135000"
    original_assessment_report: "phase-assessments/phase3/PHASE-3-ASSESSMENT-REPORT.md"
    assessment_score_before_fixes: 68
    ready_for_reassessment: true

# Move previous integration to deprecated
deprecated_integrations:
  - phase: 3
    branch: "phase3-integration-20250827-100000"
    status: "deprecated"
    deprecated_at: "2025-08-27T14:30:00Z"
    reason: "superseded by post-fixes integration"
    
error_recovery_completed:
  - phase: 3
    recovery_type: "PHASE_ASSESSMENT_NEEDS_WORK"
    completed_at: "2025-08-27T14:25:00Z"
    fixes_applied: 5
  integration_branch_created: false  # Will be true after PHASE_INTEGRATION
```

## Validation Functions

```python
def validate_phase_integration_branch():
    """Validate phase integration branch is ready for reassessment"""
    
    phase = read_yaml('orchestrator-state.yaml')['current_phase']
    
    # Check branch exists
    branch_pattern = f"phase{phase}-post-fixes-integration-*"
    branches = subprocess.check_output(
        f"git branch -r | grep '{branch_pattern}'",
        shell=True
    ).decode().strip()
    
    if not branches:
        return {
            'valid': False,
            'error': 'No phase integration branch found'
        }
    
    # Verify all waves included
    wave_count = read_yaml('orchestrator-state.yaml')['waves_per_phase'][phase-1]
    
    for wave in range(1, wave_count + 1):
        wave_branch = f"phase{phase}-wave{wave}-integration"
        merge_found = subprocess.call(
            f"git log --oneline | grep -q 'Merge.*{wave_branch}'",
            shell=True
        ) == 0
        
        if not merge_found:
            return {
                'valid': False,
                'error': f'Wave {wave} not integrated'
            }
    
    # Verify ERROR_RECOVERY fixes merged
    fixes_found = subprocess.call(
        "git log --oneline | grep -q 'ERROR_RECOVERY fixes'",
        shell=True
    ) == 0
    
    if not fixes_found:
        return {
            'valid': False,
            'error': 'ERROR_RECOVERY fixes not integrated'
        }
    
    return {
        'valid': True,
        'branch': branches.split('\n')[0].strip(),
        'ready_for_reassessment': True
    }
```

## State Transitions

From PHASE_INTEGRATION state:
- **SUCCESS** → SPAWN_ARCHITECT_PHASE_ASSESSMENT (reassessment with integrated fixes)
- **MERGE_CONFLICTS** → ERROR_RECOVERY (resolve conflicts)
- **TEST_FAILURES** → ERROR_RECOVERY (fix test issues)
- **VALIDATION_FAILED** → ERROR_RECOVERY (missing fixes)

To PHASE_INTEGRATION state:
- **ERROR_RECOVERY** → PHASE_INTEGRATION (after phase assessment fixes complete)

## Critical Success Criteria

Before transitioning to SPAWN_ARCHITECT_PHASE_ASSESSMENT:
1. ✅ Phase integration branch created and pushed
2. ✅ All wave integration branches merged
3. ✅ All ERROR_RECOVERY fix branches merged
4. ✅ No unresolved merge conflicts
5. ✅ Phase-level tests passing
6. ✅ Integration summary document created
7. ✅ State file updated with integration details
8. ✅ Original assessment report issues verified as addressed

## Common Mistakes to Avoid

1. **Not finding effort branches**
   - ❌ WRONG: Looking for branches in wrong location
   - ❌ WRONG: Assuming branches are in SF instance directory
   - ✅ RIGHT: Check `/efforts/phase${PHASE}/wave${WAVE}/${EFFORT_NAME}/`
   - ✅ RIGHT: Verify each effort directory has a git repository

2. **Creating branch from wrong base**
   - ❌ WRONG: Branch from current work branch
   - ✅ RIGHT: Branch from clean main

3. **Missing wave integrations**
   - ❌ WRONG: Only merge some waves
   - ✅ RIGHT: Merge ALL waves from the phase

4. **Skipping fix branches**
   - ❌ WRONG: Forget ERROR_RECOVERY fixes
   - ✅ RIGHT: Include all fix branches

5. **Not verifying against report**
   - ❌ WRONG: Assume fixes are complete
   - ✅ RIGHT: Verify each issue from assessment report

6. **Wrong state transition**
   - ❌ WRONG: Go directly to PHASE_COMPLETE
   - ✅ RIGHT: Go to SPAWN_ARCHITECT_PHASE_ASSESSMENT for reassessment

## Integration Checklist

- [ ] Phase integration branch created from main
- [ ] All wave integration branches identified
- [ ] Each wave branch merged successfully  
- [ ] All fix branches identified from ERROR_RECOVERY
- [ ] Each fix branch merged successfully
- [ ] Merge conflicts resolved if any
- [ ] Phase-level tests executed and passing
- [ ] Integration summary document created
- [ ] State file updated with branch details
- [ ] Branch pushed to remote repository
- [ ] Ready to spawn architect for reassessment

## Concrete Example: Phase 2 Integration

```bash
# Example: Integrating Phase 2 with 2 waves, 5 total efforts

# Step 1: Discover all effort branches
echo "Finding Phase 2 effort branches..."

# Wave 1 efforts:
ls -la /efforts/phase2/wave1/
# Output:
# drwxr-xr-x auth-system/        # effort branch: phase2-wave1-auth-system
# drwxr-xr-x user-management/    # effort branch: phase2-wave1-user-management
# drwxr-xr-x api-gateway/        # effort branch: phase2-wave1-api-gateway

# Wave 2 efforts:
ls -la /efforts/phase2/wave2/
# Output:
# drwxr-xr-x database-layer/     # effort branch: phase2-wave2-database-layer
# drwxr-xr-x cache-service/      # effort branch: phase2-wave2-cache-service

# Step 2: Verify each effort branch
for effort_dir in /efforts/phase2/wave*/*/; do
    cd "$effort_dir"
    echo "Effort: $(basename $effort_dir)"
    git branch --show-current
    git log --oneline -1
done

# Step 3: Create phase integration branch
cd ${CLAUDE_PROJECT_DIR}/efforts/phase2/phase-integration-workspace
git checkout main
git pull origin main
git checkout -b "phase2-integration-20250901-143000"

# Step 4: Merge all effort branches
git merge origin/phase2-wave1-auth-system --no-ff -m "feat: merge auth-system effort"
git merge origin/phase2-wave1-user-management --no-ff -m "feat: merge user-management effort"
git merge origin/phase2-wave1-api-gateway --no-ff -m "feat: merge api-gateway effort"
git merge origin/phase2-wave2-database-layer --no-ff -m "feat: merge database-layer effort"
git merge origin/phase2-wave2-cache-service --no-ff -m "feat: merge cache-service effort"

# Step 5: Push integrated branch
git push -u origin phase2-integration-20250901-143000
```

## Quick Reference

```bash
# Essential commands for PHASE_INTEGRATION
PHASE=$(yq '.current_phase' orchestrator-state.yaml)
BRANCH="phase${PHASE}-post-fixes-integration-$(date +%Y%m%d-%H%M%S)"

# Find all effort branches
for wave_dir in /efforts/phase${PHASE}/wave*/; do
    for effort_dir in "$wave_dir"/*/; do
        if [ -d "$effort_dir/.git" ]; then
            cd "$effort_dir"
            echo "Found: $(git branch --show-current) in $effort_dir"
        fi
    done
done

# Create integration branch
cd ${CLAUDE_PROJECT_DIR}/efforts/phase${PHASE}/phase-integration-workspace
git checkout main && git pull
git checkout -b "$BRANCH"

# Merge effort branches (NOT wave integration branches!)
for wave_dir in /efforts/phase${PHASE}/wave*/; do
    for effort_dir in "$wave_dir"/*/; do
        if [ -d "$effort_dir/.git" ]; then
            cd "$effort_dir"
            EFFORT_BRANCH=$(git branch --show-current)
            cd ${CLAUDE_PROJECT_DIR}/efforts/phase${PHASE}/phase-integration-workspace
            git merge "origin/$EFFORT_BRANCH" --no-ff
        fi
    done
done

# Push for reassessment
git push -u origin "$BRANCH"

# Update state and transition
yq -i ".phase_integration_branches += [{\"phase\": $PHASE, \"branch\": \"$BRANCH\"}]" orchestrator-state.yaml
yq -i '.current_state = "SPAWN_ARCHITECT_PHASE_ASSESSMENT"' orchestrator-state.yaml
```

### 🔴🔴🔴 RULE R233 - States Are Verbs (CRITICAL)
**SEE**: `$CLAUDE_PROJECT_DIR/rule-library/R233-all-states-immediate-action.md`

## R322 VIOLATION DETECTION

If you find yourself:
- Starting work for a new state without /continue-orchestrating
- Transitioning without stopping after state file commit
- Continuing after completing state work

**STOP IMMEDIATELY - You are violating R322!**
