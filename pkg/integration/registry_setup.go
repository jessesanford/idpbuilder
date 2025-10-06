//go:build integration

package integration

import (
	"context"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// ClusterInfo holds information about a test cluster
type ClusterInfo struct {
	Name      string     // Cluster name
	Namespace string     // Target namespace
	Context   string     // Kubeconfig context
	Cleanup   func()     // Cleanup function
}

// SetupTestRegistry starts a Gitea container with registry enabled for testing
func SetupTestRegistry(t *testing.T) (*testcontainers.Container, string) {
	ctx := context.Background()

	// Configure Gitea container with registry support
	req := testcontainers.ContainerRequest{
		Image:        "gitea/gitea:latest",
		ExposedPorts: []string{"3000/tcp"},
		Env: map[string]string{
			"GITEA__server__ROOT_URL": "http://localhost:3000",
			"GITEA__server__DISABLE_SSH": "true",
			"GITEA__log__LEVEL": "warn",
			// Enable container registry
			"GITEA__packages__ENABLED": "true",
			"GITEA__packages__CHUNKED_UPLOAD_PATH": "/tmp/gitea-packages",
		},
		WaitingFor: wait.ForAll(
			wait.ForHTTP("/").WithPort("3000").WithStartupTimeout(60*time.Second),
			wait.ForLog("Listen:").WithStartupTimeout(60*time.Second),
		),
		Networks: []string{"bridge"},
	}

	// Start the container
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:         true,
	})
	if err != nil {
		t.Fatalf("Failed to start Gitea container: %v", err)
	}

	// Get mapped port
	mappedPort, err := container.MappedPort(ctx, "3000")
	if err != nil {
		t.Fatalf("Failed to get mapped port: %v", err)
	}

	// Get container host
	host, err := container.Host(ctx)
	if err != nil {
		t.Fatalf("Failed to get container host: %v", err)
	}

	// Build registry URL
	registryURL := fmt.Sprintf("%s:%s", host, mappedPort.Port())

	// Wait for registry to be fully ready
	if err := waitForRegistryReady(registryURL, 30*time.Second); err != nil {
		// Clean up on failure
		if termErr := container.Terminate(ctx); termErr != nil {
			t.Logf("Failed to cleanup container after registry setup failure: %v", termErr)
		}
		t.Fatalf("Registry not ready: %v", err)
	}

	t.Logf("Test registry started at: %s", registryURL)
	return &container, registryURL
}

// waitForRegistryReady polls the registry endpoint until it's responding
func waitForRegistryReady(registryURL string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)

	for time.Now().Before(deadline) {
		conn, err := net.DialTimeout("tcp", registryURL, 2*time.Second)
		if err == nil {
			conn.Close()
			return nil
		}
		time.Sleep(500 * time.Millisecond)
	}

	return fmt.Errorf("registry at %s not ready within %v", registryURL, timeout)
}

// setupRegistryCredentials configures test credentials for the registry
func setupRegistryCredentials(registryURL string) (*Credentials, error) {
	return &Credentials{
		Username: "testuser",
		Password: "testpass123",
		Token:    "",
	}, nil
}

// validateRegistryConnection tests basic connectivity to the registry
func validateRegistryConnection(registryURL string) error {
	conn, err := net.DialTimeout("tcp", registryURL, 5*time.Second)
	if err != nil {
		return fmt.Errorf("cannot connect to registry at %s: %w", registryURL, err)
	}
	defer conn.Close()
	return nil
}