package certs

import (
	"fmt"
)

// CertErrorCode represents the type of certificate error that occurred.
type CertErrorCode int

const (
	// ErrInvalidCert indicates the certificate is malformed or invalid.
	ErrInvalidCert CertErrorCode = iota + 1

	// ErrExpired indicates the certificate has expired.
	ErrExpired

	// ErrNotYetValid indicates the certificate is not yet valid.
	ErrNotYetValid

	// ErrUntrusted indicates the certificate is not trusted by the system.
	ErrUntrusted

	// ErrSystemPool indicates a failure to load the system certificate pool.
	ErrSystemPool

	// ErrTLSConfig indicates a failure to create TLS configuration.
	ErrTLSConfig
)

// CertError represents a certificate-related error with additional context.
// It preserves the original error while providing structured error codes
// for programmatic error handling.
type CertError struct {
	Code    CertErrorCode
	Message string
	Err     error
}

// Error implements the error interface for CertError.
func (e *CertError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("certificate error (%d): %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("certificate error (%d): %s", e.Code, e.Message)
}

// Unwrap returns the underlying error for error wrapping support.
func (e *CertError) Unwrap() error {
	return e.Err
}

// NewCertError creates a new CertError with the specified code, message and underlying error.
func NewCertError(code CertErrorCode, message string, err error) *CertError {
	return &CertError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}