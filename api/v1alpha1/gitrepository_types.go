/*
Copyright 2024 The idpbuilder Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GitRepositorySpec defines the desired state of GitRepository
type GitRepositorySpec struct {
	// URL is the Git repository URL
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern=^(https?|git|ssh):\/\/.*$
	URL string `json:"url"`

	// Branch specifies the Git branch to use
	// +kubebuilder:default="main"
	Branch string `json:"branch,omitempty"`

	// Tag specifies a Git tag to use instead of branch
	// +optional
	Tag string `json:"tag,omitempty"`

	// Commit specifies a specific commit hash to use
	// +kubebuilder:validation:Pattern=^[a-f0-9]{40}$
	// +optional
	Commit string `json:"commit,omitempty"`

	// Path specifies a subdirectory within the repository
	// +optional
	Path string `json:"path,omitempty"`

	// Interval defines how often to check for updates
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:Pattern=^([0-9]+(\.[0-9]+)?(ms|s|m|h))+$
	// +kubebuilder:default="5m"
	Interval metav1.Duration `json:"interval,omitempty"`

	// SecretRef specifies a secret containing authentication credentials
	// +optional
	SecretRef *GitSecretReference `json:"secretRef,omitempty"`

	// Timeout specifies the timeout for Git operations
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:Pattern=^([0-9]+(\.[0-9]+)?(ms|s|m|h))+$
	// +kubebuilder:default="60s"
	Timeout metav1.Duration `json:"timeout,omitempty"`

	// Suspend tells the controller to suspend reconciliation
	// +kubebuilder:default=false
	Suspend bool `json:"suspend,omitempty"`
}

// GitSecretReference contains the secret name and keys for Git authentication
type GitSecretReference struct {
	// Name is the name of the secret
	// +kubebuilder:validation:Required
	Name string `json:"name"`

	// UsernameKey is the key in the secret containing the username
	// +kubebuilder:default="username"
	UsernameKey string `json:"usernameKey,omitempty"`

	// PasswordKey is the key in the secret containing the password/token
	// +kubebuilder:default="password"
	PasswordKey string `json:"passwordKey,omitempty"`
}

// GitRepositoryStatus defines the observed state of GitRepository
type GitRepositoryStatus struct {
	// Conditions contains details about the repository state
	// +optional
	// +patchMergeKey=type
	// +patchStrategy=merge
	Conditions []metav1.Condition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type"`

	// LastUpdateTime is when the repository was last fetched
	// +optional
	LastUpdateTime *metav1.Time `json:"lastUpdateTime,omitempty"`

	// Revision is the current revision (commit hash) of the repository
	// +optional
	Revision string `json:"revision,omitempty"`

	// Artifact contains the details of the latest fetched artifact
	// +optional
	Artifact *GitArtifact `json:"artifact,omitempty"`

	// ObservedGeneration is the generation observed by the controller
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`

	// Size is the size of the fetched repository in bytes
	// +optional
	Size int64 `json:"size,omitempty"`

	// ContentConfigMapRef references a ConfigMap containing the repository content
	// +optional
	ContentConfigMapRef *GitContentReference `json:"contentConfigMapRef,omitempty"`
}

// GitArtifact represents a fetched Git repository artifact
type GitArtifact struct {
	// Path is the local path where the artifact is stored
	// +kubebuilder:validation:Required
	Path string `json:"path"`

	// URL is the URL of the artifact
	// +kubebuilder:validation:Required
	URL string `json:"url"`

	// Revision is the revision (commit hash) of the artifact
	// +kubebuilder:validation:Required
	Revision string `json:"revision"`

	// Checksum is the checksum of the artifact
	// +optional
	Checksum string `json:"checksum,omitempty"`

	// Size is the size of the artifact in bytes
	// +optional
	Size int64 `json:"size,omitempty"`

	// LastUpdateTime is when the artifact was last updated
	// +optional
	LastUpdateTime *metav1.Time `json:"lastUpdateTime,omitempty"`
}

// GitContentReference references content stored in a ConfigMap
type GitContentReference struct {
	// Name is the name of the ConfigMap
	// +kubebuilder:validation:Required
	Name string `json:"name"`

	// Key is the key in the ConfigMap containing the content
	// +kubebuilder:default="content.tar.gz"
	Key string `json:"key,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Namespaced,categories=idpbuilder
// +kubebuilder:printcolumn:name="URL",type=string,JSONPath=".spec.url"
// +kubebuilder:printcolumn:name="Branch",type=string,JSONPath=".spec.branch"
// +kubebuilder:printcolumn:name="Revision",type=string,JSONPath=".status.revision"
// +kubebuilder:printcolumn:name="Ready",type=string,JSONPath=".status.conditions[?(@.type==\"Ready\")].status"
// +kubebuilder:printcolumn:name="Age",type=date,JSONPath=".metadata.creationTimestamp"

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

