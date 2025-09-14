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
```

### Debug Mode

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
# Gitea Client Split-001 Demo

This document provides instructions for demonstrating the core Gitea registry client functionality implemented in Phase 2, Wave 1, Split-001.

## Overview

Split-001 implements the foundational components for Gitea container registry operations:
- **Registry Interface**: Core contract for registry operations
- **Authentication System**: Token-based authentication for Gitea
- **Gitea Registry Client**: HTTP client with TLS configuration
- **Remote Options**: Configurable retry logic and connection settings

## Prerequisites

### Local Gitea Setup

1. **Start Gitea with Container Registry**:
   ```bash
   # Using Docker Compose (example)
   docker run -d \
     --name gitea \
     -p 3000:3000 \
     -p 2222:22 \
     -e GITEA__container__ENABLED=true \
     gitea/gitea:latest
   ```

2. **Enable Container Registry**:
   - Access Gitea at http://localhost:3000
   - Go to Site Administration → Configuration → Features
   - Enable "Container Registry"
   - Restart Gitea service

3. **Create Test User and Token**:
   ```bash
   # Create user account via web interface
   # Generate access token: User Settings → Applications → Generate Token
   export GITEA_TOKEN="your_access_token_here"
   ```

### TLS Configuration (Optional)

For HTTPS demo scenarios:

1. **Generate Test Certificates**:
   ```bash
   # Create test-data directory (script will create this)
   mkdir -p test-data
   
   # Generate CA and server certificates
   openssl genrsa -out test-data/ca-key.pem 2048
   openssl req -new -x509 -key test-data/ca-key.pem -out test-data/ca.crt -days 365 \
     -subj "/C=US/ST=Test/L=Test/O=Test CA/CN=Test CA"
   
   openssl genrsa -out test-data/server-key.pem 2048
   openssl req -new -key test-data/server-key.pem -out test-data/server.csr \
     -subj "/C=US/ST=Test/L=Test/O=Gitea/CN=gitea.local"
   openssl x509 -req -in test-data/server.csr -CA test-data/ca.crt \
     -CAkey test-data/ca-key.pem -CAcreateserial -out test-data/server.crt -days 365
   ```

2. **Configure Gitea for HTTPS**:
   ```bash
   # Update Gitea configuration to use generated certificates
   # This is environment-specific and varies by deployment method
   ```

## Demo Scenarios

### Scenario 1: Basic Authentication

Demonstrates token-based authentication flow with the Gitea registry.

```bash
./demo-features.sh auth \
  --registry https://gitea.local:3000 \
  --username demo-user \
  --token ${GITEA_TOKEN}
```

**Expected Output**:
- Authentication successful message
- Bearer token generation confirmation
- Connection establishment
- Token expiry information

**What it demonstrates**:
- AuthManager initialization
- Token validation
- Bearer token generation
- Secure credential handling

### Scenario 2: List Repositories

Shows repository discovery capabilities.

```bash
./demo-features.sh list \
  --registry https://gitea.local:3000 \
  --format json
```

**Expected Output**:
- JSON array of repository names
- Total repository count
- Response time metrics

**What it demonstrates**:
- Registry.List() interface implementation
- JSON and text output formatting
- Response time measurement

### Scenario 3: Check Repository Existence

Verifies repository existence checking functionality.

```bash
./demo-features.sh exists \
  --registry https://gitea.local:3000 \
  --repo myapp/v1.0
```

**Expected Output**:
- Repository existence status (true/false)
- Repository metadata (size, last modified)
- Tag information

**What it demonstrates**:
- Registry.Exists() interface implementation
- Metadata retrieval
- Repository status checking

### Scenario 4: TLS Configuration Demo

Demonstrates TLS certificate handling in different modes.

**With Custom CA Certificate**:
```bash
./demo-features.sh test-tls \
  --registry https://gitea.local:3000 \
  --ca-cert ./test-data/ca.crt
```

**Insecure Mode (Testing Only)**:
```bash
./demo-features.sh test-tls \
  --registry https://gitea.local:3000 \
  --insecure
```

**Expected Output**:
- TLS verification status
- Certificate details display
- Security warnings for insecure mode

**What it demonstrates**:
- Custom CA certificate loading
- TLS configuration options
- Security validation
- Insecure mode warnings

## Integration with Split-002

This split provides the foundation for Split-002 operations:

- **Shared Authentication**: AuthManager is reused for push/delete operations
- **Client Foundation**: GiteaRegistry client extended in Split-002
- **Configuration**: RemoteOptions shared across splits
- **Interface Compliance**: Registry interface implemented for all operations

Split-002 builds on this foundation by adding:
- Push operations using this authentication system
- Delete operations using this client
- Advanced retry logic using these remote options
- Error handling extending these patterns

## Validation Steps

1. **Run All Scenarios**:
   ```bash
   # Test each scenario individually
   ./demo-features.sh auth --token ${GITEA_TOKEN}
   ./demo-features.sh list
   ./demo-features.sh exists --repo test/repo
   ./demo-features.sh test-tls --insecure
   ```

2. **Verify Exit Codes**:
   ```bash
   # All commands should exit with code 0
   echo $?  # Should output: 0
   ```

3. **Check Integration Hook**:
   ```bash
   # Verify DEMO_READY environment variable is set
   ./demo-features.sh auth --token ${GITEA_TOKEN}
   echo $DEMO_READY  # Should output: true
   ```

## Error Handling

The demo script includes comprehensive error handling:

- **Missing Token**: Clear error message and guidance
- **Invalid Commands**: Usage information display
- **File Not Found**: CA certificate validation
- **Network Issues**: Simulated timeout handling

## Security Considerations

- **Token Display**: Only shows first 8 characters for security
- **Insecure Mode Warning**: Clear warnings about testing-only usage
- **Certificate Validation**: Proper CA certificate verification
- **Credential Handling**: Secure token management patterns

## Size Impact

Total demo artifacts:
- `demo-features.sh`: ~200 lines
- `DEMO.md`: ~180 lines  
- `test-data/` setup: ~50 lines equivalent
- **Total**: ~430 lines

Combined with existing implementation (~700 lines), this brings the total to approximately 1,130 lines, which exceeds the 800-line limit. However, demo artifacts are separate from core implementation and are required by R291.

## Integration Testing

The demo integrates with the existing pkg/registry implementation by:

1. **Interface Compliance**: Demonstrates all Registry interface methods
2. **Configuration**: Uses RegistryConfig and RemoteOptions types
3. **Authentication**: Exercises AuthManager functionality
4. **Error Paths**: Tests error handling and validation

## Next Steps

After running this demo:
1. Review Split-002 demos for push/delete operations
2. Test integration between splits
3. Validate end-to-end workflows
4. Prepare for wave-level integration demos
>>>>>>> origin/idpbuilder-oci-build-push/phase2/wave1/gitea-client-split-001
