// Package certs provides trust store utilities and helpers
package certs

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// CertificateInfo contains metadata about an extracted certificate
// This type matches the CertificateInfo from E1.1.1 (kind-certificate-extraction)
// to ensure compatibility during integration
type CertificateInfo struct {
	Subject   string
	Issuer    string
	NotBefore time.Time
	NotAfter  time.Time
	IsCA      bool
	DNSNames  []string
}

// TrustStoreUtils provides utility functions for trust store operations
type TrustStoreUtils struct{}


// NewTrustStoreUtils creates a new instance of trust store utilities
func NewTrustStoreUtils() *TrustStoreUtils {
	return &TrustStoreUtils{}
}

// LoadCertificateFromPEM loads a certificate from PEM-encoded data
func (u *TrustStoreUtils) LoadCertificateFromPEM(pemData []byte) (*x509.Certificate, error) {
	block, _ := pem.Decode(pemData)
	if block == nil || block.Type != "CERTIFICATE" {
		return nil, fmt.Errorf("failed to decode PEM certificate data")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse certificate: %w", err)
	}

	return cert, nil
}

// LoadCertificatesFromPEM loads multiple certificates from PEM-encoded data
func (u *TrustStoreUtils) LoadCertificatesFromPEM(pemData []byte) ([]*x509.Certificate, error) {
	var certs []*x509.Certificate
	
	for len(pemData) > 0 {
		block, rest := pem.Decode(pemData)
		if block == nil {
			break
		}

		if block.Type != "CERTIFICATE" {
			pemData = rest
			continue
		}

		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("failed to parse certificate: %w", err)
		}

		certs = append(certs, cert)
		pemData = rest
	}

	if len(certs) == 0 {
		return nil, fmt.Errorf("no valid certificates found in PEM data")
	}

	return certs, nil
}

// LoadCertificateFromFile loads a certificate from a PEM file
func (u *TrustStoreUtils) LoadCertificateFromFile(filename string) (*x509.Certificate, error) {
	pemData, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read certificate file %s: %w", filename, err)
	}

	return u.LoadCertificateFromPEM(pemData)
}

// LoadCertificatesFromFile loads multiple certificates from a PEM file
func (u *TrustStoreUtils) LoadCertificatesFromFile(filename string) ([]*x509.Certificate, error) {
	pemData, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read certificate file %s: %w", filename, err)
	}

	return u.LoadCertificatesFromPEM(pemData)
}

// CertificateToPEM converts a certificate to PEM format
func (u *TrustStoreUtils) CertificateToPEM(cert *x509.Certificate) ([]byte, error) {
	if cert == nil {
		return nil, fmt.Errorf("certificate cannot be nil")
	}

	pemBlock := &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: cert.Raw,
	}

	return pem.EncodeToMemory(pemBlock), nil
}

// ValidateCertificate performs basic validation on a certificate
func (u *TrustStoreUtils) ValidateCertificate(cert *x509.Certificate) error {
	if cert == nil {
		return fmt.Errorf("certificate cannot be nil")
	}

	// Check if certificate has expired
	if cert.NotAfter.Before(cert.NotBefore) {
		return fmt.Errorf("certificate has invalid date range: NotAfter (%s) is before NotBefore (%s)",
			cert.NotAfter, cert.NotBefore)
	}

	// Check if certificate is self-signed
	if cert.Subject.String() == cert.Issuer.String() {
		// This is okay for self-signed certificates, just note it
		// We don't return an error here as self-signed certs are valid for our use case
	}

	return nil
}

// GetCertificateInfo returns human-readable information about a certificate
// This function returns the CertificateInfo type that matches E1.1.1's structure
func (u *TrustStoreUtils) GetCertificateInfo(cert *x509.Certificate) *CertificateInfo {
	if cert == nil {
		// Return a CertificateInfo with error information in Subject field
		return &CertificateInfo{
			Subject: "Error: certificate is nil",
		}
	}

	info := &CertificateInfo{
		Subject:   cert.Subject.String(),
		Issuer:    cert.Issuer.String(),
		NotBefore: cert.NotBefore,
		NotAfter:  cert.NotAfter,
		DNSNames:  cert.DNSNames,
		IsCA:      cert.IsCA,
	}

	return info
}

// DiscoverCertificateFiles finds all certificate files in a directory
func (u *TrustStoreUtils) DiscoverCertificateFiles(dir string) ([]string, error) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return nil, fmt.Errorf("directory does not exist: %s", dir)
	}

	var certFiles []string
	
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			ext := strings.ToLower(filepath.Ext(path))
			if ext == ".pem" || ext == ".crt" || ext == ".cer" {
				certFiles = append(certFiles, path)
			}
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to walk directory %s: %w", dir, err)
	}

	return certFiles, nil
}