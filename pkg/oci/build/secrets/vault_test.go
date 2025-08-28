package secrets

import (
	"context"
	"os"
	"strings"
	"testing"
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
	secret := &Secret{ID: "test", Type: SecretTypeMount, Value: []byte("secret-data")}
	vault.Store(secret)
	
	// Verify encryption by checking internal storage doesn't contain plaintext
	memVault := vault.(*memoryVault)
	memVault.mu.RLock()
	encSecret := memVault.secrets["test"]
	memVault.mu.RUnlock()
	if strings.Contains(string(encSecret.encrypted), "secret-data") {
		t.Error("SECURITY VIOLATION: Secret stored in plaintext")
	}
}

func TestSanitizerSecurity(t *testing.T) {
	sanitizer := NewSanitizer()
	
	// Test registered secrets
	sanitizer.RegisterSecret("api-key", []byte("sk-1234567890abcdef"))
	input := "Key value is sk-1234567890abcdef"
	output := sanitizer.Sanitize(input)
	if strings.Contains(output, "sk-1234567890abcdef") {
		t.Error("SECURITY VIOLATION: Secret not sanitized")
	}
	if !strings.Contains(output, "***API-KEY***") {
		t.Error("Secret not replaced with mask")
	}
	
	// Test common patterns
	tests := []struct{ input, forbidden, required string }{
		{"password=secret123", "secret123", "***REDACTED***"},
		{"Bearer token123", "token123", "***TOKEN***"},
	}
	for _, test := range tests {
		out := sanitizer.Sanitize(test.input)
		if strings.Contains(out, test.forbidden) {
			t.Errorf("Forbidden pattern not removed: %s", test.forbidden)
		}
		if !strings.Contains(out, test.required) {
			t.Errorf("Required pattern missing: %s", test.required)
		}
	}
}

func TestInjectorOperations(t *testing.T) {
	vault, _ := NewVault()
	sanitizer := NewSanitizer()
	injector := NewInjector(vault, sanitizer)
	defer vault.Clear()

	// Test build arg injection
	secret := &Secret{ID: "test", Type: SecretTypeBuildArg, Value: []byte("secret"), Target: "ARG"}
	vault.Store(secret)
	ctx := context.Background()
	args, err := injector.InjectBuildArgs(ctx, map[string]string{"VER": "1.0"}, []string{"test"})
	if err != nil || args["ARG"] != "secret" || args["VER"] != "1.0" {
		t.Error("Build arg injection failed")
	}

	// Test secret mount
	secret2 := &Secret{ID: "mount", Type: SecretTypeMount, Value: []byte("content"), Mode: 0600}
	vault.Store(secret2)
	path, cleanup, err := injector.PrepareSecretMount(ctx, "mount")
	if err != nil {
		t.Error("Mount preparation failed")
	}
	defer cleanup()
	
	content, _ := os.ReadFile(path)
	if string(content) != "content" {
		t.Error("Mount content mismatch")
	}

	// Test context cancellation
	ctx2, cancel := context.WithCancel(context.Background())
	cancel()
	_, err = injector.InjectBuildArgs(ctx2, nil, []string{"test"})
	if err == nil || !strings.Contains(err.Error(), "context cancelled") {
		t.Error("Context cancellation not handled")
	}
}

func TestSecurityValidations(t *testing.T) {
	vault, _ := NewVault()
	defer vault.Clear()

	// Test empty ID rejection
	if err := vault.Store(&Secret{ID: "", Type: SecretTypeBuildArg, Value: []byte("test")}); err == nil {
		t.Error("Expected error for empty ID")
	}

	// Test nil secret rejection
	if err := vault.Store(nil); err == nil {
		t.Error("Expected error for nil secret")
	}

	// Test duplicate ID rejection
	secret1 := &Secret{ID: "dup", Type: SecretTypeBuildArg, Value: []byte("v1")}
	secret2 := &Secret{ID: "dup", Type: SecretTypeMount, Value: []byte("v2")}
	if vault.Store(secret1) != nil || vault.Store(secret2) == nil {
		t.Error("Duplicate ID validation failed")
	}
}