# 🚨🚨🚨 BLOCKING RULE R289: Orphaned State Detection and Prevention

## Purpose
Prevent confusion and system failures by detecting and eliminating orphaned states that have no valid transitions in the state machine.

## Definition
An **orphaned state** is a state that exists in the system but has:
1. No valid transitions TO it (unreachable)
2. No valid transitions FROM it (dead-end, unless terminal)
3. Both conditions (completely orphaned)

## Terminal States Exception
The following states are EXPECTED to have no outgoing transitions (terminal states):
- SUCCESS
- HARD_STOP
- COMPLETED
- BLOCKED
- DECISION

These are valid endpoints and NOT considered orphaned.

## Detection Protocol

### Automated Detection
```bash
# Find all defined states
grep '^- \*\*[A-Z_]*\*\*' SOFTWARE-FACTORY-STATE-MACHINE.md | \
    sed 's/- \*\*//' | sed 's/\*\*.*//' | sort -u > /tmp/all_states.txt

# Find states with transitions
grep -E '→' SOFTWARE-FACTORY-STATE-MACHINE.md | \
    sed 's/→/\n/g' | grep -E '^[A-Z_]+' | \
    sed 's/[^A-Z_].*//' | sort -u > /tmp/states_with_transitions.txt

# Identify orphaned states
comm -23 /tmp/all_states.txt /tmp/states_with_transitions.txt | \
    grep -v -E '^(SUCCESS|HARD_STOP|COMPLETED|BLOCKED|DECISION)$'
```

### Manual Verification
For each non-terminal state, verify:
1. At least one transition TO it exists
2. At least one transition FROM it exists
3. The state has a rules file in agent-states/[agent]/[STATE]/

## Prevention Requirements

### 1. State Addition Protocol
When adding a new state:
- ✅ MUST add at least one transition TO it
- ✅ MUST add at least one transition FROM it (unless terminal)
- ✅ MUST create corresponding rules.md file
- ✅ MUST document purpose in state description

### 2. State Removal Protocol
When removing a state:
- ✅ MUST remove ALL transitions referencing it
- ✅ MUST archive state directory (don't delete)
- ✅ MUST update all documentation
- ✅ MUST provide migration guide for existing state files

### 3. Validation Requirements
Before any commit involving state changes:
- ✅ Run orphaned state detection
- ✅ Verify no new orphaned states introduced
- ✅ Fix any detected orphaned states
- ✅ Document in commit message

## Known Deprecated States

### PLANNING State (Deprecated 2025-08-30)
**Status**: REMOVED - Was orphaned with no valid transitions
**Replacement**: Use `SPAWN_CODE_REVIEWERS_EFFORT_PLANNING` for planning activities
**Migration**: Any state file showing `current_state: PLANNING` should transition to `ERROR_RECOVERY`

## Enforcement

### Orchestrator Requirements
- MUST validate target state has valid transitions before transitioning
- MUST error if attempting to transition to orphaned state
- MUST log warning if orphaned state directories detected

### Factory Manager Requirements
- MUST run orphaned state detection during audits
- MUST report any orphaned states found
- MUST prevent merges with orphaned states

### CI/CD Requirements
- MUST fail builds if orphaned states detected
- MUST run validation on all state machine changes
- MUST enforce state directory consistency

## Error Messages

### Transition to Orphaned State
```
❌ CRITICAL ERROR: Attempted transition to orphaned state 'PLANNING'
This state has no valid transitions and is deprecated.
Use SPAWN_CODE_REVIEWERS_EFFORT_PLANNING for planning activities.
```

### Orphaned State Detection
```
⚠️ WARNING: Orphaned state detected: PLANNING
- No transitions TO this state
- No transitions FROM this state
- State directory still exists
Action: Remove or properly connect this state
```

## Recovery Procedures

### If Currently in Orphaned State
1. Transition to ERROR_RECOVERY immediately
2. Assess intended goal
3. Transition to appropriate connected state
4. Update state file

### If Orphaned State Directory Found
1. Archive with .DEPRECATED-YYYYMMDD suffix
2. Update any references in documentation
3. Commit with clear deprecation message
4. Push changes

## Validation Script Location
`/home/vscode/software-factory-template/utilities/detect-orphaned-states.sh`

## Related Rules
- R206: State machine transition validation
- R233: All states require immediate action
- R234: Mandatory state traversal

## Penalties
- Transitioning to orphaned state: -50% grade
- Creating new orphaned state: -30% grade
- Leaving orphaned states unfixed: -20% grade