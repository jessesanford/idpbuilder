// Package api defines the core types for container build operations
package api

import "fmt"

// BuildRequest represents a container build request
type BuildRequest struct {
	// DockerfilePath is the path to the Dockerfile relative to ContextDir
	DockerfilePath string `json:"dockerfilePath"`

	// ContextDir is the build context directory (absolute path)
	ContextDir string `json:"contextDir"`

	// ImageName is the target image name (without registry)
	ImageName string `json:"imageName"`

	// ImageTag is the target image tag
	ImageTag string `json:"imageTag"`
}

// BuildResponse represents the result of a build operation
type BuildResponse struct {
	// ImageID is the built image ID
	ImageID string `json:"imageID"`

	// FullTag is the complete image reference
	// Format: gitea.cnoe.localtest.me/giteaadmin/{imageName}:{imageTag}
	FullTag string `json:"fullTag"`

	// Success indicates if the build completed successfully
	Success bool `json:"success"`

	// Error contains error details if Success is false
	Error string `json:"error,omitempty"`
}

// Validate performs basic validation on BuildRequest
func (br *BuildRequest) Validate() error {
	if br.DockerfilePath == "" {
		return fmt.Errorf("DockerfilePath is required")
	}
	if br.ContextDir == "" {
		return fmt.Errorf("ContextDir is required")
	}
	if br.ImageName == "" {
		return fmt.Errorf("ImageName is required")
	}
	if br.ImageTag == "" {
		br.ImageTag = "latest" // Default tag
	}
	return nil
}