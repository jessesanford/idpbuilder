# Orchestrator - SETUP_EFFORT_INFRASTRUCTURE State Rules

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


## 🔴🔴🔴 CRITICAL: BASE BRANCH DETERMINATION (R308) 🔴🔴🔴

**VIOLATION = -100% AUTOMATIC FAILURE**

### YOU MUST DETERMINE THE CORRECT BASE BRANCH PER R308:
- Phase 1, Wave 1: from **main**
- Phase 1, Wave 2+: from **phase1-wave[N-1]-integration**  
- Phase 2+, Wave 1: from **phase[N-1]-integration** (NEVER main!)
- Phase 2+, Wave 2+: from **phase[N]-wave[N-1]-integration**

### NEVER BASE PHASE 2 ON MAIN!

Example for Phase 2, Wave 1:
```bash
# WRONG - AUTOMATIC FAILURE:
BASE_BRANCH="main"  # ❌ NEVER!

# CORRECT:
BASE_BRANCH="phase1-integration"  # ✅ From previous phase!
```

**Acknowledge: "I understand Phase 2 efforts MUST be based on phase1-integration per R308"**

## 🔴🔴🔴 STOP! STATE RULE READING IS ABSOLUTELY FIRST! 🔴🔴🔴

**YOU HAVE ENTERED SETUP_EFFORT_INFRASTRUCTURE STATE - YOU MUST READ AND ACKNOWLEDGE ALL STATE RULES BEFORE DOING ANY STATE WORK!**

## 🔴🔴🔴 R290 VERIFICATION REQUIREMENT 🔴🔴🔴

**R290 ENFORCEMENT: CREATE VERIFICATION MARKER AFTER READING**

After reading and acknowledging all state rules, you MUST create a verification marker:

```bash
# MANDATORY: Create verification marker after reading rules
touch .state_rules_read_orchestrator_SETUP_EFFORT_INFRASTRUCTURE
echo "$(date +%s) - Rules read and acknowledged for SETUP_EFFORT_INFRASTRUCTURE" > .state_rules_read_orchestrator_SETUP_EFFORT_INFRASTRUCTURE
```

**FAILURE TO CREATE MARKER = AUTOMATIC -100% PENALTY**

The system will check for this marker. No marker = Immediate failure.

### ❌ DO NOT DO ANY SETUP_EFFORT_INFRASTRUCTURE WORK UNTIL RULES ARE READ:
- ❌ Start create effort directories
- ❌ Start set up branches
- ❌ Start initialize effort tracking
- ❌ Start configure worktrees
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
   ❌ WRONG: "I acknowledge all SETUP_EFFORT_INFRASTRUCTURE rules"
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

### ✅ CORRECT PATTERN FOR SETUP_EFFORT_INFRASTRUCTURE:
```
1. READ: $CLAUDE_PROJECT_DIR/rule-library/[first-rule-file].md
2. "I acknowledge [Rule#] - [Rule Name]: [Brief description]"
3. READ: $CLAUDE_PROJECT_DIR/rule-library/[second-rule-file].md  
4. "I acknowledge [Rule#] - [Rule Name]: [Brief description]"
[Continue for EVERY rule in PRIMARY DIRECTIVES...]
5. "Ready to execute SETUP_EFFORT_INFRASTRUCTURE work"
```

### 🚨 NO WORK UNTIL ACKNOWLEDGMENT COMPLETE 🚨
**You may NOT begin ANY SETUP_EFFORT_INFRASTRUCTURE work until:**
1. ✅ ALL rules in PRIMARY DIRECTIVES have been READ
2. ✅ ALL rules have been individually ACKNOWLEDGED
3. ✅ You have stated readiness to execute SETUP_EFFORT_INFRASTRUCTURE work
1. **READ** every rule file listed in PRIMARY DIRECTIVES below
2. **ACKNOWLEDGE** each rule individually with number and description
3. **ONLY THEN** proceed with SETUP_EFFORT_INFRASTRUCTURE work

### 🚨 FAILURE TO READ STATE RULES FIRST = IMMEDIATE EXIT 🚨
**If you do ANY SETUP_EFFORT_INFRASTRUCTURE work before reading and acknowledging rules:**
- **STOP ALL WORK IMMEDIATELY**
- **EXIT WITH FAILURE STATUS**
- **YOU HAVE VIOLATED STATE COMPLIANCE**

**THE SYSTEM IS MONITORING YOUR READ TOOL CALLS!**

## 🔴🔴🔴 CRITICAL: NEVER CREATE BRANCHES IN SF REPO! 🔴🔴🔴

**BEFORE READING ANY RULES, UNDERSTAND THIS FUNDAMENTAL TRUTH:**

### ❌❌❌ EFFORT BRANCHES GO IN TARGET REPOSITORY CLONES, NOT HERE! ❌❌❌

**THIS IS THE SOFTWARE FACTORY REPO (PLANNING/ORCHESTRATION)**
- Path: `/home/vscode/software-factory-template/`
- Purpose: Rules, agents, state management ONLY
- **NEVER CREATE EFFORT BRANCHES HERE!**

**TARGET REPOSITORY IS COMPLETELY DIFFERENT**
- Defined in: `target-repo-config.yaml`
- Gets cloned to: `efforts/phaseX/waveY/effort-name/`
- **THIS IS WHERE ALL EFFORT BRANCHES GO!**

**VIOLATION = -100% AUTOMATIC FAILURE**

## 📋 PRIMARY DIRECTIVES FOR SETUP_EFFORT_INFRASTRUCTURE

**YOU MUST READ EACH RULE LISTED HERE. YOUR READ TOOL CALLS ARE BEING MONITORED.**

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

### State-Specific Rules (NOT in orchestrator.md):

1. **🚨🚨🚨 R006** - ORCHESTRATOR NEVER WRITES CODE OR PERFORMS FILE OPERATIONS (BLOCKING)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R006-orchestrator-never-writes-code.md`
   - Criticality: BLOCKING - Automatic termination, 0% grade
   - Summary: NEVER write, copy, move, or manipulate ANY code files - delegate ALL to agents
   - **CRITICAL**: Creating directories is OK, copying/moving files is FORBIDDEN!

2. **🚨🚨🚨 R315** - Infrastructure vs Implementation Boundary (BLOCKING)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R315-infrastructure-vs-implementation-boundary.md`
   - Criticality: BLOCKING - -50% for violations
   - Summary: Create ONLY empty infrastructure, NEVER copy/move code files
   - **CRITICAL**: Empty directories = infrastructure, File operations = implementation!

3. **R308** - INCREMENTAL BRANCHING STRATEGY (SUPREME LAW!)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R308-incremental-branching-strategy.md`
   - Criticality: 🔴🔴🔴 SUPREME LAW - Phase 2 NEVER from main!
   - **READ THIS FIRST - DETERMINES ALL BASE BRANCHES!**
   - **PHASE 2 FROM PHASE1-INTEGRATION, NOT MAIN!**

4. **R309** - NEVER Create Efforts in SF Repo (PARAMOUNT LAW!)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R309-never-create-efforts-in-sf-repo.md`
   - Criticality: PARAMOUNT - Automatic -100% failure for violation
   - **AVOID CATASTROPHIC FAILURE!**

5. **R312** - Git Config Immutability Protocol (SUPREME LAW!)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R312-git-config-immutability-protocol.md`
   - Criticality: 🔴🔴🔴 SUPREME LAW - Lock .git/config after setup
   - **CRITICAL FOR EFFORT ISOLATION!**
   - **PREVENTS BRANCH/REMOTE CONTAMINATION!**

6. **R191** - Target Repository Configuration
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R191-target-repo-config.md`
   - Criticality: BLOCKING - Must load config before proceeding
   
7. **R176** - Workspace Isolation  
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R176-workspace-isolation.md`
   - Criticality: BLOCKING - Ensure complete workspace isolation

8. **R271** - Single-Branch Full Checkout Protocol
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R271-single-branch-full-checkout.md`
   - Criticality: SUPREME LAW - Full single-branch checkout required

**Note**: R234, R208, R221, R287, R288 are already in your orchestrator.md Supreme Laws section.

## 📋 RULE SUMMARY FOR SETUP_EFFORT_INFRASTRUCTURE STATE

### Rules Enforced in This State:
- R308: INCREMENTAL BRANCHING [🔴 SUPREME LAW - PHASE 2 FROM PHASE1-INTEGRATION!]
- R309: NEVER Create Efforts in SF Repo [PARAMOUNT LAW - -100% FOR VIOLATION!]
- R234: Mandatory State Traversal [SUPREME LAW - NO SKIPPING!]
- R208: CD Before Spawn [SUPREME LAW - Always CD first]
- R221: Bash Directory Reset [SUPREME LAW - CD in every command]
- R191: Target Repository Configuration [BLOCKING - Must load config]
- R176: Effort Infrastructure Setup [BLOCKING - Create all directories]
- R271: Full Checkouts Only [SUPREME LAW - No sparse checkouts]
- R287: TODO Save Triggers [BLOCKING - Save within 30s]
- R288: State File Update and Commit [SUPREME LAW - Update on transition]

### Critical Requirements:
1. **DETERMINE CORRECT BASE BRANCH PER R308** - Penalty: -100%
   - Phase 2, Wave 1 MUST use phase1-integration (NOT main!)
   - Phase 2, Wave 2+ MUST use phase2-wave[N-1]-integration
2. Load target-repo-config.yaml - Penalty: -100%
3. Create ALL effort directories - Penalty: -50%
4. Clone FULL repos (no sparse) - Penalty: -100%
5. Clone from CORRECT incremental base - Penalty: -100%
6. Push all branches to remote - Penalty: -30%
7. Must transition to ANALYZE_CODE_REVIEWER_PARALLELIZATION - Penalty: -100%

### Success Criteria:
- ✅ Target config loaded and validated
- ✅ All effort directories created under /efforts/
- ✅ All repos are FULL clones (R271)
- ✅ All branches pushed with tracking
- ✅ work-log.md files initialized

### Failure Triggers:
- ❌ Skip to SPAWN_AGENTS = -100% R234 VIOLATION
- ❌ Clone wrong repository = AUTOMATIC FAILURE
- ❌ Sparse checkout detected = R271 VIOLATION
- ❌ Missing target-repo-config.yaml = Cannot proceed

## 🚨 SETUP_EFFORT_INFRASTRUCTURE IS A VERB - CREATE INFRASTRUCTURE NOW! 🚨

### IMMEDIATE ACTIONS UPON ENTERING SETUP_EFFORT_INFRASTRUCTURE

**THE MOMENT YOU ENTER THIS STATE, YOU MUST:**
1. Start creating effort directories NOW using prepare_effort_for_agent()
2. Initialize Git branches for each effort immediately
3. Set up remote tracking for all branches without delay
4. Check TodoWrite for pending infrastructure tasks
5. Create standard subdirectories (src/, tests/, docs/) immediately

**FORBIDDEN - AUTOMATIC FAILURE:**
- ❌ "STATE TRANSITION COMPLETE: Now in SETUP_EFFORT_INFRASTRUCTURE" [stops]
- ❌ "Successfully entered SETUP_EFFORT_INFRASTRUCTURE state" [waits]
- ❌ "Ready to set up infrastructure" [pauses]
- ❌ "I'm in infrastructure setup state" [does nothing]

**REQUIRED - IMMEDIATE ACTION:**
- ✅ "Entering SETUP_EFFORT_INFRASTRUCTURE, creating E3.1.1 directory now..."
- ✅ "Setting up infrastructure, cloning repo for E3.1.1..."
- ✅ "SETUP_EFFORT_INFRASTRUCTURE: Creating branch for E3.1.1..."

## State Context
You are setting up infrastructure for all efforts in the wave BEFORE spawning code reviewers to create effort plans - DO IT NOW!

## 🔴🔴🔴 SUPREME LAW R234 - MANDATORY STATE TRAVERSAL 🔴🔴🔴

**THIS IS THE HIGHEST LAW - SUPERSEDES ALL OTHER RULES!**

### MANDATORY NEXT STATE: ANALYZE_CODE_REVIEWER_PARALLELIZATION

**YOU MUST FOLLOW THIS EXACT SEQUENCE:**
```
SETUP_EFFORT_INFRASTRUCTURE (YOU ARE HERE)
    ↓ (CANNOT SKIP - MANDATORY)
ANALYZE_CODE_REVIEWER_PARALLELIZATION
    ↓ (CANNOT SKIP - MANDATORY)
SPAWN_CODE_REVIEWERS_EFFORT_PLANNING
    ↓ (CANNOT SKIP - MANDATORY)
WAITING_FOR_EFFORT_PLANS
    ↓ (CANNOT SKIP - MANDATORY)
ANALYZE_IMPLEMENTATION_PARALLELIZATION
    ↓ (CANNOT SKIP - MANDATORY)
SPAWN_AGENTS
```

**❌❌❌ FORBIDDEN TRANSITIONS (AUTOMATIC -100% FAILURE):**
- ❌ SETUP_EFFORT_INFRASTRUCTURE → SPAWN_AGENTS (skipping analysis)
- ❌ SETUP_EFFORT_INFRASTRUCTURE → SPAWN_CODE_REVIEWERS_EFFORT_PLANNING (skipping analysis)
- ❌ ANY attempt to "optimize" by skipping states

**ACKNOWLEDGMENT REQUIRED:**
"I acknowledge R234: I MUST transition to ANALYZE_CODE_REVIEWER_PARALLELIZATION next, not skip ahead."

## 🔴🔴🔴 CRITICAL: LOAD TARGET CONFIG FIRST (R191) 🔴🔴🔴

### BEFORE ANY INFRASTRUCTURE SETUP, YOU MUST:
```bash
# MANDATORY FIRST ACTION IN THIS STATE
echo "🔴 R191: Loading target repository configuration..."

# Check if already loaded
if [ -z "$TARGET_REPO_URL" ]; then
    echo "⚠️ Target config not loaded, loading now..."
    
    if [ ! -f "$SF_ROOT/target-repo-config.yaml" ]; then
        echo "🔴🔴🔴 CRITICAL: target-repo-config.yaml NOT FOUND!"
        echo "Cannot set up infrastructure without knowing WHAT to clone!"
        exit 191
    fi
    
    # Load and validate config
    TARGET_REPO_URL=$(yq '.target_repository.url' "$SF_ROOT/target-repo-config.yaml")
    BASE_BRANCH=$(yq '.target_repository.base_branch' "$SF_ROOT/target-repo-config.yaml")
    PROJECT_PREFIX=$(yq '.branch_naming.project_prefix' "$SF_ROOT/target-repo-config.yaml")
    
    if [ -z "$TARGET_REPO_URL" ] || [ "$TARGET_REPO_URL" = "null" ]; then
        echo "🔴 ERROR: No target repository URL in config!"
        exit 191
    fi
    
    export TARGET_REPO_URL
    export BASE_BRANCH
    export PROJECT_PREFIX
fi

echo "✅ Target repository: $TARGET_REPO_URL"
echo "✅ This is what will be CLONED into efforts/"
echo "⚠️ NEVER create code in the Software Factory repo itself!"
```

### UNDERSTAND THE TWO REPOSITORIES:
1. **Software Factory Repo** = Where you are now (rules, agents, state)
2. **Target Repo** = What you clone from config (actual project code)

**VIOLATIONS THAT CAUSE FAILURE:**
- Cloning the SF repo into efforts/ = WRONG REPO
- Creating code in SF repo = WRONG LOCATION
- Not having target-repo-config.yaml = CANNOT PROCEED

## 🔴🔴🔴 MANDATORY: Infrastructure BEFORE Planning 🔴🔴🔴

**THIS IS THE CORRECT SEQUENCE:**
1. Architect creates Wave Architecture Plan
2. Code Reviewer creates Wave Implementation Plan
3. **Orchestrator sets up ALL effort infrastructure** ← YOU ARE HERE
4. Code Reviewers create individual Effort Implementation Plans
5. SW Engineers implement


## Infrastructure Setup Protocol (R271 SUPREME LAW - Full Checkouts)

```bash
# 🔴🔴🔴 R308 BASE BRANCH DETERMINATION FUNCTION 🔴🔴🔴
determine_incremental_base_branch() {
    local PHASE=$1
    local WAVE=$2
    
    echo "🔴 R308: Determining incremental base branch for Phase $PHASE, Wave $WAVE"
    
    # Phase 1, Wave 1: Start from main
    if [[ $PHASE -eq 1 && $WAVE -eq 1 ]]; then
        echo "📌 R308: Phase 1, Wave 1 → base: main"
        echo "main"
        return
    fi
    
    # First wave of new phase: From previous phase integration
    if [[ $WAVE -eq 1 ]]; then
        PREV_PHASE=$((PHASE - 1))
        BASE="phase${PREV_PHASE}-integration"
        
        echo "🔴 R308: Phase $PHASE, Wave 1 → base: $BASE (NOT main!)"
        
        # Verify it exists
        if git ls-remote --heads origin "$BASE" > /dev/null 2>&1; then
            echo "✅ R308: Previous phase integration found: $BASE"
            echo "$BASE"
        else
            echo "❌ R308 FATAL: Previous phase integration not found: $BASE"
            echo "Cannot proceed without phase integration!"
            exit 308
        fi
        return
    fi
    
    # Subsequent waves: From previous wave integration
    PREV_WAVE=$((WAVE - 1))
    BASE="phase${PHASE}-wave${PREV_WAVE}-integration"
    
    echo "📌 R308: Phase $PHASE, Wave $WAVE → base: $BASE"
    
    # Verify it exists
    if git ls-remote --heads origin "$BASE" > /dev/null 2>&1; then
        echo "✅ R308: Previous wave integration found: $BASE"
        echo "$BASE"
    else
        echo "❌ R308 FATAL: Previous wave integration not found: $BASE"
        echo "Cannot proceed without wave integration!"
        exit 308
    fi
}

# MANDATORY: Use prepare_effort_for_agent() for EACH effort
prepare_effort_for_agent() {
    local PHASE=$1 WAVE=$2 EFFORT=$3
    
    echo "═══════════════════════════════════════════════════════"
    echo "🔧 SETTING UP INFRASTRUCTURE FOR: $EFFORT"
    echo "═══════════════════════════════════════════════════════"
    
    # 🔴🔴🔴 R309 VALIDATION: NEVER CREATE BRANCHES IN SF REPO! 🔴🔴🔴
    echo "🔴 R309 CHECK: Validating we're not in SF repo..."
    if [ -f ".claude/CLAUDE.md" ] || [ -f "rule-library/RULE-REGISTRY.md" ]; then
        echo "✅ Confirmed: In SF planning repo (correct location to START)"
    else
        echo "⚠️ WARNING: Not in SF root - returning to correct location"
        cd $CLAUDE_PROJECT_DIR
    fi
    
    # Verify we're NOT about to pollute the SF repo
    local CURRENT_REMOTE=$(git remote get-url origin 2>/dev/null)
    if [[ "$CURRENT_REMOTE" == *"software-factory"* ]]; then
        echo "✅ Confirmed: This is SF repo - will clone TARGET for efforts"
    fi
    
    # 🔴🔴🔴 R308: INCREMENTAL BASE BRANCH DETERMINATION 🔴🔴🔴
    echo "🔴🔴🔴 R308 ENFORCEMENT: Determining INCREMENTAL base branch"
    BASE_BRANCH=$(determine_incremental_base_branch $PHASE $WAVE)
    
    # CRITICAL VALIDATION FOR PHASE 2+
    if [[ $PHASE -gt 1 && "$BASE_BRANCH" == "main" ]]; then
        echo "🔴🔴🔴 R308 VIOLATION: Phase $PHASE CANNOT use main!"
        echo "FATAL ERROR: Must use phase$((PHASE-1))-integration"
        exit 308
    fi
    
    echo "✅ R308 VALIDATED: Using incremental base: $BASE_BRANCH"
    
    # 2. Create effort directory under efforts/ (relative to SF instance root)
    EFFORT_DIR="${CLAUDE_PROJECT_DIR}/efforts/phase${PHASE}/wave${WAVE}/${EFFORT}"
    mkdir -p "$(dirname "$EFFORT_DIR")"
    
    # 3. SINGLE-BRANCH FULL CLONE (R271 Supreme Law)
    echo "📦 Creating FULL clone from branch: $BASE_BRANCH"
    
    # CRITICAL: Get target repo URL from config (NOT the SF repo!)
    TARGET_REPO_URL=$(yq '.target_repository.url' "$SF_ROOT/target-repo-config.yaml")
    
    if [ -z "$TARGET_REPO_URL" ] || [ "$TARGET_REPO_URL" = "null" ]; then
        echo "🔴🔴🔴 R191 VIOLATION: No target repository URL!"
        echo "Cannot clone without target-repo-config.yaml!"
        exit 191
    fi
    
    echo "🎯 Cloning TARGET repository: $TARGET_REPO_URL"
    echo "⚠️ This is NOT the Software Factory repo!"
    
    # 🔴 R309 CRITICAL CHECK: Verify target is NOT SF repo!
    if [[ "$TARGET_REPO_URL" == *"software-factory"* ]]; then
        echo "🔴🔴🔴 R309 VIOLATION: Target URL is Software Factory repo!"
        echo "FATAL ERROR: You're trying to clone SF into itself!"
        echo "Fix target-repo-config.yaml to point to actual project!"
        exit 309
    fi
    
    git clone \
        --single-branch \
        --branch "$BASE_BRANCH" \
        "$TARGET_REPO_URL" \
        "$EFFORT_DIR"
    
    if [ $? -ne 0 ]; then
        echo "❌ Clone failed! Check if base branch '$BASE_BRANCH' exists"
        exit 1
    fi
    
    cd "$EFFORT_DIR"
    
    # 🔴 R309 POST-CLONE VALIDATION: Ensure we cloned the right thing!
    if [ -f ".claude/CLAUDE.md" ] || [ -f "rule-library/RULE-REGISTRY.md" ]; then
        echo "🔴🔴🔴 R309 VIOLATION: Cloned SF repo instead of target!"
        echo "FATAL ERROR: This is the wrong repository!"
        echo "You must clone the TARGET project, not SF template!"
        cd "$SF_ROOT"
        rm -rf "$EFFORT_DIR"
        exit 309
    fi
    echo "✅ R309 VALIDATED: This is TARGET repo (not SF)"
    
    # 4. Create and push effort branch (WITH PROJECT PREFIX!)
    PROJECT_PREFIX=$(yq '.branch_naming.project_prefix' "$SF_ROOT/target-repo-config.yaml")
    if [ -n "$PROJECT_PREFIX" ] && [ "$PROJECT_PREFIX" != "null" ]; then
        BRANCH="${PROJECT_PREFIX}/phase${PHASE}/wave${WAVE}/${EFFORT}"
    else
        BRANCH="phase${PHASE}/wave${WAVE}/${EFFORT}"
    fi
    
    echo "🌿 Creating effort branch: $BRANCH"
    git checkout -b "$BRANCH"
    git push -u origin "$BRANCH"
    
    # 5. Create .software-factory directory structure for plans
    echo "📁 Creating .software-factory directory structure for plans..."
    FACTORY_DIR=".software-factory/phase${PHASE}/wave${WAVE}/${EFFORT}"
    mkdir -p "$FACTORY_DIR"
    echo "✅ Created plan storage directory: $FACTORY_DIR"
    echo "   This is where Code Reviewer will save IMPLEMENTATION-PLAN-*.md"
    
    # 6. Create work-log.md with R308 compliance documentation
    echo "# Work Log for $EFFORT" > work-log.md
    echo "" >> work-log.md
    echo "## Infrastructure Details" >> work-log.md
    echo "- **Branch**: $BRANCH" >> work-log.md
    echo "- **Base Branch**: $BASE_BRANCH" >> work-log.md
    echo "- **Clone Type**: FULL (R271 compliance)" >> work-log.md
    echo "- **Created**: $(date)" >> work-log.md
    echo "" >> work-log.md
    echo "## R308 Incremental Branching Compliance" >> work-log.md
    echo "- **Phase**: $PHASE" >> work-log.md
    echo "- **Wave**: $WAVE" >> work-log.md
    echo "- **R308 Rule Applied**: $(determine_incremental_base_branch $PHASE $WAVE)" >> work-log.md
    if [[ $PHASE -eq 2 && $WAVE -eq 1 ]]; then
        echo "- **CRITICAL**: Phase 2 Wave 1 correctly based on phase1-integration (NOT main)" >> work-log.md
    elif [[ $PHASE -gt 1 ]]; then
        echo "- **Incremental**: Building on previous integration as required" >> work-log.md
    fi
    
    # 7. Verify FULL workspace (R271 compliance check)
    echo "🔍 Verifying full checkout..."
    if [ -f ".git/info/sparse-checkout" ]; then
        echo "🔴🔴🔴 SUPREME LAW VIOLATION: Sparse checkout detected!"
        exit 1
    fi
    
    # Show that we have the full codebase
    echo "✅ Full codebase available:"
    ls -la | head -10
    echo "   ... (showing first 10 entries)"
    
    # 🔴🔴🔴 R312: LOCK GIT CONFIG FOR EFFORT ISOLATION 🔴🔴🔴
    echo "🔒 R312: Applying DOUBLE PROTECTION to git config..."
    
    # Verify .git/config exists
    if [ ! -f .git/config ]; then
        echo "❌ FATAL: No .git/config found in $EFFORT_DIR"
        exit 312
    fi
    
    # Store current permissions and ownership for audit
    BEFORE_PERMS=$(stat -c %a .git/config 2>/dev/null || stat -f %A .git/config)
    BEFORE_OWNER=$(stat -c %U:%G .git/config 2>/dev/null || stat -f %Su:%Sg .git/config)
    
    # DOUBLE PROTECTION: Change ownership AND permissions
    if command -v sudo >/dev/null 2>&1; then
        # Full protection with root ownership
        echo "🔐 Applying FULL protection (root ownership + readonly)..."
        sudo chown root:root .git/config
        sudo chmod 444 .git/config
        PROTECTION_LEVEL="FULL"
    else
        # Fallback to permission-only protection
        echo "⚠️ sudo not available - applying permission-only protection"
        chmod 444 .git/config
        PROTECTION_LEVEL="PARTIAL"
    fi
    
    # Verify it's now readonly
    if [ -w .git/config ]; then
        echo "❌ R312 VIOLATION: Failed to make .git/config readonly!"
        echo "Config is still writable - effort isolation compromised!"
        exit 312
    fi
    
    # Verify ownership if sudo was available
    CURRENT_OWNER=$(stat -c %U:%G .git/config 2>/dev/null || stat -f %Su:%Sg .git/config)
    if [ "$PROTECTION_LEVEL" = "FULL" ] && [ "$CURRENT_OWNER" != "root:root" ]; then
        echo "⚠️ WARNING: Ownership not changed to root:root (got $CURRENT_OWNER)"
        echo "Protection may be weaker than intended"
    fi
    
    # Create protection marker with audit info
    cat > .git/R312_CONFIG_LOCKED << EOF
Timestamp: $(date '+%Y-%m-%d %H:%M:%S')
Locked by: orchestrator
State: SETUP_EFFORT_INFRASTRUCTURE
Effort: $EFFORT
Phase: $PHASE
Wave: $WAVE
Protection level: $PROTECTION_LEVEL
Previous ownership: $BEFORE_OWNER
Current ownership: $CURRENT_OWNER
Previous permissions: $BEFORE_PERMS
Current permissions: 444 (readonly)
Purpose: Prevent branch/remote modifications per R312
EOF
    
    echo "✅ R312 ENFORCED: Git config locked"
    echo "   Protection: $PROTECTION_LEVEL"
    echo "   Ownership: $BEFORE_OWNER → $CURRENT_OWNER"
    echo "   Permissions: $BEFORE_PERMS → 444"
    echo "📝 Protected operations now BLOCKED:"
    echo "   ❌ git checkout [other-branch]"
    echo "   ❌ git pull origin main"
    echo "   ❌ git remote add/remove"
    echo "   ❌ git branch --set-upstream-to"
    echo "📝 Allowed operations still work:"
    echo "   ✅ git add, commit, push"
    echo "   ✅ git status, diff, log"
    
    echo "✅ Infrastructure ready for $EFFORT with FULL code from $BASE_BRANCH"
    echo "🔒 Config locked per R312 - effort isolation guaranteed"
    cd "$SF_ROOT"  # Return to root
}
```

## 🔴🔴🔴 CRITICAL EXAMPLES: R308 BASE BRANCH DETERMINATION 🔴🔴🔴

### Example 1: Phase 2, Wave 1 Infrastructure (MOST COMMON MISTAKE!)
```bash
# ❌❌❌ WRONG - AUTOMATIC -100% FAILURE:
PHASE=2 WAVE=1
BASE_BRANCH="main"  # NEVER DO THIS FOR PHASE 2!

# ✅✅✅ CORRECT - MUST USE PHASE 1 INTEGRATION:
PHASE=2 WAVE=1
BASE_BRANCH=$(determine_incremental_base_branch 2 1)  # Returns: phase1-integration
echo "Phase 2 Wave 1 MUST use: $BASE_BRANCH (NOT main!)"
```

### Example 2: Phase 2, Wave 2 Infrastructure
```bash
PHASE=2 WAVE=2
BASE_BRANCH=$(determine_incremental_base_branch 2 2)  # Returns: phase2-wave1-integration
echo "Phase 2 Wave 2 uses: $BASE_BRANCH"
```

### Example 3: Phase 3, Wave 1 Infrastructure
```bash
PHASE=3 WAVE=1
BASE_BRANCH=$(determine_incremental_base_branch 3 1)  # Returns: phase2-integration
echo "Phase 3 Wave 1 uses: $BASE_BRANCH (from previous PHASE)"
```

## Example: Setting Up Wave Infrastructure

```bash
# Read current phase/wave from state
PHASE=$(yq '.current_phase' orchestrator-state.yaml)
WAVE=$(yq '.current_wave' orchestrator-state.yaml)

# 🔴 CRITICAL: Check if this is Phase 2!
if [[ $PHASE -eq 2 && $WAVE -eq 1 ]]; then
    echo "🔴🔴🔴 ATTENTION: Phase 2, Wave 1 detected!"
    echo "MUST use phase1-integration as base (R308)"
fi

# Read wave plan to get effort list
WAVE_PLAN="phase-plans/PHASE-${PHASE}-WAVE-${WAVE}-IMPLEMENTATION-PLAN.md"
EFFORTS=$(grep "^## Effort" "$WAVE_PLAN" | sed 's/## Effort [0-9]*: //')

# Setup infrastructure for ALL efforts
for effort in $EFFORTS; do
    prepare_effort_for_agent $PHASE $WAVE "$effort"
done

echo "✅ All effort infrastructure ready with R308 compliant bases"
echo "📋 Ready to spawn Code Reviewers for effort planning"
```

## State Transition

After ALL infrastructure is ready:
1. Update orchestrator-state.yaml with effort directories
2. Verify all branches pushed to remote
3. **MANDATORY: Transition to ANALYZE_CODE_REVIEWER_PARALLELIZATION (R234)**
   - DO NOT skip to SPAWN_CODE_REVIEWERS_EFFORT_PLANNING
   - DO NOT skip to SPAWN_AGENTS
   - MUST follow the mandatory sequence!

### TODO PERSISTENCE CHECKPOINT (R287-R287)
```bash
# Before state transition - MANDATORY SAVE
echo "💾 R287: Saving TODOs before state transition..."
save_todos "SETUP_EFFORT_INFRASTRUCTURE complete"

# R287: Commit within 60 seconds
cd $CLAUDE_PROJECT_DIR
git add todos/*.todo
git commit -m "todo: orchestrator - infrastructure setup complete"
git push
echo "✅ TODOs persisted before transition"
```

## 🚨🚨🚨 R209 - EFFORT DIRECTORY ISOLATION PROTOCOL (MISSION CRITICAL!)
**Source:** rule-library/R209-effort-directory-isolation-protocol.md  
**Criticality:** MISSION CRITICAL - SW Engineers MUST stay in effort directories  

### MUST INJECT METADATA AND PUSH TO REMOTE!

After Code Reviewer creates IMPLEMENTATION-PLAN.md, YOU MUST ADD DIRECTORY/BRANCH METADATA AND PUSH:

```bash
# 🚨 MANDATORY: Call this AFTER Code Reviewer creates plan, BEFORE spawning SW Engineer! 🚨
inject_r209_metadata() {
    local EFFORT_NAME="$1"
    local PHASE="$2"
    local WAVE="$3"
    local IMPL_PLAN="efforts/phase${PHASE}/wave${WAVE}/${EFFORT_NAME}/IMPLEMENTATION-PLAN.md"
    
    echo "🔧 [R209] Injecting metadata for effort: $EFFORT_NAME"
    
    # Check if plan exists
    if [ ! -f "$IMPL_PLAN" ]; then
        echo "⚠️ ERROR: Implementation plan not found at: $IMPL_PLAN"
        return 1
    fi
    
    # Source branch naming helpers
    source utilities/branch-naming-helpers.sh
    
    # Get properly formatted branch name with project prefix
    EFFORT_BRANCH=$(get_effort_branch_name "$PHASE" "$WAVE" "$EFFORT_NAME")
    
    # Add metadata header
    cat > /tmp/r209_metadata.md << EOF
<!-- ⚠️ EFFORT INFRASTRUCTURE METADATA (R209) ⚠️ -->
**EFFORT_NAME**: ${EFFORT_NAME}
**PHASE**: ${PHASE}
**WAVE**: ${WAVE}
**WORKING_DIRECTORY**: $(pwd)/efforts/phase${PHASE}/wave${WAVE}/${EFFORT_NAME}
**BRANCH**: ${EFFORT_BRANCH}
**REMOTE**: origin/${EFFORT_BRANCH}
**ISOLATION_BOUNDARY**: efforts/phase${PHASE}/wave${WAVE}/${EFFORT_NAME}

⚠️ **SW ENGINEER: YOU MUST STAY IN THIS DIRECTORY!** ⚠️
ALL work happens in: efforts/phase${PHASE}/wave${WAVE}/${EFFORT_NAME}
ALL code goes in: efforts/phase${PHASE}/wave${WAVE}/${EFFORT_NAME}/pkg/
NEVER leave this directory during implementation!
<!-- END METADATA -->

EOF
    
    # Prepend to plan
    cat /tmp/r209_metadata.md "$IMPL_PLAN" > /tmp/updated_plan.md
    mv /tmp/updated_plan.md "$IMPL_PLAN"
    echo "✅ R209: Metadata injected"
    
    # 🔴 CRITICAL: Push to remote!
    ORCH_DIR=$(pwd)
    cd "efforts/phase${PHASE}/wave${WAVE}/${EFFORT_NAME}"
    
    git add IMPLEMENTATION-PLAN.md
    git commit -m "feat: inject R209 metadata into implementation plan"
    
    if git push; then
        echo "✅ Plan with R209 metadata pushed to remote"
    else
        git push -u origin "$(git branch --show-current)"
    fi
    
    cd "$ORCH_DIR"
}
```

### MANDATORY WORKFLOW:
1. Code Reviewer creates IMPLEMENTATION-PLAN.md
2. **WAIT FOR COMPLETION**
3. **INJECT R209 METADATA** (orchestrator does this)
4. **PUSH TO REMOTE**
5. Then spawn SW Engineer

**Note:** This metadata injection typically happens in WAITING_FOR_EFFORT_PLANS state, but the infrastructure must be ready in SETUP_EFFORT_INFRASTRUCTURE.

## Common Mistakes to Avoid

❌ **WRONG:** Spawning Code Reviewers before creating directories
❌ **WRONG:** Forgetting project prefix in branch names
❌ **WRONG:** Not pushing branches to remote
❌ **WRONG:** Creating infrastructure one-by-one as needed
❌ **WRONG:** Not injecting R209 metadata before spawning SW Engineers

✅ **CORRECT:** Create ALL infrastructure first, inject metadata, then spawn agents

## R322 VIOLATION DETECTION

If you find yourself:
- Starting work for a new state without /continue-orchestrating
- Transitioning without stopping after state file commit
- Continuing after completing state work

**STOP IMMEDIATELY - You are violating R322!**
