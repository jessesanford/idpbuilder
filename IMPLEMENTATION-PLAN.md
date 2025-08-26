# Implementation Plan: Security Layer - Split 001 (Security Manager)

## <¯ Effort Overview
**Effort ID**: effort4-security-split-001
**Target Size**: 386 lines MAXIMUM
**Purpose**: Security orchestration and policy management
**Order**: IMPLEMENT AFTER split-002 (depends on crypto layer)

## =¨ CRITICAL REQUIREMENTS
1. **SIZE LIMIT**: 386 lines (well under limit)
2. **DEPENDS ON SPLIT-002**: Use Signer/Verifier from split-002
3. **ORCHESTRATION LAYER**: Coordinate security operations

## =Á Files to Implement

### 1. pkg/oci/security/manager.go (386 lines)
**Purpose**: Security orchestration, policy enforcement, and coordination

**Core Implementation**:
```go
package security

import (
    "context"
    "fmt"
    "github.com/jessesanford/idpbuilder/pkg/oci/api"
)

type SecurityManager struct {
    signer     api.Signer
    verifier   api.Verifier
    policies   []api.SecurityPolicy
    trustStore *TrustStore
}

// Key methods to implement:
- NewSecurityManager(config *SecurityConfig) (*SecurityManager, error)
- SignArtifact(artifact api.Artifact) (*api.SignedArtifact, error)
- VerifyArtifact(artifact api.SignedArtifact) error
- EnforcePolicy(artifact api.Artifact, policy api.SecurityPolicy) error
- AddPolicy(policy api.SecurityPolicy) error
- RemovePolicy(policyID string) error
- GetTrustChain(keyID string) ([]api.Certificate, error)
- RotateKeys() error
```

### 2. API Imports
**From split-002**: You'll need the crypto interfaces

## =' Implementation Steps

### Step 1: Copy API types from split-002
```bash
# Copy the crypto API from split-002
cp -r ../split-002/pkg/oci/api pkg/oci/
```

### Step 2: Copy manager.go from parent
```bash
cp ../pkg/oci/security/manager.go pkg/oci/security/
```

### Step 3: Import crypto implementations
In manager.go, ensure you're using the interfaces:
```go
import (
    "github.com/jessesanford/idpbuilder/pkg/oci/api"
    // The actual Signer and Verifier will be from split-002
)
```

### Step 4: Implement orchestration
- Use api.Signer for signing operations
- Use api.Verifier for verification
- Add policy enforcement layer
- Implement key rotation logic
- Add trust chain management

### Step 5: Verify compilation
```bash
cd pkg/oci/security
go build .
```

### Step 6: Measure size
```bash
/home/vscode/workspaces/idpbuilder-oci-mgmt/tools/line-counter.sh -c branch
# MUST be d386 lines
```

##  Success Criteria
- [ ] manager.go implements security orchestration (d386 lines)
- [ ] Uses Signer/Verifier interfaces from split-002
- [ ] Implements policy enforcement
- [ ] Handles key rotation
- [ ] Code compiles successfully
- [ ] Total d386 lines

## =¨ Critical Notes
1. **DEPENDS ON SPLIT-002**: Must use the crypto interfaces
2. **ORCHESTRATION FOCUS**: Don't re-implement crypto
3. **POLICY LAYER**: Add value on top of basic crypto
4. **SIZE COMFORTABLE**: 386 lines gives plenty of room

## Integration Points
- Uses api.Signer from split-002 for signing
- Uses api.Verifier from split-002 for verification  
- Adds SecurityPolicy enforcement on top
- Provides unified SecurityManager interface