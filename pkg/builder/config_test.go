package builder

import (
	"testing"

	"github.com/google/go-containerregistry/pkg/v1"
	"github.com/stretchr/testify/assert"
)

// TestConfigFactory tests configuration generation
func TestConfigFactory(t *testing.T) {
	platform := v1.Platform{
		OS:           "linux",
		Architecture: "amd64",
	}
	
	factory := NewConfigFactory(platform)
	assert.NotNil(t, factory)
	
	t.Run("generate basic config", func(t *testing.T) {
		opts := BuildOptions{
			Platform: platform,
			Env:      []string{"ENV_VAR=value"},
			Cmd:      []string{"echo", "hello"},
		}
		
		config, err := factory.GenerateConfig(opts)
		assert.NoError(t, err)
		assert.NotNil(t, config)
		
		assert.Equal(t, "linux", config.OS)
		assert.Equal(t, "amd64", config.Architecture)
		assert.Contains(t, config.Config.Env, "ENV_VAR=value")
		assert.Equal(t, []string{"echo", "hello"}, config.Config.Cmd)
	})
	
	t.Run("config with labels", func(t *testing.T) {
		opts := BuildOptions{
			Platform: platform,
			Labels: map[string]string{
				"org.opencontainers.image.created": "", // Should be set to build time
				"test-label":                          "test-value",
			},
		}
		
		config, err := factory.GenerateConfig(opts)
		assert.NoError(t, err)
		assert.NotEmpty(t, config.Config.Labels["org.opencontainers.image.created"])
		assert.Equal(t, "test-value", config.Config.Labels["test-label"])
	})
	
	t.Run("validation errors", func(t *testing.T) {
		// Test invalid working directory
		opts := BuildOptions{
			Platform:   platform,
			WorkingDir: "relative/path", // Should be absolute
		}
		
		_, err := factory.GenerateConfig(opts)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "working directory must be an absolute path")
	})
}

// TestConfigFactoryValidation tests more validation scenarios
func TestConfigFactoryValidation(t *testing.T) {
	platform := v1.Platform{OS: "linux", Architecture: "amd64"}
	factory := NewConfigFactory(platform)
	
	tests := []struct {
		name    string
		opts    BuildOptions
		wantErr bool
		errMsg  string
	}{
		{
			name: "invalid port format",
			opts: BuildOptions{
				Platform: platform,
				ExposedPorts: map[string]struct{}{
					"invalid": {},
				},
			},
			wantErr: true,
			errMsg:  "invalid exposed port format",
		},
		{
			name: "invalid environment variable",
			opts: BuildOptions{
				Platform: platform,
				Env:      []string{"INVALID_ENV_VAR"},
			},
			wantErr: true,
			errMsg:  "environment variable must be in KEY=value format",
		},
		{
			name: "invalid user format",
			opts: BuildOptions{
				Platform: platform,
				User:     "user:group:extra",
			},
			wantErr: true,
			errMsg:  "user specification can have at most one colon",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := factory.GenerateConfig(tt.opts)
			
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestConfigFactoryHelpers tests helper functions
func TestConfigFactoryHelpers(t *testing.T) {
	t.Run("DefaultLabels", func(t *testing.T) {
		labels := DefaultLabels("test-source")
		assert.Contains(t, labels, "org.opencontainers.image.source")
		assert.Equal(t, "test-source", labels["org.opencontainers.image.source"])
	})
	
	t.Run("MergeConfigs", func(t *testing.T) {
		config1 := &v1.ConfigFile{
			Architecture: "amd64",
			OS:           "linux",
			Config: v1.Config{
				Env: []string{"VAR1=value1"},
				Labels: map[string]string{
					"label1": "value1",
				},
			},
		}
		
		config2 := &v1.ConfigFile{
			Architecture: "arm64", // Override
			Config: v1.Config{
				Env: []string{"VAR2=value2"},
				Labels: map[string]string{
					"label2": "value2",
				},
			},
		}
		
		merged := MergeConfigs(config1, config2)
		assert.Equal(t, "arm64", merged.Architecture) // Should be overridden
		assert.Equal(t, "linux", merged.OS)           // Should be preserved
		assert.Contains(t, merged.Config.Env, "VAR1=value1")
		assert.Contains(t, merged.Config.Env, "VAR2=value2")
		assert.Equal(t, "value1", merged.Config.Labels["label1"])
		assert.Equal(t, "value2", merged.Config.Labels["label2"])
	})
}