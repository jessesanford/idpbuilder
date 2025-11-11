# Effort 1.1.4 Implementation Plan: Push Command Scaffolding

**Effort**: Effort 1.1.4 - Command Structure Definition
**Phase**: Phase 1 - Foundation & Interfaces
**Wave**: Wave 1.1 - Interface Definitions
**Created**: 2025-11-11 20:48:29 UTC
**Planner**: Code Reviewer Agent
**Status**: Ready for Implementation

---

## EFFORT INFRASTRUCTURE METADATA (FROM WAVE PLAN)

### R213 Metadata

```json
{
  "effort_id": "1.1.4",
  "effort_name": "Command Structure Definition",
  "branch_name": "idpbuilder-oci-push-rebuild/phase1-wave1-effort-1.1.4",
  "base_branch": "idpbuilder-oci-push-rebuild/phase1-wave1-effort-1.1.3",
  "parent_wave": "wave1.1",
  "parent_phase": "phase1",
  "depends_on": ["1.1.3"],
  "estimated_lines": 190,
  "complexity": "low",
  "can_parallelize": false,
  "parallel_with": [],
  "tests_required": [
    "T1.1.4-001: PushCommand struct compiles",
    "T1.1.4-002: PushFlags struct compiles",
    "T1.1.4-003: NewPushCommand creates valid Cobra command",
    "T1.1.4-004: All command flags defined correctly",
    "T1.1.4-005: Exit codes constants defined",
    "T1.1.4-006: runPush function signature valid",
    "T1.1.4-007: Command registers with Cobra successfully"
  ]
}
```

### Parallelization Info (FROM WAVE PLAN)
- **Can Parallelize**: false
- **Parallel With**: [] (none - sequential execution required)
- **Reason**: Command structure depends on all previous interface definitions (1.1.1-1.1.3)

### Size Estimate (FROM WAVE PLAN)
- **Estimated Lines**: 190 lines
- **File Count**: 1 file (cmd/push.go)
- **Test Files**: 1 file (cmd/push_test.go)
- **Within Limit**: ✅ Yes (190 < 800 line hard limit per R535)

---

## Overview

### Purpose

Create the PushCommand structure that will orchestrate all packages (Docker, Registry, Auth, TLS) in the push workflow. This is the CLI entry point for the `idpbuilder push` feature.

**This is SCAFFOLDING ONLY** - no actual push implementation logic. The command structure defines:
- Interface dependencies (DockerClient, RegistryClient, AuthProvider, TLSProvider)
- Command-line flags (registry, username, password, insecure, verbose)
- Cobra command registration
- Exit code constants
- Stub runPush() function (prints message, returns nil)

### Scope and Boundaries

**What IS in Scope** (✅):
- ✅ Command structure definition (PushCommand, PushFlags)
- ✅ Cobra command registration with flags
- ✅ runPush function signature (stub implementation - prints message only)
- ✅ Exit code constants (5 codes defined)
- ✅ Complete Go documentation comments
- ✅ Interface field declarations (no implementations)
- ✅ Flag validation (password marked as required)

**What is OUT of Scope** (❌):
- ❌ NO actual push workflow implementation (Phase 2 Wave 1)
- ❌ NO error handling logic (Phase 2 Wave 1)
- ❌ NO progress reporting (Phase 2 Wave 1)
- ❌ NO Docker client instantiation (Phase 2 Wave 1)
- ❌ NO Registry client instantiation (Phase 2 Wave 1)
- ❌ NO Auth provider instantiation (Phase 2 Wave 1)
- ❌ NO TLS provider instantiation (Phase 2 Wave 1)

### Dependencies

**Upstream Dependencies** (MUST complete before this effort):
- ✅ Effort 1.1.1 (Docker Interface) - imports `pkg/docker.DockerClient`
- ✅ Effort 1.1.2 (Registry Interface) - imports `pkg/registry.RegistryClient`
- ✅ Effort 1.1.3 (Auth & TLS Interfaces) - imports `pkg/auth.AuthProvider` and `pkg/tls.TLSProvider`

**Critical**: All three interface packages (docker, registry, auth, tls) MUST exist before implementing this effort. The code will not compile without these dependencies.

**Downstream Dependencies** (efforts that depend on this):
- None (this is the final effort in Wave 1)

**External Dependencies**:
- `github.com/spf13/cobra` - Cobra CLI framework (already in IDPBuilder go.mod)

---

## File Structure

### Files to Create

**Total Files**: 1 production file + 1 test file

#### Production Code

**File: cmd/push.go** (190 lines)
- PushCommand struct (4 interface fields: dockerClient, registryClient, authProvider, tlsProvider)
- PushFlags struct (5 flag fields: Registry, Username, Password, Insecure, Verbose)
- NewPushCommand() function (creates and configures Cobra command with flags)
- runPush() function (stub - prints "not yet implemented" message)
- Exit code constants (5 values: ExitSuccess through ExitImageNotFound)
- Complete package documentation

#### Test Code

**File: cmd/push_test.go** (created by SW Engineer in same effort)
- 7 test functions covering all requirements (see Test Requirements section)
- Validates struct compilation
- Validates Cobra command creation
- Validates all flags defined
- Validates exit codes
- Validates command registration

---

## Implementation Steps

### Step 1: Create cmd/push.go File Structure

**Action**: Create the file `cmd/push.go` with package declaration and imports.

**Exact Code to Copy** (from WAVE-1.1-ARCHITECTURE.md lines 736-749):

```go
// Package cmd provides the CLI commands for idpbuilder.
package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/jessesanford/idpbuilder/pkg/auth"
	"github.com/jessesanford/idpbuilder/pkg/docker"
	"github.com/jessesanford/idpbuilder/pkg/registry"
	"github.com/jessesanford/idpbuilder/pkg/tls"
	"github.com/spf13/cobra"
)
```

**Critical**:
- Package MUST be `cmd` (not `main`)
- Import paths MUST match exact package locations from Efforts 1.1.1-1.1.3
- All 4 interface packages MUST be imported (docker, registry, auth, tls)
- Cobra import MUST be `github.com/spf13/cobra`

---

### Step 2: Define PushCommand Struct

**Action**: Define the PushCommand struct with 4 interface fields.

**Exact Code to Copy** (from WAVE-1.1-ARCHITECTURE.md lines 751-765):

```go
// PushCommand orchestrates the image push workflow.
// It coordinates Docker, Registry, Auth, and TLS packages to push an image to a registry.
type PushCommand struct {
	// dockerClient provides access to the local Docker daemon.
	dockerClient docker.DockerClient

	// registryClient handles OCI registry operations.
	registryClient registry.RegistryClient

	// authProvider supplies authentication credentials.
	authProvider auth.AuthProvider

	// tlsProvider generates TLS configurations.
	tlsProvider tls.TLSProvider
}
```

**Critical**:
- Struct fields are INTERFACES, not concrete types
- All 4 fields are PRIVATE (lowercase first letter)
- Each field has a documentation comment
- NO constructor function for PushCommand (not needed in Wave 1)

---

### Step 3: Define PushFlags Struct

**Action**: Define the PushFlags struct with 5 flag fields.

**Exact Code to Copy** (from WAVE-1.1-ARCHITECTURE.md lines 767-783):

```go
// PushFlags holds command-line flag values for the push command.
type PushFlags struct {
	// Registry is the target registry URL (e.g., "https://gitea.cnoe.localtest.me:8443")
	Registry string

	// Username for registry authentication
	Username string

	// Password for registry authentication
	Password string

	// Insecure enables TLS certificate verification bypass (for self-signed certs)
	Insecure bool

	// Verbose enables detailed progress output
	Verbose bool
}
```

**Critical**:
- All 5 fields are PUBLIC (uppercase first letter) - needed for flag binding
- Each field has a documentation comment
- Field types MUST be: 3 strings, 2 bools

---

### Step 4: Implement NewPushCommand Function

**Action**: Create the NewPushCommand() function that builds the Cobra command with all flags.

**Exact Code to Copy** (from WAVE-1.1-ARCHITECTURE.md lines 785-835):

```go
// NewPushCommand creates and configures the push command.
//
// Example registration with Cobra:
//   rootCmd.AddCommand(cmd.NewPushCommand())
func NewPushCommand() *cobra.Command {
	flags := &PushFlags{}

	cmd := &cobra.Command{
		Use:   "push IMAGE_NAME",
		Short: "Push a container image to the IDPBuilder registry",
		Long: `Push a container image to the IDPBuilder Gitea registry.

The image must exist in your local Docker daemon. Use 'docker images' to list available images.

Examples:
  # Push with default registry
  idpbuilder push myapp:latest --password mypassword

  # Push to custom registry
  idpbuilder push myapp:v1.0 --registry https://registry.example.com --username admin --password secret

  # Push with insecure TLS (for self-signed certificates)
  idpbuilder push myapp:dev --insecure --password mypassword

  # Push with verbose output
  idpbuilder push myapp:latest --verbose --password mypassword
`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			imageName := args[0]
			return runPush(cmd.Context(), imageName, flags)
		},
	}

	// Define flags
	cmd.Flags().StringVar(&flags.Registry, "registry", "https://gitea.cnoe.localtest.me:8443",
		"Target registry URL")
	cmd.Flags().StringVar(&flags.Username, "username", "giteaadmin",
		"Registry username (default: giteaadmin)")
	cmd.Flags().StringVar(&flags.Password, "password", "",
		"Registry password (required)")
	cmd.Flags().BoolVar(&flags.Insecure, "insecure", false,
		"Skip TLS certificate verification (use only for self-signed certs)")
	cmd.Flags().BoolVar(&flags.Verbose, "verbose", false,
		"Enable verbose output")

	// Mark password as required
	cmd.MarkFlagRequired("password")

	return cmd
}
```

**Critical**:
- Function returns `*cobra.Command` (pointer)
- Uses `cobra.ExactArgs(1)` - requires exactly one image name argument
- All 5 flags defined with defaults:
  - `--registry` defaults to local Gitea URL
  - `--username` defaults to "giteaadmin"
  - `--password` is REQUIRED (no default)
  - `--insecure` defaults to false (secure mode)
  - `--verbose` defaults to false (quiet mode)
- Password flag is marked as REQUIRED via `cmd.MarkFlagRequired("password")`
- RunE closure calls `runPush()` with context, image name, and flags

---

### Step 5: Implement runPush Stub Function

**Action**: Create the runPush() stub function that prints a message and returns nil.

**Exact Code to Copy** (from WAVE-1.1-ARCHITECTURE.md lines 837-845):

```go
// runPush executes the push workflow.
// This is the main orchestration function that coordinates all packages.
func runPush(ctx context.Context, imageName string, flags *PushFlags) error {
	// Implementation will be provided in Phase 2 Wave 1
	// This skeleton shows the structure and interface usage
	fmt.Fprintf(os.Stderr, "Push command not yet implemented\n")
	fmt.Fprintf(os.Stderr, "Will push %s to %s\n", imageName, flags.Registry)
	return nil
}
```

**Critical**:
- Function is PRIVATE (lowercase `runPush`)
- Takes 3 parameters: `context.Context`, image name string, flags pointer
- Returns `error` (always nil in Wave 1 stub)
- Prints to STDERR (not STDOUT) - this is a diagnostic message
- Does NOT implement any actual push logic (that's Phase 2 Wave 1)

---

### Step 6: Define Exit Code Constants

**Action**: Define the 5 exit code constants.

**Exact Code to Copy** (from WAVE-1.1-ARCHITECTURE.md lines 847-854):

```go
// Exit codes for the push command
const (
	ExitSuccess      = 0 // Push successful
	ExitGeneralError = 1 // Invalid arguments, unexpected failures
	ExitAuthError    = 2 // Authentication failure
	ExitNetworkError = 3 // Registry unreachable, TLS failure
	ExitImageNotFound = 4 // Image not in Docker daemon
)
```

**Critical**:
- All 5 constants are PUBLIC (uppercase first letter)
- Exit codes follow standard Unix conventions:
  - 0 = success
  - 1 = general error
  - 2-4 = specific error categories
- Each constant has an inline comment explaining its purpose
- Constants are NOT USED in Wave 1 stub (will be used in Phase 2 Wave 1)

---

### Step 7: Build Validation

**Action**: Verify the file compiles successfully.

**Commands to Run**:
```bash
cd cmd
go build .
```

**Expected Result**:
- ✅ Build succeeds (no compilation errors)
- ✅ All imports resolve correctly
- ✅ All interface types recognized from previous efforts

**Common Errors to Avoid**:
- ❌ "package docker is not in GOROOT" → Effort 1.1.1 not completed
- ❌ "package registry is not in GOROOT" → Effort 1.1.2 not completed
- ❌ "package auth is not in GOROOT" → Effort 1.1.3 not completed
- ❌ "package tls is not in GOROOT" → Effort 1.1.3 not completed

---

## Test Requirements

### Test Coverage Goals

- **Target Coverage**: 100% (command structure is simple)
- **Expected Pass Rate**: 100% (7/7 tests)
- **Test Approach**: Compilation checks, struct validation, Cobra integration

### Required Tests

Create file `cmd/push_test.go` with the following 7 tests:

#### T1.1.4-001: PushCommand Struct Compiles

**Purpose**: Verify PushCommand struct is a valid Go type with correct fields.

**Test Code**:
```go
package cmd_test

import (
	"testing"

	"github.com/jessesanford/idpbuilder/cmd"
)

func TestPushCommand_StructCompiles(t *testing.T) {
	// Verify struct fields are accessible
	pc := &cmd.PushCommand{}
	_ = pc
}
```

**Expected**: Test passes (struct compiles)

---

#### T1.1.4-002: PushFlags Struct Compiles

**Purpose**: Verify PushFlags struct fields are assignable.

**Test Code**:
```go
func TestPushFlags_StructCompiles(t *testing.T) {
	flags := &cmd.PushFlags{
		Registry: "https://example.com",
		Username: "admin",
		Password: "secret",
		Insecure: true,
		Verbose:  true,
	}

	if flags.Registry != "https://example.com" {
		t.Error("PushFlags struct field assignment failed")
	}
}
```

**Expected**: Test passes (all fields assignable)

---

#### T1.1.4-003: NewPushCommand Creates Valid Cobra Command

**Purpose**: Verify NewPushCommand() returns a properly configured Cobra command.

**Test Code**:
```go
import (
	"github.com/jessesanford/idpbuilder/cmd"
)

func TestNewPushCommand_CreatesValidCommand(t *testing.T) {
	cmd := cmd.NewPushCommand()

	if cmd == nil {
		t.Fatal("NewPushCommand returned nil")
	}

	if cmd.Use != "push IMAGE_NAME" {
		t.Errorf("Expected Use to be 'push IMAGE_NAME', got %q", cmd.Use)
	}

	if cmd.Short == "" {
		t.Error("Command Short description is empty")
	}
}
```

**Expected**: Test passes (command created with correct metadata)

---

#### T1.1.4-004: All Command Flags Defined Correctly

**Purpose**: Verify all 5 flags are defined and accessible.

**Test Code**:
```go
func TestNewPushCommand_FlagsDefined(t *testing.T) {
	cmd := cmd.NewPushCommand()

	// Check required flags
	registryFlag := cmd.Flags().Lookup("registry")
	if registryFlag == nil {
		t.Error("registry flag not defined")
	}

	usernameFlag := cmd.Flags().Lookup("username")
	if usernameFlag == nil {
		t.Error("username flag not defined")
	}

	passwordFlag := cmd.Flags().Lookup("password")
	if passwordFlag == nil {
		t.Error("password flag not defined")
	}

	insecureFlag := cmd.Flags().Lookup("insecure")
	if insecureFlag == nil {
		t.Error("insecure flag not defined")
	}

	verboseFlag := cmd.Flags().Lookup("verbose")
	if verboseFlag == nil {
		t.Error("verbose flag not defined")
	}
}
```

**Expected**: Test passes (all flags found)

---

#### T1.1.4-005: Exit Codes Constants Defined

**Purpose**: Verify all 5 exit code constants have correct values.

**Test Code**:
```go
func TestExitCodes_ConstantsDefined(t *testing.T) {
	if cmd.ExitSuccess != 0 {
		t.Errorf("ExitSuccess should be 0, got %d", cmd.ExitSuccess)
	}

	if cmd.ExitGeneralError != 1 {
		t.Errorf("ExitGeneralError should be 1, got %d", cmd.ExitGeneralError)
	}

	if cmd.ExitAuthError != 2 {
		t.Errorf("ExitAuthError should be 2, got %d", cmd.ExitAuthError)
	}

	if cmd.ExitNetworkError != 3 {
		t.Errorf("ExitNetworkError should be 3, got %d", cmd.ExitNetworkError)
	}

	if cmd.ExitImageNotFound != 4 {
		t.Errorf("ExitImageNotFound should be 4, got %d", cmd.ExitImageNotFound)
	}
}
```

**Expected**: Test passes (all exit codes have correct values)

---

#### T1.1.4-006: runPush Function Signature Valid

**Purpose**: Verify runPush executes via command without errors.

**Test Code**:
```go
import (
	"context"

	"github.com/jessesanford/idpbuilder/cmd"
)

func TestNewPushCommand_RunEFunctionWorks(t *testing.T) {
	cmd := cmd.NewPushCommand()

	// Set required flag
	cmd.Flags().Set("password", "testpass")

	// Execute command with test image name
	cmd.SetArgs([]string{"testimage:latest"})
	ctx := context.Background()
	cmd.SetContext(ctx)

	err := cmd.Execute()
	if err != nil {
		t.Errorf("Command execution failed: %v", err)
	}
}
```

**Expected**: Test passes (command executes, runPush returns nil)

---

#### T1.1.4-007: Command Registers with Cobra Successfully

**Purpose**: Verify command can be registered with a Cobra root command.

**Test Code**:
```go
import (
	"github.com/jessesanford/idpbuilder/cmd"
	"github.com/spf13/cobra"
)

func TestNewPushCommand_RegistersWithCobra(t *testing.T) {
	rootCmd := &cobra.Command{Use: "idpbuilder"}
	pushCmd := cmd.NewPushCommand()

	rootCmd.AddCommand(pushCmd)

	// Verify command was added
	found := false
	for _, c := range rootCmd.Commands() {
		if c.Use == "push IMAGE_NAME" {
			found = true
			break
		}
	}

	if !found {
		t.Error("Push command not registered with root command")
	}
}
```

**Expected**: Test passes (command found in root command's children)

---

### Test Execution Plan

**During Implementation**:
```bash
# After creating cmd/push.go and cmd/push_test.go
cd cmd
go test -v -cover

# Expected output:
# PASS: TestPushCommand_StructCompiles
# PASS: TestPushFlags_StructCompiles
# PASS: TestNewPushCommand_CreatesValidCommand
# PASS: TestNewPushCommand_FlagsDefined
# PASS: TestExitCodes_ConstantsDefined
# PASS: TestNewPushCommand_RunEFunctionWorks
# PASS: TestNewPushCommand_RegistersWithCobra
# coverage: 100.0% of statements
# PASS
# ok      github.com/jessesanford/idpbuilder/cmd  0.XXXs
```

---

## Size Management

### Estimated Line Count

- **cmd/push.go**: 190 lines (estimated from Wave Plan)
  - Package declaration + imports: ~15 lines
  - PushCommand struct: ~15 lines
  - PushFlags struct: ~17 lines
  - NewPushCommand function: ~51 lines
  - runPush function: ~9 lines
  - Exit code constants: ~8 lines
  - Documentation comments: ~75 lines

- **cmd/push_test.go**: ~140 lines (7 tests)

**Total**: 330 lines (production + tests)

### Measurement Tool

**ALWAYS use the line counter tool**:
```bash
# Find project root first
PROJECT_ROOT=$(pwd)
while [ "$PROJECT_ROOT" != "/" ]; do
    [ -f "$PROJECT_ROOT/orchestrator-state-v3.json" ] && break
    PROJECT_ROOT=$(dirname "$PROJECT_ROOT")
done

# Run line counter (auto-detects base branch)
$PROJECT_ROOT/tools/line-counter.sh
```

**Tool will output**:
```
🎯 Detected base: idpbuilder-oci-push-rebuild/phase1-wave1-effort-1.1.3
📦 Analyzing branch: idpbuilder-oci-push-rebuild/phase1-wave1-effort-1.1.4
✅ Total implementation lines: ~190
⚠️  Note: Tests, demos, docs, configs NOT included
```

### Size Compliance

- **Target**: <700 lines (soft limit)
- **Hard Limit**: 800 lines (R535 enforcement threshold)
- **Estimated**: 190 lines
- **Status**: ✅ COMPLIANT (well under both limits)
- **Split Required**: ❌ NO (190 << 800)

### Check Frequency

**Measure size**:
- After completing cmd/push.go
- Before starting tests
- Before final commit

**Expected Results**:
- After cmd/push.go: ~190 lines
- Final (with tests): ~190 lines (tests excluded from count per R007)

---

## Pattern Compliance

### IDPBuilder Patterns

**Command Structure Pattern**:
- ✅ Use Cobra for CLI framework
- ✅ Place commands in `cmd/` package
- ✅ Use `*cobra.Command` return type
- ✅ Implement `RunE` for error handling
- ✅ Define flags with default values
- ✅ Mark required flags explicitly

**Interface Usage Pattern**:
- ✅ Store interfaces as struct fields (not concrete types)
- ✅ Use private fields for command internals
- ✅ Import interface packages (docker, registry, auth, tls)
- ✅ Do NOT instantiate implementations in Wave 1 (just define structure)

**Error Handling Pattern** (not implemented in Wave 1):
- Exit codes defined but NOT USED yet
- Error handling will be implemented in Phase 2 Wave 1

### Go Best Practices

- ✅ All public types have documentation comments
- ✅ Struct fields have inline comments
- ✅ Function signatures are clear and simple
- ✅ No global variables
- ✅ No init() functions
- ✅ Proper package-level comment

---

## Security Requirements

### Wave 1 Security Considerations

**Password Handling**:
- ✅ Password flag is marked as REQUIRED (cannot be empty)
- ✅ NO default password (security best practice)
- ⚠️ Password passed via command-line flag (visible in process list)
  - **NOTE**: This is acceptable for Wave 1 scaffolding
  - **TODO**: Phase 2 should add environment variable support (IDPBUILDER_REGISTRY_PASSWORD)
  - **TODO**: Phase 2 should warn users about command-line password visibility

**TLS Security**:
- ✅ Default to SECURE mode (insecure flag defaults to false)
- ✅ Insecure mode requires explicit opt-in (--insecure flag)
- ⚠️ No TLS validation in Wave 1 (just flag definition)
  - **NOTE**: TLS validation will be implemented in Phase 2 Wave 1

**Input Validation** (not in Wave 1):
- Image name validation: NOT implemented (Phase 2 Wave 1)
- Registry URL validation: NOT implemented (Phase 2 Wave 1)
- Username validation: NOT implemented (Phase 2 Wave 1)

### No Security Vulnerabilities in Wave 1

**Why Wave 1 is Safe**:
- No actual authentication logic (just interface declarations)
- No actual network connections (just flag definitions)
- No file I/O operations
- No subprocess execution
- No credential storage

---

## Integration Points

### How This Effort Connects to Others

**Imports from Previous Efforts**:
- `pkg/docker.DockerClient` (from Effort 1.1.1)
- `pkg/registry.RegistryClient` (from Effort 1.1.2)
- `pkg/auth.AuthProvider` (from Effort 1.1.3)
- `pkg/tls.TLSProvider` (from Effort 1.1.3)

**Critical**: All 4 interface packages MUST exist and compile before implementing this effort.

**Used by Future Efforts** (Phase 2 Wave 1):
- Phase 2 Wave 1 Effort 2.1.1: Docker Client Implementation
  - Will instantiate DockerClient and inject into PushCommand
- Phase 2 Wave 1 Effort 2.1.2: Registry Client Implementation
  - Will instantiate RegistryClient and inject into PushCommand
- Phase 2 Wave 1 Effort 2.1.3: Auth Provider Implementation
  - Will instantiate AuthProvider and inject into PushCommand
- Phase 2 Wave 1 Effort 2.1.4: TLS Provider Implementation
  - Will instantiate TLSProvider and inject into PushCommand
- Phase 2 Wave 1 Effort 2.1.5: Push Workflow Orchestration
  - Will implement runPush() function body
  - Will use PushCommand struct to coordinate all packages

### Integration with IDPBuilder CLI

**Current IDPBuilder Root Command**:
The IDPBuilder CLI uses `pkg/cmd/root.go` as the root command. The push command will be registered there.

**Registration Example** (NOT part of this effort - for reference only):
```go
// In pkg/cmd/root.go (existing file)
import (
	"github.com/jessesanford/idpbuilder/cmd"  // Import push command package
	"github.com/spf13/cobra"
)

func Execute() {
	rootCmd := &cobra.Command{
		Use:   "idpbuilder",
		Short: "IDPBuilder CLI tool",
	}

	// Register push command
	rootCmd.AddCommand(cmd.NewPushCommand())  // NEW - will be added in Phase 2

	// Execute
	rootCmd.Execute()
}
```

**NOTE**: The actual registration will happen in Phase 2 Wave 1. This effort ONLY creates the command structure.

---

## Acceptance Criteria

### Must-Have Requirements

Before marking this effort as COMPLETE, verify ALL of the following:

**File Creation**:
- [ ] File `cmd/push.go` created at correct location
- [ ] File `cmd/push_test.go` created with 7 tests

**Code Structure**:
- [ ] PushCommand struct defined with 4 interface fields
- [ ] PushFlags struct defined with 5 fields
- [ ] NewPushCommand() function creates configured Cobra command
- [ ] 5 command flags defined (registry, username, password, insecure, verbose)
- [ ] Password flag marked as required via `MarkFlagRequired("password")`
- [ ] runPush() function stub created (prints message, returns nil)
- [ ] 5 exit code constants defined with correct values

**Build Validation**:
- [ ] `go build ./cmd` succeeds (no compilation errors)
- [ ] All imports resolve correctly
- [ ] No linting errors (`golangci-lint run ./cmd`)

**Test Validation**:
- [ ] `go test ./cmd` succeeds (all 7 tests pass)
- [ ] Test coverage ≥95% (target 100%)
- [ ] All tests pass on first run (no flaky tests)

**Documentation**:
- [ ] All public types have documentation comments
- [ ] PushCommand struct has descriptive comment
- [ ] PushFlags struct has descriptive comment
- [ ] NewPushCommand() function has documentation with example
- [ ] runPush() function has documentation explaining stub nature
- [ ] All exit code constants have inline comments

**Size Compliance**:
- [ ] Line count measured with `tools/line-counter.sh`
- [ ] Implementation lines ≤800 (hard limit per R535)
- [ ] Line count within estimate (±15%: 162-219 lines acceptable)

**Integration Readiness**:
- [ ] Command can be registered with Cobra root command
- [ ] Command accepts image name as single argument
- [ ] Command prints expected "not yet implemented" message when executed
- [ ] No runtime panics or crashes

---

## Demo Feasibility (R630)

### Can QA Demonstrate This Feature?

**YES** - This effort CAN be demonstrated even though it's scaffolding only.

### Demo Plan

**Demo Scenario**: Register and execute the push command (stub behavior)

**Prerequisites**:
- Go development environment
- IDPBuilder source code
- All Wave 1 interface packages (Efforts 1.1.1-1.1.3) completed

**Demo Steps**:

1. **Build the command**:
   ```bash
   cd cmd
   go build .
   # Expected: Build succeeds
   ```

2. **Create test program** (temporary main.go for demo):
   ```go
   package main

   import (
       "fmt"
       "os"

       "github.com/jessesanford/idpbuilder/cmd"
       "github.com/spf13/cobra"
   )

   func main() {
       rootCmd := &cobra.Command{Use: "idpbuilder"}
       rootCmd.AddCommand(cmd.NewPushCommand())

       if err := rootCmd.Execute(); err != nil {
           fmt.Fprintf(os.Stderr, "Error: %v\n", err)
           os.Exit(1)
       }
   }
   ```

3. **Build and run demo**:
   ```bash
   go build -o idpbuilder-demo main.go
   ./idpbuilder-demo push --help
   # Expected: Shows push command help with all 5 flags
   ```

4. **Execute push command (stub)**:
   ```bash
   ./idpbuilder-demo push myapp:latest --password test123
   # Expected output to STDERR:
   # Push command not yet implemented
   # Will push myapp:latest to https://gitea.cnoe.localtest.me:8443
   ```

5. **Verify flag validation**:
   ```bash
   ./idpbuilder-demo push myapp:latest
   # Expected: Error - password flag required
   ```

6. **Verify custom flags work**:
   ```bash
   ./idpbuilder-demo push myapp:v1 --registry https://custom.com --username admin --password secret --insecure --verbose
   # Expected output to STDERR:
   # Push command not yet implemented
   # Will push myapp:v1 to https://custom.com
   ```

### Demo Success Criteria

**What QA Must Validate**:
- [ ] Command help displays correctly (push --help)
- [ ] All 5 flags are shown in help output
- [ ] Command requires exactly 1 argument (image name)
- [ ] Password flag is marked as required
- [ ] Stub message prints to STDERR (not STDOUT)
- [ ] Stub message shows correct image name and registry URL
- [ ] Exit code is 0 (ExitSuccess) when stub executes
- [ ] Custom flag values are accepted without errors

**Demo Artifacts**:
- Screenshot of help output
- Screenshot of stub execution output
- Test program source code (main.go)

### Feature Flag Analysis

**Feature Flags Required**: ❌ NO

**Rationale**:
- This is scaffolding only (no user-facing functionality)
- Stub always prints "not yet implemented" message
- No actual push workflow exists to disable
- Command registration happens in Phase 2 (not this effort)

---

## Stub Detection (R629)

### Wave 1 Stub Policy

**Stubs ALLOWED in Wave 1**: ✅ YES

**Why**:
- Wave 1 is INTERFACE DEFINITIONS ONLY (by design)
- All constructor functions MUST panic("not implemented") per architecture
- runPush() stub is INTENTIONAL (prints message, returns nil)
- This is NOT production code - it's scaffolding for Phase 2

### Stub Inventory

**Intentional Stubs in This Effort**:
1. **runPush() function** (line ~838 in cmd/push.go):
   - Status: STUB (prints "not yet implemented", returns nil)
   - Reason: Push workflow implementation is Phase 2 Wave 1
   - Will be replaced: YES (in Phase 2 Wave 1 Effort 2.1.5)
   - Blocking for Wave 1: ❌ NO (stub is expected)

### Stub Detection Command

**Run stub detection**:
```bash
cd cmd
bash $CLAUDE_PROJECT_DIR/tools/detect-stubs.sh
```

**Expected Output**:
```
⚠️ STUBS FOUND (allowed in effort directories):
cmd/push.go:838: runPush() - stub prints "not yet implemented"

Status: STUBS ALLOWED IN EFFORT DIRECTORIES
Wave 1 is interface definitions - stubs are intentional
Phase 2 Wave 1 will implement actual functionality
```

**Action**: ✅ APPROVE with stub tracking

**Tracking**:
- Document stub location in this plan
- Create tracking for Phase 2 Wave 1 to implement runPush()
- DO NOT block Wave 1 completion due to intentional stubs

---

## R629 Compliance: Stub Detection Analysis

### Policy: Effort-Level Stubs Allowed (Wave 1)

**Context**: This is an EFFORT in Phase 1 Wave 1 (interface definitions).

**R629 Policy for Efforts**:
- ✅ Stubs ALLOWED (work in progress)
- ✅ Track stubs for completion in later waves/phases
- ❌ Stubs NOT allowed in wave integration or phase integration

### Stub Analysis

**Intentional Stubs**:
1. runPush() function - returns nil, prints stub message
   - Purpose: Scaffolding for Phase 2 implementation
   - Tracking: BUG-1.1.4-STUB-001 (to be created for Phase 2)

**Review Decision**: ✅ APPROVE with stub tracking

**Stub Completion Plan**:
- Phase 2 Wave 1 Effort 2.1.5 will implement runPush() body
- Exit codes will be used in Phase 2 error handling
- Interface fields will be populated with implementations in Phase 2

---

## Known Issues and Limitations

### Wave 1 Limitations (By Design)

**What This Effort Does NOT Provide**:
1. **No Actual Push Functionality**: runPush() is a stub (returns nil)
2. **No Error Handling**: Exit codes defined but not used
3. **No Progress Reporting**: Verbose flag defined but not implemented
4. **No Validation**: No input validation for image names, URLs, credentials
5. **No Interface Implementations**: PushCommand has interface fields but no implementations

**These are NOT bugs** - they are intentional design for Wave 1 (interface scaffolding only).

### Potential Issues

**Build Issues**:
- **Issue**: Build fails with "package docker/registry/auth/tls not found"
- **Cause**: Efforts 1.1.1-1.1.3 not completed yet
- **Solution**: Ensure all upstream dependencies are merged first

**Test Issues**:
- **Issue**: TestNewPushCommand_RunEFunctionWorks fails with panic
- **Cause**: runPush() stub changed to panic instead of returning nil
- **Solution**: runPush() MUST return nil (not panic)

**Integration Issues**:
- **Issue**: Command not found when registered with root command
- **Cause**: Wrong import path or package name
- **Solution**: Verify package is `cmd` and import path is `github.com/jessesanford/idpbuilder/cmd`

---

## Timeline Estimate

### Implementation Time

**Estimated Total**: 2-3 hours for experienced Go developer

**Breakdown**:
- Step 1 (File structure): 15 minutes
- Step 2 (PushCommand struct): 15 minutes
- Step 3 (PushFlags struct): 15 minutes
- Step 4 (NewPushCommand function): 30 minutes
- Step 5 (runPush stub): 10 minutes
- Step 6 (Exit codes): 10 minutes
- Step 7 (Build validation): 10 minutes
- Test creation (7 tests): 45 minutes
- Test debugging: 15 minutes
- Documentation review: 15 minutes

**Fast Track**: 1.5 hours (if copying code exactly from architecture)

**Conservative**: 4 hours (if unfamiliar with Cobra or Go)

---

## Next Steps After Completion

### Immediate Next Steps

1. **Code Review**: Code Reviewer validates implementation
2. **Merge to Integration Branch**: `idpbuilder-oci-push-rebuild/phase1-wave1-integration`
3. **Wave 1 Completion**: All 4 efforts (1.1.1-1.1.4) now complete
4. **Wave Integration Tests**: Run full Wave 1 test suite
5. **Architect Review**: Validate interface contracts before Phase 2

### Phase 2 Planning

**Wave 1 Complete → Architect Creates Phase 2 Plan**:
- Phase 2 Wave 1: Implement all 4 interfaces (Docker, Registry, Auth, TLS)
- Phase 2 Wave 2: Implement runPush() workflow orchestration
- Phase 2 Wave 3: Add error handling and progress reporting
- Phase 2 Wave 4: Integration testing and end-to-end validation

---

## References

### Wave Plan
- **Source**: `planning/phase1/wave1/WAVE-1-IMPLEMENTATION-PLAN.md`
- **Effort 1.1.4 Section**: Lines 1150-1537
- **R213 Metadata**: Lines 1154-1177

### Architecture Document
- **Source**: `planning/phase1/wave1/WAVE-1.1-ARCHITECTURE.md`
- **cmd/push.go Code**: Lines 735-855
- **Integration Examples**: Lines 857-1050

### Test Plan
- **Source**: `planning/phase1/wave1/WAVE-1-TEST-PLAN.md`
- **Effort 1.1.4 Tests**: T1.1.4-001 through T1.1.4-007
- **Total Wave 1 Tests**: 36 tests (6 + 9 + 8 + 7 + 6 integration)

### Rule References
- **R213**: Effort Metadata Requirements (parallelization, dependencies)
- **R502**: Implementation Plan Quality Gates (EXACT fidelity)
- **R535**: Size Limits (800 line hard limit for Code Reviewers)
- **R629**: Stub Detection and Production Code Enforcement
- **R630**: Demo Validation and Success Criteria

---

**Document Status**: ✅ READY FOR IMPLEMENTATION
**Created By**: Code Reviewer Agent
**Validated Against**: WAVE-1-IMPLEMENTATION-PLAN.md, WAVE-1.1-ARCHITECTURE.md
**Fidelity Level**: EXACT (real code from architecture document)
**R213 Compliance**: ✅ COMPLETE (all metadata present)
**R502 Compliance**: ✅ COMPLETE (exact specifications provided)
**R629 Compliance**: ✅ COMPLETE (stub detection analysis included)
**R630 Compliance**: ✅ COMPLETE (demo plan validated)
