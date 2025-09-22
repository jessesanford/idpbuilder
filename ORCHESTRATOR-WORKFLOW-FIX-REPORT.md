# Orchestrator Workflow Critical Fix Report

## Date: 2025-09-07
## Manager: software-factory-manager

## 🔴🔴🔴 CRITICAL ISSUES ADDRESSED 🔴🔴🔴

### Issue 1: Orchestrator Attempting to Run line-counter.sh
**Violation**: R319 - Orchestrator NEVER measures code
**Status**: ✅ FIXED
**Solution**: 
- Clarified in MONITOR_IMPLEMENTATION that orchestrator must NEVER run line-counter.sh
- Emphasized that ONLY Code Reviewers perform size measurements
- Added explicit warnings against orchestrator measuring code

### Issue 2: Line Counter Failing on Split Branches
**Error**: "Could not determine base branch for 'idpbuilder-oci-build-push/phase1/wave2/cert-validation-split-001'"
**Status**: ✅ FIXED
**Solution**:
- Updated line-counter.sh regex to support both `--split-` and `-split-` patterns
- Fixed pattern matching to handle alternative naming formats
- Ensured split delimiter is preserved when determining previous split base

### Issue 3: Orchestrator Skipping Code Review Between Splits
**Violation**: Each split must get full review cycle
**Status**: ✅ FIXED
**Solution**:
- Added explicit workflow in MONITOR_IMPLEMENTATION state rules
- Clarified that after split implementation: MUST spawn Code Reviewer
- Enhanced CREATE_NEXT_SPLIT_INFRASTRUCTURE to only happen AFTER review passes
- Added validation that next split infrastructure only created if review passes AND more splits needed

### Issue 4: R304 Conflicts with R319 in Orchestrator States
**Violation**: R304 (use line-counter) conflicts with R319 (orchestrator never measures)
**Status**: ✅ FIXED
**Solution**:
- Removed R304 from ALL 40 orchestrator state files
- Created fix-orchestrator-r304.sh script for bulk removal
- Ensured R319 remains prominently enforced

## 📊 Files Modified

### Core Files Updated:
1. **tools/line-counter.sh**
   - Lines 309, 311-312: Added support for `-split-` pattern
   - Line 333: Use same delimiter format for previous split

2. **agent-states/orchestrator/MONITOR_IMPLEMENTATION/rules.md**
   - Added split implementation special case section
   - Clarified review requirement between splits
   - Enhanced R319 enforcement section

3. **agent-states/orchestrator/CREATE_NEXT_SPLIT_INFRASTRUCTURE/rules.md**
   - Added critical context about when to enter this state
   - Clarified valid paths (only after review)
   - Enhanced key points with review cycle requirement

4. **40 Orchestrator State Files**
   - Removed R304 references from all files
   - Maintained R319 enforcement

## 🚨 Critical Workflow Clarification

### CORRECT Split Workflow:
```
1. Implement Split-001
2. Spawn Code Reviewer for Split-001
3. Review passes → Create Split-002 infrastructure
4. Implement Split-002
5. Spawn Code Reviewer for Split-002
6. Review passes → Create Split-003 infrastructure
7. Continue...
```

### WRONG Workflow (What Was Happening):
```
1. Implement Split-001
2. Immediately create Split-002 infrastructure ❌
3. Implement Split-002
4. Immediately create Split-003 infrastructure ❌
5. Missing reviews! ❌
```

## 📋 Rule Compliance Summary

### Rules Enforced:
- ✅ **R319**: Orchestrator NEVER measures code
- ✅ **R007**: Each split gets full review cycle
- ✅ **R199**: Single reviewer for split planning
- ✅ **R202**: Sequential split implementation
- ✅ **R204**: Orchestrator creates split infrastructure
- ✅ **R290**: State rules read before state actions

### Rules Removed from Orchestrator:
- ❌ **R304**: Mandatory line counter (conflicts with R319)

## 🔍 Verification Steps

### To Verify Fixes:
1. Check orchestrator doesn't attempt line counting:
   ```bash
   grep -r "line-counter.sh" agent-states/orchestrator/*/rules.md
   ```
   Should only find references saying NOT to use it.

2. Test split branch pattern:
   ```bash
   ./tools/line-counter.sh phase1/wave2/cert-validation-split-001
   ```
   Should correctly identify base as cert-validation branch.

3. Verify R304 removal:
   ```bash
   grep -r "R304" agent-states/orchestrator/*/rules.md
   ```
   Should return no results.

## 🎯 Impact

These fixes ensure:
1. **Orchestrator compliance**: No more R319 violations
2. **Split quality**: Every split gets reviewed before proceeding
3. **Line counting works**: Handles various split naming conventions
4. **Clear workflow**: No ambiguity about review requirements

## 📝 Commit Information

**Commit Hash**: cd15157
**Branch**: orchestrator-rules-to-state-rules
**Files Changed**: 44 files
**Insertions**: 84
**Deletions**: 169

## ✅ All Critical Issues Resolved

The orchestrator workflow now properly:
- Delegates ALL measurements to Code Reviewers
- Enforces review between EVERY split
- Supports flexible split branch naming
- Maintains strict R319 compliance