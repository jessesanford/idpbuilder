# Implementation Instructions for E3.1.4 Split-004

## Split Overview
**Split**: 004 of 4 (Final)
**Target Size**: ~394 lines maximum
**Purpose**: Configuration & Integration Tests

## Files to Implement

### 1. pkg/oci/certificates/config.go (200 lines)
Configuration management:
- Implements `ConfigLoader` interface from split-001
- CertificateConfig struct with all settings
- Environment variable loading
- Config validation
- Default configurations

Key structures and methods:
- `type TrustStoreConfig struct` - main configuration
- `NewConfigLoader() *ConfigLoaderImpl`
- `LoadConfig(ctx context.Context, path string) (*Config, error)`
- `SaveConfig(ctx context.Context, config *Config, path string) error`
- `ValidateConfig(config *Config) error`
- `LoadFromEnvironment() *Config` - env var overrides
- Default config factory methods

### 2. pkg/oci/certificates/loader.go (80 lines)
Config file loading utilities:
- YAML/JSON config file parsing
- Environment variable override logic
- Path resolution and expansion
- Config merging logic

Key functions:
- `LoadYAMLConfig(path string) (*Config, error)`
- `LoadJSONConfig(path string) (*Config, error)`
- `MergeConfigs(base, override *Config) *Config`
- `ResolvePaths(config *Config) error`
- `ApplyEnvironmentOverrides(config *Config) error`

### 3. pkg/oci/certificates/integration_test.go (114 lines)
End-to-end integration tests:
- Complete system testing
- Hot-reload scenarios
- Error scenarios
- Multi-pool operations
- Configuration loading tests
- Full workflow validation

Test scenarios:
- Certificate lifecycle (add, validate, use, remove)
- Pool management with hot-reload
- Configuration changes and reloading
- Error recovery scenarios
- Concurrent operations

## Dependencies
- Import from all previous splits (001, 002, 003)
- Use gopkg.in/yaml.v3 for YAML parsing
- Standard library for JSON and environment

## Implementation Requirements

1. **Interface Compliance**: Must implement ConfigLoader interface exactly
2. **Configuration Flexibility**: Support YAML, JSON, and environment variables
3. **Integration Testing**: Test all components working together
4. **Hot-Reload Testing**: Verify configuration can be reloaded
5. **Size Compliance**: Must not exceed 394 lines total

## Branch Setup
Create branch: `idpbuilder-oci-mgmt/phase3/wave1/E3.1.4-trust-store-split-004`

## Implementation Notes
- Configuration precedence: defaults < file < environment
- Support both relative and absolute paths
- Validate all configuration values
- Integration tests should cover all splits working together
- This completes the E3.1.4 effort

## Configuration Structure
```yaml
trustStore:
  storagePath: /var/lib/certificates
  pools:
    system:
      enabled: true
      path: /etc/ssl/certs
    custom:
      enabled: true
      path: /opt/custom-certs
  validation:
    checkExpiry: true
    checkChain: true
    customRules: []
  events:
    enabled: true
    handlers: []
```

## Validation Steps
1. Run `tools/line-counter.sh` to verify size
2. Ensure ConfigLoader interface fully implemented
3. Test configuration loading from all sources
4. Run integration tests covering all splits
5. Verify hot-reload functionality