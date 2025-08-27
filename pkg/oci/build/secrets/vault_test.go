package secrets

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"
)

func TestVaultSecureStorage(t *testing.T) {
	vault, err := NewVault()
	if err != nil {
		t.Fatalf("Failed to create vault: %v", err)
	}
	defer vault.Clear()

	// Create test secret
	originalValue := []byte("super-secret-password-123")
	secret := &Secret{
		ID:     "test-secret",
		Type:   SecretTypeBuildArg,
		Value:  originalValue,
		Target: "DATABASE_PASSWORD",
		Source: "env:DB_PASS",
	}

	// Test storage
	if err := vault.Store(secret); err != nil {
		t.Errorf("Failed to store secret: %v", err)
	}

	// Verify original value was cleared from memory (security requirement)
	originalCleared := true
	for _, b := range secret.Value {
		if b != 0 {
			originalCleared = false
			break
		}
	}
	if !originalCleared {
		t.Error("SECURITY VIOLATION: Original secret value not cleared from memory")
	}

	// Test retrieval
	retrieved, err := vault.Retrieve("test-secret")
	if err != nil {
		t.Errorf("Failed to retrieve secret: %v", err)
	}

	// Verify retrieved value matches original
	if string(retrieved.Value) != "super-secret-password-123" {
		t.Error("Retrieved value doesn't match original")
	}

	// Verify metadata preserved
	if retrieved.Type != SecretTypeBuildArg {
		t.Errorf("Type mismatch: got %s, want %s", retrieved.Type, SecretTypeBuildArg)
	}
	if retrieved.Target != "DATABASE_PASSWORD" {
		t.Errorf("Target mismatch: got %s, want DATABASE_PASSWORD", retrieved.Target)
	}
}

func TestVaultEncryption(t *testing.T) {
	vault, _ := NewVault()
	defer vault.Clear()

	secret := &Secret{
		ID:    "encrypt-test",
		Type:  SecretTypeMount,
		Value: []byte("this-should-be-encrypted"),
	}

	// Store secret
	vault.Store(secret)

	// Access internal encrypted storage (testing encryption)
	memVault := vault.(*memoryVault)
	memVault.mu.RLock()
	encSecret := memVault.secrets["encrypt-test"]
	memVault.mu.RUnlock()

	// Verify encrypted data doesn't contain plaintext
	if strings.Contains(string(encSecret.encrypted), "this-should-be-encrypted") {
		t.Error("SECURITY VIOLATION: Secret appears to be stored in plaintext")
	}

	// Verify encrypted data is longer (due to nonce and authentication tag)
	if len(encSecret.encrypted) <= len("this-should-be-encrypted") {
		t.Error("Encrypted data should be longer than plaintext (nonce + auth tag)")
	}
}

func TestVaultSecureCleanup(t *testing.T) {
	vault, _ := NewVault()

	// Store multiple secrets
	for i := 0; i < 5; i++ {
		secret := &Secret{
			ID:    fmt.Sprintf("cleanup-secret-%d", i),
			Type:  SecretTypeBuildArg,
			Value: []byte(fmt.Sprintf("cleanup-value-%d", i)),
		}
		vault.Store(secret)
	}

	// Clear vault
	if err := vault.Clear(); err != nil {
		t.Errorf("Failed to clear vault: %v", err)
	}

	// Verify all secrets are gone
	for i := 0; i < 5; i++ {
		id := fmt.Sprintf("cleanup-secret-%d", i)
		if _, err := vault.Retrieve(id); err == nil {
			t.Errorf("SECURITY VIOLATION: Secret %s still retrievable after clear", id)
		}
	}
}

func TestSanitizerPreventsLeaks(t *testing.T) {
	sanitizer := NewSanitizer()

	// Register multiple secrets
	secrets := map[string]string{
		"api-key":    "sk-1234567890abcdef",
		"db-password": "postgres123!@#",
		"token":      "bearer-token-xyz789",
	}

	for id, value := range secrets {
		if err := sanitizer.RegisterSecret(id, []byte(value)); err != nil {
			t.Errorf("Failed to register secret %s: %v", id, err)
		}
	}

	// Test input with multiple secrets
	input := `
	Connecting to database with password: postgres123!@#
	Using API key: sk-1234567890abcdef
	Authorization header: bearer-token-xyz789
	`

	output := sanitizer.Sanitize(input)

	// Verify no secrets remain in output
	for id, value := range secrets {
		if strings.Contains(output, value) {
			t.Errorf("SECURITY VIOLATION: Secret %s not sanitized from output", id)
		}
		expectedMask := "***" + strings.ToUpper(id) + "***"
		if !strings.Contains(output, expectedMask) {
			t.Errorf("Secret %s not replaced with expected mask %s", id, expectedMask)
		}
	}
}

func TestSanitizerCommonPatterns(t *testing.T) {
	sanitizer := NewSanitizer()

	testCases := []struct {
		name     string
		input    string
		shouldNotContain []string
		shouldContain    []string
	}{
		{
			name:  "Password patterns",
			input: "password=mysecret123 --password secret456",
			shouldNotContain: []string{"mysecret123", "secret456"},
			shouldContain:    []string{"***REDACTED***"},
		},
		{
			name:  "Bearer tokens",
			input: "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
			shouldNotContain: []string{"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9"},
			shouldContain:    []string{"***TOKEN***"},
		},
		{
			name:  "Connection strings",
			input: "mongodb://user:password@localhost:27017/db",
			shouldNotContain: []string{"user:password"},
			shouldContain:    []string{"***USER***:***PASS***"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			output := sanitizer.Sanitize(tc.input)
			
			for _, forbidden := range tc.shouldNotContain {
				if strings.Contains(output, forbidden) {
					t.Errorf("Output still contains forbidden pattern: %s", forbidden)
				}
			}
			
			for _, required := range tc.shouldContain {
				if !strings.Contains(output, required) {
					t.Errorf("Output missing required pattern: %s", required)
				}
			}
		})
	}
}

func TestInjectorBuildArgs(t *testing.T) {
	vault, _ := NewVault()
	sanitizer := NewSanitizer()
	injector := NewInjector(vault, sanitizer)
	defer vault.Clear()

	// Store test secret
	secret := &Secret{
		ID:     "build-arg-secret",
		Type:   SecretTypeBuildArg,
		Value:  []byte("build-secret-value"),
		Target: "BUILD_SECRET",
	}
	if err := vault.Store(secret); err != nil {
		t.Fatalf("Failed to store secret: %v", err)
	}

	// Test build arg injection
	ctx := context.Background()
	baseArgs := map[string]string{
		"VERSION": "1.0.0",
		"ENV":     "production",
	}

	injectedArgs, err := injector.InjectBuildArgs(ctx, baseArgs, []string{"build-arg-secret"})
	if err != nil {
		t.Errorf("Failed to inject build args: %v", err)
	}

	// Verify injection worked
	if injectedArgs["BUILD_SECRET"] != "build-secret-value" {
		t.Error("Build arg secret not injected correctly")
	}

	// Verify existing args preserved
	if injectedArgs["VERSION"] != "1.0.0" {
		t.Error("Existing arg lost during injection")
	}
	if injectedArgs["ENV"] != "production" {
		t.Error("Existing arg lost during injection")
	}
}

func TestInjectorSecretMount(t *testing.T) {
	vault, _ := NewVault()
	sanitizer := NewSanitizer()
	injector := NewInjector(vault, sanitizer)
	defer vault.Clear()

	// Store mount secret
	secret := &Secret{
		ID:     "mount-secret",
		Type:   SecretTypeMount,
		Value:  []byte("secret-file-content"),
		Target: "secret.txt",
		Mode:   0600,
	}
	if err := vault.Store(secret); err != nil {
		t.Fatalf("Failed to store secret: %v", err)
	}

	// Test secret mount preparation
	ctx := context.Background()
	mountPath, cleanup, err := injector.PrepareSecretMount(ctx, "mount-secret")
	if err != nil {
		t.Errorf("Failed to prepare secret mount: %v", err)
	}
	defer cleanup()

	// Verify file exists and has correct content
	content, err := os.ReadFile(mountPath)
	if err != nil {
		t.Errorf("Failed to read mounted secret file: %v", err)
	}

	if string(content) != "secret-file-content" {
		t.Error("Mounted file content doesn't match secret value")
	}

	// Verify file permissions
	info, err := os.Stat(mountPath)
	if err != nil {
		t.Errorf("Failed to stat mounted file: %v", err)
	}

	if info.Mode().Perm() != 0600 {
		t.Errorf("Mounted file has wrong permissions: got %o, want %o", info.Mode().Perm(), 0600)
	}

	// Test cleanup
	cleanup()

	// Verify file is cleaned up
	if _, err := os.Stat(mountPath); !os.IsNotExist(err) {
		t.Error("Mounted file not cleaned up after cleanup()")
	}
}

func TestInjectorContextCancellation(t *testing.T) {
	vault, _ := NewVault()
	sanitizer := NewSanitizer()
	injector := NewInjector(vault, sanitizer)
	defer vault.Clear()

	// Create cancelled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	// Test that injection respects context cancellation
	_, err := injector.InjectBuildArgs(ctx, nil, []string{"nonexistent"})
	if err == nil {
		t.Error("Expected error due to cancelled context")
	}
	if !strings.Contains(err.Error(), "context cancelled") {
		t.Errorf("Expected context cancellation error, got: %v", err)
	}
}

func TestSecurityValidations(t *testing.T) {
	t.Run("Empty secret ID validation", func(t *testing.T) {
		vault, _ := NewVault()
		defer vault.Clear()

		secret := &Secret{
			ID:    "", // Empty ID should be rejected
			Type:  SecretTypeBuildArg,
			Value: []byte("test"),
		}

		if err := vault.Store(secret); err == nil {
			t.Error("Expected error for empty secret ID")
		}
	})

	t.Run("Nil secret validation", func(t *testing.T) {
		vault, _ := NewVault()
		defer vault.Clear()

		if err := vault.Store(nil); err == nil {
			t.Error("Expected error for nil secret")
		}
	})

	t.Run("Duplicate secret ID validation", func(t *testing.T) {
		vault, _ := NewVault()
		defer vault.Clear()

		secret1 := &Secret{ID: "duplicate", Type: SecretTypeBuildArg, Value: []byte("value1")}
		secret2 := &Secret{ID: "duplicate", Type: SecretTypeMount, Value: []byte("value2")}

		if err := vault.Store(secret1); err != nil {
			t.Errorf("Failed to store first secret: %v", err)
		}

		if err := vault.Store(secret2); err == nil {
			t.Error("Expected error for duplicate secret ID")
		}
	})
}