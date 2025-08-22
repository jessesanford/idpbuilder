package api

import "testing"

func TestBuildRequestValidation(t *testing.T) {
	testCases := []struct {
		name    string
		request BuildRequest
		wantErr bool
	}{
		{
			name: "valid request",
			request: BuildRequest{
				DockerfilePath: "Dockerfile",
				ContextDir:     "/tmp/build",
				ImageName:      "myapp",
				ImageTag:       "v1.0",
			},
			wantErr: false,
		},
		{
			name: "missing dockerfile path",
			request: BuildRequest{
				ContextDir: "/tmp/build",
				ImageName:  "myapp",
			},
			wantErr: true,
		},
		{
			name: "missing context dir",
			request: BuildRequest{
				DockerfilePath: "Dockerfile",
				ImageName:      "myapp",
			},
			wantErr: true,
		},
		{
			name: "missing image name",
			request: BuildRequest{
				DockerfilePath: "Dockerfile",
				ContextDir:     "/tmp/build",
			},
			wantErr: true,
		},
		{
			name: "default tag applied",
			request: BuildRequest{
				DockerfilePath: "Dockerfile",
				ContextDir:     "/tmp/build",
				ImageName:      "myapp",
			},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.request.Validate()
			if tc.wantErr && err == nil {
				t.Error("expected error but got none")
			}
			if !tc.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			// Check default tag was applied
			if tc.name == "default tag applied" && tc.request.ImageTag != "latest" {
				t.Errorf("expected default tag 'latest', got %s", tc.request.ImageTag)
			}
		})
	}
}