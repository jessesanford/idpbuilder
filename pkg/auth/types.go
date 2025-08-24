package auth

import (
	"encoding/json"
	"time"
)

// AuthType represents the type of authentication method used for registry access.
type AuthType int

const (
	// AuthTypeNone indicates no authentication is required.
	AuthTypeNone AuthType = iota
	// AuthTypeBasic indicates HTTP Basic Authentication using username/password.
	AuthTypeBasic
	// AuthTypeBearer indicates Bearer token authentication.
	AuthTypeBearer
	// AuthTypeOAuth2 indicates OAuth 2.0 authentication flow.
	AuthTypeOAuth2
)

// String returns the string representation of the AuthType.
func (a AuthType) String() string {
	switch a {
	case AuthTypeNone:
		return "none"
	case AuthTypeBasic:
		return "basic"
	case AuthTypeBearer:
		return "bearer"
	case AuthTypeOAuth2:
		return "oauth2"
	default:
		return "unknown"
	}
}

// RegistryAuth defines the interface for registry authentication methods.
// Implementations must provide credential retrieval, validation, and type identification.
type RegistryAuth interface {
	// GetCredentials returns the credentials for the given registry URL.
	// It returns an error if credentials cannot be retrieved or are invalid.
	GetCredentials(registryURL string) (*Credentials, error)

	// Validate checks if the authentication configuration is valid.
	// This includes verifying required fields and credential format.
	Validate() error

	// Type returns the authentication type used by this implementation.
	Type() AuthType

	// IsExpired checks if the current authentication credentials are expired.
	// For auth types that don't expire, this should return false.
	IsExpired() bool

	// Refresh attempts to refresh expired credentials if supported.
	// Returns an error if refresh is not supported or fails.
	Refresh() error
}

// AuthConfig represents the configuration for registry authentication.
// It supports multiple authentication methods and credential storage formats.
type AuthConfig struct {
	// Registry is the registry URL this configuration applies to.
	Registry string `json:"registry,omitempty"`

	// Username for basic authentication.
	Username string `json:"username,omitempty"`

	// Password for basic authentication.
	Password string `json:"password,omitempty"`

	// Token for bearer token authentication.
	Token string `json:"token,omitempty"`

	// IdentityToken for registry identity-based authentication.
	IdentityToken string `json:"identitytoken,omitempty"`

	// RefreshToken for OAuth2 authentication flows.
	RefreshToken string `json:"refreshtoken,omitempty"`

	// AuthType specifies the authentication method to use.
	AuthType AuthType `json:"auth_type,omitempty"`

	// TokenURL is the OAuth2 token endpoint for token refresh.
	TokenURL string `json:"token_url,omitempty"`

	// ClientID for OAuth2 authentication.
	ClientID string `json:"client_id,omitempty"`

	// ClientSecret for OAuth2 authentication.
	ClientSecret string `json:"client_secret,omitempty"`

	// Scope defines the OAuth2 access scope.
	Scope string `json:"scope,omitempty"`

	// ExpiresAt indicates when the current credentials expire.
	ExpiresAt *time.Time `json:"expires_at,omitempty"`

	// Insecure allows connections to registries with invalid certificates.
	Insecure bool `json:"insecure,omitempty"`
}

// GetCredentials extracts credentials from the AuthConfig.
func (a *AuthConfig) GetCredentials(registryURL string) (*Credentials, error) {
	creds := &Credentials{
		Registry:      registryURL,
		Username:      a.Username,
		Password:      a.Password,
		Token:         a.Token,
		IdentityToken: a.IdentityToken,
		RefreshToken:  a.RefreshToken,
		ExpiresAt:     a.ExpiresAt,
	}
	return creds, nil
}

// Validate checks if the AuthConfig has valid configuration for the specified auth type.
func (a *AuthConfig) Validate() error {
	switch a.AuthType {
	case AuthTypeBasic:
		if a.Username == "" || a.Password == "" {
			return ErrInvalidCredentials
		}
	case AuthTypeBearer:
		if a.Token == "" {
			return ErrInvalidCredentials
		}
	case AuthTypeOAuth2:
		if a.ClientID == "" || a.TokenURL == "" {
			return ErrInvalidCredentials
		}
	}
	return nil
}

// Type returns the authentication type configured.
func (a *AuthConfig) Type() AuthType {
	return a.AuthType
}

// IsExpired checks if the current credentials are expired.
func (a *AuthConfig) IsExpired() bool {
	if a.ExpiresAt == nil {
		return false
	}
	return time.Now().After(*a.ExpiresAt)
}

// Refresh attempts to refresh OAuth2 credentials using the refresh token.
func (a *AuthConfig) Refresh() error {
	if a.AuthType != AuthTypeOAuth2 || a.RefreshToken == "" {
		return ErrRefreshNotSupported
	}
	// Implementation would make HTTP request to token endpoint
	// This is a placeholder for the actual OAuth2 refresh logic
	return nil
}

// DockerConfig represents the structure of ~/.docker/config.json file.
// This is used to read existing Docker authentication configurations.
type DockerConfig struct {
	// Auths contains authentication configurations per registry.
	Auths map[string]DockerAuthConfig `json:"auths,omitempty"`

	// CredsStore specifies the credential store helper to use.
	CredsStore string `json:"credsStore,omitempty"`

	// CredHelpers maps registry hostnames to credential helper names.
	CredHelpers map[string]string `json:"credHelpers,omitempty"`
}

// DockerAuthConfig represents a single registry's authentication configuration in Docker format.
type DockerAuthConfig struct {
	// Auth is the base64 encoded username:password string.
	Auth string `json:"auth,omitempty"`

	// Username for the registry.
	Username string `json:"username,omitempty"`

	// Password for the registry.
	Password string `json:"password,omitempty"`

	// IdentityToken for registry authentication.
	IdentityToken string `json:"identitytoken,omitempty"`

	// RegistryToken for registry-specific tokens.
	RegistryToken string `json:"registrytoken,omitempty"`
}

// AuthStore defines the interface for storing and retrieving authentication configurations.
// Implementations can store credentials in memory, files, or secure storage systems.
type AuthStore interface {
	// Get retrieves authentication configuration for the specified registry.
	Get(registry string) (*AuthConfig, error)

	// Set stores authentication configuration for the specified registry.
	Set(registry string, config *AuthConfig) error

	// Delete removes authentication configuration for the specified registry.
	Delete(registry string) error

	// List returns all configured registries.
	List() ([]string, error)

	// Clear removes all stored authentication configurations.
	Clear() error
}

// RegistryAuthOptions contains options for configuring registry authentication.
type RegistryAuthOptions struct {
	// ConfigPath specifies the path to the Docker config file.
	ConfigPath string

	// UseCredentialStore enables use of system credential stores.
	UseCredentialStore bool

	// CredentialStore specifies which credential store to use.
	CredentialStore string

	// InsecureRegistries lists registries that should skip TLS verification.
	InsecureRegistries []string

	// DefaultAuth provides default authentication for registries not explicitly configured.
	DefaultAuth *AuthConfig
}