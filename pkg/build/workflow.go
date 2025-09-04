package build

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/cnoe-io/idpbuilder/pkg/builder"
	"github.com/cnoe-io/idpbuilder/pkg/logger"
	v1 "github.com/google/go-containerregistry/pkg/v1"
)

// BuildWorkflow manages the complete build process
type BuildWorkflow struct {
	context    *BuildContext
	builder    *builder.Builder
	logger     *logger.Logger
	stepCount  int
	startTime  time.Time
}

// BuildStep represents a single step in the build workflow
type BuildStep struct {
	Name        string
	Description string
	Function    func(ctx context.Context) error
	Required    bool
	Executed    bool
	Error       error
	Duration    time.Duration
}

// WorkflowOptions contains options for workflow execution
type WorkflowOptions struct {
	SkipSteps    []string
	ContinueOnError bool
	Parallel     bool
	MaxRetries   int
}

// NewBuildWorkflow creates a new build workflow
func NewBuildWorkflow(buildCtx *BuildContext) (*BuildWorkflow, error) {
	if !isCliToolsEnabled() {
		return nil, fmt.Errorf("ENABLE_CLI_TOOLS feature flag is not enabled")
	}

	return &BuildWorkflow{
		context:   buildCtx,
		logger:    logger.New(),
		startTime: time.Now(),
	}, nil
}

// Execute runs the complete build workflow
func (wf *BuildWorkflow) Execute(ctx context.Context, opts *WorkflowOptions) (*BuildResult, error) {
	wf.logger.Info("Starting build workflow")
	wf.logger.Info("Build context: %s", wf.context.ContextPath)
	
	if opts == nil {
		opts = &WorkflowOptions{
			MaxRetries: 3,
		}
	}

	// Define build steps
	steps := wf.defineBuildSteps()

	// Execute steps
	for i, step := range steps {
		if wf.shouldSkipStep(step.Name, opts.SkipSteps) {
			wf.logger.Info("Skipping step %d: %s", i+1, step.Name)
			continue
		}

		if err := wf.executeStep(ctx, &step, opts.MaxRetries); err != nil {
			step.Error = err
			wf.logger.Error("Step %d failed: %s - %v", i+1, step.Name, err)
			
			if step.Required && !opts.ContinueOnError {
				return nil, fmt.Errorf("required step failed: %s - %w", step.Name, err)
			}
		}
	}

	// Create build result
	result := &BuildResult{
		Success:     true,
		Duration:    time.Since(wf.startTime),
		Steps:       steps,
		Context:     wf.context,
		Artifacts:   wf.collectArtifacts(),
		Metadata:    wf.collectMetadata(),
	}

	wf.logger.Info("Build workflow completed in %v", result.Duration)
	return result, nil
}

// defineBuildSteps defines all the steps in the build workflow
func (wf *BuildWorkflow) defineBuildSteps() []BuildStep {
	return []BuildStep{
		{
			Name:        "validate-context",
			Description: "Validate build context and dependencies",
			Function:    wf.validateContext,
			Required:    true,
		},
		{
			Name:        "prepare-environment",
			Description: "Prepare build environment",
			Function:    wf.prepareEnvironment,
			Required:    true,
		},
		{
			Name:        "parse-dockerfile",
			Description: "Parse and analyze Dockerfile",
			Function:    wf.parseDockerfile,
			Required:    true,
		},
		{
			Name:        "build-layers",
			Description: "Build container image layers",
			Function:    wf.buildLayers,
			Required:    true,
		},
		{
			Name:        "create-image",
			Description: "Create final container image",
			Function:    wf.createImage,
			Required:    true,
		},
		{
			Name:        "tag-image",
			Description: "Apply tags to the built image",
			Function:    wf.tagImage,
			Required:    false,
		},
		{
			Name:        "cleanup",
			Description: "Clean up temporary files and resources",
			Function:    wf.cleanup,
			Required:    false,
		},
	}
}

// executeStep executes a single build step with retry logic
func (wf *BuildWorkflow) executeStep(ctx context.Context, step *BuildStep, maxRetries int) error {
	start := time.Now()
	defer func() {
		step.Duration = time.Since(start)
		step.Executed = true
	}()

	wf.logger.Info("Executing step: %s", step.Description)

	var lastErr error
	for attempt := 0; attempt <= maxRetries; attempt++ {
		if attempt > 0 {
			wf.logger.Warn("Retrying step %s (attempt %d/%d)", step.Name, attempt+1, maxRetries+1)
			time.Sleep(time.Second * time.Duration(attempt))
		}

		err := step.Function(ctx)
		if err == nil {
			wf.logger.Info("Step completed: %s", step.Name)
			return nil
		}

		lastErr = err
		if attempt < maxRetries {
			wf.logger.Warn("Step failed, will retry: %s - %v", step.Name, err)
		}
	}

	return fmt.Errorf("step failed after %d retries: %w", maxRetries+1, lastErr)
}

// Build workflow step implementations

func (wf *BuildWorkflow) validateContext(ctx context.Context) error {
	// Check if context directory exists
	if _, err := os.Stat(wf.context.ContextPath); os.IsNotExist(err) {
		return fmt.Errorf("build context does not exist: %s", wf.context.ContextPath)
	}

	// Check if Dockerfile exists
	dockerfilePath := filepath.Join(wf.context.ContextPath, wf.context.DockerfilePath)
	if _, err := os.Stat(dockerfilePath); os.IsNotExist(err) {
		return fmt.Errorf("dockerfile does not exist: %s", dockerfilePath)
	}

	// Validate build options
	if err := wf.context.BuildOptions.Validate(); err != nil {
		return fmt.Errorf("invalid build options: %w", err)
	}

	return nil
}

func (wf *BuildWorkflow) prepareEnvironment(ctx context.Context) error {
	// Create temporary directories
	tempDir, err := os.MkdirTemp("", "idp-build-*")
	if err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}
	
	wf.context.TempDir = tempDir

	// Initialize builder
	builderInstance, err := builder.New(wf.context.BuildOptions)
	if err != nil {
		return fmt.Errorf("failed to create builder: %w", err)
	}
	
	wf.builder = builderInstance

	return nil
}

func (wf *BuildWorkflow) parseDockerfile(ctx context.Context) error {
	// This would parse the Dockerfile and extract instructions
	// For now, we'll use the builder's existing parsing logic
	wf.logger.Info("Parsing Dockerfile: %s", wf.context.DockerfilePath)
	
	// The builder will handle Dockerfile parsing internally
	return nil
}

func (wf *BuildWorkflow) buildLayers(ctx context.Context) error {
	// Create layer manager
	layerManager := builder.NewLayerManager(wf.context.BuildOptions)
	
	// Build layers from context
	layer, err := layerManager.CreateLayerFromDirectory(ctx, wf.context.ContextPath)
	if err != nil {
		return fmt.Errorf("failed to create layer from context: %w", err)
	}

	// Store layer info
	layerInfo, err := layerManager.GetLayerInfo(layer)
	if err != nil {
		return fmt.Errorf("failed to get layer info: %w", err)
	}

	wf.context.Layers = []LayerInfo{*layerInfo}
	return nil
}

func (wf *BuildWorkflow) createImage(ctx context.Context) error {
	// Build the image using the builder
	result, err := wf.builder.Build()
	if err != nil {
		return fmt.Errorf("failed to build image: %w", err)
	}

	wf.context.Result = result
	return nil
}

func (wf *BuildWorkflow) tagImage(ctx context.Context) error {
	if wf.context.Result == nil {
		return fmt.Errorf("no image to tag")
	}

	// Tags are handled by the builder, so this is mostly logging
	if len(wf.context.BuildOptions.Tags) > 0 {
		wf.logger.Info("Applied tags: %v", wf.context.BuildOptions.Tags)
	}

	return nil
}

func (wf *BuildWorkflow) cleanup(ctx context.Context) error {
	// Clean up temporary directories and files
	if wf.context.TempDir != "" {
		if err := os.RemoveAll(wf.context.TempDir); err != nil {
			wf.logger.Warn("Failed to clean up temp directory: %v", err)
			// Don't fail the build for cleanup issues
		}
	}

	return nil
}

// Helper methods

func (wf *BuildWorkflow) shouldSkipStep(stepName string, skipSteps []string) bool {
	for _, skip := range skipSteps {
		if skip == stepName {
			return true
		}
	}
	return false
}

func (wf *BuildWorkflow) collectArtifacts() []Artifact {
	artifacts := []Artifact{}
	
	if wf.context.Result != nil {
		artifacts = append(artifacts, Artifact{
			Type: "container-image",
			Path: wf.context.Result.ImageID,
			Tags: wf.context.Result.Tags,
		})
	}

	return artifacts
}

func (wf *BuildWorkflow) collectMetadata() map[string]interface{} {
	metadata := map[string]interface{}{
		"build_time":    wf.startTime.Format(time.RFC3339),
		"context_path":  wf.context.ContextPath,
		"dockerfile":    wf.context.DockerfilePath,
		"platform":     wf.context.BuildOptions.Platform,
		"step_count":   wf.stepCount,
	}

	if len(wf.context.Layers) > 0 {
		metadata["layer_count"] = len(wf.context.Layers)
	}

	return metadata
}

func isCliToolsEnabled() bool {
	return os.Getenv("ENABLE_CLI_TOOLS") == "true"
}

// Result types

type BuildResult struct {
	Success   bool
	Duration  time.Duration
	Steps     []BuildStep
	Context   *BuildContext
	Artifacts []Artifact
	Metadata  map[string]interface{}
}

type Artifact struct {
	Type string
	Path string
	Tags []string
}

type LayerInfo struct {
	Digest string
	Size   int64
}