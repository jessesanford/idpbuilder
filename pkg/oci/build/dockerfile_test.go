package build

import "testing"

func TestNewDockerfileParser(t *testing.T) {
	parser := NewDockerfileParser("Dockerfile", "/tmp")
	if parser.dockerfile != "Dockerfile" {
		t.Error("dockerfile not set correctly")
	}
}

func TestParseFROM(t *testing.T) {
	parser := &DockerfileParser{}
	baseImage, stageName := parser.parseFROM("alpine:latest AS builder")

	if baseImage != "alpine:latest" {
		t.Errorf("expected alpine:latest, got %s", baseImage)
	}
	if stageName != "builder" {
		t.Errorf("expected builder, got %s", stageName)
	}
}
