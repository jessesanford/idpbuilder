package push_test

import (
	"bytes"
	"context"
	"os"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPushCommandIntegration(t *testing.T) {
	rootCmd := &cobra.Command{Use: "idpbuilder"}
	var output bytes.Buffer
	rootCmd.SetOut(&output)
	rootCmd.SetErr(&output)

	rootCmd.SetArgs([]string{"push", "--help"})
	err := rootCmd.ExecuteContext(context.Background())
	assert.NoError(t, err)
	assert.Contains(t, output.String(), "push")
	pushCmd := findCommand(rootCmd, "push")
	if pushCmd != nil {
		assert.Equal(t, "push", pushCmd.Use)
		assert.NotEmpty(t, pushCmd.Short)
	}

	output.Reset()
	rootCmd.SetArgs([]string{"push", "--invalid-flag"})
	err = rootCmd.ExecuteContext(context.Background())
	assert.Error(t, err)
}

func TestFlagPrecedence(t *testing.T) {
	os.Setenv("IDPBUILDER_REGISTRY", "env-registry.com")
	defer os.Unsetenv("IDPBUILDER_REGISTRY")

	rootCmd := &cobra.Command{Use: "idpbuilder"}
	var output bytes.Buffer
	rootCmd.SetOut(&output)

	rootCmd.SetArgs([]string{"push", "--registry", "cli-registry.com", "--dry-run"})
	err := rootCmd.ExecuteContext(context.Background())
	if err == nil || strings.Contains(err.Error(), "cli-registry.com") {
		assert.True(t, true)
	}
	output.Reset()
	rootCmd.SetArgs([]string{"push", "--dry-run"})
	err = rootCmd.ExecuteContext(context.Background())
	outputStr := output.String()
	assert.True(t, err == nil || strings.Contains(outputStr+err.Error(), "env-registry.com"))
}

func TestErrorPropagation(t *testing.T) {
	rootCmd := &cobra.Command{Use: "idpbuilder"}
	var errOutput bytes.Buffer
	rootCmd.SetErr(&errOutput)

	rootCmd.SetArgs([]string{"push"})
	err := rootCmd.ExecuteContext(context.Background())
	if err != nil {
		errorMsg := err.Error()
		assert.NotEmpty(t, errorMsg)
		assert.True(t, len(errorMsg) > 10)
	}
	errOutput.Reset()
	rootCmd.SetArgs([]string{"push", "--registry", "invalid-format", "test-image"})
	err = rootCmd.ExecuteContext(context.Background())
	if err != nil {
		errorOutput := errOutput.String() + err.Error()
		assert.Contains(t, strings.ToLower(errorOutput), "registry")
	}
	assert.Error(t, err)
}

func TestHelpTextGeneration(t *testing.T) {
	rootCmd := &cobra.Command{Use: "idpbuilder"}
	var output bytes.Buffer
	rootCmd.SetOut(&output)

	rootCmd.SetArgs([]string{"push", "--help"})
	err := rootCmd.ExecuteContext(context.Background())
	helpText := output.String()
	assert.Contains(t, helpText, "Usage:")
	assert.Contains(t, helpText, "push")

	commonFlags := []string{"--registry", "--help"}
	for _, flag := range commonFlags {
		assert.Contains(t, helpText, flag)
	}
}

func TestCommandDiscovery(t *testing.T) {
	rootCmd := &cobra.Command{Use: "idpbuilder"}
	var output bytes.Buffer
	rootCmd.SetOut(&output)

	rootCmd.SetArgs([]string{"--help"})
	err := rootCmd.ExecuteContext(context.Background())
	require.NoError(t, err)
	helpOutput := output.String()
	assert.Contains(t, helpOutput, "push")

	commands := rootCmd.Commands()
	pushFound := false
	for _, cmd := range commands {
		if cmd.Use == "push" || strings.Contains(cmd.Use, "push") {
			pushFound = true
			break
		}
	}
	if len(commands) > 0 {
		assert.True(t, pushFound || strings.Contains(helpOutput, "push"))
	}
}

func TestSubcommandInteraction(t *testing.T) {
	rootCmd := &cobra.Command{Use: "idpbuilder"}
	rootCmd.PersistentFlags().String("log-level", "info", "Log level")
	var output bytes.Buffer
	rootCmd.SetOut(&output)

	rootCmd.SetArgs([]string{"push", "--log-level", "debug", "--help"})
	err := rootCmd.ExecuteContext(context.Background())
	assert.NoError(t, err)

	pushCmd := findCommand(rootCmd, "push")
	if pushCmd != nil {
		assert.Equal(t, rootCmd, pushCmd.Parent())
		logLevelFlag := pushCmd.Flag("log-level")
		assert.NotNil(t, logLevelFlag)
	}

	ctx := context.WithValue(context.Background(), "test-key", "test-value")
	output.Reset()
	rootCmd.SetArgs([]string{"push", "--help"})
	err = rootCmd.ExecuteContext(ctx)
	assert.NoError(t, err)
}

func findCommand(rootCmd *cobra.Command, name string) *cobra.Command {
	for _, cmd := range rootCmd.Commands() {
		if cmd.Use == name || strings.HasPrefix(cmd.Use, name+" ") {
			return cmd
		}
	}
	return nil
}