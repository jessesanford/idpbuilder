package unit

import (
	"os"
	"testing"

	"github.com/cnoe-io/idpbuilder/pkg/oci/config"
)

func TestGetDefault(t *testing.T) {
	// Test default configuration
	cfg := config.GetDefault()
	
	if cfg.RegistryURL != "gitea.idpbuilder.localtest.me" {
		t.Errorf("Expected default registry URL, got %s", cfg.RegistryURL)
	}
	
	if !cfg.AutoExtractCerts {
		t.Errorf("Expected AutoExtractCerts to be true by default")
	}
	
	if cfg.InsecureMode {
		t.Errorf("Expected InsecureMode to be false by default")
	}
}

func TestEnvironmentOverrides(t *testing.T) {
	// Set environment variables
	os.Setenv("IDPBUILDER_OCI_REGISTRY", "custom.registry.com")
	os.Setenv("IDPBUILDER_OCI_INSECURE", "true")
	defer func() {
		os.Unsetenv("IDPBUILDER_OCI_REGISTRY")
		os.Unsetenv("IDPBUILDER_OCI_INSECURE")
	}()
	
	cfg := config.GetDefault()
	
	if cfg.RegistryURL != "custom.registry.com" {
		t.Errorf("Expected custom registry URL, got %s", cfg.RegistryURL)
	}
	
	if !cfg.InsecureMode {
		t.Errorf("Expected InsecureMode to be true from env var")
	}
	
	if cfg.AutoExtractCerts {
		t.Errorf("Expected AutoExtractCerts to be false in insecure mode")
	}
}