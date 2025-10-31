# ORCHESTRATOR STATE: CREATE_PHASE_INTEGRATION_BRANCH_EARLY

## STATE PURPOSE
Create phase integration branch IMMEDIATELY after test planning (R342 enforcement). Early integration branch stores tests and serves as quality gate.

## KEY RULES
- R342 (Early Integration Branch): MUST create branch immediately after test planning
- R308 (Incremental Branching): Branch from previous phase integration branch
- R288 (State File Update): Update orchestrator-state-v3.json with branch info
- R287 (TODO Persistence): Save before transitioning

## ACTIONS
1. Determine base branch (project-integration for phase 1, previous phase for phase N)
2. Create phase-integration branch from base
3. Add phase tests to integration branch
4. Commit and push tests (R342: immediate tracking)
5. Update orchestrator-state-v3.json with branch info
6. Transition to SPAWN_CODE_REVIEWER_PHASE_IMPL

## NEXT STATE
SPAWN_CODE_REVIEWER_PHASE_IMPL

## WHY R342 MATTERS
Tests need git tracking IMMEDIATELY. Creating the branch early ensures:
- Tests versioned from inception
- No orphaned test files
- Integration branch ready for effort integration
- R341 TDD compliance enforced

See: rule-library/R342-early-integration-branch-creation.md
See: docs/PROGRESSIVE-TEST-PLANNING-ARCHITECTURE.md
