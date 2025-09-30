package retry

import "fmt"

// MaxRetriesExceededError is returned when an operation fails after exhausting all retry attempts.
type MaxRetriesExceededError struct {
	Attempts int
	LastErr  error
}

// Error implements the error interface.
func (e *MaxRetriesExceededError) Error() string {
	return fmt.Sprintf("max retries exceeded after %d attempts: %v", e.Attempts, e.LastErr)
}

// Unwrap returns the underlying error for error chain unwrapping.
func (e *MaxRetriesExceededError) Unwrap() error {
	return e.LastErr
}