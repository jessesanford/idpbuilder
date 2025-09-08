# Orchestrator - CREATE_NEXT_SPLIT_INFRASTRUCTURE State Rules

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


## 🔴🔴🔴 CRITICAL: SPLITS MUST CHAIN SEQUENTIALLY! 🔴🔴🔴

**SUPREME LAW - VIOLATION = -100% IMMEDIATE FAILURE**

### SPLIT CHAINING IS MANDATORY (R308):
```
Split-001: Based on phase/wave integration branch ✓
Split-002: Based on Split-001 branch (NOT integration!) ✓
Split-003: Based on Split-002 branch (NOT integration!) ✓
```

### ❌ CATASTROPHIC VIOLATION (What's happening NOW):
```
Split-001: Based on phase1/wave1-integration ✓
Split-002: Based on phase1/wave1-integration ❌ WRONG!
Split-003: Based on phase1/wave1-integration ❌ WRONG!
Result: SPLITS LOSE EACH OTHER'S WORK!
```

### ✅ MANDATORY IMPLEMENTATION:
```bash
# For Split N+1, the base is ALWAYS Split N:
if [ $SPLIT_NUM -eq 1 ]; then
    BASE_BRANCH="phase1/wave1-integration"  # First split only
else
    BASE_BRANCH="phase1/wave1/effort-split-$(printf "%03d" $((SPLIT_NUM - 1)))"
fi
```

**THIS IS NON-NEGOTIABLE - SPLITS BUILD ON EACH OTHER!**

---

## 🔴🔴🔴 R290 ENFORCEMENT: READ THESE RULES FIRST! 🔴🔴🔴

**SUPREME LAW #3 (R290): STATE RULES MUST BE READ BEFORE STATE ACTIONS**

## 🔴🔴🔴 STOP! STATE RULE READING IS ABSOLUTELY FIRST! 🔴🔴🔴

**YOU HAVE ENTERED CREATE_NEXT_SPLIT_INFRASTRUCTURE STATE - YOU MUST READ AND ACKNOWLEDGE ALL STATE RULES BEFORE DOING ANY STATE WORK!**

## 🔴🔴🔴 R290 VERIFICATION REQUIREMENT 🔴🔴🔴

**R290 ENFORCEMENT: CREATE VERIFICATION MARKER AFTER READING**

After reading and acknowledging all state rules, you MUST create a verification marker:

```bash
# MANDATORY: Create verification marker after reading rules
touch .state_rules_read_orchestrator_CREATE_NEXT_SPLIT_INFRASTRUCTURE
echo "$(date +%s) - Rules read and acknowledged for CREATE_NEXT_SPLIT_INFRASTRUCTURE" > .state_rules_read_orchestrator_CREATE_NEXT_SPLIT_INFRASTRUCTURE
```

**FAILURE TO CREATE MARKER = AUTOMATIC -100% PENALTY**

## 🔴🔴🔴 CRITICAL: CREATE_NEXT_SPLIT_INFRASTRUCTURE IS A VERB! 🔴🔴🔴

**YOU MUST CREATE THE NEXT SPLIT INFRASTRUCTURE IMMEDIATELY!**
- ❌ NOT "Preparing to create split infrastructure"
- ❌ NOT "Ready to set up next split"
- ❌ NOT "Will now create infrastructure"
- ✅ YES "Creating split-002 infrastructure NOW"
- ✅ YES "Cloning repository for split-002 NOW"
- ✅ YES "Setting up branch for split-002 NOW"

## State Context

CREATE_NEXT_SPLIT_INFRASTRUCTURE = You ARE ACTIVELY creating the infrastructure for the NEXT split in sequence THIS INSTANT!

**🔴🔴🔴 CRITICAL: This state should ONLY be entered AFTER: 🔴🔴🔴**
1. Previous split implementation completed
2. Code Reviewer reviewed the previous split
3. Review passed AND more splits are needed

**NEVER enter this state directly after split implementation without review!**

## 📋 PRIMARY DIRECTIVES FOR CREATE_NEXT_SPLIT_INFRASTRUCTURE STATE

### 🚨🚨🚨 R006 - ORCHESTRATOR NEVER WRITES CODE OR PERFORMS FILE OPERATIONS (BLOCKING)
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R006-orchestrator-never-writes-code.md`
**Criticality**: BLOCKING - Automatic termination, 0% grade
**Summary**: NEVER write, copy, move, or manipulate ANY code files - delegate ALL to agents
**CRITICAL**: Copying files is NOT infrastructure - it's implementation work!

### 🚨🚨🚨 R315 - Infrastructure vs Implementation Boundary (BLOCKING)
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R315-infrastructure-vs-implementation-boundary.md`
**Criticality**: BLOCKING - -50% for violations
**Summary**: Create ONLY empty infrastructure, NEVER copy/move code files
**CRITICAL**: Directory creation = OK, File operations = FORBIDDEN

### 🔴🔴🔴 R251 - UNIVERSAL REPOSITORY SEPARATION LAW (PARAMOUNT)
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R251-REPOSITORY-SEPARATION-LAW.md`
**Criticality**: PARAMOUNT - Automatic -100% failure for violation
**Summary**: Software Factory = Planning ONLY, Target Repo = Code ONLY
**CRITICAL**: ALL infrastructure MUST be created in TARGET repo clones under /efforts/

### 🔴🔴🔴 R309 - NEVER Create Efforts in SF Repo (PARAMOUNT LAW)
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R309-never-create-efforts-in-sf-repo.md`
**Criticality**: PARAMOUNT - Automatic -100% failure for violation
**Summary**: NEVER create effort branches in Software Factory repo
**CRITICAL**: Validate you're in TARGET clone before ANY branch creation!

### 🔴🔴🔴 R312 - Git Config Immutability Protocol (SUPREME LAW)
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R312-git-config-immutability-protocol.md`
**Criticality**: SUPREME LAW - Lock .git/config after infrastructure setup
**Summary**: Make .git/config READONLY to prevent branch/remote changes
**CRITICAL**: Ensures split isolation - SW engineers cannot contaminate work!

### 🔴🔴🔴 R204 - Orchestrator Split Infrastructure (Just-In-Time)
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R204-orchestrator-split-infrastructure.md`
**Criticality**: BLOCKING - Must create split infrastructure correctly
**Summary**: Create ONLY the next split's infrastructure, based on previous split

### 🔴🔴🔴 R308 - Incremental Branching Strategy
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R308-incremental-branching-strategy.md`
**Criticality**: SUPREME LAW - Splits must branch sequentially
**Summary**: First split uses R308 base, subsequent splits based on previous

### 🚨🚨🚨 R302 - Comprehensive Split Tracking Protocol
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R302-comprehensive-split-tracking-protocol.md`
**Criticality**: BLOCKING - Track all split operations
**Summary**: Update split_tracking in orchestrator-state.yaml

### 🚨🚨🚨 R288 - State File Update and Commit Protocol
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R288-state-file-update-and-commit-protocol.md`
**Criticality**: SUPREME LAW - Update on every transition
**Summary**: Update orchestrator-state.yaml with split progress

### 🚨🚨🚨 R221 - Bash Reset Protocol (BLOCKING)
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R221-bash-reset-protocol.md`
**Criticality**: BLOCKING - Must reset bash state between operations
**Summary**: Reset variables and state when creating new split infrastructure

### 🚨🚨🚨 R216 - Bash Execution Syntax Protocol (BLOCKING)
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R216-bash-execution-syntax.md`
**Criticality**: BLOCKING - Incorrect syntax causes failures
**Summary**: Use parentheses for subshells, proper variable syntax

### 🚨🚨🚨 R235 - Pre-flight Verification Checklist (BLOCKING)
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R235-MANDATORY-PREFLIGHT-VERIFICATION-SUPREME-LAW.md`
**Criticality**: BLOCKING - Must verify environment before setup
**Summary**: Check directories, permissions, branches before split infrastructure

## 🚨🚨🚨 IMMEDIATE ACTIONS UPON ENTERING STATE 🚨🚨🚨

**THE INSTANT YOU ENTER THIS STATE, DO THIS:**

```bash
# ✅ CORRECT - IMMEDIATE ACTION
echo "🔧 CREATING NEXT SPLIT INFRASTRUCTURE NOW..."

# Step 1: Identify which split to create (DO NOW!)
PHASE=$(yq '.current_phase' orchestrator-state.yaml)
WAVE=$(yq '.current_wave' orchestrator-state.yaml)
EFFORT_NAME=$(yq '.split_tracking | keys | .[0]' orchestrator-state.yaml)  # Get effort with splits
CURRENT_SPLIT=$(yq ".split_tracking.\"$EFFORT_NAME\".current_split // 0" orchestrator-state.yaml)
NEXT_SPLIT=$((CURRENT_SPLIT + 1))

echo "📊 Creating infrastructure for split-$(printf "%03d" $NEXT_SPLIT)"

# Step 2: Load branch naming helpers (DO NOW!)
source "$CLAUDE_PROJECT_DIR/utilities/branch-naming-helpers.sh"

# Step 3: Determine base branch for this split (DO NOW!)
# 🔴🔴🔴 CRITICAL: SPLITS MUST CHAIN SEQUENTIALLY! 🔴🔴🔴
if [ $NEXT_SPLIT -eq 1 ]; then
    # First split ONLY: Use R308 incremental base
    echo "🔴 R308: First split uses incremental base from wave/phase"
    BASE_BRANCH=$(determine_effort_base_branch $PHASE $WAVE)
    echo "✅ Split-001 will be based on: $BASE_BRANCH"
else
    # 🔴 CRITICAL: All subsequent splits MUST be based on PREVIOUS split!
    # NOT on the integration branch! This ensures each split builds on the last!
    PREV_SPLIT=$(printf "%03d" $CURRENT_SPLIT)
    BASE_BRANCH=$(get_split_branch_name "$EFFORT_NAME" "$PREV_SPLIT")
    echo "🔴🔴🔴 SEQUENTIAL CHAINING: Split-$(printf "%03d" $NEXT_SPLIT) MUST be based on: $BASE_BRANCH"
    echo "❌ NOT based on integration branch!"
    echo "✅ YES based on previous split-$PREV_SPLIT!"
fi

# Step 4: Create the infrastructure (DO NOW!)
create_single_split_infrastructure "$EFFORT_NAME" "$NEXT_SPLIT" "$BASE_BRANCH"

# Step 5: Update state tracking (DO NOW!)
update_split_tracking "$EFFORT_NAME" "$NEXT_SPLIT" "INFRASTRUCTURE_READY"

# Step 6: Transition to SPAWN_AGENTS (DO NOW!)
echo "✅ Infrastructure ready for split-$(printf "%03d" $NEXT_SPLIT)"
transition_to_state "SPAWN_AGENTS"
```

## Core Implementation Function

```bash
create_single_split_infrastructure() {
    local EFFORT_NAME="$1"
    local SPLIT_NUM="$2"
    local BASE_BRANCH="$3"
    
    SPLIT_NAME=$(printf "%03d" $SPLIT_NUM)
    SPLIT_DIR="/efforts/phase${PHASE}/wave${WAVE}/${EFFORT_NAME}-SPLIT-${SPLIT_NAME}"
    
    echo "═══════════════════════════════════════════════════════════════"
    echo "🔧 Creating Split-${SPLIT_NAME} Infrastructure"
    echo "Directory: $SPLIT_DIR"
    echo "Base Branch: $BASE_BRANCH"
    echo "═══════════════════════════════════════════════════════════════"
    
    # Create directory
    mkdir -p "$SPLIT_DIR"
    
    # Clone with correct base
    echo "📦 Cloning from base branch: $BASE_BRANCH"
    
    # 🔴 R309 CRITICAL CHECK: Verify target is NOT SF repo!
    if [[ "$TARGET_REPO_URL" == *"software-factory"* ]]; then
        echo "🔴🔴🔴 R309 VIOLATION: Target URL is Software Factory repo!"
        echo "FATAL ERROR: You're trying to clone SF into itself!"
        echo "Fix target-repo-config.yaml to point to actual project!"
        exit 309
    fi
    
    git clone --branch "$BASE_BRANCH" --sparse "$TARGET_REPO_URL" "$SPLIT_DIR"
    
    cd "$SPLIT_DIR"
    
    # 🔴 R309 POST-CLONE VALIDATION: Ensure we cloned the right thing!
    if [ -f ".claude/CLAUDE.md" ] || [ -f "rule-library/RULE-REGISTRY.md" ]; then
        echo "🔴🔴🔴 R309 VIOLATION: Cloned SF repo instead of target!"
        echo "FATAL ERROR: This is the wrong repository!"
        exit 309
    fi
    echo "✅ R309 VALIDATED: This is TARGET repo (not SF)"
    
    # Set up sparse checkout
    git sparse-checkout init --cone
    git sparse-checkout set pkg/
    
    # Create split branch with proper naming
    SPLIT_BRANCH=$(get_split_branch_name "$EFFORT_NAME" "$SPLIT_NAME")
    git checkout -b "$SPLIT_BRANCH"
    
    # Push to remote
    git push -u origin "$SPLIT_BRANCH"
    
    # Verify remote tracking
    if git branch -vv | grep -q "$SPLIT_BRANCH.*origin/$SPLIT_BRANCH"; then
        echo "✅ Remote tracking configured for $SPLIT_BRANCH"
    else
        echo "❌ FATAL: Remote tracking failed for $SPLIT_BRANCH"
        exit 1
    fi
    
    # Copy split plan from too-large branch
    TOO_LARGE_DIR="/efforts/phase${PHASE}/wave${WAVE}/${EFFORT_NAME}"
    
    # Look for timestamped split plan (per R301)
    SPLIT_PLAN=$(ls -t "$TOO_LARGE_DIR"/SPLIT-PLAN-*-split${SPLIT_NAME}-*.md 2>/dev/null | head -1)
    
    if [ -n "$SPLIT_PLAN" ]; then
        cp "$SPLIT_PLAN" .
        echo "✅ Split plan copied: $(basename "$SPLIT_PLAN")"
    else
        # Fallback: check legacy numbered format
        if [ -f "$TOO_LARGE_DIR/SPLIT-PLAN-${SPLIT_NAME}.md" ]; then
            cp "$TOO_LARGE_DIR/SPLIT-PLAN-${SPLIT_NAME}.md" .
            echo "⚠️ WARNING: Using legacy split plan format (should be timestamped per R301)"
        else
            echo "❌ ERROR: No split plan found for split ${SPLIT_NAME}!"
            echo "   Searched for: SPLIT-PLAN-*-split${SPLIT_NAME}-*.md"
            echo "   Also checked legacy: SPLIT-PLAN-${SPLIT_NAME}.md"
            exit 1
        fi
    fi
    
    # Add metadata to split plan (handle both formats)
    SPLIT_PLAN_LOCAL=$(ls SPLIT-PLAN-*.md | head -1)
    cat >> "$SPLIT_PLAN_LOCAL" << EOF

## 🚨 SPLIT INFRASTRUCTURE METADATA (Added by Orchestrator)
**WORKING_DIRECTORY**: $(pwd)
**BRANCH**: $SPLIT_BRANCH
**REMOTE**: origin/$SPLIT_BRANCH
**BASE_BRANCH**: $BASE_BRANCH
**SPLIT_NUMBER**: $SPLIT_NAME
**CREATED_AT**: $(date '+%Y-%m-%d %H:%M:%S')

### SW Engineer Instructions
1. READ this metadata FIRST
2. cd to WORKING_DIRECTORY above
3. Verify branch matches BRANCH above
4. ONLY THEN proceed with implementation
EOF
    
    # Commit initial setup
    git add -A
    git commit -m "chore: initialize split-${SPLIT_NAME} from $BASE_BRANCH"
    git push
    
    # 🔴🔴🔴 R312: LOCK GIT CONFIG FOR SPLIT ISOLATION 🔴🔴🔴
    echo "🔒 R312: Applying DOUBLE PROTECTION to split git config..."
    
    # Verify .git/config exists
    if [ ! -f .git/config ]; then
        echo "❌ FATAL: No .git/config found in $SPLIT_DIR"
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
        echo "Split isolation compromised - config still writable!"
        exit 312
    fi
    
    # Verify ownership if sudo was available
    CURRENT_OWNER=$(stat -c %U:%G .git/config 2>/dev/null || stat -f %Su:%Sg .git/config)
    if [ "$PROTECTION_LEVEL" = "FULL" ] && [ "$CURRENT_OWNER" != "root:root" ]; then
        echo "⚠️ WARNING: Ownership not changed to root:root (got $CURRENT_OWNER)"
        echo "Protection may be weaker than intended"
    fi
    
    # Create protection marker with split details
    cat > .git/R312_CONFIG_LOCKED << EOF
Timestamp: $(date '+%Y-%m-%d %H:%M:%S')
Locked by: orchestrator
State: CREATE_NEXT_SPLIT_INFRASTRUCTURE
Split: $SPLIT_NAME
Effort: $EFFORT_NAME
Phase: $PHASE
Wave: $WAVE
Base Branch: $BASE_BRANCH
Split Branch: $SPLIT_BRANCH
Protection level: $PROTECTION_LEVEL
Previous ownership: $BEFORE_OWNER
Current ownership: $CURRENT_OWNER
Previous permissions: $BEFORE_PERMS
Current permissions: 444 (readonly)
Purpose: Prevent branch/remote modifications during split work per R312
EOF
    
    echo "✅ R312 ENFORCED: Split config locked"
    echo "   Protection: $PROTECTION_LEVEL"
    echo "   Ownership: $BEFORE_OWNER → $CURRENT_OWNER"
    echo "   Permissions: $BEFORE_PERMS → 444"
    echo "📝 Split isolation guaranteed:"
    echo "   ❌ Cannot switch to other splits"
    echo "   ❌ Cannot pull from original effort"
    echo "   ❌ Cannot merge other changes"
    echo "   ✅ Can only work on assigned split scope"
    
    echo "✅ Split $SPLIT_NAME infrastructure complete with R312 protection"
}
```

## State Transitions

From CREATE_NEXT_SPLIT_INFRASTRUCTURE state:
- **SPAWN_AGENTS** - Infrastructure ready, spawn SW Engineer for split implementation
- **ERROR_RECOVERY** - Infrastructure creation failed

**HOW YOU GOT HERE (Valid paths only):**
- MONITOR_REVIEWS → (detected split needed) → CREATE_NEXT_SPLIT_INFRASTRUCTURE
- NEVER: MONITOR_IMPLEMENTATION → CREATE_NEXT_SPLIT_INFRASTRUCTURE (missing review!)

## Common Violations to Avoid

### 🔴🔴🔴 CRITICAL VIOLATION: ALL SPLITS FROM SAME BASE 🔴🔴🔴

#### ❌ CATASTROPHIC ERROR (What orchestrator is doing WRONG now):
```bash
# WRONG - All splits branch from integration, losing each other's work!
git clone --branch phase1/wave1-integration ... split-001  ✓ (correct)
git clone --branch phase1/wave1-integration ... split-002  ❌ WRONG!
git clone --branch phase1/wave1-integration ... split-003  ❌ WRONG!

# Result: Split-002 doesn't have Split-001's work!
# Result: Split-003 doesn't have Split-001 OR Split-002's work!
# THIS IS CAUSING THE BUG RIGHT NOW!
```

#### ✅ CORRECT SEQUENTIAL CHAINING:
```bash
# RIGHT - Each split builds on the previous one
git clone --branch phase1/wave1-integration ... split-001  ✓
git clone --branch phase1/wave1/effort-split-001 ... split-002  ✓
git clone --branch phase1/wave1/effort-split-002 ... split-003  ✓

# Result: Each split has ALL previous splits' work!
```

### Visual Diagram of the Problem:
```
❌ WRONG (Current Bug):
phase1/wave1-integration
    ├── split-001 (has: base code)
    ├── split-002 (has: base code ONLY - missing split-001!)
    └── split-003 (has: base code ONLY - missing split-001 & 002!)

✅ CORRECT (How it should be):
phase1/wave1-integration
    └── split-001 (has: base code)
        └── split-002 (has: base + split-001)
            └── split-003 (has: base + split-001 + split-002)
```

### ❌ Creating All Splits at Once
```bash
# WRONG - Creating all infrastructure upfront
for split in $(seq 1 $TOTAL_SPLITS); do
    create_split_infrastructure $split
done
```

### ✅ Correct Sequential Creation
```bash
# RIGHT - One split at a time, correct base
if [ $SPLIT_NUM -eq 1 ]; then
    BASE=$(determine_effort_base_branch)  # R308 base for FIRST split only
else
    # CRITICAL: Use PREVIOUS split as base!
    BASE="phase${PHASE}/wave${WAVE}/${EFFORT_NAME}-split-$(printf "%03d" $((SPLIT_NUM - 1)))"
fi
create_single_split_infrastructure "$EFFORT" "$SPLIT_NUM" "$BASE"
```

## Split Tracking Update

```bash
update_split_tracking() {
    local EFFORT="$1"
    local SPLIT_NUM="$2"
    local STATUS="$3"
    
    # Update current split being worked
    yq -i ".split_tracking.\"$EFFORT\".current_split = $SPLIT_NUM" orchestrator-state.yaml
    
    # Add split to tracking
    yq -i ".split_tracking.\"$EFFORT\".splits[$((SPLIT_NUM - 1))].number = $SPLIT_NUM" orchestrator-state.yaml
    yq -i ".split_tracking.\"$EFFORT\".splits[$((SPLIT_NUM - 1))].status = \"$STATUS\"" orchestrator-state.yaml
    yq -i ".split_tracking.\"$EFFORT\".splits[$((SPLIT_NUM - 1))].created_at = \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\"" orchestrator-state.yaml
    
    # Commit state update
    git add orchestrator-state.yaml
    git commit -m "state: created infrastructure for split-$(printf "%03d" $SPLIT_NUM)"
    git push
}
```

## Key Points

1. **Create ONLY ONE split at a time**
2. **Base each split on the PREVIOUS split (except first)**
3. **First split uses R308 incremental base**
4. **Update split tracking immediately**
5. **Transition to SPAWN_AGENTS to implement the split**
6. **After implementation: MUST spawn Code Reviewer for review**
7. **Only create next split if review passes AND more splits needed**
8. **Each split gets FULL review cycle - no shortcuts!**

## Grading Impact

- **Creating all splits at once**: -100% (Violates sequential principle)
- **Wrong base branch**: -75% (Breaks line counting)
- **Not updating tracking**: -50% (Lost state)
- **Not creating infrastructure**: -100% (State failure)

---

**REMEMBER**: This state creates EXACTLY ONE split's infrastructure, then immediately transitions to SPAWN_AGENTS to implement it!
## R322 VIOLATION DETECTION

If you find yourself:
- Starting work for a new state without /continue-orchestrating
- Transitioning without stopping after state file commit
- Continuing after completing state work

**STOP IMMEDIATELY - You are violating R322!**
