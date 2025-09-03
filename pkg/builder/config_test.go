package builder

import (
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestDefaultBuildConfig(t *testing.T) {
	config := DefaultBuildConfig()

	// Test default values
	if config.ContextPath != "." {
		t.Errorf("expected ContextPath to be '.', got %s", config.ContextPath)
	}

	if config.Dockerfile != "Dockerfile" {
		t.Errorf("expected Dockerfile to be 'Dockerfile', got %s", config.Dockerfile)
	}

	if len(config.Tags) != 1 || config.Tags[0] != "latest" {
		t.Errorf("expected Tags to be ['latest'], got %v", config.Tags)
	}

	if config.Platform.OS != "linux" {
		t.Errorf("expected Platform.OS to be 'linux', got %s", config.Platform.OS)
	}

	if config.Platform.Architecture != "amd64" {
		t.Errorf("expected Platform.Architecture to be 'amd64', got %s", config.Platform.Architecture)
	}

	if config.Registry.Hostname != "index.docker.io" {
		t.Errorf("expected Registry.Hostname to be 'index.docker.io', got %s", config.Registry.Hostname)
	}

	if config.BuildArgs == nil {
		t.Error("expected BuildArgs to be initialized")
	}

	if config.Labels == nil {
		t.Error("expected Labels to be initialized")
	}

	if !config.Remove {
		t.Error("expected Remove to be true by default")
	}

	if config.BuildTimeout != 30*time.Minute {
		t.Errorf("expected BuildTimeout to be 30 minutes, got %v", config.BuildTimeout)
	}
}

func TestBuildConfigValidate(t *testing.T) {
	tests := []struct {
		name      string
		config    *BuildConfig
		expectErr bool
	}{
		{
			name:      "valid default config",
			config:    DefaultBuildConfig(),
			expectErr: false,
		},
		{
			name: "empty context path",
			config: &BuildConfig{
				ContextPath: "",
				Dockerfile:  "Dockerfile",
				Tags:        []string{"test:latest"},
				Platform:    DefaultPlatformConfig(),
				Registry:    DefaultRegistryConfig(),
				BuildArgs:   make(map[string]string),
				Labels:      make(map[string]string),
				BuildTimeout: 30 * time.Minute,
			},
			expectErr: true,
		},
		{
			name: "empty dockerfile",
			config: &BuildConfig{
				ContextPath: ".",
				Dockerfile:  "",
				Tags:        []string{"test:latest"},
				Platform:    DefaultPlatformConfig(),
				Registry:    DefaultRegistryConfig(),
				BuildArgs:   make(map[string]string),
				Labels:      make(map[string]string),
				BuildTimeout: 30 * time.Minute,
			},
			expectErr: true,
		},
		{
			name: "no tags",
			config: &BuildConfig{
				ContextPath: ".",
				Dockerfile:  "Dockerfile",
				Tags:        []string{},
				Platform:    DefaultPlatformConfig(),
				Registry:    DefaultRegistryConfig(),
				BuildArgs:   make(map[string]string),
				Labels:      make(map[string]string),
				BuildTimeout: 30 * time.Minute,
			},
			expectErr: true,
		},
		{
			name: "empty tag",
			config: &BuildConfig{
				ContextPath: ".",
				Dockerfile:  "Dockerfile",
				Tags:        []string{"valid:tag", ""},
				Platform:    DefaultPlatformConfig(),
				Registry:    DefaultRegistryConfig(),
				BuildArgs:   make(map[string]string),
				Labels:      make(map[string]string),
				BuildTimeout: 30 * time.Minute,
			},
			expectErr: true,
		},
		{
			name: "tag with spaces",
			config: &BuildConfig{
				ContextPath: ".",
				Dockerfile:  "Dockerfile",
				Tags:        []string{"invalid tag"},
				Platform:    DefaultPlatformConfig(),
				Registry:    DefaultRegistryConfig(),
				BuildArgs:   make(map[string]string),
				Labels:      make(map[string]string),
				BuildTimeout: 30 * time.Minute,
			},
			expectErr: true,
		},
		{
			name: "zero build timeout",
			config: &BuildConfig{
				ContextPath: ".",
				Dockerfile:  "Dockerfile",
				Tags:        []string{"test:latest"},
				Platform:    DefaultPlatformConfig(),
				Registry:    DefaultRegistryConfig(),
				BuildArgs:   make(map[string]string),
				Labels:      make(map[string]string),
				BuildTimeout: 0,
			},
			expectErr: true,
		},
		{
			name: "negative memory limit",
			config: &BuildConfig{
				ContextPath:  ".",
				Dockerfile:   "Dockerfile",
				Tags:         []string{"test:latest"},
				Platform:     DefaultPlatformConfig(),
				Registry:     DefaultRegistryConfig(),
				BuildArgs:    make(map[string]string),
				Labels:       make(map[string]string),
				BuildTimeout: 30 * time.Minute,
				MemoryLimit:  -1,
			},
			expectErr: true,
		},
		{
			name: "negative CPU limit",
			config: &BuildConfig{
				ContextPath:  ".",
				Dockerfile:   "Dockerfile",
				Tags:         []string{"test:latest"},
				Platform:     DefaultPlatformConfig(),
				Registry:     DefaultRegistryConfig(),
				BuildArgs:    make(map[string]string),
				Labels:       make(map[string]string),
				BuildTimeout: 30 * time.Minute,
				CPULimit:     -1.0,
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if tt.expectErr && err == nil {
				t.Error("expected error but got none")
			}
			if !tt.expectErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestPlatformConfigValidate(t *testing.T) {
	tests := []struct {
		name      string
		platform  PlatformConfig
		expectErr bool
	}{
		{
			name:      "valid default platform",
			platform:  DefaultPlatformConfig(),
			expectErr: false,
		},
		{
			name: "valid linux/arm64",
			platform: PlatformConfig{
				OS:           "linux",
				Architecture: "arm64",
			},
			expectErr: false,
		},
		{
			name: "valid windows/amd64",
			platform: PlatformConfig{
				OS:           "windows",
				Architecture: "amd64",
			},
			expectErr: false,
		},
		{
			name: "empty OS",
			platform: PlatformConfig{
				OS:           "",
				Architecture: "amd64",
			},
			expectErr: true,
		},
		{
			name: "empty architecture",
			platform: PlatformConfig{
				OS:           "linux",
				Architecture: "",
			},
			expectErr: true,
		},
		{
			name: "invalid OS",
			platform: PlatformConfig{
				OS:           "invalid",
				Architecture: "amd64",
			},
			expectErr: true,
		},
		{
			name: "invalid architecture",
			platform: PlatformConfig{
				OS:           "linux",
				Architecture: "invalid",
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.platform.Validate()
			if tt.expectErr && err == nil {
				t.Error("expected error but got none")
			}
			if !tt.expectErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestRegistryConfigValidate(t *testing.T) {
	tests := []struct {
		name      string
		registry  RegistryConfig
		expectErr bool
	}{
		{
			name:      "valid default registry",
			registry:  DefaultRegistryConfig(),
			expectErr: false,
		},
		{
			name: "valid with username/password",
			registry: RegistryConfig{
				Hostname: "registry.example.com",
				Username: "user",
				Password: "pass",
				Timeout:  30 * time.Second,
			},
			expectErr: false,
		},
		{
			name: "valid with token",
			registry: RegistryConfig{
				Hostname: "registry.example.com",
				Token:    "token123",
				Timeout:  30 * time.Second,
			},
			expectErr: false,
		},
		{
			name: "empty hostname",
			registry: RegistryConfig{
				Hostname: "",
				Timeout:  30 * time.Second,
			},
			expectErr: true,
		},
		{
			name: "username without password",
			registry: RegistryConfig{
				Hostname: "registry.example.com",
				Username: "user",
				Timeout:  30 * time.Second,
			},
			expectErr: true,
		},
		{
			name: "password without username",
			registry: RegistryConfig{
				Hostname: "registry.example.com",
				Password: "pass",
				Timeout:  30 * time.Second,
			},
			expectErr: true,
		},
		{
			name: "both username/password and token",
			registry: RegistryConfig{
				Hostname: "registry.example.com",
				Username: "user",
				Password: "pass",
				Token:    "token123",
				Timeout:  30 * time.Second,
			},
			expectErr: true,
		},
		{
			name: "zero timeout",
			registry: RegistryConfig{
				Hostname: "registry.example.com",
				Timeout:  0,
			},
			expectErr: true,
		},
		{
			name: "negative retry count",
			registry: RegistryConfig{
				Hostname:   "registry.example.com",
				Timeout:    30 * time.Second,
				RetryCount: -1,
			},
			expectErr: true,
		},
		{
			name: "negative retry delay",
			registry: RegistryConfig{
				Hostname:   "registry.example.com",
				Timeout:    30 * time.Second,
				RetryDelay: -1 * time.Second,
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.registry.Validate()
			if tt.expectErr && err == nil {
				t.Error("expected error but got none")
			}
			if !tt.expectErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestPlatformConfigString(t *testing.T) {
	tests := []struct {
		name     string
		platform PlatformConfig
		expected string
	}{
		{
			name: "linux/amd64",
			platform: PlatformConfig{
				OS:           "linux",
				Architecture: "amd64",
			},
			expected: "linux/amd64",
		},
		{
			name: "linux/arm64/v8",
			platform: PlatformConfig{
				OS:           "linux",
				Architecture: "arm64",
				Variant:      "v8",
			},
			expected: "linux/arm64/v8",
		},
		{
			name: "windows/amd64",
			platform: PlatformConfig{
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
				t.Errorf("expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestBuildConfigContextDir(t *testing.T) {
	config := &BuildConfig{ContextPath: "."}
	
	contextDir, err := config.ContextDir()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should return an absolute path
	if !filepath.IsAbs(contextDir) {
		t.Errorf("expected absolute path, got %s", contextDir)
	}
}

func TestBuildConfigDockerfilePath(t *testing.T) {
	config := &BuildConfig{
		ContextPath: ".",
		Dockerfile:  "Dockerfile.test",
	}
	
	dockerfilePath, err := config.DockerfilePath()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should return an absolute path
	if !filepath.IsAbs(dockerfilePath) {
		t.Errorf("expected absolute path, got %s", dockerfilePath)
	}

	// Should end with the dockerfile name
	if !strings.HasSuffix(dockerfilePath, "Dockerfile.test") {
		t.Errorf("expected path to end with 'Dockerfile.test', got %s", dockerfilePath)
	}
}

func TestBuildConfigClone(t *testing.T) {
	original := DefaultBuildConfig()
	original.Tags = []string{"original:tag"}
	original.BuildArgs["TEST"] = "value"
	original.Labels["test"] = "label"

	clone := original.Clone()

	// Modify original to ensure independence
	original.Tags[0] = "modified:tag"
	original.BuildArgs["TEST"] = "modified"
	original.Labels["test"] = "modified"

	// Clone should be unaffected
	if clone.Tags[0] != "original:tag" {
		t.Errorf("clone was affected by original modification")
	}

	if clone.BuildArgs["TEST"] != "value" {
		t.Errorf("clone build args were affected by original modification")
	}

	if clone.Labels["test"] != "label" {
		t.Errorf("clone labels were affected by original modification")
	}
}