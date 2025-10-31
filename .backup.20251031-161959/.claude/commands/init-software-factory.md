---
name: init-software-factory
description: Initialize a new Software Factory 3.0 project from idea to implementation plan
---

# /init-software-factory

╔═══════════════════════════════════════════════════════════════════════════════╗
║                        SOFTWARE FACTORY 3.0                                  ║
║                   PROJECT INITIALIZATION COMMAND                              ║
║                                                                               ║
║ Purpose: Transform a PRD/project description into a complete implementation  ║
║          plan with master → phase → wave architecture using progressive      ║
║          planning and test-driven development                                 ║
║ Flow: INIT → SPAWN_ARCHITECT_MASTER_PLANNING → TEST_PLANNING →              ║
║       PHASE_PLANNING → WAVE_START (SF 3.0 mandatory sequence)               ║
╚═══════════════════════════════════════════════════════════════════════════════╝

## 📖 USAGE

### Interactive Mode (System Prompts for Description)
```bash
/init-software-factory
# System will ask for project name and description interactively
```

### Direct Mode (Provide Description with Command)
```bash
/init-software-factory "Project Description: Build an OCI image push command..."

# For long descriptions, use heredoc:
/init-software-factory "$(cat <<'EOF'
Project Description:
Implement idpbuilder push command to upload OCI images to Gitea registry at
https://gitea.cnoe.localtest.me:8443/ with authentication via -username and
-password flags, disable default certificate checks.

Context: idpbuilder is a Go CLI tool. It currently can build images but cannot
push them. We need to integrate go-containerregistry library for OCI operations.
EOF
)"
```

### What Happens
1. **You provide**: Project name + description (PRD can be any length)
2. **Architect asks**: Clarifying questions to fill gaps
3. **System generates**: Complete PROJECT-IMPLEMENTATION-PLAN.md with:
   - All phases broken down
   - All waves within each phase
   - All efforts (each ≤800 lines) with sizing justifications
   - Dependencies mapped
   - Technology stack specified
4. **You review**: Generated plan and request changes if needed
5. **System finalizes**: All state files and prepares for development

## 🎯 AGENT IDENTITY ASSIGNMENT

**You are the orchestrator in INITIALIZATION MODE**

By invoking this command, you are now operating as the orchestrator agent in project initialization mode. You must:
- Guide the user through interactive requirements gathering
- Spawn architect agent for Q&A sessions
- Generate all required configuration files
- Create comprehensive IMPLEMENTATION-PLAN.md
- Prepare project for /continue-orchestrating

## 📋 PRE-FLIGHT CHECKS

Before beginning initialization:

1. **Verify Clean State**:
   - Check current directory is a fresh SF 3.0 template instance
   - Verify no existing implementation plan files
   - Confirm orchestrator-state-v3.json shows current_state: "INIT"
   - All 4 state files exist (orchestrator-state-v3.json, bug-tracking.json, integration-containers.json)

2. **Required User Inputs**:
   - Project name/prefix (alphanumeric only)
   - Project idea/description (can be full PRD, any length)

## 🚀 INITIALIZATION PROCESS

### Phase 1: PARSE INPUT AND GATHER BASICS

```bash
# If user provided description as argument, extract it
if [[ -n "$1" ]]; then
    PROJECT_DESCRIPTION="$1"
    echo "Project description received: ${PROJECT_DESCRIPTION:0:100}..."
else
    # Interactive mode: prompt user
    echo "Please provide the following:"
    echo "1. Project name/prefix (alphanumeric only): "
    echo "2. Project description (can be full PRD, any length): "
    read PROJECT_NAME
    read -d '' PROJECT_DESCRIPTION  # Read until EOF for multi-line
fi

# Store in temporary init state file
cat > init-state-temp.json <<EOF
{
  "project_name": "$PROJECT_NAME",
  "initial_description": "$PROJECT_DESCRIPTION",
  "state": "INIT_START",
  "timestamp": "$(date -Iseconds)"
}
EOF
```

### Phase 2: SPAWN ARCHITECT FOR MASTER PLANNING (SF 3.0)

The orchestrator enters **SPAWN_ARCHITECT_MASTER_PLANNING** state and spawns architect agent to:

**Architect goes through 3 initialization states:**
1. **INIT_REQUIREMENTS_GATHERING** (Interactive Q&A):
   - Read initial description from init-state-temp.json
   - Identify what's MISSING from description (not what's provided)
   - Ask ONLY clarifying questions for missing information
   - Gather technical requirements (language, frameworks, testing, deployment)
   - **~20 questions across 6 categories** (only ask if information missing)

2. **INIT_DECOMPOSE_PRD** (Break down project):
   - Analyze complete requirements (description + Q&A responses)
   - Identify all major features from PRD
   - Estimate lines of code per feature (using sizing guidelines)
   - Split features >800 lines into multiple efforts
   - Organize efforts into waves (3-6 efforts per wave)
   - Organize waves into phases (2-4 waves per phase)
   - Map dependencies between efforts

3. **INIT_SYNTHESIZE_PLAN** (Generate master architecture):
   - Create **MASTER-ARCHITECTURE-PLAN.md** with:
     - High-level project architecture
     - Technology choices and rationale
     - Major component breakdown
     - Phase overview (what each phase delivers)
   - **Fidelity level**: Conceptual/strategic (not pseudocode yet)

**Architect Output**: `MASTER-ARCHITECTURE-PLAN.md` + requirements bundle in state file

After architect completes, orchestrator transitions to **WAITING_FOR_MASTER_ARCHITECTURE**.

### Phase 3: PROJECT-LEVEL TEST PLANNING (SF 3.0 - NEW!)

The orchestrator enters **SPAWN_CODE_REVIEWER_PROJECT_TEST_PLANNING** and spawns code reviewer to:

**Code reviewer creates**:
- **PROJECT-TEST-PLAN.md** with:
  - Overall testing strategy
  - Test categories (unit, integration, e2e)
  - Coverage targets per phase
  - Testing tools and frameworks
  - Test data management approach
  - CI/CD integration strategy

**This is R341 TDD compliance** - test planning BEFORE any implementation.

After code reviewer completes, orchestrator transitions to **WAITING_FOR_PROJECT_TEST_PLAN**.

### Phase 4: CREATE PROJECT INTEGRATION BRANCH (SF 3.0 - NEW!)

The orchestrator enters **CREATE_PROJECT_INTEGRATION_BRANCH_EARLY** and:

**Actions**:
```bash
git checkout -b project-integration main
git push -u origin project-integration
```

**This is R342 compliance** - early integration branch creation before implementation begins.

### Phase 5: PHASE-LEVEL ARCHITECTURE PLANNING (SF 3.0)

The orchestrator enters **SPAWN_ARCHITECT_PHASE_PLANNING** and spawns architect for Phase 1:

**Architect creates**:
- **PHASE-1-ARCHITECTURE-PLAN.md** with:
  - **Pseudocode fidelity** technical approach
  - Wave breakdown with detailed descriptions
  - Library/framework choices specific to this phase
  - API contracts and interfaces
  - Data models and schemas
  - Error handling patterns

**Fidelity level**: Pseudocode - specific enough for developers to implement.

After architect completes, orchestrator transitions to **WAITING_FOR_PHASE_ARCHITECTURE**.

### Phase 6: PHASE-LEVEL TEST PLANNING (SF 3.0)

The orchestrator enters **SPAWN_CODE_REVIEWER_PHASE_TEST_PLANNING** and spawns code reviewer:

**Code reviewer creates**:
- **PHASE-1-TEST-PLAN.md** with:
  - Test cases specific to Phase 1 features
  - Test data requirements for Phase 1
  - Integration test scenarios
  - Performance test criteria
  - Security test cases

After code reviewer completes, orchestrator transitions to **WAITING_FOR_PHASE_TEST_PLAN**.

### Phase 7: CREATE PHASE INTEGRATION BRANCH (SF 3.0)

The orchestrator enters **CREATE_PHASE_INTEGRATION_BRANCH_EARLY** and:

**Actions**:
```bash
git checkout -b phase-1-integration main
git push -u origin phase-1-integration
```

### Phase 8: PHASE IMPLEMENTATION PLANNING (SF 3.0)

The orchestrator enters **SPAWN_CODE_REVIEWER_PHASE_IMPL** and spawns code reviewer:

**Code reviewer creates**:
- **PHASE-1-IMPLEMENTATION-PLAN.md** with:
  - Detailed effort breakdown
  - Effort dependencies and ordering
  - Parallelization opportunities
  - File structure and organization
  - Code review checkpoints

After code reviewer completes, orchestrator transitions to **WAITING_FOR_PHASE_IMPLEMENTATION_PLAN**.

### Phase 9: FINALIZE AND TRANSITION TO WAVE_START

**Orchestrator verifies all planning artifacts exist:**
- ✅ MASTER-ARCHITECTURE-PLAN.md
- ✅ PROJECT-TEST-PLAN.md
- ✅ project-integration branch created
- ✅ PHASE-1-ARCHITECTURE-PLAN.md
- ✅ PHASE-1-TEST-PLAN.md
- ✅ phase-1-integration branch created
- ✅ PHASE-1-IMPLEMENTATION-PLAN.md

**All information loaded into orchestrator-state-v3.json:**
```json
{
  "state_machine": {
    "current_state": "WAVE_START"
  },
  "project_progression": {
    "current_project": {
      "project_id": "project-UUID",
      "name": "PROJECT_NAME",
      "status": "IN_PROGRESS"
    },
    "current_phase": {
      "phase_number": 1,
      "name": "Phase 1",
      "status": "IN_PROGRESS"
    },
    "current_wave": {
      "wave_number": 1,
      "name": "Wave 1",
      "status": "IN_PROGRESS"
    }
  },
  "planning_artifacts": {
    "master_architecture_file": "planning/MASTER-ARCHITECTURE-PLAN.md",
    "master_architecture_status": "CREATED",
    "phase_1_architecture_file": "planning/PHASE-1-ARCHITECTURE-PLAN.md",
    "phase_1_architecture_status": "CREATED"
  }
}
```

**Handoff message**:
```
✅ Software Factory 3.0 initialization complete!

Planning artifacts created:
- MASTER-ARCHITECTURE-PLAN.md (high-level architecture)
- PROJECT-TEST-PLAN.md (overall testing strategy)
- PHASE-1-ARCHITECTURE-PLAN.md (pseudocode-level technical plan)
- PHASE-1-TEST-PLAN.md (phase-specific test cases)
- PHASE-1-IMPLEMENTATION-PLAN.md (effort breakdown and dependencies)

Integration branches created:
- project-integration (from main)
- phase-1-integration (from main)

Current state: WAVE_START
Next step: Run /continue-software-factory to begin Wave 1 development
```

## 🔴🔴🔴 SF 3.0 STATE MANAGER BOOKEND PATTERN 🔴🔴🔴

**CRITICAL**: Every state transition during initialization MUST follow this pattern:

### Before EVERY State Transition:

```bash
# 1. Complete current state work
echo "✅ Current state work complete"

# 2. Set proposed next state
PROPOSED_NEXT_STATE="NEXT_STATE_NAME"
TRANSITION_REASON="Description of why transitioning"

# 3. Spawn State Manager for SHUTDOWN_CONSULTATION
WORK_RESULTS=$(cat <<EOF
{
  "state_completed": "CURRENT_STATE",
  "work_accomplished": ["item1", "item2"],
  "proposed_next_state": "$PROPOSED_NEXT_STATE",
  "transition_reason": "$TRANSITION_REASON"
}
EOF
)

/spawn state-manager SHUTDOWN_CONSULTATION \
  --current-state "CURRENT_STATE" \
  --proposed-next-state "$PROPOSED_NEXT_STATE" \
  --transition-reason "$TRANSITION_REASON" \
  --work-results "$WORK_RESULTS"

# State Manager validates and updates orchestrator-state-v3.json atomically

# 4. Save TODOs (R287)
save_todos "STATE_TRANSITION"
git add todos/*.todo
git commit -m "todo: orchestrator - state transition [R287]"
git push

# 5. Output R405 continuation flag
echo "CONTINUE-SOFTWARE-FACTORY=TRUE REASON=STATE_COMPLETE"

# 6. Stop for context preservation (R322)
exit 0
```

**This pattern is MANDATORY for:**
- INIT → SPAWN_ARCHITECT_MASTER_PLANNING
- WAITING_FOR_MASTER_ARCHITECTURE → SPAWN_CODE_REVIEWER_PROJECT_TEST_PLANNING
- WAITING_FOR_PROJECT_TEST_PLAN → CREATE_PROJECT_INTEGRATION_BRANCH_EARLY
- CREATE_PROJECT_INTEGRATION_BRANCH_EARLY → SPAWN_ARCHITECT_PHASE_PLANNING
- WAITING_FOR_PHASE_ARCHITECTURE → SPAWN_CODE_REVIEWER_PHASE_TEST_PLANNING
- WAITING_FOR_PHASE_TEST_PLAN → CREATE_PHASE_INTEGRATION_BRANCH_EARLY
- CREATE_PHASE_INTEGRATION_BRANCH_EARLY → SPAWN_CODE_REVIEWER_PHASE_IMPL
- WAITING_FOR_PHASE_IMPLEMENTATION_PLAN → WAVE_START

## 🔴 CRITICAL RULES - SF 3.0 COMPLIANCE

1. **SF3.0_FUNDAMENTAL: State Manager Authority**:
   - State Manager is ONLY agent that updates orchestrator-state-v3.json
   - Orchestrator MUST consult State Manager for EVERY transition
   - State files updated atomically via R288 protocol

2. **R203 STATE-AWARE STARTUP**: Must follow SF 3.0 initialization mandatory sequence

3. **R287 TODO PERSISTENCE**: Save TODOs within 30s after write, before transitions

4. **R288 STATE FILE UPDATES**: All 4 state files updated through State Manager bookend

5. **R322 MANDATORY STOP**: Stop after EVERY state completion for context preservation

6. **R405 AUTOMATION FLAG**: Output CONTINUE-SOFTWARE-FACTORY=TRUE/FALSE as LAST line

7. **R341 TDD REQUIREMENTS**: Test planning BEFORE implementation at all levels

8. **R342 EARLY INTEGRATION**: Create integration branches BEFORE implementation begins

9. **R506 PRE-COMMIT ENFORCEMENT**: NEVER use --no-verify, ALWAYS let hooks validate

10. **R006 ORCHESTRATOR NEVER WRITES CODE**: Orchestrator coordinates only, spawns agents for ALL work

## 📊 PROJECT_DONE CRITERIA - SF 3.0

Initialization is complete when orchestrator-state-v3.json shows:
```json
{
  "state_machine": {
    "current_state": "WAVE_START"
  }
}
```

**And all planning artifacts exist:**
- ✅ MASTER-ARCHITECTURE-PLAN.md (conceptual architecture)
- ✅ PROJECT-TEST-PLAN.md (overall testing strategy)
- ✅ PHASE-1-ARCHITECTURE-PLAN.md (pseudocode-level technical plan)
- ✅ PHASE-1-TEST-PLAN.md (phase-specific test cases)
- ✅ PHASE-1-IMPLEMENTATION-PLAN.md (detailed effort breakdown)
- ✅ project-integration branch created (R342)
- ✅ phase-1-integration branch created (R342)
- ✅ All planning info loaded into orchestrator-state-v3.json
- ✅ Git hooks installed and active (R506)

**Ready for**: `/continue-software-factory` to begin Wave 1 execution

## 🚦 STATE TRANSITIONS - SF 3.0 MANDATORY SEQUENCE

**From state-machines/software-factory-3.0-state-machine.json → mandatory_sequences.project_initialization**

```
INIT (orchestrator: verify environment, load/create state files)
    ↓ [Consult State Manager]
SPAWN_ARCHITECT_MASTER_PLANNING (orchestrator spawns architect)
    ↓ [Architect: INIT_REQUIREMENTS_GATHERING → INIT_DECOMPOSE_PRD → INIT_SYNTHESIZE_PLAN]
    ↓ [Architect creates MASTER-ARCHITECTURE-PLAN.md]
    ↓
WAITING_FOR_MASTER_ARCHITECTURE (orchestrator: verify plan exists, load into state)
    ↓ [Consult State Manager]
SPAWN_CODE_REVIEWER_PROJECT_TEST_PLANNING (orchestrator spawns code reviewer)
    ↓ [Code reviewer creates PROJECT-TEST-PLAN.md per R341]
    ↓
WAITING_FOR_PROJECT_TEST_PLAN (orchestrator: verify test plan, load into state)
    ↓ [Consult State Manager]
CREATE_PROJECT_INTEGRATION_BRANCH_EARLY (orchestrator: create integration branch per R342)
    ↓ [Consult State Manager]
SPAWN_ARCHITECT_PHASE_PLANNING (orchestrator spawns architect for Phase 1)
    ↓ [Architect creates PHASE-1-ARCHITECTURE-PLAN.md with pseudocode fidelity]
    ↓
WAITING_FOR_PHASE_ARCHITECTURE (orchestrator: verify phase plan, load into state)
    ↓ [Consult State Manager]
SPAWN_CODE_REVIEWER_PHASE_TEST_PLANNING (orchestrator spawns code reviewer)
    ↓ [Code reviewer creates PHASE-1-TEST-PLAN.md per R341]
    ↓
WAITING_FOR_PHASE_TEST_PLAN (orchestrator: verify phase test plan, load into state)
    ↓ [Consult State Manager]
CREATE_PHASE_INTEGRATION_BRANCH_EARLY (orchestrator: create phase integration branch per R342)
    ↓ [Consult State Manager]
SPAWN_CODE_REVIEWER_PHASE_IMPL (orchestrator spawns code reviewer)
    ↓ [Code reviewer creates PHASE-1-IMPLEMENTATION-PLAN.md with effort breakdown]
    ↓
WAITING_FOR_PHASE_IMPLEMENTATION_PLAN (orchestrator: verify implementation plan, load into state)
    ↓ [Consult State Manager]
WAVE_START (ready for development - user runs /continue-software-factory)
```

**SF 3.0 KEY INNOVATIONS**:
- **3-Level Planning**: Master architecture → Phase architecture → Wave architecture (fidelity gradient)
- **Test-First**: Test planning at each level before implementation (R341, R342)
- **Early Integration**: Integration branches created before implementation begins (R342)
- **State Manager Bookends**: Every transition validated by State Manager (SF 3.0 fundamental law)

## ⚠️ COMMON PITFALLS TO AVOID

- ❌ Don't skip requirements gathering
- ❌ Don't create generic plans without user input
- ❌ Don't over-customize agent files
- ❌ Don't proceed without validating all files
- ❌ Don't forget to update state file

## 🎬 BEGIN INITIALIZATION

**HAPPY PATH WORKFLOW**:

1. **User provides PRD** (via command argument or interactively)
2. **Architect asks clarifications** (only what's missing from PRD)
3. **Architect decomposes** (phases → waves → efforts with sizing)
4. **Architect synthesizes** (complete IMPLEMENTATION-PLAN.md)
5. **User reviews plan** (can request changes)
6. **System finalizes** (config files, state files, git commit)
7. **Ready to develop** (run /continue-software-factory)

**KEY INNOVATION**: The system does the hard work of breaking down your PRD into phases, waves, and efforts. You provide the "what", the system figures out the "how" (structure).

**Start by**:
- If description provided as argument: Begin processing
- If no argument: Ask user for project name and description (can be full PRD)