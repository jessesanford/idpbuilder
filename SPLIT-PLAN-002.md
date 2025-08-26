# SPLIT-PLAN-002.md
## Split 002 of 2: Dockerfile Parser and Dependencies Completion
**Planner**: Code Reviewer (CREATE_SPLIT_PLAN state)
**Parent Effort**: dockerfile-builder
**Branch Strategy**: idpbuilder-oci-mgmt/phase2/wave1/dockerfile-builder-split-002

<!-- ORCHESTRATOR METADATA PLACEHOLDER - DO NOT REMOVE -->
<!-- The orchestrator will add infrastructure metadata below: -->
<!-- WORKING_DIRECTORY, BRANCH, REMOTE, BASE_BRANCH, etc. -->
<!-- SW Engineers MUST read this metadata to navigate to the correct directory -->
<!-- END PLACEHOLDER -->

## Split Boundaries (CRITICAL: All splits reference SAME effort!)
- **Previous Split**: Split 001 of phase2/wave1/dockerfile-builder
  - Path: efforts/phase2/wave1/dockerfile-builder/split-001/
  - Branch: idpbuilder-oci-mgmt/phase2/wave1/dockerfile-builder-split-001
  - Summary: Implemented core builder service and layer management (750 lines)
- **This Split**: Split 002 of phase2/wave1/dockerfile-builder
  - Path: efforts/phase2/wave1/dockerfile-builder/split-002/
  - Branch: idpbuilder-oci-mgmt/phase2/wave1/dockerfile-builder-split-002
- **Next Split**: None (final split of THIS effort)
  - This completes the dockerfile-builder effort

## Files in This Split (EXCLUSIVE - no overlap with Split-001)
| File | Lines | Purpose |
|------|-------|---------|
| pkg/oci/build/dockerfile.go | 290 | Dockerfile parsing and validation |
| pkg/oci/build/dockerfile_test.go | 22 | Unit tests for parser |
| pkg/oci/build/go.sum | 391 | Lines 144-534 of dependency checksums |
| Documentation updates | ~100 | Integration notes and completion |
| **TOTAL** | **~612** | Within 800-line limit |

## Functionality Scope
This split completes the dockerfile-builder effort with:
1. **Dockerfile Parser** (`dockerfile.go`):
   - Parses Dockerfile syntax
   - Validates Dockerfile instructions
   - Supports multi-stage builds
   - Extracts build metadata (base images, args, labels)
   - Handles build context validation

2. **Parser Tests** (`dockerfile_test.go`):
   - Unit tests for parser functionality
   - Edge case testing
   - Multi-stage build validation

3. **Dependency Completion** (`go.sum` completion):
   - Completes the dependency checksum file
   - Adds remaining 391 lines to go.sum
   - Ensures module integrity

## Dependencies on Split-001
- **Required from Split-001**:
  - Working `Builder` struct that will use the parser
  - `LayerManager` for processing parsed instructions
  - Base go.mod and initial go.sum (lines 1-143)
  - Module structure already established

- **Integration Points**:
  - DockerfileParser will be called by Builder.Build()
  - ParsedDockerfile will guide layer creation
  - Instructions will map to layer operations

## Implementation Instructions for SW Engineer

### 1. Setup Split Working Directory
```bash
# Start from Split-001's completed work
cd efforts/phase2/wave1/dockerfile-builder/split-002
git checkout idpbuilder-oci-mgmt/phase2/wave1/dockerfile-builder-split-001
git checkout -b idpbuilder-oci-mgmt/phase2/wave1/dockerfile-builder-split-002
```

### 2. Retrieve Split-001 Implementation
Before starting Split-002, ensure you have:
- All files from Split-001 (builder.go, layers.go, tests)
- Working go.mod
- Partial go.sum (143 lines)

### 3. File Implementation Order
1. Implement `dockerfile.go` with parser logic
2. Add `dockerfile_test.go` with comprehensive tests
3. Complete the go.sum file (append lines 144-534)
4. Update integration points in builder.go if needed
5. Update documentation

### 4. Key Implementation Points
- DockerfileParser must handle all standard Dockerfile commands
- Support multi-stage builds with proper stage references
- Validate syntax and report clear errors
- Extract metadata for build optimization
- Thread-safe if used concurrently

### 5. go.sum Completion
```bash
# After adding dockerfile.go dependencies
go mod tidy

# This will regenerate complete go.sum
# Verify it matches expected ~534 lines total

# The first 143 lines should match Split-001's version
head -143 go.sum > go.sum.split001
# Compare with Split-001's go.sum to ensure consistency

# Lines 144-534 are new for this split
tail -n +144 go.sum > go.sum.split002
wc -l go.sum.split002  # Should be ~391 lines
```

### 6. Integration Testing
```bash
# Test that parser integrates with builder
go test ./pkg/oci/build/... -v

# Ensure all components work together
# Builder can use DockerfileParser
# Parser output feeds LayerManager
```

### 7. Size Verification
```bash
# Before committing, verify size
/home/vscode/workspaces/idpbuilder-oci-mgmt/tools/line-counter.sh --base idpbuilder-oci-mgmt/phase2/wave1/dockerfile-builder-split-001
# Must show ≤612 lines for this split
```

## Quality Checklist
- [ ] dockerfile.go compiles without errors
- [ ] All parser tests pass
- [ ] Integration with Split-001 code verified
- [ ] Size under 612 lines (verified with line-counter.sh)
- [ ] go.sum completed (534 total lines)
- [ ] Parser handles all Dockerfile commands
- [ ] Multi-stage build support working
- [ ] Error messages are clear and helpful
- [ ] Documentation updated

## Final Integration Notes
After Split-002 completion:
1. Both splits will be merged back to the parent branch
2. The complete dockerfile-builder effort will have:
   - Full OCIBuildService implementation
   - Complete Dockerfile parsing
   - Layer management with caching
   - All tests passing
   - Total ~1362 lines (but delivered in compliant chunks)

## Testing the Complete Integration
```bash
# After both splits are merged
cd efforts/phase2/wave1/dockerfile-builder
go test ./... -v -cover

# Verify complete functionality
# - Can parse complex Dockerfiles
# - Can build images using buildah
# - Layer caching works
# - All Phase 1 interfaces satisfied
```

## Commit Message Template
```
feat(dockerfile-builder): add dockerfile parser and complete dependencies (split-002)

- Implement comprehensive Dockerfile parser
- Support multi-stage builds
- Complete module dependencies (go.sum)
- Add parser unit tests
- Integrate with builder from split-001

Part 2 of 2 - Parser and completion
Size: 612 lines (under 800 limit)
Completes dockerfile-builder effort
```

## Post-Split Verification
Once both splits are complete and merged:
1. Total functionality matches original IMPLEMENTATION-PLAN.md
2. All tests pass
3. No code duplication between splits
4. Clean integration with Phase 1 APIs
5. Ready for Phase 2 Wave 1 integration