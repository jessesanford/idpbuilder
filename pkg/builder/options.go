package builder

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// BuildOption defines a functional option for configuring a BuildConfig.
type BuildOption func(*BuildConfig) error

// WithContextPath sets the build context path.
func WithContextPath(path string) BuildOption {
	return func(c *BuildConfig) error {
		if path == "" {
			return fmt.Errorf("context path cannot be empty")
		}
		c.ContextPath = path
		return nil
	}
}

// WithDockerfile sets the Dockerfile path relative to the context.
func WithDockerfile(dockerfile string) BuildOption {
	return func(c *BuildConfig) error {
		if dockerfile == "" {
			return fmt.Errorf("dockerfile path cannot be empty")
		}
		c.Dockerfile = dockerfile
		return nil
	}
}

// WithTags sets the image tags.
func WithTags(tags ...string) BuildOption {
	return func(c *BuildConfig) error {
		if len(tags) == 0 {
			return fmt.Errorf("at least one tag must be specified")
		}
		for _, tag := range tags {
			if strings.TrimSpace(tag) == "" {
				return fmt.Errorf("tag cannot be empty or whitespace")
			}
		}
		c.Tags = make([]string, len(tags))
		copy(c.Tags, tags)
		return nil
	}
}

// WithPlatform sets the target platform using a platform string (e.g., "linux/amd64").
func WithPlatform(platform string) BuildOption {
	return func(c *BuildConfig) error {
		parts := strings.Split(platform, "/")
		if len(parts) < 2 {
			return fmt.Errorf("invalid platform format: %s, expected OS/Architecture[/Variant]", platform)
		}

		c.Platform.OS = parts[0]
		c.Platform.Architecture = parts[1]
		
		if len(parts) > 2 {
			c.Platform.Variant = parts[2]
		}

		return c.Platform.Validate()
	}
}

// WithPlatformConfig sets the platform configuration.
func WithPlatformConfig(platform PlatformConfig) BuildOption {
	return func(c *BuildConfig) error {
		if err := platform.Validate(); err != nil {
			return fmt.Errorf("invalid platform config: %w", err)
		}
		c.Platform = platform
		return nil
	}
}

// WithRegistry sets the registry configuration.
func WithRegistry(registry RegistryConfig) BuildOption {
	return func(c *BuildConfig) error {
		if err := registry.Validate(); err != nil {
			return fmt.Errorf("invalid registry config: %w", err)
		}
		c.Registry = registry
		return nil
	}
}

// WithRegistryAuth sets registry authentication using username and password.
func WithRegistryAuth(hostname, username, password string) BuildOption {
	return func(c *BuildConfig) error {
		if hostname == "" {
			return fmt.Errorf("registry hostname cannot be empty")
		}
		if username == "" || password == "" {
			return fmt.Errorf("both username and password must be provided")
		}
		
		c.Registry.Hostname = hostname
		c.Registry.Username = username
		c.Registry.Password = password
		
		return c.Registry.Validate()
	}
}

// WithRegistryToken sets registry authentication using a token.
func WithRegistryToken(hostname, token string) BuildOption {
	return func(c *BuildConfig) error {
		if hostname == "" {
			return fmt.Errorf("registry hostname cannot be empty")
		}
		if token == "" {
			return fmt.Errorf("token cannot be empty")
		}
		
		c.Registry.Hostname = hostname
		c.Registry.Token = token
		
		return c.Registry.Validate()
	}
}

// WithInsecureRegistry enables insecure registry connections.
func WithInsecureRegistry(insecure bool) BuildOption {
	return func(c *BuildConfig) error {
		c.Registry.Insecure = insecure
		return nil
	}
}

// WithBuildArg adds a build argument.
func WithBuildArg(key, value string) BuildOption {
	return func(c *BuildConfig) error {
		if key == "" {
			return fmt.Errorf("build arg key cannot be empty")
		}
		if c.BuildArgs == nil {
			c.BuildArgs = make(map[string]string)
		}
		c.BuildArgs[key] = value
		return nil
	}
}

// WithBuildArgs sets multiple build arguments.
func WithBuildArgs(args map[string]string) BuildOption {
	return func(c *BuildConfig) error {
		if args == nil {
			return nil
		}
		if c.BuildArgs == nil {
			c.BuildArgs = make(map[string]string)
		}
		for k, v := range args {
			if k == "" {
				return fmt.Errorf("build arg key cannot be empty")
			}
			c.BuildArgs[k] = v
		}
		return nil
	}
}

// WithLabel adds a metadata label.
func WithLabel(key, value string) BuildOption {
	return func(c *BuildConfig) error {
		if key == "" {
			return fmt.Errorf("label key cannot be empty")
		}
		if c.Labels == nil {
			c.Labels = make(map[string]string)
		}
		c.Labels[key] = value
		return nil
	}
}

// WithLabels sets multiple metadata labels.
func WithLabels(labels map[string]string) BuildOption {
	return func(c *BuildConfig) error {
		if labels == nil {
			return nil
		}
		if c.Labels == nil {
			c.Labels = make(map[string]string)
		}
		for k, v := range labels {
			if k == "" {
				return fmt.Errorf("label key cannot be empty")
			}
			c.Labels[k] = v
		}
		return nil
	}
}

// WithTarget sets the target stage for multi-stage builds.
func WithTarget(target string) BuildOption {
	return func(c *BuildConfig) error {
		c.Target = target
		return nil
	}
}

// WithNoCache disables layer caching.
func WithNoCache(noCache bool) BuildOption {
	return func(c *BuildConfig) error {
		c.NoCache = noCache
		return nil
	}
}

// WithPull forces pulling of base images.
func WithPull(pull bool) BuildOption {
	return func(c *BuildConfig) error {
		c.Pull = pull
		return nil
	}
}

// WithRemove sets whether to remove intermediate containers.
func WithRemove(remove bool) BuildOption {
	return func(c *BuildConfig) error {
		c.Remove = remove
		return nil
	}
}

// WithSquash enables layer squashing (experimental).
func WithSquash(squash bool) BuildOption {
	return func(c *BuildConfig) error {
		c.Squash = squash
		return nil
	}
}

// WithBuildTimeout sets the build timeout.
func WithBuildTimeout(timeout time.Duration) BuildOption {
	return func(c *BuildConfig) error {
		if timeout <= 0 {
			return fmt.Errorf("build timeout must be greater than zero")
		}
		c.BuildTimeout = timeout
		return nil
	}
}

// WithBuildTimeoutSeconds sets the build timeout in seconds.
func WithBuildTimeoutSeconds(seconds int) BuildOption {
	return func(c *BuildConfig) error {
		if seconds <= 0 {
			return fmt.Errorf("build timeout seconds must be greater than zero")
		}
		c.BuildTimeout = time.Duration(seconds) * time.Second
		return nil
	}
}

// WithMemoryLimit sets the memory limit for build containers.
func WithMemoryLimit(bytes int64) BuildOption {
	return func(c *BuildConfig) error {
		if bytes < 0 {
			return fmt.Errorf("memory limit cannot be negative")
		}
		c.MemoryLimit = bytes
		return nil
	}
}

// WithMemoryLimitMB sets the memory limit in megabytes.
func WithMemoryLimitMB(mb int64) BuildOption {
	return func(c *BuildConfig) error {
		if mb < 0 {
			return fmt.Errorf("memory limit cannot be negative")
		}
		c.MemoryLimit = mb * 1024 * 1024
		return nil
	}
}

// WithCPULimit sets the CPU limit for build containers.
func WithCPULimit(limit float64) BuildOption {
	return func(c *BuildConfig) error {
		if limit < 0 {
			return fmt.Errorf("CPU limit cannot be negative")
		}
		c.CPULimit = limit
		return nil
	}
}

// WithCPULimitString sets the CPU limit from a string (e.g., "0.5", "1.0").
func WithCPULimitString(limit string) BuildOption {
	return func(c *BuildConfig) error {
		if limit == "" {
			return fmt.Errorf("CPU limit string cannot be empty")
		}
		
		cpuLimit, err := strconv.ParseFloat(limit, 64)
		if err != nil {
			return fmt.Errorf("invalid CPU limit format: %s", limit)
		}
		
		if cpuLimit < 0 {
			return fmt.Errorf("CPU limit cannot be negative")
		}
		
		c.CPULimit = cpuLimit
		return nil
	}
}

// ApplyOptions applies a slice of BuildOption to a BuildConfig.
func ApplyOptions(config *BuildConfig, options ...BuildOption) error {
	for i, opt := range options {
		if err := opt(config); err != nil {
			return fmt.Errorf("option %d failed: %w", i, err)
		}
	}
	return nil
}

// NewBuildConfig creates a new BuildConfig with the given options applied to defaults.
func NewBuildConfig(options ...BuildOption) (*BuildConfig, error) {
	config := DefaultBuildConfig()
	
	if err := ApplyOptions(config, options...); err != nil {
		return nil, fmt.Errorf("failed to apply options: %w", err)
	}
	
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}
	
	return config, nil
}

// MustNewBuildConfig creates a new BuildConfig with the given options, panicking on error.
func MustNewBuildConfig(options ...BuildOption) *BuildConfig {
	config, err := NewBuildConfig(options...)
	if err != nil {
		panic(fmt.Sprintf("failed to create build config: %v", err))
	}
	return config
}

// ChainOptions chains multiple BuildOptions into a single option.
func ChainOptions(options ...BuildOption) BuildOption {
	return func(config *BuildConfig) error {
		return ApplyOptions(config, options...)
	}
}

// ConditionalOption applies an option only if the condition is true.
func ConditionalOption(condition bool, option BuildOption) BuildOption {
	return func(config *BuildConfig) error {
		if condition {
			return option(config)
		}
		return nil
	}
}