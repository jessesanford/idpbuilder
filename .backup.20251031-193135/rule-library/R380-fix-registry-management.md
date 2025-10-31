# 🚨🚨🚨 RULE R380 - Fix Registry Management [BLOCKING]
**FAIL = LOST FIX TRACKING** | Source: [R380](rule-library/RULE-REGISTRY.md#R380)

## Purpose
Ensure ALL fixes are systematically tracked from discovery through validation with NO ad-hoc fixes allowed outside the registry. This prevents orchestrators from getting lost during fix cascades and guarantees finite execution.

## Requirements

### 1. EVERY Fix Gets Unique ID
```json
{
  "fix_registry": {
    "FIX-001": {
      "discovered_at": "2025-01-21T10:00:00Z",
      "discovered_during": "initial_analysis",
      "branch": "effort-1-branch",
      "issue": "Type conflict in API handler",
      "status": "pending",
      "attempts": 0,
      "max_attempts": 3
    }
  }
}
```

### 2. Progressive Discovery Support
Fixes can be discovered at ANY point during cascade execution:
- During initial analysis
- During integration attempts
- During validation phases
- During code reviews
- During build verification

When a new fix is discovered:
1. Add to registry with unique ID
2. Update convergence metrics
3. Check if within cycle limits
4. Continue cascade or escalate

### 3. Registry States
Every fix MUST be in exactly one state:
- `pending` - Fix identified but not started
- `in_progress` - Fix being applied
- `applied` - Fix applied to branch
- `validated` - Fix verified working
- `failed` - Fix failed after max attempts
- `skipped` - Fix determined unnecessary

### 4. Idempotency Enforcement
Before applying any fix:
```bash
# Check if fix already applied
git log --oneline -n 20 | grep "FIX-${FIX_ID}" && {
  echo "Fix already applied, marking validated"
  update_registry "$FIX_ID" "validated"
  continue
}
```

### 5. Convergence Tracking
```json
{
  "convergence_metrics": {
    "fixes_pending": 5,
    "fixes_completed": 12,
    "fixes_failed": 1,
    "progress_rate": 0.71,
    "cycle_count": 7,
    "max_cycles": 20,
    "last_progress_cycle": 7,
    "stalled_cycles": 0
  }
}
```

### 6. Finite Execution Guarantees
- **MAX_CYCLES**: 20 (hard limit)
- **MAX_ATTEMPTS per fix**: 3
- **MIN_PROGRESS_RATE**: 10% per 3 cycles
- **STALL_LIMIT**: 3 cycles with no progress
- **TOTAL_TIMEOUT**: 10 hours

### 7. Checkpoint Management
After EVERY fix operation:
```json
{
  "checkpoints": {
    "before_FIX_001": "abc123def",
    "after_FIX_001": "def456ghi",
    "validation_FIX_001": "passed"
  }
}
```

### 8. Registry Persistence
- Registry saved after EVERY update
- Registry loaded on cascade restart
- Registry archived on cascade completion
- Registry included in all state commits

## Enforcement

### NO Ad-Hoc Fixes
❌ **FORBIDDEN**:
```bash
# Direct fix without registry
git cherry-pick <commit>
git commit -m "Quick fix for build"
```

✅ **REQUIRED**:
```bash
# Register fix first
add_to_registry "FIX-018" "Build configuration update" "discovered_during: validation"
# Then apply with tracking
apply_fix "FIX-018"
```

### Discovery Protocol
When discovering a new fix:
```bash
discover_fix() {
  local issue="$1"
  local context="$2"

  # Generate unique ID
  local fix_id="FIX-$(printf "%03d" $(($(get_next_fix_number))))"

  # Add to registry
  add_to_registry "$fix_id" "$issue" "$context"

  # Check convergence
  if [[ $cycle_count -ge $max_cycles ]]; then
    escalate "MAX_CYCLES reached with new fix discovered"
  fi

  # Continue cascade
  return 0
}
```

### Recovery Protocol
On cascade restart:
```bash
recover_cascade() {
  # Load registry
  load_fix_registry

  # Check each fix status
  for fix_id in $(get_all_fixes); do
    status=$(get_fix_status "$fix_id")
    case $status in
      "validated") continue ;;
      "applied") validate_fix "$fix_id" ;;
      "in_progress") revert_and_retry "$fix_id" ;;
      "pending") queue_fix "$fix_id" ;;
      "failed") report_failed "$fix_id" ;;
    esac
  done

  # Resume from next pending fix
  resume_cascade
}
```

## Example Flow with Progressive Discovery

```
START CASCADE:
  Registry: [FIX-001, FIX-002] (2 fixes known)

CYCLE 1:
  Apply FIX-001 → Success
  Apply FIX-002 → Success
  Validate → DISCOVER FIX-003 (type conflict)
  Add FIX-003 to registry
  Metrics: 2/3 complete, progress=66%

CYCLE 2:
  Apply FIX-003 → Success
  Validate → DISCOVER FIX-004, FIX-005 (integration issues)
  Add to registry
  Metrics: 3/5 complete, progress=60%
  Check: cycle_count(2) < max_cycles(20) ✓

CYCLE 3:
  Apply FIX-004 → Success
  Apply FIX-005 → Conflict → Resolve → Success
  Validate → ALL PASS
  Metrics: 5/5 complete, progress=100%

COMPLETE: 5 fixes total (3 discovered during cascade)
```

## Failure Conditions

### Automatic Escalation Triggers
1. `cycle_count >= max_cycles` (20)
2. `stalled_cycles >= stall_limit` (3)
3. `total_time >= total_timeout` (10 hours)
4. `fixes_failed >= total_fixes * 0.5` (50% failure rate)

### Manual Intervention Required
When escalation triggered:
1. Save current registry state
2. Generate CASCADE-ESCALATION-REPORT.md
3. List all pending fixes
4. Document blockers
5. Transition to FIX_CASCADE_ABORT

## Integration with State Machine

### State Updates
Every state in FIX_CASCADE must:
1. Check registry before actions
2. Update registry after actions
3. Calculate convergence metrics
4. Check escalation conditions

### New States Required
- `FIX_DISCOVERY_CONTINUOUS` - Active discovery during execution
- `FIX_REGISTRY_CHECKPOINT` - Save registry state
- `FIX_CONVERGENCE_CHECK` - Verify progress toward completion
- `FIX_CASCADE_ESCALATION` - Handle stalled cascades

## Validation Checklist

Before marking cascade complete:
- [ ] All fixes in registry are `validated` or `skipped`
- [ ] No fixes in `pending` or `in_progress` state
- [ ] All checkpoints verified
- [ ] Convergence reached 100%
- [ ] No new fixes discovered in final validation

## Recovery Assistant

```bash
#!/bin/bash
# R380 Recovery Assistant

# Check for orphaned fixes (applied but not in registry)
find_orphaned_fixes() {
  git log --oneline | grep -E "FIX-[0-9]{3}" | while read -r line; do
    fix_id=$(echo "$line" | grep -oE "FIX-[0-9]{3}")
    if ! registry_contains "$fix_id"; then
      echo "WARNING: Orphaned fix $fix_id not in registry"
      add_to_registry "$fix_id" "Recovered fix" "discovered_during: recovery"
    fi
  done
}

# Validate registry consistency
validate_registry() {
  local errors=0

  # Check for duplicate IDs
  check_duplicates || ((errors++))

  # Check for invalid states
  check_valid_states || ((errors++))

  # Check for missing timestamps
  check_timestamps || ((errors++))

  # Check checkpoint consistency
  check_checkpoints || ((errors++))

  return $errors
}

# Main recovery flow
main() {
  echo "R380 Fix Registry Recovery"

  # Load existing registry
  load_fix_registry || create_new_registry

  # Find orphaned fixes
  find_orphaned_fixes

  # Validate consistency
  validate_registry || {
    echo "Registry validation failed"
    exit 1
  }

  # Resume cascade
  recover_cascade
}
```

## Rule Interactions
- **R375**: Fix state file management (registry is part of state file)
- **R376**: Quality gates (registry tracks gate status)
- **R327**: Cascade operations (registry ensures complete cascades)
- **R352**: Overlapping cascades (separate registries per cascade chain)

## Success Criteria
✅ EVERY fix has registry entry
✅ NO fixes applied without tracking
✅ Progressive discovery handled gracefully
✅ Finite execution guaranteed
✅ Automatic recovery possible
✅ Complete audit trail maintained