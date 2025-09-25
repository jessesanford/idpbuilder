# 🚨🚨🚨 RULE R339 - Fix Grace Period Protocol

**Criticality:** BLOCKING - Prevents unnecessary splits during re-integration fixes  
**Grading Impact:** -30% for incorrect split triggers during fixes  
**Enforcement:** AUTOMATIC - Applied during FIX_ISSUES state  
**Scope:** Code Reviewer (PRIMARY), SW Engineer, Orchestrator

## 🔴🔴🔴 THE FIX GRACE PERIOD LAW 🔴🔴🔴

**When fixes are backported to an existing effort branch, a 100-line grace period applies to prevent disruptive splits during re-integration cascades!**

## Rule Statement

During FIX_ISSUES state, when fixes are applied to an existing effort branch that has already been reviewed and integrated:
- **Original effort + fix lines < 900 = NO SPLIT NEEDED** (100 line grace period)
- **Original effort + fix lines >= 900 = SPLIT REQUIRED**

This prevents the anti-pattern of splitting efforts during re-integration fixes, which disrupts the cascade and creates unnecessary complexity.

## Grace Period Application Rules

### 1. WHEN Grace Period Applies (ALL conditions must be true):
```yaml
grace_period_eligible:
  - effort_state: "completed" or "integrated"  # Already reviewed/integrated
  - current_state: "FIX_ISSUES"               # Applying fixes from integration
  - fix_source: "integration_failure"         # Fix discovered during integration
  - original_count: <= 800                    # Was originally compliant
```

### 2. WHEN Grace Period DOES NOT Apply:
```yaml
no_grace_period:
  - initial_implementation: true              # First-time development
  - feature_additions: true                   # Adding new features
  - refactoring: true                         # Major code changes
  - second_fix_round: true                    # Already used grace period once
```

### 3. Grace Period Calculation:
```python
def check_grace_period_split_requirement(effort):
    if not effort.is_fix_phase:
        # Regular 800 line limit for initial implementation
        return effort.current_count > 800
    
    if effort.grace_period_used:
        # Already used grace period, back to 800 limit
        return effort.current_count > 800
    
    # First fix gets 100 line grace period
    if effort.original_count + effort.fix_delta < 900:
        effort.grace_period_used = True
        return False  # No split needed
    else:
        return True   # Split required even with grace
```

## Required Tracking Structure

### Enhanced Line Count Tracking (R338 Extension):
```json
"line_count_tracking": {
  "initial_count": 750,              // Original implementation size
  "current_count": 820,              // After fixes applied
  "fix_delta": 70,                   // Lines added by fixes
  "grace_period_eligible": true,     // Meets grace period criteria
  "grace_period_applied": true,      // Grace period was used
  "grace_period_threshold": 900,     // Threshold with grace
  "requires_split": false,            // false because 820 < 900
  "fix_history": [
    {
      "fix_id": "FIX-001",
      "source": "wave-integration-failure",
      "lines_added": 70,
      "total_after_fix": 820,
      "within_grace": true,
      "timestamp": "2025-01-20T14:30:00Z"
    }
  ]
}
```

## Code Reviewer Responsibilities

### During Fix Review:
```bash
review_fix_with_grace_period() {
    local effort="$1"
    local original_count="$2"
    local fix_lines="$3"
    local total_count=$((original_count + fix_lines))
    
    echo "📊 FIX GRACE PERIOD ANALYSIS:"
    echo "Original implementation: ${original_count} lines"
    echo "Fix adds: ${fix_lines} lines"
    echo "Total after fix: ${total_count} lines"
    echo ""
    
    if [ "$total_count" -lt 900 ]; then
        echo "✅ WITHIN GRACE PERIOD (${total_count} < 900)"
        echo "Decision: NO SPLIT REQUIRED"
        echo "Rationale: R339 grace period prevents disruption during re-integration"
    else
        echo "🚨 EXCEEDS GRACE PERIOD (${total_count} >= 900)"
        echo "Decision: SPLIT REQUIRED"
        echo "Next Step: Create split plan for fix"
    fi
}
```

### Review Report Format:
```markdown
## SIZE COMPLIANCE WITH GRACE PERIOD

**Original Effort Size:** 750 lines  
**Fix Adds:** 70 lines  
**Total After Fix:** 820 lines  
**Grace Period Applied:** YES  
**Grace Period Threshold:** 900 lines  
**Split Required:** NO

### R339 Grace Period Justification:
- ✅ Effort was already integrated
- ✅ Fix discovered during wave integration
- ✅ Total (820) < Grace threshold (900)
- ✅ Prevents cascade disruption
```

## SW Engineer Responsibilities

### In FIX_ISSUES State:
```bash
apply_fix_with_grace_awareness() {
    echo "📋 Applying fix per R321 backport requirements"
    echo "📊 Checking R339 grace period eligibility..."
    
    # Get current size
    CURRENT_SIZE=$($CLAUDE_PROJECT_DIR/tools/line-counter.sh -c "$BRANCH")
    
    # Estimate fix size
    echo "Estimated fix size: ~70 lines"
    PROJECTED_TOTAL=$((CURRENT_SIZE + 70))
    
    if [ "$PROJECTED_TOTAL" -lt 900 ]; then
        echo "✅ Within grace period, proceeding with fix"
        # Apply fix normally
    else
        echo "⚠️ Would exceed grace period ($PROJECTED_TOTAL >= 900)"
        echo "📋 Planning minimal fix to stay under 900"
        # Adjust fix approach if possible
    fi
}
```

## Orchestrator Responsibilities

### Monitor Grace Period Usage:
```bash
track_grace_period_usage() {
    local effort="$1"
    
    # Update orchestrator-state.json
    jq ".efforts_completed[\"$effort\"].line_count_tracking |= . + {
        grace_period_applied: true,
        grace_period_threshold: 900,
        grace_period_reason: \"R321 backport during integration\"
    }" orchestrator-state.json > tmp.json && mv tmp.json orchestrator-state.json
    
    echo "📊 Grace period tracked for $effort"
}
```

### State Transition Validation:
```bash
# Before transitioning from FIX_ISSUES to next state
validate_fix_compliance() {
    local effort="$1"
    local total_lines="$2"
    
    if [ "$total_lines" -ge 900 ]; then
        echo "🚨 Fix exceeds grace period - split required"
        transition_to_state "CREATE_SPLIT_PLAN"
    else
        echo "✅ Fix within grace period - proceeding"
        transition_to_state "WAVE_INTEGRATION_RETRY"
    fi
}
```

## Practical Examples

### Example 1: Fix Within Grace Period
```
Scenario: E1.2.3 completed at 750 lines
Integration finds bug needing 70-line fix

Analysis:
- Original: 750 lines ✅
- Fix adds: 70 lines
- Total: 820 lines
- Grace threshold: 900 lines
- Decision: NO SPLIT (820 < 900) ✅

Result: Fix applied, no split, cascade continues smoothly
```

### Example 2: Fix Exceeds Grace Period
```
Scenario: E1.2.4 completed at 780 lines
Integration needs 150-line fix

Analysis:
- Original: 780 lines ✅
- Fix adds: 150 lines
- Total: 930 lines
- Grace threshold: 900 lines
- Decision: SPLIT REQUIRED (930 >= 900) 🚨

Result: Must split the fix into smaller parts
```

### Example 3: Second Fix (No Grace)
```
Scenario: E1.2.3 already used grace (now 820 lines)
Another integration issue needs 30-line fix

Analysis:
- Current: 820 lines (grace already used)
- Fix adds: 30 lines
- Total: 850 lines
- Threshold: 800 lines (no grace on second fix)
- Decision: SPLIT REQUIRED (850 > 800) 🚨

Result: Grace period only applies once per effort
```

## Why This Prevents Cascade Disruption

### WITHOUT Grace Period (Disruptive):
```
1. Effort at 750 lines → Integrated to wave
2. Integration finds 70-line fix needed
3. Total would be 820 lines
4. Forced to split into effort + effort-split1
5. Wave integration now invalid (missing split1)
6. Must recreate wave integration
7. Phase integration now invalid
8. Entire cascade disrupted for 20 lines over limit
```

### WITH Grace Period (Smooth):
```
1. Effort at 750 lines → Integrated to wave
2. Integration finds 70-line fix needed
3. Total would be 820 lines
4. Grace period allows up to 900 for fixes
5. Fix applied directly to effort branch
6. Wave integration recreated with fixed effort
7. Cascade continues smoothly
8. No unnecessary split disruption
```

## Relationship to Other Rules

### Complements R321:
- R321: Mandates immediate backport of fixes
- R339: Prevents unnecessary splits during backport

### Enhances R327:
- R327: Requires cascade re-integration after fixes
- R339: Makes cascade smoother by avoiding split disruption

### Modifies R007:
- R007: Strict 800 line limit
- R339: Adds 100 line grace for fixes only

### Extends R338:
- R338: Tracks line counts
- R339: Adds grace period tracking fields

## Quick Reference

| Situation | Original | Fix | Total | Threshold | Split? |
|-----------|----------|-----|-------|-----------|--------|
| Initial implementation | N/A | N/A | 850 | 800 | YES |
| First fix | 750 | 70 | 820 | 900 | NO |
| First fix | 750 | 200 | 950 | 900 | YES |
| Second fix | 820 | 30 | 850 | 800 | YES |
| Feature addition | 700 | 150 | 850 | 800 | YES |

## Remember

**"Grace for fixes, not for features"**  
**"100 lines to smooth the cascade"**  
**"Prevent disruption, maintain momentum"**  
**"One grace per effort, use it wisely"**

The grace period exists to prevent the re-integration cascade from being disrupted by minor size overages during critical bug fixes. It is NOT a general relaxation of size limits.