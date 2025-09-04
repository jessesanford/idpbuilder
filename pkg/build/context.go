package build

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cnoe-io/idpbuilder/pkg/builder"
)

// BuildContext represents the context for a build operation
type BuildContext struct {
	ContextPath    string
	DockerfilePath string
	BuildOptions   *builder.BuildOptions
	Environment    map[string]string
	Args           map[string]string
	Labels         map[string]string
	Layers         []LayerInfo
	TempDir        string
	Result         *builder.BuildResult
}

// ContextManager manages build contexts
type ContextManager struct {
	workingDir string
	logger     Logger
}

// ContextInfo provides information about a build context
type ContextInfo struct {
	Path            string
	DockerfilePath  string
	Size            int64
	FileCount       int
	HasDockerfile   bool
	HasDockerignore bool
	MainFiles       []string
}

// Logger interface for context operations
type Logger interface {
	Debug(msg string, args ...interface{})
	Info(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
	Error(msg string, args ...interface{})
}

// NewBuildContext creates a new build context
func NewBuildContext(contextPath string, dockerfilePath string, opts *builder.BuildOptions) (*BuildContext, error) {
	if !isCliToolsEnabled() {
		return nil, fmt.Errorf("ENABLE_CLI_TOOLS feature flag is not enabled")
	}

	// Resolve absolute paths
	absContextPath, err := filepath.Abs(contextPath)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve context path: %w", err)
	}

	absDockerfilePath := dockerfilePath
	if !filepath.IsAbs(dockerfilePath) {
		absDockerfilePath = filepath.Join(absContextPath, dockerfilePath)
	}

	// Validate paths exist
	if _, err := os.Stat(absContextPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("build context does not exist: %s", absContextPath)
	}

	if _, err := os.Stat(absDockerfilePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("dockerfile does not exist: %s", absDockerfilePath)
	}

	return &BuildContext{
		ContextPath:    absContextPath,
		DockerfilePath: absDockerfilePath,
		BuildOptions:   opts,
		Environment:    make(map[string]string),
		Args:           make(map[string]string),
		Labels:         make(map[string]string),
		Layers:         []LayerInfo{},
	}, nil
}

// NewContextManager creates a new context manager
func NewContextManager(workingDir string, logger Logger) *ContextManager {
	return &ContextManager{
		workingDir: workingDir,
		logger:     logger,
	}
}

// AnalyzeContext analyzes a build context and returns information about it
func (cm *ContextManager) AnalyzeContext(ctx context.Context, contextPath string) (*ContextInfo, error) {
	if !isCliToolsEnabled() {
		return nil, fmt.Errorf("ENABLE_CLI_TOOLS feature flag is not enabled")
	}

	cm.logInfo("Analyzing build context: %s", contextPath)

	info := &ContextInfo{
		Path: contextPath,
	}

	// Check if context exists
	contextStat, err := os.Stat(contextPath)
	if err != nil {
		return nil, fmt.Errorf("failed to stat context path: %w", err)
	}

	if !contextStat.IsDir() {
		return nil, fmt.Errorf("context path is not a directory: %s", contextPath)
	}

	// Look for Dockerfile
	dockerfilePaths := []string{"Dockerfile", "dockerfile", "Dockerfile.dev", "Dockerfile.prod"}
	for _, dockerfilePath := range dockerfilePaths {
		fullPath := filepath.Join(contextPath, dockerfilePath)
		if _, err := os.Stat(fullPath); err == nil {
			info.DockerfilePath = dockerfilePath
			info.HasDockerfile = true
			break
		}
	}

	// Check for .dockerignore
	dockerignorePath := filepath.Join(contextPath, ".dockerignore")
	if _, err := os.Stat(dockerignorePath); err == nil {
		info.HasDockerignore = true
	}

	// Calculate context size and file count
	var totalSize int64
	var fileCount int
	var mainFiles []string

	err = filepath.Walk(contextPath, func(path string, fileInfo os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !fileInfo.IsDir() {
			totalSize += fileInfo.Size()
			fileCount++

			// Record important files
			relPath, _ := filepath.Rel(contextPath, path)
			if cm.isImportantFile(relPath) {
				mainFiles = append(mainFiles, relPath)
			}
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to walk context directory: %w", err)
	}

	info.Size = totalSize
	info.FileCount = fileCount
	info.MainFiles = mainFiles

	cm.logInfo("Context analysis complete: %d files, %d bytes", fileCount, totalSize)
	return info, nil
}

// PrepareContext prepares a build context for building
func (cm *ContextManager) PrepareContext(ctx context.Context, buildCtx *BuildContext) error {
	if !isCliToolsEnabled() {
		return fmt.Errorf("ENABLE_CLI_TOOLS feature flag is not enabled")
	}

	cm.logInfo("Preparing build context: %s", buildCtx.ContextPath)

	// Read and process .dockerignore if it exists
	if err := cm.processDockerignore(buildCtx); err != nil {
		return fmt.Errorf("failed to process .dockerignore: %w", err)
	}

	// Set up environment variables
	if err := cm.setupEnvironment(buildCtx); err != nil {
		return fmt.Errorf("failed to setup environment: %w", err)
	}

	// Validate Dockerfile
	if err := cm.validateDockerfile(buildCtx); err != nil {
		return fmt.Errorf("dockerfile validation failed: %w", err)
	}

	cm.logInfo("Build context prepared successfully")
	return nil
}

// CleanupContext cleans up resources associated with a build context
func (cm *ContextManager) CleanupContext(ctx context.Context, buildCtx *BuildContext) error {
	if buildCtx.TempDir != "" {
		cm.logInfo("Cleaning up temporary directory: %s", buildCtx.TempDir)
		if err := os.RemoveAll(buildCtx.TempDir); err != nil {
			return fmt.Errorf("failed to remove temp directory: %w", err)
		}
	}

	return nil
}

// processDockerignore reads and processes .dockerignore file
func (cm *ContextManager) processDockerignore(buildCtx *BuildContext) error {
	dockerignorePath := filepath.Join(buildCtx.ContextPath, ".dockerignore")
	
	if _, err := os.Stat(dockerignorePath); os.IsNotExist(err) {
		// No .dockerignore file, nothing to process
		return nil
	}

	// Read .dockerignore file
	content, err := os.ReadFile(dockerignorePath)
	if err != nil {
		return fmt.Errorf("failed to read .dockerignore: %w", err)
	}

	// Parse ignore patterns
	lines := strings.Split(string(content), "\n")
	patterns := []string{}
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		patterns = append(patterns, line)
	}

	cm.logInfo("Loaded %d ignore patterns from .dockerignore", len(patterns))
	
	// Store patterns in environment for later use
	buildCtx.Environment["DOCKERIGNORE_PATTERNS"] = strings.Join(patterns, ",")
	
	return nil
}

// setupEnvironment sets up environment variables for the build
func (cm *ContextManager) setupEnvironment(buildCtx *BuildContext) error {
	// Copy system environment variables
	for _, env := range os.Environ() {
		parts := strings.SplitN(env, "=", 2)
		if len(parts) == 2 {
			buildCtx.Environment[parts[0]] = parts[1]
		}
	}

	// Add build-specific variables
	buildCtx.Environment["BUILD_CONTEXT"] = buildCtx.ContextPath
	buildCtx.Environment["DOCKERFILE"] = buildCtx.DockerfilePath

	return nil
}

// validateDockerfile performs basic validation on the Dockerfile
func (cm *ContextManager) validateDockerfile(buildCtx *BuildContext) error {
	content, err := os.ReadFile(buildCtx.DockerfilePath)
	if err != nil {
		return fmt.Errorf("failed to read dockerfile: %w", err)
	}

	lines := strings.Split(string(content), "\n")
	hasFromInstruction := false

	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Check for FROM instruction
		if strings.HasPrefix(strings.ToUpper(line), "FROM ") {
			hasFromInstruction = true
		}

		// Basic syntax validation could go here
		if strings.Contains(line, "\t") {
			cm.logWarn("Line %d contains tabs, consider using spaces: %s", i+1, line)
		}
	}

	if !hasFromInstruction {
		return fmt.Errorf("dockerfile must contain at least one FROM instruction")
	}

	cm.logInfo("Dockerfile validation passed")
	return nil
}

// isImportantFile determines if a file is considered important for reporting
func (cm *ContextManager) isImportantFile(relPath string) bool {
	importantFiles := []string{
		"Dockerfile",
		"dockerfile",
		".dockerignore",
		"requirements.txt",
		"package.json",
		"go.mod",
		"pom.xml",
		"Gemfile",
		"Cargo.toml",
	}

	for _, important := range importantFiles {
		if strings.Contains(relPath, important) {
			return true
		}
	}

	// Check for common config files
	if strings.HasSuffix(relPath, ".yml") || strings.HasSuffix(relPath, ".yaml") ||
		strings.HasSuffix(relPath, ".json") || strings.HasSuffix(relPath, ".toml") {
		return true
	}

	return false
}

// Logging helper methods
func (cm *ContextManager) logInfo(msg string, args ...interface{}) {
	if cm.logger != nil {
		cm.logger.Info(msg, args...)
	}
}

func (cm *ContextManager) logWarn(msg string, args ...interface{}) {
	if cm.logger != nil {
		cm.logger.Warn(msg, args...)
	}
}