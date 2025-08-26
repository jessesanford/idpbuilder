package build

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/containers/buildah"
	"github.com/containers/buildah/define"
	"github.com/containers/common/pkg/config"
	"github.com/containers/image/v5/types"
	"github.com/containers/storage"
	"github.com/opencontainers/runtime-spec/specs-go"
	"github.com/sirupsen/logrus"
)

// BuildConfig represents the main configuration for build operations
type BuildConfig struct {
	// Build configuration
	From                string            `json:"from,omitempty"`
	Tag                 string            `json:"tag,omitempty"`
	File                string            `json:"file,omitempty"`
	Context             string            `json:"context,omitempty"`
	Platform            string            `json:"platform,omitempty"`
	Target              string            `json:"target,omitempty"`
	
	// Build arguments and labels
	Args                map[string]string `json:"args,omitempty"`
	Labels              map[string]string `json:"labels,omitempty"`
	Annotations         map[string]string `json:"annotations,omitempty"`
	
	// Storage configuration
	Storage             *StoreConfig      `json:"storage,omitempty"`
	
	// Runtime configuration
	Runtime             string            `json:"runtime,omitempty"`
	RuntimeArgs         []string          `json:"runtime_args,omitempty"`
	
	// Security and isolation
	Isolation           define.Isolation  `json:"isolation,omitempty"`
	Security            SecurityConfig    `json:"security,omitempty"`
	
	// Network configuration
	Network             NetworkConfig     `json:"network,omitempty"`
	
	// Resource limits
	Resources           ResourceConfig    `json:"resources,omitempty"`
	
	// Output configuration
	Output              OutputConfig      `json:"output,omitempty"`
	
	// System context
	SystemContext       *types.SystemContext `json:"-"`
}

// SecurityConfig contains security-related build options
type SecurityConfig struct {
	NoNewPrivileges     bool     `json:"no_new_privileges,omitempty"`
	AppArmorProfile     string   `json:"apparmor_profile,omitempty"`
	SELinuxLabel        string   `json:"selinux_label,omitempty"`
	SeccompProfile      string   `json:"seccomp_profile,omitempty"`
	Capabilities        []string `json:"capabilities,omitempty"`
	DropCapabilities    []string `json:"drop_capabilities,omitempty"`
	User                string   `json:"user,omitempty"`
	Group               string   `json:"group,omitempty"`
}

// NetworkConfig contains network-related build options
type NetworkConfig struct {
	NetworkMode         string            `json:"network_mode,omitempty"`
	DNSServers          []string          `json:"dns_servers,omitempty"`
	DNSSearch           []string          `json:"dns_search,omitempty"`
	DNSOptions          []string          `json:"dns_options,omitempty"`
	ExtraHosts          map[string]string `json:"extra_hosts,omitempty"`
	HTTPProxy           string            `json:"http_proxy,omitempty"`
	HTTPSProxy          string            `json:"https_proxy,omitempty"`
	NoProxy             string            `json:"no_proxy,omitempty"`
}

// ResourceConfig contains resource limit configurations
type ResourceConfig struct {
	Memory              int64    `json:"memory,omitempty"`
	MemorySwap          int64    `json:"memory_swap,omitempty"`
	CPUShares           uint64   `json:"cpu_shares,omitempty"`
	CPUQuota            int64    `json:"cpu_quota,omitempty"`
	CPUPeriod           uint64   `json:"cpu_period,omitempty"`
	CPUSetCPUs          string   `json:"cpuset_cpus,omitempty"`
	CPUSetMems          string   `json:"cpuset_mems,omitempty"`
	ShmSize             int64    `json:"shm_size,omitempty"`
	Ulimits             []string `json:"ulimits,omitempty"`
}

// OutputConfig contains output and export options
type OutputConfig struct {
	Format              string   `json:"format,omitempty"`
	Compression         string   `json:"compression,omitempty"`
	Push                bool     `json:"push,omitempty"`
	Registry            string   `json:"registry,omitempty"`
	Manifest            string   `json:"manifest,omitempty"`
	SquashLayers        bool     `json:"squash_layers,omitempty"`
	RemoveIntermediates bool     `json:"remove_intermediates,omitempty"`
}

// ConfigManager handles configuration validation and management
type ConfigManager struct {
	config          *BuildConfig
	containerConfig *config.Config
	systemContext   *types.SystemContext
	validated       bool
}

// NewConfigManager creates a new configuration manager
func NewConfigManager(buildConfig *BuildConfig) *ConfigManager {
	if buildConfig == nil {
		buildConfig = DefaultBuildConfig()
	}
	
	return &ConfigManager{
		config:        buildConfig,
		systemContext: buildConfig.SystemContext,
	}
}

// DefaultBuildConfig returns a default build configuration
func DefaultBuildConfig() *BuildConfig {
	return &BuildConfig{
		From:      "scratch",
		Context:   ".",
		Platform:  runtime.GOARCH,
		Args:      make(map[string]string),
		Labels:    make(map[string]string),
		Annotations: make(map[string]string),
		Storage:   DefaultStoreConfig(),
		Runtime:   "runc",
		Isolation: define.IsolationDefault,
		Security: SecurityConfig{
			NoNewPrivileges: false,
			Capabilities:    []string{},
			DropCapabilities: []string{},
		},
		Network: NetworkConfig{
			NetworkMode: "bridge",
			DNSServers:  []string{},
			DNSSearch:   []string{},
			DNSOptions:  []string{},
			ExtraHosts:  make(map[string]string),
		},
		Resources: ResourceConfig{
			Memory:     0,
			MemorySwap: -1,
			CPUShares:  0,
			ShmSize:    64 * 1024 * 1024, // 64MB default
		},
		Output: OutputConfig{
			Format:              "oci",
			Compression:         "",
			Push:                false,
			SquashLayers:        false,
			RemoveIntermediates: true,
		},
		SystemContext: &types.SystemContext{},
	}
}

// Validate performs comprehensive validation of the build configuration
func (cm *ConfigManager) Validate(ctx context.Context) error {
	if cm.validated {
		return nil
	}
	
	config := cm.config
	
	// Validate basic build parameters
	if err := cm.validateBuildParams(); err != nil {
		return fmt.Errorf("build parameters validation failed: %w", err)
	}
	
	// Validate storage configuration
	if err := cm.validateStorageConfig(); err != nil {
		return fmt.Errorf("storage configuration validation failed: %w", err)
	}
	
	// Validate security configuration
	if err := cm.validateSecurityConfig(); err != nil {
		return fmt.Errorf("security configuration validation failed: %w", err)
	}
	
	// Validate network configuration  
	if err := cm.validateNetworkConfig(); err != nil {
		return fmt.Errorf("network configuration validation failed: %w", err)
	}
	
	// Validate resource configuration
	if err := cm.validateResourceConfig(); err != nil {
		return fmt.Errorf("resource configuration validation failed: %w", err)
	}
	
	// Validate output configuration
	if err := cm.validateOutputConfig(); err != nil {
		return fmt.Errorf("output configuration validation failed: %w", err)
	}
	
	// Load container configuration
	if err := cm.loadContainerConfig(); err != nil {
		return fmt.Errorf("failed to load container configuration: %w", err)
	}
	
	cm.validated = true
	logrus.Info("Build configuration validation completed successfully")
	
	return nil
}

// validateBuildParams validates basic build parameters
func (cm *ConfigManager) validateBuildParams() error {
	config := cm.config
	
	// Validate context directory
	if config.Context != "" {
		if !filepath.IsAbs(config.Context) {
			wd, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("failed to get working directory: %w", err)
			}
			config.Context = filepath.Join(wd, config.Context)
		}
		
		if _, err := os.Stat(config.Context); err != nil {
			return fmt.Errorf("context directory %s does not exist: %w", config.Context, err)
		}
	}
	
	// Validate Containerfile/Dockerfile
	if config.File != "" {
		if !filepath.IsAbs(config.File) {
			config.File = filepath.Join(config.Context, config.File)
		}
		
		if _, err := os.Stat(config.File); err != nil {
			return fmt.Errorf("containerfile %s does not exist: %w", config.File, err)
		}
	}
	
	// Validate platform
	if config.Platform != "" {
		validPlatforms := []string{"amd64", "arm64", "arm", "ppc64le", "s390x", "386"}
		valid := false
		for _, p := range validPlatforms {
			if strings.Contains(config.Platform, p) {
				valid = true
				break
			}
		}
		if !valid {
			return fmt.Errorf("unsupported platform: %s", config.Platform)
		}
	}
	
	return nil
}

// validateStorageConfig validates storage configuration
func (cm *ConfigManager) validateStorageConfig() error {
	if cm.config.Storage == nil {
		return nil
	}
	
	storage := cm.config.Storage
	
	// Validate root directory
	if storage.RootDir != "" {
		if err := os.MkdirAll(storage.RootDir, 0755); err != nil {
			return fmt.Errorf("failed to create storage root directory: %w", err)
		}
	}
	
	// Validate run root
	if storage.RunRoot != "" {
		if err := os.MkdirAll(storage.RunRoot, 0755); err != nil {
			return fmt.Errorf("failed to create storage run directory: %w", err)
		}
	}
	
	// Validate graph driver
	validDrivers := []string{"overlay", "vfs", "btrfs", "zfs", "aufs", "devicemapper"}
	if storage.GraphDriver != "" {
		valid := false
		for _, driver := range validDrivers {
			if storage.GraphDriver == driver {
				valid = true
				break
			}
		}
		if !valid {
			return fmt.Errorf("unsupported storage driver: %s", storage.GraphDriver)
		}
	}
	
	return nil
}

// validateSecurityConfig validates security configuration
func (cm *ConfigManager) validateSecurityConfig() error {
	security := &cm.config.Security
	
	// Validate user specification
	if security.User != "" {
		// Parse user:group format
		parts := strings.Split(security.User, ":")
		if len(parts) > 2 {
			return fmt.Errorf("invalid user specification: %s", security.User)
		}
		
		// Validate numeric UID if specified
		if len(parts) >= 1 && parts[0] != "" {
			if _, err := strconv.Atoi(parts[0]); err != nil {
				// Not numeric, assume it's a username (would need system lookup)
				logrus.Debugf("Using username: %s", parts[0])
			}
		}
		
		// Validate group if specified
		if len(parts) == 2 && parts[1] != "" {
			if security.Group == "" {
				security.Group = parts[1]
			}
		}
	}
	
	// Validate capabilities
	if len(security.Capabilities) > 0 {
		for _, cap := range security.Capabilities {
			if !strings.HasPrefix(cap, "CAP_") {
				security.Capabilities = append(security.Capabilities[:0], 
					append([]string{"CAP_" + strings.ToUpper(cap)}, security.Capabilities[1:]...)...)
			}
		}
	}
	
	return nil
}

// validateNetworkConfig validates network configuration
func (cm *ConfigManager) validateNetworkConfig() error {
	network := &cm.config.Network
	
	// Validate network mode
	validModes := []string{"bridge", "host", "none", "container", "private"}
	if network.NetworkMode != "" {
		valid := false
		for _, mode := range validModes {
			if network.NetworkMode == mode || strings.HasPrefix(network.NetworkMode, mode+":") {
				valid = true
				break
			}
		}
		if !valid {
			return fmt.Errorf("invalid network mode: %s", network.NetworkMode)
		}
	}
	
	// Set proxy environment variables if specified
	if network.HTTPProxy != "" {
		os.Setenv("HTTP_PROXY", network.HTTPProxy)
		os.Setenv("http_proxy", network.HTTPProxy)
	}
	if network.HTTPSProxy != "" {
		os.Setenv("HTTPS_PROXY", network.HTTPSProxy)
		os.Setenv("https_proxy", network.HTTPSProxy)
	}
	if network.NoProxy != "" {
		os.Setenv("NO_PROXY", network.NoProxy)
		os.Setenv("no_proxy", network.NoProxy)
	}
	
	return nil
}

// validateResourceConfig validates resource limits configuration
func (cm *ConfigManager) validateResourceConfig() error {
	resources := &cm.config.Resources
	
	// Validate memory settings
	if resources.Memory < 0 {
		return fmt.Errorf("memory limit cannot be negative")
	}
	if resources.MemorySwap != -1 && resources.MemorySwap < resources.Memory {
		return fmt.Errorf("memory swap limit cannot be less than memory limit")
	}
	
	// Validate CPU settings
	if resources.CPUQuota < 0 {
		return fmt.Errorf("CPU quota cannot be negative")
	}
	if resources.CPUPeriod < 0 {
		return fmt.Errorf("CPU period cannot be negative")
	}
	if resources.CPUPeriod > 0 && resources.CPUPeriod < 1000 {
		logrus.Warn("CPU period less than 1000 microseconds may cause performance issues")
	}
	
	// Validate shared memory size
	if resources.ShmSize < 0 {
		return fmt.Errorf("shared memory size cannot be negative")
	}
	
	return nil
}

// validateOutputConfig validates output and export configuration
func (cm *ConfigManager) validateOutputConfig() error {
	output := &cm.config.Output
	
	// Validate output format
	validFormats := []string{"oci", "docker"}
	if output.Format != "" {
		valid := false
		for _, format := range validFormats {
			if output.Format == format {
				valid = true
				break
			}
		}
		if !valid {
			return fmt.Errorf("unsupported output format: %s", output.Format)
		}
	}
	
	// Validate compression
	validCompression := []string{"", "gzip", "bzip2", "xz"}
	if output.Compression != "" {
		valid := false
		for _, comp := range validCompression {
			if output.Compression == comp {
				valid = true
				break
			}
		}
		if !valid {
			return fmt.Errorf("unsupported compression: %s", output.Compression)
		}
	}
	
	return nil
}

// loadContainerConfig loads the containers/common configuration
func (cm *ConfigManager) loadContainerConfig() error {
	containerConfig, err := config.Default()
	if err != nil {
		return fmt.Errorf("failed to load container configuration: %w", err)
	}
	
	cm.containerConfig = containerConfig
	return nil
}

// GetConfig returns the validated build configuration
func (cm *ConfigManager) GetConfig() *BuildConfig {
	return cm.config
}

// GetSystemContext returns the system context for image operations
func (cm *ConfigManager) GetSystemContext() *types.SystemContext {
	if cm.systemContext == nil {
		cm.systemContext = &types.SystemContext{}
	}
	return cm.systemContext
}

// CreateBuildahOptions converts the configuration to Buildah builder options
func (cm *ConfigManager) CreateBuildahOptions(ctx context.Context, store storage.Store) (buildah.BuilderOptions, error) {
	if !cm.validated {
		return buildah.BuilderOptions{}, fmt.Errorf("configuration not validated")
	}
	
	config := cm.config
	
	options := buildah.BuilderOptions{
		FromImage:         config.From,
		Container:         "",
		PullPolicy:        define.PullIfMissing,
		Registry:          config.Output.Registry,
		BlobDirectory:     "",
		SignaturePolicyPath: "",
		ReportWriter:      os.Stderr,
		Store:             store,
		SystemContext:     cm.GetSystemContext(),
		Isolation:         config.Isolation,
		NamespaceOptions:  []define.NamespaceOption{},
		ConfigureNetwork:  define.NetworkDefault,
		CNIPluginPath:     "",
		CNIConfigDir:      "",
		IDMappingOptions:  &define.IDMappingOptions{},
		Capabilities:      config.Security.Capabilities,
		Args:              config.Args,
		Format:            config.Output.Format,
		MaxPullPushRetries: 3,
		PullPushRetryDelay: 2 * time.Second,
	}
	
	// Configure namespace options based on security settings
	if config.Security.User != "" {
		options.NamespaceOptions = append(options.NamespaceOptions,
			define.NamespaceOption{Name: string(specs.UserNamespace)})
	}
	
	// Configure network settings
	switch config.Network.NetworkMode {
	case "host":
		options.ConfigureNetwork = define.NetworkDisabled
		options.NamespaceOptions = append(options.NamespaceOptions,
			define.NamespaceOption{Name: string(specs.NetworkNamespace), Host: true})
	case "none":
		options.ConfigureNetwork = define.NetworkDisabled
	}
	
	return options, nil
}

// ApplyEnvironmentVariables applies environment configuration
func (cm *ConfigManager) ApplyEnvironmentVariables() {
	config := cm.config
	network := &config.Network
	
	// Apply proxy settings
	if network.HTTPProxy != "" {
		os.Setenv("HTTP_PROXY", network.HTTPProxy)
		os.Setenv("http_proxy", network.HTTPProxy)
	}
	if network.HTTPSProxy != "" {
		os.Setenv("HTTPS_PROXY", network.HTTPSProxy)
		os.Setenv("https_proxy", network.HTTPSProxy)
	}
	if network.NoProxy != "" {
		os.Setenv("NO_PROXY", network.NoProxy)
		os.Setenv("no_proxy", network.NoProxy)
	}
}

// UpdateConfig updates the configuration with new values
func (cm *ConfigManager) UpdateConfig(updates map[string]interface{}) error {
	cm.validated = false // Mark as needing re-validation
	
	for key, value := range updates {
		switch key {
		case "from":
			if v, ok := value.(string); ok {
				cm.config.From = v
			}
		case "tag":
			if v, ok := value.(string); ok {
				cm.config.Tag = v
			}
		case "platform":
			if v, ok := value.(string); ok {
				cm.config.Platform = v
			}
		// Add more fields as needed
		default:
			logrus.Warnf("Unknown configuration key: %s", key)
		}
	}
	
	return nil
}