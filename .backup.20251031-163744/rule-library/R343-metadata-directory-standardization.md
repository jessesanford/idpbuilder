# 🔴🔴🔴 RULE R343: Metadata Directory Standardization (SUPREME LAW) 🔴🔴🔴

## Rule Identifier
**Rule ID**: R343
**Category**: Infrastructure Consistency
**Criticality**: SUPREME LAW - Violations cause agent confusion and lost work
**Priority**: ABSOLUTE

## Description

ALL effort-specific metadata, reports, and non-code artifacts MUST be stored in `.software-factory/` directories within effort workspaces. Project/phase/wave planning documents MUST be stored in the planning repository directory hierarchy (`$CLAUDE_PROJECT_DIR/planning/`). This standardizes metadata organization and maintains clean repository boundaries.

## Rationale

Without consistent metadata organization:
- Agents cannot find planning documents
- Code and metadata get mixed together
- Git operations become confusing
- Integration loses track of metadata
- Recovery after failures is impossible

## Requirements

### 1. ABSOLUTE METADATA CONTAINMENT

**NO metadata may exist outside .software-factory/ directories - ZERO EXCEPTIONS:**

- ❌ FORBIDDEN: Any .md files with metadata content in root directories
- ❌ FORBIDDEN: Work logs, reports, or plans anywhere except .software-factory/
- ❌ FORBIDDEN: Test results, validation output outside .software-factory/
- ❌ FORBIDDEN: Architecture documents, demo plans outside .software-factory/
- ✅ REQUIRED: ALL metadata MUST be in .software-factory/

### 2. METADATA DIRECTORY PATTERNS

**Planning Repository uses hierarchical directories:**

```bash
# PLANNING REPOSITORY (Software Factory Instance)
$CLAUDE_PROJECT_DIR/planning/
├── project/                              # Project-level planning
│   ├── PROJECT-ARCHITECTURE-PLAN.md
│   ├── PROJECT-IMPLEMENTATION-PLAN.md
│   ├── PROJECT-TEST-STRATEGY.md
│   ├── PROJECT-DEMO-PLAN.md
│   ├── PROJECT-PHASES.md
│   └── project-metadata.yaml
│
├── phase1/                               # Phase 1 planning
│   ├── PHASE-1-ARCHITECTURE-PLAN.md
│   ├── PHASE-1-IMPLEMENTATION-PLAN.md
│   ├── PHASE-1-TEST-PLAN.md
│   ├── PHASE-1-DEMO-PLAN.md
│   ├── PHASE-1-WAVES.md
│   ├── phase-metadata.yaml
│   │
│   ├── wave1/                            # Wave 1-1 planning
│   │   ├── WAVE-1-1-ARCHITECTURE-PLAN.md
│   │   ├── WAVE-1-1-IMPLEMENTATION-PLAN.md
│   │   ├── WAVE-1-1-TEST-PLAN.md
│   │   ├── WAVE-1-1-DEMO-PLAN.md
│   │   ├── WAVE-1-1-EFFORTS.md
│   │   └── wave-metadata.yaml
│   │
│   └── wave2/                            # Wave 1-2 planning
│       └── ...
└── phase2/                               # Phase 2 planning
    └── ...

# Effort Workspace
/efforts/phase1/wave1/effort-name/
├── .software-factory/                    # ALL metadata here
│   ├── IMPLEMENTATION-PLAN.md
│   ├── work-log.md
│   ├── CODE-REVIEW-REPORT.md
│   └── validation-results.md
├── src/                                  # Implementation code
├── tests/                                # Test code
└── README.md                            # Project files

# Split Workspace
/efforts/phase1/wave1/effort-name-split-001/
├── .software-factory/                    # Split metadata
│   ├── SPLIT-PLAN.md
│   ├── work-log.md
│   └── review-report.md
└── [implementation files]

# Integration Workspace
/efforts/phase1/wave1/integration-workspace/
├── .software-factory/                    # Integration metadata
│   ├── WAVE-MERGE-PLAN.md
│   ├── INTEGRATE_WAVE_EFFORTS-REPORT.md
│   ├── INTEGRATE_WAVE_EFFORTS-METADATA.md
│   └── integration-log.md
├── src/                                  # Merged code
└── tests/                                # Integration tests
```

### 2. FORBIDDEN PATTERNS

**NEVER place metadata in these locations:**
- ❌ Root directory of effort workspace (e.g., `/effort-name/PLAN.md`)
- ❌ Mixed with code directories (e.g., `/src/PLAN.md`)
- ❌ In a `/repo/` subdirectory (deprecated pattern)
- ❌ Planning files in effort directories (e.g., `/efforts/phase1/PHASE-PLAN.md`)
- ❌ Effort files in planning repository (e.g., `/planning/EFFORT-PLAN.md`)
- ❌ Scattered across multiple locations

### 3. COMPREHENSIVE METADATA LIST

**Effort Repository Files (in .software-factory/):**
- Effort Implementation Plans (EFFORT-IMPLEMENTATION-PLAN--*.md)
- Code Review Reports (CODE-REVIEW-REPORT--*.md)
- Split Plans (SPLIT-PLAN--*.md)
- Fix Plans (FIX-PLAN--*.md)
- Work Logs (work-log.md)
- Validation Results (validation-results.md)
- Test Results (test-results.md)
- Effort metadata files (*.yaml)

**Planning Repository Files (in planning/ hierarchy):**
- Project Architecture (PROJECT-ARCHITECTURE-PLAN.md)
- Phase Architecture (PHASE-N-ARCHITECTURE-PLAN.md)
- Wave Architecture (WAVE-N-M-ARCHITECTURE-PLAN.md)
- Phase/Wave Implementation Plans (for overall strategy)
- Merge Plans (WAVE-N-M-MERGE-PLAN.md, PHASE-N-MERGE-PLAN.md)
- Test Strategy Documents (PROJECT-TEST-STRATEGY.md)
- Demo Plans (PROJECT-DEMO-PLAN.md, PHASE-N-DEMO-PLAN.md)
- Integration Reports (PHASE-N-INTEGRATE_WAVE_EFFORTS-REPORT.md)

### 4. METADATA TYPES BY LOCATION

#### Planning Repository Structure
```bash
# Project Planning ($CLAUDE_PROJECT_DIR/planning/project/)
├── PROJECT-ARCHITECTURE-PLAN.md         # Overall system design
├── PROJECT-IMPLEMENTATION-PLAN.md       # Master implementation plan
├── PROJECT-TEST-STRATEGY.md            # End-to-end test strategy
├── PROJECT-DEMO-PLAN.md                # Final demonstration approach
├── PROJECT-PHASES.md                   # Phase breakdown and goals
└── project-metadata.yaml               # Project tracking data

# Phase Planning ($CLAUDE_PROJECT_DIR/planning/phaseN/)
├── PHASE-N-ARCHITECTURE-PLAN.md        # Phase-specific design
├── PHASE-N-IMPLEMENTATION-PLAN.md      # Phase implementation strategy
├── PHASE-N-TEST-PLAN.md                # Phase integration tests
├── PHASE-N-DEMO-PLAN.md                # Phase demonstration
├── PHASE-N-WAVES.md                    # Wave breakdown and dependencies
└── phase-metadata.yaml                 # Phase tracking data

# Wave Planning ($CLAUDE_PROJECT_DIR/planning/phaseN/waveM/)
├── WAVE-N-M-ARCHITECTURE-PLAN.md       # Wave implementation design
├── WAVE-N-M-IMPLEMENTATION-PLAN.md     # Wave implementation details
├── WAVE-N-M-TEST-PLAN.md               # Wave-specific tests
├── WAVE-N-M-DEMO-PLAN.md               # Wave demonstration
├── WAVE-N-M-EFFORTS.md                 # Effort breakdown and sizing
└── wave-metadata.yaml                  # Wave tracking data
```

#### Effort Workspace Metadata (.software-factory/)
```bash
/efforts/phase1/wave1/effort-name/.software-factory/
└── phase1/
    └── wave1/
        └── effort-name/
            ├── EFFORT-IMPLEMENTATION-PLAN--YYYYMMDD-HHMMSS.md  # R383 compliant
            ├── CODE-REVIEW-REPORT--YYYYMMDD-HHMMSS.md         # R383 compliant
            ├── FIX-PLAN--YYYYMMDD-HHMMSS.md                   # R383 compliant
            ├── work-log.md                                     # Development log
            ├── validation-results.md                           # Test results
            └── effort-metadata.yaml                            # Effort tracking
```

#### Split Workspace Metadata
```bash
.software-factory/
├── SPLIT-PLAN-*.md                      # Split strategy
├── SPLIT-IMPLEMENTATION-*.md            # Split details
├── split-review-report.md               # Review results
└── split-metadata.yaml                  # Split tracking
```

#### Integration Workspace Metadata
```bash
.software-factory/
├── WAVE-MERGE-PLAN.md                   # Merge strategy
├── INTEGRATE_WAVE_EFFORTS-REPORT.md                # Integration results
├── INTEGRATE_WAVE_EFFORTS-METADATA.md              # Infrastructure details
├── integration-log.md                   # Operation log
├── conflict-resolutions.md              # If conflicts occurred
└── integration-metadata.yaml            # Integration tracking
```

### 4. CREATION PROTOCOL

**For Planning Documents:**
```bash
# Create in planning repository
cd "$CLAUDE_PROJECT_DIR"
mkdir -p planning/phase${PHASE}/wave${WAVE}
echo "# Wave Architecture" > planning/phase${PHASE}/wave${WAVE}/WAVE-${PHASE}-${WAVE}-ARCHITECTURE-PLAN.md
```

**For Effort Workspaces:**
```bash
# Step 1: Clone target repository
git clone "$TARGET_REPO_URL" "$WORKSPACE_DIR"
cd "$WORKSPACE_DIR"

# Step 2: Create .software-factory directory structure (R383 compliant)
mkdir -p .software-factory/phase${PHASE}/wave${WAVE}/${EFFORT_NAME}

# Step 3: Add to .gitignore if needed
if ! grep -q "^\.software-factory" .gitignore 2>/dev/null; then
    echo "# Software Factory metadata (local only)" >> .gitignore
    echo ".software-factory/" >> .gitignore
fi

# Step 4: Place effort metadata in .software-factory with timestamps
TIMESTAMP=$(date +%Y%m%d-%H%M%S)
echo "# Work Log" > .software-factory/phase${PHASE}/wave${WAVE}/${EFFORT_NAME}/work-log.md
```

### 5. AGENT RESPONSIBILITIES

#### SW Engineer
- MUST check for `.software-factory/` on startup
- MUST read effort plans from `.software-factory/`
- MUST read wave plans from `$CLAUDE_PROJECT_DIR/planning/`
- MUST write logs to `.software-factory/`
- NEVER create planning documents

#### Code Reviewer
- MUST write effort plans to `.software-factory/` in efforts
- MUST write phase/wave plans to `$CLAUDE_PROJECT_DIR/planning/`
- MUST save reports to appropriate location
- NEVER mix planning and effort files

#### Architect
- MUST write architecture to `$CLAUDE_PROJECT_DIR/planning/`
- MUST read from planning repository
- NEVER create documents in effort directories
- MUST maintain planning hierarchy

#### Integration Agent
- MUST create `.software-factory/` for integration metadata
- MUST read merge plans from `$CLAUDE_PROJECT_DIR/planning/`
- MUST write integration reports to planning repository
- NEVER use `/repo/` subdirectory pattern

#### Orchestrator
- MUST verify proper structure in both repositories
- MUST track all locations in state file
- MUST enforce repository boundaries
- NEVER allow cross-repository contamination

### 6. MIGRATION FROM OLD PATTERNS

**If encountering old patterns:**

```bash
# Old pattern: /integration-workspace/repo/
if [ -d "integration-workspace/repo" ]; then
    echo "⚠️ WARNING: Deprecated /repo/ pattern detected"
    echo "New pattern: Clone directly as integration-workspace"
fi

# Old pattern: Metadata in root
if [ -f "IMPLEMENTATION-PLAN.md" ] && [ ! -d ".software-factory" ]; then
    mkdir -p .software-factory
    mv IMPLEMENTATION-PLAN*.md .software-factory/ 2>/dev/null
    mv CODE-REVIEW-REPORT*.md .software-factory/ 2>/dev/null
    echo "✅ Migrated metadata to .software-factory/"
fi
```

### 7. ENHANCED VALIDATION CHECKS

```bash
validate_repository_separation() {
    # Check planning repository structure
    if [ ! -d "$CLAUDE_PROJECT_DIR/planning" ]; then
        echo "❌ R343 VIOLATION: Missing planning directory"
        return 1
    fi

    # Check no effort files in planning repo
    if find "$CLAUDE_PROJECT_DIR/planning" -name "EFFORT-*.md" 2>/dev/null | grep -q .; then
        echo "❌ R343 VIOLATION: Effort files found in planning repository"
        return 1
    fi

    # Check no planning files in efforts
    if find "/efforts" -name "PHASE-*-ARCHITECTURE-PLAN.md" -o -name "WAVE-*-ARCHITECTURE-PLAN.md" 2>/dev/null | grep -q .; then
        echo "❌ R343 VIOLATION: Planning files found in effort repository"
        return 1
    fi

    echo "✅ Repository separation validated"
}

validate_effort_metadata_structure() {
    local workspace="$1"

    # Check .software-factory exists for efforts
    if [[ "$workspace" == */efforts/* ]] && [ ! -d "$workspace/.software-factory" ]; then
        echo "❌ R343 VIOLATION: Missing .software-factory directory in effort"
        return 1
    fi

    # Check no metadata in root
    if ls "$workspace"/*.md 2>/dev/null | grep -E "(IMPLEMENTATION-PLAN|REVIEW-REPORT|FIX-PLAN)" > /dev/null; then
        echo "❌ R343 VIOLATION: Effort metadata files in root directory"
        return 1
    fi

    echo "✅ R343 COMPLIANT: Proper metadata structure"
    return 0
}
```

## Common Violations

1. ❌ **Creating planning files in effort directories**
   - Wrong: `/efforts/phase1/PHASE-1-ARCHITECTURE-PLAN.md`
   - Right: `$CLAUDE_PROJECT_DIR/planning/phase1/PHASE-1-ARCHITECTURE-PLAN.md`

2. ❌ **Creating effort files in planning repository**
   - Wrong: `$CLAUDE_PROJECT_DIR/planning/EFFORT-IMPLEMENTATION-PLAN.md`
   - Right: `/efforts/phase1/wave1/effort/.software-factory/phase1/wave1/effort/EFFORT-IMPLEMENTATION-PLAN--20250120-100000.md`

3. ❌ **Forgetting R383 timestamp requirements**
   - Wrong: `IMPLEMENTATION-PLAN.md` (no timestamp)
   - Right: `EFFORT-IMPLEMENTATION-PLAN--20250120-100000.md`

4. ❌ **Mixing metadata with code**
   - Wrong: `/src/REVIEW-REPORT.md`
   - Right: `/.software-factory/phase1/wave1/effort/CODE-REVIEW-REPORT--20250120-120000.md`

## Enforcement

- **Grading Impact**: -50% for each violation (increased from -30%)
- **Agent Failure**: Agents cannot find plans without proper structure
- **Integration Block**: Cannot proceed without metadata organization

## Related Rules

- R344: Metadata location tracking (requires tracking in state)
- R345: Planning directory requirements (planning repository structure)
- R303: Phase/Wave document locations (planning vs effort separation)
- R340: Planning file metadata tracking (dual repository tracking)
- R383: Software Factory metadata file organization (timestamp requirements)
- R178: Effort directory structure

## Implementation Notes

1. Planning repository uses directory hierarchy, not `.software-factory/`
2. Effort repositories use `.software-factory/` for all metadata
3. R383 requires timestamps for all effort metadata files
4. Git hooks prevent cross-repository contamination
5. Clear separation enables easy agent discovery

## Key Principle

**"Planning lives in planning repository, effort metadata lives in .software-factory"**

This separation ensures:
- Clean repository boundaries
- Easy metadata discovery
- No cross-contamination
- Consistent agent behavior
- Clear ownership of documents

---

**Remember**: Planning repository uses directories, effort repositories use `.software-factory/`. NO MIXING.