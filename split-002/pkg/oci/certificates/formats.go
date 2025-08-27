package certificates

import (
	"crypto/x509"
	"encoding/asn1"
	"encoding/pem"
	"errors"
	"fmt"

	"golang.org/x/crypto/pkcs12"
	v2 "github.com/cnoe-io/idpbuilder/split-002/pkg/oci/api/v2"
)

// DetectFormat auto-detects certificate format from file content
func DetectFormat(data []byte) (v2.CertFormat, error) {
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

// ParsePEM parses PEM format certificate data
func ParsePEM(data []byte, parser v2.CertificateParser) (*v2.CertBundle, error) {
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
	
	return parser.ConvertToBundle(certs, v2.CertFormatPEM, ""), nil
}

// ValidatePEM validates PEM format data
func ValidatePEM(data []byte) error {
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

// ParseDER parses DER format certificate data
func ParseDER(data []byte, parser v2.CertificateParser) (*v2.CertBundle, error) {
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
	return parser.ConvertToBundle(certs, v2.CertFormatDER, ""), nil
}

// ValidateDER validates DER format data
func ValidateDER(data []byte) error {
	if len(data) == 0 {
		return errors.New("DER data is empty")
	}
	
	// Check for ASN.1 sequence header (DER certificates start with 0x30)
	if len(data) < 2 || data[0] != 0x30 {
		return errors.New("data does not start with ASN.1 SEQUENCE tag")
	}
	
	return nil
}

// ParsePKCS7 parses PKCS7 format certificate data
func ParsePKCS7(data []byte, parser v2.CertificateParser) (*v2.CertBundle, error) {
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
	certs, err := extractCertificatesFromPKCS7(pkcs7.Content.Bytes)
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
	
	return parser.ConvertToBundle(certs, v2.CertFormatPKCS7, ""), nil
}

// ValidatePKCS7 validates PKCS7 format data
func ValidatePKCS7(data []byte) error {
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

func extractCertificatesFromPKCS7(data []byte) ([]*x509.Certificate, error) {
	var certs []*x509.Certificate
	var certSequence []asn1.RawValue
	_, err := asn1.Unmarshal(data, &certSequence)
	if err != nil {
		return nil, err
	}
	
	for _, rawCert := range certSequence {
		cert, err := x509.ParseCertificate(rawCert.Bytes)
		if err != nil {
			continue
		}
		certs = append(certs, cert)
	}
	
	return certs, nil
}

// ParsePKCS12 parses PKCS12 format certificate data (without password)
func ParsePKCS12(data []byte, parser v2.CertificateParser) (*v2.CertBundle, error) {
	// Try with empty password first
	return ParsePKCS12WithPassword(data, "", parser)
}

func ParsePKCS12WithPassword(data []byte, password string, parser v2.CertificateParser) (*v2.CertBundle, error) {
	if len(data) == 0 {
		return nil, &v2.CertificateError{
			Code:    "EMPTY_PKCS12_DATA",
			Message: "PKCS12 data is empty",
		}
	}
	
	privateKey, cert, err := pkcs12.Decode(data, password)
	if err != nil {
		return nil, &v2.CertificateError{
			Code:    "PKCS12_DECODE_ERROR",
			Message: fmt.Sprintf("failed to decode PKCS12 data: %v", err),
			Err:     err,
		}
	}
	
	var allCerts []*x509.Certificate
	if cert != nil {
		allCerts = append(allCerts, cert)
	}
	
	if len(allCerts) == 0 {
		return nil, &v2.CertificateError{
			Code:    "NO_CERTIFICATES_IN_PKCS12",
			Message: "no certificates found in PKCS12 data",
		}
	}
	
	_ = privateKey // Not stored for security
	return parser.ConvertToBundle(allCerts, v2.CertFormatPKCS12, ""), nil
}

func ValidatePKCS12(data []byte) error {
	if len(data) == 0 {
		return errors.New("PKCS12 data is empty")
	}
	
	var pfx struct {
		Version    int
		AuthSafe   asn1.RawValue
		MacData    asn1.RawValue `asn1:"optional"`
	}
	
	_, err := asn1.Unmarshal(data, &pfx)
	if err != nil {
		return fmt.Errorf("invalid PKCS12 structure: %v", err)
	}
	
	return nil
}

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

func isPKCS12Magic(data []byte) bool {
	if len(data) < 4 {
		return false
	}
	return data[0] == 0x30 && data[2] == 0x02 && data[3] == 0x01
}

func isPKCS7Structure(data []byte) bool {
	pkcs7Patterns := [][]byte{
		{0x06, 0x09, 0x2A, 0x86, 0x48, 0x86, 0xF7, 0x0D, 0x01},
	}
	
	for _, pattern := range pkcs7Patterns {
		if containsBytes(data[:min(len(data), 50)], pattern) {
			return true
		}
	}
	return false
}

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

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}