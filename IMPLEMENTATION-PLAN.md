# Registry Authentication & Certificate Types Implementation Plan

## 🚨 CRITICAL EFFORT METADATA (FROM WAVE PLAN)
**Branch**: idpbuilder-oci-mgmt/phase1/wave1/auth-cert-types  
**Can Parallelize**: Yes (Wave 1 - all efforts can run simultaneously)  
**Parallel With**: [oci-stack-types, error-progress-types]  
**Size Estimate**: 400 lines (MUST stay within limit)  
**Dependencies**: None (foundational effort)

<!-- ⚠️ ORCHESTRATOR METADATA PLACEHOLDER - DO NOT REMOVE ⚠️ -->
<!-- The orchestrator will add infrastructure metadata below: -->
<!-- WORKING_DIRECTORY: /home/vscode/workspaces/idpbuilder-oci-mgmt/efforts/phase1/wave1/auth-cert-types -->
<!-- BRANCH: idpbuilder-oci-mgmt/phase1/wave1/auth-cert-types -->
<!-- REMOTE: origin/idpbuilder-oci-mgmt/phase1/wave1/auth-cert-types -->
<!-- BASE_BRANCH: software-factory-2.0 -->
<!-- SW Engineers MUST read this metadata to navigate to the correct directory -->
<!-- END PLACEHOLDER -->

## Overview
- **Effort**: E1.1.2 - Registry Authentication & Certificate Types
- **Phase**: 1, Wave: 1  
- **Estimated Size**: 400 lines
- **Implementation Time**: 4-6 hours
- **Purpose**: Define authentication interfaces, credential types, certificate handling, and validation logic for OCI registry operations

## File Structure
All files go under `pkg/oci/auth/` in the effort directory:

```
efforts/phase1/wave1/auth-cert-types/
└── pkg/
    └── oci/
        └── auth/
            ├── interfaces.go    # 80 lines - Authentication contracts
            ├── types.go         # 150 lines - Auth data structures  
            ├── validation.go    # 70 lines - Validation logic
            └── auth_test.go     # 100 lines - Test coverage
```

## Implementation Steps

### Step 1: Create Directory Structure (5 minutes)
```bash
# From effort root directory
mkdir -p pkg/oci/auth
```

### Step 2: Implement interfaces.go (45 minutes)
- **Purpose**: Define authentication service contracts
- **Size**: 80 lines
- **Content**:
  - `AuthProvider` interface - credential management
  - `CredentialStore` interface - persistent storage
  - `CertificateProvider` interface - certificate handling
  - `TokenManager` interface - token lifecycle
- **Key Methods**:
  - GetCredentials, ValidateCredentials, RefreshToken
  - Load/Save/Delete for credential storage
  - GetCertBundle, LoadFromKind, ValidateCertificate
  - Token validation and caching

### Step 3: Implement types.go (60 minutes)
- **Purpose**: Define all authentication data structures
- **Size**: 150 lines
- **Content**:
  - `Credentials` struct - username/password or token
  - `Token` struct - OAuth2/Bearer token support
  - `TLSConfig` struct - TLS/certificate configuration
  - `AuthConfig` struct - authentication configuration
  - `BasicAuth`, `TokenAuth`, `OAuth2Auth` structs
  - `FileStore`, `MemoryStore` structs for storage
  - `CertificateInfo` struct - certificate metadata
  - `AuthResult` struct - authentication results
  - `TLSVersionMap` - TLS version constants
- **Validation Tags**: Use struct tags for validation rules

### Step 4: Implement validation.go (30 minutes)
- **Purpose**: Validation logic for auth types
- **Size**: 70 lines
- **Content**:
  - `NewAuthValidator()` - create validator instance
  - `validateHostnamePort()` - custom hostname:port validation
  - `ValidateCredentials()` - credential validation with expiry check
  - `ValidateCertificate()` - X.509 certificate validation
  - `ValidateRegistryURL()` - registry URL format validation
- **Dependencies**: github.com/go-playground/validator/v10

### Step 5: Implement auth_test.go (45 minutes)
- **Purpose**: Comprehensive test coverage
- **Size**: 100 lines
- **Content**:
  - `TestCredentials_Validation` - test credential validation
  - `TestTLSConfig_Validation` - test TLS config validation
  - Table-driven tests with multiple scenarios
  - Edge cases: expired credentials, invalid formats
  - Target: >90% code coverage

### Step 6: Validation & Testing (30 minutes)
```bash
# Initialize go module if needed
go mod init github.com/cnoe-io/idpbuilder/efforts/phase1/wave1/auth-cert-types

# Add dependencies
go get github.com/go-playground/validator/v10@v10.15.5

# Run tests
go test ./pkg/oci/auth/... -v -cover

# Check test coverage
go test ./pkg/oci/auth/... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### Step 7: Size Verification (5 minutes)
```bash
# Find project root and use line counter
PROJECT_ROOT=/home/vscode/workspaces/idpbuilder-oci-mgmt
$PROJECT_ROOT/tools/line-counter.sh
# Must show ≤400 lines for this effort
```

## Size Management
- **Estimated Lines**: 400 total
  - interfaces.go: 80 lines
  - types.go: 150 lines  
  - validation.go: 70 lines
  - auth_test.go: 100 lines
- **Measurement Tool**: `/home/vscode/workspaces/idpbuilder-oci-mgmt/tools/line-counter.sh`
- **Check Frequency**: After each file completion
- **Split Threshold**: N/A (well within 400-line limit)
- **Buffer**: 0 lines (exactly at limit - be precise!)

## Test Requirements
- **Unit Tests**: 90% coverage minimum
  - Cover all validation functions
  - Test credential expiration logic
  - Validate TLS configuration options
  - Test custom validators (hostname_port)
- **Edge Cases**:
  - Expired credentials
  - Invalid certificate formats
  - Missing required fields
  - Mutual exclusion (username vs token)
- **Test Files**: 
  - `auth_test.go` - comprehensive unit tests

## Pattern Compliance

### Go Best Practices
- Use interfaces for abstraction
- Struct tags for validation
- Table-driven tests
- Error wrapping with context
- Proper nil checks

### IDPBuilder Patterns
- Consistent error handling
- Validation at boundaries
- Clear separation of concerns
- Testable code structure

### Security Requirements
- No credentials in logs
- Secure credential storage interface
- Certificate validation
- TLS version enforcement

## Integration Points
This effort provides foundational types that will be used by:
- Phase 1, Wave 2: Build service implementation (will use AuthProvider)
- Phase 1, Wave 3: Registry client (will use Credentials, TLSConfig)
- Phase 2: Controllers (will use authentication for registry operations)

## Dependencies Management
```go
// go.mod additions required:
require (
    github.com/go-playground/validator/v10 v10.15.5
)
```

## Success Criteria
- [ ] All 4 files created and properly structured
- [ ] Code compiles without errors
- [ ] All tests pass
- [ ] Test coverage ≥90%
- [ ] Size ≤400 lines (verified with line-counter.sh)
- [ ] All interfaces clearly documented
- [ ] Validation logic comprehensive
- [ ] No security vulnerabilities

## Common Pitfalls to Avoid
1. **Over-engineering**: Keep interfaces simple and focused
2. **Missing validation**: Every field that can be validated should be
3. **Poor test coverage**: Aim for >90% coverage
4. **Size creep**: Monitor size after each file
5. **Circular dependencies**: Keep auth package self-contained

## Review Checklist
- [ ] Interfaces follow Go conventions
- [ ] All exported types have comments
- [ ] Validation tags are correct
- [ ] Tests cover edge cases
- [ ] No hardcoded values
- [ ] Error messages are descriptive
- [ ] Code is idiomatic Go

## Notes for SW Engineer
1. Start with interfaces.go to establish contracts
2. Implement types.go with all validation tags
3. Keep validation logic separate in validation.go
4. Write tests incrementally as you go
5. Use the line counter frequently to stay within limits
6. This is types-only - no implementation logic needed
7. Focus on clean, well-documented interfaces