# Implementation Plan: Command Structure & Flag Definitions
## Effort 1.1.4 - Phase 1, Wave 1

**Created**: 2025-10-29T04:08:15Z
**Planner**: @agent-code-reviewer
**Effort ID**: 1.1.4
**Phase**: 1 (Foundation & Interfaces)
**Wave**: 1 (Interface & Contract Definitions)

---

## 🚨 EFFORT INFRASTRUCTURE METADATA (R360)

**EFFORT_NAME**: effort-4-command-structure
**BRANCH**: idpbuilder-oci-push/phase1/wave1/effort-4-command-structure
**BASE_BRANCH**: idpbuilder-oci-push/phase1/wave1/integration
**PHASE**: 1
**WAVE**: 1
**EFFORT_INDEX**: 4
**PARALLELIZATION**: sequential
**CAN_PARALLELIZE**: No
**PARALLEL_WITH**: None
**DEPENDENCIES**: ["1.1.1", "1.1.2", "1.1.3"]

---

## Overview

**Purpose**: Define the `push` command structure with all CLI flags, help text, and execution skeleton (no actual implementation).

**Scope**: Command structure ONLY - no implementations (Wave 1 contract definition)

**Estimated Size**: 130 lines (implementation code only, excluding tests)

**Expected Outcomes**:
- Complete Cobra command definition for `push`
- 5 CLI flags (--registry, --username, --password, --insecure, --verbose)
- Default constants
- Comprehensive help text with examples
- Placeholder execution function
- Command structure validation tests

---

## 🔴🔴🔴 REPOSITORY CONTEXT (R251/R309) 🔴🔴🔴

**CRITICAL UNDERSTANDING**:
- ✅ This plan is for the idpbuilder TARGET repo (https://github.com/jessesanford/idpbuilder.git)
- ✅ Implementation will happen in TARGET repo clone
- ✅ NOT in Software Factory planning repo
- ✅ Files reference TARGET repo structure: `pkg/`, `cmd/`, etc.

**Working Directory**: `/efforts/phase1/wave1/effort-4-command-structure/`
**Target Branch**: `idpbuilder-oci-push/phase1/wave1/effort-4-command-structure`
**Base Branch**: `idpbuilder-oci-push/phase1/wave1/integration`

---

## 🔴🔴🔴 EXPLICIT SCOPE CONTROL (R311 - SUPREME LAW) 🔴🔴🔴

### IMPLEMENT EXACTLY:

**Command File (1 file, ~130 lines)**:
- `cmd/push.go` with:
  - 2 constants (DefaultRegistry, DefaultUsername)
  - 5 flag variables
  - Cobra command definition with full help text
  - `init()` function with flag definitions
  - `runPushCommand()` stub function
  - Helper function signatures (stubs)

**Test File (1 file, ~70 lines - NOT counted)**:
- `cmd/push_test.go` with:
  - Command structure validation
  - Flag definition tests
  - Constants verification

**TOTAL IMPLEMENTATION**: ~130 lines (excludes tests per R007)

### DO NOT IMPLEMENT:

❌ NO actual push logic (Phase 2)
❌ NO Docker client initialization (Phase 2)
❌ NO registry client initialization (Phase 2)
❌ NO error handling logic (Phase 2)
❌ NO progress display (Phase 2)
❌ NO environment variable binding (Phase 2)
❌ NO additional helper functions beyond signatures
❌ NO logging implementation
❌ NO validation implementation
❌ NO additional commands
❌ NO root command registration (will be done in Phase 2)

---

## File Structure

### Files to Create

**Implementation File**:
1. **`cmd/push.go`** (~130 lines)
   - Package declaration and imports
   - Constants (2)
   - Flag variables (5)
   - Cobra command definition
   - Help text with examples
   - Flag definitions in init()
   - Stub functions

**Test File** (excluded from 800-line limit):
2. **`cmd/push_test.go`** (~70 lines)
   - Command structure tests
   - Flag validation tests
   - Constant verification tests

---

## Implementation Steps

### Step 1: Create cmd Directory (if not exists)

```bash
cd /efforts/phase1/wave1/effort-4-command-structure
mkdir -p cmd
```

### Step 2: Create cmd/push.go

**Complete file content** (from Wave Implementation Plan lines 966-1085):

```go
// Package cmd implements the IDPBuilder CLI commands.
package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

const (
	// DefaultRegistry is the default Gitea registry URL
	DefaultRegistry = "https://gitea.cnoe.localtest.me:8443"

	// DefaultUsername is the default registry username
	DefaultUsername = "giteaadmin"
)

var (
	// Flag variables
	registryURL string
	username    string
	password    string
	insecure    bool
	verbose     bool
)

// pushCmd represents the push command
var pushCmd = &cobra.Command{
	Use:   "push IMAGE",
	Short: "Push Docker image to OCI registry",
	Long: `Push a Docker image from local daemon to an OCI-compatible container registry.

The push command reads an image from the local Docker daemon and uploads it
to the specified registry (default: Gitea). It supports authentication with
username and password, and can bypass TLS certificate verification with the
--insecure flag for development/testing environments.

Examples:
  # Push to default Gitea registry
  idpbuilder push myapp:latest --password 'mypassword'

  # Push with custom username
  idpbuilder push myapp:latest --username developer --password 'myP@ss'

  # Push with insecure mode (bypass TLS verification)
  idpbuilder push myapp:latest -k --password 'mypassword'

  # Push to custom registry
  idpbuilder push myapp:v1.0 --registry https://registry.io --password 'pass'

  # Verbose mode for debugging
  idpbuilder push myapp:latest --verbose --password 'pass'

Environment Variables:
  IDPBUILDER_REGISTRY           Override default registry URL
  IDPBUILDER_REGISTRY_USERNAME  Override default username
  IDPBUILDER_REGISTRY_PASSWORD  Provide password (alternative to --password flag)
  IDPBUILDER_INSECURE           Set to "true" to enable insecure mode

Flag priority: CLI flags > Environment variables > Defaults`,
	Args: cobra.ExactArgs(1), // Require exactly one argument: IMAGE
	RunE: runPushCommand,
}

func init() {
	// Define flags
	pushCmd.Flags().StringVar(&registryURL, "registry", DefaultRegistry,
		"Registry URL to push to")
	pushCmd.Flags().StringVarP(&username, "username", "u", DefaultUsername,
		"Registry username for authentication")
	pushCmd.Flags().StringVarP(&password, "password", "p", "",
		"Registry password for authentication (REQUIRED)")
	pushCmd.Flags().BoolVarP(&insecure, "insecure", "k", false,
		"Skip TLS certificate verification (insecure mode)")
	pushCmd.Flags().BoolVarP(&verbose, "verbose", "v", false,
		"Enable verbose output for debugging")

	// Mark password as required
	pushCmd.MarkFlagRequired("password")

	// TODO: Add environment variable binding in Wave 2
	// viper.BindPFlag("registry", pushCmd.Flags().Lookup("registry"))
	// viper.BindEnv("registry", "IDPBUILDER_REGISTRY")

	// Register command with root command
	// rootCmd.AddCommand(pushCmd)  // Will be uncommented in Phase 2
}

// runPushCommand is the main execution function for the push command.
//
// This function orchestrates the complete push workflow:
//   1. Validate inputs (flags, image name)
//   2. Initialize Docker client
//   3. Check image exists in Docker daemon
//   4. Retrieve image from Docker
//   5. Initialize registry client with auth and TLS
//   6. Build target registry reference
//   7. Push image to registry with progress reporting
//   8. Report success or failure
//
// Implementation will be completed in Phase 2 Wave 1.
func runPushCommand(cmd *cobra.Command, args []string) error {
	// Phase 2 implementation placeholder
	imageName := args[0]

	if verbose {
		fmt.Printf("Push command invoked (not yet implemented)\n")
		fmt.Printf("  Image: %s\n", imageName)
		fmt.Printf("  Registry: %s\n", registryURL)
		fmt.Printf("  Username: %s\n", username)
		fmt.Printf("  Insecure: %v\n", insecure)
	}

	return fmt.Errorf("push command not yet implemented - interface definition only (Wave 1)")
}
```

### Step 3: Create cmd/push_test.go

**Complete file content** (from Wave Implementation Plan lines 1087-1154):

```go
package cmd

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

// TestPushCommandStructure verifies command structure is valid
func TestPushCommandStructure(t *testing.T) {
	assert.Equal(t, "push IMAGE", pushCmd.Use)
	assert.NotEmpty(t, pushCmd.Short)
	assert.NotEmpty(t, pushCmd.Long)
	assert.NotNil(t, pushCmd.RunE)
}

// TestPushCommandFlags verifies all flags are defined
func TestPushCommandFlags(t *testing.T) {
	assert.NotNil(t, pushCmd.Flags().Lookup("registry"))
	assert.NotNil(t, pushCmd.Flags().Lookup("username"))
	assert.NotNil(t, pushCmd.Flags().Lookup("password"))
	assert.NotNil(t, pushCmd.Flags().Lookup("insecure"))
	assert.NotNil(t, pushCmd.Flags().Lookup("verbose"))

	// Verify short flags
	assert.NotNil(t, pushCmd.Flags().ShorthandLookup("u"))
	assert.NotNil(t, pushCmd.Flags().ShorthandLookup("p"))
	assert.NotNil(t, pushCmd.Flags().ShorthandLookup("k"))
	assert.NotNil(t, pushCmd.Flags().ShorthandLookup("v"))
}

// TestPushCommandConstants verifies constants are defined
func TestPushCommandConstants(t *testing.T) {
	assert.Equal(t, "https://gitea.cnoe.localtest.me:8443", DefaultRegistry)
	assert.Equal(t, "giteaadmin", DefaultUsername)
}

// TestPushCommandExecution verifies command returns not implemented error
func TestPushCommandExecution(t *testing.T) {
	// Reset flags to defaults
	registryURL = DefaultRegistry
	username = DefaultUsername
	password = "test"
	insecure = false
	verbose = false

	err := runPushCommand(pushCmd, []string{"myapp:latest"})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not yet implemented")
	assert.Contains(t, err.Error(), "Wave 1")
}

// TestPushCommandVerboseMode verifies verbose flag is respected
func TestPushCommandVerboseMode(t *testing.T) {
	verbose = true
	defer func() { verbose = false }()

	err := runPushCommand(pushCmd, []string{"testimage:v1"})

	assert.Error(t, err)
	// In verbose mode, output is printed (tested via manual inspection)
}
```

### Step 4: Run Tests

```bash
cd /efforts/phase1/wave1/effort-4-command-structure
go test ./cmd/push_test.go -v
```

### Step 5: Measure Size

```bash
PROJECT_ROOT=$(pwd)
while [ "$PROJECT_ROOT" != "/" ]; do
    if [ -f "$PROJECT_ROOT/orchestrator-state-v3.json" ]; then break; fi
    PROJECT_ROOT=$(dirname "$PROJECT_ROOT")
done
cd /efforts/phase1/wave1/effort-4-command-structure
$PROJECT_ROOT/tools/line-counter.sh
```

### Step 6: Commit and Push

```bash
cd /efforts/phase1/wave1/effort-4-command-structure
git add cmd/
git commit -m "feat(cmd): add push command structure and flag definitions

Implements Effort 1.1.4 - Command Structure & Flag Definitions
Phase 1, Wave 1: Interface & Contract Definitions

Added:
- Cobra command definition for 'push' subcommand
- 5 CLI flags: --registry, --username, --password, --insecure, --verbose
- Default constants (DefaultRegistry, DefaultUsername)
- Comprehensive help text with examples
- Placeholder execution function (returns not implemented error)
- Command structure validation tests

Implementation lines: ~130
Test coverage: 100%
All tests passing

Part of Wave 1 contract definition for Phase 2 Wave 1 implementation.

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com>"

git push origin idpbuilder-oci-push/phase1/wave1/effort-4-command-structure
```

---

## Dependencies

### Upstream Dependencies

- **Effort 1.1.1**: Docker interfaces (for future imports)
- **Effort 1.1.2**: Registry interfaces (for future imports)
- **Effort 1.1.3**: Auth/TLS interfaces (for future imports)

**Note**: Wave 1 defines structures only; actual imports will be added in Phase 2

### Downstream Dependencies

**None** - This is the last effort in Wave 1

### External Library Dependencies

```go
require (
	github.com/spf13/cobra v1.8.0   // Already in idpbuilder
	github.com/spf13/viper v1.18.0  // Already in idpbuilder
	github.com/stretchr/testify v1.9.0
)
```

---

## Acceptance Criteria

- [ ] Both files created (push.go + push_test.go)
- [ ] Command compiles with Cobra
- [ ] All 5 flags defined correctly
- [ ] Help text complete with examples
- [ ] Constants defined
- [ ] All tests passing (100% pass rate)
- [ ] Test coverage = 100%
- [ ] runPushCommand returns "not implemented" error
- [ ] Line count: 130±20 lines

---

## Document Status

**Status**: ✅ READY FOR IMPLEMENTATION
**Created**: 2025-10-29T04:08:15Z
**Planner**: @agent-code-reviewer
**Effort**: 1.1.4 - Command Structure & Flag Definitions
**Phase**: 1, Wave: 1
**Fidelity**: EXACT (complete code provided)

**Lines**: 130 implementation + 70 test
**Coverage**: 100% required
**Dependencies**: ["1.1.1", "1.1.2", "1.1.3"]

---

**END OF IMPLEMENTATION PLAN - EFFORT 1.1.4**
