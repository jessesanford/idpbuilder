# Effort 1: Advanced Build Contracts & Interfaces Implementation Plan

## =¨ CRITICAL EFFORT METADATA (FROM WAVE PLAN)
**Branch**: `idpbuidler-oci-mgmt/phase2/wave2/effort1-contracts`
**Can Parallelize**: No
**Parallel With**: None (MUST BE FIRST - BLOCKS ALL!)
**Size Estimate**: 400 lines (MUST be <800)
**Dependencies**: None

## Overview
- **Effort**: Define all contracts and interfaces for Phase 2 Wave 2 advanced OCI build capabilities
- **Phase**: 2, Wave: 2
- **Estimated Size**: 400 lines (interfaces, models, and comprehensive documentation)
- **Implementation Time**: 6 hours
- **Criticality**: BLOCKING - All other efforts depend on these contracts

## <¯ Mission Critical Requirements

### Why This Effort Must Be First
1. **Defines ALL interfaces** that efforts 2-5 will implement
2. **Establishes data models** shared across all efforts
3. **Sets API contracts** for multi-stage optimization, caching, security, and registry operations
4. **Prevents integration issues** by establishing contracts upfront
5. **Enables parallel development** of efforts 2-4 once contracts are complete

### Contract Completeness Checklist
- [ ] StageOptimizer interface with all methods from architecture
- [ ] CacheManager interface with layer operations
- [ ] SecurityManager interface with signing/verification
- [ ] RegistryClient interface with push/pull/auth
- [ ] ALL data models with proper JSON/YAML tags
- [ ] Validation methods where needed
- [ ] Error types defined
- [ ] Constants for common values

## File Structure

### Core Contract Files
- `pkg/oci/api/optimizer.go`: Multi-stage optimization interfaces (80 lines)
- `pkg/oci/api/cache.go`: Cache management interfaces (70 lines)
- `pkg/oci/api/security.go`: Security and signing interfaces (90 lines)
- `pkg/oci/api/registry.go`: Registry client interfaces (80 lines)
- `pkg/oci/api/models.go`: Shared data models (80 lines)

### Test Files
- `pkg/oci/api/optimizer_test.go`: Mock implementations and interface tests
- `pkg/oci/api/cache_test.go`: Cache interface validation
- `pkg/oci/api/security_test.go`: Security interface tests
- `pkg/oci/api/registry_test.go`: Registry interface tests

## Implementation Steps

### Step 1: Create Package Structure (10 minutes)
```bash
# Create the api package directory
mkdir -p pkg/oci/api

# Create placeholder files
touch pkg/oci/api/{optimizer,cache,security,registry,models}.go
touch pkg/oci/api/{optimizer,cache,security,registry}_test.go
```

### Step 2: Define Optimizer Interfaces (60 minutes)
```go
// pkg/oci/api/optimizer.go
package api

import (
    "context"
    "time"
)

// StageOptimizer handles multi-stage build optimization
type StageOptimizer interface {
    // Analyze Dockerfile for optimization opportunities
    AnalyzeStages(dockerfile []byte) (*StageAnalysis, error)
    
    // Build stages with parallelization where possible
    BuildStages(ctx context.Context, analysis *StageAnalysis, req *BuildRequest) (*StageResult, error)
    
    // Optimize stage dependencies for parallel execution
    OptimizeDependencies(stages []*Stage) (*DependencyGraph, error)
    
    // Get build metrics for optimization analysis
    GetMetrics() *BuildMetrics
}

type StageAnalysis struct {
    Stages          []*Stage               `json:"stages" yaml:"stages"`
    Dependencies    map[string][]string    `json:"dependencies" yaml:"dependencies"`
    ParallelGroups  [][]string            `json:"parallel_groups" yaml:"parallel_groups"`
    CacheableStages []string              `json:"cacheable_stages" yaml:"cacheable_stages"`
    EstimatedTime   time.Duration         `json:"estimated_time" yaml:"estimated_time"`
}
```

### Step 3: Define Cache Interfaces (45 minutes)
```go
// pkg/oci/api/cache.go
package api

// CacheManager handles layer caching operations
type CacheManager interface {
    // Check if a layer exists in cache
    HasLayer(digest string) bool
    
    // Retrieve a cached layer
    GetLayer(digest string) (*Layer, error)
    
    // Store a new layer in cache
    StoreLayer(layer *Layer) error
    
    // Calculate cache key for an instruction
    CalculateCacheKey(instruction string, context []byte) string
    
    // Prune old cache entries
    PruneCache(before time.Time) error
    
    // Get cache statistics
    GetStats() *CacheStats
}

type Layer struct {
    Digest      string            `json:"digest" yaml:"digest"`
    Size        int64             `json:"size" yaml:"size"`
    MediaType   string            `json:"media_type" yaml:"media_type"`
    Created     time.Time         `json:"created" yaml:"created"`
    LastUsed    time.Time         `json:"last_used" yaml:"last_used"`
    RefCount    int               `json:"ref_count" yaml:"ref_count"`
    Annotations map[string]string `json:"annotations,omitempty" yaml:"annotations,omitempty"`
}
```

### Step 4: Define Security Interfaces (60 minutes)
```go
// pkg/oci/api/security.go
package api

// SecurityManager handles image security operations
type SecurityManager interface {
    // Sign an image with the provided signer
    SignImage(ctx context.Context, image string, signer Signer) (*Signature, error)
    
    // Verify image signature
    VerifySignature(ctx context.Context, image string, verifier Verifier) error
    
    // Generate Software Bill of Materials
    GenerateSBOM(ctx context.Context, image string) (*SBOM, error)
    
    // Scan for vulnerabilities
    ScanVulnerabilities(ctx context.Context, image string) (*VulnerabilityReport, error)
    
    // Attach attestation to image
    AttachAttestation(ctx context.Context, image string, attestation *Attestation) error
}

type Signer interface {
    Sign(payload []byte) ([]byte, error)
    KeyID() string
    Algorithm() string
}

type Verifier interface {
    Verify(payload []byte, signature []byte) error
    TrustedKeys() []string
    VerifyPolicy(policy *Policy) error
}
```

### Step 5: Define Registry Interfaces (60 minutes)
```go
// pkg/oci/api/registry.go
package api

// RegistryClient handles registry operations
type RegistryClient interface {
    // Push image to registry
    Push(ctx context.Context, image string, auth AuthConfig) error
    
    // Pull image from registry
    Pull(ctx context.Context, image string, auth AuthConfig) (*Image, error)
    
    // Get image manifest
    GetManifest(ctx context.Context, image string) (*Manifest, error)
    
    // List repository tags
    ListTags(ctx context.Context, repository string) ([]string, error)
    
    // Delete image from registry
    Delete(ctx context.Context, image string, auth AuthConfig) error
    
    // Check registry health
    Ping(ctx context.Context) error
}

type AuthConfig struct {
    Username      string `json:"username,omitempty" yaml:"username,omitempty"`
    Password      string `json:"password,omitempty" yaml:"password,omitempty"`
    Auth          string `json:"auth,omitempty" yaml:"auth,omitempty"`
    ServerAddress string `json:"serveraddress,omitempty" yaml:"serveraddress,omitempty"`
    IdentityToken string `json:"identitytoken,omitempty" yaml:"identitytoken,omitempty"`
    RegistryToken string `json:"registrytoken,omitempty" yaml:"registrytoken,omitempty"`
}

// Validate checks if auth config has required fields
func (a AuthConfig) Validate() error {
    // Implementation here
}
```

### Step 6: Define Shared Models (60 minutes)
```go
// pkg/oci/api/models.go
package api

// Stage represents a build stage in a multi-stage Dockerfile
type Stage struct {
    Name         string            `json:"name" yaml:"name"`
    BaseImage    string            `json:"base_image" yaml:"base_image"`
    Instructions []string          `json:"instructions" yaml:"instructions"`
    Dependencies []string          `json:"dependencies" yaml:"dependencies"`
    Cacheable    bool              `json:"cacheable" yaml:"cacheable"`
    Size         int64             `json:"size,omitempty" yaml:"size,omitempty"`
    BuildArgs    map[string]string `json:"build_args,omitempty" yaml:"build_args,omitempty"`
}

// StageResult contains the results of building a stage
type StageResult struct {
    StageID     string        `json:"stage_id" yaml:"stage_id"`
    ImageID     string        `json:"image_id" yaml:"image_id"`
    Layers      []*Layer      `json:"layers" yaml:"layers"`
    BuildTime   time.Duration `json:"build_time" yaml:"build_time"`
    CacheHit    bool          `json:"cache_hit" yaml:"cache_hit"`
    Error       string        `json:"error,omitempty" yaml:"error,omitempty"`
}

// DependencyGraph represents stage dependencies
type DependencyGraph struct {
    Nodes    map[string]*Stage   `json:"nodes" yaml:"nodes"`
    Edges    map[string][]string `json:"edges" yaml:"edges"`
    Parallel [][]string          `json:"parallel" yaml:"parallel"`
}

// Additional models for security, registry, etc.
```

### Step 7: Create Mock Implementations for Testing (45 minutes)
```go
// pkg/oci/api/optimizer_test.go
package api

import (
    "context"
    "testing"
)

// MockStageOptimizer implements StageOptimizer for testing
type MockStageOptimizer struct {
    AnalyzeStagesFunc    func([]byte) (*StageAnalysis, error)
    BuildStagesFunc      func(context.Context, *StageAnalysis, *BuildRequest) (*StageResult, error)
    OptimizeDependencies func([]*Stage) (*DependencyGraph, error)
}

// Test interface compliance
func TestStageOptimizerInterface(t *testing.T) {
    var _ StageOptimizer = &MockStageOptimizer{}
}
```

### Step 8: Add Validation and Helper Methods (30 minutes)
- Add validation methods to models
- Create helper functions for common operations
- Define error types and constants
- Add comprehensive godoc comments

## Size Management
- **Estimated Lines**: 400 lines
- **Measurement Tool**: /home/vscode/workspaces/idpbuilder-oci-mgmt/tools/line-counter.sh
- **Check Frequency**: After each interface completion
- **Split Threshold**: N/A (contracts must remain together)

## Test Requirements

### Unit Tests (80% coverage minimum)
- Interface compliance verification
- Model validation tests
- JSON/YAML marshaling tests
- Error handling tests

### Mock Implementations
- Create mocks for all interfaces
- Use for testing other efforts
- Include in test files

### Integration Points
- Verify interfaces work with Wave 1 code
- Ensure models serialize correctly
- Test with example implementations

## Pattern Compliance

### Go Best Practices
- Use context.Context for cancellation
- Return errors as last return value
- Use pointer receivers for large structs
- Follow Go naming conventions

### Interface Design
- Keep interfaces small and focused
- Use composition over inheritance
- Define behavior, not implementation
- Include method documentation

### Error Handling
- Define custom error types where needed
- Use error wrapping for context
- Never panic in library code
- Document error conditions

## Integration Strategy

### For Effort 2-4 (Parallel Implementation)
```go
import "github.com/cnoe-io/idpbuilder/pkg/oci/api"

// Implement the interfaces
type stageOptimizer struct {
    // implementation
}

// Ensure interface compliance
var _ api.StageOptimizer = (*stageOptimizer)(nil)
```

### For Effort 5 (Integration)
```go
// Import all interfaces and implementations
import (
    "github.com/cnoe-io/idpbuilder/pkg/oci/api"
    "github.com/cnoe-io/idpbuilder/pkg/oci/optimizer"
    "github.com/cnoe-io/idpbuilder/pkg/oci/cache"
    "github.com/cnoe-io/idpbuilder/pkg/oci/security"
)

// Compose into registry client
type registryClient struct {
    optimizer api.StageOptimizer
    cache     api.CacheManager
    security  api.SecurityManager
}
```

## Success Criteria

### Must Have
- [x] All interfaces from architecture plan defined
- [x] All data models with proper tags
- [x] Validation methods implemented
- [x] Mock implementations for testing
- [x] Comprehensive godoc comments

### Should Have
- [ ] Helper functions for common operations
- [ ] Custom error types
- [ ] Constants for common values
- [ ] Example usage in comments

### Could Have
- [ ] Interface compliance tests
- [ ] Benchmark scaffolding
- [ ] Integration test helpers

## Risk Mitigation

### Risk 1: Incomplete Contracts
**Mitigation**: Review with all effort implementers before marking complete

### Risk 2: Interface Changes Later
**Mitigation**: Design for extensibility, use versioning if needed

### Risk 3: Integration Issues
**Mitigation**: Test with mock implementations early

## Next Steps
After this effort is complete:
1. Code review by orchestrator
2. Efforts 2-4 can begin in parallel
3. Each effort imports these contracts
4. Effort 5 integrates all implementations

## Document History
| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 1.0 | 2025-08-26 | Code Reviewer Agent | Initial plan for contracts & interfaces |