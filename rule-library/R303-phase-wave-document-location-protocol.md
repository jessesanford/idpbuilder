# 🚨🚨🚨 BLOCKING RULE R303: Phase/Wave Document Location Protocol

## Rule Definition
ALL phase-level and wave-level planning documents MUST be stored in the `phase-plans/` directory in the Software Factory instance repository (NOT in effort branches). This centralized location ensures documents are accessible to all agents and persist across effort implementations.

## Criticality: 🚨🚨🚨 BLOCKING
Incorrect document placement causes agents to fail finding critical plans, blocks state transitions, and creates orphaned documents in effort branches that get lost.

## Requirements

### 1. Document Storage Hierarchy

#### Mandatory Directory Structure
```
${SF_INSTANCE_DIR}/
├── phase-plans/                                    # ALL phase/wave documents
│   ├── PHASE-1-ARCHITECTURE-PLAN.md              # Phase 1 architecture
│   ├── PHASE-1-IMPLEMENTATION-PLAN.md            # Phase 1 implementation  
│   ├── PHASE-1-WAVE-1-ARCHITECTURE-PLAN.md       # Wave 1-1 architecture
│   ├── PHASE-1-WAVE-1-IMPLEMENTATION-PLAN.md     # Wave 1-1 implementation
│   ├── PHASE-1-WAVE-1-MERGE-PLAN.md             # Wave 1-1 merge strategy
│   ├── PHASE-1-WAVE-1-REVIEW-REPORT.md          # Wave 1-1 review results
│   ├── PHASE-1-WAVE-2-ARCHITECTURE-PLAN.md       # Wave 1-2 architecture
│   ├── PHASE-1-WAVE-2-IMPLEMENTATION-PLAN.md     # Wave 1-2 implementation
│   └── ...
├── efforts/                                       # Effort implementations only
│   └── phase1/
│       └── wave1/
│           └── effort-001/
│               ├── .software-factory/            # Plan storage (NEW)
│               │   └── phase1/
│               │       └── wave1/
│               │           └── effort-001/
│               │               └── IMPLEMENTATION-PLAN-*.md
│               └── src/                          # Actual code
└── orchestrator-state.yaml                       # References to documents
```

### 2. Document Types and Locations

#### Phase-Level Documents (Created in SF Instance)
```bash
# ALWAYS in ${SF_INSTANCE_DIR}/phase-plans/
PHASE-${N}-ARCHITECTURE-PLAN.md         # Architect creates
PHASE-${N}-IMPLEMENTATION-PLAN.md       # Code Reviewer creates
PHASE-${N}-ASSESSMENT-REPORT.md         # Architect creates
PHASE-${N}-INTEGRATION-PLAN.md          # Integration Agent creates
PHASE-${N}-INTEGRATION-REPORT.md        # Post-integration results
```

#### Wave-Level Documents (Created in SF Instance)
```bash
# ALWAYS in ${SF_INSTANCE_DIR}/phase-plans/
PHASE-${N}-WAVE-${M}-ARCHITECTURE-PLAN.md    # Architect creates
PHASE-${N}-WAVE-${M}-IMPLEMENTATION-PLAN.md  # Code Reviewer creates
PHASE-${N}-WAVE-${M}-MERGE-PLAN.md          # Code Reviewer creates
PHASE-${N}-WAVE-${M}-REVIEW-REPORT.md       # Architect creates
PHASE-${N}-WAVE-${M}-INTEGRATION-REPORT.md  # Integration results
```

#### Effort-Level Documents (Created in Effort Branches)
```bash
# In /efforts/phase${N}/wave${M}/${effort-name}/.software-factory/phase${N}/wave${M}/${effort-name}/
IMPLEMENTATION-PLAN-YYYYMMDD-HHMMSS.md      # Effort-specific plan

# For split efforts, in .software-factory/phase${N}/wave${M}/${effort-name}-split-XXX/
SPLIT-PLAN-YYYYMMDD-HHMMSS.md              # If effort needs splitting

# Legacy locations (backward compatibility, deprecated):
# /efforts/phase${N}/wave${M}/${effort-name}/
IMPLEMENTATION-PLAN.md                     # Old location (root directory)
CODE-REVIEW-REPORT-YYYYMMDD-HHMMSS.md      # Review results
FIX-PLAN-YYYYMMDD-HHMMSS.md                # If fixes needed
```

### 3. Document Creation Protocol

#### Creating Phase Documents
```bash
create_phase_architecture_plan() {
    local PHASE="$1"
    
    # MUST be in SF instance directory
    cd "$SF_INSTANCE_DIR"
    
    # Create phase-plans directory if needed
    mkdir -p phase-plans
    
    # Create document with correct name
    local DOC_PATH="phase-plans/PHASE-${PHASE}-ARCHITECTURE-PLAN.md"
    
    cat > "$DOC_PATH" << 'EOF'
# Phase ${PHASE} Architecture Plan
Created: $(date -Iseconds)
Created By: architect
Location: ${SF_INSTANCE_DIR}/phase-plans/

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
    
    # MUST be in SF instance directory
    cd "$SF_INSTANCE_DIR"
    
    # Ensure directory exists
    mkdir -p phase-plans
    
    # Create document
    local DOC_PATH="phase-plans/PHASE-${PHASE}-WAVE-${WAVE}-IMPLEMENTATION-PLAN.md"
    
    cat > "$DOC_PATH" << 'EOF'
# Phase ${PHASE} Wave ${WAVE} Implementation Plan
Created: $(date -Iseconds)
Created By: code-reviewer
Location: ${SF_INSTANCE_DIR}/phase-plans/

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

The orchestrator-state.yaml MUST track document locations:
```yaml
phase_documents:
  phase_1:
    architecture_plan: "phase-plans/PHASE-1-ARCHITECTURE-PLAN.md"
    implementation_plan: "phase-plans/PHASE-1-IMPLEMENTATION-PLAN.md"
    assessment_report: "phase-plans/PHASE-1-ASSESSMENT-REPORT.md"
    created_at: "2025-01-20T10:00:00Z"
    
wave_documents:
  phase_1_wave_1:
    architecture_plan: "phase-plans/PHASE-1-WAVE-1-ARCHITECTURE-PLAN.md"
    implementation_plan: "phase-plans/PHASE-1-WAVE-1-IMPLEMENTATION-PLAN.md"
    merge_plan: "phase-plans/PHASE-1-WAVE-1-MERGE-PLAN.md"
    review_report: "phase-plans/PHASE-1-WAVE-1-REVIEW-REPORT.md"
    created_at: "2025-01-20T11:00:00Z"
```

### 5. Document Access Protocol

#### Reading Phase/Wave Documents
```bash
read_phase_document() {
    local PHASE="$1"
    local DOC_TYPE="$2"  # architecture, implementation, assessment
    
    # ALWAYS read from SF instance directory
    local DOC_PATH="${SF_INSTANCE_DIR}/phase-plans/PHASE-${PHASE}-${DOC_TYPE^^}-PLAN.md"
    
    if [ ! -f "$DOC_PATH" ]; then
        echo "❌ ERROR: Required document not found: $DOC_PATH"
        echo "Agent responsible for creating this document has not run yet"
        exit 1
    fi
    
    echo "📖 Reading: $DOC_PATH"
    cat "$DOC_PATH"
}

read_wave_document() {
    local PHASE="$1"
    local WAVE="$2"
    local DOC_TYPE="$3"  # architecture, implementation, merge, review
    
    # ALWAYS read from SF instance directory
    local DOC_PATH="${SF_INSTANCE_DIR}/phase-plans/PHASE-${PHASE}-WAVE-${WAVE}-${DOC_TYPE^^}-PLAN.md"
    
    if [ ! -f "$DOC_PATH" ]; then
        echo "❌ ERROR: Required document not found: $DOC_PATH"
        return 1
    fi
    
    echo "📖 Reading: $DOC_PATH"
    cat "$DOC_PATH"
}
```

### 6. Git Branch Management for Documents

#### Phase/Wave Documents Use Main Branch
```bash
# Phase/wave documents ALWAYS committed to main branch of SF instance
commit_phase_wave_document() {
    local DOC_PATH="$1"
    
    # Ensure we're in SF instance directory
    cd "$SF_INSTANCE_DIR"
    
    # Ensure we're on main branch
    CURRENT_BRANCH=$(git branch --show-current)
    if [ "$CURRENT_BRANCH" != "main" ]; then
        echo "⚠️ Switching to main branch for document commit"
        git checkout main
        git pull origin main
    fi
    
    # Add and commit
    git add "$DOC_PATH"
    git commit -m "docs: add $(basename "$DOC_PATH")"
    git push origin main
    
    echo "✅ Document committed to main branch"
}
```

#### Effort Documents Use Effort Branches
```bash
# Effort documents committed to effort branch in target repo
commit_effort_document() {
    local DOC_PATH="$1"
    local EFFORT_BRANCH="$2"
    
    # In effort directory (target repo working copy)
    cd "/efforts/phase*/wave*/${effort-name}"
    
    # Ensure correct branch
    git checkout "$EFFORT_BRANCH"
    
    # Add and commit
    git add "$DOC_PATH"
    git commit -m "docs: add $(basename "$DOC_PATH")"
    git push
}
```

### 7. Document Naming Conventions

#### Strict Naming Rules
```bash
# Phase documents: PHASE-${N}-${TYPE}-${SUFFIX}.md
PHASE-1-ARCHITECTURE-PLAN.md       # ✅ Correct
PHASE-1-IMPLEMENTATION-PLAN.md     # ✅ Correct
Phase1-Architecture.md              # ❌ Wrong format

# Wave documents: PHASE-${N}-WAVE-${M}-${TYPE}-${SUFFIX}.md
PHASE-1-WAVE-1-ARCHITECTURE-PLAN.md     # ✅ Correct
PHASE-1-WAVE-1-IMPLEMENTATION-PLAN.md   # ✅ Correct
Wave1-1-Plan.md                         # ❌ Wrong format

# Type must be one of:
ARCHITECTURE-PLAN
IMPLEMENTATION-PLAN
MERGE-PLAN
REVIEW-REPORT
ASSESSMENT-REPORT
INTEGRATION-REPORT
INTEGRATION-PLAN
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
        local DOC_PATH="${SF_INSTANCE_DIR}/phase-plans/PHASE-${PHASE}-WAVE-${WAVE}-${doc_type}.md"
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
1. **Create architecture plans** in phase-plans/
2. **Create assessment reports** in phase-plans/
3. **Never create documents** in effort directories
4. **Commit to main branch** of SF instance

### For Code Reviewer:
1. **Create phase/wave implementation plans** in phase-plans/ directory
2. **Create effort plans** in `.software-factory/phaseX/waveY/effort-name/` within effort directories
3. **Create split plans** in `.software-factory/phaseX/waveY/effort-split-XXX/` within effort directories
4. **Create merge plans** in phase-plans/
5. **Reference phase/wave plans** when creating effort plans

### For Orchestrator:
1. **Track document locations** in state file
2. **Verify documents exist** before state transitions
3. **Pass document paths** to spawned agents
4. **Create phase-plans directory** if missing

### For SW Engineer:
1. **Read wave plans** from phase-plans/
2. **Read effort plans** from `.software-factory/` subdirectory (new) or root (legacy)
3. **Read split plans** from `.software-factory/` subdirectory (new) or root (legacy)
4. **Never modify** phase/wave documents
5. **Create work logs** in effort directory root

### For Integration Agent:
1. **Read merge plans** from phase-plans/
2. **Create integration reports** in phase-plans/
3. **Reference effort documents** for context
4. **Update state file** with report locations

## Common Violations

### ❌ Creating Phase Documents in Wrong Location
```bash
# WRONG - Creating in effort directory
cd /efforts/phase1/wave1/effort-001
echo "# Phase 1 Architecture" > PHASE-1-ARCHITECTURE-PLAN.md

# CORRECT - Creating in SF instance
cd "$SF_INSTANCE_DIR"
mkdir -p phase-plans
echo "# Phase 1 Architecture" > phase-plans/PHASE-1-ARCHITECTURE-PLAN.md
```

### ❌ Wrong Branch for Document Commits
```bash
# WRONG - Committing to feature branch
git checkout phase1-wave1-integration
echo "# Wave Plan" > phase-plans/PHASE-1-WAVE-1-PLAN.md
git add . && git commit -m "add plan"

# CORRECT - Using main branch
git checkout main
echo "# Wave Plan" > phase-plans/PHASE-1-WAVE-1-IMPLEMENTATION-PLAN.md
git add . && git commit -m "docs: add wave 1-1 implementation plan"
```

### ❌ Incorrect Naming Convention
```bash
# WRONG - Various naming mistakes
phase1-plan.md                    # Missing PHASE- prefix
PHASE_1_PLAN.md                   # Underscore instead of dash
Phase-1-Plan.md                   # Wrong capitalization
WAVE-1-1-PLAN.md                  # Missing PHASE prefix

# CORRECT
PHASE-1-ARCHITECTURE-PLAN.md
PHASE-1-WAVE-1-IMPLEMENTATION-PLAN.md
```

## Error Messages

### Document Not Found
```
❌ CRITICAL: Required document not found
Expected: ${SF_INSTANCE_DIR}/phase-plans/PHASE-2-ARCHITECTURE-PLAN.md
Action: Architect must create this document before proceeding
State Transition: BLOCKED until document exists
```

### Wrong Location
```
❌ ERROR: Phase document created in wrong location
Found: /efforts/phase1/wave1/PHASE-1-PLAN.md
Required: ${SF_INSTANCE_DIR}/phase-plans/PHASE-1-IMPLEMENTATION-PLAN.md
Action: Move document to correct location and update references
```

### Naming Violation
```
❌ ERROR: Document naming convention violated
Found: Wave1-Implementation.md
Required Format: PHASE-${N}-WAVE-${M}-IMPLEMENTATION-PLAN.md
Example: PHASE-1-WAVE-1-IMPLEMENTATION-PLAN.md
```

## Verification Steps

### Verify Document Structure
```bash
verify_phase_plans_structure() {
    echo "Verifying phase-plans structure..."
    
    # Check directory exists
    if [ ! -d "${SF_INSTANCE_DIR}/phase-plans" ]; then
        echo "❌ Missing phase-plans directory"
        return 1
    fi
    
    # Check naming conventions
    find "${SF_INSTANCE_DIR}/phase-plans" -name "*.md" | while read doc; do
        basename "$doc" | grep -E "^PHASE-[0-9]+-" > /dev/null || {
            echo "❌ Invalid naming: $(basename "$doc")"
        }
    done
    
    echo "✅ Document structure verified"
}
```

## Related Rules
- R210: Architect Architecture Planning Protocol
- R211: Code Reviewer Implementation from Architecture
- R054: Implementation Plan Creation
- R212: Phase Directory Isolation Protocol
- R213: Wave and Effort Metadata Protocol

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