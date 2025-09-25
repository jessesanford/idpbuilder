package v1alpha1

import (
	"errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// CustomPackageSpec defines the desired state of CustomPackage
type CustomPackageSpec struct {
	// Source specifies the source location of the custom package
	Source CustomPackageSource `json:"source"`

	// Dependencies specifies the list of dependencies required by this package
	Dependencies []PackageDependency `json:"dependencies,omitempty"`

	// Configuration holds package-specific configuration
	Configuration map[string]string `json:"configuration,omitempty"`
}

// CustomPackageSource defines where the custom package comes from
type CustomPackageSource struct {
	// Type specifies the source type (git, oci, local)
	Type string `json:"type"`

	// URL specifies the source URL
	URL string `json:"url"`

	// Path specifies the path within the source
	Path string `json:"path,omitempty"`

	// Ref specifies the reference (branch, tag, commit) for git sources
	Ref string `json:"ref,omitempty"`
}

// PackageDependency defines a dependency of the custom package
type PackageDependency struct {
	// Name is the name of the dependency
	Name string `json:"name"`

	// Version is the version constraint
	Version string `json:"version,omitempty"`

	// Required indicates if this dependency is required
	Required bool `json:"required,omitempty"`
}

// CustomPackageStatus defines the observed state of CustomPackage
type CustomPackageStatus struct {
	// Phase represents the current phase of the custom package
	Phase string `json:"phase,omitempty"`

	// Conditions represent the latest available observations
	Conditions []metav1.Condition `json:"conditions,omitempty"`

	// Message provides additional information about the current state
	Message string `json:"message,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// CustomPackage is the Schema for the custompackages API
type CustomPackage struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CustomPackageSpec   `json:"spec,omitempty"`
	Status CustomPackageStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// CustomPackageList contains a list of CustomPackage
type CustomPackageList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []CustomPackage `json:"items"`
}

// Validate validates the CustomPackage specification
func (c *CustomPackage) Validate() error {
	if c.Spec.Source.Type == "" {
		return errors.New("source type is required")
	}

	if c.Spec.Source.URL == "" {
		return errors.New("source URL is required")
	}

	// Validate source type
	switch c.Spec.Source.Type {
	case "git", "oci", "local":
		// Valid types
	default:
		return errors.New("source type must be one of: git, oci, local")
	}

	return nil
}

// GetName returns the name of the CustomPackage
func (c *CustomPackage) GetName() string {
	return c.ObjectMeta.Name
}

// GetNamespace returns the namespace of the CustomPackage
func (c *CustomPackage) GetNamespace() string {
	return c.ObjectMeta.Namespace
}

// DeepCopyObject returns a deep copy of the CustomPackage
func (c *CustomPackage) DeepCopyObject() runtime.Object {
	if copy := c.DeepCopy(); copy != nil {
		return copy
	}
	return nil
}

// DeepCopy creates a deep copy of the CustomPackage
func (c *CustomPackage) DeepCopy() *CustomPackage {
	if c == nil {
		return nil
	}
	out := new(CustomPackage)
	c.DeepCopyInto(out)
	return out
}

// DeepCopyInto copies all properties from this object into another object
func (c *CustomPackage) DeepCopyInto(out *CustomPackage) {
	*out = *c
	out.TypeMeta = c.TypeMeta
	c.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = c.Spec
	if c.Spec.Configuration != nil {
		out.Spec.Configuration = make(map[string]string, len(c.Spec.Configuration))
		for key, val := range c.Spec.Configuration {
			out.Spec.Configuration[key] = val
		}
	}
	if c.Spec.Dependencies != nil {
		out.Spec.Dependencies = make([]PackageDependency, len(c.Spec.Dependencies))
		copy(out.Spec.Dependencies, c.Spec.Dependencies)
	}
	out.Status = c.Status
	if c.Status.Conditions != nil {
		out.Status.Conditions = make([]metav1.Condition, len(c.Status.Conditions))
		for i := range c.Status.Conditions {
			c.Status.Conditions[i].DeepCopyInto(&out.Status.Conditions[i])
		}
	}
}

// DeepCopyObject returns a deep copy of the CustomPackageList
func (c *CustomPackageList) DeepCopyObject() runtime.Object {
	if copy := c.DeepCopy(); copy != nil {
		return copy
	}
	return nil
}

// DeepCopy creates a deep copy of the CustomPackageList
func (c *CustomPackageList) DeepCopy() *CustomPackageList {
	if c == nil {
		return nil
	}
	out := new(CustomPackageList)
	c.DeepCopyInto(out)
	return out
}

// DeepCopyInto copies all properties from this object into another object
func (c *CustomPackageList) DeepCopyInto(out *CustomPackageList) {
	*out = *c
	out.TypeMeta = c.TypeMeta
	c.ListMeta.DeepCopyInto(&out.ListMeta)
	if c.Items != nil {
		out.Items = make([]CustomPackage, len(c.Items))
		for i := range c.Items {
			c.Items[i].DeepCopyInto(&out.Items[i])
		}
	}
}

func init() {
	SchemeBuilder.Register(&CustomPackage{}, &CustomPackageList{})
}