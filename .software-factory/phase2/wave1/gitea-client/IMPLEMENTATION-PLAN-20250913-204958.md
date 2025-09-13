# Gitea Registry Client Implementation Plan

**Created**: 2025-09-13T20:49:58Z
**Location**: .software-factory/phase2/wave1/gitea-client/
**Phase**: 2 - Build & Push Implementation
**Wave**: 1 - Core Build & Push
**Effort**: E2.1.2 - gitea-registry-client
**Planner**: Code Reviewer Agent
**State**: EFFORT_PLAN_CREATION

## 🚨 CRITICAL EFFORT METADATA (FROM WAVE PLAN)

**Effort ID**: E2.1.2
**Branch**: `idpbuilder-oci-build-push/phase2/wave1/gitea-client`
**Base Branch**: `idpbuilder-oci-build-push/phase1/integration`
**Can Parallelize**: Yes
**Parallel With**: [E2.1.1 - image-builder]
**Size Estimate**: 600 lines (well within 800 limit)
**Dependencies**: Phase 1 Certificate Infrastructure (certs, certvalidation, fallback)
**Directory**: `efforts/phase2/wave1/gitea-client/pkg/`
**Existing Splits**: 001 and 002 (already created)

## 🔗 PHASE 1 DEPENDENCIES (R219 COMPLIANCE)

### Dependency Analysis
Based on Phase 1 completion, this effort will leverage:

1. **registry-tls-trust** (Phase 1, Wave 1)
   - TLS trust store management
   - Certificate validation utilities
   - Custom CA handling

2. **registry-auth-types** (Phase 1, Wave 1)
   - Authentication type definitions
   - Registry credential structures
   - Token management interfaces

3. **cert-validation** (Phase 1, Wave 2)
   - Certificate chain validation
   - Trust verification logic
   - Security policy enforcement

4. **fallback-strategies** (Phase 1, Wave 2)
   - Insecure mode handling
   - Fallback authentication mechanisms
   - Error recovery patterns

### Integration Strategy
- Import Phase 1 packages directly
- Use established trust store patterns
- Leverage authentication types
- Apply fallback strategies for --insecure mode

## 🔴🔴🔴 EXPLICIT SCOPE CONTROL (R311 MANDATORY) 🔴🔴🔴

### IMPLEMENT EXACTLY

**Core Functions (5 total):**
1. `NewGiteaRegistry(config RegistryConfig) (Registry, error)` (~40 lines)
2. `Push(ctx context.Context, image v1.Image, reference string) error` (~150 lines)
3. `Authenticate(ctx context.Context) error` (~60 lines)
4. `ListRepositories(ctx context.Context) ([]string, error)` (~50 lines)
5. `GetRemoteOptions() []remote.Option` (~40 lines)

**Core Types (4 total):**
1. `Registry interface` (~15 lines)
2. `giteaRegistryImpl struct` (~20 lines)
3. `RegistryConfig struct` (~15 lines)
4. `authenticator struct` (~10 lines)

**Test Functions (8 total):**
1. `TestNewGiteaRegistry` (~40 lines)
2. `TestPushWithValidCerts` (~50 lines)
3. `TestPushWithInsecureMode` (~50 lines)
4. `TestAuthentication` (~40 lines)
5. `TestListRepositories` (~30 lines)
6. `TestRetryLogic` (~40 lines)
7. `TestProgressReporting` (~30 lines)
8. `TestPhase1Integration` (~50 lines)

**TOTAL ESTIMATED**: ~630 lines (170 lines buffer to 800 limit)

### DO NOT IMPLEMENT

- ❌ Pull operations (future effort)
- ❌ Delete/Remove operations (future effort)
- ❌ Tag management operations (future effort)
- ❌ Manifest inspection (future effort)
- ❌ Registry catalog operations (beyond basic list)
- ❌ Token caching mechanism (future optimization)
- ❌ Connection pooling (future optimization)
- ❌ Comprehensive logging framework
- ❌ Metrics collection
- ❌ Circuit breaker pattern (keep retry simple)
- ❌ Multiple registry support (Gitea only for MVP)
- ❌ OAuth/OIDC authentication (basic auth only)

## 🔄 ATOMIC PR DESIGN (R307 COMPLIANCE)

### Independent Branch Mergeability
This effort MUST be independently mergeable per R307:

**Feature Flags:**
```yaml
- flag: "GITEA_REGISTRY_ENABLED"
  purpose: "Enable Gitea registry operations"
  default: false
  location: "pkg/config/features.go"
  activation: "When image-builder effort is also merged"
```

**Graceful Degradation:**
- Registry client returns "feature disabled" when flag is off
- Push operations log warning and return nil when disabled
- Tests verify both enabled and disabled states

**Stubs for Missing Dependencies:**
```go
// pkg/registry/stubs.go - temporary until image-builder merged
type MockImageLoader struct{}

func (m *MockImageLoader) LoadImage(path string) (v1.Image, error) {
    // Returns test image for push operations
    return empty.Image, nil
}
```

## 📁 File Structure

```
efforts/phase2/wave1/gitea-client/
├── pkg/
│   ├── registry/
│   │   ├── interface.go        # Registry interface definition (~15 lines)
│   │   ├── gitea.go           # Main GiteaRegistry implementation (~200 lines)
│   │   ├── auth.go            # Authentication handling (~60 lines)
│   │   ├── push.go            # Push operation with cert integration (~150 lines)
│   │   ├── list.go            # Repository listing operations (~50 lines)
│   │   ├── retry.go           # Simple retry logic (~60 lines)
│   │   ├── remote_options.go  # Remote options configuration (~40 lines)
│   │   └── stubs.go           # Temporary stubs for missing dependencies (~30 lines)
│   ├── config/
│   │   └── features.go        # Feature flags (~20 lines)
│   └── tests/
│       ├── gitea_test.go      # Unit tests for main implementation (~100 lines)
│       ├── push_test.go       # Push operation tests (~100 lines)
│       ├── auth_test.go       # Authentication tests (~40 lines)
│       ├── integration_test.go # Integration tests with Phase 1 (~50 lines)
│       └── test_helpers.go    # Test utilities and mocks (~40 lines)
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

## 🔨 Implementation Sequence

### Step 1: Core Interface Definition (50 lines)
```go
// pkg/registry/interface.go
package registry

import (
    "context"
    v1 "github.com/google/go-containerregistry/pkg/v1"
    "github.com/google/go-containerregistry/pkg/v1/remote"
)

// Registry defines operations for OCI registry interaction
type Registry interface {
    // Push uploads an image to the registry
    Push(ctx context.Context, image v1.Image, reference string) error

    // Authenticate performs registry authentication
    Authenticate(ctx context.Context) error

    // ListRepositories returns available repositories
    ListRepositories(ctx context.Context) ([]string, error)

    // GetRemoteOptions returns configured remote options
    GetRemoteOptions() []remote.Option
}

// RegistryConfig holds registry configuration
type RegistryConfig struct {
    URL      string
    Username string
    Password string
    Insecure bool
    CABundle string // Path to custom CA bundle from Phase 1
}
```

### Step 2: Phase 1 Integration Setup (100 lines)
```go
// pkg/registry/gitea.go
package registry

import (
    "context"
    "fmt"

    // Phase 1 dependencies
    "github.com/idpbuilder/idpbuilder-oci-build-push/pkg/certs"
    "github.com/idpbuilder/idpbuilder-oci-build-push/pkg/certvalidation"
    "github.com/idpbuilder/idpbuilder-oci-build-push/pkg/fallback"
    "github.com/idpbuilder/idpbuilder-oci-build-push/pkg/authtypes"
)

type giteaRegistryImpl struct {
    config      RegistryConfig
    trustStore  certs.TrustStoreManager
    validator   certvalidation.CertValidator
    fallback    fallback.FallbackHandler
    authType    authtypes.AuthType
    authn       *authenticator
}

func NewGiteaRegistry(config RegistryConfig) (Registry, error) {
    // Initialize Phase 1 components
    trustStore, err := certs.NewTrustStoreManager(config.CABundle)
    if err != nil {
        return nil, fmt.Errorf("failed to initialize trust store: %w", err)
    }

    validator := certvalidation.NewCertValidator()
    fallbackHandler := fallback.NewFallbackHandler(config.Insecure)

    // Determine auth type
    authType := authtypes.DetermineAuthType(config.Username, config.Password)

    return &giteaRegistryImpl{
        config:     config,
        trustStore: trustStore,
        validator:  validator,
        fallback:   fallbackHandler,
        authType:   authType,
    }, nil
}
```

### Step 3: Authentication Implementation (60 lines)
```go
// pkg/registry/auth.go
package registry

import (
    "context"
    "encoding/base64"
    "fmt"
)

type authenticator struct {
    username string
    password string
    token    string
    authType string
}

func (r *giteaRegistryImpl) Authenticate(ctx context.Context) error {
    if r.config.Username == "" && r.config.Password == "" {
        // Anonymous access
        return nil
    }

    // Create authenticator
    r.authn = &authenticator{
        username: r.config.Username,
        password: r.config.Password,
        authType: string(r.authType),
    }

    // Generate basic auth token
    auth := r.config.Username + ":" + r.config.Password
    r.authn.token = base64.StdEncoding.EncodeToString([]byte(auth))

    return nil
}
```

### Step 4: Push Operation with Certificates (150 lines)
```go
// pkg/registry/push.go
package registry

import (
    "context"
    "fmt"

    "github.com/google/go-containerregistry/pkg/name"
    "github.com/google/go-containerregistry/pkg/v1"
    "github.com/google/go-containerregistry/pkg/v1/remote"
)

func (r *giteaRegistryImpl) Push(ctx context.Context, image v1.Image, reference string) error {
    // Check feature flag
    if !isGiteaRegistryEnabled() {
        return nil // Graceful degradation
    }

    // Parse reference
    ref, err := name.ParseReference(reference)
    if err != nil {
        return fmt.Errorf("invalid reference %s: %w", reference, err)
    }

    // Get remote options with certificate configuration
    opts := r.GetRemoteOptions()

    // Add progress reporting
    opts = append(opts, remote.WithProgress(makeProgressReporter()))

    // Perform push with retry logic
    err = retryWithBackoff(ctx, func() error {
        return remote.Write(ref, image, opts...)
    })

    if err != nil {
        // Use fallback handler if TLS error
        if isTLSError(err) && r.fallback.ShouldFallback() {
            opts = r.fallback.GetInsecureOptions()
            err = remote.Write(ref, image, opts...)
        }
    }

    return err
}
```

### Step 5: Remote Options Configuration (40 lines)
```go
// pkg/registry/remote_options.go
package registry

import (
    "github.com/google/go-containerregistry/pkg/authn"
    "github.com/google/go-containerregistry/pkg/v1/remote"
)

func (r *giteaRegistryImpl) GetRemoteOptions() []remote.Option {
    var opts []remote.Option

    // Add authentication if configured
    if r.authn != nil {
        opts = append(opts, remote.WithAuth(authn.FromConfig(authn.AuthConfig{
            Username: r.authn.username,
            Password: r.authn.password,
        })))
    }

    // Configure TLS using Phase 1 trust store
    if !r.config.Insecure {
        transport := r.trustStore.GetHTTPTransport()
        opts = append(opts, remote.WithTransport(transport))
    } else {
        // Use fallback for insecure mode
        opts = append(opts, r.fallback.GetInsecureOptions()...)
    }

    return opts
}
```

### Step 6: List Operations (50 lines)
```go
// pkg/registry/list.go
package registry

import (
    "context"
    "fmt"

    "github.com/google/go-containerregistry/pkg/name"
    "github.com/google/go-containerregistry/pkg/v1/remote"
)

func (r *giteaRegistryImpl) ListRepositories(ctx context.Context) ([]string, error) {
    // Parse registry URL
    registry, err := name.NewRegistry(r.config.URL)
    if err != nil {
        return nil, fmt.Errorf("invalid registry URL: %w", err)
    }

    // Get catalog using remote options
    opts := r.GetRemoteOptions()

    repos, err := remote.Catalog(ctx, registry, opts...)
    if err != nil {
        return nil, fmt.Errorf("failed to list repositories: %w", err)
    }

    return repos, nil
}
```

### Step 7: Retry Logic (60 lines)
```go
// pkg/registry/retry.go
package registry

import (
    "context"
    "time"
)

func retryWithBackoff(ctx context.Context, fn func() error) error {
    maxRetries := 3
    backoff := time.Second

    for i := 0; i < maxRetries; i++ {
        err := fn()
        if err == nil {
            return nil
        }

        if !isRetryable(err) {
            return err
        }

        if i < maxRetries-1 {
            select {
            case <-time.After(backoff):
                backoff *= 2
            case <-ctx.Done():
                return ctx.Err()
            }
        }
    }

    return fn() // Last attempt
}

func isRetryable(err error) bool {
    // Check for transient errors
    // Network timeouts, 503 errors, etc.
    return true // Simplified for now
}
```

### Step 8: Feature Flags (20 lines)
```go
// pkg/config/features.go
package config

import "os"

func IsGiteaRegistryEnabled() bool {
    return os.Getenv("GITEA_REGISTRY_ENABLED") == "true"
}
```

### Step 9: Test Stubs (30 lines)
```go
// pkg/registry/stubs.go
package registry

import (
    v1 "github.com/google/go-containerregistry/pkg/v1"
    "github.com/google/go-containerregistry/pkg/v1/empty"
)

// MockImageLoader provides test images until image-builder is merged
type MockImageLoader struct{}

func (m *MockImageLoader) LoadImage(path string) (v1.Image, error) {
    return empty.Image, nil
}

func isGiteaRegistryEnabled() bool {
    // Stub implementation for testing
    return true
}
```

## 📏 SIZE MANAGEMENT STRATEGY

### Measurement Protocol (R304/R338 Compliance)
```bash
# Regular measurement during implementation
cd /home/vscode/workspaces/idpbuilder-oci-build-push/efforts/phase2/wave1/gitea-client

# Find project root and use line counter
PROJECT_ROOT=$(pwd)
while [ "$PROJECT_ROOT" != "/" ]; do
    [ -f "$PROJECT_ROOT/orchestrator-state.json" ] && break
    PROJECT_ROOT=$(dirname "$PROJECT_ROOT")
done

# Measure implementation lines only (excludes tests/docs)
$PROJECT_ROOT/tools/line-counter.sh

# Expected output format for R338:
# 🎯 Detected base: idpbuilder-oci-build-push/phase1/integration
# 📦 Analyzing branch: idpbuilder-oci-build-push/phase2/wave1/gitea-client
# ✅ Total implementation lines: [count]
```

### Size Checkpoints
- After Step 2 (Core Setup): ~150 lines
- After Step 4 (Push Operation): ~350 lines
- After Step 6 (List Operations): ~450 lines
- After Step 9 (Complete): ~600 lines

### Split Strategy (Already Executed)
The effort has been pre-split into:
- **Split 001**: Core interfaces, authentication, main registry (635 lines)
- **Split 002**: Push/list operations, retry logic, test stubs (633 lines)

## 🧪 TESTING STRATEGY

### Unit Test Coverage (Target: 80%)
```go
// pkg/tests/gitea_test.go
func TestNewGiteaRegistry(t *testing.T) {
    // Test registry creation with various configs
    // Verify Phase 1 component initialization
}

func TestAuthentication(t *testing.T) {
    // Test basic auth
    // Test anonymous access
    // Test token generation
}
```

### Integration Tests with Phase 1
```go
// pkg/tests/integration_test.go
func TestPhase1CertificateIntegration(t *testing.T) {
    // Verify trust store usage
    // Test custom CA handling
    // Validate cert chain verification
}

func TestFallbackStrategies(t *testing.T) {
    // Test insecure mode
    // Verify fallback activation
    // Check error recovery
}
```

### E2E Tests (When Both Efforts Merged)
```go
func TestEndToEndPush(t *testing.T) {
    // Requires image-builder effort
    // Full workflow: build -> push -> verify
}
```

## 🎯 VALIDATION CHECKPOINTS

### Checkpoint 1: After Interface Definition
- ✅ Interfaces compile
- ✅ No import cycles
- ✅ Types are exported correctly

### Checkpoint 2: After Phase 1 Integration
- ✅ Phase 1 packages import successfully
- ✅ Trust store initializes
- ✅ Certificate validation works

### Checkpoint 3: After Core Implementation
- ✅ All methods implemented (no stubs per R320)
- ✅ Size under 400 lines
- ✅ Unit tests passing

### Checkpoint 4: After Complete Implementation
- ✅ Total size under 700 lines
- ✅ All tests passing
- ✅ Feature flag works both ways
- ✅ Can merge independently (R307)

## 🔄 INTEGRATION POINTS

### With image-builder (E2.1.1)
- image-builder produces v1.Image
- gitea-client consumes v1.Image for push
- Coordinate via feature flag

### With Phase 1 Components
- Import certificate management
- Use authentication types
- Apply fallback strategies

### With Phase 2 Integration Branch
- After both efforts complete
- Merge to phase2/wave1/integration
- Remove stubs, enable features

## 📋 DELIVERABLES

1. **Working Gitea Registry Client**
   - Push operations functional
   - Authentication working
   - List repositories implemented

2. **Comprehensive Tests**
   - 80% unit test coverage
   - Integration tests with Phase 1
   - E2E test stubs for future

3. **Documentation**
   - API documentation in code
   - README with usage examples
   - Integration guide

4. **Size Compliance**
   - Under 700 lines implementation
   - Measured with line-counter.sh
   - Split plans if needed (already done)

## 🚨 RISK MITIGATION

### Risk: Gitea API Changes
- **Mitigation**: Version lock Gitea API
- **Fallback**: Use standard OCI distribution API

### Risk: Certificate Issues
- **Mitigation**: Leverage Phase 1 fallback strategies
- **Fallback**: --insecure mode with warnings

### Risk: Size Overrun
- **Mitigation**: Already split into 2 parts
- **Fallback**: Further split if needed

## ✅ FINAL CHECKLIST

Before marking complete:
- [ ] All functions implemented (no stubs)
- [ ] Size under 700 lines (main effort)
- [ ] Tests passing with 80% coverage
- [ ] Feature flag tested both ways
- [ ] Can merge independently to main
- [ ] Phase 1 dependencies integrated
- [ ] Documentation complete
- [ ] No security vulnerabilities
- [ ] Performance acceptable
- [ ] Error handling comprehensive

## 📝 NOTES

- This plan builds on the existing implementation
- Splits 001 and 002 already created
- Focus on main effort coordination
- Ensure R307 compliance for independent mergeability
- Apply R320 - no stub implementations allowed
- Follow R338 for standardized size reporting

---
**Plan Created By**: Code Reviewer Agent
**Plan Version**: 2.0 (Updated for rebase requirements)
**Supersedes**: IMPLEMENTATION-PLAN-20250908-001859.md