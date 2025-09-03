# SPLIT-PLAN-004: Complete Test Suite and Integration

## Split 004 of 4: Remaining Tests and Full Integration
**Planner**: Code Reviewer code-reviewer (same for ALL splits)
**Parent Effort**: go-containerregistry-image-builder
**Target Size**: ~635 lines (well under 800 limit)

<!-- ORCHESTRATOR METADATA PLACEHOLDER - DO NOT REMOVE -->
<!-- The orchestrator will add infrastructure metadata below: -->
<!-- WORKING_DIRECTORY, BRANCH, REMOTE, BASE_BRANCH, etc. -->
<!-- SW Engineers MUST read this metadata to navigate to the correct directory -->
<!-- END PLACEHOLDER -->

## Boundaries (CRITICAL: All splits MUST reference SAME effort!)

- **Previous Split**: Split 003 of phase2/wave1/go-containerregistry-image-builder
  - Path: efforts/phase2/wave1/go-containerregistry-image-builder/split-003/
  - Branch: idpbuilder-oci-go-cr/phase2/wave1/go-containerregistry-image-builder-split-003
  - Summary: Implemented tarball operations and tests up to line 350
  
- **This Split**: Split 004 of phase2/wave1/go-containerregistry-image-builder (FINAL)
  - Path: efforts/phase2/wave1/go-containerregistry-image-builder/split-004/
  - Branch: idpbuilder-oci-go-cr/phase2/wave1/go-containerregistry-image-builder-split-004
  
- **Next Split**: None (this is the final split)

## Files in This Split (EXCLUSIVE - no overlap)

1. **pkg/builder/builder_test.go** (lines 351-674) (~324 lines)
   - Remainder of test file
   - Advanced test scenarios
   - Integration tests
   - Performance tests
   
2. **pkg/builder/testdata/** (~25 lines)
   - Dockerfile samples
   - Test content files
   - Configuration samples
   - Mock data

3. **Integration and documentation** (~286 lines)
   - Integration test scenarios
   - Package documentation
   - README updates
   - Example usage code

**Total Estimated**: ~635 lines

## Functionality Scope

### Core Responsibilities
1. Complete the comprehensive test suite
2. Add integration test scenarios
3. Provide test data and fixtures
4. Ensure full code coverage
5. Add package documentation
6. Verify end-to-end functionality

### Key Components to Implement
```go
- Advanced builder tests (lines 351-674)
- Multi-stage build tests
- Registry interaction tests
- Performance benchmarks
- Test fixtures and data
- Documentation and examples
```

## Dependencies

### From Previous Splits
- **Split 001**: All configuration types
- **Split 002**: Builder and layer implementations
- **Split 003**: Tarball operations, first part of tests

### This Split Completes
- Full test coverage for the entire package
- Integration verification
- Documentation and examples

## Implementation Instructions

### Step 1: Environment Setup
```bash
# Navigate to split directory
cd efforts/phase2/wave1/go-containerregistry-image-builder/split-004/

# Verify branch
git branch --show-current  # Should show split-004 branch

# Merge all previous splits
git merge origin/idpbuilder-oci-go-cr/phase2/wave1/go-containerregistry-image-builder-split-003
```

### Step 2: Complete builder_test.go (Lines 351-674)
```bash
# Continue the test file from line 351
# Edit pkg/builder/builder_test.go
```

Add remaining tests:
```go
// Starting from line 351...
func TestMultiStageBuild(t *testing.T)
func TestRegistryPush(t *testing.T)
func TestRegistryPull(t *testing.T)
func TestBuildCache(t *testing.T)
func TestConcurrentBuilds(t *testing.T)
func TestLargeImageBuild(t *testing.T)
func TestBuildFromDockerfile(t *testing.T)
func TestLayerReuse(t *testing.T)
func TestPlatformSpecificBuild(t *testing.T)
// ... continue to line 674
```

### Step 3: Create Test Data
```bash
# Create testdata directory and files
mkdir -p pkg/builder/testdata
mkdir -p pkg/builder/testdata/content

# Create sample Dockerfile
cat > pkg/builder/testdata/Dockerfile << 'EOF'
FROM alpine:latest
RUN apk add --no-cache curl
COPY content/ /app/
WORKDIR /app
CMD ["./run.sh"]
EOF

# Create test content files
echo "test application content" > pkg/builder/testdata/content/app.txt
cat > pkg/builder/testdata/content/config.yaml << 'EOF'
version: 1.0
settings:
  debug: true
  port: 8080
EOF
```

### Step 4: Add Integration Tests
```bash
# Create integration test file
touch pkg/builder/integration_test.go
```

Implement:
- End-to-end build scenarios
- Real registry interaction tests (with mocks)
- Multi-architecture builds
- Cache efficiency tests
- Performance benchmarks

### Step 5: Add Documentation
```bash
# Update or create package documentation
cat > pkg/builder/doc.go << 'EOF'
// Package builder provides functionality for building OCI container images
// using the go-containerregistry library.
//
// The builder supports:
// - Creating images from scratch
// - Adding layers from tarballs
// - Multi-stage builds
// - Registry push/pull operations
// - Layer caching for efficiency
//
// Example usage:
//
//   cfg := &BuildConfig{
//       Registry: "docker.io",
//       Repository: "myapp",
//       Tag: "latest",
//   }
//   
//   b := NewBuilder(cfg)
//   img, err := b.Build(ctx)
//   if err != nil {
//       log.Fatal(err)
//   }
//
package builder
EOF
```

### Step 6: Add Examples
```bash
# Create example file
touch pkg/builder/example_test.go
```

Add example functions:
```go
func ExampleBuilder_Build()
func ExampleCreateTarball()
func ExampleBuilder_Push()
```

### Step 7: Verify Full Integration
```bash
# Run all tests
go test ./pkg/builder/... -v -cover

# Ensure coverage is >80%
go test ./pkg/builder/... -coverprofile=coverage.out
go tool cover -html=coverage.out

# Run benchmarks
go test ./pkg/builder/... -bench=.

# Check for race conditions
go test ./pkg/builder/... -race
```

### Step 8: Final Size Measurement
```bash
# Measure the complete implementation
PROJECT_ROOT=$(pwd)
while [ "$PROJECT_ROOT" != "/" ]; do 
    if [ -f "$PROJECT_ROOT/orchestrator-state.yaml" ]; then break; fi
    PROJECT_ROOT=$(dirname "$PROJECT_ROOT")
done

# Check this split's addition
$PROJECT_ROOT/tools/line-counter.sh -b idpbuilder-oci-go-cr/phase2/wave1/go-containerregistry-image-builder-split-003 -c $(git branch --show-current)

# Should show ~635 new lines
```

### Step 9: Commit Final Implementation
```bash
# Add all remaining files
git add pkg/builder/
git add -A

# Final commit
git commit -m "feat(split-004): complete test suite and integration

- Complete builder_test.go (lines 351-674)
- Add comprehensive test data and fixtures
- Include integration tests and benchmarks
- Add package documentation and examples
- Achieve >80% test coverage

Part 4 of 4 in go-containerregistry-image-builder effort
FINAL SPLIT - Ready for integration"

# Push
git push origin $(git branch --show-current)
```

## Success Criteria

1. ✅ All tests complete and passing
2. ✅ Test coverage >80% for entire package
3. ✅ Integration tests validate end-to-end
4. ✅ Benchmarks show acceptable performance
5. ✅ Documentation comprehensive
6. ✅ Examples clear and runnable
7. ✅ Total lines under 700 for this split
8. ✅ Full effort ready for merge

## Testing Requirements

### Remaining Tests Must Cover:
- Multi-stage build scenarios
- Registry operations (with mocks)
- Concurrent build operations
- Large image handling
- Cache effectiveness
- Platform-specific builds
- Error recovery scenarios
- Performance characteristics

### Quality Metrics:
- Coverage: >80% overall
- No race conditions
- All benchmarks pass
- Examples compile and run
- Documentation builds correctly

### Final Validation:
```bash
# Run full test suite
go test ./... -v -cover -race

# Verify examples
go test ./pkg/builder/... -run Example

# Check lint
golangci-lint run ./...
```

## Review Checklist

Before marking complete:
- [ ] All 4 splits integrate cleanly
- [ ] Total implementation is 2585 lines (±50)
- [ ] Tests achieve >80% coverage
- [ ] No functionality lost from original
- [ ] Documentation complete
- [ ] Examples working
- [ ] Performance acceptable
- [ ] Ready for production use

## Notes for SW Engineer

- This is the final split - ensure everything integrates
- Focus on comprehensive test coverage
- Document any known limitations
- Ensure examples demonstrate key features
- Run full regression testing
- Verify no functionality was lost in splitting
- Prepare summary for review

## Final Integration Notes

### After This Split Completes:
1. All 4 splits should be reviewed
2. Branches should be merged sequentially:
   - split-001 → main effort branch
   - split-002 → main effort branch
   - split-003 → main effort branch
   - split-004 → main effort branch
3. Final effort branch ready for phase integration

### Verification Steps:
1. Full package compiles
2. All tests pass
3. Coverage meets requirements
4. Line count within limits
5. Documentation complete
6. No regressions from original

### Success Indicators:
- Clean merge of all splits
- No conflicts between splits
- Functionality preserved
- Performance maintained
- Quality metrics met