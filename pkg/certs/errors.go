package certs

import "fmt"

// CertificateError represents a certificate-related error
type CertificateError struct {
	Code    string
	Message string
}

// Error implements the error interface
func (e *CertificateError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// NewCertificateError creates a new certificate error
func NewCertificateError(code, message string) *CertificateError {
	return &CertificateError{
		Code:    code,
		Message: message,
	}
}