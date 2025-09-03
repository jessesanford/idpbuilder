# Split Plan 001A: pkg/builder Package Only

## Overview
This sub-split focuses solely on the pkg/builder package implementation.

## Estimated Size
**Target**: 572 lines (within 800 line limit)

## Scope
### Included Files
- `pkg/builder/builder.go` - Core builder implementation (~149 lines)
- `pkg/builder/options.go` - Build options and configuration (~305 lines)
- `pkg/builder/builder_test.go` - Unit tests (~118 lines)

### Total Expected: ~572 lines

## Implementation Instructions
1. Create ONLY the pkg/builder directory
2. Implement the builder interface and options
3. Add comprehensive unit tests
4. Use feature flag ENABLE_CORE_BUILDER to guard functionality
5. Ensure builds cleanly and tests pass

## Dependencies
- Base from main branch (first sub-split)
- No dependencies on pkg/imagebuilder

## Success Criteria
- ✅ Under 800 lines total
- ✅ All tests pass
- ✅ Feature flag properly guards functionality
- ✅ Builds independently
- ✅ Can be merged without imagebuilder package

## Base Branch
- Base from: main (per R308 for first split)
## 🚨 SPLIT INFRASTRUCTURE METADATA (Added by Orchestrator)
**WORKING_DIRECTORY**: /home/vscode/workspaces/idpbuilder-oci-go-cr/efforts/phase2/wave1/go-containerregistry-image-builder-SPLIT-001A
**BRANCH**: idpbuilder-oci-go-cr/phase2/wave1/go-containerregistry-image-builder--split-001a
**REMOTE**: origin/idpbuilder-oci-go-cr/phase2/wave1/go-containerregistry-image-builder--split-001a
**BASE_BRANCH**: main
**SPLIT_NUMBER**: 001A
**CREATED_AT**: 2025-09-03 22:28:00

### SW Engineer Instructions
1. READ this metadata FIRST
2. cd to WORKING_DIRECTORY above
3. Verify branch matches BRANCH above
4. ONLY THEN proceed with implementation
5. Create ONLY pkg/builder package (not imagebuilder)
6. Stay under 800 lines total
