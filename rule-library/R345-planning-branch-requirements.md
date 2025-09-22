# 🔴🔴🔴 RULE R345: Planning Branch Requirements (SUPREME LAW) 🔴🔴🔴

## Rule Identifier
**Rule ID**: R345
**Category**: Branch Organization
**Criticality**: SUPREME LAW - Violations cause planning chaos and lost architecture
**Priority**: ABSOLUTE

## Description

Project, phase, and wave planning MUST have dedicated branches and working copies. ALL planning metadata MUST be stored in .software-factory folders on these branches. Planning branches MUST be tracked in orchestrator-state.json.

## Rationale

Without dedicated planning branches:
- Architecture plans get mixed with implementation
- Test plans have no proper home
- Planning artifacts pollute main branch
- Recovery of planning decisions is impossible
- Parallel planning efforts conflict

## Requirements

### 1. PLANNING BRANCH HIERARCHY

**MANDATORY planning branches at each level:**

```bash
# PROJECT LEVEL PLANNING
Branch: ${PROJECT_PREFIX}/project-planning
Working Copy: /efforts/project-planning/
Content:
├── .software-factory/
│   ├── PROJECT-ARCHITECTURE.md
│   ├── PROJECT-TEST-PLAN.md
│   ├── PROJECT-DEMO-PLAN.md
│   ├── PROJECT-PHASES.md
│   └── project-metadata.yaml
└── [project planning artifacts]

# PHASE LEVEL PLANNING
Branch: ${PROJECT_PREFIX}/phase${N}/phase-planning
Working Copy: /efforts/phase${N}/phase-planning/
Content:
├── .software-factory/
│   ├── PHASE-ARCHITECTURE.md
│   ├── PHASE-TEST-PLAN.md
│   ├── PHASE-DEMO-PLAN.md
│   ├── PHASE-WAVES.md
│   └── phase-metadata.yaml
└── [phase planning artifacts]

# WAVE LEVEL PLANNING  
Branch: ${PROJECT_PREFIX}/phase${N}/wave${M}/wave-planning
Working Copy: /efforts/phase${N}/wave${M}/wave-planning/
Content:
├── .software-factory/
│   ├── WAVE-ARCHITECTURE.md
│   ├── WAVE-TEST-PLAN.md
│   ├── WAVE-DEMO-PLAN.md
│   ├── WAVE-EFFORTS.md
│   └── wave-metadata.yaml
└── [wave planning artifacts]
```

### 2. CREATION PROTOCOL

#### Project Planning Branch:
```bash
# Create at PROJECT_START
cd /efforts
git clone "$TARGET_REPO_URL" project-planning
cd project-planning
git checkout -b "${PROJECT_PREFIX}/project-planning"

# Create metadata structure
mkdir -p .software-factory
echo "# Project Architecture" > .software-factory/PROJECT-ARCHITECTURE.md
echo "# Project Test Plan" > .software-factory/PROJECT-TEST-PLAN.md

# Update state file
yq -i ".metadata_locations.project_planning = {
  \"branch\": \"${PROJECT_PREFIX}/project-planning\",
  \"working_copy\": \"/efforts/project-planning\",
  \"architecture\": \"/efforts/project-planning/.software-factory/PROJECT-ARCHITECTURE.md\",
  \"test_plan\": \"/efforts/project-planning/.software-factory/PROJECT-TEST-PLAN.md\",
  \"created_at\": \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\"
}" orchestrator-state.json

# Commit and push
git add -A
git commit -m "init: project planning branch with metadata structure"
git push -u origin "${PROJECT_PREFIX}/project-planning"
```

#### Phase Planning Branch:
```bash
# Create at PHASE_START
cd /efforts/phase${PHASE}
git clone "$TARGET_REPO_URL" phase-planning
cd phase-planning
git checkout -b "${PROJECT_PREFIX}/phase${PHASE}/phase-planning"

# Create metadata structure
mkdir -p .software-factory
echo "# Phase ${PHASE} Architecture" > .software-factory/PHASE-ARCHITECTURE.md
echo "# Phase ${PHASE} Test Plan" > .software-factory/PHASE-TEST-PLAN.md

# Update state file
yq -i ".metadata_locations.phase_planning.phase${PHASE} = {
  \"branch\": \"${PROJECT_PREFIX}/phase${PHASE}/phase-planning\",
  \"working_copy\": \"/efforts/phase${PHASE}/phase-planning\",
  \"architecture\": \"/efforts/phase${PHASE}/phase-planning/.software-factory/PHASE-ARCHITECTURE.md\",
  \"test_plan\": \"/efforts/phase${PHASE}/phase-planning/.software-factory/PHASE-TEST-PLAN.md\",
  \"created_at\": \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\"
}" orchestrator-state.json
```

#### Wave Planning Branch:
```bash
# Create at WAVE_START
cd /efforts/phase${PHASE}/wave${WAVE}
git clone "$TARGET_REPO_URL" wave-planning
cd wave-planning
git checkout -b "${PROJECT_PREFIX}/phase${PHASE}/wave${WAVE}/wave-planning"

# Create metadata structure
mkdir -p .software-factory
echo "# Wave ${WAVE} Architecture" > .software-factory/WAVE-ARCHITECTURE.md
echo "# Wave ${WAVE} Test Plan" > .software-factory/WAVE-TEST-PLAN.md

# Update state file
yq -i ".metadata_locations.wave_planning.phase${PHASE}_wave${WAVE} = {
  \"branch\": \"${PROJECT_PREFIX}/phase${PHASE}/wave${WAVE}/wave-planning\",
  \"working_copy\": \"/efforts/phase${PHASE}/wave${WAVE}/wave-planning\",
  \"architecture\": \"/efforts/phase${PHASE}/wave${WAVE}/wave-planning/.software-factory/WAVE-ARCHITECTURE.md\",
  \"test_plan\": \"/efforts/phase${PHASE}/wave${WAVE}/wave-planning/.software-factory/WAVE-TEST-PLAN.md\",
  \"created_at\": \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\"
}" orchestrator-state.json
```

### 3. STATE FILE TRACKING

**orchestrator-state.json MUST track all planning branches:**

```json
{
  "metadata_locations": {
    "project_planning": {
      "branch": "myproject/project-planning",
      "working_copy": "/efforts/project-planning",
      "architecture": "/efforts/project-planning/.software-factory/PROJECT-ARCHITECTURE.md",
      "test_plan": "/efforts/project-planning/.software-factory/PROJECT-TEST-PLAN.md",
      "demo_plan": "/efforts/project-planning/.software-factory/PROJECT-DEMO-PLAN.md",
      "phase_breakdown": "/efforts/project-planning/.software-factory/PROJECT-PHASES.md",
      "created_at": "2025-01-20T09:00:00Z",
      "created_by": "orchestrator",
      "status": "active"
    },
    "phase_planning": {
      "phase1": {
        "branch": "myproject/phase1/phase-planning",
        "working_copy": "/efforts/phase1/phase-planning",
        "architecture": "/efforts/phase1/phase-planning/.software-factory/PHASE-ARCHITECTURE.md",
        "test_plan": "/efforts/phase1/phase-planning/.software-factory/PHASE-TEST-PLAN.md",
        "demo_plan": "/efforts/phase1/phase-planning/.software-factory/PHASE-DEMO-PLAN.md",
        "wave_breakdown": "/efforts/phase1/phase-planning/.software-factory/PHASE-WAVES.md",
        "created_at": "2025-01-20T10:00:00Z",
        "created_by": "orchestrator",
        "status": "active"
      }
    },
    "wave_planning": {
      "phase1_wave1": {
        "branch": "myproject/phase1/wave1/wave-planning",
        "working_copy": "/efforts/phase1/wave1/wave-planning",
        "architecture": "/efforts/phase1/wave1/wave-planning/.software-factory/WAVE-ARCHITECTURE.md",
        "test_plan": "/efforts/phase1/wave1/wave-planning/.software-factory/WAVE-TEST-PLAN.md",
        "demo_plan": "/efforts/phase1/wave1/wave-planning/.software-factory/WAVE-DEMO-PLAN.md",
        "effort_breakdown": "/efforts/phase1/wave1/wave-planning/.software-factory/WAVE-EFFORTS.md",
        "created_at": "2025-01-20T11:00:00Z",
        "created_by": "orchestrator",
        "status": "active"
      }
    }
  }
}
```

### 4. AGENT RESPONSIBILITIES

#### Orchestrator:
- MUST create planning branches at each level start
- MUST track branches in orchestrator-state.json
- MUST instruct Architect/Code Reviewer to use planning branches
- MUST verify planning branches exist before spawning agents

#### Architect:
- MUST use phase-planning branch for phase architecture
- MUST use wave-planning branch for wave architecture  
- MUST store all plans in .software-factory/
- MUST update state file with document locations

#### Code Reviewer:
- MUST use phase-planning branch for phase test plans
- MUST use wave-planning branch for wave test plans
- MUST store all plans in .software-factory/
- MUST update state file with plan locations

### 5. PLANNING CONTENT REQUIREMENTS

#### Project Planning (.software-factory/):
- PROJECT-ARCHITECTURE.md - Overall system design
- PROJECT-TEST-PLAN.md - End-to-end test strategy
- PROJECT-DEMO-PLAN.md - Final demonstration approach
- PROJECT-PHASES.md - Phase breakdown and goals
- project-metadata.yaml - Project tracking data

#### Phase Planning (.software-factory/):
- PHASE-ARCHITECTURE.md - Phase-specific design
- PHASE-TEST-PLAN.md - Phase integration tests
- PHASE-DEMO-PLAN.md - Phase demonstration
- PHASE-WAVES.md - Wave breakdown and dependencies
- phase-metadata.yaml - Phase tracking data

#### Wave Planning (.software-factory/):
- WAVE-ARCHITECTURE.md - Wave implementation design
- WAVE-TEST-PLAN.md - Wave-specific tests
- WAVE-DEMO-PLAN.md - Wave demonstration
- WAVE-EFFORTS.md - Effort breakdown and sizing
- wave-metadata.yaml - Wave tracking data

### 6. BRANCH NAMING CONVENTIONS

```bash
# Project Planning
${PROJECT_PREFIX}/project-planning

# Phase Planning
${PROJECT_PREFIX}/phase${N}/phase-planning

# Wave Planning  
${PROJECT_PREFIX}/phase${N}/wave${M}/wave-planning

# Examples:
idpbuilder-oci/project-planning
idpbuilder-oci/phase1/phase-planning
idpbuilder-oci/phase1/wave1/wave-planning
```

### 7. VALIDATION PROTOCOL

```bash
validate_planning_branches() {
    local phase="$1"
    local wave="$2"
    
    # Check project planning
    local project_branch=$(yq '.metadata_locations.project_planning.branch' orchestrator-state.json)
    if [ "$project_branch" = "null" ]; then
        echo "❌ R345 VIOLATION: Missing project planning branch"
        return 1
    fi
    
    # Check phase planning
    local phase_branch=$(yq ".metadata_locations.phase_planning.phase${phase}.branch" orchestrator-state.json)
    if [ "$phase_branch" = "null" ]; then
        echo "❌ R345 VIOLATION: Missing phase${phase} planning branch"
        return 1
    fi
    
    # Check wave planning
    local wave_branch=$(yq ".metadata_locations.wave_planning.phase${phase}_wave${wave}.branch" orchestrator-state.json)
    if [ "$wave_branch" = "null" ]; then
        echo "❌ R345 VIOLATION: Missing phase${phase}/wave${wave} planning branch"
        return 1
    fi
    
    echo "✅ R345 COMPLIANT: All planning branches present"
    return 0
}
```

## Common Violations

1. ❌ **Putting planning documents in implementation branches**
   - Wrong: Architecture in effort branch
   - Right: Architecture in planning branch

2. ❌ **Not creating planning branches**
   - Wrong: Skip planning branch creation
   - Right: Create planning branch at each level start

3. ❌ **Not tracking planning branches in state**
   - Wrong: Create branch without updating state
   - Right: Update state immediately after branch creation

4. ❌ **Mixing planning with code**
   - Wrong: Test plans in src/ directory
   - Right: Test plans in planning/.software-factory/

## Enforcement

- **Grading Impact**: -50% for missing planning branches
- **Architecture Loss**: Cannot recover planning without branches
- **Test Failure**: Tests have no home without planning branches

## Related Rules

- R343: Metadata directory standardization
- R344: Metadata location tracking
- R341: TDD test plan requirements
- R342: Early integration branch creation

## Implementation Notes

1. Planning branches created BEFORE implementation starts
2. All planning in .software-factory/ directories
3. Branches persist throughout project lifecycle
4. Enable clean separation of planning and code
5. Support parallel planning efforts

## Key Principle

**"Planning deserves its own branches - never mix with implementation"**

This ensures:
- Clean separation of concerns
- Persistent planning artifacts
- Parallel planning capability
- Easy planning recovery
- No main branch pollution

---

**Remember**: Every level (project/phase/wave) gets a planning branch. NO EXCEPTIONS.