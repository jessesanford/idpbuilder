package v2

import (
	"crypto/x509"
	"time"
)

// CertFormat represents supported certificate formats
type CertFormat string

const (
	CertFormatPEM    CertFormat = "PEM"
	CertFormatDER    CertFormat = "DER"
	CertFormatPKCS7  CertFormat = "PKCS7"
	CertFormatPKCS12 CertFormat = "PKCS12"
)

// CertBundle represents a collection of certificates
type CertBundle struct {
	Certificates []*x509.Certificate `json:"certificates"`
	CAs          []*x509.Certificate `json:"cas"`
	Format       CertFormat          `json:"format"`
	LoadedAt     time.Time           `json:"loaded_at"`
	Source       string              `json:"source"`
}

// CertificateError represents certificate-specific errors
type CertificateError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Err     error  `json:"-"`
}

func (e *CertificateError) Error() string {
	if e.Err != nil {
		return e.Message + ": " + e.Err.Error()
	}
	return e.Message
}

// CertificateParser interface for certificate parsing utilities
type CertificateParser interface {
	ConvertToBundle(certs []*x509.Certificate, format CertFormat, source string) *CertBundle
}