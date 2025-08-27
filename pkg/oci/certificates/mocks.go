package certificates

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"fmt"
	"math/big"
	"net/http"
	"sync"
	"time"
)

// MockCertificateService implements CertificateService for testing
type MockCertificateService struct {
	mu sync.RWMutex

	// Configurable functions
	LoadBundleFunc        func(ctx context.Context, path string, format CertFormat) (*CertBundle, error)
	ValidateFunc          func(cert *x509.Certificate) error
	GetTLSConfigFunc      func() (*tls.Config, error)
	AddCAFunc             func(cert *x509.Certificate) error
	RemoveCAFunc          func(cert *x509.Certificate) error
	ListCertificatesFunc  func() ([]*x509.Certificate, error)
	RotateCertificateFunc func(old, new *x509.Certificate) error
	LoadGiteaFunc         func(ctx context.Context, url string) (*CertBundle, error)

	// State
	verificationMode VerificationMode
	certificates     []*x509.Certificate
	caPool           *x509.CertPool
	callCount        map[string]int
	errorOnNthCall   map[string]int

	// Test helpers
	bundleCache map[string]*CertBundle
}

// NewMockCertificateService creates a new mock certificate service
func NewMockCertificateService() *MockCertificateService {
	return &MockCertificateService{
		verificationMode: VerificationModeStrict,
		certificates:     make([]*x509.Certificate, 0),
		caPool:           x509.NewCertPool(),
		callCount:        make(map[string]int),
		errorOnNthCall:   make(map[string]int),
		bundleCache:      make(map[string]*CertBundle),
	}
}

// LoadCertificateBundle implements CertificateService
func (m *MockCertificateService) LoadCertificateBundle(ctx context.Context, path string, format CertFormat) (*CertBundle, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.callCount["LoadCertificateBundle"]++
	if errCall, exists := m.errorOnNthCall["LoadCertificateBundle"]; exists {
		if m.callCount["LoadCertificateBundle"] == errCall {
			return nil, fmt.Errorf("mock error on call %d", errCall)
		}
	}

	if m.LoadBundleFunc != nil {
		return m.LoadBundleFunc(ctx, path, format)
	}

	// Return cached bundle if available
	if bundle, exists := m.bundleCache[path]; exists {
		return bundle, nil
	}

	// Generate test bundle
	bundle := &CertBundle{
		Certificates: []*x509.Certificate{GenerateTestCertificate()},
		CAs:          []*x509.Certificate{GenerateTestCA()},
		Format:       format,
		LoadedAt:     time.Now(),
		Source:       path,
	}

	m.bundleCache[path] = bundle
	return bundle, nil
}

// SetVerificationMode implements CertificateService
func (m *MockCertificateService) SetVerificationMode(mode VerificationMode) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.verificationMode = mode
	return nil
}

// ValidateCertificate implements CertificateService
func (m *MockCertificateService) ValidateCertificate(cert *x509.Certificate) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if cert == nil {
		return fmt.Errorf("certificate is nil")
	}

	m.callCount["ValidateCertificate"]++
	if errCall, exists := m.errorOnNthCall["ValidateCertificate"]; exists {
		if m.callCount["ValidateCertificate"] == errCall {
			return fmt.Errorf("mock validation error on call %d", errCall)
		}
	}

	if m.ValidateFunc != nil {
		return m.ValidateFunc(cert)
	}

	// Mock validation logic
	if cert.NotAfter.Before(time.Now()) {
		return fmt.Errorf("certificate expired")
	}

	return nil
}

// LoadGiteaCertificate implements CertificateService
func (m *MockCertificateService) LoadGiteaCertificate(ctx context.Context, giteaURL string) (*CertBundle, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.LoadGiteaFunc != nil {
		return m.LoadGiteaFunc(ctx, giteaURL)
	}

	// Return mock Gitea certificate bundle
	return &CertBundle{
		Certificates: []*x509.Certificate{GenerateTestCertificate()},
		CAs:          []*x509.Certificate{GenerateTestCA()},
		Format:       CertFormatPEM,
		LoadedAt:     time.Now(),
		Source:       giteaURL,
	}, nil
}

// GetTLSConfig implements CertificateService
func (m *MockCertificateService) GetTLSConfig() (*tls.Config, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.GetTLSConfigFunc != nil {
		return m.GetTLSConfigFunc()
	}

	config := &tls.Config{
		RootCAs: m.caPool,
	}

	if m.verificationMode == VerificationModeSkip {
		config.InsecureSkipVerify = true
	}

	return config, nil
}

// AddCACertificate implements CertificateService
func (m *MockCertificateService) AddCACertificate(cert *x509.Certificate) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.AddCAFunc != nil {
		return m.AddCAFunc(cert)
	}

	m.caPool.AddCert(cert)
	return nil
}

// RemoveCACertificate implements CertificateService
func (m *MockCertificateService) RemoveCACertificate(cert *x509.Certificate) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.RemoveCAFunc != nil {
		return m.RemoveCAFunc(cert)
	}

	// Note: x509.CertPool doesn't support removal, so this is a no-op in mock
	return nil
}

// ListCertificates implements CertificateService
func (m *MockCertificateService) ListCertificates() ([]*x509.Certificate, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.ListCertificatesFunc != nil {
		return m.ListCertificatesFunc()
	}

	return m.certificates, nil
}

// RotateCertificate implements CertificateService
func (m *MockCertificateService) RotateCertificate(old, new *x509.Certificate) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.RotateCertificateFunc != nil {
		return m.RotateCertificateFunc(old, new)
	}

	// Mock rotation: replace old with new
	for i, cert := range m.certificates {
		if cert.Equal(old) {
			m.certificates[i] = new
			return nil
		}
	}

	return fmt.Errorf("certificate not found for rotation")
}

// GetCertificatePool implements CertificateService
func (m *MockCertificateService) GetCertificatePool() *x509.CertPool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.caPool
}

// Test helper methods
func (m *MockCertificateService) SetErrorOnNthCall(method string, callNumber int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.errorOnNthCall[method] = callNumber
}

func (m *MockCertificateService) GetCallCount(method string) int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.callCount[method]
}

func (m *MockCertificateService) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.callCount = make(map[string]int)
	m.errorOnNthCall = make(map[string]int)
	m.bundleCache = make(map[string]*CertBundle)
	m.certificates = make([]*x509.Certificate, 0)
	m.caPool = x509.NewCertPool()
}

// MockRegistryClient implements RegistryClient for testing
type MockRegistryClient struct {
	mu sync.RWMutex

	pushFunc  func(ctx context.Context, ref string, content []byte) error
	pullFunc  func(ctx context.Context, ref string) ([]byte, error)
	loginFunc func(ctx context.Context, registry, username, password string) error
	transport http.RoundTripper

	// Test state
	pushedRefs    []string
	pulledRefs    []string
	loginCalls    []string
	certValidated bool
}

// NewMockRegistryClient creates a new mock registry client
func NewMockRegistryClient() *MockRegistryClient {
	return &MockRegistryClient{
		pushedRefs: make([]string, 0),
		pulledRefs: make([]string, 0),
		loginCalls: make([]string, 0),
	}
}

// Push implements RegistryClient
func (m *MockRegistryClient) Push(ctx context.Context, ref string, content []byte) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.pushFunc != nil {
		return m.pushFunc(ctx, ref, content)
	}

	m.pushedRefs = append(m.pushedRefs, ref)
	m.certValidated = true // Simulate certificate validation
	return nil
}

// Pull implements RegistryClient
func (m *MockRegistryClient) Pull(ctx context.Context, ref string) ([]byte, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.pullFunc != nil {
		return m.pullFunc(ctx, ref)
	}

	m.pulledRefs = append(m.pulledRefs, ref)
	return []byte("mock-content"), nil
}

// Login implements RegistryClient
func (m *MockRegistryClient) Login(ctx context.Context, registry, username, password string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.loginFunc != nil {
		return m.loginFunc(ctx, registry, username, password)
	}

	m.loginCalls = append(m.loginCalls, registry)
	return nil
}

// SetTransport implements RegistryClient
func (m *MockRegistryClient) SetTransport(transport http.RoundTripper) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.transport = transport
	return nil
}

// Test helper methods
func (m *MockRegistryClient) GetPushedRefs() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return append([]string{}, m.pushedRefs...)
}

func (m *MockRegistryClient) GetPulledRefs() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return append([]string{}, m.pulledRefs...)
}

func (m *MockRegistryClient) CertificateValidated() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.certValidated
}

func (m *MockRegistryClient) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.pushedRefs = make([]string, 0)
	m.pulledRefs = make([]string, 0)
	m.loginCalls = make([]string, 0)
	m.certValidated = false
	m.transport = nil
}

// Test certificate generation utilities
func GenerateTestCA() *x509.Certificate {
	template := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization: []string{"Test CA"},
			CommonName:   "Test Root CA",
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		BasicConstraintsValid: true,
		IsCA:                  true,
	}

	priv, _ := rsa.GenerateKey(rand.Reader, 2048)
	certDER, _ := x509.CreateCertificate(rand.Reader, template, template, &priv.PublicKey, priv)
	cert, _ := x509.ParseCertificate(certDER)

	return cert
}

func GenerateTestCertificate() *x509.Certificate {
	template := &x509.Certificate{
		SerialNumber: big.NewInt(2),
		Subject: pkix.Name{
			Organization: []string{"Test Org"},
			CommonName:   "test.example.com",
		},
		NotBefore:   time.Now(),
		NotAfter:    time.Now().AddDate(1, 0, 0),
		KeyUsage:    x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		DNSNames:    []string{"test.example.com"},
	}

	ca := GenerateTestCA()
	priv, _ := rsa.GenerateKey(rand.Reader, 2048)
	caPriv, _ := rsa.GenerateKey(rand.Reader, 2048)

	certDER, _ := x509.CreateCertificate(rand.Reader, template, ca, &priv.PublicKey, caPriv)
	cert, _ := x509.ParseCertificate(certDER)

	return cert
}

func GenerateExpiredCertificate() *x509.Certificate {
	template := &x509.Certificate{
		SerialNumber: big.NewInt(3),
		Subject: pkix.Name{
			Organization: []string{"Test Org"},
			CommonName:   "expired.example.com",
		},
		NotBefore:   time.Now().AddDate(-2, 0, 0),
		NotAfter:    time.Now().AddDate(-1, 0, 0), // Expired
		KeyUsage:    x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		DNSNames:    []string{"expired.example.com"},
	}

	ca := GenerateTestCA()
	priv, _ := rsa.GenerateKey(rand.Reader, 2048)
	caPriv, _ := rsa.GenerateKey(rand.Reader, 2048)

	certDER, _ := x509.CreateCertificate(rand.Reader, template, ca, &priv.PublicKey, caPriv)
	cert, _ := x509.ParseCertificate(certDER)

	return cert
}

func GenerateSelfSignedCertificate() *x509.Certificate {
	template := &x509.Certificate{
		SerialNumber: big.NewInt(4),
		Subject: pkix.Name{
			Organization: []string{"Self Signed"},
			CommonName:   "selfsigned.example.com",
		},
		NotBefore:   time.Now(),
		NotAfter:    time.Now().AddDate(1, 0, 0),
		KeyUsage:    x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		DNSNames:    []string{"selfsigned.example.com"},
	}

	priv, _ := rsa.GenerateKey(rand.Reader, 2048)
	certDER, _ := x509.CreateCertificate(rand.Reader, template, template, &priv.PublicKey, priv)
	cert, _ := x509.ParseCertificate(certDER)

	return cert
}
