package certificates

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"time"

	v2 "github.com/cnoe-io/idpbuilder/pkg/oci/api/v2"
)

// BundleLoader integrates Split-001 types with Split-002 format detection
type BundleLoader struct {
	detector FormatDetector
}

// FormatDetector interface (from Split-002)
type FormatDetector interface {
	DetectFormat(data []byte) (v2.CertFormat, error)
}

// NewBundleLoader creates a new certificate bundle loader
func NewBundleLoader() *BundleLoader {
	return &BundleLoader{
		detector: &MagicBytesDetector{},
	}
}

// LoadBundle loads certificate bundle from file path with auto-detection
func (l *BundleLoader) LoadBundle(ctx context.Context, path string) (*v2.CertBundle, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	if path == "" {
		return nil, &v2.CertificateError{
			Code:    "EMPTY_PATH",
			Message: "certificate file path is empty",
		}
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, &v2.CertificateError{
			Code:    "FILE_READ_ERROR",
			Message: fmt.Sprintf("failed to read certificate file %s: %v", path, err),
			Err:     err,
		}
	}

	bundle, err := l.LoadFromData(ctx, data)
	if err != nil {
		return nil, err
	}

	bundle.Source = path
	return bundle, nil
}

// LoadFromPath alias for consistency
func (l *BundleLoader) LoadFromPath(ctx context.Context, path string) (*v2.CertBundle, error) {
	return l.LoadBundle(ctx, path)
}

// LoadFromReader loads certificate bundle from io.Reader
func (l *BundleLoader) LoadFromReader(ctx context.Context, reader io.Reader) (*v2.CertBundle, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	if reader == nil {
		return nil, &v2.CertificateError{
			Code:    "NULL_READER",
			Message: "certificate reader is nil",
		}
	}

	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, &v2.CertificateError{
			Code:    "READER_ERROR",
			Message: fmt.Sprintf("failed to read certificate data: %v", err),
			Err:     err,
		}
	}

	return l.LoadFromData(ctx, data)
}

// LoadFromData core integration point - uses format detection and parsing
func (l *BundleLoader) LoadFromData(ctx context.Context, data []byte) (*v2.CertBundle, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	if len(data) == 0 {
		return nil, &v2.CertificateError{
			Code:    "EMPTY_DATA",
			Message: "certificate data is empty",
		}
	}

	// Auto-detect format using Split-002 functionality
	format, err := l.detector.DetectFormat(data)
	if err != nil {
		return nil, &v2.CertificateError{
			Code:    "FORMAT_DETECTION_ERROR",
			Message: fmt.Sprintf("failed to detect certificate format: %v", err),
			Err:     err,
		}
	}

	// Create parser for detected format
	var parser FormatParser
	switch format {
	case v2.CertFormatPEM:
		parser = &PEMParser{}
	case v2.CertFormatDER:
		parser = &DERParser{}
	case v2.CertFormatPKCS7:
		parser = &PKCS7Parser{}
	case v2.CertFormatPKCS12:
		parser = &PKCS12Parser{}
	default:
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
			Message: fmt.Sprintf("failed to parse certificate data as %s: %v", format, err),
			Err:     err,
		}
	}

	// Set metadata using Split-001 types
	bundle.Format = format
	bundle.LoadedAt = time.Now()

	return bundle, nil
}

// DetectFormat exposes format detection functionality
func (l *BundleLoader) DetectFormat(data []byte) (v2.CertFormat, error) {
	return l.detector.DetectFormat(data)
}

// FormatParser interface for format-specific parsing
type FormatParser interface {
	Parse(data []byte) (*v2.CertBundle, error)
}