package certs

import (
	"context"
	"crypto/x509"
	"io/fs"
	"time"
)

// ================================
// === CERT-EXTRACTION TYPES ===
// ================================

// KindCertExtractor defines the interface for extracting certificates from Kind clusters
type KindCertExtractor interface {
	// ExtractGiteaCert extracts the Gitea certificate from the Kind cluster
	ExtractGiteaCert(ctx context.Context) (*x509.Certificate, error)
	
	// GetClusterName returns the name of the Kind cluster
	GetClusterName() string
}

// CertValidator defines the interface for validating certificates
type CertValidator interface {
	// ValidateCertificate performs validation checks on a certificate
	ValidateCertificate(cert *x509.Certificate) (*ValidationResult, error)
	
	// CheckExpiry checks if a certificate is expired or expiring soon
	CheckExpiry(cert *x509.Certificate, warnDays int) (*ExpiryResult, error)
}

// ExtractorConfig holds configuration for the certificate extractor
type ExtractorConfig struct {
	// ClusterName is the name of the Kind cluster
	ClusterName string
	
	// Namespace is the Kubernetes namespace where Gitea is deployed
	Namespace string
	
	// PodSelector is the label selector for finding Gitea pods
	PodSelector string
	
	// CertPath is the path to the certificate inside the pod
	CertPath string
	
	// OutputDir is the directory where extracted certificates will be saved
	OutputDir string
	
	// Timeout is the timeout for extraction operations
	Timeout time.Duration
}

// DefaultExtractorConfig returns a configuration with sensible defaults
func DefaultExtractorConfig() *ExtractorConfig {
	return &ExtractorConfig{
		ClusterName: "idpbuilder",
		Namespace:   "idpbuilder",
		PodSelector: "app.kubernetes.io/name=gitea",
		CertPath:    "/etc/gitea/certs/tls.crt",
		OutputDir:   "~/.idpbuilder/certs",
		Timeout:     30 * time.Second,
	}
}

// CertDiagnostics contains diagnostic information about certificate extraction
type CertDiagnostics struct {
	// ClusterName is the name of the Kind cluster
	ClusterName string
	
	// ClusterConnected indicates if the cluster connection was successful
	ClusterConnected bool
	
	// PodName is the name of the pod where the certificate was extracted
	PodName string
	
	// PodsFound is the number of pods found
	PodsFound int
	
	// CertificateFound indicates if the certificate file was found
	CertificateFound bool
	
	// CertificateParsed indicates if the certificate was successfully parsed
	CertificateParsed bool
	
	// ExtractedAt is when the certificate was extracted
	ExtractedAt time.Time
	
	// ExtractionDuration is how long the extraction took
	ExtractionDuration time.Duration
	
	// ValidationResult contains the validation results
	ValidationResult *ValidationResult
	
	// ExpiryResult contains expiry information
	ExpiryResult *ExpiryResult
	
	// SavedTo is the path where the certificate was saved
	SavedTo string
	
	// Errors contains any errors encountered during extraction
	Errors []string
	
	// Warnings contains any warnings during extraction
	Warnings []string
}

// ValidationResult contains the results of certificate validation
type ValidationResult struct {
	// Valid indicates if the certificate passed all validation checks
	Valid bool
	
	// Issues contains any validation issues found
	Issues []string
	
	// Errors contains any validation errors found
	Errors []string
	
	// Warnings contains any validation warnings
	Warnings []string
	
	// Subject is the certificate subject
	Subject string
	
	// Issuer is the certificate issuer
	Issuer string
	
	// NotBefore is the certificate validity start time
	NotBefore time.Time
	
	// NotAfter is the certificate validity end time
	NotAfter time.Time
	
	// IsCA indicates if this is a CA certificate
	IsCA bool
	
	// IsSelfSigned indicates if the certificate is self-signed
	IsSelfSigned bool
	
	// DNSNames contains the DNS names in the certificate
	DNSNames []string
	
	// IPAddresses contains the IP addresses in the certificate
	IPAddresses []string
}

// ExpiryResult contains certificate expiry information
type ExpiryResult struct {
	// Expired indicates if the certificate has already expired
	Expired bool
	
	// ExpiringSoon indicates if the certificate will expire within the warning period
	ExpiringSoon bool
	
	// DaysUntilExpiry is the number of days until the certificate expires
	DaysUntilExpiry int
	
	// ExpiryDate is the certificate expiry date
	ExpiryDate time.Time
}

// ================================
// === TRUST-STORE TYPES ===
// ================================

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