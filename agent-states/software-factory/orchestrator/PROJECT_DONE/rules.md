# PROJECT_DONE State Rules

## 🔴🔴🔴 ABSOLUTE REQUIREMENT: STATE MANAGER CONSULTATION 🔴🔴🔴

**SUPREME LAW - R517 - UNIVERSAL STATE MANAGER CONSULTATION LAW**

**BEFORE EXITING THIS STATE, YOU MUST:**

1. **SPAWN STATE MANAGER FOR SHUTDOWN_CONSULTATION** (MANDATORY - NO EXCEPTIONS)
2. **NEVER UPDATE STATE FILES DIRECTLY** (orchestrator-state-v3.json, bug-tracking.json, etc.)
3. **NEVER COMMIT STATE FILES YOURSELF** (State Manager does this atomically)

**FAILURE TO CONSULT STATE MANAGER = IMMEDIATE SYSTEM HALT (-100% GRADE)**

### Enforcement Mechanism

If you attempt to exit this state without spawning State Manager:
- ❌ Pre-commit hooks will REJECT your commit
- ❌ Validation tools will FAIL the build
- ❌ Grading system will assign -100% penalty
- ❌ System will transition to ERROR_RECOVERY

### Required Pattern (COPY THIS EXACTLY)

```bash
# At end of state work, BEFORE any state file updates:

echo "🔴 MANDATORY: Spawning State Manager for SHUTDOWN_CONSULTATION"

# Spawn State Manager (REQUIRED - NOT OPTIONAL)
# Task: state-manager
# State: SHUTDOWN_CONSULTATION
# Current State: [YOUR_CURRENT_STATE]
# Proposed Next State: [YOUR_PROPOSED_NEXT_STATE]
# Work Summary: [SUMMARY_OF_WORK_COMPLETED]

# State Manager will:
# 1. Validate proposed transition against state machine
# 2. Update all 4 state files atomically
# 3. Commit with [R288] tag
# 4. Return REQUIRED next state (may differ from proposal)

# Wait for State Manager response
# Follow State Manager's directive (REQUIRED next state)
# DO NOT proceed until State Manager confirms
```

**YOU MUST NEVER:**
- ❌ Update orchestrator-state-v3.json yourself
- ❌ Update bug-tracking.json yourself
- ❌ Update integration-containers.json yourself
- ❌ Use `jq` to modify state files
- ❌ Use `sed/awk` to modify state files
- ❌ Set `validated_by: "orchestrator"` (must be "state-manager")
- ❌ Skip State Manager consultation "just this once"
- ❌ Think "I'll validate it manually"

**ONLY State Manager may update state files. This is NON-NEGOTIABLE.**

See: `rule-library/R517-universal-state-manager-consultation-law.md`
## 🔴🔴🔴 SF 3.0 TERMINAL STATE 🔴🔴🔴

**THIS IS A TERMINAL STATE - NO FURTHER WORK REQUIRED**

### STATE BEHAVIOR (CRITICAL)
PROJECT_DONE is a **terminal state** where the Software Factory 3.0 project is complete and ready for human PR submission. There is NO work to perform in this state - it exists solely as a completion marker.

**When orchestrator enters this state:**
1. **DO NOT** attempt to perform any actions
2. **DO NOT** spawn any agents
3. **IMMEDIATELY** output CONTINUE-SOFTWARE-FACTORY=FALSE and exit

---

## State Overview
**State**: PROJECT_DONE
**Type**: terminal
**Description**: Software Factory 3.0 project complete - awaiting human PR submission per R280

## Entry Conditions
- PR_PLAN_CREATION completed successfully
- MASTER-PR-PLAN.md exists with complete PR instructions
- All integration containers converged
- Project status = "complete"
- **Production readiness validated** (see Production Readiness Checklist below)

## Exit Conditions
**NONE** - This is a terminal state with no outbound transitions

## Required Actions

### ✅ ONLY ACTION: Immediate Exit with FALSE Flag

```bash
# This is a TERMINAL state - no work to perform
# Simply acknowledge completion and exit

echo "## 🎉 PROJECT COMPLETE"
echo ""
echo "**Software Factory 3.0 project has reached terminal state PROJECT_DONE**"
echo ""
echo "**Status**: All work complete, PR plan ready for human operator"
echo "**Next Step**: Human operator should review MASTER-PR-PLAN.md and create pull requests"
echo ""
echo "Per R280, SF 3.0 never merges to main automatically - human review required."
echo ""
echo "---"
echo ""
echo "**CONTINUE-SOFTWARE-FACTORY=FALSE**"
echo ""
echo "🛑 **Project complete - no further automation**"

# Exit immediately
exit 0
```

## Valid Transitions
**NONE** - Terminal state has no outbound transitions (allowed_transitions: [])

## Critical Rules
- **R405**: CONTINUE-SOFTWARE-FACTORY=FALSE for terminal states
- **R280**: SF 3.0 never merges to main - human PR submission required
- **R322**: Mandatory stop (exit 0)

## SF 3.0 State Management

**State Manager Behavior**:
- When State Manager transitions TO this state, it sets `continue_flag = "FALSE"`
- State Manager's decision is authoritative - orchestrator must respect it
- No further SHUTDOWN_CONSULTATION needed - this is the end

**State File Status**:
- orchestrator-state-v3.json: current_state = "PROJECT_DONE"
- project_progression.project_status = "complete"
- All integration containers: status = "converged"

## Implementation Notes

### What NOT to Do (Common Mistakes)
❌ DO NOT spawn State Manager for transition (no transitions available)
❌ DO NOT attempt to "complete" work (work is already complete)
❌ DO NOT save TODOs (optional - state is terminal)
❌ DO NOT output CONTINUE-SOFTWARE-FACTORY=TRUE (violates terminal state semantics)

### Correct Behavior
✅ Acknowledge project completion
✅ Reference MASTER-PR-PLAN.md for human operator
✅ Output CONTINUE-SOFTWARE-FACTORY=FALSE
✅ Exit immediately (exit 0)

## Monitoring Requirements
**NONE** - Terminal state requires no monitoring

## Error Handling
**NOT APPLICABLE** - No work to fail

## TODO Persistence
**OPTIONAL** - Terminal state, no ongoing work to track

---

## 🚨 CRITICAL: AUTO-STOP DETECTION 🚨

**For test framework:**
When a test harness detects current_state = PROJECT_DONE and expected_final_state = PROJECT_DONE:
- The auto-stop check should trigger BEFORE the orchestrator continuation
- Framework should STOP execution and proceed to validation
- This prevents unnecessary API calls and work in a terminal state

**Test validation should verify:**
1. Final state = PROJECT_DONE
2. MASTER-PR-PLAN.md exists
3. All integration containers converged
4. Project status = complete
5. Continuation flag = FALSE

---

## Production Readiness Checklist

**NOTE**: This validation should be performed by PR_PLAN_CREATION state BEFORE transitioning to PROJECT_DONE. PROJECT_DONE is terminal - all validation must be complete by entry.

### Required Validations (All MUST Pass)

#### 1. Test Coverage & Quality
- ✅ All test suites passing (unit, integration, e2e)
- ✅ Code coverage meets project standards (typically ≥80%)
- ✅ No flaky or skipped tests without documented justification
- ✅ Performance tests passing (if applicable)
- ✅ Security tests passing (if applicable)

#### 2. Build & Artifact Quality
- ✅ Final artifact build completes successfully (R323)
- ✅ Build artifact exists and is documented (R323)
- ✅ Artifact has been tested (not just source code) (R323)
- ✅ No build warnings or errors
- ✅ All dependencies resolved correctly
- ✅ Build is reproducible

#### 3. Code Quality & Standards
- ✅ All code reviews completed and approved
- ✅ No critical or high-severity code quality issues
- ✅ Code style guidelines followed
- ✅ No TODO/FIXME comments without tracking issues
- ✅ Documentation updated for all changes
- ✅ API documentation current (if applicable)

#### 4. Dependency Security & Health
- ✅ All dependencies are secure (no known vulnerabilities)
- ✅ Dependencies are up-to-date (or documented as intentionally pinned)
- ✅ License compatibility verified
- ✅ No deprecated dependencies
- ✅ Dependency resolution clean (no conflicts)

#### 5. Integration & Compatibility
- ✅ All integration containers converged successfully
- ✅ All wave integrations complete
- ✅ All phase integrations complete
- ✅ No merge conflicts or integration errors
- ✅ Backward compatibility maintained (or breaking changes documented)
- ✅ Cross-browser/platform testing complete (if applicable)

#### 6. Documentation & Communication
- ✅ README.md updated and accurate
- ✅ CHANGELOG.md includes all changes
- ✅ Migration guide provided (if breaking changes)
- ✅ API documentation complete
- ✅ Configuration documentation current
- ✅ MASTER-PR-PLAN.md exists with complete instructions

#### 7. Deployment Readiness
- ✅ Environment configuration validated
- ✅ Database migrations tested (if applicable)
- ✅ Rollback procedure documented
- ✅ Monitoring/alerting configured (if applicable)
- ✅ Performance benchmarks within acceptable range
- ✅ Resource requirements documented

#### 8. Compliance & Security
- ✅ Security scan completed (no critical issues)
- ✅ Compliance requirements met (if applicable)
- ✅ Secrets/credentials properly managed
- ✅ Access controls verified
- ✅ Audit trail complete

### Validation Workflow

**In PR_PLAN_CREATION state (BEFORE PROJECT_DONE):**

1. **Spawn Code Reviewer for Production Validation**
   ```bash
   # Code Reviewer validates ALL checklist items above
   # Produces: PRODUCTION-READINESS-REPORT.md
   ```

2. **Review Production Readiness Report**
   ```bash
   # Check if all items passed
   # If ANY item fails → ERROR_RECOVERY (fix issues)
   # If ALL items pass → PROJECT_DONE (terminal)
   ```

3. **Document Validation in orchestrator-state-v3.json**
   ```json
   {
     "project_progression": {
       "project_status": "complete",
       "production_ready": true,
       "production_validation_date": "2025-10-22T00:45:00Z",
       "production_validation_report": "PRODUCTION-READINESS-REPORT.md"
     }
   }
   ```

### SF 2.0 Preservation Notes

This production readiness checklist preserves functionality from SF 2.0's BUILD_VALIDATION state, which was a dedicated validation gate before project completion.

**SF 3.0 Adaptation**:
- Validation performed in PR_PLAN_CREATION (not separate state)
- PROJECT_DONE entry requires production_ready = true
- Code Reviewer performs comprehensive validation
- All 8 validation categories preserved from SF 2.0

**Why Terminal State Model**:
- PROJECT_DONE has no work to perform (terminal by design)
- All validation must complete BEFORE entry
- Prevents orchestrator from "doing work" in terminal state
- Maintains clear separation: validation (PR_PLAN_CREATION) → completion (PROJECT_DONE)

---

## Summary

**PROJECT_DONE is not a state to "execute" - it's a completion marker.**

The orchestrator's role is simply to acknowledge completion and exit with FALSE flag, allowing the test framework or production system to recognize the project is complete and ready for human handoff.

**Production readiness is validated BEFORE entering PROJECT_DONE** - this state confirms all validation is complete.
