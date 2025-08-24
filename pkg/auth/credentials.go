package auth

import (
	"encoding/base64"
	"encoding/json"
	"strings"
	"sync"
	"time"
)

// Credentials represents authentication credentials for a registry.
// It supports multiple authentication methods and secure handling of sensitive data.
type Credentials struct {
	// Registry is the registry URL these credentials apply to.
	Registry string `json:"registry,omitempty"`

	// Username for basic authentication.
	Username string `json:"username,omitempty"`

	// Password for basic authentication.
	Password string `json:"password,omitempty"`

	// Token for bearer token authentication.
	Token string `json:"token,omitempty"`

	// IdentityToken for registry identity-based authentication.
	IdentityToken string `json:"identity_token,omitempty"`

	// RefreshToken for OAuth2 authentication flows.
	RefreshToken string `json:"refresh_token,omitempty"`

	// ExpiresAt indicates when these credentials expire.
	ExpiresAt *time.Time `json:"expires_at,omitempty"`

	// Scopes contains the access scopes granted for these credentials.
	Scopes []string `json:"scopes,omitempty"`
}

// IsExpired checks if the credentials are expired.
func (c *Credentials) IsExpired() bool {
	if c.ExpiresAt == nil {
		return false
	}
	return time.Now().After(*c.ExpiresAt)
}

// IsEmpty checks if the credentials contain no authentication data.
func (c *Credentials) IsEmpty() bool {
	return c.Username == "" && c.Password == "" && c.Token == "" && c.IdentityToken == ""
}

// ToBasicAuth returns the credentials as HTTP Basic Auth header value.
func (c *Credentials) ToBasicAuth() string {
	if c.Username == "" && c.Password == "" {
		return ""
	}
	auth := c.Username + ":" + c.Password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

// Clear securely clears sensitive credential data from memory.
func (c *Credentials) Clear() {
	// Overwrite sensitive fields with empty strings
	c.Password = ""
	c.Token = ""
	c.IdentityToken = ""
	c.RefreshToken = ""
}

// CredentialHelper defines the interface for external credential helpers.
// These are programs that can retrieve credentials from secure storage systems.
type CredentialHelper interface {
	// Get retrieves credentials for the specified registry.
	Get(registryURL string) (*Credentials, error)

	// Store saves credentials for the specified registry.
	Store(registryURL string, creds *Credentials) error

	// Delete removes credentials for the specified registry.
	Delete(registryURL string) error

	// List returns all registries with stored credentials.
	List() ([]string, error)
}

// CredentialStore provides in-memory storage for authentication credentials.
// It is thread-safe and supports concurrent access.
type CredentialStore struct {
	// store maps registry URLs to their credentials
	store map[string]*Credentials
	// mutex protects concurrent access to the store
	mutex sync.RWMutex
}

// NewCredentialStore creates a new empty credential store.
func NewCredentialStore() *CredentialStore {
	return &CredentialStore{
		store: make(map[string]*Credentials),
	}
}

// Get retrieves credentials for the specified registry.
func (cs *CredentialStore) Get(registry string) (*Credentials, error) {
	cs.mutex.RLock()
	defer cs.mutex.RUnlock()

	creds, exists := cs.store[registry]
	if !exists {
		return nil, ErrCredentialsNotFound
	}

	// Return a copy to prevent external modification
	return &Credentials{
		Registry:      creds.Registry,
		Username:      creds.Username,
		Password:      creds.Password,
		Token:         creds.Token,
		IdentityToken: creds.IdentityToken,
		RefreshToken:  creds.RefreshToken,
		ExpiresAt:     creds.ExpiresAt,
		Scopes:        append([]string(nil), creds.Scopes...),
	}, nil
}

// Set stores credentials for the specified registry.
func (cs *CredentialStore) Set(registry string, creds *Credentials) error {
	if registry == "" {
		return ErrInvalidRegistry
	}
	if creds == nil {
		return ErrInvalidCredentials
	}

	cs.mutex.Lock()
	defer cs.mutex.Unlock()

	// Store a copy to prevent external modification
	cs.store[registry] = &Credentials{
		Registry:      registry,
		Username:      creds.Username,
		Password:      creds.Password,
		Token:         creds.Token,
		IdentityToken: creds.IdentityToken,
		RefreshToken:  creds.RefreshToken,
		ExpiresAt:     creds.ExpiresAt,
		Scopes:        append([]string(nil), creds.Scopes...),
	}
	return nil
}

// Delete removes credentials for the specified registry.
func (cs *CredentialStore) Delete(registry string) error {
	cs.mutex.Lock()
	defer cs.mutex.Unlock()

	if _, exists := cs.store[registry]; !exists {
		return ErrCredentialsNotFound
	}

	// Clear sensitive data before deletion
	if creds, exists := cs.store[registry]; exists {
		creds.Clear()
	}

	delete(cs.store, registry)
	return nil
}

// List returns all registry URLs with stored credentials.
func (cs *CredentialStore) List() ([]string, error) {
	cs.mutex.RLock()
	defer cs.mutex.RUnlock()

	registries := make([]string, 0, len(cs.store))
	for registry := range cs.store {
		registries = append(registries, registry)
	}
	return registries, nil
}

// Clear removes all stored credentials.
func (cs *CredentialStore) Clear() error {
	cs.mutex.Lock()
	defer cs.mutex.Unlock()

	// Clear sensitive data before deletion
	for _, creds := range cs.store {
		creds.Clear()
	}

	cs.store = make(map[string]*Credentials)
	return nil
}

// TokenResponse represents the response from an OAuth2 token endpoint.
type TokenResponse struct {
	// AccessToken is the bearer token for accessing protected resources.
	AccessToken string `json:"access_token"`

	// TokenType specifies the type of token returned (usually "bearer").
	TokenType string `json:"token_type"`

	// ExpiresIn is the number of seconds until the token expires.
	ExpiresIn int `json:"expires_in"`

	// RefreshToken can be used to obtain new access tokens.
	RefreshToken string `json:"refresh_token,omitempty"`

	// Scope defines the access scope granted by the token.
	Scope string `json:"scope,omitempty"`
}

// ToCredentials converts a TokenResponse to Credentials.
func (tr *TokenResponse) ToCredentials(registry string) *Credentials {
	creds := &Credentials{
		Registry: registry,
		Token:    tr.AccessToken,
	}

	if tr.RefreshToken != "" {
		creds.RefreshToken = tr.RefreshToken
	}

	if tr.ExpiresIn > 0 {
		expiresAt := time.Now().Add(time.Duration(tr.ExpiresIn) * time.Second)
		creds.ExpiresAt = &expiresAt
	}

	if tr.Scope != "" {
		creds.Scopes = strings.Fields(tr.Scope)
	}

	return creds
}