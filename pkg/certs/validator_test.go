package certs

import (
	"crypto/x509"
	"testing"
	"time"
)

func TestNewDefaultValidator(t *testing.T) {
	validator := NewDefaultValidator()
	if validator == nil {
		t.Error("NewDefaultValidator() should not return nil")
	}
}

func TestDefaultValidator_Validate_NilCertificate(t *testing.T) {
	validator := NewDefaultValidator()
	opts := x509.VerifyOptions{}

	err := validator.Validate(nil, opts)
	if err == nil {
		t.Error("Validate() should fail with nil certificate")
	}

	certErr, ok := err.(*CertError)
	if !ok {
		t.Errorf("Validate() error type = %T, want *CertError", err)
	} else if certErr.Code != ErrInvalidCert {
		t.Errorf("Validate() error code = %v, want %v", certErr.Code, ErrInvalidCert)
	}
}

func TestDefaultValidator_Validate_ExpiredCertificate(t *testing.T) {
	validator := NewDefaultValidator()

	// Create an expired certificate
	cert := &x509.Certificate{
		NotBefore: time.Now().Add(-2 * time.Hour),
		NotAfter:  time.Now().Add(-time.Hour), // Expired 1 hour ago
	}

	opts := x509.VerifyOptions{}
	err := validator.Validate(cert, opts)
	if err == nil {
		t.Error("Validate() should fail with expired certificate")
	}

	certErr, ok := err.(*CertError)
	if !ok {
		t.Errorf("Validate() error type = %T, want *CertError", err)
	} else if certErr.Code != ErrExpired {
		t.Errorf("Validate() error code = %v, want %v", certErr.Code, ErrExpired)
	}
}

func TestDefaultValidator_Validate_NotYetValidCertificate(t *testing.T) {
	validator := NewDefaultValidator()

	// Create a certificate that's not yet valid
	cert := &x509.Certificate{
		NotBefore: time.Now().Add(time.Hour),     // Valid starting 1 hour from now
		NotAfter:  time.Now().Add(2 * time.Hour), // Valid until 2 hours from now
	}

	opts := x509.VerifyOptions{}
	err := validator.Validate(cert, opts)
	if err == nil {
		t.Error("Validate() should fail with not-yet-valid certificate")
	}

	certErr, ok := err.(*CertError)
	if !ok {
		t.Errorf("Validate() error type = %T, want *CertError", err)
	} else if certErr.Code != ErrNotYetValid {
		t.Errorf("Validate() error code = %v, want %v", certErr.Code, ErrNotYetValid)
	}
}

func TestDefaultValidator_Validate_ValidTimePeriod(t *testing.T) {
	validator := NewDefaultValidator()

	// Create a certificate with valid time period
	cert := &x509.Certificate{
		NotBefore: time.Now().Add(-time.Hour),
		NotAfter:  time.Now().Add(time.Hour),
	}

	// No roots provided, so only time validation should occur
	opts := x509.VerifyOptions{}
	err := validator.Validate(cert, opts)
	if err != nil {
		t.Errorf("Validate() should succeed with valid time period, got error: %v", err)
	}
}

func TestDefaultValidator_ValidateChain_EmptyChain(t *testing.T) {
	validator := NewDefaultValidator()
	roots := x509.NewCertPool()

	err := validator.ValidateChain([]*x509.Certificate{}, roots)
	if err == nil {
		t.Error("ValidateChain() should fail with empty chain")
	}

	certErr, ok := err.(*CertError)
	if !ok {
		t.Errorf("ValidateChain() error type = %T, want *CertError", err)
	} else if certErr.Code != ErrInvalidCert {
		t.Errorf("ValidateChain() error code = %v, want %v", certErr.Code, ErrInvalidCert)
	}
}

func TestDefaultValidator_ValidateChain_SingleCert(t *testing.T) {
	validator := NewDefaultValidator()
	roots := x509.NewCertPool()

	// Create a certificate with valid time period
	cert := &x509.Certificate{
		NotBefore: time.Now().Add(-time.Hour),
		NotAfter:  time.Now().Add(time.Hour),
	}

	// Since we're using an empty root pool and a test cert that won't verify,
	// this should pass time validation but may fail chain validation
	err := validator.ValidateChain([]*x509.Certificate{cert}, roots)
	// We expect this to potentially fail due to chain validation, but not due to empty chain
	if err != nil {
		// This is expected - the cert won't validate against an empty root pool
		// The important thing is we don't get an "empty chain" error
		certErr, ok := err.(*CertError)
		if ok && certErr.Code == ErrInvalidCert && certErr.Message == "empty certificate chain" {
			t.Error("ValidateChain() should not report empty chain for single certificate")
		}
	}
}