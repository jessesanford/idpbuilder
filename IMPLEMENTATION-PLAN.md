# Implementation Plan: Security Layer - Split 002 (Crypto Operations)

## <¯ Effort Overview
**Effort ID**: effort4-security-split-002
**Target Size**: 649 lines MAXIMUM  
**Purpose**: Core cryptographic signing and verification operations
**Order**: MUST BE IMPLEMENTED FIRST (before split-001)

## =¨ CRITICAL REQUIREMENTS
1. **SIZE LIMIT**: 649 lines (aim for 600 to have buffer)
2. **NO DEPENDENCIES**: This is the foundational layer
3. **COMPLETE INTERFACES**: Manager in split-001 depends on these

## =Á Files to Implement

### 1. pkg/oci/security/signer.go (335 lines)
**Purpose**: Digital signature operations for OCI artifacts

**Core Implementation**:
```go
package security

import (
    "crypto"
    "crypto/rsa"
    "crypto/sha256"
    "crypto/x509"
    "encoding/pem"
    "github.com/jessesanford/idpbuilder/pkg/oci/api"
)

type Signer struct {
    privateKey crypto.PrivateKey
    publicKey  crypto.PublicKey
    algorithm  api.SignatureAlgorithm
}

// Key methods to implement:
- NewSigner(privateKeyPEM []byte) (*Signer, error)
- Sign(data []byte) ([]byte, error)
- SignManifest(manifest api.Manifest) (*api.SignedManifest, error)
- GetPublicKey() ([]byte, error)
- VerifyOwnSignature(data, signature []byte) error
```

### 2. pkg/oci/security/verifier.go (314 lines)
**Purpose**: Signature verification and trust validation

**Core Implementation**:
```go
package security

import (
    "crypto"
    "crypto/rsa"
    "crypto/sha256"
    "crypto/x509"
    "github.com/jessesanford/idpbuilder/pkg/oci/api"
)

type Verifier struct {
    trustedKeys map[string]crypto.PublicKey
    trustStore  *TrustStore
}

// Key methods to implement:
- NewVerifier(trustedKeys [][]byte) (*Verifier, error)
- Verify(data, signature []byte, keyID string) error
- VerifyManifest(manifest api.SignedManifest) error
- AddTrustedKey(keyID string, publicKey []byte) error
- RemoveTrustedKey(keyID string) error
```

### 3. pkg/oci/security/trust_store.go (embedded in verifier.go)
**Purpose**: Manage trusted keys and certificates
**Note**: Keep minimal, embed in verifier.go to save lines

## =' Implementation Steps

### Step 1: Copy existing code from parent directory
```bash
cp ../pkg/oci/security/signer.go pkg/oci/security/
cp ../pkg/oci/security/verifier.go pkg/oci/security/
```

### Step 2: Ensure API types are available
```bash
# Check if api types exist, if not copy from effort1-contracts
cp -r /home/vscode/workspaces/idpbuilder-oci-mgmt/efforts/phase2/wave2/effort1-contracts/pkg/oci/api pkg/oci/
```

### Step 3: Optimize if needed
- If total exceeds 649 lines, optimize:
  - Combine helper functions
  - Remove verbose comments
  - Simplify error handling

### Step 4: Verify compilation
```bash
cd pkg/oci/security
go build .
```

### Step 5: Measure size
```bash
/home/vscode/workspaces/idpbuilder-oci-mgmt/tools/line-counter.sh -c branch
# MUST be d649 lines
```

##  Success Criteria
- [ ] signer.go implements all signing operations (~335 lines)
- [ ] verifier.go implements all verification (~314 lines)
- [ ] Total d649 lines
- [ ] Code compiles independently
- [ ] No dependency on manager.go
- [ ] Interfaces ready for split-001

## =¨ Critical Notes
1. **FOUNDATIONAL LAYER**: Split-001 depends on this
2. **PRESERVE INTERFACES**: Manager expects specific method signatures
3. **NO MANAGER REFERENCES**: This must be independent
4. **SECURITY CRITICAL**: Ensure crypto operations are correct

## Dependencies for Split-001
Split-001 (manager.go) will:
- Import Signer and Verifier types
- Use these for security orchestration
- Add policy enforcement on top