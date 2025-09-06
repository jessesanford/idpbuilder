// pkg/certs/extractor_test.go
package certs

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// createTestCertificate creates a test certificate for testing purposes
func createTestCertificate(dnsNames []string, expiry time.Time) (*x509.Certificate, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}

	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization:  []string{"Test Org"},
			Country:       []string{"US"},
			Province:      []string{""},
			Locality:      []string{"Test City"},
			StreetAddress: []string{""},
			PostalCode:    []string{""},
		},
		NotBefore:             time.Now(),
		NotAfter:              expiry,
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		DNSNames:              dnsNames,
	}

	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	if err != nil {
		return nil, err
	}

	cert, err := x509.ParseCertificate(certDER)
	if err != nil {
		return nil, err
	}

	return cert, nil
}

// createTestCertificatePEM creates a test certificate in PEM format
func createTestCertificatePEM(dnsNames []string, expiry time.Time) ([]byte, error) {
	cert, err := createTestCertificate(dnsNames, expiry)
	if err != nil {
		return nil, err
	}

	certPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: cert.Raw,
	})

	return certPEM, nil
}

func TestNewDefaultExtractor(t *testing.T) {
	clusterName := "test-cluster"
	extractor := NewDefaultExtractor(clusterName)

	if extractor.clusterName != clusterName {
		t.Errorf("Expected cluster name %s, got %s", clusterName, extractor.clusterName)
	}

	if extractor.namespace != "gitea" {
		t.Errorf("Expected namespace 'gitea', got %s", extractor.namespace)
	}
}

func TestValidateCertificate_ValidCert(t *testing.T) {
	// Create a valid certificate
	dnsNames := []string{"gitea.example.com"}
	expiry := time.Now().Add(24 * time.Hour)
	cert, err := createTestCertificate(dnsNames, expiry)
	if err != nil {
		t.Fatalf("Failed to create test certificate: %v", err)
	}

	extractor := NewDefaultExtractor("test")
	err = extractor.ValidateCertificate(cert)
	if err != nil {
		t.Errorf("Expected valid certificate, got error: %v", err)
	}
}

func TestValidateCertificate_NilCert(t *testing.T) {
	extractor := NewDefaultExtractor("test")
	err := extractor.ValidateCertificate(nil)

	if err == nil {
		t.Error("Expected error for nil certificate")
	}

	if _, ok := err.(CertificateInvalidError); !ok {
		t.Errorf("Expected CertificateInvalidError, got %T", err)
	}
}

func TestValidateCertificate_ExpiredCert(t *testing.T) {
	// Create an expired certificate
	dnsNames := []string{"gitea.example.com"}
	expiry := time.Now().Add(-24 * time.Hour) // Expired 24 hours ago
	cert, err := createTestCertificate(dnsNames, expiry)
	if err != nil {
		t.Fatalf("Failed to create test certificate: %v", err)
	}

	extractor := NewDefaultExtractor("test")
	err = extractor.ValidateCertificate(cert)

	if err == nil {
		t.Error("Expected error for expired certificate")
	}

	if _, ok := err.(CertificateInvalidError); !ok {
		t.Errorf("Expected CertificateInvalidError, got %T", err)
	}
}

func TestValidateCertificate_NotYetValidCert(t *testing.T) {
	// Create a certificate that's not yet valid
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("Failed to generate key: %v", err)
	}

	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization: []string{"Test Org"},
			CommonName:   "gitea.example.com",
		},
		NotBefore:             time.Now().Add(24 * time.Hour), // Valid 24 hours from now
		NotAfter:              time.Now().Add(48 * time.Hour),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		DNSNames:              []string{"gitea.example.com"},
	}

	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	if err != nil {
		t.Fatalf("Failed to create certificate: %v", err)
	}

	cert, err := x509.ParseCertificate(certDER)
	if err != nil {
		t.Fatalf("Failed to parse certificate: %v", err)
	}

	extractor := NewDefaultExtractor("test")
	err = extractor.ValidateCertificate(cert)

	if err == nil {
		t.Error("Expected error for not-yet-valid certificate")
	}

	if _, ok := err.(CertificateInvalidError); !ok {
		t.Errorf("Expected CertificateInvalidError, got %T", err)
	}
}

func TestValidateCertificate_NoGiteaIdentity(t *testing.T) {
	// Create a certificate without Gitea identity
	dnsNames := []string{"example.com", "test.com"}
	expiry := time.Now().Add(24 * time.Hour)

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("Failed to generate key: %v", err)
	}

	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization: []string{"Test Org"},
			CommonName:   "example.com",
		},
		NotBefore:             time.Now(),
		NotAfter:              expiry,
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		DNSNames:              dnsNames,
	}

	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	if err != nil {
		t.Fatalf("Failed to create certificate: %v", err)
	}

	cert, err := x509.ParseCertificate(certDER)
	if err != nil {
		t.Fatalf("Failed to parse certificate: %v", err)
	}

	extractor := NewDefaultExtractor("test")
	err = extractor.ValidateCertificate(cert)

	if err == nil {
		t.Error("Expected error for certificate without Gitea identity")
	}

	if _, ok := err.(CertificateInvalidError); !ok {
		t.Errorf("Expected CertificateInvalidError, got %T", err)
	}
}

func TestSaveCertificate_Success(t *testing.T) {
	// Create a valid certificate
	dnsNames := []string{"gitea.example.com"}
	expiry := time.Now().Add(24 * time.Hour)
	cert, err := createTestCertificate(dnsNames, expiry)
	if err != nil {
		t.Fatalf("Failed to create test certificate: %v", err)
	}

	// Create temp directory for test
	tempDir, err := os.MkdirTemp("", "cert_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	certPath := filepath.Join(tempDir, "test.pem")
	extractor := NewDefaultExtractor("test")

	err = extractor.SaveCertificate(cert, certPath)
	if err != nil {
		t.Errorf("Failed to save certificate: %v", err)
	}

	// Verify file exists and has correct permissions
	info, err := os.Stat(certPath)
	if err != nil {
		t.Errorf("Certificate file not created: %v", err)
	}

	if info.Mode().Perm() != 0600 {
		t.Errorf("Expected file permissions 0600, got %o", info.Mode().Perm())
	}

	// Verify file content
	data, err := os.ReadFile(certPath)
	if err != nil {
		t.Errorf("Failed to read certificate file: %v", err)
	}

	block, _ := pem.Decode(data)
	if block == nil {
		t.Error("Certificate file does not contain valid PEM data")
	}
}

func TestSaveCertificate_NilCert(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "cert_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	certPath := filepath.Join(tempDir, "test.pem")
	extractor := NewDefaultExtractor("test")

	err = extractor.SaveCertificate(nil, certPath)

	if err == nil {
		t.Error("Expected error for nil certificate")
	}

	if _, ok := err.(CertificateInvalidError); !ok {
		t.Errorf("Expected CertificateInvalidError, got %T", err)
	}
}

func TestParsePEMCertificate_Success(t *testing.T) {
	// Create test certificate PEM data
	dnsNames := []string{"gitea.example.com"}
	expiry := time.Now().Add(24 * time.Hour)
	certPEM, err := createTestCertificatePEM(dnsNames, expiry)
	if err != nil {
		t.Fatalf("Failed to create test certificate PEM: %v", err)
	}

	extractor := NewDefaultExtractor("test")
	cert, err := extractor.parsePEMCertificate(certPEM)

	if err != nil {
		t.Errorf("Failed to parse PEM certificate: %v", err)
	}

	if cert == nil {
		t.Error("Expected non-nil certificate")
	}

	if len(cert.DNSNames) != 1 || cert.DNSNames[0] != "gitea.example.com" {
		t.Errorf("Expected DNS name 'gitea.example.com', got %v", cert.DNSNames)
	}
}

func TestParsePEMCertificate_InvalidPEM(t *testing.T) {
	invalidPEM := []byte("invalid pem data")

	extractor := NewDefaultExtractor("test")
	_, err := extractor.parsePEMCertificate(invalidPEM)

	if err == nil {
		t.Error("Expected error for invalid PEM data")
	}
}

func TestParsePEMCertificate_WrongBlockType(t *testing.T) {
	// Create PEM with wrong block type
	wrongTypePEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: []byte("dummy data"),
	})

	extractor := NewDefaultExtractor("test")
	_, err := extractor.parsePEMCertificate(wrongTypePEM)

	if err == nil {
		t.Error("Expected error for wrong PEM block type")
	}
}

func TestGetCertificateInfo(t *testing.T) {
	// Create a test certificate with specific attributes
	dnsNames := []string{"gitea.example.com", "gitea.local"}
	expiry := time.Now().Add(24 * time.Hour)
	cert, err := createTestCertificate(dnsNames, expiry)
	if err != nil {
		t.Fatalf("Failed to create test certificate: %v", err)
	}

	info := GetCertificateInfo(cert)

	if info.Subject == "" {
		t.Error("Expected non-empty subject")
	}

	if info.Issuer == "" {
		t.Error("Expected non-empty issuer")
	}

	if len(info.DNSNames) != 2 {
		t.Errorf("Expected 2 DNS names, got %d", len(info.DNSNames))
	}

	if info.IsCA != false {
		t.Errorf("Expected IsCA to be false, got %v", info.IsCA)
	}
}

func TestGetCertificateInfo_NilCert(t *testing.T) {
	info := GetCertificateInfo(nil)

	// Should return zero value without panicking
	if info.Subject != "" || info.Issuer != "" || len(info.DNSNames) != 0 {
		t.Error("Expected zero value for nil certificate")
	}
}

// Test error types to improve coverage
func TestClusterNotFoundError(t *testing.T) {
	err := ClusterNotFoundError{ClusterName: "test-cluster"}
	expected := "Kind cluster not found: test-cluster"
	if err.Error() != expected {
		t.Errorf("Expected %s, got %s", expected, err.Error())
	}
}

func TestPodNotFoundError(t *testing.T) {
	err := PodNotFoundError{PodName: "test-pod", Namespace: "test-namespace"}
	expected := "pod 'test-pod' not found in namespace 'test-namespace'"
	if err.Error() != expected {
		t.Errorf("Expected %s, got %s", expected, err.Error())
	}
}

func TestCertificateInvalidError(t *testing.T) {
	err := CertificateInvalidError{Reason: "expired certificate"}
	expected := "certificate invalid: expired certificate"
	if err.Error() != expected {
		t.Errorf("Expected %s, got %s", expected, err.Error())
	}
}

func TestPermissionError(t *testing.T) {
	err := PermissionError{Path: "/tmp/test", Action: "write"}
	expected := "permission denied: cannot write file at path /tmp/test"
	if err.Error() != expected {
		t.Errorf("Expected %s, got %s", expected, err.Error())
	}
}

func TestSaveCertificate_DirectoryCreation(t *testing.T) {
	// Create a valid certificate
	dnsNames := []string{"gitea.example.com"}
	expiry := time.Now().Add(24 * time.Hour)
	cert, err := createTestCertificate(dnsNames, expiry)
	if err != nil {
		t.Fatalf("Failed to create test certificate: %v", err)
	}

	// Use a nested path that doesn't exist
	tempDir, err := os.MkdirTemp("", "cert_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	certPath := filepath.Join(tempDir, "nested", "path", "test.pem")
	extractor := NewDefaultExtractor("test")

	err = extractor.SaveCertificate(cert, certPath)
	if err != nil {
		t.Errorf("Failed to save certificate with directory creation: %v", err)
	}

	// Verify directories were created
	if _, err := os.Stat(filepath.Dir(certPath)); os.IsNotExist(err) {
		t.Error("Expected nested directories to be created")
	}
}

func TestValidateCertificate_GiteaInCommonName(t *testing.T) {
	// Create a certificate with "gitea" in the CommonName but not in DNS names
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("Failed to generate key: %v", err)
	}

	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization: []string{"Test Org"},
			CommonName:   "gitea.internal.com", // Has gitea in CN
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(24 * time.Hour),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		DNSNames:              []string{"internal.com"}, // No gitea in DNS names
	}

	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	if err != nil {
		t.Fatalf("Failed to create certificate: %v", err)
	}

	cert, err := x509.ParseCertificate(certDER)
	if err != nil {
		t.Fatalf("Failed to parse certificate: %v", err)
	}

	extractor := NewDefaultExtractor("test")
	err = extractor.ValidateCertificate(cert)

	if err != nil {
		t.Errorf("Expected certificate with gitea in CN to be valid, got error: %v", err)
	}
}
