---
name: integrate-branches
description: Spawn integration agent to merge multiple feature branches while preserving history
model: opus
instruction: |
  You are the integration agent. Your task is to integrate multiple branches into a single integration branch.
  
  ## SUPREME LAWS - NEVER VIOLATE
  1. NEVER modify original branches
  2. NEVER use cherry-pick  
  3. NEVER fix upstream bugs (only document)
  
  ## Your Grading (Acknowledge This!)
  - 50% Completeness of Integration
  - 50% Meticulous Documentation
  
  ## Required Workflow
  1. Start in INIT state - acknowledge rules and grading
  2. Move to PLANNING - create INTEGRATION-PLAN.md
  3. Move to MERGING - integrate branches per plan
  4. Move to TESTING - build and test (don't fix!)
  5. Move to REPORTING - complete all documentation
  6. Move to COMPLETED - final validation
  
  ## Target Branches
  ${TARGET_BRANCHES}
  
  ## Target Base
  ${TARGET_BASE:-main}
  
  ## Integration Goals
  ${INTEGRATION_GOALS}
  
  Follow agent definition at: .claude/agents/integration.md
  Load state rules from: agent-states/main/integration/[STATE]/rules.md
---

# Integration Agent Command

This command spawns the integration specialist to merge multiple branches.

## Usage

```bash
/integrate-branches TARGET_BRANCHES="branch1 branch2 branch3" TARGET_BASE="main" INTEGRATION_GOALS="Merge all features for release"
```

## Parameters
- `TARGET_BRANCHES` - Space-separated list of branches to integrate
- `TARGET_BASE` - Base branch to integrate into (default: main)
- `INTEGRATION_GOALS` - Description of integration objectives

## Agent Responsibilities
1. Create comprehensive integration plan
2. Merge branches in optimal order
3. Resolve all conflicts
4. Document every operation meticulously
5. Test integrated code (but don't fix)
6. Produce complete documentation

## Expected Outputs
- `INTEGRATION-PLAN.md` - Integration strategy
- `work-log.md` - Replayable operation log
- `INTEGRATION-REPORT.md` - Comprehensive report
- Integration branch pushed to remote

## Critical Rules
- Original branches remain unmodified
- No cherry-picking allowed
- Upstream bugs documented not fixed
- All documentation committed to integration branch