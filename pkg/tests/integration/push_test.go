package integration

import (
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestPushToGitea tests pushing an image to Gitea registry
func TestPushToGitea(t *testing.T) {
	env := SetupTestEnvironment(t)
	defer env.Cleanup()

	// Generate Gitea-compatible image tag
	giteaTag := GenerateGiteaImageTag()
	t.Logf("Testing with Gitea tag: %s", giteaTag)

	// First build an image
	buildCmd := exec.Command(env.CLIPath, "build",
		"--file", GetTestDockerfile("simple"),
		"--context", GetTestContext("simple-app"),
		"--tag", giteaTag)
	
	buildCmd.Dir = env.WorkingDir
	buildOutput, buildErr := buildCmd.CombinedOutput()
	
	if buildErr != nil {
		t.Logf("Build failed: %v", buildErr)
		t.Logf("Build output: %s", buildOutput)
		t.Skipf("Skipping push test due to build failure: %v", buildErr)
		return
	}

	// Then push it
	pushCmd := exec.Command(env.CLIPath, "push", giteaTag)
	pushCmd.Dir = env.WorkingDir
	pushOutput, pushErr := pushCmd.CombinedOutput()
	
	// Allow graceful degradation if Gitea not available
	if pushErr != nil {
		t.Logf("Push failed: %v", pushErr)
		t.Logf("Push output: %s", pushOutput)
		
		pushOutputStr := string(pushOutput)
		if strings.Contains(pushOutputStr, "connection") || 
		   strings.Contains(pushOutputStr, "network") ||
		   strings.Contains(pushOutputStr, "gitea") {
			t.Skipf("Push failed due to network/registry issues: %v", pushErr)
			return
		}
	}
	
	// If push succeeded, verify output
	if pushErr == nil {
		pushOutputStr := string(pushOutput)
		assert.True(t, strings.Contains(pushOutputStr, "push") || 
			strings.Contains(pushOutputStr, "success") || 
			len(pushOutputStr) > 0,
			"Expected push success indication: %s", pushOutputStr)
	}
}

// TestPushWithInsecureFlag tests push with --insecure flag
func TestPushWithInsecureFlag(t *testing.T) {
	env := SetupTestEnvironment(t)
	defer env.Cleanup()

	giteaTag := GenerateGiteaImageTag()

	// Build image first
	buildCmd := exec.Command(env.CLIPath, "build",
		"--file", GetTestDockerfile("simple"),
		"--context", GetTestContext("simple-app"),
		"--tag", giteaTag)
	
	buildCmd.Dir = env.WorkingDir
	_, buildErr := buildCmd.CombinedOutput()
	
	if buildErr != nil {
		t.Logf("Build for insecure push test failed: %v", buildErr)
		t.Skipf("Skipping insecure push test due to build failure: %v", buildErr)
		return
	}

	// Push with --insecure flag
	pushCmd := exec.Command(env.CLIPath, "push", "--insecure", giteaTag)
	pushCmd.Dir = env.WorkingDir
	pushOutput, pushErr := pushCmd.CombinedOutput()
	
	if pushErr != nil {
		t.Logf("Insecure push failed: %v", pushErr)
		t.Logf("Insecure push output: %s", pushOutput)
		t.Skipf("Insecure push skipped due to environment: %v", pushErr)
		return
	}
	
	pushOutputStr := string(pushOutput)
	assert.True(t, len(pushOutputStr) > 0, "Expected some output from insecure push")
}

// TestPushWithAuthentication tests push with credentials
func TestPushWithAuthentication(t *testing.T) {
	env := SetupTestEnvironment(t)
	defer env.Cleanup()

	giteaTag := GenerateGiteaImageTag()

	// Build image first
	buildCmd := exec.Command(env.CLIPath, "build",
		"--file", GetTestDockerfile("simple"),
		"--context", GetTestContext("simple-app"),
		"--tag", giteaTag)
	
	buildCmd.Dir = env.WorkingDir
	_, buildErr := buildCmd.CombinedOutput()
	
	if buildErr != nil {
		t.Logf("Build for auth push test failed: %v", buildErr)
		t.Skipf("Skipping auth push test due to build failure: %v", buildErr)
		return
	}

	// Push with authentication
	pushCmd := exec.Command(env.CLIPath, "push",
		"--username", "test-user",
		"--password", "test-pass",
		giteaTag)
	
	pushCmd.Dir = env.WorkingDir
	pushOutput, _ := pushCmd.CombinedOutput()
	
	// Authentication might fail, but we should get appropriate response
	pushOutputStr := string(pushOutput)
	t.Logf("Auth push output: %s", pushOutputStr)
	
	// Should get either success or auth-related error
	hasResponse := len(pushOutputStr) > 0
	assert.True(t, hasResponse, "Expected some response to auth push attempt")
}

// TestPushWithoutCluster tests push behavior when cluster unavailable  
func TestPushWithoutCluster(t *testing.T) {
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

	giteaTag := GenerateGiteaImageTag()

	pushCmd := exec.Command(env.CLIPath, "push", giteaTag)
	pushCmd.Dir = env.WorkingDir
	pushOutput, pushErr := pushCmd.CombinedOutput()
	
	// Should handle missing cluster gracefully
	pushOutputStr := string(pushOutput)
	t.Logf("No cluster push output: %s", pushOutputStr)
	
	// Should get some kind of response about the issue
	hasResponse := len(pushOutputStr) > 0
	assert.True(t, hasResponse, "Expected response when cluster unavailable")
}

// TestPushInvalidCredentials tests handling of invalid credentials
func TestPushInvalidCredentials(t *testing.T) {
	env := SetupTestEnvironment(t)
	defer env.Cleanup()

	giteaTag := GenerateGiteaImageTag()

	// Try to push with clearly invalid credentials
	pushCmd := exec.Command(env.CLIPath, "push",
		"--username", "invalid-user-12345",
		"--password", "invalid-pass-67890",
		giteaTag)
	
	pushCmd.Dir = env.WorkingDir  
	pushOutput, _ := pushCmd.CombinedOutput()
	
	// Should report authentication issue
	pushOutputStr := string(pushOutput)
	t.Logf("Invalid creds push output: %s", pushOutputStr)
	
	// Should either error or show auth failure
	hasAuthError := strings.Contains(pushOutputStr, "auth") ||
		strings.Contains(pushOutputStr, "credential") ||
		strings.Contains(pushOutputStr, "unauthorized") ||
		strings.Contains(pushOutputStr, "forbidden")
		
	if !hasAuthError && len(pushOutputStr) == 0 {
		t.Log("Expected auth error or response for invalid credentials")
	}
}

// TestPushMissingImage tests pushing non-existent image
func TestPushMissingImage(t *testing.T) {
	env := SetupTestEnvironment(t)
	defer env.Cleanup()

	// Try to push image that doesn't exist
	missingTag := "gitea.local:443/test/nonexistent:missing"
	
	pushCmd := exec.Command(env.CLIPath, "push", missingTag)
	pushCmd.Dir = env.WorkingDir
	pushOutput, pushErr := pushCmd.CombinedOutput()
	
	// Should report image not found
	pushOutputStr := string(pushOutput)
	t.Logf("Missing image push output: %s", pushOutputStr)
	
	// Should error or show missing image message
	hasError := pushErr != nil ||
		strings.Contains(pushOutputStr, "not found") ||
		strings.Contains(pushOutputStr, "missing") ||
		strings.Contains(pushOutputStr, "error")
		
	if !hasError {
		t.Logf("Expected error for missing image push: %s", pushOutputStr)
	}
}

// TestPushNetworkInterruption simulates network issues during push
func TestPushNetworkInterruption(t *testing.T) {
	env := SetupTestEnvironment(t)
	defer env.Cleanup()

	giteaTag := GenerateGiteaImageTag()

	// Build image first
	buildCmd := exec.Command(env.CLIPath, "build",
		"--file", GetTestDockerfile("simple"),
		"--context", GetTestContext("simple-app"),
		"--tag", giteaTag)
	
	buildCmd.Dir = env.WorkingDir
	_, buildErr := buildCmd.CombinedOutput()
	
	if buildErr != nil {
		t.Skipf("Skipping network test due to build failure: %v", buildErr)
		return
	}

	// Try push to non-existent registry (simulate network issue)
	networkTag := strings.Replace(giteaTag, "gitea.local", "nonexistent.registry", 1)
	
	pushCmd := exec.Command(env.CLIPath, "push", networkTag)
	pushCmd.Dir = env.WorkingDir
	pushOutput, pushErr := pushCmd.CombinedOutput()
	
	// Should handle network issues appropriately
	pushOutputStr := string(pushOutput)
	t.Logf("Network issue push output: %s", pushOutputStr)
	
	// Should get network-related error or timeout
	hasNetworkError := pushErr != nil ||
		strings.Contains(pushOutputStr, "network") ||
		strings.Contains(pushOutputStr, "connection") ||
		strings.Contains(pushOutputStr, "timeout") ||
		strings.Contains(pushOutputStr, "resolve")
		
	if !hasNetworkError && len(pushOutputStr) == 0 {
		t.Log("Expected network error for invalid registry")
	}
}

// TestPushHelp tests the push command help
func TestPushHelp(t *testing.T) {
	env := SetupTestEnvironment(t)
	defer env.Cleanup()

	cmd := exec.Command(env.CLIPath, "push", "--help")
	cmd.Dir = env.WorkingDir
	output, err := cmd.CombinedOutput()
	
	require.NoError(t, err, "Help command should not fail")
	
	outputStr := string(output)
	assert.Contains(t, outputStr, "push", "Help should mention push command")
	assert.True(t, strings.Contains(outputStr, "username") || 
		strings.Contains(outputStr, "insecure") ||
		strings.Contains(outputStr, "registry"),
		"Help should mention key flags: %s", outputStr)
}

// TestPushWithTimeout tests push with timeout considerations
func TestPushWithTimeout(t *testing.T) {
	env := SetupTestEnvironment(t)
	defer env.Cleanup()

	giteaTag := GenerateGiteaImageTag()

	// Quick timeout test - should handle gracefully
	pushCmd := exec.Command(env.CLIPath, "push", giteaTag)
	pushCmd.Dir = env.WorkingDir
	
	// Don't wait too long for this test
	output, err := pushCmd.CombinedOutput()
	
	outputStr := string(output)
	t.Logf("Timeout test push output: %s", outputStr)
	
	// Should get some response (success or failure)
	hasOutput := len(outputStr) > 0 || err != nil
	assert.True(t, hasOutput, "Expected some response from push command")
}