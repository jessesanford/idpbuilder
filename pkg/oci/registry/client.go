package registry

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/cnoe-io/idpbuilder/pkg/oci/api"
	"github.com/cnoe-io/idpbuilder/pkg/oci/security"
)

// registryClient implements the RegistryClient interface with TLS support
// for self-signed certificates, particularly for gitea.cnoe.localtest.me
type registryClient struct {
	httpClient    *http.Client
	securityMgr   api.SecurityManager
	defaultAuth   api.AuthConfig
	transportOpts TransportOptions
}

// TransportOptions configures HTTP transport behavior including TLS
type TransportOptions struct {
	MaxRetries         int
	RetryBackoff       time.Duration
	InsecureSkipVerify bool   // Critical for gitea.cnoe.localtest.me
	CACertPath         string // Custom CA certificate path
	ClientCertPath     string // Client certificate for mutual TLS
	ClientKeyPath      string // Client private key for mutual TLS
	Timeout            time.Duration
}

// Option configures the registry client
type Option func(*registryClient)

// WithSecurityManager sets the security manager for image signing/verification
func WithSecurityManager(sm api.SecurityManager) Option {
	return func(rc *registryClient) {
		rc.securityMgr = sm
	}
}

// WithDefaultAuth sets default authentication configuration
func WithDefaultAuth(auth api.AuthConfig) Option {
	return func(rc *registryClient) {
		rc.defaultAuth = auth
	}
}

// WithTransportOptions configures HTTP transport behavior
func WithTransportOptions(opts TransportOptions) Option {
	return func(rc *registryClient) {
		rc.transportOpts = opts
	}
}

// WithInsecureSkipVerify enables TLS verification skip for self-signed certificates
// This is specifically designed for gitea.cnoe.localtest.me development usage
func WithInsecureSkipVerify(insecure bool) Option {
	return func(rc *registryClient) {
		rc.transportOpts.InsecureSkipVerify = insecure
	}
}

// WithCACertificate sets custom CA certificate for registry verification
func WithCACertificate(caCertPath string) Option {
	return func(rc *registryClient) {
		rc.transportOpts.CACertPath = caCertPath
	}
}

// NewRegistryClient creates a new registry client with TLS support
// Configured specifically to handle gitea.cnoe.localtest.me self-signed certificates
func NewRegistryClient(opts ...Option) api.RegistryClient {
	rc := &registryClient{
		transportOpts: TransportOptions{
			MaxRetries:   3,
			RetryBackoff: 2 * time.Second,
			Timeout:      30 * time.Second,
		},
	}

	// Apply all options
	for _, opt := range opts {
		opt(rc)
	}

	// Configure HTTP client with TLS support
	rc.httpClient = rc.createHTTPClient()

	return rc
}

// createHTTPClient configures HTTP client with TLS settings for self-signed certs
func (rc *registryClient) createHTTPClient() *http.Client {
	tlsConfig := &tls.Config{
		InsecureSkipVerify: rc.transportOpts.InsecureSkipVerify,
	}

	// Load custom CA certificate if provided
	if rc.transportOpts.CACertPath != "" {
		caCert, err := ioutil.ReadFile(rc.transportOpts.CACertPath)
		if err == nil {
			caCertPool := x509.NewCertPool()
			if caCertPool.AppendCertsFromPEM(caCert) {
				tlsConfig.RootCAs = caCertPool
				// Disable InsecureSkipVerify when using custom CA
				tlsConfig.InsecureSkipVerify = false
			}
		}
	}

	// Load client certificates for mutual TLS if provided
	if rc.transportOpts.ClientCertPath != "" && rc.transportOpts.ClientKeyPath != "" {
		cert, err := tls.LoadX509KeyPair(
			rc.transportOpts.ClientCertPath,
			rc.transportOpts.ClientKeyPath,
		)
		if err == nil {
			tlsConfig.Certificates = []tls.Certificate{cert}
		}
	}

	// Create base transport
	baseTransport := &http.Transport{
		TLSClientConfig:       tlsConfig,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	// Wrap with retry transport for resilience
	retryTransport := newRetryTransport(
		baseTransport,
		rc.transportOpts.MaxRetries,
		rc.transportOpts.RetryBackoff,
	)

	return &http.Client{
		Transport: retryTransport,
		Timeout:   rc.transportOpts.Timeout,
	}
}

// isGiteaRegistry checks if the registry is gitea.cnoe.localtest.me
// This enables special handling for local development registries
func (rc *registryClient) isGiteaRegistry(serverAddress string) bool {
	return strings.Contains(serverAddress, "gitea.cnoe.localtest.me")
}

// getRegistryURL constructs the base registry URL from image reference
func (rc *registryClient) getRegistryURL(ref *imageReference) string {
	if ref.Registry == "" {
		return "https://registry-1.docker.io" // Default to Docker Hub
	}
	
	// For gitea.cnoe.localtest.me, ensure HTTPS but allow self-signed
	if rc.isGiteaRegistry(ref.Registry) {
		return fmt.Sprintf("https://%s", ref.Registry)
	}
	
	// For other registries, use HTTPS by default
	if !strings.HasPrefix(ref.Registry, "http://") && !strings.HasPrefix(ref.Registry, "https://") {
		return fmt.Sprintf("https://%s", ref.Registry)
	}
	
	return ref.Registry
}

// getRegistryBaseURL returns the base URL without version path
func (rc *registryClient) getRegistryBaseURL() string {
	if rc.defaultAuth.ServerAddress == "" {
		return "https://registry-1.docker.io"
	}
	return rc.getRegistryURL(&imageReference{Registry: rc.defaultAuth.ServerAddress})
}

// Close closes the registry client and cleans up resources
func (rc *registryClient) Close() error {
	if rc.httpClient != nil {
		rc.httpClient.CloseIdleConnections()
	}
	return nil
}

// imageReference represents a parsed OCI image reference
type imageReference struct {
	Registry   string
	Namespace  string
	Repository string
	Name       string
	Tag        string
	Digest     string
}

// parseImageReference parses an image reference into components
func parseImageReference(image string) (*imageReference, error) {
	ref := &imageReference{}
	
	// Handle digest references (image@sha256:...)
	if strings.Contains(image, "@") {
		parts := strings.SplitN(image, "@", 2)
		image = parts[0]
		ref.Digest = parts[1]
	}
	
	// Handle tag references (image:tag)
	if strings.Contains(image, ":") && ref.Digest == "" {
		parts := strings.SplitN(image, ":", 2)
		image = parts[0]
		ref.Tag = parts[1]
	} else if ref.Digest == "" {
		ref.Tag = "latest"
	}
	
	// Parse registry/namespace/repository
	parts := strings.Split(image, "/")
	
	if len(parts) == 1 {
		// Simple name like "ubuntu"
		ref.Repository = parts[0]
		ref.Name = parts[0]
		ref.Registry = "registry-1.docker.io"
		ref.Namespace = "library"
	} else if len(parts) == 2 {
		// Could be registry/repo or namespace/repo
		if strings.Contains(parts[0], ".") || strings.Contains(parts[0], ":") {
			// Registry with port or domain
			ref.Registry = parts[0]
			ref.Repository = parts[1]
			ref.Name = parts[1]
		} else {
			// Namespace/repo on Docker Hub
			ref.Registry = "registry-1.docker.io"
			ref.Namespace = parts[0]
			ref.Repository = parts[1]
			ref.Name = parts[1]
		}
	} else if len(parts) >= 3 {
		// Full registry/namespace/repo
		ref.Registry = parts[0]
		ref.Namespace = parts[1]
		ref.Repository = strings.Join(parts[2:], "/")
		ref.Name = strings.Join(parts[1:], "/")
	}
	
	return ref, nil
}

// String returns the string representation of the image reference
func (ref *imageReference) String() string {
	var result strings.Builder
	
	if ref.Registry != "" && ref.Registry != "registry-1.docker.io" {
		result.WriteString(ref.Registry)
		result.WriteString("/")
	}
	
	if ref.Namespace != "" && ref.Namespace != "library" {
		result.WriteString(ref.Namespace)
		result.WriteString("/")
	}
	
	result.WriteString(ref.Repository)
	
	if ref.Tag != "" && ref.Digest == "" {
		result.WriteString(":")
		result.WriteString(ref.Tag)
	}
	
	if ref.Digest != "" {
		result.WriteString("@")
		result.WriteString(ref.Digest)
	}
	
	return result.String()
}