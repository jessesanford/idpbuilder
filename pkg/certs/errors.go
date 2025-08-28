package certs

import (
	"fmt"
)

// CertificateError represents certificate-related errors with detailed context
type CertificateError struct {
	Operation string
	Cause     error
	Context   string
	Suggestions []string
}

func (e *CertificateError) Error() string {
	msg := fmt.Sprintf("certificate error in %s: %v", e.Operation, e.Cause)
	if e.Context != "" {
		msg += fmt.Sprintf(" (%s)", e.Context)
	}
	return msg
}

func (e *CertificateError) Unwrap() error {
	return e.Cause
}

// NewCertificateError creates a new certificate error with context
func NewCertificateError(operation string, cause error, context string, suggestions []string) *CertificateError {
	return &CertificateError{
		Operation:   operation,
		Cause:       cause,
		Context:     context,
		Suggestions: suggestions,
	}
}

// Common certificate error types
var (
	ErrClusterNotFound = &CertificateError{
		Operation: "cluster_detection",
		Cause:     fmt.Errorf("kind cluster not found"),
		Context:   "no running Kind cluster detected",
		Suggestions: []string{
			"Start a Kind cluster with: kind create cluster",
			"Verify cluster is running with: kind get clusters",
			"Check if idpbuilder cluster exists: kind get clusters | grep idpbuilder",
		},
	}
	
	ErrGiteaPodNotFound = &CertificateError{
		Operation: "pod_discovery",
		Cause:     fmt.Errorf("gitea pod not found"),
		Context:   "unable to locate Gitea pod in cluster",
		Suggestions: []string{
			"Verify Gitea is installed: kubectl get pods -n gitea",
			"Check pod labels match selector: kubectl get pods -n gitea -l app.kubernetes.io/name=gitea",
			"Ensure idpbuilder was created with --package-dir containing Gitea",
		},
	}
	
	ErrCertificateNotFound = &CertificateError{
		Operation: "certificate_extraction",
		Cause:     fmt.Errorf("certificate file not found"),
		Context:   "certificate file not present in expected location",
		Suggestions: []string{
			"Verify Gitea is configured with HTTPS",
			"Check certificate path in pod: kubectl exec <pod> -- ls -la /data/gitea/https/",
			"Restart Gitea pod to regenerate certificates if needed",
		},
	}
	
	ErrCertificateExpired = &CertificateError{
		Operation: "certificate_validation",
		Cause:     fmt.Errorf("certificate has expired"),
		Context:   "certificate is no longer valid",
		Suggestions: []string{
			"Recreate the cluster to generate new certificates",
			"Use --insecure flag as temporary workaround",
			"Check system clock for time synchronization issues",
		},
	}
	
	ErrCertificateInvalid = &CertificateError{
		Operation: "certificate_validation",
		Cause:     fmt.Errorf("certificate is invalid"),
		Context:   "certificate failed validation checks",
		Suggestions: []string{
			"Verify certificate format and encoding",
			"Check for certificate corruption during extraction",
			"Recreate cluster if certificate generation failed",
		},
	}
	
	ErrStoragePermission = &CertificateError{
		Operation: "certificate_storage",
		Cause:     fmt.Errorf("permission denied"),
		Context:   "unable to write certificate to storage location",
		Suggestions: []string{
			"Check write permissions for ~/.idpbuilder/certs/ directory",
			"Create directory if it doesn't exist: mkdir -p ~/.idpbuilder/certs",
			"Run with appropriate user permissions",
		},
	}
)