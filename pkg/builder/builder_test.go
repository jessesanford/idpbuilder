package builder

import (
	"context"
	"reflect"
	"testing"
)

func TestNewBuilder(t *testing.T) {
	builder, err := NewBuilder()
	if err != nil {
		t.Fatalf("NewBuilder() error = %v", err)
	}

	if builder == nil {
		t.Fatal("NewBuilder() returned nil")
	}

	// Check that feature flags are properly initialized
	if builder.IsFeatureSupported(FeatureTarballExport) {
		t.Error("FeatureTarballExport should be disabled in Split 001")
	}

	if builder.IsFeatureSupported(FeatureLayerCaching) {
		t.Error("FeatureLayerCaching should be disabled in Split 001")
	}

	if builder.IsFeatureSupported(FeatureMultiLayer) {
		t.Error("FeatureMultiLayer should be disabled in Split 001")
	}
}

func TestNewBuilder_WithOptions(t *testing.T) {
	factory := NewConfigFactory()
	opts := NewBuildOptions()
	opts.AddLabel("test", "value")

	builder, err := NewBuilder(
		WithConfigFactory(factory),
		WithFeatureFlag(FeatureTarballExport, true),
		WithDefaultBuildOptions(opts),
	)

	if err != nil {
		t.Fatalf("NewBuilder() error = %v", err)
	}

	if !builder.IsFeatureSupported(FeatureTarballExport) {
		t.Error("FeatureTarballExport should be enabled")
	}

	if builder.GetConfig() != factory {
		t.Error("Config factory should be the one we provided")
	}
}

func TestWithConfigFactory(t *testing.T) {
	// Test nil factory
	_, err := NewBuilder(WithConfigFactory(nil))
	if err == nil {
		t.Error("WithConfigFactory(nil) should return error")
	}

	// Test valid factory
	factory := NewConfigFactory()
	builder, err := NewBuilder(WithConfigFactory(factory))
	if err != nil {
		t.Fatalf("WithConfigFactory() error = %v", err)
	}

	if builder.GetConfig() != factory {
		t.Error("Config factory should match the provided one")
	}
}

func TestWithFeatureFlag(t *testing.T) {
	// Test empty feature name
	_, err := NewBuilder(WithFeatureFlag("", true))
	if err == nil {
		t.Error("WithFeatureFlag with empty name should return error")
	}

	// Test valid feature flag
	builder, err := NewBuilder(WithFeatureFlag("test_feature", true))
	if err != nil {
		t.Fatalf("WithFeatureFlag() error = %v", err)
	}

	if !builder.IsFeatureSupported("test_feature") {
		t.Error("test_feature should be enabled")
	}
}

func TestWithDefaultBuildOptions(t *testing.T) {
	// Test nil options
	_, err := NewBuilder(WithDefaultBuildOptions(nil))
	if err == nil {
		t.Error("WithDefaultBuildOptions(nil) should return error")
	}

	// Test invalid options
	invalidOpts := &BuildOptions{Platform: nil} // This will fail validation
	_, err = NewBuilder(WithDefaultBuildOptions(invalidOpts))
	if err == nil {
		t.Error("WithDefaultBuildOptions with invalid options should return error")
	}

	// Test valid options
	opts := NewBuildOptions()
	opts.AddLabel("default.label", "value")
	builder, err := NewBuilder(WithDefaultBuildOptions(opts))
	if err != nil {
		t.Fatalf("WithDefaultBuildOptions() error = %v", err)
	}

	defaults := builder.(*SimpleBuilder).GetDefaultBuildOptions()
	if defaults.Labels["default.label"] != "value" {
		t.Error("Default label should be preserved")
	}
}

func TestSimpleBuilder_Build(t *testing.T) {
	builder, err := NewBuilder()
	if err != nil {
		t.Fatalf("NewBuilder() error = %v", err)
	}

	ctx := context.Background()
	opts := *NewBuildOptions()

	// Test nil context
	_, err = builder.Build(nil, "/tmp", opts)
	if err == nil {
		t.Error("Build with nil context should return error")
	}

	// Test empty context directory
	_, err = builder.Build(ctx, "", opts)
	if err == nil {
		t.Error("Build with empty context directory should return error")
	}

	// Test invalid options
	invalidOpts := BuildOptions{Platform: nil}
	_, err = builder.Build(ctx, "/tmp", invalidOpts)
	if err == nil {
		t.Error("Build with invalid options should return error")
	}

	// Test normal case (should fail because features are disabled)
	_, err = builder.Build(ctx, "/tmp", opts)
	if err == nil {
		t.Error("Build should fail because FeatureTarballExport is disabled")
	}

	// Enable feature and test again (should still fail with stub message)
	builder, err = NewBuilder(WithFeatureFlag(FeatureTarballExport, true))
	if err != nil {
		t.Fatalf("NewBuilder() error = %v", err)
	}

	_, err = builder.Build(ctx, "/tmp", opts)
	if err == nil {
		t.Error("Build should still fail with stub implementation")
	}
}

func TestSimpleBuilder_GetSupportedFeatures(t *testing.T) {
	builder, err := NewBuilder(
		WithFeatureFlag(FeatureTarballExport, true),
		WithFeatureFlag(FeatureLayerCaching, false),
		WithFeatureFlag("custom_feature", true),
	)
	if err != nil {
		t.Fatalf("NewBuilder() error = %v", err)
	}

	features := builder.GetSupportedFeatures()

	expectedCount := 2 // FeatureTarballExport and custom_feature
	if len(features) != expectedCount {
		t.Errorf("Expected %d supported features, got %d", expectedCount, len(features))
	}

	hasFeature := func(features []string, feature string) bool {
		for _, f := range features {
			if f == feature {
				return true
			}
		}
		return false
	}

	if !hasFeature(features, FeatureTarballExport) {
		t.Error("FeatureTarballExport should be in supported features")
	}

	if !hasFeature(features, "custom_feature") {
		t.Error("custom_feature should be in supported features")
	}

	if hasFeature(features, FeatureLayerCaching) {
		t.Error("FeatureLayerCaching should not be in supported features")
	}
}

func TestSimpleBuilder_IsFeatureSupported(t *testing.T) {
	builder, err := NewBuilder(WithFeatureFlag("enabled_feature", true))
	if err != nil {
		t.Fatalf("NewBuilder() error = %v", err)
	}

	if !builder.IsFeatureSupported("enabled_feature") {
		t.Error("enabled_feature should be supported")
	}

	if builder.IsFeatureSupported("disabled_feature") {
		t.Error("disabled_feature should not be supported")
	}

	if builder.IsFeatureSupported(FeatureTarballExport) {
		t.Error("FeatureTarballExport should be disabled by default")
	}
}

func TestSimpleBuilder_ValidateContext(t *testing.T) {
	builder, err := NewBuilder()
	if err != nil {
		t.Fatalf("NewBuilder() error = %v", err)
	}

	simpleBuilder := builder.(*SimpleBuilder)

	// Test empty context directory
	err = simpleBuilder.ValidateContext("")
	if err == nil {
		t.Error("ValidateContext with empty directory should return error")
	}

	// Test valid context directory
	err = simpleBuilder.ValidateContext("/tmp")
	if err != nil {
		t.Errorf("ValidateContext with valid directory should not return error: %v", err)
	}
}

func TestSimpleBuilder_MergeBuildOptions(t *testing.T) {
	// Create builder with default options
	defaultOpts := NewBuildOptions()
	defaultOpts.AddLabel("default.label", "default-value")
	defaultOpts.AddEnvironment("DEFAULT_ENV", "default")
	defaultOpts.WorkingDir = "/default"

	builder, err := NewBuilder(WithDefaultBuildOptions(defaultOpts))
	if err != nil {
		t.Fatalf("NewBuilder() error = %v", err)
	}

	simpleBuilder := builder.(*SimpleBuilder)

	// Create override options
	overrideOpts := BuildOptions{
		WorkingDir: "/override",
		Labels:     map[string]string{"override.label": "override-value"},
		Environment: map[string]string{
			"OVERRIDE_ENV": "override",
			"DEFAULT_ENV":  "overridden", // This should override the default
		},
	}

	merged, err := simpleBuilder.MergeBuildOptions(overrideOpts)
	if err != nil {
		t.Fatalf("MergeBuildOptions() error = %v", err)
	}

	// Check that override values are used
	if merged.WorkingDir != "/override" {
		t.Errorf("Expected WorkingDir '/override', got '%s'", merged.WorkingDir)
	}

	// Check that default labels are preserved
	if merged.Labels["default.label"] != "default-value" {
		t.Error("Default label should be preserved")
	}

	// Check that override labels are applied
	if merged.Labels["override.label"] != "override-value" {
		t.Error("Override label should be applied")
	}

	// Check that environment variables are merged correctly
	if merged.Environment["OVERRIDE_ENV"] != "override" {
		t.Error("Override environment variable should be applied")
	}

	if merged.Environment["DEFAULT_ENV"] != "overridden" {
		t.Error("Default environment variable should be overridden")
	}
}

func TestSimpleBuilder_GetDefaultBuildOptions(t *testing.T) {
	originalOpts := NewBuildOptions()
	originalOpts.AddLabel("test.label", "test-value")
	originalOpts.WorkingDir = "/test"

	builder, err := NewBuilder(WithDefaultBuildOptions(originalOpts))
	if err != nil {
		t.Fatalf("NewBuilder() error = %v", err)
	}

	simpleBuilder := builder.(*SimpleBuilder)
	defaultOpts := simpleBuilder.GetDefaultBuildOptions()

	// Check that values are copied correctly
	if defaultOpts.WorkingDir != "/test" {
		t.Errorf("Expected WorkingDir '/test', got '%s'", defaultOpts.WorkingDir)
	}

	if defaultOpts.Labels["test.label"] != "test-value" {
		t.Error("Test label should be copied")
	}

	// Modify the returned options and ensure original is not affected
	defaultOpts.WorkingDir = "/modified"
	defaultOpts.AddLabel("new.label", "new-value")

	// Get defaults again and check that modifications didn't affect the original
	defaultOpts2 := simpleBuilder.GetDefaultBuildOptions()
	if defaultOpts2.WorkingDir != "/test" {
		t.Error("Modifications should not affect the original default options")
	}

	if _, exists := defaultOpts2.Labels["new.label"]; exists {
		t.Error("New label should not appear in fresh copy of defaults")
	}

	// Verify they are different objects
	if reflect.ValueOf(defaultOpts).Pointer() == reflect.ValueOf(defaultOpts2).Pointer() {
		t.Error("GetDefaultBuildOptions should return different objects")
	}
}