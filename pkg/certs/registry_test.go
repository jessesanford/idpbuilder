package certs

import (
	"strings"
	"testing"
)

func TestRegistryConfigManager_UpdateInsecureRegistry(t *testing.T) {
	tempDir := t.TempDir()
	config := TrustStoreConfig{
		Location:        UserTrustStore,
		BaseDir:         tempDir,
		DirPermissions:  0755,
		FilePermissions: 0600,
	}
	
	mgr := newRegistryConfigManager(config)
	
	// Add insecure registry
	err := mgr.UpdateInsecureRegistry("insecure.example.com", true)
	if err != nil {
		t.Fatalf("Failed to add insecure registry: %v", err)
	}
	
	// Check it was added
	insecureRegs, err := mgr.GetInsecureRegistries()
	if err != nil {
		t.Fatalf("Failed to get insecure registries: %v", err)
	}
	
	if len(insecureRegs) != 1 || insecureRegs[0] != "insecure.example.com" {
		t.Errorf("Expected [insecure.example.com], got %v", insecureRegs)
	}
	
	// Remove insecure registry
	err = mgr.UpdateInsecureRegistry("insecure.example.com", false)
	if err != nil {
		t.Fatalf("Failed to remove insecure registry: %v", err)
	}
	
	// Check it was removed
	insecureRegs, err = mgr.GetInsecureRegistries()
	if err != nil {
		t.Fatalf("Failed to get insecure registries: %v", err)
	}
	
	if len(insecureRegs) != 0 {
		t.Errorf("Expected empty list, got %v", insecureRegs)
	}
}

func TestRegistryConfigManager_SaveAndLoadConfig(t *testing.T) {
	tempDir := t.TempDir()
	config := TrustStoreConfig{
		Location:        UserTrustStore,
		BaseDir:         tempDir,
		DirPermissions:  0755,
		FilePermissions: 0600,
	}
	
	mgr := newRegistryConfigManager(config)
	
	// Add some insecure registries
	err := mgr.UpdateInsecureRegistry("reg1.example.com", true)
	if err != nil {
		t.Fatalf("Failed to add registry 1: %v", err)
	}
	
	err = mgr.UpdateInsecureRegistry("reg2.example.com", true)
	if err != nil {
		t.Fatalf("Failed to add registry 2: %v", err)
	}
	
	// Create a new manager to test loading
	mgr2 := newRegistryConfigManager(config)
	
	// Check that configuration was loaded
	insecureRegs, err := mgr2.GetInsecureRegistries()
	if err != nil {
		t.Fatalf("Failed to get insecure registries: %v", err)
	}
	
	if len(insecureRegs) != 2 {
		t.Errorf("Expected 2 registries, got %d", len(insecureRegs))
	}
}

func TestRegistryConfigManager_ParseInsecureRegistries(t *testing.T) {
	tempDir := t.TempDir()
	config := TrustStoreConfig{
		Location:        UserTrustStore,
		BaseDir:         tempDir,
		DirPermissions:  0755,
		FilePermissions: 0600,
	}
	
	mgr := &registryConfigManager{config: config}
	
	testConfig := `# Test configuration
[registries.insecure]
registries = ["insecure1.example.com", "insecure2.example.com"]

[registries.search]
registries = ["docker.io"]
`
	
	result := mgr.parseInsecureRegistries(testConfig)
	
	expected := []string{"insecure1.example.com", "insecure2.example.com"}
	if len(result) != len(expected) {
		t.Errorf("Expected %d registries, got %d", len(expected), len(result))
	}
	
	for i, reg := range expected {
		if i >= len(result) || result[i] != reg {
			t.Errorf("Expected registry %s at position %d, got %v", reg, i, result)
		}
	}
}

func TestRegistryConfigManager_GenerateRegistryConfig(t *testing.T) {
	tempDir := t.TempDir()
	config := TrustStoreConfig{
		Location:        UserTrustStore,
		BaseDir:         tempDir,
		DirPermissions:  0755,
		FilePermissions: 0600,
	}
	
	mgr := &registryConfigManager{
		config:             config,
		insecureRegistries: []string{"insecure1.example.com", "insecure2.example.com"},
	}
	
	result := mgr.generateRegistryConfig()
	
	// Check that the config contains the expected content
	if !strings.Contains(result, "[registries.insecure]") {
		t.Error("Expected config to contain [registries.insecure] section")
	}
	
	if !strings.Contains(result, `"insecure1.example.com"`) {
		t.Error("Expected config to contain insecure1.example.com")
	}
	
	if !strings.Contains(result, `"insecure2.example.com"`) {
		t.Error("Expected config to contain insecure2.example.com")
	}
}

func TestRegistryConfigManager_LoadNonExistentConfig(t *testing.T) {
	tempDir := t.TempDir()
	config := TrustStoreConfig{
		Location:        UserTrustStore,
		BaseDir:         tempDir,
		DirPermissions:  0755,
		FilePermissions: 0600,
	}
	
	mgr := newRegistryConfigManager(config)
	
	// Loading non-existent config should not error
	err := mgr.LoadConfig()
	if err != nil {
		t.Fatalf("Failed to load non-existent config: %v", err)
	}
	
	// Should have empty registry list
	regs, err := mgr.GetInsecureRegistries()
	if err != nil {
		t.Fatalf("Failed to get registries: %v", err)
	}
	
	if len(regs) != 0 {
		t.Errorf("Expected empty registry list, got %v", regs)
	}
}