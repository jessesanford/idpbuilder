# 🚨🚨🚨 BLOCKING RULE R502: Mandatory Plan Validation Gates

## Rule Statement
NO phase, wave, or effort work can begin without validated planning documents. The system MUST verify plan existence BEFORE creating infrastructure, spawning agents, or starting implementation. Missing plans BLOCK ALL PROGRESS.

## Criticality: 🚨🚨🚨 BLOCKING
Plan validation failures cause catastrophic project failures where work proceeds without architectural guidance, leading to wasted effort, incorrect implementation, and integration failures.

## Mandatory Validation Points

### 1. Phase Start Validation
**BEFORE starting any phase work:**
```bash
validate_phase_plans() {
    local PHASE="$1"

    echo "🔍 VALIDATING PHASE $PHASE PLANS..."

    # MANDATORY phase-level documents
    local REQUIRED_DOCS=(
        "$CLAUDE_PROJECT_DIR/planning/phase${PHASE}/PHASE-${PHASE}-ARCHITECTURE-PLAN.md"
        "$CLAUDE_PROJECT_DIR/planning/phase${PHASE}/PHASE-${PHASE}-IMPLEMENTATION-PLAN.md"
    )

    local ALL_EXIST=true
    for doc in "${REQUIRED_DOCS[@]}"; do
        if [ ! -f "$doc" ]; then
            echo "❌ CRITICAL: Missing required document: $doc"
            ALL_EXIST=false
        else
            echo "✅ Found: $(basename "$doc")"
        fi
    done

    if [ "$ALL_EXIST" = false ]; then
        echo "🚨🚨🚨 PHASE $PHASE CANNOT START - MISSING PLANS!"
        echo "Action Required: Architect and Code Reviewer must create plans first"
        exit 1
    fi

    echo "✅ Phase $PHASE plans validated"
}
```

### 2. Wave Start Validation
**BEFORE starting any wave work:**
```bash
validate_wave_plans() {
    local PHASE="$1"
    local WAVE="$2"

    echo "🔍 VALIDATING PHASE $PHASE WAVE $WAVE PLANS..."

    # First validate phase plans exist
    validate_phase_plans "$PHASE"

    # MANDATORY wave-level documents
    local REQUIRED_DOCS=(
        "$CLAUDE_PROJECT_DIR/planning/phase${PHASE}/wave${WAVE}/WAVE-${PHASE}-${WAVE}-ARCHITECTURE-PLAN.md"
        "$CLAUDE_PROJECT_DIR/planning/phase${PHASE}/wave${WAVE}/WAVE-${PHASE}-${WAVE}-IMPLEMENTATION-PLAN.md"
    )

    local ALL_EXIST=true
    for doc in "${REQUIRED_DOCS[@]}"; do
        if [ ! -f "$doc" ]; then
            echo "❌ CRITICAL: Missing required document: $doc"
            ALL_EXIST=false
        else
            echo "✅ Found: $(basename "$doc")"
        fi
    done

    if [ "$ALL_EXIST" = false ]; then
        echo "🚨🚨🚨 WAVE ${PHASE}-${WAVE} CANNOT START - MISSING PLANS!"
        echo "Action Required: Architect and Code Reviewer must create wave plans first"
        exit 1
    fi

    echo "✅ Wave ${PHASE}-${WAVE} plans validated"
}
```

### 3. Effort Start Validation
**BEFORE creating effort infrastructure:**
```bash
validate_effort_prerequisites() {
    local PHASE="$1"
    local WAVE="$2"
    local EFFORT="$3"

    echo "🔍 VALIDATING EFFORT $EFFORT PREREQUISITES..."

    # Wave plans MUST exist
    validate_wave_plans "$PHASE" "$WAVE"

    # Check for effort-specific plan (created by Code Reviewer)
    local EFFORT_PLAN_PATTERN="$CLAUDE_PROJECT_DIR/efforts/phase${PHASE}/wave${WAVE}/${EFFORT}/.software-factory/*/IMPLEMENTATION-PLAN--*.md"

    # Note: Effort plan may not exist yet if Code Reviewer hasn't created it
    # But phase and wave plans are MANDATORY

    echo "✅ Effort prerequisites validated (phase/wave plans exist)"
}
```

### 4. Project Start Validation
**BEFORE starting project work:**
```bash
validate_project_plans() {
    echo "🔍 VALIDATING PROJECT-LEVEL PLANS..."

    # MANDATORY project documents
    local REQUIRED_DOCS=(
        "$CLAUDE_PROJECT_DIR/planning/project/PROJECT-ARCHITECTURE-PLAN.md"
        "$CLAUDE_PROJECT_DIR/planning/project/PROJECT-IMPLEMENTATION-PLAN.md"
    )

    local ALL_EXIST=true
    for doc in "${REQUIRED_DOCS[@]}"; do
        if [ ! -f "$doc" ]; then
            echo "❌ CRITICAL: Missing required document: $doc"
            ALL_EXIST=false
        else
            echo "✅ Found: $(basename "$doc")"
        fi
    done

    if [ "$ALL_EXIST" = false ]; then
        echo "🚨🚨🚨 PROJECT CANNOT START - MISSING PLANS!"
        echo "Action Required: Create project-level plans first"
        exit 1
    fi

    echo "✅ Project plans validated"
}
```

## Integration Points

### 1. Orchestrator States That MUST Validate

#### CREATE_NEXT_INFRASTRUCTURE
```bash
# MANDATORY validation before creating any infrastructure
validate_wave_plans "$CURRENT_PHASE" "$CURRENT_WAVE"

# Proceed only if validation passes
echo "✅ Plans validated - proceeding with infrastructure creation"
```

#### SPAWN_SW_ENGINEERS
```bash
# MANDATORY validation before spawning SW Engineers
validate_effort_prerequisites "$PHASE" "$WAVE" "$EFFORT_NAME"

# Block spawn if validation fails
if [ $? -ne 0 ]; then
    echo "🚨 CANNOT SPAWN AGENTS - Plans missing!"
    update_state "ERROR_RECOVERY"
    exit 1
fi
```

#### WAVE_START
```bash
# MANDATORY validation at wave start
validate_wave_plans "$PHASE" "$WAVE"

# Cannot proceed without wave plans
```

#### START_PHASE_ITERATION (implicit in first wave)
```bash
# When starting first wave of a phase
if [ "$WAVE" = "1" ]; then
    validate_phase_plans "$PHASE"
fi
```

### 2. SW Engineer Startup Validation
```bash
# In SW Engineer initialization
echo "🔍 Validating required plans exist..."

# Check wave plan exists (MANDATORY)
WAVE_PLAN="$CLAUDE_PROJECT_DIR/planning/phase${PHASE}/wave${WAVE}/WAVE-${PHASE}-${WAVE}-IMPLEMENTATION-PLAN.md"
if [ ! -f "$WAVE_PLAN" ]; then
    echo "❌ CRITICAL: Wave plan missing!"
    echo "Expected: $WAVE_PLAN"
    echo "Cannot proceed without wave plan!"
    exit 1
fi

# Check for effort implementation plan
EFFORT_PLAN=$(find .software-factory -name "IMPLEMENTATION-PLAN--*.md" 2>/dev/null | head -1)
if [ -z "$EFFORT_PLAN" ]; then
    echo "⚠️ WARNING: No effort-specific implementation plan found"
    echo "Will use wave plan for guidance"
fi
```

### 3. Code Reviewer Validation
```bash
# Before creating effort plans
validate_wave_plans "$PHASE" "$WAVE"

echo "✅ Wave plans exist - can create effort plans"
```

### 4. Architect Validation
```bash
# Before phase/wave assessment
if [ "$ASSESSMENT_TYPE" = "phase" ]; then
    validate_phase_plans "$PHASE"
elif [ "$ASSESSMENT_TYPE" = "wave" ]; then
    validate_wave_plans "$PHASE" "$WAVE"
fi
```

## State Machine Integration

### State Transition Blocks
The following transitions are BLOCKED without plans:
```yaml
blocked_transitions:
  - from: "INIT"
    to: "WAVE_START"
    requires: "project_plans"

  - from: "WAVE_START"
    to: "CREATE_NEXT_INFRASTRUCTURE"
    requires: "wave_plans"

  - from: "CREATE_NEXT_INFRASTRUCTURE"
    to: "SPAWN_SW_ENGINEERS"
    requires: "wave_plans"

  - from: "*"
    to: "SPAWN_SW_ENGINEERS"
    requires: "phase_and_wave_plans"
```

### Error State Triggers
Missing plans trigger ERROR_RECOVERY:
```bash
if ! validate_wave_plans "$PHASE" "$WAVE"; then
    update_state "ERROR_RECOVERY"
    echo "reason: Missing required planning documents" > error_reason.txt
    exit 1
fi
```

## Common Violations and Fixes

### ❌ VIOLATION: Starting Phase 3 without Phase 1/2 plans
```bash
# WRONG - No validation
echo "Starting Phase 3 work..."
create_infrastructure "phase3/wave1/effort-001"

# CORRECT - Validate first
validate_phase_plans "1" || exit 1
validate_phase_plans "2" || exit 1
validate_phase_plans "3" || exit 1
echo "All phase plans validated - proceeding"
```

### ❌ VIOLATION: Creating infrastructure without wave plans
```bash
# WRONG - No plan check
git clone "$REPO" "efforts/phase2/wave1/effort-001"

# CORRECT - Check plans first
validate_wave_plans "2" "1" || exit 1
git clone "$REPO" "efforts/phase2/wave1/effort-001"
```

### ❌ VIOLATION: Spawning agents without validation
```bash
# WRONG - Direct spawn
/spawn sw-engineer phase2-wave1-effort-001

# CORRECT - Validate prerequisites
validate_effort_prerequisites "2" "1" "effort-001" || exit 1
/spawn sw-engineer phase2-wave1-effort-001
```

## Enforcement Mechanisms

### 1. Pre-Infrastructure Hook
```bash
# Install in CREATE_NEXT_INFRASTRUCTURE
pre_infrastructure_hook() {
    echo "🔒 Running plan validation hook..."
    validate_wave_plans "$1" "$2" || return 1
    echo "✅ Hook passed - infrastructure can be created"
}
```

### 2. Pre-Spawn Hook
```bash
# Install in SPAWN_SW_ENGINEERS
pre_spawn_hook() {
    echo "🔒 Running spawn validation hook..."
    validate_effort_prerequisites "$@" || return 1
    echo "✅ Hook passed - agents can be spawned"
}
```

### 3. State File Tracking
```json
{
  "plan_validation": {
    "project_plans_validated": true,
    "phase1_plans_validated": true,
    "phase2_plans_validated": false,
    "phase1_wave1_plans_validated": true,
    "last_validation": "2025-01-20T10:30:00Z",
    "validation_errors": []
  }
}
```

## Recovery Protocol

When plans are missing:
1. **IMMEDIATE STOP** - No further work
2. **Identify missing documents** - List exactly what's needed
3. **Spawn appropriate agents** - Architect/Code Reviewer for planning
4. **Wait for completion** - Monitor planning progress
5. **Re-validate** - Check all documents exist
6. **Resume only after validation** - Continue from stopped point

## Grading Impact

### Penalties for Violations:
- Starting phase without phase plans: **-50%**
- Starting wave without wave plans: **-40%**
- Creating infrastructure without validation: **-30%**
- Spawning agents without plan check: **-25%**
- Proceeding after validation failure: **-100% AUTOMATIC FAILURE**

### Required Evidence:
- Validation output in logs
- Plan existence checks before state transitions
- ERROR_RECOVERY triggers for missing plans
- State file tracking of validations

## Related Rules
- R109: Planning requirements
- R054: Implementation plan creation
- R303: Phase/wave document location protocol
- R210: Architect architecture planning
- R211: Code reviewer implementation planning
- R234: Mandatory state traversal
- R340: Planning file metadata tracking

## Summary

**ABSOLUTE REQUIREMENT**: No work can begin without proper planning documents!

This rule prevents the catastrophic failure where development proceeds without architectural guidance. Every phase needs plans, every wave needs plans, and the system MUST verify they exist BEFORE any implementation work begins.

---
*Rule Type*: Validation Gate
*Agents*: All Agents
*Enforcement*: Blocking validation at all work initiation points
*Created*: Emergency response to Phase 3 starting without Phase 1/2 plans
## State Manager Coordination (SF 3.0)

State Manager validates plan gates during state transitions:
- **PLANNING → SETUP_WAVE_INFRASTRUCTURE**: Requires valid plan loaded in `orchestrator-state-v3.json.project_progression`
- **Guard conditions**: Check plan completeness before infrastructure creation
- **Atomic validation**: Plan structure + state file structure must match

Orchestrator loads plan → State Manager shutdown validates plan fields → Transition allowed only if complete.

See: `orchestrator-state-v3.json` (.project_progression.phases[]), state machine guards
