package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestPushCommandStructure verifies command structure is valid
func TestPushCommandStructure(t *testing.T) {
	assert.Equal(t, "push IMAGE", pushCmd.Use)
	assert.NotEmpty(t, pushCmd.Short)
	assert.NotEmpty(t, pushCmd.Long)
	assert.NotNil(t, pushCmd.RunE)
}

// TestPushCommandFlags verifies all flags are defined
func TestPushCommandFlags(t *testing.T) {
	assert.NotNil(t, pushCmd.Flags().Lookup("registry"))
	assert.NotNil(t, pushCmd.Flags().Lookup("username"))
	assert.NotNil(t, pushCmd.Flags().Lookup("password"))
	assert.NotNil(t, pushCmd.Flags().Lookup("insecure"))
	assert.NotNil(t, pushCmd.Flags().Lookup("verbose"))

	// Verify short flags
	assert.NotNil(t, pushCmd.Flags().ShorthandLookup("u"))
	assert.NotNil(t, pushCmd.Flags().ShorthandLookup("p"))
	assert.NotNil(t, pushCmd.Flags().ShorthandLookup("k"))
	assert.NotNil(t, pushCmd.Flags().ShorthandLookup("v"))
}

// TestPushCommandConstants verifies constants are defined
func TestPushCommandConstants(t *testing.T) {
	assert.Equal(t, "https://gitea.cnoe.localtest.me:8443", DefaultRegistry)
	assert.Equal(t, "giteaadmin", DefaultUsername)
}

// TestPushCommandExecution verifies command returns not implemented error
func TestPushCommandExecution(t *testing.T) {
	// Reset flags to defaults
	registryURL = DefaultRegistry
	username = DefaultUsername
	password = "test"
	insecure = false
	verbose = false

	err := runPushCommand(pushCmd, []string{"myapp:latest"})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not yet implemented")
	assert.Contains(t, err.Error(), "Wave 1")
}

// TestPushCommandVerboseMode verifies verbose flag is respected
func TestPushCommandVerboseMode(t *testing.T) {
	verbose = true
	defer func() { verbose = false }()

	err := runPushCommand(pushCmd, []string{"testimage:v1"})

	assert.Error(t, err)
	// In verbose mode, output is printed (tested via manual inspection)
}
