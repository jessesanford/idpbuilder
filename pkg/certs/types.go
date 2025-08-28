package certs

import (
	"io/fs"
)

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