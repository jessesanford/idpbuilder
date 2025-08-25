package auth

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"time"
)

var (
	// ErrCredentialsNotFound indicates no credentials exist for a registry
	ErrCredentialsNotFound = errors.New("credentials not found for registry")

	// ErrInvalidCredentials indicates credentials are malformed
	ErrInvalidCredentials = errors.New("invalid credentials")

	// ErrExpiredCredentials indicates credentials have expired
	ErrExpiredCredentials = errors.New("credentials have expired")

	// ErrInvalidCertificate indicates certificate is invalid or malformed
	ErrInvalidCertificate = errors.New("invalid certificate")

	// ErrInvalidRegistry indicates registry URL is malformed
	ErrInvalidRegistry = errors.New("invalid registry URL")

)

// hostnamePortRegex validates hostname with optional port
var hostnamePortRegex = regexp.MustCompile(`^[a-zA-Z0-9.-]+(:[0-9]+)?$`)

// validateHostnamePort validates that a string matches the hostname:port pattern
// Supports both hostname and hostname:port formats
func validateHostnamePort(value string) bool {
	if value == "" {
		return false
	}
	return hostnamePortRegex.MatchString(value)
}

// ValidateCredentials performs comprehensive validation of credentials
func ValidateCredentials(creds *Credentials) error {
	if creds == nil {
		return ErrInvalidCredentials
	}

	// Basic registry validation
	if !validateHostnamePort(creds.Registry) {
		return fmt.Errorf("%w: invalid registry format", ErrInvalidCredentials)
	}

	// Check expiration if set
	if creds.ExpiresAt != nil && time.Now().After(*creds.ExpiresAt) {
		return ErrExpiredCredentials
	}

	// Validate authentication method-specific requirements
	switch creds.AuthMethod {
	case AuthMethodBasic:
		if creds.Username == "" || creds.Password == "" {
			return fmt.Errorf("%w: basic auth requires username and password", ErrInvalidCredentials)
		}

	case AuthMethodToken:
		if creds.Token == nil {
			return fmt.Errorf("%w: token auth requires token", ErrInvalidCredentials)
		}
		if err := ValidateToken(creds.Token); err != nil {
			return err
		}

	case AuthMethodOAuth2:
		if creds.Token == nil {
			return fmt.Errorf("%w: oauth2 auth requires token", ErrInvalidCredentials)
		}
		if err := ValidateToken(creds.Token); err != nil {
			return err
		}

	default:
		return fmt.Errorf("%w: unsupported auth method: %s", ErrInvalidCredentials, creds.AuthMethod)
	}

	return nil
}

// ValidateToken validates token structure and expiration
func ValidateToken(token *Token) error {
	if token == nil {
		return fmt.Errorf("%w: token cannot be nil", ErrInvalidCredentials)
	}

	if token.Value == "" {
		return fmt.Errorf("%w: token value cannot be empty", ErrInvalidCredentials)
	}

	if token.Type == "" {
		return fmt.Errorf("%w: token type cannot be empty", ErrInvalidCredentials)
	}

	// Check token expiration
	if time.Now().After(token.ExpiresAt) {
		return ErrExpiredCredentials
	}

	return nil
}

// ValidateCertificate validates certificate data and expiration
func ValidateCertificate(certData []byte) (*CertificateInfo, error) {
	if len(certData) == 0 {
		return nil, fmt.Errorf("%w: certificate data cannot be empty", ErrInvalidCertificate)
	}

	// Decode PEM block
	block, _ := pem.Decode(certData)
	if block == nil {
		return nil, fmt.Errorf("%w: failed to decode PEM data", ErrInvalidCertificate)
	}

	// Parse X.509 certificate
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to parse certificate: %v", ErrInvalidCertificate, err)
	}

	// Check certificate validity period
	now := time.Now()
	if now.Before(cert.NotBefore) {
		return nil, fmt.Errorf("%w: certificate not yet valid", ErrInvalidCertificate)
	}
	if now.After(cert.NotAfter) {
		return nil, fmt.Errorf("%w: certificate has expired", ErrInvalidCertificate)
	}

	// Extract IP addresses as strings
	ipStrings := make([]string, len(cert.IPAddresses))
	for i, ip := range cert.IPAddresses {
		ipStrings[i] = ip.String()
	}

	// Create certificate info
	info := &CertificateInfo{
		Subject:      cert.Subject.String(),
		Issuer:       cert.Issuer.String(),
		SerialNumber: cert.SerialNumber.String(),
		NotBefore:    cert.NotBefore,
		NotAfter:     cert.NotAfter,
		DNSNames:     cert.DNSNames,
		IPAddresses:  ipStrings,
	}

	return info, nil
}

// ValidateRegistryURL validates registry URL format
func ValidateRegistryURL(registry string) error {
	if registry == "" {
		return fmt.Errorf("%w: registry URL cannot be empty", ErrInvalidRegistry)
	}

	// Handle hostname:port format (common for registries)
	if !strings.Contains(registry, "://") {
		// Validate as hostname:port
		if !hostnamePortRegex.MatchString(registry) {
			return fmt.Errorf("%w: invalid hostname:port format: %s", ErrInvalidRegistry, registry)
		}
		return nil
	}

	// Validate as full URL
	parsedURL, err := url.Parse(registry)
	if err != nil {
		return fmt.Errorf("%w: failed to parse URL: %v", ErrInvalidRegistry, err)
	}

	if parsedURL.Scheme == "" {
		return fmt.Errorf("%w: URL must have a scheme", ErrInvalidRegistry)
	}

	if parsedURL.Host == "" {
		return fmt.Errorf("%w: URL must have a host", ErrInvalidRegistry)
	}

	return nil
}