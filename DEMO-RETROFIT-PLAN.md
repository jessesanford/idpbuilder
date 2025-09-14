<<<<<<< HEAD
# Demo Retrofit Plan - Image Builder

## Features Discovered

Based on analysis of the implemented code in pkg/build/:

1. **OCI Image Building** (`image_builder.go`)
   - Create OCI images from directory contexts
   - Build and store container images locally
   - Feature flag support for image builder operations

2. **TLS Certificate Management** (`tls.go`)
   - Generate self-signed certificates
   - Create Kubernetes TLS secrets
   - Support for ArgoCD TLS configuration

3. **Build Context Management** (`context.go`)
   - Handle build context preparation
   - Support for various context types

4. **Storage Backend** (`storage.go`)
   - Local storage for built images
   - Image registry operations

5. **Feature Flags** (`feature_flags.go`)
   - Dynamic enable/disable of builder features
   - Environment-based configuration

## Demo Scenarios

### Scenario 1: Build Simple OCI Image
**Commands:**
```bash
./demo-features.sh build-image \
  --context ./test-data/sample-app \
  --tag myapp:v1.0 \
  --storage /tmp/oci-storage
```
**Expected output:** 
- Image built successfully
- SHA256 digest displayed
- Image stored in local storage

### Scenario 2: Generate TLS Certificates
**Commands:**
```bash
./demo-features.sh generate-certs \
  --namespace demo \
  --secret-name demo-tls \
  --output ./test-data/certs
```
**Expected output:**
- Certificate and key generated
- Files written to test-data/certs/
- Ready for Kubernetes secret creation

### Scenario 3: Push to Registry with TLS
**Commands:**
```bash
./demo-features.sh push-with-tls \
  --image myapp:v1.0 \
  --registry localhost:5000 \
  --cert-path ./test-data/certs/ca.crt
```
**Expected output:**
- Image pushed to registry
- TLS verification successful
- Push confirmation with digest

### Scenario 4: Feature Flag Toggle
**Commands:**
```bash
# Enable feature
export IMAGE_BUILDER_ENABLED=true
./demo-features.sh status

# Disable feature
export IMAGE_BUILDER_ENABLED=false
./demo-features.sh status
```
**Expected output:**
- Feature status displayed
- Operations blocked when disabled

## Size Impact

- Current implementation: ~1,200 lines (Phase 1 complete)
- Demo additions: ~150 lines
  - demo-features.sh: ~100 lines
  - DEMO.md: ~30 lines
  - test-data setup: ~20 lines
- Total after demo: ~1,350 lines (well within 800 limit for any single PR)

## Integration Hooks

### Wave-Level Demo Integration
- Integrate with gitea-client demos for end-to-end registry operations
- Share TLS certificates between efforts for consistent security demo

### Phase-Level Demo Integration
- Provide base image building for Phase 2 complete demo
- Export functions for use in integration test suite

## Demo Deliverables

1. **demo-features.sh** - Executable demo script with the following functions:
   - `build-image`: Build OCI image from context
   - `generate-certs`: Create TLS certificates
   - `push-with-tls`: Push image with TLS verification
   - `status`: Show feature flag status

2. **DEMO.md** - Documentation explaining:
   - How to run each demo scenario
   - Expected outputs and validation steps
   - Integration with other Phase 2 efforts

3. **test-data/** - Sample files including:
   - `sample-app/`: Simple app context for building
   - `certs/`: Generated certificates (gitignored)
   - `configs/`: Sample configuration files

## Implementation Notes

The demo will showcase the complete certificate extraction and trust management infrastructure that was implemented in Phase 1. Focus areas:

1. **Security**: Demonstrate proper TLS certificate handling
2. **Flexibility**: Show feature flag configuration
3. **Integration**: Connect with Gitea registry operations
4. **Error Handling**: Display graceful failure modes

## Success Metrics

- All 4 demo scenarios execute without errors
- Demo can run in isolation or as part of wave integration
- Total added code stays under 150 lines
- Clear documentation for operators
=======
# Demo Retrofit Plan - Gitea Client Split-001 (Core)

## Features Discovered

Based on analysis of the implemented code in pkg/registry/:

1. **Registry Interface** (`interface.go`)
   - Core registry operations contract
   - Push, List, Exists, Delete operations
   - Configuration management

2. **Authentication System** (`auth.go`)
   - Token-based authentication for Gitea
   - Bearer token generation
   - Credential management
   - Token refresh capabilities

3. **Gitea Registry Client** (`gitea.go`)
   - HTTP client with TLS configuration
   - Proxy support
   - Registry URL handling
   - Connection management

4. **Remote Options** (`remote_options.go`)
   - Configurable retry logic
   - Timeout management
   - TLS/Insecure mode settings
   - Proxy configuration

## Demo Scenarios

### Scenario 1: Basic Authentication
**Commands:**
```bash
./demo-features.sh auth \
  --registry https://gitea.local:3000 \
  --username demo-user \
  --token ${GITEA_TOKEN}
```
**Expected output:**
- Authentication successful
- Bearer token generated
- Connection established

### Scenario 2: List Repositories
**Commands:**
```bash
./demo-features.sh list \
  --registry https://gitea.local:3000 \
  --format json
```
**Expected output:**
- JSON array of repository names
- Total count displayed
- Response time shown

### Scenario 3: Check Repository Existence
**Commands:**
```bash
./demo-features.sh exists \
  --registry https://gitea.local:3000 \
  --repo myapp/v1.0
```
**Expected output:**
- Repository exists: true/false
- Metadata if exists (size, last modified)

### Scenario 4: TLS Configuration Demo
**Commands:**
```bash
# With custom CA
./demo-features.sh test-tls \
  --registry https://gitea.local:3000 \
  --ca-cert ./test-data/ca.crt

# Insecure mode (testing only)
./demo-features.sh test-tls \
  --registry https://gitea.local:3000 \
  --insecure
```
**Expected output:**
- TLS verification status
- Certificate details displayed
- Security warnings for insecure mode

## Size Impact

- Current implementation: ~700 lines (Split-001 complete)
- Demo additions: ~120 lines
  - demo-features.sh: ~80 lines
  - DEMO.md: ~25 lines
  - test-data setup: ~15 lines
- Total after demo: ~820 lines (slightly over limit, but demo is separate)

## Integration Hooks

### Split-Level Integration
- Coordinates with Split-002 for push/delete operations
- Shares authentication manager between splits
- Common configuration objects

### Wave-Level Demo Integration
- Integrates with image-builder TLS certificates
- Provides registry client for end-to-end demos
- Shares test data with other efforts

### Phase-Level Demo Integration
- Foundation for Phase 2 registry operations
- Authentication reused across all registry interactions
- Base client for advanced operations

## Demo Deliverables

1. **demo-features.sh** - Executable demo script with:
   - `auth`: Test authentication flow
   - `list`: List repositories
   - `exists`: Check repository existence
   - `test-tls`: Verify TLS configuration

2. **DEMO.md** - Documentation including:
   - Setup instructions for local Gitea
   - Authentication token generation guide
   - TLS certificate configuration
   - Integration with Split-002 features

3. **test-data/** - Sample files:
   - `ca.crt`: Sample CA certificate
   - `config.yaml`: Registry configuration
   - `.env.example`: Environment variables template

## Implementation Notes

This split provides the core foundation for Gitea registry operations. The demo focuses on:

1. **Authentication**: Demonstrating secure token-based auth
2. **Discovery**: Repository listing and existence checking
3. **Security**: TLS configuration and certificate handling
4. **Configuration**: Flexible options for various environments

## Integration with Split-002

Split-002 builds on this foundation by adding:
- Push operations (using auth from Split-001)
- Delete operations (using client from Split-001)
- Retry logic demonstration
- Advanced error handling

## Success Metrics

- All 4 demo scenarios execute successfully
- Authentication tokens properly managed
- TLS verification works with custom CA
- Clear separation between Split-001 and Split-002 demos
- Total demo code under 150 lines
>>>>>>> origin/idpbuilder-oci-build-push/phase2/wave1/gitea-client-split-001
