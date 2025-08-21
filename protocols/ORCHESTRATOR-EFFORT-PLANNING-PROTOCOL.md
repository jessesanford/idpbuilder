# Orchestrator Effort Planning Protocol

## CRITICAL: Code Reviewer Creates Effort-Specific Plans BEFORE Implementation

### Overview
The orchestrator MUST spawn @agent-code-reviewer to create detailed, effort-specific implementation plans for EACH effort BEFORE tasking any @agent-sw-engineer to work on it.

## Planning Workflow

### Step 1: Orchestrator Tasks Code Reviewer for Planning
When starting an effort (e.g., E1.1.1), the orchestrator spawns @agent-code-reviewer with:
1. Current phase plan location (PHASE{N}-SPECIFIC-IMPL-PLAN.md)
2. Orchestrator state file with completed efforts and branches
3. Working copy path for the new effort
4. Instructions to create effort implementation plan

### Step 2: Code Reviewer Creates Effort Implementation Plan
The Code Reviewer:
1. Reads PHASE{N}-SPECIFIC-IMPL-PLAN.md for current phase
2. Examines orchestrator-state.yaml for completed efforts/branches
3. THINKS about work done so far and potential impacts
4. Creates IMPLEMENTATION-PLAN-ADDENDUM-{Date}.md if changes needed
5. Creates EFFORT-IMPLEMENTATION-PLAN-{Date}-{Time}.md in working copy

```markdown
# Location: /workspaces/efforts/phase{X}/wave{Y}/effort{Z}-{name}/EFFORT-IMPLEMENTATION-PLAN-{Date}-{Time}.md

# Effort E{X}.{Y}.{Z} Implementation Plan
Generated: {timestamp}
Created by: @agent-code-reviewer
Reviewed Phase Plan: PHASE{X}-SPECIFIC-IMPL-PLAN.md

## Context Analysis
### Completed Efforts in Current Wave
[List from orchestrator-state.yaml]
- E{X}.{Y}.1: {branch-name} - {status}
- E{X}.{Y}.2: {branch-name} - {status}

### Adjustments Based on Progress
[Analysis of what's been implemented]
- {adjustment-1}: {reason}
- {adjustment-2}: {reason}

## Effort Overview
- Phase: {X}
- Wave: {Y}  
- Effort: {Z}
- Name: {descriptive-name}
- Base Branch: {base-branch}
- Working Copy: /workspaces/efforts/phase{X}/wave{Y}/effort{Z}-{name}

## Specific Requirements
[Extract from phase plan with adjustments]
1. {requirement-1}
2. {requirement-2}
3. {requirement-3}

## Implementation Steps
1. {step-1}: {specific-action}
2. {step-2}: {specific-action}
3. {step-3}: {specific-action}

## Files to Create/Modify
- {file-1}: {purpose}
- {file-2}: {purpose}
- {file-3}: {purpose}

## Integration Instructions
```bash
# Any specific integration steps
# Merge patterns, API contracts, etc.
```

## Test Requirements
- Coverage: {X}%
- Specific tests needed:
  - {test-1}
  - {test-2}

## Size Constraints
- Target: {target} lines
- Maximum: {limit} lines (measured by line-counter.sh)
- Split strategy if exceeded: {strategy}

## Success Criteria
- [ ] All requirements implemented
- [ ] Tests pass with required coverage
- [ ] Size under limit per line-counter.sh
- [ ] No hardcoded values
- [ ] Follows project style guide
- [ ] Integrates with completed efforts

## Dependencies
- Depends on: {list-efforts}
- Blocks: {list-efforts}
```

### Step 3: Create Work Log Template
```markdown
# Location: /workspaces/efforts/phase{X}/wave{Y}/effort{Z}-{name}/work-log.md

# Work Log - Effort E{X}.{Y}.{Z}
Created: {timestamp}
Engineer: [To be assigned]

## Sessions
### Session 1: {timestamp}
- [ ] Environment verified
- [ ] Implementation plan read
- [ ] Work started

## Size Tracking
| Checkpoint | Lines | Status |
|------------|-------|--------|
| Initial | 0 | OK |

## Issues
- None yet

## Notes
- Implementation plan: EFFORT-IMPLEMENTATION-PLAN-{Date}-{Time}.md
```

### Step 4: Orchestrator Verifies Plan Creation
```yaml
# Update orchestrator-state.yaml
efforts_in_progress:
  - phase: X
    wave: Y
    effort: Z
    name: {name}
    status: "PLAN_CREATED"
    plan_file: "EFFORT-IMPLEMENTATION-PLAN-{Date}-{Time}.md"
    working_dir: "/workspaces/efforts/phase{X}/wave{Y}/effort{Z}-{name}"
```

## Implementation Task Template

After plan is created, orchestrator tasks SW Engineer:

```markdown
Task @agent-sw-engineer:

MANDATORY STARTUP: Follow /workspaces/[project]/orchestrator/SW-ENGINEER-STARTUP-REQUIREMENTS.md

Working directory: /workspaces/efforts/phase{X}/wave{Y}/effort{Z}-{name}
Branch: phase{X}/wave{Y}/effort{Z}-{name}

MANDATORY: Read these files IN ORDER:
1. EFFORT-IMPLEMENTATION-PLAN-{Date}-{Time}.md - Your specific implementation plan (created by Code Reviewer)
2. /workspaces/[project]/orchestrator/SW-ENGINEER-EXPLICIT-INSTRUCTIONS.md
3. /workspaces/[project]/orchestrator/SW-ENGINEER-STARTUP-REQUIREMENTS.md

Your EFFORT-IMPLEMENTATION-PLAN contains:
- Context analysis of completed work
- Exact requirements to implement  
- Step-by-step instructions
- Files to create/modify
- Test requirements
- Success criteria checklist

SIZE LIMIT: {limit} lines (measured by line-counter.sh)
- Measure after every logical change
- Stop if approaching limit
- Report if exceeded

Update work-log.md as you work.
```

## Benefits of This Protocol

1. **Context Awareness**: Code Reviewer sees all completed work
2. **Adaptation**: Plans adjust based on what's implemented
3. **Consistency**: All efforts get proper planning
4. **Traceability**: Plans document decisions
5. **Quality**: Clear success criteria upfront

## When Plans Need Updates

If during implementation, the SW Engineer discovers:
- Missing requirements
- Changed dependencies  
- Integration issues
- Size constraints

Then:
1. Document in work-log.md
2. Report to orchestrator
3. Orchestrator spawns Code Reviewer to update plan
4. Resume with updated plan

## Plan Versioning

Each plan is timestamped:
- `EFFORT-IMPLEMENTATION-PLAN-2025-01-21-14-30.md`
- `EFFORT-IMPLEMENTATION-PLAN-2025-01-21-16-45.md` (updated)

This allows tracking plan evolution.

## Important Notes

1. **Never skip planning** - Every effort needs a plan
2. **Plans are living documents** - Update when needed
3. **Code Reviewer owns planning** - Not orchestrator, not SW Engineer
4. **Plans drive implementation** - SW Engineer follows the plan
5. **Size limits are absolute** - Plan for splits upfront