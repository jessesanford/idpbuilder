package certs

import (
	"crypto/tls"
	"crypto/x509"
	"time"
)

// Certificate represents a wrapper around x509.Certificate with additional metadata.
type Certificate struct {
	// Cert is the underlying X.509 certificate.
	Cert *x509.Certificate

	// PEMData contains the original PEM-encoded certificate data.
	PEMData []byte

	// Fingerprint is the SHA-256 fingerprint of the certificate.
	Fingerprint string

	// Source indicates where this certificate was loaded from.
	Source string
}

// IsExpired checks if the certificate is expired.
func (c *Certificate) IsExpired() bool {
	if c.Cert == nil {
		return true
	}
	return time.Now().After(c.Cert.NotAfter)
}

// IsValidFor checks if the certificate is valid for the given hostname.
func (c *Certificate) IsValidFor(hostname string) error {
	if c.Cert == nil {
		return ErrInvalidCertificate
	}
	return c.Cert.VerifyHostname(hostname)
}

// Subject returns the certificate subject as a string.
func (c *Certificate) Subject() string {
	if c.Cert == nil {
		return ""
	}
	return c.Cert.Subject.String()
}

// CertificateBundle contains a collection of certificates for TLS configuration.
// It includes CA certificates, client certificates, and private keys.
type CertificateBundle struct {
	// CACert is the Certificate Authority certificate used for server verification.
	CACert *Certificate

	// ClientCert is the client certificate for mutual TLS authentication.
	ClientCert *Certificate

	// ClientKey contains the private key corresponding to ClientCert.
	ClientKey []byte

	// IntermediateCerts contains intermediate CA certificates in the chain.
	IntermediateCerts []*Certificate

	// ValidFrom indicates when this certificate bundle becomes valid.
	ValidFrom time.Time

	// ValidUntil indicates when this certificate bundle expires.
	ValidUntil time.Time

	// Registry is the registry URL this bundle applies to.
	Registry string
}

// IsValid checks if the certificate bundle is currently valid.
func (cb *CertificateBundle) IsValid() bool {
	now := time.Now()
	return now.After(cb.ValidFrom) && now.Before(cb.ValidUntil)
}

// IsExpired checks if any certificate in the bundle is expired.
func (cb *CertificateBundle) IsExpired() bool {
	if cb.CACert != nil && cb.CACert.IsExpired() {
		return true
	}
	if cb.ClientCert != nil && cb.ClientCert.IsExpired() {
		return true
	}
	for _, cert := range cb.IntermediateCerts {
		if cert.IsExpired() {
			return true
		}
	}
	return false
}

// TLSConfig represents TLS configuration for registry connections.
// It provides fine-grained control over TLS behavior and certificate validation.
type TLSConfig struct {
	// InsecureSkipVerify controls whether TLS certificate verification is skipped.
	// This should only be used in development or testing environments.
	InsecureSkipVerify bool

	// RootCAs defines the set of root certificate authorities for server verification.
	RootCAs *x509.CertPool

	// ClientCAs defines the set of root certificate authorities for client verification.
	ClientCAs *x509.CertPool

	// Certificates contains the client certificates for mutual TLS.
	Certificates []tls.Certificate

	// MinVersion specifies the minimum TLS version acceptable.
	MinVersion uint16

	// MaxVersion specifies the maximum TLS version acceptable.
	MaxVersion uint16

	// CipherSuites specifies the list of enabled TLS cipher suites.
	CipherSuites []uint16

	// ServerName is used to verify the hostname on the returned certificate.
	ServerName string

	// ClientAuth determines the server's policy for TLS client authentication.
	ClientAuth tls.ClientAuthType

	// Registry is the registry URL this TLS configuration applies to.
	Registry string
}

// ToStandardTLSConfig converts the TLSConfig to Go's standard tls.Config.
func (tc *TLSConfig) ToStandardTLSConfig() *tls.Config {
	return &tls.Config{
		InsecureSkipVerify: tc.InsecureSkipVerify,
		RootCAs:            tc.RootCAs,
		ClientCAs:          tc.ClientCAs,
		Certificates:       tc.Certificates,
		MinVersion:         tc.MinVersion,
		MaxVersion:         tc.MaxVersion,
		CipherSuites:       tc.CipherSuites,
		ServerName:         tc.ServerName,
		ClientAuth:         tc.ClientAuth,
	}
}

// CertificateValidator defines the interface for certificate validation operations.
type CertificateValidator interface {
	// Validate performs comprehensive validation of a certificate bundle.
	Validate(bundle *CertificateBundle) error

	// ValidateChain validates the certificate chain against root CAs.
	ValidateChain(certs []*x509.Certificate, roots *x509.CertPool) error

	// ValidateHostname checks if the certificate is valid for the given hostname.
	ValidateHostname(cert *x509.Certificate, hostname string) error

	// CheckExpiry verifies that certificates are not expired and warns of upcoming expiration.
	CheckExpiry(bundle *CertificateBundle, warningPeriod time.Duration) error
}

// CertificateStore defines the interface for certificate storage and retrieval.
type CertificateStore interface {
	// Get retrieves the certificate bundle for the specified registry.
	Get(registry string) (*CertificateBundle, error)

	// Set stores a certificate bundle for the specified registry.
	Set(registry string, bundle *CertificateBundle) error

	// Delete removes the certificate bundle for the specified registry.
	Delete(registry string) error

	// List returns all registries with stored certificate bundles.
	List() ([]string, error)

	// Refresh reloads certificate bundles from their sources.
	Refresh() error
}