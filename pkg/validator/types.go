package validator

import "fmt"

// ValidationError represents input validation failures (exit code 1)
type ValidationError struct {
	Field      string
	Message    string
	Suggestion string
	ExitCode   int
}

func (e *ValidationError) Error() string {
	if e.Suggestion != "" {
		return fmt.Sprintf("Error: %s\nSuggestion: %s", e.Message, e.Suggestion)
	}
	return fmt.Sprintf("Error: %s", e.Message)
}

// Note: SSRFWarning and SecurityWarning types have been moved to pkg/errors
// for centralized error handling. Use errors.SSRFWarning and errors.SecurityWarning instead.
