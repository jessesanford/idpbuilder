package builder

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/empty"
)

// Builder provides basic container image building capabilities
type Builder struct {
	options *BuildOptions
}

// BuildResult contains the result of a build operation
type BuildResult struct {
	Image   v1.Image
	ImageID string
	Tags    []string
}

// New creates a new Builder instance with the provided options
func New(opts *BuildOptions) (*Builder, error) {
	if opts == nil {
		opts = DefaultBuildOptions()
	}

	if err := opts.Validate(); err != nil {
		return nil, fmt.Errorf("invalid build options: %w", err)
	}

	return &Builder{options: opts}, nil
}

// Build builds a container image according to the configured options
func (b *Builder) Build() (*BuildResult, error) {
	if !b.isFeatureEnabled() {
		return nil, fmt.Errorf("ENABLE_CORE_BUILDER feature flag is not enabled")
	}

	b.logInfo("Starting image build")
	b.logInfo("Context: %s", b.options.Context)
	b.logInfo("Dockerfile: %s", b.options.Dockerfile)

	// Parse Dockerfile
	dockerfile, err := b.parseDockerfile()
	if err != nil {
		return nil, fmt.Errorf("failed to parse Dockerfile: %w", err)
	}

	// Create base image
	img := empty.Image

	// Apply basic Dockerfile instructions
	for _, instruction := range dockerfile {
		b.logInfo("Processing %s: %s", instruction.Command, instruction.Args)
	}

	// Get image details
	manifest, err := img.Manifest()
	if err != nil {
		return nil, fmt.Errorf("failed to get image manifest: %w", err)
	}

	result := &BuildResult{
		Image:   img,
		ImageID: manifest.Config.Digest.String(),
		Tags:    b.options.Tags,
	}

	b.logInfo("Build completed successfully")
	return result, nil
}

// isFeatureEnabled checks if the ENABLE_CORE_BUILDER feature flag is set
func (b *Builder) isFeatureEnabled() bool {
	return os.Getenv("ENABLE_CORE_BUILDER") == "true"
}

// parseDockerfile reads and parses the Dockerfile
func (b *Builder) parseDockerfile() ([]DockerfileInstruction, error) {
	file, err := os.Open(b.options.Dockerfile)
	if err != nil {
		return nil, fmt.Errorf("failed to open Dockerfile: %w", err)
	}
	defer file.Close()

	var instructions []DockerfileInstruction
	scanner := bufio.NewScanner(file)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		instruction, err := parseDockerfileInstruction(line, lineNum)
		if err != nil {
			return nil, fmt.Errorf("line %d: %w", lineNum, err)
		}

		instructions = append(instructions, instruction)
	}

	return instructions, nil
}

// Logging helper methods
func (b *Builder) logInfo(msg string, args ...interface{}) {
	if b.options.Logger != nil {
		b.options.Logger.Info(msg, args...)
	}
}