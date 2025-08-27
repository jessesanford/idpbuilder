package certificates

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// CertificateConfig holds configuration for the certificate management system.
type CertificateConfig struct {
	// StoragePath is the base directory for certificate storage.
	StoragePath string `yaml:"storage_path" json:"storage_path"`
	
	// AutoDiscovery enables automatic discovery of certificates from well-known locations.
	AutoDiscovery bool `yaml:"auto_discovery" json:"auto_discovery"`
	
	// WatchInterval is the minimum interval between filesystem watch events processing.
	WatchInterval time.Duration `yaml:"watch_interval" json:"watch_interval"`
	
	// ValidationMode determines how certificates are validated ("strict", "permissive", "none").
	ValidationMode string `yaml:"validation_mode" json:"validation_mode"`
	
	// PermissionCheck enables verification of file permissions on certificate files.
	PermissionCheck bool `yaml:"permission_check" json:"permission_check"`
	
	// MaxCertificates is the maximum number of certificates allowed in the store.
	MaxCertificates int `yaml:"max_certificates" json:"max_certificates"`
	
	// BackupEnabled determines whether to create backups before modifying certificates.
	BackupEnabled bool `yaml:"backup_enabled" json:"backup_enabled"`
	
	// BackupRetention is how long to keep certificate backups (in hours).
	BackupRetention int `yaml:"backup_retention" json:"backup_retention"`
	
	// TLSConfig contains TLS-specific configuration.
	TLS TLSConfig `yaml:"tls" json:"tls"`
	
	// DiscoveryPaths contains additional paths to search for certificates.
	DiscoveryPaths []string `yaml:"discovery_paths" json:"discovery_paths"`
	
	// LogLevel sets the logging level for certificate operations.
	LogLevel string `yaml:"log_level" json:"log_level"`
	
	// MetricsEnabled determines whether to collect certificate metrics.
	MetricsEnabled bool `yaml:"metrics_enabled" json:"metrics_enabled"`
}

// TLSConfig contains TLS-specific certificate configuration.
type TLSConfig struct {
	// MinVersion is the minimum TLS version to support ("1.0", "1.1", "1.2", "1.3").
	MinVersion string `yaml:"min_version" json:"min_version"`
	
	// MaxVersion is the maximum TLS version to support.
	MaxVersion string `yaml:"max_version" json:"max_version"`
	
	// CipherSuites is a list of allowed cipher suites.
	CipherSuites []string `yaml:"cipher_suites" json:"cipher_suites"`
	
	// InsecureSkipVerify disables certificate verification (for testing only).
	InsecureSkipVerify bool `yaml:"insecure_skip_verify" json:"insecure_skip_verify"`
	
	// ServerName is the server name to verify in the certificate.
	ServerName string `yaml:"server_name" json:"server_name"`
}

// ValidationMode constants for certificate validation.
const (
	ValidationModeStrict     = "strict"
	ValidationModePermissive = "permissive"
	ValidationModeNone       = "none"
)

// LogLevel constants for logging configuration.
const (
	LogLevelDebug = "debug"
	LogLevelInfo  = "info"
	LogLevelWarn  = "warn"
	LogLevelError = "error"
)

// Environment variable names for configuration.
const (
	EnvCertPath             = "IDPBUILDER_CERT_PATH"
	EnvCertAutoDiscover     = "IDPBUILDER_CERT_AUTO_DISCOVER"
	EnvCertWatchInterval    = "IDPBUILDER_CERT_WATCH_INTERVAL"
	EnvCertValidation       = "IDPBUILDER_CERT_VALIDATION"
	EnvCertPermCheck        = "IDPBUILDER_CERT_PERM_CHECK"
	EnvCertMaxCerts         = "IDPBUILDER_CERT_MAX_CERTS"
	EnvCertBackupEnabled    = "IDPBUILDER_CERT_BACKUP_ENABLED"
	EnvCertBackupRetention  = "IDPBUILDER_CERT_BACKUP_RETENTION"
	EnvCertTLSMinVersion    = "IDPBUILDER_CERT_TLS_MIN_VERSION"
	EnvCertTLSMaxVersion    = "IDPBUILDER_CERT_TLS_MAX_VERSION"
	EnvCertTLSInsecure      = "IDPBUILDER_CERT_TLS_INSECURE"
	EnvCertTLSServerName    = "IDPBUILDER_CERT_TLS_SERVER_NAME"
	EnvCertDiscoveryPaths   = "IDPBUILDER_CERT_DISCOVERY_PATHS"
	EnvCertLogLevel         = "IDPBUILDER_CERT_LOG_LEVEL"
	EnvCertMetricsEnabled   = "IDPBUILDER_CERT_METRICS_ENABLED"
)

// DefaultConfig returns a certificate configuration with sensible defaults.
func DefaultConfig() *CertificateConfig {
	homeDir := os.Getenv("HOME")
	if homeDir == "" {
		homeDir = "/tmp"
	}
	
	return &CertificateConfig{
		StoragePath:     filepath.Join(homeDir, ".idpbuilder", "certificates"),
		AutoDiscovery:   true,
		WatchInterval:   5 * time.Second,
		ValidationMode:  ValidationModeStrict,
		PermissionCheck: true,
		MaxCertificates: 1000,
		BackupEnabled:   true,
		BackupRetention: 168, // 7 days
		TLS: TLSConfig{
			MinVersion:         "1.2",
			MaxVersion:         "1.3",
			InsecureSkipVerify: false,
		},
		DiscoveryPaths: []string{
			"/etc/ssl/certs",
			"/etc/pki/tls/certs",
			"/usr/local/share/ca-certificates",
		},
		LogLevel:       LogLevelInfo,
		MetricsEnabled: false,
	}
}

// LoadConfig loads certificate configuration from multiple sources with precedence:
// 1. Environment variables (highest precedence)
// 2. Configuration file
// 3. Default values (lowest precedence)
func LoadConfig(configPath string) (*CertificateConfig, error) {
	config := DefaultConfig()
	
	// Load from configuration file if specified
	if configPath != "" {
		if err := config.loadFromFile(configPath); err != nil {
			return nil, fmt.Errorf("failed to load config from file %s: %w", configPath, err)
		}
	}
	
	// Override with environment variables
	config.loadFromEnvironment()
	
	// Validate the configuration
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}
	
	// Expand paths
	if err := config.expandPaths(); err != nil {
		return nil, fmt.Errorf("failed to expand paths: %w", err)
	}
	
	return config, nil
}

// LoadConfigFromFile loads configuration from a YAML or JSON file.
func LoadConfigFromFile(configPath string) (*CertificateConfig, error) {
	return LoadConfig(configPath)
}

// LoadConfigFromEnv loads configuration from environment variables only.
func LoadConfigFromEnv() (*CertificateConfig, error) {
	return LoadConfig("")
}

// loadFromFile loads configuration from a file.
func (c *CertificateConfig) loadFromFile(configPath string) error {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}
	
	// Try to unmarshal as YAML (which also handles JSON)
	if err := yaml.Unmarshal(data, c); err != nil {
		return fmt.Errorf("failed to parse config file: %w", err)
	}
	
	return nil
}

// loadFromEnvironment loads configuration from environment variables.
func (c *CertificateConfig) loadFromEnvironment() {
	if path := os.Getenv(EnvCertPath); path != "" {
		c.StoragePath = path
	}
	
	if autoDiscover := os.Getenv(EnvCertAutoDiscover); autoDiscover != "" {
		c.AutoDiscovery = parseBool(autoDiscover, c.AutoDiscovery)
	}
	
	if watchInterval := os.Getenv(EnvCertWatchInterval); watchInterval != "" {
		if duration, err := time.ParseDuration(watchInterval); err == nil {
			c.WatchInterval = duration
		}
	}
	
	if validation := os.Getenv(EnvCertValidation); validation != "" {
		c.ValidationMode = validation
	}
	
	if permCheck := os.Getenv(EnvCertPermCheck); permCheck != "" {
		c.PermissionCheck = parseBool(permCheck, c.PermissionCheck)
	}
	
	if maxCerts := os.Getenv(EnvCertMaxCerts); maxCerts != "" {
		if max, err := strconv.Atoi(maxCerts); err == nil && max > 0 {
			c.MaxCertificates = max
		}
	}
	
	if backupEnabled := os.Getenv(EnvCertBackupEnabled); backupEnabled != "" {
		c.BackupEnabled = parseBool(backupEnabled, c.BackupEnabled)
	}
	
	if backupRetention := os.Getenv(EnvCertBackupRetention); backupRetention != "" {
		if retention, err := strconv.Atoi(backupRetention); err == nil && retention > 0 {
			c.BackupRetention = retention
		}
	}
	
	if minVersion := os.Getenv(EnvCertTLSMinVersion); minVersion != "" {
		c.TLS.MinVersion = minVersion
	}
	
	if maxVersion := os.Getenv(EnvCertTLSMaxVersion); maxVersion != "" {
		c.TLS.MaxVersion = maxVersion
	}
	
	if insecure := os.Getenv(EnvCertTLSInsecure); insecure != "" {
		c.TLS.InsecureSkipVerify = parseBool(insecure, c.TLS.InsecureSkipVerify)
	}
	
	if serverName := os.Getenv(EnvCertTLSServerName); serverName != "" {
		c.TLS.ServerName = serverName
	}
	
	if discoveryPaths := os.Getenv(EnvCertDiscoveryPaths); discoveryPaths != "" {
		c.DiscoveryPaths = strings.Split(discoveryPaths, ":")
	}
	
	if logLevel := os.Getenv(EnvCertLogLevel); logLevel != "" {
		c.LogLevel = logLevel
	}
	
	if metricsEnabled := os.Getenv(EnvCertMetricsEnabled); metricsEnabled != "" {
		c.MetricsEnabled = parseBool(metricsEnabled, c.MetricsEnabled)
	}
}

// Validate validates the configuration for consistency and correctness.
func (c *CertificateConfig) Validate() error {
	if c.StoragePath == "" {
		return fmt.Errorf("storage_path cannot be empty")
	}
	
	if c.WatchInterval <= 0 {
		return fmt.Errorf("watch_interval must be positive")
	}
	
	// Validate validation mode
	switch c.ValidationMode {
	case ValidationModeStrict, ValidationModePermissive, ValidationModeNone:
		// Valid
	default:
		return fmt.Errorf("invalid validation_mode: %s (must be one of: strict, permissive, none)", c.ValidationMode)
	}
	
	if c.MaxCertificates < 0 {
		return fmt.Errorf("max_certificates cannot be negative")
	}
	
	if c.BackupRetention < 0 {
		return fmt.Errorf("backup_retention cannot be negative")
	}
	
	// Validate TLS configuration
	if err := c.TLS.validate(); err != nil {
		return fmt.Errorf("TLS configuration invalid: %w", err)
	}
	
	// Validate log level
	switch c.LogLevel {
	case LogLevelDebug, LogLevelInfo, LogLevelWarn, LogLevelError:
		// Valid
	default:
		return fmt.Errorf("invalid log_level: %s (must be one of: debug, info, warn, error)", c.LogLevel)
	}
	
	return nil
}

// validate validates the TLS configuration.
func (t *TLSConfig) validate() error {
	validVersions := map[string]bool{
		"1.0": true,
		"1.1": true,
		"1.2": true,
		"1.3": true,
	}
	
	if t.MinVersion != "" && !validVersions[t.MinVersion] {
		return fmt.Errorf("invalid min_version: %s (must be one of: 1.0, 1.1, 1.2, 1.3)", t.MinVersion)
	}
	
	if t.MaxVersion != "" && !validVersions[t.MaxVersion] {
		return fmt.Errorf("invalid max_version: %s (must be one of: 1.0, 1.1, 1.2, 1.3)", t.MaxVersion)
	}
	
	return nil
}

// expandPaths expands environment variables and ~ in configuration paths.
func (c *CertificateConfig) expandPaths() error {
	var err error
	
	// Expand storage path
	c.StoragePath, err = expandPath(c.StoragePath)
	if err != nil {
		return fmt.Errorf("failed to expand storage_path: %w", err)
	}
	
	// Expand discovery paths
	for i, path := range c.DiscoveryPaths {
		expandedPath, err := expandPath(path)
		if err != nil {
			return fmt.Errorf("failed to expand discovery_path[%d]: %w", i, err)
		}
		c.DiscoveryPaths[i] = expandedPath
	}
	
	return nil
}

// SaveToFile saves the configuration to a YAML file.
func (c *CertificateConfig) SaveToFile(configPath string) error {
	data, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}
	
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}
	
	return nil
}

// Clone creates a deep copy of the configuration.
func (c *CertificateConfig) Clone() *CertificateConfig {
	clone := *c
	
	// Deep copy slices
	clone.DiscoveryPaths = make([]string, len(c.DiscoveryPaths))
	copy(clone.DiscoveryPaths, c.DiscoveryPaths)
	
	clone.TLS.CipherSuites = make([]string, len(c.TLS.CipherSuites))
	copy(clone.TLS.CipherSuites, c.TLS.CipherSuites)
	
	return &clone
}

// IsStrictValidation returns true if strict validation mode is enabled.
func (c *CertificateConfig) IsStrictValidation() bool {
	return c.ValidationMode == ValidationModeStrict
}

// IsPermissiveValidation returns true if permissive validation mode is enabled.
func (c *CertificateConfig) IsPermissiveValidation() bool {
	return c.ValidationMode == ValidationModePermissive
}

// IsValidationDisabled returns true if validation is disabled.
func (c *CertificateConfig) IsValidationDisabled() bool {
	return c.ValidationMode == ValidationModeNone
}

// GetEffectiveDiscoveryPaths returns all discovery paths, including default ones if auto-discovery is enabled.
func (c *CertificateConfig) GetEffectiveDiscoveryPaths() []string {
	if !c.AutoDiscovery {
		return nil
	}
	
	paths := make([]string, len(c.DiscoveryPaths))
	copy(paths, c.DiscoveryPaths)
	
	return paths
}

// expandPath expands environment variables and ~ in a path.
func expandPath(path string) (string, error) {
	if path == "" {
		return path, nil
	}
	
	// Expand ~ to home directory
	if strings.HasPrefix(path, "~/") {
		homeDir := os.Getenv("HOME")
		if homeDir == "" {
			return path, fmt.Errorf("HOME environment variable not set")
		}
		path = filepath.Join(homeDir, path[2:])
	}
	
	// Expand environment variables
	path = os.ExpandEnv(path)
	
	// Convert to absolute path
	absPath, err := filepath.Abs(path)
	if err != nil {
		return path, fmt.Errorf("failed to get absolute path: %w", err)
	}
	
	return absPath, nil
}

// parseBool parses a string as a boolean, returning the default value on error.
func parseBool(s string, defaultValue bool) bool {
	if b, err := strconv.ParseBool(s); err == nil {
		return b
	}
	return defaultValue
}