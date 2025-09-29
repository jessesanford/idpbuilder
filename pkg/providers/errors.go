package providers

import (
	"fmt"
)

// ProviderError represents an error that occurred during a provider operation.
// It contains context about the operation that failed and the underlying error.
type ProviderError struct {
	// Op is the operation that was being performed when the error occurred
	Op string

	// Ref is the reference (registry/repo:tag) involved in the operation
	Ref string

	// Err is the underlying error that caused the failure
	Err error
}

// Error implements the error interface for ProviderError.
// It returns a formatted error message that includes the operation,
// reference, and underlying error details.
func (e ProviderError) Error() string {
	if e.Ref != "" {
		return fmt.Sprintf("provider operation %s failed for %s: %v", e.Op, e.Ref, e.Err)
	}
	return fmt.Sprintf("provider operation %s failed: %v", e.Op, e.Err)
}

// Unwrap returns the underlying error, allowing errors.Is and errors.As
// to work correctly with ProviderError.
func (e ProviderError) Unwrap() error {
	return e.Err
}