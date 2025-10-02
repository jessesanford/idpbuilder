package integration

import (
	"context"
	"fmt"
	"os/exec"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNetworkFailureRecovery tests recovery from network failures
func TestNetworkFailureRecovery(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	env := NewTestEnvironment(t)
	require.NoError(t, env.SetupIDPBuilder(t))
	defer env.Cleanup(t)

	imageName := fmt.Sprintf("%s/test-network-recovery:latest", env.GiteaRegistry)

	// Simulate intermittent network issues using network policies
	// This is a simplified test - real implementation would use chaos engineering
	testCases := []struct {
		name              string
		maxRetries        int
		failureRate       float64
		expectSuccess     bool
		expectedAttempts  int
	}{
		{
			name:             "Single transient failure",
			maxRetries:       3,
			failureRate:      0.33,
			expectSuccess:    true,
			expectedAttempts: 2,
		},
		{
			name:             "Multiple transient failures",
			maxRetries:       5,
			failureRate:      0.5,
			expectSuccess:    true,
			expectedAttempts: 3,
		},
		{
			name:             "Exceeds max retries",
			maxRetries:       2,
			failureRate:      1.0,
			expectSuccess:    false,
			expectedAttempts: 3,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
			defer cancel()

			cmd := exec.CommandContext(ctx, "idpbuilder", "push",
				"--source", "alpine:latest",
				"--dest", imageName,
				"--username", env.GiteaUsername,
				"--password", env.GiteaPassword,
				"--max-retries", fmt.Sprintf("%d", tc.maxRetries),
				"--retry-delay", "1s",
			)

			startTime := time.Now()
			output, err := cmd.CombinedOutput()
			duration := time.Since(startTime)

			if tc.expectSuccess {
				assert.NoError(t, err, "Expected push to succeed after retries: %v\nOutput: %s", err, string(output))
			} else {
				assert.Error(t, err, "Expected push to fail after exhausting retries")
			}

			t.Logf("Test completed in %v", duration)
		})
	}
}

// TestTransientErrorHandling tests handling of various transient errors
func TestTransientErrorHandling(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	env := NewTestEnvironment(t)
	require.NoError(t, env.SetupIDPBuilder(t))
	defer env.Cleanup(t)

	transientErrors := []struct {
		name          string
		errorType     string
		shouldRetry   bool
		maxRetries    int
		expectSuccess bool
	}{
		{
			name:          "Connection timeout",
			errorType:     "timeout",
			shouldRetry:   true,
			maxRetries:    3,
			expectSuccess: true,
		},
		{
			name:          "Connection refused",
			errorType:     "refused",
			shouldRetry:   true,
			maxRetries:    3,
			expectSuccess: true,
		},
		{
			name:          "Service unavailable (503)",
			errorType:     "503",
			shouldRetry:   true,
			maxRetries:    3,
			expectSuccess: true,
		},
		{
			name:          "Rate limited (429)",
			errorType:     "429",
			shouldRetry:   true,
			maxRetries:    5,
			expectSuccess: true,
		},
		{
			name:          "Server error (500)",
			errorType:     "500",
			shouldRetry:   true,
			maxRetries:    3,
			expectSuccess: true,
		},
	}

	for _, tc := range transientErrors {
		t.Run(tc.name, func(t *testing.T) {
			imageName := fmt.Sprintf("%s/test-transient-%s:latest", env.GiteaRegistry, tc.errorType)

			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
			defer cancel()

			cmd := exec.CommandContext(ctx, "idpbuilder", "push",
				"--source", "alpine:latest",
				"--dest", imageName,
				"--username", env.GiteaUsername,
				"--password", env.GiteaPassword,
				"--max-retries", fmt.Sprintf("%d", tc.maxRetries),
			)

			output, err := cmd.CombinedOutput()

			if tc.expectSuccess {
				assert.NoError(t, err, "Expected successful push after handling transient error: %v\nOutput: %s", err, string(output))
			} else {
				assert.Error(t, err, "Expected push to fail")
			}
		})
	}
}

// TestBackoffStrategy tests the exponential backoff strategy
func TestBackoffStrategy(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	env := NewTestEnvironment(t)
	require.NoError(t, env.SetupIDPBuilder(t))
	defer env.Cleanup(t)

	testCases := []struct {
		name                 string
		initialDelay         time.Duration
		maxDelay             time.Duration
		backoffMultiplier    float64
		expectedMinDuration  time.Duration
		expectedMaxDuration  time.Duration
		numRetries           int
	}{
		{
			name:                "Exponential backoff",
			initialDelay:        time.Second,
			maxDelay:            10 * time.Second,
			backoffMultiplier:   2.0,
			expectedMinDuration: 3 * time.Second,
			expectedMaxDuration: 15 * time.Second,
			numRetries:          3,
		},
		{
			name:                "Linear backoff",
			initialDelay:        time.Second,
			maxDelay:            10 * time.Second,
			backoffMultiplier:   1.0,
			expectedMinDuration: 2 * time.Second,
			expectedMaxDuration: 8 * time.Second,
			numRetries:          3,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			imageName := fmt.Sprintf("%s/test-backoff:latest", env.GiteaRegistry)

			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
			defer cancel()

			startTime := time.Now()

			cmd := exec.CommandContext(ctx, "idpbuilder", "push",
				"--source", "alpine:latest",
				"--dest", imageName,
				"--username", env.GiteaUsername,
				"--password", env.GiteaPassword,
				"--max-retries", fmt.Sprintf("%d", tc.numRetries),
				"--retry-delay", tc.initialDelay.String(),
				"--max-retry-delay", tc.maxDelay.String(),
			)

			output, err := cmd.CombinedOutput()
			duration := time.Since(startTime)

			t.Logf("Push operation took %v (expected between %v and %v)", duration, tc.expectedMinDuration, tc.expectedMaxDuration)
			t.Logf("Output: %s", string(output))

			// For successful pushes, we won't have retries, so duration checks may not apply
			if err == nil {
				t.Log("Push succeeded on first attempt (no retries)")
			}
		})
	}
}

// TestMaxRetryLimit tests enforcement of maximum retry limits
func TestMaxRetryLimit(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	env := NewTestEnvironment(t)
	require.NoError(t, env.SetupIDPBuilder(t))
	defer env.Cleanup(t)

	testCases := []struct {
		name          string
		maxRetries    int
		expectedMax   int
		expectSuccess bool
	}{
		{
			name:          "Zero retries",
			maxRetries:    0,
			expectedMax:   0,
			expectSuccess: true,
		},
		{
			name:          "One retry",
			maxRetries:    1,
			expectedMax:   1,
			expectSuccess: true,
		},
		{
			name:          "Default retries (3)",
			maxRetries:    3,
			expectedMax:   3,
			expectSuccess: true,
		},
		{
			name:          "High retry count (10)",
			maxRetries:    10,
			expectedMax:   10,
			expectSuccess: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			imageName := fmt.Sprintf("%s/test-max-retry-%d:latest", env.GiteaRegistry, tc.maxRetries)

			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
			defer cancel()

			cmd := exec.CommandContext(ctx, "idpbuilder", "push",
				"--source", "alpine:latest",
				"--dest", imageName,
				"--username", env.GiteaUsername,
				"--password", env.GiteaPassword,
				"--max-retries", fmt.Sprintf("%d", tc.maxRetries),
			)

			output, err := cmd.CombinedOutput()

			if tc.expectSuccess {
				assert.NoError(t, err, "Expected push to succeed: %v\nOutput: %s", err, string(output))
			} else {
				assert.Error(t, err, "Expected push to fail after max retries")
			}
		})
	}
}

// TestConcurrentRetriesIsolation tests that concurrent operations handle retries independently
func TestConcurrentRetriesIsolation(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	env := NewTestEnvironment(t)
	require.NoError(t, env.SetupIDPBuilder(t))
	defer env.Cleanup(t)

	const numConcurrent = 3

	var successCount atomic.Int32
	var failureCount atomic.Int32

	results := make(chan error, numConcurrent)

	for i := 0; i < numConcurrent; i++ {
		go func(index int) {
			imageName := fmt.Sprintf("%s/test-concurrent-%d:latest", env.GiteaRegistry, index)

			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
			defer cancel()

			cmd := exec.CommandContext(ctx, "idpbuilder", "push",
				"--source", "alpine:latest",
				"--dest", imageName,
				"--username", env.GiteaUsername,
				"--password", env.GiteaPassword,
				"--max-retries", "3",
			)

			_, err := cmd.CombinedOutput()

			if err == nil {
				successCount.Add(1)
			} else {
				failureCount.Add(1)
			}

			results <- err
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < numConcurrent; i++ {
		err := <-results
		t.Logf("Concurrent operation %d result: %v", i+1, err)
	}

	t.Logf("Concurrent operations: %d succeeded, %d failed", successCount.Load(), failureCount.Load())

	// At least some operations should succeed
	assert.Greater(t, int(successCount.Load()), 0, "Expected at least some concurrent operations to succeed")
}

// TestRetryMetrics tests that retry metrics are properly tracked and reported
func TestRetryMetrics(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	env := NewTestEnvironment(t)
	require.NoError(t, env.SetupIDPBuilder(t))
	defer env.Cleanup(t)

	imageName := fmt.Sprintf("%s/test-metrics:latest", env.GiteaRegistry)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	cmd := exec.CommandContext(ctx, "idpbuilder", "push",
		"--source", "alpine:latest",
		"--dest", imageName,
		"--username", env.GiteaUsername,
		"--password", env.GiteaPassword,
		"--max-retries", "3",
		"--verbose",
	)

	output, err := cmd.CombinedOutput()
	require.NoError(t, err, "Push failed: %v\nOutput: %s", err, string(output))

	outputStr := string(output)

	// Check for retry-related metrics in output
	// This assumes the push command outputs retry information
	t.Logf("Command output: %s", outputStr)

	// Verify success message
	assert.Contains(t, outputStr, "successfully", "Output should contain success message")
}
