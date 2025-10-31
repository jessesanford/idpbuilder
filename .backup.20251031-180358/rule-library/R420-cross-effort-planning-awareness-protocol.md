# 🚨🚨🚨 RULE R420: Cross-Effort Planning Awareness Protocol (BLOCKING) 🚨🚨🚨

## Category
BLOCKING - Planning Phase Requirement

## Criticality
🚨🚨🚨 BLOCKING - Violation = -50% to -100% grading penalty

## Description
Code Reviewers MUST perform comprehensive cross-effort awareness analysis BEFORE creating any effort or split plans. This analysis MUST discover and document ALL previous implementations, plans, file structures, APIs, and symbols to prevent duplicate declarations, API mismatches, and method visibility errors.

## Rationale
The integration failures that motivated this rule demonstrated that **general principles (R373/R374) are insufficient without STATE-SPECIFIC enforcement**:

### Real Failure Examples:
1. **Duplicate PushCmd Declaration** - Two files declared same symbol
2. **Retry Package API Mismatch** - Assumed DefaultBackoff existed, used wrong Config type
3. **MockRegistry Method Visibility** - Accessed unexported methods from external package

**Root Cause:** Planning happened WITHOUT reading actual implementations from previous efforts/splits.

## Relationship to Existing Rules

### R373 (Mandatory Code Reuse and Interface Compliance)
- **R373 defines WHAT**: NO duplicate implementations allowed
- **R420 defines HOW**: The process to discover what exists
- **Together**: R420 research → R373 compliance verification

### R374 (Pre-Planning Research Protocol)
- **R374 provides GENERAL framework**: Research existing code before planning
- **R420 provides SPECIFIC enforcement**: Exact steps, validation, state integration
- **Together**: R374 philosophy → R420 execution protocol

### R420's Unique Contribution
R420 is the **STATE-SPECIFIC ENFORCEMENT MECHANISM** that makes R373 and R374 actually work:
- Integrates into EFFORT_PLAN_CREATION and CREATE_SPLIT_PLAN states
- Provides step-by-step mandatory procedures
- Includes validation gates and blocking checks
- Defines specific grading penalties for incomplete research

## Requirements

### 🔴🔴🔴 PHASE 1: DISCOVERY (MANDATORY - BEFORE ANY PLANNING)

#### 1.1 Discover All Previous Efforts
```bash
# In orchestrator state file - READ efforts_completed
PREVIOUS_EFFORTS=$(jq -r '.efforts_completed[] | .effort_id' orchestrator-state-v3.json)

echo "📋 Discovering previous efforts..."
for effort_id in $PREVIOUS_EFFORTS; do
    echo "  - Found: $effort_id"
    # Track for detailed analysis
done
```

#### 1.2 Discover All Sibling Splits (for split planning)
```bash
# For CREATE_SPLIT_PLAN state only
CURRENT_EFFORT=$(jq -r '.efforts_in_progress[0].effort_id' orchestrator-state-v3.json)

# Find all splits for this effort
PREVIOUS_SPLITS=$(jq -r ".efforts_in_progress[] | select(.effort_id == \"$CURRENT_EFFORT\") | .splits_completed[]?" orchestrator-state-v3.json)

echo "📋 Discovering sibling splits..."
for split_id in $PREVIOUS_SPLITS; do
    echo "  - Found split: $split_id"
done
```

#### 1.3 Discover All Previous Plans
```bash
# Find all implementation plans
PREVIOUS_PLANS=$(find /path/to/workspace -name "*IMPLEMENTATION-PLAN*.md" 2>/dev/null)

echo "📋 Discovering previous plans..."
echo "$PREVIOUS_PLANS" | while read plan; do
    echo "  - Found plan: $plan"
done
```

### 🔴🔴🔴 PHASE 2: READ IMPLEMENTATIONS (MANDATORY)

#### 2.1 Read Code from Previous Efforts
```bash
# For EACH previous effort
for effort_id in $PREVIOUS_EFFORTS; do
    echo "📖 Reading implementation: $effort_id"

    # Get effort directory from state
    effort_dir=$(jq -r ".efforts_completed[] | select(.effort_id == \"$effort_id\") | .effort_directory" orchestrator-state-v3.json)

    if [ -d "$effort_dir" ]; then
        # Find exported symbols
        echo "  → Exported symbols:"
        grep -r "^func [A-Z]" --include="*.go" "$effort_dir" 2>/dev/null | head -10

        # Find interfaces
        echo "  → Interfaces:"
        grep -r "^type.*interface" --include="*.go" "$effort_dir" 2>/dev/null

        # Find structs
        echo "  → Structs:"
        grep -r "^type [A-Z].*struct" --include="*.go" "$effort_dir" 2>/dev/null | head -10

        # Find package structure
        echo "  → Packages:"
        find "$effort_dir" -type d -name "pkg" -o -name "cmd" 2>/dev/null
    fi
done
```

#### 2.2 Read Code from Sibling Splits (for split planning)
```bash
# For CREATE_SPLIT_PLAN only
for split_id in $PREVIOUS_SPLITS; do
    echo "📖 Reading split implementation: $split_id"

    split_dir=$(jq -r ".efforts_in_progress[0].splits_completed[] | select(.split_id == \"$split_id\") | .split_directory" orchestrator-state-v3.json)

    if [ -d "$split_dir" ]; then
        # Check for package APIs
        echo "  → Package exports from $split_id:"
        grep -r "^func [A-Z]" --include="*.go" "$split_dir/pkg" 2>/dev/null

        # Check for type definitions
        echo "  → Type definitions:"
        grep -r "^type [A-Z]" --include="*.go" "$split_dir/pkg" 2>/dev/null

        # Check method visibility
        echo "  → Exported methods (uppercase):"
        grep -r "func ([a-z]* \*[A-Z][a-zA-Z]*) [A-Z]" --include="*.go" "$split_dir" 2>/dev/null
    fi
done
```

### 🔴🔴🔴 PHASE 3: READ PLANS (MANDATORY)

#### 3.1 Read Plans for Unimplemented Efforts
```bash
# Find plans that exist but aren't implemented yet
for plan_file in $PREVIOUS_PLANS; do
    echo "📖 Reading plan: $plan_file"

    # Extract planned symbols
    echo "  → Planned interfaces:"
    grep -A5 "interface" "$plan_file" 2>/dev/null | head -20

    # Extract planned packages
    echo "  → Planned packages:"
    grep "pkg/" "$plan_file" 2>/dev/null | head -10

    # Extract planned APIs
    echo "  → Planned APIs:"
    grep -E "func|method" "$plan_file" 2>/dev/null | head -10
done
```

### 🔴🔴🔴 PHASE 4: ANALYZE FOR CONFLICTS (MANDATORY)

#### 4.1 Detect Duplicate File Paths
```bash
echo "🔍 Checking for file structure conflicts..."

# Build list of files from previous efforts
EXISTING_FILES=$(mktemp)
for effort_id in $PREVIOUS_EFFORTS; do
    effort_dir=$(jq -r ".efforts_completed[] | select(.effort_id == \"$effort_id\") | .effort_directory" orchestrator-state-v3.json)
    find "$effort_dir" -type f -name "*.go" 2>/dev/null >> $EXISTING_FILES
done

# Check if current plan would create duplicates
# (This happens during plan review)
if grep -q "pkg/cmd/push/root.go" $EXISTING_FILES; then
    echo "⚠️ WARNING: pkg/cmd/push/root.go already exists!"
    echo "   Current plan MUST NOT create duplicate file"
fi

rm $EXISTING_FILES
```

#### 4.2 Detect API Mismatches
```bash
echo "🔍 Checking for API compatibility..."

# Example: Check retry package API
if grep -q "retry.DefaultBackoff" <current-plan>; then
    # Verify it exists in previous implementation
    if ! grep -r "DefaultBackoff" --include="*.go" $PREVIOUS_SPLIT_DIR/pkg/retry/ 2>/dev/null; then
        echo "❌ ERROR: Plan assumes retry.DefaultBackoff exists"
        echo "   Actual API from split-001 uses retry.Config instead"
        echo "   Plan MUST be updated to use correct API"
        exit 1
    fi
fi
```

#### 4.3 Detect Method Visibility Errors
```bash
echo "🔍 Checking method visibility..."

# Example: Check MockRegistry methods
if grep -q "MockRegistry.*validate()" <current-plan>; then
    # Check if validate is exported (uppercase)
    if ! grep -r "func (.*MockRegistry) Validate(" --include="*.go" $PREVIOUS_EFFORT_DIR 2>/dev/null; then
        echo "❌ ERROR: Plan tries to access unexported method validate()"
        echo "   Method is lowercase = private/unexported"
        echo "   Plan MUST use only exported methods (uppercase)"
        exit 1
    fi
fi
```

### 🔴🔴🔴 PHASE 5: DOCUMENT IN PLAN (MANDATORY)

Every effort plan and split plan MUST include this section:

```markdown
## 🔍 PRIOR WORK ANALYSIS (R420 MANDATORY)

### Discovery Phase Results
- **Previous Efforts Reviewed**: [List effort IDs: E1.1.0, E1.1.1, etc.]
- **Sibling Splits Reviewed**: [For splits only: split-001, split-002, etc.]
- **Previous Plans Reviewed**: [List plan files examined]
- **Research Timestamp**: [When this research was conducted]

### File Structure Findings
| File Path | Source | Status | Action Required |
|-----------|--------|--------|-----------------|
| pkg/cmd/push/root.go | E1.2.0 | EXISTS | MUST NOT create duplicate |
| pkg/cmd/push/push.go | E1.2.1 (this effort) | NEW | Safe to create |
| pkg/retry/retry.go | split-001 | EXISTS | MUST import and use |

### Interface/API Findings
| Interface/API | Source | Signature | Action Required |
|---------------|--------|-----------|-----------------|
| Registry | E1.1.0 | Push(ctx, image string, content io.Reader) error | MUST implement exactly |
| retry.Config | split-001 | struct{MaxRetries int, BackoffDuration time.Duration} | MUST use, NOT DefaultBackoff |

### Type/Struct Findings
| Type | Source | Exported | Action Required |
|------|--------|----------|-----------------|
| GiteaClient | E1.2.0 | YES (public) | Can import and extend |
| MockRegistry | E1.1.0 | YES (public) | Can use in tests |

### Method Visibility Findings
| Method | Type | Visibility | Can Access? | Action Required |
|--------|------|------------|-------------|-----------------|
| Push() | MockRegistry | EXPORTED | YES | Safe to use in assertions |
| Pull() | MockRegistry | EXPORTED | YES | Safe to use in assertions |
| validate() | MockRegistry | UNEXPORTED | NO | CANNOT use, private method |

### Package Organization Findings
| Package | Source | Purpose | Action Required |
|---------|--------|---------|-----------------|
| pkg/cmd/push | E1.2.0 | Command hierarchy | MUST add to existing structure |
| pkg/retry | split-001 | Retry logic | MUST import, NOT recreate |

### Conflicts Detected
- ✅ NO duplicate file paths detected
- ✅ NO API mismatches detected
- ✅ NO method visibility violations detected

### Required Integrations
1. MUST import Registry interface from E1.1.0
2. MUST use retry.Config from split-001 (NOT DefaultBackoff)
3. MUST extend pkg/cmd/push structure (NOT create new root.go)
4. MUST use only exported MockRegistry methods in tests

### Forbidden Actions
- ❌ DO NOT create pkg/cmd/push/root.go (duplicate)
- ❌ DO NOT use retry.DefaultBackoff (doesn't exist)
- ❌ DO NOT access MockRegistry.validate() (private method)
- ❌ DO NOT create alternative Registry interface (R373 violation)
```

## Enforcement Mechanisms

### 🚨🚨🚨 BLOCKING Gate: Plan Cannot Be Approved Without Research

```bash
#!/bin/bash
# R420-validate-plan-research.sh

PLAN_FILE=$1

echo "Validating R420 compliance in plan..."

# Check for mandatory section
if ! grep -q "PRIOR WORK ANALYSIS (R420 MANDATORY)" "$PLAN_FILE"; then
    echo "❌ BLOCKING: Missing mandatory 'PRIOR WORK ANALYSIS (R420 MANDATORY)' section"
    echo "   Plan cannot be approved without R420 research"
    exit 1
fi

# Check for discovery results
if ! grep -q "Discovery Phase Results" "$PLAN_FILE"; then
    echo "❌ BLOCKING: Missing discovery phase documentation"
    exit 1
fi

# Check for findings tables
for section in "File Structure Findings" "Interface/API Findings" "Method Visibility Findings"; do
    if ! grep -q "$section" "$PLAN_FILE"; then
        echo "❌ BLOCKING: Missing required section: $section"
        exit 1
    fi
done

# Check that research was substantive (not empty)
FINDINGS_COUNT=$(grep -c "| .* | .* |" "$PLAN_FILE")
if [ $FINDINGS_COUNT -lt 3 ]; then
    echo "❌ BLOCKING: Research appears incomplete (too few findings documented)"
    echo "   Found $FINDINGS_COUNT finding rows, expected at least 3"
    exit 1
fi

echo "✅ Plan passes R420 validation"
exit 0
```

### Integration with Code Reviewer States

#### EFFORT_PLAN_CREATION State
**MANDATORY Pre-Flight Check:**
```bash
# BEFORE creating any effort plan
echo "🚨 R420 Pre-Planning Research Required..."

# Execute Phase 1: Discovery
# Execute Phase 2: Read Implementations
# Execute Phase 3: Read Plans
# Execute Phase 4: Analyze Conflicts
# Execute Phase 5: Document in Plan

# Validate plan before proceeding
bash tools/validate-prior-work-research.sh /path/to/plan.md || exit 1
```

#### CREATE_SPLIT_PLAN State
**MANDATORY Sibling Split Analysis:**
```bash
# BEFORE creating any split plan
echo "🚨 R420 Sibling Split Research Required..."

# Focus on previous splits in same effort
# Execute Phase 1.2: Discover Sibling Splits
# Execute Phase 2.2: Read Split Implementations
# Execute Phase 4: Analyze Conflicts (especially API compatibility)
# Execute Phase 5: Document in Plan

# Validate plan before proceeding
bash tools/validate-prior-work-research.sh /path/to/split-plan.md || exit 1
```

## Grading Penalties

### Missing Research
- **No discovery phase**: -50% (CRITICAL - shows complete neglect of R420)
- **No implementation reading**: -40% (MAJOR - planning blind to reality)
- **No plan reading**: -30% (SIGNIFICANT - ignoring future work)
- **No conflict analysis**: -50% (CRITICAL - recipe for integration failure)
- **No documentation in plan**: -40% (MAJOR - hiding incomplete work)

### Incomplete Research
- **Skipped previous efforts**: -30% (per skipped effort)
- **Skipped sibling splits**: -40% (per skipped split - critical for splits)
- **Missing findings tables**: -20% (per missing table)
- **Empty/trivial research**: -35% (shows pro forma compliance)

### Research Failures Leading to Integration Errors
- **Duplicate file declaration**: -100% (R420 + R373 violation)
- **API mismatch**: -50% (R420 failure to verify APIs)
- **Method visibility error**: -30% (R420 failure to check exports)
- **Package conflict**: -40% (R420 failure to analyze structure)

### Maximum Cumulative Penalty
- **Total possible penalty**: -100% (IMMEDIATE FAILURE)
- **Typical penalty for skipped research**: -50% to -75%
- **Penalty for research leading to duplicates**: -100% (via R373)

## Success Metrics

### Effort Plan Approval Requires:
1. ✅ R420 section present in plan
2. ✅ All previous efforts listed and reviewed
3. ✅ All findings tables populated
4. ✅ Conflicts section shows analysis occurred
5. ✅ Validation script passes

### Split Plan Approval Requires:
1. ✅ R420 section present in plan
2. ✅ All sibling splits listed and reviewed
3. ✅ API compatibility verified
4. ✅ Method visibility checked
5. ✅ Validation script passes

### Integration Success Requires:
1. ✅ Zero duplicate declarations
2. ✅ Zero API mismatches
3. ✅ Zero method visibility errors
4. ✅ Zero file structure conflicts

## Example: Correct R420 Execution

### Scenario: Planning E1.2.1 (push command implementation)

#### Phase 1: Discovery
```bash
$ jq -r '.efforts_completed[] | .effort_id' orchestrator-state-v3.json
E1.1.0
E1.2.0

$ echo "Found 2 previous efforts to review"
```

#### Phase 2: Read E1.2.0 Implementation
```bash
$ grep -r "^func [A-Z]" E1.2.0/pkg/cmd/push/
E1.2.0/pkg/cmd/push/root.go:func NewPushCmd() *cobra.Command
E1.2.0/pkg/cmd/push/root.go:var PushCmd = &cobra.Command{...}

$ echo "FOUND: PushCmd already declared in root.go"
```

#### Phase 3: Analyze Current Plan
```bash
$ grep "root.go" E1.2.1-IMPLEMENTATION-PLAN.md
- Create pkg/cmd/push/root.go with PushCmd

$ echo "CONFLICT DETECTED: Plan tries to create duplicate root.go!"
```

#### Phase 4: Document in Plan
```markdown
## 🔍 PRIOR WORK ANALYSIS (R420 MANDATORY)

### File Structure Findings
| File Path | Source | Status | Action Required |
|-----------|--------|--------|-----------------|
| pkg/cmd/push/root.go | E1.2.0 | EXISTS with PushCmd | MUST NOT create duplicate |

### Required Actions:
- E1.2.1 MUST add subcommands to EXISTING root.go from E1.2.0
- E1.2.1 MUST NOT declare new PushCmd variable
- E1.2.1 creates ONLY new subcommand files (e.g., pkg/cmd/push/container.go)

### Forbidden Actions:
- ❌ DO NOT create pkg/cmd/push/root.go (already exists in E1.2.0)
```

**Result:** Integration failure prevented! E1.2.1 extends existing structure instead of creating duplicates.

## Integration with Other Rules

### R373 (Mandatory Code Reuse and Interface Compliance)
- **R420 discovers** what interfaces/APIs exist
- **R373 enforces** that they are reused, not duplicated
- **Together**: Complete prevention of duplicate implementations

### R374 (Pre-Planning Research Protocol)
- **R374 provides** philosophical foundation and general framework
- **R420 provides** specific enforcement mechanisms and state integration
- **Together**: Principle + enforcement = actual compliance

### R362 (No Architectural Rewrites Without Approval)
- **R420 discovers** existing architecture decisions
- **R362 prevents** changing them without approval
- **Together**: Preserve architectural coherence

### R009 (Mandatory Wave/Phase Integration Protocol)
- **R420 prevents** conflicts before they happen
- **R009 validates** that no conflicts exist during integration
- **Together**: Proactive prevention + reactive validation

## State Machine Integration

### States Where R420 is MANDATORY:
1. **EFFORT_PLAN_CREATION** (Code Reviewer)
   - MUST execute R420 before creating plan
   - MUST include R420 section in plan
   - MUST pass R420 validation to approve plan

2. **CREATE_SPLIT_PLAN** (Code Reviewer)
   - MUST execute R420 sibling split analysis
   - MUST verify API compatibility with previous splits
   - MUST document method visibility findings

### States Where R420 is VALIDATED:
1. **CODE_REVIEW** (Code Reviewer)
   - MUST verify implementation matches R420 research
   - MUST confirm no forbidden duplications created
   - MUST check API compatibility

2. **REVIEW_WAVE_ARCHITECTURE** (Architect)
   - MUST verify all efforts performed R420 research
   - MUST confirm no cross-effort conflicts exist

3. **INTEGRATE_WAVE_EFFORTS** (Integration Agent)
   - MUST verify R420 prevented integration failures
   - MUST validate zero duplicate declarations

## Recovery Protocol

### If R420 Was Skipped:
1. **STOP all work immediately**
2. **Return to EFFORT_PLAN_CREATION or CREATE_SPLIT_PLAN state**
3. **Execute R420 protocol completely**
4. **Update plan with R420 findings**
5. **Revalidate plan passes R420 checks**
6. **Resume implementation ONLY after R420 complete**

### If Conflicts Discovered During Integration:
1. **STOP integration immediately**
2. **Trace back to which effort/split caused conflict**
3. **Execute R420 retroactively to find what was missed**
4. **Create fix plan based on R420 findings**
5. **Refactor to eliminate conflict**
6. **Re-integrate after fix**

## Validation Scripts

### Pre-Planning Research Executor
**File:** `tools/execute-R420-research.sh`
```bash
#!/bin/bash
# Automates R420 research protocol

ORCHESTRATOR_STATE="orchestrator-state-v3.json"
OUTPUT_FILE="R420-research-results.md"

echo "# R420 Prior Work Analysis" > $OUTPUT_FILE
echo "Generated: $(date)" >> $OUTPUT_FILE
echo "" >> $OUTPUT_FILE

# Phase 1: Discovery
echo "## Discovery Phase" >> $OUTPUT_FILE
PREVIOUS_EFFORTS=$(jq -r '.efforts_completed[] | .effort_id' $ORCHESTRATOR_STATE)
echo "Previous Efforts:" >> $OUTPUT_FILE
echo "$PREVIOUS_EFFORTS" | while read effort; do
    echo "- $effort" >> $OUTPUT_FILE
done

# Phase 2: Read Implementations
echo "" >> $OUTPUT_FILE
echo "## Implementation Analysis" >> $OUTPUT_FILE
echo "$PREVIOUS_EFFORTS" | while read effort_id; do
    effort_dir=$(jq -r ".efforts_completed[] | select(.effort_id == \"$effort_id\") | .effort_directory" $ORCHESTRATOR_STATE)

    echo "### $effort_id ($effort_dir)" >> $OUTPUT_FILE

    # Exported functions
    echo "#### Exported Functions:" >> $OUTPUT_FILE
    grep -r "^func [A-Z]" --include="*.go" "$effort_dir" 2>/dev/null | head -10 >> $OUTPUT_FILE || echo "None found" >> $OUTPUT_FILE

    # Interfaces
    echo "#### Interfaces:" >> $OUTPUT_FILE
    grep -r "^type.*interface" --include="*.go" "$effort_dir" 2>/dev/null >> $OUTPUT_FILE || echo "None found" >> $OUTPUT_FILE

    # Structs
    echo "#### Structs:" >> $OUTPUT_FILE
    grep -r "^type [A-Z].*struct" --include="*.go" "$effort_dir" 2>/dev/null | head -10 >> $OUTPUT_FILE || echo "None found" >> $OUTPUT_FILE

    echo "" >> $OUTPUT_FILE
done

echo "R420 research complete. Results in $OUTPUT_FILE"
echo "Copy relevant sections into your IMPLEMENTATION-PLAN.md under '🔍 PRIOR WORK ANALYSIS (R420 MANDATORY)'"
```

### Plan Validation Script
**File:** `tools/validate-R420-compliance.sh`
```bash
#!/bin/bash
# Validates plan includes required R420 research

PLAN_FILE=$1

if [ -z "$PLAN_FILE" ]; then
    echo "Usage: $0 <plan-file.md>"
    exit 1
fi

echo "Validating R420 compliance in $PLAN_FILE..."

# Check for mandatory section
if ! grep -q "PRIOR WORK ANALYSIS (R420 MANDATORY)" "$PLAN_FILE"; then
    echo "❌ BLOCKING: Missing 'PRIOR WORK ANALYSIS (R420 MANDATORY)' section"
    exit 1
fi

# Check for required subsections
REQUIRED_SECTIONS=(
    "Discovery Phase Results"
    "File Structure Findings"
    "Interface/API Findings"
    "Conflicts Detected"
    "Forbidden Actions"
)

for section in "${REQUIRED_SECTIONS[@]}"; do
    if ! grep -q "$section" "$PLAN_FILE"; then
        echo "❌ BLOCKING: Missing required section: $section"
        exit 1
    fi
done

# Check substantiveness
FINDINGS_COUNT=$(grep -c "| .* | .* |" "$PLAN_FILE")
if [ $FINDINGS_COUNT -lt 3 ]; then
    echo "⚠️ WARNING: Research may be incomplete (only $FINDINGS_COUNT finding rows)"
fi

echo "✅ Plan passes R420 validation"
exit 0
```

## Tags
#blocking #planning #cross-effort-awareness #r420 #research #state-enforcement #integration-failure-prevention
