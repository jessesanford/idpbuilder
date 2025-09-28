package certs

import (
	"crypto/x509"
	"time"
)

// Validator provides an interface for certificate validation operations.
// This abstraction allows for different validation strategies and enables
// easy testing through mock implementations.
type Validator interface {
	// Validate performs certificate validation using the provided options.
	// It returns an error if the certificate is invalid, expired, or cannot
	// be verified against the provided verification options.
	Validate(cert *x509.Certificate, opts x509.VerifyOptions) error
}

// DefaultValidator implements the Validator interface using standard x509 validation.
// It provides production-ready certificate validation using the crypto/x509 package.
type DefaultValidator struct{}

// NewDefaultValidator creates a new DefaultValidator instance.
func NewDefaultValidator() *DefaultValidator {
	return &DefaultValidator{}
}

// Validate implements the Validator interface using standard x509 verification.
// It performs comprehensive certificate validation including:
// - Certificate expiry checks
// - Signature verification
// - Chain validation against provided roots
// - Intermediate certificate handling
func (v *DefaultValidator) Validate(cert *x509.Certificate, opts x509.VerifyOptions) error {
	if cert == nil {
		return NewCertError(ErrInvalidCert, "certificate is nil", nil)
	}

	// Check certificate validity period
	now := time.Now()
	if now.Before(cert.NotBefore) {
		return NewCertError(ErrNotYetValid, "certificate is not yet valid", nil)
	}
	if now.After(cert.NotAfter) {
		return NewCertError(ErrExpired, "certificate has expired", nil)
	}

	// Perform x509 verification if roots are provided
	if opts.Roots != nil {
		_, err := cert.Verify(opts)
		if err != nil {
			// Determine the specific error type based on the x509 error
			if certErr, ok := err.(x509.CertificateInvalidError); ok {
				switch certErr.Reason {
				case x509.Expired:
					return NewCertError(ErrExpired, "certificate verification failed: expired", err)
				case x509.NotAuthorizedToSign:
					return NewCertError(ErrUntrusted, "certificate verification failed: not authorized", err)
				default:
					return NewCertError(ErrInvalidCert, "certificate verification failed", err)
				}
			}

			if _, ok := err.(x509.UnknownAuthorityError); ok {
				return NewCertError(ErrUntrusted, "certificate verification failed: unknown authority", err)
			}

			return NewCertError(ErrInvalidCert, "certificate verification failed", err)
		}
	}

	return nil
}

// ValidateChain validates a certificate chain where the first certificate
// is the leaf certificate and subsequent certificates are intermediates.
// The validation is performed against the provided root certificate pool.
func (v *DefaultValidator) ValidateChain(certs []*x509.Certificate, roots *x509.CertPool) error {
	if len(certs) == 0 {
		return NewCertError(ErrInvalidCert, "empty certificate chain", nil)
	}

	// Build intermediate pool from all certificates except the first (leaf)
	intermediates := x509.NewCertPool()
	for i := 1; i < len(certs); i++ {
		intermediates.AddCert(certs[i])
	}

	// Validate the leaf certificate
	opts := x509.VerifyOptions{
		Roots:         roots,
		Intermediates: intermediates,
	}

	return v.Validate(certs[0], opts)
}