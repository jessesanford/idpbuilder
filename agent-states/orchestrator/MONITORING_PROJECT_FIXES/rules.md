# Orchestrator - MONITORING_PROJECT_FIXES State Rules

## 🛑🛑🛑 R232 MONITOR STATE REQUIREMENTS 🛑🛑🛑

**CRITICAL**: Before ANY transition from MONITOR_* states, you MUST:
1. Check TodoWrite for pending items
2. Process ALL pending items IMMEDIATELY
3. NO "I will..." statements - only "I am..." with action
4. VIOLATION = AUTOMATIC FAILURE

---

## State Context

**Purpose:**
Monitor SW Engineers as they fix bugs found during project integration.

## Primary Actions

1. **Track Fix Progress**:
   ```bash
   # For each engineer fixing bugs:
   cd /efforts/phase-X-wave-Y-effort-Z
   git status
   git log --oneline -5
   
   # Check for fix commits
   git log --grep="project integration bug"
   ```

2. **Monitor Fix Completion**:
   - Check each effort directory for completed fixes
   - Verify commits match bug fix requirements
   - Track which bugs are resolved

3. **Update Tracking**:
   ```json
   {
     "project_fixes_in_progress": [
       {
         "bug_id": 1,
         "engineer": "sw-engineer-1",
         "branch": "phase-1-wave-2-effort-3",
         "status": "completed",
         "commit": "abc123"
       }
     ]
   }
   ```

4. **Spawn Code Reviewers for Fixed Bugs**:
   - Once a bug fix is complete, spawn Code Reviewer
   - Ensure fix meets requirements
   - Verify no new issues introduced

## Monitoring Protocol

```markdown
## 📊 PROJECT FIX MONITORING STATUS

### Fixes In Progress:
- Bug #1: COMPLETED ✅ (commit: abc123)
- Bug #2: IN_PROGRESS 🔄
- Bug #4: COMPLETED ✅ (commit: def456)

### Pending Reviews:
- Bug #1: Review spawned
- Bug #4: Review spawned

### Sequential Fixes Waiting:
- Bug #6: Waiting for bug #1 review to pass
- Bug #7: Waiting for bug #6 completion

### Next Actions:
- Monitor bug #2 completion
- Check review results for bugs #1 and #4
- Spawn sequential fix for bug #6 if ready
```

## Decision Logic

```python
def determine_next_state():
    all_fixes = load_project_fixes()
    
    # Check if all fixes are complete
    all_complete = all(fix.status == "completed" for fix in all_fixes)
    
    # Check if all reviews passed
    all_reviewed = all(fix.review_status == "passed" for fix in all_fixes)
    
    if not all_complete:
        return "MONITORING_PROJECT_FIXES"  # Keep monitoring
    
    if not all_reviewed:
        return "SPAWN_CODE_REVIEWERS_FOR_REVIEW"  # Review fixes
    
    # All fixes complete and reviewed
    return "PROJECT_INTEGRATION"  # Re-run integration with fixed code
```

## Valid State Transitions

- **FIXES_ONGOING** → MONITORING_PROJECT_FIXES (continue monitoring)
- **FIXES_COMPLETE** → SPAWN_CODE_REVIEWERS_FOR_REVIEW (review all fixes)
- **ALL_REVIEWED** → PROJECT_INTEGRATION (re-run with fixed code)
- **FIXES_FAILED** → ERROR_RECOVERY (unable to fix bugs)

## Critical Requirements

1. **R232 Compliance** - Process all TODOs before considering transition
2. **Track every fix** - Know status of each bug fix
3. **Verify in source branches** - Fixes must be in original branches
4. **Spawn reviews** - All fixes need code review
5. **Re-run integration** - After fixes, must re-integrate

## Integration Re-run Protocol

After all fixes complete and pass review:
1. Transition to PROJECT_INTEGRATION
2. Re-run entire project integration with fixed branches
3. Check if new bugs appear
4. If clean, proceed to validation
5. If new bugs, repeat fix cycle

## Grading Impact

- **-50%** for R232 violation (not processing TODOs)
- **-30%** if skipping code reviews for fixes
- **-40%** if not re-running integration after fixes
- **+20%** for comprehensive fix tracking
- **+15%** for efficient review spawning