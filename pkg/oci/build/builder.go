// Package build provides OCI image building orchestration using buildah.
package build

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/containers/buildah/define"
	"github.com/containers/buildah/imagebuildah"
	"github.com/containers/storage"

	// Import Phase 1 types and interfaces
	api "github.com/cnoe-io/idpbuilder/pkg/oci/api"
)

// Builder implements OCIBuildService interface for image building.
type Builder struct {
	config      *api.BuildConfig
	store       storage.Store
	builds      map[string]*api.BuildStatus
	buildsMutex sync.RWMutex
	initialized bool
}

// Ensure Builder implements OCIBuildService interface
var _ api.OCIBuildService = (*Builder)(nil)

// NewBuilder creates a new Builder instance.
func NewBuilder(config *api.BuildConfig) (*Builder, error) {
	if config == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}

	// Set defaults
	if config.MaxParallelBuilds <= 0 {
		config.MaxParallelBuilds = 3
	}
	if config.BuildTimeout <= 0 {
		config.BuildTimeout = 30 * time.Minute
	}

	return &Builder{
		config: config,
		builds: make(map[string]*api.BuildStatus),
	}, nil
}

// Initialize initializes the builder with storage.
func (b *Builder) Initialize(ctx context.Context, config *api.BuildConfig) error {
	if config != nil {
		b.config = config
	}

	if err := b.ValidateConfig(b.config); err != nil {
		return fmt.Errorf("invalid config: %w", err)
	}

	// Initialize storage
	storeOptions, err := storage.DefaultStoreOptions(false, 0)
	if err != nil {
		return fmt.Errorf("failed to get store options: %w", err)
	}

	storeOptions.GraphDriverName = b.config.StorageDriver
	storeOptions.RunRoot = b.config.RunRoot
	storeOptions.GraphRoot = b.config.GraphRoot

	store, err := storage.GetStore(storeOptions)
	if err != nil {
		return fmt.Errorf("failed to get store: %w", err)
	}

	b.store = store
	b.initialized = true
	return nil
}

// BuildImage builds an image from a build request.
func (b *Builder) BuildImage(ctx context.Context, req *api.BuildRequest) (*api.BuildResult, error) {
	if !b.initialized {
		return nil, fmt.Errorf("builder not initialized")
	}

	if err := b.validateRequest(req); err != nil {
		return nil, err
	}

	// Create build status
	status := &api.BuildStatus{
		BuildID:   req.ID,
		Status:    api.BuildPhaseInitializing,
		StartTime: time.Now(),
	}
	b.trackBuild(status)
	defer b.untrackBuild(req.ID)

	// Execute build
	start := time.Now()
	imageID, err := b.executeBuild(ctx, req, status)
	duration := time.Since(start)

	if err != nil {
		status.Status = api.BuildPhaseFailed
		status.Error = err.Error()
		endTime := time.Now()
		status.EndTime = &endTime
		return nil, err
	}

	// Get image info
	img, err := b.store.Image(imageID)
	if err != nil {
		return nil, fmt.Errorf("failed to get image info: %w", err)
	}

	// Create result
	result := &api.BuildResult{
		BuildID:  req.ID,
		ImageID:  imageID,
		Digest:   img.Digest.String(),
		Tags:     req.Tags,
		Size:     img.Size,
		Duration: duration,
		Created:  time.Now(),
	}

	status.Status = api.BuildPhaseCompleted
	status.Progress = 100
	endTime := time.Now()
	status.EndTime = &endTime

	return result, nil
}

// BuildFromDockerfile builds from a Dockerfile path.
func (b *Builder) BuildFromDockerfile(ctx context.Context, dockerfilePath string, opts *api.BuildOptions) (*api.BuildResult, error) {
	if opts == nil {
		opts = &api.BuildOptions{}
	}

	req := &api.BuildRequest{
		ID:         fmt.Sprintf("build-%d", time.Now().UnixNano()),
		Dockerfile: dockerfilePath,
		ContextDir: filepath.Dir(dockerfilePath),
		Tags:       []string{"latest"},
		NoCache:    opts.NoCache,
		Pull:       opts.Pull,
		Created:    time.Now(),
		BuildArgs:  make(map[string]string),
		Labels:     make(map[string]string),
	}

	return b.BuildImage(ctx, req)
}

// GetBuildStatus returns build status.
func (b *Builder) GetBuildStatus(ctx context.Context, buildID string) (*api.BuildStatus, error) {
	b.buildsMutex.RLock()
	defer b.buildsMutex.RUnlock()
	status, exists := b.builds[buildID]
	if !exists {
		return nil, fmt.Errorf("build %s not found", buildID)
	}
	statusCopy := *status
	return &statusCopy, nil
}

// ListBuilds returns all build statuses.
func (b *Builder) ListBuilds(ctx context.Context) ([]*api.BuildStatus, error) {
	b.buildsMutex.RLock()
	defer b.buildsMutex.RUnlock()
	var statuses []*api.BuildStatus
	for _, status := range b.builds {
		statusCopy := *status
		statuses = append(statuses, &statusCopy)
	}
	return statuses, nil
}

// CleanupBuild cleans up build resources.
func (b *Builder) CleanupBuild(ctx context.Context, buildID string) error {
	b.untrackBuild(buildID)
	return nil
}

// ValidateConfig validates build configuration.
func (b *Builder) ValidateConfig(config *api.BuildConfig) error {
	if config == nil {
		return fmt.Errorf("config cannot be nil")
	}
	if config.RunRoot == "" {
		return fmt.Errorf("run_root is required")
	}
	if config.GraphRoot == "" {
		return fmt.Errorf("graph_root is required")
	}
	if config.StorageDriver == "" {
		config.StorageDriver = "overlay"
	}
	return nil
}

// Close shuts down the builder.
func (b *Builder) Close() error {
	if b.store != nil {
		b.store.Shutdown(false)
	}
	b.initialized = false
	return nil
}

func (b *Builder) executeBuild(ctx context.Context, req *api.BuildRequest, status *api.BuildStatus) (string, error) {
	status.Status = api.BuildPhaseBuilding
	status.Progress = 50
	buildOptions := &imagebuildah.BuildOptions{
		ContextDirectory: req.ContextDir,
		PullPolicy:       define.PullIfMissing,
		Args:             req.BuildArgs,
		Labels:           req.Labels,
		NoCache:          req.NoCache,
		OutputFormat:     define.OCI,
		ReportWriter:     os.Stderr,
	}
	if req.Pull {
		buildOptions.PullPolicy = define.PullAlways
	}
	if req.Target != "" {
		buildOptions.Target = req.Target
	}
	imageID, _, err := imagebuildah.BuildDockerfiles(ctx, b.store, buildOptions, req.Dockerfile)
	return imageID, err
}

func (b *Builder) validateRequest(req *api.BuildRequest) error {
	if req.ID == "" {
		return fmt.Errorf("build ID required")
	}
	if req.Dockerfile == "" {
		return fmt.Errorf("dockerfile required")
	}
	if req.ContextDir == "" {
		return fmt.Errorf("context dir required")
	}
	if len(req.Tags) == 0 {
		return fmt.Errorf("at least one tag required")
	}
	return nil
}

func (b *Builder) trackBuild(status *api.BuildStatus) {
	b.buildsMutex.Lock()
	defer b.buildsMutex.Unlock()
	b.builds[status.BuildID] = status
}

func (b *Builder) untrackBuild(buildID string) {
	b.buildsMutex.Lock()
	defer b.buildsMutex.Unlock()
	delete(b.builds, buildID)
}
