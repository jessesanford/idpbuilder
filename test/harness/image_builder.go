package harness

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/docker/docker/api/types"
	dockerclient "github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
)

// ImageConfig defines test image characteristics for building test containers
type ImageConfig struct {
	Name   string // Image name (without tag)
	Tag    string // Image tag
	Layers int    // Number of layers to create
	SizeMB int    // Approximate total size in megabytes
	Arch   string // Architecture (amd64, arm64, etc.)
}

// BuildResult contains information about a built image
type BuildResult struct {
	ImageRef   string // Full image reference (name:tag)
	ImageID    string // Docker image ID
	LayerCount int    // Actual number of layers
	SizeBytes  int64  // Actual size in bytes
}

// BuilderTestEnvironment represents the test environment for image builder integration testing
// This type is specific to effort 3.1.2 image builder functionality
// Renamed from TestEnvironment to avoid collision with effort 3.1.1 test harness
type BuilderTestEnvironment struct {
	DockerClient *dockerclient.Client
}

// BuildTestImage creates a Docker image for testing with specified characteristics.
//
// This function generates a temporary Dockerfile and build context, builds the image
// using the Docker daemon, and returns information about the built image.
//
// The function will:
// 1. Create temporary build directory
// 2. Generate Dockerfile based on config
// 3. Generate test files for each layer
// 4. Execute Docker build
// 5. Clean up build directory
// 6. Return build result
//
// Example:
//
//	config := ImageConfig{
//	    Name: "testapp",
//	    Tag: "v1",
//	    Layers: 3,
//	    SizeMB: 10,
//	    Arch: "amd64",
//	}
//	result, err := env.BuildTestImage(ctx, config)
func (env *BuilderTestEnvironment) BuildTestImage(ctx context.Context, config ImageConfig) (*BuildResult, error) {
	// Validate configuration
	if err := validateConfig(config); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	// Create temporary build directory
	buildDir, err := os.MkdirTemp("", fmt.Sprintf("docker-build-%s-", config.Name))
	if err != nil {
		return nil, fmt.Errorf("failed to create build directory: %w", err)
	}
	defer cleanupBuildDir(buildDir)

	// Generate Dockerfile
	dockerfileContent := generateDockerfile(config)
	dockerfilePath := filepath.Join(buildDir, "Dockerfile")
	if err := os.WriteFile(dockerfilePath, []byte(dockerfileContent), 0644); err != nil {
		return nil, fmt.Errorf("failed to write Dockerfile: %w", err)
	}

	// Generate test files for layers
	if err := generateTestFiles(buildDir, config); err != nil {
		return nil, fmt.Errorf("failed to generate test files: %w", err)
	}

	// Create tar archive for build context
	buildContext, err := archive.TarWithOptions(buildDir, &archive.TarOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to create build context: %w", err)
	}
	defer buildContext.Close()

	// Build the image
	imageRef := fmt.Sprintf("%s:%s", config.Name, config.Tag)
	imageID, err := env.executeBuild(ctx, buildContext, imageRef)
	if err != nil {
		return nil, fmt.Errorf("failed to build image: %w", err)
	}

	// Verify and get image information
	result, err := env.verifyImageBuilt(ctx, imageRef, imageID)
	if err != nil {
		return nil, fmt.Errorf("failed to verify built image: %w", err)
	}

	return result, nil
}

// validateConfig checks that the image configuration is valid
func validateConfig(config ImageConfig) error {
	if config.Name == "" {
		return fmt.Errorf("image name cannot be empty")
	}
	if config.Tag == "" {
		return fmt.Errorf("image tag cannot be empty")
	}
	if config.Layers <= 0 {
		return fmt.Errorf("layer count must be positive, got %d", config.Layers)
	}
	if config.SizeMB < 0 {
		return fmt.Errorf("size cannot be negative, got %d", config.SizeMB)
	}
	if config.Arch == "" {
		config.Arch = "amd64" // Default architecture
	}
	return nil
}

// generateDockerfile creates Dockerfile content with specified layers.
//
// The generated Dockerfile will create the requested number of layers
// by using RUN and COPY commands appropriately.
func generateDockerfile(config ImageConfig) string {
	var dockerfile strings.Builder

	// Start with Alpine base
	dockerfile.WriteString("FROM alpine:latest\n\n")

	// Add RUN/COPY commands to create layers
	for i := 0; i < config.Layers; i++ {
		// Copy test file
		dockerfile.WriteString(fmt.Sprintf("# Layer %d\n", i+1))
		dockerfile.WriteString(fmt.Sprintf("COPY layer-%d.dat /data/\n", i))
		dockerfile.WriteString(fmt.Sprintf("RUN echo 'Layer %d complete' > /data/layer-%d.txt\n\n", i, i))
	}

	// Add test labels
	dockerfile.WriteString(fmt.Sprintf("LABEL test.image.name=\"%s\"\n", config.Name))
	dockerfile.WriteString(fmt.Sprintf("LABEL test.image.layers=\"%d\"\n", config.Layers))
	dockerfile.WriteString(fmt.Sprintf("LABEL test.image.arch=\"%s\"\n\n", config.Arch))

	// Add CMD instruction
	dockerfile.WriteString("CMD [\"sh\", \"-c\", \"echo 'Test image running'\"]\n")

	return dockerfile.String()
}

// generateTestFiles creates files for each layer to achieve the approximate
// total size specified in the configuration.
func generateTestFiles(buildDir string, config ImageConfig) error {
	if config.SizeMB == 0 {
		// Create minimal files if no size specified
		for i := 0; i < config.Layers; i++ {
			filename := filepath.Join(buildDir, fmt.Sprintf("layer-%d.dat", i))
			if err := os.WriteFile(filename, []byte("test data"), 0644); err != nil {
				return fmt.Errorf("failed to create test file %s: %w", filename, err)
			}
		}
		return nil
	}

	// Calculate size per layer in bytes
	totalBytes := int64(config.SizeMB) * 1024 * 1024
	bytesPerLayer := totalBytes / int64(config.Layers)

	// Create appropriately sized files
	for i := 0; i < config.Layers; i++ {
		filename := filepath.Join(buildDir, fmt.Sprintf("layer-%d.dat", i))
		if err := createRandomFile(filename, bytesPerLayer); err != nil {
			return fmt.Errorf("failed to create test file %s: %w", filename, err)
		}
	}

	return nil
}

// createRandomFile creates a file with random data of the specified size
func createRandomFile(filename string, size int64) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Use /dev/urandom for random data
	random, err := os.Open("/dev/urandom")
	if err != nil {
		return fmt.Errorf("failed to open /dev/urandom: %w", err)
	}
	defer random.Close()

	// Copy random data to file
	written, err := io.CopyN(file, random, size)
	if err != nil {
		return fmt.Errorf("failed to write random data: %w", err)
	}

	if written != size {
		return fmt.Errorf("expected to write %d bytes, wrote %d", size, written)
	}

	return nil
}

// cleanupBuildDir removes the temporary build directory
func cleanupBuildDir(buildDir string) error {
	return os.RemoveAll(buildDir)
}

// executeBuild runs the Docker build operation
func (env *BuilderTestEnvironment) executeBuild(ctx context.Context, buildContext io.Reader, imageName string) (string, error) {
	// Configure build options
	buildOptions := types.ImageBuildOptions{
		Tags:           []string{imageName},
		Remove:         true, // Remove intermediate containers
		ForceRemove:    true,
		SuppressOutput: false,
		NoCache:        false,
		PullParent:     false,
	}

	// Call Docker API ImageBuild
	response, err := env.DockerClient.ImageBuild(ctx, buildContext, buildOptions)
	if err != nil {
		return "", fmt.Errorf("Docker build failed: %w", err)
	}
	defer response.Body.Close()

	// Stream and check build output
	imageID, err := streamBuildOutput(response.Body)
	if err != nil {
		return "", fmt.Errorf("build output error: %w", err)
	}

	return imageID, nil
}

// verifyImageBuilt checks that the image exists in Docker daemon and returns build result
func (env *BuilderTestEnvironment) verifyImageBuilt(ctx context.Context, imageRef string, imageID string) (*BuildResult, error) {
	// Inspect the image
	inspect, _, err := env.DockerClient.ImageInspectWithRaw(ctx, imageID)
	if err != nil {
		return nil, fmt.Errorf("failed to inspect image %s: %w", imageRef, err)
	}

	// Count layers (RootFS layers)
	layerCount := len(inspect.RootFS.Layers)

	// Get size
	sizeBytes := inspect.Size

	// Create BuildResult
	result := &BuildResult{
		ImageRef:   imageRef,
		ImageID:    imageID,
		LayerCount: layerCount,
		SizeBytes:  sizeBytes,
	}

	return result, nil
}

// streamBuildOutput reads and logs Docker build output, extracting the image ID
func streamBuildOutput(reader io.ReadCloser) (string, error) {
	decoder := json.NewDecoder(reader)
	var imageID string

	for {
		var message struct {
			Stream string `json:"stream"`
			Error  string `json:"error"`
			Aux    struct {
				ID string `json:"ID"`
			} `json:"aux"`
		}

		if err := decoder.Decode(&message); err != nil {
			if err == io.EOF {
				break
			}
			return "", fmt.Errorf("failed to decode build output: %w", err)
		}

		// Check for errors
		if message.Error != "" {
			return "", fmt.Errorf("build error: %s", message.Error)
		}

		// Log build output
		if message.Stream != "" {
			// In production, you'd log this. For tests, we can skip or print
			// fmt.Print(message.Stream)
		}

		// Capture image ID
		if message.Aux.ID != "" {
			imageID = message.Aux.ID
		}
	}

	if imageID == "" {
		return "", fmt.Errorf("no image ID found in build output")
	}

	return imageID, nil
}
