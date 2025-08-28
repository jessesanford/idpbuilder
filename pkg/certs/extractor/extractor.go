package extractor

import (
	"context"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"time"

	"k8s.io/client-go/kubernetes"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// CertificateExtractor extracts certificates from Kind cluster pods
type CertificateExtractor interface {
	ExtractFromPod(ctx context.Context, namespace, podName, certPath string) (*x509.Certificate, error)
	GetGiteaPod(ctx context.Context) (podName, namespace string, err error)
	ValidateCluster(ctx context.Context) (clusterName string, err error)
}

// ExtractorConfig holds configuration for the certificate extractor
type ExtractorConfig struct {
	KubeConfig  string        // Path to kubeconfig file
	ClusterName string        // Expected Kind cluster name (e.g., "idpbuilder")
	Timeout     time.Duration // Operation timeout (default: 30s)
	RetryCount  int           // Number of retries for transient failures (default: 3)
}

// kindExtractor implements CertificateExtractor for Kind clusters
type kindExtractor struct {
	client      kubernetes.Interface
	clusterName string
	timeout     time.Duration
	retryCount  int
	kubeConfig  string // Store kubeconfig path for exec operations
}

// New creates a new certificate extractor
func New(config ExtractorConfig) (CertificateExtractor, error) {
	// Set defaults
	if config.Timeout == 0 {
		config.Timeout = 30 * time.Second
	}
	if config.RetryCount == 0 {
		config.RetryCount = 3
	}

	// Setup kubernetes client
	client, kubeConfigPath, err := setupKubernetesClient(config.KubeConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to setup kubernetes client: %w", err)
	}

	return &kindExtractor{
		client:      client,
		clusterName: config.ClusterName,
		timeout:     config.Timeout,
		retryCount:  config.RetryCount,
		kubeConfig:  kubeConfigPath,
	}, nil
}

// ExtractFromPod retrieves a certificate from a running pod
func (e *kindExtractor) ExtractFromPod(ctx context.Context, namespace, podName, certPath string) (*x509.Certificate, error) {
	// Add timeout to context
	ctx, cancel := context.WithTimeout(ctx, e.timeout)
	defer cancel()

	var lastErr error
	for attempt := 0; attempt < e.retryCount; attempt++ {
		// Execute cat command in pod to read certificate
		certData, err := e.execInPod(ctx, namespace, podName, []string{"cat", certPath})
		if err != nil {
			lastErr = err
			if attempt < e.retryCount-1 {
				time.Sleep(time.Duration(attempt+1) * time.Second)
				continue
			}
			return nil, fmt.Errorf("failed to read certificate from pod after %d attempts: %w", e.retryCount, err)
		}

		// Parse PEM encoded certificate
		block, _ := pem.Decode(certData)
		if block == nil {
			return nil, &CertificateParseError{
				Reason: "failed to decode PEM certificate",
			}
		}

		// Parse X.509 certificate
		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return nil, &CertificateParseError{
				Reason: fmt.Sprintf("failed to parse X.509 certificate: %v", err),
			}
		}

		return cert, nil
	}

	return nil, lastErr
}

// ValidateCluster ensures we're connected to the right Kind cluster
func (e *kindExtractor) ValidateCluster(ctx context.Context) (string, error) {
	// Query nodes to verify cluster name
	nodes, err := e.client.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		return "", fmt.Errorf("failed to list nodes: %w", err)
	}

	// Check for Kind-specific labels/annotations
	for _, node := range nodes.Items {
		if clusterName, ok := node.Labels["io.x-k8s.kind.cluster"]; ok {
			if e.clusterName != "" && clusterName != e.clusterName {
				return "", &ClusterNotFoundError{
					Expected: e.clusterName,
					Actual:   clusterName,
				}
			}
			return clusterName, nil
		}
	}

	return "", &ClusterNotFoundError{Expected: e.clusterName}
}