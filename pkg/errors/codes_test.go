package errors

import (
	"testing"
)

func TestErrorCode_String(t *testing.T) {
	tests := []struct {
		name     string
		code     ErrorCode
		expected string
	}{
		{
			name: "network timeout",
			code: ErrorNetworkTimeout,
			expected: "transient.network_timeout",
		},
		{
			name: "not found",
			code: ErrorNotFound,
			expected: "permanent.not_found",
		},
		{
			name: "invalid config",
			code: ErrorInvalidConfig,
			expected: "configuration.invalid_config",
		},
		{
			name: "invalid input",
			code: ErrorInvalidInput,
			expected: "validation.invalid_input",
		},
		{
			name: "unauthenticated",
			code: ErrorUnauthenticated,
			expected: "authentication.unauthenticated",
		},
		{
			name: "unauthorized",
			code: ErrorUnauthorized,
			expected: "authorization.unauthorized",
		},
		{
			name: "resource exhausted",
			code: ErrorResourceExhausted,
			expected: "resource.resource_exhausted",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.code.String()
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestTransientErrors(t *testing.T) {
	transientErrors := []ErrorCode{
		ErrorNetworkTimeout,
		ErrorNetworkUnavailable,
		ErrorRateLimited,
		ErrorServiceBusy,
	}

	for _, code := range transientErrors {
		t.Run("transient_"+code.Code, func(t *testing.T) {
			if code.Category != CategoryTransient {
				t.Errorf("Expected category %q, got %q", CategoryTransient, code.Category)
			}

			if !IsTransient(code) {
				t.Error("Expected error to be identified as transient")
			}

			if IsPermanent(code) {
				t.Error("Expected error to not be identified as permanent")
			}
		})
	}
}

func TestPermanentErrors(t *testing.T) {
	permanentErrors := []ErrorCode{
		ErrorNotFound,
		ErrorAlreadyExists,
		ErrorUnsupported,
		ErrorIncompatible,
	}

	for _, code := range permanentErrors {
		t.Run("permanent_"+code.Code, func(t *testing.T) {
			if code.Category != CategoryPermanent {
				t.Errorf("Expected category %q, got %q", CategoryPermanent, code.Category)
			}

			if !IsPermanent(code) {
				t.Error("Expected error to be identified as permanent")
			}

			if IsTransient(code) {
				t.Error("Expected error to not be identified as transient")
			}
		})
	}
}

func TestConfigurationErrors(t *testing.T) {
	configurationErrors := []ErrorCode{
		ErrorInvalidConfig,
		ErrorMissingConfig,
		ErrorConfigFormat,
	}

	for _, code := range configurationErrors {
		t.Run("configuration_"+code.Code, func(t *testing.T) {
			if code.Category != CategoryConfiguration {
				t.Errorf("Expected category %q, got %q", CategoryConfiguration, code.Category)
			}

			if !IsConfigurationError(code) {
				t.Error("Expected error to be identified as configuration error")
			}

			if code.Severity != SeverityHigh {
				t.Errorf("Expected severity %q, got %q", SeverityHigh, code.Severity)
			}
		})
	}
}

func TestValidationErrors(t *testing.T) {
	validationErrors := []ErrorCode{
		ErrorInvalidInput,
		ErrorMissingRequired,
		ErrorInvalidFormat,
	}

	for _, code := range validationErrors {
		t.Run("validation_"+code.Code, func(t *testing.T) {
			if code.Category != CategoryValidation {
				t.Errorf("Expected category %q, got %q", CategoryValidation, code.Category)
			}

			if !IsValidationError(code) {
				t.Error("Expected error to be identified as validation error")
			}

			if code.Severity != SeverityMedium {
				t.Errorf("Expected severity %q, got %q", SeverityMedium, code.Severity)
			}
		})
	}
}

func TestAuthenticationErrors(t *testing.T) {
	authenticationErrors := []ErrorCode{
		ErrorUnauthenticated,
		ErrorInvalidCredentials,
		ErrorTokenExpired,
	}

	for _, code := range authenticationErrors {
		t.Run("authentication_"+code.Code, func(t *testing.T) {
			if code.Category != CategoryAuthentication {
				t.Errorf("Expected category %q, got %q", CategoryAuthentication, code.Category)
			}

			if !IsAuthenticationError(code) {
				t.Error("Expected error to be identified as authentication error")
			}
		})
	}
}

func TestAuthorizationErrors(t *testing.T) {
	authorizationErrors := []ErrorCode{
		ErrorUnauthorized,
		ErrorInsufficientPermissions,
	}

	for _, code := range authorizationErrors {
		t.Run("authorization_"+code.Code, func(t *testing.T) {
			if code.Category != CategoryAuthorization {
				t.Errorf("Expected category %q, got %q", CategoryAuthorization, code.Category)
			}

			if !IsAuthorizationError(code) {
				t.Error("Expected error to be identified as authorization error")
			}

			if code.Severity != SeverityHigh {
				t.Errorf("Expected severity %q, got %q", SeverityHigh, code.Severity)
			}
		})
	}
}

func TestResourceErrors(t *testing.T) {
	resourceErrors := []ErrorCode{
		ErrorResourceExhausted,
		ErrorResourceLocked,
		ErrorResourceCorrupted,
	}

	for _, code := range resourceErrors {
		t.Run("resource_"+code.Code, func(t *testing.T) {
			if code.Category != CategoryResource {
				t.Errorf("Expected category %q, got %q", CategoryResource, code.Category)
			}

			if !IsResourceError(code) {
				t.Error("Expected error to be identified as resource error")
			}
		})
	}
}

func TestErrorSeverityLevels(t *testing.T) {
	severityTests := []struct {
		name     string
		code     ErrorCode
		severity string
		critical bool
	}{
		{
			name:     "low severity - rate limited",
			code:     ErrorRateLimited,
			severity: SeverityLow,
			critical: false,
		},
		{
			name:     "medium severity - network timeout",
			code:     ErrorNetworkTimeout,
			severity: SeverityMedium,
			critical: false,
		},
		{
			name:     "high severity - invalid config",
			code:     ErrorInvalidConfig,
			severity: SeverityHigh,
			critical: false,
		},
		{
			name:     "critical severity - resource corrupted",
			code:     ErrorResourceCorrupted,
			severity: SeverityCritical,
			critical: true,
		},
	}

	for _, tt := range severityTests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.code.Severity != tt.severity {
				t.Errorf("Expected severity %q, got %q", tt.severity, tt.code.Severity)
			}

			if IsCritical(tt.code) != tt.critical {
				t.Errorf("Expected IsCritical to return %v", tt.critical)
			}
		})
	}
}

func TestErrorCodeMessages(t *testing.T) {
	messageTests := []struct {
		name    string
		code    ErrorCode
		message string
	}{
		{
			name:    "network timeout message",
			code:    ErrorNetworkTimeout,
			message: "Network operation timed out",
		},
		{
			name:    "not found message",
			code:    ErrorNotFound,
			message: "Resource not found",
		},
		{
			name:    "invalid credentials message",
			code:    ErrorInvalidCredentials,
			message: "Invalid credentials provided",
		},
		{
			name:    "resource corrupted message",
			code:    ErrorResourceCorrupted,
			message: "Resource data is corrupted",
		},
	}

	for _, tt := range messageTests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.code.Message != tt.message {
				t.Errorf("Expected message %q, got %q", tt.message, tt.code.Message)
			}
		})
	}
}

func TestHelperFunctions(t *testing.T) {
	t.Run("IsTransient", func(t *testing.T) {
		if !IsTransient(ErrorNetworkTimeout) {
			t.Error("Expected network timeout to be transient")
		}
		
		if IsTransient(ErrorNotFound) {
			t.Error("Expected not found to not be transient")
		}
	})

	t.Run("IsPermanent", func(t *testing.T) {
		if !IsPermanent(ErrorNotFound) {
			t.Error("Expected not found to be permanent")
		}
		
		if IsPermanent(ErrorNetworkTimeout) {
			t.Error("Expected network timeout to not be permanent")
		}
	})

	t.Run("IsConfigurationError", func(t *testing.T) {
		if !IsConfigurationError(ErrorInvalidConfig) {
			t.Error("Expected invalid config to be configuration error")
		}
		
		if IsConfigurationError(ErrorNotFound) {
			t.Error("Expected not found to not be configuration error")
		}
	})

	t.Run("IsValidationError", func(t *testing.T) {
		if !IsValidationError(ErrorInvalidInput) {
			t.Error("Expected invalid input to be validation error")
		}
		
		if IsValidationError(ErrorNotFound) {
			t.Error("Expected not found to not be validation error")
		}
	})

	t.Run("IsAuthenticationError", func(t *testing.T) {
		if !IsAuthenticationError(ErrorUnauthenticated) {
			t.Error("Expected unauthenticated to be authentication error")
		}
		
		if IsAuthenticationError(ErrorNotFound) {
			t.Error("Expected not found to not be authentication error")
		}
	})

	t.Run("IsAuthorizationError", func(t *testing.T) {
		if !IsAuthorizationError(ErrorUnauthorized) {
			t.Error("Expected unauthorized to be authorization error")
		}
		
		if IsAuthorizationError(ErrorNotFound) {
			t.Error("Expected not found to not be authorization error")
		}
	})

	t.Run("IsResourceError", func(t *testing.T) {
		if !IsResourceError(ErrorResourceExhausted) {
			t.Error("Expected resource exhausted to be resource error")
		}
		
		if IsResourceError(ErrorNotFound) {
			t.Error("Expected not found to not be resource error")
		}
	})

	t.Run("IsCritical", func(t *testing.T) {
		if !IsCritical(ErrorResourceCorrupted) {
			t.Error("Expected resource corrupted to be critical")
		}
		
		if IsCritical(ErrorRateLimited) {
			t.Error("Expected rate limited to not be critical")
		}
	})
}

func TestAllErrorCodesHaveRequiredFields(t *testing.T) {
	allCodes := []ErrorCode{
		// Transient errors
		ErrorNetworkTimeout,
		ErrorNetworkUnavailable,
		ErrorRateLimited,
		ErrorServiceBusy,
		// Permanent errors
		ErrorNotFound,
		ErrorAlreadyExists,
		ErrorUnsupported,
		ErrorIncompatible,
		// Configuration errors
		ErrorInvalidConfig,
		ErrorMissingConfig,
		ErrorConfigFormat,
		// Validation errors
		ErrorInvalidInput,
		ErrorMissingRequired,
		ErrorInvalidFormat,
		// Authentication errors
		ErrorUnauthenticated,
		ErrorInvalidCredentials,
		ErrorTokenExpired,
		// Authorization errors
		ErrorUnauthorized,
		ErrorInsufficientPermissions,
		// Resource errors
		ErrorResourceExhausted,
		ErrorResourceLocked,
		ErrorResourceCorrupted,
	}

	for _, code := range allCodes {
		t.Run("code_"+code.Code, func(t *testing.T) {
			if code.Category == "" {
				t.Error("Error code missing category")
			}
			
			if code.Code == "" {
				t.Error("Error code missing code")
			}
			
			if code.Message == "" {
				t.Error("Error code missing message")
			}
			
			if code.Severity == "" {
				t.Error("Error code missing severity")
			}

			// Verify category is valid
			validCategories := []string{
				CategoryTransient,
				CategoryPermanent,
				CategoryConfiguration,
				CategoryValidation,
				CategoryNetwork,
				CategoryAuthentication,
				CategoryAuthorization,
				CategoryResource,
			}
			
			validCategory := false
			for _, valid := range validCategories {
				if code.Category == valid {
					validCategory = true
					break
				}
			}
			
			if !validCategory {
				t.Errorf("Invalid category %q", code.Category)
			}

			// Verify severity is valid
			validSeverities := []string{
				SeverityLow,
				SeverityMedium,
				SeverityHigh,
				SeverityCritical,
			}
			
			validSeverity := false
			for _, valid := range validSeverities {
				if code.Severity == valid {
					validSeverity = true
					break
				}
			}
			
			if !validSeverity {
				t.Errorf("Invalid severity %q", code.Severity)
			}
		})
	}
}

func TestErrorCodeConsistency(t *testing.T) {
	t.Run("transient errors have appropriate severity", func(t *testing.T) {
		transientCodes := []ErrorCode{
			ErrorNetworkTimeout,
			ErrorNetworkUnavailable,
			ErrorRateLimited,
			ErrorServiceBusy,
		}

		for _, code := range transientCodes {
			// Transient errors should generally be low to medium severity
			if code.Severity == SeverityCritical {
				t.Errorf("Transient error %s should not be critical severity", code.Code)
			}
		}
	})

	t.Run("configuration errors are high severity", func(t *testing.T) {
		configCodes := []ErrorCode{
			ErrorInvalidConfig,
			ErrorMissingConfig,
			ErrorConfigFormat,
		}

		for _, code := range configCodes {
			if code.Severity != SeverityHigh {
				t.Errorf("Configuration error %s should be high severity, got %s", code.Code, code.Severity)
			}
		}
	})

	t.Run("resource corrupted is critical", func(t *testing.T) {
		if ErrorResourceCorrupted.Severity != SeverityCritical {
			t.Errorf("Resource corrupted should be critical severity, got %s", ErrorResourceCorrupted.Severity)
		}
	})
}

func TestCategoryGroupings(t *testing.T) {
	categoryTests := map[string][]ErrorCode{
		CategoryTransient: {
			ErrorNetworkTimeout,
			ErrorNetworkUnavailable,
			ErrorRateLimited,
			ErrorServiceBusy,
		},
		CategoryPermanent: {
			ErrorNotFound,
			ErrorAlreadyExists,
			ErrorUnsupported,
			ErrorIncompatible,
		},
		CategoryConfiguration: {
			ErrorInvalidConfig,
			ErrorMissingConfig,
			ErrorConfigFormat,
		},
		CategoryValidation: {
			ErrorInvalidInput,
			ErrorMissingRequired,
			ErrorInvalidFormat,
		},
		CategoryAuthentication: {
			ErrorUnauthenticated,
			ErrorInvalidCredentials,
			ErrorTokenExpired,
		},
		CategoryAuthorization: {
			ErrorUnauthorized,
			ErrorInsufficientPermissions,
		},
		CategoryResource: {
			ErrorResourceExhausted,
			ErrorResourceLocked,
			ErrorResourceCorrupted,
		},
	}

	for category, codes := range categoryTests {
		t.Run("category_"+category, func(t *testing.T) {
			for _, code := range codes {
				if code.Category != category {
					t.Errorf("Code %s has category %s, expected %s", code.Code, code.Category, category)
				}
			}
		})
	}
}