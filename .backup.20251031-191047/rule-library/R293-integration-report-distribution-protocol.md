# 🚨🚨🚨 RULE R293: INTEGRATE_WAVE_EFFORTS REPORT DISTRIBUTION PROTOCOL (BLOCKING)

## PRIORITY: BLOCKING

This rule is **MANDATORY** and prevents SW Engineer confusion about which fix plans to follow during integration fix work.

## WHEN THIS RULE APPLIES

This rule applies when:
1. Any integration review completes (phase, wave, or project)
2. Integration fixes are required (INTEGRATE_WAVE_EFFORTS-REPORT.md created)
3. Orchestrator is preparing to spawn SW Engineers for fix work
4. Orchestrator is in states: INTEGRATE_WAVE_EFFORTS_FEEDBACK_REVIEW, SPAWN_SW_ENGINEERS, or similar

## MANDATORY REQUIREMENTS

### 1. REPORT DISTRIBUTION (BEFORE SPAWNING)

**The orchestrator MUST distribute the integration report BEFORE spawning SW Engineers:**

```bash
# 🚨 MANDATORY: Copy INTEGRATE_WAVE_EFFORTS-REPORT.md to ALL affected effort directories
for effort_dir in wave-*/effort-*; do
    if [[ -f "INTEGRATE_WAVE_EFFORTS-REPORT.md" ]]; then
        echo "📋 Distributing integration report to: $effort_dir"
        cp INTEGRATE_WAVE_EFFORTS-REPORT.md "$effort_dir/INTEGRATE_WAVE_EFFORTS-REPORT.md"
        
        # Verify the copy
        if [[ ! -f "$effort_dir/INTEGRATE_WAVE_EFFORTS-REPORT.md" ]]; then
            echo "❌ CRITICAL: Failed to copy INTEGRATE_WAVE_EFFORTS-REPORT.md to $effort_dir"
            exit 1
        fi
    fi
done
```

### 2. OLD PLAN ARCHIVAL (PREVENT CONFUSION)

**The orchestrator MUST archive old fix plans to prevent confusion:**

```bash
# Archive old fix plans in each effort directory
for effort_dir in wave-*/effort-*; do
    cd "$effort_dir"
    
    # Archive SPLIT-PLAN files (both old and new format)
    # Archive timestamped versions
    for split_plan in SPLIT-PLAN-*.md; do
        if [[ -f "$split_plan" && ! "$split_plan" =~ COMPLETED ]]; then
            base_name="${split_plan%.md}"
            mv "$split_plan" "${base_name}-COMPLETED-$(date +%Y%m%d-%H%M%S).md"
            echo "✅ Archived: $split_plan in $effort_dir"
        fi
    done
    
    # Archive legacy format if exists
    if [[ -f "SPLIT-PLAN.md" ]]; then
        mv SPLIT-PLAN.md SPLIT-PLAN-COMPLETED-$(date +%Y%m%d-%H%M%S).md
        echo "✅ Archived: SPLIT-PLAN.md (legacy) in $effort_dir"
    fi
    
    # Archive CODE-REVIEW-REPORT files (both old and new format)
    # Archive timestamped versions
    for review_report in CODE-REVIEW-REPORT-*.md; do
        if [[ -f "$review_report" && ! "$review_report" =~ COMPLETED ]]; then
            base_name="${review_report%.md}"
            mv "$review_report" "${base_name}-COMPLETED-$(date +%Y%m%d-%H%M%S).md"
            echo "✅ Archived: $review_report in $effort_dir"
        fi
    done
    
    # Archive legacy format if exists
    if [[ -f "CODE-REVIEW-REPORT.md" ]]; then
        mv CODE-REVIEW-REPORT.md CODE-REVIEW-REPORT-COMPLETED-$(date +%Y%m%d-%H%M%S).md
        echo "✅ Archived: CODE-REVIEW-REPORT.md (legacy) in $effort_dir"
    fi
    
    # Archive any other fix plans
    for plan in *-FIX-PLAN.md *-REVIEW.md; do
        if [[ -f "$plan" && "$plan" != "INTEGRATE_WAVE_EFFORTS-REPORT.md" ]]; then
            mv "$plan" "${plan%.md}-COMPLETED-$(date +%Y%m%d-%H%M%S).md"
            echo "✅ Archived: $plan in $effort_dir"
        fi
    done
    
    cd - > /dev/null
done
```

### 3. CLEAR SPAWN INSTRUCTIONS

**When spawning SW Engineers for integration fixes, the orchestrator MUST provide:**

```markdown
# In the spawn command message:

🔴🔴🔴 CRITICAL FIX INSTRUCTIONS:
- YOU ARE IN STATE: FIX_INTEGRATE_WAVE_EFFORTS_ISSUES
- FOLLOW ONLY: INTEGRATE_WAVE_EFFORTS-REPORT.md (located in your effort directory)
- IGNORE: Any files named *-COMPLETED-*.md (these are from previous fix cycles)
- YOUR TASK: Fix ONLY the issues listed in INTEGRATE_WAVE_EFFORTS-REPORT.md for your effort
🔴🔴🔴
```

### 4. STATE CLARITY REQUIREMENT

**The orchestrator MUST explicitly specify:**
- The EXACT state the SW Engineer is in (e.g., FIX_INTEGRATE_WAVE_EFFORTS_ISSUES)
- The EXACT plan file to follow (e.g., INTEGRATE_WAVE_EFFORTS-REPORT.md)
- A WARNING about ignoring archived plans

## VERIFICATION CHECKLIST

Before spawning SW Engineers for integration fixes:
- [ ] INTEGRATE_WAVE_EFFORTS-REPORT.md copied to ALL affected effort directories
- [ ] Old fix plans archived with COMPLETED timestamp
- [ ] Spawn instructions specify exact state
- [ ] Spawn instructions specify exact plan file
- [ ] Spawn instructions include warning about archived plans

## VIOLATIONS

**This rule is violated if:**
- ❌ SW Engineers spawned WITHOUT distributing INTEGRATE_WAVE_EFFORTS-REPORT.md
- ❌ Old fix plans NOT archived (causing confusion)
- ❌ Spawn instructions don't specify which plan to follow
- ❌ Spawn instructions don't specify the state
- ❌ Multiple active fix plans exist in effort directories

## PENALTIES

- Violation of this rule = -30% grading penalty
- SW Engineer confusion due to multiple plans = -50% penalty
- Wrong fixes applied due to following old plans = -75% penalty

## EXAMPLE COMPLIANT IMPLEMENTATION

```bash
# Orchestrator in INTEGRATE_WAVE_EFFORTS_FEEDBACK_REVIEW state:

# 1. Distribute the integration report
echo "📋 Distributing INTEGRATE_WAVE_EFFORTS-REPORT.md to all efforts..."
for effort_dir in wave-*/effort-*; do
    cp INTEGRATE_WAVE_EFFORTS-REPORT.md "$effort_dir/"
    echo "✅ Copied to $effort_dir"
done

# 2. Archive old plans
echo "📦 Archiving old fix plans..."
for effort_dir in wave-*/effort-*; do
    cd "$effort_dir"
    [[ -f "SPLIT-PLAN.md" ]] && mv SPLIT-PLAN.md SPLIT-PLAN-COMPLETED-$(date +%Y%m%d-%H%M%S).md
    [[ -f "CODE-REVIEW-REPORT.md" ]] && mv CODE-REVIEW-REPORT.md CODE-REVIEW-REPORT-COMPLETED-$(date +%Y%m%d-%H%M%S).md
    cd - > /dev/null
done

# 3. Spawn with clear instructions
/orchestrate spawn sw-engineer effort-database-schema --state FIX_INTEGRATE_WAVE_EFFORTS_ISSUES \
    --message "🔴 CRITICAL: Follow ONLY INTEGRATE_WAVE_EFFORTS-REPORT.md in your effort directory. You are in FIX_INTEGRATE_WAVE_EFFORTS_ISSUES state. Ignore any *-COMPLETED-*.md files."
```

## RELATED RULES

- R294: Fix Plan Archival Protocol
- R295: SW Engineer Spawn Clarity Protocol
- R052: Agent Spawning Protocol
- R282: Phase Integration Protocol
- R283: Project Integration Protocol

---

**REMEMBER**: Clear communication prevents wasted effort. ALWAYS distribute reports, archive old plans, and provide crystal-clear instructions!