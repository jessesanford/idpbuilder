# SPLIT-PLAN-001: Core Types and Configuration

## Split 001 of 4: Foundation Types and Build Configuration
**Planner**: Code Reviewer code-reviewer (same for ALL splits)
**Parent Effort**: go-containerregistry-image-builder
**Target Size**: ~650 lines (well under 800 limit)

<!-- ORCHESTRATOR METADATA PLACEHOLDER - DO NOT REMOVE -->
<!-- The orchestrator will add infrastructure metadata below: -->
<!-- WORKING_DIRECTORY, BRANCH, REMOTE, BASE_BRANCH, etc. -->
<!-- SW Engineers MUST read this metadata to navigate to the correct directory -->
<!-- END PLACEHOLDER -->

## Boundaries (CRITICAL: All splits MUST reference SAME effort!)

- **Previous Split**: None (first split of THIS effort)
  - Path: N/A (this is Split 001)
  - Branch: N/A
  
- **This Split**: Split 001 of phase2/wave1/go-containerregistry-image-builder
  - Path: efforts/phase2/wave1/go-containerregistry-image-builder/split-001/
  - Branch: idpbuilder-oci-go-cr/phase2/wave1/go-containerregistry-image-builder-split-001
  
- **Next Split**: Split 002 of phase2/wave1/go-containerregistry-image-builder
  - Path: efforts/phase2/wave1/go-containerregistry-image-builder/split-002/
  - Branch: idpbuilder-oci-go-cr/phase2/wave1/go-containerregistry-image-builder-split-002

## Files in This Split (EXCLUSIVE - no overlap)

1. **pkg/builder/config.go** (~319 lines)
   - Complete file implementation
   - Build configuration structures
   - Registry configuration
   - Platform settings
   
2. **pkg/builder/options.go** (~133 lines)
   - Complete file implementation
   - Functional options pattern
   - Configuration helpers
   - Validation logic

3. **Basic package structure** (~50 lines)
   - Package documentation
   - Import statements
   - Common constants
   - Error definitions

4. **Initial unit tests** (~150 lines)
   - Configuration validation tests
   - Options tests
   - Basic error cases

**Total Estimated**: ~652 lines

## Functionality Scope

### Core Responsibilities
1. Define all configuration structures for image building
2. Implement option patterns for flexible configuration
3. Registry authentication and connection settings
4. Platform-specific build configurations
5. Validation logic for all settings

### Key Types to Implement
```go
- BuildConfig struct
- ImageOptions struct
- RegistryConfig struct
- PlatformConfig struct
- BuildOptions interface
- Validation methods
```

## Dependencies
- **External**: None (foundational split)
- **Internal**: None (first split)
- **Go modules**: Standard library only for this split

## Implementation Instructions

### Step 1: Environment Setup
```bash
# Navigate to split directory (created by orchestrator)
cd efforts/phase2/wave1/go-containerregistry-image-builder/split-001/

# Verify clean workspace
git status
git branch --show-current  # Should show split-001 branch
```

### Step 2: Create Package Structure
```bash
# Create the builder package directory
mkdir -p pkg/builder

# Create initial package files
touch pkg/builder/config.go
touch pkg/builder/options.go
touch pkg/builder/doc.go
```

### Step 3: Implement config.go
- Define BuildConfig struct with all necessary fields
- Add registry configuration structures
- Implement platform-specific settings
- Add validation methods for configurations
- Include proper error handling

### Step 4: Implement options.go
- Create functional options pattern implementation
- Define option functions for each configurable field
- Add helper methods for common configurations
- Implement validation within options

### Step 5: Create Initial Tests
```bash
# Create test file
touch pkg/builder/config_test.go
touch pkg/builder/options_test.go
```
- Write unit tests for configuration validation
- Test all option functions
- Cover error cases and edge conditions
- Ensure 80%+ coverage for this split

### Step 6: Verify Compilation
```bash
# Ensure the package compiles
go build ./pkg/builder/...

# Run tests
go test ./pkg/builder/... -v

# Check test coverage
go test ./pkg/builder/... -cover
```

### Step 7: Measure Size
```bash
# Use the line counter tool
$PROJECT_ROOT/tools/line-counter.sh -b idpbuilder-oci-go-cr/phase1-integration-20250902-194557 -c $(git branch --show-current)

# Verify under 700 lines
```

### Step 8: Commit Implementation
```bash
# Add all files
git add pkg/builder/

# Commit with descriptive message
git commit -m "feat(split-001): implement core types and configuration

- Add BuildConfig and related structures
- Implement options pattern for configuration
- Add validation logic for all settings
- Include initial unit tests with 80%+ coverage

Part 1 of 4 in go-containerregistry-image-builder effort"

# Push to remote
git push origin $(git branch --show-current)
```

## Success Criteria

1. ✅ All configuration structures properly defined
2. ✅ Options pattern correctly implemented
3. ✅ Validation logic comprehensive and tested
4. ✅ Package compiles without errors
5. ✅ Tests pass with >80% coverage
6. ✅ Total lines under 700 (measured with line-counter.sh)
7. ✅ No dependencies on later splits
8. ✅ Clean git history with single commit

## Testing Requirements

### Unit Tests Must Cover:
- Configuration struct initialization
- All option functions
- Validation logic for each field
- Error cases and boundaries
- Default value handling
- Registry authentication scenarios

### Test Execution:
```bash
go test ./pkg/builder/... -v -cover
```

## Review Checklist

Before marking complete:
- [ ] Code follows Go best practices
- [ ] All exported types/functions have godoc comments
- [ ] Error messages are clear and actionable
- [ ] No TODO or FIXME comments remaining
- [ ] Line count verified under 700
- [ ] Tests achieve >80% coverage
- [ ] No compilation warnings
- [ ] Git commit follows conventional format

## Notes for SW Engineer

- This is the foundation split - take time to get the structure right
- Other splits will depend on these types, so ensure clean interfaces
- Focus on making configuration flexible but validated
- Keep external dependencies minimal
- Document all exported types thoroughly
- Consider future extensibility in the design