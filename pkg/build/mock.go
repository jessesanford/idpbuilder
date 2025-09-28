package build

import (
	"context"
	"time"
)

// MockBuilder is a mock implementation of the Builder interface for testing
type MockBuilder struct {
	// Configuration
	ShouldFailBuild    bool
	ShouldFailAddLayer bool
	ShouldFailFinalize bool

	// Tracking
	BuildCalls    []BuildCall
	AddLayerCalls []AddLayerCall
	FinalizeCalls []FinalizeCall

	// Results
	BuildResult   *BuildResult
	BuildContext  *BuildContext
	BuildError    error
	AddLayerError error
	FinalizeError error
}

// BuildCall tracks calls to Build method
type BuildCall struct {
	Context context.Context
	Config  *BuildConfig
}

// AddLayerCall tracks calls to AddLayer method
type AddLayerCall struct {
	Context  context.Context
	BuildCtx *BuildContext
	Layer    *LayerSpec
}

// FinalizeCall tracks calls to Finalize method
type FinalizeCall struct {
	Context  context.Context
	BuildCtx *BuildContext
}

// NewMockBuilder creates a new mock builder
func NewMockBuilder() *MockBuilder {
	return &MockBuilder{
		BuildCalls:    make([]BuildCall, 0),
		AddLayerCalls: make([]AddLayerCall, 0),
		FinalizeCalls: make([]FinalizeCall, 0),
	}
}

// Build implements Builder.Build for testing
func (m *MockBuilder) Build(ctx context.Context, config *BuildConfig) (*BuildContext, error) {
	// Track the call
	m.BuildCalls = append(m.BuildCalls, BuildCall{
		Context: ctx,
		Config:  config,
	})

	// Return error if configured to fail
	if m.ShouldFailBuild {
		if m.BuildError != nil {
			return nil, m.BuildError
		}
		return nil, WrapBuildError("mock_build", ErrBuildFailed, "mock build failure")
	}

	// Return pre-configured result or create default
	if m.BuildContext != nil {
		return m.BuildContext, nil
	}

	// Create default mock build context
	return &BuildContext{
		ID:         "mock-build-123",
		WorkingDir: config.WorkingDir,
		Metadata: map[string]interface{}{
			"mock":       true,
			"base_image": config.BaseImage,
			"created_at": time.Now(),
		},
	}, nil
}

// AddLayer implements Builder.AddLayer for testing
func (m *MockBuilder) AddLayer(ctx context.Context, buildCtx *BuildContext, layer *LayerSpec) error {
	// Track the call
	m.AddLayerCalls = append(m.AddLayerCalls, AddLayerCall{
		Context:  ctx,
		BuildCtx: buildCtx,
		Layer:    layer,
	})

	// Return error if configured to fail
	if m.ShouldFailAddLayer {
		if m.AddLayerError != nil {
			return m.AddLayerError
		}
		return WrapBuildError("mock_add_layer", ErrLayerAddFailed, "mock add layer failure")
	}

	return nil
}

// Finalize implements Builder.Finalize for testing
func (m *MockBuilder) Finalize(ctx context.Context, buildCtx *BuildContext) (*BuildResult, error) {
	// Track the call
	m.FinalizeCalls = append(m.FinalizeCalls, FinalizeCall{
		Context:  ctx,
		BuildCtx: buildCtx,
	})

	// Return error if configured to fail
	if m.ShouldFailFinalize {
		if m.FinalizeError != nil {
			return nil, m.FinalizeError
		}
		return nil, WrapBuildError("mock_finalize", ErrFinalizeFailed, "mock finalize failure")
	}

	// Return pre-configured result or create default
	if m.BuildResult != nil {
		return m.BuildResult, nil
	}

	// Create default mock result
	return &BuildResult{
		ImageID:   "mock-image-id-12345",
		Digest:    "sha256:mock-digest-abcdef",
		Size:      1024000,
		CreatedAt: time.Now(),
	}, nil
}

// Reset clears all tracked calls and resets failure flags
func (m *MockBuilder) Reset() {
	m.BuildCalls = make([]BuildCall, 0)
	m.AddLayerCalls = make([]AddLayerCall, 0)
	m.FinalizeCalls = make([]FinalizeCall, 0)
	m.ShouldFailBuild = false
	m.ShouldFailAddLayer = false
	m.ShouldFailFinalize = false
	m.BuildError = nil
	m.AddLayerError = nil
	m.FinalizeError = nil
	m.BuildResult = nil
	m.BuildContext = nil
}

// AssertBuildCalled checks if Build was called with expected parameters
func (m *MockBuilder) AssertBuildCalled(t interface{ Errorf(format string, args ...interface{}) }, expectedConfig *BuildConfig) {
	if len(m.BuildCalls) == 0 {
		t.Errorf("Expected Build to be called, but it wasn't")
		return
	}

	call := m.BuildCalls[len(m.BuildCalls)-1] // Get the last call
	if call.Config.BaseImage != expectedConfig.BaseImage {
		t.Errorf("Expected Build to be called with BaseImage %s, got %s", expectedConfig.BaseImage, call.Config.BaseImage)
	}
}

// AssertAddLayerCalled checks if AddLayer was called with expected parameters
func (m *MockBuilder) AssertAddLayerCalled(t interface{ Errorf(format string, args ...interface{}) }, expectedLayer *LayerSpec) {
	if len(m.AddLayerCalls) == 0 {
		t.Errorf("Expected AddLayer to be called, but it wasn't")
		return
	}

	call := m.AddLayerCalls[len(m.AddLayerCalls)-1] // Get the last call
	if call.Layer.Type != expectedLayer.Type {
		t.Errorf("Expected AddLayer to be called with Type %s, got %s", expectedLayer.Type, call.Layer.Type)
	}
}

// AssertFinalizeCalled checks if Finalize was called
func (m *MockBuilder) AssertFinalizeCalled(t interface{ Errorf(format string, args ...interface{}) }) {
	if len(m.FinalizeCalls) == 0 {
		t.Errorf("Expected Finalize to be called, but it wasn't")
	}
}

// GetBuildCallCount returns the number of times Build was called
func (m *MockBuilder) GetBuildCallCount() int {
	return len(m.BuildCalls)
}

// GetAddLayerCallCount returns the number of times AddLayer was called
func (m *MockBuilder) GetAddLayerCallCount() int {
	return len(m.AddLayerCalls)
}

// GetFinalizeCallCount returns the number of times Finalize was called
func (m *MockBuilder) GetFinalizeCallCount() int {
	return len(m.FinalizeCalls)
}

// WithBuildError configures the mock to return the specified error on Build
func (m *MockBuilder) WithBuildError(err error) *MockBuilder {
	m.ShouldFailBuild = true
	m.BuildError = err
	return m
}

// WithAddLayerError configures the mock to return the specified error on AddLayer
func (m *MockBuilder) WithAddLayerError(err error) *MockBuilder {
	m.ShouldFailAddLayer = true
	m.AddLayerError = err
	return m
}

// WithFinalizeError configures the mock to return the specified error on Finalize
func (m *MockBuilder) WithFinalizeError(err error) *MockBuilder {
	m.ShouldFailFinalize = true
	m.FinalizeError = err
	return m
}

// WithBuildResult configures the mock to return the specified result on Finalize
func (m *MockBuilder) WithBuildResult(result *BuildResult) *MockBuilder {
	m.BuildResult = result
	return m
}

// WithBuildContext configures the mock to return the specified context on Build
func (m *MockBuilder) WithBuildContext(ctx *BuildContext) *MockBuilder {
	m.BuildContext = ctx
	return m
}

// Verify that MockBuilder implements Builder interface
var _ Builder = (*MockBuilder)(nil)