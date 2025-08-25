package errors

import "fmt"

// ErrorCode represents a specific error condition with category and code
type ErrorCode struct {
	Category string
	Code     string
	Message  string
	Severity string
}

// String returns a formatted string representation of the error code
func (ec ErrorCode) String() string {
	return fmt.Sprintf("%s.%s", ec.Category, ec.Code)
}

// Error code definitions organized by category

// Transient errors - may resolve with retry
var (
	// Network-related transient errors
	ErrorNetworkTimeout = ErrorCode{
		Category: CategoryTransient,
		Code:     "network_timeout",
		Message:  "Network operation timed out",
		Severity: SeverityMedium,
	}
	
	ErrorNetworkUnavailable = ErrorCode{
		Category: CategoryTransient,
		Code:     "network_unavailable",
		Message:  "Network service temporarily unavailable",
		Severity: SeverityMedium,
	}
	
	ErrorRateLimited = ErrorCode{
		Category: CategoryTransient,
		Code:     "rate_limited",
		Message:  "Request rate limit exceeded",
		Severity: SeverityLow,
	}
	
	ErrorServiceBusy = ErrorCode{
		Category: CategoryTransient,
		Code:     "service_busy",
		Message:  "Service temporarily busy",
		Severity: SeverityMedium,
	}
)

// Permanent errors - will not resolve with retry
var (
	ErrorNotFound = ErrorCode{
		Category: CategoryPermanent,
		Code:     "not_found",
		Message:  "Resource not found",
		Severity: SeverityMedium,
	}
	
	ErrorAlreadyExists = ErrorCode{
		Category: CategoryPermanent,
		Code:     "already_exists",
		Message:  "Resource already exists",
		Severity: SeverityMedium,
	}
	
	ErrorUnsupported = ErrorCode{
		Category: CategoryPermanent,
		Code:     "unsupported",
		Message:  "Operation not supported",
		Severity: SeverityHigh,
	}
	
	ErrorIncompatible = ErrorCode{
		Category: CategoryPermanent,
		Code:     "incompatible",
		Message:  "Resource or operation incompatible",
		Severity: SeverityHigh,
	}
)

// Configuration errors
var (
	ErrorInvalidConfig = ErrorCode{
		Category: CategoryConfiguration,
		Code:     "invalid_config",
		Message:  "Invalid configuration",
		Severity: SeverityHigh,
	}
	
	ErrorMissingConfig = ErrorCode{
		Category: CategoryConfiguration,
		Code:     "missing_config",
		Message:  "Required configuration missing",
		Severity: SeverityHigh,
	}
	
	ErrorConfigFormat = ErrorCode{
		Category: CategoryConfiguration,
		Code:     "config_format",
		Message:  "Configuration format invalid",
		Severity: SeverityHigh,
	}
)

// Validation errors
var (
	ErrorInvalidInput = ErrorCode{
		Category: CategoryValidation,
		Code:     "invalid_input",
		Message:  "Input validation failed",
		Severity: SeverityMedium,
	}
	
	ErrorMissingRequired = ErrorCode{
		Category: CategoryValidation,
		Code:     "missing_required",
		Message:  "Required field missing",
		Severity: SeverityMedium,
	}
	
	ErrorInvalidFormat = ErrorCode{
		Category: CategoryValidation,
		Code:     "invalid_format",
		Message:  "Invalid format",
		Severity: SeverityMedium,
	}
)

// Authentication errors
var (
	ErrorUnauthenticated = ErrorCode{
		Category: CategoryAuthentication,
		Code:     "unauthenticated",
		Message:  "Authentication required",
		Severity: SeverityHigh,
	}
	
	ErrorInvalidCredentials = ErrorCode{
		Category: CategoryAuthentication,
		Code:     "invalid_credentials",
		Message:  "Invalid credentials provided",
		Severity: SeverityHigh,
	}
	
	ErrorTokenExpired = ErrorCode{
		Category: CategoryAuthentication,
		Code:     "token_expired",
		Message:  "Authentication token expired",
		Severity: SeverityMedium,
	}
)

// Authorization errors
var (
	ErrorUnauthorized = ErrorCode{
		Category: CategoryAuthorization,
		Code:     "unauthorized",
		Message:  "Access denied",
		Severity: SeverityHigh,
	}
	
	ErrorInsufficientPermissions = ErrorCode{
		Category: CategoryAuthorization,
		Code:     "insufficient_permissions",
		Message:  "Insufficient permissions for operation",
		Severity: SeverityHigh,
	}
)

// Resource errors
var (
	ErrorResourceExhausted = ErrorCode{
		Category: CategoryResource,
		Code:     "resource_exhausted",
		Message:  "Resource capacity exhausted",
		Severity: SeverityHigh,
	}
	
	ErrorResourceLocked = ErrorCode{
		Category: CategoryResource,
		Code:     "resource_locked",
		Message:  "Resource is locked by another operation",
		Severity: SeverityMedium,
	}
	
	ErrorResourceCorrupted = ErrorCode{
		Category: CategoryResource,
		Code:     "resource_corrupted",
		Message:  "Resource data is corrupted",
		Severity: SeverityCritical,
	}
)

// Helper functions for working with error codes

// IsTransient returns true if the error code represents a transient condition
func IsTransient(code ErrorCode) bool {
	return code.Category == CategoryTransient
}

// IsPermanent returns true if the error code represents a permanent condition
func IsPermanent(code ErrorCode) bool {
	return code.Category == CategoryPermanent
}

// IsConfigurationError returns true if the error code represents a configuration issue
func IsConfigurationError(code ErrorCode) bool {
	return code.Category == CategoryConfiguration
}

// IsValidationError returns true if the error code represents a validation issue
func IsValidationError(code ErrorCode) bool {
	return code.Category == CategoryValidation
}

// IsAuthenticationError returns true if the error code represents an authentication issue
func IsAuthenticationError(code ErrorCode) bool {
	return code.Category == CategoryAuthentication
}

// IsAuthorizationError returns true if the error code represents an authorization issue
func IsAuthorizationError(code ErrorCode) bool {
	return code.Category == CategoryAuthorization
}

// IsResourceError returns true if the error code represents a resource issue
func IsResourceError(code ErrorCode) bool {
	return code.Category == CategoryResource
}

// IsCritical returns true if the error code represents a critical severity error
func IsCritical(code ErrorCode) bool {
	return code.Severity == SeverityCritical
}