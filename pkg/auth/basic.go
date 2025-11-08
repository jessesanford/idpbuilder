package auth

import (
	"github.com/google/go-containerregistry/pkg/authn"
)

// basicAuthProvider implements the Provider interface using basic username/password authentication.
//
// This implementation uses HTTP Basic Authentication to authenticate with OCI registries.
// Credentials are transmitted via the Authorization header as base64-encoded username:password.
type basicAuthProvider struct {
	credentials Credentials
}

// NewBasicAuthProvider creates a basic authentication provider.
//
// Basic authentication uses username and password credentials transmitted
// via HTTP Basic Auth header to the registry. The credentials are base64-encoded
// and sent with each request in the Authorization header.
//
// This constructor does not validate credentials. Validation is performed
// when GetAuthenticator() or ValidateCredentials() is called.
//
// Parameters:
//   - username: Registry username (typically "giteaadmin" for Gitea)
//   - password: Registry password (supports all special characters including unicode, quotes, spaces)
//
// Returns:
//   - Provider: Authentication provider interface implementation
//
// Example:
//
//	provider := auth.NewBasicAuthProvider("giteaadmin", "myP@ssw0rd!")
//	if err := provider.ValidateCredentials(); err != nil {
//	    return fmt.Errorf("invalid credentials: %w", err)
//	}
//
// Security Considerations:
//   - Passwords can contain any characters (unicode, quotes, spaces, control chars)
//   - Usernames are validated to prevent control character injection
//   - Credentials are NOT logged or exposed in error messages
//   - Use HTTPS to prevent credential interception
func NewBasicAuthProvider(username, password string) Provider {
	return &basicAuthProvider{
		credentials: Credentials{
			Username: username,
			Password: password,
		},
	}
}

// GetAuthenticator returns an authn.Authenticator compatible with go-containerregistry.
//
// This method converts internal credentials to the authn.Basic format expected
// by go-containerregistry's remote.Push(), remote.Pull(), and other registry operations.
//
// The authenticator is created only after successful credential validation.
// If validation fails, an error is returned without creating the authenticator.
//
// Returns:
//   - authn.Authenticator: Authenticator instance for go-containerregistry
//   - error: CredentialValidationError if credentials are malformed or invalid
//
// Example:
//
//	provider := auth.NewBasicAuthProvider("giteaadmin", "password")
//	authenticator, err := provider.GetAuthenticator()
//	if err != nil {
//	    return fmt.Errorf("failed to get authenticator: %w", err)
//	}
//	// Use with go-containerregistry:
//	// remote.Push(ref, image, remote.WithAuth(authenticator))
//
// Error Handling:
//   - Returns CredentialValidationError if username is empty
//   - Returns CredentialValidationError if password is empty
//   - Returns CredentialValidationError if username contains control characters
func (p *basicAuthProvider) GetAuthenticator() (authn.Authenticator, error) {
	// Validate credentials before creating authenticator
	if err := p.ValidateCredentials(); err != nil {
		return nil, err
	}

	// Create go-containerregistry Basic authenticator
	// This type implements authn.Authenticator interface and handles
	// base64 encoding and Authorization header construction
	authenticator := &authn.Basic{
		Username: p.credentials.Username,
		Password: p.credentials.Password,
	}

	return authenticator, nil
}

// ValidateCredentials performs pre-flight validation of credentials.
//
// This method checks credential format and basic security requirements
// without contacting the registry. It ensures credentials meet minimum
// requirements before attempting authentication.
//
// Validation Rules:
//
// Username Requirements:
//   - MUST NOT be empty
//   - MUST NOT contain control characters (ASCII 0-31 or 127)
//
// Password Requirements:
//   - MUST NOT be empty
//   - MAY contain ANY characters (including unicode, quotes, spaces, control chars)
//
// The asymmetric validation (strict username, permissive password) prevents
// terminal escape sequence injection attacks while allowing maximum password flexibility.
//
// Returns:
//   - error: CredentialValidationError with details if invalid, nil if valid
//
// Example:
//
//	provider := auth.NewBasicAuthProvider("admin", "P@ss!w0rd#123")
//	if err := provider.ValidateCredentials(); err != nil {
//	    var valErr *CredentialValidationError
//	    if errors.As(err, &valErr) {
//	        log.Printf("Validation failed for %s: %s", valErr.Field, valErr.Reason)
//	    }
//	    return err
//	}
//
// Security Notes:
//   - Username validation prevents control character injection attacks
//   - Password validation only checks for empty (no character restrictions)
//   - This does NOT validate credentials with the registry
//   - Credentials are NOT logged or included in error messages
func (p *basicAuthProvider) ValidateCredentials() error {
	// Check username is not empty
	if p.credentials.Username == "" {
		return &CredentialValidationError{
			Field:  "username",
			Reason: "username cannot be empty",
		}
	}

	// Check for control characters in username
	// Control characters (ASCII 0-31 and 127) can be used for:
	// - Terminal escape sequence attacks
	// - Log injection attacks
	// - Command injection attacks
	if containsControlChars(p.credentials.Username) {
		return &CredentialValidationError{
			Field:  "username",
			Reason: "username contains control characters",
		}
	}

	// Check password is not empty
	if p.credentials.Password == "" {
		return &CredentialValidationError{
			Field:  "password",
			Reason: "password cannot be empty",
		}
	}

	// Password can contain ANY characters (including quotes, spaces, unicode, control chars)
	// HTTP Basic Auth base64-encodes the entire username:password string,
	// so all characters are safely transmitted.
	// No additional validation is needed.

	return nil
}

// Helper functions

// containsControlChars checks if a string contains control characters.
//
// Control characters are:
//   - ASCII 0-31: Null, tab, newline, escape, etc.
//   - ASCII 127: Delete character
//
// These characters are often used in terminal escape sequences and can be
// exploited for injection attacks. Regular printable characters (ASCII 32-126)
// and unicode characters are NOT considered control characters.
//
// Parameters:
//   - s: String to check
//
// Returns:
//   - bool: true if string contains control characters, false otherwise
//
// Examples:
//   - "username" -> false (no control chars)
//   - "user\n" -> true (newline is control char)
//   - "user name" -> false (space is ASCII 32, NOT a control char)
//   - "user\x1b[31m" -> true (escape sequence)
//   - "пароль" -> false (unicode is not control char)
func containsControlChars(s string) bool {
	for _, r := range s {
		// Check for control characters:
		// - ASCII 0-31: Control characters including null, tab, newline, escape
		// - ASCII 127: Delete character
		// Note: Space (ASCII 32) is NOT a control character
		if r < 32 || r == 127 {
			return true
		}
	}
	return false
}
