# IDPBUILDER PUSH COMMAND - IMPLEMENTATION PLAN

## Project Overview
Implement a `push` command for idpbuilder CLI to upload OCI images to Gitea registry using TEST-DRIVEN DEVELOPMENT methodology. The implementation will strictly follow the RED-GREEN-REFACTOR cycle with tests written FIRST for every feature.

### Core Objectives
- Add `idpbuilder push` command for OCI image uploads to Gitea registry
- Implement authentication via `--username` and `--password` flags
- Support `--insecure` flag for self-signed certificates
- Integrate with existing `get secrets` command for default credentials
- Use google/go-containerregistry library for OCI operations
- Maintain 80%+ test coverage with TDD approach

### Success Criteria
✅ Every feature has tests written FIRST (RED phase)
✅ Minimal implementation to pass tests (GREEN phase)
✅ Code refactored for quality (REFACTOR phase)
✅ 80%+ test coverage on all new code
✅ Seamless integration with existing idpbuilder
✅ OCI standards compliance (no Docker/Gitea specific code)

## Technology Stack
- **Language**: Go 1.22+
- **CLI Framework**: Cobra (existing)
- **OCI Library**: google/go-containerregistry v0.20.2
- **Testing**: stretchr/testify, gomock
- **Logging**: logr (existing)

## Architecture Design

### Component Structure
```
idpbuilder/
├── cmd/
│   └── push/
│       ├── root.go           # Push command implementation
│       └── root_test.go      # Command tests (WRITE FIRST)
└── pkg/
    └── oci/
        ├── client.go         # OCI registry client
        ├── client_test.go    # Client tests (WRITE FIRST)
        ├── auth.go           # Authentication handling
        ├── auth_test.go      # Auth tests (WRITE FIRST)
        ├── push.go           # Push operations
        ├── push_test.go      # Push tests (WRITE FIRST)
        ├── transport.go      # Transport configuration
        └── transport_test.go # Transport tests (WRITE FIRST)
```

### Key Interfaces
```go
// RegistryClient - OCI registry operations contract
type RegistryClient interface {
    Push(ctx context.Context, image, destination string, opts ...Option) error
    Authenticate(username, password string) error
    SetInsecure(insecure bool)
}

// Authenticator - Credential management contract
type Authenticator interface {
    GetCredentials(ctx context.Context) (*Credentials, error)
    ValidateCredentials(username, password string) error
}
```

## PHASE 1: FOUNDATION & COMMAND STRUCTURE
**Target: 800 LOC | Duration: 2-3 days**

### Wave 1.1: Command Skeleton with TDD
**Efforts: 3 | Total Lines: ~500**

#### Effort 1.1.1: Write Command Tests (150 LOC)
- **RED Phase**: Create comprehensive command tests
  - Test command registration with Cobra
  - Test flag parsing (--username, --password, --insecure)
  - Test argument validation (image, registry)
  - Test help text and usage documentation
- **Files**: `cmd/push/root_test.go`

#### Effort 1.1.2: Implement Command Skeleton (200 LOC)
- **GREEN Phase**: Minimal implementation to pass tests
  - Create push command structure
  - Register with root command
  - Define CLI flags
  - Basic argument parsing
- **Files**: `cmd/push/root.go`

#### Effort 1.1.3: Integration Tests (150 LOC)
- **RED-GREEN**: Test command in full CLI context
  - Test command execution flow
  - Verify flag behavior and precedence
  - Test error messages
- **Files**: `cmd/push/integration_test.go`

### Wave 1.2: Input Validation & Error Handling
**Efforts: 2 | Total Lines: ~300**

#### Effort 1.2.1: Validation Tests (150 LOC)
- **RED Phase**: Write validation test suite
  - Test image path validation
  - Test registry URL validation
  - Test authentication parameter validation
  - Test error scenarios and messages
- **Files**: `pkg/oci/validation_test.go`

#### Effort 1.2.2: Implement Validation (150 LOC)
- **GREEN-REFACTOR**: Validation implementation
  - Input sanitization
  - Error handling patterns
  - Consistent error messages
- **Files**: `pkg/oci/validation.go`

## PHASE 2: AUTHENTICATION & CREDENTIALS
**Target: 1200 LOC | Duration: 3-4 days**

### Wave 2.1: Credential Management with TDD
**Efforts: 3 | Total Lines: ~700**

#### Effort 2.1.1: Auth Interface Tests (200 LOC)
- **RED Phase**: Define auth behavior through tests
  - Test credential retrieval from secrets
  - Test credential validation
  - Test error handling
- **Files**: `pkg/oci/auth_test.go`

#### Effort 2.1.2: Implement Auth Module (300 LOC)
- **GREEN Phase**: Auth implementation
  - Secret integration with existing code
  - Credential management
  - Authenticator interface implementation
- **Files**: `pkg/oci/auth.go`

#### Effort 2.1.3: Mock Auth for Testing (200 LOC)
- **REFACTOR**: Create test infrastructure
  - Mock authenticator
  - Test doubles for secrets
  - Integration test helpers
- **Files**: `pkg/oci/auth_mock.go`, `pkg/oci/testutil/`

### Wave 2.2: Authentication Flow
**Efforts: 2 | Total Lines: ~500**

#### Effort 2.2.1: Flow Tests (200 LOC)
- **RED Phase**: Test authentication precedence
  - Test flag override behavior
  - Test default credential fallback
  - Test authentication failures
- **Files**: `pkg/oci/flow_test.go`

#### Effort 2.2.2: Implement Auth Flow (300 LOC)
- **GREEN-REFACTOR**: Complete auth flow
  - Flag override logic
  - Default credential handling from secrets
  - Error propagation
- **Files**: `pkg/oci/flow.go`

## PHASE 3: OCI REGISTRY INTEGRATION
**Target: 1500 LOC | Duration: 4-5 days**

### Wave 3.1: OCI Client Implementation with TDD
**Efforts: 3 | Total Lines: ~900**

#### Effort 3.1.1: Client Interface Tests (300 LOC)
- **RED Phase**: Define client behavior
  - Test registry connection
  - Test transport configuration
  - Test authentication integration
- **Files**: `pkg/oci/client_test.go`

#### Effort 3.1.2: Implement OCI Client (400 LOC)
- **GREEN Phase**: go-containerregistry integration
  - Registry client implementation
  - Transport configuration
  - Authentication setup
- **Files**: `pkg/oci/client.go`

#### Effort 3.1.3: Insecure Mode Handling (200 LOC)
- **RED-GREEN-REFACTOR**: TLS configuration
  - Test insecure flag behavior
  - Implement certificate validation bypass
  - Transport customization
- **Files**: `pkg/oci/transport.go`, `pkg/oci/transport_test.go`

### Wave 3.2: Push Operation
**Efforts: 3 | Total Lines: ~600**

#### Effort 3.2.1: Push Operation Tests (200 LOC)
- **RED Phase**: Define push behavior
  - Test image validation
  - Test push scenarios
  - Test progress reporting
- **Files**: `pkg/oci/push_test.go`

#### Effort 3.2.2: Implement Push (400 LOC)
- **GREEN-REFACTOR**: Push implementation
  - Image loading and validation
  - Push execution with go-containerregistry
  - Progress reporting
  - Error handling
- **Files**: `pkg/oci/push.go`

## PHASE 4: INTEGRATION & POLISH
**Target: 700 LOC | Duration: 2-3 days**

### Wave 4.1: E2E Testing & Documentation
**Efforts: 2 | Total Lines: ~700**

#### Effort 4.1.1: E2E Test Suite (400 LOC)
- **Comprehensive Testing**: Full workflow validation
  - Complete push workflows
  - Error recovery scenarios
  - Performance tests
  - Integration with test registry
- **Files**: `test/e2e/push_test.go`

#### Effort 4.1.2: Documentation & Examples (300 LOC)
- **Documentation**: User and developer docs
  - Code documentation
  - README updates
  - Usage examples
  - Troubleshooting guide
- **Files**: `docs/push-command.md`, `examples/`

## Risk Mitigation

### Technical Risks
1. **Registry Compatibility**: Extensive testing against multiple registries
2. **Large Image Handling**: Stream-based operations, memory management
3. **Network Reliability**: Retry logic with exponential backoff

### Security Risks
1. **Credential Exposure**: Masked logging, secure storage
2. **Insecure Mode**: Clear warnings, explicit confirmation

## Validation Criteria

### TDD Compliance Checklist
- [ ] Every effort starts with failing tests (RED)
- [ ] Implementation is minimal to pass tests (GREEN)
- [ ] Code is refactored after passing (REFACTOR)
- [ ] Test coverage exceeds 80%
- [ ] All tests run in CI/CD pipeline

### Functional Requirements
- [ ] Push command successfully uploads images
- [ ] Authentication works with explicit credentials
- [ ] Authentication falls back to secrets
- [ ] Insecure mode handles self-signed certificates
- [ ] Error messages are clear and actionable

### Quality Metrics
- [ ] Performance: 100MB image uploads in <30 seconds
- [ ] Memory usage: <500MB for typical operations
- [ ] Zero regression in existing functionality
- [ ] All code reviewed before merge

## Implementation Notes

### TDD Enforcement Strategy
1. **No Code Without Tests**: Every PR must show tests written first
2. **Test Commit History**: Tests committed before implementation
3. **Coverage Gates**: Builds fail if coverage drops below 80%
4. **Review Checklist**: Code reviewers verify TDD compliance

### Integration Points
1. **Command Registration**: Add to existing command tree
2. **Secret Management**: Reuse existing `get secrets` infrastructure
3. **Logging**: Use existing logr configuration
4. **Error Patterns**: Match existing error handling style

## Success Metrics
- 100% of features have tests written first
- 80%+ test coverage achieved
- Zero production bugs in first 30 days
- Implementation completed within timeline
- Full OCI standards compliance

## Next Steps
1. Initialize repository structure
2. Set up CI/CD with coverage gates
3. Begin Phase 1, Wave 1.1, Effort 1.1.1: Write Command Tests
4. Follow strict TDD methodology throughout

---
*This plan enforces TEST-DRIVEN DEVELOPMENT at every stage. Tests MUST be written first, implementation second, refactoring third.*