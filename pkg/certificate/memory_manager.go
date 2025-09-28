package certificate

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"io"
	"math/big"
	"sync"
	"time"
)

// MemoryManager is an in-memory implementation of CertificateManager.
// It stores certificates in memory and provides all certificate management operations.
// This implementation is suitable for development and testing scenarios.
type MemoryManager struct {
	// certificates stores the certificates by key
	certificates map[string]*Certificate
	// mutex protects concurrent access to the certificates map
	mutex sync.RWMutex
}

// NewMemoryManager creates a new instance of MemoryManager.
func NewMemoryManager() *MemoryManager {
	return &MemoryManager{
		certificates: make(map[string]*Certificate),
	}
}

// GenerateSelfSigned creates a new self-signed certificate with the specified options.
func (m *MemoryManager) GenerateSelfSigned(opts *GenerationOptions) (*Certificate, error) {
	// Generate private key
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("generating private key: %w", err)
	}

	// Generate certificate serial number
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return nil, fmt.Errorf("generating certificate serial number: %w", err)
	}

	// Set certificate validity period
	notBefore := time.Now()
	notAfter := notBefore.Add(opts.ValidFor)

	// Create certificate template
	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			CommonName:   opts.Subject,
			Organization: []string{opts.Organization},
		},
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		KeyUsage:              opts.KeyUsage,
		ExtKeyUsage:           opts.ExtKeyUsage,
		BasicConstraintsValid: true,
		IsCA:                  opts.IsCA,
		DNSNames:              opts.DNSNames,
	}

	// Generate the certificate
	certBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	if err != nil {
		return nil, fmt.Errorf("creating certificate: %w", err)
	}

	// Encode certificate to PEM
	certPEM, err := encodeCertificateToPEM(certBytes)
	if err != nil {
		return nil, fmt.Errorf("encoding certificate to PEM: %w", err)
	}

	// Encode private key to PEM
	keyPEM, err := encodePrivateKeyToPEM(privateKey)
	if err != nil {
		return nil, fmt.Errorf("encoding private key to PEM: %w", err)
	}

	// Create metadata
	metadata := &CertificateMetadata{
		Subject:     opts.Subject,
		DNSNames:    opts.DNSNames,
		NotBefore:   notBefore,
		NotAfter:    notAfter,
		IsCA:        opts.IsCA,
		KeyUsage:    opts.KeyUsage,
		ExtKeyUsage: opts.ExtKeyUsage,
	}

	certificate := &Certificate{
		CertPEM:  certPEM,
		KeyPEM:   keyPEM,
		Metadata: metadata,
	}

	return certificate, nil
}

// Store saves a certificate for later retrieval.
func (m *MemoryManager) Store(key string, cert *Certificate) error {
	if key == "" {
		return fmt.Errorf("certificate key cannot be empty")
	}
	if cert == nil {
		return fmt.Errorf("certificate cannot be nil")
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.certificates[key] = cert
	return nil
}

// Retrieve gets a previously stored certificate by its key.
func (m *MemoryManager) Retrieve(key string) (*Certificate, error) {
	if key == "" {
		return nil, fmt.Errorf("certificate key cannot be empty")
	}

	m.mutex.RLock()
	defer m.mutex.RUnlock()

	cert, exists := m.certificates[key]
	if !exists {
		return nil, fmt.Errorf("certificate with key '%s' not found", key)
	}

	return cert, nil
}

// IsValid checks if a certificate is currently valid.
func (m *MemoryManager) IsValid(cert *Certificate) bool {
	if cert == nil || cert.Metadata == nil {
		return false
	}

	now := time.Now()
	return now.After(cert.Metadata.NotBefore) && now.Before(cert.Metadata.NotAfter)
}

// GetExpiration returns the expiration time of a certificate.
func (m *MemoryManager) GetExpiration(cert *Certificate) (time.Time, error) {
	if cert == nil || cert.Metadata == nil {
		return time.Time{}, fmt.Errorf("certificate or metadata is nil")
	}

	return cert.Metadata.NotAfter, nil
}

// List returns all stored certificate keys.
func (m *MemoryManager) List() ([]string, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	keys := make([]string, 0, len(m.certificates))
	for key := range m.certificates {
		keys = append(keys, key)
	}

	return keys, nil
}

// Delete removes a stored certificate.
func (m *MemoryManager) Delete(key string) error {
	if key == "" {
		return fmt.Errorf("certificate key cannot be empty")
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	if _, exists := m.certificates[key]; !exists {
		return fmt.Errorf("certificate with key '%s' not found", key)
	}

	delete(m.certificates, key)
	return nil
}

// Helper function to encode certificate bytes to PEM format
func encodeCertificateToPEM(certBytes []byte) ([]byte, error) {
	var certBuffer bytes.Buffer
	err := pem.Encode(io.Writer(&certBuffer), &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	})
	if err != nil {
		return nil, err
	}

	return io.ReadAll(&certBuffer)
}

// Helper function to encode private key to PEM format
func encodePrivateKeyToPEM(privateKey *ecdsa.PrivateKey) ([]byte, error) {
	privateKeyBytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		return nil, fmt.Errorf("marshaling private key: %w", err)
	}

	var keyBuffer bytes.Buffer
	err = pem.Encode(io.Writer(&keyBuffer), &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: privateKeyBytes,
	})
	if err != nil {
		return nil, err
	}

	return io.ReadAll(&keyBuffer)
}