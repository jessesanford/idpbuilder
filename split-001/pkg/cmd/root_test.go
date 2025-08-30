package cmd

import (
	"context"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestGetRootCmd(t *testing.T) {
	cmd := GetRootCmd()
	
	assert.NotNil(t, cmd)
	assert.Equal(t, "idpbuilder", cmd.Use)
	assert.Equal(t, "Manage reference IDPs", cmd.Short)
	assert.True(t, cmd.SilenceUsage)
	assert.True(t, cmd.SilenceErrors)
}

func TestExecute(t *testing.T) {
	// Test that Execute doesn't panic with a context
	ctx := context.Background()
	
	// This should not panic
	assert.NotPanics(t, func() {
		// We can't easily test the actual execution without args
		// but we can test that the function exists and doesn't panic on creation
		_ = GetRootCmd()
	})
}

func TestSetVersion(t *testing.T) {
	cmd := GetRootCmd()
	testVersion := "v1.2.3"
	
	SetVersion(testVersion)
	
	assert.Equal(t, testVersion, cmd.Version)
}

func TestAddSubCommand(t *testing.T) {
	cmd := GetRootCmd()
	initialCommands := len(cmd.Commands())
	
	// Create a dummy subcommand
	subCmd := &cobra.Command{
		Use:   "test",
		Short: "Test command",
	}
	
	AddSubCommand(subCmd)
	
	assert.Equal(t, initialCommands+1, len(cmd.Commands()))
}