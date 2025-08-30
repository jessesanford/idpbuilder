# Certificate Integration Manager Implementation Plan

## 🚨 CRITICAL EFFORT METADATA (FROM WAVE PLAN)
**Branch**: `phase3/wave1/cert-integration-manager`
**Can Parallelize**: Yes
**Parallel With**: [security-features, error-messaging]
**Size Estimate**: 700 lines
**Dependencies**: None (uses Phase 1 and Phase 2 interfaces)

## Overview
- **Effort**: Certificate Integration Manager - Bridge between Phase 1 certificates and Phase 2 build/push
- **Phase**: 3, Wave: 1
- **Estimated Size**: 700 lines
- **Implementation Time**: 3-4 days (parallel with other efforts)

## Dependency Context

### Dependencies Analyzed
This effort has no direct dependencies on other Wave 1 efforts, allowing parallel execution. However, it integrates with:

### Phase 1 Interfaces (Existing)
- `cert-extraction.Client`: For certificate retrieval from Kind clusters
- `trust-store.Store`: For trust store management
- `validation.Validator`: For certificate validation
- `fallback.Strategy`: For fallback certificate scenarios

### Phase 2 Interfaces (Existing) 
- `buildah.Wrapper`: For build configuration with certificates
- `registry.Client`: For registry operations with authentication
- `cli.Context`: For command context and flags

### Integration Strategy
- Use existing Phase 1 extraction APIs to load certificates
- Apply certificates to Phase 2 build/push operations
- Provide unified interface for Wave 2 E2E testing

## File Structure
```
pkg/cert-integration/
├── loader.go           # Certificate loading interface and implementation (150 lines)
├── loader_test.go      # Tests for certificate loader (80 lines)
├── manager.go          # Main certificate manager implementation (200 lines)
├── manager_test.go     # Tests for certificate manager (100 lines)
├── config.go           # Configuration types and validation (100 lines)
├── config_test.go      # Tests for configuration (50 lines)
├── resolver.go         # Path resolution for certificates (100 lines)
├── resolver_test.go    # Tests for path resolver (50 lines)
├── validator.go        # Integration validation logic (150 lines)
└── validator_test.go   # Tests for validator (70 lines)
```

## Implementation Steps

### Step 1: Create Certificate Loader Interface (150 lines)
**File**: `pkg/cert-integration/loader.go`

1. Define `CertificateLoader` interface with methods:
   - `LoadFromExtraction(source string) (*CertificateSet, error)`
   - `LoadFromPath(path string) (*CertificateSet, error)`
   - `LoadFromKindCluster(clusterName string) (*CertificateSet, error)`

2. Implement `certificateLoader` struct with:
   - Integration with Phase 1 cert-extraction client
   - Certificate parsing and validation
   - Support for multiple certificate formats (PEM, DER)
   - Error handling with context

3. Define `CertificateSet` struct containing:
   - RootCA certificate
   - Intermediate certificates array
   - Server certificate
   - Client certificates map
   - Trust bundle byte array

4. Implement helper methods:
   - `buildCertificateSet()`: Organize loaded certificates
   - `loadCertificate()`: Load single certificate file
   - `categorizeCertificate()`: Determine certificate type

### Step 2: Implement Certificate Manager Core (200 lines)
**File**: `pkg/cert-integration/manager.go`

1. Create `CertificateManager` struct with:
   - Certificate loader instance
   - Path resolver instance
   - Integration validator instance
   - Manager configuration
   - Certificate cache map

2. Implement `NewCertificateManager(config *ManagerConfig)` constructor:
   - Initialize all components
   - Set up default configuration
   - Create certificate cache

3. Add `ConfigureRegistry(registry string, certs *CertificateSet)` method:
   - Map certificates to registry configuration
   - Write certificates to expected locations
   - Validate registry connectivity
   - Handle authentication setup

4. Add `ConfigureBuild(context *BuildContext, certs *CertificateSet)` method:
   - Configure buildah with certificates
   - Set up trust store for builds
   - Apply certificate overrides
   - Integrate with Phase 2 buildah wrapper

5. Implement support methods:
   - `writeCertificates()`: Write certs to filesystem
   - `setupBuildahTrust()`: Configure buildah trust
   - `cacheOperations()`: Handle certificate caching

### Step 3: Create Configuration Types (100 lines)
**File**: `pkg/cert-integration/config.go`

1. Define `ManagerConfig` struct:
   - DefaultCertPath string
   - TrustStoreLocation string
   - ValidationMode (Strict/Permissive/Disabled)
   - CacheEnabled boolean

2. Define `RegistryConfig` struct:
   - RegistryURL string
   - Insecure boolean flag
   - CertPath string
   - AuthMethod string

3. Define `BuildConfig` struct:
   - BuildContext string
   - TrustStore string
   - CertOverride map[string]string

4. Implement validation methods:
   - `Validate()` for CertificateSet
   - Certificate chain verification
   - Expiry checking
   - Format validation

### Step 4: Implement Path Resolver (100 lines)
**File**: `pkg/cert-integration/resolver.go`

1. Create `PathResolver` interface with methods:
   - `GetRegistryCertPath(registry string) string`
   - `GetBuildTrustStore() string`
   - `ResolveCertificatePath(hint string) (string, error)`

2. Implement `pathResolver` struct with:
   - Base path configuration
   - Certificate locations map
   - Fallback strategies

3. Add resolution logic:
   - Check environment variables (IDPBUILDER_CERT_PATH)
   - Try standard locations (/etc/idpbuilder/certs/)
   - Support user home directory (~/.idpbuilder/certs)
   - Handle temporary locations (/tmp/idpbuilder-certs)

4. Registry-specific paths:
   - Map Gitea to /etc/idpbuilder/certs/gitea
   - Map Harbor to /etc/idpbuilder/certs/harbor
   - Default pattern: /etc/containers/certs.d/{registry}

### Step 5: Add Integration Validator (150 lines)
**File**: `pkg/cert-integration/validator.go`

1. Define `IntegrationValidator` interface:
   - `ValidateIntegration() error`
   - `ValidateRegistryConnection(config *RegistryConfig) error`
   - `ValidateBuildConfiguration(config *BuildConfig) error`

2. Implement validation checks:
   - Phase 1 components accessibility
   - Phase 2 components compatibility
   - Certificate format validation
   - Certificate chain validation

3. Add registry connection testing:
   - TLS handshake validation
   - Certificate matching
   - Authentication verification
   - Timeout handling

4. Build configuration validation:
   - Trust store accessibility
   - Certificate permissions
   - Buildah compatibility
   - Container runtime checks

### Step 6: Testing Implementation (350 lines across all test files)

1. **Unit Tests for Loader** (`loader_test.go`):
   - Test loading from extraction
   - Test loading from filesystem
   - Test invalid certificate handling
   - Test format detection

2. **Unit Tests for Manager** (`manager_test.go`):
   - Test registry configuration
   - Test build configuration
   - Test caching behavior
   - Test error scenarios

3. **Unit Tests for Config** (`config_test.go`):
   - Test validation logic
   - Test configuration merging
   - Test default values

4. **Unit Tests for Resolver** (`resolver_test.go`):
   - Test path resolution
   - Test fallback strategies
   - Test environment variable handling

5. **Unit Tests for Validator** (`validator_test.go`):
   - Test integration validation
   - Test registry connectivity
   - Test build validation
   - Mock Phase 1/2 components

### Step 7: Integration and Documentation

1. Create integration examples
2. Document public APIs
3. Add inline code documentation
4. Create usage examples

## Size Management
- **Estimated Lines**: 700 (including tests)
- **Breakdown**:
  - Core implementation: 350 lines
  - Test implementation: 350 lines
- **Measurement Tool**: /home/vscode/workspaces/idpbuilder-oci-mvp/tools/line-counter.sh
- **Check Frequency**: After each component completion
- **Split Threshold**: 700 lines (warning), 800 lines (stop)

## Test Requirements
- **Unit Tests**: 85% coverage minimum
- **Integration Tests**: Mock Phase 1/2 components
- **E2E Tests**: Will be covered in Wave 2
- **Test Files**: 
  - loader_test.go (80 lines)
  - manager_test.go (100 lines)
  - config_test.go (50 lines)
  - resolver_test.go (50 lines)
  - validator_test.go (70 lines)

## Pattern Compliance
- **Go Patterns**: 
  - Interface-driven design
  - Dependency injection
  - Error wrapping with context
  - Structured logging
- **Security Requirements**:
  - Certificate validation
  - Permission checking
  - Secure file operations
  - Audit logging for certificate operations
- **Performance Targets**:
  - Certificate loading < 100ms
  - Registry validation < 500ms
  - Caching for repeated operations

## Integration Points

### Inputs (FROM Phase 1 & 2)
- Phase 1 cert-extraction.Client for certificate retrieval
- Phase 1 trust-store.Store for trust management
- Phase 2 buildah.Wrapper for build configuration
- Phase 2 registry.Client for registry operations

### Outputs (TO Wave 2 & Other Components)
- CertificateManager interface for E2E tests
- CertificateLoader for other components
- IntegrationValidator for validation
- Unified certificate configuration

## Risk Mitigation
1. **Certificate Format Issues**: Support multiple formats (PEM, DER)
2. **Path Resolution Failures**: Multiple fallback strategies
3. **Integration Complexity**: Clear interface boundaries
4. **Size Overrun**: Pre-identified split points if needed

## Success Criteria
- [ ] All certificates loaded successfully from Phase 1
- [ ] Certificates applied correctly to Phase 2 components
- [ ] Registry connection validated with certificates
- [ ] Build configuration includes proper trust store
- [ ] All unit tests passing with >85% coverage
- [ ] Size remains under 700 lines
- [ ] Clean integration with no compilation errors

## Implementation Order
1. Start with configuration types (foundation)
2. Implement loader (core functionality)
3. Build manager (orchestration)
4. Add resolver (support)
5. Complete validator (verification)
6. Write comprehensive tests
7. Document and integrate

---
*Implementation Plan Created: 2025-08-30*
*Code Reviewer: @agent-code-reviewer*
*State: EFFORT_PLAN_CREATION*