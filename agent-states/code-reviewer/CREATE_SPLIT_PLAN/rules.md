# Code Reviewer - CREATE_SPLIT_PLAN State Rules

## 🔴🔴🔴 ABSOLUTE REQUIREMENT: STATE MANAGER CONSULTATION 🔴🔴🔴

**SUPREME LAW - R517 - UNIVERSAL STATE MANAGER CONSULTATION LAW**

**BEFORE EXITING THIS STATE, YOU MUST:**

1. **SPAWN STATE MANAGER FOR SHUTDOWN_CONSULTATION** (MANDATORY - NO EXCEPTIONS)
2. **NEVER UPDATE STATE FILES DIRECTLY** (orchestrator-state-v3.json, bug-tracking.json, etc.)
3. **NEVER COMMIT STATE FILES YOURSELF** (State Manager does this atomically)

**FAILURE TO CONSULT STATE MANAGER = IMMEDIATE SYSTEM HALT (-100% GRADE)**

### Enforcement Mechanism

If you attempt to exit this state without spawning State Manager:
- ❌ Pre-commit hooks will REJECT your commit
- ❌ Validation tools will FAIL the build
- ❌ Grading system will assign -100% penalty
- ❌ System will transition to ERROR_RECOVERY

### Required Pattern (COPY THIS EXACTLY)

```bash
# At end of state work, BEFORE any state file updates:

echo "🔴 MANDATORY: Spawning State Manager for SHUTDOWN_CONSULTATION"

# Spawn State Manager (REQUIRED - NOT OPTIONAL)
# Task: state-manager
# State: SHUTDOWN_CONSULTATION
# Current State: [YOUR_CURRENT_STATE]
# Proposed Next State: [YOUR_PROPOSED_NEXT_STATE]
# Work Summary: [SUMMARY_OF_WORK_COMPLETED]

# State Manager will:
# 1. Validate proposed transition against state machine
# 2. Update all 4 state files atomically
# 3. Commit with [R288] tag
# 4. Return REQUIRED next state (may differ from proposal)

# Wait for State Manager response
# Follow State Manager's directive (REQUIRED next state)
# DO NOT proceed until State Manager confirms
```

**YOU MUST NEVER:**
- ❌ Update orchestrator-state-v3.json yourself
- ❌ Update bug-tracking.json yourself
- ❌ Update integration-containers.json yourself
- ❌ Use `jq` to modify state files
- ❌ Use `sed/awk` to modify state files
- ❌ Set `validated_by: "orchestrator"` (must be "state-manager")
- ❌ Skip State Manager consultation "just this once"
- ❌ Think "I'll validate it manually"

**ONLY State Manager may update state files. This is NON-NEGOTIABLE.**

See: `rule-library/R517-universal-state-manager-consultation-law.md`
## State Purpose
Create split plans when an effort exceeds size limits (>700 lines soft warning, >900 lines hard enforcement per R535). Split plans divide work into smaller, sequential splits while maintaining R420 cross-split awareness to prevent API mismatches and integration failures.

## 🚨🚨🚨 MANDATORY PRE-FLIGHT: R420 Sibling Split Analysis 🚨🚨🚨

**BEFORE creating ANY split plan, you MUST execute R420 Cross-Effort Planning Awareness Protocol with SPECIAL FOCUS on sibling splits:**

### Step 1: Discover Sibling Splits
```bash
# Read orchestrator state
ORCHESTRATOR_STATE="/path/to/orchestrator-state-v3.json"

# Get current effort
CURRENT_EFFORT=$(jq -r '.efforts_in_progress[0].effort_id' $ORCHESTRATOR_STATE)

# Discover all completed splits for this effort
PREVIOUS_SPLITS=$(jq -r ".efforts_in_progress[] | select(.effort_id == \"$CURRENT_EFFORT\") | .splits_completed[]?.split_id" $ORCHESTRATOR_STATE)

echo "📋 R420 Sibling Split Discovery:"
echo "Current effort: $CURRENT_EFFORT"
echo "Previous splits: $PREVIOUS_SPLITS"

# Count splits to review
SPLIT_COUNT=$(echo "$PREVIOUS_SPLITS" | grep -v '^$' | wc -l)
if [ $SPLIT_COUNT -gt 0 ]; then
    echo "MANDATORY: Must review $SPLIT_COUNT previous split(s)"
else
    echo "ℹ️ No previous splits (this will be split-001)"
fi
```

### Step 2: Read ALL Previous Split Implementations
```bash
# For EACH previous split (CRITICAL FOR API COMPATIBILITY)
echo "$PREVIOUS_SPLITS" | grep -v '^$' | while read split_id; do
    echo "📖 Reading split: $split_id (CRITICAL: API definitions!)"

    # Get split directory
    split_dir=$(jq -r ".efforts_in_progress[0].splits_completed[] | select(.split_id == \"$split_id\") | .split_directory" $ORCHESTRATOR_STATE)

    if [ -d "$split_dir" ]; then
        echo "  Analyzing $split_dir..."

        # CRITICAL: Find package exports (APIs for next splits)
        echo "  → Package exports (APIs this split provides):"
        grep -r "^func [A-Z]" --include="*.go" "$split_dir/pkg" 2>/dev/null

        # CRITICAL: Find type definitions
        echo "  → Type definitions:"
        grep -r "^type [A-Z]" --include="*.go" "$split_dir/pkg" 2>/dev/null

        # CRITICAL: Find exported methods
        echo "  → Exported methods (uppercase = public):"
        grep -r "func ([a-z]* \*[A-Z][a-zA-Z]*) [A-Z]" --include="*.go" "$split_dir" 2>/dev/null

        # CRITICAL: Find unexported methods (cannot be used)
        echo "  → Unexported methods (lowercase = private):"
        grep -r "func ([a-z]* \*[A-Z][a-zA-Z]*) [a-z]" --include="*.go" "$split_dir" 2>/dev/null

        # Check for specific packages created
        echo "  → Packages created:"
        find "$split_dir/pkg" -type d 2>/dev/null | sed "s|$split_dir/||"

    else
        echo "  ⚠️ WARNING: Split directory not found: $split_dir"
    fi
done
```

### Step 3: Analyze API Compatibility (CRITICAL FOR SPLITS!)
```bash
echo "🔍 R420 API Compatibility Analysis:"

# Example: If split-002 needs retry functionality
if [ -d "$SPLIT_001_DIR/pkg/retry" ]; then
    echo "  → split-001 provides retry package"

    # Check actual exports
    echo "  → Actual exports from split-001 retry package:"
    grep "^func [A-Z]\|^type [A-Z]" "$SPLIT_001_DIR/pkg/retry/retry.go" 2>/dev/null

    # CRITICAL: Document exact API available
    echo "  ✅ VERIFIED: split-002 can use these exports"
    echo "  ❌ FORBIDDEN: split-002 cannot assume other exports exist"
fi
```

### Step 4: Verify Method Visibility (CRITICAL FOR TEST CODE!)
```bash
echo "🔍 R420 Method Visibility Analysis:"

# Example: If split needs to test MockRegistry from previous split
if grep -q "MockRegistry" "$SPLIT_001_DIR/pkg/registry/mock.go" 2>/dev/null; then
    echo "  → Found MockRegistry in split-001"

    # Check which methods are exported (uppercase)
    echo "  → Exported methods (can use in tests):"
    grep "func ([a-z]* \*MockRegistry) [A-Z]" "$SPLIT_001_DIR/pkg/registry/mock.go" 2>/dev/null

    # Check which methods are unexported (CANNOT use)
    echo "  → Unexported methods (CANNOT use from other packages):"
    grep "func ([a-z]* \*MockRegistry) [a-z]" "$SPLIT_001_DIR/pkg/registry/mock.go" 2>/dev/null

    echo "  ✅ split-002 tests can ONLY use exported (uppercase) methods"
    echo "  ❌ split-002 tests CANNOT access unexported (lowercase) methods"
fi
```

## 🚨🚨🚨 MANDATORY SPLIT PLAN SECTION: Prior Split Analysis

Every split plan MUST include this section (in addition to R420 general requirements):

```markdown
## 🔍 PRIOR SPLIT ANALYSIS (R420 MANDATORY FOR SPLITS)

### Sibling Split Discovery
- **Current Effort**: [effort ID]
- **Previous Splits Reviewed**: [split-001, split-002, etc.]
- **Research Timestamp**: [ISO timestamp]
- **Split Sequence**: This will be split-[XXX]

### API Compatibility Findings (CRITICAL!)
| Package | Source Split | Exports Available | Action Required |
|---------|--------------|-------------------|-----------------|
| pkg/retry | split-001 | Config, NewConfig(), Do() | MUST use Config, NOT DefaultBackoff |
| pkg/auth | split-001 | Authenticator interface | MUST implement this interface |

### Method Visibility Findings (CRITICAL FOR TESTS!)
| Type | Method | Visibility | Source Split | Can Access? | Action |
|------|--------|------------|--------------|-------------|--------|
| MockRegistry | Push() | EXPORTED | split-001 | YES | Safe to use in tests |
| MockRegistry | Pull() | EXPORTED | split-001 | YES | Safe to use in tests |
| MockRegistry | validate() | UNEXPORTED | split-001 | NO | CANNOT use, create workaround |

### Package Structure Findings
| Package Path | Source Split | Purpose | Action Required |
|--------------|--------------|---------|-----------------|
| pkg/retry | split-001 | Retry logic | MUST import, NOT recreate |
| pkg/cmd/push | split-001 | Push commands | MUST extend, NOT replace |

### API Assumptions Verified
- ✅ VERIFIED: retry.Config exists in split-001
- ✅ VERIFIED: retry.NewConfig() function available
- ❌ INCORRECT: retry.DefaultBackoff does NOT exist (use Config instead)

### Method Access Patterns Verified
- ✅ VERIFIED: MockRegistry.Push() is exported (can use in tests)
- ❌ INCORRECT: MockRegistry.validate() is unexported (cannot access from test package)

### Required Integrations
1. MUST import retry package from split-001
2. MUST use retry.Config (NOT retry.DefaultBackoff which doesn't exist)
3. MUST test MockRegistry using ONLY exported methods (Push, Pull)
4. MUST NOT access MockRegistry.validate() (unexported)

### Forbidden Actions
- ❌ DO NOT assume retry.DefaultBackoff exists (use retry.Config)
- ❌ DO NOT create alternative retry package (use split-001's)
- ❌ DO NOT access unexported methods (validate, reset, etc.)
- ❌ DO NOT create duplicate type definitions
```

## 🚨🚨🚨 BLOCKING Validation Gate

Before split plan can be approved:

```bash
# Run R420 validation
bash /path/to/tools/validate-R420-compliance.sh <split-plan-file> || exit 1

# Verify API compatibility section exists
if ! grep -q "API Compatibility Findings" <split-plan-file>; then
    echo "❌ BLOCKING: Split plan missing API compatibility analysis"
    exit 1
fi

# Verify method visibility section exists
if ! grep -q "Method Visibility Findings" <split-plan-file>; then
    echo "❌ BLOCKING: Split plan missing method visibility analysis"
    exit 1
fi

# Verify API assumptions were verified
if ! grep -q "API Assumptions Verified" <split-plan-file>; then
    echo "❌ BLOCKING: Split plan missing API assumption verification"
    exit 1
fi
```

## State-Specific Requirements

### 1. Load Current Effort and Size Context
```bash
# From orchestrator state
CURRENT_EFFORT=$(jq -r '.efforts_in_progress[0]' orchestrator-state-v3.json)
EFFORT_ID=$(echo "$CURRENT_EFFORT" | jq -r '.effort_id')

# Get current size (why split is needed)
CURRENT_LINES=$(echo "$CURRENT_EFFORT" | jq -r '.line_count_tracking.current_line_count')

echo "Splitting effort: $EFFORT_ID"
echo "Current size: $CURRENT_LINES lines (exceeds 900-line enforcement threshold per R535)"
```

### 2. Create Split Inventory (MANDATORY - Systematic Analysis)
```bash
# CRITICAL: Before creating split plans, create systematic inventory
# This ensures comprehensive understanding of what needs to be split

INVENTORY_FILE="$EFFORT_DIR/.software-factory/SPLIT-INVENTORY-$(date +%Y%m%d-%H%M%S).md"
mkdir -p "$(dirname "$INVENTORY_FILE")"

echo "📋 Creating systematic split inventory..."

# Inventory Section 1: File Structure
echo "## File Structure Inventory" >> "$INVENTORY_FILE"
echo "" >> "$INVENTORY_FILE"
echo "**Timestamp**: $(date -u '+%Y-%m-%d %H:%M:%S UTC')" >> "$INVENTORY_FILE"
echo "**Effort**: $EFFORT_ID" >> "$INVENTORY_FILE"
echo "**Total Lines**: $CURRENT_LINES" >> "$INVENTORY_FILE"
echo "" >> "$INVENTORY_FILE"

tree -L 3 "$EFFORT_DIR" >> "$INVENTORY_FILE" 2>/dev/null || find "$EFFORT_DIR" -type f >> "$INVENTORY_FILE"

# Inventory Section 2: Package Analysis
echo "" >> "$INVENTORY_FILE"
echo "## Package Inventory" >> "$INVENTORY_FILE"
echo "" >> "$INVENTORY_FILE"

if [ -d "$EFFORT_DIR/pkg" ]; then
    echo "| Package | Files | Lines | Purpose |" >> "$INVENTORY_FILE"
    echo "|---------|-------|-------|---------|" >> "$INVENTORY_FILE"

    find "$EFFORT_DIR/pkg" -type d -mindepth 1 -maxdepth 1 | while read pkg_dir; do
        pkg_name=$(basename "$pkg_dir")
        file_count=$(find "$pkg_dir" -type f -name "*.go" | wc -l)
        line_count=$(find "$pkg_dir" -type f -name "*.go" -exec wc -l {} \; | awk '{sum+=$1} END {print sum}')

        echo "| $pkg_name | $file_count | $line_count | [Analyze purpose] |" >> "$INVENTORY_FILE"
    done
fi

# Inventory Section 3: Feature Analysis
echo "" >> "$INVENTORY_FILE"
echo "## Feature Inventory" >> "$INVENTORY_FILE"
echo "" >> "$INVENTORY_FILE"

# Identify features from test names
echo "| Feature | Test Count | Test Lines | Status |" >> "$INVENTORY_FILE"
echo "|---------|------------|------------|--------|" >> "$INVENTORY_FILE"

grep -r "^func Test" --include="*_test.go" "$EFFORT_DIR" 2>/dev/null | \
    cut -d: -f2 | sed 's/func Test//; s/(.*//; s/([^)]*)//' | sort -u | \
    while read feature; do
        test_count=$(grep -r "func Test$feature" --include="*_test.go" "$EFFORT_DIR" 2>/dev/null | wc -l)
        echo "| $feature | $test_count | [TBD] | Implemented |" >> "$INVENTORY_FILE"
    done

# Inventory Section 4: Dependency Analysis
echo "" >> "$INVENTORY_FILE"
echo "## Dependency Inventory" >> "$INVENTORY_FILE"
echo "" >> "$INVENTORY_FILE"

echo "### Internal Dependencies" >> "$INVENTORY_FILE"
grep -r "^import" --include="*.go" "$EFFORT_DIR" 2>/dev/null | \
    grep -v "\"C\"" | grep "$(basename "$EFFORT_DIR")" | \
    sort -u >> "$INVENTORY_FILE"

echo "" >> "$INVENTORY_FILE"
echo "### External Dependencies" >> "$INVENTORY_FILE"
grep -r "^import" --include="*.go" "$EFFORT_DIR" 2>/dev/null | \
    grep -v "\"C\"" | grep -v "$(basename "$EFFORT_DIR")" | \
    sort -u | head -20 >> "$INVENTORY_FILE"

# Inventory Section 5: Type and Interface Analysis
echo "" >> "$INVENTORY_FILE"
echo "## Type and Interface Inventory" >> "$INVENTORY_FILE"
echo "" >> "$INVENTORY_FILE"

echo "### Exported Types" >> "$INVENTORY_FILE"
grep -r "^type [A-Z]" --include="*.go" "$EFFORT_DIR" 2>/dev/null | \
    sed 's/:/ - /' >> "$INVENTORY_FILE"

echo "" >> "$INVENTORY_FILE"
echo "### Exported Interfaces" >> "$INVENTORY_FILE"
grep -r "^type [A-Z].*interface" --include="*.go" "$EFFORT_DIR" 2>/dev/null | \
    sed 's/:/ - /' >> "$INVENTORY_FILE"

# Inventory Section 6: Split Boundary Candidates
echo "" >> "$INVENTORY_FILE"
echo "## Potential Split Boundaries" >> "$INVENTORY_FILE"
echo "" >> "$INVENTORY_FILE"

echo "Based on the above inventory, identify logical split points:" >> "$INVENTORY_FILE"
echo "" >> "$INVENTORY_FILE"
echo "1. **Component-based**: Separate by package/module" >> "$INVENTORY_FILE"
echo "2. **Layer-based**: Separate by architectural layer (data, business, API)" >> "$INVENTORY_FILE"
echo "3. **Feature-based**: Separate by independent features" >> "$INVENTORY_FILE"
echo "" >> "$INVENTORY_FILE"

# Inventory Section 7: Size Distribution Analysis
echo "" >> "$INVENTORY_FILE"
echo "## Size Distribution Analysis" >> "$INVENTORY_FILE"
echo "" >> "$INVENTORY_FILE"

echo "| Component | Lines | % of Total | Split Candidate? |" >> "$INVENTORY_FILE"
echo "|-----------|-------|------------|------------------|" >> "$INVENTORY_FILE"

# Calculate percentages for each package
if [ -d "$EFFORT_DIR/pkg" ]; then
    find "$EFFORT_DIR/pkg" -type d -mindepth 1 -maxdepth 1 | while read pkg_dir; do
        pkg_name=$(basename "$pkg_dir")
        line_count=$(find "$pkg_dir" -type f -name "*.go" -exec wc -l {} \; | awk '{sum+=$1} END {print sum}')
        percent=$(echo "scale=1; $line_count * 100 / $CURRENT_LINES" | bc)

        if [ $line_count -gt 500 ]; then
            candidate="YES - Large enough"
        else
            candidate="Consider grouping"
        fi

        echo "| $pkg_name | $line_count | $percent% | $candidate |" >> "$INVENTORY_FILE"
    done
fi

echo "" >> "$INVENTORY_FILE"
echo "---" >> "$INVENTORY_FILE"
echo "**Inventory Complete**: $(date -u '+%Y-%m-%d %H:%M:%S UTC')" >> "$INVENTORY_FILE"

# Commit inventory to git
cd "$EFFORT_DIR" || exit 1
git add "$INVENTORY_FILE"
git commit -m "docs: Create split inventory for $EFFORT_ID

- Total lines: $CURRENT_LINES
- Inventory timestamp: $(date -u '+%Y-%m-%d %H:%M:%S UTC')
- Purpose: Systematic analysis before split planning

[R600]" && git push

echo "✅ Split inventory created: $INVENTORY_FILE"
echo ""
echo "📖 Review inventory before creating split plan!"
```

**Purpose of Split Inventory:**
- Provides systematic catalog of existing implementation
- Identifies natural split boundaries based on actual structure
- Documents dependencies and relationships
- Prevents ad-hoc split decisions
- Creates audit trail for split planning decisions

**BLOCKING:** Cannot create split plan without reviewing inventory first.

### 3. Analyze Current Implementation for Split Boundaries
```bash
# Read split inventory created in Step 2
echo "📖 Reading split inventory: $INVENTORY_FILE"
cat "$INVENTORY_FILE"

# Use inventory to inform split boundary analysis
echo ""
echo "🔍 Analyzing implementation for split boundaries based on inventory..."

# Identify logical split boundaries from inventory data
# (Use package sizes, feature lists, dependency analysis from inventory)
```

### 4. Execute R420 Sibling Split Research (MANDATORY)
See "MANDATORY PRE-FLIGHT: R420 Sibling Split Analysis" above.

**BLOCKING:** Cannot proceed without completing R420 sibling analysis.

**CRITICAL FOR SPLITS:**
- API compatibility is PARAMOUNT
- Method visibility MUST be verified
- Package structures MUST be understood

### 5. Create Split Plan
```markdown
# Split Plan Structure

## Split Overview
- Original Effort: [effort ID]
- Reason for Split: [size/complexity]
- Number of Splits: [N]
- Split Strategy: [sequential/parallel]

## 🔍 PRIOR SPLIT ANALYSIS (R420 MANDATORY FOR SPLITS)
[See mandatory section above]

## Split Definitions

### split-001: [Name]
**Scope:**
- [Feature/package 1]
- [Feature/package 2]

**Estimated Lines:** [estimate]

**Dependencies:**
- External: [list]
- Previous Splits: [none for split-001]

**API Exports:**
[What this split will provide for future splits]

**Forbidden:**
[What this split MUST NOT do]

### split-002: [Name]
**Scope:**
- [Feature/package 3]
- [Feature/package 4]

**Estimated Lines:** [estimate]

**Dependencies:**
- External: [list]
- Previous Splits: split-001
  - MUST import: pkg/retry from split-001
  - MUST use: retry.Config (NOT DefaultBackoff)

**API Compatibility Requirements:**
[Based on R420 analysis of split-001]

**Method Visibility Requirements:**
[Based on R420 analysis of split-001]

**Forbidden:**
- DO NOT recreate retry logic (exists in split-001)
- DO NOT assume retry.DefaultBackoff exists
- DO NOT access unexported methods

## Sequential Execution Order
1. split-001 (no dependencies)
2. split-002 (depends on split-001 APIs)
3. split-003 (depends on split-001 and split-002)

## Testing Strategy
- split-001: Unit tests using split-001 code only
- split-002: Integration tests using split-001 APIs
  - MUST use only exported methods from split-001
  - CANNOT access unexported methods

## Integration Plan
[How splits merge back together]
```

### 6. Validate Split Plan
```bash
# Run R420 validation (CRITICAL!)
bash tools/validate-R420-compliance.sh <split-plan-file> || exit 1

# Validate split plan structure
bash tools/validate-split-plan.sh <split-plan-file> || exit 1

# Check split size estimates
for split in split-001 split-002 split-003; do
    ESTIMATED_LINES=$(grep -A10 "^### $split:" <split-plan-file> | grep "Estimated Lines" | grep -o "[0-9]*")

    if [ $ESTIMATED_LINES -gt 900 ]; then
        echo "❌ BLOCKING: $split estimated at $ESTIMATED_LINES lines (exceeds 900-line enforcement threshold)"
        echo "   VIOLATION of R511: Cannot split a split!"
        echo "   This indicates fundamental design problem - need human architect"
        exit 1
    fi
done
```

### 7. Save Split Plan
```bash
# Save to effort metadata directory
METADATA_DIR="$EFFORT_DIR/.software-factory"
mkdir -p "$METADATA_DIR"

SPLIT_PLAN_FILE="$METADATA_DIR/${EFFORT_ID}-SPLIT-PLAN.md"
cp <split-plan-file> "$SPLIT_PLAN_FILE"

echo "Split plan saved to: $SPLIT_PLAN_FILE"

# Update orchestrator state
jq ".efforts_in_progress[0].split_plan_file = \"$SPLIT_PLAN_FILE\"" orchestrator-state-v3.json > tmp.json && mv tmp.json orchestrator-state-v3.json
jq ".efforts_in_progress[0].status = \"SPLIT_PLANNED\"" orchestrator-state-v3.json > tmp.json && mv tmp.json orchestrator-state-v3.json
```

## Critical Rules for This State

### 🚨🚨🚨 R420 - Cross-Effort Planning Awareness Protocol (BLOCKING)
**Source:** rule-library/R420-cross-effort-planning-awareness-protocol.md

**MANDATORY in CREATE_SPLIT_PLAN:**
- Execute R420 sibling split analysis BEFORE planning
- Read ALL previous split implementations
- Verify API compatibility explicitly
- Check method visibility explicitly
- Include R420 Prior Split Analysis section
- Pass R420 validation

**SPECIAL FOCUS FOR SPLITS:**
- API compatibility is CRITICAL (prevents retry.DefaultBackoff errors)
- Method visibility is CRITICAL (prevents unexported method access)
- Package structure understanding is CRITICAL

**Grading Penalty:** -50% for incomplete research, -100% if leads to API mismatches

---

### 🚨🚨🚨 R511 - Absolute Prohibition on Recursive Splits (BLOCKING)
**Source:** rule-library/R511-absolute-prohibition-recursive-splits.md

**MANDATORY in CREATE_SPLIT_PLAN:**
- If ANY split exceeds 900 lines (R535 enforcement), STOP IMMEDIATELY
- DO NOT attempt to split a split
- Request human architect intervention
- This is a DESIGN problem, not implementation problem

**Grading Penalty:** -100% for attempting recursive split

---

### 🚨🚨 R373 - Mandatory Code Reuse and Interface Compliance (BLOCKING)
**Source:** rule-library/R373-mandatory-code-reuse-and-interface-compliance.md

**MANDATORY in CREATE_SPLIT_PLAN:**
- Later splits MUST reuse APIs from earlier splits
- Later splits MUST NOT recreate packages that exist
- Later splits MUST implement existing interfaces exactly

**Grading Penalty:** -100% for duplicate implementations across splits

---

### 🚨🚨 R310 - Split Scope Strict Adherence Protocol (BLOCKING)
**Source:** rule-library/R310-split-scope-strict-adherence-protocol.md

**MANDATORY in CREATE_SPLIT_PLAN:**
- Split plans MUST be SPECIFIC and DETAILED
- Vague split plans lead to 3-5X overruns
- Each split scope MUST be clearly bounded

**Grading Penalty:** -100% for vague split plans leading to overruns

---

### 🚨 R403 - Split Test Preservation (BLOCKING)
**Source:** rule-library/R403-split-test-preservation.md

**MANDATORY in CREATE_SPLIT_PLAN:**
- Each split MUST include its tests
- Tests MUST accompany their code
- Coverage MUST be maintained across splits
- Test distribution MUST be documented in split plan

**Grading Penalty:** -50% to -75% for separating tests from code

## State Transition Criteria

### Can Transition to COMPLETED When:
1. ✅ R420 sibling split analysis executed completely
2. ✅ R420 Prior Split Analysis section included in plan
3. ✅ API compatibility verified for all previous splits
4. ✅ Method visibility checked for all previous splits
5. ✅ Plan passes R420 validation
6. ✅ All splits estimated under 900 lines (R511 + R535 enforcement)
7. ✅ Split execution order defined
8. ✅ Test distribution documented (R403)
9. ✅ Plan saved to metadata location
10. ✅ Orchestrator state updated with split plan location

### Blocked If:
- ❌ R420 sibling split research not completed
- ❌ API compatibility not verified
- ❌ Method visibility not checked
- ❌ R420 validation fails
- ❌ Any split exceeds 900 lines (R511 + R535 violation)
- ❌ Vague split scopes (R310 violation)

## Grading Criteria for This State

### Sibling Split Research Quality (40% - CRITICAL!)
- ✅ All previous splits reviewed
- ✅ API compatibility explicitly verified
- ✅ Method visibility explicitly checked
- ✅ Package structures understood
- ✅ API assumptions documented and verified

### Split Plan Quality (30%)
- ✅ Clear split boundaries
- ✅ Sequential execution order defined
- ✅ Size estimates under 900 lines each (R535 enforcement threshold)
- ✅ Dependencies explicitly stated
- ✅ Tests distributed appropriately

### R420 Compliance (20%)
- ✅ R420 Prior Split Analysis section complete
- ✅ R420 validation passes
- ✅ API Compatibility Findings table complete
- ✅ Method Visibility Findings table complete
- ✅ API Assumptions Verified section present

### Documentation Quality (10%)
- ✅ Clear instructions for each split
- ✅ Forbidden actions explicitly listed
- ✅ Required integrations specified
- ✅ Test strategy defined

### Failure Conditions:
- **R420 sibling research skipped:** -50%
- **API compatibility not verified:** -50%
- **Method visibility not checked:** -30%
- **Creating API mismatches:** -50%
- **Accessing unexported methods:** -30%
- **Recursive split attempt:** -100% (via R511)
- **Vague split scopes:** -100% (via R310)

## Exit Actions

### Before Transitioning to COMPLETED:
1. Save TODO state (R287)
2. Commit split plan to git
3. Push changes to remote
4. Update orchestrator state with split plan metadata
5. Output continuation flag: `CONTINUE-SOFTWARE-FACTORY=TRUE`

## Common Pitfalls to Avoid

### ❌ WRONG: Assuming APIs Without Verification
```markdown
## split-002 Plan
Dependencies:
- Uses retry.DefaultBackoff from split-001
- Uses MockRegistry.validate() for testing
```

**Result:**
- retry.DefaultBackoff doesn't exist → compilation error
- validate() is unexported → cannot access from test package

### ✅ CORRECT: Verified API Compatibility
```markdown
## 🔍 PRIOR SPLIT ANALYSIS (R420 MANDATORY)

### API Compatibility Findings
| Package | Exports | Source |
|---------|---------|--------|
| pkg/retry | Config, NewConfig(), Do() | split-001 |

### API Assumptions Verified
- ✅ VERIFIED: retry.Config exists
- ❌ INCORRECT: retry.DefaultBackoff does NOT exist

## split-002 Plan
Dependencies:
- MUST use retry.Config from split-001 (NOT DefaultBackoff)
- MUST use MockRegistry.Push() (exported) for testing
- CANNOT use MockRegistry.validate() (unexported)
```

**Result:** APIs match, compilation succeeds, tests work

---

### ❌ WRONG: Not Checking Method Visibility
```markdown
## Test Plan for split-002
Test cases will use:
- MockRegistry.validate() to verify state
- MockRegistry.reset() to clean up
```

**Result:** Both methods unexported → cannot access from external test package

### ✅ CORRECT: Verified Method Visibility
```markdown
## 🔍 PRIOR SPLIT ANALYSIS (R420 MANDATORY)

### Method Visibility Findings
| Method | Visibility | Can Access? |
|--------|------------|-------------|
| Push() | EXPORTED | YES |
| Pull() | EXPORTED | YES |
| validate() | UNEXPORTED | NO |

## Test Plan for split-002
Test cases will use:
- MockRegistry.Push() to test push operations (exported, safe)
- MockRegistry.Pull() to verify retrieval (exported, safe)
- Custom test helpers (since validate() is unexported)
```

**Result:** Only exported methods used, tests compile and run

---

## Summary Checklist

Before completing CREATE_SPLIT_PLAN state:

- [ ] Split inventory created and committed (Step 2 - MANDATORY)
- [ ] Inventory analyzed to identify split boundaries (Step 3)
- [ ] R420 sibling split discovery executed (Step 4)
- [ ] All previous splits reviewed and analyzed
- [ ] API compatibility explicitly verified
- [ ] Method visibility explicitly checked
- [ ] Package structures understood
- [ ] API assumptions documented and verified
- [ ] R420 Prior Split Analysis section in plan
- [ ] API Compatibility Findings table complete
- [ ] Method Visibility Findings table complete
- [ ] API Assumptions Verified section complete
- [ ] R420 validation script passes
- [ ] All splits estimated under 900 lines (R535 enforcement)
- [ ] No recursive split violations (R511)
- [ ] Split scopes clearly defined (R310)
- [ ] Test distribution documented (R403)
- [ ] Sequential execution order defined
- [ ] Plan saved to metadata location
- [ ] Orchestrator state updated
- [ ] Changes committed and pushed

**Only proceed when ALL boxes checked!**

**REMEMBER:** Splits fail from API mismatches and method visibility errors more than any other cause. R420 sibling split analysis is CRITICAL!
