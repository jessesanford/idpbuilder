package extractor

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func createTestCertificate() ([]byte, *x509.Certificate, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}

	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization: []string{"Test Org"},
			CommonName:   "test.example.com",
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(365 * 24 * time.Hour),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	if err != nil {
		return nil, nil, err
	}

	certPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certDER,
	})

	cert, err := x509.ParseCertificate(certDER)
	if err != nil {
		return nil, nil, err
	}

	return certPEM, cert, nil
}

func TestNew(t *testing.T) {
	config := ExtractorConfig{
		ClusterName: "test-cluster",
		Timeout:     10 * time.Second,
		RetryCount:  2,
	}

	// This will fail due to invalid kubeconfig, but we're testing config validation
	extractor, err := New(config)
	assert.Error(t, err) // Expected since we don't have a valid kubeconfig
	assert.Nil(t, extractor)

	// Test default values
	config = ExtractorConfig{}
	extractor, err = New(config)
	assert.Error(t, err) // Expected since we don't have a valid kubeconfig
	assert.Nil(t, extractor)
}

func TestValidateCluster(t *testing.T) {
	// Create fake kubernetes client
	fakeClient := fake.NewSimpleClientset()

	// Create a test node with Kind cluster label
	node := &corev1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test-node",
			Labels: map[string]string{
				"io.x-k8s.kind.cluster": "test-cluster",
			},
		},
	}
	fakeClient.CoreV1().Nodes().Create(context.TODO(), node, metav1.CreateOptions{})

	extractor := &kindExtractor{
		client:      fakeClient,
		clusterName: "test-cluster",
		timeout:     30 * time.Second,
		retryCount:  3,
	}

	// Test successful validation
	clusterName, err := extractor.ValidateCluster(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, "test-cluster", clusterName)

	// Test wrong cluster name
	extractor.clusterName = "wrong-cluster"
	_, err = extractor.ValidateCluster(context.Background())
	assert.Error(t, err)
	var clusterErr *ClusterNotFoundError
	assert.ErrorAs(t, err, &clusterErr)
	assert.Equal(t, "wrong-cluster", clusterErr.Expected)
	assert.Equal(t, "test-cluster", clusterErr.Actual)
}

func TestValidateCluster_NoKindLabel(t *testing.T) {
	// Create fake kubernetes client with node without Kind label
	fakeClient := fake.NewSimpleClientset()
	node := &corev1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name:   "test-node",
			Labels: map[string]string{},
		},
	}
	fakeClient.CoreV1().Nodes().Create(context.TODO(), node, metav1.CreateOptions{})

	extractor := &kindExtractor{
		client:      fakeClient,
		clusterName: "test-cluster",
		timeout:     30 * time.Second,
		retryCount:  3,
	}

	_, err := extractor.ValidateCluster(context.Background())
	assert.Error(t, err)
	var clusterErr *ClusterNotFoundError
	assert.ErrorAs(t, err, &clusterErr)
	assert.Equal(t, "test-cluster", clusterErr.Expected)
}

func TestCertificateParseError(t *testing.T) {
	err := &CertificateParseError{
		Reason: "invalid PEM format",
	}
	assert.Equal(t, "failed to parse certificate: invalid PEM format", err.Error())
}

func TestClusterNotFoundError(t *testing.T) {
	// Test with actual cluster name
	err := &ClusterNotFoundError{
		Expected: "expected-cluster",
		Actual:   "actual-cluster",
	}
	assert.Contains(t, err.Error(), "expected 'expected-cluster' but connected to 'actual-cluster'")

	// Test with no actual cluster
	err = &ClusterNotFoundError{
		Expected: "expected-cluster",
	}
	assert.Contains(t, err.Error(), "Kind cluster 'expected-cluster' not found")
}