package mapper

import (
	"context"
	"testing"
)

func TestMapStackToContainers(t *testing.T) {
	mapper := NewStackMapper()

	tests := []struct {
		name        string
		config      StackConfig
		expectError bool
		errorType   ErrorCode
	}{
		{
			name: "valid stack config",
			config: StackConfig{
				Name:    "test-stack",
				Version: "1.0.0",
				Components: []Component{
					{
						Name:   "web-server",
						Type:   "web",
						Source: "nginx:alpine",
						Config: map[string]any{
							"env": map[string]any{
								"PORT": "8080",
							},
							"labels": map[string]any{
								"service": "frontend",
							},
						},
					},
					{
						Name:   "api-server",
						Type:   "api",
						Source: "alpine:latest",
					},
				},
				Metadata: map[string]string{
					"environment": "test",
				},
			},
			expectError: false,
		},
		{
			name: "empty stack name",
			config: StackConfig{
				Name:       "",
				Version:    "1.0.0",
				Components: []Component{{Name: "test", Type: "web", Source: "nginx"}},
			},
			expectError: true,
			errorType:   ErrInvalidConfig,
		},
		{
			name: "empty stack version",
			config: StackConfig{
				Name:       "test-stack",
				Version:    "",
				Components: []Component{{Name: "test", Type: "web", Source: "nginx"}},
			},
			expectError: true,
			errorType:   ErrInvalidConfig,
		},
		{
			name: "no components",
			config: StackConfig{
				Name:       "test-stack",
				Version:    "1.0.0",
				Components: []Component{},
			},
			expectError: true,
			errorType:   ErrInvalidConfig,
		},
		{
			name: "component missing name",
			config: StackConfig{
				Name:    "test-stack",
				Version: "1.0.0",
				Components: []Component{
					{Name: "", Type: "web", Source: "nginx"},
				},
			},
			expectError: true,
			errorType:   ErrInvalidConfig,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.WithValue(context.Background(), "registry", "registry.example.com")
			result, err := mapper.MapStackToContainers(ctx, tt.config)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
					return
				}
				if mappingErr, ok := err.(*MappingError); ok {
					if mappingErr.Code != tt.errorType {
						t.Errorf("expected error type %v, got %v", tt.errorType, mappingErr.Code)
					}
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			// Validate the result
			if len(result.Containers) != len(tt.config.Components) {
				t.Errorf("expected %d containers, got %d",
					len(tt.config.Components), len(result.Containers))
			}

			if result.Manifest == nil {
				t.Error("expected manifest to be generated")
			}

			// Check metadata
			if result.Metadata["stack_name"] != tt.config.Name {
				t.Errorf("expected stack_name %s, got %s",
					tt.config.Name, result.Metadata["stack_name"])
			}

			// Validate container specs
			for i, container := range result.Containers {
				component := tt.config.Components[i]
				if container.Name != component.Name {
					t.Errorf("container %d: expected name %s, got %s",
						i, component.Name, container.Name)
				}
				if container.BaseImage == "" {
					t.Errorf("container %d: base image should not be empty", i)
				}
				if container.Labels["component.name"] != component.Name {
					t.Errorf("container %d: missing component.name label", i)
				}
			}
		})
	}
}

func TestResolveReferences(t *testing.T) {
	mapper := NewStackMapper()

	tests := []struct {
		name        string
		refs        []string
		expectError bool
		expected    map[string]ComponentRef
	}{
		{
			name: "empty references",
			refs: []string{},
			expected: map[string]ComponentRef{},
		},
		{
			name: "simple image reference",
			refs: []string{"nginx:alpine"},
			expected: map[string]ComponentRef{
				"nginx:alpine": {
					Registry:   "docker.io",
					Repository: "library/nginx",
					Tag:        "alpine",
				},
			},
		},
		{
			name: "full registry reference",
			refs: []string{"registry.example.com/myorg/myapp:v1.0.0"},
			expected: map[string]ComponentRef{
				"registry.example.com/myorg/myapp:v1.0.0": {
					Registry:   "registry.example.com",
					Repository: "myorg/myapp",
					Tag:        "v1.0.0",
				},
			},
		},
		{
			name: "digest reference",
			refs: []string{"alpine@sha256:abcd1234567890abcdef1234567890abcdef1234567890abcdef1234567890ab"},
			expected: map[string]ComponentRef{
				"alpine@sha256:abcd1234567890abcdef1234567890abcdef1234567890abcdef1234567890ab": {
					Registry:   "docker.io",
					Repository: "library/alpine",
					Digest:     "sha256:abcd1234567890abcdef1234567890abcdef1234567890abcdef1234567890ab",
				},
			},
		},
		{
			name:        "invalid digest format",
			refs:        []string{"alpine@invalid-digest"},
			expectError: true,
		},
		{
			name:        "invalid tag characters",
			refs:        []string{"nginx:tag-with-spaces spaces"},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			result, err := mapper.ResolveReferences(ctx, tt.refs)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if len(result) != len(tt.expected) {
				t.Errorf("expected %d resolved references, got %d",
					len(tt.expected), len(result))
				return
			}

			for ref, expectedRef := range tt.expected {
				actualRef, exists := result[ref]
				if !exists {
					t.Errorf("missing resolved reference for %s", ref)
					continue
				}

				if actualRef.Registry != expectedRef.Registry {
					t.Errorf("ref %s: expected registry %s, got %s",
						ref, expectedRef.Registry, actualRef.Registry)
				}
				if actualRef.Repository != expectedRef.Repository {
					t.Errorf("ref %s: expected repository %s, got %s",
						ref, expectedRef.Repository, actualRef.Repository)
				}
				if actualRef.Tag != expectedRef.Tag {
					t.Errorf("ref %s: expected tag %s, got %s",
						ref, expectedRef.Tag, actualRef.Tag)
				}
				if actualRef.Digest != expectedRef.Digest {
					t.Errorf("ref %s: expected digest %s, got %s",
						ref, expectedRef.Digest, actualRef.Digest)
				}
			}
		})
	}
}

func TestValidateMapping(t *testing.T) {
	mapper := NewStackMapper()

	tests := []struct {
		name        string
		mapping     MappingResult
		expectError bool
		errorType   ErrorCode
	}{
		{
			name: "valid mapping",
			mapping: MappingResult{
				Containers: []ContainerSpec{
					{
						Name:      "test-container",
						BaseImage: "alpine:latest",
						Env:       map[string]string{"KEY": "value"},
						Labels:    map[string]string{"app": "test"},
					},
				},
				Metadata: map[string]string{"test": "data"},
			},
			expectError: false,
		},
		{
			name: "empty containers",
			mapping: MappingResult{
				Containers: []ContainerSpec{},
			},
			expectError: true,
			errorType:   ErrValidationFailed,
		},
		{
			name: "container missing name",
			mapping: MappingResult{
				Containers: []ContainerSpec{
					{
						Name:      "",
						BaseImage: "alpine:latest",
					},
				},
			},
			expectError: true,
			errorType:   ErrValidationFailed,
		},
		{
			name: "container missing base image",
			mapping: MappingResult{
				Containers: []ContainerSpec{
					{
						Name:      "test",
						BaseImage: "",
					},
				},
			},
			expectError: true,
			errorType:   ErrValidationFailed,
		},
		{
			name: "invalid environment variable key",
			mapping: MappingResult{
				Containers: []ContainerSpec{
					{
						Name:      "test",
						BaseImage: "alpine:latest",
						Env:       map[string]string{"KEY=BAD": "value"},
					},
				},
			},
			expectError: true,
			errorType:   ErrValidationFailed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := mapper.ValidateMapping(tt.mapping)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
					return
				}
				if mappingErr, ok := err.(*MappingError); ok {
					if mappingErr.Code != tt.errorType {
						t.Errorf("expected error type %v, got %v", tt.errorType, mappingErr.Code)
					}
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestMockStackMapper(t *testing.T) {
	t.Run("default behavior", func(t *testing.T) {
		mock := NewMockStackMapper()

		config := StackConfig{
			Name:    "test",
			Version: "1.0.0",
			Components: []Component{
				{Name: "comp1", Type: "web", Source: "nginx"},
			},
		}

		result, err := mock.MapStackToContainers(context.Background(), config)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if len(result.Containers) != 1 {
			t.Errorf("expected 1 container, got %d", len(result.Containers))
		}
	})

	t.Run("error injection", func(t *testing.T) {
		mock := NewMockStackMapper().WithMappingError(newMappingError(ErrInvalidConfig, "test error"))

		config := StackConfig{Name: "test", Version: "1.0.0"}
		_, err := mock.MapStackToContainers(context.Background(), config)
		if err == nil {
			t.Error("expected error but got none")
		}
	})
}