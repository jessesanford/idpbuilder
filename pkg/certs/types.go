package certs

import (
	"context"
	"crypto/x509"
	"time"
)

// KindCertExtractor defines the interface for extracting certificates from Kind clusters
type KindCertExtractor interface {
	// ExtractGiteaCert extracts the Gitea certificate from the Kind cluster
	ExtractGiteaCert(ctx context.Context) (*x509.Certificate, error)
	
	// GetClusterName returns the name of the configured Kind cluster
	GetClusterName() string
}

// CertValidator defines the interface for validating extracted certificates
// Note: Implementation will be in Split 002
type CertValidator interface {
	// ValidateCertificate validates a certificate against requirements
	ValidateCertificate(cert *x509.Certificate) (*ValidationResult, error)
	
	// CheckExpiry checks if certificate is expired or expiring soon
	CheckExpiry(cert *x509.Certificate, warnDays int) (*ExpiryResult, error)
}

// ExtractorConfig holds configuration for the certificate extractor
type ExtractorConfig struct {
	// ClusterName is the name of the Kind cluster to extract certificates from
	ClusterName string
	
	// Namespace is the Kubernetes namespace where Gitea is deployed
	Namespace string
	
	// PodSelector is the label selector for finding Gitea pods
	PodSelector string
	
	// CertPath is the path inside the pod where certificates are stored
	CertPath string
	
	// OutputDir is the local directory where certificates will be saved
	OutputDir string
	
	// Timeout for Kubernetes operations
	Timeout time.Duration
}

// DefaultExtractorConfig returns a configuration with sensible defaults
func DefaultExtractorConfig() *ExtractorConfig {
	return &ExtractorConfig{
		ClusterName: "localdev",
		Namespace:   "gitea",
		PodSelector: "app=gitea",
		CertPath:    "/data/git/tls/cert.pem",
		OutputDir:   "~/.idpbuilder/certs",
		Timeout:     30 * time.Second,
	}
}

// CertDiagnostics contains diagnostic information about certificate extraction
type CertDiagnostics struct {
	// ClusterConnected indicates if we successfully connected to the cluster
	ClusterConnected bool
	
	// PodsFound is the number of matching pods found
	PodsFound int
	
	// CertificateFound indicates if a certificate was found in the pod
	CertificateFound bool
	
	// CertificateParsed indicates if the certificate was successfully parsed
	CertificateParsed bool
	
	// ExtractionDuration is how long the extraction took
	ExtractionDuration time.Duration
	
	// Warnings contains any warnings encountered during extraction
	Warnings []string
}

// ValidationResult contains the result of certificate validation
// Note: Full implementation will be in Split 002
type ValidationResult struct {
	// Valid indicates if the certificate passed all validation checks
	Valid bool
	
	// Issues contains any validation issues found
	Issues []string
	
	// Warnings contains validation warnings
	Warnings []string
}

// ExpiryResult contains information about certificate expiry
// Note: Full implementation will be in Split 002
type ExpiryResult struct {
	// Expired indicates if the certificate is already expired
	Expired bool
	
	// ExpiringSoon indicates if the certificate will expire within the warning period
	ExpiringSoon bool
	
	// DaysUntilExpiry is the number of days until the certificate expires
	DaysUntilExpiry int
	
	// ExpiryDate is the certificate expiry date
	ExpiryDate time.Time
}