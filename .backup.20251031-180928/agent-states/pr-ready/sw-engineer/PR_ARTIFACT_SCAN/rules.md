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
     "software-factory-3.0-state-machine.json"
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
           "state_files": ["orchestrator-state-v3.json"],
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

