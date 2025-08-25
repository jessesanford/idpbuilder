package auth

import (
	"crypto/tls"
	"time"
)

// Credentials represents authentication credentials for OCI registry access.
type Credentials struct {
	Registry  string    `json:"registry" validate:"required,url"`
	Username  string    `json:"username,omitempty" validate:"required_without=Token"`
	Password  string    `json:"password,omitempty" validate:"required_with=Username"`
	Token     *Token    `json:"token,omitempty" validate:"required_without=Username"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

// Token represents an authentication token with metadata and expiration.
type Token struct {
	AccessToken  string    `json:"access_token" validate:"required"`
	RefreshToken string    `json:"refresh_token,omitempty"`
	TokenType    string    `json:"token_type" validate:"required,oneof=Bearer OAuth2"`
	ExpiresAt    time.Time `json:"expires_at" validate:"required"`
	Scope        []string  `json:"scope,omitempty"`
	Issuer       string    `json:"issuer,omitempty" validate:"omitempty,url"`
}

// TLSConfig contains TLS configuration for secure registry connections.
type TLSConfig struct {
	InsecureSkipVerify bool     `json:"insecure_skip_verify"`
	MinVersion         uint16   `json:"min_version" validate:"min=769"`
	MaxVersion         uint16   `json:"max_version" validate:"min=769"`
	CertFile           string   `json:"cert_file,omitempty"`
	KeyFile            string   `json:"key_file,omitempty" validate:"required_with=CertFile"`
	CAFile             string   `json:"ca_file,omitempty"`
	ServerName         string   `json:"server_name,omitempty"`
	CipherSuites       []uint16 `json:"cipher_suites,omitempty"`
}

// AuthConfig represents the complete authentication configuration for a registry.
type AuthConfig struct {
	Registry    string        `json:"registry" validate:"required,url"`
	Credentials *Credentials  `json:"credentials,omitempty"`
	TLSConfig   *TLSConfig    `json:"tls_config,omitempty"`
	Timeout     time.Duration `json:"timeout" validate:"min=1s,max=5m"`
	Retries     int           `json:"retries" validate:"min=0,max=5"`
}

// BasicAuth represents username/password authentication.
type BasicAuth struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// TokenAuth represents bearer token authentication.
type TokenAuth struct {
	Token string `json:"token" validate:"required"`
}

// OAuth2Auth represents OAuth2-based authentication.
type OAuth2Auth struct {
	ClientID     string   `json:"client_id" validate:"required"`
	ClientSecret string   `json:"client_secret" validate:"required"`
	TokenURL     string   `json:"token_url" validate:"required,url"`
	Scope        []string `json:"scope,omitempty"`
}

// CertificateInfo contains metadata about an X.509 certificate.
type CertificateInfo struct {
	Subject      string    `json:"subject"`
	Issuer       string    `json:"issuer"`
	SerialNumber string    `json:"serial_number"`
	NotBefore    time.Time `json:"not_before"`
	NotAfter     time.Time `json:"not_after"`
	Fingerprint  string    `json:"fingerprint"`
}

// AuthResult represents the result of an authentication attempt.
type AuthResult struct {
	Success      bool                   `json:"success"`
	Credentials  *Credentials           `json:"credentials,omitempty"`
	Token        *Token                 `json:"token,omitempty"`
	ErrorMessage string                 `json:"error_message,omitempty"`
	ExpiresAt    *time.Time             `json:"expires_at,omitempty"`
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
}

// TLSVersionMap provides human-readable TLS version constants.
var TLSVersionMap = map[string]uint16{
	"TLS10": tls.VersionTLS10,
	"TLS11": tls.VersionTLS11,
	"TLS12": tls.VersionTLS12,
	"TLS13": tls.VersionTLS13,
}

// IsExpired checks if the credentials have expired.
func (c *Credentials) IsExpired() bool {
	if c.ExpiresAt == nil {
		return false
	}
	return time.Now().After(*c.ExpiresAt)
}

// IsExpired checks if the token has expired.
func (t *Token) IsExpired() bool {
	return time.Now().After(t.ExpiresAt)
}

// IsValid checks if the token is valid and not expired.
func (t *Token) IsValid() bool {
	return t.AccessToken != "" && !t.IsExpired()
}

// ToTLSConfig converts TLSConfig to Go's tls.Config.
func (tc *TLSConfig) ToTLSConfig() *tls.Config {
	if tc == nil {
		return &tls.Config{}
	}
	return &tls.Config{
		InsecureSkipVerify: tc.InsecureSkipVerify,
		MinVersion:         tc.MinVersion,
		MaxVersion:         tc.MaxVersion,
		ServerName:         tc.ServerName,
		CipherSuites:       tc.CipherSuites,
	}
}