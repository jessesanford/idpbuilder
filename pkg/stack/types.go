package stack

import (
	"fmt"

	"github.com/cnoe-io/idpbuilder/pkg/oci"
)

// StackConfiguration represents a multi-component application stack configuration
type StackConfiguration struct {
	// Name is the unique identifier for this stack
	Name string `json:"name" yaml:"name"`

	// Version is the stack version (e.g., "1.0.0", "v2.1.3")
	Version string `json:"version" yaml:"version"`

	// Description provides human-readable information about the stack
	Description string `json:"description,omitempty" yaml:"description,omitempty"`

	// Components is the list of components that make up this stack
	Components []StackComponent `json:"components" yaml:"components"`

	// Dependencies defines the dependencies between components in this stack
	Dependencies []StackDependency `json:"dependencies,omitempty" yaml:"dependencies,omitempty"`

	// Status tracks the current lifecycle state of the stack
	Status StackStatus `json:"status" yaml:"status"`
}

// StackComponent represents a single component within a stack
type StackComponent struct {
	// Name is the unique identifier for this component within the stack
	Name string `json:"name" yaml:"name"`

	// Type specifies the component type (e.g., "application", "database", "service")
	Type ComponentType `json:"type" yaml:"type"`

	// Version is the component version
	Version string `json:"version" yaml:"version"`

	// OCIReference points to the OCI image or artifact for this component
	OCIReference *oci.OCIReference `json:"ociReference,omitempty" yaml:"ociReference,omitempty"`

	// Configuration contains component-specific configuration as key-value pairs
	Configuration map[string]interface{} `json:"configuration,omitempty" yaml:"configuration,omitempty"`

	// Status tracks the current state of this component
	Status ComponentStatus `json:"status" yaml:"status"`
}

// StackDependency defines a dependency relationship between components
type StackDependency struct {
	// Component is the name of the component that has the dependency
	Component string `json:"component" yaml:"component"`

	// DependsOn is the name of the component this component depends on
	DependsOn string `json:"dependsOn" yaml:"dependsOn"`

	// Type specifies the dependency type (e.g., "runtime", "build", "optional")
	Type DependencyType `json:"type" yaml:"type"`

	// VersionConstraint specifies version constraints for the dependency (optional)
	VersionConstraint string `json:"versionConstraint,omitempty" yaml:"versionConstraint,omitempty"`
}

// StackStatus represents the overall status of a stack
type StackStatus string

// ComponentStatus represents the status of an individual component
type ComponentStatus string

// ComponentType represents the type of component
type ComponentType string

// DependencyType represents the type of dependency relationship
type DependencyType string

// IsValid validates the StackConfiguration
func (s *StackConfiguration) IsValid() error {
	if s.Name == "" {
		return fmt.Errorf("stack name cannot be empty")
	}
	if s.Version == "" {
		return fmt.Errorf("stack version cannot be empty")
	}
	if len(s.Components) == 0 {
		return fmt.Errorf("stack must have at least one component")
	}

	// Validate each component
	componentNames := make(map[string]bool)
	for i, component := range s.Components {
		if err := component.IsValid(); err != nil {
			return fmt.Errorf("component %d is invalid: %w", i, err)
		}
		if componentNames[component.Name] {
			return fmt.Errorf("duplicate component name: %s", component.Name)
		}
		componentNames[component.Name] = true
	}

	// Validate dependencies reference existing components
	for i, dep := range s.Dependencies {
		if !componentNames[dep.Component] {
			return fmt.Errorf("dependency %d references unknown component: %s", i, dep.Component)
		}
		if !componentNames[dep.DependsOn] {
			return fmt.Errorf("dependency %d references unknown dependency: %s", i, dep.DependsOn)
		}
		if dep.Component == dep.DependsOn {
			return fmt.Errorf("dependency %d: component cannot depend on itself", i)
		}
	}

	return nil
}

// IsValid validates the StackComponent
func (c *StackComponent) IsValid() error {
	if c.Name == "" {
		return fmt.Errorf("component name cannot be empty")
	}
	if c.Type == "" {
		return fmt.Errorf("component type cannot be empty")
	}
	if c.Version == "" {
		return fmt.Errorf("component version cannot be empty")
	}

	// Validate OCI reference if present
	if c.OCIReference != nil {
		if err := validateOCIReference(c.OCIReference); err != nil {
			return fmt.Errorf("invalid OCI reference: %w", err)
		}
	}

	return nil
}

// GetComponentByName retrieves a component by name from the stack
func (s *StackConfiguration) GetComponentByName(name string) (*StackComponent, error) {
	for i := range s.Components {
		if s.Components[i].Name == name {
			return &s.Components[i], nil
		}
	}
	return nil, fmt.Errorf("component not found: %s", name)
}

// GetDependenciesFor returns all dependencies for a given component
func (s *StackConfiguration) GetDependenciesFor(componentName string) []StackDependency {
	var deps []StackDependency
	for _, dep := range s.Dependencies {
		if dep.Component == componentName {
			deps = append(deps, dep)
		}
	}
	return deps
}

// validateOCIReference performs basic validation on an OCI reference
func validateOCIReference(ref *oci.OCIReference) error {
	if ref.Registry == "" {
		return fmt.Errorf("registry cannot be empty")
	}
	if ref.Repository == "" {
		return fmt.Errorf("repository cannot be empty")
	}
	if ref.Tag == "" && ref.Digest == "" {
		return fmt.Errorf("either tag or digest must be specified")
	}
	return nil
}