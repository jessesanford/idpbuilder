package certificates

import (
	"crypto/x509"
	"fmt"
)

// CertificateError represents a certificate-related error with context
type CertificateError struct {
	// Code is a machine-readable error code
	Code string `json:"code"`
	// Message is a human-readable error message
	Message string `json:"message"`
	// Cert is the certificate that caused the error (if applicable)
	Cert *x509.Certificate `json:"-"`
	// Err is the underlying error (if any)
	Err error `json:"-"`
}

// Error implements the error interface
func (e *CertificateError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s (underlying: %v)", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Unwrap returns the underlying error
func (e *CertificateError) Unwrap() error {
	return e.Err
}

// Common error codes
const (
	// File and I/O errors
	ErrorCodeFileReadError     = "FILE_READ_ERROR"
	ErrorCodeEmptyFile         = "EMPTY_FILE"
	ErrorCodeDirectoryReadError = "DIRECTORY_READ_ERROR"
	
	// Format detection errors
	ErrorCodeFormatDetectionError = "FORMAT_DETECTION_ERROR"
	ErrorCodeUnsupportedFormat   = "UNSUPPORTED_FORMAT"
	ErrorCodeUnknownFormat       = "UNKNOWN_FORMAT"
	ErrorCodeEmptyData           = "EMPTY_DATA"
	
	// Parse errors
	ErrorCodeParseError          = "PARSE_ERROR"
	ErrorCodePEMParseError       = "PEM_PARSE_ERROR"
	ErrorCodeDERParseError       = "DER_PARSE_ERROR"
	ErrorCodePKCS7ParseError     = "PKCS7_PARSE_ERROR"
	ErrorCodePKCS12ParseError    = "PKCS12_PARSE_ERROR"
	
	// Validation errors
	ErrorCodeCertValidationError = "CERT_VALIDATION_ERROR"
	ErrorCodeNullCertificate     = "NULL_CERTIFICATE"
	ErrorCodeNullBundle          = "NULL_BUNDLE"
	ErrorCodeEmptyBundle         = "EMPTY_BUNDLE"
	
	// Certificate status errors
	ErrorCodeCertExpired         = "CERT_EXPIRED"
	ErrorCodeCertNotYetValid     = "CERT_NOT_YET_VALID"
)

// NewCertificateError creates a new certificate error
func NewCertificateError(code, message string) *CertificateError {
	return &CertificateError{
		Code:    code,
		Message: message,
	}
}

// NewCertificateErrorWithCause creates a new certificate error wrapping another error
func NewCertificateErrorWithCause(code, message string, err error) *CertificateError {
	return &CertificateError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}