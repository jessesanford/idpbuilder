# 🚨🚨🚨 RULE R263: Integration Documentation Requirements 🚨🚨🚨

## Rule Definition
**Criticality:** BLOCKING
**Category:** Agent-Specific
**Applies To:** integration-agent

## MANDATORY DOCUMENTATION

### 1. Work Log Requirements (work-log.md)
**METICULOUS TRACKING OF EVERY OPERATION:**

The work log MUST be:
- **REPLAYABLE** - Anyone can reproduce exact same output
- **COMPLETE** - Every git command documented
- **TIMESTAMPED** - Each operation time recorded
- **ANNOTATED** - Reasoning for each decision

#### Work Log Format
```markdown
# Integration Work Log
Start Time: YYYY-MM-DD HH:MM:SS

## Operation 1: [timestamp]
Command: git checkout -b integration-branch main
Result: Switched to new branch 'integration-branch'
Reason: Creating clean integration branch from main

## Operation 2: [timestamp]
Command: git merge feature-1 --no-ff
Result: Merge successful, no conflicts
Files Changed: 5 files, +150 lines, -20 lines

## Operation 3: [timestamp]
Command: git merge feature-2 --no-ff
Result: CONFLICT in src/api/handler.go
Resolution: Kept feature-2 version (newer implementation)
Command: git add src/api/handler.go
Command: git commit -m "resolve: conflict in handler.go"
```

### 2. Integration Report Requirements (INTEGRATE_WAVE_EFFORTS-REPORT.md)
**COMPREHENSIVE FINAL REPORT:**

#### Required Sections
```markdown
# INTEGRATE_WAVE_EFFORTS REPORT
Date: YYYY-MM-DD
Integration Branch: branch-name
Target: main (or specified)

## 1. OVERVIEW
### Branches Integrated
- [ ] feature-1 (SHA: abc123)
- [ ] feature-2 (SHA: def456)
- [ ] bugfix-1 (SHA: ghi789)

### Integration Statistics
- Total Commits: X
- Files Changed: Y
- Lines Added: +Z
- Lines Removed: -W
- Conflicts Resolved: N

## 2. ERRORS AND ISSUES FOUND
### Build Errors
- Error 1: Description
  - File: path/to/file.go
  - Line: 123
  - Type: Compilation/Test/Lint
  - **STATUS: NOT FIXED** (upstream bug)

### Test Failures
- Test: TestXYZ
  - Reason: [explanation]
  - **STATUS: NOT FIXED** (document only)

## 3. COMPENSATING/REMEDIATION RECOMMENDATIONS
- Recommendation 1: [specific action needed]
- Recommendation 2: [specific fix required]
- Recommendation 3: [architectural consideration]

## 4. BUILD AND TEST RESULTS
### Build Status
- Command: make build
- Result: PROJECT_DONE/FAILURE
- Duration: Xs
- Output: [relevant output]

### Test Results  
- Command: make test
- Passed: X/Y
- Failed Tests: [list]
- Coverage: X%

## 5. UPSTREAM BUGS IDENTIFIED
**DO NOT FIX - DOCUMENTATION ONLY**
- Bug 1: [description]
  - Location: file:line
  - Impact: HIGH/MEDIUM/LOW
  - Recommended Fix: [suggestion]

## 6. INTEGRATE_WAVE_EFFORTS VERIFICATION
- [ ] All branches merged successfully
- [ ] All conflicts resolved
- [ ] Work log is complete and replayable
- [ ] No original branches were modified
- [ ] No cherry-picks were used

## 7. FINAL STATE
- Integration Branch: branch-name
- Ready for Review: YES/NO
- Blocking Issues: [list or NONE]
```

### 3. Commit Requirements
The integration report MUST be:
- Committed to the integration branch
- Pushed to remote
- Last commit before completion

```bash
# Final commit sequence
git add INTEGRATE_WAVE_EFFORTS-REPORT.md work-log.md
git commit -m "docs: complete integration report and work log"
git push origin integration-branch
```

## Enforcement

```bash
# Verify documentation completeness
verify_documentation() {
    local integration_dir="$1"
    
    # Check work-log exists and is complete
    if [[ ! -f "$integration_dir/work-log.md" ]]; then
        echo "❌ BLOCKING: Missing work-log.md!"
        return 1
    fi
    
    # Check integration report exists
    if [[ ! -f "$integration_dir/INTEGRATE_WAVE_EFFORTS-REPORT.md" ]]; then
        echo "❌ BLOCKING: Missing INTEGRATE_WAVE_EFFORTS-REPORT.md!"
        return 1
    fi
    
    # Verify work-log has commands
    local command_count=$(grep -c "Command:" "$integration_dir/work-log.md")
    if [[ $command_count -lt 3 ]]; then
        echo "❌ Work log appears incomplete (only $command_count commands)"
    fi
    
    # Verify report has all sections
    for section in "OVERVIEW" "ERRORS" "BUILD" "UPSTREAM"; do
        grep -q "$section" "$integration_dir/INTEGRATE_WAVE_EFFORTS-REPORT.md" || \
            echo "❌ Missing section: $section"
    done
}
```

## Grading Impact
- 25% - Work log completeness and replayability
- 25% - Integration report comprehensiveness
- FAILURE if documentation not committed to branch

## Related Rules
- R264 - Work Log Tracking Requirements
- R265 - Integration Testing Requirements
- R266 - Upstream Bug Documentation