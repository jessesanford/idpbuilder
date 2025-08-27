# E3.1.4 Trust Store (Certificate Storage & Management) Implementation Plan

## 🚨 CRITICAL EFFORT METADATA (FROM WAVE PLAN)
**Branch**: `idpbuidler-oci-mgmt/phase3/wave1/E3.1.4-trust-store`
**Can Parallelize**: Yes (after E3.1.1)
**Parallel With**: [E3.1.2, E3.1.3, E3.1.5]
**Size Estimate**: 600 lines
**Dependencies**: [E3.1.1]

## Overview
- **Effort**: Certificate Storage & Management
- **Phase**: 3, Wave: 1
- **Estimated Size**: 600 lines (target: 550, buffer: 50)
- **Implementation Time**: 8 hours

## Objective
Implement a robust certificate storage and management system that provides persistent certificate storage, automatic discovery, hot-reload capabilities for certificate rotation without service restart, and secure configuration management. This effort forms the foundation for production-grade certificate management in the idpbuilder-oci-mgmt system.

## Dependencies Analysis (R219 Compliance)

### E3.1.1 - Contracts & APIs (Required)
Since E3.1.1 is not yet implemented, this plan anticipates the following interfaces based on the Phase 3 plan:
- `CertificateService` interface with certificate management methods
- Certificate data structures (Certificate, CertificateBundle, VerificationMode)
- Configuration types for certificate management
- Validation interfaces for certificate verification

## File Structure

```
efforts/phase3/wave1/E3.1.4-trust-store/pkg/
└── oci/
    └── certificates/
        ├── storage.go          (250 lines)
        ├── pool.go            (200 lines)
        ├── config.go          (100 lines)
        └── storage_test.go    (50 lines)
```

### File Descriptions

1. **storage.go** (250 lines)
   - `CertificateStore` interface definition
   - `FilesystemStore` struct implementation
   - Persistence methods (Save, Load, Delete, List)
   - Atomic write operations for safety
   - File permission management (0600 for private keys, 0644 for certs)
   - Backup and recovery mechanisms

2. **pool.go** (200 lines)
   - `CertPoolManager` struct for managing certificate pools
   - Dynamic pool updates without service restart
   - Certificate lifecycle management (add, remove, rotate)
   - Watch mechanism for filesystem changes
   - Thread-safe operations with RWMutex
   - Certificate validation before pool updates

3. **config.go** (100 lines)
   - `CertificateConfig` struct definition
   - Environment variable support (IDPBUILDER_CERT_PATH, etc.)
   - Config file parsing (YAML/JSON)
   - Default configuration values
   - Config validation and sanitization
   - Path resolution and expansion

4. **storage_test.go** (50 lines)
   - Unit tests for storage operations
   - Pool management tests
   - Configuration loading tests
   - Security permission verification

## Implementation Steps

### Step 1: Define Storage Interface and Types (50 lines)
```go
// storage.go - Core interfaces and types
type CertificateStore interface {
    Save(ctx context.Context, id string, cert *Certificate) error
    Load(ctx context.Context, id string) (*Certificate, error)
    Delete(ctx context.Context, id string) error
    List(ctx context.Context) ([]string, error)
    Watch(ctx context.Context, onChange func(Event)) error
}

type Event struct {
    Type      EventType // Added, Modified, Deleted
    ID        string
    Timestamp time.Time
}
```

### Step 2: Implement Filesystem Storage (200 lines)
```go
// storage.go - Filesystem implementation
type FilesystemStore struct {
    basePath   string
    mu         sync.RWMutex
    watcher    *fsnotify.Watcher
    validators []CertificateValidator
}

// Implement atomic writes, proper permissions, backup mechanisms
// Include certificate validation before storage
// Add recovery from corrupted certificates
```

### Step 3: Create Certificate Pool Manager (200 lines)
```go
// pool.go - Dynamic pool management
type CertPoolManager struct {
    store       CertificateStore
    systemPool  *x509.CertPool
    customPool  *x509.CertPool
    mu          sync.RWMutex
    updateChan  chan Event
    stopChan    chan struct{}
}

// Implement hot-reload capabilities
// Thread-safe pool updates
// Certificate rotation without downtime
```

### Step 4: Configuration Management (100 lines)
```go
// config.go - Configuration handling
type CertificateConfig struct {
    StoragePath      string        `yaml:"storage_path" env:"IDPBUILDER_CERT_PATH"`
    AutoDiscovery    bool          `yaml:"auto_discovery" env:"IDPBUILDER_CERT_AUTO_DISCOVER"`
    WatchInterval    time.Duration `yaml:"watch_interval" env:"IDPBUILDER_CERT_WATCH_INTERVAL"`
    ValidationMode   string        `yaml:"validation_mode" env:"IDPBUILDER_CERT_VALIDATION"`
    PermissionCheck  bool          `yaml:"permission_check" env:"IDPBUILDER_CERT_PERM_CHECK"`
}

// Support multiple config sources with precedence
// Validate configuration on load
// Provide sensible defaults
```

### Step 5: Implement Auto-Discovery (included in storage.go)
```go
// Auto-discovery of certificates in well-known locations
func (fs *FilesystemStore) DiscoverCertificates() error {
    wellKnownPaths := []string{
        "/etc/ssl/certs",
        "/etc/pki/tls/certs",
        "/usr/local/share/ca-certificates",
        "$HOME/.docker/certs.d",
    }
    // Scan and import certificates from standard locations
}
```

### Step 6: Write Comprehensive Tests (50 lines)
```go
// storage_test.go
// Test persistence and recovery
// Test certificate rotation scenarios
// Test configuration loading from various sources
// Verify security permissions
// Test concurrent access patterns
```

## Size Management
- **Current Estimate**: 600 lines total
- **Measurement Tool**: `$PROJECT_ROOT/tools/line-counter.sh`
- **Check Frequency**: After each major component
- **Split Threshold**: 700 lines (warning), 800 lines (mandatory split)

### Size Breakdown by Component:
- Storage implementation: 250 lines
- Pool management: 200 lines  
- Configuration: 100 lines
- Tests: 50 lines
- **Total**: 600 lines

## Test Requirements

### Unit Tests (90% coverage target):
- **Storage Operations**:
  - Save/Load/Delete/List operations
  - Atomic write verification
  - Permission checks (0600 for keys, 0644 for certs)
  - Corruption recovery

- **Pool Management**:
  - Dynamic updates
  - Thread safety under concurrent access
  - Certificate rotation scenarios
  - Watch mechanism reliability

- **Configuration**:
  - Environment variable parsing
  - Config file loading
  - Validation and defaults
  - Path resolution

### Integration Tests:
- End-to-end certificate lifecycle
- Hot-reload simulation
- Multi-source configuration precedence
- Auto-discovery functionality

### Security Tests:
- File permission verification
- Secure storage of private keys
- Certificate validation before pool updates
- Protection against malicious certificates

## Pattern Compliance

### idpbuilder-oci-mgmt Patterns:
- Use context.Context for all operations
- Implement proper error wrapping with fmt.Errorf
- Follow Go idioms for interface design
- Use structured logging (logr)
- Implement graceful shutdown

### Security Requirements:
- Store private keys with 0600 permissions
- Store certificates with 0644 permissions
- Validate certificates before adding to pool
- Implement certificate pinning support
- Audit all certificate operations

### Performance Targets:
- Certificate load time < 100ms
- Pool update without service interruption
- Watch response time < 500ms
- Support for 1000+ certificates in pool

## Integration Points

### With E3.1.1 (Contracts & APIs):
- Implement CertificateStore interface from E3.1.1
- Use Certificate data structures from E3.1.1
- Follow validation interfaces defined in E3.1.1

### With E3.1.3 (Certificate Service):
- Provide storage backend for CertificateServiceImpl
- Enable certificate persistence for service
- Support service's verification mode requirements

### With E3.1.5 (Integration & Testing):
- Expose mock storage for testing
- Provide test certificate generation utilities
- Support integration test scenarios

## Risk Mitigation

### Technical Risks:
1. **Filesystem watch reliability**: Use polling fallback if inotify fails
2. **Certificate corruption**: Implement backup/recovery mechanisms
3. **Race conditions**: Extensive use of sync.RWMutex
4. **Memory leaks**: Proper cleanup in watch goroutines

### Security Risks:
1. **Permission vulnerabilities**: Strict permission checks on all files
2. **Certificate validation**: Multiple validation layers
3. **Path traversal**: Sanitize all filesystem paths
4. **Concurrent access**: Thread-safe operations throughout

## Success Criteria

### Functional:
- ✅ Persistent certificate storage working
- ✅ Auto-discovery finds system certificates
- ✅ Hot-reload updates pools without restart
- ✅ Configuration from multiple sources

### Non-Functional:
- ✅ All operations thread-safe
- ✅ File permissions correctly enforced
- ✅ 90% test coverage achieved
- ✅ Performance targets met

### Integration:
- ✅ Seamless integration with E3.1.3 service
- ✅ Mock implementations for E3.1.5 testing
- ✅ Compatible with future E3.1.1 contracts

## Implementation Notes

### Key Design Decisions:
1. **Filesystem-based storage**: Simple, reliable, debuggable
2. **Separate system/custom pools**: Flexibility in certificate management
3. **Watch mechanism**: Enables hot-reload without polling overhead
4. **Atomic writes**: Prevents corruption during updates

### Best Practices:
1. Always validate certificates before pool updates
2. Use context for cancellation support
3. Implement comprehensive logging for debugging
4. Design for testability with interfaces
5. Document configuration options thoroughly

### Common Pitfalls to Avoid:
1. Don't hold locks during I/O operations
2. Avoid blocking on watch channels
3. Don't ignore filesystem permission errors
4. Always cleanup watchers on shutdown
5. Handle certificate expiration gracefully

---

**Document Version**: 1.0
**Created**: 2025-08-27
**Author**: Code Reviewer Agent
**Effort**: E3.1.4-trust-store
**Status**: Ready for Implementation