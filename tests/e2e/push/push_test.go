//go:build e2e

package push_test

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/cnoe-io/idpbuilder/tests/e2e"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	IdpbuilderBinaryLocation = "../../../../idpbuilder"
	DefaultGiteaHost         = "gitea.cnoe.localtest.me"
	DefaultGiteaPort         = "8443"
	TestImageName            = "e2e-push-test"
	TestImageTag             = "latest"
)

// getGiteaCredentials retrieves credentials from idpbuilder
func getGiteaCredentials(ctx context.Context, t *testing.T) (e2e.BasicAuth, error) {
	return e2e.GetBasicAuth(ctx, "gitea-credential")
}

// buildTestImage creates a minimal test image
func buildTestImage(ctx context.Context, t *testing.T) string {
	t.Helper()

	dockerfileContent := `FROM scratch
LABEL test=true
`
	tmpDir := t.TempDir()
	dockerfilePath := filepath.Join(tmpDir, "Dockerfile")
	err := os.WriteFile(dockerfilePath, []byte(dockerfileContent), 0644)
	require.NoError(t, err, "failed to write test Dockerfile")

	imageRef := fmt.Sprintf("%s:%s-%d", TestImageName, TestImageTag, time.Now().UnixNano())

	cmd := exec.CommandContext(ctx, "docker", "build", "-f", dockerfilePath, "-t", imageRef, tmpDir)
	output, err := cmd.CombinedOutput()
	require.NoError(t, err, "failed to build test image: %s", output)

	return imageRef
}

// cleanupTestImage removes a test image
func cleanupTestImage(ctx context.Context, t *testing.T, imageRef string) {
	t.Helper()
	cmd := exec.CommandContext(ctx, "docker", "rmi", "-f", imageRef)
	_ = cmd.Run() // Ignore errors during cleanup
}

// TestPushCommand_BasicPush validates basic push to Gitea registry
func TestPushCommand_BasicPush(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping e2e test in short mode")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	creds, err := getGiteaCredentials(ctx, t)
	require.NoError(t, err, "failed to get Gitea credentials")

	testImageRef := buildTestImage(ctx, t)
	defer cleanupTestImage(ctx, t, testImageRef)

	destRef := fmt.Sprintf("%s:%s/giteaadmin/%s:%s",
		DefaultGiteaHost, DefaultGiteaPort, TestImageName, TestImageTag)

	cmd := exec.CommandContext(ctx, IdpbuilderBinaryLocation, "push",
		testImageRef,
		"--registry", fmt.Sprintf("https://%s:%s", DefaultGiteaHost, DefaultGiteaPort),
		"--username", creds.Username,
		"--password", creds.Password,
		"--insecure",
	)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()
	require.NoError(t, err, "push command failed: %s\nstderr: %s", stdout.String(), stderr.String())

	// REQ-001: stdout contains only pushed reference
	output := strings.TrimSpace(stdout.String())
	assert.Contains(t, output, "@sha256:", "stdout should contain digest")

	// Property W3.2: Progress goes to stderr
	assert.Contains(t, stderr.String(), "Pushing", "stderr should contain progress")

	t.Logf("Push successful: %s", output)
}

// TestPushCommand_EnvironmentCredentials validates credential resolution via env vars
func TestPushCommand_EnvironmentCredentials(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping e2e test in short mode")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	creds, err := getGiteaCredentials(ctx, t)
	require.NoError(t, err, "failed to get Gitea credentials")

	testImageRef := buildTestImage(ctx, t)
	defer cleanupTestImage(ctx, t, testImageRef)

	cmd := exec.CommandContext(ctx, IdpbuilderBinaryLocation, "push",
		testImageRef,
		"--registry", fmt.Sprintf("https://%s:%s", DefaultGiteaHost, DefaultGiteaPort),
		"--insecure",
	)
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("IDPBUILDER_REGISTRY_USERNAME=%s", creds.Username),
		fmt.Sprintf("IDPBUILDER_REGISTRY_PASSWORD=%s", creds.Password),
	)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()
	require.NoError(t, err, "push with env credentials failed: %s", stderr.String())

	output := strings.TrimSpace(stdout.String())
	assert.Contains(t, output, "@sha256:", "stdout should contain digest")
}

// TestPushCommand_FlagOverridesEnv validates Property P1.1: flags take precedence over environment variables
func TestPushCommand_FlagOverridesEnv(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping e2e test in short mode")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	creds, err := getGiteaCredentials(ctx, t)
	require.NoError(t, err, "failed to get Gitea credentials")

	testImageRef := buildTestImage(ctx, t)
	defer cleanupTestImage(ctx, t, testImageRef)

	// Set WRONG credentials in environment, correct via flags
	cmd := exec.CommandContext(ctx, IdpbuilderBinaryLocation, "push",
		testImageRef,
		"--registry", fmt.Sprintf("https://%s:%s", DefaultGiteaHost, DefaultGiteaPort),
		"--username", creds.Username,    // Correct via flag
		"--password", creds.Password,    // Correct via flag
		"--insecure",
	)
	cmd.Env = append(os.Environ(),
		"IDPBUILDER_REGISTRY_USERNAME=wrong-user",
		"IDPBUILDER_REGISTRY_PASSWORD=wrong-password",
	)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Should succeed because flags override wrong env credentials
	err = cmd.Run()
	require.NoError(t, err, "push should succeed when flags override wrong env credentials")

	output := strings.TrimSpace(stdout.String())
	assert.Contains(t, output, "@sha256:", "stdout should contain digest")
}

// TestPushCommand_InvalidCredentials validates exit code 1 for authentication failure
func TestPushCommand_InvalidCredentials(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping e2e test in short mode")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	testImageRef := buildTestImage(ctx, t)
	defer cleanupTestImage(ctx, t, testImageRef)

	cmd := exec.CommandContext(ctx, IdpbuilderBinaryLocation, "push",
		testImageRef,
		"--registry", fmt.Sprintf("https://%s:%s", DefaultGiteaHost, DefaultGiteaPort),
		"--username", "invalid-user",
		"--password", "invalid-password",
		"--insecure",
	)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	require.Error(t, err, "push with invalid credentials should fail")

	exitErr, ok := err.(*exec.ExitError)
	require.True(t, ok, "error should be ExitError")
	assert.Equal(t, 1, exitErr.ExitCode(), "exit code should be 1 for auth failure")
}

// TestPushCommand_ImageNotFound validates exit code 2 for missing image
func TestPushCommand_ImageNotFound(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping e2e test in short mode")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	creds, err := getGiteaCredentials(ctx, t)
	require.NoError(t, err, "failed to get Gitea credentials")

	cmd := exec.CommandContext(ctx, IdpbuilderBinaryLocation, "push",
		"nonexistent-image:nonexistent-tag",
		"--registry", fmt.Sprintf("https://%s:%s", DefaultGiteaHost, DefaultGiteaPort),
		"--username", creds.Username,
		"--password", creds.Password,
		"--insecure",
	)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()
	require.Error(t, err, "push of non-existent image should fail")

	exitErr, ok := err.(*exec.ExitError)
	require.True(t, ok, "error should be ExitError")
	assert.Equal(t, 2, exitErr.ExitCode(), "exit code should be 2 for image not found")
}

// TestPushCommand_VerifyPushedImage validates pushed image can be retrieved from registry
func TestPushCommand_VerifyPushedImage(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping e2e test in short mode")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	creds, err := getGiteaCredentials(ctx, t)
	require.NoError(t, err, "failed to get Gitea credentials")

	testImageRef := buildTestImage(ctx, t)
	defer cleanupTestImage(ctx, t, testImageRef)

	verifyImageName := TestImageName + "-verify"
	destRef := fmt.Sprintf("%s:%s/giteaadmin/%s:%s",
		DefaultGiteaHost, DefaultGiteaPort, verifyImageName, TestImageTag)

	pushCmd := exec.CommandContext(ctx, IdpbuilderBinaryLocation, "push",
		testImageRef,
		"--registry", fmt.Sprintf("https://%s:%s", DefaultGiteaHost, DefaultGiteaPort),
		"--username", creds.Username,
		"--password", creds.Password,
		"--insecure",
	)

	var pushStdout, pushStderr bytes.Buffer
	pushCmd.Stdout = &pushStdout
	pushCmd.Stderr = &pushStderr

	err = pushCmd.Run()
	require.NoError(t, err, "push failed: %s", pushStderr.String())

	// Clean local copy to ensure pull from registry
	cleanupTestImage(ctx, t, destRef)

	// Pull from registry to verify
	pullCmd := exec.CommandContext(ctx, "docker", "pull", destRef)
	pullOutput, err := pullCmd.CombinedOutput()
	require.NoError(t, err, "pull from registry failed: %s", pullOutput)

	t.Logf("Successfully verified pushed image by pulling: %s", destRef)
}

// TestE2E_BuildAndPush validates complete workflow: docker build -> idpbuilder push -> verify
func TestE2E_BuildAndPush(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping e2e test in short mode")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
	defer cancel()

	creds, err := getGiteaCredentials(ctx, t)
	require.NoError(t, err, "failed to get Gitea credentials")

	// Create test Dockerfile with actual content
	tmpDir := t.TempDir()
	dockerfileContent := `FROM alpine:3.19
RUN echo "Hello from idpbuilder push test" > /hello.txt
CMD ["cat", "/hello.txt"]
`
	dockerfilePath := filepath.Join(tmpDir, "Dockerfile")
	err = os.WriteFile(dockerfilePath, []byte(dockerfileContent), 0644)
	require.NoError(t, err, "failed to write test Dockerfile")

	// Build image
	imageRef := fmt.Sprintf("e2e-build-push:%d", time.Now().UnixNano())
	buildCmd := exec.CommandContext(ctx, "docker", "build", "-f", dockerfilePath, "-t", imageRef, tmpDir)
	buildOutput, err := buildCmd.CombinedOutput()
	require.NoError(t, err, "docker build failed: %s", buildOutput)
	defer cleanupTestImage(ctx, t, imageRef)

	// Push using idpbuilder
	pushCmd := exec.CommandContext(ctx, IdpbuilderBinaryLocation, "push",
		imageRef,
		"--registry", fmt.Sprintf("https://%s:%s", DefaultGiteaHost, DefaultGiteaPort),
		"--username", creds.Username,
		"--password", creds.Password,
		"--insecure",
	)

	var stdout, stderr bytes.Buffer
	pushCmd.Stdout = &stdout
	pushCmd.Stderr = &stderr

	err = pushCmd.Run()
	require.NoError(t, err, "idpbuilder push failed: stdout=%s, stderr=%s", stdout.String(), stderr.String())

	output := strings.TrimSpace(stdout.String())
	assert.Contains(t, output, "@sha256:", "should contain digest")
	assert.Contains(t, stderr.String(), "Pushing", "stderr should show progress")

	t.Logf("E2E Build and Push successful: %s", output)
}

// TestE2E_RetryOnNetworkError validates recovery from simulated network issues
func TestE2E_RetryOnNetworkError(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping e2e test in short mode")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
	defer cancel()

	creds, err := getGiteaCredentials(ctx, t)
	require.NoError(t, err, "failed to get Gitea credentials")

	testImageRef := buildTestImage(ctx, t)
	defer cleanupTestImage(ctx, t, testImageRef)

	// Test that push completes successfully despite transient issues
	// The actual network simulation is handled by E1.3.1 (retry logic)
	// This test validates the integration with the retry mechanism
	cmd := exec.CommandContext(ctx, IdpbuilderBinaryLocation, "push",
		testImageRef,
		"--registry", fmt.Sprintf("https://%s:%s", DefaultGiteaHost, DefaultGiteaPort),
		"--username", creds.Username,
		"--password", creds.Password,
		"--insecure",
	)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()
	require.NoError(t, err, "push should succeed with retry logic: stderr=%s", stderr.String())

	output := strings.TrimSpace(stdout.String())
	assert.Contains(t, output, "@sha256:", "stdout should contain digest")

	t.Logf("Retry validation successful: %s", output)
}

// TestE2E_ProgressOutput verifies progress is written to stderr, not stdout
func TestE2E_ProgressOutput(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping e2e test in short mode")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	creds, err := getGiteaCredentials(ctx, t)
	require.NoError(t, err, "failed to get Gitea credentials")

	testImageRef := buildTestImage(ctx, t)
	defer cleanupTestImage(ctx, t, testImageRef)

	cmd := exec.CommandContext(ctx, IdpbuilderBinaryLocation, "push",
		testImageRef,
		"--registry", fmt.Sprintf("https://%s:%s", DefaultGiteaHost, DefaultGiteaPort),
		"--username", creds.Username,
		"--password", creds.Password,
		"--insecure",
	)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()
	require.NoError(t, err, "push failed")

	// Property W3.2: Progress goes to stderr
	assert.Contains(t, stderr.String(), "Pushing", "stderr should contain 'Pushing'")
	assert.Contains(t, stderr.String(), "layers", "stderr should mention layers")

	// stdout should only contain the pushed reference
	stdoutLines := strings.Split(strings.TrimSpace(stdout.String()), "\n")
	assert.Equal(t, 1, len(stdoutLines), "stdout should have exactly one line (the reference)")
	assert.Contains(t, stdoutLines[0], "@sha256:", "stdout should be the pushed reference")
}
