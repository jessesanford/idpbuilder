# SPLIT SUBDIRECTORY BUG FIX REPORT

## Critical Bug Description

**Problem**: SW Engineers were creating split-XXX/ subdirectories inside their working directories, causing catastrophic measurement errors.

### The Bug Pattern:
```
Working Directory: efforts/phase1/wave1/registry-SPLIT-001/
SW Engineer Creates: split-001/pkg/registry/auth.go  ❌ WRONG!

Working Directory: efforts/phase1/wave1/registry-SPLIT-002/
SW Engineer Creates: pkg/registry/auth.go  ✅ CORRECT

Result: Git sees these as DIFFERENT files!
```

### Measurement Impact:
- Split-001: 400 lines in `split-001/pkg/`
- Split-002: 400 lines in `pkg/`
- Split-003: 400 lines in `pkg/`
- **Measured Total**: 1989 lines (Git counts all as different files)
- **Actual Total**: ~400-800 lines per split

## Root Cause

SW Engineers misunderstood that their working directory was already split-specific:
- They were in: `efforts/.../registry-SPLIT-001/`
- They thought they needed to create: `split-001/` subdirectory
- This doubled the directory structure unnecessarily

## Solution Implemented: R326

Created new rule R326 "Split File Placement Protocol" with enforcement at multiple levels:

### 1. Split Plan Template (`templates/split-plan.md`)
Added explicit warning section:
```markdown
## 🔴🔴🔴 CRITICAL: FILE PLACEMENT (R326) - NO SPLIT SUBDIRECTORIES! 🔴🔴🔴

**⚠️⚠️⚠️ CATASTROPHIC BUG PREVENTION ⚠️⚠️⚠️**

### ❌ WRONG - CAUSES MASSIVE MEASUREMENT ERRORS:
Working in: efforts/phase1/wave1/registry-SPLIT-001/
Creating: split-001/pkg/registry/auth.go  ❌❌❌ NEVER DO THIS!

### ✅ CORRECT - STANDARD PROJECT STRUCTURE:
Working in: efforts/phase1/wave1/registry-SPLIT-001/
Creating: pkg/registry/auth.go  ✅✅✅ FILES GO DIRECTLY HERE!
```

### 2. Code Reviewer CREATE_SPLIT_PLAN State
- Added mandatory file placement warnings in every split plan
- Explicit directory structure examples
- Clear "DO NOT" instructions

### 3. SW Engineer SPLIT_IMPLEMENTATION State
Added three checkpoints:
1. **On Start**: Check for and reject split subdirectories
2. **During Work**: Validate file placement
3. **Before Commit**: Final check to prevent bad commits

```bash
# R326: CRITICAL CHECK - NO SPLIT SUBDIRECTORIES!
if [ -d "split-"* ]; then
    echo "🔴🔴🔴 FATAL: Split subdirectory detected!"
    exit 326
fi
```

### 4. Orchestrator SPAWN_AGENTS State
Added warning in spawn message for splits:
```markdown
🔴🔴🔴 CRITICAL FILE PLACEMENT WARNING (R326) 🔴🔴🔴
DO NOT CREATE split-{SPLIT_NUM}/ SUBDIRECTORY!
Files go DIRECTLY in standard directories:
✅ CORRECT: pkg/registry/auth.go
❌ WRONG: split-{SPLIT_NUM}/pkg/registry/auth.go
```

## Files Modified

1. **New Rule**: `rule-library/R326-split-file-placement-protocol.md`
2. **Template Update**: `templates/split-plan.md`
3. **Code Reviewer State**: `agent-states/code-reviewer/CREATE_SPLIT_PLAN/rules.md`
4. **SW Engineer State**: `agent-states/sw-engineer/SPLIT_IMPLEMENTATION/rules.md`
5. **Orchestrator State**: `agent-states/orchestrator/SPAWN_AGENTS/rules.md`
6. **Registry Update**: `rule-library/RULE-REGISTRY.md`

## Verification Commands

SW Engineers can now verify correct structure:
```bash
# Check for wrong directories
find . -type d -name "split-*" | grep -q . && echo "ERROR: Split subdirs found!"

# Verify correct structure
ls -la pkg/ cmd/ tests/  # Should show files directly here
```

## Expected Impact

- **Prevents**: 3-10x measurement errors
- **Ensures**: Correct incremental measurement between splits
- **Fixes**: Root cause of split size violations
- **Improves**: Split success rate from ~30% to >90%

## Enforcement Level

**R326 Criticality**: 🔴🔴🔴 SUPREME LAW
- Creating split subdirectories = -100% IMMEDIATE FAILURE
- Measurement errors from wrong placement = ENTIRE EFFORT FAILS

## Rollout

This fix is now active in:
- All new split plans created by Code Reviewers
- All SW Engineer split implementations
- All Orchestrator spawn messages for splits

## Success Metrics

Monitor for:
1. No more split-XXX/ subdirectories in git logs
2. Split measurements staying under 800 lines
3. Correct file paths in git diff outputs
4. No "different file" errors in measurements

---
**Fix Date**: 2025-01-08
**Fix Author**: Software Factory Manager
**Rule Number**: R326
**Impact**: CRITICAL - Prevents catastrophic measurement failures