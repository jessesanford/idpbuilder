# 🔴🔴🔴 RULE R616 - BUG LIFECYCLE TRACKING PROTOCOL (SUPREME LAW)

**Criticality:** SUPREME LAW - REQUIRED FOR R615 PROGRESS TRACKING
**Grading Impact:** -100% for improper bug tracking
**Enforcement:** CONTINUOUS - Every bug discovery and closure

---

## SUPREME LAW STATEMENT

**EVERY BUG MUST HAVE A UNIQUE, STABLE IDENTIFIER AND COMPLETE LIFECYCLE TRACKING. BUGS MUST TRANSITION THROUGH DEFINED STATES (OPEN → CLOSED → REOPENED). PROGRESS MEASUREMENT DEPENDS ON ACCURATE BUG TRACKING. IMPROPER TRACKING CORRUPTS PROGRESS ANALYSIS.**

---

## 🚨🚨🚨 THE BUG TRACKING MANDATE 🚨🚨🚨

### Why Bug Tracking is Critical

**Problems Without Proper Tracking:**
- Can't tell if bugs are "same" or "new" across iterations
- Can't measure actual progress (bugs fixed)
- Can't detect flapping (reopened bugs)
- Can't apply R615 two-tiered limits
- Leads to infinite loops or premature escalation

**Requirements for Progress Tracking:**
1. **Stable Bug IDs** - Survive code changes, persist across iterations
2. **Lifecycle States** - Track bug journey (OPEN → CLOSED → REOPENED)
3. **Ownership** - Clear responsibility for bug discovery and closure
4. **Verification** - Evidence that bugs are actually fixed
5. **History** - Complete audit trail of bug lifecycle

---

## 🔴🔴🔴 BUG IDENTIFIER FORMAT 🔴🔴🔴

### Required ID Structure

**Format:** `{SCOPE}-BUG-{SEQUENCE}`

**Scope Types:**

#### Effort-Level Bugs
```
Format: E{phase}.{wave}.{effort}-BUG-{num}
Examples:
  - E1.2.3-BUG-001  (Effort 3, Wave 2, Phase 1)
  - E2.1.1-BUG-042  (Effort 1, Wave 1, Phase 2)

When to use:
  - Bugs found during effort implementation
  - Bugs found in effort code review
  - Bugs specific to single effort
```

#### Wave-Level Bugs
```
Format: W{wave}-BUG-{num}
Examples:
  - W2-BUG-001  (Wave 2 integration bug)
  - W1-BUG-042  (Wave 1 integration bug)

When to use:
  - Bugs found during wave integration
  - Bugs that span multiple efforts in wave
  - Bugs in wave integration logic
```

#### Phase-Level Bugs
```
Format: P{phase}-BUG-{num}
Examples:
  - P1-BUG-007  (Phase 1 integration bug)
  - P2-BUG-023  (Phase 2 integration bug)

When to use:
  - Bugs found during phase integration
  - Bugs that span multiple waves
  - Bugs in phase integration logic
```

#### Project-Level Bugs
```
Format: PROJ-BUG-{num}
Examples:
  - PROJ-BUG-001  (Project-level bug)
  - PROJ-BUG-099  (Project-level bug)

When to use:
  - Bugs found during final project integration
  - Bugs that affect entire project
  - Cross-cutting concerns
```

### ID Assignment Protocol

**Who Assigns IDs:**
- **Code Reviewer** assigns IDs when discovering bugs
- Reviewers maintain sequence counter per scope
- IDs are assigned sequentially within scope

**When Assigned:**
- At bug discovery time (during review)
- Before bug is reported to orchestrator
- Recorded immediately in bug-tracking.json

**ID Requirements:**
```yaml
must_be:
  - Unique across entire project
  - Stable (doesn't change with code changes)
  - Scoped appropriately (effort/wave/phase/project)
  - Sequential within scope (for tracking)
  - Greppable (simple string format)

must_not_be:
  - Based on file paths (files can be renamed)
  - Based on line numbers (lines change)
  - Based on description hashing (fragile)
  - Generated randomly (can't sequence)
```

---

## 🔴🔴🔴 BUG LIFECYCLE STATES 🔴🔴🔴

### State Definitions

#### 1. OPEN
```yaml
definition: "Bug newly discovered and not yet fixed"

transitions_from:
  - NULL (initial state - bug just found)

transitions_to:
  - CLOSED (when fix verified)
  - (stays OPEN through iterations if not fixed)

fields_required:
  - bug_id
  - status: "OPEN"
  - first_found_iteration
  - description
  - affected_files
  - affected_efforts
  - severity
  - discovered_by (reviewer_id)

example:
  bug_id: "W2-BUG-001"
  status: "OPEN"
  first_found_iteration: 3
  description: "Authentication fails on concurrent requests"
  affected_files: ["pkg/auth/handler.go"]
  affected_efforts: ["2.2.1", "2.2.3"]
  severity: "HIGH"
  discovered_by: "agent-code-reviewer-20251102-120000"
```

#### 2. CLOSED
```yaml
definition: "Bug fixed and verified by code reviewer"

transitions_from:
  - OPEN (when fix implemented and verified)

transitions_to:
  - REOPENED (if bug reappears in later iteration)
  - (stays CLOSED permanently if fix holds)

fields_required:
  - bug_id
  - status: "CLOSED"
  - first_found_iteration
  - closed_iteration
  - closed_by (reviewer_id)
  - fix_commits (array of commit hashes)
  - verification_evidence
  - description (original)

example:
  bug_id: "W2-BUG-001"
  status: "CLOSED"
  first_found_iteration: 3
  closed_iteration: 5
  closed_by: "agent-code-reviewer-20251102-140000"
  fix_commits: ["abc123ef", "def456ab"]
  verification_evidence: "Tested with 100 concurrent requests, no auth failures. Added integration test covering race condition."
  description: "Authentication fails on concurrent requests"
```

#### 3. REOPENED
```yaml
definition: "Bug previously closed but found again (flapping)"

transitions_from:
  - CLOSED (when previously fixed bug reappears)

transitions_to:
  - CLOSED (when fix verified again)
  - (stays REOPENED until re-fixed)

fields_required:
  - bug_id
  - status: "REOPENED"
  - first_found_iteration
  - closed_iteration (previous closure)
  - reopened_iteration (when it came back)
  - reopen_count (how many times reopened)
  - reopen_reason (why did it come back?)
  - description (original)

example:
  bug_id: "W2-BUG-001"
  status: "REOPENED"
  first_found_iteration: 3
  closed_iteration: 5
  reopened_iteration: 7
  reopen_count: 1
  reopen_reason: "Original fix only handled sync case, async case still fails"
  description: "Authentication fails on concurrent requests"

flags:
  - 🚨 CRITICAL: Reopened bugs indicate incomplete fixes
  - Must analyze root cause deeply
  - Consider architect review if reopen_count >= 3
```

### State Transition Rules

```
          ┌─────────────────────┐
          │   Initial State     │
          │     (NULL)          │
          └─────────┬───────────┘
                    │ Bug discovered
                    ↓
          ┌─────────────────────┐
          │       OPEN          │←──────────────┐
          │  (needs fixing)     │               │
          └─────────┬───────────┘               │
                    │ Fix verified              │
                    ↓                           │
          ┌─────────────────────┐               │ Bug reappears
          │      CLOSED         │               │ (flapping)
          │  (fix verified)     │               │
          └─────────┬───────────┘               │
                    │ Bug reappears             │
                    ↓                           │
          ┌─────────────────────┐               │
          │     REOPENED        │───────────────┘
          │   (flapping!)       │    Fix verified again
          └─────────────────────┘    (can cycle)
```

---

## 🔴🔴🔴 BUG TRACKING FILE FORMAT 🔴🔴🔴

### bug-tracking.json Schema

**Location:** `{effort_workspace}/bug-tracking.json` or `{integration_workspace}/bug-tracking.json`

```json
{
  "metadata": {
    "scope": "WAVE|PHASE|PROJECT|EFFORT",
    "scope_id": "W2|P1|PROJ|E1.2.3",
    "last_updated": "2025-11-02T14:30:00Z",
    "last_updated_by": "agent-code-reviewer-20251102-140000"
  },
  "bugs": [
    {
      "bug_id": "W2-BUG-001",
      "status": "OPEN|CLOSED|REOPENED",
      "severity": "LOW|MEDIUM|HIGH|CRITICAL",

      "lifecycle": {
        "first_found_iteration": 3,
        "first_found_date": "2025-11-02T12:00:00Z",
        "discovered_by": "agent-code-reviewer-20251102-120000",

        "closed_iteration": 5,
        "closed_date": "2025-11-02T14:00:00Z",
        "closed_by": "agent-code-reviewer-20251102-140000",

        "reopened_iteration": 7,
        "reopened_date": "2025-11-02T16:00:00Z",
        "reopen_count": 1,
        "reopen_reason": "Original fix incomplete"
      },

      "description": {
        "summary": "Authentication fails on concurrent requests",
        "details": "When multiple requests hit auth handler simultaneously, mutex is not properly acquired, leading to race condition and auth failures.",
        "reproduction_steps": [
          "Send 100 concurrent POST requests to /api/auth/login",
          "Observe ~10% failure rate",
          "Check logs for 'nil pointer dereference' errors"
        ]
      },

      "impact": {
        "affected_files": [
          "pkg/auth/handler.go",
          "pkg/auth/session.go"
        ],
        "affected_efforts": ["2.2.1", "2.2.3"],
        "affected_tests": ["TestConcurrentAuth"],
        "user_impact": "HIGH"
      },

      "fix": {
        "fix_commits": ["abc123ef", "def456ab"],
        "fix_branches": ["project-name/phase1/wave2/effort-2.2.1"],
        "fix_description": "Added proper mutex locking around session map access. Added integration test for concurrent auth.",
        "verification_evidence": "Tested with 100 concurrent requests, no auth failures. Added integration test covering race condition. Test suite passes 100 consecutive runs."
      }
    }
  ],

  "statistics": {
    "total_bugs_found": 15,
    "total_bugs_closed": 8,
    "total_bugs_open": 6,
    "total_bugs_reopened": 1,
    "bugs_by_severity": {
      "CRITICAL": 1,
      "HIGH": 4,
      "MEDIUM": 7,
      "LOW": 3
    },
    "average_time_to_close_hours": 4.5
  }
}
```

---

## 🔴🔴🔴 BUG DISCOVERY PROTOCOL 🔴🔴🔴

### When Code Reviewer Finds Bug

**Step 1: Assign Bug ID**
```bash
# Determine scope
SCOPE="W2"  # Wave 2 integration

# Get next sequence number
LAST_BUG=$(jq -r '.bugs[] | select(.bug_id | startswith("'$SCOPE'-BUG-")) | .bug_id' bug-tracking.json | tail -1)
LAST_NUM=$(echo "$LAST_BUG" | grep -o '[0-9]*$')
NEXT_NUM=$(printf "%03d" $((LAST_NUM + 1)))

BUG_ID="${SCOPE}-BUG-${NEXT_NUM}"
echo "Assigned bug ID: $BUG_ID"
```

**Step 2: Document Bug**
```bash
# Create bug entry
jq --arg id "$BUG_ID" \
   --arg summary "Authentication fails on concurrent requests" \
   --arg severity "HIGH" \
   --argjson iteration "$CURRENT_ITERATION" \
   --arg reviewer "$REVIEWER_ID" \
   '.bugs += [{
     bug_id: $id,
     status: "OPEN",
     severity: $severity,
     lifecycle: {
       first_found_iteration: $iteration,
       first_found_date: "'$(date -Iseconds)'",
       discovered_by: $reviewer
     },
     description: {
       summary: $summary,
       details: "...",
       reproduction_steps: [...]
     },
     impact: {
       affected_files: [...],
       affected_efforts: [...],
       user_impact: "HIGH"
     }
   }]' bug-tracking.json > tmp.json && mv tmp.json bug-tracking.json
```

**Step 3: Report Bug**
```bash
# Add to review report
echo "## Bug Found: $BUG_ID" >> REVIEW-REPORT.md
echo "" >> REVIEW-REPORT.md
echo "**Summary:** Authentication fails on concurrent requests" >> REVIEW-REPORT.md
echo "**Severity:** HIGH" >> REVIEW-REPORT.md
echo "**Affected Files:** pkg/auth/handler.go" >> REVIEW-REPORT.md
echo "**Reproduction:** Send 100 concurrent POST to /api/auth/login" >> REVIEW-REPORT.md
echo "" >> REVIEW-REPORT.md
```

---

## 🔴🔴🔴 BUG CLOSURE PROTOCOL 🔴🔴🔴

### When Code Reviewer Verifies Fix

**Step 1: Verify Fix (Mandatory)**
```bash
# 1. Tests must pass (prerequisite)
if ! make test; then
    echo "❌ Tests failing - cannot close bug"
    exit 1
fi

# 2. Manually verify the specific bug is fixed
echo "🔍 Verifying fix for $BUG_ID..."

# Run reproduction steps
# Check that bug no longer occurs
# Document verification evidence
```

**Step 2: Mark Bug CLOSED**
```bash
BUG_ID="W2-BUG-001"
CURRENT_ITERATION=5
REVIEWER_ID="agent-code-reviewer-20251102-140000"
FIX_COMMITS=("abc123ef" "def456ab")
VERIFICATION="Tested with 100 concurrent requests, no auth failures. Added integration test covering race condition."

jq --arg id "$BUG_ID" \
   --argjson iter "$CURRENT_ITERATION" \
   --arg reviewer "$REVIEWER_ID" \
   --argjson commits "$(printf '%s\n' "${FIX_COMMITS[@]}" | jq -R . | jq -s .)" \
   --arg evidence "$VERIFICATION" \
   '.bugs |= map(
     if .bug_id == $id then
       .status = "CLOSED" |
       .lifecycle.closed_iteration = $iter |
       .lifecycle.closed_date = "'$(date -Iseconds)'" |
       .lifecycle.closed_by = $reviewer |
       .fix.fix_commits = $commits |
       .fix.verification_evidence = $evidence
     else
       .
     end
   )' bug-tracking.json > tmp.json && mv tmp.json bug-tracking.json
```

**Step 3: Document Closure**
```bash
echo "✅ CLOSED: $BUG_ID" >> REVIEW-REPORT.md
echo "  Fixed in commits: ${FIX_COMMITS[@]}" >> REVIEW-REPORT.md
echo "  Verification: $VERIFICATION" >> REVIEW-REPORT.md
```

### Closure Requirements

```yaml
cannot_close_without:
  - All tests passing (automated)
  - Manual verification completed
  - Verification evidence documented
  - Fix commits identified
  - Code reviewer approval

verification_evidence_must_include:
  - What was tested
  - How it was tested
  - Results observed
  - Why confident bug is fixed
  - Any new tests added

example_good_evidence:
  "Tested with 100 concurrent requests over 10 runs, zero failures.
   Added integration test TestConcurrentAuth that reproduces original bug
   and passes with fix. Reviewed mutex usage, confirmed proper locking.
   No related code paths have similar issues."

example_bad_evidence:
  "Fixed it"  ❌
  "Tests pass"  ❌
  "Looks good"  ❌
```

---

## 🔴🔴🔴 BUG REOPENING PROTOCOL 🔴🔴🔴

### When Previously-Fixed Bug Reappears

**Step 1: Detect Reopening**
```bash
# Usually detected by bug-progress-analyzer.sh
# But reviewer may also notice during review

BUG_ID="W2-BUG-001"

# Check if bug was previously closed
PREVIOUSLY_CLOSED=$(jq -r ".bugs[] | select(.bug_id==\"$BUG_ID\" and .status==\"CLOSED\") | .bug_id" bug-tracking.json)

if [ -n "$PREVIOUSLY_CLOSED" ]; then
    echo "🚨 BUG FLAPPING DETECTED: $BUG_ID was previously closed!"
fi
```

**Step 2: Mark as REOPENED**
```bash
REOPEN_REASON="Original fix only handled sync case, async case still fails"

jq --arg id "$BUG_ID" \
   --argjson iter "$CURRENT_ITERATION" \
   --arg reason "$REOPEN_REASON" \
   '.bugs |= map(
     if .bug_id == $id then
       .status = "REOPENED" |
       .lifecycle.reopened_iteration = $iter |
       .lifecycle.reopened_date = "'$(date -Iseconds)'" |
       .lifecycle.reopen_count = (.lifecycle.reopen_count // 0) + 1 |
       .lifecycle.reopen_reason = $reason
     else
       .
     end
   )' bug-tracking.json > tmp.json && mv tmp.json bug-tracking.json
```

**Step 3: Analyze Root Cause**
```bash
# If reopen_count >= 3, this is critical
REOPEN_COUNT=$(jq -r ".bugs[] | select(.bug_id==\"$BUG_ID\") | .lifecycle.reopen_count" bug-tracking.json)

if [ "$REOPEN_COUNT" -ge 3 ]; then
    echo "🚨🚨🚨 CRITICAL: Bug $BUG_ID has reopened $REOPEN_COUNT times!"
    echo "This indicates systemic issue - escalating to architect"

    # Flag for architect review
    # Create root cause analysis task
fi
```

**Step 4: Document Flapping**
```bash
echo "🚨 REOPENED: $BUG_ID (reopen count: $REOPEN_COUNT)" >> REVIEW-REPORT.md
echo "  Reason: $REOPEN_REASON" >> REVIEW-REPORT.md
echo "  Original fix: $(jq -r ".bugs[] | select(.bug_id==\"$BUG_ID\") | .fix.fix_description" bug-tracking.json)" >> REVIEW-REPORT.md
echo "  ACTION REQUIRED: Deep root cause analysis" >> REVIEW-REPORT.md
```

---

## 🔴🔴🔴 BUG TRACKING VALIDATION 🔴🔴🔴

### Pre-Commit Validation

```bash
# File: .pre-commit/validate-bug-tracking.sh

validate_bug_tracking() {
    local bug_file="$1"
    local errors=0

    echo "🔍 R616: Validating bug tracking..."

    # Check all bug IDs are unique
    local duplicate_ids=$(jq -r '.bugs[].bug_id' "$bug_file" | sort | uniq -d)
    if [ -n "$duplicate_ids" ]; then
        echo "❌ R616 VIOLATION: Duplicate bug IDs found:"
        echo "$duplicate_ids"
        ((errors++))
    fi

    # Check all bug IDs follow format
    local invalid_ids=$(jq -r '.bugs[].bug_id' "$bug_file" | grep -v '^[EWP][0-9.]*-BUG-[0-9]\{3\}$\|^PROJ-BUG-[0-9]\{3\}$')
    if [ -n "$invalid_ids" ]; then
        echo "❌ R616 VIOLATION: Invalid bug ID format:"
        echo "$invalid_ids"
        ((errors++))
    fi

    # Check all bugs have required fields
    jq -r '.bugs[] | select(.status == null or .bug_id == null or .lifecycle.first_found_iteration == null) | .bug_id // "UNKNOWN"' "$bug_file" | while read bug_id; do
        echo "❌ R616 VIOLATION: Bug $bug_id missing required fields"
        ((errors++))
    done

    # Check CLOSED bugs have verification evidence
    jq -r '.bugs[] | select(.status == "CLOSED" and (.fix.verification_evidence == null or .fix.verification_evidence == "")) | .bug_id' "$bug_file" | while read bug_id; do
        echo "❌ R616 VIOLATION: Closed bug $bug_id missing verification evidence"
        ((errors++))
    done

    # Check REOPENED bugs have reopen_reason
    jq -r '.bugs[] | select(.status == "REOPENED" and (.lifecycle.reopen_reason == null or .lifecycle.reopen_reason == "")) | .bug_id' "$bug_file" | while read bug_id; do
        echo "❌ R616 VIOLATION: Reopened bug $bug_id missing reopen reason"
        ((errors++))
    done

    if [ $errors -eq 0 ]; then
        echo "✅ R616: Bug tracking validation passed"
        return 0
    else
        echo "❌ R616: Found $errors bug tracking issues"
        return 1
    fi
}
```

---

## 🔴🔴🔴 INTEGRATE_WAVE_EFFORTS WITH R615 🔴🔴🔴

### How R616 Enables R615

**R615 Dependency on R616:**
```yaml
r533_requires_r534_for:
  bug_identity:
    - Must know if bug is "same" vs "new" across iterations
    - R616 stable bug IDs enable this

  progress_measurement:
    - Must count bugs_closed accurately
    - R616 CLOSED status provides this

  flapping_detection:
    - Must detect when bugs reappear
    - R616 REOPENED status tracks this

  verification:
    - Must confirm bugs are actually fixed
    - R616 verification protocol ensures this
```

**Data Flow:**
```
R616 (Bug Tracking)
      ↓
    bug-tracking.json
      ↓
R615 (Progress Analysis)
      ↓
tools/bug-progress-analyzer.sh
      ↓
Iteration Decision (continue or ERROR_RECOVERY)
```

---

## 🔴🔴🔴 GRADING IMPACT 🔴🔴🔴

### Compliance Criteria

```yaml
r534_compliance:
  bug_id_format: 20%            # All IDs follow format
  unique_ids: 20%               # No duplicate IDs
  complete_lifecycle: 20%       # All required fields present
  verification_evidence: 20%    # Closed bugs have evidence
  reopening_detection: 10%      # Flapping tracked correctly
  ownership_clear: 10%          # Discoverer and closer recorded

violations:
  missing_bug_id: -50%              # Bug without ID
  duplicate_bug_id: -100%           # ID collision
  invalid_id_format: -30%           # Wrong format
  closed_without_evidence: -40%     # No verification
  missing_lifecycle_fields: -20%    # Incomplete data
  reopening_not_detected: -50%      # Flapping missed
```

### Automatic Failure Conditions

```yaml
automatic_failure:
  - Duplicate bug IDs in system
  - Bugs closed without verification evidence
  - Reopened bugs not detected
  - Invalid bug ID formats
  - Missing required lifecycle fields
```

---

## 🔴 SUMMARY: BUG TRACKING PRINCIPLES

```
┌─────────────────────────────────────────────────────────────┐
│            BUG LIFECYCLE TRACKING PROTOCOL                  │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  BUG IDENTITY:                                              │
│  ═════════════                                              │
│  Format: {SCOPE}-BUG-{SEQUENCE}                            │
│  Scopes: E{p}.{w}.{e}, W{w}, P{p}, PROJ                    │
│  Stable: Survives code changes                             │
│  Unique: No collisions                                      │
│                                                             │
│  BUG LIFECYCLE:                                             │
│  ══════════════                                             │
│  OPEN      → Bug discovered, needs fixing                   │
│  CLOSED    → Bug fixed and verified                         │
│  REOPENED  → Bug came back (flapping!)                      │
│                                                             │
│  BUG OWNERSHIP:                                             │
│  ══════════════                                             │
│  Discovery: Code Reviewer assigns ID                        │
│  Closure:   Code Reviewer verifies fix                      │
│  Evidence:  Required for CLOSED status                      │
│                                                             │
│  BUG FLAPPING:                                              │
│  ═════════════                                              │
│  Detection: Bug reappears after CLOSED                      │
│  Status:    Changed to REOPENED                             │
│  Action:    Deep root cause analysis                        │
│  Critical:  If reopen_count >= 3                            │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

**Remember:** Bug tracking is the foundation for progress measurement. Without stable bug IDs and lifecycle tracking, we cannot tell if we're making progress or just churning. R616 provides the data that R615 uses to make intelligent iteration decisions.

**See Also:**
- R615: Progress-Based Iteration Limits (uses R616 data)
- tools/bug-progress-analyzer.sh (analyzes R616 data)
- bug-tracking.json schema (stores R616 data)
- orchestrator-state-v3.json bugs section (tracks R616 history)
- Code Reviewer states (implement R616 protocols)
