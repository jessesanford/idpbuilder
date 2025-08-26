package build

import (
	"strings"
	"testing"
)

func TestNewLayerManager(t *testing.T) {
	manager := NewLayerManager("/tmp")
	if manager.cacheDir != "/tmp" {
		t.Error("cache dir not set correctly")
	}
}

func TestCalculateDigest(t *testing.T) {
	manager := NewLayerManager("/tmp")
	digest := manager.CalculateDigest([]byte("test"))

	if !strings.HasPrefix(digest, "sha256:") {
		t.Error("digest should start with sha256:")
	}
}

func TestIsEmptyInstruction(t *testing.T) {
	manager := NewLayerManager("/tmp")

	runInst := &Instruction{Command: "RUN"}
	if manager.isEmptyInstruction(runInst) {
		t.Error("RUN should not be empty instruction")
	}

	envInst := &Instruction{Command: "ENV"}
	if !manager.isEmptyInstruction(envInst) {
		t.Error("ENV should be empty instruction")
	}
}
