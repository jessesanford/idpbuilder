package stack

import (
	"fmt"
	"time"
	
	"github.com/cnoe-io/idpbuilder/pkg/oci"
)

// StackConfiguration represents a complete stack configuration
type StackConfiguration struct {
	// Name is the unique name of the stack
	Name string `json:"name" yaml:"name"`
	
	// Namespace is the Kubernetes namespace for the stack
	Namespace string `json:"namespace" yaml:"namespace"`
	
	// Description provides details about the stack's purpose
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
	
	// Components are the individual components that make up the stack
	Components []StackComponent `json:"components" yaml:"components"`
	
	// Dependencies define inter-component dependencies within the stack
	Dependencies []StackDependency `json:"dependencies,omitempty" yaml:"dependencies,omitempty"`
	
	// Timeout specifies the maximum time to wait for stack operations
	Timeout time.Duration `json:"timeout,omitempty" yaml:"timeout,omitempty"`
	
	// RetryAttempts specifies the number of retry attempts for failed operations
	RetryAttempts int `json:"retryAttempts,omitempty" yaml:"retryAttempts,omitempty"`
	
	// Labels provide additional metadata for the stack
	Labels map[string]string `json:"labels,omitempty" yaml:"labels,omitempty"`
	
	// Annotations provide additional metadata for the stack
	Annotations map[string]string `json:"annotations,omitempty" yaml:"annotations,omitempty"`
}

// StackComponent represents an individual component within a stack
type StackComponent struct {
	// Name is the unique name of the component within the stack
	Name string `json:"name" yaml:"name"`
	
	// Type specifies the component type (application, infrastructure, etc.)
	Type string `json:"type" yaml:"type"`
	
	// State tracks the current lifecycle state of the component
	State string `json:"state,omitempty" yaml:"state,omitempty"`
	
	// Image contains OCI image information for the component
	Image *oci.OCIImage `json:"image,omitempty" yaml:"image,omitempty"`
	
	// Config contains component-specific configuration
	Config map[string]interface{} `json:"config,omitempty" yaml:"config,omitempty"`
	
	// Resources specifies resource requirements and limits
	Resources *ComponentResources `json:"resources,omitempty" yaml:"resources,omitempty"`
	
	// HealthCheck defines health check configuration
	HealthCheck *HealthCheckConfig `json:"healthCheck,omitempty" yaml:"healthCheck,omitempty"`
	
	// Labels provide component-specific metadata
	Labels map[string]string `json:"labels,omitempty" yaml:"labels,omitempty"`
	
	// Annotations provide component-specific metadata  
	Annotations map[string]string `json:"annotations,omitempty" yaml:"annotations,omitempty"`
}

// StackDependency defines a dependency relationship between stack components
type StackDependency struct {
	// Component is the name of the component that has the dependency
	Component string `json:"component" yaml:"component"`
	
	// DependsOn is the name of the component that must be ready first
	DependsOn string `json:"dependsOn" yaml:"dependsOn"`
	
	// Type specifies the dependency type (hard, soft, optional)
	Type string `json:"type,omitempty" yaml:"type,omitempty"`
	
	// Condition specifies when the dependency is considered satisfied
	Condition string `json:"condition,omitempty" yaml:"condition,omitempty"`
}

// ComponentResources defines resource requirements for a stack component
type ComponentResources struct {
	// Requests specifies minimum required resources
	Requests map[string]string `json:"requests,omitempty" yaml:"requests,omitempty"`
	
	// Limits specifies maximum allowed resources
	Limits map[string]string `json:"limits,omitempty" yaml:"limits,omitempty"`
}

// HealthCheckConfig defines health check configuration for a component
type HealthCheckConfig struct {
	// Path is the HTTP path to check for HTTP health checks
	Path string `json:"path,omitempty" yaml:"path,omitempty"`
	
	// Port is the port to check
	Port int `json:"port,omitempty" yaml:"port,omitempty"`
	
	// InitialDelaySeconds is the delay before first health check
	InitialDelaySeconds int `json:"initialDelaySeconds,omitempty" yaml:"initialDelaySeconds,omitempty"`
	
	// PeriodSeconds is the interval between health checks
	PeriodSeconds int `json:"periodSeconds,omitempty" yaml:"periodSeconds,omitempty"`
	
	// TimeoutSeconds is the timeout for each health check
	TimeoutSeconds int `json:"timeoutSeconds,omitempty" yaml:"timeoutSeconds,omitempty"`
	
	// FailureThreshold is the number of failures before marking unhealthy
	FailureThreshold int `json:"failureThreshold,omitempty" yaml:"failureThreshold,omitempty"`
}

// Validate validates the stack configuration
func (sc *StackConfiguration) Validate() error {
	if sc.Name == "" {
		return fmt.Errorf("stack name is required")
	}
	
	if len(sc.Name) > MaxStackNameLength {
		return fmt.Errorf("stack name exceeds maximum length of %d", MaxStackNameLength)
	}
	
	if len(sc.Components) == 0 {
		return fmt.Errorf("stack must have at least one component")
	}
	
	if len(sc.Components) > MaxComponentsPerStack {
		return fmt.Errorf("stack exceeds maximum component count of %d", MaxComponentsPerStack)
	}
	
	// Validate component names are unique
	componentNames := make(map[string]bool)
	for _, comp := range sc.Components {
		if comp.Name == "" {
			return fmt.Errorf("component name is required")
		}
		if componentNames[comp.Name] {
			return fmt.Errorf("duplicate component name: %s", comp.Name)
		}
		componentNames[comp.Name] = true
		
		if err := comp.Validate(); err != nil {
			return fmt.Errorf("component %s: %w", comp.Name, err)
		}
	}
	
	// Validate dependencies reference valid components
	if err := sc.validateDependencies(componentNames); err != nil {
		return err
	}
	
	return nil
}

// validateDependencies validates that all dependencies reference valid components
func (sc *StackConfiguration) validateDependencies(componentNames map[string]bool) error {
	// Check for circular dependencies and validate component references
	dependencyGraph := make(map[string][]string)
	
	for _, dep := range sc.Dependencies {
		if !componentNames[dep.Component] {
			return fmt.Errorf("dependency references unknown component: %s", dep.Component)
		}
		if !componentNames[dep.DependsOn] {
			return fmt.Errorf("dependency references unknown component: %s", dep.DependsOn)
		}
		if dep.Component == dep.DependsOn {
			return fmt.Errorf("component cannot depend on itself: %s", dep.Component)
		}
		
		dependencyGraph[dep.Component] = append(dependencyGraph[dep.Component], dep.DependsOn)
	}
	
	// Check for circular dependencies using DFS
	visited := make(map[string]bool)
	recursionStack := make(map[string]bool)
	
	for component := range componentNames {
		if !visited[component] {
			if sc.hasCyclicDependency(component, dependencyGraph, visited, recursionStack) {
				return fmt.Errorf("circular dependency detected involving component: %s", component)
			}
		}
	}
	
	return nil
}

// hasCyclicDependency performs DFS to detect circular dependencies
func (sc *StackConfiguration) hasCyclicDependency(component string, graph map[string][]string, visited, recursionStack map[string]bool) bool {
	visited[component] = true
	recursionStack[component] = true
	
	for _, dependency := range graph[component] {
		if !visited[dependency] {
			if sc.hasCyclicDependency(dependency, graph, visited, recursionStack) {
				return true
			}
		} else if recursionStack[dependency] {
			return true
		}
	}
	
	recursionStack[component] = false
	return false
}

// Validate validates the stack component configuration
func (sc *StackComponent) Validate() error {
	if sc.Name == "" {
		return fmt.Errorf("component name is required")
	}
	
	if sc.Type == "" {
		return fmt.Errorf("component type is required")
	}
	
	// Validate component type
	validTypes := map[string]bool{
		ComponentTypeApplication:     true,
		ComponentTypeInfrastructure: true,
		ComponentTypeService:        true,
		ComponentTypeConfig:          true,
		ComponentTypeSecret:          true,
	}
	
	if !validTypes[sc.Type] {
		return fmt.Errorf("invalid component type: %s", sc.Type)
	}
	
	// Validate component state if set
	if sc.State != "" {
		validStates := map[string]bool{
			ComponentStatePending:      true,
			ComponentStateProcessing:   true,
			ComponentStateReady:        true,
			ComponentStateFailed:       true,
			ComponentStateTerminating: true,
		}
		
		if !validStates[sc.State] {
			return fmt.Errorf("invalid component state: %s", sc.State)
		}
	}
	
	return nil
}

// GetComponentByName returns a component by name
func (sc *StackConfiguration) GetComponentByName(name string) (*StackComponent, error) {
	for i, comp := range sc.Components {
		if comp.Name == name {
			return &sc.Components[i], nil
		}
	}
	return nil, fmt.Errorf("component not found: %s", name)
}

// GetDependencies returns all dependencies for a given component
func (sc *StackConfiguration) GetDependencies(componentName string) []StackDependency {
	var deps []StackDependency
	for _, dep := range sc.Dependencies {
		if dep.Component == componentName {
			deps = append(deps, dep)
		}
	}
	return deps
}

// GetDependents returns all components that depend on the given component
func (sc *StackConfiguration) GetDependents(componentName string) []StackDependency {
	var deps []StackDependency
	for _, dep := range sc.Dependencies {
		if dep.DependsOn == componentName {
			deps = append(deps, dep)
		}
	}
	return deps
}

// SetDefaults sets default values for the stack configuration
func (sc *StackConfiguration) SetDefaults() {
	if sc.Namespace == "" {
		sc.Namespace = DefaultStackNamespace
	}
	
	if sc.Timeout == 0 {
		sc.Timeout = time.Duration(DefaultStackTimeout) * time.Second
	}
	
	if sc.RetryAttempts == 0 {
		sc.RetryAttempts = DefaultRetryAttempts
	}
	
	if sc.Labels == nil {
		sc.Labels = make(map[string]string)
	}
	
	if sc.Annotations == nil {
		sc.Annotations = make(map[string]string)
	}
	
	// Set default component states
	for i := range sc.Components {
		if sc.Components[i].State == "" {
			sc.Components[i].State = ComponentStatePending
		}
		if sc.Components[i].Labels == nil {
			sc.Components[i].Labels = make(map[string]string)
		}
		if sc.Components[i].Annotations == nil {
			sc.Components[i].Annotations = make(map[string]string)
		}
	}
}