package extractor

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

func TestGetDefaultKubeconfig(t *testing.T) {
	originalKubeconfig := os.Getenv("KUBECONFIG")
	defer os.Setenv("KUBECONFIG", originalKubeconfig)

	testPath := "/test/kubeconfig"
	os.Setenv("KUBECONFIG", testPath)
	result := getDefaultKubeconfig()
	assert.Equal(t, testPath, result)

	os.Unsetenv("KUBECONFIG")
	result = getDefaultKubeconfig()
	assert.Contains(t, result, ".kube/config")
}

func TestValidateKindContext(t *testing.T) {
	// Create a temporary kubeconfig file for testing
	tempDir := t.TempDir()
	kubeconfigPath := filepath.Join(tempDir, "config")

	// Create a test kubeconfig with Kind context
	config := &clientcmdapi.Config{
		Contexts: map[string]*clientcmdapi.Context{
			"kind-test-cluster": {
				Cluster:  "kind-test-cluster",
				AuthInfo: "kind-test-cluster",
			},
		},
		CurrentContext: "kind-test-cluster",
		Clusters: map[string]*clientcmdapi.Cluster{
			"kind-test-cluster": {
				Server: "https://127.0.0.1:12345",
			},
		},
		AuthInfos: map[string]*clientcmdapi.AuthInfo{
			"kind-test-cluster": {},
		},
	}

	err := clientcmd.WriteToFile(*config, kubeconfigPath)
	require.NoError(t, err)

	// Test with valid Kind context
	err = validateKindContext(kubeconfigPath, "")
	assert.NoError(t, err)

	// Test with expected cluster name matching
	err = validateKindContext(kubeconfigPath, "test-cluster")
	assert.NoError(t, err)

	// Test with expected cluster name not matching
	err = validateKindContext(kubeconfigPath, "other-cluster")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "doesn't match expected cluster")

	// Create kubeconfig with non-Kind context
	config.CurrentContext = "minikube"
	config.Contexts["minikube"] = &clientcmdapi.Context{
		Cluster:  "minikube",
		AuthInfo: "minikube",
	}
	err = clientcmd.WriteToFile(*config, kubeconfigPath)
	require.NoError(t, err)

	// Test with non-Kind context
	err = validateKindContext(kubeconfigPath, "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "doesn't appear to be a Kind cluster")
}

func TestSetupKubernetesClient_InvalidPath(t *testing.T) {
	// Test with non-existent kubeconfig
	client, path, err := setupKubernetesClient("/non/existent/path")
	assert.Error(t, err)
	assert.Nil(t, client)
	assert.Empty(t, path)
	assert.Contains(t, err.Error(), "kubeconfig not found")
}

func TestSetupKubernetesClient_InvalidKubeconfig(t *testing.T) {
	// Create a temporary file with invalid content
	tempDir := t.TempDir()
	kubeconfigPath := filepath.Join(tempDir, "invalid-config")

	err := os.WriteFile(kubeconfigPath, []byte("invalid yaml content"), 0600)
	require.NoError(t, err)

	// Test with invalid kubeconfig
	client, path, err := setupKubernetesClient(kubeconfigPath)
	assert.Error(t, err)
	assert.Nil(t, client)
	assert.Empty(t, path)
}