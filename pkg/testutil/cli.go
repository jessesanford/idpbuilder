package testutil

import (
	"os"
	"os/exec"
	"strings"
)

// RunCLICommand runs a CLI command with given arguments
func RunCLICommand(cliPath string, args ...string) (string, error) {
	cmd := exec.Command(cliPath, args...)
	output, err := cmd.CombinedOutput()
	return string(output), err
}

// RunCLICommandWithEnv runs CLI command with custom environment
func RunCLICommandWithEnv(cliPath string, env map[string]string, args ...string) (string, error) {
	cmd := exec.Command(cliPath, args...)
	
	// Set custom environment variables
	cmd.Env = os.Environ()
	for key, value := range env {
		cmd.Env = append(cmd.Env, key+"="+value)
	}
	
	output, err := cmd.CombinedOutput()
	return string(output), err
}

// IsCLIAvailable checks if the CLI binary exists and is executable
func IsCLIAvailable(cliPath string) bool {
	if _, err := os.Stat(cliPath); err != nil {
		return false
	}
	
	// Test basic execution
	cmd := exec.Command(cliPath, "--help")
	err := cmd.Run()
	return err == nil
}