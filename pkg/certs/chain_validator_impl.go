package certs

import (
	"context"
	"crypto/x509"
	"fmt"
	"strings"
	"time"
)

// DefaultChainValidator implements the ChainValidator interface
type DefaultChainValidator struct {
	basicValidator CertValidator
	trustManager   TrustManager
	allowSelfSigned bool
}

// ChainValidatorConfig holds configuration for the chain validator
type ChainValidatorConfig struct {
	BasicValidator  CertValidator
	TrustManager    TrustManager
	AllowSelfSigned bool
}

// NewChainValidator creates a new chain validator with the specified configuration
func NewChainValidator(config *ChainValidatorConfig) ChainValidator {
	if config == nil {
		panic("ChainValidatorConfig cannot be nil")
	}
	if config.BasicValidator == nil {
		panic("BasicValidator cannot be nil")
	}
	
	return &DefaultChainValidator{
		basicValidator:  config.BasicValidator,
		trustManager:    config.TrustManager,
		allowSelfSigned: config.AllowSelfSigned,
	}
}

// ValidateChain verifies the complete certificate chain from leaf to root
func (cv *DefaultChainValidator) ValidateChain(ctx context.Context, cert *x509.Certificate, intermediates []*x509.Certificate) (*ChainValidationResult, error) {
	if cert == nil {
		return nil, fmt.Errorf("certificate cannot be nil")
	}

	result := &ChainValidationResult{
		Valid:            true,
		ChainComplete:    false,
		TrustAnchorFound: false,
		ChainLength:      0,
		Certificates:     []CertificateSummary{},
		Issues:           []ValidationIssue{},
	}

	// Build the complete chain
	chain := []*x509.Certificate{cert}
	if intermediates != nil {
		chain = append(chain, intermediates...)
	}
	result.ChainLength = len(chain)

	// Validate each certificate in the chain
	for i, c := range chain {
		summary := CertificateSummary{
			Subject:      c.Subject.String(),
			Issuer:       c.Issuer.String(),
			SerialNumber: c.SerialNumber.String(),
			NotBefore:    c.NotBefore,
			NotAfter:     c.NotAfter,
			IsCA:         c.IsCA,
			Position:     i,
		}
		result.Certificates = append(result.Certificates, summary)

		// Perform basic validation on each certificate
		basicResult, err := cv.basicValidator.ValidateCertificate(c)
		if err != nil {
			result.Valid = false
			result.Issues = append(result.Issues, ValidationIssue{
				Severity:    SeverityError,
				Code:        "BASIC_VALIDATION_FAILED",
				Message:     fmt.Sprintf("Basic validation failed for certificate at position %d: %s", i, err.Error()),
				Certificate: c.Subject.String(),
				Remediation: "Check certificate validity and format",
			})
		} else if !basicResult.Valid {
			result.Valid = false
			for _, issue := range basicResult.Issues {
				result.Issues = append(result.Issues, ValidationIssue{
					Severity:    SeverityError,
					Code:        "BASIC_VALIDATION_ISSUE",
					Message:     fmt.Sprintf("Certificate at position %d: %s", i, issue),
					Certificate: c.Subject.String(),
					Remediation: "Address the basic validation issues",
				})
			}
		}
	}

	// Validate chain ordering and linkage
	if len(chain) > 1 {
		result.ChainComplete = cv.validateChainLinkage(chain, result)
	} else {
		// Single certificate chain
		if cert.Issuer.String() == cert.Subject.String() {
			if cv.allowSelfSigned {
				result.ChainComplete = true
				result.TrustAnchorFound = true
				result.Issues = append(result.Issues, ValidationIssue{
					Severity:    SeverityWarning,
					Code:        "SELF_SIGNED_CERT",
					Message:     "Certificate is self-signed",
					Certificate: cert.Subject.String(),
					Remediation: "Self-signed certificates are acceptable for development environments",
				})
			} else {
				result.Valid = false
				result.Issues = append(result.Issues, ValidationIssue{
					Severity:    SeverityCritical,
					Code:        "SELF_SIGNED_NOT_ALLOWED",
					Message:     "Self-signed certificates are not permitted",
					Certificate: cert.Subject.String(),
					Remediation: "Obtain a certificate from a trusted Certificate Authority",
				})
			}
		} else {
			result.Issues = append(result.Issues, ValidationIssue{
				Severity:    SeverityWarning,
				Code:        "INCOMPLETE_CHAIN",
				Message:     "Certificate chain is incomplete (missing intermediate certificates)",
				Certificate: cert.Subject.String(),
				Remediation: "Provide intermediate certificates to complete the chain",
			})
		}
	}

	// Check trust store if trust manager is available
	if cv.trustManager != nil {
		cv.validateTrustAnchor(ctx, chain, result)
	}

	return result, nil
}

// validateChainLinkage validates that certificates in the chain are properly linked
func (cv *DefaultChainValidator) validateChainLinkage(chain []*x509.Certificate, result *ChainValidationResult) bool {
	chainComplete := true
	
	for i := 0; i < len(chain)-1; i++ {
		child := chain[i]
		parent := chain[i+1]
		
		// Verify that parent issued child
		if child.Issuer.String() != parent.Subject.String() {
			chainComplete = false
			result.Valid = false
			result.Issues = append(result.Issues, ValidationIssue{
				Severity:    SeverityError,
				Code:        "CHAIN_LINKAGE_ERROR",
				Message:     fmt.Sprintf("Certificate at position %d is not issued by certificate at position %d", i, i+1),
				Certificate: child.Subject.String(),
				Remediation: "Ensure certificates are provided in the correct order (leaf to root)",
			})
		}
		
		// Verify signature
		if err := child.CheckSignatureFrom(parent); err != nil {
			chainComplete = false
			result.Valid = false
			result.Issues = append(result.Issues, ValidationIssue{
				Severity:    SeverityCritical,
				Code:        "SIGNATURE_VERIFICATION_FAILED",
				Message:     fmt.Sprintf("Signature verification failed for certificate at position %d", i),
				Certificate: child.Subject.String(),
				Remediation: "Verify that the certificate chain is not corrupted",
			})
		}
	}
	
	return chainComplete
}

// validateTrustAnchor checks if the chain terminates in a trusted root
func (cv *DefaultChainValidator) validateTrustAnchor(ctx context.Context, chain []*x509.Certificate, result *ChainValidationResult) {
	// Check the root certificate against the trust store
	root := chain[len(chain)-1]
	
	// For now, we'll use a simple approach - check if root is self-signed as trust anchor
	if root.Issuer.String() == root.Subject.String() {
		result.TrustAnchorFound = true
		result.Issues = append(result.Issues, ValidationIssue{
			Severity:    SeverityWarning,
			Code:        "SELF_SIGNED_ROOT",
			Message:     "Chain terminates in a self-signed root certificate",
			Certificate: root.Subject.String(),
			Remediation: "For production use, ensure root certificate is from a trusted CA",
		})
	} else {
		result.Issues = append(result.Issues, ValidationIssue{
			Severity:    SeverityWarning,
			Code:        "NO_TRUSTED_ROOT",
			Message:     "Chain does not terminate in a trusted root certificate",
			Certificate: root.Subject.String(),
			Remediation: "Provide a complete certificate chain including the root CA",
		})
	}
}

// VerifyHostname checks if the certificate is valid for the given hostname
func (cv *DefaultChainValidator) VerifyHostname(cert *x509.Certificate, hostname string) error {
	if cert == nil {
		return fmt.Errorf("certificate cannot be nil")
	}
	if hostname == "" {
		return fmt.Errorf("hostname cannot be empty")
	}

	// Use Go's built-in hostname verification
	err := cert.VerifyHostname(hostname)
	if err != nil {
		// Create a more detailed error with context
		return fmt.Errorf("hostname verification failed for '%s': certificate is valid for %v but not for %s", 
			hostname, append([]string{cert.Subject.CommonName}, cert.DNSNames...), hostname)
	}

	return nil
}

// CheckChainExpiry verifies no certificates in the chain are expired or expiring soon
func (cv *DefaultChainValidator) CheckChainExpiry(chain []*x509.Certificate, warnDays int) (*ChainExpiryResult, error) {
	if len(chain) == 0 {
		return nil, fmt.Errorf("certificate chain cannot be empty")
	}
	
	if warnDays <= 0 {
		warnDays = 30 // Default warning period
	}
	
	now := time.Now()
	result := &ChainExpiryResult{
		ChainValid:      true,
		ExpiringCerts:   []ExpiringCertificate{},
		ExpiredCerts:    []ExpiredCertificate{},
		MinDaysToExpiry: int(chain[0].NotAfter.Sub(now).Hours() / 24),
	}
	
	for i, cert := range chain {
		daysUntilExpiry := int(cert.NotAfter.Sub(now).Hours() / 24)
		
		if cert.NotAfter.Before(now) {
			// Certificate has expired
			result.ChainValid = false
			daysExpired := int(now.Sub(cert.NotAfter).Hours() / 24)
			result.ExpiredCerts = append(result.ExpiredCerts, ExpiredCertificate{
				Subject:     cert.Subject.String(),
				Position:    i,
				DaysExpired: daysExpired,
				ExpiredDate: cert.NotAfter,
			})
		} else if daysUntilExpiry <= warnDays {
			// Certificate is expiring soon
			result.ExpiringCerts = append(result.ExpiringCerts, ExpiringCertificate{
				Subject:         cert.Subject.String(),
				Position:        i,
				DaysUntilExpiry: daysUntilExpiry,
				ExpiryDate:      cert.NotAfter,
			})
		}
		
		// Track minimum days to expiry
		if daysUntilExpiry < result.MinDaysToExpiry {
			result.MinDaysToExpiry = daysUntilExpiry
		}
	}
	
	return result, nil
}

// GenerateDiagnostics creates a comprehensive diagnostic report
func (cv *DefaultChainValidator) GenerateDiagnostics(ctx context.Context, cert *x509.Certificate, hostname string) (*CertDiagnosticsReport, error) {
	if cert == nil {
		return nil, fmt.Errorf("certificate cannot be nil")
	}

	report := &CertDiagnosticsReport{
		Timestamp: time.Now(),
		Hostname:  hostname,
		Recommendations: []Recommendation{},
	}

	// Generate certificate details
	report.CertificateDetails = cv.generateCertificateDetails(cert)

	// Perform chain analysis with just the single certificate
	chainResult, err := cv.ValidateChain(ctx, cert, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to validate chain: %w", err)
	}
	
	report.ChainAnalysis = &ChainAnalysis{
		ChainLength:      chainResult.ChainLength,
		ChainValid:       chainResult.Valid,
		TrustAnchorFound: chainResult.TrustAnchorFound,
		Certificates:     chainResult.Certificates,
		ChainIssues:      chainResult.Issues,
	}

	// Perform hostname validation if hostname provided
	if hostname != "" {
		hostnameErr := cv.VerifyHostname(cert, hostname)
		report.HostnameValidation = &HostnameValidation{
			Hostname: hostname,
			Valid:    hostnameErr == nil,
		}
		
		if hostnameErr == nil {
			report.HostnameValidation.MatchType = cv.determineMatchType(cert, hostname)
			report.HostnameValidation.MatchedValue = cv.findMatchedValue(cert, hostname)
		} else {
			report.HostnameValidation.Error = hostnameErr.Error()
		}
	}

	// Generate trust store analysis
	report.TrustStoreAnalysis = &TrustStoreAnalysis{
		TrustStoreChecked:  cv.trustManager != nil,
		CertificateTrusted: false,
	}

	// Generate recommendations based on findings
	report.Recommendations = cv.generateRecommendations(report)

	return report, nil
}

// generateCertificateDetails creates detailed certificate information
func (cv *DefaultChainValidator) generateCertificateDetails(cert *x509.Certificate) *CertificateDetails {
	ipAddresses := make([]string, len(cert.IPAddresses))
	for i, ip := range cert.IPAddresses {
		ipAddresses[i] = ip.String()
	}

	// Get basic validation result
	basicResult, _ := cv.basicValidator.ValidateCertificate(cert)

	return &CertificateDetails{
		Subject:            cert.Subject.String(),
		Issuer:             cert.Issuer.String(),
		SerialNumber:       cert.SerialNumber.String(),
		NotBefore:          cert.NotBefore,
		NotAfter:           cert.NotAfter,
		IsCA:               cert.IsCA,
		IsSelfSigned:       cert.Issuer.String() == cert.Subject.String(),
		DNSNames:           cert.DNSNames,
		IPAddresses:        ipAddresses,
		KeyUsage:           cv.keyUsageStrings(cert.KeyUsage),
		ExtKeyUsage:        cv.extKeyUsageStrings(cert.ExtKeyUsage),
		SignatureAlgorithm: cert.SignatureAlgorithm.String(),
		PublicKeyAlgorithm: cert.PublicKeyAlgorithm.String(),
		ValidationResult:   basicResult,
	}
}

// determineMatchType determines how the hostname matches the certificate
func (cv *DefaultChainValidator) determineMatchType(cert *x509.Certificate, hostname string) string {
	// Check exact CN match
	if cert.Subject.CommonName == hostname {
		return "exact"
	}
	
	// Check SAN matches
	for _, dnsName := range cert.DNSNames {
		if dnsName == hostname {
			return "san"
		}
		if strings.HasPrefix(dnsName, "*.") && cv.matchesWildcard(dnsName, hostname) {
			return "wildcard"
		}
	}
	
	// Check wildcard CN
	if strings.HasPrefix(cert.Subject.CommonName, "*.") && cv.matchesWildcard(cert.Subject.CommonName, hostname) {
		return "wildcard"
	}
	
	return "none"
}

// findMatchedValue finds which certificate value matched the hostname
func (cv *DefaultChainValidator) findMatchedValue(cert *x509.Certificate, hostname string) string {
	if cert.Subject.CommonName == hostname {
		return cert.Subject.CommonName
	}
	
	for _, dnsName := range cert.DNSNames {
		if dnsName == hostname || (strings.HasPrefix(dnsName, "*.") && cv.matchesWildcard(dnsName, hostname)) {
			return dnsName
		}
	}
	
	if strings.HasPrefix(cert.Subject.CommonName, "*.") && cv.matchesWildcard(cert.Subject.CommonName, hostname) {
		return cert.Subject.CommonName
	}
	
	return ""
}

// matchesWildcard checks if hostname matches a wildcard pattern
func (cv *DefaultChainValidator) matchesWildcard(pattern, hostname string) bool {
	if !strings.HasPrefix(pattern, "*.") {
		return false
	}
	
	domain := pattern[2:] // Remove "*."
	return strings.HasSuffix(hostname, "."+domain) || hostname == domain
}

// generateRecommendations creates actionable recommendations based on diagnostic findings
func (cv *DefaultChainValidator) generateRecommendations(report *CertDiagnosticsReport) []Recommendation {
	recommendations := []Recommendation{}

	// Check for expired certificates
	if report.CertificateDetails.ValidationResult != nil && len(report.CertificateDetails.ValidationResult.Issues) > 0 {
		recommendations = append(recommendations, Recommendation{
			Priority:    PriorityHigh,
			Title:       "Certificate Validation Issues",
			Description: "The certificate has validation issues that need to be addressed",
			Command:     "",
			Link:        "",
		})
	}

	// Check for hostname validation issues
	if report.HostnameValidation != nil && !report.HostnameValidation.Valid {
		recommendations = append(recommendations, Recommendation{
			Priority:    PriorityCritical,
			Title:       "Hostname Verification Failed",
			Description: fmt.Sprintf("Certificate is not valid for hostname '%s'", report.HostnameValidation.Hostname),
			Command:     "",
			Link:        "",
		})
	}

	return recommendations
}

// keyUsageStrings converts key usage flags to human-readable strings
func (cv *DefaultChainValidator) keyUsageStrings(usage x509.KeyUsage) []string {
	usageMap := map[x509.KeyUsage]string{
		x509.KeyUsageDigitalSignature: "DigitalSignature",
		x509.KeyUsageKeyEncipherment:  "KeyEncipherment",
		x509.KeyUsageDataEncipherment: "DataEncipherment",
		x509.KeyUsageKeyAgreement:     "KeyAgreement",
		x509.KeyUsageCertSign:         "CertSign",
		x509.KeyUsageCRLSign:          "CRLSign",
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
func (cv *DefaultChainValidator) extKeyUsageStrings(usages []x509.ExtKeyUsage) []string {
	usageMap := map[x509.ExtKeyUsage]string{
		x509.ExtKeyUsageServerAuth:      "ServerAuth",
		x509.ExtKeyUsageClientAuth:      "ClientAuth",
		x509.ExtKeyUsageCodeSigning:     "CodeSigning",
		x509.ExtKeyUsageEmailProtection: "EmailProtection",
		x509.ExtKeyUsageTimeStamping:    "TimeStamping",
		x509.ExtKeyUsageOCSPSigning:     "OCSPSigning",
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