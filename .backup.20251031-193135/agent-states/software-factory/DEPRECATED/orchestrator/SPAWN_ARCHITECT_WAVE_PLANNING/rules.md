# SPAWN_ARCHITECT_WAVE_PLANNING State Rules

## State Context

**Current Phase**: Read from `orchestrator-state-v3.json`
**Current Wave**: Read from `orchestrator-state-v3.json`
**Entry From**: WAVE_START
**Exit To**: WAITING_FOR_ARCHITECTURE_PLAN

## 🚨🚨🚨 PRIMARY OBJECTIVE 🚨🚨🚨

**SPAWN architect agent to create wave-level implementation plan per R210 (Architect Architecture Planning Protocol).**

This state ensures that BEFORE any effort planning or implementation begins, the architect creates:
1. Wave architecture design and technical approach
2. Wave implementation plan with effort breakdown

## Required Inputs

### 1. Phase Context
```bash
PHASE=$(jq -r '.current_phase' orchestrator-state-v3.json)
WAVE=$(jq -r '.current_wave' orchestrator-state-v3.json)

echo "📋 Phase: $PHASE, Wave: $WAVE"
```

### 2. Phase Plans (Prerequisites)
```bash
# Verify phase plans exist
PHASE_ARCH_PLAN=$(ls -t "$CLAUDE_PROJECT_DIR/phase-plans/phase${PHASE}/PHASE-${PHASE}-ARCHITECTURE-PLAN--"*.md 2>/dev/null | head -1)
PHASE_IMPL_PLAN=$(ls -t "$CLAUDE_PROJECT_DIR/phase-plans/phase${PHASE}/PHASE-${PHASE}-PLAN--"*.md 2>/dev/null | head -1)

if [ -z "$PHASE_ARCH_PLAN" ] || [ -z "$PHASE_IMPL_PLAN" ]; then
    echo "❌ FATAL: Cannot create wave plan without phase plans"
    echo "  Missing: Phase $PHASE architecture or implementation plan"
    exit 1
fi
```

### 3. Target Repository Configuration
```bash
# Read target repo configuration
if [ ! -f "$CLAUDE_PROJECT_DIR/target-repo-config.yaml" ]; then
    echo "❌ FATAL: target-repo-config.yaml not found"
    exit 1
fi

TARGET_REPO=$(yq -r '.target_repository' "$CLAUDE_PROJECT_DIR/target-repo-config.yaml")
PROJECT_PREFIX=$(yq -r '.project_prefix' "$CLAUDE_PROJECT_DIR/target-repo-config.yaml")
```

## Architect Spawn Requirements

### Task Prompt Template
```bash
cat << EOF
📋 ARCHITECT TASK: Wave ${PHASE}.${WAVE} Architecture & Implementation Planning

You are tasked with creating architecture and implementation plans for Wave ${WAVE} of Phase ${PHASE}.

**Context:**
- Project: ${PROJECT_PREFIX}
- Phase: ${PHASE}
- Wave: ${WAVE}
- Phase Architecture Plan: ${PHASE_ARCH_PLAN}
- Phase Implementation Plan: ${PHASE_IMPL_PLAN}

**Your Deliverables:**

1. **Wave Architecture Plan**
   Location: phase-plans/phase${PHASE}/wave${WAVE}/WAVE-${PHASE}-${WAVE}-ARCHITECTURE--$(date +%Y%m%d-%H%M%S).md

   Must contain:
   - Technical approach for this wave's functionality
   - Component designs and interfaces
   - Technology choices and patterns
   - Risk assessment and mitigation strategies

2. **Wave Implementation Plan**
   Location: phase-plans/phase${PHASE}/wave${WAVE}/WAVE-${PHASE}-${WAVE}-PLAN--$(date +%Y%m%d-%H%M%S).md

   Must contain per R210, R303, and R504:
   - Effort breakdown (each effort: 300-600 lines target)
   - Dependencies between efforts
   - **🚀 Parallelization Strategy** section with:
     - Sequential groups (blocking dependencies)
     - Parallel groups (independent work)
     - Orchestrator spawning strategy
   - Test requirements per effort
   - Acceptance criteria

3. **Pre-Calculate Infrastructure (R504 REQUIREMENT)**

   You MUST also pre-calculate ALL infrastructure and update orchestrator-state-v3.json:

   **A. For EACH effort in the wave:**
   - Calculate exact branch name (e.g., project-prefix/phase1/wave1/effort-name)
   - Calculate exact directory path (e.g., /absolute/path/efforts/phase1/wave1/effort-name/)
   - Determine remote branch name (e.g., origin/project-prefix/phase1/wave1/effort-name)
   - Specify target remote (usually "target")
   - Define split pattern (e.g., effort-name--split-NNN)
   - **CRITICAL**: Include target_repo_url for EVERY effort

   **B. For wave integration infrastructure:**
   - Calculate wave integration branch name (e.g., project-prefix/phase1-wave1-integration)
   - Determine integration workspace directory (e.g., /absolute/path/efforts/phase1/wave1/integration-workspace/)
   - Specify base branch for wave integration (previous wave integration or main)
   - **CRITICAL**: Include target_repo_url matching the target repository
   - **DEMO PRE-PLANNING (R330 + R504)**: Pre-calculate demo information:
     - demo_script_file: Full path to demo script (e.g., $CLAUDE_PROJECT_DIR/efforts/phase1/wave1/integration-phase1-wave1/.software-factory/phase1/wave1/demo/phase1/wave1/integration/demo-wave-integration.sh)
     - demo_description: Summary of what demo will showcase (update when plan changes)
     - demo_plan_file: Path to demo plan (e.g., $CLAUDE_PROJECT_DIR/efforts/phase1/wave1/integration-phase1-wave1/.software-factory/phase1/wave1/demo/demo-plan.md)

   Update orchestrator-state-v3.json with:
   ```json
   "pre_planned_infrastructure": {
     "validated": false,  // Will be validated in next state
     "target_repo_url": "https://github.com/owner/repo.git",
     "efforts": {
       "phase1_wave1_effort1": {
         "full_path": "/absolute/path/",
         "branch_name": "exact-branch-name",
         "remote_branch": "origin/exact-branch-name",
         "target_remote": "target",
         "phase": "phase1",
         "wave": "wave1",
         "effort_name": "effort1",
         "target_repo_url": "https://github.com/owner/repo.git",
         "created": false
       }
     },
     "integrations": {
       "wave_integrations": {
         "phase1_wave1": {
           "phase": "phase1",
           "wave": "wave1",
           "branch_name": "project-prefix/phase1-wave1-integration",
           "remote_branch": "project-prefix/phase1-wave1-integration",
           "base_branch": "main",  // or previous wave integration
           "target_repo_url": "https://github.com/owner/repo.git",
           "directory": "/absolute/path/efforts/phase1/wave1/integration-workspace",
           "component_efforts": ["E1.1.1", "E1.1.2"],
           "created": false,
           "validated": false,
           "demo_script_file": "$CLAUDE_PROJECT_DIR/efforts/phase1/wave1/integration-phase1-wave1/.software-factory/phase1/wave1/demo/phase1/wave1/integration/demo-wave-integration.sh",
           "demo_description": "Showcases wave 1 integration validating <feature> with <components>",
           "demo_plan_file": "$CLAUDE_PROJECT_DIR/efforts/phase1/wave1/integration-phase1-wave1/.software-factory/phase1/wave1/demo/demo-plan.md"
         }
       }
     }
   }
   ```

**Critical Requirements:**
- Follow R303 (Phase/Wave Document Location Protocol)
- Follow R504 (Pre-Infrastructure Planning) - MUST populate pre_planned_infrastructure
- Each effort must be sized 300-600 lines (split if needed)
- Clearly identify which efforts can run in parallel
- Specify cascade base branches per R308
- Include test planning per R219
- Pre-calculate ALL infrastructure details upfront

**R504 ENFORCEMENT:**
The CREATE_NEXT_INFRASTRUCTURE state will FAIL if pre_planned_infrastructure is not populated!
No infrastructure can be created without pre-planning.

**Working Directory:** $CLAUDE_PROJECT_DIR
**Branch:** $(git branch --show-current)

EOF
```

### Spawn Command
```bash
# Use Task tool to spawn architect agent
# Subagent type: architect
# Task: Wave architecture and implementation planning
# Directory: $CLAUDE_PROJECT_DIR
# Expected output location: phase-plans/phase${PHASE}/wave${WAVE}/
```

## State Update

### Update orchestrator-state-v3.json
```bash
jq --arg phase "$PHASE" \
   --arg wave "$WAVE" \
   --arg timestamp "$(date -Iseconds)" \
   '.state_machine.current_state = "WAITING_FOR_ARCHITECTURE_PLAN" |
    .state_machine.previous_state = "SPAWN_ARCHITECT_WAVE_PLANNING" |
    .current_phase = ($phase | tonumber) |
    .current_wave = ($wave | tonumber) |
    .state_transition_log += [{
        "from": "SPAWN_ARCHITECT_WAVE_PLANNING",
        "to": "WAITING_FOR_ARCHITECTURE_PLAN",
        "timestamp": $timestamp,
        "reason": "Architect spawned for Wave \($phase).\($wave) planning"
    }]' orchestrator-state-v3.json > tmp.json && mv tmp.json orchestrator-state-v3.json
```

## Validation Requirements

### Pre-Spawn Validation
- ✅ Phase plans exist
- ✅ Wave number is valid (1-N)
- ✅ No existing wave plan for this wave
- ✅ Target repo config accessible
- ✅ Working directory is correct

### Post-Spawn Validation
- ✅ Architect agent was spawned successfully
- ✅ State file updated to WAITING_FOR_ARCHITECTURE_PLAN
- ✅ Timestamp recorded

## Integration with Rules

- **R210**: Architect Architecture Planning Protocol (MISSION CRITICAL)
- **R303**: Phase/Wave Document Location Protocol
- **R308**: Cascade Base Branch Calculation
- **R219**: Dependency-Aware Effort Planning
- **R502**: Plan Validation Gates
- **R053**: Parallelization Decision Criteria

## Exit Criteria

Before transitioning to WAITING_FOR_ARCHITECTURE_PLAN:
- ✅ Architect agent spawned with complete context
- ✅ Expected deliverable locations specified
- ✅ State file updated
- ✅ All validations passed

## Common Issues

### Issue: Phase Plans Missing
**Detection**: Cannot find PHASE-X-ARCHITECTURE-PLAN or PHASE-X-PLAN files
**Resolution**: Transition back to START_PHASE_ITERATION to create phase plans first

### Issue: Wave Already Exists
**Detection**: Wave plan already exists in phase-plans/phaseX/waveY/
**Resolution**: Skip to CREATE_NEXT_INFRASTRUCTURE if wave already planned

### Issue: Invalid Wave Number
**Detection**: Wave number doesn't follow sequential ordering
**Resolution**: Validate wave sequence, correct in state file

## R505: Infrastructure Synchronization Note

```bash
# Note: Infrastructure synchronization will occur in WAITING_FOR_ARCHITECTURE_PLAN
# after the architect completes wave planning and creates the wave plan
echo "📋 R505: Infrastructure sync will occur after wave plan creation"
```

## R313 Enforcement - MANDATORY STOP

```bash
# This is the ABSOLUTE LAST thing that happens in this state
echo ""
echo "🛑 R313 ENFORCEMENT: STOPPING NOW"
echo "The orchestrator MUST NOT continue past this point."
echo "Architect has been spawned for wave planning."
echo "Wait for architect to complete before resuming."
echo ""
echo "Next state will be: WAITING_FOR_ARCHITECTURE_PLAN"
echo "Resume with: /continue-orchestrating"
```

## Exit Conditions and Continuation Flag

### ✅ CONTINUE-SOFTWARE-FACTORY=TRUE (Standard Operation - DEFAULT)

**Use TRUE when:**
- ✅ Architect agent spawned successfully for wave planning
- ✅ Wave context and requirements provided
- ✅ State file updated with spawn tracking
- ✅ Ready to transition to WAITING_FOR_ARCHITECTURE_PLAN
- ✅ Following designed workflow

**THIS IS NORMAL WORKFLOW.** Spawning the Architect for wave planning is the DESIGNED PROCESS.

### 🔴 CONTINUE-SOFTWARE-FACTORY=FALSE (Exceptional Conditions ONLY)

**Use FALSE ONLY when:**
- ❌ Cannot determine current wave
- ❌ Wave infrastructure missing or corrupt
- ❌ Cannot spawn Architect agent
- ❌ State machine corruption

**DO NOT set FALSE because:**
- ❌ Spawning Architect (NORMAL!)
- ❌ R322 requires stop (stop ≠ FALSE!)
- ❌ Waiting for planning (EXPECTED!)

**Correct pattern:** `exit 0` + `CONTINUE-SOFTWARE-FACTORY=TRUE`

## Automation Flag

```bash
# After successful spawn and state transition:
echo "CONTINUE-SOFTWARE-FACTORY=TRUE"  # Continue orchestration per R405 - system will handle monitoring
```

---

**REMEMBER**: The architect creates the wave plan that determines effort parallelization. The orchestrator IMPLEMENTS the architect's parallelization strategy, it does NOT decide it!

## 🔴🔴🔴 R405 - MANDATORY AUTOMATION CONTINUATION FLAG 🔴🔴🔴

**SUPREME LAW - PENALTY FOR VIOLATION: -100% GRADE**

**MUST output as LAST line before state completion:**
- `CONTINUE-SOFTWARE-FACTORY=TRUE` - State completed, automation can continue
- `CONTINUE-SOFTWARE-FACTORY=FALSE` - Error/block, manual intervention required

**See: $CLAUDE_PROJECT_DIR/rule-library/R405-automation-continuation-flag.md**
**See: $CLAUDE_PROJECT_DIR/rule-library/R405-CONTINUATION-FLAG-MASTER-GUIDE.md**
