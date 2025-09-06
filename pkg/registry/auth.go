package registry

import (
	"encoding/base64"
	"fmt"
	"os"
	"strings"

	"github.com/google/go-containerregistry/pkg/authn"
)

// configureAuth configures authentication for registry operations
func configureAuth(username, password string) authn.Authenticator {
	// Check for empty credentials
	if username == "" && password == "" {
		return authn.Anonymous
	}

	// Validate credentials
	if err := validateCredentials(username, password); err != nil {
		// Log warning but don't fail - fall back to anonymous
		return authn.Anonymous
	}

	return &authn.Basic{
		Username: username,
		Password: password,
	}
}

// GetAuthToken generates a base64 encoded auth token
func GetAuthToken(username, password string) string {
	if username == "" && password == "" {
		return ""
	}

	auth := fmt.Sprintf("%s:%s", username, password)
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

// ValidateCredentials checks if credentials are valid
func validateCredentials(username, password string) error {
	if username == "" {
		return fmt.Errorf("username cannot be empty")
	}
	if password == "" {
		return fmt.Errorf("password cannot be empty")
	}

	// Additional validation rules
	if len(username) < 2 {
		return fmt.Errorf("username must be at least 2 characters")
	}
	if len(password) < 4 {
		return fmt.Errorf("password must be at least 4 characters")
	}

	// Check for invalid characters
	if strings.Contains(username, ":") {
		return fmt.Errorf("username cannot contain colon character")
	}

	return nil
}

// LoadCredentialsFromEnv loads credentials from environment variables
func LoadCredentialsFromEnv() (username, password string) {
	// Standard Docker registry environment variables
	username = os.Getenv("REGISTRY_USERNAME")
	password = os.Getenv("REGISTRY_PASSWORD")

	// Gitea-specific environment variables
	if username == "" {
		username = os.Getenv("GITEA_USERNAME")
	}
	if password == "" {
		password = os.Getenv("GITEA_PASSWORD")
	}

	// IDPBuilder-specific environment variables
	if username == "" {
		username = os.Getenv("IDPBUILDER_REGISTRY_USERNAME")
	}
	if password == "" {
		password = os.Getenv("IDPBUILDER_REGISTRY_PASSWORD")
	}

	return username, password
}

// AuthConfig holds authentication configuration
type AuthConfig struct {
	Username  string
	Password  string
	Token     string
	Anonymous bool
}

// NewAuthConfig creates authentication configuration
func NewAuthConfig(username, password string) *AuthConfig {
	config := &AuthConfig{
		Username: username,
		Password: password,
	}

	// Determine if anonymous
	config.Anonymous = (username == "" && password == "")

	// Generate token if credentials provided
	if !config.Anonymous {
		config.Token = GetAuthToken(username, password)
	}

	return config
}

// ToAuthenticator converts AuthConfig to go-containerregistry authenticator
func (a *AuthConfig) ToAuthenticator() authn.Authenticator {
	if a.Anonymous {
		return authn.Anonymous
	}

	return &authn.Basic{
		Username: a.Username,
		Password: a.Password,
	}
}

// Validate checks if the authentication configuration is valid
func (a *AuthConfig) Validate() error {
	if a.Anonymous {
		return nil // Anonymous is always valid
	}

	return validateCredentials(a.Username, a.Password)
}

// Clone creates a copy of the AuthConfig
func (a *AuthConfig) Clone() *AuthConfig {
	return &AuthConfig{
		Username:  a.Username,
		Password:  a.Password,
		Token:     a.Token,
		Anonymous: a.Anonymous,
	}
}
