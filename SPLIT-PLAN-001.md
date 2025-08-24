# SPLIT-PLAN-001.md
## Split 001 of 2: OCI Package Implementation
**Planner**: Code Reviewer @agent-code-reviewer-1756072217 (same for ALL splits)
**Parent Effort**: oci-types
**Branch Name**: phase1/wave1/oci-types-split-001
**Size**: 622 lines (well under 800 limit)

## Boundaries
- **Previous Split**: None (first split)
- **This Split Coverage**: Complete OCI package (pkg/oci/*)
- **Next Split**: Stack package and documentation

## Files in This Split (EXCLUSIVE - no overlap with other splits)

### Production Code (301 lines)
1. **pkg/oci/types.go** (121 lines)
   - Core OCI type definitions
   - OCIImage struct with Name, Tag, Digest, Registry, Repository
   - OCIRepository interface for List, Get, Push, Pull operations
   - OCIReference type for image references
   - OCILayer struct for layer information
   - Validation methods for image references

2. **pkg/oci/manifest.go** (124 lines)
   - OCIManifest struct following OCI specification
   - OCIDescriptor type with MediaType, Digest, Size, URLs, Annotations
   - OCIManifestList for multi-platform support
   - OCIConfig for image configuration
   - Helper methods for manifest validation

3. **pkg/oci/constants.go** (56 lines)
   - Media type constants (application/vnd.oci.*)
   - Annotation keys from OCI specification
   - Default registry values
   - Platform defaults

### Test Code (321 lines)
4. **pkg/oci/types_test.go** (130 lines)
   - Unit tests for type creation
   - Validation tests for OCI references
   - Edge cases and error conditions

5. **pkg/oci/manifest_test.go** (191 lines)
   - Manifest marshaling/unmarshaling tests
   - Descriptor validation tests
   - Multi-platform manifest tests

## Directory Structure
```
efforts/phase1/wave1/oci-types-split-001/
├── pkg/
│   └── oci/
│       ├── types.go         (121 lines)
│       ├── manifest.go      (124 lines)
│       ├── constants.go     (56 lines)
│       ├── types_test.go    (130 lines)
│       └── manifest_test.go (191 lines)
├── go.mod
├── go.sum
└── Makefile
```

## Functionality Scope
This split implements the complete OCI (Open Container Initiative) type system:
- Image representation and metadata
- Manifest structures for single and multi-platform images
- Layer and descriptor types
- Standard OCI constants and media types
- Full validation logic for OCI references
- Comprehensive test coverage (>80%)

## Dependencies
- **External**: None (standard library only)
- **Internal**: None (foundational split)
- **Required by**: Split 002 (stack package references OCIReference)

## Implementation Instructions

### 1. Create New Branch
```bash
cd /home/vscode/workspaces/idpbuilder-oci-mgmt
git checkout main
git pull origin main
git checkout -b phase1/wave1/oci-types-split-001
```

### 2. Create Sparse Checkout (ONLY these files)
```bash
# Create effort directory structure
mkdir -p efforts/phase1/wave1/oci-types-split-001/pkg/oci

# Copy ONLY the OCI package files
cp efforts/phase1/wave1/oci-types/pkg/oci/types.go \
   efforts/phase1/wave1/oci-types-split-001/pkg/oci/

cp efforts/phase1/wave1/oci-types/pkg/oci/manifest.go \
   efforts/phase1/wave1/oci-types-split-001/pkg/oci/

cp efforts/phase1/wave1/oci-types/pkg/oci/constants.go \
   efforts/phase1/wave1/oci-types-split-001/pkg/oci/

cp efforts/phase1/wave1/oci-types/pkg/oci/types_test.go \
   efforts/phase1/wave1/oci-types-split-001/pkg/oci/

cp efforts/phase1/wave1/oci-types/pkg/oci/manifest_test.go \
   efforts/phase1/wave1/oci-types-split-001/pkg/oci/

# Copy supporting files
cp efforts/phase1/wave1/oci-types/go.mod \
   efforts/phase1/wave1/oci-types-split-001/

cp efforts/phase1/wave1/oci-types/go.sum \
   efforts/phase1/wave1/oci-types-split-001/

cp efforts/phase1/wave1/oci-types/Makefile \
   efforts/phase1/wave1/oci-types-split-001/
```

### 3. Verify Compilation
```bash
cd efforts/phase1/wave1/oci-types-split-001
go build ./pkg/oci/...
go test ./pkg/oci/... -v
```

### 4. Measure Size
```bash
/home/vscode/workspaces/idpbuilder-oci-mgmt/tools/line-counter.sh
# Must show ≤622 lines (excluding work-log.md and plan files)
```

### 5. Run Tests
```bash
go test ./pkg/oci/... -cover
# Must achieve >80% coverage
```

## Quality Requirements
- ✅ All code must compile without errors
- ✅ All tests must pass
- ✅ Test coverage must exceed 80%
- ✅ Total size must not exceed 622 lines
- ✅ All exported types must have godoc comments
- ✅ No implementation logic (types and constants only)

## Commit Message
```
feat(oci): implement OCI types package (split 001/002)

- Add core OCI type definitions (OCIImage, OCIRepository, OCIReference)
- Add manifest types (OCIManifest, OCIDescriptor, OCIManifestList)
- Add OCI specification constants and media types
- Add comprehensive unit tests (>80% coverage)
- Size: 622 lines (under 800 limit)

Part of oci-types effort split due to size constraints
```

## Review Checklist
- [ ] Only OCI package files included (no stack files)
- [ ] All 5 files present and correct
- [ ] Compilation successful
- [ ] Tests passing with >80% coverage
- [ ] Size verified with line-counter.sh
- [ ] No dependencies on Split 002
- [ ] Ready for independent merge

## Next Steps
After this split is merged:
1. Orchestrator creates Split 002 branch
2. SW Engineer implements stack package
3. Split 002 can reference types from this split
4. Both splits integrate to complete original effort