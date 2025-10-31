package auth

import (
	"testing"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TC-AUTH-IMPL-001: NewBasicAuthProvider Success
func TestNewBasicAuthProvider(t *testing.T) {
	// Given: Valid credentials
	username := "giteaadmin"
	password := "mypassword"

	// When: Creating provider
	provider := NewBasicAuthProvider(username, password)

	// Then: Provider created successfully
	assert.NotNil(t, provider)

	// Verify implements Provider interface
	var _ Provider = provider
}

// TC-AUTH-IMPL-002: GetAuthenticator Success with Valid Credentials
func TestGetAuthenticator_Success(t *testing.T) {
	// Given: Provider with valid credentials
	provider := NewBasicAuthProvider("giteaadmin", "password")

	// When: Getting authenticator
	authenticator, err := provider.GetAuthenticator()

	// Then: Returns authenticator successfully
	require.NoError(t, err)
	require.NotNil(t, authenticator)

	// Verify authenticator is correct type
	basicAuth, ok := authenticator.(*authn.Basic)
	require.True(t, ok, "authenticator should be *authn.Basic")
	assert.Equal(t, "giteaadmin", basicAuth.Username)
	assert.Equal(t, "password", basicAuth.Password)
}

// TC-AUTH-IMPL-003: GetAuthenticator Fails with Empty Username
func TestGetAuthenticator_EmptyUsername(t *testing.T) {
	// Given: Provider with empty username
	provider := NewBasicAuthProvider("", "password")

	// When: Getting authenticator
	authenticator, err := provider.GetAuthenticator()

	// Then: Returns error
	require.Error(t, err)
	assert.Nil(t, authenticator)

	// Verify error type
	var valErr *CredentialValidationError
	require.ErrorAs(t, err, &valErr)
	assert.Equal(t, "username", valErr.Field)
	assert.Contains(t, valErr.Reason, "cannot be empty")
}

// TC-AUTH-IMPL-004: ValidateCredentials Passes for Valid Credentials
func TestValidateCredentials_ValidCredentials(t *testing.T) {
	testCases := []struct {
		name     string
		username string
		password string
	}{
		{"simple", "user", "pass"},
		{"special_chars", "user", "P@ss!w0rd#123"},
		{"unicode", "user", "пароль密码🔒"},
		{"spaces", "user", "pass with spaces"},
		{"quotes_double", "user", "pass\"with\"quotes"},
		{"quotes_single", "user", "pass'with'quotes"},
		{"quotes_mixed", "user", "pass\"with'mixed"},
		{"complex_password", "user", "P@ss!w0rd#123 with \"quotes\" and 'apostrophes' пароль🔒"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Given: Provider with valid credentials
			provider := NewBasicAuthProvider(tc.username, tc.password)

			// When: Validating credentials
			err := provider.ValidateCredentials()

			// Then: Validation passes
			assert.NoError(t, err, "Password should be valid: %s", tc.password)
		})
	}
}

// TC-AUTH-IMPL-005: ValidateCredentials Fails for Empty Username
func TestValidateCredentials_EmptyUsername(t *testing.T) {
	// Given: Provider with empty username
	provider := NewBasicAuthProvider("", "password")

	// When: Validating credentials
	err := provider.ValidateCredentials()

	// Then: Validation fails
	require.Error(t, err)

	var valErr *CredentialValidationError
	require.ErrorAs(t, err, &valErr)
	assert.Equal(t, "username", valErr.Field)
	assert.Contains(t, valErr.Reason, "cannot be empty")
}

// TC-AUTH-IMPL-006: ValidateCredentials Fails for Empty Password
func TestValidateCredentials_EmptyPassword(t *testing.T) {
	// Given: Provider with empty password
	provider := NewBasicAuthProvider("username", "")

	// When: Validating credentials
	err := provider.ValidateCredentials()

	// Then: Validation fails
	require.Error(t, err)

	var valErr *CredentialValidationError
	require.ErrorAs(t, err, &valErr)
	assert.Equal(t, "password", valErr.Field)
	assert.Contains(t, valErr.Reason, "cannot be empty")
}

// TC-AUTH-IMPL-007: ValidateCredentials Fails for Control Characters in Username
func TestValidateCredentials_ControlCharactersInUsername(t *testing.T) {
	testCases := []struct {
		name     string
		username string
	}{
		{"newline", "user\n"},
		{"tab", "user\t"},
		{"null_byte", "user\x00"},
		{"escape", "user\x1b"},
		{"carriage_return", "user\r"},
		{"backspace", "user\x08"},
		{"delete", "user\x7f"},
		{"form_feed", "user\x0c"},
		{"vertical_tab", "user\x0b"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Given: Provider with control character in username
			provider := NewBasicAuthProvider(tc.username, "password")

			// When: Validating credentials
			err := provider.ValidateCredentials()

			// Then: Validation fails
			require.Error(t, err, "Username with control character should fail: %q", tc.username)

			var valErr *CredentialValidationError
			require.ErrorAs(t, err, &valErr)
			assert.Equal(t, "username", valErr.Field)
			assert.Contains(t, valErr.Reason, "control characters")
		})
	}
}

// Additional test: Control characters allowed in password
func TestValidateCredentials_ControlCharactersInPassword(t *testing.T) {
	testCases := []struct {
		name     string
		password string
	}{
		{"newline", "pass\n"},
		{"tab", "pass\t"},
		{"null_byte", "pass\x00"},
		{"escape", "pass\x1b"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Given: Provider with control character in password
			provider := NewBasicAuthProvider("user", tc.password)

			// When: Validating credentials
			err := provider.ValidateCredentials()

			// Then: Validation passes (passwords allow any characters)
			assert.NoError(t, err, "Password with control character should be valid: %q", tc.password)
		})
	}
}

// Additional test: GetAuthenticator with special character password
func TestGetAuthenticator_SpecialCharacterPassword(t *testing.T) {
	// Given: Provider with special character password
	provider := NewBasicAuthProvider("user", "P@ss!w0rd#123 with \"quotes\" 密码🔒")

	// When: Getting authenticator
	authenticator, err := provider.GetAuthenticator()

	// Then: Returns authenticator successfully
	require.NoError(t, err)
	require.NotNil(t, authenticator)

	// Verify password preserved exactly
	basicAuth, ok := authenticator.(*authn.Basic)
	require.True(t, ok)
	assert.Equal(t, "P@ss!w0rd#123 with \"quotes\" 密码🔒", basicAuth.Password)
}

// Additional test: Both username and password empty
func TestValidateCredentials_BothEmpty(t *testing.T) {
	// Given: Provider with both empty
	provider := NewBasicAuthProvider("", "")

	// When: Validating credentials
	err := provider.ValidateCredentials()

	// Then: Returns error for username (checked first)
	require.Error(t, err)

	var valErr *CredentialValidationError
	require.ErrorAs(t, err, &valErr)
	assert.Equal(t, "username", valErr.Field)
}

// Additional test: Whitespace-only username
func TestValidateCredentials_WhitespaceUsername(t *testing.T) {
	testCases := []struct {
		name     string
		username string
		valid    bool
	}{
		{"single_space", " ", true},        // Space is valid (ASCII 32)
		{"multiple_spaces", "   ", true},   // Spaces are valid
		{"space_in_middle", "user name", true}, // Spaces in middle are valid
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Given: Provider with whitespace username
			provider := NewBasicAuthProvider(tc.username, "password")

			// When: Validating credentials
			err := provider.ValidateCredentials()

			// Then: Validation result matches expectation
			if tc.valid {
				assert.NoError(t, err, "Username should be valid: %q", tc.username)
			} else {
				require.Error(t, err)
			}
		})
	}
}

// Test helper function: containsControlChars
func TestContainsControlChars(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected bool
	}{
		{"empty", "", false},
		{"normal", "username", false},
		{"with_space", "user name", false}, // Space (ASCII 32) is NOT control char
		{"with_newline", "user\n", true},
		{"with_tab", "user\t", true},
		{"with_null", "user\x00", true},
		{"with_escape", "user\x1b", true},
		{"with_delete", "user\x7f", true},
		{"ascii_31", "user\x1f", true}, // Last control char before space
		{"ascii_32", "user ", false},   // Space is NOT control char
		{"unicode", "пароль密码🔒", false},
		{"special_chars", "!@#$%^&*()", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := containsControlChars(tc.input)
			assert.Equal(t, tc.expected, result, "Input: %q", tc.input)
		})
	}
}
