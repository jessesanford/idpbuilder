# Wave Completion Architect Review Protocol

## CRITICAL: Architectural Review Required Before Next Wave

### Overview
At the completion of EACH wave, the orchestrator MUST spawn @agent-architect-reviewer to review all work completed in the current phase (all waves so far) and provide architectural guidance for the next wave. The orchestrator MUST NOT proceed to the next wave until architect feedback is addressed.

## When Architect Review is Required

### Mandatory Review Points
```yaml
wave_completion_gates:
  # After EVERY wave completion
  # Before starting ANY new wave
  # At phase boundaries
  # When integration needed
```

## Orchestrator Workflow for Wave Completion

### Step 1: Detect Wave Completion
```python
def is_wave_complete(phase, wave):
    """
    Check if all efforts in current wave are completed
    """
    wave_efforts = get_wave_efforts(phase, wave)
    completed_efforts = state['efforts_completed']
    
    for effort in wave_efforts:
        if effort not in completed_efforts:
            return False
    return True
```

### Step 2: Spawn Architect Reviewer

When wave is complete, orchestrator MUST:

```markdown
Task for @agent-architect-reviewer:

PURPOSE: Architectural review of Phase ${PHASE} Wave ${WAVE} completion

MANDATORY STARTUP:
1. Print timestamp and acknowledgment
2. State: "ARCHITECT REVIEW FOR WAVE COMPLETION"

REVIEW SCOPE:
1. READ: /workspaces/[project]/orchestrator-state.yaml
2. READ: /workspaces/[project]/phase-plans/PHASE${PHASE}-SPECIFIC-IMPL-PLAN.md

CURRENT PHASE PROGRESS:
Phase ${PHASE} - ${PHASE_NAME}
Waves Completed: ${WAVES_COMPLETED}
Current Wave Just Completed: Wave ${WAVE}

COMPLETED EFFORTS IN THIS PHASE:
${LIST_ALL_COMPLETED_EFFORTS_WITH_BRANCHES}

Example:
- E1.1.1 (data-models): branch phase1/wave1/effort1-data-models
- E1.1.2 (api-interfaces): branch phase1/wave1/effort2-api-interfaces
- E1.1.3 (database-schema): branch phase1/wave1/effort3-database-schema

ARCHITECTURAL REVIEW CHECKLIST:
1. Design Patterns Consistency
   - [ ] All efforts follow established patterns
   - [ ] No architectural drift detected
   - [ ] Interfaces properly abstracted
   
2. Integration Points
   - [ ] APIs properly versioned
   - [ ] Contracts well-defined
   - [ ] Backward compatibility maintained
   
3. Technical Debt Assessment
   - [ ] Acceptable debt levels
   - [ ] Debt documented if present
   - [ ] Mitigation plan for debt
   
4. Performance Implications
   - [ ] No architectural bottlenecks introduced
   - [ ] Scalability patterns maintained
   - [ ] Resource usage acceptable
   
5. Security Review
   - [ ] Authentication patterns consistent
   - [ ] Authorization properly implemented
   - [ ] No security anti-patterns
   
6. Testing Coverage
   - [ ] Architectural test points covered
   - [ ] Integration tests adequate
   - [ ] E2E paths validated

SPECIFIC QUESTIONS FOR THIS WAVE:
${WAVE_SPECIFIC_ARCHITECTURAL_CONCERNS}

REQUIRED OUTPUT:
Provide one of these assessments:

1. PROCEED - Architecture sound, continue to next wave
2. PROCEED_WITH_MINOR_CORRECTIONS - Continue but address notes
3. CHANGES_REQUIRED - Must fix issues before next wave
4. ARCHITECTURAL_DRIFT - Significant corrections needed
5. STOP - Critical architectural violation

If not PROCEED, create:
PHASE${PHASE}-WAVE${WAVE}-ARCHITECTURAL-CORRECTIONS.md
with specific issues and required fixes.
```

### Step 3: Process Architect Decision

```python
def process_architect_decision(decision):
    if decision == "PROCEED":
        # Continue to next wave
        return "START_NEXT_WAVE"
    
    elif decision == "PROCEED_WITH_MINOR_CORRECTIONS":
        # Document corrections, continue
        create_addendum()
        return "START_NEXT_WAVE_WITH_ADDENDUM"
    
    elif decision == "CHANGES_REQUIRED":
        # Must fix before proceeding
        create_fix_list()
        return "SPAWN_SW_ENGINEER_FOR_FIXES"
    
    elif decision == "ARCHITECTURAL_DRIFT":
        # Major corrections needed
        create_correction_plan()
        return "IMPLEMENT_CORRECTIONS"
    
    elif decision == "STOP":
        # Critical issue
        return "HALT_IMPLEMENTATION"
```

## Integration Branch Creation

After architect approval (PROCEED or PROCEED_WITH_MINOR_CORRECTIONS):

### Create Wave Integration Branch
```bash
# Orchestrator creates integration branch
git checkout main
git checkout -b phase${PHASE}/wave${WAVE}-integration

# Merge all effort branches for this wave
for effort_branch in ${WAVE_EFFORT_BRANCHES}; do
    git merge ${effort_branch}
done

# Run tests
make test

# If successful, mark wave integrated
```

## Architectural Addendum Format

If corrections needed, architect creates:

```markdown
# PHASE${PHASE}-WAVE${WAVE}-ARCHITECTURAL-CORRECTIONS.md

## Issues Identified

### Issue 1: ${ISSUE_TITLE}
- **Severity**: HIGH/MEDIUM/LOW
- **Affected Efforts**: ${LIST}
- **Description**: ${DESCRIPTION}
- **Required Fix**: ${FIX_DESCRIPTION}
- **Assigned To**: ${EFFORT_NUMBER}

### Issue 2: ${ISSUE_TITLE}
...

## Implementation Guidance for Next Wave

### Architectural Constraints
- ${CONSTRAINT_1}
- ${CONSTRAINT_2}

### Pattern Requirements
- ${PATTERN_1}
- ${PATTERN_2}

### Integration Considerations
- ${CONSIDERATION_1}
- ${CONSIDERATION_2}

## Success Criteria for Fixes
- [ ] ${CRITERIA_1}
- [ ] ${CRITERIA_2}
- [ ] ${CRITERIA_3}
```

## State Machine Integration

```yaml
# orchestrator-state.yaml update
wave_reviews:
  - phase: ${PHASE}
    wave: ${WAVE}
    decision: "${DECISION}"
    date: "${TIMESTAMP}"
    issues: [${ISSUE_LIST}]
    corrections_file: "${CORRECTIONS_FILE}"
    integration_branch: "phase${PHASE}/wave${WAVE}-integration"
```

## Recovery from CHANGES_REQUIRED

When architect requires changes:

1. **Orchestrator creates fix efforts**
2. **Spawns SW Engineers for fixes**
3. **Re-runs architect review**
4. **Only proceeds after approval**

## Important Notes

1. **Wave review is mandatory** - Cannot skip
2. **Integration before next wave** - Ensures compatibility
3. **Architect owns architecture** - Final say on patterns
4. **Document all decisions** - Maintain audit trail
5. **Fix before proceeding** - Don't accumulate debt

## Benefits

- **Early detection** of architectural drift
- **Consistent patterns** across efforts
- **Managed technical debt**
- **Clear integration points**
- **Architectural integrity** maintained

This protocol ensures architectural consistency throughout the implementation while catching issues early when they're easier to fix.