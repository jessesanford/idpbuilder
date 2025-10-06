package errors

import "fmt"

// AuthenticationError represents authentication failures with a registry
type AuthenticationError struct {
	Registry string
	Cause    error
}

// Error implements the error interface
func (e *AuthenticationError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("authentication failed for registry %s: %v", e.Registry, e.Cause)
	}
	return fmt.Sprintf("authentication failed for registry %s", e.Registry)
}

// Unwrap returns the underlying error
func (e *AuthenticationError) Unwrap() error {
	return e.Cause
}

// IsRetryable indicates whether authentication errors should be retried
func (e *AuthenticationError) IsRetryable() bool {
	// Most authentication errors are not retryable (bad credentials, etc.)
	// but some network-related auth failures might be
	if e.Cause != nil {
		// TODO: Check if underlying error is network-related
		return false
	}
	return false
}

// InsecureRegistryError indicates that a registry requires insecure mode
type InsecureRegistryError struct {
	Registry string
}

// Error implements the error interface
func (e *InsecureRegistryError) Error() string {
	return fmt.Sprintf("registry %s requires --insecure flag for self-signed certificates", e.Registry)
}

// IsRetryable indicates that insecure registry errors are not retryable
// They require configuration changes
func (e *InsecureRegistryError) IsRetryable() bool {
	return false
}

// CredentialError represents credential validation or retrieval errors
type CredentialError struct {
	Source string
	Cause  error
}

// Error implements the error interface
func (e *CredentialError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("credential error from %s: %v", e.Source, e.Cause)
	}
	return fmt.Sprintf("credential error from %s", e.Source)
}

// Unwrap returns the underlying error
func (e *CredentialError) Unwrap() error {
	return e.Cause
}

// IsRetryable indicates that credential errors are generally not retryable
func (e *CredentialError) IsRetryable() bool {
	return false
}

// RegistryError represents general registry communication errors
type RegistryError struct {
	Registry   string
	Operation  string
	StatusCode int
	Cause      error
}

// Error implements the error interface
func (e *RegistryError) Error() string {
	if e.StatusCode > 0 {
		return fmt.Sprintf("registry %s %s failed with status %d: %v",
			e.Registry, e.Operation, e.StatusCode, e.Cause)
	}
	return fmt.Sprintf("registry %s %s failed: %v", e.Registry, e.Operation, e.Cause)
}

// Unwrap returns the underlying error
func (e *RegistryError) Unwrap() error {
	return e.Cause
}

// IsRetryable determines if the registry error should be retried based on status code
func (e *RegistryError) IsRetryable() bool {
	// Retry on server errors and rate limiting
	switch e.StatusCode {
	case 429, // Too Many Requests
		 500, // Internal Server Error
		 502, // Bad Gateway
		 503, // Service Unavailable
		 504: // Gateway Timeout
		return true
	case 401, // Unauthorized
		 403: // Forbidden
		return false // Auth errors shouldn't be retried
	default:
		// For other errors, check the underlying cause
		if e.Cause != nil {
			// TODO: Implement more sophisticated retry logic based on error type
			return false
		}
		return false
	}
}

// TimeoutError represents timeout errors during registry operations
type TimeoutError struct {
	Registry  string
	Operation string
	Duration  string
	Cause     error
}

// Error implements the error interface
func (e *TimeoutError) Error() string {
	return fmt.Sprintf("timeout after %s during %s with registry %s: %v",
		e.Duration, e.Operation, e.Registry, e.Cause)
}

// Unwrap returns the underlying error
func (e *TimeoutError) Unwrap() error {
	return e.Cause
}

// IsRetryable indicates that timeout errors are generally retryable
func (e *TimeoutError) IsRetryable() bool {
	return true
}