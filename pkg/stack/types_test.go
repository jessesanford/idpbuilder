package stack

import (
	"testing"

	"github.com/cnoe-io/idpbuilder/pkg/oci"
)

func TestStackConfiguration_IsValid(t *testing.T) {
	tests := []struct {
		name    string
		stack   *StackConfiguration
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid stack configuration",
			stack: &StackConfiguration{
				Name:    "test-stack",
				Version: "1.0.0",
				Status:  StackStatusPending,
				Components: []StackComponent{
					{
						Name:    "web-app",
						Type:    ComponentTypeApplication,
						Version: "1.0.0",
						Status:  ComponentStatusPending,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "empty name",
			stack: &StackConfiguration{
				Name:    "",
				Version: "1.0.0",
				Components: []StackComponent{
					{Name: "test", Type: ComponentTypeApplication, Version: "1.0.0"},
				},
			},
			wantErr: true,
			errMsg:  "stack name cannot be empty",
		},
		{
			name: "empty version",
			stack: &StackConfiguration{
				Name:    "test-stack",
				Version: "",
				Components: []StackComponent{
					{Name: "test", Type: ComponentTypeApplication, Version: "1.0.0"},
				},
			},
			wantErr: true,
			errMsg:  "stack version cannot be empty",
		},
		{
			name: "no components",
			stack: &StackConfiguration{
				Name:       "test-stack",
				Version:    "1.0.0",
				Components: []StackComponent{},
			},
			wantErr: true,
			errMsg:  "stack must have at least one component",
		},
		{
			name: "duplicate component names",
			stack: &StackConfiguration{
				Name:    "test-stack",
				Version: "1.0.0",
				Components: []StackComponent{
					{Name: "duplicate", Type: ComponentTypeApplication, Version: "1.0.0"},
					{Name: "duplicate", Type: ComponentTypeService, Version: "1.0.0"},
				},
			},
			wantErr: true,
			errMsg:  "duplicate component name: duplicate",
		},
		{
			name: "dependency references unknown component",
			stack: &StackConfiguration{
				Name:    "test-stack",
				Version: "1.0.0",
				Components: []StackComponent{
					{Name: "app", Type: ComponentTypeApplication, Version: "1.0.0"},
				},
				Dependencies: []StackDependency{
					{Component: "app", DependsOn: "nonexistent", Type: DependencyTypeRuntime},
				},
			},
			wantErr: true,
			errMsg:  "dependency 0 references unknown dependency: nonexistent",
		},
		{
			name: "circular dependency",
			stack: &StackConfiguration{
				Name:    "test-stack",
				Version: "1.0.0",
				Components: []StackComponent{
					{Name: "app", Type: ComponentTypeApplication, Version: "1.0.0"},
				},
				Dependencies: []StackDependency{
					{Component: "app", DependsOn: "app", Type: DependencyTypeRuntime},
				},
			},
			wantErr: true,
			errMsg:  "dependency 0: component cannot depend on itself",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.stack.IsValid()
			if tt.wantErr {
				if err == nil {
					t.Errorf("StackConfiguration.IsValid() expected error but got none")
				} else if err.Error() != tt.errMsg {
					t.Errorf("StackConfiguration.IsValid() error = %v, want %v", err.Error(), tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("StackConfiguration.IsValid() unexpected error = %v", err)
				}
			}
		})
	}
}

func TestStackComponent_IsValid(t *testing.T) {
	tests := []struct {
		name      string
		component *StackComponent
		wantErr   bool
		errMsg    string
	}{
		{
			name: "valid component",
			component: &StackComponent{
				Name:    "test-component",
				Type:    ComponentTypeApplication,
				Version: "1.0.0",
				Status:  ComponentStatusPending,
			},
			wantErr: false,
		},
		{
			name: "valid component with OCI reference",
			component: &StackComponent{
				Name:    "test-component",
				Type:    ComponentTypeApplication,
				Version: "1.0.0",
				Status:  ComponentStatusPending,
				OCIReference: &oci.OCIReference{
					Registry:   "docker.io",
					Repository: "library/nginx",
					Tag:        "latest",
				},
			},
			wantErr: false,
		},
		{
			name: "empty name",
			component: &StackComponent{
				Name:    "",
				Type:    ComponentTypeApplication,
				Version: "1.0.0",
			},
			wantErr: true,
			errMsg:  "component name cannot be empty",
		},
		{
			name: "empty type",
			component: &StackComponent{
				Name:    "test-component",
				Type:    "",
				Version: "1.0.0",
			},
			wantErr: true,
			errMsg:  "component type cannot be empty",
		},
		{
			name: "empty version",
			component: &StackComponent{
				Name:    "test-component",
				Type:    ComponentTypeApplication,
				Version: "",
			},
			wantErr: true,
			errMsg:  "component version cannot be empty",
		},
		{
			name: "invalid OCI reference - empty registry",
			component: &StackComponent{
				Name:    "test-component",
				Type:    ComponentTypeApplication,
				Version: "1.0.0",
				OCIReference: &oci.OCIReference{
					Registry:   "",
					Repository: "library/nginx",
					Tag:        "latest",
				},
			},
			wantErr: true,
			errMsg:  "invalid OCI reference: registry cannot be empty",
		},
		{
			name: "invalid OCI reference - empty repository",
			component: &StackComponent{
				Name:    "test-component",
				Type:    ComponentTypeApplication,
				Version: "1.0.0",
				OCIReference: &oci.OCIReference{
					Registry:   "docker.io",
					Repository: "",
					Tag:        "latest",
				},
			},
			wantErr: true,
			errMsg:  "invalid OCI reference: repository cannot be empty",
		},
		{
			name: "invalid OCI reference - no tag or digest",
			component: &StackComponent{
				Name:    "test-component",
				Type:    ComponentTypeApplication,
				Version: "1.0.0",
				OCIReference: &oci.OCIReference{
					Registry:   "docker.io",
					Repository: "library/nginx",
				},
			},
			wantErr: true,
			errMsg:  "invalid OCI reference: either tag or digest must be specified",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.component.IsValid()
			if tt.wantErr {
				if err == nil {
					t.Errorf("StackComponent.IsValid() expected error but got none")
				} else if err.Error() != tt.errMsg {
					t.Errorf("StackComponent.IsValid() error = %v, want %v", err.Error(), tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("StackComponent.IsValid() unexpected error = %v", err)
				}
			}
		})
	}
}

func TestStackConfiguration_GetComponentByName(t *testing.T) {
	stack := &StackConfiguration{
		Name:    "test-stack",
		Version: "1.0.0",
		Components: []StackComponent{
			{Name: "web-app", Type: ComponentTypeApplication, Version: "1.0.0"},
			{Name: "database", Type: ComponentTypeDatabase, Version: "2.0.0"},
		},
	}

	tests := []struct {
		name          string
		componentName string
		wantErr       bool
		expectedType  ComponentType
	}{
		{
			name:          "existing component",
			componentName: "web-app",
			wantErr:       false,
			expectedType:  ComponentTypeApplication,
		},
		{
			name:          "another existing component",
			componentName: "database",
			wantErr:       false,
			expectedType:  ComponentTypeDatabase,
		},
		{
			name:          "non-existing component",
			componentName: "nonexistent",
			wantErr:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			component, err := stack.GetComponentByName(tt.componentName)
			if tt.wantErr {
				if err == nil {
					t.Errorf("GetComponentByName() expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("GetComponentByName() unexpected error = %v", err)
				}
				if component == nil {
					t.Errorf("GetComponentByName() expected component but got nil")
				} else if component.Type != tt.expectedType {
					t.Errorf("GetComponentByName() component type = %v, want %v", component.Type, tt.expectedType)
				}
			}
		})
	}
}

func TestStackConfiguration_GetDependenciesFor(t *testing.T) {
	stack := &StackConfiguration{
		Name:    "test-stack",
		Version: "1.0.0",
		Components: []StackComponent{
			{Name: "web-app", Type: ComponentTypeApplication, Version: "1.0.0"},
			{Name: "database", Type: ComponentTypeDatabase, Version: "2.0.0"},
			{Name: "cache", Type: ComponentTypeMiddleware, Version: "1.5.0"},
		},
		Dependencies: []StackDependency{
			{Component: "web-app", DependsOn: "database", Type: DependencyTypeRuntime},
			{Component: "web-app", DependsOn: "cache", Type: DependencyTypeOptional},
			{Component: "cache", DependsOn: "database", Type: DependencyTypeNetwork},
		},
	}

	tests := []struct {
		name              string
		componentName     string
		expectedDepsCount int
		expectedDeps      []string
	}{
		{
			name:              "component with multiple dependencies",
			componentName:     "web-app",
			expectedDepsCount: 2,
			expectedDeps:      []string{"database", "cache"},
		},
		{
			name:              "component with single dependency",
			componentName:     "cache",
			expectedDepsCount: 1,
			expectedDeps:      []string{"database"},
		},
		{
			name:              "component with no dependencies",
			componentName:     "database",
			expectedDepsCount: 0,
			expectedDeps:      []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			deps := stack.GetDependenciesFor(tt.componentName)
			if len(deps) != tt.expectedDepsCount {
				t.Errorf("GetDependenciesFor() dependencies count = %v, want %v", len(deps), tt.expectedDepsCount)
			}

			dependsOnNames := make([]string, len(deps))
			for i, dep := range deps {
				dependsOnNames[i] = dep.DependsOn
			}

			for _, expectedDep := range tt.expectedDeps {
				found := false
				for _, actualDep := range dependsOnNames {
					if actualDep == expectedDep {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("GetDependenciesFor() expected dependency %s not found", expectedDep)
				}
			}
		})
	}
}

func TestValidateOCIReference(t *testing.T) {
	tests := []struct {
		name    string
		ref     *oci.OCIReference
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid reference with tag",
			ref: &oci.OCIReference{
				Registry:   "docker.io",
				Repository: "library/nginx",
				Tag:        "latest",
			},
			wantErr: false,
		},
		{
			name: "valid reference with digest",
			ref: &oci.OCIReference{
				Registry:   "ghcr.io",
				Repository: "myorg/myapp",
				Digest:     "sha256:abc123",
			},
			wantErr: false,
		},
		{
			name: "empty registry",
			ref: &oci.OCIReference{
				Registry:   "",
				Repository: "library/nginx",
				Tag:        "latest",
			},
			wantErr: true,
			errMsg:  "registry cannot be empty",
		},
		{
			name: "empty repository",
			ref: &oci.OCIReference{
				Registry:   "docker.io",
				Repository: "",
				Tag:        "latest",
			},
			wantErr: true,
			errMsg:  "repository cannot be empty",
		},
		{
			name: "no tag or digest",
			ref: &oci.OCIReference{
				Registry:   "docker.io",
				Repository: "library/nginx",
			},
			wantErr: true,
			errMsg:  "either tag or digest must be specified",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateOCIReference(tt.ref)
			if tt.wantErr {
				if err == nil {
					t.Errorf("validateOCIReference() expected error but got none")
				} else if err.Error() != tt.errMsg {
					t.Errorf("validateOCIReference() error = %v, want %v", err.Error(), tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("validateOCIReference() unexpected error = %v", err)
				}
			}
		})
	}
}

func TestStackConfigurationWithConfiguration(t *testing.T) {
	// Test component with configuration map
	component := &StackComponent{
		Name:    "web-app",
		Type:    ComponentTypeApplication,
		Version: "1.0.0",
		Status:  ComponentStatusPending,
		Configuration: map[string]interface{}{
			ConfigKeyReplicas:     3,
			ConfigKeyMemoryLimit: "512Mi",
			ConfigKeyCPULimit:    "500m",
			ConfigKeyEnvironment: map[string]string{
				"ENV": "production",
				"LOG_LEVEL": "info",
			},
		},
	}

	if err := component.IsValid(); err != nil {
		t.Errorf("StackComponent with configuration should be valid, got error: %v", err)
	}

	// Verify configuration values
	if replicas, ok := component.Configuration[ConfigKeyReplicas].(int); !ok || replicas != 3 {
		t.Errorf("Expected replicas=3, got %v", component.Configuration[ConfigKeyReplicas])
	}
}