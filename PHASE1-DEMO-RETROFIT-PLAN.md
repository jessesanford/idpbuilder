# Phase 1 Demo Retrofit Plan

## Overview
This document provides demo requirements for all Phase 1 efforts per R330 Demo Requirements Mandate. Each effort requires demonstrable features with executable scripts and documentation.

## Size Impact Summary
- Each effort demo adds ~100-150 lines
- Total retrofit impact: ~400-600 lines across all efforts
- All efforts remain under 800 line limit with demos included

---

## E1.1.1 - Kind Certificate Extraction (pkg/certs - extraction)

### Features Discovered
- **KindCertExtractor**: Extracts certificates from Kind cluster pods
- **CertificateStorage**: Stores certificates on disk with metadata
- **KindClient**: Interacts with Kind clusters via kubectl
- **DefaultCertValidator**: Validates certificate expiry and properties
- **Retry Logic**: Configurable retry attempts with timeout

### Demo Objectives (R330 Compliant)
1. Extract certificates from a running Kind cluster
2. Validate extracted certificates for expiry
3. Store certificates with proper metadata
4. Demonstrate retry mechanism on failure
5. Show diagnostic output for troubleshooting

### Demo Scenarios

#### Scenario 1: Extract Registry Certificates from Kind
```bash
# Setup: Create Kind cluster with registry
kind create cluster --name demo-cluster
kubectl create namespace registry

# Deploy test registry pod with certificates
kubectl apply -f test-data/registry-with-certs.yaml

# Execute extraction
./demo-features.sh extract-certs \
  --cluster demo-cluster \
  --namespace registry \
  --pod-selector "app=registry" \
  --cert-path "/certs/ca.crt"

# Expected output:
# ✓ Connected to Kind cluster: demo-cluster
# ✓ Found 1 pod(s) matching selector
# ✓ Extracting certificate from pod: registry-7f6d9...
# ✓ Certificate validated: CN=Registry CA, expires 2026-01-01
# ✓ Certificate stored: ~/.idpbuilder/certs/demo-cluster/registry-ca.crt
```

#### Scenario 2: Validate Certificate Chain
```bash
# Extract full certificate chain
./demo-features.sh extract-chain \
  --cluster demo-cluster \
  --namespace registry \
  --include-intermediates

# Validate chain
./demo-features.sh validate-chain \
  --cert ~/.idpbuilder/certs/demo-cluster/registry-ca.crt \
  --show-details

# Expected output:
# Certificate Chain Validation:
# └── Root CA: CN=Demo Root CA (Valid)
#     └── Intermediate: CN=Demo Intermediate (Valid)
#         └── Leaf: CN=registry.local (Valid)
# ✓ Chain validation successful
```

#### Scenario 3: Handle Extraction Failures
```bash
# Demonstrate retry mechanism with unreachable pod
./demo-features.sh extract-certs \
  --cluster demo-cluster \
  --namespace nonexistent \
  --retry-attempts 3 \
  --timeout 5s

# Expected output:
# ⚠ Attempt 1/3: No pods found in namespace: nonexistent
# ⚠ Attempt 2/3: No pods found in namespace: nonexistent
# ⚠ Attempt 3/3: No pods found in namespace: nonexistent
# ✗ Certificate extraction failed after 3 attempts
```

### Demo Deliverables
1. `demo-features.sh` - Executable demo script (~50 lines)
2. `test-data/registry-with-certs.yaml` - Test registry deployment (~30 lines)
3. `test-data/demo-certs/` - Sample certificates for testing
4. `DEMO.md` - User documentation (~70 lines)

### Size Impact
- Current implementation: ~450 lines
- Demo additions: ~150 lines
- **Total with demos: ~600 lines** ✓ (under 800 limit)

---

## E1.1.2 - Registry TLS Trust (pkg/certs - trust)

### Features Discovered
- **TrustStoreManager**: Manages trusted certificates per registry
- **DefaultTrustStore**: Implements certificate storage and retrieval
- **HTTP Client Creation**: Configures TLS for registry connections
- **Insecure Mode**: Optional bypass for development
- **Persistence**: Saves trusted certs to disk

### Demo Objectives (R330 Compliant)
1. Add registry certificates to trust store
2. Configure HTTP client with custom CA
3. Demonstrate secure registry connection
4. Show insecure mode with warnings
5. Persist and reload trust configuration

### Demo Scenarios

#### Scenario 1: Configure Registry Trust
```bash
# Add custom CA to trust store
./demo-features.sh trust-add \
  --registry registry.example.com:5000 \
  --cert-file test-data/demo-certs/registry-ca.crt

# Verify trust configuration
./demo-features.sh trust-list

# Expected output:
# Trusted Registries:
# ├── registry.example.com:5000
# │   └── CA: CN=Example Registry CA, SHA256: abc123...
# └── Total: 1 registry, 1 certificate
```

#### Scenario 2: Test Registry Connection
```bash
# Test connection with trusted certificate
./demo-features.sh test-connection \
  --registry registry.example.com:5000 \
  --image alpine:latest

# Expected output:
# Testing registry connection...
# ✓ TLS handshake successful
# ✓ Certificate chain validated
# ✓ Registry API accessible
# ✓ Image manifest retrieved: alpine:latest
```

#### Scenario 3: Insecure Mode Demo
```bash
# Enable insecure mode for specific registry
./demo-features.sh trust-insecure \
  --registry localhost:5000 \
  --confirm

# Expected output:
# ⚠⚠⚠⚠⚠⚠⚠⚠⚠⚠
# WARNING: Certificate validation disabled for localhost:5000
# This should ONLY be used in development environments
# ⚠⚠⚠⚠⚠⚠⚠⚠⚠⚠

# Test with insecure mode
./demo-features.sh test-connection \
  --registry localhost:5000 \
  --insecure

# Shows connection without cert validation
```

#### Scenario 4: Trust Store Persistence
```bash
# Export trust configuration
./demo-features.sh trust-export \
  --output trust-config.json

# Clear and reimport
./demo-features.sh trust-clear --confirm
./demo-features.sh trust-import \
  --input trust-config.json

# Expected output:
# ✓ Exported 3 registry configurations
# ✓ Trust store cleared
# ✓ Imported 3 registry configurations
```

### Demo Deliverables
1. `demo-features.sh` - Extended with trust commands (~60 lines)
2. `test-data/trust-config.json` - Sample trust configuration (~20 lines)
3. `DEMO-TRUST.md` - Trust management guide (~50 lines)

### Size Impact
- Current implementation: ~250 lines
- Demo additions: ~130 lines
- **Total with demos: ~380 lines** ✓ (under 800 limit)

---

## E1.2.1 - Certificate Validation (pkg/certvalidation)

### Features Discovered
- **ChainValidator**: Validates complete certificate chains
- **X509 Utilities**: Certificate parsing and analysis
- **Custom Verification Options**: Time-based validation
- **Intermediate Certificate Support**: Full chain building
- **Detailed Error Reporting**: Specific validation failures

### Demo Objectives (R330 Compliant)
1. Validate certificate chains with various scenarios
2. Demonstrate expired certificate detection
3. Show intermediate certificate handling
4. Display detailed validation errors
5. Time-based validation scenarios

### Demo Scenarios

#### Scenario 1: Full Chain Validation
```bash
# Validate complete certificate chain
./demo-features.sh validate \
  --cert test-data/demo-certs/server.crt \
  --intermediate test-data/demo-certs/intermediate.crt \
  --root test-data/demo-certs/root-ca.crt \
  --verbose

# Expected output:
# Certificate Chain Analysis:
# 
# Subject: CN=server.example.com
# Issuer: CN=Example Intermediate CA
# Valid From: 2024-01-01 00:00:00 UTC
# Valid Until: 2025-01-01 00:00:00 UTC
# Key Usage: Digital Signature, Key Encipherment
# 
# Chain Building:
# ✓ Leaf certificate valid
# ✓ Intermediate certificate found
# ✓ Root certificate trusted
# ✓ Chain validation successful
```

#### Scenario 2: Detect Expired Certificates
```bash
# Test with expired certificate
./demo-features.sh validate \
  --cert test-data/demo-certs/expired.crt \
  --check-expiry

# Expected output:
# ✗ Certificate Validation Failed:
#   Certificate expired on 2023-12-31 23:59:59 UTC
#   Current time: 2024-09-10 21:30:00 UTC
#   Expired for: 253 days
```

#### Scenario 3: Missing Intermediate Detection
```bash
# Attempt validation without intermediate
./demo-features.sh validate \
  --cert test-data/demo-certs/server.crt \
  --root test-data/demo-certs/root-ca.crt \
  --strict

# Expected output:
# ✗ Chain Building Failed:
#   Cannot build path from leaf to root
#   Missing intermediate certificate:
#   Expected Issuer: CN=Example Intermediate CA
#   
#   Hint: Add intermediate certificate with --intermediate flag
```

#### Scenario 4: Time-Based Validation
```bash
# Validate certificate at specific time
./demo-features.sh validate \
  --cert test-data/demo-certs/future.crt \
  --validate-at "2025-06-01T00:00:00Z"

# Expected output:
# Validating certificate at: 2025-06-01 00:00:00 UTC
# ✓ Certificate will be valid at specified time
# Days until valid: 264
```

### Demo Deliverables
1. `demo-features.sh` - Validation commands (~70 lines)
2. `test-data/demo-certs/` - Various test certificates
3. `DEMO-VALIDATION.md` - Validation guide (~60 lines)

### Size Impact
- Current implementation: ~350 lines
- Demo additions: ~130 lines
- **Total with demos: ~480 lines** ✓ (under 800 limit)

---

## E1.2.2 - Fallback Strategies (pkg/fallback, pkg/insecure)

### Features Discovered
- **FallbackManager**: Coordinates multiple fallback strategies
- **Strategy Priority System**: Ordered execution of strategies
- **InsecureHandler**: Manages insecure mode with warnings
- **Retry Logic**: Configurable retry with delays
- **Registry-Specific Configuration**: Per-registry fallback rules

### Demo Objectives (R330 Compliant)
1. Demonstrate fallback strategy execution order
2. Show automatic retry with delays
3. Display insecure mode warnings
4. Test registry-specific configurations
5. Simulate failure recovery scenarios

### Demo Scenarios

#### Scenario 1: Fallback Strategy Chain
```bash
# Configure fallback strategies
./demo-features.sh fallback-config \
  --add-strategy "system-ca" --priority 1 \
  --add-strategy "kind-extract" --priority 2 \
  --add-strategy "insecure" --priority 3

# Test with unreachable registry
./demo-features.sh test-fallback \
  --registry secure.example.com:443 \
  --verbose

# Expected output:
# Attempting connection to secure.example.com:443
# 
# Strategy 1: System CA Trust
# ✗ No system certificate found for domain
# 
# Strategy 2: Kind Certificate Extraction
# → Searching Kind clusters for certificates...
# ✓ Found certificate in cluster 'prod-cluster'
# ✓ Connection established with extracted certificate
# 
# Result: SUCCESS (Strategy: kind-extract)
```

#### Scenario 2: Automatic Retry Mechanism
```bash
# Test retry with transient failures
./demo-features.sh test-retry \
  --registry flaky.example.com:5000 \
  --max-retries 3 \
  --retry-delay 2s

# Expected output:
# Connecting to flaky.example.com:5000
# 
# Attempt 1/3: Connection timeout
# ⏳ Waiting 2s before retry...
# 
# Attempt 2/3: TLS handshake failed
# ⏳ Waiting 2s before retry...
# 
# Attempt 3/3: Success!
# ✓ Connected after 3 attempts
```

#### Scenario 3: Insecure Fallback with Warnings
```bash
# Enable insecure as last resort
./demo-features.sh test-fallback \
  --registry untrusted.local:5000 \
  --allow-insecure \
  --require-confirmation

# Expected output:
# All secure strategies failed for untrusted.local:5000
# 
# ⚠⚠⚠⚠⚠⚠⚠⚠⚠⚠
# WARNING: About to disable certificate validation
# Registry: untrusted.local:5000
# This is INSECURE and should only be used in development
# ⚠⚠⚠⚠⚠⚠⚠⚠⚠⚠
# 
# Continue with insecure connection? [y/N]: y
# ✓ Connected (INSECURE MODE)
```

#### Scenario 4: Registry-Specific Rules
```bash
# Configure per-registry fallback rules
./demo-features.sh fallback-rules \
  --registry "*.internal" --strategy "corporate-ca" \
  --registry "localhost:*" --allow-insecure \
  --registry "prod-*.com" --strict

# Test different registries
./demo-features.sh test-rules

# Expected output:
# Testing Registry Fallback Rules:
# 
# dev.internal:5000    → Using corporate-ca (rule match)
# localhost:8080       → Insecure allowed (rule match)
# prod-api.com:443     → Strict mode, no fallback (rule match)
# other.example.com    → Default fallback chain (no rule)
```

### Demo Deliverables
1. `demo-features.sh` - Fallback testing commands (~80 lines)
2. `test-data/fallback-config.yaml` - Sample configurations (~30 lines)
3. `DEMO-FALLBACK.md` - Fallback strategy guide (~60 lines)

### Size Impact
- Current implementation: ~400 lines (combined packages)
- Demo additions: ~170 lines
- **Total with demos: ~570 lines** ✓ (under 800 limit)

---

## Integration Points for Phase 2

### Demo Integration Requirements
Each Phase 1 demo must integrate with Phase 2 features:

1. **Registry Operations**: Demos must work with actual registry push/pull
2. **Build Process**: Certificate validation during image builds
3. **Authentication**: Trust establishment before auth attempts
4. **Error Handling**: Graceful degradation when certs unavailable

### Testing Strategy
```bash
# Comprehensive integration test
./demo-features.sh integration-test \
  --setup-kind \
  --extract-certs \
  --configure-trust \
  --validate-chain \
  --test-fallback \
  --cleanup

# Expected: All components work together seamlessly
```

---

## Summary

### Total Size Impact
| Effort | Current | Demo Add | Total | Status |
|--------|---------|----------|--------|---------|
| E1.1.1 | ~450    | ~150     | ~600   | ✓ OK    |
| E1.1.2 | ~250    | ~130     | ~380   | ✓ OK    |
| E1.2.1 | ~350    | ~130     | ~480   | ✓ OK    |
| E1.2.2 | ~400    | ~170     | ~570   | ✓ OK    |

**All efforts remain under 800 line limit with demos included**

### Deliverables Checklist
- [ ] `demo-features.sh` - Master demo script
- [ ] `DEMO.md` - Main documentation
- [ ] `DEMO-TRUST.md` - Trust management guide
- [ ] `DEMO-VALIDATION.md` - Validation guide
- [ ] `DEMO-FALLBACK.md` - Fallback strategy guide
- [ ] `test-data/` - Test certificates and configurations
- [ ] Integration test suite

### Implementation Timeline
1. Week 1: Implement demo-features.sh core
2. Week 2: Add scenario scripts and test data
3. Week 3: Write documentation and guides
4. Week 4: Integration testing and refinement

---

## Approval
This plan ensures all Phase 1 efforts meet R330 Demo Requirements while maintaining size compliance.

**Created**: 2025-09-10
**Status**: Ready for Implementation
**Compliance**: R330 Demo Requirements Mandate ✓