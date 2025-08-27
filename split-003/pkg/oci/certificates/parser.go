package certificates

import (
	"crypto/x509"
	"fmt"
	"time"

	v2 "github.com/cnoe-io/idpbuilder/pkg/oci/api/v2"
)

// ValidateChain validates a complete certificate chain
func ValidateChain(chain []*x509.Certificate) error {
	if len(chain) == 0 {
		return &v2.CertificateError{
			Code:    "EMPTY_CHAIN",
			Message: "certificate chain is empty",
		}
	}

	// Validate each certificate in the chain
	for i, cert := range chain {
		if err := ValidateCertificate(cert); err != nil {
			return &v2.CertificateError{
				Code:    "CHAIN_CERT_INVALID",
				Message: fmt.Sprintf("certificate %d in chain is invalid: %v", i, err),
				Cert:    cert,
				Err:     err,
			}
		}
	}

	// Validate chain relationships
	for i := 0; i < len(chain)-1; i++ {
		current := chain[i]
		issuer := chain[i+1]
		
		if !isSignedBy(current, issuer) {
			return &v2.CertificateError{
				Code:    "BROKEN_CHAIN_LINK",
				Message: fmt.Sprintf("certificate %d is not signed by certificate %d", i, i+1),
				Cert:    current,
			}
		}
	}

	return nil
}

// ValidateCertificate performs comprehensive validation of a single certificate
func ValidateCertificate(cert *x509.Certificate) error {
	if cert == nil {
		return &v2.CertificateError{
			Code:    "NULL_CERTIFICATE",
			Message: "certificate is nil",
		}
	}

	now := time.Now()
	
	// Check certificate validity period
	if cert.NotAfter.Before(now) {
		return &v2.CertificateError{
			Code:    "CERTIFICATE_EXPIRED",
			Message: fmt.Sprintf("certificate expired on %v", cert.NotAfter),
			Cert:    cert,
		}
	}
	
	if cert.NotBefore.After(now) {
		return &v2.CertificateError{
			Code:    "CERTIFICATE_NOT_YET_VALID",
			Message: fmt.Sprintf("certificate not valid until %v", cert.NotBefore),
			Cert:    cert,
		}
	}
	
	// Validate certificate structure
	if cert.Subject.String() == "" {
		return &v2.CertificateError{
			Code:    "EMPTY_SUBJECT",
			Message: "certificate has empty subject",
			Cert:    cert,
		}
	}
	
	// Validate key usage for CA certificates
	if cert.IsCA && cert.KeyUsage&x509.KeyUsageCertSign == 0 {
		return &v2.CertificateError{
			Code:    "INVALID_CA_KEY_USAGE",
			Message: "CA certificate missing certificate signing key usage",
			Cert:    cert,
		}
	}

	return nil
}

// isSignedBy checks if the subject certificate is signed by the issuer certificate
func isSignedBy(subject, issuer *x509.Certificate) bool {
	if !subject.Issuer.Equal(issuer.Subject) {
		return false
	}
	
	err := subject.CheckSignatureFrom(issuer)
	return err == nil
}

// Format parsers and detector (stubs that delegate to Split-002)

// PEMParser handles PEM format certificates
type PEMParser struct{}

func (p *PEMParser) Parse(data []byte) (*v2.CertBundle, error) {
	// In real implementation, this would delegate to Split-002 PEM parser
	return &v2.CertBundle{Format: v2.CertFormatPEM, LoadedAt: time.Now()}, nil
}

// DERParser handles DER format certificates
type DERParser struct{}

func (p *DERParser) Parse(data []byte) (*v2.CertBundle, error) {
	return &v2.CertBundle{Format: v2.CertFormatDER, LoadedAt: time.Now()}, nil
}

// PKCS7Parser handles PKCS7 format certificates
type PKCS7Parser struct{}

func (p *PKCS7Parser) Parse(data []byte) (*v2.CertBundle, error) {
	return &v2.CertBundle{Format: v2.CertFormatPKCS7, LoadedAt: time.Now()}, nil
}

// PKCS12Parser handles PKCS12 format certificates
type PKCS12Parser struct{}

func (p *PKCS12Parser) Parse(data []byte) (*v2.CertBundle, error) {
	return &v2.CertBundle{Format: v2.CertFormatPKCS12, LoadedAt: time.Now()}, nil
}

// MagicBytesDetector detects certificate formats using magic bytes
type MagicBytesDetector struct{}

func (d *MagicBytesDetector) DetectFormat(data []byte) (v2.CertFormat, error) {
	if len(data) == 0 {
		return "", &v2.CertificateError{
			Code:    "EMPTY_DATA",
			Message: "certificate data is empty",
		}
	}
	
	// Check for PEM format
	if len(data) > 10 && string(data[:10]) == "-----BEGIN" {
		return v2.CertFormatPEM, nil
	}
	
	// Check for PKCS12 format
	if len(data) >= 4 && data[0] == 0x30 && data[2] == 0x02 && data[3] == 0x01 {
		return v2.CertFormatPKCS12, nil
	}
	
	// Check for DER format
	if len(data) >= 2 && data[0] == 0x30 {
		return v2.CertFormatDER, nil
	}
	
	return "", &v2.CertificateError{
		Code:    "UNKNOWN_FORMAT",
		Message: "unable to detect certificate format",
	}
}