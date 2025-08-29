package certs

import (
	"crypto/x509"
	"fmt"
	"strings"
	"time"
)

// Validator implements the CertValidator interface for certificate validation
type Validator struct {
	// allowSelfSigned determines if self-signed certificates are acceptable (for Kind clusters)
	allowSelfSigned bool
	
	// requiredKeyUsages specifies the key usages that must be present
	requiredKeyUsages []x509.KeyUsage
	
	// requiredExtKeyUsages specifies the extended key usages that must be present
	requiredExtKeyUsages []x509.ExtKeyUsage
}

// ValidatorConfig holds configuration options for the validator
type ValidatorConfig struct {
	// AllowSelfSigned allows self-signed certificates (useful for Kind clusters)
	AllowSelfSigned bool
	
	// RequiredKeyUsages specifies key usages that must be present
	RequiredKeyUsages []x509.KeyUsage
	
	// RequiredExtKeyUsages specifies extended key usages that must be present
	RequiredExtKeyUsages []x509.ExtKeyUsage
}

// DefaultValidatorConfig returns a validator configuration with sensible defaults for Kind clusters
func DefaultValidatorConfig() *ValidatorConfig {
	return &ValidatorConfig{
		AllowSelfSigned: true, // Kind clusters typically use self-signed certs
		RequiredKeyUsages: []x509.KeyUsage{
			x509.KeyUsageDigitalSignature,
			x509.KeyUsageKeyEncipherment,
		},
		RequiredExtKeyUsages: []x509.ExtKeyUsage{
			x509.ExtKeyUsageServerAuth,
		},
	}
}

// NewValidator creates a new certificate validator with the specified configuration
func NewValidator(config *ValidatorConfig) *Validator {
	if config == nil {
		config = DefaultValidatorConfig()
	}
	
	return &Validator{
		allowSelfSigned:      config.AllowSelfSigned,
		requiredKeyUsages:    config.RequiredKeyUsages,
		requiredExtKeyUsages: config.RequiredExtKeyUsages,
	}
}

// ValidateCertificate validates a certificate against the configured requirements
func (v *Validator) ValidateCertificate(cert *x509.Certificate) (*ValidationResult, error) {
	if cert == nil {
		return nil, NewCertificateError("INVALID_CERTIFICATE", "Certificate is nil")
	}
	
	result := &ValidationResult{
		Valid:    true,
		Issues:   []string{},
		Warnings: []string{},
	}
	
	// Check if certificate is expired
	now := time.Now()
	if cert.NotAfter.Before(now) {
		result.Valid = false
		result.Issues = append(result.Issues, fmt.Sprintf("Certificate expired on %s", cert.NotAfter.Format("2006-01-02 15:04:05")))
	}
	
	// Check if certificate is not yet valid
	if cert.NotBefore.After(now) {
		result.Valid = false
		result.Issues = append(result.Issues, fmt.Sprintf("Certificate not valid until %s", cert.NotBefore.Format("2006-01-02 15:04:05")))
	}
	
	// Check self-signed certificates
	if cert.Issuer.String() == cert.Subject.String() {
		if !v.allowSelfSigned {
			result.Valid = false
			result.Issues = append(result.Issues, "Self-signed certificates are not allowed")
		} else {
			result.Warnings = append(result.Warnings, "Certificate is self-signed (acceptable for Kind clusters)")
		}
	}
	
	// Validate key usages
	for _, requiredUsage := range v.requiredKeyUsages {
		if cert.KeyUsage&requiredUsage == 0 {
			result.Valid = false
			result.Issues = append(result.Issues, fmt.Sprintf("Required key usage missing: %v", requiredUsage))
		}
	}
	
	// Validate extended key usages
	for _, requiredExtUsage := range v.requiredExtKeyUsages {
		found := false
		for _, extUsage := range cert.ExtKeyUsage {
			if extUsage == requiredExtUsage {
				found = true
				break
			}
		}
		if !found {
			result.Valid = false
			result.Issues = append(result.Issues, fmt.Sprintf("Required extended key usage missing: %v", requiredExtUsage))
		}
	}
	
	// Validate subject information
	if cert.Subject.CommonName == "" {
		result.Warnings = append(result.Warnings, "Certificate has empty Common Name")
	}
	
	// Check for DNS SANs for web certificates
	if len(cert.DNSNames) == 0 && len(cert.IPAddresses) == 0 {
		result.Warnings = append(result.Warnings, "Certificate has no Subject Alternative Names (DNS or IP)")
	}
	
	return result, nil
}

// CheckExpiry checks if a certificate is expired or expiring soon
func (v *Validator) CheckExpiry(cert *x509.Certificate, warnDays int) (*ExpiryResult, error) {
	if cert == nil {
		return nil, NewCertificateError("INVALID_CERTIFICATE", "Certificate is nil")
	}
	
	if warnDays < 0 {
		warnDays = 30 // Default warning period of 30 days
	}
	
	now := time.Now()
	
	result := &ExpiryResult{
		ExpiryDate: cert.NotAfter,
	}
	
	// Check if certificate is already expired
	if cert.NotAfter.Before(now) {
		result.Expired = true
		result.ExpiringSoon = false
		result.DaysUntilExpiry = int(now.Sub(cert.NotAfter).Hours() / 24) * -1 // Negative for expired
		return result, nil
	}
	
	// Calculate days until expiry
	daysUntilExpiry := int(cert.NotAfter.Sub(now).Hours() / 24)
	result.DaysUntilExpiry = daysUntilExpiry
	
	// Check if certificate is expiring soon
	if daysUntilExpiry <= warnDays {
		result.ExpiringSoon = true
	}
	
	return result, nil
}

// ValidateChain validates a certificate chain for proper ordering and trust
func (v *Validator) ValidateChain(certs []*x509.Certificate) (*ValidationResult, error) {
	if len(certs) == 0 {
		return nil, NewCertificateError("EMPTY_CHAIN", "Certificate chain is empty")
	}
	
	result := &ValidationResult{
		Valid:    true,
		Issues:   []string{},
		Warnings: []string{},
	}
	
	// Validate each certificate in the chain
	for i, cert := range certs {
		certResult, err := v.ValidateCertificate(cert)
		if err != nil {
			return nil, fmt.Errorf("failed to validate certificate %d in chain: %w", i, err)
		}
		
		// Merge results
		if !certResult.Valid {
			result.Valid = false
		}
		result.Issues = append(result.Issues, certResult.Issues...)
		result.Warnings = append(result.Warnings, certResult.Warnings...)
	}
	
	// Validate chain ordering (leaf certificate should be first)
	if len(certs) > 1 {
		leaf := certs[0]
		parent := certs[1]
		
		// Check if parent issued the leaf certificate
		if leaf.Issuer.String() != parent.Subject.String() {
			result.Valid = false
			result.Issues = append(result.Issues, "Certificate chain is not properly ordered or linked")
		}
	}
	
	return result, nil
}

// VerifyHostname verifies if a certificate is valid for a given hostname
func (v *Validator) VerifyHostname(cert *x509.Certificate, hostname string) error {
	if cert == nil {
		return NewCertificateError("INVALID_CERTIFICATE", "Certificate is nil")
	}
	
	if hostname == "" {
		return NewCertificateError("INVALID_HOSTNAME", "Hostname is empty")
	}
	
	// Use Go's built-in hostname verification
	err := cert.VerifyHostname(hostname)
	if err != nil {
		// Create a more detailed error for common hostname verification failures
		return NewCertificateError("HOSTNAME_VERIFICATION_FAILED", 
			fmt.Sprintf("Certificate is not valid for hostname '%s'", hostname)).
			WithContext("hostname", hostname).
			WithContext("cert_cn", cert.Subject.CommonName).
			WithContext("cert_sans", strings.Join(cert.DNSNames, ", ")).
			WithSuggestion("Check if the hostname matches the certificate's Common Name or Subject Alternative Names").
			Wrap(err)
	}
	
	return nil
}

// GenerateDiagnostics creates detailed diagnostic information about a certificate
func (v *Validator) GenerateDiagnostics(cert *x509.Certificate) map[string]interface{} {
	if cert == nil {
		return map[string]interface{}{"error": "Certificate is nil"}
	}
	
	now := time.Now()
	diagnostics := map[string]interface{}{
		"subject": cert.Subject.String(),
		"issuer": cert.Issuer.String(),
		"serial_number": cert.SerialNumber.String(),
		"not_before": cert.NotBefore.Format(time.RFC3339),
		"not_after": cert.NotAfter.Format(time.RFC3339),
		"signature_algorithm": cert.SignatureAlgorithm.String(),
		"public_key_algorithm": cert.PublicKeyAlgorithm.String(),
		"is_ca": cert.IsCA,
		"dns_names": cert.DNSNames,
		"ip_addresses": cert.IPAddresses,
		"key_usage": v.keyUsageStrings(cert.KeyUsage),
		"ext_key_usage": v.extKeyUsageStrings(cert.ExtKeyUsage),
		"is_self_signed": cert.Issuer.String() == cert.Subject.String(),
	}
	
	// Add expiry analysis
	if cert.NotAfter.Before(now) {
		diagnostics["status"] = "expired"
		diagnostics["days_expired"] = int(now.Sub(cert.NotAfter).Hours() / 24)
	} else {
		daysUntilExpiry := int(cert.NotAfter.Sub(now).Hours() / 24)
		diagnostics["status"] = "valid"
		diagnostics["days_until_expiry"] = daysUntilExpiry
		if daysUntilExpiry <= 30 {
			diagnostics["expiry_warning"] = "Certificate expires within 30 days"
		}
	}
	
	return diagnostics
}

// keyUsageStrings converts key usage flags to human-readable strings
func (v *Validator) keyUsageStrings(usage x509.KeyUsage) []string {
	usageMap := map[x509.KeyUsage]string{
		x509.KeyUsageDigitalSignature: "DigitalSignature",
		x509.KeyUsageKeyEncipherment: "KeyEncipherment", 
		x509.KeyUsageDataEncipherment: "DataEncipherment",
		x509.KeyUsageKeyAgreement: "KeyAgreement",
		x509.KeyUsageCertSign: "CertSign",
		x509.KeyUsageCRLSign: "CRLSign",
	}
	
	var usages []string
	for flag, name := range usageMap {
		if usage&flag != 0 {
			usages = append(usages, name)
		}
	}
	return usages
}

// extKeyUsageStrings converts extended key usage values to human-readable strings
func (v *Validator) extKeyUsageStrings(usages []x509.ExtKeyUsage) []string {
	usageMap := map[x509.ExtKeyUsage]string{
		x509.ExtKeyUsageServerAuth: "ServerAuth",
		x509.ExtKeyUsageClientAuth: "ClientAuth",
		x509.ExtKeyUsageCodeSigning: "CodeSigning",
		x509.ExtKeyUsageEmailProtection: "EmailProtection",
		x509.ExtKeyUsageTimeStamping: "TimeStamping",
		x509.ExtKeyUsageOCSPSigning: "OCSPSigning",
	}
	
	var result []string
	for _, usage := range usages {
		if name, ok := usageMap[usage]; ok {
			result = append(result, name)
		} else {
			result = append(result, fmt.Sprintf("Unknown(%d)", usage))
		}
	}
	return result
}