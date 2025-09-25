package v1alpha1

import (
	"errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// LocalBuildSpec defines the desired state of LocalBuild
type LocalBuildSpec struct {
	// Source specifies the source location for the build
	Source LocalBuildSource `json:"source"`

	// BuildArgs specifies build arguments
	BuildArgs map[string]string `json:"buildArgs,omitempty"`

	// Target specifies the target stage in a multi-stage build
	Target string `json:"target,omitempty"`

	// Platform specifies the target platform
	Platform string `json:"platform,omitempty"`

	// Context specifies the build context path
	Context string `json:"context,omitempty"`

	// Dockerfile specifies the Dockerfile path
	Dockerfile string `json:"dockerfile,omitempty"`
}

// LocalBuildSource defines the source for local builds
type LocalBuildSource struct {
	// Type specifies the source type (local, git)
	Type string `json:"type"`

	// Path specifies the local path or Git URL
	Path string `json:"path"`

	// Ref specifies the Git reference (for git type)
	Ref string `json:"ref,omitempty"`
}

// LocalBuildStatus defines the observed state of LocalBuild
type LocalBuildStatus struct {
	// Phase represents the current phase of the build
	Phase string `json:"phase,omitempty"`

	// ImageRef contains the built image reference
	ImageRef string `json:"imageRef,omitempty"`

	// StartTime represents the time when the build started
	StartTime *metav1.Time `json:"startTime,omitempty"`

	// CompletionTime represents the time when the build completed
	CompletionTime *metav1.Time `json:"completionTime,omitempty"`

	// Conditions represent the latest available observations
	Conditions []metav1.Condition `json:"conditions,omitempty"`

	// Message provides additional information about the current state
	Message string `json:"message,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// LocalBuild is the Schema for the localbuilds API
type LocalBuild struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   LocalBuildSpec   `json:"spec,omitempty"`
	Status LocalBuildStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// LocalBuildList contains a list of LocalBuild
type LocalBuildList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []LocalBuild `json:"items"`
}

// Validate validates the LocalBuild specification
func (l *LocalBuild) Validate() error {
	if l.Spec.Source.Type == "" {
		return errors.New("source type is required")
	}

	if l.Spec.Source.Path == "" {
		return errors.New("source path is required")
	}

	// Validate source type
	switch l.Spec.Source.Type {
	case "local", "git":
		// Valid types
	default:
		return errors.New("source type must be one of: local, git")
	}

	return nil
}

// IsReady returns true if the LocalBuild is ready
func (l *LocalBuild) IsReady() bool {
	for _, condition := range l.Status.Conditions {
		if condition.Type == "Ready" && condition.Status == metav1.ConditionTrue {
			return true
		}
	}
	return false
}

// DeepCopyObject returns a deep copy of the LocalBuild
func (l *LocalBuild) DeepCopyObject() runtime.Object {
	if c := l.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopy creates a deep copy of the LocalBuild
func (l *LocalBuild) DeepCopy() *LocalBuild {
	if l == nil {
		return nil
	}
	out := new(LocalBuild)
	l.DeepCopyInto(out)
	return out
}

// DeepCopyInto copies all properties from this object into another object
func (l *LocalBuild) DeepCopyInto(out *LocalBuild) {
	*out = *l
	out.TypeMeta = l.TypeMeta
	l.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = l.Spec
	if l.Spec.BuildArgs != nil {
		out.Spec.BuildArgs = make(map[string]string, len(l.Spec.BuildArgs))
		for key, val := range l.Spec.BuildArgs {
			out.Spec.BuildArgs[key] = val
		}
	}
	out.Status = l.Status
	if l.Status.StartTime != nil {
		out.Status.StartTime = l.Status.StartTime.DeepCopy()
	}
	if l.Status.CompletionTime != nil {
		out.Status.CompletionTime = l.Status.CompletionTime.DeepCopy()
	}
	if l.Status.Conditions != nil {
		out.Status.Conditions = make([]metav1.Condition, len(l.Status.Conditions))
		for i := range l.Status.Conditions {
			l.Status.Conditions[i].DeepCopyInto(&out.Status.Conditions[i])
		}
	}
}

// DeepCopyObject returns a deep copy of the LocalBuildList
func (l *LocalBuildList) DeepCopyObject() runtime.Object {
	if c := l.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopy creates a deep copy of the LocalBuildList
func (l *LocalBuildList) DeepCopy() *LocalBuildList {
	if l == nil {
		return nil
	}
	out := new(LocalBuildList)
	l.DeepCopyInto(out)
	return out
}

// DeepCopyInto copies all properties from this object into another object
func (l *LocalBuildList) DeepCopyInto(out *LocalBuildList) {
	*out = *l
	out.TypeMeta = l.TypeMeta
	l.ListMeta.DeepCopyInto(&out.ListMeta)
	if l.Items != nil {
		out.Items = make([]LocalBuild, len(l.Items))
		for i := range l.Items {
			l.Items[i].DeepCopyInto(&out.Items[i])
		}
	}
}

func init() {
	SchemeBuilder.Register(&LocalBuild{}, &LocalBuildList{})
}