package errors

import (
	"fmt"
	"strings"
)

// ResolutionGuide provides step-by-step guidance for resolving certificate errors
type ResolutionGuide struct {
	Steps      []string              // Ordered resolution steps
	Examples   []string              // Command examples
	DocLinks   []string              // Documentation URLs
	Workaround string                // Quick workaround if available
	Severity   string                // How urgent is this?
	AutoFix    func() error          // Optional auto-fix function
}

// Format returns a human-readable representation of the resolution guide
func (rg *ResolutionGuide) Format() string {
	var sb strings.Builder
	
	sb.WriteString("Resolution Steps:\n")
	for i, step := range rg.Steps {
		sb.WriteString(fmt.Sprintf("  %d. %s\n", i+1, step))
	}
	
	if len(rg.Examples) > 0 {
		sb.WriteString("\nExample Commands:\n")
		for _, example := range rg.Examples {
			sb.WriteString(fmt.Sprintf("  $ %s\n", example))
		}
	}
	
	if rg.Workaround != "" {
		sb.WriteString(fmt.Sprintf("\nWorkaround: %s\n", rg.Workaround))
	}
	
	if len(rg.DocLinks) > 0 {
		sb.WriteString("\nDocumentation:\n")
		for _, link := range rg.DocLinks {
			sb.WriteString(fmt.Sprintf("  %s\n", link))
		}
	}
	
	return sb.String()
}

// resolutionGuides maps error types to their resolution guidance
var resolutionGuides = map[ErrorType]*ResolutionGuide{
	ErrCertNotFound: {
		Steps: []string{
			"Extract certificates from the cluster using kubectl",
			"Verify the certificate path in your configuration",
			"Check if the certificate was generated during cluster setup",
			"Ensure the certificate secret exists in the correct namespace",
		},
		Examples: []string{
			"kubectl get secrets -n gitea | grep tls",
			"kubectl get secret gitea-tls-cert -n gitea -o yaml",
			"ls -la /tmp/certs/",
		},
		DocLinks: []string{
			"https://github.com/idpbuilder/docs/certificate-extraction",
		},
		Workaround: "Use --insecure flag for development (NOT for production)",
		Severity:   "HIGH",
	},
	
	ErrCertInvalid: {
		Steps: []string{
			"Verify certificate format and encoding",
			"Check if certificate is properly base64 encoded",
			"Validate certificate using OpenSSL tools",
			"Regenerate certificate if validation fails",
		},
		Examples: []string{
			"openssl x509 -in /tmp/certs/registry.crt -text -noout",
			"kubectl delete secret -n gitea gitea-tls-cert",
			"kubectl rollout restart deployment -n gitea gitea",
		},
		DocLinks: []string{
			"https://github.com/idpbuilder/docs/certificate-validation",
		},
		Workaround: "Generate new certificate with correct format",
		Severity:   "HIGH",
	},
	
	ErrCertExpired: {
		Steps: []string{
			"Regenerate certificates in the cluster",
			"Update certificate rotation configuration",
			"Restart affected services to pick up new certificates",
			"Verify new certificate expiration dates",
		},
		Examples: []string{
			"kubectl delete secret -n gitea gitea-tls-cert",
			"kubectl rollout restart deployment -n gitea gitea",
			"openssl x509 -in /tmp/certs/registry.crt -dates -noout",
		},
		DocLinks: []string{
			"https://github.com/idpbuilder/docs/cert-rotation",
		},
		Workaround: "Use --insecure flag temporarily for development",
		Severity:   "HIGH",
	},
	
	ErrCertUntrusted: {
		Steps: []string{
			"Add certificate to system trust store",
			"Verify certificate chain completeness",
			"Check if CA certificate is properly installed",
			"Update application trust configuration",
		},
		Examples: []string{
			"sudo cp /tmp/certs/ca.crt /usr/local/share/ca-certificates/",
			"sudo update-ca-certificates",
			"curl -v --cacert /tmp/certs/ca.crt https://registry.local",
		},
		DocLinks: []string{
			"https://github.com/idpbuilder/docs/trust-store",
		},
		Workaround: "Use --insecure or --cacert flag with explicit CA",
		Severity:   "MEDIUM",
	},
	
	ErrCertMismatch: {
		Steps: []string{
			"Verify registry URL matches certificate Common Name",
			"Update certificate with correct domain names",
			"Check DNS resolution for registry hostname",
			"Update /etc/hosts if using local registry",
		},
		Examples: []string{
			"openssl x509 -in /tmp/certs/registry.crt -text | grep -A1 'Subject:'",
			"echo '127.0.0.1 registry.local' >> /etc/hosts",
			"nslookup registry.local",
		},
		DocLinks: []string{
			"https://github.com/idpbuilder/docs/certificate-domains",
		},
		Workaround: "Use IP address instead of hostname",
		Severity:   "MEDIUM",
	},
	
	ErrCertPermission: {
		Steps: []string{
			"Check file permissions on certificate files",
			"Verify user has read access to certificate directory",
			"Ensure certificate files are not owned by root",
			"Update file permissions if necessary",
		},
		Examples: []string{
			"ls -la /tmp/certs/",
			"chmod 644 /tmp/certs/*.crt",
			"chown $(whoami) /tmp/certs/*",
		},
		DocLinks: []string{
			"https://github.com/idpbuilder/docs/certificate-permissions",
		},
		Workaround: "Copy certificates to user-accessible location",
		Severity:   "LOW",
	},
	
	ErrCertChainBroken: {
		Steps: []string{
			"Get complete certificate chain from cluster",
			"Verify intermediate certificates are present",
			"Rebuild certificate bundle with all chain components",
			"Test certificate chain validation",
		},
		Examples: []string{
			"kubectl get secret gitea-tls-cert -n gitea -o jsonpath='{.data.tls\\.crt}' | base64 -d",
			"openssl verify -CAfile /tmp/certs/ca.crt /tmp/certs/registry.crt",
			"curl -v --cacert /tmp/certs/ca.crt https://registry.local",
		},
		DocLinks: []string{
			"https://github.com/idpbuilder/docs/certificate-chains",
		},
		Workaround: "Use CA certificate directly for validation",
		Severity:   "HIGH",
	},
	
	ErrCertFormat: {
		Steps: []string{
			"Convert certificate to the expected format",
			"Check certificate encoding (PEM vs DER)",
			"Verify certificate file extension matches format",
			"Use appropriate tools for format conversion",
		},
		Examples: []string{
			"openssl x509 -in cert.der -inform DER -out cert.pem -outform PEM",
			"openssl x509 -in cert.pem -inform PEM -out cert.der -outform DER",
			"file /tmp/certs/registry.crt",
		},
		DocLinks: []string{
			"https://github.com/idpbuilder/docs/certificate-formats",
		},
		Workaround: "Use openssl to convert between formats",
		Severity:   "LOW",
	},
}

// GetResolutionGuide returns the resolution guide for a specific error type
func GetResolutionGuide(errorType ErrorType) (*ResolutionGuide, bool) {
	guide, exists := resolutionGuides[errorType]
	return guide, exists
}

// GetResolutionGuideForError returns the resolution guide for a CertificateError
func GetResolutionGuideForError(err *CertificateError) (*ResolutionGuide, bool) {
	return GetResolutionGuide(err.Type)
}

// UpdateErrorWithGuidance adds resolution guidance to a CertificateError
func UpdateErrorWithGuidance(err *CertificateError) {
	if guide, exists := GetResolutionGuide(err.Type); exists {
		if err.Resolution == "" {
			err.Resolution = guide.Format()
		}
	}
}

// GetAllResolutionGuides returns all available resolution guides
func GetAllResolutionGuides() map[ErrorType]*ResolutionGuide {
	// Return a copy to prevent external modification
	result := make(map[ErrorType]*ResolutionGuide)
	for k, v := range resolutionGuides {
		result[k] = v
	}
	return result
}

// RegisterResolutionGuide allows registration of custom resolution guides
func RegisterResolutionGuide(errorType ErrorType, guide *ResolutionGuide) {
	resolutionGuides[errorType] = guide
}