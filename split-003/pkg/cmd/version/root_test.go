package version

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVersionInfo(t *testing.T) {
	info := VersionInfo{
		Version:   "v1.0.0",
		GitCommit: "abc123",
		GitBranch: "main",
		BuildDate: "2023-01-01T00:00:00Z",
		GoVersion: "go1.21.0",
		Platform:  "linux",
		Arch:      "amd64",
	}

	assert.Equal(t, "v1.0.0", info.Version)
	assert.Equal(t, "abc123", info.GitCommit)
	assert.Equal(t, "main", info.GitBranch)
	assert.Equal(t, "2023-01-01T00:00:00Z", info.BuildDate)
	assert.Equal(t, "go1.21.0", info.GoVersion)
	assert.Equal(t, "linux", info.Platform)
	assert.Equal(t, "amd64", info.Arch)
}

func TestVersionCommand(t *testing.T) {
	cmd := VersionCmd

	assert.Equal(t, "version", cmd.Use)
	assert.Equal(t, "Display version information", cmd.Short)
	assert.NotNil(t, cmd.RunE)

	// Test flags exist
	shortFlag := cmd.Flags().Lookup("short")
	assert.NotNil(t, shortFlag)
	assert.Equal(t, "bool", shortFlag.Value.Type())

	outputFlag := cmd.Flags().Lookup("output")
	assert.NotNil(t, outputFlag)
	assert.Equal(t, "string", outputFlag.Value.Type())
}

func TestDisplayVersionDefault(t *testing.T) {
	info := VersionInfo{
		Version:   "v1.0.0",
		GitCommit: "abc123",
		GitBranch: "main",
		BuildDate: "2023-01-01T00:00:00Z",
		GoVersion: "go1.21.0",
		Platform:  "linux",
		Arch:      "amd64",
	}

	// This test just ensures the function doesn't panic
	assert.NotPanics(t, func() {
		err := displayVersionDefault(info)
		assert.NoError(t, err)
	})
}