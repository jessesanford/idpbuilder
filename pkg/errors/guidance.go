package errors

import (
	"fmt"
	"strings"
)

// resolutionSteps provides basic resolution guidance for each error type
var resolutionSteps = map[ErrorType][]string{
	ErrCertNotFound: {
		"Extract certificates from cluster: kubectl get secrets -n gitea | grep tls",
		"Verify certificate path in configuration",
		"Use --insecure flag for development (NOT production)",
	},
	ErrCertExpired: {
		"Regenerate certificates: kubectl delete secret -n gitea gitea-tls-cert",
		"Restart services: kubectl rollout restart deployment -n gitea gitea",
		"Use --insecure flag temporarily for development",
	},
	ErrCertUntrusted: {
		"Add CA to trust store: sudo cp ca.crt /usr/local/share/ca-certificates/",
		"Update certificates: sudo update-ca-certificates",
		"Use --cacert flag with explicit CA",
	},
	ErrCertMismatch: {
		"Check certificate CN: openssl x509 -in cert.crt -text | grep Subject",
		"Update /etc/hosts for local registry",
		"Use IP address instead of hostname",
	},
	ErrCertPermission: {
		"Fix file permissions: chmod 644 /tmp/certs/*.crt",
		"Change ownership: chown $(whoami) /tmp/certs/*",
		"Copy certificates to accessible location",
	},
}

// GetResolutionSteps returns resolution steps for an error type
func GetResolutionSteps(errorType ErrorType) []string {
	steps, exists := resolutionSteps[errorType]
	if !exists {
		return []string{"No resolution steps available for this error type"}
	}
	return steps
}

// FormatResolution formats resolution steps for display
func FormatResolution(errorType ErrorType) string {
	steps := GetResolutionSteps(errorType)
	var sb strings.Builder
	
	sb.WriteString("Resolution Steps:\n")
	for i, step := range steps {
		sb.WriteString(fmt.Sprintf("  %d. %s\n", i+1, step))
	}
	
	return sb.String()
}

// AddResolutionToError adds formatted resolution guidance to an error
func AddResolutionToError(err *CertificateError) {
	resolution := FormatResolution(err.Type)
	err.Details["resolution"] = resolution
}