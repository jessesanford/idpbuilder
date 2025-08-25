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
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// Credentials represents authentication credentials for a registry
type Credentials struct {
	// Username for basic authentication
	Username string `json:"username,omitempty"`
	
	// Password for basic authentication  
	Password string `json:"password,omitempty"`
	
	// Token for bearer token authentication
	Token string `json:"token,omitempty"`
	
	// RefreshToken for OAuth2 authentication
	RefreshToken string `json:"refresh_token,omitempty"`
	
	// AccessToken for OAuth2 authentication
	AccessToken string `json:"access_token,omitempty"`
	
	// ExpiresAt indicates when the credentials expire
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
	
	// TokenType specifies the type of token (Bearer, Basic, etc.)
	TokenType string `json:"token_type,omitempty"`
	
	// Scope defines the permissions granted by the token
	Scope string `json:"scope,omitempty"`
}

// CredentialHelper defines the interface for external credential helpers
type CredentialHelper interface {
	// Get retrieves credentials from the helper
	Get(serverURL string) (*Credentials, error)
	
	// Store saves credentials using the helper
	Store(serverURL string, creds *Credentials) error
	
	// Delete removes credentials from the helper
	Delete(serverURL string) error
	
	// List returns all stored server URLs
	List() ([]string, error)
	
	// IsAvailable checks if the credential helper is available
	IsAvailable() bool
}

// CredentialStore represents an in-memory credential store
type CredentialStore struct {
	credentials map[string]*Credentials
}

// NewCredentialStore creates a new in-memory credential store
func NewCredentialStore() *CredentialStore {
	return &CredentialStore{
		credentials: make(map[string]*Credentials),
	}
}

// Get implements AuthStore interface
func (cs *CredentialStore) Get(registry string) (*Credentials, error) {
	creds, exists := cs.credentials[registry]
	if !exists {
		return nil, &AuthError{
			Type:    ErrMissingCredentials,
			Message: fmt.Sprintf("no credentials found for registry: %s", registry),
		}
	}
	
	// Check if credentials are expired
	if creds.IsExpired() {
		return nil, &AuthError{
			Type:    ErrTokenExpired,
			Message: fmt.Sprintf("credentials for registry %s have expired", registry),
		}
	}
	
	return creds, nil
}

// Set implements AuthStore interface
func (cs *CredentialStore) Set(registry string, creds *Credentials) error {
	if registry == "" {
		return &AuthError{
			Type:    ErrInvalidCredentials,
			Message: "registry cannot be empty",
		}
	}
	
	if creds == nil {
		return &AuthError{
			Type:    ErrInvalidCredentials,
			Message: "credentials cannot be nil",
		}
	}
	
	if err := creds.Validate(); err != nil {
		return err
	}
	
	cs.credentials[registry] = creds
	return nil
}

// Delete implements AuthStore interface
func (cs *CredentialStore) Delete(registry string) error {
	delete(cs.credentials, registry)
	return nil
}

// List implements AuthStore interface
func (cs *CredentialStore) List() ([]string, error) {
	registries := make([]string, 0, len(cs.credentials))
	for registry := range cs.credentials {
		registries = append(registries, registry)
	}
	return registries, nil
}

// Clear implements AuthStore interface
func (cs *CredentialStore) Clear() error {
	cs.credentials = make(map[string]*Credentials)
	return nil
}

// TokenResponse represents an OAuth2 token response
type TokenResponse struct {
	// AccessToken is the access token
	AccessToken string `json:"access_token"`
	
	// TokenType specifies the token type (usually "Bearer")
	TokenType string `json:"token_type"`
	
	// ExpiresIn specifies token lifetime in seconds
	ExpiresIn int `json:"expires_in"`
	
	// RefreshToken for refreshing the access token
	RefreshToken string `json:"refresh_token,omitempty"`
	
	// Scope defines the permissions granted
	Scope string `json:"scope,omitempty"`
	
	// IssuedAt indicates when the token was issued
	IssuedAt time.Time `json:"issued_at,omitempty"`
}

// ToCredentials converts TokenResponse to Credentials
func (tr *TokenResponse) ToCredentials() *Credentials {
	var expiresAt *time.Time
	if tr.ExpiresIn > 0 {
		expiry := tr.IssuedAt.Add(time.Duration(tr.ExpiresIn) * time.Second)
		expiresAt = &expiry
	}
	
	return &Credentials{
		AccessToken:  tr.AccessToken,
		RefreshToken: tr.RefreshToken,
		TokenType:    tr.TokenType,
		Scope:        tr.Scope,
		ExpiresAt:    expiresAt,
	}
}

// Validate checks if the credentials are valid
func (c *Credentials) Validate() error {
	// At least one form of authentication must be present
	if c.Username == "" && c.Password == "" && c.Token == "" && c.AccessToken == "" {
		return &AuthError{
			Type:    ErrInvalidCredentials,
			Message: "credentials must contain username/password, token, or access token",
		}
	}
	
	// If username is provided, password should also be provided for basic auth
	if c.Username != "" && c.Password == "" {
		return &AuthError{
			Type:    ErrInvalidCredentials,
			Message: "password is required when username is provided",
		}
	}
	
	return nil
}

// IsExpired checks if the credentials have expired
func (c *Credentials) IsExpired() bool {
	if c.ExpiresAt == nil {
		return false
	}
	return time.Now().After(*c.ExpiresAt)
}

// IsNearExpiry checks if credentials will expire soon
func (c *Credentials) IsNearExpiry(threshold time.Duration) bool {
	if c.ExpiresAt == nil {
		return false
	}
	return time.Until(*c.ExpiresAt) < threshold
}

// ToBasicAuth encodes username and password as Basic auth string
func (c *Credentials) ToBasicAuth() string {
	if c.Username == "" || c.Password == "" {
		return ""
	}
	auth := c.Username + ":" + c.Password
	return BasicAuthPrefix + base64.StdEncoding.EncodeToString([]byte(auth))
}

// ToBearerAuth formats the token as Bearer auth string
func (c *Credentials) ToBearerAuth() string {
	token := c.Token
	if token == "" {
		token = c.AccessToken
	}
	if token == "" {
		return ""
	}
	return BearerAuthPrefix + token
}

// ToAuthHeader returns the appropriate Authorization header value
func (c *Credentials) ToAuthHeader() string {
	// Prefer bearer token if available
	if bearerAuth := c.ToBearerAuth(); bearerAuth != "" {
		return bearerAuth
	}
	
	// Fall back to basic auth
	if basicAuth := c.ToBasicAuth(); basicAuth != "" {
		return basicAuth
	}
	
	return ""
}

// Clone creates a deep copy of the credentials
func (c *Credentials) Clone() *Credentials {
	clone := &Credentials{
		Username:     c.Username,
		Password:     c.Password,
		Token:        c.Token,
		RefreshToken: c.RefreshToken,
		AccessToken:  c.AccessToken,
		TokenType:    c.TokenType,
		Scope:        c.Scope,
	}
	
	if c.ExpiresAt != nil {
		expiry := *c.ExpiresAt
		clone.ExpiresAt = &expiry
	}
	
	return clone
}

// Redacted returns a copy of credentials with sensitive data redacted
func (c *Credentials) Redacted() *Credentials {
	redacted := c.Clone()
	
	if redacted.Password != "" {
		redacted.Password = "[REDACTED]"
	}
	if redacted.Token != "" {
		redacted.Token = "[REDACTED]"
	}
	if redacted.AccessToken != "" {
		redacted.AccessToken = "[REDACTED]"
	}
	if redacted.RefreshToken != "" {
		redacted.RefreshToken = "[REDACTED]"
	}
	
	return redacted
}

// MarshalJSON implements json.Marshaler with credential redaction
func (c *Credentials) MarshalJSON() ([]byte, error) {
	// Create a type alias to avoid infinite recursion
	type credentialsAlias Credentials
	return json.Marshal((*credentialsAlias)(c.Redacted()))
}

// String returns a string representation with redacted sensitive data
func (c *Credentials) String() string {
	var parts []string
	
	if c.Username != "" {
		parts = append(parts, fmt.Sprintf("username=%s", c.Username))
	}
	
	if c.Password != "" {
		parts = append(parts, "password=[REDACTED]")
	}
	
	if c.Token != "" {
		parts = append(parts, "token=[REDACTED]")
	}
	
	if c.AccessToken != "" {
		parts = append(parts, "access_token=[REDACTED]")
	}
	
	if c.RefreshToken != "" {
		parts = append(parts, "refresh_token=[REDACTED]")
	}
	
	if c.TokenType != "" {
		parts = append(parts, fmt.Sprintf("token_type=%s", c.TokenType))
	}
	
	if c.Scope != "" {
		parts = append(parts, fmt.Sprintf("scope=%s", c.Scope))
	}
	
	if c.ExpiresAt != nil {
		parts = append(parts, fmt.Sprintf("expires_at=%s", c.ExpiresAt.Format(time.RFC3339)))
	}
	
	return fmt.Sprintf("Credentials{%s}", strings.Join(parts, ", "))
}

// FromBasicAuth parses a Basic auth string into credentials
func FromBasicAuth(authHeader string) (*Credentials, error) {
	if !strings.HasPrefix(authHeader, BasicAuthPrefix) {
		return nil, &AuthError{
			Type:    ErrInvalidCredentials,
			Message: "not a Basic auth header",
		}
	}
	
	encoded := strings.TrimPrefix(authHeader, BasicAuthPrefix)
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, &AuthError{
			Type:    ErrInvalidCredentials,
			Message: "invalid base64 encoding in Basic auth",
			Cause:   err,
		}
	}
	
	parts := strings.SplitN(string(decoded), ":", 2)
	if len(parts) != 2 {
		return nil, &AuthError{
			Type:    ErrInvalidCredentials,
			Message: "invalid Basic auth format",
		}
	}
	
	return &Credentials{
		Username: parts[0],
		Password: parts[1],
	}, nil
}

// FromBearerAuth parses a Bearer auth string into credentials
func FromBearerAuth(authHeader string) (*Credentials, error) {
	if !strings.HasPrefix(authHeader, BearerAuthPrefix) {
		return nil, &AuthError{
			Type:    ErrInvalidToken,
			Message: "not a Bearer auth header",
		}
	}
	
	token := strings.TrimPrefix(authHeader, BearerAuthPrefix)
	if token == "" {
		return nil, &AuthError{
			Type:    ErrInvalidToken,
			Message: "empty Bearer token",
		}
	}
	
	return &Credentials{
		Token:     token,
		TokenType: "Bearer",
	}, nil
}