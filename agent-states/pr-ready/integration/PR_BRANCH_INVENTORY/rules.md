# PR_BRANCH_INVENTORY State Rules

## 🔴🔴🔴 ABSOLUTE REQUIREMENT: STATE MANAGER CONSULTATION 🔴🔴🔴

**SUPREME LAW - R517 - UNIVERSAL STATE MANAGER CONSULTATION LAW**

**BEFORE EXITING THIS STATE, YOU MUST:**

1. **SPAWN STATE MANAGER FOR SHUTDOWN_CONSULTATION** (MANDATORY - NO EXCEPTIONS)
2. **NEVER UPDATE STATE FILES DIRECTLY** (orchestrator-state-v3.json, bug-tracking.json, etc.)
3. **NEVER COMMIT STATE FILES YOURSELF** (State Manager does this atomically)

**FAILURE TO CONSULT STATE MANAGER = IMMEDIATE SYSTEM HALT (-100% GRADE)**

### Enforcement Mechanism

If you attempt to exit this state without spawning State Manager:
- ❌ Pre-commit hooks will REJECT your commit
- ❌ Validation tools will FAIL the build
- ❌ Grading system will assign -100% penalty
- ❌ System will transition to ERROR_RECOVERY

### Required Pattern (COPY THIS EXACTLY)

```bash
# At end of state work, BEFORE any state file updates:

echo "🔴 MANDATORY: Spawning State Manager for SHUTDOWN_CONSULTATION"

# Spawn State Manager (REQUIRED - NOT OPTIONAL)
# Task: state-manager
# State: SHUTDOWN_CONSULTATION
# Current State: [YOUR_CURRENT_STATE]
# Proposed Next State: [YOUR_PROPOSED_NEXT_STATE]
# Work Summary: [SUMMARY_OF_WORK_COMPLETED]

# State Manager will:
# 1. Validate proposed transition against state machine
# 2. Update all 4 state files atomically
# 3. Commit with [R288] tag
# 4. Return REQUIRED next state (may differ from proposal)

# Wait for State Manager response
# Follow State Manager's directive (REQUIRED next state)
# DO NOT proceed until State Manager confirms
```

**YOU MUST NEVER:**
- ❌ Update orchestrator-state-v3.json yourself
- ❌ Update bug-tracking.json yourself
- ❌ Update integration-containers.json yourself
- ❌ Use `jq` to modify state files
- ❌ Use `sed/awk` to modify state files
- ❌ Set `validated_by: "orchestrator"` (must be "state-manager")
- ❌ Skip State Manager consultation "just this once"
- ❌ Think "I'll validate it manually"

**ONLY State Manager may update state files. This is NON-NEGOTIABLE.**

See: `rule-library/R517-universal-state-manager-consultation-law.md`
## 🔴🔴🔴 STATE PURPOSE: Inventory All Effort Branches 🔴🔴🔴

### MANDATORY ACTIONS (R233 COMPLIANT)
Upon entering this state, IMMEDIATELY:

1. **Enumerate all branches**
   ```bash
   # Get all local branches
   git branch -a | grep -E "(phase|wave|effort|feature)" > branches.txt

   # Get branch creation order
   git for-each-ref --sort=committerdate refs/heads/ \
     --format='%(committerdate:short) %(refname:short)' \
     | grep -E "(phase|wave|effort)" > branch-timeline.txt
   ```

2. **Document branch dependencies**
   ```bash
   # For each branch, find merge base
   for branch in $(cat branches.txt); do
     echo "Branch: $branch"
     git merge-base main $branch
     git log --oneline main..$branch | head -5
   done > dependencies.txt
   ```

3. **Check main branch status**
   ```bash
   # Compare with upstream
   git fetch upstream
   git diff upstream/main..origin/main --stat

   # Check for SF artifacts in main
   git ls-tree -r main | grep -E "(todos/|efforts/|state.json)"
   ```

4. **Create dependency tree**
   ```json
   {
     "branch_inventory": {
       "total_branches": 16,
       "phases": {
         "phase1": {
           "wave1": ["branch1", "branch2"],
           "wave2": ["branch3", "branch4"]
         }
       },
       "dependencies": {
         "branch2": "depends_on:branch1",
         "branch3": "depends_on:branch2"
       },
       "merge_order": ["branch1", "branch2", "branch3", ...]
     }
   }
   ```

5. **Generate inventory report**
   ```markdown
   # PR Branch Inventory Report

   ## Branches Found
   - Total: 16 effort branches
   - Phases: 3
   - Waves: 6

   ## Dependency Analysis
   - Sequential dependencies: 12
   - Parallel opportunities: 4

   ## Main Branch Status
   - Contaminated: Yes/No
   - Divergence: X commits
   ```

### EXIT CRITERIA
✅ All branches inventoried
✅ Dependencies documented
✅ Main status checked
✅ Report generated

### OUTPUT FILES
- `PR-BRANCH-INVENTORY.json`
- `PR-DISCOVERY-REPORT-integration.md`
- `PR-DEPENDENCY-TREE.json`

### PROHIBITED ACTIONS
❌ Do NOT modify any branches
❌ Do NOT delete anything
❌ Do NOT merge or rebase
❌ Do NOT push changes

### ERROR HANDLING
- No branches found → Report and exit
- Upstream unavailable → Continue with local
- Malformed branch names → Document issues

## 🔴🔴🔴 R405 - MANDATORY AUTOMATION CONTINUATION FLAG 🔴🔴🔴

**SUPREME LAW - PENALTY FOR VIOLATION: -100% GRADE**

### YOU MUST OUTPUT THE CONTINUATION FLAG AS YOUR LAST ACTION

**EVERY STATE COMPLETION MUST END WITH EXACTLY ONE OF:**
```bash
# If state completed successfully and factory should continue:
echo "CONTINUE-SOFTWARE-FACTORY=TRUE"

# If error/block/manual review needed:
echo "CONTINUE-SOFTWARE-FACTORY=FALSE"
```

### CRITICAL REQUIREMENTS:
1. **ABSOLUTE LAST OUTPUT**: This MUST be the very last line of output before state completion
2. **NO TEXT AFTER**: No explanations, summaries, or any text after the flag
3. **EXACTLY AS SHOWN**: Use exact format - no variations like CONTINUE-ORCHESTRATING
4. **ALWAYS REQUIRED**: Every single state must output this flag
5. **GREPPABLE**: Must be on its own line for automation parsing

### WHEN TO USE TRUE:
- ✅ State work completed successfully
- ✅ All validations passed
- ✅ Ready for next state
- ✅ No blockers detected
- ✅ All requirements met

### WHEN TO USE FALSE:
- ❌ Any unrecoverable error occurred
- ❌ Manual intervention required
- ❌ Missing required files or configs
- ❌ Test failures blocking progress
- ❌ Ambiguous or unclear instructions
- ❌ Wrong working directory or branch
- ❌ State machine validation failed

### IMPLEMENTATION PATTERN:
```bash
# At the VERY END of state execution, after ALL work:

# Determine success/failure
if [[ "$STATE_COMPLETED_PROJECT_DONEFULLY" == "true" ]]; then
    echo "✅ State work completed successfully"
    echo "CONTINUE-SOFTWARE-FACTORY=TRUE"
else
    echo "❌ State encountered issues requiring intervention"
    echo "CONTINUE-SOFTWARE-FACTORY=FALSE"
fi

# DO NOT OUTPUT ANYTHING AFTER THE FLAG!
```

### GRADING IMPACT:
- **Missing flag**: -100% AUTOMATIC FAILURE
- **Text after flag**: -50% penalty
- **Wrong format**: -100% AUTOMATIC FAILURE
- **Multiple flags**: -50% penalty

**See: $CLAUDE_PROJECT_DIR/rule-library/R405-automation-continuation-flag.md**

