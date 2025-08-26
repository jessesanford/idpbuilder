# SPLIT-PLAN-002.md

## Split 002 of 2: Execution and Graph Components
**Planner**: Code Reviewer code-reviewer-533480 (same for ALL splits)
**Parent Effort**: effort2-optimizer

<!-- ⚠️ ORCHESTRATOR METADATA PLACEHOLDER - DO NOT REMOVE ⚠️ -->
<!-- The orchestrator will add infrastructure metadata below: -->
<!-- WORKING_DIRECTORY, BRANCH, REMOTE, BASE_BRANCH, etc. -->
<!-- SW Engineers MUST read this metadata to navigate to the correct directory -->
<!-- END PLACEHOLDER -->

## Boundaries (⚠️ CRITICAL: All splits MUST reference SAME effort!)
- **Previous Split**: Split 001 of phase2/wave2/effort2-optimizer
  - Path: efforts/phase2/wave2/effort2-optimizer/split-001/
  - Branch: idpbuidler-oci-mgmt/phase2/wave2/effort2-optimizer-split-001
  - Summary: Core optimizer with optimized analyzer and stub types
- **This Split**: Split 002 of phase2/wave2/effort2-optimizer
  - Path: efforts/phase2/wave2/effort2-optimizer/split-002/
  - Branch: idpbuidler-oci-mgmt/phase2/wave2/effort2-optimizer-split-002
- **Next Split**: None (final split of THIS effort)

## Files in This Split (EXCLUSIVE - no overlap with split-001)
- `pkg/oci/optimizer/executor.go` - Full implementation (~180 lines)
- `pkg/oci/optimizer/graph.go` - Full implementation (~120 lines)
- Test stubs for new components (~50 lines)
- **Total Target**: 350 lines (with 50-line buffer)

## Dependencies on Split 001
- Import Optimizer type that references Executor and GraphBuilder
- Replace stub implementations from split-001
- Use interfaces defined in split-001

## Implementation Instructions

### Step 1: Implement Executor (~180 lines)

```go
// pkg/oci/optimizer/executor.go
package optimizer

import (
    "context"
    "fmt"
    "sync"
    "time"
    
    "github.com/cnoe-io/idpbuilder/pkg/oci/api"
)

// Executor manages parallel stage execution with worker pools
type Executor struct {
    workers    int
    workQueue  chan *api.Stage
    resultChan chan *api.StageResult
    wg         sync.WaitGroup
    mu         sync.RWMutex
    results    map[string]*api.StageResult
}

// NewExecutor creates a new parallel executor
func NewExecutor(workers int) *Executor {
    if workers <= 0 {
        workers = api.DefaultMaxParallel
    }
    
    return &Executor{
        workers:    workers,
        workQueue:  make(chan *api.Stage, workers*2),
        resultChan: make(chan *api.StageResult, workers*2),
        results:    make(map[string]*api.StageResult),
    }
}

// Execute runs stages in parallel based on dependencies
func (e *Executor) Execute(ctx context.Context, stages []*api.Stage, deps map[string][]string) ([]*api.StageResult, error) {
    if len(stages) == 0 {
        return nil, fmt.Errorf("no stages to execute")
    }
    
    // Start worker pool
    ctx, cancel := context.WithCancel(ctx)
    defer cancel()
    
    for i := 0; i < e.workers; i++ {
        e.wg.Add(1)
        go e.worker(ctx)
    }
    
    // Start result collector
    go e.collectResults()
    
    // Schedule stages based on dependencies
    if err := e.scheduleStages(ctx, stages, deps); err != nil {
        return nil, fmt.Errorf("scheduling failed: %w", err)
    }
    
    // Wait for completion
    e.wg.Wait()
    close(e.resultChan)
    
    return e.getOrderedResults(stages), nil
}

// worker processes stages from the queue
func (e *Executor) worker(ctx context.Context) {
    defer e.wg.Done()
    
    for {
        select {
        case <-ctx.Done():
            return
        case stage, ok := <-e.workQueue:
            if !ok {
                return
            }
            e.processStage(ctx, stage)
        }
    }
}

// processStage executes a single stage
func (e *Executor) processStage(ctx context.Context, stage *api.Stage) {
    start := time.Now()
    
    result := &api.StageResult{
        StageName: stage.Name,
        StartTime: start,
        Status:    api.BuildStatusRunning,
    }
    
    // Simulate stage execution (replace with actual build logic)
    select {
    case <-ctx.Done():
        result.Status = api.BuildStatusCanceled
        result.Error = ctx.Err()
    case <-time.After(time.Duration(stage.EstimatedDuration)):
        result.Status = api.BuildStatusSuccess
        result.Duration = time.Since(start)
    }
    
    result.EndTime = time.Now()
    e.resultChan <- result
}

// Additional helper methods...
```

### Step 2: Implement GraphBuilder (~120 lines)

```go
// pkg/oci/optimizer/graph.go
package optimizer

import (
    "fmt"
    
    "github.com/cnoe-io/idpbuilder/pkg/oci/api"
)

// GraphBuilder constructs and analyzes dependency graphs
type GraphBuilder struct {
    nodes map[string]*api.Stage
    edges map[string][]string
}

// NewGraphBuilder creates a new dependency graph builder
func NewGraphBuilder() *GraphBuilder {
    return &GraphBuilder{
        nodes: make(map[string]*api.Stage),
        edges: make(map[string][]string),
    }
}

// Build constructs a dependency graph from stages
func (g *GraphBuilder) Build(stages []*api.Stage) (*api.DependencyGraph, error) {
    if len(stages) == 0 {
        return nil, fmt.Errorf("no stages provided")
    }
    
    // Reset state
    g.nodes = make(map[string]*api.Stage)
    g.edges = make(map[string][]string)
    
    // Build node map
    for _, stage := range stages {
        if stage.Name == "" {
            return nil, fmt.Errorf("stage name cannot be empty")
        }
        g.nodes[stage.Name] = stage
    }
    
    // Build edges from dependencies
    for _, stage := range stages {
        for _, dep := range stage.Dependencies {
            if _, exists := g.nodes[dep]; !exists {
                return nil, fmt.Errorf("dependency %s not found for stage %s", dep, stage.Name)
            }
            g.edges[stage.Name] = append(g.edges[stage.Name], dep)
        }
    }
    
    // Detect cycles
    if g.hasCycle() {
        return nil, fmt.Errorf("circular dependency detected")
    }
    
    // Perform topological sort
    order := g.topologicalSort()
    
    // Calculate critical path
    criticalPath := g.findCriticalPath(order)
    
    return &api.DependencyGraph{
        Nodes:        g.nodes,
        Edges:        g.edges,
        Order:        order,
        CriticalPath: criticalPath,
    }, nil
}

// hasCycle detects circular dependencies using DFS
func (g *GraphBuilder) hasCycle() bool {
    visited := make(map[string]bool)
    recStack := make(map[string]bool)
    
    var dfs func(string) bool
    dfs = func(node string) bool {
        visited[node] = true
        recStack[node] = true
        
        for _, neighbor := range g.edges[node] {
            if !visited[neighbor] {
                if dfs(neighbor) {
                    return true
                }
            } else if recStack[neighbor] {
                return true
            }
        }
        
        recStack[node] = false
        return false
    }
    
    for node := range g.nodes {
        if !visited[node] && dfs(node) {
            return true
        }
    }
    
    return false
}

// topologicalSort returns execution order
func (g *GraphBuilder) topologicalSort() []string {
    // Implementation of Kahn's algorithm
    // ...
}

// findCriticalPath identifies the longest execution path
func (g *GraphBuilder) findCriticalPath(order []string) []string {
    // Implementation of critical path method
    // ...
}
```

### Step 3: Create Test Stubs (~50 lines)

```go
// pkg/oci/optimizer/executor_test.go
package optimizer

import (
    "context"
    "testing"
    
    "github.com/cnoe-io/idpbuilder/pkg/oci/api"
)

func TestExecutorParallelExecution(t *testing.T) {
    executor := NewExecutor(2)
    
    stages := []*api.Stage{
        {Name: "base", EstimatedDuration: 100},
        {Name: "app", Dependencies: []string{"base"}},
    }
    
    deps := map[string][]string{
        "app": {"base"},
    }
    
    results, err := executor.Execute(context.Background(), stages, deps)
    if err != nil {
        t.Fatalf("execution failed: %v", err)
    }
    
    if len(results) != 2 {
        t.Errorf("expected 2 results, got %d", len(results))
    }
}

// pkg/oci/optimizer/graph_test.go
package optimizer

import (
    "testing"
    
    "github.com/cnoe-io/idpbuilder/pkg/oci/api"
)

func TestGraphBuilderCycleDetection(t *testing.T) {
    builder := NewGraphBuilder()
    
    stages := []*api.Stage{
        {Name: "a", Dependencies: []string{"b"}},
        {Name: "b", Dependencies: []string{"a"}},
    }
    
    _, err := builder.Build(stages)
    if err == nil {
        t.Error("expected cycle detection error")
    }
}
```

## Size Management
- **Target Total**: 350 lines
- **Breakdown**:
  - executor.go: ~180 lines
  - graph.go: ~120 lines
  - Test stubs: ~50 lines
- **Buffer**: 50 lines (targeting 350, could go up to 400)

## Measurement Command
```bash
cd efforts/phase2/wave2/effort2-optimizer/split-002
$PROJECT_ROOT/tools/line-counter.sh -c idpbuidler-oci-mgmt/phase2/wave2/effort2-optimizer-split-002
```

## Integration Requirements

### With Split 001
1. Import base types from split-001
2. Replace stub NewExecutor() with full implementation
3. Replace stub NewGraphBuilder() with full implementation
4. Ensure interface compatibility

### Testing Integration
```go
// After both splits complete, test full integration:
optimizer := NewOptimizer()
stages, _ := optimizer.AnalyzeStages(dockerfile)
results, _ := optimizer.BuildStages(context.Background(), &api.BuildRequest{
    Stages: stages,
})
```

## Success Criteria
- ✅ Full Executor implementation working
- ✅ Full GraphBuilder implementation working
- ✅ Size ≤350 lines (measured with line-counter.sh)
- ✅ Integration with split-001 successful
- ✅ Parallel execution demonstrable
- ✅ Graph algorithms correct

## Implementation Notes
1. Keep implementations lean - avoid over-engineering
2. Focus on core functionality first
3. Use existing api types from effort1-contracts
4. Ensure thread safety in Executor
5. Validate graph operations thoroughly

## Dependencies
- Requires split-001 to be complete (imports types)
- Uses api package from effort1-contracts

## Validation Checklist
- [ ] executor.go implements parallel execution
- [ ] graph.go implements dependency analysis
- [ ] Total size ≤350 lines
- [ ] All tests passing
- [ ] Integration with split-001 works
- [ ] No compilation errors

---
**Note**: This split completes the optimizer implementation by providing the missing execution and graph components. Ensure compatibility with split-001 interfaces.