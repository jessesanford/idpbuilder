package k8s

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"
)

// SchemeBuilder provides a builder interface for constructing runtime schemes
type SchemeBuilder struct {
	funcs []func(*runtime.Scheme) error
}

// NewSchemeBuilder creates a new SchemeBuilder with default Kubernetes schemes
func NewSchemeBuilder() *SchemeBuilder {
	return &SchemeBuilder{
		funcs: []func(*runtime.Scheme) error{},
	}
}

// Build creates a runtime.Scheme from the registered functions
func (s *SchemeBuilder) Build() *runtime.Scheme {
	sch := runtime.NewScheme()

	// Add core Kubernetes schemes
	if err := scheme.AddToScheme(sch); err != nil {
		return runtime.NewScheme()
	}

	// Apply registered functions
	for _, f := range s.funcs {
		if err := f(sch); err != nil {
			continue
		}
	}

	return sch
}