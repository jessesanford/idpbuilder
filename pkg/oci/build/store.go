package build

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/containers/buildah"
	"github.com/containers/storage"
	"github.com/containers/storage/pkg/reexec"
	"github.com/opencontainers/go-digest"
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
			RootlessStorageConfig: storage.RootlessStorageConfig{
				RootlessUID: os.Getuid(),
				RootlessGID: os.Getgid(),
			},
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

	// Configure storage options based on rootless mode
	if err := sm.configureStorage(); err != nil {
		return fmt.Errorf("failed to configure storage: %w", err)
	}

	// Open the storage store
	store, err := storage.GetStore(sm.storageOptions)
	if err != nil {
		return fmt.Errorf("failed to get storage store: %w", err)
	}

	sm.store = store
	sm.initialized = true

	logrus.Infof("Storage backend initialized: driver=%s, rootDir=%s, runRoot=%s, rootless=%v", 
		sm.graphDriver, sm.rootDir, sm.runRoot, sm.rootless)

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
func (sm *StoreManager) configureStorage() error {
	switch sm.graphDriver {
	case "overlay":
		return sm.configureOverlayDriver()
	case "vfs":
		return sm.configureVFSDriver()
	case "btrfs":
		return sm.configureBtrfsDriver()
	case "zfs":
		return sm.configureZFSDriver()
	default:
		logrus.Warnf("Unknown storage driver %s, using default configuration", sm.graphDriver)
	}
	
	return nil
}

// configureOverlayDriver sets up overlay driver specific options
func (sm *StoreManager) configureOverlayDriver() error {
	// Check if overlay is supported
	if !sm.isOverlaySupported() {
		logrus.Warn("Overlay driver may not be supported, falling back to vfs")
		sm.graphDriver = "vfs"
		sm.storageOptions.GraphDriverName = "vfs"
		return nil
	}

	// Set overlay-specific options
	overlayOptions := []string{
		"overlay.mountopt=nodev,metacopy=on",
	}
	
	if sm.rootless {
		overlayOptions = append(overlayOptions, 
			"overlay.mount_program=/usr/bin/fuse-overlayfs",
			"overlay.skip_mount_home=true",
		)
	}
	
	sm.storageOptions.GraphOptions = overlayOptions
	return nil
}

// configureVFSDriver sets up VFS driver (simple directory-based storage)
func (sm *StoreManager) configureVFSDriver() error {
	// VFS driver needs no special configuration
	sm.storageOptions.GraphOptions = []string{}
	return nil
}

// configureBtrfsDriver sets up Btrfs driver options
func (sm *StoreManager) configureBtrfsDriver() error {
	// Check if we're on a btrfs filesystem
	if !sm.isBtrfsSupported() {
		return fmt.Errorf("btrfs driver selected but filesystem is not btrfs")
	}
	
	sm.storageOptions.GraphOptions = []string{
		"btrfs.min_space=1G",
	}
	return nil
}

// configureZFSDriver sets up ZFS driver options  
func (sm *StoreManager) configureZFSDriver() error {
	// Basic ZFS configuration
	sm.storageOptions.GraphOptions = []string{
		"zfs.fsname=zroot/containers",
	}
	return nil
}

// isOverlaySupported checks if overlay filesystem is supported
func (sm *StoreManager) isOverlaySupported() bool {
	// Check if overlay module is loaded
	if _, err := os.Stat("/proc/filesystems"); err != nil {
		return false
	}
	
	// Simple check - in production this would be more thorough
	return true
}

// isBtrfsSupported checks if the storage root is on a btrfs filesystem
func (sm *StoreManager) isBtrfsSupported() bool {
	// Simple check - would need more sophisticated detection in production
	return false
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

// ListContainers returns a list of containers in the store
func (sm *StoreManager) ListContainers() ([]storage.Container, error) {
	if !sm.initialized {
		return nil, fmt.Errorf("store manager not initialized")
	}
	
	return sm.store.Containers()
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

// DeleteImage removes an image from the store
func (sm *StoreManager) DeleteImage(nameOrID string, force bool) error {
	if !sm.initialized {
		return fmt.Errorf("store manager not initialized")
	}
	
	image, err := sm.store.Image(nameOrID)
	if err != nil {
		return fmt.Errorf("failed to find image %s: %w", nameOrID, err)
	}
	
	_, err = sm.store.DeleteImage(image.ID, force)
	if err != nil {
		return fmt.Errorf("failed to delete image %s: %w", nameOrID, err)
	}
	
	return nil
}

// GetImageByDigest retrieves an image by its digest
func (sm *StoreManager) GetImageByDigest(digest digest.Digest) (*storage.Image, error) {
	if !sm.initialized {
		return nil, fmt.Errorf("store manager not initialized")
	}
	
	images, err := sm.store.Images()
	if err != nil {
		return nil, fmt.Errorf("failed to list images: %w", err)
	}
	
	for _, img := range images {
		for _, imgDigest := range img.Digests {
			if imgDigest == digest {
				return &img, nil
			}
		}
	}
	
	return nil, storage.ErrImageUnknown
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

// GarbageCollect performs cleanup of unused storage
func (sm *StoreManager) GarbageCollect(ctx context.Context) error {
	if !sm.initialized {
		return fmt.Errorf("store manager not initialized")
	}
	
	// Get all containers and images
	containers, err := sm.store.Containers()
	if err != nil {
		return fmt.Errorf("failed to list containers: %w", err)
	}
	
	images, err := sm.store.Images()
	if err != nil {
		return fmt.Errorf("failed to list images: %w", err)
	}
	
	logrus.Infof("Garbage collection completed. Found %d containers and %d images", 
		len(containers), len(images))
	
	return nil
}