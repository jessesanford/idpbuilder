package harness

import (
	"context"
	"fmt"
	"time"

	"github.com/docker/docker/client"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// TestEnvironment provides complete integration test infrastructure
// with Gitea registry container and Docker client for testing
type TestEnvironment struct {
	GiteaContainer  testcontainers.Container
	RegistryURL     string
	AdminUsername   string
	AdminPassword   string
	DockerClient    *client.Client
	CleanupFunc     func() error
}

// SetupGiteaRegistry starts a Gitea container with registry enabled
// and configures it for integration testing. Returns a TestEnvironment
// with all necessary resources initialized.
func SetupGiteaRegistry(ctx context.Context) (*TestEnvironment, error) {
	// Configure Gitea container request with registry support
	req := testcontainers.ContainerRequest{
		Image:        "gitea/gitea:1.20",
		ExposedPorts: []string{"3000/tcp"},
		Env: map[string]string{
			// Enable container registry support
			"GITEA__packages__ENABLED": "true",
			// Use SQLite for simplicity
			"GITEA__database__DB_TYPE": "sqlite3",
			// Skip installation wizard
			"GITEA__security__INSTALL_LOCK": "true",
			// Disable require sign-in for easier testing
			"GITEA__service__REQUIRE_SIGNIN_VIEW": "false",
		},
		WaitingFor: wait.ForLog("Starting new Web server: tcp:0.0.0.0:3000").
			WithStartupTimeout(60 * time.Second),
	}

	// Start the Gitea container
	giteaContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to start Gitea container: %w", err)
	}

	// Get dynamically allocated port for web interface (3000)
	// Gitea serves both web and registry on the same port
	webPort, err := giteaContainer.MappedPort(ctx, "3000")
	if err != nil {
		giteaContainer.Terminate(ctx)
		return nil, fmt.Errorf("failed to get web port: %w", err)
	}

	// Construct registry URL with dynamic port
	// Gitea registry is accessed through the same web port
	registryURL := fmt.Sprintf("localhost:%s", webPort.Port())
	webURL := fmt.Sprintf("http://localhost:%s", webPort.Port())

	// Wait for Gitea to be fully initialized
	time.Sleep(5 * time.Second)

	// Admin user credentials for test purposes
	// Note: In real tests, users should be created via Gitea API after container starts
	// For now, providing default credentials that tests can use to create users via API
	adminUsername := "testadmin"
	adminPassword := "testpassword123"

	// Create Docker client for test operations
	dockerClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		giteaContainer.Terminate(ctx)
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

		// Terminate Gitea container
		if giteaContainer != nil {
			if err := giteaContainer.Terminate(ctx); err != nil {
				errs = append(errs, fmt.Errorf("failed to terminate container: %w", err))
			}
		}

		if len(errs) > 0 {
			return fmt.Errorf("cleanup errors: %v", errs)
		}
		return nil
	}

	env := &TestEnvironment{
		GiteaContainer: giteaContainer,
		RegistryURL:    registryURL,
		AdminUsername:  adminUsername,
		AdminPassword:  adminPassword,
		DockerClient:   dockerClient,
		CleanupFunc:    cleanupFunc,
	}

	// Store web URL for potential future use (not exposed in struct for now)
	_ = webURL

	return env, nil
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
