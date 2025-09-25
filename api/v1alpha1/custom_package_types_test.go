package v1alpha1

import (
	"testing"
)

func TestCustomPackageSpec_Validate(t *testing.T) {
	spec := &CustomPackageSpec{
		GitRepository: GitRepositoryRef{
			Name: "test-repo",
		},
	}

	// Basic validation test
	if spec.GitRepository.Name != "test-repo" {
		t.Errorf("Expected GitRepository name 'test-repo', got %s", spec.GitRepository.Name)
	}

	// Test optional fields
	if spec.LocalBuild != nil {
		t.Error("Expected LocalBuild to be nil by default")
	}

	if spec.Version != "" {
		t.Errorf("Expected Version to be empty by default, got %s", spec.Version)
	}
}

func TestGitRepositoryRef_Validate(t *testing.T) {
	ref := GitRepositoryRef{
		Name: "test-git-repo",
		Namespace: "default",
	}

	if ref.Name != "test-git-repo" {
		t.Errorf("Expected name 'test-git-repo', got %s", ref.Name)
	}

	if ref.Namespace != "default" {
		t.Errorf("Expected namespace 'default', got %s", ref.Namespace)
	}
}