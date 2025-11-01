# Orchestrator - ANALYZE_BUILD_FAILURES State Rules

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

## 🔴🔴🔴 R290 ENFORCEMENT: READ THESE RULES FIRST! 🔴🔴🔴

**SUPREME LAW #3 (R290): STATE RULES MUST BE READ BEFORE STATE ACTIONS**

## 🔴🔴🔴 STOP! STATE RULE READING IS ABSOLUTELY FIRST! 🔴🔴🔴

**YOU HAVE ENTERED ANALYZE_BUILD_FAILURES STATE - YOU MUST READ AND ACKNOWLEDGE ALL STATE RULES BEFORE DOING ANY STATE WORK!**

## 🔴🔴🔴 R290 VERIFICATION REQUIREMENT 🔴🔴🔴

**R290 ENFORCEMENT: CREATE VERIFICATION MARKER AFTER READING**

After reading and acknowledging all state rules, you MUST create a verification marker:

```bash
# MANDATORY: Create verification marker after reading rules
touch .state_rules_read_orchestrator_ANALYZE_BUILD_FAILURES
echo "$(date +%s) - Rules read and acknowledged for ANALYZE_BUILD_FAILURES" > .state_rules_read_orchestrator_ANALYZE_BUILD_FAILURES
```

**FAILURE TO CREATE MARKER = AUTOMATIC -100% PENALTY**

## 📋 PRIMARY DIRECTIVES FOR ANALYZE_BUILD_FAILURES STATE

### Core Mandatory Rules

### 🚨🚨🚨 R006 - ORCHESTRATOR NEVER WRITES CODE OR PERFORMS FILE OPERATIONS (BLOCKING)
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R006-orchestrator-never-writes-code.md`
**Criticality**: BLOCKING - Automatic termination, 0% grade
**Summary**: NEVER write, copy, move, or manipulate ANY code files - delegate ALL to agents

### 🚨🚨🚨 R319 - ORCHESTRATOR NEVER MEASURES CODE (BLOCKING)
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R319-orchestrator-never-measures-code.md`
**Criticality**: BLOCKING - Orchestrator MUST NOT use line-counter.sh
**Summary**: Code Reviewers measure code size, NOT orchestrators

### ⚠️⚠️⚠️ R317 - Working Directory Restrictions (WARNING)
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R317-working-directory-restrictions.md`
**Criticality**: WARNING - -25% for violations
**Summary**: MUST NOT enter agent working directories - operate from root only

### 🚨🚨🚨 R287 - TODO Save Frequency Requirements (BLOCKING)
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R287-todo-persistence-comprehensive.md`
**Criticality**: BLOCKING - Save every 15 minutes/10 messages
**Summary**: Mandatory TODO saves during analysis

### 🔴🔴🔴 R288 - State File Update and Commit Protocol (SUPREME LAW)
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R288-state-file-update-and-commit-protocol.md`
**Criticality**: SUPREME LAW - Update on every transition
**Summary**: Update orchestrator-state-v3.json on all state changes

### State-Specific Rules

### 🔴🔴🔴 R321 - Immediate Backport During Integration (SUPREME LAW)
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R321-immediate-backport-during-integration.md`
**Criticality**: SUPREME LAW - Integration fixes must be backported immediately
**Summary**: All fixes during integration require immediate backporting

## 🔴🔴🔴 CRITICAL: ANALYZE_BUILD_FAILURES IS A VERB - START ANALYZING NOW! 🔴🔴🔴

**ANALYZE_BUILD_FAILURES MEANS ACTIVELY ANALYZING BUILD ERRORS RIGHT NOW!**
- ❌ NOT "I'm in analyze build failures state"  
- ❌ NOT "Ready to analyze failures"
- ✅ YES "I'm reading BUILD-ERRORS.txt NOW"
- ✅ YES "I'm categorizing compilation errors NOW"
- ✅ YES "I'm mapping errors to efforts NOW"

## State Context
ANALYZE_BUILD_FAILURES = You ARE ACTIVELY analyzing build failure reports and categorizing errors THIS INSTANT. Focus on understanding root causes and mapping to source efforts.

## 🎯 STATE OBJECTIVES

In the ANALYZE_BUILD_FAILURES state, you are responsible for:

1. **Reading Build Reports**
   - Read BUILD-VALIDATION-REPORT.md
   - Read BUILD-ERRORS.txt
   - Parse compilation errors
   - Identify linking issues
   - Find missing dependencies

2. **Categorizing Failures**
   - Group by error type
   - Map to source efforts
   - Identify dependencies
   - Determine fix order

3. **Creating Error Analysis**
   - Document root causes
   - Map errors to files
   - Identify affected components
   - Create fix priority list

4. **Preparing for Fix Planning**
   - Document what needs fixing
   - Identify which efforts are affected
   - Prepare to spawn Code Reviewer for fix plans
   - Track for backporting (R321)

## 📝 REQUIRED ACTIONS

### Step 1: Read Build Failure Reports
```bash
cd /efforts/integration-testing

# Read build validation report
if [ -f "BUILD-VALIDATION-REPORT.md" ]; then
    echo "📋 Reading build validation report..."
    cat BUILD-VALIDATION-REPORT.md
else
    echo "❌ No build validation report found!"
    exit 1
fi

# Read detailed build errors
if [ -f "BUILD-ERRORS.txt" ]; then
    echo "📋 Analyzing build errors..."
    cat BUILD-ERRORS.txt
else
    echo "⚠️ No specific error file, checking build output..."
    grep -i "error" BUILD-OUTPUT.log | head -50
fi
```

### Step 2: Categorize Error Types
```bash
# Create error categorization
cat > BUILD-ERROR-ANALYSIS.md << 'EOF'
# Build Error Analysis
Date: $(date)
State: ANALYZE_BUILD_FAILURES

## Error Categories

### Compilation Errors
Total: [count]

#### Type Errors
- Count: [number]
- Files affected: [list]
- Example: [sample error]

#### Import/Package Errors
- Count: [number]
- Packages missing: [list]
- Example: [sample error]

#### Syntax Errors
- Count: [number]
- Files affected: [list]
- Example: [sample error]

### Linking Errors
Total: [count]

#### Undefined References
- Count: [number]
- Symbols missing: [list]
- Example: [sample error]

#### Multiple Definitions
- Count: [number]
- Symbols duplicated: [list]
- Example: [sample error]

### Dependency Issues
Total: [count]

#### Missing Dependencies
- Packages: [list]
- Required by: [components]

#### Version Conflicts
- Packages: [list with versions]
- Conflicts: [details]

## Error Distribution by Effort
| Effort | Compilation | Linking | Dependencies | Total |
|--------|-------------|---------|--------------|-------|
| [name] | [count]     | [count] | [count]      | [sum] |

## Root Cause Analysis
1. [Primary cause]
2. [Secondary cause]
3. [Contributing factors]

## Fix Priority
1. [Dependencies first]
2. [Then compilation]
3. [Finally linking]
EOF
```

### Step 3: Map Errors to Source Efforts
```bash
# Create error-to-effort mapping
cat > ERROR-TO-EFFORT-MAP.md << 'EOF'
# Error to Effort Mapping
Date: $(date)
Analyzer: orchestrator

## Compilation Error Mapping

### Error: [Specific error message]
- File: [path/to/file]
- Line: [number]
- Original Effort: [effort-name]
- Branch: [effort-branch]
- Category: [type/import/syntax]
- Fix Strategy: [description]

### Error: [Next error]
[Continue for all errors...]

## Linking Error Mapping

### Error: [Undefined reference to X]
- Component: [component name]
- Original Effort: [effort-name]
- Branch: [effort-branch]
- Missing Symbol: [symbol]
- Fix Strategy: [description]

## Dependency Mapping

### Missing: [package/dependency]
- Required By: [component]
- Original Effort: [effort-name]
- Branch: [effort-branch]
- Fix Strategy: Add to [file]

## Effort Summary
| Effort | Errors | Priority | Estimated Complexity |
|--------|--------|----------|---------------------|
| [name] | [count]| [1-5]    | [Low/Med/High]      |

## Fix Sequencing
Based on dependencies:
1. Fix [effort-1] first (provides [what])
2. Then [effort-2] (depends on effort-1)
3. Finally [effort-3] (depends on both)
EOF
```

### Step 4: Create Analysis Summary
```bash
# Create summary for next steps
cat > BUILD-FAILURE-ANALYSIS-SUMMARY.md << 'EOF'
# Build Failure Analysis Summary
Date: $(date)
State: ANALYZE_BUILD_FAILURES

## Statistics
- Total Errors: [number]
- Affected Efforts: [count]
- Critical Blockers: [count]
- Estimated Fix Time: [hours]

## Critical Findings
1. [Most critical issue]
2. [Second critical issue]
3. [Third critical issue]

## Recommended Approach
Strategy: [Targeted fixes / Complete rebuild / Phased approach]
Rationale: [Why this approach]

## Next Steps
1. Spawn Code Reviewer to create detailed fix plans
2. Code Reviewer will analyze each error in detail
3. Code Reviewer will create FIX-PLAN documents
4. Then spawn SW Engineers to implement fixes

## Backport Requirements (R321)
Integration Context: [active/inactive]
If active, ALL fixes must be immediately backported to:
[List source branches]

## Risk Assessment
- Build Recovery Likelihood: [High/Medium/Low]
- Complexity: [Simple/Moderate/Complex]
- Dependencies: [None/Some/Many]
EOF
```

### Step 5: Prepare for Code Reviewer Spawn
```bash
# Document what Code Reviewer needs to analyze
cat > CODE-REVIEWER-FIX-PLANNING-BRIEF.md << 'EOF'
# Brief for Code Reviewer Fix Planning
Date: $(date)
From: orchestrator/ANALYZE_BUILD_FAILURES

## Your Task
Create detailed fix plans for build failures.

## Input Documents
1. BUILD-ERROR-ANALYSIS.md - Categorized errors
2. ERROR-TO-EFFORT-MAP.md - Error locations
3. BUILD-ERRORS.txt - Raw error output

## Required Outputs
For each affected effort, create:
- FIX-PLAN-[effort-name].md

Each fix plan must include:
- Specific code changes needed
- Exact files and line numbers
- Before/after code snippets
- Test requirements
- Build verification steps

## Priority Order
[Based on analysis, list priority]

## Success Criteria
- All errors have fix plans
- Plans are implementable by SW Engineers
- Dependencies resolved in correct order
- Each fix is testable
EOF
```

## ⚠️ CRITICAL REQUIREMENTS

### 🔴 R006 ENFORCEMENT
**NEVER attempt to fix code yourself!**
- ❌ Don't edit any source files
- ❌ Don't write patches
- ❌ Don't create code snippets
- ✅ Only analyze and document
- ✅ Spawn Code Reviewer for fix planning
- ✅ Spawn SW Engineers for implementation

### 🔴 R321 BACKPORT TRACKING
**During integration, track ALL fixes for immediate backporting:**
- Document source branch for each error
- Note which files need backporting
- Prepare backport manifest
- Enforce immediate backport requirement

## State Transitions

From ANALYZE_BUILD_FAILURES state:
- **ANALYSIS_COMPLETE** → SPAWN_CODE_REVIEWER_BUILD_FIX_PLAN (Need detailed fix plans)
- **NO_ERRORS_FOUND** → BUILD_VALIDATION (Rerun validation)
- **CRITICAL_STRUCTURAL_ISSUES** → ERROR_RECOVERY (Cannot proceed)

## Success Criteria
- ✅ All build errors read and analyzed
- ✅ Errors categorized by type
- ✅ Errors mapped to source efforts
- ✅ Root causes identified
- ✅ Fix priority determined
- ✅ Analysis documents created and saved

## Failure Triggers
- ❌ Auto-transition without stop = R322 VIOLATION
- ❌ Attempting to fix code = R006 VIOLATION (-100%)
- ❌ Missing error analysis = -40% penalty
- ❌ Not tracking for backports = R321 VIOLATION

## R322 VIOLATION DETECTION

If you find yourself:
- Starting work for a new state without /continue-orchestrating
- Transitioning without stopping after state file commit
- Continuing after completing state work

**STOP IMMEDIATELY - You are violating R322!**


## 🔴🔴🔴 STATE COMPLETION CHECKLIST 🔴🔴🔴

**Execute these steps IN ORDER to properly complete ANALYZE_BUILD_FAILURES:**

### ✅ Step 1: Complete State-Specific Work
**Refer to "Primary Actions" and "Mandatory Validations" sections above for state-specific tasks.**

Once all state work is complete, proceed to mandatory exit protocol:

---

### ✅ Step 2: Set Proposed Next State and Transition Reason
```bash
# Based on state work results, set variables for State Manager
PROPOSED_NEXT_STATE="[DETERMINE_FROM_STATE_LOGIC]"
TRANSITION_REASON="[REASON_FOR_TRANSITION]"
echo "Proposed next state: $PROPOSED_NEXT_STATE"
echo "Transition reason: $TRANSITION_REASON"
```

---

### ✅ Step 3: Spawn State Manager for State Transition (R288 - SUPREME LAW)
```bash
# State Manager handles ALL state file updates atomically
/spawn state-manager SHUTDOWN_CONSULTATION \
  --current-state "CURRENT_STATE_NAME" \
  --proposed-next-state "$PROPOSED_NEXT_STATE" \
  --transition-reason "$TRANSITION_REASON"

# State Manager will:
# 1. Validate the transition against state machine
# 2. Update all 4 state tracking locations atomically:
#    - orchestrator-state-v3.json
#    - orchestrator-state-demo.json
#    - .cascade-state-backup.json
#    - .orchestrator-state-v3.json
# 3. Commit and push all changes
# 4. Return control

echo "✅ State Manager completed transition"
```

---

### ✅ Step 4: Save TODOs (R287 - SUPREME LAW)
```bash
# Save TODO state before exit (R287 trigger)
save_todos "STATE_COMPLETE"

# Commit TODOs within 60 seconds (R287)
cd "$CLAUDE_PROJECT_DIR"
git add todos/*.todo

if ! git commit -m "todo: orchestrator - state complete [R287]"; then
    echo "❌ ERROR: Failed to commit TODO files"
    echo "This is non-fatal but TODOs may be lost in compaction"
    echo "Proceeding with state execution..."
    # Don't exit - TODO commit failure is not fatal
fi

git push || echo "⚠️ WARNING: TODO push failed - committed locally"
echo "✅ TODOs saved and committed"
git push
echo "✅ TODOs saved and committed"
```

---

### ✅ Step 5: Output Continuation Flag (R405 - SUPREME LAW) ⚠️ MANDATORY
```bash
# Output continuation flag as LAST action (R405)
# Use TRUE for normal completion, FALSE only for catastrophic errors

echo "CONTINUE-SOFTWARE-FACTORY=TRUE REASON=STATE_COMPLETE"
```

**⚠️ THIS MUST BE THE ABSOLUTE LAST LINE OF OUTPUT BEFORE EXIT! ⚠️**

---

### ✅ Step 6: Stop Processing (R322 - SUPREME LAW)
```bash
# Stop for context preservation (R322)
echo "🛑 Stopping for context preservation - use /continue-orchestrating to resume"
exit 0
```

---

## 🚨 CHECKLIST ENFORCEMENT 🚨

**Skipping ANY step in this checklist = FAILURE:**
- Missing Step 2: No proposed next state = State Manager can't proceed
- Missing Step 3: No State Manager consultation = bypassing bookend pattern (R288 violation, -100%)
- Missing Step 4: No TODO save = work lost (R287 violation, -20% to -100%)
- Missing Step 5: No continuation flag = automation stops (R405 violation, -100%)
- Missing Step 6: No exit = R322 violation (-100%)

**ALL 6 STEPS ARE MANDATORY - NO EXCEPTIONS**
**NOTE**: State file validation and commits now handled by State Manager (SF 3.0 pattern)

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
PROPOSED_NEXT_STATE="SPAWN_SW_ENGINEERS"
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

