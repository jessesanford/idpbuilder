package stack

import (
	"testing"
	"time"
	
	"github.com/cnoe-io/idpbuilder/pkg/oci"
)

func TestStackConfiguration_Validate(t *testing.T) {
	tests := []struct {
		name    string
		stack   *StackConfiguration
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid stack configuration",
			stack: &StackConfiguration{
				Name:      "test-stack",
				Namespace: "test-namespace",
				Components: []StackComponent{
					{
						Name: "component1",
						Type: ComponentTypeApplication,
					},
					{
						Name: "component2", 
						Type: ComponentTypeService,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "empty stack name",
			stack: &StackConfiguration{
				Components: []StackComponent{
					{Name: "comp1", Type: ComponentTypeApplication},
				},
			},
			wantErr: true,
			errMsg:  "stack name is required",
		},
		{
			name: "stack name too long",
			stack: &StackConfiguration{
				Name: "this-stack-name-is-way-too-long-and-exceeds-the-maximum-allowed-length-for-stack-names",
				Components: []StackComponent{
					{Name: "comp1", Type: ComponentTypeApplication},
				},
			},
			wantErr: true,
			errMsg:  "stack name exceeds maximum length",
		},
		{
			name: "no components",
			stack: &StackConfiguration{
				Name:       "test-stack",
				Components: []StackComponent{},
			},
			wantErr: true,
			errMsg:  "stack must have at least one component",
		},
		{
			name: "duplicate component names",
			stack: &StackConfiguration{
				Name: "test-stack",
				Components: []StackComponent{
					{Name: "duplicate", Type: ComponentTypeApplication},
					{Name: "duplicate", Type: ComponentTypeService},
				},
			},
			wantErr: true,
			errMsg:  "duplicate component name",
		},
		{
			name: "invalid dependency component",
			stack: &StackConfiguration{
				Name: "test-stack",
				Components: []StackComponent{
					{Name: "comp1", Type: ComponentTypeApplication},
				},
				Dependencies: []StackDependency{
					{Component: "nonexistent", DependsOn: "comp1"},
				},
			},
			wantErr: true,
			errMsg:  "dependency references unknown component",
		},
		{
			name: "circular dependency",
			stack: &StackConfiguration{
				Name: "test-stack",
				Components: []StackComponent{
					{Name: "comp1", Type: ComponentTypeApplication},
					{Name: "comp2", Type: ComponentTypeService},
				},
				Dependencies: []StackDependency{
					{Component: "comp1", DependsOn: "comp2"},
					{Component: "comp2", DependsOn: "comp1"},
				},
			},
			wantErr: true,
			errMsg:  "circular dependency detected",
		},
		{
			name: "self dependency",
			stack: &StackConfiguration{
				Name: "test-stack",
				Components: []StackComponent{
					{Name: "comp1", Type: ComponentTypeApplication},
				},
				Dependencies: []StackDependency{
					{Component: "comp1", DependsOn: "comp1"},
				},
			},
			wantErr: true,
			errMsg:  "component cannot depend on itself",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.stack.Validate()
			
			if (err != nil) != tt.wantErr {
				t.Errorf("StackConfiguration.Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			
			if err != nil && tt.errMsg != "" {
				if !contains(err.Error(), tt.errMsg) {
					t.Errorf("StackConfiguration.Validate() error = %q, want to contain %q", err.Error(), tt.errMsg)
				}
			}
		})
	}
}

func TestStackComponent_Validate(t *testing.T) {
	tests := []struct {
		name      string
		component *StackComponent
		wantErr   bool
		errMsg    string
	}{
		{
			name: "valid component",
			component: &StackComponent{
				Name: "test-component",
				Type: ComponentTypeApplication,
			},
			wantErr: false,
		},
		{
			name: "valid component with image",
			component: &StackComponent{
				Name: "test-component",
				Type: ComponentTypeApplication,
				Image: &oci.OCIImage{
					Reference: &oci.OCIReference{
						Registry:   "docker.io",
						Namespace:  "library",
						Repository: "nginx",
						Tag:        "alpine",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "empty component name",
			component: &StackComponent{
				Type: ComponentTypeApplication,
			},
			wantErr: true,
			errMsg:  "component name is required",
		},
		{
			name: "empty component type",
			component: &StackComponent{
				Name: "test-component",
			},
			wantErr: true,
			errMsg:  "component type is required",
		},
		{
			name: "invalid component type",
			component: &StackComponent{
				Name: "test-component",
				Type: "invalid-type",
			},
			wantErr: true,
			errMsg:  "invalid component type",
		},
		{
			name: "invalid component state",
			component: &StackComponent{
				Name:  "test-component",
				Type:  ComponentTypeApplication,
				State: "invalid-state",
			},
			wantErr: true,
			errMsg:  "invalid component state",
		},
		{
			name: "valid component with state",
			component: &StackComponent{
				Name:  "test-component",
				Type:  ComponentTypeApplication,
				State: ComponentStateReady,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.component.Validate()
			
			if (err != nil) != tt.wantErr {
				t.Errorf("StackComponent.Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			
			if err != nil && tt.errMsg != "" {
				if !contains(err.Error(), tt.errMsg) {
					t.Errorf("StackComponent.Validate() error = %q, want to contain %q", err.Error(), tt.errMsg)
				}
			}
		})
	}
}

func TestStackConfiguration_GetComponentByName(t *testing.T) {
	stack := &StackConfiguration{
		Name: "test-stack",
		Components: []StackComponent{
			{Name: "comp1", Type: ComponentTypeApplication},
			{Name: "comp2", Type: ComponentTypeService},
		},
	}

	tests := []struct {
		name          string
		componentName string
		wantFound     bool
		wantName      string
	}{
		{
			name:          "existing component",
			componentName: "comp1",
			wantFound:     true,
			wantName:      "comp1",
		},
		{
			name:          "non-existing component",
			componentName: "comp3",
			wantFound:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := stack.GetComponentByName(tt.componentName)
			
			if tt.wantFound {
				if err != nil {
					t.Errorf("GetComponentByName() error = %v, want nil", err)
					return
				}
				if got.Name != tt.wantName {
					t.Errorf("GetComponentByName() name = %v, want %v", got.Name, tt.wantName)
				}
			} else {
				if err == nil {
					t.Errorf("GetComponentByName() error = nil, want error")
				}
				if got != nil {
					t.Errorf("GetComponentByName() = %v, want nil", got)
				}
			}
		})
	}
}

func TestStackConfiguration_GetDependencies(t *testing.T) {
	stack := &StackConfiguration{
		Name: "test-stack",
		Components: []StackComponent{
			{Name: "comp1", Type: ComponentTypeApplication},
			{Name: "comp2", Type: ComponentTypeService},
			{Name: "comp3", Type: ComponentTypeConfig},
		},
		Dependencies: []StackDependency{
			{Component: "comp1", DependsOn: "comp2"},
			{Component: "comp1", DependsOn: "comp3"},
			{Component: "comp2", DependsOn: "comp3"},
		},
	}

	tests := []struct {
		name          string
		componentName string
		wantCount     int
		wantDeps      []string
	}{
		{
			name:          "component with multiple dependencies",
			componentName: "comp1",
			wantCount:     2,
			wantDeps:      []string{"comp2", "comp3"},
		},
		{
			name:          "component with single dependency",
			componentName: "comp2",
			wantCount:     1,
			wantDeps:      []string{"comp3"},
		},
		{
			name:          "component with no dependencies",
			componentName: "comp3",
			wantCount:     0,
			wantDeps:      []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := stack.GetDependencies(tt.componentName)
			
			if len(got) != tt.wantCount {
				t.Errorf("GetDependencies() count = %v, want %v", len(got), tt.wantCount)
				return
			}
			
			for i, dep := range got {
				if i < len(tt.wantDeps) && dep.DependsOn != tt.wantDeps[i] {
					t.Errorf("GetDependencies()[%d].DependsOn = %v, want %v", i, dep.DependsOn, tt.wantDeps[i])
				}
			}
		})
	}
}

func TestStackConfiguration_GetDependents(t *testing.T) {
	stack := &StackConfiguration{
		Name: "test-stack",
		Components: []StackComponent{
			{Name: "comp1", Type: ComponentTypeApplication},
			{Name: "comp2", Type: ComponentTypeService},
			{Name: "comp3", Type: ComponentTypeConfig},
		},
		Dependencies: []StackDependency{
			{Component: "comp1", DependsOn: "comp2"},
			{Component: "comp1", DependsOn: "comp3"},
			{Component: "comp2", DependsOn: "comp3"},
		},
	}

	tests := []struct {
		name          string
		componentName string
		wantCount     int
		wantDeps      []string
	}{
		{
			name:          "component with multiple dependents",
			componentName: "comp3",
			wantCount:     2,
			wantDeps:      []string{"comp1", "comp2"},
		},
		{
			name:          "component with single dependent",
			componentName: "comp2",
			wantCount:     1,
			wantDeps:      []string{"comp1"},
		},
		{
			name:          "component with no dependents",
			componentName: "comp1",
			wantCount:     0,
			wantDeps:      []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := stack.GetDependents(tt.componentName)
			
			if len(got) != tt.wantCount {
				t.Errorf("GetDependents() count = %v, want %v", len(got), tt.wantCount)
				return
			}
			
			for i, dep := range got {
				if i < len(tt.wantDeps) && dep.Component != tt.wantDeps[i] {
					t.Errorf("GetDependents()[%d].Component = %v, want %v", i, dep.Component, tt.wantDeps[i])
				}
			}
		})
	}
}

func TestStackConfiguration_SetDefaults(t *testing.T) {
	stack := &StackConfiguration{
		Name: "test-stack",
		Components: []StackComponent{
			{Name: "comp1", Type: ComponentTypeApplication},
			{Name: "comp2", Type: ComponentTypeService, State: ComponentStateReady},
		},
	}

	stack.SetDefaults()

	// Test default values
	if stack.Namespace != DefaultStackNamespace {
		t.Errorf("SetDefaults() namespace = %v, want %v", stack.Namespace, DefaultStackNamespace)
	}
	
	expectedTimeout := time.Duration(DefaultStackTimeout) * time.Second
	if stack.Timeout != expectedTimeout {
		t.Errorf("SetDefaults() timeout = %v, want %v", stack.Timeout, expectedTimeout)
	}
	
	if stack.RetryAttempts != DefaultRetryAttempts {
		t.Errorf("SetDefaults() retryAttempts = %v, want %v", stack.RetryAttempts, DefaultRetryAttempts)
	}
	
	if stack.Labels == nil {
		t.Errorf("SetDefaults() labels should be initialized")
	}
	
	if stack.Annotations == nil {
		t.Errorf("SetDefaults() annotations should be initialized")
	}
	
	// Test component defaults
	if stack.Components[0].State != ComponentStatePending {
		t.Errorf("SetDefaults() component[0] state = %v, want %v", stack.Components[0].State, ComponentStatePending)
	}
	
	if stack.Components[1].State != ComponentStateReady {
		t.Errorf("SetDefaults() component[1] state should not change from %v", ComponentStateReady)
	}
	
	for i, comp := range stack.Components {
		if comp.Labels == nil {
			t.Errorf("SetDefaults() component[%d] labels should be initialized", i)
		}
		if comp.Annotations == nil {
			t.Errorf("SetDefaults() component[%d] annotations should be initialized", i)
		}
	}
}

func TestStackConfiguration_SetDefaults_PreserveExisting(t *testing.T) {
	customTimeout := 120 * time.Second
	stack := &StackConfiguration{
		Name:          "test-stack",
		Namespace:     "custom-namespace",
		Timeout:       customTimeout,
		RetryAttempts: 5,
		Labels:        map[string]string{"existing": "label"},
		Annotations:   map[string]string{"existing": "annotation"},
		Components: []StackComponent{
			{Name: "comp1", Type: ComponentTypeApplication},
		},
	}

	stack.SetDefaults()

	// Test that existing values are preserved
	if stack.Namespace != "custom-namespace" {
		t.Errorf("SetDefaults() should preserve existing namespace")
	}
	
	if stack.Timeout != customTimeout {
		t.Errorf("SetDefaults() should preserve existing timeout")
	}
	
	if stack.RetryAttempts != 5 {
		t.Errorf("SetDefaults() should preserve existing retryAttempts")
	}
	
	if stack.Labels["existing"] != "label" {
		t.Errorf("SetDefaults() should preserve existing labels")
	}
	
	if stack.Annotations["existing"] != "annotation" {
		t.Errorf("SetDefaults() should preserve existing annotations")
	}
}

func TestComponentTypes(t *testing.T) {
	validTypes := []string{
		ComponentTypeApplication,
		ComponentTypeInfrastructure,
		ComponentTypeService,
		ComponentTypeConfig,
		ComponentTypeSecret,
	}

	for _, componentType := range validTypes {
		t.Run("component type "+componentType, func(t *testing.T) {
			component := &StackComponent{
				Name: "test-component",
				Type: componentType,
			}
			
			if err := component.Validate(); err != nil {
				t.Errorf("Valid component type %s should not produce error: %v", componentType, err)
			}
		})
	}
}

func TestComponentStates(t *testing.T) {
	validStates := []string{
		ComponentStatePending,
		ComponentStateProcessing,
		ComponentStateReady,
		ComponentStateFailed,
		ComponentStateTerminating,
	}

	for _, state := range validStates {
		t.Run("component state "+state, func(t *testing.T) {
			component := &StackComponent{
				Name:  "test-component",
				Type:  ComponentTypeApplication,
				State: state,
			}
			
			if err := component.Validate(); err != nil {
				t.Errorf("Valid component state %s should not produce error: %v", state, err)
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