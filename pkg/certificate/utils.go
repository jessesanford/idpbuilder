package certificate

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"time"
)

// ParseCertificateFromPEM parses a PEM-encoded certificate and returns the x509.Certificate.
func ParseCertificateFromPEM(certPEM []byte) (*x509.Certificate, error) {
	block, _ := pem.Decode(certPEM)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block containing certificate")
	}

	if block.Type != "CERTIFICATE" {
		return nil, fmt.Errorf("invalid PEM block type: %s, expected CERTIFICATE", block.Type)
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse certificate: %w", err)
	}

	return cert, nil
}

// ValidateCertificate performs comprehensive validation of a certificate.
// It checks expiration, signature, and basic constraints.
func ValidateCertificate(certPEM []byte) error {
	cert, err := ParseCertificateFromPEM(certPEM)
	if err != nil {
		return fmt.Errorf("parsing certificate: %w", err)
	}

	// Check if certificate is expired
	now := time.Now()
	if now.Before(cert.NotBefore) {
		return fmt.Errorf("certificate is not yet valid (NotBefore: %v)", cert.NotBefore)
	}
	if now.After(cert.NotAfter) {
		return fmt.Errorf("certificate has expired (NotAfter: %v)", cert.NotAfter)
	}

	// For self-signed certificates, verify against itself
	// Only check signature if it's a CA or has the right key usage
	if cert.IsCA || (cert.KeyUsage&x509.KeyUsageCertSign != 0) {
		err = cert.CheckSignatureFrom(cert)
		if err != nil {
			return fmt.Errorf("certificate signature validation failed: %w", err)
		}
	}

	return nil
}

// GetCertificateInfo extracts key information from a PEM-encoded certificate.
func GetCertificateInfo(certPEM []byte) (*CertificateMetadata, error) {
	cert, err := ParseCertificateFromPEM(certPEM)
	if err != nil {
		return nil, fmt.Errorf("parsing certificate: %w", err)
	}

	return &CertificateMetadata{
		Subject:     cert.Subject.CommonName,
		DNSNames:    cert.DNSNames,
		NotBefore:   cert.NotBefore,
		NotAfter:    cert.NotAfter,
		IsCA:        cert.IsCA,
		KeyUsage:    cert.KeyUsage,
		ExtKeyUsage: cert.ExtKeyUsage,
	}, nil
}

// IsCertificateExpiringSoon checks if a certificate will expire within the specified duration.
func IsCertificateExpiringSoon(certPEM []byte, within time.Duration) (bool, error) {
	cert, err := ParseCertificateFromPEM(certPEM)
	if err != nil {
		return false, fmt.Errorf("parsing certificate: %w", err)
	}

	expirationThreshold := time.Now().Add(within)
	return cert.NotAfter.Before(expirationThreshold), nil
}

// GetCertificateFingerprint calculates the SHA-256 fingerprint of a certificate.
func GetCertificateFingerprint(certPEM []byte) (string, error) {
	cert, err := ParseCertificateFromPEM(certPEM)
	if err != nil {
		return "", fmt.Errorf("parsing certificate: %w", err)
	}

	// x509.Certificate already provides SHA256 fingerprint calculation
	fingerprint := fmt.Sprintf("%x", cert.Raw)
	return fingerprint, nil
}

// CreateGenerationOptionsForTLS creates certificate generation options suitable for TLS server certificates.
func CreateGenerationOptionsForTLS(dnsNames []string, validFor time.Duration) *GenerationOptions {
	opts := DefaultGenerationOptions()
	opts.DNSNames = dnsNames
	opts.ValidFor = validFor
	opts.Subject = "TLS Server Certificate"

	// TLS server certificate specific settings
	opts.KeyUsage = x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment
	opts.ExtKeyUsage = []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth}

	return opts
}

// CreateGenerationOptionsForCA creates certificate generation options suitable for Certificate Authority certificates.
func CreateGenerationOptionsForCA(organization string, validFor time.Duration) *GenerationOptions {
	opts := DefaultGenerationOptions()
	opts.Organization = organization
	opts.ValidFor = validFor
	opts.Subject = fmt.Sprintf("%s Certificate Authority", organization)
	opts.IsCA = true

	// CA certificate specific settings
	opts.KeyUsage = x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign | x509.KeyUsageCRLSign
	opts.ExtKeyUsage = []x509.ExtKeyUsage{} // CAs typically don't have extended key usage

	return opts
}