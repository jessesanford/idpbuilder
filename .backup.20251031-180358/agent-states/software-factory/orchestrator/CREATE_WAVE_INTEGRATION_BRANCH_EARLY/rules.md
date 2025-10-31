# ORCHESTRATOR STATE: CREATE_WAVE_INTEGRATION_BRANCH_EARLY

## STATE PURPOSE
Create wave integration branch IMMEDIATELY after test planning (R342 enforcement). Early integration branch stores wave tests.

## KEY RULES
- R342 (Early Integration Branch): MUST create branch immediately after test planning
- R336 (Wave Integration): Branch from previous wave or phase integration branch
- R288 (State File Update): Update orchestrator-state-v3.json with branch info
- R287 (TODO Persistence): Save before transitioning

## ACTIONS
1. Determine base branch (phase-integration for wave 1, previous wave for wave N)
2. Create wave-integration branch from base
3. Add wave tests to integration branch
4. Commit and push tests (R342: immediate tracking)
5. Update orchestrator-state-v3.json with branch info
6. Transition to SPAWN_CODE_REVIEWER_WAVE_IMPL

## NEXT STATE
SPAWN_CODE_REVIEWER_WAVE_IMPL

## BRANCH NAMING
Format: `phase-${PHASE}-wave-${WAVE}-integration`
Example: `phase-1-wave-2-integration`

## WHY R342 MATTERS
Wave tests must be tracked immediately to ensure:
- Tests available for effort integration
- No test loss between waves
- R341 TDD compliance enforced
- Cumulative test accumulation (R308/R336)

See: rule-library/R342-early-integration-branch-creation.md
See: docs/PROGRESSIVE-TEST-PLANNING-ARCHITECTURE.md
