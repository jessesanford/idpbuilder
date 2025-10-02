package integration

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestE2EBasicPush tests a complete end-to-end push operation
func TestE2EBasicPush(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	env := NewTestEnvironment(t)
	require.NoError(t, env.SetupIDPBuilder(t))
	defer env.Cleanup(t)

	imageName := fmt.Sprintf("%s/test-basic-e2e:v1.0.0", env.GiteaRegistry)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	cmd := exec.CommandContext(ctx, "idpbuilder", "push",
		"--source", "alpine:latest",
		"--dest", imageName,
		"--username", env.GiteaUsername,
		"--password", env.GiteaPassword,
	)

	output, err := cmd.CombinedOutput()
	require.NoError(t, err, "E2E push failed: %v\nOutput: %s", err, string(output))

	// Verify the image was pushed
	assert.Contains(t, string(output), "successfully", "Output should confirm successful push")

	// Verify image exists in registry
	verifyImageExists(t, env, imageName)
}

// TestE2EMultiArchPush tests pushing multi-architecture images
func TestE2EMultiArchPush(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	env := NewTestEnvironment(t)
	require.NoError(t, env.SetupIDPBuilder(t))
	defer env.Cleanup(t)

	imageName := fmt.Sprintf("%s/test-multiarch:latest", env.GiteaRegistry)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// Use a known multi-arch image
	cmd := exec.CommandContext(ctx, "idpbuilder", "push",
		"--source", "alpine:latest",
		"--dest", imageName,
		"--username", env.GiteaUsername,
		"--password", env.GiteaPassword,
		"--all-platforms",
	)

	output, err := cmd.CombinedOutput()
	require.NoError(t, err, "Multi-arch push failed: %v\nOutput: %s", err, string(output))

	// Verify manifest list was created
	verifyManifestList(t, env, imageName)
}

// TestE2ELargeImagePush tests pushing large images with multiple layers
func TestE2ELargeImagePush(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	env := NewTestEnvironment(t)
	require.NoError(t, env.SetupIDPBuilder(t))
	defer env.Cleanup(t)

	// Use a larger image for this test
	imageName := fmt.Sprintf("%s/test-large:latest", env.GiteaRegistry)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	cmd := exec.CommandContext(ctx, "idpbuilder", "push",
		"--source", "ubuntu:22.04",
		"--dest", imageName,
		"--username", env.GiteaUsername,
		"--password", env.GiteaPassword,
	)

	startTime := time.Now()
	output, err := cmd.CombinedOutput()
	duration := time.Since(startTime)

	require.NoError(t, err, "Large image push failed: %v\nOutput: %s", err, string(output))
	t.Logf("Large image push completed in %v", duration)

	// Verify the image
	verifyImageExists(t, env, imageName)
}

// TestE2ETagValidation tests various tag formats
func TestE2ETagValidation(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	env := NewTestEnvironment(t)
	require.NoError(t, env.SetupIDPBuilder(t))
	defer env.Cleanup(t)

	testTags := []struct {
		tag           string
		expectSuccess bool
	}{
		{"latest", true},
		{"v1.0.0", true},
		{"v1.0.0-alpha.1", true},
		{"v1.0.0+build.123", true},
		{"feature_branch-test", true},
		{"sha-abc123def456", true},
		{"UPPERCASE", true},
		{"with-dashes-123", true},
	}

	for _, tc := range testTags {
		t.Run(fmt.Sprintf("tag_%s", tc.tag), func(t *testing.T) {
			imageName := fmt.Sprintf("%s/test-tags:%s", env.GiteaRegistry, tc.tag)

			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
			defer cancel()

			cmd := exec.CommandContext(ctx, "idpbuilder", "push",
				"--source", "alpine:latest",
				"--dest", imageName,
				"--username", env.GiteaUsername,
				"--password", env.GiteaPassword,
			)

			output, err := cmd.CombinedOutput()

			if tc.expectSuccess {
				assert.NoError(t, err, "Tag %s should be valid: %v\nOutput: %s", tc.tag, err, string(output))
				verifyImageExists(t, env, imageName)
			} else {
				assert.Error(t, err, "Tag %s should be invalid", tc.tag)
			}
		})
	}
}

// TestE2EDigestValidation tests digest-based image references
func TestE2EDigestValidation(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	env := NewTestEnvironment(t)
	require.NoError(t, env.SetupIDPBuilder(t))
	defer env.Cleanup(t)

	imageName := fmt.Sprintf("%s/test-digest:latest", env.GiteaRegistry)

	// First, push an image
	ctx1, cancel1 := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel1()

	cmd1 := exec.CommandContext(ctx1, "idpbuilder", "push",
		"--source", "alpine:latest",
		"--dest", imageName,
		"--username", env.GiteaUsername,
		"--password", env.GiteaPassword,
	)

	output1, err1 := cmd1.CombinedOutput()
	require.NoError(t, err1, "Initial push failed: %v\nOutput: %s", err1, string(output1))

	// Get the digest
	digest := extractDigest(t, string(output1))
	require.NotEmpty(t, digest, "Digest should be present in output")

	t.Logf("Pushed image digest: %s", digest)

	// Verify we can reference by digest
	verifyImageByDigest(t, env, imageName, digest)
}

// TestE2ECompleteWorkflow tests a complete workflow with multiple operations
func TestE2ECompleteWorkflow(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	env := NewTestEnvironment(t)
	require.NoError(t, env.SetupIDPBuilder(t))
	defer env.Cleanup(t)

	baseImage := fmt.Sprintf("%s/workflow-base:v1", env.GiteaRegistry)

	// Step 1: Push base image
	t.Log("Step 1: Pushing base image")
	ctx1, cancel1 := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel1()

	cmd1 := exec.CommandContext(ctx1, "idpbuilder", "push",
		"--source", "alpine:latest",
		"--dest", baseImage,
		"--username", env.GiteaUsername,
		"--password", env.GiteaPassword,
	)

	_, err1 := cmd1.CombinedOutput()
	require.NoError(t, err1, "Base image push failed: %v", err1)
	t.Log("Base image pushed successfully")

	// Step 2: Push with different tag
	t.Log("Step 2: Pushing with different tag")
	taggedImage := fmt.Sprintf("%s/workflow-base:v2", env.GiteaRegistry)

	ctx2, cancel2 := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel2()

	cmd2 := exec.CommandContext(ctx2, "idpbuilder", "push",
		"--source", "alpine:latest",
		"--dest", taggedImage,
		"--username", env.GiteaUsername,
		"--password", env.GiteaPassword,
	)

	_, err2 := cmd2.CombinedOutput()
	require.NoError(t, err2, "Tagged image push failed: %v", err2)
	t.Log("Tagged image pushed successfully")

	// Step 3: Verify both images exist
	t.Log("Step 3: Verifying images")
	verifyImageExists(t, env, baseImage)
	verifyImageExists(t, env, taggedImage)

	t.Log("Complete workflow test passed")
}

// TestE2EStreamingProgress tests that progress is reported during push
func TestE2EStreamingProgress(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	env := NewTestEnvironment(t)
	require.NoError(t, env.SetupIDPBuilder(t))
	defer env.Cleanup(t)

	imageName := fmt.Sprintf("%s/test-progress:latest", env.GiteaRegistry)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	cmd := exec.CommandContext(ctx, "idpbuilder", "push",
		"--source", "alpine:latest",
		"--dest", imageName,
		"--username", env.GiteaUsername,
		"--password", env.GiteaPassword,
		"--verbose",
	)

	// Capture output to verify progress reporting
	output, err := cmd.CombinedOutput()
	require.NoError(t, err, "Push failed: %v\nOutput: %s", err, string(output))

	outputStr := string(output)

	// Check for progress indicators
	t.Logf("Command output: %s", outputStr)
	assert.True(t, len(outputStr) > 0, "Should have progress output")
}

// Helper functions

// verifyImageExists verifies that an image exists in the registry
func verifyImageExists(t *testing.T, env *TestEnvironment, imageName string) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "crane", "manifest", imageName)
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("CRANE_USERNAME=%s", env.GiteaUsername),
		fmt.Sprintf("CRANE_PASSWORD=%s", env.GiteaPassword),
	)

	output, err := cmd.CombinedOutput()
	assert.NoError(t, err, "Failed to get image manifest: %v\nOutput: %s", err, string(output))
	assert.True(t, len(output) > 0, "Manifest should not be empty")

	t.Logf("Verified image exists: %s", imageName)
}

// verifyManifestList verifies that a manifest list exists for multi-arch images
func verifyManifestList(t *testing.T, env *TestEnvironment, imageName string) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "crane", "manifest", imageName)
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("CRANE_USERNAME=%s", env.GiteaUsername),
		fmt.Sprintf("CRANE_PASSWORD=%s", env.GiteaPassword),
	)

	output, err := cmd.CombinedOutput()
	assert.NoError(t, err, "Failed to get manifest list: %v", err)

	// Check if it's a manifest list (contains manifests array)
	outputStr := string(output)
	assert.Contains(t, outputStr, "manifests", "Should be a manifest list")

	t.Logf("Verified manifest list exists: %s", imageName)
}

// verifyImageByDigest verifies an image can be referenced by digest
func verifyImageByDigest(t *testing.T, env *TestEnvironment, imageName, digest string) {
	// Extract registry and repo from imageName
	parts := strings.Split(imageName, "/")
	if len(parts) < 2 {
		t.Fatalf("Invalid image name: %s", imageName)
	}

	registry := parts[0]
	repo := strings.Split(parts[1], ":")[0]
	imageByDigest := fmt.Sprintf("%s/%s@%s", registry, repo, digest)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "crane", "manifest", imageByDigest)
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("CRANE_USERNAME=%s", env.GiteaUsername),
		fmt.Sprintf("CRANE_PASSWORD=%s", env.GiteaPassword),
	)

	output, err := cmd.CombinedOutput()
	assert.NoError(t, err, "Failed to get image by digest: %v", err)
	assert.True(t, len(output) > 0, "Manifest should not be empty")

	t.Logf("Verified image accessible by digest: %s", imageByDigest)
}

// extractDigest extracts the digest from command output
func extractDigest(t *testing.T, output string) string {
	// Look for sha256:... in the output
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if strings.Contains(line, "sha256:") {
			parts := strings.Fields(line)
			for _, part := range parts {
				if strings.HasPrefix(part, "sha256:") {
					return part
				}
			}
		}
	}

	t.Log("No digest found in output, trying to extract differently")
	// Alternative: look for digest: prefix
	for _, line := range lines {
		if strings.Contains(line, "digest:") {
			parts := strings.Split(line, "digest:")
			if len(parts) > 1 {
				digest := strings.TrimSpace(parts[1])
				if strings.HasPrefix(digest, "sha256:") {
					return digest
				}
			}
		}
	}

	return ""
}

// TestE2EErrorRecovery tests that the system recovers from various error conditions
func TestE2EErrorRecovery(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	env := NewTestEnvironment(t)
	require.NoError(t, env.SetupIDPBuilder(t))
	defer env.Cleanup(t)

	// Test 1: Recover from invalid source
	t.Run("invalid_source", func(t *testing.T) {
		imageName := fmt.Sprintf("%s/test-recovery:latest", env.GiteaRegistry)

		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
		defer cancel()

		cmd := exec.CommandContext(ctx, "idpbuilder", "push",
			"--source", "nonexistent-image:invalid-tag",
			"--dest", imageName,
			"--username", env.GiteaUsername,
			"--password", env.GiteaPassword,
		)

		_, err := cmd.CombinedOutput()
		assert.Error(t, err, "Should fail with invalid source")
	})

	// Test 2: Valid push after error
	t.Run("valid_after_error", func(t *testing.T) {
		imageName := fmt.Sprintf("%s/test-recovery:latest", env.GiteaRegistry)

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cancel()

		cmd := exec.CommandContext(ctx, "idpbuilder", "push",
			"--source", "alpine:latest",
			"--dest", imageName,
			"--username", env.GiteaUsername,
			"--password", env.GiteaPassword,
		)

		output, err := cmd.CombinedOutput()
		require.NoError(t, err, "Should succeed after previous error: %v\nOutput: %s", err, string(output))
	})
}
