# Gitea Registry Client Implementation Plan

**Effort ID**: E2.1.2  
**Effort Name**: gitea-registry-client  
**Phase**: 2 - Build & Push Implementation  
**Wave**: 1 - Core OCI Operations  
**Created By**: Code Reviewer (code-reviewer)  
**Date Created**: 2025-09-02 22:52:06 UTC  

## 🚨 CRITICAL EFFORT METADATA (FROM WAVE PLAN)
**Branch**: `idpbuilder-oci-go-cr/phase2/wave1/gitea-registry-client`  
**Can Parallelize**: Yes  
**Parallel With**: E2.1.1 (go-containerregistry-image-builder)  
**Size Estimate**: 600 lines  
**Dependencies**: Phase 1 TrustStoreManager, E2.1.1 interfaces  

<!-- ⚠️ ORCHESTRATOR METADATA PLACEHOLDER - DO NOT REMOVE ⚠️ -->
<!-- The orchestrator will add infrastructure metadata below: -->
<!-- WORKING_DIRECTORY: efforts/phase2/wave1/gitea-registry-client -->
<!-- BRANCH: idpbuilder-oci-go-cr/phase2/wave1/gitea-registry-client -->
<!-- REMOTE: origin -->
<!-- BASE_BRANCH: idpbuilder-oci-go-cr/phase1-integration-20250902-194557 -->
<!-- SW Engineers MUST read this metadata to navigate to the correct directory -->
<!-- END PLACEHOLDER -->

## 🎯 Effort Overview

### Description
This effort implements a Gitea-specific registry client for OCI image push and pull operations. It integrates with Phase 1's certificate infrastructure to handle self-signed certificates properly, provides authentication mechanisms for Gitea's container registry, and enables seamless image pushing with progress tracking. This is a critical component for the IDPBuilder's ability to push locally built images to the internal Gitea registry at gitea.cnoe.localtest.me.

### Size Estimate
- **Estimated Lines**: 600 lines (within limit)
- **Actual Measurement Method**: `$PROJECT_ROOT/tools/line-counter.sh`
- **Check Frequency**: Every 200 lines
- **Warning Threshold**: 550 lines
- **Stop Threshold**: 800 lines

### Implementation Time
- **Estimated Duration**: 8-10 hours
- **Parallelization**: Can be developed simultaneously with E2.1.1

## 📁 File Structure

```
efforts/phase2/wave1/gitea-registry-client/
├── pkg/
│   ├── registry/
│   │   ├── client.go          (150 lines) - Registry client interface definition
│   │   ├── gitea_client.go    (200 lines) - Gitea-specific implementation
│   │   ├── auth.go            (80 lines)  - Authentication handling
│   │   ├── transport.go       (100 lines) - TLS transport configuration with Phase 1 integration
│   │   └── options.go         (40 lines)  - Client configuration options
│   └── registry/tests/
│       └── gitea_client_test.go (30 lines) - Unit tests for Gitea client
└── go.mod                      - Module dependencies
```

## 🔗 Dependencies

### Phase 1 Dependencies (R219 Compliance)
From analyzing Phase 1 efforts, this effort depends on:

1. **TrustStoreManager** (from phase1/phase-integration-workspace/pkg/certs/)
   - `NewTrustStoreManager()` - Create trust store
   - `AddCertificate()` - Add trusted certificates
   - `ConfigureTransport()` - Configure HTTP transport with certificates
   - `CreateHTTPClient()` - Create HTTP client with certificate handling
   - `SetInsecureRegistry()` - Handle insecure mode when needed

2. **KindCertExtractor** (optional, for testing)
   - `ExtractGiteaCert()` - Extract Gitea certificates from Kind cluster

### External Dependencies
```go
require (
    github.com/google/go-containerregistry v0.19.0
    github.com/stretchr/testify v1.9.0
)
```

### Import Paths
```go
// Phase 1 certificate infrastructure
import "github.com/cnoe-io/idpbuilder/efforts/phase1/phase-integration-workspace/pkg/certs"

// go-containerregistry for OCI operations
import "github.com/google/go-containerregistry/pkg/v1"
import "github.com/google/go-containerregistry/pkg/v1/remote"
import "github.com/google/go-containerregistry/pkg/authn"
import "github.com/google/go-containerregistry/pkg/name"
```

## 🛠️ Implementation Sequence

### Step 1: Define Registry Client Interface (client.go)
**File**: `pkg/registry/client.go`
**Lines**: ~150

```go
package registry

import (
    "context"
    v1 "github.com/google/go-containerregistry/pkg/v1"
)

// Client defines the interface for registry operations
type Client interface {
    // Push pushes an image to the registry
    Push(ctx context.Context, image v1.Image, ref string, opts PushOptions) error
    
    // Pull pulls an image from the registry
    Pull(ctx context.Context, ref string, opts PullOptions) (v1.Image, error)
    
    // Catalog lists repositories in the registry
    Catalog(ctx context.Context) ([]string, error)
    
    // Tags lists tags for a repository
    Tags(ctx context.Context, repository string) ([]string, error)
}

// PushOptions contains options for pushing images
type PushOptions struct {
    Insecure bool           // Skip TLS verification
    Progress ProgressFunc   // Progress callback
}

// PullOptions contains options for pulling images
type PullOptions struct {
    Platform *v1.Platform  // Target platform
    Insecure bool          // Skip TLS verification
}

// ProgressFunc reports progress during operations
type ProgressFunc func(current, total int64)
```

**Key Implementation Points**:
- Define clear interfaces for registry operations
- Support both push and pull operations
- Include progress reporting capabilities
- Provide options for insecure mode (R307 - feature flag)

### Step 2: Implement Gitea Client (gitea_client.go)
**File**: `pkg/registry/gitea_client.go`
**Lines**: ~200

```go
package registry

import (
    "context"
    "fmt"
    "net/http"
    
    "github.com/cnoe-io/idpbuilder/efforts/phase1/phase-integration-workspace/pkg/certs"
    v1 "github.com/google/go-containerregistry/pkg/v1"
    "github.com/google/go-containerregistry/pkg/authn"
    "github.com/google/go-containerregistry/pkg/name"
    "github.com/google/go-containerregistry/pkg/v1/remote"
)

// GiteaClient implements Client for Gitea registries
type GiteaClient struct {
    baseURL    string
    username   string
    password   string
    trustStore certs.TrustStoreManager // Phase 1 integration
    transport  http.RoundTripper
    auth       authn.Authenticator
}

// NewGiteaClient creates a new Gitea registry client
func NewGiteaClient(baseURL, username, password string, trustStore certs.TrustStoreManager) (*GiteaClient, error) {
    // Implementation details...
}

// Push pushes an image to Gitea registry
func (c *GiteaClient) Push(ctx context.Context, image v1.Image, ref string, opts PushOptions) error {
    // Parse reference
    // Configure transport with Phase 1 certificates
    // Use go-containerregistry remote.Write
    // Handle progress reporting
}

// Pull pulls an image from Gitea registry
func (c *GiteaClient) Pull(ctx context.Context, ref string, opts PullOptions) (v1.Image, error) {
    // Parse reference
    // Configure transport
    // Use go-containerregistry remote.Image
}
```

**Key Implementation Points**:
- Integrate Phase 1's TrustStoreManager for certificate handling
- Use go-containerregistry's remote package for OCI operations
- Support authentication with username/password
- Proper error handling and context propagation

### Step 3: Authentication Handling (auth.go)
**File**: `pkg/registry/auth.go`
**Lines**: ~80

```go
package registry

import (
    "encoding/base64"
    "fmt"
    "github.com/google/go-containerregistry/pkg/authn"
)

// ConfigureAuth configures authentication for registry operations
func ConfigureAuth(username, password string) authn.Authenticator {
    if username == "" && password == "" {
        return authn.Anonymous
    }
    
    return &authn.Basic{
        Username: username,
        Password: password,
    }
}

// GetAuthToken generates a base64 encoded auth token
func GetAuthToken(username, password string) string {
    auth := fmt.Sprintf("%s:%s", username, password)
    return base64.StdEncoding.EncodeToString([]byte(auth))
}

// ValidateCredentials checks if credentials are valid
func ValidateCredentials(username, password string) error {
    if username == "" {
        return fmt.Errorf("username cannot be empty")
    }
    if password == "" {
        return fmt.Errorf("password cannot be empty")
    }
    return nil
}
```

**Key Implementation Points**:
- Support both authenticated and anonymous access
- Generate proper auth tokens for Gitea
- Validate credentials before use

### Step 4: Transport Configuration (transport.go)
**File**: `pkg/registry/transport.go`
**Lines**: ~100

```go
package registry

import (
    "crypto/tls"
    "fmt"
    "net/http"
    "time"
    
    "github.com/cnoe-io/idpbuilder/efforts/phase1/phase-integration-workspace/pkg/certs"
)

// configureTransport sets up HTTP transport with certificate handling
func (c *GiteaClient) configureTransport() error {
    // Use Phase 1's TrustStoreManager to create HTTP client
    httpClient, err := c.trustStore.CreateHTTPClient(c.baseURL)
    if err != nil {
        return fmt.Errorf("failed to create HTTP client: %w", err)
    }
    
    c.transport = httpClient.Transport
    return nil
}

// createInsecureTransport creates transport that skips TLS verification
func createInsecureTransport() *http.Transport {
    return &http.Transport{
        TLSClientConfig: &tls.Config{
            InsecureSkipVerify: true,
        },
        Proxy:                 http.ProxyFromEnvironment,
        MaxIdleConns:          100,
        IdleConnTimeout:       90 * time.Second,
        TLSHandshakeTimeout:   10 * time.Second,
        ExpectContinueTimeout: 1 * time.Second,
    }
}

// createSecureTransport creates transport with proper certificate validation
func createSecureTransport(certPool *tls.CertPool) *http.Transport {
    return &http.Transport{
        TLSClientConfig: &tls.Config{
            RootCAs: certPool,
        },
        Proxy:                 http.ProxyFromEnvironment,
        MaxIdleConns:          100,
        IdleConnTimeout:       90 * time.Second,
        TLSHandshakeTimeout:   10 * time.Second,
        ExpectContinueTimeout: 1 * time.Second,
    }
}
```

**Key Implementation Points**:
- Leverage Phase 1's certificate infrastructure
- Support both secure and insecure modes
- Configure proper timeouts and connection pooling
- Handle proxy settings from environment

### Step 5: Client Options (options.go)
**File**: `pkg/registry/options.go`
**Lines**: ~40

```go
package registry

import "time"

// ClientOptions contains configuration for the Gitea client
type ClientOptions struct {
    BaseURL         string
    Username        string
    Password        string
    Insecure        bool
    Timeout         time.Duration
    MaxRetries      int
    RetryDelay      time.Duration
    UserAgent       string
}

// DefaultClientOptions returns default client options
func DefaultClientOptions() ClientOptions {
    return ClientOptions{
        BaseURL:    "https://gitea.cnoe.localtest.me:443",
        Username:   "gitea_admin",
        Timeout:    30 * time.Second,
        MaxRetries: 3,
        RetryDelay: 1 * time.Second,
        UserAgent:  "idpbuilder-oci/1.0",
    }
}

// Validate checks if options are valid
func (o ClientOptions) Validate() error {
    if o.BaseURL == "" {
        return fmt.Errorf("base URL cannot be empty")
    }
    if o.Timeout <= 0 {
        o.Timeout = 30 * time.Second
    }
    return nil
}
```

### Step 6: Unit Tests (gitea_client_test.go)
**File**: `pkg/registry/tests/gitea_client_test.go`
**Lines**: ~30

```go
package registry_test

import (
    "context"
    "testing"
    
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    
    "github.com/cnoe-io/idpbuilder/efforts/phase2/wave1/gitea-registry-client/pkg/registry"
)

func TestGiteaClientCreation(t *testing.T) {
    // Test client creation with valid options
    // Test client creation with invalid options
    // Test authentication configuration
}

func TestPushOperation(t *testing.T) {
    // Mock HTTP transport
    // Test successful push
    // Test push with authentication failure
    // Test push with network error
}

func TestPullOperation(t *testing.T) {
    // Mock HTTP transport
    // Test successful pull
    // Test pull of non-existent image
}
```

## 🧪 Testing Requirements

### Unit Test Coverage (80% minimum)
- Client creation and configuration
- Authentication handling
- Transport configuration with certificates
- Push/pull operations with mocked transport
- Error handling for various failure scenarios

### Integration Test Scenarios
```go
// Integration test with actual Gitea registry (requires test environment)
func TestGiteaIntegration(t *testing.T) {
    // Skip if not in integration test mode
    // Extract certificates from Kind cluster (Phase 1)
    // Create client with real certificates
    // Push test image
    // Pull image back and verify
}
```

### Test Coverage Areas
- **Authentication**: Valid/invalid credentials, token generation
- **Certificate Handling**: Phase 1 integration, insecure mode
- **Network Operations**: Timeouts, retries, connection failures
- **Image Operations**: Push, pull, catalog, tags
- **Error Scenarios**: Network errors, auth failures, invalid references

## 🔒 Security Requirements

### Certificate Validation (R307 - Feature Flags)
```go
// Feature flag for insecure mode
const (
    EnvInsecureRegistry = "IDPBUILDER_INSECURE_REGISTRY"
)

// Check feature flag
if os.Getenv(EnvInsecureRegistry) == "true" {
    opts.Insecure = true
}
```

### Authentication Security
- Never log credentials
- Support environment variables for passwords
- Clear credentials from memory after use
- Use secure token transmission

### TLS Configuration
- Default to secure mode with certificate validation
- Require explicit flag for insecure mode
- Log warnings when operating in insecure mode
- Integrate Phase 1's certificate validation

## 📏 Size Management

### Monitoring Strategy
```bash
# Run measurement every 200 lines of development
cd /home/vscode/workspaces/idpbuilder-oci-go-cr/efforts/phase2/wave1/gitea-registry-client

# Find project root
PROJECT_ROOT=$(pwd)
while [ "$PROJECT_ROOT" != "/" ]; do
    if [ -f "$PROJECT_ROOT/orchestrator-state.yaml" ]; then break; fi
    PROJECT_ROOT=$(dirname "$PROJECT_ROOT")
done

# Measure with correct base branch (R308)
BASE_BRANCH="idpbuilder-oci-go-cr/phase1-integration-20250902-194557"
CURRENT_BRANCH=$(git branch --show-current)
$PROJECT_ROOT/tools/line-counter.sh -b $BASE_BRANCH -c $CURRENT_BRANCH
```

### Size Checkpoints
- After Step 2 (client implementation): ~350 lines
- After Step 4 (transport configuration): ~530 lines
- After Step 6 (tests): ~600 lines
- **Warning at 550 lines**
- **Stop at 700 lines**

## 🎯 Pattern Compliance

### Error Handling Pattern
```go
// Always wrap errors with context
if err != nil {
    return fmt.Errorf("operation failed: %w", err)
}

// Certificate-specific errors
if cert == nil && !opts.Insecure {
    return fmt.Errorf("certificate required (use --insecure to skip)")
}
```

### Logging Pattern
```go
// Use structured logging
log.WithFields(log.Fields{
    "registry": c.baseURL,
    "ref":      ref,
    "insecure": opts.Insecure,
}).Info("Pushing image to registry")
```

### Phase 1 Integration Pattern
```go
// Always use Phase 1 interfaces
trustStore := certs.NewTrustStoreManager()
trustStore.AddCertificate(registryURL, cert)
transportOpt, err := trustStore.ConfigureTransport(registryURL)
```

## 🚀 Integration Points

### With E2.1.1 (go-containerregistry-image-builder)
- Uses v1.Image interface from builder
- Builder creates images, this client pushes them
- Shared v1.Platform specifications

### With Phase 1 Certificate Infrastructure
- Import: `github.com/cnoe-io/idpbuilder/efforts/phase1/phase-integration-workspace/pkg/certs`
- Use TrustStoreManager for all certificate operations
- Leverage existing certificate extraction and validation

### With Future CLI Commands (E2.2.1)
- This client will be used by the push command
- Provides progress reporting for CLI output
- Returns proper error messages for user display

## ✅ Completion Criteria

### Functional Requirements
- [ ] Client can authenticate with Gitea registry
- [ ] Push operations work with self-signed certificates
- [ ] Pull operations retrieve images correctly
- [ ] Catalog and tags operations return correct data
- [ ] Insecure mode works when explicitly enabled

### Quality Requirements
- [ ] 80% unit test coverage achieved
- [ ] All tests pass
- [ ] No security vulnerabilities
- [ ] Proper error handling throughout
- [ ] Integration with Phase 1 verified

### Size Requirements
- [ ] Total implementation under 600 lines
- [ ] Measurement done with designated tool
- [ ] No manual line counting used

## 🔄 Rollback Strategy

If implementation encounters issues:

1. **Certificate Issues**: Fall back to Phase 1's insecure mode temporarily
2. **Authentication Issues**: Use docker CLI with credential helpers as fallback
3. **Size Violations**: Pre-identified split points:
   - Split 1: Client interface and basic implementation (300 lines)
   - Split 2: Authentication and transport (300 lines)

## 📝 Notes for Software Engineers

### Getting Started
1. Checkout branch: `idpbuilder-oci-go-cr/phase2/wave1/gitea-registry-client`
2. Review Phase 1 certificate code in `efforts/phase1/phase-integration-workspace/pkg/certs/`
3. Install go-containerregistry: `go get github.com/google/go-containerregistry`
4. Follow the file structure exactly as specified
5. Measure frequently with `$PROJECT_ROOT/tools/line-counter.sh`

### Key Implementation Tips
- Always use Phase 1's TrustStoreManager, don't reimplement certificate handling
- Use go-containerregistry's remote package for all OCI operations
- Test with actual Gitea registry in Kind cluster when possible
- Keep authentication simple - basic auth is sufficient for MVP
- Focus on push operation first, it's the critical path

### Common Pitfalls to Avoid
- Don't bypass Phase 1's certificate infrastructure
- Don't forget to handle context cancellation
- Don't log sensitive information (passwords, tokens)
- Don't exceed 600 lines total
- Don't implement features not in the plan

### Resources
- Phase 1 Certificate Code: `efforts/phase1/phase-integration-workspace/pkg/certs/`
- go-containerregistry docs: https://pkg.go.dev/github.com/google/go-containerregistry
- Gitea Registry API: https://docs.gitea.io/en-us/packages/container/
- OCI Distribution Spec: https://github.com/opencontainers/distribution-spec

---

*This implementation plan provides concrete, executable instructions for implementing the Gitea registry client with proper Phase 1 integration and strict size compliance.*