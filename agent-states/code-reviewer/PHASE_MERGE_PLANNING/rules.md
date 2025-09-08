# Code Reviewer - PHASE_MERGE_PLANNING State Rules

## 🚨 CRITICAL: PHASE MERGE PLAN CREATION ONLY - NO EXECUTION! 🚨

### PRIMARY PURPOSE
Create a detailed PHASE MERGE PLAN document for integrating an entire phase worth of work.
DO NOT execute any merges - only plan them!

### 🔴🔴🔴 CRITICAL RULE: USE ORIGINAL BRANCHES ONLY! 🔴🔴🔴

**MERGE FROM ORIGINAL EFFORT AND FIX BRANCHES ONLY!**
- ✅ CORRECT: phase3/wave1/effort1-api-types (original effort)
- ✅ CORRECT: phase3-fix-kcp-patterns-20250827 (ERROR_RECOVERY fix)
- ❌ WRONG: phase3-wave1-integration-20250827 (wave integration)
- ❌ WRONG: phase3-integration-20250827 (phase integration)

Wave/Phase integration branches are intermediate artifacts, not sources!

## State Context
You are creating a phase-level merge plan after ERROR_RECOVERY has created fix branches. The orchestrator has set up the phase integration infrastructure. Your role is to plan the merging of ALL wave efforts plus ERROR_RECOVERY fixes.

## PHASE MERGE PLAN REQUIREMENTS

### 1. Inventory All Branches to Merge
```bash
#!/bin/bash
# Find all branches for phase integration

PHASE=$(jq '.current_phase' orchestrator-state.json)
echo "📊 Analyzing branches for Phase $PHASE integration..."

# A. Original effort branches from all waves
echo "=== EFFORT BRANCHES ==="
for wave in $(seq 1 4); do  # Adjust based on actual wave count
    echo "Wave $wave efforts:"
    for effort_dir in /efforts/phase${PHASE}/wave${wave}/*/; do
        [[ "$effort_dir" == *"integration-workspace"* ]] && continue
        
        effort=$(basename "$effort_dir")
        cd "$effort_dir"
        
        current_branch=$(git branch --show-current)
        
        # Include splits, exclude "too large" originals
        if [[ "$current_branch" == *"-split"* ]]; then
            echo "  ✅ $current_branch (split)"
        elif git log --oneline | grep -q "too large"; then
            echo "  ❌ $current_branch (too large, use splits)"
        else
            echo "  ✅ $current_branch"
        fi
    done
done

# B. ERROR_RECOVERY fix branches
echo "=== FIX BRANCHES ==="
FIX_BRANCHES=$(jq '.error_recovery_fixes[].branch' orchestrator-state.json)
if [ -n "$FIX_BRANCHES" ]; then
    for branch in $FIX_BRANCHES; do
        echo "  ✅ $branch (ERROR_RECOVERY fix)"
    done
else
    # Check by pattern
    git branch -r | grep "origin/phase${PHASE}-fix-" | while read branch; do
        echo "  ✅ ${branch#origin/} (ERROR_RECOVERY fix)"
    done
fi
```

### 2. Determine Phase-Level Merge Order
```markdown
## Phase Merge Order Analysis

### Merge Groups (IN ORDER):

#### Group 1: Wave 1 Efforts
All Wave 1 effort branches (no dependencies on other waves)
- phase3/wave1/effort1-api-types
- phase3/wave1/effort2-controller-split1
- phase3/wave1/effort2-controller-split2

#### Group 2: Wave 2 Efforts  
All Wave 2 effort branches (may depend on Wave 1)
- phase3/wave2/effort1-webhooks
- phase3/wave2/effort2-validation

#### Group 3: Wave 3 Efforts
All Wave 3 effort branches (may depend on Waves 1-2)
- phase3/wave3/effort1-reconciler
- phase3/wave3/effort2-status-split1
- phase3/wave3/effort2-status-split2

#### Group 4: Wave 4 Efforts
All Wave 4 effort branches (may depend on Waves 1-3)
- phase3/wave4/effort1-metrics
- phase3/wave4/effort2-observability

#### Group 5: ERROR_RECOVERY Fixes
All fix branches addressing phase assessment issues
- phase3-fix-kcp-patterns-20250827
- phase3-fix-api-compatibility-20250827
- phase3-fix-test-coverage-20250827

### Dependency Analysis:
- Waves are merged in sequence (1→2→3→4)
- Within a wave, splits must maintain order (split1→split2→split3)
- ERROR_RECOVERY fixes are applied LAST to address assessment issues
```

### 3. Reference Phase Assessment Report
```markdown
## Phase Assessment Issues Being Addressed

From: phase-assessments/phase${PHASE}/PHASE-${PHASE}-ASSESSMENT-REPORT.md

### Priority 1 Issues (Must Fix):
1. KCP pattern violations in controller logic
   - Fix Branch: phase3-fix-kcp-patterns-20250827
2. API backward compatibility issues
   - Fix Branch: phase3-fix-api-compatibility-20250827
3. Insufficient test coverage (<80%)
   - Fix Branch: phase3-fix-test-coverage-20250827

### Verification:
Each fix branch addresses specific assessment findings and will be merged last.
```

### 4. Create PHASE-MERGE-PLAN.md

**Location:** `${INTEGRATION_DIR}/PHASE-MERGE-PLAN.md`

```markdown
# Phase ${PHASE} Integration Merge Plan

**Generated:** $(date -u +%Y-%m-%dT%H:%M:%SZ)
**Code Reviewer:** code-reviewer
**State:** PHASE_MERGE_PLANNING
**Type:** Post-Assessment Integration

## Target Integration Branch
- **Branch Name:** phase${PHASE}-post-fixes-integration-${timestamp}
- **Base:** main at ${commit}
- **Purpose:** Complete phase integration with assessment fixes
- **Location:** /efforts/phase${PHASE}/phase-integration-workspace

## Phase Scope Summary
- **Total Waves:** 4
- **Total Efforts:** 12 (3 splits replaced originals)
- **ERROR_RECOVERY Fixes:** 3
- **Assessment Score Before Fixes:** 68/100
- **Target Score After Integration:** >85/100

## Branches to Merge (IN STRICT ORDER)

### Wave 1 Efforts (Foundation)

#### 1. phase3/wave1/effort1-api-types
- **Type:** Original effort branch
- **Base:** main at abc123
- **Size:** 542 lines
- **Purpose:** Core API type definitions
- **Merge Command:**
  ```bash
  git fetch origin phase3/wave1/effort1-api-types
  git merge origin/phase3/wave1/effort1-api-types --no-ff \
    -m "Phase integration: Wave 1 effort1-api-types"
  ```

#### 2. phase3/wave1/effort2-controller-split1
- **Type:** Split branch (1 of 2)
- **Base:** main at abc123
- **Size:** 398 lines
- **Purpose:** Controller base implementation
- **Merge Command:**
  ```bash
  git fetch origin phase3/wave1/effort2-controller-split1
  git merge origin/phase3/wave1/effort2-controller-split1 --no-ff \
    -m "Phase integration: Wave 1 effort2-controller-split1"
  ```

#### 3. phase3/wave1/effort2-controller-split2
- **Type:** Split branch (2 of 2)
- **Base:** effort2-controller-split1 at def456
- **Size:** 412 lines
- **Purpose:** Controller reconciliation logic
- **Merge Command:**
  ```bash
  git fetch origin phase3/wave1/effort2-controller-split2
  git merge origin/phase3/wave1/effort2-controller-split2 --no-ff \
    -m "Phase integration: Wave 1 effort2-controller-split2"
  ```

### Wave 2 Efforts (Webhooks & Validation)

#### 4. phase3/wave2/effort1-webhooks
[Continue pattern for all waves...]

### Wave 3 Efforts (Reconciliation)
[Continue...]

### Wave 4 Efforts (Observability)
[Continue...]

### ERROR_RECOVERY Fix Branches (Assessment Remediation)

#### 15. phase3-fix-kcp-patterns-20250827
- **Type:** ERROR_RECOVERY fix branch
- **Base:** phase3/wave4 completed state
- **Purpose:** Fix KCP pattern violations (Priority 1)
- **Assessment Issue:** #1 - Controller patterns non-compliant
- **Merge Command:**
  ```bash
  git fetch origin phase3-fix-kcp-patterns-20250827
  git merge origin/phase3-fix-kcp-patterns-20250827 --no-ff \
    -m "Phase integration: Fix KCP pattern violations per assessment"
  ```

#### 16. phase3-fix-api-compatibility-20250827
- **Type:** ERROR_RECOVERY fix branch
- **Base:** Previous fix applied
- **Purpose:** Restore API backward compatibility (Priority 1)
- **Assessment Issue:** #2 - Breaking API changes detected
- **Merge Command:**
  ```bash
  git fetch origin phase3-fix-api-compatibility-20250827
  git merge origin/phase3-fix-api-compatibility-20250827 --no-ff \
    -m "Phase integration: Fix API compatibility per assessment"
  ```

#### 17. phase3-fix-test-coverage-20250827
- **Type:** ERROR_RECOVERY fix branch
- **Base:** Previous fixes applied
- **Purpose:** Increase test coverage to >80% (Priority 1)
- **Assessment Issue:** #3 - Test coverage at 67%
- **Merge Command:**
  ```bash
  git fetch origin phase3-fix-test-coverage-20250827
  git merge origin/phase3-fix-test-coverage-20250827 --no-ff \
    -m "Phase integration: Fix test coverage per assessment"
  ```

## Excluded Branches (DO NOT MERGE)
These branches are superseded or intermediate:
- phase3-wave1-integration-* (intermediate wave integration)
- phase3-wave2-integration-* (intermediate wave integration)
- phase3-wave3-integration-* (intermediate wave integration)
- phase3-wave4-integration-* (intermediate wave integration)
- phase3/wave1/effort2-controller (original, too large, use splits)

## Merge Strategy
1. **Wave-by-Wave:** Complete each wave before moving to next
2. **Splits in Sequence:** Maintain split order within efforts
3. **Fixes Last:** Apply ERROR_RECOVERY fixes after all efforts
4. **Conflict Resolution:** Favor fix branches over original implementation
5. **Testing:** Phase-level tests after all merges

## Expected Integration Challenges
1. **Controller Conflicts:** Waves 1 & 3 both modify controller logic
2. **API Evolution:** Each wave extends API, fixes ensure compatibility
3. **Test Interactions:** New tests from fixes may reveal issues
4. **Large Integration:** ~15,000 lines total, significant complexity

## Phase-Level Validation
```bash
# After ALL merges complete:

# 1. Run phase-specific test suite
make test-phase3

# 2. Verify assessment issues fixed
./scripts/verify-assessment-fixes.sh phase3

# 3. Check combined size
$PROJECT_ROOT/tools/line-counter.sh -c $(git branch --show-current)

# 4. Performance benchmarks
make benchmark-phase3

# 5. Integration smoke tests
make test-integration
```

## Risk Mitigation
- **High Risk:** Large number of branches increases conflict probability
- **Mitigation 1:** Test after each wave group
- **Mitigation 2:** Save state between wave groups
- **Mitigation 3:** Document all conflicts in work-log.md

## Integration Agent Instructions
1. This is a PHASE-level integration (larger scope than wave)
2. Execute merges in EXACT order specified
3. Test after each WAVE GROUP (not just each merge)
4. Document conflicts with detail about resolution
5. Verify assessment fixes are effective
6. Create comprehensive PHASE-INTEGRATION-REPORT.md

## Success Criteria
- All 17 branches merged successfully
- No test failures after integration
- Assessment issues verified as resolved
- Performance benchmarks pass
- Ready for architect reassessment
```

## Validation Before Completion

```bash
#!/bin/bash
# Validate phase merge plan completeness

validate_phase_merge_plan() {
    local plan_file="$1"
    local phase="$2"
    
    echo "🔍 Validating Phase $phase Merge Plan..."
    
    # Check wave coverage
    for wave in $(seq 1 4); do
        if ! grep -q "Wave $wave Efforts" "$plan_file"; then
            echo "❌ Missing Wave $wave in merge plan"
            return 1
        fi
    done
    
    # Check ERROR_RECOVERY fixes included
    if ! grep -q "ERROR_RECOVERY Fix Branches" "$plan_file"; then
        echo "❌ Missing ERROR_RECOVERY fixes section"
        return 1
    fi
    
    # Verify no integration branches as sources
    if grep -E "phase[0-9]-wave[0-9]-integration" "$plan_file" | grep "git merge"; then
        echo "❌ ERROR: Wave integration branches used as sources!"
        return 1
    fi
    
    # Check assessment reference
    if ! grep -q "PHASE-${phase}-ASSESSMENT-REPORT.md" "$plan_file"; then
        echo "⚠️ Warning: No reference to phase assessment report"
    fi
    
    # Count total merges planned
    merge_count=$(grep -c "git merge origin/" "$plan_file")
    echo "📊 Total branches to merge: $merge_count"
    
    if [[ $merge_count -lt 10 ]]; then
        echo "⚠️ Warning: Fewer merges than expected for a full phase"
    fi
    
    echo "✅ Phase merge plan validation passed"
    return 0
}

# Run validation
validate_phase_merge_plan "PHASE-MERGE-PLAN.md" "$PHASE"
```

## State Transitions

From PHASE_MERGE_PLANNING state:
- **PLAN_COMPLETE** → Return to orchestrator
- **VALIDATION_FAILED** → Fix and re-validate

## Critical Success Criteria

1. ✅ PHASE-MERGE-PLAN.md created in phase integration directory
2. ✅ All wave efforts catalogued (excluding "too large" originals)
3. ✅ All ERROR_RECOVERY fix branches included
4. ✅ Wave-by-wave merge order established
5. ✅ NO wave integration branches used as sources
6. ✅ Assessment report issues mapped to fix branches
7. ✅ Phase-level validation steps included
8. ✅ Clear instructions for Integration Agent

## Common Mistakes to Avoid

1. **Merging from wave integration branches**
   - ❌ WRONG: Use phase3-wave1-integration as source
   - ✅ RIGHT: Use original effort branches only

2. **Missing ERROR_RECOVERY fixes**
   - ❌ WRONG: Only merge effort branches
   - ✅ RIGHT: Include all fix branches at the end

3. **Wrong wave order**
   - ❌ WRONG: Random wave ordering
   - ✅ RIGHT: Sequential wave progression

4. **Executing merges**
   - ❌ WRONG: Running git merge commands
   - ✅ RIGHT: Only documenting merge plan

5. **Ignoring assessment report**
   - ❌ WRONG: Generic integration plan
   - ✅ RIGHT: Reference specific assessment issues

---
### ⚠️⚠️⚠️ RULE R261 - Code Reviewer Merge Plan No Execution
**Source:** rule-library/R261-code-reviewer-merge-plan-no-execution.md
**Criticality:** WARNING - Violation = Role confusion

Code Reviewer creates merge plans ONLY. NEVER executes merges. That's the Integration Agent's job.
---

---
### 🔴🔴🔴 RULE R262 - No Integration Branches as Sources
**Source:** rule-library/R262-no-integration-branches-as-sources.md
**Criticality:** SUPREME - Violation = Recursive integration chaos

CRITICAL: Only original effort/fix branches in merge plans. Integration branches are TARGETS not SOURCES.
---