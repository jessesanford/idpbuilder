package api

import (
	"strings"
	"testing"
	"time"
)

func TestValidateBuildConfig(t *testing.T) {
	validConfig := &BuildConfig{
		StorageDriver: "overlay", RuntimePath: "/usr/bin/runc",
		RunRoot: "/run/containers", GraphRoot: "/var/lib/containers",
		MaxParallelBuilds: 5, BuildTimeout: 5 * time.Minute, LogLevel: "info",
	}

	if err := ValidateBuildConfig(validConfig); err != nil {
		t.Errorf("Expected valid config to pass, got: %v", err)
	}
	
	if err := ValidateBuildConfig(nil); err == nil {
		t.Error("Expected nil config to fail")
	}

	rootlessOverlay := *validConfig
	rootlessOverlay.Rootless = true
	if err := ValidateBuildConfig(&rootlessOverlay); err == nil {
		t.Error("Expected rootless with overlay to fail")
	}

	rootlessVFS := *validConfig
	rootlessVFS.StorageDriver = "vfs"
	rootlessVFS.Rootless = true
	if err := ValidateBuildConfig(&rootlessVFS); err != nil {
		t.Error("Expected rootless with vfs to pass")
	}

	shortTimeout := *validConfig
	shortTimeout.BuildTimeout = 30 * time.Second
	if err := ValidateBuildConfig(&shortTimeout); err == nil {
		t.Error("Expected short timeout to fail")
	}
}

func TestValidateRegistryConfig(t *testing.T) {
	validConfig := &RegistryConfig{
		URL: "https://registry.example.com", Username: "user", Password: "pass",
		TLSVerify: true, Timeout: 30 * time.Second, MaxRetries: 3, RetryDelay: 5 * time.Second,
	}

	if err := ValidateRegistryConfig(validConfig); err != nil {
		t.Errorf("Expected valid config to pass, got: %v", err)
	}
	
	if err := ValidateRegistryConfig(nil); err == nil {
		t.Error("Expected nil config to fail")
	}

	localhost := &RegistryConfig{
		URL: "http://localhost:5000", TLSVerify: false,
		Timeout: 30 * time.Second, MaxRetries: 3, RetryDelay: 5 * time.Second,
	}
	if err := ValidateRegistryConfig(localhost); err != nil {
		t.Error("Expected localhost without auth to pass")
	}

	noAuth := *validConfig
	noAuth.Username = ""
	noAuth.Password = ""
	if err := ValidateRegistryConfig(&noAuth); err == nil {
		t.Error("Expected remote without auth to fail")
	}

	bothAuth := *validConfig
	bothAuth.Token = "token"
	if err := ValidateRegistryConfig(&bothAuth); err == nil {
		t.Error("Expected both auth types to fail")
	}
}

func TestValidateStackConfig(t *testing.T) {
	validConfig := &StackOCIConfig{
		StackName: "mystack", Version: "1.0.0", Platform: "linux/amd64",
		BaseImage: "ubuntu:20.04", Repository: "mystack", Tags: []string{"latest"},
		Dockerfile: "/Dockerfile", ContextDir: "/context", Created: time.Now(), Updated: time.Now(),
	}

	if err := ValidateStackConfig(validConfig); err != nil {
		t.Errorf("Expected valid config to pass, got: %v", err)
	}
	
	if err := ValidateStackConfig(nil); err == nil {
		t.Error("Expected nil config to fail")
	}

	invalidName := *validConfig
	invalidName.StackName = "My-Stack!"
	if err := ValidateStackConfig(&invalidName); err == nil {
		t.Error("Expected invalid stack name to fail")
	}

	invalidPlatform := *validConfig
	invalidPlatform.Platform = "unsupported/arch"
	if err := ValidateStackConfig(&invalidPlatform); err == nil {
		t.Error("Expected unsupported platform to fail")
	}
}

func TestValidateBuildRequest(t *testing.T) {
	validReq := &BuildRequest{
		ID: "build-12345678", Dockerfile: "/Dockerfile",
		ContextDir: "/context", Tags: []string{"latest"}, Created: time.Now(),
	}

	if err := ValidateBuildRequest(validReq); err != nil {
		t.Errorf("Expected valid request to pass, got: %v", err)
	}
	
	if err := ValidateBuildRequest(nil); err == nil {
		t.Error("Expected nil request to fail")
	}

	invalidDockerfile := *validReq
	invalidDockerfile.Dockerfile = "/somefile.txt"
	if err := ValidateBuildRequest(&invalidDockerfile); err == nil {
		t.Error("Expected invalid dockerfile to fail")
	}

	shortID := *validReq
	shortID.ID = "short"
	if err := ValidateBuildRequest(&shortID); err == nil {
		t.Error("Expected short ID to fail")
	}
}

func TestCustomValidators(t *testing.T) {
	// Test image tags
	validTags := []string{"latest", "v1.0.0", "stable-branch"}
	for _, tag := range validTags {
		if !imageTagRegex.MatchString(tag) {
			t.Errorf("Expected %s to be valid", tag)
		}
	}

	invalidTags := []string{"-invalid", "invalid-", ".invalid", ""}
	for _, tag := range invalidTags {
		if imageTagRegex.MatchString(tag) {
			t.Errorf("Expected %s to be invalid", tag)
		}
	}

	// Test semver
	validVersions := []string{"1.0.0", "v2.1.3", "1.0.0-alpha.1"}
	for _, v := range validVersions {
		if !semverRegex.MatchString(v) {
			t.Errorf("Expected %s to be valid semver", v)
		}
	}

	// Test platforms
	validPlatforms := []string{"linux/amd64", "windows/amd64"}
	for _, p := range validPlatforms {
		if !platformRegex.MatchString(p) {
			t.Errorf("Expected %s to be valid platform", p)
		}
	}
}

func TestHelperFunctions(t *testing.T) {
	// DNS labels
	if !isValidDNSLabel("valid") || !isValidDNSLabel("valid-123") {
		t.Error("Expected valid DNS labels to pass")
	}
	if isValidDNSLabel("") || isValidDNSLabel("-invalid") || isValidDNSLabel(strings.Repeat("a", 64)) {
		t.Error("Expected invalid DNS labels to fail")
	}

	// Repositories
	if !isValidRepository("myrepo") || !isValidRepository("namespace/repo") {
		t.Error("Expected valid repositories to pass")
	}
	if isValidRepository("") || isValidRepository("/invalid") {
		t.Error("Expected invalid repositories to fail")
	}

	// Label keys
	if !isValidLabelKey("valid") || !isValidLabelKey("app.version") {
		t.Error("Expected valid label keys to pass")
	}
	if isValidLabelKey("") || isValidLabelKey("-invalid") {
		t.Error("Expected invalid label keys to fail")
	}

	// Annotation keys
	if !isValidAnnotationKey("example.com") || !isValidAnnotationKey("io.kubernetes.pod-name") {
		t.Error("Expected valid annotation keys to pass")
	}
	if isValidAnnotationKey("") || isValidAnnotationKey("-invalid") {
		t.Error("Expected invalid annotation keys to fail")
	}
}