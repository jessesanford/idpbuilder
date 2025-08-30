package testutil

import (
	"os/exec"
	"strings"
)

// IsKindClusterRunning checks if a Kind cluster is running
func IsKindClusterRunning(clusterName string) bool {
	cmd := exec.Command("kind", "get", "clusters")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return false
	}
	
	clusters := strings.Split(strings.TrimSpace(string(output)), "\n")
	for _, cluster := range clusters {
		if strings.Contains(cluster, clusterName) || strings.Contains(clusterName, cluster) {
			return true
		}
	}
	return false
}

// GetGiteaURL returns the Gitea registry URL
func GetGiteaURL() string {
	return "gitea.local:443"
}

// IsDockerAvailable checks if Docker is available for image operations
func IsDockerAvailable() bool {
	cmd := exec.Command("docker", "version")
	return cmd.Run() == nil
}