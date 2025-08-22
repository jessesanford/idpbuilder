package buildah

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// validateBuildOptions validates the build options for correctness
func (c *client) validateBuildOptions(opts BuildOptions) error {
	if opts.ContextDir == "" {
		return fmt.Errorf("context directory is required")
	}

	if opts.Tag == "" {
		return fmt.Errorf("image tag is required")
	}

	// Default dockerfile if not specified
	if opts.Dockerfile == "" {
		opts.Dockerfile = "Dockerfile"
	}

	// Check if dockerfile exists in context directory
	dockerfilePath := filepath.Join(opts.ContextDir, opts.Dockerfile)
	if !fileExists(dockerfilePath) {
		return fmt.Errorf("dockerfile not found at %s", dockerfilePath)
	}

	// Validate context directory exists
	if !dirExists(opts.ContextDir) {
		return fmt.Errorf("context directory does not exist: %s", opts.ContextDir)
	}

	return nil
}

// buildBuildArgs constructs the command line arguments for buildah build
func (c *client) buildBuildArgs(opts BuildOptions) ([]string, error) {
	args := []string{"build"}

	// Add log level
	args = append(args, "--log-level", c.logLevel)

	// Add dockerfile if specified and not default
	if opts.Dockerfile != "" && opts.Dockerfile != "Dockerfile" {
		args = append(args, "--file", opts.Dockerfile)
	}

	// Add tag
	args = append(args, "--tag", opts.Tag)

	// Add build args
	for key, value := range opts.BuildArgs {
		args = append(args, "--build-arg", key+"="+value)
	}

	// Add labels
	for key, value := range opts.Labels {
		args = append(args, "--label", key+"="+value)
	}

	// Add cache control
	if opts.NoCache {
		args = append(args, "--no-cache")
	}

	// Add pull policy
	if opts.Pull {
		args = append(args, "--pull")
	}

	// Add context directory as the last argument
	args = append(args, opts.ContextDir)

	return args, nil
}

// runBuildahCommand executes a buildah command with the given arguments
func (c *client) runBuildahCommand(ctx context.Context, args []string, opts BuildOptions) ([]byte, error) {
	// Create context with timeout
	timeoutCtx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	// Build the command
	cmd := exec.CommandContext(timeoutCtx, c.buildahPath, args...)
	cmd.Dir = c.workDir

	// Set up output capture
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// If progress writer is provided, also write to it
	if opts.Progress != nil {
		cmd.Stdout = io.MultiWriter(&stdout, opts.Progress)
		cmd.Stderr = io.MultiWriter(&stderr, opts.Progress)
	}

	c.logger.V(1).Info("Executing buildah command", 
		"args", args,
		"dir", cmd.Dir,
		"timeout", c.timeout)

	// Execute the command
	err := cmd.Run()
	
	output := stdout.Bytes()
	errorOutput := stderr.Bytes()
	
	if err != nil {
		// Get exit code if available
		exitCode := -1
		if exitError, ok := err.(*exec.ExitError); ok {
			exitCode = exitError.ExitCode()
		}

		// Combine stdout and stderr for error output
		combinedOutput := string(output)
		if len(errorOutput) > 0 {
			if len(combinedOutput) > 0 {
				combinedOutput += "\n"
			}
			combinedOutput += string(errorOutput)
		}

		return nil, &BuildError{
			Op:       "buildah_command",
			Err:      err,
			Output:   combinedOutput,
			ExitCode: exitCode,
		}
	}

	c.logger.V(1).Info("Buildah command completed successfully")
	return output, nil
}

// parseBuildResult extracts build information from buildah output
func (c *client) parseBuildResult(opts BuildOptions, output []byte, duration time.Duration) (*BuildResult, error) {
	result := &BuildResult{
		Tag:       opts.Tag,
		Duration:  duration,
		BuildArgs: make(map[string]string),
		Labels:    make(map[string]string),
	}

	// Copy build args
	for k, v := range opts.BuildArgs {
		result.BuildArgs[k] = v
	}

	// Copy labels
	for k, v := range opts.Labels {
		result.Labels[k] = v
	}

	// Parse output to extract image ID
	outputStr := string(output)
	lines := strings.Split(outputStr, "\n")
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		
		// Look for image ID in various formats buildah might output
		if strings.HasPrefix(line, "COMMIT") || strings.Contains(line, "Writing image") {
			// Try to extract image ID from commit line
			if imageID := extractImageID(line); imageID != "" {
				result.ImageID = imageID
				break
			}
		}
		
		// Look for size information
		if strings.Contains(line, "size") && strings.Contains(line, "bytes") {
			if size := extractImageSize(line); size > 0 {
				result.Size = size
			}
		}
	}

	// If we couldn't parse an image ID from output, try to get it another way
	if result.ImageID == "" {
		c.logger.V(1).Info("Could not parse image ID from build output, using tag as fallback")
		result.ImageID = opts.Tag
	}

	return result, nil
}

// extractImageID attempts to extract an image ID from a buildah output line
func extractImageID(line string) string {
	// Common patterns for image IDs in buildah output
	patterns := []string{
		"COMMIT ",
		"Writing image ",
		"sha256:",
	}

	for _, pattern := range patterns {
		if idx := strings.Index(line, pattern); idx != -1 {
			// Extract the part after the pattern
			remaining := strings.TrimSpace(line[idx+len(pattern):])
			
			// Take the first word/token which should be the image ID
			if parts := strings.Fields(remaining); len(parts) > 0 {
				imageID := parts[0]
				
				// Clean up common suffixes/prefixes
				imageID = strings.TrimPrefix(imageID, "sha256:")
				if len(imageID) >= 12 {
					return imageID
				}
			}
		}
	}

	return ""
}

// extractImageSize attempts to extract image size from a buildah output line
func extractImageSize(line string) int64 {
	// Look for size information in the line
	words := strings.Fields(strings.ToLower(line))
	
	for i, word := range words {
		if word == "size" && i > 0 {
			// The size should be the previous word
			if sizeStr := words[i-1]; sizeStr != "" {
				if size, err := strconv.ParseInt(sizeStr, 10, 64); err == nil {
					return size
				}
			}
		}
		
		if strings.Contains(word, "bytes") && i > 0 {
			// Try to extract number from previous words
			for j := i - 1; j >= 0; j-- {
				if size, err := strconv.ParseInt(words[j], 10, 64); err == nil {
					return size
				}
			}
		}
	}

	return 0
}

// fileExists checks if a file exists at the given path
func fileExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// dirExists checks if a directory exists at the given path
func dirExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}