# CLI Commands Integration Implementation Plan

## 🚨 CRITICAL EFFORT METADATA (FROM WAVE PLAN)
**Branch**: `idpbuilder-oci-mvp/phase2/wave2/cli-commands`
**Can Parallelize**: No
**Parallel With**: None
**Size Estimate**: 300 lines
**Dependencies**: Phase 2 Wave 1 (buildah-build-wrapper, gitea-registry-client)

## Overview
- **Effort**: Integrate buildah build and push functionality into existing idpbuilder CLI
- **Phase**: 2, Wave: 2
- **Estimated Size**: 300 lines (modifications to existing code)
- **Implementation Time**: 2-3 hours

## 🚨 CRITICAL REQUIREMENTS
This effort MUST:
- **MODIFY** existing idpbuilder CLI code (NOT create new standalone CLI)
- **ADD** commands to the existing command structure
- **INTEGRATE** with existing patterns and utilities
- **REUSE** existing configuration and flag handling

## File Structure (MODIFICATIONS ONLY)
- `pkg/cmd/build/` - NEW directory for build command
  - `root.go` - Build command definition (NEW ~80 lines)
- `pkg/cmd/push/` - NEW directory for push command
  - `root.go` - Push command definition (NEW ~80 lines)
- `pkg/cmd/root.go` - MODIFY to register new commands (~10 lines)
- `pkg/build/` - Integration wrapper for buildah functionality
  - `integration.go` - Wrapper to connect to Phase 2 Wave 1 code (NEW ~60 lines)
- `pkg/registry/` - Integration wrapper for registry functionality
  - `integration.go` - Wrapper to connect to Phase 2 Wave 1 code (NEW ~60 lines)
- Tests for new commands (~20 lines)

## Implementation Steps

### Step 1: Create Build Command Structure
**Task**: Add build command to existing CLI
**Files to create**:
- `pkg/cmd/build/root.go`

**Implementation**:
```go
// pkg/cmd/build/root.go
package build

import (
    "github.com/spf13/cobra"
    "github.com/cnoe-io/idpbuilder/pkg/cmd/helpers"
    // Import Phase 2 Wave 1 buildah wrapper
)

var BuildCmd = &cobra.Command{
    Use:   "build [dockerfile]",
    Short: "Build container images using Buildah",
    Long:  "Build OCI container images from Dockerfiles with certificate support",
    RunE:  runBuild,
}

var (
    dockerfilePath string
    contextDir     string
    tag           string
    insecure      bool
)

func init() {
    BuildCmd.Flags().StringVarP(&dockerfilePath, "file", "f", "Dockerfile", "Path to Dockerfile")
    BuildCmd.Flags().StringVar(&contextDir, "context", ".", "Build context directory")
    BuildCmd.Flags().StringVarP(&tag, "tag", "t", "", "Image tag")
    BuildCmd.Flags().BoolVar(&insecure, "insecure", false, "Skip certificate validation")
}
```

### Step 2: Create Push Command Structure
**Task**: Add push command to existing CLI
**Files to create**:
- `pkg/cmd/push/root.go`

**Implementation**:
```go
// pkg/cmd/push/root.go
package push

import (
    "github.com/spf13/cobra"
    // Import Phase 2 Wave 1 registry client
)

var PushCmd = &cobra.Command{
    Use:   "push [image] [registry]",
    Short: "Push container images to registry",
    Long:  "Push OCI container images to Gitea or other registries with certificate support",
    RunE:  runPush,
}

var (
    insecure bool
    username string
    password string
)

func init() {
    PushCmd.Flags().BoolVar(&insecure, "insecure", false, "Skip TLS certificate validation")
    PushCmd.Flags().StringVar(&username, "username", "", "Registry username")
    PushCmd.Flags().StringVar(&password, "password", "", "Registry password")
}
```

### Step 3: Integrate Commands into Root
**Task**: Modify existing root.go to register new commands
**Files to modify**:
- `pkg/cmd/root.go`

**Changes**:
```go
// Add imports
import (
    "github.com/cnoe-io/idpbuilder/pkg/cmd/build"
    "github.com/cnoe-io/idpbuilder/pkg/cmd/push"
    // ... existing imports
)

// In init() function, add:
func init() {
    // ... existing flags
    rootCmd.AddCommand(build.BuildCmd)
    rootCmd.AddCommand(push.PushCmd)
    // ... existing commands
}
```

### Step 4: Create Integration Wrappers
**Task**: Create thin wrappers to connect to Phase 2 Wave 1 functionality
**Files to create**:
- `pkg/build/integration.go`
- `pkg/registry/integration.go`

**Purpose**: These files will import and use the buildah-build-wrapper and gitea-registry-client from Phase 2 Wave 1, providing a clean interface for the CLI commands.

### Step 5: Add Tests
**Task**: Create tests for new commands
**Files to create**:
- `pkg/cmd/build/root_test.go` (basic flag parsing tests)
- `pkg/cmd/push/root_test.go` (basic flag parsing tests)

## Size Management
- **Estimated Lines**: 300
- **Measurement Tool**: ${PROJECT_ROOT}/tools/line-counter.sh
- **Check Frequency**: After each step
- **Split Threshold**: 700 lines (warning), 800 lines (stop)

## Breakdown by Component:
- Build command: ~80 lines
- Push command: ~80 lines
- Root.go modifications: ~10 lines
- Integration wrappers: ~120 lines total
- Tests: ~20 lines
- **Total**: ~310 lines (within 300 line estimate with margin)

## Test Requirements
- **Unit Tests**: Command flag parsing
- **Integration Tests**: Mock build/push operations
- **Manual Tests**: 
  - Build with local Dockerfile
  - Push to local Gitea registry
  - Test --insecure flag
  - Test certificate validation

## Pattern Compliance
- **idpbuilder Patterns**: 
  - Use cobra for command structure
  - Follow existing command directory structure
  - Use existing helpers package
  - Match existing error handling patterns
- **Integration Requirements**:
  - Import Phase 2 Wave 1 packages correctly
  - Handle certificate trust manager integration
  - Maintain backward compatibility

## Dependencies from Phase 2 Wave 1
The implementation will import and use:
- `/efforts/phase2/wave1/buildah-build-wrapper/pkg/build/` - Buildah functionality
- Certificate trust management interfaces
- Registry client functionality

## Success Criteria Checklist
- [ ] Build command added to existing CLI
- [ ] Push command added to existing CLI
- [ ] Commands follow existing idpbuilder patterns
- [ ] --insecure flag works for both commands
- [ ] Certificate handling properly integrated
- [ ] All existing tests still pass
- [ ] New tests for new commands pass
- [ ] Total modifications under 400 lines
- [ ] Manual testing successful with local Gitea

## Split Contingency Plan
If implementation exceeds 300 lines:
1. **Split Point 1** (if >400 lines): Separate build and push into two efforts
2. **Split Point 2** (if >500 lines): Extract integration wrappers as separate effort
3. **Split Point 3** (if >600 lines): Move tests to separate effort

## Next Steps for SW Engineer
1. Verify workspace is in correct directory
2. Study existing command patterns in detail
3. Implement build command first
4. Test build command thoroughly
5. Implement push command
6. Test push command thoroughly
7. Create integration wrappers
8. Run all tests
9. Measure final size with line-counter.sh
10. Request code review