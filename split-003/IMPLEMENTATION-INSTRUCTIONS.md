# Implementation Instructions for E3.1.4 Split-003

## Split Overview
**Split**: 003 of 4
**Target Size**: ~390 lines maximum
**Purpose**: Pool Management

## Files to Implement

### 1. pkg/oci/certificates/pool.go (200 lines)
Implement CertPoolManager:
- Implements `CertPoolManager` interface from split-001
- System/custom pool separation
- Hot-reload logic for certificate updates
- Pool operations (add, remove, get certificates)
- Thread-safe pool access

Key methods:
- `NewCertPoolManager(store CertificateStore) (*PoolManager, error)`
- `CreatePool(ctx context.Context, name string) error`
- `DeletePool(ctx context.Context, name string) error`
- `AddCertificateToPool(ctx context.Context, poolName, certID string) error`
- `RemoveCertificateFromPool(ctx context.Context, poolName, certID string) error`
- `GetPoolCertificates(ctx context.Context, poolName string) ([]*Certificate, error)`
- `ListPools(ctx context.Context) ([]string, error)`
- `RefreshPools(ctx context.Context) error` - hot-reload

### 2. pkg/oci/certificates/validation.go (100 lines)
Certificate validation pipeline:
- Implements `CertificateValidator` interface from split-001
- Chain validation
- Expiry checking
- Permission validation
- Custom validation rules

Key methods:
- `NewValidator() *Validator`
- `ValidateCertificate(cert *Certificate) (*ValidationResult, error)`
- `ValidatePEM(pemData []byte) (*ValidationResult, error)`
- `ValidateCertificateChain(certs []*Certificate) (*ValidationResult, error)`
- `AddValidationRule(rule ValidationRule) error`
- `RemoveValidationRule(id string) error`
- `ListValidationRules() []ValidationRule`

### 3. pkg/oci/certificates/pool_test.go (90 lines)
Pool management tests:
- Test pool creation/deletion
- Test certificate add/remove from pools
- Test hot-reload functionality
- Test concurrent pool access
- Validation pipeline tests

## Dependencies
- Import types and interfaces from split-001
- Import storage implementation from split-002
- Use standard library crypto/x509 for certificate operations

## Implementation Requirements

1. **Interface Compliance**: Must implement CertPoolManager and CertificateValidator exactly as defined
2. **Thread Safety**: Use mutex for all pool operations
3. **Hot-Reload**: Support dynamic pool updates without restart
4. **Validation Pipeline**: Chain multiple validation rules
5. **Size Compliance**: Must not exceed 390 lines total

## Branch Setup
Create branch: `idpbuilder-oci-mgmt/phase3/wave1/E3.1.4-trust-store-split-003`

## Implementation Notes
- System pool: certificates from system trust store
- Custom pool: user-provided certificates
- Pool data structure: map[string][]string (pool name to cert IDs)
- Validation rules should be composable
- Include file watching integration if space permits

## Validation Steps
1. Run `tools/line-counter.sh` to verify size
2. Ensure interfaces fully implemented
3. Test hot-reload with concurrent access
4. Verify validation pipeline works correctly