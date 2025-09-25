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

// CustomPackageSpec defines the desired state of CustomPackage
type CustomPackageSpec struct {
	// GitRepository references a GitRepository resource
	// +kubebuilder:validation:Required
	GitRepository GitRepositoryRef `json:"gitRepository"`

	// LocalBuild references a LocalBuild resource
	// +optional
	LocalBuild *LocalBuildRef `json:"localBuild,omitempty"`

	// Version specifies the version of the package
	// +optional
	Version string `json:"version,omitempty"`
}

// GitRepositoryRef references a GitRepository resource
type GitRepositoryRef struct {
	// Name is the name of the GitRepository resource
	// +kubebuilder:validation:Required
	Name string `json:"name"`

	// Namespace is the namespace of the GitRepository resource
	// If not specified, the same namespace as the CustomPackage is used
	// +optional
	Namespace string `json:"namespace,omitempty"`
}

// LocalBuildRef references a LocalBuild resource
type LocalBuildRef struct {
	// Name is the name of the LocalBuild resource
	// +kubebuilder:validation:Required
	Name string `json:"name"`

	// Namespace is the namespace of the LocalBuild resource
	// If not specified, the same namespace as the CustomPackage is used
	// +optional
	Namespace string `json:"namespace,omitempty"`
}


// CustomPackageStatus defines the observed state of CustomPackage
type CustomPackageStatus struct {
	// Phase represents the current phase of the package installation
	// +kubebuilder:validation:Enum=Pending;Installing;Ready;Failed;Updating
	Phase CustomPackagePhase `json:"phase,omitempty"`

	// Conditions contains details about the package state
	// +optional
	// +patchMergeKey=type
	// +patchStrategy=merge
	Conditions []metav1.Condition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type"`

	// InstalledVersion is the currently installed version
	// +optional
	InstalledVersion string `json:"installedVersion,omitempty"`

	// LastUpdateTime is when the package was last updated
	// +optional
	LastUpdateTime *metav1.Time `json:"lastUpdateTime,omitempty"`

	// ObservedGeneration is the generation observed by the controller
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`
}

// CustomPackagePhase represents the phase of package installation
type CustomPackagePhase string

const (
	// CustomPackagePhasePending means the package installation is pending
	CustomPackagePhasePending CustomPackagePhase = "Pending"
	// CustomPackagePhaseInstalling means the package is being installed
	CustomPackagePhaseInstalling CustomPackagePhase = "Installing"
	// CustomPackagePhaseReady means the package is ready
	CustomPackagePhaseReady CustomPackagePhase = "Ready"
	// CustomPackagePhaseFailed means the package installation failed
	CustomPackagePhaseFailed CustomPackagePhase = "Failed"
	// CustomPackagePhaseUpdating means the package is being updated
	CustomPackagePhaseUpdating CustomPackagePhase = "Updating"
)

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Namespaced,categories=idpbuilder
// +kubebuilder:printcolumn:name="Phase",type=string,JSONPath=".status.phase"
// +kubebuilder:printcolumn:name="Version",type=string,JSONPath=".status.installedVersion"
// +kubebuilder:printcolumn:name="Age",type=date,JSONPath=".metadata.creationTimestamp"

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

