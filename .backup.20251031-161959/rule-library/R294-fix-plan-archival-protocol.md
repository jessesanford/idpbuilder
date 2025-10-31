# 🚨🚨🚨 RULE R294: FIX PLAN ARCHIVAL PROTOCOL (BLOCKING)

## PRIORITY: BLOCKING

This rule ensures that completed fix plans are properly archived to prevent confusion and maintain clear task tracking.

## WHEN THIS RULE APPLIES

This rule applies when:
1. SW Engineers complete fixes from any fix plan
2. Code Reviewer completes review of fixes
3. Orchestrator transitions out of fix-related states
4. New fix plans are being distributed (must archive old ones first)

## MANDATORY ARCHIVAL REQUIREMENTS

### 1. SW ENGINEER COMPLETION ARCHIVAL

**When SW Engineers complete fixes, they MUST archive the plan:**

```bash
# In FIX_ISSUES state completion:
if [[ -f "CODE-REVIEW-REPORT.md" ]]; then
    mv CODE-REVIEW-REPORT.md CODE-REVIEW-REPORT-COMPLETED-$(date +%Y%m%d-%H%M%S).md
    echo "✅ Archived completed review report"
fi

# In SPLIT_IMPLEMENTATION state completion:
if [[ -f "SPLIT-PLAN.md" ]]; then
    mv SPLIT-PLAN.md SPLIT-PLAN-COMPLETED-$(date +%Y%m%d-%H%M%S).md
    echo "✅ Archived completed split plan"
fi

# In FIX_INTEGRATE_WAVE_EFFORTS_ISSUES state completion:
if [[ -f "INTEGRATE_WAVE_EFFORTS-REPORT.md" ]]; then
    mv INTEGRATE_WAVE_EFFORTS-REPORT.md INTEGRATE_WAVE_EFFORTS-REPORT-COMPLETED-$(date +%Y%m%d-%H%M%S).md
    echo "✅ Archived completed integration report"
fi
```

### 2. ORCHESTRATOR DISTRIBUTION ARCHIVAL

**Before distributing new fix plans, orchestrator MUST archive existing ones:**

```bash
# Function to archive all active fix plans
archive_active_fix_plans() {
    local effort_dir="$1"
    local timestamp=$(date +%Y%m%d-%H%M%S)
    
    cd "$effort_dir"
    
    # List of plan file patterns to archive (supports both old and new formats)
    local plan_patterns=(
        "SPLIT-PLAN*.md"           # Includes SPLIT-PLAN.md and SPLIT-PLAN-YYYYMMDD-HHMMSS.md
        "CODE-REVIEW-REPORT*.md"   # Includes all review reports
        "INTEGRATE_WAVE_EFFORTS-REPORT*.md"   # Includes all integration reports  
        "FIX-INSTRUCTIONS*.md"     # Includes all fix instructions
        "IMPLEMENTATION-PLAN*.md"  # Includes all implementation plans
        "PHASE-FIX-PLAN*.md"       # Phase fix plans
        "WAVE-FIX-PLAN*.md"        # Wave fix plans
    )
    
    # Archive files matching each pattern
    for pattern in "${plan_patterns[@]}"; do
        for plan_file in $pattern; do
            # Skip if file doesn't exist or is already archived
            if [[ -f "$plan_file" && ! "$plan_file" =~ COMPLETED ]]; then
                local archived_name="${plan_file%.md}-COMPLETED-${timestamp}.md"
                mv "$plan_file" "$archived_name"
                echo "📦 Archived: $plan_file → $archived_name"
            
                # Add archival note to the file
                echo -e "\n---\n📦 ARCHIVED: $(date '+%Y-%m-%d %H:%M:%S')\nReason: New fix cycle starting\n---" >> "$archived_name"
            fi
        done
    done
    
    cd - > /dev/null
}
```

### 3. ARCHIVAL NAMING CONVENTION

**All archived plans MUST follow this naming pattern:**
```
[ORIGINAL-NAME]-COMPLETED-[YYYYMMDD-HHMMSS].md
```

Examples:
- `SPLIT-PLAN-COMPLETED-20250131-143022.md`
- `CODE-REVIEW-REPORT-COMPLETED-20250131-145500.md`
- `INTEGRATE_WAVE_EFFORTS-REPORT-COMPLETED-20250131-151230.md`

### 4. ARCHIVAL VERIFICATION

**After archiving, MUST verify:**

```bash
verify_archival() {
    local effort_dir="$1"
    
    # Check no active fix plans remain
    local active_plans=$(ls "$effort_dir"/*.md 2>/dev/null | grep -E "(SPLIT-PLAN|CODE-REVIEW-REPORT|INTEGRATE_WAVE_EFFORTS-REPORT|FIX-PLAN)\.md$" | grep -v COMPLETED)
    
    if [[ -n "$active_plans" ]]; then
        echo "❌ ERROR: Active fix plans still exist after archival:"
        echo "$active_plans"
        return 1
    fi
    
    echo "✅ All fix plans properly archived"
    return 0
}
```

## STATE-SPECIFIC ARCHIVAL RULES

### SW Engineer States:
- **FIX_ISSUES**: Archive CODE-REVIEW-REPORT.md on completion
- **SPLIT_IMPLEMENTATION**: Archive SPLIT-PLAN.md on completion
- **FIX_INTEGRATE_WAVE_EFFORTS_ISSUES**: Archive INTEGRATE_WAVE_EFFORTS-REPORT.md on completion

### Orchestrator States:
- **INTEGRATE_WAVE_EFFORTS_FEEDBACK_REVIEW**: Archive old plans before distributing new
- **SPAWN_SW_ENGINEERS**: Verify plans archived before spawning
- **WAVE_COMPLETE**: Archive all effort fix plans
- **COMPLETE_PHASE**: Archive all wave and effort fix plans

## CLEANUP PROTOCOL

**Periodically clean up old archived plans (optional):**

```bash
# Keep only last 5 archived versions of each plan type
cleanup_old_archives() {
    local effort_dir="$1"
    
    cd "$effort_dir"
    
    # For each plan type, keep only the 5 most recent
    for prefix in "SPLIT-PLAN" "CODE-REVIEW-REPORT" "INTEGRATE_WAVE_EFFORTS-REPORT"; do
        ls -t ${prefix}-COMPLETED-*.md 2>/dev/null | tail -n +6 | xargs -r rm -v
    done
    
    cd - > /dev/null
}
```

## VIOLATIONS

**This rule is violated if:**
- ❌ Fix plans not archived after completion
- ❌ New plans distributed without archiving old ones
- ❌ Wrong naming convention used for archives
- ❌ Active and archived plans coexist (causing confusion)

## PENALTIES

- Not archiving completed plans = -20% penalty
- Confusion from multiple active plans = -40% penalty
- SW Engineer follows wrong plan due to missing archival = -60% penalty

## EXAMPLE COMPLIANT WORKFLOW

```bash
# 1. SW Engineer completes fixes
echo "✅ All fixes from CODE-REVIEW-REPORT.md completed"
mv CODE-REVIEW-REPORT.md CODE-REVIEW-REPORT-COMPLETED-$(date +%Y%m%d-%H%M%S).md
git add -A
git commit -m "fix: complete review fixes and archive plan"
git push

# 2. Orchestrator distributes new integration fixes
for effort_dir in wave-*/effort-*; do
    # Archive old plans first
    archive_active_fix_plans "$effort_dir"
    
    # Then copy new plan
    cp INTEGRATE_WAVE_EFFORTS-REPORT.md "$effort_dir/"
    
    # Verify no conflicts
    verify_archival "$effort_dir"
done

# 3. Clear spawn message
/orchestrate spawn sw-engineer effort-001 \
    --message "Follow ONLY INTEGRATE_WAVE_EFFORTS-REPORT.md. All previous plans have been archived as COMPLETED."
```

## RELATED RULES

- R293: Integration Report Distribution Protocol
- R295: SW Engineer Spawn Clarity Protocol
- R287: TODO Persistence Comprehensive

---

**REMEMBER**: Completed work should be clearly marked as completed. Archive fix plans immediately upon completion!