# PROJECT ARCHITECTURE PLAN
## IDPBuilder OCI Push Command

---

## Document Metadata

| Field | Value |
|-------|-------|
| **Project Name** | idpbuilder-oci-push-command |
| **Project Type** | CLI Enhancement |
| **Primary Language** | Go |
| **Architecture Style** | Modular CLI with Interface-First Design |
| **Document Version** | v1.0 |
| **Created** | 2025-10-28 |
| **Status** | Ready for Implementation |
| **Based on PRD** | prd/idpbuilder-oci-push-command-prd.md |

---

## 1. Executive Summary

This architecture plan defines the complete implementation strategy for adding OCI image push functionality to IDPBuilder. The design follows Software Factory 3.0 best practices with interface-first development, enabling maximum parallelization through clear contracts established in Phase 1 Wave 1.

**Key Architectural Decisions:**
- **Interface-First Design**: All interfaces defined upfront in Phase 1 Wave 1
- **Modular Components**: Clear separation of concerns (Docker, Registry, Auth, TLS)
- **Dependency on go-containerregistry**: Leverage industry-standard OCI library
- **No Certificate Management**: Rely on --insecure flag only (per PRD exclusions)
- **Docker Daemon Integration**: Read images from local Docker, no build functionality

**Implementation Scope:**
- **3 Phases**: Foundation → Core Features → Polish & Integration
- **7 Total Waves**: Organized for maximum parallelization
- **13 Total Efforts**: Each scoped to ~400-700 lines for safety margin
- **No breaking changes**: Pure additive enhancement to IDPBuilder

---

## 2. Architectural Principles & Constraints

### 2.1 Design Principles

**Interface-First Development (R307 Compliance):**
- All package interfaces defined in Phase 1 Wave 1
- Enables parallel implementation in Phase 1 Wave 2
- Ensures independent branch mergeability
- Contracts frozen before implementation begins

**Modularity & Separation of Concerns:**
- Docker integration isolated in `pkg/docker`
- Registry operations isolated in `pkg/registry`
- Authentication isolated in `pkg/auth`
- TLS configuration isolated in `pkg/tls`
- Command layer in `cmd/push`

**Incremental Branching (R308 Compliance):**
- Phase 1 Wave 2 branches from Phase 1 Wave 1 integration
- Phase 2 builds on Phase 1 integration
- Phase 3 builds on Phase 2 integration
- Each wave incrementally adds functionality

**Size Compliance (R220/R221):**
- Target 400-700 lines per effort (safety margin below 800 hard limit)
- Each effort focused on single responsibility
- Integration tests separate from implementation
- Documentation efforts separate from code

### 2.2 Technical Constraints

**MUST Requirements:**
- ✅ Use go-containerregistry library (user-specified)
- ✅ Integrate with existing IDPBuilder cobra CLI
- ✅ Support --insecure flag for TLS bypass
- ✅ Default registry: https://gitea.cnoe.localtest.me:8443/
- ✅ Default username: giteaadmin
- ✅ Support special characters in passwords
- ✅ Read images from Docker daemon (no build)

**MUST NOT Requirements:**
- ❌ No image build functionality
- ❌ No certificate bundle management
- ❌ No automatic certificate export from Gitea
- ❌ No credential storage (deferred to future)

### 2.3 Dependency Analysis

**External Libraries (go.mod additions):**
```go
require (
    github.com/google/go-containerregistry v0.19.0 // OCI registry client
    github.com/docker/docker v24.0.0+incompatible  // Docker Engine API
    github.com/spf13/cobra v1.x.x                  // Already in IDPBuilder
    github.com/spf13/viper v1.x.x                  // Already in IDPBuilder
)
```

**System Dependencies:**
- Docker daemon running locally
- Network access to target registry
- HTTPS support (standard Go crypto/tls)

---

## 3. System Architecture

### 3.1 Component Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                     IDPBuilder CLI                          │
│                    (cobra framework)                        │
└────────────────────────┬────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│                  cmd/push (Command Layer)                   │
│  • Flag parsing & validation                                │
│  • User interaction & progress display                      │
│  • Orchestration of push workflow                           │
└───┬──────────────┬──────────────┬──────────────┬────────────┘
    │              │              │              │
    ▼              ▼              ▼              ▼
┌─────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐
│  Docker │  │ Registry │  │   Auth   │  │   TLS    │
│   pkg   │  │   pkg    │  │   pkg    │  │   pkg    │
└────┬────┘  └────┬─────┘  └────┬─────┘  └────┬─────┘
     │            │             │             │
     ▼            ▼             ▼             ▼
┌─────────┐  ┌──────────────────────────────────┐
│ Docker  │  │  go-containerregistry library    │
│ Daemon  │  │  (OCI image operations)          │
└─────────┘  └──────────────────────────────────┘
```

### 3.2 Package Structure

```
idpbuilder/
├── cmd/
│   └── push.go                    # Cobra command definition
│
├── pkg/
│   ├── docker/                    # Docker daemon integration
│   │   ├── interface.go           # Docker client interface
│   │   ├── client.go              # Docker Engine API implementation
│   │   ├── image_validator.go    # Image existence validation
│   │   └── client_test.go
│   │
│   ├── registry/                  # OCI registry operations
│   │   ├── interface.go           # Registry client interface
│   │   ├── client.go              # go-containerregistry wrapper
│   │   ├── pusher.go              # Push orchestration
│   │   ├── progress.go            # Progress reporting
│   │   └── client_test.go
│   │
│   ├── auth/                      # Authentication
│   │   ├── interface.go           # Auth provider interface
│   │   ├── basic.go               # Basic auth implementation
│   │   ├── validator.go           # Credential validation
│   │   └── basic_test.go
│   │
│   └── tls/                       # TLS configuration
│       ├── interface.go           # TLS config provider interface
│       ├── config.go              # TLS setup (insecure mode)
│       └── config_test.go
│
└── test/
    └── integration/
        └── push_test.go           # E2E integration tests
```

### 3.3 Data Flow Architecture

```
User executes: idpbuilder push myapp:latest --insecure

1. Command Layer (cmd/push.go)
   ├─ Parse flags: --registry, --username, --password, --insecure
   ├─ Validate arguments: image name format
   └─ Initialize components

2. Docker Integration (pkg/docker)
   ├─ Connect to Docker daemon
   ├─ Validate image exists: myapp:latest
   └─ Prepare image for export

3. Authentication (pkg/auth)
   ├─ Build basic auth credentials
   └─ Validate username/password format

4. TLS Configuration (pkg/tls)
   ├─ Configure InsecureSkipVerify if --insecure
   └─ Or use system cert pool

5. Registry Client (pkg/registry)
   ├─ Authenticate to registry
   ├─ Push image layers (with progress)
   ├─ Push manifest
   └─ Verify push success

6. Progress Reporting
   ├─ Display layer upload progress
   └─ Report final success/failure

7. Return to user with exit code
```

---

## 4. Interface Contracts (Phase 1 Wave 1)

### 4.1 Docker Client Interface

**File:** `pkg/docker/interface.go`

```go
package docker

import (
    "context"
    "io"
    v1 "github.com/google/go-containerregistry/pkg/v1"
)

// Client defines operations for interacting with Docker daemon
type Client interface {
    // ImageExists checks if an image exists in the local Docker daemon
    // Returns: exists (bool), error
    ImageExists(ctx context.Context, imageName string) (bool, error)

    // GetImage retrieves an image from Docker daemon as v1.Image
    // Returns: OCI image object, error
    GetImage(ctx context.Context, imageName string) (v1.Image, error)

    // ValidateImageName checks if image name follows OCI spec
    // Returns: error if invalid, nil if valid
    ValidateImageName(imageName string) error

    // Close cleans up Docker client resources
    Close() error
}

// NewClient creates a new Docker client
// Returns: Docker client, error
func NewClient() (Client, error) {
    // Implementation in client.go
}
```

### 4.2 Registry Client Interface

**File:** `pkg/registry/interface.go`

```go
package registry

import (
    "context"
    v1 "github.com/google/go-containerregistry/pkg/v1"
)

// Client defines operations for OCI registry interactions
type Client interface {
    // Push pushes an image to the registry
    // Parameters:
    //   - ctx: context for cancellation
    //   - image: OCI v1.Image to push
    //   - targetRef: fully qualified image reference (registry/namespace/image:tag)
    //   - progressCallback: optional callback for progress updates
    // Returns: error
    Push(ctx context.Context, image v1.Image, targetRef string, progressCallback ProgressCallback) error

    // BuildImageReference constructs registry image reference
    // Parameters:
    //   - registryURL: base registry URL
    //   - imageName: image name with optional tag
    // Returns: fully qualified reference string
    BuildImageReference(registryURL, imageName string) (string, error)

    // ValidateRegistry checks if registry is accessible
    // Returns: error if unreachable
    ValidateRegistry(ctx context.Context, registryURL string) error
}

// ProgressCallback is invoked during layer uploads
type ProgressCallback func(update ProgressUpdate)

// ProgressUpdate contains progress information
type ProgressUpdate struct {
    LayerDigest  string
    LayerSize    int64
    BytesPushed  int64
    Status       string // "uploading", "complete", "exists"
}

// NewClient creates a new registry client with authentication
func NewClient(authProvider AuthProvider, tlsConfig TLSConfig) (Client, error) {
    // Implementation in client.go
}
```

### 4.3 Authentication Interface

**File:** `pkg/auth/interface.go`

```go
package auth

import (
    "github.com/google/go-containerregistry/pkg/authn"
)

// Provider defines authentication credential provision
type Provider interface {
    // GetAuthenticator returns an authn.Authenticator for go-containerregistry
    // Returns: authenticator, error
    GetAuthenticator() (authn.Authenticator, error)

    // ValidateCredentials checks if credentials are well-formed
    // Returns: error if invalid
    ValidateCredentials() error
}

// Credentials holds authentication information
type Credentials struct {
    Username string
    Password string
}

// NewBasicAuthProvider creates a basic auth provider
func NewBasicAuthProvider(username, password string) Provider {
    // Implementation in basic.go
}
```

### 4.4 TLS Configuration Interface

**File:** `pkg/tls/interface.go`

```go
package tls

import (
    "crypto/tls"
)

// ConfigProvider defines TLS configuration provision
type ConfigProvider interface {
    // GetTLSConfig returns a tls.Config for HTTP transport
    // Returns: TLS configuration
    GetTLSConfig() *tls.Config

    // IsInsecure returns whether InsecureSkipVerify is enabled
    IsInsecure() bool
}

// Config holds TLS configuration options
type Config struct {
    InsecureSkipVerify bool
}

// NewConfigProvider creates a TLS config provider
func NewConfigProvider(insecure bool) ConfigProvider {
    // Implementation in config.go
}
```

---

## 5. Phase Breakdown

### Phase 1: Foundation & Interfaces
**Goal:** Establish all interfaces and core package structure
**Duration Estimate:** 2 waves
**Outcome:** All contracts defined, basic implementations ready

### Phase 2: Core Push Functionality
**Goal:** Implement complete push workflow with all features
**Duration Estimate:** 3 waves
**Outcome:** Fully functional push command with auth and TLS support

### Phase 3: Testing & Integration
**Goal:** Comprehensive testing and IDPBuilder integration
**Duration Estimate:** 2 waves
**Outcome:** Production-ready, tested, documented feature

---

## 6. Detailed Phase & Wave Breakdown

---

### **PHASE 1: FOUNDATION & INTERFACES**

**Objective:** Define all interfaces and establish package structure for parallel development

---

#### **PHASE 1, WAVE 1: Interface & Contract Definitions**

**Goal:** Define ALL interfaces upfront to enable Phase 1 Wave 2 parallelization

**Efforts:**

**Effort 1.1.1: Docker Client Interface Definition**
- **File:** `pkg/docker/interface.go`
- **Scope:**
  - Define `Client` interface with ImageExists, GetImage, ValidateImageName, Close
  - Define `NewClient()` constructor signature
  - Add comprehensive GoDoc comments
  - Define error types for Docker operations
- **Size Estimate:** ~150 lines
- **Acceptance Criteria:**
  - Interface compiles without implementation
  - All method signatures documented
  - Error types defined
  - No implementation code (interface only)

**Effort 1.1.2: Registry Client Interface Definition**
- **File:** `pkg/registry/interface.go`
- **Scope:**
  - Define `Client` interface with Push, BuildImageReference, ValidateRegistry
  - Define `ProgressCallback` function type
  - Define `ProgressUpdate` struct
  - Define `NewClient()` constructor signature
  - Add comprehensive GoDoc comments
- **Size Estimate:** ~180 lines
- **Acceptance Criteria:**
  - Interface compiles
  - Progress reporting types defined
  - All methods documented
  - Integration points with auth/tls defined

**Effort 1.1.3: Auth & TLS Interface Definitions**
- **Files:** `pkg/auth/interface.go`, `pkg/tls/interface.go`
- **Scope:**
  - Define `auth.Provider` interface with GetAuthenticator, ValidateCredentials
  - Define `auth.Credentials` struct
  - Define `tls.ConfigProvider` interface with GetTLSConfig, IsInsecure
  - Define `tls.Config` struct
  - Add GoDoc comments
- **Size Estimate:** ~120 lines
- **Acceptance Criteria:**
  - Both interfaces compile
  - Compatible with go-containerregistry's authn package
  - Clear separation of concerns

**Effort 1.1.4: Command Structure & Flag Definitions**
- **File:** `cmd/push.go` (skeleton)
- **Scope:**
  - Define cobra command structure for `idpbuilder push`
  - Define all flags: --registry, --username, --password, -k/--insecure, --verbose
  - Define flag validation functions (signatures only)
  - Add help text and examples
  - Define command execution function signature
- **Size Estimate:** ~200 lines
- **Acceptance Criteria:**
  - Command registers with cobra
  - All flags defined with defaults
  - Help text complete
  - Compiles (no implementation yet)

**Wave 1 Integration Branch:** `phase1-wave1-integration`
- Merge all 4 efforts
- Verify all interfaces compile together
- Run `go build` to ensure no conflicts
- Document interface contracts in architecture docs

**Wave 1 Size Check:**
- Total new lines: ~650 (well under 800 limit per effort)
- All efforts mergeable independently to main

---

#### **PHASE 1, WAVE 2: Core Package Implementations** (Parallel)

**Goal:** Implement all packages in parallel using Phase 1 Wave 1 contracts

**Branch Strategy:** All efforts branch from `phase1-wave1-integration`

**Effort 1.2.1: Docker Client Implementation**
- **File:** `pkg/docker/client.go`, `pkg/docker/image_validator.go`
- **Scope:**
  - Implement `Client` interface using Docker Engine API
  - Implement `ImageExists()` with daemon connection
  - Implement `GetImage()` to convert Docker image to v1.Image
  - Implement `ValidateImageName()` per OCI spec
  - Implement `Close()` for resource cleanup
  - Add unit tests using mocked Docker client
- **Size Estimate:** ~500 lines
- **Acceptance Criteria:**
  - All interface methods implemented
  - Unit tests with 85%+ coverage
  - Handles Docker daemon connection errors
  - Validates image names per OCI spec

**Effort 1.2.2: Registry Client Implementation**
- **File:** `pkg/registry/client.go`, `pkg/registry/pusher.go`
- **Scope:**
  - Implement `Client` interface using go-containerregistry
  - Implement `Push()` with layer upload and manifest push
  - Implement `BuildImageReference()` for registry URL construction
  - Implement `ValidateRegistry()` with registry ping
  - Integrate with auth.Provider and tls.ConfigProvider
  - Add unit tests
- **Size Estimate:** ~550 lines
- **Acceptance Criteria:**
  - Push workflow implemented
  - Integrates with go-containerregistry
  - Uses auth and TLS providers
  - Unit tests with 85%+ coverage

**Effort 1.2.3: Authentication Implementation**
- **File:** `pkg/auth/basic.go`, `pkg/auth/validator.go`
- **Scope:**
  - Implement `Provider` interface for basic auth
  - Implement `GetAuthenticator()` returning authn.Authenticator
  - Implement `ValidateCredentials()` for username/password validation
  - Support special characters in passwords
  - Add unit tests for credential validation
- **Size Estimate:** ~300 lines
- **Acceptance Criteria:**
  - Basic auth implemented
  - Handles special characters (including quotes, spaces, unicode)
  - Password length supports 256+ characters
  - Unit tests with 90%+ coverage

**Effort 1.2.4: TLS Configuration Implementation**
- **File:** `pkg/tls/config.go`
- **Scope:**
  - Implement `ConfigProvider` interface
  - Implement `GetTLSConfig()` with InsecureSkipVerify support
  - Implement system cert pool fallback when secure
  - Implement `IsInsecure()` method
  - Add unit tests
- **Size Estimate:** ~200 lines
- **Acceptance Criteria:**
  - TLS config provider implemented
  - Insecure mode bypasses verification
  - Secure mode uses system certs
  - Unit tests with 90%+ coverage

**Wave 2 Integration Branch:** `phase1-wave2-integration`
- Merge all 4 parallel efforts
- Verify all packages work together
- Run all unit tests
- Verify ~1550 total new lines across efforts

**Wave 2 Verification:**
- All interfaces implemented
- Unit tests passing
- No integration with command layer yet (Phase 2)

---

### **PHASE 2: CORE PUSH FUNCTIONALITY**

**Objective:** Connect all components and implement complete push command

**Branch Strategy:** Phase 2 efforts branch from `phase1-integration`

---

#### **PHASE 2, WAVE 1: Command Implementation & Integration**

**Goal:** Wire all packages together in the push command

**Effort 2.1.1: Push Command Core Logic**
- **File:** `cmd/push.go` (main implementation)
- **Scope:**
  - Implement command execution function (RunE)
  - Initialize Docker client
  - Initialize registry client with auth and TLS
  - Orchestrate push workflow: validate → retrieve → push
  - Implement flag validation logic
  - Add error handling and user feedback
- **Size Estimate:** ~450 lines
- **Acceptance Criteria:**
  - Command executes end-to-end
  - All flags functional
  - Error handling comprehensive
  - User-friendly error messages

**Effort 2.1.2: Progress Reporting Implementation**
- **File:** `pkg/registry/progress.go`, `cmd/push_progress.go`
- **Scope:**
  - Implement `ProgressCallback` handler
  - Implement real-time layer upload progress display
  - Implement progress bar or status updates
  - Add verbose mode for detailed output
  - Format layer sizes and percentages
- **Size Estimate:** ~300 lines
- **Acceptance Criteria:**
  - Real-time progress updates during push
  - Layer-by-layer status displayed
  - Verbose mode shows detailed logs
  - Clean output formatting

**Wave 1 Integration Branch:** `phase2-wave1-integration`
- Merge both efforts
- Test complete push workflow
- Verify progress reporting works

---

#### **PHASE 2, WAVE 2: Advanced Features** (Parallel)

**Goal:** Implement registry override and environment variable support

**Branch Strategy:** Branch from `phase2-wave1-integration`

**Effort 2.2.1: Custom Registry Override**
- **File:** `cmd/push.go` (flag processing), `pkg/registry/reference_builder.go`
- **Scope:**
  - Implement --registry flag handling
  - Override default registry in image reference
  - Validate custom registry URLs
  - Update BuildImageReference to handle overrides
  - Add unit tests for registry override logic
- **Size Estimate:** ~250 lines
- **Acceptance Criteria:**
  - --registry flag overrides default
  - Custom registry URLs validated
  - Works with both HTTP and HTTPS
  - Unit tests for edge cases

**Effort 2.2.2: Environment Variable Support**
- **File:** `cmd/push.go` (env var integration)
- **Scope:**
  - Implement environment variable reading for:
    - IDPBUILDER_REGISTRY
    - IDPBUILDER_REGISTRY_USERNAME
    - IDPBUILDER_REGISTRY_PASSWORD
    - IDPBUILDER_INSECURE
  - Implement priority: flags > env vars > defaults
  - Add validation for env var values
  - Document environment variables in help text
- **Size Estimate:** ~200 lines
- **Acceptance Criteria:**
  - All env vars functional
  - Correct priority order
  - Documented in --help
  - Unit tests for precedence

**Wave 2 Integration Branch:** `phase2-wave2-integration`
- Merge both efforts
- Test registry override with env vars
- Verify flag priority works

---

#### **PHASE 2, WAVE 3: Error Handling & Validation**

**Goal:** Comprehensive error handling and input validation

**Effort 2.3.1: Input Validation & Sanitization**
- **File:** `pkg/validator/image.go`, `pkg/validator/credentials.go`
- **Scope:**
  - Implement comprehensive image name validation (OCI spec)
  - Implement registry URL validation
  - Implement credential sanitization
  - Prevent command injection in image names
  - Add validation for special characters
  - Comprehensive unit tests for all edge cases
- **Size Estimate:** ~400 lines
- **Acceptance Criteria:**
  - All inputs validated per OCI spec
  - Command injection prevented
  - Clear validation error messages
  - 95%+ test coverage

**Effort 2.3.2: Error Handling & Exit Codes**
- **File:** `cmd/push_errors.go`, `pkg/errors/types.go`
- **Scope:**
  - Define error types for different failure modes
  - Implement proper exit codes (0, 1, 2, 3, 4 per PRD)
  - Implement actionable error messages with suggestions
  - Add error recovery suggestions
  - Add unit tests for error formatting
- **Size Estimate:** ~350 lines
- **Acceptance Criteria:**
  - Exit codes match PRD specification
  - Error messages are actionable
  - Suggestions provided for common errors
  - Error types clearly categorized

**Wave 3 Integration Branch:** `phase2-wave3-integration`
- Merge both efforts
- Test all error scenarios
- Verify exit codes correct

**Phase 2 Final Integration:** `phase2-integration`
- Merge all Wave 1, 2, 3 integrations
- Complete push functionality ready
- All features implemented

---

### **PHASE 3: TESTING & INTEGRATION**

**Objective:** Comprehensive testing and IDPBuilder integration

**Branch Strategy:** Phase 3 efforts branch from `phase2-integration`

---

#### **PHASE 3, WAVE 1: Integration Testing**

**Goal:** E2E and integration tests

**Effort 3.1.1: Integration Tests - Core Workflow**
- **File:** `test/integration/push_test.go`
- **Scope:**
  - E2E test: local Docker image → push → verify in registry
  - Test with Gitea registry (Docker container)
  - Test authentication success/failure
  - Test insecure mode TLS bypass
  - Test large image push (100MB+)
  - Setup/teardown test infrastructure
- **Size Estimate:** ~500 lines
- **Acceptance Criteria:**
  - All E2E workflows tested
  - Test Gitea registry in Docker
  - Tests run in CI
  - Cleanup after tests

**Effort 3.1.2: Integration Tests - Edge Cases**
- **File:** `test/integration/push_edge_cases_test.go`
- **Scope:**
  - Test network failure scenarios
  - Test missing image errors
  - Test invalid credentials
  - Test registry unreachable
  - Test multi-architecture images
  - Test tag override scenarios
- **Size Estimate:** ~400 lines
- **Acceptance Criteria:**
  - All error paths tested
  - Edge cases covered
  - Flaky test detection/mitigation
  - Runs in CI pipeline

**Wave 1 Integration Branch:** `phase3-wave1-integration`
- All integration tests passing
- Test coverage report generated

---

#### **PHASE 3, WAVE 2: Documentation & IDPBuilder Integration**

**Goal:** Documentation and final IDPBuilder integration

**Effort 3.2.1: User Documentation**
- **Files:** `docs/push-command.md`, updates to main README
- **Scope:**
  - Comprehensive command documentation
  - Usage examples for all scenarios
  - Troubleshooting guide
  - FAQ for common issues
  - Update IDPBuilder main README with push command
- **Size Estimate:** ~300 lines (documentation)
- **Acceptance Criteria:**
  - All features documented
  - Examples tested and verified
  - Troubleshooting covers common issues
  - README updated

**Effort 3.2.2: IDPBuilder Integration & Build System**
- **Files:** `Makefile` updates, `go.mod` updates, CI integration
- **Scope:**
  - Update IDPBuilder Makefile to include push command
  - Update go.mod with new dependencies
  - Update CI/CD pipelines for new tests
  - Verify build system includes all new files
  - Add push command to IDPBuilder help
- **Size Estimate:** ~200 lines (build configs)
- **Acceptance Criteria:**
  - `make build` includes push command
  - `make test` runs all new tests
  - CI pipeline updated
  - Dependencies properly versioned

**Wave 2 Integration Branch:** `phase3-wave2-integration`
- Documentation complete
- Build system verified

**Phase 3 Final Integration:** `phase3-integration`
- All phases complete
- Ready for PR to main IDPBuilder repo

---

## 7. Effort Summary Table

| Phase | Wave | Effort | Description | Est. Lines | Parallelizable |
|-------|------|--------|-------------|-----------|----------------|
| 1 | 1 | 1.1.1 | Docker Client Interface | ~150 | No (sequential) |
| 1 | 1 | 1.1.2 | Registry Client Interface | ~180 | No (sequential) |
| 1 | 1 | 1.1.3 | Auth & TLS Interfaces | ~120 | No (sequential) |
| 1 | 1 | 1.1.4 | Command Structure & Flags | ~200 | No (sequential) |
| 1 | 2 | 1.2.1 | Docker Client Implementation | ~500 | **YES** |
| 1 | 2 | 1.2.2 | Registry Client Implementation | ~550 | **YES** |
| 1 | 2 | 1.2.3 | Authentication Implementation | ~300 | **YES** |
| 1 | 2 | 1.2.4 | TLS Configuration Implementation | ~200 | **YES** |
| 2 | 1 | 2.1.1 | Push Command Core Logic | ~450 | No |
| 2 | 1 | 2.1.2 | Progress Reporting | ~300 | No |
| 2 | 2 | 2.2.1 | Custom Registry Override | ~250 | **YES** |
| 2 | 2 | 2.2.2 | Environment Variable Support | ~200 | **YES** |
| 2 | 3 | 2.3.1 | Input Validation & Sanitization | ~400 | No |
| 2 | 3 | 2.3.2 | Error Handling & Exit Codes | ~350 | No |
| 3 | 1 | 3.1.1 | Integration Tests - Core | ~500 | No |
| 3 | 1 | 3.1.2 | Integration Tests - Edge Cases | ~400 | **YES** |
| 3 | 2 | 3.2.1 | User Documentation | ~300 | **YES** |
| 3 | 2 | 3.2.2 | IDPBuilder Integration | ~200 | **YES** |

**Total Estimated Lines:** ~5,650 lines (across all efforts)

**Parallelization Opportunities:**
- Phase 1 Wave 2: 4 parallel efforts (Docker, Registry, Auth, TLS)
- Phase 2 Wave 2: 2 parallel efforts (Registry override, Env vars)
- Phase 3 Wave 1: 2 parallel efforts (Core tests, Edge cases)
- Phase 3 Wave 2: 2 parallel efforts (Docs, Build integration)

**Safety Margin:** All efforts under 600 lines (well below 800 hard limit)

---

## 8. Interface-First Development Benefits

**Why Phase 1 Wave 1 Defines All Interfaces:**

1. **Maximum Parallelization (R307):**
   - Phase 1 Wave 2: 4 teams can implement Docker, Registry, Auth, TLS simultaneously
   - All teams code against frozen interfaces
   - No coordination needed during implementation

2. **Independent Mergeability:**
   - Each implementation can merge to main independently
   - No breaking changes across the wave
   - Build always green (interfaces already compiled)

3. **Incremental Branching (R308):**
   - Phase 2 builds on complete Phase 1 foundation
   - Each wave adds functionality incrementally
   - No "big bang" integration at end

4. **Clear Contracts:**
   - All method signatures documented upfront
   - Error types defined before implementation
   - No interface changes mid-development

---

## 9. Integration Strategy

### 9.1 Wave Integration Process

**After each wave:**
1. Create integration branch (e.g., `phase1-wave1-integration`)
2. Merge all wave efforts sequentially
3. Run full test suite
4. Verify build succeeds
5. Tag integration point
6. Next wave branches from integration branch

**Integration Branch Naming:**
- Phase 1 Wave 1: `phase1-wave1-integration`
- Phase 1 Wave 2: `phase1-wave2-integration`
- Phase 1 Complete: `phase1-integration`
- (Continue for all phases)

### 9.2 Phase Integration Process

**After each phase:**
1. Create phase integration branch
2. Merge all wave integrations
3. Run comprehensive tests
4. Performance benchmarks
5. Architect review
6. Tag phase completion

### 9.3 Final Integration to IDPBuilder

**After Phase 3 complete:**
1. Create PR from `phase3-integration` to IDPBuilder main
2. Include all tests, documentation, build system updates
3. IDPBuilder team review
4. Merge and release with next IDPBuilder version

---

## 10. Testing Strategy by Phase

### Phase 1 Testing
- **Unit tests:** Each package (docker, registry, auth, tls)
- **Coverage target:** 85% minimum
- **Mock dependencies:** Docker daemon, registry API
- **Focus:** Interface compliance, error handling

### Phase 2 Testing
- **Unit tests:** Command logic, validation, error handling
- **Coverage target:** 85% minimum
- **Integration tests:** Components working together
- **Focus:** Workflow orchestration, user interaction

### Phase 3 Testing
- **E2E tests:** Complete workflows with real Gitea
- **Integration tests:** All components together
- **Performance tests:** Large image pushes
- **Coverage target:** 80% overall, 95% for critical paths

---

## 11. Dependency Management

### 11.1 Go Module Updates

**New dependencies to add:**
```go
require (
    github.com/google/go-containerregistry v0.19.0
    github.com/docker/docker v24.0.0+incompatible
    github.com/docker/docker/client v24.0.0+incompatible
)
```

**Version selection rationale:**
- go-containerregistry: Latest stable (v0.19.0 as of architecture planning)
- docker/docker: Latest stable API client
- Use `go get` to fetch exact latest versions during implementation

### 11.2 Dependency Review

**License compliance:**
- go-containerregistry: Apache 2.0 ✅
- docker/docker: Apache 2.0 ✅
- No GPL or restrictive licenses

**Security:**
- All dependencies will be scanned via Dependabot
- Regular updates as part of IDPBuilder maintenance

---

## 12. Performance Considerations

### 12.1 Performance Targets (from PRD)

| Metric | Target | Implementation Strategy |
|--------|--------|-------------------------|
| Command startup | <500ms | Lazy initialization of Docker client |
| Memory footprint | <200MB | Stream layers, avoid buffering entire images |
| Progress reporting | Real-time | Use go-containerregistry's progress callbacks |
| Large image (100MB) | <30s | Network-dependent, chunked uploads |

### 12.2 Optimization Strategies

**Memory:**
- Stream image layers (no full buffering)
- Close Docker client connections promptly
- Reuse HTTP transport for registry operations

**Network:**
- Use go-containerregistry's chunked upload support
- Parallel layer uploads where supported by registry
- Connection pooling for multi-layer pushes

---

## 13. Security Architecture

### 13.1 Security Requirements Implementation

**Input Validation (Phase 2 Wave 3):**
- Image name validation: Prevent command injection
- Registry URL validation: Prevent SSRF attacks
- Credential validation: Ensure well-formed inputs

**Credential Handling:**
- Never log passwords or auth tokens
- Use go-containerregistry's authn package (secure)
- Future: OS keychain integration (out of scope v1)

**TLS Security:**
- Default: Strict certificate verification
- Insecure mode: Clear warning to user
- No automatic trust of self-signed certs

**Secrets in Logs:**
- Redact credentials from all log output
- Sanitize error messages
- Use structured logging with secret filtering

### 13.2 Threat Model

**Threats Mitigated:**
- ✅ Command injection via image names
- ✅ MITM attacks (with TLS verification)
- ✅ Credential exposure in logs
- ✅ SSRF via registry URL manipulation

**Threats NOT Mitigated (Out of Scope):**
- ❌ Image content scanning (future)
- ❌ Image signing verification (future)
- ❌ Registry compromise detection

---

## 14. Error Handling Architecture

### 14.1 Error Categories & Exit Codes

| Category | Exit Code | Example |
|----------|-----------|---------|
| Success | 0 | Image pushed successfully |
| General error | 1 | Invalid flag combination |
| Auth failure | 2 | Incorrect username/password |
| Network/Registry error | 3 | Registry unreachable, TLS error |
| Image not found | 4 | Image not in Docker daemon |

**Implementation:** Phase 2 Wave 3 (Error Handling & Exit Codes effort)

### 14.2 Error Message Design

**Format:**
```
Error: <what went wrong>
Suggestion: <actionable next step>
Context: <relevant details>
```

**Example:**
```
Error: Image 'myapp:latest' not found in Docker daemon
Suggestion: Run 'docker images' to list available images or build the image first with 'docker build'
Context: Searched for image: myapp:latest
```

---

## 15. Future Enhancements (Out of Scope v1)

**Explicitly Excluded from This Plan:**
- ❌ Credential storage (OS keychain)
- ❌ Custom certificate bundles
- ❌ Automatic Gitea certificate export
- ❌ Multi-registry push
- ❌ Image signing/verification
- ❌ Layer deduplication
- ❌ Resume partial uploads
- ❌ Image build functionality

**Rationale:** Keep v1 focused on core push workflow. Gather user feedback before adding complexity.

---

## 16. Compliance with Software Factory 3.0

### 16.1 Rule Compliance Matrix

| Rule | Description | Compliance Strategy |
|------|-------------|---------------------|
| R307 | Independent Branch Mergeability | Interface-first design enables parallel development |
| R308 | Incremental Branching Strategy | Each wave builds on previous integration |
| R359 | No Code Deletion | Pure additive enhancement (no deletions) |
| R383 | Metadata File Organization | All metadata in `.software-factory/` with timestamps |
| R220/R221 | Size Limits | All efforts 400-600 lines (safety margin) |

### 16.2 Verification Checkpoints

**Wave Completion Checklist:**
- [ ] All efforts merged to wave integration branch
- [ ] Build succeeds (`make build`)
- [ ] All tests passing
- [ ] No effort exceeds 700 lines (measurement via line-counter.sh)
- [ ] Metadata organized per R383

**Phase Completion Checklist:**
- [ ] All waves integrated
- [ ] Architect review passed
- [ ] Phase integration branch tagged
- [ ] Documentation updated

---

## 17. Risk Assessment & Mitigation

### 17.1 Technical Risks

| Risk | Probability | Impact | Mitigation |
|------|-------------|--------|------------|
| Docker daemon API changes | Low | Medium | Use stable Docker client library version, comprehensive tests |
| go-containerregistry breaking changes | Medium | High | Pin to stable version, monitor releases, test upgrades |
| Gitea registry incompatibilities | Low | Medium | Use standard OCI spec, test against real Gitea |
| TLS certificate issues | High | Low | Clear --insecure documentation, warn users |
| Size limit violations | Low | Low | Conservative effort scoping (400-600 lines), regular measurement |

### 17.2 Process Risks

| Risk | Probability | Impact | Mitigation |
|------|-------------|--------|------------|
| Interface changes after Wave 1 | Low | High | Architect review of interfaces before Wave 2 starts |
| Parallel efforts conflict | Low | Medium | Clear interface contracts, integration testing |
| Effort underestimation | Medium | Low | 200-line safety margin per effort (target 600, limit 800) |

---

## 18. Success Criteria

### 18.1 Phase Completion Criteria

**Phase 1 Complete:**
- ✅ All interfaces defined and documented
- ✅ All packages implemented with unit tests
- ✅ Unit test coverage ≥85%
- ✅ All code compiles and integrates

**Phase 2 Complete:**
- ✅ Push command fully functional
- ✅ All flags working (--registry, --username, --password, --insecure)
- ✅ Environment variable support
- ✅ Error handling comprehensive
- ✅ Exit codes correct

**Phase 3 Complete:**
- ✅ E2E tests passing
- ✅ Integration tests with real Gitea
- ✅ Documentation complete
- ✅ IDPBuilder build system updated
- ✅ Ready for PR to main repo

### 18.2 Overall Project Success

**Technical Success:**
- Push command works with local Docker and Gitea registry
- Supports all PRD requirements
- Code quality: 85%+ test coverage
- No efforts exceed 700 lines

**Process Success:**
- Follows Software Factory 3.0 workflow
- All rule compliance verified (R307, R308, R383, etc.)
- Maximum parallelization achieved
- Clean git history with independent merges

---

## 19. Architecture Decision Records (ADRs)

### ADR-001: Use go-containerregistry for OCI Operations
**Decision:** Use Google's go-containerregistry library for all OCI image operations
**Rationale:** User-recommended, industry-standard, well-maintained, handles OCI spec complexity
**Alternatives Considered:** Direct OCI API implementation (too complex), docker/distribution library (less mature)
**Status:** Approved

### ADR-002: Interface-First Development in Phase 1 Wave 1
**Decision:** Define all interfaces before any implementation
**Rationale:** Enables maximum parallelization in Wave 2, ensures independent mergeability (R307)
**Alternatives Considered:** Iterative interface evolution (violates R307)
**Status:** Approved

### ADR-003: No Certificate Management in v1
**Decision:** Only support --insecure flag, no certificate bundle management
**Rationale:** Explicitly excluded per PRD, reduces scope, user can manage certs externally
**Alternatives Considered:** Certificate bundle support (deferred to v2)
**Status:** Approved

### ADR-004: Basic Auth Only for v1
**Decision:** Only username/password authentication, no token auth or keychain
**Rationale:** Meets Gitea registry requirements, keeps scope focused
**Alternatives Considered:** OS keychain integration (deferred to v2)
**Status:** Approved

### ADR-005: 400-600 Line Effort Target
**Decision:** Target 400-600 lines per effort (200-line safety margin below 800 limit)
**Rationale:** Provides buffer for unexpected complexity, reduces split risk
**Alternatives Considered:** Target 700 lines (less safety margin)
**Status:** Approved

---

## 20. Appendix: Command Examples

### Basic Usage
```bash
# Push to default Gitea registry with default username
idpbuilder push myapp:latest --password 'mypassword'

# Push with custom username
idpbuilder push myapp:latest --username developer --password 'myP@ss'

# Push with insecure mode (bypass TLS verification)
idpbuilder push -k myapp:latest --password 'mypassword'
```

### Advanced Usage
```bash
# Push to custom registry
idpbuilder push --registry https://custom-registry.example.com myapp:v1.0.0 --password 'pass'

# Using environment variables
export IDPBUILDER_REGISTRY_USERNAME=developer
export IDPBUILDER_REGISTRY_PASSWORD='complex!P@ssw0rd#123'
export IDPBUILDER_INSECURE=true
idpbuilder push myapp:latest

# Verbose mode
idpbuilder push --verbose -k myapp:latest --password 'pass'
```

### Expected Output
```
Pushing myapp:latest to gitea.cnoe.localtest.me:8443/giteaadmin/myapp:latest
✓ Layer 1/3: sha256:abc123... (12.5 MB) - Complete
✓ Layer 2/3: sha256:def456... (45.2 MB) - Complete
✓ Layer 3/3: sha256:ghi789... (5.1 MB) - Complete
✓ Manifest pushed successfully

Successfully pushed myapp:latest to gitea.cnoe.localtest.me:8443/giteaadmin/myapp:latest
```

---

## 21. Glossary

- **OCI:** Open Container Initiative - standardized container formats
- **go-containerregistry:** Google's Go library for OCI registry operations
- **Gitea:** Git service with built-in container registry
- **TLS:** Transport Layer Security
- **Basic Auth:** Username/password authentication
- **Insecure Mode:** Bypass TLS certificate verification (--insecure flag)
- **Wave:** Group of related efforts within a phase
- **Effort:** Single focused implementation unit (<800 lines)
- **Integration Branch:** Branch merging all efforts in a wave
- **Interface-First:** Define all interfaces before implementation

---

## 22. Document Status

**Status:** APPROVED FOR IMPLEMENTATION
**Author:** Architect Agent (@agent-architect)
**Created:** 2025-10-28
**Version:** v1.0
**Next Step:** Hand off to Orchestrator for phase/wave execution

**Compliance Verified:**
- ✅ R307: Independent Branch Mergeability (interface-first design)
- ✅ R308: Incremental Branching Strategy (wave integration chain)
- ✅ R359: No Code Deletion (pure additive enhancement)
- ✅ R383: Metadata Organization (all metadata in `.software-factory/`)
- ✅ R220/R221: Size Limits (all efforts 400-600 line target)

**Ready for:**
- Phase 1 Wave 1 implementation
- Orchestrator to create working copies and spawn SW Engineers

---

**END OF ARCHITECTURE PLAN**
