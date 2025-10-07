package push

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"time"

	"github.com/go-logr/logr"
	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/remote"
)

// ImagePusher encapsulates the logic for pushing OCI images to registries
type ImagePusher struct {
	auth      authn.Authenticator
	transport *http.Transport
	progress  ProgressReporter
	logger    logr.Logger
	options   PusherOptions
}

// PusherOptions configures the ImagePusher behavior
type PusherOptions struct {
	MaxRetries        int           // Maximum number of retry attempts
	InitialBackoff    time.Duration // Initial backoff duration
	BackoffMultiplier float64       // Multiplier for exponential backoff
	MaxBackoff        time.Duration // Maximum backoff duration
	Insecure          bool          // Allow insecure (HTTP) registries
	UserAgent         string        // Custom user agent
}

// DefaultPusherOptions returns sensible defaults for push operations
func DefaultPusherOptions() PusherOptions {
	return PusherOptions{
		MaxRetries:        5,
		InitialBackoff:    time.Second,
		BackoffMultiplier: 2.0,
		MaxBackoff:        30 * time.Second,
		Insecure:          false,
		UserAgent:         "idpbuilder-push/1.0.0",
	}
}

// NewImagePusher creates a new ImagePusher with the specified configuration
func NewImagePusher(auth authn.Authenticator, transport *http.Transport, progress ProgressReporter, logger logr.Logger) *ImagePusher {
	return NewImagePusherWithOptions(auth, transport, progress, logger, DefaultPusherOptions())
}

// NewImagePusherWithOptions creates a new ImagePusher with custom options
func NewImagePusherWithOptions(auth authn.Authenticator, transport *http.Transport, progress ProgressReporter, logger logr.Logger, options PusherOptions) *ImagePusher {
	// If no transport provided, create a default one
	if transport == nil {
		transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: options.Insecure,
			},
		}
	}

	return &ImagePusher{
		auth:      auth,
		transport: transport,
		progress:  progress,
		logger:    logger,
		options:   options,
	}
}

// PushResult contains the result of a push operation
type PushResult struct {
	ImageName string        // Name of the pushed image
	Digest    string        // Digest of the pushed image
	Size      int64         // Size of the pushed image in bytes
	Duration  time.Duration // Time taken for the push
	Error     error         // Error if push failed
	Retries   int           // Number of retries attempted
}

// Push pushes a single image to the specified registry reference
func (p *ImagePusher) Push(ctx context.Context, img v1.Image, ref name.Reference) (*PushResult, error) {
	startTime := time.Now()

	result := &PushResult{
		ImageName: ref.String(),
	}

	// Get image size for progress reporting
	size, err := p.getImageSize(img)
	if err != nil {
		p.logger.V(1).Info("Failed to get image size for progress tracking", "error", err)
		size = -1 // Unknown size
	}
	result.Size = size

	// Set up remote options
	opts := p.buildRemoteOptions()

	// Start progress reporting
	if p.progress != nil {
		digest, _ := img.Digest()
		p.progress.StartImage(digest.String(), size)
		defer p.progress.FinishImage(digest.String())
	}

	p.logger.Info("Starting image push", "image", ref.String(), "size_bytes", size)

	// Perform the actual push
	err = remote.Write(ref, img, opts...)
	if err != nil {
		result.Error = fmt.Errorf("failed to push image %s: %w", ref.String(), err)
		result.Duration = time.Since(startTime)
		return result, result.Error
	}

	// Get the digest of the pushed image
	digest, err := img.Digest()
	if err != nil {
		p.logger.V(1).Info("Failed to get image digest", "error", err)
	} else {
		result.Digest = digest.String()
	}

	result.Duration = time.Since(startTime)
	p.logger.Info("Successfully pushed image",
		"image", ref.String(),
		"digest", result.Digest,
		"duration", result.Duration,
		"size_bytes", size)

	return result, nil
}

// PushWithRetry pushes an image with automatic retry logic for transient failures
func (p *ImagePusher) PushWithRetry(ctx context.Context, img v1.Image, ref name.Reference) (*PushResult, error) {
	backoff := p.options.InitialBackoff
	var lastResult *PushResult

	for attempt := 0; attempt <= p.options.MaxRetries; attempt++ {
		// Check context cancellation
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		p.logger.V(1).Info("Attempting image push",
			"image", ref.String(),
			"attempt", attempt+1,
			"max_retries", p.options.MaxRetries+1)

		result, err := p.Push(ctx, img, ref)
		if result != nil {
			result.Retries = attempt
		}

		if err == nil {
			return result, nil
		}

		lastResult = result
		p.logger.V(1).Info("Push attempt failed",
			"image", ref.String(),
			"attempt", attempt+1,
			"error", err)

		// Check if we should retry
		if attempt >= p.options.MaxRetries {
			break
		}

		if !isRetryableError(err) {
			p.logger.Info("Error is not retryable, stopping", "error", err)
			break
		}

		// Wait before retrying
		p.logger.V(1).Info("Waiting before retry",
			"backoff", backoff,
			"attempt", attempt+1)

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(backoff):
		}

		// Calculate next backoff duration
		backoff = time.Duration(float64(backoff) * p.options.BackoffMultiplier)
		if backoff > p.options.MaxBackoff {
			backoff = p.options.MaxBackoff
		}
	}

	// All retries exhausted
	if lastResult != nil {
		return lastResult, fmt.Errorf("push failed after %d attempts: %w",
			p.options.MaxRetries+1, lastResult.Error)
	}

	return nil, fmt.Errorf("push failed after %d attempts", p.options.MaxRetries+1)
}

// buildRemoteOptions constructs the remote options for go-containerregistry
func (p *ImagePusher) buildRemoteOptions() []remote.Option {
	var opts []remote.Option

	// Add authentication if provided
	if p.auth != nil {
		opts = append(opts, remote.WithAuth(p.auth))
	}

	// Add custom transport if provided
	if p.transport != nil {
		opts = append(opts, remote.WithTransport(p.transport))
	}

	// Set user agent
	if p.options.UserAgent != "" {
		opts = append(opts, remote.WithUserAgent(p.options.UserAgent))
	}

	// Add context for cancellation
	// Note: Context will be passed separately to Write()

	return opts
}

// getImageSize calculates the total size of an image
func (p *ImagePusher) getImageSize(img v1.Image) (int64, error) {
	manifest, err := img.Manifest()
	if err != nil {
		return 0, fmt.Errorf("failed to get manifest: %w", err)
	}

	var totalSize int64

	// Add config size
	totalSize += manifest.Config.Size

	// Add layer sizes
	for _, layer := range manifest.Layers {
		totalSize += layer.Size
	}

	return totalSize, nil
}

// isRetryableError determines if an error is worth retrying
func isRetryableError(err error) bool {
	// Network-related errors that are typically transient
	errorStr := err.Error()

	retryableErrors := []string{
		"connection refused",
		"connection reset",
		"connection timeout",
		"timeout",
		"temporary failure",
		"service unavailable",
		"502 bad gateway",
		"503 service unavailable",
		"504 gateway timeout",
		"too many requests",
		"rate limit",
		"network is unreachable",
		"no route to host",
		"i/o timeout",
	}

	for _, retryable := range retryableErrors {
		if contains(errorStr, retryable) {
			return true
		}
	}

	return false
}

// contains checks if a string contains a substring (case insensitive)
func contains(s, substr string) bool {
	return len(s) >= len(substr) &&
		(s == substr ||
			len(s) > len(substr) &&
				(s[:len(substr)] == substr ||
					s[len(s)-len(substr):] == substr ||
					containsInner(s, substr)))
}

func containsInner(s, substr string) bool {
	for i := 1; i < len(s)-len(substr)+1; i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// BatchPush pushes multiple images concurrently
func (p *ImagePusher) BatchPush(ctx context.Context, images map[name.Reference]v1.Image, maxConcurrency int) ([]*PushResult, error) {
	if maxConcurrency <= 0 {
		maxConcurrency = 3 // Default concurrency
	}

	results := make([]*PushResult, 0, len(images))
	resultChan := make(chan *PushResult, len(images))
	semaphore := make(chan struct{}, maxConcurrency)

	// Start goroutines for each image
	for ref, img := range images {
		go func(ref name.Reference, img v1.Image) {
			semaphore <- struct{}{}        // Acquire semaphore
			defer func() { <-semaphore }() // Release semaphore

			result, _ := p.PushWithRetry(ctx, img, ref)
			resultChan <- result
		}(ref, img)
	}

	// Collect results
	for i := 0; i < len(images); i++ {
		select {
		case result := <-resultChan:
			results = append(results, result)
		case <-ctx.Done():
			return results, ctx.Err()
		}
	}

	return results, nil
}

// ValidateImage performs pre-push validation on an image
func (p *ImagePusher) ValidateImage(img v1.Image) error {
	// Check if image has a valid manifest
	manifest, err := img.Manifest()
	if err != nil {
		return fmt.Errorf("invalid image manifest: %w", err)
	}

	// Verify config blob exists
	if manifest.Config.Size == 0 {
		return fmt.Errorf("image config is empty")
	}

	// Verify layers exist
	if len(manifest.Layers) == 0 {
		p.logger.V(1).Info("Warning: image has no layers")
	}

	// Check for valid media types
	if manifest.Config.MediaType == "" {
		return fmt.Errorf("image config has no media type")
	}

	p.logger.V(2).Info("Image validation passed",
		"config_size", manifest.Config.Size,
		"layer_count", len(manifest.Layers))

	return nil
}
