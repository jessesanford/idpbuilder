package multistage

import (
	"context"
	"io"
)

// MultistageBuilder provides multi-stage Dockerfile build capabilities using Buildah
type MultistageBuilder struct {
	parser       *DockerfileParser
	stageManager *StageManager
}

// BuildStage represents a single stage in a multi-stage Dockerfile
type BuildStage struct {
	Name         string            `json:"name"`
	BaseImage    string            `json:"baseImage"`
	Dependencies []string          `json:"dependencies"`
	Commands     []Command         `json:"commands"`
	BuildArgs    map[string]string `json:"buildArgs"`
	Context      *BuildContext     `json:"context,omitempty"`
}

// Command represents a Dockerfile command within a stage
type Command struct {
	Type string   `json:"type"`     // RUN, COPY, ADD, etc.
	Args []string `json:"args"`     // Command arguments
	From string   `json:"from"`     // Source stage for COPY --from
}

// StageGraph represents the dependency graph of all stages
type StageGraph struct {
	Stages         []BuildStage      `json:"stages"`
	Dependencies   map[string][]string `json:"dependencies"`
	ExecutionOrder []string          `json:"executionOrder"`
}

// BuildContext contains the build context information for a stage
type BuildContext struct {
	WorkingDir string            `json:"workingDir"`
	Files      map[string][]byte `json:"files"`
	Args       map[string]string `json:"args"`
}

// BuildOptions contains options for multi-stage builds
type BuildOptions struct {
	Target        string            `json:"target,omitempty"`        // Target stage to build
	BuildArgs     map[string]string `json:"buildArgs,omitempty"`     // Build-time variables
	ContextDir    string            `json:"contextDir"`              // Build context directory
	Dockerfile    string            `json:"dockerfile"`              // Path to Dockerfile
	CacheEnabled  bool             `json:"cacheEnabled"`            // Enable intermediate image caching
	Output        io.Writer        `json:"-"`                       // Build output destination
}

// MultistageBuilderInterface defines the contract for multi-stage builders
type MultistageBuilderInterface interface {
	// Build executes a multi-stage build
	Build(ctx context.Context, options *BuildOptions) error
	
	// ParseDockerfile parses a multi-stage Dockerfile
	ParseDockerfile(dockerfile io.Reader) (*StageGraph, error)
	
	// GetStageGraph returns the stage dependency graph
	GetStageGraph() *StageGraph
	
	// SetBuildArgs sets build-time arguments
	SetBuildArgs(args map[string]string)
}