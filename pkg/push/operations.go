package push

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-logr/logr"
	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1"
	"github.com/spf13/cobra"
)

// PushOperation represents a complete push operation configuration
type PushOperation struct {
	Registry     string                // Target registry URL
	Username     string                // Registry username
	Password     string                // Registry password
	Insecure     bool                  // Allow insecure (HTTP) connections
	BuildPath    string                // Path to search for images
	UserAgent    string                // Custom user agent
	MaxRetries   int                   // Maximum retry attempts
	Concurrency  int                   // Maximum concurrent pushes
	Auth         authn.Authenticator   // Authentication configuration
	Transport    *http.Transport       // HTTP transport configuration
	Progress     ProgressReporter      // Progress reporting
	Logger       *PushLogger           // Push-specific logger
	FilterCriteria *FilterCriteria     // Image filtering criteria
}

// PushOperationResult contains the results of a push operation
type PushOperationResult struct {
	StartTime     time.Time      // When the operation started
	EndTime       time.Time      // When the operation completed
	TotalDuration time.Duration  // Total time taken
	ImagesFound   int            // Number of images discovered
	ImagesPushed  int            // Number of images successfully pushed
	ImagesFailed  int            // Number of images that failed to push
	TotalBytes    int64          // Total bytes pushed
	Results       []*PushResult  // Individual push results
	Errors        []error        // Any errors encountered
	Metrics       *PushMetrics   // Detailed performance metrics
}

// NewPushOperationFromCommand creates a PushOperation from cobra command flags
func NewPushOperationFromCommand(cmd *cobra.Command, logger logr.Logger) (*PushOperation, error) {
	// Extract command flags (these would be set by E1.2.1 Command Structure)
	registry, _ := cmd.Flags().GetString("registry")
	username, _ := cmd.Flags().GetString("username")
	password, _ := cmd.Flags().GetString("password")
	insecure, _ := cmd.Flags().GetBool("insecure")
	buildPath, _ := cmd.Flags().GetString("build-path")
	maxRetries, _ := cmd.Flags().GetInt("max-retries")
	concurrency, _ := cmd.Flags().GetInt("concurrency")

	// Use default build path if not specified
	if buildPath == "" {
		buildPath = "."
	}

	// Set sensible defaults
	if maxRetries == 0 {
		maxRetries = 3
	}
	if concurrency == 0 {
		concurrency = 3
	}

	// Create the operation
	op := &PushOperation{
		Registry:    registry,
		Username:    username,
		Password:    password,
		Insecure:    insecure,
		BuildPath:   buildPath,
		UserAgent:   "idpbuilder-push/1.0.0",
		MaxRetries:  maxRetries,
		Concurrency: concurrency,
		Logger:      NewPushLogger(logger),
	}

	// Set up authentication (this would integrate with E1.2.2 Registry Authentication)
	err := op.setupAuthentication()
	if err != nil {
		return nil, fmt.Errorf("failed to setup authentication: %w", err)
	}

	// Set up transport
	op.setupTransport()

	// Set up progress reporting
	op.Progress = NewConsoleProgressReporter(os.Stdout)

	return op, nil
}

// setupAuthentication configures registry authentication
func (op *PushOperation) setupAuthentication() error {
	if op.Username == "" && op.Password == "" {
		// Use anonymous authentication
		op.Auth = authn.Anonymous
		op.Logger.LogAuthentication(op.Registry, "anonymous")
		return nil
	}

	// Use basic authentication
	op.Auth = &authn.Basic{
		Username: op.Username,
		Password: op.Password,
	}
	op.Logger.LogAuthentication(op.Registry, "basic")

	return nil
}

// setupTransport configures the HTTP transport
func (op *PushOperation) setupTransport() {
	op.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: op.Insecure,
		},
	}
}

// Execute performs the complete push operation
func (op *PushOperation) Execute(ctx context.Context) (*PushOperationResult, error) {
	result := &PushOperationResult{
		StartTime: time.Now(),
		Results:   make([]*PushResult, 0),
		Errors:    make([]error, 0),
		Metrics:   &PushMetrics{},
	}

	op.Logger.LogInfo("Starting push operation",
		"registry", op.Registry,
		"build_path", op.BuildPath,
		"insecure", op.Insecure)

	// Step 1: Discover images
	discoveryStart := time.Now()
	op.Logger.LogDiscoveryStart(op.BuildPath)

	images, err := op.discoverImages()
	if err != nil {
		result.Errors = append(result.Errors, fmt.Errorf("image discovery failed: %w", err))
		return result, err
	}

	discoveryDuration := time.Since(discoveryStart)
	result.ImagesFound = len(images)
	result.Metrics.DiscoveryDuration = discoveryDuration

	op.Logger.LogDiscoveryComplete(len(images), discoveryDuration)

	if len(images) == 0 {
		op.Logger.LogInfo("No images found to push")
		result.EndTime = time.Now()
		result.TotalDuration = result.EndTime.Sub(result.StartTime)
		return result, nil
	}

	// Step 2: Validate images
	err = op.validateImages(images)
	if err != nil {
		result.Errors = append(result.Errors, fmt.Errorf("image validation failed: %w", err))
		return result, err
	}

	// Step 3: Push images
	pushStart := time.Now()
	pushResults, err := op.pushImages(ctx, images)
	if err != nil {
		result.Errors = append(result.Errors, fmt.Errorf("push operation failed: %w", err))
	}

	pushDuration := time.Since(pushStart)
	result.Metrics.PushDuration = pushDuration

	// Step 4: Aggregate results
	result.Results = pushResults
	for _, pushResult := range pushResults {
		if pushResult.Error == nil {
			result.ImagesPushed++
			result.TotalBytes += pushResult.Size
		} else {
			result.ImagesFailed++
			result.Errors = append(result.Errors, pushResult.Error)
		}
		result.Metrics.RetryCount += pushResult.Retries
	}

	result.EndTime = time.Now()
	result.TotalDuration = result.EndTime.Sub(result.StartTime)
	result.Metrics.TotalDuration = result.TotalDuration
	result.Metrics.ImagesPushed = result.ImagesPushed
	result.Metrics.TotalBytes = result.TotalBytes

	// Log final metrics
	op.Logger.LogPerformanceMetrics(result.Metrics)

	op.Logger.LogInfo("Push operation completed",
		"total_duration", result.TotalDuration,
		"images_found", result.ImagesFound,
		"images_pushed", result.ImagesPushed,
		"images_failed", result.ImagesFailed,
		"total_bytes", result.TotalBytes)

	return result, nil
}

// discoverImages finds all pushable images in the build path
func (op *PushOperation) discoverImages() ([]*LocalImage, error) {
	images, err := DiscoverLocalImages(op.BuildPath)
	if err != nil {
		return nil, fmt.Errorf("failed to discover images: %w", err)
	}

	// Apply filtering if criteria is set
	if op.FilterCriteria != nil {
		images = FilterPushTargets(images, op.FilterCriteria)
	}

	// Log discovered images
	for _, img := range images {
		size := int64(-1)
		if info, err := os.Stat(img.Path); err == nil {
			size = info.Size()
		}
		op.Logger.LogImageDiscovered(img.Name, img.Format, img.Path, size)
	}

	return images, nil
}

// validateImages performs pre-push validation on all discovered images
func (op *PushOperation) validateImages(images []*LocalImage) error {
	pusher := NewImagePusherWithOptions(
		op.Auth,
		op.Transport,
		op.Progress,
		op.Logger.logger,
		PusherOptions{
			MaxRetries:        op.MaxRetries,
			InitialBackoff:    time.Second,
			BackoffMultiplier: 2.0,
			MaxBackoff:        30 * time.Second,
			Insecure:          op.Insecure,
			UserAgent:         op.UserAgent,
		},
	)

	var validationIssues []string

	for _, img := range images {
		err := pusher.ValidateImage(img.Image)
		if err != nil {
			validationIssues = append(validationIssues,
				fmt.Sprintf("%s: %v", img.Name, err))
		}

		op.Logger.LogImageValidation(img.Name, err == nil,
			func() []string {
				if err != nil {
					return []string{err.Error()}
				}
				return nil
			}())
	}

	if len(validationIssues) > 0 {
		return fmt.Errorf("image validation failed: %v", validationIssues)
	}

	return nil
}

// pushImages performs the actual push operations for all images
func (op *PushOperation) pushImages(ctx context.Context, images []*LocalImage) ([]*PushResult, error) {
	pusher := NewImagePusherWithOptions(
		op.Auth,
		op.Transport,
		op.Progress,
		op.Logger.logger,
		PusherOptions{
			MaxRetries:        op.MaxRetries,
			InitialBackoff:    time.Second,
			BackoffMultiplier: 2.0,
			MaxBackoff:        30 * time.Second,
			Insecure:          op.Insecure,
			UserAgent:         op.UserAgent,
		},
	)

	// Convert images to reference map for batch push
	imageMap := make(map[name.Reference]v1.Image)
	for _, img := range images {
		var refString string

		// Check if img.Name already contains a full reference (registry/repo:tag)
		if strings.Contains(img.Name, "/") || strings.Contains(img.Name, ":") {
			// Name includes registry/repo/tag - use it directly
			refString = img.Name
		} else {
			// Name is just image name - prepend registry
			refString = fmt.Sprintf("%s/%s:latest", op.Registry, img.Name)
		}

		ref, err := name.ParseReference(refString)
		if err != nil {
			op.Logger.LogError(err, "Failed to parse reference", "ref", refString)
			continue
		}
		imageMap[ref] = img.Image
	}

	// Perform batch push with concurrency control
	results, err := pusher.BatchPush(ctx, imageMap, op.Concurrency)
	if err != nil {
		return results, fmt.Errorf("batch push failed: %w", err)
	}

	return results, nil
}

// PushImages is the main entry point for the push command (called by E1.2.1)
func PushImages(cmd *cobra.Command, args []string) error {
	// Get logger from command context or create a default one
	logger := logr.Discard() // Default no-op logger
	if loggerValue := cmd.Context().Value("logger"); loggerValue != nil {
		if cmdLogger, ok := loggerValue.(logr.Logger); ok {
			logger = cmdLogger
		}
	}

	// Create push operation from command
	operation, err := NewPushOperationFromCommand(cmd, logger)
	if err != nil {
		return fmt.Errorf("failed to create push operation: %w", err)
	}

	// Execute the push operation
	ctx := cmd.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	result, err := operation.Execute(ctx)
	if err != nil {
		return fmt.Errorf("push operation failed: %w", err)
	}

	// Check if any pushes failed
	if result.ImagesFailed > 0 {
		return fmt.Errorf("push completed with %d failures out of %d images",
			result.ImagesFailed, result.ImagesFound)
	}

	return nil
}

// Summary provides a human-readable summary of the push operation
func (result *PushOperationResult) Summary() string {
	successRate := float64(result.ImagesPushed) / float64(result.ImagesFound) * 100

	return fmt.Sprintf(`Push Operation Summary:
  Duration: %s
  Images Found: %d
  Images Pushed: %d
  Images Failed: %d
  Success Rate: %.1f%%
  Total Bytes: %s
  Average Throughput: %.2f MB/s`,
		result.TotalDuration,
		result.ImagesFound,
		result.ImagesPushed,
		result.ImagesFailed,
		successRate,
		formatBytes(result.TotalBytes),
		result.Metrics.AverageThroughputMBps())
}

// formatBytes formats bytes as human-readable string
func formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}

	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}

	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}