<<<<<<< HEAD
# Image Builder Demo Guide

This document provides comprehensive instructions for running the Image Builder feature demonstrations, validating outputs, and understanding integration points.

## Overview

The Image Builder demo showcases the complete certificate extraction and trust management infrastructure implemented in the idpbuilder-oci-build-push project. It demonstrates:

1. **OCI Image Building** - Creating container images from directory contexts
2. **TLS Certificate Management** - Generating and managing self-signed certificates  
3. **Secure Registry Operations** - Pushing images with custom TLS verification
4. **Feature Flag Control** - Dynamic enable/disable of builder features

## Quick Start

```bash
# Enable the image builder feature
export ENABLE_IMAGE_BUILDER=true

# Run all demo scenarios
./demo-features.sh build-image --context ./test-data/sample-app --tag myapp:v1.0
./demo-features.sh generate-certs --namespace demo --secret-name demo-tls
./demo-features.sh push-with-tls --image myapp:v1.0 --registry localhost:5000
./demo-features.sh status
```

## Demo Scenarios

### Scenario 1: Build Simple OCI Image

**Purpose**: Demonstrate OCI image creation from directory contexts

**Command**:
```bash
./demo-features.sh build-image \
  --context ./test-data/sample-app \
  --tag myapp:v1.0 \
  --storage /tmp/oci-storage
```

**What it demonstrates**:
- Directory context scanning and packaging
- OCI layer creation and assembly
- Image metadata and labeling
- Local storage management
- Build result reporting with digest and size

**Expected Output**:
```
🎬 Demo: Image Builder Features
Timestamp: 2025-09-11 00:25:30
================================
📦 Building OCI Image
====================
Context: ./test-data/sample-app
Tag: myapp:v1.0
Storage: /tmp/oci-storage

🔨 Building image...
   - Creating context archive...
   - Generating OCI layers...
   - Adding metadata and labels...
   - Calculating digest...
✅ Image built successfully!
   Image ID: sha256:a1b2c3d4e5f6...
   Size: 1024 bytes
   Storage: /tmp/oci-storage/myapp_v1.0_

🔍 Verification:
   ✅ Image file exists
   ✅ Storage location accessible
   ✅ Build completed without errors
✅ Demo complete - ready for integration
```

**Validation Steps**:
1. ✅ Check that storage file exists at specified location
2. ✅ Verify image has proper SHA256 digest format
3. ✅ Confirm size is reasonable (>0 bytes)
4. ✅ Validate exit code is 0

### Scenario 2: Generate TLS Certificates

**Purpose**: Showcase TLS certificate generation for secure operations

**Command**:
```bash
./demo-features.sh generate-certs \
  --namespace demo \
  --secret-name demo-tls \
  --output ./test-data/certs
```

**What it demonstrates**:
- Self-signed certificate generation
- CA and server certificate creation
- Proper file permissions and security
- Kubernetes secret preparation
- Certificate validation and verification

**Expected Output**:
```
🔐 Generating TLS Certificates
==============================
Namespace: demo
Secret: demo-tls
Output: ./test-data/certs

🔨 Generating certificates...
   - Creating CA private key...
   - Generating CA certificate...
   - Creating server private key...
   - Generating server certificate...
✅ Certificates generated successfully!
   CA Certificate: ./test-data/certs/ca.crt
   Server Certificate: ./test-data/certs/server.crt
   Private Key: ./test-data/certs/server.key

🔍 Verification:
   ✅ All certificate files created
   ✅ Proper file permissions set
   ✅ Ready for Kubernetes secret creation

📋 Next steps:
   kubectl create secret tls demo-tls \
     --cert=./test-data/certs/server.crt \
     --key=./test-data/certs/server.key \
     --namespace=demo
✅ Demo complete - ready for integration
```

**Validation Steps**:
1. ✅ Verify ca.crt file exists with valid PEM format
2. ✅ Check server.crt contains certificate data
3. ✅ Confirm server.key has restricted permissions (600)
4. ✅ Validate kubectl command syntax is correct

### Scenario 3: Push to Registry with TLS

**Purpose**: Demonstrate secure registry operations with custom certificates

**Command**:
```bash
./demo-features.sh push-with-tls \
  --image myapp:v1.0 \
  --registry localhost:5000 \
  --cert-path ./test-data/certs/ca.crt
```

**What it demonstrates**:
- Image loading from local storage
- TLS certificate validation
- Secure registry connection establishment
- Layer upload and manifest pushing
- End-to-end secure container distribution

**Expected Output**:
```
🚀 Push to Registry with TLS
============================
Image: myapp:v1.0
Registry: localhost:5000
CA Certificate: ./test-data/certs/ca.crt

🔨 Pushing to registry...
   - Loading image from local storage...
   - Validating TLS certificate...
   - Establishing secure connection to localhost:5000...
   - Uploading layers...
   - Pushing manifest...
✅ Image pushed successfully!
   Registry: localhost:5000
   Tag: myapp:v1.0
   Digest: sha256:b2c3d4e5f6a7...
   TLS: Verified with custom CA

🔍 Verification:
   ✅ TLS certificate validation successful
   ✅ Secure connection established
   ✅ Image uploaded without errors
   ✅ Manifest pushed to registry

📋 Pull command:
   docker pull localhost:5000/myapp:v1.0
✅ Demo complete - ready for integration
```

**Validation Steps**:
1. ✅ Confirm TLS verification shows success
2. ✅ Check push digest is valid SHA256
3. ✅ Verify all upload steps completed
4. ✅ Validate pull command syntax

### Scenario 4: Feature Flag Toggle

**Purpose**: Display feature flag configuration and control

**Command**:
```bash
# Show current status
./demo-features.sh status

# Demonstrate enabling
export ENABLE_IMAGE_BUILDER=true
./demo-features.sh status

# Demonstrate disabling  
export ENABLE_IMAGE_BUILDER=false
./demo-features.sh status
```

**What it demonstrates**:
- Environment variable based feature control
- Feature status reporting
- Operation blocking when disabled
- Dynamic configuration changes

**Expected Output (Enabled)**:
```
🏁 Feature Flag Status
=====================
Current Environment:
   ENABLE_IMAGE_BUILDER: true

✅ Image Builder: ENABLED
   - OCI image building available
   - TLS certificate generation available
   - Registry operations enabled

📊 Feature Status Summary:
   Build Images: ✅ Available
   Generate Certs: ✅ Available
   Registry Push: ✅ Available

🔄 Toggle Examples:
   Enable:  export ENABLE_IMAGE_BUILDER=true
   Disable: export ENABLE_IMAGE_BUILDER=false
   Check:   echo $ENABLE_IMAGE_BUILDER
✅ Demo complete - ready for integration
```

**Expected Output (Disabled)**:
```
🏁 Feature Flag Status
=====================
Current Environment:
   ENABLE_IMAGE_BUILDER: <not set>

❌ Image Builder: DISABLED
   - All operations will be blocked
   - To enable: export ENABLE_IMAGE_BUILDER=true

📊 Feature Status Summary:
   Build Images: ❌ Blocked
   Generate Certs: ❌ Blocked
   Registry Push: ❌ Blocked
```

**Validation Steps**:
1. ✅ Verify status reflects environment variable correctly
2. ✅ Check that disabled state shows blocking
3. ✅ Confirm toggle examples are accurate
4. ✅ Validate environment detection logic

## Environment Setup

### Prerequisites

- Bash shell (version 4.0+)
- Write access to `/tmp` directory for image storage
- (Optional) `kubectl` for Kubernetes secret operations
- (Optional) `docker` for image pull verification

### Environment Variables

```bash
# Required for enabling features
export ENABLE_IMAGE_BUILDER=true

# Optional configuration
export DEMO_STORAGE_DIR="/tmp/oci-storage"
export DEMO_CERT_DIR="./test-data/certs" 
export DEMO_REGISTRY="localhost:5000"
```

### File Permissions

The demo will create files with appropriate permissions:
- Certificate files (`.crt`): 644 (readable)
- Private key files (`.key`): 600 (owner only)
- Image storage files: 644 (readable)
- Demo script: 755 (executable)

## Directory Structure

After running all demos, your directory will look like:

```
image-builder/
├── demo-features.sh          # Main demo script
├── DEMO.md                   # This documentation
├── test-data/                # Demo test files
│   ├── sample-app/           # Sample application context
│   │   ├── Dockerfile        # Container build instructions
│   │   ├── app.py           # Sample Python app
│   │   ├── config.json      # App configuration
│   │   ├── requirements.txt # Python dependencies
│   │   └── README.md        # App documentation
│   ├── configs/             # Configuration examples
│   │   ├── registry-config.yaml  # Registry settings
│   │   └── .gitignore       # Ignore patterns
│   └── certs/               # Generated certificates
│       ├── ca.crt           # Certificate Authority
│       ├── server.crt       # Server certificate
│       └── server.key       # Private key (600 perms)
└── /tmp/oci-storage/        # Image storage (external)
    └── myapp_v1.0_          # Built image file
```

## Integration with Other Efforts

### Wave-Level Integration

The Image Builder demos integrate with other Phase 2 Wave 1 efforts:

**Gitea Client Integration**:
- Use certificates generated by image builder for gitea registry access
- Share TLS configuration between efforts
- Demonstrate end-to-end registry operations

**Example Integration Command**:
```bash
# Generate certs with image builder
./demo-features.sh generate-certs --output ./shared-certs

# Use certs with gitea client  
cd ../gitea-client
./demo-features.sh registry-login --ca-cert ../image-builder/shared-certs/ca.crt
```

### Phase-Level Integration

For Phase 2 complete demonstration:

1. **Base Image Building**: Image builder provides container images for other components
2. **Security Infrastructure**: TLS certificates used across all Phase 2 efforts  
3. **Registry Operations**: Demonstrate complete push/pull workflows
4. **Feature Coordination**: Show unified feature flag management

## Troubleshooting Guide

### Common Issues

**Issue**: Demo script not executable
```bash
# Solution:
chmod +x demo-features.sh
```

**Issue**: Permission denied creating certificates
```bash
# Solution: Ensure write access to test-data directory
mkdir -p test-data/certs
chmod 755 test-data test-data/certs
```

**Issue**: Feature appears disabled
```bash
# Solution: Check environment variable
echo $ENABLE_IMAGE_BUILDER
export ENABLE_IMAGE_BUILDER=true
```

**Issue**: Storage directory not found
```bash
# Solution: Create storage directory
mkdir -p /tmp/oci-storage
# Or use custom location:
./demo-features.sh build-image --storage ./local-storage
=======
# Gitea Client Split-002 Demo Guide

This demo showcases the advanced registry operations implemented in Split-002, including push operations with progress tracking, repository listing with pagination, retry logic with exponential backoff, and repository deletion capabilities.

## Overview

Split-002 completes the Gitea registry client with operational features:

1. **Performance**: Chunked uploads for large images
2. **Reliability**: Retry logic for network issues  
3. **Observability**: Progress tracking and logging
4. **Testing**: Comprehensive stubs for unit tests

## Prerequisites

- Access to a Gitea registry (or use simulation mode)
- Split-001 authentication components available
- `demo-features.sh` script in the current directory

## Demo Scenarios

### 1. Push Image with Progress Tracking

Demonstrates multi-layer image push with real-time progress reporting.

```bash
./demo-features.sh push \
  --registry https://gitea.local:3000 \
  --image myapp:v1.0 \
  --source ./test-data/image.tar \
  --progress
```

**Features Demonstrated:**
- Chunked upload with configurable size (5MB default)
- Progress callbacks with layer and byte tracking
- SHA256 digest verification per layer
- Upload speed calculation
- Final manifest push

**Expected Output:**
- Layer upload progress bars
- Bytes transferred counter and upload speed
- SHA256 digest confirmations
- Total upload time and size metrics

### 2. Repository Listing with Pagination

Shows repository discovery with pagination support and metadata retrieval.

```bash
./demo-features.sh list-repos \
  --registry https://gitea.local:3000 \
  --page 1 \
  --per-page 10 \
  --format table
```

**Features Demonstrated:**
- Registry catalog API integration
- Pagination controls (page/per-page)
- Tag information retrieval
- Formatted table output
- Repository metadata display

**Expected Output:**
- Formatted table of repositories with tags and timestamps
- Pagination information (page X of Y)
- Total repository count
- Repository metadata (size, last push time)

### 3. Retry Logic with Exponential Backoff

Illustrates network resilience through configurable retry mechanisms.

```bash
./demo-features.sh push-with-retry \
  --registry https://gitea.local:3000 \
  --image stress-test:v1.0 \
  --simulate-failures 3 \
  --max-retries 5
```

**Features Demonstrated:**
- Exponential backoff algorithm
- Configurable retry attempts and delays
- Error classification (retryable vs permanent)
- Request timeout handling
- Retry attempt logging

**Expected Output:**
- Individual retry attempts with failure reasons
- Backoff delay calculations (1s, 2s, 4s, etc.)
- Final success after configured retries
- Total time including retry delays

### 4. Repository Deletion

Shows repository cleanup operations with confirmation safeguards.

```bash
./demo-features.sh delete \
  --registry https://gitea.local:3000 \
  --repo myapp \
  --confirm
```

**Features Demonstrated:**
- Repository existence verification
- Confirmation requirement for safety
- Manifest and blob cleanup
- Catalog update procedures
- Verification of successful deletion

**Expected Output:**
- Confirmation prompt handling
- Step-by-step deletion progress
- Cleanup verification results
- Success confirmation message

## Integration with Split-001

Split-002 seamlessly integrates with Split-001 authentication:

- **Shared Authentication**: Uses auth manager from Split-001
- **Common Registry Client**: Extends base client with operational features
- **TLS Configuration**: Leverages Split-001 certificate handling
- **Remote Options**: Compatible with Split-001 transport settings

### Integration Example

```bash
# Authenticate using Split-001 components
export GITEA_USERNAME="demo-user"
export GITEA_PASSWORD="demo-pass"

# Use Split-002 operations with Split-001 auth
./demo-features.sh push --registry https://gitea.local:3000 --image demo:latest
```

## Performance Tuning

### Chunk Size Optimization

Configure upload chunk size based on network conditions:

```bash
# For fast networks
export CHUNK_SIZE="10MB"

# For slower networks
export CHUNK_SIZE="1MB"

# For very slow networks
export CHUNK_SIZE="512KB"
```

### Retry Policy Tuning

Adjust retry behavior for different environments:

```bash
# Aggressive retries for unstable networks
export MAX_RETRIES=10
export INITIAL_DELAY="500ms"
export BACKOFF_FACTOR=1.5

# Conservative retries for stable networks
export MAX_RETRIES=3
export INITIAL_DELAY="1s"
export BACKOFF_FACTOR=2.0
```

## Troubleshooting

### Common Issues

#### 1. Authentication Failures

```bash
# Verify Split-001 auth setup
echo "Checking authentication..."
curl -H "Authorization: Bearer $GITEA_TOKEN" https://gitea.local:3000/v2/

# Solution: Ensure Split-001 auth components are configured
export GITEA_TOKEN=$(get-auth-token.sh)
```

#### 2. Network Timeouts

```bash
# Increase timeout values
export REQUEST_TIMEOUT="60s"
export TLS_HANDSHAKE_TIMEOUT="30s"

# Enable retry logic
./demo-features.sh push-with-retry --max-retries 10
```

#### 3. Large Image Upload Failures

```bash
# Use smaller chunk sizes
export CHUNK_SIZE="1MB"

# Enable progress tracking
./demo-features.sh push --progress
```

#### 4. TLS Certificate Issues

```bash
# Use insecure mode for testing (NOT for production)
export TLS_INSECURE=true

# Or add certificate to trust store (recommended)
cp gitea-cert.pem /usr/local/share/ca-certificates/
update-ca-certificates
>>>>>>> gitea-split-002/idpbuilder-oci-build-push/phase2/wave1/gitea-client-split-002
```

### Debug Mode

<<<<<<< HEAD
Enable verbose output by setting:
```bash
export DEMO_DEBUG=true
./demo-features.sh <scenario>
```

### Cleanup

To clean up demo artifacts:
```bash
# Remove generated files
rm -rf test-data/certs/*.crt test-data/certs/*.key
rm -rf /tmp/oci-storage/*

# Reset environment
unset ENABLE_IMAGE_BUILDER
unset DEMO_DEBUG
```

## Success Metrics

The demo is considered successful when:

- ✅ All 4 scenarios execute without errors (exit code 0)
- ✅ Expected output files are created in correct locations
- ✅ File permissions are set appropriately for security
- ✅ Feature flag toggling works as expected
- ✅ Integration hooks are properly exported
- ✅ All validation steps pass

## Security Considerations

### Certificate Security

- Private keys are generated with restricted permissions (600)
- Demo certificates are clearly marked as test/demo use only
- Production deployments should use proper CA-signed certificates
- All certificate operations are logged with timestamps

### Feature Flag Security

- Feature flags provide defense-in-depth security
- Disabled features are completely blocked, not just warned
- Feature state is clearly visible in all operations
- Environment-based control allows runtime configuration

### Demo Limitations

- Mock certificates are used for demonstration only
- Actual registry operations are simulated for safety
- Real deployments require proper secret management
- TLS verification is demonstrated but not cryptographically validated

## Next Steps

After running the demos:

1. **Review the code**: Examine `pkg/build/` for implementation details
2. **Integration testing**: Run with actual Gitea registry
3. **Production setup**: Replace demo certificates with real ones
4. **Monitoring**: Add metrics and logging for production use
5. **Documentation**: Update operational runbooks

## Support

For issues with the demo:

1. Check the troubleshooting section above
2. Verify all prerequisites are met
3. Review the implementation plan in `IMPLEMENTATION-PLAN.md`
4. Check the work log in `work-log.md` for recent changes

---

**Demo Version**: 1.0  
**Compatible with**: Image Builder implementation Phase 2 Wave 1  
**Last Updated**: 2025-09-11  
**Integration Ready**: ✅
=======
Enable detailed logging for troubleshooting:

```bash
export DEBUG_MODE=true
export LOG_LEVEL=debug

./demo-features.sh push --registry https://gitea.local:3000 --image debug:latest
```

## Testing Configuration

### Mock Registry Setup

For testing without a real Gitea instance:

```bash
# Use test stubs from stubs.go
export DEMO_MODE="simulation"
export MOCK_REGISTRY=true

./demo-features.sh push --registry mock://localhost
```

### Performance Testing

Benchmark upload performance:

```bash
# Create test image
dd if=/dev/zero of=./test-data/large-image.tar bs=1M count=100

# Test with different chunk sizes
for chunk in 1MB 5MB 10MB; do
    export CHUNK_SIZE=$chunk
    time ./demo-features.sh push --source ./test-data/large-image.tar
done
```

## Security Considerations

### Authentication Security

- Always use HTTPS in production
- Store credentials securely (environment variables, not files)
- Rotate authentication tokens regularly
- Use least-privilege access controls

### TLS Security

- Verify server certificates in production
- Use `--insecure` only for development/testing
- Keep CA certificates updated
- Monitor for certificate expiration

### Network Security

- Use private networks when possible
- Implement network policies for registry access
- Monitor for suspicious upload patterns
- Log all registry operations for audit

## Advanced Usage

### Batch Operations

Process multiple images:

```bash
#!/bin/bash
for image in $(cat image-list.txt); do
    ./demo-features.sh push --image "$image" --progress
done
```

### Integration with CI/CD

```yaml
# Example GitHub Actions step
- name: Push to Gitea Registry
  run: |
    ./demo-features.sh push \
      --registry ${{ secrets.GITEA_REGISTRY_URL }} \
      --image ${{ github.repository }}:${{ github.sha }} \
      --progress
```

### Monitoring Integration

```bash
# Export metrics for monitoring
./demo-features.sh push --image monitor:latest --progress | \
  grep -E "(uploaded|complete)" | \
  curl -X POST http://metrics-collector/api/registry-ops
```

## API Reference

The demo script supports these environment variables:

- `REGISTRY_URL`: Default registry URL
- `DEMO_MODE`: Set to "simulation" for testing
- `CHUNK_SIZE`: Upload chunk size
- `MAX_RETRIES`: Maximum retry attempts
- `INITIAL_DELAY`: Initial retry delay
- `BACKOFF_FACTOR`: Exponential backoff multiplier
- `DEBUG_MODE`: Enable debug logging
- `TLS_INSECURE`: Skip TLS verification (testing only)

## Support

For issues or questions:

1. Check the troubleshooting section above
2. Review logs with debug mode enabled
3. Verify Split-001 integration is working
4. Test with simulation mode to isolate issues

## Version Compatibility

- Requires Split-001 (authentication) to be implemented
- Compatible with go-containerregistry v0.19.0+
- Tested with Gitea 1.20+
- Supports Docker Registry API v2
>>>>>>> gitea-split-002/idpbuilder-oci-build-push/phase2/wave1/gitea-client-split-002
