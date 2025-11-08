# SPAWN_ARCHITECT_PHASE_PLANNING State Rules

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
## State Purpose
Spawn architect agent to create phase-level architecture and implementation plans per R210.

## Entry Conditions
- Arrived from START_PHASE_ITERATION
- Phase planning needed (no existing plans)
- orchestrator-state-v3.json has current_phase set

## Spawning Protocol

### 1. Prepare Architect Context
```bash
# Get current phase
PHASE=$(jq -r '.current_phase' orchestrator-state-v3.json)
echo "Preparing to spawn architect for Phase ${PHASE} planning..."

# Create planning directory if needed
PLANNING_DIR="$CLAUDE_PROJECT_DIR/planning/phase${PHASE}"
mkdir -p "$PLANNING_DIR"

# Gather context for architect
echo "Gathering context for architect..."

# Check for previous phase assessments
PREV_PHASE=$((PHASE - 1))
if [ "$PREV_PHASE" -gt 0 ]; then
    PREV_ASSESSMENT=$(ls -t "$CLAUDE_PROJECT_DIR/assessments/phase${PREV_PHASE}/PHASE-${PREV_PHASE}-ASSESSMENT--"*.md 2>/dev/null | head -1)
    if [ -n "$PREV_ASSESSMENT" ]; then
        echo "Found previous phase assessment: $(basename "$PREV_ASSESSMENT")"
    fi
fi
```

### 2. Create Architect Instructions (R210)
```bash
cat > /tmp/architect-phase-${PHASE}-instructions.md << EOF
# Phase ${PHASE} Architecture Planning Instructions

You are the Architect agent being spawned to create phase-level plans per R210.

## State Context
- Current State: PHASE_ARCHITECTURE_PLANNING
- Phase: ${PHASE}
- Your Role: Create comprehensive phase architecture and implementation plans

## R210: Architecture Planning Protocol
Per R210, you MUST create phase-level plans BEFORE any wave planning can begin.

## Required Deliverables

### 1. Phase Architecture Plan
**Location**: \`$CLAUDE_PROJECT_DIR/planning/phase${PHASE}/PHASE-${PHASE}-ARCHITECTURE-PLAN--\$(date +%Y%m%d-%H%M%S).md\`

Must include:
- Phase vision alignment with master plan
- Analysis of previous phases (what was built, what was learned)
- Core architectural decisions for this phase
- APIs and contracts (MUST BE FIRST!)
- Abstractions and interfaces
- Shared libraries from previous phases
- Parallelization strategy
- Integration points

### 2. Phase Implementation Plan
**Location**: \`$CLAUDE_PROJECT_DIR/planning/phase${PHASE}/PHASE-${PHASE}-PLAN--\$(date +%Y%m%d-%H%M%S).md\`

Must include:
- Phase overview and objectives
- Wave breakdown (how many waves, what each accomplishes)
- Effort sequencing within waves
- Dependencies between efforts
- Parallelization opportunities
- Risk mitigation strategies
- Success criteria

### 3. Pre-Calculate Phase Integration Infrastructure (R504 REQUIREMENT)

You MUST also pre-calculate phase-level integration infrastructure and update orchestrator-state-v3.json:

**For the phase integration branch:**
- Calculate exact branch name (e.g., project-prefix/phase1-integration)
- Calculate exact remote branch name (e.g., project-prefix/phase1-integration)
- Determine base branch (last wave integration of this phase, e.g., phase1-wave2-integration)
- Specify target repository URL (CRITICAL - must match target-repo-config.yaml)
- Calculate directory path (e.g., /absolute/path/efforts/phase1/integration)
- **DEMO PRE-PLANNING (R330 + R504)**: Pre-calculate demo information:
  - demo_script_file: Full path to demo script (e.g., $CLAUDE_PROJECT_DIR/efforts/phase1/integration-phase1/.software-factory/phase1/demo/phase1/integration/demo-phase-integration.sh)
  - demo_description: Summary of what demo will showcase (update when plan changes)
  - demo_plan_file: Path to demo plan (e.g., $CLAUDE_PROJECT_DIR/efforts/phase1/integration-phase1/.software-factory/phase1/demo/demo-plan.md)

Update orchestrator-state-v3.json with:
\`\`\`json
"pre_planned_infrastructure": {
  "integrations": {
    "phase_integrations": {
      "phase${PHASE}": {
        "phase": "phase${PHASE}",
        "branch_name": "project-prefix/phase${PHASE}-integration",
        "remote_branch": "project-prefix/phase${PHASE}-integration",
        "base_branch": "project-prefix/phase${PHASE}-wave${LAST_WAVE}-integration",
        "target_repo_url": "https://github.com/owner/repo.git",  // MUST MATCH TARGET REPO!
        "directory": "/absolute/path/efforts/phase${PHASE}/integration",
        "component_waves": ["phase${PHASE}_wave1", "phase${PHASE}_wave2", ...],
        "created": false,
        "validated": false,
        "demo_script_file": "$CLAUDE_PROJECT_DIR/efforts/phase${PHASE}/integration-phase${PHASE}/.software-factory/phase${PHASE}/demo/phase${PHASE}/integration/demo-phase-integration.sh",
        "demo_description": "Showcases phase ${PHASE} integration validating <phase objective> with <all wave components>",
        "demo_plan_file": "$CLAUDE_PROJECT_DIR/efforts/phase${PHASE}/integration-phase${PHASE}/.software-factory/phase${PHASE}/demo/demo-plan.md"
      }
    }
  }
}
\`\`\`

**NOTE**: Wave-level infrastructure will be pre-planned during wave planning, but phase integration infrastructure must be pre-planned here!

## Required Reading
1. Master plan: \`$CLAUDE_PROJECT_DIR/PROJECT-IMPLEMENTATION-PLAN.md\`
2. Current state: \`$CLAUDE_PROJECT_DIR/orchestrator-state-v3.json\`
3. Previous phase assessment (if exists): \`${PREV_ASSESSMENT}\`
4. Existing codebase from completed phases

## State Machine Requirements
- You are in state: PHASE_ARCHITECTURE_PLANNING
- After creating both plans, stop and wait
- The orchestrator will transition to WAITING_FOR_PHASE_PLANS
- Do NOT proceed to wave planning - that comes later

## Critical Success Factors
- Both plans must be created (not just one)
- Plans must follow the exact naming convention
- Architecture decisions drive all subsequent implementation
- Wave breakdown must be realistic and achievable
- Parallelization opportunities must be identified

## R313: Stop After Completion
Per R313, after creating both plans you MUST STOP.
Do not continue to wave planning or any other state.
EOF

echo "Instructions prepared for architect"
```

### 3. Spawn the Architect
```bash
echo "═══════════════════════════════════════════════════════"
echo "Spawning architect for Phase ${PHASE} planning..."
echo "═══════════════════════════════════════════════════════"

# Update orchestrator state before spawning
jq ".phases.phase_${PHASE}.planning_started = \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\" |
    .phases.phase_${PHASE}.status = \"planning\"" \
   orchestrator-state-v3.json > orchestrator-state.tmp && \
   mv orchestrator-state.tmp orchestrator-state-v3.json

# Spawn the architect with instructions
/spawn architect PHASE_ARCHITECTURE_PLANNING /tmp/architect-phase-${PHASE}-instructions.md

echo "✅ Architect spawned successfully"
```

### 4. Update State and Stop (R313)
```bash
# R313: Must stop after spawning
echo "═══════════════════════════════════════════════════════"
echo "R313: Stopping after spawn - DO NOT CONTINUE"
echo "═══════════════════════════════════════════════════════"

# Record spawn event
jq ".phases.phase_${PHASE}.architect_spawned = \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\"" \
   orchestrator-state-v3.json > orchestrator-state.tmp && \
   mv orchestrator-state.tmp orchestrator-state-v3.json

# Transition to waiting state
update_state "WAITING_FOR_PHASE_PLANS" "Architect spawned for Phase ${PHASE} planning, waiting for completion"

echo "State updated to WAITING_FOR_PHASE_PLANS"
echo "Orchestrator MUST STOP NOW per R313"
```

## Exit Conditions
- Architect spawned with complete instructions
- orchestrator-state-v3.json updated with spawn event
- State transitioned to WAITING_FOR_PHASE_PLANS
- Orchestrator stops per R313

## State Transitions
- **WAITING_FOR_PHASE_PLANS**: Always (after successful spawn)
- **ERROR_RECOVERY**: Only if spawn fails

## Error Handling
```bash
handle_spawn_error() {
    local error_type="$1"
    local error_details="$2"

    echo "ERROR: Failed to spawn architect"
    echo "Type: $error_type"
    echo "Details: $error_details"

    # Record error in state
    jq --arg err "$error_type: $error_details" \
       '.phases.phase_'$PHASE'.spawn_error = $err' \
       orchestrator-state-v3.json > orchestrator-state.tmp && \
       mv orchestrator-state.tmp orchestrator-state-v3.json

    update_state "ERROR_RECOVERY" "Architect spawn failed: $error_details"
}
```

## R313 Enforcement
```bash
# This is the ABSOLUTE LAST thing that happens
echo ""
echo "🛑 R313 ENFORCEMENT: STOPPING NOW"
echo "The orchestrator MUST NOT continue past this point."
echo "Wait for /continue command to resume from WAITING_FOR_PHASE_PLANS"
```

## R505: Infrastructure Synchronization
```bash
# After architect completes phase planning, synchronize infrastructure
sync_phase_infrastructure_after_planning() {
    local PHASE=$1

    echo "═══════════════════════════════════════════════════════"
    echo "R505: Synchronizing pre_planned_infrastructure from phase plans"
    echo "═══════════════════════════════════════════════════════"

    # This will be executed in WAITING_FOR_PHASE_PLANS state
    # after architect completes planning
}
```

## Exit Conditions and Continuation Flag

### ✅ CONTINUE-SOFTWARE-FACTORY=TRUE (Standard Operation - DEFAULT)

**Use TRUE when:**
- ✅ Architect agent spawned successfully
- ✅ Context provided about current phase
- ✅ Phase planning requirements specified
- ✅ State file updated with spawn tracking
- ✅ Ready to transition to waiting state
- ✅ Following designed workflow

**THIS IS NORMAL WORKFLOW.** Spawning the Architect for phase planning is the
DESIGNED PROCESS. This is automation working correctly.

### 🔴 CONTINUE-SOFTWARE-FACTORY=FALSE (Exceptional Conditions ONLY)

**Use FALSE ONLY when:**
- ❌ Cannot determine current phase
- ❌ Phase infrastructure missing or corrupt
- ❌ State machine corruption detected
- ❌ Cannot spawn Architect agent
- ❌ Required planning metadata missing

**DO NOT set FALSE because:**
- ❌ Spawning Architect (NORMAL workflow!)
- ❌ Waiting for planning (EXPECTED process)
- ❌ R322 requires stop (stop ≠ FALSE!)
- ❌ Transitioning to waiting state (NORMAL)

### Critical Distinction: R322 Stop vs Continuation Flag

**Correct pattern:**
```bash
# Spawn architect
# Update state
exit 0  # R322 stop
```

**Last line:**
```bash
echo "CONTINUE-SOFTWARE-FACTORY=TRUE"  # Normal operation
```

**Grading Impact:** Using FALSE for normal spawning: -20%

## Automation Flag
```bash
# After successful spawn and state transition:
echo "CONTINUE-SOFTWARE-FACTORY=TRUE"  # Continue orchestration per R405 - system will handle monitoring
```
## 🔴🔴🔴 R405 - MANDATORY AUTOMATION CONTINUATION FLAG 🔴🔴🔴

**SUPREME LAW - PENALTY FOR VIOLATION: -100% GRADE**

**MUST output as LAST line before state completion:**
- `CONTINUE-SOFTWARE-FACTORY=TRUE` - State completed, automation can continue
- `CONTINUE-SOFTWARE-FACTORY=FALSE` - Error/block, manual intervention required

**See: $CLAUDE_PROJECT_DIR/rule-library/R405-automation-continuation-flag.md**
**See: $CLAUDE_PROJECT_DIR/rule-library/R405-CONTINUATION-FLAG-MASTER-GUIDE.md**
