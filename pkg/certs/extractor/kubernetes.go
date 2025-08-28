package extractor

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

// setupKubernetesClient creates a Kubernetes client from kubeconfig
// Returns the client and the actual kubeconfig path used
func setupKubernetesClient(kubeconfig string) (kubernetes.Interface, string, error) {
	// If no kubeconfig specified, try default locations
	if kubeconfig == "" {
		kubeconfig = getDefaultKubeconfig()
	}

	// Check if file exists
	if _, err := os.Stat(kubeconfig); err != nil {
		return nil, "", fmt.Errorf("kubeconfig not found at %s: %w", kubeconfig, err)
	}

	// Validate this is a Kind context
	if err := validateKindContext(kubeconfig, ""); err != nil {
		return nil, "", fmt.Errorf("kubeconfig validation failed: %w", err)
	}

	// Build config from kubeconfig file
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, "", fmt.Errorf("failed to build config: %w", err)
	}

	// Create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, "", fmt.Errorf("failed to create clientset: %w", err)
	}

	// Test connection
	_, err = clientset.ServerVersion()
	if err != nil {
		return nil, "", fmt.Errorf("failed to connect to cluster: %w", err)
	}

	return clientset, kubeconfig, nil
}

// getDefaultKubeconfig returns the default kubeconfig path
func getDefaultKubeconfig() string {
	// First check KUBECONFIG environment variable
	if kubeconfig := os.Getenv("KUBECONFIG"); kubeconfig != "" {
		return kubeconfig
	}

	// Fall back to ~/.kube/config
	if home := homedir.HomeDir(); home != "" {
		return filepath.Join(home, ".kube", "config")
	}

	return ""
}

// validateKindContext ensures the kubeconfig context points to a Kind cluster
func validateKindContext(kubeconfig string, expectedCluster string) error {
	config, err := clientcmd.LoadFromFile(kubeconfig)
	if err != nil {
		return fmt.Errorf("failed to load kubeconfig: %w", err)
	}

	currentContext := config.CurrentContext
	if expectedCluster != "" && !strings.Contains(currentContext, expectedCluster) {
		return fmt.Errorf("current context %s doesn't match expected cluster %s",
			currentContext, expectedCluster)
	}

	// Check if it's a Kind cluster (context usually contains "kind-")
	if !strings.Contains(currentContext, "kind-") {
		return fmt.Errorf("current context %s doesn't appear to be a Kind cluster", currentContext)
	}

	return nil
}