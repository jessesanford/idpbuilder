package api

import (
	"fmt"
	"time"
)

// Stage represents a build stage in a multi-stage Dockerfile.
type Stage struct {
	Name               string            `json:"name" yaml:"name"`
	BaseImage          string            `json:"base_image" yaml:"base_image"`
	Instructions       []string          `json:"instructions" yaml:"instructions"`
	Dependencies       []string          `json:"dependencies" yaml:"dependencies"`
	Cacheable          bool              `json:"cacheable" yaml:"cacheable"`
	Size               int64             `json:"size,omitempty" yaml:"size,omitempty"`
	BuildArgs          map[string]string `json:"build_args,omitempty" yaml:"build_args,omitempty"`
	Platform           string            `json:"platform,omitempty" yaml:"platform,omitempty"`
	EstimatedBuildTime time.Duration     `json:"estimated_build_time,omitempty" yaml:"estimated_build_time,omitempty"`
	Tags               []string          `json:"tags,omitempty" yaml:"tags,omitempty"`
}

// Validate checks if the stage configuration is valid.
func (s *Stage) Validate() error {
	if s.Name == "" {
		return fmt.Errorf("stage name cannot be empty")
	}
	if s.BaseImage == "" {
		return fmt.Errorf("stage base image cannot be empty")
	}
	if len(s.Instructions) == 0 {
		return fmt.Errorf("stage must have at least one instruction")
	}
	return nil
}

// HasDependency checks if this stage depends on another stage.
func (s *Stage) HasDependency(stageName string) bool {
	for _, dep := range s.Dependencies {
		if dep == stageName {
			return true
		}
	}
	return false
}

// StageResult contains the results of building a stage.
type StageResult struct {
	StageID       string                `json:"stage_id" yaml:"stage_id"`
	ImageID       string                `json:"image_id" yaml:"image_id"`
	Layers        []*Layer              `json:"layers" yaml:"layers"`
	BuildTime     time.Duration         `json:"build_time" yaml:"build_time"`
	CacheHit      bool                  `json:"cache_hit" yaml:"cache_hit"`
	Error         string                `json:"error,omitempty" yaml:"error,omitempty"`
	StartTime     time.Time             `json:"start_time" yaml:"start_time"`
	EndTime       time.Time             `json:"end_time" yaml:"end_time"`
	ResourceUsage *StageResourceUsage   `json:"resource_usage,omitempty" yaml:"resource_usage,omitempty"`
	Warnings      []string              `json:"warnings,omitempty" yaml:"warnings,omitempty"`
}

// IsSuccess returns true if the stage build completed without errors.
func (sr *StageResult) IsSuccess() bool {
	return sr.Error == ""
}

// StageResourceUsage tracks resource consumption during stage building.
type StageResourceUsage struct {
	PeakMemoryMB     int64   `json:"peak_memory_mb" yaml:"peak_memory_mb"`
	AverageCPU       float64 `json:"average_cpu" yaml:"average_cpu"`
	DiskIOBytes      int64   `json:"disk_io_bytes" yaml:"disk_io_bytes"`
	NetworkIOBytes   int64   `json:"network_io_bytes" yaml:"network_io_bytes"`
}

// DependencyGraph represents stage dependencies and execution order.
type DependencyGraph struct {
	Nodes                map[string]*Stage `json:"nodes" yaml:"nodes"`
	Edges                map[string][]string `json:"edges" yaml:"edges"`
	Parallel             [][]string        `json:"parallel" yaml:"parallel"`
	ExecutionOrder       []string          `json:"execution_order" yaml:"execution_order"`
	CriticalPath         []string          `json:"critical_path" yaml:"critical_path"`
	EstimatedTotalTime   time.Duration     `json:"estimated_total_time" yaml:"estimated_total_time"`
}

// Validate checks if the dependency graph is valid.
func (dg *DependencyGraph) Validate() error {
	if len(dg.Nodes) == 0 {
		return fmt.Errorf("dependency graph must have at least one node")
	}

	for node, deps := range dg.Edges {
		if _, exists := dg.Nodes[node]; !exists {
			return fmt.Errorf("edge references non-existent node: %s", node)
		}
		for _, dep := range deps {
			if _, exists := dg.Nodes[dep]; !exists {
				return fmt.Errorf("dependency references non-existent node: %s", dep)
			}
		}
	}

	return nil
}

// GetRootNodes returns stages with no dependencies.
func (dg *DependencyGraph) GetRootNodes() []string {
	var roots []string
	for name := range dg.Nodes {
		if len(dg.Edges[name]) == 0 {
			roots = append(roots, name)
		}
	}
	return roots
}

// GetLeafNodes returns stages that no other stages depend on.
func (dg *DependencyGraph) GetLeafNodes() []string {
	dependedUpon := make(map[string]bool)
	for _, deps := range dg.Edges {
		for _, dep := range deps {
			dependedUpon[dep] = true
		}
	}

	var leaves []string
	for name := range dg.Nodes {
		if !dependedUpon[name] {
			leaves = append(leaves, name)
		}
	}
	return leaves
}

// BuildError represents errors that can occur during image building.
type BuildError struct {
	Stage       string                 `json:"stage,omitempty" yaml:"stage,omitempty"`
	Instruction string                 `json:"instruction,omitempty" yaml:"instruction,omitempty"`
	Message     string                 `json:"message" yaml:"message"`
	Code        string                 `json:"code,omitempty" yaml:"code,omitempty"`
	Details     map[string]interface{} `json:"details,omitempty" yaml:"details,omitempty"`
	Timestamp   time.Time              `json:"timestamp" yaml:"timestamp"`
	Recoverable bool                   `json:"recoverable" yaml:"recoverable"`
}

// Error implements the error interface.
func (be *BuildError) Error() string {
	if be.Stage != "" {
		return fmt.Sprintf("stage %s: %s", be.Stage, be.Message)
	}
	return be.Message
}

// BuildContext contains the context information for a build operation.
type BuildContext struct {
	ContextPath    string   `json:"context_path" yaml:"context_path"`
	DockerfilePath string   `json:"dockerfile_path" yaml:"dockerfile_path"`
	Excludes       []string `json:"excludes,omitempty" yaml:"excludes,omitempty"`
	Size           int64    `json:"size" yaml:"size"`
	FileCount      int64    `json:"file_count" yaml:"file_count"`
	Hash           string   `json:"hash" yaml:"hash"`
}

// ProgressEvent represents a build progress event.
type ProgressEvent struct {
	ID        string                 `json:"id" yaml:"id"`
	Type      string                 `json:"type" yaml:"type"`
	Stage     string                 `json:"stage,omitempty" yaml:"stage,omitempty"`
	Message   string                 `json:"message" yaml:"message"`
	Progress  int                    `json:"progress,omitempty" yaml:"progress,omitempty"`
	Timestamp time.Time              `json:"timestamp" yaml:"timestamp"`
	Details   map[string]interface{} `json:"details,omitempty" yaml:"details,omitempty"`
}

// Constants for common values and configurations.
const (
	// MediaType constants for OCI images
	MediaTypeImageManifest     = "application/vnd.oci.image.manifest.v1+json"
	MediaTypeImageIndex        = "application/vnd.oci.image.index.v1+json"
	MediaTypeImageConfig       = "application/vnd.oci.image.config.v1+json"
	MediaTypeImageLayerTarGzip = "application/vnd.oci.image.layer.v1.tar+gzip"

	// Docker compatibility media types
	MediaTypeDockerManifest     = "application/vnd.docker.distribution.manifest.v2+json"
	MediaTypeDockerManifestList = "application/vnd.docker.distribution.manifest.list.v2+json"
	MediaTypeDockerConfig       = "application/vnd.docker.container.image.v1+json"
	MediaTypeDockerLayerTarGzip = "application/vnd.docker.image.rootfs.diff.tar.gzip"

	// Build event types
	EventTypeBuildStart    = "build_start"
	EventTypeStageStart    = "stage_start"
	EventTypeStageProgress = "stage_progress"
	EventTypeStageComplete = "stage_complete"
	EventTypeBuildComplete = "build_complete"
	EventTypeError         = "error"

	// Error codes
	ErrorCodeInvalidDockerfile = "INVALID_DOCKERFILE"
	ErrorCodeMissingBase       = "MISSING_BASE_IMAGE"
	ErrorCodeInstructionFailed = "INSTRUCTION_FAILED"
	ErrorCodeResourceLimit     = "RESOURCE_LIMIT_EXCEEDED"
	ErrorCodeTimeout           = "BUILD_TIMEOUT"
	ErrorCodeCacheError        = "CACHE_ERROR"
	ErrorCodeRegistryError     = "REGISTRY_ERROR"

	// Cache constants
	DefaultCacheExpiry    = 24 * time.Hour
	DefaultMaxCacheSize   = 10 * 1024 * 1024 * 1024 // 10GB
	DefaultPruneThreshold = 0.8                      // 80%

	// Build constants
	DefaultBuildTimeout = 30 * time.Minute
	DefaultMaxParallel  = 4
	MaxDockerfileSize   = 1024 * 1024           // 1MB
	MaxContextSize      = 100 * 1024 * 1024 * 1024 // 100GB
)