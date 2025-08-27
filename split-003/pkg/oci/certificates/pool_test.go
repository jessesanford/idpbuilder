package certificates

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"sync"
	"testing"
	"time"
)

// mockStore implements CertificateStore for testing
type mockStore struct {
	certificates map[string]*Certificate
	mu           sync.RWMutex
}

func newMockStore() *mockStore {
	return &mockStore{certificates: make(map[string]*Certificate)}
}

func (m *mockStore) AddCertificate(ctx context.Context, cert *Certificate) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.certificates[cert.ID] = cert
	return nil
}

func (m *mockStore) GetCertificate(ctx context.Context, id string) (*Certificate, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if cert, exists := m.certificates[id]; exists {
		return cert, nil
	}
	return nil, NewStorageError(ErrCodeCertNotFound, "certificate not found", nil)
}

func (m *mockStore) UpdateCertificate(ctx context.Context, cert *Certificate) error {
	return m.AddCertificate(ctx, cert)
}

func (m *mockStore) DeleteCertificate(ctx context.Context, id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.certificates, id)
	return nil
}

func (m *mockStore) ListCertificates(ctx context.Context, filter *CertificateFilter) ([]*Certificate, error) {
	var certs []*Certificate
	for _, cert := range m.certificates {
		certs = append(certs, cert)
	}
	return certs, nil
}

func (m *mockStore) SearchCertificates(ctx context.Context, query string) ([]*Certificate, error) {
	return m.ListCertificates(ctx, nil)
}

func (m *mockStore) GetCertificatesByStatus(ctx context.Context, status CertificateStatus) ([]*Certificate, error) {
	return m.ListCertificates(ctx, nil)
}

// generateTestCertificate creates a test certificate
func generateTestCertificate(id, subject string) (*Certificate, error) {
	privateKey, _ := rsa.GenerateKey(rand.Reader, 2048)
	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject:      pkix.Name{CommonName: subject},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().Add(365 * 24 * time.Hour),
		KeyUsage:     x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
	}

	certDER, _ := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	cert, _ := x509.ParseCertificate(certDER)
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})

	return &Certificate{
		ID:          id,
		Name:        subject,
		Certificate: cert,
		PEM:         certPEM,
		Status:      CertificateStatusActive,
		ValidFrom:   cert.NotBefore,
		ValidTo:     cert.NotAfter,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Tags:        make(map[string]string),
	}, nil
}

func TestPoolManager_CreatePool(t *testing.T) {
	pm, _ := NewPoolManager(newMockStore())
	ctx := context.Background()
	
	err := pm.CreatePool(ctx, "test-pool")
	if err != nil {
		t.Errorf("Failed to create pool: %v", err)
	}

	err = pm.CreatePool(ctx, "test-pool")
	if err == nil {
		t.Error("Expected error when creating duplicate pool")
	}
}

func TestPoolManager_AddRemoveCertificate(t *testing.T) {
	store := newMockStore()
	pm, _ := NewPoolManager(store)
	ctx := context.Background()
	
	cert, _ := generateTestCertificate("test-cert", "test.example.com")
	store.AddCertificate(ctx, cert)
	pm.CreatePool(ctx, "test-pool")

	err := pm.AddCertificateToPool(ctx, "test-pool", "test-cert")
	if err != nil {
		t.Errorf("Failed to add certificate to pool: %v", err)
	}

	certs, _ := pm.GetPool(ctx, "test-pool")
	if len(certs) != 1 {
		t.Errorf("Expected 1 certificate in pool, got %d", len(certs))
	}

	pm.RemoveCertificateFromPool(ctx, "test-pool", "test-cert")
	certs, _ = pm.GetPool(ctx, "test-pool")
	if len(certs) != 0 {
		t.Errorf("Expected 0 certificates in pool, got %d", len(certs))
	}
}

func TestValidator_ValidateCertificate(t *testing.T) {
	validator := NewValidator()
	cert, _ := generateTestCertificate("valid-cert", "valid.example.com")

	result, err := validator.ValidateCertificate(context.Background(), cert)
	if err != nil {
		t.Errorf("Failed to validate certificate: %v", err)
	}
	if !result.Valid {
		t.Errorf("Expected certificate to be valid, got invalid")
	}
}