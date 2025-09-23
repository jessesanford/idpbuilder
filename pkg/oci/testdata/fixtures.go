package testdata

import (
	"fmt"
	"os"
)

// Test registry URLs for testing different scenarios
const (
	TestRegistryHTTPS    = "https://my-registry.example.com"
	TestRegistryHTTP     = "http://insecure-registry.example.com"
	TestRegistryWithPort = "registry.example.com:5000"
	TestRegistryNoScheme = "registry.example.com"
	TestGiteaRegistry    = "gitea.cnoe.localtest.me:8443"
)

// Valid test credentials
const (
	ValidUsername = "testuser"
	ValidPassword = "testpass123"
	ValidToken    = "ghp_1234567890abcdef1234567890abcdef12345678"
)

// Invalid test credentials for error testing
const (
	EmptyUsername = ""
	EmptyPassword = ""
	InvalidToken  = "invalid-token-format"
)

// MockCredentials represents test credential data
type MockCredentials struct {
	Username string
	Password string
	Token    string
	Source   string
}

// ValidMockCredentials returns a set of valid test credentials
func ValidMockCredentials() MockCredentials {
	return MockCredentials{
		Username: ValidUsername,
		Password: ValidPassword,
		Token:    ValidToken,
		Source:   "test",
	}
}

// InvalidMockCredentials returns invalid credentials for error testing
func InvalidMockCredentials() MockCredentials {
	return MockCredentials{
		Username: EmptyUsername,
		Password: EmptyPassword,
		Token:    InvalidToken,
		Source:   "test",
	}
}

// SetupTestEnv sets up environment variables for testing
func SetupTestEnv() {
	os.Setenv("GITEA_USERNAME", ValidUsername)
	os.Setenv("GITEA_PASSWORD", ValidPassword)
}

// CleanupTestEnv removes test environment variables
func CleanupTestEnv() {
	os.Unsetenv("GITEA_USERNAME")
	os.Unsetenv("GITEA_PASSWORD")
}

// MockSecretData simulates Kubernetes secret data
func MockSecretData() map[string][]byte {
	return map[string][]byte{
		"username": []byte(ValidUsername),
		"password": []byte(ValidPassword),
	}
}

// EmptySecretData simulates missing or empty secret data
func EmptySecretData() map[string][]byte {
	return map[string][]byte{}
}

// ExpectedErrorMessage returns expected error messages for testing
func ExpectedErrorMessage(errorType string) string {
	switch errorType {
	case "empty_username":
		return "username cannot be empty"
	case "empty_password":
		return "password cannot be empty"
	case "missing_secret":
		return "credential secret not found"
	case "auth_failed":
		return "authentication failed"
	default:
		return fmt.Sprintf("unknown error type: %s", errorType)
	}
}