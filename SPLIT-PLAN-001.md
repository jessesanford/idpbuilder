# SPLIT-PLAN-001.md

## Split 001 of 2: Core Optimizer with Optimized Analyzer
**Planner**: Code Reviewer code-reviewer-533480 (same for ALL splits)
**Parent Effort**: effort2-optimizer

<!-- ⚠️ ORCHESTRATOR METADATA PLACEHOLDER - DO NOT REMOVE ⚠️ -->
<!-- The orchestrator will add infrastructure metadata below: -->
<!-- WORKING_DIRECTORY, BRANCH, REMOTE, BASE_BRANCH, etc. -->
<!-- SW Engineers MUST read this metadata to navigate to the correct directory -->
<!-- END PLACEHOLDER -->

## Boundaries (⚠️ CRITICAL: All splits MUST reference SAME effort!)
- **Previous Split**: None (first split of THIS effort)
  - Path: N/A (this is Split 001)
  - Branch: N/A
- **This Split**: Split 001 of phase2/wave2/effort2-optimizer
  - Path: efforts/phase2/wave2/effort2-optimizer/split-001/
  - Branch: idpbuidler-oci-mgmt/phase2/wave2/effort2-optimizer-split-001
- **Next Split**: Split 002 of phase2/wave2/effort2-optimizer
  - Path: efforts/phase2/wave2/effort2-optimizer/split-002/
  - Branch: idpbuidler-oci-mgmt/phase2/wave2/effort2-optimizer-split-002
- **File Boundaries**:
  - This Split: analyzer.go (optimized), optimizer.go (fixed with stubs)
  - Next Split: executor.go (full impl), graph.go (full impl)

## Files in This Split (EXCLUSIVE - no overlap with other splits)
- `pkg/oci/optimizer/analyzer.go` - Optimize from 496 to ~380 lines
- `pkg/oci/optimizer/optimizer.go` - Fix and reduce from 368 to ~320 lines
- **Total Target**: 700 lines (with buffer)

## Current Issues to Fix
1. **Build Broken**: optimizer.go references undefined NewExecutor() and NewGraphBuilder()
2. **Size Issue**: analyzer.go is too verbose (496 lines)
3. **Missing Types**: Executor and GraphBuilder don't exist

## Implementation Instructions

### Step 1: Fix Compilation (CRITICAL - Do This First!)
Add these stub types to optimizer.go to fix the build:

```go
// Executor manages parallel stage execution (stub for split-001)
type Executor struct {
    workers int
}

// NewExecutor creates a new executor (full implementation in split-002)
func NewExecutor(workers int) *Executor {
    return &Executor{workers: workers}
}

// GraphBuilder constructs dependency graphs (stub for split-001)
type GraphBuilder struct{}

// NewGraphBuilder creates a new graph builder (full implementation in split-002)
func NewGraphBuilder() *GraphBuilder {
    return &GraphBuilder{}
}
```

### Step 2: Optimize analyzer.go (~116 line reduction)

#### 2.1: Extract Constants (save ~20 lines)
```go
// Move to top of file
const (
    defaultBuildTimeSeconds = 30
    minCacheableSize = 10
    maxStageNameLength = 50
    defaultComplexityScore = 1.0
)

var (
    stageRegex    = regexp.MustCompile(`(?i)^FROM\s+([^\s]+)(?:\s+AS\s+([^\s]+))?`)
    copyFromRegex = regexp.MustCompile(`(?i)COPY\s+--from=([^\s]+)`)
    argRegex      = regexp.MustCompile(`(?i)^ARG\s+([^=\s]+)(?:=([^\s]+))?`)
)
```

#### 2.2: Consolidate Helper Functions (save ~30 lines)
```go
// Combine validation functions
func (a *Analyzer) validateStage(stage *api.Stage) error {
    if stage.Name == "" {
        return fmt.Errorf("stage name cannot be empty")
    }
    if len(stage.Name) > maxStageNameLength {
        return fmt.Errorf("stage name too long: %d > %d", len(stage.Name), maxStageNameLength)
    }
    // Combine other validations
    return nil
}

// Merge similar parsing helpers
func (a *Analyzer) parseStageInfo(line string) (*stageInfo, error) {
    // Combine parseStages subfunctions
}
```

#### 2.3: Simplify Error Handling (save ~40 lines)
```go
// Create error helper
func wrapError(err error, msg string) error {
    if err == nil {
        return nil
    }
    return fmt.Errorf("%s: %w", msg, err)
}

// Use throughout:
// Before: 
if err != nil {
    return nil, fmt.Errorf("failed to parse stages: %w", err)
}
// After:
if err := wrapError(parseErr, "failed to parse stages"); err != nil {
    return nil, err
}
```

#### 2.4: Remove Verbose Comments and Redundant Code (save ~26 lines)
- Remove obvious comments
- Consolidate duplicate validation loops
- Simplify string building operations

### Step 3: Optimize optimizer.go (~48 line reduction)

#### 3.1: Simplify Method Implementations
- Combine similar error checking patterns
- Extract common validation logic
- Reduce verbose struct initializations

#### 3.2: Add Interface Methods for Missing Components
```go
// Add these methods to Executor stub
func (e *Executor) Execute(ctx context.Context, stages []*api.Stage) error {
    // TODO: Implement in split-002
    return fmt.Errorf("executor not implemented - see split-002")
}

// Add these methods to GraphBuilder stub
func (g *GraphBuilder) Build(stages []*api.Stage) (*api.DependencyGraph, error) {
    // TODO: Implement in split-002
    return nil, fmt.Errorf("graph builder not implemented - see split-002")
}
```

## Size Management
- **Current Total**: 864 lines (496 + 368)
- **Target Total**: 700 lines (380 + 320)
- **Reduction Needed**: 164 lines
- **Buffer**: 100 lines (targeting 650, limit is 700)

## Measurement Command
```bash
cd efforts/phase2/wave2/effort2-optimizer/split-001
$PROJECT_ROOT/tools/line-counter.sh -c idpbuidler-oci-mgmt/phase2/wave2/effort2-optimizer-split-001
```

## Test Requirements
- All existing tests must still pass
- Add tests for stub implementations (return expected errors)
- Maintain 85% coverage on optimized code

## Dependencies
- None (foundational split)
- Defines interfaces for split-002

## Integration with Split 002
Split 002 will:
1. Import types from this split
2. Replace stub implementations with real ones
3. Add executor.go and graph.go files

## Success Criteria
- ✅ Code compiles without errors
- ✅ Size ≤700 lines (measured with line-counter.sh)
- ✅ All optimizer interface methods present
- ✅ Tests pass with coverage ≥85%
- ✅ Stub implementations documented with TODO comments

## Validation Checklist
- [ ] Build succeeds (no undefined types)
- [ ] analyzer.go reduced by ~116 lines
- [ ] optimizer.go reduced by ~48 lines
- [ ] Total size ≤700 lines
- [ ] All tests passing
- [ ] Interfaces well-defined for split-002

---
**Note**: Focus on optimization without breaking functionality. The stubs are temporary and will be replaced in split-002.