package certificates

import (
	"crypto/x509"
	"time"
)

// CertFormat represents the format of a certificate
type CertFormat string

const (
	// CertFormatPEM represents PEM format certificates
	CertFormatPEM CertFormat = "PEM"
	// CertFormatDER represents DER format certificates  
	CertFormatDER CertFormat = "DER"
	// CertFormatPKCS7 represents PKCS#7 format certificates
	CertFormatPKCS7 CertFormat = "PKCS7"
	// CertFormatPKCS12 represents PKCS#12 format certificates
	CertFormatPKCS12 CertFormat = "PKCS12"
)

// ValidationStatus represents the validation status of a certificate
type ValidationStatus string

const (
	// ValidationStatusValid indicates the certificate is valid
	ValidationStatusValid ValidationStatus = "valid"
	// ValidationStatusExpired indicates the certificate is expired
	ValidationStatusExpired ValidationStatus = "expired"
	// ValidationStatusNotYetValid indicates the certificate is not yet valid
	ValidationStatusNotYetValid ValidationStatus = "not_yet_valid"
	// ValidationStatusInvalid indicates the certificate is invalid
	ValidationStatusInvalid ValidationStatus = "invalid"
)

// CertBundle represents a collection of certificates and CA certificates
type CertBundle struct {
	// Certificates contains the end-entity certificates
	Certificates []*x509.Certificate `json:"certificates"`
	// CAs contains the Certificate Authority certificates
	CAs []*x509.Certificate `json:"cas"`
	// Format indicates the original format of the certificate bundle
	Format CertFormat `json:"format"`
	// Source contains information about where the bundle was loaded from
	Source string `json:"source"`
	// LoadedAt indicates when this bundle was loaded
	LoadedAt time.Time `json:"loaded_at"`
}

// Certificate represents an enhanced x509 certificate with additional metadata
type Certificate struct {
	// Certificate is the underlying X.509 certificate
	*x509.Certificate
	// ValidationStatus indicates the current validation status
	ValidationStatus ValidationStatus `json:"validation_status"`
	// IsCA indicates if this is a Certificate Authority certificate
	IsCA bool `json:"is_ca"`
	// IsSelfSigned indicates if this certificate is self-signed
	IsSelfSigned bool `json:"is_self_signed"`
}

// LoaderConfig represents configuration for certificate loaders
type LoaderConfig struct {
	// StrictMode enables strict certificate validation
	StrictMode bool `json:"strict_mode"`
	// MaxChainDepth sets the maximum allowed certificate chain depth
	MaxChainDepth int `json:"max_chain_depth"`
	// AllowExpired allows loading of expired certificates
	AllowExpired bool `json:"allow_expired"`
}

// NewCertBundle creates a new certificate bundle
func NewCertBundle(format CertFormat) *CertBundle {
	return &CertBundle{
		Certificates: make([]*x509.Certificate, 0),
		CAs:          make([]*x509.Certificate, 0),
		Format:       format,
		LoadedAt:     time.Now(),
	}
}

// AddCertificate adds a certificate to the appropriate collection in the bundle
func (b *CertBundle) AddCertificate(cert *x509.Certificate) {
	if cert.IsCA {
		b.CAs = append(b.CAs, cert)
	} else {
		b.Certificates = append(b.Certificates, cert)
	}
}

// IsEmpty returns true if the bundle contains no certificates
func (b *CertBundle) IsEmpty() bool {
	return len(b.Certificates) == 0 && len(b.CAs) == 0
}