# Kind Certificate Extraction Implementation Plan

## 🚨 CRITICAL EFFORT METADATA (FROM WAVE PLAN)
**Branch**: `idpbuidler-oci-go-cr/phase1/wave1/kind-certificate-extraction`
**Can Parallelize**: Yes
**Parallel With**: Effort 1.1.2 (Registry TLS Trust Integration)
**Size Estimate**: ~500 lines
**Dependencies**: None (foundational effort)

<!-- ⚠️ ORCHESTRATOR METADATA PLACEHOLDER - DO NOT REMOVE ⚠️ -->
<!-- The orchestrator will add infrastructure metadata below: -->
<!-- WORKING_DIRECTORY: /home/vscode/workspaces/idpbuilder-oci-go-cr/efforts/phase1/wave1/kind-certificate-extraction -->
<!-- BRANCH: idpbuidler-oci-go-cr/phase1/wave1/kind-certificate-extraction -->
<!-- REMOTE: origin (https://github.com/jessesanford/idpbuilder.git) -->
<!-- BASE_BRANCH: software-factory-2.0 -->
<!-- SW Engineers MUST read this metadata to navigate to the correct directory -->
<!-- END PLACEHOLDER -->

## Overview
- **Effort**: Extract and manage certificates from Kind/Gitea
- **Phase**: 1 (Certificate Infrastructure), Wave: 1 (Certificate Management Core)
- **Estimated Size**: ~500 lines
- **Implementation Time**: 1 day (Day 1 of MVP)

## Detailed Description
This effort implements the core functionality to extract self-signed certificates from the Gitea instance running inside the Kind cluster. The extracted certificates will be stored locally and made available for the registry client to establish secure TLS connections.

## File Structure
```
pkg/
└── certs/
    ├── extractor.go         # Main certificate extraction logic (~250 lines)
    ├── extractor_test.go     # Unit tests for extractor (~150 lines)
    ├── types.go              # Interface and type definitions (~50 lines)
    └── errors.go             # Custom error types and handling (~50 lines)
```

## Implementation Steps

### Step 1: Define Core Interfaces and Types (types.go)
```go
// pkg/certs/types.go
package certs

import (
    "context"
    "crypto/x509"
)

// KindCertExtractor defines the interface for extracting certificates from Kind clusters
type KindCertExtractor interface {
    // ExtractGiteaCert extracts the Gitea certificate from the Kind cluster
    ExtractGiteaCert(ctx context.Context) (*x509.Certificate, error)
    
    // GetClusterName returns the name of the Kind cluster
    GetClusterName() (string, error)
    
    // ValidateCertificate performs basic validation on the extracted certificate
    ValidateCertificate(cert *x509.Certificate) error
    
    // SaveCertificate saves the certificate to the local trust store
    SaveCertificate(cert *x509.Certificate, path string) error
}

// CertificateInfo contains metadata about an extracted certificate
type CertificateInfo struct {
    Subject    string
    Issuer     string
    NotBefore  time.Time
    NotAfter   time.Time
    IsCA       bool
    DNSNames   []string
}
```

### Step 2: Implement Error Handling (errors.go)
```go
// pkg/certs/errors.go
package certs

// Define custom error types for clear diagnostics
type ClusterNotFoundError struct {
    ClusterName string
}

type PodNotFoundError struct {
    PodName   string
    Namespace string
}

type CertificateInvalidError struct {
    Reason string
}

type PermissionError struct {
    Path   string
    Action string
}
```

### Step 3: Implement Certificate Extraction (extractor.go)
Key implementation points:
1. **Detect Kind cluster**: Use kubectl to check for Kind clusters
2. **Locate Gitea pod**: Find the Gitea pod in the cluster
3. **Extract certificate**: Copy cert from `/data/gitea/https/cert.pem` in pod
4. **Parse certificate**: Convert PEM to x509.Certificate
5. **Validate certificate**: Check expiry, subject, and basic validity
6. **Store locally**: Save to `~/.idpbuilder/certs/gitea.pem` with proper permissions

Implementation approach:
- Use `exec.Command` to run kubectl commands
- Handle errors gracefully with clear messages
- Support both docker and podman as container runtimes
- Create certificate directory if it doesn't exist
- Set appropriate file permissions (0600 for cert files)

### Step 4: Write Comprehensive Tests (extractor_test.go)
Test scenarios to cover:
1. **Happy path**: Successful certificate extraction
2. **Missing cluster**: Handle when Kind cluster doesn't exist
3. **Missing pod**: Handle when Gitea pod is not found
4. **Invalid certificate**: Handle malformed certificate data
5. **Permission issues**: Handle file system permission errors
6. **Certificate validation**: Test expiry and validity checks
7. **Mock kubectl**: Mock kubectl commands for unit testing

## Size Management
- **Estimated Lines**: ~500 lines total
- **Measurement Tool**: `/home/vscode/workspaces/idpbuilder-oci-go-cr/tools/line-counter.sh`
- **Check Frequency**: After each major component (every ~100 lines)
- **Split Threshold**: 700 lines (warning), 800 lines (stop immediately)
- **Current Status**: New effort, starting from 0 lines

## Test Requirements
- **Unit Tests**: 80% minimum coverage
- **Integration Tests**: Not required for this effort (Phase 2 responsibility)
- **Test Strategy**: 
  - Mock kubectl commands using interface injection
  - Test all error scenarios
  - Validate certificate parsing and storage
  - Use testify for assertions
  - Use golang/mock for kubectl mocking

### Test Files Expected:
- `pkg/certs/extractor_test.go` - Main unit tests
- Test fixtures for mock certificates
- Mock implementations for kubectl interface

## Pattern Compliance
- **Go Patterns**: 
  - Follow Go standard project layout
  - Use interfaces for testability
  - Return errors, don't panic
  - Use context for cancellation
  - Clear error messages with context
  
- **Security Requirements**:
  - Never expose private keys in logs
  - Set restrictive file permissions (0600)
  - Validate certificates before trusting
  - Clear audit trail for security operations
  - Support --insecure flag but log warnings

- **Code Style**:
  - gofmt compliant
  - golint clean
  - Meaningful variable names
  - Comprehensive comments for exported functions
  - Example usage in comments

## Dependencies
External packages required:
- `k8s.io/client-go` - For Kubernetes client operations
- `k8s.io/apimachinery` - For Kubernetes API types
- Standard library packages:
  - `crypto/x509` - Certificate handling
  - `encoding/pem` - PEM encoding/decoding
  - `os/exec` - For kubectl commands
  - `path/filepath` - Path manipulation
  - `io/ioutil` or `os` - File operations

## Success Criteria
- ✅ Successfully extracts certificate from Gitea pod
- ✅ Stores certificate in local trust store
- ✅ Clear error messages for all failure scenarios
- ✅ 80% test coverage achieved
- ✅ No security vulnerabilities
- ✅ Code stays under 500 lines
- ✅ All tests pass
- ✅ Integrates cleanly with Effort 1.1.2

## Integration Points
This effort provides the foundation for:
- **Effort 1.1.2**: Registry TLS Trust Integration will use the extracted certificates
- **Phase 2**: Build & Push operations will rely on the certificate infrastructure

The `KindCertExtractor` interface will be consumed by the trust store manager in Effort 1.1.2.

## Notes for SW Engineer
1. Start with the interface definitions to establish the contract
2. Implement error types early for consistent error handling
3. Use dependency injection for kubectl to enable testing
4. Consider using k8s.io/client-go instead of exec.Command for better error handling
5. Ensure all file operations check permissions first
6. Log security-relevant operations for audit trail
7. Keep the implementation focused - don't add features not in spec
8. Measure size frequently with the line-counter tool
9. Stop immediately if approaching 800 lines