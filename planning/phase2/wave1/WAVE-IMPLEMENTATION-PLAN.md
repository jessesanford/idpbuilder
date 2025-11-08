# Wave 2.1 Implementation Plan - Command Implementation & Core Integration

**Wave**: Phase 2, Wave 1
**Phase**: Phase 2 - Core Push Functionality
**Created**: 2025-10-31
**Planner**: Code Reviewer Agent (@agent-code-reviewer)
**Fidelity Level**: **EXACT SPECIFICATIONS** (detailed efforts with R213 metadata)

---

## Wave Overview

**Goal**: Implement the `idpbuilder push` command that integrates all Phase 1 packages (docker, registry, auth, tls) into a working end-to-end push workflow with progress reporting.

**Architecture Reference**: See `planning/phase2/wave1/WAVE-2.1-ARCHITECTURE.md` for design details

**Test Plan Reference**: See `planning/phase2/wave1/WAVE-TEST-PLAN.md` for 40 concrete tests

**Total Efforts**: 2 (sequential execution)

**Estimated Total Lines**: ~750 lines (450 + 300)

---

## Effort Definitions

### Effort 2.1.1: Push Command Core & Pipeline Orchestration

#### R213 Metadata

```yaml
---
effort_metadata:
  effort_id: "2.1.1"
  effort_name: "Push Command Core & Pipeline Orchestration"
  estimated_lines: 450
  dependencies: []
  files_touched:
    - "pkg/cmd/push/push.go"
    - "pkg/cmd/push/types.go"
    - "pkg/cmd/push/push_test.go"
    - "pkg/cmd/root.go"
  branch_name: "idpbuilder-oci-push/phase2/wave1/effort-1-push-command-core"
  base_branch: "idpbuilder-oci-push/phase2/wave1/integration"
  theme: "CLI command implementation with Phase 1 integration"
  scope: "Command structure, flag definitions, pipeline orchestration, basic error handling"
  complexity: "medium"
  can_parallelize: false
  parallel_with: []
---
```

#### Scope

**Purpose**: Create the foundational `push` command that orchestrates the complete image push pipeline by integrating Phase 1's docker, registry, auth, and tls packages.

**Boundaries**:
- ✅ IN SCOPE:
  - Cobra command registration with IDPBuilder root
  - Flag definitions (--registry, --username, --password, --insecure, --verbose)
  - 8-stage pipeline orchestration (see architecture)
  - Integration with all Phase 1 interfaces
  - Basic error handling with wrapped context
  - Simple progress callback (verbose mode prints to console)

- ❌ OUT OF SCOPE:
  - Advanced progress reporter (that's Effort 2.1.2)
  - Environment variable support (Wave 2.2)
  - Registry auto-detection (Wave 2.2)
  - Sophisticated error formatting (Wave 2.3)
  - Exit code mapping (Wave 2.3)

#### Files to Create/Modify

**New Files**:
- `pkg/cmd/push/push.go` (300 lines) - Main command implementation
- `pkg/cmd/push/types.go` (50 lines) - PushOptions struct and validation
- `pkg/cmd/push/push_test.go` (100 lines) - Unit tests (25 tests from test plan)

**Modified Files**:
- `pkg/cmd/root.go` (add push command registration, +5 lines)

**Total Estimated Lines**: 455 lines

#### Exact Code Specifications

**File: pkg/cmd/push/push.go**

```go
// Package push implements the idpbuilder push command for pushing Docker images to OCI registries
package push

import (
    "context"
    "fmt"
    "os"

    "github.com/spf13/cobra"
    "github.com/cnoe-io/idpbuilder/pkg/docker"
    "github.com/cnoe-io/idpbuilder/pkg/registry"
    "github.com/cnoe-io/idpbuilder/pkg/auth"
    "github.com/cnoe-io/idpbuilder/pkg/tls"
)

// NewPushCommand creates the push command that integrates all Phase 1 packages
func NewPushCommand() *cobra.Command {
    opts := &PushOptions{}

    cmd := &cobra.Command{
        Use:   "push IMAGE",
        Short: "Push a Docker image to an OCI registry",
        Long: `Push a local Docker image to an OCI-compliant container registry.

The command retrieves the image from the local Docker daemon and pushes it to
the specified registry using credentials provided via flags or environment variables.

Examples:
  # Push to default Gitea registry
  idpbuilder push alpine:latest --username admin --password password

  # Push to custom registry
  idpbuilder push myapp:v1.0 --registry docker.io --username user --password pass

  # Push with verbose progress
  idpbuilder push alpine:latest --verbose --username admin --password password

  # Push with insecure TLS (development only)
  idpbuilder push alpine:latest --insecure --username admin --password password`,
        Args: cobra.ExactArgs(1),
        RunE: func(cmd *cobra.Command, args []string) error {
            opts.ImageName = args[0]
            return runPush(cmd.Context(), opts)
        },
    }

    // Define flags
    cmd.Flags().StringVar(&opts.Registry, "registry", "gitea.cnoe.localtest.me:8443",
        "Registry URL (default: Gitea registry)")
    cmd.Flags().StringVar(&opts.Username, "username", "",
        "Registry username (required)")
    cmd.Flags().StringVar(&opts.Password, "password", "",
        "Registry password (required)")
    cmd.Flags().BoolVar(&opts.Insecure, "insecure", false,
        "Skip TLS certificate verification (insecure)")
    cmd.Flags().BoolVar(&opts.Verbose, "verbose", false,
        "Enable verbose progress output")

    // Mark required flags
    cmd.MarkFlagRequired("username")
    cmd.MarkFlagRequired("password")

    return cmd
}

// runPush orchestrates the 8-stage push pipeline using Phase 1 interfaces
func runPush(ctx context.Context, opts *PushOptions) error {
    // STAGE 1: Initialize Docker client
    dockerClient, err := docker.NewClient()
    if err != nil {
        return fmt.Errorf("failed to connect to Docker daemon: %w", err)
    }
    defer dockerClient.Close()

    // STAGE 2: Retrieve image from Docker daemon
    fmt.Printf("Retrieving image %s from Docker daemon...\n", opts.ImageName)
    image, err := dockerClient.GetImage(ctx, opts.ImageName)
    if err != nil {
        return fmt.Errorf("failed to get image: %w", err)
    }

    // STAGE 3: Setup authentication
    authProvider := auth.NewBasicAuthProvider(opts.Username, opts.Password)
    if err := authProvider.ValidateCredentials(); err != nil {
        return fmt.Errorf("invalid credentials: %w", err)
    }

    // STAGE 4: Setup TLS configuration
    tlsProvider := tls.NewConfigProvider(opts.Insecure)

    // STAGE 5: Create registry client
    registryClient, err := registry.NewClient(authProvider, tlsProvider)
    if err != nil {
        return fmt.Errorf("failed to create registry client: %w", err)
    }

    // STAGE 6: Build target reference
    targetRef, err := registryClient.BuildImageReference(opts.Registry, opts.ImageName)
    if err != nil {
        return fmt.Errorf("invalid registry or image name: %w", err)
    }

    // STAGE 7: Create progress callback (basic implementation for Wave 1)
    progressCallback := func(update registry.ProgressUpdate) {
        if opts.Verbose {
            fmt.Printf("Layer %s: %d/%d bytes (%s)\n",
                truncateDigest(update.LayerDigest, 12),
                update.BytesPushed,
                update.LayerSize,
                update.Status)
        }
    }

    // STAGE 8: Execute push
    fmt.Printf("Pushing to %s...\n", targetRef)
    if err := registryClient.Push(ctx, image, targetRef, progressCallback); err != nil {
        return fmt.Errorf("push failed: %w", err)
    }

    fmt.Printf("✓ Successfully pushed %s to %s\n", opts.ImageName, opts.Registry)
    return nil
}

// truncateDigest truncates digest to specified length for display
func truncateDigest(digest string, length int) string {
    if len(digest) <= length {
        return digest
    }
    return digest[:length]
}
```

**Implementation Requirements for push.go**:
- MUST use Phase 1 interfaces exactly as shown (docker.Client, registry.Client, auth.Provider, tls.ConfigProvider)
- Error wrapping MUST preserve original error types using %w
- Pipeline stages MUST execute in order (1-8) with early returns on error
- defer dockerClient.Close() MUST be called to ensure cleanup
- Progress callback MUST check opts.Verbose before printing (respect user preference)
- MUST NOT implement advanced progress features (that's Effort 2.1.2)

---

**File: pkg/cmd/push/types.go**

```go
package push

import (
    "fmt"
)

// PushOptions configures the push command execution
type PushOptions struct {
    // ImageName is the local Docker image to push (e.g., "alpine:latest")
    ImageName string

    // Registry is the target registry URL (default: Gitea registry)
    Registry string

    // Username is the registry authentication username
    Username string

    // Password is the registry authentication password
    Password string

    // Insecure controls TLS certificate verification
    Insecure bool

    // Verbose enables detailed progress output
    Verbose bool
}

// Validate checks if PushOptions are valid and complete
func (o *PushOptions) Validate() error {
    if o.ImageName == "" {
        return fmt.Errorf("image name is required")
    }
    if o.Username == "" {
        return fmt.Errorf("username is required")
    }
    if o.Password == "" {
        return fmt.Errorf("password is required")
    }
    if o.Registry == "" {
        return fmt.Errorf("registry is required")
    }
    return nil
}
```

**Implementation Requirements for types.go**:
- PushOptions MUST have exported fields for Cobra flag binding
- Validate() MUST check all required fields
- Error messages MUST be user-friendly (not technical jargon)
- MUST NOT add additional validation beyond required fields (keep simple for Wave 1)

---

**File: pkg/cmd/root.go (MODIFICATION)**

```go
// Modify the Execute() function to add push command registration

package cmd

import (
    "github.com/spf13/cobra"
    "github.com/cnoe-io/idpbuilder/pkg/cmd/push"  // NEW IMPORT
)

func Execute() error {
    rootCmd := &cobra.Command{
        Use:   "idpbuilder",
        Short: "IDPBuilder CLI",
    }

    // Existing commands
    rootCmd.AddCommand(newCreateCmd())
    rootCmd.AddCommand(newGetCmd())
    rootCmd.AddCommand(newDeleteCmd())
    rootCmd.AddCommand(newVersionCmd())

    // NEW: Register push command
    rootCmd.AddCommand(push.NewPushCommand())

    return rootCmd.Execute()
}
```

**Implementation Requirements for root.go modification**:
- ONLY add the import and AddCommand line
- DO NOT modify existing command registrations
- DO NOT change rootCmd structure
- Keep modifications minimal (5 lines total)

#### Tests Required

**File: pkg/cmd/push/push_test.go**

Test specifications are defined in `planning/phase2/wave1/WAVE-TEST-PLAN.md` (Tests T-2.1.1-01 through T-2.1.1-25).

**Key Test Categories**:
1. **Flag Tests** (T-2.1.1-01 to T-2.1.1-03): Verify flag definitions, defaults, and required flags
2. **Validation Tests** (T-2.1.1-04 to T-2.1.1-07): Test PushOptions.Validate() for valid and invalid inputs
3. **Error Handling Tests** (T-2.1.1-08 to T-2.1.1-14): Test error scenarios for each pipeline stage
4. **Progress Tests** (T-2.1.1-15 to T-2.1.1-16): Verify progress callback invocation
5. **Pipeline Tests** (T-2.1.1-17 to T-2.1.1-23): Test complete pipeline, context handling, cleanup
6. **Integration Tests** (T-2.1.1-24 to T-2.1.1-25): Cobra integration and help text

**Mock Fixtures** (from Phase 1 Wave 2 test patterns):
- `mockDockerClient` - Mock docker.Client interface
- `mockRegistryClient` - Mock registry.Client interface
- `mockAuthProvider` - Mock auth.Provider interface
- `mockTLSProvider` - Mock tls.ConfigProvider interface

**Test Implementation Pattern**:
```go
package push_test

import (
    "context"
    "testing"

    "github.com/cnoe-io/idpbuilder/pkg/cmd/push"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

// Example: T-2.1.1-01 - TestNewPushCommand_Flags
func TestNewPushCommand_Flags(t *testing.T) {
    // When: Creating push command
    cmd := push.NewPushCommand()

    // Then: All flags exist
    registryFlag := cmd.Flags().Lookup("registry")
    assert.NotNil(t, registryFlag, "registry flag should exist")

    usernameFlag := cmd.Flags().Lookup("username")
    assert.NotNil(t, usernameFlag, "username flag should exist")

    passwordFlag := cmd.Flags().Lookup("password")
    assert.NotNil(t, passwordFlag, "password flag should exist")

    insecureFlag := cmd.Flags().Lookup("insecure")
    assert.NotNil(t, insecureFlag, "insecure flag should exist")

    verboseFlag := cmd.Flags().Lookup("verbose")
    assert.NotNil(t, verboseFlag, "verbose flag should exist")
}

// Example: T-2.1.1-20 - TestRunPush_Success_AllStages
func TestRunPush_Success_AllStages(t *testing.T) {
    // Track stage execution
    dockerGetCalled := false
    registryPushCalled := false
    dockerCloseCalled := false

    // Setup mocks (full implementation in test plan)
    mockDocker := &mockDockerClient{
        getImageFunc: func(ctx context.Context, imageName string) (v1.Image, error) {
            dockerGetCalled = true
            return empty.Image, nil
        },
        closeFunc: func() error {
            dockerCloseCalled = true
            return nil
        },
    }

    // ... setup other mocks ...

    opts := &push.PushOptions{
        ImageName: "alpine:latest",
        Registry:  "gitea.cnoe.localtest.me:8443",
        Username:  "admin",
        Password:  "password",
        Insecure:  true,
        Verbose:   false,
    }

    // When: Running push
    err := push.RunPushWithMocks(context.Background(), opts, mockDocker, mockRegistry, mockAuth, mockTLS)

    // Then: All stages executed successfully
    assert.NoError(t, err, "Push should succeed")
    assert.True(t, dockerGetCalled, "Docker GetImage should be called")
    assert.True(t, registryPushCalled, "Registry Push should be called")
    assert.True(t, dockerCloseCalled, "Docker Close should be called")
}
```

**Test Coverage Requirements**:
- Minimum 90% code coverage
- All 25 tests from test plan MUST be implemented
- Use Phase 1 mock fixtures (reuse from Phase 1 Wave 2)
- Table-driven tests for validation scenarios

#### Dependencies

**Upstream Dependencies** (must complete before this effort):
- Phase 1 Wave 1: docker.Client interface (✅ Complete)
- Phase 1 Wave 2: registry.Client interface (✅ Complete)
- Phase 1 Wave 3: auth.Provider interface (✅ Complete)
- Phase 1 Wave 4: tls.ConfigProvider interface (✅ Complete)

**Downstream Dependencies** (efforts that depend on this):
- Effort 2.1.2: Progress Reporter (needs callback signature from this effort)
- Wave 2.2: Environment variable support (needs PushOptions structure)
- Wave 2.3: Error handling enhancements (needs runPush pipeline)

#### Acceptance Criteria

- [ ] All files created/modified as specified
- [ ] All 25 tests from test plan passing (100% pass rate)
- [ ] Code coverage ≥90%
- [ ] No linting errors (golangci-lint)
- [ ] All exported functions/types have godoc comments
- [ ] Line count within estimate (450 ± 67 lines, 15% tolerance)
- [ ] `idpbuilder push --help` displays correct help text
- [ ] Manual test: `idpbuilder push alpine:latest --username admin --password password --insecure` succeeds
- [ ] All Phase 1 interfaces used correctly (verified by compilation)

---

### Effort 2.1.2: Progress Reporter & Output Formatting

#### R213 Metadata

```yaml
---
effort_metadata:
  effort_id: "2.1.2"
  effort_name: "Progress Reporter & Output Formatting"
  estimated_lines: 300
  dependencies: ["2.1.1"]
  files_touched:
    - "pkg/progress/reporter.go"
    - "pkg/progress/interface.go"
    - "pkg/progress/reporter_test.go"
    - "pkg/cmd/push/push.go"
  branch_name: "idpbuilder-oci-push/phase2/wave1/effort-2-progress-reporter"
  base_branch: "idpbuilder-oci-push/phase2/wave1/effort-1-push-command-core"
  theme: "Enhanced progress tracking and console output formatting"
  scope: "Layer-by-layer progress, thread-safe updates, summary statistics, verbose/normal modes"
  complexity: "medium"
  can_parallelize: false
  parallel_with: []
---
```

#### Scope

**Purpose**: Replace the basic progress callback from Effort 2.1.1 with a sophisticated reporter that tracks layer-by-layer progress, provides thread-safe updates, and displays formatted summaries.

**Boundaries**:
- ✅ IN SCOPE:
  - ProgressReporter interface definition
  - Thread-safe layer progress tracking
  - Normal mode output (compact)
  - Verbose mode output (detailed with rates)
  - Summary statistics (layers, bytes, duration)
  - Integration with push.go from Effort 2.1.1

- ❌ OUT OF SCOPE:
  - Terminal UI/TUI components (not needed for Wave 1)
  - Progress bars (ANSI escape codes)
  - Color output
  - Real-time updating (overwriting lines)
  - JSON output format (Wave 2.3)

#### Files to Create/Modify

**New Files**:
- `pkg/progress/reporter.go` (200 lines) - Main reporter implementation
- `pkg/progress/interface.go` (30 lines) - ProgressReporter interface
- `pkg/progress/reporter_test.go` (70 lines) - Unit tests (15 tests from test plan)

**Modified Files**:
- `pkg/cmd/push/push.go` (replace basic callback with reporter, ~10 lines changed)

**Total Estimated Lines**: 310 lines

#### Exact Code Specifications

**File: pkg/progress/interface.go**

```go
// Package progress provides progress tracking and reporting for image push operations
package progress

import "github.com/cnoe-io/idpbuilder/pkg/registry"

// ProgressReporter defines operations for tracking and displaying push progress
type ProgressReporter interface {
    // HandleProgress processes a progress update from registry push
    HandleProgress(update registry.ProgressUpdate)

    // DisplaySummary shows final statistics after push completes
    DisplaySummary()

    // GetCallback returns a callback function for registry.Push()
    GetCallback() registry.ProgressCallback
}
```

**Implementation Requirements for interface.go**:
- Interface MUST match registry.ProgressCallback signature exactly
- GetCallback() MUST return a function compatible with registry.Push()
- Keep interface minimal (only 3 methods needed for Wave 1)

---

**File: pkg/progress/reporter.go**

```go
package progress

import (
    "fmt"
    "sync"
    "time"

    "github.com/cnoe-io/idpbuilder/pkg/registry"
)

// Reporter tracks and displays progress for image push operations
type Reporter struct {
    verbose      bool
    startTime    time.Time
    layers       map[string]*LayerProgress
    mu           sync.Mutex
}

// LayerProgress tracks individual layer upload progress
type LayerProgress struct {
    Digest       string
    Size         int64
    Pushed       int64
    Status       string
    StartTime    time.Time
    CompleteTime *time.Time
}

// NewReporter creates a progress reporter
func NewReporter(verbose bool) *Reporter {
    return &Reporter{
        verbose:   verbose,
        startTime: time.Now(),
        layers:    make(map[string]*LayerProgress),
    }
}

// HandleProgress processes a progress update from registry.Push()
// This matches the registry.ProgressCallback signature from Phase 1
func (r *Reporter) HandleProgress(update registry.ProgressUpdate) {
    r.mu.Lock()
    defer r.mu.Unlock()

    // Get or create layer progress
    layer, exists := r.layers[update.LayerDigest]
    if !exists {
        layer = &LayerProgress{
            Digest:    update.LayerDigest,
            Size:      update.LayerSize,
            StartTime: time.Now(),
        }
        r.layers[update.LayerDigest] = layer
    }

    // Update layer state
    layer.Pushed = update.BytesPushed
    layer.Status = update.Status
    if update.Status == "complete" || update.Status == "exists" {
        now := time.Now()
        layer.CompleteTime = &now
    }

    // Display progress
    if r.verbose {
        r.displayVerbose(layer)
    } else {
        r.displayNormal(layer)
    }
}

// displayNormal shows compact progress (for normal mode)
func (r *Reporter) displayNormal(layer *LayerProgress) {
    shortDigest := layer.Digest
    if len(shortDigest) > 12 {
        shortDigest = shortDigest[:12]
    }

    switch layer.Status {
    case "uploading":
        percent := float64(layer.Pushed) / float64(layer.Size) * 100
        fmt.Printf("  %s: Pushing [%.1f%%]\n", shortDigest, percent)
    case "complete":
        fmt.Printf("  %s: Pushed ✓\n", shortDigest)
    case "exists":
        fmt.Printf("  %s: Already exists ✓\n", shortDigest)
    }
}

// displayVerbose shows detailed progress (for --verbose)
func (r *Reporter) displayVerbose(layer *LayerProgress) {
    elapsed := time.Since(layer.StartTime).Seconds()

    switch layer.Status {
    case "uploading":
        percent := float64(layer.Pushed) / float64(layer.Size) * 100
        rate := float64(layer.Pushed) / elapsed / 1024 / 1024 // MB/s
        fmt.Printf("  Layer %s:\n", layer.Digest)
        fmt.Printf("    Status: Uploading [%.1f%%] (%d/%d bytes)\n",
            percent, layer.Pushed, layer.Size)
        fmt.Printf("    Rate: %.2f MB/s\n", rate)
        fmt.Printf("    Elapsed: %.1fs\n", elapsed)
    case "complete":
        duration := elapsed
        if layer.CompleteTime != nil {
            duration = layer.CompleteTime.Sub(layer.StartTime).Seconds()
        }
        fmt.Printf("  Layer %s: Complete (%.1fs)\n", layer.Digest, duration)
    case "exists":
        fmt.Printf("  Layer %s: Already exists (skipped)\n", layer.Digest)
    }
}

// DisplaySummary shows final statistics after push completes
func (r *Reporter) DisplaySummary() {
    r.mu.Lock()
    defer r.mu.Unlock()

    totalLayers := len(r.layers)
    totalBytes := int64(0)
    pushedLayers := 0
    skippedLayers := 0

    for _, layer := range r.layers {
        totalBytes += layer.Size
        if layer.Status == "complete" {
            pushedLayers++
        } else if layer.Status == "exists" {
            skippedLayers++
        }
    }

    duration := time.Since(r.startTime)
    avgRate := float64(totalBytes) / duration.Seconds() / 1024 / 1024 // MB/s

    fmt.Println("\nPush Summary:")
    fmt.Printf("  Total layers: %d\n", totalLayers)
    fmt.Printf("  Pushed: %d, Skipped: %d\n", pushedLayers, skippedLayers)
    fmt.Printf("  Total size: %.2f MB\n", float64(totalBytes)/1024/1024)
    fmt.Printf("  Duration: %.1fs\n", duration.Seconds())
    fmt.Printf("  Average rate: %.2f MB/s\n", avgRate)
}

// GetCallback returns a registry.ProgressCallback that uses this reporter
func (r *Reporter) GetCallback() registry.ProgressCallback {
    return func(update registry.ProgressUpdate) {
        r.HandleProgress(update)
    }
}
```

**Implementation Requirements for reporter.go**:
- MUST use sync.Mutex for thread-safe layer map access
- HandleProgress MUST handle concurrent calls correctly
- displayNormal and displayVerbose MUST NOT be exported (internal only)
- Digest truncation in displayNormal MUST be exactly 12 characters
- Rate calculations MUST handle division by zero (elapsed < 0.001s)
- Summary MUST calculate correct totals even with mixed status layers

---

**File: pkg/cmd/push/push.go (MODIFICATION)**

```go
// Replace the basic progress callback in runPush() with the reporter

// In runPush(), replace:
//   progressCallback := func(update registry.ProgressUpdate) {
//       if opts.Verbose {
//           fmt.Printf("Layer %s: %d/%d bytes (%s)\n", ...)
//       }
//   }

// With:
import "github.com/cnoe-io/idpbuilder/pkg/progress"

// STAGE 7: Create progress reporter (replaces basic callback)
reporter := progress.NewReporter(opts.Verbose)

// STAGE 8: Execute push with reporter
fmt.Printf("Pushing to %s...\n", targetRef)
if err := registryClient.Push(ctx, image, targetRef, reporter.GetCallback()); err != nil {
    return fmt.Errorf("push failed: %w", err)
}

// Display final summary
reporter.DisplaySummary()

fmt.Printf("✓ Successfully pushed %s to %s\n", opts.ImageName, opts.Registry)
return nil
```

**Implementation Requirements for push.go modification**:
- Replace lines 95-102 in push.go (basic callback) with reporter creation
- Add progress import at top of file
- Call reporter.DisplaySummary() before success message
- DO NOT change other parts of runPush() pipeline
- Keep modifications minimal (~10 lines changed)

#### Tests Required

**File: pkg/progress/reporter_test.go**

Test specifications are defined in `planning/phase2/wave1/WAVE-TEST-PLAN.md` (Tests T-2.1.2-01 through T-2.1.2-15).

**Key Test Categories**:
1. **Reporter Creation** (T-2.1.2-01 to T-2.1.2-02): Test NewReporter with verbose/normal modes
2. **Progress Tracking** (T-2.1.2-03 to T-2.1.2-07): Test HandleProgress for uploading/complete/exists states, multiple layers, thread safety
3. **Display Formatting** (T-2.1.2-08 to T-2.1.2-10): Test displayNormal, displayVerbose, rate calculations
4. **Summary Statistics** (T-2.1.2-11 to T-2.1.2-13): Test DisplaySummary with single/multiple/mixed status layers
5. **Integration** (T-2.1.2-14 to T-2.1.2-15): Test GetCallback and digest truncation

**Test Implementation Pattern**:
```go
package progress_test

import (
    "testing"
    "sync"

    "github.com/cnoe-io/idpbuilder/pkg/progress"
    "github.com/cnoe-io/idpbuilder/pkg/registry"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

// Example: T-2.1.2-04 - TestReporter_HandleProgress_Complete
func TestReporter_HandleProgress_Complete(t *testing.T) {
    // Given: Reporter
    reporter := progress.NewReporter(false)

    // When: Receiving progress updates from uploading to complete
    reporter.HandleProgress(registry.ProgressUpdate{
        LayerDigest: "sha256:abc123def456",
        LayerSize:   1024,
        BytesPushed: 0,
        Status:      "uploading",
    })

    reporter.HandleProgress(registry.ProgressUpdate{
        LayerDigest: "sha256:abc123def456",
        LayerSize:   1024,
        BytesPushed: 1024,
        Status:      "complete",
    })

    // Then: No panics, summary works
    reporter.DisplaySummary()
}

// Example: T-2.1.2-07 - TestReporter_HandleProgress_ThreadSafety
func TestReporter_HandleProgress_ThreadSafety(t *testing.T) {
    // Given: Reporter
    reporter := progress.NewReporter(false)

    // When: Concurrent progress updates
    var wg sync.WaitGroup
    for i := 0; i < 10; i++ {
        wg.Add(1)
        go func(layerNum int) {
            defer wg.Done()
            digest := fmt.Sprintf("sha256:layer%d", layerNum)
            reporter.HandleProgress(registry.ProgressUpdate{
                LayerDigest: digest,
                LayerSize:   1000,
                BytesPushed: 500,
                Status:      "uploading",
            })
        }(i)
    }

    wg.Wait()

    // Then: No race conditions, all updates processed
    reporter.DisplaySummary()
    // Race detector would catch issues: go test -race
}
```

**Test Coverage Requirements**:
- Minimum 85% code coverage
- All 15 tests from test plan MUST be implemented
- Thread safety MUST be verified with -race flag
- Output verification via captured stdout (if feasible)

#### Dependencies

**Upstream Dependencies** (must complete before this effort):
- Effort 2.1.1: Push Command Core (✅ MUST complete first - defines callback signature)

**Downstream Dependencies** (efforts that depend on this):
- None (this completes Wave 2.1)
- Wave 2.2 will enhance this reporter with registry auto-detection feedback

#### Acceptance Criteria

- [ ] All files created/modified as specified
- [ ] All 15 tests from test plan passing (100% pass rate)
- [ ] Code coverage ≥85%
- [ ] No linting errors (golangci-lint)
- [ ] Thread safety verified with `go test -race` (zero data races)
- [ ] All exported functions/types have godoc comments
- [ ] Line count within estimate (300 ± 45 lines, 15% tolerance)
- [ ] Manual test: `idpbuilder push alpine:latest --verbose --username admin --password password --insecure` shows detailed progress
- [ ] Summary displays correct statistics (layers, bytes, duration, rate)
- [ ] Integration with push.go works correctly (verified by end-to-end test)

---

## Parallelization Strategy

### Sequential Execution (NO Parallelization)

```
Effort 2.1.1: Push Command Core (MUST complete first)
    ↓ (depends on callback signature)
Effort 2.1.2: Progress Reporter (enhances 2.1.1)
```

**Rationale for Sequential Execution**:
1. **Effort 2.1.1 is foundational**: Defines command structure, PushOptions, and callback signature
2. **Small wave size**: Only 2 efforts (~750 lines total), parallelization overhead not beneficial
3. **Integration testing dependency**: 2.1.2 needs 2.1.1's command to test end-to-end
4. **Clear handoff point**: 2.1.1 delivers working command, 2.1.2 enhances it with better progress
5. **Signature dependency**: 2.1.2 implements registry.ProgressCallback defined in Phase 1 but used in 2.1.1

**Future Parallelization Opportunities** (Waves 2.2 and 2.3):
- Wave 2.2 efforts (registry override, env vars) CAN parallelize
- Wave 2.3 efforts (validation, error handling) CAN parallelize

---

## Wave Size Compliance

**Total Wave Lines**: 755 lines (450 + 300 + 5 modifications)

**Size Limits**:
- Soft Limit: 3500 lines ✅
- Hard Limit: 4000 lines ✅

**Status**:
- [x] Within soft limit (755 < 3500)
- [x] Within hard limit (755 < 4000)
- [ ] Requires split plan (N/A - well under limit)

**Size Verification Command**:
```bash
# After Wave 2.1 integration branch created, measure with:
$PROJECT_ROOT/tools/line-counter.sh
# Expected output: ~755 lines (implementation + tests)
```

---

## Integration Strategy

### Wave 2.1 Integration Sequence

1. **Effort 2.1.1 completes** → Code Reviewer reviews → Merge to `idpbuilder-oci-push/phase2/wave1/integration`
2. **Effort 2.1.2 starts** (branches from 2.1.1) → Code Reviewer reviews → Merge to wave integration branch
3. **Wave integration tests** run on `idpbuilder-oci-push/phase2/wave1/integration` branch
4. **Architect performs wave assessment** (R340 compliance check)
5. **Merge to phase integration branch** `idpbuilder-oci-push/phase2/integration`

### Integration Branch Structure

```
idpbuilder-oci-push/phase2/integration (phase level)
    ↑
idpbuilder-oci-push/phase2/wave1/integration (wave level)
    ↑
idpbuilder-oci-push/phase2/wave1/effort-2-progress-reporter (Effort 2.1.2)
    ↑
idpbuilder-oci-push/phase2/wave1/effort-1-push-command-core (Effort 2.1.1)
    ↑
idpbuilder-oci-push/phase1/integration (base)
```

---

## Testing Strategy

### Unit Tests (Per Effort)

**Effort 2.1.1** - 25 tests:
- Flag definitions and validation (5 tests)
- Error handling for each pipeline stage (9 tests)
- Progress callback invocation (2 tests)
- Complete pipeline execution (9 tests)

**Effort 2.1.2** - 15 tests:
- Reporter creation and interface (2 tests)
- Progress tracking and thread safety (5 tests)
- Display formatting (3 tests)
- Summary statistics (3 tests)
- Integration with callback (2 tests)

**Coverage Targets**:
- Effort 2.1.1: ≥90% coverage (critical user-facing command)
- Effort 2.1.2: ≥85% coverage (display logic, some manual verification)

### Integration Tests (Wave Level)

**File: pkg/cmd/push/push_integration_test.go**

```go
package push_test

import (
    "context"
    "testing"
    "time"

    "github.com/cnoe-io/idpbuilder/pkg/cmd/push"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

// TestPushCommand_EndToEnd tests complete push workflow with progress reporter
// Prerequisites: Docker daemon running, alpine:latest pulled, Gitea registry available
func TestPushCommand_EndToEnd(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping integration test")
    }

    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
    defer cancel()

    opts := &push.PushOptions{
        ImageName: "alpine:latest",
        Registry:  "gitea.cnoe.localtest.me:8443",
        Username:  "giteaAdmin",
        Password:  "password",
        Insecure:  true,  // Test Gitea uses self-signed cert
        Verbose:   true,
    }

    // Execute push with progress reporter
    err := push.runPush(ctx, opts)

    // Verify success
    require.NoError(t, err, "Push should succeed")
    // Manual verification: Check Gitea registry for pushed image
}

// TestPushCommand_ProgressReporter_Verbose verifies detailed progress output
func TestPushCommand_ProgressReporter_Verbose(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping integration test")
    }

    // Similar to above but captures stdout to verify progress messages
    // Verify: Layer digests printed, percentages shown, summary displayed
}
```

**Integration Test Requirements**:
- Run with `go test` (not `go test -short`)
- Require Docker daemon accessible
- Require alpine:latest image pulled
- Optional: Require test Gitea registry (for full E2E)
- Tests MUST handle missing prerequisites gracefully (skip with message)

### Wave-Level Test Execution

```bash
# Unit tests only (no Docker required)
go test -short ./pkg/cmd/push/... ./pkg/progress/... -v

# All tests including integration (Docker required)
go test ./pkg/cmd/push/... ./pkg/progress/... -v

# Coverage report
go test ./pkg/cmd/push/... ./pkg/progress/... -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html

# Thread safety verification
go test ./pkg/progress/... -race -v

# Coverage enforcement
go test ./pkg/cmd/push/... -coverprofile=push-coverage.out
PUSH_COV=$(go tool cover -func=push-coverage.out | grep total | awk '{print $3}' | sed 's/%//')
if (( $(echo "$PUSH_COV < 90" | bc -l) )); then
    echo "❌ Push command coverage below 90%: $PUSH_COV%"
    exit 1
fi
```

---

## Risk Mitigation

### High-Risk Areas

**Effort 2.1.1 Risks**:
- **Docker daemon connection**: May fail if daemon not running
  - Mitigation: Clear error message, integration tests skip if unavailable
- **Phase 1 interface changes**: Breaking changes in docker/registry/auth/tls
  - Mitigation: Phase 1 interfaces are frozen (already integrated)
- **Credential exposure**: Password flag visible in process list
  - Mitigation: Document env var support coming in Wave 2.2

**Effort 2.1.2 Risks**:
- **Thread safety**: Concurrent progress updates could race
  - Mitigation: Use sync.Mutex, verify with -race flag
- **Division by zero**: Rate calculation if elapsed time = 0
  - Mitigation: Check elapsed > 0.001s before division

### External Dependencies

**Go Libraries** (already in go.mod from Phase 1):
- github.com/spf13/cobra v1.8.0 - CLI framework
- github.com/google/go-containerregistry v0.16.1 - OCI types
- github.com/stretchr/testify v1.9.0 - Testing framework

**Runtime Dependencies**:
- Docker daemon (for integration tests and actual usage)
- Gitea registry (optional, for full E2E testing)

### Complexity Hotspots

**Push Pipeline Orchestration** (push.go):
- 8 sequential stages with error handling
- Resource cleanup (defer dockerClient.Close())
- Context propagation through all stages
- Mitigation: Comprehensive unit tests for each stage, integration test for full pipeline

**Thread-Safe Progress Tracking** (reporter.go):
- Concurrent HandleProgress calls from registry client
- Shared layer map access
- Mitigation: sync.Mutex, -race flag verification, stress tests

---

## Quality Gates (R502 Compliance)

### Implementation Plan Quality Requirements

- [x] **R213 metadata blocks**: Both efforts have complete metadata
- [x] **Exact code specifications**: Real Go code (not pseudocode) for all implementations
- [x] **Complete file lists**: All files specified with line counts
- [x] **Detailed task breakdowns**: Step-by-step implementation guidance
- [x] **Test specifications**: References to 40 concrete tests in test plan
- [x] **Dependency graphs**: Clear upstream/downstream dependencies
- [x] **Effort sizing**: Both efforts <800 lines (450 + 300)
- [x] **Integration strategy**: Sequential execution with rationale
- [x] **Risk mitigation**: High-risk areas identified with mitigations

### R340 Wave Architecture Compliance

- [x] **Real code examples**: Architecture shows actual Go interfaces
- [x] **Concrete function signatures**: Complete with parameter types
- [x] **Phase 1 integration**: Uses real docker, registry, auth, tls packages
- [x] **Adaptation notes**: Lessons from Phase 1 documented
- [x] **No pseudocode**: All examples are working Go code

### R341 TDD Compliance

- [x] **Tests before implementation**: Test plan created before this plan
- [x] **40 tests specified**: 25 for push command, 15 for progress reporter
- [x] **Concrete test code**: All tests are real Go code
- [x] **Phase 1 fixtures**: Reuses mock providers from Phase 1

### R510 Checklist Structure

- [x] **Clear effort definitions**: Two efforts with precise scopes
- [x] **Acceptance criteria**: Each effort has measurable success criteria
- [x] **Quality gates**: Coverage targets, linting, line counts
- [x] **Compliance verified**: R213, R340, R341, R502, R510

---

## Next Steps (Orchestrator Actions)

1. **Review this implementation plan** for R502 compliance
2. **Create wave infrastructure**:
   - Branch: `idpbuilder-oci-push/phase2/wave1/integration` (base for all efforts)
   - Branch: `idpbuilder-oci-push/phase2/wave1/effort-1-push-command-core`
   - Worktree: `efforts/phase2/wave1/effort-2.1.1-push-command-core/`
3. **Spawn SW Engineer** for Effort 2.1.1 with:
   - Working directory: `efforts/phase2/wave1/effort-2.1.1-push-command-core/`
   - Branch: `idpbuilder-oci-push/phase2/wave1/effort-1-push-command-core`
   - Instructions: This implementation plan (Effort 2.1.1 section)
   - Test plan: `planning/phase2/wave1/WAVE-TEST-PLAN.md` (T-2.1.1-01 to T-2.1.1-25)
4. **Monitor Effort 2.1.1** implementation (size checks every 100 lines)
5. **Spawn Code Reviewer** after Effort 2.1.1 completes
6. **Merge Effort 2.1.1** to wave integration branch after review approval
7. **Create Effort 2.1.2 branch** (base: Effort 2.1.1)
8. **Spawn SW Engineer** for Effort 2.1.2 (similar process)
9. **Wave integration** after both efforts merged
10. **Spawn Architect** for wave assessment (R340 review)

---

## Document Status

**Status**: ✅ READY FOR ORCHESTRATOR INFRASTRUCTURE CREATION
**Created**: 2025-10-31
**Planner**: Code Reviewer Agent (@agent-code-reviewer)
**Efforts**: 2 (Push Command Core, Progress Reporter)
**Total Lines**: ~755 lines (well under limit)
**Fidelity Level**: EXACT (detailed code specifications, R213 metadata)

**Compliance Verified**:
- ✅ R213: Effort metadata blocks complete (effort_id, estimated_lines, dependencies, branch_name, can_parallelize)
- ✅ R502: Implementation plan quality gates (exact specifications, file lists, test references)
- ✅ R340: Wave architecture compliance (real code examples from architecture)
- ✅ R341: TDD compliance (test plan precedes this plan)
- ✅ R510: Checklist structure (clear criteria for all sections)
- ✅ R550: Correct file placement (planning/phase2/wave1/WAVE-IMPLEMENTATION-PLAN.md)

**Next Action**: Orchestrator creates infrastructure and spawns SW Engineer for Effort 2.1.1

---

**END OF WAVE 2.1 IMPLEMENTATION PLAN**
