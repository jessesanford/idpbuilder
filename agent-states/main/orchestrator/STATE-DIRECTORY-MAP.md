# Orchestrator State Directory Map

## Complete List of Valid Orchestrator States

Each state MUST have a corresponding directory with rules.md file.

### ✅ States with Rules Files:

1. **INIT** - `/agent-states/orchestrator/INIT/rules.md`
   - Initial startup and state detection

2. **WAVE_START** - `/agent-states/orchestrator/WAVE_START/rules.md`
   - Beginning a new wave of efforts

3. **SETUP_EFFORT_INFRASTRUCTURE** - `/agent-states/orchestrator/SETUP_EFFORT_INFRASTRUCTURE/rules.md`
   - Creating effort directories and clones

4. **ANALYZE_CODE_REVIEWER_PARALLELIZATION** - `/agent-states/orchestrator/ANALYZE_CODE_REVIEWER_PARALLELIZATION/rules.md`
   - MANDATORY - Analyzing wave plan to determine Code Reviewer spawn strategy

5. **SPAWN_CODE_REVIEWERS_EFFORT_PLANNING** - `/agent-states/orchestrator/SPAWN_CODE_REVIEWERS_EFFORT_PLANNING/rules.md`
   - Spawning code reviewers to create implementation plans

6. **WAITING_FOR_EFFORT_PLANS** - `/agent-states/orchestrator/WAITING_FOR_EFFORT_PLANS/rules.md`
   - Waiting for code reviewers to complete plans

7. **ANALYZE_IMPLEMENTATION_PARALLELIZATION** - `/agent-states/orchestrator/ANALYZE_IMPLEMENTATION_PARALLELIZATION/rules.md`
   - MANDATORY - Analyzing effort plans to determine SW Engineer spawn strategy

8. **SPAWN_AGENTS** - `/agent-states/orchestrator/SPAWN_AGENTS/rules.md`
   - Spawning SW engineers for implementation

9. **MONITOR_IMPLEMENTATION** - `/agent-states/orchestrator/MONITOR_IMPLEMENTATION/rules.md`
   - Monitoring SW Engineer implementation progress

9a. **MONITOR_REVIEWS** - `/agent-states/orchestrator/MONITOR_REVIEWS/rules.md`
   - Monitoring Code Reviewer progress

9b. **MONITOR_FIXES** - `/agent-states/orchestrator/MONITOR_FIXES/rules.md`
   - Monitoring fix implementation progress

10. **WAVE_COMPLETE** - `/agent-states/orchestrator/WAVE_COMPLETE/rules.md`
    - All wave efforts completed and reviewed

11. **INTEGRATION** - `/agent-states/orchestrator/INTEGRATION/rules.md`
    - Creating integration branches in target repo

12. **WAVE_REVIEW** - `/agent-states/orchestrator/WAVE_REVIEW/rules.md`
    - Architect reviewing wave integration

13. **ERROR_RECOVERY** - `/agent-states/orchestrator/ERROR_RECOVERY/rules.md`
    - Handling errors and failures

14. **SPAWN_ARCHITECT_PHASE_ASSESSMENT** - `/agent-states/orchestrator/SPAWN_ARCHITECT_PHASE_ASSESSMENT/rules.md`
    - Request architect to assess complete phase (last wave only)

15. **WAITING_FOR_PHASE_ASSESSMENT** - `/agent-states/orchestrator/WAITING_FOR_PHASE_ASSESSMENT/rules.md`
    - Waiting for architect phase assessment decision

16. **PHASE_COMPLETE** - `/agent-states/orchestrator/PHASE_COMPLETE/rules.md`
    - Phase assessment passed, handling phase-level integration

17. **SUCCESS** - `/agent-states/orchestrator/SUCCESS/rules.md`
    - Project successfully completed (only after phase assessment)

18. **HARD_STOP** - `/agent-states/orchestrator/HARD_STOP/rules.md`
    - Critical failure requiring manual intervention

19. **SPAWN_ARCHITECT_PHASE_PLANNING** - `/agent-states/orchestrator/SPAWN_ARCHITECT_PHASE_PLANNING/rules.md`
    - Request architect to create phase architecture (R210)

20. **SPAWN_ARCHITECT_WAVE_PLANNING** - `/agent-states/orchestrator/SPAWN_ARCHITECT_WAVE_PLANNING/rules.md`
    - Request architect to create wave architecture (R210)

21. **SPAWN_CODE_REVIEWER_PHASE_IMPL** - `/agent-states/orchestrator/SPAWN_CODE_REVIEWER_PHASE_IMPL/rules.md`
    - Request code reviewer to create phase implementation from architecture (R211)

22. **SPAWN_CODE_REVIEWER_WAVE_IMPL** - `/agent-states/orchestrator/SPAWN_CODE_REVIEWER_WAVE_IMPL/rules.md`
    - Request code reviewer to create wave implementation from architecture (R211)

23. **WAITING_FOR_ARCHITECTURE_PLAN** - `/agent-states/orchestrator/WAITING_FOR_ARCHITECTURE_PLAN/rules.md`
    - Waiting for architect to complete architecture plan

24. **WAITING_FOR_IMPLEMENTATION_PLAN** - `/agent-states/orchestrator/WAITING_FOR_IMPLEMENTATION_PLAN/rules.md`
    - Waiting for code reviewer to complete implementation plan

25. **INJECT_WAVE_METADATA** - `/agent-states/orchestrator/INJECT_WAVE_METADATA/rules.md`
    - Injecting R213 wave metadata into plans

26. **CREATE_NEXT_SPLIT_INFRASTRUCTURE** - `/agent-states/orchestrator/CREATE_NEXT_SPLIT_INFRASTRUCTURE/rules.md`
    - Creating infrastructure for the next split in sequence (R204 just-in-time)

27. **SPAWN_CODE_REVIEWER_MERGE_PLAN** - `/agent-states/orchestrator/SPAWN_CODE_REVIEWER_MERGE_PLAN/rules.md`
    - Spawning Code Reviewer to create merge plan

28. **WAITING_FOR_MERGE_PLAN** - `/agent-states/orchestrator/WAITING_FOR_MERGE_PLAN/rules.md`
    - Waiting for Code Reviewer merge plan completion

29. **SPAWN_INTEGRATION_AGENT** - `/agent-states/orchestrator/SPAWN_INTEGRATION_AGENT/rules.md`
    - Spawning Integration Agent to execute merges

30. **MONITORING_INTEGRATION** - `/agent-states/orchestrator/MONITORING_INTEGRATION/rules.md`
    - Monitoring Integration Agent progress and checking for reports (R238)

31. **SPAWN_CODE_REVIEWERS_FOR_REVIEW** - `/agent-states/orchestrator/SPAWN_CODE_REVIEWERS_FOR_REVIEW/rules.md`
    - Spawning Code Reviewers to review fixed code

32. **SPAWN_ENGINEERS_FOR_FIXES** - `/agent-states/orchestrator/SPAWN_ENGINEERS_FOR_FIXES/rules.md`
    - Spawning SW Engineers to implement integration fixes

33. **MONITORING_FIX_PROGRESS** - `/agent-states/orchestrator/MONITORING_FIX_PROGRESS/rules.md`
    - Monitoring engineers implementing fixes

34. **PHASE_INTEGRATION** - `/agent-states/orchestrator/PHASE_INTEGRATION/rules.md`
    - Setting up phase integration infrastructure

35. **PROJECT_INTEGRATION** - `/agent-states/orchestrator/PROJECT_INTEGRATION/rules.md`
    - Setting up project-level integration for all phases (R283)

36. **SPAWN_CODE_REVIEWER_PROJECT_MERGE_PLAN** - `/agent-states/orchestrator/SPAWN_CODE_REVIEWER_PROJECT_MERGE_PLAN/rules.md`
    - Spawning Code Reviewer to create project merge plan

37. **WAITING_FOR_PROJECT_MERGE_PLAN** - `/agent-states/orchestrator/WAITING_FOR_PROJECT_MERGE_PLAN/rules.md`
    - Waiting for Code Reviewer project merge plan

38. **SPAWN_INTEGRATION_AGENT_PROJECT** - `/agent-states/orchestrator/SPAWN_INTEGRATION_AGENT_PROJECT/rules.md`
    - Spawning Integration Agent to merge all phases

39. **MONITORING_PROJECT_INTEGRATION** - `/agent-states/orchestrator/MONITORING_PROJECT_INTEGRATION/rules.md`
    - Monitoring project-level integration progress

40. **SPAWN_CODE_REVIEWER_PROJECT_VALIDATION** - `/agent-states/orchestrator/SPAWN_CODE_REVIEWER_PROJECT_VALIDATION/rules.md`
    - Spawning Code Reviewer for project validation

41. **WAITING_FOR_PROJECT_VALIDATION** - `/agent-states/orchestrator/WAITING_FOR_PROJECT_VALIDATION/rules.md`
    - Waiting for project validation results

42. **CREATE_INTEGRATION_TESTING** - `/agent-states/orchestrator/CREATE_INTEGRATION_TESTING/rules.md`
    - Creating integration-testing branch from project integration (R272)

43. **INTEGRATION_TESTING** - `/agent-states/orchestrator/INTEGRATION_TESTING/rules.md`
    - Final validation in integration-testing branch (R271)

44. **PRODUCTION_READY_VALIDATION** - `/agent-states/orchestrator/PRODUCTION_READY_VALIDATION/rules.md`
    - Validating software is production-ready (R273-R275)

45. **BUILD_VALIDATION** - `/agent-states/orchestrator/BUILD_VALIDATION/rules.md`
    - Final build and deployment verification (R277)

46. **ANALYZE_BUILD_FAILURES** - `/agent-states/orchestrator/ANALYZE_BUILD_FAILURES/rules.md`
    - Orchestrator analyzing build errors and categorizing failures

47. **COORDINATE_BUILD_FIXES** - `/agent-states/orchestrator/COORDINATE_BUILD_FIXES/rules.md`
    - Orchestrator distributing fix work to SW Engineers with proper tracking

48. **IMMEDIATE_BACKPORT_REQUIRED** - `/agent-states/orchestrator/IMMEDIATE_BACKPORT_REQUIRED/rules.md`
    - R321 enforcement: fixing source branches immediately when integration issues found

49. **SPAWN_CODE_REVIEWER_BACKPORT_PLAN** - `/agent-states/orchestrator/SPAWN_CODE_REVIEWER_BACKPORT_PLAN/rules.md`
    - Spawn Code Reviewer to create backport plan

50. **WAITING_FOR_BACKPORT_PLAN** - `/agent-states/orchestrator/WAITING_FOR_BACKPORT_PLAN/rules.md`
    - Waiting for Code Reviewer to complete backport plan

51. **SPAWN_SW_ENGINEER_BACKPORT_FIXES** - `/agent-states/orchestrator/SPAWN_SW_ENGINEER_BACKPORT_FIXES/rules.md`
    - Spawn SW Engineers to implement backport fixes

52. **MONITORING_BACKPORT_PROGRESS** - `/agent-states/orchestrator/MONITORING_BACKPORT_PROGRESS/rules.md`
    - Monitor SW Engineers implementing backports

53. **PR_PLAN_CREATION** - `/agent-states/orchestrator/PR_PLAN_CREATION/rules.md`
    - Generating MASTER-PR-PLAN.md for human PRs (R279)

54. **SPAWN_CODE_REVIEWER_PHASE_MERGE_PLAN** - `/agent-states/orchestrator/SPAWN_CODE_REVIEWER_PHASE_MERGE_PLAN/rules.md`
    - Spawning Code Reviewer for phase merge plan

55. **WAITING_FOR_PHASE_MERGE_PLAN** - `/agent-states/orchestrator/WAITING_FOR_PHASE_MERGE_PLAN/rules.md`
    - Waiting for Code Reviewer phase merge plan

56. **SPAWN_INTEGRATION_AGENT_PHASE** - `/agent-states/orchestrator/SPAWN_INTEGRATION_AGENT_PHASE/rules.md`
    - Spawning Integration Agent for phase merges

57. **MONITORING_PHASE_INTEGRATION** - `/agent-states/orchestrator/MONITORING_PHASE_INTEGRATION/rules.md`
    - Monitoring Integration Agent phase progress and checking reports (R238)

58. **SPAWN_CODE_REVIEWER_PHASE_FIX_PLAN** - `/agent-states/orchestrator/SPAWN_CODE_REVIEWER_PHASE_FIX_PLAN/rules.md`
    - Spawning Code Reviewer for phase-level fix plans

59. **WAITING_FOR_PHASE_FIX_PLANS** - `/agent-states/orchestrator/WAITING_FOR_PHASE_FIX_PLANS/rules.md`
    - Waiting for phase-level fix plans

60. **SPAWN_CODE_REVIEWER_FIX_PLAN** - `/agent-states/orchestrator/SPAWN_CODE_REVIEWER_FIX_PLAN/rules.md`
    - Spawning Code Reviewer to create fix plans for failures

61. **WAITING_FOR_FIX_PLANS** - `/agent-states/orchestrator/WAITING_FOR_FIX_PLANS/rules.md`
    - Waiting for Code Reviewer to complete fix plans

62. **DISTRIBUTE_FIX_PLANS** - `/agent-states/orchestrator/DISTRIBUTE_FIX_PLANS/rules.md`
    - Distributing fix plans to effort directories (R239)

63. **IMMEDIATE_BACKPORT_REQUIRED** - `/agent-states/orchestrator/IMMEDIATE_BACKPORT_REQUIRED/rules.md`
    - R321 enforcement: Fixing integration issues in source branches immediately

## State Transition Flow

```
INIT 
  ↓
SPAWN_ARCHITECT_PHASE_PLANNING or WAVE_START
  ↓
SETUP_EFFORT_INFRASTRUCTURE
  ↓
ANALYZE_CODE_REVIEWER_PARALLELIZATION  ← MANDATORY GATE
  ↓
SPAWN_CODE_REVIEWERS_EFFORT_PLANNING
  ↓
WAITING_FOR_EFFORT_PLANS
  ↓
ANALYZE_IMPLEMENTATION_PARALLELIZATION  ← MANDATORY GATE
  ↓
SPAWN_AGENTS
  ↓
MONITOR
  ↓
WAVE_COMPLETE
  ↓
INTEGRATION
  ↓
WAVE_REVIEW
  ↓
(Next wave WAVE_START or SPAWN_ARCHITECT_PHASE_ASSESSMENT if last wave)
  ↓
SPAWN_ARCHITECT_PHASE_ASSESSMENT (last wave only)
  ↓
WAITING_FOR_PHASE_ASSESSMENT
  ↓
PHASE_COMPLETE
  ↓
SUCCESS
```

## Common State Transition Errors

1. **Missing State Directory**: If transitioning to a state without a rules.md file
   - Solution: Create the directory and rules.md file

2. **Invalid State Name**: Transitioning to a state not in the list
   - Solution: Check spelling and use exact state names

3. **Wrong Agent State**: Using a state from another agent type
   - Example: IMPLEMENTATION is for SW engineers, not orchestrator

## Verification Command

```bash
# Verify all state directories exist
for state in INIT WAVE_START SETUP_EFFORT_INFRASTRUCTURE \
  ANALYZE_CODE_REVIEWER_PARALLELIZATION \
  SPAWN_CODE_REVIEWERS_EFFORT_PLANNING WAITING_FOR_EFFORT_PLANS \
  ANALYZE_IMPLEMENTATION_PARALLELIZATION \
  SPAWN_AGENTS MONITOR_IMPLEMENTATION MONITOR_REVIEWS MONITOR_FIXES \
  WAVE_COMPLETE INTEGRATION WAVE_REVIEW \
  SPAWN_ARCHITECT_PHASE_ASSESSMENT WAITING_FOR_PHASE_ASSESSMENT \
  PHASE_COMPLETE ERROR_RECOVERY SUCCESS HARD_STOP \
  SPAWN_ARCHITECT_PHASE_PLANNING SPAWN_ARCHITECT_WAVE_PLANNING \
  SPAWN_CODE_REVIEWER_PHASE_IMPL SPAWN_CODE_REVIEWER_WAVE_IMPL \
  WAITING_FOR_ARCHITECTURE_PLAN WAITING_FOR_IMPLEMENTATION_PLAN \
  INJECT_WAVE_METADATA CREATE_NEXT_SPLIT_INFRASTRUCTURE \
  SPAWN_CODE_REVIEWER_MERGE_PLAN WAITING_FOR_MERGE_PLAN \
  SPAWN_INTEGRATION_AGENT MONITORING_INTEGRATION \
  SPAWN_CODE_REVIEWERS_FOR_REVIEW SPAWN_ENGINEERS_FOR_FIXES \
  MONITORING_FIX_PROGRESS PHASE_INTEGRATION PROJECT_INTEGRATION \
  SPAWN_CODE_REVIEWER_PROJECT_MERGE_PLAN WAITING_FOR_PROJECT_MERGE_PLAN \
  SPAWN_INTEGRATION_AGENT_PROJECT MONITORING_PROJECT_INTEGRATION \
  SPAWN_CODE_REVIEWER_PROJECT_VALIDATION WAITING_FOR_PROJECT_VALIDATION \
  CREATE_INTEGRATION_TESTING INTEGRATION_TESTING \
  PRODUCTION_READY_VALIDATION BUILD_VALIDATION \
  ANALYZE_BUILD_FAILURES COORDINATE_BUILD_FIXES \
  IMMEDIATE_BACKPORT_REQUIRED SPAWN_CODE_REVIEWER_BACKPORT_PLAN \
  WAITING_FOR_BACKPORT_PLAN SPAWN_SW_ENGINEER_BACKPORT_FIXES \
  MONITORING_BACKPORT_PROGRESS PR_PLAN_CREATION \
  SPAWN_CODE_REVIEWER_PHASE_MERGE_PLAN WAITING_FOR_PHASE_MERGE_PLAN \
  SPAWN_INTEGRATION_AGENT_PHASE MONITORING_PHASE_INTEGRATION \
  SPAWN_CODE_REVIEWER_PHASE_FIX_PLAN \
  WAITING_FOR_PHASE_FIX_PLANS SPAWN_CODE_REVIEWER_FIX_PLAN \
  WAITING_FOR_FIX_PLANS DISTRIBUTE_FIX_PLANS \
  IMMEDIATE_BACKPORT_REQUIRED; do
  
  if [ -f "agent-states/orchestrator/$state/rules.md" ]; then
    echo "✅ $state"
  else
    echo "❌ $state - MISSING!"
  fi
done
```

## R217 Compliance

When transitioning to any state, the orchestrator MUST:
1. Update state file (R288)
2. Read the corresponding rules.md file (R217)
3. Acknowledge the rules before proceeding

Example:
```bash
# Transitioning to WAVE_REVIEW
update_orchestrator_state "WAVE_REVIEW" "Integration complete, requesting review"
# READ: agent-states/orchestrator/WAVE_REVIEW/rules.md
# Then proceed with WAVE_REVIEW work
```