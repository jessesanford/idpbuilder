package k8s

import (
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// ResourceHelper provides utility functions for working with Kubernetes resources
type ResourceHelper struct {
	scheme *runtime.Scheme
}

// NewResourceHelper creates a new ResourceHelper with a default scheme
func NewResourceHelper() *ResourceHelper {
	return &ResourceHelper{
		scheme: NewSchemeBuilder().Build(),
	}
}

// GetGVK returns the GroupVersionKind for a runtime object
func (r *ResourceHelper) GetGVK(obj runtime.Object) schema.GroupVersionKind {
	if obj == nil {
		return schema.GroupVersionKind{}
	}

	// Get the type info from the scheme
	gvks, _, err := r.scheme.ObjectKinds(obj)
	if err != nil || len(gvks) == 0 {
		return schema.GroupVersionKind{}
	}

	return gvks[0]
}

// IsKnownType checks if the given object type is known to the scheme
func (r *ResourceHelper) IsKnownType(obj runtime.Object) bool {
	if obj == nil {
		return false
	}

	gvks, _, err := r.scheme.ObjectKinds(obj)
	return err == nil && len(gvks) > 0
}

// GetAPIVersion returns the API version string for an object
func (r *ResourceHelper) GetAPIVersion(obj runtime.Object) string {
	gvk := r.GetGVK(obj)
	if gvk.Group == "" {
		return gvk.Version
	}
	return fmt.Sprintf("%s/%s", gvk.Group, gvk.Version)
}

// GetKind returns the Kind string for an object
func (r *ResourceHelper) GetKind(obj runtime.Object) string {
	return r.GetGVK(obj).Kind
}

// SetScheme updates the scheme used by this helper
func (r *ResourceHelper) SetScheme(scheme *runtime.Scheme) {
	r.scheme = scheme
}