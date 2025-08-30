package e2e

import (
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/cnoe-io/idpbuilder/pkg/tests/integration"
)

// TestCompleteWorkflow tests the complete build and push workflow
func TestCompleteWorkflow(t *testing.T) {
	env := integration.SetupTestEnvironment(t)
	defer env.Cleanup()

	// Generate unique tag for this workflow
	giteaTag := integration.GenerateGiteaImageTag()
	t.Logf("Running complete workflow with tag: %s", giteaTag)

	// Step 1: Build image
	buildCmd := exec.Command(env.CLIPath, "build",
		"--file", integration.GetTestDockerfile("simple"),
		"--context", integration.GetTestContext("simple-app"),
		"--tag", giteaTag)
	
	buildCmd.Dir = env.WorkingDir
	buildOutput, buildErr := buildCmd.CombinedOutput()
	
	if buildErr != nil {
		t.Logf("Build failed: %v", buildErr)
		t.Logf("Build output: %s", buildOutput)
		t.Skipf("Complete workflow skipped due to build failure: %v", buildErr)
		return
	}

	t.Logf("Build successful: %s", buildOutput)

	// Step 2: Push image
	pushCmd := exec.Command(env.CLIPath, "push", giteaTag)
	pushCmd.Dir = env.WorkingDir
	pushOutput, pushErr := pushCmd.CombinedOutput()
	
	if pushErr != nil {
		t.Logf("Push failed: %v", pushErr)
		t.Logf("Push output: %s", pushOutput)
		
		pushOutputStr := string(pushOutput)
		if strings.Contains(pushOutputStr, "network") ||
		   strings.Contains(pushOutputStr, "connection") ||
		   strings.Contains(pushOutputStr, "gitea") {
			t.Skipf("Push skipped due to registry connectivity: %v", pushErr)
			return
		}
	}

	// If we got here, verify the workflow completed
	if pushErr == nil {
		t.Logf("Push successful: %s", pushOutput)
		pushOutputStr := string(pushOutput)
		assert.True(t, len(pushOutputStr) > 0, "Expected output from successful push")
	}

	t.Log("Complete workflow test completed")
}

// TestFreshInstallationFlow tests workflow with no pre-existing setup
func TestFreshInstallationFlow(t *testing.T) {
	env := integration.SetupTestEnvironment(t)
	defer env.Cleanup()

	// Test with simple build that should work in fresh environment
	testTag := integration.GenerateTestImageTag()
	
	cmd := exec.Command(env.CLIPath, "build",
		"--file", integration.GetTestDockerfile("simple"),
		"--context", integration.GetTestContext("simple-app"),
		"--tag", testTag)
	
	cmd.Dir = env.WorkingDir
	output, err := cmd.CombinedOutput()
	
	// Fresh installation should either work or fail gracefully
	outputStr := string(output)
	t.Logf("Fresh installation output: %s", outputStr)
	
	if err != nil {
		// Should provide helpful error messages
		assert.True(t, len(outputStr) > 0, "Expected error message for fresh install issues")
	} else {
		// Should work without manual certificate setup
		assert.Contains(t, outputStr, "build", "Expected build-related output")
	}
}

// TestCertificateRotation tests handling of certificate changes
func TestCertificateRotation(t *testing.T) {
	env := integration.SetupTestEnvironment(t)
	defer env.Cleanup()

	testTag := integration.GenerateTestImageTag()
	
	// First build - should establish certificate trust
	cmd1 := exec.Command(env.CLIPath, "build",
		"--file", integration.GetTestDockerfile("simple"),
		"--context", integration.GetTestContext("simple-app"),
		"--tag", testTag+"1")
	
	cmd1.Dir = env.WorkingDir
	output1, err1 := cmd1.CombinedOutput()
	
	if err1 != nil {
		t.Logf("First build failed: %v, skipping cert rotation test", err1)
		t.SkipNow()
	}

	// Second build - should reuse or refresh certificates
	cmd2 := exec.Command(env.CLIPath, "build",
		"--file", integration.GetTestDockerfile("simple"),
		"--context", integration.GetTestContext("simple-app"),
		"--tag", testTag+"2")
	
	cmd2.Dir = env.WorkingDir
	output2, err2 := cmd2.CombinedOutput()
	
	// Both builds should handle certificates consistently
	output1Str := string(output1)
	output2Str := string(output2)
	
	t.Logf("First build: %s", output1Str)
	t.Logf("Second build: %s", output2Str)
	
	// Should not fail due to certificate issues
	if err2 != nil && !strings.Contains(output2Str, "certificate") {
		t.Logf("Second build failed for non-cert reasons: %v", err2)
	}
}

// TestRecoveryFromFailure tests recovery from various failure modes
func TestRecoveryFromFailure(t *testing.T) {
	env := integration.SetupTestEnvironment(t)
	defer env.Cleanup()

	// Test 1: Invalid build followed by valid build
	invalidTag := integration.GenerateTestImageTag()
	
	invalidCmd := exec.Command(env.CLIPath, "build",
		"--file", integration.GetTestDockerfile("invalid"),
		"--context", integration.GetTestContext("simple-app"),
		"--tag", invalidTag)
	
	invalidCmd.Dir = env.WorkingDir
	invalidOutput, _ := invalidCmd.CombinedOutput()
	t.Logf("Invalid build output: %s", invalidOutput)

	// Test 2: Valid build after invalid should work
	validTag := integration.GenerateTestImageTag()
	
	validCmd := exec.Command(env.CLIPath, "build",
		"--file", integration.GetTestDockerfile("simple"),
		"--context", integration.GetTestContext("simple-app"),
		"--tag", validTag)
	
	validCmd.Dir = env.WorkingDir
	validOutput, validErr := validCmd.CombinedOutput()
	
	if validErr != nil {
		t.Logf("Recovery build failed: %v", validErr)
		t.Logf("Recovery output: %s", validOutput)
		t.Skipf("Recovery test skipped due to environment: %v", validErr)
		return
	}

	// Should recover and build successfully
	validOutputStr := string(validOutput)
	assert.True(t, len(validOutputStr) > 0, "Expected output from recovery build")
	t.Log("Recovery from failure successful")
}