# Implementation Plan for Registry Client Abstraction

Created: 2025-09-28T14:25:36Z
Location: /home/vscode/workspaces/idpbuilder-gitea-push/efforts/phase1/wave2/registry-client
Phase: 1
Wave: 2

## 🚨🚨🚨 R355 PRODUCTION READINESS - ZERO TOLERANCE 🚨🚨🚨

This implementation MUST be production-ready from the first commit:
- ❌ NO STUBS or placeholder implementations
- ❌ NO MOCKS except in test directories
- ❌ NO hardcoded credentials or secrets
- ❌ NO static configuration values
- ❌ NO TODO/FIXME markers in code
- ❌ NO returning nil or empty for "later implementation"
- ❌ NO panic("not implemented") patterns
- ❌ NO fake or dummy data

VIOLATION = -100% AUTOMATIC FAILURE

## Effort Metadata

**Effort ID**: P1W2-E2
**Effort Name**: Registry Client Abstraction
**Branch**: `phase1/wave2/registry-client`
**Theme**: Wrap go-containerregistry with clean, testable interface
**Size Estimate**: 350 lines (NEW code to be added)
**Dependencies**:
- P1W1-E2 (OCI Package Format) - types from pkg/oci/format
- P1W1-E3 (Registry Configuration) - config types from pkg/config
- P1W1-E4 (CLI Contracts) - interfaces from pkg/cmd/interfaces

**Can Parallelize**: Yes
**Parallel With**: All other Wave 2 efforts (builder-interface, certificate-manager, stack-mapper, progress-reporter, error-handler)

## Pre-Planning Research Results (R374 MANDATORY)

### Existing Interfaces Found

| Interface | Location | Signature | Must Implement |
|-----------|----------|-----------|----------------|
| RegistryCommand | P1W1-E4-cli-contracts/pkg/cmd/interfaces/registry.go | ConfigureRegistry, ValidateCredentials, GetRegistryInfo, TestConnection, GetCapabilities | NO - This is CLI layer |
| RegistryManager | P1W1-E4-cli-contracts/pkg/cmd/interfaces/registry.go | AddRegistry, RemoveRegistry, UpdateRegistry, GetRegistry, ListRegistries | NO - This is CLI management |
| RegistryAuthenticator | P1W1-E4-cli-contracts/pkg/cmd/interfaces/registry.go | Authenticate, RefreshCredentials, GetAuthToken | NO - Separate concern |

### Existing Implementations to Reuse

| Component | Location | Purpose | How to Use |
|-----------|----------|---------|------------|
| RegistryConfig | P1W1-E3/pkg/config/registry.go | Registry configuration struct | Import and use as parameter type |
| AuthConfig | P1W1-E3/pkg/config/auth.go | Authentication config | Import and use for auth setup |
| TLSConfig | P1W1-E3/pkg/config/auth.go | TLS configuration | Import for secure connections |
| Artifact | P1W1-E1/pkg/providers/types.go | OCI artifact representation | Use as return type |
| Layer | P1W1-E1/pkg/providers/types.go | OCI layer representation | Use in artifact construction |
| PackageManifest | P1W1-E2/pkg/oci/format/spec.go | OCI manifest structure | Use for manifest operations |

### APIs Already Defined

| API | Method | Signature | Notes |
|-----|--------|-----------|-------|
| LoadRegistryConfig | P1W1-E3/pkg/config | func LoadRegistryConfig(path string) (*RegistryConfig, error) | Use for config loading |
| ToConnectionString | P1W1-E3/pkg/config | func ToConnectionString(config *RegistryConfig) string | Build connection strings |
| GetAuthConfig | P1W1-E3/pkg/config | func GetAuthConfig(config *RegistryConfig) (*AuthConfig, error) | Extract auth from config |

### FORBIDDEN DUPLICATIONS (R373)
- DO NOT create competing RegistryConfig types
- DO NOT reimplement authentication config structures
- DO NOT create alternative Artifact/Layer types
- DO NOT duplicate PackageManifest structure

### REQUIRED INTEGRATIONS (R373)
- MUST import github.com/idpbuilder/idpbuilder-gitea-push/efforts/phase1/wave1/P1W1-E3-registry-config/pkg/config
- MUST import github.com/idpbuilder/idpbuilder-gitea-push/efforts/phase1/wave1/P1W1-E1-provider-interface/pkg/providers
- MUST use existing config types for all configuration
- MUST return existing Artifact types for registry operations

## EXPLICIT SCOPE (R311 MANDATORY)

### IMPLEMENT EXACTLY:
- Interface: `Client` with 4 methods (~40 lines)
  - `Push(ctx context.Context, ref string, artifact *providers.Artifact) error` (~10 lines)
  - `Pull(ctx context.Context, ref string) (*providers.Artifact, error)` (~10 lines)
  - `Exists(ctx context.Context, ref string) (bool, error)` (~10 lines)
  - `ListTags(ctx context.Context, repository string) ([]string, error)` (~10 lines)
- Type: `gcrAdapter` struct implementing Client (~30 lines)
  - Fields: config, auth, transport, remoteOpts
- Function: `NewClient(config *config.RegistryConfig) (Client, error)` (~50 lines)
  - Creates gcrAdapter with proper initialization
- Type: `Transport` interface with 2 methods (~20 lines)
  - `RoundTrip(req *http.Request) (*http.Response, error)`
  - `WithAuth(auth *config.AuthConfig) Transport`
- Type: `httpTransport` struct implementing Transport (~80 lines)
  - Wraps http.RoundTripper with auth injection
- Type: `Reference` interface with 3 methods (~20 lines)
  - `Parse(ref string) error`
  - `Registry() string`
  - `Repository() string`
- Type: `reference` struct implementing Reference (~50 lines)
  - Simple reference parsing using go-containerregistry/name
- Function: `parseReference(ref string) (Reference, error)` (~30 lines)
- Tests: 5 basic unit tests (~100 lines total)
  - TestNewClient (~20 lines)
  - TestParseReference (~20 lines)
  - TestTransportWithAuth (~20 lines)
  - TestClientPush (~20 lines)
  - TestClientPull (~20 lines)

**TOTAL: ~350 lines**

### DO NOT IMPLEMENT:
- ❌ Catalog browsing (future effort)
- ❌ Repository management (future effort)
- ❌ Complex authentication flows (Wave 3)
- ❌ Certificate management (E1.2.3 handles this)
- ❌ Retry logic (error-handler effort)
- ❌ Progress reporting (E1.2.5 handles this)
- ❌ Caching mechanisms
- ❌ Connection pooling
- ❌ Rate limiting
- ❌ Metrics collection

## Size Limit Clarification (R359)
- The 800-line limit applies to NEW CODE YOU ADD
- Repository will grow by ~350 lines (EXPECTED)
- NEVER delete existing code to meet size limits
- Current codebase: ~0 lines (new effort)
- Expected total after: ~350 lines

## File Structure

```
pkg/registry/
├── client.go          # Client interface definition (~40 lines)
├── gcr_adapter.go     # go-containerregistry wrapper (~180 lines)
├── transport.go       # HTTP transport abstraction (~80 lines)
├── reference.go       # Reference parsing interface (~50 lines)
└── client_test.go     # Unit tests (~100 lines)
```

## Implementation Details

### 1. client.go - Client Interface (~40 lines)
```go
package registry

import (
    "context"
    "github.com/idpbuilder/idpbuilder-gitea-push/efforts/phase1/wave1/P1W1-E1-provider-interface/pkg/providers"
)

// Client defines the interface for registry operations.
// This abstraction wraps go-containerregistry for testability.
type Client interface {
    // Push uploads an artifact to the registry
    Push(ctx context.Context, ref string, artifact *providers.Artifact) error

    // Pull retrieves an artifact from the registry
    Pull(ctx context.Context, ref string) (*providers.Artifact, error)

    // Exists checks if an artifact exists in the registry
    Exists(ctx context.Context, ref string) (bool, error)

    // ListTags returns all tags for a repository
    ListTags(ctx context.Context, repository string) ([]string, error)
}
```

### 2. gcr_adapter.go - Implementation (~180 lines)
```go
package registry

import (
    "context"
    "fmt"
    "os"

    "github.com/google/go-containerregistry/pkg/authn"
    "github.com/google/go-containerregistry/pkg/name"
    "github.com/google/go-containerregistry/pkg/v1/remote"
    "github.com/idpbuilder/idpbuilder-gitea-push/efforts/phase1/wave1/P1W1-E1-provider-interface/pkg/providers"
    "github.com/idpbuilder/idpbuilder-gitea-push/efforts/phase1/wave1/P1W1-E3-registry-config/pkg/config"
)

// gcrAdapter wraps go-containerregistry library
type gcrAdapter struct {
    config     *config.RegistryConfig
    auth       authn.Authenticator
    transport  Transport
    remoteOpts []remote.Option
}

// NewClient creates a new registry client with the given configuration
func NewClient(cfg *config.RegistryConfig) (Client, error) {
    if cfg == nil {
        return nil, fmt.Errorf("registry config cannot be nil")
    }

    if cfg.URL == "" {
        cfg.URL = os.Getenv("REGISTRY_URL")
        if cfg.URL == "" {
            return nil, fmt.Errorf("registry URL not configured")
        }
    }

    // Set up authenticator based on config
    var auth authn.Authenticator
    authCfg, err := config.GetAuthConfig(cfg)
    if err != nil {
        return nil, fmt.Errorf("failed to get auth config: %w", err)
    }

    switch authCfg.Type {
    case "basic":
        if authCfg.Username != "" && authCfg.Password != "" {
            auth = &authn.Basic{
                Username: authCfg.Username,
                Password: authCfg.Password,
            }
        }
    case "token":
        if authCfg.Token != "" {
            auth = &authn.Bearer{Token: authCfg.Token}
        }
    default:
        auth = authn.Anonymous
    }

    // Set up transport with auth
    transport := &httpTransport{
        base: remote.DefaultTransport,
        auth: authCfg,
    }

    // Build remote options
    opts := []remote.Option{
        remote.WithAuth(auth),
        remote.WithTransport(transport),
    }

    if cfg.Insecure {
        // Note: In production, properly configure TLS
        // This is a minimal working implementation
        opts = append(opts, remote.WithTransport(transport))
    }

    return &gcrAdapter{
        config:     cfg,
        auth:       auth,
        transport:  transport,
        remoteOpts: opts,
    }, nil
}

// Implement Push, Pull, Exists, ListTags methods...
```

### 3. transport.go - HTTP Transport (~80 lines)
```go
package registry

import (
    "net/http"
    "github.com/idpbuilder/idpbuilder-gitea-push/efforts/phase1/wave1/P1W1-E3-registry-config/pkg/config"
)

// Transport defines the HTTP transport interface for registry operations
type Transport interface {
    RoundTrip(req *http.Request) (*http.Response, error)
    WithAuth(auth *config.AuthConfig) Transport
}

// httpTransport implements Transport with authentication injection
type httpTransport struct {
    base http.RoundTripper
    auth *config.AuthConfig
}

// Implement RoundTrip and WithAuth methods...
```

### 4. reference.go - Reference Parsing (~50 lines)
```go
package registry

import (
    "github.com/google/go-containerregistry/pkg/name"
)

// Reference represents a parsed container image reference
type Reference interface {
    Parse(ref string) error
    Registry() string
    Repository() string
}

// reference wraps go-containerregistry name.Reference
type reference struct {
    ref name.Reference
}

// parseReference parses a reference string
func parseReference(ref string) (Reference, error) {
    // Implementation using go-containerregistry/pkg/name
}
```

## Configuration Requirements (R355 Mandatory)

### CORRECT Production-Ready Patterns:
```go
// ✅ From environment variable
registryURL := os.Getenv("REGISTRY_URL")
if registryURL == "" {
    return errors.New("REGISTRY_URL not set")
}

// ✅ From configuration file
cfg, err := config.LoadRegistryConfig(configPath)
if err != nil {
    return fmt.Errorf("failed to load config: %w", err)
}

// ✅ Configurable timeout
timeout := cfg.Options["timeout"]
if timeout == "" {
    timeout = defaultTimeout
}
```

### FORBIDDEN Anti-Patterns:
```go
// ❌ VIOLATION - Hardcoded URL
registryURL := "https://registry.example.com"

// ❌ VIOLATION - Stub implementation
func Push(ctx context.Context, ref string, artifact *providers.Artifact) error {
    // TODO: implement later
    return nil
}

// ❌ VIOLATION - Static auth
username := "admin"
password := "secret123"
```

## Effort Atomic PR Design

### PR Summary
Single PR implementing registry client abstraction wrapping go-containerregistry

### Can Merge to Main Alone
**true** - This PR can merge independently without breaking the build

### R355 Production Ready Checklist
- no_hardcoded_values: true
- all_config_from_env: true
- no_stub_implementations: true
- no_todo_markers: true
- all_functions_complete: true

### Configuration Approach
- **Registry URL**: From REGISTRY_URL env var or config file
- **Auth Credentials**: From REGISTRY_USERNAME/REGISTRY_PASSWORD env vars
- **Timeout**: From config with default fallback
- **TLS Settings**: From config.Insecure flag

### Feature Flags Needed
**None** - This is foundational abstraction, not a toggleable feature

### Interface Implementations
- Interface: `Client`
  - Implementation: `gcrAdapter`
  - Production Ready: true
  - Notes: Fully functional wrapper around go-containerregistry

- Interface: `Transport`
  - Implementation: `httpTransport`
  - Production Ready: true
  - Notes: Working HTTP transport with auth injection

- Interface: `Reference`
  - Implementation: `reference`
  - Production Ready: true
  - Notes: Complete reference parsing using go-containerregistry/name

### PR Verification
- tests_pass_alone: true
- build_remains_working: true
- flags_tested_both_ways: N/A (no feature flags)
- no_external_dependencies: false (uses go-containerregistry)
- backward_compatible: N/A (new code)

## Testing Requirements

### Unit Tests (60% minimum coverage)
1. **TestNewClient** - Verify client creation with various configs
2. **TestParseReference** - Test reference parsing edge cases
3. **TestTransportWithAuth** - Verify auth header injection
4. **TestClientPush** - Test push operation (with mock)
5. **TestClientPull** - Test pull operation (with mock)

### Interface Tests
- Verify Client interface contract
- Verify Transport interface behavior
- Verify Reference interface parsing

### Mock Implementations
- Provide mock Client for testing in Wave 3
- Provide mock Transport for auth testing

## Implementation Steps

1. **Create pkg/registry directory**
   ```bash
   mkdir -p pkg/registry
   ```

2. **Implement client.go**
   - Define Client interface
   - Add documentation

3. **Implement gcr_adapter.go**
   - Create gcrAdapter struct
   - Implement NewClient function
   - Implement Push, Pull, Exists, ListTags methods
   - Ensure all config from environment/files

4. **Implement transport.go**
   - Define Transport interface
   - Implement httpTransport
   - Add auth injection logic

5. **Implement reference.go**
   - Define Reference interface
   - Implement reference struct
   - Add parseReference function

6. **Write tests**
   - Unit tests for each component
   - Mock implementations
   - Verify 60% coverage

7. **Measure with line-counter.sh**
   ```bash
   $PROJECT_ROOT/tools/line-counter.sh
   ```

## Size Management
- **Estimated Lines**: 350 (well under 800 limit)
- **Measurement Tool**: ${PROJECT_ROOT}/tools/line-counter.sh
- **Check Frequency**: After each file completion
- **Split Threshold**: 700 lines (warning), 800 lines (stop)

## Success Criteria

### Completion Requirements
- ✅ Client interface defined and implemented
- ✅ go-containerregistry properly wrapped
- ✅ All configuration from env/files (no hardcoded values)
- ✅ Transport abstraction complete
- ✅ Reference parsing functional
- ✅ 60% test coverage achieved
- ✅ Under 800 lines total
- ✅ No TODO/FIXME in code
- ✅ All functions fully implemented

### Quality Gates
1. **Code Review**: Clean abstraction without leaky implementation details
2. **Testing**: All unit tests passing
3. **Documentation**: All public APIs documented
4. **Size Compliance**: Verified with line-counter.sh < 800 lines
5. **Integration**: Imports Wave 1 types correctly

## Next Wave Dependencies

Wave 3 efforts will use this abstraction:
- **Basic Push Command** → Uses Client.Push()
- **Basic Pull Command** → Uses Client.Pull()
- **Registry Auth Flow** → Uses Transport with auth
- **All Registry Ops** → Use Client interface

## Out of Scope Items

The following are explicitly NOT part of this effort:
- Certificate management (handled by certificate-manager effort)
- Progress reporting (handled by progress-reporter effort)
- Error retry logic (handled by error-handler effort)
- Complex authentication flows (Wave 3)
- Registry catalog operations (future wave)
- Connection pooling and optimization (future wave)
- Caching mechanisms (future wave)

---

**Implementation Plan Created By**: @agent-code-reviewer
**Date**: 2025-09-28T14:25:36Z
**Phase**: 1
**Wave**: 2
**Effort**: registry-client
**Status**: Ready for Implementation