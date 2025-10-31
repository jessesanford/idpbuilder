# 🔴🔴🔴 RULE R332 - MANDATORY BUG FILING PROTOCOL (SUPREME LAW) 🔴🔴🔴

**Classification**: SUPREME LAW
**Severity**: BLOCKING
**Penalty**: -100% (IMMEDIATE FAILURE)
**Applies To**: Code Reviewer, Integration Agent, All Review States
**Replaces**: "Pre-existing bug" excuses across all reviews

---

## 🚨🚨🚨 PURPOSE 🚨🚨🚨

Eliminate the "pre-existing bug" excuse that creates systemic failure where bugs are NEVER FIXED.

**THE CRISIS THIS SOLVES:**
```
Phase 1 Review (E1.1.2):
  // TODO: Implement actual push logic in Phase 2
  Review: ✅ TODO comment is properly scoped for Phase 2 (acceptable)

Phase 2 Review (E2.2.2):
  ./pkg/cmd/push/root.go: TODO: Implement actual push logic (pre-existing)
  Review: ✅ PASS - Pre-existing from earlier efforts

Result:
  - TODO never fixed
  - Feature never implemented
  - 2 integration checkpoints PASSED
  - Production code has stub functionality
```

**THE PATTERN OF FAILURE:**
1. Phase 1: "Acceptable for Phase 2" (defer)
2. Phase 2: "Pre-existing" (defer again)
3. Result: Bug NEVER gets fixed

---

## 🚨🚨🚨 THE ABSOLUTE RULE 🚨🚨🚨

### "PRE-EXISTING" IS NOT AN ACCEPTABLE REASON TO PASS A REVIEW

**If a bug exists, it MUST be fixed. The age of the bug is IRRELEVANT.**

Every bug found MUST:
1. ✅ Be filed in bug tracking system
2. ✅ Be tracked through fix cascade
3. ✅ Block review until addressed
4. ✅ Have verification that fix is in progress

**NO EXCEPTIONS. NO DEFERRALS. NO "PRE-EXISTING" PASSES.**

---

## 📋 WHEN CODE REVIEWER FINDS ANY ISSUE

**Regardless of when the bug was introduced:**

### Step 1: File Bug Immediately

```bash
# Use bug tracking system (create if doesn't exist)
cd $CLAUDE_PROJECT_DIR

# Create bug report
BUG_ID="BUG-$(date +%Y%m%d-%H%M%S)-${ISSUE_SUMMARY}"

cat > bugs/filed/${BUG_ID}.md << 'EOF'
# Bug Report

**Bug ID**: BUG-YYYYMMDD-HHMMSS-short-description
**Severity**: [CRITICAL/HIGH/MEDIUM/LOW]
**Component**: [path/to/affected/code.go]
**Found During**: [effort ID or integration level]
**Found By**: code-reviewer
**Filed At**: [ISO-8601 timestamp]

## Issue Description
[Clear description of the bug]

## Evidence
```
[Code snippet showing bug]
```

## Impact
[What breaks because of this bug]

## Root Cause
[Why this bug exists]

## Required Fix
[What needs to be done to fix it]

## Affected Files
- [List all affected files]

## Verification Steps
[How to verify the fix works]

## Related Bugs
[List any related or duplicate bugs]

## Status
**Current**: FILED
**Assigned To**: [Next effort or ERROR_RECOVERY]
**Fix Cascade**: [FC-ID if created]
**Priority**: [Based on severity]

## Timeline
- Filed: [timestamp]
EOF

git add bugs/filed/${BUG_ID}.md
git commit -m "bug: File ${BUG_ID} - ${ISSUE_SUMMARY}"
git push
```

### Step 2: Check Fix Cascade Status

```bash
# Check if bug is already being fixed
grep -r "${BUG_ID}" fix-cascade/ orchestrator-state-v3.json 2>/dev/null

# If in active fix cascade:
if jq -e ".fix_cascades[] | select(.bugs[] | contains(\"${BUG_ID}\"))" orchestrator-state-v3.json >/dev/null; then
    #   - Verify fix is in progress
    #   - Verify fix scope includes this issue
    #   - Allow review to continue if fix verified
    echo "✅ Bug ${BUG_ID} is being addressed in active fix cascade"

    # Get fix details
    FIX_EFFORT=$(jq -r ".fix_cascades[] | select(.bugs[] | contains(\"${BUG_ID}\")) | .assigned_to" orchestrator-state-v3.json)
    FIX_STATUS=$(jq -r ".fix_cascades[] | select(.bugs[] | contains(\"${BUG_ID}\")) | .status" orchestrator-state-v3.json)

    echo "Effort: ${FIX_EFFORT}"
    echo "Status: ${FIX_STATUS}"
else
    # If NOT in fix cascade:
    #   - Create fix cascade entry
    #   - Assign to appropriate effort
    #   - Block review until cascade planned
    echo "🔴 Bug ${BUG_ID} is NOT being addressed!"
    echo "Creating fix cascade..."
fi
```

### Step 3: Verify Fix or Block Review

```bash
if jq -e ".fix_cascades[] | select(.bugs[] | contains(\"${BUG_ID}\"))" orchestrator-state-v3.json >/dev/null; then
    echo "✅ Bug ${BUG_ID} is being addressed in active fix cascade"

    # Verify cascade details
    CASCADE_ID=$(jq -r ".fix_cascades[] | select(.bugs[] | contains(\"${BUG_ID}\")) | .cascade_id" orchestrator-state-v3.json)
    FIX_EFFORT=$(jq -r ".fix_cascades[] | select(.bugs[] | contains(\"${BUG_ID}\")) | .assigned_to" orchestrator-state-v3.json)
    FIX_STATUS=$(jq -r ".fix_cascades[] | select(.bugs[] | contains(\"${BUG_ID}\")) | .status" orchestrator-state-v3.json)

    echo "Cascade: ${CASCADE_ID}"
    echo "Effort: ${FIX_EFFORT}"
    echo "Status: ${FIX_STATUS}"

    # Allow review to continue
else
    echo "🔴🔴🔴 Bug ${BUG_ID} is NOT being addressed!"
    echo "BLOCKING review until fix cascade created"

    # Create fix cascade
    CASCADE_ID="FC-$(date +%Y%m%d-%H%M%S)"
    SEVERITY=$(extract_bug_severity "bugs/filed/${BUG_ID}.md")

    jq ".fix_cascades += [{
        \"cascade_id\": \"${CASCADE_ID}\",
        \"bugs\": [\"${BUG_ID}\"],
        \"severity\": \"${SEVERITY}\",
        \"component\": \"$(extract_component bugs/filed/${BUG_ID}.md)\",
        \"status\": \"FILED\",
        \"assigned_to\": \"ERROR_RECOVERY\",
        \"created_at\": \"$(date -Iseconds)\"
    }]" orchestrator-state-v3.json > tmp && mv tmp orchestrator-state-v3.json

    # Update state file
    git add orchestrator-state-v3.json
    git commit -m "state: Add fix cascade ${CASCADE_ID} for ${BUG_ID}"
    git push

    exit 332  # BLOCK REVIEW - R332 violation
fi
```

---

## 🐛 BUG CATEGORIES (ALL Require Filing)

### 1. TODO Comments - VERIFICATION REQUIRED

```go
// TODO: Implement actual push logic
// FIXME: This needs proper error handling
// XXX: Temporary workaround
// HACK: Quick fix for demo
```

**Action**: **MUST verify against plan files first (see TODO Acceptance Criteria below)**, OR file bug immediately if no valid plan

#### TODO Acceptance Criteria (R332 Exception)

**TODOs are ACCEPTABLE IF AND ONLY IF Code Reviewer can provide ALL evidence:**

1. **Future Implementation Explicitly Planned**
   - Functionality MUST be in a documented plan file
   - Plan file MUST specify wave/phase for implementation
   - Plan MUST be for FUTURE wave/phase (not current or past)

2. **Code Reviewer MUST Provide Evidence**
   ```markdown
   ## TODO Accepted: Implementation Planned for Future Phase

   **File**: pkg/feature/handler.go:45
   **TODO**: // TODO: Implement retry logic (Phase 3)

   **Evidence**:
   - Plan File: `phase-plans/phase3-plan.md`
   - Line Numbers: 156-178
   - Planned Effort: E3.2.1 - Retry Mechanism Implementation
   - Scheduled Wave: Phase 3, Wave 2
   - Status: ✅ VERIFIED - Implementation planned for future

   **Verification**:
   ```bash
   grep -n "retry.*mechanism\|Retry.*Implementation" phase-plans/phase3-plan.md
   # Output:
   # 156: ## E3.2.1 - Retry Mechanism Implementation
   # 157: Implement retry logic with exponential backoff
   # 158: - Location: pkg/feature/handler.go
   # 159: - TODO removal: Line 45
   ```

   ✅ TODO ACCEPTED - Will be addressed in Phase 3 Wave 2 per plan
   ```

3. **Plan Must Be Specific**
   - Plan MUST identify exact file with TODO
   - Plan MUST specify what will replace TODO
   - Plan MUST be in upcoming wave/phase (not vague "future")

#### When TODO is a BUG (Must File Immediately):

1. **No Plan Exists**
   ```bash
   grep -r "retry.*logic\|Retry.*Implementation" phase-plans/ wave-plans/
   # No results found

   🔴 BUG: No plan for TODO implementation
   ```

2. **Plan Was for Current/Past Phase**
   ```markdown
   **File**: pkg/push/operations.go:69
   **TODO**: // TODO: Implement actual push logic in Phase 2

   **Current Phase**: Phase 2 ← BUG! Should be implemented NOW!

   🔴 BUG: TODO deferred to current phase but not implemented
   ```

3. **Plan Doesn't Match TODO**
   ```bash
   # TODO says: "Implement push logic"
   # Plan says: "Add authentication"

   🔴 BUG: Plan doesn't address this TODO
   ```

4. **Reviewer Cannot Provide Evidence**
   ```markdown
   Code Reviewer Response:
   "TODO is acceptable for future implementation"

   Evidence Provided: NONE ❌

   🔴 FAIL: Must provide exact plan file + line numbers (R332)
   ```

#### Verification Protocol for Code Reviewer:

When encountering TODO, Code Reviewer MUST:

```bash
#!/bin/bash
# R332 TODO Verification Protocol

TODO_FILE="$1"
TODO_LINE="$2"
TODO_TEXT="$3"

# Extract what functionality is needed
FUNCTIONALITY=$(echo "$TODO_TEXT" | sed 's/.*TODO: //')

# Search plan files
PLAN_FILES=$(find phase-plans wave-plans -name "*.md" -o -name "*.txt" 2>/dev/null)

echo "=== Searching for planned implementation of: $FUNCTIONALITY ==="

for plan in $PLAN_FILES; do
    # Search for related functionality
    MATCHES=$(grep -in "$FUNCTIONALITY" "$plan" 2>/dev/null)

    if [ -n "$MATCHES" ]; then
        echo "✅ Found in: $plan"
        echo "$MATCHES"

        # Extract phase/wave from plan filename
        PHASE=$(echo "$plan" | grep -o "phase[0-9]*" | head -1)
        WAVE=$(echo "$plan" | grep -o "wave[0-9]*" | head -1)

        # Get current phase/wave from orchestrator-state-v3.json
        CURRENT_PHASE=$(jq -r '.current_phase' orchestrator-state-v3.json 2>/dev/null || echo "0")
        CURRENT_WAVE=$(jq -r '.current_wave' orchestrator-state-v3.json 2>/dev/null || echo "0")

        # Compare
        if is_future_phase_wave "$PHASE" "$WAVE" "$CURRENT_PHASE" "$CURRENT_WAVE"; then
            echo "✅ ACCEPT TODO - Planned for future: $PHASE $WAVE"
            echo "   Plan: $plan (lines shown above)"
            exit 0  # TODO acceptable
        else
            echo "🔴 BUG - Planned for current/past: $PHASE $WAVE"
            echo "   Should have been implemented already!"
            exit 1  # BUG - file it
        fi
    fi
done

echo "🔴 BUG - No plan found for: $FUNCTIONALITY"
echo "   No future implementation documented"
exit 1  # BUG - file it

is_future_phase_wave() {
    local planned_phase="${1#phase}"
    local planned_wave="${2#wave}"
    local current_phase="$3"
    local current_wave="$4"

    # Handle empty values
    [ -z "$planned_phase" ] && return 1
    [ -z "$planned_wave" ] && return 1

    if [ "$planned_phase" -gt "$current_phase" ]; then
        return 0  # Future phase
    elif [ "$planned_phase" -eq "$current_phase" ] && [ "$planned_wave" -gt "$current_wave" ]; then
        return 0  # Future wave in current phase
    else
        return 1  # Current or past
    fi
}
```

#### Evidence Requirements:

Code Reviewer MUST document:

```markdown
## TODO Analysis: [File:Line]

**TODO Text**: [exact TODO comment]

**Search Results**:
```bash
grep -rn "functionality keywords" phase-plans/ wave-plans/
```

**Plan File Found**: phase-plans/phase3-wave2-plan.md
**Line Numbers**: 156-178
**Effort ID**: E3.2.1
**Effort Name**: Retry Mechanism Implementation

**Plan Excerpt**:
```
156: ## E3.2.1 - Retry Mechanism Implementation
157:
158: Implement retry logic with exponential backoff for network operations.
159:
160: **Files to Modify**:
161: - pkg/feature/handler.go (remove TODO at line 45, implement retry)
162:
163: **Implementation Details**:
164: - Exponential backoff: 1s, 2s, 4s, 8s, 16s
165: - Max retries: 5
166: - Configurable via flags
```

**Phase/Wave Comparison**:
- Current: Phase 2, Wave 3
- Planned: Phase 3, Wave 2
- Status: ✅ FUTURE (acceptable)

**Decision**: ✅ TODO ACCEPTED
**Reason**: Explicitly planned for Phase 3 Wave 2 per documented plan

**Tracking**: Added note to orchestrator-state-v3.json for verification in Phase 3
```

#### Mandatory Fields for TODO Acceptance:

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
 Evidence: [grep output]
 Decision: ACCEPTED"
```

#### Penalty Updates for TODOs:

| Scenario | Penalty |
|----------|---------|
| TODO with NO plan evidence | File BUG (-100% if not filed) |
| TODO planned for past phase | File BUG (-100% if not filed) |
| TODO accepted WITHOUT evidence | -100% IMMEDIATE FAILURE |
| TODO accepted with VAGUE reasoning | -100% IMMEDIATE FAILURE |
| Reviewer says "will be addressed later" | -100% IMMEDIATE FAILURE (R332 violation) |
| TODO properly verified with evidence | ✅ ACCEPTABLE (no penalty) |

### 2. Stub Code

```go
func DoThing() error {
    fmt.Println("Success")
    return nil  // Returns without doing work!
}

func PushImages(cmd *cobra.Command, args []string) error {
    // TODO: Implement actual push logic
    return nil
}
```

**Action**: File bug, create fix cascade, implement immediately

### 3. Missing Integration

```go
// pkg/cmd/push/root.go
func RunE(cmd *cobra.Command, args []string) error {
    return nil  // Should call push.PushImages(cmd, args)
}

// pkg/push/operations.go
func PushImages(...) { ... }  // Exists but never called!
```

**Action**: File bug, create fix cascade, wire immediately

### 4. Non-Existent Flags

```bash
# Demo uses:
./command --registry localhost:5000

# But --help shows:
# No --registry flag exists!
```

**Action**: File bug, fix demo or implement flag

### 5. Incomplete Features

```go
// Feature partially implemented
// Some functions exist, others missing
// Integration incomplete
```

**Action**: File bug, complete implementation

---

## ❌ PROHIBITED REVIEWER RESPONSES

### 🚨🚨🚨 FORBIDDEN - Will Result in -100% FAILURE 🚨🚨🚨

1. **"Pre-existing bug"**
   - ❌ NO EXCUSE - file bug and track fix
   - **Penalty**: -100% IMMEDIATE FAILURE

2. **"Acceptable for Phase X"**
   - ❌ NO DEFERRALS - fix now or create fix cascade
   - **Penalty**: -100% IMMEDIATE FAILURE

3. **"Well-documented future enhancement"**
   - ❌ ENHANCEMENT ≠ BUG - if broken, file bug
   - **Penalty**: -100% IMMEDIATE FAILURE

4. **"Standard pattern"**
   - ❌ BROKEN PATTERN - file bug
   - **Penalty**: -100% IMMEDIATE FAILURE

5. **"Out of scope for this effort"**
   - ❌ NOT AN EXCUSE - file bug for fix cascade
   - **Penalty**: -100% IMMEDIATE FAILURE

### ✅ REQUIRED - Proper Response

```markdown
## Issue Found: TODO in Production Code

**File**: pkg/cmd/push/root.go:69
**Issue**: Stub code with TODO instead of implementation
**Severity**: CRITICAL

**Actions Taken**:
1. ✅ Filed BUG-20251006-143000-push-stub-not-implemented.md
2. ✅ Verified bug NOT in active fix cascade
3. ✅ Created fix cascade entry in orchestrator-state-v3.json
4. ✅ Assigned to ERROR_RECOVERY effort
5. 🔴 BLOCKING review until fix planned

**Verification**:
- Bug file: bugs/filed/BUG-20251006-143000-push-stub-not-implemented.md
- Fix cascade: orchestrator-state-v3.json (cascade_id: FC-20251006-143000)
- Status: Pending ERROR_RECOVERY

**Review Status**: BLOCKED (exit 332)
```

---

## 🔄 FIX CASCADE INTEGRATE_WAVE_EFFORTS

When bug filed, Code Reviewer MUST:

### 1. Check Orchestrator State

```bash
# List all fix cascades
jq '.fix_cascades[]' orchestrator-state-v3.json

# Check if specific bug in cascade
jq -e ".fix_cascades[] | select(.bugs[] | contains(\"${BUG_ID}\"))" orchestrator-state-v3.json
```

### 2. Create Fix Cascade Entry (if not exists)

```json
{
  "fix_cascades": [
    {
      "cascade_id": "FC-20251006-143000",
      "bugs": ["BUG-20251006-143000-push-stub-not-implemented"],
      "severity": "CRITICAL",
      "component": "pkg/cmd/push",
      "status": "FILED",
      "assigned_to": "ERROR_RECOVERY",
      "created_at": "2025-10-06T14:30:00Z"
    }
  ]
}
```

### 3. Update State File

```bash
git add orchestrator-state-v3.json
git commit -m "state: Add fix cascade for ${BUG_ID}"
git push
```

### 4. Transition to ERROR_RECOVERY (if blocking)

```bash
# Update current state
jq '.state_machine.current_state = "ERROR_RECOVERY"' orchestrator-state-v3.json > tmp && mv tmp orchestrator-state-v3.json

# Add error context
jq ".error_context = {
    \"error_type\": \"UNRESOLVED_BUG\",
    \"bug_id\": \"${BUG_ID}\",
    \"fix_cascade_id\": \"${CASCADE_ID}\",
    \"blocking_review\": true,
    \"severity\": \"CRITICAL\"
}" orchestrator-state-v3.json > tmp && mv tmp orchestrator-state-v3.json

git add orchestrator-state-v3.json
git commit -m "state: Transition to ERROR_RECOVERY for ${BUG_ID}"
git push
```

---

## 🛡️ ENFORCEMENT

### Code Reviewer MUST:

- ✅ File bug for EVERY issue found (no exceptions)
- ✅ Check fix cascade status
- ✅ Block review if bug not being addressed
- ✅ Create BLOCKED-BUGS.md listing all blocking bugs
- ✅ FAIL review with exit code 332 if bugs unfixed

### Orchestrator MUST:

- ✅ Verify all bugs filed have fix cascades
- ✅ Transition to ERROR_RECOVERY if critical bugs unfixed
- ✅ Track fix cascade progress
- ✅ Prevent integration completion with open bugs

### Integration Agent MUST:

- ✅ Check for bugs/ directory entries
- ✅ Verify all bugs addressed before reporting success
- ✅ Include bug status in integration report

---

## ⚖️ PENALTIES

| Violation | Penalty |
|-----------|---------|
| Using "pre-existing" excuse | -100% IMMEDIATE FAILURE |
| Not filing bug when found | -100% IMMEDIATE FAILURE |
| Passing review with known bug | -100% IMMEDIATE FAILURE |
| Not creating fix cascade | -75% |
| Not checking cascade status | -50% |
| Deferring bug to future phase | -100% IMMEDIATE FAILURE |

---

## ✅ PROJECT_DONE CRITERIA

Review can ONLY pass if:

- ✅ No bugs found, OR
- ✅ All bugs filed, AND
- ✅ All bugs in active fix cascade, AND
- ✅ All critical bugs fixed in current effort, AND
- ✅ All non-critical bugs tracked for next effort

---

## 🔗 RELATED RULES

- **R355**: TODO Comment Management (now requires R332 filing)
- **R331**: Demo Validation Protocol (catches bugs demos should find)
- **R259**: Phase Integration After ERROR_RECOVERY (fix cascade execution)
- **R206**: State Machine Validation (ERROR_RECOVERY transitions)

---

## 📝 EXAMPLE WORKFLOW

### Scenario: Reviewer finds TODO in code

```bash
# 1. File bug
BUG_ID="BUG-20251006-push-stub"
cat > bugs/filed/${BUG_ID}.md << 'EOF'
# Bug: Push command stub not implemented
**Severity**: CRITICAL
**File**: pkg/cmd/push/root.go:69
**Issue**: Function returns nil without implementing actual push logic
EOF

git add bugs/filed/${BUG_ID}.md
git commit -m "bug: File ${BUG_ID} - push command stub"
git push

# 2. Check fix cascade
if ! jq -e ".fix_cascades[] | select(.bugs[] | contains(\"${BUG_ID}\"))" orchestrator-state-v3.json >/dev/null; then
    echo "Bug not in fix cascade, creating..."

    # 3. Create cascade
    CASCADE_ID="FC-20251006-001"
    jq ".fix_cascades += [{
        \"cascade_id\": \"${CASCADE_ID}\",
        \"bugs\": [\"${BUG_ID}\"],
        \"severity\": \"CRITICAL\",
        \"status\": \"FILED\",
        \"assigned_to\": \"ERROR_RECOVERY\",
        \"created_at\": \"$(date -Iseconds)\"
    }]" orchestrator-state-v3.json > tmp && mv tmp orchestrator-state-v3.json

    git add orchestrator-state-v3.json
    git commit -m "state: Add fix cascade ${CASCADE_ID} for ${BUG_ID}"
    git push

    # 4. Block review
    echo "🔴🔴🔴 REVIEW BLOCKED: Critical bug unfixed 🔴🔴🔴"
    echo "Bug: ${BUG_ID}"
    echo "Cascade: ${CASCADE_ID}"
    echo "Status: Filed, awaiting ERROR_RECOVERY"

    exit 332
fi
```

**Result**: Review blocked until bug fixed. No "pre-existing" excuse accepted.

---

## 📚 TODO VERIFICATION EXAMPLES

### Example 1: TODO Properly Verified and Accepted

**File**: pkg/retry/backoff.go:89
**TODO**: // TODO: Add jitter to backoff (Phase 4 Wave 1)

**Verification**:
```bash
$ grep -rn "jitter.*backoff" phase-plans/
phase-plans/phase4-wave1-plan.md:67: E4.1.3 - Add Backoff Jitter
phase-plans/phase4-wave1-plan.md:68: Implement jitter in pkg/retry/backoff.go line 89
```

**Plan Excerpt**:
```
67: ## E4.1.3 - Add Backoff Jitter
68: Implement jitter in pkg/retry/backoff.go line 89
69: Add randomization to prevent thundering herd
```

**Analysis**:
- Current: Phase 2, Wave 3
- Planned: Phase 4, Wave 1
- Status: FUTURE ✅

**Decision**: ✅ ACCEPTED - Explicitly planned for Phase 4 Wave 1

---

### Example 2: TODO Without Plan - BUG Filed

**File**: pkg/push/operations.go:69
**TODO**: // TODO: Implement actual push logic in Phase 2

**Search**:
```bash
$ grep -rn "push.*logic\|actual.*push" phase-plans/ wave-plans/
(no results)
```

**Analysis**:
- Current: Phase 2, Wave 3
- TODO says: "in Phase 2" ← We're IN Phase 2!
- No plan found for implementation

**Decision**: 🔴 BUG - Filed as BUG-20251006-push-not-implemented
**Reason**: TODO references current phase but not implemented

**Bug Filing**:
```bash
BUG_ID="BUG-20251006-143000-push-not-implemented"
cat > bugs/filed/${BUG_ID}.md << 'EOF'
# Bug: Push Logic Not Implemented

**Bug ID**: BUG-20251006-143000-push-not-implemented
**Severity**: CRITICAL
**Component**: pkg/push/operations.go
**Line**: 69
**Found During**: Phase 2 Wave 3 Code Review
**Found By**: code-reviewer
**Filed At**: 2025-10-06T14:30:00Z

## Issue Description
TODO comment indicates push logic should be implemented in Phase 2,
but we are currently IN Phase 2 Wave 3 and it's still not implemented.

## Evidence
```go
// TODO: Implement actual push logic in Phase 2
func PushImages(cmd *cobra.Command, args []string) error {
    return nil
}
```

## Impact
Push command is completely non-functional - returns without doing anything.

## Root Cause
Implementation deferred but never completed in planned phase.

## Required Fix
Implement actual push logic immediately.

## Status
**Current**: FILED
**Assigned To**: ERROR_RECOVERY
**Fix Cascade**: FC-20251006-143000
**Priority**: CRITICAL
EOF

git add bugs/filed/${BUG_ID}.md
git commit -m "bug: File ${BUG_ID} - push logic not implemented"
git push
```

---

### Example 3: TODO for Past Phase - BUG Filed

**File**: pkg/auth/handler.go:45
**TODO**: // TODO: Add authentication (Phase 1 Wave 2)

**Search**:
```bash
$ grep -rn "authentication" phase-plans/phase1-wave2-plan.md
123: E1.2.2 - Authentication Implementation
124: Add authentication to pkg/auth/handler.go
```

**Analysis**:
- Current: Phase 2, Wave 3
- Planned: Phase 1, Wave 2 ← PAST!
- Status: SHOULD HAVE BEEN DONE

**Decision**: 🔴 BUG - Filed as BUG-20251006-auth-missing
**Reason**: Planned for past phase but never implemented

**Bug Filing**:
```bash
BUG_ID="BUG-20251006-143100-auth-missing"
cat > bugs/filed/${BUG_ID}.md << 'EOF'
# Bug: Authentication Not Implemented (Past Due)

**Bug ID**: BUG-20251006-143100-auth-missing
**Severity**: CRITICAL
**Component**: pkg/auth/handler.go
**Line**: 45
**Found During**: Phase 2 Wave 3 Code Review
**Found By**: code-reviewer
**Filed At**: 2025-10-06T14:31:00Z

## Issue Description
TODO comment indicates authentication should have been implemented
in Phase 1 Wave 2, but we are now in Phase 2 Wave 3 and it's still missing.

## Evidence
```go
// TODO: Add authentication (Phase 1 Wave 2)
func Authenticate(token string) error {
    return nil  // Stub - not implemented
}
```

**Plan Evidence**:
```
phase-plans/phase1-wave2-plan.md:123: E1.2.2 - Authentication Implementation
phase-plans/phase1-wave2-plan.md:124: Add authentication to pkg/auth/handler.go
```

## Impact
Security vulnerability - no authentication enforcement.

## Root Cause
Planned implementation in Phase 1 Wave 2 was never completed.
TODO carried forward through multiple phases without being addressed.

## Required Fix
Implement authentication immediately as critical security issue.

## Status
**Current**: FILED
**Assigned To**: ERROR_RECOVERY
**Fix Cascade**: FC-20251006-143100
**Priority**: CRITICAL (Security)
EOF

git add bugs/filed/${BUG_ID}.md
git commit -m "bug: File ${BUG_ID} - authentication not implemented (past due)"
git push
```

---

### Example 4: Vague TODO Acceptance (PROHIBITED)

**WRONG Approach (R332 Violation = -100%):**

```markdown
## Code Review: pkg/feature/handler.go

**TODO Found**: Line 89 - "TODO: Add retry logic"

**Assessment**: ✅ PASS
**Reason**: "This TODO is acceptable because it will be addressed in a future phase."

Evidence: None provided
Plan File: Not specified
Phase/Wave: "Later"
```

**PENALTY: -100% IMMEDIATE FAILURE**

**Violations**:
- No plan file specified
- No line numbers from plan
- No grep evidence
- Vague "future phase" (not specific)
- No phase/wave comparison

---

**CORRECT Approach:**

```markdown
## TODO Analysis: pkg/feature/handler.go:89

**TODO Text**: `// TODO: Add retry logic`

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
156: ## E3.2.1 - Retry Mechanism Implementation
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

**Decision**: ✅ ACCEPTED
**Reason**: Explicitly planned for Phase 3 Wave 2 per phase3-wave2-plan.md lines 156-178

**Tracking**: Added to orchestrator-state-v3.json todo_tracking section for verification in Phase 3
```

---

## 🚨🚨🚨 FINAL WARNING 🚨🚨🚨

**"PRE-EXISTING" IS ELIMINATED FROM THE SOFTWARE FACTORY.**

Every bug found MUST be:
1. Filed immediately
2. Tracked in fix cascade
3. Fixed or explicitly managed
4. Never dismissed as "someone else's problem"

**ZERO TOLERANCE FOR UNTRACKED BUGS.**

The age of a bug is IRRELEVANT. If it exists, it's OUR bug, and we WILL fix it.

---

**Last Updated**: 2025-10-07
**Rule Owner**: Factory Manager
**Enforcement Level**: ABSOLUTE
**Exit Code**: 332 (R332 Violation - Untracked Bug)
