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

// LocalBuildSpec defines the desired state of LocalBuild
type LocalBuildSpec struct {
	// Source specifies the source configuration for the build
	// +kubebuilder:validation:Required
	Source LocalBuildSource `json:"source"`

	// Builder specifies the build configuration
	// +kubebuilder:validation:Required
	Builder LocalBuildConfiguration `json:"builder"`

	// Output specifies where to push the built artifacts
	// +kubebuilder:validation:Required
	Output LocalBuildOutput `json:"output"`

	// BuildArgs contains build-time variables
	// +optional
	BuildArgs map[string]string `json:"buildArgs,omitempty"`

	// Labels contains labels to apply to the built image
	// +optional
	Labels map[string]string `json:"labels,omitempty"`

	// Annotations contains annotations to apply to the built image
	// +optional
	Annotations map[string]string `json:"annotations,omitempty"`

	// Timeout specifies the timeout for the build process
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:Pattern=^([0-9]+(\.[0-9]+)?(ms|s|m|h))+$
	// +kubebuilder:default="10m"
	Timeout metav1.Duration `json:"timeout,omitempty"`

	// Suspend tells the controller to suspend reconciliation
	// +kubebuilder:default=false
	Suspend bool `json:"suspend,omitempty"`
}

// LocalBuildSource defines the source for a local build
type LocalBuildSource struct {
	// GitRepository references a GitRepository resource
	// +optional
	GitRepository *GitRepositoryRef `json:"gitRepository,omitempty"`

	// Path specifies a local filesystem path
	// +optional
	Path string `json:"path,omitempty"`
}

// LocalBuildConfiguration defines how to build the source
type LocalBuildConfiguration struct {
	// Strategy specifies the build strategy
	// +kubebuilder:validation:Enum=dockerfile;buildpacks;kaniko;custom
	// +kubebuilder:default="dockerfile"
	Strategy BuildStrategy `json:"strategy"`

	// Dockerfile specifies the Dockerfile configuration
	// +optional
	Dockerfile *DockerfileBuild `json:"dockerfile,omitempty"`

	// Resources specifies resource requirements for the build
	// +optional
	Resources *BuildResources `json:"resources,omitempty"`
}

// BuildStrategy represents different build strategies
type BuildStrategy string

const (
	// BuildStrategyDockerfile uses a Dockerfile for building
	BuildStrategyDockerfile BuildStrategy = "dockerfile"
	// BuildStrategyBuildpacks uses Cloud Native Buildpacks
	BuildStrategyBuildpacks BuildStrategy = "buildpacks"
	// BuildStrategyKaniko uses Kaniko for building
	BuildStrategyKaniko BuildStrategy = "kaniko"
	// BuildStrategyCustom uses a custom build process
	BuildStrategyCustom BuildStrategy = "custom"
)

// DockerfileBuild defines Dockerfile-based build configuration
type DockerfileBuild struct {
	// Path is the path to the Dockerfile relative to the source root
	// +kubebuilder:default="Dockerfile"
	Path string `json:"path,omitempty"`

	// Context is the build context path relative to the source root
	// +kubebuilder:default="."
	Context string `json:"context,omitempty"`

	// Target specifies the target stage in a multi-stage Dockerfile
	// +optional
	Target string `json:"target,omitempty"`
}


// BuildResources defines resource requirements for builds
type BuildResources struct {
	// CPU specifies CPU requirements
	// +optional
	CPU string `json:"cpu,omitempty"`

	// Memory specifies memory requirements
	// +optional
	Memory string `json:"memory,omitempty"`

	// Storage specifies storage requirements
	// +optional
	Storage string `json:"storage,omitempty"`
}

// LocalBuildOutput defines where to push built artifacts
type LocalBuildOutput struct {
	// Registry specifies the container registry configuration
	// +kubebuilder:validation:Required
	Registry LocalBuildRegistry `json:"registry"`
}

// LocalBuildRegistry defines container registry configuration
type LocalBuildRegistry struct {
	// URL is the registry URL
	// +kubebuilder:validation:Required
	URL string `json:"url"`

	// Repository is the repository within the registry
	// +kubebuilder:validation:Required
	Repository string `json:"repository"`
}

// LocalBuildStatus defines the observed state of LocalBuild
type LocalBuildStatus struct {
	// Phase represents the current phase of the build
	// +kubebuilder:validation:Enum=Pending;Building;Ready;Failed
	Phase LocalBuildPhase `json:"phase,omitempty"`

	// Conditions contains details about the build state
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// LocalBuildPhase represents the phase of a local build
type LocalBuildPhase string

const (
	// LocalBuildPhasePending means the build is pending
	LocalBuildPhasePending LocalBuildPhase = "Pending"
	// LocalBuildPhaseBuilding means the build is in progress
	LocalBuildPhaseBuilding LocalBuildPhase = "Building"
	// LocalBuildPhaseReady means the build completed successfully
	LocalBuildPhaseReady LocalBuildPhase = "Ready"
	// LocalBuildPhaseFailed means the build failed
	LocalBuildPhaseFailed LocalBuildPhase = "Failed"
)


// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Namespaced,categories=idpbuilder
// +kubebuilder:printcolumn:name="Strategy",type=string,JSONPath=".spec.builder.strategy"
// +kubebuilder:printcolumn:name="Phase",type=string,JSONPath=".status.phase"
// +kubebuilder:printcolumn:name="Age",type=date,JSONPath=".metadata.creationTimestamp"

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

