# Orchestrator - SPAWN_INTEGRATION_AGENT State Rules

# PRIMARY DIRECTIVES

You MUST read and acknowledge these rules:

1. **R006** - Orchestrator cannot write code (BLOCKING)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R006-orchestrator-never-writes-code.md`

2. **R361** - Integration Conflict Resolution Only (SUPREME LAW)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R361-integration-conflict-resolution-only.md`
   - Integration Agent can ONLY resolve conflicts, NO new code/packages

3. **R362** - No Architectural Rewrites Without Approval (SUPREME LAW #7)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R362-no-architectural-rewrites.md`
   - Integration Agent MUST NOT change architecture, remove libraries, or deviate from plan

4. **R269** - WAVE Integration Merge Plan Protocol (SUPREME LAW)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R269-WAVE-integration-merge-plan-protocol.md`

5. **R270** - PHASE Integration Merge Plan Protocol (SUPREME LAW)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R270-PHASE-integration-merge-plan-protocol.md`

6. **R287** - TODO Persistence Requirements (SUPREME LAW)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R287-todo-persistence-comprehensive.md`

7. **R288** - State File Update Requirements (SUPREME LAW)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R288-state-file-update-requirements.md`

8. **R304** - Mandatory Line Counter Usage (SUPREME LAW)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R304-mandatory-line-counter-usage.md`

9. **R322** - Mandatory Stop After Spawn States (SUPREME LAW)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R322-mandatory-stop-after-spawn.md`

10. **R324** - State Transition Validation (SUPREME LAW)
    - File: `$CLAUDE_PROJECT_DIR/rule-library/R324-state-transition-validation.md`


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

## 🔴🔴🔴 R322 MANDATORY: STOP BEFORE STATE TRANSITIONS 🔴🔴🔴

**CRITICAL REQUIREMENT PER R322:**
After spawning ANY agents in this state, you MUST:
1. Record what was spawned in state file
2. Save TODOs per R287
3. Commit and push state changes
4. Display stop message with continuation instructions
5. EXIT immediately with code 0

**VIOLATION = AUTOMATIC -100% FAILURE**

See: `$CLAUDE_PROJECT_DIR/rule-library/R322-mandatory-stop-before-state-transitions.md`

---


## 🔴🔴🔴 STOP! STATE RULE READING IS ABSOLUTELY FIRST! 🔴🔴🔴

**YOU HAVE ENTERED SPAWN_INTEGRATION_AGENT STATE - YOU MUST READ AND ACKNOWLEDGE ALL STATE RULES BEFORE DOING ANY STATE WORK!**

## 🔴🔴🔴 R290 VERIFICATION REQUIREMENT 🔴🔴🔴

**R290 ENFORCEMENT: CREATE VERIFICATION MARKER AFTER READING**

After reading and acknowledging all state rules, you MUST create a verification marker:

```bash
# MANDATORY: Create verification marker after reading rules
touch .state_rules_read_orchestrator_SPAWN_INTEGRATION_AGENT
echo "$(date +%s) - Rules read and acknowledged for SPAWN_INTEGRATION_AGENT" > .state_rules_read_orchestrator_SPAWN_INTEGRATION_AGENT
```

**FAILURE TO CREATE MARKER = AUTOMATIC -100% PENALTY**

The system will check for this marker. No marker = Immediate failure.

### ❌ DO NOT DO ANY SPAWN_INTEGRATION_AGENT WORK UNTIL RULES ARE READ:
- ❌ Start spawn integration specialist
- ❌ Start coordinate merging
- ❌ Start manage integration
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
   ❌ WRONG: "I acknowledge all SPAWN_INTEGRATION_AGENT rules"
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

### ✅ CORRECT PATTERN FOR SPAWN_INTEGRATION_AGENT:
```
1. READ: $CLAUDE_PROJECT_DIR/rule-library/[first-rule-file].md
2. "I acknowledge [Rule#] - [Rule Name]: [Brief description]"
3. READ: $CLAUDE_PROJECT_DIR/rule-library/[second-rule-file].md  
4. "I acknowledge [Rule#] - [Rule Name]: [Brief description]"
### 🚨 NO WORK UNTIL ACKNOWLEDGMENT COMPLETE 🚨
**You may NOT begin ANY SPAWN_INTEGRATION_AGENT work until:**
### 🚨 FAILURE TO READ STATE RULES FIRST = IMMEDIATE EXIT 🚨
**If you do ANY SPAWN_INTEGRATION_AGENT work before reading and acknowledging rules:**
- **STOP ALL WORK IMMEDIATELY**
- **EXIT WITH FAILURE STATUS**
- **YOU HAVE VIOLATED STATE COMPLIANCE**

**THE SYSTEM IS MONITORING YOUR READ TOOL CALLS!**

## 🔴🔴🔴 R208 SUPREME LAW: CD BEFORE SPAWN 🔴🔴🔴

**YOU MUST CD TO INTEGRATION DIRECTORY BEFORE SPAWNING THE INTEGRATION AGENT!**
- Violation = -100% GRADE = AUTOMATIC FAILURE
- NO EXCEPTIONS, NO SHORTCUTS, NO WORKAROUNDS

## State Definition
The orchestrator spawns an Integration Agent to execute the merge plan created by Code Reviewer. The merge plan must already exist.

## Required Actions

### 1. Verify Merge Plan Exists
```bash
# Check that Code Reviewer completed the merge plan
PHASE=$(jq '.current_phase' orchestrator-state.json)
WAVE=$(jq '.current_wave' orchestrator-state.json)
INTEGRATION_DIR="/efforts/phase${PHASE}/wave${WAVE}/integration-workspace"

cd "$INTEGRATION_DIR"

if [ ! -f "WAVE-MERGE-PLAN.md" ]; then
    echo "❌ Cannot spawn Integration Agent - no merge plan!"
    echo "🔍 Code Reviewer must complete WAVE-MERGE-PLAN.md first"
    exit 1
fi

echo "✅ Found WAVE-MERGE-PLAN.md"

# 🔴🔴🔴 R312 EXCEPTION: UNLOCK CONFIG FOR INTEGRATION 🔴🔴🔴
echo ""
echo "🔓 R312 EXCEPTION: Unlocking git config for INTEGRATION agent"
echo "Integration agents NEED to pull from multiple branches"

# Check if config is locked
if [ ! -w .git/config ]; then
    echo "📋 Config is currently locked (as expected for efforts)"
    
    # Store current permissions and ownership for audit
    BEFORE_PERMS=$(stat -c %a .git/config 2>/dev/null || stat -f %A .git/config)
    BEFORE_OWNER=$(stat -c %U:%G .git/config 2>/dev/null || stat -f %Su:%Sg .git/config)
    
    # Check current ownership
    CURRENT_OWNER=$(stat -c %U:%G .git/config 2>/dev/null || stat -f %Su:%Sg .git/config)
    
    # Unlock for integration work - handle root ownership
    if [ "$CURRENT_OWNER" = "root:root" ]; then
        # Need sudo to change from root ownership
        if command -v sudo >/dev/null 2>&1; then
            echo "🔓 Restoring user ownership from root..."
            sudo chown $(id -u):$(id -g) .git/config
            sudo chmod 644 .git/config
        else
            echo "❌ ERROR: Config is root-owned but sudo not available!"
            echo "Cannot unlock config for integration"
            exit 312
        fi
    else
        # Simple permission change
        chmod 644 .git/config
    fi
    
    # Verify unlock succeeded
    if [ ! -w .git/config ]; then
        echo "❌ ERROR: Failed to unlock config for integration!"
        echo "Integration agent needs writable config to merge branches"
        exit 312
    fi
    
    # Create exception marker
    cat > .git/R312_INTEGRATION_EXCEPTION << EOF
Timestamp: $(date '+%Y-%m-%d %H:%M:%S')
Unlocked by: orchestrator
State: SPAWN_INTEGRATION_AGENT
Phase: ${PHASE}
Wave: ${WAVE}
Previous ownership: $BEFORE_OWNER
Current ownership: $(stat -c %U:%G .git/config 2>/dev/null || echo 'unknown')
Previous permissions: $BEFORE_PERMS
Current permissions: 644 (writable)
Purpose: Integration requires ability to merge from multiple branches
EOF
    
    echo "✅ R312 Exception Applied: Config unlocked for integration"
    echo "   Owner: $BEFORE_OWNER → $(stat -c %U:%G .git/config 2>/dev/null || echo 'unknown')"
    echo "   Permissions: $BEFORE_PERMS → 644"
    echo "📝 Integration agent can now:"
    echo "   ✅ Pull from multiple effort branches"
    echo "   ✅ Create integration branches"
    echo "   ✅ Merge efforts together"
else
    echo "✅ Config already writable (integration workspace)"
fi

# Quick validation of merge plan
MERGE_COUNT=$(grep -c "git merge origin/" WAVE-MERGE-PLAN.md || echo "0")
echo "📊 Merge plan contains $MERGE_COUNT merge operations"

if [[ $MERGE_COUNT -eq 0 ]]; then
    echo "⚠️ Warning: No merge commands found in plan!"
fi
```

### 2. Spawn Integration Agent
```bash
# Prepare spawn command for Integration Agent
CURRENT_BRANCH=$(git branch --show-current)

cat > /tmp/integration-agent-task.md << EOF
Execute integration merges for Phase ${PHASE} Wave ${WAVE}.

🔴🔴🔴 R361 SUPREME LAW: CONFLICT RESOLUTION ONLY 🔴🔴🔴
- NO new packages or directories
- NO adapter or wrapper code
- NO "glue code" or compatibility layers
- Maximum 50 lines of changes total (excluding merges)
- Integration = conflict resolution ONLY

CRITICAL REQUIREMENTS (R260):
1. You are in INTEGRATION_DIR: ${INTEGRATION_DIR}
2. IMMEDIATELY acknowledge and set INTEGRATION_DIR variable
3. Verify you're in the correct directory
4. Read and follow WAVE-MERGE-PLAN.md EXACTLY
5. Execute merges in specified order
6. Handle conflicts as directed in plan (R361: choose versions, don't create new code)
7. Run tests after each merge
8. Document everything in work-log.md

🎬 DEMO REQUIREMENTS (R291/R330):
9. Execute effort demos after each merge (see merge plan)
10. Run integrated wave demo after all merges
11. Capture all demo outputs in demo-results/
12. Document demo status in INTEGRATION_REPORT.md
13. If ANY demo fails, mark Demo Status: FAILED (triggers ERROR_RECOVERY)

R291 GATES YOU MUST ENFORCE:
- BUILD GATE: Code must compile
- TEST GATE: All tests must pass
- DEMO GATE: Demo scripts must execute
- ARTIFACT GATE: Build outputs must exist

GRADING CRITERIA (R267):
- 50% Completeness of Integration (including demos)
- 50% Meticulous Tracking and Documentation

Your working directory: ${INTEGRATION_DIR}
Current branch: ${CURRENT_BRANCH}
Merge plan to follow: WAVE-MERGE-PLAN.md

You are spawned into state: INIT
EOF

# 🔴🔴🔴 R208 SUPREME LAW: CD BEFORE SPAWN 🔴🔴🔴
echo "🔴 R208 SUPREME LAW: CD'ing to integration directory"
cd "$INTEGRATION_DIR" || {
    echo "❌ R208 VIOLATION: Failed to CD to $INTEGRATION_DIR"
    echo "❌ GRADE: -100% (AUTOMATIC FAILURE)"
    exit 1
}

echo "📍 R208 PWD VERIFICATION: $(pwd)"  # MUST show integration directory
echo "✅ R208: Confirmed in correct directory for Integration Agent spawn"

echo "🚀 Spawning Integration Agent for merge execution..."

/spawn integration-agent INIT "$(cat /tmp/integration-agent-task.md)"

# R208: Return to orchestrator directory after spawn
cd /workspaces/project
echo "📍 R208: Returned to orchestrator directory"
```

### 3. Update State Tracking
```yaml
# Update orchestrator-state.json
integration_status:
  phase: ${PHASE}
  wave: ${WAVE}
  merge_plan_ready: true
  integration_agent_spawned: true
  integration_started_at: "$(date -u +%Y-%m-%dT%H:%M:%SZ)"
  waiting_for: "integration-agent-completion"
```

## 🔴🔴🔴 MANDATORY MEASUREMENT RULE - R304 🔴🔴🔴

**ABSOLUTE REQUIREMENTS:**
- ✅ MUST use `$CLAUDE_PROJECT_DIR/tools/line-counter.sh` for ALL line counting
- ❌ NEVER use `wc -l` or any manual counting method
- ❌ NEVER count lines any other way - this is a -100% automatic failure
- ✅ MUST specify both -b (base branch) and -c (current branch) parameters
- ✅ Base branch MUST be phase integration branch (NOT "main")

**Failure to use the line counter tool = AUTOMATIC -100% GRADE**

## Critical Requirements

### R260 Compliance - INTEGRATION_DIR Acknowledgment
The Integration Agent MUST:
1. Acknowledge the INTEGRATION_DIR immediately upon startup
2. Set INTEGRATION_DIR environment variable
3. Verify current directory matches INTEGRATION_DIR
4. Exit with error if in wrong directory

### Working Directory Setup
**CRITICAL**: The orchestrator MUST cd into INTEGRATION_DIR before spawning!
```bash
cd "$INTEGRATION_DIR"  # MANDATORY before spawn
```

## Transition Rules

### 🔴🔴🔴 CRITICAL: Update State BEFORE Stopping! 🔴🔴🔴
Per R322, you MUST update `current_state` to the next state BEFORE stopping:

```bash
# After spawning integration agent successfully:
echo "📝 Updating state file for transition..."
jq '.current_state = "MONITORING_INTEGRATION"' orchestrator-state.json > tmp.json && mv tmp.json orchestrator-state.json
jq '.previous_state = "SPAWN_INTEGRATION_AGENT"' orchestrator-state.json > tmp.json && mv tmp.json orchestrator-state.json
jq '.transition_time = \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\"' orchestrator-state.json > tmp.json && mv tmp.json orchestrator-state.json
git add orchestrator-state.json
git commit -m "state: transition from SPAWN_INTEGRATION_AGENT to MONITORING_INTEGRATION"
git push

# THEN stop per R322
echo "🛑 Stopping before MONITORING_INTEGRATION state (per R322)"
```

- Next state: MONITORING_INTEGRATION (UPDATE STATE FIRST!)
- Cannot transition if: No merge plan exists
- Must be in integration directory when spawning

## Success Criteria
- Merge plan verified to exist
- Integration Agent spawned with INTEGRATION_DIR
- Working directory set correctly
- Clear grading criteria communicated
- State tracking updated



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
