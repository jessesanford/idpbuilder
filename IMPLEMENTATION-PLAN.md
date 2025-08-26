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
    "fmt"
    "github.com/jessesanford/idpbuilder/pkg/oci/api"
)

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