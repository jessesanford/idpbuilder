package build

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/containers/buildah/define"
	"github.com/opencontainers/runtime-spec/specs-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestBuildConfig tests the BuildConfig struct and its methods
func TestBuildConfig(t *testing.T) {
	t.Run("ValidConfiguration", func(t *testing.T) {
		config := &BuildConfig{
			StorageDriver:  "overlay2",
			StorageRoot:    "/tmp/test-storage",
			RunRoot:        "/tmp/test-run",
			Isolation:      "oci",
			CgroupManager:  "systemd",
			RuntimePath:    "/usr/bin/runc",
			ConmonPath:     "/usr/bin/conmon",
			UIDMap: []specs.LinuxIDMapping{{ContainerID: 0, HostID: 1000, Size: 1}},
			GIDMap: []specs.LinuxIDMapping{{ContainerID: 0, HostID: 1000, Size: 1}},
			Capabilities: []string{"CAP_CHOWN", "CAP_SETUID"},
			StorageOptions: map[string]string{"overlay2.size": "20G"},
			DefaultMounts: []string{"/proc:/proc:rw", "/sys:/sys:ro"},
		}

		assert.Equal(t, "overlay2", config.StorageDriver)
		assert.Equal(t, "/tmp/test-storage", config.StorageRoot)
		assert.Len(t, config.UIDMap, 1)
		assert.Len(t, config.Capabilities, 2)
		assert.Len(t, config.StorageOptions, 1)
		assert.Len(t, config.DefaultMounts, 2)
	})
}

// TestRuntimeManager tests the RuntimeManager functionality
func TestRuntimeManager(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "runtime-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)
	storageRoot := filepath.Join(tmpDir, "storage")
	runRoot := filepath.Join(tmpDir, "run")

	t.Run("NewRuntimeManager", func(t *testing.T) {
		config := &BuildConfig{StorageDriver: "vfs", StorageRoot: storageRoot, RunRoot: runRoot, Isolation: "oci"}
		rm, err := NewRuntimeManager(config)
		require.NoError(t, err)
		assert.Equal(t, config, rm.GetConfig())
		assert.False(t, rm.IsInitialized())
		assert.Equal(t, os.Geteuid() != 0, rm.IsRootless())
		
		// Test nil config
		rm2, err2 := NewRuntimeManager(nil)
		assert.Error(t, err2)
		assert.Nil(t, rm2)
	})

	t.Run("ValidateConfiguration", func(t *testing.T) {
		// Test invalid configs
		_, err := NewRuntimeManager(&BuildConfig{StorageRoot: storageRoot, RunRoot: runRoot})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "storage driver is required")
		
		_, err = NewRuntimeManager(&BuildConfig{StorageDriver: "vfs", RunRoot: runRoot})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "storage root is required")
		
		_, err = NewRuntimeManager(&BuildConfig{StorageDriver: "vfs", StorageRoot: storageRoot})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "run root is required")
		
		// Test invalid paths
		_, err = NewRuntimeManager(&BuildConfig{StorageDriver: "vfs", StorageRoot: storageRoot, RunRoot: runRoot, RuntimePath: "/nonexistent/runtime"})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "runtime path does not exist")
	})

	t.Run("Initialize", func(t *testing.T) {
		config := &BuildConfig{StorageDriver: "vfs", StorageRoot: storageRoot, RunRoot: runRoot, Isolation: "oci"}
		rm, err := NewRuntimeManager(config)
		require.NoError(t, err)
		
		ctx := context.Background()
		err = rm.Initialize(ctx)
		require.NoError(t, err)
		assert.True(t, rm.IsInitialized())
		assert.NotNil(t, rm.GetStore())
		
		// Test double initialization
		err = rm.Initialize(ctx)
		assert.NoError(t, err)
		
		// Cleanup
		err = rm.Cleanup()
		assert.NoError(t, err)
		assert.False(t, rm.IsInitialized())
	})
}

// TestRuntimeManagerRootless tests rootless-specific functionality
func TestRuntimeManagerRootless(t *testing.T) {
	if os.Geteuid() == 0 {
		t.Skip("Skipping rootless tests when running as root")
	}
	tmpDir, err := os.MkdirTemp("", "rootless-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)
	
	config := &BuildConfig{
		StorageDriver: "vfs",
		StorageRoot: filepath.Join(tmpDir, "storage"),
		RunRoot: filepath.Join(tmpDir, "run"),
		UIDMap: []specs.LinuxIDMapping{{ContainerID: 0, HostID: uint32(os.Getuid()), Size: 1}},
		GIDMap: []specs.LinuxIDMapping{{ContainerID: 0, HostID: uint32(os.Getgid()), Size: 1}},
	}
	
	rm, err := NewRuntimeManager(config)
	require.NoError(t, err)
	assert.True(t, rm.IsRootless())
	
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err = rm.Initialize(ctx); err == nil {
		assert.True(t, rm.IsInitialized())
		nsOptions, err := rm.GetNamespaceOptions()
		assert.NoError(t, err)
		assert.NotNil(t, nsOptions)
		rm.Cleanup()
	}
}

// TestBuildOptions tests build option generation
func TestBuildOptions(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "buildoptions-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)
	
	config := &BuildConfig{
		StorageDriver: "vfs",
		StorageRoot: filepath.Join(tmpDir, "storage"),
		RunRoot: filepath.Join(tmpDir, "run"),
		Isolation: "oci",
		CgroupManager: "systemd",
		DefaultMounts: []string{"/proc:/proc:rw", "/sys:/sys:ro"},
		Capabilities: []string{"CAP_CHOWN", "CAP_SETUID"},
	}
	
	rm, err := NewRuntimeManager(config)
	require.NoError(t, err)
	assert.Nil(t, rm.CreateBuildOptions()) // Before initialization
	
	err = rm.Initialize(context.Background())
	require.NoError(t, err)
	
	options := rm.CreateBuildOptions()
	require.NotNil(t, options)
	assert.Equal(t, define.IsolationOCI, options.Isolation)
	// Note: Mounts and CapAdd are not directly available in BuildOptions
	// They would be handled during the build process
	rm.Cleanup()
}

// TestValidateRuntimeEnvironment tests runtime environment validation
func TestValidateRuntimeEnvironment(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "validate-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)
	
	config := &BuildConfig{StorageDriver: "vfs", StorageRoot: filepath.Join(tmpDir, "storage"), RunRoot: filepath.Join(tmpDir, "run")}
	rm, err := NewRuntimeManager(config)
	require.NoError(t, err)
	
	// Test validation before initialization
	err = rm.ValidateRuntimeEnvironment()
	assert.Error(t, err)
	
	// Initialize and test validation
	err = rm.Initialize(context.Background())
	require.NoError(t, err)
	err = rm.ValidateRuntimeEnvironment()
	assert.NoError(t, err)
	rm.Cleanup()
}

// TestIDMappingValidation tests UID/GID mapping validation
func TestIDMappingValidation(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "idmapping-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)
	
	// Valid mappings
	config := &BuildConfig{
		StorageDriver: "vfs",
		StorageRoot: filepath.Join(tmpDir, "storage"),
		RunRoot: filepath.Join(tmpDir, "run"),
		UIDMap: []specs.LinuxIDMapping{{ContainerID: 0, HostID: 1000, Size: 1}},
		GIDMap: []specs.LinuxIDMapping{{ContainerID: 0, HostID: 1000, Size: 1}},
	}
	rm, err := NewRuntimeManager(config)
	require.NoError(t, err)
	err = rm.Initialize(context.Background())
	require.NoError(t, err)
	rm.Cleanup()
	
	// Invalid mappings
	config.UIDMap = []specs.LinuxIDMapping{{ContainerID: 0, HostID: 1000, Size: 0}}
	rm2, err := NewRuntimeManager(config)
	require.NoError(t, err)
	err = rm2.Initialize(context.Background())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ID mapping size cannot be zero")
}

// TestSecurityCapabilities tests security capability setup
func TestSecurityCapabilities(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "capabilities-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)
	
	// Test custom capabilities
	config := &BuildConfig{
		StorageDriver: "vfs",
		StorageRoot: filepath.Join(tmpDir, "storage"),
		RunRoot: filepath.Join(tmpDir, "run"),
		Capabilities: []string{"CAP_NET_ADMIN", "CAP_SYS_ADMIN"},
	}
	rm, err := NewRuntimeManager(config)
	require.NoError(t, err)
	err = rm.Initialize(context.Background())
	require.NoError(t, err)
	
	capabilities := rm.GetConfig().Capabilities
	assert.Len(t, capabilities, 2)
	assert.Contains(t, capabilities, "CAP_NET_ADMIN")
	rm.Cleanup()
	
	// Test default capabilities (rootless only)
	if os.Geteuid() != 0 {
		config.Capabilities = nil
		rm2, _ := NewRuntimeManager(config)
		if rm2.Initialize(context.Background()) == nil {
			assert.NotEmpty(t, rm2.GetConfig().Capabilities)
			rm2.Cleanup()
		}
	}
}

// TestStorageOperations tests storage-related operations
func TestStorageOperations(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "storage-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	storageRoot := filepath.Join(tmpDir, "storage")
	runRoot := filepath.Join(tmpDir, "run")

	config := &BuildConfig{
		StorageDriver: "vfs",
		StorageRoot:   storageRoot,
		RunRoot:       runRoot,
		StorageOptions: map[string]string{
			"vfs.mountopt": "nodev,nosuid",
		},
	}

	rm, err := NewRuntimeManager(config)
	require.NoError(t, err)

	ctx := context.Background()
	err = rm.Initialize(ctx)
	require.NoError(t, err)

	// Test storage access
	store := rm.GetStore()
	require.NotNil(t, store)

	// Test that we can perform basic storage operations
	images, err := store.Images()
	assert.NoError(t, err)
	assert.NotNil(t, images)

	containers, err := store.Containers()
	assert.NoError(t, err)
	assert.NotNil(t, containers)

	err = rm.Cleanup()
	assert.NoError(t, err)
}

// TestConcurrentOperations tests thread safety of the RuntimeManager
func TestConcurrentOperations(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "concurrent-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	storageRoot := filepath.Join(tmpDir, "storage")
	runRoot := filepath.Join(tmpDir, "run")

	config := &BuildConfig{
		StorageDriver: "vfs",
		StorageRoot:   storageRoot,
		RunRoot:       runRoot,
	}

	rm, err := NewRuntimeManager(config)
	require.NoError(t, err)

	ctx := context.Background()
	err = rm.Initialize(ctx)
	require.NoError(t, err)

	// Test concurrent access to configuration
	done := make(chan bool, 3)

	go func() {
		for i := 0; i < 100; i++ {
			_ = rm.GetConfig()
			_ = rm.IsInitialized()
			_ = rm.IsRootless()
		}
		done <- true
	}()

	go func() {
		for i := 0; i < 100; i++ {
			_ = rm.CreateBuildOptions()
			_, _ = rm.GetNamespaceOptions()
		}
		done <- true
	}()

	go func() {
		for i := 0; i < 100; i++ {
			_ = rm.ValidateRuntimeEnvironment()
		}
		done <- true
	}()

	// Wait for all goroutines to complete
	for i := 0; i < 3; i++ {
		<-done
	}

	err = rm.Cleanup()
	assert.NoError(t, err)
}

// BenchmarkRuntimeManagerInitialization benchmarks runtime manager initialization
func BenchmarkRuntimeManagerInitialization(b *testing.B) {
	tmpDir, err := os.MkdirTemp("", "benchmark-test-*")
	require.NoError(b, err)
	defer os.RemoveAll(tmpDir)

	config := &BuildConfig{
		StorageDriver: "vfs",
		StorageRoot:   filepath.Join(tmpDir, "storage"),
		RunRoot:       filepath.Join(tmpDir, "run"),
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		rm, err := NewRuntimeManager(config)
		require.NoError(b, err)

		ctx := context.Background()
		err = rm.Initialize(ctx)
		require.NoError(b, err)

		err = rm.Cleanup()
		require.NoError(b, err)
	}
}

// TestErrorConditions tests various error conditions
func TestErrorConditions(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "error-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	t.Run("InvalidStorageDriver", func(t *testing.T) {
		config := &BuildConfig{
			StorageDriver: "nonexistent",
			StorageRoot:   filepath.Join(tmpDir, "storage"),
			RunRoot:       filepath.Join(tmpDir, "run"),
		}

		rm, err := NewRuntimeManager(config)
		require.NoError(t, err)

		ctx := context.Background()
		err = rm.Initialize(ctx)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to initialize storage")
	})

	t.Run("ReadOnlyStorageRoot", func(t *testing.T) {
		readOnlyDir := filepath.Join(tmpDir, "readonly")
		err := os.MkdirAll(readOnlyDir, 0555) // Read-only directory
		require.NoError(t, err)

		config := &BuildConfig{
			StorageDriver: "vfs",
			StorageRoot:   filepath.Join(readOnlyDir, "storage"),
			RunRoot:       filepath.Join(tmpDir, "run"),
		}

		rm, err := NewRuntimeManager(config)
		require.NoError(t, err)

		ctx := context.Background()
		err = rm.Initialize(ctx)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to setup directories")
	})
}

// Helper function to create a test executable file
func createTestExecutable(t *testing.T, dir, name string) string {
	t.Helper()
	execPath := filepath.Join(dir, name)
	err := os.WriteFile(execPath, []byte("#!/bin/bash\necho 'test'\n"), 0755)
	require.NoError(t, err)
	return execPath
}

// TestExecutableValidation tests executable path validation
func TestExecutableValidation(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "executable-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	storageRoot := filepath.Join(tmpDir, "storage")
	runRoot := filepath.Join(tmpDir, "run")

	t.Run("ValidExecutablePaths", func(t *testing.T) {
		runtimePath := createTestExecutable(t, tmpDir, "test-runtime")
		conmonPath := createTestExecutable(t, tmpDir, "test-conmon")

		config := &BuildConfig{
			StorageDriver: "vfs",
			StorageRoot:   storageRoot,
			RunRoot:       runRoot,
			RuntimePath:   runtimePath,
			ConmonPath:    conmonPath,
		}

		rm, err := NewRuntimeManager(config)
		require.NoError(t, err)

		ctx := context.Background()
		err = rm.Initialize(ctx)
		assert.NoError(t, err)

		err = rm.Cleanup()
		assert.NoError(t, err)
	})

	t.Run("NonExecutablePaths", func(t *testing.T) {
		// Create non-executable file
		nonExecPath := filepath.Join(tmpDir, "non-exec")
		err := os.WriteFile(nonExecPath, []byte("test"), 0644) // Not executable
		require.NoError(t, err)

		config := &BuildConfig{
			StorageDriver: "vfs",
			StorageRoot:   storageRoot,
			RunRoot:       runRoot,
			RuntimePath:   nonExecPath,
		}

		_, err = NewRuntimeManager(config)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "file is not executable")
	})
}