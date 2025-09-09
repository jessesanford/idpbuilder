# WAITING_FOR_PROJECT_FIX_PLANS State Rules

## State Purpose
Actively monitor Code Reviewer creating fix plans for project-level integration bugs documented per R266. These are bugs found during final project integration that need to be fixed in source branches.

## Critical Rules

### 🔴🔴🔴 RULE R322: MANDATORY STOP BEFORE STATE TRANSITION (SUPREME LAW)
- **STOP** and save state before ANY transition
- **READ** orchestrator-state.json to verify current state
- **VALIDATE** next state exists in SOFTWARE-FACTORY-STATE-MACHINE.md
- **VIOLATION = IMMEDIATE FAILURE**

### 🔴🔴🔴 RULE R233: IMMEDIATE ACTION REQUIRED (SUPREME LAW)
- **NO PASSIVE WAITING** - Must actively check for completion
- **IMMEDIATE ACTION** - Start checking within first response
- **CONTINUOUS MONITORING** - Check every 30-60 seconds
- **States are VERBS** - "WAITING" means "ACTIVELY CHECKING"

### 🔴🔴🔴 RULE R321: IMMEDIATE BACKPORT REQUIRED (SUPREME LAW)
- Project fixes MUST be applied to source branches
- Integration branches are READ-ONLY for code
- Fix plans must target phase/wave/effort branches
- After fixes, project integration MUST be re-run

### 🚨🚨🚨 RULE R290: STATE RULE VERIFICATION (BLOCKING)
- **MUST** verify this rules file exists and is loaded
- **MUST** acknowledge all rules before proceeding
- **MUST** validate state transitions against state machine

### 🚨🚨🚨 RULE R232: MONITOR STATE REQUIREMENTS (BLOCKING)
- **MUST** check TodoWrite for pending items BEFORE transition
- **MUST** process ALL pending items immediately
- **NO** "I will..." statements - only "I am..." with action
- **VIOLATION = AUTOMATIC FAILURE**

### ⚠️⚠️⚠️ RULE R266: BUG DOCUMENTATION REQUIREMENTS (WARNING)
- Fix plans MUST reference bugs documented in PROJECT-INTEGRATION-REPORT.md
- Each bug from R266 documentation MUST have corresponding fix
- Plans MUST maintain bug numbering from original documentation
- Plans MUST specify exact source branch for each fix

### ⚠️⚠️⚠️ RULE R287: TODO PERSISTENCE (WARNING)
- **MUST** save TODOs every 10 messages or 15 minutes
- **MUST** save before state transition
- **MUST** commit and push TODO state

## Required Actions

1. **Initial Check (IMMEDIATE)**
   ```bash
   # Verify Code Reviewer was spawned for fix planning
   grep "spawned_agents" orchestrator-state.json | grep -i "project_fix"
   
   # Check project integration directory
   ls -la project-integration/
   
   # Look for bug documentation per R266
   grep "UPSTREAM BUGS IDENTIFIED" project-integration/PROJECT-INTEGRATION-REPORT.md
   ```

2. **Active Monitoring Loop**
   ```bash
   # Monitor for fix plan creation
   while true; do
     # Check for fix plan file
     if [ -f "project-integration/PROJECT-FIX-PLAN.md" ]; then
       echo "✓ Project fix plan detected at $(date)"
       break
     fi
     
     # Check Code Reviewer status if available
     if [ -f "project-integration/fix-planning-status.yaml" ]; then
       STATUS=$(grep "status:" project-integration/fix-planning-status.yaml | awk '{print $2}')
       echo "Fix planning status: $STATUS at $(date)"
       
       if [ "$STATUS" = "COMPLETED" ]; then
         echo "✓ Fix planning completed"
         break
       elif [ "$STATUS" = "BLOCKED" ]; then
         echo "✗ Fix planning blocked - need to handle"
         # May need to transition to ERROR_RECOVERY
         break
       fi
     fi
     
     echo "Waiting for project fix plan... checking again in 30s"
     sleep 30
   done
   ```

3. **Validate Fix Plan**
   ```bash
   PLAN="project-integration/PROJECT-FIX-PLAN.md"
   
   if [ -f "$PLAN" ]; then
     echo "Validating project fix plan..."
     
     # Check required sections
     grep -q "## Bug Summary" "$PLAN" || echo "✗ Missing bug summary"
     grep -q "## Fix Strategy" "$PLAN" || echo "✗ Missing fix strategy"
     grep -q "### Bug #" "$PLAN" || echo "✗ Missing individual bug fixes"
     
     # Verify R266 bug references
     BUG_COUNT=$(grep "### Bug #" project-integration/PROJECT-INTEGRATION-REPORT.md | wc -l)
     FIX_COUNT=$(grep "#### Bug #" "$PLAN" | wc -l)
     
     if [ "$BUG_COUNT" != "$FIX_COUNT" ]; then
       echo "✗ ERROR: Bug count mismatch! R266 docs: $BUG_COUNT, Fix plan: $FIX_COUNT"
     fi
     
     # Verify targets source branches (R321)
     if grep -q "project-integration" "$PLAN" | grep -v "after\|from"; then
       echo "✗ ERROR: Fix plan targets integration branch (violates R321)"
     fi
     
     # Check for parallelization strategy
     grep -q "Parallel Fix Group" "$PLAN" && echo "✓ Has parallel fix strategy"
     grep -q "Sequential Fix" "$PLAN" && echo "✓ Has sequential fix handling"
     
     # Verify SW Engineer spawn instructions
     grep -q "## SW Engineer Spawn Instructions" "$PLAN" || \
       echo "⚠ Missing spawn instructions section"
   else
     echo "✗ PROJECT-FIX-PLAN.md not found yet"
   fi
   ```

4. **Extract Fix Information**
   ```bash
   # Parse fix plan for orchestrator use
   if [ -f "$PLAN" ]; then
     # Count fixes by type
     PARALLEL_FIXES=$(grep -c "### Parallel Fix Group" "$PLAN")
     SEQUENTIAL_FIXES=$(grep -c "### Sequential Fix" "$PLAN")
     
     # Extract engineer assignments
     ENGINEERS_NEEDED=$(grep "Engineer [0-9]:" "$PLAN" | wc -l)
     
     echo "Fix Plan Summary:"
     echo "- Parallel fix groups: $PARALLEL_FIXES"
     echo "- Sequential fixes: $SEQUENTIAL_FIXES"
     echo "- Engineers needed: $ENGINEERS_NEEDED"
     
     # Update state with plan information
     jq ".project_fixes.plan_received = true" -i orchestrator-state.json
     jq ".project_fixes.parallel_groups = $PARALLEL_FIXES" -i orchestrator-state.json
     jq ".project_fixes.engineers_needed = $ENGINEERS_NEEDED" -i orchestrator-state.json
   fi
   ```

5. **Check for Timeout**
   ```bash
   # Get spawn time
   SPAWN_TIME=$(grep "project_fix_plan" orchestrator-state.json -A 2 | \
                grep "timestamp" | tail -1 | cut -d'"' -f2)
   
   # Calculate elapsed
   ELAPSED=$(($(date +%s) - $(date -d "$SPAWN_TIME" +%s)))
   
   # Timeout after 30 minutes (project fixes should be straightforward)
   if [ $ELAPSED -gt 1800 ]; then
     echo "✗ Timeout waiting for project fix plans"
     # Transition to SPAWN_SW_ENGINEER_PROJECT_FIXES with error handling
   fi
   ```

6. **Prepare for Next State**
   ```bash
   # Once plan is ready and validated
   if [ -f "project-integration/PROJECT-FIX-PLAN.md" ]; then
     echo "✅ Project fix plan ready"
     
     # Update state for transition
     jq '.current_state = "SPAWN_SW_ENGINEER_PROJECT_FIXES"' -i orchestrator-state.json
     jq '.project_fixes.plan_location = "project-integration/PROJECT-FIX-PLAN.md"' \
        -i orchestrator-state.json
     
     git add orchestrator-state.json
     git commit -m "state: project fix plan received, ready to spawn engineers"
     git push
     
     echo "Ready to transition to SPAWN_SW_ENGINEER_PROJECT_FIXES"
   fi
   ```

## Transition Rules

### Valid Next States
- **SPAWN_SW_ENGINEER_PROJECT_FIXES** - Fix plan ready, spawn engineers
- **ERROR_RECOVERY** - Planning failed or blocked

### Invalid Transitions
- ❌ Direct to MONITORING_PROJECT_FIXES (must spawn engineers first)
- ❌ Back to PROJECT_INTEGRATION without fixes
- ❌ Skipping fix implementation

## Common Violations to Avoid

1. **Passive waiting** - Violates R233, must actively monitor
2. **Not verifying R266 compliance** - Plan doesn't match documented bugs
3. **Not checking R321 compliance** - Plans target integration branch
4. **Missing timeout handling** - Waiting forever
5. **Not validating plan completeness** - Missing bugs or instructions
6. **Forgetting TODO persistence** - Violates R287

## Monitoring Pattern

```bash
# CORRECT: Active monitoring with comprehensive checks
echo "Monitoring project fix planning at $(date)"
START_TIME=$(date +%s)

while true; do
  # Check for plan file
  if [ -f "project-integration/PROJECT-FIX-PLAN.md" ]; then
    echo "✓ Project fix plan found!"
    
    # Validate it has all bugs from R266
    R266_BUGS=$(grep -c "### Bug #" project-integration/PROJECT-INTEGRATION-REPORT.md)
    PLAN_BUGS=$(grep -c "#### Bug #" project-integration/PROJECT-FIX-PLAN.md)
    
    if [ "$R266_BUGS" = "$PLAN_BUGS" ]; then
      echo "✓ All $R266_BUGS bugs addressed in plan"
      break
    else
      echo "⚠ Plan incomplete: $PLAN_BUGS/$R266_BUGS bugs"
    fi
  fi
  
  # Check status file
  test -f project-integration/fix-planning-status.yaml && \
    grep "progress:" project-integration/fix-planning-status.yaml
  
  # Check timeout
  ELAPSED=$(($(date +%s) - START_TIME))
  if [ $ELAPSED -gt 1800 ]; then
    echo "✗ Timeout after 30 minutes"
    break
  fi
  
  echo "Check at $(date): No complete plan yet"
  sleep 30
done

# WRONG: Just waiting without checks
sleep 1800  # Wait 30 minutes
ls project-integration/PROJECT-FIX-PLAN.md  # Check once at end
```

## Verification Commands

```bash
# Verify state entry
echo "Entered WAITING_FOR_PROJECT_FIX_PLANS at $(date)"
echo "Project integration found bugs - monitoring fix plan creation"

# Check R266 bug documentation
echo "Bugs documented per R266:"
grep "### Bug #" project-integration/PROJECT-INTEGRATION-REPORT.md | head -5

# Active monitoring with validation
for i in {1..60}; do  # 30 minute timeout
  if [ -f "project-integration/PROJECT-FIX-PLAN.md" ]; then
    echo "✓ Fix plan found after $i checks"
    
    # Validate R321 compliance
    if grep -q "fix.*project-integration" project-integration/PROJECT-FIX-PLAN.md; then
      echo "✗ R321 VIOLATION: Plan suggests fixing integration branch!"
    else
      echo "✓ R321 compliant: All fixes target source branches"
    fi
    
    break
  fi
  echo "Check $i: No plan yet at $(date)"
  sleep 30
done

# Prepare transition
if [ -f "project-integration/PROJECT-FIX-PLAN.md" ]; then
  echo "Transitioning to SPAWN_SW_ENGINEER_PROJECT_FIXES"
fi
```

## R322 Checkpoint Requirements

Before transitioning to SPAWN_SW_ENGINEER_PROJECT_FIXES:
1. **STOP** and display checkpoint message
2. Show the PROJECT-FIX-PLAN.md for user review
3. Summarize fix strategy (parallel vs sequential)
4. List affected source branches
5. Wait for user /continue-orchestrating

```markdown
## 🛑 STATE TRANSITION CHECKPOINT: WAITING_FOR_PROJECT_FIX_PLANS → SPAWN_SW_ENGINEER_PROJECT_FIXES

### ✅ Fix Plan Received:
- Location: project-integration/PROJECT-FIX-PLAN.md
- Bugs addressed: X/Y from R266 documentation
- Parallel fix groups: N
- Engineers needed: M

### 📋 Fix Strategy Summary:
[Show key parts of fix plan]

### ⏸️ STOPPED - Awaiting User Review
Please review the fix plan before spawning engineers.
Use /continue-orchestrating to proceed.
```

## References
- R232: rule-library/R232-monitor-state-requirements.md
- R233: rule-library/R233-all-states-immediate-action.md
- R266: rule-library/R266-upstream-bug-documentation.md
- R287: rule-library/R287-todo-persistence-comprehensive.md
- R290: rule-library/R290-state-rule-verification.md
- R321: rule-library/R321-immediate-backport-during-integration.md
- R322: rule-library/R322-mandatory-stop-before-transition.md