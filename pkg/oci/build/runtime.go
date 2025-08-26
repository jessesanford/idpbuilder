// Package build provides OCI build capabilities using Buildah
// This file implements runtime management and rootless operation setup
package build

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/containers/buildah/define"
	"github.com/containers/storage"
	"github.com/containers/storage/pkg/idtools"
	"github.com/opencontainers/runtime-spec/specs-go"
	"github.com/pkg/errors"
)

// BuildConfig represents the configuration for OCI builds
// This should match the definition from split-001
type BuildConfig struct {
	// Storage configuration
	StorageDriver   string            `json:"storageDriver" yaml:"storageDriver"`
	StorageRoot     string            `json:"storageRoot" yaml:"storageRoot"`
	RunRoot         string            `json:"runRoot" yaml:"runRoot"`
	StorageOptions  map[string]string `json:"storageOptions" yaml:"storageOptions"`
	
	// Build configuration
	Isolation       string            `json:"isolation" yaml:"isolation"`
	CgroupManager   string            `json:"cgroupManager" yaml:"cgroupManager"`
	DefaultMounts   []string          `json:"defaultMounts" yaml:"defaultMounts"`
	
	// Security configuration
	UIDMap          []specs.LinuxIDMapping `json:"uidMap" yaml:"uidMap"`
	GIDMap          []specs.LinuxIDMapping `json:"gidMap" yaml:"gidMap"`
	Capabilities    []string          `json:"capabilities" yaml:"capabilities"`
	
	// Runtime paths
	RuntimePath     string            `json:"runtimePath" yaml:"runtimePath"`
	ConmonPath      string            `json:"conmonPath" yaml:"conmonPath"`
	CNIConfigDir    string            `json:"cniConfigDir" yaml:"cniConfigDir"`
	CNIPluginDir    string            `json:"cniPluginDir" yaml:"cniPluginDir"`
}

// RuntimeManager manages the container runtime environment for Buildah operations
type RuntimeManager struct {
	config      *BuildConfig
	store       storage.Store
	initialized bool
	rootless    bool
	userNS      bool
}

// NewRuntimeManager creates a new RuntimeManager with the provided configuration
func NewRuntimeManager(config *BuildConfig) (*RuntimeManager, error) {
	if config == nil {
		return nil, errors.New("configuration cannot be nil")
	}

	rm := &RuntimeManager{
		config:   config,
		rootless: os.Geteuid() != 0,
	}

	if err := rm.validateConfiguration(); err != nil {
		return nil, errors.Wrap(err, "invalid configuration")
	}

	return rm, nil
}

// validateConfiguration validates the runtime configuration
func (rm *RuntimeManager) validateConfiguration() error {
	if rm.config.StorageDriver == "" {
		return errors.New("storage driver is required")
	}
	if rm.config.StorageRoot == "" {
		return errors.New("storage root is required")
	}
	if rm.config.RunRoot == "" {
		return errors.New("run root is required")
	}

	// Validate runtime paths if specified
	for path, name := range map[string]string{rm.config.RuntimePath: "runtime", rm.config.ConmonPath: "conmon"} {
		if path != "" {
			if _, err := os.Stat(path); os.IsNotExist(err) {
				return errors.Errorf("%s path does not exist: %s", name, path)
			}
		}
	}
	return nil
}

// Initialize initializes the runtime environment
func (rm *RuntimeManager) Initialize(ctx context.Context) error {
	if rm.initialized {
		return nil
	}

	// Setup directories
	if err := rm.setupDirectories(); err != nil {
		return errors.Wrap(err, "failed to setup directories")
	}

	// Initialize storage
	if err := rm.initializeStorage(); err != nil {
		return errors.Wrap(err, "failed to initialize storage")
	}

	// Setup rootless environment if needed
	if rm.rootless {
		if err := rm.setupRootlessEnvironment(); err != nil {
			return errors.Wrap(err, "failed to setup rootless environment")
		}
	}

	// Setup security capabilities
	if err := rm.setupSecurityCapabilities(); err != nil {
		return errors.Wrap(err, "failed to setup security capabilities")
	}

	rm.initialized = true
	return nil
}

// setupDirectories creates and validates required directories
func (rm *RuntimeManager) setupDirectories() error {
	dirs := []string{rm.config.StorageRoot, rm.config.RunRoot, rm.config.CNIConfigDir, rm.config.CNIPluginDir}
	for _, dir := range dirs {
		if dir != "" {
			if err := os.MkdirAll(dir, 0755); err != nil {
				return errors.Wrapf(err, "failed to create directory: %s", dir)
			}
		}
	}
	return nil
}

// initializeStorage initializes the container storage
func (rm *RuntimeManager) initializeStorage() error {
	storeOptions := storage.StoreOptions{
		RunRoot:         rm.config.RunRoot,
		GraphRoot:       rm.config.StorageRoot,
		GraphDriverName: rm.config.StorageDriver,
		GraphDriverOptions: []string{},
	}

	// Apply storage options
	for key, value := range rm.config.StorageOptions {
		storeOptions.GraphDriverOptions = append(storeOptions.GraphDriverOptions, fmt.Sprintf("%s=%s", key, value))
	}

	// Setup UID/GID mappings for rootless
	if rm.rootless {
		// Convert specs-go ID mappings to idtools format
		uidMap := make([]idtools.IDMap, len(rm.config.UIDMap))
		for i, mapping := range rm.config.UIDMap {
			uidMap[i] = idtools.IDMap{
				ContainerID: int(mapping.ContainerID),
				HostID:      int(mapping.HostID),
				Size:        int(mapping.Size),
			}
		}
		gidMap := make([]idtools.IDMap, len(rm.config.GIDMap))
		for i, mapping := range rm.config.GIDMap {
			gidMap[i] = idtools.IDMap{
				ContainerID: int(mapping.ContainerID),
				HostID:      int(mapping.HostID),
				Size:        int(mapping.Size),
			}
		}
		storeOptions.UIDMap = uidMap
		storeOptions.GIDMap = gidMap
	}

	store, err := storage.GetStore(storeOptions)
	if err != nil {
		return errors.Wrap(err, "failed to initialize storage")
	}

	rm.store = store
	return nil
}

// setupRootlessEnvironment configures the environment for rootless operation
func (rm *RuntimeManager) setupRootlessEnvironment() error {
	// Check if user namespaces are available
	if err := rm.checkUserNamespaceSupport(); err != nil {
		return errors.Wrap(err, "user namespaces not supported")
	}

	// Setup UID/GID mappings if not provided
	if len(rm.config.UIDMap) == 0 || len(rm.config.GIDMap) == 0 {
		if err := rm.setupDefaultIDMappings(); err != nil {
			return errors.Wrap(err, "failed to setup default ID mappings")
		}
	}

	// Validate ID mappings
	if err := rm.validateIDMappings(); err != nil {
		return errors.Wrap(err, "invalid ID mappings")
	}

	rm.userNS = true
	return nil
}

// checkUserNamespaceSupport checks if user namespaces are available
func (rm *RuntimeManager) checkUserNamespaceSupport() error {
	if _, err := os.Stat("/proc/self/uid_map"); os.IsNotExist(err) {
		return errors.New("user namespaces not available")
	}
	if err := exec.Command("unshare", "-U", "true").Run(); err != nil {
		return errors.New("cannot create user namespaces")
	}
	return nil
}

// setupDefaultIDMappings sets up default UID/GID mappings for rootless operation
func (rm *RuntimeManager) setupDefaultIDMappings() error {
	uid := os.Getuid()
	gid := os.Getgid()

	// Read subuid and subgid mappings
	subuidMappings, err := readSubIDFile("/etc/subuid", uid)
	if err != nil {
		return errors.Wrap(err, "failed to read subuid mappings")
	}

	subgidMappings, err := readSubIDFile("/etc/subgid", gid)
	if err != nil {
		return errors.Wrap(err, "failed to read subgid mappings")
	}

	// Setup UID mappings
	rm.config.UIDMap = []specs.LinuxIDMapping{
		{
			ContainerID: 0,
			HostID:      uint32(uid),
			Size:        1,
		},
	}
	rm.config.UIDMap = append(rm.config.UIDMap, subuidMappings...)

	// Setup GID mappings
	rm.config.GIDMap = []specs.LinuxIDMapping{
		{
			ContainerID: 0,
			HostID:      uint32(gid),
			Size:        1,
		},
	}
	rm.config.GIDMap = append(rm.config.GIDMap, subgidMappings...)

	return nil
}

// readSubIDFile reads subuid or subgid file and returns ID mappings
func readSubIDFile(filename string, id int) ([]specs.LinuxIDMapping, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read %s", filename)
	}

	var mappings []specs.LinuxIDMapping
	user := os.Getenv("USER")
	for _, line := range strings.Split(string(content), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		fields := strings.Split(line, ":")
		if len(fields) == 3 && fields[0] == user {
			if startID, err := strconv.ParseUint(fields[1], 10, 32); err == nil {
				if count, err := strconv.ParseUint(fields[2], 10, 32); err == nil {
					mappings = append(mappings, specs.LinuxIDMapping{
						ContainerID: 1, HostID: uint32(startID), Size: uint32(count),
					})
				}
			}
		}
	}
	return mappings, nil
}

// validateIDMappings validates the configured UID/GID mappings
func (rm *RuntimeManager) validateIDMappings() error {
	for _, mappings := range [][]specs.LinuxIDMapping{rm.config.UIDMap, rm.config.GIDMap} {
		for _, mapping := range mappings {
			if mapping.Size == 0 {
				return errors.New("ID mapping size cannot be zero")
			}
		}
	}
	return nil
}

// setupSecurityCapabilities configures security capabilities for the runtime
func (rm *RuntimeManager) setupSecurityCapabilities() error {
	if len(rm.config.Capabilities) == 0 {
		if rm.rootless {
			rm.config.Capabilities = []string{"CAP_CHOWN", "CAP_DAC_OVERRIDE", "CAP_FOWNER", "CAP_FSETID", "CAP_KILL", "CAP_NET_BIND_SERVICE", "CAP_SETFCAP", "CAP_SETGID", "CAP_SETPCAP", "CAP_SETUID", "CAP_SYS_CHROOT"}
		} else {
			rm.config.Capabilities = []string{"CAP_AUDIT_WRITE", "CAP_CHOWN", "CAP_DAC_OVERRIDE", "CAP_FOWNER", "CAP_FSETID", "CAP_KILL", "CAP_MKNOD", "CAP_NET_BIND_SERVICE", "CAP_NET_RAW", "CAP_SETFCAP", "CAP_SETGID", "CAP_SETPCAP", "CAP_SETUID", "CAP_SYS_CHROOT"}
		}
	}
	return nil
}

// GetStore returns the initialized storage store
func (rm *RuntimeManager) GetStore() storage.Store {
	return rm.store
}

// IsInitialized returns whether the runtime manager has been initialized
func (rm *RuntimeManager) IsInitialized() bool {
	return rm.initialized
}

// IsRootless returns whether the runtime is running in rootless mode
func (rm *RuntimeManager) IsRootless() bool {
	return rm.rootless
}

// GetConfig returns the runtime configuration
func (rm *RuntimeManager) GetConfig() *BuildConfig {
	return rm.config
}

// CreateBuildOptions creates Buildah build options from the runtime configuration
func (rm *RuntimeManager) CreateBuildOptions() *define.BuildOptions {
	if !rm.initialized {
		return nil
	}
	
	// Parse isolation type
	var isolation define.Isolation
	switch strings.ToLower(rm.config.Isolation) {
	case "chroot":
		isolation = define.IsolationChroot
	case "rootless":
		isolation = define.IsolationOCIRootless
	case "oci":
		isolation = define.IsolationOCI
	default:
		isolation = define.IsolationDefault
	}
	
	options := &define.BuildOptions{
		Isolation: isolation,
	}
	
	// Note: CgroupParent, Mounts, and CapAdd are not directly available in BuildOptions
	// These would typically be handled during the build process or in the build context
	// For this split, we're focusing on the core runtime management functionality
	
	return options
}

// Cleanup performs cleanup operations
func (rm *RuntimeManager) Cleanup() error {
	if rm.store != nil {
		if _, err := rm.store.Shutdown(false); err != nil {
			return errors.Wrap(err, "failed to shutdown storage")
		}
	}

	rm.initialized = false
	return nil
}

// GetNamespaceOptions returns namespace options for container creation
func (rm *RuntimeManager) GetNamespaceOptions() (*define.NamespaceOptions, error) {
	if !rm.initialized {
		return nil, errors.New("runtime manager not initialized")
	}
	nsOptions := &define.NamespaceOptions{}
	for _, ns := range []string{"net", "pid", "ipc", "uts"} {
		nsOptions.AddOrReplace(define.NamespaceOption{Name: ns, Host: false})
	}
	if rm.rootless && rm.userNS {
		nsOptions.AddOrReplace(define.NamespaceOption{Name: "user", Host: false})
	}
	return nsOptions, nil
}

// ValidateRuntimeEnvironment validates that the runtime environment is properly configured
func (rm *RuntimeManager) ValidateRuntimeEnvironment() error {
	if !rm.initialized {
		return errors.New("runtime manager not initialized")
	}

	// Check storage accessibility
	if rm.store == nil {
		return errors.New("storage not initialized")
	}

	// Test storage operations
	if err := rm.testStorageOperations(); err != nil {
		return errors.Wrap(err, "storage validation failed")
	}

	// Validate runtime paths
	if err := rm.validateRuntimePaths(); err != nil {
		return errors.Wrap(err, "runtime paths validation failed")
	}

	return nil
}

// testStorageOperations performs basic storage validation
func (rm *RuntimeManager) testStorageOperations() error {
	if _, err := rm.store.Images(); err != nil {
		return errors.Wrap(err, "failed to list images from storage")
	}
	if _, err := rm.store.Containers(); err != nil {
		return errors.Wrap(err, "failed to list containers from storage")
	}
	return nil
}

// validateRuntimePaths validates configured runtime paths
func (rm *RuntimeManager) validateRuntimePaths() error {
	for path, name := range map[string]string{rm.config.RuntimePath: "runtime", rm.config.ConmonPath: "conmon"} {
		if path != "" {
			if err := rm.validateExecutable(path); err != nil {
				return errors.Wrapf(err, "invalid %s path: %s", name, path)
			}
		}
	}
	return nil
}

// validateExecutable checks if a file exists and is executable
func (rm *RuntimeManager) validateExecutable(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}
	if info.IsDir() {
		return errors.New("path is a directory, not an executable")
	}
	if info.Mode()&0111 == 0 {
		return errors.New("file is not executable")
	}
	return nil
}