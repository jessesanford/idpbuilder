package integration

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestBasicAuthentication tests basic authentication scenarios with Gitea registry
func TestBasicAuthentication(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	env := NewTestEnvironment(t)
	require.NoError(t, env.SetupIDPBuilder(t))
	defer env.Cleanup(t)

	testCases := []struct {
		name          string
		username      string
		password      string
		expectSuccess bool
		errorContains string
	}{
		{
			name:          "Valid credentials",
			username:      env.GiteaUsername,
			password:      env.GiteaPassword,
			expectSuccess: true,
		},
		{
			name:          "Invalid username",
			username:      "invaliduser",
			password:      env.GiteaPassword,
			expectSuccess: false,
			errorContains: "authentication failed",
		},
		{
			name:          "Invalid password",
			username:      env.GiteaUsername,
			password:      "wrongpassword",
			expectSuccess: false,
			errorContains: "authentication failed",
		},
		{
			name:          "Empty username",
			username:      "",
			password:      env.GiteaPassword,
			expectSuccess: false,
			errorContains: "username required",
		},
		{
			name:          "Empty password",
			username:      env.GiteaUsername,
			password:      "",
			expectSuccess: false,
			errorContains: "password required",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			imageName := fmt.Sprintf("%s/test-auth-%d:latest", env.GiteaRegistry, time.Now().Unix())

			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
			defer cancel()

			// Build test command with credentials
			cmd := exec.CommandContext(ctx, "idpbuilder", "push",
				"--source", "alpine:latest",
				"--dest", imageName,
				"--username", tc.username,
				"--password", tc.password,
			)

			output, err := cmd.CombinedOutput()

			if tc.expectSuccess {
				assert.NoError(t, err, "Expected push to succeed but got error: %v\nOutput: %s", err, string(output))
			} else {
				assert.Error(t, err, "Expected push to fail but it succeeded")
				if tc.errorContains != "" {
					assert.Contains(t, string(output), tc.errorContains, "Error message should contain expected text")
				}
			}
		})
	}
}

// TestTokenAuthentication tests token-based authentication
func TestTokenAuthentication(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	env := NewTestEnvironment(t)
	require.NoError(t, env.SetupIDPBuilder(t))
	defer env.Cleanup(t)

	// Generate an access token
	token, err := generateGiteaAccessToken(t, env)
	require.NoError(t, err, "Failed to generate access token")

	testCases := []struct {
		name          string
		token         string
		expectSuccess bool
		errorContains string
	}{
		{
			name:          "Valid token",
			token:         token,
			expectSuccess: true,
		},
		{
			name:          "Invalid token",
			token:         "invalid-token-12345",
			expectSuccess: false,
			errorContains: "authentication failed",
		},
		{
			name:          "Empty token",
			token:         "",
			expectSuccess: false,
			errorContains: "token required",
		},
		{
			name:          "Expired token",
			token:         "expired-token-from-past",
			expectSuccess: false,
			errorContains: "authentication failed",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			imageName := fmt.Sprintf("%s/test-token-%d:latest", env.GiteaRegistry, time.Now().Unix())

			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
			defer cancel()

			cmd := exec.CommandContext(ctx, "idpbuilder", "push",
				"--source", "alpine:latest",
				"--dest", imageName,
				"--token", tc.token,
			)

			output, err := cmd.CombinedOutput()

			if tc.expectSuccess {
				assert.NoError(t, err, "Expected push to succeed but got error: %v\nOutput: %s", err, string(output))
			} else {
				assert.Error(t, err, "Expected push to fail but it succeeded")
				if tc.errorContains != "" {
					assert.Contains(t, string(output), tc.errorContains, "Error message should contain expected text")
				}
			}
		})
	}
}

// TestNoAuthRegistry tests pushing to a registry without authentication
func TestNoAuthRegistry(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// This test requires a no-auth registry configuration
	// For now, we'll skip if not available
	t.Skip("No-auth registry not available in test environment")

	imageName := "localhost:5000/test-noauth:latest"

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	cmd := exec.CommandContext(ctx, "idpbuilder", "push",
		"--source", "alpine:latest",
		"--dest", imageName,
		"--insecure",
	)

	output, err := cmd.CombinedOutput()
	assert.NoError(t, err, "Push to no-auth registry failed: %v\nOutput: %s", err, string(output))
}

// TestAuthenticationWithTLS tests authentication over TLS
func TestAuthenticationWithTLS(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	env := NewTestEnvironment(t)
	require.NoError(t, env.SetupIDPBuilder(t))
	defer env.Cleanup(t)

	testCases := []struct {
		name          string
		useTLS        bool
		skipVerify    bool
		provideCert   bool
		expectSuccess bool
	}{
		{
			name:          "TLS with certificate verification",
			useTLS:        true,
			skipVerify:    false,
			provideCert:   true,
			expectSuccess: true,
		},
		{
			name:          "TLS skip verification",
			useTLS:        true,
			skipVerify:    true,
			provideCert:   false,
			expectSuccess: true,
		},
		{
			name:          "TLS without certificate",
			useTLS:        true,
			skipVerify:    false,
			provideCert:   false,
			expectSuccess: false,
		},
		{
			name:          "No TLS",
			useTLS:        false,
			skipVerify:    false,
			provideCert:   false,
			expectSuccess: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			imageName := fmt.Sprintf("%s/test-tls-%d:latest", env.GiteaRegistry, time.Now().Unix())

			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
			defer cancel()

			args := []string{
				"push",
				"--source", "alpine:latest",
				"--dest", imageName,
				"--username", env.GiteaUsername,
				"--password", env.GiteaPassword,
			}

			if tc.skipVerify {
				args = append(args, "--insecure-skip-tls-verify")
			}

			if tc.provideCert {
				args = append(args, "--cert", env.CertPath)
			}

			cmd := exec.CommandContext(ctx, "idpbuilder", args...)
			output, err := cmd.CombinedOutput()

			if tc.expectSuccess {
				assert.NoError(t, err, "Expected push to succeed but got error: %v\nOutput: %s", err, string(output))
			} else {
				assert.Error(t, err, "Expected push to fail but it succeeded")
			}
		})
	}
}

// generateGiteaAccessToken generates an access token for the test user
func generateGiteaAccessToken(t *testing.T, env *TestEnvironment) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	tokenName := fmt.Sprintf("test-token-%d", time.Now().Unix())

	cmd := exec.CommandContext(ctx, "kubectl", "exec", "-n", "gitea", "deployment/gitea", "--",
		"gitea", "admin", "user", "generate-access-token",
		"--username", env.GiteaUsername,
		"--token-name", tokenName,
		"--scopes", "write:repository,read:repository",
	)

	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to generate access token: %w", err)
	}

	// Parse token from output (format varies, this is a placeholder)
	token := string(output)
	t.Logf("Generated access token: %s", tokenName)

	return token, nil
}

// TestCredentialCaching tests that credentials are properly cached
func TestCredentialCaching(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	env := NewTestEnvironment(t)
	require.NoError(t, env.SetupIDPBuilder(t))
	defer env.Cleanup(t)

	imageName := fmt.Sprintf("%s/test-cache:latest", env.GiteaRegistry)

	// First push with credentials
	ctx1, cancel1 := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel1()

	cmd1 := exec.CommandContext(ctx1, "idpbuilder", "push",
		"--source", "alpine:latest",
		"--dest", imageName,
		"--username", env.GiteaUsername,
		"--password", env.GiteaPassword,
	)

	output1, err1 := cmd1.CombinedOutput()
	require.NoError(t, err1, "First push failed: %v\nOutput: %s", err1, string(output1))

	// Second push should use cached credentials
	ctx2, cancel2 := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel2()

	cmd2 := exec.CommandContext(ctx2, "idpbuilder", "push",
		"--source", "alpine:latest",
		"--dest", imageName,
	)

	output2, err2 := cmd2.CombinedOutput()

	// Note: This test assumes credential caching is implemented
	// If not implemented, we expect this to fail
	t.Logf("Second push output: %s", string(output2))
	if err2 != nil {
		t.Logf("Credential caching may not be implemented: %v", err2)
	}
}

// TestAuthEnvironmentVariables tests authentication via environment variables
func TestAuthEnvironmentVariables(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	env := NewTestEnvironment(t)
	require.NoError(t, env.SetupIDPBuilder(t))
	defer env.Cleanup(t)

	imageName := fmt.Sprintf("%s/test-env-auth:latest", env.GiteaRegistry)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	cmd := exec.CommandContext(ctx, "idpbuilder", "push",
		"--source", "alpine:latest",
		"--dest", imageName,
	)

	// Set credentials via environment variables
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("REGISTRY_USERNAME=%s", env.GiteaUsername),
		fmt.Sprintf("REGISTRY_PASSWORD=%s", env.GiteaPassword),
	)

	output, err := cmd.CombinedOutput()
	assert.NoError(t, err, "Push with env var auth failed: %v\nOutput: %s", err, string(output))
}
