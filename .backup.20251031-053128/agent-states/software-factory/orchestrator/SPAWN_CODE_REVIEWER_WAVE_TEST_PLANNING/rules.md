# ORCHESTRATOR STATE: SPAWN_CODE_REVIEWER_WAVE_TEST_PLANNING

## STATE PURPOSE
Wave Test Planning with Progressive Realism - spawns Code Reviewer to create tests using ACTUAL implementations from completed waves/efforts.

## KEY RULES
- R341 (TDD): Tests before implementation
- R342 (Early Branch): Integration branch created after test planning
- R287 (TODO Persistence): Save before stopping
- R322 (Stop After Spawn): Context preservation
- R405 (Continuation Flag): CONTINUE-SOFTWARE-FACTORY=TRUE

## PROGRESSIVE TEST PLANNING
Tests use REAL imports, fixtures, mocks, and stubs from completed waves (not assumptions).

## TIMING
After wave architecture, BEFORE wave implementation planning.

## NEXT STATE
WAITING_FOR_WAVE_TEST_PLAN

## CONTEXT FOR CODE REVIEWER
- Wave architecture plan (WHAT to build)
- Preplanned test infrastructure (WHERE to put tests)
- Completed wave/effort implementations (WHAT actual code exists)
- Previous wave test fixtures (WHAT test infrastructure exists)

See: docs/PROGRESSIVE-TEST-PLANNING-ARCHITECTURE.md
