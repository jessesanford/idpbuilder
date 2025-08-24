package certs

import (
	"crypto/tls"
	"errors"
	"time"
)

// Default paths for certificate and key files.
const (
	// DefaultCACertPath is the default path for CA certificates.
	DefaultCACertPath = "/etc/ssl/certs/ca-certificates.crt"

	// DefaultClientCertPath is the default path for client certificates.
	DefaultClientCertPath = "~/.docker/certs.d"

	// DefaultSystemCertPath is the system certificate store path on Linux.
	DefaultSystemCertPath = "/etc/ssl/certs"

	// DefaultUserCertPath is the user certificate store path.
	DefaultUserCertPath = "~/.certs"
)

// Certificate file extensions and formats.
const (
	// CertExtensionPEM is the file extension for PEM-encoded certificates.
	CertExtensionPEM = ".pem"

	// CertExtensionCRT is the file extension for certificate files.
	CertExtensionCRT = ".crt"

	// CertExtensionCER is the file extension for certificate files.
	CertExtensionCER = ".cer"

	// KeyExtensionKey is the file extension for private key files.
	KeyExtensionKey = ".key"

	// CertFilePattern is the glob pattern for certificate files.
	CertFilePattern = "*.{pem,crt,cer}"

	// KeyFilePattern is the glob pattern for key files.
	KeyFilePattern = "*.key"
)

// TLS version constants.
const (
	// MinTLSVersion is the minimum acceptable TLS version.
	MinTLSVersion = tls.VersionTLS12

	// MaxTLSVersion is the maximum TLS version to use.
	MaxTLSVersion = tls.VersionTLS13

	// DefaultTLSVersion is the default TLS version preference.
	DefaultTLSVersion = tls.VersionTLS13
)

// Certificate validation periods and thresholds.
const (
	// CertExpiryWarningPeriod is the period before expiry to start warning.
	CertExpiryWarningPeriod = 30 * 24 * time.Hour // 30 days

	// CertRenewalThreshold is when to automatically renew certificates.
	CertRenewalThreshold = 7 * 24 * time.Hour // 7 days

	// CertValidationInterval is how often to check certificate validity.
	CertValidationInterval = 24 * time.Hour // 24 hours

	// DefaultCertValidityPeriod is the default certificate validity period.
	DefaultCertValidityPeriod = 365 * 24 * time.Hour // 1 year
)

// Certificate type identifiers.
const (
	// CertTypeCA identifies Certificate Authority certificates.
	CertTypeCA = "ca"

	// CertTypeClient identifies client certificates for mutual TLS.
	CertTypeClient = "client"

	// CertTypeServer identifies server certificates.
	CertTypeServer = "server"

	// CertTypeIntermediate identifies intermediate CA certificates.
	CertTypeIntermediate = "intermediate"

	// CertTypeSelfSigned identifies self-signed certificates.
	CertTypeSelfSigned = "self-signed"
)

// Certificate extension OIDs (Object Identifiers).
const (
	// OIDKeyUsage is the OID for Key Usage extension.
	OIDKeyUsage = "2.5.29.15"

	// OIDExtendedKeyUsage is the OID for Extended Key Usage extension.
	OIDExtendedKeyUsage = "2.5.29.37"

	// OIDSubjectAltName is the OID for Subject Alternative Name extension.
	OIDSubjectAltName = "2.5.29.17"

	// OIDBasicConstraints is the OID for Basic Constraints extension.
	OIDBasicConstraints = "2.5.29.19"
)

// Common certificate validation error variables.
var (
	// ErrInvalidCertificate indicates that the certificate is invalid or malformed.
	ErrInvalidCertificate = errors.New("invalid or malformed certificate")

	// ErrCertificateExpired indicates that the certificate has expired.
	ErrCertificateExpired = errors.New("certificate has expired")

	// ErrCertificateNotYetValid indicates that the certificate is not yet valid.
	ErrCertificateNotYetValid = errors.New("certificate is not yet valid")

	// ErrUntrustedCertificate indicates that the certificate is not trusted.
	ErrUntrustedCertificate = errors.New("certificate is not trusted")

	// ErrHostnameMismatch indicates that the certificate hostname doesn't match.
	ErrHostnameMismatch = errors.New("certificate hostname mismatch")

	// ErrInvalidCertificateChain indicates that the certificate chain is invalid.
	ErrInvalidCertificateChain = errors.New("invalid certificate chain")

	// ErrCertificateNotFound indicates that no certificate was found.
	ErrCertificateNotFound = errors.New("certificate not found")

	// ErrPrivateKeyNotFound indicates that the private key was not found.
	ErrPrivateKeyNotFound = errors.New("private key not found")

	// ErrCertificateKeyMismatch indicates that the certificate and key don't match.
	ErrCertificateKeyMismatch = errors.New("certificate and private key do not match")

	// ErrUnsupportedCertificateFormat indicates an unsupported certificate format.
	ErrUnsupportedCertificateFormat = errors.New("unsupported certificate format")
)