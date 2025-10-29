package auth

import "fmt"

// CredentialValidationError represents an error that occurs during credential validation.
//
// This error type is returned when credentials fail validation checks such as
// empty fields, invalid formats, or control characters in usernames.
type CredentialValidationError struct {
	// Field identifies which credential field failed validation (e.g., "username", "password")
	Field string

	// Reason provides a human-readable explanation of why validation failed
	Reason string
}

// Error implements the error interface.
func (e *CredentialValidationError) Error() string {
	return fmt.Sprintf("credential validation failed for %s: %s", e.Field, e.Reason)
}
