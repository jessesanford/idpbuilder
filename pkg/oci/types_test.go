package oci

import (
	"testing"
)

func TestOCIReference_String(t *testing.T) {
	tests := []struct {
		name     string
		ref      *OCIReference
		expected string
	}{
		{
			name: "full reference with registry and namespace",
			ref: &OCIReference{
				Registry:   "registry.example.com",
				Namespace:  "myorg",
				Repository: "myapp",
				Tag:        "v1.0.0",
			},
			expected: "registry.example.com/myorg/myapp:v1.0.0",
		},
		{
			name: "default registry with namespace",
			ref: &OCIReference{
				Registry:   DefaultRegistry,
				Namespace:  "myorg",
				Repository: "myapp",
				Tag:        "latest",
			},
			expected: "myorg/myapp:latest",
		},
		{
			name: "default registry and namespace",
			ref: &OCIReference{
				Registry:   DefaultRegistry,
				Namespace:  DefaultNamespace,
				Repository: "nginx",
				Tag:        "alpine",
			},
			expected: "nginx:alpine",
		},
		{
			name: "with digest instead of tag",
			ref: &OCIReference{
				Registry:   "quay.io",
				Namespace:  "prometheus",
				Repository: "prometheus",
				Digest:     "sha256:abc123def456",
			},
			expected: "quay.io/prometheus/prometheus@sha256:abc123def456",
		},
		{
			name: "digest takes precedence over tag",
			ref: &OCIReference{
				Registry:   "gcr.io",
				Namespace:  "project",
				Repository: "image",
				Tag:        "v2.0.0",
				Digest:     "sha256:xyz789",
			},
			expected: "gcr.io/project/image@sha256:xyz789",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.ref.String()
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestOCIReference_IsDigest(t *testing.T) {
	tests := []struct {
		name     string
		ref      *OCIReference
		expected bool
	}{
		{
			name: "has digest",
			ref: &OCIReference{
				Repository: "app",
				Digest:     "sha256:abc123",
			},
			expected: true,
		},
		{
			name: "no digest",
			ref: &OCIReference{
				Repository: "app",
				Tag:        "latest",
			},
			expected: false,
		},
		{
			name: "empty digest",
			ref: &OCIReference{
				Repository: "app",
				Digest:     "",
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.ref.IsDigest()
			if result != tt.expected {
				t.Errorf("expected %t, got %t", tt.expected, result)
			}
		})
	}
}

func TestOCIPlatform_String(t *testing.T) {
	tests := []struct {
		name     string
		platform *OCIPlatform
		expected string
	}{
		{
			name: "basic platform",
			platform: &OCIPlatform{
				OS:           "linux",
				Architecture: "amd64",
			},
			expected: "linux/amd64",
		},
		{
			name: "platform with variant",
			platform: &OCIPlatform{
				OS:           "linux",
				Architecture: "arm",
				Variant:      "v7",
			},
			expected: "linux/arm/v7",
		},
		{
			name: "windows platform",
			platform: &OCIPlatform{
				OS:           "windows",
				Architecture: "amd64",
			},
			expected: "windows/amd64",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.platform.String()
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestParseOCIReference(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *OCIReference
		wantErr bool
	}{
		{
			name:  "simple repository name",
			input: "nginx",
			want: &OCIReference{
				Registry:   DefaultRegistry,
				Namespace:  DefaultNamespace,
				Repository: "nginx",
				Tag:        DefaultTag,
			},
		},
		{
			name:  "repository with tag",
			input: "nginx:alpine",
			want: &OCIReference{
				Registry:   DefaultRegistry,
				Namespace:  DefaultNamespace,
				Repository: "nginx",
				Tag:        "alpine",
			},
		},
		{
			name:  "namespace and repository",
			input: "myorg/myapp",
			want: &OCIReference{
				Registry:   DefaultRegistry,
				Namespace:  "myorg",
				Repository: "myapp",
				Tag:        DefaultTag,
			},
		},
		{
			name:  "full reference with registry",
			input: "registry.example.com/myorg/myapp:v1.0.0",
			want: &OCIReference{
				Registry:   "registry.example.com",
				Namespace:  "myorg",
				Repository: "myapp",
				Tag:        "v1.0.0",
			},
		},
		{
			name:  "registry with port",
			input: "localhost:5000/test/app:latest",
			want: &OCIReference{
				Registry:   "localhost:5000",
				Namespace:  "test",
				Repository: "app",
				Tag:        "latest",
			},
		},
		{
			name:  "with digest",
			input: "nginx@sha256:abc123def456",
			want: &OCIReference{
				Registry:   DefaultRegistry,
				Namespace:  DefaultNamespace,
				Repository: "nginx",
				Digest:     "sha256:abc123def456",
			},
		},
		{
			name:  "full reference with digest",
			input: "quay.io/prometheus/prometheus@sha256:xyz789",
			want: &OCIReference{
				Registry:   "quay.io",
				Namespace:  "prometheus",
				Repository: "prometheus",
				Digest:     "sha256:xyz789",
			},
		},
		{
			name:    "empty reference",
			input:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseOCIReference(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseOCIReference() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				return
			}

			if got.Registry != tt.want.Registry {
				t.Errorf("Registry: got %q, want %q", got.Registry, tt.want.Registry)
			}
			if got.Namespace != tt.want.Namespace {
				t.Errorf("Namespace: got %q, want %q", got.Namespace, tt.want.Namespace)
			}
			if got.Repository != tt.want.Repository {
				t.Errorf("Repository: got %q, want %q", got.Repository, tt.want.Repository)
			}
			if got.Tag != tt.want.Tag {
				t.Errorf("Tag: got %q, want %q", got.Tag, tt.want.Tag)
			}
			if got.Digest != tt.want.Digest {
				t.Errorf("Digest: got %q, want %q", got.Digest, tt.want.Digest)
			}
		})
	}
}

func TestNewOCIImage(t *testing.T) {
	tests := []struct {
		name    string
		ref     string
		wantErr bool
	}{
		{
			name: "valid reference",
			ref:  "nginx:alpine",
		},
		{
			name: "full reference",
			ref:  "registry.example.com/myorg/myapp:v1.0.0",
		},
		{
			name:    "empty reference",
			ref:     "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewOCIImage(tt.ref)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewOCIImage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				return
			}

			if got.Reference == nil {
				t.Errorf("NewOCIImage() reference is nil")
				return
			}
			if got.MediaType != MediaTypeManifest {
				t.Errorf("NewOCIImage() mediaType = %q, want %q", got.MediaType, MediaTypeManifest)
			}
			if got.Annotations == nil {
				t.Errorf("NewOCIImage() annotations is nil")
			}
		})
	}
}

func TestOCIManifest_Validate(t *testing.T) {
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
			wantErr:  false,
		},
		{
			name: "invalid schema version",
			manifest: &OCIManifest{
				SchemaVersion: 1,
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
			},
			wantErr: true,
			errMsg:  "invalid schema version",
		},
		{
			name: "missing mediaType",
			manifest: &OCIManifest{
				SchemaVersion: 2,
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
			},
			wantErr: true,
			errMsg:  "mediaType is required",
		},
		{
			name: "missing config digest",
			manifest: &OCIManifest{
				SchemaVersion: 2,
				MediaType:     MediaTypeManifest,
				Config: OCIDescriptor{
					MediaType: MediaTypeImageConfig,
					Size:      1234,
				},
				Layers: []OCIDescriptor{
					{
						MediaType: MediaTypeImageLayer,
						Digest:    "sha256:layer123",
						Size:      5678,
					},
				},
			},
			wantErr: true,
			errMsg:  "config digest is required",
		},
		{
			name: "no layers",
			manifest: &OCIManifest{
				SchemaVersion: 2,
				MediaType:     MediaTypeManifest,
				Config: OCIDescriptor{
					MediaType: MediaTypeImageConfig,
					Digest:    "sha256:config123",
					Size:      1234,
				},
				Layers: []OCIDescriptor{},
			},
			wantErr: true,
			errMsg:  "manifest must have at least one layer",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.manifest.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("OCIManifest.Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && tt.errMsg != "" {
				if !contains(err.Error(), tt.errMsg) {
					t.Errorf("OCIManifest.Validate() error = %q, want to contain %q", err.Error(), tt.errMsg)
				}
			}
		})
	}
}

// Helper function to check if string contains substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 || 
		(len(s) > len(substr) && findSubstring(s, substr)))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		match := true
		for j := 0; j < len(substr); j++ {
			if s[i+j] != substr[j] {
				match = false
				break
			}
		}
		if match {
			return true
		}
	}
	return false
}