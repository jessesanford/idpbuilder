# Code Reviewer - CREATE_SPLIT_PLAN State Rules

## State Context
You are creating individual split plan files for an oversized implementation (>800 lines). These plans will be saved in the too-large branch and later copied to split effort directories by the orchestrator.

## 🔴🔴🔴 CRITICAL: YOU MUST MEASURE - R319 DOES NOT APPLY TO YOU! 🔴🔴🔴
**R319 (Orchestrator Never Measures) applies ONLY to orchestrators!**
**As a Code Reviewer creating split plans, you MUST measure to determine split sizes!**
**Use line-counter.sh tool to verify actual implementation size before planning splits!**

## 🔴🔴🔴 PARAMOUNT: Repository Location (R251 & R309) 🔴🔴🔴

### R251: Universal Repository Separation Law
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R251-REPOSITORY-SEPARATION-LAW.md`
**Criticality**: PARAMOUNT - Automatic -100% failure for violation
**YOU ARE IN**: TARGET REPOSITORY CLONE (the too-large effort)
**NOT IN**: Software Factory planning repo

### R309: Never Create Efforts in SF Repo
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R309-never-create-efforts-in-sf-repo.md`
**Criticality**: PARAMOUNT - Automatic -100% failure for violation
**CRITICAL**: Split plans are created in TARGET repo effort directories

**VERIFY CORRECT REPOSITORY:**
```bash
echo "🔴 R251/R309: Verifying we're in TARGET repo effort..."
if [ -f "orchestrator-state.yaml" ] || [ -f ".claude/CLAUDE.md" ]; then
    echo "🔴🔴🔴 FATAL: In Software Factory repo!"
    echo "Split plans must be created in TARGET effort directory!"
    exit 309
fi

if [[ "$(pwd)" != *"/efforts/"* ]]; then
    echo "🔴 WARNING: Not clearly in /efforts/ directory"
    echo "Verify you're in the too-large effort directory"
    pwd
fi

echo "✅ Creating split plans in TARGET repo effort directory"
```

## 🔴🔴🔴 CRITICAL: EXPLICIT BOUNDARIES TO PREVENT OVER-ENGINEERING 🔴🔴🔴

**ACTUAL FAILURE FROM TRANSCRIPT: 3.4X OVERRUN DUE TO VAGUE LANGUAGE**

### The Actual Problem That Just Happened:
```
Vague Plan: "Implement selected command files and utilities"
SW Result: Implemented ALL commands, complete utility modules, full tests
Outcome: 2,215 lines instead of 650 (3.4X OVERRUN) = COMPLETE FAILURE

Root Cause Analysis:
- "selected" was interpreted as "all important ones"
- "utilities" was interpreted as "complete utility module"
- "handlers" was interpreted as "comprehensive handler system"
- No test count specified, so engineer added 430 lines of tests
```

### Your Solution: ELIMINATE ALL AMBIGUOUS WORDS
```
Explicit Plan: "Implement EXACTLY these 3 files:
1. cmd/push.go - ONLY HandlePush() function (200 lines)
2. pkg/push/validator.go - ONLY ValidatePushArgs() (150 lines)
3. cmd/push_test.go - EXACTLY 2 tests (100 lines)
DO NOT implement other command files."

SW Result: Exactly 3 files, 2 functions, 2 tests
Outcome: 450 lines as planned = SUCCESS
```

### FORBIDDEN WORDS IN SPLIT PLANS:
- ❌ "selected" - Say which ones exactly
- ❌ "necessary" - Define what's necessary
- ❌ "utilities" - Name each utility function
- ❌ "handlers" - List each handler by name
- ❌ "components" - Specify exact components
- ❌ "some" - Give exact count
- ❌ "various" - List them all
- ❌ "etc." - Never use this

## 🔴🔴🔴 CRITICAL: Sequential Branching Strategy (Within R308 Incremental) 🔴🔴🔴

**ALL SPLITS MUST BE PLANNED FOR SEQUENTIAL BRANCHING!**

### Understanding Split Branching vs Incremental Branching:
- **R308 Incremental**: Efforts branch from previous wave/phase integration
- **Split Sequential**: Splits branch from each other sequentially
- **The Base**: Original effort's base is determined by R308

### The Mandatory Pattern:
```
Original effort base (per R308): phase1-wave2-integration (example)
    ↓
Split-001: Based on same incremental base as original
    ↓ (becomes base for next)
Split-002: Based on Split-001 (NOT the integration!)
    ↓ (becomes base for next)
Split-003: Based on Split-002 (NOT the integration!)
```

### Why This is CRITICAL:
1. **Line Counting**: Each split measures ONLY its additions (400 lines, not cumulative)
2. **Dependencies**: Later splits can use earlier split code
3. **Clean Integration**: No merge conflicts between splits
4. **Progressive Building**: Each split extends the previous

### 🚨 CRITICAL MEASUREMENT - TOOL AUTO-DETECTS! 🚨
**THE TOOL NOW AUTOMATICALLY MEASURES AGAINST THE CORRECT BASE!**

```bash
# ✅ NEW TOOL - Automatically detects split predecessors:
./tools/line-counter.sh split-003
# Tool output: 🎯 Detected base: split-002
# CORRECT: Shows 280 lines (ONLY split-003's incremental work)

# ❌ OLD PROBLEM (NOW IMPOSSIBLE):
# Manual base selection led to errors like:
# ./tools/line-counter.sh -b main -c split-003  # OLD SYNTAX
# Would show 5,584 lines (included ALL splits!)
# THIS ERROR CAN NO LONGER HAPPEN!
```

**If reviewers measure against main, they will:**
- See inflated line counts (3-10X actual size)
- Reject valid splits as "violations"
- Cause unnecessary re-splitting
- Create cascading failures

**ALWAYS DOCUMENT THE MEASUREMENT BASE IN SPLIT PLANS!**

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

### Step 1: Verify You're in the Too-Large Branch and Create .software-factory Structure
```bash
echo "═══════════════════════════════════════════════════════"
echo "🔍 VERIFYING TOO-LARGE BRANCH CONTEXT"
echo "═══════════════════════════════════════════════════════"

# Verify current directory and branch
CURRENT_DIR=$(pwd)
CURRENT_BRANCH=$(git branch --show-current)
EFFORT_NAME=$(basename "$CURRENT_DIR")

echo "Current directory: $CURRENT_DIR"
echo "Current branch: $CURRENT_BRANCH"
echo "Effort name: $EFFORT_NAME"

# Extract phase and wave from path
if [[ "$CURRENT_DIR" =~ phase([0-9]+)/wave([0-9]+) ]]; then
    PHASE="${BASH_REMATCH[1]}"
    WAVE="${BASH_REMATCH[2]}"
    echo "Phase: $PHASE, Wave: $WAVE"
else
    echo "⚠️ WARNING: Cannot extract phase/wave from path"
    PHASE="X"
    WAVE="Y"
fi

# Create .software-factory structure for split plans
PLAN_BASE_DIR=".software-factory/phase${PHASE}/wave${WAVE}"
mkdir -p "$PLAN_BASE_DIR"
echo "📁 Created plan directory structure: $PLAN_BASE_DIR"

# Verify this is the too-large branch that needs splitting
if [ ! -f "IMPLEMENTATION-PLAN.md" ] && [ ! -d ".software-factory" ]; then
    echo "❌ ERROR: No IMPLEMENTATION-PLAN.md or .software-factory found"
    exit 1
fi

# Check if we already have split plans in new location
if find .software-factory -name "SPLIT-PLAN-*.md" 2>/dev/null | head -1; then
    echo "⚠️ WARNING: Previous split plans exist in .software-factory:"
    find .software-factory -name "SPLIT-PLAN-*.md" -exec ls -la {} \;
    echo "New plans will be created with current timestamp: $(date '+%Y%m%d-%H%M%S')"
    echo "This prevents collisions per R301"
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

## Infrastructure Requirements
Each split requires:
- **Separate Directory**: ${EFFORT_NAME}-SPLIT-00Z format
- **Separate Clone**: Own git repository in each directory
- **Sequential Branches**: Each based on previous split

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

### Step 3: Create Individual Split Plan Files in .software-factory Structure

**MANDATORY: Use timestamps per R301 to prevent collisions!**
```bash
# CRITICAL: Copy template to ensure explicit boundaries (if it exists)
if [ -f "$CLAUDE_PROJECT_DIR/templates/split-plan.md" ]; then
    cp "$CLAUDE_PROJECT_DIR/templates/split-plan.md" SPLIT-PLAN-TEMPLATE.md
fi

# Generate timestamp for this split planning session
TIMESTAMP=$(date '+%Y%m%d-%H%M%S')

# For each split, create timestamped plan in .software-factory subdirectory
for SPLIT_NUM in 001 002 003; do
    # Create subdirectory for this split
    SPLIT_DIR="${PLAN_BASE_DIR}/${EFFORT_NAME}-split-${SPLIT_NUM}"
    mkdir -p "$SPLIT_DIR"
    
    # Create filename with timestamp in the split's subdirectory
    PLAN_FILE="${SPLIT_DIR}/SPLIT-PLAN-${TIMESTAMP}.md"
    cat > "${PLAN_FILE}" << 'EOF'
# Split Plan ${SPLIT_NUM} - ${SPLIT_NAME}

## 🔴 MINDSET REMINDER FOR SW ENGINEER 🔴
**This is a PARTIAL implementation. It SHOULD feel incomplete.**
**Your job is to stay within budget, NOT to complete the feature.**
**If you think "this needs X to work properly" - STOP. X is probably in another split.**

## Split Metadata
- **Split Number**: ${SPLIT_NUM}
- **Parent Effort**: ${PARENT_EFFORT_NAME}
- **Original Branch**: ${TOO_LARGE_BRANCH}
- **Target Size**: ${TARGET_SIZE} lines (max 800)
- **Created**: $(date '+%Y-%m-%d %H:%M:%S')

## Infrastructure Requirements
- **Directory Name**: ${PARENT_EFFORT_NAME}-SPLIT-${SPLIT_NUM}
- **Location**: efforts/phase${PHASE}/wave${WAVE}/${PARENT_EFFORT_NAME}-SPLIT-${SPLIT_NUM}/
- **Clone Required**: Yes - separate clone of target repository
- **Branch Base**: $([ ${SPLIT_NUM} = "001" ] && echo "Same as original (e.g., phase-integration)" || echo "Previous split branch (split-$(printf "%03d" $((10#${SPLIT_NUM} - 1))))")

## Implementation Scope

### Files to Create/Modify
${FILE_LIST_FOR_SPLIT}

### 🚨 EXPLICIT SCOPE DEFINITION (MANDATORY PER R310)

#### MINIMUM VIABLE SCOPE (Your Exact Contract)

**FILES TO CREATE/MODIFY (COMPLETE LIST):**
```
1. path/to/file1.go (CREATE) - 200 lines MAX
2. path/to/file2.go (MODIFY) - Add 100 lines MAX
3. path/to/file1_test.go (CREATE) - 100 lines MAX
TOTAL: 3 files (implementing any other file = AUTOMATIC FAILURE)
```

**FUNCTIONS TO IMPLEMENT (BY EXACT SIGNATURE):**
```go
// In file1.go:
1. func HandlePush(args []string) error  // 150 lines
2. func validateImageRef(ref string) error // 30 lines

// In file2.go:
3. func GetPushOptions() *PushOptions // 50 lines

// TOTAL: 3 functions, 230 lines
// DO NOT ADD: Helper functions, utility functions, convenience wrappers
```

**METHODS TO IMPLEMENT (EXACT COUNT):**
```
NONE - Zero methods in this split
(OR list exact method signatures if needed)
```

**TESTS TO WRITE (EXACT COUNT OR ZERO):**
```
EXACTLY 2 tests in file1_test.go:
1. TestHandlePushValidArgs - 50 lines max
2. TestHandlePushInvalidArgs - 50 lines max
DO NOT ADD: Edge cases, table tests, benchmarks, examples
```

#### EXACTLY What Types/Structs to Create:
```
type StructName struct {  // EXACTLY N fields, NO methods in this split
    Field1 Type
    Field2 Type
    [LIST EVERY FIELD]
}
```

### 🛑 STOP BOUNDARIES - DO NOT IMPLEMENT
**EXPLICITLY FORBIDDEN IN THIS SPLIT:**
- ❌ DO NOT add validation methods
- ❌ DO NOT add Clone() or Copy() methods  
- ❌ DO NOT add converters or transformers
- ❌ DO NOT implement [specific things to avoid]
- ❌ DO NOT write comprehensive tests (basic only)
- ❌ DO NOT refactor existing code
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

## Size Management with REALISTIC Calculations
- Target: ${TARGET_SIZE} lines (MAX 400-600 to leave buffer)
- Hard Stop: 700 lines (better incomplete than oversized)
- Measurement: Use line-counter.sh before committing

### Realistic Line Estimates (MANDATORY):
```
Go function with validation: 30-50 lines
Simple option function: 10-20 lines
Struct definition: 5-10 lines per 5 fields
Basic test: 20-30 lines
Comprehensive test: 50-100 lines (AVOID)

This Split Calculation:
- N functions × Y lines = X lines
- N structs × Y lines = X lines
- N basic tests × Y lines = X lines
TOTAL: [must be under 600]
```

## Success Criteria
- [ ] All specified files implemented
- [ ] Size under 800 lines (measured)
- [ ] Tests passing
- [ ] Integrates with previous splits (if applicable)
- [ ] No functionality regression

## Notes for SW Engineer
${SPECIAL_NOTES}

### 🔴 ADHERENCE REMINDER (R310):
- Implement EXACTLY what's listed - no more
- If it seems incomplete, that's intentional
- Do NOT add "helpful" extras
- Do NOT "complete" the implementation
- STOP at the boundaries specified above
EOF
done

echo "✅ Created timestamped split plan files in .software-factory structure:"
find .software-factory -name "SPLIT-PLAN-${TIMESTAMP}.md" -exec ls -la {} \;
```

### Step 4: Commit and Push Split Plans to Too-Large Branch
```bash
echo "═══════════════════════════════════════════════════════"
echo "📦 COMMITTING SPLIT PLANS TO TOO-LARGE BRANCH"
echo "═══════════════════════════════════════════════════════"

# Add all split-related files in .software-factory structure
git add SPLIT-INVENTORY.md
git add ".software-factory/"

# Show what we're committing
echo "Files to commit:"
git status --short

# Commit with descriptive message
git commit -m "feat: add timestamped split plans for oversized effort

Created split plans with timestamp ${TIMESTAMP} to prevent collisions (R301):
- SPLIT-INVENTORY.md: Overview of all splits
- .software-factory/phase${PHASE}/wave${WAVE}/${EFFORT_NAME}-split-001/SPLIT-PLAN-${TIMESTAMP}.md
- .software-factory/phase${PHASE}/wave${WAVE}/${EFFORT_NAME}-split-002/SPLIT-PLAN-${TIMESTAMP}.md
- .software-factory/phase${PHASE}/wave${WAVE}/${EFFORT_NAME}-split-003/SPLIT-PLAN-${TIMESTAMP}.md

Plans are stored in .software-factory/ structure for organization.
This branch will be abandoned after splits are implemented and merged."

# Push to remote
echo "Pushing to remote..."
git push

echo "✅ Split plans committed and pushed to: $CURRENT_BRANCH"
echo "   Plans are stored in .software-factory/ subdirectories"
echo "   SW Engineers will find them at:"
for SPLIT_NUM in 001 002 003; do
    echo "     .software-factory/phase${PHASE}/wave${WAVE}/${EFFORT_NAME}-split-${SPLIT_NUM}/SPLIT-PLAN-*.md"
done
```

## Required Split Plan Elements

### SPLIT-INVENTORY.md Must Include:
1. **Overview** - Why the split was needed
2. **Split Structure** - Table of all splits with sizes
3. **File Distribution** - Which files go in which split
4. **Integration Strategy** - How splits will merge
5. **Dependencies** - Order of implementation

### Each SPLIT-PLAN-{effort}-split{num}-{timestamp}.md Must Include:
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
3. Copy appropriate timestamped split plan to each split directory
4. Spawn SW Engineer with the split plan

## Example Split Plan Creation

```bash
# Real example for api-types effort that's 1247 lines
create_api_types_split_plans() {
    local EFFORT_NAME="api-types"
    local TIMESTAMP=$(date '+%Y%m%d-%H%M%S')
    
    # Create inventory with timestamp reference
    cat > SPLIT-INVENTORY.md << EOF
# Split Inventory for ${EFFORT_NAME}

## Overview
${EFFORT_NAME} effort exceeded 800 lines (actual: 1247) and requires splitting.

## Metadata
- **Split Date**: $(date '+%Y-%m-%d %H:%M:%S')
- **Split Plans Timestamp**: ${TIMESTAMP}
- **Split By**: Code Reviewer Agent

## Split Structure
| Split # | Name | Description | Est. Lines | Plan File |
|---------|------|-------------|------------|----------|
| 001 | core-types | Core API type definitions | 420 | SPLIT-PLAN-${EFFORT_NAME}-split001-${TIMESTAMP}.md |
| 002 | validators | Validation logic | 380 | SPLIT-PLAN-${EFFORT_NAME}-split002-${TIMESTAMP}.md |
| 003 | converters | Type converters and helpers | 400 | SPLIT-PLAN-${EFFORT_NAME}-split003-${TIMESTAMP}.md |

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

    # Create individual timestamped plans
    for split in 001 002 003; do
        create_individual_split_plan "${EFFORT_NAME}" "${split}" "${TIMESTAMP}"
    done
    
    # Commit and push with timestamp reference
    git add SPLIT-INVENTORY.md
    git add SPLIT-PLAN-*-${TIMESTAMP}.md
    git commit -m "feat: split ${EFFORT_NAME} (1247 lines) into 3 compliant efforts

Split plans created with timestamp ${TIMESTAMP} per R301.
This prevents collisions if re-splitting is needed."
    git push
}
```

## 🔴🔴🔴 MANDATORY: Preventing Over-Engineering (R310 & R314) 🔴🔴🔴

### CRITICAL CHECKLIST - Every Split Plan MUST Have:

1. **EXACT FILE LIST** - Not "command files", but "cmd/push.go, cmd/build.go, cmd/tag.go"
2. **EXACT FUNCTION NAMES** - Not "options", but "WithImage(), WithContext(), WithPlatform()"
3. **EXACT LINE COUNTS PER FILE** - Not just total, but "cmd/push.go: 200 lines MAX"
4. **EXPLICIT TEST COUNT** - Not "write tests", but "3 tests" or "0 tests"
5. **COMPREHENSIVE FORBIDDEN LIST** - List everything they might add
6. **REALISTIC CALCULATIONS** - Show your math:
   ```
   3 functions × 40 lines = 120 lines
   1 struct × 25 lines = 25 lines
   3 tests × 30 lines = 90 lines
   TOTAL: 235 lines (well under 400)
   ```
7. **MINDSET GUIDANCE** - Explicitly state:
   ```
   MINDSET: This is a PARTIAL implementation.
   It SHOULD feel incomplete.
   DO NOT "complete" it.
   ```

### Examples Based on ACTUAL FAILURES:

❌ **ACTUAL BAD (From Transcript - Caused 3.4X overrun):**
```markdown
## Implementation
- Implement selected command files
- Add utilities and handlers as needed
- Create workflow components
```

✅ **CORRECTED VERSION (Would have prevented overrun):**
```markdown
## Implementation - EXACTLY 3 Files (450 lines MAX)

### FILES TO CREATE (NO OTHER FILES):
1. cmd/push.go (200 lines MAX)
   - ONLY: func HandlePush(args []string) error
   - DO NOT: Add other commands

2. pkg/push/validator.go (150 lines MAX)
   - ONLY: func ValidatePushArgs(args []string) error
   - DO NOT: Add other validators

3. cmd/push_test.go (100 lines MAX)
   - EXACTLY 2 tests:
     * TestHandlePushValidArgs (50 lines)
     * TestHandlePushInvalidArgs (50 lines)
   - DO NOT: Add edge cases, benchmarks, or examples

### FORBIDDEN (DO NOT IMPLEMENT):
- ❌ Other command files (build.go, tag.go, etc.)
- ❌ Workflow orchestration (allocated to Split-004)
- ❌ Comprehensive utility module (only validator needed)
- ❌ Integration tests (allocated to Split-005)
- ❌ Error handling beyond basic returns
```

❌ **BAD (No boundaries):**
"Complete the config.go file implementation"

✅ **GOOD (Clear boundaries):**
"Implement ONLY:
- Config struct with 5 fields (no methods)
- NewConfig() constructor (30 lines max)
DO NOT add methods, validation, or Clone()"

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