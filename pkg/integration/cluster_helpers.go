//go:build integration

package integration

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// Credentials holds authentication information for registry access
type Credentials struct {
	Username string // Registry username
	Password string // Registry password
	Token    string // Optional token
}

// CreateTestCluster creates an idpbuilder test cluster
func CreateTestCluster(t *testing.T) (*ClusterInfo, error) {
	// Generate unique cluster name
	clusterName := fmt.Sprintf("test-cluster-%d", time.Now().Unix())

	// Prepare idpbuilder create command
	cmd := exec.Command("idpbuilder", "create", "--name", clusterName)
	cmd.Env = os.Environ()

	// Execute cluster creation
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to create test cluster %s: %w\nOutput: %s", clusterName, err, string(output))
	}

	// Parse output to extract cluster information
	namespace := "idpbuilder-system"
	kubeContext := fmt.Sprintf("kind-%s", clusterName)

	// Create cleanup function
	cleanupFunc := func() {
		cleanupCmd := exec.Command("idpbuilder", "delete", "--name", clusterName)
		if cleanupErr := cleanupCmd.Run(); cleanupErr != nil {
			t.Logf("Warning: Failed to cleanup cluster %s: %v", clusterName, cleanupErr)
		}
	}

	// Verify cluster is accessible
	if err := verifyClusterAccess(kubeContext); err != nil {
		cleanupFunc() // Cleanup on failure
		return nil, fmt.Errorf("cluster verification failed: %w", err)
	}

	clusterInfo := &ClusterInfo{
		Name:      clusterName,
		Namespace: namespace,
		Context:   kubeContext,
		Cleanup:   cleanupFunc,
	}

	t.Logf("Test cluster created: %s (context: %s)", clusterName, kubeContext)
	return clusterInfo, nil
}

// CleanupTestCluster safely destroys a test cluster
func CleanupTestCluster(t *testing.T, cluster *ClusterInfo) {
	if cluster == nil {
		t.Log("No cluster to cleanup")
		return
	}

	t.Logf("Cleaning up test cluster: %s", cluster.Name)

	// Execute cleanup function if available
	if cluster.Cleanup != nil {
		cluster.Cleanup()
	} else {
		// Fallback cleanup using idpbuilder
		cmd := exec.Command("idpbuilder", "delete", "--name", cluster.Name)
		if err := cmd.Run(); err != nil {
			t.Logf("Warning: Failed to cleanup cluster %s: %v", cluster.Name, err)
		}
	}

	t.Logf("Cluster cleanup completed: %s", cluster.Name)
}

// GenerateTestCredentials creates test credentials for registry authentication
func GenerateTestCredentials() (*Credentials, error) {
	// Generate deterministic but unique credentials for testing
	timestamp := time.Now().Unix()

	creds := &Credentials{
		Username: fmt.Sprintf("testuser-%d", timestamp),
		Password: fmt.Sprintf("testpass-%d", timestamp),
		Token:    "", // Will be populated if token auth is needed
	}

	return creds, nil
}

// verifyClusterAccess checks if the cluster is accessible via kubectl
func verifyClusterAccess(kubeContext string) error {
	// Check if kubectl is available
	if _, err := exec.LookPath("kubectl"); err != nil {
		return fmt.Errorf("kubectl not found: %w", err)
	}

	// Try to access cluster info
	cmd := exec.Command("kubectl", "--context", kubeContext, "cluster-info")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("cannot access cluster %s: %w\nOutput: %s", kubeContext, err, string(output))
	}

	return nil
}

// getKubeconfig returns the path to the kubeconfig file for the cluster
func getKubeconfig(clusterName string) (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("cannot get home directory: %w", err)
	}
	return filepath.Join(homeDir, ".kube", "config"), nil
}

// validateCredentials checks if the provided credentials are valid
func validateCredentials(creds *Credentials) error {
	if creds == nil {
		return fmt.Errorf("credentials cannot be nil")
	}
	if strings.TrimSpace(creds.Username) == "" {
		return fmt.Errorf("username cannot be empty")
	}
	if strings.TrimSpace(creds.Password) == "" && strings.TrimSpace(creds.Token) == "" {
		return fmt.Errorf("either password or token must be provided")
	}
	return nil
}