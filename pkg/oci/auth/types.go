package auth

import (
	"time"
)

// Credentials represents authentication credentials for registry access.
// It supports multiple authentication methods and includes lifecycle tracking.
type Credentials struct {
	// Registry is the registry hostname or URL these credentials apply to
	Registry string `json:"registry"`

	// Username for basic authentication (optional for token auth)
	Username string `json:"username,omitempty"`

	// Password for basic authentication (optional for token auth, omitempty for security)
	Password string `json:"password,omitempty"`

	// Token for token-based authentication
	Token *Token `json:"token,omitempty"`

	// AuthMethod specifies the authentication method to use
	AuthMethod AuthMethod `json:"auth_method"`

	// ExpiresAt indicates when these credentials expire (optional)
	ExpiresAt *time.Time `json:"expires_at,omitempty"`

	// CreatedAt tracks when these credentials were created
	CreatedAt time.Time `json:"created_at"`

	// Metadata stores additional provider-specific information
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// Token represents an authentication token with lifecycle and metadata.
type Token struct {
	// Value is the actual token string
	Value string `json:"value"`

	// Type specifies the token type (Bearer, Basic, etc.)
	Type TokenType `json:"type"`

	// IssuedAt indicates when the token was issued
	IssuedAt time.Time `json:"issued_at"`

	// ExpiresAt indicates when the token expires
	ExpiresAt time.Time `json:"expires_at"`

	// RefreshToken for OAuth2 token refresh
	RefreshToken string `json:"refresh_token,omitempty"`

	// Scope defines the permissions granted by this token
	Scope []string `json:"scope,omitempty"`
}

// TLSConfig represents complete TLS configuration for registry connections.
type TLSConfig struct {
	// InsecureSkipVerify disables certificate verification (dev only)
	InsecureSkipVerify bool `json:"insecure_skip_verify"`

	// ServerName for certificate verification (SNI)
	ServerName string `json:"server_name,omitempty"`

	// MinVersion specifies minimum TLS version
	MinVersion uint16 `json:"min_version,omitempty"`

	// MaxVersion specifies maximum TLS version
	MaxVersion uint16 `json:"max_version,omitempty"`

	// CipherSuites specifies allowed cipher suites
	CipherSuites []uint16 `json:"cipher_suites,omitempty"`

	// CertificateBundle contains CA certificates and client certificates
	CertificateBundle *CertificateBundle `json:"certificate_bundle,omitempty"`
}

// CertificateBundle contains certificate data and metadata.
type CertificateBundle struct {
	// CACerts contains PEM-encoded CA certificates
	CACerts []byte `json:"ca_certs,omitempty"`

	// ClientCert contains PEM-encoded client certificate
	ClientCert []byte `json:"client_cert,omitempty"`

	// ClientKey contains PEM-encoded client private key
	ClientKey []byte `json:"client_key,omitempty"`

	// Bundle contains combined certificate bundle
	Bundle []byte `json:"bundle,omitempty"`
}

// CertificateInfo contains metadata about a certificate.
type CertificateInfo struct {
	// Subject contains the certificate subject information
	Subject string `json:"subject"`

	// Issuer contains the certificate issuer information
	Issuer string `json:"issuer"`

	// SerialNumber is the certificate serial number
	SerialNumber string `json:"serial_number"`

	// NotBefore indicates when the certificate becomes valid
	NotBefore time.Time `json:"not_before"`

	// NotAfter indicates when the certificate expires
	NotAfter time.Time `json:"not_after"`

	// DNSNames contains subject alternative names
	DNSNames []string `json:"dns_names,omitempty"`

	// IPAddresses contains IP subject alternative names
	IPAddresses []string `json:"ip_addresses,omitempty"`
}

// BasicAuth represents simple username/password authentication.
type BasicAuth struct {
	// Username for authentication
	Username string `json:"username"`

	// Password for authentication
	Password string `json:"password"`
}

// TokenAuth represents token-based authentication.
type TokenAuth struct {
	// Token value
	Token string `json:"token"`

	// Type specifies token type
	Type TokenType `json:"type"`
}

// OAuth2Auth represents OAuth2 authentication configuration.
type OAuth2Auth struct {
	// ClientID for OAuth2 client credentials flow
	ClientID string `json:"client_id"`

	// ClientSecret for OAuth2 client credentials flow
	ClientSecret string `json:"client_secret"`

	// TokenURL is the OAuth2 token endpoint
	TokenURL string `json:"token_url"`

	// Scopes specify requested permissions
	Scopes []string `json:"scopes,omitempty"`

	// AccessToken is the current access token
	AccessToken string `json:"access_token,omitempty"`

	// RefreshToken for token refresh
	RefreshToken string `json:"refresh_token,omitempty"`

	// ExpiresAt indicates token expiration
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
}

// FileStore represents file-based credential storage configuration.
type FileStore struct {
	// Path to the credentials file
	Path string `json:"path"`

	// Encrypted indicates if the file should be encrypted
	Encrypted bool `json:"encrypted"`

	// FileMode specifies file permissions
	FileMode uint32 `json:"file_mode,omitempty"`
}

// MemoryStore represents in-memory credential storage configuration.
type MemoryStore struct {
	// TTL specifies how long to keep credentials in memory
	TTL time.Duration `json:"ttl"`

	// MaxSize limits the number of stored credentials
	MaxSize int `json:"max_size"`
}

// AuthResult represents the result of an authentication attempt.
type AuthResult struct {
	// Success indicates if authentication succeeded
	Success bool `json:"success"`

	// Token contains the authentication token if successful
	Token *Token `json:"token,omitempty"`

	// Error contains error information if authentication failed
	Error string `json:"error,omitempty"`

	// RetryAfter indicates when to retry authentication (for rate limiting)
	RetryAfter *time.Time `json:"retry_after,omitempty"`
}

// AuthConfig represents overall authentication configuration.
type AuthConfig struct {
	// Method specifies the authentication method to use
	Method AuthMethod `json:"method"`

	// StoreType specifies the credential storage backend
	StoreType StoreType `json:"store_type"`

	// FileStore configuration (when StoreType is File)
	FileStore *FileStore `json:"file_store,omitempty"`

	// MemoryStore configuration (when StoreType is Memory)
	MemoryStore *MemoryStore `json:"memory_store,omitempty"`

	// TLSConfig for secure connections
	TLSConfig *TLSConfig `json:"tls_config,omitempty"`
}

// AuthMethod defines supported authentication methods
type AuthMethod string

const (
	// AuthMethodBasic for username/password authentication
	AuthMethodBasic AuthMethod = "basic"

	// AuthMethodToken for token-based authentication
	AuthMethodToken AuthMethod = "token"

	// AuthMethodOAuth2 for OAuth2 authentication
	AuthMethodOAuth2 AuthMethod = "oauth2"
)

// TokenType defines supported token types
type TokenType string

const (
	// TokenTypeBearer for Bearer tokens
	TokenTypeBearer TokenType = "Bearer"

	// TokenTypeBasic for Basic tokens
	TokenTypeBasic TokenType = "Basic"

	// TokenTypeOAuth2 for OAuth2 access tokens
	TokenTypeOAuth2 TokenType = "OAuth2"
)

// StoreType defines supported storage backends
type StoreType string

const (
	// StoreTypeFile for file-based storage
	StoreTypeFile StoreType = "file"

	// StoreTypeMemory for in-memory storage
	StoreTypeMemory StoreType = "memory"

	// StoreTypeKeyring for OS keyring storage (future)
	StoreTypeKeyring StoreType = "keyring"
)