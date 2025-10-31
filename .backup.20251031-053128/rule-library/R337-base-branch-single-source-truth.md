# 🔴🔴🔴 SUPREME RULE R337: Orchestrator State as Single Source of Base Branch Truth

## Criticality: SUPREME LAW
**Violation = -100% AUTOMATIC FAILURE**

## Description
The orchestrator-state-v3.json is the SOLE ARBITER of all base branch decisions. NO agent, planner, or system component may determine, guess, or calculate base branches independently. Every base branch decision MUST be recorded in and read from the state file. This ensures deterministic, traceable, and cascadable base branch management.

## 🔴🔴🔴 THE ABSOLUTE BASE BRANCH LAW 🔴🔴🔴

**NO ONE DECIDES BASE BRANCHES EXCEPT THE ORCHESTRATOR-STATE.JSON!**

### 🚨🚨🚨 MANDATORY TRACKING STRUCTURE 🚨🚨🚨

Every effort entry in orchestrator-state-v3.json MUST contain:

```json
{
  "name": "effort-name",
  "phase": 1,
  "wave": 2,
  "branch": "project/phase1/wave2/effort-name",
  "base_branch_tracking": {
    "planned_base": "phase1-wave1-integration",    // What planner specified
    "actual_base": "phase1-wave1-integration",      // What was actually used
    "base_commit": "abc123def456",                  // Exact commit SHA
    "branched_at": "2025-01-20T10:00:00Z",         // When branched
    "dependent_efforts": ["effort-2", "effort-3"],  // Who depends on this
    "depends_on": [],                               // What this depends on
    "requires_rebase": false,                       // Cascade rebase flag
    "rebase_reason": null,                          // Why rebase needed
    "last_rebase": null,                            // Last rebase timestamp
    "integration_eligible": true                    // Can be integrated
  }
}
```

## 🔴 MANDATORY WORKFLOWS 🔴

### 1. PLANNER WORKFLOW (Code Reviewer Creating Plans)
```markdown
## Base Branch Specification (MANDATORY)
**Base Branch**: phase1-wave1-integration
**Reason**: Wave 2 efforts must build on Wave 1 integrated work per R308
**Dependent Efforts**: [List any efforts that will depend on this]
**Cascade Group**: Wave 2 efforts
```

**PLANNERS MUST:**
- Explicitly specify base branch in EVERY effort plan
- Document WHY that base was chosen
- List dependent efforts that will be affected
- Never leave base branch implicit or undefined

### 2. ORCHESTRATOR WORKFLOW (Recording Decisions)
```bash
# When creating effort infrastructure
record_base_branch() {
    local EFFORT=$1
    local PLANNED_BASE=$2
    local ACTUAL_BASE=$3
    local BASE_COMMIT=$4
    
    # Record in state file (MANDATORY)
    jq --arg effort "$EFFORT" \
       --arg planned "$PLANNED_BASE" \
       --arg actual "$ACTUAL_BASE" \
       --arg commit "$BASE_COMMIT" \
       --arg timestamp "$(date -Iseconds)" \
       '.efforts_in_progress += [{
          "name": $effort,
          "base_branch_tracking": {
            "planned_base": $planned,
            "actual_base": $actual,
            "base_commit": $commit,
            "branched_at": $timestamp,
            "dependent_efforts": [],
            "depends_on": [],
            "requires_rebase": false,
            "integration_eligible": true
          }
       }]' orchestrator-state-v3.json > tmp.json && mv tmp.json orchestrator-state-v3.json
    
    echo "✅ Base branch recorded: $EFFORT from $ACTUAL_BASE ($BASE_COMMIT)"
}
```

### 3. ENGINEER WORKFLOW (Reading Base from State)
```bash
# Engineers MUST read base from state file
get_effort_base() {
    local EFFORT=$1
    
    # NEVER calculate, ALWAYS read from state
    BASE=$(jq -r --arg e "$EFFORT" '
        (.efforts_in_progress[] | select(.name == $e) | .base_branch_tracking.actual_base) //
        (.efforts_completed[] | select(.name == $e) | .base_branch_tracking.actual_base) //
        (.split_tracking[$e].splits[0].base_branch)
    ' orchestrator-state-v3.json)
    
    if [ -z "$BASE" ] || [ "$BASE" = "null" ]; then
        echo "❌ FATAL: No base branch in state for $EFFORT"
        echo "State file MUST contain base branch tracking!"
        exit 1
    fi
    
    echo "$BASE"
}

# Usage
BASE=$(get_effort_base "my-effort")
git checkout -b new-branch "$BASE"
```

## 🔴 CASCADE REBASE TRACKING 🔴

When a base branch changes (e.g., integration recreated per R327):

### Automatic Cascade Detection
```bash
mark_cascade_rebase() {
    local CHANGED_BRANCH=$1
    
    # Find all efforts depending on this branch
    jq --arg branch "$CHANGED_BRANCH" '
        .efforts_in_progress |= map(
            if .base_branch_tracking.actual_base == $branch then
                .base_branch_tracking.requires_rebase = true |
                .base_branch_tracking.rebase_reason = "Base branch \($branch) was recreated"
            else . end
        ) |
        .efforts_completed |= map(
            if .base_branch_tracking.actual_base == $branch then
                .base_branch_tracking.requires_rebase = true |
                .base_branch_tracking.rebase_reason = "Base branch \($branch) was recreated"
            else . end
        )
    ' orchestrator-state-v3.json > tmp.json && mv tmp.json orchestrator-state-v3.json
    
    # Mark dependent integrations as stale
    mark_integrations_stale "$CHANGED_BRANCH"
}
```

### Rebase Enforcement
```bash
# Before ANY work on effort
check_rebase_required() {
    local EFFORT=$1
    
    NEEDS_REBASE=$(jq -r --arg e "$EFFORT" '
        (.efforts_in_progress[] | select(.name == $e) | .base_branch_tracking.requires_rebase) //
        false
    ' orchestrator-state-v3.json)
    
    if [ "$NEEDS_REBASE" = "true" ]; then
        echo "🚨 BLOCKING: Effort $EFFORT requires rebase!"
        echo "Base branch has changed - must rebase before continuing"
        exit 1
    fi
}
```

## 🔴 VALIDATION REQUIREMENTS 🔴

### 1. Planner Validation
```bash
# Validate effort plan has base branch
validate_effort_plan() {
    local PLAN_FILE=$1
    
    if ! grep -q "Base Branch:" "$PLAN_FILE"; then
        echo "❌ FATAL: Effort plan missing base branch specification!"
        exit 1
    fi
    
    if ! grep -q "Reason:" "$PLAN_FILE"; then
        echo "❌ FATAL: Effort plan missing base branch reason!"
        exit 1
    fi
}
```

### 2. State File Validation
```bash
# Validate all efforts have base tracking
validate_state_bases() {
    # Check every effort has base tracking
    MISSING=$(jq -r '
        [.efforts_in_progress[], .efforts_completed[]] |
        map(select(.base_branch_tracking == null or 
                   .base_branch_tracking.actual_base == null)) |
        map(.name) | .[]
    ' orchestrator-state-v3.json)
    
    if [ -n "$MISSING" ]; then
        echo "❌ FATAL: Efforts missing base branch tracking:"
        echo "$MISSING"
        exit 1
    fi
}
```

### 3. Consistency Validation
```bash
# Validate planned matches actual (unless documented)
validate_base_consistency() {
    INCONSISTENT=$(jq -r '
        [.efforts_in_progress[], .efforts_completed[]] |
        map(select(.base_branch_tracking.planned_base != .base_branch_tracking.actual_base and
                   .base_branch_tracking.rebase_reason == null)) |
        map(.name) | .[]
    ' orchestrator-state-v3.json)
    
    if [ -n "$INCONSISTENT" ]; then
        echo "⚠️ WARNING: Base branch inconsistencies detected:"
        echo "$INCONSISTENT"
        echo "Document rebase_reason if intentional!"
    fi
}
```

## 🔴 INTEGRATE_WAVE_EFFORTS WITH OTHER RULES 🔴

### R308 (Incremental Branching)
- R308 defines the STRATEGY (what base to use when)
- R337 enforces the STORAGE (state file is truth)
- R308 logic MUST write to state file, not return values

### R327 (Cascade Re-integration)
- R327 triggers cascade when fixes applied
- R337 tracks the cascade requirements in state
- All affected efforts marked with requires_rebase

### R196 (Base Branch Selection)
- R196 orchestrator creates infrastructure
- R337 ensures base is read from state, not calculated
- R196 MUST update state file with base decisions

## 🚨 COMMON VIOLATIONS (AUTOMATIC FAILURE) 🚨

### ❌ VIOLATION 1: Calculating Base Branch
```bash
# WRONG - Calculating base
if [ $WAVE -eq 1 ]; then
    BASE="main"
else
    BASE="phase${PHASE}-wave$((WAVE-1))-integration"
fi
```

### ✅ CORRECT: Reading from State
```bash
# RIGHT - Reading from state
BASE=$(jq -r '.current_effort_base' orchestrator-state-v3.json)
if [ -z "$BASE" ]; then
    echo "❌ FATAL: No base branch in state!"
    exit 1
fi
```

### ❌ VIOLATION 2: Implicit Base in Plans
```markdown
## Effort Plan
Create authentication module...
[No base branch specified]
```

### ✅ CORRECT: Explicit Base
```markdown
## Effort Plan
**Base Branch**: phase1-wave1-integration
**Reason**: Wave 2 builds on Wave 1 per R308
Create authentication module...
```

### ❌ VIOLATION 3: Ignoring Rebase Requirements
```bash
# WRONG - Working without checking rebase
cd efforts/phase2/wave1/my-effort
git pull
# Start working...
```

### ✅ CORRECT: Check Rebase First
```bash
# RIGHT - Check rebase requirement
check_rebase_required "my-effort"
cd efforts/phase2/wave1/my-effort
```

## 🔴 GRADING IMPACT 🔴

- **Missing base in state**: -100% (SUPREME LAW violation)
- **Calculating instead of reading**: -100% (Determinism failure)  
- **Plan without base specification**: -50% (Planning failure)
- **Ignoring requires_rebase**: -50% (Cascade violation)
- **Inconsistent planned/actual without reason**: -25% (Documentation failure)

## 🔴 THE SINGLE SOURCE PRINCIPLE 🔴

**REMEMBER: The orchestrator-state-v3.json is the ONLY place where base branch decisions live!**

- Planners PROPOSE bases (in plans)
- Orchestrator RECORDS bases (in state)  
- Engineers READ bases (from state)
- Changes CASCADE through state
- Validation CHECKS state

**NEVER decide base branches anywhere else!**