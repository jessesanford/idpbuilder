# Implementation Plan: Multi-Stage Build Optimizer - Split 001

## <¯ Effort Overview
**Effort ID**: effort2-optimizer-split-001
**Target Size**: 700 lines MAXIMUM
**Purpose**: Core optimizer with optimized analyzer (fixing build issues)

## =¨ CRITICAL REQUIREMENTS
1. **SIZE LIMIT**: 700 lines HARD LIMIT (use line-counter.sh)
2. **MUST COMPILE**: Fix all compilation errors
3. **OPTIMIZE AGGRESSIVELY**: Reduce analyzer.go from 496 to ~380 lines

## =Á Files to Implement

### 1. pkg/oci/optimizer/analyzer.go (~380 lines)
**Current**: 496 lines
**Target**: ~380 lines (REDUCE BY 116+ LINES)

**Optimization Strategy**:
1. **Extract Constants** (save ~20 lines):
   ```go
   const (
       defaultBuildTime = 30 * time.Second
       minCacheableSize = 10
       maxParallelStages = 10
   )
   ```

2. **Combine Similar Functions** (save ~30 lines):
   - Merge parseStages helper functions
   - Combine validation methods into one

3. **Simplify Error Handling** (save ~40 lines):
   ```go
   // Create error helper
   func wrapErr(result interface{}, err error, msg string) (interface{}, error) {
       if err != nil {
           return nil, fmt.Errorf("%s: %w", msg, err)
       }
       return result, nil
   }
   ```

4. **Remove Redundant Validation** (save ~26 lines):
   - Consolidate duplicate dependency checks
   - Merge stage validation logic

### 2. pkg/oci/optimizer/optimizer.go (~320 lines)
**Current**: 368 lines (broken - references missing types)
**Target**: ~320 lines (FIXED AND WORKING)

**Fix Strategy**:
1. **Add Stub Types** to make it compile:
   ```go
   // Stub types for split-002 implementation
   type Executor struct {
       workers int
   }
   
   func NewExecutor(workers int) *Executor {
       return &Executor{workers: workers}
   }
   
   func (e *Executor) Execute(stages []Stage) error {
       // TODO: Implement in split-002
       return nil
   }
   
   type GraphBuilder struct{}
   
   func NewGraphBuilder() *GraphBuilder {
       return &GraphBuilder{}
   }
   
   func (g *GraphBuilder) Build(stages []Stage) (*DependencyGraph, error) {
       // TODO: Implement in split-002
       return &DependencyGraph{}, nil
   }
   
   type DependencyGraph struct {
       // TODO: Implement in split-002
   }
   ```

2. **Complete Core Logic**:
   - Fix all compilation errors
   - Ensure interfaces are properly defined
   - Add proper error handling

## =' Implementation Steps

### Step 1: Optimize analyzer.go
1. Extract all constants to top of file
2. Create helper functions to reduce duplication
3. Combine similar validation logic
4. Use table-driven tests where applicable
5. Remove verbose comments

### Step 2: Fix optimizer.go
1. Add stub type definitions
2. Fix all undefined references
3. Ensure compilation succeeds
4. Define clear interfaces for split-002

### Step 3: Verify Size
```bash
/home/vscode/workspaces/idpbuilder-oci-mgmt/tools/line-counter.sh -c branch
# MUST be d700 lines
```

##  Success Criteria
- [ ] analyzer.go reduced to ~380 lines
- [ ] optimizer.go fixed and ~320 lines
- [ ] Total implementation d700 lines
- [ ] Code compiles without errors
- [ ] All tests pass
- [ ] Clear interfaces for split-002

## =¨ Critical Notes
1. **DO NOT EXCEED 700 LINES** - This is a hard limit
2. **MUST COMPILE** - Current code is broken, fix it
3. **PRESERVE FUNCTIONALITY** - Only optimize, don't break features
4. **DOCUMENT INTERFACES** - Split-002 depends on these

## Dependencies for Split-002
The following will be implemented in split-002:
- Full Executor implementation (~180 lines)
- Full GraphBuilder implementation (~120 lines)
- DependencyGraph structure and methods (~50 lines)