# SPLIT-PLAN-001: Core Builder Interface and Configuration

## Split 001 of 4: Foundation Components
**Planner**: Code Reviewer Agent (same for ALL splits)
**Parent Effort**: go-containerregistry-image-builder
**Created**: 2025-09-03 04:56:00

<!-- ⚠️ ORCHESTRATOR METADATA PLACEHOLDER - DO NOT REMOVE ⚠️ -->
<!-- The orchestrator will add infrastructure metadata below: -->
<!-- WORKING_DIRECTORY, BRANCH, REMOTE, BASE_BRANCH, etc. -->
<!-- SW Engineers MUST read this metadata to navigate to the correct directory -->
<!-- END PLACEHOLDER -->

## Boundaries (⚠️⚠️⚠️ CRITICAL: All splits MUST reference SAME effort!)

- **Previous Split**: None (first split of THIS effort)
  - Path: N/A (this is Split 001)
  - Branch: N/A
  
- **This Split**: Split 001 of phase2/wave1/go-containerregistry-image-builder
  - Path: efforts/phase2/wave1/go-containerregistry-image-builder/split-001/
  - Branch: idpbuilder-oci-go-cr/phase2/wave1/go-containerregistry-image-builder-split-001
  
- **Next Split**: Split 002 of phase2/wave1/go-containerregistry-image-builder
  - Path: efforts/phase2/wave1/go-containerregistry-image-builder/split-002/
  - Branch: idpbuilder-oci-go-cr/phase2/wave1/go-containerregistry-image-builder-split-002

## Files in This Split (EXCLUSIVE - no overlap with other splits)

### Core Implementation Files (~650 lines total)
1. **pkg/builder/builder.go** (163 lines)
   - Builder interface definition
   - SimpleBuilder struct
   - NewBuilder constructor
   - Build method implementation (partial - interface only)

2. **pkg/builder/options.go** (132 lines)
   - BuildOptions struct
   - Platform configuration
   - Feature flags
   - Validation methods

3. **pkg/builder/config.go** (318 lines)
   - ConfigFactory implementation
   - OCI config generation
   - Platform-specific configurations
   - Label and environment handling

4. **pkg/builder/doc.go** (~37 lines - new file)
   - Package documentation
   - Interface usage examples
   - Feature flag documentation

## Functionality in This Split

### Primary Components:
- **Builder Interface**: Core contract for image building
- **Configuration Management**: OCI image config generation
- **Build Options**: Comprehensive option handling
- **Platform Support**: Multi-architecture configuration

### Feature Flags (R307 Compliance):
```go
// Feature flags for incomplete functionality
const (
    FeatureTarballExport = "tarball_export" // Disabled in Split 001
    FeatureLayerCaching  = "layer_caching"   // Disabled in Split 001
    FeatureMultiLayer    = "multi_layer"     // Disabled in Split 001
)
```

## Dependencies

### External:
- github.com/google/go-containerregistry v0.19.0
- Standard library (context, fmt, os, path/filepath)

### From Phase 1 (imported, not reimplemented):
- Certificate infrastructure (will import in Split 003)

### Internal (this split):
- None - this is the foundation

## Implementation Instructions

### Step 1: Infrastructure Setup
```bash
# Create split directory structure
mkdir -p efforts/phase2/wave1/go-containerregistry-image-builder/split-001/pkg/builder
cd efforts/phase2/wave1/go-containerregistry-image-builder/split-001

# Initialize go module
go mod init github.com/cnoe-io/idpbuilder/pkg/builder
go get github.com/google/go-containerregistry@v0.19.0
```

### Step 2: Implement Core Files
1. Create `pkg/builder/doc.go` with package documentation
2. Implement `pkg/builder/options.go` with BuildOptions struct
3. Implement `pkg/builder/config.go` with ConfigFactory
4. Implement `pkg/builder/builder.go` with interface and basic structure

### Step 3: Add Stubs for Incomplete Features
```go
// In builder.go - stub implementation
func (b *SimpleBuilder) Build(ctx context.Context, contextDir string, opts BuildOptions) (v1.Image, error) {
    if !b.featureFlags[FeatureTarballExport] {
        return nil, fmt.Errorf("tarball export not enabled in this split")
    }
    // Stub for now - will be completed in Split 002
    return nil, fmt.Errorf("not implemented in Split 001")
}
```

### Step 4: Unit Tests
Create basic unit tests for:
- Options validation
- Config generation
- Platform detection
- Feature flag checking

### Step 5: Verification
```bash
# Ensure compilation
go build ./...

# Run tests
go test ./... -v

# Measure size (MUST be <700 lines)
PROJECT_ROOT=/home/vscode/workspaces/idpbuilder-oci-go-cr
$PROJECT_ROOT/tools/line-counter.sh
```

## Expected Outcome

After Split 001 completion:
- ✅ Core interfaces defined
- ✅ Configuration management working
- ✅ Build options validated
- ✅ Compiles independently
- ✅ Basic tests passing
- ✅ Size <700 lines

## Integration Notes

This split provides the foundation that Split 002-004 will build upon:
- Split 002 will implement layer and tarball functionality
- Split 003 will add build utilities
- Split 004 will complete test coverage

## Risk Mitigation

- All methods that depend on future splits return clear error messages
- Feature flags prevent accidental use of incomplete functionality
- Interfaces allow for mock implementations in tests
## 🚨 SPLIT INFRASTRUCTURE METADATA (Added by Orchestrator)
**WORKING_DIRECTORY**: /home/vscode/workspaces/idpbuilder-oci-go-cr/efforts/phase2/wave1/go-containerregistry-image-builder-SPLIT-001
**BRANCH**: idpbuilder-oci-go-cr/phase2/wave1/go-containerregistry-image-builder--split-001
**REMOTE**: origin/idpbuilder-oci-go-cr/phase2/wave1/go-containerregistry-image-builder--split-001
**BASE_BRANCH**: idpbuilder-oci-go-cr/phase1-integration-20250902-194557
**SPLIT_NUMBER**: 001
**CREATED_AT**: 2025-09-03 05:05:00

### SW Engineer Instructions
1. READ this metadata FIRST
2. cd to WORKING_DIRECTORY above
3. Verify branch matches BRANCH above
4. ONLY THEN proceed with implementation
