# SPLIT-PLAN-001.md
## Split 001 of 2: Core Builder and Layer Management
**Planner**: Code Reviewer (CREATE_SPLIT_PLAN state) 
**Parent Effort**: dockerfile-builder
**Branch Strategy**: idpbuilder-oci-mgmt/phase2/wave1/dockerfile-builder-split-001

<!-- ORCHESTRATOR METADATA PLACEHOLDER - DO NOT REMOVE -->
<!-- The orchestrator will add infrastructure metadata below: -->
<!-- WORKING_DIRECTORY, BRANCH, REMOTE, BASE_BRANCH, etc. -->
<!-- SW Engineers MUST read this metadata to navigate to the correct directory -->
<!-- END PLACEHOLDER -->

## Split Boundaries (CRITICAL: All splits reference SAME effort!)
- **Previous Split**: None (first split of THIS effort)
  - Path: N/A (this is Split 001)
  - Branch: N/A
- **This Split**: Split 001 of phase2/wave1/dockerfile-builder
  - Path: efforts/phase2/wave1/dockerfile-builder/split-001/
  - Branch: idpbuilder-oci-mgmt/phase2/wave1/dockerfile-builder-split-001
- **Next Split**: Split 002 of phase2/wave1/dockerfile-builder
  - Path: efforts/phase2/wave1/dockerfile-builder/split-002/
  - Branch: idpbuilder-oci-mgmt/phase2/wave1/dockerfile-builder-split-002
- **File Boundaries**:
  - This Split: Core builder, layers, and initial dependencies
  - Next Split: Dockerfile parser and remaining dependencies

## Files in This Split (EXCLUSIVE - no overlap with Split-002)
| File | Lines | Purpose |
|------|-------|---------|
| pkg/oci/build/builder.go | 264 | Main builder orchestration implementing OCIBuildService |
| pkg/oci/build/builder_test.go | 33 | Unit tests for builder |
| pkg/oci/build/layers.go | 132 | Layer management and caching |
| pkg/oci/build/layers_test.go | 36 | Unit tests for layer manager |
| pkg/oci/build/go.mod | 142 | Module definition with core dependencies |
| pkg/oci/build/go.sum | 143 | First 143 lines of dependency checksums |
| **TOTAL** | **750** | Within 800-line limit |

## Functionality Scope
This split implements the core building infrastructure:
1. **Builder Service** (`builder.go`):
   - Implements `OCIBuildService` interface from Phase 1
   - Manages build configuration and state
   - Orchestrates the build process
   - Handles build status tracking

2. **Layer Management** (`layers.go`):
   - Creates and manages image layers
   - Implements layer caching
   - Calculates layer digests
   - Manages layer storage

3. **Module Setup** (`go.mod`, partial `go.sum`):
   - Establishes module structure
   - Defines core dependencies (buildah, storage)
   - Sets up Phase 1 integration imports

## Dependencies
- **External**: None (foundational split)
- **Phase 1 Integration**: 
  - Imports api types from `github.com/idpbuilder/oci-mgmt/integrations/phase1/wave1/integration-workspace/pkg/oci/api`
  - Implements `OCIBuildService` interface
- **Libraries**:
  - `github.com/containers/buildah` for build operations
  - `github.com/containers/storage` for layer storage

## Implementation Instructions for SW Engineer

### 1. Setup Split Working Directory
```bash
# Orchestrator will provide exact paths
cd efforts/phase2/wave1/dockerfile-builder/split-001
git checkout -b idpbuilder-oci-mgmt/phase2/wave1/dockerfile-builder-split-001
```

### 2. File Implementation Order
1. First implement `layers.go` - standalone component
2. Then implement `builder.go` - uses layers
3. Add tests for both components
4. Setup go.mod with required dependencies
5. Generate initial go.sum (first 143 lines)

### 3. Key Implementation Points
- Ensure `Builder` struct properly implements `OCIBuildService` interface
- LayerManager should handle concurrent access (use mutexes)
- Include proper error handling and validation
- Follow Phase 1 API contracts exactly

### 4. Testing Requirements
- Unit tests for all public methods
- Test concurrent layer access
- Test builder initialization and configuration
- Mock external dependencies appropriately

### 5. go.sum Handling
Since go.sum is being split:
```bash
# After implementing and testing, generate full go.sum
go mod tidy

# Extract only first 143 lines for this split
head -143 go.sum > go.sum.split001

# Move the partial file to go.sum
mv go.sum.split001 go.sum

# Document in work-log that go.sum is partial
```

### 6. Size Verification
```bash
# Before committing, verify size
/home/vscode/workspaces/idpbuilder-oci-mgmt/tools/line-counter.sh --base origin/software-factory-2.0
# Must show ≤750 lines
```

## Quality Checklist
- [ ] All files compile without errors
- [ ] Tests pass for builder and layers
- [ ] Size under 750 lines (verified with line-counter.sh)
- [ ] Follows Phase 1 interface contracts
- [ ] Proper error handling throughout
- [ ] Concurrent access handled safely
- [ ] go.sum partial (143 lines) is valid
- [ ] Documentation comments on all exports

## Integration Notes for Split-002
- Split-002 will add dockerfile.go and its tests
- Split-002 will complete the go.sum file (adding remaining lines)
- Builder in Split-001 sets foundation for dockerfile parser integration
- Layer manager will be used by dockerfile processor in Split-002

## Commit Message Template
```
feat(dockerfile-builder): implement core builder and layer management (split-001)

- Implement OCIBuildService interface
- Add layer management with caching
- Setup initial module dependencies
- Add comprehensive unit tests

Part 1 of 2 - Core infrastructure
Size: 750 lines (under 800 limit)
```