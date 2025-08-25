package api

import (
	"encoding/json"
	"testing"
	"time"
)

func TestBuildConfig_JSONSerialization(t *testing.T) {
	config := &BuildConfig{
		StorageDriver:     "overlay",
		RuntimePath:       "/usr/bin/runc",
		RunRoot:           "/run/containers",
		GraphRoot:         "/var/lib/containers",
		Rootless:          false,
		MaxParallelBuilds: 5,
		BuildTimeout:      5 * time.Minute,
		LogLevel:          "info",
		CacheDir:          "/var/cache/buildah",
		TempDir:           "/tmp/buildah",
	}

	data, err := json.Marshal(config)
	if err != nil {
		t.Errorf("Failed to marshal BuildConfig: %v", err)
	}

	var unmarshaled BuildConfig
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Errorf("Failed to unmarshal BuildConfig: %v", err)
	}

	if unmarshaled.StorageDriver != config.StorageDriver {
		t.Errorf("StorageDriver mismatch: got %s, want %s", unmarshaled.StorageDriver, config.StorageDriver)
	}
	if unmarshaled.MaxParallelBuilds != config.MaxParallelBuilds {
		t.Errorf("MaxParallelBuilds mismatch: got %d, want %d", unmarshaled.MaxParallelBuilds, config.MaxParallelBuilds)
	}
}

func TestRegistryConfig_JSONSerialization(t *testing.T) {
	config := &RegistryConfig{
		URL:        "https://registry.example.com",
		Username:   "testuser",
		Password:   "testpass",
		TLSVerify:  true,
		Insecure:   false,
		Timeout:    30 * time.Second,
		MaxRetries: 3,
		RetryDelay: 5 * time.Second,
		UserAgent:  "idpbuilder/1.0",
	}

	data, err := json.Marshal(config)
	if err != nil {
		t.Errorf("Failed to marshal RegistryConfig: %v", err)
	}

	var unmarshaled RegistryConfig
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Errorf("Failed to unmarshal RegistryConfig: %v", err)
	}

	if unmarshaled.URL != config.URL {
		t.Errorf("URL mismatch: got %s, want %s", unmarshaled.URL, config.URL)
	}
	if unmarshaled.MaxRetries != config.MaxRetries {
		t.Errorf("MaxRetries mismatch: got %d, want %d", unmarshaled.MaxRetries, config.MaxRetries)
	}
}

func TestBuildRequest_Creation(t *testing.T) {
	req := &BuildRequest{
		ID:         "build-123",
		Dockerfile: "/path/to/Dockerfile",
		ContextDir: "/path/to/context",
		Tags:       []string{"myapp:latest", "myapp:v1.0"},
		Platform:   "linux/amd64",
		BuildArgs: map[string]string{
			"VERSION": "1.0",
			"ENV":     "production",
		},
		NoCache:      false,
		Pull:         true,
		SquashLayers: false,
		Created:      time.Now(),
	}

	if req.ID == "" {
		t.Error("BuildRequest ID should not be empty")
	}
	if len(req.Tags) == 0 {
		t.Error("BuildRequest should have at least one tag")
	}
	if req.BuildArgs["VERSION"] != "1.0" {
		t.Error("BuildArgs not properly set")
	}
}

func TestBuildResult_Creation(t *testing.T) {
	layers := []*LayerInfo{
		{
			Digest:     "sha256:abc123",
			Size:       1024,
			MediaType:  "application/vnd.docker.image.rootfs.diff.tar.gzip",
			Created:    time.Now(),
			EmptyLayer: false,
		},
		{
			Digest:     "sha256:def456",
			Size:       2048,
			MediaType:  "application/vnd.docker.image.rootfs.diff.tar.gzip",
			Created:    time.Now(),
			EmptyLayer: false,
		},
	}

	result := &BuildResult{
		BuildID:  "build-123",
		ImageID:  "sha256:image123",
		Digest:   "sha256:result123",
		Tags:     []string{"myapp:latest"},
		Size:     3072,
		Duration: 2 * time.Minute,
		Layers:   layers,
		Warnings: []string{"Layer cache miss"},
		Created:  time.Now(),
	}

	if result.BuildID == "" {
		t.Error("BuildResult BuildID should not be empty")
	}
	if len(result.Layers) != 2 {
		t.Errorf("Expected 2 layers, got %d", len(result.Layers))
	}
	if result.Size != 3072 {
		t.Errorf("Expected size 3072, got %d", result.Size)
	}
}

func TestBuildStatus_StatusValidation(t *testing.T) {
	status := &BuildStatus{
		BuildID:     "build-123",
		Status:      BuildPhaseBuilding,
		Progress:    50,
		CurrentStep: "RUN apt-get update",
		StartTime:   time.Now().Add(-5 * time.Minute),
	}

	if status.Status != BuildPhaseBuilding {
		t.Errorf("Expected status %s, got %s", BuildPhaseBuilding, status.Status)
	}
	if status.Progress < 0 || status.Progress > 100 {
		t.Errorf("Progress should be between 0-100, got %d", status.Progress)
	}
}

func TestBuildPhaseConstants(t *testing.T) {
	phases := []BuildPhase{
		BuildPhaseInitializing,
		BuildPhaseDownloading,
		BuildPhaseBuilding,
		BuildPhaseFinishing,
		BuildPhaseCompleted,
		BuildPhaseFailed,
		BuildPhaseCancelled,
	}

	expected := []string{
		"initializing",
		"downloading",
		"building",
		"finishing",
		"completed",
		"failed",
		"cancelled",
	}

	for i, phase := range phases {
		if string(phase) != expected[i] {
			t.Errorf("Phase constant mismatch: got %s, want %s", string(phase), expected[i])
		}
	}
}

func TestBuildOptions_Creation(t *testing.T) {
	opts := &BuildOptions{
		Quiet:       false,
		NoCache:     true,
		Pull:        true,
		Remove:      true,
		ForceRemove: false,
		Memory:      1024 * 1024 * 1024, // 1GB
		CPUShares:   512,
		CPUQuota:    50000,
		CPUPeriod:   100000,
	}

	if opts.Memory <= 0 {
		t.Error("Memory should be positive")
	}
	if opts.CPUShares <= 0 {
		t.Error("CPUShares should be positive")
	}
}

func TestPushOptions_Creation(t *testing.T) {
	opts := &PushOptions{
		All:                 false,
		Compress:            true,
		DisableContentTrust: false,
		Quiet:               false,
	}

	// Test JSON serialization
	data, err := json.Marshal(opts)
	if err != nil {
		t.Errorf("Failed to marshal PushOptions: %v", err)
	}

	var unmarshaled PushOptions
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Errorf("Failed to unmarshal PushOptions: %v", err)
	}

	if unmarshaled.Compress != opts.Compress {
		t.Errorf("Compress mismatch: got %v, want %v", unmarshaled.Compress, opts.Compress)
	}
}

func TestPullOptions_Creation(t *testing.T) {
	opts := &PullOptions{
		All:                 true,
		DisableContentTrust: false,
		Platform:            "linux/amd64",
		Quiet:               false,
	}

	if opts.Platform == "" {
		t.Error("Platform should be set")
	}
}

func TestLayerInfo_Creation(t *testing.T) {
	layer := &LayerInfo{
		Digest:     "sha256:abc123def456",
		Size:       1024,
		MediaType:  "application/vnd.docker.image.rootfs.diff.tar.gzip",
		Created:    time.Now(),
		CreatedBy:  "RUN apt-get update",
		Comment:    "Install packages",
		EmptyLayer: false,
	}

	if layer.Digest == "" {
		t.Error("LayerInfo Digest should not be empty")
	}
	if layer.Size <= 0 {
		t.Error("LayerInfo Size should be positive")
	}
	if layer.MediaType == "" {
		t.Error("LayerInfo MediaType should not be empty")
	}
}

func TestImageInfo_Creation(t *testing.T) {
	layers := []*LayerInfo{
		{
			Digest:    "sha256:layer1",
			Size:      512,
			MediaType: "application/vnd.docker.image.rootfs.diff.tar.gzip",
			Created:   time.Now(),
		},
	}

	info := &ImageInfo{
		ID:           "sha256:image123",
		Digest:       "sha256:digest123",
		Tags:         []string{"myapp:latest", "myapp:v1.0"},
		Size:         2048,
		Created:      time.Now(),
		Labels:       map[string]string{"version": "1.0", "env": "prod"},
		Architecture: "amd64",
		OS:           "linux",
		Layers:       layers,
	}

	if info.ID == "" {
		t.Error("ImageInfo ID should not be empty")
	}
	if len(info.Tags) == 0 {
		t.Error("ImageInfo should have at least one tag")
	}
	if info.Architecture == "" {
		t.Error("ImageInfo Architecture should not be empty")
	}
	if info.OS == "" {
		t.Error("ImageInfo OS should not be empty")
	}
	if len(info.Layers) != 1 {
		t.Errorf("Expected 1 layer, got %d", len(info.Layers))
	}
}