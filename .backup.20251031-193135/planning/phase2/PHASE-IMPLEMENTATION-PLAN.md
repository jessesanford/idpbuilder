# Phase 2 Implementation Plan - Core Push Functionality

**Phase**: Phase 2 - Core Push Functionality
**Created**: 2025-10-31
**Planner**: @agent-code-reviewer
**Fidelity Level**: **WAVE LIST ONLY** (high-level descriptions, no detailed plans)

---

## Phase Overview

**Goal**: Implement the `idpbuilder push` command that orchestrates Phase 1 packages to push Docker images to OCI registries with full CLI integration, progress reporting, and comprehensive error handling.

**Scope**: Command layer implementation, user interaction, configuration management, input validation, and security features. Does NOT include new Docker/Registry/Auth/TLS implementations (those are from Phase 1).

**Dependencies**:
- **Phase 1 (Core Libraries)**: All packages (docker, registry, auth, tls) must be complete
- **IDPBuilder CLI Framework**: Cobra and Viper libraries already in IDPBuilder
- **Test Infrastructure**: Docker daemon and test registry for integration testing

---

## Wave List

This section provides a **high-level roadmap** of waves in this phase. Detailed effort definitions will be created **just-in-time** during wave planning.

### Wave 1: Command Implementation & Core Integration

**Description**: Create the foundational `idpbuilder push` command that integrates all Phase 1 packages into a functional CLI tool with basic progress reporting. This wave establishes the pipeline architecture and demonstrates end-to-end push capability.

**Key Features**:
- Cobra command registration with all required flags (--registry, --username, --password, --insecure, --verbose)
- Pipeline orchestration across Docker → Auth → TLS → Registry stages
- Real-time progress reporting with normal and verbose modes
- Basic error handling and exit code mapping
- Integration with Phase 1's docker, registry, auth, and tls packages

**Efforts**: **TBD** (Will be defined during Wave 1 planning)

**Estimated Complexity**: **Medium** (integrating existing packages, new CLI surface)

**Success Criteria**:
- ✅ Push command successfully registered in IDPBuilder
- ✅ Can push alpine:latest to test Gitea registry
- ✅ Progress updates displayed layer-by-layer
- ✅ Basic error messages shown on failure

---

### Wave 2: Advanced Configuration Features

**Description**: Add registry override capability and environment variable support with proper precedence handling. This wave makes the command production-ready for diverse deployment scenarios and CI/CD integration.

**Key Features**:
- Custom registry URL override (--registry flag)
- Environment variable support (IDPBUILDER_REGISTRY, IDPBUILDER_REGISTRY_USERNAME, etc.)
- Configuration precedence logic (flags > env vars > defaults)
- Viper integration for unified configuration management
- Image reference override for custom registries

**Efforts**: **TBD** (Will be defined during Wave 2 planning)

**Estimated Complexity**: **Low** (mostly configuration plumbing)

**Parallelization**: **YES** (both efforts can be implemented independently)

**Success Criteria**:
- ✅ Can push to custom registries (DockerHub, Quay.io, localhost:5000)
- ✅ All flags bindable to environment variables
- ✅ Precedence correctly applied in all scenarios
- ✅ Help text documents environment variables

---

### Wave 3: Security & Validation

**Description**: Implement comprehensive input validation, command injection prevention, SSRF protection, and production-grade error handling with actionable user feedback. This wave ensures the command is secure and user-friendly.

**Key Features**:
- Input validation (OCI image name format, registry URL format, credential validation)
- Command injection prevention (shell metacharacter filtering)
- SSRF prevention (private IP range detection with warnings)
- Error type system with specific exit codes (1=validation, 2=auth, 3=network, 4=image not found)
- Actionable error messages with suggestions (Error: X, Suggestion: Y format)

**Efforts**: **TBD** (Will be defined during Wave 3 planning)

**Estimated Complexity**: **Medium** (security requires thorough testing)

**Success Criteria**:
- ✅ All dangerous inputs rejected with clear errors
- ✅ Error messages include actionable suggestions
- ✅ Exit codes correctly mapped to error types
- ✅ 95%+ test coverage on validation code

---

## Wave Dependencies

```
Wave 1 (Foundation & Core Integration)
  ↓
Wave 2 (Configuration Features) - Builds on Wave 1 command structure
  ↓
Wave 3 (Security & Validation) - Validates configuration from Wave 2
```

**Critical Path**:
- Wave 1 must complete first (establishes command foundation)
- Wave 2 can parallelize some efforts (registry override vs env vars)
- Wave 3 depends on Wave 2 (validates the full configuration surface)

**Branching Strategy**:
- Wave 1 branches from: Phase 1 integration branch
- Wave 2 branches from: Wave 1 integration branch
- Wave 3 branches from: Wave 2 integration branch

---

## Risk Assessment

### High-Risk Waves

- **Wave 1 (Command Integration)**: **MEDIUM RISK** - First time integrating all Phase 1 packages together. Potential issues: Phase 1 interfaces may not work together as expected, go-containerregistry integration complexities, progress callback threading issues.
  - **Mitigation**: Extensive integration testing with real Docker daemon and registry, mock-based unit tests for pipeline stages, early smoke tests with alpine:latest

- **Wave 3 (Security & Validation)**: **MEDIUM RISK** - Security validation is critical and easy to get wrong. Potential issues: Missing command injection patterns, SSRF bypass techniques, validation regex errors.
  - **Mitigation**: Reuse Phase 1's command injection tests (from docker package), security-focused code review, penetration testing with dangerous inputs, 95%+ coverage requirement

### Dependencies on External Systems

- **Docker Daemon**: Required for Wave 1 integration testing (must be running)
- **Test Gitea Registry**: Required for all integration tests (can use container)
- **IDPBuilder Cobra Framework**: Required for Wave 1 (already in IDPBuilder)
- **Viper Configuration Library**: Required for Wave 2 (already in IDPBuilder)

---

## Test Strategy

### TDD Approach (R341 Compliance)

**CRITICAL**: Phase 2 Test Plan already created (see `planning/phase2/PHASE-2-TEST-PLAN.md`) **BEFORE** implementation.

**Test-First Workflow**:
1. ✅ **NOW**: Test plan exists with 155 defined tests
2. ⏭️ **Wave 1**: Implement to pass command + progress tests (50 tests)
3. ⏭️ **Wave 2**: Implement to pass configuration tests (covered in Wave 1 tests)
4. ⏭️ **Wave 3**: Implement to pass validation + error tests (65 tests)
5. ⏭️ **Integration**: Run full suite (40 integration tests)

**Coverage Targets**:
- Phase 2 Overall: **85%+**
- Command Layer (cmd/push): **90%+**
- Input Validator (pkg/validator): **95%+**
- Progress Reporter (pkg/progress): **85%+**

**Test Categories** (from test plan):
- Unit Tests - Command Layer: 35 tests
- Unit Tests - Progress Reporter: 15 tests
- Unit Tests - Input Validator: 45 tests
- Unit Tests - Error Handling: 20 tests
- Integration Tests - CLI: 25 tests
- Integration Tests - Component: 15 tests

---

## Success Criteria

**Phase Complete When**:
- [ ] All 3 waves integrated and tested
- [ ] 155 tests passing (from test plan)
- [ ] Coverage ≥85% overall, ≥90% for critical components
- [ ] Can push to Gitea, DockerHub, and custom registries
- [ ] All error types tested with correct exit codes
- [ ] Documentation complete (help text, environment variables)
- [ ] Architect phase assessment approved
- [ ] Phase integration tests passing
- [ ] No stub implementations remaining
- [ ] Security validation comprehensive

---

## Progressive Planning Notes

### Why Wave List Only?

This phase implementation plan provides **only a wave list** because:

1. **Adaptive Planning**: Wave 1 implementation may reveal better patterns for Wave 2/3
2. **Fidelity Gradient**: Effort definitions require **real code examples** from Wave 1 integration
3. **Just-In-Time Planning**: Detailed effort plans created when wave starts (not months ahead)
4. **Progressive Refinement**: Each wave's planning benefits from previous wave's lessons

### When Are Effort Plans Created?

**Effort plans** are created **just-in-time** at **wave start**:

1. **Orchestrator transitions to WAVE_START** (e.g., start Wave 1)
2. **Architect creates wave architecture plan** (`WAVE-1-ARCHITECTURE.md`) with **real code examples**
3. **Code Reviewer creates wave implementation plan** (`WAVE-1-IMPLEMENTATION.md`) with **detailed efforts**:
   - Exact file lists (e.g., cmd/push.go, cmd/push_test.go)
   - Real code specifications (actual function signatures)
   - R213 metadata (effort_id, estimated_lines, dependencies, branch_name, can_parallelize)
   - Task breakdowns (step-by-step implementation)

**Example**: Wave 1 Implementation Plan will specify:
- Effort 2.1.1: Push Command Core (~450 lines) - cmd/push.go, pipeline orchestration, flag definitions
- Effort 2.1.2: Progress Reporter (~300 lines) - pkg/progress/reporter.go, callback handling, console output

---

## Lessons from Phase 1

### What Worked Well (Carry Forward)
- **Interface-first approach**: Phase 1's stable interfaces enable easy integration
- **Mock-based testing**: Will reuse mock providers from Phase 1 tests
- **Package separation**: Phase 1's clean boundaries make command layer simple
- **Go standard patterns**: Context, defer, error wrapping worked great

### What We Learned
- **go-containerregistry is straightforward**: Integration should be smooth
- **Basic auth handles special characters**: No escaping needed (from Phase 1 testing)
- **InsecureSkipVerify works as expected**: Wave 1 can use directly
- **Phase 1 coverage is 85%+**: Same bar for Phase 2

### Patterns to Continue
- **Given/When/Then test structure**: Proven in Phase 1, will reuse
- **Table-driven tests**: Effective for flag/validation testing
- **Error type checking with ErrorAs()**: Will use for Phase 2 error types
- **httptest.NewServer() for mocking**: Will mock registry endpoints

---

## Integration with Phase 1

### Phase 1 Packages Available

**All Phase 1 packages are complete and tested**:

1. **pkg/docker** (DockerClient interface)
   - `GetImage(ctx, imageName) -> (v1.Image, error)`
   - Already handles command injection prevention
   - Tests: 31 test cases, 85%+ coverage

2. **pkg/registry** (RegistryClient interface)
   - `Push(ctx, image, targetRef, progressCallback) -> error`
   - Progress callback pattern already defined
   - Tests: 31 test cases, integration with go-containerregistry verified

3. **pkg/auth** (AuthProvider interface)
   - `NewBasicAuthProvider(username, password) -> AuthProvider`
   - Special character support verified (unicode, quotes, spaces)
   - Tests: 31 test cases, special character validation complete

4. **pkg/tls** (TLSProvider interface)
   - `NewConfigProvider(insecure bool) -> TLSProvider`
   - Both secure and insecure modes tested
   - Tests: 10 test cases, HTTP client integration verified

### Phase 1 Test Fixtures for Reuse

**Mock Providers** (from pkg/registry/client_test.go):
- `mockAuthProvider` - Reuse in Wave 1 pipeline tests
- `mockTLSProvider` - Reuse in Wave 1 pipeline tests

**Test Images**:
- `alpine:latest` - Phase 1 prerequisite, available for all integration tests
- `v1/empty.Image` - go-containerregistry's empty image for unit tests

**Test Patterns**:
- Table-driven tests (pkg/docker, pkg/auth)
- Command injection prevention (pkg/docker)
- Error type testing with `assert.ErrorAs()`
- httptest.NewServer() for registry mocking

---

## Next Steps

1. **Phase Architecture Review**: Architect reviews `PHASE-2-ARCHITECTURE.md` (already complete)
2. **Phase Test Plan Review**: Review `PHASE-2-TEST-PLAN.md` (already complete, 155 tests defined)
3. **Orchestrator Approval**: Orchestrator validates phase plan structure
4. **Wave 1 Start**: When approved, orchestrator proceeds to Wave 1
5. **Wave 1 Architecture**: Architect creates `planning/phase2/wave1/WAVE-1-ARCHITECTURE.md` with **real Go code**
6. **Wave 1 Implementation**: Code Reviewer creates `planning/phase2/wave1/WAVE-1-IMPLEMENTATION.md` with **exact effort specifications** including R213 metadata

**Note**: This document is intentionally **high-level**. Detailed planning happens progressively during wave execution following SF 3.0's fidelity gradient (Phase = waves only, Wave = efforts with real code).

---

## R502 Quality Gates (Phase Implementation)

### Phase Implementation Plan Requirements

- ✅ **Wave list only**: No effort definitions (correct for phase level)
- ✅ **High-level descriptions**: Each wave has 2-3 sentence summary
- ✅ **Wave dependencies documented**: Explicit sequence and parallelization notes
- ✅ **No R213 metadata**: Metadata comes in wave implementation plans (not phase)
- ✅ **No file paths**: File specifications come in wave implementation plans
- ✅ **Success criteria defined**: Clear phase completion criteria
- ✅ **3-8 waves**: 3 waves (appropriate scope)

### What This Plan DOES NOT Include (Intentionally)

- ❌ **Effort definitions**: Will be created during wave planning
- ❌ **R213 metadata**: Added in wave implementation plans
- ❌ **File paths**: Specified in wave implementation plans
- ❌ **Real code examples**: Provided in wave architecture plans
- ❌ **Detailed task breakdowns**: Created just-in-time per effort

### Fidelity Verification

**This is a PHASE Implementation Plan (HIGH-LEVEL)**:
- Uses wave names and descriptions
- No effort specifications
- No R213 metadata blocks
- No file lists
- **Correct fidelity level per R502**

---

## Compliance Checklist

### R502 Quality Gates
- ✅ Phase-level fidelity (wave list only)
- ✅ 3 waves defined (appropriate scope)
- ✅ Wave descriptions are high-level (2-3 sentences each)
- ✅ No effort definitions (correct - comes later)
- ✅ No R213 metadata (correct - wave level only)

### R341 TDD Compliance
- ✅ Test plan created BEFORE implementation
- ✅ 155 tests defined in PHASE-2-TEST-PLAN.md
- ✅ Coverage targets specified (85%+ overall)
- ✅ Test-first workflow documented

### R308 Incremental Branching
- ✅ Wave 1 branches from Phase 1 integration
- ✅ Wave 2 branches from Wave 1 integration
- ✅ Wave 3 branches from Wave 2 integration
- ✅ No parallel phase development

### R510 Checklist Structure
- ✅ All sections have clear criteria
- ✅ Success criteria checkboxes included
- ✅ Quality gates verified
- ✅ Compliance checklist present

---

## Document Status

**Status**: ✅ READY FOR ORCHESTRATOR HANDOFF
**Planner**: @agent-code-reviewer
**Created**: 2025-10-31
**Waves**: 3 (Command Integration, Configuration, Security)
**Fidelity Level**: HIGH-LEVEL (wave list only, no efforts)

**Next Step**:
- Orchestrator validates phase plan structure
- Orchestrator transitions to Phase 2, Wave 1 start
- Architect creates Wave 1 architecture plan (real code examples)
- Code Reviewer creates Wave 1 implementation plan (detailed efforts with R213)

**Compliance Verified**:
- ✅ R502: Phase-level fidelity (wave list only)
- ✅ R341: TDD compliance (test plan exists)
- ✅ R510: Checklist structure followed
- ✅ R308: Incremental branching defined

---

**END OF PHASE 2 IMPLEMENTATION PLAN**
