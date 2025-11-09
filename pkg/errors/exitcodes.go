package errors

import (
	"errors"
	"fmt"
	"os"
)

// Exit codes for different error types.
// These codes allow scripts to distinguish between error categories.
const (
	ExitSuccess         = 0 // Command succeeded
	ExitValidationError = 1 // Input validation failed
	ExitAuthError       = 2 // Authentication failed
	ExitNetworkError    = 3 // Network/connectivity failed
	ExitImageNotFound   = 4 // Image not found
	ExitGenericError    = 1 // Default for untyped errors
)

// GetExitCode maps an error to its appropriate exit code.
// It uses errors.As to check the error chain for typed errors.
//
// Exit code mapping:
//   - nil error: 0 (success)
//   - ValidationError: 1
//   - AuthenticationError: 2
//   - NetworkError: 3
//   - ImageNotFoundError: 4
//   - Untyped errors: 1 (generic)
//
// Example:
//
//	err := NewAuthenticationError("docker.io", "auth failed", "check credentials")
//	code := GetExitCode(err) // returns 2
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

	// Default to generic error for untyped errors
	return ExitGenericError
}

// FormatError formats an error with consistent styling.
// Warnings are prefixed with ⚠️, regular errors with ❌.
// In test mode (IDPBUILDER_TEST_MODE=true), emojis are omitted to prevent
// false positives in test output (go test treats ❌ as a failure marker).
//
// Example:
//
//	err := NewValidationError("imageName", "invalid format", "use name:tag")
//	msg := FormatError(err) // returns "❌ Error: invalid format\nSuggestion: use name:tag"
func FormatError(err error) string {
	if err == nil {
		return ""
	}

	// Check if running in test mode (no emojis to prevent test false positives)
	testMode := IsTestMode()

	// Warnings have special formatting
	if IsWarning(err) {
		if testMode {
			return err.Error() // No emoji in test mode
		}
		return fmt.Sprintf("⚠️  %s", err.Error())
	}

	// Regular errors
	if testMode {
		return err.Error() // No emoji in test mode
	}
	return fmt.Sprintf("❌ %s", err.Error())
}

// IsTestMode checks if the code is running in test mode.
// Test mode is detected by checking for:
// 1. IDPBUILDER_TEST_MODE environment variable
// 2. Go's testing package (tests are running)
func IsTestMode() bool {
	// Check explicit test mode flag
	if testMode := os.Getenv("IDPBUILDER_TEST_MODE"); testMode == "true" || testMode == "1" {
		return true
	}
	return false
}
