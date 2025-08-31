package push

import (
	"testing"
)

func TestPushCmdFlags(t *testing.T) {
	// Test that the command has the expected flags
	flags := PushCmd.Flags()

	if flags.Lookup("insecure") == nil {
		t.Errorf("expected --insecure flag to be present")
	}

	if flags.Lookup("username") == nil {
		t.Errorf("expected --username flag to be present")
	}

	if flags.Lookup("password") == nil {
		t.Errorf("expected --password flag to be present")
	}

	if flags.Lookup("authfile") == nil {
		t.Errorf("expected --authfile flag to be present")
	}
}

func TestPushCmdMetadata(t *testing.T) {
	if PushCmd.Use != "push [image] [registry]" {
		t.Errorf("expected Use to be 'push [image] [registry]', got %s", PushCmd.Use)
	}

	if PushCmd.Short != "Push container images to registry" {
		t.Errorf("expected Short description to match, got %s", PushCmd.Short)
	}
}
