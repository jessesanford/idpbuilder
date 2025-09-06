// Package fallback provides certificate problem detection and fallback strategies
// for IDPBuilder OCI registry operations with comprehensive error analysis
package fallback

import (
	"crypto/x509"
	"fmt"
	"regexp"
	"strings"

	"github.com/cnoe-io/idpbuilder/pkg/certs"
)

// ProblemType represents specific certificate validation issues
type ProblemType string

const (
	ProblemSelfSigned       ProblemType = "self-signed"
	ProblemExpired          ProblemType = "expired"
	ProblemNotYetValid      ProblemType = "not-yet-valid"
	ProblemHostnameMismatch ProblemType = "hostname-mismatch"
	ProblemUntrustedCA      ProblemType = "untrusted-ca"
	ProblemUnknownAuthority ProblemType = "unknown-authority"
	ProblemUnknown          ProblemType = "unknown"
)

// CertProblem contains detailed information about a certificate validation issue
type CertProblem struct {
	Type        ProblemType                // The specific problem identified
	Certificate *x509.Certificate          // The problematic certificate
	Error       error                      // Original validation error
	Details     map[string]interface{}     // Additional context
	Suggestions []string                   // Quick fix suggestions
}

// ProblemDetector analyzes certificate validation errors and identifies specific issues
type ProblemDetector interface {
	// DetectProblem analyzes a validation error and identifies the specific problem type
	DetectProblem(validationErr error, cert *x509.Certificate) (*CertProblem, error)
	
	// AnalyzeCertChain provides comprehensive analysis of certificate chain issues
	AnalyzeCertChain(certs []*x509.Certificate) ([]*CertProblem, error)
	
	// DetectFromDiagnostics analyzes problems from CertDiagnostics output
	DetectFromDiagnostics(diag *certs.CertDiagnostics) ([]*CertProblem, error)
}

// DefaultDetector implements ProblemDetector with integration to CertValidator from E1.2.1
type DefaultDetector struct {
	validator certs.CertValidator // CertValidator from E1.2.1
}

// NewDetector creates a detector integrated with the CertValidator from E1.2.1
func NewDetector(validator certs.CertValidator) *DefaultDetector {
	return &DefaultDetector{
		validator: validator,
	}
}

// DetectProblem analyzes validation errors using pattern matching to identify specific issues
func (d *DefaultDetector) DetectProblem(validationErr error, cert *x509.Certificate) (*CertProblem, error) {
	if validationErr == nil {
		return nil, nil // No error, no problem
	}
	
	if cert == nil {
		return nil, fmt.Errorf("certificate cannot be nil")
	}

	errStr := validationErr.Error()
	errStrLower := strings.ToLower(errStr)

	problem := &CertProblem{
		Certificate: cert,
		Error:       validationErr,
		Details:     make(map[string]interface{}),
		Suggestions: make([]string, 0),
	}

	// Pattern matching for different x509 error types
	switch {
	case strings.Contains(errStrLower, "certificate has expired") || strings.Contains(errStrLower, "expired"):
		problem.Type = ProblemExpired
		problem.Details["expired_on"] = cert.NotAfter
		problem.Suggestions = append(problem.Suggestions, "Renew the certificate", "Use --insecure flag for testing")
		
	case strings.Contains(errStrLower, "certificate is not valid until") || strings.Contains(errStrLower, "not yet valid"):
		problem.Type = ProblemNotYetValid
		problem.Details["valid_from"] = cert.NotBefore
		problem.Suggestions = append(problem.Suggestions, "Check system clock synchronization")
		
	case strings.Contains(errStrLower, "certificate signed by unknown authority") || strings.Contains(errStrLower, "unknown authority"):
		// Could be self-signed or untrusted CA
		if cert.Subject.String() == cert.Issuer.String() {
			problem.Type = ProblemSelfSigned
			problem.Details["issuer"] = cert.Issuer.String()
			problem.Suggestions = append(problem.Suggestions, 
				"Add certificate to trust store", 
				"Use --insecure flag for development",
				"Import the CA certificate")
		} else {
			problem.Type = ProblemUntrustedCA
			problem.Details["ca_issuer"] = cert.Issuer.String()
			problem.Suggestions = append(problem.Suggestions,
				"Import the root CA certificate",
				"Add intermediate CA certificates",
				"Use --insecure flag for testing")
		}
		
	case strings.Contains(errStrLower, "hostname") && (strings.Contains(errStrLower, "doesn't match") || strings.Contains(errStrLower, "does not match")):
		problem.Type = ProblemHostnameMismatch
		problem.Details["valid_hostnames"] = cert.DNSNames
		problem.Details["common_name"] = cert.Subject.CommonName
		
		// Extract requested hostname from error if possible
		hostnameRegex := regexp.MustCompile(`hostname ['\"]?([^'\"]+)['\"]?`)
		if matches := hostnameRegex.FindStringSubmatch(errStr); len(matches) > 1 {
			problem.Details["requested_hostname"] = matches[1]
		}
		
		problem.Suggestions = append(problem.Suggestions,
			"Use the correct hostname from certificate",
			"Update certificate with correct Subject Alternative Names",
			"Add hostname to /etc/hosts for testing")
			
	default:
		problem.Type = ProblemUnknown
		problem.Details["error_type"] = fmt.Sprintf("%T", validationErr)
		problem.Suggestions = append(problem.Suggestions,
			"Check certificate validity manually",
			"Enable debug logging for more details")
	}

	return problem, nil
}

// AnalyzeCertChain examines multiple certificates and identifies chain-specific issues
func (d *DefaultDetector) AnalyzeCertChain(certs []*x509.Certificate) ([]*CertProblem, error) {
	if len(certs) == 0 {
		return nil, fmt.Errorf("certificate chain cannot be empty")
	}

	problems := make([]*CertProblem, 0)

	// Analyze each certificate in the chain
	for i, cert := range certs {
		// Try validation through our validator
		err := d.validator.ValidateChain(cert)
		if err != nil {
			problem, detectionErr := d.DetectProblem(err, cert)
			if detectionErr != nil {
				return nil, fmt.Errorf("failed to detect problem for cert %d: %w", i, detectionErr)
			}
			if problem != nil {
				problem.Details["chain_position"] = i
				problem.Details["total_certs"] = len(certs)
				problems = append(problems, problem)
			}
		}
	}

	return problems, nil
}

// DetectFromDiagnostics analyzes problems from CertDiagnostics generated by E1.2.1
func (d *DefaultDetector) DetectFromDiagnostics(diag *certs.CertDiagnostics) ([]*CertProblem, error) {
	if diag == nil {
		return nil, fmt.Errorf("diagnostics cannot be nil")
	}

	problems := make([]*CertProblem, 0)

	// Convert each ValidationError from diagnostics to a CertProblem
	for _, valErr := range diag.ValidationErrors {
		problem := &CertProblem{
			Error:   fmt.Errorf("%s: %s", valErr.Message, valErr.Detail),
			Details: make(map[string]interface{}),
			Suggestions: make([]string, 0),
		}

		// Map validation error type to problem type
		switch valErr.Type {
		case "chain":
			if strings.Contains(strings.ToLower(valErr.Detail), "self-signed") {
				problem.Type = ProblemSelfSigned
			} else if strings.Contains(strings.ToLower(valErr.Detail), "unknown authority") {
				problem.Type = ProblemUntrustedCA
			} else {
				problem.Type = ProblemUnknownAuthority
			}
		case "expiry":
			if strings.Contains(strings.ToLower(valErr.Detail), "expired") {
				problem.Type = ProblemExpired
			} else {
				problem.Type = ProblemNotYetValid
			}
		case "hostname":
			problem.Type = ProblemHostnameMismatch
		default:
			problem.Type = ProblemUnknown
		}

		// Add diagnostic context to details
		problem.Details["subject"] = diag.Subject
		problem.Details["issuer"] = diag.Issuer
		problem.Details["dns_names"] = diag.DNSNames
		problem.Details["not_before"] = diag.NotBefore
		problem.Details["not_after"] = diag.NotAfter

		problems = append(problems, problem)
	}

	return problems, nil
}

// GetProblemSummary returns a human-readable summary of the certificate problem
func (cp *CertProblem) GetProblemSummary() string {
	switch cp.Type {
	case ProblemSelfSigned:
		return "Certificate is self-signed and not trusted by system"
	case ProblemExpired:
		return "Certificate has expired"
	case ProblemNotYetValid:
		return "Certificate is not yet valid"
	case ProblemHostnameMismatch:
		return "Certificate hostname does not match requested hostname"
	case ProblemUntrustedCA:
		return "Certificate was issued by an untrusted Certificate Authority"
	case ProblemUnknownAuthority:
		return "Certificate authority is not recognized"
	default:
		return "Unknown certificate validation problem"
	}
}

// GetDetailedDescription provides comprehensive information about the problem
func (cp *CertProblem) GetDetailedDescription() string {
	summary := cp.GetProblemSummary()
	
	switch cp.Type {
	case ProblemSelfSigned:
		if issuer, ok := cp.Details["issuer"].(string); ok {
			return fmt.Sprintf("%s. Issuer: %s. This is common in development environments and local Kind clusters.", summary, issuer)
		}
	case ProblemExpired:
		if expiredOn, ok := cp.Details["expired_on"]; ok {
			return fmt.Sprintf("%s on %v. The certificate needs to be renewed.", summary, expiredOn)
		}
	case ProblemHostnameMismatch:
		if validHosts, ok := cp.Details["valid_hostnames"].([]string); ok && len(validHosts) > 0 {
			return fmt.Sprintf("%s. Certificate is valid for: %s", summary, strings.Join(validHosts, ", "))
		}
		if cn, ok := cp.Details["common_name"].(string); ok && cn != "" {
			return fmt.Sprintf("%s. Certificate Common Name: %s", summary, cn)
		}
	}
	
	return summary
}