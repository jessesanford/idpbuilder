package api

import (
	"context"
	"testing"
)

func TestBuilderInterface(t *testing.T) {
	// Test that interface can be implemented
	var builder Builder = &mockBuilder{}

	req := BuildRequest{
		DockerfilePath: "Dockerfile",
		ContextDir:     "/tmp",
		ImageName:      "test",
		ImageTag:       "latest",
	}

	resp, err := builder.BuildAndPush(context.Background(), req)
	if err != nil {
		t.Errorf("interface implementation failed: %v", err)
	}
	if resp == nil {
		t.Error("expected response, got nil")
	}
	if !resp.Success {
		t.Error("expected success")
	}
}

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()
	
	if config.Registry != "gitea.cnoe.localtest.me" {
		t.Errorf("expected registry gitea.cnoe.localtest.me, got %s", config.Registry)
	}
	if config.Namespace != "giteaadmin" {
		t.Errorf("expected namespace giteaadmin, got %s", config.Namespace)
	}
	if !config.InsecureSkipTLSVerify {
		t.Error("expected InsecureSkipTLSVerify to be true")
	}
}

type mockBuilder struct{}

func (m *mockBuilder) BuildAndPush(ctx context.Context, req BuildRequest) (*BuildResponse, error) {
	return &BuildResponse{
		ImageID: "test-image-id",
		FullTag: "gitea.cnoe.localtest.me/giteaadmin/test:latest",
		Success: true,
	}, nil
}