# 🔴🔴🔴 RULE R353: CASCADE FOCUS PROTOCOL 🔴🔴🔴

## SUPREME LAW - NO DIVERSIONS DURING CASCADE OPERATIONS

**CRITICALITY:** SUPREME LAW - Violation = -100% AUTOMATIC FAILURE
**PRIORITY:** P0 - HIGHEST
**ENFORCEMENT:** MANDATORY - NO EXCEPTIONS

## 🚨🚨🚨 THE ABSOLUTE CASCADE FOCUS LAW 🚨🚨🚨

**DURING CASCADE_REINTEGRATION, NO DIVERSIONS ARE ALLOWED!**

When the system is in CASCADE mode (propagating fixes through integration levels), ALL agents MUST maintain absolute focus on the cascade operation. NO size checks, NO split evaluations, NO quality assessments that could divert from the cascade.

## 🔴 CORE PRINCIPLE: CASCADE FOCUS IS SACRED

### What CASCADE Mode Means:
- System is propagating critical fixes through integration levels
- Every moment of delay impacts the entire project
- Diversions break the cascade chain and corrupt state
- Focus must be maintained until fixes reach project level

### What MUST Happen in CASCADE Mode:
1. ✅ ONLY validate that rebases were successful
2. ✅ ONLY check for merge conflicts
3. ✅ ONLY verify builds still pass
4. ✅ ONLY identify NEW fixes that need cascading
5. ✅ **R354 ENFORCEMENT**: Spawn Code Reviewers for post-rebase reviews
6. ✅ Continue cascade without ANY diversions

### What MUST NOT Happen in CASCADE Mode:
1. ❌ NO line count measurements for split decisions
2. ❌ NO split requirement evaluations
3. ❌ NO quality assessments beyond "does it build"
4. ❌ NO transitions to split-related states
5. ❌ NO spawning of split planning operations
6. ❌ NO "let's also fix this while we're here"

## 🚨 CASCADE MODE CONTEXT PROPAGATION

### When Orchestrator Enters CASCADE_REINTEGRATION:
```json
{
  "cascade_coordination": {
    "cascade_mode": true,
    "cascade_purpose": "fix_propagation",
    "cascade_focus": {
      "allow_splits": false,
      "allow_size_checks": false,
      "only_validate_rebases": true,
      "deferred_evaluations": [
        "size_checks",
        "split_requirements",
        "code_quality",
        "performance_analysis"
      ]
    }
  }
}
```

### When Spawning ANY Agent During CASCADE:
```bash
# Orchestrator MUST pass cascade context to ALL spawned agents
if [[ "$(jq -r '.cascade_coordination.cascade_mode' orchestrator-state.json)" == "true" ]]; then
    CASCADE_CONTEXT="--cascade-mode=true --skip-size-checks --skip-split-evaluation --rebase-validation-only"
    echo "🔴 CASCADE MODE ACTIVE - Passing focus context to agent"
fi

# Spawn with CASCADE context
/spawn-agent code-reviewer $CASCADE_CONTEXT --effort="$EFFORT"
```

## 🔴🔴🔴 R354 EXCEPTION: POST-REBASE REVIEWS ARE MANDATORY 🔴🔴🔴

**CRITICAL CLARIFICATION: R354 SUPERSEDES SOME R353 RESTRICTIONS!**

While R353 says "no diversions" during cascade, R354 creates ONE mandatory exception:
- **AFTER EVERY REBASE**: Must spawn Code Reviewer for post-rebase review
- **This is NOT a diversion**: It's a mandatory cascade safety check
- **Review scope remains limited**: Only integration validation per R353
- **No quality/size checks**: Still skip those per R353

### The R353/R354 Interaction:
```
R353: "No diversions during cascade!"
R354: "Except post-rebase reviews - those are MANDATORY"
Result: Post-rebase reviews happen IN cascade mode with R353 restrictions
```

## 🔴 CODE REVIEWER CASCADE BEHAVIOR

### When Code Reviewer Receives CASCADE_MODE=true:

**MUST DO:**
```bash
# 1. Acknowledge CASCADE mode
echo "🔴 CASCADE MODE ACTIVE - Maintaining cascade focus per R353"

# 2. SKIP all size measurements
echo "📊 Size checks DEFERRED during cascade (R353)"

# 3. ONLY validate rebase success
cd $EFFORT_DIR && git log --oneline -5
echo "✅ Rebase successful - no conflicts"

# 4. Check for build issues ONLY
make test || echo "⚠️ Build issues detected - needs fix"

# 5. Return CASCADE-focused verdict
echo "CASCADE_VALIDATION: REBASE_VALID | REBASE_FAILED | NEEDS_FIX"
```

**MUST NOT DO:**
```bash
# ❌ NEVER do this during CASCADE:
$PROJECT_ROOT/tools/line-counter.sh  # FORBIDDEN during cascade!

# ❌ NEVER check for splits during CASCADE:
if [[ "$LINES" -gt 700 ]]; then  # FORBIDDEN during cascade!
    echo "SPLIT_REQUIRED"  # ABSOLUTELY NOT!
fi

# ❌ NEVER evaluate code quality during CASCADE:
grep -r "TODO" .  # NOT during cascade!
```

## 🚨 ORCHESTRATOR CASCADE TRANSITIONS

### FORBIDDEN Transitions During CASCADE:
```
❌ CASCADE_REINTEGRATION → CREATE_NEXT_SPLIT_INFRASTRUCTURE
❌ CASCADE_REINTEGRATION → SPLIT_PLANNING
❌ CASCADE_REINTEGRATION → ANALYZE_SPLIT_REQUIREMENTS
❌ MONITOR_REVIEWS (cascade) → CREATE_NEXT_SPLIT_INFRASTRUCTURE
❌ Any state → Split-related state (when cascade_mode=true)
```

### ALLOWED Transitions During CASCADE:
```
✅ CASCADE_REINTEGRATION → INTEGRATION (continue cascade)
✅ CASCADE_REINTEGRATION → PHASE_INTEGRATION (continue cascade)
✅ CASCADE_REINTEGRATION → PROJECT_INTEGRATION (continue cascade)
✅ CASCADE_REINTEGRATION → SPAWN_ENGINEERS_FOR_FIXES (new fixes found)
✅ CASCADE_REINTEGRATION → INTEGRATION_CODE_REVIEW (cascade complete)
```

## 🔴 ENFORCEMENT MECHANISM

### State Machine Enforcement:
```bash
# In any orchestrator state transition logic:
NEXT_STATE="$1"
CASCADE_MODE=$(jq -r '.cascade_coordination.cascade_mode' orchestrator-state.json)

if [[ "$CASCADE_MODE" == "true" ]]; then
    # Check if transition violates cascade focus
    if [[ "$NEXT_STATE" =~ SPLIT|split ]]; then
        echo "🚨🚨🚨 R353 VIOLATION: Cannot transition to $NEXT_STATE during CASCADE!"
        echo "CASCADE FOCUS MUST BE MAINTAINED!"
        exit 353
    fi
fi
```

### Code Reviewer Enforcement:
```bash
# At the start of any Code Reviewer operation:
CASCADE_MODE="${CASCADE_MODE:-false}"

if [[ "$CASCADE_MODE" == "true" ]]; then
    echo "🔴 R353 CASCADE FOCUS PROTOCOL ACTIVE"
    echo "- Size checks: DISABLED"
    echo "- Split evaluation: DISABLED"
    echo "- Quality analysis: MINIMAL"
    echo "- Focus: REBASE VALIDATION ONLY"
    
    # Override any size checking commands
    alias line-counter.sh="echo 'R353: Size checks disabled during CASCADE'"
fi
```

## 🚨 DEFERRED EVALUATIONS

### What Gets Deferred During CASCADE:
1. **Size Measurements** - Check AFTER cascade completes
2. **Split Requirements** - Evaluate AFTER fixes propagate
3. **Code Quality** - Assess AFTER cascade stability
4. **Performance Testing** - Run AFTER cascade done
5. **Documentation Updates** - Write AFTER cascade

### Post-CASCADE Evaluation:
```bash
# After CASCADE_REINTEGRATION completes:
if [[ "$PREVIOUS_STATE" == "CASCADE_REINTEGRATION" ]]; then
    echo "📋 CASCADE COMPLETE - Running deferred evaluations:"
    echo "1. Size compliance checks"
    echo "2. Split requirement analysis"
    echo "3. Full code quality review"
    echo "4. Performance validation"
    # NOW these checks are allowed
fi
```

## 🔴 GRADING CRITERIA

### PASS Conditions (+100%):
- ✅ NO size checks during cascade operations
- ✅ NO split evaluations while cascading
- ✅ ALL agents respect cascade_mode flag
- ✅ Focus maintained until cascade complete
- ✅ Deferred evaluations run AFTER cascade

### FAIL Conditions (-100%):
- ❌ Code Reviewer demands split during CASCADE
- ❌ Orchestrator transitions to split state during CASCADE
- ❌ Size measurements taken during CASCADE
- ❌ CASCADE interrupted for "quick split"
- ❌ Cascade_mode context not propagated to agents

## 🚨 CRITICAL EXAMPLE - THE FAILURE CASE

### ❌ WHAT WENT WRONG (From Transcript):
```
During CASCADE_REINTEGRATION:
1. Orchestrator spawns Code Reviewers for rebased efforts
2. Code Reviewers measure: "3646 lines - SPLIT_REQUIRED!"
3. Orchestrator: "Oh, need to handle splits"
4. Transitions to CREATE_NEXT_SPLIT_INFRASTRUCTURE
5. CASCADE DISRUPTED! User intervention required
```

### ✅ WHAT SHOULD HAPPEN (With R353):
```
During CASCADE_REINTEGRATION:
1. Orchestrator spawns Code Reviewers with cascade_mode=true
2. Code Reviewers: "CASCADE MODE - skipping size checks"
3. Code Reviewers: "Rebase successful, no conflicts"
4. Orchestrator continues CASCADE to next level
5. CASCADE COMPLETES! Size checks deferred
```

## 🔴 IMPLEMENTATION CHECKLIST

### For Orchestrator:
- [ ] Set cascade_mode=true when entering CASCADE_REINTEGRATION
- [ ] Pass cascade_mode to ALL spawned agents
- [ ] Block transitions to split states during CASCADE
- [ ] Clear cascade_mode only after CASCADE complete
- [ ] Schedule deferred evaluations post-CASCADE

### For Code Reviewer:
- [ ] Check for cascade_mode in context/environment
- [ ] SKIP line counting when cascade_mode=true
- [ ] SKIP split evaluation when cascade_mode=true  
- [ ] ONLY validate rebase and build success
- [ ] Return simple CASCADE verdict (not split requirements)

### For State Machine:
- [ ] Add cascade_mode checks to transition validation
- [ ] Prevent split-related transitions during CASCADE
- [ ] Ensure CASCADE states maintain focus
- [ ] Document deferred evaluation states

## 🚨 THE CASCADE FOCUS MANTRA

```
When fixes cascade through the tree,
No diversions shall there be.
No counts of lines, no splits to plan,
Just propagate as fast as you can.

Focus sharp, the cascade flows,
Through wave and phase and project it goes.
Only when it's fully done,
Can size checks and splits begun.
```

## RELATED RULES

- **R327**: Mandatory Re-Integration After Fixes (triggers cascade)
- **R348**: Cascade State Transitions (state flow)
- **R352**: Overlapping Cascade Protocol (multiple chains)
- **R354**: Post-Rebase Review Requirement (mandatory reviews after rebases)
- **R220**: Line Count Limits (DEFERRED during cascade)
- **R221**: Split Requirements (DEFERRED during cascade)

---

**REMEMBER:** CASCADE FOCUS IS ABSOLUTE. When propagating fixes, NOTHING else matters until the cascade completes. Size can wait. Splits can wait. Quality can wait. The CASCADE CANNOT WAIT.

**Violation of R353 = -100% AUTOMATIC FAILURE**