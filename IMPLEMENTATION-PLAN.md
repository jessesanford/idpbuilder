<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
# Phase 2 Wave 2 Integration Plan

## INTEGRATION WORKSPACE OVERVIEW
**Purpose**: Merge all Wave 2 efforts into integrated build system
**Integration Branch**: `idpbuilder-oci-mgmt/phase2/wave2-integration`
**Status**: RECOVERY - Completing incomplete integration

## COMPLETED INTEGRATIONS
### Effort 1: Advanced Build Contracts & Interfaces ✓
- **Branch**: `idpbuidler-oci-mgmt/phase2/wave2/effort1-contracts`
- **Status**: ✅ MERGED
- **Files**: Core API contracts in pkg/oci/api/

## INTEGRATIONS IN PROGRESS

### Effort 2: Multi-Stage Build Optimizer (Split Implementation)
- **Split 001**: `idpbuilder-oci-mgmt/phase2/wave2/effort2-optimizer-split-001` (728 lines) ✅ MERGED
  - Core optimizer with optimized analyzer
  - Fixed compilation issues with stub types
- **Split 002**: `idpbuilder-oci-mgmt/phase2/wave2/effort2-optimizer-split-002` (350 lines) 🔄 MERGING
  - Full Executor and GraphBuilder implementation
  - Completes the optimizer implementation

### Effort 3: Cache Manager
- **Branch**: `idpbuidler-oci-mgmt/phase2/wave2/effort3-cache` (834 lines) 🔄 MERGING
- **Purpose**: Layer caching operations and optimization
- **Note**: Originally exceeded size limit (834 > 800) but implementing as-is for integration

### Effort 4: Security Manager (Split Implementation)  
- **Split 001**: `idpbuilder-oci-mgmt/phase2/wave2/effort4-security-split-001` (809 lines) 🔄 MERGING
  - Security orchestration and policy management
  - Note: Slightly over original 762 line estimate
- **Split 002**: `idpbuilder-oci-mgmt/phase2/wave2/effort4-security-split-002` (744 lines) 🔄 MERGING
  - Cryptographic operations layer (signer/verifier implementations)
  - Foundational crypto layer for split-001

### Effort 5: Registry Client
- **Branch**: `idpbuidler-oci-mgmt/phase2/wave2/effort5-registry` (793 lines) 🔄 MERGING
- **Purpose**: Registry operations with self-signed certificate support
- **Key Feature**: Specialized handling for gitea.cnoe.localtest.me

## INTEGRATION STRATEGY

### Sequential Merge Order
1. effort2-optimizer-split-001 (IN PROGRESS)
2. effort2-optimizer-split-002  
3. effort3-cache
4. effort4-security-split-001
5. effort4-security-split-002
6. effort5-registry

### Conflict Resolution
- Merge conflicts expected in IMPLEMENTATION-PLAN.md and work-log.md
- Preserve integration workspace structure
- Maintain effort-specific details in logs
- Ensure no code functionality conflicts

### Integration Verification
- All packages properly structured under `/pkg/oci/`
- Cross-effort dependencies resolved
- No circular dependencies
- Compilation successful
- Tests passing

## FILE STRUCTURE POST-INTEGRATION
```
pkg/
├── oci/
│   ├── api/          # Effort 1 contracts
│   ├── optimizer/    # Effort 2 implementation  
│   ├── cache/        # Effort 3 implementation
│   ├── security/     # Effort 4 implementation
│   └── registry/     # Effort 5 implementation
└── k8s/              # Wave 1 integration
```
=======
# Implementation Plan: Multi-Stage Build Optimizer - Split 002

## 🎯 Effort Overview
**Effort ID**: effort2-optimizer-split-002
**Target Size**: 350 lines MAXIMUM
**Purpose**: Complete Executor and GraphBuilder implementations

## 🚨 CRITICAL REQUIREMENTS
1. **SIZE LIMIT**: 350 lines HARD LIMIT
2. **MUST INTEGRATE**: Work with split-001's interfaces
3. **COMPLETE FUNCTIONALITY**: Implement all stub methods from split-001

## 📁 Files to Implement

### 1. pkg/oci/optimizer/executor.go (~180 lines)
**Purpose**: Parallel execution engine for build stages

**Required Implementation**:
```go
package optimizer

import (
    "context"
    "sync"
    "time"
    "github.com/jessesanford/idpbuilder/pkg/oci/api"
)

type Executor struct {
    workers int
    pool    chan struct{}
}

func NewExecutor(workers int) *Executor {
    return &Executor{
        workers: workers,
        pool:    make(chan struct{}, workers),
    }
}

func (e *Executor) Execute(stages []api.Stage) error {
    // Implement:
    // 1. Worker pool management
    // 2. Stage scheduling based on dependencies
    // 3. Parallel execution with proper synchronization
    // 4. Result collection and error handling
}

func (e *Executor) executeStage(stage api.Stage, wg *sync.WaitGroup) {
    // Implement stage execution logic
}

func (e *Executor) scheduleStages(stages []api.Stage) [][]api.Stage {
    // Group stages by dependency level for parallel execution
}
```

### 2. pkg/oci/optimizer/graph.go (~120 lines)
**Purpose**: Dependency graph builder and analysis

**Required Implementation**:
```go
package optimizer

import (
=======
# Implementation Plan: Security Layer - Split 001 (Security Manager)

## <� Effort Overview
**Effort ID**: effort4-security-split-001
**Target Size**: 386 lines MAXIMUM
**Purpose**: Security orchestration and policy management
**Order**: IMPLEMENT AFTER split-002 (depends on crypto layer)

## =� CRITICAL REQUIREMENTS
1. **SIZE LIMIT**: 386 lines (well under limit)
2. **DEPENDS ON SPLIT-002**: Use Signer/Verifier from split-002
3. **ORCHESTRATION LAYER**: Coordinate security operations

## =� Files to Implement

### 1. pkg/oci/security/manager.go (386 lines)
**Purpose**: Security orchestration, policy enforcement, and coordination
=======
# Implementation Plan: Security Layer - Split 002 (Crypto Operations)

## <� Effort Overview
**Effort ID**: effort4-security-split-002
**Target Size**: 649 lines MAXIMUM  
**Purpose**: Core cryptographic signing and verification operations
**Order**: MUST BE IMPLEMENTED FIRST (before split-001)

## =� CRITICAL REQUIREMENTS
1. **SIZE LIMIT**: 649 lines (aim for 600 to have buffer)
2. **NO DEPENDENCIES**: This is the foundational layer
3. **COMPLETE INTERFACES**: Manager in split-001 depends on these

## =� Files to Implement

### 1. pkg/oci/security/signer.go (335 lines)
**Purpose**: Digital signature operations for OCI artifacts
>>>>>>> origin/idpbuilder-oci-mgmt/phase2/wave2/effort4-security-split-002

**Core Implementation**:
```go
package security

import (
<<<<<<< HEAD
    "context"
>>>>>>> origin/idpbuilder-oci-mgmt/phase2/wave2/effort4-security-split-001
    "fmt"
    "github.com/jessesanford/idpbuilder/pkg/oci/api"
)

<<<<<<< HEAD
type GraphBuilder struct {
    nodes map[string]*Node
    edges map[string][]string
}

type Node struct {
    Stage    api.Stage
    Level    int
    Visited  bool
    Children []string
}

type DependencyGraph struct {
    Nodes map[string]*Node
    Levels [][]string
}

func NewGraphBuilder() *GraphBuilder {
    return &GraphBuilder{
        nodes: make(map[string]*Node),
        edges: make(map[string][]string),
    }
}

func (g *GraphBuilder) Build(stages []api.Stage) (*DependencyGraph, error) {
    // Implement:
    // 1. Build node map from stages
    // 2. Create edge relationships
    // 3. Perform topological sort
    // 4. Calculate critical path
    // 5. Return structured graph
}

func (g *GraphBuilder) topologicalSort() ([]string, error) {
    // Implement Kahn's algorithm for topological sorting
}

func (g *GraphBuilder) calculateLevels() [][]string {
    // Group nodes by dependency level
}
```

### 3. pkg/oci/optimizer/executor_test.go (~25 lines)
**Purpose**: Basic test stubs

### 4. pkg/oci/optimizer/graph_test.go (~25 lines)
**Purpose**: Basic test stubs

## 🔧 Implementation Steps

### Step 1: Copy API types from split-001
```bash
# Copy the api package from split-001
cp -r ../split-001/pkg/oci/api pkg/oci/
```

### Step 2: Implement executor.go
1. Create worker pool mechanism
2. Implement stage scheduling logic
3. Add parallel execution with sync.WaitGroup
4. Handle errors and timeouts
5. Collect execution results

### Step 3: Implement graph.go
1. Build node structure from stages
2. Create edge relationships from dependencies
3. Implement topological sorting (Kahn's algorithm)
4. Calculate execution levels
5. Identify critical path

### Step 4: Add test stubs
1. Create basic test files
2. Add placeholder test functions
3. Ensure package builds

### Step 5: Verify Size
```bash
/home/vscode/workspaces/idpbuilder-oci-mgmt/tools/line-counter.sh -c branch
# MUST be ≤350 lines
```

## ✅ Success Criteria
- [ ] executor.go implements all required methods (~180 lines)
- [ ] graph.go implements dependency analysis (~120 lines)
- [ ] Test stubs present (~50 lines total)
- [ ] Total implementation ≤350 lines
- [ ] Code compiles with split-001
- [ ] All interfaces satisfied

## 🚨 Critical Notes
1. **DO NOT EXCEED 350 LINES** - Be extremely concise
2. **MUST INTEGRATE** - Use exact types from split-001
3. **FOCUS ON CORE** - Implement minimum viable functionality
4. **NO EXTRAS** - Skip nice-to-haves, focus on essentials

## Integration with Split-001
Split-001 provides:
- `api.Stage`, `api.BuildResult` types
- Stub `Executor` and `GraphBuilder` interfaces
- `Optimizer` that calls these components

Your implementation must satisfy these interfaces exactly.
>>>>>>> origin/idpbuilder-oci-mgmt/phase2/wave2/effort2-optimizer-split-002
=======
type SecurityManager struct {
    signer     api.Signer
    verifier   api.Verifier
    policies   []api.SecurityPolicy
    trustStore *TrustStore
}

// Key methods to implement:
- NewSecurityManager(config *SecurityConfig) (*SecurityManager, error)
- SignArtifact(artifact api.Artifact) (*api.SignedArtifact, error)
- VerifyArtifact(artifact api.SignedArtifact) error
- EnforcePolicy(artifact api.Artifact, policy api.SecurityPolicy) error
- AddPolicy(policy api.SecurityPolicy) error
- RemovePolicy(policyID string) error
- GetTrustChain(keyID string) ([]api.Certificate, error)
- RotateKeys() error
```

### 2. API Imports
**From split-002**: You'll need the crypto interfaces

## =' Implementation Steps

### Step 1: Copy API types from split-002
```bash
# Copy the crypto API from split-002
cp -r ../split-002/pkg/oci/api pkg/oci/
```

### Step 2: Copy manager.go from parent
```bash
cp ../pkg/oci/security/manager.go pkg/oci/security/
```

### Step 3: Import crypto implementations
In manager.go, ensure you're using the interfaces:
```go
import (
    "github.com/jessesanford/idpbuilder/pkg/oci/api"
    // The actual Signer and Verifier will be from split-002
)
```

### Step 4: Implement orchestration
- Use api.Signer for signing operations
- Use api.Verifier for verification
- Add policy enforcement layer
- Implement key rotation logic
- Add trust chain management

### Step 5: Verify compilation
=======
    "crypto"
    "crypto/rsa"
    "crypto/sha256"
    "crypto/x509"
    "encoding/pem"
    "github.com/jessesanford/idpbuilder/pkg/oci/api"
)

type Signer struct {
    privateKey crypto.PrivateKey
    publicKey  crypto.PublicKey
    algorithm  api.SignatureAlgorithm
}

// Key methods to implement:
- NewSigner(privateKeyPEM []byte) (*Signer, error)
- Sign(data []byte) ([]byte, error)
- SignManifest(manifest api.Manifest) (*api.SignedManifest, error)
- GetPublicKey() ([]byte, error)
- VerifyOwnSignature(data, signature []byte) error
```

### 2. pkg/oci/security/verifier.go (314 lines)
**Purpose**: Signature verification and trust validation

**Core Implementation**:
```go
package security

import (
    "crypto"
    "crypto/rsa"
    "crypto/sha256"
    "crypto/x509"
    "github.com/jessesanford/idpbuilder/pkg/oci/api"
)

type Verifier struct {
    trustedKeys map[string]crypto.PublicKey
    trustStore  *TrustStore
}

// Key methods to implement:
- NewVerifier(trustedKeys [][]byte) (*Verifier, error)
- Verify(data, signature []byte, keyID string) error
- VerifyManifest(manifest api.SignedManifest) error
- AddTrustedKey(keyID string, publicKey []byte) error
- RemoveTrustedKey(keyID string) error
```

### 3. pkg/oci/security/trust_store.go (embedded in verifier.go)
**Purpose**: Manage trusted keys and certificates
**Note**: Keep minimal, embed in verifier.go to save lines

## =' Implementation Steps

### Step 1: Copy existing code from parent directory
```bash
cp ../pkg/oci/security/signer.go pkg/oci/security/
cp ../pkg/oci/security/verifier.go pkg/oci/security/
```

### Step 2: Ensure API types are available
```bash
# Check if api types exist, if not copy from effort1-contracts
cp -r /home/vscode/workspaces/idpbuilder-oci-mgmt/efforts/phase2/wave2/effort1-contracts/pkg/oci/api pkg/oci/
```

### Step 3: Optimize if needed
- If total exceeds 649 lines, optimize:
  - Combine helper functions
  - Remove verbose comments
  - Simplify error handling

### Step 4: Verify compilation
>>>>>>> origin/idpbuilder-oci-mgmt/phase2/wave2/effort4-security-split-002
```bash
cd pkg/oci/security
go build .
```

<<<<<<< HEAD
### Step 6: Measure size
```bash
/home/vscode/workspaces/idpbuilder-oci-mgmt/tools/line-counter.sh -c branch
# MUST be d386 lines
```

##  Success Criteria
- [ ] manager.go implements security orchestration (d386 lines)
- [ ] Uses Signer/Verifier interfaces from split-002
- [ ] Implements policy enforcement
- [ ] Handles key rotation
- [ ] Code compiles successfully
- [ ] Total d386 lines

## =� Critical Notes
1. **DEPENDS ON SPLIT-002**: Must use the crypto interfaces
2. **ORCHESTRATION FOCUS**: Don't re-implement crypto
3. **POLICY LAYER**: Add value on top of basic crypto
4. **SIZE COMFORTABLE**: 386 lines gives plenty of room

## Integration Points
- Uses api.Signer from split-002 for signing
- Uses api.Verifier from split-002 for verification  
- Adds SecurityPolicy enforcement on top
- Provides unified SecurityManager interface
>>>>>>> origin/idpbuilder-oci-mgmt/phase2/wave2/effort4-security-split-001
=======
### Step 5: Measure size
```bash
/home/vscode/workspaces/idpbuilder-oci-mgmt/tools/line-counter.sh -c branch
# MUST be d649 lines
```

##  Success Criteria
- [ ] signer.go implements all signing operations (~335 lines)
- [ ] verifier.go implements all verification (~314 lines)
- [ ] Total d649 lines
- [ ] Code compiles independently
- [ ] No dependency on manager.go
- [ ] Interfaces ready for split-001

## =� Critical Notes
1. **FOUNDATIONAL LAYER**: Split-001 depends on this
2. **PRESERVE INTERFACES**: Manager expects specific method signatures
3. **NO MANAGER REFERENCES**: This must be independent
4. **SECURITY CRITICAL**: Ensure crypto operations are correct

## Dependencies for Split-001
Split-001 (manager.go) will:
- Import Signer and Verifier types
- Use these for security orchestration
- Add policy enforcement on top
>>>>>>> origin/idpbuilder-oci-mgmt/phase2/wave2/effort4-security-split-002
=======
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
>>>>>>> origin/idpbuidler-oci-mgmt/phase2/wave2/effort5-registry
