package validator

import (
	"fmt"
	"strings"
)

// ValidateCredentials validates username and password for security
//
// Validation Rules:
// 1. Username cannot be empty
// 2. Password cannot be empty
// 3. Username must not contain shell metacharacters
// 4. Password CAN contain special characters (no restrictions)
// 5. Weak credential warning (non-blocking)
//
// Returns:
//   - nil if valid
//   - ValidationError if required fields missing or username has dangerous chars
//   - SecurityWarning if weak credentials detected (allows continuation)
func ValidateCredentials(username, password string) error {
	// Validate username
	if username == "" {
		return &ValidationError{
			Field:      "username",
			Message:    "username is required",
			Suggestion: "provide username via --username flag or IDPBUILDER_USERNAME environment variable",
			ExitCode:   1,
		}
	}

	// Check username for command injection patterns
	if containsAnyChar(username, dangerousChars) {
		return &ValidationError{
			Field:      "username",
			Message:    fmt.Sprintf("username contains shell metacharacters: %s", username),
			Suggestion: "use alphanumeric characters, hyphens, underscores, dots, and @ symbols only",
			ExitCode:   1,
		}
	}

	// Validate password
	if password == "" {
		return &ValidationError{
			Field:      "password",
			Message:    "password is required",
			Suggestion: "provide password via --password flag or IDPBUILDER_PASSWORD environment variable",
			ExitCode:   1,
		}
	}

	// Note: We do NOT validate password characters (passwords can contain special chars)
	// But we DO escape them properly when used in authentication

	// Warn if credentials appear in plain sight (basic security check)
	if strings.Contains(username, "admin") && strings.Contains(password, "password") {
		return &SecurityWarning{
			Message:    "using default or weak credentials detected",
			Suggestion: "consider using stronger credentials for production registries",
		}
	}

	return nil
}
