# Registry Operations Implementation Plan

## CRITICAL EFFORT METADATA (FROM WAVE PLAN)
**Branch**: `idpbuilder-oci-mgmt/phase2/wave2/effort5-registry`
**Can Parallelize**: No
**Parallel With**: None (depends on Efforts 1-4)
**Size Estimate**: 700 lines (MUST be <800)
**Dependencies**: Efforts 1 (contracts), 2 (optimizer), 3 (cache), 4 (security)

## Overview
- **Effort**: Registry Operations & Authentication for OCI image management
- **Phase**: 2, Wave: 2
- **Estimated Size**: 700 lines
- **Implementation Time**: 2 days

## CRITICAL: Self-Signed Certificate Support
This effort has specific requirements for handling **gitea.cnoe.localtest.me** with self-signed certificates:
- MUST support `InsecureSkipVerify` option for TLS
- MUST handle custom CA certificates
- MUST work with local development registries
- MUST maintain security for production registries

## File Structure

### Core Implementation Files
- `pkg/oci/registry/client.go` (200 lines): Main RegistryClient implementation
- `pkg/oci/registry/auth.go` (150 lines): Authentication handling with TLS config
- `pkg/oci/registry/push_pull.go` (180 lines): Push/pull operations with security
- `pkg/oci/registry/manifest.go` (120 lines): Manifest manipulation
- `pkg/oci/registry/transport.go` (50 lines): HTTP transport with retry and TLS

### Test Files
- `pkg/oci/registry/client_test.go`: Registry client tests
- `pkg/oci/registry/auth_test.go`: Authentication tests including self-signed
- `pkg/oci/registry/push_pull_test.go`: Push/pull operation tests
- `pkg/oci/registry/manifest_test.go`: Manifest handling tests
- `pkg/oci/registry/testdata/`: Mock responses and test certificates

## Implementation Steps

### Step 1: Create Registry Client Base (200 lines)
```go
// pkg/oci/registry/client.go
package registry

import (
    "context"
    "crypto/tls"
    "crypto/x509"
    "fmt"
    "io/ioutil"
    "net/http"
    
    "github.com/cnoe-io/idpbuilder/pkg/oci/api"
    "github.com/cnoe-io/idpbuilder/pkg/oci/security"
)

type registryClient struct {
    httpClient    *http.Client
    securityMgr   api.SecurityManager
    defaultAuth   api.AuthConfig
    transportOpts TransportOptions
}

type TransportOptions struct {
    MaxRetries         int
    RetryBackoff       time.Duration
    InsecureSkipVerify bool  // Critical for gitea.cnoe.localtest.me
    CACertPath         string
    ClientCertPath     string
    ClientKeyPath      string
}

func NewRegistryClient(opts ...Option) api.RegistryClient {
    rc := &registryClient{
        transportOpts: TransportOptions{
            MaxRetries:   3,
            RetryBackoff: 2 * time.Second,
        },
    }
    
    // Apply options
    for _, opt := range opts {
        opt(rc)
    }
    
    // Configure TLS for self-signed certificates
    rc.httpClient = rc.createHTTPClient()
    
    return rc
}

func (rc *registryClient) createHTTPClient() *http.Client {
    tlsConfig := &tls.Config{
        InsecureSkipVerify: rc.transportOpts.InsecureSkipVerify,
    }
    
    // Load custom CA if provided
    if rc.transportOpts.CACertPath != "" {
        caCert, err := ioutil.ReadFile(rc.transportOpts.CACertPath)
        if err == nil {
            caCertPool := x509.NewCertPool()
            caCertPool.AppendCertsFromPEM(caCert)
            tlsConfig.RootCAs = caCertPool
        }
    }
    
    // Load client certificates if provided
    if rc.transportOpts.ClientCertPath != "" && rc.transportOpts.ClientKeyPath != "" {
        cert, err := tls.LoadX509KeyPair(
            rc.transportOpts.ClientCertPath,
            rc.transportOpts.ClientKeyPath,
        )
        if err == nil {
            tlsConfig.Certificates = []tls.Certificate{cert}
        }
    }
    
    transport := &http.Transport{
        TLSClientConfig: tlsConfig,
        // Add retry logic wrapper
    }
    
    return &http.Client{
        Transport: transport,
        Timeout:   30 * time.Second,
    }
}
```

### Step 2: Implement Authentication (150 lines)
```go
// pkg/oci/registry/auth.go
package registry

import (
    "encoding/base64"
    "encoding/json"
    "fmt"
    "net/http"
    "strings"
    
    "github.com/cnoe-io/idpbuilder/pkg/oci/api"
)

type authHandler struct {
    config api.AuthConfig
    tokens map[string]*tokenResponse
}

type tokenResponse struct {
    Token       string    `json:"token"`
    AccessToken string    `json:"access_token"`
    ExpiresIn   int       `json:"expires_in"`
    IssuedAt    time.Time `json:"issued_at"`
}

func newAuthHandler(config api.AuthConfig) *authHandler {
    return &authHandler{
        config: config,
        tokens: make(map[string]*tokenResponse),
    }
}

// Handle authentication for gitea.cnoe.localtest.me
func (ah *authHandler) authenticate(req *http.Request, challenge *authChallenge) error {
    // Special handling for local gitea with self-signed cert
    if strings.Contains(ah.config.ServerAddress, "gitea.cnoe.localtest.me") {
        // Use basic auth for local gitea
        if ah.config.Username != "" && ah.config.Password != "" {
            auth := base64.StdEncoding.EncodeToString(
                []byte(fmt.Sprintf("%s:%s", ah.config.Username, ah.config.Password)),
            )
            req.Header.Set("Authorization", fmt.Sprintf("Basic %s", auth))
            return nil
        }
    }
    
    // Standard Docker Registry V2 authentication flow
    if challenge != nil {
        return ah.handleAuthChallenge(req, challenge)
    }
    
    // Basic authentication fallback
    if ah.config.Username != "" && ah.config.Password != "" {
        auth := base64.StdEncoding.EncodeToString(
            []byte(fmt.Sprintf("%s:%s", ah.config.Username, ah.config.Password)),
        )
        req.Header.Set("Authorization", fmt.Sprintf("Basic %s", auth))
    }
    
    // Bearer token if available
    if ah.config.RegistryToken != "" {
        req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", ah.config.RegistryToken))
    }
    
    return nil
}

func (ah *authHandler) handleAuthChallenge(req *http.Request, challenge *authChallenge) error {
    // Parse WWW-Authenticate header
    // Handle Bearer token flow
    // Support OAuth2 if needed
    // Cache tokens for reuse
    return nil
}
```

### Step 3: Implement Push/Pull Operations (180 lines)
```go
// pkg/oci/registry/push_pull.go
package registry

import (
    "bytes"
    "context"
    "crypto/sha256"
    "encoding/hex"
    "fmt"
    "io"
    
    "github.com/cnoe-io/idpbuilder/pkg/oci/api"
    "github.com/cnoe-io/idpbuilder/pkg/oci/security"
)

// Push uploads an image with automatic signing
func (rc *registryClient) Push(ctx context.Context, image string, auth api.AuthConfig) error {
    // Parse image reference
    ref, err := parseImageReference(image)
    if err != nil {
        return fmt.Errorf("invalid image reference: %w", err)
    }
    
    // Setup authentication
    authHandler := newAuthHandler(auth)
    
    // Get image manifest and layers from local store
    manifest, layers, err := rc.getLocalImage(ref)
    if err != nil {
        return fmt.Errorf("failed to get local image: %w", err)
    }
    
    // Upload layers with progress tracking
    for _, layer := range layers {
        if err := rc.uploadLayer(ctx, ref, layer, authHandler); err != nil {
            return fmt.Errorf("failed to upload layer: %w", err)
        }
    }
    
    // Upload manifest
    if err := rc.uploadManifest(ctx, ref, manifest, authHandler); err != nil {
        return fmt.Errorf("failed to upload manifest: %w", err)
    }
    
    // Sign image if security manager is available
    if rc.securityMgr != nil {
        signature, err := rc.securityMgr.SignImage(ctx, image, nil)
        if err != nil {
            return fmt.Errorf("failed to sign image: %w", err)
        }
        // Attach signature as additional layer or attestation
        if err := rc.attachSignature(ctx, ref, signature, authHandler); err != nil {
            return fmt.Errorf("failed to attach signature: %w", err)
        }
    }
    
    return nil
}

// Pull downloads an image with automatic verification
func (rc *registryClient) Pull(ctx context.Context, image string, auth api.AuthConfig) (*api.Image, error) {
    // Parse image reference
    ref, err := parseImageReference(image)
    if err != nil {
        return nil, fmt.Errorf("invalid image reference: %w", err)
    }
    
    // Setup authentication
    authHandler := newAuthHandler(auth)
    
    // Download manifest
    manifest, err := rc.downloadManifest(ctx, ref, authHandler)
    if err != nil {
        return nil, fmt.Errorf("failed to download manifest: %w", err)
    }
    
    // Verify signature if security manager is available
    if rc.securityMgr != nil {
        if err := rc.securityMgr.VerifySignature(ctx, image, nil); err != nil {
            return nil, fmt.Errorf("signature verification failed: %w", err)
        }
    }
    
    // Download layers with progress tracking
    layers := make([]*api.LayerInfo, 0, len(manifest.Layers))
    for _, descriptor := range manifest.Layers {
        layer, err := rc.downloadLayer(ctx, ref, descriptor, authHandler)
        if err != nil {
            return nil, fmt.Errorf("failed to download layer: %w", err)
        }
        layers = append(layers, layer)
    }
    
    // Construct image object
    return &api.Image{
        Name:      ref.Name,
        Tag:       ref.Tag,
        Digest:    manifest.Config.Digest,
        MediaType: manifest.MediaType,
        Layers:    layers,
    }, nil
}

func (rc *registryClient) uploadLayer(ctx context.Context, ref *imageReference, layer *api.LayerInfo, auth *authHandler) error {
    // Check if layer exists (HEAD request)
    // If not, initiate upload (POST)
    // Upload chunks (PATCH)
    // Finalize upload (PUT)
    return nil
}

func (rc *registryClient) downloadLayer(ctx context.Context, ref *imageReference, descriptor *api.Descriptor, auth *authHandler) (*api.LayerInfo, error) {
    // Download layer blob
    // Verify digest
    // Save to local store
    return nil, nil
}
```

### Step 4: Implement Manifest Handling (120 lines)
```go
// pkg/oci/registry/manifest.go
package registry

import (
    "encoding/json"
    "fmt"
    
    "github.com/cnoe-io/idpbuilder/pkg/oci/api"
)

// GetManifest retrieves and parses an image manifest
func (rc *registryClient) GetManifest(ctx context.Context, image string) (*api.Manifest, error) {
    ref, err := parseImageReference(image)
    if err != nil {
        return nil, fmt.Errorf("invalid image reference: %w", err)
    }
    
    // Create request with proper Accept headers
    url := fmt.Sprintf("%s/v2/%s/manifests/%s", 
        rc.getRegistryURL(ref), ref.Name, ref.Tag)
    
    req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
    if err != nil {
        return nil, err
    }
    
    // Accept both OCI and Docker manifest types
    req.Header.Set("Accept", "application/vnd.oci.image.manifest.v1+json")
    req.Header.Add("Accept", "application/vnd.docker.distribution.manifest.v2+json")
    req.Header.Add("Accept", "application/vnd.docker.distribution.manifest.list.v2+json")
    
    // Authenticate request
    authHandler := newAuthHandler(rc.defaultAuth)
    if err := authHandler.authenticate(req, nil); err != nil {
        return nil, fmt.Errorf("authentication failed: %w", err)
    }
    
    resp, err := rc.httpClient.Do(req)
    if err != nil {
        return nil, fmt.Errorf("request failed: %w", err)
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("manifest fetch failed: %s", resp.Status)
    }
    
    var manifest api.Manifest
    if err := json.NewDecoder(resp.Body).Decode(&manifest); err != nil {
        return nil, fmt.Errorf("failed to decode manifest: %w", err)
    }
    
    return &manifest, nil
}

// ListTags returns all tags for a repository
func (rc *registryClient) ListTags(ctx context.Context, repository string) ([]string, error) {
    url := fmt.Sprintf("%s/v2/%s/tags/list", rc.getRegistryBaseURL(), repository)
    
    req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
    if err != nil {
        return nil, err
    }
    
    authHandler := newAuthHandler(rc.defaultAuth)
    if err := authHandler.authenticate(req, nil); err != nil {
        return nil, fmt.Errorf("authentication failed: %w", err)
    }
    
    resp, err := rc.httpClient.Do(req)
    if err != nil {
        return nil, fmt.Errorf("request failed: %w", err)
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("tags list failed: %s", resp.Status)
    }
    
    var result struct {
        Tags []string `json:"tags"`
    }
    
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return nil, fmt.Errorf("failed to decode tags: %w", err)
    }
    
    return result.Tags, nil
}

// Helper to create multi-arch manifest
func (rc *registryClient) createManifestList(manifests []*api.Manifest) (*manifestList, error) {
    // Create manifest list for multi-architecture support
    // Combine individual manifests
    // Set proper media types
    return nil, nil
}
```

### Step 5: Implement Transport with Retry (50 lines)
```go
// pkg/oci/registry/transport.go
package registry

import (
    "context"
    "net/http"
    "time"
)

type retryTransport struct {
    base       http.RoundTripper
    maxRetries int
    backoff    time.Duration
}

func newRetryTransport(base http.RoundTripper, maxRetries int, backoff time.Duration) *retryTransport {
    return &retryTransport{
        base:       base,
        maxRetries: maxRetries,
        backoff:    backoff,
    }
}

func (rt *retryTransport) RoundTrip(req *http.Request) (*http.Response, error) {
    var lastErr error
    
    for i := 0; i <= rt.maxRetries; i++ {
        if i > 0 {
            // Exponential backoff
            time.Sleep(rt.backoff * time.Duration(1<<uint(i-1)))
        }
        
        resp, err := rt.base.RoundTrip(req)
        if err != nil {
            lastErr = err
            continue
        }
        
        // Don't retry on success or client errors
        if resp.StatusCode < 500 {
            return resp, nil
        }
        
        // Server error, close response and retry
        resp.Body.Close()
        lastErr = fmt.Errorf("server error: %d", resp.StatusCode)
    }
    
    return nil, fmt.Errorf("max retries exceeded: %w", lastErr)
}
```

## Dependencies and Imports

### From Effort 1 (Contracts)
```go
import "github.com/cnoe-io/idpbuilder/pkg/oci/api"
```
- Use all interfaces: `RegistryClient`, `SecurityManager`
- Use all models: `AuthConfig`, `Image`, `Manifest`, `Layer`, etc.

### From Effort 4 (Security)
```go
import "github.com/cnoe-io/idpbuilder/pkg/oci/security"
```
- Use for image signing on push
- Use for signature verification on pull
- Integration with SBOM attachment

### From Wave 1 (Build Components)
```go
import "github.com/cnoe-io/idpbuilder/pkg/oci/build"
```
- Reuse runtime configuration
- Access local image storage

## Size Management
- **Estimated Lines**: 700 lines
- **Measurement Tool**: `${PROJECT_ROOT}/tools/line-counter.sh` (find project root first)
- **Check Frequency**: After each file implementation
- **Split Threshold**: 700 lines (warning), 800 lines (stop)

## Test Requirements

### Unit Tests
- **Coverage Target**: 80%
- **Key Test Areas**:
  - Authentication with self-signed certificates
  - Push/pull with retry logic
  - Manifest parsing and creation
  - Error handling and edge cases

### Integration Tests
- Test against local gitea.cnoe.localtest.me registry
- Verify self-signed certificate handling
- Test security integration (signing/verification)
- Multi-arch manifest support

### Test Files Structure
```yaml
test_files:
  - pkg/oci/registry/client_test.go: Main client tests
  - pkg/oci/registry/auth_test.go: Authentication including self-signed
  - pkg/oci/registry/push_pull_test.go: Push/pull operations
  - pkg/oci/registry/manifest_test.go: Manifest handling
  - pkg/oci/registry/testdata/:
    - test-cert.pem: Self-signed certificate for testing
    - test-key.pem: Private key for testing
    - manifest-v2.json: Sample Docker manifest
    - manifest-oci.json: Sample OCI manifest
```

## Pattern Compliance

### Go Patterns
- Use interfaces for all public APIs
- Implement options pattern for configuration
- Use context for cancellation and timeouts
- Proper error wrapping with fmt.Errorf

### Security Patterns
- Always validate TLS certificates in production
- Support InsecureSkipVerify ONLY for development
- Integrate security manager for signing/verification
- Never log sensitive authentication data

### Registry Patterns
- Follow Docker Registry V2 API specification
- Support both OCI and Docker manifest formats
- Implement proper content negotiation
- Handle authentication challenges correctly

## Success Criteria

### Functionality
- Successfully push images to gitea.cnoe.localtest.me
- Successfully pull images with verification
- Handle self-signed certificates correctly
- Support all authentication methods

### Performance
- Push/pull operations complete within reasonable time
- Efficient layer deduplication
- Proper connection pooling
- Retry logic prevents transient failures

### Security
- Images signed on push when security manager available
- Signatures verified on pull
- TLS configuration properly handled
- No credential leakage in logs

## Implementation Order

1. **Day 1 Morning**: 
   - Create client.go with TLS configuration
   - Implement transport.go with retry logic
   
2. **Day 1 Afternoon**:
   - Complete auth.go with self-signed cert support
   - Test authentication against gitea.cnoe.localtest.me
   
3. **Day 2 Morning**:
   - Implement push_pull.go with security integration
   - Add progress tracking and error handling
   
4. **Day 2 Afternoon**:
   - Complete manifest.go with multi-arch support
   - Write comprehensive tests
   - Verify size compliance (<800 lines)

## Critical Implementation Notes

### Self-Signed Certificate Handling
```go
// CRITICAL: For gitea.cnoe.localtest.me development registry
tlsConfig := &tls.Config{
    InsecureSkipVerify: auth.Insecure, // Set to true for local dev
}

// Production registries should use proper CA validation
if auth.CACertPath != "" {
    // Load and validate CA certificate
}
```

### Authentication Priority
1. Check for gitea.cnoe.localtest.me - use basic auth
2. Try bearer token if available
3. Fall back to basic authentication
4. Handle OAuth2 flow if required

### Error Handling
- Distinguish between transient and permanent errors
- Retry on network failures and 5xx errors
- Don't retry on authentication failures (401)
- Provide detailed error messages for debugging

## Risk Mitigation

### Risk: Size Limit Exceeded
**Mitigation**: 
- Target 650 lines implementation
- Monitor size after each file
- Ready to defer advanced features if needed

### Risk: Self-Signed Cert Issues
**Mitigation**:
- Test early with gitea.cnoe.localtest.me
- Provide clear configuration examples
- Include troubleshooting documentation

### Risk: Security Integration Complexity
**Mitigation**:
- Make signing/verification optional
- Graceful degradation if security manager unavailable
- Clear separation of concerns

## Document Metadata
- **Created By**: Code Reviewer Agent
- **Date**: 2025-08-26
- **Phase**: 2, Wave: 2, Effort: 5
- **Dependencies Verified**: Yes (efforts 1-4 complete)
- **Size Validated**: 700 lines (under 800 limit)