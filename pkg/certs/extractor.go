package certs

import (
	"context"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/cnoe-io/idpbuilder/pkg/kind"
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

var extractorLog = log.Log.WithName("cert-extractor")

// KindExtractor implements KindCertExtractor for extracting certificates from Kind clusters
type KindExtractor struct {
	config     *ExtractorConfig
	kubeClient kubernetes.Interface
	restConfig *rest.Config
	logger     logr.Logger
}

// NewKindExtractor creates a new KindExtractor with the provided configuration
func NewKindExtractor(config *ExtractorConfig) (*KindExtractor, error) {
	if config == nil {
		config = DefaultExtractorConfig()
	}

	// Get kubeconfig for the Kind cluster
	kubeConfig, err := getKubeConfig(config.ClusterName)
	if err != nil {
		return nil, NewCertificateError("kubeconfig_setup", err, 
			"failed to get kubeconfig for cluster", []string{
				"Verify cluster exists: kind get clusters",
				"Export kubeconfig: kind export kubeconfig --name " + config.ClusterName,
			})
	}

	// Create Kubernetes client
	kubeClient, err := kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		return nil, NewCertificateError("client_creation", err,
			"failed to create Kubernetes client", []string{
				"Check cluster connectivity",
				"Verify kubeconfig is valid",
			})
	}

	return &KindExtractor{
		config:     config,
		kubeClient: kubeClient,
		restConfig: kubeConfig,
		logger:     extractorLog,
	}, nil
}

// ExtractGiteaCert extracts the Gitea certificate from the Kind cluster
func (e *KindExtractor) ExtractGiteaCert(ctx context.Context) (*x509.Certificate, error) {
	e.logger.Info("Starting Gitea certificate extraction", 
		"cluster", e.config.ClusterName, "namespace", e.config.GiteaNamespace)

	// Find Gitea pod
	pod, err := e.findGiteaPod(ctx)
	if err != nil {
		return nil, err
	}

	e.logger.Info("Found Gitea pod", "pod", pod.Name, "namespace", pod.Namespace)

	// Extract certificate from pod
	certData, err := e.extractCertFromPod(ctx, pod)
	if err != nil {
		return nil, err
	}

	// Parse certificate
	cert, err := e.parseCertificate(certData)
	if err != nil {
		return nil, err
	}

	// Validate certificate
	if err := e.ValidateCertificate(cert); err != nil {
		e.logger.V(1).Info("Certificate validation failed", "error", err)
		return cert, err // Return cert even if validation fails for diagnostic purposes
	}

	// Store certificate
	if err := e.storeCertificate(cert, certData); err != nil {
		e.logger.Error(err, "Failed to store certificate, but extraction succeeded")
		// Don't return error here - extraction succeeded even if storage failed
	}

	e.logger.Info("Successfully extracted Gitea certificate", 
		"subject", cert.Subject.String(), "expires", cert.NotAfter)

	return cert, nil
}

// GetClusterName returns the configured cluster name
func (e *KindExtractor) GetClusterName() (string, error) {
	if e.config.ClusterName == "" {
		return "", NewCertificateError("configuration", fmt.Errorf("cluster name not configured"),
			"no cluster name specified", []string{"Set cluster name in configuration"})
	}
	return e.config.ClusterName, nil
}

// ValidateCertificate validates the extracted certificate
func (e *KindExtractor) ValidateCertificate(cert *x509.Certificate) error {
	if cert == nil {
		return NewCertificateError("validation", fmt.Errorf("certificate is nil"),
			"no certificate provided for validation", []string{"Ensure certificate was extracted properly"})
	}

	now := time.Now()
	
	// Check if certificate is expired
	if now.After(cert.NotAfter) {
		return NewCertificateError("validation", ErrCertificateExpired.Cause,
			fmt.Sprintf("certificate expired on %v", cert.NotAfter), ErrCertificateExpired.Suggestions)
	}

	// Check if certificate is not yet valid
	if now.Before(cert.NotBefore) {
		return NewCertificateError("validation", fmt.Errorf("certificate not yet valid"),
			fmt.Sprintf("certificate valid from %v", cert.NotBefore), []string{
				"Check system clock synchronization",
				"Wait until certificate becomes valid",
			})
	}

	// Warn if certificate expires soon (within 30 days)
	if now.Add(30 * 24 * time.Hour).After(cert.NotAfter) {
		e.logger.Info("Certificate expires soon", "expires", cert.NotAfter, 
			"days_remaining", int(cert.NotAfter.Sub(now).Hours()/24))
	}

	return nil
}

// findGiteaPod locates the Gitea pod in the cluster
func (e *KindExtractor) findGiteaPod(ctx context.Context) (*corev1.Pod, error) {
	pods, err := e.kubeClient.CoreV1().Pods(e.config.GiteaNamespace).List(ctx, metav1.ListOptions{
		LabelSelector: e.config.GiteaPodLabelSelector,
	})
	if err != nil {
		return nil, NewCertificateError("pod_listing", err,
			"failed to list pods in namespace", []string{
				"Check if namespace exists: kubectl get namespace " + e.config.GiteaNamespace,
				"Verify cluster connectivity",
			})
	}

	if len(pods.Items) == 0 {
		return nil, ErrGiteaPodNotFound
	}

	// Find a running pod
	for _, pod := range pods.Items {
		if pod.Status.Phase == corev1.PodRunning {
			return &pod, nil
		}
	}

	return nil, NewCertificateError("pod_discovery", fmt.Errorf("no running gitea pods found"),
		"found pods but none are in running state", []string{
			"Check pod status: kubectl get pods -n " + e.config.GiteaNamespace,
			"Wait for pods to start or troubleshoot startup issues",
		})
}

// extractCertFromPod extracts the certificate file from the Gitea pod
func (e *KindExtractor) extractCertFromPod(ctx context.Context, pod *corev1.Pod) ([]byte, error) {
	// Create exec command
	req := e.kubeClient.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(pod.Name).
		Namespace(pod.Namespace).
		SubResource("exec")

	req.VersionedParams(&corev1.PodExecOptions{
		Command: []string{"cat", e.config.CertPath},
		Stdout:  true,
		Stderr:  true,
	}, scheme.ParameterCodec)

	executor, err := remotecommand.NewSPDYExecutor(e.restConfig, "POST", req.URL())
	if err != nil {
		return nil, NewCertificateError("exec_setup", err,
			"failed to create command executor", []string{
				"Check cluster connectivity",
				"Verify pod is accessible",
			})
	}

	// Capture output
	var stdout, stderr strings.Builder
	err = executor.Stream(remotecommand.StreamOptions{
		Stdout: &stdout,
		Stderr: &stderr,
	})

	if err != nil {
		stderrStr := stderr.String()
		if strings.Contains(stderrStr, "No such file") || strings.Contains(stderrStr, "not found") {
			return nil, ErrCertificateNotFound
		}
		return nil, NewCertificateError("certificate_read", err,
			fmt.Sprintf("failed to read certificate from pod, stderr: %s", stderrStr), []string{
				"Verify certificate path: " + e.config.CertPath,
				"Check pod has certificate configured",
			})
	}

	certData := stdout.String()
	if certData == "" {
		return nil, NewCertificateError("certificate_read", fmt.Errorf("empty certificate data"),
			"certificate file exists but is empty", []string{
				"Check Gitea configuration",
				"Verify HTTPS is enabled in Gitea",
			})
	}

	return []byte(certData), nil
}

// parseCertificate parses the PEM-encoded certificate data
func (e *KindExtractor) parseCertificate(data []byte) (*x509.Certificate, error) {
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, NewCertificateError("certificate_parsing", fmt.Errorf("no PEM data found"),
			"certificate data is not in PEM format", []string{
				"Verify certificate file format",
				"Check for certificate corruption",
			})
	}

	if block.Type != "CERTIFICATE" {
		return nil, NewCertificateError("certificate_parsing", 
			fmt.Errorf("invalid PEM block type: %s", block.Type),
			"PEM block is not a certificate", []string{
				"Verify you're reading the certificate file, not a key file",
				"Check certificate generation process",
			})
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, NewCertificateError("certificate_parsing", err,
			"failed to parse certificate data", []string{
				"Verify certificate is valid X.509 format",
				"Check for certificate corruption",
			})
	}

	return cert, nil
}

// storeCertificate stores the certificate to the configured location
func (e *KindExtractor) storeCertificate(cert *x509.Certificate, data []byte) error {
	storagePath := e.config.StoragePath
	
	// Expand home directory
	if strings.HasPrefix(storagePath, "~") {
		home, err := os.UserHomeDir()
		if err != nil {
			return NewCertificateError("path_resolution", err,
				"failed to get user home directory", []string{
					"Specify absolute path instead of ~",
					"Check user environment",
				})
		}
		storagePath = filepath.Join(home, storagePath[1:])
	}

	// Create directory if it doesn't exist
	dir := filepath.Dir(storagePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return NewCertificateError("directory_creation", err,
			"failed to create certificate storage directory", []string{
				"Check write permissions for parent directory",
				"Verify disk space availability",
			})
	}

	// Write certificate file
	if err := os.WriteFile(storagePath, data, 0644); err != nil {
		return ErrStoragePermission
	}

	e.logger.Info("Certificate stored successfully", "path", storagePath)
	return nil
}

// getKubeConfig gets the kubeconfig for the specified Kind cluster
func getKubeConfig(clusterName string) (*rest.Config, error) {
	// First try to use the Kind cluster's kubeconfig
	kubeConfigPath := filepath.Join(os.TempDir(), fmt.Sprintf("kind-config-%s", clusterName))
	
	// Export kubeconfig for the Kind cluster (this would be handled by kind package)
	// For now, try to load from default locations
	
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	configOverrides := &clientcmd.ConfigOverrides{}
	
	if clusterName != "" {
		configOverrides.CurrentContext = fmt.Sprintf("kind-%s", clusterName)
	}
	
	config := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)
	return config.ClientConfig()
}