// Package certificates provides Gitea-specific certificate discovery and management.
// This file implements automatic discovery of Gitea certificates from various sources
// including configuration files, registry endpoints, and CA certificates.
package certificates

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// GiteaDiscovery handles automatic discovery of Gitea certificates
// from various sources including configuration files and API endpoints.
type GiteaDiscovery struct {
	// Configuration
	configPaths []string      // Paths to search for Gitea config
	registryURL string        // Registry URL for certificate discovery
	timeout     time.Duration // HTTP timeout for discovery requests

	// Caching
	discoveryCache map[string]*x509.Certificate // Cache discovered certificates
	lastDiscovery  time.Time                    // Last discovery timestamp
	cacheTTL       time.Duration                // Cache time-to-live

	// Thread safety
	mu sync.RWMutex
}

// GiteaConfig represents Gitea configuration structure
type GiteaConfig struct {
	Server struct {
		CertFile string `json:"CERT_FILE"`
		KeyFile  string `json:"KEY_FILE"`
		Protocol string `json:"PROTOCOL"`
	} `json:"server"`
	Security struct {
		RootPath string `json:"ROOT_PATH"`
	} `json:"security"`
}

// NewGiteaDiscovery creates a new Gitea discovery instance with default configuration
func NewGiteaDiscovery() (*GiteaDiscovery, error) {
	return &GiteaDiscovery{
		configPaths: []string{
			"/etc/gitea/app.ini",
			"/var/lib/gitea/app.ini",
			"./gitea/app.ini",
			"./config/app.ini",
		},
		timeout:        10 * time.Second,
		discoveryCache: make(map[string]*x509.Certificate),
		cacheTTL:       5 * time.Minute,
	}, nil
}

// DiscoverGiteaCertificates discovers certificates associated with a Gitea instance
func (g *GiteaDiscovery) DiscoverGiteaCertificates(ctx context.Context, giteaURL string) ([]*x509.Certificate, error) {
	g.mu.Lock()
	defer g.mu.Unlock()

	var allCerts []*x509.Certificate

	// Check cache first
	if time.Since(g.lastDiscovery) < g.cacheTTL {
		var cachedCerts []*x509.Certificate
		for _, cert := range g.discoveryCache {
			cachedCerts = append(cachedCerts, cert)
		}
		if len(cachedCerts) > 0 {
			return cachedCerts, nil
		}
	}

	// Discover from configuration files
	configCerts, err := g.discoverFromConfig(ctx)
	if err == nil {
		allCerts = append(allCerts, configCerts...)
	}

	// Discover from registry endpoint
	registryCerts, err := g.discoverFromRegistry(ctx, giteaURL)
	if err == nil {
		allCerts = append(allCerts, registryCerts...)
	}

	// Discover root CA certificates
	rootCACerts, err := g.discoverRootCA(ctx, giteaURL)
	if err == nil {
		allCerts = append(allCerts, rootCACerts...)
	}

	// Update cache
	g.discoveryCache = make(map[string]*x509.Certificate)
	for i, cert := range allCerts {
		cacheKey := fmt.Sprintf("cert_%d", i)
		g.discoveryCache[cacheKey] = cert
	}
	g.lastDiscovery = time.Now()

	return allCerts, nil
}

// LoadGiteaRootCA loads the root CA certificate used by Gitea
func (g *GiteaDiscovery) LoadGiteaRootCA(ctx context.Context, rootCAPath string) (*x509.Certificate, error) {
	// Check if file exists
	if _, err := os.Stat(rootCAPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("root CA file not found: %s", rootCAPath)
	}

	// Read CA file
	data, err := ioutil.ReadFile(rootCAPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read root CA file: %w", err)
	}

	// Parse certificate
	certs, err := g.parseCertificateData(data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse root CA certificate: %w", err)
	}

	if len(certs) == 0 {
		return nil, fmt.Errorf("no valid certificates found in root CA file")
	}

	// Validate CA constraints
	rootCA := certs[0]
	if !rootCA.IsCA {
		return nil, fmt.Errorf("certificate is not a CA certificate")
	}

	return rootCA, nil
}

// LoadGiteaRegistryCert loads the certificate used by Gitea's built-in registry
func (g *GiteaDiscovery) LoadGiteaRegistryCert(ctx context.Context, registryURL string) (*x509.Certificate, error) {
	if registryURL == "" {
		return nil, fmt.Errorf("registry URL cannot be empty")
	}

	// Create HTTP client with custom transport to capture certificates
	var serverCert *x509.Certificate
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
			VerifyConnection: func(cs tls.ConnectionState) error {
				if len(cs.PeerCertificates) > 0 {
					serverCert = cs.PeerCertificates[0]
				}
				return nil
			},
		},
	}

	client := &http.Client{
		Transport: tr,
		Timeout:   g.timeout,
	}

	// Make request to capture certificate
	resp, err := client.Get(registryURL + "/v2/")
	if err != nil {
		// Even if request fails, we might have captured the certificate
		if serverCert != nil {
			return serverCert, nil
		}
		return nil, fmt.Errorf("failed to connect to registry: %w", err)
	}
	defer resp.Body.Close()

	if serverCert == nil {
		return nil, fmt.Errorf("no certificate found from registry connection")
	}

	return serverCert, nil
}

// ValidateGiteaCertChain validates a certificate chain against Gitea's root CA
func (g *GiteaDiscovery) ValidateGiteaCertChain(ctx context.Context, cert *x509.Certificate, rootCA *x509.Certificate) error {
	if rootCA == nil {
		return fmt.Errorf("root CA certificate is required for chain validation")
	}

	// Create certificate pool with root CA
	roots := x509.NewCertPool()
	roots.AddCert(rootCA)

	// Create verification options
	opts := x509.VerifyOptions{
		Roots: roots,
	}

	// Perform chain verification
	chains, err := cert.Verify(opts)
	if err != nil {
		return fmt.Errorf("certificate chain validation failed: %w", err)
	}

	if len(chains) == 0 {
		return fmt.Errorf("no valid certificate chains found")
	}

	return nil
}

// ExtractCertFromGiteaConfig extracts certificate path from Gitea configuration
func (g *GiteaDiscovery) ExtractCertFromGiteaConfig(configPath string) (string, error) {
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return "", fmt.Errorf("failed to read config file: %w", err)
	}

	config := string(data)
	lines := strings.Split(config, "\n")

	var certFile string
	inServerSection := false

	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Check for section header
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			section := strings.Trim(line, "[]")
			inServerSection = (section == "server")
			continue
		}

		// Look for CERT_FILE in server section
		if inServerSection && strings.Contains(line, "CERT_FILE") {
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				certFile = strings.TrimSpace(parts[1])
				break
			}
		}
	}

	if certFile == "" {
		return "", fmt.Errorf("CERT_FILE not found in config")
	}

	return certFile, nil
}

// discoverFromConfig discovers certificates from Gitea configuration files
func (g *GiteaDiscovery) discoverFromConfig(ctx context.Context) ([]*x509.Certificate, error) {
	var allCerts []*x509.Certificate

	for _, configPath := range g.configPaths {
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			continue
		}

		certPath, err := g.ExtractCertFromGiteaConfig(configPath)
		if err != nil {
			continue
		}

		// Make path absolute if relative
		if !filepath.IsAbs(certPath) {
			configDir := filepath.Dir(configPath)
			certPath = filepath.Join(configDir, certPath)
		}

		cert, err := g.loadCertificateFromFile(certPath)
		if err != nil {
			continue
		}

		allCerts = append(allCerts, cert)
	}

	return allCerts, nil
}

// discoverFromRegistry discovers certificates from registry endpoint
func (g *GiteaDiscovery) discoverFromRegistry(ctx context.Context, giteaURL string) ([]*x509.Certificate, error) {
	if giteaURL == "" {
		return nil, fmt.Errorf("gitea URL is required")
	}

	registryURL := giteaURL
	if !strings.HasSuffix(registryURL, "/") {
		registryURL += "/"
	}
	registryURL += "v2/"

	cert, err := g.LoadGiteaRegistryCert(ctx, registryURL)
	if err != nil {
		return nil, err
	}

	return []*x509.Certificate{cert}, nil
}

// discoverRootCA discovers root CA certificates
func (g *GiteaDiscovery) discoverRootCA(ctx context.Context, giteaURL string) ([]*x509.Certificate, error) {
	commonRootPaths := []string{
		"/etc/ssl/certs/gitea-root-ca.crt",
		"/var/lib/gitea/certs/root-ca.crt",
		"./certs/root-ca.crt",
		"./gitea/certs/root-ca.crt",
	}

	var allCerts []*x509.Certificate

	for _, rootPath := range commonRootPaths {
		cert, err := g.LoadGiteaRootCA(ctx, rootPath)
		if err != nil {
			continue
		}
		allCerts = append(allCerts, cert)
	}

	return allCerts, nil
}

// loadCertificateFromFile loads a certificate from a file
func (g *GiteaDiscovery) loadCertificateFromFile(filePath string) (*x509.Certificate, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read certificate file: %w", err)
	}

	certs, err := g.parseCertificateData(data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse certificate: %w", err)
	}

	if len(certs) == 0 {
		return nil, fmt.Errorf("no valid certificates found in file")
	}

	return certs[0], nil
}

// parseCertificateData parses certificate data in various formats
func (g *GiteaDiscovery) parseCertificateData(data []byte) ([]*x509.Certificate, error) {
	var certs []*x509.Certificate

	// Try parsing as PEM first
	rest := data
	for len(rest) > 0 {
		block, remaining := pem.Decode(rest)
		if block == nil {
			break
		}

		if block.Type == "CERTIFICATE" {
			cert, err := x509.ParseCertificate(block.Bytes)
			if err == nil {
				certs = append(certs, cert)
			}
		}

		rest = remaining
	}

	// If no PEM certificates found, try DER format
	if len(certs) == 0 {
		cert, err := x509.ParseCertificate(data)
		if err == nil {
			certs = append(certs, cert)
		}
	}

	return certs, nil
}
