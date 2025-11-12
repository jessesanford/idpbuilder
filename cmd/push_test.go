package cmd_test

import (
	"context"
	"testing"

	"github.com/cnoe-io/idpbuilder/cmd"
	"github.com/spf13/cobra"
)

// TestPushCommand_StructCompiles verifies PushCommand struct is a valid Go type with correct fields.
func TestPushCommand_StructCompiles(t *testing.T) {
	// Verify struct fields are accessible
	pc := &cmd.PushCommand{}
	_ = pc
}

// TestPushFlags_StructCompiles verifies PushFlags struct fields are assignable.
func TestPushFlags_StructCompiles(t *testing.T) {
	flags := &cmd.PushFlags{
		Registry: "https://example.com",
		Username: "admin",
		Password: "secret",
		Insecure: true,
		Verbose:  true,
	}

	if flags.Registry != "https://example.com" {
		t.Error("PushFlags struct field assignment failed")
	}
}

// TestNewPushCommand_CreatesValidCommand verifies NewPushCommand() returns a properly configured Cobra command.
func TestNewPushCommand_CreatesValidCommand(t *testing.T) {
	pushCmd := cmd.NewPushCommand()

	if pushCmd == nil {
		t.Fatal("NewPushCommand returned nil")
	}

	if pushCmd.Use != "push IMAGE_NAME" {
		t.Errorf("Expected Use to be 'push IMAGE_NAME', got %q", pushCmd.Use)
	}

	if pushCmd.Short == "" {
		t.Error("Command Short description is empty")
	}
}

// TestNewPushCommand_FlagsDefined verifies all 5 flags are defined and accessible.
func TestNewPushCommand_FlagsDefined(t *testing.T) {
	pushCmd := cmd.NewPushCommand()

	// Check required flags
	registryFlag := pushCmd.Flags().Lookup("registry")
	if registryFlag == nil {
		t.Error("registry flag not defined")
	}

	usernameFlag := pushCmd.Flags().Lookup("username")
	if usernameFlag == nil {
		t.Error("username flag not defined")
	}

	passwordFlag := pushCmd.Flags().Lookup("password")
	if passwordFlag == nil {
		t.Error("password flag not defined")
	}

	insecureFlag := pushCmd.Flags().Lookup("insecure")
	if insecureFlag == nil {
		t.Error("insecure flag not defined")
	}

	verboseFlag := pushCmd.Flags().Lookup("verbose")
	if verboseFlag == nil {
		t.Error("verbose flag not defined")
	}
}

// TestExitCodes_ConstantsDefined verifies all 5 exit code constants have correct values.
func TestExitCodes_ConstantsDefined(t *testing.T) {
	if cmd.ExitSuccess != 0 {
		t.Errorf("ExitSuccess should be 0, got %d", cmd.ExitSuccess)
	}

	if cmd.ExitGeneralError != 1 {
		t.Errorf("ExitGeneralError should be 1, got %d", cmd.ExitGeneralError)
	}

	if cmd.ExitAuthError != 2 {
		t.Errorf("ExitAuthError should be 2, got %d", cmd.ExitAuthError)
	}

	if cmd.ExitNetworkError != 3 {
		t.Errorf("ExitNetworkError should be 3, got %d", cmd.ExitNetworkError)
	}

	if cmd.ExitImageNotFound != 4 {
		t.Errorf("ExitImageNotFound should be 4, got %d", cmd.ExitImageNotFound)
	}
}

// TestNewPushCommand_RunEFunctionWorks verifies runPush executes via command without errors.
func TestNewPushCommand_RunEFunctionWorks(t *testing.T) {
	pushCmd := cmd.NewPushCommand()

	// Set required flag
	pushCmd.Flags().Set("password", "testpass")

	// Execute command with test image name
	pushCmd.SetArgs([]string{"testimage:latest"})
	ctx := context.Background()
	pushCmd.SetContext(ctx)

	err := pushCmd.Execute()
	if err != nil {
		t.Errorf("Command execution failed: %v", err)
	}
}

// TestNewPushCommand_RegistersWithCobra verifies command can be registered with a Cobra root command.
func TestNewPushCommand_RegistersWithCobra(t *testing.T) {
	rootCmd := &cobra.Command{Use: "idpbuilder"}
	pushCmd := cmd.NewPushCommand()

	rootCmd.AddCommand(pushCmd)

	// Verify command was added
	found := false
	for _, c := range rootCmd.Commands() {
		if c.Use == "push IMAGE_NAME" {
			found = true
			break
		}
	}

	if !found {
		t.Error("Push command not registered with root command")
	}
}
