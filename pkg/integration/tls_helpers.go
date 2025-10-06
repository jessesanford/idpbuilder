//go:build integration

package integration

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"net/http"
	"time"
)

// IntegrationTestConfig holds configuration for integration tests
type IntegrationTestConfig struct {
	RegistryURL    string        // Registry endpoint
	InsecureMode   bool          // Allow insecure connections
	TestImagePath  string        // Path to test image
	Timeout        time.Duration // Test timeout
}

// SetupInsecureCertTest creates a TLS configuration for testing with self-signed certificates
func SetupInsecureCertTest() (*tls.Config, error) {
	// Generate a self-signed certificate for testing
	cert, key, err := generateSelfSignedCert()
	if err != nil {
		return nil, fmt.Errorf("failed to generate self-signed certificate: %w", err)
	}

	// Create TLS certificate from generated cert and key
	tlsCert, err := tls.X509KeyPair(cert, key)
	if err != nil {
		return nil, fmt.Errorf("failed to create TLS certificate: %w", err)
	}

	// Configure TLS with the self-signed certificate
	tlsConfig := &tls.Config{
		Certificates:       []tls.Certificate{tlsCert},
		InsecureSkipVerify: true, // For testing purposes
		MinVersion:         tls.VersionTLS12,
	}

	return tlsConfig, nil
}

// generateSelfSignedCert creates a self-signed certificate for testing
func generateSelfSignedCert() ([]byte, []byte, error) {
	// Generate RSA private key
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate private key: %w", err)
	}

	// Create certificate template
	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization:  []string{"Test Organization"},
			Country:       []string{"US"},
			Province:      []string{""},
			Locality:      []string{"Test City"},
			StreetAddress: []string{""},
			PostalCode:    []string{""},
		},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().Add(365 * 24 * time.Hour), // Valid for 1 year
		KeyUsage:     x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		IPAddresses:  []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},
		DNSNames:     []string{"localhost"},
	}

	// Create certificate
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create certificate: %w", err)
	}

	// Encode certificate to PEM
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})

	// Encode private key to PEM
	privateKeyDER, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal private key: %w", err)
	}
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: privateKeyDER})

	return certPEM, keyPEM, nil
}

// validateTLSConfig checks if a TLS configuration is suitable for testing
func validateTLSConfig(config *tls.Config) error {
	if config == nil {
		return fmt.Errorf("TLS config cannot be nil")
	}
	if config.MinVersion < tls.VersionTLS12 {
		return fmt.Errorf("minimum TLS version should be 1.2 or higher")
	}
	return nil
}

// createTestConfig creates a default integration test configuration
func createTestConfig(registryURL string) *IntegrationTestConfig {
	return &IntegrationTestConfig{
		RegistryURL:    registryURL,
		InsecureMode:   true,
		TestImagePath:  "test/image:latest",
		Timeout:        30 * time.Second,
	}
}

// configureInsecureTransport creates an HTTP transport with insecure TLS
func configureInsecureTransport() *http.Transport {
	return &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
}