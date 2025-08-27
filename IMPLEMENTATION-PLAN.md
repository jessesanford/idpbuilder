# E3.1.2 Certificate Bundle Loader Implementation Plan

## 🚨 CRITICAL EFFORT METADATA (FROM WAVE PLAN)
**Branch**: idpbuidler-oci-mgmt/phase3/wave1/E3.1.2-bundle-loader  
**Can Parallelize**: Yes (after E3.1.1)  
**Parallel With**: [E3.1.3, E3.1.4, E3.1.5]  
**Size Estimate**: 700 lines  
**Dependencies**: [E3.1.1]  
**Working Directory**: /home/vscode/workspaces/idpbuilder-oci-mgmt/efforts/phase3/wave1/E3.1.2-bundle-loader

## Overview
- **Effort**: Certificate Bundle Loader - Multi-format certificate loading with auto-detection
- **Phase**: 3, Wave: 1
- **Estimated Size**: 700 lines total
- **Implementation Time**: 6 hours
- **Purpose**: Implement comprehensive certificate loading supporting PEM, DER, PKCS7, and PKCS12 formats with automatic format detection

## Dependency Context from E3.1.1
Based on the interfaces defined in E3.1.1-certificate-contracts:
- **CertificateService Interface**: We implement the LoadCertificateBundle method
- **CertBundle Structure**: Use the defined structure for returning loaded certificates
- **CertFormat Types**: Support all four defined formats (PEM, DER, PKCS7, PKCS12)
- **CertificateError**: Use the standardized error type for certificate-related errors

## File Structure
```
pkg/oci/certificates/
├── loader.go (250 lines)
│   └── MultiFormatLoader struct implementing certificate loading
│   └── LoadPEM(), LoadDER(), LoadPKCS7(), LoadPKCS12() methods
│   └── Format detection logic with magic bytes
├── parser.go (200 lines)
│   └── Certificate parsing utilities
│   └── Chain building and validation logic
│   └── Format conversion helpers
├── formats.go (150 lines)
│   └── Format-specific handlers and validators
│   └── Conversion utilities between formats
│   └── Certificate extraction from bundles
└── loader_test.go (100 lines)
    └── Comprehensive format tests
    └── Invalid certificate handling tests
    └── Performance benchmarks
```

## Implementation Steps

### Step 1: Create Package Structure and Base Loader (250 lines)
**File**: `pkg/oci/certificates/loader.go`

```go
package certificates

import (
    "context"
    "crypto/x509"
    "encoding/pem"
    "errors"
    "fmt"
    "io/ioutil"
    "path/filepath"
    
    v2 "github.com/cnoe-io/idpbuilder/pkg/oci/api/v2"
)

// MultiFormatLoader handles loading certificates from multiple formats
type MultiFormatLoader struct {
    parsers map[v2.CertFormat]FormatParser
    detector FormatDetector
}

// NewMultiFormatLoader creates a new multi-format certificate loader
func NewMultiFormatLoader() *MultiFormatLoader {
    return &MultiFormatLoader{
        parsers: map[v2.CertFormat]FormatParser{
            v2.CertFormatPEM:    &PEMParser{},
            v2.CertFormatDER:    &DERParser{},
            v2.CertFormatPKCS7:  &PKCS7Parser{},
            v2.CertFormatPKCS12: &PKCS12Parser{},
        },
        detector: &MagicBytesDetector{},
    }
}

// LoadBundle loads a certificate bundle with auto-detection
func (l *MultiFormatLoader) LoadBundle(ctx context.Context, path string) (*v2.CertBundle, error)

// LoadPEM loads certificates from PEM format
func (l *MultiFormatLoader) LoadPEM(ctx context.Context, data []byte) (*v2.CertBundle, error)

// LoadDER loads certificates from DER format  
func (l *MultiFormatLoader) LoadDER(ctx context.Context, data []byte) (*v2.CertBundle, error)

// LoadPKCS7 loads certificates from PKCS7 format
func (l *MultiFormatLoader) LoadPKCS7(ctx context.Context, data []byte) (*v2.CertBundle, error)

// LoadPKCS12 loads certificates from PKCS12 format with password
func (l *MultiFormatLoader) LoadPKCS12(ctx context.Context, data []byte, password string) (*v2.CertBundle, error)

// DetectFormat auto-detects certificate format from file content
func (l *MultiFormatLoader) DetectFormat(data []byte) (v2.CertFormat, error)
```

**Key Implementation Details**:
1. Auto-detection using magic bytes and content analysis
2. Graceful fallback between formats
3. Context support for cancellation
4. Comprehensive error handling with v2.CertificateError
5. Support for certificate chains

### Step 2: Implement Parser Utilities (200 lines)
**File**: `pkg/oci/certificates/parser.go`

```go
package certificates

import (
    "crypto/x509"
    "errors"
    "fmt"
    "time"
    
    v2 "github.com/cnoe-io/idpbuilder/pkg/oci/api/v2"
)

// FormatParser interface for format-specific parsers
type FormatParser interface {
    Parse(data []byte) (*v2.CertBundle, error)
    Validate(data []byte) error
}

// CertificateParser provides parsing utilities
type CertificateParser struct {
    strictMode bool
    maxChainDepth int
}

// ParseCertificateChain parses and validates a certificate chain
func (p *CertificateParser) ParseCertificateChain(certs []*x509.Certificate) ([]*x509.Certificate, error)

// BuildChain builds a complete certificate chain from a collection
func (p *CertificateParser) BuildChain(leaf *x509.Certificate, intermediates []*x509.Certificate) ([]*x509.Certificate, error)

// ValidateCertificate performs comprehensive certificate validation
func (p *CertificateParser) ValidateCertificate(cert *x509.Certificate) error

// ExtractCAs separates CA certificates from end-entity certificates
func (p *CertificateParser) ExtractCAs(certs []*x509.Certificate) ([]*x509.Certificate, []*x509.Certificate)

// SortByHierarchy sorts certificates in chain order
func (p *CertificateParser) SortByHierarchy(certs []*x509.Certificate) []*x509.Certificate

// ValidateChain validates a complete certificate chain
func (p *CertificateParser) ValidateChain(chain []*x509.Certificate) error

// ConvertToBundle creates a CertBundle from parsed certificates
func (p *CertificateParser) ConvertToBundle(certs []*x509.Certificate, format v2.CertFormat, source string) *v2.CertBundle
```

**Key Implementation Details**:
1. Chain building with proper hierarchy
2. Certificate validation (expiry, key usage, extensions)
3. CA certificate identification
4. Chain sorting and validation
5. Error reporting with detailed context

### Step 3: Implement Format-Specific Handlers (150 lines)
**File**: `pkg/oci/certificates/formats.go`

```go
package certificates

import (
    "crypto/x509"
    "encoding/pem"
    "errors"
    
    "golang.org/x/crypto/pkcs12"
    v2 "github.com/cnoe-io/idpbuilder/pkg/oci/api/v2"
)

// PEMParser handles PEM format certificates
type PEMParser struct {
    parser *CertificateParser
}

func (p *PEMParser) Parse(data []byte) (*v2.CertBundle, error) {
    var certs []*x509.Certificate
    
    for len(data) > 0 {
        block, rest := pem.Decode(data)
        if block == nil {
            break
        }
        
        if block.Type == "CERTIFICATE" {
            cert, err := x509.ParseCertificate(block.Bytes)
            if err != nil {
                return nil, err
            }
            certs = append(certs, cert)
        }
        data = rest
    }
    
    return p.parser.ConvertToBundle(certs, v2.CertFormatPEM, ""), nil
}

// DERParser handles DER format certificates
type DERParser struct {
    parser *CertificateParser
}

func (p *DERParser) Parse(data []byte) (*v2.CertBundle, error)
func (p *DERParser) Validate(data []byte) error

// PKCS7Parser handles PKCS7 format certificates
type PKCS7Parser struct {
    parser *CertificateParser
}

func (p *PKCS7Parser) Parse(data []byte) (*v2.CertBundle, error)
func (p *PKCS7Parser) Validate(data []byte) error

// PKCS12Parser handles PKCS12 format certificates
type PKCS12Parser struct {
    parser *CertificateParser
}

func (p *PKCS12Parser) Parse(data []byte) (*v2.CertBundle, error)
func (p *PKCS12Parser) ParseWithPassword(data []byte, password string) (*v2.CertBundle, error)

// FormatDetector auto-detects certificate format
type FormatDetector interface {
    DetectFormat(data []byte) (v2.CertFormat, error)
}

// MagicBytesDetector uses magic bytes for format detection
type MagicBytesDetector struct{}

func (d *MagicBytesDetector) DetectFormat(data []byte) (v2.CertFormat, error) {
    // Check for PEM format (starts with "-----BEGIN")
    if len(data) > 10 && string(data[:10]) == "-----BEGIN" {
        return v2.CertFormatPEM, nil
    }
    
    // Check for DER format (ASN.1 sequence)
    if len(data) > 2 && data[0] == 0x30 {
        return v2.CertFormatDER, nil
    }
    
    // Check for PKCS7 (specific OID patterns)
    // Check for PKCS12 (specific magic bytes)
    
    return "", errors.New("unable to detect certificate format")
}

// ConvertFormat converts certificates between formats
func ConvertFormat(bundle *v2.CertBundle, targetFormat v2.CertFormat) (*v2.CertBundle, error)
```

**Key Implementation Details**:
1. Format-specific parsing logic
2. Magic byte detection for auto-detection
3. Support for encrypted PKCS12 with passwords
4. Format conversion capabilities
5. Validation before parsing

### Step 4: Implement Comprehensive Tests (100 lines)
**File**: `pkg/oci/certificates/loader_test.go`

```go
package certificates

import (
    "context"
    "testing"
    "time"
    
    v2 "github.com/cnoe-io/idpbuilder/pkg/oci/api/v2"
)

func TestMultiFormatLoader_LoadPEM(t *testing.T) {
    // Test loading valid PEM certificates
    // Test loading PEM with multiple certificates
    // Test loading PEM with mixed content
}

func TestMultiFormatLoader_LoadDER(t *testing.T) {
    // Test loading valid DER certificate
    // Test loading DER with invalid data
}

func TestMultiFormatLoader_LoadPKCS7(t *testing.T) {
    // Test loading PKCS7 bundle
    // Test extracting certificate chain
}

func TestMultiFormatLoader_LoadPKCS12(t *testing.T) {
    // Test loading PKCS12 with password
    // Test loading PKCS12 with wrong password
}

func TestMultiFormatLoader_DetectFormat(t *testing.T) {
    // Test auto-detection for each format
    // Test detection with corrupted data
}

func TestCertificateParser_BuildChain(t *testing.T) {
    // Test chain building with complete chain
    // Test chain building with missing intermediate
}

func TestCertificateParser_ValidateCertificate(t *testing.T) {
    // Test validation of expired certificate
    // Test validation of self-signed certificate
    // Test validation of invalid key usage
}

func BenchmarkMultiFormatLoader_LoadPEM(b *testing.B) {
    // Benchmark PEM loading performance
}

func BenchmarkMultiFormatLoader_DetectFormat(b *testing.B) {
    // Benchmark format detection
}
```

**Test Coverage Requirements**:
- All supported formats tested
- Invalid certificate handling
- Chain validation scenarios
- Performance benchmarks
- Edge cases (empty files, corrupted data, mixed formats)

## Size Management
- **Target Lines**: 700 lines
- **Current Breakdown**:
  - `loader.go`: 250 lines
  - `parser.go`: 200 lines
  - `formats.go`: 150 lines
  - `loader_test.go`: 100 lines
- **Measurement Tool**: `$PROJECT_ROOT/tools/line-counter.sh`
- **Check Frequency**: After each file implementation
- **Split Threshold**: 700 lines (warning), 800 lines (stop)

## Test Requirements
- **Unit Tests**: 85% coverage minimum
- **Integration Tests**: Test with real certificate files
- **Performance Tests**: Benchmark all format loaders
- **Test Scenarios**:
  - Valid certificates in all formats
  - Invalid/corrupted certificates
  - Certificate chains with missing intermediates
  - Expired certificates
  - Self-signed certificates
  - Password-protected PKCS12 files

## Pattern Compliance
- **idpbuilder Patterns**:
  - Use context for all operations
  - Return v2.CertificateError for certificate-specific errors
  - Follow existing package structure under pkg/oci/
  - Use interfaces for extensibility
  
- **Security Requirements**:
  - Validate all certificates before returning
  - Support certificate expiry checking
  - Verify certificate chains properly
  - Handle password-protected formats securely
  - Clear sensitive data from memory after use
  
- **Performance Targets**:
  - Load 100 certificates < 100ms
  - Format detection < 1ms
  - Memory efficient for large bundles

## Integration Points
- **E3.1.1 Dependencies**: Import and use v2.CertBundle, v2.CertFormat, v2.CertificateError
- **E3.1.3 Integration**: Certificate Service will use this loader
- **E3.1.4 Storage**: Storage component will persist loaded bundles
- **Phase 2 Registry**: Registry client will use loaded certificates for TLS

## Implementation Guidelines

### Critical Implementation Order
1. **Create pkg/oci/certificates/ directory first**
2. **Import v2 package from E3.1.1**: `v2 "github.com/cnoe-io/idpbuilder/pkg/oci/api/v2"`
3. **Start with loader.go** - Core loading logic
4. **Then parser.go** - Parsing utilities needed by format handlers
5. **Then formats.go** - Format-specific implementations
6. **Finally loader_test.go** - Comprehensive tests

### Key Implementation Details
1. **Auto-detection MUST be reliable** - Use multiple heuristics
2. **Support partial chains** - Don't fail if intermediate CAs missing
3. **Graceful error handling** - Provide clear error messages
4. **Thread-safe implementation** - Loaders may be used concurrently
5. **Memory efficiency** - Stream large files when possible

### Error Handling Strategy
- Use v2.CertificateError for all certificate-specific errors
- Include certificate details in error messages
- Distinguish between format errors and validation errors
- Provide actionable error messages

## Success Criteria
- All 4 certificate formats supported (PEM, DER, PKCS7, PKCS12)
- Auto-detection works reliably for all formats
- Certificate chains handled properly
- Implementation stays under 700 lines
- Tests achieve >85% coverage
- Performance benchmarks meet targets
- Clean integration with E3.1.1 interfaces

## Next Steps After Implementation
1. Run line counter to verify size compliance
2. Run tests to verify functionality
3. Run benchmarks to verify performance
4. Document usage examples in code comments
5. Update work-log.md with completion status
6. Prepare for integration with E3.1.3 Certificate Service

## Usage Example
```go
loader := NewMultiFormatLoader()

// Auto-detect and load
bundle, err := loader.LoadBundle(ctx, "/path/to/cert.pem")

// Or load with known format
pemData, _ := ioutil.ReadFile("/path/to/cert.pem")
bundle, err := loader.LoadPEM(ctx, pemData)

// Load PKCS12 with password
p12Data, _ := ioutil.ReadFile("/path/to/cert.p12")
bundle, err := loader.LoadPKCS12(ctx, p12Data, "password")
```