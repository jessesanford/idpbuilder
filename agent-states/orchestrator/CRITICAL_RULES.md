# CRITICAL RULES THAT APPLY TO MULTIPLE STATES

These rules are so important they must be available across multiple orchestrator states.

## 🔴🔴🔴 SUPREME LAWS (HIGHEST PRIORITY) 🔴🔴🔴

---
### 🔴 SUPREME LAW #1 (HIGHEST): R234 - Mandatory State Traversal
**Source:** rule-library/R234-mandatory-state-traversal-supreme-law.md
**Criticality:** THE HIGHEST SUPREME LAW - NO OTHER RULE CAN OVERRIDE
**Used in states:** ALL - ESPECIALLY STATE TRANSITIONS
**Penalty:** -100% GRADE (AUTOMATIC FAIL) + IMMEDIATE TERMINATION

**YOU MUST TRAVERSE EVERY STATE IN MANDATORY SEQUENCES - NO SKIPPING EVER!**

Critical Mandatory Sequence:
```
SETUP_EFFORT_INFRASTRUCTURE
    ↓ (MANDATORY)
ANALYZE_CODE_REVIEWER_PARALLELIZATION
    ↓ (MANDATORY)
SPAWN_CODE_REVIEWERS_EFFORT_PLANNING
    ↓ (MANDATORY)
WAITING_FOR_EFFORT_PLANS
    ↓ (MANDATORY)
ANALYZE_IMPLEMENTATION_PARALLELIZATION
    ↓ (MANDATORY)
SPAWN_AGENTS
```

**FORBIDDEN TRANSITIONS:**
- ❌ SETUP_EFFORT_INFRASTRUCTURE → SPAWN_AGENTS (INSTANT FAIL!)
- ❌ Skipping ANY state for "efficiency" or "continuous operation"
- ❌ Rationalizing that R021/R231 allow skipping states

**REMEMBER:** 
- "Continuous operation" = FLOW THROUGH all states, not skip them
- Every state has CRITICAL work that MUST be done
- NO EXCEPTIONS, NO WORKAROUNDS, NO NEGOTIATIONS

---
### 🔴 SUPREME LAW #2: R208 - CD BEFORE SPAWN
**Source:** rule-library/R208-orchestrator-spawn-directory-protocol.md
**Criticality:** SUPREME LAW #2 - NO OTHER RULE CAN OVERRIDE (except R234)
**Used in states:** ALL SPAWN STATES - NO EXCEPTIONS
**Penalty:** -100% GRADE (AUTOMATIC FAIL) + IMMEDIATE TERMINATION

**YOU MUST ALWAYS CD TO THE CORRECT DIRECTORY BEFORE SPAWNING ANY AGENT - NO EXCEPTIONS!**

**MANDATORY SPAWN PROTOCOL:**
```bash
# BEFORE EVERY SPAWN - THIS EXACT SEQUENCE:
1. DETERMINE target directory for agent
2. CD to that directory (MANDATORY - NO EXCEPTIONS!)
3. VERIFY with pwd (MANDATORY - NO SKIPPING!)
4. SPAWN the agent (inherits your current directory)
5. RETURN to orchestrator directory
```

**FORBIDDEN ACTIONS (INSTANT -100% FAILURE):**
- ❌ Spawning without CD'ing first = AUTOMATIC FAIL
- ❌ Assuming agent will CD itself = AUTOMATIC FAIL
- ❌ Using --working-directory instead of CD = AUTOMATIC FAIL
- ❌ Spawning from wrong directory "for efficiency" = AUTOMATIC FAIL
- ❌ Skipping pwd verification = AUTOMATIC FAIL

**THE ONLY ACCEPTABLE SPAWN PATTERN:**
```bash
# SW ENGINEER IN EFFORT:
cd efforts/phase1/wave1/effort-api-types || exit 1
pwd  # MUST verify correct directory
task spawn sw-engineer "implement effort"
cd /workspaces/project  # Return

# CODE REVIEWER IN SPLIT:
cd efforts/phase1/wave1/effort-api-types/split-001 || exit 1
pwd  # MUST verify correct directory
task spawn code-reviewer "review split"
cd /workspaces/project  # Return
```

**REMEMBER:**
- R208 is SUPREME LAW #2 - Cannot be overridden by ANY rule except R234
- "Continuous operation" does NOT mean skip the CD step
- "Efficiency" is NOT an excuse to violate this law
- EVERY spawn MUST follow this protocol - NO EXCEPTIONS
- Agents inherit spawner's directory - wrong dir = total failure

**SPAWNING WITHOUT CD = -100% GRADE = AUTOMATIC FAILURE**

---
### 🔴 SUPREME LAW: R288 - Mandatory State File Commit/Push
**Source:** rule-library/R288-state-file-update-and-commit-protocol.md
**Criticality:** SUPREME LAW - AUTOMATIC FAILURE IF VIOLATED
**Used in states:** ALL

EVERY SINGLE EDIT to orchestrator-state.yaml MUST be IMMEDIATELY:
1. Staged with `git add orchestrator-state.yaml`
2. Committed with `git commit -m "state: description [R288]"`
3. Pushed with `git push`

NO BATCHING! NO DEFERRALS! NO EXCEPTIONS!

```bash
# EVERY state edit MUST follow this pattern:
yq -i '.current_state = "NEW_STATE"' orchestrator-state.yaml
git add orchestrator-state.yaml
git commit -m "state: transition to NEW_STATE [R288]"
git push
```

VIOLATION = AUTOMATIC ORCHESTRATION FAILURE!

---
### 🔴 SUPREME LAW: R288 - State File Update and Commit Protocol
**Source:** rule-library/R288-state-file-update-and-commit-protocol.md
**Criticality:** SUPREME LAW - State file is single source of truth
**Used in states:** ALL

EVERY state transition MUST update orchestrator-state.yaml IMMEDIATELY:
- Update current_state field
- Add transition timestamp
- Record transition reason
- Update any relevant state-specific fields

```bash
# MANDATORY on every transition:
update_orchestrator_state "NEW_STATE" "reason for transition"
mark_wave_complete "$PHASE" "$WAVE"  # When applicable
```

---
### 🔴 SUPREME LAW: R217 - Immediate Rule Reloading
**Source:** agent-states/orchestrator/SUPREME_LAW_R217.md
**Criticality:** SUPREME LAW - NO DEFERRALS
**Used in states:** ALL

After EVERY state transition, IMMEDIATELY read rules before ANY other action:
1. Read agent-states/orchestrator/SUPREME_LAW_R217.md
2. Read orchestrator.md
3. Read SOFTWARE-FACTORY-STATE-MACHINE.md
4. Read CRITICAL_RULES.md
5. Read agent-states/orchestrator/{NEW_STATE}/rules.md

NO "BUT FIRST"! NO EXCEPTIONS!

---
### 🔴 SUPREME LAW: R021 - Orchestrator Never Stops
**Source:** rule-library/R021-orchestrator-never-stops.md
**Criticality:** SUPREME LAW - AUTOMATIC FAILURE IF VIOLATED
**Used in states:** ALL

THE ORCHESTRATOR MUST NEVER STOP WORKING UNTIL:
- All tasks are complete, OR
- User explicitly says "stop", OR
- HARD_STOP state is reached

**FORBIDDEN REASONS TO STOP:**
- ❌ Context constraints ("context getting long")
- ❌ Time concerns ("been working for hours")
- ❌ Complexity ("too many tasks")
- ❌ Summary requests ("let me provide a summary")
- ❌ Autonomous judgment ("I think it's better to stop")

**VIOLATION = AUTOMATIC -100% FAILURE**

```bash
# Before ANY stop decision:
if ! all_tasks_complete() && ! user_said_stop() && ! hard_stop_state(); then
    echo "🔴🔴🔴 R021 VIOLATION PREVENTED 🔴🔴🔴"
    echo "NO VALID REASON TO STOP - CONTINUING..."
    # MUST CONTINUE WORKING
fi
```

**THE ORCHESTRATOR'S CONFIDENCE:**
- "Context concerns don't matter - R287 saves my TODOs every 10 messages."
- "Compaction doesn't scare me - R288 preserves my state continuously."
- "Recovery is automatic - I read state, load TODOs, and continue."
- "The system protects my progress - I can work indefinitely."
- "Only completion or user 'stop' matters."

**CONTEXT RECOVERY PROTOCOL:**
```bash
# After any compaction event:
1. Check /tmp/compaction_marker.txt
2. Read todos/orchestrator-*.todo (latest)
3. Load TODOs with TodoWrite tool
4. Read orchestrator-state.yaml
5. Continue from current_state
# SEAMLESS RECOVERY - NO PROGRESS LOST
```

---
### 🚨 RULE R203 - State-Aware Startup Protocol (NO EXCUSES!)
**Source:** rule-library/RULE-REGISTRY.md#R203
**Criticality:** BLOCKING - Must follow startup sequence
**Used in states:** ALL (on every startup)

MANDATORY STARTUP SEQUENCE:
1. Read orchestrator.md (core config)
2. Read SOFTWARE-FACTORY-STATE-MACHINE.md
3. Read CRITICAL_RULES.md
4. Determine current state from orchestrator-state.yaml
5. Read state-specific rules from agent-states/orchestrator/{STATE}/rules.md
6. Acknowledge all rules
7. Only then proceed with orchestration

⚠️⚠️⚠️ **NO EXCUSES POLICY** ⚠️⚠️⚠️
❌ **NOT ACCEPTABLE**: "I already read the state rules earlier"
❌ **NOT ACCEPTABLE**: "Let me acknowledge what I read before"
❌ **NOT ACCEPTABLE**: "I remember the rules from earlier"
❌ **NOT ACCEPTABLE**: "Per R203, I must re-acknowledge"
✅ **ONLY ACCEPTABLE**: Actually use the READ tool for EVERY file

**YOU MUST PHYSICALLY READ THE FILES WITH THE READ TOOL!**
- NO shortcuts
- NO assumptions  
- NO cached knowledge
- NO "I already know this"

**VIOLATION = AUTOMATIC FAILURE**

---
### 🚨 RULE R206 - State Machine Validation
**Source:** rule-library/RULE-REGISTRY.md#R206
**Criticality:** ABSOLUTE - Every transition must be validated
**Used in states:** ALL

EVERY state transition MUST be validated against SOFTWARE-FACTORY-STATE-MACHINE.md:
```bash
validate_state_transition() {
    local from_state="$1"
    local to_state="$2"
    
    # Validate against SOFTWARE-FACTORY-STATE-MACHINE.md
    # Check transition is allowed
    # Verify prerequisites met
    # Update orchestrator-state.yaml
}
```

---
### 🚨 RULE R216 - Bash Execution Syntax
**Source:** rule-library/R216-bash-execution-syntax.md
**Criticality:** CRITICAL - Proper command formatting
**Used in states:** ALL

- Use multi-line format when executing bash commands
- If single-line needed, use semicolons (`;`) between statements
- Do NOT include backslashes (`\`) from documentation in actual execution
- Backslashes are ONLY for documentation line continuation

---
### 🔴🔴🔴 RULE R217 - SUPREME LAW: IMMEDIATE POST-TRANSITION RULE RELOADING 🔴🔴🔴
**Source:** rule-library/RULE-REGISTRY.md#R217
**Criticality:** SUPREME LAW - ABSOLUTE HIGHEST PRIORITY
**Used in states:** ALL (IMMEDIATELY after EVERY state transition)

**⚠️ THIS IS NON-NEGOTIABLE - NO DEFERRALS ALLOWED ⚠️**

🚫 **ABSOLUTELY FORBIDDEN EXCUSES:**
- ❌ "I already read the MONITOR state rules earlier, but per R217 I must re-acknowledge them"
- ❌ "I read these rules in a previous message"
- ❌ "Let me complete the acknowledgment"
- ❌ "I'm familiar with these rules"
- ❌ "As I recall from earlier"

✅ **THE ONLY WAY**: USE THE READ TOOL. PERIOD.

**WE ALLOW NO EXCUSES!** Every transition = Fresh READ of ALL files!

AFTER EVERY STATE TRANSITION - DO THIS IMMEDIATELY OR FAIL:
```bash
perform_state_transition() {
    local OLD_STATE="$1"
    local NEW_STATE="$2"
    
    # Step 1: Validate transition
    validate_state_transition "$OLD_STATE" "$NEW_STATE"
    
    # Step 2: Update state
    update_state "current_state" "$NEW_STATE"
    
    # Step 3: 🔴🔴🔴 STOP EVERYTHING - READ RULES NOW 🔴🔴🔴
    echo "🔴🔴🔴 R217 SUPREME LAW ACTIVATED 🔴🔴🔴"
    echo "State transition complete. HALTING ALL WORK."
    echo "MUST read rules IMMEDIATELY - no deferrals allowed!"
    
    # MANDATORY IMMEDIATE READS:
    # READ: agent-states/orchestrator/SUPREME_LAW_R217.md
    # READ: .claude/agents/orchestrator.md
    # READ: SOFTWARE-FACTORY-STATE-MACHINE.md
    # READ: agent-states/orchestrator/CRITICAL_RULES.md
    # READ: agent-states/orchestrator/$NEW_STATE/rules.md
    
    # Step 4: Only after ALL reads complete
    echo "✅ R217 COMPLETE: Rules reloaded for state: $NEW_STATE"
    echo "NOW authorized to proceed with work."
}
```

**FORBIDDEN PHRASES AFTER STATE TRANSITION:**
- ❌ "But first, let me..."
- ❌ "According to R217, I must... But first..."
- ❌ "I'll acknowledge rules after..."

**VIOLATION = AUTOMATIC GRADING FAILURE**

## 🚨🚨🚨 RULE R209 - EFFORT DIRECTORY ISOLATION PROTOCOL (MISSION CRITICAL!)
**Source:** rule-library/R209-effort-directory-isolation-protocol.md  
**Criticality:** MISSION CRITICAL - SW Engineers MUST stay in effort directories  
**Used in states:** SPAWN_CODE_REVIEWERS_EFFORT_PLANNING, WAITING_FOR_EFFORT_PLANS, SPAWN_AGENTS

### MUST INJECT METADATA AND PUSH TO REMOTE!

YOU MUST ADD DIRECTORY/BRANCH METADATA TO IMPLEMENTATION PLANS AND PUSH:
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

## 🚨🚨🚨 RULE R218 - Parallel Code Reviewer Spawning
**Source:** rule-library/R218-orchestrator-parallel-code-reviewer-spawning.md
**Criticality:** MANDATORY - Must read wave plan parallelization headers
**Used in states:** SPAWN_CODE_REVIEWERS_EFFORT_PLANNING

The orchestrator MUST:
1. **USE READ TOOL** to read wave implementation plan
2. **ACKNOWLEDGE** reading the parallelization headers
3. **STATE THE FILE PATH** of the wave plan
4. **SPAWN** Code Reviewers according to parallelization metadata

```bash
# MANDATORY OUTPUT:
echo "📖 R218: Using Read tool to read wave plan..."
READ: phase-plans/PHASE-X-WAVE-Y-IMPLEMENTATION-PLAN.md

echo "📖 R218: I have READ the parallelization headers from:"
echo "   phase-plans/PHASE-X-WAVE-Y-IMPLEMENTATION-PLAN.md"

echo "✅ Spawning Code Reviewers as allowed by the parallelization headers"
```

---
### 🚨 RULE R181 - Orchestrator Workspace Setup Responsibility
**Source:** rule-library/RULE-REGISTRY.md#R181
**Criticality:** CRITICAL - Workspace must be ready before agents spawn
**Used in states:** SETUP_EFFORT_INFRASTRUCTURE
The orchestrator MUST prepare complete working environments:
1. Create effort directory structure
2. Perform FULL single-branch clone of target repository (R271 SUPREME LAW)
3. Create properly named git branch
4. Verify workspace before spawning agent

---
### 🔴🔴🔴 RULE R271 - Single-Branch Full Checkout Protocol (SUPREME LAW)
**Source:** rule-library/R271-single-branch-full-checkout.md
**Criticality:** SUPREME - Supersedes R193, NO SPARSE CHECKOUTS
**Used in states:** SETUP_EFFORT_INFRASTRUCTURE

Each effort needs its own FULL single-branch clone:
```bash
# R271: THINK about base branch first
echo "🧠 THINKING: What base branch should this effort use?"
BASE_BRANCH="main"  # Or from dependencies

EFFORT_DIR="efforts/phase1/wave1/core-types"
mkdir -p "$(dirname "$EFFORT_DIR")"

# FULL single-branch clone (NO SPARSE!)
git clone --single-branch --branch "$BASE_BRANCH" [URL] "$EFFORT_DIR"
cd "$EFFORT_DIR"
```

---
### 🚨 RULE R183 - Branch Creation Protocol
**Source:** rule-library/RULE-REGISTRY.md#R183
**Criticality:** CRITICAL - Required for remote tracking
**Used in states:** SETUP_EFFORT_INFRASTRUCTURE

Create branch with project prefix using branch naming helpers:
```bash
# Source the branch naming helpers
source utilities/branch-naming-helpers.sh

# Get properly formatted branch name with project prefix
EFFORT_BRANCH=$(get_effort_branch_name "$PHASE" "$WAVE" "$EFFORT")
# Example with prefix: tmc-workspace/phase1/wave1/core-types
# Example without: phase1/wave1/core-types

git checkout -b "$EFFORT_BRANCH"
git push -u origin "$EFFORT_BRANCH"
```

Or manually if helpers unavailable:
```bash
PROJECT_PREFIX=$(yq '.branch_naming.project_prefix' target-repo-config.yaml)
if [ -n "$PROJECT_PREFIX" ] && [ "$PROJECT_PREFIX" != "null" ]; then
    BRANCH="${PROJECT_PREFIX}/phase${PHASE}/wave${WAVE}/${EFFORT}"
else
    BRANCH="phase${PHASE}/wave${WAVE}/${EFFORT}"
fi
```

---
### 🚨 RULE R184 - Effort Branch Naming
**Source:** rule-library/RULE-REGISTRY.md#R184
**Criticality:** CRITICAL - Consistent naming required
**Used in states:** SETUP_EFFORT_INFRASTRUCTURE
- Format: `[prefix/]phase{X}/wave{Y}/effort-{name}`
- Use hyphens for multi-word names
- Never use spaces or special characters

---
### 🚨 RULE R185 - Workspace Verification Before Spawn
**Source:** rule-library/RULE-REGISTRY.md#R185
**Criticality:** BLOCKING - Cannot spawn without verified workspace
**Used in states:** SETUP_EFFORT_INFRASTRUCTURE, SPAWN_AGENTS
```bash
verify_effort_workspace() {
    local effort_dir="$1"
    
    [ ! -d "$effort_dir" ] && echo "❌ Directory missing" && return 1
    [ ! -d "$effort_dir/.git" ] && echo "❌ Not a git repo" && return 1
    
    cd "$effort_dir"
    local branch=$(git branch --show-current)
    [[ "$branch" != *"phase"*"wave"*"effort"* ]] && echo "❌ Wrong branch" && return 1
    
    echo "✅ Workspace ready"
    return 0
}
```

---
### 🚨 RULE R191 - Target Repository Configuration
**Source:** rule-library/RULE-REGISTRY.md#R191
**Criticality:** BLOCKING - Cannot proceed without config
**Used in states:** INIT, PLANNING, SETUP_EFFORT_INFRASTRUCTURE
```bash
load_target_config() {
    [ ! -f "target-repo-config.yaml" ] && echo "❌ Config missing!" && exit 1
    
    export TARGET_REPO_URL=$(yq '.target_repository.url' target-repo-config.yaml)
    export BASE_BRANCH=$(yq '.target_repository.base_branch' target-repo-config.yaml)
    export PROJECT_PREFIX=$(yq '.branch_naming.project_prefix' target-repo-config.yaml)
}
```

---
### 🚨 RULE R192 - Repository Separation
**Source:** rule-library/RULE-REGISTRY.md#R192
**Criticality:** CRITICAL - Maintain clean separation
**Used in states:** ALL
- SF Instance: Planning, configs, state (NO CODE!)
- Target Clones: Source code, tests, builds (NO CONFIGS!)

---
### 🔴🔴🔴 RULE R271 - Single-Branch Full Checkout (SUPREME LAW)
**Source:** rule-library/R271-single-branch-full-checkout.md
**Criticality:** SUPREME - Supersedes old R193
**Used in states:** SETUP_EFFORT_INFRASTRUCTURE
Create FULL single-branch clone for EVERY effort before spawning agents.
NO SPARSE CHECKOUTS ALLOWED!

---
### 🚨 RULE R194 - Remote Tracking Required
**Source:** rule-library/RULE-REGISTRY.md#R194
**Criticality:** CRITICAL - Enables collaboration
**Used in states:** SETUP_EFFORT_INFRASTRUCTURE
Every branch must have remote tracking set up.

---
### 🚨 RULE R195 - Push Verification
**Source:** rule-library/RULE-REGISTRY.md#R195
**Criticality:** CRITICAL - Verify before spawning
**Used in states:** SETUP_EFFORT_INFRASTRUCTURE
Verify push succeeded before spawning agents.

---
### 🚨 RULE R196 - Integration Branch Management
**Source:** rule-library/RULE-REGISTRY.md#R196
**Criticality:** CRITICAL - Required for wave/phase completion
**Used in states:** WAVE_COMPLETE, SUCCESS

Create integration branches at wave and phase boundaries WITH PROJECT PREFIX:
```bash
# Source branch naming helpers
source utilities/branch-naming-helpers.sh

# Wave integration (includes project prefix)
WAVE_INTEGRATION=$(get_wave_integration_branch_name "$PHASE" "$WAVE")
# Example: tmc-workspace/phase1/wave1-integration
git checkout -b "$WAVE_INTEGRATION"

# Merge effort branches (also with prefix)
for effort in effort1 effort2 effort3; do
    EFFORT_BRANCH=$(get_effort_branch_name "$PHASE" "$WAVE" "$effort")
    git merge "origin/$EFFORT_BRANCH" --no-ff
done

# Phase integration (includes project prefix)
PHASE_INTEGRATION=$(get_phase_integration_branch_name "$PHASE")
# Example: tmc-workspace/phase1-integration
git checkout -b "$PHASE_INTEGRATION"
git merge "$WAVE_INTEGRATION" # And other wave integrations
```

---
### 🚨 RULE R197 - Merge Order Protocol
**Source:** rule-library/RULE-REGISTRY.md#R197
**Criticality:** CRITICAL - Prevents conflicts
**Used in states:** WAVE_COMPLETE

Merge efforts in dependency order, not alphabetical.

---
### 🚨 RULE R198 - Conflict Resolution
**Source:** rule-library/RULE-REGISTRY.md#R198
**Criticality:** CRITICAL - Maintain functionality
**Used in states:** WAVE_COMPLETE, ERROR_RECOVERY

Resolve conflicts preserving all functionality from both branches.

---
### 🚨 RULE R199 - Integration Testing
**Source:** rule-library/RULE-REGISTRY.md#R199
**Criticality:** BLOCKING - Must pass before proceeding
**Used in states:** WAVE_COMPLETE

Run full test suite on integration branches before proceeding.

---
### 🚨 RULE R001 - Never Implement Code
**Source:** rule-library/RULE-REGISTRY.md#R001
**Criticality:** BLOCKING - Automatic failure
**Used in states:** ALL

Orchestrator NEVER writes code. Delegate ALL implementation.

---
### 🚨 RULE R010 - State Machine Authority
**Source:** rule-library/RULE-REGISTRY.md#R010
**Criticality:** SUPREME - Highest authority
**Used in states:** ALL

SOFTWARE-FACTORY-STATE-MACHINE.md is absolute authority.

---
### 🚨 RULE R171-180 - Workspace Management Suite
**Source:** rule-library/RULE-REGISTRY.md#R171
**Criticality:** CRITICAL - Foundation for all work
**Used in states:** SETUP_EFFORT_INFRASTRUCTURE

Complete suite of workspace management rules:
- R171: Clean workspace for each effort
- R172: Proper directory structure
- R173: Git configuration
- R174: Branch protection
- R175: Access permissions
- R176: Infrastructure verification (detailed in state file)
- R177: Workspace isolation
- R178: Resource allocation
- R179: Cleanup protocols
- R180: Backup procedures

---
### 🚨 RULE R186 - Monitoring Frequency
**Source:** rule-library/RULE-REGISTRY.md#R186
**Criticality:** CRITICAL - Catch issues early
**Used in states:** MONITOR

Check all spawned agents every 5 messages.

---
### 🚨 RULE R202 - State File Management
**Source:** rule-library/RULE-REGISTRY.md#R202
**Criticality:** CRITICAL - Maintain consistency
**Used in states:** ALL

Update orchestrator-state.yaml after EVERY state transition.

---
### 🔴🔴🔴 RULE R288 - MANDATORY STATE FILE UPDATES AT EVERY TRANSITION
**Source:** rule-library/R288-state-file-update-and-commit-protocol.md
**Criticality:** SUPREME LAW - AUTOMATIC FAILURE IF VIOLATED
**Used in states:** ALL (EVERY state transition)

**NO STATE TRANSITION WITHOUT STATE FILE UPDATE!**

```bash
# MANDATORY for EVERY transition:
update_orchestrator_state() {
    local NEW_STATE="$1"
    local REASON="$2"
    
    # Update IMMEDIATELY - not later!
    yq -i ".state_machine.current_state = \"$NEW_STATE\"" orchestrator-state.yaml
    yq -i ".state_machine.previous_state = \"$OLD_STATE\"" orchestrator-state.yaml
    yq -i ".state_machine.transition_time = \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\"" orchestrator-state.yaml
    yq -i ".state_machine.transition_reason = \"$REASON\"" orchestrator-state.yaml
}

# MANDATORY for WAVE_COMPLETE:
mark_wave_complete() {
    local PHASE="$1"
    local WAVE="$2"
    
    yq -i ".waves_completed.phase${PHASE}.wave${WAVE}.completed_at = \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\"" orchestrator-state.yaml
    yq -i ".waves_completed.phase${PHASE}.wave${WAVE}.status = \"COMPLETE\"" orchestrator-state.yaml
    # ... all required fields
}
```

**State-specific required updates:**
- WAVE_COMPLETE: Must add waves_completed entry
- INTEGRATION: Must add current_integration entry
- SUCCESS: Must add phase_completion entry
- ERROR_RECOVERY: Must add error_context entry
- SPAWN_AGENTS: Must add agents_spawned entries

**The state file is the single source of truth!**

---
### 🚨 RULE R204 - Error Recovery Protocol
**Source:** rule-library/RULE-REGISTRY.md#R204
**Criticality:** CRITICAL - Handle failures gracefully
**Used in states:** ERROR_RECOVERY

Document error, save state, attempt recovery, escalate if needed.

---
### 🚨 RULE R205 - Progress Tracking
**Source:** rule-library/RULE-REGISTRY.md#R205
**Criticality:** INFO - Best practice
**Used in states:** MONITOR

Update progress every 5 orchestrator messages.

---
### 🚨 RULE R208 - Directory Context for Spawning
**Source:** rule-library/RULE-REGISTRY.md#R208
**Criticality:** CRITICAL - Agents need correct context
**Used in states:** SPAWN_AGENTS, SPAWN_CODE_REVIEWERS_EFFORT_PLANNING

CD to effort directory before spawning agent:
```bash
cd efforts/phase${PHASE}/wave${WAVE}/${EFFORT}
pwd  # Verify location
Task: spawn-agent
```

---
### 🚨 RULE R212 - Wave Completion Requirements
**Source:** rule-library/RULE-REGISTRY.md#R212
**Criticality:** BLOCKING - Cannot proceed without
**Used in states:** WAVE_COMPLETE

Before marking wave complete:
1. All efforts implemented and reviewed
2. All size limits verified (<800 lines)
3. All tests passing
4. Integration branch created
5. No blocking issues

---
### 🚨 RULE R213 - Phase Completion Requirements
**Source:** rule-library/RULE-REGISTRY.md#R213
**Criticality:** BLOCKING - Cannot proceed without
**Used in states:** SUCCESS

Before marking phase complete:
1. All waves integrated
2. Phase integration branch created
3. Architect review passed
4. All documentation updated
5. No technical debt

---
### 🔴 RULE R254 - Sequential Split Fix Processing
**Source:** rule-library/R254-sequential-split-fix-processing.md
**Criticality:** BLOCKING - Parallel split fixes = AUTOMATIC FAILURE
**Used in states:** MONITOR, FIX_ISSUES, SPAWN_AGENTS

⚠️⚠️⚠️ **NEVER spawn SW Engineers in parallel for review fixes on splits!** ⚠️⚠️⚠️

SPLIT FIX SEQUENCING REQUIREMENTS:
1. Process split fixes ONE AT A TIME
2. Complete each split's fix-review cycle before starting next
3. NEVER spawn multiple SW Engineers for splits simultaneously
4. Wait for each SW Engineer to complete before spawning next
5. Violation = AUTOMATIC ORCHESTRATION FAILURE

```bash
# CORRECT: Sequential processing
for split in "${SPLITS[@]}"; do
    spawn_sw_engineer_for_split "$split"
    wait_for_completion "$split"
    spawn_code_reviewer "$split"  
    wait_for_review "$split"
done

# WRONG: Parallel processing (AUTOMATIC FAILURE!)
for split in "${SPLITS[@]}"; do
    spawn_sw_engineer_for_split "$split" &  # NO! Parallel = FAILURE!
done
```

---
### 🔴🔴🔴 RULE R255 - Post-Agent Work Verification with Intelligent Recovery
**Source:** rule-library/R255-POST-AGENT-WORK-VERIFICATION.md
**Criticality:** BLOCKING - Wrong location = SALVAGE or DELETE AND RESTART
**Used in states:** MONITOR, WAITING_FOR_EFFORT_PLANS, all post-completion

⚠️⚠️⚠️ **VERIFY EVERY AGENT'S WORK OR FACE AUTOMATIC FAILURE!** ⚠️⚠️⚠️

AFTER EVERY AGENT COMPLETION, YOU MUST:
1. **VERIFY** work is in `/efforts/phase{X}/wave{Y}/{effort}` directory
2. **VERIFY** work is committed to correct branch
3. **VERIFY** work is pushed to remote
4. **TRY TO SALVAGE** if verification fails
5. **DELETE AND RESTART** only if salvage impossible

```bash
# MANDATORY after EVERY agent reports complete
verify_agent_work() {
    EXPECTED_DIR="efforts/phase${PHASE}/wave${WAVE}/${EFFORT}"
    EXPECTED_BRANCH="phase${PHASE}/wave${WAVE}/${EFFORT}"
    
    # Check directory
    if [ ! -d "$EXPECTED_DIR" ]; then
        echo "❌ R255 VIOLATION: Wrong directory!"
        attempt_salvage_then_recover  # TRY SALVAGE FIRST
        return 1
    fi
    
    cd "$EXPECTED_DIR"
    
    # Check branch
    if [ "$(git branch --show-current)" != "$EXPECTED_BRANCH" ]; then
        echo "❌ R255 VIOLATION: Wrong branch!"
        attempt_salvage_then_recover  # TRY SALVAGE FIRST
        return 1
    fi
    
    # Check pushed
    if ! git branch -r | grep -q "origin/$EXPECTED_BRANCH"; then
        echo "❌ R255 VIOLATION: Not pushed!"
        attempt_salvage_then_recover  # TRY SALVAGE FIRST
        return 1
    fi
}
```

**INTELLIGENT RECOVERY PROTOCOL:**
1. **ASSESS** if work is salvageable
2. **SALVAGE** via: move files, cherry-pick, or commit/push
3. **DELETE** only if salvage impossible
4. **RECREATE** clean infrastructure
5. **RESPAWN** with ULTRA-SPECIFIC path instructions

**SALVAGEABLE:** Wrong directory, wrong branch, uncommitted work
**NOT SALVAGEABLE:** Wrong repo, merge conflicts, wrong effort code

**GRADING:**
- Not verifying = -100% AUTOMATIC FAILURE
- Each violation = -20% to -100% escalating

---
### 🔴 RULE R250 - Integration Must Use Separate Target Repo Clone
**Source:** rule-library/R250-INTEGRATION-ISOLATION.md
**Criticality:** BLOCKING - AUTOMATIC FAILURE IF VIOLATED
**Used in states:** INTEGRATION

**THE SOFTWARE FACTORY INSTANCE IS NOT THE TARGET REPOSITORY!**

Integration MUST happen in a FRESH CLONE of the target repository:
```bash
# ✅ CORRECT: Integration in new target repo clone WITH PROPER BRANCH NAMING
source utilities/branch-naming-helpers.sh
mkdir -p integrations/phase${PHASE}/wave${WAVE}
cd integrations/phase${PHASE}/wave${WAVE}
git clone "$TARGET_REPO_URL" integration-workspace
cd integration-workspace

# Use helper function to get branch name with project prefix
INTEGRATION_BRANCH=$(get_wave_integration_branch_name "$PHASE" "$WAVE")
# Example with prefix: tmc-workspace/phase1/wave1-integration
# Example without: phase1/wave1-integration
git checkout -b "$INTEGRATION_BRANCH"

# ❌ WRONG: Never integrate in SF instance directory!
# The SF instance has orchestrator-state.yaml, NOT source code!
```

The Software Factory instance orchestrates work.
The Target Repository contains the actual code.
NEVER confuse the two!