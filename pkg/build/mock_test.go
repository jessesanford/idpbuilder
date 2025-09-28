package build

import (
	"context"
	"testing"
	"time"
)

func TestMockBuilder(t *testing.T) {
	mock := NewMockBuilder()
	ctx := context.Background()

	// Test successful build workflow
	config := &BuildConfig{
		BaseImage:  "alpine:latest",
		WorkingDir: "/app",
	}

	// Test Build
	buildCtx, err := mock.Build(ctx, config)
	if err != nil {
		t.Errorf("Mock Build() unexpected error: %v", err)
	}

	if buildCtx == nil {
		t.Errorf("Mock Build() returned nil context")
	}

	if buildCtx.ID != "mock-build-123" {
		t.Errorf("Mock Build() ID = %v, want mock-build-123", buildCtx.ID)
	}

	// Test AddLayer
	layer := &LayerSpec{
		Source:      "/src/app",
		Destination: "/app",
		Type:        LayerTypeCopy,
	}

	err = mock.AddLayer(ctx, buildCtx, layer)
	if err != nil {
		t.Errorf("Mock AddLayer() unexpected error: %v", err)
	}

	// Test Finalize
	result, err := mock.Finalize(ctx, buildCtx)
	if err != nil {
		t.Errorf("Mock Finalize() unexpected error: %v", err)
	}

	if result == nil {
		t.Errorf("Mock Finalize() returned nil result")
	}

	if result.ImageID != "mock-image-id-12345" {
		t.Errorf("Mock Finalize() ImageID = %v, want mock-image-id-12345", result.ImageID)
	}

	// Test call tracking
	if mock.GetBuildCallCount() != 1 {
		t.Errorf("Build call count = %d, want 1", mock.GetBuildCallCount())
	}

	if mock.GetAddLayerCallCount() != 1 {
		t.Errorf("AddLayer call count = %d, want 1", mock.GetAddLayerCallCount())
	}

	if mock.GetFinalizeCallCount() != 1 {
		t.Errorf("Finalize call count = %d, want 1", mock.GetFinalizeCallCount())
	}
}

func TestMockBuilderFailures(t *testing.T) {
	mock := NewMockBuilder()
	ctx := context.Background()

	config := &BuildConfig{BaseImage: "alpine:latest"}

	// Test Build failure
	mock.ShouldFailBuild = true
	_, err := mock.Build(ctx, config)
	if err == nil {
		t.Errorf("Expected Build to fail, but it didn't")
	}

	if !IsBuildError(err) {
		t.Errorf("Expected BuildError, got %T", err)
	}

	// Reset and test AddLayer failure
	mock.Reset()
	buildCtx, _ := mock.Build(ctx, config)

	mock.ShouldFailAddLayer = true
	layer := &LayerSpec{Type: LayerTypeCopy}
	err = mock.AddLayer(ctx, buildCtx, layer)
	if err == nil {
		t.Errorf("Expected AddLayer to fail, but it didn't")
	}

	// Reset and test Finalize failure
	mock.Reset()
	buildCtx, _ = mock.Build(ctx, config)

	mock.ShouldFailFinalize = true
	_, err = mock.Finalize(ctx, buildCtx)
	if err == nil {
		t.Errorf("Expected Finalize to fail, but it didn't")
	}
}

func TestMockBuilderConfiguration(t *testing.T) {
	mock := NewMockBuilder()
	ctx := context.Background()

	// Test with custom build context
	customBuildCtx := &BuildContext{
		ID:         "custom-id",
		WorkingDir: "/custom",
		Metadata:   map[string]interface{}{"custom": true},
	}

	mock.WithBuildContext(customBuildCtx)

	config := &BuildConfig{BaseImage: "alpine:latest"}
	result, err := mock.Build(ctx, config)
	if err != nil {
		t.Errorf("Build() unexpected error: %v", err)
	}

	if result.ID != "custom-id" {
		t.Errorf("Build() ID = %v, want custom-id", result.ID)
	}

	// Test with custom build result
	customResult := &BuildResult{
		ImageID:   "custom-image-id",
		Digest:    "custom-digest",
		Size:      2048000,
		CreatedAt: time.Now(),
	}

	mock.WithBuildResult(customResult)

	finalResult, err := mock.Finalize(ctx, result)
	if err != nil {
		t.Errorf("Finalize() unexpected error: %v", err)
	}

	if finalResult.ImageID != "custom-image-id" {
		t.Errorf("Finalize() ImageID = %v, want custom-image-id", finalResult.ImageID)
	}
}

func TestMockBuilderChaining(t *testing.T) {
	mock := NewMockBuilder().
		WithBuildError(ErrBuildFailed).
		WithAddLayerError(ErrLayerAddFailed).
		WithFinalizeError(ErrFinalizeFailed)

	ctx := context.Background()
	config := &BuildConfig{BaseImage: "alpine:latest"}

	// All operations should fail
	_, err := mock.Build(ctx, config)
	if err == nil {
		t.Errorf("Expected Build to fail")
	}

	err = mock.AddLayer(ctx, nil, nil)
	if err == nil {
		t.Errorf("Expected AddLayer to fail")
	}

	_, err = mock.Finalize(ctx, nil)
	if err == nil {
		t.Errorf("Expected Finalize to fail")
	}
}

func TestMockBuilderReset(t *testing.T) {
	mock := NewMockBuilder()
	ctx := context.Background()

	// Make some calls
	config := &BuildConfig{BaseImage: "alpine:latest"}
	buildCtx, _ := mock.Build(ctx, config)
	mock.AddLayer(ctx, buildCtx, &LayerSpec{Type: LayerTypeCopy})
	mock.Finalize(ctx, buildCtx)

	// Verify calls were tracked
	if mock.GetBuildCallCount() != 1 {
		t.Errorf("Build call count before reset = %d, want 1", mock.GetBuildCallCount())
	}

	// Reset
	mock.Reset()

	// Verify calls were cleared
	if mock.GetBuildCallCount() != 0 {
		t.Errorf("Build call count after reset = %d, want 0", mock.GetBuildCallCount())
	}

	if mock.GetAddLayerCallCount() != 0 {
		t.Errorf("AddLayer call count after reset = %d, want 0", mock.GetAddLayerCallCount())
	}

	if mock.GetFinalizeCallCount() != 0 {
		t.Errorf("Finalize call count after reset = %d, want 0", mock.GetFinalizeCallCount())
	}

	// Verify failure flags were reset
	if mock.ShouldFailBuild || mock.ShouldFailAddLayer || mock.ShouldFailFinalize {
		t.Errorf("Failure flags not reset properly")
	}
}