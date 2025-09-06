# Orchestrator - ANALYZE_BUILD_FAILURES State Rules

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
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R287-todo-save-frequency.md`
**Criticality**: BLOCKING - Save every 15 minutes/10 messages
**Summary**: Mandatory TODO saves during analysis

### 🔴🔴🔴 R288 - State File Update and Commit Protocol (SUPREME LAW)
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R288-state-file-update-and-commit-protocol.md`
**Criticality**: SUPREME LAW - Update on every transition
**Summary**: Update orchestrator-state.yaml on all state changes

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