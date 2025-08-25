package oci

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestManifestReader_ReadManifest(t *testing.T) {
	validManifest := &OCIManifest{
		SchemaVersion: 2,
		MediaType:     MediaTypeManifest,
		Config: OCIDescriptor{
			MediaType: MediaTypeImageConfig,
			Digest:    "sha256:config123abc",
			Size:      1234,
		},
		Layers: []OCIDescriptor{
			{
				MediaType: MediaTypeImageLayer,
				Digest:    "sha256:layer123abc",
				Size:      5678,
			},
		},
		Annotations: map[string]string{
			"org.opencontainers.image.created": "2023-01-01T00:00:00Z",
		},
	}

	validJSON, _ := json.Marshal(validManifest)

	tests := []struct {
		name    string
		input   string
		wantErr bool
		errMsg  string
	}{
		{
			name:  "valid manifest",
			input: string(validJSON),
		},
		{
			name:    "empty input",
			input:   "",
			wantErr: true,
			errMsg:  "empty manifest",
		},
		{
			name:    "invalid JSON",
			input:   `{"schemaVersion": 2, "mediaType": "invalid json"`,
			wantErr: true,
			errMsg:  "failed to parse manifest JSON",
		},
		{
			name: "invalid manifest - no layers",
			input: `{
				"schemaVersion": 2,
				"mediaType": "application/vnd.oci.image.manifest.v1+json",
				"config": {
					"mediaType": "application/vnd.oci.image.config.v1+json",
					"digest": "sha256:config123",
					"size": 1234
				},
				"layers": []
			}`,
			wantErr: true,
			errMsg:  "invalid manifest",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := NewManifestReader()
			result, err := reader.ReadManifest(strings.NewReader(tt.input))
			
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadManifest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			
			if err != nil && tt.errMsg != "" {
				if !contains(err.Error(), tt.errMsg) {
					t.Errorf("ReadManifest() error = %q, want to contain %q", err.Error(), tt.errMsg)
				}
				return
			}
			
			if result != nil {
				if result.SchemaVersion != 2 {
					t.Errorf("Expected schema version 2, got %d", result.SchemaVersion)
				}
				if result.MediaType != MediaTypeManifest {
					t.Errorf("Expected mediaType %s, got %s", MediaTypeManifest, result.MediaType)
				}
			}
		})
	}
}

func TestParseManifestJSON(t *testing.T) {
	validManifest := &OCIManifest{
		SchemaVersion: 2,
		MediaType:     MediaTypeManifest,
		Config: OCIDescriptor{
			MediaType: MediaTypeImageConfig,
			Digest:    "sha256:config123",
			Size:      1234,
		},
		Layers: []OCIDescriptor{
			{
				MediaType: MediaTypeImageLayer,
				Digest:    "sha256:layer123",
				Size:      5678,
			},
		},
	}

	validJSON, _ := json.Marshal(validManifest)

	tests := []struct {
		name    string
		data    []byte
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid manifest",
			data: validJSON,
		},
		{
			name:    "empty data",
			data:    []byte{},
			wantErr: true,
			errMsg:  "empty manifest data",
		},
		{
			name:    "too large",
			data:    make([]byte, MaxManifestSize+1),
			wantErr: true,
			errMsg:  "manifest size",
		},
		{
			name:    "invalid JSON",
			data:    []byte(`{invalid json`),
			wantErr: true,
			errMsg:  "failed to parse manifest JSON",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseManifestJSON(tt.data)
			
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseManifestJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			
			if err != nil && tt.errMsg != "" {
				if !contains(err.Error(), tt.errMsg) {
					t.Errorf("ParseManifestJSON() error = %q, want to contain %q", err.Error(), tt.errMsg)
				}
				return
			}
			
			if result != nil && result.SchemaVersion != 2 {
				t.Errorf("Expected schema version 2, got %d", result.SchemaVersion)
			}
		})
	}
}

func TestManifestWriter_WriteManifest(t *testing.T) {
	validManifest := &OCIManifest{
		SchemaVersion: 2,
		MediaType:     MediaTypeManifest,
		Config: OCIDescriptor{
			MediaType: MediaTypeImageConfig,
			Digest:    "sha256:config123",
			Size:      1234,
		},
		Layers: []OCIDescriptor{
			{
				MediaType: MediaTypeImageLayer,
				Digest:    "sha256:layer123",
				Size:      5678,
			},
		},
		Annotations: map[string]string{
			"test": "annotation",
		},
	}

	tests := []struct {
		name     string
		manifest *OCIManifest
		wantErr  bool
		errMsg   string
	}{
		{
			name:     "valid manifest",
			manifest: validManifest,
		},
		{
			name:     "nil manifest",
			manifest: nil,
			wantErr:  true,
			errMsg:   "manifest cannot be nil",
		},
		{
			name: "invalid manifest",
			manifest: &OCIManifest{
				SchemaVersion: 1, // Invalid
				MediaType:     MediaTypeManifest,
			},
			wantErr: true,
			errMsg:  "invalid manifest",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writer := NewManifestWriter()
			var buf bytes.Buffer
			
			err := writer.WriteManifest(&buf, tt.manifest)
			
			if (err != nil) != tt.wantErr {
				t.Errorf("WriteManifest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			
			if err != nil && tt.errMsg != "" {
				if !contains(err.Error(), tt.errMsg) {
					t.Errorf("WriteManifest() error = %q, want to contain %q", err.Error(), tt.errMsg)
				}
				return
			}
			
			if !tt.wantErr && buf.Len() == 0 {
				t.Errorf("WriteManifest() produced empty output")
			}
		})
	}
}

func TestMarshalManifest(t *testing.T) {
	validManifest := &OCIManifest{
		SchemaVersion: 2,
		MediaType:     MediaTypeManifest,
		Config: OCIDescriptor{
			MediaType: MediaTypeImageConfig,
			Digest:    "sha256:config123",
			Size:      1234,
		},
		Layers: []OCIDescriptor{
			{
				MediaType: MediaTypeImageLayer,
				Digest:    "sha256:layer123",
				Size:      5678,
			},
		},
	}

	tests := []struct {
		name     string
		manifest *OCIManifest
		wantErr  bool
		errMsg   string
	}{
		{
			name:     "valid manifest",
			manifest: validManifest,
		},
		{
			name:     "nil manifest",
			manifest: nil,
			wantErr:  true,
			errMsg:   "manifest cannot be nil",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := MarshalManifest(tt.manifest)
			
			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalManifest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			
			if err != nil && tt.errMsg != "" {
				if !contains(err.Error(), tt.errMsg) {
					t.Errorf("MarshalManifest() error = %q, want to contain %q", err.Error(), tt.errMsg)
				}
				return
			}
			
			if !tt.wantErr {
				if len(data) == 0 {
					t.Errorf("MarshalManifest() produced empty data")
				}
				
				// Verify we can parse it back
				var parsed OCIManifest
				if err := json.Unmarshal(data, &parsed); err != nil {
					t.Errorf("MarshalManifest() produced invalid JSON: %v", err)
				}
			}
		})
	}
}

func TestDetectManifestMediaType(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		want    string
		wantErr bool
		errMsg  string
	}{
		{
			name: "OCI manifest",
			data: []byte(`{"mediaType": "application/vnd.oci.image.manifest.v1+json"}`),
			want: MediaTypeManifest,
		},
		{
			name: "OCI image index",
			data: []byte(`{"mediaType": "application/vnd.oci.image.index.v1+json"}`),
			want: MediaTypeImageIndex,
		},
		{
			name: "Docker manifest",
			data: []byte(`{"mediaType": "application/vnd.docker.distribution.manifest.v2+json"}`),
			want: DockerMediaTypeManifest,
		},
		{
			name:    "empty data",
			data:    []byte{},
			wantErr: true,
			errMsg:  "empty data",
		},
		{
			name:    "invalid JSON",
			data:    []byte(`{invalid`),
			wantErr: true,
			errMsg:  "not valid JSON",
		},
		{
			name:    "missing mediaType",
			data:    []byte(`{"schemaVersion": 2}`),
			wantErr: true,
			errMsg:  "mediaType field is missing",
		},
		{
			name:    "unknown mediaType",
			data:    []byte(`{"mediaType": "unknown/type"}`),
			wantErr: true,
			errMsg:  "unknown manifest media type",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DetectManifestMediaType(tt.data)
			
			if (err != nil) != tt.wantErr {
				t.Errorf("DetectManifestMediaType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			
			if err != nil && tt.errMsg != "" {
				if !contains(err.Error(), tt.errMsg) {
					t.Errorf("DetectManifestMediaType() error = %q, want to contain %q", err.Error(), tt.errMsg)
				}
				return
			}
			
			if got != tt.want {
				t.Errorf("DetectManifestMediaType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsOCIManifest(t *testing.T) {
	tests := []struct {
		name string
		data []byte
		want bool
	}{
		{
			name: "OCI manifest",
			data: []byte(`{"mediaType": "application/vnd.oci.image.manifest.v1+json"}`),
			want: true,
		},
		{
			name: "OCI image index",
			data: []byte(`{"mediaType": "application/vnd.oci.image.index.v1+json"}`),
			want: true,
		},
		{
			name: "Docker manifest",
			data: []byte(`{"mediaType": "application/vnd.docker.distribution.manifest.v2+json"}`),
			want: false,
		},
		{
			name: "invalid data",
			data: []byte(`invalid`),
			want: false,
		},
		{
			name: "empty data",
			data: []byte{},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsOCIManifest(tt.data)
			if got != tt.want {
				t.Errorf("IsOCIManifest() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsDockerManifest(t *testing.T) {
	tests := []struct {
		name string
		data []byte
		want bool
	}{
		{
			name: "Docker manifest",
			data: []byte(`{"mediaType": "application/vnd.docker.distribution.manifest.v2+json"}`),
			want: true,
		},
		{
			name: "Docker manifest list",
			data: []byte(`{"mediaType": "application/vnd.docker.distribution.manifest.list.v2+json"}`),
			want: true,
		},
		{
			name: "OCI manifest",
			data: []byte(`{"mediaType": "application/vnd.oci.image.manifest.v1+json"}`),
			want: false,
		},
		{
			name: "invalid data",
			data: []byte(`invalid`),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsDockerManifest(tt.data)
			if got != tt.want {
				t.Errorf("IsDockerManifest() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNormalizeManifestMediaType(t *testing.T) {
	tests := []struct {
		name      string
		mediaType string
		want      string
	}{
		{
			name:      "Docker manifest to OCI",
			mediaType: DockerMediaTypeManifest,
			want:      MediaTypeManifest,
		},
		{
			name:      "Docker manifest list to OCI index",
			mediaType: DockerMediaTypeManifestList,
			want:      MediaTypeImageIndex,
		},
		{
			name:      "Docker config to OCI config",
			mediaType: DockerMediaTypeConfig,
			want:      MediaTypeImageConfig,
		},
		{
			name:      "Docker layer to OCI layer",
			mediaType: DockerMediaTypeLayer,
			want:      MediaTypeImageLayer,
		},
		{
			name:      "OCI manifest unchanged",
			mediaType: MediaTypeManifest,
			want:      MediaTypeManifest,
		},
		{
			name:      "Unknown type unchanged",
			mediaType: "unknown/type",
			want:      "unknown/type",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NormalizeManifestMediaType(tt.mediaType)
			if got != tt.want {
				t.Errorf("NormalizeManifestMediaType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidateDigest(t *testing.T) {
	tests := []struct {
		name    string
		digest  string
		wantErr bool
		errMsg  string
	}{
		{
			name:   "valid sha256",
			digest: "sha256:abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789",
		},
		{
			name:   "valid sha512",
			digest: "sha512:abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789",
		},
		{
			name:    "empty digest",
			digest:  "",
			wantErr: true,
			errMsg:  "digest cannot be empty",
		},
		{
			name:    "no colon",
			digest:  "sha256abc123",
			wantErr: true,
			errMsg:  "digest must contain algorithm and hex",
		},
		{
			name:    "empty algorithm",
			digest:  ":abc123",
			wantErr: true,
			errMsg:  "digest algorithm cannot be empty",
		},
		{
			name:    "empty hex",
			digest:  "sha256:",
			wantErr: true,
			errMsg:  "digest hex cannot be empty",
		},
		{
			name:    "sha256 wrong length",
			digest:  "sha256:abc123",
			wantErr: true,
			errMsg:  "sha256 digest must be 64 hex characters",
		},
		{
			name:    "invalid hex character",
			digest:  "sha256:abcdef0123456789abcdef0123456789abcdef0123456789abcdef012345xy89",
			wantErr: true,
			errMsg:  "digest hex contains invalid character",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateDigest(tt.digest)
			
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateDigest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			
			if err != nil && tt.errMsg != "" {
				if !contains(err.Error(), tt.errMsg) {
					t.Errorf("ValidateDigest() error = %q, want to contain %q", err.Error(), tt.errMsg)
				}
			}
		})
	}
}

func TestCreateManifest(t *testing.T) {
	config := OCIDescriptor{
		MediaType: MediaTypeImageConfig,
		Digest:    "sha256:config123",
		Size:      1234,
	}
	
	layers := []OCIDescriptor{
		{
			MediaType: MediaTypeImageLayer,
			Digest:    "sha256:layer123",
			Size:      5678,
		},
	}
	
	annotations := map[string]string{
		"test": "annotation",
	}
	
	manifest := CreateManifest(config, layers, annotations)
	
	if manifest.SchemaVersion != 2 {
		t.Errorf("Expected schema version 2, got %d", manifest.SchemaVersion)
	}
	
	if manifest.MediaType != MediaTypeManifest {
		t.Errorf("Expected mediaType %s, got %s", MediaTypeManifest, manifest.MediaType)
	}
	
	if manifest.Config.Digest != config.Digest {
		t.Errorf("Config digest mismatch: expected %s, got %s", config.Digest, manifest.Config.Digest)
	}
	
	if len(manifest.Layers) != len(layers) {
		t.Errorf("Expected %d layers, got %d", len(layers), len(manifest.Layers))
	}
	
	if manifest.Annotations == nil {
		t.Errorf("Annotations should not be nil")
	}
	
	if manifest.Annotations["test"] != "annotation" {
		t.Errorf("Annotation mismatch: expected 'annotation', got %s", manifest.Annotations["test"])
	}
}

func TestCreateManifest_NilAnnotations(t *testing.T) {
	config := OCIDescriptor{
		MediaType: MediaTypeImageConfig,
		Digest:    "sha256:config123",
		Size:      1234,
	}
	
	layers := []OCIDescriptor{
		{
			MediaType: MediaTypeImageLayer,
			Digest:    "sha256:layer123",
			Size:      5678,
		},
	}
	
	manifest := CreateManifest(config, layers, nil)
	
	if manifest.Annotations == nil {
		t.Errorf("Annotations should be initialized even when nil is passed")
	}
}