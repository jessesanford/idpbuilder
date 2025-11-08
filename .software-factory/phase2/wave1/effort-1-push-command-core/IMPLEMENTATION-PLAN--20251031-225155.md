# Effort 2.1.1: Push Command Core & Pipeline Orchestration - Implementation Plan

**Status**: READY FOR IMPLEMENTATION
**Created**: 2025-10-31T22:51:55Z
**Planner**: Code Reviewer Agent (@agent-code-reviewer)
**Fidelity Level**: EXACT SPECIFICATIONS

---

## 🚨 CRITICAL EFFORT METADATA (R213 - FROM WAVE PLAN)

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

---

## Overview

**Purpose**: Create the foundational `idpbuilder push` command that orchestrates the complete image push pipeline by integrating all Phase 1 packages (docker, registry, auth, tls).

**Wave Context**: Phase 2, Wave 1 - First effort in command implementation wave
**Dependencies**: Phase 1 complete (all interfaces frozen and tested)
**Downstream**: Effort 2.1.2 (Progress Reporter) depends on callback signature defined here

---

## Scope Definition

### ✅ IN SCOPE (What this effort MUST deliver)

1. **Cobra Command Registration**
   - Create `NewPushCommand()` function
   - Register with IDPBuilder root command
   - Define command structure (Use, Short, Long, Args, RunE)

2. **Flag Definitions**
   - `--registry` (default: gitea.cnoe.localtest.me:8443)
   - `--username` (required)
   - `--password` (required)
   - `--insecure/-k` (boolean, default: false)
   - `--verbose` (boolean, default: false)

3. **8-Stage Pipeline Orchestration**
   - Stage 1: Initialize Docker client
   - Stage 2: Retrieve image from Docker daemon
   - Stage 3: Setup authentication
   - Stage 4: Setup TLS configuration
   - Stage 5: Create registry client
   - Stage 6: Build target reference
   - Stage 7: Create progress callback (simple implementation)
   - Stage 8: Execute push

4. **Phase 1 Integration**
   - Use `docker.Client` interface
   - Use `registry.Client` interface
   - Use `auth.Provider` interface
   - Use `tls.ConfigProvider` interface

5. **Error Handling**
   - Wrap errors with context using %w
   - Early returns on failure
   - Resource cleanup (defer dockerClient.Close())

6. **Progress Callback**
   - Simple callback that prints to console in verbose mode
   - Defines callback signature for Effort 2.1.2 to enhance

7. **Testing**
   - 25 unit tests (T-2.1.1-01 through T-2.1.1-25)
   - Mock-based testing using Phase 1 fixtures
   - ≥90% code coverage

### ❌ OUT OF SCOPE (What this effort MUST NOT include)

1. **Advanced Progress Reporter** - That's Effort 2.1.2
2. **Environment Variable Support** - Wave 2.2
3. **Registry Auto-Detection** - Wave 2.2
4. **Sophisticated Error Formatting** - Wave 2.3
5. **Exit Code Mapping** - Wave 2.3
6. **JSON Output Format** - Wave 2.3
7. **Configuration File Support** - Future wave

---

## File Structure

### New Files to Create

#### 1. pkg/cmd/push/push.go (300 lines)
**Purpose**: Main command implementation with 8-stage pipeline

**Imports Required**:
```go
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
```

**Functions to Implement**:
- `NewPushCommand() *cobra.Command` - Command factory
- `runPush(ctx context.Context, opts *PushOptions) error` - Pipeline orchestrator
- `truncateDigest(digest string, length int) string` - Helper for display

#### 2. pkg/cmd/push/types.go (50 lines)
**Purpose**: PushOptions struct and validation

**Content**:
```go
package push

import "fmt"

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

#### 3. pkg/cmd/push/push_test.go (100 lines)
**Purpose**: Unit tests for push command (25 tests)

**Test Structure**:
```go
package push_test

import (
    "context"
    "testing"

    "github.com/cnoe-io/idpbuilder/pkg/cmd/push"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

// Mock fixtures from Phase 1 Wave 2 test patterns
type mockDockerClient struct { /* ... */ }
type mockRegistryClient struct { /* ... */ }
type mockAuthProvider struct { /* ... */ }
type mockTLSProvider struct { /* ... */ }

// Test categories:
// - Flag tests (T-2.1.1-01 to T-2.1.1-03)
// - Validation tests (T-2.1.1-04 to T-2.1.1-07)
// - Error handling tests (T-2.1.1-08 to T-2.1.1-14)
// - Progress tests (T-2.1.1-15 to T-2.1.1-16)
// - Pipeline tests (T-2.1.1-17 to T-2.1.1-23)
// - Integration tests (T-2.1.1-24 to T-2.1.1-25)
```

### Modifications to Existing Files

#### 4. pkg/cmd/root.go (+5 lines)
**Purpose**: Register push command with IDPBuilder root

**Changes**:
```go
// ADD THIS IMPORT:
import "github.com/cnoe-io/idpbuilder/pkg/cmd/push"

// IN Execute() function, ADD THIS LINE after existing commands:
rootCmd.AddCommand(push.NewPushCommand())
```

**CRITICAL**:
- ONLY add the import and AddCommand line
- DO NOT modify existing command registrations
- DO NOT change rootCmd structure
- Keep modifications minimal (exactly 5 lines total)

---

## Implementation Steps (Sequential)

### Step 1: Create Package Structure (5 minutes)

```bash
# Create package directory
mkdir -p pkg/cmd/push

# Verify directory structure
ls -la pkg/cmd/push
```

**Verification**: Directory exists and is empty

---

### Step 2: Implement types.go (15 minutes)

**File**: `pkg/cmd/push/types.go`

**Implementation Order**:
1. Create package declaration
2. Define PushOptions struct with all 6 fields
3. Add godoc comments for each field
4. Implement Validate() method
5. Test each validation case

**Code Template**:
```go
package push

import "fmt"

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

**Verification**:
- All fields exported (capitalized)
- Godoc comments present
- Validate() covers all required fields

---

### Step 3: Implement push.go - Command Factory (30 minutes)

**File**: `pkg/cmd/push/push.go`

**Implementation Order**:
1. Package declaration and imports
2. NewPushCommand() function skeleton
3. Define command structure (Use, Short, Long, Args)
4. Add flag definitions
5. Mark required flags
6. Wire RunE to runPush (stub for now)

**Code Template**:
```go
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
```

**Verification**:
- Command compiles
- All 5 flags defined
- Required flags marked
- Help text complete

---

### Step 4: Implement push.go - Pipeline Orchestrator (60 minutes)

**File**: `pkg/cmd/push/push.go` (add runPush function)

**Implementation Order**:
1. STAGE 1: Initialize Docker client (with error handling)
2. STAGE 2: Retrieve image (with defer Close)
3. STAGE 3: Setup authentication (with validation)
4. STAGE 4: Setup TLS
5. STAGE 5: Create registry client
6. STAGE 6: Build target reference
7. STAGE 7: Create progress callback (simple version)
8. STAGE 8: Execute push

**Code Template**:
```go
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

**Critical Requirements**:
- MUST use Phase 1 interfaces exactly as shown
- Error wrapping MUST preserve original error types using %w
- Pipeline stages MUST execute in order (1-8) with early returns on error
- defer dockerClient.Close() MUST be called to ensure cleanup
- Progress callback MUST check opts.Verbose before printing
- MUST NOT implement advanced progress features (that's Effort 2.1.2)

**Verification**:
- All 8 stages implemented
- Error wrapping consistent
- Resource cleanup present
- Progress callback functional

---

### Step 5: Implement Tests (90 minutes)

**File**: `pkg/cmd/push/push_test.go`

**Test Implementation Order**:

1. **Setup Mock Fixtures** (20 minutes)
   - Copy mock providers from Phase 1 Wave 2 test patterns
   - `mockDockerClient`
   - `mockRegistryClient`
   - `mockAuthProvider`
   - `mockTLSProvider`

2. **Flag Tests** (10 minutes) - T-2.1.1-01 to T-2.1.1-03
   ```go
   func TestNewPushCommand_Flags(t *testing.T) { /* ... */ }
   func TestNewPushCommand_FlagDefaults(t *testing.T) { /* ... */ }
   func TestNewPushCommand_RequiredFlags(t *testing.T) { /* ... */ }
   ```

3. **Validation Tests** (15 minutes) - T-2.1.1-04 to T-2.1.1-07
   ```go
   func TestPushOptions_Validate_Valid(t *testing.T) { /* ... */ }
   func TestPushOptions_Validate_MissingImage(t *testing.T) { /* ... */ }
   func TestPushOptions_Validate_MissingUsername(t *testing.T) { /* ... */ }
   func TestPushOptions_Validate_MissingPassword(t *testing.T) { /* ... */ }
   ```

4. **Error Handling Tests** (25 minutes) - T-2.1.1-08 to T-2.1.1-14
   - Docker connection errors
   - Image not found errors
   - Authentication errors
   - Registry client errors
   - Push failures

5. **Progress Tests** (10 minutes) - T-2.1.1-15 to T-2.1.1-16
   - Callback invocation
   - Nil callback handling

6. **Pipeline Tests** (15 minutes) - T-2.1.1-17 to T-2.1.1-23
   - Context cancellation
   - Insecure mode
   - Verbose mode
   - Success path
   - Docker close called
   - Error wrapping
   - Custom registry

7. **Integration Tests** (5 minutes) - T-2.1.1-24 to T-2.1.1-25
   - Cobra integration
   - Help text

**Test Pattern Example**:
```go
package push_test

import (
    "context"
    "testing"

    "github.com/cnoe-io/idpbuilder/pkg/cmd/push"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

// Example: T-2.1.1-01
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

// Example: T-2.1.1-20
func TestRunPush_Success_AllStages(t *testing.T) {
    // Track stage execution
    dockerGetCalled := false
    registryPushCalled := false
    dockerCloseCalled := false

    // Setup mocks
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

    mockRegistry := &mockRegistryClient{
        buildRefFunc: func(registryURL, imageName string) (string, error) {
            return registryURL + "/" + imageName, nil
        },
        pushFunc: func(ctx context.Context, image v1.Image, targetRef string, callback registry.ProgressCallback) error {
            registryPushCalled = true
            return nil
        },
    }

    mockAuth := &mockAuthProvider{username: "admin", password: "password"}
    mockTLS := &mockTLSProvider{insecure: true}

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

**Verification**:
```bash
# Run tests
go test ./pkg/cmd/push/... -v

# Check coverage
go test ./pkg/cmd/push/... -coverprofile=coverage.out
go tool cover -func=coverage.out | grep total
# Must show ≥90%
```

---

### Step 6: Register Command with Root (10 minutes)

**File**: `pkg/cmd/root.go`

**Changes**:
1. Add import: `"github.com/cnoe-io/idpbuilder/pkg/cmd/push"`
2. In Execute() function, add: `rootCmd.AddCommand(push.NewPushCommand())`

**Location**: After existing commands (create, get, delete, version)

**Verification**:
```bash
# Build and check help
go build -o idpbuilder ./...
./idpbuilder push --help
# Should display push command help text
```

---

## Size Management

### Line Count Tracking

**Tool**: `$PROJECT_ROOT/tools/line-counter.sh`

**Measurement Points**:
1. After implementing types.go (~50 lines)
2. After implementing NewPushCommand (~100 lines)
3. After implementing runPush (~200 lines)
4. After implementing tests (~100 lines)
5. Final measurement before review

**Commands**:
```bash
# Find project root
PROJECT_ROOT=$(pwd)
while [ "$PROJECT_ROOT" != "/" ]; do
    [ -f "$PROJECT_ROOT/orchestrator-state-v3.json" ] && break
    PROJECT_ROOT=$(dirname "$PROJECT_ROOT")
done

# Measure implementation lines
cd /path/to/effort/directory
$PROJECT_ROOT/tools/line-counter.sh

# Expected output: ~450 lines
```

**Thresholds**:
- ✅ <700 lines: COMPLIANT (continue)
- ⚠️ 700-800 lines: WARNING (approaching limit, finish and review)
- ❌ ≥800 lines: STOP IMMEDIATELY (requires split plan)

---

## Testing Strategy

### Unit Tests (No External Dependencies)

**Run Command**:
```bash
go test -short ./pkg/cmd/push/... -v
```

**Expected Results**:
- 25 tests passing
- 0 failures
- <5 seconds execution time

### Integration Tests (Requires Docker Daemon)

**Prerequisites**:
```bash
# Verify Docker daemon
docker info

# Pull test image
docker pull alpine:latest
```

**Run Command**:
```bash
go test ./pkg/cmd/push/... -v
```

**Expected Results**:
- 25 tests passing (including integration)
- Docker connection successful
- Image retrieval successful

### Coverage Verification

**Command**:
```bash
go test ./pkg/cmd/push/... -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html

# Check coverage percentage
COVERAGE=$(go tool cover -func=coverage.out | grep total | awk '{print $3}' | sed 's/%//')
if (( $(echo "$COVERAGE < 90" | bc -l) )); then
    echo "❌ Coverage below 90%: $COVERAGE%"
    exit 1
else
    echo "✅ Coverage: $COVERAGE%"
fi
```

**Minimum**: 90% coverage required

---

## Phase 1 Integration

### Docker Client Integration

**Interface** (from Phase 1):
```go
type docker.Client interface {
    GetImage(ctx context.Context, imageName string) (v1.Image, error)
    Close() error
}
```

**Usage in push.go**:
```go
dockerClient, err := docker.NewClient()
if err != nil {
    return fmt.Errorf("failed to connect to Docker daemon: %w", err)
}
defer dockerClient.Close()

image, err := dockerClient.GetImage(ctx, opts.ImageName)
if err != nil {
    // Error types from Phase 1:
    // - docker.ImageNotFoundError
    // - docker.DaemonConnectionError
    return fmt.Errorf("failed to get image: %w", err)
}
```

---

### Registry Client Integration

**Interface** (from Phase 1):
```go
type registry.Client interface {
    Push(ctx context.Context, image v1.Image, targetRef string, callback ProgressCallback) error
    BuildImageReference(registryURL, imageName string) (string, error)
}
```

**Usage in push.go**:
```go
registryClient, err := registry.NewClient(authProvider, tlsProvider)
if err != nil {
    return fmt.Errorf("failed to create registry client: %w", err)
}

targetRef, err := registryClient.BuildImageReference(opts.Registry, opts.ImageName)
if err != nil {
    return fmt.Errorf("invalid registry or image name: %w", err)
}

err = registryClient.Push(ctx, image, targetRef, progressCallback)
if err != nil {
    // Error types from Phase 1:
    // - registry.AuthenticationError
    // - registry.NetworkError
    // - registry.PushFailedError
    return fmt.Errorf("push failed: %w", err)
}
```

---

### Authentication Integration

**Interface** (from Phase 1):
```go
type auth.Provider interface {
    GetAuthenticator() (authn.Authenticator, error)
    ValidateCredentials() error
}
```

**Usage in push.go**:
```go
authProvider := auth.NewBasicAuthProvider(opts.Username, opts.Password)
if err := authProvider.ValidateCredentials(); err != nil {
    // Error types from Phase 1:
    // - auth.ValidationError
    return fmt.Errorf("invalid credentials: %w", err)
}

// Pass to registry client
registryClient, err := registry.NewClient(authProvider, tlsProvider)
```

---

### TLS Configuration Integration

**Interface** (from Phase 1):
```go
type tls.ConfigProvider interface {
    GetTLSConfig() *tls.Config
    IsInsecure() bool
}
```

**Usage in push.go**:
```go
tlsProvider := tls.NewConfigProvider(opts.Insecure)

// Pass to registry client
registryClient, err := registry.NewClient(authProvider, tlsProvider)
```

---

## Dependencies

### External Libraries (Already in go.mod from Phase 1)

```go
require (
    github.com/spf13/cobra v1.8.0
    github.com/google/go-containerregistry v0.16.1
    github.com/stretchr/testify v1.9.0
)
```

### Internal Dependencies (Phase 1 - Complete and Frozen)

- `pkg/docker` - Docker client interface (31 tests, 85%+ coverage)
- `pkg/registry` - Registry client interface (31 tests, 85%+ coverage)
- `pkg/auth` - Authentication provider interface (31 tests, 85%+ coverage)
- `pkg/tls` - TLS configuration provider interface (10 tests, 90%+ coverage)

---

## Acceptance Criteria

### Implementation Checklist

- [ ] All 4 files created/modified as specified
  - [ ] pkg/cmd/push/push.go (300 lines)
  - [ ] pkg/cmd/push/types.go (50 lines)
  - [ ] pkg/cmd/push/push_test.go (100 lines)
  - [ ] pkg/cmd/root.go (+5 lines modification)

- [ ] All 25 tests from test plan implemented and passing
  - [ ] T-2.1.1-01 to T-2.1.1-03: Flag tests
  - [ ] T-2.1.1-04 to T-2.1.1-07: Validation tests
  - [ ] T-2.1.1-08 to T-2.1.1-14: Error handling tests
  - [ ] T-2.1.1-15 to T-2.1.1-16: Progress tests
  - [ ] T-2.1.1-17 to T-2.1.1-23: Pipeline tests
  - [ ] T-2.1.1-24 to T-2.1.1-25: Integration tests

- [ ] Code quality
  - [ ] Code coverage ≥90%
  - [ ] No linting errors (golangci-lint)
  - [ ] All exported functions/types have godoc comments
  - [ ] Line count within estimate (450 ± 67 lines, 15% tolerance)

- [ ] Manual testing
  - [ ] `idpbuilder push --help` displays correct help text
  - [ ] `idpbuilder push alpine:latest --username admin --password password --insecure` succeeds
  - [ ] All Phase 1 interfaces used correctly (verified by compilation)

- [ ] Integration verification
  - [ ] Command registered with root
  - [ ] All flags functional
  - [ ] Error messages clear and actionable
  - [ ] Progress callback invoked (verbose mode)

---

## Risk Mitigation

### High-Risk Areas

1. **Docker Daemon Connection**
   - **Risk**: May fail if daemon not running
   - **Mitigation**: Clear error message, integration tests skip if unavailable
   - **Detection**: Test early with `docker info`

2. **Phase 1 Interface Changes**
   - **Risk**: Breaking changes in docker/registry/auth/tls
   - **Mitigation**: Phase 1 interfaces are frozen (already integrated)
   - **Detection**: Compile-time errors

3. **Credential Exposure**
   - **Risk**: Password flag visible in process list
   - **Mitigation**: Document env var support coming in Wave 2.2
   - **Detection**: Manual inspection of running processes

4. **Test Coverage**
   - **Risk**: May not achieve 90% target
   - **Mitigation**: Write tests FIRST (TDD), use table-driven tests
   - **Detection**: Coverage tool after each test batch

---

## Quality Gates (R502 Compliance)

### Implementation Plan Quality

- [x] **R213 metadata**: Complete effort metadata from wave plan
- [x] **Exact code specifications**: Real Go code (not pseudocode) for all implementations
- [x] **Complete file lists**: All 4 files specified with line counts
- [x] **Detailed task breakdowns**: 6 sequential implementation steps
- [x] **Test specifications**: References to 25 concrete tests in test plan
- [x] **Dependency documentation**: All Phase 1 integrations documented
- [x] **Scope clarity**: Clear IN SCOPE and OUT OF SCOPE sections
- [x] **Size management**: Line counting strategy with thresholds

---

## Downstream Dependencies

**Effort 2.1.2 (Progress Reporter) depends on:**
1. **Callback signature**: `func(update registry.ProgressUpdate)`
2. **PushOptions structure**: Verbose flag for mode selection
3. **runPush pipeline**: Stage 7 integration point for enhanced reporter

**Wave 2.2 efforts will enhance:**
- Environment variable support for credentials
- Registry auto-detection logic
- Additional flag options

**Wave 2.3 efforts will enhance:**
- Validation improvements
- Error message formatting
- Exit code mapping

---

## Next Steps (After Implementation)

1. **Size Measurement**: Run line-counter.sh to verify <800 lines
2. **Code Review**: Spawn Code Reviewer agent for review
3. **Integration**: Merge to wave integration branch after approval
4. **Handoff**: Effort 2.1.2 can begin (branches from this effort)

---

## Document Status

**Status**: ✅ READY FOR SOFTWARE ENGINEER
**Created**: 2025-10-31T22:51:55Z
**Planner**: Code Reviewer Agent (@agent-code-reviewer)
**Estimated Lines**: 450 lines
**Fidelity Level**: EXACT (step-by-step implementation guidance)

**Compliance Verified**:
- ✅ R213: Effort metadata complete
- ✅ R303: Timestamped filename (IMPLEMENTATION-PLAN--20251031-225155.md)
- ✅ R383: Stored in .software-factory/phase2/wave1/effort-1-push-command-core/
- ✅ R502: Implementation plan quality gates (exact specifications)
- ✅ R341: References test plan (25 tests)

**Next Action**: Orchestrator spawns Software Engineer with this plan

---

**🚨 CRITICAL R405 AUTOMATION FLAG 🚨**

CONTINUE-SOFTWARE-FACTORY=TRUE

---

**END OF IMPLEMENTATION PLAN**
