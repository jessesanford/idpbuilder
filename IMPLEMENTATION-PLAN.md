<<<<<<< HEAD
<<<<<<< HEAD
# Phase 2 Wave 2 Integration Plan

## INTEGRATION WORKSPACE OVERVIEW
**Purpose**: Merge all Wave 2 efforts into integrated build system
**Integration Branch**: `idpbuilder-oci-mgmt/phase2/wave2-integration`
**Status**: RECOVERY - Completing incomplete integration

## COMPLETED INTEGRATIONS
### Effort 1: Advanced Build Contracts & Interfaces ✓
- **Branch**: `idpbuidler-oci-mgmt/phase2/wave2/effort1-contracts`
- **Status**: ✅ MERGED
- **Files**: Core API contracts in pkg/oci/api/

## INTEGRATIONS IN PROGRESS

### Effort 2: Multi-Stage Build Optimizer (Split Implementation)
- **Split 001**: `idpbuilder-oci-mgmt/phase2/wave2/effort2-optimizer-split-001` (728 lines) ✅ MERGED
  - Core optimizer with optimized analyzer
  - Fixed compilation issues with stub types
- **Split 002**: `idpbuilder-oci-mgmt/phase2/wave2/effort2-optimizer-split-002` (350 lines) 🔄 MERGING
  - Full Executor and GraphBuilder implementation
  - Completes the optimizer implementation

### Effort 3: Cache Manager
- **Branch**: `idpbuidler-oci-mgmt/phase2/wave2/effort3-cache` (834 lines) 🔄 MERGING
- **Purpose**: Layer caching operations and optimization
- **Note**: Originally exceeded size limit (834 > 800) but implementing as-is for integration

### Effort 4: Security Manager (Split Implementation)  
- **Split 001**: `idpbuilder-oci-mgmt/phase2/wave2/effort4-security-split-001` (809 lines) 🔄 MERGING
  - Security orchestration and policy management
  - Note: Slightly over original 762 line estimate
- **Split 002**: `idpbuilder-oci-mgmt/phase2/wave2/effort4-security-split-002` (744 lines)
  - Crypto implementations and additional security features

### Effort 5: Registry Client
- **Branch**: `idpbuidler-oci-mgmt/phase2/wave2/effort5-registry` (793 lines)
- **Purpose**: Registry operations integration

## INTEGRATION STRATEGY

### Sequential Merge Order
1. effort2-optimizer-split-001 (IN PROGRESS)
2. effort2-optimizer-split-002  
3. effort3-cache
4. effort4-security-split-001
5. effort4-security-split-002
6. effort5-registry

### Conflict Resolution
- Merge conflicts expected in IMPLEMENTATION-PLAN.md and work-log.md
- Preserve integration workspace structure
- Maintain effort-specific details in logs
- Ensure no code functionality conflicts

### Integration Verification
- All packages properly structured under `/pkg/oci/`
- Cross-effort dependencies resolved
- No circular dependencies
- Compilation successful
- Tests passing

## FILE STRUCTURE POST-INTEGRATION
```
pkg/
├── oci/
│   ├── api/          # Effort 1 contracts
│   ├── optimizer/    # Effort 2 implementation  
│   ├── cache/        # Effort 3 implementation
│   ├── security/     # Effort 4 implementation
│   └── registry/     # Effort 5 implementation
└── k8s/              # Wave 1 integration
```
=======
# Implementation Plan: Multi-Stage Build Optimizer - Split 002

## 🎯 Effort Overview
**Effort ID**: effort2-optimizer-split-002
**Target Size**: 350 lines MAXIMUM
**Purpose**: Complete Executor and GraphBuilder implementations

## 🚨 CRITICAL REQUIREMENTS
1. **SIZE LIMIT**: 350 lines HARD LIMIT
2. **MUST INTEGRATE**: Work with split-001's interfaces
3. **COMPLETE FUNCTIONALITY**: Implement all stub methods from split-001

## 📁 Files to Implement

### 1. pkg/oci/optimizer/executor.go (~180 lines)
**Purpose**: Parallel execution engine for build stages

**Required Implementation**:
```go
package optimizer

import (
    "context"
    "sync"
    "time"
    "github.com/jessesanford/idpbuilder/pkg/oci/api"
)

type Executor struct {
    workers int
    pool    chan struct{}
}

func NewExecutor(workers int) *Executor {
    return &Executor{
        workers: workers,
        pool:    make(chan struct{}, workers),
    }
}

func (e *Executor) Execute(stages []api.Stage) error {
    // Implement:
    // 1. Worker pool management
    // 2. Stage scheduling based on dependencies
    // 3. Parallel execution with proper synchronization
    // 4. Result collection and error handling
}

func (e *Executor) executeStage(stage api.Stage, wg *sync.WaitGroup) {
    // Implement stage execution logic
}

func (e *Executor) scheduleStages(stages []api.Stage) [][]api.Stage {
    // Group stages by dependency level for parallel execution
}
```

### 2. pkg/oci/optimizer/graph.go (~120 lines)
**Purpose**: Dependency graph builder and analysis

**Required Implementation**:
```go
package optimizer

import (
=======
# Implementation Plan: Security Layer - Split 001 (Security Manager)

## <� Effort Overview
**Effort ID**: effort4-security-split-001
**Target Size**: 386 lines MAXIMUM
**Purpose**: Security orchestration and policy management
**Order**: IMPLEMENT AFTER split-002 (depends on crypto layer)

## =� CRITICAL REQUIREMENTS
1. **SIZE LIMIT**: 386 lines (well under limit)
2. **DEPENDS ON SPLIT-002**: Use Signer/Verifier from split-002
3. **ORCHESTRATION LAYER**: Coordinate security operations

## =� Files to Implement

### 1. pkg/oci/security/manager.go (386 lines)
**Purpose**: Security orchestration, policy enforcement, and coordination

**Core Implementation**:
```go
package security

import (
    "context"
>>>>>>> origin/idpbuilder-oci-mgmt/phase2/wave2/effort4-security-split-001
    "fmt"
    "github.com/jessesanford/idpbuilder/pkg/oci/api"
)

<<<<<<< HEAD
type GraphBuilder struct {
    nodes map[string]*Node
    edges map[string][]string
}

type Node struct {
    Stage    api.Stage
    Level    int
    Visited  bool
    Children []string
}

type DependencyGraph struct {
    Nodes map[string]*Node
    Levels [][]string
}

func NewGraphBuilder() *GraphBuilder {
    return &GraphBuilder{
        nodes: make(map[string]*Node),
        edges: make(map[string][]string),
    }
}

func (g *GraphBuilder) Build(stages []api.Stage) (*DependencyGraph, error) {
    // Implement:
    // 1. Build node map from stages
    // 2. Create edge relationships
    // 3. Perform topological sort
    // 4. Calculate critical path
    // 5. Return structured graph
}

func (g *GraphBuilder) topologicalSort() ([]string, error) {
    // Implement Kahn's algorithm for topological sorting
}

func (g *GraphBuilder) calculateLevels() [][]string {
    // Group nodes by dependency level
}
```

### 3. pkg/oci/optimizer/executor_test.go (~25 lines)
**Purpose**: Basic test stubs

### 4. pkg/oci/optimizer/graph_test.go (~25 lines)
**Purpose**: Basic test stubs

## 🔧 Implementation Steps

### Step 1: Copy API types from split-001
```bash
# Copy the api package from split-001
cp -r ../split-001/pkg/oci/api pkg/oci/
```

### Step 2: Implement executor.go
1. Create worker pool mechanism
2. Implement stage scheduling logic
3. Add parallel execution with sync.WaitGroup
4. Handle errors and timeouts
5. Collect execution results

### Step 3: Implement graph.go
1. Build node structure from stages
2. Create edge relationships from dependencies
3. Implement topological sorting (Kahn's algorithm)
4. Calculate execution levels
5. Identify critical path

### Step 4: Add test stubs
1. Create basic test files
2. Add placeholder test functions
3. Ensure package builds

### Step 5: Verify Size
```bash
/home/vscode/workspaces/idpbuilder-oci-mgmt/tools/line-counter.sh -c branch
# MUST be ≤350 lines
```

## ✅ Success Criteria
- [ ] executor.go implements all required methods (~180 lines)
- [ ] graph.go implements dependency analysis (~120 lines)
- [ ] Test stubs present (~50 lines total)
- [ ] Total implementation ≤350 lines
- [ ] Code compiles with split-001
- [ ] All interfaces satisfied

## 🚨 Critical Notes
1. **DO NOT EXCEED 350 LINES** - Be extremely concise
2. **MUST INTEGRATE** - Use exact types from split-001
3. **FOCUS ON CORE** - Implement minimum viable functionality
4. **NO EXTRAS** - Skip nice-to-haves, focus on essentials

## Integration with Split-001
Split-001 provides:
- `api.Stage`, `api.BuildResult` types
- Stub `Executor` and `GraphBuilder` interfaces
- `Optimizer` that calls these components

Your implementation must satisfy these interfaces exactly.
>>>>>>> origin/idpbuilder-oci-mgmt/phase2/wave2/effort2-optimizer-split-002
=======
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

## =� Critical Notes
1. **DEPENDS ON SPLIT-002**: Must use the crypto interfaces
2. **ORCHESTRATION FOCUS**: Don't re-implement crypto
3. **POLICY LAYER**: Add value on top of basic crypto
4. **SIZE COMFORTABLE**: 386 lines gives plenty of room

## Integration Points
- Uses api.Signer from split-002 for signing
- Uses api.Verifier from split-002 for verification  
- Adds SecurityPolicy enforcement on top
- Provides unified SecurityManager interface
>>>>>>> origin/idpbuilder-oci-mgmt/phase2/wave2/effort4-security-split-001
