<<<<<<< HEAD
# Registry Authentication Implementation Plan

## 📌 Effort Overview

**Effort ID**: E1.1.2B
**Effort Name**: registry-auth
**Phase**: 1, Wave: 1
**Branch**: `idpbuilder-oci-build-push/phase1/wave1/registry-auth`
**Base Branch**: `idpbuilder-oci-build-push/phase1/wave1/registry-types`
**Estimated Lines**: 350
**Can Parallelize**: No (depends on registry-types)
**Parallel With**: None
**Dependencies**: [registry-types (E1.1.2A)]
**Created**: 2025-09-18

## 🎯 Mission

Implement authentication handlers and middleware for OCI registry operations. This effort provides the authentication layer that sits between the registry types and actual registry operations, handling credential management, token negotiation, and authentication flow.

## 📦 Scope and Boundaries

### What This Effort Includes
- Authentication handler implementations
- Token management and refresh logic
- Credential validation
- Authentication middleware for registry operations
- Basic auth and token auth implementations

### What This Effort Excludes
- Core registry types (in registry-types)
- Helper utilities and convenience functions (in registry-helpers)
- Test implementations (in registry-tests)
- Actual registry client operations
=======
# E1.2.2 - Fallback Strategies Implementation Plan

## 🚨 CRITICAL EFFORT METADATA (FROM WAVE PLAN)
**Effort ID**: E1.2.2  
**Effort Name**: Fallback Strategies  
**Branch**: `phase1/wave2/fallback-strategies`  
**Can Parallelize**: Yes  
**Parallel With**: [E1.2.1 - Certificate Validation]  
**Size Estimate**: 700 lines  
**Dependencies**: [E1.1.1 - Kind Certificate Extraction, E1.1.2 - Registry TLS Trust]  
**Feature Flag**: `FALLBACK_STRATEGIES_ENABLED`  

## Overview
- **Effort**: Implement fallback handling mechanisms and --insecure flag support for certificate validation
- **Phase**: 1 (Certificate Infrastructure), Wave: 2 (Certificate Validation & Fallback)
- **Estimated Size**: 700 lines (under 800 line limit)
- **Implementation Time**: 6-8 hours

## 🎯 Mission Statement
Establish robust fallback mechanisms for certificate validation failures, implement the --insecure flag functionality to bypass certificate checks when explicitly requested, and provide graceful degradation when certificates cannot be validated. This ensures the system remains usable in development environments while maintaining security in production.

## 📋 Technical Architecture

### Core Components
1. **Fallback Strategy Manager**: Orchestrates fallback mechanisms based on configuration
2. **Insecure Mode Handler**: Manages --insecure flag state and behavior
3. **Retry Logic**: Implements exponential backoff for transient failures
4. **Graceful Degradation**: Provides progressive fallback options
5. **Warning System**: Clear user notifications about security implications

### Integration Points
- **Wave 1 E1.1.1**: Uses KindCertExtractor for certificate retrieval attempts
- **Wave 1 E1.1.2**: Extends DefaultTrustStore with fallback capabilities
- **Wave 2 E1.2.1**: Coordinates with validation logic for failure handling
>>>>>>> origin/idpbuilder-oci-build-push/phase2/wave1/gitea-client-split-001

## 📁 File Structure

```
<<<<<<< HEAD
pkg/
└── registry/
    └── auth/
        ├── authenticator.go      (~80 lines) - Core authenticator interface and factory
        ├── basic.go              (~60 lines) - Basic auth implementation
        ├── token.go              (~90 lines) - Token auth implementation
        ├── middleware.go         (~70 lines) - Auth middleware for HTTP clients
        └── manager.go            (~50 lines) - Auth manager for credential lifecycle
```

## 🔧 Implementation Details

### 1. Core Authenticator Interface (`pkg/registry/auth/authenticator.go`)

```go
package auth

import (
    "context"
    "net/http"
    "github.com/cnoe-io/idpbuilder/pkg/registry/types"
)

// Authenticator defines the interface for registry authentication
type Authenticator interface {
    // Authenticate adds authentication to the request
    Authenticate(ctx context.Context, req *http.Request) error

    // Refresh refreshes authentication credentials if needed
    Refresh(ctx context.Context) error

    // IsValid checks if current auth is still valid
    IsValid() bool
}

// NewAuthenticator creates an authenticator based on config
func NewAuthenticator(config *types.AuthConfig) (Authenticator, error) {
    switch config.AuthType {
    case types.AuthTypeBasic:
        return NewBasicAuthenticator(config)
    case types.AuthTypeToken:
        return NewTokenAuthenticator(config)
    case types.AuthTypeNone:
        return NewNoOpAuthenticator(), nil
    default:
        return nil, fmt.Errorf("unsupported auth type: %s", config.AuthType)
    }
}

// NoOpAuthenticator for registries without auth
type NoOpAuthenticator struct{}

func NewNoOpAuthenticator() *NoOpAuthenticator {
    return &NoOpAuthenticator{}
}

func (n *NoOpAuthenticator) Authenticate(ctx context.Context, req *http.Request) error {
    return nil
}

func (n *NoOpAuthenticator) Refresh(ctx context.Context) error {
    return nil
}

func (n *NoOpAuthenticator) IsValid() bool {
    return true
}
```

### 2. Basic Authentication (`pkg/registry/auth/basic.go`)

```go
package auth

import (
    "context"
    "encoding/base64"
    "fmt"
    "net/http"
    "github.com/cnoe-io/idpbuilder/pkg/registry/types"
)

// BasicAuthenticator implements basic authentication
type BasicAuthenticator struct {
    username string
    password string
    encoded  string
}

// NewBasicAuthenticator creates a new basic authenticator
func NewBasicAuthenticator(config *types.AuthConfig) (*BasicAuthenticator, error) {
    if config.Username == "" || config.Password == "" {
        return nil, fmt.Errorf("username and password required for basic auth")
    }

    auth := &BasicAuthenticator{
        username: config.Username,
        password: config.Password,
    }
    auth.encoded = base64.StdEncoding.EncodeToString(
        []byte(fmt.Sprintf("%s:%s", config.Username, config.Password)),
    )

    return auth, nil
}

func (b *BasicAuthenticator) Authenticate(ctx context.Context, req *http.Request) error {
    req.Header.Set("Authorization", "Basic "+b.encoded)
    return nil
}

func (b *BasicAuthenticator) Refresh(ctx context.Context) error {
    // Basic auth doesn't need refresh
    return nil
}

func (b *BasicAuthenticator) IsValid() bool {
    return b.encoded != ""
}
```

### 3. Token Authentication (`pkg/registry/auth/token.go`)

```go
package auth

import (
    "context"
    "fmt"
    "net/http"
    "sync"
    "time"
    "github.com/cnoe-io/idpbuilder/pkg/registry/types"
)

// TokenAuthenticator implements bearer token authentication
type TokenAuthenticator struct {
    mu          sync.RWMutex
    token       string
    expiresAt   time.Time
    authConfig  *types.AuthConfig
    tokenClient TokenClient
}

// TokenClient interface for token operations
type TokenClient interface {
    RequestToken(ctx context.Context, config *types.AuthConfig) (*types.TokenResponse, error)
}

// NewTokenAuthenticator creates a new token authenticator
func NewTokenAuthenticator(config *types.AuthConfig, client TokenClient) (*TokenAuthenticator, error) {
    if config.Token == "" && client == nil {
        return nil, fmt.Errorf("token or token client required")
    }

    auth := &TokenAuthenticator{
        authConfig:  config,
        tokenClient: client,
    }

    if config.Token != "" {
        auth.token = config.Token
        // Set a default expiry if token is provided directly
        auth.expiresAt = time.Now().Add(1 * time.Hour)
    }

    return auth, nil
}

func (t *TokenAuthenticator) Authenticate(ctx context.Context, req *http.Request) error {
    t.mu.RLock()
    token := t.token
    t.mu.RUnlock()

    if token == "" {
        if err := t.Refresh(ctx); err != nil {
            return fmt.Errorf("failed to obtain token: %w", err)
        }
        t.mu.RLock()
        token = t.token
        t.mu.RUnlock()
    }

    req.Header.Set("Authorization", "Bearer "+token)
    return nil
}

func (t *TokenAuthenticator) Refresh(ctx context.Context) error {
    if t.tokenClient == nil {
        return fmt.Errorf("no token client configured for refresh")
    }

    resp, err := t.tokenClient.RequestToken(ctx, t.authConfig)
    if err != nil {
        return fmt.Errorf("token request failed: %w", err)
    }

    t.mu.Lock()
    t.token = resp.Token
    if resp.ExpiresIn > 0 {
        t.expiresAt = time.Now().Add(time.Duration(resp.ExpiresIn) * time.Second)
    } else {
        t.expiresAt = time.Now().Add(1 * time.Hour)
    }
    t.mu.Unlock()

    return nil
}

func (t *TokenAuthenticator) IsValid() bool {
    t.mu.RLock()
    defer t.mu.RUnlock()

    if t.token == "" {
        return false
    }

    // Check if token is expired with 30-second buffer
    return time.Now().Add(30 * time.Second).Before(t.expiresAt)
}
```

### 4. Authentication Middleware (`pkg/registry/auth/middleware.go`)

```go
package auth

import (
    "context"
    "net/http"
    "github.com/cnoe-io/idpbuilder/pkg/registry/types"
)

// Transport wraps an http.RoundTripper with authentication
type Transport struct {
    Base          http.RoundTripper
    Authenticator Authenticator
}

// NewTransport creates a new authenticated transport
func NewTransport(base http.RoundTripper, auth Authenticator) *Transport {
    if base == nil {
        base = http.DefaultTransport
    }

    return &Transport{
        Base:          base,
        Authenticator: auth,
    }
}

// RoundTrip implements http.RoundTripper
func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
    // Clone the request to avoid modifying the original
    clonedReq := req.Clone(req.Context())

    // Apply authentication
    if t.Authenticator != nil {
        // Check if auth needs refresh
        if !t.Authenticator.IsValid() {
            if err := t.Authenticator.Refresh(req.Context()); err != nil {
                return nil, fmt.Errorf("auth refresh failed: %w", err)
            }
        }

        if err := t.Authenticator.Authenticate(req.Context(), clonedReq); err != nil {
            return nil, fmt.Errorf("authentication failed: %w", err)
        }
    }

    // Execute the request
    resp, err := t.Base.RoundTrip(clonedReq)
    if err != nil {
        return nil, err
    }

    // Handle 401 Unauthorized by refreshing auth and retrying once
    if resp.StatusCode == http.StatusUnauthorized && t.Authenticator != nil {
        resp.Body.Close()

        if err := t.Authenticator.Refresh(req.Context()); err != nil {
            return nil, fmt.Errorf("auth refresh after 401 failed: %w", err)
        }

        // Retry with refreshed auth
        retryReq := req.Clone(req.Context())
        if err := t.Authenticator.Authenticate(req.Context(), retryReq); err != nil {
            return nil, fmt.Errorf("re-authentication failed: %w", err)
        }

        return t.Base.RoundTrip(retryReq)
    }

    return resp, nil
}
```

### 5. Auth Manager (`pkg/registry/auth/manager.go`)

```go
package auth

import (
    "context"
    "fmt"
    "sync"
    "github.com/cnoe-io/idpbuilder/pkg/registry/types"
)

// Manager manages authentication for multiple registries
type Manager struct {
    mu      sync.RWMutex
    auths   map[string]Authenticator
    store   types.CredentialStore
}

// NewManager creates a new auth manager
func NewManager(store types.CredentialStore) *Manager {
    return &Manager{
        auths: make(map[string]Authenticator),
        store: store,
    }
}

// GetAuthenticator gets or creates an authenticator for a registry
func (m *Manager) GetAuthenticator(ctx context.Context, registry string) (Authenticator, error) {
    m.mu.RLock()
    auth, exists := m.auths[registry]
    m.mu.RUnlock()

    if exists && auth.IsValid() {
        return auth, nil
    }

    // Load credentials from store
    creds, err := m.store.GetCredentials(registry)
    if err != nil {
        return nil, fmt.Errorf("failed to get credentials: %w", err)
    }

    // Create new authenticator
    auth, err = NewAuthenticator(creds)
    if err != nil {
        return nil, fmt.Errorf("failed to create authenticator: %w", err)
    }

    // Cache the authenticator
    m.mu.Lock()
    m.auths[registry] = auth
    m.mu.Unlock()

    return auth, nil
}

// Clear removes cached authenticator for a registry
func (m *Manager) Clear(registry string) {
    m.mu.Lock()
    delete(m.auths, registry)
    m.mu.Unlock()
}
```

## 📊 Size Estimation Breakdown

| File | Estimated Lines | Purpose |
|------|-----------------|---------|
| authenticator.go | 80 | Core interface and factory |
| basic.go | 60 | Basic auth implementation |
| token.go | 90 | Token auth with refresh |
| middleware.go | 70 | HTTP transport wrapper |
| manager.go | 50 | Multi-registry auth management |
| **TOTAL** | **350** | Within limit |

## 🔗 Dependencies and Integration

### Internal Dependencies
- `github.com/cnoe-io/idpbuilder/pkg/registry/types` - From registry-types effort (E1.1.2A)
  - Uses: `AuthConfig`, `AuthType`, `TokenResponse`, `CredentialStore`

### External Dependencies
- Standard library only (`net/http`, `context`, `sync`, `time`, `encoding/base64`)

### Integration Points
1. **Registry Types**: Import and use types from registry-types package
2. **Registry Client**: Will be used by registry client (future effort) via middleware
3. **Registry Helpers**: Helpers will use these authenticators for operations

## ⚡ Implementation Strategy

### Phase 1: Core Structure (50 lines)
1. Create package structure
2. Define Authenticator interface
3. Implement NoOpAuthenticator

### Phase 2: Basic Auth (60 lines)
1. Implement BasicAuthenticator
2. Add header generation
3. Test with mock requests

### Phase 3: Token Auth (90 lines)
1. Define TokenClient interface
2. Implement TokenAuthenticator
3. Add refresh logic with expiry

### Phase 4: Middleware (70 lines)
1. Create Transport wrapper
2. Add authentication injection
3. Implement 401 retry logic

### Phase 5: Manager (50 lines)
1. Implement auth caching
2. Add credential store integration
3. Handle multi-registry scenarios

### Phase 6: Integration (30 lines)
1. Wire components together
2. Add error handling
3. Final testing and validation
=======
efforts/phase1/wave2/fallback-strategies/
├── pkg/
│   ├── fallback/
│   │   ├── manager.go           # Core fallback strategy manager (200 lines)
│   │   ├── manager_test.go      # Unit tests for manager (150 lines)
│   │   ├── strategies.go        # Fallback strategy implementations (150 lines)
│   │   └── strategies_test.go   # Strategy tests (100 lines)
│   └── insecure/
│       ├── handler.go           # Insecure mode implementation (50 lines)
│       └── handler_test.go      # Insecure mode tests (50 lines)
└── IMPLEMENTATION-PLAN-20250907-064500.md
```

## 🔧 Implementation Steps

### Step 1: Create Fallback Manager Core (200 lines)
**File**: `pkg/fallback/manager.go`

```go
package fallback

import (
    "context"
    "crypto/x509"
    "fmt"
    "time"
    
    "github.com/jessesanford/idpbuilder/pkg/certs"
)

// FallbackStrategy defines the interface for fallback mechanisms
type FallbackStrategy interface {
    Name() string
    Priority() int
    Execute(ctx context.Context, registry string) error
    ShouldRetry(err error) bool
}

// FallbackManager coordinates fallback strategies
type FallbackManager struct {
    strategies      []FallbackStrategy
    trustStore      certs.TrustStoreManager
    insecureMode    bool
    maxRetries      int
    retryDelay      time.Duration
    warningCallback func(string)
}

// NewFallbackManager creates a new fallback manager
func NewFallbackManager(trustStore certs.TrustStoreManager, opts ...Option) *FallbackManager {
    fm := &FallbackManager{
        trustStore:  trustStore,
        strategies:  make([]FallbackStrategy, 0),
        maxRetries:  3,
        retryDelay:  time.Second,
        warningCallback: defaultWarning,
    }
    
    // Apply options
    for _, opt := range opts {
        opt(fm)
    }
    
    // Initialize default strategies
    fm.initDefaultStrategies()
    return fm
}

// Option configures the FallbackManager
type Option func(*FallbackManager)

// WithInsecureMode enables insecure mode
func WithInsecureMode(insecure bool) Option {
    return func(fm *FallbackManager) {
        fm.insecureMode = insecure
    }
}

// WithMaxRetries sets maximum retry attempts
func WithMaxRetries(max int) Option {
    return func(fm *FallbackManager) {
        fm.maxRetries = max
    }
}

// HandleValidationFailure processes certificate validation failures
func (fm *FallbackManager) HandleValidationFailure(ctx context.Context, registry string, err error) error {
    // Check if insecure mode is enabled
    if fm.insecureMode {
        fm.warningCallback(fmt.Sprintf("⚠️  INSECURE MODE: Bypassing certificate validation for %s", registry))
        return fm.trustStore.SetInsecure(registry, true)
    }
    
    // Try fallback strategies in order of priority
    for _, strategy := range fm.strategies {
        select {
        case <-ctx.Done():
            return ctx.Err()
        default:
            if err := fm.executeWithRetry(ctx, strategy, registry); err == nil {
                return nil
            }
        }
    }
    
    return fmt.Errorf("all fallback strategies failed for %s: %w", registry, err)
}

// executeWithRetry executes a strategy with retry logic
func (fm *FallbackManager) executeWithRetry(ctx context.Context, strategy FallbackStrategy, registry string) error {
    var lastErr error
    
    for attempt := 0; attempt < fm.maxRetries; attempt++ {
        if attempt > 0 {
            // Exponential backoff
            delay := fm.retryDelay * time.Duration(1<<uint(attempt-1))
            select {
            case <-time.After(delay):
            case <-ctx.Done():
                return ctx.Err()
            }
        }
        
        lastErr = strategy.Execute(ctx, registry)
        if lastErr == nil {
            return nil
        }
        
        if !strategy.ShouldRetry(lastErr) {
            break
        }
    }
    
    return lastErr
}

// initDefaultStrategies sets up the default fallback strategies
func (fm *FallbackManager) initDefaultStrategies() {
    fm.strategies = []FallbackStrategy{
        NewSystemCertStrategy(fm.trustStore),
        NewCachedCertStrategy(fm.trustStore),
        NewSelfSignedAcceptStrategy(fm.trustStore),
    }
    
    // Sort by priority
    sortStrategies(fm.strategies)
}
```

### Step 2: Implement Fallback Strategies (150 lines)
**File**: `pkg/fallback/strategies.go`

```go
package fallback

import (
    "context"
    "crypto/x509"
    "fmt"
    "os"
    "path/filepath"
    
    "github.com/jessesanford/idpbuilder/pkg/certs"
)

// SystemCertStrategy tries to use system certificate store
type SystemCertStrategy struct {
    trustStore certs.TrustStoreManager
    priority   int
}

func NewSystemCertStrategy(ts certs.TrustStoreManager) *SystemCertStrategy {
    return &SystemCertStrategy{
        trustStore: ts,
        priority:   1,
    }
}

func (s *SystemCertStrategy) Name() string { return "system-cert-fallback" }
func (s *SystemCertStrategy) Priority() int { return s.priority }

func (s *SystemCertStrategy) Execute(ctx context.Context, registry string) error {
    // Try to load certificates from system store
    systemPool, err := x509.SystemCertPool()
    if err != nil {
        return fmt.Errorf("failed to load system cert pool: %w", err)
    }
    
    // Add system certs to trust store for this registry
    for _, cert := range systemPool.Subjects() {
        // Parse and add certificate
        // Implementation details...
    }
    
    return nil
}

func (s *SystemCertStrategy) ShouldRetry(err error) bool {
    // Don't retry system cert loading failures
    return false
}

// CachedCertStrategy uses previously cached certificates
type CachedCertStrategy struct {
    trustStore certs.TrustStoreManager
    cacheDir   string
    priority   int
}

func NewCachedCertStrategy(ts certs.TrustStoreManager) *CachedCertStrategy {
    homeDir, _ := os.UserHomeDir()
    return &CachedCertStrategy{
        trustStore: ts,
        cacheDir:   filepath.Join(homeDir, ".idpbuilder", "cert-cache"),
        priority:   2,
    }
}

func (c *CachedCertStrategy) Name() string { return "cached-cert-fallback" }
func (c *CachedCertStrategy) Priority() int { return c.priority }

func (c *CachedCertStrategy) Execute(ctx context.Context, registry string) error {
    // Look for cached certificates for this registry
    cacheFile := filepath.Join(c.cacheDir, fmt.Sprintf("%s.pem", sanitizeFilename(registry)))
    
    if _, err := os.Stat(cacheFile); os.IsNotExist(err) {
        return fmt.Errorf("no cached certificate found for %s", registry)
    }
    
    // Load and validate cached certificate
    certData, err := os.ReadFile(cacheFile)
    if err != nil {
        return fmt.Errorf("failed to read cached cert: %w", err)
    }
    
    // Parse and add to trust store
    // Implementation details...
    
    return nil
}

func (c *CachedCertStrategy) ShouldRetry(err error) bool {
    // Retry on transient file system errors
    return os.IsTimeout(err)
}

// SelfSignedAcceptStrategy accepts self-signed certificates with user warning
type SelfSignedAcceptStrategy struct {
    trustStore certs.TrustStoreManager
    priority   int
}

func NewSelfSignedAcceptStrategy(ts certs.TrustStoreManager) *SelfSignedAcceptStrategy {
    return &SelfSignedAcceptStrategy{
        trustStore: ts,
        priority:   10, // Lowest priority
    }
}

func (s *SelfSignedAcceptStrategy) Name() string { return "self-signed-accept" }
func (s *SelfSignedAcceptStrategy) Priority() int { return s.priority }

func (s *SelfSignedAcceptStrategy) Execute(ctx context.Context, registry string) error {
    // Warn user about accepting self-signed certificate
    fmt.Printf("⚠️  WARNING: Accepting self-signed certificate for %s\n", registry)
    fmt.Println("This reduces security. Use --insecure flag to suppress this warning.")
    
    // Configure trust store to accept self-signed for this registry
    return s.trustStore.SetInsecure(registry, true)
}

func (s *SelfSignedAcceptStrategy) ShouldRetry(err error) bool {
    return false
}

// Helper functions
func sortStrategies(strategies []FallbackStrategy) {
    // Sort strategies by priority (lower number = higher priority)
    // Implementation...
}

func sanitizeFilename(s string) string {
    // Sanitize registry name for filesystem
    // Implementation...
}
```

### Step 3: Implement Insecure Mode Handler (50 lines)
**File**: `pkg/insecure/handler.go`

```go
package insecure

import (
    "fmt"
    "os"
    "strings"
)

// InsecureHandler manages the --insecure flag behavior
type InsecureHandler struct {
    enabled     bool
    registries  map[string]bool
    warnOnce    map[string]bool
}

// NewInsecureHandler creates a new insecure mode handler
func NewInsecureHandler() *InsecureHandler {
    return &InsecureHandler{
        enabled:    false,
        registries: make(map[string]bool),
        warnOnce:   make(map[string]bool),
    }
}

// Enable activates insecure mode
func (h *InsecureHandler) Enable(registries ...string) {
    h.enabled = true
    
    if len(registries) == 0 {
        // Global insecure mode
        h.WarnGlobal()
    } else {
        // Registry-specific insecure mode
        for _, reg := range registries {
            h.registries[reg] = true
            h.WarnRegistry(reg)
        }
    }
}

// IsInsecure checks if insecure mode is enabled for a registry
func (h *InsecureHandler) IsInsecure(registry string) bool {
    if h.enabled && len(h.registries) == 0 {
        // Global insecure mode
        return true
    }
    return h.registries[registry]
}

// WarnGlobal displays a warning for global insecure mode
func (h *InsecureHandler) WarnGlobal() {
    if !h.warnOnce["_global"] {
        fmt.Fprintln(os.Stderr, strings.Repeat("⚠", 10))
        fmt.Fprintln(os.Stderr, "WARNING: Running in INSECURE mode")
        fmt.Fprintln(os.Stderr, "Certificate validation is DISABLED for ALL registries")
        fmt.Fprintln(os.Stderr, "This should ONLY be used in development environments")
        fmt.Fprintln(os.Stderr, strings.Repeat("⚠", 10))
        h.warnOnce["_global"] = true
    }
}

// WarnRegistry displays a warning for registry-specific insecure mode
func (h *InsecureHandler) WarnRegistry(registry string) {
    if !h.warnOnce[registry] {
        fmt.Fprintf(os.Stderr, "⚠️  WARNING: Certificate validation disabled for %s\n", registry)
        h.warnOnce[registry] = true
    }
}
```

### Step 4: Create Unit Tests for Manager (150 lines)
**File**: `pkg/fallback/manager_test.go`

```go
package fallback

import (
    "context"
    "errors"
    "testing"
    "time"
    
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

// MockTrustStore for testing
type MockTrustStore struct {
    mock.Mock
}

func (m *MockTrustStore) SetInsecure(registry string, insecure bool) error {
    args := m.Called(registry, insecure)
    return args.Error(0)
}

// Additional mock methods...

func TestFallbackManager_InsecureMode(t *testing.T) {
    tests := []struct {
        name         string
        insecureMode bool
        expectBypass bool
    }{
        {
            name:         "insecure mode enabled",
            insecureMode: true,
            expectBypass: true,
        },
        {
            name:         "insecure mode disabled",
            insecureMode: false,
            expectBypass: false,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            mockStore := new(MockTrustStore)
            fm := NewFallbackManager(mockStore, WithInsecureMode(tt.insecureMode))
            
            if tt.expectBypass {
                mockStore.On("SetInsecure", "test.registry", true).Return(nil)
            }
            
            ctx := context.Background()
            err := fm.HandleValidationFailure(ctx, "test.registry", errors.New("cert error"))
            
            if tt.expectBypass {
                assert.NoError(t, err)
                mockStore.AssertExpectations(t)
            } else {
                assert.Error(t, err)
            }
        })
    }
}

func TestFallbackManager_RetryLogic(t *testing.T) {
    mockStore := new(MockTrustStore)
    fm := NewFallbackManager(mockStore, 
        WithMaxRetries(3),
        WithRetryDelay(10*time.Millisecond))
    
    // Create a mock strategy that fails twice then succeeds
    mockStrategy := &MockStrategy{
        failCount: 2,
    }
    fm.strategies = []FallbackStrategy{mockStrategy}
    
    ctx := context.Background()
    err := fm.HandleValidationFailure(ctx, "test.registry", errors.New("initial error"))
    
    assert.NoError(t, err)
    assert.Equal(t, 3, mockStrategy.attempts)
}

func TestFallbackManager_ContextCancellation(t *testing.T) {
    mockStore := new(MockTrustStore)
    fm := NewFallbackManager(mockStore)
    
    // Create a context that's already cancelled
    ctx, cancel := context.WithCancel(context.Background())
    cancel()
    
    err := fm.HandleValidationFailure(ctx, "test.registry", errors.New("cert error"))
    
    assert.Equal(t, context.Canceled, err)
}

// Additional test cases...
```

### Step 5: Create Strategy Tests (100 lines)
**File**: `pkg/fallback/strategies_test.go`

```go
package fallback

import (
    "context"
    "os"
    "path/filepath"
    "testing"
    
    "github.com/stretchr/testify/assert"
)

func TestSystemCertStrategy(t *testing.T) {
    mockStore := new(MockTrustStore)
    strategy := NewSystemCertStrategy(mockStore)
    
    assert.Equal(t, "system-cert-fallback", strategy.Name())
    assert.Equal(t, 1, strategy.Priority())
    
    // Test execution
    ctx := context.Background()
    err := strategy.Execute(ctx, "test.registry")
    
    // System cert loading may fail in test environment
    // Just verify it attempts the operation
    assert.NotNil(t, err)
}

func TestCachedCertStrategy(t *testing.T) {
    // Create temp directory for cache
    tmpDir := t.TempDir()
    
    mockStore := new(MockTrustStore)
    strategy := &CachedCertStrategy{
        trustStore: mockStore,
        cacheDir:   tmpDir,
        priority:   2,
    }
    
    // Test with no cached cert
    ctx := context.Background()
    err := strategy.Execute(ctx, "test.registry")
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "no cached certificate found")
    
    // Create a cached cert file
    cacheFile := filepath.Join(tmpDir, "test.registry.pem")
    testCert := []byte("-----BEGIN CERTIFICATE-----\ntest\n-----END CERTIFICATE-----")
    assert.NoError(t, os.WriteFile(cacheFile, testCert, 0644))
    
    // Test with cached cert
    err = strategy.Execute(ctx, "test.registry")
    // Will fail on parse but shows it reads the file
    assert.Error(t, err)
}

func TestSelfSignedAcceptStrategy(t *testing.T) {
    mockStore := new(MockTrustStore)
    mockStore.On("SetInsecure", "test.registry", true).Return(nil)
    
    strategy := NewSelfSignedAcceptStrategy(mockStore)
    
    assert.Equal(t, "self-signed-accept", strategy.Name())
    assert.Equal(t, 10, strategy.Priority()) // Lowest priority
    
    ctx := context.Background()
    err := strategy.Execute(ctx, "test.registry")
    
    assert.NoError(t, err)
    mockStore.AssertExpectations(t)
}

// Additional test cases...
```

### Step 6: Create Insecure Handler Tests (50 lines)
**File**: `pkg/insecure/handler_test.go`

```go
package insecure

import (
    "testing"
    
    "github.com/stretchr/testify/assert"
)

func TestInsecureHandler_GlobalMode(t *testing.T) {
    handler := NewInsecureHandler()
    
    // Initially disabled
    assert.False(t, handler.IsInsecure("any.registry"))
    
    // Enable global insecure mode
    handler.Enable()
    
    // Should be insecure for all registries
    assert.True(t, handler.IsInsecure("registry1.example.com"))
    assert.True(t, handler.IsInsecure("registry2.example.com"))
}

func TestInsecureHandler_RegistrySpecific(t *testing.T) {
    handler := NewInsecureHandler()
    
    // Enable for specific registries
    handler.Enable("registry1.example.com", "registry2.example.com")
    
    // Should be insecure only for specified registries
    assert.True(t, handler.IsInsecure("registry1.example.com"))
    assert.True(t, handler.IsInsecure("registry2.example.com"))
    assert.False(t, handler.IsInsecure("registry3.example.com"))
}

func TestInsecureHandler_WarnOnce(t *testing.T) {
    handler := NewInsecureHandler()
    
    // First warning should set the flag
    handler.WarnRegistry("test.registry")
    assert.True(t, handler.warnOnce["test.registry"])
    
    // Subsequent calls should not change state
    handler.WarnRegistry("test.registry")
    assert.True(t, handler.warnOnce["test.registry"])
}
```

## 📊 Size Management Strategy

### Line Count Breakdown
- `pkg/fallback/manager.go`: 200 lines
- `pkg/fallback/strategies.go`: 150 lines
- `pkg/fallback/manager_test.go`: 150 lines
- `pkg/fallback/strategies_test.go`: 100 lines
- `pkg/insecure/handler.go`: 50 lines
- `pkg/insecure/handler_test.go`: 50 lines
- **Total**: 700 lines (under 800 limit)

### Size Control Measures
1. **Regular Measurement**: Use `${PROJECT_ROOT}/tools/line-counter.sh` after each file
2. **Checkpoint at 500 lines**: Verify trajectory
3. **Warning at 650 lines**: Consider optimization
4. **Stop at 700 lines**: Complete testing and documentation

### Measurement Commands
```bash
# Navigate to effort directory
cd /home/vscode/workspaces/idpbuilder-oci-build-push/efforts/phase1/wave2/fallback-strategies

# Find project root
PROJECT_ROOT=$(pwd)
while [ "$PROJECT_ROOT" != "/" ]; do 
    [ -f "$PROJECT_ROOT/orchestrator-state.yaml" ] && break
    PROJECT_ROOT=$(dirname "$PROJECT_ROOT")
done

# Measure lines (after implementation)
$PROJECT_ROOT/tools/line-counter.sh
```

## 🧪 Testing Requirements
>>>>>>> origin/idpbuilder-oci-build-push/phase2/wave1/gitea-client-split-001

### Unit Test Coverage
- **Target**: 85% coverage minimum
- **Critical Paths**: 100% coverage required
  - Insecure mode activation
  - Fallback strategy execution
  - Retry logic
  - Context cancellation

<<<<<<< HEAD
Testing will be handled in the separate registry-tests effort (E1.1.2D), but key test scenarios include:

1. **Unit Tests**:
   - Each authenticator type
   - Token refresh logic
   - Middleware behavior
   - Manager caching

2. **Integration Tests**:
   - Full auth flow with mock registry
   - Credential store integration
   - Multi-registry scenarios

## ✅ Success Criteria

1. **Functional Requirements**:
   - ✅ Basic authentication works
   - ✅ Token authentication with refresh
   - ✅ Middleware properly injects auth
   - ✅ Manager handles multiple registries

2. **Non-Functional Requirements**:
   - ✅ Thread-safe operations
   - ✅ Proper error handling
   - ✅ Clean separation from other efforts
   - ✅ Under 350 lines limit

## 🚀 Next Steps

After this effort is complete:
1. **registry-helpers (E1.1.2C)** will build convenience functions using these authenticators
2. **registry-tests (E1.1.2D)** will provide comprehensive test coverage
3. Future efforts will use this auth layer for actual registry operations

## 📝 Notes

- This effort focuses ONLY on authentication logic
- No actual registry operations are implemented here
- Clean interfaces allow for easy extension (OAuth2, etc.)
- Thread safety is critical for concurrent operations
- Error messages should be clear and actionable
=======
### Integration Test Scenarios
1. **Insecure Flag**: Verify --insecure bypasses all validation
2. **Fallback Chain**: Test strategy execution order
3. **Retry Behavior**: Verify exponential backoff
4. **Warning Display**: Ensure warnings shown appropriately
5. **Registry-Specific**: Test per-registry insecure settings

### Test Execution
```bash
# Run unit tests
go test ./pkg/fallback/... -v -cover
go test ./pkg/insecure/... -v -cover

# Run with race detection
go test ./... -race

# Generate coverage report
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## ✅ Validation Checkpoints

### Checkpoint 1: After Core Implementation (Step 1-3)
- [ ] FallbackManager compiles and passes basic tests
- [ ] Strategies implement interface correctly
- [ ] InsecureHandler manages state properly
- [ ] Line count under 400

### Checkpoint 2: After Test Implementation (Step 4-6)
- [ ] All unit tests pass
- [ ] Coverage exceeds 85%
- [ ] No race conditions detected
- [ ] Line count under 700

### Checkpoint 3: Integration Verification
- [ ] Works with Wave 1 trust store
- [ ] Coordinates with E1.2.1 validation
- [ ] --insecure flag functions correctly
- [ ] Warnings display appropriately

### Checkpoint 4: Final Review
- [ ] All tests pass
- [ ] Documentation complete
- [ ] Line count verified with tool
- [ ] Ready for code review

## 🔗 Dependencies and Integration

### From Wave 1:
- **E1.1.1 (Kind Certificate Extraction)**:
  - Import: `github.com/jessesanford/idpbuilder/pkg/certs`
  - Use: Certificate retrieval for fallback attempts

- **E1.1.2 (Registry TLS Trust)**:
  - Import: `github.com/jessesanford/idpbuilder/pkg/certs`
  - Use: `TrustStoreManager` interface for trust configuration

### Coordination with E1.2.1:
- **Certificate Validation**: Fallback manager activated on validation failures
- **Shared Types**: May share error types and validation interfaces
- **Parallel Development**: Can proceed independently, integrate during testing

## 🚀 Implementation Sequence

1. **Setup** (30 min)
   - Create directory structure
   - Initialize go module if needed
   - Set up development environment

2. **Core Implementation** (3 hours)
   - Implement FallbackManager
   - Create fallback strategies
   - Add insecure handler

3. **Testing** (2 hours)
   - Write comprehensive unit tests
   - Achieve coverage targets
   - Fix any issues found

4. **Integration** (1 hour)
   - Test with Wave 1 components
   - Coordinate with E1.2.1
   - Verify end-to-end flow

5. **Documentation** (30 min)
   - Update code comments
   - Create usage examples
   - Document security implications

6. **Review Preparation** (30 min)
   - Run line counter
   - Ensure all tests pass
   - Prepare for code review

## 🔒 Security Considerations

### Critical Security Points
1. **Insecure Mode Warnings**: Must be prominent and clear
2. **No Silent Failures**: Always notify user of fallback usage
3. **Audit Logging**: Log all certificate validation bypasses
4. **Configuration Safety**: Insecure mode should not persist
5. **Production Guards**: Consider environment-based restrictions

### Warning Messages
- Global insecure: Full-screen warning with multiple ⚠️ symbols
- Registry-specific: Clear indication of affected registry
- Self-signed acceptance: Explain security implications

## 📝 Notes for SW Engineer

1. **Start with the manager.go** - It's the core orchestrator
2. **Keep strategies simple** - Each does one thing well
3. **Test retry logic thoroughly** - Edge cases are important
4. **Make warnings impossible to miss** - Security is critical
5. **Coordinate with E1.2.1** - Share error types if beneficial
6. **Use feature flag** - `FALLBACK_STRATEGIES_ENABLED` for safe rollout
7. **Monitor line count** - Check after each major component

## 🎯 Success Criteria

- [ ] All fallback strategies implemented and tested
- [ ] --insecure flag works globally and per-registry
- [ ] Retry logic with exponential backoff functional
- [ ] Clear security warnings displayed
- [ ] 85%+ test coverage achieved
- [ ] Under 800 lines total (target: 700)
- [ ] Integrates cleanly with Wave 1 components
- [ ] Code review passed on first submission
>>>>>>> origin/idpbuilder-oci-build-push/phase2/wave1/gitea-client-split-001
