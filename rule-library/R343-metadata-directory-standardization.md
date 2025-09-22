# 🔴🔴🔴 RULE R343: Metadata Directory Standardization (SUPREME LAW) 🔴🔴🔴

## Rule Identifier
**Rule ID**: R343
**Category**: Infrastructure Consistency
**Criticality**: SUPREME LAW - Violations cause agent confusion and lost work
**Priority**: ABSOLUTE

## Description

ALL metadata, planning documents, reports, and non-code artifacts MUST be stored in `.software-factory/` directories within their respective workspaces. This standardizes metadata organization across efforts, splits, and integration workspaces.

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

### 2. UNIVERSAL METADATA DIRECTORY PATTERN

**EVERY workspace type MUST use `.software-factory/` for metadata:**

```bash
# Planning Workspace (Project Level)
/efforts/project-planning/
├── .software-factory/                    # ALL planning metadata here
│   ├── PROJECT-ARCHITECTURE.md
│   ├── PROJECT-TEST-PLAN.md
│   ├── PROJECT-DEMO-PLAN.md
│   ├── PROJECT-PHASES.md
│   └── project-metadata.yaml
└── [NO code - planning only]

# Planning Workspace (Phase Level)
/efforts/phase1/phase-planning/
├── .software-factory/                    # Phase planning metadata
│   ├── PHASE-ARCHITECTURE.md
│   ├── PHASE-TEST-PLAN.md
│   ├── PHASE-DEMO-PLAN.md
│   ├── PHASE-WAVES.md
│   └── phase-metadata.yaml
└── [NO code - planning only]

# Planning Workspace (Wave Level)
/efforts/phase1/wave1/wave-planning/
├── .software-factory/                    # Wave planning metadata
│   ├── WAVE-ARCHITECTURE.md
│   ├── WAVE-TEST-PLAN.md
│   ├── WAVE-DEMO-PLAN.md
│   ├── WAVE-EFFORTS.md
│   └── wave-metadata.yaml
└── [NO code - planning only]

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
│   ├── INTEGRATION-REPORT.md
│   ├── INTEGRATION-METADATA.md
│   └── integration-log.md
├── src/                                  # Merged code
└── tests/                                # Integration tests
```

### 2. FORBIDDEN PATTERNS

**NEVER place metadata in these locations:**
- ❌ Root directory of workspace (e.g., `/effort-name/PLAN.md`)
- ❌ Mixed with code directories (e.g., `/src/PLAN.md`)
- ❌ In a `/repo/` subdirectory (deprecated pattern)
- ❌ In the Software Factory instance directory
- ❌ Scattered across multiple locations

### 3. COMPREHENSIVE METADATA LIST

**ALL of these file types MUST be in .software-factory/:**

- Implementation Plans (IMPLEMENTATION-PLAN-*.md)
- Code Review Reports (CODE-REVIEW-REPORT-*.md)
- Split Plans (SPLIT-PLAN-*.md)
- Fix Plans (FIX-PLAN-*.md)
- Work Logs (work-log.md)
- Validation Results (validation-results.md)
- Test Results (test-results.md)
- Integration Reports (INTEGRATION-REPORT.md)
- Merge Plans (WAVE-MERGE-PLAN.md, PHASE-MERGE-PLAN.md)
- Architecture Documents (WAVE-ARCHITECTURE.md, PHASE-ARCHITECTURE.md)
- Demo Plans (WAVE-DEMO-PLAN.md, PHASE-DEMO-PLAN.md)
- Test Plans (WAVE-TEST-PLAN.md, PHASE-TEST-PLAN.md)
- Metadata YAML files (*.yaml)
- ALL other non-code documentation

### 4. METADATA TYPES BY WORKSPACE

#### Planning Workspace Metadata (Project/Phase/Wave)
```bash
# Project Planning
.software-factory/
├── PROJECT-ARCHITECTURE.md               # Overall system design
├── PROJECT-TEST-PLAN.md                 # End-to-end test strategy
├── PROJECT-DEMO-PLAN.md                 # Final demonstration approach
├── PROJECT-PHASES.md                    # Phase breakdown and goals
└── project-metadata.yaml                # Project tracking data

# Phase Planning
.software-factory/
├── PHASE-ARCHITECTURE.md                # Phase-specific design
├── PHASE-TEST-PLAN.md                   # Phase integration tests
├── PHASE-DEMO-PLAN.md                   # Phase demonstration
├── PHASE-WAVES.md                       # Wave breakdown and dependencies
└── phase-metadata.yaml                  # Phase tracking data

# Wave Planning
.software-factory/
├── WAVE-ARCHITECTURE.md                 # Wave implementation design
├── WAVE-TEST-PLAN.md                    # Wave-specific tests
├── WAVE-DEMO-PLAN.md                    # Wave demonstration
├── WAVE-EFFORTS.md                      # Effort breakdown and sizing
└── wave-metadata.yaml                   # Wave tracking data
```

#### Effort Workspace Metadata
```bash
.software-factory/
├── IMPLEMENTATION-PLAN-*.md              # From Code Reviewer
├── CODE-REVIEW-REPORT-*.md              # Review results
├── FIX-PLAN-*.md                        # Fix requirements
├── work-log.md                          # Development log
├── validation-results.md                # Test results
└── effort-metadata.yaml                 # Effort tracking
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
├── INTEGRATION-REPORT.md                # Integration results
├── INTEGRATION-METADATA.md              # Infrastructure details
├── integration-log.md                   # Operation log
├── conflict-resolutions.md              # If conflicts occurred
└── integration-metadata.yaml            # Integration tracking
```

### 4. CREATION PROTOCOL

**When creating ANY workspace:**

```bash
# Step 1: Clone target repository
git clone "$TARGET_REPO_URL" "$WORKSPACE_DIR"
cd "$WORKSPACE_DIR"

# Step 2: IMMEDIATELY create .software-factory directory
mkdir -p .software-factory

# Step 3: Add to .gitignore if needed
if ! grep -q "^\.software-factory" .gitignore 2>/dev/null; then
    echo "# Software Factory metadata (local only)" >> .gitignore
    echo ".software-factory/" >> .gitignore
fi

# Step 4: Place ALL metadata in .software-factory
echo "# Work Log" > .software-factory/work-log.md
```

### 5. AGENT RESPONSIBILITIES

#### SW Engineer
- MUST check for `.software-factory/` on startup
- MUST read plans from `.software-factory/`
- MUST write logs to `.software-factory/`
- NEVER create metadata in root directory

#### Code Reviewer
- MUST write all plans to `.software-factory/`
- MUST read code from proper directories
- MUST save reports to `.software-factory/`
- NEVER mix metadata with code

#### Integration Agent
- MUST create `.software-factory/` immediately after clone
- MUST store merge plans in `.software-factory/`
- MUST write reports to `.software-factory/`
- NEVER use `/repo/` subdirectory pattern

#### Orchestrator
- MUST verify `.software-factory/` exists in workspaces
- MUST track metadata locations in state file
- MUST instruct agents about metadata location
- NEVER place metadata in wrong location

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
validate_metadata_structure() {
    local workspace="$1"
    
    # Check .software-factory exists
    if [ ! -d "$workspace/.software-factory" ]; then
        echo "❌ R343 VIOLATION: Missing .software-factory directory"
        return 1
    fi
    
    # Check no metadata in root
    if ls "$workspace"/*.md 2>/dev/null | grep -E "(PLAN|REPORT|REVIEW|ARCHITECTURE|TEST|DEMO|LOG)" > /dev/null; then
        echo "❌ R343 VIOLATION: Metadata files in root directory"
        return 1
    fi
    
    # Check no metadata in code directories
    if find "$workspace/src" "$workspace/tests" -name "*.md" 2>/dev/null | grep -E "(PLAN|REPORT|REVIEW)" > /dev/null; then
        echo "❌ R343 VIOLATION: Metadata mixed with code"
        return 1
    fi
    
    # Check no /repo subdirectory
    if [ -d "$workspace/repo" ]; then
        echo "❌ R343 VIOLATION: Deprecated /repo/ subdirectory pattern"
        return 1
    fi
    
    # Verify all metadata is in .software-factory
    local metadata_count=$(find "$workspace" -name "*PLAN*.md" -o -name "*REPORT*.md" -o -name "*work-log*.md" 2>/dev/null | grep -v ".software-factory" | wc -l)
    if [ "$metadata_count" -gt 0 ]; then
        echo "❌ R343 VIOLATION: Found $metadata_count metadata files outside .software-factory"
        return 1
    fi
    
    echo "✅ R343 COMPLIANT: Proper metadata structure"
    return 0
}
```

## Common Violations

1. ❌ **Creating metadata in root directory**
   - Wrong: `/effort-name/IMPLEMENTATION-PLAN.md`
   - Right: `/effort-name/.software-factory/IMPLEMENTATION-PLAN.md`

2. ❌ **Using /repo/ subdirectory for integration**
   - Wrong: `/integration-workspace/repo/`
   - Right: `/integration-workspace/` (clone directly)

3. ❌ **Forgetting to create .software-factory**
   - Wrong: Start work without metadata directory
   - Right: Create `.software-factory/` immediately after clone

4. ❌ **Mixing metadata with code**
   - Wrong: `/src/REVIEW-REPORT.md`
   - Right: `/.software-factory/REVIEW-REPORT.md`

## Enforcement

- **Grading Impact**: -50% for each violation (increased from -30%)
- **Agent Failure**: Agents cannot find plans without proper structure
- **Integration Block**: Cannot proceed without metadata organization

## Related Rules

- R344: Metadata location tracking (NEW - requires tracking in state)
- R345: Planning branch requirements (NEW - dedicated planning branches)
- R303: Phase/Wave document locations
- R340: Planning file metadata tracking
- R250: Integration isolation (updated to remove /repo/)
- R178: Effort directory structure

## Implementation Notes

1. `.software-factory/` is created immediately after workspace setup
2. Can be added to .gitignore for local-only metadata
3. Can be committed for persistent metadata
4. Consistent across ALL workspace types
5. Enables easy agent discovery of metadata

## Key Principle

**"Code lives in code directories, metadata lives in .software-factory"**

This separation ensures:
- Clean git history
- Easy metadata discovery
- No confusion about file purposes
- Consistent agent behavior
- Simple recovery patterns

---

**Remember**: EVERY workspace has `.software-factory/`. NO EXCEPTIONS.