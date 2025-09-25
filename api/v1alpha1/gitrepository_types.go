package v1alpha1

import (
	"errors"
	"strings"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// SecretReference contains the name of a secret in the same namespace
type SecretReference struct {
	// Name of the secret
	Name string `json:"name"`
}

// GitRepositorySpec defines the desired state of GitRepository
type GitRepositorySpec struct {
	// URL specifies the Git repository URL
	URL string `json:"url"`

	// Branch specifies the Git branch to use
	Branch string `json:"branch,omitempty"`

	// Tag specifies the Git tag to use
	Tag string `json:"tag,omitempty"`

	// Commit specifies the Git commit to use
	Commit string `json:"commit,omitempty"`

	// Path specifies the path within the repository
	Path string `json:"path,omitempty"`

	// Credentials specifies authentication credentials
	Credentials *GitCredentials `json:"credentials,omitempty"`
}

// GitCredentials defines authentication credentials for Git repositories
type GitCredentials struct {
	// SecretRef references a secret containing credentials
	SecretRef *SecretReference `json:"secretRef,omitempty"`

	// Username for basic authentication
	Username string `json:"username,omitempty"`

	// Token for token-based authentication
	Token string `json:"token,omitempty"`
}

// GitRepositoryStatus defines the observed state of GitRepository
type GitRepositoryStatus struct {
	// LastSyncedCommit represents the last synced commit SHA
	LastSyncedCommit string `json:"lastSyncedCommit,omitempty"`

	// Conditions represent the latest available observations
	Conditions []metav1.Condition `json:"conditions,omitempty"`

	// Message provides additional information about the current state
	Message string `json:"message,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// GitRepository is the Schema for the gitrepositories API
type GitRepository struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   GitRepositorySpec   `json:"spec,omitempty"`
	Status GitRepositoryStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// GitRepositoryList contains a list of GitRepository
type GitRepositoryList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []GitRepository `json:"items"`
}

// Validate validates the GitRepository specification
func (g *GitRepository) Validate() error {
	if g.Spec.URL == "" {
		return errors.New("repository URL is required")
	}

	// Validate URL format
	if !strings.HasPrefix(g.Spec.URL, "http://") &&
	   !strings.HasPrefix(g.Spec.URL, "https://") &&
	   !strings.HasPrefix(g.Spec.URL, "git@") {
		return errors.New("repository URL must be a valid Git URL")
	}

	// Validate that only one of branch, tag, or commit is specified
	refCount := 0
	if g.Spec.Branch != "" {
		refCount++
	}
	if g.Spec.Tag != "" {
		refCount++
	}
	if g.Spec.Commit != "" {
		refCount++
	}

	if refCount > 1 {
		return errors.New("only one of branch, tag, or commit can be specified")
	}

	return nil
}

// GetURL returns the repository URL
func (g *GitRepository) GetURL() string {
	return g.Spec.URL
}

// DeepCopyObject returns a deep copy of the GitRepository
func (g *GitRepository) DeepCopyObject() runtime.Object {
	if c := g.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopy creates a deep copy of the GitRepository
func (g *GitRepository) DeepCopy() *GitRepository {
	if g == nil {
		return nil
	}
	out := new(GitRepository)
	g.DeepCopyInto(out)
	return out
}

// DeepCopyInto copies all properties from this object into another object
func (g *GitRepository) DeepCopyInto(out *GitRepository) {
	*out = *g
	out.TypeMeta = g.TypeMeta
	g.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = g.Spec
	if g.Spec.Credentials != nil {
		out.Spec.Credentials = new(GitCredentials)
		*out.Spec.Credentials = *g.Spec.Credentials
		if g.Spec.Credentials.SecretRef != nil {
			out.Spec.Credentials.SecretRef = new(SecretReference)
			*out.Spec.Credentials.SecretRef = *g.Spec.Credentials.SecretRef
		}
	}
	out.Status = g.Status
	if g.Status.Conditions != nil {
		out.Status.Conditions = make([]metav1.Condition, len(g.Status.Conditions))
		for i := range g.Status.Conditions {
			g.Status.Conditions[i].DeepCopyInto(&out.Status.Conditions[i])
		}
	}
}

// DeepCopyObject returns a deep copy of the GitRepositoryList
func (g *GitRepositoryList) DeepCopyObject() runtime.Object {
	if c := g.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopy creates a deep copy of the GitRepositoryList
func (g *GitRepositoryList) DeepCopy() *GitRepositoryList {
	if g == nil {
		return nil
	}
	out := new(GitRepositoryList)
	g.DeepCopyInto(out)
	return out
}

// DeepCopyInto copies all properties from this object into another object
func (g *GitRepositoryList) DeepCopyInto(out *GitRepositoryList) {
	*out = *g
	out.TypeMeta = g.TypeMeta
	g.ListMeta.DeepCopyInto(&out.ListMeta)
	if g.Items != nil {
		out.Items = make([]GitRepository, len(g.Items))
		for i := range g.Items {
			g.Items[i].DeepCopyInto(&out.Items[i])
		}
	}
}

func init() {
	SchemeBuilder.Register(&GitRepository{}, &GitRepositoryList{})
}