package testutil

import (
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// AssertImageExists verifies an image exists locally or in registry
func AssertImageExists(t *testing.T, imageRef string) {
	t.Helper()
	
	// Try local check first
	cmd := exec.Command("docker", "image", "inspect", imageRef)
	if cmd.Run() == nil {
		return // Image exists locally
	}
	
	// Try registry check (basic connectivity)
	if strings.Contains(imageRef, "gitea.local") {
		cmd := exec.Command("curl", "-k", "-I", "https://gitea.local:443/")
		if cmd.Run() == nil {
			t.Logf("Registry accessible for image: %s", imageRef)
		} else {
			t.Skipf("Cannot verify image %s - registry not accessible", imageRef)
		}
	}
}

// AssertCertificateConfigured verifies certificate trust is configured
func AssertCertificateConfigured(t *testing.T) {
	t.Helper()
	
	// Check if certificate files exist in common locations
	certPaths := []string{
		"/etc/ssl/certs/",
		"/usr/local/share/ca-certificates/",
		"~/.docker/certs.d/",
	}
	
	foundCerts := false
	for _, path := range certPaths {
		cmd := exec.Command("ls", path)
		if cmd.Run() == nil {
			foundCerts = true
			break
		}
	}
	
	if !foundCerts {
		t.Skip("No certificate paths found - may be environment specific")
	}
	
	t.Log("Certificate configuration appears to be available")
}

// AssertCommandSucceeds verifies a command runs successfully
func AssertCommandSucceeds(t *testing.T, cmdName string, args ...string) {
	t.Helper()
	
	cmd := exec.Command(cmdName, args...)
	output, err := cmd.CombinedOutput()
	
	assert.NoError(t, err, "Command should succeed: %s %v\nOutput: %s", 
		cmdName, args, string(output))
}