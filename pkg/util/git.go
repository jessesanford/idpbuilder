// pkg/util/git.go
package util

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/cnoe-io/idpbuilder/api/v1alpha1"
	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/storage/memory"
)

// GitGetRemoteURL returns the remote URL for origin
func GitGetRemoteURL() (string, error) {
	cmd := exec.Command("git", "remote", "get-url", "origin")
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

// GitGetCurrentBranch returns current branch name
func GitGetCurrentBranch() (string, error) {
	cmd := exec.Command("git", "branch", "--show-current")
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

// CloneRemoteRepoToDir clones a remote repository to a directory
func CloneRemoteRepoToDir(ctx context.Context, spec v1alpha1.RemoteRepositorySpec, depth int, progress bool, dir string, username string) (*git.Repository, string, error) {
	cloneOptions := &git.CloneOptions{
		URL: spec.Url,
	}

	if depth > 0 {
		cloneOptions.Depth = depth
		cloneOptions.ShallowSubmodules = true
	}

	if spec.Ref != "" {
		cloneOptions.ReferenceName = plumbing.ReferenceName(fmt.Sprintf("refs/tags/%s", spec.Ref))
		cloneOptions.SingleBranch = true
	}

	if progress {
		cloneOptions.Progress = os.Stdout
	}

	// Check if directory exists and has .git
	if _, err := os.Stat(filepath.Join(dir, ".git")); err == nil {
		// Repository already exists, open it and potentially update
		repo, err := git.PlainOpen(dir)
		if err != nil {
			return nil, "", err
		}

		// If we have a specific ref, check it out
		if spec.Ref != "" {
			w, err := repo.Worktree()
			if err != nil {
				return nil, "", err
			}

			// Try to fetch the ref first
			err = repo.FetchContext(ctx, &git.FetchOptions{
				RefSpecs: []config.RefSpec{"refs/tags/*:refs/tags/*"},
			})
			if err != nil && err != git.NoErrAlreadyUpToDate {
				// Ignore fetch errors, might be tags already exist
			}

			// Checkout the specific tag
			err = w.Checkout(&git.CheckoutOptions{
				Branch: plumbing.ReferenceName(fmt.Sprintf("refs/tags/%s", spec.Ref)),
			})
			if err != nil {
				return nil, "", err
			}
		}

		return repo, dir, nil
	}

	// Create directory if it doesn't exist
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, "", err
	}

	repo, err := git.PlainCloneContext(ctx, dir, false, cloneOptions)
	if err != nil {
		return nil, "", err
	}

	return repo, dir, nil
}

// CloneRemoteRepoToMemory clones a remote repository to memory
func CloneRemoteRepoToMemory(ctx context.Context, spec v1alpha1.RemoteRepositorySpec, depth int, progress bool) (billy.Filesystem, string, error) {
	cloneOptions := &git.CloneOptions{
		URL: spec.Url,
	}

	if depth > 0 {
		cloneOptions.Depth = depth
		cloneOptions.ShallowSubmodules = true
	}

	if spec.Ref != "" {
		cloneOptions.ReferenceName = plumbing.ReferenceName(fmt.Sprintf("refs/tags/%s", spec.Ref))
		cloneOptions.SingleBranch = true
	}

	if progress {
		cloneOptions.Progress = os.Stdout
	}

	fs := memfs.New()
	repo, err := git.CloneContext(ctx, memory.NewStorage(), fs, cloneOptions)
	if err != nil {
		return nil, "", err
	}

	head, err := repo.Head()
	if err != nil {
		return nil, "", err
	}

	return fs, head.Hash().String(), nil
}

// CopyTreeToTree copies files from source filesystem to destination filesystem
func CopyTreeToTree(src, dst billy.Filesystem, srcPath, dstPath string) error {
	files, err := src.ReadDir(srcPath)
	if err != nil {
		return err
	}

	// Create destination directory
	if err := dst.MkdirAll(dstPath, 0755); err != nil {
		return err
	}

	for _, file := range files {
		srcFilePath := filepath.Join(srcPath, file.Name())
		dstFilePath := filepath.Join(dstPath, file.Name())

		if file.IsDir() {
			// Recursively copy subdirectory
			if err := CopyTreeToTree(src, dst, srcFilePath, dstFilePath); err != nil {
				return err
			}
		} else {
			// Copy file
			srcFile, err := src.Open(srcFilePath)
			if err != nil {
				return err
			}
			defer srcFile.Close()

			dstFile, err := dst.Create(dstFilePath)
			if err != nil {
				return err
			}
			defer dstFile.Close()

			_, err = io.Copy(dstFile, srcFile)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// ReadWorktreeFile reads a file from a worktree filesystem
func ReadWorktreeFile(fs billy.Filesystem, path string) ([]byte, error) {
	file, err := fs.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return io.ReadAll(file)
}

// GetWorktreeYamlFiles gets YAML files from a worktree filesystem
func GetWorktreeYamlFiles(startPath string, fs billy.Filesystem, recursive bool) ([]string, error) {
	var yamlFiles []string

	err := walkFilesystem(fs, startPath, recursive, func(path string, info os.FileInfo) error {
		if !info.IsDir() && (strings.HasSuffix(path, ".yaml") || strings.HasSuffix(path, ".yml")) {
			yamlFiles = append(yamlFiles, path)
		}
		return nil
	})

	return yamlFiles, err
}

// walkFilesystem walks through a billy filesystem
func walkFilesystem(fs billy.Filesystem, startPath string, recursive bool, fn func(string, os.FileInfo) error) error {
	files, err := fs.ReadDir(startPath)
	if err != nil {
		return err
	}

	for _, file := range files {
		path := filepath.Join(startPath, file.Name())

		if err := fn(path, file); err != nil {
			return err
		}

		if file.IsDir() && recursive {
			if err := walkFilesystem(fs, path, recursive, fn); err != nil {
				return err
			}
		}
	}

	return nil
}
