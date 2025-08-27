package secrets

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestNewSecretManager(t *testing.T) {
	tests := []struct {
		name      string
		k8sClient *fake.Clientset
		namespace string
		logger    *logrus.Logger
		wantErr   bool
	}{
		{
			name:      "valid_creation",
			k8sClient: fake.NewSimpleClientset(),
			namespace: "default",
			logger:    logrus.New(),
			wantErr:   false,
		},
		{
			name:      "nil_logger_should_work",
			k8sClient: fake.NewSimpleClientset(),
			namespace: "test-ns",
			logger:    nil,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sm, err := NewSecretManager(tt.k8sClient, tt.namespace, tt.logger)
			
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.NotNil(t, sm)
			assert.Equal(t, tt.namespace, sm.namespace)
			assert.NotEmpty(t, sm.tempDir)
			assert.NotNil(t, sm.secretFiles)
			assert.NotNil(t, sm.cleanupTicker)

			// Verify temp directory exists and has correct permissions
			stat, err := os.Stat(sm.tempDir)
			require.NoError(t, err)
			assert.True(t, stat.IsDir())
			assert.Equal(t, os.FileMode(0700), stat.Mode().Perm())

			// Cleanup
			sm.Cleanup()
		})
	}
}

func TestSanitizeBuildArgs(t *testing.T) {
	sm := createTestSecretManager(t)
	defer sm.Cleanup()

	tests := []struct {
		name          string
		args          map[string]string
		expectedArgs  map[string]string
		expectedCount int
	}{
		{
			name: "no_secrets",
			args: map[string]string{
				"VERSION": "1.0.0",
				"ENV":     "production",
			},
			expectedArgs: map[string]string{
				"VERSION": "1.0.0",
				"ENV":     "production",
			},
			expectedCount: 0,
		},
		{
			name: "password_secret",
			args: map[string]string{
				"VERSION":  "1.0.0",
				"PASSWORD": "secret123",
			},
			expectedArgs: map[string]string{
				"VERSION": "1.0.0",
			},
			expectedCount: 1,
		},
		{
			name: "multiple_secrets",
			args: map[string]string{
				"VERSION":     "1.0.0",
				"DB_PASSWORD": "dbpass123",
				"API_KEY":     "abc123def456",
				"SECRET_KEY":  "secretvalue",
			},
			expectedArgs: map[string]string{
				"VERSION": "1.0.0",
			},
			expectedCount: 3,
		},
		{
			name: "jwt_like_token",
			args: map[string]string{
				"VERSION": "1.0.0",
				"TOKEN":   "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIn0.signature",
			},
			expectedArgs: map[string]string{
				"VERSION": "1.0.0",
			},
			expectedCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buildArgs, err := sm.SanitizeBuildArgs(tt.args)
			
			require.NoError(t, err)
			assert.Equal(t, tt.expectedArgs, buildArgs.Args)
			assert.Len(t, buildArgs.SecretIDs, tt.expectedCount)

			// Verify secret files were created
			for _, secretID := range buildArgs.SecretIDs {
				path, err := sm.GetSecretPath(secretID)
				assert.NoError(t, err)
				assert.FileExists(t, path)
			}
		})
	}
}

func TestIsSecretArg(t *testing.T) {
	sm := createTestSecretManager(t)
	defer sm.Cleanup()

	tests := []struct {
		name     string
		key      string
		value    string
		expected bool
	}{
		{"password_key", "PASSWORD", "secret123", true},
		{"api_key", "API_KEY", "abc123", true},
		{"token_key", "ACCESS_TOKEN", "token123", true},
		{"secret_key", "SECRET_KEY", "secret", true},
		{"regular_env", "VERSION", "1.0.0", false},
		{"regular_env2", "ENV", "production", false},
		{"jwt_like", "AUTH", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIn0.sig", true},
		{"base64_like", "DATA", "YWJjZGVmZ2hpamtsbW5vcHFyc3R1dnd4eXo=", true},
		{"short_value", "BUILD_VERSION", "abc", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sm.isSecretArg(tt.key, tt.value)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestMountFileSecret(t *testing.T) {
	sm := createTestSecretManager(t)
	defer sm.Cleanup()

	// Create a test file
	tempFile, err := os.CreateTemp("", "test-secret-*")
	require.NoError(t, err)
	defer os.Remove(tempFile.Name())

	testContent := "secret-content-123"
	_, err = tempFile.WriteString(testContent)
	require.NoError(t, err)
	tempFile.Close()

	mount := &SecretMount{
		ID:     "test-file-secret",
		Source: tempFile.Name(),
		Target: "/tmp/secret",
		Type:   SecretTypeFile,
	}

	path, err := sm.mountFileSecret(mount)
	require.NoError(t, err)
	assert.NotEmpty(t, path)
	assert.FileExists(t, path)

	// Verify content
	content, err := os.ReadFile(path)
	require.NoError(t, err)
	assert.Equal(t, testContent, string(content))

	// Verify permissions
	stat, err := os.Stat(path)
	require.NoError(t, err)
	assert.Equal(t, os.FileMode(0600), stat.Mode().Perm())
}

func TestMountKubernetesSecret(t *testing.T) {
	tests := []struct {
		name        string
		secretName  string
		secretData  map[string][]byte
		source      string
		expectError bool
		expectData  string
	}{
		{
			name:       "single_key_secret",
			secretName: "test-secret",
			secretData: map[string][]byte{
				"password": []byte("secret123"),
			},
			source:      "test-secret",
			expectError: false,
			expectData:  "secret123",
		},
		{
			name:       "specific_key",
			secretName: "multi-key-secret",
			secretData: map[string][]byte{
				"username": []byte("admin"),
				"password": []byte("secret456"),
			},
			source:      "multi-key-secret/password",
			expectError: false,
			expectData:  "secret456",
		},
		{
			name:       "missing_key",
			secretName: "test-secret",
			secretData: map[string][]byte{
				"password": []byte("secret123"),
			},
			source:      "test-secret/missing",
			expectError: true,
		},
		{
			name:       "multiple_keys_no_spec",
			secretName: "multi-key-secret",
			secretData: map[string][]byte{
				"username": []byte("admin"),
				"password": []byte("secret456"),
			},
			source:      "multi-key-secret",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k8sClient := fake.NewSimpleClientset()
			sm, err := NewSecretManager(k8sClient, "default", logrus.New())
			require.NoError(t, err)
			defer sm.Cleanup()

			// Create the secret
			secret := &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      tt.secretName,
					Namespace: "default",
				},
				Data: tt.secretData,
			}
			_, err = k8sClient.CoreV1().Secrets("default").Create(
				context.TODO(), secret, metav1.CreateOptions{})
			require.NoError(t, err)

			mount := &SecretMount{
				ID:     "test-k8s-secret",
				Source: tt.source,
				Target: "/tmp/secret",
				Type:   SecretTypeKubernetes,
			}

			path, err := sm.mountKubernetesSecret(mount)
			
			if tt.expectError {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.NotEmpty(t, path)
			assert.FileExists(t, path)

			// Verify content
			content, err := os.ReadFile(path)
			require.NoError(t, err)
			assert.Equal(t, tt.expectData, string(content))
		})
	}
}

func TestMountEnvSecret(t *testing.T) {
	sm := createTestSecretManager(t)
	defer sm.Cleanup()

	// Set test environment variable
	testEnvVar := "TEST_SECRET_VAR"
	testValue := "env-secret-value-123"
	os.Setenv(testEnvVar, testValue)
	defer os.Unsetenv(testEnvVar)

	mount := &SecretMount{
		ID:     "test-env-secret",
		Source: testEnvVar,
		Target: "/tmp/secret",
		Type:   SecretTypeEnv,
	}

	path, err := sm.mountEnvSecret(mount)
	require.NoError(t, err)
	assert.NotEmpty(t, path)
	assert.FileExists(t, path)

	// Verify content
	content, err := os.ReadFile(path)
	require.NoError(t, err)
	assert.Equal(t, testValue, string(content))

	// Test missing environment variable
	mount.Source = "NON_EXISTENT_VAR"
	_, err = sm.mountEnvSecret(mount)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found or empty")
}

func TestRedactLogMessage(t *testing.T) {
	sm := createTestSecretManager(t)
	defer sm.Cleanup()

	tests := []struct {
		name     string
		message  string
		expected string
	}{
		{
			name:     "password_redaction",
			message:  "Setting password=secret123 for user",
			expected: "Setting password[REDACTED] for user",
		},
		{
			name:     "bearer_token",
			message:  "Authorization: bearer abc123def456",
			expected: "Authorization: bearer [REDACTED]",
		},
		{
			name:     "api_key",
			message:  "Using API_KEY=sk-1234567890abcdef",
			expected: "Using API_KEY[REDACTED]",
		},
		{
			name:     "no_secrets",
			message:  "Starting application version 1.0.0",
			expected: "Starting application version 1.0.0",
		},
		{
			name:     "multiple_secrets",
			message:  "Config: password=secret123 token=abc456",
			expected: "Config: password[REDACTED] token[REDACTED]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sm.RedactLogMessage(tt.message)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCleanupExpiredSecrets(t *testing.T) {
	sm := createTestSecretManager(t)
	defer sm.Cleanup()

	// Create a secret with short TTL
	secretID, err := sm.createSecretFile("test", "value")
	require.NoError(t, err)

	// Verify secret exists
	path, err := sm.GetSecretPath(secretID)
	require.NoError(t, err)
	assert.FileExists(t, path)

	// Manually expire the secret
	sm.mutex.Lock()
	sm.secretFiles[secretID].TTL = 1 * time.Millisecond
	sm.secretFiles[secretID].CreatedAt = time.Now().Add(-1 * time.Second)
	sm.mutex.Unlock()

	// Run cleanup
	sm.cleanupExpiredSecrets()

	// Verify secret is cleaned up
	_, err = sm.GetSecretPath(secretID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestGetStats(t *testing.T) {
	sm := createTestSecretManager(t)
	defer sm.Cleanup()

	// Initially no secrets
	stats := sm.GetStats()
	assert.Equal(t, 0, stats["active_secrets"])
	assert.Equal(t, 0, stats["secrets_under_5min"])
	assert.Equal(t, 0, stats["secrets_over_5min"])

	// Create some secrets
	_, err := sm.createSecretFile("test1", "value1")
	require.NoError(t, err)
	_, err = sm.createSecretFile("test2", "value2")
	require.NoError(t, err)

	stats = sm.GetStats()
	assert.Equal(t, 2, stats["active_secrets"])
	assert.Equal(t, 2, stats["secrets_under_5min"])
	assert.Equal(t, 0, stats["secrets_over_5min"])
}

// Helper function to create test secret manager
func createTestSecretManager(t *testing.T) *SecretManager {
	k8sClient := fake.NewSimpleClientset()
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)
	
	sm, err := NewSecretManager(k8sClient, "default", logger)
	require.NoError(t, err)
	
	return sm
}