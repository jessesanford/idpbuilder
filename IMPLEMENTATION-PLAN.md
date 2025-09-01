# Registry TLS Trust Integration Implementation Plan

## 🚨 CRITICAL EFFORT METADATA (FROM WAVE PLAN)
**Branch**: `idpbuidler-oci-go-cr/phase1/wave1/registry-tls-trust-integration`
**Can Parallelize**: Yes (foundational effort)
**Parallel With**: E1.1.1 (Kind Certificate Extraction)
**Size Estimate**: ~600 lines
**Dependencies**: None (foundational effort - can run parallel with E1.1.1)

## Overview
- **Effort**: E1.1.2 - Registry TLS Trust Integration
- **Phase**: 1, Wave: 1
- **Estimated Size**: ~600 lines
- **Implementation Time**: 1 day (Day 2 of implementation)
- **Purpose**: Configure TLS trust for registry operations with go-containerregistry

## Context & Goals

### Primary Objective
Implement a robust TLS trust management system that enables go-containerregistry (ggcr) to work with Gitea's self-signed certificates. This is a critical component of the MVP that solves the certificate problem preventing reliable OCI image operations.

### Key Success Criteria
- ✅ Load custom CA certificates into x509.CertPool
- ✅ Configure ggcr remote transport with proper TLS settings
- ✅ Support certificate rotation without restart
- ✅ Provide --insecure override for testing
- ✅ Clear error messages for certificate issues
- ✅ Zero certificate errors during normal operation

## File Structure
```
efforts/phase1/wave1/registry-tls-trust-integration/
└── pkg/
    └── certs/
        ├── trust.go              # Main TrustStoreManager implementation (~250 lines)
        ├── transport.go          # GGCR transport configuration (~150 lines)
        ├── trust_store.go        # Trust store persistence (~100 lines)
        └── trust_test.go         # Comprehensive test suite (~100 lines)
```

## Implementation Steps

### Step 1: Core Trust Store Manager (250 lines)
**File**: `pkg/certs/trust.go`

Create the main TrustStoreManager interface and implementation:

```go
// pkg/certs/trust.go
package certs

import (
    "crypto/tls"
    "crypto/x509"
    "fmt"
    "io/ioutil"
    "os"
    "path/filepath"
    "sync"
)

// TrustStoreManager manages trusted certificates for registry operations
type TrustStoreManager interface {
    // AddCertificate adds a certificate for a specific registry
    AddCertificate(registry string, cert *x509.Certificate) error
    
    // RemoveCertificate removes the certificate for a registry
    RemoveCertificate(registry string) error
    
    // SetInsecureRegistry marks a registry as insecure (skip TLS verification)
    SetInsecureRegistry(registry string, insecure bool) error
    
    // GetTrustedCerts returns all trusted certificates for a registry
    GetTrustedCerts(registry string) ([]*x509.Certificate, error)
    
    // GetCertPool returns a configured cert pool for a registry
    GetCertPool(registry string) (*x509.CertPool, error)
    
    // IsInsecure checks if a registry is marked as insecure
    IsInsecure(registry string) bool
}

type trustStoreManager struct {
    certsDir        string
    trustedCerts    map[string][]*x509.Certificate
    insecureRegistries map[string]bool
    mu              sync.RWMutex
}
```

Key implementation details:
- Thread-safe operations with sync.RWMutex
- Persistent storage at `~/.idpbuilder/certs/`
- Certificate files named by registry hostname
- Support for multiple certificates per registry
- Automatic reload of certificates on access

### Step 2: GGCR Transport Configuration (150 lines)
**File**: `pkg/certs/transport.go`

Configure go-containerregistry remote options with TLS settings:

```go
// pkg/certs/transport.go
package certs

import (
    "crypto/tls"
    "net/http"
    "github.com/google/go-containerregistry/pkg/v1/remote"
)

// ConfigureTransport creates remote.Option with proper TLS configuration
func (m *trustStoreManager) ConfigureTransport(registry string) (remote.Option, error) {
    if m.IsInsecure(registry) {
        return remote.WithTransport(&http.Transport{
            TLSClientConfig: &tls.Config{
                InsecureSkipVerify: true,
            },
        }), nil
    }
    
    certPool, err := m.GetCertPool(registry)
    if err != nil {
        return nil, fmt.Errorf("failed to get cert pool: %w", err)
    }
    
    return remote.WithTransport(&http.Transport{
        TLSClientConfig: &tls.Config{
            RootCAs: certPool,
        },
    }), nil
}
```

Key features:
- Seamless integration with go-containerregistry
- Support for --insecure flag via InsecureSkipVerify
- Custom CA pool configuration
- Error handling with clear messages

### Step 3: Trust Store Persistence (100 lines)
**File**: `pkg/certs/trust_store.go`

Implement certificate persistence and loading:

```go
// pkg/certs/trust_store.go
package certs

// LoadFromDisk loads certificates from disk into memory
func (m *trustStoreManager) LoadFromDisk() error {
    // Implementation details:
    // 1. Scan ~/.idpbuilder/certs/ directory
    // 2. Load all .pem files
    // 3. Parse certificates
    // 4. Populate in-memory maps
    // 5. Handle file permissions errors gracefully
}

// SaveToDisk persists a certificate to disk
func (m *trustStoreManager) SaveToDisk(registry string, cert *x509.Certificate) error {
    // Implementation details:
    // 1. Create directory if not exists
    // 2. Write certificate to PEM file
    // 3. Set appropriate file permissions (0600)
    // 4. Handle write errors with clear messages
}
```

### Step 4: Comprehensive Test Suite (100 lines)
**File**: `pkg/certs/trust_test.go`

Test all functionality with high coverage:

```go
// pkg/certs/trust_test.go
package certs

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestTrustStoreManager(t *testing.T) {
    t.Run("AddCertificate", func(t *testing.T) {
        // Test adding certificates
    })
    
    t.Run("LoadFromPEM", func(t *testing.T) {
        // Test loading PEM files
    })
    
    t.Run("CertRotation", func(t *testing.T) {
        // Test certificate rotation support
    })
    
    t.Run("InsecureMode", func(t *testing.T) {
        // Test --insecure flag behavior
    })
    
    t.Run("TransportConfiguration", func(t *testing.T) {
        // Test GGCR transport setup
    })
}
```

## Integration Points

### With E1.1.1 (Kind Certificate Extraction)
While this effort can run in parallel with E1.1.1, they will integrate:
- E1.1.1 extracts certificates from Kind/Gitea
- E1.1.2 (this effort) provides the trust store to use them
- Both efforts share the `pkg/certs` package namespace

### With Phase 2 Efforts
This trust store will be consumed by:
- E2.1.2: Gitea Registry Client (uses ConfigureTransport)
- All registry operations will use the trust configuration

## Size Management
- **Estimated Lines**: 600
- **Current Structure**: 250 + 150 + 100 + 100 = 600 lines
- **Measurement Tool**: ${PROJECT_ROOT}/tools/line-counter.sh
- **Check Frequency**: After each major component
- **Split Threshold**: 700 lines (warning), 800 lines (stop)

## Test Requirements
- **Unit Tests**: 80% coverage minimum
- **Integration Tests**: Test with actual certificates
- **E2E Tests**: Not required in Phase 1
- **Test Files**: 
  - `pkg/certs/trust_test.go` - Main test suite
  - Test fixtures in `testdata/` if needed

### Test Coverage Areas
1. **Certificate Management**
   - Adding/removing certificates
   - Loading from PEM files
   - Certificate validation
   - Error handling

2. **Transport Configuration**
   - TLS configuration with custom CA
   - Insecure mode handling
   - Connection testing

3. **Persistence**
   - Save/load from disk
   - File permission handling
   - Directory creation

4. **Certificate Rotation**
   - Reload certificates without restart
   - Handle expired certificates
   - Update notifications

## Pattern Compliance
- **Go Best Practices**: 
  - Interface-driven design (TrustStoreManager)
  - Error wrapping with context
  - Thread-safe operations
  - Clear package boundaries

- **Security Requirements**:
  - Never silently ignore certificate errors
  - Require explicit --insecure flag
  - Log all security decisions
  - Secure file permissions (0600)

- **Performance Targets**:
  - Certificate loading < 100ms
  - In-memory caching for performance
  - Lazy loading on first use

## Error Handling Strategy

### Certificate Errors
```go
// Clear, actionable error messages
return fmt.Errorf("certificate verification failed for %s: %w\n" +
    "To fix this issue:\n" +
    "1. Ensure the certificate is valid\n" +
    "2. Check certificate expiry: %s\n" +
    "3. Or use --insecure flag for testing", 
    registry, err, cert.NotAfter)
```

### Permission Errors
```go
// Handle permission issues gracefully
if os.IsPermission(err) {
    return fmt.Errorf("permission denied accessing certificate store at %s\n" +
        "Please check file permissions or run with appropriate privileges",
        m.certsDir)
}
```

## CLI Integration

The trust store will integrate with CLI commands:

```bash
# Normal operation with certificates
idpbuilder-oci push myimage:latest

# Testing with --insecure flag
idpbuilder-oci push --insecure myimage:latest

# Certificate management (future)
idpbuilder-oci cert add gitea.local ~/certs/gitea.pem
idpbuilder-oci cert list
idpbuilder-oci cert remove gitea.local
```

## Success Metrics
- ✅ Zero certificate errors during normal operation
- ✅ Clear error messages when certificate issues occur
- ✅ Successful integration with go-containerregistry
- ✅ 80%+ test coverage
- ✅ Implementation within 600 line budget

## Dependencies & Imports
```go
import (
    // Standard library
    "crypto/tls"
    "crypto/x509"
    "encoding/pem"
    "fmt"
    "io/ioutil"
    "os"
    "path/filepath"
    "sync"
    
    // External dependencies
    "github.com/google/go-containerregistry/pkg/v1/remote"
    
    // Testing
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)
```

## Notes for SW Engineer

### Priority Order
1. Implement TrustStoreManager interface first
2. Add transport configuration for GGCR
3. Implement persistence layer
4. Write comprehensive tests
5. Integrate with CLI commands (coordination with other efforts)

### Key Considerations
- This is a foundational component - quality is critical
- Focus on clear error messages for debugging
- Certificate rotation must work without restarts
- Security decisions must be explicit and logged
- Coordinate with E1.1.1 on shared package structure

### Testing Approach
- Use self-signed certificates in test fixtures
- Test both happy path and error cases
- Verify thread safety with concurrent tests
- Mock filesystem for permission testing

## Completion Checklist
- [ ] TrustStoreManager interface implemented
- [ ] GGCR transport configuration working
- [ ] Certificate persistence functional
- [ ] 80%+ test coverage achieved
- [ ] Error messages are clear and actionable
- [ ] --insecure flag properly implemented
- [ ] Code review passed
- [ ] Size limit compliance (<800 lines)