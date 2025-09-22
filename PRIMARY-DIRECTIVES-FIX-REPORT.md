# PRIMARY DIRECTIVES Fix Report

## Problem Identified

The orchestrator reported: "I didn't see a clear PRIMARY DIRECTIVES section listing specific rule files" when entering SPAWN_CODE_REVIEWER_MERGE_PLAN state.

Investigation revealed that PRIMARY DIRECTIVES sections were buried deep in the state files (typically around line 94) after multiple other sections like R322 and R290 enforcement blocks.

## Root Cause

The PRIMARY DIRECTIVES section was present but not easily discoverable because:
1. It appeared after 90+ lines of other content
2. R322 and R290 sections appeared first, causing the orchestrator to potentially stop reading
3. The orchestrator needs PRIMARY DIRECTIVES immediately to know which rules to acknowledge

## Scope of Issue

- **Total SPAWN states affected**: 26
- **SPAWN_CODE_REVIEWER states**: 16 files
- **Other SPAWN states**: 10 files (1 file SPAWN_SW_ENGINEER_PHASE_FIXES was not found)

## Solution Implemented

### 1. Moved PRIMARY DIRECTIVES to Top
- PRIMARY DIRECTIVES now appears immediately after the state title
- Clear, consistent format across all SPAWN states
- Each state lists its specific rules plus common orchestrator rules

### 2. Standard Format Applied
```markdown
# [State Name] State Rules

# PRIMARY DIRECTIVES

You MUST read and acknowledge these rules:

1. **R006** - Orchestrator cannot write code (BLOCKING)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R006-orchestrator-never-writes-code.md`

[State-specific rules...]

[Common rules: R287, R288, R304, R322, R324...]
```

### 3. State-Specific Rules Included
- **MERGE_PLAN states**: Added R269, R270
- **FIX_PLAN states**: Added R256
- **TEST_PLANNING states**: Added R355
- **EFFORT_PLANNING states**: Added R251, R309
- **ARCHITECT states**: Added R308
- **INTEGRATION states**: Added R269, R270
- **SW_ENGINEER states**: Added R232, R220
- **SPAWN_AGENTS**: Added R151, R208

## Files Fixed

### SPAWN_CODE_REVIEWER States (16):
1. SPAWN_CODE_REVIEWER_BACKPORT_PLAN
2. SPAWN_CODE_REVIEWER_FIX_PLAN
3. SPAWN_CODE_REVIEWER_INTEGRATION_FIX_PLAN
4. SPAWN_CODE_REVIEWER_MERGE_PLAN ✓ (original issue)
5. SPAWN_CODE_REVIEWER_PHASE_FIX_PLAN
6. SPAWN_CODE_REVIEWER_PHASE_IMPL
7. SPAWN_CODE_REVIEWER_PHASE_MERGE_PLAN
8. SPAWN_CODE_REVIEWER_PHASE_TEST_PLANNING
9. SPAWN_CODE_REVIEWER_PROJECT_FIX_PLANNING
10. SPAWN_CODE_REVIEWER_PROJECT_MERGE_PLAN
11. SPAWN_CODE_REVIEWER_PROJECT_TEST_PLANNING
12. SPAWN_CODE_REVIEWER_PROJECT_VALIDATION
13. SPAWN_CODE_REVIEWERS_EFFORT_PLANNING
14. SPAWN_CODE_REVIEWERS_FOR_REVIEW
15. SPAWN_CODE_REVIEWER_WAVE_IMPL
16. SPAWN_CODE_REVIEWER_WAVE_TEST_PLANNING

### Other SPAWN States (10):
1. SPAWN_AGENTS
2. SPAWN_ARCHITECT_MASTER_PLANNING
3. SPAWN_ARCHITECT_PHASE_PLANNING
4. SPAWN_ARCHITECT_WAVE_PLANNING
5. SPAWN_ENGINEERS_FOR_FIXES
6. SPAWN_INTEGRATION_AGENT
7. SPAWN_INTEGRATION_AGENT_PHASE
8. SPAWN_INTEGRATION_AGENT_PROJECT
9. SPAWN_SW_ENGINEER_BACKPORT_FIXES
10. SPAWN_SW_ENGINEER_PROJECT_FIXES

## Verification

After fixes, all SPAWN states now have:
- ✓ PRIMARY DIRECTIVES in the first 30 lines
- ✓ Clear list of rules to acknowledge
- ✓ Proper file paths for each rule
- ✓ Consistent formatting

## Impact

The orchestrator will now:
1. Immediately see PRIMARY DIRECTIVES when entering any SPAWN state
2. Know exactly which rules to read and acknowledge
3. Have clear file paths for each rule
4. Not need to search through 90+ lines to find requirements

## Scripts Created

1. `fix-spawn-code-reviewer-directives.sh` - Fixed 16 SPAWN_CODE_REVIEWER states
2. `fix-remaining-spawn-directives.sh` - Fixed 10 other SPAWN states

## Commits

- Commit 1: `25ef5ee` - Fixed all SPAWN_CODE_REVIEWER states
- Commit 2: `1bd0f05` - Fixed remaining SPAWN states

## Recommendation

Consider applying similar fixes to non-SPAWN states if they have the same issue with buried PRIMARY DIRECTIVES sections.