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
