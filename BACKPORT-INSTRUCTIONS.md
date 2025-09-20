# R321 BACKPORT INSTRUCTIONS - Test Fixes

## Overview
These fixes were discovered during project integration testing and MUST be backported to source branches per R321 (Supreme Law).
All test failures found in integration MUST be fixed in the original effort branches to maintain trunk-based development.

## Context
- Integration Branch: `idpbuilder-oci-build-push/project-integration`
- Source Effort Branch: `idpbuilder-oci-build-push/phase2/wave1/image-builder`
- CLI Commands Branch: `idpbuilder-oci-build-push/phase2/wave2/cli-commands`

## Fixes for image-builder effort

### FIX-TEST-001: Add EnableImageBuilderFlag constant
**STATUS: Already exists - NO ACTION NEEDED**
**File:** `pkg/build/feature_flags.go`
**Current State:** The constant already exists at line 7:
```go
const (
	// EnableImageBuilderFlag controls whether the OCI image builder is enabled
	EnableImageBuilderFlag = "IDPBUILDER_ENABLE_IMAGE_BUILDER"
)
```

### FIX-TEST-002: Add ErrFeatureDisabled error
**STATUS: Already exists - NO ACTION NEEDED**
**File:** `pkg/build/image_builder.go`
**Current State:** The error already exists at lines 18-20:
```go
var (
	// ErrFeatureDisabled is returned when the image builder feature is disabled
	ErrFeatureDisabled = errors.New("image builder feature is disabled")
)
```

### FIX-TEST-003: Feature flag check
**STATUS: Already exists - NO ACTION NEEDED**
**File:** `pkg/build/image_builder.go`
**Current State:** The feature flag check already exists in the BuildImage method at lines 41-44:
```go
// Check if the image builder feature is enabled
if !IsImageBuilderEnabled() {
    return nil, ErrFeatureDisabled
}
```

### FIX-TEST-005: Fix ErrorIs assertion
**File:** `pkg/build/image_builder_test.go`
**Location:** Line 39 in TestBuildImageFeatureDisabled function
**REQUIRED CHANGE:**

**BEFORE (current incorrect code):**
```go
func TestBuildImageFeatureDisabled(t *testing.T) {
	os.Unsetenv(EnableImageBuilderFlag)
	tempDir := t.TempDir()
	builder, _ := NewBuilder(tempDir)
	contextDir := createTestContext(t)

	result, err := builder.BuildImage(context.Background(), BuildOptions{
		ContextPath: contextDir,
		Tag:         "test:latest",
	})

	assert.Error(t, err)
	assert.Equal(t, ErrFeatureDisabled, err)  // LINE 39 - THIS IS WRONG
	assert.Nil(t, result)
}
```

**AFTER (corrected code):**
```go
func TestBuildImageFeatureDisabled(t *testing.T) {
	os.Unsetenv(EnableImageBuilderFlag)
	tempDir := t.TempDir()
	builder, _ := NewBuilder(tempDir)
	contextDir := createTestContext(t)

	result, err := builder.BuildImage(context.Background(), BuildOptions{
		ContextPath: contextDir,
		Tag:         "test:latest",
	})

	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrFeatureDisabled)  // LINE 39 - CORRECTED
	assert.Nil(t, result)
}
```

**Note:** You may also need to add the `errors` import at the top of the test file if using `errors.Is()` directly instead of `assert.ErrorIs()`.

## Fixes for cli-commands effort

### FIX-TEST-004: Remove unused import
**File:** `pkg/cmd_test/build_test.go`
**Location:** Line 6 in the imports section
**REQUIRED CHANGE:**

**BEFORE (current code with unused import):**
```go
package cmd_test

import (
	"testing"

	"github.com/cnoe-io/idpbuilder/pkg/cmd"  // LINE 6 - UNUSED IMPORT
	"github.com/spf13/cobra"
)
```

**AFTER (corrected code):**
```go
package cmd_test

import (
	"testing"

	"github.com/spf13/cobra"
)
```

**Explanation:** The `"github.com/cnoe-io/idpbuilder/pkg/cmd"` import is not used anywhere in the test. The variable `cmd` used throughout the test is a local variable of type `*cobra.Command`, not from the imported package.

## Verification Steps

### For image-builder effort:
1. Navigate to the effort directory:
   ```bash
   cd efforts/phase2/wave1/image-builder
   ```

2. Apply the fix to `pkg/build/image_builder_test.go` (only FIX-TEST-005 needed)

3. Run tests to verify:
   ```bash
   go test ./pkg/build/...
   ```

4. Expected result: All tests pass, specifically `TestBuildImageFeatureDisabled` should pass

### For cli-commands effort:
1. Navigate to the effort directory:
   ```bash
   cd efforts/phase2/wave2/cli-commands
   ```

2. Apply the fix to `pkg/cmd_test/build_test.go` (remove unused import)

3. Run tests to verify:
   ```bash
   go test ./pkg/cmd_test/...
   ```

4. Check for unused imports:
   ```bash
   go vet ./...
   ```

5. Expected result: No unused import warnings, all tests pass

## Summary

Out of the 5 fixes identified:
- **3 fixes (FIX-TEST-001, 002, 003)** are already implemented in the image-builder effort
- **1 fix (FIX-TEST-005)** needs to be applied to image-builder effort - change error assertion method
- **1 fix (FIX-TEST-004)** needs to be applied to cli-commands effort - remove unused import

These are minor fixes that ensure test consistency and clean code. The main functional code is already correct.