package errors

import "fmt"

// ErrorCode represents a numeric error code for programmatic error handling.
// Error codes are organized into ranges by functional area to avoid conflicts.
type ErrorCode int

const (
	// Authentication Errors (1000-1999)
	CodeAuthenticationFailed    ErrorCode = 1001
	CodeInvalidCredentials      ErrorCode = 1002
	CodeTokenExpired           ErrorCode = 1003
	CodeInvalidToken           ErrorCode = 1004
	CodeAuthorizationDenied    ErrorCode = 1005
	CodeCertificateInvalid     ErrorCode = 1006
	CodeTLSHandshakeFailed     ErrorCode = 1007

	// Registry Errors (2000-2999)
	CodeRegistryUnavailable    ErrorCode = 2001
	CodeRegistryTimeout        ErrorCode = 2002
	CodeRegistryNotFound       ErrorCode = 2003
	CodeRegistryUnauthorized   ErrorCode = 2004
	CodeRegistryServerError    ErrorCode = 2005
	CodeRegistryRateLimit      ErrorCode = 2006
	CodeRegistryMaintenance    ErrorCode = 2007

	// Manifest Errors (3000-3999)
	CodeManifestNotFound       ErrorCode = 3001
	CodeManifestInvalid        ErrorCode = 3002
	CodeManifestUnsupported    ErrorCode = 3003
	CodeManifestCorrupted      ErrorCode = 3004
	CodeDigestMismatch         ErrorCode = 3005
	CodeManifestTooLarge       ErrorCode = 3006

	// Network Errors (4000-4999)
	CodeNetworkTimeout         ErrorCode = 4001
	CodeConnectionRefused      ErrorCode = 4002
	CodeNetworkUnreachable     ErrorCode = 4003
	CodeDNSResolutionFailed    ErrorCode = 4004
	CodeProxyError             ErrorCode = 4005
	CodeSSLError               ErrorCode = 4006

	// Validation Errors (5000-5999)
	CodeInvalidInput           ErrorCode = 5001
	CodeMissingRequiredField   ErrorCode = 5002
	CodeInvalidFormat          ErrorCode = 5003
	CodeValueOutOfRange        ErrorCode = 5004
	CodeConflictingValues      ErrorCode = 5005
	CodeInvalidConfiguration   ErrorCode = 5006
)

// ErrorCategory represents the category of an error, which determines
// the appropriate handling strategy (retry, user action, etc.).
type ErrorCategory string

const (
	// CategoryTransient indicates errors that may succeed if retried
	CategoryTransient ErrorCategory = "transient"

	// CategoryPermanent indicates errors that will not succeed if retried
	CategoryPermanent ErrorCategory = "permanent"

	// CategoryConfiguration indicates errors caused by incorrect configuration
	CategoryConfiguration ErrorCategory = "configuration"

	// CategoryAuthentication indicates authentication/authorization errors
	CategoryAuthentication ErrorCategory = "authentication"

	// CategoryValidation indicates input validation errors
	CategoryValidation ErrorCategory = "validation"

	// CategoryNetwork indicates network-related errors
	CategoryNetwork ErrorCategory = "network"
)

// String returns the string representation of the error category.
func (ec ErrorCategory) String() string {
	return string(ec)
}

// IsRetryable returns true if errors in this category can be retried.
func (ec ErrorCategory) IsRetryable() bool {
	return ec == CategoryTransient || ec == CategoryNetwork
}

// codeCategories maps error codes to their categories
var codeCategories = map[ErrorCode]ErrorCategory{
	// Authentication errors - permanent until credentials fixed
	CodeAuthenticationFailed: CategoryAuthentication,
	CodeInvalidCredentials:   CategoryAuthentication,
	CodeTokenExpired:        CategoryAuthentication,
	CodeInvalidToken:        CategoryAuthentication,
	CodeAuthorizationDenied: CategoryAuthentication,
	CodeCertificateInvalid:  CategoryAuthentication,
	CodeTLSHandshakeFailed:  CategoryAuthentication,

	// Registry errors - mix of transient and permanent
	CodeRegistryUnavailable:  CategoryTransient,
	CodeRegistryTimeout:      CategoryTransient,
	CodeRegistryNotFound:     CategoryPermanent,
	CodeRegistryUnauthorized: CategoryAuthentication,
	CodeRegistryServerError:  CategoryTransient,
	CodeRegistryRateLimit:    CategoryTransient,
	CodeRegistryMaintenance:  CategoryTransient,

	// Manifest errors - mostly permanent
	CodeManifestNotFound:    CategoryPermanent,
	CodeManifestInvalid:     CategoryPermanent,
	CodeManifestUnsupported: CategoryPermanent,
	CodeManifestCorrupted:   CategoryPermanent,
	CodeDigestMismatch:      CategoryPermanent,
	CodeManifestTooLarge:    CategoryPermanent,

	// Network errors - transient
	CodeNetworkTimeout:      CategoryNetwork,
	CodeConnectionRefused:   CategoryNetwork,
	CodeNetworkUnreachable: CategoryNetwork,
	CodeDNSResolutionFailed: CategoryNetwork,
	CodeProxyError:          CategoryNetwork,
	CodeSSLError:            CategoryNetwork,

	// Validation errors - configuration
	CodeInvalidInput:           CategoryValidation,
	CodeMissingRequiredField:   CategoryValidation,
	CodeInvalidFormat:          CategoryValidation,
	CodeValueOutOfRange:        CategoryValidation,
	CodeConflictingValues:      CategoryValidation,
	CodeInvalidConfiguration:   CategoryConfiguration,
}

// GetCategory returns the category for a given error code.
func (ec ErrorCode) GetCategory() ErrorCategory {
	if category, exists := codeCategories[ec]; exists {
		return category
	}
	return CategoryPermanent // default to permanent if unknown
}

// String returns the string representation of the error code.
func (ec ErrorCode) String() string {
	return fmt.Sprintf("ErrorCode(%d)", int(ec))
}

// IsValid returns true if the error code is defined.
func (ec ErrorCode) IsValid() bool {
	_, exists := codeCategories[ec]
	return exists
}