# Orchestrator - MONITORING_PROJECT_FIXES State Rules

## 🔴🔴🔴 MANDATORY STATE RULE READING AND ACKNOWLEDGMENT 🔴🔴🔴

### ⚠️⚠️⚠️ YOU MUST READ EACH RULE FILE LISTED IN PRIMARY DIRECTIVES. **I AM WATCHING YOUR TOOL CALLS FOR READ OPERATIONS** *YOU WILL FAIL* IF YOU DO NOT MAKE A READ FILE CALL FOR EACH RULE FILE IN PRIMARY DIRECTIVES!!! ⚠️⚠️⚠️

**AFTER READING, YOU MUST ACKNOWLEDGE ALL THE STATE RULES AND STATE THAT YOU WILL ABIDE BY THEM ONE AT A TIME GIVING THE RULE NUMBER AND DESCRIPTION.**

## 📋 PRIMARY DIRECTIVES FOR MONITORING_PROJECT_FIXES STATE

### Core Mandatory Rules (ALL orchestrator states must have these):

1. **🚨🚨🚨 R006** - ORCHESTRATOR NEVER WRITES CODE OR PERFORMS FILE OPERATIONS (BLOCKING)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R006-orchestrator-never-writes-code.md`
   - Criticality: BLOCKING - Automatic termination, 0% grade
   - Summary: NEVER write, copy, move, or manipulate ANY code files - delegate ALL to agents

2. **🔴🔴🔴 R287** - TODO PERSISTENCE COMPREHENSIVE (SUPREME LAW)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R287-todo-persistence-comprehensive.md`
   - Criticality: SUPREME - -20% to -100% penalty for violations
   - Summary: MUST save TODOs within 30s after write, every 10 messages, before transitions

3. **🔴🔴🔴 R288** - STATE FILE UPDATE REQUIREMENTS (SUPREME LAW)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R288-state-file-update-requirements.md`
   - Criticality: SUPREME - State updates required for all transitions
   - Summary: MUST update orchestrator-state.json before EVERY state transition

4. **🔴🔴🔴 R322 Part A** - Mandatory Stop After Spawn States
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R322 Part A-mandatory-stop-after-spawn.md`
   - Criticality: SUPREME LAW - Must stop after spawning
   - Summary: ALL spawn states require STOP after spawning agents

### State-Specific Rules:

5. **🔴🔴🔴 R233** - Immediate Action On State Entry
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R233-immediate-action-on-state-entry.md`
   - Criticality: SUPREME LAW - Must act immediately on entering state
   - Summary: Monitoring states require active checking, not passive waiting

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
    
    # 🔴🔴🔴 CRITICAL: MUST RE-RUN FULL INTEGRATION 🔴🔴🔴
    # All fixes complete and reviewed - MUST go back to PROJECT_INTEGRATION
    # to re-run the ENTIRE integration with the fixed source branches!
    return "PROJECT_INTEGRATION"  # Re-run FULL integration with fixed code
```

## Valid State Transitions

- **FIXES_ONGOING** → MONITORING_PROJECT_FIXES (continue monitoring)
- **FIXES_COMPLETE** → SPAWN_CODE_REVIEWERS_FOR_REVIEW (review all fixes)
- **🔴 ALL_REVIEWED** → PROJECT_INTEGRATION (MANDATORY: Re-run FULL integration with fixed code)
- **FIXES_FAILED** → ERROR_RECOVERY (unable to fix bugs)

## Critical Requirements

1. **R232 Compliance** - Process all TODOs before considering transition
2. **Track every fix** - Know status of each bug fix
3. **Verify in source branches** - Fixes must be in original branches
4. **Spawn reviews** - All fixes need code review
5. **Re-run integration** - After fixes, must re-integrate

## 🔴🔴🔴 MANDATORY INTEGRATION RE-RUN PROTOCOL 🔴🔴🔴

**CRITICAL: After fixes are complete, you MUST re-run the ENTIRE project integration!**

### Why Re-Integration Is MANDATORY:
- Fixes were applied to UPSTREAM branches (phase/wave/effort branches)
- The project-integration branch still has the BROKEN code
- You MUST re-merge all branches to get the fixed code into integration
- Without re-integration, the binary CANNOT be built

### The CORRECT Re-Integration Cycle:
```
MONITORING_PROJECT_FIXES (all fixes complete & reviewed)
    ↓
PROJECT_INTEGRATION (delete old integration, create fresh)
    ↓
SPAWN_CODE_REVIEWER_PROJECT_MERGE_PLAN (create new merge plan)
    ↓
SPAWN_INTEGRATION_AGENT_PROJECT (re-merge ALL branches with fixes)
    ↓
MONITORING_PROJECT_INTEGRATION (check if NOW it works)
    ↓
If bugs found → SPAWN_CODE_REVIEWER_PROJECT_FIX_PLANNING → WAITING_FOR_PROJECT_FIX_PLANS → SPAWN_SW_ENGINEER_PROJECT_FIXES → MONITORING_PROJECT_FIXES
If clean → SPAWN_CODE_REVIEWER_PROJECT_VALIDATION → SUCCESS
```

### What Happens During Re-Integration:
1. **Delete old broken integration branch** (it has unfixed code)
2. **Create fresh project-integration from main**
3. **Re-run ENTIRE merge plan** (all phases with their fixed branches)
4. **All fixes from upstream branches now merged in**
5. **Binary can finally be built with working code**

### NEVER DO THESE (AUTOMATIC FAILURE):
- ❌ Skip re-integration and claim "fixes are done"
- ❌ Manually copy fixes to integration branch
- ❌ Proceed to validation with broken integration branch
- ❌ Cherry-pick fixes instead of full re-merge

## Grading Impact

- **-50%** for R232 violation (not processing TODOs)
- **-30%** if skipping code reviews for fixes
- **-40%** if not re-running integration after fixes
- **+20%** for comprehensive fix tracking
- **+15%** for efficient review spawning