# 🚨🚨🚨 MANDATORY LINE COUNT TRACKING ENFORCEMENT PROMPT 🚨🚨🚨

## 🔴🔴🔴 IMMEDIATE ACTION REQUIRED: IMPLEMENT R338 LINE COUNT TRACKING 🔴🔴🔴

**YOU ARE THE ORCHESTRATOR AGENT AND YOU MUST IMMEDIATELY IMPLEMENT COMPREHENSIVE LINE COUNT TRACKING!**

## 📊 CRITICAL FAILURES DETECTED IN SYSTEM

### ❌ CURRENT STATE: CRITICAL COMPLIANCE VIOLATIONS
- **ZERO** efforts have line_count_tracking structures
- **ZERO** line counts captured from Code Reviewer reports
- **ZERO** compliance with R338 requirements
- **100% FAILURE RATE** on size tracking

### 🚨 GRADING IMPACT: -100% IMMEDIATE FAILURE
Without line count tracking, you are operating BLIND and will receive:
- **-50%** for EACH effort missing line_count_tracking
- **-100%** for NO tracking at all (current state)
- **AUTOMATIC FAILURE** for size violations you didn't track

## 🔴 MANDATORY IMPLEMENTATION TASKS

### TASK 1: AUDIT ALL EXISTING EFFORTS
```bash
# Check orchestrator-state.json for all efforts
cat orchestrator-state.json | jq '.efforts_in_progress, .efforts_completed'

# For EACH effort, check if line_count_tracking exists:
cat orchestrator-state.json | jq '.efforts_in_progress[].line_count_tracking'
```

### TASK 2: ADD LINE_COUNT_TRACKING TO EVERY EFFORT
For EVERY effort in orchestrator-state.json, add this structure:

```json
"line_count_tracking": {
  "initial_count": null,
  "current_count": null,
  "last_measured": null,
  "measured_by": null,
  "measurement_command": null,
  "auto_detected_base": null,
  "implementation_only": true,
  "within_limit": null,
  "requires_split": false,
  "split_plan_path": null,
  "measurement_history": []
}
```

### TASK 3: RETRIEVE LINE COUNTS FROM EXISTING REPORTS
For each effort that has a CODE-REVIEW-REPORT:

```bash
# Find all code review reports
find efforts/ -name "CODE-REVIEW-REPORT*.md" -type f

# Extract line counts from each report
for report in $(find efforts/ -name "CODE-REVIEW-REPORT*.md"); do
    echo "Processing: $report"
    
    # Look for the new standardized format
    grep "Implementation Lines:" "$report" || \
    grep "Current Lines:" "$report" || \
    grep "Total.*lines:" "$report"
    
    # Extract the command used
    grep "Command:" "$report" || \
    grep "Tool Used:" "$report"
    
    # Extract the base branch
    grep "Auto-detected Base:" "$report" || \
    grep "Detected base:" "$report"
done
```

### TASK 4: UPDATE STATE FILE WITH LINE COUNTS
For each effort with a review report, update orchestrator-state.json:

```bash
# Example update for an effort
EFFORT_NAME="registry-auth"
LINE_COUNT=687
COMMAND="./tools/line-counter.sh phase1/wave1/registry-auth"
BASE="phase1-wave1-integration"
TIMESTAMP=$(date -u +%Y-%m-%dT%H:%M:%SZ)

# Update the state file
jq --arg effort "$EFFORT_NAME" \
   --argjson count "$LINE_COUNT" \
   --arg cmd "$COMMAND" \
   --arg base "$BASE" \
   --arg ts "$TIMESTAMP" \
   '.efforts_completed = [.efforts_completed[] | 
    if .name == $effort then
      .line_count_tracking = {
        initial_count: $count,
        current_count: $count,
        last_measured: $ts,
        measured_by: "code-reviewer",
        measurement_command: $cmd,
        auto_detected_base: $base,
        implementation_only: true,
        within_limit: ($count <= 800),
        requires_split: ($count > 800),
        split_plan_path: (if $count > 800 then "efforts/\(.phase)/\(.wave)/\(.name)/SPLIT-PLAN.md" else null end),
        measurement_history: [{
          timestamp: $ts,
          count: $count,
          measured_by: "code-reviewer",
          reason: "initial review",
          command: $cmd,
          auto_detected_base: $base,
          excludes: "tests/demos/docs per R007",
          within_limit: ($count <= 800)
        }]
      }
    else . end]' orchestrator-state.json > tmp.json && mv tmp.json orchestrator-state.json
```

### TASK 5: VALIDATE ALL TRACKING IS COMPLETE
```bash
# Count efforts without tracking
TOTAL_EFFORTS=$(jq '[.efforts_in_progress[], .efforts_completed[]] | length' orchestrator-state.json)
WITH_TRACKING=$(jq '[.efforts_in_progress[], .efforts_completed[]] | map(select(.line_count_tracking != null)) | length' orchestrator-state.json)

echo "Total efforts: $TOTAL_EFFORTS"
echo "With tracking: $WITH_TRACKING"
echo "Missing tracking: $((TOTAL_EFFORTS - WITH_TRACKING))"

# List efforts missing tracking
jq -r '[.efforts_in_progress[], .efforts_completed[]] | 
       map(select(.line_count_tracking == null)) | 
       .[].name' orchestrator-state.json
```

### TASK 6: CREATE TRACKING COMPLIANCE REPORT
```markdown
# LINE COUNT TRACKING COMPLIANCE REPORT

## Summary
- Total Efforts: [X]
- With Tracking: [Y]
- Missing Tracking: [Z]
- Compliance Rate: [Y/X * 100]%

## Efforts WITH Line Count Tracking
1. [effort-name]: [line_count] lines (last measured: [date])
2. ...

## Efforts MISSING Line Count Tracking
1. [effort-name] - ACTION REQUIRED: Need Code Reviewer measurement
2. ...

## Size Violations Detected
1. [effort-name]: [line_count] lines (EXCEEDS 800 limit)
   - Split Plan: [path-to-split-plan]
   - Status: [split status]

## Next Steps
1. Spawn Code Reviewers for unmeasured efforts
2. Update tracking for efforts with stale measurements
3. Create split plans for violations
```

### TASK 7: IMPLEMENT AUTOMATED CAPTURE GOING FORWARD
Add this function to your monitoring workflow:

```bash
capture_line_count_from_review() {
    local report_path="$1"
    local effort_name="$2"
    
    # Extract standardized format per R338
    LINE_COUNT=$(grep "Implementation Lines:" "$report_path" | awk '{print $3}')
    COMMAND=$(grep "Command:" "$report_path" | cut -d':' -f2- | xargs)
    BASE=$(grep "Auto-detected Base:" "$report_path" | cut -d':' -f2- | xargs)
    TIMESTAMP=$(date -u +%Y-%m-%dT%H:%M:%SZ)
    
    if [ -z "$LINE_COUNT" ]; then
        echo "⚠️ WARNING: No line count found in report!"
        return 1
    fi
    
    echo "✅ Captured line count for $effort_name: $LINE_COUNT lines"
    
    # Update orchestrator-state.json immediately
    # [Insert jq update command here]
    
    # Commit the update
    git add orchestrator-state.json
    git commit -m "tracking: update line count for $effort_name - $LINE_COUNT lines"
    git push
}

# Use in MONITOR_REVIEWS state
for effort in $(jq -r '.efforts_in_progress[].name' orchestrator-state.json); do
    REPORT="efforts/phase${PHASE}/wave${WAVE}/${effort}/CODE-REVIEW-REPORT.md"
    if [ -f "$REPORT" ]; then
        capture_line_count_from_review "$REPORT" "$effort"
    fi
done
```

## 🔴 VALIDATION CHECKPOINTS

### Before ANY State Transition:
```bash
validate_line_count_tracking() {
    # Check for missing tracking
    MISSING=$(jq '[.efforts_in_progress[], .efforts_completed[]] | 
                  map(select(.line_count_tracking == null)) | length' orchestrator-state.json)
    
    if [ "$MISSING" -gt 0 ]; then
        echo "🚨 R338 VIOLATION: $MISSING efforts missing line count tracking!"
        echo "❌ CANNOT proceed to next state!"
        return 1
    fi
    
    echo "✅ All efforts have line count tracking"
    return 0
}
```

### Before Wave Completion:
```bash
# ALL efforts in wave must have current line counts
for effort in $(jq -r '.efforts_completed[] | select(.wave == '$WAVE') | .name' orchestrator-state.json); do
    TRACKING=$(jq ".efforts_completed[] | select(.name == \"$effort\") | .line_count_tracking" orchestrator-state.json)
    if [ "$TRACKING" = "null" ]; then
        echo "❌ Cannot complete wave - $effort missing line count tracking!"
        exit 1
    fi
done
```

## 🚨 ENFORCEMENT MESSAGE

**THIS IS NOT OPTIONAL!**

Per R338, you MUST:
1. **IMMEDIATELY** add line_count_tracking to ALL efforts
2. **CAPTURE** line counts from ALL Code Reviewer reports
3. **UPDATE** tracking after EVERY measurement
4. **VALIDATE** tracking before EVERY state transition
5. **REPORT** compliance status regularly

**Failure to implement = -100% AUTOMATIC FAILURE**

## 📋 SUCCESS CRITERIA

You have successfully implemented R338 when:
- ✅ EVERY effort has line_count_tracking structure
- ✅ ALL line counts are captured and current
- ✅ State file is the single source of truth for sizes
- ✅ Automated capture is working for new reviews
- ✅ Validation gates prevent missing tracking
- ✅ Compliance report shows 100% tracking

## 🔴 START IMMEDIATELY

**DO NOT PROCEED WITH ANY OTHER WORK UNTIL LINE COUNT TRACKING IS FULLY IMPLEMENTED!**

The system is watching. Your compliance will be measured. Your grade depends on this.

---

**R338 Enforcement Prompt v1.0**
**Criticality: BLOCKING**
**Penalty for Non-Compliance: -100%**