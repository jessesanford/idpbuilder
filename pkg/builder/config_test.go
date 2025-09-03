package builder

import (
	"reflect"
	"testing"

	v1 "github.com/google/go-containerregistry/pkg/v1"
)

func TestNewConfigFactory(t *testing.T) {
	factory := NewConfigFactory()

	if factory == nil {
		t.Fatal("NewConfigFactory returned nil")
	}

	if factory.BaseConfig == nil {
		t.Error("BaseConfig should not be nil")
	}

	if factory.DefaultLabels == nil {
		t.Error("DefaultLabels should not be nil")
	}

	if len(factory.DefaultLabels) == 0 {
		t.Error("DefaultLabels should contain default labels")
	}
}

func TestConfigFactory_CreateConfig(t *testing.T) {
	factory := NewConfigFactory()

	tests := []struct {
		name    string
		opts    *BuildOptions
		wantErr bool
	}{
		{
			name:    "nil options",
			opts:    nil,
			wantErr: true,
		},
		{
			name:    "valid options",
			opts:    NewBuildOptions(),
			wantErr: false,
		},
		{
			name: "invalid options",
			opts: &BuildOptions{
				Platform: nil, // This will cause validation to fail
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config, err := factory.CreateConfig(tt.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConfigFactory.CreateConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && config == nil {
				t.Error("ConfigFactory.CreateConfig() returned nil config for valid input")
			}
		})
	}
}

func TestConfigFactory_CreateConfig_Values(t *testing.T) {
	factory := NewConfigFactory()
	opts := NewBuildOptions()

	// Set some custom values
	opts.WorkingDir = "/app"
	opts.User = "appuser"
	opts.Entrypoint = []string{"/app/entrypoint.sh"}
	opts.Cmd = []string{"start"}
	opts.AddLabel("app.name", "test-app")
	opts.AddEnvironment("APP_ENV", "test")
	opts.ExposedPorts = []string{"8080/tcp", "9090/tcp"}

	config, err := factory.CreateConfig(opts)
	if err != nil {
		t.Fatalf("ConfigFactory.CreateConfig() error = %v", err)
	}

	if config.WorkingDir != "/app" {
		t.Errorf("Expected WorkingDir '/app', got '%s'", config.WorkingDir)
	}

	if config.User != "appuser" {
		t.Errorf("Expected User 'appuser', got '%s'", config.User)
	}

	if !reflect.DeepEqual(config.Entrypoint, []string{"/app/entrypoint.sh"}) {
		t.Errorf("Expected Entrypoint ['/app/entrypoint.sh'], got %v", config.Entrypoint)
	}

	if !reflect.DeepEqual(config.Cmd, []string{"start"}) {
		t.Errorf("Expected Cmd ['start'], got %v", config.Cmd)
	}

	if config.Labels["app.name"] != "test-app" {
		t.Errorf("Expected label 'app.name' with value 'test-app', got '%s'", config.Labels["app.name"])
	}

	// Check that environment was merged properly
	foundAppEnv := false
	for _, env := range config.Env {
		if env == "APP_ENV=test" {
			foundAppEnv = true
			break
		}
	}
	if !foundAppEnv {
		t.Error("Expected environment variable 'APP_ENV=test' not found")
	}

	// Check exposed ports
	if len(config.ExposedPorts) != 2 {
		t.Errorf("Expected 2 exposed ports, got %d", len(config.ExposedPorts))
	}

	if _, exists := config.ExposedPorts["8080/tcp"]; !exists {
		t.Error("Expected port '8080/tcp' to be exposed")
	}
}

func TestConfigFactory_CreatePlatformConfig(t *testing.T) {
	factory := NewConfigFactory()

	tests := []struct {
		name     string
		opts     *BuildOptions
		expected *v1.Platform
		wantErr  bool
	}{
		{
			name:     "nil options",
			opts:     nil,
			expected: &v1.Platform{Architecture: "amd64", OS: "linux"},
			wantErr:  false,
		},
		{
			name: "valid platform",
			opts: &BuildOptions{
				Platform: &v1.Platform{
					Architecture: "arm64",
					OS:           "linux",
				},
			},
			expected: &v1.Platform{Architecture: "arm64", OS: "linux"},
			wantErr:  false,
		},
		{
			name: "invalid platform",
			opts: &BuildOptions{
				Platform: &v1.Platform{
					Architecture: "invalid",
					OS:           "linux",
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			platform, err := factory.CreatePlatformConfig(tt.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConfigFactory.CreatePlatformConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if platform.Architecture != tt.expected.Architecture {
					t.Errorf("Expected architecture '%s', got '%s'", tt.expected.Architecture, platform.Architecture)
				}
				if platform.OS != tt.expected.OS {
					t.Errorf("Expected OS '%s', got '%s'", tt.expected.OS, platform.OS)
				}
			}
		})
	}
}

func TestValidatePlatform(t *testing.T) {
	tests := []struct {
		name     string
		platform *v1.Platform
		wantErr  bool
	}{
		{
			name: "valid linux/amd64",
			platform: &v1.Platform{
				OS:           "linux",
				Architecture: "amd64",
			},
			wantErr: false,
		},
		{
			name: "valid linux/arm64",
			platform: &v1.Platform{
				OS:           "linux",
				Architecture: "arm64",
			},
			wantErr: false,
		},
		{
			name: "empty OS",
			platform: &v1.Platform{
				OS:           "",
				Architecture: "amd64",
			},
			wantErr: true,
		},
		{
			name: "empty architecture",
			platform: &v1.Platform{
				OS:           "linux",
				Architecture: "",
			},
			wantErr: true,
		},
		{
			name: "unsupported OS",
			platform: &v1.Platform{
				OS:           "unsupported",
				Architecture: "amd64",
			},
			wantErr: true,
		},
		{
			name: "unsupported architecture",
			platform: &v1.Platform{
				OS:           "linux",
				Architecture: "unsupported",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validatePlatform(tt.platform)
			if (err != nil) != tt.wantErr {
				t.Errorf("validatePlatform() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestParsePlatform(t *testing.T) {
	tests := []struct {
		name        string
		platformStr string
		expected    *v1.Platform
		wantErr     bool
	}{
		{
			name:        "empty string",
			platformStr: "",
			expected:    &v1.Platform{OS: "linux", Architecture: "amd64"},
			wantErr:     false,
		},
		{
			name:        "linux/amd64",
			platformStr: "linux/amd64",
			expected:    &v1.Platform{OS: "linux", Architecture: "amd64"},
			wantErr:     false,
		},
		{
			name:        "linux/arm64/v8",
			platformStr: "linux/arm64/v8",
			expected:    &v1.Platform{OS: "linux", Architecture: "arm64", Variant: "v8"},
			wantErr:     false,
		},
		{
			name:        "invalid format",
			platformStr: "invalid",
			wantErr:     true,
		},
		{
			name:        "too many parts",
			platformStr: "linux/amd64/v1/extra",
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			platform, err := ParsePlatform(tt.platformStr)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParsePlatform() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if platform.OS != tt.expected.OS {
					t.Errorf("Expected OS '%s', got '%s'", tt.expected.OS, platform.OS)
				}
				if platform.Architecture != tt.expected.Architecture {
					t.Errorf("Expected architecture '%s', got '%s'", tt.expected.Architecture, platform.Architecture)
				}
				if platform.Variant != tt.expected.Variant {
					t.Errorf("Expected variant '%s', got '%s'", tt.expected.Variant, platform.Variant)
				}
			}
		})
	}
}

func TestFormatPlatform(t *testing.T) {
	tests := []struct {
		name     string
		platform *v1.Platform
		expected string
	}{
		{
			name:     "nil platform",
			platform: nil,
			expected: "linux/amd64",
		},
		{
			name: "linux/amd64",
			platform: &v1.Platform{
				OS:           "linux",
				Architecture: "amd64",
			},
			expected: "linux/amd64",
		},
		{
			name: "linux/arm64/v8",
			platform: &v1.Platform{
				OS:           "linux",
				Architecture: "arm64",
				Variant:      "v8",
			},
			expected: "linux/arm64/v8",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatPlatform(tt.platform)
			if result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

func TestGenerateConfigDigest(t *testing.T) {
	config1 := &v1.Config{
		User:       "root",
		WorkingDir: "/app",
		Entrypoint: []string{"/entrypoint.sh"},
		Cmd:        []string{"start"},
	}

	config2 := &v1.Config{
		User:       "root",
		WorkingDir: "/app",
		Entrypoint: []string{"/entrypoint.sh"},
		Cmd:        []string{"start"},
	}

	config3 := &v1.Config{
		User:       "appuser", // Different user
		WorkingDir: "/app",
		Entrypoint: []string{"/entrypoint.sh"},
		Cmd:        []string{"start"},
	}

	digest1, err := GenerateConfigDigest(config1)
	if err != nil {
		t.Fatalf("GenerateConfigDigest() error = %v", err)
	}

	digest2, err := GenerateConfigDigest(config2)
	if err != nil {
		t.Fatalf("GenerateConfigDigest() error = %v", err)
	}

	digest3, err := GenerateConfigDigest(config3)
	if err != nil {
		t.Fatalf("GenerateConfigDigest() error = %v", err)
	}

	// Same configs should produce same digest
	if digest1 != digest2 {
		t.Error("Same configs should produce same digest")
	}

	// Different configs should produce different digests
	if digest1 == digest3 {
		t.Error("Different configs should produce different digests")
	}

	// Test nil config
	_, err = GenerateConfigDigest(nil)
	if err == nil {
		t.Error("GenerateConfigDigest() should error for nil config")
	}
}