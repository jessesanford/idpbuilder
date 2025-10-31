# ORCHESTRATOR CONTINUATION ISSUES - ROOT CAUSE ANALYSIS

**Project**: idpbuilder-oci-push-planning
**Date**: 2025-10-31
**Analyzed By**: Software Factory Manager Agent
**Analysis Type**: Configuration and State File Structural Issues

---

## EXECUTIVE SUMMARY

The orchestrator's continuation command is exhibiting filesystem search behavior (`ls` commands) instead of consulting the orchestrator-state-v3.json metadata. This analysis identifies **5 critical root causes** across 3 categories: continuation command logic, state file structure, and metadata mismatches.

### Critical Findings:
1. **Continuation command hardcoded to look for wrong filename**: `./IMPLEMENTATION-PLAN.md` vs `planning/phase2/PHASE-2-IMPLEMENTATION.md`
2. **State file has incorrect metadata field name**: Stores `phase2_implementation_plan` but should be `implementation_plan_path`
3. **State file structural inconsistency**: Empty `.phases[]` array but populated `phase2_waves` field
4. **State file actual filename mismatch**: Points to `PHASE-2-IMPLEMENTATION.md` but file is `PHASE-IMPLEMENTATION-PLAN.md`
5. **Continuation command lacks state-awareness**: No logic to read phase-specific metadata from state file

---

## PART 1: ROOT CAUSE ANALYSIS

### ROOT CAUSE #1: Continuation Command Hardcoded Path (CRITICAL)

**Location**: `/home/vscode/workspaces/idpbuilder-oci-push-planning/.claude/commands/continue-orchestrating.md`
**Lines**: 138-156

**Problem**:
```bash
# Line 139 - Hardcoded path that doesn't exist
if [ ! -f "./IMPLEMENTATION-PLAN.md" ]; then
    echo "⚠️ No Master Implementation Plan found!"

    # Line 143 - Falls back to filesystem search
    if [ -f "./ARCHITECT-PROMPT-IDPBUILDER-OCI.md" ]; then
        echo "📋 Found architect prompt - will spawn architect to create plan"
        echo "Decision: NEED_ARCHITECT_FOR_PLAN"
    else
        echo "❌ No implementation plan or architect prompt found"
        echo "Please create IMPLEMENTATION-PLAN.md or provide architect requirements"
        echo "Checking for other planning files..."
        ls -la *.md 2>/dev/null | head -10 || echo "No .md files found"  # ❌ FILESYSTEM SEARCH!
        echo "Decision: MISSING_REQUIREMENTS"
    fi
```

**Why This is Wrong**:
- The continuation command looks for `./IMPLEMENTATION-PLAN.md` in project root
- This is SF 2.0 legacy behavior (pre-multi-phase support)
- SF 3.0 stores plans in `planning/phase{N}/` subdirectories
- When file not found, falls back to `ls -la *.md` filesystem search

**Expected Behavior**:
```bash
# Should read from state file first
CURRENT_PHASE=$(jq -r '.current_phase' orchestrator-state-v3.json)
IMPL_PLAN=$(jq -r ".phase${CURRENT_PHASE}_implementation_plan" orchestrator-state-v3.json)

if [ "$IMPL_PLAN" != "null" ] && [ -f "$IMPL_PLAN" ]; then
    echo "✅ Found implementation plan: $IMPL_PLAN"
else
    # THEN fall back to discovery
    echo "⚠️ No implementation plan in state file, searching..."
fi
```

---

### ROOT CAUSE #2: State File Metadata Field Name Inconsistency (CRITICAL)

**Location**: `orchestrator-state-v3.json`
**Lines**: Check with `jq 'keys' orchestrator-state-v3.json`

**Problem**:
The state file stores phase-specific metadata as:
- `phase2_implementation_plan`
- `phase2_waves`
- Presumably `phase1_implementation_plan`, `phase3_implementation_plan`, etc.

This violates the **Generic Metadata Access Pattern** principle:

```json
// ❌ WRONG - Phase-specific field names
{
  "current_phase": 2,
  "phase2_implementation_plan": "planning/phase2/PHASE-2-IMPLEMENTATION.md",
  "phase2_waves": [...]
}

// ✅ CORRECT - Generic access pattern
{
  "current_phase": 2,
  "phases": [
    {
      "phase_number": 1,
      "implementation_plan_path": "planning/phase1/...",
      "waves": [...]
    },
    {
      "phase_number": 2,
      "implementation_plan_path": "planning/phase2/PHASE-IMPLEMENTATION-PLAN.md",
      "waves": [...]
    }
  ]
}
```

**Why This Matters**:
- Orchestrator must know phase 2's plan path
- With phase-specific fields, must construct variable name: `phase${current_phase}_implementation_plan`
- This is fragile, error-prone, and doesn't scale
- Generic array allows: `phases[current_phase - 1].implementation_plan_path`

**Current Workaround**:
```bash
# What the code currently must do (brittle)
PHASE_NUM=$(jq -r '.current_phase' state.json)
PLAN_PATH=$(jq -r ".phase${PHASE_NUM}_implementation_plan" state.json)

# What it should be able to do (robust)
PLAN_PATH=$(jq -r '.phases[] | select(.phase_number == $PHASE_NUM) | .implementation_plan_path' state.json)
```

---

### ROOT CAUSE #3: State File Structural Inconsistency (ARCHITECTURAL)

**Location**: `orchestrator-state-v3.json`

**Problem - The `.phases[]` Array is Empty**:
```bash
$ jq '.phases | length' orchestrator-state-v3.json
0  # ❌ EMPTY!

$ jq '.phase2_waves[0] | keys' orchestrator-state-v3.json
[
  "depends_on",
  "wave_description",
  "wave_name",
  "wave_number"
]  # ✅ But phase2_waves exists!
```

**The Inconsistency**:
- `.phases` array is empty (length 0)
- But `.phase2_waves` exists as a separate top-level field
- And `.phase2_implementation_plan` exists as a separate top-level field
- This creates two separate data structures for the same conceptual information

**Why This is Problematic**:
1. **Schema Confusion**: Are phases in `.phases[]` or as `phaseN_*` fields?
2. **Incomplete Migration**: Suggests upgrade from SF 2.0 to SF 3.0 was partial
3. **Discovery Complexity**: Code must check BOTH locations to find phase data
4. **Validation Challenges**: Pre-commit hooks can't validate schema properly

**Expected Structure** (SF 3.0):
```json
{
  "schema_version": "3.0",
  "current_phase": 2,
  "current_wave": 1,
  "phases": [
    {
      "phase_number": 1,
      "phase_name": "Foundation",
      "status": "complete",
      "implementation_plan_path": "planning/phase1/PHASE-IMPLEMENTATION-PLAN.md",
      "architecture_plan_path": "planning/phase1/PHASE-ARCHITECTURE-PLAN.md",
      "waves": [...]
    },
    {
      "phase_number": 2,
      "phase_name": "Core Features",
      "status": "in_progress",
      "implementation_plan_path": "planning/phase2/PHASE-IMPLEMENTATION-PLAN.md",
      "architecture_plan_path": "planning/phase2/PHASE-ARCHITECTURE-PLAN.md",
      "waves": [
        {
          "wave_number": 1,
          "wave_name": "Initial Setup",
          "wave_description": "...",
          "depends_on": []
        }
      ]
    }
  ]
}
```

**Current Structure** (Broken):
```json
{
  "current_phase": 2,
  "current_wave": 1,
  "phases": [],  // ❌ EMPTY!
  "phase2_implementation_plan": "planning/phase2/PHASE-2-IMPLEMENTATION.md",
  "phase2_waves": [
    {
      "wave_number": 1,
      "wave_name": "...",
      "wave_description": "...",
      "depends_on": []
    }
  ]
}
```

---

### ROOT CAUSE #4: Filename Mismatch in State Metadata (DATA INTEGRITY)

**Problem**:
```bash
# State file says:
$ jq -r '.phase2_implementation_plan' orchestrator-state-v3.json
planning/phase2/PHASE-2-IMPLEMENTATION.md

# But actual file is:
$ ls -la planning/phase2/PHASE-IMPLEMENTATION-PLAN.md
-rw-rw-r-- 1 vscode vscode 15337 Oct 31 07:10 PHASE-IMPLEMENTATION-PLAN.md

# The file doesn't exist at the stated path:
$ ls -la planning/phase2/PHASE-2-IMPLEMENTATION.md
ls: cannot access 'planning/phase2/PHASE-2-IMPLEMENTATION.md': No such file or directory
```

**Why This Happened**:
- State file was created/updated with incorrect filename
- Likely during initialization or upgrade process
- File was either:
  - Renamed after state file update
  - Created with different name than expected
  - Never validated after state update

**Impact**:
- Even if continuation command reads state file correctly, it will fail
- `[ -f "$IMPL_PLAN" ]` check will return false
- Triggers fallback to filesystem search
- User sees "file not found" errors

**Required Fix**:
```bash
# Update state file with correct filename
jq '.phase2_implementation_plan = "planning/phase2/PHASE-IMPLEMENTATION-PLAN.md"' \
   orchestrator-state-v3.json > temp.json && mv temp.json orchestrator-state-v3.json
```

---

### ROOT CAUSE #5: Continuation Command Not State-Aware (DESIGN FLAW)

**Location**: `.claude/commands/continue-orchestrating.md` lines 267-417

**Problem - State Machine Navigation Section**:
```bash
# Line 350-361: WAVE_START case
case "$CURRENT_STATE" in
    "INIT")
        echo "Starting new orchestration"
        ACTION: Read master plan
        ACTION: Initialize Phase 1
        NEXT_STATE: "WAVE_START"
        ;;

    "WAVE_START")
        echo "Starting new wave"
        ACTION: Set up wave infrastructure if needed
        ACTION: Prepare for effort planning
        NEXT_STATE: "CREATE_NEXT_INFRASTRUCTURE"
        ;;
```

**The Issue**:
- Case statement shows what to do in each state
- But provides no guidance on WHERE to find metadata
- Doesn't instruct: "In WAVE_START, read wave plan from `.phases[current_phase].waves[current_wave]`"
- Orchestrator must infer this from general knowledge

**Expected Pattern**:
```bash
"WAVE_START")
    echo "Starting new wave"

    # READ STATE METADATA FIRST
    CURRENT_PHASE=$(jq -r '.current_phase' orchestrator-state-v3.json)
    CURRENT_WAVE=$(jq -r '.current_wave' orchestrator-state-v3.json)
    IMPL_PLAN=$(jq -r ".phases[] | select(.phase_number == $CURRENT_PHASE) | .implementation_plan_path" orchestrator-state-v3.json)

    # VALIDATE METADATA
    if [ ! -f "$IMPL_PLAN" ]; then
        echo "❌ Implementation plan not found: $IMPL_PLAN"
        echo "Decision: ERROR_STATE_FILE_CORRUPTION"
        exit 1
    fi

    # THEN proceed with state work
    ACTION: Read implementation plan from $IMPL_PLAN
    ACTION: Set up wave infrastructure
    NEXT_STATE: "CREATE_NEXT_INFRASTRUCTURE"
    ;;
```

**Why This is Critical**:
- The continuation command is the orchestrator's "startup script"
- It sets the pattern for how to access metadata
- Without explicit guidance, orchestrator falls back to intuition
- Intuition leads to filesystem searches instead of state file queries

---

## PART 2: ALL STRUCTURAL ISSUES IN STATE FILE

### Issue 2.1: Empty `.phases[]` Array
- **Severity**: ARCHITECTURAL
- **Impact**: Primary phase metadata storage is unused
- **Current State**: `.phases` has length 0
- **Expected**: Array with objects for each phase

### Issue 2.2: Phase-Specific Top-Level Fields
- **Severity**: ARCHITECTURAL
- **Impact**: Non-scalable metadata access pattern
- **Current State**: `phase2_implementation_plan`, `phase2_waves` as top-level keys
- **Expected**: All phase data in `.phases[]` array

### Issue 2.3: Missing Phase Metadata Schema
- **Severity**: CRITICAL
- **Impact**: No standardized phase object structure
- **Current State**: Phase metadata scattered across multiple fields
- **Expected**: Consistent phase object schema with:
  - `phase_number`
  - `phase_name`
  - `status`
  - `implementation_plan_path`
  - `architecture_plan_path`
  - `waves[]`

### Issue 2.4: Wave Structure Location Inconsistency
- **Severity**: HIGH
- **Impact**: Wave metadata in wrong location
- **Current State**: `phase2_waves` at top level
- **Expected**: `phases[1].waves` (nested under phase)

### Issue 2.5: Filename Mismatch in Metadata
- **Severity**: DATA INTEGRITY
- **Impact**: State file points to non-existent file
- **Current Value**: `"planning/phase2/PHASE-2-IMPLEMENTATION.md"`
- **Actual File**: `planning/phase2/PHASE-IMPLEMENTATION-PLAN.md`

### Issue 2.6: Missing Schema Version
- **Severity**: MEDIUM
- **Impact**: Can't determine if SF 2.0 or SF 3.0 structure
- **Current State**: No `schema_version` field
- **Expected**: `"schema_version": "3.0"`

### Issue 2.7: No Phase Status Tracking
- **Severity**: MEDIUM
- **Impact**: Can't determine which phases are complete
- **Current State**: Only `current_phase` number
- **Expected**: Phase objects with `status: "complete"/"in_progress"/"pending"`

---

## PART 3: CONFIGURATION AND RULE ISSUES

### Config Issue 3.1: Continuation Command Master Plan Check
- **File**: `.claude/commands/continue-orchestrating.md`
- **Lines**: 138-156
- **Severity**: CRITICAL
- **Issue**: Hardcoded `./IMPLEMENTATION-PLAN.md` path
- **Should Check**: State file metadata first, then fall back

### Config Issue 3.2: No Phase-Aware Metadata Guidance
- **File**: `.claude/commands/continue-orchestrating.md`
- **Lines**: 267-417 (state machine navigation)
- **Severity**: HIGH
- **Issue**: Case statements don't show how to read phase metadata
- **Should Include**: Example jq queries for each state's metadata needs

### Config Issue 3.3: State File Validation Not Checking Metadata
- **File**: Likely in pre-commit hooks or validation scripts
- **Severity**: HIGH
- **Issue**: Validation allows state file with wrong filename in metadata
- **Should Validate**: All file paths in state metadata actually exist

### Rule Issue 3.4: WAVE_START Rules Don't Mandate Metadata Read
- **File**: `agent-states/software-factory/orchestrator/WAVE_START/rules.md`
- **Lines**: 1-205
- **Severity**: MEDIUM
- **Issue**: Rules show WHAT to do, not WHERE to find metadata
- **Should Include**: "Step 0: Read implementation plan path from state file"

---

## PART 4: SPECIFIC ACTIONABLE FIXES

### Priority Order (MUST FIX IN THIS ORDER):

#### FIX #1: Correct the Filename in State Metadata (IMMEDIATE)
**Priority**: P0 (Blocking all other work)
**Impact**: High - Enables state file to be usable
**Effort**: 30 seconds

```bash
cd /home/vscode/workspaces/idpbuilder-oci-push-planning

# Update the filename in state file
jq '.phase2_implementation_plan = "planning/phase2/PHASE-IMPLEMENTATION-PLAN.md"' \
   orchestrator-state-v3.json > temp.json && mv temp.json orchestrator-state-v3.json

# Verify the fix
jq -r '.phase2_implementation_plan' orchestrator-state-v3.json
# Should output: planning/phase2/PHASE-IMPLEMENTATION-PLAN.md

# Verify file exists
ls -la "$(jq -r '.phase2_implementation_plan' orchestrator-state-v3.json)"

# Commit the fix
git add orchestrator-state-v3.json
git commit -m "fix: correct phase2 implementation plan filename in state"
git push
```

**Validation**:
- `jq -r '.phase2_implementation_plan' orchestrator-state-v3.json` returns correct path
- `ls -la` on that path shows file exists

---

#### FIX #2: Update Continuation Command to Read State File First (URGENT)
**Priority**: P0 (Prevents filesystem search behavior)
**Impact**: High - Fixes the root cause of reported issue
**Effort**: 10 minutes

**File**: `.claude/commands/continue-orchestrating.md`

**Change Section "Master Implementation Plan Check" (lines 136-157)**:

```bash
### 5. Master Implementation Plan Check (State-Aware)
# 🔴🔴🔴 CRITICAL: Check STATE FILE first, THEN fall back to filesystem 🔴🔴🔴

# STEP 1: Read current phase from state file
CURRENT_PHASE=$(jq -r '.current_phase' orchestrator-state-v3.json 2>/dev/null || echo "1")

# STEP 2: Read phase-specific implementation plan path from state
# NOTE: During SF 2.0 → 3.0 migration, check both locations
IMPL_PLAN=$(jq -r ".phase${CURRENT_PHASE}_implementation_plan // .phases[] | select(.phase_number == $CURRENT_PHASE) | .implementation_plan_path // empty" orchestrator-state-v3.json 2>/dev/null)

# STEP 3: Validate state file metadata
if [ -n "$IMPL_PLAN" ] && [ "$IMPL_PLAN" != "null" ]; then
    if [ -f "$IMPL_PLAN" ]; then
        echo "✅ Found implementation plan from state file: $IMPL_PLAN"
        echo "Decision: PROCEED_WITH_ORCHESTRATION"
    else
        echo "❌ State file points to non-existent file: $IMPL_PLAN"
        echo "This indicates state file corruption or stale metadata"
        echo "Decision: STATE_FILE_CORRUPTION"
        exit 1
    fi
else
    # STEP 4: Fall back to legacy discovery (SF 2.0 compatibility)
    echo "⚠️ No implementation plan in state file for phase $CURRENT_PHASE"

    if [ -f "./IMPLEMENTATION-PLAN.md" ]; then
        echo "✅ Found legacy master plan: ./IMPLEMENTATION-PLAN.md"
        echo "Decision: PROCEED_WITH_ORCHESTRATION"
    else
        echo "❌ No implementation plan found in state file or filesystem"
        echo "Please ensure orchestrator-state-v3.json contains implementation_plan_path"
        echo "Decision: MISSING_REQUIREMENTS"
        exit 1
    fi
fi
```

**Additional Change**: Update WAVE_START case in state machine navigation (lines 362-368):

```bash
"WAVE_START")
    echo "Starting new wave"

    # Read implementation plan path from state file
    IMPL_PLAN=$(jq -r ".phase${CURRENT_PHASE}_implementation_plan // .phases[] | select(.phase_number == $CURRENT_PHASE) | .implementation_plan_path" orchestrator-state-v3.json)

    echo "📋 Using implementation plan: $IMPL_PLAN"
    ACTION: Read plan from $IMPL_PLAN
    ACTION: Set up wave infrastructure if needed
    ACTION: Prepare for effort planning
    NEXT_STATE: "CREATE_NEXT_INFRASTRUCTURE"
    ;;
```

**Validation**:
- Restart orchestrator with `/continue-orchestrating`
- Should output: "✅ Found implementation plan from state file: planning/phase2/PHASE-IMPLEMENTATION-PLAN.md"
- Should NOT run `ls -la *.md` filesystem search

---

#### FIX #3: Restructure State File to Use `.phases[]` Array (ARCHITECTURAL)
**Priority**: P1 (Important but not blocking immediate work)
**Impact**: High - Enables scalable multi-phase support
**Effort**: 30 minutes

**Create Migration Script**: `tools/migrate-state-to-phases-array.sh`

```bash
#!/bin/bash
set -euo pipefail

STATE_FILE="${1:-orchestrator-state-v3.json}"
BACKUP_FILE="${STATE_FILE}.backup-$(date +%Y%m%d-%H%M%S)"

echo "Migrating state file to phases array structure..."

# Backup original
cp "$STATE_FILE" "$BACKUP_FILE"
echo "✅ Backed up to: $BACKUP_FILE"

# Create new structure with phases array
jq '
{
  schema_version: "3.0",
  current_state: .current_state,
  current_phase: .current_phase,
  current_wave: .current_wave,
  phases: [
    # Migrate phase 2 data (extend for other phases as needed)
    {
      phase_number: 2,
      phase_name: "Core Implementation",
      status: (if .current_phase == 2 then "in_progress" else "complete" end),
      implementation_plan_path: .phase2_implementation_plan,
      waves: .phase2_waves
    }
  ],
  efforts_in_progress: .efforts_in_progress,
  efforts_completed: .efforts_completed,
  split_tracking: .split_tracking,
  state_machine: .state_machine,
  transition_time: .transition_time
}
' "$STATE_FILE" > "${STATE_FILE}.migrated"

# Validate migrated structure
if jq -e '.phases | length > 0' "${STATE_FILE}.migrated" >/dev/null; then
    echo "✅ Migration successful - phases array populated"
    mv "${STATE_FILE}.migrated" "$STATE_FILE"
    echo "✅ State file updated"
else
    echo "❌ Migration failed - phases array still empty"
    rm "${STATE_FILE}.migrated"
    exit 1
fi

# Show new structure
echo ""
echo "New structure:"
jq '.phases[] | {phase_number, status, implementation_plan_path, wave_count: (.waves | length)}' "$STATE_FILE"
```

**Validation**:
```bash
# Run migration
bash tools/migrate-state-to-phases-array.sh orchestrator-state-v3.json

# Verify phases array is populated
jq '.phases | length' orchestrator-state-v3.json  # Should be > 0

# Verify phase data is accessible
jq '.phases[] | select(.phase_number == 2) | .implementation_plan_path' orchestrator-state-v3.json
# Should output: planning/phase2/PHASE-IMPLEMENTATION-PLAN.md

# Commit migration
git add orchestrator-state-v3.json
git commit -m "refactor: migrate state file to phases array structure (SF 3.0)"
git push
```

---

#### FIX #4: Add State File Metadata Validation (QUALITY)
**Priority**: P2 (Prevents future issues)
**Impact**: Medium - Catches metadata problems early
**Effort**: 20 minutes

**Create Validation Script**: `tools/validate-state-metadata.sh`

```bash
#!/bin/bash
set -euo pipefail

STATE_FILE="${1:-orchestrator-state-v3.json}"

echo "Validating state file metadata..."

ERRORS=0

# Check 1: Validate all file paths in metadata exist
echo "Checking file paths..."

# Check implementation plan paths (both old and new structure)
for PLAN_PATH in $(jq -r '.phase*_implementation_plan // .phases[].implementation_plan_path // empty' "$STATE_FILE" 2>/dev/null); do
    if [ -f "$PLAN_PATH" ]; then
        echo "  ✅ $PLAN_PATH"
    else
        echo "  ❌ MISSING: $PLAN_PATH"
        ((ERRORS++))
    fi
done

# Check 2: Validate phases array consistency
PHASES_COUNT=$(jq '.phases | length' "$STATE_FILE")
if [ "$PHASES_COUNT" -eq 0 ]; then
    echo "⚠️  Warning: .phases array is empty (may be legacy SF 2.0 structure)"
fi

# Check 3: Validate current_phase has metadata
CURRENT_PHASE=$(jq -r '.current_phase' "$STATE_FILE")
HAS_METADATA=$(jq -e ".phases[] | select(.phase_number == $CURRENT_PHASE) // .phase${CURRENT_PHASE}_implementation_plan // empty" "$STATE_FILE" >/dev/null 2>&1 && echo "yes" || echo "no")

if [ "$HAS_METADATA" = "no" ]; then
    echo "❌ Current phase $CURRENT_PHASE has no metadata in state file"
    ((ERRORS++))
fi

# Report results
echo ""
if [ $ERRORS -eq 0 ]; then
    echo "✅ State file metadata validation passed"
    exit 0
else
    echo "❌ State file validation failed with $ERRORS errors"
    exit 1
fi
```

**Integration**: Add to pre-commit hook

```bash
# In .git/hooks/pre-commit or .pre-commit-config.yaml
bash tools/validate-state-metadata.sh orchestrator-state-v3.json || {
    echo "❌ State file metadata validation failed"
    echo "Run: bash tools/validate-state-metadata.sh orchestrator-state-v3.json"
    exit 1
}
```

---

#### FIX #5: Update WAVE_START State Rules (DOCUMENTATION)
**Priority**: P2 (Clarifies expected behavior)
**Impact**: Medium - Helps orchestrator understand what to do
**Effort**: 10 minutes

**File**: `agent-states/software-factory/orchestrator/WAVE_START/rules.md`

**Add new section after line 36**:

```markdown
---

### ✅ Step 0: Read State Metadata (MANDATORY - DO THIS FIRST!)

**CRITICAL: Before ANY wave work, read metadata from state file**

```bash
# Read current phase and wave
CURRENT_PHASE=$(jq -r '.current_phase' orchestrator-state-v3.json)
CURRENT_WAVE=$(jq -r '.current_wave' orchestrator-state-v3.json)

echo "📊 Starting Phase $CURRENT_PHASE, Wave $CURRENT_WAVE"

# Read implementation plan path (supports both SF 2.0 and SF 3.0 structure)
IMPL_PLAN=$(jq -r "
  .phases[] | select(.phase_number == $CURRENT_PHASE) | .implementation_plan_path
  // .phase${CURRENT_PHASE}_implementation_plan
  // empty
" orchestrator-state-v3.json)

# Validate plan path exists
if [ -z "$IMPL_PLAN" ] || [ "$IMPL_PLAN" = "null" ]; then
    echo "❌ CRITICAL: No implementation plan found in state file for phase $CURRENT_PHASE"
    echo "State file may be corrupted or incomplete"
    echo "CONTINUE-SOFTWARE-FACTORY=FALSE REASON=MISSING_METADATA"
    exit 1
fi

if [ ! -f "$IMPL_PLAN" ]; then
    echo "❌ CRITICAL: Implementation plan does not exist: $IMPL_PLAN"
    echo "State file metadata points to non-existent file"
    echo "CONTINUE-SOFTWARE-FACTORY=FALSE REASON=STATE_FILE_CORRUPTION"
    exit 1
fi

echo "✅ Implementation plan: $IMPL_PLAN"

# Read wave metadata from state file
WAVE_METADATA=$(jq -r "
  .phases[] | select(.phase_number == $CURRENT_PHASE) | .waves[] | select(.wave_number == $CURRENT_WAVE)
  // .phase${CURRENT_PHASE}_waves[] | select(.wave_number == $CURRENT_WAVE)
  // empty
" orchestrator-state-v3.json)

if [ -z "$WAVE_METADATA" ] || [ "$WAVE_METADATA" = "null" ]; then
    echo "⚠️  Warning: No wave metadata found in state file"
    echo "Will read wave details from implementation plan"
fi

echo "✅ State metadata loaded successfully"
```

**Why This Matters**:
- Provides explicit, step-by-step instructions
- Shows EXACT jq queries to use
- Handles both SF 2.0 and SF 3.0 state structures
- Validates data before proceeding

---

## PART 5: VERIFICATION CHECKLIST

After implementing fixes, verify the following:

### Verification #1: State File Integrity
```bash
cd /home/vscode/workspaces/idpbuilder-oci-push-planning

# Check filename is correct
jq -r '.phase2_implementation_plan' orchestrator-state-v3.json
# Expected: planning/phase2/PHASE-IMPLEMENTATION-PLAN.md

# Verify file exists at that path
ls -la "$(jq -r '.phase2_implementation_plan' orchestrator-state-v3.json)"
# Expected: File found with recent timestamp

# Check phases array (after migration)
jq '.phases | length' orchestrator-state-v3.json
# Expected: > 0 (at least 1 phase)
```

### Verification #2: Continuation Command Behavior
```bash
# Test continuation command directly
bash -c 'source .claude/commands/continue-orchestrating.md (lines 136-157)'

# Should output:
# ✅ Found implementation plan from state file: planning/phase2/PHASE-IMPLEMENTATION-PLAN.md
# Decision: PROCEED_WITH_ORCHESTRATION

# Should NOT output:
# Checking for other planning files...
# ls -la *.md
```

### Verification #3: Orchestrator Startup
```bash
# Run orchestrator continuation
/continue-orchestrating

# Monitor first 20 lines of output - should see:
# ✅ Implementation plan: planning/phase2/PHASE-IMPLEMENTATION-PLAN.md
# ✅ State metadata loaded successfully

# Should NOT see:
# ⚠️ No Master Implementation Plan found!
# Checking for other planning files...
```

### Verification #4: Metadata Validation
```bash
# Run validation script
bash tools/validate-state-metadata.sh orchestrator-state-v3.json

# Expected output:
# ✅ planning/phase2/PHASE-IMPLEMENTATION-PLAN.md
# ✅ State file metadata validation passed

# Should NOT output:
# ❌ MISSING: ...
# ❌ Current phase X has no metadata
```

---

## PART 6: PREVENTION STRATEGIES

### Strategy #1: State File Schema Enforcement
- Add JSON schema validation to pre-commit hooks
- Enforce that all file paths in state metadata exist
- Validate `.phases[]` array is populated (SF 3.0)

### Strategy #2: Metadata-First Design Pattern
- All continuation commands must check state file FIRST
- Filesystem searches only as last resort
- Document this pattern in all state rules

### Strategy #3: Migration Path Documentation
- Create clear SF 2.0 → SF 3.0 migration guide
- Provide detection script to identify structure version
- Support both structures during transition period

### Strategy #4: Automated Metadata Testing
- Add test that spawns orchestrator and verifies no `ls` calls
- Test that state file queries succeed before filesystem operations
- Validate all metadata paths resolve to real files

---

## PART 7: LONG-TERM RECOMMENDATIONS

### Recommendation #1: Complete SF 3.0 Migration
- Migrate ALL state files to use `.phases[]` array structure
- Deprecate `phaseN_*` top-level fields
- Update all documentation to reflect new structure

### Recommendation #2: State File Versioning
- Add `schema_version` field to all state files
- Version: "2.0" vs "3.0"
- Code can adapt behavior based on version

### Recommendation #3: Metadata Accessor Functions
Create helper functions in continuation command:

```bash
get_implementation_plan() {
    local phase="$1"
    local state_file="$2"

    # Try SF 3.0 structure first
    local plan=$(jq -r ".phases[] | select(.phase_number == $phase) | .implementation_plan_path" "$state_file" 2>/dev/null)

    # Fall back to SF 2.0 structure
    if [ -z "$plan" ] || [ "$plan" = "null" ]; then
        plan=$(jq -r ".phase${phase}_implementation_plan" "$state_file" 2>/dev/null)
    fi

    echo "$plan"
}

# Usage:
IMPL_PLAN=$(get_implementation_plan "$CURRENT_PHASE" orchestrator-state-v3.json)
```

### Recommendation #4: State File Documentation
- Create `docs/STATE-FILE-SCHEMA.md` documenting structure
- Include examples of both SF 2.0 and SF 3.0 formats
- Show all required fields and their purposes

---

## SUMMARY OF FIXES BY PRIORITY

| Priority | Fix | Impact | Effort | Blocks Other Fixes? |
|----------|-----|--------|--------|---------------------|
| **P0** | #1: Correct filename in state | High | 30 sec | Yes |
| **P0** | #2: Update continuation command | High | 10 min | No |
| **P1** | #3: Migrate to phases array | High | 30 min | No |
| **P2** | #4: Add metadata validation | Medium | 20 min | No |
| **P2** | #5: Update WAVE_START rules | Medium | 10 min | No |

**Recommended Order of Execution**:
1. Fix #1 (blocks everything)
2. Fix #2 (fixes user-visible issue)
3. Fix #3 (architectural improvement)
4. Fix #4 (prevents regression)
5. Fix #5 (improves documentation)

**Total Time**: ~70 minutes for all fixes

---

## CONCLUSION

The orchestrator's filesystem search behavior is caused by **5 interconnected issues**:

1. ✅ **Immediate Fix**: State file points to wrong filename → Fix in 30 seconds
2. ✅ **User-Visible Fix**: Continuation command ignores state file → Fix in 10 minutes
3. ✅ **Architectural Fix**: Empty `.phases[]` array → Migrate in 30 minutes
4. ✅ **Quality Fix**: No metadata validation → Add in 20 minutes
5. ✅ **Documentation Fix**: State rules don't mandate metadata read → Update in 10 minutes

**All issues are fixable with clear, actionable steps provided above.**

The root cause is a **partial migration from SF 2.0 to SF 3.0** where:
- State file structure was updated (added phase-specific fields)
- But `.phases[]` array was never populated
- Continuation command still looks for legacy `./IMPLEMENTATION-PLAN.md`
- No validation ensures metadata accuracy

**Implementing Fix #1 and Fix #2 will resolve the immediate user issue within 11 minutes.**

---

**END OF ANALYSIS**
