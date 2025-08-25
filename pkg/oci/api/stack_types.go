// Package api provides stack-specific types for OCI operations in IDPBuilder.
// This file contains configuration structures and data types specific to
// stack management, history tracking, and progress reporting.
package api

import (
	"time"
)

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