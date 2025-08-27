package certificates

import (
	"context"
	"crypto/x509"
)

// Loader defines the interface for loading certificate bundles
type Loader interface {
	// LoadBundle loads a certificate bundle from the specified path
	LoadBundle(ctx context.Context, path string) (*CertBundle, error)
	// LoadPEM loads certificates from PEM format data
	LoadPEM(ctx context.Context, data []byte) (*CertBundle, error)
	// LoadDER loads certificates from DER format data
	LoadDER(ctx context.Context, data []byte) (*CertBundle, error)
	// LoadPKCS7 loads certificates from PKCS7 format data
	LoadPKCS7(ctx context.Context, data []byte) (*CertBundle, error)
	// LoadPKCS12 loads certificates from PKCS12 format data with password
	LoadPKCS12(ctx context.Context, data []byte, password string) (*CertBundle, error)
}

// Parser defines the interface for parsing certificates and chains
type Parser interface {
	// ParseCertificateChain parses and validates a certificate chain
	ParseCertificateChain(certs []*x509.Certificate) ([]*x509.Certificate, error)
	// ValidateCertificate performs comprehensive certificate validation
	ValidateCertificate(cert *x509.Certificate) error
	// ConvertToBundle creates a CertBundle from parsed certificates
	ConvertToBundle(certs []*x509.Certificate, format CertFormat, source string) *CertBundle
}

// FormatParser defines the interface for format-specific parsers
type FormatParser interface {
	// Parse parses certificate data in a specific format
	Parse(data []byte) (*CertBundle, error)
	// Validate validates that data is in the expected format
	Validate(data []byte) error
}

// FormatDetector defines the interface for auto-detecting certificate formats
type FormatDetector interface {
	// DetectFormat auto-detects certificate format from file content
	DetectFormat(data []byte) (CertFormat, error)
}

// Validator defines the interface for certificate validation
type Validator interface {
	// ValidateBundle performs comprehensive validation of a certificate bundle
	ValidateBundle(bundle *CertBundle) error
	// ValidateCertificate validates a single certificate
	ValidateCertificate(cert *x509.Certificate) error
}