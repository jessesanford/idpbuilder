package build

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/containers/buildah"
	"github.com/containers/storage"
	"github.com/containers/storage/pkg/reexec"
	"github.com/sirupsen/logrus"
)

// StoreManager handles the lifecycle and operations of the Buildah storage backend
type StoreManager struct {
	store           storage.Store
	rootDir         string
	runRoot         string
	graphDriver     string
	storageOptions  storage.StoreOptions
	initialized     bool
	rootless        bool
}

// StoreConfig contains configuration options for initializing the storage backend
type StoreConfig struct {
	RootDir       string            // Root directory for storage
	RunRoot       string            // Runtime directory for storage
	GraphDriver   string            // Storage driver (vfs, overlay, btrfs, zfs)
	GraphOptions  []string          // Driver-specific options  
	StorageOpts   map[string]string // Additional storage options
	Rootless      bool              // Whether running in rootless mode
}

// DefaultStoreConfig returns a default storage configuration
func DefaultStoreConfig() *StoreConfig {
	homeDir, _ := os.UserHomeDir()
	return &StoreConfig{
		RootDir:     filepath.Join(homeDir, ".local", "share", "containers", "storage"),
		RunRoot:     filepath.Join(os.TempDir(), fmt.Sprintf("containers-%d", os.Getuid())),
		GraphDriver: "overlay",
		GraphOptions: []string{
			"overlay.mount_program=/usr/bin/fuse-overlayfs",
		},
		StorageOpts: make(map[string]string),
		Rootless:    os.Getuid() != 0,
	}
}

// NewStoreManager creates a new StoreManager instance
func NewStoreManager(config *StoreConfig) *StoreManager {
	if config == nil {
		config = DefaultStoreConfig()
	}
	
	return &StoreManager{
		rootDir:     config.RootDir,
		runRoot:     config.RunRoot,
		graphDriver: config.GraphDriver,
		rootless:    config.Rootless,
		storageOptions: storage.StoreOptions{
			RunRoot:         config.RunRoot,
			GraphRoot:       config.RootDir,
			GraphDriverName: config.GraphDriver,
			GraphOptions:    config.GraphOptions,
		},
	}
}

// Initialize sets up the storage backend and opens the store
func (sm *StoreManager) Initialize(ctx context.Context) error {
	if sm.initialized {
		return nil
	}

	// Initialize reexec for containers/storage
	if reexec.Init() {
		return nil
	}

	// Create necessary directories
	if err := sm.createDirectories(); err != nil {
		return fmt.Errorf("failed to create storage directories: %w", err)
	}

	// Configure storage options
	sm.configureStorage()

	// Open the storage store
	store, err := storage.GetStore(sm.storageOptions)
	if err != nil {
		return fmt.Errorf("failed to get storage store: %w", err)
	}

	sm.store = store
	sm.initialized = true

	logrus.Infof("Storage backend initialized: driver=%s, rootless=%v", 
		sm.graphDriver, sm.rootless)

	return nil
}

// createDirectories ensures all required directories exist
func (sm *StoreManager) createDirectories() error {
	dirs := []string{sm.rootDir, sm.runRoot}
	
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}
	
	return nil
}

// configureStorage sets up driver-specific storage options
func (sm *StoreManager) configureStorage() {
	switch sm.graphDriver {
	case "overlay":
		if sm.rootless {
			sm.storageOptions.GraphOptions = []string{
				"overlay.mount_program=/usr/bin/fuse-overlayfs",
			}
		}
	case "vfs":
		sm.storageOptions.GraphOptions = []string{}
	default:
		logrus.Warnf("Unknown storage driver %s, using default", sm.graphDriver)
	}
}

// GetStore returns the underlying storage.Store
func (sm *StoreManager) GetStore() storage.Store {
	if !sm.initialized {
		return nil
	}
	return sm.store
}

// GetStoreOptions returns the storage options used for this store
func (sm *StoreManager) GetStoreOptions() storage.StoreOptions {
	return sm.storageOptions
}

// CreateBuilder creates a new buildah.Builder using this store
func (sm *StoreManager) CreateBuilder(ctx context.Context, options buildah.BuilderOptions) (*buildah.Builder, error) {
	if !sm.initialized {
		return nil, fmt.Errorf("store manager not initialized")
	}
	
	// Set the store in builder options
	options.Store = sm.store
	
	return buildah.NewBuilder(ctx, options)
}

// ListImages returns a list of images in the store
func (sm *StoreManager) ListImages() ([]storage.Image, error) {
	if !sm.initialized {
		return nil, fmt.Errorf("store manager not initialized")
	}
	
	return sm.store.Images()
}

// ImageExists checks if an image exists in the store
func (sm *StoreManager) ImageExists(nameOrID string) (bool, error) {
	if !sm.initialized {
		return false, fmt.Errorf("store manager not initialized")
	}
	
	_, err := sm.store.Image(nameOrID)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	
	return true, nil
}

// Shutdown gracefully shuts down the store manager
func (sm *StoreManager) Shutdown() error {
	if !sm.initialized {
		return nil
	}
	
	if sm.store != nil {
		_, err := sm.store.Shutdown(false)
		if err != nil {
			logrus.Errorf("Error shutting down storage store: %v", err)
			return err
		}
	}
	
	sm.initialized = false
	logrus.Info("Storage backend shut down successfully")
	
	return nil
}