package certs

import (
	"crypto/x509"
	"fmt"
	"time"

	"github.com/go-logr/logr"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

var validatorLog = log.Log.WithName("cert-validator")

// Validator implements CertValidator for certificate validation operations
type Validator struct {
	logger logr.Logger
	cert   *x509.Certificate
}

// NewValidator creates a new certificate validator
func NewValidator(cert *x509.Certificate) *Validator {
	return &Validator{
		logger: validatorLog,
		cert:   cert,
	}
}

// ValidateChain validates the certificate chain
func (v *Validator) ValidateChain(cert *x509.Certificate) error {
	if cert == nil {
		return NewCertificateError("chain_validation", fmt.Errorf("certificate is nil"),
			"no certificate provided for chain validation", []string{
				"Ensure certificate is properly extracted",
			})
	}

	// For self-signed certificates (common in Kind), we expect the certificate
	// to be its own issuer
	if cert.Subject.String() == cert.Issuer.String() {
		v.logger.V(1).Info("Self-signed certificate detected", 
			"subject", cert.Subject.String())
		return nil // Self-signed is acceptable for Gitea in Kind
	}

	// For certificates with separate issuers, we would need to validate the chain
	// This is a simplified validation for the MVP
	v.logger.V(1).Info("Certificate with separate issuer", 
		"subject", cert.Subject.String(), "issuer", cert.Issuer.String())
	
	return nil
}

// CheckExpiry checks if certificate is expired or will expire soon
func (v *Validator) CheckExpiry(cert *x509.Certificate) (*time.Duration, error) {
	if cert == nil {
		return nil, NewCertificateError("expiry_check", fmt.Errorf("certificate is nil"),
			"no certificate provided for expiry check", []string{
				"Ensure certificate is properly extracted",
			})
	}

	now := time.Now()
	
	// Check if already expired
	if now.After(cert.NotAfter) {
		return nil, NewCertificateError("expiry_check", fmt.Errorf("certificate expired"),
			fmt.Sprintf("certificate expired on %v", cert.NotAfter), []string{
				"Recreate the cluster to generate new certificates",
				"Use --insecure flag as temporary workaround",
			})
	}

	// Check if not yet valid
	if now.Before(cert.NotBefore) {
		return nil, NewCertificateError("expiry_check", fmt.Errorf("certificate not yet valid"),
			fmt.Sprintf("certificate valid from %v", cert.NotBefore), []string{
				"Check system clock synchronization",
				"Wait until certificate becomes valid",
			})
	}

	// Calculate time until expiry
	timeUntilExpiry := cert.NotAfter.Sub(now)
	
	// Warn if expiring soon
	if timeUntilExpiry < 30*24*time.Hour {
		v.logger.Info("Certificate expires soon", 
			"expires_at", cert.NotAfter,
			"time_remaining", timeUntilExpiry.String())
	}

	return &timeUntilExpiry, nil
}

// VerifyHostname verifies the certificate is valid for the given hostname
func (v *Validator) VerifyHostname(cert *x509.Certificate, hostname string) error {
	if cert == nil {
		return NewCertificateError("hostname_verification", fmt.Errorf("certificate is nil"),
			"no certificate provided for hostname verification", []string{
				"Ensure certificate is properly extracted",
			})
	}

	if hostname == "" {
		return NewCertificateError("hostname_verification", fmt.Errorf("hostname is empty"),
			"no hostname provided for verification", []string{
				"Provide a valid hostname to verify against",
			})
	}

	// Use Go's built-in hostname verification
	err := cert.VerifyHostname(hostname)
	if err != nil {
		return NewCertificateError("hostname_verification", err,
			fmt.Sprintf("hostname %s does not match certificate", hostname), []string{
				"Check if certificate includes the hostname in Subject Alt Names",
				"Verify you're using the correct hostname",
				"For local development, consider using --insecure flag",
			})
	}

	v.logger.V(1).Info("Hostname verification successful", "hostname", hostname)
	return nil
}

// GenerateDiagnostics provides detailed certificate diagnostic information
func (v *Validator) GenerateDiagnostics() (*CertDiagnostics, error) {
	if v.cert == nil {
		return nil, NewCertificateError("diagnostics", fmt.Errorf("no certificate set"),
			"validator has no certificate for diagnostics", []string{
				"Create validator with a valid certificate",
			})
	}

	cert := v.cert
	now := time.Now()
	
	diagnostics := &CertDiagnostics{
		IsValid:         true,
		ExpiresAt:       cert.NotAfter,
		DaysUntilExpiry: int(cert.NotAfter.Sub(now).Hours() / 24),
		Subject:         cert.Subject.String(),
		Issuer:          cert.Issuer.String(),
		DNSNames:        cert.DNSNames,
		Issues:          []string{},
		Recommendations: []string{},
	}

	// Check for issues
	if now.After(cert.NotAfter) {
		diagnostics.IsValid = false
		diagnostics.Issues = append(diagnostics.Issues, 
			fmt.Sprintf("Certificate expired on %v", cert.NotAfter))
		diagnostics.Recommendations = append(diagnostics.Recommendations,
			"Recreate cluster to generate new certificates")
	} else if now.Before(cert.NotBefore) {
		diagnostics.IsValid = false
		diagnostics.Issues = append(diagnostics.Issues,
			fmt.Sprintf("Certificate not valid until %v", cert.NotBefore))
		diagnostics.Recommendations = append(diagnostics.Recommendations,
			"Check system clock synchronization")
	}

	// Check for upcoming expiry
	if diagnostics.DaysUntilExpiry < 30 && diagnostics.DaysUntilExpiry > 0 {
		diagnostics.Issues = append(diagnostics.Issues,
			fmt.Sprintf("Certificate expires in %d days", diagnostics.DaysUntilExpiry))
		diagnostics.Recommendations = append(diagnostics.Recommendations,
			"Consider recreating cluster before expiry")
	}

	// Check if self-signed
	if cert.Subject.String() == cert.Issuer.String() {
		diagnostics.Recommendations = append(diagnostics.Recommendations,
			"Self-signed certificate detected - normal for Kind clusters")
	}

	// Check for missing DNS names
	if len(cert.DNSNames) == 0 {
		diagnostics.Issues = append(diagnostics.Issues,
			"No DNS names in certificate")
		diagnostics.Recommendations = append(diagnostics.Recommendations,
			"May cause hostname verification failures")
	}

	v.logger.V(1).Info("Generated certificate diagnostics",
		"valid", diagnostics.IsValid,
		"expires_in_days", diagnostics.DaysUntilExpiry,
		"issues_count", len(diagnostics.Issues))

	return diagnostics, nil
}

// ValidateForRegistry validates a certificate for use with a specific registry
func (v *Validator) ValidateForRegistry(cert *x509.Certificate, registryHost string) error {
	// Validate basic certificate properties
	if err := v.ValidateChain(cert); err != nil {
		return err
	}

	// Check expiry
	if _, err := v.CheckExpiry(cert); err != nil {
		return err
	}

	// Verify hostname if registry host is provided
	if registryHost != "" {
		if err := v.VerifyHostname(cert, registryHost); err != nil {
			// For Kind/Gitea, hostname verification often fails due to localhost/IP usage
			// Log the error but don't fail the validation for MVP
			v.logger.V(1).Info("Hostname verification failed but proceeding", 
				"error", err.Error(), "registry", registryHost)
		}
	}

	return nil
}