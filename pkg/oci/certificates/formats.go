package certificates

import (
	"crypto/x509"
	"encoding/asn1"
	"encoding/pem"
	"errors"
	"fmt"

	v2 "github.com/cnoe-io/idpbuilder/pkg/oci/api/v2"
	"golang.org/x/crypto/pkcs12"
)

// PEMParser handles PEM format certificates
type PEMParser struct {
	parser *CertificateParser
}

// Parse parses PEM format certificate data
func (p *PEMParser) Parse(data []byte) (*v2.CertBundle, error) {
	var certs []*x509.Certificate

	remaining := data
	for len(remaining) > 0 {
		block, rest := pem.Decode(remaining)
		if block == nil {
			break
		}

		if block.Type == "CERTIFICATE" {
			cert, err := x509.ParseCertificate(block.Bytes)
			if err != nil {
				return nil, &v2.CertificateError{
					Code:    "PEM_PARSE_ERROR",
					Message: fmt.Sprintf("failed to parse PEM certificate: %v", err),
					Err:     err,
				}
			}
			certs = append(certs, cert)
		}

		remaining = rest
	}

	if len(certs) == 0 {
		return nil, &v2.CertificateError{
			Code:    "NO_CERTIFICATES_IN_PEM",
			Message: "no valid certificates found in PEM data",
		}
	}

	return p.parser.ConvertToBundle(certs, v2.CertFormatPEM, ""), nil
}

// Validate validates PEM format data
func (p *PEMParser) Validate(data []byte) error {
	if len(data) == 0 {
		return errors.New("PEM data is empty")
	}

	// Check for PEM header
	if !containsPEMHeader(data) {
		return errors.New("data does not contain valid PEM headers")
	}

	// Try to decode at least one PEM block
	block, _ := pem.Decode(data)
	if block == nil {
		return errors.New("failed to decode any PEM blocks")
	}

	return nil
}

// DERParser handles DER format certificates
type DERParser struct {
	parser *CertificateParser
}

// Parse parses DER format certificate data
func (p *DERParser) Parse(data []byte) (*v2.CertBundle, error) {
	if len(data) == 0 {
		return nil, &v2.CertificateError{
			Code:    "EMPTY_DER_DATA",
			Message: "DER data is empty",
		}
	}

	cert, err := x509.ParseCertificate(data)
	if err != nil {
		return nil, &v2.CertificateError{
			Code:    "DER_PARSE_ERROR",
			Message: fmt.Sprintf("failed to parse DER certificate: %v", err),
			Err:     err,
		}
	}

	certs := []*x509.Certificate{cert}
	return p.parser.ConvertToBundle(certs, v2.CertFormatDER, ""), nil
}

// Validate validates DER format data
func (p *DERParser) Validate(data []byte) error {
	if len(data) == 0 {
		return errors.New("DER data is empty")
	}

	// Check for ASN.1 sequence header (DER certificates start with 0x30)
	if len(data) < 2 || data[0] != 0x30 {
		return errors.New("data does not start with ASN.1 SEQUENCE tag")
	}

	return nil
}

// PKCS7Parser handles PKCS7 format certificates
type PKCS7Parser struct {
	parser *CertificateParser
}

// Parse parses PKCS7 format certificate data
func (p *PKCS7Parser) Parse(data []byte) (*v2.CertBundle, error) {
	if len(data) == 0 {
		return nil, &v2.CertificateError{
			Code:    "EMPTY_PKCS7_DATA",
			Message: "PKCS7 data is empty",
		}
	}

	// Parse PKCS7 structure
	var pkcs7 struct {
		ContentType asn1.ObjectIdentifier
		Content     asn1.RawValue `asn1:"explicit,tag:0"`
	}

	_, err := asn1.Unmarshal(data, &pkcs7)
	if err != nil {
		return nil, &v2.CertificateError{
			Code:    "PKCS7_PARSE_ERROR",
			Message: fmt.Sprintf("failed to parse PKCS7 structure: %v", err),
			Err:     err,
		}
	}

	// Extract certificates from PKCS7 content
	certs, err := p.extractCertificatesFromPKCS7(pkcs7.Content.Bytes)
	if err != nil {
		return nil, &v2.CertificateError{
			Code:    "PKCS7_CERT_EXTRACTION_ERROR",
			Message: fmt.Sprintf("failed to extract certificates from PKCS7: %v", err),
			Err:     err,
		}
	}

	if len(certs) == 0 {
		return nil, &v2.CertificateError{
			Code:    "NO_CERTIFICATES_IN_PKCS7",
			Message: "no certificates found in PKCS7 data",
		}
	}

	return p.parser.ConvertToBundle(certs, v2.CertFormatPKCS7, ""), nil
}

// Validate validates PKCS7 format data
func (p *PKCS7Parser) Validate(data []byte) error {
	if len(data) == 0 {
		return errors.New("PKCS7 data is empty")
	}

	// Basic ASN.1 structure validation
	var pkcs7 struct {
		ContentType asn1.ObjectIdentifier
		Content     asn1.RawValue `asn1:"explicit,tag:0"`
	}

	_, err := asn1.Unmarshal(data, &pkcs7)
	if err != nil {
		return fmt.Errorf("invalid PKCS7 structure: %v", err)
	}

	return nil
}

// extractCertificatesFromPKCS7 extracts certificates from PKCS7 content
func (p *PKCS7Parser) extractCertificatesFromPKCS7(data []byte) ([]*x509.Certificate, error) {
	// This is a simplified PKCS7 certificate extraction
	// In a production environment, you might want to use a specialized PKCS7 library

	var certs []*x509.Certificate

	// Parse the content as a sequence of certificates
	var certSequence []asn1.RawValue
	_, err := asn1.Unmarshal(data, &certSequence)
	if err != nil {
		return nil, err
	}

	for _, rawCert := range certSequence {
		cert, err := x509.ParseCertificate(rawCert.Bytes)
		if err != nil {
			// Skip invalid certificates but continue
			continue
		}
		certs = append(certs, cert)
	}

	return certs, nil
}

// PKCS12Parser handles PKCS12 format certificates
type PKCS12Parser struct {
	parser *CertificateParser
}

// Parse parses PKCS12 format certificate data (without password)
func (p *PKCS12Parser) Parse(data []byte) (*v2.CertBundle, error) {
	// Try with empty password first
	return p.ParseWithPassword(data, "")
}

// ParseWithPassword parses PKCS12 format certificate data with password
func (p *PKCS12Parser) ParseWithPassword(data []byte, password string) (*v2.CertBundle, error) {
	if len(data) == 0 {
		return nil, &v2.CertificateError{
			Code:    "EMPTY_PKCS12_DATA",
			Message: "PKCS12 data is empty",
		}
	}

	// Decode PKCS12 data
	privateKey, cert, caCerts, err := pkcs12.DecodeChain(data, password)
	if err != nil {
		return nil, &v2.CertificateError{
			Code:    "PKCS12_DECODE_ERROR",
			Message: fmt.Sprintf("failed to decode PKCS12 data: %v", err),
			Err:     err,
		}
	}

	var allCerts []*x509.Certificate

	// Add the main certificate
	if cert != nil {
		allCerts = append(allCerts, cert)
	}

	// Add CA certificates
	if caCerts != nil {
		allCerts = append(allCerts, caCerts...)
	}

	if len(allCerts) == 0 {
		return nil, &v2.CertificateError{
			Code:    "NO_CERTIFICATES_IN_PKCS12",
			Message: "no certificates found in PKCS12 data",
		}
	}

	// Note: We're not storing the private key in the bundle for security
	// If needed, this could be extended to handle private keys separately
	_ = privateKey

	return p.parser.ConvertToBundle(allCerts, v2.CertFormatPKCS12, ""), nil
}

// Validate validates PKCS12 format data
func (p *PKCS12Parser) Validate(data []byte) error {
	if len(data) == 0 {
		return errors.New("PKCS12 data is empty")
	}

	// PKCS12 files typically start with specific ASN.1 structure
	// Basic validation - try to parse the outer structure
	var pfx struct {
		Version  int
		AuthSafe asn1.RawValue
		MacData  asn1.RawValue `asn1:"optional"`
	}

	_, err := asn1.Unmarshal(data, &pfx)
	if err != nil {
		return fmt.Errorf("invalid PKCS12 structure: %v", err)
	}

	return nil
}

// MagicBytesDetector uses magic bytes for format detection
type MagicBytesDetector struct{}

// DetectFormat auto-detects certificate format from file content
func (d *MagicBytesDetector) DetectFormat(data []byte) (v2.CertFormat, error) {
	if len(data) == 0 {
		return "", &v2.CertificateError{
			Code:    "EMPTY_DATA",
			Message: "certificate data is empty",
		}
	}

	// Check for PEM format (starts with "-----BEGIN")
	if containsPEMHeader(data) {
		return v2.CertFormatPEM, nil
	}

	// Check for PKCS12 format (specific magic bytes)
	if len(data) >= 4 && isPKCS12Magic(data[:4]) {
		return v2.CertFormatPKCS12, nil
	}

	// Check for DER format (ASN.1 sequence)
	if len(data) >= 2 && data[0] == 0x30 && data[1] > 0 {
		// Could be DER or PKCS7, need to distinguish
		if isPKCS7Structure(data) {
			return v2.CertFormatPKCS7, nil
		}
		return v2.CertFormatDER, nil
	}

	return "", &v2.CertificateError{
		Code:    "UNKNOWN_FORMAT",
		Message: "unable to detect certificate format from data",
	}
}

// containsPEMHeader checks if data contains PEM headers
func containsPEMHeader(data []byte) bool {
	pemHeaders := []string{
		"-----BEGIN CERTIFICATE-----",
		"-----BEGIN X509 CERTIFICATE-----",
		"-----BEGIN TRUSTED CERTIFICATE-----",
	}

	dataStr := string(data[:min(len(data), 100)])
	for _, header := range pemHeaders {
		if len(dataStr) >= len(header) && dataStr[:len(header)] == header {
			return true
		}
	}

	return false
}

// isPKCS12Magic checks PKCS12 magic bytes
func isPKCS12Magic(data []byte) bool {
	// PKCS12 files start with ASN.1 sequence with specific version
	if len(data) < 4 {
		return false
	}

	// Look for PKCS12 ASN.1 structure pattern
	return data[0] == 0x30 && data[2] == 0x02 && data[3] == 0x01
}

// isPKCS7Structure checks if the ASN.1 structure indicates PKCS7
func isPKCS7Structure(data []byte) bool {
	// This is a heuristic - PKCS7 structures tend to be more complex
	// than simple DER certificates and contain specific OIDs

	// Look for PKCS7 OID patterns in the first few bytes
	// This is simplified - a full implementation would parse the ASN.1
	pkcs7Patterns := [][]byte{
		{0x06, 0x09, 0x2A, 0x86, 0x48, 0x86, 0xF7, 0x0D, 0x01}, // PKCS#7 OID prefix
	}

	for _, pattern := range pkcs7Patterns {
		if containsBytes(data[:min(len(data), 50)], pattern) {
			return true
		}
	}

	return false
}

// containsBytes checks if data contains pattern
func containsBytes(data, pattern []byte) bool {
	if len(pattern) > len(data) {
		return false
	}

	for i := 0; i <= len(data)-len(pattern); i++ {
		match := true
		for j := 0; j < len(pattern); j++ {
			if data[i+j] != pattern[j] {
				match = false
				break
			}
		}
		if match {
			return true
		}
	}

	return false
}

// ConvertFormat converts certificates between formats
func ConvertFormat(bundle *v2.CertBundle, targetFormat v2.CertFormat) (*v2.CertBundle, error) {
	if bundle == nil {
		return nil, &v2.CertificateError{
			Code:    "NULL_BUNDLE",
			Message: "certificate bundle is nil",
		}
	}

	if bundle.Format == targetFormat {
		// Already in target format, return copy
		newBundle := *bundle
		return &newBundle, nil
	}

	// Create new bundle with target format
	newBundle := &v2.CertBundle{
		Certificates: bundle.Certificates,
		CAs:          bundle.CAs,
		Format:       targetFormat,
		LoadedAt:     bundle.LoadedAt,
		Source:       bundle.Source,
	}

	return newBundle, nil
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
