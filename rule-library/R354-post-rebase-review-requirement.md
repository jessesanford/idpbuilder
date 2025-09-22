# 🔴🔴🔴 RULE R354: POST-REBASE AND FIX CASCADE REVIEW REQUIREMENT 🔴🔴🔴

## SUPREME LAW - EVERY REBASE AND FIX REQUIRES IMMEDIATE REVIEW

**CRITICALITY:** SUPREME LAW - Violation = -100% AUTOMATIC FAILURE
**PRIORITY:** P0 - HIGHEST
**ENFORCEMENT:** MANDATORY - NO EXCEPTIONS

## 🚨🚨🚨 THE ABSOLUTE POST-CHANGE REVIEW LAW 🚨🚨🚨

**EVERY CODE CHANGE DURING CASCADE OR FIX CASCADE MUST BE REVIEWED!**

This rule applies to:
1. **REBASE OPERATIONS** - When code is rebased onto a new base branch
2. **FIX CASCADE OPERATIONS** - When fixes are backported or forward-ported
3. **CONFLICT RESOLUTIONS** - When merge conflicts are resolved
4. **ANY CODE MODIFICATION** - During cascade or fix cascade operations

The review validates that:
- Changes were applied correctly
- No integration issues were introduced
- Build and tests still pass
- The fix actually resolves the problem (for fix cascades)

## 🔴 CORE PRINCIPLE: REBASES CAN BREAK CODE

### Why Post-Rebase Reviews Are MANDATORY:
1. **Merge Conflicts** - Auto-resolution can choose wrong code
2. **Semantic Conflicts** - Code compiles but behavior changed
3. **Lost Changes** - Rebase can accidentally drop commits
4. **Integration Issues** - New base may have incompatible changes
5. **Build Failures** - Dependencies may have changed
6. **Test Failures** - Integration tests may reveal issues

### What Makes This DIFFERENT from Regular Reviews:
- **Focus**: Integration correctness, not code quality
- **Scope**: Changes introduced by rebase, not original implementation
- **Priority**: Build/test success over style issues
- **Context**: CASCADE_MODE changes review behavior

## 🔴🔴🔴 FIX CASCADE ENFORCEMENT PROTOCOL 🔴🔴🔴

### CRITICAL: FIX CASCADE QUALITY GATES

During fix cascade operations (`/fix-cascade` command), EVERY code change MUST be reviewed:

1. **AFTER BACKPORT TO ANY BRANCH:**
   - Code review REQUIRED before continuing
   - Focus: Verify fix applied correctly
   - Check: Build passes, tests pass
   - Verify: Fix actually resolves the issue

2. **AFTER FORWARD PORT TO ANY BRANCH:**
   - Code review REQUIRED before continuing
   - Focus: Verify fix integrates properly
   - Check: No new conflicts introduced
   - Verify: Fix still works in new context

3. **AFTER CONFLICT RESOLUTION:**
   - Code review REQUIRED for resolved conflicts
   - Focus: Correct resolution chosen
   - Check: Both sides properly merged
   - Verify: No code lost or broken

4. **COMPREHENSIVE FINAL REVIEW:**
   - Review ALL accumulated changes
   - Verify fix solves original problem
   - Check for unintended side effects
   - Ensure code quality standards met

### FIX CASCADE REVIEW ENFORCEMENT:

```bash
# During BACKPORT state
after_backport() {
    local branch="$1"
    local fix_commits="$2"

    echo "🔴 R354: Fix cascade backport review REQUIRED"

    # Add to pending reviews
    jq --arg branch "$branch" --arg commits "$fix_commits" '
        .fix_cascade.pending_reviews += [{
            "branch": $branch,
            "operation": "backport",
            "commits": $commits,
            "review_status": "pending",
            "r354_mandated": true,
            "timestamp": now | todate
        }]' orchestrator-state.json > tmp.json && mv tmp.json orchestrator-state.json

    # MUST spawn Code Reviewer
    echo "🚨 BLOCKING: Cannot continue fix cascade until review completes"
    /spawn-code-reviewer --fix-cascade-mode=true --review-type=fix-backport --branch=$branch

    # Wait for review completion
    wait_for_review_completion "$branch"
}

# During FORWARD_PORT state
after_forward_port() {
    local branch="$1"
    local fix_commits="$2"

    echo "🔴 R354: Fix cascade forward-port review REQUIRED"

    # Similar enforcement as backport
    jq --arg branch "$branch" --arg commits "$fix_commits" '
        .fix_cascade.pending_reviews += [{
            "branch": $branch,
            "operation": "forward_port",
            "commits": $commits,
            "review_status": "pending",
            "r354_mandated": true,
            "timestamp": now | todate
        }]' orchestrator-state.json > tmp.json && mv tmp.json orchestrator-state.json

    # MUST spawn Code Reviewer
    echo "🚨 BLOCKING: Cannot continue fix cascade until review completes"
    /spawn-code-reviewer --fix-cascade-mode=true --review-type=fix-forward-port --branch=$branch

    wait_for_review_completion "$branch"
}

# During VALIDATION state
comprehensive_fix_review() {
    echo "🔴 R354: Comprehensive fix cascade review REQUIRED"

    # Review ALL changes made during fix cascade
    jq '.fix_cascade | {
        "branches_modified": .branches_modified,
        "total_changes": .total_commits,
        "conflicts_resolved": .conflicts_resolved
    }' orchestrator-state.json

    # Spawn comprehensive reviewer
    /spawn-code-reviewer --fix-cascade-mode=true --review-type=comprehensive-fix-validation

    echo "📊 Comprehensive review checklist:"
    echo "  ✓ All branches build successfully"
    echo "  ✓ All tests pass on all branches"
    echo "  ✓ Fix resolves original issue"
    echo "  ✓ No regressions introduced"
    echo "  ✓ Code quality maintained"
}
```

### FIX CASCADE BLOCKING CONDITIONS:

```bash
# Cannot proceed if reviews pending
check_fix_cascade_reviews() {
    PENDING=$(jq '.fix_cascade.pending_reviews | length' orchestrator-state.json)

    if [[ "$PENDING" -gt 0 ]]; then
        echo "🛑 R354 VIOLATION: $PENDING fix cascade reviews pending"
        echo "❌ CANNOT CONTINUE fix cascade until ALL reviews complete"

        jq -r '.fix_cascade.pending_reviews[] |
            "  - Branch: \(.branch) (\(.operation)) - Status: \(.review_status)"' \
            orchestrator-state.json

        return 354  # R354 violation code
    fi
}
```

## 🚨 CASCADE MODE POST-REBASE PROTOCOL

### During CASCADE_REINTEGRATION:

**AFTER EVERY REBASE:**
```bash
# 1. Complete the rebase operation
cd $EFFORT_DIR
git rebase $NEW_BASE
git push --force-with-lease

# 2. MANDATORY: Mark review required
jq --arg effort "$EFFORT" '
    .cascade_coordination.pending_reviews += [{
        "effort": $effort,
        "type": "post_rebase",
        "review_required": true,
        "review_status": "pending",
        "rebased_at": now | todate,
        "cascade_mode": true
    }]' orchestrator-state.json > tmp.json && mv tmp.json orchestrator-state.json

# 3. MANDATORY: Spawn Code Reviewer
echo "🔴 R354: Post-rebase review REQUIRED"
/spawn-code-reviewer --cascade-mode=true --review-type=post-rebase --effort=$EFFORT

# 4. CANNOT PROCEED until review completes
echo "⏸️ CASCADE PAUSED: Awaiting post-rebase review per R354"
```

### Code Reviewer Behavior in Post-Rebase Review:

**WHEN review-type=post-rebase AND cascade-mode=true:**
```bash
echo "🔴 R354 POST-REBASE REVIEW PROTOCOL ACTIVE"
echo "Focus: Integration validation, NOT quality assessment"

# 1. Validate rebase success
cd $EFFORT_DIR
git log --oneline -10  # Check commit history preserved
git status  # Ensure clean state

# 2. Check for integration issues ONLY
echo "🔍 Checking rebase integration..."
# ✅ Build still passes
make build || echo "❌ BUILD FAILURE - Rebase broke build!"

# ✅ Tests still pass
make test || echo "❌ TEST FAILURE - Rebase broke tests!"

# ✅ No conflict markers
grep -r "<<<<<<" . && echo "❌ CONFLICT MARKERS - Rebase incomplete!"

# ✅ Dependencies resolved
go mod tidy || npm install || echo "⚠️ Dependency issues"

# 3. Skip quality checks per R353
echo "📊 Size checks SKIPPED (R353 - cascade focus)"
echo "🎨 Style checks SKIPPED (R353 - cascade focus)"
echo "📝 Documentation checks SKIPPED (R353 - cascade focus)"

# 4. Return verdict
if [[ $BUILD_PASSES && $TESTS_PASS && $NO_CONFLICTS ]]; then
    echo "✅ POST_REBASE_VALID - Integration successful"
    VERDICT="REBASE_VALID"
else
    echo "❌ POST_REBASE_FAILED - Issues found"
    VERDICT="FIXES_NEEDED"
    # Document specific fixes needed
fi
```

## 🔴 CASCADE CONTINUATION WITH FIXES

### If Post-Rebase Review Finds Issues:

**CRITICAL: Fixes discovered during post-rebase review BECOME PART OF THE CASCADE!**

```bash
# 1. Code Reviewer identifies issues
echo "❌ Post-rebase review found issues:"
echo "- Build failure in module X"
echo "- Test failure in integration test Y"

# 2. Create fix requirements
cat > POST-REBASE-FIXES.md << EOF
## Post-Rebase Issues Found
- Build failure: Missing import after rebase
- Test failure: API signature changed in base
- Action: SW Engineer must fix before cascade continues
EOF

# 3. Orchestrator spawns SW Engineer for fixes
/spawn-sw-engineer --cascade-mode=true --fix-type=post-rebase --effort=$EFFORT

# 4. After fixes applied, fixes JOIN the cascade
jq --arg effort "$EFFORT" --arg fixes "$FIX_COMMITS" '
    .cascade_coordination.cascade_fixes += [{
        "source": "post_rebase_review",
        "effort": $effort,
        "fixes": ($fixes | split(" ")),
        "must_cascade": true,
        "added_to_cascade": now | todate
    }]' orchestrator-state.json > tmp.json && mv tmp.json orchestrator-state.json

# 5. CASCADE CONTINUES with original fixes + new fixes
echo "🔄 CASCADE EXPANDED: Now includes post-rebase fixes"
```

## 🚨 ORCHESTRATOR ENFORCEMENT

### In CASCADE_REINTEGRATION State:

```bash
# Function to enforce R354
enforce_post_rebase_review() {
    local effort="$1"
    local rebased_to="$2"

    echo "🔴 R354 ENFORCEMENT: Post-rebase review required"

    # Add to pending reviews
    jq --arg e "$effort" --arg base "$rebased_to" '
        .cascade_coordination.pending_reviews += [{
            "effort": $e,
            "rebased_to": $base,
            "review_type": "post_rebase",
            "r354_mandated": true,
            "timestamp": now | todate
        }]' orchestrator-state.json > tmp.json && mv tmp.json orchestrator-state.json

    # Set state to spawn reviewer
    jq '.cascade_coordination.r354_enforcement = {
        "active": true,
        "blocking_effort": $e,
        "waiting_for": "post_rebase_review"
    }' orchestrator-state.json > tmp.json && mv tmp.json orchestrator-state.json

    # Transition to review spawn state
    NEXT_STATE="SPAWN_CODE_REVIEWERS_FOR_REVIEW"
    echo "🔄 Transitioning to spawn Code Reviewer for post-rebase review"
}

# After EVERY rebase operation
if [[ "$REBASE_COMPLETED" == "true" ]]; then
    enforce_post_rebase_review "$EFFORT_NAME" "$NEW_BASE"
fi
```

### Transition Blocking:

```bash
# Cannot continue cascade without review
PENDING_REVIEWS=$(jq '.cascade_coordination.pending_reviews | length' orchestrator-state.json)

if [[ "$PENDING_REVIEWS" -gt 0 ]]; then
    echo "🛑 R354: Cannot continue cascade - $PENDING_REVIEWS reviews pending"
    echo "Reviews required for:"
    jq -r '.cascade_coordination.pending_reviews[] |
           "  - \(.effort) (rebased to \(.rebased_to))"' orchestrator-state.json

    # MUST spawn reviewers
    NEXT_STATE="SPAWN_CODE_REVIEWERS_FOR_REVIEW"
fi
```

## 🔴 MONITORING_INTEGRATION ENFORCEMENT

### Check for Pending Post-Rebase Reviews:

```bash
# In MONITORING_INTEGRATION state
check_r354_compliance() {
    echo "🔍 R354 Compliance Check: Post-rebase reviews"

    # Find efforts that were rebased but not reviewed
    for effort_dir in /efforts/*/*/effort-*; do
        if [[ -d "$effort_dir" ]]; then
            cd "$effort_dir"

            # Check git reflog for recent rebases
            RECENT_REBASE=$(git reflog --since="1 hour ago" | grep -c "rebase:")

            if [[ "$RECENT_REBASE" -gt 0 ]]; then
                EFFORT_NAME=$(basename "$effort_dir")

                # Check if review exists
                REVIEW_EXISTS=$(jq --arg e "$EFFORT_NAME" '
                    .cascade_coordination.completed_reviews[]? |
                    select(.effort == $e and .type == "post_rebase") |
                    .effort' orchestrator-state.json)

                if [[ -z "$REVIEW_EXISTS" ]]; then
                    echo "❌ R354 VIOLATION: $EFFORT_NAME rebased but not reviewed!"
                    echo "🔴 MANDATORY ACTION: Spawn Code Reviewer immediately"

                    # Add to violations
                    jq --arg e "$EFFORT_NAME" '
                        .r354_violations += [{
                            "effort": $e,
                            "violation": "rebased_without_review",
                            "detected_at": now | todate
                        }]' orchestrator-state.json > tmp.json && mv tmp.json orchestrator-state.json

                    return 354  # R354 violation code
                fi
            fi
        fi
    done

    echo "✅ R354 Compliance: All rebases have reviews"
}

# Run check
check_r354_compliance || {
    echo "🚨 R354 VIOLATION DETECTED - MUST RESOLVE!"
    NEXT_STATE="SPAWN_CODE_REVIEWERS_FOR_REVIEW"
}
```

## 🔴 GRADING CRITERIA

### PASS Conditions (+100%):
- ✅ EVERY rebase followed by review
- ✅ Code Reviewer spawned with cascade_mode=true
- ✅ Review focuses on integration, not quality (during cascade)
- ✅ Fixes from review added to cascade
- ✅ Cascade continues with expanded fix set

### FAIL Conditions (-100%):
- ❌ Rebase without review
- ❌ Proceeding with cascade before review completes
- ❌ Ignoring post-rebase review failures
- ❌ Not adding post-rebase fixes to cascade
- ❌ Losing cascade context during review

## 🚨 CRITICAL EXAMPLE - THE MISSING REVIEW

### ❌ WHAT WENT WRONG (From Transcript):
```
During CASCADE_REINTEGRATION:
1. Orchestrator rebased cli-commands onto Phase 2 Wave 1 integration
2. Orchestrator marked it complete and ready for WAVE_COMPLETE
3. User: "but we didn't have a code reviewer review the rebased cli-commands"
4. Orchestrator: "You're absolutely right! After the rebase, we need a fresh code review"
5. CASCADE INCOMPLETE - User intervention required
```

### ✅ WHAT SHOULD HAPPEN (With R354):
```
During CASCADE_REINTEGRATION:
1. Orchestrator rebases cli-commands onto Phase 2 Wave 1 integration
2. R354 TRIGGERS: "Post-rebase review required"
3. Orchestrator spawns Code Reviewer with cascade_mode=true
4. Code Reviewer validates rebase (skip quality per R353)
5. If issues found, fixes added to cascade
6. CASCADE CONTINUES with all fixes
```

## 🔴 IMPLEMENTATION CHECKLIST

### For Orchestrator:
- [ ] After EVERY rebase, add to pending_reviews
- [ ] Block cascade continuation until review completes
- [ ] Spawn Code Reviewer with cascade_mode AND review_type=post_rebase
- [ ] Add any fixes discovered to cascade tracking
- [ ] Continue cascade with expanded fix set

### For Code Reviewer:
- [ ] Detect review_type=post_rebase in context
- [ ] Focus on integration validation, not quality
- [ ] Check builds, tests, conflicts only
- [ ] Skip size/style/documentation checks (R353)
- [ ] Return simple verdict: REBASE_VALID or FIXES_NEEDED

### For State Machine:
- [ ] CASCADE_REINTEGRATION must check pending_reviews
- [ ] Block transitions if reviews pending
- [ ] MONITORING_INTEGRATION must detect R354 violations
- [ ] Track post-rebase reviews separately

## 🚨 THE POST-REBASE MANTRA

```
After every rebase done,
A review must be run.
Check the build, check the test,
Skip the style, skip the rest.

If the rebase broke the code,
Fix it now, don't let it erode.
Add those fixes to cascade's flow,
Until to project level they go.

Never skip, never forget,
Post-rebase review is set.
R354 stands guard at the gate,
Reviews required - cascade must wait!
```

## CASCADE TRACKING METADATA

### Required State Tracking:
```json
{
  "cascade_coordination": {
    "pending_reviews": [
      {
        "effort": "cli-commands",
        "type": "post_rebase",
        "rebased_to": "phase2/wave1/integration",
        "review_required": true,
        "review_status": "pending",
        "r354_mandated": true
      }
    ],
    "completed_reviews": [
      {
        "effort": "api-server",
        "type": "post_rebase",
        "verdict": "REBASE_VALID",
        "reviewed_at": "2024-01-20T10:30:00Z"
      }
    ],
    "cascade_fixes": [
      {
        "source": "post_rebase_review",
        "effort": "cli-commands",
        "fixes": ["abc123", "def456"],
        "must_cascade": true
      }
    ]
  }
}
```

## RELATED RULES

- **R353**: CASCADE FOCUS PROTOCOL (skip quality during cascade)
- **R352**: Overlapping Cascade Protocol (multiple cascade chains)
- **R327**: Mandatory Re-Integration After Fixes (triggers cascade)
- **R348**: Cascade State Transitions (state flow)
- **R351**: Cascade Execution Protocol (execution order)

---

**REMEMBER:** Every rebase can introduce subtle integration issues that weren't present in the original code. Post-rebase reviews are your safety net to catch these issues before they cascade through the entire system.

**Violation of R354 = -100% AUTOMATIC FAILURE**