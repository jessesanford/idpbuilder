package build

import (
	"testing"
)

func TestBuildCmdFlags(t *testing.T) {
	// Test that the command has the expected flags
	flags := BuildCmd.Flags()

	if flags.Lookup("file") == nil {
		t.Errorf("expected --file flag to be present")
	}

	if flags.Lookup("tag") == nil {
		t.Errorf("expected --tag flag to be present")
	}

	if flags.Lookup("context") == nil {
		t.Errorf("expected --context flag to be present")
	}

	if flags.Lookup("insecure") == nil {
		t.Errorf("expected --insecure flag to be present")
	}
}

func TestBuildCmdMetadata(t *testing.T) {
	if BuildCmd.Use != "build [context]" {
		t.Errorf("expected Use to be 'build [context]', got %s", BuildCmd.Use)
	}

	if BuildCmd.Short != "Build container images using Buildah" {
		t.Errorf("expected Short description to match, got %s", BuildCmd.Short)
	}
}
