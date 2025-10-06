package format

import (
	"testing"
)

// TestNewPackageManifest tests the creation of new package manifests
func TestNewPackageManifest(t *testing.T) {
	tests := []struct {
		name     string
		pkgName  string
		version  string
		expected bool
	}{
		{"Valid manifest", "test-package", "1.0.0", true},
		{"Empty name", "", "1.0.0", false},
		{"Empty version", "test-package", "", false},
		{"Empty both", "", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manifest := NewPackageManifest(tt.pkgName, tt.version)
			if tt.expected && manifest == nil {
				t.Error("Expected valid manifest, got nil")
			}
			if !tt.expected && manifest != nil {
				t.Error("Expected nil manifest, got valid manifest")
			}
			if manifest != nil {
				if manifest.SchemaVersion != SchemaVersion {
					t.Errorf("Expected schema version %d, got %d", SchemaVersion, manifest.SchemaVersion)
				}
				if manifest.MediaType != MediaTypes.ManifestV2 {
					t.Errorf("Expected media type %s, got %s", MediaTypes.ManifestV2, manifest.MediaType)
				}
			}
		})
	}
}

// TestValidateManifest tests manifest validation logic
func TestValidateManifest(t *testing.T) {
	validManifest := &PackageManifest{
		SchemaVersion: SchemaVersion,
		MediaType:     MediaTypes.ManifestV2,
		Config: Descriptor{
			MediaType: MediaTypes.ConfigV1,
			Digest:    "sha256:test",
			Size:      1024,
		},
	}

	tests := []struct {
		name     string
		manifest *PackageManifest
		valid    bool
	}{
		{"Valid manifest", validManifest, true},
		{"Nil manifest", nil, false},
		{"Invalid schema version", &PackageManifest{SchemaVersion: 1}, false},
		{"Empty media type", &PackageManifest{SchemaVersion: SchemaVersion}, false},
		{"Invalid config", &PackageManifest{SchemaVersion: SchemaVersion, MediaType: MediaTypes.ManifestV2}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateManifest(tt.manifest)
			if tt.valid && err != nil {
				t.Errorf("Expected valid manifest, got error: %v", err)
			}
			if !tt.valid && err == nil {
				t.Error("Expected validation error, got nil")
			}
		})
	}
}

// TestAddLayer tests adding layers to manifests
func TestAddLayer(t *testing.T) {
	manifest := NewPackageManifest("test", "1.0.0")
	validLayer := &LayerDescriptor{
		Descriptor: Descriptor{
			MediaType: MediaTypes.LayerTarGzip,
			Digest:    "sha256:layer",
			Size:      2048,
		},
	}

	tests := []struct {
		name     string
		manifest *PackageManifest
		layer    *LayerDescriptor
		valid    bool
	}{
		{"Valid layer", manifest, validLayer, true},
		{"Nil manifest", nil, validLayer, false},
		{"Nil layer", manifest, nil, false},
		{"Invalid layer", manifest, &LayerDescriptor{}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := AddLayer(tt.manifest, tt.layer)
			if tt.valid && err != nil {
				t.Errorf("Expected success, got error: %v", err)
			}
			if !tt.valid && err == nil {
				t.Error("Expected error, got success")
			}
		})
	}
}

// TestSerializeParseManifest tests serialization and parsing of manifests
func TestSerializeParseManifest(t *testing.T) {
	manifest := NewPackageManifest("test-package", "1.0.0")
	manifest.Config = Descriptor{
		MediaType: MediaTypes.ConfigV1,
		Digest:    "sha256:config",
		Size:      512,
	}

	// Test serialization
	data, err := SerializeManifest(manifest)
	if err != nil {
		t.Fatalf("Serialization failed: %v", err)
	}

	if len(data) == 0 {
		t.Error("Expected serialized data, got empty")
	}

	// Test parsing
	parsed, err := ParseManifest(data)
	if err != nil {
		t.Fatalf("Parsing failed: %v", err)
	}

	if parsed.SchemaVersion != manifest.SchemaVersion {
		t.Errorf("Schema version mismatch: expected %d, got %d", manifest.SchemaVersion, parsed.SchemaVersion)
	}

	if parsed.MediaType != manifest.MediaType {
		t.Errorf("Media type mismatch: expected %s, got %s", manifest.MediaType, parsed.MediaType)
	}

	// Test error cases
	_, err = SerializeManifest(nil)
	if err == nil {
		t.Error("Expected error for nil manifest serialization")
	}

	_, err = ParseManifest([]byte{})
	if err == nil {
		t.Error("Expected error for empty data parsing")
	}
}