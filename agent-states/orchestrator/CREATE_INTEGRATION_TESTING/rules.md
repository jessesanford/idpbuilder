# Orchestrator - CREATE_INTEGRATION_TESTING State Rules

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


## 🔴🔴🔴 STOP! STATE RULE READING IS ABSOLUTELY FIRST! 🔴🔴🔴

**YOU HAVE ENTERED CREATE_INTEGRATION_TESTING STATE - YOU MUST READ AND ACKNOWLEDGE ALL STATE RULES BEFORE DOING ANY STATE WORK!**

## 🔴🔴🔴 R290 VERIFICATION REQUIREMENT 🔴🔴🔴

**R290 ENFORCEMENT: CREATE VERIFICATION MARKER AFTER READING**

After reading and acknowledging all state rules, you MUST create a verification marker:

```bash
# MANDATORY: Create verification marker after reading rules
touch .state_rules_read_orchestrator_CREATE_INTEGRATION_TESTING
echo "$(date +%s) - Rules read and acknowledged for CREATE_INTEGRATION_TESTING" > .state_rules_read_orchestrator_CREATE_INTEGRATION_TESTING
```

**FAILURE TO CREATE MARKER = AUTOMATIC -100% PENALTY**

The system will check for this marker. No marker = Immediate failure.

### ❌ DO NOT DO ANY INTEGRATION TESTING SETUP WORK UNTIL RULES ARE READ:
- ❌ Start creating integration-testing branch
- ❌ Start cloning repositories
- ❌ Start setting up workspaces
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
   ❌ WRONG: "I acknowledge R272, R271, R273..."
   (No Read tool calls detected = AUTOMATIC FAILURE)
   ```

2. **Bulk Acknowledgment**:
   ```
   ❌ WRONG: "I acknowledge all CREATE_INTEGRATION_TESTING rules"
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
   ❌ WRONG: "I know R272 requires creating from main HEAD..."
   (Must READ from file, not recall from memory)
   ```

5. **Skipping Rules in PRIMARY DIRECTIVES**:
   ```
   ❌ WRONG: Reading only some rules from the list
   (ALL rules in PRIMARY DIRECTIVES are MANDATORY)
   ```

### ✅ CORRECT PATTERN FOR CREATE_INTEGRATION_TESTING:
```
1. READ: $CLAUDE_PROJECT_DIR/rule-library/[first-rule-file].md
2. "I acknowledge [Rule#] - [Rule Name]: [Brief description]"
3. READ: $CLAUDE_PROJECT_DIR/rule-library/[second-rule-file].md  
4. "I acknowledge [Rule#] - [Rule Name]: [Brief description]"
[Continue for EVERY rule in PRIMARY DIRECTIVES...]
5. "Ready to execute CREATE_INTEGRATION_TESTING work"
```

### 🚨 NO WORK UNTIL ACKNOWLEDGMENT COMPLETE 🚨
**You may NOT begin ANY CREATE_INTEGRATION_TESTING work until:**
1. ✅ ALL rules in PRIMARY DIRECTIVES have been READ
2. ✅ ALL rules have been individually ACKNOWLEDGED
3. ✅ You have stated readiness to execute CREATE_INTEGRATION_TESTING work

### 🚨 FAILURE TO READ STATE RULES FIRST = IMMEDIATE EXIT 🚨
**If you do ANY CREATE_INTEGRATION_TESTING work before reading and acknowledging rules:**
- **STOP ALL WORK IMMEDIATELY**
- **EXIT WITH FAILURE STATUS**
- **YOU HAVE VIOLATED STATE COMPLIANCE**

**THE SYSTEM IS MONITORING YOUR READ TOOL CALLS!**

## 📋 PRIMARY DIRECTIVES FOR CREATE_INTEGRATION_TESTING STATE

### 🚨🚨🚨 R006 - Orchestrator NEVER Writes Code [BLOCKING]
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R006-orchestrator-never-writes-code.md`
**Criticality**: BLOCKING - Any code operation = -100% IMMEDIATE FAILURE
**Summary**: Orchestrator is a COORDINATOR ONLY - never writes, edits, or modifies code

**⚠️ R006 WARNING FOR CREATE_INTEGRATION_TESTING STATE:**
- DO NOT create test files yourself!
- DO NOT write integration test code!
- DO NOT modify any source files!
- You only set up infrastructure - SW Engineers write ALL code

### 🚨🚨🚨 R272 - Integration Testing Branch Requirement
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R272-integration-testing-branch.md`
**Criticality**: BLOCKING - Must create from main HEAD
**Summary**: Create dedicated integration-testing branch from main's current HEAD

### 🚨🚨🚨 R271 - Mandatory Production-Ready Validation
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R271-mandatory-production-ready-validation.md`
**Criticality**: BLOCKING - Full checkouts required
**Summary**: Must use full repository clones, no sparse checkouts

### 🚨🚨🚨 R014 - Branch Naming Convention
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R014-branch-naming-convention.md`
**Criticality**: BLOCKING - Project prefix required
**Summary**: Use project prefix for integration-testing branch

### 🚨🚨🚨 R251 - Repository Separation Law
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R251-REPOSITORY-SEPARATION-LAW.md`
**Criticality**: BLOCKING - SF/target repo isolation
**Summary**: Integration testing happens in target repository only

### 🚨🚨🚨 R280 - Main Branch Protection
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R280-main-branch-protection.md`
**Criticality**: SUPREME LAW - Never modify main
**Summary**: Software Factory NEVER pushes to main branch

## 🚨 CREATE_INTEGRATION_TESTING IS A VERB - CREATE INFRASTRUCTURE NOW! 🚨

### IMMEDIATE ACTIONS UPON ENTERING CREATE_INTEGRATION_TESTING

**THE MOMENT YOU ENTER THIS STATE, YOU MUST:**
1. Checkout main branch and pull latest
2. Create timestamped integration-testing branch from main HEAD
3. Setup integration testing workspace
4. Document branch creation details
5. Transition to INTEGRATION_TESTING state

**FORBIDDEN - AUTOMATIC FAILURE:**
- ❌ "STATE TRANSITION COMPLETE: Now in CREATE_INTEGRATION_TESTING" [stops]
- ❌ "Successfully entered CREATE_INTEGRATION_TESTING state" [waits]
- ❌ "Ready to create integration testing branch" [pauses]
- ❌ "I'm in CREATE_INTEGRATION_TESTING state" [does nothing]

**REQUIRED - IMMEDIATE ACTION:**
- ✅ "Entering CREATE_INTEGRATION_TESTING, creating branch from main HEAD NOW..."
- ✅ "START INTEGRATION TESTING SETUP, pulling latest main..."
- ✅ "CREATE_INTEGRATION_TESTING: Creating timestamped branch..."

## State Context
You are creating the final integration testing infrastructure where the PROJECT INTEGRATION branch will be merged to prove the entire project works. This is the ultimate validation before generating the MASTER-PR-PLAN for humans.

## 🔴🔴🔴 CRITICAL: INTEGRATION TESTING USES PROJECT INTEGRATION (R283) 🔴🔴🔴

**Prerequisites:**
- PROJECT_INTEGRATION has been completed (R283)
- All phases have been merged into project integration branch
- Code Reviewer has validated the project integration
- Ready for final integration testing

**This is NOT merging individual efforts - this is PROJECT-WIDE:**
1. **START** from main branch's current HEAD (not any integration branch)
2. **CREATE** timestamped integration-testing branch
3. **PREPARE** to merge the PROJECT INTEGRATION branch (which contains ALL phases)
4. **VALIDATE** entire project functionality
5. **PROVE** software is production-ready

**MERGE APPROACH (R283 compliant):**
- ✅ Merge the single project-integration branch (contains all phases)
- ❌ DO NOT merge individual effort branches
- ❌ DO NOT merge phase integration branches separately
- ❌ DO NOT merge wave integration branches

**ALWAYS BASE ON:**
- ✅ main branch HEAD (current state)
- ✅ Then merge project-integration branch into it

## Integration Testing Infrastructure Setup

### 🚨🚨🚨 R272 Compliance - Create From Main HEAD
**SEE**: `$CLAUDE_PROJECT_DIR/rule-library/R272-integration-testing-branch.md`

```bash
# 🔴 CRITICAL: Integration testing setup script
create_integration_testing_infrastructure() {
    echo "🏭 CREATE_INTEGRATION_TESTING: Starting infrastructure setup..."
    
    # 0. Save SF instance directory
    SF_INSTANCE_DIR=$(pwd)
    echo "📁 SF Instance: $SF_INSTANCE_DIR"
    
    # 1. Source branch naming helpers
    source "$SF_INSTANCE_DIR/utilities/branch-naming-helpers.sh"
    PROJECT_PREFIX=$(jq '.branch_naming.project_prefix' "$SF_INSTANCE_DIR/target-repo-config.yaml")
    
    # 2. Get target repository configuration
    TARGET_REPO_URL=$(jq '.target_repository.url' "$SF_INSTANCE_DIR/target-repo-config.yaml")
    TARGET_REPO_NAME=$(basename "$TARGET_REPO_URL" .git)
    
    # 3. Create integration testing workspace
    TIMESTAMP=$(date +%Y%m%d-%H%M%S)
    INTEGRATION_DIR="${SF_INSTANCE_DIR}/integration-testing-${TIMESTAMP}"
    echo "📦 Creating integration testing workspace: $INTEGRATION_DIR"
    
    # 4. Clone target repository from main HEAD (R271 - full clone)
    echo "🔄 Cloning target repository from main branch..."
    git clone \
        --single-branch \
        --branch main \
        "$TARGET_REPO_URL" \
        "$INTEGRATION_DIR"
    
    if [ $? -ne 0 ]; then
        echo "❌ Failed to clone target repository!"
        exit 1
    fi
    
    cd "$INTEGRATION_DIR"
    
    # 5. Verify full checkout (R271 compliance)
    if [ -f ".git/info/sparse-checkout" ]; then
        echo "🔴🔴🔴 SUPREME LAW VIOLATION: Sparse checkout detected!"
        exit 1
    fi
    echo "✅ Full repository checkout confirmed"
    
    # 6. Create integration-testing branch with project prefix
    INTEGRATION_BRANCH="${PROJECT_PREFIX}integration-testing-${TIMESTAMP}"
    echo "🌿 Creating integration testing branch: $INTEGRATION_BRANCH"
    git checkout -b "$INTEGRATION_BRANCH"
    
    # 7. Document branch creation (R272 requirement)
    cat > INTEGRATION-INFO.md << EOF
# Integration Testing Branch Information

## Branch Details
- **Branch Name**: $INTEGRATION_BRANCH
- **Created From**: main @ $(git rev-parse HEAD)
- **Created At**: $(date -u +%Y-%m-%dT%H:%M:%SZ)
- **SF Instance**: $SF_INSTANCE_DIR
- **Target Repository**: $TARGET_REPO_URL

## Purpose
This branch serves as the final integration testing ground for ALL Software Factory efforts.
It will receive merges from all phase and wave efforts to validate the complete project.

## Integration Sequence
All effort branches will be merged in dependency order:
1. Phase 1 efforts (foundation)
2. Phase 2 efforts (dependent on Phase 1)
3. Phase 3 efforts (if applicable)
... continuing for all phases

## Validation Protocol
After each merge:
- Build validation
- Test execution
- Conflict resolution documentation
- Feature verification

## Important Notes
- This branch is ephemeral - not pushed to origin main
- Used only to prove integration works
- Basis for MASTER-PR-PLAN generation
- Humans will execute actual PRs to main
EOF
    
    git add INTEGRATION-INFO.md
    git commit -m "doc: initialize integration testing branch documentation"
    
    # 8. Update orchestrator state
    cd "$SF_INSTANCE_DIR"
    jq '.integration_testing.branch = \"$INTEGRATION_BRANCH\"' orchestrator-state.json > tmp.json && mv tmp.json orchestrator-state.json
    jq '.integration_testing.workspace = \"$INTEGRATION_DIR\"' orchestrator-state.json > tmp.json && mv tmp.json orchestrator-state.json
    jq '.integration_testing.created_at = \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\"' orchestrator-state.json > tmp.json && mv tmp.json orchestrator-state.json
    jq '.integration_testing.base_commit = \"$(cd $INTEGRATION_DIR && git rev-parse HEAD)\"' orchestrator-state.json > tmp.json && mv tmp.json orchestrator-state.json
    
    # 9. Prepare effort inventory
    echo "📊 Preparing effort inventory for integration..."
    cat > "$INTEGRATION_DIR/INTEGRATION-PLAN.md" << 'EOF'
# Integration Testing Plan

## Effort Inventory
[To be populated with all effort branches]

## Merge Sequence
[To be determined based on dependencies]

## Validation Checkpoints
- [ ] All Phase 1 efforts merged
- [ ] All Phase 2 efforts merged  
- [ ] All Phase 3 efforts merged
- [ ] Build successful
- [ ] All tests passing
- [ ] Feature demos complete

## Conflict Resolution Log
[Conflicts will be documented here]
EOF
    
    echo "✅ Integration testing infrastructure created successfully!"
    echo "📁 Workspace: $INTEGRATION_DIR"
    echo "🌿 Branch: $INTEGRATION_BRANCH"
    echo "📋 Next: Transitioning to INTEGRATION_TESTING state..."
}

# Execute immediately upon state entry
create_integration_testing_infrastructure
```

## Gathering Effort Branches

```bash
# Collect all effort branches for integration
gather_effort_branches() {
    echo "📊 Gathering all effort branches for integration..."
    
    EFFORTS_FILE="$INTEGRATION_DIR/ALL-EFFORTS.txt"
    > "$EFFORTS_FILE"
    
    # Iterate through all phases and waves
    for phase_dir in $SF_INSTANCE_DIR/efforts/phase*/; do
        [ -d "$phase_dir" ] || continue
        PHASE=$(basename "$phase_dir")
        
        for wave_dir in $phase_dir/wave*/; do
            [ -d "$wave_dir" ] || continue
            WAVE=$(basename "$wave_dir")
            
            for effort_dir in $wave_dir/effort-*/; do
                [ -d "$effort_dir" ] || continue
                EFFORT=$(basename "$effort_dir")
                
                # Get branch name from effort
                if [ -d "$effort_dir/.git" ]; then
                    cd "$effort_dir"
                    BRANCH=$(git branch --show-current)
                    echo "$PHASE/$WAVE/$EFFORT:$BRANCH" >> "$EFFORTS_FILE"
                fi
            done
        done
    done
    
    echo "📋 Found $(wc -l < $EFFORTS_FILE) effort branches to integrate"
}
```

## ⚠️ CRITICAL REQUIREMENTS

### 🔴 R006 VIOLATION DETECTION 
If you find yourself about to:
- Create test files with code
- Write integration test implementations
- Modify any source files
- Add test cases directly
- Make excuses like "just setting up test infrastructure"

**STOP IMMEDIATELY - You are violating R006!**
You create directories and infrastructure ONLY - SW Engineers write ALL code!

### Integration Testing Setup Rules
- Create branch from main HEAD only
- Use full repository clones (no sparse)
- Set up directory structure only
- Document what tests are needed
- Spawn SW Engineers for ALL code writing

## 🚫 FORBIDDEN ACTIONS

1. **NEVER write any test code yourself** - R006 VIOLATION = -100%
2. **NEVER create test files with implementations** - R006 VIOLATION = -100%
3. **NEVER modify source files** - R006 VIOLATION = -100%
4. **NEVER use sparse checkouts for integration**
5. **NEVER create branch from non-HEAD commits**
6. **NEVER proceed without full workspace setup**

## State Transitions

From CREATE_INTEGRATION_TESTING state:
- **SUCCESS** → INTEGRATION_TESTING (branch and workspace ready)
- **FAILURE** → ERROR_RECOVERY (setup failed)

The next state (INTEGRATION_TESTING) will:
1. Merge all effort branches in dependency order
2. Validate after each merge
3. Document all conflicts and resolutions
4. Run comprehensive test suite
5. Verify production readiness

## 🏗️ MANDATORY VERIFICATION BEFORE TRANSITION

**Before transitioning to INTEGRATION_TESTING:**
1. ✅ Integration-testing branch created from main HEAD
2. ✅ Full repository clone (no sparse checkout)
3. ✅ Branch has project prefix
4. ✅ INTEGRATION-INFO.md committed
5. ✅ Workspace path recorded in state
6. ✅ All phases completed and integrated

```bash
# Verify readiness for integration testing
verify_integration_testing_ready() {
    echo "🔍 Verifying integration testing readiness..."
    
    # Check branch exists
    cd "$INTEGRATION_DIR"
    CURRENT_BRANCH=$(git branch --show-current)
    if [[ ! "$CURRENT_BRANCH" =~ integration-testing ]]; then
        echo "❌ Not on integration-testing branch!"
        exit 1
    fi
    
    # Check based on main
    BASE_BRANCH=$(git merge-base HEAD main)
    MAIN_HEAD=$(git rev-parse main)
    if [ "$BASE_BRANCH" != "$MAIN_HEAD" ]; then
        echo "❌ Not based on main HEAD!"
        exit 1
    fi
    
    # Check documentation exists
    if [ ! -f "INTEGRATION-INFO.md" ]; then
        echo "❌ Missing INTEGRATION-INFO.md!"
        exit 1
    fi
    
    echo "✅ All prerequisites verified"
    echo "Ready to transition to INTEGRATION_TESTING state"
}
```

## Summary

The CREATE_INTEGRATION_TESTING state is responsible for:
1. Creating a clean integration-testing branch from main HEAD
2. Setting up the integration testing workspace
3. Documenting the integration testing setup
4. Preparing for merging ALL effort branches
5. Ensuring R272 compliance (integration from main)
6. Maintaining R280 (never push to main)

This state is the gateway to final validation - proving the entire Software Factory output works as a cohesive whole before generating the MASTER-PR-PLAN for human execution.
## R322 VIOLATION DETECTION

If you find yourself:
- Starting work for a new state without /continue-orchestrating
- Transitioning without stopping after state file commit
- Continuing after completing state work

**STOP IMMEDIATELY - You are violating R322!**
