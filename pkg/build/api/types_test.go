package api

import "testing"

func TestBuildRequestValidation(t *testing.T) {
	testCases := []struct {
		name    string
		request BuildRequest
		wantErr bool
	}{
		{"valid request", BuildRequest{"Dockerfile", "/tmp/build", "myapp", "v1.0"}, false},
		{"missing dockerfile path", BuildRequest{"", "/tmp/build", "myapp", ""}, true},
		{"missing context dir", BuildRequest{"Dockerfile", "", "myapp", ""}, true},
		{"missing image name", BuildRequest{"Dockerfile", "/tmp/build", "", ""}, true},
		{"default tag applied", BuildRequest{"Dockerfile", "/tmp/build", "myapp", ""}, false},
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
			if tc.name == "default tag applied" && tc.request.ImageTag != "latest" {
				t.Errorf("expected default tag 'latest', got %s", tc.request.ImageTag)
			}
		})
	}
}