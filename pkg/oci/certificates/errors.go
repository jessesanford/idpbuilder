package certificates

import "fmt"

// ValidationError represents errors that occur during certificate validation.
type ValidationError struct {
	Code          string `json:"code"`
	Message       string `json:"message"`
	Rule          string `json:"rule,omitempty"`
	CertificateID string `json:"certificate_id,omitempty"`
}

func (e ValidationError) Error() string {
	if e.Rule != "" {
		return fmt.Sprintf("validation error [%s:%s]: %s", e.Code, e.Rule, e.Message)
	}
	return fmt.Sprintf("validation error [%s]: %s", e.Code, e.Message)
}

func NewValidationError(code, message string) ValidationError {
	return ValidationError{Code: code, Message: message}
}

// StorageError represents errors that occur during storage operations.
type StorageError struct {
	Code          string `json:"code"`
	Message       string `json:"message"`
	Operation     string `json:"operation,omitempty"`
	CertificateID string `json:"certificate_id,omitempty"`
	Cause         error  `json:"-"`
}

func (e StorageError) Error() string {
	if e.Operation != "" {
		return fmt.Sprintf("storage error [%s:%s]: %s", e.Code, e.Operation, e.Message)
	}
	return fmt.Sprintf("storage error [%s]: %s", e.Code, e.Message)
}

func (e StorageError) Unwrap() error { return e.Cause }

func NewStorageError(code, message string, cause error) StorageError {
	return StorageError{Code: code, Message: message, Cause: cause}
}

// ConfigError represents errors that occur during configuration operations.
type ConfigError struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Field   string      `json:"field,omitempty"`
	Value   interface{} `json:"value,omitempty"`
}

func (e ConfigError) Error() string {
	if e.Field != "" {
		return fmt.Sprintf("config error [%s:%s]: %s", e.Code, e.Field, e.Message)
	}
	return fmt.Sprintf("config error [%s]: %s", e.Code, e.Message)
}

func NewConfigError(code, message string) ConfigError {
	return ConfigError{Code: code, Message: message}
}

// Common error constants for consistent error handling.
const (
	// Validation error codes
	ErrCodeCertNotFound         = "CERT_NOT_FOUND"
	ErrCodeCertExpired          = "CERT_EXPIRED"
	ErrCodeInvalidPEM           = "INVALID_PEM"
	ErrCodeUnsupportedKeyUsage  = "UNSUPPORTED_KEY_USAGE"
	ErrCodeInvalidCertChain     = "INVALID_CERT_CHAIN"
	ErrCodeValidationRuleExists = "VALIDATION_RULE_EXISTS"

	// Storage error codes
	ErrCodeStorageUnavailable = "STORAGE_UNAVAILABLE"
	ErrCodeDuplicateCert      = "DUPLICATE_CERT"
	ErrCodeCertInUse          = "CERT_IN_USE"
	ErrCodePoolNotFound       = "POOL_NOT_FOUND"
	ErrCodePoolExists         = "POOL_EXISTS"

	// Configuration error codes
	ErrCodeInvalidConfig      = "INVALID_CONFIG"
	ErrCodeMissingConfig      = "MISSING_CONFIG"
	ErrCodeInvalidStorageType = "INVALID_STORAGE_TYPE"
	ErrCodeInvalidEventConfig = "INVALID_EVENT_CONFIG"
)