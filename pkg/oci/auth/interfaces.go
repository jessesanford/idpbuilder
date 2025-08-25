package auth

import (
	"context"
	"time"
)

// AuthProvider defines the contract for providing authentication credentials
// to OCI registries. Implementations handle credential management, validation,
// and token lifecycle operations.
type AuthProvider interface {
	// GetCredentials retrieves authentication credentials for the given registry
	GetCredentials(ctx context.Context, registry string) (*Credentials, error)
	
	// ValidateCredentials checks if the provided credentials are valid and not expired
	ValidateCredentials(ctx context.Context, creds *Credentials) (*AuthResult, error)
	
	// RefreshToken attempts to refresh an expired token using refresh credentials
	RefreshToken(ctx context.Context, token *Token) (*Token, error)
	
	// SupportsRegistry indicates whether this provider can authenticate to the registry
	SupportsRegistry(registry string) bool
}

// CredentialStore defines the contract for persistent storage of authentication
// credentials. Implementations may store credentials in files, system keychains,
// or other secure storage mechanisms.
type CredentialStore interface {
	// Load retrieves stored credentials for the specified registry
	Load(ctx context.Context, registry string) (*Credentials, error)
	
	// Save stores credentials for the specified registry securely
	Save(ctx context.Context, registry string, creds *Credentials) error
	
	// Delete removes stored credentials for the specified registry
	Delete(ctx context.Context, registry string) error
	
	// List returns all registries that have stored credentials
	List(ctx context.Context) ([]string, error)
	
	// Exists checks if credentials are stored for the specified registry
	Exists(ctx context.Context, registry string) bool
}

// CertificateProvider defines the contract for managing TLS certificates
// used in registry authentication. Handles certificate loading, validation,
// and bundle creation for secure connections.
type CertificateProvider interface {
	// GetCertBundle retrieves a certificate bundle for TLS authentication
	GetCertBundle(ctx context.Context, registry string) (interface{}, error)
	
	// LoadFromKind loads certificates from Kubernetes secrets or other sources
	LoadFromKind(ctx context.Context, kind, namespace, name string) (*TLSConfig, error)
	
	// ValidateCertificate validates a certificate against policy requirements
	ValidateCertificate(ctx context.Context, cert interface{}) error
	
	// GetCABundle retrieves the Certificate Authority bundle for verification
	GetCABundle(ctx context.Context, registry string) ([]interface{}, error)
}

// TokenManager defines the contract for managing authentication tokens,
// including caching, validation, and lifecycle management.
type TokenManager interface {
	// CacheToken stores a token for future use with expiration handling
	CacheToken(ctx context.Context, registry string, token *Token) error
	
	// GetCachedToken retrieves a cached token if valid and not expired
	GetCachedToken(ctx context.Context, registry string) (*Token, error)
	
	// InvalidateToken removes a token from cache (e.g., on authentication failure)
	InvalidateToken(ctx context.Context, registry string) error
	
	// ValidateToken checks if a token is valid and not expired
	ValidateToken(ctx context.Context, token *Token) bool
	
	// GetTokenExpiry returns the expiration time for a token
	GetTokenExpiry(token *Token) time.Time
}