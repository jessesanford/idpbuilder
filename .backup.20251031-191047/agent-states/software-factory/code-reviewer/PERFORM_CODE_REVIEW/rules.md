# Code Reviewer - CODE_REVIEW State Rules

## State Context
You are performing a comprehensive code review of an implementation that has already been measured and confirmed to be within size limits (≤800 lines).

## 🔴🔴🔴 CRITICAL: MEASUREMENT VERIFICATION 🔴🔴🔴

**BEFORE ANY REVIEW, VERIFY SIZE WAS MEASURED CORRECTLY:**
- ✅ Measurement MUST have used `line-counter.sh` tool
- ✅ Tool MUST have auto-detected the base (no -b parameter)
- ✅ Output MUST show "🎯 Detected base:" line
- ❌ If measurement used manual counting (wc -l, git diff, etc.) = STOP AND REMEASURE
- ❌ If shows 11,876 lines for ~500 line effort = WRONG BASE, REMEASURE!

**If previous measurement was incorrect, YOU MUST remeasure using line-counter.sh!**

---

### 🚨🚨🚨 RULE R332 - MANDATORY BUG FILING PROTOCOL (SUPREME LAW) 🚨🚨🚨
**Source:** rule-library/R332-mandatory-bug-filing-protocol.md
**Criticality:** BLOCKING - "Pre-existing" excuse = -100% FAILURE

**ABSOLUTE RULE: "PRE-EXISTING" IS NOT AN ACCEPTABLE EXCUSE**

**ALL BUGS FOUND MUST:**
1. ✅ Be filed immediately in bugs/filed/
2. ✅ Be checked against fix cascades
3. ✅ Create fix cascade if not tracked
4. ✅ Block review if critical and unfixed
5. ✅ Exit 332 if untracked bugs found

**PROHIBITED RESPONSES (R332 Violation = -100%):**
- ❌ "Pre-existing bug"
- ❌ "Acceptable for Phase X"
- ❌ "Out of scope for this effort"
- ❌ "Well-documented future enhancement"

**REQUIRED RESPONSE:**
```markdown
## Bug Found: [Description]

**R332 Actions Taken:**
1. ✅ Filed BUG-YYYYMMDD-HHMMSS-[description].md
2. ✅ Checked fix cascade status
3. ✅ Created fix cascade (if needed)
4. ✅ Assigned to [effort or ERROR_RECOVERY]
5. 🔴 BLOCKING review until addressed

**Bug File**: bugs/filed/BUG-YYYYMMDD-HHMMSS-[description].md
**Review Status**: BLOCKED (exit 332)
```

**PENALTIES:**
- Using "pre-existing": -100% IMMEDIATE FAILURE
- Not filing bug: -100% IMMEDIATE FAILURE
- Passing with known bug: -100% IMMEDIATE FAILURE

---

### 🚨🚨🚨 RULE R320 - No Stub Implementations (CRITICAL BLOCKER) 🚨🚨🚨
**Source:** rule-library/R320-no-stub-implementations.md
**Criticality:** BLOCKING - Any stub = FAILED REVIEW
**Integration:** ALL R320 violations MUST trigger R332 bug filing

**MANDATORY STUB DETECTION PROTOCOL:**
1. Search for ALL "not implemented" patterns
2. Check for TODO in function bodies
3. Verify each function has ACTUAL logic
4. Any stub found = CRITICAL BLOCKER
5. **Stub found = FILE R332 BUG IMMEDIATELY**

**Common stub patterns to detect:**
- `return fmt.Errorf("not implemented")`
- `panic("TODO")` or `panic("unimplemented")`
- `raise NotImplementedError`
- Empty function bodies with just return
- `throw new Error("Not implemented")`

**GRADING PENALTIES:**
- **-50%**: Passing ANY stub implementation
- **-30%**: Classifying stub as "minor issue"
- **-40%**: Marking stub code as "properly implemented"
- **-100%**: Not filing R332 bug for stub

---

### 🚨🚨🚨 RULE R332 - TODO Verification Protocol 🚨🚨🚨
**Source:** rule-library/R332-mandatory-bug-filing-protocol.md (TODO Acceptance Criteria section)
**Criticality:** BLOCKING - Vague TODO acceptance = -100% FAILURE

**MANDATORY TODO VERIFICATION:**

When TODO comment found, Code Reviewer MUST execute verification protocol:

#### Step 1: Extract TODO Context
```bash
TODO_FILE="pkg/feature/handler.go"
TODO_LINE="89"
TODO_TEXT="// TODO: Implement retry logic"

# Extract functionality needed
FUNCTIONALITY="retry logic"
```

#### Step 2: Search Plan Files
```bash
# Search all plan files for this functionality
grep -rn "retry.*logic\|Retry.*Implementation" phase-plans/ wave-plans/ IMPLEMENTATION-PLAN.md

# Example output:
# phase3-wave2-plan.md:156: E3.2.1 - Retry Logic Implementation
# phase3-wave2-plan.md:157: Complete retry mechanism in pkg/feature/handler.go
```

#### Step 3: Verify Plan Details
```bash
# Read the specific plan section
sed -n '156,170p' phase-plans/phase3-wave2-plan.md

# Output must show:
# - Exact file with TODO
# - What will replace TODO
# - When it will be implemented
```

#### Step 4: Compare Phase/Wave
```bash
# Get current phase/wave
CURRENT_PHASE=$(jq -r '.project_progression.current_phase' orchestrator-state-v3.json)
CURRENT_WAVE=$(jq -r '.project_progression.current_wave' orchestrator-state-v3.json)

# Extract planned phase/wave from plan filename
PLANNED_PHASE="3"  # from phase3-wave2-plan.md
PLANNED_WAVE="2"

# Compare
if [ "$PLANNED_PHASE" -gt "$CURRENT_PHASE" ]; then
    echo "✅ Future phase - TODO acceptable with evidence"
elif [ "$PLANNED_PHASE" -eq "$CURRENT_PHASE" ] && [ "$PLANNED_WAVE" -gt "$CURRENT_WAVE" ]; then
    echo "✅ Future wave - TODO acceptable with evidence"
else
    echo "🔴 BUG - Should be implemented now or in past!"
    file_bug "TODO for current/past phase not implemented"
fi
```

#### Step 5: Document Evidence or File Bug

**If plan found and is future:**
```markdown
## TODO Verification: pkg/feature/handler.go:89

**TODO Text**: `// TODO: Implement retry logic`

**Search Command**:
```bash
grep -rn "retry.*logic\|Retry.*Implementation" phase-plans/ wave-plans/
```

**Results**:
- Plan File: `phase-plans/phase3-wave2-plan.md`
- Line Numbers: `156-178`
- Effort ID: `E3.2.1`

**Plan Excerpt**:
```
156: ## E3.2.1 - Retry Logic Implementation
157:
158: Implement retry logic with exponential backoff for network operations.
159:
160: **Files to Modify**:
161: - pkg/feature/handler.go (remove TODO at line 89, implement retry)
```

**Phase/Wave Analysis**:
- Current: Phase 2, Wave 3
- Planned: Phase 3, Wave 2
- Comparison: FUTURE ✅

**Decision**: ✅ TODO ACCEPTED
**Reason**: Explicitly planned for Phase 3 Wave 2 per phase3-wave2-plan.md lines 156-178

**Tracking**: Added to bug-tracking.json for verification in Phase 3
```

**If no plan or plan is past:**
```bash
# File bug per R332
BUG_ID="BUG-$(date +%Y%m%d-%H%M%S)-todo-not-implemented"

cat > bugs/filed/${BUG_ID}.md << 'EOF'
# Bug: TODO Without Valid Plan

**Severity**: HIGH
**File**: pkg/feature/handler.go:89
**TODO**: Implement retry logic

**Issue**: TODO exists but:
- No plan found for implementation, OR
- Plan was for current/past phase

**Search Results**:
```bash
grep -rn "retry.*logic" phase-plans/ wave-plans/
# No relevant results
```

**Current Phase**: Phase 2, Wave 3
**Expected**: Should be implemented or have future plan

**Required Action**: Implement immediately or create explicit plan
EOF

git add bugs/filed/${BUG_ID}.md
git commit -m "bug: File ${BUG_ID} - TODO without valid plan"
git push

# Create fix cascade and block review
exit 332
```

#### Mandatory Fields for TODO Acceptance:

Code Reviewer MUST provide ALL of these:

1. **Plan File Path**: Exact path to plan file
2. **Line Numbers**: Exact lines showing planned implementation
3. **Effort ID**: Which effort will implement this
4. **Phase/Wave**: When it will be implemented
5. **Evidence**: Grep command output showing plan text
6. **Comparison**: Current vs planned phase/wave
7. **Decision**: ACCEPT or BUG with clear reasoning

**If ANY field missing**: MUST file bug per R332 (cannot accept TODO without evidence)

#### Anti-Pattern - PROHIBITED:

```markdown
❌ WRONG:
"TODO is acceptable for future implementation"
"Will be addressed in a later phase"
"Standard pattern for deferred work"

✅ CORRECT:
"TODO verified in phase3-wave2-plan.md lines 156-178
 Planned for E3.2.1 (Phase 3 Wave 2)
 Current: Phase 2 Wave 3
 Evidence: [grep output shown above]
 Decision: ACCEPTED"
```

**GRADING PENALTIES:**
- **-100%**: TODO accepted WITHOUT complete evidence
- **-100%**: TODO accepted with VAGUE reasoning ("will be addressed later")
- **-100%**: Missing plan file path or line numbers
- **-100%**: Missing phase/wave comparison
- **-100%**: Not filing bug when plan doesn't exist or is past
- **✅ 0%**: TODO properly verified with ALL evidence fields

**See R332 TODO Acceptance Criteria section for complete examples**

---

### ℹ️ RULE R108.0.0 - CODE_REVIEW Rules
**Source:** rule-library/RULE-REGISTRY.md#R108
**Criticality:** INFO - Best practice

CODE REVIEW PROTOCOL:
1. **CHECK FOR STUBS FIRST (R320)** - Any stub = FAILED REVIEW
2. Validate implementation against plan
3. Verify test coverage requirements
4. Validate KCP/Kubernetes patterns
5. Check multi-tenancy implementation
6. Assess security and performance
7. Provide detailed feedback

---

## Review Focus Areas

### 1. R332 Bug Filing + R320 Stub Detection - HIGHEST PRIORITY

**EXECUTION PROTOCOL:**
```bash
# Step 1: Scan for all bug categories (R332 + R320 + R355)
cd $EFFORT_DIR

echo "🔍 Scanning for R332 bug categories..."

# Detect TODO comments (R355 violation → R332 bug)
grep -rn "TODO\|FIXME\|XXX\|HACK" --include="*.go" --include="*.py" --include="*.js" --exclude-dir=test > /tmp/todos.txt || true

# Detect stubs (R320 violation → R332 bug)
grep -rn "not.*implemented\|NotImplementedError\|panic.*TODO\|return nil.*TODO" \
    --include="*.go" --include="*.py" --include="*.js" --exclude-dir=test > /tmp/stubs.txt || true

# Detect hardcoded values (R355 violation → R332 bug)
grep -rn "password.*=.*['\"]\\|token.*=.*['\"]\\|localhost:.*[0-9]" \
    --include="*.go" --include="*.py" --include="*.js" --exclude-dir=test > /tmp/hardcoded.txt || true

# Step 2: File R332 bug for EACH issue found
BUGS_FILED=()

# Process TODOs
while IFS= read -r line; do
    [[ -z "$line" ]] && continue

    FILE=$(echo "$line" | cut -d: -f1)
    LINENO=$(echo "$line" | cut -d: -f2)
    CONTENT=$(echo "$line" | cut -d: -f3-)

    BUG_ID="BUG-$(date +%Y%m%d-%H%M%S)-todo-$(basename $FILE | sed 's/\./_/g')-L${LINENO}"

    cat > bugs/filed/${BUG_ID}.md << EOF
# Bug: TODO Comment in Production Code

**Bug ID**: ${BUG_ID}
**Severity**: HIGH
**Component**: ${FILE}
**Line**: ${LINENO}
**Found During**: ${EFFORT_ID}
**Found By**: code-reviewer
**Filed At**: $(date -Iseconds)

## Issue Description
TODO comment found in production code violates R355 Production Readiness.

## Evidence
\`\`\`
${CONTENT}
\`\`\`

## Impact
Incomplete work marker indicates unfinished implementation.

## Required Fix
Complete implementation and remove TODO comment.

## Status
**Current**: FILED
**Assigned To**: ERROR_RECOVERY
**Priority**: HIGH
EOF

    git add bugs/filed/${BUG_ID}.md
    BUGS_FILED+=("${BUG_ID}")
done < /tmp/todos.txt

# Process stubs (similar pattern)
while IFS= read -r line; do
    [[ -z "$line" ]] && continue

    FILE=$(echo "$line" | cut -d: -f1)
    LINENO=$(echo "$line" | cut -d: -f2)
    CONTENT=$(echo "$line" | cut -d: -f3-)

    BUG_ID="BUG-$(date +%Y%m%d-%H%M%S)-stub-$(basename $FILE | sed 's/\./_/g')-L${LINENO}"

    cat > bugs/filed/${BUG_ID}.md << EOF
# Bug: Stub Implementation (R320 Violation)

**Bug ID**: ${BUG_ID}
**Severity**: CRITICAL
**Component**: ${FILE}
**Line**: ${LINENO}
**Found During**: ${EFFORT_ID}
**Found By**: code-reviewer
**Filed At**: $(date -Iseconds)

## Issue Description
Stub implementation found - violates R320 No Stub Implementations.

## Evidence
\`\`\`
${CONTENT}
\`\`\`

## Impact
Non-functional code in production - feature broken.

## Required Fix
Implement actual functionality.

## Status
**Current**: FILED
**Assigned To**: ERROR_RECOVERY
**Priority**: CRITICAL
EOF

    git add bugs/filed/${BUG_ID}.md
    BUGS_FILED+=("${BUG_ID}")
done < /tmp/stubs.txt

# Step 3: Commit all bug files
if [ ${#BUGS_FILED[@]} -gt 0 ]; then
    git commit -m "bug: File ${#BUGS_FILED[@]} bugs found during ${EFFORT_ID} review (R332)

Bugs filed:
$(printf '%s\n' "${BUGS_FILED[@]}")

R332: All bugs MUST be filed, no 'pre-existing' excuse."
    git push
fi

# Step 4: Check/create fix cascades for each bug
cd $CLAUDE_PROJECT_DIR

for BUG_ID in "${BUGS_FILED[@]}"; do
    # Check if bug already tracked
    if ! jq -e ".bugs[] | select(.bug_id == \"${BUG_ID}\")" bug-tracking.json >/dev/null 2>&1; then
        echo "🔴 Bug ${BUG_ID} NOT tracked, adding to bug-tracking.json..."

        # Extract severity from bug file
        SEVERITY=$(grep "^\\*\\*Severity\\*\\*:" bugs/filed/${BUG_ID}.md | cut -d: -f2 | tr -d ' ')
        COMPONENT=$(grep "^\\*\\*Component\\*\\*:" bugs/filed/${BUG_ID}.md | cut -d: -f2 | tr -d ' ')

        # Add bug to bug-tracking.json
        jq ".bugs += [{
            \"bug_id\": \"${BUG_ID}\",
            \"severity\": \"${SEVERITY}\",
            \"component\": \"${COMPONENT}\",
            \"status\": \"FILED\",
            \"assigned_to\": \"ERROR_RECOVERY\",
            \"filed_at\": \"$(date -Iseconds)\",
            \"fix_cascade_id\": null
        }]" bug-tracking.json > tmp && mv tmp bug-tracking.json

        git add bug-tracking.json
        git commit -m "bugs: Add ${BUG_ID} to bug tracking [R288]"
        git push
    fi
done

# Step 5: Block review if critical bugs found
CRITICAL_COUNT=$(grep -l "\\*\\*Severity\\*\\*: CRITICAL" bugs/filed/BUG-*.md 2>/dev/null | wc -l)

if [ "$CRITICAL_COUNT" -gt 0 ]; then
    echo "🔴🔴🔴 REVIEW BLOCKED 🔴🔴🔴"
    echo "Critical bugs found: $CRITICAL_COUNT"
    echo "Bugs must be fixed before review passes"
    echo "Filed bugs: ${#BUGS_FILED[@]}"
    echo ""
    echo "R332: 'Pre-existing' is NOT an excuse"
    echo "ALL bugs MUST be fixed or tracked in fix cascades"
    echo ""

    # Create BLOCKED-BUGS report
    cat > BLOCKED-BUGS.md << EOF
# Review Blocked - Critical Bugs Must Be Fixed

**Review**: ${EFFORT_ID}
**Blocked At**: $(date -Iseconds)
**Critical Bugs**: $CRITICAL_COUNT
**Total Bugs**: ${#BUGS_FILED[@]}

## Critical Bugs

$(for bug in "${BUGS_FILED[@]}"; do
    if grep -q "\\*\\*Severity\\*\\*: CRITICAL" bugs/filed/${bug}.md; then
        echo "- ${bug}"
        grep "^## Issue Description" -A 2 bugs/filed/${bug}.md
    fi
done)

## All Bugs Filed

$(printf '- %s\n' "${BUGS_FILED[@]}")

## R332 Requirement

ALL bugs MUST be:
1. Filed (✅ DONE)
2. In fix cascades (✅ DONE)
3. Fixed or in active fix (🔴 REQUIRED)

**Review cannot pass until critical bugs fixed.**
EOF

    git add BLOCKED-BUGS.md
    git commit -m "review: Block ${EFFORT_ID} - critical bugs must be fixed"
    git push

    exit 332  # R332 violation - untracked/unfixed bugs
fi

echo "✅ R332 Compliance: ${#BUGS_FILED[@]} bugs filed and tracked"
```

### 2. Test Coverage Validation

---
### 🚨🚨 RULE R032.0.0 - Test Coverage Requirements
**Source:** rule-library/RULE-REGISTRY.md#R032
**Criticality:** MANDATORY - Required for approval

MANDATORY COVERAGE VALIDATION:
- Unit Tests: 90% line coverage minimum
- Integration Tests: All API endpoints covered
- Multi-tenant Tests: Cross-workspace scenarios tested
- Error Cases: All error paths validated
- Performance: Resource usage within limits
---

```python
def validate_test_coverage(effort_dir):
    """Validate test coverage meets requirements"""
    
    coverage_results = {
        'unit_test_coverage': measure_unit_test_coverage(effort_dir),
        'integration_test_coverage': assess_integration_tests(effort_dir),
        'multi_tenant_test_coverage': assess_multi_tenant_tests(effort_dir),
        'error_case_coverage': assess_error_case_coverage(effort_dir),
        'performance_test_coverage': assess_performance_tests(effort_dir)
    }
    
    # Check for critical coverage gaps
    unit_coverage = coverage_results['unit_test_coverage'].get('percentage', 0)
    if unit_coverage < 90:
        return {
            'critical_issue': 'INSUFFICIENT_UNIT_COVERAGE',
            'blocking': True
        }
    
    return coverage_results
```

### 3. KCP/Kubernetes Pattern Validation

---
### ℹ️ RULE R037.0.0 - Pattern Compliance
**Source:** rule-library/RULE-REGISTRY.md#R037
**Criticality:** INFO - Best practice

KCP PATTERN CHECKLIST:
✅ Multi-tenancy: Logical cluster awareness
✅ APIExport: Proper integration patterns
✅ Virtual Workspace: Compliance with VW model
✅ Syncer: Compatible with syncer patterns
✅ RBAC: Workspace-scoped permissions
✅ Resource Quotas: Tenant isolation enforcement
---

### 4. Security Review

---
### ℹ️ RULE R038.0.0 - Security Review
**Source:** rule-library/RULE-REGISTRY.md#R038
**Criticality:** INFO - Best practice

SECURITY CHECKLIST:
✅ Input validation on all external data
✅ Workspace isolation properly enforced
✅ RBAC permissions correctly implemented
✅ No hardcoded credentials or secrets
✅ Error messages don't leak sensitive information
✅ Resource access properly authorized
---

### 5. Architecture Compliance

Review implementation against architectural plan:
- Component structure matches design
- Interfaces properly implemented
- Design patterns correctly applied
- Dependencies appropriately managed

## Review Decision Framework

```python
def make_review_decision(review_data):
    """Make final review decision based on all validation results"""
    
    # Critical blocking issues
    blocking_issues = []
    
    # STUB IMPLEMENTATIONS (HIGHEST PRIORITY - R320)
    stub_result = review_data.get('stub_detection', {})
    if stub_result.get('stubs_found', False):
        blocking_issues.append({
            'type': 'STUB_IMPLEMENTATION_DETECTED',
            'severity': 'CRITICAL_BLOCKER',
            'description': f"Found {stub_result.get('stub_count', 0)} stub implementations",
            'action_required': 'COMPLETE_IMPLEMENTATION'
        })
    
    # Test coverage (CRITICAL)
    coverage_result = review_data.get('test_coverage', {})
    if not coverage_result.get('meets_requirements', False):
        blocking_issues.append({
            'type': 'INSUFFICIENT_COVERAGE',
            'description': f"Coverage {coverage_result.get('coverage_score', 0)}% < 90%",
            'action_required': 'IMPROVE_TESTS'
        })
    
    # Make decision
    if blocking_issues:
        return {
            'result': 'FAILED',
            'blocking_issues': blocking_issues,
            'can_proceed': False
        }
    else:
        warnings = collect_review_warnings(review_data)
        
        if len(warnings) == 0:
            return {'result': 'PASSED'}
        elif len(warnings) <= 3:
            return {'result': 'PASSED_WITH_WARNINGS'}
        else:
            return {'result': 'CHANGES_RECOMMENDED'}
```

## Review Documentation (R383 COMPLIANT)

### MANDATORY: Use sf_metadata_path helper from R383

```bash
# Include R383 helper function (MANDATORY)
sf_metadata_path() {
    local phase="$1"
    local wave="$2"
    local effort="$3"
    local filename="$4"
    local ext="$5"

    if [[ -z "$phase" || -z "$wave" || -z "$effort" || -z "$filename" || -z "$ext" ]]; then
        echo "❌ R383 VIOLATION: Missing parameters to sf_metadata_path" >&2
        exit 1
    fi

    local dir=".software-factory/phase${phase}/wave${wave}/${effort}"
    mkdir -p "$dir"

    local timestamp=$(date +%Y%m%d-%H%M%S)
    local full_path="${dir}/${filename}--${timestamp}.${ext}"

    echo "$full_path"
}

# Determine context
PHASE="${PHASE:-1}"  # Get from environment or context
WAVE="${WAVE:-1}"    # Get from environment or context
EFFORT_NAME="${EFFORT_NAME:-$(basename $(pwd))}"  # Get from context

# Create review report using R383-compliant path
REPORT_FILE=$(sf_metadata_path "$PHASE" "$WAVE" "$EFFORT_NAME" "CODE-REVIEW-REPORT" "md")
echo "📝 R383 COMPLIANT: Creating review report at: $REPORT_FILE"
```

### Create CODE-REVIEW-REPORT with timestamp:
```bash
cat > "$REPORT_FILE" << 'EOF'
# Code Review Report
effort_id: "${EFFORT_NAME}"
reviewed_at: "$(date -Iseconds)"
reviewer: "code-reviewer-agent"
report_path: "${REPORT_FILE}"

## Pre-Review Verification
size_measurement:
  completed: true
  compliant: true
  lines: [number]
  
## Review Results

### 1. Stub Detection (R320)
stubs_found: [true/false]
stub_count: [number]
stub_locations: []
result: [PASSED/FAILED]

### 2. Test Coverage
unit_test_coverage: [percentage]
integration_tests: [count]
multi_tenant_scenarios: [count]
meets_requirements: [true/false]

### 3. KCP Compliance
multi_tenancy_score: [percentage]
api_export_integration: [percentage]
workspace_isolation: [percentage]
overall_compliance: [percentage]

### 4. Security Review
input_validation: [PASS/FAIL]
workspace_isolation: [PASS/FAIL]
secret_management: [PASS/FAIL]
critical_issues: [count]

### 5. Architecture Review
plan_adherence: [percentage]
design_patterns: [percentage]
interface_compliance: [percentage]

## Final Decision
REVIEW_STATUS: [PASSED/FAILED/PASSED_WITH_WARNINGS]
blocking_issues: []
warnings: []
recommendations: []

## Required Actions
[List any required fixes if FAILED]
EOF

# Commit the review report
git add "$REPORT_FILE"
git commit -m "review: add code review report for ${EFFORT_NAME}

Report location: $REPORT_FILE
Phase $PHASE, Wave $WAVE
R383 compliant: timestamp included"
git push

echo "✅ R383 COMPLIANT: Review report created with timestamp"

# R340: Report review report location to Orchestrator
echo ""
echo "📋 Code Review Report: $REPORT_FILE"
echo "Effort ID: $EFFORT_NAME"
echo "Phase: $PHASE"
echo "Wave: $WAVE"
echo "R340: Review report tracked for lookup"
echo ""
```

## State Transitions

From CODE_REVIEW state:
- **REVIEW_PASSED** → COMPLETED (Implementation approved)
- **REVIEW_FAILED** → CREATE_FIX_PLAN (Issues need fixing)
- **CRITICAL_STUBS_FOUND** → CREATE_FIX_PLAN (R320 violation)

## Success Criteria
- ✅ Thoroughly checked for stub implementations
- ✅ Validated test coverage ≥90%
- ✅ Verified KCP patterns compliance
- ✅ Completed security assessment
- ✅ Created comprehensive review report
- ✅ Made clear pass/fail decision

## Failure Triggers
- ❌ Missing stub detection = R320 VIOLATION (-50%)
- ❌ Passing stub code = R320 VIOLATION (-50%)
- ❌ Incomplete review = -30% penalty
- ❌ Missing review report = -40% penalty

## 🔴🔴🔴 R405 - MANDATORY AUTOMATION CONTINUATION FLAG 🔴🔴🔴

**SUPREME LAW - PENALTY FOR VIOLATION: -100% GRADE**

### YOU MUST OUTPUT THE CONTINUATION FLAG AS YOUR LAST ACTION

**EVERY STATE COMPLETION MUST END WITH EXACTLY ONE OF:**
```bash
# If state completed successfully and factory should continue:
echo "CONTINUE-SOFTWARE-FACTORY=TRUE"

# If error/block/manual review needed:
echo "CONTINUE-SOFTWARE-FACTORY=FALSE"
```

### CRITICAL REQUIREMENTS:
1. **ABSOLUTE LAST OUTPUT**: This MUST be the very last line of output before state completion
2. **NO TEXT AFTER**: No explanations, summaries, or any text after the flag
3. **EXACTLY AS SHOWN**: Use exact format - no variations like CONTINUE-ORCHESTRATING
4. **ALWAYS REQUIRED**: Every single state must output this flag
5. **GREPPABLE**: Must be on its own line for automation parsing

### WHEN TO USE TRUE:
- ✅ State work completed successfully
- ✅ All validations passed
- ✅ Ready for next state
- ✅ No blockers detected
- ✅ All requirements met

### WHEN TO USE FALSE:
- ❌ Any unrecoverable error occurred
- ❌ Manual intervention required
- ❌ Missing required files or configs
- ❌ Test failures blocking progress
- ❌ Ambiguous or unclear instructions
- ❌ Wrong working directory or branch
- ❌ State machine validation failed

### IMPLEMENTATION PATTERN:
```bash
# At the VERY END of state execution, after ALL work:

# Determine success/failure
if [[ "$STATE_COMPLETED_PROJECT_DONEFULLY" == "true" ]]; then
    echo "✅ State work completed successfully"
    echo "CONTINUE-SOFTWARE-FACTORY=TRUE"
else
    echo "❌ State encountered issues requiring intervention"
    echo "CONTINUE-SOFTWARE-FACTORY=FALSE"
fi

# DO NOT OUTPUT ANYTHING AFTER THE FLAG!
```

### GRADING IMPACT:
- **Missing flag**: -100% AUTOMATIC FAILURE
- **Text after flag**: -50% penalty
- **Wrong format**: -100% AUTOMATIC FAILURE
- **Multiple flags**: -50% penalty

**See: $CLAUDE_PROJECT_DIR/rule-library/R405-automation-continuation-flag.md**

