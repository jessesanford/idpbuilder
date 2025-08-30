package certs

import (
	"errors"
	"fmt"
	"strings"
)

// CertificateError represents a certificate extraction or validation error with context
type CertificateError struct {
	// Code is a unique error code for programmatic handling
	Code string
	
	// Message is the human-readable error message
	Message string
	
	// Context provides additional context about the error
	Context map[string]interface{}
	
	// Suggestions provides actionable suggestions to resolve the error
	Suggestions []string
	
	// Underlying is the wrapped underlying error if any
	Underlying error
}

// Error implements the error interface
func (e *CertificateError) Error() string {
	var parts []string
	
	if e.Code != "" {
		parts = append(parts, fmt.Sprintf("[%s]", e.Code))
	}
	
	parts = append(parts, e.Message)
	
	if len(e.Context) > 0 {
		contextParts := make([]string, 0, len(e.Context))
		for key, value := range e.Context {
			contextParts = append(contextParts, fmt.Sprintf("%s=%v", key, value))
		}
		parts = append(parts, fmt.Sprintf("context: %s", strings.Join(contextParts, ", ")))
	}
	
	if e.Underlying != nil {
		parts = append(parts, fmt.Sprintf("underlying: %v", e.Underlying))
	}
	
	return strings.Join(parts, " | ")
}

// Unwrap returns the underlying error for error unwrapping
func (e *CertificateError) Unwrap() error {
	return e.Underlying
}

// WithContext adds context information to the error
func (e *CertificateError) WithContext(key string, value interface{}) *CertificateError {
	if e.Context == nil {
		e.Context = make(map[string]interface{})
	}
	e.Context[key] = value
	return e
}

// WithSuggestion adds a suggestion to resolve the error
func (e *CertificateError) WithSuggestion(suggestion string) *CertificateError {
	e.Suggestions = append(e.Suggestions, suggestion)
	return e
}

// Wrap wraps an existing error as the underlying error
func (e *CertificateError) Wrap(err error) *CertificateError {
	e.Underlying = err
	return e
}

// NewCertificateError creates a new CertificateError
func NewCertificateError(code, message string) *CertificateError {
	return &CertificateError{
		Code:    code,
		Message: message,
		Context: make(map[string]interface{}),
	}
}

// Predefined error instances for common scenarios
var (
	// ErrClusterNotFound indicates the specified Kind cluster was not found
	ErrClusterNotFound = &CertificateError{
		Code:    "CLUSTER_NOT_FOUND",
		Message: "Kind cluster not found or not accessible",
		Suggestions: []string{
			"Verify the cluster name is correct",
			"Ensure the Kind cluster is running with 'kind get clusters'",
			"Check your kubeconfig context with 'kubectl config current-context'",
		},
	}
	
	// ErrClusterConnection indicates failure to connect to the cluster
	ErrClusterConnection = &CertificateError{
		Code:    "CLUSTER_CONNECTION_FAILED",
		Message: "Failed to establish connection to Kind cluster",
		Suggestions: []string{
			"Check if the cluster is running and accessible",
			"Verify your kubeconfig is properly configured",
			"Try running 'kubectl get nodes' to test connectivity",
		},
	}
	
	// ErrGiteaPodNotFound indicates no Gitea pods were found
	ErrGiteaPodNotFound = &CertificateError{
		Code:    "GITEA_POD_NOT_FOUND",
		Message: "No Gitea pods found in the cluster",
		Suggestions: []string{
			"Check if Gitea is deployed with 'kubectl get pods -n gitea'",
			"Verify the namespace and pod selector configuration",
			"Ensure idpbuilder has been run successfully",
		},
	}
	
	// ErrMultipleGiteaPods indicates multiple Gitea pods were found
	ErrMultipleGiteaPods = &CertificateError{
		Code:    "MULTIPLE_GITEA_PODS",
		Message: "Multiple Gitea pods found, cannot determine which to use",
		Suggestions: []string{
			"Refine the pod selector to match a single pod",
			"Check pod labels with 'kubectl describe pods -n gitea'",
			"Consider using a more specific selector",
		},
	}
	
	// ErrCertificateNotFound indicates the certificate file was not found in the pod
	ErrCertificateNotFound = &CertificateError{
		Code:    "CERTIFICATE_NOT_FOUND",
		Message: "Certificate file not found in Gitea pod",
		Suggestions: []string{
			"Verify the certificate path configuration",
			"Check if Gitea TLS is properly configured",
			"Look for certificates in alternative paths",
		},
	}
	
	// ErrCertificateRead indicates failure to read the certificate from the pod
	ErrCertificateRead = &CertificateError{
		Code:    "CERTIFICATE_READ_FAILED",
		Message: "Failed to read certificate from Gitea pod",
		Suggestions: []string{
			"Check pod permissions and file accessibility",
			"Verify kubectl exec permissions",
			"Ensure the certificate file is readable",
		},
	}
	
	// ErrCertificateParse indicates failure to parse the certificate
	ErrCertificateParse = &CertificateError{
		Code:    "CERTIFICATE_PARSE_FAILED",
		Message: "Failed to parse certificate data",
		Suggestions: []string{
			"Verify the certificate is in valid PEM format",
			"Check if the certificate data is complete",
			"Ensure the certificate is not corrupted",
		},
	}
	
	// ErrCertificateStore indicates failure to store the certificate locally
	ErrCertificateStore = &CertificateError{
		Code:    "CERTIFICATE_STORE_FAILED",
		Message: "Failed to store certificate to local filesystem",
		Suggestions: []string{
			"Check write permissions to the output directory",
			"Ensure the output directory exists or can be created",
			"Verify sufficient disk space is available",
		},
	}
)

// WrapError wraps a standard error with additional certificate-specific context
func WrapError(err error, code, message string) *CertificateError {
	return &CertificateError{
		Code:       code,
		Message:    message,
		Underlying: err,
		Context:    make(map[string]interface{}),
	}
}

// IsErrorCode checks if an error matches a specific certificate error code
func IsErrorCode(err error, code string) bool {
	var certErr *CertificateError
	if !errors.As(err, &certErr) {
		return false
	}
	return certErr.Code == code
}