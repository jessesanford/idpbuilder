# Phase 2 Wave 2 Integration Fix Plan

## Error Summary

The integration build is failing with the following errors in `pkg/cmd/push.go`:

1. **Unused import**: `github.com/cnoe-io/idpbuilder/pkg/cmd/helpers` imported but not used
2. **Undefined function**: `certs.ExtractCertificate` doesn't exist
3. **Wrong parameters**: `gitea.NewClient` expects 2 parameters, got 4
4. **Undefined function**: `gitea.NewClientAutoDetect` doesn't exist
5. **Undefined method**: `ValidateCredentials` not defined on `gitea.Client`
6. **Undefined method**: `PushImage` not defined on `gitea.Client`

## Root Cause Analysis

The three efforts developed independently without proper interface coordination:

1. **cli-commands** effort:
   - Expects functions that don't exist in other efforts
   - Uses outdated interface assumptions
   - Tries to call methods that were never implemented

2. **credential-management** effort:
   - Provides `NewTrustStore()` returning `*DefaultTrustStore`
   - Does NOT provide `ExtractCertificate` function
   - Has certificate management but different interface

3. **image-operations** effort:
   - `gitea.Client` has `Push()` method, not `PushImage()`
   - `NewClient()` takes 2 params: (registryURL, *DefaultTrustStore)
   - No `ValidateCredentials()` method
   - No `NewClientAutoDetect()` function

## Affected Efforts

### 1. cli-commands (Primary Changes Needed)
**Branch**: `phase2/wave2/cli-commands`
**Status**: NEEDS MAJOR FIXES

### 2. credential-management (Supporting Changes)
**Branch**: `phase2/wave2/credential-management`
**Status**: NEEDS MINOR ADDITIONS

### 3. image-operations (No Changes)
**Branch**: `phase2/wave2/image-operations`
**Status**: INTERFACES OK AS-IS

## Fix Requirements

### For effort: cli-commands
**Branch**: `phase2/wave2/cli-commands`
**File**: `pkg/cmd/push.go`

Changes needed:

1. **Remove the unused import** (line 9):
   - Remove: `"github.com/cnoe-io/idpbuilder/pkg/cmd/helpers"`
   - This import is in the integration but not used

2. **Fix certificate extraction** (lines 62-75):
   - Replace `certs.ExtractCertificate(cmd.Context())`
   - With: Direct use of `certs.NewTrustStore()` since certificate extraction is handled internally

3. **Fix NewClient parameters** (line 87):
   - Current: `gitea.NewClient(pushRegistry, pushUsername, pushToken, certPath)`
   - Fix to: `gitea.NewClient(pushRegistry, certManager)`
   - Remove username/token from NewClient call

4. **Remove NewClientAutoDetect** (lines 89-91):
   - Replace entire auto-detect logic
   - With: Use credential manager's built-in credential detection

5. **Remove ValidateCredentials call** (lines 98-101):
   - Delete the entire validation block
   - Credentials are validated during Push operation

6. **Fix PushImage to Push** (line 105):
   - Replace: `client.PushImage(cmd.Context(), imageName)`
   - With: `client.Push(imageName, nil)` or with progress channel

Updated implementation for lines 51-111:
```go
func runPush(cmd *cobra.Command, args []string) error {
	imageName := args[0]

	// Ensure image has a tag
	if !strings.Contains(imageName, ":") {
		imageName = imageName + ":latest"
	}

	// Create client with appropriate configuration
	var client *gitea.Client
	var err error

	if pushInsecure {
		fmt.Printf("⚠ Warning: Running in insecure mode - certificate verification disabled\n")
		client, err = gitea.NewInsecureClient(pushRegistry)
	} else {
		// Use certificate trust store
		certManager := certs.NewTrustStore()
		client, err = gitea.NewClient(pushRegistry, certManager)
		if err != nil {
			return fmt.Errorf("failed to create registry client: %w", err)
		}
		fmt.Printf("✓ Certificate configuration enabled\n")
	}

	if err != nil {
		return fmt.Errorf("failed to create registry client: %w", err)
	}

	// Set credentials if provided via flags
	if pushUsername != "" && pushToken != "" {
		client.SetCredentials(pushUsername, pushToken)
		fmt.Printf("✓ Using provided credentials\n")
	} else {
		fmt.Printf("✓ Using auto-detected credentials\n")
	}

	fmt.Printf("✓ Registry: %s\n", pushRegistry)

	// Push the image (with optional progress tracking)
	fmt.Printf("Pushing %s to %s...\n", imageName, pushRegistry)
	if err := client.Push(imageName, nil); err != nil {
		return fmt.Errorf("failed to push image: %w", err)
	}

	fmt.Printf("✓ Image pushed successfully\n")
	return nil
}
```

### For effort: credential-management
**Branch**: `phase2/wave2/credential-management`
**File**: `pkg/certs/extractor.go` (or new file)

Changes needed:

1. **Add ExtractCertificate helper** (optional - for backward compatibility):
   ```go
   // ExtractCertificate provides a simplified interface for certificate extraction
   // Returns a Certificate struct with the extracted certificate data
   func ExtractCertificate(ctx context.Context) (*Certificate, error) {
       // Use existing extraction logic
       extractor := NewKindCertExtractor(/* default config */)
       return extractor.Extract(ctx)
   }
   ```

   Note: This is OPTIONAL since cli-commands can be fixed to not need this.

### For effort: image-operations
**Branch**: `phase2/wave2/image-operations`
**File**: No changes needed

The image-operations effort has the correct interfaces:
- `NewClient(registryURL string, certManager *certs.DefaultTrustStore)` ✓
- `NewInsecureClient(registryURL string)` ✓
- `Push(imageRef string, progressChan chan<- PushProgress)` ✓
- `SetCredentials(username, password string)` ✓

## Interface Alignment Strategy

1. **Primary Fix Strategy**: Modify cli-commands to use the actual interfaces
   - This is the least invasive approach
   - Only one effort needs changes
   - Maintains existing functionality

2. **Alternative Strategy** (NOT RECOMMENDED): Add wrapper functions
   - Would require changes to multiple efforts
   - More complex and error-prone
   - Could introduce new bugs

## Verification Steps

After fixes are applied to effort branches:

1. **Build each effort independently**:
   ```bash
   cd efforts/phase2/wave2/cli-commands
   go build ./...

   cd efforts/phase2/wave2/credential-management
   go build ./...

   cd efforts/phase2/wave2/image-operations
   go build ./...
   ```

2. **Run unit tests in each effort**:
   ```bash
   go test ./pkg/... -v
   ```

3. **Re-attempt integration**:
   ```bash
   # After fixes are applied and pushed to effort branches
   # Re-run integration to verify build success
   ```

## R300 Compliance Statement

✅ **ALL fixes will be applied to effort branches, NOT the integration branch.**

As per R300 (Comprehensive Fix Management Protocol - SUPREME LAW):
- cli-commands fixes will be applied to: `phase2/wave2/cli-commands` branch
- credential-management fixes (if needed) to: `phase2/wave2/credential-management` branch
- NO changes will be made to: `phase2-wave2-integration` branch
- Integration branch will be rebuilt after effort fixes are complete

## Implementation Priority

1. **IMMEDIATE** (Required for build): Fix cli-commands/pkg/cmd/push.go
2. **OPTIONAL** (Nice to have): Add ExtractCertificate helper to credential-management
3. **NO CHANGE**: image-operations already has correct interfaces

## Time Estimate

- Analysis complete: 5 minutes ✓
- Fix implementation: 10-15 minutes per effort
- Verification: 5 minutes
- Total: Within 30-minute R156 target

## Next Steps for Orchestrator

1. Spawn SW Engineer agents to apply fixes to effort branches
2. Each agent gets specific fix instructions for their effort
3. After fixes complete, re-run integration build
4. Verify R291 build gate passes
5. Continue with integration process

---

**Created**: 2025-09-15 23:37:51 UTC
**Agent**: code-reviewer
**State**: FIX_PLAN_CREATION
**Recovery Type**: R291 Build Gate Failure