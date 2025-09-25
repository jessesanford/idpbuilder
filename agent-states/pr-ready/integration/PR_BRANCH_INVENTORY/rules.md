# PR_BRANCH_INVENTORY State Rules

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