# PR_ARTIFACT_SCAN State Rules

## 🔴🔴🔴 STATE PURPOSE: Scan Branches for SF Artifacts 🔴🔴🔴

### 🚨🚨🚨 RULE R365 - PR Artifact Detection and Inventory Protocol [BLOCKING]
**Source:** `$CLAUDE_PROJECT_DIR/rule-library/R365-pr-artifact-detection-protocol.md`
**Criticality:** BLOCKING - Undetected artifacts = PR contamination

This state implements R365 requirements for comprehensive artifact detection.

### MANDATORY ACTIONS (R233 + R365 COMPLIANT)
Upon entering this state, IMMEDIATELY:

1. **Define artifact patterns**
   ```bash
   # Software Factory artifact patterns
   SF_DIRS=(
     "todos"
     "efforts"
     "agent-states"
     "rule-library"
     "templates"
     "utilities"
     "phase-plans"
     "wave-plans"
     "protocols"
     ".claude/agents"
     ".claude/commands"
   )

   SF_FILES=(
     "*-state.json"
     "*.todo"
     "CODE-REVIEW-REPORT.md"
     "SPLIT-PLAN*.md"
     "PROJECT-IMPLEMENTATION-PLAN.md"
     "EFFORT-IMPLEMENTATION-PLAN.md"
     "SOFTWARE-FACTORY-STATE-MACHINE.md"
     "RECOVERY-*.md"
     "CURRENT-TODO-STATE.md"
     "CRITICAL-*.md"
     "FINAL-*.md"
   )
   ```

2. **Scan each assigned branch**
   ```bash
   for BRANCH in $ASSIGNED_BRANCHES; do
     echo "Scanning branch: $BRANCH"
     git checkout $BRANCH

     # Count directories
     for dir in "${SF_DIRS[@]}"; do
       if [ -d "$dir" ]; then
         FILE_COUNT=$(find "$dir" -type f | wc -l)
         echo "  Found: $dir/ ($FILE_COUNT files)"
       fi
     done

     # Count files
     for pattern in "${SF_FILES[@]}"; do
       FILES=$(ls $pattern 2>/dev/null)
       if [ ! -z "$FILES" ]; then
         echo "  Found: $pattern files"
       fi
     done
   done
   ```

3. **Create artifact inventory**
   ```json
   {
     "artifact_scan_results": {
       "branch1": {
         "directories": {
           "todos": 45,
           "efforts": 12,
           "agent-states": 78
         },
         "files": {
           "state_files": ["orchestrator-state.json"],
           "todo_files": ["file1.todo", "file2.todo"],
           "planning_docs": ["PROJECT-IMPLEMENTATION-PLAN.md"]
         },
         "total_artifacts": 156,
         "size_mb": 2.3
       }
     }
   }
   ```

4. **Identify high-risk contamination**
   ```bash
   # Check for deeply embedded artifacts
   find . -name "*.todo" -o -name "*-state.json" | head -20

   # Check for SF references in code
   grep -r "SOFTWARE-FACTORY" --include="*.go" --include="*.js"
   ```

5. **Generate scan report**
   ```markdown
   # SF Artifact Scan Report

   ## Summary
   - Branches scanned: 16
   - Total artifacts: 2,345
   - Contaminated branches: 16/16

   ## Per-Branch Breakdown
   [Detailed breakdown by branch]

   ## High-Risk Findings
   - Branch8: Contains 90,000 line deletion
   - Branch12: SF references in source code
   ```

### EXIT CRITERIA
✅ All branches scanned
✅ Artifacts inventoried
✅ High-risk items flagged
✅ Report generated

### OUTPUT FILES
- `PR-ARTIFACT-INVENTORY.json`
- `PR-ARTIFACT-SCAN-REPORT.md`
- `PR-HIGH-RISK-BRANCHES.txt`

### PROHIBITED ACTIONS
❌ Do NOT delete any artifacts
❌ Do NOT modify files
❌ Do NOT make commits
❌ Do NOT push changes

### SCAN REQUIREMENTS
- Check every directory
- Use exact pattern matching
- Document file counts
- Calculate storage impact