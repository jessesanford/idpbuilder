package mapper

import "fmt"

// MappingError represents errors during mapping operations
type MappingError struct {
	Code    ErrorCode
	Message string
	Details map[string]string
}

// ErrorCode represents different types of mapping errors
type ErrorCode int

const (
	// ErrInvalidConfig indicates the stack configuration is invalid
	ErrInvalidConfig ErrorCode = iota
	// ErrMissingComponent indicates a required component is missing
	ErrMissingComponent
	// ErrInvalidReference indicates a component reference is invalid
	ErrInvalidReference
	// ErrValidationFailed indicates mapping validation failed
	ErrValidationFailed
)

// String returns a string representation of the error code
func (e ErrorCode) String() string {
	switch e {
	case ErrInvalidConfig:
		return "INVALID_CONFIG"
	case ErrMissingComponent:
		return "MISSING_COMPONENT"
	case ErrInvalidReference:
		return "INVALID_REFERENCE"
	case ErrValidationFailed:
		return "VALIDATION_FAILED"
	default:
		return "UNKNOWN_ERROR"
	}
}

// newMappingError creates a new mapping error with the specified code and message
func newMappingError(code ErrorCode, msg string) *MappingError {
	return &MappingError{
		Code:    code,
		Message: msg,
		Details: make(map[string]string),
	}
}

// Error implements the error interface
func (e *MappingError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code.String(), e.Message)
}

// WithDetail adds additional detail to the error
func (e *MappingError) WithDetail(key, value string) *MappingError {
	e.Details[key] = value
	return e
}