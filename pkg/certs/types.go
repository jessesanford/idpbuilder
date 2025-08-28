package certs

import (
	"context"
	"crypto/x509"
	"time"
)

// KindCertExtractor defines the interface for extracting certificates from Kind clusters
type KindCertExtractor interface {
	// ExtractGiteaCert extracts the Gitea certificate from a Kind cluster
	ExtractGiteaCert(ctx context.Context) (*x509.Certificate, error)
	
	// GetClusterName returns the name of the Kind cluster
	GetClusterName() (string, error)
	
	// ValidateCertificate validates that a certificate is suitable for use
	ValidateCertificate(cert *x509.Certificate) error
}

// CertValidator defines certificate validation operations
type CertValidator interface {
	// ValidateChain validates the certificate chain
	ValidateChain(cert *x509.Certificate) error
	
	// CheckExpiry checks if certificate is expired or will expire soon
	CheckExpiry(cert *x509.Certificate) (*time.Duration, error)
	
	// VerifyHostname verifies the certificate is valid for the given hostname
	VerifyHostname(cert *x509.Certificate, hostname string) error
	
	// GenerateDiagnostics provides detailed certificate diagnostic information
	GenerateDiagnostics() (*CertDiagnostics, error)
}

// CertDiagnostics contains certificate diagnostic information
type CertDiagnostics struct {
	IsValid       bool      `json:"is_valid"`
	ExpiresAt     time.Time `json:"expires_at"`
	DaysUntilExpiry int     `json:"days_until_expiry"`
	Subject       string    `json:"subject"`
	Issuer        string    `json:"issuer"`
	DNSNames      []string  `json:"dns_names"`
	Issues        []string  `json:"issues"`
	Recommendations []string `json:"recommendations"`
}

// ExtractorConfig contains configuration for certificate extraction
type ExtractorConfig struct {
	ClusterName     string `json:"cluster_name"`
	GiteaNamespace  string `json:"gitea_namespace"`
	GiteaPodLabelSelector string `json:"gitea_pod_label_selector"`
	CertPath        string `json:"cert_path"`
	StoragePath     string `json:"storage_path"`
	Timeout         time.Duration `json:"timeout"`
}

// DefaultExtractorConfig returns a default configuration for certificate extraction
func DefaultExtractorConfig() *ExtractorConfig {
	return &ExtractorConfig{
		ClusterName:           "idpbuilder",
		GiteaNamespace:        "gitea",
		GiteaPodLabelSelector: "app.kubernetes.io/name=gitea",
		CertPath:              "/data/gitea/https/cert.pem",
		StoragePath:           "~/.idpbuilder/certs/gitea.pem",
		Timeout:               30 * time.Second,
	}
}