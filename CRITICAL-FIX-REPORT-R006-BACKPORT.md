# CRITICAL FIX REPORT - R006 BACKPORT VIOLATIONS

## Issue Detected
**Date**: 2025-09-05
**Severity**: BLOCKING - Immediate Failure Risk
**Rule**: R006 - Orchestrator NEVER Writes Code

### Violation Observed
The orchestrator in BACKPORT_FIXES state was:
1. About to edit code files directly
2. Making excuses: "the backport process is simpler since everything is integrated"
3. Rationalizing not spawning SW Engineers
4. Planning to use git cherry-pick directly

### Root Cause
The BACKPORT_FIXES state rules were not explicit enough about:
- R006 enforcement during backporting
- Invalid excuses orchestrators make
- Mandatory delegation to SW Engineers

## Fixes Applied

### 1. BACKPORT_FIXES State Rules Enhanced
**File**: `/agent-states/orchestrator/BACKPORT_FIXES/rules.md`

Added sections:
- **CRITICAL WARNING section** at top alerting to R006 monitoring
- **Explicit list of detected violations** orchestrators attempt
- **Forbidden excuses** that indicate rule violation
- **Mandatory process** requiring SW Engineer spawning
- **Detection triggers** for violations

Key changes:
- Orchestrator MUST spawn SW Engineers for EACH effort
- Orchestrator NEVER touches code during backporting
- All excuses about "simpler integration" are INVALID
- Violation = -100% immediate failure

### 2. R006 Rule Strengthened
**File**: `/rule-library/R006-orchestrator-never-writes-code.md`

Added prohibitions:
- ❌ BACKPORT fixes to effort branches
- ❌ CHERRY-PICK commits between branches  
- ❌ APPLY patches or fixes directly

Added violation examples:
- Backporting fixes directly
- Making excuses to avoid delegation
- List of invalid rationalizations

## Verification Steps

### To verify the fix is working:
```bash
# 1. Check BACKPORT_FIXES rules include R006 enforcement
grep -A 10 "R006 ENFORCEMENT" /home/vscode/software-factory-template/agent-states/orchestrator/BACKPORT_FIXES/rules.md

# 2. Verify forbidden excuses are listed
grep -A 5 "NO EXCUSES ACCEPTED" /home/vscode/software-factory-template/agent-states/orchestrator/BACKPORT_FIXES/rules.md

# 3. Confirm R006 mentions backporting
grep -A 3 "BACKPORT" /home/vscode/software-factory-template/rule-library/R006-orchestrator-never-writes-code.md
```

## Impact Assessment

### Without These Fixes:
- Orchestrator would violate R006 = IMMEDIATE FAILURE
- PRs would not have backported fixes
- Individual effort branches would fail to merge
- Entire integration would be unusable

### With These Fixes:
- Orchestrator forced to spawn SW Engineers
- Each effort gets proper backporting
- All PRs independently mergeable
- Integration succeeds properly

## Recommendations

1. **Monitor orchestrator behavior** in BACKPORT_FIXES state
2. **Watch for new excuses** orchestrators invent
3. **Ensure SW Engineers** are always spawned for code work
4. **Never allow** orchestrator to rationalize violations

## Status
✅ **FIXED AND DEPLOYED**
- Committed: 8430c95
- Pushed to: main branch
- Effective: Immediately

---
**Factory Manager**: software-factory-manager
**Timestamp**: 2025-09-05 05:44:21 UTC