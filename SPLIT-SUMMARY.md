# Split Plan: error-reporting-types

## Original Effort
- **Size**: 892 lines (92 lines OVER the 800 line limit!)
- **Reason for Split**: Size compliance - exceeds mandatory 800 line limit
- **Measurement Tool**: /home/vscode/workspaces/idpbuilder-oci-mgmt/tools/line-counter.sh

## Split Strategy

### Split 1: Error Types and Codes (error-types-core)
- **Files**:
  - `pkg/errors/codes.go` (148 lines)
  - `pkg/errors/types.go` (191 lines)
  - `pkg/errors/constants.go` (65 lines)
  - `pkg/errors/codes_test.go` (323 lines)
- **Functionality**: Core error type definitions, error codes, categories, and base error implementation
- **Dependencies**: None (foundational split)
- **Estimated Size**: ~404 lines implementation + tests
- **Total with tests**: ~727 lines (UNDER 800 limit)

### Split 2: Progress Tracking (progress-tracking)
- **Files**:
  - `pkg/progress/types.go` (358 lines)
  - `pkg/progress/constants.go` (68 lines)
  - `pkg/errors/types_test.go` (275 lines) - error types integration tests
- **Functionality**: Progress tracking types, event management, and reporter interfaces
- **Dependencies**: Imports from Split 1 (error types for progress error handling)
- **Estimated Size**: ~426 lines implementation + tests  
- **Total with tests**: ~701 lines (UNDER 800 limit)

### Split 3: Progress Tests (progress-tests)
- **Files**:
  - `pkg/progress/types_test.go` (373 lines)
- **Functionality**: Comprehensive progress tracking tests
- **Dependencies**: Imports from Split 2 (progress types)
- **Estimated Size**: ~373 lines
- **Total**: ~373 lines (WELL UNDER 800 limit)

## Implementation Order
1. **Split 1 (error-types-core)**: Implement first as foundational types
   - Branch: `phase1/wave1/error-reporting-types-split-1`
   - Create error codes, categories, and base error types
   - Include unit tests for error functionality
   
2. **Split 2 (progress-tracking)**: Implement after Split 1 merges
   - Branch: `phase1/wave1/error-reporting-types-split-2`  
   - Depends on error types from Split 1
   - Create progress tracking implementation
   - Include error integration tests
   
3. **Split 3 (progress-tests)**: Implement after Split 2 merges
   - Branch: `phase1/wave1/error-reporting-types-split-3`
   - Add comprehensive progress tracking tests
   - Validate full functionality

## Split Branch Strategy
- Original branch: `phase1/wave1/error-reporting-types`
- Split 1: `phase1/wave1/error-reporting-types-split-1`
- Split 2: `phase1/wave1/error-reporting-types-split-2`
- Split 3: `phase1/wave1/error-reporting-types-split-3`
- Integration: Merge all splits back to `phase1/wave1/error-reporting-types`

## Key Instructions for SW Engineers

### For Split 1:
1. Create new branch from main
2. Copy only the specified files for Split 1
3. Ensure all error types and codes are complete
4. Run tests to verify functionality
5. Measure with line counter before commit

### For Split 2:
1. Create branch from Split 1 after it merges
2. Add progress tracking implementation
3. Import error types from Split 1
4. Include error integration tests
5. Verify size compliance

### For Split 3:
1. Create branch from Split 2 after it merges
2. Add comprehensive progress tests
3. Ensure full test coverage
4. Validate all progress functionality

## Size Compliance Verification
Each split MUST be measured using:
```bash
/home/vscode/workspaces/idpbuilder-oci-mgmt/tools/line-counter.sh [branch-name]
```

## Critical Rules
- Each split MUST be under 800 lines
- Splits MUST be implemented sequentially, not in parallel
- Each split gets full code review before next split starts
- If any split still exceeds limit, it must be split recursively