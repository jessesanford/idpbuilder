package contexts

import "testing"

func TestNewArchiveContext(t *testing.T) {
	// Valid archive
	ctx, err := NewArchiveContext("test.tar.gz")
	if err != nil {
		t.Errorf("NewArchiveContext with valid archive failed: %v", err)
	}
	if ctx == nil {
		t.Error("NewArchiveContext returned nil")
	}

	// Invalid archive
	_, err = NewArchiveContext("test.zip")
	if err == nil {
		t.Error("NewArchiveContext should fail with invalid format")
	}
}