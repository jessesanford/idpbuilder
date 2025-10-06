package integration

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
)

// CleanupManager handles cleanup of test resources
type CleanupManager struct {
	tempDirs      []string
	registryImages []string
	kubeResources []kubeResource
	logFiles      []string
}

type kubeResource struct {
	kind      string
	name      string
	namespace string
}

// NewCleanupManager creates a new cleanup manager
func NewCleanupManager() *CleanupManager {
	return &CleanupManager{
		tempDirs:      make([]string, 0),
		registryImages: make([]string, 0),
		kubeResources: make([]kubeResource, 0),
		logFiles:      make([]string, 0),
	}
}

// AddTempDir registers a temporary directory for cleanup
func (cm *CleanupManager) AddTempDir(dir string) {
	cm.tempDirs = append(cm.tempDirs, dir)
}

// AddRegistryImage registers a registry image for cleanup
func (cm *CleanupManager) AddRegistryImage(image string) {
	cm.registryImages = append(cm.registryImages, image)
}

// AddKubeResource registers a Kubernetes resource for cleanup
func (cm *CleanupManager) AddKubeResource(kind, name, namespace string) {
	cm.kubeResources = append(cm.kubeResources, kubeResource{
		kind:      kind,
		name:      name,
		namespace: namespace,
	})
}

// AddLogFile registers a log file for cleanup
func (cm *CleanupManager) AddLogFile(path string) {
	cm.logFiles = append(cm.logFiles, path)
}

// CleanAll performs all cleanup operations
func (cm *CleanupManager) CleanAll(t *testing.T) error {
	var errors []error

	// Clean registry images first
	if err := cm.cleanRegistryImages(t); err != nil {
		errors = append(errors, fmt.Errorf("registry cleanup failed: %w", err))
	}

	// Clean Kubernetes resources
	if err := cm.cleanKubeResources(t); err != nil {
		errors = append(errors, fmt.Errorf("kubernetes cleanup failed: %w", err))
	}

	// Clean temporary directories
	if err := cm.cleanTempDirs(t); err != nil {
		errors = append(errors, fmt.Errorf("temp dir cleanup failed: %w", err))
	}

	// Clean log files
	if err := cm.cleanLogFiles(t); err != nil {
		errors = append(errors, fmt.Errorf("log file cleanup failed: %w", err))
	}

	if len(errors) > 0 {
		return fmt.Errorf("cleanup had %d errors: %v", len(errors), errors)
	}

	t.Log("All cleanup operations completed successfully")
	return nil
}

// cleanRegistryImages removes test images from the registry
func (cm *CleanupManager) cleanRegistryImages(t *testing.T) error {
	t.Log("Cleaning registry images...")

	for _, image := range cm.registryImages {
		t.Logf("Removing registry image: %s", image)

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		// Use crane to delete the image
		cmd := exec.CommandContext(ctx, "crane", "delete", image)
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Logf("Warning: failed to delete image %s: %v (output: %s)", image, err, string(output))
			continue
		}

		t.Logf("Successfully removed image: %s", image)
	}

	return nil
}

// cleanKubeResources removes Kubernetes test resources
func (cm *CleanupManager) cleanKubeResources(t *testing.T) error {
	t.Log("Cleaning Kubernetes resources...")

	for _, res := range cm.kubeResources {
		t.Logf("Deleting %s/%s in namespace %s", res.kind, res.name, res.namespace)

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		cmd := exec.CommandContext(ctx, "kubectl", "delete", res.kind, res.name, "-n", res.namespace, "--ignore-not-found")
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Logf("Warning: failed to delete %s/%s: %v (output: %s)", res.kind, res.name, err, string(output))
			continue
		}

		t.Logf("Successfully deleted %s/%s", res.kind, res.name)
	}

	return nil
}

// cleanTempDirs removes temporary directories
func (cm *CleanupManager) cleanTempDirs(t *testing.T) error {
	t.Log("Cleaning temporary directories...")

	for _, dir := range cm.tempDirs {
		t.Logf("Removing directory: %s", dir)

		if err := os.RemoveAll(dir); err != nil {
			t.Logf("Warning: failed to remove directory %s: %v", dir, err)
			continue
		}

		t.Logf("Successfully removed directory: %s", dir)
	}

	return nil
}

// cleanLogFiles removes test log files
func (cm *CleanupManager) cleanLogFiles(t *testing.T) error {
	t.Log("Cleaning log files...")

	for _, logFile := range cm.logFiles {
		t.Logf("Removing log file: %s", logFile)

		if err := os.Remove(logFile); err != nil {
			if os.IsNotExist(err) {
				continue
			}
			t.Logf("Warning: failed to remove log file %s: %v", logFile, err)
			continue
		}

		t.Logf("Successfully removed log file: %s", logFile)
	}

	return nil
}

// CleanTestArtifacts removes common test artifacts from the filesystem
func CleanTestArtifacts(t *testing.T, baseDir string) error {
	t.Logf("Cleaning test artifacts in %s", baseDir)

	patterns := []string{
		"*.log",
		"*.tmp",
		"test-output-*",
		".test-*",
	}

	for _, pattern := range patterns {
		matches, err := filepath.Glob(filepath.Join(baseDir, pattern))
		if err != nil {
			t.Logf("Warning: failed to glob pattern %s: %v", pattern, err)
			continue
		}

		for _, match := range matches {
			if err := os.RemoveAll(match); err != nil {
				t.Logf("Warning: failed to remove %s: %v", match, err)
			} else {
				t.Logf("Removed test artifact: %s", match)
			}
		}
	}

	return nil
}

// CleanGiteaRegistry removes all test images from Gitea registry
func CleanGiteaRegistry(t *testing.T, registryURL, username, password string) error {
	t.Logf("Cleaning Gitea registry: %s", registryURL)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	// List all repositories
	cmd := exec.CommandContext(ctx, "crane", "catalog", registryURL)
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("CRANE_USERNAME=%s", username),
		fmt.Sprintf("CRANE_PASSWORD=%s", password),
	)

	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to list registry repositories: %w", err)
	}

	// TODO: Parse output and delete each repository
	t.Logf("Registry catalog output: %s", string(output))

	return nil
}

// WaitForCleanup waits for cleanup operations to complete with a timeout
func WaitForCleanup(t *testing.T, cleanupFn func() error, timeout time.Duration) error {
	done := make(chan error, 1)

	go func() {
		done <- cleanupFn()
	}()

	select {
	case err := <-done:
		return err
	case <-time.After(timeout):
		return fmt.Errorf("cleanup timed out after %v", timeout)
	}
}
