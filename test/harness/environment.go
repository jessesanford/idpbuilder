package harness

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/docker/docker/client"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// TestEnvironment provides complete integration test infrastructure
// with Gitea registry container and Docker client for testing
type TestEnvironment struct {
	GiteaContainer testcontainers.Container
	RegistryURL    string
	AdminUsername  string
	AdminPassword  string
	DockerClient   *client.Client
	CleanupFunc    func() error
}

// SetupGiteaRegistry starts a registry container (Docker registry:2)
// and configures it for integration testing. Returns a TestEnvironment
// with all necessary resources initialized.
//
// Note: Despite the name "Gitea" for backward compatibility, this now uses
// the official Docker registry:2 image which is simpler and more reliable
// for integration testing of OCI push functionality.
func SetupGiteaRegistry(ctx context.Context) (*TestEnvironment, error) {
	// Use official Docker registry:2 for simpler, more reliable testing
	req := testcontainers.ContainerRequest{
		Image:        "registry:2",
		ExposedPorts: []string{"5000/tcp"},
		Env: map[string]string{
			// Disable authentication for test simplicity
			"REGISTRY_AUTH": "none",
			// Enable deletion for cleanup
			"REGISTRY_STORAGE_DELETE_ENABLED": "true",
		},
		// Wait for HTTP health check on registry port
		WaitingFor: wait.ForHTTP("/v2/").
			WithPort("5000/tcp").
			WithStartupTimeout(60 * time.Second).
			WithPollInterval(500 * time.Millisecond),
	}

	// Start the registry container
	registryContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to start registry container: %w", err)
	}

	// Get dynamically allocated port for registry (5000)
	registryPort, err := registryContainer.MappedPort(ctx, "5000")
	if err != nil {
		registryContainer.Terminate(ctx)
		return nil, fmt.Errorf("failed to get registry port: %w", err)
	}

	// Construct registry URL with dynamic port
	registryURL := fmt.Sprintf("localhost:%s", registryPort.Port())

	// No authentication required for registry:2 in test mode
	// But provide dummy credentials for code that requires them
	// The registry:2 container will ignore these and allow anonymous access
	adminUsername := "testuser"
	adminPassword := "testpass"

	// Create Docker client for test operations
	dockerClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		registryContainer.Terminate(ctx)
		return nil, fmt.Errorf("failed to create Docker client: %w", err)
	}

	// Create cleanup function
	cleanupFunc := func() error {
		var errs []error

		// Close Docker client
		if dockerClient != nil {
			if err := dockerClient.Close(); err != nil {
				errs = append(errs, fmt.Errorf("failed to close Docker client: %w", err))
			}
		}

		// Terminate registry container
		if registryContainer != nil {
			if err := registryContainer.Terminate(ctx); err != nil {
				errs = append(errs, fmt.Errorf("failed to terminate container: %w", err))
			}
		}

		if len(errs) > 0 {
			return fmt.Errorf("cleanup errors: %v", errs)
		}
		return nil
	}

	env := &TestEnvironment{
		GiteaContainer: registryContainer,
		RegistryURL:    registryURL,
		AdminUsername:  adminUsername,
		AdminPassword:  adminPassword,
		DockerClient:   dockerClient,
		CleanupFunc:    cleanupFunc,
	}

	return env, nil
}

// createGiteaUser creates a new user in Gitea via the API
func createGiteaUser(webURL, username, password, email string) error {
	// Gitea user creation payload
	payload := map[string]interface{}{
		"username":             username,
		"email":                email,
		"password":             password,
		"must_change_password": false,
		"send_notify":          false,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	// Use user/sign_up endpoint (doesn't require authentication for first user)
	signupURL := fmt.Sprintf("%s/user/sign_up", webURL)
	req, err := http.NewRequest("POST", signupURL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Accept both 201 (created) and 302/303 (redirect after creation)
	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusFound &&
		resp.StatusCode != http.StatusSeeOther && resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

// initializeGiteaAdmin would create an admin user inside the Gitea container
// However, this requires running commands as non-root which conflicts with
// testcontainers execution model. In practice, tests should create users
// via Gitea's API after the container starts.
//
// Example API-based user creation (for future reference):
//   POST /api/v1/admin/users with JSON body
//
// Keeping this as a comment for documentation purposes.

// Cleanup cleans up the test environment by closing connections
// and terminating containers. This should be called in a defer
// statement after SetupGiteaRegistry.
func (env *TestEnvironment) Cleanup() error {
	if env.CleanupFunc != nil {
		return env.CleanupFunc()
	}
	return nil
}
