# 🔴🔴🔴 RULE R345: Planning Branch Requirements (SUPREME LAW) 🔴🔴🔴

## Rule Identifier
**Rule ID**: R345
**Category**: Branch Organization
**Criticality**: SUPREME LAW - Violations cause planning chaos and lost architecture
**Priority**: ABSOLUTE

## Description

Project, phase, and wave planning MUST be done in the planning repository ($CLAUDE_PROJECT_DIR/planning/) with proper directory hierarchy. ALL planning documents MUST be stored in the appropriate planning directories. Planning locations MUST be tracked in orchestrator-state-v3.json.

## Rationale

Without proper planning directory structure:
- Architecture plans get mixed with implementation
- Test plans have no proper home
- Planning artifacts pollute effort branches
- Recovery of planning decisions is impossible
- Repository boundaries get violated

## Requirements

### 1. PLANNING DIRECTORY HIERARCHY

**MANDATORY planning directory structure:**

```bash
# PLANNING REPOSITORY STRUCTURE
$CLAUDE_PROJECT_DIR/planning/
├── project/                             # PROJECT LEVEL PLANNING
│   ├── PROJECT-ARCHITECTURE-PLAN.md    # Overall architecture
│   ├── PROJECT-IMPLEMENTATION-PLAN.md  # Master implementation plan
│   ├── PROJECT-TEST-STRATEGY.md       # Test strategy
│   ├── PROJECT-DEMO-PLAN.md           # Demo plan
│   ├── PROJECT-PHASES.md              # Phase breakdown
│   └── project-metadata.yaml         # Project metadata
│
├── phase1/                             # PHASE 1 PLANNING
│   ├── PHASE-1-ARCHITECTURE-PLAN.md   # Phase architecture
│   ├── PHASE-1-IMPLEMENTATION-PLAN.md # Phase implementation
│   ├── PHASE-1-TEST-PLAN.md          # Phase tests
│   ├── PHASE-1-DEMO-PLAN.md          # Phase demo
│   ├── PHASE-1-WAVES.md               # Wave breakdown
│   ├── phase-metadata.yaml           # Phase metadata
│   │
│   ├── wave1/                         # WAVE 1-1 PLANNING
│   │   ├── WAVE-1-1-ARCHITECTURE-PLAN.md
│   │   ├── WAVE-1-1-IMPLEMENTATION-PLAN.md
│   │   ├── WAVE-1-1-TEST-PLAN.md
│   │   ├── WAVE-1-1-DEMO-PLAN.md
│   │   ├── WAVE-1-1-EFFORTS.md
│   │   └── wave-metadata.yaml
│   │
│   └── wave2/                         # WAVE 1-2 PLANNING
│       └── ...
│
└── phase2/                             # PHASE 2 PLANNING
    └── ...
```

### 2. CREATION PROTOCOL

#### Project Planning Directory:
```bash
# Create at PROJECT_START
cd "$CLAUDE_PROJECT_DIR"

# Create planning directory structure
mkdir -p planning/project

# Create planning documents
echo "# Project Architecture Plan" > planning/project/PROJECT-ARCHITECTURE-PLAN.md
echo "# Project Test Strategy" > planning/project/PROJECT-TEST-STRATEGY.md
echo "# Project Implementation Plan" > planning/project/PROJECT-IMPLEMENTATION-PLAN.md

# Update state file
yq -i ".metadata_locations.project_planning = {
  \"directory\": \"$CLAUDE_PROJECT_DIR/planning/project\",
  \"architecture\": \"$CLAUDE_PROJECT_DIR/planning/project/PROJECT-ARCHITECTURE-PLAN.md\",
  \"test_plan\": \"$CLAUDE_PROJECT_DIR/planning/project/PROJECT-TEST-STRATEGY.md\",
  \"implementation_plan\": \"$CLAUDE_PROJECT_DIR/planning/project/PROJECT-IMPLEMENTATION-PLAN.md\",
  \"created_at\": \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\"
}" orchestrator-state-v3.json

# Commit to planning repository
git add planning/project/
git add orchestrator-state-v3.json
git commit -m "planning: initialize project planning documents"
git push
```

#### Phase Planning Directory:
```bash
# Create at START_PHASE_ITERATION
cd "$CLAUDE_PROJECT_DIR"

# Create phase planning directory
mkdir -p planning/phase${PHASE}

# Create phase planning documents
echo "# Phase ${PHASE} Architecture Plan" > planning/phase${PHASE}/PHASE-${PHASE}-ARCHITECTURE-PLAN.md
echo "# Phase ${PHASE} Test Plan" > planning/phase${PHASE}/PHASE-${PHASE}-TEST-PLAN.md
echo "# Phase ${PHASE} Implementation Plan" > planning/phase${PHASE}/PHASE-${PHASE}-IMPLEMENTATION-PLAN.md

# Update state file
yq -i ".metadata_locations.phase_planning.phase${PHASE} = {
  \"directory\": \"$CLAUDE_PROJECT_DIR/planning/phase${PHASE}\",
  \"architecture\": \"$CLAUDE_PROJECT_DIR/planning/phase${PHASE}/PHASE-${PHASE}-ARCHITECTURE-PLAN.md\",
  \"test_plan\": \"$CLAUDE_PROJECT_DIR/planning/phase${PHASE}/PHASE-${PHASE}-TEST-PLAN.md\",
  \"implementation_plan\": \"$CLAUDE_PROJECT_DIR/planning/phase${PHASE}/PHASE-${PHASE}-IMPLEMENTATION-PLAN.md\",
  \"created_at\": \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\"
}" orchestrator-state-v3.json

# Commit to planning repository
git add planning/phase${PHASE}/
git add orchestrator-state-v3.json
git commit -m "planning: initialize phase ${PHASE} planning documents"
git push
```

#### Wave Planning Directory:
```bash
# Create at WAVE_START
cd "$CLAUDE_PROJECT_DIR"

# Create wave planning directory
mkdir -p planning/phase${PHASE}/wave${WAVE}

# Create wave planning documents
echo "# Wave ${PHASE}-${WAVE} Architecture Plan" > planning/phase${PHASE}/wave${WAVE}/WAVE-${PHASE}-${WAVE}-ARCHITECTURE-PLAN.md
echo "# Wave ${PHASE}-${WAVE} Test Plan" > planning/phase${PHASE}/wave${WAVE}/WAVE-${PHASE}-${WAVE}-TEST-PLAN.md
echo "# Wave ${PHASE}-${WAVE} Implementation Plan" > planning/phase${PHASE}/wave${WAVE}/WAVE-${PHASE}-${WAVE}-IMPLEMENTATION-PLAN.md

# Update state file
yq -i ".metadata_locations.wave_planning.phase${PHASE}_wave${WAVE} = {
  \"directory\": \"$CLAUDE_PROJECT_DIR/planning/phase${PHASE}/wave${WAVE}\",
  \"architecture\": \"$CLAUDE_PROJECT_DIR/planning/phase${PHASE}/wave${WAVE}/WAVE-${PHASE}-${WAVE}-ARCHITECTURE-PLAN.md\",
  \"test_plan\": \"$CLAUDE_PROJECT_DIR/planning/phase${PHASE}/wave${WAVE}/WAVE-${PHASE}-${WAVE}-TEST-PLAN.md\",
  \"implementation_plan\": \"$CLAUDE_PROJECT_DIR/planning/phase${PHASE}/wave${WAVE}/WAVE-${PHASE}-${WAVE}-IMPLEMENTATION-PLAN.md\",
  \"created_at\": \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\"
}" orchestrator-state-v3.json

# Commit to planning repository
git add planning/phase${PHASE}/wave${WAVE}/
git add orchestrator-state-v3.json
git commit -m "planning: initialize phase ${PHASE} wave ${WAVE} planning documents"
git push
```

### 3. STATE FILE TRACKING

**orchestrator-state-v3.json MUST track all planning locations:**

```json
{
  "metadata_locations": {
    "project_planning": {
      "directory": "$CLAUDE_PROJECT_DIR/planning/project",
      "architecture": "$CLAUDE_PROJECT_DIR/planning/project/PROJECT-ARCHITECTURE-PLAN.md",
      "implementation_plan": "$CLAUDE_PROJECT_DIR/planning/project/PROJECT-IMPLEMENTATION-PLAN.md",
      "test_strategy": "$CLAUDE_PROJECT_DIR/planning/project/PROJECT-TEST-STRATEGY.md",
      "demo_plan": "$CLAUDE_PROJECT_DIR/planning/project/PROJECT-DEMO-PLAN.md",
      "phase_breakdown": "$CLAUDE_PROJECT_DIR/planning/project/PROJECT-PHASES.md",
      "created_at": "2025-01-20T09:00:00Z",
      "created_by": "orchestrator",
      "status": "active"
    },
    "phase_planning": {
      "phase1": {
        "directory": "$CLAUDE_PROJECT_DIR/planning/phase1",
        "architecture": "$CLAUDE_PROJECT_DIR/planning/phase1/PHASE-1-ARCHITECTURE-PLAN.md",
        "implementation_plan": "$CLAUDE_PROJECT_DIR/planning/phase1/PHASE-1-IMPLEMENTATION-PLAN.md",
        "test_plan": "$CLAUDE_PROJECT_DIR/planning/phase1/PHASE-1-TEST-PLAN.md",
        "demo_plan": "$CLAUDE_PROJECT_DIR/planning/phase1/PHASE-1-DEMO-PLAN.md",
        "wave_breakdown": "$CLAUDE_PROJECT_DIR/planning/phase1/PHASE-1-WAVES.md",
        "created_at": "2025-01-20T10:00:00Z",
        "created_by": "orchestrator",
        "status": "active"
      }
    },
    "wave_planning": {
      "phase1_wave1": {
        "directory": "$CLAUDE_PROJECT_DIR/planning/phase1/wave1",
        "architecture": "$CLAUDE_PROJECT_DIR/planning/phase1/wave1/WAVE-1-1-ARCHITECTURE-PLAN.md",
        "implementation_plan": "$CLAUDE_PROJECT_DIR/planning/phase1/wave1/WAVE-1-1-IMPLEMENTATION-PLAN.md",
        "test_plan": "$CLAUDE_PROJECT_DIR/planning/phase1/wave1/WAVE-1-1-TEST-PLAN.md",
        "demo_plan": "$CLAUDE_PROJECT_DIR/planning/phase1/wave1/WAVE-1-1-DEMO-PLAN.md",
        "effort_breakdown": "$CLAUDE_PROJECT_DIR/planning/phase1/wave1/WAVE-1-1-EFFORTS.md",
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
- MUST create planning directory structure at each level start
- MUST track planning locations in orchestrator-state-v3.json
- MUST instruct Architect/Code Reviewer to use planning directories
- MUST verify planning directories exist before spawning agents

#### Architect:
- MUST use planning/phaseN/ directory for phase architecture
- MUST use planning/phaseN/waveM/ directory for wave architecture
- MUST store all plans in appropriate planning directories
- MUST update state file with document locations

#### Code Reviewer:
- MUST use planning/phaseN/ directory for phase-level plans
- MUST use planning/phaseN/waveM/ directory for wave-level plans
- MUST create effort-specific plans in effort .software-factory/ directories
- MUST update state file with all plan locations

### 5. PLANNING CONTENT REQUIREMENTS

#### Project Planning (planning/project/):
- PROJECT-ARCHITECTURE-PLAN.md - Overall system design
- PROJECT-IMPLEMENTATION-PLAN.md - Master implementation plan
- PROJECT-TEST-STRATEGY.md - End-to-end test strategy
- PROJECT-DEMO-PLAN.md - Final demonstration approach
- PROJECT-PHASES.md - Phase breakdown and goals
- project-metadata.yaml - Project tracking data

#### Phase Planning (planning/phaseN/):
- PHASE-N-ARCHITECTURE-PLAN.md - Phase-specific design
- PHASE-N-IMPLEMENTATION-PLAN.md - Phase implementation strategy
- PHASE-N-TEST-PLAN.md - Phase integration tests
- PHASE-N-DEMO-PLAN.md - Phase demonstration
- PHASE-N-WAVES.md - Wave breakdown and dependencies
- phase-metadata.yaml - Phase tracking data

#### Wave Planning (planning/phaseN/waveM/):
- WAVE-N-M-ARCHITECTURE-PLAN.md - Wave implementation design
- WAVE-N-M-IMPLEMENTATION-PLAN.md - Wave implementation details
- WAVE-N-M-TEST-PLAN.md - Wave-specific tests
- WAVE-N-M-DEMO-PLAN.md - Wave demonstration
- WAVE-N-M-EFFORTS.md - Effort breakdown and sizing
- wave-metadata.yaml - Wave tracking data

### 6. DIRECTORY NAMING CONVENTIONS

```bash
# Project Planning
$CLAUDE_PROJECT_DIR/planning/project/

# Phase Planning
$CLAUDE_PROJECT_DIR/planning/phase${N}/

# Wave Planning
$CLAUDE_PROJECT_DIR/planning/phase${N}/wave${M}/

# Examples:
/home/user/sf-instance/planning/project/
/home/user/sf-instance/planning/phase1/
/home/user/sf-instance/planning/phase1/wave1/
```

### 7. VALIDATION PROTOCOL

```bash
validate_planning_structure() {
    local phase="$1"
    local wave="$2"

    # Check project planning directory
    if [ ! -d "$CLAUDE_PROJECT_DIR/planning/project" ]; then
        echo "❌ R345 VIOLATION: Missing project planning directory"
        return 1
    fi

    # Check phase planning directory
    if [ ! -d "$CLAUDE_PROJECT_DIR/planning/phase${phase}" ]; then
        echo "❌ R345 VIOLATION: Missing phase${phase} planning directory"
        return 1
    fi

    # Check wave planning directory
    if [ ! -d "$CLAUDE_PROJECT_DIR/planning/phase${phase}/wave${wave}" ]; then
        echo "❌ R345 VIOLATION: Missing phase${phase}/wave${wave} planning directory"
        return 1
    fi

    # Check state file tracking
    local project_dir=$(yq '.metadata_locations.project_planning.directory' orchestrator-state-v3.json)
    if [ "$project_dir" = "null" ]; then
        echo "❌ R345 VIOLATION: Project planning not tracked in state"
        return 1
    fi

    echo "✅ R345 COMPLIANT: All planning directories present and tracked"
    return 0
}
```

## Common Violations

1. ❌ **Putting planning documents in effort directories**
   - Wrong: Phase architecture in /efforts/phase1/
   - Right: Phase architecture in $CLAUDE_PROJECT_DIR/planning/phase1/

2. ❌ **Not creating planning directory structure**
   - Wrong: Skip planning directory creation
   - Right: Create planning directories at each level start

3. ❌ **Not tracking planning locations in state**
   - Wrong: Create directories without updating state
   - Right: Update state immediately after directory creation

4. ❌ **Mixing planning repository with effort repository**
   - Wrong: Planning files in effort branches
   - Right: Planning files in planning repository only

## Enforcement

- **Grading Impact**: -50% for missing planning directories
- **Architecture Loss**: Cannot recover planning without proper structure
- **Repository Violation**: -100% for mixing planning and effort repositories

## Related Rules

- R343: Metadata directory standardization
- R344: Metadata location tracking
- R341: TDD test plan requirements
- R342: Early integration branch creation

## Implementation Notes

1. Planning directories created BEFORE implementation starts
2. All planning in $CLAUDE_PROJECT_DIR/planning/ hierarchy
3. Planning repository separate from effort repository
4. Enable clean separation of planning and code
5. Git hooks prevent cross-repository contamination

## Key Principle

**"Planning deserves its own repository structure - never mix with implementation"**

This ensures:
- Clean separation of repositories
- Persistent planning artifacts
- Clear directory hierarchy
- Easy planning recovery
- No cross-repository pollution

---

**Remember**: Every level (project/phase/wave) gets a planning directory in the planning repository. NO EXCEPTIONS.