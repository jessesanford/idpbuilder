package builder

import (
	"testing"
	"time"
)

func TestWithContextPath(t *testing.T) {
	config := DefaultBuildConfig()
	
	// Test valid path
	opt := WithContextPath("/path/to/context")
	if err := opt(config); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	
	if config.ContextPath != "/path/to/context" {
		t.Errorf("expected ContextPath to be '/path/to/context', got %s", config.ContextPath)
	}

	// Test empty path
	opt = WithContextPath("")
	if err := opt(config); err == nil {
		t.Error("expected error for empty context path")
	}
}

func TestWithDockerfile(t *testing.T) {
	config := DefaultBuildConfig()
	
	// Test valid dockerfile
	opt := WithDockerfile("Dockerfile.custom")
	if err := opt(config); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	
	if config.Dockerfile != "Dockerfile.custom" {
		t.Errorf("expected Dockerfile to be 'Dockerfile.custom', got %s", config.Dockerfile)
	}

	// Test empty dockerfile
	opt = WithDockerfile("")
	if err := opt(config); err == nil {
		t.Error("expected error for empty dockerfile")
	}
}

func TestWithTags(t *testing.T) {
	config := DefaultBuildConfig()
	
	// Test valid tags
	opt := WithTags("app:v1.0", "app:latest")
	if err := opt(config); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	
	expectedTags := []string{"app:v1.0", "app:latest"}
	if len(config.Tags) != len(expectedTags) {
		t.Fatalf("expected %d tags, got %d", len(expectedTags), len(config.Tags))
	}
	
	for i, expected := range expectedTags {
		if config.Tags[i] != expected {
			t.Errorf("expected tag %d to be '%s', got '%s'", i, expected, config.Tags[i])
		}
	}

	// Test no tags
	opt = WithTags()
	if err := opt(config); err == nil {
		t.Error("expected error for no tags")
	}

	// Test empty tag
	opt = WithTags("valid:tag", "")
	if err := opt(config); err == nil {
		t.Error("expected error for empty tag")
	}
}

func TestWithPlatform(t *testing.T) {
	config := DefaultBuildConfig()
	
	// Test valid platform
	opt := WithPlatform("linux/arm64")
	if err := opt(config); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	
	if config.Platform.OS != "linux" {
		t.Errorf("expected Platform.OS to be 'linux', got %s", config.Platform.OS)
	}
	
	if config.Platform.Architecture != "arm64" {
		t.Errorf("expected Platform.Architecture to be 'arm64', got %s", config.Platform.Architecture)
	}

	// Test platform with variant
	opt = WithPlatform("linux/arm/v7")
	if err := opt(config); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	
	if config.Platform.Variant != "v7" {
		t.Errorf("expected Platform.Variant to be 'v7', got %s", config.Platform.Variant)
	}

	// Test invalid platform format
	opt = WithPlatform("linux")
	if err := opt(config); err == nil {
		t.Error("expected error for invalid platform format")
	}
}

func TestWithRegistryAuth(t *testing.T) {
	config := DefaultBuildConfig()
	
	// Test valid auth
	opt := WithRegistryAuth("registry.example.com", "user", "pass")
	if err := opt(config); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	
	if config.Registry.Hostname != "registry.example.com" {
		t.Errorf("expected Registry.Hostname to be 'registry.example.com', got %s", config.Registry.Hostname)
	}
	
	if config.Registry.Username != "user" {
		t.Errorf("expected Registry.Username to be 'user', got %s", config.Registry.Username)
	}
	
	if config.Registry.Password != "pass" {
		t.Errorf("expected Registry.Password to be 'pass', got %s", config.Registry.Password)
	}

	// Test empty hostname
	opt = WithRegistryAuth("", "user", "pass")
	if err := opt(config); err == nil {
		t.Error("expected error for empty hostname")
	}

	// Test empty username
	opt = WithRegistryAuth("registry.example.com", "", "pass")
	if err := opt(config); err == nil {
		t.Error("expected error for empty username")
	}

	// Test empty password
	opt = WithRegistryAuth("registry.example.com", "user", "")
	if err := opt(config); err == nil {
		t.Error("expected error for empty password")
	}
}

func TestWithRegistryToken(t *testing.T) {
	config := DefaultBuildConfig()
	
	// Test valid token
	opt := WithRegistryToken("registry.example.com", "token123")
	if err := opt(config); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	
	if config.Registry.Hostname != "registry.example.com" {
		t.Errorf("expected Registry.Hostname to be 'registry.example.com', got %s", config.Registry.Hostname)
	}
	
	if config.Registry.Token != "token123" {
		t.Errorf("expected Registry.Token to be 'token123', got %s", config.Registry.Token)
	}

	// Test empty hostname
	opt = WithRegistryToken("", "token123")
	if err := opt(config); err == nil {
		t.Error("expected error for empty hostname")
	}

	// Test empty token
	opt = WithRegistryToken("registry.example.com", "")
	if err := opt(config); err == nil {
		t.Error("expected error for empty token")
	}
}

func TestWithBuildArg(t *testing.T) {
	config := DefaultBuildConfig()
	
	// Test valid build arg
	opt := WithBuildArg("VERSION", "1.0")
	if err := opt(config); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	
	if config.BuildArgs["VERSION"] != "1.0" {
		t.Errorf("expected BuildArgs['VERSION'] to be '1.0', got %s", config.BuildArgs["VERSION"])
	}

	// Test empty key
	opt = WithBuildArg("", "value")
	if err := opt(config); err == nil {
		t.Error("expected error for empty build arg key")
	}
}

func TestWithBuildArgs(t *testing.T) {
	config := DefaultBuildConfig()
	
	// Test valid build args
	args := map[string]string{
		"VERSION": "1.0",
		"ENV":     "production",
	}
	
	opt := WithBuildArgs(args)
	if err := opt(config); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	
	for key, expectedValue := range args {
		if config.BuildArgs[key] != expectedValue {
			t.Errorf("expected BuildArgs['%s'] to be '%s', got %s", key, expectedValue, config.BuildArgs[key])
		}
	}

	// Test nil args
	opt = WithBuildArgs(nil)
	if err := opt(config); err != nil {
		t.Fatalf("unexpected error for nil args: %v", err)
	}

	// Test empty key
	invalidArgs := map[string]string{
		"":      "value",
		"VALID": "value",
	}
	
	opt = WithBuildArgs(invalidArgs)
	if err := opt(config); err == nil {
		t.Error("expected error for empty build arg key")
	}
}

func TestWithLabel(t *testing.T) {
	config := DefaultBuildConfig()
	
	// Test valid label
	opt := WithLabel("version", "1.0")
	if err := opt(config); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	
	if config.Labels["version"] != "1.0" {
		t.Errorf("expected Labels['version'] to be '1.0', got %s", config.Labels["version"])
	}

	// Test empty key
	opt = WithLabel("", "value")
	if err := opt(config); err == nil {
		t.Error("expected error for empty label key")
	}
}

func TestWithBuildTimeout(t *testing.T) {
	config := DefaultBuildConfig()
	
	// Test valid timeout
	timeout := 45 * time.Minute
	opt := WithBuildTimeout(timeout)
	if err := opt(config); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	
	if config.BuildTimeout != timeout {
		t.Errorf("expected BuildTimeout to be %v, got %v", timeout, config.BuildTimeout)
	}

	// Test zero timeout
	opt = WithBuildTimeout(0)
	if err := opt(config); err == nil {
		t.Error("expected error for zero timeout")
	}

	// Test negative timeout
	opt = WithBuildTimeout(-1 * time.Second)
	if err := opt(config); err == nil {
		t.Error("expected error for negative timeout")
	}
}

func TestWithBuildTimeoutSeconds(t *testing.T) {
	config := DefaultBuildConfig()
	
	// Test valid timeout in seconds
	opt := WithBuildTimeoutSeconds(1800) // 30 minutes
	if err := opt(config); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	
	expected := 30 * time.Minute
	if config.BuildTimeout != expected {
		t.Errorf("expected BuildTimeout to be %v, got %v", expected, config.BuildTimeout)
	}

	// Test zero timeout
	opt = WithBuildTimeoutSeconds(0)
	if err := opt(config); err == nil {
		t.Error("expected error for zero timeout")
	}

	// Test negative timeout
	opt = WithBuildTimeoutSeconds(-1)
	if err := opt(config); err == nil {
		t.Error("expected error for negative timeout")
	}
}

func TestWithMemoryLimit(t *testing.T) {
	config := DefaultBuildConfig()
	
	// Test valid memory limit
	opt := WithMemoryLimit(1024 * 1024 * 1024) // 1GB
	if err := opt(config); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	
	expected := int64(1024 * 1024 * 1024)
	if config.MemoryLimit != expected {
		t.Errorf("expected MemoryLimit to be %d, got %d", expected, config.MemoryLimit)
	}

	// Test negative memory limit
	opt = WithMemoryLimit(-1)
	if err := opt(config); err == nil {
		t.Error("expected error for negative memory limit")
	}
}

func TestWithMemoryLimitMB(t *testing.T) {
	config := DefaultBuildConfig()
	
	// Test valid memory limit in MB
	opt := WithMemoryLimitMB(512) // 512MB
	if err := opt(config); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	
	expected := int64(512 * 1024 * 1024)
	if config.MemoryLimit != expected {
		t.Errorf("expected MemoryLimit to be %d, got %d", expected, config.MemoryLimit)
	}

	// Test negative memory limit
	opt = WithMemoryLimitMB(-1)
	if err := opt(config); err == nil {
		t.Error("expected error for negative memory limit")
	}
}

func TestWithCPULimit(t *testing.T) {
	config := DefaultBuildConfig()
	
	// Test valid CPU limit
	opt := WithCPULimit(1.5)
	if err := opt(config); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	
	if config.CPULimit != 1.5 {
		t.Errorf("expected CPULimit to be 1.5, got %f", config.CPULimit)
	}

	// Test negative CPU limit
	opt = WithCPULimit(-1.0)
	if err := opt(config); err == nil {
		t.Error("expected error for negative CPU limit")
	}
}

func TestWithCPULimitString(t *testing.T) {
	config := DefaultBuildConfig()
	
	// Test valid CPU limit string
	opt := WithCPULimitString("1.5")
	if err := opt(config); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	
	if config.CPULimit != 1.5 {
		t.Errorf("expected CPULimit to be 1.5, got %f", config.CPULimit)
	}

	// Test empty string
	opt = WithCPULimitString("")
	if err := opt(config); err == nil {
		t.Error("expected error for empty CPU limit string")
	}

	// Test invalid format
	opt = WithCPULimitString("invalid")
	if err := opt(config); err == nil {
		t.Error("expected error for invalid CPU limit format")
	}

	// Test negative CPU limit
	opt = WithCPULimitString("-1.0")
	if err := opt(config); err == nil {
		t.Error("expected error for negative CPU limit")
	}
}

func TestNewBuildConfig(t *testing.T) {
	// Test with valid options
	config, err := NewBuildConfig(
		WithContextPath("/app"),
		WithDockerfile("Dockerfile.prod"),
		WithTags("app:v1.0", "app:latest"),
		WithPlatform("linux/arm64"),
	)
	
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	
	if config.ContextPath != "/app" {
		t.Errorf("expected ContextPath to be '/app', got %s", config.ContextPath)
	}
	
	if config.Dockerfile != "Dockerfile.prod" {
		t.Errorf("expected Dockerfile to be 'Dockerfile.prod', got %s", config.Dockerfile)
	}

	// Test with invalid options
	_, err = NewBuildConfig(
		WithContextPath(""),
	)
	
	if err == nil {
		t.Error("expected error for invalid options")
	}
}

func TestMustNewBuildConfig(t *testing.T) {
	// Test with valid options (should not panic)
	config := MustNewBuildConfig(
		WithContextPath("/app"),
		WithTags("app:latest"),
	)
	
	if config.ContextPath != "/app" {
		t.Errorf("expected ContextPath to be '/app', got %s", config.ContextPath)
	}

	// Test with invalid options (should panic)
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for invalid options")
		}
	}()
	
	MustNewBuildConfig(WithContextPath(""))
}

func TestApplyOptions(t *testing.T) {
	config := DefaultBuildConfig()
	
	options := []BuildOption{
		WithContextPath("/app"),
		WithDockerfile("Dockerfile.test"),
		WithTags("test:latest"),
	}
	
	err := ApplyOptions(config, options...)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	
	if config.ContextPath != "/app" {
		t.Errorf("expected ContextPath to be '/app', got %s", config.ContextPath)
	}
	
	if config.Dockerfile != "Dockerfile.test" {
		t.Errorf("expected Dockerfile to be 'Dockerfile.test', got %s", config.Dockerfile)
	}

	// Test with invalid option
	invalidOptions := []BuildOption{
		WithContextPath("/app"),
		WithContextPath(""), // This should fail
	}
	
	err = ApplyOptions(config, invalidOptions...)
	if err == nil {
		t.Error("expected error for invalid option")
	}
}

func TestChainOptions(t *testing.T) {
	config := DefaultBuildConfig()
	
	chainedOption := ChainOptions(
		WithContextPath("/app"),
		WithDockerfile("Dockerfile.test"),
		WithTags("test:latest"),
	)
	
	err := chainedOption(config)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	
	if config.ContextPath != "/app" {
		t.Errorf("expected ContextPath to be '/app', got %s", config.ContextPath)
	}
	
	if config.Dockerfile != "Dockerfile.test" {
		t.Errorf("expected Dockerfile to be 'Dockerfile.test', got %s", config.Dockerfile)
	}
}

func TestConditionalOption(t *testing.T) {
	config := DefaultBuildConfig()
	
	// Test with condition true
	opt := ConditionalOption(true, WithContextPath("/app"))
	err := opt(config)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	
	if config.ContextPath != "/app" {
		t.Errorf("expected ContextPath to be '/app', got %s", config.ContextPath)
	}
	
	// Test with condition false
	opt = ConditionalOption(false, WithContextPath("/other"))
	err = opt(config)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	
	// Should not have changed
	if config.ContextPath != "/app" {
		t.Errorf("expected ContextPath to remain '/app', got %s", config.ContextPath)
	}
}