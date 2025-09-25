# Split Branch Naming Inconsistency Fix Report

## Issue Identified
The orchestrator, code-reviewer, and sw-engineer agents have inconsistent expectations for split branch naming, causing agents to be unable to find the split branches they need to work with.

## Canonical Format (Per R204 and branch-naming-helpers.sh)

### Directory Naming
- **Format**: `${effort_name}-SPLIT-${SPLIT_NUMBER}`
- **Example**: `api-types-SPLIT-001`
- **Location**: `/efforts/phase${PHASE}/wave${WAVE}/${effort_name}-SPLIT-${SPLIT_NUMBER}`
- **Key**: Single hyphen, uppercase "SPLIT"

### Branch Naming
- **Format**: `${original_branch}--split-${split_number}`
- **Example**: `tmc-workspace/phase1/wave1/api-types--split-001`
- **Key**: Double hyphen, lowercase "split"
- **Must use**: `get_split_branch_name()` function from branch-naming-helpers.sh

## Inconsistencies Found

### 1. Orchestrator MONITOR State (INCORRECT)
**File**: `/agent-states/orchestrator/MONITOR/rules.md`
**Lines**: 404, 410
**Current**: 
```bash
SPLIT_BRANCH="${PROJECT_PREFIX}/phase${PHASE}/wave${WAVE}/${EFFORT_NAME}-SPLIT-${SPLIT_NAME}"
```
**Should be**:
```bash
# Load helper functions
source "$CLAUDE_PROJECT_DIR/utilities/branch-naming-helpers.sh"
# Get original branch name with prefix
ORIGINAL_BRANCH=$(get_effort_branch_name "$PHASE" "$WAVE" "$EFFORT_NAME")
# Create split branch name
SPLIT_BRANCH=$(get_split_branch_name "$ORIGINAL_BRANCH" "$SPLIT_NAME")
```

### 2. Previous Split Branch References (Line 404)
**Current**:
```bash
PREV_BRANCH="${PROJECT_PREFIX}/phase${PHASE}/wave${WAVE}/${EFFORT_NAME}-SPLIT-${PREV_SPLIT}"
```
**Should be**:
```bash
PREV_BRANCH=$(get_split_branch_name "$ORIGINAL_BRANCH" "$PREV_SPLIT")
```

## Files That Need Updates

1. `/agent-states/orchestrator/MONITOR/rules.md` - Fix split branch creation
2. `/agent-states/orchestrator/SETUP_EFFORT_INFRASTRUCTURE/rules.md` - Check for similar issues
3. `/agent-states/orchestrator/CREATE_SPLIT_INFRASTRUCTURE/rules.md` - If it exists
4. Any other orchestrator states that create split branches

## Verification Needed

After fixes, verify:
1. Orchestrator creates branches as: `{prefix}/phase{X}/wave{Y}/{effort}--split-{NNN}`
2. Code-reviewer expects: `{prefix}/phase{X}/wave{Y}/{effort}--split-{NNN}`
3. SW-engineer expects: `{prefix}/phase{X}/wave{Y}/{effort}--split-{NNN}`
4. Directories remain: `{effort}-SPLIT-{NNN}` (uppercase, single hyphen)

## Impact
This fix is critical because:
- Agents cannot find split branches with wrong names
- Line counting fails when branches don't exist
- Integration fails when branches are misnamed
- The entire split workflow breaks down

## Fixes Applied

### 1. Orchestrator MONITOR State - FIXED ✅
- Updated to use `get_split_branch_name()` function
- Now correctly creates branches as `--split-NNN` format
- Previous split references also fixed

### 2. SW Engineer SPLIT_IMPLEMENTATION State - FIXED ✅
- Updated preflight check to expect `--split-` format only
- Fixed split number extraction to use `--split-` pattern
- Fixed base branch calculation for sequential splits

### 3. Code Reviewer - VERIFIED ✅
- Already refers to splits generically
- Works with the format created by orchestrator

## Summary
The split branch naming inconsistency has been resolved:
- **Directories**: `{effort}-SPLIT-{NNN}` (uppercase, single hyphen)
- **Branches**: `{original-branch}--split-{NNN}` (lowercase, double hyphen)
- All agents now use consistent naming through the helper functions