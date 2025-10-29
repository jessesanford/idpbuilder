package auth

import "fmt"

// CredentialValidationError indicates credential validation failed.
type CredentialValidationError struct {
	Field  string // "username" or "password"
	Reason string
}

func (e *CredentialValidationError) Error() string {
	return fmt.Sprintf("credential validation failed (%s): %s", e.Field, e.Reason)
}
