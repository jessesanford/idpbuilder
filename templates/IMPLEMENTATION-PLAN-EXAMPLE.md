# IMPLEMENTATION PLAN: Container Registry for IDP Builder

## 1. Project Overview

This project adds container registry capabilities to the CNOE IDP Builder, enabling it to serve as a local OCI registry for development environments. The IDP Builder currently provides a complete Internal Developer Platform but lacks integrated container image management. This enhancement will allow developers to push, pull, and manage container images locally without external dependencies, making the IDP Builder a more complete solution for cloud-native development.

The implementation will integrate a lightweight, embedded registry server that starts alongside other IDP Builder components. It will provide full OCI compliance, support for multi-architecture images, and integrate with the existing Kind cluster for seamless Kubernetes deployments.

### Key Features
- **Embedded Registry Server**: Lightweight registry that starts with `idpbuilder create`
- **OCI Compliance**: Full support for OCI image and distribution specifications
- **Kind Integration**: Automatic configuration for pulling images from local registry
- **Security**: Built-in TLS support with automatic certificate generation

### Target Users
- Platform engineers setting up development environments
- Developers needing local container registries
- CI/CD pipelines requiring fast, local image storage

## 2. Goals and Objectives

### Primary Objectives
1. **Add Registry Server**: Embed a registry server that lifecycle manages with IDP Builder
2. **Ensure OCI Compliance**: Support standard docker push/pull operations
3. **Integrate with Kind**: Configure Kind cluster to use local registry
4. **Maintain Simplicity**: Zero additional configuration for basic use cases

### Secondary Goals
- Support multi-architecture images (AMD64, ARM64)
- Provide registry UI for browsing images
- Enable registry persistence across restarts
- Support image garbage collection

### Non-Goals (Out of Scope)
- Multi-node registry clustering
- External authentication providers (LDAP, OIDC)
- Image scanning or vulnerability assessment
- Registry replication or mirroring

## 3. Technical Architecture

### Technology Stack
- **Primary Language**: Go 1.21+
- **Framework**: Cobra (CLI), Gin (HTTP)
- **Build System**: Make
- **Testing**: Go test with testify
- **Registry**: distribution/distribution v3
- **Deployment**: Binary distribution, Docker

### Architecture Pattern
The registry will follow IDP Builder's plugin architecture, implementing the `Component` interface for lifecycle management. It will run as a separate goroutine within the main process, sharing the configuration and state management system.

### System Components
1. **Registry Server**: Core OCI registry using distribution/distribution
2. **Configuration Manager**: Handles registry config and Kind integration
3. **Certificate Manager**: Generates and manages TLS certificates
4. **CLI Extensions**: New commands for registry management

### External Dependencies
- distribution/distribution v3: OCI registry implementation
- go-containerregistry: Client library for OCI operations
- Kind API: For cluster configuration updates

## 4. Implementation Phases

### Phase 1: Foundation and Core Registry
**Goal**: Establish basic registry server with lifecycle management
**Duration Estimate**: 2 weeks

#### Wave 1.1: Project Setup and Registry Integration
- **Effort 1.1.1**: Add registry dependencies and vendor management
  - Add distribution/distribution to go.mod
  - Update vendor directory
  - Configure build tags for registry

- **Effort 1.1.2**: Implement registry component interface
  - Create pkg/registry package structure
  - Implement Component interface
  - Add configuration structures

- **Effort 1.1.3**: Basic registry server startup
  - Implement Start() and Stop() methods
  - Configure in-memory storage backend
  - Add health check endpoint

#### Wave 1.2: CLI and Configuration
- **Effort 1.2.1**: Add registry CLI commands
  - Implement `idpbuilder create --registry` flag
  - Add `idpbuilder registry status` command
  - Add `idpbuilder registry list` for images

- **Effort 1.2.2**: Registry configuration management
  - Define registry configuration schema
  - Implement configuration persistence
  - Add port and storage configuration

- **Effort 1.2.3**: Integration tests for basic operations
  - Test registry lifecycle
  - Test basic push/pull operations
  - Validate configuration persistence

### Phase 2: Kind Integration and Security
**Goal**: Integrate registry with Kind cluster and add security features
**Duration Estimate**: 2 weeks

#### 📝 PHASE 2 VALIDATION TESTS (WRITE FIRST!)
```go
// Test 1: Complete Kind integration
func TestPhase2_KindIntegration(t *testing.T) {
    builder := setupBuilderWithRegistry(t)

    // Push test image to registry
    testImage := "localhost:5000/test:v1"
    pushTestImage(t, testImage)

    // Deploy to Kind using local registry
    deployment := newDeployment(testImage)
    err := applyToKind(builder, deployment)
    require.NoError(t, err, "Kind should use local registry")

    // Verify pods running with local image
    pods := getPodsWithImage(builder, testImage)
    require.Len(t, pods, 1, "Pod should be running with local image")
    require.Equal(t, "Running", pods[0].Status)
}

// Test 2: TLS security
func TestPhase2_TLSSecurity(t *testing.T) {
    registry := setupSecureRegistry(t)

    // Test HTTPS endpoint
    resp, err := https.Get(registry.GetSecureURL() + "/v2/")
    require.NoError(t, err)
    require.Equal(t, 200, resp.StatusCode)

    // Test certificate validation
    cert := registry.GetCertificate()
    require.NotEmpty(t, cert.Subject)
    require.True(t, cert.NotAfter.After(time.Now()))
}
```

#### Wave 2.1: Kind Cluster Integration
- **Effort 2.1.1**: Kind configuration patches
  - Generate containerd config for registry
  - Add registry to Kind's known registries
  - Configure insecure registry for local development

- **Effort 2.1.2**: Automatic registry discovery
  - Implement registry service in Kind cluster
  - Add DNS configuration for registry.local
  - Create ConfigMap with registry details

- **Effort 2.1.3**: Image preloading capabilities
  - Add `idpbuilder registry push` command
  - Implement image loading into Kind
  - Support for multi-architecture images

#### Wave 2.2: Security and TLS
- **Effort 2.2.1**: TLS certificate generation
  - Implement self-signed certificate generation
  - Add certificate rotation logic
  - Store certificates securely

- **Effort 2.2.2**: Registry authentication
  - Implement basic auth for registry
  - Add htpasswd file management
  - Create credential helper for Docker

- **Effort 2.2.3**: Security hardening
  - Add rate limiting
  - Implement access logging
  - Add security headers

### Phase 3: Production Features and Polish
**Goal**: Add advanced features and prepare for production use
**Duration Estimate**: 1.5 weeks

#### 📝 PHASE 3 VALIDATION TESTS (WRITE FIRST!)
```go
// Test 1: Registry UI functionality
func TestPhase3_RegistryUI(t *testing.T) {
    registry := setupRegistryWithUI(t)

    // Test UI accessibility
    resp, err := http.Get(registry.GetUIURL())
    require.NoError(t, err)
    require.Equal(t, 200, resp.StatusCode)

    // Test image browsing
    images := getImagesFromUI(registry)
    require.Greater(t, len(images), 0, "UI should show images")
}

// Test 2: Production monitoring
func TestPhase3_Monitoring(t *testing.T) {
    registry := setupRegistryWithMonitoring(t)

    // Test metrics endpoint
    metrics := registry.GetMetrics()
    require.Contains(t, metrics, "registry_storage_blob_upload_total")
    require.Contains(t, metrics, "registry_http_request_duration_seconds")

    // Test health monitoring
    health := registry.GetDetailedHealth()
    require.Equal(t, "healthy", health.Status)
    require.NotEmpty(t, health.Checks)
}

// Test 3: Performance and scale
func TestPhase3_Performance(t *testing.T) {
    registry := setupProductionRegistry(t)

    // Load test: Push 100 images concurrently
    var wg sync.WaitGroup
    errors := make(chan error, 100)

    start := time.Now()
    for i := 0; i < 100; i++ {
        wg.Add(1)
        go func(n int) {
            defer wg.Done()
            err := registry.Push(fmt.Sprintf("load-test:v%d", n))
            if err != nil {
                errors <- err
            }
        }(i)
    }

    wg.Wait()
    close(errors)

    // Assert performance
    duration := time.Since(start)
    require.Less(t, duration, 30*time.Second, "Should handle 100 images in 30s")
    require.Len(t, errors, 0, "No errors during load test")
}
```

#### Wave 3.1: Advanced Features
- **Effort 3.1.1**: Persistent storage backend
  - Implement filesystem storage driver
  - Add garbage collection command
  - Configure storage quotas

- **Effort 3.1.2**: Registry UI integration
  - Add optional registry UI component
  - Configure proxy routing
  - Add UI authentication

- **Effort 3.1.3**: Monitoring and metrics
  - Add Prometheus metrics endpoint
  - Implement usage statistics
  - Add performance monitoring

#### Wave 3.2: Documentation and Examples
- **Effort 3.2.1**: User documentation
  - Write registry usage guide
  - Add troubleshooting section
  - Create migration guide

- **Effort 3.2.2**: Example applications
  - Create sample application with registry usage
  - Add CI/CD pipeline examples
  - Provide multi-arch build examples

## 5. Success Criteria

### Phase 1 Completion Criteria
- [ ] Registry starts and stops with idpbuilder
- [ ] Basic push/pull operations work
- [ ] CLI commands implemented and tested
- [ ] Unit tests >80% coverage

### Phase 2 Completion Criteria
- [ ] Kind cluster can pull images from registry
- [ ] TLS certificates generated automatically
- [ ] Authentication working for push operations
- [ ] Integration tests passing

### Phase 3 Completion Criteria
- [ ] Persistent storage working across restarts
- [ ] Registry UI accessible and functional
- [ ] Documentation complete with examples
- [ ] Performance benchmarks met (<100ms push latency)

### Overall Project Success Metrics
- **Adoption**: >50% of IDP Builder users enable registry
- **Performance**: Support 100 concurrent operations
- **Reliability**: 99.9% uptime in development use
- **Size**: Registry adds <10MB to binary size

## 6. Risk Mitigation

### Technical Risks
1. **Risk**: Distribution/distribution library compatibility issues
   - **Impact**: High
   - **Mitigation**: Vendor specific version, maintain compatibility matrix

2. **Risk**: Kind networking complications with registry
   - **Impact**: Medium
   - **Mitigation**: Implement multiple connection methods, provide debugging tools

3. **Risk**: Certificate management complexity
   - **Impact**: Medium
   - **Mitigation**: Use established libraries, provide manual override options

### Schedule Risks
1. **Risk**: Upstream IDP Builder changes during development
   - **Impact**: Medium
   - **Mitigation**: Regular rebasing, maintain feature flag for isolation

### External Dependencies
1. **Risk**: Distribution/distribution maintenance concerns
   - **Impact**: Low
   - **Mitigation**: Consider alternative implementations if needed

## 7. Appendices

### A. Glossary
- **OCI**: Open Container Initiative
- **Kind**: Kubernetes in Docker
- **IDP**: Internal Developer Platform
- **CNOE**: Cloud Native Operational Excellence

### B. References
- Target Project Repository (configured as $PROJECT-TARGET-REPO in target-repo-config.yaml)
- [OCI Distribution Spec](https://github.com/opencontainers/distribution-spec)
- [Distribution/Distribution](https://github.com/distribution/distribution)

### C. Assumptions
- IDP Builder continues using Cobra for CLI
- Kind remains the primary Kubernetes provider
- Go version compatibility maintained at 1.21+

---

*This plan demonstrates a realistic SF2.0 project for enhancing an existing codebase with new capabilities while maintaining compatibility and simplicity.*