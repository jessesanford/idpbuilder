// Package push contains tests for the idpbuilder push command.
// This file contains TDD RED phase tests that MUST fail initially.
package push

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestPushCommandRegistration verifies the push command is properly registered
func TestPushCommandRegistration(t *testing.T) {
	require.NotNil(t, pushCmd, "push command should be registered")
	assert.Equal(t, "push", pushCmd.Use, "command should have 'push' as Use")
	assert.NotEmpty(t, pushCmd.Short, "command should have Short description")
	assert.NotEmpty(t, pushCmd.Long, "command should have Long description")
	assert.Equal(t, "push", pushCmd.Name(), "command name should be 'push'")
}

// TestPushCommandFlags verifies all required flags are registered with correct types
func TestPushCommandFlags(t *testing.T) {
	require.NotNil(t, pushCmd, "push command must exist to test flags")
	flags := pushCmd.Flags()
	usernameFlag := flags.Lookup("username")
	require.NotNil(t, usernameFlag, "--username flag should be registered")
	assert.Equal(t, "string", usernameFlag.Value.Type(), "--username should be string type")

	passwordFlag := flags.Lookup("password")
	require.NotNil(t, passwordFlag, "--password flag should be registered")

	namespaceFlag := flags.Lookup("namespace")
	require.NotNil(t, namespaceFlag, "--namespace flag should be registered")

	dirFlag := flags.Lookup("dir")
	require.NotNil(t, dirFlag, "--dir flag should be registered")

	insecureFlag := flags.Lookup("insecure")
	require.NotNil(t, insecureFlag, "--insecure flag should be registered")
	assert.Equal(t, "bool", insecureFlag.Value.Type(), "--insecure should be bool type")

	plainHTTPFlag := flags.Lookup("plain-http")
	require.NotNil(t, plainHTTPFlag, "--plain-http flag should be registered")
}

// TestPushCommandArgValidation verifies argument validation behavior
func TestPushCommandArgValidation(t *testing.T) {
	require.NotNil(t, pushCmd, "push command must exist to test args")

	// Test with no arguments - should fail
	err := pushCmd.Args(pushCmd, []string{})
	assert.Error(t, err, "command should require registry URL argument")

	// Test with valid registry URL
	err = pushCmd.Args(pushCmd, []string{"localhost:5000"})
	assert.NoError(t, err, "command should accept valid registry URL")

	// Test with multiple arguments - should fail
	err = pushCmd.Args(pushCmd, []string{"localhost:5000", "extra"})
	assert.Error(t, err, "command should reject multiple arguments")
}

// TestPushCommandHelp verifies help output contains expected information
func TestPushCommandHelp(t *testing.T) {
	require.NotNil(t, pushCmd, "push command must exist to test help")

	var buf bytes.Buffer
	pushCmd.SetOut(&buf)
	pushCmd.SetArgs([]string{"--help"})
	err := pushCmd.Execute()
	assert.NoError(t, err, "help command should execute without error")

	helpOutput := buf.String()
	assert.Contains(t, helpOutput, "push", "help should contain command name")
	assert.Contains(t, helpOutput, "--username", "help should show username flag")
	assert.Contains(t, helpOutput, "registry", "help should mention registry")
}

// TestPushCommandFlagShorthands verifies shorthand flags work correctly
func TestPushCommandFlagShorthands(t *testing.T) {
	require.NotNil(t, pushCmd, "push command must exist to test shorthands")
	flags := pushCmd.Flags()

	usernameFlag := flags.ShorthandLookup("u")
	require.NotNil(t, usernameFlag, "-u shorthand should map to username")
	assert.Equal(t, "username", usernameFlag.Name, "-u should be shorthand for username")

	passwordFlag := flags.ShorthandLookup("p")
	require.NotNil(t, passwordFlag, "-p shorthand should map to password")

	namespaceFlag := flags.ShorthandLookup("n")
	require.NotNil(t, namespaceFlag, "-n shorthand should map to namespace")

	dirFlag := flags.ShorthandLookup("d")
	require.NotNil(t, dirFlag, "-d shorthand should map to dir")
}

// TestPushCommandEnvVariables verifies environment variable integration
func TestPushCommandEnvVariables(t *testing.T) {
	require.NotNil(t, pushCmd, "push command must exist to test env vars")

	originalVars := []string{
		os.Getenv("IDPBUILDER_USERNAME"),
		os.Getenv("IDPBUILDER_PASSWORD"),
		os.Getenv("IDPBUILDER_NAMESPACE"),
	}
	defer func() {
		os.Setenv("IDPBUILDER_USERNAME", originalVars[0])
		os.Setenv("IDPBUILDER_PASSWORD", originalVars[1])
		os.Setenv("IDPBUILDER_NAMESPACE", originalVars[2])
	}()

	os.Setenv("IDPBUILDER_USERNAME", "testuser")
	os.Setenv("IDPBUILDER_PASSWORD", "testpass")
	os.Setenv("IDPBUILDER_NAMESPACE", "test-namespace")
}

// TestPushCommandDefaults verifies default values are set correctly
func TestPushCommandDefaults(t *testing.T) {
	require.NotNil(t, pushCmd, "push command must exist to test defaults")
	flags := pushCmd.Flags()

	namespaceFlag := flags.Lookup("namespace")
	require.NotNil(t, namespaceFlag, "namespace flag should exist")
	assert.Equal(t, "idpbuilder", namespaceFlag.DefValue, "default namespace should be 'idpbuilder'")

	dirFlag := flags.Lookup("dir")
	require.NotNil(t, dirFlag, "dir flag should exist")
	assert.Equal(t, ".", dirFlag.DefValue, "default dir should be '.'")

	insecureFlag := flags.Lookup("insecure")
	require.NotNil(t, insecureFlag, "insecure flag should exist")
	assert.Equal(t, "false", insecureFlag.DefValue, "default insecure should be false")

	plainHTTPFlag := flags.Lookup("plain-http")
	require.NotNil(t, plainHTTPFlag, "plain-http flag should exist")
	assert.Equal(t, "false", plainHTTPFlag.DefValue, "default plain-http should be false")
}
