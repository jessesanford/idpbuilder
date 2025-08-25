// Copyright 2024 The IDP Builder Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package auth

import (
	"encoding/json"
	"time"
)

// RegistryAuth defines the interface for registry authentication mechanisms
type RegistryAuth interface {
	// GetCredentials returns the credentials for the given registry
	GetCredentials(registry string) (*Credentials, error)
	
	// Validate checks if the authentication configuration is valid
	Validate() error
	
	// Type returns the authentication type
	Type() AuthType
	
	// Refresh refreshes authentication credentials if supported
	Refresh() error
	
	// IsExpired checks if the authentication has expired
	IsExpired() bool
}

// AuthConfig represents the registry authentication configuration
type AuthConfig struct {
	// Registry is the registry URL or hostname
	Registry string `json:"registry"`
	
	// Type specifies the authentication type
	Type AuthType `json:"type"`
	
	// Username for basic authentication
	Username string `json:"username,omitempty"`
	
	// Password for basic authentication
	Password string `json:"password,omitempty"`
	
	// Token for bearer token authentication
	Token string `json:"token,omitempty"`
	
	// RefreshToken for OAuth2 token refresh
	RefreshToken string `json:"refresh_token,omitempty"`
	
	// ExpiresAt indicates when the token expires
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
	
	// InsecureSkipTLSVerify allows skipping TLS verification
	InsecureSkipTLSVerify bool `json:"insecure_skip_tls_verify,omitempty"`
	
	// CredentialHelper specifies external credential helper
	CredentialHelper string `json:"credential_helper,omitempty"`
}

// DockerConfig represents the structure of Docker's config.json file
type DockerConfig struct {
	// Auths contains registry authentication configurations
	Auths map[string]DockerAuthEntry `json:"auths,omitempty"`
	
	// CredStore specifies the default credential store
	CredStore string `json:"credStore,omitempty"`
	
	// CredHelpers maps registries to specific credential helpers
	CredHelpers map[string]string `json:"credHelpers,omitempty"`
	
	// HttpHeaders contains additional HTTP headers to send
	HttpHeaders map[string]string `json:"HttpHeaders,omitempty"`
	
	// PsFormat specifies the default format for docker ps
	PsFormat string `json:"psFormat,omitempty"`
	
	// ImagesFormat specifies the default format for docker images
	ImagesFormat string `json:"imagesFormat,omitempty"`
}

// DockerAuthEntry represents an entry in Docker's auths section
type DockerAuthEntry struct {
	// Username for authentication
	Username string `json:"username,omitempty"`
	
	// Password for authentication
	Password string `json:"password,omitempty"`
	
	// Email associated with the account (deprecated but kept for compatibility)
	Email string `json:"email,omitempty"`
	
	// Auth contains base64 encoded username:password
	Auth string `json:"auth,omitempty"`
	
	// IdentityToken for bearer token authentication
	IdentityToken string `json:"identitytoken,omitempty"`
	
	// RegistryToken for registry-specific tokens
	RegistryToken string `json:"registrytoken,omitempty"`
}

// AuthStore defines the interface for credential storage systems
type AuthStore interface {
	// Get retrieves credentials for the specified registry
	Get(registry string) (*Credentials, error)
	
	// Set stores credentials for the specified registry
	Set(registry string, creds *Credentials) error
	
	// Delete removes credentials for the specified registry
	Delete(registry string) error
	
	// List returns all stored registries
	List() ([]string, error)
	
	// Clear removes all stored credentials
	Clear() error
}

// RegistryAuthOptions contains options for registry authentication
type RegistryAuthOptions struct {
	// ConfigDir specifies the directory containing Docker config
	ConfigDir string
	
	// ConfigFile specifies the path to Docker config file
	ConfigFile string
	
	// CredentialStore specifies the default credential store to use
	CredentialStore string
	
	// AllowInsecure allows insecure registries
	AllowInsecure bool
	
	// Timeout for authentication operations
	Timeout time.Duration
	
	// RetryAttempts specifies number of retry attempts
	RetryAttempts int
	
	// UserAgent to send in HTTP requests
	UserAgent string
}

// BasicAuth implements RegistryAuth for HTTP Basic authentication
type BasicAuth struct {
	username string
	password string
	registry string
}

// NewBasicAuth creates a new BasicAuth instance
func NewBasicAuth(registry, username, password string) *BasicAuth {
	return &BasicAuth{
		registry: registry,
		username: username,
		password: password,
	}
}

// GetCredentials implements RegistryAuth interface
func (b *BasicAuth) GetCredentials(registry string) (*Credentials, error) {
	if registry != b.registry {
		return nil, &AuthError{
			Type:    ErrMissingCredentials,
			Message: "no credentials for registry: " + registry,
		}
	}
	
	return &Credentials{
		Username: b.username,
		Password: b.password,
	}, nil
}

// Validate implements RegistryAuth interface
func (b *BasicAuth) Validate() error {
	if b.username == "" || b.password == "" {
		return &AuthError{
			Type:    ErrInvalidCredentials,
			Message: "username and password are required for basic auth",
		}
	}
	return nil
}

// Type implements RegistryAuth interface
func (b *BasicAuth) Type() AuthType {
	return AuthTypeBasic
}

// Refresh implements RegistryAuth interface
func (b *BasicAuth) Refresh() error {
	// Basic auth doesn't support refresh
	return nil
}

// IsExpired implements RegistryAuth interface
func (b *BasicAuth) IsExpired() bool {
	// Basic auth doesn't expire
	return false
}

// BearerAuth implements RegistryAuth for Bearer token authentication
type BearerAuth struct {
	token     string
	registry  string
	expiresAt *time.Time
}

// NewBearerAuth creates a new BearerAuth instance
func NewBearerAuth(registry, token string, expiresAt *time.Time) *BearerAuth {
	return &BearerAuth{
		registry:  registry,
		token:     token,
		expiresAt: expiresAt,
	}
}

// GetCredentials implements RegistryAuth interface
func (b *BearerAuth) GetCredentials(registry string) (*Credentials, error) {
	if registry != b.registry {
		return nil, &AuthError{
			Type:    ErrMissingCredentials,
			Message: "no credentials for registry: " + registry,
		}
	}
	
	if b.IsExpired() {
		return nil, &AuthError{
			Type:    ErrTokenExpired,
			Message: "bearer token has expired",
		}
	}
	
	return &Credentials{
		Token: b.token,
	}, nil
}

// Validate implements RegistryAuth interface
func (b *BearerAuth) Validate() error {
	if b.token == "" {
		return &AuthError{
			Type:    ErrInvalidToken,
			Message: "token is required for bearer auth",
		}
	}
	return nil
}

// Type implements RegistryAuth interface
func (b *BearerAuth) Type() AuthType {
	return AuthTypeBearer
}

// Refresh implements RegistryAuth interface
func (b *BearerAuth) Refresh() error {
	// Token refresh would require additional OAuth2 flow
	return &AuthError{
		Type:    ErrUnsupportedAuthType,
		Message: "token refresh not supported for bearer auth",
	}
}

// IsExpired implements RegistryAuth interface
func (b *BearerAuth) IsExpired() bool {
	if b.expiresAt == nil {
		return false
	}
	return time.Now().After(*b.expiresAt)
}

// AuthError represents an authentication error
type AuthError struct {
	Type    string
	Message string
	Cause   error
}

// Error implements the error interface
func (e *AuthError) Error() string {
	if e.Cause != nil {
		return e.Message + ": " + e.Cause.Error()
	}
	return e.Message
}

// Unwrap implements error unwrapping
func (e *AuthError) Unwrap() error {
	return e.Cause
}

// String implements fmt.Stringer for AuthType
func (a AuthType) String() string {
	return string(a)
}

// IsValid checks if the AuthType is valid
func (a AuthType) IsValid() bool {
	switch a {
	case AuthTypeBasic, AuthTypeBearer, AuthTypeOAuth2, AuthTypeAnonymous:
		return true
	default:
		return false
	}
}

// MarshalJSON implements json.Marshaler for AuthType
func (a AuthType) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(a))
}

// UnmarshalJSON implements json.Unmarshaler for AuthType
func (a *AuthType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	*a = AuthType(s)
	return nil
}