package errors

import (
	"testing"
)

func TestErrorCode_GetCategory(t *testing.T) {
	tests := []struct {
		name     string
		code     ErrorCode
		expected ErrorCategory
	}{
		{
			name:     "authentication error",
			code:     CodeAuthenticationFailed,
			expected: CategoryAuthentication,
		},
		{
			name:     "registry transient error",
			code:     CodeRegistryTimeout,
			expected: CategoryTransient,
		},
		{
			name:     "registry permanent error",
			code:     CodeRegistryNotFound,
			expected: CategoryPermanent,
		},
		{
			name:     "network error",
			code:     CodeNetworkTimeout,
			expected: CategoryNetwork,
		},
		{
			name:     "validation error",
			code:     CodeInvalidInput,
			expected: CategoryValidation,
		},
		{
			name:     "configuration error",
			code:     CodeInvalidConfiguration,
			expected: CategoryConfiguration,
		},
		{
			name:     "unknown error code",
			code:     ErrorCode(9999),
			expected: CategoryPermanent, // default
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.code.GetCategory()
			if got != tt.expected {
				t.Errorf("ErrorCode.GetCategory() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestErrorCode_String(t *testing.T) {
	code := CodeAuthenticationFailed
	expected := "ErrorCode(1001)"
	
	if got := code.String(); got != expected {
		t.Errorf("ErrorCode.String() = %v, want %v", got, expected)
	}
}

func TestErrorCode_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		code     ErrorCode
		expected bool
	}{
		{
			name:     "valid authentication code",
			code:     CodeAuthenticationFailed,
			expected: true,
		},
		{
			name:     "valid registry code",
			code:     CodeRegistryUnavailable,
			expected: true,
		},
		{
			name:     "valid network code",
			code:     CodeNetworkTimeout,
			expected: true,
		},
		{
			name:     "invalid code",
			code:     ErrorCode(9999),
			expected: false,
		},
		{
			name:     "zero code",
			code:     ErrorCode(0),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.code.IsValid()
			if got != tt.expected {
				t.Errorf("ErrorCode.IsValid() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestErrorCategory_String(t *testing.T) {
	tests := []struct {
		name     string
		category ErrorCategory
		expected string
	}{
		{
			name:     "transient category",
			category: CategoryTransient,
			expected: "transient",
		},
		{
			name:     "permanent category",
			category: CategoryPermanent,
			expected: "permanent",
		},
		{
			name:     "configuration category",
			category: CategoryConfiguration,
			expected: "configuration",
		},
		{
			name:     "authentication category",
			category: CategoryAuthentication,
			expected: "authentication",
		},
		{
			name:     "validation category",
			category: CategoryValidation,
			expected: "validation",
		},
		{
			name:     "network category",
			category: CategoryNetwork,
			expected: "network",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.category.String()
			if got != tt.expected {
				t.Errorf("ErrorCategory.String() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestErrorCategory_IsRetryable(t *testing.T) {
	tests := []struct {
		name     string
		category ErrorCategory
		expected bool
	}{
		{
			name:     "transient is retryable",
			category: CategoryTransient,
			expected: true,
		},
		{
			name:     "network is retryable",
			category: CategoryNetwork,
			expected: true,
		},
		{
			name:     "permanent is not retryable",
			category: CategoryPermanent,
			expected: false,
		},
		{
			name:     "authentication is not retryable",
			category: CategoryAuthentication,
			expected: false,
		},
		{
			name:     "validation is not retryable",
			category: CategoryValidation,
			expected: false,
		},
		{
			name:     "configuration is not retryable",
			category: CategoryConfiguration,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.category.IsRetryable()
			if got != tt.expected {
				t.Errorf("ErrorCategory.IsRetryable() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestErrorCodeRanges(t *testing.T) {
	tests := []struct {
		name      string
		code      ErrorCode
		minRange  int
		maxRange  int
		category  ErrorCategory
	}{
		{
			name:     "authentication range",
			code:     CodeAuthenticationFailed,
			minRange: 1000,
			maxRange: 1999,
			category: CategoryAuthentication,
		},
		{
			name:     "registry range",
			code:     CodeRegistryUnavailable,
			minRange: 2000,
			maxRange: 2999,
			category: CategoryTransient,
		},
		{
			name:     "manifest range",
			code:     CodeManifestNotFound,
			minRange: 3000,
			maxRange: 3999,
			category: CategoryPermanent,
		},
		{
			name:     "network range",
			code:     CodeNetworkTimeout,
			minRange: 4000,
			maxRange: 4999,
			category: CategoryNetwork,
		},
		{
			name:     "validation range",
			code:     CodeInvalidInput,
			minRange: 5000,
			maxRange: 5999,
			category: CategoryValidation,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			codeInt := int(tt.code)
			if codeInt < tt.minRange || codeInt > tt.maxRange {
				t.Errorf("code %d is outside expected range %d-%d", codeInt, tt.minRange, tt.maxRange)
			}

			if got := tt.code.GetCategory(); got != tt.category {
				t.Errorf("code category = %v, want %v", got, tt.category)
			}
		})
	}
}

func TestAllErrorCodesMapped(t *testing.T) {
	// Define all expected error codes
	expectedCodes := []ErrorCode{
		// Authentication Errors
		CodeAuthenticationFailed,
		CodeInvalidCredentials,
		CodeTokenExpired,
		CodeInvalidToken,
		CodeAuthorizationDenied,
		CodeCertificateInvalid,
		CodeTLSHandshakeFailed,

		// Registry Errors
		CodeRegistryUnavailable,
		CodeRegistryTimeout,
		CodeRegistryNotFound,
		CodeRegistryUnauthorized,
		CodeRegistryServerError,
		CodeRegistryRateLimit,
		CodeRegistryMaintenance,

		// Manifest Errors
		CodeManifestNotFound,
		CodeManifestInvalid,
		CodeManifestUnsupported,
		CodeManifestCorrupted,
		CodeDigestMismatch,
		CodeManifestTooLarge,

		// Network Errors
		CodeNetworkTimeout,
		CodeConnectionRefused,
		CodeNetworkUnreachable,
		CodeDNSResolutionFailed,
		CodeProxyError,
		CodeSSLError,

		// Validation Errors
		CodeInvalidInput,
		CodeMissingRequiredField,
		CodeInvalidFormat,
		CodeValueOutOfRange,
		CodeConflictingValues,
		CodeInvalidConfiguration,
	}

	// Verify all codes have category mappings
	for _, code := range expectedCodes {
		if !code.IsValid() {
			t.Errorf("error code %d is not mapped to a category", code)
		}
	}

	// Verify category mappings exist for expected number of codes
	if len(codeCategories) < len(expectedCodes) {
		t.Errorf("codeCategories has %d mappings, expected at least %d", len(codeCategories), len(expectedCodes))
	}
}