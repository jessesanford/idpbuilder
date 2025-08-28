package certs

import (
	"context"
	"crypto/x509"
	"io/fs"
	"time"
)

// ========== Certificate Extraction Types (from cert-extraction) ==========

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

// ========== Trust Store Types (from trust-store) ==========

// TrustStoreLocation defines where trust store certificates are stored
type TrustStoreLocation int

const (
	UserTrustStore TrustStoreLocation = iota
	SystemTrustStore
)

// String returns a string representation of the trust store location
func (l TrustStoreLocation) String() string {
	switch l {
	case UserTrustStore:
		return "user"
	case SystemTrustStore:
		return "system"
	default:
		return "unknown"
	}
}

// CertificateInfo represents metadata about a certificate
type CertificateInfo struct {
	// Subject contains the certificate subject information
	Subject string
	// Issuer contains the certificate issuer information
	Issuer string
	// SerialNumber is the certificate serial number as a string
	SerialNumber string
	// NotBefore is the certificate validity start time
	NotBefore string
	// NotAfter is the certificate validity end time
	NotAfter string
	// Fingerprint is the SHA256 fingerprint of the certificate
	Fingerprint string
}

// Certificate represents a trust store certificate
type Certificate struct {
	// Data contains the raw PEM-encoded certificate data
	Data []byte
	// Info contains parsed certificate metadata
	Info CertificateInfo
	// FilePath is the path where the certificate is stored
	FilePath string
}

// RegistryConfig represents configuration for a container registry
type RegistryConfig struct {
	// Registry is the registry hostname and optional port (e.g., "registry.example.com:5000")
	Registry string
	// Insecure indicates whether to skip TLS verification
	Insecure bool
	// Certificates contains trusted certificates for this registry
	Certificates []Certificate
}

// TrustStoreConfig represents the overall trust store configuration
type TrustStoreConfig struct {
	// Location specifies whether to use user or system trust store
	Location TrustStoreLocation
	// BaseDir is the base directory for certificate storage
	BaseDir string
	// DirPermissions are the permissions for certificate directories
	DirPermissions fs.FileMode
	// FilePermissions are the permissions for certificate files
	FilePermissions fs.FileMode
}