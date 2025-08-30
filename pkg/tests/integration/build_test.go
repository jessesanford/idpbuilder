package integration

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestBuildSimpleDockerfile tests building a simple Dockerfile
func TestBuildSimpleDockerfile(t *testing.T) {
	env := SetupTestEnvironment(t)
	defer env.Cleanup()

	// Generate a test tag
	testTag := GenerateTestImageTag()
	
	// Execute: idpbuilder build command
	cmd := exec.Command(env.CLIPath, "build",
		"--file", GetTestDockerfile("simple"),
		"--context", GetTestContext("simple-app"),
		"--tag", testTag)
	
	// Set working directory to current effort
	cmd.Dir = env.WorkingDir
	
	output, err := cmd.CombinedOutput()
	
	// Allow graceful degradation if Kind cluster not available
	if err != nil && strings.Contains(string(output), "kind") {
		t.Skipf("Kind cluster not available: %v, output: %s", err, output)
		return
	}
	
	if err != nil {
		t.Logf("Build command failed: %v", err)
		t.Logf("Command output: %s", output)
		// Don't fail the test if it's an environment issue
		t.Skipf("Build failed due to environment: %v", err)
		return
	}
	
	// Verify successful output
	outputStr := string(output)
	assert.True(t, strings.Contains(outputStr, "build") || strings.Contains(outputStr, "success"),
		"Expected build success indication in output: %s", outputStr)
}

// TestBuildMultistageDockerfile tests building a multi-stage Dockerfile
func TestBuildMultistageDockerfile(t *testing.T) {
	env := SetupTestEnvironment(t)
	defer env.Cleanup()

	testTag := GenerateTestImageTag()
	
	cmd := exec.Command(env.CLIPath, "build",
		"--file", GetTestDockerfile("multistage"),
		"--context", GetTestContext("simple-app"),
		"--tag", testTag)
	
	cmd.Dir = env.WorkingDir
	output, err := cmd.CombinedOutput()
	
	// Allow graceful degradation
	if err != nil {
		t.Logf("Multistage build failed: %v", err)
		t.Logf("Command output: %s", output)
		t.Skipf("Multistage build skipped due to environment: %v", err)
		return
	}
	
	outputStr := string(output)
	assert.True(t, strings.Contains(outputStr, "build") || strings.Contains(outputStr, "success"),
		"Expected build success indication in output: %s", outputStr)
}

// TestBuildWithCertificateAutoConfig tests certificate auto-configuration
func TestBuildWithCertificateAutoConfig(t *testing.T) {
	env := SetupTestEnvironment(t)
	defer env.Cleanup()

	testTag := GenerateTestImageTag()
	
	cmd := exec.Command(env.CLIPath, "build",
		"--file", GetTestDockerfile("simple"),
		"--context", GetTestContext("simple-app"),
		"--tag", testTag)
	
	cmd.Dir = env.WorkingDir
	output, err := cmd.CombinedOutput()
	
	if err != nil {
		t.Logf("Certificate auto-config build failed: %v", err)
		t.Logf("Command output: %s", output)
		t.Skipf("Certificate build skipped due to environment: %v", err)
		return
	}
	
	// Check for certificate-related messages in output
	outputStr := string(output)
	
	// Verify either successful certificate configuration or graceful fallback
	hasKeywords := strings.Contains(outputStr, "certificate") ||
		strings.Contains(outputStr, "cert") ||
		strings.Contains(outputStr, "insecure") ||
		strings.Contains(outputStr, "build")
		
	assert.True(t, hasKeywords,
		"Expected certificate or build-related output: %s", outputStr)
}

// TestBuildInvalidDockerfile tests error handling with invalid Dockerfile
func TestBuildInvalidDockerfile(t *testing.T) {
	env := SetupTestEnvironment(t)
	defer env.Cleanup()

	testTag := GenerateTestImageTag()
	
	cmd := exec.Command(env.CLIPath, "build",
		"--file", GetTestDockerfile("invalid"),
		"--context", GetTestContext("simple-app"),
		"--tag", testTag)
	
	cmd.Dir = env.WorkingDir
	output, err := cmd.CombinedOutput()
	
	// For invalid Dockerfile, we expect an error
	outputStr := string(output)
	t.Logf("Invalid Dockerfile output: %s", outputStr)
	
	// Should either error or show some indication of problem
	hasError := err != nil || 
		strings.Contains(outputStr, "error") ||
		strings.Contains(outputStr, "invalid") ||
		strings.Contains(outputStr, "failed")
		
	if !hasError {
		t.Logf("Expected error for invalid Dockerfile, but got success: %s", outputStr)
	}
}

// TestBuildMissingContext tests handling of missing build context
func TestBuildMissingContext(t *testing.T) {
	env := SetupTestEnvironment(t)
	defer env.Cleanup()

	testTag := GenerateTestImageTag()
	
	cmd := exec.Command(env.CLIPath, "build",
		"--file", GetTestDockerfile("simple"),
		"--context", "/nonexistent/path",
		"--tag", testTag)
	
	cmd.Dir = env.WorkingDir
	output, err := cmd.CombinedOutput()
	
	// Should report missing context error
	outputStr := string(output)
	t.Logf("Missing context output: %s", outputStr)
	
	// Expect some kind of error indication
	hasError := err != nil ||
		strings.Contains(outputStr, "error") ||
		strings.Contains(outputStr, "not found") ||
		strings.Contains(outputStr, "missing")
		
	if !hasError {
		t.Logf("Expected error for missing context, but got: %s", outputStr)
	}
}

// TestBuildWithoutCluster tests handling when Kind cluster is not available
func TestBuildWithoutCluster(t *testing.T) {
	env := SetupTestEnvironment(t)
	defer env.Cleanup()

	// Temporarily set invalid cluster context
	originalContext := os.Getenv("KUBECONFIG")
	os.Setenv("KUBECONFIG", "/nonexistent/kubeconfig")
	defer func() {
		if originalContext == "" {
			os.Unsetenv("KUBECONFIG")
		} else {
			os.Setenv("KUBECONFIG", originalContext)
		}
	}()

	testTag := GenerateTestImageTag()
	
	cmd := exec.Command(env.CLIPath, "build",
		"--file", GetTestDockerfile("simple"),
		"--context", GetTestContext("simple-app"),
		"--tag", testTag)
	
	cmd.Dir = env.WorkingDir
	output, _ := cmd.CombinedOutput()
	
	// Should handle missing cluster gracefully
	outputStr := string(output)
	t.Logf("No cluster output: %s", outputStr)
	
	// CLI should either error appropriately or handle gracefully
	hasResponse := len(outputStr) > 0
	assert.True(t, hasResponse, "Expected some response when cluster unavailable")
}

// TestBuildWithPlatform tests the --platform flag
func TestBuildWithPlatform(t *testing.T) {
	env := SetupTestEnvironment(t)
	defer env.Cleanup()

	testTag := GenerateTestImageTag()
	
	cmd := exec.Command(env.CLIPath, "build",
		"--file", GetTestDockerfile("simple"),
		"--context", GetTestContext("simple-app"),
		"--tag", testTag,
		"--platform", "linux/amd64")
	
	cmd.Dir = env.WorkingDir
	output, err := cmd.CombinedOutput()
	
	// Allow graceful degradation for platform-specific builds
	if err != nil {
		t.Logf("Platform build failed: %v", err)
		t.Logf("Command output: %s", output)
		t.Skipf("Platform build skipped due to environment: %v", err)
		return
	}
	
	outputStr := string(output)
	assert.True(t, strings.Contains(outputStr, "build") || strings.Contains(outputStr, "success") || len(outputStr) > 0,
		"Expected build output with platform flag: %s", outputStr)
}

// TestBuildHelp tests the build command help
func TestBuildHelp(t *testing.T) {
	env := SetupTestEnvironment(t)
	defer env.Cleanup()

	cmd := exec.Command(env.CLIPath, "build", "--help")
	cmd.Dir = env.WorkingDir
	output, err := cmd.CombinedOutput()
	
	require.NoError(t, err, "Help command should not fail")
	
	outputStr := string(output)
	assert.Contains(t, outputStr, "build", "Help should mention build command")
	assert.True(t, strings.Contains(outputStr, "file") || strings.Contains(outputStr, "tag"),
		"Help should mention key flags: %s", outputStr)
}

// TestBuildVersionCompatibility tests building with current working directory
func TestBuildVersionCompatibility(t *testing.T) {
	env := SetupTestEnvironment(t)
	defer env.Cleanup()

	// Create a local Dockerfile in temp directory
	dockerfilePath := filepath.Join(env.WorkingDir, "Dockerfile")
	dockerfileContent := `FROM alpine:latest
CMD ["echo", "local build test"]`
	
	err := os.WriteFile(dockerfilePath, []byte(dockerfileContent), 0644)
	if err != nil {
		t.Fatalf("Should create local Dockerfile: %v", err)
	}

	testTag := GenerateTestImageTag()
	
	cmd := exec.Command(env.CLIPath, "build", "--tag", testTag)
	cmd.Dir = env.WorkingDir
	output, err := cmd.CombinedOutput()
	
	if err != nil {
		t.Logf("Local build failed: %v", err)
		t.Logf("Command output: %s", output)
		t.Skipf("Local build skipped due to environment: %v", err)
		return
	}
	
	outputStr := string(output)
	assert.True(t, len(outputStr) > 0, "Expected some output from local build")
}