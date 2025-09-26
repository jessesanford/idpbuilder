package buildah

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewBuildOptions(t *testing.T) {
	opts := NewBuildOptions()

	assert.NotNil(t, opts)
	assert.NotNil(t, opts.BuildArgs)
	assert.NotNil(t, opts.Labels)
	assert.Empty(t, opts.Env)
	assert.False(t, opts.NoCache)
	assert.False(t, opts.Squash)
}

func TestBuildOptions_WithBuildArg(t *testing.T) {
	opts := NewBuildOptions()

	opts.WithBuildArg("VERSION", "1.0.0").
		WithBuildArg("BUILD_DATE", "2024-01-01")

	assert.Equal(t, "1.0.0", opts.BuildArgs["VERSION"])
	assert.Equal(t, "2024-01-01", opts.BuildArgs["BUILD_DATE"])
}

func TestBuildOptions_WithPlatform(t *testing.T) {
	tests := []struct {
		name     string
		platform string
		wantOS   string
		wantArch string
	}{
		{"linux/amd64", "linux/amd64", "linux", "amd64"},
		{"linux/arm64", "linux/arm64", "linux", "arm64"},
		{"darwin/amd64", "darwin/amd64", "darwin", "amd64"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := NewBuildOptions().WithPlatform(tt.platform)

			assert.Equal(t, tt.platform, opts.Platform)
			assert.Equal(t, tt.wantOS, opts.OS)
			assert.Equal(t, tt.wantArch, opts.Arch)
		})
	}
}

func TestBuildOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		setup   func() *BuildOptions
		wantErr bool
	}{
		{
			name: "valid options",
			setup: func() *BuildOptions {
				return NewBuildOptions().
					WithPlatform("linux/amd64").
					WithEnv("FOO=bar")
			},
			wantErr: false,
		},
		{
			name: "invalid platform",
			setup: func() *BuildOptions {
				return NewBuildOptions().WithPlatform("invalid")
			},
			wantErr: true,
		},
		{
			name: "invalid env var",
			setup: func() *BuildOptions {
				opts := NewBuildOptions()
				opts.Env = []string{"INVALID_FORMAT"}
				return opts
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := tt.setup()
			err := opts.Validate()

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestBuildOptions_ToBuildahArgs(t *testing.T) {
	opts := NewBuildOptions().
		WithBuildArg("VERSION", "2.0.0").
		WithBuildArg("AUTHOR", "test").
		WithPlatform("linux/arm64").
		WithLabel("version", "2.0.0").
		WithLabel("maintainer", "test@example.com")

	opts.NoCache = true
	opts.Squash = true

	args := opts.ToBuildahArgs()

	// Check all expected arguments are present
	assert.Contains(t, args, "--build-arg")
	assert.Contains(t, args, "VERSION=2.0.0")
	assert.Contains(t, args, "AUTHOR=test")
	assert.Contains(t, args, "--platform")
	assert.Contains(t, args, "linux/arm64")
	assert.Contains(t, args, "--label")
	assert.Contains(t, args, "version=2.0.0")
	assert.Contains(t, args, "--no-cache")
	assert.Contains(t, args, "--squash")
}

func TestBuildOptions_ChainedMethods(t *testing.T) {
	opts := NewBuildOptions().
		WithBuildArg("KEY1", "value1").
		WithBuildArg("KEY2", "value2").
		WithEnv("ENV1=val1").
		WithEnv("ENV2=val2").
		WithPlatform("linux/amd64").
		WithLabel("label1", "value1")

	assert.Len(t, opts.BuildArgs, 2)
	assert.Len(t, opts.Env, 2)
	assert.Equal(t, "linux/amd64", opts.Platform)
	assert.Len(t, opts.Labels, 1)
}