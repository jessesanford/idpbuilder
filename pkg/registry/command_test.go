package registry

import (
	"context"
	"testing"

	"github.com/cnoe-io/idpbuilder/pkg/cmd/interfaces"
	"github.com/cnoe-io/idpbuilder/pkg/config"
)

func TestNewCommand(t *testing.T) {
	tests := []struct {
		name        string
		config      *config.RegistryConfig
		expectError bool
		errorMsg    string
	}{
		{
			name:        "nil config",
			config:      nil,
			expectError: true,
			errorMsg:    "registry configuration cannot be nil",
		},
		{
			name: "valid config",
			config: &config.RegistryConfig{
				URL: "https://registry.example.com",
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd, err := NewCommand(tt.config)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
					return
				}
				if err.Error() != tt.errorMsg {
					t.Errorf("expected error %q, got %q", tt.errorMsg, err.Error())
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if cmd == nil {
				t.Error("expected non-nil command")
				return
			}

			// Verify command implements RegistryCommand interface
			var _ interfaces.RegistryCommand = cmd

			// Verify internal state
			if cmd.client == nil {
				t.Error("client should be initialized")
			}

			if cmd.config != tt.config {
				t.Error("config should match input")
			}
		})
	}
}

func TestCommand_ConfigureRegistry(t *testing.T) {
	// Create initial command
	initialConfig := &config.RegistryConfig{
		URL: "https://initial.example.com",
	}

	cmd, err := NewCommand(initialConfig)
	if err != nil {
		t.Fatalf("failed to create command: %v", err)
	}

	// Test configuration update
	newConfig := interfaces.RegistryConfig{
		Name:     "test-registry",
		URL:      "https://new.example.com",
		Type:     "generic",
		AuthType: "basic",
		Insecure: true,
	}

	err = cmd.ConfigureRegistry(newConfig)
	if err != nil {
		t.Errorf("ConfigureRegistry() error = %v", err)
		return
	}

	// Verify configuration was updated
	if cmd.config.URL != newConfig.URL {
		t.Errorf("expected URL %q, got %q", newConfig.URL, cmd.config.URL)
	}

	if cmd.config.Type != newConfig.Type {
		t.Errorf("expected Type %q, got %q", newConfig.Type, cmd.config.Type)
	}

	if cmd.config.Insecure != newConfig.Insecure {
		t.Errorf("expected Insecure %t, got %t", newConfig.Insecure, cmd.config.Insecure)
	}

	// Verify client was recreated
	if cmd.client == nil {
		t.Error("client should still be initialized after configuration")
	}
}

func TestCommand_GetRegistryInfo(t *testing.T) {
	config := &config.RegistryConfig{
		URL:  "https://registry.example.com",
		Type: "generic",
		Auth: config.AuthConfig{
			Username: "testuser",
		},
	}

	cmd, err := NewCommand(config)
	if err != nil {
		t.Fatalf("failed to create command: %v", err)
	}

	ctx := context.Background()
	info, err := cmd.GetRegistryInfo(ctx)
	if err != nil {
		t.Errorf("GetRegistryInfo() error = %v", err)
		return
	}

	if info == nil {
		t.Error("expected non-nil registry info")
		return
	}

	// Verify info fields
	if info.Name != config.Type {
		t.Errorf("expected Name %q, got %q", config.Type, info.Name)
	}

	if info.URL != config.URL {
		t.Errorf("expected URL %q, got %q", config.URL, info.URL)
	}

	if len(info.Capabilities) == 0 {
		t.Error("expected non-empty capabilities")
	}

	// Check for expected capabilities
	expectedCaps := []string{"push", "pull", "list", "delete"}
	for _, expectedCap := range expectedCaps {
		found := false
		for _, cap := range info.Capabilities {
			if cap == expectedCap {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected capability %q not found", expectedCap)
		}
	}

	// Verify metadata
	if info.Metadata == nil {
		t.Error("expected non-nil metadata")
	}

	if authenticated, exists := info.Metadata["authenticated"]; !exists {
		t.Error("expected authenticated metadata field")
	} else if authenticated != "true" {
		t.Errorf("expected authenticated=true, got %q", authenticated)
	}
}

func TestCommand_GetCapabilities(t *testing.T) {
	config := &config.RegistryConfig{
		URL: "https://registry.example.com",
	}

	cmd, err := NewCommand(config)
	if err != nil {
		t.Fatalf("failed to create command: %v", err)
	}

	ctx := context.Background()
	capabilities, err := cmd.GetCapabilities(ctx)
	if err != nil {
		t.Errorf("GetCapabilities() error = %v", err)
		return
	}

	if len(capabilities) == 0 {
		t.Error("expected non-empty capabilities")
		return
	}

	// Check for core capabilities
	expectedCaps := []string{"push", "pull", "list", "delete"}
	for _, expectedCap := range expectedCaps {
		found := false
		for _, cap := range capabilities {
			if cap == expectedCap {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected capability %q not found", expectedCap)
		}
	}
}

func TestCommand_NilClient(t *testing.T) {
	cmd := &Command{
		client: nil,
		config: &config.RegistryConfig{
			URL: "https://registry.example.com",
		},
	}

	ctx := context.Background()

	// Test all methods that should fail with nil client
	t.Run("ValidateCredentials", func(t *testing.T) {
		err := cmd.ValidateCredentials(ctx)
		if err == nil {
			t.Error("expected error with nil client")
		}
		if err.Error() != "registry client not configured" {
			t.Errorf("expected 'registry client not configured', got %q", err.Error())
		}
	})

	t.Run("GetRegistryInfo", func(t *testing.T) {
		_, err := cmd.GetRegistryInfo(ctx)
		if err == nil {
			t.Error("expected error with nil client")
		}
		if err.Error() != "registry client not configured" {
			t.Errorf("expected 'registry client not configured', got %q", err.Error())
		}
	})

	t.Run("TestConnection", func(t *testing.T) {
		err := cmd.TestConnection(ctx)
		if err == nil {
			t.Error("expected error with nil client")
		}
		if err.Error() != "registry client not configured" {
			t.Errorf("expected 'registry client not configured', got %q", err.Error())
		}
	})

	t.Run("GetCapabilities", func(t *testing.T) {
		_, err := cmd.GetCapabilities(ctx)
		if err == nil {
			t.Error("expected error with nil client")
		}
		if err.Error() != "registry client not configured" {
			t.Errorf("expected 'registry client not configured', got %q", err.Error())
		}
	})
}

// TestCommand_RegistryCommandInterface verifies that Command implements the RegistryCommand interface
func TestCommand_RegistryCommandInterface(t *testing.T) {
	var _ interfaces.RegistryCommand = (*Command)(nil)
}