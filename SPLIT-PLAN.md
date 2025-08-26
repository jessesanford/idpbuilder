# Split Plan: Multi-Stage Build Optimizer

## 🚨 SIZE VIOLATION DETECTED

**Review Date**: 2025-08-26
**Reviewer**: Code Reviewer Agent (ID: code-reviewer-533480)
**Current Size**: 864 lines (EXCEEDS LIMIT)
**Size Limit**: 800 lines
**Overage**: 64 lines (8% over limit)
**Measurement Tool**: /home/vscode/workspaces/idpbuilder-oci-mgmt/tools/line-counter.sh

## Current Implementation Analysis

### File Breakdown
| File | Lines | Status | Purpose |
|------|-------|--------|---------|
| pkg/oci/optimizer/analyzer.go | 496 | Implemented | Stage dependency analysis |
| pkg/oci/optimizer/optimizer.go | 368 | Partially Implemented | Main optimizer (references missing types) |
| pkg/oci/optimizer/executor.go | 0 | Not Implemented | Parallel execution engine |
| pkg/oci/optimizer/graph.go | 0 | Not Implemented | Dependency graph builder |
| **Total** | **864** | **Broken Build** | **Missing critical components** |

### Critical Issues Found
1. **Compilation Failure**: optimizer.go references `NewExecutor()` and `NewGraphBuilder()` which don't exist
2. **Size Violation**: Total implementation exceeds 800 line limit by 64 lines
3. **Missing Components**: executor and graph packages are referenced but not implemented

## Split Strategy: REDUCTION AND COMPLETION

Given the current state (partially implemented with missing components), the optimal strategy is:
1. **REDUCE** analyzer.go through optimization
2. **FIX** compilation issues in Split 001
3. **IMPLEMENT** missing components in Split 002

## Split Configuration

### Split 001: Core Optimizer with Analyzer (Target: 700 lines)
**Description**: Core optimization framework with complete analyzer
**Strategy**: Optimize analyzer.go and complete optimizer.go with stub implementations

**Files**:
- `pkg/oci/optimizer/analyzer.go` - Reduce from 496 to ~380 lines through:
  - Combine similar helper functions
  - Extract constants and regex patterns
  - Simplify verbose error handling
  - Remove redundant validation
- `pkg/oci/optimizer/optimizer.go` - Fix and complete at ~320 lines:
  - Add stub Executor and GraphBuilder types
  - Implement core optimizer logic
  - Add proper error handling

**Size Calculation**:
- Optimized analyzer.go: ~380 lines (reduced by ~116 lines)
- Fixed optimizer.go: ~320 lines (reduced by ~48 lines)
- **Total Split 001**: ~700 lines ✅

### Split 002: Execution and Graph Components (Target: 350 lines)
**Description**: Parallel execution engine and dependency graph builder
**Strategy**: Implement the missing components that are currently breaking the build

**Files**:
- `pkg/oci/optimizer/executor.go` - Implement ~180 lines:
  - Worker pool management
  - Stage scheduling
  - Result collection
- `pkg/oci/optimizer/graph.go` - Implement ~120 lines:
  - Dependency graph construction
  - Topological sorting
  - Critical path analysis
- Test stubs: ~50 lines

**Size Calculation**:
- executor.go: ~180 lines
- graph.go: ~120 lines
- Test stubs: ~50 lines
- **Total Split 002**: ~350 lines ✅

## Implementation Instructions for Orchestrator

### Phase 1: Create Split Infrastructure
```bash
# Create split directories
mkdir -p efforts/phase2/wave2/effort2-optimizer/split-001
mkdir -p efforts/phase2/wave2/effort2-optimizer/split-002

# Create branches
git checkout -b idpbuidler-oci-mgmt/phase2/wave2/effort2-optimizer-split-001
git checkout -b idpbuidler-oci-mgmt/phase2/wave2/effort2-optimizer-split-002
```

### Phase 2: Execute Split 001 (SEQUENTIAL - BLOCKING)
1. **Spawn SW Engineer** for split-001 with optimization focus
2. **Optimize analyzer.go**:
   - Consolidate helper functions
   - Simplify error handling patterns
   - Extract common patterns to constants
3. **Fix optimizer.go**:
   - Add stub implementations for Executor and GraphBuilder
   - Ensure compilation succeeds
4. **Measure and Verify**: Must be ≤700 lines
5. **Review**: Code Reviewer validates optimization

### Phase 3: Execute Split 002 (AFTER Split 001 Complete)
1. **Spawn SW Engineer** for split-002 
2. **Implement executor.go**: Complete parallel execution engine
3. **Implement graph.go**: Complete dependency graph builder
4. **Test Integration**: Ensure all components work together
5. **Measure and Verify**: Must be ≤400 lines
6. **Review**: Code Reviewer validates functionality

## Optimization Opportunities for Split 001

### analyzer.go Reduction Strategies
1. **Combine Similar Functions** (save ~30 lines):
   - Merge `parseStages` helper functions
   - Combine validation methods

2. **Simplify Error Handling** (save ~40 lines):
   ```go
   // Before: Multiple error checks
   if err != nil {
       return nil, fmt.Errorf("failed to parse: %w", err)
   }
   
   // After: Use error helper
   return checkError(result, err, "failed to parse")
   ```

3. **Extract Constants** (save ~20 lines):
   ```go
   // Move regex patterns and constants to top
   const (
       defaultBuildTime = 30 * time.Second
       minCacheableSize = 10
   )
   ```

4. **Consolidate Duplicate Logic** (save ~26 lines):
   - Merge similar dependency checking loops
   - Combine stage validation logic

### optimizer.go Fixes
1. **Add Stub Types** (minimum viable):
   ```go
   type Executor struct{}
   func NewExecutor(workers int) *Executor { return &Executor{} }
   
   type GraphBuilder struct{}
   func NewGraphBuilder() *GraphBuilder { return &GraphBuilder{} }
   ```

2. **Placeholder Methods**: Implement interface with TODOs for split-002

## Risk Mitigation

### Risk: Further Size Overrun
**Mitigation**: Aggressive optimization in split-001, target 650 lines with 50-line buffer

### Risk: Breaking Functionality During Optimization
**Mitigation**: Preserve all tests, validate interface compliance

### Risk: Integration Issues Between Splits
**Mitigation**: Define clear interfaces in split-001, implement in split-002

## Success Criteria
- ✅ Split 001 compiles successfully and is ≤700 lines
- ✅ Split 002 implements missing components and is ≤400 lines
- ✅ Combined functionality matches original plan
- ✅ All tests pass
- ✅ No duplicate code between splits

## Validation Checklist
- [ ] Split 001 size verified with line-counter.sh
- [ ] Split 001 compiles without errors
- [ ] Split 002 size verified with line-counter.sh
- [ ] Integration between splits tested
- [ ] Original functionality preserved
- [ ] Test coverage maintained at 85%+

## Next Steps for Orchestrator
1. **IMMEDIATE**: Stop current effort implementation
2. **CREATE**: Split working directories
3. **SPAWN**: SW Engineer for split-001 with this plan
4. **SEQUENTIAL**: Complete split-001 before starting split-002
5. **INTEGRATE**: Merge splits back after both complete

---

**Note**: This is a MANDATORY split due to size violation. The implementation cannot proceed without splitting. The current code is also broken (references undefined types), making this an ideal splitting point.