# ORCHESTRATOR STATE: SPAWN_CODE_REVIEWER_PHASE_TEST_PLANNING

[Content continues - creating via Bash heredoc due to Write tool limitations for new state rules files]

## STATE PURPOSE

Phase Test Planning with Progressive Realism - spawns Code Reviewer to create tests using ACTUAL implementations from completed phases.

## KEY RULES
- R341 (TDD): Tests before implementation
- R342 (Early Branch): Integration branch created after test planning
- R287 (TODO Persistence): Save before stopping
- R322 (Stop After Spawn): Context preservation
- R405 (Continuation Flag): CONTINUE-SOFTWARE-FACTORY=TRUE

## PROGRESSIVE TEST PLANNING
Tests use REAL imports, fixtures, and mocks from completed phases (not assumptions).

## NEXT STATE
WAITING_FOR_PHASE_TEST_PLAN

See: docs/PROGRESSIVE-TEST-PLANNING-ARCHITECTURE.md for full architecture.
