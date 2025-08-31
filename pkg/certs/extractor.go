package certs

import (
	"context"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-logr/logr"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// KindExtractor implements the KindCertExtractor interface for extracting certificates from Kind clusters
type KindExtractor struct {
	// client is the Kubernetes client for interacting with the cluster
	client kubernetes.Interface

	// config holds the extractor configuration
	config *ExtractorConfig

	// logger provides structured logging
	logger logr.Logger
}

// NewKindExtractor creates a new KindExtractor with the specified configuration
func NewKindExtractor(config *ExtractorConfig) (*KindExtractor, error) {
	if config == nil {
		config = DefaultExtractorConfig()
	}

	// Create Kubernetes client using current context
	kubeConfig, err := clientcmd.BuildConfigFromFlags("", "")
	if err != nil {
		return nil, WrapError(err, "KUBECONFIG_LOAD_FAILED", "Failed to load kubeconfig")
	}

	client, err := kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		return nil, WrapError(err, "CLIENT_CREATE_FAILED", "Failed to create Kubernetes client")
	}

	// Verify cluster connectivity
	ctx, cancel := context.WithTimeout(context.Background(), config.Timeout)
	defer cancel()

	_, err = client.CoreV1().Nodes().List(ctx, metav1.ListOptions{Limit: 1})
	if err != nil {
		return nil, ErrClusterConnection.WithContext("cluster", config.ClusterName).Wrap(err)
	}

	logger := log.FromContext(context.Background()).WithName("cert-extractor")

	return &KindExtractor{
		client: client,
		config: config,
		logger: logger,
	}, nil
}

// ExtractGiteaCert extracts the Gitea certificate from the Kind cluster
func (k *KindExtractor) ExtractGiteaCert(ctx context.Context) (*x509.Certificate, error) {
	startTime := time.Now()
	diagnostics := &CertDiagnostics{
		ClusterConnected: true,
		Warnings:         make([]string, 0),
	}

	k.logger.Info("Starting Gitea certificate extraction", "cluster", k.config.ClusterName)

	// Find Gitea pod
	podName, err := k.findGiteaPod(ctx)
	if err != nil {
		return nil, err
	}
	diagnostics.PodsFound = 1

	k.logger.Info("Found Gitea pod", "pod", podName, "namespace", k.config.Namespace)

	// Extract certificate from pod
	certData, err := k.extractCertFromPod(ctx, podName)
	if err != nil {
		return nil, err
	}
	diagnostics.CertificateFound = true

	// Parse certificate
	cert, err := k.parseCertificate(certData)
	if err != nil {
		return nil, err
	}
	diagnostics.CertificateParsed = true

	// Store certificate locally
	err = k.storeCertificate(cert, certData)
	if err != nil {
		return nil, err
	}

	diagnostics.ExtractionDuration = time.Since(startTime)

	k.logger.Info("Certificate extraction completed successfully",
		"duration", diagnostics.ExtractionDuration,
		"subject", cert.Subject.String(),
		"expires", cert.NotAfter)

	return cert, nil
}

// GetClusterName returns the name of the configured Kind cluster
func (k *KindExtractor) GetClusterName() string {
	return k.config.ClusterName
}

// findGiteaPod locates the Gitea pod in the cluster using the configured selector
func (k *KindExtractor) findGiteaPod(ctx context.Context) (string, error) {
	k.logger.Info("Searching for Gitea pod", "namespace", k.config.Namespace, "selector", k.config.PodSelector)

	// Parse label selector
	labelSelector := k.config.PodSelector

	// List pods with the specified selector
	pods, err := k.client.CoreV1().Pods(k.config.Namespace).List(ctx, metav1.ListOptions{
		LabelSelector: labelSelector,
	})
	if err != nil {
		return "", ErrClusterConnection.WithContext("namespace", k.config.Namespace).
			WithContext("selector", labelSelector).Wrap(err)
	}

	// Filter running pods
	var runningPods []v1.Pod
	for _, pod := range pods.Items {
		if pod.Status.Phase == v1.PodRunning {
			runningPods = append(runningPods, pod)
		}
	}

	// Check pod count
	if len(runningPods) == 0 {
		return "", ErrGiteaPodNotFound.WithContext("namespace", k.config.Namespace).
			WithContext("selector", labelSelector).
			WithContext("total_pods", len(pods.Items)).
			WithContext("running_pods", len(runningPods))
	}

	if len(runningPods) > 1 {
		podNames := make([]string, len(runningPods))
		for i, pod := range runningPods {
			podNames[i] = pod.Name
		}
		return "", ErrMultipleGiteaPods.WithContext("namespace", k.config.Namespace).
			WithContext("selector", labelSelector).
			WithContext("pod_names", strings.Join(podNames, ", "))
	}

	return runningPods[0].Name, nil
}

// extractCertFromPod executes a command in the pod to retrieve the certificate
func (k *KindExtractor) extractCertFromPod(ctx context.Context, podName string) ([]byte, error) {
	k.logger.Info("Extracting certificate from pod", "pod", podName, "path", k.config.CertPath)

	// Build kubectl command to read the certificate file
	cmd := exec.CommandContext(ctx, "kubectl", "exec",
		"-n", k.config.Namespace,
		podName,
		"--",
		"cat", k.config.CertPath)

	// Execute command
	output, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			stderr := string(exitErr.Stderr)
			if strings.Contains(stderr, "No such file") || strings.Contains(stderr, "not found") {
				return nil, ErrCertificateNotFound.WithContext("pod", podName).
					WithContext("path", k.config.CertPath).
					WithContext("stderr", stderr)
			}
		}
		return nil, ErrCertificateRead.WithContext("pod", podName).
			WithContext("path", k.config.CertPath).Wrap(err)
	}

	if len(output) == 0 {
		return nil, ErrCertificateNotFound.WithContext("pod", podName).
			WithContext("path", k.config.CertPath).
			WithContext("reason", "empty file")
	}

	k.logger.Info("Successfully read certificate data", "pod", podName, "size", len(output))

	return output, nil
}

// parseCertificate parses PEM-encoded certificate data into an x509.Certificate
func (k *KindExtractor) parseCertificate(certData []byte) (*x509.Certificate, error) {
	k.logger.Info("Parsing certificate data", "size", len(certData))

	// Decode PEM block
	block, _ := pem.Decode(certData)
	if block == nil {
		return nil, ErrCertificateParse.WithContext("reason", "no PEM block found").
			WithContext("data_preview", string(certData[:min(100, len(certData))]))
	}

	if block.Type != "CERTIFICATE" {
		return nil, ErrCertificateParse.WithContext("reason", "not a certificate block").
			WithContext("block_type", block.Type)
	}

	// Parse X.509 certificate
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, ErrCertificateParse.WithContext("reason", "invalid certificate data").Wrap(err)
	}

	k.logger.Info("Successfully parsed certificate",
		"subject", cert.Subject.String(),
		"issuer", cert.Issuer.String(),
		"not_before", cert.NotBefore,
		"not_after", cert.NotAfter)

	return cert, nil
}

// storeCertificate saves the certificate to the local filesystem
func (k *KindExtractor) storeCertificate(cert *x509.Certificate, certData []byte) error {
	// Expand tilde in output directory
	outputDir := k.config.OutputDir
	if strings.HasPrefix(outputDir, "~/") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return ErrCertificateStore.WithContext("reason", "cannot determine home directory").Wrap(err)
		}
		outputDir = filepath.Join(homeDir, outputDir[2:])
	}

	k.logger.Info("Storing certificate", "output_dir", outputDir)

	// Ensure output directory exists
	err := os.MkdirAll(outputDir, 0755)
	if err != nil {
		return ErrCertificateStore.WithContext("directory", outputDir).
			WithContext("reason", "failed to create directory").Wrap(err)
	}

	// Generate filename based on certificate subject
	filename := k.generateCertFilename(cert)
	certPath := filepath.Join(outputDir, filename)

	// Write certificate to file
	err = os.WriteFile(certPath, certData, 0644)
	if err != nil {
		return ErrCertificateStore.WithContext("file", certPath).
			WithContext("reason", "failed to write file").Wrap(err)
	}

	k.logger.Info("Certificate stored successfully", "path", certPath)

	return nil
}

// generateCertFilename generates a filename for the certificate based on its properties
func (k *KindExtractor) generateCertFilename(cert *x509.Certificate) string {
	// Use common name or first DNS name as base
	name := cert.Subject.CommonName
	if name == "" && len(cert.DNSNames) > 0 {
		name = cert.DNSNames[0]
	}
	if name == "" {
		name = "gitea"
	}

	// Sanitize filename
	name = strings.ReplaceAll(name, "*", "wildcard")
	name = strings.ReplaceAll(name, ".", "_")
	name = strings.ReplaceAll(name, ":", "_")

	// Add timestamp to avoid conflicts
	timestamp := time.Now().Format("20060102-150405")

	return fmt.Sprintf("%s-%s.crt", name, timestamp)
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
