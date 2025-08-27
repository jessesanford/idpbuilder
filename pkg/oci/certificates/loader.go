package certificates

import (
	"context"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"time"

	v2 "github.com/cnoe-io/idpbuilder/pkg/oci/api/v2"
)

// MultiFormatLoader handles loading certificates from multiple formats
type MultiFormatLoader struct {
	parsers map[v2.CertFormat]FormatParser
	detector FormatDetector
	parser   *CertificateParser
}

// FormatParser interface for format-specific parsers
type FormatParser interface {
	Parse(data []byte) (*v2.CertBundle, error)
	Validate(data []byte) error
}

// FormatDetector auto-detects certificate format from file content
type FormatDetector interface {
	DetectFormat(data []byte) (v2.CertFormat, error)
}

// NewMultiFormatLoader creates a new multi-format certificate loader
func NewMultiFormatLoader() *MultiFormatLoader {
	parser := &CertificateParser{
		strictMode:    false,
		maxChainDepth: 10,
	}
	
	return &MultiFormatLoader{
		parsers: map[v2.CertFormat]FormatParser{
			v2.CertFormatPEM:    &PEMParser{parser: parser},
			v2.CertFormatDER:    &DERParser{parser: parser},
			v2.CertFormatPKCS7:  &PKCS7Parser{parser: parser},
			v2.CertFormatPKCS12: &PKCS12Parser{parser: parser},
		},
		detector: &MagicBytesDetector{},
		parser:   parser,
	}
}

// LoadBundle loads a certificate bundle with auto-detection
func (l *MultiFormatLoader) LoadBundle(ctx context.Context, path string) (*v2.CertBundle, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	// Read the file
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, &v2.CertificateError{
			Code:    "FILE_READ_ERROR",
			Message: fmt.Sprintf("failed to read certificate file %s: %v", path, err),
			Err:     err,
		}
	}

	if len(data) == 0 {
		return nil, &v2.CertificateError{
			Code:    "EMPTY_FILE",
			Message: fmt.Sprintf("certificate file %s is empty", path),
		}
	}

	// Auto-detect format
	format, err := l.detector.DetectFormat(data)
	if err != nil {
		return nil, &v2.CertificateError{
			Code:    "FORMAT_DETECTION_ERROR",
			Message: fmt.Sprintf("failed to detect certificate format for %s: %v", path, err),
			Err:     err,
		}
	}

	// Get appropriate parser
	parser, exists := l.parsers[format]
	if !exists {
		return nil, &v2.CertificateError{
			Code:    "UNSUPPORTED_FORMAT",
			Message: fmt.Sprintf("unsupported certificate format: %s", format),
		}
	}

	// Parse the certificate bundle
	bundle, err := parser.Parse(data)
	if err != nil {
		return nil, &v2.CertificateError{
			Code:    "PARSE_ERROR",
			Message: fmt.Sprintf("failed to parse certificate %s as %s: %v", path, format, err),
			Err:     err,
		}
	}

	// Set additional metadata
	bundle.Source = path
	bundle.LoadedAt = time.Now()
	bundle.Format = format

	return bundle, nil
}

// LoadPEM loads certificates from PEM format
func (l *MultiFormatLoader) LoadPEM(ctx context.Context, data []byte) (*v2.CertBundle, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	parser := l.parsers[v2.CertFormatPEM]
	bundle, err := parser.Parse(data)
	if err != nil {
		return nil, err
	}

	bundle.LoadedAt = time.Now()
	bundle.Format = v2.CertFormatPEM
	return bundle, nil
}

// LoadDER loads certificates from DER format  
func (l *MultiFormatLoader) LoadDER(ctx context.Context, data []byte) (*v2.CertBundle, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	parser := l.parsers[v2.CertFormatDER]
	bundle, err := parser.Parse(data)
	if err != nil {
		return nil, err
	}

	bundle.LoadedAt = time.Now()
	bundle.Format = v2.CertFormatDER
	return bundle, nil
}

// LoadPKCS7 loads certificates from PKCS7 format
func (l *MultiFormatLoader) LoadPKCS7(ctx context.Context, data []byte) (*v2.CertBundle, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	parser := l.parsers[v2.CertFormatPKCS7]
	bundle, err := parser.Parse(data)
	if err != nil {
		return nil, err
	}

	bundle.LoadedAt = time.Now()
	bundle.Format = v2.CertFormatPKCS7
	return bundle, nil
}

// LoadPKCS12 loads certificates from PKCS12 format with password
func (l *MultiFormatLoader) LoadPKCS12(ctx context.Context, data []byte, password string) (*v2.CertBundle, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	parser := l.parsers[v2.CertFormatPKCS12]
	
	// Check if parser supports password-based loading
	if p12Parser, ok := parser.(*PKCS12Parser); ok {
		bundle, err := p12Parser.ParseWithPassword(data, password)
		if err != nil {
			return nil, err
		}
		
		bundle.LoadedAt = time.Now()
		bundle.Format = v2.CertFormatPKCS12
		return bundle, nil
	}

	return nil, &v2.CertificateError{
		Code:    "PARSER_ERROR",
		Message: "PKCS12 parser does not support password-based loading",
	}
}

// DetectFormat auto-detects certificate format from file content
func (l *MultiFormatLoader) DetectFormat(data []byte) (v2.CertFormat, error) {
	return l.detector.DetectFormat(data)
}

// SetStrictMode enables or disables strict certificate validation
func (l *MultiFormatLoader) SetStrictMode(strict bool) {
	l.parser.strictMode = strict
}

// SetMaxChainDepth sets the maximum allowed certificate chain depth
func (l *MultiFormatLoader) SetMaxChainDepth(depth int) {
	if depth > 0 {
		l.parser.maxChainDepth = depth
	}
}

// ValidateBundle performs comprehensive validation of a certificate bundle
func (l *MultiFormatLoader) ValidateBundle(bundle *v2.CertBundle) error {
	if bundle == nil {
		return &v2.CertificateError{
			Code:    "NULL_BUNDLE",
			Message: "certificate bundle is nil",
		}
	}

	if len(bundle.Certificates) == 0 && len(bundle.CAs) == 0 {
		return &v2.CertificateError{
			Code:    "EMPTY_BUNDLE",
			Message: "certificate bundle contains no certificates",
		}
	}

	// Validate each certificate in the bundle
	for i, cert := range bundle.Certificates {
		if err := l.parser.ValidateCertificate(cert); err != nil {
			return &v2.CertificateError{
				Code:    "CERT_VALIDATION_ERROR",
				Message: fmt.Sprintf("certificate %d validation failed: %v", i, err),
				Cert:    cert,
				Err:     err,
			}
		}
	}

	// Validate CA certificates
	for i, cert := range bundle.CAs {
		if err := l.parser.ValidateCertificate(cert); err != nil {
			return &v2.CertificateError{
				Code:    "CA_VALIDATION_ERROR",
				Message: fmt.Sprintf("CA certificate %d validation failed: %v", i, err),
				Cert:    cert,
				Err:     err,
			}
		}
	}

	// If strict mode, validate certificate chains
	if l.parser.strictMode {
		for _, cert := range bundle.Certificates {
			if err := l.parser.ValidateChain([]*x509.Certificate{cert}); err != nil {
				return &v2.CertificateError{
					Code:    "CHAIN_VALIDATION_ERROR",
					Message: fmt.Sprintf("certificate chain validation failed: %v", err),
					Cert:    cert,
					Err:     err,
				}
			}
		}
	}

	return nil
}

// LoadFromDirectory loads all certificate files from a directory
func (l *MultiFormatLoader) LoadFromDirectory(ctx context.Context, dirPath string) ([]*v2.CertBundle, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return nil, &v2.CertificateError{
			Code:    "DIRECTORY_READ_ERROR",
			Message: fmt.Sprintf("failed to read directory %s: %v", dirPath, err),
			Err:     err,
		}
	}

	var bundles []*v2.CertBundle
	supportedExts := map[string]bool{
		".pem": true, ".crt": true, ".cert": true,
		".der": true, ".cer": true,
		".p7b": true, ".p7c": true,
		".p12": true, ".pfx": true,
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		ext := filepath.Ext(file.Name())
		if !supportedExts[ext] {
			continue
		}

		filePath := filepath.Join(dirPath, file.Name())
		bundle, err := l.LoadBundle(ctx, filePath)
		if err != nil {
			// Log error but continue with other files
			continue
		}

		bundles = append(bundles, bundle)
	}

	if len(bundles) == 0 {
		return nil, &v2.CertificateError{
			Code:    "NO_CERTIFICATES_FOUND",
			Message: fmt.Sprintf("no valid certificate files found in directory %s", dirPath),
		}
	}

	return bundles, nil
}