# Orchestrator - SPAWN_CODE_REVIEWER_PROJECT_FIX_PLANNING State Rules

# PRIMARY DIRECTIVES

You MUST read and acknowledge these rules:

1. **R006** - Orchestrator cannot write code (BLOCKING)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R006-orchestrator-never-writes-code.md`

2. **R256** - Fix Planning Protocol
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R256-fix-planning-protocol.md`

4. **R287** - TODO Persistence Requirements (SUPREME LAW)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R287-todo-persistence-comprehensive.md`

5. **R288** - State File Update Requirements (SUPREME LAW)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R288-state-file-update-requirements.md`

6. **R304** - Mandatory Line Counter Usage (SUPREME LAW)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R304-mandatory-line-counter-usage.md`

7. **R322** - Mandatory Stop After Spawn States (SUPREME LAW)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R322-mandatory-stop-after-spawn.md`

8. **R324** - State Transition Validation (SUPREME LAW)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R324-state-transition-validation.md`


## 🛑🛑🛑 R322 MANDATORY STOP BEFORE STATE TRANSITIONS 🛑🛑🛑

**SUPREME LAW - VIOLATION = -100% IMMEDIATE FAILURE**

### Part A: MANDATORY STOP AFTER SPAWNING AGENTS
**WHEN IN ANY SPAWN STATE, YOU MUST:**
1. ✅ Spawn the agent(s) as required
2. ✅ Update state file with spawn information
3. ✅ Commit and push the state file
4. ✅ STOP IMMEDIATELY - Do not continue

### Part B: MANDATORY STOP BEFORE STATE TRANSITIONS
**WHEN TRANSITIONING FROM ANY STATE:**
1. ✅ Complete all work for current state
2. ✅ Update state file with next state
3. ✅ Commit and push the state file
4. ✅ STOP and wait for user to continue

### STOP PROTOCOL:
```markdown
## 🛑 STATE TRANSITION CHECKPOINT: SPAWN_CODE_REVIEWER_PROJECT_FIX_PLANNING → WAITING_FOR_PROJECT_FIX_PLANS

### ✅ Current State Work Completed:
- Spawned Code Reviewer for project fix planning
- Provided bug documentation from integration report
- Updated state file with spawn information

### 📊 Current Status:
- Current State: SPAWN_CODE_REVIEWER_PROJECT_FIX_PLANNING
- Next State: WAITING_FOR_PROJECT_FIX_PLANS
- Code Reviewer Spawned: ✅
- State Files: Updated and committed ✅

### ⏸️ STOPPED - Awaiting User Continuation
Ready to transition to WAITING_FOR_PROJECT_FIX_PLANS. Please use /continue-orchestrating.
```

**STOP MEANS STOP - Exit and wait for /continue-orchestrating**

---

## 🔴🔴🔴 R290 ENFORCEMENT: READ THESE RULES FIRST! 🔴🔴🔴

**SUPREME LAW #3 (R290): STATE RULES MUST BE READ BEFORE STATE ACTIONS**

## 🔴🔴🔴 STOP! STATE RULE READING IS ABSOLUTELY FIRST! 🔴🔴🔴

**YOU HAVE ENTERED SPAWN_CODE_REVIEWER_PROJECT_FIX_PLANNING STATE - YOU MUST READ AND ACKNOWLEDGE ALL STATE RULES BEFORE DOING ANY STATE WORK!**

## 🔴🔴🔴 R290 VERIFICATION REQUIREMENT 🔴🔴🔴

**R290 ENFORCEMENT: CREATE VERIFICATION MARKER AFTER READING**

After reading and acknowledging all state rules, you MUST create a verification marker:

```bash
# MANDATORY: Create verification marker after reading rules
touch .state_rules_read_orchestrator_SPAWN_CODE_REVIEWER_PROJECT_FIX_PLANNING
echo "$(date +%s) - Rules read and acknowledged for SPAWN_CODE_REVIEWER_PROJECT_FIX_PLANNING" > .state_rules_read_orchestrator_SPAWN_CODE_REVIEWER_PROJECT_FIX_PLANNING
```

**FAILURE TO CREATE MARKER = AUTOMATIC -100% PENALTY**

The system will check for this marker. No marker = Immediate failure.

### ❌ DO NOT DO ANY SPAWNING WORK UNTIL RULES ARE READ:
- ❌ Read project integration report
- ❌ Extract bug information
- ❌ Create command files
- ❌ Spawn Code Reviewer agents
- ❌ Update state files
- ❌ Continue to next state
- ❌ Think about what to do in this state

### ✅ YOU MUST IMMEDIATELY:

## 🔴🔴🔴 MANDATORY STATE RULE READING AND ACKNOWLEDGMENT 🔴🔴🔴

### 🛑 RULE R322 - Mandatory Stop Before State Transitions
**Source:** $CLAUDE_PROJECT_DIR/rule-library/R322-mandatory-stop-before-state-transitions.md
**Criticality:** 🔴🔴🔴 SUPREME LAW - Violation = -100% FAILURE

After completing state work and committing state file:
1. STOP IMMEDIATELY after spawning
2. Do NOT continue to next state
3. Do NOT start new work
4. Exit and wait for user
---

**USE THESE EXACT READ COMMANDS (IN THIS ORDER):**
1. Read: $CLAUDE_PROJECT_DIR/rule-library/R234-mandatory-state-traversal-supreme-law.md
2. Read: $CLAUDE_PROJECT_DIR/rule-library/R006-orchestrator-never-writes-code.md
3. Read: $CLAUDE_PROJECT_DIR/rule-library/R290-state-rule-reading-verification-supreme-law.md
4. Read: $CLAUDE_PROJECT_DIR/rule-library/R266-upstream-bug-documentation.md
5. Read: $CLAUDE_PROJECT_DIR/rule-library/R321-immediate-backport-during-integration.md
6. Read: $CLAUDE_PROJECT_DIR/rule-library/R219-code-reviewer-dependency-aware-effort-planning.md
7. Read: $CLAUDE_PROJECT_DIR/rule-library/R206-state-machine-transition-validation.md

**WE ARE WATCHING EACH READ TOOL CALL**

### ❌ ANTI-PATTERNS THAT WILL CAUSE FAILURE:

1. **Fake Acknowledgment Without Reading**:
   ```
   ❌ WRONG: "I acknowledge R234, R006, R290..."
   (No Read tool calls detected = AUTOMATIC FAILURE)
   ```

2. **Bulk Acknowledgment**:
   ```
   ❌ WRONG: "I acknowledge all SPAWN_CODE_REVIEWER_PROJECT_FIX_PLANNING rules"
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
   ❌ WRONG: "I know R266 requires bug documentation..."
   (Must READ from file, not recall from memory)
   ```

### ✅ CORRECT PATTERN FOR SPAWN_CODE_REVIEWER_PROJECT_FIX_PLANNING:
```
1. READ: $CLAUDE_PROJECT_DIR/rule-library/R234-mandatory-state-traversal-supreme-law.md
2. "I acknowledge R234 - Mandatory State Traversal: [Brief description]"
3. READ: $CLAUDE_PROJECT_DIR/rule-library/R006-orchestrator-never-writes-code.md  
4. "I acknowledge R006 - Orchestrator Never Writes Code: [Brief description]"
### 🚨 NO WORK UNTIL ACKNOWLEDGMENT COMPLETE 🚨
**You may NOT begin ANY spawning work until:**
### 🚨 FAILURE TO READ STATE RULES FIRST = IMMEDIATE EXIT 🚨
**If you do ANY spawning work before reading and acknowledging rules:**
- **STOP ALL WORK IMMEDIATELY**
- **EXIT WITH FAILURE STATUS**
- **YOU HAVE VIOLATED STATE COMPLIANCE**

---

## 🔴🔴🔴 SUPREME DIRECTIVE: CREATE PROJECT-LEVEL FIX PLANS 🔴🔴🔴

**SPAWN CODE REVIEWER FOR PROJECT INTEGRATION FIX PLANNING!**

## State Overview

In SPAWN_CODE_REVIEWER_PROJECT_FIX_PLANNING, you spawn a Code Reviewer to analyze project integration bugs documented per R266 and create comprehensive fix plans.

## State Context

**Purpose:**
Spawn Code Reviewer to create fix plans for bugs found and documented during project integration per R266. This follows the proper separation of concerns:
- Orchestrator: Coordinates and spawns agents
- Code Reviewer: Analyzes bugs and creates fix plans
- SW Engineer: Implements the fixes

## Required Actions

### 1. Read Bug Documentation
```bash
# Read the project integration report with bug documentation
PROJECT_REPORT="project-integration/PROJECT-INTEGRATION-REPORT.md"

if [ ! -f "$PROJECT_REPORT" ]; then
    echo "❌ ERROR: Project integration report not found!"
    exit 1
fi

# Extract bug section per R266
BUG_SECTION=$(sed -n '/## UPSTREAM BUGS IDENTIFIED/,/## /p' "$PROJECT_REPORT")

if [ -z "$BUG_SECTION" ]; then
    echo "❌ ERROR: No bug documentation found in report!"
    echo "Should not be in PROJECT_FIX_PLANNING state without bugs"
    exit 1
fi

echo "✅ Found bug documentation per R266"
```

### 2. Spawn Code Reviewer for Fix Planning
```bash
TIMESTAMP=$(date +%Y%m%d-%H%M%S)
INTEGRATION_PATH="project-integration"

# Create command for Code Reviewer
COMMAND_FILE="$INTEGRATION_PATH/code-reviewer-project-fix-command.md"

cat > "$COMMAND_FILE" << 'EOF'
# CODE REVIEWER PROJECT FIX PLAN CREATION

## Your State
You are in state: CREATE_PROJECT_FIX_PLAN

## Project Integration Bug Analysis

1. **Read project integration report**:
   - Location: `project-integration/PROJECT-INTEGRATION-REPORT.md`
   - Focus on: "UPSTREAM BUGS IDENTIFIED" section (R266)

2. **Analyze each documented bug**:
   - Bug number and title
   - Source branch (which phase/wave/effort)
   - Specific file and line numbers
   - Root cause analysis
   - Severity assessment

3. **Create comprehensive fix plan**:
   - Output: `PROJECT-FIX-PLAN.md`
   - For EACH bug documented per R266:
     * Identify exact source branch containing bug
     * Provide specific fix instructions
     * Include code snippets showing before/after
     * Specify verification steps
   - Group fixes by parallelization capability
   - Respect R321: ALL fixes go to source branches

4. **Fix Plan Structure**:
   ```markdown
   # PROJECT INTEGRATION FIX PLAN
   
   ## Bug Summary
   - Total Bugs Found: X
   - Critical: Y / High: Z / Medium: A / Low: B
   
   ## Fix Strategy
   
   ### Parallel Fix Group 1
   #### Bug #1: [Title from R266 documentation]
   - **Source Branch**: phase-X-wave-Y-effort-Z
   - **File**: path/to/file.ext
   - **Lines**: 234-237
   - **Current Code**:
     ```language
     [buggy code]
     ```
   - **Fixed Code**:
     ```language
     [corrected code]
     ```
   - **Verification**: [How to verify fix]
   
   ### Sequential Fix Group
   [Bugs with dependencies]
   
   ## SW Engineer Spawn Instructions
   - Engineer 1: Bugs #1, #3, #5
   - Engineer 2: Bugs #2, #4
   - Sequential: Bug #6 after parallel complete
   ```

5. **CRITICAL R321 Compliance**:
   - NEVER suggest fixes to integration branch
   - ALL fixes MUST target source branches
   - Integration will be re-run after fixes

Remember: Create the plan, do NOT execute fixes!
EOF

echo "🚀 Spawning Code Reviewer for project fix planning"
echo "@agent-code-reviewer Please execute: $COMMAND_FILE"

# Update state with spawn information
jq ".agents_spawned += [{
    \"type\": \"code-reviewer\",
    \"task\": \"project_fix_plan\",
    \"state\": \"CREATE_PROJECT_FIX_PLAN\",
    \"command_file\": \"$COMMAND_FILE\",
    \"timestamp\": \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\"
}]" -i orchestrator-state.json

# Transition to waiting state
jq ".current_state = \"WAITING_FOR_PROJECT_FIX_PLANS\"" -i orchestrator-state.json

# Record bug information for tracking
jq ".project_fixes = {
    \"bug_count\": $(echo "$BUG_SECTION" | grep -c "^### Bug #"),
    \"report_location\": \"$PROJECT_REPORT\",
    \"fix_plan_requested\": true,
    \"timestamp\": \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\"
}" -i orchestrator-state.json

git add orchestrator-state.json "$COMMAND_FILE"
git commit -m "spawn: Code Reviewer for project fix planning per R266"
git push
```

## Valid Transitions

1. **ALWAYS**: `SPAWN_CODE_REVIEWER_PROJECT_FIX_PLANNING` → `WAITING_FOR_PROJECT_FIX_PLANS`
   - Transition after spawning Code Reviewer
   - MUST STOP after spawn per R322 Part A

## Spawn Requirements

1. Verify bug documentation exists per R266
2. Create clear command file with fix instructions
3. Reference the PROJECT-INTEGRATION-REPORT.md
4. Specify CREATE_PROJECT_FIX_PLAN state
5. Emphasize R321 compliance (source branch fixes)
6. Record spawn in state file
7. Transition to waiting state
8. STOP immediately after spawn

## Grading Criteria

- ✅ **+25%**: Spawn Code Reviewer correctly
- ✅ **+25%**: Reference R266 bug documentation
- ✅ **+25%**: Emphasize R321 compliance
- ✅ **+25%**: Stop after spawn per R322

## Common Violations

- ❌ **-100%**: Not spawning Code Reviewer (orchestrator creating plans)
- ❌ **-100%**: Continuing after spawn (R322 violation)
- ❌ **-50%**: Missing R266 bug documentation reference
- ❌ **-50%**: Not emphasizing source branch fixes (R321)
- ❌ **-30%**: Not recording spawn in state file

## Related Rules

- R266: Upstream Bug Documentation Protocol
- R321: Immediate Backport During Integration
- R219: Code Reviewer Dependency-Aware Planning
- R206: State Machine Transition Validation
- R322: Mandatory Stop Before State Transitions

## 🔴🔴🔴 MANDATORY MEASUREMENT RULE - R304 🔴🔴🔴

**ABSOLUTE REQUIREMENTS:**
- ✅ MUST use `$CLAUDE_PROJECT_DIR/tools/line-counter.sh` for ALL line counting
- ❌ NEVER use `wc -l` or any manual counting method
- ❌ NEVER count lines any other way - this is a -100% automatic failure
- ✅ Base branch for project fixes is typically main or project-integration

**Failure to use the line counter tool = AUTOMATIC -100% GRADE**

## R322 VIOLATION DETECTION

If you find yourself:
- Starting work for WAITING_FOR_PROJECT_FIX_PLANS without stopping
- Transitioning without stopping after spawning
- Creating fix plans yourself instead of spawning Code Reviewer

**STOP IMMEDIATELY - You are violating R322 and proper separation of concerns!**

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
