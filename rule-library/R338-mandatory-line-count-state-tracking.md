# 🚨🚨🚨 RULE R338 - Mandatory Line Count State Tracking

**Criticality:** BLOCKING - Cannot proceed without tracking  
**Grading Impact:** -50% per missing tracking, -100% if not tracked at all  
**Enforcement:** CONTINUOUS - Every effort and split MUST have line count tracking
**Scope:** ORCHESTRATOR (PRIMARY), Code Reviewer (REPORTER)

## 🔴🔴🔴 ABSOLUTE REQUIREMENT 🔴🔴🔴

**EVERY effort and split in orchestrator-state.json MUST have complete line_count_tracking!**

## Rule Statement

The Orchestrator MUST maintain comprehensive line count tracking for EVERY effort and split in orchestrator-state.json. Code Reviewers MUST report line counts in a standardized format that the Orchestrator captures and records. This tracking is the SINGLE SOURCE OF TRUTH for size compliance.

## Mandatory Tracking Structure

### For EVERY Effort and Split:
```json
"line_count_tracking": {
  "initial_count": null,              // First measurement from Code Reviewer
  "current_count": null,              // Latest measurement (after fixes/updates)
  "last_measured": null,              // ISO timestamp of last measurement
  "measured_by": null,                // Which agent measured (always "code-reviewer")
  "measurement_command": null,        // Exact command used for measurement
  "auto_detected_base": null,         // Base branch auto-detected by tool
  "implementation_only": true,        // Per R007: excludes tests/demos/docs
  "within_limit": null,               // Boolean: true if ≤800, false if >800
  "requires_split": false,            // True if >800 lines detected (or >900 with grace)
  "split_plan_path": null,           // Path to SPLIT-PLAN.md if splitting
  
  // Fix Grace Period Tracking (R339)
  "fix_delta": null,                  // Lines added by fixes (not initial implementation)
  "grace_period_eligible": false,     // True if eligible for 900-line grace period
  "grace_period_applied": false,      // True if grace period was used
  "grace_period_threshold": 800,      // 900 if grace applied, 800 otherwise
  "fix_history": [],                  // History of fixes applied
  
  "measurement_history": [            // Track all measurements
    {
      "timestamp": "2025-01-20T10:30:00Z",
      "count": 687,
      "measured_by": "code-reviewer",
      "reason": "initial review",
      "command": "./tools/line-counter.sh phase1/wave1/effort1",
      "auto_detected_base": "phase1-wave1-integration",
      "excludes": "tests/demos/docs per R007",
      "within_limit": true
    }
  ]
}
```

## 🚨 ORCHESTRATOR RESPONSIBILITIES (MANDATORY)

### 1. Capture Line Counts from Code Reviewer Reports
```bash
# When Code Reviewer completes review, extract line count:
parse_code_review_report() {
    local report_path="$1"
    
    # Extract standardized line count section
    LINE_COUNT=$(grep -A5 "SIZE MEASUREMENT REPORT" "$report_path" | grep "Implementation Lines:" | awk '{print $3}')
    COMMAND=$(grep -A5 "SIZE MEASUREMENT REPORT" "$report_path" | grep "Command:" | cut -d':' -f2-)
    BASE=$(grep -A5 "SIZE MEASUREMENT REPORT" "$report_path" | grep "Auto-detected Base:" | cut -d':' -f2-)
    
    # Update orchestrator-state.json
    update_line_count_tracking "$effort_name" "$LINE_COUNT" "$COMMAND" "$BASE"
}
```

### 2. Update State File with Line Count Tracking
```bash
update_line_count_tracking() {
    local effort="$1"
    local count="$2"
    local command="$3"
    local base="$4"
    
    # Create tracking structure if not exists
    jq ".efforts_in_progress[] |= 
        if .name == \"$effort\" then
            .line_count_tracking = {
                initial_count: (if .line_count_tracking.initial_count == null then $count else .line_count_tracking.initial_count end),
                current_count: $count,
                last_measured: \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\",
                measured_by: \"code-reviewer\",
                measurement_command: \"$command\",
                auto_detected_base: \"$base\",
                implementation_only: true,
                within_limit: ($count <= 800),
                requires_split: ($count > 800),
                measurement_history: (.line_count_tracking.measurement_history + [{
                    timestamp: \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\",
                    count: $count,
                    measured_by: \"code-reviewer\",
                    reason: \"code review\",
                    command: \"$command\",
                    auto_detected_base: \"$base\",
                    excludes: \"tests/demos/docs per R007\",
                    within_limit: ($count <= 800)
                }])
            }
        else . end" orchestrator-state.json > tmp.json && mv tmp.json orchestrator-state.json
}
```

### 3. Track Fix Grace Period Applications (R339)
```bash
track_fix_with_grace_period() {
    local effort="$1"
    local original_count="$2"
    local fix_lines="$3"
    local total_count=$((original_count + fix_lines))
    
    # Determine if grace period applies
    local grace_eligible="false"
    local grace_applied="false"
    local threshold=800
    
    if [ "$total_count" -lt 900 ] && [ "$fix_lines" -gt 0 ]; then
        grace_eligible="true"
        grace_applied="true"
        threshold=900
    fi
    
    # Update tracking with grace period info
    jq ".efforts_completed[\"$effort\"].line_count_tracking |= . + {
        fix_delta: $fix_lines,
        grace_period_eligible: $grace_eligible,
        grace_period_applied: $grace_applied,
        grace_period_threshold: $threshold,
        fix_history: .fix_history + [{
            fix_id: \"FIX-$(date +%s)\",
            source: \"integration-failure\",
            lines_added: $fix_lines,
            total_after_fix: $total_count,
            within_grace: ($total_count < $threshold),
            timestamp: \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\"
        }]
    }" orchestrator-state.json > tmp.json && mv tmp.json orchestrator-state.json
}
```

### 4. Validate Line Count Tracking Before State Transitions
```bash
validate_line_count_tracking() {
    local effort="$1"
    
    # Check if tracking exists
    TRACKING=$(jq ".efforts_in_progress[] | select(.name == \"$effort\") | .line_count_tracking" orchestrator-state.json)
    
    if [ "$TRACKING" = "null" ]; then
        echo "🚨 R338 VIOLATION: No line count tracking for $effort!"
        echo "❌ CANNOT proceed without line count tracking"
        return 1
    fi
    
    # Check if measurement is current
    LAST_MEASURED=$(echo "$TRACKING" | jq -r '.last_measured')
    CURRENT_COUNT=$(echo "$TRACKING" | jq -r '.current_count')
    
    if [ "$CURRENT_COUNT" = "null" ]; then
        echo "🚨 R338 VIOLATION: No line count recorded for $effort!"
        return 1
    fi
    
    echo "✅ Line count tracking valid: $CURRENT_COUNT lines (measured: $LAST_MEASURED)"
    return 0
}
```

## 🚨 CODE REVIEWER RESPONSIBILITIES (MANDATORY)

### Standardized Line Count Reporting Format
Every Code Reviewer report MUST include:

```markdown
## 📊 SIZE MEASUREMENT REPORT
**Implementation Lines:** 687
**Command:** ./tools/line-counter.sh phase1/wave1/effort1  
**Auto-detected Base:** phase1-wave1-integration
**Timestamp:** 2025-01-20T10:30:00Z
**Within Limit:** ✅ Yes (687 < 800)
**Excludes:** tests/demos/docs per R007

### Raw Output:
```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
📊 Line Counter - Software Factory 2.0
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
📌 Analyzing branch: phase1/wave1/effort1
🎯 Detected base:    phase1-wave1-integration
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
✅ Total non-generated lines: 687
```
```

## 🔴 VALIDATION GATES

### Gate 1: Before Marking Effort Complete
```bash
before_marking_complete() {
    local effort="$1"
    
    # MUST have line count tracking
    if ! validate_line_count_tracking "$effort"; then
        echo "❌ CANNOT mark complete without line count tracking!"
        return 1
    fi
    
    # MUST be within limit
    WITHIN_LIMIT=$(jq ".efforts_in_progress[] | select(.name == \"$effort\") | .line_count_tracking.within_limit" orchestrator-state.json)
    
    if [ "$WITHIN_LIMIT" != "true" ]; then
        echo "❌ CANNOT mark complete - exceeds 800 line limit!"
        echo "🔄 Must split first!"
        return 1
    fi
}
```

### Gate 2: Before Wave Integration
```bash
before_wave_integration() {
    # ALL efforts must have line count tracking
    for effort in $(jq -r '.efforts_completed[].name' orchestrator-state.json); do
        TRACKING=$(jq ".efforts_completed[] | select(.name == \"$effort\") | .line_count_tracking" orchestrator-state.json)
        
        if [ "$TRACKING" = "null" ]; then
            echo "❌ R338 VIOLATION: Effort $effort missing line count tracking!"
            echo "🚨 CANNOT integrate wave without complete tracking!"
            return 1
        fi
    done
}
```

## 🚨 COMMON VIOLATIONS

### VIOLATION 1: Missing Line Count Tracking
```json
// ❌ WRONG - No tracking structure
{
  "name": "effort1",
  "status": "complete",
  "lines": 650  // Informal, not structured tracking!
}
```

### VIOLATION 2: Incomplete Tracking
```json
// ❌ WRONG - Missing required fields
{
  "line_count_tracking": {
    "current_count": 687
    // Missing: last_measured, measured_by, command, etc.
  }
}
```

### VIOLATION 3: Not Updating After Changes
```json
// ❌ WRONG - Stale measurement
{
  "line_count_tracking": {
    "initial_count": 500,
    "current_count": 500,  // Never updated after fixes!
    "last_measured": "2025-01-15T10:00:00Z"  // 5 days old!
  }
}
```

## ✅ CORRECT PATTERN

### Complete Tracking Example:
```json
{
  "name": "registry-auth",
  "phase": 1,
  "wave": 2,
  "line_count_tracking": {
    "initial_count": 850,
    "current_count": 687,
    "last_measured": "2025-01-20T14:30:00Z",
    "measured_by": "code-reviewer",
    "measurement_command": "./tools/line-counter.sh phase1/wave2/registry-auth",
    "auto_detected_base": "phase1-wave1-integration",
    "implementation_only": true,
    "within_limit": true,
    "requires_split": false,
    "split_plan_path": null,
    "measurement_history": [
      {
        "timestamp": "2025-01-20T10:00:00Z",
        "count": 850,
        "measured_by": "code-reviewer",
        "reason": "initial review",
        "command": "./tools/line-counter.sh phase1/wave2/registry-auth",
        "auto_detected_base": "phase1-wave1-integration",
        "excludes": "tests/demos/docs per R007",
        "within_limit": false
      },
      {
        "timestamp": "2025-01-20T14:30:00Z",
        "count": 687,
        "measured_by": "code-reviewer",
        "reason": "after refactoring",
        "command": "./tools/line-counter.sh phase1/wave2/registry-auth",
        "auto_detected_base": "phase1-wave1-integration",
        "excludes": "tests/demos/docs per R007",
        "within_limit": true
      }
    ]
  }
}
```

## 🔴 SPLIT TRACKING REQUIREMENTS

When an effort requires splitting:

```json
{
  "name": "large-effort",
  "line_count_tracking": {
    "initial_count": 1250,
    "current_count": 1250,
    "within_limit": false,
    "requires_split": true,
    "split_plan_path": "efforts/phase1/wave1/large-effort/SPLIT-PLAN.md"
  },
  "splits": [
    {
      "number": 1,
      "branch": "phase1/wave1/large-effort-split-001",
      "line_count_tracking": {
        "initial_count": null,
        "current_count": 650,
        "last_measured": "2025-01-20T16:00:00Z",
        "within_limit": true
      }
    },
    {
      "number": 2,
      "branch": "phase1/wave1/large-effort-split-002",
      "line_count_tracking": {
        "initial_count": null,
        "current_count": 600,
        "last_measured": "2025-01-20T18:00:00Z",
        "within_limit": true
      }
    }
  ]
}
```

## GRADING FORMULA

```python
def calculate_line_tracking_grade(state_file):
    total_efforts = 0
    missing_tracking = 0
    incomplete_tracking = 0
    stale_tracking = 0
    
    for effort in state_file['efforts_in_progress'] + state_file['efforts_completed']:
        total_efforts += 1
        
        if 'line_count_tracking' not in effort:
            missing_tracking += 1
        elif effort['line_count_tracking']['current_count'] is None:
            incomplete_tracking += 1
        elif is_stale(effort['line_count_tracking']['last_measured']):
            stale_tracking += 1
    
    # Base grade
    grade = 100
    
    # Deductions
    grade -= missing_tracking * 50      # -50% per missing
    grade -= incomplete_tracking * 30   # -30% per incomplete
    grade -= stale_tracking * 20        # -20% per stale
    
    # Total failure if no tracking at all
    if missing_tracking == total_efforts:
        grade = 0
    
    return max(grade, 0)
```

## Integration with Other Rules

- **R007**: Line counts exclude tests/demos/docs
- **R304**: Must use line-counter.sh tool
- **R319**: Orchestrator never measures (only records)
- **R108**: Code Reviewer must measure and report
- **R302**: Split tracking requirements

## 🔴 THE LINE COUNT LAW

```
EVERY effort gets tracked - NO EXCEPTIONS
Code Reviewer measures - Orchestrator records
State file is truth - ALWAYS current
No tracking = No progress - BLOCKED
Complete tracking = Compliant system
```

---

**REMEMBER:** The orchestrator-state.json file is the SINGLE SOURCE OF TRUTH for all line counts. Without proper tracking, the entire system fails!