package build

import (
	"testing"

	// Import Phase 1 types for tests
	api "github.com/cnoe-io/idpbuilder/pkg/oci/api"
)

func TestNewBuilder(t *testing.T) {
	config := &api.BuildConfig{RunRoot: "/tmp", GraphRoot: "/tmp"}
	builder, err := NewBuilder(config)
	if err != nil {
		t.Fatal(err)
	}
	if builder.config.MaxParallelBuilds != 3 {
		t.Error("default MaxParallelBuilds should be 3")
	}
}

func TestValidateConfig(t *testing.T) {
	builder := &Builder{}
	err := builder.ValidateConfig(nil)
	if err == nil {
		t.Error("should error on nil config")
	}

	config := &api.BuildConfig{RunRoot: "/tmp", GraphRoot: "/tmp"}
	err = builder.ValidateConfig(config)
	if err != nil {
		t.Errorf("valid config should not error: %v", err)
	}
}
