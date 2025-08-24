package oci

import (
	"testing"
	"time"
)

func TestOCIPlatform_String(t *testing.T) {
	tests := []struct {
		name     string
		platform OCIPlatform
		expected string
	}{
		{
			name:     "basic linux/amd64",
			platform: OCIPlatform{OS: "linux", Architecture: "amd64"},
			expected: "linux/amd64",
		},
		{
			name:     "with variant",
			platform: OCIPlatform{OS: "linux", Architecture: "arm", Variant: "v7"},
			expected: "linux/arm/v7",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.platform.String()
			if result != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestOCIReference_String(t *testing.T) {
	tests := []struct {
		name     string
		ref      OCIReference
		expected string
	}{
		{
			name:     "with registry and tag",
			ref:      OCIReference{Registry: "registry.io", Repository: "myapp/image", Tag: "v1.0.0"},
			expected: "registry.io/myapp/image:v1.0.0",
		},
		{
			name:     "default registry",
			ref:      OCIReference{Registry: DefaultRegistry, Repository: "nginx", Tag: "latest"},
			expected: "nginx:latest",
		},
		{
			name:     "with digest",
			ref:      OCIReference{Repository: "myapp/image", Digest: "sha256:abc123"},
			expected: "myapp/image@sha256:abc123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.ref.String()
			if result != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestOCIReference_IsValid(t *testing.T) {
	tests := []struct {
		name  string
		ref   OCIReference
		valid bool
	}{
		{
			name:  "valid with tag",
			ref:   OCIReference{Repository: "myapp/image", Tag: "v1.0.0"},
			valid: true,
		},
		{
			name:  "valid with digest",
			ref:   OCIReference{Repository: "myapp/image", Digest: "sha256:abc123"},
			valid: true,
		},
		{
			name:  "invalid - no repository",
			ref:   OCIReference{Tag: "v1.0.0"},
			valid: false,
		},
		{
			name:  "invalid - no tag or digest",
			ref:   OCIReference{Repository: "myapp/image"},
			valid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.ref.IsValid()
			if result != tt.valid {
				t.Errorf("expected %t, got %t", tt.valid, result)
			}
		})
	}
}

func TestOCIImage_Creation(t *testing.T) {
	now := time.Now()
	image := OCIImage{
		Name:       "nginx",
		Tag:        "latest",
		Registry:   "docker.io",
		Repository: "library/nginx",
		Platform: OCIPlatform{
			OS:           OSLinux,
			Architecture: ArchAMD64,
		},
		CreatedAt: now,
		Size:      1024,
	}

	if image.Name != "nginx" {
		t.Errorf("expected name nginx, got %s", image.Name)
	}
	if image.Platform.OS != OSLinux {
		t.Errorf("expected OS linux, got %s", image.Platform.OS)
	}
	if image.Size != 1024 {
		t.Errorf("expected size 1024, got %d", image.Size)
	}
}