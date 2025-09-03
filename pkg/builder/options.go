package builder

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// BuildOptions contains all configuration for building container images
type BuildOptions struct {
	// Core image configuration
	ImageName   string            `json:"image_name" yaml:"image_name"`
	Tag         string            `json:"tag" yaml:"tag"`
	BaseImage   string            `json:"base_image" yaml:"base_image"`
	Platform    string            `json:"platform" yaml:"platform"`
	
	// Build context
	ContextPath string     `json:"context_path" yaml:"context_path"`
	Files       []FileSpec `json:"files" yaml:"files"`
	
	// Image configuration
	Labels     map[string]string `json:"labels" yaml:"labels"`
	Env        []string          `json:"env" yaml:"env"`
	Cmd        []string          `json:"cmd" yaml:"cmd"`
	Entrypoint []string          `json:"entrypoint" yaml:"entrypoint"`
	WorkingDir string            `json:"working_dir" yaml:"working_dir"`
	User       string            `json:"user" yaml:"user"`
	
	// Build configuration
	NoCache       bool              `json:"no_cache" yaml:"no_cache"`
	PullPolicy    PullPolicy        `json:"pull_policy" yaml:"pull_policy"`
	NetworkMode   string            `json:"network_mode" yaml:"network_mode"`
	BuildArgs     map[string]string `json:"build_args" yaml:"build_args"`
	Target        string            `json:"target" yaml:"target"`
	
	// Registry configuration
	Registry       string              `json:"registry" yaml:"registry"`
	RegistryAuth   *RegistryAuth       `json:"registry_auth" yaml:"registry_auth"`
	Insecure       bool                `json:"insecure" yaml:"insecure"`
	SkipTLS        bool                `json:"skip_tls" yaml:"skip_tls"`
	
	// Output configuration
	OutputFormat   OutputFormat        `json:"output_format" yaml:"output_format"`
	Compression    CompressionType     `json:"compression" yaml:"compression"`
	OutputFile     string              `json:"output_file" yaml:"output_file"`
	
	// Advanced options
	Squash         bool                `json:"squash" yaml:"squash"`
	Annotations    map[string]string   `json:"annotations" yaml:"annotations"`
	CreatedBy      string              `json:"created_by" yaml:"created_by"`
	Author         string              `json:"author" yaml:"author"`
	
	// Performance options
	MaxRetries     int                 `json:"max_retries" yaml:"max_retries"`
	Timeout        time.Duration       `json:"timeout" yaml:"timeout"`
	Parallelism    int                 `json:"parallelism" yaml:"parallelism"`
	
	// Cache configuration
	CacheFrom      []string            `json:"cache_from" yaml:"cache_from"`
	CacheTo        []string            `json:"cache_to" yaml:"cache_to"`
	CacheDir       string              `json:"cache_dir" yaml:"cache_dir"`
}

// FileSpec specifies a file to include in the container image
type FileSpec struct {
	Source      string      `json:"source" yaml:"source"`
	Destination string      `json:"destination" yaml:"destination"`
	Mode        os.FileMode `json:"mode" yaml:"mode"`
	Owner       string      `json:"owner" yaml:"owner"`
	Group       string      `json:"group" yaml:"group"`
}

// BuilderOptions contains configuration for the builder itself
type BuilderOptions struct {
	Registry    string        `json:"registry" yaml:"registry"`
	Insecure    bool          `json:"insecure" yaml:"insecure"`
	Timeout     time.Duration `json:"timeout" yaml:"timeout"`
	MaxRetries  int           `json:"max_retries" yaml:"max_retries"`
	LogLevel    string        `json:"log_level" yaml:"log_level"`
	CacheDir    string        `json:"cache_dir" yaml:"cache_dir"`
	TempDir     string        `json:"temp_dir" yaml:"temp_dir"`
}

// RegistryAuth contains registry authentication information
type RegistryAuth struct {
	Username      string `json:"username" yaml:"username"`
	Password      string `json:"password" yaml:"password"`
	Token         string `json:"token" yaml:"token"`
	IdentityToken string `json:"identity_token" yaml:"identity_token"`
	Auth          string `json:"auth" yaml:"auth"`
	ServerAddress string `json:"server_address" yaml:"server_address"`
}

// PullPolicy defines when to pull base images
type PullPolicy string

const (
	PullPolicyAlways       PullPolicy = "always"
	PullPolicyIfNotPresent PullPolicy = "if-not-present"
	PullPolicyNever        PullPolicy = "never"
)

// OutputFormat defines the output format for built images
type OutputFormat string

const (
	OutputFormatDocker OutputFormat = "docker"
	OutputFormatOCI    OutputFormat = "oci"
	OutputFormatTar    OutputFormat = "tar"
)

// CompressionType defines compression algorithms for layers
type CompressionType string

const (
	CompressionGzip     CompressionType = "gzip"
	CompressionZstd     CompressionType = "zstd"
	CompressionNone     CompressionType = "none"
	CompressionSnappy   CompressionType = "snappy"
)

// NewDefaultBuildOptions returns a BuildOptions struct with sensible defaults
func NewDefaultBuildOptions() BuildOptions {
	return BuildOptions{
		BaseImage:      "gcr.io/distroless/static:nonroot",
		Platform:       "linux/amd64",
		Tag:           "latest",
		Labels:        make(map[string]string),
		Env:           []string{},
		BuildArgs:     make(map[string]string),
		Annotations:   make(map[string]string),
		PullPolicy:    PullPolicyIfNotPresent,
		OutputFormat:  OutputFormatDocker,
		Compression:   CompressionGzip,
		MaxRetries:    3,
		Timeout:       30 * time.Minute,
		Parallelism:   1,
		CreatedBy:     "idpbuilder-oci-go-cr",
	}
}

// NewDefaultBuilderOptions returns a BuilderOptions struct with sensible defaults
func NewDefaultBuilderOptions() BuilderOptions {
	return BuilderOptions{
		Registry:   "localhost:5000",
		Insecure:   false,
		Timeout:    30 * time.Minute,
		MaxRetries: 3,
		LogLevel:   "info",
		CacheDir:   "/tmp/builder-cache",
		TempDir:    "/tmp",
	}
}

// Validate validates the build options
func (opts *BuildOptions) Validate() error {
	if err := opts.validateRequired(); err != nil {
		return err
	}
	
	if err := opts.validateImageName(); err != nil {
		return err
	}
	
	if err := opts.validateTag(); err != nil {
		return err
	}
	
	if err := opts.validatePlatform(); err != nil {
		return err
	}
	
	if err := opts.validateFiles(); err != nil {
		return err
	}
	
	if err := opts.validatePaths(); err != nil {
		return err
	}
	
	if err := opts.validateNumericFields(); err != nil {
		return err
	}
	
	return nil
}

// validateRequired checks that required fields are present
func (opts *BuildOptions) validateRequired() error {
	if opts.ImageName == "" {
		return fmt.Errorf("image name is required")
	}
	
	if opts.Tag == "" {
		return fmt.Errorf("tag is required")
	}
	
	if opts.BaseImage == "" {
		return fmt.Errorf("base image is required")
	}
	
	return nil
}

// validateImageName validates the image name format
func (opts *BuildOptions) validateImageName() error {
	// Basic validation for image name format
	validName := regexp.MustCompile(`^[a-z0-9]+(([._]|__|-+)[a-z0-9]+)*(/[a-z0-9]+(([._]|__|-+)[a-z0-9]+)*)*$`)
	if !validName.MatchString(opts.ImageName) {
		return fmt.Errorf("invalid image name format: %s", opts.ImageName)
	}
	
	return nil
}

// validateTag validates the tag format
func (opts *BuildOptions) validateTag() error {
	// Tags must be valid ASCII and may contain lowercase and uppercase letters, digits, underscores, periods and hyphens
	// They may not start with a period or hyphen and may not contain consecutive periods
	validTag := regexp.MustCompile(`^[a-zA-Z0-9_][a-zA-Z0-9_.-]*[a-zA-Z0-9_]$|^[a-zA-Z0-9_]$`)
	if !validTag.MatchString(opts.Tag) {
		return fmt.Errorf("invalid tag format: %s", opts.Tag)
	}
	
	if len(opts.Tag) > 128 {
		return fmt.Errorf("tag too long (max 128 characters): %s", opts.Tag)
	}
	
	return nil
}

// validatePlatform validates the platform string
func (opts *BuildOptions) validatePlatform() error {
	if opts.Platform == "" {
		return nil // Platform is optional
	}
	
	// Platform should be in format os/arch or os/arch/variant
	parts := strings.Split(opts.Platform, "/")
	if len(parts) < 2 || len(parts) > 3 {
		return fmt.Errorf("invalid platform format, expected os/arch[/variant]: %s", opts.Platform)
	}
	
	validOS := map[string]bool{
		"linux":   true,
		"windows": true,
		"darwin":  true,
		"freebsd": true,
	}
	
	validArch := map[string]bool{
		"amd64": true,
		"arm64": true,
		"arm":   true,
		"386":   true,
		"ppc64": true,
		"s390x": true,
	}
	
	if !validOS[parts[0]] {
		return fmt.Errorf("unsupported OS: %s", parts[0])
	}
	
	if !validArch[parts[1]] {
		return fmt.Errorf("unsupported architecture: %s", parts[1])
	}
	
	return nil
}

// validateFiles validates file specifications
func (opts *BuildOptions) validateFiles() error {
	for i, file := range opts.Files {
		if file.Source == "" {
			return fmt.Errorf("file[%d]: source path is required", i)
		}
		
		if file.Destination == "" {
			return fmt.Errorf("file[%d]: destination path is required", i)
		}
		
		if !filepath.IsAbs(file.Destination) {
			return fmt.Errorf("file[%d]: destination must be absolute path: %s", i, file.Destination)
		}
		
		// If context path is set, validate source exists
		if opts.ContextPath != "" {
			fullSource := filepath.Join(opts.ContextPath, file.Source)
			if _, err := os.Stat(fullSource); err != nil {
				return fmt.Errorf("file[%d]: source file not found: %s", i, fullSource)
			}
		}
	}
	
	return nil
}

// validatePaths validates directory paths
func (opts *BuildOptions) validatePaths() error {
	if opts.ContextPath != "" {
		if !filepath.IsAbs(opts.ContextPath) {
			return fmt.Errorf("context path must be absolute: %s", opts.ContextPath)
		}
		
		if _, err := os.Stat(opts.ContextPath); err != nil {
			return fmt.Errorf("context path not accessible: %w", err)
		}
	}
	
	if opts.OutputFile != "" {
		dir := filepath.Dir(opts.OutputFile)
		if _, err := os.Stat(dir); err != nil {
			return fmt.Errorf("output directory not accessible: %w", err)
		}
	}
	
	if opts.CacheDir != "" {
		if !filepath.IsAbs(opts.CacheDir) {
			return fmt.Errorf("cache directory must be absolute: %s", opts.CacheDir)
		}
	}
	
	return nil
}

// validateNumericFields validates numeric configuration
func (opts *BuildOptions) validateNumericFields() error {
	if opts.MaxRetries < 0 {
		return fmt.Errorf("max retries cannot be negative: %d", opts.MaxRetries)
	}
	
	if opts.MaxRetries > 10 {
		return fmt.Errorf("max retries too high (max 10): %d", opts.MaxRetries)
	}
	
	if opts.Parallelism < 1 {
		return fmt.Errorf("parallelism must be at least 1: %d", opts.Parallelism)
	}
	
	if opts.Parallelism > 16 {
		return fmt.Errorf("parallelism too high (max 16): %d", opts.Parallelism)
	}
	
	if opts.Timeout <= 0 {
		return fmt.Errorf("timeout must be positive: %v", opts.Timeout)
	}
	
	if opts.Timeout > 4*time.Hour {
		return fmt.Errorf("timeout too long (max 4h): %v", opts.Timeout)
	}
	
	return nil
}

// SetDefaults applies default values to unset fields
func (opts *BuildOptions) SetDefaults() {
	defaults := NewDefaultBuildOptions()
	
	if opts.BaseImage == "" {
		opts.BaseImage = defaults.BaseImage
	}
	
	if opts.Platform == "" {
		opts.Platform = defaults.Platform
	}
	
	if opts.Tag == "" {
		opts.Tag = defaults.Tag
	}
	
	if opts.Labels == nil {
		opts.Labels = make(map[string]string)
	}
	
	if opts.BuildArgs == nil {
		opts.BuildArgs = make(map[string]string)
	}
	
	if opts.Annotations == nil {
		opts.Annotations = make(map[string]string)
	}
	
	if opts.PullPolicy == "" {
		opts.PullPolicy = defaults.PullPolicy
	}
	
	if opts.OutputFormat == "" {
		opts.OutputFormat = defaults.OutputFormat
	}
	
	if opts.Compression == "" {
		opts.Compression = defaults.Compression
	}
	
	if opts.MaxRetries == 0 {
		opts.MaxRetries = defaults.MaxRetries
	}
	
	if opts.Timeout == 0 {
		opts.Timeout = defaults.Timeout
	}
	
	if opts.Parallelism == 0 {
		opts.Parallelism = defaults.Parallelism
	}
	
	if opts.CreatedBy == "" {
		opts.CreatedBy = defaults.CreatedBy
	}
	
	// Set standard labels
	opts.Labels["org.opencontainers.image.created"] = time.Now().UTC().Format(time.RFC3339)
	opts.Labels["org.opencontainers.image.source"] = "idpbuilder-oci-go-cr"
}

// Clone creates a deep copy of the build options
func (opts *BuildOptions) Clone() BuildOptions {
	clone := *opts
	
	// Deep copy maps
	if opts.Labels != nil {
		clone.Labels = make(map[string]string, len(opts.Labels))
		for k, v := range opts.Labels {
			clone.Labels[k] = v
		}
	}
	
	if opts.BuildArgs != nil {
		clone.BuildArgs = make(map[string]string, len(opts.BuildArgs))
		for k, v := range opts.BuildArgs {
			clone.BuildArgs[k] = v
		}
	}
	
	if opts.Annotations != nil {
		clone.Annotations = make(map[string]string, len(opts.Annotations))
		for k, v := range opts.Annotations {
			clone.Annotations[k] = v
		}
	}
	
	// Deep copy slices
	if opts.Env != nil {
		clone.Env = make([]string, len(opts.Env))
		copy(clone.Env, opts.Env)
	}
	
	if opts.Cmd != nil {
		clone.Cmd = make([]string, len(opts.Cmd))
		copy(clone.Cmd, opts.Cmd)
	}
	
	if opts.Entrypoint != nil {
		clone.Entrypoint = make([]string, len(opts.Entrypoint))
		copy(clone.Entrypoint, opts.Entrypoint)
	}
	
	if opts.Files != nil {
		clone.Files = make([]FileSpec, len(opts.Files))
		copy(clone.Files, opts.Files)
	}
	
	if opts.CacheFrom != nil {
		clone.CacheFrom = make([]string, len(opts.CacheFrom))
		copy(clone.CacheFrom, opts.CacheFrom)
	}
	
	if opts.CacheTo != nil {
		clone.CacheTo = make([]string, len(opts.CacheTo))
		copy(clone.CacheTo, opts.CacheTo)
	}
	
	// Deep copy registry auth
	if opts.RegistryAuth != nil {
		authCopy := *opts.RegistryAuth
		clone.RegistryAuth = &authCopy
	}
	
	return clone
}

// ToEnvVars converts BuildOptions to environment variables for compatibility
func (opts *BuildOptions) ToEnvVars() map[string]string {
	env := make(map[string]string)
	
	env["BUILDER_IMAGE_NAME"] = opts.ImageName
	env["BUILDER_TAG"] = opts.Tag
	env["BUILDER_BASE_IMAGE"] = opts.BaseImage
	env["BUILDER_PLATFORM"] = opts.Platform
	env["BUILDER_CONTEXT_PATH"] = opts.ContextPath
	env["BUILDER_WORKING_DIR"] = opts.WorkingDir
	env["BUILDER_USER"] = opts.User
	env["BUILDER_REGISTRY"] = opts.Registry
	env["BUILDER_INSECURE"] = strconv.FormatBool(opts.Insecure)
	env["BUILDER_SKIP_TLS"] = strconv.FormatBool(opts.SkipTLS)
	env["BUILDER_NO_CACHE"] = strconv.FormatBool(opts.NoCache)
	env["BUILDER_SQUASH"] = strconv.FormatBool(opts.Squash)
	env["BUILDER_OUTPUT_FORMAT"] = string(opts.OutputFormat)
	env["BUILDER_COMPRESSION"] = string(opts.Compression)
	env["BUILDER_OUTPUT_FILE"] = opts.OutputFile
	env["BUILDER_MAX_RETRIES"] = strconv.Itoa(opts.MaxRetries)
	env["BUILDER_TIMEOUT"] = opts.Timeout.String()
	env["BUILDER_PARALLELISM"] = strconv.Itoa(opts.Parallelism)
	env["BUILDER_CACHE_DIR"] = opts.CacheDir
	env["BUILDER_CREATED_BY"] = opts.CreatedBy
	env["BUILDER_AUTHOR"] = opts.Author
	
	return env
}