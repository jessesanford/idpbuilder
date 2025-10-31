# 🚨🚨🚨 BLOCKING RULE R303: Phase/Wave Document Location Protocol

## Rule Definition
ALL project-level, phase-level and wave-level planning documents MUST be stored in the `$CLAUDE_PROJECT_DIR/phase-plans/` directory hierarchy (NOT in effort branches). Effort-specific implementation plans MUST be stored in `.software-factory/` directories within effort branches. This separation ensures clean repository boundaries and prevents cross-contamination between planning and implementation.

## Criticality: 🚨🚨🚨 BLOCKING
Incorrect document placement causes agents to fail finding critical plans, blocks state transitions, and creates orphaned documents in effort branches that get lost.

## Requirements

### 1. Document Storage Hierarchy

#### Mandatory Directory Structure
```
# PLANNING REPOSITORY (Software Factory Instance)
$CLAUDE_PROJECT_DIR/
├── phase-plans/                                   # ALL project/phase/wave planning
│   ├── project/                                  # Project-level planning
│   │   ├── PROJECT-ARCHITECTURE-PLAN--TIMESTAMP.md         # Overall architecture
│   │   ├── PROJECT-IMPLEMENTATION-PLAN--TIMESTAMP.md       # Master implementation plan
│   │   └── PROJECT-TEST-STRATEGY--TIMESTAMP.md            # Test strategy
│   ├── phase1/                                  # Phase 1 planning
│   │   ├── PHASE-1-ARCHITECTURE-PLAN--TIMESTAMP.md        # Phase 1 architecture
│   │   ├── PHASE-1-PLAN--TIMESTAMP.md                     # Phase 1 implementation
│   │   ├── PHASE-1-ASSESSMENT-REPORT--TIMESTAMP.md        # Phase 1 assessment
│   │   ├── wave1/                              # Wave 1-1 planning
│   │   │   ├── WAVE-1-1-ARCHITECTURE-PLAN--TIMESTAMP.md   # Wave architecture
│   │   │   ├── WAVE-1-1-PLAN--TIMESTAMP.md                # Wave implementation
│   │   │   ├── WAVE-1-1-MERGE-PLAN--TIMESTAMP.md         # Wave merge strategy
│   │   │   └── WAVE-1-1-REVIEW-REPORT--TIMESTAMP.md      # Wave review results
│   │   └── wave2/                              # Wave 1-2 planning
│   │       └── ...
│   └── phase2/                                  # Phase 2 planning
│       └── ...
└── orchestrator-state-v3.json                      # References to all documents

# EFFORT REPOSITORY (Target Implementation)
efforts/
└── phase1/
    └── wave1/
        └── effort-001/
            ├── .software-factory/               # Effort-specific metadata ONLY
            │   └── phase1/
            │       └── wave1/
            │           └── effort-001/
            │               ├── EFFORT-IMPLEMENTATION-PLAN--YYYYMMDD-HHMMSS.md
            │               ├── CODE-REVIEW-REPORT--YYYYMMDD-HHMMSS.md
            │               └── FIX-PLAN--YYYYMMDD-HHMMSS.md
            ├── src/                             # Actual implementation code
            └── tests/                           # Test code
```

### 2. Document Types and Locations

#### Project-Level Documents (Planning Repository)
```bash
# ALWAYS in $CLAUDE_PROJECT_DIR/phase-plans/project/
PROJECT-ARCHITECTURE-PLAN--YYYYMMDD-HHMMSS.md            # Architect creates
PROJECT-IMPLEMENTATION-PLAN--YYYYMMDD-HHMMSS.md          # Master plan (initial)
PROJECT-TEST-STRATEGY--YYYYMMDD-HHMMSS.md               # Overall test strategy
PROJECT-INTEGRATE_WAVE_EFFORTS-PLAN--YYYYMMDD-HHMMSS.md            # Final integration plan
```

#### Phase-Level Documents (Planning Repository)
```bash
# ALWAYS in $CLAUDE_PROJECT_DIR/phase-plans/phaseN/
PHASE-${N}-ARCHITECTURE-PLAN--YYYYMMDD-HHMMSS.md         # Architect creates
PHASE-${N}-PLAN--YYYYMMDD-HHMMSS.md                      # Code Reviewer creates
PHASE-${N}-ASSESSMENT-REPORT--YYYYMMDD-HHMMSS.md         # Architect creates
PHASE-${N}-INTEGRATE_WAVE_EFFORTS-PLAN--YYYYMMDD-HHMMSS.md          # Integration Agent creates
PHASE-${N}-INTEGRATE_WAVE_EFFORTS-REPORT--YYYYMMDD-HHMMSS.md        # Post-integration results
```

#### Wave-Level Documents (Planning Repository)
```bash
# ALWAYS in $CLAUDE_PROJECT_DIR/phase-plans/phaseN/waveM/
WAVE-${N}-${M}-ARCHITECTURE-PLAN--YYYYMMDD-HHMMSS.md     # Architect creates
WAVE-${N}-${M}-PLAN--YYYYMMDD-HHMMSS.md                  # Code Reviewer creates
WAVE-${N}-${M}-MERGE-PLAN--YYYYMMDD-HHMMSS.md           # Code Reviewer creates
WAVE-${N}-${M}-REVIEW-REPORT--YYYYMMDD-HHMMSS.md        # Architect creates
WAVE-${N}-${M}-INTEGRATE_WAVE_EFFORTS-REPORT--YYYYMMDD-HHMMSS.md   # Integration results
```

#### Effort-Level Documents (Created in Effort Branches)
```bash
# 🔴🔴🔴 R383 MANDATORY: All metadata files MUST have timestamps
# Use sf_metadata_path helper function from utilities/sf-metadata-path.sh

# In /efforts/phase${N}/wave${M}/${effort-name}/.software-factory/phase${N}/wave${M}/${effort-name}/
IMPLEMENTATION-PLAN--YYYYMMDD-HHMMSS.md     # Effort-specific plan (R383 format)
CODE-REVIEW-REPORT--YYYYMMDD-HHMMSS.md      # Review results (R383 format)
FIX-PLAN--YYYYMMDD-HHMMSS.md                # If fixes needed (R383 format)

# For split efforts, in .software-factory/phase${N}/wave${M}/${effort-name}-split-XXX/
SPLIT-PLAN--YYYYMMDD-HHMMSS.md              # If effort needs splitting (R383 format)

# ❌ FORBIDDEN (R383 violations):
IMPLEMENTATION-PLAN.md                      # NO TIMESTAMP - VIOLATION!
CODE-REVIEW-REPORT.md                       # NO TIMESTAMP - VIOLATION!
```

### 3. Document Creation Protocol

#### Creating Phase Documents
```bash
create_phase_architecture_plan() {
    local PHASE="$1"
    local TIMESTAMP=$(date +%Y%m%d-%H%M%S)

    # MUST be in planning directory
    cd "$CLAUDE_PROJECT_DIR"

    # Create planning directory structure if needed
    mkdir -p phase-plans/phase${PHASE}

    # Create document with correct name and timestamp
    local DOC_PATH="phase-plans/phase${PHASE}/PHASE-${PHASE}-ARCHITECTURE-PLAN--${TIMESTAMP}.md"

    cat > "$DOC_PATH" << 'EOF'
# Phase ${PHASE} Architecture Plan
Created: $(date -Iseconds)
Created By: architect
Location: $CLAUDE_PROJECT_DIR/phase-plans/phase${PHASE}/

## Phase Vision
[Architecture content here]
EOF
    
    # Commit to SF instance repo
    git add "$DOC_PATH"
    git commit -m "architect: create Phase ${PHASE} architecture plan"
    git push
    
    echo "✅ Created: $DOC_PATH"
}
```

#### Creating Wave Documents
```bash
create_wave_implementation_plan() {
    local PHASE="$1"
    local WAVE="$2"
    local TIMESTAMP=$(date +%Y%m%d-%H%M%S)

    # MUST be in planning directory
    cd "$CLAUDE_PROJECT_DIR"

    # Ensure directory exists
    mkdir -p phase-plans/phase${PHASE}/wave${WAVE}

    # Create document with timestamp
    local DOC_PATH="phase-plans/phase${PHASE}/wave${WAVE}/WAVE-${PHASE}-${WAVE}-PLAN--${TIMESTAMP}.md"

    cat > "$DOC_PATH" << 'EOF'
# Phase ${PHASE} Wave ${WAVE} Implementation Plan
Created: $(date -Iseconds)
Created By: code-reviewer
Location: $CLAUDE_PROJECT_DIR/phase-plans/phase${PHASE}/wave${WAVE}/

## Wave Implementation Strategy
[Implementation content here]
EOF
    
    # Commit to SF instance repo
    git add "$DOC_PATH"
    git commit -m "code-reviewer: create Phase ${PHASE} Wave ${WAVE} implementation plan"
    git push
    
    echo "✅ Created: $DOC_PATH"
}
```

### 4. Document Reference in State File

The orchestrator-state-v3.json MUST track document locations:
```yaml
project_documents:
    architecture_plan: "phase-plans/project/PROJECT-ARCHITECTURE-PLAN--20250120-090000.md"
    implementation_plan: "phase-plans/project/PROJECT-IMPLEMENTATION-PLAN--20250120-091500.md"
    test_strategy: "phase-plans/project/PROJECT-TEST-STRATEGY--20250120-093000.md"
    created_at: "2025-01-20T09:00:00Z"

phase_documents:
  phase_1:
    architecture_plan: "phase-plans/phase1/PHASE-1-ARCHITECTURE-PLAN--20250120-100000.md"
    implementation_plan: "phase-plans/phase1/PHASE-1-PLAN--20250120-101500.md"
    assessment_report: "phase-plans/phase1/PHASE-1-ASSESSMENT-REPORT--20250120-103000.md"
    created_at: "2025-01-20T10:00:00Z"

wave_documents:
  phase_1_wave_1:
    architecture_plan: "phase-plans/phase1/wave1/WAVE-1-1-ARCHITECTURE-PLAN--20250120-110000.md"
    implementation_plan: "phase-plans/phase1/wave1/WAVE-1-1-PLAN--20250120-111500.md"
    merge_plan: "phase-plans/phase1/wave1/WAVE-1-1-MERGE-PLAN--20250120-113000.md"
    review_report: "phase-plans/phase1/wave1/WAVE-1-1-REVIEW-REPORT--20250120-114500.md"
    created_at: "2025-01-20T11:00:00Z"
```

### 5. Document Access Protocol

#### Reading Phase/Wave Documents
```bash
read_phase_document() {
    local PHASE="$1"
    local DOC_TYPE="$2"  # architecture, plan, assessment

    # Use wildcard to find latest timestamped document
    local DOC_PATH=$(ls -t "${CLAUDE_PROJECT_DIR}/phase-plans/phase${PHASE}/PHASE-${PHASE}-${DOC_TYPE^^}"--*.md 2>/dev/null | head -1)

    if [ -z "$DOC_PATH" ] || [ ! -f "$DOC_PATH" ]; then
        echo "❌ ERROR: Required document not found: PHASE-${PHASE}-${DOC_TYPE^^}--*.md"
        echo "Agent responsible for creating this document has not run yet"
        exit 1
    fi

    echo "📖 Reading: $DOC_PATH"
    cat "$DOC_PATH"
}

read_wave_document() {
    local PHASE="$1"
    local WAVE="$2"
    local DOC_TYPE="$3"  # architecture, plan, merge, review

    # Use wildcard to find latest timestamped document
    local DOC_PATH=$(ls -t "${CLAUDE_PROJECT_DIR}/phase-plans/phase${PHASE}/wave${WAVE}/WAVE-${PHASE}-${WAVE}-${DOC_TYPE^^}"--*.md 2>/dev/null | head -1)

    if [ -z "$DOC_PATH" ] || [ ! -f "$DOC_PATH" ]; then
        echo "❌ ERROR: Required document not found: WAVE-${PHASE}-${WAVE}-${DOC_TYPE^^}--*.md"
        return 1
    fi

    echo "📖 Reading: $DOC_PATH"
    cat "$DOC_PATH"
}
```

### 6. Git Management for Documents

#### Planning Documents Are Committed to Planning Repository
```bash
# Planning documents ALWAYS committed to planning repository (SF instance)
commit_planning_document() {
    local DOC_PATH="$1"

    # Ensure we're in planning repository
    cd "$CLAUDE_PROJECT_DIR"

    # Add and commit to current branch (usually main or a planning branch)
    git add "$DOC_PATH"
    git commit -m "planning: add $(basename "$DOC_PATH")"
    git push

    echo "✅ Planning document committed to planning repository"
}
```

#### Effort Documents Are Committed to Effort Branches
```bash
# Effort documents committed to effort branch in target repo
commit_effort_document() {
    local DOC_PATH="$1"
    local EFFORT_BRANCH="$2"

    # In effort directory (target repo working copy)
    cd "/efforts/phase*/wave*/${effort-name}"

    # Ensure correct branch
    git checkout "$EFFORT_BRANCH"

    # Add and commit (with R383 timestamp)
    git add "$DOC_PATH"
    git commit -m "effort: add $(basename "$DOC_PATH")"
    git push

    echo "✅ Effort document committed to effort branch"
}
```

### 7. Document Naming Conventions

#### Strict Naming Rules
```bash
# Phase documents: PHASE-${N}-${TYPE}--TIMESTAMP.md
PHASE-1-ARCHITECTURE-PLAN--20250120-100000.md       # ✅ Correct
PHASE-1-PLAN--20250120-101500.md                    # ✅ Correct
Phase1-Architecture.md                              # ❌ Wrong format (no timestamp, wrong format)

# Wave documents: WAVE-${N}-${M}-${TYPE}--TIMESTAMP.md
WAVE-1-1-ARCHITECTURE-PLAN--20250120-110000.md     # ✅ Correct
WAVE-1-1-PLAN--20250120-111500.md                  # ✅ Correct
Wave1-1-Plan.md                                    # ❌ Wrong format (no timestamp, wrong prefix)

# Type must be one of:
ARCHITECTURE-PLAN
PLAN
MERGE-PLAN
REVIEW-REPORT
ASSESSMENT-REPORT
INTEGRATE_WAVE_EFFORTS-REPORT
INTEGRATE_WAVE_EFFORTS-PLAN
```

### 8. Cross-Agent Document Discovery

#### Orchestrator Finding Documents
```bash
discover_wave_documents() {
    local PHASE="$1"
    local WAVE="$2"
    
    echo "🔍 Discovering documents for Phase $PHASE Wave $WAVE"
    
    # Check for expected documents
    local DOCS=(
        "ARCHITECTURE-PLAN"
        "IMPLEMENTATION-PLAN"
        "MERGE-PLAN"
        "REVIEW-REPORT"
    )
    
    for doc_type in "${DOCS[@]}"; do
        local DOC_PATH=$(ls -t "${CLAUDE_PROJECT_DIR}/phase-plans/phase${PHASE}/wave${WAVE}/WAVE-${PHASE}-${WAVE}-${doc_type}"--*.md 2>/dev/null | head -1)
        if [ -f "$DOC_PATH" ]; then
            echo "  ✅ Found: $(basename "$DOC_PATH")"
        else
            echo "  ⚠️ Missing: $(basename "$DOC_PATH")"
        fi
    done
}
```

## Implementation Requirements

### For Architect:
1. **Create architecture plans** in $CLAUDE_PROJECT_DIR/phase-plans/phaseN/ and phase-plans/phaseN/waveM/
2. **Create assessment reports** in $CLAUDE_PROJECT_DIR/phase-plans/phaseN/
3. **Never create documents** in effort directories
4. **Commit to planning repository** (SF instance)

### For Code Reviewer:
1. **Create phase/wave implementation plans** in $CLAUDE_PROJECT_DIR/phase-plans/phaseN/ and phase-plans/phaseN/waveM/
2. **Create effort-specific plans** in `.software-factory/phaseX/waveY/effort-name/` within effort directories
3. **Create split plans** in `.software-factory/phaseX/waveY/effort-split-XXX/` within effort directories
4. **Create merge plans** in $CLAUDE_PROJECT_DIR/phase-plans/phaseN/waveM/
5. **Reference phase/wave plans** when creating effort plans

### For Orchestrator:
1. **Track document locations** in state file
2. **Verify documents exist** before state transitions
3. **Pass document paths** to spawned agents
4. **Create planning directory structure** if missing

### For SW Engineer:
1. **Read wave plans** from $CLAUDE_PROJECT_DIR/phase-plans/phaseN/waveM/
2. **Read effort plans** from `.software-factory/` subdirectory
3. **Read split plans** from `.software-factory/` subdirectory
4. **Never modify** phase/wave documents in planning repository
5. **Create work logs** in `.software-factory/` directory

### For Integration Agent:
1. **Read merge plans** from $CLAUDE_PROJECT_DIR/phase-plans/phaseN/waveM/
2. **Create integration reports** in $CLAUDE_PROJECT_DIR/phase-plans/phaseN/waveM/
3. **Reference effort documents** for context
4. **Update state file** with report locations

## Common Violations

### ❌ Creating Phase Documents in Wrong Location
```bash
# WRONG - Creating in effort directory
cd /efforts/phase1/wave1/effort-001
echo "# Phase 1 Architecture" > PHASE-1-ARCHITECTURE-PLAN.md

# CORRECT - Creating in planning repository with timestamp
cd "$CLAUDE_PROJECT_DIR"
mkdir -p phase-plans/phase1
echo "# Phase 1 Architecture" > phase-plans/phase1/PHASE-1-ARCHITECTURE-PLAN--$(date +%Y%m%d-%H%M%S).md
```

### ❌ Mixing Planning and Effort Files
```bash
# WRONG - Planning file in effort directory
cd /efforts/phase1/wave1/effort-001
echo "# Wave Plan" > WAVE-1-1-IMPLEMENTATION-PLAN.md

# CORRECT - Planning files in planning repository with timestamp
cd "$CLAUDE_PROJECT_DIR"
mkdir -p phase-plans/phase1/wave1
echo "# Wave Plan" > phase-plans/phase1/wave1/WAVE-1-1-PLAN--$(date +%Y%m%d-%H%M%S).md
```

### ❌ Incorrect Naming Convention
```bash
# WRONG - Various naming mistakes
phase1-plan.md                    # Missing PHASE- prefix
PHASE_1_PLAN.md                   # Underscore instead of dash
Phase-1-Plan.md                   # Wrong capitalization
PHASE-1-WAVE-1-PLAN.md            # Old naming format

# CORRECT (with timestamps)
PHASE-1-ARCHITECTURE-PLAN--20250120-100000.md      # Phase architecture document
PHASE-1-PLAN--20250120-101500.md                   # Phase implementation document
WAVE-1-1-ARCHITECTURE-PLAN--20250120-110000.md     # Wave architecture document
WAVE-1-1-PLAN--20250120-111500.md                  # Wave implementation document
```

## Error Messages

### Document Not Found
```
❌ CRITICAL: Required document not found
Expected: $CLAUDE_PROJECT_DIR/phase-plans/phase2/PHASE-2-ARCHITECTURE-PLAN--*.md
Action: Architect must create this document before proceeding
State Transition: BLOCKED until document exists
```

### Wrong Location
```
❌ ERROR: Phase document created in wrong location
Found: /efforts/phase1/wave1/PHASE-1-PLAN.md
Required: $CLAUDE_PROJECT_DIR/phase-plans/phase1/PHASE-1-PLAN--YYYYMMDD-HHMMSS.md
Action: Move document to correct location with timestamp and update references
```

### Naming Violation
```
❌ ERROR: Document naming convention violated
Found: Wave1-Implementation.md
Required Format: WAVE-${N}-${M}-PLAN--YYYYMMDD-HHMMSS.md
Example: WAVE-1-1-PLAN--20250120-111500.md
```

## Verification Steps

### Verify Document Structure
```bash
verify_planning_structure() {
    echo "Verifying planning structure..."

    # Check directory exists
    if [ ! -d "${CLAUDE_PROJECT_DIR}/phase-plans" ]; then
        echo "❌ Missing phase-plans directory"
        return 1
    fi

    # Check naming conventions for phase docs (with timestamp)
    find "${CLAUDE_PROJECT_DIR}/phase-plans" -path "*/phase*/PHASE-*.md" | while read doc; do
        basename "$doc" | grep -E "^PHASE-[0-9]+-(ARCHITECTURE-PLAN|PLAN|ASSESSMENT-REPORT)--[0-9]{8}-[0-9]{6}\.md$" > /dev/null || {
            echo "❌ Invalid phase doc naming: $(basename "$doc")"
        }
    done

    # Check naming conventions for wave docs (with timestamp)
    find "${CLAUDE_PROJECT_DIR}/phase-plans" -path "*/wave*/WAVE-*.md" | while read doc; do
        basename "$doc" | grep -E "^WAVE-[0-9]+-[0-9]+-(ARCHITECTURE-PLAN|PLAN|MERGE-PLAN|REVIEW-REPORT)--[0-9]{8}-[0-9]{6}\.md$" > /dev/null || {
            echo "❌ Invalid wave doc naming: $(basename "$doc")"
        }
    done

    echo "✅ Planning structure verified"
}
```

## Related Rules
- R210: Architect Architecture Planning Protocol
- R211: Code Reviewer Implementation from Architecture
- R054: Implementation Plan Creation
- R212: Phase Directory Isolation Protocol
- R213: Wave and Effort Metadata Protocol
- R343: Metadata Directory Standardization
- R383: Software Factory Metadata File Organization (MANDATORY timestamps)

## Penalties
- Creating documents in wrong location: -30%
- Wrong branch for commits: -20%
- Incorrect naming convention: -15%
- Missing document references in state: -10%
- Not creating phase-plans directory: -25%

---
*Rule Type*: Protocol
*Agents*: All Agents
*Enforcement*: Directory structure validation and state file tracking