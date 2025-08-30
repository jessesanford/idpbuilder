package integration

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestEnvironment represents the test environment for integration tests
type TestEnvironment struct {
	ClusterName   string
	RegistryURL   string
	TestNamespace string
	CLIPath       string
	WorkingDir    string
}

// SetupTestEnvironment initializes the test environment for integration tests
func SetupTestEnvironment(t *testing.T) *TestEnvironment {
	t.Helper()

	env := &TestEnvironment{
		ClusterName:   "kind-idp-test",
		RegistryURL:   "gitea.local:443",
		TestNamespace: "integration-test",
		WorkingDir:    t.TempDir(),
	}

	// 1. Verify Kind cluster exists or create one
	env.ensureKindCluster(t)

	// 2. Verify Gitea is running and accessible
	env.verifyGiteaAccess(t)

	// 3. Locate or build CLI binary
	env.locateCLIBinary(t)

	// 4. Create test namespace/workspace
	env.setupTestNamespace(t)

	return env
}

// ensureKindCluster verifies the Kind cluster exists
func (env *TestEnvironment) ensureKindCluster(t *testing.T) {
	t.Helper()

	// Check if cluster exists
	cmd := exec.Command("kind", "get", "clusters")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Skipf("Kind not available: %v", err)
		return
	}

	clusters := strings.Split(strings.TrimSpace(string(output)), "\n")
	for _, cluster := range clusters {
		if strings.Contains(cluster, "idp") {
			env.ClusterName = fmt.Sprintf("kind-%s", cluster)
			t.Logf("Using existing Kind cluster: %s", env.ClusterName)
			return
		}
	}

	t.Skipf("No IDP Kind cluster found. Available clusters: %v", clusters)
}

// verifyGiteaAccess verifies Gitea registry is accessible
func (env *TestEnvironment) verifyGiteaAccess(t *testing.T) {
	t.Helper()

	// Simple connectivity check - try to resolve the hostname
	cmd := exec.Command("nslookup", "gitea.local")
	if err := cmd.Run(); err != nil {
		t.Logf("Warning: Cannot resolve gitea.local, tests may fail: %v", err)
	}
}

// locateCLIBinary finds the CLI binary to test
func (env *TestEnvironment) locateCLIBinary(t *testing.T) {
	t.Helper()

	// Check if binary exists in cli-commands effort directory
	cliCommandsPath := "/home/vscode/workspaces/idpbuilder-oci-mvp/efforts/phase2/wave2/cli-commands/idpbuilder"
	if _, err := os.Stat(cliCommandsPath); err == nil {
		env.CLIPath = cliCommandsPath
		t.Logf("Using CLI binary from cli-commands effort: %s", env.CLIPath)
		return
	}

	// Try to find in current directory
	currentPath := "./idpbuilder"
	if _, err := os.Stat(currentPath); err == nil {
		abs, _ := filepath.Abs(currentPath)
		env.CLIPath = abs
		t.Logf("Using CLI binary from current directory: %s", env.CLIPath)
		return
	}

	// Try to build it
	t.Log("CLI binary not found, attempting to build...")
	buildCmd := exec.Command("go", "build", "-o", "idpbuilder", ".")
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("Failed to build CLI: %v", err)
	}

	abs, _ := filepath.Abs("idpbuilder")
	env.CLIPath = abs
	t.Logf("Built CLI binary: %s", env.CLIPath)
}

// setupTestNamespace creates test namespace if needed
func (env *TestEnvironment) setupTestNamespace(t *testing.T) {
	t.Helper()

	// Create kubectl context for the Kind cluster
	cmd := exec.Command("kubectl", "config", "use-context", env.ClusterName)
	if err := cmd.Run(); err != nil {
		t.Logf("Warning: Failed to set kubectl context: %v", err)
	}

	// Create test namespace
	cmd = exec.Command("kubectl", "create", "namespace", env.TestNamespace, "--dry-run=client", "-o", "yaml")
	output, _ := cmd.Output()
	
	cmd = exec.Command("kubectl", "apply", "-f", "-")
	cmd.Stdin = strings.NewReader(string(output))
	if err := cmd.Run(); err != nil {
		t.Logf("Warning: Failed to create test namespace: %v", err)
	}
}

// Cleanup cleans up test artifacts
func (env *TestEnvironment) Cleanup() {
	// Clean up test artifacts
	if env.TestNamespace != "" {
		cmd := exec.Command("kubectl", "delete", "namespace", env.TestNamespace, "--ignore-not-found=true")
		cmd.Run() // Ignore errors during cleanup
	}

	// Remove test images (if we created any)
	cmd := exec.Command("docker", "images", "--format", "{{.Repository}}:{{.Tag}}", "--filter", "reference=integration-test-*")
	if output, err := cmd.Output(); err == nil {
		images := strings.Split(strings.TrimSpace(string(output)), "\n")
		for _, image := range images {
			if image != "" {
				removeCmd := exec.Command("docker", "rmi", image)
				removeCmd.Run() // Ignore errors during cleanup
			}
		}
	}
}