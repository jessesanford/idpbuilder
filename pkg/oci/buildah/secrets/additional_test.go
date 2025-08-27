package secrets

import (
	"os"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMountSecret(t *testing.T) {
	tests := []struct {
		name        string
		mount       *SecretMount
		setupEnv    func()
		cleanupEnv  func()
		expectError bool
		errorMsg    string
	}{
		{
			name: "unsupported_secret_type",
			mount: &SecretMount{
				ID:     "test",
				Source: "test",
				Target: "/tmp/secret",
				Type:   SecretType("unsupported"),
			},
			expectError: true,
			errorMsg:    "unsupported secret type",
		},
		{
			name: "vault_secret_not_implemented",
			mount: &SecretMount{
				ID:     "test",
				Source: "vault:secret/data/test",
				Target: "/tmp/secret",
				Type:   SecretTypeVault,
			},
			expectError: true,
			errorMsg:    "vault secret mounting not yet implemented",
		},
		{
			name: "file_secret_missing_source",
			mount: &SecretMount{
				ID:     "test",
				Source: "/non/existent/file",
				Target: "/tmp/secret",
				Type:   SecretTypeFile,
			},
			expectError: true,
			errorMsg:    "source file not accessible",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sm := createTestSecretManager(t)
			defer sm.Cleanup()

			if tt.setupEnv != nil {
				tt.setupEnv()
			}
			if tt.cleanupEnv != nil {
				defer tt.cleanupEnv()
			}

			path, err := sm.MountSecret(tt.mount)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
				assert.Empty(t, path)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, path)
			}
		})
	}
}

func TestGetSecretPath(t *testing.T) {
	sm := createTestSecretManager(t)
	defer sm.Cleanup()

	// Test non-existent secret
	_, err := sm.GetSecretPath("non-existent")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")

	// Create a secret and test retrieval
	secretID, err := sm.createSecretFile("test", "value")
	require.NoError(t, err)

	path, err := sm.GetSecretPath(secretID)
	assert.NoError(t, err)
	assert.NotEmpty(t, path)
	assert.FileExists(t, path)

	// Remove the file manually to test stale reference cleanup
	os.Remove(path)
	
	_, err = sm.GetSecretPath(secretID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no longer exists")
}

func TestCleanupSecret(t *testing.T) {
	sm := createTestSecretManager(t)
	defer sm.Cleanup()

	// Test cleanup of non-existent secret
	err := sm.CleanupSecret("non-existent")
	assert.NoError(t, err) // Should not error

	// Create a secret and test cleanup
	secretID, err := sm.createSecretFile("test", "value")
	require.NoError(t, err)

	path, err := sm.GetSecretPath(secretID)
	require.NoError(t, err)
	assert.FileExists(t, path)

	// Clean up the secret
	err = sm.CleanupSecret(secretID)
	assert.NoError(t, err)

	// Verify it's gone
	_, err = sm.GetSecretPath(secretID)
	assert.Error(t, err)
	assert.NoFileExists(t, path)
}

func TestSanitizeEnvVars(t *testing.T) {
	sm := createTestSecretManager(t)
	defer sm.Cleanup()

	envVars := map[string]string{
		"VERSION":     "1.0.0",
		"PASSWORD":    "secret123",
		"API_KEY":     "abc123def",
		"PORT":        "8080",
		"SECRET_KEY":  "supersecret",
	}

	sanitized := sm.SanitizeEnvVars(envVars)

	assert.Equal(t, "1.0.0", sanitized["VERSION"])
	assert.Equal(t, "8080", sanitized["PORT"])
	assert.Equal(t, "[REDACTED]", sanitized["PASSWORD"])
	assert.Equal(t, "[REDACTED]", sanitized["API_KEY"])
	assert.Equal(t, "[REDACTED]", sanitized["SECRET_KEY"])
}

func TestGenerateSecretID(t *testing.T) {
	// Test that we get different IDs
	id1 := generateSecretID()
	id2 := generateSecretID()
	
	assert.NotEmpty(t, id1)
	assert.NotEmpty(t, id2)
	assert.NotEqual(t, id1, id2)
	assert.True(t, len(id1) >= 16) // Should be hex encoded, so at least 16 chars
}

func TestIsBase64Like(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"valid_base64", "YWJjZGVmZ2hpamtsbW5vcHFyc3R1dnd4eXo=", true},
		{"valid_base64_no_padding", "YWJjZGVmZ2hpamtsbW5vcA", false},
		{"short_string", "abc", false},
		{"invalid_chars", "abc#def", false},
		{"wrong_length", "abc", false},
		{"empty_string", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isBase64Like(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestRedactPath(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"no_secrets", "/path/to/file", "/path/to/file"},
		{"password_in_path", "/secret/password=abc123/file", "/secret/password[REDACTED]"},
		{"token_in_path", "/config/token:xyz789", "/config/token[REDACTED]"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := redactPath(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestSecretManagerWithoutK8sClient(t *testing.T) {
	sm, err := NewSecretManager(nil, "default", logrus.New())
	require.NoError(t, err)
	defer sm.Cleanup()

	// Test kubernetes secret mounting without client
	mount := &SecretMount{
		ID:     "test-k8s-secret",
		Source: "test-secret",
		Target: "/tmp/secret",
		Type:   SecretTypeKubernetes,
	}

	_, err = sm.MountSecret(mount)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "kubernetes client not configured")
}

func TestCleanupFailureHandling(t *testing.T) {
	sm := createTestSecretManager(t)
	
	// Create a secret file
	secretID, err := sm.createSecretFile("test", "value")
	require.NoError(t, err)

	// Get the path and make it unremovable by changing directory permissions
	_, err = sm.GetSecretPath(secretID)
	require.NoError(t, err)
	
	// Make parent directory read-only to simulate cleanup failure
	parentDir := sm.tempDir
	os.Chmod(parentDir, 0400)
	defer os.Chmod(parentDir, 0700)

	// This should handle the error gracefully
	err = sm.CleanupSecret(secretID)
	// Should still return error for cleanup failure but not panic
	assert.Error(t, err)
}

func TestSecretManagerContextCancellation(t *testing.T) {
	sm := createTestSecretManager(t)
	
	// Cancel the context to stop cleanup routine
	sm.cancel()
	
	// Give it a moment for the routine to notice
	time.Sleep(10 * time.Millisecond)
	
	// Should handle gracefully
	sm.Cleanup()
}