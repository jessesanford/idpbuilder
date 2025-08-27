# E3.1.1 Certificate Contracts & APIs Implementation Plan

## 🚨 CRITICAL EFFORT METADATA (FROM WAVE PLAN)
**Branch**: idpbuidler-oci-mgmt/phase3/wave1/E3.1.1-certificate-contracts  
**Can Parallelize**: No (blocks all other efforts)  
**Parallel With**: None  
**Size Estimate**: 400 lines  
**Dependencies**: None (foundation effort)  
**Working Directory**: /home/vscode/workspaces/idpbuilder-oci-mgmt/efforts/phase3/wave1/E3.1.1-certificate-contracts

## Overview
- **Effort**: Phase 3 Contracts & APIs - Foundation interfaces for certificate management, resilience, optimization, and observability
- **Phase**: 3, Wave: 1
- **Estimated Size**: 400 lines total
- **Implementation Time**: 4 hours
- **Criticality**: BLOCKING - All other Phase 3 efforts depend on these interfaces

## File Structure
```
pkg/oci/api/v2/
├── certificate_service.go (150 lines)
│   └── CertificateService interface
│   └── Certificate data structures
│   └── Verification modes and configurations
├── resilience_service.go (120 lines)
│   └── ResilienceService interface
│   └── Circuit breaker configurations
│   └── Retry policies and error types
├── optimization_service.go (100 lines)
│   └── OptimizationService interface
│   └── Cache configurations
│   └── Performance metrics structures
└── observability_service.go (30 lines)
    └── ObservabilityService interface
    └── Span and Alert type definitions
```

## Implementation Steps

### Step 1: Create Package Structure
1. Create `pkg/oci/api/v2/` directory
2. Initialize module with proper imports
3. Set up common types file for shared structures

### Step 2: Implement CertificateService Interface (150 lines)
**File**: `pkg/oci/api/v2/certificate_service.go`

```go
package v2

import (
    "context"
    "crypto/tls"
    "crypto/x509"
    "time"
)

// CertificateService defines the contract for certificate management
type CertificateService interface {
    // LoadCertificateBundle loads certificates from various formats
    LoadCertificateBundle(ctx context.Context, path string, format CertFormat) (*CertBundle, error)
    
    // SetVerificationMode configures how certificates are verified
    SetVerificationMode(mode VerificationMode) error
    
    // ValidateCertificate validates a single certificate
    ValidateCertificate(cert *x509.Certificate) error
    
    // LoadGiteaCertificate auto-discovers and loads Gitea certificates
    LoadGiteaCertificate(ctx context.Context, giteaURL string) (*CertBundle, error)
    
    // GetTLSConfig returns configured TLS configuration
    GetTLSConfig() (*tls.Config, error)
    
    // AddCACertificate adds a CA to the certificate pool
    AddCACertificate(cert *x509.Certificate) error
    
    // RemoveCACertificate removes a CA from the pool
    RemoveCACertificate(cert *x509.Certificate) error
    
    // ListCertificates returns all managed certificates
    ListCertificates() ([]*x509.Certificate, error)
    
    // RotateCertificate handles certificate rotation
    RotateCertificate(old, new *x509.Certificate) error
    
    // GetCertificatePool returns the current certificate pool
    GetCertificatePool() *x509.CertPool
}

// CertFormat represents supported certificate formats
type CertFormat string

const (
    CertFormatPEM    CertFormat = "pem"
    CertFormatDER    CertFormat = "der"
    CertFormatPKCS7  CertFormat = "pkcs7"
    CertFormatPKCS12 CertFormat = "pkcs12"
)

// VerificationMode defines certificate verification strategies
type VerificationMode string

const (
    VerificationModeStrict     VerificationMode = "strict"
    VerificationModePermissive VerificationMode = "permissive"
    VerificationModeSkip       VerificationMode = "skip"
    VerificationModeCustomCA   VerificationMode = "custom-ca"
)

// CertBundle represents a collection of certificates
type CertBundle struct {
    Certificates []*x509.Certificate
    CAs          []*x509.Certificate
    Format       CertFormat
    LoadedAt     time.Time
    Source       string
}

// CertificateConfig holds certificate service configuration
type CertificateConfig struct {
    BundlePath           string
    VerificationMode     VerificationMode
    AutoDiscoverGitea    bool
    SkipVerifyFallback   bool
    RefreshInterval      time.Duration
    CustomCAPath         string
}

// CertificateError represents certificate-specific errors
type CertificateError struct {
    Code    string
    Message string
    Cert    *x509.Certificate
    Err     error
}

func (e *CertificateError) Error() string {
    return e.Message
}
```

### Step 3: Implement ResilienceService Interface (120 lines)
**File**: `pkg/oci/api/v2/resilience_service.go`

```go
package v2

import (
    "context"
    "time"
)

// ResilienceService defines the contract for resilience patterns
type ResilienceService interface {
    // CreateCircuitBreaker creates a new circuit breaker
    CreateCircuitBreaker(name string, config CircuitBreakerConfig) (CircuitBreaker, error)
    
    // GetCircuitBreaker retrieves an existing circuit breaker
    GetCircuitBreaker(name string) (CircuitBreaker, error)
    
    // CreateRetryPolicy creates a new retry policy
    CreateRetryPolicy(name string, config RetryConfig) (RetryPolicy, error)
    
    // ExecuteWithResilience wraps an operation with resilience patterns
    ExecuteWithResilience(ctx context.Context, op Operation) error
    
    // RegisterHealthCheck registers a health check
    RegisterHealthCheck(name string, check HealthCheck) error
    
    // GetHealthStatus returns current health status
    GetHealthStatus() HealthStatus
    
    // ResetCircuitBreaker manually resets a circuit breaker
    ResetCircuitBreaker(name string) error
    
    // GetStatistics returns resilience statistics
    GetStatistics() ResilienceStats
}

// CircuitBreaker interface for circuit breaker pattern
type CircuitBreaker interface {
    Execute(ctx context.Context, op Operation) error
    GetState() CircuitState
    Reset() error
    GetStatistics() CircuitStats
}

// CircuitBreakerConfig defines circuit breaker configuration
type CircuitBreakerConfig struct {
    FailureThreshold   int
    SuccessThreshold   int
    Timeout           time.Duration
    ResetTimeout      time.Duration
    HalfOpenRequests  int
}

// CircuitState represents circuit breaker states
type CircuitState string

const (
    CircuitStateClosed   CircuitState = "closed"
    CircuitStateOpen     CircuitState = "open"
    CircuitStateHalfOpen CircuitState = "half-open"
)

// RetryPolicy interface for retry strategies
type RetryPolicy interface {
    Execute(ctx context.Context, op Operation) error
    GetConfig() RetryConfig
}

// RetryConfig defines retry policy configuration
type RetryConfig struct {
    MaxAttempts     int
    InitialDelay    time.Duration
    MaxDelay        time.Duration
    Multiplier      float64
    Jitter          float64
    RetryableErrors []string
}

// Operation represents a resilient operation
type Operation func(ctx context.Context) error

// HealthCheck represents a health check function
type HealthCheck func(ctx context.Context) error

// HealthStatus represents overall health status
type HealthStatus struct {
    Status     string
    Checks     map[string]CheckResult
    LastUpdate time.Time
}

// CheckResult represents individual check result
type CheckResult struct {
    Status  string
    Message string
    Error   error
}

// ResilienceStats provides resilience statistics
type ResilienceStats struct {
    CircuitBreakers map[string]CircuitStats
    RetryPolicies   map[string]RetryStats
}

// CircuitStats provides circuit breaker statistics
type CircuitStats struct {
    State            CircuitState
    FailureCount     int64
    SuccessCount     int64
    LastStateChange  time.Time
}

// RetryStats provides retry policy statistics
type RetryStats struct {
    TotalAttempts   int64
    SuccessfulRetries int64
    FailedRetries    int64
}
```

### Step 4: Implement OptimizationService Interface (100 lines)
**File**: `pkg/oci/api/v2/optimization_service.go`

```go
package v2

import (
    "context"
    "time"
)

// OptimizationService defines the contract for performance optimization
type OptimizationService interface {
    // AnalyzeBuild analyzes a build for optimization opportunities
    AnalyzeBuild(ctx context.Context, dockerfile string) (*BuildAnalysis, error)
    
    // OptimizeBuildOrder optimizes build step ordering
    OptimizeBuildOrder(steps []BuildStep) []BuildStep
    
    // PredictCacheHit predicts cache hit probability
    PredictCacheHit(layer LayerInfo) float64
    
    // GetCacheStatistics returns cache performance statistics
    GetCacheStatistics() CacheStats
    
    // EnablePersistentCache enables persistent caching
    EnablePersistentCache(config CacheConfig) error
    
    // ClearCache clears optimization caches
    ClearCache(cacheType CacheType) error
    
    // GetPerformanceMetrics returns performance metrics
    GetPerformanceMetrics() PerformanceMetrics
}

// BuildAnalysis represents build optimization analysis
type BuildAnalysis struct {
    OptimizationOpportunities []Optimization
    EstimatedTimeSaved       time.Duration
    CacheableStages          []string
    ParallelizableStages     [][]string
}

// Optimization represents an optimization opportunity
type Optimization struct {
    Type        OptimizationType
    Description string
    Impact      ImpactLevel
    Stage       string
}

// OptimizationType defines types of optimizations
type OptimizationType string

const (
    OptimizationTypeLayerReorder    OptimizationType = "layer-reorder"
    OptimizationTypeParallelization OptimizationType = "parallelization"
    OptimizationTypeCacheReuse      OptimizationType = "cache-reuse"
    OptimizationTypeMultiStage      OptimizationType = "multi-stage"
)

// ImpactLevel represents optimization impact
type ImpactLevel string

const (
    ImpactLevelHigh   ImpactLevel = "high"
    ImpactLevelMedium ImpactLevel = "medium"
    ImpactLevelLow    ImpactLevel = "low"
)

// BuildStep represents a build step
type BuildStep struct {
    Command      string
    Dependencies []string
    Cacheable    bool
    EstimatedTime time.Duration
}

// LayerInfo represents layer information
type LayerInfo struct {
    Digest      string
    Size        int64
    Created     time.Time
    Command     string
}

// CacheConfig defines cache configuration
type CacheConfig struct {
    PersistentPath string
    MaxSize        int64
    TTL            time.Duration
    CompressionEnabled bool
}

// CacheType defines cache types
type CacheType string

const (
    CacheTypeLayer    CacheType = "layer"
    CacheTypeBuild    CacheType = "build"
    CacheTypeMetadata CacheType = "metadata"
)

// CacheStats provides cache statistics
type CacheStats struct {
    HitRate      float64
    MissRate     float64
    TotalHits    int64
    TotalMisses  int64
    CacheSize    int64
    EvictionCount int64
}

// PerformanceMetrics provides performance metrics
type PerformanceMetrics struct {
    AverageBuildTime   time.Duration
    CacheHitRate       float64
    OptimizationSavings time.Duration
    ResourceUsage      ResourceMetrics
}

// ResourceMetrics tracks resource usage
type ResourceMetrics struct {
    CPUUsage    float64
    MemoryUsage int64
    DiskIO      int64
}
```

### Step 5: Implement ObservabilityService Interface (30 lines)
**File**: `pkg/oci/api/v2/observability_service.go`

```go
package v2

import (
    "context"
)

// ObservabilityService defines the contract for observability
type ObservabilityService interface {
    // StartSpan starts a new trace span
    StartSpan(ctx context.Context, name string) (Span, context.Context)
    
    // RecordMetric records a metric value
    RecordMetric(name string, value float64, labels map[string]string)
    
    // CreateAlert creates an alert
    CreateAlert(alert Alert) error
    
    // GetMetrics returns current metrics
    GetMetrics() map[string]interface{}
}

// Span represents a trace span
type Span interface {
    End()
    SetAttribute(key string, value interface{})
    RecordError(err error)
}

// Alert represents an alert configuration
type Alert struct {
    Name        string
    Condition   string
    Threshold   float64
    Duration    string
    Labels      map[string]string
}
```

## Implementation Order and Guidelines

### Implementation Priority
1. **MUST implement all interfaces first** - Other efforts depend on these contracts
2. Start with data structures and types
3. Add comprehensive godoc comments for each interface method
4. Include validation methods where appropriate
5. Create mock implementations for testing (in separate mock file)

### Key Implementation Guidelines
1. **Strict adherence to 400-line limit** (currently at exactly 400 lines)
2. Use standard Go idioms and patterns
3. Ensure all interfaces are well-documented
4. Design for extensibility and backward compatibility
5. Include error types for each service domain

### File Creation Order
1. Create `pkg/oci/api/v2/` directory structure
2. Implement `certificate_service.go` (150 lines)
3. Implement `resilience_service.go` (120 lines)
4. Implement `optimization_service.go` (100 lines)
5. Implement `observability_service.go` (30 lines)

## Size Management
- **Target Lines**: 400 lines (exactly as estimated)
- **Current Breakdown**:
  - Certificate Service: 150 lines
  - Resilience Service: 120 lines
  - Optimization Service: 100 lines
  - Observability Service: 30 lines
- **Measurement Tool**: Use line counter after implementation
- **Check Frequency**: After each file creation
- **Split Threshold**: N/A (well within limit)

## Test Requirements
- **Unit Tests**: 90% coverage required
- **Test Files**:
  - `certificate_service_test.go` - Validate data structures
  - `resilience_service_test.go` - Test configuration validation
  - `optimization_service_test.go` - Test metric calculations
  - `observability_service_test.go` - Test span operations
- **Mock Implementations**: Create `mocks/` subdirectory with mock implementations

## Pattern Compliance
- **Go Patterns**: 
  - Use interfaces for contracts
  - Embed error information in custom error types
  - Use context for cancellation and timeouts
  - Follow Go naming conventions
- **Security Requirements**:
  - Certificate validation must be thorough
  - Support for custom CA certificates
  - Secure defaults (strict verification mode)
- **Performance Targets**:
  - Interfaces should not impose performance overhead
  - Design for efficient implementation

## Integration Points
- These interfaces will be imported by ALL other Phase 3 efforts
- Phase 2 services will be wrapped with these new interfaces
- Mock implementations enable parallel development of dependent efforts

## Success Criteria
- All 4 interface files created and properly documented
- Total implementation exactly 400 lines
- All data structures well-defined
- Mock implementations available for testing
- Can be imported successfully by other efforts
- Enables parallel development of E3.1.2 through E3.1.5

## Next Steps After Implementation
1. Verify all interfaces compile correctly
2. Create basic mock implementations
3. Document usage examples
4. Enable other efforts to begin development
5. Update work-log.md with completion status