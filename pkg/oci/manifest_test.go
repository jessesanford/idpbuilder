package oci

import (
	"encoding/json"
	"testing"
)

func TestOCIDescriptor_IsValid(t *testing.T) {
	tests := []struct {
		name       string
		descriptor OCIDescriptor
		expectErr  bool
	}{
		{
			name: "valid descriptor",
			descriptor: OCIDescriptor{
				MediaType: MediaTypeLayer,
				Digest:    "sha256:abc123",
				Size:      1024,
			},
			expectErr: false,
		},
		{
			name: "missing media type",
			descriptor: OCIDescriptor{
				Digest: "sha256:abc123",
				Size:   1024,
			},
			expectErr: true,
		},
		{
			name: "missing digest",
			descriptor: OCIDescriptor{
				MediaType: MediaTypeLayer,
				Size:      1024,
			},
			expectErr: true,
		},
		{
			name: "negative size",
			descriptor: OCIDescriptor{
				MediaType: MediaTypeLayer,
				Digest:    "sha256:abc123",
				Size:      -1,
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.descriptor.IsValid()
			if tt.expectErr && err == nil {
				t.Error("expected error but got none")
			}
			if !tt.expectErr && err != nil {
				t.Errorf("expected no error but got: %v", err)
			}
		})
	}
}

func TestOCIManifest_IsValid(t *testing.T) {
	validDescriptor := OCIDescriptor{
		MediaType: MediaTypeLayer,
		Digest:    "sha256:abc123",
		Size:      1024,
	}

	tests := []struct {
		name      string
		manifest  OCIManifest
		expectErr bool
	}{
		{
			name: "valid manifest",
			manifest: OCIManifest{
				SchemaVersion: 2,
				MediaType:     MediaTypeManifest,
				Config:        validDescriptor,
				Layers:        []OCIDescriptor{validDescriptor},
			},
			expectErr: false,
		},
		{
			name: "invalid schema version",
			manifest: OCIManifest{
				SchemaVersion: 1,
				MediaType:     MediaTypeManifest,
				Config:        validDescriptor,
				Layers:        []OCIDescriptor{validDescriptor},
			},
			expectErr: true,
		},
		{
			name: "missing media type",
			manifest: OCIManifest{
				SchemaVersion: 2,
				Config:        validDescriptor,
				Layers:        []OCIDescriptor{validDescriptor},
			},
			expectErr: true,
		},
		{
			name: "no layers",
			manifest: OCIManifest{
				SchemaVersion: 2,
				MediaType:     MediaTypeManifest,
				Config:        validDescriptor,
				Layers:        []OCIDescriptor{},
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.manifest.IsValid()
			if tt.expectErr && err == nil {
				t.Error("expected error but got none")
			}
			if !tt.expectErr && err != nil {
				t.Errorf("expected no error but got: %v", err)
			}
		})
	}
}

func TestOCIManifest_ToJSON(t *testing.T) {
	manifest := OCIManifest{
		SchemaVersion: 2,
		MediaType:     MediaTypeManifest,
		Config: OCIDescriptor{
			MediaType: MediaTypeConfig,
			Digest:    "sha256:config123",
			Size:      512,
		},
		Layers: []OCIDescriptor{
			{
				MediaType: MediaTypeLayer,
				Digest:    "sha256:layer123",
				Size:      1024,
			},
		},
	}

	data, err := manifest.ToJSON()
	if err != nil {
		t.Fatalf("ToJSON failed: %v", err)
	}

	var parsed OCIManifest
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	if parsed.SchemaVersion != manifest.SchemaVersion {
		t.Errorf("expected schema version %d, got %d", manifest.SchemaVersion, parsed.SchemaVersion)
	}
}

func TestOCIManifest_FromJSON(t *testing.T) {
	jsonData := []byte(`{
		"schemaVersion": 2,
		"mediaType": "application/vnd.oci.image.manifest.v1+json",
		"config": {
			"mediaType": "application/vnd.oci.image.config.v1+json",
			"digest": "sha256:config123",
			"size": 512
		},
		"layers": [
			{
				"mediaType": "application/vnd.oci.image.layer.v1.tar+gzip",
				"digest": "sha256:layer123",
				"size": 1024
			}
		]
	}`)

	var manifest OCIManifest
	err := manifest.FromJSON(jsonData)
	if err != nil {
		t.Fatalf("FromJSON failed: %v", err)
	}

	if manifest.SchemaVersion != 2 {
		t.Errorf("expected schema version 2, got %d", manifest.SchemaVersion)
	}
	if len(manifest.Layers) != 1 {
		t.Errorf("expected 1 layer, got %d", len(manifest.Layers))
	}
}