# Code Reviewer - CREATE_SPLIT_PLAN State Rules

## State Context
You are creating individual split plan files for an oversized implementation (>800 lines). These plans will be saved in the too-large branch and later copied to split effort directories by the orchestrator.

## 🔴🔴🔴 CRITICAL: Sequential Branching Strategy 🔴🔴🔴

**ALL SPLITS MUST BE PLANNED FOR SEQUENTIAL BRANCHING!**

### The Mandatory Pattern:
```
Split-001: Based on phase-integration (same as original)
    ↓ (becomes base for next)
Split-002: Based on Split-001 (NOT phase-integration!)
    ↓ (becomes base for next)
Split-003: Based on Split-002 (NOT phase-integration!)
```

### Why This is CRITICAL:
1. **Line Counting**: Each split measures ONLY its additions (400 lines, not cumulative)
2. **Dependencies**: Later splits can use earlier split code
3. **Clean Integration**: No merge conflicts between splits
4. **Progressive Building**: Each split extends the previous

### Include in Every Split Plan:
```markdown
## Branching Strategy
- **Split-001**: Branches from `phase-integration` (same as original)
- **Split-002**: Branches from `split-001` (NOT phase-integration!)
- **Split-003**: Branches from `split-002` (NOT phase-integration!)

This SEQUENTIAL branching ensures each split measures only its own additions.
```

## 🔴🔴🔴 CRITICAL: Split Plan File Management 🔴🔴🔴

**YOU MUST CREATE AND COMMIT SPLIT PLANS IN THE TOO-LARGE BRANCH**

### Step 1: Verify You're in the Too-Large Branch
```bash
echo "═══════════════════════════════════════════════════════"
echo "🔍 VERIFYING TOO-LARGE BRANCH CONTEXT"
echo "═══════════════════════════════════════════════════════"

# Verify current directory and branch
CURRENT_DIR=$(pwd)
CURRENT_BRANCH=$(git branch --show-current)

echo "Current directory: $CURRENT_DIR"
echo "Current branch: $CURRENT_BRANCH"

# Verify this is the too-large branch that needs splitting
if [ ! -f "IMPLEMENTATION-PLAN.md" ]; then
    echo "❌ ERROR: No IMPLEMENTATION-PLAN.md found"
    exit 1
fi

# Check if we already have split plans
if ls SPLIT-PLAN-*.md 2>/dev/null | head -1; then
    echo "⚠️ WARNING: Split plans already exist:"
    ls -la SPLIT-PLAN-*.md
    echo "Will overwrite if proceeding..."
fi
```

### Step 2: Create Split Inventory File
```bash
# Create SPLIT-INVENTORY.md that lists all splits
cat > SPLIT-INVENTORY.md << 'EOF'
# Split Inventory for ${EFFORT_NAME}

## Overview
This effort exceeded 800 lines and has been split into manageable sub-efforts.

- **Original Size**: ${ORIGINAL_SIZE} lines
- **Number of Splits**: ${NUM_SPLITS}
- **Date Split**: $(date '+%Y-%m-%d')
- **Split By**: Code Reviewer Agent

## Split Structure

| Split # | Name | Description | Est. Lines | Status |
|---------|------|-------------|------------|--------|
| 001 | ${SPLIT_001_NAME} | ${SPLIT_001_DESC} | ${SPLIT_001_LINES} | Planned |
| 002 | ${SPLIT_002_NAME} | ${SPLIT_002_DESC} | ${SPLIT_002_LINES} | Planned |
| 003 | ${SPLIT_003_NAME} | ${SPLIT_003_DESC} | ${SPLIT_003_LINES} | Planned |

## Integration Strategy
${INTEGRATION_STRATEGY}

## Files Distribution

### Split-001
- Files to implement:
  - ${FILE_LIST_001}

### Split-002
- Files to implement:
  - ${FILE_LIST_002}

### Split-003
- Files to implement:
  - ${FILE_LIST_003}

## Dependencies
- Split-002 depends on Split-001 (and MUST be branched from Split-001!)
- Split-003 depends on Split-002 (and MUST be branched from Split-002!)
- All splits must be completed sequentially with SEQUENTIAL BRANCHING

## Validation
Each split must:
- Stay under 800 lines
- Pass all tests independently
- Integrate cleanly with previous splits
EOF
```

### Step 3: Create Individual Split Plan Files
```bash
# For each split, create SPLIT-PLAN-XXX.md
for SPLIT_NUM in 001 002 003; do
    cat > SPLIT-PLAN-${SPLIT_NUM}.md << 'EOF'
# Split Plan ${SPLIT_NUM} - ${SPLIT_NAME}

## Split Metadata
- **Split Number**: ${SPLIT_NUM}
- **Parent Effort**: ${PARENT_EFFORT_NAME}
- **Original Branch**: ${TOO_LARGE_BRANCH}
- **Target Size**: ${TARGET_SIZE} lines (max 800)
- **Created**: $(date '+%Y-%m-%d %H:%M:%S')

## Implementation Scope

### Files to Create/Modify
${FILE_LIST_FOR_SPLIT}

### Functionality to Implement
${FUNCTIONALITY_DESCRIPTION}

### Excluded from This Split
${EXCLUDED_ITEMS}
(These will be handled in other splits)

## Technical Requirements

### Dependencies
- External dependencies:
  ${EXTERNAL_DEPS}
- From previous splits:
  ${PREVIOUS_SPLIT_DEPS}

### Interfaces to Provide
${PROVIDED_INTERFACES}

### Interfaces to Consume
${CONSUMED_INTERFACES}

## Implementation Instructions

### Step 1: Setup
1. Verify you're in the split directory (not the too-large directory)
2. Confirm branch is ${SPLIT_BRANCH_NAME}
3. Verify base files from previous splits (if any)

### Step 2: Implementation
${DETAILED_IMPLEMENTATION_STEPS}

### Step 3: Testing
${TEST_REQUIREMENTS}

### Step 4: Integration
${INTEGRATION_STEPS}

## Size Management
- Target: ${TARGET_SIZE} lines
- Buffer: 100 lines (implement up to 700 lines)
- Measurement: Use line-counter.sh before committing

## Success Criteria
- [ ] All specified files implemented
- [ ] Size under 800 lines (measured)
- [ ] Tests passing
- [ ] Integrates with previous splits (if applicable)
- [ ] No functionality regression

## Notes for SW Engineer
${SPECIAL_NOTES}
EOF
done

echo "✅ Created split plan files:"
ls -la SPLIT-PLAN-*.md
```

### Step 4: Commit and Push Split Plans to Too-Large Branch
```bash
echo "═══════════════════════════════════════════════════════"
echo "📦 COMMITTING SPLIT PLANS TO TOO-LARGE BRANCH"
echo "═══════════════════════════════════════════════════════"

# Add all split-related files
git add SPLIT-INVENTORY.md
git add SPLIT-PLAN-*.md

# Show what we're committing
echo "Files to commit:"
git status --short

# Commit with descriptive message
git commit -m "feat: add split plans for oversized effort

Created split plans to break down ${ORIGINAL_SIZE}-line implementation:
- SPLIT-INVENTORY.md: Overview of all splits
- SPLIT-PLAN-001.md: ${SPLIT_001_NAME} (${SPLIT_001_LINES} lines)
- SPLIT-PLAN-002.md: ${SPLIT_002_NAME} (${SPLIT_002_LINES} lines)
- SPLIT-PLAN-003.md: ${SPLIT_003_NAME} (${SPLIT_003_LINES} lines)

This branch will be abandoned after splits are implemented and merged."

# Push to remote
echo "Pushing to remote..."
git push

echo "✅ Split plans committed and pushed to: $CURRENT_BRANCH"
echo "   These files will be copied by orchestrator to split directories"
```

## Required Split Plan Elements

### SPLIT-INVENTORY.md Must Include:
1. **Overview** - Why the split was needed
2. **Split Structure** - Table of all splits with sizes
3. **File Distribution** - Which files go in which split
4. **Integration Strategy** - How splits will merge
5. **Dependencies** - Order of implementation

### Each SPLIT-PLAN-XXX.md Must Include:
1. **Metadata** - Split number, parent effort, size target
2. **Scope** - Exactly what to implement
3. **Exclusions** - What NOT to implement (handled elsewhere)
4. **Dependencies** - What this split needs from others
5. **Instructions** - Clear steps for SW Engineer
6. **Success Criteria** - Checklist for completion

## Integration with R199 and R204

### R199 Compliance - Single Reviewer
As the sole code reviewer for this effort:
```bash
echo "═══════════════════════════════════════════════════════"
echo "SPLIT PLANNING ASSIGNMENT CONFIRMATION"
echo "═══════════════════════════════════════════════════════"
echo "✅ I am the SOLE split planner per R199"
echo "✅ I will create ALL split plans"
echo "✅ Plans will be saved in branch: $CURRENT_BRANCH"
echo "✅ Orchestrator will copy to split directories per R204"
```

### R204 Integration - Orchestrator Will Use These Plans
The orchestrator will:
1. Read SPLIT-INVENTORY.md from too-large branch
2. Create split directories
3. Copy appropriate SPLIT-PLAN-XXX.md to each split directory
4. Spawn SW Engineer with the split plan

## Example Split Plan Creation

```bash
# Real example for api-types effort that's 1247 lines
create_api_types_split_plans() {
    # Create inventory
    cat > SPLIT-INVENTORY.md << 'EOF'
# Split Inventory for api-types

## Overview
api-types effort exceeded 800 lines (actual: 1247) and requires splitting.

## Split Structure
| Split # | Name | Description | Est. Lines | Status |
|---------|------|-------------|------------|--------|
| 001 | core-types | Core API type definitions | 420 | Planned |
| 002 | validators | Validation logic | 380 | Planned |
| 003 | converters | Type converters and helpers | 400 | Planned |

## File Distribution
### Split-001 (core-types)
- pkg/apis/v1alpha1/types.go
- pkg/apis/v1alpha1/register.go
- pkg/apis/v1alpha1/doc.go

### Split-002 (validators)
- pkg/apis/v1alpha1/validation.go
- pkg/apis/v1alpha1/webhook_validation.go

### Split-003 (converters)
- pkg/apis/v1alpha1/conversion.go
- pkg/apis/v1alpha1/helpers.go
EOF

    # Create individual plans
    for split in 001 002 003; do
        create_individual_split_plan $split
    done
    
    # Commit and push
    git add SPLIT-*.md
    git commit -m "feat: split api-types (1247 lines) into 3 compliant efforts"
    git push
}
```

## Common Patterns

### Sequential Dependency Pattern
```yaml
Split-001: Foundation (no dependencies)
Split-002: Extensions (depends on Split-001)
Split-003: Integration (depends on Split-001 and Split-002)
```

### Parallel Implementation Pattern
```yaml
Split-001: Component A (independent)
Split-002: Component B (independent)
Split-003: Integration Layer (depends on both)
```

### Layer-Based Pattern
```yaml
Split-001: Data Layer
Split-002: Business Logic Layer
Split-003: API Layer
```

## Validation Before Committing

```bash
validate_split_plans() {
    echo "Validating split plans..."
    
    # Check inventory exists
    if [ ! -f "SPLIT-INVENTORY.md" ]; then
        echo "❌ Missing SPLIT-INVENTORY.md"
        return 1
    fi
    
    # Check all referenced plans exist
    PLAN_COUNT=$(grep -c "^| [0-9]" SPLIT-INVENTORY.md)
    ACTUAL_PLANS=$(ls SPLIT-PLAN-*.md 2>/dev/null | wc -l)
    
    if [ "$PLAN_COUNT" != "$ACTUAL_PLANS" ]; then
        echo "❌ Mismatch: Inventory lists $PLAN_COUNT splits but found $ACTUAL_PLANS plan files"
        return 1
    fi
    
    # Verify each plan has required sections
    for plan in SPLIT-PLAN-*.md; do
        if ! grep -q "## Split Metadata" "$plan"; then
            echo "❌ $plan missing Split Metadata section"
            return 1
        fi
        if ! grep -q "## Implementation Scope" "$plan"; then
            echo "❌ $plan missing Implementation Scope section"
            return 1
        fi
        if ! grep -q "## Success Criteria" "$plan"; then
            echo "❌ $plan missing Success Criteria section"
            return 1
        fi
    done
    
    echo "✅ All split plans validated"
    return 0
}
```

## State Transition

After creating and committing split plans:
1. Verify all plans committed to too-large branch
2. Confirm plans are accessible from remote
3. Report completion to orchestrator
4. Transition to COMPLETED state

The orchestrator will then:
1. Fetch split plans from too-large branch
2. Create split effort directories
3. Copy plans to appropriate directories
4. Spawn SW Engineers for implementation

## Remember

- **ALWAYS** create plans in the too-large branch
- **ALWAYS** commit and push to remote
- **ALWAYS** include SPLIT-INVENTORY.md
- **NEVER** create split directories (orchestrator does this)
- **NEVER** switch branches (stay in too-large branch)
- **ALWAYS** validate plans before committing