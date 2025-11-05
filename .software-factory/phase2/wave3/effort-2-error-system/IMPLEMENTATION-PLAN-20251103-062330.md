# Effort Implementation Plan

**Effort ID**: E2.3.2
**Effort Name**: Error Type System & Exit Code Mapping
**Phase**: 2 - Core Push Functionality
**Wave**: 3 - Error Handling & Validation
**Created By**: Code Reviewer (agent-code-reviewer)
**Date Created**: 2025-11-03
**Assigned To**: SW Engineer (TBD)

## 📋 Effort Overview

### Description
Implement comprehensive error type system with exit code mapping, error wrapping with context preservation, and integration with Wave 2.1/2.2 push command for actionable error messages. This effort moves error types from pkg/validator/types.go to pkg/errors/ and extends them with full error wrapping support, exit code constants, and integration with the runPush function.

### Size Estimate
- **Estimated Lines**: 350 (implementation only, tests excluded per R007)
- **Confidence Level**: High
- **Split Risk**: Low

### Dependencies
- **Requires**: Effort 2.3.1 (Input Validation & Security Checks) - MUST COMPLETE FIRST
  - Imports: ValidationError, SSRFWarning, SecurityWarning types from pkg/validator/types.go
  - Imports: ValidateImageName, ValidateRegistryURL, ValidateCredentials functions
- **Blocks**: None (final effort in Wave 2.3)
- **External**:
  - Wave 2.1: Push command core (runPush function, PushOptions)
  - Wave 2.2: Configuration system (PushConfig)
  - Standard library: errors, fmt, strings, os

## 🚨 R213 EFFORT METADATA (MANDATORY)

```yaml
effort_id: "2.3.2"
effort_name: "Error Type System & Exit Code Mapping"
branch_name: "idpbuilder-oci-push/phase2/wave3/effort-2.3.2-error-handling"
base_branch: "idpbuilder-oci-push/phase2/wave3/effort-2.3.1-input-validation"
parent_wave: "wave-2.3"
parent_phase: "phase-2"
depends_on:
  - "effort:2.3.1"  # Needs ValidationError types
dependencies_detail:
  - type: "effort"
    target: "2.3.1"
    reason: "Imports ValidationError, SSRFWarning, SecurityWarning types"
  - type: "code"
    target: "pkg/validator/types.go"
    reason: "Moves error types to pkg/errors/ and extends them"
  - type: "integration"
    target: "pkg/cmd/push/push.go"
    reason: "Wraps errors in runPush function"
estimated_lines: 350
complexity: "medium"
can_parallelize: false
parallel_with: []
files_touched:
  - "pkg/errors/types.go"              # new (moves from pkg/validator/types.go)
  - "pkg/errors/exitcodes.go"          # new
  - "pkg/cmd/push/errors.go"           # new
  - "pkg/errors/types_test.go"         # new
  - "pkg/cmd/push/push_errors_test.go" # new
  - "pkg/cmd/push/push.go"             # modified (add error wrapping)
  - "pkg/validator/imagename.go"       # modified (update imports)
  - "pkg/validator/registry.go"        # modified (update imports)
  - "pkg/validator/credentials.go"     # modified (update imports)
risk_level: "medium"
risk_reason: "Changes exit code behavior of push command"
test_count: 30
```

## 🎯 Requirements

### Functional Requirements
- [ ] Move error types from pkg/validator/types.go to pkg/errors/types.go
- [ ] Extend error types with BaseError struct and Unwrap support
- [ ] Implement exit code mapping (1=validation, 2=auth, 3=network, 4=image not found)
- [ ] Create error wrapping functions (wrapDockerError, wrapRegistryError)
- [ ] Integrate error wrapping into runPush function
- [ ] Add exit code handling to Cobra command
- [ ] Format errors with "Error: X, Suggestion: Y" pattern
- [ ] Update imports in pkg/validator/*.go files

### Non-Functional Requirements
- [ ] Performance: Error wrapping adds <1ms overhead per operation
- [ ] Security: Error messages must not leak sensitive information (credentials)
- [ ] Maintainability: Error types follow standard Go error wrapping patterns (errors.As, errors.Unwrap)
- [ ] Scalability: Exit code mapping extensible for future error types

### Acceptance Criteria
- [ ] All 30 unit tests passing (18 error types tests + 12 wrapping tests)
- [ ] Test coverage ≥95% statement, ≥90% branch for pkg/errors
- [ ] Code review approved
- [ ] Size ≤350 lines (measured with line-counter.sh, tests excluded per R007)
- [ ] No critical TODOs
- [ ] All public functions have godoc comments
- [ ] Exit code mapping verified (returns 1, 2, 3, 4, 0)
- [ ] Error wrapping preserves context and type information
- [ ] Integration with Wave 2.1 runPush verified
- [ ] Backward compatibility with Wave 2.2 configuration maintained
- [ ] Scope boundaries followed (R311)

## 🚨 EXPLICIT SCOPE DEFINITION (R311 MANDATORY)

### IMPLEMENT EXACTLY

#### Functions (EXACTLY 8 functions)
```
1. NewValidationError(field, message, suggestion) *ValidationError      // ~8 lines  - Create validation error
2. NewAuthenticationError(registry, message, suggestion) *AuthenticationError  // ~8 lines  - Create auth error
3. NewNetworkError(target, message, suggestion) *NetworkError           // ~8 lines  - Create network error
4. NewImageNotFoundError(imageName, message, suggestion) *ImageNotFoundError   // ~8 lines  - Create image error
5. GetExitCode(err error) int                                           // ~30 lines - Map error to exit code
6. FormatError(err error) string                                        // ~10 lines - Format error message
7. wrapDockerError(err error, imageName string) error                   // ~40 lines - Wrap Docker errors
8. wrapRegistryError(err error, registry string) error                  // ~45 lines - Wrap registry errors
9. validatePushOptions(opts *PushOptions) error                         // ~30 lines - Validate push options
TOTAL: 9 functions, ~187 lines
```

#### Types/Models (EXACTLY 5 types)
```
1. BaseError - Base error struct with Message, Suggestion, Cause fields - 3 methods (Error, Unwrap)
2. ValidationError - Embeds BaseError, adds Field and ExitCode - NO additional methods
3. AuthenticationError - Embeds BaseError, adds Registry and ExitCode - NO additional methods
4. NetworkError - Embeds BaseError, adds Target and ExitCode - NO additional methods
5. ImageNotFoundError - Embeds BaseError, adds ImageName and ExitCode - NO additional methods
TOTAL: 5 types, ~100 lines (including methods)
```

#### Constants (EXACTLY 5 constants)
```
ExitSuccess         = 0
ExitValidationError = 1
ExitAuthError       = 2
ExitNetworkError    = 3
ExitImageNotFound   = 4
ExitGenericError    = 1
TOTAL: 6 constants, ~10 lines
```

### 🛑 DO NOT IMPLEMENT
**CRITICAL: These are FORBIDDEN in this effort**

- ❌ DO NOT implement validation logic - already in Effort 2.3.1
- ❌ DO NOT implement SSRF detection - already in Effort 2.3.1
- ❌ DO NOT implement command injection prevention - already in Effort 2.3.1
- ❌ DO NOT add retry logic or error recovery
- ❌ DO NOT implement progress reporting (handled in Wave 2.1)
- ❌ DO NOT modify configuration system (handled in Wave 2.2)
- ❌ DO NOT add logging infrastructure
- ❌ DO NOT create custom error types beyond the 5 specified
- ❌ DO NOT add error rate limiting or circuit breakers
- ❌ DO NOT implement error telemetry or metrics
- ❌ DO NOT refactor existing Wave 2.1/2.2 code beyond error wrapping
- ❌ DO NOT add comprehensive edge case handling beyond string matching

### SIZE BREAKDOWN
```
Production Code:
- pkg/errors/types.go:      100 lines (5 types + constructors)
- pkg/errors/exitcodes.go:   60 lines (constants + mapping)
- pkg/cmd/push/errors.go:   100 lines (wrapping functions)
- pkg/cmd/push/push.go:     +80 lines (modifications)
- pkg/validator imports:     +10 lines (3 files × ~3 lines)
Subtotal:                    350 lines

Test Code (EXCLUDED per R007):
- pkg/errors/types_test.go:        150 lines (18 tests)
- pkg/cmd/push/push_errors_test.go: 200 lines (12 tests)
Subtotal:                           350 lines (NOT counted in size limit)

GRAND TOTAL: 350 lines implementation (within 900-line enforcement threshold per R535)
             + 350 lines tests (excluded per R007)
             = 700 lines total code
```

## 📁 Implementation Details

### Files to Create
| File Path | Purpose | Estimated Lines |
|-----------|---------|-----------------|
| `pkg/errors/types.go` | Error type definitions with BaseError and constructors | 100 |
| `pkg/errors/exitcodes.go` | Exit code constants and mapping function | 60 |
| `pkg/cmd/push/errors.go` | Error wrapping functions for Docker/registry errors | 100 |
| `pkg/errors/types_test.go` | Unit tests for error types and exit codes (18 tests) | 150 |
| `pkg/cmd/push/push_errors_test.go` | Integration tests for error wrapping (12 tests) | 200 |
| **Total** | | 610 |

### Files to Modify
| File Path | Changes | Estimated Lines |
|-----------|---------|-----------------|
| `pkg/cmd/push/push.go` | Add error wrapping to runPush, add exit code handling to Cobra command | +80 |
| `pkg/validator/imagename.go` | Update import from pkg/validator types to pkg/errors | +3 |
| `pkg/validator/registry.go` | Update import from pkg/validator types to pkg/errors | +3 |
| `pkg/validator/credentials.go` | Update import from pkg/validator types to pkg/errors | +4 |
| **Total** | | +90 |

### Key Components

#### Component 1: Error Type System (pkg/errors/types.go)
**Purpose**: Define all custom error types with BaseError embedding and constructor functions
**Location**: `pkg/errors/types.go`
**Lines**: ~100

```go
package errors

import (
	"fmt"
)

// BaseError is the base type for all custom errors
type BaseError struct {
	Message    string
	Suggestion string
	Cause      error
}

func (e *BaseError) Error() string {
	if e.Suggestion != "" {
		return fmt.Sprintf("Error: %s\nSuggestion: %s", e.Message, e.Suggestion)
	}
	return fmt.Sprintf("Error: %s", e.Message)
}

func (e *BaseError) Unwrap() error {
	return e.Cause
}

// ValidationError represents input validation failures (exit code 1)
type ValidationError struct {
	BaseError
	Field    string
	ExitCode int
}

func NewValidationError(field, message, suggestion string) *ValidationError {
	return &ValidationError{
		BaseError: BaseError{
			Message:    message,
			Suggestion: suggestion,
		},
		Field:    field,
		ExitCode: 1,
	}
}

// AuthenticationError represents authentication failures (exit code 2)
type AuthenticationError struct {
	BaseError
	Registry string
	ExitCode int
}

func NewAuthenticationError(registry, message, suggestion string) *AuthenticationError {
	return &AuthenticationError{
		BaseError: BaseError{
			Message:    message,
			Suggestion: suggestion,
		},
		Registry: registry,
		ExitCode: 2,
	}
}

// NetworkError represents network/connectivity failures (exit code 3)
type NetworkError struct {
	BaseError
	Target   string
	ExitCode int
}

func NewNetworkError(target, message, suggestion string) *NetworkError {
	return &NetworkError{
		BaseError: BaseError{
			Message:    message,
			Suggestion: suggestion,
		},
		Target:   target,
		ExitCode: 3,
	}
}

// ImageNotFoundError represents missing image errors (exit code 4)
type ImageNotFoundError struct {
	BaseError
	ImageName string
	ExitCode  int
}

func NewImageNotFoundError(imageName, message, suggestion string) *ImageNotFoundError {
	return &ImageNotFoundError{
		BaseError: BaseError{
			Message:    message,
			Suggestion: suggestion,
		},
		ImageName: imageName,
		ExitCode:  4,
	}
}

// SSRFWarning represents a potential SSRF risk (warning, not error)
type SSRFWarning struct {
	Target     string
	Message    string
	Suggestion string
}

func (w *SSRFWarning) Error() string {
	return fmt.Sprintf("Warning: %s\nSuggestion: %s", w.Message, w.Suggestion)
}

// SecurityWarning represents security concerns (warning, not error)
type SecurityWarning struct {
	Message    string
	Suggestion string
}

func (w *SecurityWarning) Error() string {
	return fmt.Sprintf("Warning: %s\nSuggestion: %s", w.Message, w.Suggestion)
}

// IsWarning returns true if the error is a warning (should not stop execution)
func IsWarning(err error) bool {
	_, isSSRF := err.(*SSRFWarning)
	_, isSecurity := err.(*SecurityWarning)
	return isSSRF || isSecurity
}
```

#### Component 2: Exit Code Mapping (pkg/errors/exitcodes.go)
**Purpose**: Define exit code constants and provide mapping function
**Location**: `pkg/errors/exitcodes.go`
**Lines**: ~60

```go
package errors

import (
	"errors"
	"fmt"
)

// Exit codes for different error types
const (
	ExitSuccess         = 0
	ExitValidationError = 1
	ExitAuthError       = 2
	ExitNetworkError    = 3
	ExitImageNotFound   = 4
	ExitGenericError    = 1
)

// GetExitCode maps an error to its appropriate exit code
func GetExitCode(err error) int {
	if err == nil {
		return ExitSuccess
	}

	// Check for typed errors using errors.As
	var validationErr *ValidationError
	if errors.As(err, &validationErr) {
		return ExitValidationError
	}

	var authErr *AuthenticationError
	if errors.As(err, &authErr) {
		return ExitAuthError
	}

	var networkErr *NetworkError
	if errors.As(err, &networkErr) {
		return ExitNetworkError
	}

	var imageNotFoundErr *ImageNotFoundError
	if errors.As(err, &imageNotFoundErr) {
		return ExitImageNotFound
	}

	// Default to validation error for untyped errors
	return ExitGenericError
}

// FormatError formats an error with consistent styling
func FormatError(err error) string {
	if err == nil {
		return ""
	}

	// Warnings have special formatting
	if IsWarning(err) {
		return fmt.Sprintf("⚠️  %s", err.Error())
	}

	// Regular errors
	return fmt.Sprintf("❌ %s", err.Error())
}
```

#### Component 3: Error Wrapping (pkg/cmd/push/errors.go)
**Purpose**: Wrap Docker and registry client errors with typed errors
**Location**: `pkg/cmd/push/errors.go`
**Lines**: ~100

```go
package push

import (
	"fmt"
	"strings"

	"github.com/cnoe-io/idpbuilder/pkg/errors"
	"github.com/cnoe-io/idpbuilder/pkg/validator"
)

// validatePushOptions validates all push options with security checks
//
// This extends Wave 2.2's PushConfig.Validate() with additional
// security validation from Wave 2.3.
func validatePushOptions(opts *PushOptions) error {
	// Validate image name
	if err := validator.ValidateImageName(opts.ImageName); err != nil {
		// validator returns typed errors, just pass through
		return err
	}

	// Validate registry URL
	if err := validator.ValidateRegistryURL(opts.Registry); err != nil {
		// Check if it's a warning
		if errors.IsWarning(err) {
			// Log warning but continue
			fmt.Println(errors.FormatError(err))
		} else {
			return err
		}
	}

	// Validate credentials
	if err := validator.ValidateCredentials(opts.Username, opts.Password); err != nil {
		if errors.IsWarning(err) {
			fmt.Println(errors.FormatError(err))
		} else {
			return err
		}
	}

	return nil
}

// wrapDockerError wraps Docker client errors with appropriate types
func wrapDockerError(err error, imageName string) error {
	if err == nil {
		return nil
	}

	errMsg := err.Error()

	// Check for "image not found" errors
	if strings.Contains(errMsg, "No such image") {
		return errors.NewImageNotFoundError(
			imageName,
			fmt.Sprintf("image '%s' not found in local Docker daemon", imageName),
			"pull the image first with: docker pull "+imageName,
		)
	}

	// Check for Docker daemon connection errors
	if strings.Contains(errMsg, "connection refused") || strings.Contains(errMsg, "Cannot connect") {
		return errors.NewNetworkError(
			"docker daemon",
			"cannot connect to Docker daemon",
			"ensure Docker daemon is running: systemctl start docker or start Docker Desktop",
		)
	}

	// Generic Docker error
	return fmt.Errorf("Docker error: %w", err)
}

// wrapRegistryError wraps registry client errors with appropriate types
func wrapRegistryError(err error, registry string) error {
	if err == nil {
		return nil
	}

	errMsg := err.Error()

	// Check for authentication errors
	if strings.Contains(errMsg, "401") || strings.Contains(errMsg, "unauthorized") {
		return errors.NewAuthenticationError(
			registry,
			fmt.Sprintf("authentication failed for registry %s", registry),
			"check your username and password, or verify registry credentials",
		)
	}

	// Check for network errors
	if strings.Contains(errMsg, "connection refused") || strings.Contains(errMsg, "timeout") {
		return errors.NewNetworkError(
			registry,
			fmt.Sprintf("cannot connect to registry %s", registry),
			"verify registry URL and network connectivity, or try with --insecure if using self-signed certificates",
		)
	}

	// TLS certificate errors
	if strings.Contains(errMsg, "x509") || strings.Contains(errMsg, "certificate") {
		return errors.NewNetworkError(
			registry,
			fmt.Sprintf("TLS certificate verification failed for %s", registry),
			"use --insecure flag to skip certificate verification (not recommended for production)",
		)
	}

	// Generic registry error
	return fmt.Errorf("registry error: %w", err)
}
```

## 🧪 Testing Strategy

### Unit Tests Required (COMPREHENSIVE)
- [ ] Test file 1: `pkg/errors/types_test.go` (18 tests)
- [ ] Test file 2: `pkg/cmd/push/push_errors_test.go` (12 tests)
- [ ] Coverage target: 95% statement, 90% branch
- [ ] Test cases:
  - [ ] Error type creation (4 constructor tests)
  - [ ] Error formatting (6 format tests)
  - [ ] Error unwrapping (2 unwrap tests)
  - [ ] Exit code mapping (8 mapping tests)
  - [ ] Docker error wrapping (4 wrapping tests)
  - [ ] Registry error wrapping (5 wrapping tests)
  - [ ] Context preservation (3 integration tests)

### Test Suites

#### Test Suite 1: Error Type Creation & Formatting (10 tests)
From WAVE-TEST-PLAN.md Test Suite 4:
- T-2.3.4-01: TestNewValidationError
- T-2.3.4-02: TestNewAuthenticationError
- T-2.3.4-03: TestNewNetworkError
- T-2.3.4-04: TestNewImageNotFoundError
- T-2.3.4-05: TestValidationError_Format
- T-2.3.4-06: TestAuthenticationError_Format
- T-2.3.4-07: TestSSRFWarning_Format
- T-2.3.4-08: TestBaseError_Unwrap
- T-2.3.4-09: TestErrorChain_Unwrap
- T-2.3.4-10: TestIsWarning_Detection

#### Test Suite 2: Exit Code Mapping (8 tests)
From WAVE-TEST-PLAN.md Test Suite 5:
- T-2.3.5-01: TestGetExitCode_ValidationError
- T-2.3.5-02: TestGetExitCode_AuthenticationError
- T-2.3.5-03: TestGetExitCode_NetworkError
- T-2.3.5-04: TestGetExitCode_ImageNotFoundError
- T-2.3.5-05: TestGetExitCode_GenericError
- T-2.3.5-06: TestGetExitCode_NilError
- T-2.3.5-07: TestGetExitCode_WrappedValidationError
- T-2.3.5-08: TestGetExitCode_WrappedAuthError

#### Test Suite 3: Error Wrapping Integration (12 tests)
From WAVE-TEST-PLAN.md Test Suite 6:
- T-2.3.6-01: TestWrapDockerError_ImageNotFound
- T-2.3.6-02: TestWrapDockerError_ConnectionRefused
- T-2.3.6-03: TestWrapDockerError_CannotConnect
- T-2.3.6-04: TestWrapDockerError_GenericError
- T-2.3.6-05: TestWrapRegistryError_Unauthorized
- T-2.3.6-06: TestWrapRegistryError_ConnectionRefused
- T-2.3.6-07: TestWrapRegistryError_Timeout
- T-2.3.6-08: TestWrapRegistryError_TLSError
- T-2.3.6-09: TestWrapRegistryError_GenericError
- T-2.3.6-10: TestWrapDockerError_PreservesImageName
- T-2.3.6-11: TestWrapRegistryError_PreservesRegistry
- T-2.3.6-12: TestWrapErrors_ChainUnwraps

### Test Data Requirements
- Mock Docker client errors (image not found, connection refused)
- Mock registry client errors (401, timeout, x509)
- Error chain test fixtures
- No external service stubs needed (unit tests only)

## 🔄 Implementation Approach

### Step 1: Create Error Type System
- [ ] Create pkg/errors/types.go with BaseError
- [ ] Add ValidationError, AuthenticationError, NetworkError, ImageNotFoundError
- [ ] Add constructor functions (NewValidationError, etc.)
- [ ] Add SSRFWarning and SecurityWarning
- [ ] Add IsWarning helper function
- **Estimated Time**: 2 hours
- **Lines**: ~100

### Step 2: Create Exit Code Mapping
- [ ] Create pkg/errors/exitcodes.go
- [ ] Define exit code constants (0-4)
- [ ] Implement GetExitCode function with errors.As
- [ ] Implement FormatError function
- **Estimated Time**: 1 hour
- **Lines**: ~60

### Step 3: Create Error Wrapping Functions
- [ ] Create pkg/cmd/push/errors.go
- [ ] Implement validatePushOptions function
- [ ] Implement wrapDockerError function
- [ ] Implement wrapRegistryError function
- **Estimated Time**: 2 hours
- **Lines**: ~100

### Step 4: Integrate with runPush
- [ ] Modify pkg/cmd/push/push.go
- [ ] Add validatePushOptions call at start of runPush
- [ ] Wrap Docker client errors with wrapDockerError
- [ ] Wrap registry client errors with wrapRegistryError
- [ ] Add exit code handling to Cobra command
- **Estimated Time**: 2 hours
- **Lines**: +80

### Step 5: Update Validator Imports
- [ ] Update pkg/validator/imagename.go import
- [ ] Update pkg/validator/registry.go import
- [ ] Update pkg/validator/credentials.go import
- [ ] Verify validator tests still pass
- **Estimated Time**: 0.5 hours
- **Lines**: +10

### Step 6: Write Tests (TDD Red Phase - BEFORE implementation)
- [ ] Write pkg/errors/types_test.go (18 tests)
- [ ] Write pkg/cmd/push/push_errors_test.go (12 tests)
- [ ] Run tests (should all fail initially - TDD Red phase)
- **Estimated Time**: 3 hours
- **Lines**: ~350 (NOT counted in size limit per R007)

### Step 7: Implement to Pass Tests (TDD Green Phase)
- [ ] Implement error types to pass tests
- [ ] Implement exit code mapping to pass tests
- [ ] Implement error wrapping to pass tests
- [ ] Achieve 95% coverage
- **Estimated Time**: 2 hours
- **Lines**: ~0 (already counted in steps 1-5)

### Step 8: Measure Size
- [ ] Run line-counter.sh to verify size
- [ ] Verify ≤350 lines implementation (tests excluded)
- [ ] Document actual line count
- **Estimated Time**: 0.5 hours

## 📊 Size Management

### Size Breakdown
```
Core Logic (error types):         100 lines (29%)
Exit code mapping:                 60 lines (17%)
Error wrapping:                   100 lines (29%)
Integration (push.go changes):     80 lines (23%)
Import updates:                    10 lines (2%)
---------------------------------------------------
Total Estimated:                  350 lines (100%)
Buffer (15%):                      53 lines
---------------------------------------------------
Final Estimate with Buffer:       403 lines (within 900-line enforcement threshold per R535)
```

### Split Contingency Plan
**Not needed** - effort is well under 900-line enforcement threshold (R535).

If effort unexpectedly exceeds 700 lines during implementation:

**Split Option 1**: By functionality
- Split 1: Error types + exit codes (~160 lines)
- Split 2: Error wrapping + integration (~190 lines)

**Split Option 2**: By package
- Split 1: pkg/errors package complete (~160 lines)
- Split 2: pkg/cmd/push integration (~190 lines)

## 🚨 Risk Assessment

### Technical Risks
| Risk | Probability | Impact | Mitigation |
|------|------------|--------|------------|
| Size overrun | Low | Medium | Early measurement at 200 lines, contingency splits defined |
| String matching brittle | Medium | Medium | Comprehensive error wrapping tests, document error patterns |
| Lost error context | Low | High | Error unwrapping tests, verify errors.As works |
| Exit code breaking change | Low | High | Backward compatibility tests, document exit code changes |

### Implementation Risks
- **Integration with Wave 2.1/2.2**: Error wrapping must not break existing push pipeline
  - Mitigation: Integration tests with mock Docker/registry clients
  - Mitigation: Verify all Wave 2.1/2.2 tests still pass
- **Error Type Detection**: String matching for error wrapping is brittle
  - Mitigation: Comprehensive error message tests
  - Mitigation: Document all error patterns checked
  - Future: Use typed errors from Docker/registry libraries
- **Exit Code Changes**: Scripts may depend on specific exit codes
  - Mitigation: Document all exit code changes
  - Mitigation: Provide backward compatibility mode if needed

## 🔗 Integration Points

### APIs/Interfaces
- **Consumes**:
  - pkg/validator: ValidateImageName, ValidateRegistryURL, ValidateCredentials
  - pkg/cmd/push: runPush function, PushOptions struct
  - Docker client errors (from Wave 2.1)
  - Registry client errors (from Wave 2.1)
- **Provides**:
  - pkg/errors: Error types, exit code mapping, error formatting
  - pkg/cmd/push: validatePushOptions, wrapDockerError, wrapRegistryError
- **Modifies**:
  - pkg/cmd/push/push.go: Adds error wrapping and exit code handling
  - pkg/validator/*.go: Changes imports from local types to pkg/errors

### Data Structures
- **Uses**:
  - PushOptions (from Wave 2.1/2.2)
  - Error interfaces (from Wave 2.1 Docker/registry clients)
- **Creates**:
  - BaseError, ValidationError, AuthenticationError, NetworkError, ImageNotFoundError
  - SSRFWarning, SecurityWarning (moved from pkg/validator)
- **Modifies**: None

### Dependencies Graph
```
This Effort (2.3.2)
    ↓ imports
Effort 2.3.1 (validator package)
    ↓ validation functions
Wave 2.2 (configuration system)
    ↓ PushConfig, PushOptions
Wave 2.1 (push command core)
    ↓ runPush, Docker client, Registry client
```

## 📈 Success Metrics

### Code Quality Metrics
- Cyclomatic complexity: <10 per function
- Code duplication: <5%
- Linting errors: 0 (golangci-lint)
- Type checking: 100% pass (go vet)

### Performance Metrics
- Error wrapping overhead: <1ms per operation
- Memory allocation: <100 bytes per error
- Error formatting: <0.5ms per error

### Coverage Metrics
- Statement coverage: ≥95%
- Branch coverage: ≥90%
- Function coverage: 100%

## 📚 Documentation Requirements

### Code Documentation
- [ ] All public functions have godoc comments
- [ ] Error type constructors documented with usage examples
- [ ] Exit code constants documented
- [ ] Error wrapping patterns documented with examples
- [ ] Integration points documented in push.go

### Integration Documentation
- [ ] Document error wrapping changes in Wave 2.1/2.2
- [ ] Document exit code behavior for scripts
- [ ] Document error message format changes
- [ ] Provide migration guide if needed

## ✅ Review Checklist

### For Implementation
- [ ] All 9 functions implemented
- [ ] All 5 error types created
- [ ] All 30 tests written and passing
- [ ] Size measured with line-counter.sh (≤350 lines implementation)
- [ ] No hardcoded error messages (all contextual)
- [ ] Error wrapping preserves context
- [ ] Exit code mapping tested for all error types
- [ ] Integration with Wave 2.1/2.2 verified
- [ ] Backward compatibility maintained
- [ ] All godoc comments complete

### For Code Review
- [ ] Measure size: cd $EFFORT_DIR && $PROJECT_ROOT/tools/line-counter.sh
- [ ] Verify test coverage: go test -cover ./pkg/errors ./pkg/cmd/push
- [ ] Check linting: golangci-lint run ./pkg/errors ./pkg/cmd/push
- [ ] Review error handling: All errors wrapped appropriately
- [ ] Validate integration: Wave 2.1/2.2 tests still pass
- [ ] Check exit codes: All 5 exit codes tested
- [ ] Approve or request changes

## 📝 Notes

### Implementation Notes
- Error types follow standard Go error wrapping patterns (errors.As, errors.Unwrap)
- String matching for error wrapping is intentionally simple (future: use typed errors from libraries)
- SSRFWarning and SecurityWarning are non-blocking (allow continuation)
- Exit code mapping is extensible (can add new error types in future waves)
- Error messages must not leak sensitive information (credentials masked)

### Assumptions
- Wave 2.1 push command core is complete and stable
- Wave 2.2 configuration system is complete and stable
- Effort 2.3.1 validation system will be complete before this effort starts
- Docker and registry client errors have predictable error messages
- Exit code changes will not break critical user scripts

### Open Questions
- None - all requirements clear from Wave Implementation Plan

### Integration with Effort 2.3.1
This effort has a HARD DEPENDENCY on Effort 2.3.1:
- **Imports**: ValidationError, SSRFWarning, SecurityWarning from pkg/validator/types.go
- **Moves**: These types from pkg/validator/types.go to pkg/errors/types.go
- **Extends**: Adds BaseError embedding and Unwrap support
- **Uses**: ValidateImageName, ValidateRegistryURL, ValidateCredentials in validatePushOptions

**Critical Path**: Effort 2.3.1 → Effort 2.3.2 (MUST BE SEQUENTIAL)

---

**Status**: Planning Complete
**Ready for Implementation**: Yes (pending Effort 2.3.1 completion)
**Approved By**: Code Reviewer (agent-code-reviewer)
**Approval Date**: 2025-11-03

**Orchestrator Instructions**:
1. Verify Effort 2.3.1 is COMPLETE before spawning SW Engineer for 2.3.2
2. SW Engineer must follow TDD protocol (R341): Write tests BEFORE implementation
3. SW Engineer must measure size with line-counter.sh after every 100 lines
4. SW Engineer must verify Wave 2.1/2.2 tests still pass after integration
5. Code Reviewer must verify exit code mapping and error wrapping in review

**Remember**: This plan is the contract between Code Reviewer and SW Engineer!
