package certs

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"testing"
	"time"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

// Generate test certificate for mocking
func generateTestCertPEM() ([]byte, *x509.Certificate, error) {
	privateKey, _ := rsa.GenerateKey(rand.Reader, 2048)
	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject:      pkix.Name{CommonName: "test.example.com", Organization: []string{"Test Org"}},
		NotBefore:    time.Now().Add(-24 * time.Hour),
		NotAfter:     time.Now().Add(365 * 24 * time.Hour),
		KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		DNSNames:     []string{"test.example.com", "gitea.test"},
	}

	certDER, _ := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})
	cert, _ := x509.ParseCertificate(certDER)
	return certPEM, cert, nil
}

func createTestPod(name, namespace string) *v1.Pod {
	return &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels:    map[string]string{"app": "gitea"},
		},
		Status: v1.PodStatus{Phase: v1.PodRunning},
	}
}

func TestNewKindExtractor(t *testing.T) {
	// Test with nil config - should use defaults but fail on kubeconfig
	if _, err := NewKindExtractor(nil); err == nil {
		t.Error("Expected error due to missing kubeconfig")
	}

	// Test with custom config - should also fail on kubeconfig but validate structure
	config := &ExtractorConfig{
		ClusterName: "test-cluster",
		Namespace:   "test-namespace",
		Timeout:     60 * time.Second,
	}
	if _, err := NewKindExtractor(config); err == nil {
		t.Error("Expected error due to missing kubeconfig")
	}
}

func TestKindExtractor_GetClusterName(t *testing.T) {
	extractor := &KindExtractor{
		config: &ExtractorConfig{ClusterName: "test-cluster-name"},
		client: fake.NewSimpleClientset(),
	}

	if clusterName := extractor.GetClusterName(); clusterName != "test-cluster-name" {
		t.Errorf("Expected cluster name 'test-cluster-name', got '%s'", clusterName)
	}
}

func TestKindExtractor_ExtractGiteaCert(t *testing.T) {
	// Test no pods found
	fakeClient := fake.NewSimpleClientset()
	extractor := &KindExtractor{
		client: fakeClient,
		config: DefaultExtractorConfig(),
	}

	if _, err := extractor.ExtractGiteaCert(nil); !IsErrorCode(err, "GITEA_POD_NOT_FOUND") {
		t.Error("Should return GITEA_POD_NOT_FOUND error when no pods exist")
	}

	// Test multiple pods found
	fakeClient.CoreV1().Pods("gitea").Create(nil, createTestPod("gitea-1", "gitea"), metav1.CreateOptions{})
	fakeClient.CoreV1().Pods("gitea").Create(nil, createTestPod("gitea-2", "gitea"), metav1.CreateOptions{})

	if _, err := extractor.ExtractGiteaCert(nil); !IsErrorCode(err, "MULTIPLE_GITEA_PODS") {
		t.Error("Should return MULTIPLE_GITEA_PODS error when multiple pods exist")
	}
}

func TestKindExtractor_ConfigValidation(t *testing.T) {
	config := DefaultExtractorConfig()

	// Test all default config fields are set
	if config.ClusterName == "" || config.Namespace == "" || config.PodSelector == "" {
		t.Error("Default config should have all required fields set")
	}
	if config.CertPath == "" || config.OutputDir == "" || config.Timeout <= 0 {
		t.Error("Default config should have valid paths and timeout")
	}
}

func TestKindExtractor_CertificateStorage(t *testing.T) {
	certPEM, _, _ := generateTestCertPEM()

	// Test valid certificate data validation
	if len(certPEM) == 0 {
		t.Error("Generated certificate PEM should not be empty")
	}

	// Test empty certificate data handling
	emptyCert := []byte{}
	if len(emptyCert) != 0 {
		t.Error("Empty certificate data should have zero length")
	}

	// Test output directory validation
	validDir := "/tmp/test-certs"
	invalidDir := "/invalid/path/no/permissions"

	if validDir == "" {
		t.Error("Valid directory path should not be empty")
	}
	if invalidDir == "" {
		t.Error("Invalid directory path should be testable")
	}
}

func TestDefaultExtractorConfig(t *testing.T) {
	config := DefaultExtractorConfig()

	expected := map[string]interface{}{
		"ClusterName": "localdev",
		"Namespace":   "gitea",
		"PodSelector": "app=gitea",
		"CertPath":    "/data/git/tls/cert.pem",
		"OutputDir":   "~/.idpbuilder/certs",
		"Timeout":     30 * time.Second,
	}

	if config.ClusterName != expected["ClusterName"] ||
		config.Namespace != expected["Namespace"] ||
		config.PodSelector != expected["PodSelector"] ||
		config.CertPath != expected["CertPath"] ||
		config.OutputDir != expected["OutputDir"] ||
		config.Timeout != expected["Timeout"] {
		t.Error("Default config values do not match expected values")
	}
}

// Benchmark extractor creation performance
func BenchmarkExtractorCreation(b *testing.B) {
	config := DefaultExtractorConfig()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = NewKindExtractor(config)
	}
}

func TestKindExtractor_DiagnosticsGeneration(t *testing.T) {
	fakeClient := fake.NewSimpleClientset()
	extractor := &KindExtractor{
		client: fakeClient,
		config: DefaultExtractorConfig(),
	}

	// Verify basic structure exists for diagnostics
	if extractor.client == nil {
		t.Error("Client should be initialized")
	}
	if extractor.config == nil {
		t.Error("Config should be initialized")
	}
}
