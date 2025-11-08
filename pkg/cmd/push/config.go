package push

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// ConfigSource represents the source of a configuration value
type ConfigSource int

const (
	// SourceDefault indicates the value came from the default
	SourceDefault ConfigSource = iota
	// SourceEnv indicates the value came from an environment variable
	SourceEnv
	// SourceFlag indicates the value came from a command-line flag
	SourceFlag
)

// String returns a string representation of the ConfigSource
func (s ConfigSource) String() string {
	switch s {
	case SourceDefault:
		return "default"
	case SourceEnv:
		return "environment"
	case SourceFlag:
		return "flag"
	default:
		return "unknown"
	}
}

// ConfigValue represents a configuration value with its source
type ConfigValue struct {
	Value  string
	Source ConfigSource
}

// PushConfig holds the push command configuration with source tracking
type PushConfig struct {
	ImageName ConfigValue
	Registry  ConfigValue
	Username  ConfigValue
	Password  ConfigValue
	Insecure  ConfigValue
	Verbose   ConfigValue
}

// Environment variable names for configuration
const (
	EnvRegistry = "IDPBUILDER_REGISTRY"
	EnvUsername = "IDPBUILDER_USERNAME"
	EnvPassword = "IDPBUILDER_PASSWORD"
	EnvInsecure = "IDPBUILDER_INSECURE"
	EnvVerbose  = "IDPBUILDER_VERBOSE"
)

// DefaultRegistry is the default registry URL
const DefaultRegistry = "gitea.cnoe.localtest.me:8443"

// LoadConfig loads configuration from flags, environment variables, and defaults
// following the precedence: Flags > Environment > Defaults
func LoadConfig(cmd *cobra.Command, args []string, v *viper.Viper) (*PushConfig, error) {
	if len(args) < 1 {
		return nil, fmt.Errorf("image name is required as first argument")
	}

	config := &PushConfig{
		// Image name always comes from args
		ImageName: ConfigValue{
			Value:  args[0],
			Source: SourceFlag,
		},
	}

	// Resolve string configurations
	config.Registry = resolveStringConfig(cmd, "registry", EnvRegistry, DefaultRegistry)
	config.Username = resolveStringConfig(cmd, "username", EnvUsername, "")
	config.Password = resolveStringConfig(cmd, "password", EnvPassword, "")

	// Resolve boolean configurations
	config.Insecure = resolveBoolConfig(cmd, "insecure", EnvInsecure, false)
	config.Verbose = resolveBoolConfig(cmd, "verbose", EnvVerbose, false)

	return config, nil
}

// resolveStringConfig resolves a string configuration value
// following precedence: Flag > Environment > Default
func resolveStringConfig(cmd *cobra.Command, flagName, envName, defaultValue string) ConfigValue {
	// Check if flag was explicitly set
	flag := cmd.Flags().Lookup(flagName)
	if flag != nil && flag.Changed {
		return ConfigValue{
			Value:  flag.Value.String(),
			Source: SourceFlag,
		}
	}

	// Check environment variable
	if envValue := os.Getenv(envName); envValue != "" {
		return ConfigValue{
			Value:  envValue,
			Source: SourceEnv,
		}
	}

	// Use default value
	return ConfigValue{
		Value:  defaultValue,
		Source: SourceDefault,
	}
}

// resolveBoolConfig resolves a boolean configuration value
// following precedence: Flag > Environment > Default
// Supports formats: true/false, 1/0, yes/no (case-insensitive)
func resolveBoolConfig(cmd *cobra.Command, flagName, envName string, defaultValue bool) ConfigValue {
	// Check if flag was explicitly set
	flag := cmd.Flags().Lookup(flagName)
	if flag != nil && flag.Changed {
		return ConfigValue{
			Value:  flag.Value.String(),
			Source: SourceFlag,
		}
	}

	// Check environment variable
	if envValue := os.Getenv(envName); envValue != "" {
		// Normalize the environment value
		normalized := strings.ToLower(strings.TrimSpace(envValue))

		switch normalized {
		case "true", "1", "yes":
			return ConfigValue{
				Value:  "true",
				Source: SourceEnv,
			}
		case "false", "0", "no":
			return ConfigValue{
				Value:  "false",
				Source: SourceEnv,
			}
		default:
			// Invalid value, fall back to default
			return ConfigValue{
				Value:  fmt.Sprintf("%t", defaultValue),
				Source: SourceDefault,
			}
		}
	}

	// Use default value
	return ConfigValue{
		Value:  fmt.Sprintf("%t", defaultValue),
		Source: SourceDefault,
	}
}

// ToPushOptions converts PushConfig to PushOptions for Wave 2.1 compatibility
func (c *PushConfig) ToPushOptions() *PushOptions {
	return &PushOptions{
		ImageName: c.ImageName.Value,
		Registry:  c.Registry.Value,
		Username:  c.Username.Value,
		Password:  c.Password.Value,
		Insecure:  c.Insecure.Value == "true",
		Verbose:   c.Verbose.Value == "true",
	}
}

// Validate checks if the configuration is valid and complete
func (c *PushConfig) Validate() error {
	if c.ImageName.Value == "" {
		return fmt.Errorf("image name is required")
	}
	if c.Username.Value == "" {
		return fmt.Errorf("username is required (use --username flag or %s environment variable)", EnvUsername)
	}
	if c.Password.Value == "" {
		return fmt.Errorf("password is required (use --password flag or %s environment variable)", EnvPassword)
	}
	if c.Registry.Value == "" {
		return fmt.Errorf("registry is required (use --registry flag or %s environment variable)", EnvRegistry)
	}
	return nil
}

// DisplaySources displays the source of each configuration value (for verbose mode)
func (c *PushConfig) DisplaySources() {
	fmt.Println("Configuration sources:")
	fmt.Printf("  Image:    %s (from %s)\n", c.ImageName.Value, c.ImageName.Source)
	fmt.Printf("  Registry: %s (from %s)\n", c.Registry.Value, c.Registry.Source)
	fmt.Printf("  Username: %s (from %s)\n", c.Username.Value, c.Username.Source)
	fmt.Printf("  Password: %s (from %s)\n", "***", c.Password.Source)
	fmt.Printf("  Insecure: %s (from %s)\n", c.Insecure.Value, c.Insecure.Source)
	fmt.Printf("  Verbose:  %s (from %s)\n", c.Verbose.Value, c.Verbose.Source)
}
