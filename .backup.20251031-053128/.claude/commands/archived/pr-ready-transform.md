---
name: pr-ready-transform
description: Continue PR-ready transformation of effort branches into clean upstream PRs. Removes all Software Factory artifacts, consolidates commits, and prepares branches for production pull requests.
model: opus
---

# PR-READY TRANSFORMATION CONTINUATION

## 🎯 PURPOSE
Transform Software Factory effort branches into clean, production-ready branches suitable for upstream pull requests.

## 📊 PROCESS OVERVIEW
The PR-ready transformation follows 7 phases:
1. **Discovery** - Find all SF artifacts
2. **Cleanup** - Remove all SF metadata
3. **Consolidation** - Squash commits
4. **Verification** - Ensure core files intact
5. **Rebase** - Establish proper sequence
6. **Validation** - Test merge compatibility
7. **Preparation** - Create PR documentation

## 🚀 STARTUP SEQUENCE

1. **Check current PR state**
   ```bash
   cat pr-ready-state.json 2>/dev/null || echo "No PR state found"
   ```

2. **Determine current phase**
   - If no state → Start with `PR_READY_INIT`
   - If state exists → Continue from `current_state`

3. **Load state machine**
   ```bash
   cat SOFTWARE-FACTORY-PR-READY-STATE-MACHINE.md
   ```

4. **Resume as orchestrator**
   ```bash
   /spawn agent-orchestrator [CURRENT_PR_STATE]
   ```

## 📋 COMMON CONTINUATIONS

### Starting Fresh
```bash
# Initialize PR transformation
/spawn agent-orchestrator PR_READY_INIT \
  --branches "all-effort-branches" \
  --upstream "upstream/main"
```

### After Discovery
```bash
# Review discovery results and plan cleanup
/spawn agent-orchestrator PR_CLEANUP_PLANNING
```

### After Cleanup
```bash
# Verify cleanup and consolidate commits
/spawn agent-orchestrator PR_CONSOLIDATION_PLANNING
```

### After Issues Found
```bash
# Fix destructive branches
/spawn agent-orchestrator PR_FIX_DESTRUCTIVE_CHANGES \
  --branches "branch8,branch12"
```

### Final Steps
```bash
# Complete transformation
/spawn agent-orchestrator PR_FINAL_PREPARATION
```

## ⚠️ CRITICAL REQUIREMENTS

1. **ZERO SF ARTIFACTS** - Not a single SF file may remain
2. **CORE INTEGRITY** - No application files deleted
3. **CLEAN HISTORY** - Commits properly consolidated
4. **MERGE READY** - Branches test clean to upstream
5. **BUILD VERIFIED** - Code compiles and tests pass

## 🔍 VALIDATION TOOLS

```bash
# Check for SF artifacts
./tools/pr-ready/detect-sf-artifacts.sh

# Verify core files preserved
./tools/pr-ready/verify-core-files.sh

# Test merge compatibility
./tools/pr-ready/test-merge-compatibility.sh
```

## 📊 PROGRESS TRACKING

Monitor transformation progress:
```bash
# Check current state
jq '.current_state' pr-ready-state.json

# View discovered artifacts
jq '.artifacts_discovered' pr-ready-state.json

# Check cleanup status
jq '.cleanup_status' pr-ready-state.json

# Review validation results
jq '.validation_results' pr-ready-state.json
```

## 🚨 ERROR RECOVERY

If transformation fails:
1. Check error state: `jq '.errors' pr-ready-state.json`
2. Review abort reason: `cat PR-ABORT-REASON.md`
3. Fix identified issues
4. Resume from recovery state

## ✅ PROJECT_DONE INDICATORS

Transformation complete when:
- All branches cleaned (0 SF artifacts)
- All branches consolidated (1 commit each)
- All branches rebased (proper sequence)
- All merges tested (conflicts documented)
- Documentation created (PR-MERGE-GUIDE.md)

---

**Command Type**: Continuation
**Agent**: Orchestrator
**State Machine**: SOFTWARE-FACTORY-PR-READY-STATE-MACHINE.md