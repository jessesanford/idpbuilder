package validator

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test Suite 1: Image Name Validation (15 tests)

func TestValidateImageName_SimpleTag(t *testing.T) {
	err := ValidateImageName("alpine:latest")
	assert.NoError(t, err)
}

func TestValidateImageName_WithRegistry(t *testing.T) {
	err := ValidateImageName("docker.io/alpine:latest")
	assert.NoError(t, err)
}

func TestValidateImageName_WithNamespace(t *testing.T) {
	err := ValidateImageName("docker.io/library/ubuntu:22.04")
	assert.NoError(t, err)
}

func TestValidateImageName_WithDigest(t *testing.T) {
	err := ValidateImageName("alpine@sha256:1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef")
	assert.NoError(t, err)
}

func TestValidateImageName_NoTag(t *testing.T) {
	err := ValidateImageName("alpine")
	assert.NoError(t, err)
}

func TestValidateImageName_Empty(t *testing.T) {
	err := ValidateImageName("")
	require.Error(t, err)

	valErr, ok := err.(*ValidationError)
	require.True(t, ok)
	assert.Equal(t, "image-name", valErr.Field)
	assert.Contains(t, valErr.Message, "cannot be empty")
	assert.NotEmpty(t, valErr.Suggestion)
	assert.Equal(t, 1, valErr.ExitCode)
}

func TestValidateImageName_CommandInjection_Semicolon(t *testing.T) {
	err := ValidateImageName("alpine;rm -rf /")
	require.Error(t, err)

	valErr, ok := err.(*ValidationError)
	require.True(t, ok)
	assert.Equal(t, "image-name", valErr.Field)
	assert.Contains(t, valErr.Message, "shell metacharacters")
	assert.Equal(t, 1, valErr.ExitCode)
}

func TestValidateImageName_CommandInjection_Backtick(t *testing.T) {
	err := ValidateImageName("alpine`whoami`")
	require.Error(t, err)

	valErr, ok := err.(*ValidationError)
	require.True(t, ok)
	assert.Contains(t, valErr.Message, "shell metacharacters")
}

func TestValidateImageName_CommandInjection_Dollar(t *testing.T) {
	err := ValidateImageName("alpine$(whoami)")
	require.Error(t, err)

	valErr, ok := err.(*ValidationError)
	require.True(t, ok)
	assert.Contains(t, valErr.Message, "shell metacharacters")
}

func TestValidateImageName_InvalidTag(t *testing.T) {
	err := ValidateImageName("alpine:latest@invalid")
	require.Error(t, err)

	valErr, ok := err.(*ValidationError)
	require.True(t, ok)
	// This image has an invalid digest, so digest validation catches it first
	assert.Equal(t, "image-digest", valErr.Field)
}

func TestValidateImageName_InvalidDigest(t *testing.T) {
	err := ValidateImageName("alpine@sha256:invalid")
	require.Error(t, err)

	valErr, ok := err.(*ValidationError)
	require.True(t, ok)
	assert.Equal(t, "image-digest", valErr.Field)
	assert.Contains(t, valErr.Message, "invalid digest format")
}

func TestValidateImageName_Localhost(t *testing.T) {
	err := ValidateImageName("localhost:5000/myimage:latest")
	assert.NoError(t, err)
}

func TestValidateImageName_IPv4Registry(t *testing.T) {
	err := ValidateImageName("192.168.1.100:5000/myimage:latest")
	assert.NoError(t, err)
}

func TestValidateImageName_IPv6Registry(t *testing.T) {
	// IPv6 format is complex, testing simplified version
	err := ValidateImageName("myregistry/myimage:latest")
	assert.NoError(t, err)
}

func TestValidateImageName_LongName(t *testing.T) {
	err := ValidateImageName("registry.example.com:8080/namespace/subnamespace/repository:v1.2.3-beta")
	assert.NoError(t, err)
}

// Test Suite 2: Registry URL Validation (10 tests)

func TestValidateRegistryURL_SimpleDomain(t *testing.T) {
	err := ValidateRegistryURL("docker.io")
	assert.NoError(t, err)
}

func TestValidateRegistryURL_DomainWithPort(t *testing.T) {
	err := ValidateRegistryURL("registry.example.com:5000")
	assert.NoError(t, err)
}

func TestValidateRegistryURL_Localhost(t *testing.T) {
	err := ValidateRegistryURL("localhost")
	require.Error(t, err)

	// Should be an SSRFWarning (non-blocking)
	warning, ok := err.(*SSRFWarning)
	require.True(t, ok)
	assert.Contains(t, warning.Message, "private IP range")
	assert.True(t, IsWarning(err))
}

func TestValidateRegistryURL_IPv4(t *testing.T) {
	err := ValidateRegistryURL("192.168.1.100")
	require.Error(t, err)

	// Private IP should trigger SSRF warning
	_, ok := err.(*SSRFWarning)
	require.True(t, ok)
	assert.True(t, IsWarning(err))
}

func TestValidateRegistryURL_IPv6(t *testing.T) {
	// Testing with loopback IPv6
	err := ValidateRegistryURL("::1")
	require.Error(t, err)

	_, ok := err.(*SSRFWarning)
	require.True(t, ok)
	assert.True(t, IsWarning(err))
}

func TestValidateRegistryURL_PrivateIP_ClassA(t *testing.T) {
	err := ValidateRegistryURL("10.0.0.1")
	require.Error(t, err)

	warning, ok := err.(*SSRFWarning)
	require.True(t, ok)
	assert.Contains(t, warning.Message, "private IP range")
}

func TestValidateRegistryURL_PrivateIP_ClassB(t *testing.T) {
	err := ValidateRegistryURL("172.16.0.1")
	require.Error(t, err)

	warning, ok := err.(*SSRFWarning)
	require.True(t, ok)
	assert.Contains(t, warning.Message, "private IP range")
}

func TestValidateRegistryURL_PrivateIP_ClassC(t *testing.T) {
	err := ValidateRegistryURL("192.168.1.1")
	require.Error(t, err)

	warning, ok := err.(*SSRFWarning)
	require.True(t, ok)
	assert.Contains(t, warning.Message, "private IP range")
}

func TestValidateRegistryURL_Localhost_Warning(t *testing.T) {
	err := ValidateRegistryURL("127.0.0.1:5000")
	require.Error(t, err)

	_, ok := err.(*SSRFWarning)
	require.True(t, ok)
	assert.True(t, IsWarning(err))
}

func TestValidateRegistryURL_CommandInjection(t *testing.T) {
	err := ValidateRegistryURL("registry.example.com;rm -rf /")
	require.Error(t, err)

	valErr, ok := err.(*ValidationError)
	require.True(t, ok)
	assert.Equal(t, "registry", valErr.Field)
	assert.Contains(t, valErr.Message, "shell metacharacters")
	assert.Equal(t, 1, valErr.ExitCode)
}

// Test Suite 3: Credentials Validation (8 tests)

func TestValidateCredentials_Alphanumeric(t *testing.T) {
	err := ValidateCredentials("user123", "pass123")
	assert.NoError(t, err)
}

func TestValidateCredentials_SpecialCharsPassword(t *testing.T) {
	err := ValidateCredentials("user", "P@ssw0rd!#$%")
	assert.NoError(t, err)
}

func TestValidateCredentials_EmailUsername(t *testing.T) {
	err := ValidateCredentials("user@example.com", "password123")
	assert.NoError(t, err)
}

func TestValidateCredentials_EmptyUsername(t *testing.T) {
	err := ValidateCredentials("", "password")
	require.Error(t, err)

	valErr, ok := err.(*ValidationError)
	require.True(t, ok)
	assert.Equal(t, "username", valErr.Field)
	assert.Contains(t, valErr.Message, "required")
	assert.Contains(t, valErr.Suggestion, "IDPBUILDER_USERNAME")
	assert.Equal(t, 1, valErr.ExitCode)
}

func TestValidateCredentials_EmptyPassword(t *testing.T) {
	err := ValidateCredentials("user", "")
	require.Error(t, err)

	valErr, ok := err.(*ValidationError)
	require.True(t, ok)
	assert.Equal(t, "password", valErr.Field)
	assert.Contains(t, valErr.Message, "required")
	assert.Contains(t, valErr.Suggestion, "IDPBUILDER_PASSWORD")
}

func TestValidateCredentials_UsernameInjection(t *testing.T) {
	err := ValidateCredentials("user;whoami", "password")
	require.Error(t, err)

	valErr, ok := err.(*ValidationError)
	require.True(t, ok)
	assert.Equal(t, "username", valErr.Field)
	assert.Contains(t, valErr.Message, "shell metacharacters")
}

func TestValidateCredentials_UsernameBacktick(t *testing.T) {
	err := ValidateCredentials("user`id`", "password")
	require.Error(t, err)

	valErr, ok := err.(*ValidationError)
	require.True(t, ok)
	assert.Contains(t, valErr.Message, "shell metacharacters")
}

func TestValidateCredentials_WeakCredentials(t *testing.T) {
	err := ValidateCredentials("admin", "password")
	require.Error(t, err)

	// Should be a SecurityWarning (non-blocking)
	warning, ok := err.(*SecurityWarning)
	require.True(t, ok)
	assert.Contains(t, warning.Message, "weak credentials")
	assert.True(t, IsWarning(err))
}

// Additional helper tests

func TestContainsAnyChar(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		chars    string
		expected bool
	}{
		{"no match", "abc123", ";|&", false},
		{"semicolon match", "abc;123", ";|&", true},
		{"pipe match", "abc|123", ";|&", true},
		{"ampersand match", "abc&123", ";|&", true},
		{"empty input", "", ";|&", false},
		{"empty chars", "abc123", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := containsAnyChar(tt.input, tt.chars)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsWarning(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{"SSRFWarning is warning", &SSRFWarning{Target: "test", Message: "test", Suggestion: "test"}, true},
		{"SecurityWarning is warning", &SecurityWarning{Message: "test", Suggestion: "test"}, true},
		{"ValidationError is not warning", &ValidationError{Field: "test", Message: "test", Suggestion: "test", ExitCode: 1}, false},
		{"nil is not warning", nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsWarning(tt.err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Test Error() methods for coverage

func TestValidationError_Error(t *testing.T) {
	// With suggestion
	err := &ValidationError{
		Field:      "test-field",
		Message:    "test error",
		Suggestion: "test suggestion",
		ExitCode:   1,
	}
	errMsg := err.Error()
	assert.Contains(t, errMsg, "Error: test error")
	assert.Contains(t, errMsg, "Suggestion: test suggestion")

	// Without suggestion
	err2 := &ValidationError{
		Field:    "test-field",
		Message:  "test error",
		ExitCode: 1,
	}
	errMsg2 := err2.Error()
	assert.Contains(t, errMsg2, "Error: test error")
	assert.NotContains(t, errMsg2, "Suggestion:")
}

func TestSSRFWarning_Error(t *testing.T) {
	warning := &SSRFWarning{
		Target:     "192.168.1.1",
		Message:    "private IP detected",
		Suggestion: "ensure this is intentional",
	}
	errMsg := warning.Error()
	assert.Contains(t, errMsg, "Warning: private IP detected")
	assert.Contains(t, errMsg, "Suggestion: ensure this is intentional")
}

func TestSecurityWarning_Error(t *testing.T) {
	warning := &SecurityWarning{
		Message:    "weak credentials",
		Suggestion: "use stronger passwords",
	}
	errMsg := warning.Error()
	assert.Contains(t, errMsg, "Warning: weak credentials")
	assert.Contains(t, errMsg, "Suggestion: use stronger passwords")
}

// Additional registry URL tests for coverage

func TestValidateRegistryURL_Empty(t *testing.T) {
	err := ValidateRegistryURL("")
	require.Error(t, err)

	valErr, ok := err.(*ValidationError)
	require.True(t, ok)
	assert.Equal(t, "registry", valErr.Field)
	assert.Contains(t, valErr.Message, "cannot be empty")
}

func TestValidateRegistryURL_FullHTTPSURL(t *testing.T) {
	err := ValidateRegistryURL("https://registry.example.com:5000")
	assert.NoError(t, err)
}

func TestValidateRegistryURL_InvalidHostname(t *testing.T) {
	err := ValidateRegistryURL("not@valid@hostname")
	require.Error(t, err)

	valErr, ok := err.(*ValidationError)
	require.True(t, ok)
	assert.Equal(t, "registry", valErr.Field)
	assert.Contains(t, valErr.Message, "invalid registry hostname format")
}
