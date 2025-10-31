# Phase 1 Implementation Plan
# IDPBuilder OCI Push Command - Foundation & Interfaces

**Phase**: Phase 1 - Foundation & Interfaces
**Created**: 2025-10-29
**Planner**: Code Reviewer Agent
**Fidelity Level**: **WAVE LIST ONLY** (high-level descriptions, no detailed plans)

---

## Phase Overview

**Goal**: Establish all interfaces and core package implementations to enable parallel development in subsequent waves and phases.

**Scope**:
- Define clean contracts between Docker, Registry, Auth, and TLS components
- Implement foundational packages using go-containerregistry library
- Create interface-first architecture for maximum parallelization
- Prepare command structure skeleton for Phase 2 integration

**Dependencies**:
- Existing IDPBuilder cobra framework
- Docker daemon running locally
- go-containerregistry library (to be added)
- Docker Engine API client library (to be added)

**Phase Outcomes**:
- All package interfaces documented and compiled
- Four core packages fully implemented with unit tests
- Unit test coverage ≥85%
- Foundation ready for Phase 2 command integration
- No breaking changes to existing IDPBuilder code

---

## Wave List

This section provides a **high-level roadmap** of waves in Phase 1. Detailed effort definitions will be created **just-in-time** during wave planning.

### Wave 1: Interface & Contract Definitions

**Description**: Define ALL interfaces upfront to freeze contracts and enable Phase 1 Wave 2 parallel implementation. This wave establishes the complete API surface for Docker operations, registry interactions, authentication, and TLS configuration.

**Key Features**:
- Docker client interface (ImageExists, GetImage, ValidateImageName, Close)
- Registry client interface (Push, BuildImageReference, ValidateRegistry)
- Authentication provider interface (GetAuthenticator, ValidateCredentials)
- TLS configuration provider interface (GetTLSConfig, IsInsecure)
- Progress reporting types (ProgressCallback, ProgressUpdate)
- Command structure skeleton with flag definitions
- Comprehensive error type definitions
- All GoDoc documentation

**Efforts**: **TBD** (Will be defined during Wave 1 planning - expect ~4 interface definition efforts)

**Estimated Complexity**: **Low**
- Pure interface definitions, no implementations
- Approximately 650 total lines across all efforts
- Sequential development (interfaces must be coordinated)
- Build guaranteed green (interfaces compile)

**Dependencies**: None (first wave of phase)

**Outcomes**:
- All interfaces frozen and documented
- go build succeeds with all interfaces
- Contracts ready for parallel Wave 2 implementation
- No interface changes allowed after Wave 1 integration

---

### Wave 2: Core Package Implementations

**Description**: Implement all four core packages in parallel using frozen Wave 1 interfaces. Teams work independently on Docker client, Registry client, Authentication, and TLS configuration with no coordination needed during implementation.

**Key Features**:
- Docker client implementation using Docker Engine API
- Registry client implementation using go-containerregistry
- Basic authentication implementation with special character support
- TLS configuration with insecure mode support
- Comprehensive unit tests for all packages (85%+ coverage)
- Error handling and validation logic
- Integration with go-containerregistry's authn package
- System certificate pool fallback for secure mode

**Efforts**: **TBD** (Will be defined during Wave 2 planning - expect 4 parallel implementation efforts)

**Estimated Complexity**: **Medium**
- Approximately 1,550 total lines across all efforts
- **FULLY PARALLELIZABLE** (4 independent teams)
- Each effort implements frozen interface
- No cross-effort dependencies
- Unit tests with mocked external dependencies

**Dependencies**:
- Phase 1 Wave 1 integration (interfaces frozen)
- go-containerregistry library added to go.mod
- Docker Engine API client library added to go.mod

**Outcomes**:
- All interfaces implemented
- Docker client handles daemon connection and image retrieval
- Registry client integrates with go-containerregistry for push operations
- Authentication supports username/password with special characters
- TLS configuration supports both secure and insecure modes
- Unit test coverage ≥85% across all packages
- All packages compile and integrate together
- Ready for Phase 2 command integration

---

## Wave Dependencies

```
Wave 1 (Interface Definitions)
  ↓
Wave 2 (Parallel Implementations)
  ├─ Docker Client Implementation
  ├─ Registry Client Implementation
  ├─ Authentication Implementation
  └─ TLS Configuration Implementation
```

**Critical Path**:
- Wave 1 must complete before Wave 2 (interfaces freeze contracts)
- Wave 2: All 4 efforts can proceed in parallel (maximum parallelization)
- No sequential dependencies within Wave 2

**Parallelization Strategy**:
- Wave 1: Sequential (4 efforts coordinating interface design)
- Wave 2: **PARALLEL** (4 independent implementation teams)

---

## Phase Architecture Principles

### Interface-First Design (R307 Compliance)

**Why Wave 1 Defines All Interfaces:**
1. Enables maximum parallelization in Wave 2 (4 teams simultaneously)
2. Ensures independent branch mergeability (R307)
3. Prevents interface changes mid-implementation
4. Clear contracts between all components
5. Build always green (interfaces compile before implementations)

### Modular Package Structure

**Separation of Concerns:**
- `pkg/docker/` - Docker daemon integration (isolated)
- `pkg/registry/` - OCI registry operations (isolated)
- `pkg/auth/` - Authentication provision (isolated)
- `pkg/tls/` - TLS configuration (isolated)
- `cmd/push.go` - Command layer skeleton (Phase 2 completion)

### Incremental Branching (R308)

**Branch Strategy:**
- Wave 2 branches from Wave 1 integration
- Phase 2 will branch from Phase 1 integration
- Each wave adds functionality incrementally
- No "big bang" integration at end

---

## Risk Assessment

### High-Risk Waves

**Wave 1 - Interface Design:**
- **Risk**: Interface changes after Wave 2 starts would break parallelization
- **Mitigation**: Architect review of all interfaces before Wave 2 spawns
- **Mitigation**: Comprehensive GoDoc to clarify intent
- **Mitigation**: Consider integration points during design

**Wave 2 - Parallel Implementation:**
- **Risk**: Implementations diverge from interface contracts
- **Mitigation**: Unit tests verify interface compliance
- **Mitigation**: Integration testing at wave boundary
- **Mitigation**: Code review focus on contract adherence

### Dependencies on External Systems

**Docker Daemon:**
- Required for Wave 2 Docker client implementation
- Required for integration testing
- Mitigation: Mocked for unit tests, real daemon for integration tests

**go-containerregistry Library:**
- Required for Wave 2 Registry client implementation
- Pin to stable version (v0.19.0+)
- Monitor for breaking changes

**Gitea Registry:**
- NOT required for Phase 1 (testing deferred to Phase 3)
- Wave 2 uses mocked registry responses

---

## Success Criteria

**Phase 1 Complete When**:
- [ ] All interfaces defined and documented (Wave 1)
- [ ] All four packages implemented with unit tests (Wave 2)
- [ ] Unit test coverage ≥85% across all packages
- [ ] go build succeeds for all new code
- [ ] All efforts within size limits (<700 lines each)
- [ ] Phase 1 integration tests passing (packages work together)
- [ ] Metadata organized per R383 in .software-factory/
- [ ] Architect phase assessment approved

**Build Verification**:
- [ ] `go build ./pkg/docker` succeeds
- [ ] `go build ./pkg/registry` succeeds
- [ ] `go build ./pkg/auth` succeeds
- [ ] `go build ./pkg/tls` succeeds
- [ ] `go build ./cmd` succeeds (skeleton only)
- [ ] `go test ./...` all tests pass

**Test Coverage Targets**:
- [ ] pkg/docker: ≥85% coverage
- [ ] pkg/registry: ≥85% coverage
- [ ] pkg/auth: ≥90% coverage (critical security)
- [ ] pkg/tls: ≥90% coverage (critical security)

---

## Progressive Planning Notes

### Why Wave List Only?

This phase implementation plan provides **only a wave list** because:

1. **Adaptive Planning**: Detailed effort plans must adapt based on what we learn during implementation
2. **Fidelity Gradient**: Effort definitions require **real code examples** from the architecture plan (already provided in PHASE-1-ARCHITECTURE.md)
3. **Just-In-Time Planning**: Creating detailed plans upfront wastes effort when requirements change
4. **Progressive Refinement**: Each wave's planning benefits from architect's pseudocode guidance

### When Are Effort Plans Created?

**Effort plans** are created **just-in-time** at **wave start**:

1. Orchestrator transitions to **WAVE_START**
2. Architect has already created **PHASE-1-ARCHITECTURE.md** with detailed pseudocode
3. Code Reviewer creates **WAVE-N-IMPLEMENTATION.md** with detailed effort definitions
4. Efforts include:
   - Exact file lists
   - Concrete specifications based on architecture pseudocode
   - R213 metadata
   - Dependencies and parallelization info
   - Size estimates

### Architecture Guidance Available

**Phase 1 has comprehensive pseudocode in PHASE-1-ARCHITECTURE.md:**
- Complete interface pseudocode for all contracts
- Implementation pseudocode for all packages
- Data flow patterns
- Error handling strategies
- Testing approaches
- All design decisions documented

**This provides strong foundation for wave planning!**

---

## Sizing Estimates

**Conservative Effort Sizing (Safety Margin):**

| Wave | Total Lines | Efforts | Avg per Effort | Safety Margin |
|------|------------|---------|----------------|---------------|
| Wave 1 | ~650 | 4 | ~160 | 640 lines below limit |
| Wave 2 | ~1,550 | 4 | ~390 | 410 lines below limit |
| **Phase Total** | **~2,200** | **8** | **~275** | **Well under limits** |

**All efforts designed with 200-500 line buffer below 800 hard limit**

---

## Next Steps

1. **Wave 1 Start**: When orchestrator approves, proceed to Wave 1
2. **Wave 1 Planning**: Code Reviewer creates `WAVE-1-IMPLEMENTATION.md` with exact effort specs
3. **Wave 1 Execution**: Spawn 4 SW Engineers for interface definition efforts (sequential)
4. **Wave 1 Integration**: Create `phase1-wave1-integration` branch
5. **Wave 2 Planning**: Code Reviewer creates `WAVE-2-IMPLEMENTATION.md` with exact effort specs
6. **Wave 2 Execution**: Spawn 4 SW Engineers in PARALLEL for implementation efforts
7. **Wave 2 Integration**: Create `phase1-wave2-integration` branch
8. **Phase 1 Integration**: Merge all waves, create `phase1-integration` branch
9. **Architect Review**: Phase assessment before Phase 2

**Note**: This document is intentionally **high-level**. Detailed planning happens progressively during wave execution, guided by comprehensive pseudocode in PHASE-1-ARCHITECTURE.md.

---

## Compliance Summary

**Software Factory 3.0 Rules:**

| Rule | Description | Compliance Strategy |
|------|-------------|---------------------|
| R307 | Independent Branch Mergeability | Interface-first design enables parallel Wave 2 |
| R308 | Incremental Branching Strategy | Wave 2 from Wave 1, Phase 2 from Phase 1 |
| R359 | No Code Deletion | Pure additive enhancement (no deletions) |
| R383 | Metadata File Organization | All metadata in .software-factory/ with timestamps |
| R220/R221 | Size Limits | All efforts 150-550 lines (safety margin) |
| R502 | Wave List Fidelity | This document = wave list only (correct level) |

**All compliance verified at wave planning time.**

---

## Document Status

**Status**: ✅ READY FOR WAVE PLANNING
**Phase**: Phase 1 of 3
**Fidelity**: WAVE LIST (R502 compliant)
**Created By**: Code Reviewer Agent
**Date**: 2025-10-29

**Next Action**: Orchestrator proceeds to Wave 1 start
- Architect has already provided PHASE-1-ARCHITECTURE.md (pseudocode)
- Code Reviewer will create WAVE-1-IMPLEMENTATION.md (detailed efforts)
- SW Engineers execute efforts

---

**END OF PHASE 1 IMPLEMENTATION PLAN**
