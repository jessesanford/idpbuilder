# Split Plan for Effort 3: Cache Manager & Layer Optimization

## Executive Summary
**Current Status**: 838 lines (38 lines over 800 limit)
**Recommendation**: OPTIMIZATION PREFERRED, SPLIT AS CONTINGENCY
**Reason**: Temporary placeholder code accounts for ~20 lines, bringing actual implementation to ~818 lines

## Current Size Analysis

### Detailed Breakdown (from line-counter.sh)
```
pkg/oci/cache/key_calculator.go |  49 lines
pkg/oci/cache/layer_db.go       | 267 lines  
pkg/oci/cache/manager.go        | 344 lines
pkg/oci/cache/strategies.go     | 178 lines
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Total:                             838 lines
```

## Option A: OPTIMIZATION (RECOMMENDED)

### Identified Optimization Opportunities
1. **Remove Temporary Interface Placeholders** (~20 lines)
   - Lines 13-33 in manager.go contain temporary Layer and CacheManager definitions
   - These will be replaced with imports from effort1-contracts
   - Immediate reduction: 20 lines → **818 lines**

2. **Additional Optimizations** (~20-30 lines possible)
   - Consolidate error messages to constants
   - Combine similar struct definitions
   - Use type aliases where appropriate
   - Target: Get below 790 lines for safety margin

### Optimization Instructions for SW Engineer
```go
// REMOVE these temporary definitions from manager.go (lines 13-33):
type Layer struct { ... }  // Delete - will import from effort1
type CacheManager interface { ... }  // Delete - will import from effort1

// REPLACE WITH:
import "github.com/cnoe-io/idpbuilder/pkg/oci/api"
// Use api.Layer and api.CacheManager throughout
```

### Expected Result After Optimization
- Remove 20 lines of placeholders: 838 → 818 lines
- Minor consolidations: 818 → ~790 lines
- **COMPLIANT with 800-line limit**
- Maintains 10-line safety buffer

## Option B: SPLIT PLAN (IF OPTIMIZATION INSUFFICIENT)

If optimization cannot achieve <800 lines, implement the following split:

### Split Strategy Overview
- **Total Splits Required**: 2
- **Split Approach**: Functional separation
- **Each Split Target**: <700 lines for safety

### SPLIT-001: Core Cache Management
**Size**: ~450 lines
**Branch**: `idpbuidler-oci-mgmt/phase2/wave2/effort3-cache-split-001`

**Files**:
```yaml
pkg/oci/cache/
  manager.go         # 344 lines - Main CacheManager implementation
  key_calculator.go  # 49 lines - Cache key generation
  interfaces.go      # ~60 lines - Shared types and interfaces
```

**Functionality**:
- Core cache manager implementation
- Cache key calculation logic
- Basic CRUD operations for layers
- Statistics tracking

**Dependencies**: 
- effort1-contracts (for api.Layer and api.CacheManager interfaces)

### SPLIT-002: Storage and Eviction
**Size**: ~390 lines
**Branch**: `idpbuidler-oci-mgmt/phase2/wave2/effort3-cache-split-002`

**Files**:
```yaml
pkg/oci/cache/
  layer_db.go     # 267 lines - Layer database with metadata
  strategies.go   # 178 lines - Cache eviction strategies
```

**Functionality**:
- Layer database implementation
- Metadata management and indexing
- All eviction strategies (LRU, TTL, Size)
- Reference counting logic

**Dependencies**:
- Split-001 (imports shared interfaces)
- effort1-contracts

### Split Boundaries
```
┌─────────────────────────────────────────────┐
│ SPLIT-001: Core Cache (450 lines)          │
│ - CacheManager implementation               │
│ - Cache key calculation                     │
│ - Public API surface                        │
├─────────────────────────────────────────────┤
│ SPLIT-002: Storage & Eviction (390 lines)  │
│ - Layer database                            │
│ - Eviction strategies                       │
│ - Metadata indexing                         │
└─────────────────────────────────────────────┘
```

### Deduplication Verification
| Component | Split-001 | Split-002 |
|-----------|-----------|-----------|
| manager.go | ✅ | ❌ |
| key_calculator.go | ✅ | ❌ |
| layer_db.go | ❌ | ✅ |
| strategies.go | ❌ | ✅ |
| interfaces.go (new) | ✅ | imports only |

## Recommendation Decision Matrix

### Choose OPTIMIZATION if:
- ✅ Temporary placeholders can be removed (confirmed: yes)
- ✅ Total reduction achieves <800 lines (expected: yes)
- ✅ No new features needed immediately (current scope complete)

### Choose SPLIT if:
- ❌ Optimization cannot achieve <800 lines
- ❌ Additional features required before merge
- ❌ Review identifies need for more code

## Next Steps

### For SW Engineer:
1. **ATTEMPT OPTIMIZATION FIRST**:
   ```bash
   # Remove temporary interface definitions
   # Replace with imports from effort1-contracts
   # Run line counter to verify
   cd /home/vscode/workspaces/idpbuilder-oci-mgmt/efforts/phase2/wave2/effort3-cache
   /home/vscode/workspaces/idpbuilder-oci-mgmt/tools/line-counter.sh
   ```

2. **If optimization achieves <800 lines**:
   - Commit optimizations
   - Request final review
   - Mark effort complete

3. **If optimization insufficient**:
   - Implement Split-001 first
   - Then implement Split-002
   - Each split gets separate review

### For Orchestrator:
1. Direct SW Engineer to attempt optimization
2. Monitor line count after optimization
3. If still over limit, initiate split protocol
4. Ensure sequential split execution (not parallel)

## Risk Mitigation
- **Primary Risk**: Optimization might not achieve sufficient reduction
- **Mitigation**: Split plan ready as contingency
- **Secondary Risk**: Dependencies on effort1-contracts not ready
- **Mitigation**: Keep temporary interfaces until effort1 complete

## Compliance Verification
- Target: <800 lines per effort/split
- Safety margin: Target 700-750 lines
- Measurement tool: `/home/vscode/workspaces/idpbuilder-oci-mgmt/tools/line-counter.sh`
- Validation: Must pass before merge

## Document Metadata
- **Created**: 2025-08-26 15:27:00 UTC
- **Author**: Code Reviewer Agent (SPLIT_PLANNING state)
- **Decision**: OPTIMIZATION PREFERRED
- **Contingency**: SPLIT PLAN PROVIDED