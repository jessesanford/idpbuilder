# Code Reviewer Assignment: Create Backport Plan

## Assignment Overview
You are being spawned to create a comprehensive BACKPORT-PLAN.md that will guide SW Engineers in implementing the missing authentication constructor functions that were identified during TDD GREEN phase testing.

## Your Role
As Code Reviewer, you will:
1. Analyze the integration test failures
2. Map required fixes to specific effort branches
3. Create detailed implementation instructions
4. Define verification criteria

## Current Situation
- **Phase 2 Wave 1 Integration**: Complete but tests failing (expected for TDD)
- **Issue**: Authentication constructor functions are undefined
- **Resolution**: Need to backport implementations to source branches per R321

## Input Documentation
1. **FIX-MANIFEST-FOR-BACKPORT.md** - Contains details of missing functions
   - Location: `/home/vscode/workspaces/idpbuilder-push/efforts/phase2/wave1/integration-workspace/FIX-MANIFEST-FOR-BACKPORT.md`
   - Content: Test failures, missing functions, affected branches

2. **INTEGRATION-REPORT.md** - Shows integration status
   - Location: `/home/vscode/workspaces/idpbuilder-push/efforts/phase2/wave1/integration-workspace/.software-factory/INTEGRATION-REPORT.md`
   - Content: Build success, test failures detail

3. **Integration Branch State**
   - Branch: `idpbuilderpush/phase2/wave1/integration`
   - Tests exist but constructors undefined

## Required Output: BACKPORT-PLAN.md

Your plan MUST include:

### 1. Executive Summary
- Brief overview of what needs to be backported
- Why these changes are required
- Expected outcome after implementation

### 2. Effort-Specific Instructions

For **Effort 2.1.2 (auth-implementation)**:
- **Branch**: `idpbuilderpush/phase2/wave1/auth-implementation`
- **Working Directory**: `/home/vscode/workspaces/idpbuilder-push/efforts/phase2/wave1/auth-implementation`
- **File to Modify**: `pkg/oci/auth.go`
- **Functions to Add**:
  - `NewAuthenticatorFromSecrets(secretData map[string][]byte) (*Authenticator, error)`
  - `NewAuthenticatorFromFlags(username, password string) (*Authenticator, error)`
  - `NewAuthenticatorFromEnv() (*Authenticator, error)`
- **Implementation Guidelines**: Based on test expectations
- **Size Constraints**: Keep within effort limits

### 3. Implementation Sequence
1. Order of implementation (if dependencies exist)
2. Which functions to implement first
3. Any prerequisites

### 4. Verification Steps
For each backport:
- Build command: `go build ./pkg/oci/...`
- Test command: `go test ./pkg/oci/...`
- Expected results after implementation
- How to verify no regressions

### 5. Success Criteria
- All constructor functions implemented
- Build succeeds
- Tests progress from "undefined" errors
- No existing code deleted (R359)
- Implementation within size limits (R304)

### 6. Risk Assessment
- Potential issues
- Mitigation strategies
- Rollback plan if needed

## Constraints and Requirements

### Mandatory Rules
- **R321**: Fixes must go to source branches, not integration
- **R359**: No deletion of existing code
- **R304**: Use line counter to verify size compliance
- **R006**: Don't write code yourself, only planning

### Quality Standards
- Clear, unambiguous instructions
- Exact file paths and function signatures
- Verification steps that can be automated

## Working Directory
Start in: `/home/vscode/workspaces/idpbuilder-push/efforts/phase2/wave1/integration-workspace`

## State Management
Create and maintain: `code-reviewer-state.json` with:
```json
{
  "current_state": "BACKPORT_PLAN_CREATION",
  "plan_status": "in_progress",
  "efforts_analyzed": [],
  "plan_location": ""
}
```

When complete, update to:
```json
{
  "current_state": "BACKPORT_PLAN_COMPLETE",
  "plan_status": "complete",
  "efforts_analyzed": ["2.1.2"],
  "plan_location": "efforts/phase2/wave1/integration-workspace/BACKPORT-PLAN.md"
}
```

## Deliverables Checklist
- [ ] BACKPORT-PLAN.md created
- [ ] All missing functions mapped to branches
- [ ] Implementation instructions clear and complete
- [ ] Verification steps defined
- [ ] State file updated
- [ ] Plan ready for SW Engineer consumption

## Timeline
- Start immediately upon spawn
- Complete analysis within reasonable time
- Focus on clarity over speed

## Success Metrics
Your plan will be successful if:
1. SW Engineers can implement without ambiguity
2. All test failures are addressed
3. No questions about what goes where
4. Verification is straightforward
5. Complies with all Software Factory rules

## Notes
- This is TDD GREEN phase - tests exist, implementation needed
- Constructor functions are the primary gap
- Keep instructions specific to auth-implementation effort
- Remember R321 - we're fixing source, not integration

---
*Assigned by Orchestrator*
*State: SPAWN_CODE_REVIEWER_BACKPORT_PLAN*
*Date: 2025-09-24*