// Package api defines core types for OCI operations in IDPBuilder.
// This file contains configuration structures, request/response types,
// and operational data structures for OCI build and registry operations.
package api

import (
	"time"
)

// BuildConfig holds the configuration for OCI build operations.
// It encapsulates buildah-specific settings and runtime parameters.
type BuildConfig struct {
	// StorageDriver specifies the storage driver for buildah operations.
	// Common values: "vfs", "overlay", "btrfs"
	StorageDriver string `json:"storage_driver" validate:"required,oneof=vfs overlay btrfs zfs"`

	// StorageOptions contains driver-specific storage options.
	StorageOptions map[string]string `json:"storage_options,omitempty"`

	// RuntimePath specifies the path to the container runtime.
	RuntimePath string `json:"runtime_path" validate:"required"`

	// RunRoot defines the directory for runtime state.
	RunRoot string `json:"run_root" validate:"required"`

	// GraphRoot defines the directory for storage state.
	GraphRoot string `json:"graph_root" validate:"required"`

	// Rootless indicates whether to run in rootless mode.
	Rootless bool `json:"rootless"`

	// DefaultMountsFilePath specifies the path to default mounts configuration.
	DefaultMountsFilePath string `json:"default_mounts_file_path,omitempty"`

	// SignaturePolicyPath specifies the path to signature policy configuration.
	SignaturePolicyPath string `json:"signature_policy_path,omitempty"`

	// NetworkConfigPath specifies the path to network configuration.
	NetworkConfigPath string `json:"network_config_path,omitempty"`

	// MaxParallelBuilds limits the number of concurrent build operations.
	MaxParallelBuilds int `json:"max_parallel_builds" validate:"min=1,max=10"`

	// BuildTimeout specifies the maximum duration for a build operation.
	BuildTimeout time.Duration `json:"build_timeout" validate:"min=1m"`

	// LogLevel controls the verbosity of build operations.
	LogLevel string `json:"log_level" validate:"required,oneof=debug info warn error"`

	// CacheDir specifies the directory for build cache storage.
	CacheDir string `json:"cache_dir,omitempty"`

	// TempDir specifies the directory for temporary build files.
	TempDir string `json:"temp_dir,omitempty"`
}

// RegistryConfig defines configuration for OCI registry operations.
type RegistryConfig struct {
	// URL is the base URL of the OCI registry.
	URL string `json:"url" validate:"required,url"`

	// Username for registry authentication.
	Username string `json:"username,omitempty"`

	// Password for registry authentication.
	Password string `json:"password,omitempty"`

	// Token for token-based authentication.
	Token string `json:"token,omitempty"`

	// TLSVerify controls TLS certificate verification.
	TLSVerify bool `json:"tls_verify"`

	// CertDir specifies the directory containing TLS certificates.
	CertDir string `json:"cert_dir,omitempty"`

	// Insecure allows insecure HTTP connections.
	Insecure bool `json:"insecure"`

	// Timeout specifies the timeout for registry operations.
	Timeout time.Duration `json:"timeout" validate:"min=10s"`

	// Retry configuration for failed operations.
	MaxRetries int           `json:"max_retries" validate:"min=0,max=10"`
	RetryDelay time.Duration `json:"retry_delay" validate:"min=1s"`

	// UserAgent string for registry requests.
	UserAgent string `json:"user_agent,omitempty"`
}

// StackOCIConfig defines stack-specific OCI configuration.
type StackOCIConfig struct {
	// StackName is the unique identifier for this stack.
	StackName string `json:"stack_name" validate:"required,min=1,max=100"`

	// Version specifies the stack version using semantic versioning.
	Version string `json:"version" validate:"required,semver"`

	// Description provides a human-readable description of the stack.
	Description string `json:"description,omitempty" validate:"max=500"`

	// Labels are key-value pairs applied to the built images.
	Labels map[string]string `json:"labels,omitempty"`

	// Annotations are key-value pairs for image annotations.
	Annotations map[string]string `json:"annotations,omitempty"`

	// Platform specifies the target platform (os/arch).
	Platform string `json:"platform" validate:"required,platform"`

	// BaseImage specifies the base image for this stack.
	BaseImage string `json:"base_image" validate:"required"`

	// Maintainer information for the stack images.
	Maintainer string `json:"maintainer,omitempty" validate:"max=200"`

	// Repository specifies the target repository for pushing images.
	Repository string `json:"repository" validate:"required"`

	// Tags defines the tags to apply to built images.
	Tags []string `json:"tags" validate:"required,min=1,dive,required,image_tag"`

	// BuildArgs provides build-time variables.
	BuildArgs map[string]string `json:"build_args,omitempty"`

	// Dockerfile specifies the path to the Dockerfile.
	Dockerfile string `json:"dockerfile" validate:"required"`

	// ContextDir specifies the build context directory.
	ContextDir string `json:"context_dir" validate:"required"`

	// IgnoreFile specifies the path to .dockerignore file.
	IgnoreFile string `json:"ignore_file,omitempty"`

	// Created timestamp for stack configuration.
	Created time.Time `json:"created"`

	// Updated timestamp for stack configuration.
	Updated time.Time `json:"updated"`
}

// BuildRequest represents a request to build an OCI image.
type BuildRequest struct {
	// ID is a unique identifier for this build request.
	ID string `json:"id" validate:"required"`

	// Dockerfile specifies the path to the Dockerfile.
	Dockerfile string `json:"dockerfile" validate:"required"`

	// ContextDir specifies the build context directory.
	ContextDir string `json:"context_dir" validate:"required"`

	// Tags specifies the tags to apply to the built image.
	Tags []string `json:"tags" validate:"required,min=1,dive,required,image_tag"`

	// Platform specifies the target platform.
	Platform string `json:"platform,omitempty" validate:"platform"`

	// BuildArgs provides build-time variables.
	BuildArgs map[string]string `json:"build_args,omitempty"`

	// Labels provides image labels.
	Labels map[string]string `json:"labels,omitempty"`

	// Target specifies the target stage in multi-stage builds.
	Target string `json:"target,omitempty"`

	// NoCache disables build cache usage.
	NoCache bool `json:"no_cache"`

	// Pull forces pulling of base images.
	Pull bool `json:"pull"`

	// SquashLayers combines all layers into a single layer.
	SquashLayers bool `json:"squash_layers"`

	// Created timestamp for the request.
	Created time.Time `json:"created"`
}

// BuildResult represents the result of a build operation.
type BuildResult struct {
	// BuildID identifies the build operation.
	BuildID string `json:"build_id"`

	// ImageID is the ID of the built image.
	ImageID string `json:"image_id"`

	// Digest is the content digest of the built image.
	Digest string `json:"digest"`

	// Tags are the tags applied to the image.
	Tags []string `json:"tags"`

	// Size is the size of the built image in bytes.
	Size int64 `json:"size"`

	// Duration is the time taken to complete the build.
	Duration time.Duration `json:"duration"`

	// Layers contains information about image layers.
	Layers []*LayerInfo `json:"layers,omitempty"`

	// Warnings contains any warnings from the build process.
	Warnings []string `json:"warnings,omitempty"`

	// Created timestamp for the build result.
	Created time.Time `json:"created"`
}

// BuildStatus represents the current status of a build operation.
type BuildStatus struct {
	// BuildID identifies the build operation.
	BuildID string `json:"build_id"`

	// Status indicates the current build phase.
	Status BuildPhase `json:"status"`

	// Progress indicates build completion percentage (0-100).
	Progress int `json:"progress" validate:"min=0,max=100"`

	// CurrentStep describes the current build step.
	CurrentStep string `json:"current_step,omitempty"`

	// StartTime is when the build started.
	StartTime time.Time `json:"start_time"`

	// EndTime is when the build completed (if finished).
	EndTime *time.Time `json:"end_time,omitempty"`

	// Error contains error information if the build failed.
	Error string `json:"error,omitempty"`

	// LogPath specifies the path to detailed build logs.
	LogPath string `json:"log_path,omitempty"`
}

// BuildPhase represents the current phase of a build operation.
type BuildPhase string

const (
	// BuildPhaseInitializing indicates the build is being set up.
	BuildPhaseInitializing BuildPhase = "initializing"

	// BuildPhaseDownloading indicates base images are being pulled.
	BuildPhaseDownloading BuildPhase = "downloading"

	// BuildPhaseBuilding indicates the image is being built.
	BuildPhaseBuilding BuildPhase = "building"

	// BuildPhaseFinishing indicates final processing is occurring.
	BuildPhaseFinishing BuildPhase = "finishing"

	// BuildPhaseCompleted indicates the build completed successfully.
	BuildPhaseCompleted BuildPhase = "completed"

	// BuildPhaseFailed indicates the build failed.
	BuildPhaseFailed BuildPhase = "failed"

	// BuildPhaseCancelled indicates the build was cancelled.
	BuildPhaseCancelled BuildPhase = "cancelled"
)

// BuildOptions provides additional options for build operations.
type BuildOptions struct {
	// Quiet suppresses build output.
	Quiet bool `json:"quiet"`

	// NoCache disables cache usage.
	NoCache bool `json:"no_cache"`

	// Pull forces pulling of base images.
	Pull bool `json:"pull"`

	// Remove removes intermediate containers after build.
	Remove bool `json:"remove"`

	// ForceRemove removes intermediate containers even on failure.
	ForceRemove bool `json:"force_remove"`

	// Memory sets memory limit for build operations.
	Memory int64 `json:"memory,omitempty" validate:"min=0"`

	// MemorySwap sets memory swap limit.
	MemorySwap int64 `json:"memory_swap,omitempty" validate:"min=0"`

	// CPUShares sets CPU shares for build operations.
	CPUShares int64 `json:"cpu_shares,omitempty" validate:"min=0"`

	// CPUQuota sets CPU quota for build operations.
	CPUQuota int64 `json:"cpu_quota,omitempty" validate:"min=0"`

	// CPUPeriod sets CPU period for build operations.
	CPUPeriod int64 `json:"cpu_period,omitempty" validate:"min=0"`
}

// PushOptions provides options for registry push operations.
type PushOptions struct {
	// All pushes all tags of the image.
	All bool `json:"all"`

	// Compress compresses layers during push.
	Compress bool `json:"compress"`

	// DisableContentTrust disables content trust for push.
	DisableContentTrust bool `json:"disable_content_trust"`

	// Quiet suppresses push output.
	Quiet bool `json:"quiet"`
}

// PullOptions provides options for registry pull operations.
type PullOptions struct {
	// All pulls all tags of the image.
	All bool `json:"all"`

	// DisableContentTrust disables content trust for pull.
	DisableContentTrust bool `json:"disable_content_trust"`

	// Platform specifies the platform for multi-platform images.
	Platform string `json:"platform,omitempty" validate:"platform"`

	// Quiet suppresses pull output.
	Quiet bool `json:"quiet"`
}

// LayerInfo contains information about an image layer.
type LayerInfo struct {
	// Digest is the content digest of the layer.
	Digest string `json:"digest"`

	// Size is the size of the layer in bytes.
	Size int64 `json:"size"`

	// MediaType specifies the layer media type.
	MediaType string `json:"media_type"`

	// Created timestamp for the layer.
	Created time.Time `json:"created"`

	// CreatedBy contains the command that created this layer.
	CreatedBy string `json:"created_by,omitempty"`

	// Comment provides additional information about the layer.
	Comment string `json:"comment,omitempty"`

	// EmptyLayer indicates if this is an empty layer.
	EmptyLayer bool `json:"empty_layer"`
}

// ImageInfo contains detailed information about an OCI image.
type ImageInfo struct {
	// ID is the unique identifier for the image.
	ID string `json:"id"`

	// Digest is the content digest of the image.
	Digest string `json:"digest"`

	// Tags are the tags associated with the image.
	Tags []string `json:"tags"`

	// Size is the total size of the image in bytes.
	Size int64 `json:"size"`

	// Created timestamp for the image.
	Created time.Time `json:"created"`

	// Labels are key-value pairs associated with the image.
	Labels map[string]string `json:"labels,omitempty"`

	// Architecture specifies the image architecture.
	Architecture string `json:"architecture"`

	// OS specifies the target operating system.
	OS string `json:"os"`

	// Layers contains information about image layers.
	Layers []*LayerInfo `json:"layers,omitempty"`
}

// StackImageInfo extends ImageInfo with stack-specific information.
type StackImageInfo struct {
	*ImageInfo

	// StackName is the name of the stack this image belongs to.
	StackName string `json:"stack_name"`

	// StackVersion is the version of the stack.
	StackVersion string `json:"stack_version"`

	// BuildConfig contains the build configuration used.
	BuildConfig *StackOCIConfig `json:"build_config,omitempty"`
}

// StackHistoryEntry represents a single entry in stack build/push history.
type StackHistoryEntry struct {
	// ID uniquely identifies this history entry.
	ID string `json:"id"`

	// StackName is the name of the stack.
	StackName string `json:"stack_name"`

	// Version is the stack version for this entry.
	Version string `json:"version"`

	// Action describes what action was performed.
	Action string `json:"action" validate:"required,oneof=build push update delete"`

	// ImageID is the ID of the affected image.
	ImageID string `json:"image_id"`

	// Status indicates the result of the action.
	Status string `json:"status" validate:"required,oneof=success failed cancelled"`

	// Timestamp records when this action occurred.
	Timestamp time.Time `json:"timestamp"`

	// Duration records how long the action took.
	Duration time.Duration `json:"duration"`

	// Details provides additional information about the action.
	Details map[string]interface{} `json:"details,omitempty"`
}

// ProgressEvent represents a progress update during build or push operations.
type ProgressEvent struct {
	// ID uniquely identifies the operation being tracked.
	ID string `json:"id"`

	// Action describes the current action being performed.
	Action string `json:"action"`

	// Status provides the current status of the action.
	Status string `json:"status"`

	// Progress indicates completion percentage (0-100).
	Progress int `json:"progress" validate:"min=0,max=100"`

	// Total indicates the total units of work (bytes, steps, etc).
	Total int64 `json:"total,omitempty"`

	// Current indicates the current units completed.
	Current int64 `json:"current,omitempty"`

	// Message provides a human-readable progress message.
	Message string `json:"message,omitempty"`

	// Timestamp records when this progress update occurred.
	Timestamp time.Time `json:"timestamp"`

	// Details provides additional context-specific information.
	Details map[string]interface{} `json:"details,omitempty"`
}